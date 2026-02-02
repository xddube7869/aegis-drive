package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"os"
)

// Encrypt function: File ko lock karne ke liye
func encryptFile(filename string, key []byte) error {
	plaintext, _ := os.ReadFile(filename)
	
	block, _ := aes.NewCipher(key)
	gcm, _ := cipher.NewGCM(block)
	
	nonce := make([]byte, gcm.NonceSize())
	io.ReadFull(rand.Reader, nonce)
	
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	
	// Nayi encrypted file banana (.aegis extension ke sath)
	return os.WriteFile(filename+".aegis", ciphertext, 0644)
}

// Decrypt function: File ko wapas asli roop mein laane ke liye
func decryptFile(filename string, key []byte) error {
	ciphertext, _ := os.ReadFile(filename)
	
	block, _ := aes.NewCipher(key)
	gcm, _ := cipher.NewGCM(block)
	
	nonceSize := gcm.NonceSize()
	nonce, actualCiphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	
	plaintext, _ := gcm.Open(nil, nonce, actualCiphertext, nil)
	
	// Asli file wapas nikalna (bin_ prefix ke sath testing ke liye)
	return os.WriteFile("decrypted_"+filename[:len(filename)-6], plaintext, 0644)
}

func main() {
	// Ye tumhara secret 32-byte key hai (Password)
	// Real project mein hum ise user se input lenge
	key := []byte("thisisa32bitsecretkeyforproject!") 

	fmt.Println("🔒 Aegis Engine Starting...")
	
	// Testing ke liye: Ek file ko encrypt karo
	// Pehle apne folder mein 'test.txt' bana lena
	err := encryptFile("test.txt", key)
	if err == nil {
		fmt.Println("✅ File Encrypted: test.txt.aegis")
	} else {
		fmt.Println("❌ Error:", err)
	}
}
