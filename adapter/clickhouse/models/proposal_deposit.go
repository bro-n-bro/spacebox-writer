package models

type ProposalDeposit struct {
	DepositorAddress string `json:"depositor_address"`
	Coins            string `json:"coins"`
	ProposalID       uint64 `json:"proposal_id"`
	Height           int64  `json:"height"`
}
