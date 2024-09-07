package ui

import (
	"errors"

	"github.com/siddeshwarnavink/UTA/embeded"
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
	Algo     string
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

func AlgorithmFromString(name string) (*embeded.CryptoAlgo, error) {
	for _, algo := range embeded.CryptoList {
		if algo.Name == name {
			return &algo, nil
		}
	}
	return nil, errors.New("Crypto algorithm not found")
}

func KeyProtocolFromString(s string) keyExchange.Protocol {
	switch s {
	case "Diffie Hellman Key Exchange":
		return keyExchange.DiffieHellman
	default:
		return keyExchange.DiffieHellman
	}
}
