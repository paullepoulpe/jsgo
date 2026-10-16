// Port of com/google/javascript/jscomp/parsing/parser/trees/DefaultParameterTree.java @ closure-compiler v20260415

package parser

// DefaultParameterTree is the parse tree for ParseTreeType DEFAULT_PARAMETER.
type DefaultParameterTree struct {
	ParseTreeBase

	Lhs          ParseTree
	DefaultValue ParseTree
}

// NewDefaultParameterTree constructs a DefaultParameterTree.
//
// Java: DefaultParameterTree(SourceRange, ParseTree, ParseTree)
func NewDefaultParameterTree(location *SourceRange, lhs ParseTree, defaultValue ParseTree) *DefaultParameterTree {
	t := &DefaultParameterTree{}
	t.InitParseTreeBase(t, ParseTreeType_DEFAULT_PARAMETER, location)
	t.Lhs = lhs
	t.DefaultValue = defaultValue
	return t
}
