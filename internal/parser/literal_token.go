// Port of com/google/javascript/jscomp/parsing/parser/LiteralToken.java @ closure-compiler v20260415

package parser

// LiteralToken is a token representing a javascript literal. Includes
// string, regexp, and number literals. Boolean and null literals are
// represented as regular keyword tokens.
//
// The value just includes the raw lexeme. For string literals it includes
// the beginning and ending delimiters.
//
// TODO: Regexp literals should have their own token type.
// TODO: A way to get the processed value, rather than the raw value.
//
// Java: class LiteralToken (with subclasses StringLiteralToken and
// TemplateLiteralToken).
type LiteralToken interface {
	Token
	GetLiteralTokenBase() *LiteralTokenBase
}

// LiteralTokenBase holds LiteralToken's fields; plain literal tokens are
// *LiteralTokenBase values.
type LiteralTokenBase struct {
	TokenBase

	Value string
}

// NewLiteralToken constructs a plain LiteralToken.
//
// Java: LiteralToken(TokenType, String, SourceRange)
func NewLiteralToken(type_ TokenType, value string, location *SourceRange) *LiteralTokenBase {
	t := &LiteralTokenBase{}
	t.InitLiteralTokenBase(t, type_, value, location)
	return t
}

// InitLiteralTokenBase is the Go form of the LiteralToken constructor for
// subclasses.
//
// Java: LiteralToken(TokenType, String, SourceRange)
func (t *LiteralTokenBase) InitLiteralTokenBase(self LiteralToken, type_ TokenType, value string, location *SourceRange) {
	t.InitTokenBase(self, type_, location)
	t.Value = value
}

// GetLiteralTokenBase returns the base struct (field access through the
// interface).
func (t *LiteralTokenBase) GetLiteralTokenBase() *LiteralTokenBase {
	return t
}

// String returns the raw lexeme.
//
// Java: LiteralToken.toString()
func (t *LiteralTokenBase) String() string {
	return t.Value
}
