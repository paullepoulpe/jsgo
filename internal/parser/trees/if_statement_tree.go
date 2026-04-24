// Port of com/google/javascript/jscomp/parsing/parser/trees/IfStatementTree.java
// @ closure-compiler v20260415

package trees

type IfStatementTree struct {
	ParseTree
	Condition  *ParseTree
	IfClause   *ParseTree
	ElseClause *ParseTree
}

func NewIfStatementTree(location SourceRange, condition *ParseTree, ifClause *ParseTree, elseClause *ParseTree) *IfStatementTree {
	return &IfStatementTree{
		ParseTree:  ParseTree{Type: ParseTreeTypeIfStatement, Location: location},
		Condition:  condition,
		IfClause:   ifClause,
		ElseClause: elseClause,
	}
}
