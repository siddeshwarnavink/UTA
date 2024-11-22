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

type WsDataType int

const (
	PeerTableData    WsDataType = 1
	TransmissionData WsDataType = 2
)

type WsData struct {
	PeerTable    *map[string]p2p.Peer `json:"peerTable,omitempty"`
	Transmission *p2p.TransmissionMsg `json:"transmission,omitempty"`
}

func main() {
	peerTable := p2p.NewPeerTable()

	peerConn, err := p2p.GetMulticastConn()
	if err != nil {
		panic(err)
	}
	defer peerConn.Close()

	go p2p.AnnouncePresence(*peerConn, p2p.Wizard, "", "")
	ch := p2p.ListenForPeers(peerTable)

	r := gin.Default()

	r.GET("/ws", func(c *gin.Context) {
		if !websocket.IsWebSocketUpgrade(c.Request) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Not a WebSocket request"})
			return
		}

		msgChan := make(chan WsData)

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade connection"})
			return
		}
		defer conn.Close()

		go func() {
			for val := range ch {
				msgChan <- WsData{Transmission: &val, PeerTable: nil}
			}
		}()

		go func() {
			for {
				routingTable := peerTable.GetRoutingTable()
				msgChan <- WsData{PeerTable: &routingTable, Transmission: nil}
				time.Sleep(5 * time.Second)
			}
		}()

		for {
			select {
			case msg := <-msgChan:
				json, err := json.Marshal(msg)
				if err != nil {
					fmt.Printf("Error in converting to JSON: %v\n", err)
					break
				}

				if err := conn.WriteMessage(websocket.TextMessage, json); err != nil {
					fmt.Printf("Error in writing message: %v\n", err)
					break
				}
			}
		}
	})

	r.Static("/static", "./wizard/dist/static")

	r.NoRoute(func(c *gin.Context) {
		c.File("./wizard/dist/index.html")
	})

	r.Run(fmt.Sprintf(":%d", wizardPort))

	fmt.Printf("Wizard running on http://localhost:%d\n", wizardPort)
}
