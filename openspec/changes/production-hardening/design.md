# Design

## Context

futuapi4go-demo is a CLI demo repo with 128 example directories, one test file (`examples/pkg/connect/connect_test.go`), and existing CI (`ci.yml`, `codeql.yml`, `govulncheck.yml`, `pages.yml`). Dependencies are clean: no xerrors, no golang/protobuf. `google.golang.org/protobuf` v1.36.12 is indirect via otel/prometheus.

## Goals / Non-Goals

**Goals:**
- Initialize OpenSpec root
- Add a smoke test proving SDK wiring without credentials
- Document transitive dependency risks

**Non-Goals:**
- No test coverage for the 128 examples (each requires a live OpenD instance)
- No dependency upgrades (already current)
- No CI changes (already comprehensive)

## Decisions

### 1. Smoke test uses examples/pkg/connect

The existing `examples/pkg/connect` package already has a test file. A smoke test there proves the SDK's `client.New()` and `Connect` wiring without requiring a live OpenD instance — connection failures are expected and tolerated.

### 2. Dependency notes are minimal

No xerrors, no golang/protobuf. The only transitive concern is `google.golang.org/protobuf` v1.36.12 (indirect, current version, no known CVEs). Documented for completeness.

## Risks / Trade-offs

- **Smoke test cannot verify live connectivity**: The SDK requires a running OpenD instance. A wiring test only proves the client can be constructed. Accepted — the same limitation applies to all demo repos.

## Migration Plan

No migration needed — all changes are additive.

## Open Questions

None.
