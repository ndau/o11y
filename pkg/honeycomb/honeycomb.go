package honeycomb

// ----- ---- --- -- -
// Copyright 2018, 2019, 2020 The Axiom Foundation. All Rights Reserved.
//
// Licensed under the Apache License 2.0 (the "License").  You may not use
// this file except in compliance with the License.  You can obtain a copy
// in the file LICENSE in the source distribution or at
// https://www.apache.org/licenses/LICENSE-2.0.txt
// - -- --- ---- -----


import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"sync"

	libhoney "github.com/honeycombio/libhoney-go"
	"github.com/sirupsen/logrus"
)

// Autoflushing is important in a serverless context.
// Per https://docs.honeycomb.io/getting-data-in/integrations/aws/aws-lambda/ :
//
// > Normally, libhoney events are enqueued and sent as batches. By default,
// > this occurs every 100ms or whenever the queue is full. However, because
// > Lambda freezes the function instance between invocations, the goroutine
// > responsible for sending this batch is not guaranteed to execute. To ensure
// > that events are sent, call Flush before your function returns.
//
// However, we don't want to do that unconditionally: in non-serverless
// contexts, the batching behavior is preferable. Therefore, an environment-
// controlled variable.
var autoflush = false

func init() {
	if os.Getenv("HONEYCOMB_AUTOFLUSH") == "1" {
		autoflush = true
	}
}

////////////////////////////////////////////////////////////////////////////////
// Honeycomb.io Logrus hook
////////////////////////////////////////////////////////////////////////////////

// A HoneycombHook is a hook compatible with logrus which dispatches logged
// messages to honeycomb.
type HoneycombHook struct {
}

// Fire implements logrus.Hook
func (hook *HoneycombHook) Fire(entry *logrus.Entry) error {
	eventBuilder := libhoney.NewBuilder()
	honeycombEvent := eventBuilder.NewEvent()
	const binKey string = "bin"
	const levelKey string = "level"
	foundBin := false
	foundLevel := false
	for eachKey, eachValue := range entry.Data {
		honeycombEvent.AddField(eachKey, eachValue)
		switch eachKey {
		case binKey:
			foundBin = true
		case levelKey:
			foundLevel = true
		}
	}
	if !foundLevel {
		honeycombEvent.AddField("level", entry.Level.String())
	}
	if !foundBin {
		honeycombEvent.AddField(binKey, filepath.Base(os.Args[0]))
	}
	// Use cryptic values for these common fields, so there's
	// less of a chance to conflict with any keys in entry.Data.
	honeycombEvent.AddField("_ts", entry.Time)
	honeycombEvent.AddField("_txt", entry.Message)
	honeycombEvent.Send()
	if autoflush {
		libhoney.Flush()
	}
	return nil
}

// Levels implements logrus.Hook
func (hook *HoneycombHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.InfoLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
	}
}

// Flush ensures that all queued messages are dispatched immediately to Honeycomb.
func (*HoneycombHook) Flush() {
	libhoney.Flush()
}

// There are two things we should only do once -- one is initialize the libhoney library,
// and the other is registration of the logrus hook. Consequently, we need two instances
// of a sync.Once.
var setupOnce sync.Once
var registerLogrusOnce sync.Once

// If we tried to setup honeycomb and failed, don't ever try again,
// and make sure we still return the same error.
var setupError error

// setup initializes honeycomb by calling the Init function only once.
func setup() error {
	setupOnce.Do(func() {
		writeKey := os.Getenv("HONEYCOMB_KEY")
		datasetName := os.Getenv("HONEYCOMB_DATASET")
		cfg := libhoney.Config{
			WriteKey: writeKey,
			Dataset:  datasetName,
		}
		setupError = libhoney.Init(cfg)
		if setupError != nil {
			return
		}
		_, setupError = libhoney.VerifyWriteKey(cfg)
	})
	return setupError
}

// newLogrusHook returns a new Honeycomb.io logrus hook
func newLogrusHook() (logrus.Hook, error) {
	err := setup()
	if err != nil {
		return nil, err
	}
	return &HoneycombHook{}, nil
}

// Setup sets up a logrus logger to send its data to honeycomb instead of
// sending it to stdout.
func Setup(logger *logrus.Logger) *logrus.Logger {
	honeycombLoggingHook, err := newLogrusHook()
	if err != nil {
		logger.Warn(err)
		logger.Warn("Honeycomb failed to initialize properly - did you set HONEYCOMB_KEY and HONEYCOMB_DATASET? Not changing logger.")
		return logger
	}

	registerLogrusOnce.Do(func() {
		logger.Hooks.Add(honeycombLoggingHook)
	})
	return logger
}

type honeycombWriter struct{}

// Write implements io.Writer for honeycombWriter; it assumes that b is a JSON blob
// and unmarshals it into an interface{}, then simply sends that as a new event
// to honeycomb.
func (h *honeycombWriter) Write(b []byte) (int, error) {
	var data map[string]interface{}
	err := json.Unmarshal(b, &data)
	if err != nil {
		return 0, err
	}

	data = expandFieldsIn(data, "_msg")
	evt := libhoney.NewBuilder().NewEvent()
	err = evt.Add(data)
	if err != nil {
		return 0, err
	}

	err = evt.Send()
	if err != nil {
		return 0, err
	}

	return len(b), nil
}

var _ io.Writer = (*honeycombWriter)(nil)

// NewWriter constructs a writer that assumes its input is JSON and
// sends it to Honeycomb.
func NewWriter() (io.Writer, error) {
	err := setup()
	if err != nil {
		return nil, err
	}
	return &honeycombWriter{}, nil
}

// Tendermint seems to shove a blob of badly-formatted data into _msg, so we
// check for that case and try to extract key/value pairs from it.
// But not everything matches that way so we also keep _msg around
func expandFieldsIn(data map[string]interface{}, field string) map[string]interface{} {
	if m, ok := data[field]; ok {
		// pattern for matching lines that have key: value
		lpat := regexp.MustCompile(`^([A-Z][A-Za-z0-9]+):[ \t]*(.*[^{])$`)
		// pattern for splitting up lines
		spat := regexp.MustCompile(`[ \t]*\n[ \t]*`)
		ss := spat.Split(m.(string), -1)
		for _, s := range ss {
			r := lpat.FindStringSubmatch(s)
			if r != nil {
				n, err := strconv.Atoi(r[2])
				if err != nil {
					data[r[1]] = r[2]
				} else {
					data[r[1]] = n
				}
			}
		}
	}
	return data
}
