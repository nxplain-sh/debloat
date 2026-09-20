# debloat — Go CLI for debloating Vivo (BBK) phones over adb.

BINARY  := debloat
CMD     := ./cmd/debloat
BINDIR  := bin
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -s -w -X github.com/nxplain-sh/debloat/internal/cli.version=$(VERSION)

.DEFAULT_GOAL := help
.PHONY: help build install run test vet fmt fmt-check lint check tidy clean \
        format format-check install-adb hooks

help:
	@echo "debloat"
	@echo ""
	@echo "  make build          Build $(BINDIR)/$(BINARY)"
	@echo "  make install        go install $(CMD)"
	@echo "  make run ARGS=...    go run $(CMD) with ARGS (e.g. ARGS='uninstall --dry-run')"
	@echo "  make test           go test -race ./..."
	@echo "  make vet            go vet ./..."
	@echo "  make fmt            gofmt -w ."
	@echo "  make fmt-check      fail if gofmt would change files"
	@echo "  make lint           golangci-lint run (requires golangci-lint v2)"
	@echo "  make check          fmt-check + vet + lint + test"
	@echo "  make tidy           go mod tidy"
	@echo "  make clean          remove build artifacts"
	@echo "  make format         format docs/config with Prettier"
	@echo "  make format-check   fail if Prettier would change files"
	@echo "  make install-adb    install adb via Homebrew (macOS)"
	@echo "  make hooks          enable the repo's git pre-commit hook"

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
	@command -v golangci-lint >/dev/null 2>&1 || { \
		echo "golangci-lint not found: https://golangci-lint.run/welcome/install/"; exit 1; }
	golangci-lint run ./...

check: fmt-check vet lint test

tidy:
	go mod tidy

clean:
	rm -rf $(BINDIR) dist

format:
	npx --yes prettier --write . --ignore-unknown

format-check:
	npx --yes prettier --check . --ignore-unknown

hooks:
	git config core.hooksPath .githooks
	@echo "pre-commit hook enabled (fmt-check + vet + test)"

install-adb:
	@command -v brew >/dev/null 2>&1 || { echo "Homebrew not found. Install from https://brew.sh"; exit 1; }
	brew install --cask android-platform-tools
	@echo "adb should be on your PATH (restart the shell if needed). Verify: adb version"
