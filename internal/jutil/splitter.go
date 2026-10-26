package jutil

// Go-only: language-adapter — Guava Splitter and CharMatcher.whitespace().

// charMatcherWhitespaceTable is CharMatcher.Whitespace.TABLE: the characters
// CharMatcher.whitespace() matches (Unicode White_Space, per Guava; this is
// not the set that String.trim() strips).
const charMatcherWhitespaceTable = " 　\r\u0085   　" +
	" \u000B　   　 " +
	"\u0009     \u000C " +
	"　 　　 \n 　"

// CharMatcher_Whitespace_Matches is CharMatcher.whitespace().matches(c).
func CharMatcher_Whitespace_Matches(c uint16) bool {
	// Guava: TABLE.charAt((MULTIPLIER * c) >>> SHIFT) == c, where SHIFT is
	// Integer.numberOfLeadingZeros(TABLE.length() - 1) = 27.
	const multiplier = 1682554634
	idx := uint32(int32(multiplier)*int32(c)) >> 27
	return charMatcherWhitespaceUnits[idx] == c
}

var charMatcherWhitespaceUnits = StringToUTF16(charMatcherWhitespaceTable)

// CharMatcher_Whitespace_TrimFrom is CharMatcher.whitespace().trimFrom(s).
func CharMatcher_Whitespace_TrimFrom(s string) string {
	u := StringToUTF16(s)
	first, last := 0, len(u)-1
	for first < len(u) && CharMatcher_Whitespace_Matches(u[first]) {
		first++
	}
	for last > first && CharMatcher_Whitespace_Matches(u[last]) {
		last--
	}
	return StringFromUTF16(u[first : last+1])
}

// Splitter_OnCharTrimResults is
// Splitter.on(separator).trimResults().split(s), as a list: every piece
// between separators (including empty ones), trimmed of
// CharMatcher.whitespace().
func Splitter_OnCharTrimResults(s string, separator uint16) []string {
	u := StringToUTF16(s)
	var out []string
	start := 0
	for i := 0; i <= len(u); i++ {
		if i == len(u) || u[i] == separator {
			out = append(out, CharMatcher_Whitespace_TrimFrom(StringFromUTF16(u[start:i])))
			start = i + 1
		}
	}
	return out
}
