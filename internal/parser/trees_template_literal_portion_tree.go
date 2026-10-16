// Port of com/google/javascript/jscomp/parsing/parser/trees/TemplateLiteralPortionTree.java @ closure-compiler v20260415

package parser

// A production representing a literal portion of a template literal.
// Eg: `hello`.
type TemplateLiteralPortionTree struct {
	ParseTreeBase

	Value Token
}

// NewTemplateLiteralPortionTree constructs a TemplateLiteralPortionTree.
//
// Java: TemplateLiteralPortionTree(SourceRange, Token)
func NewTemplateLiteralPortionTree(location *SourceRange, value Token) *TemplateLiteralPortionTree {
	t := &TemplateLiteralPortionTree{}
	t.InitParseTreeBase(t, ParseTreeType_TEMPLATE_LITERAL_PORTION, location)
	t.Value = value
	return t
}
