package jutil

// Go-only: language-adapter — String.toLowerCase(Locale.ROOT) and
// Character.toLowerCase(int), written from the Unicode Standard (Unicode
// 15.0, the version of both JDK 21 and Go 1.24's unicode package) and
// checked black-box against the JDK by string_case_test.go (vectors from
// tools/oracle/java/JutilStringCaseVectors.java).
//
// A string is lowercased code point by code point (a surrogate pair is one
// code point, a lone surrogate stays as it is) with the full lowercase
// mappings that apply in every language:
//   - the simple mapping of UnicodeData.txt, as for Character.toLowerCase;
//   - the unconditional entries of SpecialCasing.txt, of which only U+0130
//     LATIN CAPITAL LETTER I WITH DOT ABOVE differs from its simple mapping:
//     it becomes "i\u0307";
//   - the one conditional entry that is not language-specific: U+03A3
//     GREEK CAPITAL LETTER SIGMA becomes U+03C2 (final sigma) in the
//     Final_Sigma context and U+03C3 elsewhere.
// Apart from final sigma, this matches the JDK on every code point and on
// all the random strings of the vectors.
//
// Final sigma follows the Unicode Standard (§3.13, Table 3-17): the sigma
// comes after a cased character and then only case-ignorable ones, and is
// not followed by case-ignorable characters and then a cased one (Cased and
// Case_Ignorable as D135 and D136 define them, from Go's unicode tables and
// the Word_Break values of UAX #29; they match DerivedCoreProperties.txt on
// every code point). This is a deliberate choice: the JDK decides final
// sigma with its own word-boundary logic, which has no public
// specification; this file does not try to replicate it. The two disagree
// in these contexts (measured black-box on JDK 21, see string_case_test.go):
//   - Ordinary text agrees: a sigma next to BMP letters with case (Lu, Ll,
//     Lt), nonspacing marks, spaces, and punctuation such as , ; ! ? ( ),
//     or one apostrophe or full stop. In random Greek and Latin prose no
//     sigma differs.
//   - The JDK keeps in the word, where Unicode ends it: digits, '_', '-' and
//     the other dashes, '"', spacing marks (Mc), and the BMP letters without
//     case except kana and the CJK ideographs of Unicode 1.1 (Hebrew,
//     Arabic, Hangul, ...). "ΟΔΟΣ-ΑΘΗΝΩΝ" is "οδοσ-αθηνων" in Java and
//     "οδος-αθηνων" here; "Α1Σ" is "α1ς" in Java and "α1σ" here.
//   - The JDK stops, where Unicode skips case-ignorable characters, at
//     modifier symbols (^ ` ´ ...), the iteration and length marks
//     々 ゝ ゞ ー ヽ ヾ, the word-internal punctuation : · ’ ‘ and the
//     like, runs of several apostrophes, full stops or format characters
//     (Unicode skips "'." and two soft hyphens; one alone agrees), and every
//     character outside the BMP ("Α\U00010400Σ" ends in σ in Java and in ς
//     here).
//   - Cased characters that are not letters with case (ª, º, some modifier
//     letters such as ᵢ, circled letters such as Ⓐ, U+0345) are cased here
//     and not for the JDK.
// Over every code point c, the contexts "AcΣ", "cΣ", "AΣca" and "AΣc" give
// a different sigma for 28720, 328, 27228 and 327 values of c. On random
// strings drawn from such characters about 13% of the sigmas differ, and on
// random Greek words joined by - _ : ’ " digits and the like, about 19%.

import "unicode"

// StringToLowerCaseRoot is String.toLowerCase(Locale.ROOT).
func StringToLowerCaseRoot(s string) string {
	if isASCII(s) {
		for i := 0; i < len(s); i++ {
			if c := s[i]; c >= 'A' && c <= 'Z' {
				b := []byte(s)
				for j := i; j < len(b); j++ {
					if c := b[j]; c >= 'A' && c <= 'Z' {
						b[j] = c + ('a' - 'A')
					}
				}
				return string(b)
			}
		}
		return s
	}
	u := StringToUTF16(s)
	var out []uint16 // stays nil until the first code point that changes
	for i := 0; i < len(u); {
		cp, n := utf16CodePointAt(u, i)
		var lower rune
		switch cp {
		case 0x0130:
			if out == nil {
				out = append(make([]uint16, 0, len(u)+8), u[:i]...)
			}
			out = append(out, 'i', 0x0307)
			i += n
			continue
		case 0x03A3:
			if isFinalSigma(u, i) {
				lower = 0x03C2
			} else {
				lower = 0x03C3
			}
		default:
			lower = rune(CharacterToLowerCaseCodePoint(int(cp)))
		}
		if lower != cp && out == nil {
			out = append(make([]uint16, 0, len(u)+8), u[:i]...)
		}
		if out != nil {
			out = appendUTF16CodePoint(out, lower)
		}
		i += n
	}
	if out == nil {
		return s
	}
	return StringFromUTF16(out)
}

