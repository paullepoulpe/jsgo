// Port of com/google/javascript/jscomp/parsing/parser/trees/TemplateLiteralExpressionTree.java
// @ closure-compiler v20260415

package trees

// TemplateLiteralExpressionTree represents a template literal expression (ES6).
// When Operand is non-nil this is a tagged template literal.
type TemplateLiteralExpressionTree struct {
	ParseTree
	// Operand is nil for untagged template literals.
	Operand  *ParseTree
	Elements []*ParseTree
}

func NewTemplateLiteralExpressionTree(
	location SourceRange,
	operand *ParseTree,
	elements []*ParseTree,
) *TemplateLiteralExpressionTree {
	return &TemplateLiteralExpressionTree{
		ParseTree: NewParseTree(ParseTreeTypeTemplateLiteralExpression, location),
		Operand:   operand,
		Elements:  elements,
	}
}
