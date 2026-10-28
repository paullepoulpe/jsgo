// Port of com/google/javascript/jscomp/parsing/parser/trees/ArgumentListTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/ArrayLiteralExpressionTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/ArrayPatternTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/AwaitExpressionTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/BinaryOperatorTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/CallExpressionTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/CommaExpressionTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/ConditionalExpressionTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/DefaultParameterTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/DynamicImportTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/IdentifierExpressionTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/IterRestTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/IterSpreadTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/LiteralExpressionTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/MemberExpressionTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/MemberLookupExpressionTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/NewExpressionTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/NewTargetExpressionTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/ObjectLiteralExpressionTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/ObjectPatternTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/ObjectRestTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/ObjectSpreadTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/OptChainCallExpressionTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/OptionalMemberExpressionTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/OptionalMemberLookupExpressionTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/ParenExpressionTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/PropertyNameAssignmentTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/SuperExpressionTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/TemplateLiteralExpressionTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/TemplateLiteralPortionTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/TemplateSubstitutionTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/ThisExpressionTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/UnaryExpressionTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/UpdateExpressionTree.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/trees/YieldExpressionTree.java @ closure-compiler v20261021

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
 * Copyright 2016 The Closure Compiler Authors.
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

// ArgumentListTree is the parse tree for ParseTreeType ARGUMENT_LIST.
type ArgumentListTree struct {
	ParseTreeBase
	// Go-only: performance — before the other fields, in the base's padding.
	HasTrailingComma bool

	Arguments      []ParseTree
	CommaPositions []SourcePosition
}

// NewArgumentListTree constructs a ArgumentListTree.
//
// Java: ArgumentListTree(SourceRange, ImmutableList, boolean, ImmutableList)
func NewArgumentListTree(location SourceRange, arguments []ParseTree, hasTrailingComma bool, commaPositions []SourcePosition) *ArgumentListTree {
	t := &ArgumentListTree{}
	t.InitParseTreeBase(ParseTreeType_ARGUMENT_LIST, location)
	t.Arguments = arguments
	t.HasTrailingComma = hasTrailingComma
	t.CommaPositions = commaPositions
	return t
}

// ArrayLiteralExpressionTree is the parse tree for ParseTreeType ARRAY_LITERAL_EXPRESSION.
type ArrayLiteralExpressionTree struct {
	ParseTreeBase

	Elements         []ParseTree
	HasTrailingComma bool
}

// NewArrayLiteralExpressionTree constructs a ArrayLiteralExpressionTree.
//
// Java: ArrayLiteralExpressionTree(SourceRange, ImmutableList, boolean)
func NewArrayLiteralExpressionTree(location SourceRange, elements []ParseTree, hasTrailingComma bool) *ArrayLiteralExpressionTree {
	t := &ArrayLiteralExpressionTree{}
	t.InitParseTreeBase(ParseTreeType_ARRAY_LITERAL_EXPRESSION, location)
	t.Elements = elements
	t.HasTrailingComma = hasTrailingComma
	return t
}

// ArrayPatternTree is the parse tree for ParseTreeType ARRAY_PATTERN.
type ArrayPatternTree struct {
	ParseTreeBase

	Elements []ParseTree
}

// NewArrayPatternTree constructs a ArrayPatternTree.
//
// Java: ArrayPatternTree(SourceRange, ImmutableList)
func NewArrayPatternTree(location SourceRange, elements []ParseTree) *ArrayPatternTree {
	t := &ArrayPatternTree{}
	t.InitParseTreeBase(ParseTreeType_ARRAY_PATTERN, location)
	t.Elements = elements
	return t
}

// Represents an await expression.
//
//	AwaitExpression : await UnaryExpression
type AwaitExpressionTree struct {
	ParseTreeBase

	Expression ParseTree
}

// NewAwaitExpressionTree constructs a AwaitExpressionTree.
//
// Java: AwaitExpressionTree(SourceRange, ParseTree)
func NewAwaitExpressionTree(location SourceRange, expression ParseTree) *AwaitExpressionTree {
	t := &AwaitExpressionTree{}
	t.InitParseTreeBase(ParseTreeType_AWAIT_EXPRESSION, location)
	t.Expression = expression
	return t
}

