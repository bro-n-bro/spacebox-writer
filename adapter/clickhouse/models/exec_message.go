package models

type (
	// ExecMessage represents exec message
	ExecMessage struct {
		Grantee  string `json:"grantee"`
		Msgs     string `json:"msgs"`
		TxHash   string `json:"tx_hash"`
		Height   int64  `json:"height"`
		MsgIndex int64  `json:"msg_index"`
	}
)
