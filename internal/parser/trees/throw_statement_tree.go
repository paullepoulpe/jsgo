// Port of com/google/javascript/jscomp/parsing/parser/trees/ThrowStatementTree.java
// @ closure-compiler v20260415

package trees

type ThrowStatementTree struct {
	ParseTree
	Value *ParseTree
}

func NewThrowStatementTree(location SourceRange, value *ParseTree) *ThrowStatementTree {
	return &ThrowStatementTree{
		ParseTree: ParseTree{Type: ParseTreeTypeThrowStatement, Location: location},
		Value:     value,
	}
}
