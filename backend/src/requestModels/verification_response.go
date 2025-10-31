package requestModels

type VerificationResponse struct {
	Verified         bool     `json:"verified"`
	TotalUpdates     int      `json:"total_updates"`
	VerifiedUpdates  int      `json:"verified_updates"`
	BlockchainHashes int      `json:"blockchain_hashes"`
	Status           string   `json:"status"`
	Message          string   `json:"message"`
	Mismatches       []string `json:"mismatches,omitempty"`
}