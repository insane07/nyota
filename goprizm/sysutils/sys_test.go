package sysutils

import (
	"testing"
)

func TestNewUUID(t *testing.T) {
	uuid, err := NewUUID()
	if err != nil {
		t.Fatalf("failed to generate UUID err: %v", err)
	}

	if uuid == "" {
		t.Fatalf("expected non empty UUID")
	}

	if len(uuid) != 36 {
		t.Fatalf("invalid len %d for UUID(%s)", len(uuid), uuid)
	}
}
