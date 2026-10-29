// Port of test/com/google/javascript/jscomp/regex/RegExpTreeTest.java @ closure-compiler v20261021

/*
 * Copyright 2014 The Closure Compiler Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package regex

import (
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

func parseRegExpAndPrintPattern(regex, flags string) string {
	tree, err := ParseRegExp(regex, flags)
	if err != nil {
		panic(err.toThrowable())
	}
	var sb strings.Builder
	tree.AppendSourceCode(&sb)
	return sb.String()
}

func exceptionFrom(t *testing.T, regex, flags string) jutil.IRuntimeException {
	t.Helper()
	tree, err := ParseRegExp(regex, flags)
	if err != nil {
		return err.toThrowable()
	}
	var sb strings.Builder
	tree.AppendSourceCode(&sb)
	t.Fatalf("Expected exception, but none was thrown. Instead got back: %s", sb.String())
	return nil
}

func assertRegexCompilesTo(t *testing.T, regex, flags, expected string) {
	t.Helper()
	if got := parseRegExpAndPrintPattern(regex, flags); got != expected {
		t.Errorf("parse(%q, %q) = %q, want %q", regex, flags, got, expected)
	}
}

func assertRegexThrowsExceptionWithMessage(t *testing.T, regex, flags, message string) {
	t.Helper()
	if got := exceptionFrom(t, regex, flags).GetMessage(); got != message {
		t.Errorf("parse(%q, %q) threw %q, want %q", regex, flags, got, message)
	}
}

func assertRegexCompilesToSame(t *testing.T, regex, flags string) {
	t.Helper()
	assertRegexCompilesTo(t, regex, flags, regex)
}

func TestValidEs2018LookbehindAssertions(t *testing.T) {
	assertRegexCompilesToSame(t, "(?<=asdf)", "")
	assertRegexCompilesToSame(t, "(?<!asdf)", "")
	assertRegexCompilesToSame(t, "(?<=(?<!asdf))", "")
	assertRegexCompilesToSame(t, "(?<=(?<!asdf))", "")
}

func TestInvalidEs2018LookbehindAssertions(t *testing.T) {
	assertRegexThrowsExceptionWithMessage(t, "(?<asdf)", "", "Invalid capture group name: <asdf)")
}

func TestValidEs2018UnicodePropertyEscapes(t *testing.T) {
	assertRegexCompilesToSame(t, "\\p{Script=Greek}", "u")
	assertRegexCompilesToSame(t, "\\P{Script=Greek}", "u")
	assertRegexCompilesToSame(t, "\\p{Letter}", "u")
	assertRegexCompilesToSame(t, "\\P{Letter}", "u")
	assertRegexCompilesTo(t, "\\p", "", "p") // Without 'u' flag, '\p' is just 'p'
	assertRegexCompilesTo(t, "\\P", "", "P")
}

func TestInvalidEs2018UnicodePropertyEscapes(t *testing.T) {
	assertRegexThrowsExceptionWithMessage(t, "\\p{", "u",
		"Malformed Unicode Property Escape: expected '=' or '}' after \\p{")

	assertRegexThrowsExceptionWithMessage(t, "\\P{", "u",
		"Malformed Unicode Property Escape: expected '=' or '}' after \\P{")

	assertRegexThrowsExceptionWithMessage(t, "\\p{=Greek}", "u",
		"if '=' is present in a unicode property escape, the name cannot be empty")

	assertRegexThrowsExceptionWithMessage(t, "\\P{=Greek}", "u",
		"if '=' is present in a unicode property escape, the name cannot be empty")

	assertRegexThrowsExceptionWithMessage(t, "\\p{}", "u",
		"unicode property escape value cannot be empty")

	assertRegexThrowsExceptionWithMessage(t, "\\P{}", "u",
		"unicode property escape value cannot be empty")
}

func TestValidEs2018RegexNamedCaptureGroups(t *testing.T) {
	assertRegexCompilesToSame(t, "(?<name>)", "")
	assertRegexCompilesToSame(t, "(?<h$h1h_>)", "u")
	assertRegexCompilesToSame(t, "(?<$var_name>blah)", "")
	assertRegexCompilesToSame(t, "(?<_var_name>>>>)", "")
}

func TestInvalidEs2018RegexNamedCaptureGroups(t *testing.T) {
	assertRegexThrowsExceptionWithMessage(t, "(?<name)", "", "Invalid capture group name: <name)")
	assertRegexThrowsExceptionWithMessage(t, "(?<1b>)", "", "Invalid capture group name: <1b>)")
	assertRegexThrowsExceptionWithMessage(t, "(?<>)", "", "Invalid capture group name: <>)")
	assertRegexThrowsExceptionWithMessage(t, "(?<.name>)", "", "Invalid capture group name: <.name>)")
}

func TestNumCapturingGroups(t *testing.T) {
	assertRegexCompilesToSame(t, "(h(i))\\2", "")
	// TODO(b/116048051): reference to non-existent capture group should be an error.
	assertRegexCompilesTo(t, "(h(i))\\3", "", "(h(i))\\x03")

	assertRegexCompilesToSame(t, "(?<foo>.*(?<bar>))", "")
}

func TestValidEs2018CaptureNameBackreferencing(t *testing.T) {
	assertRegexCompilesToSame(t, "(?<name>)\\k<name>", "")
	// Note that (?: ) only used for printing purposes to
	// indicate the nesting structure of Concatenation nodes.
	// It is not actually what the compiler prints out as source code.
	assertRegexCompilesTo(t,
		"(?<foo>(?<bar>))\\k<foo>\\k<bar>", "", "(?:(?<foo>(?<bar>))\\k<foo>)\\k<bar>")
	assertRegexCompilesTo(t,
		"(?<foo>(?<bar>)\\k<bar>)\\k<foo>", "", "(?<foo>(?<bar>)\\k<bar>)\\k<foo>")

	// The below examples where the backreference comes before the definition of named groups is
	// allowed syntactically, although it is not able to reference the original group semantically
	assertRegexCompilesToSame(t, "\\k<foo>(?<foo>)", "")
	assertRegexCompilesToSame(t, "\\k<foo>(?<foo>\\k<bar>(?<bar>))", "")

	// Backreferencing the name in the group it is defined is also allowed
	assertRegexCompilesToSame(t, "(?<foo>\\k<foo>)", "")
}

func TestInvalidEs2018CaptureNameBackreferencing(t *testing.T) {
	assertRegexThrowsExceptionWithMessage(t, "(?<foo>)\\k<bar>", "",
		"Invalid named capture referenced: \\k<bar>")
	assertRegexThrowsExceptionWithMessage(t, "(?<foo>)\\k<foo", "",
		"Malformed named capture group: <foo")
	assertRegexThrowsExceptionWithMessage(t, "\\k<1b>(?<foo>)", "",
		"Invalid capture group name: <1b>(?<foo>)")

	// Even though enclosed in (?<>), 'foo' not a capture name definition
	assertRegexThrowsExceptionWithMessage(t, "[(?<foo>)](?<bar>)\\k<foo>", "",
		"Invalid named capture referenced: \\k<foo>")
}

func TestBackreferencingTreatedAsStringIfNoGroup(t *testing.T) {
	// Backreferencing without named group definitions is just treated as normal string
	assertRegexCompilesTo(t, "\\k<foo>", "", "k<foo>")

	// Note that the (?: ) that is wrapped around "k<" at expected output is only used for printing
	// purposes to indicate the nesting structure of Concatenation nodes.
	// It is not actually what the compiler prints out as source code.
	assertRegexCompilesTo(t, "\\k<.", "", "(?:k<).")

	// Even though enclosed in (?<>), 'foo' not a capture name definition
	// (?: ) in expected output serves same purpose as above test
	assertRegexCompilesTo(t, "[(?<foo>)]\\k<foo>", "", "(?:[()<>?fo]k)<foo>")
}

func TestValidUnicodeEscape(t *testing.T) {
	assertRegexCompilesTo(t, "\\u0061", "", "a")
	assertRegexCompilesTo(t, "\\u10b1", "u", "\\u10b1")
	assertRegexCompilesTo(t, "\\u{61}", "u", "a")
	assertRegexCompilesTo(t, "\\u{10b1}", "u", "\\u10b1")
	assertRegexCompilesTo(t, "\\u{1bc}", "u", "\\u01bc")
	assertRegexCompilesTo(t, "\\u{100A3}", "u", "\\ud800\\udca3")
}

func TestInvalidUnicodeEscape(t *testing.T) {
	assertRegexThrowsExceptionWithMessage(t, "\\u{a012", "u",
		"Malformed unicode escape: expected '}' after {a012")
	assertRegexThrowsExceptionWithMessage(t, "\\u{}", "u", //
		"Empty unicode escape")
	assertRegexThrowsExceptionWithMessage(t, "\\u{10za}", "u", //
		"za}")
	assertRegexThrowsExceptionWithMessage(t, "\\u{FFFFFF}", "u",
		"Unicode must be at most 0x10FFFF: FFFFFF")
	assertRegexThrowsExceptionWithMessage(t, "\\u{FF00FFFF}", "u",
		"Cannot parse hexadecimal encoding wider than 28 bits: FF00FFFF")
}

// Java: RegExpTreeTest.testDecomposedCharsetEqualsAndHashCode()
func TestDecomposedCharsetEqualsAndHashCode(t *testing.T) {
	empty1 := NewRegExpTree_DecomposedCharset(false, CharRanges_EMPTY, "")
	empty2 := NewRegExpTree_DecomposedCharset(false, CharRanges_EMPTY, "")
	invertedEmpty := NewRegExpTree_DecomposedCharset(true, CharRanges_EMPTY, "")

	rangesA := CharRanges_WithMembers('a', 'b')
	rangesB := CharRanges_WithMembers('a', 'c')
	charsetA1 := NewRegExpTree_DecomposedCharset(false, rangesA, "\\d")
	charsetA2 := NewRegExpTree_DecomposedCharset(false, rangesA, "\\d")
	charsetAInverted := NewRegExpTree_DecomposedCharset(true, rangesA, "\\d")
	charsetDiffRanges := NewRegExpTree_DecomposedCharset(false, rangesB, "\\d")
	charsetDiffGroups := NewRegExpTree_DecomposedCharset(false, rangesA, "\\s")

	isEqualTo := func(what string, a, b *RegExpTree_DecomposedCharset) {
		t.Helper()
		if !a.Equals(b) {
			t.Errorf("%s: not equal", what)
		}
	}
	isNotEqualTo := func(what string, a, b *RegExpTree_DecomposedCharset) {
		t.Helper()
		if a.Equals(b) {
			t.Errorf("%s: equal", what)
		}
	}
	sameHash := func(what string, a, b *RegExpTree_DecomposedCharset) {
		t.Helper()
		if a.HashCode() != b.HashCode() {
			t.Errorf("%s: hash codes %d and %d", what, a.HashCode(), b.HashCode())
		}
	}

	// Symmetric equality
	isEqualTo("empty1, empty2", empty1, empty2)
	isEqualTo("empty2, empty1", empty2, empty1)
	sameHash("empty1, empty2", empty1, empty2)

	isEqualTo("charsetA1, charsetA2", charsetA1, charsetA2)
	isEqualTo("charsetA2, charsetA1", charsetA2, charsetA1)
	sameHash("charsetA1, charsetA2", charsetA1, charsetA2)

	// Inverted true equality
	charsetAInverted2 := NewRegExpTree_DecomposedCharset(true, rangesA, "\\d")
	isEqualTo("charsetAInverted, charsetAInverted2", charsetAInverted, charsetAInverted2)
	sameHash("charsetAInverted, charsetAInverted2", charsetAInverted, charsetAInverted2)

	// Inequality across inverted, ranges, and namedGroups
	isNotEqualTo("empty1, invertedEmpty", empty1, invertedEmpty)
	isNotEqualTo("invertedEmpty, empty1", invertedEmpty, empty1)
	isNotEqualTo("charsetA1, charsetAInverted", charsetA1, charsetAInverted)
	isNotEqualTo("charsetAInverted, charsetA1", charsetAInverted, charsetA1)
	isNotEqualTo("charsetA1, charsetDiffRanges", charsetA1, charsetDiffRanges)
	isNotEqualTo("charsetA1, charsetDiffGroups", charsetA1, charsetDiffGroups)
	if charsetA1 == nil {
		t.Error("charsetA1 is null")
	}
}

func TestParseRegExpReturnsErrorWithoutPanic(t *testing.T) {
	cases := []struct {
		pattern string
		flags   string
		kind    RegExpParseErrorKind
		wantMsg string
	}{
		{"foo\\", "", RegExpParseErrorStringIndexOutOfBounds, "Index 4 out of bounds for length 4"},
		{"[\\", "", RegExpParseErrorStringIndexOutOfBounds, "Index 2 out of bounds for length 2"},
		{"(?<", "", RegExpParseErrorStringIndexOutOfBounds, "Index 3 out of bounds for length 3"},
		{"[z-a]", "", RegExpParseErrorIndexOutOfBounds, "122 > 97"},
		{"(foo", "", RegExpParseErrorIllegalArgument, "Unclosed parenthetical group: (foo"},
		{"foo)", "", RegExpParseErrorIllegalArgument, ")"},
		{"\\p{}", "u", RegExpParseErrorIllegalArgument, "unicode property escape value cannot be empty"},
		{"\\p{=Greek}", "u", RegExpParseErrorIllegalArgument, "if '=' is present in a unicode property escape, the name cannot be empty"},
	}
	for _, tc := range cases {
		tree, err := ParseRegExp(tc.pattern, tc.flags)
		if err == nil || tree != nil {
			t.Fatalf("ParseRegExp(%q, %q) = (%v, %v), want error %q", tc.pattern, tc.flags, tree, err, tc.wantMsg)
		}
		if err.Kind != tc.kind || err.Message != tc.wantMsg {
			t.Errorf("ParseRegExp(%q, %q) = (%v, %q), want (%v, %q)", tc.pattern, tc.flags, err.Kind, err.Message, tc.kind, tc.wantMsg)
		}
	}
}
