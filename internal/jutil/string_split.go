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
