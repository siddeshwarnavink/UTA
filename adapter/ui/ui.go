package ui

import (
	"errors"
	"fmt"
	"os"

	"github.com/siddeshwarnavink/UTA/adapter/embeded"
)

func ParseFlags() (*Flags, error) {
	args := os.Args[1:]

	var mode AdapterMode
	var enc, dec string
	var algo string
	var protocol string

	i := 0
	for i < len(args) {
		arg := args[i]
		switch arg {
		case "--client":
			mode = Client
		case "--server":
			mode = Server
		case "-enc":
			if i+1 < len(args) {
				enc = args[i+1]
				i++
			} else if dec != "" {
				return nil, errors.New("missing value for -enc")
			}
		case "-dec":
			if i+1 < len(args) {
				dec = args[i+1]
				i++
			} else if enc != "" {
				return nil, errors.New("missing value for -dec")
			}
		case "--algo":
			if i+1 < len(args) {
				algo = args[i+1]
				i++
			} else if enc != "" {
				return nil, errors.New("missing value for -algo")
			}
		case "--prot":
			if i+1 < len(args) {
				protocol = args[i+1]
				i++
			} else if enc != "" {
				return nil, errors.New("missing value for -prot")
			}
		default:
			return nil, fmt.Errorf("unknown flag: %s", arg)
		}
		i++
	}

	flags := Flags{
		Mode:     mode,
		Enc:      enc,
		Dec:      dec,
		Algo:     algo,
		Protocol: protocol,
	}

	finalFlag, err := RenderForm(flags)
	if err != nil {
		return nil, err
	}

	return &finalFlag, nil
}

func RenderForm(parsedFlags Flags) (Flags, error) {
	const (
		Primary = "\033[38;5;205m"
		Reset   = "\033[0m"
	)

	fmt.Println(Primary +
		`
  __  ___________
 / / / /_  __/ _ |
/ /_/ / / / / __ |
\____/ /_/ /_/ |_|
` + Reset)
	fmt.Println(Primary + "\033[1m" + "By Code Factory Unlimited" + "\033[1m" + Reset)

	if parsedFlags.Mode == "" {
		modeResult, err := MCQ("Which mode is this system on?", []string{"Client", "Server"})
		if err != nil {
			return parsedFlags, err
		}
		fmt.Println("---")
		parsedFlags.Mode = ModeFromString(modeResult)
	}

	if parsedFlags.Enc == "" && parsedFlags.Dec == "" {
		portResult, err := Form("Enter the Connection Addresses", []string{"Unencrypted Connection's Address", "Encrypted Connection's Address"})
		if err != nil {
			return parsedFlags, err
		}
		fmt.Println("---")
		parsedFlags.Dec = portResult[0]
		parsedFlags.Enc = portResult[1]
	}

	if parsedFlags.Protocol == "" {
		var KeyExchangeProtocol = []string{}
		for _, entry := range embeded.KeyExchangeList {
			KeyExchangeProtocol = append(KeyExchangeProtocol, entry.Name)
		}
		keyProtoResult, err := MCQ("Which Key Exchange Protocol is being used?", KeyExchangeProtocol)
		if err != nil {
			return parsedFlags, err
		}
		fmt.Println("---")
		parsedFlags.Protocol = keyProtoResult
	}

	if parsedFlags.Algo == "" {
		var Algorithms = []string{}
		for _, entry := range embeded.CryptoList {
			Algorithms = append(Algorithms, entry.Name)
		}
		algoResult, err := MCQ("Which Cryptographic Algorithm to be used?", Algorithms)
		if err != nil {
			return parsedFlags, err
		}
		fmt.Println("---")
		parsedFlags.Algo = algoResult
	}
	return parsedFlags, nil
}
