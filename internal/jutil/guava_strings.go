package jutil

// Go-only: language-adapter — Guava com.google.common.base.Strings; a Java
// String reference that may be null is a *string (nil is null).

// Strings_IsNullOrEmpty is Strings.isNullOrEmpty(String): true if the
// string is null or the empty string.
func Strings_IsNullOrEmpty(s *string) bool {
	return s == nil || *s == ""
}
