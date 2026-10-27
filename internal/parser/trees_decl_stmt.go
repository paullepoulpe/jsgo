// Port of com/google/javascript/jscomp/parsing/parser/trees/ParseTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/ParseTreeType.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/NullTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/MissingPrimaryExpressionTree.java @ closure-compiler v20261021
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
// Port of com/google/javascript/jscomp/parsing/parser/trees/BlockTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/IfStatementTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/ForAwaitOfStatementTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/ForInStatementTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/ForOfStatementTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/ForStatementTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/WhileStatementTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/DoWhileStatementTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/SwitchStatementTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/CaseClauseTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/DefaultClauseTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/TryStatementTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/CatchTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/FinallyTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/ReturnStatementTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/BreakStatementTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/ContinueStatementTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/ThrowStatementTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/LabelledStatementTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/WithStatementTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/VariableStatementTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/VariableDeclarationListTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/VariableDeclarationTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/ExpressionStatementTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/EmptyStatementTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/DebuggerStatementTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/FormalParameterListTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/Comment.java @ closure-compiler v20261021

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

// ParseTree is an abstract syntax tree for JavaScript parse trees.
// Immutable.
// A plain old data structure. Should include data members and simple
// accessors only.
//
// Derived classes should have a 'Tree' suffix. Each concrete derived class
// should have a ParseTreeType whose name matches the derived class name.
//
// A parse tree derived from source should have a non-null location. A parse
// tree that is synthesized by the compiler may have a null location.
//
// When adding a new subclass of ParseTree you must also do the following:
//   - add a new entry to ParseTreeType
//   - add ParseTree.asXTree()
//
// Java: class ParseTree. Every tree embeds ParseTreeBase; fields are read
// through GetParseTreeBase(), and the asX() casts are Go type assertions
// (tree.(*XTree)).
type ParseTree interface {
	GetParseTreeBase() *ParseTreeBase
	GetStart() SourcePosition
	GetEnd() SourcePosition
	IsPattern() bool
	IsValidNonVanillaAssignmentTarget() bool
	IsValidAssignmentTarget() bool
	IsRestParameter() bool
	String() string
}

// ParseTreeBase holds ParseTree's fields and methods.
type ParseTreeBase struct {
	Type     ParseTreeType
	Location SourceRange
}

// InitParseTreeBase is the Go form of the (protected) ParseTree
// constructor.
//
// Java: ParseTree(ParseTreeType, SourceRange)
func (t *ParseTreeBase) InitParseTreeBase(type_ ParseTreeType, location SourceRange) {
	t.Type = type_
	t.Location = location
}

// GetParseTreeBase returns the base struct (field access through the
// interface).
//
// Go-only: language-adapter — field access through the interface (PORTING.md §4).
func (t *ParseTreeBase) GetParseTreeBase() *ParseTreeBase {
	return t
}

// GetStart returns the start position.
//
// Java: ParseTree.getStart()
func (t *ParseTreeBase) GetStart() SourcePosition {
	return t.Location.Start
}

// GetEnd returns the end position.
//
// Java: ParseTree.getEnd()
func (t *ParseTreeBase) GetEnd() SourcePosition {
	return t.Location.End
}

// IsPattern reports whether the tree (inside any parentheses) is an array or
// object pattern.
//
// Java: ParseTree.isPattern() — Go-only: language-adapter — the loop over
// parentheses is ParenExpressionTree.IsPattern.
func (t *ParseTreeBase) IsPattern() bool {
	switch t.Type {
	case ParseTreeType_ARRAY_PATTERN, ParseTreeType_OBJECT_PATTERN:
		return true
	default:
		return false
	}
}

// IsValidNonVanillaAssignmentTarget: is valid assignment target for
// non-vanilla assignment operators like `+=`, `-+`, `**=`, etc
//
// Java: ParseTree.isValidNonVanillaAssignmentTarget() — Go-only:
// language-adapter — the loop over parentheses is
// ParenExpressionTree.IsValidNonVanillaAssignmentTarget.
func (t *ParseTreeBase) IsValidNonVanillaAssignmentTarget() bool {
	switch t.Type {
	case ParseTreeType_IDENTIFIER_EXPRESSION, ParseTreeType_MEMBER_EXPRESSION, ParseTreeType_MEMBER_LOOKUP_EXPRESSION, ParseTreeType_DEFAULT_PARAMETER:
		return true
	case ParseTreeType_ARRAY_PATTERN, ParseTreeType_OBJECT_PATTERN:
		return false
	default:
		return false
	}
}

