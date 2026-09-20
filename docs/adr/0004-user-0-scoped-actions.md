# 0004. User 0 scoped actions

Date: 2026-09-20

## Status

Accepted

## Context

Removing bloat from a stock ROM without root has several degrees of severity:

- Full uninstall of a system app (`pm uninstall <pkg>`) — not possible without
  root for preinstalled packages.
- Per-user removal (`pm uninstall --user 0 <pkg>`) — the app disappears for
  user 0 (the owner), but the system APK stays on the device.
- Restriction (`freeze`/`disable`) — the app stays installed but is blocked
  from running or is disabled for the user.

The default action must be recoverable: users run this on their daily phone.

## Decision

Scope all mutating commands to user 0 and keep each one reversible by design.
`uninstall` maps to `pm uninstall --user 0`, `disable` to
`pm disable-user --user 0`, `freeze` to an app-ops background restriction, and
`reinstall` to `cmd package install-existing`. Every mutating command supports
`--dry-run`, which prints the exact `adb` commands without touching the device.

## Consequences

- No root is required, and a wrong removal is usually recoverable with
  `reinstall` as long as the ROM still carries the package.
- Storage is not actually reclaimed: the APK remains on the system image.
- `reinstall` fails after a ROM update drops the package, so the documented
  recovery path ends at a factory reset.
- Other users and work profiles are unaffected; this is a feature for shared
  devices and a limitation for users who expected a device-wide removal.
