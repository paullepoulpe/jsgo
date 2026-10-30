// Port of com/google/javascript/jscomp/PeepholeMinimizeConditions.java @ closure-compiler v20261021

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
	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Go-only: this pass reads the shapes it rewrites (if, hook, &&, ||, the
// blocks of an if) through the typed AST view of package ast (ifs.Cond(),
// ifs.Else(), hook.Then(), block.OnlyStatement()...) instead of child
// positions. The Java isXBlock/getXBlock pairs are single comma-ok
// functions that return the typed statement or expression. The edits still
// use the rhino API, in the Java's order.

// PeepholeMinimizeConditions is a peephole optimization that minimizes
// conditional expressions according to De Morgan's laws. Also rewrites
// conditional statements as expressions by replacing them with HOOKs and
// short-circuit binary operators.
//
// Based on PeepholeSubstituteAlternateSyntax
type PeepholeMinimizeConditions struct {
	AbstractPeepholeOptimizationBase
	late      bool
	scratchMC MinimizedCondition
}

// NewPeepholeMinimizeConditions creates the optimization.
//
// late: When late is false, this mean we are currently running before most
// of the other optimizations. In this case we would avoid optimizations that
// would make the code harder to analyze (such as using string splitting,
// merging statements with commas, etc). When late is true, we would do
// anything to minimize for size.
//
// Java: PeepholeMinimizeConditions(boolean)
func NewPeepholeMinimizeConditions(late bool) *PeepholeMinimizeConditions {
	return &PeepholeMinimizeConditions{late: late}
}

var peepholeMinimizeConditionsTokens = peepholeTokenSet(
	rhino.Token_THROW,
	rhino.Token_RETURN,
	rhino.Token_NOT,
	rhino.Token_IF,
	rhino.Token_EXPR_RESULT,
	rhino.Token_HOOK,
	rhino.Token_WHILE,
	rhino.Token_DO,
	rhino.Token_FOR,
	rhino.Token_BLOCK,
)

func (p *PeepholeMinimizeConditions) peepholeConfig() peepholeOptConfig {
	return peepholeOptConfig{
		tokens:     peepholeMinimizeConditionsTokens,
		isStandard: true,
	}
}

func (p *PeepholeMinimizeConditions) cloneForWorker() AbstractPeepholeOptimization {
	return NewPeepholeMinimizeConditions(p.late)
}

// OptimizeSubtree tries to apply our various peephole minimizations on the
// passed in node.
//
// Java: PeepholeMinimizeConditions.optimizeSubtree(Node)
func (p *PeepholeMinimizeConditions) OptimizeSubtree(node rhino.Node) rhino.Node {
	switch node.GetToken() {
	case rhino.Token_THROW, rhino.Token_RETURN:
		result := p.TryRemoveRedundantExit(node)
		if result != node {
			return result
		}
		return p.TryReplaceExitWithBreak(node)
		// TODO(johnlenz): Maybe remove redundant BREAK and CONTINUE. Overlaps
		// with MinimizeExitPoints.
	case rhino.Token_NOT:
		not := ast.MustUnary(node)
		p.TryMinimizeCondition(not.Operand().Node())
		return p.TryMinimizeNot(node)
	case rhino.Token_IF:
		ifs := ast.MustIf(node)
		p.PerformConditionSubstitutions(ifs.Cond().Node())
		return p.TryMinimizeIf(ifs)
	case rhino.Token_EXPR_RESULT:
		expr := ast.MustExprStmt(node)
		inner := expr.Expr().Node()
		if peepholeExprResultCannotMinimize(inner) {
			return node
		}
		p.PerformConditionSubstitutions(inner)
		return p.TryMinimizeExprResult(expr)
	case rhino.Token_HOOK:
		hook := ast.MustHook(node)
		p.PerformConditionSubstitutions(hook.Cond().Node())
		return p.TryMinimizeHook(hook)
	case rhino.Token_WHILE:
		while := ast.MustWhile(node)
		p.TryMinimizeCondition(while.Cond().Node())
		return node
	case rhino.Token_DO:
		do := ast.MustDo(node)
		p.TryMinimizeCondition(do.Cond().Node())
		return node
	case rhino.Token_FOR:
		forNode := ast.MustFor(node)
		p.TryJoinForCondition(forNode)
		// (Go: for (;;) has no condition to minimize; the Java tries its EMPTY
		// placeholder, which has no boolean value.)
		if cond, ok := forNode.Cond(); ok {
			p.TryMinimizeCondition(cond.Node())
		}
		return node
	case rhino.Token_BLOCK:
		block := ast.MustBlock(node)
		return p.TryReplaceIf(block)
	default:
		return node
		// Nothing changed
	}
}

func peepholeExprResultCannotMinimize(expr rhino.Node) bool {
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
		rhino.Token_ASSIGN_COALESCE:
		target := expr.GetFirstChild()
		if target.IsName() {
			return true
		}
		if !NodeUtil_IsNormalGet(target) {
			return false
		}
		current := target.GetFirstChild()
		if NodeUtil_EvaluatesToLocalValue(current) {
			return false
		}
		for NodeUtil_IsNormalGet(current) {
			current = current.GetFirstChild()
		}
		return !NodeUtil_IsLiteralValue(current, true)
	case rhino.Token_INC,
		rhino.Token_DEC,
		rhino.Token_DELPROP,
		rhino.Token_YIELD,
		rhino.Token_AWAIT,
		rhino.Token_DYNAMIC_IMPORT:
		return true
	default:
		return false
	}
}

