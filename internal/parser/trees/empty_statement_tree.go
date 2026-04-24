// Port of com/google/javascript/jscomp/parsing/parser/trees/EmptyStatementTree.java
// @ closure-compiler v20260415

package trees

type EmptyStatementTree struct{}

func NewEmptyStatementTree(location SourceRange) *EmptyStatementTree {
	return &EmptyStatementTree{}
}
