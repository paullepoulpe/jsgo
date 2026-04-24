// Port of com/google/javascript/jscomp/parsing/parser/trees/ProgramTree.java
// @ closure-compiler v20260415

package trees

// ProgramTree is the root of a parsed JavaScript program.
type ProgramTree struct {
	ParseTree
	SourceElements []*ParseTree
	SourceComments []*Comment
}

func NewProgramTree(
	location SourceRange,
	sourceElements []*ParseTree,
	sourceComments []*Comment,
) *ProgramTree {
	return &ProgramTree{
		ParseTree:      NewParseTree(ParseTreeTypeProgram, location),
		SourceElements: sourceElements,
		SourceComments: sourceComments,
	}
}
