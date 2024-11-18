package utils

import "fmt"

//AdapterHeader is the header of the adapter message
//First 4 bytes of the message are the version
//Next bit is if the message is for adapter or not(1 for adapter, 0 for not)
//Next 3 bits are reserved for future use

// 0001 - version 1
// 1 - adapter message
// 000 - reserved
var AdapterHeader byte = 0b00011000

// 0001 - version 1
// 0 - data message
// 000 - reserved
var DataHeader byte = 0b00010000

func GetAdapterMessage(buf []byte) (string, error) {
	header := buf[0]
	if header == AdapterHeader {
		return string(buf[1:]), nil
	}
	return "", fmt.Errorf("invalid adapter message header: %x", header)
}

func GenerateAdapterMessage(message string) []byte {
	buf := make([]byte, len(message)+1)
	buf[0] = AdapterHeader
	copy(buf[1:], message)
	return buf
}

func GetDataMessage(buf []byte) (string, error) {
	header := buf[0]
	if header == DataHeader {
		return string(buf[1:]), nil
	}
	return "", fmt.Errorf("invalid data message header: %x", header)
}

func GenerateDataMessage(message string) []byte {
	buf := make([]byte, len(message)+1)
	buf[0] = DataHeader
	copy(buf[1:], message)
	return buf
}
