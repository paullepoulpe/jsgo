// Port of com/google/javascript/jscomp/parsing/parser/trees/NewTargetExpressionTree.java @ closure-compiler v20260415

package parser

// Represents 'new.target'.
type NewTargetExpressionTree struct {
	ParseTreeBase
}

// NewNewTargetExpressionTree constructs a NewTargetExpressionTree.
//
// Java: NewTargetExpressionTree(SourceRange)
func NewNewTargetExpressionTree(location *SourceRange) *NewTargetExpressionTree {
	t := &NewTargetExpressionTree{}
	t.InitParseTreeBase(t, ParseTreeType_NEW_TARGET_EXPRESSION, location)
	return t
}
