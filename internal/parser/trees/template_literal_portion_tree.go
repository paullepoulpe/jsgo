// Port of com/google/javascript/jscomp/parsing/parser/trees/TemplateLiteralPortionTree.java
// @ closure-compiler v20260415

package trees

// TemplateLiteralPortionTree represents a literal portion of a template literal,
// e.g. "hello" in `hello ${name}`.
type TemplateLiteralPortionTree struct {
	ParseTree
	Value *Token
}

// NewTemplateLiteralPortionTree constructs a TemplateLiteralPortionTree.
func NewTemplateLiteralPortionTree(location SourceRange, value *Token) *TemplateLiteralPortionTree {
	return &TemplateLiteralPortionTree{
		ParseTree: ParseTree{Type: ParseTreeTypeTemplateLiteralPortion, Location: location},
		Value:     value,
	}
}
