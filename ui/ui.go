package ui

import (
	"errors"
	"fmt"
	"os"

	"github.com/siddeshwarnavink/UTA/keyExchange"
)

func ParseFlags() (*Flags, error) {
	args := os.Args[1:]

	var mode AdapterMode
	var enc, dec string
	var algo string
	var protocol keyExchange.Protocol

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
				switch args[i+1] {
				case string(keyExchange.DiffieHellman):
					protocol = keyExchange.DiffieHellman
				}
				i++
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

	// Assuming RenderForm is a function that validates or processes the flags
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
	fmt.Println(Primary + "\033[1m" + "By Code Factort Unlimited" + "\033[1m" + Reset)

	if parsedFlags.Mode == "" {
		modeChan := make(chan string)
		go RenderModeForm(modeChan)
		modeResult := <-modeChan
		if modeResult == "error" {
			return parsedFlags, errors.New("mode not selected")
		}
		fmt.Println("---")
		parsedFlags.Mode = ModeFromString(modeResult)
	}

	if parsedFlags.Enc == "" && parsedFlags.Dec == "" {
		portChan := make(chan []string)
		go RenderPortForm(portChan)
		portResult := <-portChan
		if portResult[0] == "error" {
			return parsedFlags, errors.New("encrypted end's address not mentioned and unencrypted end's address not mentioned")
		}
		fmt.Println("---")
		parsedFlags.Enc = portResult[1]
		parsedFlags.Dec = portResult[0]
	}

	if parsedFlags.Protocol == "" {
		keyProtoChan := make(chan string)
		go RenderKeyProtoForm(keyProtoChan)
		keyProtoResult := <-keyProtoChan
		if keyProtoResult == "error" {
			return parsedFlags, errors.New("key Exchange Protocol not selected")
		}
		fmt.Println("---")
		parsedFlags.Protocol = KeyProtocolFromString(keyProtoResult)
	}

	if parsedFlags.Algo == "" {
		algoChan := make(chan string)
		go RenderAlgoForm(algoChan)
		algoResult := <-algoChan
		if algoResult == "error" {
			return parsedFlags, errors.New("algorithm not selected")
		}
		fmt.Println("---")
		parsedFlags.Algo = algoResult
	}
	return parsedFlags, nil
}
