// Port of com/google/javascript/jscomp/parsing/parser/trees/ObjectSpreadTree.java @ closure-compiler v20260415

package parser

// ObjectSpreadTree is the parse tree for ParseTreeType OBJECT_SPREAD.
type ObjectSpreadTree struct {
	ParseTreeBase

	Expression ParseTree
}

// NewObjectSpreadTree constructs a ObjectSpreadTree.
//
// Java: ObjectSpreadTree(SourceRange, ParseTree)
func NewObjectSpreadTree(location *SourceRange, expression ParseTree) *ObjectSpreadTree {
	t := &ObjectSpreadTree{}
	t.InitParseTreeBase(t, ParseTreeType_OBJECT_SPREAD, location)
	t.Expression = expression
	return t
}
