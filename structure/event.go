package structure

import "time"

type Transaction struct {
	BlockId       string      `json:"block_id"`
	BlockHeight   uint64      `json:"block_height"`
	BlockTime     time.Time   `json:"block_time"`
	BlockProducer string      `json:"block_producer"`
	Creator       string      `json:"creator"`
	TrxHash       string      `json:"trx_hash"`
	Invoice       interface{} `json:"invoice"`
	Operations    interface{} `json:"operations"`
}
