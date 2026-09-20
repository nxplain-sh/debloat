# debloat

A small, no-root CLI to list system packages on a Vivo (BBK) phone running
**OriginOS** (older models may still show **Funtouch** in settings) and remove
or restrict apps over `adb`. It's a single self-contained Go binary that runs on
**Linux, macOS, and Windows**. The default package list is based on
[Technastic — Vivo bloatware list](https://technastic.com/vivo-bloatware-preinstalled-apps-list/).

> `debloat` shells out to `adb`; it never talks to the phone directly, so
> you still need `adb` installed. It does **not** require root.

## Install

Pick one:

- **Prebuilt binary** — download the archive for your OS/arch from the
  [Releases](https://github.com/nxplain-sh/debloat/releases) page,
  extract it, and put `debloat` on your `PATH`.
- **Go toolchain**:
  ```bash
  go install github.com/nxplain-sh/debloat/cmd/debloat@latest
  ```
- **From source**:
  ```bash
  git clone https://github.com/nxplain-sh/debloat.git
  cd debloat
  make build      # produces ./bin/debloat
  ```

## Prerequisites

- **`adb`** on your `PATH`:
  - **macOS** — `make install-adb`, or `brew install --cask android-platform-tools`.
  - **Linux** — from your package manager (`sudo apt install adb`,
    `sudo dnf install android-tools`, or `sudo pacman -S android-tools`), or
    Google's [platform-tools](https://developer.android.com/tools/releases/platform-tools).
  - **Windows** — see below.
- On the phone: **Developer options** → **USB debugging** enabled; accept the
  computer's RSA fingerprint when prompted.
- One phone connected over USB (or Wi‑Fi debugging with `adb connect`). If
  several devices show up in `adb devices`, set **`ANDROID_SERIAL`** to the
  serial you want.

### Windows

1. **Install `adb`** — either
   `winget install --id Google.PlatformTools` (or `choco install adb`), or
   download Google's
   [platform-tools](https://developer.android.com/tools/releases/platform-tools)
   ZIP, extract it (e.g. to `C:\platform-tools`), and add that folder to your
   **PATH** (Settings → _Edit the system environment variables_ → **Environment
   Variables**).
2. **Install vivo's USB driver** so Windows can see the phone — install
   **vivo PC Suite** / **EasyShare**, or let **Windows Update** fetch the
   driver, then confirm the phone shows up in **Device Manager** with no
   warning icon. If `adb devices` lists the device as `unauthorized`, unlock
   the phone and accept the RSA prompt.
3. **Run the commands** from **PowerShell** or **Windows Terminal**. The binary
   is `debloat.exe`; to target a specific device, set the environment
   variable the PowerShell way:
   ```powershell
   $env:ANDROID_SERIAL = "ABC123"
   .\debloat.exe list
   ```

## Quick start

```bash
# Snapshot what the ROM ships (keep dumps in backups/, which is gitignored)
debloat show -o backups/my-phone-system.txt

# Edit packages.txt: comment out (#) anything you want to keep
debloat list                      # check install status vs packages.txt
debloat uninstall --dry-run       # preview adb commands (no changes)
debloat uninstall                 # apply — or: freeze | disable
```

By default the package list is `packages.txt` in the **current directory**, so
run these from the repo (or point `-f` at your own list from anywhere).

## Commands

| Command     | Purpose                                                                                                                          |
| ----------- | -------------------------------------------------------------------------------------------------------------------------------- |
| `show`      | Print **system** package names from the device (`pm list packages -s`, with `--user 0` when supported). One plain name per line. |
| `list`      | For each entry in the package list, print `installed` or `not installed` for user 0.                                             |
| `uninstall` | `pm uninstall --user 0 <package>` — removes the app **for user 0** (system APK usually stays on the device).                     |
| `freeze`    | `cmd appops set <package> RUN_IN_BACKGROUND ignore` — restricts background behavior (lighter than uninstall).                    |
| `disable`   | `pm disable-user --user 0 <package>`.                                                                                            |
| `reinstall` | `cmd package install-existing <package>` — brings back a system app for the user when the ROM still has it.                      |
| `version`   | Print the version.                                                                                                               |
| `help`      | Show usage.                                                                                                                      |

### Flags

- **`-f` / `--file`** — package list file (default: **`packages.txt`** in the
  current directory). Used by `list`, `uninstall`, `freeze`, `disable`, and
  `reinstall`.
- **`-o` / `--output`** — output file for **`show`** only. If omitted, the list
  is printed to stdout. Parent directories are created if missing.
- **`-n` / `--dry-run`** — for **`uninstall`**, **`freeze`**, **`disable`**, and
  **`reinstall`** only: print the `adb` commands that would run; **no device
  access** and no changes.

### Workflow

Typical order: snapshot the ROM, curate the list, verify, preview, then apply
(or recover).

1. **Save a system package dump** — full list from the phone; use `backups/` so
   it stays out of git.
   ```bash
   debloat show -o backups/vivo-x200-system.txt
   ```
2. **Edit `packages.txt`** — comment out (`#`) every package you want to keep.
   Optionally compare with your dump so you only target packages that exist on
   your ROM.
3. **Check install status**:
   ```bash
   debloat list
   ```
4. **Preview changes** — prints the `adb` commands only; does not touch the
   device:
   ```bash
   debloat uninstall --dry-run
   ```
   Same pattern for `freeze --dry-run` or `disable --dry-run`.
5. **Apply** — run the same command without `--dry-run` (pick **one** of
   uninstall, freeze, or disable depending on how aggressive you want to be):
   ```bash
   debloat uninstall
   ```
6. **Recover** — if you removed or disabled something you still need, and the
   system image still has the package:
   ```bash
   debloat reinstall
   ```

**Multiple devices** — set the serial when more than one device is connected:

```bash
ANDROID_SERIAL=ABC123 debloat list
```

## `packages.txt`

- One **package name** per line (e.g. `com.vivo.browser`).
- Lines starting with **`#`** are comments; blank lines are ignored.
- Start from the bundled list, compare with your own `show` export, and
  **comment out** anything you rely on (dialer, MMS, cloud backup you use, etc.).

## Project layout

Follows the [Standard Go Project Layout](https://github.com/golang-standards/project-layout):

- `cmd/debloat/` — `main`; thin entry point.
- `internal/cli/` — command dispatch, flags, and output.
- `internal/adb/` — adb invocation, device detection, and output parsing.
- `internal/packages/` — parser for `packages.txt`.
- `packages.txt` — the default debloat list (data).
- `docs/adr/` — architecture decision records.
- `.golangci.yaml`, `.goreleaser.yaml`, `.github/workflows/` — lint, release, CI.
- **`backups/`** (gitignored) — a good place for `show -o` dumps so you don't
  commit device-specific lists by mistake.

## Development

```bash
make build      # ./bin/debloat
make test       # go test -race ./...
make check      # gofmt check + go vet + golangci-lint + tests
make run ARGS='uninstall --dry-run'
```

## Contributing

Issues and pull requests are welcome — see [CONTRIBUTING.md](CONTRIBUTING.md)
and the [Code of Conduct](CODE_OF_CONDUCT.md). Most changes can be developed and
tested without a phone using `--dry-run`.

## Safety

- Removing or disabling the wrong package can break **calling, SMS/MMS, updates,
  or security features**. Prefer `list` and a full **`show`** dump before
  changing the device.
- **`reinstall`** only works if the package is still on the system image; keep a
  copy of your package list and know how to recover (e.g. factory reset) if
  something goes wrong.

## License

Released under the [MIT License](LICENSE). The software is provided "as is",
without warranty of any kind. **Use at your own risk — you are responsible for
what you run on your device.**
