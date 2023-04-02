package models

type (
	// Delegation represents delegation
	Delegation struct {
		OperatorAddress  string `json:"operator_address"`
		DelegatorAddress string `json:"delegator_address"`
		Coin             string `json:"coin"`
		Height           int64  `json:"height"`
	}
)
