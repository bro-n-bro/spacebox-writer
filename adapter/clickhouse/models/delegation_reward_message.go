package models

type (
	// DelegationRewardMessage represents delegation reward message
	DelegationRewardMessage struct {
		DelegatorAddress string `json:"delegator_address"`
		ValidatorAddress string `json:"validator_address"`
		Coins            string `json:"coins"`
		TxHash           string `json:"tx_hash"`
		Height           int64  `json:"height"`
		MsgIndex         int64  `json:"msg_index"`
	}
)
