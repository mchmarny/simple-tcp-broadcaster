package commons

import (
	"time"
)

// NewRequest creates a new reqquest for this client ID
func NewRequest(clientID string) SimpleRequest {
	return SimpleRequest{
		ID:        GetUUIDv4(),
		CreatedAt: time.Now(),
		ClientID:  clientID,
	}
}

// SimpleRequest represents simple client request
type SimpleRequest struct {
	ID        string
	CreatedAt time.Time
	ClientID  string
	Data      []byte
}
