// Port of com/google/javascript/jscomp/parsing/parser/trees/SuperExpressionTree.java @ closure-compiler v20260415

package parser

// SuperExpressionTree is the parse tree for ParseTreeType SUPER_EXPRESSION.
type SuperExpressionTree struct {
	ParseTreeBase
}

// NewSuperExpressionTree constructs a SuperExpressionTree.
//
// Java: SuperExpressionTree(SourceRange)
func NewSuperExpressionTree(location *SourceRange) *SuperExpressionTree {
	t := &SuperExpressionTree{}
	t.InitParseTreeBase(t, ParseTreeType_SUPER_EXPRESSION, location)
	return t
}
