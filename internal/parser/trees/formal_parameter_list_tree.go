// Port of com/google/javascript/jscomp/parsing/parser/trees/FormalParameterListTree.java
// @ closure-compiler v20260415

package trees

// FormalParameterListTree represents a function's formal parameter list.
type FormalParameterListTree struct {
	ParseTree
	Parameters       []*ParseTree
	HasTrailingComma bool
	CommaPositions   []SourcePosition
}

func NewFormalParameterListTree(
	location SourceRange,
	parameters []*ParseTree,
	hasTrailingComma bool,
	commaPositions []SourcePosition,
) *FormalParameterListTree {
	return &FormalParameterListTree{
		ParseTree:        NewParseTree(ParseTreeTypeFormalParameterList, location),
		Parameters:       parameters,
		HasTrailingComma: hasTrailingComma,
		CommaPositions:   commaPositions,
	}
}