// BinaryOperatorTree is the parse tree for ParseTreeType BINARY_OPERATOR.
type BinaryOperatorTree struct {
	ParseTreeBase

	Left     ParseTree
	Operator Token
	Right    ParseTree
}

// NewBinaryOperatorTree constructs a BinaryOperatorTree.
//
// Java: BinaryOperatorTree(SourceRange, ParseTree, Token, ParseTree)
func NewBinaryOperatorTree(location SourceRange, left ParseTree, operator Token, right ParseTree) *BinaryOperatorTree {
	t := &BinaryOperatorTree{}
	t.InitParseTreeBase(ParseTreeType_BINARY_OPERATOR, location)
	t.Left = left
	t.Operator = operator
	t.Right = right
	return t
}

// CallExpressionTree is the parse tree for ParseTreeType CALL_EXPRESSION.
type CallExpressionTree struct {
	ParseTreeBase

	Operand   ParseTree
	Arguments *ArgumentListTree
}

// NewCallExpressionTree constructs a CallExpressionTree.
//
// Java: CallExpressionTree(SourceRange, ParseTree, ArgumentListTree)
func NewCallExpressionTree(location SourceRange, operand ParseTree, arguments *ArgumentListTree) *CallExpressionTree {
	t := &CallExpressionTree{}
	t.InitParseTreeBase(ParseTreeType_CALL_EXPRESSION, location)
	t.Operand = operand
	t.Arguments = arguments
	return t
}

// CommaExpressionTree is the parse tree for ParseTreeType COMMA_EXPRESSION.
type CommaExpressionTree struct {
	ParseTreeBase

	Expressions []ParseTree
}

// NewCommaExpressionTree constructs a CommaExpressionTree.
//
// Java: CommaExpressionTree(SourceRange, ImmutableList)
func NewCommaExpressionTree(location SourceRange, expressions []ParseTree) *CommaExpressionTree {
	t := &CommaExpressionTree{}
	t.InitParseTreeBase(ParseTreeType_COMMA_EXPRESSION, location)
	t.Expressions = expressions
	return t
}

// ConditionalExpressionTree is the parse tree for ParseTreeType CONDITIONAL_EXPRESSION.
type ConditionalExpressionTree struct {
	ParseTreeBase

	Condition ParseTree
	Left      ParseTree
	Right     ParseTree
}

// NewConditionalExpressionTree constructs a ConditionalExpressionTree.
//
// Java: ConditionalExpressionTree(SourceRange, ParseTree, ParseTree, ParseTree)
func NewConditionalExpressionTree(location SourceRange, condition ParseTree, left ParseTree, right ParseTree) *ConditionalExpressionTree {
	t := &ConditionalExpressionTree{}
	t.InitParseTreeBase(ParseTreeType_CONDITIONAL_EXPRESSION, location)
	t.Condition = condition
	t.Left = left
	t.Right = right
	return t
}

// DefaultParameterTree is the parse tree for ParseTreeType DEFAULT_PARAMETER.
type DefaultParameterTree struct {
	ParseTreeBase

	Lhs          ParseTree
	DefaultValue ParseTree
}

// NewDefaultParameterTree constructs a DefaultParameterTree.
//
// Java: DefaultParameterTree(SourceRange, ParseTree, ParseTree)
func NewDefaultParameterTree(location SourceRange, lhs ParseTree, defaultValue ParseTree) *DefaultParameterTree {
	t := &DefaultParameterTree{}
	t.InitParseTreeBase(ParseTreeType_DEFAULT_PARAMETER, location)
	t.Lhs = lhs
	t.DefaultValue = defaultValue
	return t
}

// Represents a dynamic import expression.
//
// e.g. `import(myModuleName)` in `const modulePromise = import(myModuleName)`.
type DynamicImportTree struct {
	ParseTreeBase

	Argument ParseTree
}

// NewDynamicImportTree constructs a DynamicImportTree.
//
// Java: DynamicImportTree(SourceRange, ParseTree)
func NewDynamicImportTree(location SourceRange, argument ParseTree) *DynamicImportTree {
	t := &DynamicImportTree{}
	t.InitParseTreeBase(ParseTreeType_DYNAMIC_IMPORT_EXPRESSION, location)
	t.Argument = argument
	return t
}

