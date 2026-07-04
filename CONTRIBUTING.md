# Contributing to vivo-debloater

Thanks for your interest in improving vivo-debloater! This is a small,
no-root ADB helper for debloating Vivo (BBK) phones. Contributions of all
sizes are welcome — bug fixes, package-list corrections, and documentation.

## Ground rules

- Be kind and constructive. See the [Code of Conduct](CODE_OF_CONDUCT.md).
- Keep changes focused. One logical change per pull request.
- The tooling stays **dependency-light**: POSIX-ish `bash`, `adb`, and a
  `Makefile` facade. Please don't introduce heavyweight runtimes.

## Getting set up

```bash
git clone git@github.com:miguelmartens/vivo-debloater.git
cd vivo-debloater
make install-adb        # optional: adb via Homebrew (android-platform-tools)
make install-prettier   # optional: Prettier for formatting
```

You can develop and test most changes **without a phone** using `--dry-run`,
which prints the `adb` commands that would run:

```bash
make uninstall ARGS='--dry-run'
```

## Before you open a pull request

Run the same checks CI runs:

```bash
make shellcheck       # lint scripts/debloat.sh (requires shellcheck on PATH)
make prettier-check   # verify formatting
```

- **`scripts/debloat.sh`** must pass `shellcheck` with no warnings.
- All files must be Prettier-clean (`make prettier` to auto-format).
- If you touch behavior, update `README.md` to match.

## Editing the package list

`packages.txt` is the default debloat list. When proposing changes:

- One package name per line; use `#` to comment out entries.
- Explain **why** a package is safe to remove (or why one should be removed
  from the default list because it risks breaking calling, SMS/MMS, updates,
  or security features).
- Prefer commenting out risky entries over deleting them, so users can opt in.

## Commit messages

Short, imperative subject lines (e.g. `add com.vivo.weather to default list`).
Reference an issue number where relevant.

## Reporting bugs

Open an issue using the templates in `.github/ISSUE_TEMPLATE/`. Please include
your device model, OriginOS/Funtouch version, and the exact command and output
(redact serials).

## Safety reminder

This tool changes what runs on real devices. When in doubt, test with
`--dry-run`, keep a `show` dump of your ROM, and know how to recover
(`reinstall` or, worst case, a factory reset).
