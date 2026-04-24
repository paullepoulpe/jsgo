// Port of com/google/javascript/jscomp/parsing/parser/trees/ComprehensionForTree.java
// @ closure-compiler v20260415

package trees

// ComprehensionForTree represents a for-of clause in a comprehension.
type ComprehensionForTree struct {
	ParseTree
	Initializer *ParseTree
	Collection  *ParseTree
}

func NewComprehensionForTree(
	location SourceRange,
	initializer *ParseTree,
	collection *ParseTree,
) *ComprehensionForTree {
	return &ComprehensionForTree{
		ParseTree:   NewParseTree(ParseTreeTypeComprehensionFor, location),
		Initializer: initializer,
		Collection:  collection,
	}
}
