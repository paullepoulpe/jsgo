package sourcemap

import (
	"math"
	"strconv"
	"testing"
)

func TestSourceMapJsonLexer_StreamStringsAndNumbers(t *testing.T) {
	var lex SourceMapJsonLexer
	lex.initStream(`[0, -42, 2147483647, -2147483648, "hello\n\u0041\u03b1\ud83d\ude00"]`)
	if err := lex.beginArrayStream(); err != nil {
		t.Fatal(err)
	}
	wants := []string{
		"0",
		"-42",
		strconv.Itoa(math.MaxInt32),
		strconv.Itoa(math.MinInt32),
		"hello\nA\u03b1\U0001f600",
	}
	for _, want := range wants {
		got, err := lex.nextStringStream()
		if err != nil {
			t.Fatalf("nextStringStream() unexpected error: %v", err)
		}
		if got != want {
			t.Fatalf("nextStringStream() = %q, want %q", got, want)
		}
	}
	if err := lex.endArrayStream(); err != nil {
		t.Fatal(err)
	}

	for _, bad := range []string{`"\u12g4"`, `"\u12"`} {
		lex.initStream(bad)
		if _, err := lex.nextStringStream(); err == nil {
			t.Fatalf("nextStringStream(%s) succeeded, want error", bad)
		}
	}
}
