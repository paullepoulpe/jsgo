// Port of com/google/javascript/jscomp/lint/CheckNullabilityModifiers.java @ closure-compiler v20261015

/*
 * Copyright 2018 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// CheckNullabilityModifiers checks for missing or redundant nullability
// modifiers.
//
// Primitive and literal types should not be preceded by a `!` modifier.
// Reference types must be preceded by a `?` or `!` modifier.
type CheckNullabilityModifiers struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler AbstractCompiler

	// Store the candidate warnings and template types found while traversing a single script node.
	redundantCandidates   *jutil.LinkedHashSet[rhino.Node]
	missingCandidates     *jutil.LinkedHashSet[rhino.Node]
	nullMissingCandidates *jutil.LinkedHashSet[rhino.Node]
	templateTypeNames     *jutil.LinkedHashSet[string]
}

// CheckNullabilityModifiers_MISSING_NULLABILITY_MODIFIER_JSDOC is the
// diagnostic for a missing nullability modifier.
var CheckNullabilityModifiers_MISSING_NULLABILITY_MODIFIER_JSDOC = DiagnosticType_Disabled(
	"JSC_MISSING_NULLABILITY_MODIFIER_JSDOC",
	"{0} is a reference type with no nullability modifier, "+
		"which is disallowed by the style guide.\n"+
		"Please add a '!' to make it explicitly non-nullable, "+
		"or a '?' to make it explicitly nullable.")

// CheckNullabilityModifiers_NULL_MISSING_NULLABILITY_MODIFIER_JSDOC is the
// diagnostic for a missing nullability modifier, where the value is clearly
// nullable.
var CheckNullabilityModifiers_NULL_MISSING_NULLABILITY_MODIFIER_JSDOC = DiagnosticType_Disabled(
	"JSC_NULL_MISSING_NULLABILITY_MODIFIER_JSDOC",
	"{0} is a reference type with no nullability modifier that is explicitly set to null.\n"+
		"Add a '?' to make it explicitly nullable.")

// CheckNullabilityModifiers_REDUNDANT_NULLABILITY_MODIFIER_JSDOC is the
// diagnostic for a redundant nullability modifier.
var CheckNullabilityModifiers_REDUNDANT_NULLABILITY_MODIFIER_JSDOC = DiagnosticType_Disabled(
	"JSC_REDUNDANT_NULLABILITY_MODIFIER_JSDOC",
	"{0} is a non-reference type which is already non-nullable.\n"+
		"Please remove the redundant '!', which is disallowed by the style guide.")

// The set of primitive type names. Note that `void` is a synonym for `undefined`.
var checkNullabilityModifiers_PRIMITIVE_TYPE_NAMES = jutil.NewLinkedHashSet(
	"boolean", "number", "bigint", "string", "symbol", "undefined", "void", "null")

// NewCheckNullabilityModifiers is new CheckNullabilityModifiers(compiler).
//
// Java: CheckNullabilityModifiers(AbstractCompiler)
func NewCheckNullabilityModifiers(compiler AbstractCompiler) *CheckNullabilityModifiers {
	c := &CheckNullabilityModifiers{
		redundantCandidates:   jutil.NewLinkedHashSet[rhino.Node](),
		missingCandidates:     jutil.NewLinkedHashSet[rhino.Node](),
		nullMissingCandidates: jutil.NewLinkedHashSet[rhino.Node](),
		templateTypeNames:     jutil.NewLinkedHashSet[string](),
	}
	c.compiler = compiler
	return c
}

// Java: CheckNullabilityModifiers.process(Node, Node)
func (c *CheckNullabilityModifiers) Process(externs rhino.Node, root rhino.Node) {
	NodeTraversal_TraverseRoots(c.compiler, c, externs, root)
}

// Java: CheckNullabilityModifiers.visit(NodeTraversal, Node, Node)
func (c *CheckNullabilityModifiers) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	info := n.GetJSDocInfo()
	if info != nil {
		c.templateTypeNames.AddAll(info.GetTemplateTypeNames()...)
		for name := range info.GetTypeTransformations().Keys() {
			c.templateTypeNames.Add(name)
		}
		if info.HasType() {
			c.HandleHasType(info.GetType(), n)
		}
		for param := range info.GetParameterNames().All() {
			if info.HasParameterType(param) {
				c.VisitTypeExpression(info.GetParameterType(param), false)
			}
		}
		if info.HasReturnType() {
			c.VisitTypeExpression(info.GetReturnType(), false)
		}
		if info.HasEnumParameterType() {
			// JSDocInfoParser wraps the @enum type in an artificial '!' if it's not a primitive type.
			// Thus a missing modifier warning can never trigger, but a redundant modifier warning can.
			c.VisitTypeExpression(info.GetEnumParameterType(), false)
		}
		if info.HasTypedefType() {
			c.VisitTypeExpression(info.GetTypedefType(), false)
		}
		if info.HasThisType() {
			// JSDocInfoParser wraps the @this type in an artificial '!'. Thus a missing modifier
			// warning can never trigger, but a top-level '!' should be ignored if it precedes a
			// primitive or literal type; otherwise it will trigger a spurious redundant modifier
			// warning.
			c.VisitTypeExpression(info.GetThisType(), true)
		}
		// JSDocInfoParser enforces the @extends and @implements types to be unqualified in the source
		// code, so we don't need to check them.
	}

	if n.IsScript() {
		// If exiting a script node, report applicable warnings and clear the maps.
		c.Report(t)
		c.redundantCandidates.Clear()
		c.missingCandidates.Clear()
		c.nullMissingCandidates.Clear()
		c.templateTypeNames.Clear()
		return
	}
}

// Java: CheckNullabilityModifiers.handleHasType(JSTypeExpression, Node)
func (c *CheckNullabilityModifiers) HandleHasType(expr *rhino.JSTypeExpression, n rhino.Node) {
	// Check if the type is explicitly set to null, if so use NULL_MISSING_NULLABILITY diagnostic.
	if NodeUtil_IsNameDeclOrSimpleAssignLhs(n.GetFirstChild(), n) {
		rValue := NodeUtil_GetRValueOfLValue(n.GetFirstChild())
		if !rValue.IsNil() && rValue.IsNull() {
			c.VisitTypeExpressionWithRValue(expr, false, rValue)
			return
		}
	}
	c.VisitTypeExpression(expr, false)
}

// Java: CheckNullabilityModifiers.report(NodeTraversal)
func (c *CheckNullabilityModifiers) Report(t *NodeTraversal) {
	for n := range c.missingCandidates.All() {
		if c.ShouldReport(n) {
			t.Report(n, CheckNullabilityModifiers_MISSING_NULLABILITY_MODIFIER_JSDOC, CheckNullabilityModifiers_GetReportedTypeName(n))
		}
	}
	for n := range c.nullMissingCandidates.All() {
		if c.ShouldReport(n) {
			t.Report(n, CheckNullabilityModifiers_NULL_MISSING_NULLABILITY_MODIFIER_JSDOC, CheckNullabilityModifiers_GetReportedTypeName(n))
		}
	}
	for n := range c.redundantCandidates.All() {
		if c.ShouldReport(n) {
			// Report on parent so that modifier is also highlighted.
			t.Report(n.GetParent(), CheckNullabilityModifiers_REDUNDANT_NULLABILITY_MODIFIER_JSDOC, CheckNullabilityModifiers_GetReportedTypeName(n))
		}
	}
}

// Java: CheckNullabilityModifiers.shouldReport(Node)
func (c *CheckNullabilityModifiers) ShouldReport(n rhino.Node) bool {
	// Ignore type names that appear in @template clauses in the same source file. In rare cases,
	// this may cause a false negative (if the name is used in a non-template capacity in the same
	// file) or a false positive (if the name is used in a template capacity in a separate file),
	// but it makes this check possible in the absence of type information. If the style guide ever
	// mandates template types (and nothing else) to be all-caps, we can use that assumption to make
	// this check more precise.
	return !n.IsStringLit() || !c.templateTypeNames.Contains(n.GetString())
}

// Java: CheckNullabilityModifiers.visitTypeExpression(JSTypeExpression, boolean)
func (c *CheckNullabilityModifiers) VisitTypeExpression(expr *rhino.JSTypeExpression, hasArtificialTopLevelBang bool) {
	c.VisitTypeExpressionWithRValue(expr, hasArtificialTopLevelBang, rhino.NilNode /* rValue= */)
}

