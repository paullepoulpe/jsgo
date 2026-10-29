// Port of com/google/javascript/jscomp/DiagnosticType.java @ closure-compiler v20261021

/*
 * Copyright 2008 The Closure Compiler Authors.
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

package jscomp

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

func TestDiagnosticType_Format(t *testing.T) {
	cases := []struct {
		pattern string
		args    []string
		want    string
	}{
		{"plain literal message", nil, "plain literal message"},
		{"plain literal message", []string{"unused"}, "plain literal message"},
		{"''", nil, "'"},
		{"''{0}''", nil, "'{0}'"},
		{"''{0}''", []string{"foo"}, "'foo'"},
		{"''{0}''", []string{"x'y"}, "'x'y'"},
		{"'{0}'", []string{"foo"}, "{0}"},
		{"'quoted {0} text", []string{"foo"}, "quoted {0} text"},
		{"'use strict' is unnecessary", nil, "use strict is unnecessary"},
		{"Try var {0} = class '{'...'}' (metadata {1})", []string{"C", "M"}, "Try var C = class {...} (metadata M)"},
		{"Try var {0} = class '{'...'}' (metadata {1})", []string{"C"}, "Try var C = class {...} (metadata {1})"},
		{"{0} and {1} and {2}", []string{"a", "b"}, "a and b and {2}"},
		{"{1} cannot extend {0}s ({0})", []string{"Foo", "Bar"}, "Bar cannot extend Foos (Foo)"},
		{"{00}", []string{"zero"}, "zero"},
		{"{0,}", []string{"zero"}, "zero"},
		{"{0}}", []string{"a"}, "a}"},
		{"é{0}😀{1}", []string{"A", "B"}, "éA😀B"},
		{"{9999}", []string{"a"}, "{9999}"},
	}

	for _, tc := range cases {
		dt := DiagnosticType_Error("JSC_TEST", tc.pattern)
		var got string
		if tc.args == nil {
			got = dt.Format(nil...)
		} else {
			got = dt.Format(tc.args...)
		}
		if got != tc.want {
			t.Errorf("DiagnosticType(%q).Format(%v) = %q, want %q", tc.pattern, tc.args, got, tc.want)
		}
		// Verify cached msgFormat path on second call and Platform_FormatMessage parity.
		if got2 := dt.Format(tc.args...); got2 != tc.want {
			t.Errorf("cached DiagnosticType(%q).Format(%v) = %q, want %q", tc.pattern, tc.args, got2, tc.want)
		}
		if gotPlat := Platform_FormatMessage(tc.pattern, tc.args...); gotPlat != tc.want {
			t.Errorf("Platform_FormatMessage(%q, %v) = %q, want %q", tc.pattern, tc.args, gotPlat, tc.want)
		}
	}

	// Also verify explicit zero-arg call t.Format().
	dt := DiagnosticType_Warning("JSC_NO_ARGS", "missing ''{0}'' and {1}")
	if got, want := dt.Format(), "missing '{0}' and {1}"; got != want {
		t.Errorf("dt.Format() = %q, want %q", got, want)
	}
}

func TestDiagnosticType_FormatInvalidPatterns(t *testing.T) {
	badPatterns := []struct {
		pattern string
		wantMsg string
	}{
		{"{0", "Unmatched braces in the pattern."},
		{"{x}", "can't parse argument number: x"},
		{"{ 0}", "can't parse argument number:  0"},
		{"{+0}", "can't parse argument number: +0"},
		{"{{0}}", "can't parse argument number: {0}"},
		{"{-1}", "negative argument number: -1"},
		{"{10000}", "10000 exceeds the ArgumentIndex implementation limit"},
	}
	for _, tc := range badPatterns {
		var gotMsg string
		jutil.Try(func() {
			DiagnosticType_Error("JSC_BAD", tc.pattern).Format("a")
		}, func(e jutil.IIllegalArgumentException) {
			gotMsg = e.GetMessage()
		})
		if gotMsg != tc.wantMsg {
			t.Errorf("DiagnosticType(%q).Format: got exception %q, want %q", tc.pattern, gotMsg, tc.wantMsg)
		}
	}
}
