//go:build !wasm

package time

import (
	"time"

	. "github.com/tinywasm/fmt"
)

// NewTimeProvider returns the correct implementation based on the build environment.
func NewTimeProvider() TimeProvider {
	return &timeServer{}
}

// timeServer implements TimeProvider for standard Go.
type timeServer struct{}

func (ts *timeServer) UnixNano() int64 {
	return time.Now().UTC().UnixNano()
}

func (ts *timeServer) FormatDate(value any) string {
	switch v := value.(type) {
	case int64:
		return time.Unix(0, v).UTC().Format("2006-01-02")
	case string:
		if _, err := time.Parse("2006-01-02", v); err == nil {
			return v
		}
	}
	return ""
}

func (ts *timeServer) FormatTime(value any) string {
	switch v := value.(type) {
	case int64: // UnixNano
		return time.Unix(0, v).UTC().Format("15:04:05")
	case int16: // Minutes since midnight
		hours := v / 60
		minutes := v % 60
		return Fmt("%02d:%02d", hours, minutes)
	case string:
		// Try to parse as numeric timestamp first (e.g., from unixid.GetNewID())
		if nano, err := Convert(v).Int64(); err == nil {
			return time.Unix(0, nano).UTC().Format("15:04:05")
		}
		// Otherwise check if already formatted
		if Count(v, ":") >= 1 {
			return v
		}
	}
	return ""
}

func (ts *timeServer) FormatDateTime(value any) string {
	switch v := value.(type) {
	case int64:
		return time.Unix(0, v).UTC().Format("2006-01-02 15:04:05")
	case string:
		if _, err := time.Parse("2006-01-02 15:04:05", v); err == nil {
			return v
		}
	}
	return ""
}

func (ts *timeServer) FormatDateTimeShort(value any) string {
	switch v := value.(type) {
	case int64:
		return time.Unix(0, v).UTC().Format("2006-01-02 15:04")
	case string:
		if _, err := time.Parse("2006-01-02 15:04", v); err == nil {
			return v
		}
	}
	return ""
}

func (ts *timeServer) ParseDate(dateStr string) (int64, error) {
	t, err := time.ParseInLocation("2006-01-02", dateStr, time.UTC)
	if err != nil {
		return 0, err
	}
	return t.UnixNano(), nil
}

func (ts *timeServer) ParseTime(timeStr string) (int16, error) {
	return parseTime(timeStr)
}

func (ts *timeServer) ParseDateTime(dateStr, timeStr string) (int64, error) {
	layout := "2006-01-02 15:04:05"
	if len(timeStr) == 5 {
		layout = "2006-01-02 15:04"
	}

	t, err := time.ParseInLocation(layout, dateStr+" "+timeStr, time.UTC)
	if err != nil {
		return 0, err
	}
	return t.UnixNano(), nil
}

func (ts *timeServer) IsToday(nano int64) bool {
	t := time.Unix(0, nano).UTC()
	now := time.Now().UTC()
	return t.Year() == now.Year() && t.YearDay() == now.YearDay()
}

func (ts *timeServer) IsPast(nano int64) bool {
	return nano < ts.UnixNano()
}

func (ts *timeServer) IsFuture(nano int64) bool {
	return nano > ts.UnixNano()
}

func (ts *timeServer) DaysBetween(nano1, nano2 int64) int {
	return daysBetween(nano1, nano2)
}

// timerWrapper wraps time.Timer to implement Timer interface
type timerWrapper struct {
	timer *time.Timer
}

func (tw *timerWrapper) Stop() bool {
	return tw.timer.Stop()
}

func (ts *timeServer) AfterFunc(milliseconds int, f func()) Timer {
	t := time.AfterFunc(time.Duration(milliseconds)*time.Millisecond, f)
	return &timerWrapper{timer: t}
}
