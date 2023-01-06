package models

import "time"

type (
	Blocks []*Block

	Block struct {
		Created time.Time ``
		From    string    ``
		To      string    ``
		Value   float64   ``
	}
)
