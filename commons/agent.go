package commons

import (
	"encoding/gob"
	"log"
	"net"
	"time"
)

const (
	defaultClientTimeout  = 60
	inspectionPeriod      = 10
	keepAliveEverySeconds = 10

	// UndefinedAgentMode is default
	UndefinedAgentMode AgentMode = 0
	// ServerAgentMode client
	ServerAgentMode AgentMode = 1
	// ClientAgentMode server
	ClientAgentMode AgentMode = 2

	// ClientPrefix is the bit before client id
	ClientPrefix = "client"
	//ServerPrefix is the bit before server id
	ServerPrefix = "server"
)

// AgentMode indicated client or server connection mode
type AgentMode int

// Agent represents common client
type Agent struct {
	Socket      net.Conn
	Message     chan *SimpleMessage
	IdleTimeout time.Duration
	Mode        AgentMode
	Encoder     *gob.Encoder
	Decoder     *gob.Decoder
}

// NewSeverAgent creates server connection
func NewSeverAgent(conn net.Conn) *Agent {
	return newAgent(conn, ServerAgentMode)
}

// NewClientAgent creates client connection
func NewClientAgent(serverAddress string) (agent *Agent, err error) {
	c, e := net.Dial("tcp", serverAddress)
	if e != nil {
		log.Fatalf("Error on dial: %v", e)
		return nil, err
	}

	log.Printf("Connected to server %s from %s", c.RemoteAddr(), c.LocalAddr())

	tcpConn := c.(*net.TCPConn)
	tcpConn.SetKeepAlive(true)
	tcpConn.SetKeepAlivePeriod(time.Second * time.Duration(keepAliveEverySeconds))

	return newAgent(c, ClientAgentMode), nil
}

//newAgent sets all the necessary defaults
func newAgent(conn net.Conn, mode AgentMode) *Agent {
	c := &Agent{
		Socket:      conn,
		Mode:        mode,
		Message:     make(chan *SimpleMessage),
		IdleTimeout: time.Minute * time.Duration(defaultClientTimeout),
		Encoder:     gob.NewEncoder(conn),
		Decoder:     gob.NewDecoder(conn),
	}
	c.updateDeadline()
	return c
}

// GetLocalClientID returns local address as string
func (c *Agent) GetLocalClientID() string {
	return ParseID(ClientPrefix, c.Socket.LocalAddr().String())
}

// GetLocalServerID returns local address as string
func (c *Agent) GetLocalServerID() string {
	return ParseID(ServerPrefix, c.Socket.LocalAddr().String())
}

// GetRemoteClientID returns remote address as string
func (c *Agent) GetRemoteClientID() string {
	return ParseID(ClientPrefix, c.Socket.RemoteAddr().String())
}

// GetRemoteServerID returns remote address as string
func (c *Agent) GetRemoteServerID() string {
	return ParseID(ServerPrefix, c.Socket.RemoteAddr().String())
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
func (c *Agent) updateDeadline() {
	idleDeadline := time.Now().Add(c.IdleTimeout)
	c.Socket.SetDeadline(idleDeadline)
}

// Stop closes the connection
func (c *Agent) Stop() error {
	nowish := time.Now().Add(time.Millisecond * 10)
	c.Socket.SetDeadline(nowish)
	return c.Socket.Close()
}

func (c *Agent) Write(msg *SimpleMessage) error {
	c.updateDeadline()
	if msg != nil {
		return c.Encoder.Encode(msg)
	}
	return nil
}

// Read processes client messages
func (c *Agent) Read() {
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
