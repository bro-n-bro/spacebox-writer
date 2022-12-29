package models

type Transaction struct {
	Hash        string `json:"hash"`
	Height      int64  `json:"height"`
	Success     bool   `json:"success"`
	Messages    string `json:"messages"`
	Memo        string `json:"memo"`
	Signatures  string `json:"signatures"`
	Signer      string `json:"signer"`
	GasWanted   int64  `json:"gas_wanted"`
	GasUsed     int64  `json:"gas_used"`
	RawLog      string `json:"raw_log"`
	Logs        string `json:"logs"`
	SignerInfos string `json:"signer_infos"`
	Fee         string `json:"fee,omitempty"`
}
