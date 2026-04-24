// Port of com/google/javascript/jscomp/parsing/parser/trees/ObjectSpreadTree.java
// @ closure-compiler v20260415

package trees

// ObjectSpreadTree represents an object spread element ({...expression}).
type ObjectSpreadTree struct {
	ParseTree
	Expression *ParseTree
}

func NewObjectSpreadTree(location SourceRange, expression *ParseTree) *ObjectSpreadTree {
	return &ObjectSpreadTree{
		ParseTree:  NewParseTree(ParseTreeTypeObjectSpread, location),
		Expression: expression,
	}
}
