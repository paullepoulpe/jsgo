// Port of com/google/javascript/jscomp/PeepholeSubstituteAlternateSyntax.java @ closure-compiler v20261021

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

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Go-only: this pass uses the typed AST view (internal/ast) to read the
// shapes it rewrites (call.Callee(), call.OnlyArg(), ret.Value(),
// assign.Value(), getProp.Object()...) instead of child positions, and each
// tryX method takes the kind it rewrites (ast.Call, ast.Return...), so the
// Java checks of the node's token are gone. It no longer follows the Java
// line by line; the Go-only notes below say where the logic differs. Code
// that moves children around (splitComma, the bound function call, the
// array literal built from a call's arguments) and reduceTrueFalse, which
// only looks at its parent's token, use the rhino API.

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
func (p *PeepholeSubstituteAlternateSyntax) OptimizeSubtree(node rhino.Node) rhino.Node {
	switch node.GetToken() {
	case rhino.Token_ASSIGN_SUB:
		assign, _ := ast.AsAssignOp(node) // the case proves the kind
		return p.ReduceSubstractionAssignment(assign)
	case rhino.Token_TRUE, rhino.Token_FALSE:
		return p.ReduceTrueFalse(node)
	case rhino.Token_NEW:
		newNode, _ := ast.AsNew(node) // the case proves the kind
		result := p.TryFoldStandardConstructors(newNode)
		// tryFoldStandardConstructors() may convert a NEW node into a CALL node
		if call, ok := ast.AsCall(result); ok {
			return p.TryReduceCall(call)
		}
		return result
	case rhino.Token_CALL:
		call, _ := ast.AsCall(node) // the case proves the kind
		return p.TryReduceCall(call)
	case rhino.Token_RETURN:
		ret, _ := ast.AsReturn(node) // the case proves the kind
		return p.TryReduceReturn(ret)
	case rhino.Token_EXPR_RESULT:
		stmt, _ := ast.AsExprStmt(node) // the case proves the kind
		if stmt.Expr().Token() == rhino.Token_COMMA {
			return p.TrySplitComma(stmt)
		}
		return node
	case rhino.Token_NAME:
		name, _ := ast.AsName(node) // the case proves the kind
		if name.Ident() != "undefined" {
			return node
		}
		return p.TryReplaceUndefined(name)
	case rhino.Token_ARRAYLIT:
		array, _ := ast.AsArrayLit(node) // the case proves the kind
		return p.TryMinimizeArrayLiteral(array)
	case rhino.Token_GETPROP:
		getProp, _ := ast.AsGetProp(node) // the case proves the kind
		return p.TryMinimizeWindowRefs(getProp)
	case rhino.Token_TEMPLATELIT:
		template, _ := ast.AsTemplateLit(node) // the case proves the kind
		return p.TryTurnTemplateStringsToStrings(template)
	case rhino.Token_AND, rhino.Token_OR, rhino.Token_BITOR, rhino.Token_BITXOR, rhino.Token_BITAND, rhino.Token_COALESCE:
		op, _ := ast.AsBinary(node) // the case proves the kind
		return p.TryRotateAssociativeOperator(op)
	case rhino.Token_MUL:
		op, _ := ast.AsBinary(node) // the case proves the kind
		if !p.late || p.MayHaveSideEffects(node) {
			return node
		}
		return p.TryRotateCommutativeOperator(op)
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
func (p *PeepholeSubstituteAlternateSyntax) TryMinimizeWindowRefs(getProp ast.GetProp) rhino.Node {
	node := getProp.Node()
	window, ok := getProp.Object().AsName()
	if !ok || window.Ident() != "window" {
		return node
	}
	// Normalization needs to be done to ensure there's no shadowing.
	// Since normalization has run we know we're referring to the global window.
	if !p.IsASTNormalized() || !peepholeSubstituteAlternateSyntax_BUILTIN_EXTERNS.Contains(getProp.Property()) {
		return node
	}
	newNameNode := rhino.IR_Name(getProp.Property())
	parentNode := node.GetParent()

	newNameNode.Srcref(node)
	node.ReplaceWith(newNameNode)

	// Go-only note: like the Java, this marks the parent call as free whatever
	// the role of the window ref in it, so `obj.f(window.Array)` prints as
	// `(0,obj.f)(Array)` and calls f without obj as this. Kept for identical
	// output.
	if _, ok := ast.AsNormalOrOptChainCall(parentNode); ok {
		// e.g. when converting `window.Array?.()` to `Array?.()`, ensure that the
		// OPTCHAIN_CALL gets marked as `FREE_CALL`
		parentNode.PutBooleanProp(rhino.Node_FREE_CALL, true)
	}
	p.ReportChangeToEnclosingScope(parentNode)
	return newNameNode
}

// Java: PeepholeSubstituteAlternateSyntax.tryRotateCommutativeOperator(Node)
func (p *PeepholeSubstituteAlternateSyntax) TryRotateCommutativeOperator(n ast.Binary) rhino.Node {
	if !p.late {
		return n.Node()
	}
	// Transform a * (b / c) to b / c * a
	rhs := n.Right()
	lhs := n.Left()
	for lhs.Token() == n.Token() && NodeUtil_IsAssociative(n.Token()) {
		lhsOp, _ := lhs.AsBinary() // it has the token of n, a binary operator
		lhs = lhsOp.Left()
	}
	precedence := NodeUtil_Precedence(n.Token())
	lhsPrecedence := NodeUtil_Precedence(lhs.Token())
	rhsPrecedence := NodeUtil_Precedence(rhs.Token())
	if rhsPrecedence == precedence && lhsPrecedence != precedence {
		rhs.Detach()
		lhs.ReplaceWith(rhs)
		// (n has lost its right operand: SetRight would replace the left one.)
		n.Node().AddChildToBack(lhs.Node())
		p.ReportChangeToEnclosingScope(n.Node())
	}
	return n.Node()
}

// Java: PeepholeSubstituteAlternateSyntax.tryRotateAssociativeOperator(Node)
func (p *PeepholeSubstituteAlternateSyntax) TryRotateAssociativeOperator(n ast.Binary) rhino.Node {
	if !p.late {
		return n.Node()
	}
	// All commutative operators are also associative
	jutil.CheckArgument(NodeUtil_IsAssociative(n.Token()))
	if rhs, ok := n.Right().AsBinary(); ok && rhs.Token() == n.Token() {
		// Transform a * (b * c) to a * b * c
		first, second, third := n.Left(), rhs.Left(), rhs.Right()
		newLhs := ast.NewBinary(n.Token(), first.Detach(), second.Detach())
		newLhs.Node().SrcrefIfMissing(n.Node())
		newRoot := ast.NewBinary(rhs.Token(), newLhs.AsExpr(), third.Detach())
		newRoot.Node().SrcrefIfMissing(rhs.Node())
		n.Node().ReplaceWith(newRoot.Node())
		p.ReportChangeToEnclosingScope(newRoot.Node())
		return newRoot.Node()
	} else if NodeUtil_IsCommutative(n.Token()) && !p.MayHaveSideEffects(n.Node()) {
		// Transform a * (b / c) to b / c * a
		return p.TryRotateCommutativeOperator(n)
	}
	return n.Node()
}

// Java: PeepholeSubstituteAlternateSyntax.tryFoldSimpleFunctionCall(Node)
func (p *PeepholeSubstituteAlternateSyntax) TryFoldSimpleFunctionCall(call ast.Call) rhino.Node {
	n := call.Node()
	callTarget, ok := call.Callee().AsName()
	if !ok {
		return n
	}
	switch callTarget.Ident() {
	case "Boolean":
		// Fold Boolean(a) to !!a
		// http://www.ecma-international.org/ecma-262/6.0/index.html#sec-boolean-constructor-boolean-value
		// and
		// http://www.ecma-international.org/ecma-262/6.0/index.html#sec-logical-not-operator-runtime-semantics-evaluation
		// only handle the single known parameter case
		arg, ok := call.OnlyArg()
		if !ok {
			break
		}
		// Go-only: a spread argument is left alone. Java folds Boolean(...a) to
		// !!...a, which IR.not rejects: an internal compiler error.
		value, ok := arg.AsExpr()
		if !ok {
			break
		}
		value.Detach()
		var replacement rhino.Node
		if NodeUtil_IsBooleanResult(value.Node()) {
			// If it is already a boolean do nothing.
			replacement = value.Node()
		} else {
			// Replace it with a "!!value"
			replacement = rhino.IR_Not(rhino.IR_Not(value.Node()).Srcref(n))
		}
		n.ReplaceWith(replacement)
		p.ReportChangeToEnclosingScope(replacement)
		// (Like the Java, this returns n, now out of the tree, not the
		// replacement.)
	case "String":
		// Fold String(a) to '' + (a) on immutable literals,
		// which allows further optimizations
		//
		// We can't do this in the general case, because String(a) has
		// slightly different semantics than '' + (a). See
		// https://blickly.github.io/closure-compiler-issues/#759
		if value, ok := call.OnlyArg(); ok && NodeUtil_IsImmutableValue(value.Node()) {
			addition := rhino.IR_Add(rhino.IR_String("").Srcref(callTarget.Node()), value.Detach().Node())
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
func (p *PeepholeSubstituteAlternateSyntax) TryFoldImmediateCallToBoundFunction(call ast.Call) rhino.Node {
	// Rewriting "(fn.bind(a,b))()" to "fn.call(a,b)" makes it inlinable
	n := call.Node()
	boundCall, ok := call.Callee().AsCall()
	if !ok {
		return n
	}
	bind := p.GetCodingConvention().DescribeFunctionBind(boundCall.Node(), false)
	if bind != nil {
		// replace the call target
		bind.target.Detach()
		boundCall.Node().ReplaceWith(bind.target)
		callTarget := bind.target

		// push the parameters
		PeepholeSubstituteAlternateSyntax_AddParameterAfter(bind.parameters, callTarget)

		// add the this value before the parameters if necessary
		if !bind.thisValue.IsNil() && !NodeUtil_IsUndefined(bind.thisValue) {
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
func PeepholeSubstituteAlternateSyntax_AddParameterAfter(parameterList rhino.Node, after rhino.Node) {
	if !parameterList.IsNil() {
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
func (p *PeepholeSubstituteAlternateSyntax) SplitComma(n rhino.Node, insert bool, insertAfter rhino.Node) rhino.Node {
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
func (p *PeepholeSubstituteAlternateSyntax) TrySplitComma(stmt ast.ExprStmt) rhino.Node {
	n := stmt.Node()
	if p.late {
		return n
	}
	if n.GetParent().IsLabel() {
		// Do not split labeled comma expressions.
		return n
	}
	comma := stmt.Expr()
	jutil.CheckState(comma.Token() == rhino.Token_COMMA)
	leftmost := p.SplitComma(comma.Node(), false, n)
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
func (p *PeepholeSubstituteAlternateSyntax) TryReplaceUndefined(name ast.Name) rhino.Node {
	n := name.Node()
	// TODO(johnlenz): consider doing this as a normalization.
	if name.Ident() == "undefined" &&
		p.IsASTNormalized() &&
		!NodeUtil_IsLValue(n) {
		replacement := NodeUtil_NewUndefinedNode(n)
		n.ReplaceWith(replacement)
		p.ReportChangeToEnclosingScope(replacement)
		return replacement
	}
	return n
}

// Java: PeepholeSubstituteAlternateSyntax.tryReduceCall(Node)
func (p *PeepholeSubstituteAlternateSyntax) TryReduceCall(call ast.Call) rhino.Node {
	switch call.Callee().Token() {
	case rhino.Token_NAME:
		result := p.TryFoldLiteralConstructor(call)
		if call.Is(result) {
			result = p.TryFoldSimpleFunctionCall(call)
		}
		return result
	case rhino.Token_CALL:
		return p.TryFoldImmediateCallToBoundFunction(call)
	default:
		return call.Node()
	}
}

// TryReduceReturn reduces "return undefined" or "return void 0" to simply
// "return". Returns the original node, maybe simplified.
//
// Java: PeepholeSubstituteAlternateSyntax.tryReduceReturn(Node)
func (p *PeepholeSubstituteAlternateSyntax) TryReduceReturn(ret ast.Return) rhino.Node {
	n := ret.Node()
	result, ok := ret.Value()
	if !ok {
		return n
	}
	if void, ok := result.AsUnary(); ok && void.Token() == rhino.Token_VOID {
		if !p.MayHaveSideEffects(void.Operand().Node()) {
			ret.RemoveValue()
			p.ReportChangeToEnclosingScope(n)
		}
	} else if name, ok := result.AsName(); ok && name.Ident() == "undefined" {
		ret.RemoveValue()
		p.ReportChangeToEnclosingScope(n)
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
func (p *PeepholeSubstituteAlternateSyntax) TryFoldStandardConstructors(newNode ast.New) rhino.Node {
	n := newNode.Node()
	if p.CanFoldStandardConstructors(newNode) {
		// The node becomes a call in place: same children, new token.
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
func (p *PeepholeSubstituteAlternateSyntax) CanFoldStandardConstructors(newNode ast.New) bool {
	// If name normalization has been run then we know that
	// new Object() does in fact refer to what we think it is
	// and not some custom-defined Object().
	constructor, ok := newNode.Callee().AsName()
	if !ok || !p.IsASTNormalized() {
		return false
	}
	className := constructor.Ident()
	if peepholeSubstituteAlternateSyntax_STANDARD_OBJECT_CONSTRUCTORS.Contains(className) {
		return true
	}
	if "RegExp" == className {
		// Fold "new RegExp()" to "RegExp()", but only if the argument is a string.
		// See issue 1260.
		arg, ok := newNode.FirstArg()
		if !ok {
			return true
		}
		_, isString := arg.AsStringLit()
		return isString
	}
	return false
}

// TryFoldLiteralConstructor replaces a new Array, Object, or RegExp node with
// a literal, unless the call is to a local constructor function with the same
// name.
//
// (Java also accepts a new; it is only called on calls.)
//
// Java: PeepholeSubstituteAlternateSyntax.tryFoldLiteralConstructor(Node)
func (p *PeepholeSubstituteAlternateSyntax) TryFoldLiteralConstructor(call ast.Call) rhino.Node {
	n := call.Node()

	// We require the AST to be normalized to ensure that, say,
	// Object() really refers to the built-in Object constructor
	// and not a user-defined constructor with the same name.
	constructorName, ok := call.Callee().AsName()
	if !ok {
		return n
	}
	className := constructorName.Ident()
	if ("Object" != className && "Array" != className) || !p.IsASTNormalized() {
		return n
	}

	var newLiteralNode rhino.Node
	if "Object" == className && !call.HasArgs() {
		// "Object()" --> "{}"
		newLiteralNode = rhino.IR_Objectlit()
	} else if "Array" == className {
		// "Array(arg0, arg1, ...)" --> "[arg0, arg1, ...]"
		action := PeepholeSubstituteAlternateSyntax_IsSafeToFoldArrayConstructor(call)

		if action == peepholeSubstituteAlternateSyntax_FoldArrayAction_SAFE_TO_FOLD_WITH_ARGS ||
			action == peepholeSubstituteAlternateSyntax_FoldArrayAction_SAFE_TO_FOLD_WITHOUT_ARGS {
			newLiteralNode = rhino.IR_Arraylit()
			constructorName.Detach() // discard the function name
			elements := n.RemoveChildren()
			if action == peepholeSubstituteAlternateSyntax_FoldArrayAction_SAFE_TO_FOLD_WITH_ARGS {
				newLiteralNode.AddChildrenToFront(elements)
			}
		}
	}

	if newLiteralNode.IsNil() {
		return n
	}
	n.ReplaceWith(newLiteralNode)
	p.ReportChangeToEnclosingScope(newLiteralNode)
	return newLiteralNode
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
// it is safe to fold the Array() constructor call into []. It can be
// obviously done, if the initial constructor has either no arguments or at
// least two. The remaining case may be unsafe since Array(number) actually
// reserves memory for an empty array which contains number elements.
//
// (Go: takes the call; Java takes its first argument, or null.)
//
// Java: PeepholeSubstituteAlternateSyntax.isSafeToFoldArrayConstructor(Node)
func PeepholeSubstituteAlternateSyntax_IsSafeToFoldArrayConstructor(call ast.Call) peepholeSubstituteAlternateSyntax_FoldArrayAction {
	arg, ok := call.FirstArg()
	if !ok {
		return peepholeSubstituteAlternateSyntax_FoldArrayAction_SAFE_TO_FOLD_WITHOUT_ARGS
	}
	if _, ok := arg.Next(); ok {
		return peepholeSubstituteAlternateSyntax_FoldArrayAction_SAFE_TO_FOLD_WITH_ARGS
	}
	if _, ok := arg.AsStringLit(); ok {
		// "Array('a')" --> "['a']"
		return peepholeSubstituteAlternateSyntax_FoldArrayAction_SAFE_TO_FOLD_WITH_ARGS
	}
	if number, ok := arg.AsNumber(); ok && number.Value() == 0 {
		// "Array(0)" --> "[]"
		return peepholeSubstituteAlternateSyntax_FoldArrayAction_SAFE_TO_FOLD_WITHOUT_ARGS
	}
	if _, ok := arg.AsArrayLit(); ok {
		// "Array([args])" --> "[[args]]"
		return peepholeSubstituteAlternateSyntax_FoldArrayAction_SAFE_TO_FOLD_WITH_ARGS
	}
	return peepholeSubstituteAlternateSyntax_FoldArrayAction_NOT_SAFE_TO_FOLD
}

// Java: PeepholeSubstituteAlternateSyntax.reduceSubstractionAssignment(Node)
func (p *PeepholeSubstituteAlternateSyntax) ReduceSubstractionAssignment(assign ast.AssignOp) rhino.Node {
	n := assign.Node()
	right := assign.Value()
	isNegative := false
	if neg, ok := right.AsUnary(); ok && neg.Token() == rhino.Token_NEG {
		isNegative = true
		right = neg.Operand()
	}

	if number, ok := right.AsNumber(); ok && number.Value() == 1 {
		left := assign.Target().Detach().Node()
		var newNode rhino.Node
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

// (n is a TRUE or a FALSE, two kinds without a common category in package
// ast, and the rewrite depends on its parent's token: it stays on the rhino
// API.)
//
// Java: PeepholeSubstituteAlternateSyntax.reduceTrueFalse(Node)
func (p *PeepholeSubstituteAlternateSyntax) ReduceTrueFalse(n rhino.Node) rhino.Node {
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
func (p *PeepholeSubstituteAlternateSyntax) TryMinimizeArrayLiteral(array ast.ArrayLit) rhino.Node {
	if !p.late {
		return array.Node()
	}
	for elem := range array.Elements() {
		if _, ok := elem.AsStringLit(); !ok {
			return array.Node()
		}
	}
	return p.TryMinimizeStringArrayLiteral(array)
}

// Java: PeepholeSubstituteAlternateSyntax.tryMinimizeStringArrayLiteral(Node)
func (p *PeepholeSubstituteAlternateSyntax) TryMinimizeStringArrayLiteral(array ast.ArrayLit) rhino.Node {
	n := array.Node()
	if !p.late {
		return n
	}

	numElements := n.GetChildCount()
	// We save two bytes per element.
	saving := numElements*2 - peepholeSubstituteAlternateSyntax_STRING_SPLIT_OVERHEAD
	if saving <= 0 {
		return n
	}

	strs := make([]string, 0, numElements)
	for elem := range array.Elements() {
		str, ok := elem.AsStringLit()
		jutil.CheckState(ok) // tryMinimizeArrayLiteral checked every element
		strs = append(strs, str.Value())
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
func (p *PeepholeSubstituteAlternateSyntax) TryTurnTemplateStringsToStrings(template ast.TemplateLit) rhino.Node {
	n := template.Node()
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
