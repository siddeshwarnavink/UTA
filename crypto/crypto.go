package crypto

import (
	"fmt"
)

type Algorithm string

const (
	AlgoAES     Algorithm = "AES"
	AlgoChaCha  Algorithm = "ChaCha20"
	AlgoTwoFish Algorithm = "TwoFish"
)

// Encrypt encrypts data based on the chosen algorithm.
func Encrypt(data []byte, key []byte, algoName Algorithm) ([]byte, error) {

	switch algoName {
	case AlgoAES:
		return EncryptAES(data, key)

	case AlgoChaCha:
		return EncryptChaCha(data, key)

	case AlgoTwoFish:
		return EncryptTwoFish(data, key)

	default:
		return nil, fmt.Errorf("unsupported algorithm: %s", algoName)
	}
}

// Decrypt decrypts data based on the chosen algorithm.
func Decrypt(data []byte, key []byte, algoName Algorithm) ([]byte, error) {

	switch algoName {
	case AlgoAES:
		return DecryptAES(data, key)

	case AlgoChaCha:
		return DecryptChaCha(data, key)

	case AlgoTwoFish:
		return DecryptTwoFish(data, key)

	default:
		return nil, fmt.Errorf("unsupported algorithm: %s", algoName)
	}
}
