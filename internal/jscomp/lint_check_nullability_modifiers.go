// Port of com/google/javascript/jscomp/lint/CheckNullabilityModifiers.java @ closure-compiler v20261021

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
	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

const (
	nullabilityCandMissing uint8 = 1 << iota
	nullabilityCandNullMissing
	nullabilityCandRedundant
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
	candidateMask         map[rhino.Node]uint8
	redundantCandidates   []rhino.Node
	missingCandidates     []rhino.Node
	nullMissingCandidates []rhino.Node
	templateTypeNames     map[string]bool
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
func checkNullabilityModifiers_isPrimitiveTypeName(name string) bool {
	switch name {
	case "boolean", "number", "bigint", "string", "symbol", "undefined", "void", "null":
		return true
	default:
		return false
	}
}

// NewCheckNullabilityModifiers is new CheckNullabilityModifiers(compiler).
//
// Java: CheckNullabilityModifiers(AbstractCompiler)
func NewCheckNullabilityModifiers(compiler AbstractCompiler) *CheckNullabilityModifiers {
	return &CheckNullabilityModifiers{
		compiler: compiler,
	}
}

func (*CheckNullabilityModifiers) cloneForParallel(c *Compiler) NodeTraversal_Callback {
	return NewCheckNullabilityModifiers(c)
}

// Java: CheckNullabilityModifiers.process(Node, Node)
func (c *CheckNullabilityModifiers) Process(externs rhino.Node, root rhino.Node) {
	NodeTraversal_TraverseRoots(c.compiler, c, externs, root)
}

func (c *CheckNullabilityModifiers) addTemplateTypeName(name string) {
	if c.templateTypeNames == nil {
		c.templateTypeNames = make(map[string]bool)
	}
	c.templateTypeNames[name] = true
}

func (c *CheckNullabilityModifiers) recordCandidate(node rhino.Node, bit uint8) bool {
	if c.candidateMask == nil {
		c.candidateMask = make(map[rhino.Node]uint8)
	}
	if c.candidateMask[node]&bit != 0 {
		return false
	}
	c.candidateMask[node] |= bit
	return true
}

// Java: CheckNullabilityModifiers.visit(NodeTraversal, Node, Node)
func (c *CheckNullabilityModifiers) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	info := n.GetJSDocInfo()
	if info != nil {
		for _, name := range info.GetTemplateTypeNames() {
			c.addTemplateTypeName(name)
		}
		for name := range info.GetTypeTransformations().Keys() {
			c.addTemplateTypeName(name)
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

	if _, ok := ast.AsScript(n); ok {
		// If exiting a script node, report applicable warnings and clear the maps.
		c.Report(t)
		clear(c.candidateMask)
		c.redundantCandidates = c.redundantCandidates[:0]
		c.missingCandidates = c.missingCandidates[:0]
		c.nullMissingCandidates = c.nullMissingCandidates[:0]
		clear(c.templateTypeNames)
	}
}

// Java: CheckNullabilityModifiers.handleHasType(JSTypeExpression, Node)
func (c *CheckNullabilityModifiers) HandleHasType(expr *rhino.JSTypeExpression, n rhino.Node) {
	// Check if the type is explicitly set to null, if so use NULL_MISSING_NULLABILITY diagnostic.
	firstChild := n.GetFirstChild()
	if NodeUtil_IsNameDeclOrSimpleAssignLhs(firstChild, n) {
		rValue := NodeUtil_GetRValueOfLValue(firstChild)
		if _, isNull := ast.AsNull(rValue); isNull {
			c.VisitTypeExpressionWithRValue(expr, false, rValue)
			return
		}
	}
	c.VisitTypeExpression(expr, false)
}

// Java: CheckNullabilityModifiers.report(NodeTraversal)
func (c *CheckNullabilityModifiers) Report(t *NodeTraversal) {
	for _, n := range c.missingCandidates {
		if c.ShouldReport(n) {
			t.Report(n, CheckNullabilityModifiers_MISSING_NULLABILITY_MODIFIER_JSDOC, CheckNullabilityModifiers_GetReportedTypeName(n))
		}
	}
	for _, n := range c.nullMissingCandidates {
		if c.ShouldReport(n) {
			t.Report(n, CheckNullabilityModifiers_NULL_MISSING_NULLABILITY_MODIFIER_JSDOC, CheckNullabilityModifiers_GetReportedTypeName(n))
		}
	}
	for _, n := range c.redundantCandidates {
		if c.ShouldReport(n) {
			// Report on parent so that modifier is also highlighted.
			t.Report(n.GetParent(), CheckNullabilityModifiers_REDUNDANT_NULLABILITY_MODIFIER_JSDOC, CheckNullabilityModifiers_GetReportedTypeName(n))
		}
	}
}

// Java: CheckNullabilityModifiers.shouldReport(Node)
func (c *CheckNullabilityModifiers) ShouldReport(n rhino.Node) bool {
	if len(c.templateTypeNames) == 0 {
		return true
	}
	// Ignore type names that appear in @template clauses in the same source file.
	strLit, ok := ast.AsStringLit(n)
	return !ok || !c.templateTypeNames[strLit.Value()]
}

// Java: CheckNullabilityModifiers.visitTypeExpression(JSTypeExpression, boolean)
func (c *CheckNullabilityModifiers) VisitTypeExpression(expr *rhino.JSTypeExpression, hasArtificialTopLevelBang bool) {
	c.VisitTypeExpressionWithRValue(expr, hasArtificialTopLevelBang, rhino.NilNode /* rValue= */)
}

// Java: CheckNullabilityModifiers.visitTypeExpression(JSTypeExpression, boolean, Node)
func (c *CheckNullabilityModifiers) VisitTypeExpressionWithRValue(
	expr *rhino.JSTypeExpression, hasArtificialTopLevelBang bool, rValue rhino.Node) {
	root := expr.GetRoot()
	_, rValueIsNull := ast.AsNull(rValue)
	NodeUtil_VisitPreOrder(
		root,
		func(node rhino.Node) {
			var isPrimitiveOrLiteral, isReference bool
			switch node.GetToken() {
			case rhino.Token_STRINGLIT:
				if checkNullabilityModifiers_isPrimitiveTypeName(node.GetString()) {
					isPrimitiveOrLiteral = true
				} else {
					isReference = true
				}
			case rhino.Token_FUNCTION, rhino.Token_LC:
				isPrimitiveOrLiteral = true
			default:
				return
			}

			parent := node.GetParent()

			// Whether the node is preceded by a '!' or '?'.
			var parentTok rhino.Token
			hasParent := !parent.IsNil()
			if hasParent {
				parentTok = parent.GetToken()
			}
			hasBang := hasParent && parentTok == rhino.Token_BANG
			hasQmark := hasParent && parentTok == rhino.Token_QMARK

			// Whether the node is preceded by a '!' that wasn't artificially added.
			hasNonArtificialBang := hasBang && !(hasArtificialTopLevelBang && parent == root)

			// Whether the node is the type in function(new:T) or function(this:T).
			isNewOrThis := hasParent && (parentTok == rhino.Token_NEW || parentTok == rhino.Token_THIS)

			// Whether the node is the type name in a typeof expression.
			isTypeOfType := hasParent && parentTok == rhino.Token_TYPEOF

			// Whether the node is definitely a possible type for the rValue e.g. 'Type' in @type
			// {Type} or @type {Type|number} but not in @type {!Array<Type>}
			isAppliedToRValue :=
				node == root || (hasParent && parentTok == rhino.Token_PIPE)

			if isReference && !hasBang && !hasQmark && !isNewOrThis && !isTypeOfType {
				if isAppliedToRValue && rValueIsNull {
					if c.recordCandidate(node, nullabilityCandNullMissing) {
						c.nullMissingCandidates = append(c.nullMissingCandidates, node)
					}
				} else {
					if c.recordCandidate(node, nullabilityCandMissing) {
						c.missingCandidates = append(c.missingCandidates, node)
					}
				}
			} else if isPrimitiveOrLiteral && hasNonArtificialBang {
				if c.recordCandidate(node, nullabilityCandRedundant) {
					c.redundantCandidates = append(c.redundantCandidates, node)
				}
			}
		})
}

// Java: CheckNullabilityModifiers.isPrimitiveType(Node)
func CheckNullabilityModifiers_IsPrimitiveType(node rhino.Node) bool {
	strLit, ok := ast.AsStringLit(node)
	return ok && checkNullabilityModifiers_isPrimitiveTypeName(strLit.Value())
}

// Java: CheckNullabilityModifiers.isReferenceType(Node)
func CheckNullabilityModifiers_IsReferenceType(node rhino.Node) bool {
	strLit, ok := ast.AsStringLit(node)
	return ok && !checkNullabilityModifiers_isPrimitiveTypeName(strLit.Value())
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
	strLit, ok := ast.AsStringLit(node)
	jutil.CheckState(ok)
	return strLit.Value()
}
