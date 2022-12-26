package models

import "time"

type Redelegation struct {
	CompletionTime      time.Time `json:"completion_time"`
	Coin                string    `json:"coin"`
	DelegatorAddress    string    `json:"delegator_address"`
	SrcValidatorAddress string    `json:"src_validator"`
	DstValidatorAddress string    `json:"dst_validator"`
	Height              int64     `json:"height"`
}