// Java: PeepholeMinimizeConditions.tryJoinForCondition(Node)
func (p *PeepholeMinimizeConditions) TryJoinForCondition(n ast.For) {
	if !p.late {
		return
	}

	block := n.Body()
	first, ok := block.FirstStatement()
	if !ok {
		return
	}
	maybeIf, ok := first.AsIf()
	if !ok {
		return
	}
	thenBlock := maybeIf.Then()
	maybeBreak, ok := thenBlock.FirstStatement()
	if !ok {
		return
	}
	brk, ok := maybeBreak.AsBreak()
	if !ok {
		return
	}
	if _, hasLabel := brk.Label(); hasLabel {
		return
	}

	ifCondition := maybeIf.Cond()
	// Preserve the IF ELSE expression is there is one.
	if els, ok := maybeIf.Else(); ok {
		maybeIf.Node().ReplaceWith(els.Node().Detach())
	} else {
		NodeUtil_RedeclareVarsInsideBranch(thenBlock.Node())
		maybeIf.Node().Detach()
	}

	ifCondition.Node().Detach()
	fixedIfCondition := ast.NewUnary(rhino.Token_NOT, ifCondition)
	fixedIfCondition.Node().Srcref(ifCondition.Node())

	// OK, join the IF expression with the FOR expression
	if forCondition, ok := n.Cond(); ok {
		replacement := rhino.NewNode(rhino.Token_AND)
		forCondition.Node().ReplaceWith(replacement)
		replacement.AddChildToBack(forCondition.Node())
		replacement.AddChildToBack(fixedIfCondition.Node())
		p.ReportChangeToEnclosingScope(replacement)
	} else {
		// for (init; ; update): the condition replaces the EMPTY placeholder.
		n.SetCond(fixedIfCondition.AsExpr())
		p.ReportChangeToEnclosingScope(fixedIfCondition.Node())
	}
}

