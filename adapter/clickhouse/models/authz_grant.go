package models

import "database/sql"

type (
	// AuthzGrant represents authz grant message
	AuthzGrant struct {
		Expiration     sql.NullTime `json:"expiration"`
		GranterAddress string       `json:"granter_address"`
		GranteeAddress string       `json:"grantee_address"`
		MsgType        string       `json:"msg_type"`
		Height         int64        `json:"height"`
	}
)
