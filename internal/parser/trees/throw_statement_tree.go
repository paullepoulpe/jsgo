// Port of com/google/javascript/jscomp/parsing/parser/trees/ThrowStatementTree.java
// @ closure-compiler v20260415

package trees

// ThrowStatementTree represents a throw statement.
type ThrowStatementTree struct {
	ParseTree
	Value *ParseTree
}

func NewThrowStatementTree(location SourceRange, value *ParseTree) *ThrowStatementTree {
	return &ThrowStatementTree{
		ParseTree: NewParseTree(ParseTreeTypeThrowStatement, location),
		Value:     value,
	}
}
