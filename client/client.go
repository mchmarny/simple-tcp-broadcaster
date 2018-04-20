package client

import (
	"bufio"
	"log"
	"net"
	"os"
	"time"

	"github.com/mchmarny/simple-tcp-broadcaster/commons"
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

	log.Printf("Connected to server %s from %s",
		conn.RemoteAddr(), conn.LocalAddr())

	tcpConn := conn.(*net.TCPConn)
	tcpConn.SetKeepAlive(true)
	tcpConn.SetKeepAlivePeriod(time.Second * time.Duration(keepAliveEverySeconds))

	agent := commons.NewClientAgent(conn)
	log.Printf("Client ID: %s", agent.GetLocalClientID())

	go agent.Read()

	for {
		message, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		msg := commons.NewMessage(agent.GetLocalClientID())
		msg.Data = []byte(message)
		if err := agent.Write(msg); err != nil {
			log.Fatalf("Error on write: %v", err)
		}

	}

}
