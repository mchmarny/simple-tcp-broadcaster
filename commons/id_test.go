package commons

import "testing"

func TestParseID(t *testing.T) {
	rawLocalAddress := "127.0.0.1:45645"
	prefix := "client"
	expectedID := "client-127-0-0-1-45645"

	id := ParseID(prefix, rawLocalAddress)

	if id != expectedID {
		t.Error("Failed to parse expected ID")
	}
}
