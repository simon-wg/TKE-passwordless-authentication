package structs

type VerifyRequest struct {
	Username  string `json:"username"`
	Signature []byte `json:"signature"`
}

type LoginRequest struct {
	Username string `json:"username"`
}

type LoginResponse struct {
	Challenge string `json:"challenge"`
	Signature string `json:"signature"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Pubkey   []byte `json:"pubkey"`
}
