package models

type BlockchainStatusResponse struct {
	Connected       bool   `json:"connected"`
	Network         string `json:"network"`
	WalletAddress   string `json:"wallet_address,omitempty"`
	WalletBalance   string `json:"wallet_balance,omitempty"`
	BlockNumber     uint64 `json:"block_number,omitempty"`
	ContractAddress string `json:"contract_address,omitempty"`
	Error           string `json:"error,omitempty"`
}
