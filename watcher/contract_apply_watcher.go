package watcher

import (
	"encoding/json"
	"github.com/pote/philote-go"
	"github.com/sirupsen/logrus"
	"reflect"
	"sorcerer/constants"
	"sorcerer/prototype"
	"sorcerer/structure"
	"sorcerer/util"
)

type ContractApplyWatcher struct {
	client  *philote.Client
	log     *logrus.Logger
	channel string
	token   string
}

func NewContractApplyWatcher(log *logrus.Logger) (*ContractApplyWatcher, error) {
	c, err := philote.NewClient("ws://localhost:6380", constants.TokenContractRW)
	if err != nil {
		return nil, err
	}
	return &ContractApplyWatcher{token: constants.TokenContractRW, client: c, channel: "contract-event", log: log}, nil
}

func (w *ContractApplyWatcher) typeAssert(trx *structure.Transaction) *prototype.ContractApplyOperation {
	switch reflect.TypeOf(trx.Operations).Kind() {
	case reflect.Slice:
		ops := reflect.ValueOf(trx.Operations)
		// ignore operations except the first
		opDesc := ops.Index(0)
		op, ok := opDesc.Interface().(util.Op)
		if !ok {
			return nil
		}
		baseOpInterface, ok := op["contract_apply"]
		if !ok {
			return nil
		}
		contractApplyOp, ok := reflect.ValueOf(baseOpInterface).Interface().(*prototype.ContractApplyOperation)
		if !ok {
			return nil
		}
		return contractApplyOp
	default:
		return nil
	}
}

func (w *ContractApplyWatcher) Agree(trx *structure.Transaction) bool {
	contractApplyOp := w.typeAssert(trx)
	if contractApplyOp == nil {
		return false
	}
	if contractApplyOp.Contract == "hello" && contractApplyOp.Owner.Value == "initminer" {
		return true
	} else {
		return false
	}
}

func (w *ContractApplyWatcher) Publish(trx *structure.Transaction) {
	contractApplyOp := w.typeAssert(trx)
	if data, err := json.Marshal(contractApplyOp); err == nil {
		_ = w.client.Publish(&philote.Message{
			Channel: w.channel,
			Data:    string(data),
		})
	} else {
		w.log.Errorf("Marshal %v to json failed", contractApplyOp)
	}

}
