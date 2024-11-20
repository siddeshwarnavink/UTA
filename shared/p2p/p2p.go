package p2p

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/jedib0t/go-pretty/table"
)

const multicastAddr = "224.0.0.1:9999" // TODO: Make this dynamic
const heartbeatInterval = 5 * time.Second
const peerTimeout = 10 * time.Second // ideally it should be 30 sec

type Peer struct {
	IP       string // UDP IP of that peer
	FromIP   string // i.e Dec flag
	ToIP     string // i.e Enc flag
	Role     string // "adapter-client", "adapter-server", "wizard"
	LastSeen time.Time
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

func AnnouncePresence(role string, fromIP, toIP string) {
	addr, err := net.ResolveUDPAddr("udp", multicastAddr)
	if err != nil {
		fmt.Printf("Error resolving UDP address: %v\n", err)
		return
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Printf("Error dialing UDP connection: %v\n", err)
		return
	}
	defer conn.Close()

	for {
		message := fmt.Sprintf("%s,%s,%s", role, fromIP, toIP)
		messageBytes := []byte(message)

		conn.Write(messageBytes)
		time.Sleep(heartbeatInterval)
	}
}

func ListenForPeers(peerTable *PeerTable) {
	go peerTable.cleanupInactivePeers()

	addr, err := net.ResolveUDPAddr("udp", multicastAddr)
	if err != nil {
		fmt.Printf("Error resolving UDP address: %v\n", err)
		return
	}

	conn, err := net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		fmt.Printf("Error listening multicast: %v\n", err)
		return
	}

	defer conn.Close()

	buf := make([]byte, 1024)
	for {
		n, src, _ := conn.ReadFromUDP(buf)
		peerTable.updatePeerTable(src.String(), buf[:n])
	}
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

func extractMessage(message string) (string, string, string, error) {
	parts := strings.Split(message, ",")
	if len(parts) != 3 {
		return "", "", "", errors.New("invalid message format: expected 3 parts")
	}

	role := parts[0]
	fromIP := parts[1]
	toIP := parts[2]

	return role, fromIP, toIP, nil
}

func (pt *PeerTable) updatePeerTable(address string, message []byte) {
	pt.mu.Lock()
	defer pt.mu.Unlock()

	role, fromIP, toIP, err := extractMessage(string(message))

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
