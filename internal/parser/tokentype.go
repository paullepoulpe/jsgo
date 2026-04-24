// Port of com/google/javascript/jscomp/parsing/parser/TokenType.java,
//         com/google/javascript/jscomp/parsing/parser/Token.java,
//         com/google/javascript/jscomp/parsing/parser/IdentifierToken.java,
//         com/google/javascript/jscomp/parsing/parser/LiteralToken.java,
//         com/google/javascript/jscomp/parsing/parser/TemplateLiteralToken.java,
//         com/google/javascript/jscomp/parsing/parser/Keywords.java
// @ closure-compiler v20260415.

package parser

// SourcePosition mirrors com.google.javascript.jscomp.parsing.parser.util.SourcePosition.
// offset is a rune-index into the source, line is 0-based (Java convention),
// column is 0-based.
type SourcePosition struct {
	Offset int
	Line   int
	Column int
}

// SourceRange mirrors com.google.javascript.jscomp.parsing.parser.util.SourceRange.
// Start is inclusive, End is exclusive.
type SourceRange struct {
	Start SourcePosition
	End   SourcePosition
}

// Comment mirrors com.google.javascript.jscomp.parsing.parser.trees.Comment.
type Comment struct {
	Type  CommentType
	Range SourceRange
	Value string
}

// CommentType mirrors Comment.Type.
type CommentType int

const (
	CommentLine      CommentType = iota // //
	CommentBlock                        // /* */
	CommentJSDoc                        // /** */
	CommentImportant                    // /*! */
	CommentShebang                      // #!
)

// TokenType mirrors com.google.javascript.jscomp.parsing.parser.TokenType.
// 7.5 Tokens
type TokenType int

const (
	// sentinel / error
	TokenTypeEndOfFile TokenType = iota
	TokenTypeError

	// 7.6 Identifier Names and Identifiers
	TokenTypeIdentifier

	// 7.6.1.1 Keywords
	TokenTypeBreak
	TokenTypeCase
	TokenTypeCatch
	TokenTypeContinue
	TokenTypeDebugger
	TokenTypeDefault
	TokenTypeDelete
	TokenTypeDo
	TokenTypeElse
	TokenTypeFinally
	TokenTypeFor
	TokenTypeFunction
	TokenTypeIf
	TokenTypeIn
	TokenTypeInstanceof
	TokenTypeNew
	TokenTypeReturn
	TokenTypeSwitch
	TokenTypeThis
	TokenTypeThrow
	TokenTypeTry
	TokenTypeTypeof
	TokenTypeVar
	TokenTypeVoid
	TokenTypeWhile
	TokenTypeWith

	// 7.6.1.2 Future reserved words
	TokenTypeClass
	TokenTypeConst
	TokenTypeEnum
	TokenTypeExport
	TokenTypeExtends
	TokenTypeImport
	TokenTypeSuper

	// Future reserved words in strict mode
	TokenTypeImplements
	TokenTypeInterface
	TokenTypeLet
	TokenTypePackage
	TokenTypePrivate
	TokenTypeProtected
	TokenTypePublic
	TokenTypeStatic
	TokenTypeYield

	// 7.7 Punctuators
	TokenTypeOpenCurly
	TokenTypeCloseCurly
	TokenTypeOpenParen
	TokenTypeCloseParen
	TokenTypeOpenSquare
	TokenTypeCloseSquare
	TokenTypePeriod
	TokenTypeSemiColon
	TokenTypeComma
	TokenTypeOpenAngle
	TokenTypeCloseAngle
	TokenTypeLessEqual
	TokenTypeGreaterEqual
	TokenTypeArrow
	TokenTypeEqualEqual
	TokenTypeNotEqual
	TokenTypeEqualEqualEqual
	TokenTypeNotEqualEqual
	TokenTypePlus
	TokenTypeMinus
	TokenTypeStar
	TokenTypeStarStar
	TokenTypePercent
	TokenTypePlusPlus
	TokenTypeMinusMinus
	TokenTypeLeftShift
	TokenTypeRightShift
	TokenTypeUnsignedRightShift
	TokenTypeAmpersand
	TokenTypeBar
	TokenTypeCaret
	TokenTypeBang
	TokenTypeTilde
	TokenTypeAnd
	TokenTypeOr
	TokenTypeQuestion
	TokenTypeQuestionQuestion
	TokenTypeQuestionDot
	TokenTypeColon
	TokenTypeEqual
	TokenTypePlusEqual
	TokenTypeMinusEqual
	TokenTypeStarEqual
	TokenTypeStarStarEqual
	TokenTypePercentEqual
	TokenTypeLeftShiftEqual
	TokenTypeRightShiftEqual
	TokenTypeUnsignedRightShiftEqual
	TokenTypeAmpersandEqual
	TokenTypeBarEqual
	TokenTypeCaretEqual
	TokenTypeSlash
	TokenTypeSlashEqual

	// Logical Assignment Punctuators
	TokenTypeAndEqual
	TokenTypeOrEqual
	TokenTypeQuestionQuestionEqual

	// 7.8 Literals
	TokenTypeNull
	TokenTypeTrue
	TokenTypeFalse
	TokenTypeNumber
	TokenTypeString
	TokenTypeBigInt
	TokenTypeRegularExpression

	// Harmony extensions
	TokenTypeEllipsis

	// 12.2.9 Template Literals
	TokenTypeTemplateHead
	TokenTypeTemplateMiddle
	TokenTypeTemplateTail
	TokenTypeNoSubstitutionTemplate

	// TypeScript
	TokenTypeType
	TokenTypeDeclare
	TokenTypeModule
	TokenTypeNamespace
)

