package auth

import (
	"DataLake/internal/logger"
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

// TokenStorage интерфейс для хранения токенов
type TokenStorage interface {
	SaveToken(providerName string, token TokenResponse) error
	LoadToken(providerName string) (TokenResponse, error)
}

type FileTokenStorage struct {
	filepath   string
	encryption *Encryption
}

// NewFileTokenStorage создает новое хранилище токенов в файле с шифрованием
// encryptionKey должен быть 32 байта для AES-256
func NewFileTokenStorage(filepath string, encryptionKey string) (*FileTokenStorage, error) {
	encryption, err := NewEncryption(encryptionKey)
	if err != nil {
		return nil, err
	}
	return &FileTokenStorage{
		filepath:   filepath,
		encryption: encryption,
	}, nil
}

// encryptedToken структура для хранения зашифрованных данных
type encryptedToken struct {
	EncryptedData string `json:"encrypted_data"`
}

// NewFileTokenStorageFromEnv создает хранилище токенов с ключом шифрования из переменной окружения
// Если ENCRYPTION_KEY не задан, возвращает ошибку
func NewFileTokenStorageFromEnv(filepath string) (*FileTokenStorage, error) {
	encryptionKey := os.Getenv("ENCRYPTION_KEY")
	if encryptionKey == "" {
		return nil, fmt.Errorf("ENCRYPTION_KEY environment variable is required for token encryption")
	}

	if len(encryptionKey) != 32 {
		return nil, fmt.Errorf("ENCRYPTION_KEY must be exactly 32 bytes for AES-256, current length: %d", len(encryptionKey))
	}

	return NewFileTokenStorage(filepath, encryptionKey)
}

// SaveToken сохраняет токен для указанного провайдера с шифрованием
func (s *FileTokenStorage) SaveToken(providerName string, token TokenResponse) error {
	log := logger.Get()

	tokens, _ := s.loadAll()
	if tokens == nil {
		tokens = make(map[string]TokenResponse)
	}

	tokens[providerName] = token

	// Сериализуем все токены в JSON
	data, err := json.Marshal(tokens)
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal tokens")
		return err
	}

	// Шифруем данные
	encryptedData, err := s.encryption.Encrypt(data)
	if err != nil {
		log.Error().Err(err).Msg("failed to encrypt tokens")
		return err
	}

	// Сохраняем зашифрованные данные
	encToken := encryptedToken{EncryptedData: encryptedData}
	encData, err := json.MarshalIndent(encToken, "", "  ")
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal encrypted token")
		return err
	}

	err = os.WriteFile(s.filepath, encData, 0600) // 0600 вместо 0644 для лучшей безопасности
	if err != nil {
		log.Error().Err(err).Msg("failed to write tokens to file")
		return err
	}

	log.Info().Str("provider", providerName).Msg("tokens saved successfully (encrypted)")
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

// loadAll загружает все токены из файла с расшифровкой
func (s *FileTokenStorage) loadAll() (map[string]TokenResponse, error) {
	data, err := os.ReadFile(s.filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]TokenResponse), nil
		}
		return nil, err
	}

	// Пытаемся загрузить как зашифрованные данные (новый формат)
	var encToken encryptedToken
	err = json.Unmarshal(data, &encToken)
	if err == nil && encToken.EncryptedData != "" {
		// Расшифровываем данные
		decryptedData, err := s.encryption.Decrypt(encToken.EncryptedData)
		if err != nil {
			return nil, errors.New("failed to decrypt tokens: " + err.Error())
		}

		var tokens map[string]TokenResponse
		err = json.Unmarshal(decryptedData, &tokens)
		if err != nil {
			return nil, err
		}

		return tokens, nil
	}

	// Если не получилось, пытаемся загрузить как незашифрованные (старый формат для миграции)
	var tokens map[string]TokenResponse
	err = json.Unmarshal(data, &tokens)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}
