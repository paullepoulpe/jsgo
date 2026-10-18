package parser

import "testing"

func TestTokenPositionCacheValuesAndRewind(t *testing.T) {
	source := NewSourceFile("positions.js", "x\r\ny\u2028z")
	cached, plain := NewLineNumberScanner(source), NewLineNumberScanner(source)
	for _, offset := range []int{0, 0, 1, 1, 3, 3, 4, 5, 6, 6} {
		got, want := cached.getTokenSourcePosition(offset), plain.GetSourcePosition(offset)
		if *got != *want {
			t.Fatalf("offset %d: %#v != %#v", offset, got, want)
		}
		// A long-lived comment must never alias a token-slab position.
		if comment := cached.GetSourcePosition(offset); comment == got || *comment != *got {
			t.Fatal("comment position aliases token storage")
		}
	}
	rewind := NewSourcePosition(source, 0, 0, 0)
	cached.RewindTo(rewind)
	if got := cached.getTokenSourcePosition(0); got.Offset != 0 || got.Line != 0 || got.Column != 0 {
		t.Fatal(got)
	}
}
