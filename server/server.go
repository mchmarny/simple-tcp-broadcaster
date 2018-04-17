package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

const (
	responseMsgTemplate = "Hi there, I got your message: %s"
	errorTemplate       = "Closing client %d connection returned error: %v\n"
	successTemplate     = "Client %d\n received: %s\n returned: %s\n"
)

func handleConnection(c net.Conn, i int) {
	defer func() {
		if err := c.Close(); err != nil {
			log.Printf(errorTemplate, i, err)
		}
	}()

	scanner := bufio.NewScanner(c)
	for scanner.Scan() {
		msg := scanner.Text()
		msgNew := fmt.Sprintf(responseMsgTemplate, msg)
		if _, err := c.Write([]byte(msgNew + "\n")); err != nil {
			log.Printf("Client %d has closed writer: %v\n", i, err)
			break
		}
		log.Printf(successTemplate, i, msg, msgNew)
	}
	log.Printf("Client %d disconnected...\n", i)
}

// StartServer start TCP server on a given port
func StartServer(port int) error {
	log.Println("Launching sever...")
	// listen on all interfaces
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	log.Printf("Server launched on port:%d, waiting for clients...\n", port)

	i := 0
	for {
		// accept connection on port
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Connection error: %v", err)
		}
		i++
		log.Printf("Client connected: %d", i)
		go handleConnection(conn, i)
	}
}
