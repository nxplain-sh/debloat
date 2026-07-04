# vivo-debloater — Go CLI for debloating Vivo (BBK) phones over adb.

BINARY  := vivo-debloater
CMD     := ./cmd/vivo-debloater
BINDIR  := bin
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -s -w -X github.com/miguelmartens/vivo-debloater/internal/cli.version=$(VERSION)

.DEFAULT_GOAL := help
.PHONY: help build install run test vet fmt fmt-check lint check tidy clean \
        shellcheck prettier prettier-check install-adb

help:
	@echo "vivo-debloater"
	@echo ""
	@echo "  make build          Build $(BINDIR)/$(BINARY)"
	@echo "  make install        go install $(CMD)"
	@echo "  make run ARGS=...    go run $(CMD) with ARGS (e.g. ARGS='uninstall --dry-run')"
	@echo "  make test           go test -race ./..."
	@echo "  make vet            go vet ./..."
	@echo "  make fmt            gofmt -w ."
	@echo "  make fmt-check      fail if gofmt would change files"
	@echo "  make lint           staticcheck ./... (requires staticcheck)"
	@echo "  make check          fmt-check + vet + lint + test"
	@echo "  make tidy           go mod tidy"
	@echo "  make clean          remove build artifacts"
	@echo "  make shellcheck     lint the legacy scripts/debloat.sh"
	@echo "  make prettier       format docs/config with Prettier"
	@echo "  make install-adb    install adb via Homebrew (macOS)"

build:
	@mkdir -p $(BINDIR)
	go build -ldflags '$(LDFLAGS)' -o $(BINDIR)/$(BINARY) $(CMD)

install:
	go install -ldflags '$(LDFLAGS)' $(CMD)

run:
	go run $(CMD) $(ARGS)

test:
	go test -race ./...

vet:
	go vet ./...

fmt:
	gofmt -w .

fmt-check:
	@out="$$(gofmt -l .)"; \
	if [ -n "$$out" ]; then echo "gofmt needed on:"; echo "$$out"; exit 1; fi

lint:
	@command -v staticcheck >/dev/null 2>&1 || { \
		echo "staticcheck not found: go install honnef.co/go/tools/cmd/staticcheck@latest"; exit 1; }
	staticcheck ./...

check: fmt-check vet lint test

tidy:
	go mod tidy

clean:
	rm -rf $(BINDIR) dist

shellcheck:
	@command -v shellcheck >/dev/null 2>&1 || { echo "shellcheck not found"; exit 1; }
	shellcheck scripts/debloat.sh

prettier:
	npx --yes prettier@3 --write . --ignore-unknown

prettier-check:
	npx --yes prettier@3 --check . --ignore-unknown

install-adb:
	@command -v brew >/dev/null 2>&1 || { echo "Homebrew not found. Install from https://brew.sh"; exit 1; }
	brew install --cask android-platform-tools
	@echo "adb should be on your PATH (restart the shell if needed). Verify: adb version"
