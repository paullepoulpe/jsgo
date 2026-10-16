package jutil

import (
	"bufio"
	"os"
	"strings"
	"testing"
)

// TestPatternGroupDollarMatchesJava checks find()/matches() of regexes with a
// `$` inside a group or before other alternatives against vectors printed by
// JDK 21 (tools/oracle/java/JutilPatternGroupDollarVectors.java). The inputs
// of the find() vectors do not end with a line terminator (where such a
// find() is unsupported); TestPatternGroupDollarUnsupportedFind checks that.
func TestPatternGroupDollarMatchesJava(t *testing.T) {
	f, err := os.Open("testdata/pattern_group_dollar.tsv")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	n := 0
	for sc.Scan() {
		parts := strings.Split(sc.Text(), "\t")
		n++
		p := Pattern_Compile(dollarUnits(t, parts[1]))
		in := dollarUnits(t, parts[2])
		m := p.Matcher(in)
		switch parts[0] {
		case "find":
			var got []string
			for k := 0; k < 50 && m.Find(); k++ {
				got = append(got, dollarMatch(m))
			}
			if g := strings.Join(got, ";"); g != parts[3] {
				t.Errorf("find %q in %q: got %s, want %s", p.Pattern(), in, g, parts[3])
			}
		case "matches":
			got := "-"
			if m.Matches() {
				got = dollarMatch(m)
			}
			if got != parts[3] {
				t.Errorf("matches %q on %q: got %s, want %s", p.Pattern(), in, got, parts[3])
			}
		default:
			t.Fatalf("bad line %q", sc.Text())
		}
	}
	if n < 500 {
		t.Fatalf("only %d vectors", n)
	}
}

// TestPatternGroupDollarFindBeforeTerminator checks find() of "(code$)" on
// inputs ending with a line terminator, where `$` matches before it: JDK 21
// reports one match at 0-4, group 1 at 0-4, for each of these inputs.
func TestPatternGroupDollarFindBeforeTerminator(t *testing.T) {
	for _, in := range []string{"code\n", "code\r\n", "code\r", "code\u0085", "code\u2028", "code\u2029"} {
		m := Pattern_Compile("(code$)").Matcher(in)
		var got []string
		for m.Find() {
			got = append(got, dollarMatch(m))
		}
		if len(got) != 1 {
			t.Errorf("find on %q: got %v, want one match", in, got)
		}
		m = Pattern_Compile("(code$)").Matcher(in)
		if !m.Find() || m.Start() != 0 || m.End() != 4 || m.StartGroup(1) != 0 || m.EndGroup(1) != 4 {
			t.Errorf("find on %q: want 0-4, group 1 0-4", in)
		}
	}
}
