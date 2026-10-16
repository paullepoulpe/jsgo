// Port of com/google/javascript/jscomp/parsing/parser/trees/ObjectSpreadTree.java
// @ closure-compiler v20260415

package trees

// ObjectSpreadTree represents an object spread element, e.g. { ...expr }.
type ObjectSpreadTree struct {
	ParseTree
	Expression *ParseTree
}

// NewObjectSpreadTree constructs an ObjectSpreadTree.
// Go-only: language-adapter
func NewObjectSpreadTree(location SourceRange, expression *ParseTree) *ObjectSpreadTree {
	return &ObjectSpreadTree{
		ParseTree:  ParseTree{Type: ParseTreeTypeObjectSpread, Location: location},
		Expression: expression,
	}
}
