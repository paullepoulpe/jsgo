// Port of com/google/javascript/jscomp/parsing/parser/trees/ArgumentListTree.java
// @ closure-compiler v20260415

package trees

// ArgumentListTree represents an argument list in a call or new expression.
type ArgumentListTree struct {
	ParseTree
	Arguments        []*ParseTree
	HasTrailingComma bool
	CommaPositions   []SourcePosition
}

func NewArgumentListTree(
	location SourceRange,
	arguments []*ParseTree,
	hasTrailingComma bool,
	commaPositions []SourcePosition,
) *ArgumentListTree {
	return &ArgumentListTree{
		ParseTree:        NewParseTree(ParseTreeTypeArgumentList, location),
		Arguments:        arguments,
		HasTrailingComma: hasTrailingComma,
		CommaPositions:   commaPositions,
	}
}
