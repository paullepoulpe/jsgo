package parsing

// FeatureCollector has no Java unit test; it is checked against the Java
// class itself (tools/oracle/java/FeatureCollectorDump.java): each file is
// parsed with ParserRunner.parse (the compiler's default configuration, as
// in parsingdump_test.go) and FeatureCollector.visitScript runs over the
// SCRIPT.
//
//   - TestFeatureCollectorGolden compares testdata/featurecollector.tsv
//     (the Java output for testdata/parse/src).
//   - TestFeatureCollectorOracle compares a whole set when FEATURES_ROOT (the
//     root of the sources) and FEATURES_TSV (the FeatureCollectorDump output
//     for them) are set.

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// collectFeatures is FeatureCollectorDump's per-file line.
func collectFeatures(rel, src string) string {
	var ast *rhino.Node
	jutil.Try(func() {
		ast = ParserRunner_Parse(
			rhino.NewSimpleSourceFile(rel, rhino.StaticSourceFile_SourceKind_STRONG), src, defaultDumpConfig(),
			rhino.ErrorReporter_NULL_INSTANCE).Ast
	}, func(e jutil.IRuntimeException) { ast = nil })
	if ast == nil {
		return rel + "\tnull"
	}
	fc := NewFeatureCollector()
	fc.VisitScript(ast)
	return rel + "\t" + fc.AllFeatures().String()
}

func compareFeatureCollector(t *testing.T, root, tsv string) {
	f, err := os.Open(tsv)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 1<<20), 1<<20)
	n, differ := 0, 0
	for sc.Scan() {
		want := sc.Text()
		rel, _, _ := strings.Cut(want, "\t")
		src, err := os.ReadFile(filepath.Join(root, rel))
		if err != nil {
			t.Fatal(err)
		}
		n++
		if got := collectFeatures(rel, string(src)); got != want {
			differ++
			if differ <= 20 {
				t.Errorf("FeatureCollector differs from Java:\n  java: %s\n  go:   %s", want, got)
			}
		}
	}
	t.Logf("compared %d files: %d differ", n, differ)
	if n == 0 {
		t.Fatal("no files compared")
	}
}

func TestFeatureCollectorGolden(t *testing.T) {
	compareFeatureCollector(t, "testdata/parse/src", "testdata/featurecollector.tsv")
}

func TestFeatureCollectorOracle(t *testing.T) {
	root, tsv := os.Getenv("FEATURES_ROOT"), os.Getenv("FEATURES_TSV")
	if root == "" || tsv == "" {
		t.Skip("set FEATURES_ROOT and FEATURES_TSV to run the FeatureCollector oracle")
	}
	compareFeatureCollector(t, root, tsv)
}
