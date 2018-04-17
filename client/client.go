package client

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

// StartClient starts a client and connects to server
// example "127.0.0.1:9999"
func StartClient(serverAddress string) error {

	i := 0 // connection index
	// connect to server
	for {
	connect:
		c, errConn := net.Dial("tcp", serverAddress)
		if errConn != nil {
			continue
		} else {
			i++
			if i <= 1 {
				log.Println("Connected to server...")
				fmt.Println("---")
			} else {
				log.Println("Reconnected to server...")
				fmt.Println("---")
			}
		}
		for {
			// read in input from stdin
			scannerStdin := bufio.NewScanner(os.Stdin)
			fmt.Print("Server message: ")
			for scannerStdin.Scan() {
				text := scannerStdin.Text()
				fmt.Println("---")
				// send to server
				_, errWrite := fmt.Fprintf(c, text+"\n")
				if errWrite != nil {
					log.Println("Server offline, attempting to reconnect...")
					goto connect
				}
				log.Print("Server receives: " + text)
				break
			}
			// listen for reply
			scannerConn := bufio.NewScanner(c)
			for scannerConn.Scan() {
				log.Println("Server sends: " + scannerConn.Text())
				break
			}
			if errReadConn := scannerStdin.Err(); errReadConn != nil {
				log.Printf("Read error: %T %+v", errReadConn, errReadConn)
				return errReadConn
			}
			fmt.Println("---")
		}
	}

}
