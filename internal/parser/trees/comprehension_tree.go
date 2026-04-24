// Port of com/google/javascript/jscomp/parsing/parser/trees/ComprehensionTree.java
// @ closure-compiler v20260415

package trees

// ComprehensionKind identifies whether a comprehension produces an array or generator.
type ComprehensionKind int

const (
	ComprehensionKindArray     ComprehensionKind = iota
	ComprehensionKindGenerator
)

// ComprehensionTree represents an array or generator comprehension.
type ComprehensionTree struct {
	ParseTree
	// ComprehensionType holds the comprehension kind (ARRAY or GENERATOR).
	// The field is named ComprehensionType to distinguish it from the embedded
	// ParseTree.Type field.
	ComprehensionType ComprehensionKind
	// Children holds the for-of and if clauses.
	Children       []*ParseTree
	TailExpression *ParseTree
}

func NewComprehensionTree(
	location SourceRange,
	comprehensionType ComprehensionKind,
	children []*ParseTree,
	tailExpression *ParseTree,
) *ComprehensionTree {
	return &ComprehensionTree{
		ParseTree:         NewParseTree(ParseTreeTypeComprehension, location),
		ComprehensionType: comprehensionType,
		Children:          children,
		TailExpression:    tailExpression,
	}
}
