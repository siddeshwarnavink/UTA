package main

import (
	"fmt"
	"io"
	"log"
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

	//doing lua stuff
	configPath := ui.GetConfigFile()
	fmt.Printf("Using config file: %s\n", configPath)
	embeded.HandleLua(l, configPath)

	//render TUI if config not in lua
	ui.RenderForm()

	peerTable := p2p.NewPeerTable()

	peerConn, err := p2p.GetMulticastConn()
	if err != nil {
		panic(err)
	}
	defer peerConn.Close()

	if embeded.CurrentFlags.Mode == embeded.Client {
		go p2p.AnnouncePresence(*peerConn, p2p.ClientProxy, embeded.CurrentFlags.Dec, embeded.CurrentFlags.Enc)
	} else {
		go p2p.AnnouncePresence(*peerConn, p2p.ServerProxy, embeded.CurrentFlags.Dec, embeded.CurrentFlags.Enc)
	}
	p2p.ListenForPeers(peerTable)

	switch embeded.CurrentFlags.Mode {
	case embeded.Client:
		proxy.ClientProxy(l, *peerConn)
	case embeded.Server:
		proxy.ServerProxy(l, *peerConn)
	}
}
