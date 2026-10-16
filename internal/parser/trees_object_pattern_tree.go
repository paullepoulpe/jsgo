// Port of com/google/javascript/jscomp/parsing/parser/trees/ObjectPatternTree.java @ closure-compiler v20260415

package parser

// ObjectPatternTree is the parse tree for ParseTreeType OBJECT_PATTERN.
type ObjectPatternTree struct {
	ParseTreeBase

	Fields []ParseTree
}

// NewObjectPatternTree constructs a ObjectPatternTree.
//
// Java: ObjectPatternTree(SourceRange, ImmutableList)
func NewObjectPatternTree(location *SourceRange, fields []ParseTree) *ObjectPatternTree {
	t := &ObjectPatternTree{}
	t.InitParseTreeBase(t, ParseTreeType_OBJECT_PATTERN, location)
	t.Fields = fields
	return t
}
