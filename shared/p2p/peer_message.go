/*

Peer Message format,
by Sideshwar (who only thinks like web developer btw)

+------+-----------------+------------+
| S.no | Content         | Size(bits) |
+------+-----------------+------------+
| 1    | Message type    | 4          |
| 2    | Peer type       | 4          |
| 3    | Akshual message | -          |
+------+-----------------+------------+

+--------------+------+
| Message type | Bits |
+--------------+------+
| Discovery    | 0000 |
| Transmission | 0001 |
| String       | 0010 |
+--------------+------+

+----------------+------+
| Peer type      | Bits |
+----------------+------+
| Client adapter | 0000 |
| Server adapter | 0001 |
| Wizard         | 0010 |
+----------------+------+

Discovery message contains 2 IPv4 - fromIP, toIP

Transmission message contains only one bit,
0 - sent, 1 - received

*/

package p2p

import (
	"encoding/json"
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

func GetPeerMsgType(bits []byte) (PeerMsgType, error) {
	// Extract the first 4 bits
	typeBits := bits[0] >> 4

	switch typeBits {
	case 0b0000:
		return Discovery, nil
	case 0b0001:
		return Transmission, nil
	case 0b0010:
		strmsg := string(bits[1:])
		lastchar := strmsg[len(strmsg)-1:] 
		if lastchar == "0" {
			return StringRequestMessageType, nil
		} else if lastchar == "1"  {
			return StringResponseMessageType, nil
		}
		return StringMessageType, nil
	default:
		return Invalid, fmt.Errorf("Invalid peer message type: %08b", typeBits)
	}
}

func bytesToIpv4(data []byte) (string, error) {
	// Check if the data is exactly 6 bytes
	if len(data) != 6 {
		return "", fmt.Errorf("invalid byte slice length, expected 6 bytes")
	}

	// Extract the 4-byte IPv4 address
	ipv4 := data[:4]
	// Extract the 2-byte port number (big-endian format)
	port := int(data[4])<<8 | int(data[5])

	// Convert the 4-byte slice back to an IP address string
	ip := net.IPv4(ipv4[0], ipv4[1], ipv4[2], ipv4[3]).String()

	// Return the address and port in the format "address:port"
	return fmt.Sprintf("%s:%d", ip, port), nil
}

func splitAddressAndPort(input string) (string, int, error) {
	// Split the string at the colon (':')
	parts := strings.Split(input, ":")
	if len(parts) != 2 {
		return "", 0, fmt.Errorf("invalid address:port format")
	}

	// The first part is the address (string), the second part is the port (string)
	address := parts[0]
	portStr := parts[1]

	// Convert the port to an integer
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return "", 0, fmt.Errorf("invalid port: %v", err)
	}

	return address, port, nil
}

func ipv4ToBytes(addr string) ([]byte, error) {
	ip, port, err := splitAddressAndPort(addr)
	if err != nil {
		return nil, err
	}

	// Parse the IPv4 address
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return nil, fmt.Errorf("invalid IPv4 address")
	}

	// Extract the 4-byte representation of the IPv4 address
	ipv4 := parsedIP.To4()
	if ipv4 == nil {
		return nil, fmt.Errorf("not an IPv4 address")
	}

	// Check if the port is in the valid range
	if port < 0 || port > 65535 {
		return nil, fmt.Errorf("invalid port number")
	}

	// Create a 6-byte array: 4 bytes for IP and 2 bytes for port
	result := make([]byte, 6)

	// Copy the 4-byte IPv4 address into the result
	copy(result[:4], ipv4)

	// Add the 2-byte port number (big-endian encoding)
	result[4] = byte(port >> 8)   // High byte
	result[5] = byte(port & 0xFF) // Low byte

	return result, nil
}

func DiscoveryMessage(role PeerRole, fromIP string, toIP string) ([]byte, error) {
	msg := []byte{0b00000000}

	roleBits, err := getRoleBits(role)
	if err != nil {
		return nil, err
	}

	msg[0] |= roleBits

	if roleBits != 0b10 { // If not "Wizard"
		fromIPBytes, err := ipv4ToBytes(fromIP)
		if err != nil {
			return nil, err
		}

		toIPBytes, err := ipv4ToBytes(toIP)
		if err != nil {
			return nil, err
		}

		msg = append(msg, fromIPBytes...)
		msg = append(msg, toIPBytes...)
	}

	return msg, nil
}

