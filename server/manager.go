package server

import (
	"log"

	"github.com/mchmarny/simple-server/commons"
)

const (
	maxBuffer = 4096
)

// ClientManager manages server slide client connection
type ClientManager struct {
	clients    map[*commons.Connection]bool
	broadcast  chan []byte
	register   chan *commons.Connection
	unregister chan *commons.Connection
}

// Start distributes coommands
func (m *ClientManager) Start() {
	for {
		select {
		case conn := <-m.register:
			m.clients[conn] = true
			log.Printf("New connection: %s", conn.GetID())
		case conn := <-m.unregister:
			if _, ok := m.clients[conn]; ok {
				connID := conn.GetID()
				close(conn.Data)
				delete(m.clients, conn)
				log.Printf("Connection terminated: %s", connID)
			}
		case msg := <-m.broadcast:
			log.Println("Broadcasting...")
			for conn := range m.clients {
				select {
				case conn.Data <- msg:
				default:
					close(conn.Data)
					delete(m.clients, conn)
				}
			}
		}
	}
}

// Receive processes manager connecitons
func (m *ClientManager) Receive(c *commons.Connection) {
	for {
		msg := make([]byte, maxBuffer)
		length, err := c.Socket.Read(msg)
		if err != nil {
			m.unregister <- c
			c.Socket.Close()
			break
		}
		if length > 0 {
			log.Println("Client message: " + string(msg))
			m.broadcast <- msg
		}
	}
}

// Send sends data back to the client
func (m *ClientManager) Send(c *commons.Connection) {
	defer c.Socket.Close()
	for {
		select {
		case msg, ok := <-c.Data:
			if !ok {
				return
			}
			c.Socket.Write(msg)
		}
	}
}
