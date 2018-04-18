package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

var (
	// timeout for silent connecitons
	timeoutDuration = 25 * time.Minute
)

// StopServer gracefully shuts down server
func StopServer() {
	log.Println("Server stopped")
}

func handleConnection(c net.Conn) {

	defer func() {
		if err := c.Close(); err != nil {
			log.Printf("Error closing client connection: %v\n", err)
		}
	}()

	scanner := bufio.NewScanner(c)
	for scanner.Scan() {
		// Set timeout, err if received after deadline
		c.SetReadDeadline(time.Now().Add(timeoutDuration))

		msg := scanner.Text()
		msgNew := fmt.Sprintf("Hi there, I got your message: %s", msg)

		if _, err := c.Write([]byte(msgNew + "\n")); err != nil {
			log.Printf("Client closed writer: %v\n", err)
			break
		}

		log.Printf("Client\n received: %s\n returned: %s\n", msg, msgNew)

	}
	log.Printf("Client disconnected\n")
}

// StartServer start TCP server on a given port
func StartServer(port int) error {

	log.Println("Launching sever...")
	// listen on all interfaces
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	defer func() {
		if err := ln.Close(); err != nil {
			log.Printf("Error shutting down server: %v\n", err)
		}
	}()

	log.Printf("Server launched on port:%d, waiting for clients...\n", port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Client connection error: %v", err)
		}
		go handleConnection(conn)
	}

}
