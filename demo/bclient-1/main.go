package main

// import (
// 	"fmt"
// 	"net"
// 	"os"
// 	"reflect"
// 	"sync"
// 	"time"

// 	"github.com/bxcodec/faker/v4"
// )

// func randomName() string {
// 	firstName := faker.FirstName()
// 	lastName := faker.LastName()

// 	fullName := fmt.Sprintf("%s %s", firstName, lastName)
// 	return fullName
// }

// func handleConnection(conn net.Conn, wg *sync.WaitGroup) {
// 	defer wg.Done()
// 	defer conn.Close()

// 	buffer := make([]byte, 1024)
// 	for {
// 		n, err := conn.Read(buffer)
// 		if err != nil {
// 			fmt.Println("Error reading from connection:", err)
// 			return
// 		}
// 		fmt.Printf("Received from server: %s\n", string(buffer[:n]))

// 		randomName := randomName()
// 		_, err = conn.Write([]byte(randomName + "\n"))
// 		if err != nil {
// 			fmt.Println("Error writing to connection:", err)
// 			return
// 		}
// 	}
// }

// func main() {
// 	// local := flag.Bool("local", false, "Run outside docker")
// 	// flag.Parse()
// 	var port = "8888"
// 	var addr = "0.0.0.0:"

// 	args := os.Args[1:]
// 	for i, arg := range args {
// 		if arg == "--local" {
// 			addr = "127.0.0.1:"
// 		}
// 		if arg == "--port" {
// 			port = args[i+1]
// 			fmt.Println(reflect.TypeOf(port))
// 		}
// 	}

// 	addr = addr + port

// 	var listener net.Listener
// 	var err error

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

// 	fmt.Printf("Client is listening to %s\n", addr)

// 	var wg sync.WaitGroup

// 	for {
// 		conn, err := listener.Accept()
// 		if err != nil {
// 			fmt.Println("Error accepting connection:", err)
// 			continue
// 		}

// 		wg.Add(1)
// 		go handleConnection(conn, &wg)
// 	}

// 	wg.Wait()
// }
