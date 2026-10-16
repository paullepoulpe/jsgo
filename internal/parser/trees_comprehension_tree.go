// Port of com/google/javascript/jscomp/parsing/parser/trees/ComprehensionTree.java @ closure-compiler v20260415

package parser

// ComprehensionTree is the parse tree for ParseTreeType COMPREHENSION.
type ComprehensionTree struct {
	ParseTreeBase

	// Hides ParseTreeBase.Type, exactly like the Java field hides ParseTree.type.
	Type ComprehensionTree_ComprehensionType
	// The "for of" and "if" trees.
	Children []ParseTree
	// The AssignmentExpression at the end of the ComprehensionTail.
	TailExpression ParseTree
}

// ComprehensionTree_ComprehensionType is ComprehensionTree.ComprehensionType.
type ComprehensionTree_ComprehensionType int32

const (
	ComprehensionTree_ComprehensionType_ARRAY ComprehensionTree_ComprehensionType = iota
	ComprehensionTree_ComprehensionType_GENERATOR
)

var comprehensionTree_ComprehensionType_names = [...]string{"ARRAY", "GENERATOR"}

// ComprehensionTree_ComprehensionType_values is ComprehensionType.values().
func ComprehensionTree_ComprehensionType_values() []ComprehensionTree_ComprehensionType {
	return []ComprehensionTree_ComprehensionType{ComprehensionTree_ComprehensionType_ARRAY, ComprehensionTree_ComprehensionType_GENERATOR}
}

// Name is Enum.name().
func (t ComprehensionTree_ComprehensionType) Name() string {
	return comprehensionTree_ComprehensionType_names[t]
}

// String is Enum.toString().
func (t ComprehensionTree_ComprehensionType) String() string { return t.Name() }

// NewComprehensionTree constructs a ComprehensionTree.
//
// Java: ComprehensionTree(SourceRange, ComprehensionType, ImmutableList, ParseTree)
func NewComprehensionTree(location *SourceRange, type_ ComprehensionTree_ComprehensionType, children []ParseTree, tailExpression ParseTree) *ComprehensionTree {
	t := &ComprehensionTree{}
	t.InitParseTreeBase(t, ParseTreeType_COMPREHENSION, location)
	t.Type = type_
	t.Children = children
	t.TailExpression = tailExpression
	return t
}
