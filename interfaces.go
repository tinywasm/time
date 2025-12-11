package time

// TimeProvider defines the interface for time utilities, implemented for both standard Go and WASM/JS environments.
type TimeProvider interface {
	// UnixNano retrieves the current Unix timestamp in nanoseconds.
	// e.g., 1624397134562544800
	UnixNano() int64

	// FormatDate formats a value into a date string: "YYYY-MM-DD".
	// Accepts: int64 (UnixNano), string ("2024-01-15").
	FormatDate(value any) string

	// FormatTime formats a value into a time string.
	// Accepts: int64 (UnixNano) -> "HH:MM:SS", int16 (minutes) -> "HH:MM", string ("08:30").
	FormatTime(value any) string

	// FormatDateTime formats a value into a date-time string: "YYYY-MM-DD HH:MM:SS".
	// Accepts: int64 (UnixNano), string ("2024-01-15 08:30:45").
	FormatDateTime(value any) string

	// FormatDateTimeShort formats a value into a short date-time string: "YYYY-MM-DD HH:MM" (without seconds).
	// Accepts: int64 (UnixNano), string ("2024-01-15 08:30").
	FormatDateTimeShort(value any) string

	// ParseDate parses a date string ("YYYY-MM-DD") into a UnixNano timestamp (at midnight UTC).
	ParseDate(dateStr string) (int64, error)

	// ParseTime parses a time string ("HH:MM" or "HH:MM:SS") into minutes since midnight.
	ParseTime(timeStr string) (int16, error)

	// ParseDateTime combines date and time strings into a single UnixNano timestamp (UTC).
	ParseDateTime(dateStr, timeStr string) (int64, error)

	// IsToday checks if the given UnixNano timestamp is today (UTC).
	IsToday(nano int64) bool

	// IsPast checks if the given UnixNano timestamp is in the past.
	IsPast(nano int64) bool

	// IsFuture checks if the given UnixNano timestamp is in the future.
	IsFuture(nano int64) bool

	// DaysBetween calculates the number of full days between two UnixNano timestamps.
	DaysBetween(nano1, nano2 int64) int

	// AfterFunc waits for the specified milliseconds then calls f.
	// Returns a Timer that can be used to cancel the call.
	// WARNING: In WASM, callback runs in JS event loop - keep it lightweight.
	AfterFunc(milliseconds int, f func()) Timer
}

// Timer represents a cancelable timer
type Timer interface {
	// Stop prevents the timer from firing. Returns true if the timer was active.
	Stop() bool
}
