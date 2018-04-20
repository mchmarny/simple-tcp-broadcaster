package server

import (
	"log"
	"sync"

	"github.com/mchmarny/simple-server/commons"
)

// ClientManager manages server slide client connection
type ClientManager struct {
	clients    map[*commons.Connection]bool
	broadcast  chan *commons.SimpleMessage
	register   chan *commons.Connection
	unregister chan *commons.Connection
	mutex      *sync.Mutex
	stopping   bool
}

func (m *ClientManager) deleteConn(conn *commons.Connection) {
	defer m.mutex.Unlock()
	m.mutex.Lock()
	close(conn.Message)
	delete(m.clients, conn)
}

func (m *ClientManager) addConn(conn *commons.Connection) {
	if !m.stopping {
		defer m.mutex.Unlock()
		m.mutex.Lock()
		m.clients[conn] = true
	}
}

// Stop cleans up all connections
func (m *ClientManager) Stop() {
	defer m.mutex.Unlock()
	m.mutex.Lock()
	m.stopping = true
	for c := range m.clients {
		log.Printf("Disconnecting: %s", c.Socket.RemoteAddr().String())
		c.Socket.Close()
		c.Socket = nil
	}
}

// Start distributes coommands
func (m *ClientManager) Start() {
	m.stopping = false
	for {
		select {
		case conn := <-m.register:
			m.addConn(conn)
			log.Printf("New connection: %s", conn.Socket.RemoteAddr().String())
		case conn := <-m.unregister:
			if _, ok := m.clients[conn]; ok {
				connID := conn.Socket.RemoteAddr().String()
				m.deleteConn(conn)
				log.Printf("Connection terminated: %s", connID)
			}
		case msg := <-m.broadcast:
			log.Println("Broadcasting...")
			for conn := range m.clients {
				select {
				case conn.Message <- msg:
				default:
					m.deleteConn(conn)
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
