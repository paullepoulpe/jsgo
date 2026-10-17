// Port of test/com/google/debugging/sourcemap/UtilTest.java @ closure-compiler v20260415

package sourcemap_test

import (
	"strings"
	"testing"

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
}
