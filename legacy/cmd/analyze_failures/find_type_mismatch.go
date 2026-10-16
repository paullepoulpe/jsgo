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

	"github.com/paullepoulpe/jsgo/legacy/internal/jscomp"
)

type tuple struct {
	path string
	line int
	col  int
	code string
}

var diagRe = regexp.MustCompile(`^(.+?):(\d+):(\d+): (?:ERROR|WARNING) - \[([A-Z_]+)\]`)

func parseTuples2(text string) []tuple {
	var out []tuple
	for _, line := range strings.Split(text, "\n") {
		m := diagRe.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		ln, _ := strconv.Atoi(m[2])
		cl, _ := strconv.Atoi(m[3])
		out = append(out, tuple{path: m[1], line: ln, col: cl, code: m[4]})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].line != out[j].line {
			return out[i].line < out[j].line
		}
		if out[i].col != out[j].col {
			return out[i].col < out[j].col
		}
		return out[i].code < out[j].code
	})
	return out
}

func main2() {
	root := "testdata/fixtures/extracted"
	count := 0
	_ = filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(p, ".js") {
			return nil
		}
		if count >= 5 {
			return nil
		}

		fixRel := p
		oraclePath := filepath.Join("testdata/oracle",
			strings.TrimSuffix(fixRel[len("testdata/fixtures/"):], ".js")+".txt")
		oracleBytes, err := os.ReadFile(oraclePath)
		if err != nil {
			return nil
		}

		expected := parseTuples2(string(oracleBytes))
		hasMismatch := false
		for _, e := range expected {
			if e.code == "JSC_TYPE_MISMATCH" {
				hasMismatch = true
				break
			}
		}
		if !hasMismatch {
			return nil
		}

		src, _ := os.ReadFile(fixRel)
		diags, _ := jscomp.CheckSourceWithOptions(fixRel, string(src), jscomp.NewCompilerOptions())
		var actual []tuple
		for _, d := range diags {
			actual = append(actual, tuple{path: d.Path, line: d.Line, col: d.Col, code: d.Code})
		}
		sort.Slice(actual, func(i, j int) bool {
			if actual[i].line != actual[j].line {
				return actual[i].line < actual[j].line
			}
			if actual[i].col != actual[j].col {
				return actual[i].col < actual[j].col
			}
			return actual[i].code < actual[j].code
		})

		// Check for missed mismatches
		for _, e := range expected {
			if e.code != "JSC_TYPE_MISMATCH" {
				continue
			}
			found := false
			for _, a := range actual {
				if a == e {
					found = true
					break
				}
			}
			if !found {
				fmt.Printf("\n=== %s (missed %s at %d:%d)\n", fixRel, e.code, e.line, e.col)
				fmt.Printf("Expected: %v\n", expected)
				fmt.Printf("Got: %v\n", actual)
				// Show the JS source around that line
				lines := strings.Split(string(src), "\n")
				start := e.line - 3
				if start < 0 {
					start = 0
				}
				end := e.line + 1
				if end > len(lines) {
					end = len(lines)
				}
				for i := start; i < end; i++ {
					prefix := "  "
					if i+1 == e.line {
						prefix = ">>"
					}
					fmt.Printf("%s%4d: %s\n", prefix, i+1, lines[i])
				}
				count++
				break
			}
		}
		return nil
	})
}
