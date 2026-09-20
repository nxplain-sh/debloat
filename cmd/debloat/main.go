// Command debloat lists and removes preinstalled apps on Vivo (BBK)
// phones over adb, without root. It is a thin wrapper: the real work lives in
// internal packages.
package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/nxplain-sh/debloat/internal/cli"
)

func main() {
	os.Exit(run())
}

// run wires up signal-based cancellation and dispatches to the CLI, returning
// the process exit code. It exists so `defer stop()` runs before os.Exit.
func run() int {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	return cli.Run(ctx, os.Args[1:], os.Stdout, os.Stderr)
}
