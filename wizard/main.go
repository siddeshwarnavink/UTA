package main

import (
	"fmt"

	"github.com/siddeshwarnavink/UTA/shared/p2p"

	"github.com/gin-gonic/gin"
)

const wizardPort = 3300 // TODO: Make this dynamic via flag

func main() {
	peerTable := p2p.NewPeerTable()
	go p2p.AnnouncePresence(fmt.Sprintf("127.0.0.1:%d", wizardPort))
	go p2p.ListenForPeers(peerTable)

	r := gin.Default()
	r.Static("/", "./wizard/dist")

	r.Run(fmt.Sprintf(":%d", wizardPort))
}
