package models

type AccountBalance struct {
	Address string `json:"address"`
	Height  int64  `json:"height"`
	Coins   string `json:"coins"`
}
