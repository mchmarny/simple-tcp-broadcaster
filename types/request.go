package types

import (
	"time"

	"github.com/mchmarny/simple-server/utils"
)

// NewRequest creates a new reqquest for this client ID
func NewRequest(clientID string) SimpleRequest {
	return SimpleRequest{
		ID:        utils.GetUUIDv4(),
		CreatedAt: time.Now(),
		Client:    clientID,
	}
}

// SimpleRequest represents simple client request
type SimpleRequest struct {
	ID        string
	CreatedAt time.Time
	Client    string
	Data      []byte
}
