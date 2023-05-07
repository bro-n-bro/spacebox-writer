package models

type (
	// DistributionReward represents distribution reward
	DistributionReward struct {
		Validator string `json:"validator"`
		Amount    string `json:"amount"`
		Height    int64  `json:"height"`
	}
)
