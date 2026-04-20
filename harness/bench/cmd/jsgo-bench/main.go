// jsgo-bench is a process-level helper used only by the type-check
// benchmark harness (harness/bench). It consumes a newline-delimited
// list of .js file paths from stdin and runs jsgo's type checker on
// each file in a single process, printing Java-shaped diagnostic
// lines to stdout. This mirrors how `java -jar closure.jar --js a.js
// --js b.js ...` batches many files per JVM start, giving a fair
// wall-clock comparison against Closure's cold-start cost.
//
// NOTE: This binary is bench-infra only; it is not part of the MVP
// CLI surface (see SPEC.md). It lives under harness/ which SPEC's
// guardrail #1 explicitly exempts from the pure-Go runtime rule's
// stricter reading (os/exec, net, etc). This binary itself stays
// pure-Go anyway.
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/diag"
	"github.com/paullepoulpe/jsgo/internal/jscomp"
)

func main() {
	os.Exit(run(os.Stdin, os.Stdout, os.Stderr, os.Args[1:]))
}

func run(stdin io.Reader, stdout, stderr io.Writer, args []string) int {
	// Flags are intentionally minimal. The harness passes every file
	// on stdin; we don't parse per-file CLI flags because we want to
	// amortise one Go-process startup across the whole corpus.
	strict := false
	for _, a := range args {
		switch a {
		case "--strict":
			strict = true
		case "-h", "--help":
			fmt.Fprintln(stderr, "usage: jsgo-bench [--strict] < filelist")
			return 0
		default:
			fmt.Fprintln(stderr, "unknown arg:", a)
			return 2
		}
	}

	paths := readLines(stdin)
	exit := 0
	optsFn := func() *jscomp.CompilerOptions { return jscomp.NewCompilerOptions() }
	if strict {
		optsFn = func() *jscomp.CompilerOptions {
			o := jscomp.NewCompilerOptions()
			o.SetWarningLevel(jscomp.GroupStrictCheckTypes, jscomp.LevelWarning)
			return o
		}
	}
	for _, p := range paths {
		src, err := os.ReadFile(p)
		if err != nil {
			fmt.Fprintln(stderr, err)
			exit = 2
			continue
		}
		diags, err := jscomp.CheckSourceWithOptions(p, string(src), optsFn())
		if err != nil {
			// Parse failure. Java would log a warning and keep going; we
			// match that shape and keep going too.
			fmt.Fprintf(stderr, "%s: parse error: %v\n", p, err)
			continue
		}
		for _, d := range diags {
			fmt.Fprintln(stdout, d.FormatJavaLike())
			if d.Level == diag.Error {
				exit = 1
			}
		}
	}
	return exit
}

func readLines(r io.Reader) []string {
	var out []string
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 64*1024), 16*1024*1024)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		out = append(out, line)
	}
	return out
}
