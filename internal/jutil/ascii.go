package jutil

// Go-only: language-adapter — com.google.common.base.Ascii.

// asciiCaseMask is Ascii.CASE_MASK.
const asciiCaseMask = 0x20

// AsciiEqualsIgnoreCase is Ascii.equalsIgnoreCase(CharSequence,
// CharSequence): equality of UTF-16 code units, where only the ASCII letters
// a-z and A-Z are compared case-insensitively.
func AsciiEqualsIgnoreCase(s1, s2 string) bool {
	if s1 == s2 {
		return true
	}
	u1 := StringToUTF16(s1)
	u2 := StringToUTF16(s2)
	if len(u1) != len(u2) {
		return false
	}
	for i := range u1 {
		c1 := u1[i]
		c2 := u2[i]
		if c1 == c2 {
			continue
		}
		alphaIndex := asciiGetAlphaIndex(c1)
		if alphaIndex < 26 && alphaIndex == asciiGetAlphaIndex(c2) {
			continue
		}
		return false
	}
	return true
}

// asciiGetAlphaIndex is Ascii.getAlphaIndex(char): folds upper-case ASCII to
// lower-case and makes it zero-indexed and unsigned (by casting to char).
func asciiGetAlphaIndex(c uint16) uint16 {
	return (c | asciiCaseMask) - 'a'
}
