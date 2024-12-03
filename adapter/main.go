package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/siddeshwarnavink/UTA/adapter/embeded"
	"github.com/siddeshwarnavink/UTA/adapter/proxy"
	"github.com/siddeshwarnavink/UTA/adapter/ui"
	"github.com/siddeshwarnavink/UTA/shared/p2p"

	lua "github.com/yuin/gopher-lua"
)

func main() {
	// logging
	logFile, err := os.OpenFile("logs/adapter.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Printf("Error opening log file: %v\n", err)
		os.Exit(1)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)

	// lua stack
	l := lua.NewState()
	defer l.Close()

	embeded.HandleLua(l)

	ui.RenderForm()

	peerTable := p2p.NewPeerTable()

	peerConn, err := p2p.GetMulticastConn()
	if err != nil {
		panic(err)
	}
	defer peerConn.Close()

	if embeded.currentFlags.Mode == embeded.Client {
		go p2p.AnnouncePresence(*peerConn, p2p.ClientProxy, embeded.currentFlags.Dec, embeded.currentFlags.Enc)
	} else {
		go p2p.AnnouncePresence(*peerConn, p2p.ServerProxy, embeded.currentFlags.Dec, embeded.currentFlags.Enc)
	}
	p2p.ListenForPeers(peerTable)

	switch embeded.currentFlags.Mode {
	case embeded.Client:
		ClientProxy(l, *peerConn)
	case embeded.Server:
		ServerProxy(l, *peerConn)
	}
}

func ClientProxy(l *lua.LState, peerConn net.UDPConn) {
	fromAddress := embeded.currentFlags.Dec
	toAddress := embeded.currentFlags.Enc

	listener, err := net.Listen("tcp", fromAddress)
	if err != nil {
		log.Fatalf("Could not start client adapter: %v", err)
	}
	defer listener.Close()

	fmt.Printf("Client adapter listening on %s, forwarding to server %s\n",
		fromAddress, toAddress)

	for {
		plainConn, err := listener.Accept()
		if err != nil {
			log.Printf("Could not accept connection: %v", err)
			continue
		}

		encryptedConn, err := net.Dial("tcp", toAddress)
		if err != nil {
			log.Printf("Could not connect to server: %v", err)
			return
		}
		defer encryptedConn.Close()

		keyalgo, err := ui.KeyAlgorithmFromString(embeded.currentFlags.Protocol)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		keyalgo.Key = keyalgo.Generate(encryptedConn)

		if !proxy.IsUninitialized(keyalgo.Key) {
			fmt.Printf("\nGot shared key %x\n", keyalgo.Key)

			algo, err := ui.AlgorithmFromString(embeded.currentFlags.Algo)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			go proxy.ProxyHandler(plainConn, encryptedConn, keyalgo.Key, algo, peerConn)
		}
	}
}

func ServerProxy(l *lua.LState, peerConn net.UDPConn) {
	fromAddress := embeded.currentFlags.Enc
	toAddress := embeded.currentFlags.Dec

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

		keyalgo, err := ui.KeyAlgorithmFromString(embeded.currentFlags.Protocol)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		keyalgo.Key = keyalgo.Generate(encryptedConn)

		if !proxy.IsUninitialized(keyalgo.Key) {
			fmt.Printf("\nGot shared key %x\n", keyalgo.Key)

			algo, err := ui.AlgorithmFromString(embeded.currentFlags.Algo)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			go proxy.ProxyHandler(plainConn, encryptedConn, keyalgo.Key, algo, peerConn)
		}
	}
}
