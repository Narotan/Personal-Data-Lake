package auth

import (
	"context"
	"fmt"
	"time"

	"DataLake/internal/logger"
)

type TokenManager struct {
	storage  TokenStorage
	provider Provider
}

func NewTokenManager(storage TokenStorage, provider Provider) *TokenManager {
	return &TokenManager{
		storage:  storage,
		provider: provider,
	}
}

// GetValidToken возвращает валидный токен, обновляя его при необходимости
func (tm *TokenManager) GetValidToken(ctx context.Context, providerName string) (TokenResponse, error) {
	log := logger.Get()

	token, err := tm.storage.LoadToken(providerName)
	if err != nil {
		log.Error().Err(err).Str("provider", providerName).Msg("failed to load token")
		return TokenResponse{}, fmt.Errorf("failed to load token: %w", err)
	}

	if !tm.isTokenExpired(token) {
		log.Debug().Str("provider", providerName).Msg("token is still valid")
		return token, nil
	}

	log.Info().Str("provider", providerName).Msg("token expired, refreshing...")

	newToken, err := tm.provider.RefreshToken(ctx, token.RefreshToken)
	if err != nil {
		log.Error().Err(err).Str("provider", providerName).Msg("failed to refresh token")
		return TokenResponse{}, fmt.Errorf("failed to refresh token: %w", err)
	}

	if newToken.RefreshToken == "" {
		newToken.RefreshToken = token.RefreshToken
	}

	if err := tm.storage.SaveToken(providerName, newToken); err != nil {
		log.Error().Err(err).Str("provider", providerName).Msg("failed to save refreshed token")
		return TokenResponse{}, fmt.Errorf("failed to save refreshed token: %w", err)
	}

	log.Info().Str("provider", providerName).Msg("token successfully refreshed")

	return newToken, nil
}

// isTokenExpired проверяет истёк ли токен
func (tm *TokenManager) isTokenExpired(token TokenResponse) bool {
	if token.ExpiresAt == "" {
		return true
	}

	expiresAt, err := time.Parse(time.RFC3339, token.ExpiresAt)
	if err != nil {
		return true
	}

	// Добавляем буфер в 5 минут для безопасности
	return time.Now().Add(5 * time.Minute).After(expiresAt)
}
