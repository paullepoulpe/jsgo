// Port of com/google/javascript/jscomp/parsing/parser/trees/ParseTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/ParseTreeType.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/NullTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/MissingPrimaryExpressionTree.java @ closure-compiler v20261021
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
// Port of com/google/javascript/jscomp/parsing/parser/trees/ExpressionStatementTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/EmptyStatementTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/DebuggerStatementTree.java @ closure-compiler v20261021
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

func unwrapParens(tree ParseTree) ParseTreeType {
	for {
		paren, ok := tree.(*ParenExpressionTree)
		if !ok {
			return tree.GetParseTreeBase().Type
		}
		tree = paren.Expression
	}
}

// IsPattern reports whether the tree (inside any parentheses) is an array or
// object pattern.
//
// Java: ParseTree.isPattern()
func IsPattern(tree ParseTree) bool {
	switch unwrapParens(tree) {
	case ParseTreeType_ARRAY_PATTERN, ParseTreeType_OBJECT_PATTERN:
		return true
	default:
		return false
	}
}

// IsValidNonVanillaAssignmentTarget: is valid assignment target for
// non-vanilla assignment operators like `+=`, `-+`, `**=`, etc
//
// Java: ParseTree.isValidNonVanillaAssignmentTarget()
func IsValidNonVanillaAssignmentTarget(tree ParseTree) bool {
	switch unwrapParens(tree) {
	case ParseTreeType_IDENTIFIER_EXPRESSION, ParseTreeType_MEMBER_EXPRESSION, ParseTreeType_MEMBER_LOOKUP_EXPRESSION, ParseTreeType_DEFAULT_PARAMETER:
		return true
	default:
		return false
	}
}

// IsValidAssignmentTarget: is valid assignment target for any assignment
// operator like `=`, `+=`, `-+`, `**=`, etc
//
// Java: ParseTree.isValidAssignmentTarget()
func IsValidAssignmentTarget(tree ParseTree) bool {
	switch unwrapParens(tree) {
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
	return t.Type.String() + "@" + t.Location.String()
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

	Name Token
}

// NewBreakStatementTree constructs a BreakStatementTree.
//
// Java: BreakStatementTree(SourceRange, IdentifierToken)
func NewBreakStatementTree(location SourceRange, name Token) *BreakStatementTree {
	t := &BreakStatementTree{}
	t.InitParseTreeBase(ParseTreeType_BREAK_STATEMENT, location)
	t.Name = name
	return t
}

// HasLabel reports whether the statement names a label.
//
// Java: BreakStatementTree.hasLabel()
func (t *BreakStatementTree) HasLabel() bool {
	return t.Name.IsPresent()
}

// ContinueStatementTree is the parse tree for ParseTreeType CONTINUE_STATEMENT.
type ContinueStatementTree struct {
	ParseTreeBase

	Name Token
}

// NewContinueStatementTree constructs a ContinueStatementTree.
//
// Java: ContinueStatementTree(SourceRange, IdentifierToken)
func NewContinueStatementTree(location SourceRange, name Token) *ContinueStatementTree {
	t := &ContinueStatementTree{}
	t.InitParseTreeBase(ParseTreeType_CONTINUE_STATEMENT, location)
	t.Name = name
	return t
}

// HasLabel reports whether the statement names a label.
//
// Java: ContinueStatementTree.hasLabel()
func (t *ContinueStatementTree) HasLabel() bool {
	return t.Name.IsPresent()
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

	Name      Token
	Statement ParseTree
}

// NewLabelledStatementTree constructs a LabelledStatementTree.
//
// Java: LabelledStatementTree(SourceRange, IdentifierToken, ParseTree)
func NewLabelledStatementTree(location SourceRange, name Token, statement ParseTree) *LabelledStatementTree {
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
func NewComment(value string, location LineRange, type_ Comment_Type) Comment {
	return Comment{
		Value:    value,
		Location: location,
		Type:     type_,
	}
}

// IsJsDoc reports whether this is a JSDoc comment.
//
// Java: Comment.isJsDoc()
func (c Comment) IsJsDoc() bool {
	return c.Type == Comment_Type_JSDOC
}

// GetAbsolutePosition returns the start offset.
//
// Java: Comment.getAbsolutePosition()
func (c Comment) GetAbsolutePosition() int {
	return int(c.Location.Start.Offset)
}

// GetLength returns the length in UTF-16 code units.
//
// Java: Comment.getLength()
func (c Comment) GetLength() int {
	return int(c.Location.End.Offset - c.Location.Start.Offset)
}
