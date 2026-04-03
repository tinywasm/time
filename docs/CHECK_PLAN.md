# PLAN: Añadir FormatCompact a tinywasm/time

## Motivación
`tinywasm/pdf` necesita formatear fechas como `"YYYYMMDDHHmmss"` para PDF metadata
(`/CreationDate D:20260402153045`). Actualmente no existe un formato compacto sin separadores.

## Estado actual de formatos

| Función              | Formato salida             | Ejemplo                |
|----------------------|---------------------------|------------------------|
| FormatDate           | `YYYY-MM-DD`              | `2026-04-02`           |
| FormatTime           | `HH:MM:SS`                | `15:30:45`             |
| FormatDateTime       | `YYYY-MM-DD HH:MM:SS`    | `2026-04-02 15:30:45`  |
| FormatDateTimeShort  | `YYYY-MM-DD HH:MM`       | `2026-04-02 15:30`     |
| FormatISO8601        | `YYYY-MM-DDTHH:MM:SSZ`   | `2026-04-02T15:30:45Z` |
| **FormatCompact**    | **`YYYYMMDDHHmmss`**      | **`20260402153045`**   |

---

## Stage 1: Implementación

**Riesgo**: bajo | **Complejidad**: baja

### 1.1 API pública

Añadir a `api.go`:
```go
// FormatCompact formats a UnixNano timestamp into a compact string "YYYYMMDDHHmmss" (UTC).
// Useful for PDF metadata dates, file naming, and compact timestamps.
func FormatCompact(nano int64) string {
    return provider.FormatCompact(nano)
}
```

### 1.2 Interface

Añadir método a `timeProvider` en `api.go`:
```go
FormatCompact(nano int64) string
```

### 1.3 Backend (`backStlib.go`)

```go
func (ts *timeServer) FormatCompact(nano int64) string {
    t := time.Unix(0, nano).UTC()
    return t.Format("20060102150405")
}
```

### 1.4 Frontend WASM (`frontWasm.go`)

```go
func (tc *timeClient) FormatCompact(nano int64) string {
    jsDate := tc.dateCtor.New(float64(nano) / 1e6)
    iso := jsDate.Call("toISOString").String()
    // iso = "YYYY-MM-DDTHH:mm:ss.sssZ"
    // extraer: YYYY(0:4) MM(5:7) DD(8:10) HH(11:13) mm(14:16) ss(17:19)
    return iso[0:4] + iso[5:7] + iso[8:10] + iso[11:13] + iso[14:16] + iso[17:19]
}
```

### 1.5 Archivos a modificar

| Archivo          | Cambio |
|------------------|--------|
| `api.go`         | Añadir `FormatCompact()` función pública + método en interface `timeProvider` |
| `backStlib.go`   | Implementar `FormatCompact` en `timeServer` |
| `frontWasm.go`   | Implementar `FormatCompact` en `timeClient` |

---

## Stage 2: Tests

Seguir el patrón compartido existente del proyecto:
- Tests compartidos como funciones exportadas `XxxShared(t)` en `data_test.go`
- Registrados en `provider_test.go` dentro de `RunAPITests()`
- Ejecutados por `backStlib_test.go` (!wasm) y `frontWasm_test.go` (wasm)

### 2.1 Añadir test compartido en `data_test.go`

```go
// Test FormatCompact
func FormatCompactShared(t *testing.T) {
    // Known timestamp: 2024-01-15 15:30:45 UTC
    nano := int64(1705332645000000000)
    expected := "20240115153045"

    result := time.FormatCompact(nano)
    if result != expected {
        t.Errorf("FormatCompact(%d) = %s; want %s", nano, result, expected)
    }

    // Test with zero value (epoch)
    result = time.FormatCompact(int64(0))
    if result != "19700101000000" {
        t.Errorf("FormatCompact(0) = %s; want 19700101000000", result)
    }

    // Test length is always 14 characters
    if len(result) != 14 {
        t.Errorf("FormatCompact length = %d; want 14", len(result))
    }

    // Test with current time
    currentNano := time.Now()
    result = time.FormatCompact(currentNano)
    if len(result) != 14 {
        t.Errorf("FormatCompact(current) length = %d; want 14", len(result))
    }

    t.Logf("FormatCompact tests passed")
}
```

### 2.2 Registrar en `provider_test.go`

Añadir dentro de `RunAPITests()`:
```go
t.Run("FormatCompact", func(t *testing.T) { FormatCompactShared(t) })
```

### 2.3 Archivos a modificar

| Archivo             | Cambio |
|---------------------|--------|
| `data_test.go`      | Añadir `FormatCompactShared(t)` |
| `provider_test.go`  | Registrar en `RunAPITests()` |

No se necesita modificar `backStlib_test.go` ni `frontWasm_test.go` — ya llaman `RunAPITests()`.

---

## Validación

### Requisito: gotest
Instalar si no está disponible:
```bash
go install github.com/nicholasgasior/gopher-wasm/cmd/gotest@latest
```

### Ejecutar tests

**Backend (stdlib)**:
```bash
go test ./...
```

**WASM (browser)**:
```bash
gotest
```

Ambos deben pasar, incluyendo el nuevo test `FormatCompact`.

### Criterios de aceptación
1. `FormatCompact` retorna exactamente 14 caracteres
2. Formato es `YYYYMMDDHHmmss` sin separadores
3. Usa UTC (sin timezone offset)
4. Tests pasan en backend (`go test`) y WASM (`gotest`)

---

## Stage 3: Documentación

### 3.1 Actualizar `README.md`

Añadir en la sección **Display Formatting** (después de `FormatISO8601`):

```markdown
#### `FormatCompact(nano int64) string`
Formats a UnixNano timestamp into a compact string: "YYYYMMDDHHmmss".
Outputs **UTC time**, ignoring timezone offsets. Useful for PDF metadata dates, file naming, and compact timestamps.
```

### 3.2 Actualizar tabla de formatos en `README.md`

Si existe una tabla de formatos, añadir la fila de `FormatCompact`.

### 3.3 Actualizar ejemplo en `README.md` sección Quick Start

Añadir al ejemplo existente:
```go
    // Compact format (UTC, no separators)
    compact := time.FormatCompact(nano)
    println("Compact:", compact) // "20260402153045"
```

### 3.4 Archivos a modificar

| Archivo     | Cambio |
|-------------|--------|
| `README.md` | Añadir `FormatCompact` en sección Display Formatting (después de `FormatISO8601`) y ejemplo en Quick Start |
