package watcher

import "sorcerer/structure"

type Watcher interface {
	Agree(trx *structure.Transaction) bool
	Do(trx *structure.Transaction) string
	Channel() string
}
