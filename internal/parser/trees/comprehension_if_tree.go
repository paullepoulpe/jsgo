// Port of com/google/javascript/jscomp/parsing/parser/trees/ComprehensionIfTree.java
// @ closure-compiler v20260415

package trees

// ComprehensionIfTree represents an if clause in a comprehension.
type ComprehensionIfTree struct {
	ParseTree
	Expression *ParseTree
}

func NewComprehensionIfTree(location SourceRange, expression *ParseTree) *ComprehensionIfTree {
	return &ComprehensionIfTree{
		ParseTree:  NewParseTree(ParseTreeTypeComprehensionIf, location),
		Expression: expression,
	}
}
