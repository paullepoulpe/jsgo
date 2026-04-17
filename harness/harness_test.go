package harness_test

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/diag"
	"github.com/paullepoulpe/jsgo/internal/jscomp"
)

// tuple is (path, line, col, code) as specified in SPEC.md. Message
// strings are deliberately excluded.
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
		return out[i].col < out[j].col
	})
	return out
}

func parseOracleExit(text string) int {
	// Oracle format: raw stderr, then a line "---EXIT---", then the numeric exit code.
	lines := strings.Split(strings.TrimRight(text, "\n"), "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		if lines[i] == "---EXIT---" && i+1 < len(lines) {
			v, _ := strconv.Atoi(strings.TrimSpace(lines[i+1]))
			return v
		}
	}
	return -1
}

func jsgoTuples(t *testing.T, path string) ([]tuple, int) {
	t.Helper()
	src, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	diags, err := jscomp.CheckSource(path, string(src))
	if err != nil {
		t.Fatalf("jscomp.CheckSource(%s): %v", path, err)
	}
	var out []tuple
	exit := 0
	for _, d := range diags {
		out = append(out, tuple{path: d.Path, line: d.Line, col: d.Col, code: d.Code})
		if d.Level == diag.Error {
			exit = 1
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].line != out[j].line {
			return out[i].line < out[j].line
		}
		return out[i].col < out[j].col
	})
	return out, exit
}

// TestConformance walks testdata/fixtures/**/*.js and verifies jsgo's
// diagnostic tuples + exit-code class match the cached Java oracle.
func TestConformance(t *testing.T) {
	// testdata paths are relative to the repo root; tests run from the
	// package dir, so step up once.
	repoRoot, err := filepath.Abs("..")
	if err != nil {
		t.Fatal(err)
	}
	fixturesRoot := filepath.Join(repoRoot, "testdata", "fixtures")
	oracleRoot := filepath.Join(repoRoot, "testdata", "oracle")

	var fixtures []string
	_ = filepath.Walk(fixturesRoot, func(p string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && strings.HasSuffix(p, ".js") {
			fixtures = append(fixtures, p)
		}
		return nil
	})
	if len(fixtures) == 0 {
		t.Fatalf("no fixtures found under %s", fixturesRoot)
	}

	for _, fix := range fixtures {
		fix := fix
		rel, _ := filepath.Rel(fixturesRoot, fix)
		t.Run(rel, func(t *testing.T) {
			oraclePath := filepath.Join(oracleRoot, strings.TrimSuffix(rel, ".js")+".txt")
			oracleBytes, err := os.ReadFile(oraclePath)
			if err != nil {
				t.Fatalf("oracle missing: %v", err)
			}

			// Invoke jsgo with the same path string Java was given, so
			// the path component of diagnostics matches. That means
			// running from the repo root and passing a repo-relative path.
			fixRel, err := filepath.Rel(repoRoot, fix)
			if err != nil {
				t.Fatal(err)
			}
			if err := os.Chdir(repoRoot); err != nil {
				t.Fatal(err)
			}

			expected := parseTuples(string(oracleBytes))
			actual, exit := jsgoTuples(t, fixRel)

			if diff := diffTuples(expected, actual); diff != "" {
				t.Fatalf("tuple mismatch:\n%s", diff)
			}
			wantExit := parseOracleExit(string(oracleBytes))
			gotClass := exitClass(exit)
			wantClass := exitClass(wantExit)
			if gotClass != wantClass {
				t.Fatalf("exit-code class mismatch: jsgo=%d (class %s), oracle=%d (class %s)",
					exit, gotClass, wantExit, wantClass)
			}
		})
	}
}

func exitClass(e int) string {
	if e == 0 {
		return "zero"
	}
	return "nonzero"
}

func diffTuples(want, got []tuple) string {
	if len(want) != len(got) {
		return tupleDump("want", want) + tupleDump("got", got)
	}
	for i := range want {
		if want[i] != got[i] {
			return tupleDump("want", want) + tupleDump("got", got)
		}
	}
	return ""
}

func tupleDump(label string, ts []tuple) string {
	var b strings.Builder
	b.WriteString(label + ":\n")
	for _, t := range ts {
		b.WriteString("  ")
		b.WriteString(t.path)
		b.WriteString(":")
		b.WriteString(strconv.Itoa(t.line))
		b.WriteString(":")
		b.WriteString(strconv.Itoa(t.col))
		b.WriteString(" [")
		b.WriteString(t.code)
		b.WriteString("]\n")
	}
	return b.String()
}
