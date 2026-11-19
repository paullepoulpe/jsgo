// Command gopasses runs modular, independently toggleable Go-to-Go AST + go/types
// post-step optimization passes over mechanically transpiled Go packages.
//
// Usage:
//
//	gopasses [flags] [packages / directories / files...]
//
// Examples:
//
//	gopasses -passes=all -stats ./transpiled/gen/...
//	gopasses -passes=dead_cast_cleanup,j2gort_peephole -stats ./transpiled/gen/...
//	gopasses -passes=none -stats ./transpiled/gen/...
//	gopasses -passes=all -dry-run -stats ./transpiled/gen/...
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/paullepoulpe/jsgo/transpiled/translator/passes"
)

func main() {
	os.Exit(runMain(os.Args[1:], os.Stdout, os.Stderr))
}

func runMain(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("gopasses", flag.ContinueOnError)
	fs.SetOutput(stderr)

	passesSpec := fs.String("passes", "all",
		fmt.Sprintf("comma-separated passes to run (%s), 'all', or 'none'", strings.Join(passes.PassNames(), ",")))
	printStats := fs.Bool("stats", false, "print exact counts of transformed sites per pass")
	jsonStats := fs.Bool("json-stats", false, "print pass statistics as JSON")
	dryRun := fs.Bool("dry-run", false, "analyze, verify, and report stats without modifying files on disk")
	writeFiles := fs.Bool("w", true, "write transformed files in-place (set -w=false for dry-run)")
	includeTests := fs.Bool("tests", false, "include in-package *_test.go files")
	listPasses := fs.Bool("list", false, "list available passes and exit")

	fs.Usage = func() {
		fmt.Fprintf(stderr, "Usage: gopasses [flags] [packages / directories / files...]\n\n")
		fmt.Fprintf(stderr, "Available passes:\n")
		for _, p := range passes.AllPasses() {
			fmt.Fprintf(stderr, "  %-20s %s\n", p.Name(), p.Description())
		}
		fmt.Fprintf(stderr, "\nFlags:\n")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 2
	}

	if *listPasses {
		for _, p := range passes.AllPasses() {
			fmt.Fprintf(stdout, "%-20s %s\n", p.Name(), p.Description())
		}
		return 0
	}

	selected, err := passes.ParsePassesFlag(*passesSpec)
	if err != nil {
		fmt.Fprintf(stderr, "gopasses: %v\n", err)
		return 2
	}

	patterns := fs.Args()
	if len(patterns) == 0 {
		patterns = []string{"."}
	}

	opts := passes.Options{
		DryRun:       *dryRun || !*writeFiles,
		IncludeTests: *includeTests,
	}

	stats, err := passes.Run(patterns, selected, opts)
	if err != nil {
		fmt.Fprintf(stderr, "gopasses: %v\n", err)
		return 1
	}

	if *jsonStats {
		if err := stats.WriteJSON(stdout); err != nil {
			fmt.Fprintf(stderr, "gopasses: %v\n", err)
			return 1
		}
	} else if *printStats {
		if err := stats.WriteText(stdout); err != nil {
			fmt.Fprintf(stderr, "gopasses: %v\n", err)
			return 1
		}
	}

	return 0
}
