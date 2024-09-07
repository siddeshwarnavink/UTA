package main

import (
	"fmt"
	"net"
	"os"
	"sync"
)

func handleConnection(conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	defer conn.Close()

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading from connection:", err)
			return
		}
		fmt.Printf("Received from client: %s\n", string(buffer[:n]))

		// Echo the message back to the client
		_, err = conn.Write(buffer[:n])
		if err != nil {
			fmt.Println("Error writing to connection:", err)
			return
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", "localhost:10000")
	if err != nil {
		fmt.Println("Error starting the server:", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("Server is listening on port 10000")

	var wg sync.WaitGroup

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		wg.Add(1)
		go handleConnection(conn, &wg)
	}

	wg.Wait()
}
