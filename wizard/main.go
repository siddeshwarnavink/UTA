package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/siddeshwarnavink/UTA/shared/p2p"

	"github.com/gin-gonic/gin"
)

const wizardPort = 3300 // TODO: Make this dynamic via flag

type WsDataType int

const (
	PeerTableData    WsDataType = 1
	TransmissionData WsDataType = 2
	ResponseData     WsDataType = 3
)

type WsData struct {
	PeerTable    *map[string]p2p.Peer `json:"peerTable,omitempty"`
	Transmission *p2p.TransmissionMsg `json:"transmission,omitempty"`
	Response     *p2p.WsResponseMsg   `json:"response,omitempty"`
}

type RequestMessageJson struct {
	Type    p2p.RequestMessageType `json:"request"`
	ReqId   string                 `json:"reqId"`
	Payload string                 `json:"payload"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // WARN: This is a bad idea
	},
}

var broadcastChan = make(chan WsData)
var clients = make(map[*websocket.Conn]bool)
var clients_mutex = &sync.Mutex{}

func main() {
	peerTable := p2p.NewPeerTable()

	peerConn, err := p2p.GetMulticastConn()
	if err != nil {
		panic(err)
	}
	defer peerConn.Close()

	go p2p.AnnouncePresence(*peerConn, p2p.Wizard, "", "")
	transCh, resChan := p2p.ListenForPeers(*peerConn, p2p.Wizard, peerTable, "")

	r := gin.Default()

	// Send routing table every 5 seconds
	go func() {
		for {
			routingTable := peerTable.GetRoutingTable()
			broadcastChan <- WsData{PeerTable: &routingTable, Transmission: nil}
			time.Sleep(5 * time.Second)
		}
	}()

	// Send data transfers
	go func() {
		for val := range transCh {
			broadcastChan <- WsData{Transmission: &val, PeerTable: nil}
		}
	}()

	// Send response for client's request
	go func() {
		for val := range resChan {
			broadcastChan <- WsData{Response: &val}
			time.Sleep(1 * time.Second)
		}
	}()

	r.GET("/ws", func(c *gin.Context) {
		if !websocket.IsWebSocketUpgrade(c.Request) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Not a WebSocket request"})
			return
		}

		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade connection"})
			return
		}
		defer ws.Close()

		clients_mutex.Lock()
		clients[ws] = true
		clients_mutex.Unlock()

		defer delete(clients, ws)

		// Listen data from clients
		go func() {
			for {
				_, msg, err := ws.ReadMessage()
				if err != nil {
					if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
						fmt.Println("Unexpected WebSocket closure")
					}
					time.Sleep(5 * time.Second)
					continue
				}

				var reqObj RequestMessageJson
				err = json.Unmarshal(msg, &reqObj)
				if err != nil {
					fmt.Println("Invalid json request")
					continue
				}

				udpReq, err := p2p.RequestMessage(p2p.Wizard, reqObj.Type, reqObj.ReqId, reqObj.Payload)
				fmt.Printf("sending request %s\n", udpReq)
				peerConn.Write([]byte(udpReq))
				time.Sleep(1 * time.Second)
			}
		}()

		// Send data to all clients
		for {
			msg := <-broadcastChan
			json, err := json.Marshal(msg)
			if err != nil {
				fmt.Printf("Error in converting to JSON: %v\n", err)
				continue
			}

			clients_mutex.Lock()
			for client := range clients {
				err := client.WriteMessage(websocket.TextMessage, json)
				if err != nil {
					client.Close()
					delete(clients, client)
				}
			}
			clients_mutex.Unlock()
		}
	})

	r.Static("/static", "./wizard/dist/static")

	r.NoRoute(func(c *gin.Context) {
		c.File("./wizard/dist/index.html")
	})

	r.Run(fmt.Sprintf(":%d", wizardPort))

	fmt.Printf("Wizard running on http://localhost:%d\n", wizardPort)
}
