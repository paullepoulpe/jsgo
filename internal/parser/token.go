// Port of com/google/javascript/jscomp/parsing/parser/Token.java @ closure-compiler v20260415

package parser

// Token is a Token in a javascript file.
// Immutable.
// A plain old data structure. Should contain data members and simple accessors only.
//
// Java: class Token (with subclasses IdentifierToken and LiteralToken).
type Token interface {
	GetTokenBase() *TokenBase
	GetStart() *SourcePosition
	String() string
	AsIdentifier() *IdentifierToken
	AsLiteral() LiteralToken
	AsTemplateLiteral() *TemplateLiteralToken
}

// TokenBase holds Token's fields and base methods; plain tokens are
// *TokenBase values.
type TokenBase struct {
	self Token

	Location *SourceRange
	Type     TokenType
}

// NewToken constructs a plain Token.
//
// Java: Token(TokenType, SourceRange)
func NewToken(type_ TokenType, location *SourceRange) *TokenBase {
	t := &TokenBase{}
	t.InitTokenBase(t, type_, location)
	return t
}

// InitTokenBase is the Go form of the Token constructor for subclasses.
//
// Java: Token(TokenType, SourceRange)
func (t *TokenBase) InitTokenBase(self Token, type_ TokenType, location *SourceRange) {
	t.self = self
	t.Type = type_
	t.Location = location
}

// GetTokenBase returns the base struct (field access through the interface).
func (t *TokenBase) GetTokenBase() *TokenBase {
	return t
}

// GetStart returns the start position.
//
// Java: Token.getStart()
func (t *TokenBase) GetStart() *SourcePosition {
	return t.Location.Start
}

// String returns the token type's string.
//
// Java: Token.toString()
func (t *TokenBase) String() string {
	return t.Type.String()
}

// AsIdentifier casts to IdentifierToken.
//
// Java: Token.asIdentifier()
func (t *TokenBase) AsIdentifier() *IdentifierToken {
	return t.self.(*IdentifierToken)
}

// AsLiteral casts to LiteralToken.
//
// Java: Token.asLiteral()
func (t *TokenBase) AsLiteral() LiteralToken {
	return t.self.(LiteralToken)
}

// AsTemplateLiteral casts to TemplateLiteralToken.
//
// Java: Token.asTemplateLiteral()
func (t *TokenBase) AsTemplateLiteral() *TemplateLiteralToken {
	return t.self.(*TemplateLiteralToken)
}
