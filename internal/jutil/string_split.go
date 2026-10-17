package jutil

import "strings"

// Go-only: language-adapter — String.split(String regex) (limit 0) for a
// regex that matches only the literal, non-empty string sep (e.g. "/" or
// "/node_modules/"): the input is cut at each non-overlapping occurrence,
// the input itself is returned when there is no occurrence, and trailing
// empty strings are removed (so the result may be empty).

// StringSplitLiteral is s.split(regex) where regex matches exactly sep.
func StringSplitLiteral(s, sep string) []string {
	if sep == "" {
		panic(NewIllegalArgumentExceptionMsg("StringSplitLiteral: empty separator"))
	}
	if !strings.Contains(s, sep) {
		return []string{s}
	}
	parts := strings.Split(s, sep)
	n := len(parts)
	for n > 0 && parts[n-1] == "" {
		n--
	}
	return parts[:n]
}

// StringSplitLiteralLimit is s.split(regex, limit) where regex matches
// exactly the non-empty literal sep. With limit > 0 there are at most limit
// parts, the last holding the rest of the input, and trailing empty strings
// are kept; with limit == 0 it is StringSplitLiteral; with limit < 0 every
// part is returned, trailing empty strings included. Without an occurrence
// of sep the result is s itself.
func StringSplitLiteralLimit(s, sep string, limit int) []string {
	if limit == 0 {
		return StringSplitLiteral(s, sep)
	}
	if sep == "" {
		panic(NewIllegalArgumentExceptionMsg("StringSplitLiteralLimit: empty separator"))
	}
	if limit < 0 {
		return strings.Split(s, sep)
	}
	return strings.SplitN(s, sep, limit)
}
