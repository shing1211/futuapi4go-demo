package test

import (
	"testing"

	"github.com/shing1211/futuapi4go/client"
)

// TestSmokeWiring constructs an SDK client and verifies the wiring is intact.
// client.New() makes no network call. The test catches a dependency upgrade
// that silently breaks the SDK constructor or the transport layer.
func TestSmokeWiring(t *testing.T) {
	cli := client.New()
	if cli == nil {
		t.Fatal("client.New() returned a nil client")
	}
	if cli.IsConnected() {
		t.Fatal("a freshly constructed client must not be connected")
	}
}