// IdentifierExpressionTree is the parse tree for ParseTreeType IDENTIFIER_EXPRESSION.
type IdentifierExpressionTree struct {
	ParseTreeBase

	IdentifierToken Token
}

// NewIdentifierExpressionTree constructs a IdentifierExpressionTree.
//
// Java: IdentifierExpressionTree(SourceRange, IdentifierToken)
func NewIdentifierExpressionTree(location SourceRange, identifierToken Token) *IdentifierExpressionTree {
	t := &IdentifierExpressionTree{}
	t.InitParseTreeBase(ParseTreeType_IDENTIFIER_EXPRESSION, location)
	t.IdentifierToken = identifierToken
	return t
}

// IterRestTree is the parse tree for ParseTreeType ITER_REST.
type IterRestTree struct {
	ParseTreeBase

	AssignmentTarget ParseTree
}

// NewIterRestTree constructs a IterRestTree.
//
// Java: IterRestTree(SourceRange, ParseTree)
func NewIterRestTree(location SourceRange, assignmentTarget ParseTree) *IterRestTree {
	t := &IterRestTree{}
	t.InitParseTreeBase(ParseTreeType_ITER_REST, location)
	t.AssignmentTarget = assignmentTarget
	return t
}

// IterSpreadTree is the parse tree for ParseTreeType ITER_SPREAD.
type IterSpreadTree struct {
	ParseTreeBase

	Expression ParseTree
}

// NewIterSpreadTree constructs a IterSpreadTree.
//
// Java: IterSpreadTree(SourceRange, ParseTree)
func NewIterSpreadTree(location SourceRange, expression ParseTree) *IterSpreadTree {
	t := &IterSpreadTree{}
	t.InitParseTreeBase(ParseTreeType_ITER_SPREAD, location)
	t.Expression = expression
	return t
}

// LiteralExpressionTree is the parse tree for ParseTreeType LITERAL_EXPRESSION.
type LiteralExpressionTree struct {
	ParseTreeBase

	LiteralToken Token
}

// NewLiteralExpressionTree constructs a LiteralExpressionTree.
//
// Java: LiteralExpressionTree(SourceRange, Token)
func NewLiteralExpressionTree(location SourceRange, literalToken Token) *LiteralExpressionTree {
	t := &LiteralExpressionTree{}
	t.InitParseTreeBase(ParseTreeType_LITERAL_EXPRESSION, location)
	t.LiteralToken = literalToken
	return t
}

// MemberExpressionTree is the parse tree for ParseTreeType MEMBER_EXPRESSION.
type MemberExpressionTree struct {
	ParseTreeBase

	Operand    ParseTree
	MemberName Token
}

// NewMemberExpressionTree constructs a MemberExpressionTree.
//
// Java: MemberExpressionTree(SourceRange, ParseTree, IdentifierToken)
func NewMemberExpressionTree(location SourceRange, operand ParseTree, memberName Token) *MemberExpressionTree {
	t := &MemberExpressionTree{}
	t.InitParseTreeBase(ParseTreeType_MEMBER_EXPRESSION, location)
	t.Operand = operand
	t.MemberName = memberName
	return t
}

// MemberLookupExpressionTree is the parse tree for ParseTreeType MEMBER_LOOKUP_EXPRESSION.
type MemberLookupExpressionTree struct {
	ParseTreeBase

	Operand          ParseTree
	MemberExpression ParseTree
}

// NewMemberLookupExpressionTree constructs a MemberLookupExpressionTree.
//
// Java: MemberLookupExpressionTree(SourceRange, ParseTree, ParseTree)
func NewMemberLookupExpressionTree(location SourceRange, operand ParseTree, memberExpression ParseTree) *MemberLookupExpressionTree {
	t := &MemberLookupExpressionTree{}
	t.InitParseTreeBase(ParseTreeType_MEMBER_LOOKUP_EXPRESSION, location)
	t.Operand = operand
	t.MemberExpression = memberExpression
	return t
}

// NewExpressionTree is the parse tree for ParseTreeType NEW_EXPRESSION.
type NewExpressionTree struct {
	ParseTreeBase

	Operand          ParseTree
	Arguments        *ArgumentListTree
	HasTrailingComma bool
}

