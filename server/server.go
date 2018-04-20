package server

import (
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/mchmarny/simple-tcp-broadcaster/commons"
)

var (
	manager *ClientManager
)

// StopServer cleans up the connected clietns
func StopServer() {
	manager.Stop()
}

// StartServer starts TCP server on specified port
func StartServer(port int) error {
	log.Printf("Starting server on port:%d ...", port)
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	manager = &ClientManager{
		port:       port,
		clients:    make(map[*commons.Agent]bool),
		broadcast:  make(chan *commons.SimpleMessage),
		register:   make(chan *commons.Agent),
		unregister: make(chan *commons.Agent),
		mutex:      &sync.Mutex{},
	}

	go manager.Start()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Connect error: %v", err)
			continue
		}

		c := commons.NewSeverAgent(conn)
		manager.register <- c
		go manager.Receive(c)
		go manager.Send(c)
	}
}
