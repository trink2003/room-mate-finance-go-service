package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"room-mate-finance-go-service/constant"
)

func GenerateCheckSumUsingSha256New(input string) (string, error) {
	if input == "" {
		return constant.EmptyString, errors.New("input can not be empty")
	}
	h := sha256.New()
	h.Write([]byte(input))
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func GenerateCheckSumUsingSha256Sum256(input string) (string, error) {
	if input == "" {
		return constant.EmptyString, errors.New("input can not be empty")
	}
	sum := sha256.Sum256([]byte("this is a password"))
	return fmt.Sprintf("%x", sum), nil
}

func AesDecryption(input string, key string) (string, error) {
	if input == "" || key == "" {
		return constant.EmptyString, errors.New("input can not be empty")
	}
	if len(key) < 32 {
		return constant.EmptyString, errors.New("key length must be greater than 32 bytes")
	}
	gcm, nonce, err := generateAesNecessaryComponent([]byte(key))
	if err != nil {
		return constant.EmptyString, err
	}
	plaintext, decryptionError := gcm.Open(nil, nonce, []byte(input), nil)
	if decryptionError != nil {
		return constant.EmptyString, decryptionError
	}
	// return fmt.Sprintf("%s", plaintext), nil
	return string(plaintext), nil
}

func AesEncryption(input string, key string) (string, error) {
	if input == "" || key == "" {
		return constant.EmptyString, errors.New("input can not be empty")
	}
	if len(key) < 32 {
		return constant.EmptyString, errors.New("key length must be greater than 32 bytes")
	}

	gcm, nonce, err := generateAesNecessaryComponent([]byte(key))
	if err != nil {
		return constant.EmptyString, err
	}

	// Encrypt the input
	ciphertext := gcm.Seal(nonce, nonce, []byte(input), nil)

	return fmt.Sprintf("%x", ciphertext), nil
}

func generateAesNecessaryComponent(inputByte []byte) (cipher.AEAD, []byte, error) {
	// Generate a new AES cipher using our len(key)-byte long key
	cipherResult, cipherCreationError := aes.NewCipher(inputByte)

	if cipherCreationError != nil {
		return nil, []byte(constant.EmptyString), cipherCreationError
	}

	// Galois/Counter Mode (GCM) is a mode of operation for symmetric key cryptographic block ciphers
	gcm, gcmCreationError := cipher.NewGCM(cipherResult)

	if gcmCreationError != nil {
		return nil, []byte(constant.EmptyString), gcmCreationError
	}

	// Generate a random nonce
	nonce := make([]byte, gcm.NonceSize())

	if _, ioReadFullError := io.ReadFull(rand.Reader, nonce); ioReadFullError != nil {
		return nil, []byte(constant.EmptyString), ioReadFullError
	}

	return gcm, nonce, nil
}
