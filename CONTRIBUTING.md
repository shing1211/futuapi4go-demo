# Contributing to this project

## Workflow

All changes travel through GitHub PRs:

1. **Create an issue first** — before any branch or code
2. **Branch** from 
3. **Commit** with a conventional message referencing the issue: `fix(#N): description`, `feat(#N): description`
4. **Open a PR** with `Closes #N` in the body
5. **Self-review** the diff before requesting review
6. **Merge** via squash merge (fast-forward merge for hotfixes)
7. **Delete** the branch after merge

## Code style

- Go: run `go fmt` and `go vet` before committing
- Python: run `black` and `ruff check`
- Shell: shellcheck for shell scripts

## Tests

- New features must include tests
- Bug fixes must include a regression test
- Run the full test suite before opening a PR:
  ```bash
  go test ./...   # Go projects
  python3 -m pytest  # Python projects
  ```

## Commit message format

```
<type>(<scope>): <short description>

[optional body]

Closes #<issue-number>
```

Types: `fix`, `feat`, `docs`, `chore`, `refactor`, `test`

## Branch protection

- `main` is protected; force-push is disabled
- All PRs require at least one review before merge
- CI must pass before merge

## Questions

Open an issue for discussion before starting work on non-trivial changes.

