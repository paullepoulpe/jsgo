// Port of com/google/javascript/jscomp/parsing/parser/trees/FinallyTree.java
// @ closure-compiler v20260415

package trees

// FinallyTree represents the finally clause in a try statement.
type FinallyTree struct {
	ParseTree
	Block *ParseTree
}

func NewFinallyTree(location SourceRange, block *ParseTree) *FinallyTree {
	return &FinallyTree{
		ParseTree: NewParseTree(ParseTreeTypeFinally, location),
		Block:     block,
	}
}
