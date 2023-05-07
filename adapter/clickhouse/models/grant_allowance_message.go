package models

import "database/sql"

type (
	// GrantAllowanceMessage represents grant allowance message
	GrantAllowanceMessage struct {
		Expiration sql.NullTime `json:"expiration"`
		Granter    string       `json:"granter"`
		Grantee    string       `json:"grantee"`
		Allowance  string       `json:"allowance"`
		TxHash     string       `json:"tx_hash"`
		Height     int64        `json:"height"`
		MsgIndex   int64        `json:"msg_index"`
	}
)
