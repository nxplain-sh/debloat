# 0002. Standard library only, single static binary

Date: 2026-09-20

## Status

Accepted

## Context

The CLI is small: parse flags, run a handful of `adb` commands, read a text
file, print output. The Go module ecosystem offers convenience packages for
each of these (flag frameworks, output formatting, semver, testing helpers),
but each one adds supply-chain surface, version churn, and a `go.sum` to a
tool whose value is being auditable and easy to build.

Cross-platform distribution (Linux, macOS, Windows) is a first-class goal, and
`CGO_ENABLED=0` static binaries from the standard library alone satisfy it.

## Decision

Use the Go standard library exclusively. Do not add third-party module
dependencies without a strongly justified reason; keep the module
dependency-free so `go build` needs only a Go toolchain.

## Consequences

- One `go build` produces a self-contained binary; releases and local builds
  behave identically.
- No third-party CVE monitoring or dependency updates for Go code; Renovate
  only tracks GitHub Actions.
- Some conveniences are hand-rolled (argument parsing, output alignment,
  assertions in tests). This is accepted because the surface is small.
- The rule is deliberately strict: any future dependency needs an ADR
  explaining why the standard library is insufficient.