// ToString mirrors TokenType.toString() — returns the value string for the token type.
func (t TokenType) ToString() string {
	return tokenTypeValue(t)
}

// Go-only: language-adapter — tokenTypeValue is a helper backing ToString(); Java stores the value
// directly as an enum field `value` rather than a function.
func tokenTypeValue(t TokenType) string {
	switch t {
	case TokenTypeEndOfFile:
		return "End of File"
	case TokenTypeError:
		return "error"
	case TokenTypeIdentifier:
		return "identifier"
	case TokenTypeOpenCurly:
		return "{"
	case TokenTypeCloseCurly:
		return "}"
	case TokenTypeOpenParen:
		return "("
	case TokenTypeCloseParen:
		return ")"
	case TokenTypeOpenSquare:
		return "["
	case TokenTypeCloseSquare:
		return "]"
	case TokenTypePeriod:
		return "."
	case TokenTypeSemiColon:
		return ";"
	case TokenTypeComma:
		return ","
	case TokenTypeOpenAngle:
		return "<"
	case TokenTypeCloseAngle:
		return ">"
	case TokenTypeLessEqual:
		return "<="
	case TokenTypeGreaterEqual:
		return ">="
	case TokenTypeArrow:
		return "=>"
	case TokenTypeEqualEqual:
		return "=="
	case TokenTypeNotEqual:
		return "!="
	case TokenTypeEqualEqualEqual:
		return "==="
	case TokenTypeNotEqualEqual:
		return "!=="
	case TokenTypePlus:
		return "+"
	case TokenTypeMinus:
		return "-"
	case TokenTypeStar:
		return "*"
	case TokenTypeStarStar:
		return "**"
	case TokenTypePercent:
		return "%"
	case TokenTypePlusPlus:
		return "++"
	case TokenTypeMinusMinus:
		return "--"
	case TokenTypeLeftShift:
		return "<<"
	case TokenTypeRightShift:
		return ">>"
	case TokenTypeUnsignedRightShift:
		return ">>>"
	case TokenTypeAmpersand:
		return "&"
	case TokenTypeBar:
		return "|"
	case TokenTypeCaret:
		return "^"
	case TokenTypeBang:
		return "!"
	case TokenTypeTilde:
		return "~"
	case TokenTypeAnd:
		return "&&"
	case TokenTypeOr:
		return "||"
	case TokenTypeQuestion:
		return "?"
	case TokenTypeQuestionQuestion:
		return "??"
	case TokenTypeQuestionDot:
		return "?."
	case TokenTypeColon:
		return ":"
	case TokenTypeEqual:
		return "="
	case TokenTypePlusEqual:
		return "+="
	case TokenTypeMinusEqual:
		return "-="
	case TokenTypeStarEqual:
		return "*="
	case TokenTypeStarStarEqual:
		return "**="
	case TokenTypePercentEqual:
		return "%="
	case TokenTypeLeftShiftEqual:
		return "<<="
	case TokenTypeRightShiftEqual:
		return ">>="
	case TokenTypeUnsignedRightShiftEqual:
		return ">>>="
	case TokenTypeAmpersandEqual:
		return "&="
	case TokenTypeBarEqual:
		return "|="
	case TokenTypeCaretEqual:
		return "^="
	case TokenTypeSlash:
		return "/"
	case TokenTypeSlashEqual:
		return "/="
	case TokenTypeAndEqual:
		return "&&="
	case TokenTypeOrEqual:
		return "||="
	case TokenTypeQuestionQuestionEqual:
		return "??="
	case TokenTypeNumber:
		return "number literal"
	case TokenTypeString:
		return "string literal"
	case TokenTypeBigInt:
		return "bigint literal"
	case TokenTypeRegularExpression:
		return "regular expression literal"
	case TokenTypeEllipsis:
		return "..."
	case TokenTypeTemplateHead:
		return "template head"
	case TokenTypeTemplateMiddle:
		return "template middle"
	case TokenTypeTemplateTail:
		return "template tail"
	case TokenTypeNoSubstitutionTemplate:
		return "no substitution template"
	default:
		// keyword tokens: return the keyword string via keywordForType
		if s := keywordForType(t); s != "" {
			return s
		}
		return ""
	}
}

