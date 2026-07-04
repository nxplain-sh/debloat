#!/usr/bin/env bash
# Debloat Vivo (and similar BBK) devices via ADB — package list from Technastic:
# https://technastic.com/vivo-bloatware-preinstalled-apps-list/
#
# Usage:
#   ./scripts/debloat.sh show [-o file]   # system packages from phone (pm list packages -s)
#   ./scripts/debloat.sh list [-f packages.txt]
#   ./scripts/debloat.sh uninstall [-f packages.txt]
#   ./scripts/debloat.sh freeze [-f packages.txt]      # appops RUN_IN_BACKGROUND ignore
#   ./scripts/debloat.sh disable [-f packages.txt]     # pm disable-user --user 0
#   ./scripts/debloat.sh reinstall [-f packages.txt]   # cmd package install-existing
#
#   uninstall, freeze, disable, reinstall accept -n / --dry-run (print shell commands; no changes).
#
# Requires: adb in PATH, USB debugging, one device (or ANDROID_SERIAL).

set -euo pipefail

_script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
if [[ "$(basename "$_script_dir")" == "scripts" ]]; then
  PROJECT_ROOT="$(cd "$_script_dir/.." && pwd)"
else
  PROJECT_ROOT="$_script_dir"
fi
PACKAGES_FILE="$PROJECT_ROOT/packages.txt"
OUTPUT_FILE=""
DRY_RUN=0
ADB=(adb)

usage() {
  sed -n '1,21p' "$0" | tail -n +2
  exit "${1:-0}"
}

pkg_installed_for_user0() {
  local pkg="$1" haystack="$2"
  grep -Fx "package:${pkg}" <<<"$haystack" >/dev/null
}

classify_failure() {
  local out="$1"
  if echo "$out" | grep -qiE 'not installed|not found|Unknown package|does not exist'; then
    echo skip
  else
    echo fail
  fi
}

# pm often prints "Failure [not installed for 0]" — not a real error; clarify for the user.
tally_cmd_result() {
  local out="$1"
  case "$(classify_failure "$out")" in
    skip)
      echo "  → skipped (not installed for user 0 on this device, or already removed)" >&2
      ((count_skip++)) || true
      ;;
    fail)
      ((count_fail++)) || true
      ;;
  esac
}

ensure_adb_device() {
  if ! command -v adb >/dev/null 2>&1; then
    echo "adb not found. Install Android platform-tools and add to PATH." >&2
    exit 1
  fi

  mapfile -t devices < <(adb devices | awk 'NR>1 && $2=="device" {print $1}')
  if [[ ${#devices[@]} -eq 0 ]]; then
    echo "No authorized device in 'adb devices'. Enable USB debugging and accept the RSA prompt." >&2
    exit 1
  fi
  if [[ ${#devices[@]} -gt 1 && -z "${ANDROID_SERIAL:-}" ]]; then
    echo "Multiple devices connected; set ANDROID_SERIAL to one of:" >&2
    printf '  %s\n' "${devices[@]}" >&2
    exit 1
  fi
}

fetch_system_package_list() {
  local raw
  raw="$("${ADB[@]}" shell pm list packages -s --user 0 2>/dev/null | tr -d '\r')" || raw=""
  if [[ -n "$(echo "$raw" | tr -d '[:space:]')" ]]; then
    printf '%s\n' "$raw"
    return 0
  fi
  "${ADB[@]}" shell pm list packages -s 2>/dev/null | tr -d '\r'
}

if [[ $# -lt 1 ]]; then
  echo "Missing command." >&2
  usage 1
fi

CMD="$1"
shift

while [[ $# -gt 0 ]]; do
  case "$1" in
    -f|--file)
      PACKAGES_FILE="$2"
      shift 2
      ;;
    -o|--output)
      OUTPUT_FILE="$2"
      shift 2
      ;;
    -n|--dry-run)
      DRY_RUN=1
      shift
      ;;
    -h|--help) usage 0 ;;
    *)
      echo "Unknown option: $1" >&2
      usage 1
      ;;
  esac
done

case "$CMD" in
  list|uninstall|freeze|disable|reinstall|show) ;;
  *)
    echo "Unknown command: $CMD" >&2
    echo "Valid: show, list, uninstall, freeze, disable, reinstall" >&2
    exit 1
    ;;
esac

if [[ "$DRY_RUN" -eq 1 ]]; then
  case "$CMD" in
    uninstall|freeze|disable|reinstall) ;;
    *)
      echo "Note: --dry-run only applies to uninstall, freeze, disable, reinstall (ignored for $CMD)." >&2
      ;;
  esac
fi

case "$CMD" in
  list|uninstall|freeze|disable|reinstall)
    if [[ ! -f "$PACKAGES_FILE" ]]; then
      echo "Packages file not found: $PACKAGES_FILE" >&2
      exit 1
    fi
    ;;
