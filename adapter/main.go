package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/siddeshwarnavink/UTA/adapter/embeded"
	"github.com/siddeshwarnavink/UTA/adapter/keyExchange"
	"github.com/siddeshwarnavink/UTA/adapter/proxy"
	"github.com/siddeshwarnavink/UTA/adapter/ui"
	"github.com/siddeshwarnavink/UTA/shared/p2p"
	"github.com/siddeshwarnavink/UTA/shared/utils"

	"github.com/coreos/go-iptables/iptables"
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

	// UDP Peer table for discovery
	peerTable := p2p.NewPeerTable()

	if flags.Mode == ui.Client {
		go p2p.AnnouncePresence("adapter-client", flags.Dec, flags.Enc)
	} else {
		go p2p.AnnouncePresence("adapter-server", flags.Dec, flags.Enc)
	}
	go p2p.ListenForPeers(peerTable)

	// Setup TCP Transparent proxy
	setupIPTables(flags)

	switch flags.Mode {
	case ui.Client:
		ClientProxy(l, flags)
	case ui.Server:
		ServerProxy(l, flags)
	}
}

func setupIPTables(flags *ui.Flags) {
	ipt, err := iptables.New()
	if err != nil {
		panic(err)
	}

	chain := "UTA_PROXY"
	exists, err := ipt.ChainExists("mangle", chain)
	if err != nil {
		fmt.Println("Error checking if chain exists:", err)
		os.Exit(1)
	}

	if !exists {
		// Create our chain if it doesn't exist
		err := ipt.NewChain("mangle", chain)
		if err != nil {
			fmt.Println("Error creating new chain:", err)
			os.Exit(1)
		}
	}

	err = ipt.AppendUnique("mangle", "PREROUTING", "-j", chain)
	if err != nil {
		fmt.Println("Error appending rule to iptables:", err)
		os.Exit(1)
	}


	err = ipt.AppendUnique("mangle", "PREROUTING", "-j", chain)
	if err != nil {
		panic(err)
	}

	var encPort int
	var decPort int

	encPort, err = utils.ExtractPort(flags.Enc)
	if err != nil {
		panic(err)
	}

	decPort, err = utils.ExtractPort(flags.Dec)
	if err != nil {
		panic(err)
	}

	if flags.Mode == ui.Server {
		// Mark packets coming to server
		err = ipt.AppendUnique("mangle", "PREROUTING", "-p", "tcp", "--dport", strconv.Itoa(encPort), "-j", chain, "--tproxy-mark", "1/1")
		if err != nil {
			panic(err)
		}
		// Redirect them to our adapter
		err = ipt.AppendUnique("mangle", "UTA_PROXY", "-j", chain, "--tproxy-match", "MARK", "set", "1", "1", "--on-port", strconv.Itoa(decPort))
		if err != nil {
			panic(err)
		}
	} else {
		// Mark packets coming from the server IP
		err = ipt.AppendUnique("mangle", "PREROUTING", "-p", "tcp", "-s", strconv.Itoa(encPort), "-j", "MARK", "--set-mark", "1/1")
		if err != nil {
			panic(err)
		}
		// Redirect them to our adapter
		err = ipt.AppendUnique("mangle", "PREROUTING", "-p", "tcp", "-m", "mark", "--mark", "1/1", "-j", chain, "--on-port", strconv.Itoa(decPort))
		if err != nil {
			panic(err)
		}
	}
}

func ClientProxy(l *lua.LState, flags *ui.Flags) {
	decAddress := flags.Dec
	encAddress := flags.Enc

	decListener, err := net.Listen("tcp", decAddress)
	if err != nil {
		log.Fatalf("Could not start client adapter: %v", err)
	}
	defer decListener.Close()

	fmt.Printf("Client adapter listening on %s, forwarding to server %s\n",
		decAddress, encAddress)

	for {
		decConn, err := decListener.Accept()
		if err != nil {
			log.Printf("Could not accept connection: %v", err)
			continue
		}

		encConn, err := net.Dial("tcp", encAddress)
		if err != nil {
			log.Printf("Could not connect to server: %v", err)
			return
		}
		defer encConn.Close()

		derivedKey, err := keyExchange.ClientKeyExchange(encConn, flags.Protocol)

		if !proxy.IsUninitialized(derivedKey) {
			fmt.Printf("\nGot p2p key %x\n", derivedKey)

			algo, err := ui.AlgorithmFromString(flags.Algo)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			go proxy.ProxyHandler(decConn, encConn, derivedKey, algo)
		}
	}
}

func ServerProxy(l *lua.LState, flags *ui.Flags) {
	encAddress := flags.Enc
	decAddress := flags.Dec

	encListener, err := net.Listen("tcp", encAddress)
	if err != nil {
		log.Fatalf("Could not start server proxy: %v", err)
	}
	defer encListener.Close()

	fmt.Printf("Server adapter listening on %s, forwarding to server %s\n",
		encAddress, decAddress)

	for {
		encConn, err := encListener.Accept()
		if err != nil {
			log.Printf("Could not accept connection: %v", err)
			continue
		}

		decConn, err := net.Dial("tcp", decAddress)
		if err != nil {
			log.Printf("Could not connect to server: %v", err)
			return
		}
		defer decConn.Close()

		derivedKey, err := keyExchange.ServerKeyExchange(encConn, flags.Protocol)

		if !proxy.IsUninitialized(derivedKey) {
			fmt.Printf("\nGot p2p key %x\n", derivedKey)

			algo, err := ui.AlgorithmFromString(flags.Algo)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			go proxy.ProxyHandler(decConn, encConn, derivedKey, algo)
		}
	}
}
