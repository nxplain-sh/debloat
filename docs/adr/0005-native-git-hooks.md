# 0005. Native git hooks over a hook framework

Date: 2026-09-20

## Status

Accepted

## Context

Contributors should run the cheap correctness checks before committing:
`gofmt`, `go vet`, and tests. Options considered:

- Husky + lint-staged (requires Node tooling in the repo).
- The `pre-commit` framework (requires Python on every contributor machine
  and manages its own tool environments).
- Native git hooks: a committed `.githooks/` directory enabled with
  `git config core.hooksPath .githooks`.

Only Prettier already needs Node, and it is optional for Go-only
contributions. Adding a second runtime and a hook manager to run three
already-installed Go commands is disproportionate for a dependency-free
repository.

## Decision

Maintain the hook as a small shell script in `.githooks/pre-commit`, enabled
per clone with `make hooks` (`git config core.hooksPath .githooks`). The hook
runs `make fmt-check vet test` only; golangci-lint and Prettier stay in CI and
in the pre-pull-request checklist, because they are slower or need extra
tooling.

## Consequences

- No new runtime or hook manager; the hook is readable shell plus make.
- Enabling hooks is opt-in per clone; forgetting `make hooks` means no local
  guard, and CI still catches the same problems.
- Commit-time checks are a subset of `make check`; a green commit is not proof
  of a green CI run.
- Hook content lives in the repository and is versioned with the code, so
  changes to it are reviewed like any other change.
