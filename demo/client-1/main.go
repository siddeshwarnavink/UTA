package main

import (
	"flag"
	"fmt"
	"net"
	"time"
)

func main() {
	local := flag.Bool("local", false, "Run outside docker")
	raw := flag.Bool("raw", false, "Run without adapter")
	flag.Parse()

	serverIP := "client-1"
	serverPort := 8888

	if *local {
		serverIP = "127.0.0.1"
	}

	if *raw {
		serverIP = "localhost"
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

	message := []byte("I'm client!")
	totalMessages := 10
	successfulMessages := 0
	var totalLatency time.Duration

	startTime := time.Now()

	for i := 0; i < totalMessages; i++ {
		writeStart := time.Now()
		_, err := conn.Write(message)
		writeTime := time.Since(writeStart)

		if err != nil {
			fmt.Println("Error sending message:", err)
			continue
		}

		buffer := make([]byte, 1024)
		readStart := time.Now()
		n, err := conn.Read(buffer)
		readTime := time.Since(readStart)

		if err != nil {
			fmt.Println("Error reading response:", err)
			continue
		}

		roundTripTime := writeTime + readTime
		totalLatency += roundTripTime
		successfulMessages++

		response := string(buffer[:n])
		fmt.Printf("Sent: %s | Received: %s | Write Time: %v | Read Time: %v | RTT: %v\n",
			message, response, writeTime, readTime, roundTripTime)

		time.Sleep(1 * time.Second)
	}

	endTime := time.Now()
	totalDuration := endTime.Sub(startTime)
	throughput := float64(successfulMessages) / totalDuration.Seconds()
	packetLoss := float64(totalMessages-successfulMessages) / float64(totalMessages) * 100

	fmt.Println("=== Metrics ===")
	fmt.Printf("Total Messages Sent: %d\n", totalMessages)
	fmt.Printf("Successful Messages: %d\n", successfulMessages)
	fmt.Printf("Packet Loss: %.2f%%\n", packetLoss)
	fmt.Printf("Average RTT: %v\n", totalLatency/time.Duration(successfulMessages))
	fmt.Printf("Throughput: %.2f messages/second\n", throughput)
	fmt.Println("================")

}
