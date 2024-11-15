package proxy

import (
	"bytes"
)

func IsUninitialized(arr []byte) bool {
	return bytes.Equal(arr, make([]byte, len(arr)))
}
