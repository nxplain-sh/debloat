// Command vivo-debloater lists and removes preinstalled apps on Vivo (BBK)
// phones over adb, without root. It is a thin wrapper: the real work lives in
// internal packages.
package main

import (
	"os"

	"github.com/miguelmartens/vivo-debloater/internal/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:], os.Stdout, os.Stderr))
}