// TryReplaceIf uses "return x?1:2;" in place of "if(x)return 1;return 2;".
//
// Java: PeepholeMinimizeConditions.tryReplaceIf(Node)
func (p *PeepholeMinimizeConditions) TryReplaceIf(n ast.Block) rhino.Node {
	// The next statement is read before the body runs: the body may detach
	// the current one.
	next, hasNext := n.FirstStatement()
	for hasNext {
		child := next
		next, hasNext = child.Next()
		ifs, ok := child.AsIf()
		if !ok {
			continue
		}
		cond := ifs.Cond().Node()
		thenBranch := ifs.Then()
		elseBranch, hasElse := ifs.Else()

		if hasNext && !hasElse && PeepholeMinimizeConditions_IsReturnBlock(thenBranch) {
			if nextIf, ok := next.AsIf(); ok {
				nextCond := nextIf.Cond().Node()
				nextThen := nextIf.Then()
				nextElse, nextHasElse := nextIf.Else()
				if p.AreNodesEqualForInlining(thenBranch.Node(), nextThen.Node()) {
					// Transform
					//   if (x) return 1; if (y) return 1;
					// to
					//   if (x||y) return 1;
					ifs.Node().Detach()
					ifs.Node().DetachChildren()
					newCond := rhino.NewNodeWithChild(rhino.Token_OR, cond)
					nextCond.ReplaceWith(newCond)
					newCond.AddChildToBack(nextCond)
					p.ReportChangeToEnclosingScope(newCond)
				} else if nextHasElse && p.AreNodesEqualForInlining(thenBranch.Node(), nextElse.Node()) {
					// Transform
					//   if (x) return 1; if (y) foo() else return 1;
					// to
					//   if (!x&&y) foo() else return 1;
					ifs.Node().Detach()
					ifs.Node().DetachChildren()
					newCond := rhino.NewNodeWithChild(rhino.Token_AND,
						rhino.IR_Not(cond).Srcref(cond))
					nextCond.ReplaceWith(newCond)
					newCond.AddChildToBack(nextCond)
					p.ReportChangeToEnclosingScope(newCond)
				}
			} else if elseExpr, ok := PeepholeMinimizeConditions_GetReturnExpression(next); ok {
				var thenExpr rhino.Node
				// if(x)return; return 1 -> return x?void 0:1
				if value, ok := PeepholeMinimizeConditions_GetBlockReturnExpression(thenBranch); ok {
					thenExpr = value.Node()
					thenExpr.Detach()
				} else {
					thenExpr = NodeUtil_NewUndefinedNode(ifs.Node())
				}

				cond.Detach()
				elseExpr.Node().Detach()

				returnNode := rhino.IR_ReturnNodeWithExpr(
					rhino.IR_Hook(cond, thenExpr, elseExpr.Node()).
						Srcref(ifs.Node()))
				ifs.Node().ReplaceWith(returnNode)
				next.Node().Detach()
				p.ReportChangeToEnclosingScope(n.Node())
				// everything else in the block is dead code.
				break
			}
		} else if hasElse && PeepholeMinimizeConditions_StatementMustExitParent(thenBranch.AsStatement()) {
			ifs.RemoveElse()
			elseBranch.Node().InsertAfter(ifs.Node())
			p.ReportChangeToEnclosingScope(n.Node())
		}
	}
	return n.Node()
}

// Java: PeepholeMinimizeConditions.statementMustExitParent(Node)
func PeepholeMinimizeConditions_StatementMustExitParent(n ast.Statement) bool {
	switch n.Node().GetToken() {
	case rhino.Token_THROW, rhino.Token_RETURN:
		return true
	case rhino.Token_BLOCK:
		block, _ := n.AsBlock() // the case proves the kind
		if child, ok := block.LastStatement(); ok {
			return PeepholeMinimizeConditions_StatementMustExitParent(child)
		}
		return false
	// TODO(johnlenz): handle TRY/FINALLY
	case rhino.Token_FUNCTION:
		return false
	default:
		return false
	}
}

