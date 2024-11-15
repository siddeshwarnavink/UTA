package embeded

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
	lua "github.com/yuin/gopher-lua"
)

func EncryptAES(key, text []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	ciphertext := make([]byte, aes.BlockSize+len(text))
	iv := ciphertext[:aes.BlockSize]
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], text)
	return hex.EncodeToString(ciphertext), nil
}

func encryptAES(L *lua.LState) int {
	key := []byte(L.ToString(1))
	text := []byte(L.ToString(2))
	result, err := EncryptAES(key, text)
	if err != nil {
		L.Push(lua.LString(err.Error()))
	} else {
		L.Push(lua.LString(result))
	}
	return 1
}

func DecryptAES(key []byte, encryptedText string) (string, error) {
    ciphertext, err := hex.DecodeString(encryptedText)
    if err != nil {
        return "", err
    }
    if len(ciphertext) < aes.BlockSize {
        return "", errors.New("ciphertext too short")
    }

    block, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }
    iv := ciphertext[:aes.BlockSize]
    ciphertext = ciphertext[aes.BlockSize:]
    stream := cipher.NewCFBDecrypter(block, iv)
    stream.XORKeyStream(ciphertext, ciphertext)
    return string(ciphertext), nil
}

func decryptAES(L *lua.LState) int {
    key := []byte(L.ToString(1))
    encryptedText := L.ToString(2)
    result, err := DecryptAES(key, encryptedText)
    if err != nil {
        L.Push(lua.LString(err.Error()))
    } else {
        L.Push(lua.LString(result))
    }
    return 1
}

func AlogAesLoader(l *lua.LState) int {
	var exports = map[string]lua.LGFunction{
		"encrypt": encryptAES,
		"decrypt": decryptAES,
	}

	mod := l.SetFuncs(l.NewTable(), exports)
	l.Push(mod)
	return 1
}
