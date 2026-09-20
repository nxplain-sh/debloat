# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to
[Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- `package.json` with Prettier as a dev dependency; run `npm install` once,
  then `make format` / `make format-check`.
- Architecture decision records under `docs/adr/`.
- Git pre-commit hook (`.githooks/pre-commit`) running `fmt-check`, `vet`, and
  tests; enable it with `make hooks`.

### Changed

- Repository moved to <https://github.com/nxplain-sh/debloat>; the Go module
  path is now `github.com/nxplain-sh/debloat`.
- Renamed the command and binary from `vivo-debloater` to `debloat`
  (`cmd/debloat`, `./bin/debloat`).
- Prettier config renamed from `.prettierrc.yaml` to `.prettierrc`.
- Renamed Makefile targets `prettier`/`prettier-check` to
  `format`/`format-check`; Prettier now resolves from `package.json`.

## [0.1.0] - 2026-07-04

### Added

- No-root CLI for listing and removing preinstalled apps on Vivo (BBK) phones
  running OriginOS/Funtouch over `adb`, built as a single self-contained Go
  binary for Linux, macOS, and Windows.
- Commands: `show`, `list`, `uninstall`, `freeze`, `disable`, `reinstall`,
  `version`, and `help`.
- `-f`/`--file` package-list flag, `-o`/`--output` for `show`, and
  `-n`/`--dry-run` preview for the mutating commands.
- Default debloat list in `packages.txt`, sourced from Technastic.
- Windows setup docs (platform-tools, vivo USB driver, PowerShell).
- golangci-lint config, Prettier formatting, EditorConfig, and CI/release
  workflows (GoReleaser).

[Unreleased]: https://github.com/nxplain-sh/debloat/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/nxplain-sh/debloat/releases/tag/v0.1.0
