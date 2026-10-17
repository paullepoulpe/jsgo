package jscomp

// Go-only: Compiler.CodeBuilder's line and column bookkeeping, which counts
// UTF-16 chars as Java's String does, and stays linear in the size of the
// appended code. It was quadratic in its number of lines: in the held-out
// sweep, pretty-printing webpack's 1.4 MB lib/javascript/syntax-printer.js
// took 47 s against Java's 31 s, and 7 s once fixed.

import (
	"strings"
	"testing"
)

func TestCompilerCodeBuilder_LineAndColumn(t *testing.T) {
	b := NewCompiler_CodeBuilder()
	check := func(line, col int) {
		t.Helper()
		if b.GetLineIndex() != line || b.GetColumnIndex() != col {
			t.Fatalf("line %d col %d, want line %d col %d", b.GetLineIndex(), b.GetColumnIndex(), line, col)
		}
	}
	check(0, 0)
	b.Append("var a;")
	check(0, 6)
	b.Append("é𝒳") // one char, then two (a surrogate pair)
	check(0, 9)
	b.Append("\n")
	check(1, 0)
	b.Append("x\n\nyz")
	check(3, 2)
	b.Append("𝒳\nα𝒳") // the last line: 1 + 2 chars
	check(4, 3)
	b.Reset()
	check(4, 3) // reset leaves the counts unchanged
	b.Append("q")
	check(4, 4)
}

func TestCompilerCodeBuilder_ManyLines(t *testing.T) {
	const lines = 200000
	code := strings.Repeat("  f(\"é\");\n", lines) + "end;"
	b := NewCompiler_CodeBuilder()
	b.Append(code)
	if b.GetLineIndex() != lines || b.GetColumnIndex() != 4 {
		t.Fatalf("line %d col %d, want line %d col 4", b.GetLineIndex(), b.GetColumnIndex(), lines)
	}
	if b.GetLength() != lines*10+4 {
		t.Fatalf("length %d, want %d", b.GetLength(), lines*10+4)
	}
}