// NewNewExpressionTree constructs a NewExpressionTree.
//
// Java: NewExpressionTree(SourceRange, ParseTree, ArgumentListTree, boolean)
func NewNewExpressionTree(location SourceRange, operand ParseTree, arguments *ArgumentListTree, hasTrailingComma bool) *NewExpressionTree {
	t := &NewExpressionTree{}
	t.InitParseTreeBase(ParseTreeType_NEW_EXPRESSION, location)
	t.Operand = operand
	t.Arguments = arguments
	t.HasTrailingComma = hasTrailingComma
	return t
}

// Represents 'new.target'.
type NewTargetExpressionTree struct {
	ParseTreeBase
}

// NewNewTargetExpressionTree constructs a NewTargetExpressionTree.
//
// Java: NewTargetExpressionTree(SourceRange)
func NewNewTargetExpressionTree(location SourceRange) *NewTargetExpressionTree {
	t := &NewTargetExpressionTree{}
	t.InitParseTreeBase(ParseTreeType_NEW_TARGET_EXPRESSION, location)
	return t
}

// ObjectLiteralExpressionTree is the parse tree for ParseTreeType OBJECT_LITERAL_EXPRESSION.
type ObjectLiteralExpressionTree struct {
	ParseTreeBase

	PropertyNameAndValues []ParseTree
	HasTrailingComma      bool
}

// NewObjectLiteralExpressionTree constructs a ObjectLiteralExpressionTree.
//
// Java: ObjectLiteralExpressionTree(SourceRange, ImmutableList, boolean)
func NewObjectLiteralExpressionTree(location SourceRange, propertyNameAndValues []ParseTree, hasTrailingComma bool) *ObjectLiteralExpressionTree {
	t := &ObjectLiteralExpressionTree{}
	t.InitParseTreeBase(ParseTreeType_OBJECT_LITERAL_EXPRESSION, location)
	t.PropertyNameAndValues = propertyNameAndValues
	t.HasTrailingComma = hasTrailingComma
	return t
}

// ObjectPatternTree is the parse tree for ParseTreeType OBJECT_PATTERN.
type ObjectPatternTree struct {
	ParseTreeBase

	Fields []ParseTree
}

// NewObjectPatternTree constructs a ObjectPatternTree.
//
// Java: ObjectPatternTree(SourceRange, ImmutableList)
func NewObjectPatternTree(location SourceRange, fields []ParseTree) *ObjectPatternTree {
	t := &ObjectPatternTree{}
	t.InitParseTreeBase(ParseTreeType_OBJECT_PATTERN, location)
	t.Fields = fields
	return t
}

// ObjectRestTree is the parse tree for ParseTreeType OBJECT_REST.
type ObjectRestTree struct {
	ParseTreeBase

	AssignmentTarget ParseTree
}

// NewObjectRestTree constructs a ObjectRestTree.
//
// Java: ObjectRestTree(SourceRange, ParseTree)
func NewObjectRestTree(location SourceRange, assignmentTarget ParseTree) *ObjectRestTree {
	t := &ObjectRestTree{}
	t.InitParseTreeBase(ParseTreeType_OBJECT_REST, location)
	t.AssignmentTarget = assignmentTarget
	return t
}

// ObjectSpreadTree is the parse tree for ParseTreeType OBJECT_SPREAD.
type ObjectSpreadTree struct {
	ParseTreeBase

	Expression ParseTree
}

// NewObjectSpreadTree constructs a ObjectSpreadTree.
//
// Java: ObjectSpreadTree(SourceRange, ParseTree)
func NewObjectSpreadTree(location SourceRange, expression ParseTree) *ObjectSpreadTree {
	t := &ObjectSpreadTree{}
	t.InitParseTreeBase(ParseTreeType_OBJECT_SPREAD, location)
	t.Expression = expression
	return t
}

// Creates a OptChainCallExpressionTree that represents a call expression within an optional chain.
// Has an isOptionalChainStart field to indicate whether it is the start of an optional chain.
type OptChainCallExpressionTree struct {
	ParseTreeBase

	Operand                ParseTree
	Arguments              *ArgumentListTree
	IsStartOfOptionalChain bool
	HasTrailingComma       bool
}

