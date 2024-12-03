/*

Peer Message format,
by Sideshwar (who only thinks like web developer btw)

+----------+-----------------+------------+
| S.no     | Content         | Size(bits) |
+----------+-----------------+------------+
| 1        | Message type    | 2          |
| 2        | Peer type       | 2          |
| 3        | Message length  | 8          |
| 4        | Akshual message | -          |
+----------+-----------------+------------+

+--------------+------+
| Message type | Bits |
+--------------+------+
| Discovery    | 00   |
| Transmission | 01   |
| String       | 10   |
+--------------+------+

+----------------+------+
| Peer type      | Bits |
+----------------+------+
| Client adapter | 00   |
| Server adapter | 01   |
| Wizard         | 10   |
+----------------+------+

Discovery message contains 2 IPv4 - fromIP, toIP

Transmission message contains only one bit,
0 - sent, 1 - received

*/

package p2p

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type PeerMsgType int

const (
	Invalid                   PeerMsgType = -1
	Discovery                 PeerMsgType = 0
	Transmission              PeerMsgType = 1
	StringMessageType         PeerMsgType = 2
	StringRequestMessageType  PeerMsgType = 3
	StringResponseMessageType PeerMsgType = 4
)

func convertIPv4ToBits(address string) (string, error) {
	host, portStr, err := net.SplitHostPort(address)
	if err != nil {
		host = address
	}

	ips, err := net.LookupIP(host)
	if err != nil {
		return "", fmt.Errorf("failed to resolve hostname '%s': %v", host, err)
	}
	if len(ips) == 0 {
		return "", fmt.Errorf("no IPs found for hostname '%s'", host)
	}
	parsedIP := ips[0]

	var bitString strings.Builder
	if ipv4 := parsedIP.To4(); ipv4 != nil {
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

func binaryStringToInt(binary string) int {
	var result int
	for i, bit := range binary {
		if bit == '1' {
			result += (1 << (7 - i))
		}
	}
	return result
}

func GetPeerMsgType(bits string) (PeerMsgType, error) {
	typeBits := bits[:2]
	switch typeBits {
	case "00":
		return Discovery, nil
	case "01":
		return Transmission, nil
	case "10":
		if bits[len(bits)-1:] == "0" {
			return StringRequestMessageType, nil
		} else if bits[len(bits)-1:] == "1" {
			return StringResponseMessageType, nil
		}
		return StringMessageType, nil
	default:
		return Invalid, fmt.Errorf("Invalid peer message type: %s", typeBits)
	}
}

func DiscoveryMessage(role PeerRole, fromIP string, toIP string) (string, error) {
	msg := "00"

	roleBits, err := getRoleBits(role)
	if err != nil {
		return "", err
	}

	msg += roleBits

	if roleBits != "10" {
		msg += "01000000" // 64-bits

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
	} else {
		msg += "00000000" // 0-bits
	}

	return msg, nil
}

func ExtractDiscoveryMessageDetails(msg string) (PeerRole, string, string, error) {
	if len(msg) != 108 && len(msg) != 12 {
		return "", "", "", fmt.Errorf("invalid message size")
	}

	msgtype, err := GetPeerMsgType(msg)
	if err != nil || msgtype != Discovery {
		return "", "", "", fmt.Errorf("not discovery type message")
	}

	role, err := getRoleFromBits(msg)
	if err != nil {
		return "", "", "", err
	}

	if role != Wizard {
		fromIPBits := msg[12:60]
		fromIP, err := convertBitsToIPv4(fromIPBits)
		if err != nil {
			return "", "", "", fmt.Errorf("failed to convert from IP bits: %v", err)
		}

		toIPBits := msg[60:108]
		toIP, err := convertBitsToIPv4(toIPBits)
		if err != nil {
			return "", "", "", fmt.Errorf("failed to convert to IP bits: %v", err)
		}

		return role, fromIP, toIP, nil
	} else {
		return role, "", "", nil
	}
}

func getRoleBits(role PeerRole) (string, error) {
	roleMap := map[PeerRole]string{
		ClientProxy: "00",
		ServerProxy: "01",
		Wizard:      "10",
	}

	roleBits, ok := roleMap[role]
	if !ok {
		return "", fmt.Errorf("Invalid role: %s", role)
	}
	return roleBits, nil
}

func getRoleFromBits(msg string) (PeerRole, error) {
	roleBits := msg[2:4]

	var role PeerRole
	switch roleBits {
	case "00":
		role = ClientProxy
	case "01":
		role = ServerProxy
	case "10":
		role = Wizard
	default:
		return InvalidRole, fmt.Errorf("unknown role for bit: %s", roleBits)
	}

	return role, nil
}

func TransmissionMessage(role PeerRole, sent bool) (string, error) {
	msg := "01"

	roleBits, err := getRoleBits(role)
	if err != nil {
		return "", err
	}

	msg += roleBits

	msg += "000000001" // 1-bit

	if sent {
		msg += "0"
	} else {
		msg += "1"
	}

	return msg, nil
}

func ExtractTransmissionMessageDetails(msg string) (PeerRole, bool, error) {
	if len(msg) != 14 {
		return InvalidRole, false, fmt.Errorf("invalid message size")
	}

	msgtype, err := GetPeerMsgType(msg)
	if err != nil || msgtype != Transmission {
		return InvalidRole, false, fmt.Errorf("not transmission type message")
	}

	role, err := getRoleFromBits(msg)
	if err != nil {
		return InvalidRole, false, err
	}

	transmissionBit := string(msg[len(msg)-1])

	return role, transmissionBit == "0", err
}

func StringMessage(role PeerRole, message string) (string, error) {
	msg := "10"

	roleBits, err := getRoleBits(role)
	if err != nil {
		return "", err
	}

	msg += roleBits
	msg += "00000000"
	msg += message

	return msg, nil
}

func ExtractStringMessage(msg string) (PeerRole, string, error) {
	msgtype, err := GetPeerMsgType(msg)
	if err != nil || !(msgtype == StringMessageType || msgtype == StringRequestMessageType || msgtype == StringResponseMessageType) {
		return InvalidRole, "", fmt.Errorf("not string message")
	}

	role, err := getRoleFromBits(msg)
	if err != nil {
		return InvalidRole, "", err
	}

	strmsg := msg[12:]
	return role, strmsg, nil
}

type RequestMessageType int

const (
	RequestMessageTypeInvalid RequestMessageType = -1
	RequestTypeConfig         RequestMessageType = 0
	RequestTypeLogs           RequestMessageType = 1
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
