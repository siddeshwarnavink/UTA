package p2p

import (
	"fmt"
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
	IP       string // UDP IP of that peer
	FromIP   string // i.e Dec flag
	ToIP     string // i.e Enc flag
	Role     PeerRole
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

func AnnouncePresence(role PeerRole, fromIP, toIP string) {
	conn, err := GetMulticastConn()
	if err != nil {
		panic(err)
	}
	defer conn.Close()

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

func ListenForPeers(peerTable *PeerTable) {
	addr, err := net.ResolveUDPAddr("udp", multicastAddr)
	if err != nil {
		panic(err)
	}

	conn, err := net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	buf := make([]byte, 1024)
	for {
		n, src, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Printf("Error reading from UDP: %v\n", err)
			continue
		}
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
	msgtype, err := GetPeerMsgType(string(message))
	if err != nil {
		fmt.Printf("Invalid peer message: %s", string(message))
	}

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
	} else if msgtype == Transmission  {
		fmt.Printf("Got Transmission message %s\n", string(message))
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
