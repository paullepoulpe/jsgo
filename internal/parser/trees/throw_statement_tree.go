// Port of com/google/javascript/jscomp/parsing/parser/trees/ThrowStatementTree.java
// @ closure-compiler v20260415

package trees

type ThrowStatementTree struct {
	Value ParseTree
}

func NewThrowStatementTree(location SourceRange, value ParseTree) *ThrowStatementTree {
	return &ThrowStatementTree{
		Value: value,
	}
}
