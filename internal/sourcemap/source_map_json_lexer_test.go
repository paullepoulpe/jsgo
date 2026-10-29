package sourcemap

import (
	"math"
	"testing"
)

func TestSourceMapJsonLexer_ParseIntAndUnescape(t *testing.T) {
	cases := []struct {
		in    string
		radix int
		want  int32
		ok    bool
	}{
		{"0", 10, 0, true},
		{"-0", 10, 0, true},
		{"+42", 10, 42, true},
		{"2147483647", 10, math.MaxInt32, true},
		{"-2147483648", 10, math.MinInt32, true},
		{"2147483648", 10, 0, false},
		{"-2147483649", 10, 0, false},
		{"99999999999999999999", 10, 0, false},
		{"-", 10, 0, false},
		{"+", 10, 0, false},
		{"", 10, 0, false},
		{"12a", 10, 0, false},
		{"\u0661\u0662", 10, 12, true},
		{"007f", 16, 0x7f, true},
		{"ABCD", 16, 0xabcd, true},
		{"ffff", 16, 0xffff, true},
		{"12g4", 16, 0, false},
	}
	for _, tc := range cases {
		got, ok := sourceMapJsonLexer_parseInt(tc.in, 0, len(tc.in), tc.radix)
		if ok != tc.ok || got != tc.want {
			t.Errorf("parseInt(%q, %d) = (%d, %v), want (%d, %v)", tc.in, tc.radix, got, ok, tc.want, tc.ok)
		}
	}

	lex := NewSourceMapJsonLexer("3 -42 2147483647 -2147483648")
	for _, want := range []int{3, -42, math.MaxInt32, math.MinInt32} {
		if got := lex.NextInt(); got != want {
			t.Fatalf("NextInt() = %d, want %d", got, want)
		}
	}

	assertParseException := func(fn func()) {
		t.Helper()
		defer func() {
			r := recover()
			if r == nil {
				t.Fatal("expected SourceMapParseException panic")
			}
			if _, ok := r.(*SourceMapParseException); !ok {
				t.Fatalf("expected *SourceMapParseException, got %T: %v", r, r)
			}
		}()
		fn()
	}

	assertParseException(func() { NewSourceMapJsonLexer("-").NextInt() })
	assertParseException(func() { NewSourceMapJsonLexer("2147483648").NextInt() })
	assertParseException(func() { NewSourceMapJsonLexer("-2147483649").NextInt() })

	strLex := NewSourceMapJsonLexer(`"hello\n\u0041\u03b1\ud83d\ude00"`)
	if got := strLex.NextString(); got != "hello\nA\u03b1\U0001f600" {
		t.Fatalf("NextString() = %q, want %q", got, "hello\nA\u03b1\U0001f600")
	}
	assertParseException(func() { NewSourceMapJsonLexer(`"\u12g4"`).NextString() })
	assertParseException(func() { NewSourceMapJsonLexer(`"\u12"`).NextString() })
}
