package models

type (
	// ProposerReward represents proposer reward
	ProposerReward struct {
		Validator string `json:"validator"`
		Reward    string `json:"reward"`
		Height    int64  `json:"height"`
	}
)
