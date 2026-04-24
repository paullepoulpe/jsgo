// Port of com/google/javascript/jscomp/parsing/parser/trees/ComprehensionForTree.java
// @ closure-compiler v20260415

package trees

// ComprehensionForTree represents a for clause in a comprehension, e.g. for (x of iterable).
type ComprehensionForTree struct {
	ParseTree
	Initializer *ParseTree
	Collection  *ParseTree
}

// NewComprehensionForTree constructs a ComprehensionForTree.
func NewComprehensionForTree(location SourceRange, initializer *ParseTree, collection *ParseTree) *ComprehensionForTree {
	return &ComprehensionForTree{
		ParseTree:   ParseTree{Type: ParseTreeTypeComprehensionFor, Location: location},
		Initializer: initializer,
		Collection:  collection,
	}
}
