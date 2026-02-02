package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256" // Naya import: Hashing ke liye
	"fmt"
	"io"
	"os"
)

// Password ko 32-byte key mein convert karne ke liye
func createKey(password string) []byte {
	hash := sha256.Sum256([]byte(password))
	return hash[:]
}

func encryptFile(filename string, password string) error {
	plaintext, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	key := createKey(password) // Password se key banayi
	block, _ := aes.NewCipher(key)
	gcm, _ := cipher.NewGCM(block)

	nonce := make([]byte, gcm.NonceSize())
	io.ReadFull(rand.Reader, nonce)

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return os.WriteFile(filename+".aegis", ciphertext, 0644)
}

func decryptFile(filename string, password string) error {
	ciphertext, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	key := createKey(password) // Wahi password se key generate ki
	block, _ := aes.NewCipher(key)
	gcm, _ := cipher.NewGCM(block)

	nonceSize := gcm.NonceSize()
	nonce, actualCiphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, actualCiphertext, nil)
	if err != nil {
		return fmt.Errorf("galat password ya file kharab hai")
	}

	// Extension hatane ke liye
	outputName := "unlocked_" + filename[:len(filename)-6]
	return os.WriteFile(outputName, plaintext, 0644)
}
