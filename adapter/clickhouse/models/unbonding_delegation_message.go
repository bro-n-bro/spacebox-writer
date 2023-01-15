package models

import "time"

type UnbondingDelegationMessage struct {
	CompletionTimestamp time.Time `json:"completion_timestamp"`
	Coin                string    `json:"coin"`
	DelegatorAddress    string    `json:"delegator_address"`
	ValidatorAddress    string    `json:"validator_oper_addr"`
	TxHash              string    `json:"tx_hash"`
	Height              int64     `json:"height"`
	MsgIndex            int64     `json:"msg_index"`
}
