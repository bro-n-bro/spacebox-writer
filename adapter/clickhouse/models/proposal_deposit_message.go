package models

type ProposalDepositMessage struct {
	DepositorAddress string `json:"depositor_address"`
	Coins            string `json:"coins"`
	TxHash           string `json:"tx_hash"`
	ProposalID       int64  `json:"proposal_id"`
	Height           int64  `json:"height"`
	MsgIndex         int64  `json:"msg_index"`
}
