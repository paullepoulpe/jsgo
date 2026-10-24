// Port of com/google/javascript/jscomp/RewriteLogicalAssignmentOperatorsHelper.java @ closure-compiler v20261021

/*
 * Copyright 2021 The Closure Compiler Authors.
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

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// RewriteLogicalAssignmentOperatorsHelper is a helper class for
// RewriteLogicalAssignmentOperators.
type RewriteLogicalAssignmentOperatorsHelper struct {
	compiler         AbstractCompiler
	astFactory       *AstFactory
	uniqueIdSupplier *UniqueIdSupplier
}

// Java: RewriteLogicalAssignmentOperatorsHelper.TEMP_VAR_NAME_PREFIX
const RewriteLogicalAssignmentOperatorsHelper_TEMP_VAR_NAME_PREFIX = "$jscomp$logical$assign$tmp"

// Java: RewriteLogicalAssignmentOperatorsHelper.TEMP_INDEX_VAR_NAME_PREFIX
const RewriteLogicalAssignmentOperatorsHelper_TEMP_INDEX_VAR_NAME_PREFIX = "$jscomp$logical$assign$tmpindex"

// Java: RewriteLogicalAssignmentOperatorsHelper(AbstractCompiler, AstFactory, UniqueIdSupplier)
func NewRewriteLogicalAssignmentOperatorsHelper(
	compiler AbstractCompiler, astFactory *AstFactory, uniqueIdSupplier *UniqueIdSupplier) *RewriteLogicalAssignmentOperatorsHelper {
	h := &RewriteLogicalAssignmentOperatorsHelper{}
	h.compiler = compiler
	h.astFactory = astFactory
	h.uniqueIdSupplier = uniqueIdSupplier
	return h
}

// Java: RewriteLogicalAssignmentOperatorsHelper.visitLogicalAssignmentOperator(NodeTraversal, Node)
func (h *RewriteLogicalAssignmentOperatorsHelper) VisitLogicalAssignmentOperator(t *NodeTraversal, logicalAssignment *rhino.Node) {
	enclosingFunction := NodeUtil_GetEnclosingFunction(logicalAssignment)
	if enclosingFunction != nil && !NodeUtil_GetFunctionBody(enclosingFunction).IsBlock() {
		returnValue := NodeUtil_GetFunctionBody(enclosingFunction)
		body := rhino.IR_BlockOfNode(rhino.IR_ReturnNodeWithExpr(returnValue.Detach())).SrcrefTreeIfMissing(returnValue)
		enclosingFunction.AddChildToBack(body)
		h.compiler.ReportChangeToEnclosingScope(body)
	}

	enclosingStatement := NodeUtil_GetEnclosingStatement(logicalAssignment)

	left := logicalAssignment.RemoveFirstChild()
	right := logicalAssignment.GetLastChild().Detach()

	var replacement *rhino.Node

	for left.IsCast() {
		// This pass runs after type checking, which is what requires the CAST nodes.
		// They aren't needed after that, but they are still in the AST until
		// Normalization removes them. So, it's safe to just remove this CAST now to simplify things.
		// If this pass moves after Normalization, then it will no longer be necessary to check for
		// CAST nodes.
		left = left.RemoveFirstChild()
	}

	if left.IsName() {
		replacement = h.HandleLHSName(logicalAssignment, left, right)

	} else {
		jutil.CheckStateMsg1(left.IsGetProp() || left.IsGetElem(), "%s", left)
		replacement = h.HandleLHSPropertyReference(t, logicalAssignment, left, right, enclosingStatement)
	}

	logicalAssignment.ReplaceWith(replacement)

	h.compiler.ReportChangeToEnclosingScope(enclosingStatement)
}

// Java: RewriteLogicalAssignmentOperatorsHelper.handleLHSName(Node, Node, Node)
func (h *RewriteLogicalAssignmentOperatorsHelper) HandleLHSName(logicalAssignment *rhino.Node, left *rhino.Node, right *rhino.Node) *rhino.Node {
	// handle name case
	// e.g. convert `name ??= something()` to
	// `name ?? (name = something())`
	assignToRHS := h.astFactory.CreateAssign(left, right).Srcref(right)

	// TODO(bradfordcsmith): We should have an AstFactory method for this.
	return rhino.NewNodeWith2Children(
		NodeUtil_GetOpFromAssignmentOp(logicalAssignment), left.CloneNode(), assignToRHS).
		CopyTypeFrom(logicalAssignment).
		Srcref(logicalAssignment)
}

// Java: RewriteLogicalAssignmentOperatorsHelper.handleLHSPropertyReference(NodeTraversal, Node, Node, Node, Node)
func (h *RewriteLogicalAssignmentOperatorsHelper) HandleLHSPropertyReference(
	t *NodeTraversal, logicalAssignment *rhino.Node, left *rhino.Node, right *rhino.Node, enclosingStatement *rhino.Node) *rhino.Node {
	// handle getprop case and getelem case
	input := t.GetInput()
	uniqueId := h.uniqueIdSupplier.GetUniqueId(input)
	tempVarName := RewriteLogicalAssignmentOperatorsHelper_TEMP_VAR_NAME_PREFIX + uniqueId

	var assignToRHS *rhino.Node
	var newLHS *rhino.Node

	objectNode := left.RemoveFirstChild()

	let := h.astFactory.CreateSingleLetNameDeclaration(tempVarName).SrcrefTree(logicalAssignment)
	let.InsertBefore(enclosingStatement)

	tempName := h.astFactory.CreateName(tempVarName, AstFactory_Type_(objectNode)).Srcref(objectNode)
	assignTemp :=
		h.astFactory.CreateAssign(tempName, objectNode).Srcref(objectNode) // (tmp = someExpression)

	if left.IsGetProp() {
		// handle getprop case
		// e.g. convert `(someExpression).property ??= something()` to
		// let tmp;
		// (tmp = someExpression).property ?? (tmp.property = something());
		propertyName := left.GetString()

		tempProp := h.astFactory.
			CreateGetProp(tempName.CloneNode(), propertyName, AstFactory_Type_(left)).
			Srcref(right) // (tmp.property)
		assignToRHS =
			h.astFactory.CreateAssign(tempProp, right).Srcref(right) // (tmp.property = something())
		newLHS = h.astFactory.
			CreateGetProp(assignTemp, propertyName, AstFactory_Type_(left)).
			Srcref(left) // ((tmp = someExpression).property)
	} else {
		// handle getelem case
		// e.g. convert `someExpression[indexExpression] ??= something()` to
		// let tmp;
		// let tmpIndex;
		// (tmp = someExpression)[tmpIndex = indexExpression] ?? (tmp[tmpIndex] = something());
		jutil.CheckStateMsg1(left.IsGetElem(), "%s", left)
		tempIndexVarName := RewriteLogicalAssignmentOperatorsHelper_TEMP_INDEX_VAR_NAME_PREFIX + uniqueId

		indexExprNode := left.GetLastChild().Detach()

		letIndex := h.astFactory.CreateSingleLetNameDeclaration(tempIndexVarName).SrcrefTree(logicalAssignment)
		letIndex.InsertBefore(enclosingStatement)

		tempIndexName := h.astFactory.
			CreateName(tempIndexVarName, AstFactory_Type_(indexExprNode)).
			Srcref(indexExprNode) // tmpIndex

		assignTempIndex := h.astFactory.
			CreateAssign(tempIndexName, indexExprNode).
			Srcref(indexExprNode) // [tmpIndex = indexExpression]

		tempElem := h.astFactory.
			CreateGetElem(tempName.CloneNode(), tempIndexName.CloneNode()).
			CopyTypeFrom(left).
			Srcref(right) // (tmp[tmpIndex])

		assignToRHS =
			h.astFactory.CreateAssign(tempElem, right).Srcref(right) // (tmp[tmpIndex] = something())
		newLHS = h.astFactory.
			CreateGetElem(assignTemp, assignTempIndex).
			CopyTypeFrom(left).
			Srcref(left) // (tmp = someExpression)[tmpIndex = indexExpression]
	}

	NodeUtil_AddFeatureToScript(t.GetCurrentScript(), parser.FeatureSet_Feature_LET_DECLARATIONS, h.compiler)

	// TODO(bradfordcsmith): We should have an AstFactory method for this.
	return rhino.NewNodeWith2Children(NodeUtil_GetOpFromAssignmentOp(logicalAssignment), newLHS, assignToRHS).
		CopyTypeFrom(logicalAssignment).
		Srcref(logicalAssignment)
}
