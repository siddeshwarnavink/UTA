package proxy

// import (
// 	"fmt"
// 	"log"
// 	"net"
// 	"os"
// 	"sync"

// 	"github.com/siddeshwarnavink/UTA/adapter/embeded"
// 	"github.com/siddeshwarnavink/UTA/adapter/ui"

// 	lua "github.com/yuin/gopher-lua"
// )

// func BroadcastServerProxy(l *lua.LState, peerConn net.UDPConn) {
// 	fromAddress := embeded.CurrentFlags.Dec
// 	toAddress := embeded.CurrentFlags.Enc

// 	listener, err := net.Listen("tcp", fromAddress)
// 	if err != nil {
// 		log.Fatalf("Could not start client adapter: %v", err)
// 	}
// 	defer listener.Close()

// 	fmt.Printf("Broadcast Server adapter listening on %s, forwarding to client %s\n",
// 		fromAddress, toAddress)

// 	// Keep track of all connected clients
// 	var clientConns []net.Conn
// 	var mu sync.Mutex

// 	for {
// 		plainConn, err := listener.Accept()
// 		if err != nil {
// 			log.Printf("Could not accept connection: %v", err)
// 			continue
// 		}

// 		// Add client to the list of connections
// 		mu.Lock()
// 		clientConns = append(clientConns, plainConn)
// 		mu.Unlock()

// 		// Connect to the server
// 		encryptedConn, err := net.Dial("tcp", toAddress)
// 		if err != nil {
// 			log.Printf("Could not connect to client: %v", err)
// 			return
// 		}
// 		defer encryptedConn.Close()

// 		keyalgo, err := ui.KeyAlgorithmFromString(embeded.CurrentFlags.KeyAlgo)
// 		if err != nil {
// 			fmt.Println(err)
// 			os.Exit(1)
// 		}
// 		keyalgo.Key = keyalgo.Generate(encryptedConn)

// 		if !IsUninitialized(keyalgo.Key) {
// 			fmt.Printf("\nGot shared key %x\n", keyalgo.Key)

// 			algo, err := ui.AlgorithmFromString(embeded.CurrentFlags.CryptoAlgo)
// 			if err != nil {
// 				fmt.Println(err)
// 				os.Exit(1)
// 			}

// 			go ProxyHandlerForMultipleClients(plainConn, encryptedConn, keyalgo.Key, algo, peerConn, &clientConns, &mu)
// 		}
// 	}
// }
