package dto

// LoginRequest represents the request payload for user authentication.
// It contains the necessary fields for logging in a user.
type LoginRequest struct {
	// Email is the user's email address, required for authentication.
	Email string `json:"email" binding:"required"`
	// Password is the user's password, required for authentication.
	Password string `json:"password" binding:"required"`
}
