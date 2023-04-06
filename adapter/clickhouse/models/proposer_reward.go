package models

type (
	// ProposerReward represents proposer reward
	ProposerReward struct {
		Height    int64  `json:"height"`
		Validator string `json:"validator"`
		Reward    string `json:"reward"`
	}
)
