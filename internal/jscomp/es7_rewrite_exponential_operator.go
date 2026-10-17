// Port of com/google/javascript/jscomp/Es7RewriteExponentialOperator.java @ closure-compiler v20261015

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

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Es7RewriteExponentialOperator replaces the ES7 `**` and `**=` operators to
// calls to `Math.pow`.
type Es7RewriteExponentialOperator struct {
	compiler         AbstractCompiler
	astFactory       *AstFactory
	uniqueIdSupplier *UniqueIdSupplier
	// This node should only ever be cloned, not directly inserted.
	//
	// Go-only: language-adapter — Suppliers.memoize(this::createMathPowCall)
	// is the lazily set mathPowCallValue, read through mathPowCall().
	mathPowCallValue *rhino.Node
}

// Java: Es7RewriteExponentialOperator.TRANSPILE_EXPONENT_USING_BIGINT
var Es7RewriteExponentialOperator_TRANSPILE_EXPONENT_USING_BIGINT = DiagnosticType_Error(
	"JSC_TRANSPILE_EXPONENT_USING_BIGINT",
	"Cannot transpile `**` operator applied to BigInt operands.")

const (
	// Java: Es7RewriteExponentialOperator.TEMP_VAR_NAME_PREFIX
	Es7RewriteExponentialOperator_TEMP_VAR_NAME_PREFIX = "$jscomp$exp$assign$tmp"
	// Java: Es7RewriteExponentialOperator.TEMP_INDEX_VAR_NAME_PREFIX
	Es7RewriteExponentialOperator_TEMP_INDEX_VAR_NAME_PREFIX = "$jscomp$exp$assign$tmpindex"
)

// Java: Es7RewriteExponentialOperator(AbstractCompiler)
func NewEs7RewriteExponentialOperator(compiler AbstractCompiler) *Es7RewriteExponentialOperator {
	e := &Es7RewriteExponentialOperator{}
	e.compiler = compiler
	e.astFactory = compiler.CreateAstFactory()
	e.uniqueIdSupplier = compiler.GetUniqueIdSupplier()
	return e
}

// mathPowCall is mathPowCall.get() of the memoizing supplier.
//
// Go-only: language-adapter — Suppliers.memoize.
func (e *Es7RewriteExponentialOperator) mathPowCall() *rhino.Node {
	if e.mathPowCallValue == nil {
		e.mathPowCallValue = e.CreateMathPowCall()
	}
	return e.mathPowCallValue
}

// Java: Es7RewriteExponentialOperator.process(Node, Node)
func (e *Es7RewriteExponentialOperator) Process(externs *rhino.Node, root *rhino.Node) {
	NodeTraversal_Traverse(e.compiler, root, e)
	TranspilationPasses_MaybeMarkFeatureAsTranspiledAway(e.compiler, root, parser.FeatureSet_Feature_EXPONENT_OP)
}

// Java: Es7RewriteExponentialOperator.shouldTraverse(NodeTraversal, Node, Node)
func (e *Es7RewriteExponentialOperator) ShouldTraverse(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
	if n.IsScript() {
		scriptFeatures := NodeUtil_GetFeatureSetOfScript(n)
		// Runs only if features exist in the script
		return scriptFeatures == nil || scriptFeatures.ContainsFeature(parser.FeatureSet_Feature_EXPONENT_OP)
	}
	return true
}

// Java: Es7RewriteExponentialOperator.visit(NodeTraversal, Node, Node)
func (e *Es7RewriteExponentialOperator) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_EXPONENT:
		if e.CheckOperatorType(n) {
			e.VisitExponentiationOperator(n)
		}
	case rhino.Token_ASSIGN_EXPONENT:
		if e.CheckOperatorType(n) {
			e.VisitExponentiationAssignmentOperator(t, n)
		}
	default:
	}
}

// Java: Es7RewriteExponentialOperator.visitExponentiationOperator(Node)
func (e *Es7RewriteExponentialOperator) VisitExponentiationOperator(operator *rhino.Node) {
	callClone := e.mathPowCall().CloneTree()
	callClone.AddChildToBack(operator.RemoveFirstChild()) // Base argument.
	callClone.AddChildToBack(operator.RemoveFirstChild()) // Exponent argument.

	callClone.SrcrefTreeIfMissing(operator)
	operator.ReplaceWith(callClone)

	e.compiler.ReportChangeToEnclosingScope(callClone)
}

// Java: Es7RewriteExponentialOperator.visitExponentiationAssignmentOperator(NodeTraversal, Node)
func (e *Es7RewriteExponentialOperator) VisitExponentiationAssignmentOperator(t *NodeTraversal, operator *rhino.Node) {
	enclosingStatement := NodeUtil_GetEnclosingStatement(operator)

	left := operator.RemoveFirstChild()

	var replacement *rhino.Node

	if left.IsName() {
		replacement = e.HandleLHSName(operator, left)
	} else {
		jutil.CheckStateMsg1(left.IsGetProp() || left.IsGetElem(), "%s", left)
		replacement = e.HandleLHSPropertyReference(t, operator, left, enclosingStatement)
	}

	operator.ReplaceWith(replacement)

	e.compiler.ReportChangeToEnclosingScope(enclosingStatement)
}

