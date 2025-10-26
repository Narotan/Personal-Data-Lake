package auth

import (
	"context"
)

// Provider определяет общий интерфейс для работы с OAuth2 провайдерами
type Provider interface {
	// GetAuthURL возвращает URL для авторизации пользователя
	GetAuthURL(state string) string

	// ExchangeToken обменивает код авторизации на токены
	ExchangeToken(ctx context.Context, code string) (TokenResponse, error)

	// RefreshToken обновляет access token используя refresh token
	RefreshToken(ctx context.Context, refreshToken string) (TokenResponse, error)
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    string `json:"expires_at"`
	TokenType    string `json:"token_type,omitempty"`
	Scope        string `json:"scope,omitempty"`
	UID          string `json:"uid,omitempty"`
}
