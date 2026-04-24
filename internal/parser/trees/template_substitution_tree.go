// Port of com/google/javascript/jscomp/parsing/parser/trees/TemplateSubstitutionTree.java
// @ closure-compiler v20260415

package trees

// TemplateSubstitutionTree represents an expression to be evaluated and substituted
// into a template literal, e.g. ${hello}.
type TemplateSubstitutionTree struct {
	ParseTree
	Expression *ParseTree
}

// NewTemplateSubstitutionTree constructs a TemplateSubstitutionTree.
func NewTemplateSubstitutionTree(location SourceRange, expression *ParseTree) *TemplateSubstitutionTree {
	return &TemplateSubstitutionTree{
		ParseTree:  ParseTree{Type: ParseTreeTypeTemplateSubstitution, Location: location},
		Expression: expression,
	}
}
