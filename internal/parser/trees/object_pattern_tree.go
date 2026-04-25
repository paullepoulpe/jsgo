// Port of com/google/javascript/jscomp/parsing/parser/trees/ObjectPatternTree.java
// @ closure-compiler v20260415

package trees

// ObjectPatternTree represents a destructuring object pattern, e.g. { a, b }.
type ObjectPatternTree struct {
	ParseTree
	Fields []*ParseTree
}

// NewObjectPatternTree constructs an ObjectPatternTree.
// Go-only: language-adapter
func NewObjectPatternTree(location SourceRange, fields []*ParseTree) *ObjectPatternTree {
	return &ObjectPatternTree{
		ParseTree: ParseTree{Type: ParseTreeTypeObjectPattern, Location: location},
		Fields:    fields,
	}
}
