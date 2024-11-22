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
	"github.com/siddeshwarnavink/UTA/shared/p2p"
	"github.com/siddeshwarnavink/UTA/shared/utils"
)

func ProxyHandler(plainConn net.Conn,
	encryptedConn net.Conn,
	derivedKey []byte,
	algo *embeded.CryptoAlgo) {

	var wg sync.WaitGroup
	wg.Add(2)

	peerConn, err := p2p.GetMulticastConn()
	if err != nil {
		panic(err)
	}
	defer peerConn.Close()

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
			formatedData := utils.GenerateDataMessage(string(encryptedData))
			_, err = encryptedConn.Write(formatedData)

			if err != nil {
				if errors.Is(err, os.ErrClosed) {
					log.Printf("Encrypted connection closed, closing plain connection")
				} else {
					fmt.Printf("Error reading from client: %v", err)
				}

				plainConn.Close()
				return
			}

			peerMsg, err :=	p2p.TransmissionMessage(p2p.ClientProxy, true)
			peerMsgBytes := []byte(peerMsg)
			peerConn.Write(peerMsgBytes)
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

			unformatedData, err := utils.GetDataMessage(buf[:n])
			decryptedData := algo.Decrypt(derivedKey, []byte(unformatedData))

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

			peerMsg, err :=	p2p.TransmissionMessage(p2p.ClientProxy, false)
			peerMsgBytes := []byte(peerMsg)
			peerConn.Write(peerMsgBytes)
		}
	}()

	wg.Wait()
}
