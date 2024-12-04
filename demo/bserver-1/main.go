package main

// import (
// 	"fmt"
// 	"net"
// 	"os"
// 	"sync"
// 	"time"
// )

// func handleConnection(conn net.Conn, wg *sync.WaitGroup, clients *sync.Map) {
// 	defer wg.Done()
// 	defer conn.Close()

// 	clients.Store(conn.RemoteAddr().String(), conn)
// 	defer clients.Delete(conn.RemoteAddr().String())

// 	go func() {
// 		for {
// 			message := "Hello, this is the broadcast server!\n"
// 			_, err := conn.Write([]byte(message))
// 			if err != nil {
// 				fmt.Println("Error writing to client:", err)
// 				return
// 			}
// 			time.Sleep(5 * time.Second) // send every 5 seconds
// 		}
// 	}()

// 	buffer := make([]byte, 1024)
// 	for {
// 		Received, err := conn.Read(buffer)
// 		if err != nil {
// 			fmt.Println("Error reading from client:", err)
// 			return
// 		}
// 		fmt.Printf("Received from client: %s\n", string(buffer[:Received]))
// 	}
// }

// func broadcastToClients(clients *sync.Map, message string) {
// 	// Broadcast message to all connected clients
// 	clients.Range(func(key, value interface{}) bool {
// 		clientConn := value.(net.Conn)
// 		_, err := clientConn.Write([]byte(message))
// 		if err != nil {
// 			fmt.Println("Error broadcasting to client:", err)
// 			return false
// 		}
// 		return true
// 	})
// }

// func main() {
// 	var port = "8888"
// 	var addr = "0.0.0.0:" + port

// 	args := os.Args[1:]
// 	for i, arg := range args {
// 		if arg == "--local" {
// 			addr = "127.0.0.1:" + port
// 		}
// 		if arg == "--port" {
// 			port = args[i+1]
// 			addr = "0.0.0.0:" + port
// 		}
// 	}

// 	var listener net.Listener
// 	var err error

// 	// Start listening on the specified address and port
// 	for {
// 		listener, err = net.Listen("tcp", addr)
// 		if err != nil {
// 			fmt.Println("Error starting the server:", err)
// 			fmt.Println("Retrying in 5 seconds...")
// 			time.Sleep(5 * time.Second)
// 			continue
// 		}
// 		break
// 	}
// 	defer listener.Close()

// 	fmt.Printf("Broadcast server is listening on %s\n", addr)

// 	var wg sync.WaitGroup
// 	clients := &sync.Map{} // Thread-safe map to store connected clients

// 	// Start a goroutine to broadcast message every 10 seconds
// 	go func() {
// 		for {
// 			message := "Broadcast message to all clients!\n"
// 			broadcastToClients(clients, message)
// 			time.Sleep(10 * time.Second) // broadcast every 10 seconds
// 		}
// 	}()

// 	// Accept incoming client connections
// 	for {
// 		conn, err := listener.Accept()
// 		if err != nil {
// 			fmt.Println("Error accepting connection:", err)
// 			continue
// 		}

// 		wg.Add(1)
// 		go handleConnection(conn, &wg, clients)
// 	}

// 	wg.Wait()
// }
