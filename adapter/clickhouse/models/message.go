package models

type Message struct {
	TransactionHash           string `json:"transaction_hash"`
	Type                      string `json:"type"`
	Value                     string `json:"value"`
	InvolvedAccountsAddresses string `json:"involved_accounts_addresses"`
	Signer                    string `json:"signer"`
	Index                     int    `json:"index"`
}
