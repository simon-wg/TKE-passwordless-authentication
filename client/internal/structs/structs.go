package structs

// RegisterRequest represents the data required to register a new user
// It includes the username and the user's public key
type RegisterRequest struct {
	Username string `json:"username"`
	Pubkey   []byte `json:"pubkey"`
}

// LoginRequest represents the payload for a login request
// It contains the username of the user attempting to log in
type LoginRequest struct {
	Username string `json:"username"`
}

// LoginResponse represents the response received after a login attempt
// It contains a challenge string and a signature string, both of which are used to verify the authenticity of the login request
type LoginResponse struct {
	Challenge string `json:"challenge"`
	Signature string `json:"signature"`
}

// VerifyRequest represents a request to verify a user's identity
// It contains the username of the user and a cryptographic signature to authenticate the request
type VerifyRequest struct {
	Username  string `json:"username"`
	Signature []byte `json:"signature"`
}

type VerifyResponse struct {
}

// AddPublicKeyRequest represents a request to add a new public key for a user
// It contains the username of the user and the new public key to be added
type AddPublicKeyRequest struct {
	Username string `json:"username"`
	Pubkey   []byte `json:"pubkey"`
}

// RemovePublicKeyRequest represents a request to remove a public key for a user
// It contains the username of the user and the public key to be removed
type RemovePublicKeyRequest struct {
	Username string `json:"username"`
	Pubkey   []byte `json:"pubkey"`
}
