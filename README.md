# vivo-debloater

Small bash helpers to list system packages on a Vivo (BBK) phone running **OriginOS** (older models may still show **Funtouch** in settings) and remove or restrict apps **without root**, using `adb`. The default package list is based on [Technastic — Vivo bloatware list](https://technastic.com/vivo-bloatware-preinstalled-apps-list/).

## Prerequisites

- **`adb`** on your `PATH` — e.g. install with **Homebrew**: `make install-adb` (installs the [android-platform-tools](https://developer.android.com/tools/releases/platform-tools) cask), or install manually from Google’s [platform-tools](https://developer.android.com/tools/releases/platform-tools) page
- On the phone: **Developer options** → **USB debugging** enabled; accept the computer’s RSA fingerprint when prompted
- One phone connected over USB (or Wi‑Fi debugging with `adb connect`). If several devices show up in `adb devices`, set **`ANDROID_SERIAL`** to the serial you want

## Layout

- `scripts/debloat.sh` — main script (defaults to `packages.txt` in the **repo root**)
- `packages.txt` — package list to list / uninstall / freeze / disable / reinstall
- `Makefile` — optional wrappers (`make show`, `make list`, …)
- **`backups/`** (optional) — a good place for `show -o` dumps: keeps the repo root tidy and stays **gitignored** so you do not commit device-specific lists by mistake. **Do not use a temp directory** for these files; you want to keep them for diffing and recovery. The script creates missing parent directories for `-o` paths (for example `backups/phone.txt`).

## Quick start

```bash
make show ARGS='-o backups/my-phone-system.txt'   # snapshot what the ROM ships
# Or: ./scripts/debloat.sh show -o backups/my-phone-system.txt

# Edit packages.txt: comment out (#) anything you want to keep
make list                                   # check install status vs packages.txt
make uninstall ARGS='--dry-run'             # preview adb commands (no changes)
make uninstall                                # apply — or: make freeze | make disable
```

## Commands

| Command     | Purpose                                                                                                                          |
| ----------- | -------------------------------------------------------------------------------------------------------------------------------- |
| `show`      | Print **system** package names from the device (`pm list packages -s`, with `--user 0` when supported). One plain name per line. |
| `list`      | For each entry in the package list, prints `installed` or `not installed` for the current user.                                  |
| `uninstall` | `pm uninstall --user 0 <package>` — removes the app **for user 0** (system APK usually stays on the device).                     |
| `freeze`    | `cmd appops set <package> RUN_IN_BACKGROUND ignore` — restricts background behavior (lighter than uninstall).                    |
| `disable`   | `pm disable-user --user 0 <package>`.                                                                                            |
| `reinstall` | `cmd package install-existing <package>` — brings back a system app for the user when the ROM still has it.                      |

The command is named **`show`** (not `check`) because **`list`** already _checks_ each entry in `packages.txt` against the device; **`show`** is only for dumping the phone’s full system package list.

### Options

- **`-f` / `--file`** — package list file (default: **`packages.txt` in the repository root**). Used by `list`, `uninstall`, `freeze`, `disable`, and `reinstall`.
- **`-o` / `--output`** — output file for **`show`** only. If omitted, the list is printed to stdout. Parent directories are created if missing.
- **`-n` / `--dry-run`** — for **`uninstall`**, **`freeze`**, **`disable`**, and **`reinstall`** only: print the `adb` shell commands that would run; **no device access** and no changes. Ignored for `show` and `list` (a note is printed if you pass it there).

### Workflow

Typical order: snapshot the ROM, curate the list, verify, preview, then apply (or recover).

1. **Save a system package dump** — full list from the phone; use `backups/` so it stays out of git (see Layout).

   ```bash
   make show ARGS='-o backups/vivo-x200-system.txt'
   ```

2. **Edit `packages.txt`** — comment out (`#`) every package you want to keep. Optionally compare lines with your dump so you only target packages that exist on your ROM.

3. **Check install status** — for each uncommented line, see `installed` or `not installed` for user 0:

   ```bash
   make list
   ```

4. **Preview changes** — prints the `adb` shell commands only; does not touch the device:

   ```bash
   make uninstall ARGS='--dry-run'
   ```

   Use the same pattern for lighter actions: `make freeze ARGS='--dry-run'` or `make disable ARGS='--dry-run'`.

5. **Apply** — run the same command without `--dry-run` (pick **one** of uninstall, freeze, or disable depending on how aggressive you want to be):

   ```bash
   make uninstall
   ```

   Or: `make freeze` / `make disable`.

6. **Recover** — if you removed or disabled something you still need, and the system image still has the package:

   ```bash
   make reinstall
   ```

**Multiple devices** — set the serial when more than one device is connected:

```bash
ANDROID_SERIAL=ABC123 make list
```

**Calling the script directly** (equivalent to the `make` targets above):

```bash
./scripts/debloat.sh show -o backups/vivo-x200-system.txt
./scripts/debloat.sh list -f ./packages.txt
./scripts/debloat.sh uninstall --dry-run -f ./packages.txt
./scripts/debloat.sh uninstall -f ./packages.txt
./scripts/debloat.sh freeze -f ./packages.txt
./scripts/debloat.sh disable -f ./packages.txt
./scripts/debloat.sh reinstall -f ./packages.txt
```

**Lint** (optional):

```bash
make shellcheck
```

## `packages.txt`

- One **package name** per line (e.g. `com.vivo.browser`).
- Lines starting with **`#`** are comments; blank lines are ignored.
- Start from the bundled list, compare with your own `show` export, and **comment out** anything you rely on (dialer, MMS, cloud backup you use, etc.).

## Safety

- Removing or disabling the wrong package can break **calling, SMS/MMS, updates, or security features**. Prefer `list` and a full **`show`** dump before changing the device.
- **`reinstall`** only works if the package is still on the system image; keep a copy of your package list and know how to recover (e.g. factory reset) if something goes wrong.

## Contributing

Issues and pull requests are welcome — see [CONTRIBUTING.md](CONTRIBUTING.md)
and the [Code of Conduct](CODE_OF_CONDUCT.md). Most changes can be developed and
tested without a phone using `--dry-run`.

## License

Released under the [MIT License](LICENSE). The software is provided "as is",
without warranty of any kind. **Use at your own risk — you are responsible for
what you run on your device.**
