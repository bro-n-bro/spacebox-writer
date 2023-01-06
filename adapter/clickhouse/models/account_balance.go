package models

type AccountBalance struct {
	Address string `json:"address"`
	Coins   string `json:"coins"`
	Height  int64  `json:"height"`
}
