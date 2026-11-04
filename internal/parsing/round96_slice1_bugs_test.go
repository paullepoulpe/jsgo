package parsing

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

func TestRound96JsDocParseTypeNameNilNode(t *testing.T) {
	const wantMsg = "Bad type annotation. expecting a type name. See https://github.com/google/closure-compiler/wiki/Annotating-JavaScript-for-the-Closure-Compiler for more information."

	pt := newParserTest(t)
	pt.parseWarning("/** @extends {<T>} */ class C {}", wantMsg)
	pt.parseWarning("/** @implements {<T>} */ class C {}", wantMsg)

	jt := newJsDocInfoParserTest(t)
	jt.parse("@extends {<T>} */", wantMsg)
}

func TestRound96BracedUnicodeCodePointOverflow(t *testing.T) {
	pt := newParserTest(t)
	pt.parseError("var i = '\\u{80000000}'", "Undefined Unicode code-point")
	pt.parseError("var i = '\\u{100000000}'", "Undefined Unicode code-point")
	pt.parseError("var i = '\\u{ffffffffffffffffffff}'", "Undefined Unicode code-point")
	pt.parseError("var i = `\\u{80000000}`", "Undefined Unicode code-point")
	pt.parseError("var i = `\\u{100000000}`", "Undefined Unicode code-point")
	pt.parse("var ok = '\\u{00000000000000000061}'")
}

func TestRound96ParseUTF16HexCodePoint(t *testing.T) {
	cases := []struct {
		hex  string
		want int32
	}{
		{"00000000000000000061", 0x61},
		{"10ffff", 0x10ffff},
		{"110000", 0x110000},
		{"80000000", 0x110000},
		{"100000000", 0x110000},
		{"ffffffffffffffffffff", 0x110000},
	}
	for _, tc := range cases {
		units := jutil.StringToUTF16(tc.hex)
		if got := parseUTF16HexCodePoint(units, 0, len(units)); got != tc.want {
			t.Errorf("parseUTF16HexCodePoint(%q) = 0x%x, want 0x%x", tc.hex, got, tc.want)
		}
	}

	for _, bad := range []string{"", "100000000g", "g"} {
		units := jutil.StringToUTF16(bad)
		func() {
			defer func() {
				r := recover()
				if _, ok := r.(jutil.INumberFormatException); !ok {
					t.Errorf("parseUTF16HexCodePoint(%q) panic = %T (%v), want INumberFormatException", bad, r, r)
				}
			}()
			parseUTF16HexCodePoint(units, 0, len(units))
		}()
	}
}
