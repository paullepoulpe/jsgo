// Port of com/google/javascript/jscomp/CheckSuspiciousCode.java @ closure-compiler v20261021

/*
 * Copyright 2012 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Go-only: this pass uses the typed AST view (internal/ast): each check
// takes the kind it reads (the branches of an if, the body of a loop, the
// operands of a comparison or a logical operator) instead of testing the
// token and reading children by position. The optional chain check stays on
// rhino: it reads the node's parent, of any kind.

// CheckSuspiciousCode checks for common errors, such as misplaced semicolons:
//
//	if (x); act_now();
//
// or comparison against NaN:
//
//	if (x === NaN) act();
//
// and generates warnings.
type CheckSuspiciousCode struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	// Go-only: language-adapter — distinct identity (the Java class has no
	// fields; pointers to zero-size Go structs may compare equal).
	_ byte
}

// Java: CheckSuspiciousCode.SUSPICIOUS_SEMICOLON
var CheckSuspiciousCode_SUSPICIOUS_SEMICOLON = DiagnosticType_Warning(
	"JSC_SUSPICIOUS_SEMICOLON",
	"If this if/for/while really shouldn''t have a body, use '{}'")

// Java: CheckSuspiciousCode.SUSPICIOUS_COMPARISON_WITH_NAN
var CheckSuspiciousCode_SUSPICIOUS_COMPARISON_WITH_NAN = DiagnosticType_Warning(
	"JSC_SUSPICIOUS_NAN", "Comparison against NaN is always false. Did you mean isNaN()?")

// Java: CheckSuspiciousCode.SUSPICIOUS_IN_OPERATOR
var CheckSuspiciousCode_SUSPICIOUS_IN_OPERATOR = DiagnosticType_Warning(
	"JSC_SUSPICIOUS_IN",
	"Use of the \"in\" keyword on non-object types throws an exception.")

// Java: CheckSuspiciousCode.SUSPICIOUS_INSTANCEOF_LEFT_OPERAND
var CheckSuspiciousCode_SUSPICIOUS_INSTANCEOF_LEFT_OPERAND = DiagnosticType_Warning(
	"JSC_SUSPICIOUS_INSTANCEOF_LEFT",
	"\"instanceof\" with left non-object operand is always false.")

// Java: CheckSuspiciousCode.SUSPICIOUS_LEFT_OPERAND_OF_LOGICAL_OPERATOR
var CheckSuspiciousCode_SUSPICIOUS_LEFT_OPERAND_OF_LOGICAL_OPERATOR = DiagnosticType_Warning(
	"JSC_SUSPICIOUS_LEFT_OPERAND_OF_LOGICAL_OPERATOR",
	"Left operand of {0} operator is always {1}.")

// Java: CheckSuspiciousCode.SUSPICIOUS_NEGATED_LEFT_OPERAND_OF_IN_OPERATOR
var CheckSuspiciousCode_SUSPICIOUS_NEGATED_LEFT_OPERAND_OF_IN_OPERATOR = DiagnosticType_Warning(
	"JSC_SUSPICIOUS_NEGATED_LEFT_OPERAND_OF_IN_OPERATOR",
	"Suspicious negated left operand of 'in' operator.")

// Java: CheckSuspiciousCode.SUSPICIOUS_BREAKING_OUT_OF_OPTIONAL_CHAIN
var CheckSuspiciousCode_SUSPICIOUS_BREAKING_OUT_OF_OPTIONAL_CHAIN = DiagnosticType_Warning(
	"SUSPICIOUS_BREAKING_OUT_OF_OPTIONAL_CHAIN",
	"Suspicious breaking out of optional chain. May result in TypeError if optional chain is"+
		" undefined.")

// Java: CheckSuspiciousCode()
func NewCheckSuspiciousCode() *CheckSuspiciousCode {
	return &CheckSuspiciousCode{}
}

var checkSuspiciousCodeVisitTokens = func() ast.TokenSet {
	var s ast.TokenSet
	for _, t := range [...]rhino.Token{
		rhino.Token_IF,
		rhino.Token_WHILE, rhino.Token_FOR, rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF,
		rhino.Token_EQ, rhino.Token_GE, rhino.Token_GT, rhino.Token_LE, rhino.Token_LT, rhino.Token_NE,
		rhino.Token_SHEQ, rhino.Token_SHNE,
		rhino.Token_IN,
		rhino.Token_INSTANCEOF,
		rhino.Token_OR, rhino.Token_AND,
		rhino.Token_OPTCHAIN_GETPROP, rhino.Token_OPTCHAIN_GETELEM, rhino.Token_OPTCHAIN_CALL,
	} {
		s.Add(t)
	}
	return s
}()

func (*CheckSuspiciousCode) visitTokens() *ast.TokenSet        { return &checkSuspiciousCodeVisitTokens }
func (*CheckSuspiciousCode) combinedVisitTokens() ast.TokenSet { return checkSuspiciousCodeVisitTokens }

// Java: CheckSuspiciousCode.visit(NodeTraversal, Node, Node)
//
// Go-only: Visit switches on the token once and hands each check the kind
// it reads; Java calls every check on every node and each check tests the
// token. A node has one token, so the checks run in the same order (for
// `in`: the right operand first, then the negated left one).
func (c *CheckSuspiciousCode) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_IF:
		c.CheckMissingSemicolon(t, ast.MustIf(n))
	case rhino.Token_WHILE, rhino.Token_FOR, rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF:
		CheckSuspiciousCode_ReportIfWasEmpty(t, ast.MustLoop(n).Body())
	case rhino.Token_EQ, rhino.Token_GE, rhino.Token_GT, rhino.Token_LE, rhino.Token_LT, rhino.Token_NE,
		rhino.Token_SHEQ, rhino.Token_SHNE:
		c.CheckNaN(t, ast.MustBinary(n))
	case rhino.Token_IN:
		in := ast.MustBinary(n)
		c.CheckInvalidIn(t, in)
		c.CheckNegatedLeftOperandOfInOperator(t, in)
	case rhino.Token_INSTANCEOF:
		c.CheckNonObjectInstanceOf(t, ast.MustBinary(n))
	case rhino.Token_OR, rhino.Token_AND:
		c.CheckLeftOperandOfLogicalOperator(t, ast.MustBinary(n))
	case rhino.Token_OPTCHAIN_GETPROP, rhino.Token_OPTCHAIN_GETELEM, rhino.Token_OPTCHAIN_CALL:
		CheckSuspiciousCode_CheckSuspiciousBreakingOutOfOptionalChain(t, n)
	}
}

// CheckMissingSemicolon reports an if whose branch is a lone semicolon.
//
// Go-only: it takes the if; Visit checks the body of a loop directly.
//
// Java: CheckSuspiciousCode.checkMissingSemicolon(NodeTraversal, Node)
func (c *CheckSuspiciousCode) CheckMissingSemicolon(t *NodeTraversal, ifs ast.If) {
	CheckSuspiciousCode_ReportIfWasEmpty(t, ifs.Then())
	if elseCase, ok := ifs.Else(); ok {
		CheckSuspiciousCode_ReportIfWasEmpty(t, elseCase)
	}
}

// Java: CheckSuspiciousCode.reportIfWasEmpty(NodeTraversal, Node)
func CheckSuspiciousCode_ReportIfWasEmpty(t *NodeTraversal, block ast.Block) {
	// A semicolon is distinguished from a block without children by
	// annotating it with EMPTY_BLOCK.  Blocks without children are
	// usually intentional, especially with loops.
	if !block.HasStatements() && block.Node().IsAddedBlock() {
		t.GetCompiler().Report(JSError_MakeWithNode(block.Node(), CheckSuspiciousCode_SUSPICIOUS_SEMICOLON))
	}
}

// Java: CheckSuspiciousCode.checkNaN(NodeTraversal, Node)
func (c *CheckSuspiciousCode) CheckNaN(t *NodeTraversal, comparison ast.Binary) {
	CheckSuspiciousCode_ReportIfNaN(t, comparison, comparison.Left())
	CheckSuspiciousCode_ReportIfNaN(t, comparison, comparison.Right())
}

// Go-only: it takes the comparison it reports on, which Java reads as the
// operand's parent.
//
// Java: CheckSuspiciousCode.reportIfNaN(NodeTraversal, Node)
func CheckSuspiciousCode_ReportIfNaN(t *NodeTraversal, comparison ast.Binary, operand ast.Expr) {
	if NodeUtil_IsNaN(operand.Node()) {
		t.GetCompiler().Report(JSError_MakeWithNode(comparison.Node(), CheckSuspiciousCode_SUSPICIOUS_COMPARISON_WITH_NAN))
	}
}

// Java: CheckSuspiciousCode.checkInvalidIn(NodeTraversal, Node)
func (c *CheckSuspiciousCode) CheckInvalidIn(t *NodeTraversal, in ast.Binary) {
	CheckSuspiciousCode_ReportIfNonObject(t, in, in.Right(), CheckSuspiciousCode_SUSPICIOUS_IN_OPERATOR)
}

// Java: CheckSuspiciousCode.checkNonObjectInstanceOf(NodeTraversal, Node)
func (c *CheckSuspiciousCode) CheckNonObjectInstanceOf(t *NodeTraversal, instanceOf ast.Binary) {
	CheckSuspiciousCode_ReportIfNonObject(t, instanceOf, instanceOf.Left(), CheckSuspiciousCode_SUSPICIOUS_INSTANCEOF_LEFT_OPERAND)
}

// Go-only: it takes the operator it reports on, which Java reads as the
// operand's parent.
//
// Java: CheckSuspiciousCode.reportIfNonObject(NodeTraversal, Node, DiagnosticType)
func CheckSuspiciousCode_ReportIfNonObject(t *NodeTraversal, operator ast.Binary, operand ast.Expr, diagnosticType *DiagnosticType) bool {
	if operand.Token() == rhino.Token_ADD || !NodeUtil_MayBeObject(operand.Node()) {
		t.Report(operator.Node(), diagnosticType)
		return true
	}
	return false
}

// Java: CheckSuspiciousCode.checkNegatedLeftOperandOfInOperator(NodeTraversal, Node)
func (c *CheckSuspiciousCode) CheckNegatedLeftOperandOfInOperator(t *NodeTraversal, in ast.Binary) {
	if not, ok := in.Left().AsUnaryOf(rhino.Token_NOT); ok {
		t.Report(not.Node(), CheckSuspiciousCode_SUSPICIOUS_NEGATED_LEFT_OPERAND_OF_IN_OPERATOR)
	}
}

// CheckLeftOperandOfLogicalOperator checks for the LHS of a logical operator
// (&& and ||) being deterministically truthy or falsy, using both syntactic
// and type information. This is always suspicious (though for different
// reasons: "truthy and" means the LHS is always ignored and should be removed,
// "falsy and" means the RHS is dead code, and vice versa for "or". However,
// there are a number of legitimate use cases where we need to back off from
// using type information (see GetBooleanValueWithTypes for more details on
// these back offs).
//
// Java: CheckSuspiciousCode.checkLeftOperandOfLogicalOperator(NodeTraversal, Node)
func (c *CheckSuspiciousCode) CheckLeftOperandOfLogicalOperator(t *NodeTraversal, op ast.Binary) {
	operator := "&&"
	if op.Token() == rhino.Token_OR {
		operator = "||"
	}
	v := c.GetBooleanValueWithTypes(op.Left())
	if v != base.Tri_UNKNOWN {
		result := "falsy"
		if v == base.Tri_TRUE {
			result = "truthy"
		}
		t.Report(op.Node(), CheckSuspiciousCode_SUSPICIOUS_LEFT_OPERAND_OF_LOGICAL_OPERATOR, operator, result)
	}
}

// CheckSuspiciousCode_CheckSuspiciousBreakingOutOfOptionalChain checks for
// breaking out of optional chain. (e.g. `(a?.b).c)` There is no reason to
// break out of optional chains and doing so may cause a TypeError if `a` is
// nullish. Using `a?.b.c` is safer and will have the same effect when `a` is
// not nullish.
//
// Java: CheckSuspiciousCode.checkSuspiciousBreakingOutOfOptionalChain(NodeTraversal, Node)
func CheckSuspiciousCode_CheckSuspiciousBreakingOutOfOptionalChain(t *NodeTraversal, n rhino.Node) {
	if NodeUtil_IsOptChainNode(n) {
		parent := n.GetParent()
		if n.IsFirstChildOf(parent) &&
			(parent.IsGetProp() || parent.IsGetElem() || parent.IsCall()) {
			t.Report(n, CheckSuspiciousCode_SUSPICIOUS_BREAKING_OUT_OF_OPTIONAL_CHAIN)
		}
	}
}

// GetBooleanValueWithTypes returns the possible boolean values of a node.
// This is a combination of NodeUtil.getBooleanValue and
// JSType.getPossibleToBooleanOutcomes, with some additional backoff for
// situations that are known to be less accurate (to wit, qualified names and
// truthy return types, which return UNKNOWN even if the type information
// seems conclusive). Another difference from the NodeUtil and JSTyoe methods
// is that we include some specific optimization to avoid quadratic behavior
// of large nested logical operations.
//
// The specifics of when this method returns UNKNOWN instead of TRUE or FALSE
// is as follows:
//
//   - We should not determine that a qualified name (expanded slightly to
//     include computed properties) is always truthy or always falsy. There
//     are many cases where an extern defines a name to be truthy, but it
//     still makes sense to feature-test in the browser.
//   - We should not determine that a getelem or a function call is always
//     truthy (though we expand it to simply never complain about
//     always-truthy based only on types) since the standard externs lie
//     about the return type of Map.prototype.get and array accesses, which
//     can both return undefined despite what the externs say.
//
// We do not back off from boolean literals (e.g. "true &&"), though they
// appear to be common in generated code. Instead, such code should suppress
// "suspiciousCode". We also do not back off from always-falsy function call
// results, since it provides a valuable check and lies in this direction are
// much less common.
//
// Java: CheckSuspiciousCode.getBooleanValueWithTypes(Node)
func (c *CheckSuspiciousCode) GetBooleanValueWithTypes(e ast.Expr) base.Tri {
	n := e.Node()
	switch n.GetToken() {
	case rhino.Token_ASSIGN:
		return c.GetBooleanValueWithTypes(ast.MustAssign(n).Value())
	case rhino.Token_COMMA:
		return c.GetBooleanValueWithTypes(ast.MustBinary(n).Right())
	case rhino.Token_NOT:
		return c.GetBooleanValueWithTypes(ast.MustUnary(n).Operand()).Not()
	case rhino.Token_AND:
		// Assume the left-hand side is unknown. If it's not then we'll report it elsewhere. This
		// prevents revisiting deeper nodes repeatedly, which would result in O(n^2) performance.
		return base.Tri_UNKNOWN.And(c.GetBooleanValueWithTypes(ast.MustBinary(n).Right()))
	case rhino.Token_OR:
		// Assume the left-hand side is unknown. If it's not then we'll report it elsewhere. This
		// prevents revisiting deeper nodes repeatedly, which would result in O(n^2) performance.
		return base.Tri_UNKNOWN.Or(c.GetBooleanValueWithTypes(ast.MustBinary(n).Right()))
	case rhino.Token_HOOK:
		hook := ast.MustHook(n)
		trueValue := c.GetBooleanValueWithTypes(hook.Then())
		falseValue := c.GetBooleanValueWithTypes(hook.Else())
		if trueValue == falseValue {
			return trueValue
		}
		return base.Tri_UNKNOWN
	case rhino.Token_FUNCTION, rhino.Token_CLASS, rhino.Token_NEW, rhino.Token_ARRAYLIT, rhino.Token_OBJECTLIT:
		return base.Tri_TRUE
	case rhino.Token_VOID:
		return base.Tri_FALSE
	case rhino.Token_GETPROP, rhino.Token_GETELEM, rhino.Token_OPTCHAIN_GETELEM, rhino.Token_OPTCHAIN_GETPROP:
		// Assume that type information on getprops and getelems are likely to be wrong.  This
		// prevents spurious warnings from not including undefined in getelem's return value,
		// from existence checks of symbols the externs define as certainly true, or from default
		// initialization of globals ({@code x.y = x.y || {}}).
		return base.Tri_UNKNOWN
	default:
	}
	// If we reach this point then all the composite structures that we can decompose have
	// already been handled, leaving only qualified names and type-aware checks to handle below.
	// Note that much of the switch above in fact duplicates the logic in getImpureBooleanValue,
	// though with some subtle differences.  Important differences include (1) avoiding recursion
	// into the left-hand-side of nested logical operators, instead treating them as unknown since
	// they would have already been reported elsewhere in the traversal had they been otherwise
	// (this guarantees we visit each node once, rather than quadratically repeating work);
	// (2) it propagates our unique amalgam of syntax-based and type-based checks to work when more
	// deeply nested (i.e. recursively).  These differences rely on assumptions that are very
	// specific to this use case, so it does not make sense to upstream them.
	literalValue := NodeUtil_GetBooleanValue(n)
	if literalValue != base.Tri_UNKNOWN || n.IsName() {
		// If the truthiness is determinstic from the syntax then return that immediately.
		// Alternatively, NAME nodes also get a pass since we don't trust the type information.
		return literalValue
	}
	type_ := n.GetJSType()
	if type_ != nil {
		// Distrust types we think are always truthy, since sometimes the types lie, even for results
		// of function calls (e.g. Map.prototype.get), so it's still important to check.  But
		// always-falsy values are a little more obviously wrong and there should be no reason for
		// those type annotations to be lies.  ANDing with UNKNOWN ensures we never return TRUE.
		return base.Tri_UNKNOWN.And(type_.GetPossibleToBooleanOutcomes().ToTri())
	}
	return base.Tri_UNKNOWN
}
