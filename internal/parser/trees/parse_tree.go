// Port of com/google/javascript/jscomp/parsing/parser/trees/ParseTree.java
// @ closure-compiler v20260415

package trees

import (
	"fmt"
	"unsafe"
)

// ParseTree is the common base for all JavaScript parse-tree nodes.
// Immutable; a plain-old-data structure.
//
// Derived types embed ParseTree (mirroring Java's class hierarchy via
// Go embedding). Every concrete type has a matching ParseTreeType constant
// whose name corresponds to the struct name.
//
// A tree derived from source has a non-zero Location; a synthesised tree
// may have a zero Location.
type ParseTree struct {
	Type     ParseTreeType
	Location SourceRange
}

// NewParseTree constructs the base ParseTree value used by concrete ctor functions.
func NewParseTree(treeType ParseTreeType, location SourceRange) ParseTree {
	return ParseTree{Type: treeType, Location: location}
}

func (t *ParseTree) GetStart() SourcePosition {
	return t.Location.Start
}

func (t *ParseTree) GetEnd() SourcePosition {
	return t.Location.End
}

func (t *ParseTree) IsPattern() bool {
	tree := t
	for tree.Type == ParseTreeTypeParenExpression {
		tree = &tree.AsParenExpression().ParseTree
	}
	switch tree.Type {
	case ParseTreeTypeArrayPattern, ParseTreeTypeObjectPattern:
		return true
	}
	return false
}

// IsValidNonVanillaAssignmentTarget reports whether this tree is a valid
// assignment target for non-vanilla operators like +=, -=, **=, etc.
func (t *ParseTree) IsValidNonVanillaAssignmentTarget() bool {
	tree := t
	for tree.Type == ParseTreeTypeParenExpression {
		tree = &tree.AsParenExpression().ParseTree
	}
	switch tree.Type {
	case ParseTreeTypeIdentifierExpression,
		ParseTreeTypeMemberExpression,
		ParseTreeTypeMemberLookupExpression,
		ParseTreeTypeDefaultParameter:
		return true
	}
	return false
}

// IsValidAssignmentTarget reports whether this tree is a valid assignment
// target for any assignment operator (=, +=, -=, **=, etc.).
func (t *ParseTree) IsValidAssignmentTarget() bool {
	tree := t
	for tree.Type == ParseTreeTypeParenExpression {
		tree = &tree.AsParenExpression().ParseTree
	}
	switch tree.Type {
	case ParseTreeTypeIdentifierExpression,
		ParseTreeTypeMemberExpression,
		ParseTreeTypeMemberLookupExpression,
		ParseTreeTypeArrayPattern,
		ParseTreeTypeObjectPattern,
		ParseTreeTypeDefaultParameter:
		return true
	}
	return false
}

func (t *ParseTree) IsRestParameter() bool {
	return t.Type == ParseTreeTypeIterRest
}

func (t *ParseTree) String() string {
	return fmt.Sprintf("%v@%v", t.Type, t.Location)
}

// ---------------------------------------------------------------------------
// As-casts – mirror Java's downcasting helpers. Each panics if the node's
// Type does not match the expected constant (mirroring Java ClassCastException).
//
// Because every concrete tree type starts with an embedded ParseTree at
// offset 0, unsafe.Pointer conversion is sound: the *ParseTree points at
// the same address as the outer concrete struct.
// ---------------------------------------------------------------------------

