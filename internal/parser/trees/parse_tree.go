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
// Stub — full implementation lives in a future port of
// com/google/javascript/jscomp/parsing/parser/Token.java.
type Token struct {
	Type     TokenType
	Location SourceRange
}

// TokenType identifies the lexical kind of a Token.
// Stub — full implementation lives in a future port of
// com/google/javascript/jscomp/parsing/parser/TokenType.java.
type TokenType int

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
