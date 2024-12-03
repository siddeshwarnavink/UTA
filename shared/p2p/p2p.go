package p2p

import (
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"time"

	"github.com/jedib0t/go-pretty/table"
)

type PeerRole string

const (
	InvalidRole PeerRole = "invalid"
	ClientProxy PeerRole = "adapter-client"
	ServerProxy PeerRole = "adapter-server"
	Wizard      PeerRole = "wizard"
)

const multicastAddr = "224.0.0.1:9999" // TODO: Make this dynamic
const heartbeatInterval = 5 * time.Second
const peerTimeout = 10 * time.Second // ideally it should be 30 sec

type Peer struct {
	IP       string    `json:"ip"`
	FromIP   string    `json:"from_ip"`
	ToIP     string    `json:"to_ip"`
	Role     PeerRole  `json:"role"`
	LastSeen time.Time `json:"last_seen"`
}

type TransmissionMsg struct {
	IP   string `json:"ip"`
	Send bool   `json:"sent"`
}

type WsResponseMsg struct {
	RequestId string `json:"reqId"`
	Data      string `json:"data"`
}

type PeerTable struct {
	mu    sync.Mutex
	peers map[string]Peer
}

func NewPeerTable() *PeerTable {
	return &PeerTable{
		peers: make(map[string]Peer),
	}
}

// use this for dialing only
func GetMulticastConn() (*net.UDPConn, error) {
	addr, err := net.ResolveUDPAddr("udp", multicastAddr)
	if err != nil {
		return nil, fmt.Errorf("Error resolving UDP address: %v\n", err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, fmt.Errorf("Error dialing UDP connection: %v\n", err)
	}

	return conn, nil
}

func AnnouncePresence(conn net.UDPConn, role PeerRole, fromIP, toIP string) {
	for {
		message, err := DiscoveryMessage(role, fromIP, toIP)
		if err != nil {
			fmt.Printf("Error in encode discovery message: %v\n", err)
			return
		}

		messageBytes := []byte(message)

		conn.Write(messageBytes)
		time.Sleep(heartbeatInterval)
	}
}

func ListenForPeers(peerConn net.UDPConn, role PeerRole, peerTable *PeerTable) (chan TransmissionMsg, chan WsResponseMsg) {
	transCh := make(chan TransmissionMsg)
	resCh := make(chan WsResponseMsg)

	go func() {
		addr, err := net.ResolveUDPAddr("udp", multicastAddr)
		if err != nil {
			panic(err)
		}

		conn, err := net.ListenMulticastUDP("udp", nil, addr)
		if err != nil {
			panic(err)
		}
		defer conn.Close()

		buf := make([]byte, 6040)
		for {
			n, src, err := conn.ReadFromUDP(buf)
			if err != nil {
				fmt.Printf("Error reading from UDP: %v\n", err)
				continue
			}
			address := src.String()

			message := string(buf[:n])

			msgtype, err := GetPeerMsgType(message)
			if err != nil {
				fmt.Printf("Invalid peer message: %s", message)
			}

			// fmt.Printf("UDP(%d)=%s\n", msgtype, message)

			peerTable.updatePeerTable(address, message, msgtype)

			// Adapter getting a request
			if role != Wizard && msgtype == StringRequestMessageType {
				_, reqType, reqId, payload, err := ExtractRequestMessage(message)
				if err == nil {
					// if that is me, send response
					if payload == peerConn.LocalAddr().String() {
						fmt.Println("That request is for me")
						// TODO: validate if valid wizard
						// send response
						switch reqType {
						case RequestTypeConfig:
							fmt.Println("asking for my config")

							// TODO: Read config file from flag
							data, err := os.ReadFile("./adapter/config/init.lua")
							if err != nil {
								fmt.Println("Error reading config file:", err)
								return
							}
							configContent := string(data)

							resMsg, err := ResponseMessage(role, reqId, configContent)
							if err != nil {
								fmt.Println("Invalid response message:", err)
								return
							}

							peerConn.Write([]byte(resMsg))
							break
						case RequestTypeLogs:
							fmt.Println("asking for my logs")

							file, err := os.Open("./logs/adapter.log")
							if err != nil {
								fmt.Println(err)
								return
							}
							defer file.Close()

							const chunkSize = 500
							buf := make([]byte, chunkSize)

							_, err = file.Read(buf)
							if err != nil && err != io.EOF {
								fmt.Println(err)
								return
							}

							resMsg, err := ResponseMessage(role, reqId, string(buf))

							fmt.Printf("my response=%s\n", resMsg)
							if err != nil {
								fmt.Println("Invalid response message:", err)
								return
							}

							peerConn.Write([]byte(resMsg))
							break
						default:
							fmt.Println("Invalid request sent to me")
							break
						}
					} else {
						fmt.Printf("Not for me, I am %s\n", peerConn.LocalAddr().String())
					}
				} else {
					fmt.Println("It was a invalid request")
				}
			}

			// Wizard getting response from adapters
			if role == Wizard && msgtype == StringResponseMessageType {
				fmt.Printf("Got response in UDP %s\n", message)
				_, reqId, data, err := ExtractResponseMessage(message)
				if err == nil {
					obj := WsResponseMsg{
						RequestId: reqId,
						Data:      data,
					}
					resCh <- obj
				} else {
					fmt.Printf("That response had errors %v\n", err)
				}
			}

			// Wizard getting transmission message
			if role == Wizard && msgtype == Transmission {
				_, sent, err := ExtractTransmissionMessageDetails(message)
				if err == nil {
					chmsg := TransmissionMsg{
						IP:   address,
						Send: sent,
					}
					transCh <- chmsg
				}
			}
		}
	}()

	return transCh, resCh
}

func (pt *PeerTable) cleanupInactivePeers() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		pt.mu.Lock()
		for addr, peer := range pt.peers {
			if time.Since(peer.LastSeen) > peerTimeout {
				delete(pt.peers, addr)
				pt.PrintRoutingTable()
			}
		}
		pt.mu.Unlock()
	}
}

