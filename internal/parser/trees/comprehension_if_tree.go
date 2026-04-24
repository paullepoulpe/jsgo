// Port of com/google/javascript/jscomp/parsing/parser/trees/ComprehensionIfTree.java
// @ closure-compiler v20260415

package trees

// ComprehensionIfTree represents an if clause in a comprehension, e.g. if (condition).
type ComprehensionIfTree struct {
	ParseTree
	Expression *ParseTree
}

// NewComprehensionIfTree constructs a ComprehensionIfTree.
func NewComprehensionIfTree(location SourceRange, expression *ParseTree) *ComprehensionIfTree {
	return &ComprehensionIfTree{
		ParseTree:  ParseTree{Type: ParseTreeTypeComprehensionIf, Location: location},
		Expression: expression,
	}
}
