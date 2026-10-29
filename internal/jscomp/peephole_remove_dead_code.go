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
}

// Java: PeepholeRemoveDeadCode()
func NewPeepholeRemoveDeadCode() *PeepholeRemoveDeadCode {
	return &PeepholeRemoveDeadCode{}
}

var peepholeRemoveDeadCodeTokens = peepholeTokenSet(
	rhino.Token_ASSIGN,
	rhino.Token_COMMA,
	rhino.Token_SCRIPT,
	rhino.Token_BLOCK,
	rhino.Token_EXPR_RESULT,
	rhino.Token_HOOK,
	rhino.Token_SWITCH,
	rhino.Token_IF,
	rhino.Token_FOR,
	rhino.Token_DO,
	rhino.Token_TRY,
	rhino.Token_LABEL,
	rhino.Token_ARRAY_PATTERN,
	rhino.Token_OBJECT_PATTERN,
	rhino.Token_VAR,
	rhino.Token_CONST,
	rhino.Token_LET,
	rhino.Token_DEFAULT_VALUE,
	rhino.Token_OPTCHAIN_GETPROP,
	rhino.Token_OPTCHAIN_CALL,
	rhino.Token_OPTCHAIN_GETELEM,
)

func (p *PeepholeRemoveDeadCode) peepholeConfig() peepholeOptConfig {
	return peepholeOptConfig{
		tokens:     peepholeRemoveDeadCodeTokens,
		isStandard: true,
	}
}

