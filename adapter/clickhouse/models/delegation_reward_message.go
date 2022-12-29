package models

type DelegationRewardMessage struct {
	Coins            string `json:"coins"`
	Height           int64  `json:"height"`
	DelegatorAddress string `json:"delegator_address"`
	ValidatorAddress string `json:"validator_address"`
	TxHash           string `json:"tx_hash"`
}
