/*

Peer Message format,
by Sideshwar (who only thinks like web developer btw)

+----------+-----------------+------------+
| S.no     | Content         | Size(bits) |
+----------+-----------------+------------+
| 1        | Message type    | 4          |
| 2        | Peer type       | 4          |
| 4        | Akshual message | -          |
+----------+-----------------+------------+

+--------------+-----+
| Message type | Bits|
+--------------+-----+
| Discovery    | 0x1 |
| Transmission | 0x2 |
| String       | 0x3 |
+--------------+-----+

+----------------+-----+
| Peer type      | Bits|
+----------------+-----+
| Client adapter | 0x1 |
| Server adapter | 0x2 |
| Wizard         | 0x3 |
+----------------+-----+

Discovery message contains 2 IPv4 - fromIP, toIP

Transmission message contains only one byte,
00000001 - sent, 00000010 - received

*/

package p2p

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type PeerMsgType string

const (
	Invalid      string = "invalid"
	Discovery    string = "discovery"
	Transmission string = "transmission"
)

const (
	disc  string = "0001"
	trans string = "0010"
	str   string = "0011"

	StringRequestMessageType  string = "00000001"
	StringResponseMessageType string = "00000010"
)

const (
	cli string = "0001"
	srv string = "0010"
	wiz string = "0011"
)

const (
	sen string = "00000001"
	rec string = "00000010"
)

const (
	MessageTypeStart int = 0
	MessageTypeLen   int = 4
	PeerTypeStart    int = 4
	PeerTypeLen      int = 4
	PayloadStart     int = 8
)

func convertIPv4ToBits(address string) (string, error) {
	host, portStr, err := net.SplitHostPort(address)
	if err != nil {
		host = address
	}

	var ip = net.ParseIP(host)

	var bitString strings.Builder
	if ipv4 := ip.To4(); ipv4 != nil {
		for _, b := range ipv4 {
			bitString.WriteString(fmt.Sprintf("%08b", b))
		}
	} else {
		return "", fmt.Errorf("not IPv4 '%s'", host)
	}

	port, err := strconv.Atoi(portStr)
	if err != nil || port < 0 || port > 65535 {
		return "", fmt.Errorf("invalid port: %v", err)
	}

	portBits := fmt.Sprintf("%016b", port)
	bitString.WriteString(portBits)

	return bitString.String(), nil
}

func convertBitsToIPv4(bits string) (string, error) {
	if len(bits) != 48 {
		return "", errors.New("invalid bit string length; must be 32 bits")
	}

	var ipParts []string
	for i := 0; i < 4; i++ {
		bitSegment := bits[i*8 : (i+1)*8]
		part, err := strconv.ParseInt(bitSegment, 2, 32)
		if err != nil {
			return "", fmt.Errorf("failed to parse bits to integer: %v", err)
		}
		ipParts = append(ipParts, fmt.Sprintf("%d", part))
	}

	ipAddress := strings.Join(ipParts, ".")

	portBits := bits[32:48]
	port, err := strconv.ParseInt(portBits, 2, 16)
	if err != nil {
		return "", fmt.Errorf("failed to parse port bits to integer: %v", err)
	}

	return fmt.Sprintf("%s:%d", ipAddress, port), nil
}

func MsgType(bits string) (string, error) {
	typeBits := bits[MessageTypeStart : MessageTypeStart+MessageTypeLen]
	switch typeBits {
	case disc:
		return Discovery, nil
	case trans:
		return Transmission, nil
	case str:
		switch bits[len(bits)-8:] {
		case StringRequestMessageType:
			return StringRequestMessageType, nil
		case StringResponseMessageType:
			return StringResponseMessageType, nil
		default:
			return Invalid, fmt.Errorf("Invalid peer message type: %s", typeBits)
		}
	default:
		return Invalid, fmt.Errorf("Invalid peer message type: %s", typeBits)
	}
}

func getPeerBits(role PeerRole) (string, error) {
	roleMap := map[PeerRole]string{
		ClientProxy: cli,
		ServerProxy: srv,
		Wizard:      wiz,
	}

	roleBits, ok := roleMap[role]
	if !ok {
		return "", fmt.Errorf("Invalid role: %s", role)
	}
	return roleBits, nil
}

func getPeerFromBits(msg string) (PeerRole, error) {
	roleBits := msg[PeerTypeStart : PeerTypeStart+PeerTypeLen]

	var role PeerRole
	switch roleBits {
	case cli:
		role = ClientProxy
	case srv:
		role = ServerProxy
	case wiz:
		role = Wizard
	default:
		return InvalidRole, fmt.Errorf("unknown role for bit: %s", roleBits)
	}
	return role, nil
}

func DiscoveryMessage(role PeerRole, fromIP string, toIP string) (string, error) {
	msg := disc

	roleBits, err := getPeerBits(role)
	if err != nil {
		return "", err
	}

	msg += roleBits

	if roleBits != wiz {
		fromIPBits, err := convertIPv4ToBits(fromIP)
		if err != nil {
			return "", fmt.Errorf("invalid from IP: %v", err)
		}
		msg += fromIPBits
		toIPBits, err := convertIPv4ToBits(toIP)
		if err != nil {
			return "", fmt.Errorf("invalid to IP: %v", err)
		}
		msg += toIPBits
	}

	return msg, nil
}

func ExtractDiscoveryMessage(msg string) (PeerRole, string, string, error) {
	role, err := getPeerFromBits(msg)
	if err != nil {
		return "", "", "", err
	}

	if role != Wizard {
		fromIPBits := msg[PayloadStart : PayloadStart+32]
		fromIP, err := convertBitsToIPv4(fromIPBits)
		if err != nil {
			return "", "", "", fmt.Errorf("failed to convert from IP bits: %v", err)
		}

		toIPBits := msg[PayloadStart+32 : PeerTypeStart+64]
		toIP, err := convertBitsToIPv4(toIPBits)
		if err != nil {
			return "", "", "", fmt.Errorf("failed to convert to IP bits: %v", err)
		}

		return role, fromIP, toIP, nil
	} else {
		return role, "", "", nil
	}
}

func TransmissionMessage(role PeerRole, sent bool) (string, error) {
	msg := trans
	roleBits, err := getPeerBits(role)
	if err != nil {
		return "", err
	}

	msg += roleBits
	if sent {
		msg += sen
	} else {
		msg += rec
	}
	return msg, nil
}

func ExtractTransmissionMessage(msg string) (PeerRole, bool, error) {
	msgtype, err := MsgType(msg)
	if err != nil || msgtype != Transmission {
		return InvalidRole, false, fmt.Errorf("not transmission type message")
	}

	role, err := getPeerFromBits(msg)
	if err != nil {
		return InvalidRole, false, err
	}

	transmissionBit := string(msg[len(msg)-8])

	return role, transmissionBit == sen, err
}

func StringMessage(role PeerRole, message string) (string, error) {
	msg := str

	roleBits, err := getPeerBits(role)
	if err != nil {
		return "", err
	}

	msg += roleBits
	msg += message

	return msg, nil
}

func ExtractStringMessage(msg string) (PeerRole, string, error) {
	msgtype, err := MsgType(msg)
	if err != nil || !(msgtype == StringRequestMessageType || msgtype == StringResponseMessageType) {
		return InvalidRole, "", fmt.Errorf("not string message")
	}

	role, err := getPeerFromBits(msg)
	if err != nil {
		return InvalidRole, "", err
	}

	strmsg := msg[PayloadStart:]
	return role, strmsg, nil
}
