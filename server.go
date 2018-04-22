package main

import (
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/mchmarny/simple-tcp-broadcaster/commons"
)

var (
	manager *commons.ConnectionManager
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

	manager = &commons.ConnectionManager{
		Port:       port,
		Clients:    make(map[*commons.Connection]bool),
		Broadcast:  make(chan *commons.SimpleMessage),
		Register:   make(chan *commons.Connection),
		Unregister: make(chan *commons.Connection),
		Mutex:      &sync.Mutex{},
	}

	go manager.Start()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Connect error: %v", err)
			continue
		}

		c := commons.NewSeverConnection(conn)
		manager.Register <- c
		go manager.Receive(c)
		go manager.Send(c)
	}
}
