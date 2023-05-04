package models

type (
	// WithdrawValidatorCommissionMessage represents withdraw validator commission message
	WithdrawValidatorCommissionMessage struct {
		ValidatorAddress   string `json:"validator_address"`
		TxHash             string `json:"tx_hash"`
		WithdrawCommission string `json:"withdraw_commission"`
		MsgIndex           int64  `json:"msg_index"`
		Height             int64  `json:"height"`
	}
)
