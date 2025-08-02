// internal/models/keycloak.go - UPDATED
package models

type LoginParams struct {
	Username string `json:"username"` // Email yerine username
	Password string `json:"password"`
}

type RegisterParams struct {
	Firstname string      `json:"firstname"`
	Lastname  string      `json:"lastname"`
	Username  string      `json:"username"`
	Email     string      `json:"email"`    // Email ayrı field olarak
	Password  string      `json:"password"` // Password direkt olarak
}

type UserPayload struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Username  string `json:"username"`
	Email     string `json:"email"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}
