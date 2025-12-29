//go:build wasm

package time_test

import (
	"testing"

	"github.com/tinywasm/time"
)

// TestTimeProviderWasm tests the WASM/JS time provider implementation.
func TestTimeProviderWasm(t *testing.T) {
	tp := time.NewTimeProvider()
	RunProviderTests(t, tp)
}
