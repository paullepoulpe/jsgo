// Port of com/google/javascript/jscomp/parsing/parser/trees/ThisExpressionTree.java @ closure-compiler v20260415

package parser

// ThisExpressionTree is the parse tree for ParseTreeType THIS_EXPRESSION.
type ThisExpressionTree struct {
	ParseTreeBase
}

// NewThisExpressionTree constructs a ThisExpressionTree.
//
// Java: ThisExpressionTree(SourceRange)
func NewThisExpressionTree(location *SourceRange) *ThisExpressionTree {
	t := &ThisExpressionTree{}
	t.InitParseTreeBase(t, ParseTreeType_THIS_EXPRESSION, location)
	return t
}
