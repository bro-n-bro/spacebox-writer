package models

type (
	DistributionCommission struct {
		Validator string `json:"validator"`
		Amount    string `json:"amount"`
		Height    int64  `json:"height"`
	}
)
