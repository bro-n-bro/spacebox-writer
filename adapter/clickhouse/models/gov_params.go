package models

type GovParams struct {
	DepositParams string `json:"deposit_params"`
	VotingParams  string `json:"voting_params"`
	TallyParams   string `json:"tally_params"`
	Height        int64  `json:"height"`
}
