package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	serverIP := "server" // name in docker-compose
	serverPort := 10000
	serverAddress := fmt.Sprintf("%s:%d", serverIP, serverPort)

	// Connect to the server
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	fmt.Printf("Connected to server at %s\n", serverAddress)

	message := "Hello from client!"

	for i := 0; i < 10; i++ {
		_, err := conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Error sending message:", err)
			break
		}
		fmt.Printf("Sent: %s\n", message)

		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading response:", err)
			break
		}
		response := string(buffer[:n])
		fmt.Printf("Received: %s\n", response)

		time.Sleep(1 * time.Second)
	}

	fmt.Println("Connection closed")
}

