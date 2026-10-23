package parser

import "testing"

func TestLineNumberScannerPositionsAndRewind(t *testing.T) {
	l := NewLineNumberScanner(NewSourceFile("positions.js", "x\r\ny z"))
	want := []SourcePosition{
		{0, 0, 0}, {1, 0, 1}, {2, 0, 2}, {3, 1, 0}, {4, 1, 1}, {5, 2, 0}, {6, 2, 1},
	}
	for offset, w := range want {
		if got := l.GetSourcePosition(offset); got != w {
			t.Fatalf("offset %d: %+v != %+v", offset, got, w)
		}
	}
	l.RewindTo(NewSourcePosition(0, 0, 0))
	if got := l.GetSourcePosition(3); got != want[3] {
		t.Fatalf("after rewind: %+v", got)
	}
}
