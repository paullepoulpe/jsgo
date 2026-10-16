// Port of com/google/javascript/rhino/TokenUtil.java @ closure-compiler v20260415

package rhino

import (
	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// TokenUtil holds helper methods for parsing JavaScript.
type TokenUtil struct{}

// TokenUtil_IsJSSpace is as defined in ECMA. jsscan.c uses C isspace()
// (which allows \v, I think.) note that code in getChar() implicitly accepts
// '\r' == \u000D as well.
//
// Java: TokenUtil.isJSSpace(int)
func TokenUtil_IsJSSpace(c int) bool {
	if c <= 127 {
		return c == 0x20 || c == 0x9 || c == 0xC || c == 0xB
	} else {
		return c == 0xA0 ||
			jutil.CharacterGetType(uint16(c)) == jutil.Character_SPACE_SEPARATOR
	}
}

// Java: TokenUtil.isJSFormatChar(int)
func TokenUtil_IsJSFormatChar(c int) bool {
	return c > 127 && jutil.CharacterGetTypeOfCodePoint(c) == jutil.Character_FORMAT
}

// Java: TokenUtil.isWhitespace(int)
func TokenUtil_IsWhitespace(c int) bool {
	return jutil.CharacterIsWhitespaceCodePoint(c)
}

// TokenUtil_IsStrWhiteSpaceChar is copied from Rhino's ScriptRuntime.
//
// Java: TokenUtil.isStrWhiteSpaceChar(int)
func TokenUtil_IsStrWhiteSpaceChar(c int) base.Tri {
	switch c {
	case 0x0B: // <VT>
		return base.Tri_UNKNOWN // IE says "no", ECMAScript says "yes"
	case ' ', // <SP>
		'\n',     // <LF>
		'\r',     // <CR>
		'\t',     // <TAB>
		0xA0,     // <NBSP>
		'\u000C', // <FF>
		0x2028,   // <LS>
		0x2029,   // <PS>
		0xFEFF:   // <BOM>
		return base.Tri_TRUE
	default:
		if jutil.CharacterGetTypeOfCodePoint(c) == jutil.Character_SPACE_SEPARATOR {
			return base.Tri_TRUE
		}
		return base.Tri_FALSE
	}
}

// Java: TokenUtil()
func NewTokenUtil() *TokenUtil { return &TokenUtil{} }
