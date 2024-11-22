package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/siddeshwarnavink/UTA/shared/p2p"

	"github.com/gin-gonic/gin"
)

const wizardPort = 3300 // TODO: Make this dynamic via flag

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // WARN: This is a bad idea
	},
}

func main() {
	peerTable := p2p.NewPeerTable()

	go p2p.AnnouncePresence(p2p.Wizard, "", "")
	go p2p.ListenForPeers(peerTable)

	r := gin.Default()

	r.GET("/ws", func(c *gin.Context) {
		if !websocket.IsWebSocketUpgrade(c.Request) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Not a WebSocket request"})
			return
		}

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade connection"})
			return
		}
		defer conn.Close()

		for {
			routingTable := peerTable.GetRoutingTable()
			data, err := json.Marshal(routingTable)
			if err != nil {
				fmt.Printf("Error in converting to JSON: %v\n", err)
				break
			}

			if err := conn.WriteMessage(websocket.TextMessage, []byte(data)); err != nil {
				fmt.Printf("Error in writing message: %v\n", err)
				break
			}
			time.Sleep(5 * time.Second)
		}

	})

	r.Static("/static", "./wizard/dist/static")

	r.NoRoute(func(c *gin.Context) {
		c.File("./wizard/dist/index.html")
	})

	r.Run(fmt.Sprintf(":%d", wizardPort))

	fmt.Printf("Wizard running on http://localhost:%d\n", wizardPort)
}
