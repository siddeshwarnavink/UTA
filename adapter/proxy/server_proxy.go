package proxy

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/siddeshwarnavink/UTA/adapter/embeded"
	"github.com/siddeshwarnavink/UTA/adapter/ui"

	lua "github.com/yuin/gopher-lua"
)

func ServerProxy(l *lua.LState, peerConn net.UDPConn) {
	fromAddress := embeded.CurrentFlags.Enc
	toAddress := embeded.CurrentFlags.Dec

	listener, err := net.Listen("tcp", fromAddress)
	if err != nil {
		log.Fatalf("Could not start server proxy: %v", err)
	}
	defer listener.Close()

	fmt.Printf("Server adapter listening on %s, forwarding to server %s\n",
		fromAddress, toAddress)

	for {
		encryptedConn, err := listener.Accept()
		if err != nil {
			log.Printf("Could not accept connection: %v", err)
			continue
		}

		plainConn, err := net.Dial("tcp", toAddress)
		if err != nil {
			log.Printf("Could not connect to server: %v", err)
			return
		}
		defer plainConn.Close()

		keyalgo, err := ui.KeyAlgorithmFromString(embeded.CurrentFlags.KeyAlgo)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		keyalgo.Key = keyalgo.Generate(encryptedConn)

		if !IsUninitialized(keyalgo.Key) {
			fmt.Printf("\nGot shared key %x\n", keyalgo.Key)

			algo, err := ui.AlgorithmFromString(embeded.CurrentFlags.CryptoAlgo)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			go ProxyHandler(plainConn, encryptedConn, keyalgo.Key, algo, peerConn)
		}
	}
}
