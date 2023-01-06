package models

type DelegationMessage struct {
	OperatorAddress  string `json:"operator_address"`
	DelegatorAddress string `json:"delegator_address"`
	Coin             string `json:"coin"`
	TxHash           string `json:"tx_hash"`
	Height           int64  `json:"height"`
}
