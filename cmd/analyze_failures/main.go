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

	"github.com/paullepoulpe/jsgo/internal/jscomp"
)

type tuple struct {
	path string
	line int
	col  int
	code string
}

var diagRe = regexp.MustCompile(`^(.+?):(\d+):(\d+): (?:ERROR|WARNING) - \[([A-Z_]+)\]`)

func parseTuples(text string) []tuple {
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

func main() {
	root := "testdata/fixtures/extracted"
	failCodes := map[string]int{}
	missCodes := map[string]int{}
	failFiles := 0
	passFiles := 0

	_ = filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(p, ".js") {
			return nil
		}

		fixRel := p
		oraclePath := filepath.Join("testdata/oracle",
			strings.TrimSuffix(fixRel[len("testdata/fixtures/"):], ".js")+".txt")
		oracleBytes, err := os.ReadFile(oraclePath)
		if err != nil {
			return nil
		}

		expected := parseTuples(string(oracleBytes))
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

		match := len(expected) == len(actual)
		if match {
			for i := range expected {
				if expected[i] != actual[i] {
					match = false
					break
				}
			}
		}

		if match {
			passFiles++
			return nil
		}
		failFiles++

		// Count false positives (extra diags)
		for _, a := range actual {
			found := false
			for _, e := range expected {
				if a == e {
					found = true
					break
				}
			}
			if !found {
				failCodes[a.code]++
			}
		}
		// Count false negatives (missed diags)
		for _, e := range expected {
			found := false
			for _, a := range actual {
				if a == e {
					found = true
					break
				}
			}
			if !found {
				missCodes[e.code]++
			}
		}
		return nil
	})

	fmt.Printf("pass=%d fail=%d\n\n", passFiles, failFiles)

	fmt.Println("False positives (extra diags generated):")
	type kv struct {
		k string
		v int
	}
	var fps []kv
	for k, v := range failCodes {
		fps = append(fps, kv{k, v})
	}
	sort.Slice(fps, func(i, j int) bool { return fps[i].v > fps[j].v })
	for _, kv := range fps {
		fmt.Printf("  %s: %d\n", kv.k, kv.v)
	}

	fmt.Println("\nFalse negatives (missed expected diags):")
	var fns []kv
	for k, v := range missCodes {
		fns = append(fns, kv{k, v})
	}
	sort.Slice(fns, func(i, j int) bool { return fns[i].v > fns[j].v })
	for _, kv := range fns {
		fmt.Printf("  %s: %d\n", kv.k, kv.v)
	}
}
