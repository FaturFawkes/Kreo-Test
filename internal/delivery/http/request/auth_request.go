package request

// AuthRequest represents the auth token request payload
type AuthRequest struct {
	APIKey string `json:"api_key" binding:"required"`
}
