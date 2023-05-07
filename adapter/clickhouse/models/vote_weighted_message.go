package models

type (
	// VoteWeightedMessage represents vote weighted message
	VoteWeightedMessage struct {
		WeightedVoteOption string `json:"weighted_vote_option"`
		TxHash             string `json:"tx_hash"`
		Voter              string `json:"voter"`
		MsgIndex           int64  `json:"msg_index"`
		ProposalID         int64  `json:"proposal_id"`
		Height             int64  `json:"height"`
	}
)
