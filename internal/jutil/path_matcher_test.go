package jutil

import (
	"os"
	"strings"
	"testing"
)

// TestPathMatcherGlobVectors checks the glob PathMatcher against
// testdata/glob_vectors.tsv (tools/oracle/java/JutilGlobVectors.java).
func TestPathMatcherGlobVectors(t *testing.T) {
	data, err := os.ReadFile("testdata/glob_vectors.tsv")
	if err != nil {
		t.Fatal(err)
	}
	paths := []string{"a.js", "b.js", "c.js", "d.js", "-.js", "./a.js", "./x/a.js", "x/y/a.js", "ab.js", "a*.js", "./src/x_test.js", "./src/d/y_test.js", "a/b", "]].js", "].js", "x*y", "x?y", "x\\y", "axb", ".js"}
	for _, line := range strings.Split(strings.TrimRight(string(data), "\n"), "\n") {
		parts := strings.SplitN(line, "\t", 3)
		glob, kind, want := parts[0], parts[1], parts[2]
		var got string
		Try(func() {
			m := FileSystem_GetPathMatcher("glob:" + glob)
			var matched []string
			for _, p := range paths {
				if m.Matches(Path_Of(p)) {
					matched = append(matched, p)
				}
			}
			got = "OK\t" + strings.Join(matched, " ")
		}, func(e Throwable) { got = "ERR\t" + strings.ReplaceAll(e.Error(), "\n", "\\n") })
		if strings.Contains(want, "Unclosed character class") {
			// Java's message quotes its internal regex (see path_matcher.go).
			if !strings.HasPrefix(got, "ERR\tjava.util.regex.PatternSyntaxException: Unclosed character class") {
				t.Errorf("glob %q: got %s, want an Unclosed character class error", glob, got)
			}
			continue
		}
		if got != kind+"\t"+want {
			t.Errorf("glob %q:\n got %s\nwant %s\t%s", glob, got, kind, want)
		}
	}
}
