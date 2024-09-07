package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/siddeshwarnavink/UTA/embeded"
	"github.com/siddeshwarnavink/UTA/keyExchange"
	"github.com/siddeshwarnavink/UTA/proxy"
	"github.com/siddeshwarnavink/UTA/ui"
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

	// lua stack
	l := lua.NewState()
	defer l.Close()

	embeded.HandleLua(l)

	flags, err := ui.ParseFlags()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	switch flags.Mode {
	case ui.Client:
		ClientProxy(l, flags)
	case ui.Server:
		ServerProxy(l, flags)
	}
}

func ClientProxy(l *lua.LState, flags *ui.Flags) {
	fromAddress := flags.Dec
	toAddress := flags.Enc

	listener, err := net.Listen("tcp", fromAddress)
	if err != nil {
		log.Fatalf("Could not start client adapter: %v", err)
	}
	defer listener.Close()

	fmt.Printf("Client adapter listening on %s, forwarding to server %s",
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

		derivedKey, err := keyExchange.ClientKeyExchange(encryptedConn, flags.Protocol)

		if !proxy.IsUninitialized(derivedKey) {
			fmt.Printf("\nGot shared key %x\n", derivedKey)

			algo, err := ui.AlgorithmFromString(flags.Algo)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			go proxy.ProxyHandler(plainConn, encryptedConn, derivedKey, algo)
		}
	}
}

func ServerProxy(l *lua.LState, flags *ui.Flags) {
	fromAddress := flags.Enc
	toAddress := flags.Dec

	listener, err := net.Listen("tcp", fromAddress)
	if err != nil {
		log.Fatalf("Could not start server proxy: %v", err)
	}
	defer listener.Close()

	fmt.Printf("Server adapter listening on %s, forwarding to server %s",
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

		derivedKey, err := keyExchange.ServerKeyExchange(encryptedConn, flags.Protocol)

		if !proxy.IsUninitialized(derivedKey) {
			fmt.Printf("\nGot shared key %x\n", derivedKey)

			algo, err := ui.AlgorithmFromString(flags.Algo)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			go proxy.ProxyHandler(plainConn, encryptedConn, derivedKey, algo)
		}
	}
}
