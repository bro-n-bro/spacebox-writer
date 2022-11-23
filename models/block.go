package models

import "time"

type (
	Blocks []*Block

	Block struct {
		From    string    `` // From
		To      string    `` // To
		Value   float64   `` // Value
		Created time.Time `` // Created
	}
)
