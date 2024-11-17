package embeded

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net"

	"github.com/siddeshwarnavink/UTA/shared/utils"
	lua "github.com/yuin/gopher-lua"
	"golang.org/x/crypto/curve25519"
)

func ClientDiffieHellman(encryptedConn net.Conn) (string, error) {
	return diffieHellman(encryptedConn)
}

func ServerDiffieHellman(encryptedConn net.Conn) (string, error) {
	return diffieHellman(encryptedConn)
}

func diffieHellman(encryptedConn net.Conn) (string, error) {
	var privateKey, publicKey [32]byte

	// Generate a random private key
	_, err := rand.Read(privateKey[:])
	if err != nil {
		fmt.Println("Error generating private key:", err)
		return "", err
	}

	// Generate the public key from the private key using curve25519
	curve25519.ScalarBaseMult(&publicKey, &privateKey)

	// Send public key to the other party
	jsonMsg := utils.GenerateAdapterMessage(base64.StdEncoding.EncodeToString(publicKey[:]))

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
	payload, err := utils.GetAdapterMessage(buff)
	if err != nil {
		fmt.Printf("Invalid adapter message: %s", string(buff))
		return "", errors.New("invalid adapter message")
	}

	otherPublicKey, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		fmt.Println("Failed to decode base64 public key: " + err.Error())
		return "", err
	}

	sharedSecret, _ := curve25519.X25519(privateKey[:], otherPublicKey)

	return base64.StdEncoding.EncodeToString(sharedSecret[:]), nil
}

// Lua bindings for the Diffie-Hellman functions

func ClientDiffieHellmanLua(L *lua.LState) int {
	encConn := L.ToUserData(1).Value.(net.Conn)

	sharedSecret, err := ClientDiffieHellman(encConn)
	if err != nil {
		L.Push(lua.LString(err.Error()))
	} else {
		L.Push(lua.LString(sharedSecret))
	}
	return 1
}

func ServerDiffieHellmanLua(L *lua.LState) int {
	encConn := L.ToUserData(1).Value.(net.Conn)

	sharedSecret, err := ServerDiffieHellman(encConn)
	if err != nil {
		L.Push(lua.LString(err.Error()))
	} else {
		L.Push(lua.LString(sharedSecret))
	}
	return 1
}

func DiffieHellmanLoader(l *lua.LState) int {
	var exports = map[string]lua.LGFunction{
		"clientDiffieHellman": ClientDiffieHellmanLua,
		"serverDiffieHellman": ServerDiffieHellmanLua,
	}

	mod := l.SetFuncs(l.NewTable(), exports)
	l.Push(mod)
	return 1
}
