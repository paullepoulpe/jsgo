// Port of com/google/javascript/jscomp/parsing/parser/trees/ProgramTree.java
// @ closure-compiler v20260415

package trees

type ProgramTree struct {
	SourceElements []ParseTree
	SourceComments []Comment
}

func NewProgramTree(location SourceRange, sourceElements []ParseTree, sourceComments []Comment) *ProgramTree {
	return &ProgramTree{
		SourceElements: sourceElements,
		SourceComments: sourceComments,
	}
}
