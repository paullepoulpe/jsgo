// Port of com/google/javascript/jscomp/parsing/parser/trees/TemplateSubstitutionTree.java
// @ closure-compiler v20260415

package trees

// TemplateSubstitutionTree represents the expression evaluated within a template literal
// substitution (e.g. ${hello}).
type TemplateSubstitutionTree struct {
	ParseTree
	Expression *ParseTree
}

func NewTemplateSubstitutionTree(location SourceRange, expression *ParseTree) *TemplateSubstitutionTree {
	return &TemplateSubstitutionTree{
		ParseTree:  NewParseTree(ParseTreeTypeTemplateSubstitution, location),
		Expression: expression,
	}
}
