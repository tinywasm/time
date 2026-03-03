# Plan: Add UTC / ISO8601 formatting to tinywasm/time

**Context:** The `tinywasm/time` package currently generates strings formatted for end-users, applying local time offsets unconditionally. To support systems like FHIR which require raw UTC timestamps in the standard `RFC3339`/`ISO8601` format (`YYYY-MM-DDTHH:MM:SSZ`), we need a dedicated formatter that ignores the offset and returns strict UTC.

**Goal:** Implement `FormatISO8601(nano int64) string` to output `YYYY-MM-DDTHH:MM:SSZ` in strict UTC (no offset), working natively across `!wasm` (stdlib time) and `wasm` (Vanilla JS Date).

---

## Stage 1: API Surface

### 1. Update `timeProvider` interface in `api.go`

Add the new method to the `timeProvider` interface:

```go
type timeProvider interface {
	// ... existing methods
	FormatISO8601(nano int64) string
}
```

### 2. Export `FormatISO8601` in `api.go`

Add the public function wrapper that delegates to the provider:

```go
// FormatISO8601 formats a UnixNano timestamp into an ISO 8601 string (UTC).
// Format: "YYYY-MM-DDTHH:MM:SSZ"
func FormatISO8601(nano int64) string {
	return provider.FormatISO8601(nano)
}
```

---

## Stage 2: Backend Implementation (`backStlib.go`)

Implement the method for the standard library provider.

### Update `timeServer {}` in `backStlib.go`:

```go
func (ts *timeServer) FormatISO8601(nano int64) string {
	t := time.Unix(0, nano).UTC()
	return t.Format(time.RFC3339) 
}
```
*(Note: `time.RFC3339` exactly maps to "2006-01-02T15:04:05Z" under UTC)*

---

## Stage 3: Frontend Implementation (`frontWasm.go`)

Implement the method for the WebAssembly provider using JavaScript's native Date.

### Update `timeClient {}` in `frontWasm.go`:

```go
func (tc *timeClient) FormatISO8601(nano int64) string {
	// We MUST NOT use applyOffset here. We need raw UTC.
	jsDate := tc.dateCtor.New(float64(nano) / 1e6) 
	
	// 'toISOString' natively outputs 'YYYY-MM-DDTHH:mm:ss.sssZ'
	iso := jsDate.Call("toISOString").String()
	
	// We slice out the milliseconds to strictly match 'YYYY-MM-DDTHH:MM:SSZ'
	return iso[0:19] + "Z"
}
```

---

## Stage 4: Testing & Documentation

### 1. Update `data_test.go`
Add a dedicated test case for `FormatISO8601` to ensure both UTC formatting and correct slicing behavior.

```go
func TestFormatISO8601(t *testing.T) {
	// Unix timestamp for 2024-01-15 15:30:45 UTC
	nano := int64(1705332645000000000) 
	expected := "2024-01-15T15:30:45Z"

	result := time.FormatISO8601(nano)
	if result != expected {
		t.Errorf("FormatISO8601 failed: expected %%s, got %%s", expected, result)
	}
}
```

### 2. Update `README.md`

**A. Clarify existing formatting functions:**
Update the description under **Display Formatting** to make it extremely clear that existing functions apply local time offsets.

```markdown
### Display Formatting
**⚠️ IMPORTANT:** All formatting functions below (`FormatDate`, `FormatTime`, `FormatDateTime`, etc.) automatically apply the current timezone offset to display **local time**. If you need raw UTC time, use `FormatISO8601` instead.
```

**B. Add documentation for the new `FormatISO8601` function:**

```markdown
#### `FormatISO8601(nano int64) string`
Formats a UnixNano timestamp into an ISO 8601 string: "YYYY-MM-DDTHH:MM:SSZ".
Unlike other formatting functions, this strictly outputs **UTC time** and ignores any local timezone offsets. Often used for DB records, HL7/FHIR, and REST APIs.
```

---

## Execution Command

```bash
gotest
gopush 'feat: add FormatISO8601 for strict UTC timezone-free formatting'
```
