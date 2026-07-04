# vivo-debloater — package list + adb helpers
# Logic lives in scripts/; Makefile is a thin facade.

SHELL := /bin/bash
SCRIPT := scripts/debloat.sh
.DEFAULT_GOAL := help

.PHONY: help install-adb install-prettier prettier prettier-check shellcheck show list uninstall freeze disable reinstall reinstall-one

help:
	@echo "vivo-debloater"
	@echo ""
	@echo "  make install-adb       Install adb via Homebrew (android-platform-tools cask)"
	@echo "  make install-prettier  Install Prettier via Homebrew"
	@echo "  make prettier          Format with Prettier"
	@echo "  make prettier-check    Check formatting (CI)"
	@echo "  make shellcheck        Lint $(SCRIPT) (requires shellcheck on PATH)"
	@echo "  make show              Dump system packages from device (needs adb + phone)"
	@echo "  make list              Compare packages.txt to device"
	@echo "  make uninstall | freeze | disable | reinstall"
	@echo "  make reinstall-one PKG=com.example.app   # single package (install-existing)"
	@echo "  (add ARGS='--dry-run' or ARGS='-n' to print adb commands only; no device needed)"
	@echo ""
	@echo "  make show ARGS='-o my-phone-system.txt'"
	@echo "  make list ARGS='-f packages.txt'"
	@echo ""
	@echo "Or: bash $(SCRIPT) <command> [...]"
	@echo ""

install-adb:
	@command -v brew >/dev/null 2>&1 || { echo "Homebrew not found. Install from https://brew.sh"; exit 1; }
	brew install --cask android-platform-tools
	@echo "adb should be on your PATH (restart the shell if needed). Verify: adb version"

install-prettier:
	@command -v brew >/dev/null 2>&1 || { echo "Homebrew not found. Install from https://brew.sh"; exit 1; }
	brew install prettier
	@echo "Verify: prettier --version"

prettier:
	@command -v prettier >/dev/null 2>&1 || { echo "prettier not found. Run: make install-prettier"; exit 1; }
	@prettier --write . --ignore-unknown $(ARGS)

prettier-check:
	@command -v prettier >/dev/null 2>&1 || { echo "prettier not found. Run: make install-prettier"; exit 1; }
	@prettier --check . --ignore-unknown $(ARGS)

shellcheck:
	@command -v shellcheck >/dev/null 2>&1 || { echo "shellcheck not found"; exit 1; }
	shellcheck $(SCRIPT)

show list uninstall freeze disable reinstall:
	@bash $(SCRIPT) $(@) $(ARGS)

# PKG= is required — same as: adb shell cmd package install-existing <PKG>
reinstall-one:
ifndef PKG
	$(error PKG is required, e.g. make reinstall-one PKG=com.vivo.weather)
endif
	@if [ -n "$(filter --dry-run -n,$(ARGS))" ]; then \
		echo "[dry-run] adb shell cmd package install-existing $(PKG)"; \
	else \
		adb shell cmd package install-existing "$(PKG)"; \
	fi
