package models

type (
	// DelegationMessage represents delegation message
	DelegationMessage struct {
		OperatorAddress  string `json:"operator_address"`
		DelegatorAddress string `json:"delegator_address"`
		Coin             string `json:"coin"`
		TxHash           string `json:"tx_hash"`
		Height           int64  `json:"height"`
		MsgIndex         int64  `json:"msg_index"`
	}
)
