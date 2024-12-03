package ui

import (
	"errors"
	"fmt"
	"os"

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

func GetConfigFile() string {
	args := os.Args[1:]
	var filePath string
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "--config":
			if i+1 < len(args) {
				filePath = args[i+1]
				fmt.Println(filePath)
				i++
			} else {
				fmt.Println("Please provide a config file path")
				os.Exit(1)
			}
		default:
			continue
		}
	}
	return filePath
}
