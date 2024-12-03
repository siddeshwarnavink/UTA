package ui

import (
	"fmt"

	"github.com/siddeshwarnavink/UTA/adapter/embeded"
)

func RenderForm() {
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

	if embeded.CurrentFlags.Mode == "" {
		for _, i := range embeded.UIQuestionList {
			if i.Name == "MODE" {
				modeResult, err := i.RenderFunc(i.Question, i.Options, i.PlaceHolder)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println("---")
				embeded.CurrentFlags.Mode = ModeFromString(modeResult)
			}
		}
	}

	if embeded.CurrentFlags.Dec == "" {
		for _, i := range embeded.UIQuestionList {
			if i.Name == "UNENCRYPTED_ADDRESS" {
				decResult, err := i.RenderFunc(i.Question, i.Options, i.PlaceHolder)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println("---")
				embeded.CurrentFlags.Dec = decResult
			}
		}
	}

	if embeded.CurrentFlags.Enc == "" {
		for _, i := range embeded.UIQuestionList {
			if i.Name == "ENCRYPTED_ADDRESS" {
				encResult, err := i.RenderFunc(i.Question, i.Options, i.PlaceHolder)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println("---")
				embeded.CurrentFlags.Enc = encResult
			}
		}
	}

	if embeded.CurrentFlags.KeyAlgo == "" {
		for _, i := range embeded.UIQuestionList {
			if i.Name == "KEY_EXCHANGE" {
				keyProtoResult, err := i.RenderFunc(i.Question, i.Options, i.PlaceHolder)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println("---")
				embeded.CurrentFlags.KeyAlgo = keyProtoResult
			}
		}
	}

	if embeded.CurrentFlags.CryptoAlgo == "" {
		for _, i := range embeded.UIQuestionList {
			if i.Name == "ENCRYPTION" {
				algoResult, err := i.RenderFunc(i.Question, i.Options, i.PlaceHolder)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println("---")
				embeded.CurrentFlags.CryptoAlgo = algoResult
			}
		}
	}

}
