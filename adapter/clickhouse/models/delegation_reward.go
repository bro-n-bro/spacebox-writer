package models

type (
	// DelegationReward represents delegation reward
	DelegationReward struct {
		OperatorAddress  string `json:"operator_address"`
		DelegatorAddress string `json:"delegator_address"`
		WithdrawAddress  string `json:"withdraw_address"`
		Coins            string `json:"coins"`
		Height           int64  `json:"height"`
	}
)
