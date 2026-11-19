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
// Port of com/google/javascript/jscomp/parsing/parser/trees/ParseTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java polymorphic class interface: com.google.javascript.jscomp.parsing.parser.trees.ParseTree
type ParseTree interface {
	Implements_parser_ParseTree()
	IsJavaNil_() bool
	GetParseTreeBase() *ParseTreeBase
	FieldGet_type_() ParseTreeType
	FieldGet_location() SourceRange
	GetStart() SourcePosition
	GetEnd() SourcePosition
	AsArrayLiteralExpression() *ArrayLiteralExpressionTree
	AsArrayPattern() *ArrayPatternTree
	AsBinaryOperator() *BinaryOperatorTree
	AsBlock() *BlockTree
	AsBreakStatement() *BreakStatementTree
	AsCallExpression() *CallExpressionTree
	AsOptChainCallExpression() *OptChainCallExpressionTree
	AsCaseClause() *CaseClauseTree
	AsCatch() *CatchTree
	AsClassDeclaration() *ClassDeclarationTree
	AsCommaExpression() *CommaExpressionTree
	AsComputedPropertyDefinition() *ComputedPropertyDefinitionTree
	AsComputedPropertyField() *ComputedPropertyFieldTree
	AsComputedPropertyGetter() *ComputedPropertyGetterTree
	AsComputedPropertyMethod() *ComputedPropertyMethodTree
	AsComputedPropertySetter() *ComputedPropertySetterTree
	AsConditionalExpression() *ConditionalExpressionTree
	AsContinueStatement() *ContinueStatementTree
	AsDebuggerStatement() *DebuggerStatementTree
	AsDefaultClause() *DefaultClauseTree
	AsDefaultParameter() *DefaultParameterTree
	AsDoWhileStatement() *DoWhileStatementTree
	AsEmptyStatement() *EmptyStatementTree
	AsExportDeclaration() *ExportDeclarationTree
	AsExportSpecifier() *ExportSpecifierTree
	AsExpressionStatement() *ExpressionStatementTree
	AsFinally() *FinallyTree
	AsForOfStatement() *ForOfStatementTree
	AsForInStatement() *ForInStatementTree
	AsFormalParameterList() *FormalParameterListTree
	AsForStatement() *ForStatementTree
	AsFunctionDeclaration() *FunctionDeclarationTree
	AsGetAccessor() *GetAccessorTree
	AsIdentifierExpression() *IdentifierExpressionTree
	AsIfStatement() *IfStatementTree
	AsImportDeclaration() *ImportDeclarationTree
	AsImportSpecifier() *ImportSpecifierTree
	AsDynamicImportExpression() *DynamicImportTree
	AsImportMetaExpression() *ImportMetaExpressionTree
	AsLabelledStatement() *LabelledStatementTree
	AsLiteralExpression() *LiteralExpressionTree
	AsMemberExpression() *MemberExpressionTree
	AsOptionalMemberExpression() *OptionalMemberExpressionTree
	AsMemberLookupExpression() *MemberLookupExpressionTree
	AsOptionalMemberLookupExpression() *OptionalMemberLookupExpressionTree
	AsMissingPrimaryExpression() *MissingPrimaryExpressionTree
	AsNewExpression() *NewExpressionTree
	AsNull() *NullTree
	AsObjectLiteralExpression() *ObjectLiteralExpressionTree
	AsObjectPattern() *ObjectPatternTree
	AsParenExpression() *ParenExpressionTree
	AsProgram() *ProgramTree
	AsPropertyNameAssignment() *PropertyNameAssignmentTree
	AsIterRest() *IterRestTree
	AsObjectRest() *ObjectRestTree
	AsReturnStatement() *ReturnStatementTree
	AsSetAccessor() *SetAccessorTree
	AsIterSpread() *IterSpreadTree
	AsObjectSpread() *ObjectSpreadTree
	AsSuperExpression() *SuperExpressionTree
	AsSwitchStatement() *SwitchStatementTree
	AsTemplateLiteralExpression() *TemplateLiteralExpressionTree
	AsTemplateLiteralPortion() *TemplateLiteralPortionTree
	AsTemplateSubstitution() *TemplateSubstitutionTree
	AsThisExpression() *ThisExpressionTree
	AsThrowStatement() *ThrowStatementTree
	AsTryStatement() *TryStatementTree
	AsUnaryExpression() *UnaryExpressionTree
	AsVariableDeclarationList() *VariableDeclarationListTree
	AsVariableDeclaration() *VariableDeclarationTree
	AsVariableStatement() *VariableStatementTree
	AsWhileStatement() *WhileStatementTree
	AsWithStatement() *WithStatementTree
	AsYieldStatement() *YieldExpressionTree
	AsAwaitExpression() *AwaitExpressionTree
	AsNewTargetExpression() *NewTargetExpressionTree
	AsUpdateExpression() *UpdateExpressionTree
	AsForAwaitOfStatement() *ForAwaitOfStatementTree
	AsFieldDeclaration() *FieldDeclarationTree
	IsPattern() bool
	IsValidNonVanillaAssignmentTarget() bool
	IsValidAssignmentTarget() bool
	IsRestParameter() bool
	String() string
}

