package crypto

import (
	"crypto/rand"
	"fmt"
	"golang.org/x/crypto/chacha20"
	"io"
)

// EncryptChaCha encrypts data using ChaCha20.
func EncryptChaCha(data []byte, key []byte) ([]byte, error) {
	var nonce [12]byte
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		return nil, err
	}
	c, err := chacha20.NewUnauthenticatedCipher(key, nonce[:])
	if err != nil {
		return nil, fmt.Errorf("ChaCha20 initialization error: %v", err)
	}
	ciphertext := make([]byte, len(data))
	c.XORKeyStream(ciphertext, data)
	return append(nonce[:], ciphertext...), nil
}

// DecryptChaCha decrypts data using ChaCha20.
func DecryptChaCha(data []byte, key []byte) ([]byte, error) {
	if len(data) < 12 {
		return nil, fmt.Errorf("data too short for nonce")
	}
	var nonce [12]byte
	copy(nonce[:], data[:12])
	ciphertext := data[12:]
	c, err := chacha20.NewUnauthenticatedCipher(key, nonce[:])
	if err != nil {
		return nil, fmt.Errorf("ChaCha20 initialization error: %v", err)
	}
	plaintext := make([]byte, len(ciphertext))
	c.XORKeyStream(plaintext, ciphertext)
	return plaintext, nil
}
