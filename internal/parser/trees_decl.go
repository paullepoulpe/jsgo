// Port of com/google/javascript/jscomp/parsing/parser/trees/ProgramTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/ImportDeclarationTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/ImportSpecifierTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/ImportMetaExpressionTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/ExportDeclarationTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/ExportSpecifierTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/ClassDeclarationTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/FunctionDeclarationTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/FieldDeclarationTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/ComputedPropertyDefinitionTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/ComputedPropertyFieldTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/ComputedPropertyGetterTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/ComputedPropertyMethodTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/ComputedPropertySetterTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/GetAccessorTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/SetAccessorTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/VariableDeclarationListTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/VariableDeclarationTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/FormalParameterListTree.java @ closure-compiler v20261021

/*
 * Copyright 2011 The Closure Compiler Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

/*
 * Copyright 2019 The Closure Compiler Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

/*
 * Copyright 2014 The Closure Compiler Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package parser

import "github.com/paullepoulpe/jsgo/internal/jutil"

// ProgramTree is the parse tree for ParseTreeType PROGRAM.
type ProgramTree struct {
	ParseTreeBase

	SourceElements []ParseTree
	SourceComments []Comment
}

// NewProgramTree constructs a ProgramTree.
//
// Java: ProgramTree(SourceRange, ImmutableList, ImmutableList)
func NewProgramTree(location SourceRange, sourceElements []ParseTree, sourceComments []Comment) *ProgramTree {
	t := &ProgramTree{}
	t.InitParseTreeBase(ParseTreeType_PROGRAM, location)
	t.SourceElements = sourceElements
	t.SourceComments = sourceComments
	return t
}

// ImportDeclarationTree is the parse tree for ParseTreeType IMPORT_DECLARATION.
type ImportDeclarationTree struct {
	ParseTreeBase

	// The identifier for the "default" import from the module.
	DefaultBindingIdentifier Token
	// The list of names imported from the module when using the
	// `import {...} from 'specifier'` form.
	ImportSpecifierList []ParseTree
	// The namespace into which imported names are put when using the
	// `import * as nameSpaceImportIdentifier from 'specifier'` form.
	NameSpaceImportIdentifier Token
	// The string identifying the module.
	ModuleSpecifier Token
}

// NewImportDeclarationTree constructs a ImportDeclarationTree.
//
// Java: ImportDeclarationTree(SourceRange, IdentifierToken, ImmutableList, IdentifierToken, LiteralToken)
func NewImportDeclarationTree(location SourceRange, defaultBindingIdentifier Token, importSpecifierList []ParseTree, nameSpaceImportIdentifier Token, moduleSpecifier Token) *ImportDeclarationTree {
	t := &ImportDeclarationTree{}
	t.InitParseTreeBase(ParseTreeType_IMPORT_DECLARATION, location)
	t.DefaultBindingIdentifier = defaultBindingIdentifier
	t.ImportSpecifierList = importSpecifierList
	t.NameSpaceImportIdentifier = nameSpaceImportIdentifier
	t.ModuleSpecifier = moduleSpecifier
	return t
}

// ImportSpecifierTree is the parse tree for ParseTreeType IMPORT_SPECIFIER.
type ImportSpecifierTree struct {
	ParseTreeBase

	ImportedName    Token
	DestinationName Token
}

// NewImportSpecifierTree constructs a ImportSpecifierTree.
//
// Java: ImportSpecifierTree(SourceRange, IdentifierToken, IdentifierToken)
func NewImportSpecifierTree(location SourceRange, importedName Token, destinationName Token) *ImportSpecifierTree {
	t := &ImportSpecifierTree{}
	t.InitParseTreeBase(ParseTreeType_IMPORT_SPECIFIER, location)
	t.ImportedName = importedName
	t.DestinationName = destinationName
	return t
}

// Represents 'new.target'.
type ImportMetaExpressionTree struct {
	ParseTreeBase
}

// NewImportMetaExpressionTree constructs a ImportMetaExpressionTree.
//
// Java: ImportMetaExpressionTree(SourceRange)
func NewImportMetaExpressionTree(location SourceRange) *ImportMetaExpressionTree {
	t := &ImportMetaExpressionTree{}
	t.InitParseTreeBase(ParseTreeType_IMPORT_META_EXPRESSION, location)
	return t
}

// ExportDeclarationTree is the parse tree for ParseTreeType EXPORT_DECLARATION.
type ExportDeclarationTree struct {
	ParseTreeBase

	// VariableStatement, Declaration (let, const, function or class),
	// or default with assignmentExpression
	IsDefault           bool
	IsExportAll         bool
	Declaration         ParseTree
	ExportSpecifierList []ParseTree
	From                Token
}

// NewExportDeclarationTree constructs a ExportDeclarationTree.
//
// Java: ExportDeclarationTree(SourceRange, boolean, boolean, ParseTree, ImmutableList, LiteralToken)
func NewExportDeclarationTree(location SourceRange, isDefault bool, isExportAll bool, declaration ParseTree, exportSpecifierList []ParseTree, from Token) *ExportDeclarationTree {
	t := &ExportDeclarationTree{}
	t.InitParseTreeBase(ParseTreeType_EXPORT_DECLARATION, location)
	t.IsDefault = isDefault
	t.IsExportAll = isExportAll
	t.Declaration = declaration
	t.ExportSpecifierList = exportSpecifierList
	t.From = from
	return t
}

// ExportSpecifierTree is the parse tree for ParseTreeType EXPORT_SPECIFIER.
type ExportSpecifierTree struct {
	ParseTreeBase

	ImportedName    Token
	DestinationName Token
}

// NewExportSpecifierTree constructs a ExportSpecifierTree.
//
// Java: ExportSpecifierTree(SourceRange, IdentifierToken, IdentifierToken)
func NewExportSpecifierTree(location SourceRange, importedName Token, destinationName Token) *ExportSpecifierTree {
	t := &ExportSpecifierTree{}
	t.InitParseTreeBase(ParseTreeType_EXPORT_SPECIFIER, location)
	t.ImportedName = importedName
	t.DestinationName = destinationName
	return t
}

// ClassDeclarationTree is the parse tree for ParseTreeType CLASS_DECLARATION.
type ClassDeclarationTree struct {
	ParseTreeBase

	Name       Token
	SuperClass ParseTree
	Elements   []ParseTree
}

// NewClassDeclarationTree constructs a ClassDeclarationTree.
//
// Java: ClassDeclarationTree(SourceRange, IdentifierToken, ParseTree, ImmutableList)
func NewClassDeclarationTree(location SourceRange, name Token, superClass ParseTree, elements []ParseTree) *ClassDeclarationTree {
	t := &ClassDeclarationTree{}
	t.InitParseTreeBase(ParseTreeType_CLASS_DECLARATION, location)
	t.Name = name
	t.SuperClass = superClass
	t.Elements = elements
	return t
}

// FunctionDeclarationTree is the parse tree for ParseTreeType
// FUNCTION_DECLARATION (declarations, expressions, methods and arrows).
type FunctionDeclarationTree struct {
	ParseTreeBase

	Name                Token
	FormalParameterList *FormalParameterListTree
	FunctionBody        ParseTree
	IsClassMember       bool
	IsStatic            bool
	IsGenerator         bool
	IsOptional          bool
	IsAsync             bool
	Kind                FunctionDeclarationTree_Kind
}

// FunctionDeclarationTree_Kind is FunctionDeclarationTree.Kind.
type FunctionDeclarationTree_Kind int32

const (
	FunctionDeclarationTree_Kind_DECLARATION FunctionDeclarationTree_Kind = iota
	FunctionDeclarationTree_Kind_EXPRESSION
	FunctionDeclarationTree_Kind_MEMBER
	FunctionDeclarationTree_Kind_ARROW
)

var functionDeclarationTree_Kind_names = [...]string{"DECLARATION", "EXPRESSION", "MEMBER", "ARROW"}

// FunctionDeclarationTree_Kind_values is Kind.values().
//
// Java: FunctionDeclarationTree.Kind.values()
func FunctionDeclarationTree_Kind_values() []FunctionDeclarationTree_Kind {
	return []FunctionDeclarationTree_Kind{
		FunctionDeclarationTree_Kind_DECLARATION,
		FunctionDeclarationTree_Kind_EXPRESSION,
		FunctionDeclarationTree_Kind_MEMBER,
		FunctionDeclarationTree_Kind_ARROW,
	}
}

// Name is Enum.name().
//
// Java: FunctionDeclarationTree.Kind.name()
func (k FunctionDeclarationTree_Kind) Name() string { return functionDeclarationTree_Kind_names[k] }

// String is Enum.toString().
//
// Java: FunctionDeclarationTree.Kind.toString()
func (k FunctionDeclarationTree_Kind) String() string { return k.Name() }

// FunctionDeclarationTree_NewBuilder returns a builder for a function of the
// given kind.
//
// Go-only: language-adapter — the mechanical name FunctionDeclarationTree_Builder
// is taken by the nested Builder type.
//
// Java: FunctionDeclarationTree.builder(Kind)
func FunctionDeclarationTree_NewBuilder(kind FunctionDeclarationTree_Kind) *FunctionDeclarationTree_Builder {
	return NewFunctionDeclarationTree_Builder(kind)
}

// NewFunctionDeclarationTree constructs a FunctionDeclarationTree from a
// builder.
//
// Java: FunctionDeclarationTree(Builder)
func NewFunctionDeclarationTree(builder *FunctionDeclarationTree_Builder) *FunctionDeclarationTree {
	t := &FunctionDeclarationTree{}
	t.InitParseTreeBase(ParseTreeType_FUNCTION_DECLARATION, builder.location)

	t.Name = builder.name
	t.IsClassMember = builder.isClassMember
	t.IsStatic = builder.isStatic
	t.IsGenerator = builder.isGenerator
	t.IsOptional = builder.isOptional
	t.Kind = builder.kind // Java: checkNotNull(builder.kind); an enum value is never null in Go.
	if builder.formalParameterList == nil || builder.functionBody == nil {
		panic(jutil.NewNullPointerException())
	}
	t.FormalParameterList = builder.formalParameterList
	t.FunctionBody = builder.functionBody
	t.IsAsync = builder.isAsync
	return t
}

// FunctionDeclarationTree_Builder builds a FunctionDeclarationTree.
type FunctionDeclarationTree_Builder struct {
	kind FunctionDeclarationTree_Kind

	name                Token
	formalParameterList *FormalParameterListTree
	functionBody        ParseTree
	isClassMember       bool
	isStatic            bool
	isGenerator         bool
	isOptional          bool
	isAsync             bool
	location            SourceRange
}

// NewFunctionDeclarationTree_Builder constructs a Builder.
//
// Java: FunctionDeclarationTree.Builder(Kind)
func NewFunctionDeclarationTree_Builder(kind FunctionDeclarationTree_Kind) *FunctionDeclarationTree_Builder {
	b := &FunctionDeclarationTree_Builder{}
	b.kind = kind
	return b
}

// SetName sets the optional function name.
//
// Default is nil.
//
// Java: FunctionDeclarationTree.Builder.setName(IdentifierToken)
func (b *FunctionDeclarationTree_Builder) SetName(name Token) *FunctionDeclarationTree_Builder {
	b.name = name
	return b
}

// SetFormalParameterList sets the required parameter list.
//
// Java: FunctionDeclarationTree.Builder.setFormalParameterList(FormalParameterListTree)
func (b *FunctionDeclarationTree_Builder) SetFormalParameterList(formalParameterList *FormalParameterListTree) *FunctionDeclarationTree_Builder {
	b.formalParameterList = formalParameterList
	return b
}

// SetFunctionBody sets the required function body.
//
// Java: FunctionDeclarationTree.Builder.setFunctionBody(ParseTree)
func (b *FunctionDeclarationTree_Builder) SetFunctionBody(functionBody ParseTree) *FunctionDeclarationTree_Builder {
	b.functionBody = functionBody
	return b
}

// SetIsClassMember: is the method an ES6 class member?
//
// Default is false. Only relevant for class method member declarations.
//
// Java: FunctionDeclarationTree.Builder.setIsClassMember(boolean)
func (b *FunctionDeclarationTree_Builder) SetIsClassMember(isClassMember bool) *FunctionDeclarationTree_Builder {
	b.isClassMember = isClassMember
	return b
}

// SetStatic: is the method static?
//
// Default is false. Only relevant for class method member declarations.
//
// Java: FunctionDeclarationTree.Builder.setStatic(boolean)
func (b *FunctionDeclarationTree_Builder) SetStatic(isStatic bool) *FunctionDeclarationTree_Builder {
	b.isStatic = isStatic
	return b
}

// SetGenerator: is this a generator function?
//
// Default is false.
//
// Java: FunctionDeclarationTree.Builder.setGenerator(boolean)
func (b *FunctionDeclarationTree_Builder) SetGenerator(isGenerator bool) *FunctionDeclarationTree_Builder {
	b.isGenerator = isGenerator
	return b
}

// SetOptional: is this the declaration of an optional function parameter?
// Default is false.
//
// Only relevant for function declaration as a parameter to another
// function.
//
// Java: FunctionDeclarationTree.Builder.setOptional(boolean)
func (b *FunctionDeclarationTree_Builder) SetOptional(isOptional bool) *FunctionDeclarationTree_Builder {
	b.isOptional = isOptional
	return b
}

// SetAsync: is this an asynchronous function?
//
// Default is false.
//
// Java: FunctionDeclarationTree.Builder.setAsync(boolean)
func (b *FunctionDeclarationTree_Builder) SetAsync(isAsync bool) *FunctionDeclarationTree_Builder {
	b.isAsync = isAsync
	return b
}

// Build returns a new FunctionDeclarationTree.
//
// The location is provided at this point because it cannot be correctly
// calculated until the whole function has been parsed.
//
// Java: FunctionDeclarationTree.Builder.build(SourceRange)
func (b *FunctionDeclarationTree_Builder) Build(location SourceRange) *FunctionDeclarationTree {
	b.location = location
	return NewFunctionDeclarationTree(b)
}

// Tree for public property fields
type FieldDeclarationTree struct {
	ParseTreeBase

	Name        Token
	IsStatic    bool
	Initializer ParseTree
}

// NewFieldDeclarationTree constructs a FieldDeclarationTree.
//
// Java: FieldDeclarationTree(SourceRange, IdentifierToken, boolean, ParseTree)
func NewFieldDeclarationTree(location SourceRange, name Token, isStatic bool, initializer ParseTree) *FieldDeclarationTree {
	t := &FieldDeclarationTree{}
	t.InitParseTreeBase(ParseTreeType_FIELD_DECLARATION, location)
	t.Name = name
	t.IsStatic = isStatic
	t.Initializer = initializer
	return t
}

// Represents a computed property in an object literal.
type ComputedPropertyDefinitionTree struct {
	ParseTreeBase

	Property ParseTree
	Value    ParseTree
}

// NewComputedPropertyDefinitionTree constructs a ComputedPropertyDefinitionTree.
//
// Java: ComputedPropertyDefinitionTree(SourceRange, ParseTree, ParseTree)
func NewComputedPropertyDefinitionTree(location SourceRange, property ParseTree, value ParseTree) *ComputedPropertyDefinitionTree {
	t := &ComputedPropertyDefinitionTree{}
	t.InitParseTreeBase(ParseTreeType_COMPUTED_PROPERTY_DEFINITION, location)
	t.Property = property
	t.Value = value
	return t
}

// Tree for public comptuted property fields
type ComputedPropertyFieldTree struct {
	ParseTreeBase
	// Go-only: performance — before the other fields, in the base's padding.
	IsStatic bool

	Property    ParseTree
	Initializer ParseTree
}

// NewComputedPropertyFieldTree constructs a ComputedPropertyFieldTree.
//
// Java: ComputedPropertyFieldTree(SourceRange, ParseTree, boolean, ParseTree)
func NewComputedPropertyFieldTree(location SourceRange, property ParseTree, isStatic bool, initializer ParseTree) *ComputedPropertyFieldTree {
	t := &ComputedPropertyFieldTree{}
	t.InitParseTreeBase(ParseTreeType_COMPUTED_PROPERTY_FIELD, location)
	t.Property = property
	t.IsStatic = isStatic
	t.Initializer = initializer
	return t
}

// ComputedPropertyGetterTree is the parse tree for ParseTreeType COMPUTED_PROPERTY_GETTER.
type ComputedPropertyGetterTree struct {
	ParseTreeBase

	Property ParseTree
	IsStatic bool
	Body     *BlockTree
}

// NewComputedPropertyGetterTree constructs a ComputedPropertyGetterTree.
//
// Java: ComputedPropertyGetterTree(SourceRange, ParseTree, boolean, BlockTree)
func NewComputedPropertyGetterTree(location SourceRange, property ParseTree, isStatic bool, body *BlockTree) *ComputedPropertyGetterTree {
	t := &ComputedPropertyGetterTree{}
	t.InitParseTreeBase(ParseTreeType_COMPUTED_PROPERTY_GETTER, location)
	t.Property = property
	t.IsStatic = isStatic
	t.Body = body
	return t
}

// ComputedPropertyMethodTree is the parse tree for ParseTreeType COMPUTED_PROPERTY_METHOD.
type ComputedPropertyMethodTree struct {
	ParseTreeBase

	Property ParseTree
	Method   ParseTree
}

// NewComputedPropertyMethodTree constructs a ComputedPropertyMethodTree.
//
// Java: ComputedPropertyMethodTree(SourceRange, ParseTree, ParseTree)
func NewComputedPropertyMethodTree(location SourceRange, property ParseTree, method ParseTree) *ComputedPropertyMethodTree {
	t := &ComputedPropertyMethodTree{}
	t.InitParseTreeBase(ParseTreeType_COMPUTED_PROPERTY_METHOD, location)
	t.Property = property
	t.Method = method
	return t
}

// ComputedPropertySetterTree is the parse tree for ParseTreeType COMPUTED_PROPERTY_SETTER.
type ComputedPropertySetterTree struct {
	ParseTreeBase
	// Go-only: performance — before the other fields, in the base's padding.
	IsStatic bool

	Property  ParseTree
	Parameter *FormalParameterListTree
	Body      *BlockTree
}

// NewComputedPropertySetterTree constructs a ComputedPropertySetterTree.
//
// Java: ComputedPropertySetterTree(SourceRange, ParseTree, boolean, FormalParameterListTree, BlockTree)
func NewComputedPropertySetterTree(location SourceRange, property ParseTree, isStatic bool, parameter *FormalParameterListTree, body *BlockTree) *ComputedPropertySetterTree {
	t := &ComputedPropertySetterTree{}
	t.InitParseTreeBase(ParseTreeType_COMPUTED_PROPERTY_SETTER, location)
	t.Property = property
	t.IsStatic = isStatic
	t.Parameter = parameter
	t.Body = body
	return t
}

// GetAccessorTree is the parse tree for ParseTreeType GET_ACCESSOR.
type GetAccessorTree struct {
	ParseTreeBase

	PropertyName  Token
	IsClassMember bool
	IsStatic      bool
	Body          *BlockTree
}

// NewGetAccessorTree constructs a GetAccessorTree.
//
// Java: GetAccessorTree(SourceRange, Token, boolean, boolean, BlockTree)
func NewGetAccessorTree(location SourceRange, propertyName Token, isClassMember bool, isStatic bool, body *BlockTree) *GetAccessorTree {
	t := &GetAccessorTree{}
	t.InitParseTreeBase(ParseTreeType_GET_ACCESSOR, location)
	t.PropertyName = propertyName
	t.IsClassMember = isClassMember
	t.IsStatic = isStatic
	t.Body = body
	return t
}

// SetAccessorTree is the parse tree for ParseTreeType SET_ACCESSOR.
type SetAccessorTree struct {
	ParseTreeBase
	// Go-only: performance — before the other fields, in the base's padding.
	IsClassMember bool
	IsStatic      bool

	PropertyName Token
	Parameter    *FormalParameterListTree
	Body         *BlockTree
}

// NewSetAccessorTree constructs a SetAccessorTree.
//
// Java: SetAccessorTree(SourceRange, Token, boolean, boolean, FormalParameterListTree, BlockTree)
func NewSetAccessorTree(location SourceRange, propertyName Token, isClassMember bool, isStatic bool, parameter *FormalParameterListTree, body *BlockTree) *SetAccessorTree {
	t := &SetAccessorTree{}
	t.InitParseTreeBase(ParseTreeType_SET_ACCESSOR, location)
	t.PropertyName = propertyName
	t.IsClassMember = isClassMember
	t.IsStatic = isStatic
	t.Parameter = parameter
	t.Body = body
	return t
}

// VariableDeclarationListTree is the parse tree for ParseTreeType VARIABLE_DECLARATION_LIST.
type VariableDeclarationListTree struct {
	ParseTreeBase

	DeclarationType TokenType // const, var, let
	Declarations    []*VariableDeclarationTree
}

// NewVariableDeclarationListTree constructs a VariableDeclarationListTree.
//
// Java: VariableDeclarationListTree(SourceRange, TokenType, ImmutableList)
func NewVariableDeclarationListTree(location SourceRange, declarationType TokenType, declarations []*VariableDeclarationTree) *VariableDeclarationListTree {
	t := &VariableDeclarationListTree{}
	t.InitParseTreeBase(ParseTreeType_VARIABLE_DECLARATION_LIST, location)
	t.DeclarationType = declarationType
	t.Declarations = declarations
	return t
}

// VariableDeclarationTree is the parse tree for ParseTreeType VARIABLE_DECLARATION.
type VariableDeclarationTree struct {
	ParseTreeBase

	// An IdentifierExpressionTree, ArrayPatternTree or ObjectLiteralTree
	Lvalue      ParseTree
	Initializer ParseTree
}

// NewVariableDeclarationTree constructs a VariableDeclarationTree.
//
// Java: VariableDeclarationTree(SourceRange, ParseTree, ParseTree)
func NewVariableDeclarationTree(location SourceRange, lvalue ParseTree, initializer ParseTree) *VariableDeclarationTree {
	t := &VariableDeclarationTree{}
	t.InitParseTreeBase(ParseTreeType_VARIABLE_DECLARATION, location)
	t.Lvalue = lvalue
	t.Initializer = initializer
	return t
}

// FormalParameterListTree is the parse tree for ParseTreeType FORMAL_PARAMETER_LIST.
type FormalParameterListTree struct {
	ParseTreeBase
	// Go-only: performance — before the other fields, in the base's padding.
	HasTrailingComma bool

	Parameters     []ParseTree
	CommaPositions []SourcePosition
}

// NewFormalParameterListTree constructs a FormalParameterListTree.
//
// Java: FormalParameterListTree(SourceRange, ImmutableList, boolean, ImmutableList)
func NewFormalParameterListTree(location SourceRange, parameters []ParseTree, hasTrailingComma bool, commaPositions []SourcePosition) *FormalParameterListTree {
	t := &FormalParameterListTree{}
	t.InitParseTreeBase(ParseTreeType_FORMAL_PARAMETER_LIST, location)
	numParams := len(parameters)
	numCommas := len(commaPositions)
	if !(numCommas <= numParams && numCommas >= numParams-1) {
		// Java: checkArgument(cond, "...", parameters, commaPositions); the lists
		// are rendered with AbstractCollection.toString only on failure.
		panic(jutil.NewIllegalArgumentExceptionMsg(jutil.LenientFormat(
			"Unexpected # of comma and formal params.\nparams: %s\ncomma positions: %s",
			jutil.CollectionToString(parameters),
			jutil.CollectionToString(commaPositions))))
	}
	t.Parameters = parameters
	t.HasTrailingComma = hasTrailingComma
	t.CommaPositions = commaPositions
	return t
}
