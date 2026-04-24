// Port of com/google/javascript/jscomp/parsing/parser/trees/ParseTree.java
// @ closure-compiler v20260415

package trees

// SourcePosition represents an offset within a source file.
// Stub — full implementation lives in a future port of
// com/google/javascript/jscomp/parsing/parser/util/SourcePosition.java.
type SourcePosition struct {
	Line   int
	Column int
	Offset int
}

// SourceRange is a half-open [Start, End) span within a source file.
// Stub — full implementation lives in a future port of
// com/google/javascript/jscomp/parsing/parser/util/SourceRange.java.
type SourceRange struct {
	Start SourcePosition
	End   SourcePosition
}

// Token is the base token type emitted by the scanner.
// Mirrors com/google/javascript/jscomp/parsing/parser/Token.java.
type Token struct {
	Type     TokenType
	Location SourceRange
	// Value holds the string value for identifier, literal, and string tokens.
	// Empty for punctuator/keyword tokens.
	// Mirrors Token.java's LiteralToken / IdentifierToken value field.
	Value string
}

// TokenType identifies the lexical kind of a Token.
// Stub — full implementation lives in a future port of
// com/google/javascript/jscomp/parsing/parser/TokenType.java.
//
// The integer values below are kept in sync with
// internal/parser.TokenType (the authoritative source). The conversion
// trees.TokenType(int(parserToken.Type)) is used by parser_new.go to bridge
// the two packages without creating an import cycle.
type TokenType int

const (
	// sentinel / error
	TokenTypeEndOfFile TokenType = iota // 0
	TokenTypeError                      // 1
	// 7.6 Identifier Names and Identifiers
	TokenTypeIdentifier // 2
	// 7.6.1.1 Keywords
	TokenTypeBreak    // 3
	TokenTypeCase     // 4
	TokenTypeCatch    // 5
	TokenTypeContinue // 6
	TokenTypeDebugger // 7
	TokenTypeDefault  // 8
	TokenTypeDelete   // 9
	TokenTypeDo       // 10
	TokenTypeElse     // 11
	TokenTypeFinally  // 12
	TokenTypeFor      // 13
	TokenTypeFunction // 14
	TokenTypeIf       // 15
	TokenTypeIn       // 16
	TokenTypeInstanceof // 17
	TokenTypeNew      // 18
	TokenTypeReturn   // 19
	TokenTypeSwitch   // 20
	TokenTypeThis     // 21
	TokenTypeThrow    // 22
	TokenTypeTry      // 23
	TokenTypeTypeof   // 24
	TokenTypeVar      // 25
	TokenTypeVoid     // 26
	TokenTypeWhile    // 27
	TokenTypeWith     // 28
	// 7.6.1.2 Future reserved words
	TokenTypeClass     // 29
	TokenTypeConst     // 30
	TokenTypeEnum      // 31
	TokenTypeExport    // 32
	TokenTypeExtends   // 33
	TokenTypeImport    // 34
	TokenTypeSuper     // 35
	// Future reserved words in strict mode
	TokenTypeImplements // 36
	TokenTypeInterface  // 37
	TokenTypeLet        // 38
	TokenTypePackage    // 39
	TokenTypePrivate    // 40
	TokenTypeProtected  // 41
	TokenTypePublic     // 42
	TokenTypeStatic     // 43
	TokenTypeYield      // 44
	// 7.7 Punctuators
	TokenTypeOpenCurly            // 45
	TokenTypeCloseCurly           // 46
	TokenTypeOpenParen            // 47
	TokenTypeCloseParen           // 48
	TokenTypeOpenSquare           // 49
	TokenTypeCloseSquare          // 50
	TokenTypePeriod               // 51
	TokenTypeSemiColon            // 52
	TokenTypeComma                // 53
	TokenTypeOpenAngle            // 54
	TokenTypeCloseAngle           // 55
	TokenTypeLessEqual            // 56
	TokenTypeGreaterEqual         // 57
	TokenTypeArrow                // 58
	TokenTypeEqualEqual           // 59
	TokenTypeNotEqual             // 60
	TokenTypeEqualEqualEqual      // 61
	TokenTypeNotEqualEqual        // 62
	TokenTypePlus                 // 63
	TokenTypeMinus                // 64
	TokenTypeStar                 // 65
	TokenTypeStarStar             // 66
	TokenTypePercent              // 67
	TokenTypePlusPlus             // 68
	TokenTypeMinusMinus           // 69
	TokenTypeLeftShift            // 70
	TokenTypeRightShift           // 71
	TokenTypeUnsignedRightShift   // 72
	TokenTypeAmpersand            // 73
	TokenTypeBar                  // 74
	TokenTypeCaret                // 75
	TokenTypeBang                 // 76
	TokenTypeTilde                // 77
	TokenTypeAnd                  // 78
	TokenTypeOr                   // 79
	TokenTypeQuestion             // 80
	TokenTypeQuestionQuestion     // 81
	TokenTypeQuestionDot          // 82
	TokenTypeColon                // 83
	TokenTypeEqual                // 84
	TokenTypePlusEqual            // 85
	TokenTypeMinusEqual           // 86
	TokenTypeStarEqual            // 87
	TokenTypeStarStarEqual        // 88
	TokenTypePercentEqual         // 89
	TokenTypeLeftShiftEqual       // 90
	TokenTypeRightShiftEqual      // 91
	TokenTypeUnsignedRightShiftEqual // 92
	TokenTypeAmpersandEqual       // 93
	TokenTypeBarEqual             // 94
	TokenTypeCaretEqual           // 95
	TokenTypeSlash                // 96
	TokenTypeSlashEqual           // 97
	// Logical Assignment Punctuators
	TokenTypeAndEqual             // 98
	TokenTypeOrEqual              // 99
	TokenTypeQuestionQuestionEqual // 100
	// 7.8 Literals
	TokenTypeNull                // 101
	TokenTypeTrue                // 102
	TokenTypeFalse               // 103
	TokenTypeNumber              // 104
	TokenTypeString              // 105
	TokenTypeBigInt              // 106
	TokenTypeRegularExpression   // 107
	// Harmony extensions
	TokenTypeEllipsis            // 108
	// 12.2.9 Template Literals
	TokenTypeTemplateHead           // 109
	TokenTypeTemplateMiddle         // 110
	TokenTypeTemplateTail           // 111
	TokenTypeNoSubstitutionTemplate // 112
	// TypeScript
	TokenTypeType      // 113
	TokenTypeDeclare   // 114
	TokenTypeModule    // 115
	TokenTypeNamespace // 116
)

// IdentifierToken is a Token subclass carrying the identifier string.
// Stub — full implementation lives in a future port of
// com/google/javascript/jscomp/parsing/parser/IdentifierToken.java.
type IdentifierToken struct {
	Token
	Value string
}

// ParseTree is the abstract base for all parse-tree nodes.
// Immutable plain-old-data; derived types add fields but no logic.
//
// Port of com/google/javascript/jscomp/parsing/parser/trees/ParseTree.java
// @ closure-compiler v20260415
type ParseTree struct {
	Type     ParseTreeType
	Location SourceRange
}

// NewParseTree constructs a ParseTree node with the given type and location.
func NewParseTree(typ ParseTreeType, location SourceRange) *ParseTree {
	return &ParseTree{Type: typ, Location: location}
}

// LiteralToken is a stub placeholder for the full Java LiteralToken class.
// The scanner emits Token instances carrying literal values; once the
// parser layer is migrated to consume internal/parser's token stream, this
// alias will be replaced by a direct reference to parser.Token.
type LiteralToken struct {
	Token
}
