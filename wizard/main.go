package main

import (
	"fmt"
	"net/http"

	"github.com/siddeshwarnavink/UTA/shared/p2p"

	"github.com/gin-gonic/gin"
)

const wizardPort = 3300 // TODO: Make this dynamic via flag

func main() {
	peerTable := p2p.NewPeerTable()
	go p2p.AnnouncePresence("wizard", "", "")
	go p2p.ListenForPeers(peerTable)

	r := gin.Default()

	apiGroup := r.Group("/api")
	{
		apiGroup.GET("/peer-table", func(c *gin.Context) {
			routingTable := peerTable.GetRoutingTable()
			c.IndentedJSON(http.StatusOK, routingTable)
		})
	}

	r.Static("/static", "./wizard/dist/static")

	r.NoRoute(func(c *gin.Context) {
		c.File("./wizard/dist/index.html")
	})

	r.Run(fmt.Sprintf(":%d", wizardPort))

	fmt.Printf("Wizard running on http://localhost:%d\n", wizardPort)
}
