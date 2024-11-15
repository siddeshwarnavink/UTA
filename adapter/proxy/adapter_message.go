package proxy

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type AdapterMessage struct {
	To        string `json:"to"`
	Payload   string `json:"payload"`
	Timestamp int64  `json:"timestamp"`
}

func GetAdapterMessage(buf []byte) (string, error) {

	var obj map[string]interface{}
	err := json.Unmarshal(buf, &obj)

	if err != nil {
		return "not json", err
	}

	val, ok := obj["to"]
	if !ok || val != "adapter" {
		return "", errors.New("not to adapter")
	}

	payload, ok := obj["payload"].(string)
	if !ok {
		return "", errors.New("no payload")
	}
	return payload, nil
}

func GenerateAdapterMessage(payload string) []byte {
	unixTimestamp := time.Now().Unix()
	message := AdapterMessage{
		To:        "adapter",
		Timestamp: unixTimestamp,
		Payload:   payload,
	}

	responseJson, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Error:", err)
	}

	return responseJson
}
