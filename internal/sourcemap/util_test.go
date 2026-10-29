// Port of test/com/google/debugging/sourcemap/UtilTest.java @ closure-compiler v20261021

/*
 * Copyright 2017 The Closure Compiler Authors.
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

package sourcemap_test

import (
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/sourcemap"
)

func utilTestAppendHexJavaScriptRepresentation(t *testing.T, ch uint16, expectedOut string) {
	t.Helper()
	var sb strings.Builder
	sourcemap.Util_AppendHexJavaScriptRepresentation(&sb, ch)
	if sb.String() != expectedOut {
		t.Errorf("appendHexJavaScriptRepresentation(%q) = %q, want %q", rune(ch), sb.String(), expectedOut)
	}
}

func TestAppendHexJavaScriptRepresentation(t *testing.T) {
	utilTestAppendHexJavaScriptRepresentation(t, 'a', "\\u0061")
	utilTestAppendHexJavaScriptRepresentation(t, 'z', "\\u007a")
	utilTestAppendHexJavaScriptRepresentation(t, 0, "\\u0000")
	utilTestAppendHexJavaScriptRepresentation(t, 0xa1, "\\u00a1")

	var sb strings.Builder
	sourcemap.Util_AppendHexJavaScriptRepresentationToAppendable(0x1f600, &sb)
	if sb.String() != `\ud83d\ude00` {
		t.Errorf("supplementary append = %q, want %q", sb.String(), `\ud83d\ude00`)
	}
}

func TestEscapeStringDirectUTF8(t *testing.T) {
	loneHigh := jutil.StringFromUTF16([]uint16{'a', 0xd800, 'b'})
	loneLow := jutil.StringFromUTF16([]uint16{0xdfff})
	cases := []struct {
		in   string
		want string
	}{
		{"plain ascii 'ok'", `"plain ascii 'ok'"`},
		{"a\nb\rc\td\\e\"f\x00\x1f\x7f", `"a\nb\rc\td\\e\"f\u0000\u001f` + "\x7f" + `"`},
		{"a > b && c-->d && e]]>f", `"a > b && c--\u003ed && e]]\u003ef"`},
		{"</script> </ScRiPt> <!-- <div", `"\u003c/script> \u003c/ScRiPt> \u003c!-- <div"`},
		{"</\u017fcript>", `"\u003c/\u017fcript>"`},
		{"α😀Ω", `"\u03b1\ud83d\ude00\u03a9"`},
		{loneHigh, `"a\ud800b"`},
		{loneLow, `"\udfff"`},
	}
	for _, tc := range cases {
		if got := sourcemap.Util_EscapeString(tc.in); got != tc.want {
			t.Errorf("Util_EscapeString(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}
