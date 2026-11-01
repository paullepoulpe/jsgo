// Command ast_isomorphism verifies lockstep AST + lexical-scope alpha-isomorphism
// between two compiled JS files (e.g. compat vs. fast_mode output) using jsgo's
// parser and SyntacticScopeCreator.
//
//	ast_isomorphism <compat.js> <fast.js>
//	    whole-program check (fast_mode_gate.py layer 2): exit 0 when isomorphic
//	ast_isomorphism -diff [-props exact|renamed] [-fixed names.txt] <compat.js> <fast.js>
//	    where the builds differ, per function (diff.go; the behavioural
//	    oracle's coverage rule): JSON on stdout, exit 0 unless an input fails
//	    to parse
//	ast_isomorphism -extern-names
//	    the identifiers of the default externs, one per line (the property
//	    names that ADVANCED renaming leaves alone; the oracle passes them back
//	    as -fixed)
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/cli"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type output struct {
	Isomorphic bool                 `json:"isomorphic"`
	Error      string               `json:"error,omitempty"`
	Stats      *AstIsomorphismStats `json:"stats,omitempty"`
}

func main() {
	diff := flag.Bool("diff", false, "report where the builds differ, per function")
	props := flag.String("props", "exact", "with -diff: exact (property names compared as they are) or renamed (up to a consistent renaming)")
	fixed := flag.String("fixed", "", "with -diff -props renamed: a file of names (one per line) that keep their spelling")
	externNames := flag.Bool("extern-names", false, "print the identifiers of the default externs")
	flag.Parse()
	if *externNames {
		printExternNames()
		return
	}
	if flag.NArg() != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s [-diff [-props exact|renamed] [-fixed names.txt]] <compat.js> <fast.js>\n", os.Args[0])
		os.Exit(2)
	}
	bA, err := os.ReadFile(flag.Arg(0))
	if err != nil {
		fmt.Fprintf(os.Stderr, "read %s: %v\n", flag.Arg(0), err)
		os.Exit(2)
	}
	bB, err := os.ReadFile(flag.Arg(1))
	if err != nil {
		fmt.Fprintf(os.Stderr, "read %s: %v\n", flag.Arg(1), err)
		os.Exit(2)
	}
	enc := json.NewEncoder(os.Stdout)

	if *diff {
		opts := DiffOptions{RenamedProps: *props == "renamed", Fixed: map[string]bool{}}
		if *props != "exact" && *props != "renamed" {
			fmt.Fprintf(os.Stderr, "-props: exact or renamed, not %q\n", *props)
			os.Exit(2)
		}
		if *fixed != "" {
			f, err := os.Open(*fixed)
			if err != nil {
				fmt.Fprintf(os.Stderr, "read %s: %v\n", *fixed, err)
				os.Exit(2)
			}
			sc := bufio.NewScanner(f)
			sc.Buffer(make([]byte, 1<<20), 1<<20)
			for sc.Scan() {
				if w := strings.TrimSpace(sc.Text()); w != "" {
					opts.Fixed[w] = true
				}
			}
			f.Close()
		}
		res, err := DiffBuilds(string(bA), string(bB), opts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "diff: %v\n", err)
			os.Exit(1)
		}
		_ = enc.Encode(res)
		return
	}

	stats, isoErr := VerifyAstAlphaIsomorphism(string(bA), string(bB))
	res := output{
		Isomorphic: isoErr == nil,
		Stats:      stats,
	}
	if isoErr != nil {
		res.Error = isoErr.Error()
	}
	_ = enc.Encode(res)
	if isoErr != nil {
		os.Exit(1)
	}
}

var identRE = regexp.MustCompile(`[A-Za-z_$][A-Za-z0-9_$]*`)
var typeTagRE = regexp.MustCompile(`@[A-Za-z]+\s*\{`)

// printExternNames prints the names that ADVANCED's property renaming takes
// from the default externs: the property names of their code, and the
// identifiers inside their JSDoc type expressions (record types name
// properties there). Words of the prose comments are left out, since short
// renamed names (a, b, ...) occur there too.
func printExternNames() {
	seen := map[string]bool{}
	for _, f := range cli.CommandLineRunner_GetDefaultExterns() {
		code := f.GetCode()
		root, _, err := parseForIsomorphism(f.GetName(), code)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", f.GetName(), err)
			os.Exit(1)
		}
		collectPropNames(root, seen)
		for _, loc := range typeTagRE.FindAllStringIndex(code, -1) {
			depth, i := 0, loc[1]-1
			for ; i < len(code); i++ {
				if code[i] == '{' {
					depth++
				} else if code[i] == '}' {
					depth--
					if depth == 0 {
						break
					}
				}
			}
			if i > len(code) {
				i = len(code)
			}
			for _, w := range identRE.FindAllString(code[loc[1]:i], -1) {
				seen[w] = true
			}
		}
	}
	names := make([]string, 0, len(seen))
	// order-independent: sorted below.
	for w := range seen {
		names = append(names, w)
	}
	sort.Strings(names)
	w := bufio.NewWriter(os.Stdout)
	for _, n := range names {
		fmt.Fprintln(w, n)
	}
	w.Flush()
}

func collectPropNames(n rhino.Node, seen map[string]bool) {
	if isPropToken(n.GetToken()) {
		seen[n.GetString()] = true
	}
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		collectPropNames(c, seen)
	}
}
