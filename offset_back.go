//go:build !wasm

package time

import "time"

// detectOffset auto-detects the timezone offset from the system.
func detectOffset() {
	_, offset := time.Now().Zone()
	tzOffsetMinutes.Store(int32(offset / 60))
}
