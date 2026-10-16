// Port of test/com/google/debugging/sourcemap/UtilTest.java @ closure-compiler v20260415

package sourcemap

import (
	"strings"
	"testing"
)

func utilTestAppendHexJavaScriptRepresentation(t *testing.T, ch uint16, expectedOut string) {
	t.Helper()
	var sb strings.Builder
	Util_AppendHexJavaScriptRepresentation(&sb, ch)
	if sb.String() != expectedOut {
		t.Errorf("appendHexJavaScriptRepresentation(%q) = %q, want %q", rune(ch), sb.String(), expectedOut)
	}
}

func TestAppendHexJavaScriptRepresentation(t *testing.T) {
	utilTestAppendHexJavaScriptRepresentation(t, 'a', "\\u0061")
	utilTestAppendHexJavaScriptRepresentation(t, 'z', "\\u007a")
	utilTestAppendHexJavaScriptRepresentation(t, 0, "\\u0000")
	utilTestAppendHexJavaScriptRepresentation(t, 0xa1, "\\u00a1")
}
