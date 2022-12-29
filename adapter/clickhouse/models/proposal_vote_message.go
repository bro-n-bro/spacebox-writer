package models

type ProposalVoteMessage struct {
	ProposalID   int64  `json:"proposal_id"`
	VoterAddress string `json:"voter"`
	Option       string `json:"option"`
	Height       int64  `json:"height"`
}
