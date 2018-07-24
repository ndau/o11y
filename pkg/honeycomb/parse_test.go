package honeycomb

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParserReal(t *testing.T) {
	result, err := Parse("Parser1", []byte("Block{\n  Header{\n    ChainID:        test-chain-HMwr8o\n    Height:         23449\n    Time:           2018-07-22 23:44:32.6592659 +0000 UTC\n    NumTxs:         0\n    TotalTxs:       0\n    LastBlockID:    715A95922538778F136A65D86D974A3B273EA66E:1:2E904B3E22B2\n    LastCommit:     215A957565E17E4AEE20CB3E1667FBF15697A140\n    Data:           \n    Validators:     967B6BC72DD16F1695C763BFD27C8B64A81EA519\n    App:            D7F0E5DABF8CFA67182FFAC50526DF4F010841BA\n    Consensus:       D6B74BB35BDFFD8392340F2A379173548AE188FE\n    Results:        \n    Evidence:       \n  }#CF6F909CF3EE0316723FE043F0E50F6D0F4C166C\n  Data{\n    \n  }#\n  EvidenceData{\n    \n  }#\n  Commit{\n    BlockID:    715A95922538778F136A65D86D974A3B273EA66E:1:2E904B3E22B2\n    Precommits: Vote{0:EA0DD2EB887E 23448/00/2(Precommit) 715A95922538 /7301270C6CD1.../ @ 2018-07-22T23:44:31.652Z}\n  }#215A957565E17E4AEE20CB3E1667FBF15697A140\n}#CF6F909CF3EE0316723FE043F0E50F6D0F4C166C"),
		Debug(false))
	if err != nil {
		fmt.Println(err)
	}
	assert.Nil(t, err)
	j, err := json.MarshalIndent(result, "", "  ")
	assert.Nil(t, err)
	fmt.Println(string(j))
}

func TestParserReal2(t *testing.T) {
	s := `Block{
  Header{
    ChainID:        test-chain-kqCFfl
    Height:         21
    Time:           2018-07-23 23:35:38.9258088 +0000 UTC
    NumTxs:         0
    TotalTxs:       0
    LastBlockID:    B769E9117DE45551306EE9AAA3EDEDE840FF4A23:1:18E18ACF84EF
    LastCommit:     6B2EDB94EBAC877EF302290B7752C8CA5C2E05E4
    Data:
    Validators:     3E59692101DB6D12A9FF970FD85A231718ED1E92
    App:            C9BAC9796940AC8CE88AD2C13D19CDCEED462C3A
    Consensus:       D6B74BB35BDFFD8392340F2A379173548AE188FE
    Results:
    Evidence:
  }#4C56C110266D0D94CAC250963EC704C8C1864E02
  Data{

  }#
  EvidenceData{

  }#
  Commit{
    BlockID:    B769E9117DE45551306EE9AAA3EDEDE840FF4A23:1:18E18ACF84EF
    Precommits: Vote{0:EBC3C5DD8A4B 20/00/2(Precommit) B769E9117DE4 /A04C3646CFAF.../ @ 2018-07-23T23:35:37.918Z}
  }#6B2EDB94EBAC877EF302290B7752C8CA5C2E05E4
}#4C56C110266D0D94CAC250963EC704C8C1864E02`

	result, err := Parse("Parser1", []byte(s), Debug(false))
	if err != nil {
		fmt.Println(err)
	}
	assert.Nil(t, err)
	j, err := json.MarshalIndent(result, "", "  ")
	assert.Nil(t, err)
	fmt.Println(string(j))
}

func TestParser1(t *testing.T) {
	result, err := Parse("Parser1", []byte("Block{\nValue: 12\n}#"), Debug(false))
	if err != nil {
		fmt.Println(err)
	}
	assert.Nil(t, err)
	j, err := json.MarshalIndent(result, "", "  ")
	assert.Nil(t, err)
	fmt.Println(string(j))
}
