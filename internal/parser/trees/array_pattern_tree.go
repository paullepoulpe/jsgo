// Port of com/google/javascript/jscomp/parsing/parser/trees/ArrayPatternTree.java
// @ closure-compiler v20260415

package trees

// ArrayPatternTree represents an array destructuring pattern.
type ArrayPatternTree struct {
	ParseTree
	Elements []*ParseTree
}

func NewArrayPatternTree(location SourceRange, elements []*ParseTree) *ArrayPatternTree {
	return &ArrayPatternTree{
		ParseTree: NewParseTree(ParseTreeTypeArrayPattern, location),
		Elements:  elements,
	}
}
