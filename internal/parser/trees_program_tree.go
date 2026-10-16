// Port of com/google/javascript/jscomp/parsing/parser/trees/ProgramTree.java @ closure-compiler v20260415

package parser

// ProgramTree is the parse tree for ParseTreeType PROGRAM.
type ProgramTree struct {
	ParseTreeBase

	SourceElements []ParseTree
	SourceComments []*Comment
}

// NewProgramTree constructs a ProgramTree.
//
// Java: ProgramTree(SourceRange, ImmutableList, ImmutableList)
func NewProgramTree(location *SourceRange, sourceElements []ParseTree, sourceComments []*Comment) *ProgramTree {
	t := &ProgramTree{}
	t.InitParseTreeBase(t, ParseTreeType_PROGRAM, location)
	t.SourceElements = sourceElements
	t.SourceComments = sourceComments
	return t
}
