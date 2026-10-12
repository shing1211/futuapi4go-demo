# Production Hardening

## Why

futuapi4go-demo lacks an OpenSpec root, has no smoke test, and has no dependency documentation. Bringing it in line with the other demo repos (tiger-go-demo, longbridge-go-demo) reduces maintenance friction and documents transitive dependency risks.

## What Changes

- Initialize OpenSpec (`openspec init`)
- Add `test/smoke_test.go` — credentials-free wiring test
- Add `docs/dependency-notes.md` — document transitive dependency risks

## Capabilities

None — this is a pure tooling change with no spec-level behaviour changes. `skip_specs: true`.

## Impact

- Affected code: none (new files only)
- New files: `test/smoke_test.go`, `docs/dependency-notes.md`, `openspec/` root
- New dependencies: none
- Affected specs: none
- Affected workflows: none
