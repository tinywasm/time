//go:build wasm

package time_test

import (
	"testing"
)

// TestTimeAPIWasm tests the WASM/JS time API implementation.
func TestTimeAPIWasm(t *testing.T) {
	RunAPITests(t)
}