func ExtractDiscoveryMessageDetails(msg []byte) (PeerRole, string, string, error) {
	msgtype, err := GetPeerMsgType(msg)
	if err != nil || msgtype != Discovery {
		return "", "", "", fmt.Errorf("not discovery type message")
	}

	role, err := getRoleFromBits(msg)
	if err != nil {
		return "", "", "", err
	}

	if role != Wizard {
		fromIPBytes := msg[1:7]
		fromIP, err := bytesToIpv4(fromIPBytes)
		if err != nil {
			return "", "", "", fmt.Errorf("failed to convert from IP bits: %v", err)
		}

		toIPBits := msg[7:13]
		toIP, err := bytesToIpv4(toIPBits)
		if err != nil {
			return "", "", "", fmt.Errorf("failed to convert to IP bits: %v", err)
		}

		return role, fromIP, toIP, nil
	} else {
		return role, "", "", nil
	}
}

// Map roles to their corresponding binary values
func getRoleBits(role PeerRole) (uint8, error) {
	roleMap := map[PeerRole]uint8{
		ClientProxy: 0b00,
		ServerProxy: 0b01,
		Wizard:      0b10,
	}

	roleBits, ok := roleMap[role]
	if !ok {
		return 0, fmt.Errorf("Invalid role: %s", role)
	}
	return roleBits, nil
}

func getRoleFromBits(msg []byte) (PeerRole, error) {
	last4Bits := msg[0] & 0x0F

	var role PeerRole

	switch last4Bits {
	case 0x00:
		role = ClientProxy
	case 0x01:
		role = ServerProxy
	case 0x02:
		role = Wizard
	default:
		return InvalidRole, fmt.Errorf("unknown role for bit: %08b", last4Bits)
	}

	return role, nil
}

func TransmissionMessage(role PeerRole, sent bool) ([]byte, error) {
	msg := []byte{0b00010000}

	roleBits, err := getRoleBits(role)
	if err != nil {
		return nil, err
	}

	msg[0] |= roleBits

	if sent {
		msg = append(msg, 0x00)
	} else {
		msg = append(msg, 0x01)
	}

	return msg, nil
}

func ExtractTransmissionMessageDetails(msg []byte) (PeerRole, bool, error) {
	msgtype, err := GetPeerMsgType(msg)
	if err != nil || msgtype != Discovery {
		return InvalidRole, false, fmt.Errorf("not discovery type message")
	}

	role, err := getRoleFromBits(msg)
	if err != nil {
		return InvalidRole, false, err
	}

	last4Bits := msg[0] & 0x0F

	return role, last4Bits == 0x00, err
}

func StringMessage(role PeerRole, message string) ([]byte, error) {
	msg := []byte{0b00100000}

	roleBits, err := getRoleBits(role)
	if err != nil {
		return nil, err
	}

	msg[0] |= roleBits

	msg = append(msg, []byte(message)...)

	return msg, nil
}

func ExtractStringMessage(msg []byte) (PeerRole, string, error) {
	msgtype, err := GetPeerMsgType(msg)
	if err != nil || !(msgtype == StringMessageType || msgtype == StringRequestMessageType || msgtype == StringResponseMessageType) {
		return InvalidRole, "", fmt.Errorf("not string message")
	}

	role, err := getRoleFromBits(msg)
	if err != nil {
		return InvalidRole, "", err
	}

	strmsg := string(msg[1:])
	return role, strmsg, nil
}

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

func RequestMessage(role PeerRole, reqType RequestMessageType, requestId string, payload string) ([]byte, error) {
	reqObj := RequestMessageJson{
		Type:      reqType,
		RequestId: requestId,
		Payload:   payload,
	}

	json, err := json.Marshal(reqObj)
	if err != nil {
		return nil, err
	}

	return StringMessage(role, string(json)+"0")
}

func ExtractRequestMessage(message []byte) (PeerRole, RequestMessageType, string, string, error) {
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

func ResponseMessage(role PeerRole, requestId string, data string) ([]byte, error) {
	reqObj := ResponseMessageJson{
		RequestId: requestId,
		Data:      data,
	}

	json, err := json.Marshal(reqObj)
	if err != nil {
		return nil, err
	}

	return StringMessage(role, string(json)+"1")
}

func ExtractResponseMessage(message []byte) (PeerRole, string, string, error) {
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
