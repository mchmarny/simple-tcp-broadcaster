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
	broadcast  chan *commons.SimpleMessage
	register   chan *commons.Connection
	unregister chan *commons.Connection
}

// Start distributes coommands
func (m *ClientManager) Start() {
	for {
		select {
		case conn := <-m.register:
			m.clients[conn] = true
			log.Printf("New connection: %s", conn.Socket.RemoteAddr().String())
		case conn := <-m.unregister:
			if _, ok := m.clients[conn]; ok {
				connID := conn.Socket.RemoteAddr().String()
				close(conn.Message)
				delete(m.clients, conn)
				log.Printf("Connection terminated: %s", connID)
			}
		case msg := <-m.broadcast:
			log.Println("Broadcasting...")
			for conn := range m.clients {
				select {
				case conn.Message <- msg:
				default:
					close(conn.Message)
					delete(m.clients, conn)
				}
			}
		}
	}
}

// Receive processes manager connecitons
func (m *ClientManager) Receive(c *commons.Connection) {
	for {
		msg := &commons.SimpleMessage{}
		err := c.Decoder.Decode(msg)
		if err != nil {
			m.unregister <- c
			c.Socket.Close()
			break
		}
		log.Printf("Client message: %+v", msg)
		m.broadcast <- msg
	}
}

// Send sends data back to the client
func (m *ClientManager) Send(c *commons.Connection) {
	defer c.Socket.Close()
	for {
		select {
		case msg, ok := <-c.Message:
			if !ok {
				return
			}
			c.Encoder.Encode(msg)
		}
	}
}
