// Port of com/google/javascript/jscomp/parsing/parser/trees/TemplateLiteralExpressionTree.java
// @ closure-compiler v20260415

package trees

// TemplateLiteralExpressionTree represents a template literal production in ES6.
// Operand is nil for untagged template literals.
type TemplateLiteralExpressionTree struct {
	ParseTree
	Operand  *ParseTree
	Elements []*ParseTree
}

// NewTemplateLiteralExpressionTree constructs a TemplateLiteralExpressionTree.
// Go-only: language-adapter
func NewTemplateLiteralExpressionTree(location SourceRange, operand *ParseTree, elements []*ParseTree) *TemplateLiteralExpressionTree {
	return &TemplateLiteralExpressionTree{
		ParseTree: ParseTree{Type: ParseTreeTypeTemplateLiteralExpression, Location: location},
		Operand:   operand,
		Elements:  elements,
	}
}
