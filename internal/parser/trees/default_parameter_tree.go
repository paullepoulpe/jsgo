// Port of com/google/javascript/jscomp/parsing/parser/trees/DefaultParameterTree.java
// @ closure-compiler v20260415

package trees

// DefaultParameterTree represents a parameter with a default value.
type DefaultParameterTree struct {
	ParseTree
	Lhs          *ParseTree
	DefaultValue *ParseTree
}

func NewDefaultParameterTree(
	location SourceRange,
	lhs *ParseTree,
	defaultValue *ParseTree,
) *DefaultParameterTree {
	return &DefaultParameterTree{
		ParseTree:    NewParseTree(ParseTreeTypeDefaultParameter, location),
		Lhs:          lhs,
		DefaultValue: defaultValue,
	}
}
