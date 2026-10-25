// Port of com/google/javascript/jscomp/PeepholeRemoveDeadCode.java @ closure-compiler v20261021

/*
 * Copyright 2004 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Go-only: this pass uses the typed AST view (internal/ast) to read
// statements (ifs.Else(), try.Catch(), sw.Body().Cases(), do.Cond()...)
// instead of child positions, so it no longer follows the Java line by
// line; the Go-only notes below say where the logic differs. Code that
// treats children alike (the statements of a SCRIPT or a BLOCK, the
// operands of any expression) or moves nodes around uses the rhino API.

// PeepholeRemoveDeadCode is a peephole optimization to remove useless code
// such as IF's with false guard conditions, comma operator's left hand sides
// with no side effects, etc.
type PeepholeRemoveDeadCode struct {
	AbstractPeepholeOptimizationBase
	dequePool  []*peepholeRemoveDeadCode_nodeDeque
	dequeDepth int
}

// Java: PeepholeRemoveDeadCode()
func NewPeepholeRemoveDeadCode() *PeepholeRemoveDeadCode {
	p := &PeepholeRemoveDeadCode{}
	p.InitAbstractPeepholeOptimizationBase(p) // the implicit super()
	return p
}

// Java: PeepholeRemoveDeadCode.optimizeSubtree(Node)
func (p *PeepholeRemoveDeadCode) OptimizeSubtree(subtree rhino.Node) rhino.Node {
	switch subtree.GetToken() {
	case rhino.Token_ASSIGN:
		return p.TryFoldAssignment(subtree)
	case rhino.Token_COMMA:
		return p.TryFoldComma(subtree)
	case rhino.Token_SCRIPT, rhino.Token_BLOCK:
		return p.TryOptimizeBlock(subtree)
	case rhino.Token_EXPR_RESULT:
		return p.TryFoldExpr(subtree)
	case rhino.Token_HOOK:
		return p.TryFoldHook(subtree)
	case rhino.Token_SWITCH:
		return p.TryOptimizeSwitch(subtree)
	case rhino.Token_IF:
		return p.TryFoldIf(subtree)
	case rhino.Token_WHILE:
		// This pass gets run both before and after denormalize. Hence, the AST could potentially
		// contain WHILE (denormalized).
		// TODO: Ideally, we should optimize this case instead of returning
		return subtree
	case rhino.Token_FOR:
		forNode := ast.MustFor(subtree)
		if cond, ok := forNode.Cond(); ok {
			p.TryFoldForCondition(cond.Node())
		}
		return p.TryFoldFor(subtree)
	case rhino.Token_DO:
		foldedDo := p.TryFoldDoAway(subtree)
		if foldedDo.IsDo() {
			return p.TryFoldEmptyDo(foldedDo)
		}
		return foldedDo
	case rhino.Token_TRY:
		return p.TryFoldTry(subtree)
	case rhino.Token_LABEL:
		return p.TryFoldLabel(subtree)
	case rhino.Token_ARRAY_PATTERN:
		return p.TryOptimizeArrayPattern(subtree)
	case rhino.Token_OBJECT_PATTERN:
		return p.TryOptimizeObjectPattern(subtree)
	case rhino.Token_VAR, rhino.Token_CONST, rhino.Token_LET:
		return p.TryOptimizeNameDeclaration(subtree)
	case rhino.Token_DEFAULT_VALUE:
		return p.TryRemoveDefaultValue(subtree)
	case rhino.Token_OPTCHAIN_GETPROP, rhino.Token_OPTCHAIN_CALL, rhino.Token_OPTCHAIN_GETELEM:
		return p.TryRemoveOptionalChaining(subtree)
	default:
		return subtree
	}
}

// Java: PeepholeRemoveDeadCode.tryRemoveDefaultValue(Node)
func (p *PeepholeRemoveDeadCode) TryRemoveDefaultValue(n rhino.Node) rhino.Node {
	defaultValue, ok := ast.AsDefaultValue(n)
	jutil.CheckArgumentMsg1(ok, "%s", n)

	lValue := defaultValue.Target().Node()
	val := defaultValue.Value()
	removeVal := false

	// If the default is `undefined` always remove the value
	if name, ok := val.AsName(); ok && name.Ident() == "undefined" {
		removeVal = true
	}

	// If the `void` application is pure, remove the value
	if void, ok := val.AsUnary(); ok && void.Token() == rhino.Token_VOID {
		removeVal = !p.MayHaveSideEffects(void.Operand().Node())
	}

	if removeVal {
		n.ReplaceWith(lValue.Detach())
		p.ReportChangeToEnclosingScope(lValue)
		return lValue
	}

	return n
}

// Java: PeepholeRemoveDeadCode.tryFoldLabel(Node)
func (p *PeepholeRemoveDeadCode) TryFoldLabel(n rhino.Node) rhino.Node {
	label, ok := ast.AsLabel(n)
	jutil.CheckState(ok)
	labelName := label.Name().Ident()
	stmt := label.Body()
	if _, ok := stmt.AsEmpty(); ok {
		p.ReportChangeToEnclosingScope(n)
		n.Detach()
		return rhino.NilNode
	}

	if block, ok := stmt.AsBlock(); ok && !block.HasStatements() {
		p.ReportChangeToEnclosingScope(n)
		if n.GetParent().IsLabel() {
			// If the parent is itself a label, replace this label
			// with its contained block to keep the AST in a valid state.
			n.ReplaceWith(block.Node().Detach())
		} else {
			n.Detach()
		}
		return rhino.NilNode
	}

	if child, ok := PeepholeRemoveDeadCode_GetOnlyInterestingChild(stmt); ok {
		stmt = child
	}
	// Go-only: Java reads the label of the break without checking that it has
	// one, and throws on an unlabeled break such as `a: break;` in a switch
	// case. Label() makes the missing label explicit: such a label is kept.
	brk, ok := stmt.AsBreak()
	if !ok {
		return n
	}
	if target, ok := brk.Label(); ok && target.Ident() == labelName {
		if n.GetParent().IsLabel() {
			replacement := rhino.IR_Block().Srcref(n)
			n.ReplaceWith(replacement)
			p.ReportChangeToEnclosingScope(replacement)
			return replacement
		} else {
			parent := n.GetParent()
			n.Detach()
			p.ReportChangeToEnclosingScope(parent)
			return rhino.NilNode
		}
	}
	return n
}

// PeepholeRemoveDeadCode_GetOnlyInterestingChild returns the only
// "interesting" child of block, if it has exactly one interesting child;
// ok is false otherwise (or if stmt is not a block). For purposes of this
// method, a node is considered "interesting" unless it is an empty
// synthetic block.
//
// Java: PeepholeRemoveDeadCode.getOnlyInterestingChild(Node)
func PeepholeRemoveDeadCode_GetOnlyInterestingChild(stmt ast.Statement) (ast.Statement, bool) {
	block, ok := stmt.AsBlock()
	if !ok {
		return ast.Statement{}, false
	}
	if block.Node().HasOneChild() {
		return block.FirstStatement()
	}

	var ret ast.Statement
	found := false
	for child := range block.Statements() {
		if child.Node().IsSyntheticBlock() && !child.Node().HasChildren() {
			// Uninteresting child.
		} else if found {
			// Found more than one interesting child.
			return ast.Statement{}, false
		} else {
			ret, found = child, true
		}
	}
	return ret, found
}

// TryFoldTry removes try blocks without catch blocks and with empty or not
// existent finally blocks. Or, only leave the finally blocks if try body
// blocks are empty.
//
// Returns the replacement node, if changed, or the original if not.
//
// Java: PeepholeRemoveDeadCode.tryFoldTry(Node)
func (p *PeepholeRemoveDeadCode) TryFoldTry(n rhino.Node) rhino.Node {
	try, ok := ast.AsTry(n)
	jutil.CheckStateMsg1(ok, "%s", n)
	body := try.Body()
	catch, hasCatch := try.Catch()
	finally, hasFinally := try.Finally()

	// Removes TRYs that had its CATCH removed and/or empty FINALLY.
	if !hasCatch && (!hasFinally || !finally.HasStatements()) {
		jutil.CheckState(!n.GetParent().IsLabel())
		body.Node().Detach()
		n.ReplaceWith(body.Node())
		p.ReportChangeToEnclosingScope(body.Node())
		return body.Node()
	}

	// Only leave FINALLYs if TRYs are not empty
	if !body.HasStatements() {
		if hasCatch {
			NodeUtil_RedeclareVarsInsideBranch(catch.Node())
		}
		p.ReportChangeToEnclosingScope(n)
		if hasFinally {
			finally.Node().Detach()
			jutil.CheckState(!n.GetParent().IsLabel())
			n.ReplaceWith(finally.Node())
			return finally.Node()
		} else {
			jutil.CheckState(!n.GetParent().IsLabel())
			n.Detach()
			return rhino.NilNode
		}
	}

	return n
}

// TryFoldAssignment tries removing identity assignments and empty
// destructuring pattern assignments.
//
// Returns the replacement node, if changed, or the original if not.
//
// Java: PeepholeRemoveDeadCode.tryFoldAssignment(Node)
func (p *PeepholeRemoveDeadCode) TryFoldAssignment(subtree rhino.Node) rhino.Node {
	assign, ok := ast.AsAssign(subtree)
	jutil.CheckState(ok)
	left := assign.Target()
	right := assign.Value().Node()
	leftName, leftIsName := left.AsName()
	rightName, rightIsName := assign.Value().AsName()
	if leftIsName && rightIsName && leftName.Ident() == rightName.Ident() {
		// Only names
		subtree.ReplaceWith(right.Detach())
		p.ReportChangeToEnclosingScope(right)
		return right
	} else if pattern, ok := left.AsPattern(); ok && !pattern.Node().HasChildren() {
		// `[] = <expr>` becomes `<expr>`
		// Note that this does potentially change behavior. If `<expr>` is not iterable and this
		// code originally threw, it will no longer throw.
		subtree.ReplaceWith(right.Detach())
		p.ReportChangeToEnclosingScope(right)
		return right
	}
	return subtree
}

// TryOptimizeNameDeclaration tries removing identity assignments and empty
// destructuring pattern assignments.
//
// Returns the replacement node, if changed, or the original if not.
//
// Java: PeepholeRemoveDeadCode.tryOptimizeNameDeclaration(Node)
func (p *PeepholeRemoveDeadCode) TryOptimizeNameDeclaration(subtree rhino.Node) rhino.Node {
	decl, ok := ast.AsNameDeclaration(subtree)
	jutil.CheckState(ok)
	first, ok := decl.FirstDeclarator()
	if !ok {
		return subtree
	}
	lhs, ok := first.AsDestructuringLhs()
	if !ok {
		return subtree
	}
	if init, ok := lhs.Init(); ok && !lhs.Pattern().Node().HasChildren() {
		// `var [] = foo();` becomes `foo();`
		value := init.Node()
		subtree.ReplaceWith(rhino.IR_ExprResult(value.Detach()).Srcref(value))
		p.ReportChangeToEnclosingScope(value)
	}
	return subtree
}

// TryFoldExpr tries folding EXPR_RESULT nodes by removing useless Ops and
// expressions.
//
// Returns the replacement node, if changed, or the original if not.
//
// Java: PeepholeRemoveDeadCode.tryFoldExpr(Node)
func (p *PeepholeRemoveDeadCode) TryFoldExpr(subtree rhino.Node) rhino.Node {
	stmt, ok := ast.AsExprStmt(subtree)
	jutil.CheckState(ok)
	expr := stmt.Expr().Node()
	switch expr.GetToken() {
	case rhino.Token_ASSIGN,
		rhino.Token_ASSIGN_BITOR,
		rhino.Token_ASSIGN_BITXOR,
		rhino.Token_ASSIGN_BITAND,
		rhino.Token_ASSIGN_LSH,
		rhino.Token_ASSIGN_RSH,
		rhino.Token_ASSIGN_URSH,
		rhino.Token_ASSIGN_ADD,
		rhino.Token_ASSIGN_SUB,
		rhino.Token_ASSIGN_MUL,
		rhino.Token_ASSIGN_DIV,
		rhino.Token_ASSIGN_MOD,
		rhino.Token_ASSIGN_EXPONENT,
		rhino.Token_ASSIGN_OR,
		rhino.Token_ASSIGN_AND,
		rhino.Token_ASSIGN_COALESCE,
		rhino.Token_INC,
		rhino.Token_DEC,
		rhino.Token_DELPROP,
		rhino.Token_YIELD,
		rhino.Token_AWAIT,
		rhino.Token_DYNAMIC_IMPORT:
		return subtree
	case rhino.Token_CALL, rhino.Token_NEW, rhino.Token_TAGGED_TEMPLATELIT:
		if p.NodeTypeMayHaveSideEffects(expr) {
			return subtree
		}
	}
	result := p.TrySimplifyUnusedResult(expr)
	if result.IsNil() {
		parent := subtree.GetParent()
		// If the EXPR_RESULT no longer has any children, remove it as well.
		if parent.IsLabel() {
			replacement := rhino.IR_Block().Srcref(subtree)
			subtree.ReplaceWith(replacement)
			subtree = replacement
		} else {
			subtree.Detach()
			subtree = rhino.NilNode
		}
	}
	return subtree
}

// TrySimplifyUnusedResult replaces expression with an expression that
// contains only side-effects of the original.
//
// This replacement is made under the assumption that the result of
// expression is unused and therefore it is correct to eliminate
// non-side-effectful nodes.
//
// Returns the replacement expression, or null if there were no side-effects
// to preserve.
func (p *PeepholeRemoveDeadCode) acquireDeque() *peepholeRemoveDeadCode_nodeDeque {
	idx := p.dequeDepth
	p.dequeDepth++
	if idx < len(p.dequePool) {
		d := p.dequePool[idx]
		clear(d.elems)
		d.elems = d.elems[:0]
		d.head = 0
		return d
	}
	d := &peepholeRemoveDeadCode_nodeDeque{
		elems: make([]rhino.Node, 0, 8),
	}
	p.dequePool = append(p.dequePool, d)
	return d
}

// Java: PeepholeRemoveDeadCode.trySimplifyUnusedResult(Node)
func (p *PeepholeRemoveDeadCode) TrySimplifyUnusedResult(expression rhino.Node) rhino.Node {
	sideEffectRoots := p.acquireDeque()
	atFixedPoint := p.TrySimplifyUnusedResultInternal(expression, sideEffectRoots)
	p.dequeDepth--

	if atFixedPoint {
		// `expression` is in a form that cannot be further optimized.
		return expression
	} else if sideEffectRoots.isEmpty() {
		p.DeleteNode(expression)
		return rhino.NilNode
	} else if sideEffectRoots.peekFirst() == expression {
		// Expression was a conditional that was transformed. There can't be any other side-effects,
		// but we also can't detach the transformed root.
		jutil.CheckStateMsg1(sideEffectRoots.size() == 1, "%s", sideEffectRoots)
		p.ReportChangeToEnclosingScope(expression)
		return expression
	} else {
		sideEffects := PeepholeRemoveDeadCode_AsDetachedExpression(sideEffectRoots.pollFirst())

		// Assemble a tree of comma expressions for all the side-effects. The tree must execute the
		// side-effects in FIFO order with respect to the queue. It must also be left leaning to match
		// the parser's preferred structure.
		for !sideEffectRoots.isEmpty() {
			next := PeepholeRemoveDeadCode_AsDetachedExpression(sideEffectRoots.pollFirst())
			sideEffects = rhino.IR_Comma(sideEffects, next).Srcref(next)
		}

		sideEffects.InsertBefore(expression)
		p.DeleteNode(expression)
		return sideEffects
	}
}

// peepholeRemoveDeadCode_nodeDeque is the ArrayDeque<Node> used by
// trySimplifyUnusedResult (only addLast, peekFirst, pollFirst, size and
// isEmpty are used).
//
// Go-only: language-adapter — java.util.ArrayDeque.
type peepholeRemoveDeadCode_nodeDeque struct {
	elems []rhino.Node
	head  int
}

func (d *peepholeRemoveDeadCode_nodeDeque) addLast(n rhino.Node) {
	rhino.CheckNotNull(n)
	d.elems = append(d.elems, n)
}

func (d *peepholeRemoveDeadCode_nodeDeque) peekFirst() rhino.Node {
	if d.head == len(d.elems) {
		return rhino.NilNode
	}
	return d.elems[d.head]
}

func (d *peepholeRemoveDeadCode_nodeDeque) pollFirst() rhino.Node {
	if d.head == len(d.elems) {
		return rhino.NilNode
	}
	n := d.elems[d.head]
	d.head++
	return n
}

func (d *peepholeRemoveDeadCode_nodeDeque) size() int { return len(d.elems) - d.head }

func (d *peepholeRemoveDeadCode_nodeDeque) isEmpty() bool { return d.head == len(d.elems) }

// String is AbstractCollection.toString().
func (d *peepholeRemoveDeadCode_nodeDeque) String() string {
	return jutil.CollectionToString(d.elems[d.head:])
}

// TrySimplifyUnusedResultInternal collects any potentially side-effectful
// subtrees within tree into sideEffectRoots.
//
// When a node is determined to have side-effects its descendants are not
// explored. This method assumes the entire subtree of such a node must be
// preserved. As a corollary, the contents of sideEffectRoots are a forest.
//
// This operation generally does not mutate tree; however, exceptions are
// made for expressions that alter control-flow. Such expression will be
// pruned of their side-effectless branches. Even in this case, tree is never
// detached.
//
// sideEffectRoots are the roots of subtrees determined to have side-effects,
// in execution order. Returns true iff there is no code to be removed from
// within tree; it is already at a fixed point for code removal.
//
// Java: PeepholeRemoveDeadCode.trySimplifyUnusedResultInternal(Node, ArrayDeque<Node>)
func (p *PeepholeRemoveDeadCode) TrySimplifyUnusedResultInternal(tree rhino.Node, sideEffectRoots *peepholeRemoveDeadCode_nodeDeque) bool {
	// Special cases for conditional expressions that may be using results.
	switch tree.GetToken() {
	case rhino.Token_OPTCHAIN_GETELEM, rhino.Token_OPTCHAIN_GETPROP, rhino.Token_OPTCHAIN_CALL:
		if p.MayHaveSideEffects(tree) {
			sideEffectRoots.addLast(tree)
			return PeepholeRemoveDeadCode_HasFixedPointParent(tree)
		} else {
			return false
		}
	case rhino.Token_HOOK:
		// Try to remove one or more of the conditional children and transform the HOOK to an
		// equivalent operation. Remember that if either value branch still exists, the result of
		// the predicate expression is being used, and so cannot be removed.
		//    x() ? foo() : 1 --> x() && foo()
		//    x() ? 1 : foo() --> x() || foo()
		//    x() ? 1 : 1 --> x()
		//    x ? 1 : 1 --> null
		hook := ast.MustHook(tree)
		// Read the operands first: simplifying a branch can remove it from the hook.
		cond, thenBranch, elseBranch := hook.Cond(), hook.Then(), hook.Else()
		trueNode := p.TrySimplifyUnusedResult(thenBranch.Node())
		falseNode := p.TrySimplifyUnusedResult(elseBranch.Node())
		if trueNode.IsNil() && !falseNode.IsNil() {
			jutil.CheckStateMsg1(tree.HasTwoChildren(), "%s", tree)

			tree.SetToken(rhino.Token_OR)
			sideEffectRoots.addLast(tree)
			return false // The node type was changed.
		} else if !trueNode.IsNil() && falseNode.IsNil() {
			jutil.CheckStateMsg1(tree.HasTwoChildren(), "%s", tree)

			tree.SetToken(rhino.Token_AND)
			sideEffectRoots.addLast(tree)
			return false // The node type was changed.
		} else if trueNode.IsNil() && falseNode.IsNil() {
			// Don't bother adding true and false branch children to make the AST valid; this HOOK is
			// going to be deleted. We just need to collect any side-effects from the predicate
			// expression.
			p.TrySimplifyUnusedResultInternal(cond.Node(), sideEffectRoots)
			return false // This HOOK must be cleaned up.
		} else {
			sideEffectRoots.addLast(tree)
			return PeepholeRemoveDeadCode_HasFixedPointParent(tree)
		}
	case rhino.Token_AND, rhino.Token_OR, rhino.Token_COALESCE:
		// Try to remove the second operand from a AND, OR, and COALESCE operations. Remember that
		// if the second
		// child still exists, the result of the first expression is being used, and so cannot be
		// removed.
		//    x() ?? f --> x()
		//    x() || f --> x()
		//    x() && f --> x()
		op := ast.MustBinary(tree)
		// Read the operands first: simplifying the right one can remove it.
		left, right := op.Left(), op.Right()
		conditionalResultNode := p.TrySimplifyUnusedResult(right.Node())
		if conditionalResultNode.IsNil() {
			// Don't bother adding a second child to make the AST valid; this op is going to be
			// deleted. We just need to collect any side-effects from the predicate first child.
			p.TrySimplifyUnusedResultInternal(left.Node(), sideEffectRoots)
			return false // This op must be cleaned up.
		} else {
			sideEffectRoots.addLast(tree)
			return PeepholeRemoveDeadCode_HasFixedPointParent(tree)
		}
	case rhino.Token_FUNCTION:
		// Functions that aren't being invoked are dead. If they were invoked we'd see the CALL
		// before arriving here. We don't want to look at any children since they'll never execute.
	case rhino.Token_MEMBER_FIELD_DEF:
		field := ast.MustMemberFieldDef(tree)
		if value, ok := field.Value(); ok && tree.IsStaticMember() {
			p.TrySimplifyUnusedResultInternal(value.Node(), sideEffectRoots)
		}
	case rhino.Token_COMPUTED_FIELD_DEF:
		field := ast.MustComputedFieldDef(tree)
		value, hasValue := field.Value()
		// Computed key initializers run once when the class is defined, even on non-static fields.
		p.TrySimplifyUnusedResultInternal(field.Key().Node(), sideEffectRoots)
		if hasValue && tree.IsStaticMember() {
			p.TrySimplifyUnusedResultInternal(value.Node(), sideEffectRoots)
		}
	default:
		// This is the meat of this function. It covers the general case of nodes which are unused
		if p.NodeTypeMayHaveSideEffects(tree) {
			sideEffectRoots.addLast(tree)
			return PeepholeRemoveDeadCode_HasFixedPointParent(tree)
		} else if !tree.HasChildren() {
			return false // A node must have children or side-effects to be at fixed-point.
		}

		atFixedPoint := PeepholeRemoveDeadCode_HasFixedPointParent(tree)
		for child := tree.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
			atFixedPoint = p.TrySimplifyUnusedResultInternal(child, sideEffectRoots) && atFixedPoint
		}
		return atFixedPoint
	}
	return false
}

// PeepholeRemoveDeadCode_AsDetachedExpression returns an expression
// executing expr which is legal in any expression context.
//
// expr is an attached expression. Returns a detached expression.
//
// Java: PeepholeRemoveDeadCode.asDetachedExpression(Node)
func PeepholeRemoveDeadCode_AsDetachedExpression(expr rhino.Node) rhino.Node {
	switch expr.GetToken() {
	case rhino.Token_ITER_SPREAD, rhino.Token_OBJECT_SPREAD:
		switch expr.GetParent().GetToken() {
		case rhino.Token_ARRAYLIT, rhino.Token_NEW,
			rhino.Token_CALL,          // `Math.sin(...c)`
			rhino.Token_OPTCHAIN_CALL: // `Math?.sin(...c)`
			expr = rhino.IR_Arraylit(expr.Detach()).Srcref(expr)
		case rhino.Token_OBJECTLIT:
			expr = rhino.IR_Objectlit(expr.Detach()).Srcref(expr)
		default:
			panic(jutil.NewIllegalStateExceptionMsg(expr.ToStringTree()))
		}
	default:
	}

	if expr.HasParent() {
		expr.Detach()
	}

	jutil.CheckStateMsg1(rhino.IR_MayBeExpression(expr), "%s", expr)
	return expr
}

// PeepholeRemoveDeadCode_HasFixedPointParent returns true iff expr is
// parented such that it is valid in a fixed-point representation of an
// unused expression tree.
//
// A fixed-point representation is one in which no further nodes should be
// changed or removed when removing unused code. This method assumes that the
// expression tree in question is unused, so only side-effects are relevant.
//
// Java: PeepholeRemoveDeadCode.hasFixedPointParent(Node)
func PeepholeRemoveDeadCode_HasFixedPointParent(expr rhino.Node) bool {
	// Most kinds of nodes shouldn't be branches in the fixed-point tree of an unused
	// expression. Those listed below are the only valid kinds.
	switch expr.GetParent().GetToken() {
	case rhino.Token_AND, rhino.Token_COMMA, rhino.Token_HOOK, rhino.Token_OR, rhino.Token_COALESCE:
		return true
	case rhino.Token_ARRAYLIT, rhino.Token_OBJECTLIT:
		// Make a special allowance for SPREADs so they remain in a legal context. Parent types
		// other than ARRAYLIT and OBJECTLIT are not fixed-point because they are the tersest
		// legal
		// parents and are known to be side-effect free.
		return expr.IsSpread()
	default:
		// Statments are always fixed-point parents. All other expressions are not.
		return NodeUtil_IsStatement(expr.GetParent())
	}
}

// peepholeRemoveDeadCode_isUnnamedBreak reports whether n is a break without
// a label.
//
// Go-only: Java writes n.isBreak() && !n.hasChildren() at each use.
func peepholeRemoveDeadCode_isUnnamedBreak(n rhino.Node) bool {
	brk, ok := ast.AsBreak(n)
	if !ok {
		return false
	}
	_, hasLabel := brk.Label()
	return !hasLabel
}

// peepholeRemoveDeadCode_isUnnamedContinue reports whether n is a continue
// without a label.
//
// Go-only: Java writes n.isContinue() && !n.hasChildren().
func peepholeRemoveDeadCode_isUnnamedContinue(n rhino.Node) bool {
	cont, ok := ast.AsContinue(n)
	if !ok {
		return false
	}
	_, hasLabel := cont.Label()
	return !hasLabel
}

// peepholeRemoveDeadCode_MatchUnnamedBreak is a predicate for matching
// anything except function nodes.
type peepholeRemoveDeadCode_MatchUnnamedBreak struct {
	// Go-only: language-adapter — distinct identity
	_ byte
}

// Java: PeepholeRemoveDeadCode.MatchUnnamedBreak.apply(Node)
func (m *peepholeRemoveDeadCode_MatchUnnamedBreak) Apply(n rhino.Node) bool {
	return peepholeRemoveDeadCode_isUnnamedBreak(n)
}

// Java: PeepholeRemoveDeadCode.MATCH_UNNAMED_BREAK
var PeepholeRemoveDeadCode_MATCH_UNNAMED_BREAK func(rhino.Node) bool = (&peepholeRemoveDeadCode_MatchUnnamedBreak{}).Apply

// RemoveIfUnnamedBreak removes maybeBreak if it is a break without a label.
//
// Go-only: Java also accepts null; callers here only pass a statement.
//
// Java: PeepholeRemoveDeadCode.removeIfUnnamedBreak(Node)
func (p *PeepholeRemoveDeadCode) RemoveIfUnnamedBreak(maybeBreak ast.Statement) {
	if peepholeRemoveDeadCode_isUnnamedBreak(maybeBreak.Node()) {
		p.ReportChangeToEnclosingScope(maybeBreak.Node())
		maybeBreak.Node().Detach()
	}
}

// Java: PeepholeRemoveDeadCode.tryRemoveSwitchWithSingleCase(Node, boolean)
func (p *PeepholeRemoveDeadCode) TryRemoveSwitchWithSingleCase(sw ast.Switch, shouldHoistCondition bool) rhino.Node {
	n := sw.Node()
	cases := sw.Body()
	jutil.CheckState(cases.Node().HasOneChild())
	onlyCase, _ := cases.FirstCase() // the check above proves there is one
	caseBlock := onlyCase.Body()
	if last, ok := caseBlock.LastStatement(); ok {
		p.RemoveIfUnnamedBreak(last)
	}
	// Back off if the switch contains statements like "if (a) { break; }"
	if NodeUtil_Has(caseBlock.Node(), PeepholeRemoveDeadCode_MATCH_UNNAMED_BREAK, NodeUtil_MATCH_NOT_FUNCTION) {
		return n
	}
	if shouldHoistCondition {
		cond := sw.Discriminant().Node()
		rhino.IR_ExprResult(cond.Detach()).Srcref(n).InsertBefore(n)
	}
	n.ReplaceWith(caseBlock.Node().Detach())
	p.ReportChangeToEnclosingScope(caseBlock.Node())
	return caseBlock.Node()
}

// Java: PeepholeRemoveDeadCode.tryRemoveSwitch(Node)
func (p *PeepholeRemoveDeadCode) TryRemoveSwitch(sw ast.Switch) rhino.Node {
	n := sw.Node()
	cases := sw.Body()
	if !cases.HasCases() {
		// Remove the switch if there are no remaining cases
		condition := sw.Discriminant().Node()
		replacement := rhino.IR_ExprResult(condition.Detach()).Srcref(n)
		n.ReplaceWith(replacement)
		p.ReportChangeToEnclosingScope(replacement)
		return replacement
	}
	if cases.Node().HasOneChild() {
		onlyCase, _ := cases.FirstCase() // the check above proves there is one
		if _, ok := onlyCase.AsDefaultCase(); ok {
			// Before removing switch, we must preserve the switch condition if it has side effects
			shouldHoistCondition := p.MayHaveSideEffects(sw.Discriminant().Node())
			return p.TryRemoveSwitchWithSingleCase(sw, shouldHoistCondition)
		}
	}
	return n
}

// peepholeRemoveDeadCode_caseTagOrBody returns the tag of a case, or the
// body of a default case.
//
// Go-only: Java reads the first child of each case, which for a default case
// (it has no tag) is its body: the switch optimizations below then compare
// that body with the condition and check it for side effects. This helper
// keeps that behavior, for the same output, and makes it visible.
func peepholeRemoveDeadCode_caseTagOrBody(c ast.SwitchCase) rhino.Node {
	if cs, ok := c.AsCase(); ok {
		return cs.Value().Node()
	}
	return c.Body().Node()
}

// TryOptimizeSwitch removes useless switches and cases.
//
// Java: PeepholeRemoveDeadCode.tryOptimizeSwitch(Node)
func (p *PeepholeRemoveDeadCode) TryOptimizeSwitch(n rhino.Node) rhino.Node {
	sw, ok := ast.AsSwitch(n)
	jutil.CheckStateMsg1(ok, "%s", n)

	cases := sw.Body()
	defaultCase := p.TryOptimizeDefaultCase(sw)

	// Generally, it is unsafe to remove other cases when the default case is not the last one.
	// (defaultCase, when there is one, is the only default case.)
	if (defaultCase.IsNil() || defaultCase.GetNext().IsNil()) && p.AreAllCaseTagsLiterals(cases) {
		cond := sw.Discriminant().Node()

		// First, remove empty cases where possible: always empty default cases; or when there is
		// no default case, other empty cases that are not the first matching case, may be
		// removable as well.
		var prev rhino.Node // the last case kept, if any
		foundMatchingCase := false
		for cur, ok := cases.FirstCase(); ok; {
			next, hasNext := cur.Next() // read it first: cur may be removed
			tag := peepholeRemoveDeadCode_caseTagOrBody(cur)
			foundMatchingCase = p.IsFirstSwitchMatch(foundMatchingCase, cond, tag)
			if !foundMatchingCase &&
				!p.MayHaveSideEffects(tag) &&
				p.IsUselessCase(cur, prev, defaultCase) {
				p.RemoveCase(n, cur.Node())
			} else {
				prev = cur.Node()
			}
			cur, ok = next, hasNext
		}

		// Next, optimize switches with constant condition
		if NodeUtil_IsLiteralValue(cond, false) {
			caseMatches := base.Tri_TRUE
			// Remove cases until you find one that may match
			cur, ok := cases.FirstCase()
			for ok {
				caseMatches = PeepholeFoldConstants_EvaluateComparison(
					p, rhino.Token_SHEQ, cond, peepholeRemoveDeadCode_caseTagOrBody(cur))
				if caseMatches == base.Tri_TRUE || caseMatches == base.Tri_UNKNOWN {
					break
				}
				next, hasNext := cur.Next()
				p.RemoveCase(n, cur.Node())
				cur, ok = next, hasNext
			}
			if ok && caseMatches == base.Tri_TRUE {
				// Skip cases until you find one whose last stm is a removable break
				matchingCase := cur
				matchingCaseBlock := matchingCase.Body()
				for ok {
					block := cur.Body()
					isLastStmRemovableBreak := false
					if lastStm, ok := block.LastStatement(); ok && PeepholeRemoveDeadCode_IsExit(lastStm.Node()) {
						p.RemoveIfUnnamedBreak(lastStm)
						isLastStmRemovableBreak = true
					}
					next, hasNext := cur.Next()
					// Remove the fallthrough case labels
					if cur != matchingCase {
						for block.HasStatements() {
							matchingCaseBlock.Node().AddChildToBack(block.Node().RemoveFirstChild())
						}
						p.ReportChangeToEnclosingScope(cur.Node())
						cur.Node().Detach()
					}
					cur, ok = next, hasNext
					if isLastStmRemovableBreak {
						break
					}
				}

				// Remove any remaining cases
				for ok {
					next, hasNext := cur.Next()
					p.RemoveCase(n, cur.Node())
					cur, ok = next, hasNext
				}
				// If there is one case left, we may be able to fold it
				// Go-only: Java tests cond.getNext(), which is the switch body (not
				// a case), so its test always passes; only the matching case is
				// left at this point anyway.
				if cases.Node().HasOneChild() {
					return p.TryRemoveSwitchWithSingleCase(sw, false)
				}
			}
		}
	}

	// Last, try to remove the entire switch if possible
	return p.TryRemoveSwitch(sw)
}

// TryOptimizeDefaultCase returns the default case node or null if there is
// no default case or if the default case is removed.
//
// Java: PeepholeRemoveDeadCode.tryOptimizeDefaultCase(Node)
func (p *PeepholeRemoveDeadCode) TryOptimizeDefaultCase(sw ast.Switch) rhino.Node {
	n := sw.Node()
	cases := sw.Body()
	var lastNonRemovable rhino.Node // the most recently iterated case known to not be removable.

	for c := range cases.Cases() {
		if cs, ok := c.AsCase(); ok {
			if cs.Body().HasStatements() || p.MayHaveSideEffects(cs.Value().Node()) {
				lastNonRemovable = c.Node()
			}
			continue
		}

		// c is the default case.
		// Remove any cases that fall-through to the default case
		caseToRemove := cases.Node().GetFirstChild()
		if !lastNonRemovable.IsNil() {
			caseToRemove = lastNonRemovable.GetNext()
		}
		for next := rhino.NilNode; !c.Is(caseToRemove); caseToRemove = next {
			next = caseToRemove.GetNext()
			p.RemoveCase(n, caseToRemove)
		}

		// Remove the default case if we can
		if p.IsUselessCase(c, lastNonRemovable, c.Node()) {
			p.RemoveCase(n, c.Node())
			return rhino.NilNode
		}
		return c.Node()
	}
	return rhino.NilNode
}

// RemoveCase removes the case from the switch redeclaring any variables
// declared in it.
//
// caseNode is the case to remove.
//
// Java: PeepholeRemoveDeadCode.removeCase(Node, Node)
func (p *PeepholeRemoveDeadCode) RemoveCase(switchNode rhino.Node, caseNode rhino.Node) {
	NodeUtil_RedeclareVarsInsideBranch(caseNode)
	caseNode.Detach()
	p.ReportChangeToEnclosingScope(switchNode)
}

// IsUselessCase reports whether the CASE or DEFAULT_CASE block does anything
// useful.
//
// The function assumes that when checking a CASE node there is no
// DEFAULT_CASE node in the SWITCH, or the DEFAULT_CASE is the last case in
// the SWITCH.
//
// previousCase (the case before caseNode) and defaultCase are the nil node
// when there is none.
//
// Java: PeepholeRemoveDeadCode.isUselessCase(Node, Node, Node)
func (p *PeepholeRemoveDeadCode) IsUselessCase(
	caseNode ast.SwitchCase, previousCase rhino.Node, defaultCase rhino.Node) bool {
	jutil.CheckState(previousCase.IsNil() || previousCase.GetNext() == caseNode.Node())
	// A case isn't useless if a previous case falls through to it unless it happens to be the last
	// case in the switch.
	if prev, ok := ast.AsSwitchCase(previousCase); ok {
		if _, hasNext := caseNode.Next(); hasNext {
			previousLast, ok := prev.Body().LastStatement()
			if !ok || !PeepholeRemoveDeadCode_IsExit(previousLast.Node()) {
				return false
			}
		}
	}

	for executingCase, ok := caseNode, true; ok; executingCase, ok = executingCase.Next() {
		// We only expect a DEFAULT case if the case we are checking is the
		// DEFAULT case.  Otherwise, we assume the DEFAULT case has already
		// been removed.
		c, isCase := executingCase.AsCase()
		jutil.CheckState(caseNode == executingCase || isCase)
		if isCase && p.MayHaveSideEffects(c.Value().Node()) {
			// The case falls thru to a case whose condition has a potential side-effect,
			// removing the candidate case would skip that side-effect, so don't.
			return false
		}
		for blockChild := range executingCase.Body().Statements() {
			// If this is a block with a labelless break, it is useless.
			switch blockChild.Node().GetToken() {
			case rhino.Token_BREAK:
				// A case with a single labelless break is useless if it is the default case or if
				// there is no default case. A break to a different control structure isn't useless.
				return peepholeRemoveDeadCode_isUnnamedBreak(blockChild.Node()) &&
					(defaultCase.IsNil() || executingCase.Is(defaultCase))
			case rhino.Token_VAR:
				// Variable declarations without initializations are OK: one declarator,
				// and no child under it (a destructuring declarator has its pattern).
				if blockChild.Node().HasOneChild() && blockChild.Node().GetFirstFirstChild().IsNil() {
					continue
				}
				return false
			default:
				return false
			}
		}
		// Look at the fallthrough case
	}
	return true
}

// Java: PeepholeRemoveDeadCode.isFirstSwitchMatch(boolean, Node, Node)
func (p *PeepholeRemoveDeadCode) IsFirstSwitchMatch(foundMatchingCase bool, condition rhino.Node, tag rhino.Node) bool {
	if foundMatchingCase {
		return false
	}
	return PeepholeFoldConstants_EvaluateComparison(p, rhino.Token_SHEQ, condition, tag) == base.Tri_TRUE
}

// Java: PeepholeRemoveDeadCode.areAllCaseTagsLiterals(Node)
func (p *PeepholeRemoveDeadCode) AreAllCaseTagsLiterals(cases ast.SwitchBody) bool {
	for c := range cases.Cases() {
		// A default case has no tag.
		if cs, ok := c.AsCase(); ok && !NodeUtil_IsLiteralValue(cs.Value().Node(), false) {
			return false
		}
	}
	return true
}

// PeepholeRemoveDeadCode_IsExit returns whether the node is a control flow
// exit from the current block.
//
// Java: PeepholeRemoveDeadCode.isExit(Node)
func PeepholeRemoveDeadCode_IsExit(n rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_BREAK, rhino.Token_CONTINUE, rhino.Token_RETURN, rhino.Token_THROW:
		return true
	case rhino.Token_SWITCH:
		return PeepholeRemoveDeadCode_IsSwitchExit(n)
	case rhino.Token_TRY:
		return PeepholeRemoveDeadCode_IsTryExit(n)
	default:
		return false
	}
}

// PeepholeRemoveDeadCode_IsUnconditionalBlockExit returns whether the block
// is a control flow exit from the block containing current switch or
// try..catch statement.
//
// Java: PeepholeRemoveDeadCode.isUnconditionalBlockExit(Node)
func PeepholeRemoveDeadCode_IsUnconditionalBlockExit(block ast.Block) bool {
	jutil.CheckStateMsg1(!block.Node().GetParent().IsLabel(), "%s", block.Node())

	// Last statement must lead out of the block.
	lastStm, ok := block.LastStatement()
	if !ok {
		return false
	}
	switch lastStm.Node().GetToken() {
	case rhino.Token_BREAK:
		if peepholeRemoveDeadCode_isUnnamedBreak(lastStm.Node()) {
			return false
		}
		// Last statement is OK - continue with checking others.
	case rhino.Token_RETURN, rhino.Token_THROW:
		// Last statement is OK - continue with checking others.
	default:
		return false
	}

	// Other statements can be anything except for unlabeled "break". But for simplicity, don't go
	// into inner blocks and complex constructs - instead, allow only the simplest statements.
	for child := range block.Statements() {
		if child == lastStm {
			break
		}
		switch child.Node().GetToken() {
		case rhino.Token_BREAK:
			if peepholeRemoveDeadCode_isUnnamedBreak(child.Node()) {
				return false
			}
			// This break is OK - continue with checking others.
		case rhino.Token_RETURN, rhino.Token_THROW, rhino.Token_FUNCTION, rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST, rhino.Token_EXPR_RESULT:
			// This statement is OK - continue with checking others.
		default:
			return false
		}
	}

	return true
}

// PeepholeRemoveDeadCode_IsSwitchExit returns true if the switch always
// "exits" (return, throw, etc).
//
// Java: PeepholeRemoveDeadCode.isSwitchExit(Node)
func PeepholeRemoveDeadCode_IsSwitchExit(n rhino.Node) bool {
	sw, ok := ast.AsSwitch(n)
	jutil.CheckStateMsg1(ok, "%s", n)

	hasDefaultCase := false

	for switchCase := range sw.Body().Cases() {
		if _, ok := switchCase.AsDefaultCase(); ok {
			hasDefaultCase = true
		}
		block := switchCase.Body()
		if !block.HasStatements() {
			if _, hasNext := switchCase.Next(); hasNext {
				continue // An empty case falls through to the next one.
			}
		}
		if !PeepholeRemoveDeadCode_IsUnconditionalBlockExit(block) {
			return false
		}
	}

	return hasDefaultCase
}

// PeepholeRemoveDeadCode_IsTryExit returns true if the try..catch always
// "exits" (return, throw, etc).
//
// Java: PeepholeRemoveDeadCode.isTryExit(Node)
func PeepholeRemoveDeadCode_IsTryExit(n rhino.Node) bool {
	try, ok := ast.AsTry(n)
	jutil.CheckStateMsg1(ok, "%s", n)

	// finally - regardless of the behavior of the other blocks,
	// an exit from the finally with guarantee that behavior.
	if finally, ok := try.Finally(); ok && PeepholeRemoveDeadCode_IsUnconditionalBlockExit(finally) {
		return true
	}
	// try
	if !PeepholeRemoveDeadCode_IsUnconditionalBlockExit(try.Body()) {
		return false
	}
	// catch
	catch, ok := try.Catch()
	return !ok || PeepholeRemoveDeadCode_IsUnconditionalBlockExit(catch.Body())
}

// Java: PeepholeRemoveDeadCode.tryFoldComma(Node)
func (p *PeepholeRemoveDeadCode) TryFoldComma(n rhino.Node) rhino.Node {
	comma, ok := ast.AsBinary(n)
	jutil.CheckState(ok)
	// If the left side does nothing replace the comma with the result.
	parent := n.GetParent()
	right := comma.Right().Node()

	left := p.TrySimplifyUnusedResult(comma.Left().Node())
	if left.IsNil() || !p.MayHaveSideEffects(left) {
		// Fold it!
		right.Detach()
		n.ReplaceWith(right)
		p.ReportChangeToEnclosingScope(parent)
		return right
	}
	return n
}

// TryOptimizeBlock tries removing unneeded block nodes and their useless
// children.
//
// Go-only: n is a SCRIPT or a BLOCK, whose children this walks alike, so it
// uses the rhino API.
//
// Java: PeepholeRemoveDeadCode.tryOptimizeBlock(Node)
func (p *PeepholeRemoveDeadCode) TryOptimizeBlock(n rhino.Node) rhino.Node {
	// Remove any useless children
	for c := n.GetFirstChild(); !c.IsNil(); {
		next := c.GetNext() // save c.next, since 'c' may be removed
		switch c.GetToken() {
		case rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST:
			p.TryOptimizeConditionalAfterAssign(c)
			c = next
			continue
		case rhino.Token_BREAK, rhino.Token_CONTINUE, rhino.Token_RETURN, rhino.Token_THROW:
			p.RemoveFollowingNodes(c)
			c = rhino.NilNode
			continue
		case rhino.Token_FUNCTION,
			rhino.Token_FOR,
			rhino.Token_FOR_IN,
			rhino.Token_FOR_OF,
			rhino.Token_FOR_AWAIT_OF,
			rhino.Token_WHILE,
			rhino.Token_DO,
			rhino.Token_LABEL,
			rhino.Token_IMPORT,
			rhino.Token_EXPORT,
			rhino.Token_DEBUGGER:
			c = next
			continue
		case rhino.Token_EXPR_RESULT:
			if expr := c.GetFirstChild(); expr.IsAssign() && expr.GetFirstChild().IsName() {
				p.TryOptimizeConditionalAfterAssign(c)
				c = next
				continue
			}
		}
		if !PeepholeRemoveDeadCode_IsUnremovableNode(c) && !p.MayHaveSideEffects(c) {
			// TODO(johnlenz): determine what this is actually removing. Candidates
			//    include: EMPTY nodes, control structures without children
			//    (removing infinite loops), empty try blocks.  What else?
			c.Detach()
			p.ReportChangeToEnclosingScope(n)
			p.MarkFunctionsDeleted(c)
		} else if PeepholeRemoveDeadCode_IsExit(c) {
			p.RemoveFollowingNodes(c)
			break
		} else {
			p.TryOptimizeConditionalAfterAssign(c)
		}
		c = next
	}

	if n.IsSyntheticBlock() || n.IsScript() || n.GetParent().IsNil() {
		return n
	}

	// Try to merge the block with its parent, or remove it if it is an empty class static block.
	parent := n.GetParent()
	if NodeUtil_TryMergeBlock(n, p.IsASTNormalized()) {
		p.ReportChangeToEnclosingScope(parent)
		return rhino.NilNode
	} else if parent.IsClassMembers() && !n.HasChildren() {
		n.Detach()
		p.ReportChangeToEnclosingScope(parent)
		return rhino.NilNode
	}

	return n
}

// Java: PeepholeRemoveDeadCode.removeFollowingNodes(Node)
func (p *PeepholeRemoveDeadCode) RemoveFollowingNodes(n rhino.Node) {
	parent := n.GetParent()
	if parent.GetLastChild() != n {
		changed := false
		for !n.GetNext().IsNil() {
			dead := n.GetNext()
			if NodeUtil_IsFunctionDeclaration(dead) {
				// Don't remove function declarations as they are hoisted
				n = dead
				continue
			}
			changed = true
			NodeUtil_RedeclareVarsInsideBranch(dead)
			p.RedeclareIfBlockScopedVar(dead, parent)
			dead.Detach()
			p.MarkFunctionsDeleted(dead)
		}
		if changed {
			p.ReportChangeToEnclosingScope(parent)
		}
	}
}

// RedeclareIfBlockScopedVar redeclares the given node's names in the parent
// scope if they are block-scoped declarations; otherwise does nothing.
//
// Java: PeepholeRemoveDeadCode.redeclareIfBlockScopedVar(Node, Node)
func (p *PeepholeRemoveDeadCode) RedeclareIfBlockScopedVar(decl rhino.Node, parent rhino.Node) {
	// This isn't called on function declarations ever: this method is only used in
	// removeFollowingNodes, which preserves function declarations as they're hoisted. This
	// preconditions check is to make sure nothing else tries reusing this code, since in that case
	// maybe they do want to handle function declarations.
	jutil.CheckStateMsg1(!NodeUtil_IsFunctionDeclaration(decl), "Unexpected function declaration %s", decl)
	if !NodeUtil_IsBlockScopedDeclaration(decl) {
		return
	}
	var blockScopedVars []rhino.Node
	if class, ok := ast.AsClass(decl); ok {
		if name, ok := class.Name(); ok {
			blockScopedVars = append(blockScopedVars, name.Node())
		}
	} else {
		NodeUtil_VisitLhsNodesInNode(
			decl,
			func(n rhino.Node) {
				if n.IsName() {
					blockScopedVars = append(blockScopedVars, n)
				}
			})
	}

	addBefore := NodeUtil_GetInsertionPointAfterAllInnerFunctionDeclarations(parent)
	for _, nameNode := range blockScopedVars {
		var_ := rhino.IR_Let(rhino.IR_Name(nameNode.GetString()).Srcref(nameNode)).Srcref(nameNode)
		NodeUtil_CopyNameAnnotations(nameNode, var_.GetFirstChild())
		if !addBefore.IsNil() {
			var_.InsertBefore(addBefore)
		} else {
			parent.AddChildToBack(var_)
		}
	}
	p.AddFeatureToEnclosingScript(parser.FeatureSet_Feature_LET_DECLARATIONS)
}

// PeepholeRemoveDeadCode_IsUnremovableNode reports whether n is one of the
// nodes that are unremovable but don't have side effects so they aren't
// caught by mayHaveSideEffects.
//
// Java: PeepholeRemoveDeadCode.isUnremovableNode(Node)
func PeepholeRemoveDeadCode_IsUnremovableNode(n rhino.Node) bool {
	return (n.IsBlock() && n.IsSyntheticBlock()) || n.IsScript()
}

// TryOptimizeConditionalAfterAssign attempts to replace the condition of if
// or hook immediately that is a reference to a name that is assigned
// immediately before.
//
// TODO(johnlenz): Consider moving this to a separate peephole pass.
//
// Java: PeepholeRemoveDeadCode.tryOptimizeConditionalAfterAssign(Node)
func (p *PeepholeRemoveDeadCode) TryOptimizeConditionalAfterAssign(n rhino.Node) {
	next := n.GetNext()
	if next.IsNil() {
		return
	}

	// Look for patterns like the following and replace the if-condition with
	// a constant value so it can later be folded:
	//   var a = /a/;
	//   if (a) {foo(a)}
	// or
	//   a = 0;
	//   a ? foo(a) : c;
	// or
	//   a = 0;
	//   a || foo(a);
	// or
	//   a = 0;
	//   a && foo(a)
	// or
	//   a = 0;
	//   a ?? foo(a)
	// TODO(johnlenz): This would be better handled by control-flow sensitive
	// constant propagation. As the other case that I want to handle is:
	//   i=0; for(;i<0;i++){}
	// as right now nothing facilitates removing a loop like that.
	// This is here simply to remove the cruft left behind goog.userAgent and
	// similar cases.

	lhsAssign, rhsAssign, ok := peepholeRemoveDeadCode_simpleAssignment(n)
	if !ok {
		return
	}
	conditionalRoot := p.GetConditionalRoot(next)
	if conditionalRoot.IsNil() {
		return
	}

	condition := p.GetConditionalStatementCondition(next)
	if !lhsAssign.Node().MatchesNameNode(condition) {
		return
	}
	switch conditionalRoot.GetToken() {
	case rhino.Token_AND, rhino.Token_OR, rhino.Token_IF, rhino.Token_HOOK:
		// conditionals that coerce their condition to a boolean
		value := NodeUtil_GetBooleanValue(rhsAssign.Node())
		if value != base.Tri_UNKNOWN {
			replacementConditionNode := NodeUtil_BooleanNode(value.ToBoolean(true))
			condition.ReplaceWith(replacementConditionNode)
			p.ReportChangeToEnclosingScope(replacementConditionNode)
		}
		return
	case rhino.Token_COALESCE:
		// conditional that checks whether its operand is nullish
		valueType := NodeUtil_GetKnownValueType(rhsAssign.Node())
		switch valueType {
		case NodeUtil_ValueType_NULL, NodeUtil_ValueType_VOID:
			condition.ReplaceWith(NodeUtil_NewUndefinedNode(condition))
			p.ReportChangeToEnclosingScope(conditionalRoot)
		case NodeUtil_ValueType_NUMBER, NodeUtil_ValueType_BIGINT, NodeUtil_ValueType_STRING, NodeUtil_ValueType_BOOLEAN, NodeUtil_ValueType_OBJECT:
			// semi-arbitrarily use '0' as a short non-nullish conditional
			condition.ReplaceWith(rhino.IR_Number(0).Srcref(condition))
			p.ReportChangeToEnclosingScope(conditionalRoot)
		case NodeUtil_ValueType_UNDETERMINED:
		}
		return
	default:
		panic(jutil.NewAssertionError("Unhandled condition " + conditionalRoot.String()))
	}
}

// peepholeRemoveDeadCode_simpleAssignment returns the name and the value of
// the simple assignment n; ok is false if n is not one.
//
// For our purposes we define a simple assignment to be a assignment to a
// NAME node (`a = value;`), or a VAR declaration with one child and a
// initializer (`var a = value;`, also with let or const).
//
// Go-only: Java has a predicate and a getter for the name and for the value;
// here one typed helper returns both. Java also accepts a declaration with
// one destructuring declarator (`var [a] = value;`), whose "name" never
// matches a condition: the result is the same.
//
// Java: PeepholeRemoveDeadCode.isSimpleAssignment(Node),
// getSimpleAssignmentName(Node) and getSimpleAssignmentValue(Node)
func peepholeRemoveDeadCode_simpleAssignment(n rhino.Node) (name ast.Name, value ast.Expr, ok bool) {
	if stmt, ok := ast.AsExprStmt(n); ok {
		assign, ok := stmt.Expr().AsAssign()
		if !ok {
			return ast.Name{}, ast.Expr{}, false
		}
		name, ok := assign.Target().AsName()
		return name, assign.Value(), ok
	}
	decl, ok := ast.AsNameDeclaration(n)
	if !ok {
		return ast.Name{}, ast.Expr{}, false
	}
	only, ok := decl.FirstDeclarator()
	if !ok {
		return ast.Name{}, ast.Expr{}, false
	}
	if _, more := only.Next(); more {
		return ast.Name{}, ast.Expr{}, false
	}
	name, ok = only.AsName()
	if !ok {
		return ast.Name{}, ast.Expr{}, false
	}
	value, ok = name.Init()
	return name, value, ok
}

// GetConditionalRoot returns the root node if a conditional statement or
// else null.
//
// Java: PeepholeRemoveDeadCode.getConditionalRoot(Node)
func (p *PeepholeRemoveDeadCode) GetConditionalRoot(n rhino.Node) rhino.Node {
	// We defined a conditional statement to be a IF or EXPR_RESULT rooted with
	// a HOOK, AND, or OR node.
	if n.IsNil() {
		return rhino.NilNode
	}
	if n.IsIf() {
		return n
	} else if PeepholeRemoveDeadCode_IsExprConditional(n) {
		return n.GetFirstChild()
	}
	return rhino.NilNode
}

// PeepholeRemoveDeadCode_IsExprConditional returns whether the node is a
// rooted with a HOOK, AND, OR, or COALESCE node.
//
// Java: PeepholeRemoveDeadCode.isExprConditional(Node)
func PeepholeRemoveDeadCode_IsExprConditional(n rhino.Node) bool {
	if n.IsExprResult() {
		switch n.GetFirstChild().GetToken() {
		case rhino.Token_HOOK, rhino.Token_AND, rhino.Token_OR, rhino.Token_COALESCE:
			return true
		default:
		}
	}
	return false
}

// GetConditionalStatementCondition returns the condition of a conditional
// statement.
//
// Java: PeepholeRemoveDeadCode.getConditionalStatementCondition(Node)
func (p *PeepholeRemoveDeadCode) GetConditionalStatementCondition(n rhino.Node) rhino.Node {
	if ifs, ok := ast.AsIf(n); ok {
		return ifs.Cond().Node()
	} else {
		jutil.CheckState(PeepholeRemoveDeadCode_IsExprConditional(n))
		// The first operand of the HOOK, AND, OR or COALESCE: the tested one.
		return n.GetFirstFirstChild()
	}
}

// TryFoldIf tries folding IF nodes by removing dead branches.
//
// Returns the replacement node, if changed, or the original if not.
//
// Java: PeepholeRemoveDeadCode.tryFoldIf(Node)
func (p *PeepholeRemoveDeadCode) TryFoldIf(n rhino.Node) rhino.Node {
	ifs, ok := ast.AsIf(n)
	jutil.CheckStateMsg1(ok, "%s", n)
	parent := n.GetParent()
	rhino.CheckNotNull(parent)
	cond := ifs.Cond().Node()
	thenBody := ifs.Then()
	elseBody, hasElse := ifs.Else()

	// if (x) { .. } else { } --> if (x) { ... }
	if hasElse && !p.MayHaveSideEffects(elseBody.Node()) {
		ifs.RemoveElse()
		p.ReportChangeToEnclosingScope(n)
		hasElse = false
	}

	// if (x) { } else { ... } --> if (!x) { ... }
	if !p.MayHaveSideEffects(thenBody.Node()) && hasElse {
		ifs.RemoveElse()
		ifs.SetThen(elseBody)
		notCond := rhino.NewNode(rhino.Token_NOT)
		cond.ReplaceWith(notCond)
		p.ReportChangeToEnclosingScope(n)
		notCond.AddChildToFront(cond)
		cond = notCond
		thenBody = elseBody
		hasElse = false
	}

	// `if (x()) { }` or `if (x?.()) { }`
	if !p.MayHaveSideEffects(thenBody.Node()) && !hasElse {
		if p.MayHaveSideEffects(cond) {
			// `x()` or `x?.()` has side effects, just leave the condition on its own.
			cond.Detach()
			replacement := NodeUtil_NewExpr(cond)
			n.ReplaceWith(replacement)
			p.ReportChangeToEnclosingScope(parent)
			return replacement
		} else {
			// `x()` or `x?.()` has no side effects, the whole tree is useless now.
			NodeUtil_RemoveChild(parent, n)
			p.ReportChangeToEnclosingScope(parent)
			return rhino.NilNode
		}
	}

	// Try transforms that apply to both IF and HOOK.
	condValue := NodeUtil_GetBooleanValue(cond)
	if condValue == base.Tri_UNKNOWN {
		return n // We can't remove branches otherwise!
	}

	if p.MayHaveSideEffects(cond) {
		// Transform "if (a = 2) {x =2}" into "if (true) {a=2;x=2}"
		newConditionValue := condValue == base.Tri_TRUE
		// Add an elseBody if it is needed.
		if !newConditionValue && !hasElse {
			elseBody = ast.NewBlock()
			elseBody.Node().Srcref(n)
			ifs.SetElse(elseBody)
			hasElse = true
		}
		cond.ReplaceWith(NodeUtil_BooleanNode(newConditionValue))
		branchToKeep := thenBody
		if !newConditionValue {
			branchToKeep = elseBody
		}
		branchToKeep.Node().AddChildToFront(rhino.IR_ExprResult(cond).Srcref(cond))
		p.ReportChangeToEnclosingScope(branchToKeep.Node())
	}

	condTrue := condValue.ToBoolean(true)
	if !hasElse {
		if condTrue {
			// Replace "if (true) { X }" with "X".
			thenStmt := thenBody.Node()
			thenStmt.Detach()
			n.ReplaceWith(thenStmt)
			p.ReportChangeToEnclosingScope(thenStmt)
			return thenStmt
		} else {
			// Remove "if (false) { X }" completely.
			NodeUtil_RedeclareVarsInsideBranch(n)
			NodeUtil_RemoveChild(parent, n)
			p.ReportChangeToEnclosingScope(parent)
			p.MarkFunctionsDeleted(n)
			return rhino.NilNode
		}
	} else {
		// Replace "if (true) { X } else { Y }" with X, or
		// replace "if (false) { X } else { Y }" with Y.
		branchToKeep, branchToRemove := thenBody, elseBody
		if !condTrue {
			branchToKeep, branchToRemove = elseBody, thenBody
		}
		NodeUtil_RedeclareVarsInsideBranch(branchToRemove.Node())
		branchToKeep.Node().Detach()
		n.ReplaceWith(branchToKeep.Node())
		p.ReportChangeToEnclosingScope(branchToKeep.Node())
		p.MarkFunctionsDeleted(n)
		return branchToKeep.Node()
	}
}

// TryFoldHook tries folding HOOK (?:) if the condition results of the
// condition is known.
//
// Returns the replacement node, if changed, or the original if not.
//
// Java: PeepholeRemoveDeadCode.tryFoldHook(Node)
func (p *PeepholeRemoveDeadCode) TryFoldHook(n rhino.Node) rhino.Node {
	hook, ok := ast.AsHook(n)
	jutil.CheckStateMsg1(ok, "%s", n)
	parent := n.GetParent()
	rhino.CheckNotNull(parent)
	cond := hook.Cond().Node()
	thenBody := hook.Then().Node()
	elseBody := hook.Else().Node()

	condValue := NodeUtil_GetBooleanValue(cond)
	if condValue == base.Tri_UNKNOWN {
		// If the result nodes are equivalent, then one of the nodes can be
		// removed and it doesn't matter which.
		if !p.AreNodesEqualForInlining(thenBody, elseBody) {
			return n // We can't remove branches otherwise!
		}
	}

	// Transform "(a = 2) ? x =2 : y" into "a=2,x=2"
	var branchToKeep rhino.Node
	var branchToRemove rhino.Node
	if condValue.ToBoolean(true) {
		branchToKeep = thenBody
		branchToRemove = elseBody
	} else {
		branchToKeep = elseBody
		branchToRemove = thenBody
	}

	var replacement rhino.Node
	condHasSideEffects := p.MayHaveSideEffects(cond)
	// Must detach after checking for side effects, to ensure that the parents
	// of nodes are set correctly.
	n.DetachChildren()
	if condHasSideEffects {
		replacement = rhino.IR_Comma(cond, branchToKeep).Srcref(n)
	} else {
		replacement = branchToKeep
		p.MarkFunctionsDeleted(cond)
	}

	n.ReplaceWith(replacement)
	p.ReportChangeToEnclosingScope(replacement)
	p.MarkFunctionsDeleted(branchToRemove)
	return replacement
}

// TryFoldFor removes FORs that always evaluate to false.
//
// Java: PeepholeRemoveDeadCode.tryFoldFor(Node)
func (p *PeepholeRemoveDeadCode) TryFoldFor(n rhino.Node) rhino.Node {
	forNode, ok := ast.AsFor(n)
	jutil.CheckArgument(ok)

	// Read the parts first: removing the init or the update below shifts the
	// children that follow it.
	init, hasInit := forNode.Init()
	update, hasUpdate := forNode.Update()
	body := forNode.Body()

	if hasInit {
		if _, isDecl := init.AsNameDeclaration(); !isDecl && p.TrySimplifyUnusedResult(init.Node()).IsNil() {
			// The init was removed: put an EMPTY in its place.
			n.AddChildToFront(rhino.IR_Empty().Srcref(n))
		}
	}

	if hasUpdate && p.TrySimplifyUnusedResult(update.Node()).IsNil() {
		// The update was removed: put an EMPTY in its place, before the body.
		rhino.IR_Empty().Srcref(n).InsertBefore(body.Node())
	}

	// There is an initializer skip it
	if _, ok := forNode.Init(); ok {
		return n
	}

	condExpr, ok := forNode.Cond()
	if !ok || NodeUtil_GetBooleanValue(condExpr.Node()) != base.Tri_FALSE {
		return n
	}
	cond := condExpr.Node()

	parent := n.GetParent()
	NodeUtil_RedeclareVarsInsideBranch(n)

	if !p.MayHaveSideEffects(cond) {
		// Remove the entire loop and any associated labels.
		for parent.IsLabel() {
			n = parent
			parent = parent.GetParent()
		}
		n.Detach()
	} else {
		statement := rhino.IR_ExprResult(cond.Detach()).SrcrefIfMissing(cond)
		if parent.IsLabel() {
			block := rhino.IR_Block()
			block.SrcrefIfMissing(statement)
			block.AddChildToFront(statement)
			statement = block
		}
		n.ReplaceWith(statement)
	}
	p.ReportChangeToEnclosingScope(parent)

	return rhino.NilNode
}

// TryFoldDoAway removes DOs that always evaluate to false. This leaves the
// statements that were in the loop in a BLOCK node. The block will be removed
// in a later pass, if possible.
//
// Java: PeepholeRemoveDeadCode.tryFoldDoAway(Node)
func (p *PeepholeRemoveDeadCode) TryFoldDoAway(n rhino.Node) rhino.Node {
	do, ok := ast.AsDo(n)
	jutil.CheckArgument(ok)

	cond := do.Cond().Node()
	if NodeUtil_GetBooleanValue(cond) != base.Tri_FALSE {
		return n
	}

	block := do.Body().Node()
	if n.GetParent().IsLabel() || PeepholeRemoveDeadCode_HasUnnamedBreakOrContinue(block) {
		return n
	}

	parent := n.GetParent()
	n.ReplaceWith(block.Detach())
	if p.MayHaveSideEffects(cond) {
		condStatement := rhino.IR_ExprResult(cond.Detach()).Srcref(cond)
		condStatement.InsertAfter(block)
	}
	p.ReportChangeToEnclosingScope(parent)

	return block
}

// TryFoldEmptyDo removes DOs that have empty bodies into FORs, which are
// much easier for the CFA to analyze.
//
// Java: PeepholeRemoveDeadCode.tryFoldEmptyDo(Node)
func (p *PeepholeRemoveDeadCode) TryFoldEmptyDo(n rhino.Node) rhino.Node {
	do, ok := ast.AsDo(n)
	jutil.CheckArgument(ok)

	// (Java also checks that the body is a block: it always is.)
	body := do.Body()
	if !body.HasStatements() {
		cond := do.Cond().Node()
		forNode :=
			rhino.IR_ForNode(rhino.IR_Empty().Srcref(n), cond.Detach(), rhino.IR_Empty().Srcref(n), body.Node().Detach())
		n.ReplaceWith(forNode)
		p.ReportChangeToEnclosingScope(forNode)
		return forNode
	}
	return n
}

// TryOptimizeObjectPattern removes string keys with an empty pattern as
// their child.
//
// Java: PeepholeRemoveDeadCode.tryOptimizeObjectPattern(Node)
func (p *PeepholeRemoveDeadCode) TryOptimizeObjectPattern(n rhino.Node) rhino.Node {
	pattern, ok := ast.AsObjectPattern(n)
	jutil.CheckArgumentMsg1(ok, "%s", n)

	if last, ok := pattern.LastProp(); ok {
		if _, isRest := last.AsObjectRest(); isRest {
			// don't remove any elements in `const {f: [], ...rest} = obj` because that affects what's
			// assigned to `rest`. only the last element can be object rest.
			return n
		}
	}

	// remove trailing EMPTY nodes and empty destructuring patterns
	for prop, ok := pattern.FirstProp(); ok; {
		next, hasNext := prop.Next() // read it first: prop may be removed

		// don't try to remove rest or computed properties, since they might have side effects
		if key, isKey := prop.AsStringKey(); isKey && p.IsRemovableDestructuringTarget(key.Value().Node()) {
			// e.g. `const {f: {}} = obj;`
			key.Node().Detach()
			p.ReportChangeToEnclosingScope(n)
		}
		prop, ok = next, hasNext
	}
	return n
}

// TryOptimizeArrayPattern removes trailing EMPTY nodes and empty array
// patterns.
//
// Java: PeepholeRemoveDeadCode.tryOptimizeArrayPattern(Node)
func (p *PeepholeRemoveDeadCode) TryOptimizeArrayPattern(n rhino.Node) rhino.Node {
	pattern, ok := ast.AsArrayPattern(n)
	jutil.CheckArgumentMsg1(ok, "%s", n)

	for last, ok := pattern.LastElement(); ok; last, ok = pattern.LastElement() {
		if _, isHole := last.AsEmpty(); !isHole && !p.IsRemovableDestructuringTarget(last.Node()) {
			// don't remove any non-trailing empty nodes because that will change the ordering of the
			// other assignments
			// note that this case also covers array pattern rest, which must be the final element
			break
		}
		last.Node().Detach()
		p.ReportChangeToEnclosingScope(n)
	}
	return n
}

// IsRemovableDestructuringTarget reports whether the element of a pattern
// (an array pattern element, or the value of an object pattern key) is an
// empty pattern, with a default value without side effects if it has one.
//
// Java: PeepholeRemoveDeadCode.isRemovableDestructuringTarget(Node)
func (p *PeepholeRemoveDeadCode) IsRemovableDestructuringTarget(destructuringElement rhino.Node) bool {
	target := destructuringElement
	var defaultValue rhino.Node // the nil node if there is none
	if dv, ok := ast.AsDefaultValue(destructuringElement); ok {
		target = dv.Target().Node()
		defaultValue = dv.Value().Node()
	}
	if pattern, ok := ast.AsPattern(target); !ok || pattern.Node().HasChildren() {
		return false
	}
	// only remove default values without side effects
	return defaultValue.IsNil() || !p.MayHaveSideEffects(defaultValue)
}

// PeepholeRemoveDeadCode_HasUnnamedBreakOrContinue returns whether a node
// has any unhandled breaks or continue.
//
// Java: PeepholeRemoveDeadCode.hasUnnamedBreakOrContinue(Node)
func PeepholeRemoveDeadCode_HasUnnamedBreakOrContinue(n rhino.Node) bool {
	return NodeUtil_Has(
		n,
		// Check for unlabeled breaks
		peepholeRemoveDeadCode_isUnnamedBreak,
		// ...inside contexts that can contain breaks.
		func(node rhino.Node) bool {
			return !rhino.IR_MayBeExpression(node) && // Functions are not visited
				!NodeUtil_IsLoopStructure(node) &&
				!node.IsSwitch()
		}) ||
		NodeUtil_Has(
			n,
			// Check for unlabeled continues
			peepholeRemoveDeadCode_isUnnamedContinue,
			// ...inside contexts that can contain continues.
			func(node rhino.Node) bool {
				return !rhino.IR_MayBeExpression(node) && // Functions are not visited
					!NodeUtil_IsLoopStructure(node)
			})
}

// TryFoldForCondition removes always true loop conditions.
//
// Java: PeepholeRemoveDeadCode.tryFoldForCondition(Node)
func (p *PeepholeRemoveDeadCode) TryFoldForCondition(forCondition rhino.Node) {
	if p.GetSideEffectFreeBooleanValue(forCondition) == base.Tri_TRUE {
		p.ReportChangeToEnclosingScope(forCondition)
		forCondition.ReplaceWith(rhino.IR_Empty())
	}
}

// TryRemoveOptionalChaining removes an optional chain whose object (or
// callee) is null or undefined.
//
// Go-only: the three optional chain kinds have no common typed view (the
// first child is the object of a property access but the callee of a
// call), so this uses the rhino API.
//
// Java: PeepholeRemoveDeadCode.tryRemoveOptionalChaining(Node)
func (p *PeepholeRemoveDeadCode) TryRemoveOptionalChaining(optionalChain rhino.Node) rhino.Node {
	callee := optionalChain.GetFirstChild()
	if !NodeUtil_IsNullOrUndefined(callee) {
		return optionalChain
	}
	var result rhino.Node
	if p.MayHaveSideEffects(callee) {
		// Simplify `(void sideEffectFunction())?.()` to `(void sideEffectFunction())`
		// The optional chain call won't execute but sideEffectFunction() is still evaluated.
		optionalChain.ReplaceWith(callee.Detach())
		result = callee
	} else {
		// Remove `(void 0)?.()` and (null)?.() and simplify `(void 0)?.x and null?.x` to void 0
		result = NodeUtil_NewUndefinedNode(callee)
		optionalChain.ReplaceWith(result)
	}
	p.MarkFunctionsDeleted(optionalChain)
	p.ReportChangeToEnclosingScope(result)
	return result
}
