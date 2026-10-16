// Port of com/google/javascript/jscomp/parsing/parser/trees/FormalParameterListTree.java
// @ closure-compiler v20260415

package trees

// FormalParameterListTree represents the formal parameter list of a function.
type FormalParameterListTree struct {
	ParseTree
	Parameters       []*ParseTree
	HasTrailingComma bool
	CommaPositions   []SourcePosition
}

// NewFormalParameterListTree constructs a FormalParameterListTree.
// Go-only: language-adapter
func NewFormalParameterListTree(location SourceRange, parameters []*ParseTree, hasTrailingComma bool, commaPositions []SourcePosition) *FormalParameterListTree {
	return &FormalParameterListTree{
		ParseTree:        ParseTree{Type: ParseTreeTypeFormalParameterList, Location: location},
		Parameters:       parameters,
		HasTrailingComma: hasTrailingComma,
		CommaPositions:   commaPositions,
	}
}
