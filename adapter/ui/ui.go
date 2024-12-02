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
		for _, i := range embeded.UIQuestionList {
			if i.Name == "MODE" {
				modeResult, err := i.RenderFunc(i.Question, i.Options, i.PlaceHolder)
				if err != nil {
					return parsedFlags, err
				}
				fmt.Println("---")
				parsedFlags.Mode = ModeFromString(modeResult)
			}
		}
	}

	if parsedFlags.Dec == "" {
		for _, i := range embeded.UIQuestionList {
			if i.Name == "UNENCRYPTED_ADDRESS" {
				decResult, err := i.RenderFunc(i.Question, i.Options, i.PlaceHolder)
				if err != nil {
					return parsedFlags, err
				}
				fmt.Println("---")
				parsedFlags.Dec = decResult
			}
		}
	}

	if parsedFlags.Enc == "" {
		for _, i := range embeded.UIQuestionList {
			if i.Name == "ENCRYPTED_ADDRESS" {
				encResult, err := i.RenderFunc(i.Question, i.Options, i.PlaceHolder)
				if err != nil {
					return parsedFlags, err
				}
				fmt.Println("---")
				parsedFlags.Enc = encResult
			}
		}
	}

	if parsedFlags.Protocol == "" {
		for _, i := range embeded.UIQuestionList {
			if i.Name == "KEY_EXCHANGE" {
				keyProtoResult, err := i.RenderFunc(i.Question, i.Options, i.PlaceHolder)
				if err != nil {
					return parsedFlags, err
				}
				fmt.Println("---")
				parsedFlags.Protocol = keyProtoResult
			}
		}
	}

	if parsedFlags.Algo == "" {
		for _, i := range embeded.UIQuestionList {
			if i.Name == "ENCRYPTION" {
				algoResult, err := i.RenderFunc(i.Question, i.Options, i.PlaceHolder)
				if err != nil {
					return parsedFlags, err
				}
				fmt.Println("---")
				parsedFlags.Algo = algoResult
			}
		}
	}

	return parsedFlags, nil
}
