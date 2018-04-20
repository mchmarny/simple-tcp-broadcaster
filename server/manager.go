package server

import (
	"log"
	"sync"

	"github.com/mchmarny/simple-tcp-broadcaster/commons"
)

// ClientManager manages server slide client connection
type ClientManager struct {
	port       int
	clients    map[*commons.Agent]bool
	broadcast  chan *commons.SimpleMessage
	register   chan *commons.Agent
	unregister chan *commons.Agent
	mutex      *sync.Mutex
	stopping   bool
}

func (m *ClientManager) deleteConn(conn *commons.Agent) {
	defer m.mutex.Unlock()
	m.mutex.Lock()
	close(conn.Message)
	delete(m.clients, conn)
}

func (m *ClientManager) addConn(conn *commons.Agent) {
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
		log.Printf("Disconnecting: %s", c.GetRemoteClientID())
		c.Socket.Close()
		c.Socket = nil
	}
}

// Start distributes coommands
func (m *ClientManager) Start() {
	m.stopping = false
	for {
		select {
		case c := <-m.register:
			m.addConn(c)
			log.Printf("New connection: %s", c.GetRemoteClientID())
		case c := <-m.unregister:
			if _, ok := m.clients[c]; ok {
				connID := c.GetRemoteClientID()
				m.deleteConn(c)
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
func (m *ClientManager) Receive(c *commons.Agent) {
	for {
		msg := &commons.SimpleMessage{}
		err := c.Decoder.Decode(msg)
		if err != nil {
			m.unregister <- c
			c.Socket.Close()
			break
		}
		log.Printf("Client %s message: %+v", c.GetRemoteClientID(), msg)
		m.broadcast <- msg
	}
}

// Send sends data back to the client
func (m *ClientManager) Send(c *commons.Agent) {
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
