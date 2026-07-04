// Package packages parses the vivo-debloater package-list file: one Android
// package name per line, with `#` comments and blank lines ignored.
package packages

import (
	"bufio"
	"io"
	"os"
	"strings"
)

// Parse reads a package list from r and returns the package names in file
// order. A `#` begins a comment that runs to the end of the line; blank lines
// and surrounding whitespace are ignored. Duplicates are preserved for the
// caller to handle.
func Parse(r io.Reader) ([]string, error) {
	var pkgs []string
	sc := bufio.NewScanner(r)
	// Allow long lines without the default 64 KiB token limit tripping us up.
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for sc.Scan() {
		line := sc.Text()
		if i := strings.IndexByte(line, '#'); i >= 0 {
			line = line[:i]
		}
		// Package names carry no whitespace; drop it all (matches the shell
		// version's `${line//[[:space:]]/}`).
		line = strings.Join(strings.Fields(line), "")
		if line == "" {
			continue
		}
		pkgs = append(pkgs, line)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return pkgs, nil
}

// ParseFile reads and parses the package list at path.
func ParseFile(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Parse(f)
}
