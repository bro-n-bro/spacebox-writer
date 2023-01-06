package models

type ProposalVoteMessage struct {
	VoterAddress string `json:"voter"`
	Option       string `json:"option"`
	ProposalID   int64  `json:"proposal_id"`
	Height       int64  `json:"height"`
}
