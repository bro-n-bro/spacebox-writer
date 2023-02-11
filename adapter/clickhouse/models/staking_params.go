package models

type (
	// StakingParams represents staking params
	StakingParams struct {
		Params string `json:"params"`
		Height int64  `json:"height"`
	}
)
