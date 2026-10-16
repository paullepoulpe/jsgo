package jutil

// StringReplaceFirst is String.replaceFirst(String regex, String
// replacement).
func StringReplaceFirst(s string, regex string, replacement string) string {
	return Pattern_Compile(regex).Matcher(s).ReplaceFirst(replacement)
}

// StringReplaceAll is String.replaceAll(String regex, String replacement).
func StringReplaceAll(s string, regex string, replacement string) string {
	return Pattern_Compile(regex).Matcher(s).ReplaceAll(replacement)
}