func (p *PeepholeRemoveDeadCode) cloneForWorker() AbstractPeepholeOptimization {
	return NewPeepholeRemoveDeadCode()
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
	if name, ok := val.AsName(); ok && name.Node().GetStringID() == nodeUtilIDUndefined {
		removeVal = true
	}

	// If the `void` application is pure, remove the value
	if void, ok := val.AsUnaryOf(rhino.Token_VOID); ok {
		removeVal = !p.MayHaveSideEffects(void.Operand().Node())
	}

	if removeVal {
		n.ReplaceWith(lValue.Detach())
		p.ReportChangeToEnclosingScope(lValue)
		return lValue
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
	if leftIsName && rightIsName && leftName.Node().GetStringID() == rightName.Node().GetStringID() {
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
//
// Java: PeepholeRemoveDeadCode.trySimplifyUnusedResult(Node)
func (p *PeepholeRemoveDeadCode) TrySimplifyUnusedResult(expression rhino.Node) rhino.Node {
	var rootsBuf [8]rhino.Node
	sideEffectRoots, atFixedPoint := p.trySimplifyUnusedResultInternal(expression, rootsBuf[:0])

	if atFixedPoint {
		// `expression` is in a form that cannot be further optimized.
		return expression
	}
	if len(sideEffectRoots) == 0 {
		p.DeleteNode(expression)
		return rhino.NilNode
	}
	if sideEffectRoots[0] == expression {
		// Expression was a conditional that was transformed. There can't be any other side-effects,
		// but we also can't detach the transformed root.
		if len(sideEffectRoots) != 1 {
			panic(jutil.NewIllegalStateExceptionMsg(jutil.CollectionToString(sideEffectRoots)))
		}
		p.ReportChangeToEnclosingScope(expression)
		return expression
	}
	sideEffects := PeepholeRemoveDeadCode_AsDetachedExpression(sideEffectRoots[0])

	// Assemble a tree of comma expressions for all the side-effects. The tree must execute the
	// side-effects in FIFO order with respect to the queue. It must also be left leaning to match
	// the parser's preferred structure.
	for _, root := range sideEffectRoots[1:] {
		next := PeepholeRemoveDeadCode_AsDetachedExpression(root)
		sideEffects = rhino.IR_Comma(sideEffects, next).Srcref(next)
	}

	sideEffects.InsertBefore(expression)
	p.DeleteNode(expression)
	return sideEffects
}

// trySimplifyUnusedResultInternal collects any potentially side-effectful
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
func (p *PeepholeRemoveDeadCode) trySimplifyUnusedResultInternal(tree rhino.Node, sideEffectRoots []rhino.Node) ([]rhino.Node, bool) {
	// Special cases for conditional expressions that may be using results.
	switch tree.GetToken() {
	case rhino.Token_OPTCHAIN_GETELEM, rhino.Token_OPTCHAIN_GETPROP, rhino.Token_OPTCHAIN_CALL:
		if p.MayHaveSideEffects(tree) {
			return append(sideEffectRoots, tree), PeepholeRemoveDeadCode_HasFixedPointParent(tree)
		}
		return sideEffectRoots, false
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
			return append(sideEffectRoots, tree), false // The node type was changed.
		}
		if !trueNode.IsNil() && falseNode.IsNil() {
			jutil.CheckStateMsg1(tree.HasTwoChildren(), "%s", tree)
			tree.SetToken(rhino.Token_AND)
			return append(sideEffectRoots, tree), false // The node type was changed.
		}
		if trueNode.IsNil() && falseNode.IsNil() {
			// Don't bother adding true and false branch children to make the AST valid; this HOOK is
			// going to be deleted. We just need to collect any side-effects from the predicate
			// expression.
			sideEffectRoots, _ = p.trySimplifyUnusedResultInternal(cond.Node(), sideEffectRoots)
			return sideEffectRoots, false // This HOOK must be cleaned up.
		}
		return append(sideEffectRoots, tree), PeepholeRemoveDeadCode_HasFixedPointParent(tree)
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
		if p.TrySimplifyUnusedResult(right.Node()).IsNil() {
			// Don't bother adding a second child to make the AST valid; this op is going to be
			// deleted. We just need to collect any side-effects from the predicate first child.
			sideEffectRoots, _ = p.trySimplifyUnusedResultInternal(left.Node(), sideEffectRoots)
			return sideEffectRoots, false // This op must be cleaned up.
		}
		return append(sideEffectRoots, tree), PeepholeRemoveDeadCode_HasFixedPointParent(tree)
	case rhino.Token_FUNCTION:
		// Functions that aren't being invoked are dead. If they were invoked we'd see the CALL
		// before arriving here. We don't want to look at any children since they'll never execute.
	case rhino.Token_MEMBER_FIELD_DEF:
		field := ast.MustMemberFieldDef(tree)
		if value, ok := field.Value(); ok && tree.IsStaticMember() {
			sideEffectRoots, _ = p.trySimplifyUnusedResultInternal(value.Node(), sideEffectRoots)
		}
	case rhino.Token_COMPUTED_FIELD_DEF:
		field := ast.MustComputedFieldDef(tree)
		value, hasValue := field.Value()
		// Computed key initializers run once when the class is defined, even on non-static fields.
		sideEffectRoots, _ = p.trySimplifyUnusedResultInternal(field.Key().Node(), sideEffectRoots)
		if hasValue && tree.IsStaticMember() {
			sideEffectRoots, _ = p.trySimplifyUnusedResultInternal(value.Node(), sideEffectRoots)
		}
	default:
		// This is the meat of this function. It covers the general case of nodes which are unused
		if p.NodeTypeMayHaveSideEffects(tree) {
			return append(sideEffectRoots, tree), PeepholeRemoveDeadCode_HasFixedPointParent(tree)
		}
		if !tree.HasChildren() {
			return sideEffectRoots, false // A node must have children or side-effects to be at fixed-point.
		}

		atFixedPoint := PeepholeRemoveDeadCode_HasFixedPointParent(tree)
		for child := tree.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
			var childFixed bool
			sideEffectRoots, childFixed = p.trySimplifyUnusedResultInternal(child, sideEffectRoots)
			atFixedPoint = childFixed && atFixedPoint
		}
		return sideEffectRoots, atFixedPoint
	}
	return sideEffectRoots, false
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
	var varsBuf [4]rhino.Node
	blockScopedVars := varsBuf[:0]
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
	}
	jutil.CheckState(PeepholeRemoveDeadCode_IsExprConditional(n))
	// The first operand of the HOOK, AND, OR or COALESCE: the tested one.
	return n.GetFirstFirstChild()
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
