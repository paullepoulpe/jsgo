// Port of com/google/javascript/jscomp/parsing/parser/trees/ComprehensionTree.java
// @ closure-compiler v20260415

package trees

// ComprehensionType indicates whether a comprehension is an array or generator comprehension.
type ComprehensionType int

const (
	ComprehensionTypeArray     ComprehensionType = iota
	ComprehensionTypeGenerator
)

// ComprehensionTree represents an array or generator comprehension.
type ComprehensionTree struct {
	ParseTree
	ComprehensionType ComprehensionType
	Children          []*ParseTree
	TailExpression    *ParseTree
}

// NewComprehensionTree constructs a ComprehensionTree.
// Go-only: language-adapter
func NewComprehensionTree(location SourceRange, comprehensionType ComprehensionType, children []*ParseTree, tailExpression *ParseTree) *ComprehensionTree {
	return &ComprehensionTree{
		ParseTree:         ParseTree{Type: ParseTreeTypeComprehension, Location: location},
		ComprehensionType: comprehensionType,
		Children:          children,
		TailExpression:    tailExpression,
	}
}
