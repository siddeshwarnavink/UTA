package embeded

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/siddeshwarnavink/UTA/shared/utils"
	lua "github.com/yuin/gopher-lua"
)

type KeyHolder struct {
	key       *rsa.PrivateKey
	createdAt time.Time
}

var currentKey KeyHolder

func getKey() (*rsa.PrivateKey, error) {
	now := time.Now()
	twoDaysAgo := now.Add(-time.Hour * 2 * 24)

	if (currentKey == KeyHolder{}) || currentKey.createdAt.Before(twoDaysAgo) {
		temp, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			return nil, err
		}
		currentKey.createdAt = now
		currentKey.key = temp
	}
	return currentKey.key, nil
}

func ClientRSA(encryptedConn net.Conn) (string, error) {
	return rsaKeyExchange(encryptedConn)
}

func ServerRSA(encryptedConn net.Conn) (string, error) {
	return rsaKeyExchange(encryptedConn)
}

func rsaKeyExchange(encryptedConn net.Conn) (string, error) {
	// Generate private RSA key
	key, err := getKey()
	if err != nil {
		fmt.Println("Error generating RSA key:", err)
		return "", err
	}

	// Send public key
	pemBlock := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(&key.PublicKey),
	}
	publicKeyPem := pem.EncodeToMemory(pemBlock)
	jsonMsg := utils.GenerateAdapterMessage(base64.StdEncoding.EncodeToString(publicKeyPem))

	_, err = encryptedConn.Write(jsonMsg)
	if err != nil {
		fmt.Println("Error sending public key:", err)
		return "", err
	}

	// Receive the other party's public key
	buff := make([]byte, 1080)
	n, err := encryptedConn.Read(buff)
	if err != nil {
		fmt.Println("Error receiving public key:", err)
		return "", err
	}
	buff = buff[:n]

	// Extract the payload and decode the session key
	payload, err := utils.GetAdapterMessage(buff)
	if err != nil {
		fmt.Printf("Invalid adapter message: %s", string(buff))
		return "", errors.New("Invalid adapter message")
	}

	sessionKeyPem, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		fmt.Println("Failed to decode base64 public key: " + err.Error())
		return "", err
	}

	block, _ := pem.Decode(sessionKeyPem)
	if block == nil || block.Type != "SESSION KEY" {
		fmt.Println("Failed to decode session key")
		return "", errors.New("Invalid session key")
	}

	// Decrypt the session key with the RSA private key
	encryptedSessionKey := block.Bytes
	sessionKey, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, key, encryptedSessionKey, nil)
	if err != nil {
		fmt.Println("Error decrypting session key:", err)
		return "", err
	}

	return base64.StdEncoding.EncodeToString(sessionKey), nil
}

// Lua bindings for the RSA key exchange functions

func ClientRSALua(L *lua.LState) int {
	encConn := L.ToUserData(1).Value.(net.Conn)

	sessionKey, err := ClientRSA(encConn)
	if err != nil {
		L.Push(lua.LString(err.Error()))
	} else {
		L.Push(lua.LString(sessionKey))
	}
	return 1
}

func ServerRSALua(L *lua.LState) int {
	encConn := L.ToUserData(1).Value.(net.Conn)

	sessionKey, err := ServerRSA(encConn)
	if err != nil {
		L.Push(lua.LString(err.Error()))
	} else {
		L.Push(lua.LString(sessionKey))
	}
	return 1
}

func RSAKeyExchangeLoader(L *lua.LState) int {
	var exports = map[string]lua.LGFunction{
		"clientRSA": ClientRSALua,
		"serverRSA": ServerRSALua,
	}

	mod := L.SetFuncs(L.NewTable(), exports)
	L.Push(mod)
	return 1
}
