package googlecalendarauth

import (
	"DataLake/auth"
	"DataLake/internal/logger"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	googleAuthEndpoint  = "https://accounts.google.com/o/oauth2/v2/auth"
	googleTokenEndpoint = "https://oauth2.googleapis.com/token"
)

var defaultScopes = []string{
	"https://www.googleapis.com/auth/calendar.readonly",
	"https://www.googleapis.com/auth/calendar.events.readonly",
}

type Provider struct {
	clientID     string
	clientSecret string
	redirectURI  string
}

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
}

func NewProvider(clientID, clientSecret, redirectURI string) *Provider {
	return &Provider{
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURI:  redirectURI,
	}
}

// NewProviderFromEnv создаёт провайдер из переменных окружения
func NewProviderFromEnv() *Provider {
	return &Provider{
		clientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		clientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		redirectURI:  os.Getenv("GOOGLE_CALENDAR_REDIRECT_URI"),
	}
}

func (p *Provider) GetAuthURL(state string) string {
	u, _ := url.Parse(googleAuthEndpoint)

	q := url.Values{}
	q.Set("client_id", p.clientID)
	q.Set("redirect_uri", p.redirectURI)
	q.Set("response_type", "code")
	q.Set("scope", strings.Join(defaultScopes, " "))
	q.Set("access_type", "offline")
	q.Set("prompt", "consent")
	q.Set("include_granted_scopes", "true")
	if state != "" {
		q.Set("state", state)
	}

	u.RawQuery = q.Encode()
	return u.String()
}

func (p *Provider) ExchangeToken(ctx context.Context, code string) (auth.TokenResponse, error) {
	log := logger.Get()

	data := url.Values{}
	data.Set("code", code)
	data.Set("client_id", p.clientID)
	data.Set("client_secret", p.clientSecret)
	data.Set("redirect_uri", p.redirectURI)
	data.Set("grant_type", "authorization_code")

	log.Info().Str("url", googleTokenEndpoint).Msg("exchanging code for token")

	req, err := http.NewRequestWithContext(ctx, "POST", googleTokenEndpoint, strings.NewReader(data.Encode()))
	if err != nil {
		log.Error().Err(err).Msg("failed to create request")
		return auth.TokenResponse{}, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("failed to execute request")
		return auth.TokenResponse{}, fmt.Errorf("failed to exchange token: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("failed to read response body")
		return auth.TokenResponse{}, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Error().Str("status", resp.Status).Str("body", string(body)).Msg("token exchange failed")
		return auth.TokenResponse{}, fmt.Errorf("token exchange failed: %s - %s", resp.Status, string(body))
	}

	var tokenResp tokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		log.Error().Err(err).Msg("failed to parse response")
		return auth.TokenResponse{}, fmt.Errorf("failed to parse token response: %w", err)
	}

	log.Info().Msg("successfully exchanged code for token")

	return auth.TokenResponse{
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		ExpiresAt:    time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second).Format(time.RFC3339),
	}, nil
}

func (p *Provider) RefreshToken(ctx context.Context, refreshToken string) (auth.TokenResponse, error) {
	log := logger.Get()

	data := url.Values{}
	data.Set("refresh_token", refreshToken)
	data.Set("client_id", p.clientID)
	data.Set("client_secret", p.clientSecret)
	data.Set("grant_type", "refresh_token")

	log.Info().Str("url", googleTokenEndpoint).Msg("refreshing token")

	req, err := http.NewRequestWithContext(ctx, "POST", googleTokenEndpoint, strings.NewReader(data.Encode()))
	if err != nil {
		log.Error().Err(err).Msg("failed to create request")
		return auth.TokenResponse{}, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("failed to execute request")
		return auth.TokenResponse{}, fmt.Errorf("failed to refresh token: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("failed to read response body")
		return auth.TokenResponse{}, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Error().Str("status", resp.Status).Str("body", string(body)).Msg("token refresh failed")
		return auth.TokenResponse{}, fmt.Errorf("token refresh failed: %s - %s", resp.Status, string(body))
	}

	var tokenResp tokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		log.Error().Err(err).Msg("failed to parse response")
		return auth.TokenResponse{}, fmt.Errorf("failed to parse token response: %w", err)
	}

	log.Info().Msg("successfully refreshed token")

	return auth.TokenResponse{
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		ExpiresAt:    time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second).Format(time.RFC3339),
	}, nil
}
