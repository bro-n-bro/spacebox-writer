package models

type SendMessage struct {
	Coins       string `json:"coins"`
	AddressFrom string `json:"address_from"`
	AddressTo   string `json:"address_to"`
	TxHash      string `json:"tx_hash"`
	Height      int64  `json:"height"`
	MsgIndex    int64  `json:"msg_index"`
}
