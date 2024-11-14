package p2p

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"github.com/jedib0t/go-pretty/table"
)

const multicastAddr = "224.0.0.1:9999" // TODO: Make this dynamic
const heartbeatInterval = 5 * time.Second
const peerTimeout = 30 * time.Second

type Peer struct {
	IP       string // UDP IP of that peer
	Address  string // TCP address of the system it is attached to (i.e the Dec flag)
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

func AnnouncePresence(port string) {
	addr, _ := net.ResolveUDPAddr("udp", multicastAddr)
	conn, _ := net.DialUDP("udp", nil, addr)
	defer conn.Close()

	for {
		message := []byte(port)
		conn.Write(message)
		time.Sleep(heartbeatInterval)
	}
}

func ListenForPeers(peerTable *PeerTable) {
	go peerTable.cleanupInactivePeers()

	addr, _ := net.ResolveUDPAddr("udp", multicastAddr)
	conn, _ := net.ListenMulticastUDP("udp", nil, addr)
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

func (pt *PeerTable) updatePeerTable(address string, message []byte) {
	pt.mu.Lock()
	defer pt.mu.Unlock()

	peerInfo := string(message)
	if peerInfo != "" {
		if _, exists := pt.peers[address]; !exists {
			pt.peers[address] = Peer{
				IP:       address,
				Address:  peerInfo,
				LastSeen: time.Now(),
			}
			pt.PrintRoutingTable()
		} else {
			pt.peers[address] = Peer{
				IP:       address,
				Address:  pt.peers[address].Address,
				LastSeen: time.Now(),
			}
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

	fmt.Print("\033[H\033[2J") // clear screen
	fmt.Println("Peer Table")

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{
		"Address",
	})

	for _, peer := range pt.peers {
		t.AppendRow(table.Row{
			peer.Address,
		})
	}
	t.Render()
}
