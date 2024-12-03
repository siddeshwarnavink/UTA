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

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // WARN: This is a bad idea
	},
}

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

var transmissions = sync.Map{} // ip: boolean

func main() {
	peerTable := p2p.NewPeerTable()

	peerConn, err := p2p.GetMulticastConn()
	if err != nil {
		panic(err)
	}
	defer peerConn.Close()

	go p2p.AnnouncePresence(*peerConn, p2p.Wizard, "", "")
	transCh, resChan := p2p.ListenForPeers(*peerConn, p2p.Wizard, peerTable)

	r := gin.Default()

	r.GET("/ws", func(c *gin.Context) {
		if !websocket.IsWebSocketUpgrade(c.Request) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Not a WebSocket request"})
			return
		}

		transmissionChan := make(chan p2p.TransmissionMsg)
		msgChan := make(chan WsData)

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade connection"})
			return
		}
		defer conn.Close()

		// if transmissionChan is idle for 5 seconds flush cache
		idleTimeout := 5 * time.Second
		go func() {
			timer := time.NewTimer(idleTimeout)
			defer timer.Stop()

			for {
				select {
				case <-transmissionChan:
					// if !timer.Stop() {
					// 	<-timer.C
					// }
					timer.Reset(idleTimeout)
				case <-timer.C:
					cleanupIp := []string{}
					transmissions.Range(func(ip, _ interface{}) bool {
						val := p2p.TransmissionMsg{IP: ip.(string), Send: false}
						msgChan <- WsData{Transmission: &val, PeerTable: nil}

						cleanupIp = append(cleanupIp, ip.(string))
						return true
					})
					for _, ip := range cleanupIp {
						transmissions.Delete(ip)
					}
					timer.Reset(idleTimeout)
					return
				}
			}
		}()

		// send data transfers
		go func() {
			for val := range transCh {
				transmissions.Store(val.IP, val.Send)
				transmissionChan <- val
				msgChan <- WsData{Transmission: &val, PeerTable: nil}
			}
		}()

		// listen for request message
		go func() {
			for {
				_, msg, err := conn.ReadMessage()
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

		// send response
		go func() {
			for val := range resChan {
				msgChan <- WsData{Response: &val}
				time.Sleep(5 * time.Second)
			}
		}()

		// send routing table every 5 seconds
		go func() {
			for {
				routingTable := peerTable.GetRoutingTable()
				msgChan <- WsData{PeerTable: &routingTable, Transmission: nil}
				time.Sleep(5 * time.Second)
			}
		}()

		// sender junction
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
