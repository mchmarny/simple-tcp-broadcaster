package commons

import (
	"log"
	"sync"
	"time"
)

// ConnectionManager manages server slide client connection
type ConnectionManager struct {
	Port       int
	Clients    map[*Connection]bool
	Broadcast  chan *SimpleMessage
	Register   chan *Connection
	Unregister chan *Connection
	Mutex      *sync.Mutex
	Stopping   bool
}

func (m *ConnectionManager) deleteConn(conn *Connection) {
	defer m.Mutex.Unlock()
	m.Mutex.Lock()
	close(conn.Message)
	delete(m.Clients, conn)
}

func (m *ConnectionManager) addConn(conn *Connection) {
	if !m.Stopping {
		defer m.Mutex.Unlock()
		m.Mutex.Lock()
		m.Clients[conn] = true
	}
}

// Stop cleans up all connections
func (m *ConnectionManager) Stop() {
	defer m.Mutex.Unlock()
	m.Mutex.Lock()
	m.Stopping = true
	for c := range m.Clients {
		log.Printf("Disconnecting: %s", c.GetRemoteConnectorID())
		c.Socket.SetDeadline(time.Now().Add(time.Millisecond * 10))
		c.Socket.Close()
		c.Socket = nil
	}
}

// Start distributes coommands
func (m *ConnectionManager) Start() {
	m.Stopping = false
	for {
		select {
		case c := <-m.Register:
			m.addConn(c)
			log.Printf("New connection: %s", c.GetRemoteConnectorID())
		case c := <-m.Unregister:
			if _, ok := m.Clients[c]; ok {
				connID := c.GetRemoteConnectorID()
				m.deleteConn(c)
				log.Printf("Connection terminated: %s", connID)
			}
		case msg := <-m.Broadcast:
			log.Println("Broadcasting...")
			for conn := range m.Clients {
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
func (m *ConnectionManager) Receive(c *Connection) {
	for {
		msg := &SimpleMessage{}
		err := c.Decoder.Decode(msg)
		if err != nil {
			m.Unregister <- c
			c.Socket.Close()
			break
		}
		if msg.Type == DataMessageTypeCode {
			log.Printf("Client %s message: %+v", c.GetRemoteConnectorID(), msg)
			m.Broadcast <- msg
		}
	}
}

// Send sends data back to the client
func (m *ConnectionManager) Send(c *Connection) {
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
