# PLAN: Add Weekday, MidnightUTC, LocalMinutesToUnixUTC to tinywasm/time

## Motivation

The `appointment-booking` module (velty_modules) requires three functions from
`github.com/tinywasm/time` that do not exist in v0.4.0. A previous agent worked
around the gap by forking the entire package into `internal/tinytime` and
redirecting imports with a `go.mod replace` directive — this must be undone.
The correct fix is to add the three functions here, publish a new version, and
let the consumer drop its local fork.

## Missing functions

| Function | Where computation lives | Provider needed? |
|---|---|---|
| `Weekday(unixSec int64) int` | pure arithmetic | No — add to `api.go` |
| `MidnightUTC(unixSec int64) int64` | pure arithmetic | No — add to `api.go` |
| `LocalMinutesToUnixUTC(dateSec int64, localMinutes int, tz string) int64` | IANA lookup (back) / global offset (front) | Yes — add to interface + both impls |

---

## Execution plan

### Step 1 — `api.go`: add pure functions and public wrapper

Add after `DaysBetween` and before the `timeProvider` interface:

```go
// Weekday returns the day of the week (0=Sunday … 6=Saturday) for a Unix
// timestamp in seconds (UTC). Based on the fact that 1970-01-01 was a Thursday (4).
func Weekday(unixSec int64) int {
	return int((unixSec/86400 + 4) % 7)
}

// MidnightUTC returns the Unix timestamp in seconds for midnight UTC of the
// day that contains the given Unix timestamp in seconds.
func MidnightUTC(unixSec int64) int64 {
	return (unixSec / 86400) * 86400
}

// LocalMinutesToUnixUTC converts minutes-from-midnight expressed in a local
// timezone into a UTC Unix timestamp in seconds for the given date.
// dateSec is a Unix timestamp in seconds (UTC) that identifies the target date.
// localMinutes is minutes elapsed since midnight in the local timezone.
// tz is an IANA timezone name (e.g. "America/New_York"); falls back to UTC if invalid.
func LocalMinutesToUnixUTC(dateSec int64, localMinutes int, tz string) int64 {
	return provider.LocalMinutesToUnixUTC(dateSec, localMinutes, tz)
}
```

Also add `LocalMinutesToUnixUTC` to the `timeProvider` interface:

```go
LocalMinutesToUnixUTC(dateSec int64, localMinutes int, tz string) int64
```

### Step 2 — `backStlib.go`: stdlib implementation

Add method to `timeServer`. Uses `time.LoadLocation` (stdlib use is legitimate here,
this is the backend implementation layer of the ecosystem time package):

```go
func (ts *timeServer) LocalMinutesToUnixUTC(dateSec int64, localMinutes int, tz string) int64 {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		loc = time.UTC
	}
	utcDate := time.Unix(dateSec, 0).UTC()
	hour := localMinutes / 60
	minute := localMinutes % 60
	localTime := time.Date(utcDate.Year(), utcDate.Month(), utcDate.Day(), hour, minute, 0, 0, loc)
	return localTime.Unix()
}
```

### Step 3 — `frontWasm.go`: WASM implementation

WASM has no access to the IANA timezone database. Use the ecosystem's global
offset (set by `SetTimeZoneOffset`) as the approximation:

```go
func (tc *timeClient) LocalMinutesToUnixUTC(dateSec int64, localMinutes int, tz string) int64 {
	offsetSec := int64(getOffsetMinutes()) * 60
	midnight := (dateSec / 86400) * 86400
	return midnight + int64(localMinutes)*60 - offsetSec
}
```

---

## Step 4 — Tests

### 4.1 Add shared test functions to `data_test.go`

