# Tasks

## 1. OpenSpec setup

- [x] 1.1 `openspec init` — done

## 2. Smoke test

- [x] 2.1 Create `test/smoke_test.go` with SDK wiring test
- [x] 2.2 Verify `go vet ./test/` and `go test ./test/` pass

## 3. Dependency notes

- [x] 3.1 Create `docs/dependency-notes.md` documenting transitive deps
- [x] 3.2 Verify `go mod why` chain matches documented claims

## 4. Verify and PR

- [x] 4.1 `go vet ./test/` clean (full `./...` vet requires GOCACHE on disk — tmpfs constraint)
- [x] 4.2 `gofmt -l test/` returns nothing
- [ ] 4.3 `go build ./...` passes (deferred — disk constraint)
- [ ] 4.4 `go test ./...` passes (deferred — disk constraint)
- [ ] 4.5 Branch, commit, push, open PR for Mary review
- [ ] 4.6 Close issue #14 after merge
