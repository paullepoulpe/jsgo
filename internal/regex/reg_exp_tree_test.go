// Port of test/com/google/javascript/jscomp/regex/RegExpTreeTest.java @ closure-compiler v20260415

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
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

func parseRegExpAndPrintPattern(regex, flags string) string {
	tree := RegExpTree_ParseRegExp(regex, flags)
	sb := jutil.NewStringBuilder()
	tree.AppendSourceCode(sb)
	return sb.String()
}

func exceptionFrom(t *testing.T, regex, flags string) jutil.IRuntimeException {
	t.Helper()
	var thrown jutil.IRuntimeException
	jutil.Try(func() {
		printed := parseRegExpAndPrintPattern(regex, flags)
		t.Fatalf("Expected exception, but none was thrown. Instead got back: %s", printed)
	}, func(e jutil.IRuntimeException) {
		thrown = e
	})
	return thrown
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
