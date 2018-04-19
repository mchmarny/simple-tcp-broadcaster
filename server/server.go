package server

import (
	"fmt"
	"log"
	"net"

	"github.com/mchmarny/simple-server/commons"
)

var (
	manager *ClientManager
)

// StartServerMode starts TCP server on specified port
func StartServerMode(port int) error {
	log.Println("Starting server...")
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	manager = &ClientManager{
		clients:    make(map[*commons.Connection]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *commons.Connection),
		unregister: make(chan *commons.Connection),
	}

	go manager.Start()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Client connect error: %v", err)
			continue
		}
		client := &commons.Connection{
			Socket: conn,
			Data:   make(chan []byte),
		}
		manager.register <- client
		go manager.Receive(client)
		go manager.Send(client)
	}
}

// 	sc := make(chan bool)
// 	deadline := time.After(conn.IdleTimeout)
// 	for {
// 		go func(s chan bool) {
// 			s <- scanr.Scan()
// 		}(sc)
// 		select {
// 		case <-deadline:
// 			return nil
// 		case scanned := <-sc:
// 			if !scanned {
// 				if err := scanr.Err(); err != nil {
// 					return err
// 				}
// 				return nil
// 			}
// 			val := scanr.Text()
// 			log.Printf("Client said: %s", val)
// 			w.WriteString(strings.ToUpper(val) + "\n")
// 			w.Flush()
// 			deadline = time.After(conn.IdleTimeout)
// 		}
// 	}
// }

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
