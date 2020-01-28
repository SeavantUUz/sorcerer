package watcher

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"sorcerer/constants"
	grpcpb "sorcerer/rpc"
	"sorcerer/util"
	"testing"
)

func TestContractApplyWatcher_TypeAssert(t *testing.T) {
	assert := assert.New(t)
	conn, err := util.Dial(constants.Node)
	client := grpcpb.NewApiServiceClient(conn)
	if err != nil {
		t.Error(err)
	}
	signBlockResponse, err := client.GetSignedBlock(context.Background(), &grpcpb.GetSignedBlockRequest{Start: 428})
	if err != nil {
		t.Error(err)
	}
	block := signBlockResponse.Block
	transactions := util.ExtractTransactions(block)
	transaction := transactions[0]
	w, err := NewContractApplyWatcher(logrus.New())
	if err != nil {
		t.Error(err)
	}
	op := w.typeAssert(transaction)
	assert.Equal(op.Owner.Value, "initminer")
	assert.Equal(op.Contract, "richman")
	assert.Equal(op.Amount.Value, uint64(2000000))
}
