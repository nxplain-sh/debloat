// Package cli implements the vivo-debloater command-line interface. Run is the
// single entry point invoked from cmd/vivo-debloater.
package cli

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/miguelmartens/vivo-debloater/internal/adb"
	"github.com/miguelmartens/vivo-debloater/internal/packages"
)

// version is overridden at release time via -ldflags
// "-X .../internal/cli.version=<tag>".
var version = "dev"

// defaultPackagesFile is the package list used when -f/--file is not given.
// It is resolved relative to the current working directory.
const defaultPackagesFile = "packages.txt"

// mutation describes a mutating command's human-readable verb and the adb
// arguments it runs for a single package.
type mutation struct {
	verb string
	args func(pkg string) []string
}

var mutations = map[string]mutation{
	"uninstall": {
		verb: "Uninstalling (user 0)",
		args: func(p string) []string { return []string{"shell", "pm", "uninstall", "--user", "0", p} },
	},
	"freeze": {
		verb: "Freezing (RUN_IN_BACKGROUND=ignore)",
		args: func(p string) []string {
			return []string{"shell", "cmd", "appops", "set", p, "RUN_IN_BACKGROUND", "ignore"}
		},
	},
	"disable": {
		verb: "Disabling (user 0)",
		args: func(p string) []string { return []string{"shell", "pm", "disable-user", "--user", "0", p} },
	},
	"reinstall": {
		verb: "Re-installing (install-existing)",
		args: func(p string) []string { return []string{"shell", "cmd", "package", "install-existing", p} },
	},
}

// Run parses args (the program arguments, excluding the program name) and
// executes the selected command, returning a process exit code. The context
// cancels in-flight adb calls (e.g. on Ctrl-C).
func Run(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "missing command")
		usage(stderr)
		return 2
	}
	cmd, rest := args[0], args[1:]
	switch cmd {
	case "help", "-h", "--help":
		usage(stdout)
		return 0
	case "version", "-v", "--version":
		fmt.Fprintf(stdout, "vivo-debloater %s\n", version)
		return 0
	case "show":
		return runShow(ctx, rest, stdout, stderr)
	case "list":
		return runList(ctx, rest, stdout, stderr)
	case "uninstall", "freeze", "disable", "reinstall":
		return runMutate(ctx, cmd, rest, stdout, stderr)
	default:
		fmt.Fprintf(stderr, "unknown command: %s\n", cmd)
		fmt.Fprintln(stderr, "valid: show, list, uninstall, freeze, disable, reinstall")
		return 2
	}
}

func runShow(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	fs := newFlagSet("show", stderr)
	var output string
	stringFlag(fs, &output, "", "write the package list to FILE instead of stdout", "o", "output")
	if err := fs.Parse(args); err != nil {
		return 2
	}

	client := adb.New("")
	if err := client.EnsureDevice(ctx, os.Getenv("ANDROID_SERIAL")); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	pkgs, err := client.SystemPackages(ctx)
	if err != nil || len(pkgs) == 0 {
		fmt.Fprintln(stderr, "failed to read system packages: is the device authorized for adb?")
		return 1
	}

	body := strings.Join(pkgs, "\n") + "\n"
	if output == "" {
		fmt.Fprint(stdout, body)
		return 0
	}
	if dir := filepath.Dir(output); dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			fmt.Fprintf(stderr, "creating %s: %v\n", dir, err)
			return 1
		}
	}
	if err := os.WriteFile(output, []byte(body), 0o644); err != nil {
		fmt.Fprintf(stderr, "writing %s: %v\n", output, err)
		return 1
	}
	fmt.Fprintf(stderr, "Wrote %d system package names to %s\n", len(pkgs), output)
	return 0
}

func runList(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	fs := newFlagSet("list", stderr)
	var file string
	stringFlag(fs, &file, defaultPackagesFile, "package list file", "f", "file")
	if err := fs.Parse(args); err != nil {
		return 2
	}

	pkgs, err := packages.ParseFile(file)
	if err != nil {
		fmt.Fprintf(stderr, "packages file: %v\n", err)
		return 1
	}

	client := adb.New("")
	if err := client.EnsureDevice(ctx, os.Getenv("ANDROID_SERIAL")); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	installed, err := client.InstalledUser0(ctx)
	if err != nil {
		fmt.Fprintf(stderr, "listing installed packages: %v\n", err)
		return 1
	}

	for _, p := range pkgs {
		status := "not installed"
		if installed[p] {
			status = "installed"
		}
		fmt.Fprintf(stdout, "%s\t%s\n", p, status)
	}
	return 0
}

