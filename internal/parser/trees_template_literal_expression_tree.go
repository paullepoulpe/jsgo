// Port of com/google/javascript/jscomp/parsing/parser/trees/TemplateLiteralExpressionTree.java @ closure-compiler v20260415

package parser

// Template literal production in ES6.
type TemplateLiteralExpressionTree struct {
	ParseTreeBase

	Operand  ParseTree
	Elements []ParseTree
}

// NewTemplateLiteralExpressionTree constructs a TemplateLiteralExpressionTree.
//
// Java: TemplateLiteralExpressionTree(SourceRange, ParseTree, ImmutableList)
func NewTemplateLiteralExpressionTree(location *SourceRange, operand ParseTree, elements []ParseTree) *TemplateLiteralExpressionTree {
	t := &TemplateLiteralExpressionTree{}
	t.InitParseTreeBase(t, ParseTreeType_TEMPLATE_LITERAL_EXPRESSION, location)
	t.Operand = operand
	t.Elements = elements
	return t
}
