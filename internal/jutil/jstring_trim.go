package jutil

// Go-only: language-adapter — more java.lang.String methods on WTF-8 strings.

// StringTrim is Java's String.trim(): strips leading and trailing chars
// <= U+0020. Those are all ASCII, so trimming bytes of the WTF-8 form is
// exact (no multi-byte sequence contains a byte <= 0x20).
func StringTrim(s string) string {
	start, end := 0, len(s)
	for start < end && s[start] <= ' ' {
		start++
	}
	for start < end && s[end-1] <= ' ' {
		end--
	}
	return s[start:end]
}
