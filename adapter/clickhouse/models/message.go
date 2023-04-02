package models

type (
	// Message represents message
	Message struct {
		TransactionHash           string   `json:"transaction_hash"`
		Type                      string   `json:"type"`
		Value                     string   `json:"value"`
		Signer                    string   `json:"signer"`
		InvolvedAccountsAddresses []string `json:"involved_accounts_addresses"`
		MsgIndex                  int64    `json:"msg_index"`
	}
)
