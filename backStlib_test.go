//go:build !wasm
// +build !wasm

package time_test

import (
	"testing"

	"github.com/tinywasm/time"
)

// Tests for standard Go environment using shared validation functions

func TestAllShared(t *testing.T) {
	tp := time.NewTimeProvider()

	t.Run("UnixNano", func(t *testing.T) { UnixNanoShared(t, tp) })
	t.Run("FormatDate", func(t *testing.T) { FormatDateShared(t, tp) })
	t.Run("FormatTime", func(t *testing.T) { FormatTimeShared(t, tp) })
	t.Run("FormatDateTime", func(t *testing.T) { FormatDateTimeShared(t, tp) })
	t.Run("FormatDateTimeShort", func(t *testing.T) { FormatDateTimeShortShared(t, tp) })
	t.Run("ParseDate", func(t *testing.T) { ParseDateShared(t, tp) })
	t.Run("ParseTime", func(t *testing.T) { ParseTimeShared(t, tp) })
	t.Run("ParseDateTime", func(t *testing.T) { ParseDateTimeShared(t, tp) })
	t.Run("IsToday", func(t *testing.T) { IsTodayShared(t, tp) })
	t.Run("IsPast", func(t *testing.T) { IsPastShared(t, tp) })
	t.Run("IsFuture", func(t *testing.T) { IsFutureShared(t, tp) })
	t.Run("DaysBetween", func(t *testing.T) { DaysBetweenShared(t, tp) })
}
