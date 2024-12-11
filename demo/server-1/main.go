package main

import (
	"flag"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/bxcodec/faker/v4"
)

func randomName() string {
	firstName := faker.FirstName()
	lastName := faker.LastName()

	fullName := fmt.Sprintf("%s %s", firstName, lastName)
	return fullName
}

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

		randomName := randomName()
		_, err = conn.Write([]byte(randomName + "\n"))
		if err != nil {
			fmt.Println("Error writing to connection:", err)
			return
		}
	}
}

func main() {
	local := flag.Bool("local", false, "Run outside docker")
	raw := flag.Bool("raw", false, "Run without adapter")
	flag.Parse()

	addr := "0.0.0.0:10000"

	if *local {
		addr = "127.0.0.1:10000"
	}


	if *raw {
		addr = "localhost:8888"
	}


	var listener net.Listener
	var err error

	for {
		listener, err = net.Listen("tcp", addr)
		if err != nil {
			fmt.Println("Error starting the server:", err)
			fmt.Println("Retrying in 5 seconds...")
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}
	defer listener.Close()

	fmt.Printf("Server is listening to %s\n", addr)

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
