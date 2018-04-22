package commons

import (
	"time"
)

const (
	undefinedMessageTypeText = "Undefined"
	// UndefinedMessageTypeCode is default
	UndefinedMessageTypeCode MessageTypeCode = 0
	// DataMessageTypeCode when data
	DataMessageTypeCode MessageTypeCode = 1
	// HeartbeatMessageTypeCode when heartbeat
	HeartbeatMessageTypeCode MessageTypeCode = 2
)

// MessageTypeCode indicates what type of message this is
type MessageTypeCode int

// String returns string representation of the enum
func (t MessageTypeCode) String() string {

	names := [...]string{
		undefinedMessageTypeText,
		"Data",
		"Heartbeat",
	}

	if t < DataMessageTypeCode || t > HeartbeatMessageTypeCode {
		return undefinedMessageTypeText
	}

	return names[t]
}

// NewHeartbeatMessage creates a heartbeat message
func NewHeartbeatMessage(clientID string) *SimpleMessage {
	return &SimpleMessage{
		ID:        GetUUIDv4(),
		Source:    clientID,
		CreatedAt: time.Now(),
		Type:      HeartbeatMessageTypeCode,
	}
}

// NewMessage creates a new response for specific request ID
func NewMessage(clientID string, data []byte) *SimpleMessage {
	return &SimpleMessage{
		ID:        GetUUIDv4(),
		Source:    clientID,
		CreatedAt: time.Now(),
		Type:      DataMessageTypeCode,
		Data:      data,
	}
}

// SimpleMessage represents server response to simple request
type SimpleMessage struct {
	ID        string
	Source    string
	CreatedAt time.Time
	Type      MessageTypeCode
	Data      []byte
}

// GetDataString returns message data as string
func (m *SimpleMessage) GetDataString() string {
	if m.Data != nil {
		return string(m.Data)
	}
	return ""
}