// TryReplaceExitWithBreak replaces duplicate exits in control structures.
// If the node following the exit node expression has the same effect as exit
// node, the node can be replaced or removed. For example:
//
//	"while (a) {return f()} return f();" ==> "while (a) {break} return f();"
//	"while (a) {throw 'ow'} throw 'ow';" ==> "while (a) {break} throw 'ow';"
//
// n is an follow control exit expression (a THROW or RETURN node). Returns
// the replacement for n, or the original if no change was made.
//
// Java: PeepholeMinimizeConditions.tryReplaceExitWithBreak(Node)
func (p *PeepholeMinimizeConditions) TryReplaceExitWithBreak(n rhino.Node) rhino.Node {
	result := n.GetFirstChild()

	// Find the enclosing control structure, if any, that a "break" would exit
	// from.
	breakTarget := n
	for ; !cfaIsBreakTargetByID(breakTarget, 0, false); breakTarget = breakTarget.GetParent() {
		if breakTarget.IsFunction() || breakTarget.IsScript() {
			// No break target.
			return n
		}
	}

	follow := ControlFlowAnalysis_ComputeFollowNode(breakTarget)

	// Skip pass all the finally blocks because both the break and return will
	// also trigger all the finally blocks. However, the order of execution is
	// slightly changed. Consider:
	//
	// return a() -> finally { b() } -> return a()
	//
	// which would call a() first. However, changing the first return to a
	// break will result in calling b().

	prefinallyFollows := follow
	follow = PeepholeMinimizeConditions_SkipFinallyNodes(follow)

	innerFollow := ControlFlowAnalysis_ComputeFollowNode(n)
	if prefinallyFollows != follow || innerFollow != PeepholeMinimizeConditions_SkipFinallyNodes(innerFollow) {
		// There were finally clauses
		if !p.IsPure(result) {
			// Can't defer the exit
			return n
		}
	}

	if follow.IsNil() && (n.IsThrow() || !result.IsNil()) {
		// Can't complete remove a throw here or a return with a result.
		return n
	}

	// When follow is null, this mean the follow of a break target is the
	// end of a function. This means a break is same as return.
	if follow.IsNil() || p.AreMatchingExits(n, follow) {
		replacement := rhino.IR_BreakNode()
		n.ReplaceWith(replacement)
		p.ReportChangeToEnclosingScope(replacement)
		return replacement
	}

	return n
}

// TryRemoveRedundantExit removes duplicate exits. If the node following the
// exit node expression has the same effect as exit node, the node can be
// removed. For example: "if (a) {return f()} return f();" ==> "if (a) {}
// return f();" "if (a) {throw 'ow'} throw 'ow';" ==> "if (a) {} throw 'ow';"
//
// n is an follow control exit expression (a THROW or RETURN node). Returns
// the replacement for n, or the original if no change was made.
//
// Java: PeepholeMinimizeConditions.tryRemoveRedundantExit(Node)
func (p *PeepholeMinimizeConditions) TryRemoveRedundantExit(n rhino.Node) rhino.Node {
	exitExpr := n.GetFirstChild()

	follow := ControlFlowAnalysis_ComputeFollowNode(n)

	// Skip pass all the finally blocks because both the fall through and return
	// will also trigger all the finally blocks.
	prefinallyFollows := follow
	follow = PeepholeMinimizeConditions_SkipFinallyNodes(follow)
	if prefinallyFollows != follow {
		// There were finally clauses
		if !p.IsPure(exitExpr) {
			// Can't replace the return
			return n
		}
	}

	if follow.IsNil() && (n.IsThrow() || !exitExpr.IsNil()) {
		// Can't complete remove a throw here or a return with a result.
		return n
	}

	// When follow is null, this mean the follow of a break target is the
	// end of a function. This means a break is same as return.
	if follow.IsNil() || p.AreMatchingExits(n, follow) {
		p.ReportChangeToEnclosingScope(n)
		n.Detach()
		return rhino.NilNode
	}

	return n
}

// IsPure returns whether the expression does not produces and can not be
// affected by side-effects.
//
// Java: PeepholeMinimizeConditions.isPure(Node)
func (p *PeepholeMinimizeConditions) IsPure(n rhino.Node) bool {
	return n.IsNil() ||
		(!NodeUtil_CanBeSideEffected(n) &&
			!p.MayHaveSideEffects(n))
}

