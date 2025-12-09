package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

// Encryption предоставляет методы для шифрования и дешифрования данных
type Encryption struct {
	key []byte
}

// NewEncryption создает новый экземпляр Encryption с заданным ключом
// Ключ должен быть 16, 24 или 32 байта для AES-128, AES-192 или AES-256
func NewEncryption(key string) (*Encryption, error) {
	keyBytes := []byte(key)
	if len(keyBytes) != 16 && len(keyBytes) != 24 && len(keyBytes) != 32 {
		return nil, errors.New("encryption key must be 16, 24, or 32 bytes")
	}
	return &Encryption{key: keyBytes}, nil
}

// Encrypt шифрует данные с использованием AES-GCM
func (e *Encryption) Encrypt(plaintext []byte) (string, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Генерируем случайный nonce
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Шифруем данные и добавляем nonce в начало
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)

	// Кодируем в base64 для безопасного хранения
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt расшифровывает данные, зашифрованные с помощью Encrypt
func (e *Encryption) Decrypt(ciphertext string) ([]byte, error) {
	// Декодируем из base64
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesGCM.NonceSize()
	if len(data) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	// Извлекаем nonce и зашифрованные данные
	nonce, ciphertextBytes := data[:nonceSize], data[nonceSize:]

	// Расшифровываем
	plaintext, err := aesGCM.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
