// Port of com/google/javascript/jscomp/parsing/parser/trees/ObjectPatternTree.java
// @ closure-compiler v20260415

package trees

// ObjectPatternTree represents an object destructuring pattern.
type ObjectPatternTree struct {
	ParseTree
	Fields []*ParseTree
}

func NewObjectPatternTree(location SourceRange, fields []*ParseTree) *ObjectPatternTree {
	return &ObjectPatternTree{
		ParseTree: NewParseTree(ParseTreeTypeObjectPattern, location),
		Fields:    fields,
	}
}