// IsValidAssignmentTarget: is valid assignment target for any assignment
// operator like `=`, `+=`, `-+`, `**=`, etc
//
// Java: ParseTree.isValidAssignmentTarget() — Go-only: language-adapter —
// the loop over parentheses is ParenExpressionTree.IsValidAssignmentTarget.
func (t *ParseTreeBase) IsValidAssignmentTarget() bool {
	switch t.Type {
	case ParseTreeType_IDENTIFIER_EXPRESSION,
		ParseTreeType_MEMBER_EXPRESSION,
		ParseTreeType_MEMBER_LOOKUP_EXPRESSION,
		ParseTreeType_ARRAY_PATTERN,
		ParseTreeType_OBJECT_PATTERN,
		ParseTreeType_DEFAULT_PARAMETER:
		return true
	default:
		return false
	}
}

// IsRestParameter reports whether this is an ITER_REST tree.
//
// Java: ParseTree.isRestParameter()
func (t *ParseTreeBase) IsRestParameter() bool {
	return t.Type == ParseTreeType_ITER_REST
}

// String renders "<type>@<location>".
//
// Java: ParseTree.toString()
func (t *ParseTreeBase) String() string {
	return t.Type.String() + "@" + jutil.ToJavaString(t.Location)
}

// ParseTreeType: the types of concrete parse trees.
//
// The name of the ParseTreeType must match the name of the class that it
// applies to. For example the DerivedTree class should use
// ParseTreeType.DERIVED.
type ParseTreeType int32

const (
	ParseTreeType_PROGRAM ParseTreeType = iota
	ParseTreeType_FUNCTION_DECLARATION
	ParseTreeType_BLOCK
	ParseTreeType_VARIABLE_STATEMENT
	ParseTreeType_VARIABLE_DECLARATION
	ParseTreeType_EMPTY_STATEMENT
	ParseTreeType_EXPRESSION_STATEMENT
	ParseTreeType_IF_STATEMENT
	ParseTreeType_DO_WHILE_STATEMENT
	ParseTreeType_WHILE_STATEMENT
	ParseTreeType_FOR_IN_STATEMENT
	ParseTreeType_FOR_STATEMENT
	ParseTreeType_VARIABLE_DECLARATION_LIST
	ParseTreeType_CONTINUE_STATEMENT
	ParseTreeType_BREAK_STATEMENT
	ParseTreeType_RETURN_STATEMENT
	ParseTreeType_WITH_STATEMENT
	ParseTreeType_CASE_CLAUSE
	ParseTreeType_DEFAULT_CLAUSE
	ParseTreeType_SWITCH_STATEMENT
	ParseTreeType_LABELLED_STATEMENT
	ParseTreeType_THROW_STATEMENT
	ParseTreeType_CATCH
	ParseTreeType_TRY_STATEMENT
	ParseTreeType_DEBUGGER_STATEMENT
	ParseTreeType_THIS_EXPRESSION
	ParseTreeType_IDENTIFIER_EXPRESSION
	ParseTreeType_LITERAL_EXPRESSION
	ParseTreeType_ARRAY_LITERAL_EXPRESSION
	ParseTreeType_OBJECT_LITERAL_EXPRESSION
	ParseTreeType_GET_ACCESSOR
	ParseTreeType_SET_ACCESSOR
	ParseTreeType_PROPERTY_NAME_ASSIGNMENT
	ParseTreeType_COMPUTED_PROPERTY_DEFINITION
	ParseTreeType_COMPUTED_PROPERTY_GETTER
	ParseTreeType_COMPUTED_PROPERTY_METHOD
	ParseTreeType_COMPUTED_PROPERTY_SETTER
	ParseTreeType_MISSING_PRIMARY_EXPRESSION
	ParseTreeType_COMMA_EXPRESSION
	ParseTreeType_BINARY_OPERATOR
	ParseTreeType_CONDITIONAL_EXPRESSION
	ParseTreeType_UNARY_EXPRESSION
	ParseTreeType_UPDATE_EXPRESSION
	ParseTreeType_MEMBER_EXPRESSION
	ParseTreeType_OPT_CHAIN_MEMBER_EXPRESSION
	ParseTreeType_NEW_EXPRESSION
	ParseTreeType_ARGUMENT_LIST
	ParseTreeType_CALL_EXPRESSION
	ParseTreeType_OPT_CHAIN_CALL_EXPRESSION
	ParseTreeType_CLASS_DECLARATION
	ParseTreeType_MEMBER_LOOKUP_EXPRESSION
	ParseTreeType_OPT_CHAIN_MEMBER_LOOKUP_EXPRESSION
	ParseTreeType_PAREN_EXPRESSION
	ParseTreeType_FINALLY
	ParseTreeType_SUPER_EXPRESSION
	ParseTreeType_ARRAY_PATTERN
	ParseTreeType_OBJECT_PATTERN
	ParseTreeType_FORMAL_PARAMETER_LIST
	ParseTreeType_ITER_SPREAD
	ParseTreeType_OBJECT_SPREAD
	ParseTreeType_NULL
	ParseTreeType_ITER_REST
	ParseTreeType_OBJECT_REST
	ParseTreeType_MODULE_IMPORT
	ParseTreeType_EXPORT_DECLARATION
	ParseTreeType_EXPORT_SPECIFIER
	ParseTreeType_IMPORT_DECLARATION
	ParseTreeType_IMPORT_SPECIFIER
	ParseTreeType_FOR_OF_STATEMENT
	ParseTreeType_FOR_AWAIT_OF_STATEMENT
	ParseTreeType_YIELD_EXPRESSION
	ParseTreeType_DEFAULT_PARAMETER
	ParseTreeType_TEMPLATE_LITERAL_EXPRESSION
	ParseTreeType_TEMPLATE_LITERAL_PORTION
	ParseTreeType_TEMPLATE_SUBSTITUTION
	ParseTreeType_OPTIONAL_PARAMETER
	ParseTreeType_NEW_TARGET_EXPRESSION
	ParseTreeType_AWAIT_EXPRESSION
	ParseTreeType_DYNAMIC_IMPORT_EXPRESSION
	ParseTreeType_IMPORT_META_EXPRESSION
	ParseTreeType_COMPUTED_PROPERTY_FIELD
	ParseTreeType_FIELD_DECLARATION
)

