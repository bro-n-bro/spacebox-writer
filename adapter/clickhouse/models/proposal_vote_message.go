package models

type (
	// ProposalVoteMessage represents proposal vote message
	ProposalVoteMessage struct {
		VoterAddress string `json:"voter"`
		Option       string `json:"option"`
		TxHash       string `json:"tx_hash"`
		ProposalID   int64  `json:"proposal_id"`
		Height       int64  `json:"height"`
		MsgIndex     int64  `json:"msg_index"`
	}
)