// Java: CheckNullabilityModifiers.visitTypeExpression(JSTypeExpression, boolean, Node)
func (c *CheckNullabilityModifiers) VisitTypeExpressionWithRValue(
	expr *rhino.JSTypeExpression, hasArtificialTopLevelBang bool, rValue rhino.Node) {
	root := expr.GetRoot()
	NodeUtil_VisitPreOrder(
		root,
		NodeUtil_VisitorFunc(func(node rhino.Node) {
			parent := node.GetParent()

			isPrimitiveOrLiteral :=
				CheckNullabilityModifiers_IsPrimitiveType(node) || CheckNullabilityModifiers_IsFunctionLiteral(node) || CheckNullabilityModifiers_IsRecordLiteral(node)
			isReference := CheckNullabilityModifiers_IsReferenceType(node)

			// Whether the node is preceded by a '!' or '?'.
			hasBang := !parent.IsNil() && parent.GetToken() == rhino.Token_BANG
			hasQmark := !parent.IsNil() && parent.GetToken() == rhino.Token_QMARK

			// Whether the node is preceded by a '!' that wasn't artificially added.
			hasNonArtificialBang := hasBang && !(hasArtificialTopLevelBang && parent == root)

			// Whether the node is the type in function(new:T) or function(this:T).
			isNewOrThis := !parent.IsNil() && (parent.IsNew() || parent.IsThis())

			// Whether the node is the type name in a typeof expression.
			isTypeOfType := !parent.IsNil() && parent.IsTypeOf()

			// Whether the node is definitely a possible type for the rValue e.g. 'Type' in @type
			// {Type} or @type {Type|number} but not in @type {!Array<Type>}
			isAppliedToRValue :=
				node == root || (!parent.IsNil() && parent.GetToken() == rhino.Token_PIPE)

			if isReference && !hasBang && !hasQmark && !isNewOrThis && !isTypeOfType {
				if isAppliedToRValue && !rValue.IsNil() && rValue.IsNull() {
					c.nullMissingCandidates.Add(node)
				} else {
					c.missingCandidates.Add(node)
				}
			} else if isPrimitiveOrLiteral && hasNonArtificialBang {
				c.redundantCandidates.Add(node)
			}
		}))
}

