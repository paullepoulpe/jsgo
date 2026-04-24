// Port of com/google/javascript/jscomp/parsing/parser/trees/TemplateLiteralPortionTree.java
// @ closure-compiler v20260415

package trees

// TemplateLiteralPortionTree represents a literal portion of a template literal (e.g. `hello`).
type TemplateLiteralPortionTree struct {
	ParseTree
	Value *Token
}

func NewTemplateLiteralPortionTree(location SourceRange, value *Token) *TemplateLiteralPortionTree {
	return &TemplateLiteralPortionTree{
		ParseTree: NewParseTree(ParseTreeTypeTemplateLiteralPortion, location),
		Value:     value,
	}
}
