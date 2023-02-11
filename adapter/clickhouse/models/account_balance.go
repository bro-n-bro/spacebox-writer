package models

type (
	// AccountBalance represents account balance
	AccountBalance struct {
		Address string `json:"address"`
		Coins   string `json:"coins"`
		Height  int64  `json:"height"`
	}
)
