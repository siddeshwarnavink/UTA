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

	if embeded.currentFlags.Mode == embeded.Client {
		fmt.Println("Mode: Client")
	} else {
		fmt.Println("Mode: Server")
	}
	// if embeded.currentFlags.Mode == nil {
	// 	for _, i := range embeded.UIQuestionList {
	// 		if i.Name == "MODE" {
	// 			modeResult, err := i.RenderFunc(i.Question, i.Options, i.PlaceHolder)
	// 			if err != nil {
	// 				return embeded.currentFlags, err
	// 			}
	// 			fmt.Println("---")
	// 			embeded.currentFlags.Mode = ModeFromString(modeResult)
	// 		}
	// 	}
	// }

	// if embeded.currentFlags.Dec == "" {
	// 	for _, i := range embeded.UIQuestionList {
	// 		if i.Name == "UNENCRYPTED_ADDRESS" {
	// 			decResult, err := i.RenderFunc(i.Question, i.Options, i.PlaceHolder)
	// 			if err != nil {
	// 				return embeded.currentFlags, err
	// 			}
	// 			fmt.Println("---")
	// 			embeded.currentFlags.Dec = decResult
	// 		}
	// 	}
	// }

	// if embeded.currentFlags.Enc == "" {
	// 	for _, i := range embeded.UIQuestionList {
	// 		if i.Name == "ENCRYPTED_ADDRESS" {
	// 			encResult, err := i.RenderFunc(i.Question, i.Options, i.PlaceHolder)
	// 			if err != nil {
	// 				return embeded.currentFlags, err
	// 			}
	// 			fmt.Println("---")
	// 			embeded.currentFlags.Enc = encResult
	// 		}
	// 	}
	// }

	// if embeded.currentFlags.Protocol == "" {
	// 	for _, i := range embeded.UIQuestionList {
	// 		if i.Name == "KEY_EXCHANGE" {
	// 			keyProtoResult, err := i.RenderFunc(i.Question, i.Options, i.PlaceHolder)
	// 			if err != nil {
	// 				return embeded.currentFlags, err
	// 			}
	// 			fmt.Println("---")
	// 			embeded.currentFlags.Protocol = keyProtoResult
	// 		}
	// 	}
	// }

	// if embeded.currentFlags.Algo == "" {
	// 	for _, i := range embeded.UIQuestionList {
	// 		if i.Name == "ENCRYPTION" {
	// 			algoResult, err := i.RenderFunc(i.Question, i.Options, i.PlaceHolder)
	// 			if err != nil {
	// 				return embeded.currentFlags, err
	// 			}
	// 			fmt.Println("---")
	// 			embeded.currentFlags.Algo = algoResult
	// 		}
	// 	}
	// }

}
