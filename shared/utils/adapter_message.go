package utils

import "fmt"

func GetAdapterMessage(buf []byte) (string, error) {
	header := buf[0]
	if header == 0x18 {
		return string(buf[1:]), nil
	}
	return "", fmt.Errorf("invalid adapter message header: %x", header)
}

func GenerateAdapterMessage(message string) []byte {
	buf := make([]byte, len(message)+1)
	buf[0] = 0x18
	copy(buf[1:], message)
	return buf
}
