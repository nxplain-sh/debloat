# Security Policy

## Scope

vivo-debloater is a set of local `bash` + `adb` helpers. It has no server, no
network service, and stores no credentials. The main risk surface is:

- shell handling of package names and file paths, and
- the `adb` commands the tool runs against a connected device.

## Reporting a vulnerability

Please **do not** open a public issue for security problems. Instead, email
**mail@miguelmartens.com** with:

- a description of the issue and its impact,
- steps to reproduce (redact device serials), and
- any suggested fix.

You can expect an acknowledgement within a few days. Coordinated disclosure is
appreciated.

## Good practice for users

- Review `packages.txt` before applying changes.
- Use `--dry-run` to preview the exact `adb` commands.
- Keep a `show` dump of your ROM so you can recover with `reinstall`.
