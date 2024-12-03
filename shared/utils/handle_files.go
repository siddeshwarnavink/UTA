package utils

import (
	"fmt"
	"io"
	"os"
)

func PaginateFile(filePath string, pageNumber, chunkSize int) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if pageNumber <= 0 {
		return nil, fmt.Errorf("Invalid page number: %d", pageNumber)
	}

	offset := (pageNumber - 1) * chunkSize

	_, err = file.Seek(int64(offset), io.SeekStart)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, chunkSize)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return buf[:n], nil
}
