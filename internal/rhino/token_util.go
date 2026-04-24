// Port of com/google/javascript/rhino/TokenUtil.java @ closure-compiler v20260415

package rhino

import "unicode"

// IsJSSpace reports whether c is a JavaScript whitespace character as defined
// by ECMA. Mirrors TokenUtil.isJSSpace(int).
func IsJSSpace(c rune) bool {
	if c <= 127 {
		return c == 0x20 || c == 0x9 || c == 0xC || c == 0xB
	}
	return c == 0xA0 || unicode.Is(unicode.Zs, c)
}

// IsJSFormatChar reports whether c is a Unicode format character.
// Mirrors TokenUtil.isJSFormatChar(int).
func IsJSFormatChar(c rune) bool {
	return c > 127 && unicode.Is(unicode.Cf, c)
}

// IsWhitespace reports whether c is a whitespace character.
// Mirrors TokenUtil.isWhitespace(int).
func IsWhitespace(c rune) bool {
	return unicode.IsSpace(c)
}

// IsStrWhiteSpaceChar classifies a character as a string whitespace character.
// Returns TriUnknown for U+000B (VT, IE differs from ECMAScript),
// TriTrue for standard whitespace characters, TriFalse otherwise.
// Mirrors TokenUtil.isStrWhiteSpaceChar(int).
func IsStrWhiteSpaceChar(c rune) Tri {
	switch c {
	case 0x000B: // <VT>
		return TriUnknown // IE says "no", ECMAScript says "yes"
	case 0x0020, // <SP>
		'\n',   // <LF>
		'\r',   // <CR>
		'\t',   // <TAB>
		0x00A0, // <NBSP>
		0x000C, // <FF>
		0x2028, // <LS>
		0x2029, // <PS>
		0xFEFF: // <BOM>
		return TriTrue
	default:
		if unicode.Is(unicode.Zs, c) {
			return TriTrue
		}
		return TriFalse
	}
}
