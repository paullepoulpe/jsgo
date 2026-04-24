//go:build !skeletons

// Minimal type stubs for parser/trees types that are ported in sibling PRs
// (ParseTree base + token types). Replace when those ports land.

package trees

// ParseTree is the abstract base for all parse trees.
// Full port is in parse_tree.go (//go:build skeletons until PR #135 lands).
type ParseTree interface{}

// ParseTreeType identifies the concrete type of a ParseTree.
// Full port is in parse_tree_type.go (//go:build skeletons until PR #135 lands).
type ParseTreeType int

// SourceRange represents a range of positions in a source string.
// Full port is in internal/parser/skel/util/sourcerange.go.
type SourceRange struct{}

// SourcePosition represents a position within a source file.
// Full port is in internal/parser/skel/util/sourceposition.go.
type SourcePosition struct{}

// IdentifierToken is a token representing an identifier.
// Full port is in internal/parser/skel/*.go (//go:build skeletons).
type IdentifierToken struct {
	Value string
}

// LiteralToken is a token representing a literal value.
// Full port is in internal/parser/skel/*.go (//go:build skeletons).
type LiteralToken struct {
	Value string
}

// TokenType identifies the type of a lexical token.
// Full port is in internal/parser/skel/*.go (//go:build skeletons).
type TokenType int

// FormalParameterListTree represents a function's formal parameter list.
// Full port is in formal_parameter_list_tree.go (//go:build skeletons until PR #135 lands).
type FormalParameterListTree struct {
	Parameters      []ParseTree
	HasTrailingComma bool
	CommaPositions  []SourcePosition
}

// Comment represents a comment node in the parse tree.
// Full port is in comment.go (//go:build skeletons until PR #135 lands).
type Comment struct {
	Value    string
	Location SourceRange
}