// PeepholeMinimizeConditions_SkipFinallyNodes returns n or the node
// following any following finally nodes.
//
// Java: PeepholeMinimizeConditions.skipFinallyNodes(Node)
func PeepholeMinimizeConditions_SkipFinallyNodes(n rhino.Node) rhino.Node {
	for !n.IsNil() && NodeUtil_IsTryFinallyNode(n.GetParent(), n) {
		n = ControlFlowAnalysis_ComputeFollowNode(n)
	}
	return n
}

// AreMatchingExits checks whether one exit can be replaced with another.
// Verify:
//  1. They are identical expressions
//  2. If an exception is possible that the statements, the original and the
//     potential replacement are in the same exception handler.
//
// Java: PeepholeMinimizeConditions.areMatchingExits(Node, Node)
func (p *PeepholeMinimizeConditions) AreMatchingExits(nodeThis rhino.Node, nodeThat rhino.Node) bool {
	if !p.IsASTNormalized() &&
		(nodeThis.IsThrow() || nodeThis.IsReturn()) &&
		nodeThis.HasChildren() {
		// if the ast isn't normalized "return a" or "throw a" may not mean the same thing in
		// different blocks.
		return false
	}
	return nodeThis.IsEquivalentTo(nodeThat) &&
		(!PeepholeMinimizeConditions_IsExceptionPossible(nodeThis) ||
			PeepholeMinimizeConditions_GetExceptionHandler(nodeThis) == PeepholeMinimizeConditions_GetExceptionHandler(nodeThat))
}

// Java: PeepholeMinimizeConditions.isExceptionPossible(Node)
func PeepholeMinimizeConditions_IsExceptionPossible(n rhino.Node) bool {
	// TODO(johnlenz): maybe use ControlFlowAnalysis.mayThrowException?
	if _, ok := ast.AsThrow(n); ok {
		return true
	}
	ret, ok := ast.AsReturn(n)
	jutil.CheckStateMsg1(ok, "%s", n)
	value, ok := ret.Value()
	return ok && !NodeUtil_IsLiteralValue(value.Node(), true)
}

// Java: PeepholeMinimizeConditions.getExceptionHandler(Node)
func PeepholeMinimizeConditions_GetExceptionHandler(n rhino.Node) rhino.Node {
	return ControlFlowAnalysis_GetExceptionHandler(n)
}

// TryRemoveRepeatedStatements tries to remove duplicate statements from IF
// blocks. For example:
//
//	if (a) {
//	  x = 1;
//	  return true;
//	} else {
//	  x = 2;
//	  return true;
//	}
//
// becomes:
//
//	if (a) {
//	  x = 1;
//	} else {
//	  x = 2;
//	}
//	return true;
//
// n is the IF node to examine.
//
// Java: PeepholeMinimizeConditions.tryRemoveRepeatedStatements(Node)
func (p *PeepholeMinimizeConditions) TryRemoveRepeatedStatements(n ast.If) {
	// Only run this if variable names are guaranteed to be unique. Otherwise bad things can happen:
	// see PeepholeMinimizeConditionsTest#testDontRemoveDuplicateStatementsWithoutNormalization
	if !p.IsASTNormalized() {
		return
	}

	parent := n.Node().GetParent()
	if !NodeUtil_IsStatementBlock(parent) {
		// If the immediate parent is something like a label, we
		// can't move the statement, so bail.
		return
	}

	trueBranch := n.Then()
	falseBranch, ok := n.Else()
	jutil.CheckState(ok)

	for {
		lastTrue, ok1 := trueBranch.LastStatement()
		lastFalse, ok2 := falseBranch.LastStatement()
		if !ok1 || !ok2 ||
			!p.AreNodesEqualForInlining(lastTrue.Node(), lastFalse.Node()) {
			break
		}
		lastTrue.Node().Detach()
		lastFalse.Node().Detach()
		lastTrue.Node().InsertAfter(n.Node())
		p.ReportChangeToEnclosingScope(parent)
	}
}

