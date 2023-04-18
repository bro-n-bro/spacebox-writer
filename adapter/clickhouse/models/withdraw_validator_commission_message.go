package models

type (
	// WithdrawValidatorCommissionMessage represents withdraw validator commission message
	WithdrawValidatorCommissionMessage struct {
		SenderAddress      string `json:"sender_address"`
		TxHash             string `json:"tx_hash"`
		WithdrawCommission int64  `json:"withdraw_commission"`
		MsgIndex           int64  `json:"msg_index"`
		Height             int64  `json:"height"`
	}
)
