// Port of com/google/javascript/rhino/JSIdentifier.java @ closure-compiler v20260415

package rhino

import "github.com/paullepoulpe/jsgo/internal/jutil"

// JSIdentifier is a utility class to hold isJSIdentifier.
//
// IMPORTANT: As of 2018-03-09 it is still not possible to use Java 8
// features in this file due to limitations on some internal Google projects
// that depend on it.
type JSIdentifier struct{}

// Java: JSIdentifier.isJSIdentifier(String)
func JSIdentifier_IsJSIdentifier(s string) bool {
	chars := jutil.StringToUTF16(s) // s.charAt(i) is chars[i]
	length := len(chars)

	if length == 0 ||
		jutil.CharacterIsIdentifierIgnorable(chars[0]) ||
		!jutil.CharacterIsJavaIdentifierStart(chars[0]) {
		return false
	}

	for i := 1; i < length; i++ {
		if jutil.CharacterIsIdentifierIgnorable(chars[i]) ||
			!jutil.CharacterIsJavaIdentifierPart(chars[i]) {
			return false
		}
	}
	return true
}

// Java: JSIdentifier()
func NewJSIdentifier() *JSIdentifier { return &JSIdentifier{} }
