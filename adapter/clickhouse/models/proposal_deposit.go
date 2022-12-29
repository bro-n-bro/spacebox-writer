package models

type ProposalDeposit struct {
	ProposalID       uint64 `json:"proposal_id"`
	DepositorAddress string `json:"depositor_address"`
	Coins            string `json:"coins"`
	Height           int64  `json:"height"`
}