func (t *ParseTree) AsArrayLiteralExpression() *ArrayLiteralExpressionTree {
	t.assertType(ParseTreeTypeArrayLiteralExpression)
	return (*ArrayLiteralExpressionTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsArrayPattern() *ArrayPatternTree {
	t.assertType(ParseTreeTypeArrayPattern)
	return (*ArrayPatternTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsBinaryOperator() *BinaryOperatorTree {
	t.assertType(ParseTreeTypeBinaryOperator)
	return (*BinaryOperatorTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsBlock() *BlockTree {
	t.assertType(ParseTreeTypeBlock)
	return (*BlockTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsBreakStatement() *BreakStatementTree {
	t.assertType(ParseTreeTypeBreakStatement)
	return (*BreakStatementTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsCallExpression() *CallExpressionTree {
	t.assertType(ParseTreeTypeCallExpression)
	return (*CallExpressionTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsOptChainCallExpression() *OptChainCallExpressionTree {
	t.assertType(ParseTreeTypeOptChainCallExpression)
	return (*OptChainCallExpressionTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsCaseClause() *CaseClauseTree {
	t.assertType(ParseTreeTypeCaseClause)
	return (*CaseClauseTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsCatch() *CatchTree {
	t.assertType(ParseTreeTypeCatch)
	return (*CatchTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsClassDeclaration() *ClassDeclarationTree {
	t.assertType(ParseTreeTypeClassDeclaration)
	return (*ClassDeclarationTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsCommaExpression() *CommaExpressionTree {
	t.assertType(ParseTreeTypeCommaExpression)
	return (*CommaExpressionTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsComprehensionIf() *ComprehensionIfTree {
	t.assertType(ParseTreeTypeComprehensionIf)
	return (*ComprehensionIfTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsComprehensionFor() *ComprehensionForTree {
	t.assertType(ParseTreeTypeComprehensionFor)
	return (*ComprehensionForTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsComprehension() *ComprehensionTree {
	t.assertType(ParseTreeTypeComprehension)
	return (*ComprehensionTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsComputedPropertyDefinition() *ComputedPropertyDefinitionTree {
	t.assertType(ParseTreeTypeComputedPropertyDefinition)
	return (*ComputedPropertyDefinitionTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsComputedPropertyField() *ComputedPropertyFieldTree {
	t.assertType(ParseTreeTypeComputedPropertyField)
	return (*ComputedPropertyFieldTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsComputedPropertyGetter() *ComputedPropertyGetterTree {
	t.assertType(ParseTreeTypeComputedPropertyGetter)
	return (*ComputedPropertyGetterTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsComputedPropertyMethod() *ComputedPropertyMethodTree {
	t.assertType(ParseTreeTypeComputedPropertyMethod)
	return (*ComputedPropertyMethodTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsComputedPropertySetter() *ComputedPropertySetterTree {
	t.assertType(ParseTreeTypeComputedPropertySetter)
	return (*ComputedPropertySetterTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsConditionalExpression() *ConditionalExpressionTree {
	t.assertType(ParseTreeTypeConditionalExpression)
	return (*ConditionalExpressionTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsContinueStatement() *ContinueStatementTree {
	t.assertType(ParseTreeTypeContinueStatement)
	return (*ContinueStatementTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsDebuggerStatement() *DebuggerStatementTree {
	t.assertType(ParseTreeTypeDebuggerStatement)
	return (*DebuggerStatementTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsDefaultClause() *DefaultClauseTree {
	t.assertType(ParseTreeTypeDefaultClause)
	return (*DefaultClauseTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsDefaultParameter() *DefaultParameterTree {
	t.assertType(ParseTreeTypeDefaultParameter)
	return (*DefaultParameterTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsDoWhileStatement() *DoWhileStatementTree {
	t.assertType(ParseTreeTypeDoWhileStatement)
	return (*DoWhileStatementTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsEmptyStatement() *EmptyStatementTree {
	t.assertType(ParseTreeTypeEmptyStatement)
	return (*EmptyStatementTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsExportDeclaration() *ExportDeclarationTree {
	t.assertType(ParseTreeTypeExportDeclaration)
	return (*ExportDeclarationTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsExportSpecifier() *ExportSpecifierTree {
	t.assertType(ParseTreeTypeExportSpecifier)
	return (*ExportSpecifierTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsExpressionStatement() *ExpressionStatementTree {
	t.assertType(ParseTreeTypeExpressionStatement)
	return (*ExpressionStatementTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsFinally() *FinallyTree {
	t.assertType(ParseTreeTypeFinally)
	return (*FinallyTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsForOfStatement() *ForOfStatementTree {
	t.assertType(ParseTreeTypeForOfStatement)
	return (*ForOfStatementTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsForInStatement() *ForInStatementTree {
	t.assertType(ParseTreeTypeForInStatement)
	return (*ForInStatementTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsFormalParameterList() *FormalParameterListTree {
	t.assertType(ParseTreeTypeFormalParameterList)
	return (*FormalParameterListTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsForStatement() *ForStatementTree {
	t.assertType(ParseTreeTypeForStatement)
	return (*ForStatementTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsFunctionDeclaration() *FunctionDeclarationTree {
	t.assertType(ParseTreeTypeFunctionDeclaration)
	return (*FunctionDeclarationTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsGetAccessor() *GetAccessorTree {
	t.assertType(ParseTreeTypeGetAccessor)
	return (*GetAccessorTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsIdentifierExpression() *IdentifierExpressionTree {
	t.assertType(ParseTreeTypeIdentifierExpression)
	return (*IdentifierExpressionTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsIfStatement() *IfStatementTree {
	t.assertType(ParseTreeTypeIfStatement)
	return (*IfStatementTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsImportDeclaration() *ImportDeclarationTree {
	t.assertType(ParseTreeTypeImportDeclaration)
	return (*ImportDeclarationTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsImportSpecifier() *ImportSpecifierTree {
	t.assertType(ParseTreeTypeImportSpecifier)
	return (*ImportSpecifierTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsDynamicImportExpression() *DynamicImportTree {
	t.assertType(ParseTreeTypeDynamicImportExpression)
	return (*DynamicImportTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsImportMetaExpression() *ImportMetaExpressionTree {
	t.assertType(ParseTreeTypeImportMetaExpression)
	return (*ImportMetaExpressionTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsLabelledStatement() *LabelledStatementTree {
	t.assertType(ParseTreeTypeLabelledStatement)
	return (*LabelledStatementTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsLiteralExpression() *LiteralExpressionTree {
	t.assertType(ParseTreeTypeLiteralExpression)
	return (*LiteralExpressionTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsMemberExpression() *MemberExpressionTree {
	t.assertType(ParseTreeTypeMemberExpression)
	return (*MemberExpressionTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsOptionalMemberExpression() *OptionalMemberExpressionTree {
	t.assertType(ParseTreeTypeOptChainMemberExpression)
	return (*OptionalMemberExpressionTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsMemberLookupExpression() *MemberLookupExpressionTree {
	t.assertType(ParseTreeTypeMemberLookupExpression)
	return (*MemberLookupExpressionTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsOptionalMemberLookupExpression() *OptionalMemberLookupExpressionTree {
	t.assertType(ParseTreeTypeOptChainMemberLookupExpression)
	return (*OptionalMemberLookupExpressionTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsMissingPrimaryExpression() *MissingPrimaryExpressionTree {
	t.assertType(ParseTreeTypeMissingPrimaryExpression)
	return (*MissingPrimaryExpressionTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsNewExpression() *NewExpressionTree {
	t.assertType(ParseTreeTypeNewExpression)
	return (*NewExpressionTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsNull() *NullTree {
	t.assertType(ParseTreeTypeNull)
	return (*NullTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsObjectLiteralExpression() *ObjectLiteralExpressionTree {
	t.assertType(ParseTreeTypeObjectLiteralExpression)
	return (*ObjectLiteralExpressionTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsObjectPattern() *ObjectPatternTree {
	t.assertType(ParseTreeTypeObjectPattern)
	return (*ObjectPatternTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsParenExpression() *ParenExpressionTree {
	t.assertType(ParseTreeTypeParenExpression)
	return (*ParenExpressionTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsProgram() *ProgramTree {
	t.assertType(ParseTreeTypeProgram)
	return (*ProgramTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsPropertyNameAssignment() *PropertyNameAssignmentTree {
	t.assertType(ParseTreeTypePropertyNameAssignment)
	return (*PropertyNameAssignmentTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsIterRest() *IterRestTree {
	t.assertType(ParseTreeTypeIterRest)
	return (*IterRestTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsObjectRest() *ObjectRestTree {
	t.assertType(ParseTreeTypeObjectRest)
	return (*ObjectRestTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsReturnStatement() *ReturnStatementTree {
	t.assertType(ParseTreeTypeReturnStatement)
	return (*ReturnStatementTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsSetAccessor() *SetAccessorTree {
	t.assertType(ParseTreeTypeSetAccessor)
	return (*SetAccessorTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsIterSpread() *IterSpreadTree {
	t.assertType(ParseTreeTypeIterSpread)
	return (*IterSpreadTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsObjectSpread() *ObjectSpreadTree {
	t.assertType(ParseTreeTypeObjectSpread)
	return (*ObjectSpreadTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsSuperExpression() *SuperExpressionTree {
	t.assertType(ParseTreeTypeSuperExpression)
	return (*SuperExpressionTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsSwitchStatement() *SwitchStatementTree {
	t.assertType(ParseTreeTypeSwitchStatement)
	return (*SwitchStatementTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsTemplateLiteralExpression() *TemplateLiteralExpressionTree {
	t.assertType(ParseTreeTypeTemplateLiteralExpression)
	return (*TemplateLiteralExpressionTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsTemplateLiteralPortion() *TemplateLiteralPortionTree {
	t.assertType(ParseTreeTypeTemplateLiteralPortion)
	return (*TemplateLiteralPortionTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsTemplateSubstitution() *TemplateSubstitutionTree {
	t.assertType(ParseTreeTypeTemplateSubstitution)
	return (*TemplateSubstitutionTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsThisExpression() *ThisExpressionTree {
	t.assertType(ParseTreeTypeThisExpression)
	return (*ThisExpressionTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsThrowStatement() *ThrowStatementTree {
	t.assertType(ParseTreeTypeThrowStatement)
	return (*ThrowStatementTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsTryStatement() *TryStatementTree {
	t.assertType(ParseTreeTypeTryStatement)
	return (*TryStatementTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsUnaryExpression() *UnaryExpressionTree {
	t.assertType(ParseTreeTypeUnaryExpression)
	return (*UnaryExpressionTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsVariableDeclarationList() *VariableDeclarationListTree {
	t.assertType(ParseTreeTypeVariableDeclarationList)
	return (*VariableDeclarationListTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsVariableDeclaration() *VariableDeclarationTree {
	t.assertType(ParseTreeTypeVariableDeclaration)
	return (*VariableDeclarationTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsVariableStatement() *VariableStatementTree {
	t.assertType(ParseTreeTypeVariableStatement)
	return (*VariableStatementTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsWhileStatement() *WhileStatementTree {
	t.assertType(ParseTreeTypeWhileStatement)
	return (*WhileStatementTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsWithStatement() *WithStatementTree {
	t.assertType(ParseTreeTypeWithStatement)
	return (*WithStatementTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsYieldStatement() *YieldExpressionTree {
	t.assertType(ParseTreeTypeYieldExpression)
	return (*YieldExpressionTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsAwaitExpression() *AwaitExpressionTree {
	t.assertType(ParseTreeTypeAwaitExpression)
	return (*AwaitExpressionTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsNewTargetExpression() *NewTargetExpressionTree {
	t.assertType(ParseTreeTypeNewTargetExpression)
	return (*NewTargetExpressionTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsUpdateExpression() *UpdateExpressionTree {
	t.assertType(ParseTreeTypeUpdateExpression)
	return (*UpdateExpressionTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsForAwaitOfStatement() *ForAwaitOfStatementTree {
	t.assertType(ParseTreeTypeForAwaitOfStatement)
	return (*ForAwaitOfStatementTree)(unsafe.Pointer(t))
}
func (t *ParseTree) AsFieldDeclaration() *FieldDeclarationTree {
	t.assertType(ParseTreeTypeFieldDeclaration)
	return (*FieldDeclarationTree)(unsafe.Pointer(t))
}

func (t *ParseTree) assertType(expected ParseTreeType) {
	if t.Type != expected {
		panic(fmt.Sprintf("ParseTree.As*: expected %v, got %v", expected, t.Type))
	}
}
