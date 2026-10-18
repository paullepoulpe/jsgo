package jutil

// Go-only: language-adapter — String.toUpperCase(Locale.ROOT)
// (String.toLowerCase(Locale.ROOT) is StringToLowerCaseRoot in
// string_to_lower_case.go).
//
// For the root locale it is context-free per code point. The simple mapping
// is Character.toUpperCase(int), which matches Go's unicode.ToUpper on every
// code point (Unicode 15.0 in JDK 21 and Go 1.24); the full mappings of
// SpecialCasing.txt are in jstring_case_table.go. Both are checked against
// the JDK by TestStringCaseRootMatchesJava.

import (
	unicode "github.com/paullepoulpe/jsgo/internal/unicode15"
	"strings"
)

// StringToUpperCaseRoot is String.toUpperCase(Locale.ROOT).
func StringToUpperCaseRoot(s string) string {
	if isASCII(s) {
		return StringToUpperCaseASCII(s)
	}
	var sb strings.Builder
	sb.Grow(len(s))
	for i := 0; i < len(s); {
		r, n := decodeWTF8(s[i:])
		i += n
		if special, ok := stringUpperCaseRootSpecial[r]; ok {
			sb.WriteString(special)
			continue
		}
		appendWTF8(&sb, unicode.ToUpper(r))
	}
	return sb.String()
}
