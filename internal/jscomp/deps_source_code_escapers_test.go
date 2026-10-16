// Port of com/google/javascript/jscomp/deps/SourceCodeEscapersTest.java @ closure-compiler v20260415

package jscomp

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

func TestSourceCodeEscapers_GoogModuleCode(t *testing.T) {
	input := "goog.module('Foo');\nclass Foo {}"
	escaped := SourceCodeEscapers_JavascriptEscaper().Escape(input)
	assertEqual(t, "escaped", escaped, "goog.module(\\x27Foo\\x27);\\nclass Foo {}")
}

func TestSourceCodeEscapers_ExplicitCharacterEscaping(t *testing.T) {
	input := []uint16{'\'', '"', '<', '=', '>', '&', '\b', '\t', '\n', '\f', '\r', '\\'}
	escaped := SourceCodeEscapers_JavascriptEscaper().Escape(jutil.StringFromUTF16(input))
	assertEqual(t, "escaped", escaped, "\\x27\\x22\\x3c\\x3d\\x3e\\x26\\b\\t\\n\\f\\r\\\\")
}

func TestSourceCodeEscapers_LowUnicodeValueEscaping(t *testing.T) {
	input := []uint16{0x00, 0x10, 0x1F, 0x20, 0x7E, 0x7F, 0xC0, 0xFF}
	escaped := SourceCodeEscapers_JavascriptEscaper().Escape(jutil.StringFromUTF16(input))
	assertEqual(t, "escaped", escaped, "\\x00\\x10\\x1f ~\\x7f\\xc0\\xff")
}

func TestSourceCodeEscapers_HighUnicodeValueEscaping(t *testing.T) {
	input := []uint16{0x100, 0xF00, 0xFFFF}
	escaped := SourceCodeEscapers_JavascriptEscaper().Escape(jutil.StringFromUTF16(input))
	assertEqual(t, "escaped", escaped, "\\u0100\\u0f00\\uffff")
}
