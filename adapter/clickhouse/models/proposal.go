package models

import (
	"database/sql"
)

type (
	// Proposal represents proposal
	Proposal struct {
		SubmitTime      sql.NullTime `json:"submit_time"`
		DepositEndTime  sql.NullTime `json:"deposit_end_time"`
		VotingStartTime sql.NullTime `json:"voting_start_time"`
		VotingEndTime   sql.NullTime `json:"voting_end_time"`
		Title           string       `json:"title"`
		Description     string       `json:"description"`
		ProposalRoute   string       `json:"proposal_route"`
		ProposalType    string       `json:"proposal_type"`
		ProposerAddress string       `json:"proposer_address"`
		Status          string       `json:"status"`
		Content         string       `json:"content"`
		ID              int64        `json:"id"`
	}
)
