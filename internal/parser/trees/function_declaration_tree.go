// Port of com/google/javascript/jscomp/parsing/parser/trees/FunctionDeclarationTree.java
// @ closure-compiler v20260415

package trees

// FunctionDeclarationKind mirrors FunctionDeclarationTree.Kind in Java.
type FunctionDeclarationKind int

const (
	FunctionDeclarationKindDeclaration FunctionDeclarationKind = iota
	FunctionDeclarationKindExpression
	FunctionDeclarationKindMember
	FunctionDeclarationKindArrow
)

// FunctionDeclarationTree represents a function declaration, expression, arrow, or member.
type FunctionDeclarationTree struct {
	ParseTree
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

// FunctionDeclarationBuilder mirrors FunctionDeclarationTree.Builder in Java.
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

func (b *FunctionDeclarationBuilder) SetFormalParameterList(fpl *FormalParameterListTree) *FunctionDeclarationBuilder {
	b.formalParameterList = fpl
	return b
}

func (b *FunctionDeclarationBuilder) SetFunctionBody(body *ParseTree) *FunctionDeclarationBuilder {
	b.functionBody = body
	return b
}

func (b *FunctionDeclarationBuilder) SetIsClassMember(isClassMember bool) *FunctionDeclarationBuilder {
	b.isClassMember = isClassMember
	return b
}

func (b *FunctionDeclarationBuilder) SetStatic(isStatic bool) *FunctionDeclarationBuilder {
	b.isStatic = isStatic
	return b
}

func (b *FunctionDeclarationBuilder) SetGenerator(isGenerator bool) *FunctionDeclarationBuilder {
	b.isGenerator = isGenerator
	return b
}

func (b *FunctionDeclarationBuilder) SetOptional(isOptional bool) *FunctionDeclarationBuilder {
	b.isOptional = isOptional
	return b
}

func (b *FunctionDeclarationBuilder) SetAsync(isAsync bool) *FunctionDeclarationBuilder {
	b.isAsync = isAsync
	return b
}

func (b *FunctionDeclarationBuilder) Build(location SourceRange) *FunctionDeclarationTree {
	return &FunctionDeclarationTree{
		ParseTree:           ParseTree{Type: ParseTreeTypeFunctionDeclaration, Location: location},
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
