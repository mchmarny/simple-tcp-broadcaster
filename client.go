package main

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

var (
	agent            *commons.Connection
	inspectionPeriod = time.Second * keepAliveEverySeconds
)

// StopClient stops client
func StopClient() error {
	if agent != nil {
		return agent.Stop()
	}
	return nil
}

func watchAgent() {
	inspectoin := time.After(inspectionPeriod)
	for {
		select {
		case <-inspectoin:
			if agent != nil {
				err := agent.WriteHeartbeat()
				if err != nil {
					log.Fatalf("Error on heartbeat: %v", err)
				}
			}
			inspectoin = time.After(inspectionPeriod)
		}
	}
}

// StartClient starts a client and connects to server
func StartClient(serverAddress string) error {

	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		log.Fatalf("Error on dial: %v", err)
		return err
	}

	log.Printf("Connected to server %s from %s", conn.RemoteAddr(), conn.LocalAddr())

	tcpConn := conn.(*net.TCPConn)
	tcpConn.SetKeepAlive(true)
	tcpConn.SetKeepAlivePeriod(time.Second * time.Duration(keepAliveEverySeconds))

	agent := commons.NewClientConnection(conn)

	log.Printf("Client ID: %s", agent.GetLocalClientID())
	go agent.Read()
	go watchAgent()

	for {
		message, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		msg := commons.NewMessage(
			agent.GetLocalClientID(),
			[]byte(message),
		)
		if err := agent.Write(msg); err != nil {
			log.Fatalf("Error on write: %v", err)
		}

	}

}
