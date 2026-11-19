package javaregex

import (
	"reflect"
	"strings"
	"testing"
)

func TestPatternAndMatcher(t *testing.T) {
	p := Pattern_Compile(`([a-z]+)=(\d+)`)
	m := p.Matcher("foo=12, bar=345")
	if !m.Find() {
		t.Fatal("expected first match")
	}
	if g1, ok := m.Group(1); !ok || g1 != "foo" {
		t.Fatalf("group 1 = %q, %v", g1, ok)
	}
	if g2, ok := m.Group(2); !ok || g2 != "12" {
		t.Fatalf("group 2 = %q, %v", g2, ok)
	}
	if m.ReplaceAll("$2:$1") != "12:foo, 345:bar" {
		t.Fatalf("ReplaceAll got %q", m.ReplaceAll("$2:$1"))
	}

	// Java dot excludes \r and \u2028 unless DOTALL is set
	dotP := Pattern_Compile(`a.b`)
	if dotP.Matches("a\rb") || dotP.Matches("a\u2028b") {
		t.Fatal("Java '.' without DOTALL must not match \\r or \\u2028")
	}
	dotAllP := Pattern_CompileFlags(`a.b`, Pattern_DOTALL)
	if !dotAllP.Matches("a\rb") || !dotAllP.Matches("a\nb") {
		t.Fatal("Java '.' with DOTALL must match line terminators")
	}

	// Split with limit == 0 trims trailing empties, limit == -1 preserves them
	sp := Pattern_Compile(`,\s*`)
	if got := sp.Split("a, b, , "); !reflect.DeepEqual(got, []string{"a", "b"}) {
		t.Fatalf("Split limit=0 got %#v, want [a b]", got)
	}
	if got := sp.SplitLimit("a, b, , ", -1); !reflect.DeepEqual(got, []string{"a", "b", "", ""}) {
		t.Fatalf("SplitLimit -1 got %#v, want [a b '' '']", got)
	}

	// AppendReplacement / AppendTail
	var sb strings.Builder
	m.Reset()
	for m.Find() {
		m.AppendReplacement(&sb, Matcher_QuoteReplacement("$literal$"))
	}
	m.AppendTail(&sb)
	if sb.String() != "$literal$, $literal$" {
		t.Fatalf("AppendReplacement/Tail = %q", sb.String())
	}
}
