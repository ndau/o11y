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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReal1(t *testing.T) {
	data := map[string]interface{}{
		"_msg": `Block{
  Header{
    Version:        {9 0}
    ChainID:        test-chain-OuwpVw
    Height:         54
    Time:           2019-02-07 23:01:59.268908973 +0000 UTC
    NumTxs:         0
    TotalTxs:       32
    LastBlockID:    8D5BF9FB560629C5A769FC0B81E8344CCCF09D3BE910D8AB6F04365DA0170692:1:FB253C748504
    LastCommit:     FC9BAB94965F7C00A896AF484610B7869973E7D06E537274407F70359B353E7A
    Data:
    Validators:     E3E0F4F720C407CE1D4BF9247F7894637E446EB66DE8ED6007B377D0D37C3DC4
    NextValidators: E3E0F4F720C407CE1D4BF9247F7894637E446EB66DE8ED6007B377D0D37C3DC4
    App:            EA0E1107DE0E660B17134C622CA1EE2130C74D1D
    Consensus:       048091BC7DDC283F77BFBF91D73C44DA58C3DF8A9CBC867405D8B7F3DAADA22F
    Results:        6E340B9CFFB37A989CA544E6BB780A2C78901D3FB33738768511A30617AFA01D
    Evidence:
    Proposer:       B349357FA7F27B7086029D9BB190FFEAE0642974
  }#338B7D2D130B01C7A40EB5A36D5E55F7F99AED066F3E8AF8EA77777CE9C9F9D9
  Data{

  }#
  EvidenceData{

  }#
  Commit{
    BlockID:    8D5BF9FB560629C5A769FC0B81E8344CCCF09D3BE910D8AB6F04365DA0170692:1:FB253C748504
    Precommits:
      Vote{0:96B8EFE99E02 53/00/2(Precommit) 8D5BF9FB5606 141F776838A3 @ 2019-02-07T23:01:59.268908973Z}
      Vote{1:B349357FA7F2 53/00/2(Precommit) 8D5BF9FB5606 B294DFE03F69 @ 2019-02-07T23:01:59.268908973Z}
  }#FC9BAB94965F7C00A896AF484610B7869973E7D06E537274407F70359B353E7A
}#338B7D2D130B01C7A40EB5A36D5E55F7F99AED066F3E8AF8EA77777CE9C9F9D9`}

	result := expandFieldsIn(data, "_msg")
	assert.Equal(t, result["Height"], 54)
	assert.Equal(t, result["Version"], "{9 0}")
	assert.Equal(t, result["ChainID"], "test-chain-OuwpVw")
	assert.Equal(t, result["Height"], 54)
	assert.Equal(t, result["Time"], "2019-02-07 23:01:59.268908973 +0000 UTC")
	assert.Equal(t, result["NumTxs"], 0)
	assert.Equal(t, result["TotalTxs"], 32)
	assert.Equal(t, result["LastBlockID"], "8D5BF9FB560629C5A769FC0B81E8344CCCF09D3BE910D8AB6F04365DA0170692:1:FB253C748504")
	assert.Equal(t, result["LastCommit"], "FC9BAB94965F7C00A896AF484610B7869973E7D06E537274407F70359B353E7A")
	assert.Equal(t, result["BlockID"], "8D5BF9FB560629C5A769FC0B81E8344CCCF09D3BE910D8AB6F04365DA0170692:1:FB253C748504")
}
