package jutil

import (
	"strconv"
	"strings"
	"testing"
)

// Black-box JDK 21 results (printed by scratch programs that call only the
// public java.util.regex API) for what pattern_dollar_test.go's vectors do
// not reach: find() after failures and after matches()/lookingAt(), and `$`
// before the rarer line terminators.

// matcherScript runs ops on a matcher: f find(), m matches(), l
// lookingAt(), r reset(); a successful find prints "start,end", a failed one
// "-", matches()/lookingAt() print their result.
func matcherScript(regex, input, ops string) string {
	m := Pattern_Compile(regex).Matcher(input)
	var out []string
	for _, op := range ops {
		switch op {
		case 'f':
			if m.Find() {
				out = append(out, strconv.Itoa(m.Start())+","+strconv.Itoa(m.End()))
			} else {
				out = append(out, "-")
			}
		case 'm':
			out = append(out, strconv.FormatBool(m.Matches()))
		case 'l':
			out = append(out, strconv.FormatBool(m.LookingAt()))
		case 'r':
			m.Reset()
			out = append(out, "reset")
		}
	}
	return strings.Join(out, " ")
}

func TestMatcherFindStateObserved(t *testing.T) {
	for _, c := range []struct{ regex, input, ops, want string }{
		{"a", "a", "ffff", "0,1 - - -"},
		{"a", "aXa", "ffff", "0,1 2,3 - -"},
		{"a*", "baab", "fffffff", "0,0 1,3 3,3 4,4 - - -"},
		{"a+", "aa a", "mff", "false 0,2 3,4"},
		{"a+", "aa", "mff", "true - -"},
		{"a+", "aa aa", "lff", "true 3,5 -"},
		{"a*", "ba", "lfff", "true 1,2 2,2 -"},
		{"a*", "", "mff", "true - -"},
		{"b", "ab", "lff", "false 1,2 -"},
		{"a", "aa", "fffrf", "0,1 1,2 - reset 0,1"},
		{"x*", "ab", "fffff", "0,0 1,1 2,2 - -"},
		{"a|", "aab", "ffffff", "0,1 1,2 2,2 3,3 - -"},
		{"a", "aXa", "fmff", "0,1 false 2,3 -"},
		{"a", "aXaXa", "fflff", "0,1 2,3 true 2,3 4,5"},
		{"a*", "bab", "fmff", "0,0 false 0,0 1,2"},
		{"X", "aXa", "flff", "1,2 false - -"},
		{"a", "a", "ffmf", "0,1 - true -"},
		{"a", "aa", "fffl" + "ff", "0,1 1,2 - true 1,2 -"},
		{"a*", "aa", "fff", "0,2 2,2 -"},
		{"a*", "ab", "ffff", "0,1 1,1 2,2 -"},
		{"x*", "ab", "ffffmff", "0,0 1,1 2,2 - false 2,2 -"},
		{"x*", "ab", "fffflff", "0,0 1,1 2,2 - true 1,1 2,2"},
		{"x*", "ab", "ffffflf", "0,0 1,1 2,2 - - true 1,1"},
		{"b|", "ab", "ffffmf", "0,0 1,2 2,2 - false 2,2"},
		{"a|", "a", "fffmff", "0,1 1,1 - true 1,1 -"},
		{"", "", "ffmflf", "0,0 - true - true -"},
		{"^", "ab", "ffff", "0,0 - 0,0 -"},
		{"^", "ab", "ffmff", "0,0 - false 0,0 -"},
		{"a*", "ab", "fmf", "0,1 false 1,1"},
	} {
		if got := matcherScript(c.regex, c.input, c.ops); got != c.want {
			t.Errorf("%q on %q, %s: got %s, want %s", c.regex, c.input, c.ops, got, c.want)
		}
	}
	m := Pattern_Compile("a").Matcher("ab")
	m.Find()
	m.Find()
	Try(func() { m.Start(); t.Error("start() after a failed find(): no exception") }, func(e IRuntimeException) {
		if e.GetMessage() != "No match found" {
			t.Errorf("message %q", e.GetMessage())
		}
	})
}

