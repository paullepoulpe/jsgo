// Port of test/com/google/debugging/sourcemap/Base64Test.java @ closure-compiler v20260415

package sourcemap

import "testing"

func TestBase64(t *testing.T) {
	for i := int32(0); i < 64; i++ {
		base64TestValue(t, i)
	}
}

func TestBase64EncodeInt(t *testing.T) {
	for _, tc := range []struct {
		v    int32
		want string
	}{
		{0, "AAAAAA"},
		{1, "AAAAAQ"},
		{42, "AAAAKg"},
		{-100, "////nA"},
		{-1, "/////w"}, // 0xffffffff
	} {
		if got := Base64_Base64EncodeInt(tc.v); got != tc.want {
			t.Errorf("base64EncodeInt(%d) = %q, want %q", tc.v, got, tc.want)
		}
	}
}

func base64TestValue(t *testing.T, value int32) {
	t.Helper()
	if got := Base64_FromBase64(Base64_ToBase64(value)); got != value {
		t.Errorf("fromBase64(toBase64(%d)) = %d", value, got)
	}
}
