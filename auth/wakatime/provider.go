package wakatimeauth

import (
	"DataLake/auth"
	"DataLake/internal/logger"
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Provider реализация OAuth2 провайдера для WakaTime
type Provider struct {
	clientID     string
	clientSecret string
	redirectURI  string
}

// NewProvider создает новый провайдер WakaTime
func NewProvider(clientID, clientSecret, redirectURI string) *Provider {
	return &Provider{
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURI:  redirectURI,
	}
}

// GetAuthURL возвращает URL для авторизации пользователя
func (p *Provider) GetAuthURL(state string) string {
	baseURL := "https://wakatime.com/oauth/authorize"
	params := url.Values{}
	params.Set("client_id", p.clientID)
	params.Set("redirect_uri", p.redirectURI)
	params.Set("response_type", "code")
	if state != "" {
		params.Set("state", state)
	}

	return baseURL + "?" + params.Encode()
}

// ExchangeToken обменивает код авторизации на токены
func (p *Provider) ExchangeToken(ctx context.Context, code string) (auth.TokenResponse, error) {
	log := logger.Get()

	baseURL := "https://wakatime.com/oauth/token"
	data := url.Values{}
	data.Set("client_id", p.clientID)
	data.Set("client_secret", p.clientSecret)
	data.Set("redirect_uri", p.redirectURI)
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)

	log.Info().Str("url", baseURL).Msg("exchanging code for token")

	req, err := http.NewRequestWithContext(ctx, "POST", baseURL, strings.NewReader(data.Encode()))
	if err != nil {
		log.Error().Err(err).Msg("failed to create request")
		return auth.TokenResponse{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("failed to execute request")
		return auth.TokenResponse{}, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("failed to read response body")
		return auth.TokenResponse{}, err
	}

	values, err := url.ParseQuery(string(bodyBytes))
	if err != nil {
		log.Error().Err(err).Msg("failed to parse response")
		return auth.TokenResponse{}, err
	}

	token := auth.TokenResponse{
		AccessToken:  values.Get("access_token"),
		RefreshToken: values.Get("refresh_token"),
		Scope:        values.Get("scope"),
		TokenType:    values.Get("token_type"),
		UID:          values.Get("uid"),
		ExpiresAt:    values.Get("expires_at"),
	}

	log.Info().Str("uid", token.UID).Msg("successfully exchanged code for token")

	return token, nil
}

// RefreshToken обновляет access token используя refresh token
func (p *Provider) RefreshToken(ctx context.Context, refreshToken string) (auth.TokenResponse, error) {
	log := logger.Get()

	baseURL := "https://wakatime.com/oauth/token"
	data := url.Values{}
	data.Set("client_id", p.clientID)
	data.Set("client_secret", p.clientSecret)
	data.Set("redirect_uri", p.redirectURI)
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)

	log.Info().Str("url", baseURL).Msg("refreshing token")

	req, err := http.NewRequestWithContext(ctx, "POST", baseURL, strings.NewReader(data.Encode()))
	if err != nil {
		log.Error().Err(err).Msg("failed to create request")
		return auth.TokenResponse{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("failed to execute request")
		return auth.TokenResponse{}, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("failed to read response body")
		return auth.TokenResponse{}, err
	}

	values, err := url.ParseQuery(string(bodyBytes))
	if err != nil {
		log.Error().Err(err).Msg("failed to parse response")
		return auth.TokenResponse{}, err
	}

	token := auth.TokenResponse{
		AccessToken:  values.Get("access_token"),
		RefreshToken: values.Get("refresh_token"),
		Scope:        values.Get("scope"),
		TokenType:    values.Get("token_type"),
		UID:          values.Get("uid"),
		ExpiresAt:    values.Get("expires_at"),
	}

	log.Info().Msg("successfully refreshed token")

	return token, nil
}
