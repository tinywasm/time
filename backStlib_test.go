//go:build !wasm

package time_test

import (
	"testing"

	"github.com/tinywasm/time"
)

// TestTimeProviderStdlib tests the standard Go time provider implementation.
func TestTimeProviderStdlib(t *testing.T) {
	tp := time.NewTimeProvider()
	RunProviderTests(t, tp)
}
