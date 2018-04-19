package client

import (
	"bufio"
	"log"
	"net"
	"os"

	"github.com/mchmarny/simple-server/commons"
)

// StartClient starts a client and connects to server
func StartClient(serverAddress string) error {

	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		log.Fatalf("Error on dial: %v", err)
		return err
	}

	log.Printf("Connected to server: %s", conn.RemoteAddr())

	client := commons.NewConnection(conn)
	go client.Read()

	for {
		message, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		msg := commons.NewMessage(client.Socket.LocalAddr().String())
		msg.Data = []byte(message)
		client.Write(msg)
	}

}
