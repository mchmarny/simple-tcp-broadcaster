package commons

import "testing"

func TestMessage(t *testing.T) {

	sourceID := "12345"
	msg := NewMessage(sourceID, []byte("test"))

	if msg.Source != sourceID {
		t.Error("Failed to create expected source")
	}
}
