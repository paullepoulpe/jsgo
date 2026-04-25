// Port of com/google/javascript/jscomp/parsing/parser/trees/IterSpreadTree.java
// @ closure-compiler v20260415

package trees

// IterSpreadTree represents an iterable spread element in an array literal or
// argument list, e.g. ...arr.
type IterSpreadTree struct {
	ParseTree
	Expression *ParseTree
}

// NewIterSpreadTree constructs an IterSpreadTree.
// Go-only: language-adapter
func NewIterSpreadTree(location SourceRange, expression *ParseTree) *IterSpreadTree {
	return &IterSpreadTree{
		ParseTree:  ParseTree{Type: ParseTreeTypeIterSpread, Location: location},
		Expression: expression,
	}
}