var parseTreeType_names = [...]string{
	"PROGRAM",
	"FUNCTION_DECLARATION",
	"BLOCK",
	"VARIABLE_STATEMENT",
	"VARIABLE_DECLARATION",
	"EMPTY_STATEMENT",
	"EXPRESSION_STATEMENT",
	"IF_STATEMENT",
	"DO_WHILE_STATEMENT",
	"WHILE_STATEMENT",
	"FOR_IN_STATEMENT",
	"FOR_STATEMENT",
	"VARIABLE_DECLARATION_LIST",
	"CONTINUE_STATEMENT",
	"BREAK_STATEMENT",
	"RETURN_STATEMENT",
	"WITH_STATEMENT",
	"CASE_CLAUSE",
	"DEFAULT_CLAUSE",
	"SWITCH_STATEMENT",
	"LABELLED_STATEMENT",
	"THROW_STATEMENT",
	"CATCH",
	"TRY_STATEMENT",
	"DEBUGGER_STATEMENT",
	"THIS_EXPRESSION",
	"IDENTIFIER_EXPRESSION",
	"LITERAL_EXPRESSION",
	"ARRAY_LITERAL_EXPRESSION",
	"OBJECT_LITERAL_EXPRESSION",
	"GET_ACCESSOR",
	"SET_ACCESSOR",
	"PROPERTY_NAME_ASSIGNMENT",
	"COMPUTED_PROPERTY_DEFINITION",
	"COMPUTED_PROPERTY_GETTER",
	"COMPUTED_PROPERTY_METHOD",
	"COMPUTED_PROPERTY_SETTER",
	"MISSING_PRIMARY_EXPRESSION",
	"COMMA_EXPRESSION",
	"BINARY_OPERATOR",
	"CONDITIONAL_EXPRESSION",
	"UNARY_EXPRESSION",
	"UPDATE_EXPRESSION",
	"MEMBER_EXPRESSION",
	"OPT_CHAIN_MEMBER_EXPRESSION",
	"NEW_EXPRESSION",
	"ARGUMENT_LIST",
	"CALL_EXPRESSION",
	"OPT_CHAIN_CALL_EXPRESSION",
	"CLASS_DECLARATION",
	"MEMBER_LOOKUP_EXPRESSION",
	"OPT_CHAIN_MEMBER_LOOKUP_EXPRESSION",
	"PAREN_EXPRESSION",
	"FINALLY",
	"SUPER_EXPRESSION",
	"ARRAY_PATTERN",
	"OBJECT_PATTERN",
	"FORMAL_PARAMETER_LIST",
	"ITER_SPREAD",
	"OBJECT_SPREAD",
	"NULL",
	"ITER_REST",
	"OBJECT_REST",
	"MODULE_IMPORT",
	"EXPORT_DECLARATION",
	"EXPORT_SPECIFIER",
	"IMPORT_DECLARATION",
	"IMPORT_SPECIFIER",
	"FOR_OF_STATEMENT",
	"FOR_AWAIT_OF_STATEMENT",
	"YIELD_EXPRESSION",
	"DEFAULT_PARAMETER",
	"TEMPLATE_LITERAL_EXPRESSION",
	"TEMPLATE_LITERAL_PORTION",
	"TEMPLATE_SUBSTITUTION",
	"OPTIONAL_PARAMETER",
	"NEW_TARGET_EXPRESSION",
	"AWAIT_EXPRESSION",
	"DYNAMIC_IMPORT_EXPRESSION",
	"IMPORT_META_EXPRESSION",
	"COMPUTED_PROPERTY_FIELD",
	"FIELD_DECLARATION",
}

