// Port of com/google/javascript/jscomp/parsing/parser/trees/IfStatementTree.java
// @ closure-compiler v20260415

package trees

// IfStatementTree represents an if statement with an optional else clause.
type IfStatementTree struct {
	ParseTree
	Condition  *ParseTree
	IfClause   *ParseTree
	// ElseClause is nil when there is no else branch.
	ElseClause *ParseTree
}

func NewIfStatementTree(
	location SourceRange,
	condition *ParseTree,
	ifClause *ParseTree,
	elseClause *ParseTree,
) *IfStatementTree {
	return &IfStatementTree{
		ParseTree:  NewParseTree(ParseTreeTypeIfStatement, location),
		Condition:  condition,
		IfClause:   ifClause,
		ElseClause: elseClause,
	}
}
