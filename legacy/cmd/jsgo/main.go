package main

import (
	"fmt"
	"io"
	"os"

	"github.com/paullepoulpe/jsgo/legacy/internal/jscomp"
)

func main() {
	code := run(os.Args[1:], os.Stdout, os.Stderr)
	os.Exit(code)
}

func run(args []string, stdout, stderr io.Writer) int {
	if len(args) < 2 || args[0] != "typecheck" {
		fmt.Fprintln(stderr, "usage: jsgo typecheck <file.js>")
		return 2
	}
	path := args[1]
	src, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	diags, err := jscomp.CheckSource(path, string(src))
	if err != nil {
		// Parse errors are not part of the Tier-0 conformance check; we
		// surface them to stderr without a Closure error code and exit
		// non-zero. Guardrail #6 (novel-input sanity) only requires no
		// panic.
		fmt.Fprintln(stderr, err)
		return 2
	}
	exit := 0
	for _, d := range diags {
		fmt.Fprintln(stdout, d.FormatJavaLike())
		exit = 1
	}
	return exit
}