```go
// Test Weekday
func WeekdayShared(t *testing.T) {
	// 2021-01-01 00:00:00 UTC was a Friday (5)
	result := time.Weekday(1609459200)
	if result != 5 {
		t.Errorf("Weekday(2021-01-01) = %d; want 5 (Friday)", result)
	}

	// 1970-01-01 was a Thursday (4)
	result = time.Weekday(0)
	if result != 4 {
		t.Errorf("Weekday(1970-01-01) = %d; want 4 (Thursday)", result)
	}

	// 2024-01-07 was a Sunday (0)
	result = time.Weekday(1704585600)
	if result != 0 {
		t.Errorf("Weekday(2024-01-07) = %d; want 0 (Sunday)", result)
	}

	t.Logf("Weekday tests passed")
}

// Test MidnightUTC
func MidnightUTCShared(t *testing.T) {
	// Midpoint of 2021-01-01: should return start of that day
	midnight := time.MidnightUTC(1609459200 + 3600) // 01:00:00 UTC
	if midnight != 1609459200 {
		t.Errorf("MidnightUTC = %d; want 1609459200", midnight)
	}

	// Already midnight — idempotent
	midnight = time.MidnightUTC(1609459200)
	if midnight != 1609459200 {
		t.Errorf("MidnightUTC(already midnight) = %d; want 1609459200", midnight)
	}

	// Last second of a day: 2021-01-01 23:59:59
	midnight = time.MidnightUTC(1609545599)
	if midnight != 1609459200 {
		t.Errorf("MidnightUTC(23:59:59) = %d; want 1609459200", midnight)
	}

	t.Logf("MidnightUTC tests passed")
}

// Test LocalMinutesToUnixUTC
func LocalMinutesToUnixUTCShared(t *testing.T) {
	// UTC offset 0: 09:00 local == 09:00 UTC
	// date: 2021-01-01 00:00:00 UTC (1609459200)
	result := time.LocalMinutesToUnixUTC(1609459200, 9*60, "UTC")
	expected := int64(1609459200 + 9*3600)
	if result != expected {
		t.Errorf("LocalMinutesToUnixUTC(UTC, 09:00) = %d; want %d", result, expected)
	}

	// UTC-5 (America/New_York winter): 09:00 local == 14:00 UTC
	result = time.LocalMinutesToUnixUTC(1609459200, 9*60, "America/New_York")
	expected = int64(1609459200 + 14*3600)
	if result != expected {
		t.Errorf("LocalMinutesToUnixUTC(America/New_York, 09:00) = %d; want %d", result, expected)
	}

	// Invalid timezone falls back to UTC
	result = time.LocalMinutesToUnixUTC(1609459200, 9*60, "Invalid/Zone")
	expectedFallback := int64(1609459200 + 9*3600)
	if result != expectedFallback {
		t.Errorf("LocalMinutesToUnixUTC(invalid tz) = %d; want %d (UTC fallback)", result, expectedFallback)
	}

	t.Logf("LocalMinutesToUnixUTC tests passed")
}
```

### 4.2 Register in `provider_test.go`

Add to `RunAPITests`:

```go
t.Run("Weekday", func(t *testing.T) { WeekdayShared(t) })
t.Run("MidnightUTC", func(t *testing.T) { MidnightUTCShared(t) })
t.Run("LocalMinutesToUnixUTC", func(t *testing.T) { LocalMinutesToUnixUTCShared(t) })
```

### 4.3 WASM test note

`LocalMinutesToUnixUTC` in WASM uses the global offset instead of IANA lookup.
The existing `frontWasm_test.go` calls `RunAPITests` — the UTC case will pass
unchanged. The timezone-specific cases will diverge (expected behavior since
WASM cannot resolve IANA names); skip them with a build tag if needed.

---

## Step 5 — Documentation

Update `docs/BROWSER_TEST.md` or create `docs/CALENDAR_FUNCTIONS.md` documenting:

- `Weekday`: pure arithmetic, no provider dependency, identical result on back and front
- `MidnightUTC`: pure arithmetic, identical on back and front
- `LocalMinutesToUnixUTC`: back uses full IANA resolution; front uses global offset set via `SetTimeZoneOffset`

---

## Step 6 — Publish

> **Prerequisite**: use `gotest` (not `go test`) to run tests in both environments.
> Install once: `go install github.com/tinywasm/devflow/cmd/gotest@latest`
> `gotest` runs vet, race detection, stdlib tests, and WASM tests automatically.

```bash
gotest .                   # backend tests
gotest .                   # all tests pass
git add api.go backStlib.go frontWasm.go data_test.go provider_test.go
git commit -m "feat: add Weekday, MidnightUTC, LocalMinutesToUnixUTC"
git tag v0.5.0
git push && git push --tags
```

---

## Closure criteria

```bash
gotest .                                         # no errors (backend + WASM)
grep -n "Weekday\|MidnightUTC\|LocalMinutes" api.go  # 3 matches
go build ./...                                   # no errors
```

---

## Files to modify

| File | Action |
|------|--------|
| `api.go` | Add `Weekday`, `MidnightUTC`, `LocalMinutesToUnixUTC` public functions + extend `timeProvider` interface |
| `backStlib.go` | Add `LocalMinutesToUnixUTC` method to `timeServer` |
| `frontWasm.go` | Add `LocalMinutesToUnixUTC` method to `timeClient` |
| `data_test.go` | Add `WeekdayShared`, `MidnightUTCShared`, `LocalMinutesToUnixUTCShared` |
| `provider_test.go` | Register 3 new test runners in `RunAPITests` |
| `docs/CALENDAR_FUNCTIONS.md` | Document the 3 new functions and their back/front behavior |

## Execution order

```
Step 1 — api.go (pure functions + interface extension)
Step 2 — backStlib.go (stdlib implementation)
Step 3 — frontWasm.go (WASM implementation)
Step 4 — data_test.go + provider_test.go (tests)
Step 5 — docs/CALENDAR_FUNCTIONS.md (documentation)
Step 6 — go test ./... + tag v0.5.0 + push
```
