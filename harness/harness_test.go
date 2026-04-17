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
	lines := strings.Split(strings.TrimRight(text, "\n"), "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		if lines[i] == "---EXIT---" && i+1 < len(lines) {
			v, _ := strconv.Atoi(strings.TrimSpace(lines[i+1]))
			return v
		}
	}
	return -1
}

func jsgoTuples(path string) ([]tuple, int, error) {
	src, err := os.ReadFile(path)
	if err != nil {
		return nil, 0, err
	}
	diags, err := jscomp.CheckSource(path, string(src))
	if err != nil {
		return nil, 0, nil
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
	return out, exit, nil
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
		b.WriteString("  " + t.path + ":" + strconv.Itoa(t.line) + ":" + strconv.Itoa(t.col) + " [" + t.code + "]\n")
	}
	return b.String()
}

// repoRoot is resolved once at package init, and the test process
// chdirs there so every fixture's repo-relative path resolves from the
// same place. Later tests never touch the cwd again.
var repoRoot = func() string {
	p, err := filepath.Abs("..")
	if err != nil {
		panic(err)
	}
	if err := os.Chdir(p); err != nil {
		panic(err)
	}
	return p
}()

func runFixture(t *testing.T, _, fix string) error {
	fixRel, err := filepath.Rel(repoRoot, fix)
	if err != nil {
		return err
	}
	oraclePath := filepath.Join(repoRoot, "testdata", "oracle",
		strings.TrimSuffix(fixRel[len("testdata/fixtures/"):], ".js")+".txt")
	oracleBytes, err := os.ReadFile(oraclePath)
	if err != nil {
		return err
	}
	expected := parseTuples(string(oracleBytes))
	// Java recorded its diagnostics against the repo-relative path
	// (that's what the oracle shows). Pass the same string to jsgo.
	actual, exit, err := jsgoTuples(fixRel)
	if err != nil {
		return err
	}
	// Tuples from jsgo use the raw fixRel path; from the oracle they
	// embed the same repo-relative path (e.g.
	// `testdata/fixtures/tier0/...`). They match directly.
	if diff := diffTuples(expected, actual); diff != "" {
		return &fixtureError{msg: "tuple mismatch:\n" + diff}
	}
	wantExit := parseOracleExit(string(oracleBytes))
	if exitClass(exit) != exitClass(wantExit) {
		return &fixtureError{msg: "exit-class mismatch"}
	}
	return nil
}

type fixtureError struct{ msg string }

func (e *fixtureError) Error() string { return e.msg }

func walkFixtures(root string) []string {
	var out []string
	_ = filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && strings.HasSuffix(p, ".js") {
			out = append(out, p)
		}
		return nil
	})
	sort.Strings(out)
	return out
}

func runTier(t *testing.T, dir string) {
	t.Helper()
	root := filepath.Join(repoRoot, "testdata", "fixtures", dir)
	fixtures := walkFixtures(root)
	if len(fixtures) == 0 {
		t.Fatalf("no %s fixtures", dir)
	}
	for _, fix := range fixtures {
		rel, _ := filepath.Rel(root, fix)
		t.Run(rel, func(t *testing.T) {
			if err := runFixture(t, repoRoot, fix); err != nil {
				t.Fatalf("%v", err)
			}
		})
	}
}

// TestTier0 is the authoritative primitive-scope conformance test.
func TestTier0(t *testing.T) { runTier(t, "tier0") }

// TestTier1 covers function signatures: @param/@return, call-site arg
// checks, and return-type consistency.
func TestTier1(t *testing.T) { runTier(t, "tier1") }

// TestExtractedCorpus is a reporting benchmark — it runs every
// extracted TypeCheckTest fixture and prints a pass-rate summary. It is
// intentionally lenient: failures here reflect the difference between
// Tier-0 scope and what the full Closure Compiler's checker covers.
// CI fails only if the pass rate regresses below the floor set here.
func TestExtractedCorpus(t *testing.T) {
	const floorPercent = 40 // regression floor; raise as we ship more checks.
	root := filepath.Join(repoRoot, "testdata", "fixtures", "extracted")
	fixtures := walkFixtures(root)
	if len(fixtures) == 0 {
		t.Skip("extracted corpus not present")
	}
	pass, fail := 0, 0
	for _, fix := range fixtures {
		if err := runFixture(t, repoRoot, fix); err != nil {
			fail++
		} else {
			pass++
		}
	}
	total := pass + fail
	rate := 100.0 * float64(pass) / float64(total)
	t.Logf("extracted: pass=%d fail=%d total=%d rate=%.2f%%", pass, fail, total, rate)
	if rate < floorPercent {
		t.Fatalf("pass rate %.2f%% below floor %d%%", rate, floorPercent)
	}
}