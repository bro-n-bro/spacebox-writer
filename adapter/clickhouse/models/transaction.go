package models

type Transaction struct {
	RawLog      string   `json:"raw_log"`
	SignerInfos string   `json:"signer_infos"`
	Logs        string   `json:"logs"`
	Messages    string   `json:"messages"`
	Memo        string   `json:"memo"`
	Signer      string   `json:"signer"`
	Hash        string   `json:"hash"`
	Fee         string   `json:"fee,omitempty"`
	Signatures  []string `json:"signatures"`
	GasUsed     int64    `json:"gas_used"`
	GasWanted   int64    `json:"gas_wanted"`
	Height      int64    `json:"height"`
	Success     bool     `json:"success"`
}
