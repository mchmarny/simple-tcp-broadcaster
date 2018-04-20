package commons

import (
	"encoding/gob"
	"log"
	"net"
	"time"
)

const (
	defaultClientTimeout = 60
	inspectionPeriod     = 10

	// UndefinedConnectionMode is default
	UndefinedConnectionMode ConnectionMode = 0
	// ServerConnectionMode client
	ServerConnectionMode ConnectionMode = 1
	// ClientConnectionMode server
	ClientConnectionMode ConnectionMode = 2
)

// ConnectionMode indicated client or server connection mode
type ConnectionMode int

// Connection represents common client
type Connection struct {
	Socket      net.Conn
	Message     chan *SimpleMessage
	IdleTimeout time.Duration
	Mode        ConnectionMode
	Encoder     *gob.Encoder
	Decoder     *gob.Decoder
}

// NewSeverConnection creates server connection
func NewSeverConnection(conn net.Conn) *Connection {
	return newConnection(conn, ServerConnectionMode)
}

// NewClientConnection creates client connection
func NewClientConnection(conn net.Conn) *Connection {
	return newConnection(conn, ClientConnectionMode)
}

//NewConnection sets all the necessary defaults
func newConnection(conn net.Conn, mode ConnectionMode) *Connection {
	c := &Connection{
		Socket:      conn,
		Mode:        mode,
		Message:     make(chan *SimpleMessage),
		IdleTimeout: time.Minute * time.Duration(defaultClientTimeout),
		Encoder:     gob.NewEncoder(conn),
		Decoder:     gob.NewDecoder(conn),
	}
	c.updateDeadline()

	if mode == ServerConnectionMode {
		go c.watch()
	} else if mode == ClientConnectionMode {
		c.Socket.(*net.TCPConn).SetKeepAlive(true)
		c.Socket.(*net.TCPConn).SetKeepAlivePeriod(time.Second * time.Duration(inspectionPeriod))
	}

	return c
}

func (c *Connection) watch() {
	inspectoin := time.After(time.Second * inspectionPeriod)
	for {
		select {
		case <-inspectoin:
			log.Printf("Inspected: %s", c.Socket.RemoteAddr().String())
			inspectoin = time.After(time.Second * inspectionPeriod)
		}
	}
}

// updateDeadline resets the connection timeout
func (c *Connection) updateDeadline() {
	idleDeadline := time.Now().Add(c.IdleTimeout)
	c.Socket.SetDeadline(idleDeadline)
}

// Disconnect sets imediate connection and closes it
func (c *Connection) Disconnect() {
	idleDeadline := time.Now().Add(time.Microsecond * 10)
	c.Socket.SetDeadline(idleDeadline)
	c.Socket.Close()
}

func (c *Connection) Write(msg *SimpleMessage) error {
	c.updateDeadline()
	if msg != nil {
		return c.Encoder.Encode(msg)
	}
	return nil
}

// Read processes client messages
func (c *Connection) Read() {
	for {
		msg := &SimpleMessage{}
		err := c.Decoder.Decode(msg)
		if err != nil {
			c.Socket.Close()
			break
		}
		log.Printf("Client message: %+v", msg)
	}
}
