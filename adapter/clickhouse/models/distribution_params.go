package models

type (
	// DistributionParams represents distribution params
	DistributionParams struct {
		Params string `json:"params"`
		Height int64  `json:"height"`
	}
)