// Token mirrors com.google.javascript.jscomp.parsing.parser.Token (and its
// subclasses IdentifierToken, LiteralToken, StringLiteralToken,
// TemplateLiteralToken). In Java these are separate classes connected by
// inheritance; in Go we carry the union of all fields in one struct and
// leave unused fields at their zero values. The Type field determines which
// fields are meaningful (same contract as Java's `.asLiteral()` casts).
type Token struct {
	Type     TokenType
	Location SourceRange

	// Value is set for LiteralToken (NUMBER, BIGINT, STRING, REGULAR_EXPRESSION,
	// TEMPLATE_*, NO_SUBSTITUTION_TEMPLATE) and IdentifierToken (IDENTIFIER).
	// For plain punctuator/keyword tokens it is "".
	Value string

	// HasUnescapedUnicodeLineOrParagraphSeparator is set for STRING tokens
	// (StringLiteralToken.hasUnescapedUnicodeLineOrParagraphSeparator).
	HasUnescapedUnicodeLineOrParagraphSeparator bool

	// Template error fields (TemplateLiteralToken).
	ErrorMessage  string
	ErrorLevel    TemplateLiteralErrorLevel
	ErrorPosition SourcePosition
}

// Go-only: language-adapter — GetStart consolidates Token.getStart() from Token.java into this file;
// Java's Token is a separate class not counted in TokenType.java's method surface.
func (t *Token) GetStart() SourcePosition { return t.Location.Start }

// ToString mirrors Token.toString() — returns the token type's string representation.
// Go-only: language-adapter — Token.java is a secondary source; the name does not appear in TokenType.java.
func (t *Token) ToString() string { return t.Type.ToString() }

// Go-only: language-adapter — AsLiteral consolidates Token.asLiteral() from Token.java into this file.
func (t *Token) AsLiteral() *Token { return t }

// Go-only: language-adapter — AsIdentifier consolidates Token.asIdentifier() from Token.java into this file.
func (t *Token) AsIdentifier() *Token { return t }

// Go-only: language-adapter — AsTemplateLiteral consolidates Token.asTemplateLiteral() from Token.java into this file.
func (t *Token) AsTemplateLiteral() *Token { return t }

// Go-only: language-adapter — IsPrivateIdentifier has no Java equivalent in TokenType.java or Token.java;
// private-identifier detection is handled inline in the Java scanner.
func (t *Token) IsPrivateIdentifier() bool {
	return t.Type == TokenTypeIdentifier && len(t.Value) > 0 && t.Value[0] == '#'
}

// Go-only: language-adapter — HasError consolidates TemplateLiteralToken.hasError() from TemplateLiteralToken.java.
func (t *Token) HasError() bool {
	return t.ErrorLevel == TemplateLiteralError && t.ErrorMessage != ""
}

// TemplateLiteralErrorLevel mirrors TemplateLiteralToken.ErrorLevel.
type TemplateLiteralErrorLevel int

const (
	TemplateLiteralWarning TemplateLiteralErrorLevel = iota
	TemplateLiteralError
)

// ---- Convenience constructors ---------------------------------------

// Go-only: language-adapter — newToken replaces Java's `new Token(type, location)` constructor.
func newToken(typ TokenType, location SourceRange) *Token {
	return &Token{Type: typ, Location: location}
}

// Go-only: language-adapter — newLiteralToken replaces Java's `new LiteralToken(type, value, location)` constructor.
func newLiteralToken(typ TokenType, value string, location SourceRange) *Token {
	return &Token{Type: typ, Location: location, Value: value}
}

// Go-only: language-adapter — newStringLiteralToken replaces Java's `new StringLiteralToken(value, location, hasUnescaped)` constructor.
func newStringLiteralToken(value string, location SourceRange, hasUnescaped bool) *Token {
	return &Token{
		Type:     TokenTypeString,
		Location: location,
		Value:    value,
		HasUnescapedUnicodeLineOrParagraphSeparator: hasUnescaped,
	}
}

// Go-only: language-adapter — newIdentifierToken replaces Java's `new IdentifierToken(location, value)` constructor.
func newIdentifierToken(location SourceRange, value string) *Token {
	return &Token{Type: TokenTypeIdentifier, Location: location, Value: value}
}

