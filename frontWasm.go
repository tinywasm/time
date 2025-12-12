//go:build wasm

package time

import (
	"syscall/js"

	. "github.com/tinywasm/fmt"
)

// timeClient implements TimeProvider for WASM/JS environments using the JavaScript Date API.
type timeClient struct {
	dateCtor js.Value
}

// NewTimeProvider returns the correct implementation for WASM.
func NewTimeProvider() TimeProvider {
	return &timeClient{
		dateCtor: js.Global().Get("Date"),
	}
}

func (tc *timeClient) UnixNano() int64 {
	jsDate := tc.dateCtor.New()
	msTimestamp := jsDate.Call("getTime").Float()
	// Convert milliseconds to nanoseconds
	return int64(msTimestamp) * 1000000
}

func (tc *timeClient) FormatDate(value any) string {
	switch v := value.(type) {
	case int64:
		jsDate := tc.dateCtor.New(float64(v) / 1e6)
		return jsDate.Call("toISOString").String()[0:10]
	case string:
		// Validate date format: YYYY-MM-DD (10 chars with dashes at positions 4 and 7)
		if len(v) == 10 && v[4] == '-' && v[7] == '-' {
			return v
		}
	}
	return ""
}

func (tc *timeClient) FormatTime(value any) string {
	switch v := value.(type) {
	case int64: // UnixNano
		jsDate := tc.dateCtor.New(float64(v) / 1e6)
		hours := jsDate.Call("getUTCHours").Int()
		minutes := jsDate.Call("getUTCMinutes").Int()
		seconds := jsDate.Call("getUTCSeconds").Int()
		return Fmt("%02d:%02d:%02d", hours, minutes, seconds)
	case int16: // Minutes since midnight
		hours := v / 60
		minutes := v % 60
		return Fmt("%02d:%02d", hours, minutes)
	case string:
		if Count(v, ":") >= 1 {
			return v
		}
	}
	return ""
}

func (tc *timeClient) FormatDateTime(value any) string {
	switch v := value.(type) {
	case int64:
		jsDate := tc.dateCtor.New(float64(v) / 1e6)
		iso := jsDate.Call("toISOString").String()
		return iso[0:10] + " " + iso[11:19]
	case string:
		// Validate datetime format: YYYY-MM-DD HH:MM:SS (19 chars)
		if len(v) == 19 && v[4] == '-' && v[7] == '-' && v[10] == ' ' && v[13] == ':' && v[16] == ':' {
			return v
		}
	}
	return ""
}

func (tc *timeClient) FormatDateTimeShort(value any) string {
	switch v := value.(type) {
	case int64:
		jsDate := tc.dateCtor.New(float64(v) / 1e6)
		iso := jsDate.Call("toISOString").String()
		return iso[0:10] + " " + iso[11:16]
	case string:
		// Validate short datetime format: YYYY-MM-DD HH:MM (16 chars)
		if len(v) == 16 && v[4] == '-' && v[7] == '-' && v[10] == ' ' && v[13] == ':' {
			return v
		}
	}
	return ""
}

func (tc *timeClient) ParseDate(dateStr string) (int64, error) {
	// Validate format: YYYY-MM-DD (10 chars with dashes at positions 4 and 7)
	if len(dateStr) != 10 || dateStr[4] != '-' || dateStr[7] != '-' {
		return 0, Errf("invalid date format: %s (expected YYYY-MM-DD)", dateStr)
	}

	jsDate := tc.dateCtor.New(dateStr + "T00:00:00Z")
	if jsDate.Call("toString").String() == "Invalid Date" {
		return 0, Errf("invalid date format: %s", dateStr)
	}

	// Verify date components match (JS Date auto-corrects invalid dates like Feb 30)
	year := jsDate.Call("getUTCFullYear").Int()
	month := jsDate.Call("getUTCMonth").Int() + 1
	day := jsDate.Call("getUTCDate").Int()
	expected := Fmt("%04d-%02d-%02d", year, month, day)
	if expected != dateStr {
		return 0, Errf("invalid date: %s (auto-corrected to %s)", dateStr, expected)
	}

	ms := jsDate.Call("getTime").Float()
	return int64(ms) * 1000000, nil
}

func (tc *timeClient) ParseTime(timeStr string) (int16, error) {
	return parseTime(timeStr)
}

func (tc *timeClient) ParseDateTime(dateStr, timeStr string) (int64, error) {
	if len(timeStr) == 5 {
		timeStr += ":00"
	}
	isoStr := dateStr + "T" + timeStr + "Z"
	jsDate := tc.dateCtor.New(isoStr)
	if jsDate.Call("toString").String() == "Invalid Date" {
		return 0, Errf("invalid date/time format: %s %s", dateStr, timeStr)
	}
	ms := jsDate.Call("getTime").Float()
	return int64(ms) * 1000000, nil
}

func (tc *timeClient) IsToday(nano int64) bool {
	jsDate := tc.dateCtor.New(float64(nano) / 1e6)
	now := tc.dateCtor.New()
	return jsDate.Call("toDateString").String() == now.Call("toDateString").String()
}

func (tc *timeClient) IsPast(nano int64) bool {
	return nano < tc.UnixNano()
}

func (tc *timeClient) IsFuture(nano int64) bool {
	return nano > tc.UnixNano()
}

func (tc *timeClient) DaysBetween(nano1, nano2 int64) int {
	return daysBetween(nano1, nano2)
}

// WasmTimer implements Timer for WASM using setTimeout
type WasmTimer struct {
	id     int
	active bool
	jsFunc js.Func // Store to release later
	f      func()  // Store callback to execute
}

func (wt *WasmTimer) Stop() bool {
	if !wt.active {
		return false
	}
	js.Global().Call("clearTimeout", wt.id)
	wt.active = false
	wt.jsFunc.Release() // Free memory
	return true
}

func (wt *WasmTimer) Fire() {
	if !wt.active {
		return
	}
	wt.active = false
	wt.jsFunc.Release() // Free memory after execution
	if wt.f != nil {
		wt.f()
	}
}

func (tc *timeClient) AfterFunc(milliseconds int, f func()) Timer {
	wt := &WasmTimer{
		active: true,
		f:      f,
	}

	wt.jsFunc = js.FuncOf(func(this js.Value, args []js.Value) any {
		wt.Fire()
		return nil
	})

	wt.id = js.Global().Call("setTimeout", wt.jsFunc, milliseconds).Int()
	return wt
}
