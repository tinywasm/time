//go:build !wasm

package time_test

import (
	"testing"
)

// TestTimeAPIBackend tests the standard Go time API implementation.
func TestTimeAPIBackend(t *testing.T) {
	RunAPITests(t)
}
