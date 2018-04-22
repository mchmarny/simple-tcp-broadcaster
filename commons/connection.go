package commons

import (
	"encoding/gob"
	"log"
	"net"
	"time"
)

const (
	defaultConnectionTimeout   = 60
	connectionInspectionPeriod = 10

	undefinedConnectionModeText = "Undefined"
	// UndefinedConnectionMode is default
	UndefinedConnectionMode ConnectionMode = 0
	// ServerConnectionMode client
	ServerConnectionMode ConnectionMode = 1
	// ClientConnectionMode server
	ClientConnectionMode ConnectionMode = 2

	// ClientConnectionModePrefix is the bit before client id
	ClientConnectionModePrefix = "client"
	//ServerConnectionModePrefix is the bit before server id
	ServerConnectionModePrefix = "server"
)

// ConnectionMode indicated client or server connection mode
type ConnectionMode int

// String returns string representation of the enum
func (t ConnectionMode) String() string {

	names := [...]string{
		undefinedConnectionModeText,
		"Server",
		"Client",
	}

	if t < ServerConnectionMode || t > ClientConnectionMode {
		return undefinedConnectionModeText
	}

	return names[t]
}

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

//newConnection sets all the necessary defaults
func newConnection(conn net.Conn, mode ConnectionMode) *Connection {
	c := &Connection{
		Socket:      conn,
		Mode:        mode,
		Message:     make(chan *SimpleMessage),
		IdleTimeout: time.Minute * time.Duration(defaultConnectionTimeout),
		Encoder:     gob.NewEncoder(conn),
		Decoder:     gob.NewDecoder(conn),
	}
	c.updateDeadline()
	return c
}

// GetLocalClientID returns local address as string
func (c *Connection) GetLocalClientID() string {
	return ParseID(ClientConnectionModePrefix, c.Socket.LocalAddr().String())
}

// GetLocalServerID returns local address as string
func (c *Connection) GetLocalServerID() string {
	return ParseID(ServerConnectionModePrefix, c.Socket.LocalAddr().String())
}

// GetRemoteConnectorID returns remote address as string
func (c *Connection) GetRemoteConnectorID() string {
	return ParseID(ClientConnectionModePrefix, c.Socket.RemoteAddr().String())
}

// GetRemoteServerID returns remote address as string
func (c *Connection) GetRemoteServerID() string {
	return ParseID(ServerConnectionModePrefix, c.Socket.RemoteAddr().String())
}

// func (c *Connection) watch() {
// 	inspectoin := time.After(time.Second * inspectionPeriod)
// 	for {
// 		select {
// 		case <-inspectoin:
// 			log.Printf("Inspected: %s", c.Socket.RemoteAddr().String())
// 			inspectoin = time.After(time.Second * inspectionPeriod)
// 		}
// 	}
// }

// updateDeadline resets the connection timeout
func (c *Connection) updateDeadline() {
	idleDeadline := time.Now().Add(c.IdleTimeout)
	c.Socket.SetDeadline(idleDeadline)
}

// Stop closes the connection
func (c *Connection) Stop() error {
	nowish := time.Now().Add(time.Millisecond * 10)
	c.Socket.SetDeadline(nowish)
	return c.Socket.Close()
}

func (c *Connection) Write(msg *SimpleMessage) error {
	c.updateDeadline()
	if msg != nil {
		return c.Encoder.Encode(msg)
	}
	return nil
}

// WriteHeartbeat sends heartbeat no data
func (c *Connection) WriteHeartbeat() error {
	c.updateDeadline()
	return c.Encoder.Encode(
		NewHeartbeatMessage(c.GetLocalClientID()),
	)
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
		log.Printf("Connection message: %+v", msg)
	}
}