// CharacterToLowerCaseCodePoint is Character.toLowerCase(int): the simple
// lowercase mapping of UnicodeData.txt. Ints that are not code points, and
// surrogates, are returned unchanged (observed behaviour, black-box: -1,
// Integer.MIN_VALUE and 0x110041 map to themselves).
func CharacterToLowerCaseCodePoint(codePoint int) int {
	if codePoint < 0x80 {
		if codePoint >= 'A' && codePoint <= 'Z' {
			return codePoint + ('a' - 'A')
		}
		return codePoint
	}
	if codePoint > unicode.MaxRune || codePoint >= 0xD800 && codePoint <= 0xDFFF {
		return codePoint
	}
	return int(unicode.ToLower(rune(codePoint)))
}

// isFinalSigma reports whether the sigma at u[i] is in the Final_Sigma
// context of the Unicode Standard (§3.13, Table 3-17).
func isFinalSigma(u []uint16, i int) bool {
	// Before: a cased character, then zero or more case-ignorable ones. A
	// character can be both (U+02B0 MODIFIER LETTER SMALL H is Lm and
	// Other_Lowercase); it then counts as the cased one.
	for j := i; ; {
		if j == 0 {
			return false
		}
		cp, n := utf16CodePointBefore(u, j)
		if isCased(cp) {
			break
		}
		if !isCaseIgnorable(cp) {
			return false
		}
		j -= n
	}
	// After: not zero or more case-ignorable characters and then a cased
	// one.
	for j := i + 1; j < len(u); {
		cp, n := utf16CodePointAt(u, j)
		if isCased(cp) {
			return false
		}
		if !isCaseIgnorable(cp) {
			break
		}
		j += n
	}
	return true
}

// isCased is the Cased property (Unicode §3.13, D135): Lowercase or
// Uppercase (the general categories Ll and Lu plus Other_Lowercase and
// Other_Uppercase), or the general category Lt.
func isCased(cp rune) bool {
	if cp < 0x80 {
		return cp >= 'a' && cp <= 'z' || cp >= 'A' && cp <= 'Z'
	}
	return unicode.In(cp, unicode.Ll, unicode.Lu, unicode.Lt, unicode.Other_Lowercase, unicode.Other_Uppercase)
}

// isCaseIgnorable is the Case_Ignorable property (Unicode §3.13, D136): the
// general categories Mn, Me, Cf, Lm and Sk, plus the characters whose
// Word_Break property (UAX #29) is MidLetter, MidNumLet or Single_Quote.
func isCaseIgnorable(cp rune) bool {
	switch cp {
	case
		// Word_Break=MidLetter
		0x003A, 0x00B7, 0x0387, 0x055F, 0x05F4, 0x2027, 0xFE13, 0xFE55, 0xFF1A,
		// Word_Break=MidNumLet
		0x002E, 0x2018, 0x2019, 0x2024, 0xFE52, 0xFF07, 0xFF0E,
		// Word_Break=Single_Quote
		0x0027:
		return true
	}
	if cp < 0x80 {
		return cp == '^' || cp == '`' // the ASCII Sk characters
	}
	return unicode.In(cp, unicode.Mn, unicode.Me, unicode.Cf, unicode.Lm, unicode.Sk)
}

// utf16CodePointAt decodes the code point starting at u[i]: a surrogate pair,
// or a single unit (a lone surrogate stands for itself).
func utf16CodePointAt(u []uint16, i int) (rune, int) {
	c := u[i]
	if c >= 0xD800 && c <= 0xDBFF && i+1 < len(u) {
		if lo := u[i+1]; lo >= 0xDC00 && lo <= 0xDFFF {
			return (rune(c)-0xD800)<<10 + (rune(lo) - 0xDC00) + 0x10000, 2
		}
	}
	return rune(c), 1
}

// utf16CodePointBefore decodes the code point that ends just before u[j].
func utf16CodePointBefore(u []uint16, j int) (rune, int) {
	c := u[j-1]
	if c >= 0xDC00 && c <= 0xDFFF && j >= 2 {
		if hi := u[j-2]; hi >= 0xD800 && hi <= 0xDBFF {
			return (rune(hi)-0xD800)<<10 + (rune(c) - 0xDC00) + 0x10000, 2
		}
	}
	return rune(c), 1
}

func appendUTF16CodePoint(out []uint16, cp rune) []uint16 {
	if cp < 0x10000 {
		return append(out, uint16(cp))
	}
	cp -= 0x10000
	return append(out, uint16(0xD800+(cp>>10)), uint16(0xDC00+(cp&0x3FF)))
}
