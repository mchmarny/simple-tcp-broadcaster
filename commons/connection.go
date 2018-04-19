package commons

import (
	"encoding/gob"
	"log"
	"net"
)

// Connection represents common client
type Connection struct {
	Socket  net.Conn
	Message chan *SimpleMessage
	//IdleTimeout   time.Duration
	//MaxReadBuffer int64
	Encoder *gob.Encoder
	Decoder *gob.Decoder
}

//NewConnection sets all the necessary defaults
func NewConnection(conn net.Conn) *Connection {
	return &Connection{
		Socket:  conn,
		Message: make(chan *SimpleMessage),
		Encoder: gob.NewEncoder(conn),
		Decoder: gob.NewDecoder(conn),
	}
}

func (c *Connection) Write(msg *SimpleMessage) {
	if msg != nil {
		c.Encoder.Encode(msg)
	}
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
