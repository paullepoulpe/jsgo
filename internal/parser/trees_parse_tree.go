// Port of com/google/javascript/jscomp/parsing/parser/trees/ParseTree.java @ closure-compiler v20261015

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
// through GetParseTreeBase().
type ParseTree interface {
	GetParseTreeBase() *ParseTreeBase
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
	AsComprehensionIf() *ComprehensionIfTree
	AsComprehensionFor() *ComprehensionForTree
	AsComprehension() *ComprehensionTree
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

// ParseTreeBase holds ParseTree's fields and methods.
type ParseTreeBase struct {
	self ParseTree

	Type     ParseTreeType
	Location SourceRange
}

// InitParseTreeBase is the Go form of the (protected) ParseTree
// constructor; self is the concrete tree.
//
// Java: ParseTree(ParseTreeType, SourceRange)
func (t *ParseTreeBase) InitParseTreeBase(self ParseTree, type_ ParseTreeType, location SourceRange) {
	t.self = self
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

// AsArrayLiteralExpression casts to ArrayLiteralExpressionTree.
//
// Java: ParseTree.asArrayLiteralExpression()
func (t *ParseTreeBase) AsArrayLiteralExpression() *ArrayLiteralExpressionTree {
	return t.self.(*ArrayLiteralExpressionTree)
}

// AsArrayPattern casts to ArrayPatternTree.
//
// Java: ParseTree.asArrayPattern()
func (t *ParseTreeBase) AsArrayPattern() *ArrayPatternTree {
	return t.self.(*ArrayPatternTree)
}

// AsBinaryOperator casts to BinaryOperatorTree.
//
// Java: ParseTree.asBinaryOperator()
func (t *ParseTreeBase) AsBinaryOperator() *BinaryOperatorTree {
	return t.self.(*BinaryOperatorTree)
}

// AsBlock casts to BlockTree.
//
// Java: ParseTree.asBlock()
func (t *ParseTreeBase) AsBlock() *BlockTree {
	return t.self.(*BlockTree)
}

// AsBreakStatement casts to BreakStatementTree.
//
// Java: ParseTree.asBreakStatement()
func (t *ParseTreeBase) AsBreakStatement() *BreakStatementTree {
	return t.self.(*BreakStatementTree)
}

// AsCallExpression casts to CallExpressionTree.
//
// Java: ParseTree.asCallExpression()
func (t *ParseTreeBase) AsCallExpression() *CallExpressionTree {
	return t.self.(*CallExpressionTree)
}

// AsOptChainCallExpression casts to OptChainCallExpressionTree.
//
// Java: ParseTree.asOptChainCallExpression()
func (t *ParseTreeBase) AsOptChainCallExpression() *OptChainCallExpressionTree {
	return t.self.(*OptChainCallExpressionTree)
}

// AsCaseClause casts to CaseClauseTree.
//
// Java: ParseTree.asCaseClause()
func (t *ParseTreeBase) AsCaseClause() *CaseClauseTree {
	return t.self.(*CaseClauseTree)
}

// AsCatch casts to CatchTree.
//
// Java: ParseTree.asCatch()
func (t *ParseTreeBase) AsCatch() *CatchTree {
	return t.self.(*CatchTree)
}

// AsClassDeclaration casts to ClassDeclarationTree.
//
// Java: ParseTree.asClassDeclaration()
func (t *ParseTreeBase) AsClassDeclaration() *ClassDeclarationTree {
	return t.self.(*ClassDeclarationTree)
}

// AsCommaExpression casts to CommaExpressionTree.
//
// Java: ParseTree.asCommaExpression()
func (t *ParseTreeBase) AsCommaExpression() *CommaExpressionTree {
	return t.self.(*CommaExpressionTree)
}

// AsComprehensionIf casts to ComprehensionIfTree.
//
// Java: ParseTree.asComprehensionIf()
func (t *ParseTreeBase) AsComprehensionIf() *ComprehensionIfTree {
	return t.self.(*ComprehensionIfTree)
}

// AsComprehensionFor casts to ComprehensionForTree.
//
// Java: ParseTree.asComprehensionFor()
func (t *ParseTreeBase) AsComprehensionFor() *ComprehensionForTree {
	return t.self.(*ComprehensionForTree)
}

// AsComprehension casts to ComprehensionTree.
//
// Java: ParseTree.asComprehension()
func (t *ParseTreeBase) AsComprehension() *ComprehensionTree {
	return t.self.(*ComprehensionTree)
}

// AsComputedPropertyDefinition casts to ComputedPropertyDefinitionTree.
//
// Java: ParseTree.asComputedPropertyDefinition()
func (t *ParseTreeBase) AsComputedPropertyDefinition() *ComputedPropertyDefinitionTree {
	return t.self.(*ComputedPropertyDefinitionTree)
}

// AsComputedPropertyField casts to ComputedPropertyFieldTree.
//
// Java: ParseTree.asComputedPropertyField()
func (t *ParseTreeBase) AsComputedPropertyField() *ComputedPropertyFieldTree {
	return t.self.(*ComputedPropertyFieldTree)
}

// AsComputedPropertyGetter casts to ComputedPropertyGetterTree.
//
// Java: ParseTree.asComputedPropertyGetter()
func (t *ParseTreeBase) AsComputedPropertyGetter() *ComputedPropertyGetterTree {
	return t.self.(*ComputedPropertyGetterTree)
}

// AsComputedPropertyMethod casts to ComputedPropertyMethodTree.
//
// Java: ParseTree.asComputedPropertyMethod()
func (t *ParseTreeBase) AsComputedPropertyMethod() *ComputedPropertyMethodTree {
	return t.self.(*ComputedPropertyMethodTree)
}

// AsComputedPropertySetter casts to ComputedPropertySetterTree.
//
// Java: ParseTree.asComputedPropertySetter()
func (t *ParseTreeBase) AsComputedPropertySetter() *ComputedPropertySetterTree {
	return t.self.(*ComputedPropertySetterTree)
}

// AsConditionalExpression casts to ConditionalExpressionTree.
//
// Java: ParseTree.asConditionalExpression()
func (t *ParseTreeBase) AsConditionalExpression() *ConditionalExpressionTree {
	return t.self.(*ConditionalExpressionTree)
}

// AsContinueStatement casts to ContinueStatementTree.
//
// Java: ParseTree.asContinueStatement()
func (t *ParseTreeBase) AsContinueStatement() *ContinueStatementTree {
	return t.self.(*ContinueStatementTree)
}

// AsDebuggerStatement casts to DebuggerStatementTree.
//
// Java: ParseTree.asDebuggerStatement()
func (t *ParseTreeBase) AsDebuggerStatement() *DebuggerStatementTree {
	return t.self.(*DebuggerStatementTree)
}

// AsDefaultClause casts to DefaultClauseTree.
//
// Java: ParseTree.asDefaultClause()
func (t *ParseTreeBase) AsDefaultClause() *DefaultClauseTree {
	return t.self.(*DefaultClauseTree)
}

// AsDefaultParameter casts to DefaultParameterTree.
//
// Java: ParseTree.asDefaultParameter()
func (t *ParseTreeBase) AsDefaultParameter() *DefaultParameterTree {
	return t.self.(*DefaultParameterTree)
}

// AsDoWhileStatement casts to DoWhileStatementTree.
//
// Java: ParseTree.asDoWhileStatement()
func (t *ParseTreeBase) AsDoWhileStatement() *DoWhileStatementTree {
	return t.self.(*DoWhileStatementTree)
}

// AsEmptyStatement casts to EmptyStatementTree.
//
// Java: ParseTree.asEmptyStatement()
func (t *ParseTreeBase) AsEmptyStatement() *EmptyStatementTree {
	return t.self.(*EmptyStatementTree)
}

// AsExportDeclaration casts to ExportDeclarationTree.
//
// Java: ParseTree.asExportDeclaration()
func (t *ParseTreeBase) AsExportDeclaration() *ExportDeclarationTree {
	return t.self.(*ExportDeclarationTree)
}

// AsExportSpecifier casts to ExportSpecifierTree.
//
// Java: ParseTree.asExportSpecifier()
func (t *ParseTreeBase) AsExportSpecifier() *ExportSpecifierTree {
	return t.self.(*ExportSpecifierTree)
}

// AsExpressionStatement casts to ExpressionStatementTree.
//
// Java: ParseTree.asExpressionStatement()
func (t *ParseTreeBase) AsExpressionStatement() *ExpressionStatementTree {
	return t.self.(*ExpressionStatementTree)
}

// AsFinally casts to FinallyTree.
//
// Java: ParseTree.asFinally()
func (t *ParseTreeBase) AsFinally() *FinallyTree {
	return t.self.(*FinallyTree)
}

// AsForOfStatement casts to ForOfStatementTree.
//
// Java: ParseTree.asForOfStatement()
func (t *ParseTreeBase) AsForOfStatement() *ForOfStatementTree {
	return t.self.(*ForOfStatementTree)
}

// AsForInStatement casts to ForInStatementTree.
//
// Java: ParseTree.asForInStatement()
func (t *ParseTreeBase) AsForInStatement() *ForInStatementTree {
	return t.self.(*ForInStatementTree)
}

// AsFormalParameterList casts to FormalParameterListTree.
//
// Java: ParseTree.asFormalParameterList()
func (t *ParseTreeBase) AsFormalParameterList() *FormalParameterListTree {
	return t.self.(*FormalParameterListTree)
}

// AsForStatement casts to ForStatementTree.
//
// Java: ParseTree.asForStatement()
func (t *ParseTreeBase) AsForStatement() *ForStatementTree {
	return t.self.(*ForStatementTree)
}

// AsFunctionDeclaration casts to FunctionDeclarationTree.
//
// Java: ParseTree.asFunctionDeclaration()
func (t *ParseTreeBase) AsFunctionDeclaration() *FunctionDeclarationTree {
	return t.self.(*FunctionDeclarationTree)
}

// AsGetAccessor casts to GetAccessorTree.
//
// Java: ParseTree.asGetAccessor()
func (t *ParseTreeBase) AsGetAccessor() *GetAccessorTree {
	return t.self.(*GetAccessorTree)
}

// AsIdentifierExpression casts to IdentifierExpressionTree.
//
// Java: ParseTree.asIdentifierExpression()
func (t *ParseTreeBase) AsIdentifierExpression() *IdentifierExpressionTree {
	return t.self.(*IdentifierExpressionTree)
}

// AsIfStatement casts to IfStatementTree.
//
// Java: ParseTree.asIfStatement()
func (t *ParseTreeBase) AsIfStatement() *IfStatementTree {
	return t.self.(*IfStatementTree)
}

// AsImportDeclaration casts to ImportDeclarationTree.
//
// Java: ParseTree.asImportDeclaration()
func (t *ParseTreeBase) AsImportDeclaration() *ImportDeclarationTree {
	return t.self.(*ImportDeclarationTree)
}

// AsImportSpecifier casts to ImportSpecifierTree.
//
// Java: ParseTree.asImportSpecifier()
func (t *ParseTreeBase) AsImportSpecifier() *ImportSpecifierTree {
	return t.self.(*ImportSpecifierTree)
}

// AsDynamicImportExpression casts to DynamicImportTree.
//
// Java: ParseTree.asDynamicImportExpression()
func (t *ParseTreeBase) AsDynamicImportExpression() *DynamicImportTree {
	return t.self.(*DynamicImportTree)
}

// AsImportMetaExpression casts to ImportMetaExpressionTree.
//
// Java: ParseTree.asImportMetaExpression()
func (t *ParseTreeBase) AsImportMetaExpression() *ImportMetaExpressionTree {
	return t.self.(*ImportMetaExpressionTree)
}

// AsLabelledStatement casts to LabelledStatementTree.
//
// Java: ParseTree.asLabelledStatement()
func (t *ParseTreeBase) AsLabelledStatement() *LabelledStatementTree {
	return t.self.(*LabelledStatementTree)
}

// AsLiteralExpression casts to LiteralExpressionTree.
//
// Java: ParseTree.asLiteralExpression()
func (t *ParseTreeBase) AsLiteralExpression() *LiteralExpressionTree {
	return t.self.(*LiteralExpressionTree)
}

// AsMemberExpression casts to MemberExpressionTree.
//
// Java: ParseTree.asMemberExpression()
func (t *ParseTreeBase) AsMemberExpression() *MemberExpressionTree {
	return t.self.(*MemberExpressionTree)
}

// AsOptionalMemberExpression casts to OptionalMemberExpressionTree.
//
// Java: ParseTree.asOptionalMemberExpression()
func (t *ParseTreeBase) AsOptionalMemberExpression() *OptionalMemberExpressionTree {
	return t.self.(*OptionalMemberExpressionTree)
}

// AsMemberLookupExpression casts to MemberLookupExpressionTree.
//
// Java: ParseTree.asMemberLookupExpression()
func (t *ParseTreeBase) AsMemberLookupExpression() *MemberLookupExpressionTree {
	return t.self.(*MemberLookupExpressionTree)
}

// AsOptionalMemberLookupExpression casts to OptionalMemberLookupExpressionTree.
//
// Java: ParseTree.asOptionalMemberLookupExpression()
func (t *ParseTreeBase) AsOptionalMemberLookupExpression() *OptionalMemberLookupExpressionTree {
	return t.self.(*OptionalMemberLookupExpressionTree)
}

// AsMissingPrimaryExpression casts to MissingPrimaryExpressionTree.
//
// Java: ParseTree.asMissingPrimaryExpression()
func (t *ParseTreeBase) AsMissingPrimaryExpression() *MissingPrimaryExpressionTree {
	return t.self.(*MissingPrimaryExpressionTree)
}

// AsNewExpression casts to NewExpressionTree.
//
// Java: ParseTree.asNewExpression()
func (t *ParseTreeBase) AsNewExpression() *NewExpressionTree {
	return t.self.(*NewExpressionTree)
}

// AsNull casts to NullTree.
//
// Java: ParseTree.asNull()
func (t *ParseTreeBase) AsNull() *NullTree {
	return t.self.(*NullTree)
}

// AsObjectLiteralExpression casts to ObjectLiteralExpressionTree.
//
// Java: ParseTree.asObjectLiteralExpression()
func (t *ParseTreeBase) AsObjectLiteralExpression() *ObjectLiteralExpressionTree {
	return t.self.(*ObjectLiteralExpressionTree)
}

// AsObjectPattern casts to ObjectPatternTree.
//
// Java: ParseTree.asObjectPattern()
func (t *ParseTreeBase) AsObjectPattern() *ObjectPatternTree {
	return t.self.(*ObjectPatternTree)
}

// AsParenExpression casts to ParenExpressionTree.
//
// Java: ParseTree.asParenExpression()
func (t *ParseTreeBase) AsParenExpression() *ParenExpressionTree {
	return t.self.(*ParenExpressionTree)
}

// AsProgram casts to ProgramTree.
//
// Java: ParseTree.asProgram()
func (t *ParseTreeBase) AsProgram() *ProgramTree {
	return t.self.(*ProgramTree)
}

// AsPropertyNameAssignment casts to PropertyNameAssignmentTree.
//
// Java: ParseTree.asPropertyNameAssignment()
func (t *ParseTreeBase) AsPropertyNameAssignment() *PropertyNameAssignmentTree {
	return t.self.(*PropertyNameAssignmentTree)
}

// AsIterRest casts to IterRestTree.
//
// Java: ParseTree.asIterRest()
func (t *ParseTreeBase) AsIterRest() *IterRestTree {
	return t.self.(*IterRestTree)
}

// AsObjectRest casts to ObjectRestTree.
//
// Java: ParseTree.asObjectRest()
func (t *ParseTreeBase) AsObjectRest() *ObjectRestTree {
	return t.self.(*ObjectRestTree)
}

// AsReturnStatement casts to ReturnStatementTree.
//
// Java: ParseTree.asReturnStatement()
func (t *ParseTreeBase) AsReturnStatement() *ReturnStatementTree {
	return t.self.(*ReturnStatementTree)
}

// AsSetAccessor casts to SetAccessorTree.
//
// Java: ParseTree.asSetAccessor()
func (t *ParseTreeBase) AsSetAccessor() *SetAccessorTree {
	return t.self.(*SetAccessorTree)
}

// AsIterSpread casts to IterSpreadTree.
//
// Java: ParseTree.asIterSpread()
func (t *ParseTreeBase) AsIterSpread() *IterSpreadTree {
	return t.self.(*IterSpreadTree)
}

// AsObjectSpread casts to ObjectSpreadTree.
//
// Java: ParseTree.asObjectSpread()
func (t *ParseTreeBase) AsObjectSpread() *ObjectSpreadTree {
	return t.self.(*ObjectSpreadTree)
}

// AsSuperExpression casts to SuperExpressionTree.
//
// Java: ParseTree.asSuperExpression()
func (t *ParseTreeBase) AsSuperExpression() *SuperExpressionTree {
	return t.self.(*SuperExpressionTree)
}

// AsSwitchStatement casts to SwitchStatementTree.
//
// Java: ParseTree.asSwitchStatement()
func (t *ParseTreeBase) AsSwitchStatement() *SwitchStatementTree {
	return t.self.(*SwitchStatementTree)
}

// AsTemplateLiteralExpression casts to TemplateLiteralExpressionTree.
//
// Java: ParseTree.asTemplateLiteralExpression()
func (t *ParseTreeBase) AsTemplateLiteralExpression() *TemplateLiteralExpressionTree {
	return t.self.(*TemplateLiteralExpressionTree)
}

// AsTemplateLiteralPortion casts to TemplateLiteralPortionTree.
//
// Java: ParseTree.asTemplateLiteralPortion()
func (t *ParseTreeBase) AsTemplateLiteralPortion() *TemplateLiteralPortionTree {
	return t.self.(*TemplateLiteralPortionTree)
}

// AsTemplateSubstitution casts to TemplateSubstitutionTree.
//
// Java: ParseTree.asTemplateSubstitution()
func (t *ParseTreeBase) AsTemplateSubstitution() *TemplateSubstitutionTree {
	return t.self.(*TemplateSubstitutionTree)
}

// AsThisExpression casts to ThisExpressionTree.
//
// Java: ParseTree.asThisExpression()
func (t *ParseTreeBase) AsThisExpression() *ThisExpressionTree {
	return t.self.(*ThisExpressionTree)
}

// AsThrowStatement casts to ThrowStatementTree.
//
// Java: ParseTree.asThrowStatement()
func (t *ParseTreeBase) AsThrowStatement() *ThrowStatementTree {
	return t.self.(*ThrowStatementTree)
}

// AsTryStatement casts to TryStatementTree.
//
// Java: ParseTree.asTryStatement()
func (t *ParseTreeBase) AsTryStatement() *TryStatementTree {
	return t.self.(*TryStatementTree)
}

// AsUnaryExpression casts to UnaryExpressionTree.
//
// Java: ParseTree.asUnaryExpression()
func (t *ParseTreeBase) AsUnaryExpression() *UnaryExpressionTree {
	return t.self.(*UnaryExpressionTree)
}

// AsVariableDeclarationList casts to VariableDeclarationListTree.
//
// Java: ParseTree.asVariableDeclarationList()
func (t *ParseTreeBase) AsVariableDeclarationList() *VariableDeclarationListTree {
	return t.self.(*VariableDeclarationListTree)
}

// AsVariableDeclaration casts to VariableDeclarationTree.
//
// Java: ParseTree.asVariableDeclaration()
func (t *ParseTreeBase) AsVariableDeclaration() *VariableDeclarationTree {
	return t.self.(*VariableDeclarationTree)
}

// AsVariableStatement casts to VariableStatementTree.
//
// Java: ParseTree.asVariableStatement()
func (t *ParseTreeBase) AsVariableStatement() *VariableStatementTree {
	return t.self.(*VariableStatementTree)
}

// AsWhileStatement casts to WhileStatementTree.
//
// Java: ParseTree.asWhileStatement()
func (t *ParseTreeBase) AsWhileStatement() *WhileStatementTree {
	return t.self.(*WhileStatementTree)
}

// AsWithStatement casts to WithStatementTree.
//
// Java: ParseTree.asWithStatement()
func (t *ParseTreeBase) AsWithStatement() *WithStatementTree {
	return t.self.(*WithStatementTree)
}

// AsYieldStatement casts to YieldExpressionTree.
//
// Java: ParseTree.asYieldStatement()
func (t *ParseTreeBase) AsYieldStatement() *YieldExpressionTree {
	return t.self.(*YieldExpressionTree)
}

// AsAwaitExpression casts to AwaitExpressionTree.
//
// Java: ParseTree.asAwaitExpression()
func (t *ParseTreeBase) AsAwaitExpression() *AwaitExpressionTree {
	return t.self.(*AwaitExpressionTree)
}

// AsNewTargetExpression casts to NewTargetExpressionTree.
//
// Java: ParseTree.asNewTargetExpression()
func (t *ParseTreeBase) AsNewTargetExpression() *NewTargetExpressionTree {
	return t.self.(*NewTargetExpressionTree)
}

// AsUpdateExpression casts to UpdateExpressionTree.
//
// Java: ParseTree.asUpdateExpression()
func (t *ParseTreeBase) AsUpdateExpression() *UpdateExpressionTree {
	return t.self.(*UpdateExpressionTree)
}

// AsForAwaitOfStatement casts to ForAwaitOfStatementTree.
//
// Java: ParseTree.asForAwaitOfStatement()
func (t *ParseTreeBase) AsForAwaitOfStatement() *ForAwaitOfStatementTree {
	return t.self.(*ForAwaitOfStatementTree)
}

// AsFieldDeclaration casts to FieldDeclarationTree.
//
// Java: ParseTree.asFieldDeclaration()
func (t *ParseTreeBase) AsFieldDeclaration() *FieldDeclarationTree {
	return t.self.(*FieldDeclarationTree)
}

// IsPattern reports whether the tree (inside any parentheses) is an array or
// object pattern.
//
// Java: ParseTree.isPattern()
func (t *ParseTreeBase) IsPattern() bool {
	parseTree := t.self
	for parseTree.GetParseTreeBase().Type == ParseTreeType_PAREN_EXPRESSION {
		parseTree = parseTree.AsParenExpression().Expression
	}

	switch parseTree.GetParseTreeBase().Type {
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
func (t *ParseTreeBase) IsValidNonVanillaAssignmentTarget() bool {
	parseTree := t.self
	for parseTree.GetParseTreeBase().Type == ParseTreeType_PAREN_EXPRESSION {
		parseTree = parseTree.AsParenExpression().Expression
	}
	switch parseTree.GetParseTreeBase().Type {
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
// Java: ParseTree.isValidAssignmentTarget()
func (t *ParseTreeBase) IsValidAssignmentTarget() bool {
	parseTree := t.self
	for parseTree.GetParseTreeBase().Type == ParseTreeType_PAREN_EXPRESSION {
		parseTree = parseTree.AsParenExpression().Expression
	}

	switch parseTree.GetParseTreeBase().Type {
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
