# Implementation Plan

This document tracks planned enhancements for the futuapi4go-demo project.
It is a living document; items are added, refined, or closed as work progresses.

---

## Legend

| Priority | Label | Description |
|----------|-------|-------------|
| P0 | Critical | Blocks demo from running cleanly |
| P1 | High | Significant usability or robustness improvement |
| P2 | Medium | Nice-to-have quality-of-life improvement |
| P3 | Low | Future idea, nice but not urgent |

| Status | Meaning |
|--------|---------|
| `todo` | Not started |
| `in_progress` | Being worked on |
| `done` | Completed |

---

## Items

### 1. Structured Error Handling — `must(label, err)`
**Priority:** P1 | **Status:** `todo`

Replace the `must(err)` helper with a labeled version that prints which
demo section and API call failed:

```go
func must(label string, err error) {
    if err != nil {
        yellow(fmt.Sprintf("  [%s] WARNING: %v (continuing anyway)\n", label, err))
    }
}
```

All call sites in `cmd/demo/main.go` need updating, e.g.:
```go
must("GetGlobalState", sys.GetGlobalState(cli.Inner()))
```

---

### 2. Context with Timeout
**Priority:** P1 | **Status:** `todo`

Wrap all API calls with `context.WithTimeout` to prevent indefinite hangs,
especially on network issues or when OpenD is unresponsive.

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
// pass ctx to calls that accept it
```

Affected areas:
- All `sys.*`, `qot.*`, `trd.*` calls
- Especially `demoPushSubscriptions` (currently blocks forever)

---

### 3. Graceful Shutdown
**Priority:** P2 | **Status:** `todo`

Add a timeout to the `signal.Notify` in `demoPushSubscriptions` so the
program exits cleanly if something goes wrong:

```go
sig := make(chan os.Signal, 1)
signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
select {
case <-sig:
    // user pressed Ctrl+C
case <-time.After(60 * time.Second):
    // timeout — forced exit
}
```

---

### 4. Env-Var Config for Demo Securities
**Priority:** P2 | **Status:** `todo`

Move hardcoded market constants and stock codes to environment variables
with sensible defaults:

| Variable | Default | Description |
|----------|---------|-------------|
| `FUTU_ADDR` | `127.0.0.1:11111` | OpenD address (already supported) |
| `FUTU_SKIP_TRADING` | `0` | Set to `1` to skip all trading API calls |
| `FUTU_SEC_HK` | `00700` | Hong Kong demo security |
| `FUTU_SEC_US` | `AAPL` | US demo security |

Replace all hardcoded `sec(MarketHK, "00700")` calls with values read from env.

---

### 5. Config File Support
**Priority:** P2 | **Status:** `todo`

Add optional `config.yaml` (or `.env`) to allow users to configure demo
behavior without recompiling:

```yaml
# config.yaml
address: "127.0.0.1:11111"
securities:
  HK: "00700"
  US: "AAPL"
  skip_trading: false