// Java: CheckNullabilityModifiers.isPrimitiveType(Node)
func CheckNullabilityModifiers_IsPrimitiveType(node rhino.Node) bool {
	return node.IsStringLit() && checkNullabilityModifiers_PRIMITIVE_TYPE_NAMES.Contains(node.GetString())
}

// Java: CheckNullabilityModifiers.isReferenceType(Node)
func CheckNullabilityModifiers_IsReferenceType(node rhino.Node) bool {
	return node.IsStringLit() && !checkNullabilityModifiers_PRIMITIVE_TYPE_NAMES.Contains(node.GetString())
}

// Java: CheckNullabilityModifiers.isFunctionLiteral(Node)
func CheckNullabilityModifiers_IsFunctionLiteral(node rhino.Node) bool {
	return node.IsFunction()
}

// Java: CheckNullabilityModifiers.isRecordLiteral(Node)
func CheckNullabilityModifiers_IsRecordLiteral(node rhino.Node) bool {
	return node.GetToken() == rhino.Token_LC
}

// Java: CheckNullabilityModifiers.getReportedTypeName(Node)
func CheckNullabilityModifiers_GetReportedTypeName(node rhino.Node) string {
	if CheckNullabilityModifiers_IsFunctionLiteral(node) {
		return "Function"
	}
	if CheckNullabilityModifiers_IsRecordLiteral(node) {
		return "Record literal"
	}
	jutil.CheckState(CheckNullabilityModifiers_IsPrimitiveType(node) || CheckNullabilityModifiers_IsReferenceType(node))
	return node.GetString()
}
