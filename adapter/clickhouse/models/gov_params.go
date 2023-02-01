package models

type GovParams struct {
	DepositParams interface{} `json:"deposit_params"`
	VotingParams  interface{} `json:"voting_params"`
	TallyParams   interface{} `json:"tally_params"`
	Height        int64       `json:"height"`
}

//type GovParams struct {
//	DepositParams JSONString `json:"deposit_params"`
//	VotingParams  JSONString `json:"voting_params"`
//	TallyParams   JSONString `json:"tally_params"`
//	Height        int64      `json:"height"`
//}
