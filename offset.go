package time

import "sync/atomic"

// tzOffsetMinutes stores the timezone offset in minutes from UTC.
// Positive values mean ahead of UTC (e.g., +60 for UTC+1).
// Negative values mean behind UTC (e.g., -180 for UTC-3).
var tzOffsetMinutes atomic.Int32

// SetTimeZoneOffset manually sets the timezone offset in hours.
func SetTimeZoneOffset(hours int) {
	tzOffsetMinutes.Store(int32(hours * 60))
}

// GetTimeZoneOffset returns the current timezone offset in hours.
func GetTimeZoneOffset() int {
	return int(tzOffsetMinutes.Load() / 60)
}

// getOffsetMinutes returns the internal offset in minutes.
func getOffsetMinutes() int32 {
	return tzOffsetMinutes.Load()
}

func init() {
	detectOffset()
}
