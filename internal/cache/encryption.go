package cache

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const (
	// KeySize is the size of the encryption key in bytes (AES-256)
	KeySize = 32

	// KeyFileName is the name of the encryption key file
	KeyFileName = ".key"
)

// GenerateKey generates a new random encryption key
func GenerateKey() ([]byte, error) {
	key := make([]byte, KeySize)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, fmt.Errorf("failed to generate key: %w", err)
	}
	return key, nil
}

// SaveKey saves the encryption key to a file with restricted permissions
func SaveKey(configDir string, key []byte) error {
	keyPath := filepath.Join(configDir, KeyFileName)

	// Create file with restricted permissions (600 = owner read/write only)
	if err := os.WriteFile(keyPath, key, 0600); err != nil {
		return fmt.Errorf("failed to save key: %w", err)
	}

	return nil
}

// LoadKey loads the encryption key from file
func LoadKey(configDir string) ([]byte, error) {
	keyPath := filepath.Join(configDir, KeyFileName)

	key, err := os.ReadFile(keyPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("encryption key not found (run 'sreq init' to create)")
		}
		return nil, fmt.Errorf("failed to read key: %w", err)
	}

	if len(key) != KeySize {
		return nil, fmt.Errorf("invalid key size: expected %d, got %d", KeySize, len(key))
	}

	return key, nil
}

// KeyExists checks if the encryption key file exists
func KeyExists(configDir string) bool {
	keyPath := filepath.Join(configDir, KeyFileName)
	_, err := os.Stat(keyPath)
	return err == nil
}

// Encrypt encrypts data using AES-256-GCM
func Encrypt(key, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Create nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt and prepend nonce
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// Decrypt decrypts data using AES-256-GCM
func Decrypt(key, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("decryption failed: %w", err)
	}

	return plaintext, nil
}
