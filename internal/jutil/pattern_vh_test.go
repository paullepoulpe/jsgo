package jutil

import (
	"os"
	"strconv"
	"strings"
	"testing"
)

// Go-only: the \v, \V, \h and \H classes of jutil.Pattern against
// java.util.regex (testdata/pattern_vh.tsv, written by
// tools/oracle/java/JutilPatternVhVectors.java).
func TestPatternVerticalHorizontalWhitespace(t *testing.T) {
	data, err := os.ReadFile("testdata/pattern_vh.tsv")
	if err != nil {
		t.Fatal(err)
	}
	for _, line := range strings.Split(strings.TrimSuffix(string(data), "\n"), "\n") {
		f := strings.Split(line, "\t")
		regex, chars, bits := f[0], strings.Split(f[1], ","), f[2]
		p := Pattern_Compile(regex)
		for i, h := range chars {
			c, _ := strconv.ParseUint(h, 16, 16)
			in := StringFromChar(uint16(c))
			if strings.HasPrefix(regex, "x") {
				in = "x" + in + "y"
			}
			if got, want := p.Matches(in), bits[i] == '1'; got != want {
				t.Errorf("Pattern %q matches U+%04X: got %v, want %v", regex, c, got, want)
			}
		}
	}
}
