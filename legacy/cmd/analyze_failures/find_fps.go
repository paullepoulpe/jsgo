//go:build ignore

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/paullepoulpe/jsgo/legacy/internal/diag"
	"github.com/paullepoulpe/jsgo/legacy/internal/jscomp"
)

type tuple3 struct {
	path string
	line int
	col  int
	code string
}

var diagRe3 = regexp.MustCompile(`^(.+?):(\d+):(\d+): (?:ERROR|WARNING) - \[([A-Z_]+)\]`)

func parseTuples3(text string) []tuple3 {
	var out []tuple3
	for _, line := range strings.Split(text, "\n") {
		m := diagRe3.FindStringSubmatch(line)
		if m == nil { continue }
		ln, _ := strconv.Atoi(m[2])
		cl, _ := strconv.Atoi(m[3])
		out = append(out, tuple3{path: m[1], line: ln, col: cl, code: m[4]})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].line != out[j].line { return out[i].line < out[j].line }
		if out[i].col != out[j].col { return out[i].col < out[j].col }
		return out[i].code < out[j].code
	})
	return out
}

func parseOracleExit3(text string) int {
	lines := strings.Split(strings.TrimRight(text, "\n"), "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		if lines[i] == "---EXIT---" && i+1 < len(lines) {
			v, _ := strconv.Atoi(strings.TrimSpace(lines[i+1]))
			return v
		}
	}
	return -1
}

func main() {
	root := "testdata/fixtures/extracted"
	count := 0
	_ = filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(p, ".js") { return nil }
		if count >= 20 { return nil }

		fixRel := p
		oraclePath := filepath.Join("testdata/oracle",
			strings.TrimSuffix(fixRel[len("testdata/fixtures/"):], ".js")+".txt")
		oracleBytes, err := os.ReadFile(oraclePath)
		if err != nil { return nil }

		expected := parseTuples3(string(oracleBytes))
		wantExit := parseOracleExit3(string(oracleBytes))

		src, _ := os.ReadFile(fixRel)
		ds, _ := jscomp.CheckSourceWithOptions(fixRel, string(src), jscomp.NewCompilerOptions())
		var actual []tuple3
		gotExit := 0
		for _, d := range ds {
			actual = append(actual, tuple3{path: d.Path, line: d.Line, col: d.Col, code: d.Code})
			if d.Level == diag.Error { gotExit = 1 }
		}
		sort.Slice(actual, func(i, j int) bool {
			if actual[i].line != actual[j].line { return actual[i].line < actual[j].line }
			if actual[i].col != actual[j].col { return actual[i].col < actual[j].col }
			return actual[i].code < actual[j].code
		})

		// Only report false positives
		hasFP := false
		for _, a := range actual {
			found := false
			for _, e := range expected { if a == e { found = true; break } }
			if !found { hasFP = true; break }
		}
		if !hasFP { return nil }

		exitMatch := (wantExit == 0) == (gotExit == 0)
		match := len(expected) == len(actual)
		if match {
			for i := range expected {
				if expected[i] != actual[i] { match = false; break }
			}
		}
		if match && exitMatch { return nil }

		count++
		fmt.Printf("FP %s\n", fixRel)
		for _, a := range actual {
			found := false
			for _, e := range expected { if a == e { found = true; break } }
			if !found {
				fmt.Printf("  EXTRA: %v\n", a)
			}
		}
		for _, e := range expected {
			found := false
			for _, a := range actual { if a == e { found = true; break } }
			if !found {
				fmt.Printf("  MISS:  %v\n", e)
			}
		}
		// Show JS content
		src2, _ := os.ReadFile(fixRel)
		fmt.Printf("  JS: %s\n", strings.TrimSpace(string(src2)))
		return nil
	})
}
