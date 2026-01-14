//go:build wasm

package time

import "syscall/js"

// detectOffset auto-detects the timezone offset from the browser.
func detectOffset() {
	// JS Date.getTimezoneOffset() returns minutes between UTC and local time.
	// Note: JS returns positive for UTC- (e.g. +180 for UTC-3), so we invert it.
	jsDate := js.Global().Get("Date").New()
	offsetMinutes := jsDate.Call("getTimezoneOffset").Int()
	tzOffsetMinutes.Store(int32(-offsetMinutes))
}
