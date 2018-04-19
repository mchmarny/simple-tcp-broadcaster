package commons

import (
	"time"
)

const (
	// UndefinedResponseCode is default
	UndefinedResponseCode ResponseCode = 0
	// SuccessResponseCode when success
	SuccessResponseCode ResponseCode = 1
	// ErrorResponseCode when error
	ErrorResponseCode ResponseCode = 2
)

// ResponseCode represents server response code
type ResponseCode int

// NewMessage creates a new response for specific request ID
func NewMessage(sourceID string) *SimpleMessage {
	return &SimpleMessage{
		ID:        GetUUIDv4(),
		Source:    sourceID,
		CreatedAt: time.Now(),
		Status:    UndefinedResponseCode,
	}
}

// SimpleMessage represents server response to simple request
type SimpleMessage struct {
	ID        string
	Source    string
	CreatedAt time.Time
	Status    ResponseCode
	Data      []byte
}
