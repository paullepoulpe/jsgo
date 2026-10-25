package parser

import "testing"

// Go-only: lines and columns come from the file's line table (see
// SourcePosition): \r\n is one line terminator, and U+2028 is one too.
func TestSourceFileLinePositions(t *testing.T) {
	f := NewSourceFile("positions.js", "x\r\ny z")
	want := []LinePosition{
		{0, 0, 0}, {1, 0, 1}, {2, 0, 2}, {3, 1, 0}, {4, 1, 1}, {5, 2, 0}, {6, 2, 1},
	}
	for offset, w := range want {
		if got := f.LinePosition(NewSourcePosition(offset)); got != w {
			t.Fatalf("offset %d: %+v != %+v", offset, got, w)
		}
	}
	// Lookups in any order (the parser rewinds).
	for _, offset := range []int{6, 0, 3, 5, 1, 4, 2} {
		if got := f.LinePosition(NewSourcePosition(offset)); got != want[offset] {
			t.Fatalf("offset %d out of order: %+v != %+v", offset, got, want[offset])
		}
	}
}
