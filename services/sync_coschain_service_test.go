package services

import (
	"github.com/sirupsen/logrus"
	"sorcerer/structure"
	"sorcerer/util"
	"testing"
)

func TestSyncServiceBlockInfo(t *testing.T) {
	log := logrus.New()
	ch := make(chan *structure.Transaction)
	s, err := NewSyncService(log, ch)
	if err != nil {
		t.Error(err)
	}
	err = s.initRpc()
	if err != nil {
		t.Error(err)
	}
	signBlockResponse, err := s.blockInfo(428)
	if err != nil {
		t.Error(err)
	}
	block := signBlockResponse.Block
	transactions := util.ExtractTransactions(block)
	for _, tx := range transactions {
		t.Log(tx)
	}
}