output_format: "text"  # or "json"
```

Priority: load from env vars first, then config file, then hardcoded defaults.

---

### 6. Auto-Detect OpenD Port
**Priority:** P3 | **Status:** `todo`

Try common OpenD ports in sequence instead of requiring the user to
set `FUTU_ADDR`:

```go
ports := []string{"11111", "22222", "33333"}
for _, port := range ports {
    addr := "127.0.0.1:" + port
    if cli.Connect(addr) == nil {
        return
    }
}
```

Show a clear error listing all ports tried.

---

### 7. Simulator Detection — Skip Trading Automatically
**Priority:** P1 | **Status:** `done`

`demoTrading` checks for `realAccID == 0` after `GetAccList` and skips
the entire trading section gracefully with a clear message.
No further changes needed.

---

### 8. Table-Formatted Output
**Priority:** P2 | **Status:** `todo`

Replace bare `fmt.Printf` rows with a simple ASCII table for multi-row
results such as positions, order fills, and stock filter results.

Helper idea:
```go
type Table struct {
    headers []string
    rows    [][]string
}
func (t *Table) Print()
```

Replace output in `demoTrading` (positions, orders, fills) and
`demoStockFilter`.

---

### 9. JSON Export Option
**Priority:** P3 | **Status:** `todo`

Add a `--json` flag or `FUTU_OUTPUT=json` env var to output all demo
results as machine-readable JSON for scripting and integration testing.

```bash
FUTU_OUTPUT=json go run ./cmd/demo/main.go <<< "0"
```

---

### 10. Parallel API Calls in Sections 2–3
**Priority:** P3 | **Status:** `todo`

In `demoMarketData` and `demoMarketAnalysis`, run independent API calls
concurrently using goroutines + `errgroup` for faster demo execution:

```go
eg, ctx := errgroup.WithContext(context.Background())
eg.Go(func() error { return callGetBasicQot(ctx) })
eg.Go(func() error { return callGetKL(ctx) })
// ...
if err := eg.Wait(); err != nil {
    must("parallel", err)
}
```

---

### 11. Pagination Helper for RequestHistoryKL
**Priority:** P2 | **Status:** `todo`

The `RequestHistoryKL` API returns `NextReqKey` for pagination. Add a
helper that fetches all pages automatically:

```go
func fetchAllHistoryKL(cli *client.Client, req *qot.RequestHistoryKLRequest) ([]*qotcommon.KLine, error)
```

Update `demoHistoricalData` to use the helper instead of manually
handling one page.

---

### 12. Unit Tests for Helper Functions
**Priority:** P2 | **Status:** `todo`

Add tests in `cmd/demo/main_test.go` (or a dedicated `pkg/util/`) for:

| Function | Test cases |
|----------|-----------|
| `formatVolume` | 0, 999, 1_000, 999_999, 1_000_000, 1_000_000_000 |
| `formatMoney` | same as above |
| `trdSideName` | all known enum values + unknown |
| `sec` | nil inputs, empty code |
| `ptrStr`, `ptrInt32`, etc. | verify pointer values |

---

### 13. GitHub Actions CI Pipeline
**Priority:** P2 | **Status:** `todo`

Add `.github/workflows/ci.yml`:

```yaml
name: CI
on: [push, pull_request]
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.26'
      - run: go vet ./...
      - run: go build ./...
      - run: go test ./...
```

---

### 14. Colored Price Change Output
**Priority:** P3 | **Status:** `todo`

In `demoMarketData`, color-code price changes instead of showing raw
numbers:

```go
chgRate := (q.CurPrice - q.LastClosePrice) / q.LastClosePrice * 100
if chgRate > 0 {
    green(fmt.Sprintf("+%.2f%%", chgRate))
} else if chgRate < 0 {
    red(fmt.Sprintf("%.2f%%", chgRate))
} else {
    fmt.Print(" 0.00%")
}
```

---

## Progress Summary

| # | Item | Priority | Status |
|---|------|----------|--------|
| 1 | Structured error handling (`must(label, err)`) | P1 | `todo` |
| 2 | Context with timeout | P1 | `todo` |
| 3 | Graceful shutdown | P2 | `todo` |
| 4 | Env-var config for demo securities | P2 | `todo` |
| 5 | Config file support | P2 | `todo` |
| 6 | Auto-detect OpenD port | P3 | `todo` |
| 7 | Simulator detection — skip trading | P1 | `done` |
| 8 | Table-formatted output | P2 | `todo` |
| 9 | JSON export option | P3 | `todo` |
| 10 | Parallel API calls | P3 | `todo` |
| 11 | Pagination helper | P2 | `todo` |
| 12 | Unit tests | P2 | `todo` |
| 13 | GitHub Actions CI | P2 | `todo` |
| 14 | Colored price changes | P3 | `todo` |

**P0 items:** none (project currently builds and runs cleanly)
**P1 items:** 2 open
**P2 items:** 8 open
**P3 items:** 3 open
**Total:** 14 items | **Done:** 1
