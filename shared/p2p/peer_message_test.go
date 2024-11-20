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
	msg, err := DiscoveryMessage("client-adapter", "192.168.1.1:3000", "192.168.1.2:4000")

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
	msg, err := DiscoveryMessage("wizard", "", "")

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

	if role != "client-adapter" {
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

	if role != "wizard" {
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
