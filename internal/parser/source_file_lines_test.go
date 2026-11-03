package parser

import (
	"strings"
	"testing"
)

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

func TestSourceFileNegativeOffsetDoesNotPanicOrPoisonLastLine(t *testing.T) {
	f := NewSourceFile("neg.js", "abc\ndef\nghi")

	// Call LineColumn(-1) before lineStarts is even built.
	line, col := f.LineColumn(-1)
	if line != 0 || col != -1 {
		t.Fatalf("LineColumn(-1) = (%d, %d), want (0, -1)", line, col)
	}
	if last := f.lastLine.Load(); last < 0 {
		t.Fatalf("LineColumn(-1) poisoned f.lastLine = %d", last)
	}

	// Move cursor forward to line 2 ("ghi" at offset 8), then call LineColumn(-1)
	// and verify f.lastLine is neither poisoned nor out of bounds for subsequent lookups.
	if got := f.LinePosition(NewSourcePosition(9)); got != (LinePosition{Offset: 9, Line: 2, Column: 1}) {
		t.Fatalf("unexpected LinePosition(9): %+v", got)
	}
	if l2, c2 := f.LineColumn(-42); l2 != 0 || c2 != -42 {
		t.Fatalf("LineColumn(-42) = (%d, %d), want (0, -42)", l2, c2)
	}
	if last := f.lastLine.Load(); last != 2 {
		t.Fatalf("expected f.lastLine preserved at 2 after negative lookup, got %d", last)
	}
	if got := f.LinePosition(NewSourcePosition(5)); got != (LinePosition{Offset: 5, Line: 1, Column: 1}) {
		t.Fatalf("unexpected LinePosition(5) after negative lookup: %+v", got)
	}
	if gotLine := FindLine(f.LineStarts(), -1); gotLine != 0 {
		t.Fatalf("FindLine(starts, -1) = %d, want 0", gotLine)
	}
	if gotEmpty := FindLine(nil, -1); gotEmpty != 0 {
		t.Fatalf("FindLine(nil, -1) = %d, want 0", gotEmpty)
	}
}

func TestScannerTemplateErrorOverflowBeyondUint16(t *testing.T) {
	// Generate 65,540 template tokens, each containing an invalid octal escape
	// `\1` at offset start+1 so every token records an entry in Scanner.templateErrors.
	// Previously:
	//   - token #65,535 (1-based 65,536 == 0x10000) wrapped Token.aux to 0 (losing its error!)
	//   - token #65,536 (1-based 65,537 == 0x10001) wrapped Token.aux to 1 (aliasing token #0!)
	const count = 65540
	const unit = "`\\1`;"
	src := strings.Repeat(unit, count)

	sf := NewSourceFile("many_templates.js", src)
	reporter := &recordingErrorReporter{}
	sc := NewScanner(reporter, nil, sf, 0)

	seen := 0
	for {
		tok := sc.NextToken()
		if tok.Type() == TokenType_END_OF_FILE {
			break
		}
		if tok.Type() != TokenType_NO_SUBSTITUTION_TEMPLATE {
			continue
		}
		errInfo := sc.TemplateError(tok)
		// In "`\1`;" (len 5), SkipTemplateLiteralEscapeSequence records s.GetPosition()
		// right after consuming `\1` (offset i*5 + 3), while tok.Location.Start.Offset is i*5.
		if !errInfo.HasMsg || errInfo.Message == "" {
			t.Fatalf("template token #%d (offset %d, aux=%#x) lost its TemplateTokenError: %+v",
				seen, tok.Location.Start.Offset, tok.aux, errInfo)
		}
		wantPosOffset := int32(seen*len(unit) + 3)
		if errInfo.Position.Offset != wantPosOffset {
			t.Fatalf("template token #%d aliased wrong position: got %d, want %d (aux=%#x)",
				seen, errInfo.Position.Offset, wantPosOffset, tok.aux)
		}
		if seen == 65534 || seen == 65535 || seen == 65536 {
			if tok.aux != tokenAuxOverflow {
				t.Fatalf("template token #%d expected aux == %#x, got %#x", seen, tokenAuxOverflow, tok.aux)
			}
		}
		seen++
	}
	if seen != count {
		t.Fatalf("scanned %d template tokens, want %d", seen, count)
	}
}
