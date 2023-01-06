package models

import "time"

type Block struct {
	Timestamp       time.Time `json:"timestamp"`
	Hash            string    `json:"hash"`
	ProposerAddress string    `json:"proposer_address"`
	Height          int64     `json:"height"`
	NumTXS          int64     `json:"num_txs"`
	TotalGas        int64     `json:"total_gas"`
}
