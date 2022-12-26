package models

import "time"

type UnbondingDelegation struct {
	CompletionTimestamp time.Time `json:"completion_timestamp"`
	Coin                string    `json:"coin"`
	DelegatorAddress    string    `json:"delegator_address"`
	ValidatorAddress    string    `json:"validator_oper_addr"`
	Height              int64     `json:"height"`
}