// Java: Es7RewriteExponentialOperator.handleLHSName(Node, Node)
func (e *Es7RewriteExponentialOperator) HandleLHSName(operator *rhino.Node, left *rhino.Node) *rhino.Node {
	// handle name case
	// e.g. convert `name **= value` to
	// `name = Math.pow(name, value)`
	callClone := e.mathPowCall().CloneTree()
	callClone.AddChildToBack(left.CloneTree())            // Base argument.
	callClone.AddChildToBack(operator.RemoveFirstChild()) // Exponent argument.

	return e.astFactory.CreateAssign(left, callClone).SrcrefTreeIfMissing(operator)
}

// Java: Es7RewriteExponentialOperator.handleLHSPropertyReference(NodeTraversal, Node, Node, Node)
func (e *Es7RewriteExponentialOperator) HandleLHSPropertyReference(
	t *NodeTraversal, operator *rhino.Node, left *rhino.Node, enclosingStatement *rhino.Node) *rhino.Node {
	// handle getprop case and getelem case
	input := t.GetInput()
	uniqueId := e.uniqueIdSupplier.GetUniqueId(input)
	tempVarName := Es7RewriteExponentialOperator_TEMP_VAR_NAME_PREFIX + uniqueId

	var newLHS *rhino.Node
	var tempPropOrElem *rhino.Node

	objectNode := left.RemoveFirstChild()

	let := e.astFactory.CreateSingleLetNameDeclaration(tempVarName).SrcrefTree(operator)
	let.InsertBefore(enclosingStatement)

	tempName := e.astFactory.CreateName(tempVarName, AstFactory_Type_(objectNode)).Srcref(objectNode)
	assignTemp :=
		e.astFactory.CreateAssign(tempName, objectNode).Srcref(objectNode) // (tmp = someExpression)

	if left.IsGetProp() {
		// handle getprop case
		// e.g. convert `(someExpression).property **= value` to
		// let tmp;
		// (tmp = someExpression).property = Math.pow(tmp.property, value);
		propertyName := left.GetString()

		tempPropOrElem =
			e.astFactory.
				CreateGetProp(tempName.CloneNode(), propertyName, AstFactory_Type_(left)).
				Srcref(left) // (tmp.property)
		newLHS =
			e.astFactory.
				CreateGetProp(assignTemp, propertyName, AstFactory_Type_(left)).
				Srcref(left) // ((tmp = someExpression).property)
	} else {
		// handle getelem case
		// e.g. convert `someExpression[indexExpression] **= value` to
		// let tmp;
		// let tmpIndex;
		// `(tmp = someExpression)[(tmpIndex = indexExpression)] = Math.pow(tmp[tmpIndex], value);`
		jutil.CheckStateMsg1(left.IsGetElem(), "%s", left)
		tempIndexVarName := Es7RewriteExponentialOperator_TEMP_INDEX_VAR_NAME_PREFIX + uniqueId

		indexExprNode := left.GetLastChild().Detach()

		letIndex :=
			e.astFactory.CreateSingleLetNameDeclaration(tempIndexVarName).SrcrefTree(operator)
		letIndex.InsertBefore(enclosingStatement)

		tempIndexName :=
			e.astFactory.
				CreateName(tempIndexVarName, AstFactory_Type_(indexExprNode)).
				Srcref(indexExprNode) // tmpIndex

		assignTempIndex :=
			e.astFactory.
				CreateAssign(tempIndexName, indexExprNode).
				Srcref(indexExprNode) // [tmpIndex = indexExpression]

		tempPropOrElem =
			e.astFactory.
				CreateGetElem(tempName.CloneNode(), tempIndexName.CloneNode()).
				CopyTypeFrom(left).
				Srcref(left) // (tmp[tmpIndex])

		newLHS =
			e.astFactory.
				CreateGetElem(assignTemp, assignTempIndex).
				CopyTypeFrom(left).
				Srcref(left) // (tmp = someExpression)[tmpIndex = indexExpression]
	}
	callClone := e.mathPowCall().CloneTree()
	callClone.AddChildToBack(tempPropOrElem.CloneTree())  // Base argument.
	callClone.AddChildToBack(operator.RemoveFirstChild()) // Exponent argument.

	NodeUtil_AddFeatureToScript(t.GetCurrentScript(), parser.FeatureSet_Feature_LET_DECLARATIONS, e.compiler)

	return e.astFactory.CreateAssign(newLHS, callClone).SrcrefTreeIfMissing(operator)
}

// Java: Es7RewriteExponentialOperator.createMathPowCall()
func (e *Es7RewriteExponentialOperator) CreateMathPowCall() *rhino.Node {
	return e.astFactory.CreateCall(
		e.astFactory.CreateQName(e.compiler.GetTranspilationNamespace(), "Math.pow"),
		AstFactory_Type_OfColor(colors.StandardColors_NUMBER))
}

// CheckOperatorType reports an error if the `**` getting transpiled to
// `Math.pow()`is of BIGINT type
//
// Java: Es7RewriteExponentialOperator.checkOperatorType(Node)
func (e *Es7RewriteExponentialOperator) CheckOperatorType(operator *rhino.Node) bool {
	jutil.CheckArgumentMsg1(operator.IsExponent() || operator.IsAssignExponent(), "%s", operator)
	if e.compiler.HasTypeCheckingRun() {
		if operator.GetColor().Equals(colors.StandardColors_BIGINT) {
			e.compiler.Report(JSError_MakeWithNode(operator, Es7RewriteExponentialOperator_TRANSPILE_EXPONENT_USING_BIGINT))
			return false
		}
	}
	return true
}
