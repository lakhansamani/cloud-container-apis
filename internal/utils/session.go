package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"strings"

	"github.com/lakhansamani/cloud-container/internal/db/models"
)

const (
	// EncryptionKey is the key used for encryption
	EncryptionKey = "cc-encryption-key-2024-13-05"
)

var bytes = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 0o5}

// EncryptB64 encrypts data into base64 string
func EncryptB64(text string) string {
	return base64.StdEncoding.EncodeToString([]byte(text))
}

// DecryptB64 decrypts from base64 string to readable string
func DecryptB64(s string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// EncryptAES method is to encrypt or hide any classified text
func EncryptAES(text string) (string, error) {
	key := []byte(EncryptionKey)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	plainText := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, bytes)
	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)
	return EncryptB64(string(cipherText)), nil
}

// DecryptAES method is to extract back the encrypted text
func DecryptAES(text string) (string, error) {
	key := []byte(EncryptionKey)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	cipherText, err := DecryptB64(text)
	if err != nil {
		return "", err
	}
	cfb := cipher.NewCFBDecrypter(block, bytes)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, []byte(cipherText))
	return string(plainText), nil
}

// GenerateSession token using user id and email address from user object
// Encrypt this token before sending to client
func GenerateSession(user *models.User, nonce string) (string, error) {
	token := user.ID + ":" + nonce
	// Encrypt token using AES encryption
	encryptedToken, err := EncryptAES(token)
	if err != nil {
		return "", err
	}
	return encryptedToken, nil
}

// DecryptSession token using user id and email address from user object
// Decrypt this token before sending to client
func DecryptSession(token string) (string, string, error) {
	// Decrypt token using AES encryption
	decryptedToken, err := DecryptAES(token)
	if err != nil {
		return "", "", err
	}
	// Split session
	splitSession := strings.Split(decryptedToken, ":")
	if len(splitSession) != 2 {
		return "", "", errors.New("invalid session")
	}
	// Get user info from session
	userID := splitSession[0]
	nonce := splitSession[2]
	return userID, nonce, nil
}
