package models

type ProposalVoteMessage struct {
	VoterAddress string `json:"voter"`
	Option       string `json:"option"`
	ProposalID   uint64 `json:"proposal_id"`
	Height       int64  `json:"height"`
	MsgIndex     int64  `json:"msg_index"`
	TxHash       string `json:"tx_hash"`
}
