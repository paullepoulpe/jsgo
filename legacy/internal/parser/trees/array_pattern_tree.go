// Port of com/google/javascript/jscomp/parsing/parser/trees/ArrayPatternTree.java
// @ closure-compiler v20260415

package trees

// ArrayPatternTree represents a destructuring array pattern, e.g. [a, b].
type ArrayPatternTree struct {
	ParseTree
	Elements []*ParseTree
}

// NewArrayPatternTree constructs an ArrayPatternTree.
// Go-only: language-adapter
func NewArrayPatternTree(location SourceRange, elements []*ParseTree) *ArrayPatternTree {
	return &ArrayPatternTree{
		ParseTree: ParseTree{Type: ParseTreeTypeArrayPattern, Location: location},
		Elements:  elements,
	}
}
