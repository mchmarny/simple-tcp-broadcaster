package commons

import (
	"fmt"
	"log"
	"strings"

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

// ParseID combines parts into id
func ParseID(prefix, val string) string {
	val = strings.Replace(val, ".", "-", -1)
	val = strings.Replace(val, ":", "-", -1)
	return fmt.Sprintf("%s-%s", prefix, val)
}
