// Port of com/google/javascript/jscomp/J2clEqualitySameRewriterPass.java @ closure-compiler v20260415

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

package jscomp

import (
	"math"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// J2clEqualitySameRewriterPass is an optimization pass to re-write J2CL
// Equality.$same.
//
// Java: public class J2clEqualitySameRewriterPass extends AbstractPeepholeOptimization
type J2clEqualitySameRewriterPass struct {
	AbstractPeepholeOptimizationBase
	useTypes            bool
	shouldRunJ2clPasses bool
}

// Java: J2clEqualitySameRewriterPass(boolean)
func NewJ2clEqualitySameRewriterPass(useTypes bool) *J2clEqualitySameRewriterPass {
	p := &J2clEqualitySameRewriterPass{}
	p.InitAbstractPeepholeOptimizationBase(p)
	p.useTypes = useTypes
	return p
}

// Java: J2clEqualitySameRewriterPass.beginTraversal(AbstractCompiler)
func (p *J2clEqualitySameRewriterPass) BeginTraversal(compiler AbstractCompiler) {
	p.AbstractPeepholeOptimizationBase.BeginTraversal(compiler)
	p.shouldRunJ2clPasses = J2clSourceFileChecker_ShouldRunJ2clPasses(compiler)
}

// Java: J2clEqualitySameRewriterPass.optimizeSubtree(Node)
func (p *J2clEqualitySameRewriterPass) OptimizeSubtree(node *rhino.Node) *rhino.Node {
	if !p.shouldRunJ2clPasses {
		return node
	}

	// Do not optimize if any parameters were removed during optimizations
	if !node.HasXChildren(3) {
		return node
	}

	var replacement *rhino.Node
	if j2clEqualitySameRewriterPass_isStringEqualsMethod(node) {
		replacement = p.trySubstituteStringEquals(node)
	} else if j2clEqualitySameRewriterPass_isEqualitySameCall(node) {
		replacement = p.trySubstituteEqualitySame(node)
	}

	if replacement == nil {
		return node
	}

	replacement.SrcrefIfMissing(node)
	node.ReplaceWith(replacement)
	p.ReportChangeToEnclosingScope(replacement)
	return replacement
}

// Java: J2clEqualitySameRewriterPass.trySubstituteStringEquals(Node)
func (p *J2clEqualitySameRewriterPass) trySubstituteStringEquals(callNode *rhino.Node) *rhino.Node {
	firstExprValue := j2clEqualitySameRewriterPass_getKnownLiteralValue(callNode.GetSecondChild())
	if firstExprValue == J2clEqualitySameRewriterPass_NodeValue_UNKNOWN ||
		firstExprValue == J2clEqualitySameRewriterPass_NodeValue_NULL_OR_UNDEFINED {
		// Potential NPE, don't optimize.
		return nil
	}

	return p.trySubstituteEqualitySame(callNode)
}

// Java: J2clEqualitySameRewriterPass.trySubstituteEqualitySame(Node)
func (p *J2clEqualitySameRewriterPass) trySubstituteEqualitySame(callNode *rhino.Node) *rhino.Node {
	firstExpr := callNode.GetSecondChild()
	firstExprValue := j2clEqualitySameRewriterPass_getKnownLiteralValue(firstExpr)
	secondExpr := callNode.GetLastChild()
	secondExprValue := j2clEqualitySameRewriterPass_getKnownLiteralValue(secondExpr)

	if firstExprValue == J2clEqualitySameRewriterPass_NodeValue_UNKNOWN &&
		secondExprValue == J2clEqualitySameRewriterPass_NodeValue_UNKNOWN {
		return nil
	}

	if firstExprValue == J2clEqualitySameRewriterPass_NodeValue_NULL_OR_UNDEFINED {
		return p.rewriteNullCheck(secondExpr, firstExpr)
	}

	if secondExprValue == J2clEqualitySameRewriterPass_NodeValue_NULL_OR_UNDEFINED {
		return p.rewriteNullCheck(firstExpr, secondExpr)
	}

	if firstExprValue == J2clEqualitySameRewriterPass_NodeValue_NON_NULL ||
		secondExprValue == J2clEqualitySameRewriterPass_NodeValue_NON_NULL {
		// There is a coercion danger (e.g. 0 == null) but since at least one side is not null, we can
		// safely use === that will not trigger any coercion.
		return p.rewriteAsStrictEq(firstExpr, secondExpr)
	}

	jutil.CheckState(firstExprValue == J2clEqualitySameRewriterPass_NodeValue_NUMBER ||
		secondExprValue == J2clEqualitySameRewriterPass_NodeValue_NUMBER)
	return p.rewriteNumberCheck(firstExpr, secondExpr)
}

// Java: J2clEqualitySameRewriterPass.rewriteNullCheck(Node, Node)
func (p *J2clEqualitySameRewriterPass) rewriteNullCheck(expr *rhino.Node, nullExpression *rhino.Node) *rhino.Node {
	expr.Detach()
	nullExpression.Detach()
	if p.useTypes && p.canOnlyBeObject(expr) {
		return rhino.IR_Not(expr)
	}
	// At least one side is null or undefined so no coercion danger with ==.
	return rhino.IR_Eq(expr, nullExpression)
}

// Java: J2clEqualitySameRewriterPass.canOnlyBeObject(Node)
func (p *J2clEqualitySameRewriterPass) canOnlyBeObject(n *rhino.Node) bool {
	color := n.GetColor()
	// Safe as long as the color is a) not the UNKNOWN native color and b) not any primitive
	if color == nil {
		return false
	}

	if color.IsUnion() {
		// ignore null/undefined
		color = color.SubtractNullOrVoid()
	}

	// In theory we could allow unions of multiple objects here
	return !color.IsUnion() && !color.IsPrimitive() && !color.Equals(colors.StandardColors_UNKNOWN)
}

// Java: J2clEqualitySameRewriterPass.rewriteAsStrictEq(Node, Node)
func (p *J2clEqualitySameRewriterPass) rewriteAsStrictEq(firstExpr *rhino.Node, secondExpr *rhino.Node) *rhino.Node {
	firstExpr.Detach()
	secondExpr.Detach()
	return rhino.IR_Sheq(firstExpr, secondExpr)
}

// Java: J2clEqualitySameRewriterPass.rewriteNumberCheck(Node, Node)
func (p *J2clEqualitySameRewriterPass) rewriteNumberCheck(firstExpr *rhino.Node, secondExpr *rhino.Node) *rhino.Node {
	firstValue := NodeUtil_GetNumberValue(firstExpr)
	secondValue := NodeUtil_GetNumberValue(secondExpr)

	if firstValue != nil && secondValue != nil {
		firstExpr.Detach()
		secondExpr.Detach()
		// Double.equals compares doubleToLongBits.
		return NodeUtil_BooleanNode(jutil.DoubleToLongBits(*firstValue) == jutil.DoubleToLongBits(*secondValue))
	}

	if j2clEqualitySameRewriterPass_isSafeNumber(firstValue) || j2clEqualitySameRewriterPass_isSafeNumber(secondValue) {
		// Since one side is not 0, -0 or NaN, there is no risk of -0 vs 0 or NaN vs NaN comparison.
		return p.rewriteAsStrictEq(firstExpr, secondExpr)
	}

	if p.useTypes && (p.canOnlyBeObject(firstExpr) || p.canOnlyBeObject(secondExpr)) {
		// Since one side is not number, there is no risk of -0 vs 0 or NaN vs NaN comparision.
		return p.rewriteAsStrictEq(firstExpr, secondExpr)
	}

	return nil
}

// Java: J2clEqualitySameRewriterPass.isSafeNumber(Double)
func j2clEqualitySameRewriterPass_isSafeNumber(d *float64) bool {
	return d != nil && *d != 0 && !math.IsNaN(*d)
}

// Java: private enum J2clEqualitySameRewriterPass.NodeValue
type J2clEqualitySameRewriterPass_NodeValue int32

const (
	J2clEqualitySameRewriterPass_NodeValue_UNKNOWN J2clEqualitySameRewriterPass_NodeValue = iota
	J2clEqualitySameRewriterPass_NodeValue_NULL_OR_UNDEFINED
	// JavaScript number. Needs special treatment to preserve Object.is semantics.
	J2clEqualitySameRewriterPass_NodeValue_NUMBER
	// Non-null and also known to be not a number.
	J2clEqualitySameRewriterPass_NodeValue_NON_NULL
)

var j2clEqualitySameRewriterPass_NodeValue_names = [...]string{
	"UNKNOWN", "NULL_OR_UNDEFINED", "NUMBER", "NON_NULL",
}

// J2clEqualitySameRewriterPass_NodeValue_values is NodeValue.values().
func J2clEqualitySameRewriterPass_NodeValue_values() []J2clEqualitySameRewriterPass_NodeValue {
	return []J2clEqualitySameRewriterPass_NodeValue{
		J2clEqualitySameRewriterPass_NodeValue_UNKNOWN,
		J2clEqualitySameRewriterPass_NodeValue_NULL_OR_UNDEFINED,
		J2clEqualitySameRewriterPass_NodeValue_NUMBER,
		J2clEqualitySameRewriterPass_NodeValue_NON_NULL,
	}
}

// String is NodeValue.name().
func (v J2clEqualitySameRewriterPass_NodeValue) String() string {
	return j2clEqualitySameRewriterPass_NodeValue_names[v]
}

// Java: J2clEqualitySameRewriterPass.getKnownLiteralValue(Node)
func j2clEqualitySameRewriterPass_getKnownLiteralValue(n *rhino.Node) J2clEqualitySameRewriterPass_NodeValue {
	switch NodeUtil_GetKnownValueType(n) {
	case NodeUtil_ValueType_VOID:
		if NodeUtil_CanBeSideEffected(n) {
			return J2clEqualitySameRewriterPass_NodeValue_UNKNOWN
		}
		return J2clEqualitySameRewriterPass_NodeValue_NULL_OR_UNDEFINED
	case NodeUtil_ValueType_NULL:
		return J2clEqualitySameRewriterPass_NodeValue_NULL_OR_UNDEFINED
	case NodeUtil_ValueType_STRING, NodeUtil_ValueType_BOOLEAN, NodeUtil_ValueType_OBJECT, NodeUtil_ValueType_BIGINT:
		return J2clEqualitySameRewriterPass_NodeValue_NON_NULL
	case NodeUtil_ValueType_NUMBER:
		return J2clEqualitySameRewriterPass_NodeValue_NUMBER
	case NodeUtil_ValueType_UNDETERMINED:
		return J2clEqualitySameRewriterPass_NodeValue_UNKNOWN
	}
	// Go-only: language-adapter — unreachable for a valid enum value: the
	// Java switch expression is exhaustive.
	panic(jutil.NewIllegalStateException())
}

// Java: J2clEqualitySameRewriterPass.isEqualitySameCall(Node)
func j2clEqualitySameRewriterPass_isEqualitySameCall(node *rhino.Node) bool {
	return node.IsCall() && j2clEqualitySameRewriterPass_hasName(node.GetFirstChild(), "Equality", "$same")
}

// Java: J2clEqualitySameRewriterPass.isStringEqualsMethod(Node)
func j2clEqualitySameRewriterPass_isStringEqualsMethod(node *rhino.Node) bool {
	return node.IsCall() &&
		j2clEqualitySameRewriterPass_hasName(
			node.GetFirstChild(),
			"String",
			"m_equals__java_lang_String__java_lang_Object__boolean")
}

// Java: J2clEqualitySameRewriterPass.hasName(Node, String, String)
func j2clEqualitySameRewriterPass_hasName(fnName *rhino.Node, className string, methodName string) bool {
	if !fnName.IsQualifiedName() {
		return false
	}
	// NOTE: This should be rewritten to use method name + file name of definition site
	// like other J2CL passes, which is more precise.
	originalQname := *fnName.GetOriginalQualifiedName()
	return strings.HasSuffix(originalQname, methodName) && strings.Contains(originalQname, className)
}
