package auth

import (
	"encoding/json"
	"os"
)

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    string `json:"expires_at"`
}

func SaveTokens(token Tokens) error {
	data, err := json.Marshal(token)
	if err != nil {
		return err
	}
	return os.WriteFile("token.json", data, 0644)
}

func LoadTokens() (Tokens, error) {
	data, err := os.ReadFile("token.json")
	if err != nil {
		return Tokens{}, err
	}
	var token Tokens
	err = json.Unmarshal(data, &token)
	return token, err
}
