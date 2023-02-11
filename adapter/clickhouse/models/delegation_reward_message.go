package models

type (
	// DelegationRewardMessage represents delegation reward message
	DelegationRewardMessage struct {
		Coins            string `json:"coins"`
		DelegatorAddress string `json:"delegator_address"`
		ValidatorAddress string `json:"validator_address"`
		TxHash           string `json:"tx_hash"`
		Height           int64  `json:"height"`
		MsgIndex         int64  `json:"msg_index"`
	}
)
