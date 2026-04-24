// Port of com/google/javascript/jscomp/parsing/parser/trees/FunctionDeclarationTree.java
// @ closure-compiler v20260415

package trees

// FunctionDeclarationKind identifies how a function is declared.
type FunctionDeclarationKind int

const (
	FunctionDeclarationKindDeclaration FunctionDeclarationKind = iota
	FunctionDeclarationKindExpression
	FunctionDeclarationKindMember
	FunctionDeclarationKindArrow
)

// FunctionDeclarationTree represents a function declaration, expression, method, or arrow.
type FunctionDeclarationTree struct {
	ParseTree
	// Name is nil for anonymous function expressions.
	Name                *IdentifierToken
	FormalParameterList *FormalParameterListTree
	FunctionBody        *ParseTree
	IsClassMember       bool
	IsStatic            bool
	IsGenerator         bool
	IsOptional          bool
	IsAsync             bool
	Kind                FunctionDeclarationKind
}

// FunctionDeclarationBuilder mirrors the Java builder pattern.
type FunctionDeclarationBuilder struct {
	kind                FunctionDeclarationKind
	name                *IdentifierToken
	formalParameterList *FormalParameterListTree
	functionBody        *ParseTree
	isClassMember       bool
	isStatic            bool
	isGenerator         bool
	isOptional          bool
	isAsync             bool
}

func NewFunctionDeclarationBuilder(kind FunctionDeclarationKind) *FunctionDeclarationBuilder {
	return &FunctionDeclarationBuilder{kind: kind}
}

func (b *FunctionDeclarationBuilder) SetName(name *IdentifierToken) *FunctionDeclarationBuilder {
	b.name = name
	return b
}

func (b *FunctionDeclarationBuilder) SetFormalParameterList(p *FormalParameterListTree) *FunctionDeclarationBuilder {
	b.formalParameterList = p
	return b
}

func (b *FunctionDeclarationBuilder) SetFunctionBody(body *ParseTree) *FunctionDeclarationBuilder {
	b.functionBody = body
	return b
}

func (b *FunctionDeclarationBuilder) SetIsClassMember(v bool) *FunctionDeclarationBuilder {
	b.isClassMember = v
	return b
}

func (b *FunctionDeclarationBuilder) SetStatic(v bool) *FunctionDeclarationBuilder {
	b.isStatic = v
	return b
}

func (b *FunctionDeclarationBuilder) SetGenerator(v bool) *FunctionDeclarationBuilder {
	b.isGenerator = v
	return b
}

func (b *FunctionDeclarationBuilder) SetOptional(v bool) *FunctionDeclarationBuilder {
	b.isOptional = v
	return b
}

func (b *FunctionDeclarationBuilder) SetAsync(v bool) *FunctionDeclarationBuilder {
	b.isAsync = v
	return b
}

func (b *FunctionDeclarationBuilder) Build(location SourceRange) *FunctionDeclarationTree {
	return &FunctionDeclarationTree{
		ParseTree:           NewParseTree(ParseTreeTypeFunctionDeclaration, location),
		Name:                b.name,
		FormalParameterList: b.formalParameterList,
		FunctionBody:        b.functionBody,
		IsClassMember:       b.isClassMember,
		IsStatic:            b.isStatic,
		IsGenerator:         b.isGenerator,
		IsOptional:          b.isOptional,
		IsAsync:             b.isAsync,
		Kind:                b.kind,
	}
}
