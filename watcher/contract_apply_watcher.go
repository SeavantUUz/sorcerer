package watcher

import "sorcerer/structure"

type ContractApplyWatcher struct {
	token string
}

func NewContractApplyWatcher() *ContractApplyWatcher {
	return &ContractApplyWatcher{token: "hha"}
}

func (w ContractApplyWatcher) Agree(trx *structure.Transaction) bool {
	return true
}

func (w ContractApplyWatcher) Do(trx *structure.Transaction) string {
	return "hello"
}

func (w ContractApplyWatcher) Channel() string {
	return "hello"
}
