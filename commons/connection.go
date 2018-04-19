package commons

import (
	"log"
	"net"
	"strings"
	"time"
)

// Connection represents common client
type Connection struct {
	Socket        net.Conn
	Data          chan []byte
	IdleTimeout   time.Duration
	MaxReadBuffer int64
}

// GetID returns connectionID
func (c *Connection) GetID() string {
	if c.Socket != nil {
		return c.Socket.RemoteAddr().String()
	} else {
		return ""
	}
}

func (c *Connection) Write(msg string) {
	if msg != "" {
		c.Socket.Write([]byte(strings.TrimRight(msg, "\n")))
	}
}

// Read processes client messages
func (c *Connection) Read() {
	for {
		msg := make([]byte, 4096)
		length, err := c.Socket.Read(msg)
		if err != nil {
			c.Socket.Close()
			break
		}
		if length > 0 {
			log.Printf("Client message: %s\n", msg)
		}
	}
}
