// Copyright (c) 2018 Oneiro NA, Inc. All rights reserved.

package honeycomb

import (
	"encoding/json"
	"io"
	"os"
	"regexp"
	"strconv"
	"sync"

	libhoney "github.com/honeycombio/libhoney-go"
	"github.com/sirupsen/logrus"
)

////////////////////////////////////////////////////////////////////////////////
// Honeycomb.io Logrus hook
////////////////////////////////////////////////////////////////////////////////
type honeycombHook struct {
}

func (hook *honeycombHook) Fire(entry *logrus.Entry) error {
	eventBuilder := libhoney.NewBuilder()
	honeycombEvent := eventBuilder.NewEvent()
	for eachKey, eachValue := range entry.Data {
		honeycombEvent.AddField(eachKey, eachValue)
	}
	honeycombEvent.AddField("ts", entry.Time)
	honeycombEvent.Send()
	return nil
}

func (hook *honeycombHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.InfoLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
	}
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
	return &honeycombHook{}, nil
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
