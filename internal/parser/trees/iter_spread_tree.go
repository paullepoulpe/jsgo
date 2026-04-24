// Port of com/google/javascript/jscomp/parsing/parser/trees/IterSpreadTree.java
// @ closure-compiler v20260415

package trees

// IterSpreadTree represents an iterable spread element (...expression).
type IterSpreadTree struct {
	ParseTree
	Expression *ParseTree
}

func NewIterSpreadTree(location SourceRange, expression *ParseTree) *IterSpreadTree {
	return &IterSpreadTree{
		ParseTree:  NewParseTree(ParseTreeTypeIterSpread, location),
		Expression: expression,
	}
}
