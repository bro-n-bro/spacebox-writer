package models

type (
	// EditValidatorMessage represents edit validator message
	EditValidatorMessage struct {
		TxHash      string `json:"tx_hash"`
		Description string `json:"description"`
		Height      int64  `json:"height"`
		MsgIndex    int64  `json:"msg_index"`
	}
)
