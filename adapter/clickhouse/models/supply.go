package models

type (
	// Supply represents supply
	Supply struct {
		Coins  string `json:"coins"`
		Height int64  `json:"height"`
	}
)