// NewOptChainCallExpressionTree constructs a OptChainCallExpressionTree.
//
// Java: OptChainCallExpressionTree(SourceRange, ParseTree, ArgumentListTree, boolean, boolean)
func NewOptChainCallExpressionTree(location SourceRange, operand ParseTree, arguments *ArgumentListTree, isStartOfOptionalChain bool, hasTrailingComma bool) *OptChainCallExpressionTree {
	t := &OptChainCallExpressionTree{}
	t.InitParseTreeBase(ParseTreeType_OPT_CHAIN_CALL_EXPRESSION, location)
	t.Operand = operand
	t.Arguments = arguments
	t.IsStartOfOptionalChain = isStartOfOptionalChain
	t.HasTrailingComma = hasTrailingComma
	return t
}

// Creates a OptionalMemberExpressionTree that represents a getprop expression within an optional
// chain. Has an isOptionalChainStart field to indicate whether it is the start of an optional
// chain.
type OptionalMemberExpressionTree struct {
	ParseTreeBase

	Operand                ParseTree
	MemberName             Token
	IsStartOfOptionalChain bool
}

// NewOptionalMemberExpressionTree constructs a OptionalMemberExpressionTree.
//
// Java: OptionalMemberExpressionTree(SourceRange, ParseTree, IdentifierToken, boolean)
func NewOptionalMemberExpressionTree(location SourceRange, operand ParseTree, memberName Token, isStartOfOptionalChain bool) *OptionalMemberExpressionTree {
	t := &OptionalMemberExpressionTree{}
	t.InitParseTreeBase(ParseTreeType_OPT_CHAIN_MEMBER_EXPRESSION, location)
	t.Operand = operand
	t.MemberName = memberName
	t.IsStartOfOptionalChain = isStartOfOptionalChain
	return t
}

// Creates a OptionalMemberExpressionTree that represents a getprop expression within an optional
// chain. Has an isOptionalChainStart field to indicate whether it is the start of an optional
// chain.
type OptionalMemberLookupExpressionTree struct {
	ParseTreeBase
	// Go-only: performance — before the other fields, in the base's padding.
	IsStartOfOptionalChain bool

	Operand          ParseTree
	MemberExpression ParseTree
}

// NewOptionalMemberLookupExpressionTree constructs a OptionalMemberLookupExpressionTree.
//
// Java: OptionalMemberLookupExpressionTree(SourceRange, ParseTree, ParseTree, boolean)
func NewOptionalMemberLookupExpressionTree(location SourceRange, operand ParseTree, memberExpression ParseTree, isStartOfOptionalChain bool) *OptionalMemberLookupExpressionTree {
	t := &OptionalMemberLookupExpressionTree{}
	t.InitParseTreeBase(ParseTreeType_OPT_CHAIN_MEMBER_LOOKUP_EXPRESSION, location)
	t.Operand = operand
	t.MemberExpression = memberExpression
	t.IsStartOfOptionalChain = isStartOfOptionalChain
	return t
}

// ParenExpressionTree is the parse tree for ParseTreeType PAREN_EXPRESSION.
type ParenExpressionTree struct {
	ParseTreeBase

	Expression ParseTree
}

// NewParenExpressionTree constructs a ParenExpressionTree.
//
// Java: ParenExpressionTree(SourceRange, ParseTree)
func NewParenExpressionTree(location SourceRange, expression ParseTree) *ParenExpressionTree {
	t := &ParenExpressionTree{}
	t.InitParseTreeBase(ParseTreeType_PAREN_EXPRESSION, location)
	t.Expression = expression
	return t
}

// IsPattern looks through the parentheses.
//
// Java: ParseTree.isPattern()
func (t *ParenExpressionTree) IsPattern() bool {
	return t.Expression.IsPattern()
}

// IsValidNonVanillaAssignmentTarget looks through the parentheses.
//
// Java: ParseTree.isValidNonVanillaAssignmentTarget()
func (t *ParenExpressionTree) IsValidNonVanillaAssignmentTarget() bool {
	return t.Expression.IsValidNonVanillaAssignmentTarget()
}

// IsValidAssignmentTarget looks through the parentheses.
//
// Java: ParseTree.isValidAssignmentTarget()
func (t *ParenExpressionTree) IsValidAssignmentTarget() bool {
	return t.Expression.IsValidAssignmentTarget()
}

