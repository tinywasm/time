//go:build wasm
// +build wasm

package time

// FireTimer triggers the timer callback manually for testing purposes.
// This allows verifying the callback logic without waiting for the browser event loop.
func FireTimer(t Timer) {
	if wt, ok := t.(*wasmTimer); ok {
		wt.fire()
	}
}
