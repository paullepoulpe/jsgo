// Package pin reads tools/oracle/pin.env, the one place that names the
// pinned Closure Compiler release and its oracle jar, for the Go tools
// (tools/clie2e). The shell scripts source the same file.
package pin

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Path is the pin file, relative to the repository root.
const Path = "tools/oracle/pin.env"

// Pin is the content of the pin file.
type Pin struct {
	// Version is the release tag (CLOSURE_VERSION), e.g. "v20261021".
	Version string
	// Jar is the oracle jar's file name in .cache/ (CLOSURE_JAR).
	Jar string
	// JarSHA256 are the sha256 digests fetch.sh accepts for the jar
	// (CLOSURE_JAR_SHA256).
	JarSHA256 []string
}

// Load reads the pin file, found by walking up from the current directory
// (a test runs in its package's directory, a tool at the repository root).
func Load() (Pin, error) {
	dir, err := os.Getwd()
	if err != nil {
		return Pin{}, err
	}
	for {
		b, err := os.ReadFile(filepath.Join(dir, filepath.FromSlash(Path)))
		if err == nil {
			return Parse(b)
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return Pin{}, fmt.Errorf("%s not found above the current directory", Path)
		}
		dir = parent
	}
}

// Parse parses the KEY=value lines of a pin file; '#' starts a comment
// line, and a value may be double-quoted.
func Parse(b []byte) (Pin, error) {
	vars := map[string]string{}
	sc := bufio.NewScanner(bytes.NewReader(b))
	for n := 1; sc.Scan(); n++ {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		k, v, ok := strings.Cut(line, "=")
		if !ok {
			return Pin{}, fmt.Errorf("%s:%d: not KEY=value: %q", Path, n, line)
		}
		if len(v) >= 2 && v[0] == '"' && v[len(v)-1] == '"' {
			v = v[1 : len(v)-1]
		}
		vars[k] = v
	}
	p := Pin{
		Version:   vars["CLOSURE_VERSION"],
		Jar:       vars["CLOSURE_JAR"],
		JarSHA256: strings.Fields(vars["CLOSURE_JAR_SHA256"]),
	}
	if p.Version == "" || p.Jar == "" || len(p.JarSHA256) == 0 {
		return Pin{}, fmt.Errorf("%s: CLOSURE_VERSION, CLOSURE_JAR and CLOSURE_JAR_SHA256 must all be set", Path)
	}
	return p, nil
}

// MustLoad is Load for tools: it exits on error.
func MustLoad() Pin {
	p, err := Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, "pin:", err)
		os.Exit(1)
	}
	return p
}
