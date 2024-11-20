/*

Peer Message format,
by Sideshwar (who only thinks like web developer btw)

+------+-----------------+------------+
| S.no | Content         | Size(bits) |
+------+-----------------+------------+
| 1    | Message type    | 2          |
| 2    | Peer type       | 2          |
| 3    | Message length  | 8          |
| 4    | Akshual message | -          |
+------------------------+------------+
| MIN SIZE			     | 12         |
+------------------------+------------+

+--------------+------+
| Message type | Bits |
+--------------+------+
| Discovery    | 00   |
| Transmission | 01   |
+--------------+------+

+----------------+------+
| Peer type      | Bits |
+----------------+------+
| Client adapter | 00   |
| Server adapter | 01   |
| Wizard         | 10   |
+----------------+------+

*/

package p2p

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
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

func DiscoveryMessage(role PeerRole, fromIP string, toIP string) (string, error) {
	msg := "00"

	roleMap := map[PeerRole]string{
		ClientProxy: "00",
		ServerProxy: "01",
		Wizard:      "10",
	}

	roleBits, ok := roleMap[role]
	if !ok {
		roleBits = roleMap["wizard"]
	}
	msg += roleBits

	if roleBits != "10" {
		msg += "01000000" // 64-bits

		fromIPBits, err := convertIPv4ToBits(fromIP)
		if err != nil {
			return "", fmt.Errorf("invalid fromIP: %v", err)
		}
		msg += fromIPBits

		toIPBits, err := convertIPv4ToBits(toIP)
		if err != nil {
			return "", fmt.Errorf("invalid toIP: %v", err)
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
		return "", "", "", fmt.Errorf("unknown role bits: %s", roleBits)
	}

	if role != "wizard" {
		fromIPBits := msg[12:60]
		fromIP, err := convertBitsToIPv4(fromIPBits)
		if err != nil {
			return "", "", "", fmt.Errorf("failed to convert fromIP bits: %v", err)
		}

		toIPBits := msg[60:108]
		toIP, err := convertBitsToIPv4(toIPBits)
		if err != nil {
			return "", "", "", fmt.Errorf("failed to convert toIP bits: %v", err)
		}

		return role, fromIP, toIP, nil
	} else {
		return role, "", "", nil
	}
}
