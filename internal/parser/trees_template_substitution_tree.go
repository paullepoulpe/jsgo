// Port of com/google/javascript/jscomp/parsing/parser/trees/TemplateSubstitutionTree.java @ closure-compiler v20260415

package parser

// A production representing the expression to be evaluated and substituted
// into a template literal. Eg: ${hello}.
type TemplateSubstitutionTree struct {
	ParseTreeBase

	Expression ParseTree
}

// NewTemplateSubstitutionTree constructs a TemplateSubstitutionTree.
//
// Java: TemplateSubstitutionTree(SourceRange, ParseTree)
func NewTemplateSubstitutionTree(location *SourceRange, expression ParseTree) *TemplateSubstitutionTree {
	t := &TemplateSubstitutionTree{}
	t.InitParseTreeBase(t, ParseTreeType_TEMPLATE_SUBSTITUTION, location)
	t.Expression = expression
	return t
}
