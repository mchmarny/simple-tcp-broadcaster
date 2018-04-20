package server

import (
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/mchmarny/simple-server/commons"
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
	log.Println("Starting server...")
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	manager = &ClientManager{
		clients:    make(map[*commons.Connection]bool),
		broadcast:  make(chan *commons.SimpleMessage),
		register:   make(chan *commons.Connection),
		unregister: make(chan *commons.Connection),
		mutex:      &sync.Mutex{},
	}

	go manager.Start()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Client connect error: %v", err)
			continue
		}

		c := commons.NewSeverConnection(conn)
		manager.register <- c
		go manager.Receive(c)
		go manager.Send(c)
	}
}

// func (srv *Server) deleteConn(conn *conn) {
// 	defer srv.mu.Unlock()
// 	srv.mu.Lock()
// 	delete(srv.conns, conn)
// }

// func (srv *Server) Shutdown() {
// 	// should be guarded by mu
// 	srv.inShutdown = true
// 	log.Println("shutting down...")
// 	srv.listener.Close()
// 	ticker := time.NewTicker(500 * time.Millisecond)
// 	defer ticker.Stop()
// 	for {
// 		select {
// 		case <-ticker.C:
// 			log.Printf("waiting on %v connections", len(srv.conns))
// 		}
// 		if len(srv.conns) == 0 {
// 			return
// 		}
// 	}
// }
