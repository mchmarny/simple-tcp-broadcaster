package client

import (
	"bufio"
	"log"
	"net"
	"os"
	"time"

	"github.com/mchmarny/simple-server/commons"
)

const (
	keepAliveEverySeconds = 10
)

// StartClient starts a client and connects to server
func StartClient(serverAddress string) error {

	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		log.Fatalf("Error on dial: %v", err)
		return err
	}

	log.Printf("Connected to server: %s", conn.RemoteAddr())

	tcpConn := conn.(*net.TCPConn)
	tcpConn.SetKeepAlive(true)
	tcpConn.SetKeepAlivePeriod(time.Second * time.Duration(keepAliveEverySeconds))

	client := commons.NewClientConnection(conn)

	go client.Read()

	for {
		message, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		msg := commons.NewMessage(client.Socket.LocalAddr().String())
		msg.Data = []byte(message)
		if err := client.Write(msg); err != nil {
			log.Fatalf("Error on write: %v", err)
		}

	}

}
