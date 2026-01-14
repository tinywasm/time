package time

// Timer represents a cancelable timer.
type Timer interface {
	// Stop prevents the timer from firing. Returns true if the timer was active.
	Stop() bool
}

// Now retrieves the current Unix timestamp in nanoseconds in UTC.
func Now() int64 {
	return provider.UnixNano()
}

// FormatTime formats a value into a time string "HH:MM:SS" applying the timezone offset.
func FormatTime(value any) string {
	return provider.FormatTime(value)
}

// FormatDate formats a value into a date string "YYYY-MM-DD" applying the timezone offset.
func FormatDate(value any) string {
	return provider.FormatDate(value)
}

// FormatDateTime formats a value into a date-time string "YYYY-MM-DD HH:MM:SS" applying the timezone offset.
func FormatDateTime(value any) string {
	return provider.FormatDateTime(value)
}

// FormatDateTimeShort formats a value into a short date-time string "YYYY-MM-DD HH:MM".
func FormatDateTimeShort(value any) string {
	return provider.FormatDateTimeShort(value)
}

// ParseDate parses a date string ("YYYY-MM-DD") into a UnixNano timestamp (UTC).
func ParseDate(dateStr string) (int64, error) {
	return provider.ParseDate(dateStr)
}

// ParseTime parses a time string into minutes since midnight (UTC).
func ParseTime(timeStr string) (int16, error) {
	return provider.ParseTime(timeStr)
}

// ParseDateTime combines date and time strings into a single UnixNano timestamp (UTC).
func ParseDateTime(dateStr, timeStr string) (int64, error) {
	return provider.ParseDateTime(dateStr, timeStr)
}

// IsToday checks if the given UnixNano timestamp is today according to the current timezone offset.
func IsToday(nano int64) bool {
	return provider.IsToday(nano)
}

// IsPast checks if the given UnixNano timestamp is in the past.
func IsPast(nano int64) bool {
	return provider.IsPast(nano)
}

// IsFuture checks if the given UnixNano timestamp is in the future.
func IsFuture(nano int64) bool {
	return provider.IsFuture(nano)
}

// DaysBetween calculates the number of full days between two UnixNano timestamps.
func DaysBetween(nano1, nano2 int64) int {
	return DaysBetweenShared(nano1, nano2)
}

// AfterFunc waits for the specified milliseconds then calls f.
func AfterFunc(milliseconds int, f func()) Timer {
	return provider.AfterFunc(milliseconds, f)
}

// Internal interface for the singleton provider
type timeProvider interface {
	UnixNano() int64
	FormatDate(value any) string
	FormatTime(value any) string
	FormatDateTime(value any) string
	FormatDateTimeShort(value any) string
	ParseDate(dateStr string) (int64, error)
	ParseTime(timeStr string) (int16, error)
	ParseDateTime(dateStr, timeStr string) (int64, error)
	IsToday(nano int64) bool
	IsPast(nano int64) bool
	IsFuture(nano int64) bool
	AfterFunc(milliseconds int, f func()) Timer
}

var provider timeProvider
