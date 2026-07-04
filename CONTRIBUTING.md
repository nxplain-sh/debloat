# Contributing to vivo-debloater

Thanks for your interest in improving vivo-debloater! This is a small,
no-root ADB helper for debloating Vivo (BBK) phones. Contributions of all
sizes are welcome — bug fixes, package-list corrections, and documentation.

## Ground rules

- Be kind and constructive. See the [Code of Conduct](CODE_OF_CONDUCT.md).
- Keep changes focused. One logical change per pull request.
- The CLI is written in **Go with no external dependencies** (standard library
  only). `adb` stays an external runtime dependency. Please don't add third-party
  modules without a good reason.

## Getting set up

```bash
git clone https://github.com/miguelmartens/vivo-debloater.git
cd vivo-debloater
make install-adb        # optional: adb via Homebrew (android-platform-tools)
make build              # ./bin/vivo-debloater
```

You can develop and test most changes **without a phone** using `--dry-run`,
which prints the `adb` commands that would run:

```bash
make run ARGS='uninstall --dry-run'
```

## Before you open a pull request

Run the same checks CI runs:

```bash
make check            # gofmt check + go vet + golangci-lint + go test -race
make shellcheck       # lint the legacy scripts/debloat.sh (requires shellcheck)
make prettier-check   # verify docs/config formatting
```

- Go code must be **`gofmt`-clean** and pass **`go vet`** and **`golangci-lint`**
  (config in `.golangci.yaml`; install it from
  [golangci-lint.run](https://golangci-lint.run/welcome/install/)).
- Add or update **tests** for behavior changes (`internal/...`).
- The legacy **`scripts/debloat.sh`** must still pass `shellcheck`.
- Docs and config must be **Prettier-clean** (`make prettier` to auto-format).
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
