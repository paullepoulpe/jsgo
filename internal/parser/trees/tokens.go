// Port of com/google/javascript/jscomp/parsing/parser/{Token,IdentifierToken,LiteralToken}.java
// and com/google/javascript/jscomp/parsing/parser/util/{SourceRange,SourcePosition}.java
// @ closure-compiler v20260415
//
// These supporting types live in the trees package so that the tree structs
// remain self-contained (no cross-package import cycle at this layer).

package trees

import "fmt"

// SourceFile is a minimal stand-in for
// com.google.javascript.jscomp.parsing.parser.SourceFile.
type SourceFile struct {
	Name     string
	Contents string
}

// SourcePosition mirrors com.google.javascript.jscomp.parsing.parser.util.SourcePosition.
type SourcePosition struct {
	// Source is the source file (nil for synthesised positions).
	Source *SourceFile
	Offset int
	Line   int
	Column int
}

func (p SourcePosition) String() string {
	name := ""
	if p.Source != nil {
		name = p.Source.Name
	}
	return fmt.Sprintf("%s(%d, %d)", name, p.Line+1, p.Column+1)
}

// SourceRange mirrors com.google.javascript.jscomp.parsing.parser.util.SourceRange.
// Start is inclusive, End is exclusive.
type SourceRange struct {
	Start SourcePosition
	End   SourcePosition
}

func (r SourceRange) String() string {
	return fmt.Sprintf("<%s - %s>", r.Start, r.End)
}

// TokenType mirrors com.google.javascript.jscomp.parsing.parser.TokenType.
// Only the entries relevant to the trees layer are listed here; additional
// values will be added as more of the parser is ported.
type TokenType int

const (
	TokenTypeError TokenType = iota
	TokenTypeIdentifier
	TokenTypeNumber
	TokenTypeString
	TokenTypeRegularExpression
	TokenTypeNoSubstitutionTemplate
	TokenTypeTemplateHead
	TokenTypeTemplateMiddle
	TokenTypeTemplateTail
	TokenTypeNull
	TokenTypeTrue
	TokenTypeFalse
	TokenTypeVar
	TokenTypeLet
	TokenTypeConst
)

// Token mirrors com.google.javascript.jscomp.parsing.parser.Token.
type Token struct {
	Type     TokenType
	Location SourceRange
}

func NewToken(tokenType TokenType, location SourceRange) *Token {
	return &Token{Type: tokenType, Location: location}
}

func (t *Token) GetStart() SourcePosition {
	return t.Location.Start
}

func (t *Token) String() string {
	return fmt.Sprintf("%v", t.Type)
}

// IdentifierToken mirrors com.google.javascript.jscomp.parsing.parser.IdentifierToken.
type IdentifierToken struct {
	Token
	value            string
	privateIdentifier bool
}

func NewIdentifierToken(location SourceRange, value string) *IdentifierToken {
	t := &IdentifierToken{
		Token: Token{Type: TokenTypeIdentifier, Location: location},
		value: value,
	}
	t.privateIdentifier = len(value) > 0 && value[0] == '#'
	return t
}

func (t *IdentifierToken) String() string { return t.value }

func (t *IdentifierToken) ValueEquals(s string) bool { return t.value == s }

func (t *IdentifierToken) GetValue() string { return t.value }

func (t *IdentifierToken) GetMaybePrivateValue() string { return t.value }

func (t *IdentifierToken) IsPrivateIdentifier() bool { return t.privateIdentifier }

// LiteralToken mirrors com.google.javascript.jscomp.parsing.parser.LiteralToken.
type LiteralToken struct {
	Token
	Value string
}

func NewLiteralToken(tokenType TokenType, value string, location SourceRange) *LiteralToken {
	return &LiteralToken{
		Token: Token{Type: tokenType, Location: location},
		Value: value,
	}
}

func (t *LiteralToken) String() string { return t.Value }
