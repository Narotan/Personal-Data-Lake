package auth

import (
	"DataLake/internal/logger"
	"encoding/json"
	"os"
)

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    string `json:"expires_at"`
}

func SaveTokens(token Tokens) error {
	log := logger.Get()

	data, err := json.Marshal(token)
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal tokens")
		return err
	}

	err = os.WriteFile("token.json", data, 0644)
	if err != nil {
		log.Error().Err(err).Msg("failed to write tokens to file")
		return err
	}

	log.Info().Msg("tokens saved successfully")
	return nil
}

func LoadTokens() (Tokens, error) {
	log := logger.Get()

	data, err := os.ReadFile("token.json")
	if err != nil {
		log.Error().Err(err).Msg("failed to read tokens from file")
		return Tokens{}, err
	}

	var token Tokens
	err = json.Unmarshal(data, &token)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal tokens")
		return Tokens{}, err
	}

	log.Info().Msg("tokens loaded successfully")
	return token, err
}