func ToIface_ParseTree[T interface {
	ParseTree
	comparable
}](v T) ParseTree {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_ParseTree[T ParseTree](v T) ParseTree {
	var vi ParseTree = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.jscomp.parsing.parser.trees.ParseTree
type ParseTreeBase struct {
	self_    ParseTree
	location SourceRange   `gson:"location"`
	type_    ParseTreeType `gson:"type"`
}

func (this_ *ParseTreeBase) BindSelf_(s any) {
	_ = s
	this_.self_ = s.(ParseTree)
}
func (this_ *ParseTreeBase) IsJavaNil_() bool { return this_ == nil }
func (*ParseTreeBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.ParseTree", "ParseTree"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.ParseTree", "com.google.javascript.jscomp.parsing.parser.trees.ParseTree", func(v_ any) bool {
		return j2gort.InstanceOf[ParseTree](v_)
	})
}
func (b *ParseTreeBase) GetParseTreeBase() *ParseTreeBase   { return b }
func (*ParseTreeBase) Implements_parser_ParseTree()         {}
func (this_ *ParseTreeBase) FieldGet_type_() ParseTreeType  { return this_.type_ }
func (this_ *ParseTreeBase) FieldGet_location() SourceRange { return this_.location }

func (this_ *ParseTreeBase) initFields_ParseTree() {
	_ = this_
}

// Java constructor: ParseTree(com.google.javascript.jscomp.parsing.parser.trees.ParseTreeType type, com.google.javascript.jscomp.parsing.parser.util.SourceRange location)
func NewParseTree(type_ ParseTreeType, location SourceRange) ParseTree {
	this_ := &ParseTreeBase{}
	this_.BindSelf_(this_)
	InitParseTree(this_, type_, location)
	return this_
}

func InitParseTree(this_ *ParseTreeBase, type_ ParseTreeType, location SourceRange) {
	_ = type_
	_ = location
	this_.initFields_ParseTree()
	this_.type_ = type_
	this_.location = location
}

// Java method: ParseTree.getStart()
func (this_ *ParseTreeBase) GetStart() SourcePosition {
	return this_.location.start
}

// Java method: ParseTree.getEnd()
func (this_ *ParseTreeBase) GetEnd() SourcePosition {
	return this_.location.end
}

// Java method: ParseTree.asArrayLiteralExpression()
func (this_ *ParseTreeBase) AsArrayLiteralExpression() *ArrayLiteralExpressionTree {
	return j2gort.CastPtr[ArrayLiteralExpressionTree](this_.self_)
}

// Java method: ParseTree.asArrayPattern()
func (this_ *ParseTreeBase) AsArrayPattern() *ArrayPatternTree {
	return j2gort.CastPtr[ArrayPatternTree](this_.self_)
}

// Java method: ParseTree.asBinaryOperator()
func (this_ *ParseTreeBase) AsBinaryOperator() *BinaryOperatorTree {
	return j2gort.CastPtr[BinaryOperatorTree](this_.self_)
}

// Java method: ParseTree.asBlock()
func (this_ *ParseTreeBase) AsBlock() *BlockTree {
	return j2gort.CastPtr[BlockTree](this_.self_)
}

// Java method: ParseTree.asBreakStatement()
func (this_ *ParseTreeBase) AsBreakStatement() *BreakStatementTree {
	return j2gort.CastPtr[BreakStatementTree](this_.self_)
}

// Java method: ParseTree.asCallExpression()
func (this_ *ParseTreeBase) AsCallExpression() *CallExpressionTree {
	return j2gort.CastPtr[CallExpressionTree](this_.self_)
}

// Java method: ParseTree.asOptChainCallExpression()
func (this_ *ParseTreeBase) AsOptChainCallExpression() *OptChainCallExpressionTree {
	return j2gort.CastPtr[OptChainCallExpressionTree](this_.self_)
}

// Java method: ParseTree.asCaseClause()
func (this_ *ParseTreeBase) AsCaseClause() *CaseClauseTree {
	return j2gort.CastPtr[CaseClauseTree](this_.self_)
}

// Java method: ParseTree.asCatch()
func (this_ *ParseTreeBase) AsCatch() *CatchTree {
	return j2gort.CastPtr[CatchTree](this_.self_)
}

// Java method: ParseTree.asClassDeclaration()
func (this_ *ParseTreeBase) AsClassDeclaration() *ClassDeclarationTree {
	return j2gort.CastPtr[ClassDeclarationTree](this_.self_)
}

// Java method: ParseTree.asCommaExpression()
func (this_ *ParseTreeBase) AsCommaExpression() *CommaExpressionTree {
	return j2gort.CastPtr[CommaExpressionTree](this_.self_)
}

// Java method: ParseTree.asComputedPropertyDefinition()
func (this_ *ParseTreeBase) AsComputedPropertyDefinition() *ComputedPropertyDefinitionTree {
	return j2gort.CastPtr[ComputedPropertyDefinitionTree](this_.self_)
}

// Java method: ParseTree.asComputedPropertyField()
func (this_ *ParseTreeBase) AsComputedPropertyField() *ComputedPropertyFieldTree {
	return j2gort.CastPtr[ComputedPropertyFieldTree](this_.self_)
}

// Java method: ParseTree.asComputedPropertyGetter()
func (this_ *ParseTreeBase) AsComputedPropertyGetter() *ComputedPropertyGetterTree {
	return j2gort.CastPtr[ComputedPropertyGetterTree](this_.self_)
}

// Java method: ParseTree.asComputedPropertyMethod()
func (this_ *ParseTreeBase) AsComputedPropertyMethod() *ComputedPropertyMethodTree {
	return j2gort.CastPtr[ComputedPropertyMethodTree](this_.self_)
}

// Java method: ParseTree.asComputedPropertySetter()
func (this_ *ParseTreeBase) AsComputedPropertySetter() *ComputedPropertySetterTree {
	return j2gort.CastPtr[ComputedPropertySetterTree](this_.self_)
}

// Java method: ParseTree.asConditionalExpression()
func (this_ *ParseTreeBase) AsConditionalExpression() *ConditionalExpressionTree {
	return j2gort.CastPtr[ConditionalExpressionTree](this_.self_)
}

// Java method: ParseTree.asContinueStatement()
func (this_ *ParseTreeBase) AsContinueStatement() *ContinueStatementTree {
	return j2gort.CastPtr[ContinueStatementTree](this_.self_)
}

// Java method: ParseTree.asDebuggerStatement()
func (this_ *ParseTreeBase) AsDebuggerStatement() *DebuggerStatementTree {
	return j2gort.CastPtr[DebuggerStatementTree](this_.self_)
}

// Java method: ParseTree.asDefaultClause()
func (this_ *ParseTreeBase) AsDefaultClause() *DefaultClauseTree {
	return j2gort.CastPtr[DefaultClauseTree](this_.self_)
}

// Java method: ParseTree.asDefaultParameter()
func (this_ *ParseTreeBase) AsDefaultParameter() *DefaultParameterTree {
	return j2gort.CastPtr[DefaultParameterTree](this_.self_)
}

// Java method: ParseTree.asDoWhileStatement()
func (this_ *ParseTreeBase) AsDoWhileStatement() *DoWhileStatementTree {
	return j2gort.CastPtr[DoWhileStatementTree](this_.self_)
}

// Java method: ParseTree.asEmptyStatement()
func (this_ *ParseTreeBase) AsEmptyStatement() *EmptyStatementTree {
	return j2gort.CastPtr[EmptyStatementTree](this_.self_)
}

// Java method: ParseTree.asExportDeclaration()
func (this_ *ParseTreeBase) AsExportDeclaration() *ExportDeclarationTree {
	return j2gort.CastPtr[ExportDeclarationTree](this_.self_)
}

// Java method: ParseTree.asExportSpecifier()
func (this_ *ParseTreeBase) AsExportSpecifier() *ExportSpecifierTree {
	return j2gort.CastPtr[ExportSpecifierTree](this_.self_)
}

// Java method: ParseTree.asExpressionStatement()
func (this_ *ParseTreeBase) AsExpressionStatement() *ExpressionStatementTree {
	return j2gort.CastPtr[ExpressionStatementTree](this_.self_)
}

// Java method: ParseTree.asFinally()
func (this_ *ParseTreeBase) AsFinally() *FinallyTree {
	return j2gort.CastPtr[FinallyTree](this_.self_)
}

// Java method: ParseTree.asForOfStatement()
func (this_ *ParseTreeBase) AsForOfStatement() *ForOfStatementTree {
	return j2gort.CastPtr[ForOfStatementTree](this_.self_)
}

// Java method: ParseTree.asForInStatement()
func (this_ *ParseTreeBase) AsForInStatement() *ForInStatementTree {
	return j2gort.CastPtr[ForInStatementTree](this_.self_)
}

// Java method: ParseTree.asFormalParameterList()
func (this_ *ParseTreeBase) AsFormalParameterList() *FormalParameterListTree {
	return j2gort.CastPtr[FormalParameterListTree](this_.self_)
}

// Java method: ParseTree.asForStatement()
func (this_ *ParseTreeBase) AsForStatement() *ForStatementTree {
	return j2gort.CastPtr[ForStatementTree](this_.self_)
}

// Java method: ParseTree.asFunctionDeclaration()
func (this_ *ParseTreeBase) AsFunctionDeclaration() *FunctionDeclarationTree {
	return j2gort.CastPtr[FunctionDeclarationTree](this_.self_)
}

// Java method: ParseTree.asGetAccessor()
func (this_ *ParseTreeBase) AsGetAccessor() *GetAccessorTree {
	return j2gort.CastPtr[GetAccessorTree](this_.self_)
}

// Java method: ParseTree.asIdentifierExpression()
func (this_ *ParseTreeBase) AsIdentifierExpression() *IdentifierExpressionTree {
	return j2gort.CastPtr[IdentifierExpressionTree](this_.self_)
}

// Java method: ParseTree.asIfStatement()
func (this_ *ParseTreeBase) AsIfStatement() *IfStatementTree {
	return j2gort.CastPtr[IfStatementTree](this_.self_)
}

// Java method: ParseTree.asImportDeclaration()
func (this_ *ParseTreeBase) AsImportDeclaration() *ImportDeclarationTree {
	return j2gort.CastPtr[ImportDeclarationTree](this_.self_)
}

// Java method: ParseTree.asImportSpecifier()
func (this_ *ParseTreeBase) AsImportSpecifier() *ImportSpecifierTree {
	return j2gort.CastPtr[ImportSpecifierTree](this_.self_)
}

// Java method: ParseTree.asDynamicImportExpression()
func (this_ *ParseTreeBase) AsDynamicImportExpression() *DynamicImportTree {
	return j2gort.CastPtr[DynamicImportTree](this_.self_)
}

// Java method: ParseTree.asImportMetaExpression()
func (this_ *ParseTreeBase) AsImportMetaExpression() *ImportMetaExpressionTree {
	return j2gort.CastPtr[ImportMetaExpressionTree](this_.self_)
}

// Java method: ParseTree.asLabelledStatement()
func (this_ *ParseTreeBase) AsLabelledStatement() *LabelledStatementTree {
	return j2gort.CastPtr[LabelledStatementTree](this_.self_)
}

// Java method: ParseTree.asLiteralExpression()
func (this_ *ParseTreeBase) AsLiteralExpression() *LiteralExpressionTree {
	return j2gort.CastPtr[LiteralExpressionTree](this_.self_)
}

// Java method: ParseTree.asMemberExpression()
func (this_ *ParseTreeBase) AsMemberExpression() *MemberExpressionTree {
	return j2gort.CastPtr[MemberExpressionTree](this_.self_)
}

// Java method: ParseTree.asOptionalMemberExpression()
func (this_ *ParseTreeBase) AsOptionalMemberExpression() *OptionalMemberExpressionTree {
	return j2gort.CastPtr[OptionalMemberExpressionTree](this_.self_)
}

// Java method: ParseTree.asMemberLookupExpression()
func (this_ *ParseTreeBase) AsMemberLookupExpression() *MemberLookupExpressionTree {
	return j2gort.CastPtr[MemberLookupExpressionTree](this_.self_)
}

// Java method: ParseTree.asOptionalMemberLookupExpression()
func (this_ *ParseTreeBase) AsOptionalMemberLookupExpression() *OptionalMemberLookupExpressionTree {
	return j2gort.CastPtr[OptionalMemberLookupExpressionTree](this_.self_)
}

// Java method: ParseTree.asMissingPrimaryExpression()
func (this_ *ParseTreeBase) AsMissingPrimaryExpression() *MissingPrimaryExpressionTree {
	return j2gort.CastPtr[MissingPrimaryExpressionTree](this_.self_)
}

// Java method: ParseTree.asNewExpression()
func (this_ *ParseTreeBase) AsNewExpression() *NewExpressionTree {
	return j2gort.CastPtr[NewExpressionTree](this_.self_)
}

// Java method: ParseTree.asNull()
func (this_ *ParseTreeBase) AsNull() *NullTree {
	return j2gort.CastPtr[NullTree](this_.self_)
}

// Java method: ParseTree.asObjectLiteralExpression()
func (this_ *ParseTreeBase) AsObjectLiteralExpression() *ObjectLiteralExpressionTree {
	return j2gort.CastPtr[ObjectLiteralExpressionTree](this_.self_)
}

// Java method: ParseTree.asObjectPattern()
func (this_ *ParseTreeBase) AsObjectPattern() *ObjectPatternTree {
	return j2gort.CastPtr[ObjectPatternTree](this_.self_)
}

// Java method: ParseTree.asParenExpression()
func (this_ *ParseTreeBase) AsParenExpression() *ParenExpressionTree {
	return j2gort.CastPtr[ParenExpressionTree](this_.self_)
}

// Java method: ParseTree.asProgram()
func (this_ *ParseTreeBase) AsProgram() *ProgramTree {
	return j2gort.CastPtr[ProgramTree](this_.self_)
}

// Java method: ParseTree.asPropertyNameAssignment()
func (this_ *ParseTreeBase) AsPropertyNameAssignment() *PropertyNameAssignmentTree {
	return j2gort.CastPtr[PropertyNameAssignmentTree](this_.self_)
}

// Java method: ParseTree.asIterRest()
func (this_ *ParseTreeBase) AsIterRest() *IterRestTree {
	return j2gort.CastPtr[IterRestTree](this_.self_)
}

// Java method: ParseTree.asObjectRest()
func (this_ *ParseTreeBase) AsObjectRest() *ObjectRestTree {
	return j2gort.CastPtr[ObjectRestTree](this_.self_)
}

// Java method: ParseTree.asReturnStatement()
func (this_ *ParseTreeBase) AsReturnStatement() *ReturnStatementTree {
	return j2gort.CastPtr[ReturnStatementTree](this_.self_)
}

// Java method: ParseTree.asSetAccessor()
func (this_ *ParseTreeBase) AsSetAccessor() *SetAccessorTree {
	return j2gort.CastPtr[SetAccessorTree](this_.self_)
}

// Java method: ParseTree.asIterSpread()
func (this_ *ParseTreeBase) AsIterSpread() *IterSpreadTree {
	return j2gort.CastPtr[IterSpreadTree](this_.self_)
}

// Java method: ParseTree.asObjectSpread()
func (this_ *ParseTreeBase) AsObjectSpread() *ObjectSpreadTree {
	return j2gort.CastPtr[ObjectSpreadTree](this_.self_)
}

// Java method: ParseTree.asSuperExpression()
func (this_ *ParseTreeBase) AsSuperExpression() *SuperExpressionTree {
	return j2gort.CastPtr[SuperExpressionTree](this_.self_)
}

// Java method: ParseTree.asSwitchStatement()
func (this_ *ParseTreeBase) AsSwitchStatement() *SwitchStatementTree {
	return j2gort.CastPtr[SwitchStatementTree](this_.self_)
}

// Java method: ParseTree.asTemplateLiteralExpression()
func (this_ *ParseTreeBase) AsTemplateLiteralExpression() *TemplateLiteralExpressionTree {
	return j2gort.CastPtr[TemplateLiteralExpressionTree](this_.self_)
}

// Java method: ParseTree.asTemplateLiteralPortion()
func (this_ *ParseTreeBase) AsTemplateLiteralPortion() *TemplateLiteralPortionTree {
	return j2gort.CastPtr[TemplateLiteralPortionTree](this_.self_)
}

// Java method: ParseTree.asTemplateSubstitution()
func (this_ *ParseTreeBase) AsTemplateSubstitution() *TemplateSubstitutionTree {
	return j2gort.CastPtr[TemplateSubstitutionTree](this_.self_)
}

// Java method: ParseTree.asThisExpression()
func (this_ *ParseTreeBase) AsThisExpression() *ThisExpressionTree {
	return j2gort.CastPtr[ThisExpressionTree](this_.self_)
}

// Java method: ParseTree.asThrowStatement()
func (this_ *ParseTreeBase) AsThrowStatement() *ThrowStatementTree {
	return j2gort.CastPtr[ThrowStatementTree](this_.self_)
}

// Java method: ParseTree.asTryStatement()
func (this_ *ParseTreeBase) AsTryStatement() *TryStatementTree {
	return j2gort.CastPtr[TryStatementTree](this_.self_)
}

// Java method: ParseTree.asUnaryExpression()
func (this_ *ParseTreeBase) AsUnaryExpression() *UnaryExpressionTree {
	return j2gort.CastPtr[UnaryExpressionTree](this_.self_)
}

// Java method: ParseTree.asVariableDeclarationList()
func (this_ *ParseTreeBase) AsVariableDeclarationList() *VariableDeclarationListTree {
	return j2gort.CastPtr[VariableDeclarationListTree](this_.self_)
}

// Java method: ParseTree.asVariableDeclaration()
func (this_ *ParseTreeBase) AsVariableDeclaration() *VariableDeclarationTree {
	return j2gort.CastPtr[VariableDeclarationTree](this_.self_)
}

// Java method: ParseTree.asVariableStatement()
func (this_ *ParseTreeBase) AsVariableStatement() *VariableStatementTree {
	return j2gort.CastPtr[VariableStatementTree](this_.self_)
}

// Java method: ParseTree.asWhileStatement()
func (this_ *ParseTreeBase) AsWhileStatement() *WhileStatementTree {
	return j2gort.CastPtr[WhileStatementTree](this_.self_)
}

// Java method: ParseTree.asWithStatement()
func (this_ *ParseTreeBase) AsWithStatement() *WithStatementTree {
	return j2gort.CastPtr[WithStatementTree](this_.self_)
}

// Java method: ParseTree.asYieldStatement()
func (this_ *ParseTreeBase) AsYieldStatement() *YieldExpressionTree {
	return j2gort.CastPtr[YieldExpressionTree](this_.self_)
}

// Java method: ParseTree.asAwaitExpression()
func (this_ *ParseTreeBase) AsAwaitExpression() *AwaitExpressionTree {
	return j2gort.CastPtr[AwaitExpressionTree](this_.self_)
}

// Java method: ParseTree.asNewTargetExpression()
func (this_ *ParseTreeBase) AsNewTargetExpression() *NewTargetExpressionTree {
	return j2gort.CastPtr[NewTargetExpressionTree](this_.self_)
}

// Java method: ParseTree.asUpdateExpression()
func (this_ *ParseTreeBase) AsUpdateExpression() *UpdateExpressionTree {
	return j2gort.CastPtr[UpdateExpressionTree](this_.self_)
}

// Java method: ParseTree.asForAwaitOfStatement()
func (this_ *ParseTreeBase) AsForAwaitOfStatement() *ForAwaitOfStatementTree {
	return j2gort.CastPtr[ForAwaitOfStatementTree](this_.self_)
}

// Java method: ParseTree.asFieldDeclaration()
func (this_ *ParseTreeBase) AsFieldDeclaration() *FieldDeclarationTree {
	return j2gort.CastPtr[FieldDeclarationTree](this_.self_)
}

// Java method: ParseTree.isPattern()
func (this_ *ParseTreeBase) IsPattern() bool {
	var parseTree ParseTree = this_.self_
	_ = parseTree
	for parseTree.GetParseTreeBase().type_ == StaticGet_ParseTreeType_PAREN_EXPRESSION() {
		parseTree = parseTree.GetParseTreeBase().AsParenExpression().expression
	}
	switch (parseTree.GetParseTreeBase().type_).Ordinal() {
	case int32(55) /* ParseTreeType.ARRAY_PATTERN */, int32(56) /* ParseTreeType.OBJECT_PATTERN */ :
		return true
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java method: ParseTree.isValidNonVanillaAssignmentTarget()
func (this_ *ParseTreeBase) IsValidNonVanillaAssignmentTarget() bool {
	var parseTree ParseTree = this_.self_
	_ = parseTree
	for parseTree.GetParseTreeBase().type_ == StaticGet_ParseTreeType_PAREN_EXPRESSION() {
		parseTree = parseTree.GetParseTreeBase().AsParenExpression().expression
	}
	switch (parseTree.GetParseTreeBase().type_).Ordinal() {
	case int32(26) /* ParseTreeType.IDENTIFIER_EXPRESSION */, int32(43) /* ParseTreeType.MEMBER_EXPRESSION */, int32(50) /* ParseTreeType.MEMBER_LOOKUP_EXPRESSION */, int32(71) /* ParseTreeType.DEFAULT_PARAMETER */ :
		return true
	case int32(55) /* ParseTreeType.ARRAY_PATTERN */, int32(56) /* ParseTreeType.OBJECT_PATTERN */ :
		return false
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java method: ParseTree.isValidAssignmentTarget()
func (this_ *ParseTreeBase) IsValidAssignmentTarget() bool {
	var parseTree ParseTree = this_.self_
	_ = parseTree
	for parseTree.GetParseTreeBase().type_ == StaticGet_ParseTreeType_PAREN_EXPRESSION() {
		parseTree = parseTree.GetParseTreeBase().AsParenExpression().expression
	}
	switch (parseTree.GetParseTreeBase().type_).Ordinal() {
	case int32(26) /* ParseTreeType.IDENTIFIER_EXPRESSION */, int32(43) /* ParseTreeType.MEMBER_EXPRESSION */, int32(50) /* ParseTreeType.MEMBER_LOOKUP_EXPRESSION */, int32(55) /* ParseTreeType.ARRAY_PATTERN */, int32(56) /* ParseTreeType.OBJECT_PATTERN */, int32(71) /* ParseTreeType.DEFAULT_PARAMETER */ :
		return true
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java method: ParseTree.isRestParameter()
func (this_ *ParseTreeBase) IsRestParameter() bool {
	return (this_.type_ == StaticGet_ParseTreeType_ITER_REST())
}

// Java method: ParseTree.toString()
func (this_ *ParseTreeBase) String() string {
	return (j2gort.StrVal(this_.type_) + "@" + j2gort.StrVal(this_.location))
}
