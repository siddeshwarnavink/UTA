package embeded

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
	"net"

	"github.com/siddeshwarnavink/UTA/shared/utils"
	lua "github.com/yuin/gopher-lua"
)

func ClientECDH(encryptedConn net.Conn) ([]byte, error) {
	return ecdhKeyExchange(encryptedConn)
}

func ServerECDH(encryptedConn net.Conn) ([]byte, error) {
	return ecdhKeyExchange(encryptedConn)
}

func ecdhKeyExchange(encryptedConn net.Conn) ([]byte, error) {
	// Generate the ECDH key pair using P-256 curve
	curve := elliptic.P256()
	priv, x, y, err := generateECDHKey(curve)
	if err != nil {
		return nil, err
	}

	// Send the public key to the other party
	publicKey := elliptic.Marshal(curve, x, y)
	jsonMsg := utils.GenerateAdapterMessage(base64.StdEncoding.EncodeToString(publicKey))
	_, err = encryptedConn.Write(jsonMsg)
	if err != nil {
		return nil, fmt.Errorf("error sending public key: %w", err)
	}

	// Receive the other party's public key
	buff := make([]byte, 1080)
	n, err := encryptedConn.Read(buff)
	if err != nil {
		return nil, fmt.Errorf("error receiving public key: %w", err)
	}
	buff = buff[:n]

	// Parse the received public key
	payload, err := utils.GetAdapterMessage(buff)
	if err != nil {
		return nil, fmt.Errorf("invalid adapter message: %w", err)
	}

	otherPublicKey, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 public key: %w", err)
	}

	// Parse the received public key and calculate the shared secret
	x2, y2 := elliptic.Unmarshal(curve, otherPublicKey)
	if x2 == nil || y2 == nil {
		return nil, errors.New("invalid public key received")
	}

	// Compute the shared secret
	sharedX, _ := curve.ScalarMult(x2, y2, priv.D.Bytes())
	sharedSecret := sha256.Sum256(sharedX.Bytes())

	// Store the shared secret in the key exchange list
	for _, keyExchange := range KeyExchangeList {
		if keyExchange.Name == "ECDH" {
			keyExchange.Key = sharedSecret[:]
		}
	}

	return sharedSecret[:], nil
}

func generateECDHKey(curve elliptic.Curve) (*ecdsa.PrivateKey, *big.Int, *big.Int, error) {
	priv, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to generate ECDH private key: %w", err)
	}
	x, y := curve.ScalarBaseMult(priv.D.Bytes())
	return priv, x, y, nil
}

// Lua bindings for the ECDH functions

func ClientECDHLua(L *lua.LState) int {
	encConn := L.ToUserData(1).Value.(net.Conn)

	sharedSecret, err := ClientECDH(encConn)
	if err != nil {
		L.Push(lua.LString(err.Error()))
	} else {
		L.Push(lua.LString(base64.StdEncoding.EncodeToString(sharedSecret)))
	}
	return 1
}

func ServerECDHLua(L *lua.LState) int {
	encConn := L.ToUserData(1).Value.(net.Conn)

	sharedSecret, err := ServerECDH(encConn)
	if err != nil {
		L.Push(lua.LString(err.Error()))
	} else {
		L.Push(lua.LString(base64.StdEncoding.EncodeToString(sharedSecret)))
	}
	return 1
}

func ECDHLoader(l *lua.LState) int {
	var exports = map[string]lua.LGFunction{
		"clientECDH": ClientECDHLua,
		"serverECDH": ServerECDHLua,
	}

	mod := l.SetFuncs(l.NewTable(), exports)
	l.Push(mod)
	return 1
}
