package watcher

import "sorcerer/structure"

type Watcher interface {
	Agree(trx *structure.Transaction) bool
	Publish(trx *structure.Transaction)
}