// ParseTreeType_values is ParseTreeType.values().
//
// Java: ParseTreeType.values()
func ParseTreeType_values() []ParseTreeType {
	values := make([]ParseTreeType, len(parseTreeType_names))
	for i := range values {
		values[i] = ParseTreeType(i)
	}
	return values
}

// Name is Enum.name().
//
// Java: ParseTreeType.name()
func (t ParseTreeType) Name() string { return parseTreeType_names[t] }

// String is Enum.toString().
//
// Java: ParseTreeType.toString()
func (t ParseTreeType) String() string { return t.Name() }

// Represents a null parse tree element.
type NullTree struct {
	ParseTreeBase
}

// NewNullTree constructs a NullTree.
//
// Java: NullTree(SourceRange)
func NewNullTree(location SourceRange) *NullTree {
	t := &NullTree{}
	t.InitParseTreeBase(ParseTreeType_NULL, location)
	return t
}

// MissingPrimaryExpressionTree is the parse tree for ParseTreeType MISSING_PRIMARY_EXPRESSION.
type MissingPrimaryExpressionTree struct {
	ParseTreeBase
}

// NewMissingPrimaryExpressionTree constructs a MissingPrimaryExpressionTree.
//
// Java: MissingPrimaryExpressionTree(SourceRange)
func NewMissingPrimaryExpressionTree(location SourceRange) *MissingPrimaryExpressionTree {
	t := &MissingPrimaryExpressionTree{}
	t.InitParseTreeBase(ParseTreeType_MISSING_PRIMARY_EXPRESSION, location)
	return t
}

// ProgramTree is the parse tree for ParseTreeType PROGRAM.
type ProgramTree struct {
	ParseTreeBase

	SourceElements []ParseTree
	SourceComments []*Comment
}

// NewProgramTree constructs a ProgramTree.
//
// Java: ProgramTree(SourceRange, ImmutableList, ImmutableList)
func NewProgramTree(location SourceRange, sourceElements []ParseTree, sourceComments []*Comment) *ProgramTree {
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
	DefaultBindingIdentifier *IdentifierToken
	// The list of names imported from the module when using the
	// `import {...} from 'specifier'` form.
	ImportSpecifierList []ParseTree
	// The namespace into which imported names are put when using the
	// `import * as nameSpaceImportIdentifier from 'specifier'` form.
	NameSpaceImportIdentifier *IdentifierToken
	// The string identifying the module.
	ModuleSpecifier LiteralToken
}

// NewImportDeclarationTree constructs a ImportDeclarationTree.
//
// Java: ImportDeclarationTree(SourceRange, IdentifierToken, ImmutableList, IdentifierToken, LiteralToken)
func NewImportDeclarationTree(location SourceRange, defaultBindingIdentifier *IdentifierToken, importSpecifierList []ParseTree, nameSpaceImportIdentifier *IdentifierToken, moduleSpecifier LiteralToken) *ImportDeclarationTree {
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

	ImportedName    *IdentifierToken
	DestinationName *IdentifierToken
}

// NewImportSpecifierTree constructs a ImportSpecifierTree.
//
// Java: ImportSpecifierTree(SourceRange, IdentifierToken, IdentifierToken)
func NewImportSpecifierTree(location SourceRange, importedName *IdentifierToken, destinationName *IdentifierToken) *ImportSpecifierTree {
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
	From                LiteralToken
}

// NewExportDeclarationTree constructs a ExportDeclarationTree.
//
// Java: ExportDeclarationTree(SourceRange, boolean, boolean, ParseTree, ImmutableList, LiteralToken)
func NewExportDeclarationTree(location SourceRange, isDefault bool, isExportAll bool, declaration ParseTree, exportSpecifierList []ParseTree, from LiteralToken) *ExportDeclarationTree {
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

	ImportedName    *IdentifierToken
	DestinationName *IdentifierToken
}

// NewExportSpecifierTree constructs a ExportSpecifierTree.
//
// Java: ExportSpecifierTree(SourceRange, IdentifierToken, IdentifierToken)
func NewExportSpecifierTree(location SourceRange, importedName *IdentifierToken, destinationName *IdentifierToken) *ExportSpecifierTree {
	t := &ExportSpecifierTree{}
	t.InitParseTreeBase(ParseTreeType_EXPORT_SPECIFIER, location)
	t.ImportedName = importedName
	t.DestinationName = destinationName
	return t
}

// ClassDeclarationTree is the parse tree for ParseTreeType CLASS_DECLARATION.
type ClassDeclarationTree struct {
	ParseTreeBase

	Name       *IdentifierToken
	SuperClass ParseTree
	Elements   []ParseTree
}

