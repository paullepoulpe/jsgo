// Port of com/google/javascript/jscomp/parsing/parser/trees/ProgramTree.java
// @ closure-compiler v20260415

package trees

// ProgramTree is the root node of a parsed JavaScript program.
type ProgramTree struct {
	ParseTree
	SourceElements []*ParseTree
	SourceComments []Comment
}

// NewProgramTree constructs a ProgramTree.
func NewProgramTree(location SourceRange, sourceElements []*ParseTree, sourceComments []Comment) *ProgramTree {
	return &ProgramTree{
		ParseTree:      ParseTree{Type: ParseTreeTypeProgram, Location: location},
		SourceElements: sourceElements,
		SourceComments: sourceComments,
	}
}