func TestPatternDollarTerminatorsObserved(t *testing.T) {
	for _, c := range []struct{ regex, input, want string }{
		{".*$", "x\u0085", "0,1 1,1 2,2"},
		{"x$", "x\u0085", "0,1"},
		{"[^/]+$", "x\u0085", "0,2"},
		{".*$", "x ", "0,1 1,1 2,2"},
		{"x$", "x ", "0,1"},
		{".*$", "x\n\r", "2,2 3,3"},
		{"x$", "x\n\r", ""},
		{"[^/]+$", "x\n\r", "0,3"},
		{".*$", "x\r\r", "2,2 3,3"},
		{".*$", "\n", "0,0 1,1"},
		{".*$", "\r", "0,0 1,1"},
		{".*$", "x\x0B", "0,2 2,2"},
		{"x$", "x\x0B", ""},
		{".*$", "x\x0C", "0,2 2,2"},
		{"(x$)|(y)", "x\n", "0,1"},
	} {
		var got []string
		for m := Pattern_Compile(c.regex).Matcher(c.input); m.Find(); {
			got = append(got, strconv.Itoa(m.Start())+","+strconv.Itoa(m.End()))
		}
		if g := strings.Join(got, " "); g != c.want {
			t.Errorf("find %q in %q: got %s, want %s", c.regex, c.input, g, c.want)
		}
	}
	m := Pattern_Compile("(x)($)").Matcher("x\r\n")
	if !m.Find() || m.Start() != 0 || m.End() != 1 || m.startGroup(1) != 0 || m.endGroup(1) != 1 || m.startGroup(2) != 1 || m.endGroup(2) != 1 {
		t.Errorf("(x)($) on x\\r\\n: %v", m.groups)
	}
	for _, regex := range []string{"(a$)*", "a$\\n", "(a$)?", "a$?", "a$(?s).", "a$[b]"} {
		Try(func() { Pattern_Compile(regex); t.Errorf("%q compiled", regex) }, func(e IRuntimeException) {
			if _, ok := e.(*UnsupportedOperationException); !ok {
				t.Errorf("%q: %v", regex, e)
			}
		})
	}
}

func TestPatternWordBoundariesAndPostDollarLiterals(t *testing.T) {
	wb := Pattern_Compile(`\bfoo\b`)
	if !wb.Matches("foo") || wb.Matches("foobar") || !wb.Find("a foo bar") || wb.Find("a foobar") {
		t.Fatalf(`\bfoo\b mismatch`)
	}
	nwb := Pattern_Compile(`\Bfoo\B`)
	if nwb.Matches("foo") || !nwb.Find("xfooy") || nwb.Find("foo bar") {
		t.Fatalf(`\Bfoo\B mismatch`)
	}
	for _, cls := range []string{`[\b]`, `[\B]`} {
		Try(func() {
			Pattern_Compile(cls)
			t.Errorf("%q compiled; expected PatternSyntaxException", cls)
		}, func(e IRuntimeException) {
			if _, ok := e.(*PatternSyntaxException); !ok {
				t.Errorf("%q: got %T (%v), want *PatternSyntaxException", cls, e, e)
			}
		})
	}

	for _, unmatchable := range []string{"a$b", "foo$bar", "(foo$)bar", "^$a", `foo$.`} {
		p := Pattern_Compile(unmatchable)
		for _, in := range []string{"", "a", "ab", "a\n", "a\r\n", "foo", "foobar", "foo\n", "foo\nbar"} {
			if p.Matches(in) || p.Find(in) {
				t.Errorf("%q unexpectedly matched %q", unmatchable, in)
			}
		}
	}

	alt := Pattern_Compile("foo$bar|baz")
	if !alt.Matches("baz") || !alt.Find("xxbazyy") || alt.Find("foobar") || alt.Find("foo\n") {
		t.Fatalf("foo$bar|baz mismatch")
	}

	grp := Pattern_Compile("(a$|b)c")
	m := grp.Matcher("bc")
	if !m.Matches() {
		t.Fatalf("(a$|b)c on bc mismatch")
	}
	if g1, ok := m.Group(1); !ok || g1 != "b" {
		t.Fatalf("(a$|b)c group(1) = (%q, %v), want (\"b\", true)", g1, ok)
	}
	if grp.Find("ac") || grp.Find("a\n") || grp.Find("a\nc") {
		t.Fatalf("(a$|b)c unexpectedly matched")
	}

	optTail := Pattern_Compile("foo$b*")
	if !optTail.Matches("foo") || optTail.Matches("foo\n") {
		t.Fatalf("foo$b* matches() mismatch")
	}
	m2 := optTail.Matcher("foo\n")
	if !m2.Find() || m2.Start() != 0 || m2.End() != 3 {
		t.Fatalf("foo$b* find() on foo\\n got %v", m2.groups)
	}
}

func TestPatternMultilineAndDotallMatches(t *testing.T) {
	pmEmbedded := Pattern_Compile("(?m)^foo$")
	if pmEmbedded.Matches("foo\nbar") {
		t.Errorf("(?m) Matches should anchor to full string, matched 'foo\\nbar'")
	}
	if !pmEmbedded.Find("foo\nbar") {
		t.Errorf("(?m) Find should match line in 'foo\\nbar'")
	}
	ps := Pattern_Compile("(?s)foo.*bar")
	if !ps.Matches("foo\nbar") {
		t.Errorf("(?s) Matches should match across newline in 'foo\\nbar'")
	}
}
