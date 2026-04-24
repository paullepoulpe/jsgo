// Port of com/google/javascript/jscomp/parsing/parser/trees/ForInStatementTree.java
// @ closure-compiler v20260415

package trees

type ForInStatementTree struct {
	Initializer ParseTree
	Collection  ParseTree
	Body        ParseTree
}

func NewForInStatementTree(location SourceRange, initializer ParseTree, collection ParseTree, body ParseTree) *ForInStatementTree {
	return &ForInStatementTree{
		Initializer: initializer,
		Collection:  collection,
		Body:        body,
	}
}