// NewClassDeclarationTree constructs a ClassDeclarationTree.
//
// Java: ClassDeclarationTree(SourceRange, IdentifierToken, ParseTree, ImmutableList)
func NewClassDeclarationTree(location SourceRange, name *IdentifierToken, superClass ParseTree, elements []ParseTree) *ClassDeclarationTree {
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

	Name                *IdentifierToken
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
	t.FormalParameterList = jutil.CheckNotNull(builder.formalParameterList)
	t.FunctionBody = jutil.CheckNotNull(builder.functionBody)
	t.IsAsync = builder.isAsync
	return t
}

// FunctionDeclarationTree_Builder builds a FunctionDeclarationTree.
type FunctionDeclarationTree_Builder struct {
	kind FunctionDeclarationTree_Kind

	name                *IdentifierToken
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
func (b *FunctionDeclarationTree_Builder) SetName(name *IdentifierToken) *FunctionDeclarationTree_Builder {
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

	Name        *IdentifierToken
	IsStatic    bool
	Initializer ParseTree
}

// NewFieldDeclarationTree constructs a FieldDeclarationTree.
//
// Java: FieldDeclarationTree(SourceRange, IdentifierToken, boolean, ParseTree)
func NewFieldDeclarationTree(location SourceRange, name *IdentifierToken, isStatic bool, initializer ParseTree) *FieldDeclarationTree {
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

// BlockTree is the parse tree for ParseTreeType BLOCK.
type BlockTree struct {
	ParseTreeBase

	Statements []ParseTree
}

// NewBlockTree constructs a BlockTree.
//
// Java: BlockTree(SourceRange, ImmutableList)
func NewBlockTree(location SourceRange, statements []ParseTree) *BlockTree {
	t := &BlockTree{}
	t.InitParseTreeBase(ParseTreeType_BLOCK, location)
	t.Statements = statements
	return t
}

// IfStatementTree is the parse tree for ParseTreeType IF_STATEMENT.
type IfStatementTree struct {
	ParseTreeBase

	Condition  ParseTree
	IfClause   ParseTree
	ElseClause ParseTree
}

// NewIfStatementTree constructs a IfStatementTree.
//
// Java: IfStatementTree(SourceRange, ParseTree, ParseTree, ParseTree)
func NewIfStatementTree(location SourceRange, condition ParseTree, ifClause ParseTree, elseClause ParseTree) *IfStatementTree {
	t := &IfStatementTree{}
	t.InitParseTreeBase(ParseTreeType_IF_STATEMENT, location)
	t.Condition = condition
	t.IfClause = ifClause
	t.ElseClause = elseClause
	return t
}

// Represents a for-await-of statement, which was added in ES2018
type ForAwaitOfStatementTree struct {
	ParseTreeBase

	Initializer ParseTree
	Collection  ParseTree
	Body        ParseTree
}

// NewForAwaitOfStatementTree constructs a ForAwaitOfStatementTree.
//
// Java: ForAwaitOfStatementTree(SourceRange, ParseTree, ParseTree, ParseTree)
func NewForAwaitOfStatementTree(location SourceRange, initializer ParseTree, collection ParseTree, body ParseTree) *ForAwaitOfStatementTree {
	t := &ForAwaitOfStatementTree{}
	t.InitParseTreeBase(ParseTreeType_FOR_AWAIT_OF_STATEMENT, location)
	t.Initializer = initializer
	t.Collection = collection
	t.Body = body
	return t
}

// ForInStatementTree is the parse tree for ParseTreeType FOR_IN_STATEMENT.
type ForInStatementTree struct {
	ParseTreeBase

	Initializer ParseTree
	Collection  ParseTree
	Body        ParseTree
}

// NewForInStatementTree constructs a ForInStatementTree.
//
// Java: ForInStatementTree(SourceRange, ParseTree, ParseTree, ParseTree)
func NewForInStatementTree(location SourceRange, initializer ParseTree, collection ParseTree, body ParseTree) *ForInStatementTree {
	t := &ForInStatementTree{}
	t.InitParseTreeBase(ParseTreeType_FOR_IN_STATEMENT, location)
	t.Initializer = initializer
	t.Collection = collection
	t.Body = body
	return t
}

// ForOfStatementTree is the parse tree for ParseTreeType FOR_OF_STATEMENT.
type ForOfStatementTree struct {
	ParseTreeBase

	Initializer ParseTree
	Collection  ParseTree
	Body        ParseTree
}

// NewForOfStatementTree constructs a ForOfStatementTree.
//
// Java: ForOfStatementTree(SourceRange, ParseTree, ParseTree, ParseTree)
func NewForOfStatementTree(location SourceRange, initializer ParseTree, collection ParseTree, body ParseTree) *ForOfStatementTree {
	t := &ForOfStatementTree{}
	t.InitParseTreeBase(ParseTreeType_FOR_OF_STATEMENT, location)
	t.Initializer = initializer
	t.Collection = collection
	t.Body = body
	return t
}

// ForStatementTree is the parse tree for ParseTreeType FOR_STATEMENT.
type ForStatementTree struct {
	ParseTreeBase

	Initializer ParseTree
	Condition   ParseTree
	Increment   ParseTree
	Body        ParseTree
}

// NewForStatementTree constructs a ForStatementTree.
//
// Java: ForStatementTree(SourceRange, ParseTree, ParseTree, ParseTree, ParseTree)
func NewForStatementTree(location SourceRange, variables ParseTree, condition ParseTree, increment ParseTree, body ParseTree) *ForStatementTree {
	t := &ForStatementTree{}
	t.InitParseTreeBase(ParseTreeType_FOR_STATEMENT, location)
	t.Initializer = variables
	t.Condition = condition
	t.Increment = increment
	t.Body = body
	return t
}

// WhileStatementTree is the parse tree for ParseTreeType WHILE_STATEMENT.
type WhileStatementTree struct {
	ParseTreeBase

	Condition ParseTree
	Body      ParseTree
}

// NewWhileStatementTree constructs a WhileStatementTree.
//
// Java: WhileStatementTree(SourceRange, ParseTree, ParseTree)
func NewWhileStatementTree(location SourceRange, condition ParseTree, body ParseTree) *WhileStatementTree {
	t := &WhileStatementTree{}
	t.InitParseTreeBase(ParseTreeType_WHILE_STATEMENT, location)
	t.Condition = condition
	t.Body = body
	return t
}

// DoWhileStatementTree is the parse tree for ParseTreeType DO_WHILE_STATEMENT.
type DoWhileStatementTree struct {
	ParseTreeBase

	Body      ParseTree
	Condition ParseTree
}

// NewDoWhileStatementTree constructs a DoWhileStatementTree.
//
// Java: DoWhileStatementTree(SourceRange, ParseTree, ParseTree)
func NewDoWhileStatementTree(location SourceRange, body ParseTree, condition ParseTree) *DoWhileStatementTree {
	t := &DoWhileStatementTree{}
	t.InitParseTreeBase(ParseTreeType_DO_WHILE_STATEMENT, location)
	t.Body = body
	t.Condition = condition
	return t
}

// SwitchStatementTree is the parse tree for ParseTreeType SWITCH_STATEMENT.
type SwitchStatementTree struct {
	ParseTreeBase

	Expression  ParseTree
	CaseClauses []ParseTree
}

// NewSwitchStatementTree constructs a SwitchStatementTree.
//
// Java: SwitchStatementTree(SourceRange, ParseTree, ImmutableList)
func NewSwitchStatementTree(location SourceRange, expression ParseTree, caseClauses []ParseTree) *SwitchStatementTree {
	t := &SwitchStatementTree{}
	t.InitParseTreeBase(ParseTreeType_SWITCH_STATEMENT, location)
	t.Expression = expression
	t.CaseClauses = caseClauses
	return t
}

// CaseClauseTree is the parse tree for ParseTreeType CASE_CLAUSE.
type CaseClauseTree struct {
	ParseTreeBase

	Expression ParseTree
	Statements []ParseTree
}

// NewCaseClauseTree constructs a CaseClauseTree.
//
// Java: CaseClauseTree(SourceRange, ParseTree, ImmutableList)
func NewCaseClauseTree(location SourceRange, expression ParseTree, statements []ParseTree) *CaseClauseTree {
	t := &CaseClauseTree{}
	t.InitParseTreeBase(ParseTreeType_CASE_CLAUSE, location)
	t.Expression = expression
	t.Statements = statements
	return t
}

// DefaultClauseTree is the parse tree for ParseTreeType DEFAULT_CLAUSE.
type DefaultClauseTree struct {
	ParseTreeBase

	Statements []ParseTree
}

// NewDefaultClauseTree constructs a DefaultClauseTree.
//
// Java: DefaultClauseTree(SourceRange, ImmutableList)
func NewDefaultClauseTree(location SourceRange, statements []ParseTree) *DefaultClauseTree {
	t := &DefaultClauseTree{}
	t.InitParseTreeBase(ParseTreeType_DEFAULT_CLAUSE, location)
	t.Statements = statements
	return t
}

// TryStatementTree is the parse tree for ParseTreeType TRY_STATEMENT.
type TryStatementTree struct {
	ParseTreeBase

	Body         ParseTree
	CatchBlock   ParseTree
	FinallyBlock ParseTree
}

// NewTryStatementTree constructs a TryStatementTree.
//
// Java: TryStatementTree(SourceRange, ParseTree, ParseTree, ParseTree)
func NewTryStatementTree(location SourceRange, body ParseTree, catchBlock ParseTree, finallyBlock ParseTree) *TryStatementTree {
	t := &TryStatementTree{}
	t.InitParseTreeBase(ParseTreeType_TRY_STATEMENT, location)
	t.Body = body
	t.CatchBlock = catchBlock
	t.FinallyBlock = finallyBlock
	return t
}

// CatchTree is the parse tree for ParseTreeType CATCH.
type CatchTree struct {
	ParseTreeBase

	Exception ParseTree
	CatchBody ParseTree
}

// NewCatchTree constructs a CatchTree.
//
// Java: CatchTree(SourceRange, ParseTree, ParseTree)
func NewCatchTree(location SourceRange, exception ParseTree, catchBody ParseTree) *CatchTree {
	t := &CatchTree{}
	t.InitParseTreeBase(ParseTreeType_CATCH, location)
	t.Exception = exception
	t.CatchBody = catchBody
	return t
}

// FinallyTree is the parse tree for ParseTreeType FINALLY.
type FinallyTree struct {
	ParseTreeBase

	Block ParseTree
}

// NewFinallyTree constructs a FinallyTree.
//
// Java: FinallyTree(SourceRange, ParseTree)
func NewFinallyTree(location SourceRange, block ParseTree) *FinallyTree {
	t := &FinallyTree{}
	t.InitParseTreeBase(ParseTreeType_FINALLY, location)
	t.Block = block
	return t
}

// ReturnStatementTree is the parse tree for ParseTreeType RETURN_STATEMENT.
type ReturnStatementTree struct {
	ParseTreeBase

	Expression ParseTree
}

// NewReturnStatementTree constructs a ReturnStatementTree.
//
// Java: ReturnStatementTree(SourceRange, ParseTree)
func NewReturnStatementTree(location SourceRange, expression ParseTree) *ReturnStatementTree {
	t := &ReturnStatementTree{}
	t.InitParseTreeBase(ParseTreeType_RETURN_STATEMENT, location)
	t.Expression = expression
	return t
}

// BreakStatementTree is the parse tree for ParseTreeType BREAK_STATEMENT.
type BreakStatementTree struct {
	ParseTreeBase

	Name *IdentifierToken
}

// NewBreakStatementTree constructs a BreakStatementTree.
//
// Java: BreakStatementTree(SourceRange, IdentifierToken)
func NewBreakStatementTree(location SourceRange, name *IdentifierToken) *BreakStatementTree {
	t := &BreakStatementTree{}
	t.InitParseTreeBase(ParseTreeType_BREAK_STATEMENT, location)
	t.Name = name
	return t
}

// HasLabel reports whether the statement names a label.
//
// Java: BreakStatementTree.hasLabel()
func (t *BreakStatementTree) HasLabel() bool {
	return t.Name != nil
}

// ContinueStatementTree is the parse tree for ParseTreeType CONTINUE_STATEMENT.
type ContinueStatementTree struct {
	ParseTreeBase

	Name *IdentifierToken
}

// NewContinueStatementTree constructs a ContinueStatementTree.
//
// Java: ContinueStatementTree(SourceRange, IdentifierToken)
func NewContinueStatementTree(location SourceRange, name *IdentifierToken) *ContinueStatementTree {
	t := &ContinueStatementTree{}
	t.InitParseTreeBase(ParseTreeType_CONTINUE_STATEMENT, location)
	t.Name = name
	return t
}

// HasLabel reports whether the statement names a label.
//
// Java: ContinueStatementTree.hasLabel()
func (t *ContinueStatementTree) HasLabel() bool {
	return t.Name != nil
}

// ThrowStatementTree is the parse tree for ParseTreeType THROW_STATEMENT.
type ThrowStatementTree struct {
	ParseTreeBase

	Value ParseTree
}

// NewThrowStatementTree constructs a ThrowStatementTree.
//
// Java: ThrowStatementTree(SourceRange, ParseTree)
func NewThrowStatementTree(location SourceRange, value ParseTree) *ThrowStatementTree {
	t := &ThrowStatementTree{}
	t.InitParseTreeBase(ParseTreeType_THROW_STATEMENT, location)
	t.Value = value
	return t
}

// LabelledStatementTree is the parse tree for ParseTreeType LABELLED_STATEMENT.
type LabelledStatementTree struct {
	ParseTreeBase

	Name      *IdentifierToken
	Statement ParseTree
}

// NewLabelledStatementTree constructs a LabelledStatementTree.
//
// Java: LabelledStatementTree(SourceRange, IdentifierToken, ParseTree)
func NewLabelledStatementTree(location SourceRange, name *IdentifierToken, statement ParseTree) *LabelledStatementTree {
	t := &LabelledStatementTree{}
	t.InitParseTreeBase(ParseTreeType_LABELLED_STATEMENT, location)
	t.Name = name
	t.Statement = statement
	return t
}

// WithStatementTree is the parse tree for ParseTreeType WITH_STATEMENT.
type WithStatementTree struct {
	ParseTreeBase

	Expression ParseTree
	Body       ParseTree
}

// NewWithStatementTree constructs a WithStatementTree.
//
// Java: WithStatementTree(SourceRange, ParseTree, ParseTree)
func NewWithStatementTree(location SourceRange, expression ParseTree, body ParseTree) *WithStatementTree {
	t := &WithStatementTree{}
	t.InitParseTreeBase(ParseTreeType_WITH_STATEMENT, location)
	t.Expression = expression
	t.Body = body
	return t
}

// VariableStatementTree is the parse tree for ParseTreeType VARIABLE_STATEMENT.
type VariableStatementTree struct {
	ParseTreeBase

	Declarations *VariableDeclarationListTree
}

// NewVariableStatementTree constructs a VariableStatementTree.
//
// Java: VariableStatementTree(SourceRange, VariableDeclarationListTree)
func NewVariableStatementTree(location SourceRange, declarations *VariableDeclarationListTree) *VariableStatementTree {
	t := &VariableStatementTree{}
	t.InitParseTreeBase(ParseTreeType_VARIABLE_STATEMENT, location)
	t.Declarations = declarations
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

// ExpressionStatementTree is the parse tree for ParseTreeType EXPRESSION_STATEMENT.
type ExpressionStatementTree struct {
	ParseTreeBase

	Expression ParseTree
}

// NewExpressionStatementTree constructs a ExpressionStatementTree.
//
// Java: ExpressionStatementTree(SourceRange, ParseTree)
func NewExpressionStatementTree(location SourceRange, expression ParseTree) *ExpressionStatementTree {
	t := &ExpressionStatementTree{}
	t.InitParseTreeBase(ParseTreeType_EXPRESSION_STATEMENT, location)
	t.Expression = expression
	return t
}

// EmptyStatementTree is the parse tree for ParseTreeType EMPTY_STATEMENT.
type EmptyStatementTree struct {
	ParseTreeBase
}

// NewEmptyStatementTree constructs a EmptyStatementTree.
//
// Java: EmptyStatementTree(SourceRange)
func NewEmptyStatementTree(location SourceRange) *EmptyStatementTree {
	t := &EmptyStatementTree{}
	t.InitParseTreeBase(ParseTreeType_EMPTY_STATEMENT, location)
	return t
}

// DebuggerStatementTree is the parse tree for ParseTreeType DEBUGGER_STATEMENT.
type DebuggerStatementTree struct {
	ParseTreeBase
}

// NewDebuggerStatementTree constructs a DebuggerStatementTree.
//
// Java: DebuggerStatementTree(SourceRange)
func NewDebuggerStatementTree(location SourceRange) *DebuggerStatementTree {
	t := &DebuggerStatementTree{}
	t.InitParseTreeBase(ParseTreeType_DEBUGGER_STATEMENT, location)
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

// Comment is a placeholder class.
type Comment struct {
	Value    string
	Location LineRange // Go-only: with lines and columns (see SourcePosition)
	Type     Comment_Type
}

// Comment_Type is Comment.Type.
type Comment_Type int32

const (
	// /* comment */
	Comment_Type_BLOCK Comment_Type = iota

	// // comment
	Comment_Type_LINE

	// /** comment */
	Comment_Type_JSDOC

	// #!/usr/bin/node
	// Only valid at the start of a file.
	Comment_Type_SHEBANG

	// /*! comment */
	Comment_Type_IMPORTANT
)

var comment_Type_names = [...]string{"BLOCK", "LINE", "JSDOC", "SHEBANG", "IMPORTANT"}

// Comment_Type_values is Comment.Type.values().
//
// Java: Comment.Type.values()
func Comment_Type_values() []Comment_Type {
	return []Comment_Type{Comment_Type_BLOCK, Comment_Type_LINE, Comment_Type_JSDOC, Comment_Type_SHEBANG, Comment_Type_IMPORTANT}
}

// Name is Enum.name().
//
// Java: Comment.Type.name()
func (t Comment_Type) Name() string { return comment_Type_names[t] }

// String is Enum.toString().
//
// Java: Comment.Type.toString()
func (t Comment_Type) String() string { return t.Name() }

// NewComment constructs a Comment.
//
// Java: Comment(String, SourceRange, Type)
func NewComment(value string, location LineRange, type_ Comment_Type) *Comment {
	c := &Comment{}
	c.Value = value
	c.Location = location
	c.Type = type_
	return c
}

// IsJsDoc reports whether this is a JSDoc comment.
//
// Java: Comment.isJsDoc()
func (c *Comment) IsJsDoc() bool {
	return c.Type == Comment_Type_JSDOC
}

// GetAbsolutePosition returns the start offset.
//
// Java: Comment.getAbsolutePosition()
func (c *Comment) GetAbsolutePosition() int {
	return int(c.Location.Start.Offset)
}

// GetLength returns the length in UTF-16 code units.
//
// Java: Comment.getLength()
func (c *Comment) GetLength() int {
	return int(c.Location.End.Offset - c.Location.Start.Offset)
}
