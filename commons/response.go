package commons

import (
	"time"
)

const (
	// NotSetResponseCode is default
	NotSetResponseCode ResponseCode = 0
	// SuccessResponseCode when success
	SuccessResponseCode ResponseCode = 1
	// ErrorResponseCode when error
	ErrorResponseCode ResponseCode = 2
)

// ResponseCode represents server response code
type ResponseCode int

// NewResponse creates a new response for specific request ID
func NewResponse(requestID string) *SimpleResponse {
	return &SimpleResponse{
		ID:        GetUUIDv4(),
		RequestID: requestID,
		CreatedAt: time.Now(),
		Status:    NotSetResponseCode,
	}
}

// SimpleResponse represents server response to simple request
type SimpleResponse struct {
	ID        string
	RequestID string
	CreatedAt time.Time
	Status    ResponseCode
	Data      []byte
}
