package models

type (
	// CommunityPool represents community pool
	CommunityPool struct {
		Coins  string `json:"coins"`
		Height int64  `json:"height"`
	}
)
