package crypto

import (
	"fmt"
	"golang.org/x/crypto/twofish"
)

// EncryptTwoFish encrypts data using Twofish.
func EncryptTwoFish(data []byte, key []byte) ([]byte, error) {
	block, err := twofish.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("Twofish encryption error: %v", err)
	}
	ciphertext := make([]byte, len(data))
	block.Encrypt(ciphertext, data)
	return ciphertext, nil
}

// DecryptTwoFish decrypts data using Twofish.
func DecryptTwoFish(data []byte, key []byte) ([]byte, error) {
	block, err := twofish.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("Twofish decryption error: %v", err)
	}
	plaintext := make([]byte, len(data))
	block.Decrypt(plaintext, data)
	return plaintext, nil
}
