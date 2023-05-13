package models

import "database/sql"

type (
	// GrantMessage represents grant message
	GrantMessage struct {
		Expiration sql.NullTime `json:"expiration"`
		Granter    string       `json:"granter"`
		Grantee    string       `json:"grantee"`
		TxHash     string       `json:"tx_hash"`
		MsgType    string       `json:"msg_hash"`
		Height     int64        `json:"height"`
		MsgIndex   int64        `json:"msg_index"`
	}
)
