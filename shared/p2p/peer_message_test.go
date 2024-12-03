package p2p

import (
	"testing"
)

func TestIPToBits(t *testing.T) {
	bits, err := convertIPv4ToBits("192.168.1.1:3000")

	if err != nil {
		t.Errorf("convertIPv4ToBits() error = %v", err)
		return
	}

	if bits != "110000001010100000000001000000010000101110111000" {
		t.Errorf("convertIPv4ToBits() incorrect bits %s", bits)
		return
	}
}

func TestBitsToIP(t *testing.T) {
	ip, err := convertBitsToIPv4("110000001010100000000001000000010000101110111000")

	if err != nil {
		t.Errorf("convertBitsToIPv4() error = %v", err)
		return
	}

	if ip != "192.168.1.1:3000" {
		t.Errorf("convertBitsToIPv4() incorrect ip %s", ip)
		return
	}
}

func TestDiscoveryMessageForClientAdapter(t *testing.T) {
	msg, err := DiscoveryMessage(ClientProxy, "192.168.1.1:3000", "192.168.1.2:4000")

	if err != nil {
		t.Errorf("DiscoveryMessage() error = %v", err)
		return
	}

	// 00-00-01000000-110000001010100000000001000000010000101110111000-110000001010100000000001000000100000111110100000
	if msg != "000001000000110000001010100000000001000000010000101110111000110000001010100000000001000000100000111110100000" {
		t.Errorf("DiscoveryMessage() incorrect message %s", msg)
		return
	}
}

func TestDiscoveryMessageForWizard(t *testing.T) {
	msg, err := DiscoveryMessage(Wizard, "", "")

	if err != nil {
		t.Errorf("DiscoveryMessage() error = %v", err)
		return
	}

	// 00-10-00000000
	if msg != "001000000000" {
		t.Errorf("DiscoveryMessage() incorrect message %s", msg)
		return
	}
}

func TestExtractDiscoveryMessageDetailsForClientProxy(t *testing.T) {
	role, fromIP, toIP, err := ExtractDiscoveryMessageDetails("000001000000110000001010100000000001000000010000101110111000110000001010100000000001000000100000111110100000")

	if err != nil {
		t.Errorf("ExtractDiscoveryMessageDetails() error = %v", err)
		return
	}

	if role != ClientProxy {
		t.Errorf("ExtractDiscoveryMessageDetails() incorrect role = %s", role)
		return
	}

	if fromIP != "192.168.1.1:3000" {
		t.Errorf("ExtractDiscoveryMessageDetails() incorrect fromIP = %s", role)
		return
	}

	if toIP != "192.168.1.2:4000" {
		t.Errorf("ExtractDiscoveryMessageDetails() incorrect toIP = %s", role)
		return
	}
}

func TestExtractDiscoveryMessageDetailsForWizard(t *testing.T) {
	role, fromIP, toIP, err := ExtractDiscoveryMessageDetails("001000000000")

	if err != nil {
		t.Errorf("ExtractDiscoveryMessageDetails() error = %v", err)
		return
	}

	if role != Wizard {
		t.Errorf("ExtractDiscoveryMessageDetails() incorrect role = %s", role)
		return
	}

	if fromIP != "" {
		t.Errorf("ExtractDiscoveryMessageDetails() incorrect fromIP = %s", role)
		return
	}

	if toIP != "" {
		t.Errorf("ExtractDiscoveryMessageDetails() incorrect toIP = %s", role)
		return
	}
}

func TestTransmissionMessage(t *testing.T) {
	msg, err := TransmissionMessage(ClientProxy, true)

	if err != nil {
		t.Errorf("TransmissionMessage() error = %v", err)
		return
	}

	// 01-00-000000001-0
	if msg != "01000000000010" {
		t.Errorf("TransmissionMessage() incorrect message = %s", msg)
		return
	}
}

