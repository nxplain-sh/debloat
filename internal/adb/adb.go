// Package adb wraps the Android Debug Bridge (adb) command-line tool with the
// small set of operations vivo-debloater needs. adb itself remains an external
// dependency and must be on PATH.
package adb

import (
	"fmt"
	"os/exec"
	"regexp"
	"sort"
	"strings"
)

// DefaultBinary is the adb executable name looked up on PATH.
const DefaultBinary = "adb"

// Client runs adb commands. Create one with New.
type Client struct {
	bin string
}

// New returns a Client that invokes the given adb binary. An empty name means
// the default ("adb").
func New(bin string) *Client {
	if bin == "" {
		bin = DefaultBinary
	}
	return &Client{bin: bin}
}

// Run executes `adb <args...>` and returns its combined stdout+stderr with
// carriage returns and the trailing newline stripped. adb honors the
// ANDROID_SERIAL environment variable to select a device.
func (c *Client) Run(args ...string) (string, error) {
	out, err := exec.Command(c.bin, args...).CombinedOutput()
	s := strings.ReplaceAll(string(out), "\r", "")
	return strings.TrimRight(s, "\n"), err
}

// output runs adb and returns stdout only (stderr discarded), CR-stripped.
func (c *Client) output(args ...string) (string, error) {
	out, err := exec.Command(c.bin, args...).Output()
	return strings.ReplaceAll(string(out), "\r", ""), err
}

// EnsureDevice verifies adb is on PATH and that exactly one usable device is
// selected. serial should be the value of ANDROID_SERIAL (may be empty). With
// several devices connected and no serial set, it returns an error listing
// the available serials.
func (c *Client) EnsureDevice(serial string) error {
	if _, err := exec.LookPath(c.bin); err != nil {
		return fmt.Errorf("adb not found on PATH: install Android platform-tools")
	}
	out, err := c.output("devices")
	if err != nil {
		return fmt.Errorf("running `adb devices`: %w", err)
	}
	devices := parseDevices(out)
	switch {
	case len(devices) == 0:
		return fmt.Errorf("no authorized device in `adb devices`: enable USB debugging and accept the RSA prompt")
	case len(devices) > 1 && serial == "":
		return fmt.Errorf("multiple devices connected; set ANDROID_SERIAL to one of: %s", strings.Join(devices, ", "))
	}
	return nil
}

// parseDevices extracts the serials of authorized ("device") entries from the
// output of `adb devices`. The header line and states such as "unauthorized"
// or "offline" are ignored.
func parseDevices(out string) []string {
	var devices []string
	for _, line := range strings.Split(out, "\n") {
		fields := strings.Fields(line)
		if len(fields) >= 2 && fields[1] == "device" {
			devices = append(devices, fields[0])
		}
	}
	return devices
}

// SystemPackages returns the sorted, de-duplicated system package names on the
// device (`pm list packages -s`), preferring the current user (`--user 0`)
// and falling back when that flag is unsupported.
func (c *Client) SystemPackages() ([]string, error) {
	out, err := c.output("shell", "pm", "list", "packages", "-s", "--user", "0")
	if err != nil || strings.TrimSpace(out) == "" {
		out, err = c.output("shell", "pm", "list", "packages", "-s")
		if err != nil {
			return nil, err
		}
	}
	return normalizePackageList(out), nil
}

// InstalledUser0 returns the set of packages installed for user 0.
func (c *Client) InstalledUser0() (map[string]bool, error) {
	out, err := c.output("shell", "pm", "list", "packages", "--user", "0")
	if err != nil {
		return nil, err
	}
	set := make(map[string]bool)
	for _, p := range normalizePackageList(out) {
		set[p] = true
	}
	return set, nil
}

// normalizePackageList turns raw `pm list packages` output into plain package
// names: it strips the `package:` prefix, drops blanks, then sorts and
// de-duplicates.
func normalizePackageList(out string) []string {
	seen := make(map[string]struct{})
	var pkgs []string
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		line = strings.TrimPrefix(line, "package:")
		if line == "" {
			continue
		}
		if _, ok := seen[line]; ok {
			continue
		}
		seen[line] = struct{}{}
		pkgs = append(pkgs, line)
	}
	sort.Strings(pkgs)
	return pkgs
}

// Result categorizes the outcome of a mutating adb command.
type Result int

const (
	// ResultOK means the command succeeded.
	ResultOK Result = iota
	// ResultSkipped means the package was absent / already removed for user 0.
	ResultSkipped
	// ResultFailed means the command failed for another reason.
	ResultFailed
)

var absentRE = regexp.MustCompile(`(?i)not installed|not found|unknown package|does not exist`)

// Classify inspects the output of a failed mutating command and decides
// whether the failure merely means the package was absent (ResultSkipped) or a
// real error (ResultFailed). Call it only when the command returned an error.
func Classify(out string) Result {
	if absentRE.MatchString(out) {
		return ResultSkipped
	}
	return ResultFailed
}
