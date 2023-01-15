package models

import "time"

type RedelegationMessage struct {
	CompletionTime      time.Time `json:"completion_time"`
	Coin                string    `json:"coin"`
	DelegatorAddress    string    `json:"delegator_address"`
	SrcValidatorAddress string    `json:"src_validator"`
	DstValidatorAddress string    `json:"dst_validator"`
	TxHash              string    `json:"tx_hash"`
	Height              int64     `json:"height"`
	MsgIndex            int64     `json:"msg_index"`
}