// Go-only: language-adapter — newTemplateLiteralToken replaces Java's `new TemplateLiteralToken(...)` constructor.
func newTemplateLiteralToken(
	typ TokenType,
	value string,
	errorMsg string,
	errorLevel TemplateLiteralErrorLevel,
	position SourcePosition,
	location SourceRange,
) *Token {
	return &Token{
		Type:          typ,
		Location:      location,
		Value:         value,
		ErrorMessage:  errorMsg,
		ErrorLevel:    errorLevel,
		ErrorPosition: position,
	}
}

// ---- Legacy type aliases (kept for callers that used the old separate types) ---

// LiteralToken and IdentifierToken are now just *Token — aliases for documentation.
// The scan methods that used to return these types now return *Token.
type LiteralToken = Token
type IdentifierToken = Token
type StringLiteralToken = Token
type TemplateLiteralToken = Token

// ---- Keywords -------------------------------------------------------

// keywordEntry holds the string → TokenType mapping, mirroring Keywords.java.
type keywordEntry struct {
	value string
	typ   TokenType
}

// keywords is the full set from Keywords.java, in the same order.
var keywords = []keywordEntry{
	// 7.6.1.1 Keywords
	{"break", TokenTypeBreak},
	{"case", TokenTypeCase},
	{"catch", TokenTypeCatch},
	{"continue", TokenTypeContinue},
	{"debugger", TokenTypeDebugger},
	{"default", TokenTypeDefault},
	{"delete", TokenTypeDelete},
	{"do", TokenTypeDo},
	{"else", TokenTypeElse},
	{"finally", TokenTypeFinally},
	{"for", TokenTypeFor},
	{"function", TokenTypeFunction},
	{"if", TokenTypeIf},
	{"in", TokenTypeIn},
	{"instanceof", TokenTypeInstanceof},
	{"new", TokenTypeNew},
	{"return", TokenTypeReturn},
	{"switch", TokenTypeSwitch},
	{"this", TokenTypeThis},
	{"throw", TokenTypeThrow},
	{"try", TokenTypeTry},
	{"typeof", TokenTypeTypeof},
	{"var", TokenTypeVar},
	{"void", TokenTypeVoid},
	{"while", TokenTypeWhile},
	{"with", TokenTypeWith},
	// 7.6.1.2 Future Reserved Words
	{"class", TokenTypeClass},
	{"const", TokenTypeConst},
	{"enum", TokenTypeEnum},
	{"export", TokenTypeExport},
	{"extends", TokenTypeExtends},
	{"import", TokenTypeImport},
	{"super", TokenTypeSuper},
	// Future Reserved Words in a strict context
	{"implements", TokenTypeImplements},
	{"interface", TokenTypeInterface},
	{"let", TokenTypeLet},
	{"package", TokenTypePackage},
	{"private", TokenTypePrivate},
	{"protected", TokenTypeProtected},
	{"public", TokenTypePublic},
	{"static", TokenTypeStatic},
	{"yield", TokenTypeYield},
	// 7.8 Literals
	{"null", TokenTypeNull},
	{"true", TokenTypeTrue},
	{"false", TokenTypeFalse},
}

// keywordsByName mirrors Keywords.KEYWORDS_BY_NAME.
var keywordsByName map[string]TokenType

// keywordsByType mirrors Keywords.KEYWORDS_BY_TYPE.
var keywordsByType map[TokenType]string

// Go-only: language-adapter — init populates the keyword lookup maps; Java uses a static initializer block.
func init() {
	keywordsByName = make(map[string]TokenType, len(keywords))
	keywordsByType = make(map[TokenType]string, len(keywords))
	for _, kw := range keywords {
		keywordsByName[kw.value] = kw.typ
		keywordsByType[kw.typ] = kw.value
	}
}

// Go-only: language-adapter — getKeywordType ports Keywords.get(String) into the parser package;
// Java resolves this via a separate Keywords enum.
func getKeywordType(value string) (TokenType, bool) {
	t, ok := keywordsByName[value]
	return t, ok
}

// Go-only: language-adapter — keywordForType ports Keywords.get(TokenType) into the parser package;
// Java resolves this via a separate Keywords enum.
func keywordForType(t TokenType) string {
	return keywordsByType[t]
}

// Go-only: language-adapter — isKeywordType ports Keywords.isKeyword(TokenType) into the parser package;
// Java resolves this via a separate Keywords enum.
func isKeywordType(t TokenType) bool {
	_, ok := keywordsByType[t]
	return ok
}

// Go-only: language-adapter — isStrictKeyword ports Keywords.isStrictKeyword(TokenType) into the parser package;
// Java resolves this via a separate Keywords enum.
func isStrictKeyword(t TokenType) bool {
	switch t {
	case TokenTypeImplements, TokenTypeInterface, TokenTypeLet, TokenTypePackage,
		TokenTypePrivate, TokenTypeProtected, TokenTypePublic, TokenTypeStatic, TokenTypeYield:
		return true
	}
	return false
}
