package ui

import (
	"github.com/siddeshwarnavink/UTA/crypto"
	"github.com/siddeshwarnavink/UTA/keyExchange"
)

type AdapterMode string

const (
	Client AdapterMode = "Client"
	Server AdapterMode = "Server"
)

type Flags struct {
	Mode     AdapterMode
	Enc      string
	Dec      string
	Algo     crypto.Algorithm
	Protocol keyExchange.Protocol
}

func ModeFromString(s string) AdapterMode {
	switch s {
	case "Client":
		return Client
	case "Server":
		return Server
	default:
		return Client
	}
}

func AlgorithmFromString(s string) crypto.Algorithm {
	switch s {
	case "Advanced Encryption Standard(AES)":
		return crypto.AlgoAES
	case "ChaCha20":
		return crypto.AlgoChaCha
	case "TwoFish":
		return crypto.AlgoTwoFish
	default:
		return crypto.AlgoAES
	}
}

func KeyProtocolFromString(s string) keyExchange.Protocol {
	switch s {
	case "Diffie Hellman Key Exchange":
		return keyExchange.DiffieHellman
	default:
		return keyExchange.DiffieHellman
	}
}