// PropertyNameAssignmentTree is the parse tree for ParseTreeType PROPERTY_NAME_ASSIGNMENT.
type PropertyNameAssignmentTree struct {
	ParseTreeBase

	Name Token
	// May be null in ES6 and above.
	Value ParseTree
}

// NewPropertyNameAssignmentTree constructs a PropertyNameAssignmentTree.
//
// Java: PropertyNameAssignmentTree(SourceRange, Token, ParseTree)
func NewPropertyNameAssignmentTree(location SourceRange, name Token, value ParseTree) *PropertyNameAssignmentTree {
	t := &PropertyNameAssignmentTree{}
	t.InitParseTreeBase(ParseTreeType_PROPERTY_NAME_ASSIGNMENT, location)
	t.Name = name
	t.Value = value
	return t
}

// SuperExpressionTree is the parse tree for ParseTreeType SUPER_EXPRESSION.
type SuperExpressionTree struct {
	ParseTreeBase
}

// NewSuperExpressionTree constructs a SuperExpressionTree.
//
// Java: SuperExpressionTree(SourceRange)
func NewSuperExpressionTree(location SourceRange) *SuperExpressionTree {
	t := &SuperExpressionTree{}
	t.InitParseTreeBase(ParseTreeType_SUPER_EXPRESSION, location)
	return t
}

// Template literal production in ES6.
type TemplateLiteralExpressionTree struct {
	ParseTreeBase

	Operand  ParseTree
	Elements []ParseTree
}

// NewTemplateLiteralExpressionTree constructs a TemplateLiteralExpressionTree.
//
// Java: TemplateLiteralExpressionTree(SourceRange, ParseTree, ImmutableList)
func NewTemplateLiteralExpressionTree(location SourceRange, operand ParseTree, elements []ParseTree) *TemplateLiteralExpressionTree {
	t := &TemplateLiteralExpressionTree{}
	t.InitParseTreeBase(ParseTreeType_TEMPLATE_LITERAL_EXPRESSION, location)
	t.Operand = operand
	t.Elements = elements
	return t
}

// A production representing a literal portion of a template literal.
// Eg: `hello`.
type TemplateLiteralPortionTree struct {
	ParseTreeBase

	Value Token
}

// NewTemplateLiteralPortionTree constructs a TemplateLiteralPortionTree.
//
// Java: TemplateLiteralPortionTree(SourceRange, Token)
func NewTemplateLiteralPortionTree(location SourceRange, value Token) *TemplateLiteralPortionTree {
	t := &TemplateLiteralPortionTree{}
	t.InitParseTreeBase(ParseTreeType_TEMPLATE_LITERAL_PORTION, location)
	t.Value = value
	return t
}

// A production representing the expression to be evaluated and substituted
// into a template literal. Eg: ${hello}.
type TemplateSubstitutionTree struct {
	ParseTreeBase

	Expression ParseTree
}

// NewTemplateSubstitutionTree constructs a TemplateSubstitutionTree.
//
// Java: TemplateSubstitutionTree(SourceRange, ParseTree)
func NewTemplateSubstitutionTree(location SourceRange, expression ParseTree) *TemplateSubstitutionTree {
	t := &TemplateSubstitutionTree{}
	t.InitParseTreeBase(ParseTreeType_TEMPLATE_SUBSTITUTION, location)
	t.Expression = expression
	return t
}

// ThisExpressionTree is the parse tree for ParseTreeType THIS_EXPRESSION.
type ThisExpressionTree struct {
	ParseTreeBase
}

// NewThisExpressionTree constructs a ThisExpressionTree.
//
// Java: ThisExpressionTree(SourceRange)
func NewThisExpressionTree(location SourceRange) *ThisExpressionTree {
	t := &ThisExpressionTree{}
	t.InitParseTreeBase(ParseTreeType_THIS_EXPRESSION, location)
	return t
}

// UnaryExpressionTree is the parse tree for ParseTreeType UNARY_EXPRESSION.
type UnaryExpressionTree struct {
	ParseTreeBase

	Operator Token
	Operand  ParseTree
}

// NewUnaryExpressionTree constructs a UnaryExpressionTree.
//
// Java: UnaryExpressionTree(SourceRange, Token, ParseTree)
func NewUnaryExpressionTree(location SourceRange, operator Token, operand ParseTree) *UnaryExpressionTree {
	t := &UnaryExpressionTree{}
	t.InitParseTreeBase(ParseTreeType_UNARY_EXPRESSION, location)
	t.Operator = operator
	t.Operand = operand
	return t
}

