package models

type (
	SubmitProposalMessage struct {
		TxHash         string `json:"tx_hash"`
		Proposer       string `json:"proposer"`
		InitialDeposit string `json:"initial_deposit"`
		Title          string `json:"title"`
		Description    string `json:"description"`
		Type           string `json:"type"`
		ProposalID     int64  `json:"proposal_id"`
		Height         int64  `json:"height"`
		MsgIndex       int64  `json:"msg_index"`
	}
)
