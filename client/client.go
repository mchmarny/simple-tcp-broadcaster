package client

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/mchmarny/simple-server/types"
)

// StartClient starts a client and connects to server
// example "127.0.0.1:9999"
func StartClient(serverAddress string) error {

	c, err := net.Dial("tcp", serverAddress)
	if err != nil {
		return err
	}
	log.Printf("Connected to server: %s", serverAddress)

	for {
		sendScanner := bufio.NewScanner(os.Stdin)
		for sendScanner.Scan() {

			req := types.NewRequest(c.LocalAddr().String())
			req.Data = sendScanner.Bytes()

			_, err := fmt.Fprint(c, req)

			// _, err := fmt.Fprintf(c, text+"\n")
			if err != nil {
				return err
			}
			log.Printf("Server received: %s", req.Data)
			break
		}

		// listen
		listenScanner := bufio.NewScanner(c)
		for listenScanner.Scan() {
			log.Println("Server sends: " + listenScanner.Text())
			break
		}
		if err := listenScanner.Err(); err != nil {
			return err
		}
	}
}
