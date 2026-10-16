package jutil

// Go-only: language-adapter — Guava's com.google.common.base.CaseFormat
// conversions not covered by CaseFormat_UpperUnderscoreToLowerCamel.

// CaseFormat_LowerCamelToUpperUnderscore is Guava's
// CaseFormat.LOWER_CAMEL.to(CaseFormat.UPPER_UNDERSCORE, s): words start at
// each ASCII uppercase letter and are upper-cased (Ascii.toUpperCase) and
// joined with '_'. This is Guava's CaseFormat.convert, step for step.
func CaseFormat_LowerCamelToUpperUnderscore(s string) string {
	units := StringToUTF16(s)
	isBoundary := func(c uint16) bool { return c >= 'A' && c <= 'Z' }
	indexIn := func(from int) int {
		for k := from; k < len(units); k++ {
			if isBoundary(units[k]) {
				return k
			}
		}
		return -1
	}
	upper := func(word []uint16) string {
		out := make([]uint16, len(word))
		for k, c := range word {
			if c >= 'a' && c <= 'z' {
				c -= 'a' - 'A'
			}
			out[k] = c
		}
		return StringFromUTF16(out)
	}
	var out *StringBuilder
	i := 0
	j := -1
	for {
		j++
		j = indexIn(j)
		if j == -1 {
			break
		}
		if i == 0 {
			// include some extra space for separators
			out = NewStringBuilder()
			out.Append(upper(units[i:j]))
		} else {
			out.Append(upper(units[i:j]))
		}
		out.Append("_")
		i = j // + wordSeparator.length() == 0 for LOWER_CAMEL
	}
	if i == 0 {
		return upper(units)
	}
	return out.Append(upper(units[i:])).String()
}
