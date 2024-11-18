package utils

import (
	"errors"
	"strconv"
	"strings"
)

func ExtractPort(ipv4 string) (int, error) {
	parts := strings.Split(ipv4, ":")
	if len(parts) != 2 {
		return 0, errors.New("invalid IPv4")
	}

	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, errors.New("invalid port")
	}

	return port, nil
}
