package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jscomp"
)

type tuple struct {
	line int
	col  int
	code string
}

var diagRe = regexp.MustCompile(`^.+?:(\d+):(\d+): (?:ERROR|WARNING) - \[([A-Z_]+)\]`)

func parseTuples(text string) []tuple {
	var out []tuple
	for _, line := range strings.Split(text, "\n") {
		m := diagRe.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		ln := 0
		fmt.Sscanf(m[1], "%d", &ln)
		cl := 0
		fmt.Sscanf(m[2], "%d", &cl)
		out = append(out, tuple{line: ln, col: cl, code: m[3]})
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
	// If a code filter arg is given, only print mismatches that have
	// extras or missings of that code, and show the fixture name.
	filterCode := ""
	if len(os.Args) > 1 {
		filterCode = os.Args[1]
	}

	root := "/home/user/jsgo/testdata/fixtures/extracted"
	oracleRoot := "/home/user/jsgo/testdata/oracle/extracted"
	extraByCode := map[string]int{}
	missingByCode := map[string]int{}
	totalMismatch := 0
	entries, _ := ioutil.ReadDir(root)
	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), ".js") {
			continue
		}
		fixPath := filepath.Join(root, e.Name())
		oracPath := filepath.Join(oracleRoot, strings.TrimSuffix(e.Name(), ".js")+".txt")
		src, err := ioutil.ReadFile(fixPath)
		if err != nil {
			continue
		}
		orac, err := ioutil.ReadFile(oracPath)
		if err != nil {
			continue
		}
		diags, _ := jscomp.CheckSource(fixPath, string(src))
		var actual []string
		for _, d := range diags {
			actual = append(actual, d.FormatJavaLike())
		}
		aT := parseTuples(strings.Join(actual, "\n"))
		eT := parseTuples(string(orac))
		if len(aT) != len(eT) || !tuplesEqual(aT, eT) {
			totalMismatch++
			var extras, missings []tuple
			em := map[tuple]int{}
			for _, t := range eT {
				em[t]++
			}
			for _, t := range aT {
				if em[t] > 0 {
					em[t]--
				} else {
					extraByCode[t.code]++
					extras = append(extras, t)
				}
			}
			am := map[tuple]int{}
			for _, t := range aT {
				am[t]++
			}
			for _, t := range eT {
				if am[t] > 0 {
					am[t]--
				} else {
					missingByCode[t.code]++
					missings = append(missings, t)
				}
			}
			if filterCode != "" {
				matched := false
				for _, x := range extras {
					if x.code == filterCode {
						matched = true
					}
				}
				for _, x := range missings {
					if x.code == filterCode {
						matched = true
					}
				}
				if matched {
					fmt.Printf("== %s\n", e.Name())
					for _, x := range extras {
						if x.code == filterCode {
							fmt.Printf("  EXTRA %d:%d %s\n", x.line, x.col, x.code)
						}
					}
					for _, x := range missings {
						if x.code == filterCode {
							fmt.Printf("  MISS  %d:%d %s\n", x.line, x.col, x.code)
						}
					}
				}
			}
		}
	}
	fmt.Printf("Total mismatched: %d\n", totalMismatch)
	fmt.Println("Extra (we emit but oracle doesn't) by code:")
	for k, v := range extraByCode {
		fmt.Printf("  %s: %d\n", k, v)
	}
	fmt.Println("Missing (oracle expects but we don't emit) by code:")
	for k, v := range missingByCode {
		fmt.Printf("  %s: %d\n", k, v)
	}
}

func tuplesEqual(a, b []tuple) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
