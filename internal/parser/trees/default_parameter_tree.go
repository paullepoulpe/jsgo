// Port of com/google/javascript/jscomp/parsing/parser/trees/DefaultParameterTree.java
// @ closure-compiler v20260415

package trees

// DefaultParameterTree represents a parameter with a default value, e.g. x = 0.
type DefaultParameterTree struct {
	ParseTree
	Lhs          *ParseTree
	DefaultValue *ParseTree
}

// NewDefaultParameterTree constructs a DefaultParameterTree.
// Go-only: language-adapter
func NewDefaultParameterTree(location SourceRange, lhs *ParseTree, defaultValue *ParseTree) *DefaultParameterTree {
	return &DefaultParameterTree{
		ParseTree:    ParseTree{Type: ParseTreeTypeDefaultParameter, Location: location},
		Lhs:          lhs,
		DefaultValue: defaultValue,
	}
}
