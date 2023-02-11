package models

import "database/sql"

type (
	// Redelegation represents redelegation
	Redelegation struct {
		CompletionTime      sql.NullTime `json:"completion_time"`
		Coin                string       `json:"coin"`
		DelegatorAddress    string       `json:"delegator_address"`
		SrcValidatorAddress string       `json:"src_validator"`
		DstValidatorAddress string       `json:"dst_validator"`
		Height              int64        `json:"height"`
	}
)
