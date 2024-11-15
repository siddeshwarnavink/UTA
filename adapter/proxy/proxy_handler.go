package proxy

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"

	"github.com/siddeshwarnavink/UTA/adapter/embeded"
)

func ProxyHandler(plainConn net.Conn,
	encryptedConn net.Conn,
	derivedKey []byte,
	algo *embeded.CryptoAlgo) {

	var wg sync.WaitGroup
	wg.Add(2)

	// plain -> encrypted
	go func() {
		defer wg.Done()

		buf := make([]byte, 1024)
		for {
			n, err := plainConn.Read(buf)

			if err != nil {
				if errors.Is(err, io.EOF) {
					log.Printf("Plain connection closed, closing encrypted Conn")
				} else {
					log.Printf("Error reading from encrypted Conn: %v", err)
				}

				encryptedConn.Close()
				return
			}

			encryptedData := algo.Encrypt(derivedKey, buf[:n])
			_, err = encryptedConn.Write(encryptedData)

			if err != nil {
				if errors.Is(err, os.ErrClosed) {
					log.Printf("Encrypted connection closed, closing plain connection")
				} else {
					fmt.Printf("Error reading from client: %v", err)
				}

				plainConn.Close()
				return
			}
		}
	}()

	// encrypted -> plain
	go func() {
		defer wg.Done()

		buf := make([]byte, 1024)
		for {
			n, err := encryptedConn.Read(buf)

			if err != nil {
				if errors.Is(err, io.EOF) {
					log.Printf("Encrypted connection closed, closing plain Conn")
				} else {
					log.Printf("Error reading from plain Conn: %v", err)
				}

				plainConn.Close()
				return
			}

			decryptedData := algo.Decrypt(derivedKey, buf[:n])

			_, err = plainConn.Write(decryptedData)

			if err != nil {
				if errors.Is(err, os.ErrClosed) {
					log.Printf("Plain connection closed, closing encrypted connection")
				} else {
					fmt.Printf("Error reading from server: %v", err)
				}

				encryptedConn.Close()
				return
			}
		}
	}()

	wg.Wait()
}
