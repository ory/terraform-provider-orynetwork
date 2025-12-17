package provider

import (
	"testing"
)

func TestProvider(t *testing.T) {
	// Verify the provider can be instantiated
	p := New("test")()
	if p == nil {
		t.Fatal("provider should not be nil")
	}
}