func (pt *PeerTable) updatePeerTable(address string, message string, msgtype PeerMsgType) {
	if msgtype == Discovery {
		pt.mu.Lock()
		defer pt.mu.Unlock()

		role, fromIP, toIP, err := ExtractDiscoveryMessageDetails(string(message))

		if err == nil {
			_, exists := pt.peers[address]

			pt.peers[address] = Peer{
				IP:       address,
				Role:     role,
				FromIP:   fromIP,
				ToIP:     toIP,
				LastSeen: time.Now(),
			}

			if !exists {
				fmt.Printf("Discovered new peer: %s\n", address)
				pt.PrintRoutingTable()
			}
		} else {
			fmt.Print(err)
		}
	}
}

func (pt *PeerTable) GetPeers() []Peer {
	pt.mu.Lock()
	defer pt.mu.Unlock()

	var peers []Peer
	for _, peer := range pt.peers {
		if time.Since(peer.LastSeen) < 10*heartbeatInterval {
			peers = append(peers, peer)
		}
	}
	return peers
}

func (pt *PeerTable) PrintRoutingTable() {
	// pt.mu.Lock()
	// defer pt.mu.Unlock()

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{
		"IP",
		"Role",
		"From",
		"To",
	})

	for _, peer := range pt.peers {
		t.AppendRow(table.Row{
			peer.IP,
			peer.Role,
			peer.FromIP,
			peer.ToIP,
		})
	}
	t.Render()
}

func (pt *PeerTable) GetRoutingTable() map[string]Peer {
	pt.mu.Lock()
	defer pt.mu.Unlock()

	routingTable := make(map[string]Peer)
	for k, v := range pt.peers {
		routingTable[k] = v
	}
	return routingTable
}
