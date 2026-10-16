// Port of com/google/javascript/jscomp/parsing/parser/StringLiteralToken.java @ closure-compiler v20260415

package parser

// StringLiteralToken is a single or double quoted JavaScript string literal.
type StringLiteralToken struct {
	LiteralTokenBase

	hasUnescapedUnicodeLineOrParagraphSeparator bool
}

// NewStringLiteralToken constructs a StringLiteralToken.
//
// Java: StringLiteralToken(String, SourceRange, boolean)
func NewStringLiteralToken(value string, location *SourceRange, hasUnescapedUnicodeLineOrParagraphSeparator bool) *StringLiteralToken {
	t := &StringLiteralToken{}
	t.InitLiteralTokenBase(t, TokenType_STRING, value, location)
	t.hasUnescapedUnicodeLineOrParagraphSeparator = hasUnescapedUnicodeLineOrParagraphSeparator
	return t
}

// HasUnescapedUnicodeLineOrParagraphSeparator reports whether the literal
// contains an unescaped U+2028 or U+2029.
//
// Java: StringLiteralToken.hasUnescapedUnicodeLineOrParagraphSeparator()
func (t *StringLiteralToken) HasUnescapedUnicodeLineOrParagraphSeparator() bool {
	return t.hasUnescapedUnicodeLineOrParagraphSeparator
}
