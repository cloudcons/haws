package customtags

import "testing"

func TestNew(t *testing.T) {
	tags := New()
	if tags == nil {
		t.Fatal("New returned nil")
	}
}
