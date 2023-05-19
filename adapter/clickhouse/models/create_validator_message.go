package models

type (
	// CreateValidatorMessage represents create validator message
	CreateValidatorMessage struct {
		TxHash            string  `json:"tx_hash"`
		DelegatorAddress  string  `json:"delegator_address"`
		ValidatorAddress  string  `json:"validator_address"`
		Description       string  `json:"description"`
		Pubkey            string  `json:"pubkey"`
		Height            int64   `json:"height"`
		MsgIndex          int64   `json:"msg_index"`
		MinSelfDelegation int64   `json:"min_self_delegation"`
		CommissionRates   float64 `json:"commission_rates"`
	}
)
