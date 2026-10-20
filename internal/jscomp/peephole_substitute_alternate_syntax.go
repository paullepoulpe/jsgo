// Port of com/google/javascript/jscomp/PeepholeSubstituteAlternateSyntax.java @ closure-compiler v20261015

/*
 * Copyright 2010 The Closure Compiler Authors.
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
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// PeepholeSubstituteAlternateSyntax is a peephole optimization that
// minimizes code by simplifying conditional expressions, replacing IFs with
// HOOKs, replacing object constructors with literals, and simplifying
// returns.
type PeepholeSubstituteAlternateSyntax struct {
	AbstractPeepholeOptimizationBase
	late bool
}

var peepholeSubstituteAlternateSyntax_STRING_SPLIT_OVERHEAD = jutil.StringLength(".split('.')")

// NewPeepholeSubstituteAlternateSyntax creates the optimization.
//
// late: When late is false, this mean we are currently running before most
// of the other optimizations. In this case we would avoid optimizations that
// would make the code harder to analyze (such as using string splitting,
// merging statements with commas, etc). When this is true, we would do
// anything to minimize for size.
//
// Java: PeepholeSubstituteAlternateSyntax(boolean)
func NewPeepholeSubstituteAlternateSyntax(late bool) *PeepholeSubstituteAlternateSyntax {
	p := &PeepholeSubstituteAlternateSyntax{}
	p.InitAbstractPeepholeOptimizationBase(p) // the implicit super()
	p.late = late
	return p
}

// OptimizeSubtree tries apply our various peephole minimizations on the
// passed in node.
//
// Java: PeepholeSubstituteAlternateSyntax.optimizeSubtree(Node)
func (p *PeepholeSubstituteAlternateSyntax) OptimizeSubtree(node *rhino.Node) *rhino.Node {
	switch node.GetToken() {
	case rhino.Token_ASSIGN_SUB:
		return p.ReduceSubstractionAssignment(node)
	case rhino.Token_TRUE, rhino.Token_FALSE:
		return p.ReduceTrueFalse(node)
	case rhino.Token_NEW:
		result := p.TryFoldStandardConstructors(node)
		if !result.IsCall() {
			return result
		}
		// tryFoldStandardConstructors() may convert a NEW node into a CALL node
		return p.TryReduceCall(result)
	case rhino.Token_CALL:
		return p.TryReduceCall(node)
	case rhino.Token_RETURN:
		return p.TryReduceReturn(node)
	case rhino.Token_EXPR_RESULT:
		if node.GetFirstChild().IsComma() {
			return p.TrySplitComma(node)
		}
		return node
	case rhino.Token_NAME:
		return p.TryReplaceUndefined(node)
	case rhino.Token_ARRAYLIT:
		return p.TryMinimizeArrayLiteral(node)
	case rhino.Token_GETPROP:
		return p.TryMinimizeWindowRefs(node)
	case rhino.Token_TEMPLATELIT:
		return p.TryTurnTemplateStringsToStrings(node)
	case rhino.Token_AND, rhino.Token_OR, rhino.Token_BITOR, rhino.Token_BITXOR, rhino.Token_BITAND, rhino.Token_COALESCE:
		return p.TryRotateAssociativeOperator(node)
	case rhino.Token_MUL:
		if p.MayHaveSideEffects(node) {
			return node
		}
		return p.TryRotateCommutativeOperator(node)
	default:
		return node
	}
}

var peepholeSubstituteAlternateSyntax_BUILTIN_EXTERNS = jutil.NewLinkedHashSet(
	"Object",
	"Array",
	"Error",
	"RegExp",
	"Math")

// Java: PeepholeSubstituteAlternateSyntax.tryMinimizeWindowRefs(Node)
func (p *PeepholeSubstituteAlternateSyntax) TryMinimizeWindowRefs(node *rhino.Node) *rhino.Node {
	// Normalization needs to be done to ensure there's no shadowing.
	if !p.IsASTNormalized() {
		return node
	}

	jutil.CheckArgument(node.IsGetProp())

	if node.GetFirstChild().IsName() {
		nameNode := node.GetFirstChild()

		// Since normalization has run we know we're referring to the global window.
		if "window" == nameNode.GetString() && peepholeSubstituteAlternateSyntax_BUILTIN_EXTERNS.Contains(node.GetString()) {
			newNameNode := rhino.IR_Name(node.GetString())
			parentNode := node.GetParent()

			newNameNode.Srcref(node)
			node.ReplaceWith(newNameNode)

			if parentNode.IsCall() || parentNode.IsOptChainCall() {
				// e.g. when converting `window.Array?.()` to `Array?.()`, ensure that the
				// OPTCHAIN_CALL gets marked as `FREE_CALL`
				parentNode.PutBooleanProp(rhino.Node_FREE_CALL, true)
			}
			p.ReportChangeToEnclosingScope(parentNode)
			return newNameNode
		}
	}

	return node
}

// Java: PeepholeSubstituteAlternateSyntax.tryRotateCommutativeOperator(Node)
func (p *PeepholeSubstituteAlternateSyntax) TryRotateCommutativeOperator(n *rhino.Node) *rhino.Node {
	if !p.late {
		return n
	}
	// Transform a * (b / c) to b / c * a
	rhs := n.GetLastChild()
	lhs := n.GetFirstChild()
	for lhs.GetToken() == n.GetToken() && NodeUtil_IsAssociative(n.GetToken()) {
		lhs = lhs.GetFirstChild()
	}
	precedence := NodeUtil_Precedence(n.GetToken())
	lhsPrecedence := NodeUtil_Precedence(lhs.GetToken())
	rhsPrecedence := NodeUtil_Precedence(rhs.GetToken())
	if rhsPrecedence == precedence && lhsPrecedence != precedence {
		rhs.Detach()
		lhs.ReplaceWith(rhs)
		n.AddChildToBack(lhs)
		p.ReportChangeToEnclosingScope(n)
	}
	return n
}

// Java: PeepholeSubstituteAlternateSyntax.tryRotateAssociativeOperator(Node)
func (p *PeepholeSubstituteAlternateSyntax) TryRotateAssociativeOperator(n *rhino.Node) *rhino.Node {
	if !p.late {
		return n
	}
	// All commutative operators are also associative
	jutil.CheckArgument(NodeUtil_IsAssociative(n.GetToken()))
	rhs := n.GetLastChild()
	if n.GetToken() == rhs.GetToken() {
		// Transform a * (b * c) to a * b * c
		first := n.RemoveFirstChild()
		second := rhs.RemoveFirstChild()
		third := rhs.GetLastChild().Detach()
		newLhs := rhino.NewNodeWith2Children(n.GetToken(), first, second).SrcrefIfMissing(n)
		newRoot := rhino.NewNodeWith2Children(rhs.GetToken(), newLhs, third).SrcrefIfMissing(rhs)
		n.ReplaceWith(newRoot)
		p.ReportChangeToEnclosingScope(newRoot)
		return newRoot
	} else if NodeUtil_IsCommutative(n.GetToken()) && !p.MayHaveSideEffects(n) {
		// Transform a * (b / c) to b / c * a
		return p.TryRotateCommutativeOperator(n)
	}
	return n
}

// Java: PeepholeSubstituteAlternateSyntax.tryFoldSimpleFunctionCall(Node)
func (p *PeepholeSubstituteAlternateSyntax) TryFoldSimpleFunctionCall(n *rhino.Node) *rhino.Node {
	jutil.CheckStateMsg1(n.IsCall(), "%s", n)
	callTarget := n.GetFirstChild()
	if callTarget == nil || !callTarget.IsName() {
		return n
	}
	targetName := callTarget.GetString()
	switch targetName {
	case "Boolean":
		// Fold Boolean(a) to !!a
		// http://www.ecma-international.org/ecma-262/6.0/index.html#sec-boolean-constructor-boolean-value
		// and
		// http://www.ecma-international.org/ecma-262/6.0/index.html#sec-logical-not-operator-runtime-semantics-evaluation
		paramCount := n.GetChildCount() - 1
		// only handle the single known parameter case
		if paramCount == 1 {
			value := n.GetLastChild().Detach()
			var replacement *rhino.Node
			if NodeUtil_IsBooleanResult(value) {
				// If it is already a boolean do nothing.
				replacement = value
			} else {
				// Replace it with a "!!value"
				replacement = rhino.IR_Not(rhino.IR_Not(value).Srcref(n))
			}
			n.ReplaceWith(replacement)
			p.ReportChangeToEnclosingScope(replacement)
		}
	case "String":
		// Fold String(a) to '' + (a) on immutable literals,
		// which allows further optimizations
		//
		// We can't do this in the general case, because String(a) has
		// slightly different semantics than '' + (a). See
		// https://blickly.github.io/closure-compiler-issues/#759
		value := callTarget.GetNext()
		if value != nil && value.GetNext() == nil && NodeUtil_IsImmutableValue(value) {
			addition := rhino.IR_Add(rhino.IR_String("").Srcref(callTarget), value.Detach())
			n.ReplaceWith(addition)
			p.ReportChangeToEnclosingScope(addition)
			return addition
		}
	default:
		// nothing.
	}
	return n
}

// Java: PeepholeSubstituteAlternateSyntax.tryFoldImmediateCallToBoundFunction(Node)
func (p *PeepholeSubstituteAlternateSyntax) TryFoldImmediateCallToBoundFunction(n *rhino.Node) *rhino.Node {
	// Rewriting "(fn.bind(a,b))()" to "fn.call(a,b)" makes it inlinable
	jutil.CheckState(n.IsCall())
	callTarget := n.GetFirstChild()
	bind := p.GetCodingConvention().DescribeFunctionBind(callTarget, false)
	if bind != nil {
		// replace the call target
		bind.target.Detach()
		callTarget.ReplaceWith(bind.target)
		callTarget = bind.target

		// push the parameters
		PeepholeSubstituteAlternateSyntax_AddParameterAfter(bind.parameters, callTarget)

		// add the this value before the parameters if necessary
		if bind.thisValue != nil && !NodeUtil_IsUndefined(bind.thisValue) {
			// rewrite from "fn(a, b)" to "fn.call(thisValue, a, b)"
			newCallTarget := rhino.IR_Getprop(callTarget.CloneTree(), "call")
			p.MarkNewScopesChanged(newCallTarget)
			callTarget.ReplaceWith(newCallTarget)
			p.MarkFunctionsDeleted(callTarget)
			bind.thisValue.CloneTree().InsertAfter(newCallTarget)
			n.PutBooleanProp(rhino.Node_FREE_CALL, false)
		} else {
			n.PutBooleanProp(rhino.Node_FREE_CALL, true)
		}
		p.ReportChangeToEnclosingScope(n)
	}
	return n
}

// Java: PeepholeSubstituteAlternateSyntax.addParameterAfter(Node, Node)
func PeepholeSubstituteAlternateSyntax_AddParameterAfter(parameterList *rhino.Node, after *rhino.Node) {
	if parameterList != nil {
		// push the last parameter to the head of the list first.
		PeepholeSubstituteAlternateSyntax_AddParameterAfter(parameterList.GetNext(), after)
		parameterList.CloneTree().InsertAfter(after)
	}
}

// SplitComma converts expressions of a potentially nested comma expression
// into a sequence of expression result statements and inserts them into the
// AST.
//
// insert is whether or not the leftmost expression is inserted into the AST.
// Returns the leftmost expression.
//
// Java: PeepholeSubstituteAlternateSyntax.splitComma(Node, boolean, Node)
func (p *PeepholeSubstituteAlternateSyntax) SplitComma(n *rhino.Node, insert bool, insertAfter *rhino.Node) *rhino.Node {
	for n.IsComma() {
		left := n.GetFirstChild()
		right := n.GetLastChild()
		n.DetachChildren()
		if right.IsComma() {
			p.SplitComma(right, true, insertAfter)
		} else {
			// Add the right expression after the optimized expression.
			newStatement := rhino.IR_ExprResult(right)
			newStatement.SrcrefIfMissing(right)
			newStatement.InsertAfter(insertAfter)
		}
		n = left
	}
	if insert {
		newStatement := rhino.IR_ExprResult(n)
		newStatement.SrcrefIfMissing(n)
		newStatement.InsertAfter(insertAfter)
		return newStatement
	}
	return n
}

// Java: PeepholeSubstituteAlternateSyntax.trySplitComma(Node)
func (p *PeepholeSubstituteAlternateSyntax) TrySplitComma(n *rhino.Node) *rhino.Node {
	if p.late {
		return n
	}
	jutil.CheckState(n.IsExprResult())
	if n.GetParent().IsLabel() {
		// Do not split labeled comma expressions.
		return n
	}
	comma := n.GetFirstChild()
	jutil.CheckState(comma.IsComma())
	leftmost := p.SplitComma(comma, false, n)
	// Replace original expression with leftmost comma expression.
	n.RemoveChildren()
	n.AddChildToFront(leftmost)
	n.Srcref(leftmost)
	p.ReportChangeToEnclosingScope(leftmost)
	return leftmost
}

// TryReplaceUndefined uses "void 0" in place of "undefined".
//
// Java: PeepholeSubstituteAlternateSyntax.tryReplaceUndefined(Node)
func (p *PeepholeSubstituteAlternateSyntax) TryReplaceUndefined(n *rhino.Node) *rhino.Node {
	// TODO(johnlenz): consider doing this as a normalization.
	if p.IsASTNormalized() &&
		NodeUtil_IsUndefined(n) &&
		!NodeUtil_IsLValue(n) {
		replacement := NodeUtil_NewUndefinedNode(n)
		n.ReplaceWith(replacement)
		p.ReportChangeToEnclosingScope(replacement)
		return replacement
	}
	return n
}

// Java: PeepholeSubstituteAlternateSyntax.tryReduceCall(Node)
func (p *PeepholeSubstituteAlternateSyntax) TryReduceCall(node *rhino.Node) *rhino.Node {
	result := p.TryFoldLiteralConstructor(node)
	if result == node {
		result = p.TryFoldSimpleFunctionCall(node)
		if result == node {
			result = p.TryFoldImmediateCallToBoundFunction(node)
		}
	}
	return result
}

// TryReduceReturn reduces "return undefined" or "return void 0" to simply
// "return". Returns the original node, maybe simplified.
//
// Java: PeepholeSubstituteAlternateSyntax.tryReduceReturn(Node)
func (p *PeepholeSubstituteAlternateSyntax) TryReduceReturn(n *rhino.Node) *rhino.Node {
	result := n.GetFirstChild()

	if result != nil {
		switch result.GetToken() {
		case rhino.Token_VOID:
			operand := result.GetFirstChild()
			if !p.MayHaveSideEffects(operand) {
				n.RemoveFirstChild()
				p.ReportChangeToEnclosingScope(n)
			}
		case rhino.Token_NAME:
			name := result.GetString()
			if name == "undefined" {
				n.RemoveFirstChild()
				p.ReportChangeToEnclosingScope(n)
			}
		default:
		}
	}

	return n
}

// String, Number, and Boolean functions return non-object types, whereas
// new String, new Number, and new Boolean return object types, so don't
// include them here.
var peepholeSubstituteAlternateSyntax_STANDARD_OBJECT_CONSTRUCTORS = jutil.NewLinkedHashSet(
	"Object",
	"Array",
	"Error")

// TryFoldStandardConstructors folds "new Object()" to "Object()".
//
// Java: PeepholeSubstituteAlternateSyntax.tryFoldStandardConstructors(Node)
func (p *PeepholeSubstituteAlternateSyntax) TryFoldStandardConstructors(n *rhino.Node) *rhino.Node {
	jutil.CheckState(n.IsNew())

	if p.CanFoldStandardConstructors(n) {
		n.SetToken(rhino.Token_CALL)
		n.PutBooleanProp(rhino.Node_FREE_CALL, true)
		p.ReportChangeToEnclosingScope(n)
	}

	return n
}

// CanFoldStandardConstructors returns whether "new Object()" can be folded
// to "Object()" on n.
//
// Java: PeepholeSubstituteAlternateSyntax.canFoldStandardConstructors(Node)
func (p *PeepholeSubstituteAlternateSyntax) CanFoldStandardConstructors(n *rhino.Node) bool {
	// If name normalization has been run then we know that
	// new Object() does in fact refer to what we think it is
	// and not some custom-defined Object().
	if p.IsASTNormalized() && n.GetFirstChild().IsName() {
		className := n.GetFirstChild().GetString()
		if peepholeSubstituteAlternateSyntax_STANDARD_OBJECT_CONSTRUCTORS.Contains(className) {
			return true
		}
		if "RegExp" == className {
			// Fold "new RegExp()" to "RegExp()", but only if the argument is a string.
			// See issue 1260.
			if n.GetSecondChild() == nil || n.GetSecondChild().IsStringLit() {
				return true
			}
		}
	}

	return false
}

// TryFoldLiteralConstructor replaces a new Array, Object, or RegExp node with
// a literal, unless the call is to a local constructor function with the same
// name.
//
// Java: PeepholeSubstituteAlternateSyntax.tryFoldLiteralConstructor(Node)
func (p *PeepholeSubstituteAlternateSyntax) TryFoldLiteralConstructor(n *rhino.Node) *rhino.Node {
	jutil.CheckArgument(n.IsCall() || n.IsNew())

	constructorNameNode := n.GetFirstChild()

	var newLiteralNode *rhino.Node

	// We require the AST to be normalized to ensure that, say,
	// Object() really refers to the built-in Object constructor
	// and not a user-defined constructor with the same name.

	if p.IsASTNormalized() && constructorNameNode.IsName() {

		className := constructorNameNode.GetString()

		constructorHasArgs := constructorNameNode.GetNext() != nil

		if "Object" == className && !constructorHasArgs {
			// "Object()" --> "{}"
			newLiteralNode = rhino.IR_Objectlit()
		} else if "Array" == className {
			// "Array(arg0, arg1, ...)" --> "[arg0, arg1, ...]"
			arg0 := constructorNameNode.GetNext()
			action := PeepholeSubstituteAlternateSyntax_IsSafeToFoldArrayConstructor(arg0)

			if action == peepholeSubstituteAlternateSyntax_FoldArrayAction_SAFE_TO_FOLD_WITH_ARGS ||
				action == peepholeSubstituteAlternateSyntax_FoldArrayAction_SAFE_TO_FOLD_WITHOUT_ARGS {
				newLiteralNode = rhino.IR_Arraylit()
				n.RemoveFirstChild() // discard the function name
				elements := n.RemoveChildren()
				if action == peepholeSubstituteAlternateSyntax_FoldArrayAction_SAFE_TO_FOLD_WITH_ARGS {
					newLiteralNode.AddChildrenToFront(elements)
				}
			}
		}

		if newLiteralNode != nil {
			n.ReplaceWith(newLiteralNode)
			p.ReportChangeToEnclosingScope(newLiteralNode)
			return newLiteralNode
		}
	}
	return n
}

// Java: PeepholeSubstituteAlternateSyntax.FoldArrayAction
type peepholeSubstituteAlternateSyntax_FoldArrayAction int32

const (
	peepholeSubstituteAlternateSyntax_FoldArrayAction_NOT_SAFE_TO_FOLD peepholeSubstituteAlternateSyntax_FoldArrayAction = iota
	peepholeSubstituteAlternateSyntax_FoldArrayAction_SAFE_TO_FOLD_WITH_ARGS
	peepholeSubstituteAlternateSyntax_FoldArrayAction_SAFE_TO_FOLD_WITHOUT_ARGS
)

var peepholeSubstituteAlternateSyntax_FoldArrayAction_names = [...]string{
	"NOT_SAFE_TO_FOLD",
	"SAFE_TO_FOLD_WITH_ARGS",
	"SAFE_TO_FOLD_WITHOUT_ARGS",
}

// Java: PeepholeSubstituteAlternateSyntax.FoldArrayAction.values()
func peepholeSubstituteAlternateSyntax_FoldArrayAction_values() []peepholeSubstituteAlternateSyntax_FoldArrayAction {
	return []peepholeSubstituteAlternateSyntax_FoldArrayAction{
		peepholeSubstituteAlternateSyntax_FoldArrayAction_NOT_SAFE_TO_FOLD,
		peepholeSubstituteAlternateSyntax_FoldArrayAction_SAFE_TO_FOLD_WITH_ARGS,
		peepholeSubstituteAlternateSyntax_FoldArrayAction_SAFE_TO_FOLD_WITHOUT_ARGS,
	}
}

// Java: PeepholeSubstituteAlternateSyntax.FoldArrayAction.name()
func (a peepholeSubstituteAlternateSyntax_FoldArrayAction) String() string {
	return peepholeSubstituteAlternateSyntax_FoldArrayAction_names[a]
}

// PeepholeSubstituteAlternateSyntax_IsSafeToFoldArrayConstructor checks if
// it is safe to fold Array() constructor into []. It can be obviously done,
// if the initial constructor has either no arguments or at least two. The
// remaining case may be unsafe since Array(number) actually reserves memory
// for an empty array which contains number elements.
//
// Java: PeepholeSubstituteAlternateSyntax.isSafeToFoldArrayConstructor(Node)
func PeepholeSubstituteAlternateSyntax_IsSafeToFoldArrayConstructor(arg *rhino.Node) peepholeSubstituteAlternateSyntax_FoldArrayAction {
	action := peepholeSubstituteAlternateSyntax_FoldArrayAction_NOT_SAFE_TO_FOLD

	if arg == nil {
		action = peepholeSubstituteAlternateSyntax_FoldArrayAction_SAFE_TO_FOLD_WITHOUT_ARGS
	} else if arg.GetNext() != nil {
		action = peepholeSubstituteAlternateSyntax_FoldArrayAction_SAFE_TO_FOLD_WITH_ARGS
	} else {
		switch arg.GetToken() {
		case rhino.Token_STRINGLIT:
			// "Array('a')" --> "['a']"
			action = peepholeSubstituteAlternateSyntax_FoldArrayAction_SAFE_TO_FOLD_WITH_ARGS
		case rhino.Token_NUMBER:
			// "Array(0)" --> "[]"
			if arg.GetDouble() == 0 {
				action = peepholeSubstituteAlternateSyntax_FoldArrayAction_SAFE_TO_FOLD_WITHOUT_ARGS
			}
		case rhino.Token_ARRAYLIT:
			// "Array([args])" --> "[[args]]"
			action = peepholeSubstituteAlternateSyntax_FoldArrayAction_SAFE_TO_FOLD_WITH_ARGS
		default:
		}
	}
	return action
}

// Java: PeepholeSubstituteAlternateSyntax.reduceSubstractionAssignment(Node)
func (p *PeepholeSubstituteAlternateSyntax) ReduceSubstractionAssignment(n *rhino.Node) *rhino.Node {
	right := n.GetLastChild()
	isNegative := false
	if right.IsNeg() {
		isNegative = true
		right = right.GetOnlyChild()
	}

	if right.IsNumber() && right.GetDouble() == 1 {
		left := n.RemoveFirstChild()
		var newNode *rhino.Node
		if isNegative {
			newNode = rhino.IR_Inc(left, false)
		} else {
			newNode = rhino.IR_Dec(left, false)
		}
		n.ReplaceWith(newNode)
		p.ReportChangeToEnclosingScope(newNode)
		return newNode
	}

	return n
}

// Java: PeepholeSubstituteAlternateSyntax.reduceTrueFalse(Node)
func (p *PeepholeSubstituteAlternateSyntax) ReduceTrueFalse(n *rhino.Node) *rhino.Node {
	if p.late {
		switch n.GetParent().GetToken() {
		case rhino.Token_EQ, rhino.Token_GT, rhino.Token_GE, rhino.Token_LE, rhino.Token_LT, rhino.Token_NE:
			var value float64
			if n.IsTrue() {
				value = 1
			} else {
				value = 0
			}
			number := rhino.IR_Number(value)
			n.ReplaceWith(number)
			p.ReportChangeToEnclosingScope(number)
			return number
		default:
		}

		var value float64
		if n.IsTrue() {
			value = 0
		} else {
			value = 1
		}
		not := rhino.IR_Not(rhino.IR_Number(value))
		not.SrcrefTreeIfMissing(n)
		n.ReplaceWith(not)
		p.ReportChangeToEnclosingScope(not)
		return not
	}
	return n
}

// Java: PeepholeSubstituteAlternateSyntax.tryMinimizeArrayLiteral(Node)
func (p *PeepholeSubstituteAlternateSyntax) TryMinimizeArrayLiteral(n *rhino.Node) *rhino.Node {
	if !p.late {
		return n
	}
	for cur := n.GetFirstChild(); cur != nil; cur = cur.GetNext() {
		if !cur.IsStringLit() {
			return n
		}
	}
	return p.TryMinimizeStringArrayLiteral(n)
}

// Java: PeepholeSubstituteAlternateSyntax.tryMinimizeStringArrayLiteral(Node)
func (p *PeepholeSubstituteAlternateSyntax) TryMinimizeStringArrayLiteral(n *rhino.Node) *rhino.Node {
	if !p.late {
		return n
	}

	numElements := n.GetChildCount()
	// We save two bytes per element.
	saving := numElements*2 - peepholeSubstituteAlternateSyntax_STRING_SPLIT_OVERHEAD
	if saving <= 0 {
		return n
	}

	strs := make([]string, numElements)
	idx := 0
	for cur := n.GetFirstChild(); cur != nil; cur = cur.GetNext() {
		strs[idx] = cur.GetString()
		idx++
	}

	// These delimiters are chars that appears a lot in the program therefore
	// probably have a small Huffman encoding.
	delimiter := PeepholeSubstituteAlternateSyntax_PickDelimiter(strs)
	if delimiter != nil {
		template := strings.Join(strs, *delimiter)
		call := rhino.IR_Call(rhino.IR_Getprop(rhino.IR_String(template), "split"), rhino.IR_String(""+*delimiter))
		call.SrcrefTreeIfMissing(n)
		n.ReplaceWith(call)
		p.ReportChangeToEnclosingScope(call)
		return call
	}
	return n
}

// Java: PeepholeSubstituteAlternateSyntax.tryTurnTemplateStringsToStrings(Node)
func (p *PeepholeSubstituteAlternateSyntax) TryTurnTemplateStringsToStrings(n *rhino.Node) *rhino.Node {
	jutil.CheckStateMsg1(n.IsTemplateLit(), "%s", n)
	if n.GetParent().IsTaggedTemplateLit() {
		return n
	}
	str := p.GetSideEffectFreeStringValue(n)
	if str == nil {
		return n
	}
	stringNode := rhino.IR_String(*str).Srcref(n)
	n.ReplaceWith(stringNode)
	p.ReportChangeToEnclosingScope(stringNode)
	return stringNode
}

// PeepholeSubstituteAlternateSyntax_PickDelimiter finds a delimiter that
// does not occur in the given strings. strs are the strings that must be
// separated. Returns a delimiter string or null.
//
// Java: PeepholeSubstituteAlternateSyntax.pickDelimiter(String[])
func PeepholeSubstituteAlternateSyntax_PickDelimiter(strs []string) *string {
	allLength1 := true
	for _, s := range strs {
		if jutil.StringLength(s) != 1 {
			allLength1 = false
			break
		}
	}

	if allLength1 {
		empty := ""
		return &empty
	}

	space, semicolon, comma, lbrace, rbrace := " ", ";", ",", "{", "}"
	delimiters := []*string{&space, &semicolon, &comma, &lbrace, &rbrace, nil}
	i := 0
NEXT_DELIMITER:
	for ; delimiters[i] != nil; i++ {
		for _, cur := range strs {
			if strings.Contains(cur, *delimiters[i]) {
				continue NEXT_DELIMITER
			}
		}
		break
	}
	return delimiters[i]
}
