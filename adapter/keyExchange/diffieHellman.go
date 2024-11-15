package keyExchange

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net"

	"github.com/siddeshwarnavink/UTA/adapter/proxy"

	"golang.org/x/crypto/curve25519"
)

func ClientDiffieHellman(encryptedConn net.Conn) ([]byte, error) {
	return diffieHellman(encryptedConn)
}

func ServerDiffieHellman(encryptedConn net.Conn) ([]byte, error) {
	return diffieHellman(encryptedConn)
}

func diffieHellman(encryptedConn net.Conn) ([]byte, error) {
	var privateKey, publicKey [32]byte

	// private key
	_, err := rand.Read(privateKey[:])
	if err != nil {
		fmt.Println("Error generating private key:", err)
		return nil, err
	}

	// public key
	curve25519.ScalarBaseMult(&publicKey, &privateKey)

	// send adapter public key
	jsonMsg := proxy.GenerateAdapterMessage(base64.StdEncoding.EncodeToString(publicKey[:]))

	_, err = encryptedConn.Write(jsonMsg)
	if err != nil {
		fmt.Println("Error sending public key:", err)
		return nil, err
	}

	buff := make([]byte, 1080)
	n, err := encryptedConn.Read(buff)
	if err != nil {
		fmt.Println("Error receiving public key:", err)
		return nil, err
	}

	buff = buff[:n]
	payload, err := proxy.GetAdapterMessage(buff)
	if err != nil {
		fmt.Printf("Invalid adapter message: %s", string(buff))
		return nil, errors.New("Invalid adapter message")
	}

	otherPublicKey, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		fmt.Println("Failed to decode base64 public key: " + err.Error())
		return nil, err
	}

	// get shared key
	sharedSecret, _ := curve25519.X25519(privateKey[:], otherPublicKey)

	return sharedSecret[:], nil
}
