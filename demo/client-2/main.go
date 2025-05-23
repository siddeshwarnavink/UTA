package main

import (
	"flag"
	"fmt"
	"net"
	"time"
)

func main() {
	local := flag.Bool("local", false, "Run outside docker")
	flag.Parse()

	serverIP := "client-2"
	serverPort := 8888

	if *local {
		serverIP = "127.0.0.1"
	}

	serverAddress := fmt.Sprintf("%s:%d", serverIP, serverPort)

	var conn net.Conn
	var err error

	for {
		conn, err = net.Dial("tcp", serverAddress)
		if err != nil {
			fmt.Println("Error connecting to server:", err)
			fmt.Println("Retrying in 5 seconds...")
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}
	defer conn.Close()

	fmt.Printf("Connected to server at %s\n", serverAddress)

	message := "I'm client!"

	for i := 0; i < 100; i++ {
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