func runMutate(ctx context.Context, cmd string, args []string, stdout, stderr io.Writer) int {
	fs := newFlagSet(cmd, stderr)
	var (
		file   string
		dryRun bool
	)
	stringFlag(fs, &file, defaultPackagesFile, "package list file", "f", "file")
	boolFlag(fs, &dryRun, "print the adb commands without running them", "n", "dry-run")
	if err := fs.Parse(args); err != nil {
		return 2
	}

	pkgs, err := packages.ParseFile(file)
	if err != nil {
		fmt.Fprintf(stderr, "packages file: %v\n", err)
		return 1
	}

	m := mutations[cmd]
	client := adb.New("")
	if !dryRun {
		if err := client.EnsureDevice(ctx, os.Getenv("ANDROID_SERIAL")); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
	}

	var ok, skipped, failed int
	for _, p := range pkgs {
		cmdArgs := m.args(p)
		if dryRun {
			fmt.Fprintf(stdout, "[dry-run] %s %s\n", adb.DefaultBinary, strings.Join(cmdArgs, " "))
			ok++
			continue
		}

		fmt.Fprintf(stdout, "%s: %s\n", m.verb, p)
		out, runErr := client.Run(ctx, cmdArgs...)
		if out != "" {
			fmt.Fprintf(stdout, "  %s\n", strings.ReplaceAll(out, "\n", "\n  "))
		}
		if runErr == nil {
			ok++
			continue
		}
		if adb.Classify(out) == adb.ResultSkipped {
			fmt.Fprintln(stderr, "  → skipped (not installed for user 0 on this device, or already removed)")
			skipped++
		} else {
			failed++
		}
	}

	fmt.Fprintln(stdout, "---")
	if dryRun {
		fmt.Fprintf(stdout, "dry-run: %d package(s); no changes made.\n", ok)
	} else {
		fmt.Fprintf(stdout, "ok: %d  skipped/absent: %d  failed: %d\n", ok, skipped, failed)
	}
	if failed > 0 {
		return 1
	}
	return 0
}

// newFlagSet builds a FlagSet that reports errors to w and suppresses the
// default usage dump (Run prints its own help).
func newFlagSet(name string, w io.Writer) *flag.FlagSet {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(w)
	return fs
}

// stringFlag registers a string flag under several names (e.g. short and long
// forms) bound to the same variable.
func stringFlag(fs *flag.FlagSet, p *string, def, usage string, names ...string) {
	for _, n := range names {
		fs.StringVar(p, n, def, usage)
	}
}

// boolFlag registers a bool flag under several names bound to the same variable.
func boolFlag(fs *flag.FlagSet, p *bool, usage string, names ...string) {
	for _, n := range names {
		fs.BoolVar(p, n, false, usage)
	}
}

func usage(w io.Writer) {
	fmt.Fprint(w, `vivo-debloater — debloat Vivo (BBK) phones over adb, without root

Usage:
  vivo-debloater <command> [flags]

Commands:
  show        Print system package names from the device (pm list packages -s)
  list        Show installed / not installed for each entry in the package list
  uninstall   pm uninstall --user 0 <pkg>
  freeze      cmd appops set <pkg> RUN_IN_BACKGROUND ignore
  disable     pm disable-user --user 0 <pkg>
  reinstall   cmd package install-existing <pkg>
  version     Print the version
  help        Show this help

Flags:
  -f, --file FILE     package list (default "packages.txt");
                      used by list, uninstall, freeze, disable, reinstall
  -o, --output FILE   write output to FILE instead of stdout (show only)
  -n, --dry-run       print adb commands without running them
                      (uninstall, freeze, disable, reinstall only)

Environment:
  ANDROID_SERIAL      target device when several are connected

Examples:
  vivo-debloater show -o backups/phone.txt
  vivo-debloater list -f packages.txt
  vivo-debloater uninstall --dry-run
  vivo-debloater uninstall
`)
}