// PeepholeMinimizeConditions_IsReturnBlock returns whether the node is a
// block with a single statement that is an return with or without an
// expression.
//
// Java: PeepholeMinimizeConditions.isReturnBlock(Node)
func PeepholeMinimizeConditions_IsReturnBlock(n ast.Block) bool {
	first, ok := n.OnlyStatement()
	if !ok {
		return false
	}
	_, ok = first.AsReturn()
	return ok
}

// PeepholeMinimizeConditions_GetBlockReturnExpression returns the
// expression of n if n is a block with a single statement that is a return
// with an expression; ok is false otherwise.
//
// Java: PeepholeMinimizeConditions.isReturnExpressBlock(Node) and
// getBlockReturnExpression(Node) (Go: one comma-ok function).
func PeepholeMinimizeConditions_GetBlockReturnExpression(n ast.Block) (ast.Expr, bool) {
	first, ok := n.OnlyStatement()
	if !ok {
		return ast.Expr{}, false
	}
	return PeepholeMinimizeConditions_GetReturnExpression(first)
}

// PeepholeMinimizeConditions_GetReturnExpression returns the expression of
// n if n is a return statement with an expression; ok is false otherwise.
//
// Java: PeepholeMinimizeConditions.isReturnExpression(Node) (Go: comma-ok,
// with the expression the callers read next).
func PeepholeMinimizeConditions_GetReturnExpression(n ast.Statement) (ast.Expr, bool) {
	ret, ok := n.AsReturn()
	if !ok {
		return ast.Expr{}, false
	}
	return ret.Value()
}

// PeepholeMinimizeConditions_GetBlockVar returns the var statement of n and
// the name it declares if n is a block with a single statement that is a
// VAR declaration of a single variable; ok is false otherwise.
//
// Java: PeepholeMinimizeConditions.isVarBlock(Node) and getBlockVar(Node)
// (Go: one comma-ok function. A var that declares a destructuring pattern
// is not one: the callers compare the declared name, where the Java would
// throw.)
func PeepholeMinimizeConditions_GetBlockVar(n ast.Block) (ast.NameDeclaration, ast.Name, bool) {
	first, ok := n.OnlyStatement()
	if !ok {
		return ast.NameDeclaration{}, ast.Name{}, false
	}
	decl, ok := first.AsNameDeclaration()
	if !ok || decl.Token() != rhino.Token_VAR {
		return ast.NameDeclaration{}, ast.Name{}, false
	}
	declarator, ok := decl.OnlyDeclarator()
	if !ok {
		return ast.NameDeclaration{}, ast.Name{}, false
	}
	name, ok := declarator.AsName()
	return decl, name, ok
}

// PeepholeMinimizeConditions_ConsumesDanglingElse reports whether a statement
// consumes a 'dangling else'. A statement consumes a 'dangling else' if an
// 'else' token following the statement would be considered by the parser to
// be part of the statement.
//
// Java: PeepholeMinimizeConditions.consumesDanglingElse(Node)
func PeepholeMinimizeConditions_ConsumesDanglingElse(n ast.Statement) bool {
	for {
		switch n.Node().GetToken() {
		case rhino.Token_IF:
			ifs, _ := n.AsIf() // the case proves the kind
			els, ok := ifs.Else()
			if !ok {
				return true
			}
			// This IF node has an else clause: look into it.
			n = els.AsStatement()
		case rhino.Token_BLOCK:
			block, _ := n.AsBlock() // the case proves the kind
			only, ok := block.OnlyStatement()
			if !ok {
				return false
			}
			// This BLOCK has no curly braces.
			n = only
		case rhino.Token_WITH:
			with, _ := n.AsWith() // the case proves the kind
			n = with.Body().AsStatement()
		case rhino.Token_WHILE, rhino.Token_FOR, rhino.Token_FOR_IN:
			loop, _ := n.AsLoop() // the case proves the kind
			n = loop.Body().AsStatement()
		default:
			return false
		}
	}
}
