package p2p

import "encoding/json"

type RequestMessageType int

const (
	RequestMessageTypeInvalid RequestMessageType = -1
	RequestTypeConfig         RequestMessageType = 0
	RequestTypeLogs           RequestMessageType = 1
	RequestTypeSaveConfig     RequestMessageType = 2
)

type RequestMessageJson struct {
	Type      RequestMessageType `json:"t"`
	RequestId string             `json:"i"`
	Payload   string             `json:"p"`
}

func RequestMessage(role PeerRole, reqType RequestMessageType, requestId string, payload string) (string, error) {
	reqObj := RequestMessageJson{
		Type:      reqType,
		RequestId: requestId,
		Payload:   payload,
	}

	json, err := json.Marshal(reqObj)
	if err != nil {
		return "", err
	}

	return StringMessage(role, string(json)+"0")
}

func ExtractRequestMessage(message string) (PeerRole, RequestMessageType, string, string, error) {
	role, rawMsg, err := ExtractStringMessage(message)
	if err != nil {
		return InvalidRole, RequestMessageTypeInvalid, "", "", err
	}

	// TODO: Check if valid request message

	var obj RequestMessageJson
	err = json.Unmarshal([]byte(rawMsg[:len(rawMsg)-1]), &obj)
	if err != nil {
		return role, RequestMessageTypeInvalid, "", "", err
	}

	return role, obj.Type, obj.RequestId, obj.Payload, nil
}

type ResponseMessageJson struct {
	RequestId string `json:"i"`
	Data      string `json:"d"`
}

func ResponseMessage(role PeerRole, requestId string, data string) (string, error) {
	reqObj := ResponseMessageJson{
		RequestId: requestId,
		Data:      data,
	}

	json, err := json.Marshal(reqObj)
	if err != nil {
		return "", err
	}

	return StringMessage(role, string(json)+"1")
}

func ExtractResponseMessage(message string) (PeerRole, string, string, error) {
	role, rawMsg, err := ExtractStringMessage(message)
	if err != nil {
		return InvalidRole, "", "", err
	}

	// TODO: Check if valid response message

	var obj ResponseMessageJson
	err = json.Unmarshal([]byte(rawMsg[:len(rawMsg)-1]), &obj)
	if err != nil {
		return role, "", "", err
	}

	return role, obj.RequestId, obj.Data, nil
}
