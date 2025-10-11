package auth

import (
	"DataLake/internal/logger"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func BuildAuthRequest(cfg Config) string {
	baseUrl := "https://wakatime.com/oauth/authorize"
	params := url.Values{}
	params.Set("client_id", cfg.ClientID)
	params.Set("redirect_uri", cfg.RedirectURI)
	params.Set("response_type", "code")

	fullURL := baseUrl + "?" + params.Encode()
	return fullURL
}

type AccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
	UID          string `json:"uid"`
	ExpiresAt    string `json:"expires_at"`
}

func ExchangeCodeForToken(cfg Config, code string) (AccessTokenResponse, error) {
	log := logger.Get()

	baseUrl := "https://wakatime.com/oauth/token"
	data := url.Values{}
	data.Set("client_id", cfg.ClientID)
	data.Set("client_secret", cfg.ClientSecret)
	data.Set("redirect_uri", cfg.RedirectURI)
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)

	log.Info().Str("url", baseUrl).Msg("exchanging code for token")

	req, err := http.NewRequest("POST", baseUrl, strings.NewReader(data.Encode()))
	if err != nil {
		log.Error().Err(err).Msg("failed to create request")
		return AccessTokenResponse{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("failed to execute request")
		return AccessTokenResponse{}, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("failed to read response body")
		return AccessTokenResponse{}, err
	}

	values, err := url.ParseQuery(string(bodyBytes))
	if err != nil {
		log.Error().Err(err).Msg("failed to parse response")
		return AccessTokenResponse{}, err
	}

	token := AccessTokenResponse{
		AccessToken:  values.Get("access_token"),
		RefreshToken: values.Get("refresh_token"),
		Scope:        values.Get("scope"),
		TokenType:    values.Get("token_type"),
		UID:          values.Get("uid"),
		ExpiresAt:    values.Get("expires_at"),
	}

	log.Info().Str("uid", token.UID).Msg("successfully exchanged code for token")

	return token, err
}
