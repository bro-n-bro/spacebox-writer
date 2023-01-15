package models

type MultiSendMessage struct {
	Coins       string `json:"coins"`
	AddressFrom string `json:"address_from"`
	AddressesTo string `json:"addresses_to"`
	TxHash      string `json:"tx_hash"`
	Height      int64  `json:"height"`
	MsgIndex    int64  `json:"msg_index"`
}
