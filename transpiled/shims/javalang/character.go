// Clean-room specification based on java.lang.Character, Integer, Long, Double Javadoc.

package javalang

import (
	"math"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Character radix limits (Character.MIN_RADIX, Character.MAX_RADIX).
const (
	CharacterMinRadix = 2
	CharacterMaxRadix = 36
)

// CharacterIsLetter implements Character.isLetter(char).
func CharacterIsLetter(ch uint16) bool {
	if ch < 0x80 {
		return (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z')
	}
	if ch >= 0xD800 && ch <= 0xDFFF {
		return false
	}
	return unicode.IsLetter(rune(ch))
}

// CharacterIsDigit implements Character.isDigit(char).
func CharacterIsDigit(ch uint16) bool {
	if ch < 0x80 {
		return ch >= '0' && ch <= '9'
	}
	return unicode.Is(unicode.Nd, rune(ch))
}

// CharacterIsLetterOrDigit implements Character.isLetterOrDigit(char).
func CharacterIsLetterOrDigit(ch uint16) bool {
	return CharacterIsLetter(ch) || CharacterIsDigit(ch)
}

// CharacterIsUpperCase implements Character.isUpperCase(char).
func CharacterIsUpperCase(ch uint16) bool {
	if ch < 0x80 {
		return ch >= 'A' && ch <= 'Z'
	}
	return unicode.In(rune(ch), unicode.Lu, unicode.Other_Uppercase)
}

// CharacterIsLowerCase implements Character.isLowerCase(char).
func CharacterIsLowerCase(ch uint16) bool {
	if ch < 0x80 {
		return ch >= 'a' && ch <= 'z'
	}
	return unicode.In(rune(ch), unicode.Ll, unicode.Other_Lowercase)
}

// CharacterIsWhitespace implements Character.isWhitespace(char).
func CharacterIsWhitespace(ch uint16) bool {
	return CharacterIsWhitespaceCodePoint(int(ch))
}

// CharacterIsWhitespaceCodePoint implements Character.isWhitespace(int).
func CharacterIsWhitespaceCodePoint(codePoint int) bool {
	switch codePoint {
	case '\t', '\n', 0x0B, '\f', '\r', 0x1C, 0x1D, 0x1E, 0x1F:
		return true
	case 0x00A0, 0x2007, 0x202F:
		return false
	}
	if codePoint < 0 || codePoint > unicode.MaxRune {
		return false
	}
	return unicode.In(rune(codePoint), unicode.Zs, unicode.Zl, unicode.Zp)
}

// CharacterToUpperCase implements Character.toUpperCase(char).
func CharacterToUpperCase(ch uint16) uint16 {
	if ch < 0x80 {
		if ch >= 'a' && ch <= 'z' {
			return ch - 'a' + 'A'
		}
		return ch
	}
	if ch >= 0xD800 && ch <= 0xDFFF {
		return ch
	}
	u := unicode.ToUpper(rune(ch))
	if u > 0xFFFF {
		return ch
	}
	return uint16(u)
}

// CharacterToLowerCase implements Character.toLowerCase(char).
func CharacterToLowerCase(ch uint16) uint16 {
	if ch < 0x80 {
		if ch >= 'A' && ch <= 'Z' {
			return ch - 'A' + 'a'
		}
		return ch
	}
	if ch >= 0xD800 && ch <= 0xDFFF {
		return ch
	}
	l := unicode.ToLower(rune(ch))
	if l > 0xFFFF {
		return ch
	}
	return uint16(l)
}

// CharacterDigit implements Character.digit(char, int).
func CharacterDigit(ch uint16, radix int) int {
	if radix < CharacterMinRadix || radix > CharacterMaxRadix {
		return -1
	}
	val := -1
	switch {
	case ch >= '0' && ch <= '9':
		val = int(ch - '0')
	case ch >= 'A' && ch <= 'Z':
		val = int(ch-'A') + 10
	case ch >= 'a' && ch <= 'z':
		val = int(ch-'a') + 10
	case ch >= 0xFF21 && ch <= 0xFF3A:
		val = int(ch-0xFF21) + 10
	case ch >= 0xFF41 && ch <= 0xFF5A:
		val = int(ch-0xFF41) + 10
	case CharacterIsDigit(ch):
		val = decimalDigitValue(rune(ch))
	}
	if val >= 0 && val < radix {
		return val
	}
	return -1
}

func decimalDigitValue(r rune) int {
	for _, rg := range unicode.Nd.R16 {
		if r >= rune(rg.Lo) && r <= rune(rg.Hi) && rg.Stride == 1 {
			return int(r-rune(rg.Lo)) % 10
		}
	}
	for _, rg := range unicode.Nd.R32 {
		if r >= rune(rg.Lo) && r <= rune(rg.Hi) && rg.Stride == 1 {
			return int(r-rune(rg.Lo)) % 10
		}
	}
	return -1
}

// CharacterIsSupplementaryCodePoint implements Character.isSupplementaryCodePoint(int).
func CharacterIsSupplementaryCodePoint(codePoint int) bool {
	return codePoint >= 0x10000 && codePoint <= 0x10FFFF
}

// CharacterIsIdentifierIgnorable implements Character.isIdentifierIgnorable(char).
func CharacterIsIdentifierIgnorable(ch uint16) bool {
	if ch <= 0x08 || (ch >= 0x0E && ch <= 0x1B) || (ch >= 0x7F && ch <= 0x9F) {
		return true
	}
	return unicode.Is(unicode.Cf, rune(ch))
}

// CharacterIsJavaIdentifierStart implements Character.isJavaIdentifierStart(char).
func CharacterIsJavaIdentifierStart(ch uint16) bool {
	if ch < 0x80 {
		return (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z') || ch == '$' || ch == '_'
	}
	if ch >= 0xD800 && ch <= 0xDFFF {
		return false
	}
	return unicode.In(rune(ch), unicode.Lu, unicode.Ll, unicode.Lt, unicode.Lm, unicode.Lo, unicode.Nl, unicode.Sc, unicode.Pc)
}

// CharacterIsJavaIdentifierPart implements Character.isJavaIdentifierPart(char).
func CharacterIsJavaIdentifierPart(ch uint16) bool {
	if CharacterIsJavaIdentifierStart(ch) || CharacterIsDigit(ch) {
		return true
	}
	if ch >= 0xD800 && ch <= 0xDFFF {
		return false
	}
	if unicode.In(rune(ch), unicode.Mc, unicode.Mn) {
		return true
	}
	return CharacterIsIdentifierIgnorable(ch)
}

func checkRadix(radix int) {
	if radix < CharacterMinRadix {
		panic(NewNumberFormatException("radix " + strconv.Itoa(radix) + " less than Character.MIN_RADIX"))
	}
	if radix > CharacterMaxRadix {
		panic(NewNumberFormatException("radix " + strconv.Itoa(radix) + " greater than Character.MAX_RADIX"))
	}
}

func BadNumberInput(s string, radix int) *NumberFormatException {
	msg := "For input string: \"" + s + "\""
	if radix != 10 {
		msg += " under radix " + strconv.Itoa(radix)
	}
	return NewNumberFormatException(msg)
}

func parseRadixInt(s string, radix int, maxPos uint64) (int64, bool) {
	i := 0
	neg := false
	if len(s) > 0 && (s[0] == '-' || s[0] == '+') {
		neg = s[0] == '-'
		i = 1
	}
	if i == len(s) {
		return 0, false
	}
	limit := maxPos
	if neg {
		limit++
	}
	base := uint64(radix)
	var mag uint64
	for i < len(s) {
		var digit int
		if c := s[i]; c < utf8.RuneSelf {
			digit = CharacterDigit(uint16(c), radix)
			i++
		} else {
			r, sz := utf8.DecodeRuneInString(s[i:])
			if sz == 1 || r > 0xFFFF {
				return 0, false
			}
			digit = CharacterDigit(uint16(r), radix)
			i += sz
		}
		if digit < 0 || mag > (limit-uint64(digit))/base {
			return 0, false
		}
		mag = mag*base + uint64(digit)
	}
	if neg {
		return -int64(mag), true
	}
	return int64(mag), true
}

// IntegerParseInt implements Integer.parseInt(String).
func IntegerParseInt(s string) int32 {
	return IntegerParseIntRadix(s, 10)
}

// IntegerParseIntRadix implements Integer.parseInt(String, int).
func IntegerParseIntRadix(s string, radix int) int32 {
	checkRadix(radix)
	v, ok := parseRadixInt(s, radix, math.MaxInt32)
	if !ok {
		panic(BadNumberInput(s, radix))
	}
	return int32(v)
}

// LongParseLong implements Long.parseLong(String).
func LongParseLong(s string) int64 {
	return LongParseLongRadix(s, 10)
}

// LongParseLongRadix implements Long.parseLong(String, int).
func LongParseLongRadix(s string, radix int) int64 {
	checkRadix(radix)
	v, ok := parseRadixInt(s, radix, math.MaxInt64)
	if !ok {
		panic(BadNumberInput(s, radix))
	}
	return v
}

// DoubleParseDouble implements Double.parseDouble(String).
func DoubleParseDouble(s string) float64 {
	trimmed := strings.TrimFunc(s, func(r rune) bool { return r <= 0x20 })
	if trimmed == "" || strings.ContainsRune(trimmed, '_') {
		panic(NewNumberFormatException("For input string: \"" + s + "\""))
	}
	// Strip optional Java trailing float type suffix ('f', 'F', 'd', 'D') if not hex/NaN/Infinity
	clean := trimmed
	switch clean {
	case "NaN", "+NaN", "-NaN":
		if clean != "NaN" {
			panic(NewNumberFormatException("For input string: \"" + s + "\""))
		}
		return math.NaN()
	case "Infinity", "+Infinity":
		return math.Inf(1)
	case "-Infinity":
		return math.Inf(-1)
	}
	if last := clean[len(clean)-1]; (last == 'd' || last == 'D' || last == 'f' || last == 'F') && !strings.ContainsAny(clean[:len(clean)-1], "xX") {
		clean = clean[:len(clean)-1]
	}
	v, err := strconv.ParseFloat(clean, 64)
	if err != nil && err.(*strconv.NumError).Err != strconv.ErrRange {
		panic(NewNumberFormatException("For input string: \"" + s + "\""))
	}
	return v
}

// DoubleCompare implements Double.compare(double, double).
func DoubleCompare(a, b float64) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	na, nb := math.IsNaN(a), math.IsNaN(b)
	if na && nb {
		return 0
	}
	if na {
		return 1
	}
	if nb {
		return -1
	}
	sa, sb := math.Signbit(a), math.Signbit(b)
	if sa == sb {
		return 0
	}
	if sa {
		return -1
	}
	return 1
}

// DoubleIsNaN implements Double.isNaN(double).
func DoubleIsNaN(d float64) bool { return math.IsNaN(d) }

// DoubleIsInfinite implements Double.isInfinite(double).
func DoubleIsInfinite(d float64) bool { return math.IsInf(d, 0) }

// DoubleIsFinite implements Double.isFinite(double).
func DoubleIsFinite(d float64) bool { return !math.IsNaN(d) && !math.IsInf(d, 0) }

// CharacterToChars implements Character.toChars(int codePoint).
func CharacterToChars(codePoint int) []uint16 {
	if codePoint < 0 || codePoint > 0x10FFFF {
		panic(NewIllegalArgumentException("Not a valid Unicode code point: 0x" + strings.ToUpper(strconv.FormatInt(int64(codePoint), 16))))
	}
	if codePoint < 0x10000 {
		return []uint16{uint16(codePoint)}
	}
	u := codePoint - 0x10000
	return []uint16{uint16(0xD800 + (u >> 10)), uint16(0xDC00 + (u & 0x3FF))}
}

// CharacterGetType implements Character.getType(int codePoint) per Unicode General Category constants
// documented in java.lang.Character Javadoc Constant Field Values.
func CharacterGetType(codePoint int) int32 {
	if codePoint < 0 || codePoint > 0x10FFFF {
		return 0 // UNASSIGNED (Cn)
	}
	r := rune(codePoint)
	if r < 0x80 {
		switch {
		case r <= 0x1F || r == 0x7F:
			return 15 // CONTROL (Cc)
		case r == ' ':
			return 12 // SPACE_SEPARATOR (Zs)
		case r >= '0' && r <= '9':
			return 9 // DECIMAL_DIGIT_NUMBER (Nd)
		case r >= 'A' && r <= 'Z':
			return 1 // UPPERCASE_LETTER (Lu)
		case r >= 'a' && r <= 'z':
			return 2 // LOWERCASE_LETTER (Ll)
		}
	}
	if codePoint >= 0xD800 && codePoint <= 0xDFFF {
		return 19 // SURROGATE (Cs)
	}
	switch {
	case unicode.Is(unicode.Zs, r):
		return 12 // SPACE_SEPARATOR
	case unicode.Is(unicode.Cf, r):
		return 16 // FORMAT
	case unicode.Is(unicode.Lu, r):
		return 1 // UPPERCASE_LETTER
	case unicode.Is(unicode.Ll, r):
		return 2 // LOWERCASE_LETTER
	case unicode.Is(unicode.Lt, r):
		return 3 // TITLECASE_LETTER
	case unicode.Is(unicode.Lm, r):
		return 4 // MODIFIER_LETTER
	case unicode.Is(unicode.Lo, r):
		return 5 // OTHER_LETTER
	case unicode.Is(unicode.Mn, r):
		return 6 // NON_SPACING_MARK
	case unicode.Is(unicode.Me, r):
		return 7 // ENCLOSING_MARK
	case unicode.Is(unicode.Mc, r):
		return 8 // COMBINING_SPACING_MARK
	case unicode.Is(unicode.Nd, r):
		return 9 // DECIMAL_DIGIT_NUMBER
	case unicode.Is(unicode.Nl, r):
		return 10 // LETTER_NUMBER
	case unicode.Is(unicode.No, r):
		return 11 // OTHER_NUMBER
	case unicode.Is(unicode.Zl, r):
		return 13 // LINE_SEPARATOR
	case unicode.Is(unicode.Zp, r):
		return 14 // PARAGRAPH_SEPARATOR
	case unicode.Is(unicode.Cc, r):
		return 15 // CONTROL
	case unicode.Is(unicode.Co, r):
		return 18 // PRIVATE_USE
	case unicode.Is(unicode.Pd, r):
		return 20 // DASH_PUNCTUATION
	case unicode.Is(unicode.Ps, r):
		return 21 // START_PUNCTUATION
	case unicode.Is(unicode.Pe, r):
		return 22 // END_PUNCTUATION
	case unicode.Is(unicode.Pc, r):
		return 23 // CONNECTOR_PUNCTUATION
	case unicode.Is(unicode.Pi, r):
		return 29 // INITIAL_QUOTE_PUNCTUATION
	case unicode.Is(unicode.Pf, r):
		return 30 // FINAL_QUOTE_PUNCTUATION
	case unicode.Is(unicode.Po, r):
		return 24 // OTHER_PUNCTUATION
	case unicode.Is(unicode.Sm, r):
		return 25 // MATH_SYMBOL
	case unicode.Is(unicode.Sc, r):
		return 26 // CURRENCY_SYMBOL
	case unicode.Is(unicode.Sk, r):
		return 27 // MODIFIER_SYMBOL
	case unicode.Is(unicode.So, r):
		return 28 // OTHER_SYMBOL
	default:
		return 0 // UNASSIGNED
	}
}
