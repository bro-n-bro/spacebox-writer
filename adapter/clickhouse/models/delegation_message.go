package models

type DelegationMessage struct {
	OperatorAddress  string `json:"operator_address"`
	DelegatorAddress string `json:"delegator_address"`
	Coin             string `json:"coin"`
	Height           int64  `json:"height"`
	TxHash           string `json:"tx_hash"`
}
