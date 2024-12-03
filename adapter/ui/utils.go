package ui

import (
	"errors"

	"github.com/siddeshwarnavink/UTA/adapter/embeded"
)

func ModeFromString(s string) embeded.AdapterMode {
	switch s {
	case "Client":
		return embeded.Client
	case "Server":
		return embeded.Server
	default:
		return embeded.Client
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
