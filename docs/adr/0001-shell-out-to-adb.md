# 0001. Shell out to the `adb` CLI

Date: 2026-09-20

## Status

Accepted

## Context

The tool must run package queries and mutations on a connected phone. The
alternatives were:

- Reimplement the ADB wire protocol (and USB/TCP transport) in Go.
- Link against an Android library such as `mobiledevice`/`adbkit` equivalents.
- Invoke the `adb` binary that every target user already needs installed.

Users of a debloat tool already have `adb` (it is the documented prerequisite),
and the commands we need (`pm list packages`, `pm uninstall --user 0`, `cmd
appops`, `pm disable-user`, `cmd package install-existing`) are stable
shell-level interfaces available without root.

## Decision

Shell out to `adb` via `os/exec`, pass a fixed set of subcommands, and parse
plain stdout. Device selection stays in `adb`'s own mechanism: honor
`ANDROID_SERIAL` rather than managing device state ourselves.

## Consequences

- No protocol, driver, or USB code to maintain; bugs in transport are `adb`'s
  to fix, and `adb` features (Wi-Fi debugging, `adb connect`) work for free.
- Output parsing must tolerate `adb` version and OEM formatting differences;
  the parser favors plain name-per-line output.
- Runtime dependency on an external binary on `PATH`; error messages must
  point users at installing `adb`.
- No root is required or attempted, because every command used is a shell
  command scoped to the user.
