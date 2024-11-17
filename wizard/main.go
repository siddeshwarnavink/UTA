package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/siddeshwarnavink/UTA/shared/p2p"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const wizardPort = 3300 // TODO: Make this dynamic via flag

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096, // read buffer size to 4KB
	WriteBufferSize: 4096, // write buffer size to 4KB
	CheckOrigin: func(r *http.Request) bool {
		return true // WARNING: this is a bad idea man
	},
}

func main() {
	peerTable := p2p.NewPeerTable()
	go p2p.AnnouncePresence("wizard", "", "")
	go p2p.ListenForPeers(peerTable)

	r := gin.Default()

	wsGroup := r.Group("/ws")
	{
		wsGroup.GET("/peer-table", func(c *gin.Context) {
			handleWebSocket(c.Writer, c.Request, peerTable)
		})
	}

	r.Static("/static", "./wizard/dist/static")

	r.NoRoute(func(c *gin.Context) {
		c.File("./wizard/dist/index.html")
	})

	r.Run(fmt.Sprintf(":%d", wizardPort))
}

func handleWebSocket(w http.ResponseWriter, r *http.Request, pt *p2p.PeerTable) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}
	defer conn.Close()

	err = sendRoutingTable(conn, pt)
	if err != nil {
		log.Println("Error sending routing table:", err)
		return
	}

	for {
		time.Sleep(5 * time.Second)

		err := sendRoutingTable(conn, pt)
		if err != nil {
			log.Println("Error sending updated routing table:", err)
			break
		}
	}
}

func sendRoutingTable(conn *websocket.Conn, pt *p2p.PeerTable) error {
	routingTable := pt.GetRoutingTable()

	err := conn.WriteJSON(routingTable)
	if err != nil {
		log.Println("Error sending routing table:", err)
	}
	return err
}