// Represents UpdateExpression productions from the spec.
//
// UpdateExpression :=
//
//	{ ++ | -- } UnaryExpression
//	LeftHandSideExpression [no LineTerminator here] { ++ | -- }
type UpdateExpressionTree struct {
	ParseTreeBase
	// Go-only: performance — before the other fields, in the base's padding.
	OperatorPosition UpdateExpressionTree_OperatorPosition

	Operator Token
	Operand  ParseTree
}

// UpdateExpressionTree_OperatorPosition is UpdateExpressionTree.OperatorPosition:
// position of the operator relative to the operand.
type UpdateExpressionTree_OperatorPosition int32

const (
	UpdateExpressionTree_OperatorPosition_PREFIX UpdateExpressionTree_OperatorPosition = iota
	UpdateExpressionTree_OperatorPosition_POSTFIX
)

var updateExpressionTree_OperatorPosition_names = [...]string{"PREFIX", "POSTFIX"}

// UpdateExpressionTree_OperatorPosition_values is OperatorPosition.values().
//
// Java: UpdateExpressionTree.OperatorPosition.values()
func UpdateExpressionTree_OperatorPosition_values() []UpdateExpressionTree_OperatorPosition {
	return []UpdateExpressionTree_OperatorPosition{UpdateExpressionTree_OperatorPosition_PREFIX, UpdateExpressionTree_OperatorPosition_POSTFIX}
}

// Name is Enum.name().
//
// Java: UpdateExpressionTree.OperatorPosition.name()
func (p UpdateExpressionTree_OperatorPosition) Name() string {
	return updateExpressionTree_OperatorPosition_names[p]
}

// String is Enum.toString().
//
// Java: UpdateExpressionTree.OperatorPosition.toString()
func (p UpdateExpressionTree_OperatorPosition) String() string { return p.Name() }

// NewUpdateExpressionTree constructs a UpdateExpressionTree.
//
// Java: UpdateExpressionTree(SourceRange, Token, OperatorPosition, ParseTree)
func NewUpdateExpressionTree(location SourceRange, operator Token, operatorPosition UpdateExpressionTree_OperatorPosition, operand ParseTree) *UpdateExpressionTree {
	t := &UpdateExpressionTree{}
	t.InitParseTreeBase(ParseTreeType_UPDATE_EXPRESSION, location)
	t.Operator = operator
	t.OperatorPosition = operatorPosition
	t.Operand = operand
	return t
}

// UpdateExpressionTree_Prefix creates a prefix update expression.
//
// Java: UpdateExpressionTree.prefix(SourceRange, Token, ParseTree)
func UpdateExpressionTree_Prefix(location SourceRange, operator Token, operand ParseTree) *UpdateExpressionTree {
	return NewUpdateExpressionTree(location, operator, UpdateExpressionTree_OperatorPosition_PREFIX, operand)
}

// UpdateExpressionTree_Postfix creates a postfix update expression.
//
// Java: UpdateExpressionTree.postfix(SourceRange, Token, ParseTree)
func UpdateExpressionTree_Postfix(location SourceRange, operator Token, operand ParseTree) *UpdateExpressionTree {
	return NewUpdateExpressionTree(location, operator, UpdateExpressionTree_OperatorPosition_POSTFIX, operand)
}

// YieldExpressionTree is the parse tree for ParseTreeType YIELD_EXPRESSION.
type YieldExpressionTree struct {
	ParseTreeBase
	// Go-only: performance — before the other fields, in the base's padding.
	IsYieldAll bool

	Expression ParseTree
}

// NewYieldExpressionTree constructs a YieldExpressionTree.
//
// Java: YieldExpressionTree(SourceRange, boolean, ParseTree)
func NewYieldExpressionTree(location SourceRange, isYieldAll bool, expression ParseTree) *YieldExpressionTree {
	t := &YieldExpressionTree{}
	t.InitParseTreeBase(ParseTreeType_YIELD_EXPRESSION, location)
	t.IsYieldAll = isYieldAll
	t.Expression = expression
	return t
}