esac

# Mutating commands in dry-run only print what would run; no adb required.
if ! [[ "$DRY_RUN" -eq 1 && "$CMD" =~ ^(uninstall|freeze|disable|reinstall)$ ]]; then
  ensure_adb_device
fi

if [[ "$CMD" == "show" ]]; then
  raw="$(fetch_system_package_list || true)"
  if [[ -z "$(echo "$raw" | tr -d '[:space:]')" ]]; then
    echo "Failed to read system packages. Is the device authorized for adb?" >&2
    exit 1
  fi
  # Plain names, one per line (matches packages.txt style)
  plain="$(printf '%s\n' "$raw" | sed '/^$/d' | sed 's/^package://' | LC_ALL=C sort -u)"
  n="$(printf '%s\n' "$plain" | grep -c . || true)"
  if [[ -n "$OUTPUT_FILE" ]]; then
    mkdir -p "$(dirname "$OUTPUT_FILE")"
    printf '%s\n' "$plain" >"$OUTPUT_FILE"
    echo "Wrote $n system package names to $OUTPUT_FILE" >&2
  else
    printf '%s\n' "$plain"
  fi
  exit 0
fi

count_ok=0
count_skip=0
count_fail=0

installed_packages_user0=""
if [[ "$CMD" == "list" ]]; then
  installed_packages_user0="$("${ADB[@]}" shell pm list packages --user 0 2>/dev/null || true)"
fi

# Read the package file on fd 3, not stdin — otherwise `adb shell` inside the loop
# inherits stdin from the file and can consume the rest of packages.txt (only first line runs).
while IFS= read -r line <&3 || [[ -n "$line" ]]; do
  line="${line%%#*}"
  line="${line//[[:space:]]/}"
  [[ -z "$line" ]] && continue

  pkg="$line"

  case "$CMD" in
    list)
      if pkg_installed_for_user0 "$pkg" "$installed_packages_user0"; then
        printf '%s\tinstalled\n' "$pkg"
      else
        printf '%s\tnot installed\n' "$pkg"
      fi
      ;;
    uninstall)
      if [[ "$DRY_RUN" -eq 1 ]]; then
        echo "[dry-run] ${ADB[*]} shell pm uninstall --user 0 $pkg"
        ((count_ok++)) || true
      else
        echo "Uninstalling (user 0): $pkg"
        if out="$("${ADB[@]}" shell pm uninstall --user 0 "$pkg" 2>&1)"; then
          echo "  $out"
          ((count_ok++)) || true
        else
          echo "  $out" >&2
          tally_cmd_result "$out"
        fi
      fi
      ;;
    freeze)
      if [[ "$DRY_RUN" -eq 1 ]]; then
        echo "[dry-run] ${ADB[*]} shell cmd appops set $pkg RUN_IN_BACKGROUND ignore"
        ((count_ok++)) || true
      else
        echo "Freezing (RUN_IN_BACKGROUND=ignore): $pkg"
        if out="$("${ADB[@]}" shell cmd appops set "$pkg" RUN_IN_BACKGROUND ignore 2>&1)"; then
          echo "  $out"
          ((count_ok++)) || true
        else
          echo "  $out" >&2
          tally_cmd_result "$out"
        fi
      fi
      ;;
    disable)
      if [[ "$DRY_RUN" -eq 1 ]]; then
        echo "[dry-run] ${ADB[*]} shell pm disable-user --user 0 $pkg"
        ((count_ok++)) || true
      else
        echo "Disabling (user 0): $pkg"
        if out="$("${ADB[@]}" shell pm disable-user --user 0 "$pkg" 2>&1)"; then
          echo "  $out"
          ((count_ok++)) || true
        else
          echo "  $out" >&2
          tally_cmd_result "$out"
        fi
      fi
      ;;
    reinstall)
      if [[ "$DRY_RUN" -eq 1 ]]; then
        echo "[dry-run] ${ADB[*]} shell cmd package install-existing $pkg"
        ((count_ok++)) || true
      else
        echo "Re-installing (install-existing): $pkg"
        if out="$("${ADB[@]}" shell cmd package install-existing "$pkg" 2>&1)"; then
          echo "  $out"
          ((count_ok++)) || true
        else
          echo "  $out" >&2
          tally_cmd_result "$out"
        fi
      fi
      ;;
  esac
done 3< "$PACKAGES_FILE"

if [[ "$CMD" != "list" ]]; then
  echo "---"
  if [[ "$DRY_RUN" -eq 1 && "$CMD" =~ ^(uninstall|freeze|disable|reinstall)$ ]]; then
    echo "dry-run: $count_ok package(s); no changes made."
  else
    echo "ok: $count_ok  skipped/absent: $count_skip  failed: $count_fail"
  fi
fi
