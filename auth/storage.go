package auth

import (
	"DataLake/internal/logger"
	"encoding/json"
	"os"
)

// TokenStorage интерфейс для хранения токенов
type TokenStorage interface {
	SaveToken(providerName string, token TokenResponse) error
	LoadToken(providerName string) (TokenResponse, error)
}

type FileTokenStorage struct {
	filepath string
}

// NewFileTokenStorage создает новое хранилище токенов в файле
func NewFileTokenStorage(filepath string) *FileTokenStorage {
	return &FileTokenStorage{filepath: filepath}
}

// SaveToken сохраняет токен для указанного провайдера
func (s *FileTokenStorage) SaveToken(providerName string, token TokenResponse) error {
	log := logger.Get()

	tokens, _ := s.loadAll()
	if tokens == nil {
		tokens = make(map[string]TokenResponse)
	}

	tokens[providerName] = token

	data, err := json.MarshalIndent(tokens, "", "  ")
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal tokens")
		return err
	}

	err = os.WriteFile(s.filepath, data, 0644)
	if err != nil {
		log.Error().Err(err).Msg("failed to write tokens to file")
		return err
	}

	log.Info().Str("provider", providerName).Msg("tokens saved successfully")
	return nil
}

// LoadToken загружает токен для указанного провайдера
func (s *FileTokenStorage) LoadToken(providerName string) (TokenResponse, error) {
	log := logger.Get()

	tokens, err := s.loadAll()
	if err != nil {
		log.Error().Err(err).Msg("failed to read tokens from file")
		return TokenResponse{}, err
	}

	token, exists := tokens[providerName]
	if !exists {
		log.Warn().Str("provider", providerName).Msg("token not found for provider")
		return TokenResponse{}, os.ErrNotExist
	}

	log.Info().Str("provider", providerName).Msg("tokens loaded successfully")
	return token, nil
}

// loadAll загружает все токены из файла
func (s *FileTokenStorage) loadAll() (map[string]TokenResponse, error) {
	data, err := os.ReadFile(s.filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]TokenResponse), nil
		}
		return nil, err
	}

	var tokens map[string]TokenResponse
	err = json.Unmarshal(data, &tokens)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}
