package util

import (
	"encoding/hex"
	prototype "sorcerer/goproto"
	"sorcerer/structure"
	"time"
)

func FindCreator(operation *prototype.Operation) (name string) {
	signers := make(map[string]bool)
	prototype.GetBaseOperation(operation).GetSigner(&signers)
	if len(signers) > 0 {
		for s := range signers {
			name = s
			break
		}
	}
	return
}

type Op map[string]interface{}

func PurgeOperation(operations []*prototype.Operation) []Op {
	var ops []Op
	for _, operation := range operations {
		ops = append(ops, Op{prototype.GetGenericOperationName(operation): prototype.GetBaseOperation(operation)})
	}
	return ops
}

func ExtractTransactions(block *prototype.SignedBlock) []*structure.Transaction {
	header := block.SignedHeader
	transactions := block.Transactions
	var trxs []*structure.Transaction
	for _, t := range transactions {
		invoice := t.Receipt
		trxOperations := t.SigTrx.Trx.GetOperations()
		trxHash, _ := t.SigTrx.Id()
		trxId := hex.EncodeToString(trxHash.Hash)
		creator := FindCreator(t.SigTrx.GetTrx().GetOperations()[0])
		data := block.Id().Data
		blockId := hex.EncodeToString(data[:])
		blockProducer := header.Header.BlockProducer.Value
		blockTime := header.Header.Timestamp.UtcSeconds
		blockHeight := header.Number()
		operations := PurgeOperation(trxOperations)
		trx := &structure.Transaction{
			BlockId:       blockId,
			BlockHeight:   blockHeight,
			BlockTime:     time.Unix(int64(blockTime), 0),
			BlockProducer: blockProducer,
			Creator:       creator,
			TrxHash:       trxId,
			Invoice:       invoice,
			Operations:    operations,
		}
		trxs = append(trxs, trx)
	}
	return trxs
}
