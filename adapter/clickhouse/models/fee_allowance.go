package models

import "database/sql"

type (
	// FeeAllowance represents fee allowance
	FeeAllowance struct {
		Expiration sql.NullTime `json:"expiration"`
		Granter    string       `json:"granter"`
		Grantee    string       `json:"grantee"`
		Allowance  string       `json:"allowance"`
		Height     int64        `json:"height"`
	}
)
