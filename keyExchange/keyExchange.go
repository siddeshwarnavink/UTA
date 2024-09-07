package keyExchange

import (
	"fmt"
	"net"
)

type Protocol string

const (
	DiffieHellman Protocol = "dhkc"
)

func ServerKeyExchange(encryptedConn net.Conn, protocolName Protocol) ([]byte, error) {

	switch protocolName {
	case DiffieHellman:
		return ServerDiffieHellman(encryptedConn)

	default:
		return nil, fmt.Errorf("unsupported algorithm: %s", protocolName)
	}

}

func ClientKeyExchange(encryptedConn net.Conn, protocolName Protocol) ([]byte, error) {

	switch protocolName {
	case DiffieHellman:
		return ClientDiffieHellman(encryptedConn)

	default:
		return nil, fmt.Errorf("unsupported algorithm: %s", protocolName)
	}

}