func TestExtractTransmissionMessageDetails(t *testing.T) {
	role, sent, err := ExtractTransmissionMessageDetails("01000000000010")

	if err != nil {
		t.Errorf("ExtractTransmissionMessageDetails() error = %v", err)
		return
	}

	if role != ClientProxy {
		t.Errorf("ExtractTransmissionMessageDetails() incorrect role = %s", role)
		return
	}

	if !sent {
		t.Errorf("ExtractTransmissionMessageDetails() incorrect transmission state")
		return
	}
}

func TestStringMessage(t *testing.T) {
	msg, err := StringMessage(ClientProxy, "yo mama")
	if err != nil {
		t.Errorf("TestStringMessage() error = %v", err)
		return
	}
	if msg != "100000000000yo mama" {
		t.Errorf("TestStringMessage() error = %v", msg)
		return
	}

}

func TestExtractStringMessage(t *testing.T) {
	role, msg, err := ExtractStringMessage("100000000000yo mama")
	if err != nil {
		t.Errorf("TestExtractStringMessage() error = %v", err)
		return
	}
	if role != ClientProxy {
		t.Errorf("TestStringMessage() invalid role = %v", role)
		return
	}
	if msg != "yo mama" {
		t.Errorf("TestStringMessage() invalid message = %v", msg)
		return
	}

}

func TestRequestMessage(t *testing.T) {
	msg, err := RequestMessage(Wizard, RequestTypeConfig, "abc", "def")
	if err != nil {
		t.Errorf("TestRequestMessage() error = %v", err)
		return
	}
	if msg != "101000000000{\"t\":0,\"i\":\"abc\",\"p\":\"def\"}0" {
		t.Errorf("TestRequestMessage() invalid message = %s", msg)
		return
	}
}

func TestStringMessageType(t *testing.T) {
	bits := "101000000000{\"t\":1,\"i\":\"d4692450-9775-4d5d-a346-4fbf916904c2\",\"p\":\"192.168.1.7:41474\"}0"
	msgtype, err := GetPeerMsgType(bits)
	if err != nil {
		t.Errorf("TestStringMessageType() error = %v", err)
		return
	}

	if msgtype != StringRequestMessageType {
		t.Errorf("TestStringMessageType() invalid type = %d", msgtype)
	}
}

func TestExtractMessageRequest(t *testing.T) {
	bits := "101000000000{\"t\":0,\"i\":\"c6ac0727-9469-403e-ac5b-cc2a73723199\",\"p\":\"192.168.1.7:52026\"}0"
	role, reqType, reqId, payload, err := ExtractRequestMessage(bits)

	if err != nil {
		t.Errorf("TestExtractMessageRequest() error = %v", err)
		return
	}

	if role != Wizard {
		t.Errorf("TestExtractMessageRequest() invalid role = %s", role)
		return
	}

	if reqType != RequestTypeConfig {
		t.Errorf("TestExtractMessageRequest() invalid request type = %d", reqType)
		return
	}

	if reqId != "c6ac0727-9469-403e-ac5b-cc2a73723199" {
		t.Errorf("TestExtractMessageRequest() invalid request id = %s", reqId)
		return
	}

	if payload != "192.168.1.7:52026" {
		t.Errorf("TestExtractMessageRequest() invalid payload = %s", payload)
		return
	}
}

func TestExtractResponse(t *testing.T) {
	bits := "100100000000{\"i\":\"68d8535e-ca09-44c2-9781-87b6247162a0\",\"d\":\"deeznuts\"}1"

	role, reqId, data, err := ExtractResponseMessage(bits)

	if err != nil {
		t.Errorf("TestExtractResponse() error = %v", err)
		return
	}

	if role != ServerProxy {
		t.Errorf("TestExtractResponse() invalid role = %s", role)
		return
	}

	if data != "deeznuts" {
		t.Errorf("TestExtractResponse() invalid data = %s", data)
		return
	}

	if reqId != "68d8535e-ca09-44c2-9781-87b6247162a0" {
		t.Errorf("TestExtractResponse() invalid request id = %s", reqId)
		return
	}
}
