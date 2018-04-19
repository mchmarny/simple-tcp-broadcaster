package commons

import (
	"log"

	"github.com/google/uuid"
)

// GetUUIDv4 generates string formatted UUIDv4
func GetUUIDv4() string {
	id, err := uuid.NewRandom()
	if err != nil {
		log.Fatalf("Error while getting id: %v\n", err)
	}
	return id.String()
}
