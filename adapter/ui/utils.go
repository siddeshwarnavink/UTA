package ui

import (
	"errors"

	"github.com/siddeshwarnavink/UTA/adapter/embeded"
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
	Protocol string
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
	return nil, errors.New("crypto algorithm not found")
}

func KeyAlgorithmFromString(s string) (*embeded.KeyExchangeAlgo, error) {
	for _, protocol := range embeded.KeyExchangeList {
		if protocol.Name == s {
			return &protocol, nil
		}
	}
	return nil, errors.New("key exchange algorithm not found")
}
