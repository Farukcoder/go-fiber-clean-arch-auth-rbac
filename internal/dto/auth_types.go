package dto

import "github.com/golang-jwt/jwt/v5"

// LoginRequest is the JSON payload for the login endpoint.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterRequest is the JSON payload for the registration endpoint.
type RegisterRequest struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	Phone           string `json:"phone"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

// LoginInput supports login by email OR phone number.
type LoginInput struct {
	EmailOrPhone string `json:"email_or_phone"`
	Password     string `json:"password"`
}

// TokenResponse is the JSON payload returned after a successful login or refresh.
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}

// RefreshRequest is the JSON payload for the token-refresh and logout endpoints.
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// AuthClaims carries the custom JWT claims used by both access and refresh tokens.
type AuthClaims struct {
	UserID    int64  `json:"user_id"`
	RoleID    int64  `json:"role_id"`
	RoleName  string `json:"role_name"`
	Email     string `json:"email"`
	TokenType string `json:"token_type"`
	JTI       string `json:"jti,omitempty"`
	jwt.RegisteredClaims
}
