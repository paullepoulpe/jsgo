// Port of com/google/javascript/jscomp/NodeUtil.java @ closure-compiler v20261021

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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// NodeUtil_GetInsertionPointAfterAllInnerFunctionDeclarations returns the
// node within the given block at which something can be inserted such that
// it's after all inner function declarations in that block. We want this
// because normalization expects all inner function declarations to be
// hoisted. If there's not good insertion point (e.g. the block is empty or
// only contains inner function declarations), return null.
//
// Java: NodeUtil.getInsertionPointAfterAllInnerFunctionDeclarations(Node)
func NodeUtil_GetInsertionPointAfterAllInnerFunctionDeclarations(block rhino.Node) rhino.Node {
	jutil.CheckState(block.IsBlock())
	current := block.GetFirstChild()

	for !current.IsNil() && NodeUtil_IsFunctionDeclaration(current) {
		current = current.GetNext()
	}
	return current
}

// NodeUtil_IsSomeCompileTimeConstStringValue returns true iff the value
// associated with the node is a JS string literal, a concatenation thereof or
// a ternary operator choosing between string literals.
//
// Java: NodeUtil.isSomeCompileTimeConstStringValue(Node)
func NodeUtil_IsSomeCompileTimeConstStringValue(node rhino.Node) bool {
	// TODO(bangert): Support constants, using a Scope argument. See ConstParamCheck
	if node.IsStringLit() || (node.IsTemplateLit() && node.HasOneChild()) {
		return true
	} else if node.IsAdd() {
		jutil.CheckStateMsg1(node.HasTwoChildren(), "%s", node)
		add := ast.MustBinary(node)
		return NodeUtil_IsSomeCompileTimeConstStringValue(add.Left().Node()) &&
			NodeUtil_IsSomeCompileTimeConstStringValue(add.Right().Node())
	} else if hook, ok := ast.AsHook(node); ok {
		// Ternary operator a ? b : c
		return NodeUtil_IsSomeCompileTimeConstStringValue(hook.Then().Node()) &&
			NodeUtil_IsSomeCompileTimeConstStringValue(hook.Else().Node())
	}
	return false
}

// Java: NodeUtil.getEnclosingNonArrowFunction(Node)
func NodeUtil_GetEnclosingNonArrowFunction(n rhino.Node) rhino.Node {
	return NodeUtil_GetEnclosingNode(n, NodeUtil_IsNonArrowFunction)
}

// NodeUtil_GetEnclosingClass finds the class containing the given node.
//
// Java: NodeUtil.getEnclosingClass(Node)
func NodeUtil_GetEnclosingClass(n rhino.Node) rhino.Node {
	return NodeUtil_GetEnclosingNode(n, (rhino.Node).IsClass)
}

// Java: NodeUtil.getEnclosingModuleIfPresent(Node)
func NodeUtil_GetEnclosingModuleIfPresent(n rhino.Node) rhino.Node {
	return NodeUtil_GetEnclosingNode(n, (rhino.Node).IsModuleBody)
}

// NodeUtil_GetEnclosingFunction finds the function containing the given
// node.
//
// Java: NodeUtil.getEnclosingFunction(Node)
func NodeUtil_GetEnclosingFunction(n rhino.Node) rhino.Node {
	return NodeUtil_GetEnclosingNode(n, (rhino.Node).IsFunction)
}

// NodeUtil_GetEnclosingScript finds the script containing the given node.
//
// Java: NodeUtil.getEnclosingScript(Node)
func NodeUtil_GetEnclosingScript(n rhino.Node) rhino.Node {
	return NodeUtil_GetEnclosingNode(n, (rhino.Node).IsScript)
}

// NodeUtil_GetEnclosingBlock finds the block containing the given node.
//
// Java: NodeUtil.getEnclosingBlock(Node)
func NodeUtil_GetEnclosingBlock(n rhino.Node) rhino.Node {
	return NodeUtil_GetEnclosingNode(n, (rhino.Node).IsBlock)
}

// Java: NodeUtil.getEnclosingBlockScopeRoot(Node)
func NodeUtil_GetEnclosingBlockScopeRoot(n rhino.Node) rhino.Node {
	return NodeUtil_GetEnclosingNode(n, NodeUtil_CreatesBlockScope)
}

// Java: NodeUtil.getEnclosingScopeRoot(Node)
func NodeUtil_GetEnclosingScopeRoot(n rhino.Node) rhino.Node {
	return NodeUtil_GetEnclosingNode(n, NodeUtil_CreatesScope)
}

// NodeUtil_GetEnclosingHoistScopeRoot returns the nearest enclosing hoist
// scope root node, null for the global scope. There are currently 4 such
// roots to consider: the global, module, function, and class static blocks.
//
// Java: NodeUtil.getEnclosingHoistScopeRoot(Node)
func NodeUtil_GetEnclosingHoistScopeRoot(n rhino.Node) rhino.Node {
	return NodeUtil_GetEnclosingNode(n, NodeUtil_IsHoistScopeRoot)
}

// Java: NodeUtil.isHoistScopeRoot(Node)
func NodeUtil_IsHoistScopeRoot(n rhino.Node) bool {
	// The "global" hoist scope has multiple meaning so
	return n.IsFunction() || n.IsModuleBody() || NodeUtil_IsClassStaticBlock(n)
}

// Java: NodeUtil.isInFunction(Node)
func NodeUtil_IsInFunction(n rhino.Node) bool {
	return !NodeUtil_GetEnclosingFunction(n).IsNil()
}

// Java: NodeUtil.getEnclosingStatement(Node)
func NodeUtil_GetEnclosingStatement(n rhino.Node) rhino.Node {
	return NodeUtil_GetEnclosingNode(n, NodeUtil_IsStatement)
}

// Java: NodeUtil.getEnclosingNode(Node, Predicate<Node>)
func NodeUtil_GetEnclosingNode(n rhino.Node, pred func(rhino.Node) bool) rhino.Node {
	curr := n
	for !curr.IsNil() && !pred(curr) {
		curr = curr.GetParent()
	}
	return curr
}

// NodeUtil_GetFirstPropMatchingKey returns the first property in the objlit
// or class members, that matches the key.
//
// Java: NodeUtil.getFirstPropMatchingKey(Node, String)
func NodeUtil_GetFirstPropMatchingKey(n rhino.Node, keyName string) rhino.Node {
	jutil.CheckState(n.IsObjectLit() || n.IsClassMembers())
	for keyNode := n.GetFirstChild(); !keyNode.IsNil(); keyNode = keyNode.GetNext() {
		if (keyNode.IsStringKey() || keyNode.IsMemberFunctionDef()) &&
			keyNode.GetString() == keyName {
			return ast.MustPlainKeyMember(keyNode).Value().Node()
		} else if comp, ok := ast.AsComputedProp(keyNode); ok {
			if keyLit, ok := ast.AsStringLit(comp.Key().Node()); ok && keyLit.Value() == keyName {
				if val, ok := comp.Value(); ok {
					return val.Node()
				}
				return comp.Key().Node()
			}
		}
	}
	return rhino.NilNode
}

// Java: NodeUtil.isAnyFor(Node)
func NodeUtil_IsAnyFor(n rhino.Node) bool {
	return n.IsVanillaFor() || n.IsForIn() || n.IsForOf() || n.IsForAwaitOf()
}

// NodeUtil_IsLoopStructure determines whether the given node is a FOR, DO,
// or WHILE node.
//
// Java: NodeUtil.isLoopStructure(Node)
func NodeUtil_IsLoopStructure(n rhino.Node) bool {
	_, ok := ast.AsLoop(n)
	return ok
}

// NodeUtil_GetLoopCodeBlock: n is the node to inspect. If the node, is a
// FOR, WHILE, or DO, it returns the node for the code BLOCK, null otherwise.
//
// Java: NodeUtil.getLoopCodeBlock(Node)
func NodeUtil_GetLoopCodeBlock(n rhino.Node) rhino.Node {
	if loop, ok := ast.AsLoop(n); ok {
		return loop.Body().Node()
	}
	return rhino.NilNode
}

// NodeUtil_IsWithinLoop returns whether the specified node has a loop parent
// that is within the current scope.
//
// Java: NodeUtil.isWithinLoop(Node)
func NodeUtil_IsWithinLoop(n rhino.Node) bool {
	for parent := range n.GetAncestors().All() {
		if NodeUtil_IsLoopStructure(parent) {
			return true
		}

		if parent.IsFunction() {
			break
		}
	}
	return false
}

// NodeUtil_IsControlStructure determines whether the given node is a FOR,
// DO, WHILE, WITH, or IF node.
//
// Java: NodeUtil.isControlStructure(Node)
func NodeUtil_IsControlStructure(n rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_FOR,
		rhino.Token_FOR_IN,
		rhino.Token_FOR_OF,
		rhino.Token_FOR_AWAIT_OF,
		rhino.Token_DO,
		rhino.Token_WHILE,
		rhino.Token_WITH,
		rhino.Token_IF,
		rhino.Token_LABEL,
		rhino.Token_TRY,
		rhino.Token_CATCH,
		rhino.Token_SWITCH,
		rhino.Token_CASE,
		rhino.Token_DEFAULT_CASE:
		return true
	default:
		return false
	}
}

// NodeUtil_IsControlStructureCodeBlock determines whether the given node is
// code node for FOR, DO, WHILE, WITH, or IF node.
//
// Java: NodeUtil.isControlStructureCodeBlock(Node, Node)
func NodeUtil_IsControlStructureCodeBlock(parent rhino.Node, n rhino.Node) bool {
	if parent.IsTry() {
		return parent.GetFirstChild() == n || parent.GetLastChild() == n
	}
	switch ast.RoleIn(n, parent) {
	case ast.RoleDoBody,
		ast.RoleForBody,
		ast.RoleForInBody,
		ast.RoleForOfBody,
		ast.RoleForAwaitOfBody,
		ast.RoleWhileBody,
		ast.RoleLabelBody,
		ast.RoleWithBody,
		ast.RoleCatchBody,
		ast.RoleIfThen,
		ast.RoleIfElse,
		ast.RoleSwitchBody,
		ast.RoleCaseBody,
		ast.RoleDefaultCaseBody:
		return true
	default:
		if parent.IsIf() || parent.IsSwitch() || parent.IsCase() {
			return parent.GetFirstChild() != n // exclude the condition
		}
		if parent.IsDefaultCase() {
			return true
		}
		jutil.CheckStateMsg1(NodeUtil_IsControlStructure(parent), "%s", parent)
		return false
	}
}

// NodeUtil_GetConditionExpression gets the condition of an ON_TRUE /
// ON_FALSE CFG edge.
//
// n is a node with an outgoing conditional CFG edge; returns the condition
// node or null if the condition is not obviously a node.
//
// Java: NodeUtil.getConditionExpression(Node)
func NodeUtil_GetConditionExpression(n rhino.Node) rhino.Node {
	switch n.GetToken() {
	case rhino.Token_IF:
		return ast.MustIf(n).Cond().Node()
	case rhino.Token_WHILE:
		return ast.MustWhile(n).Cond().Node()
	case rhino.Token_DO:
		return ast.MustDo(n).Cond().Node()
	case rhino.Token_FOR:
		return n.GetSecondChild()
	case rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF, rhino.Token_CASE:
		return rhino.NilNode
	default:
		panic(jutil.NewIllegalArgumentExceptionMsg(n.String() + " does not have a condition."))
	}
}

// NodeUtil_IsStatementBlock returns whether the node is of a type that
// contain other statements.
//
// Java: NodeUtil.isStatementBlock(Node)
func NodeUtil_IsStatementBlock(n rhino.Node) bool {
	return n.IsRoot() || n.IsScript() || n.IsBlock() || n.IsModuleBody()
}

// NodeUtil_CreatesBlockScope: a block scope is created by a non-synthetic
// block node, a for loop node, or a for-of loop node.
//
// Note: for functions, we use two separate scopes for parameters and
// declarations in the body. We need to make sure default parameters cannot
// reference var / function declarations in the body.
//
// Returns whether the node creates a block scope.
//
// Java: NodeUtil.createsBlockScope(Node)
func NodeUtil_CreatesBlockScope(n rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_BLOCK:
		parent := n.GetParent()
		if !parent.IsNil() && (parent.IsBlock() || parent.IsScript()) && n.IsSyntheticBlock() {
			// Don't create block scope for synthetic blocks.
			return false
		}
		// Don't create block scope for switch cases or catch blocks.
		return !parent.IsNil() && !NodeUtil_IsSwitchCase(parent) && !parent.IsCatch()
	case rhino.Token_FOR, rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF, rhino.Token_SWITCH_BODY,
		rhino.Token_CLASS:
		return true
	default:
		return false
	}
}

// Java: NodeUtil.createsScope(Node)
func NodeUtil_CreatesScope(n rhino.Node) bool {
	return NodeUtil_CreatesBlockScope(n) ||
		n.IsFunction() ||
		n.IsModuleBody() ||
		n.IsMemberFieldDef() ||
		n.IsComputedFieldDef() ||
		// The ROOT nodes that are the root of the externs tree or main JS tree do not
		// create scopes. The parent of those two, which is the root of the entire AST and
		// therefore has no parent, is the only ROOT node that creates a scope.
		(n.IsRoot() && n.GetParent().IsNil())
}

// Java: NodeUtil.isValidCfgRoot(Node)
func NodeUtil_IsValidCfgRoot(n rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_FUNCTION, rhino.Token_SCRIPT, rhino.Token_MODULE_BODY, rhino.Token_ROOT:
		return true
	default:
		return NodeUtil_IsClassStaticBlock(n)
	}
}

// NodeUtil_IsStatement returns whether the node is used as a statement.
//
// Java: NodeUtil.isStatement(Node)
func NodeUtil_IsStatement(n rhino.Node) bool {
	return !n.IsModuleBody() && !n.IsScript() && !n.IsRoot() && NodeUtil_IsStatementParent(n.GetParent())
}

// Java: NodeUtil.isStatementParent(Node)
func NodeUtil_IsStatementParent(parent rhino.Node) bool {
	// It is not possible to determine definitely if a node is a statement
	// or not if it is not part of the AST.  A FUNCTION node can be
	// either part of an expression or a statement.
	switch parent.GetToken() {
	case rhino.Token_SCRIPT,
		rhino.Token_MODULE_BODY,
		rhino.Token_BLOCK,
		rhino.Token_LABEL,
		rhino.Token_NAMESPACE_ELEMENTS,
		rhino.Token_INTERFACE_MEMBERS:
		return true
	default:
		return false
	}
}

// Java: NodeUtil.isDeclarationParent(Node)
func NodeUtil_IsDeclarationParent(parent rhino.Node) bool {
	switch parent.GetToken() {
	case rhino.Token_DECLARE, rhino.Token_EXPORT:
		return true
	default:
		return NodeUtil_IsStatementParent(parent)
	}
}

// NodeUtil_IsSwitchCase: whether the node is part of a switch statement.
//
// Java: NodeUtil.isSwitchCase(Node)
func NodeUtil_IsSwitchCase(n rhino.Node) bool {
	return n.IsCase() || n.IsDefaultCase()
}

// NodeUtil_IsReferenceName returns whether the node is a reference to a
// variable, function, class or function parameter (not a label or an empty
// function expression name).
//
// Java: NodeUtil.isReferenceName(Node)
func NodeUtil_IsReferenceName(n rhino.Node) bool {
	return n.IsName() && n.GetString() != ""
}

// NodeUtil_IsTryFinallyNode: whether the child node is the FINALLY block of
// a try.
//
// Java: NodeUtil.isTryFinallyNode(Node, Node)
func NodeUtil_IsTryFinallyNode(parent rhino.Node, child rhino.Node) bool {
	if tryNode, ok := ast.AsTry(parent); ok {
		fin, hasFin := tryNode.Finally()
		return hasFin && fin.Node() == child
	}
	return false
}

// NodeUtil_IsTryCatchNodeContainer: whether the node is a CATCH container
// BLOCK.
//
// Java: NodeUtil.isTryCatchNodeContainer(Node)
func NodeUtil_IsTryCatchNodeContainer(n rhino.Node) bool {
	if tryNode, ok := ast.AsTry(n.GetParent()); ok {
		return tryNode.Catches().Node() == n
	}
	return false
}

// Java: NodeUtil.canContainHoistedVarsDecls(Node)
func NodeUtil_CanContainHoistedVarsDecls(n rhino.Node) bool {
	if NodeUtil_IsStatement(n) {
		switch n.GetToken() {
		case rhino.Token_EXPR_RESULT, rhino.Token_RETURN, rhino.Token_THROW, rhino.Token_BREAK, rhino.Token_CONTINUE,
			rhino.Token_EMPTY, rhino.Token_DEBUGGER:
			return false
		default:
			return true // assume anything else can
		}
	}
	return true
}

// NodeUtil_CopyNameAnnotations copies any annotations that follow a named
// value.
//
// Java: NodeUtil.copyNameAnnotations(Node, Node)
func NodeUtil_CopyNameAnnotations(source rhino.Node, destination rhino.Node) {
	if source.GetBooleanProp(rhino.Node_IS_CONSTANT_NAME) {
		destination.PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, true)
	}
}

// NodeUtil_GetAddingRoot gets a Node at the top of the current scope where
// we can add new var declarations as children.
//
// Java: NodeUtil.getAddingRoot(Node)
func NodeUtil_GetAddingRoot(n rhino.Node) rhino.Node {
	var addingRoot rhino.Node
	ancestor := n
crawl_ancestors:
	for {
		ancestor = ancestor.GetParent()
		if ancestor.IsNil() {
			break
		}
		switch ancestor.GetToken() {
		case rhino.Token_SCRIPT, rhino.Token_MODULE_BODY:
			addingRoot = ancestor
			break crawl_ancestors
		case rhino.Token_FUNCTION:
			addingRoot = ancestor.GetLastChild()
			break crawl_ancestors
		default:
			continue crawl_ancestors
		}
	}

	// make sure that the adding root looks ok
	jutil.CheckState(addingRoot.IsBlock() || addingRoot.IsModuleBody() || addingRoot.IsScript())
	jutil.CheckState(!addingRoot.HasChildren() || !addingRoot.GetFirstChild().IsScript())
	return addingRoot
}

// Java: NodeUtil.newDeclaration(Node, Node, Token)
func NodeUtil_NewDeclaration(lhs rhino.Node, rhs rhino.Node, declarationType rhino.Token) rhino.Node {
	if rhs.IsNil() {
		return rhino.IR_Declaration(lhs, declarationType)
	}
	return rhino.IR_DeclarationWithValue(lhs, rhs, declarationType)
}

// NodeUtil_HasFinally returns whether a TRY node has a finally block.
//
// Java: NodeUtil.hasFinally(Node)
func NodeUtil_HasFinally(n rhino.Node) bool {
	jutil.CheckArgument(n.IsTry())
	_, ok := ast.MustTry(n).Finally()
	return ok
}

// NodeUtil_GetCatchBlock returns the BLOCK node containing the CATCH node
// (if any) of a TRY.
//
// Java: NodeUtil.getCatchBlock(Node)
func NodeUtil_GetCatchBlock(n rhino.Node) rhino.Node {
	jutil.CheckArgument(n.IsTry())
	return ast.MustTry(n).Catches().Node()
}

// NodeUtil_HasCatchHandler returns whether BLOCK (from a TRY node) contains
// a CATCH.
//
// See NodeUtil_GetCatchBlock.
//
// Java: NodeUtil.hasCatchHandler(Node)
func NodeUtil_HasCatchHandler(n rhino.Node) bool {
	jutil.CheckArgument(n.IsBlock())
	return n.HasChildren() && n.GetFirstChild().IsCatch()
}

// NodeUtil_IsExpressionResultUsed returns true iff the result of the
// expression is consumed.
//
// Java: NodeUtil.isExpressionResultUsed(Node)
func NodeUtil_IsExpressionResultUsed(expr rhino.Node) bool {
	parent := expr.GetParent()
	switch parent.GetToken() {
	case rhino.Token_BLOCK, rhino.Token_EXPR_RESULT:
		return false
	case rhino.Token_CAST:
		return NodeUtil_IsExpressionResultUsed(parent)
	case rhino.Token_HOOK:
		return ast.MustHook(parent).Cond().Node() == expr || NodeUtil_IsExpressionResultUsed(parent)
	case rhino.Token_AND, rhino.Token_OR, rhino.Token_COALESCE:
		return ast.MustBinary(parent).Left().Node() == expr || NodeUtil_IsExpressionResultUsed(parent)
	case rhino.Token_COMMA:
		comma := ast.MustBinary(parent)
		grandparent := parent.GetParent()
		if (grandparent.IsCall() || grandparent.IsTaggedTemplateLit()) &&
			ast.MustInvocation(grandparent).Callee().Node() == parent &&
			expr == comma.Left().Node() &&
			parent.HasTwoChildren() {
			// Special case the indirect function call pattern, e.g. (0, myFn)(arg1, arg2).
			// The indirect call pattern has two use cases:
			calledFn := comma.Right().Node()
			if calledFn.IsName() && calledFn.MatchesName("eval") {
				// 1) eval
				//    Semantically, a direct call to eval is different from an indirect
				//    call to an eval. See ECMA-262 S15.1.2.1.
				return true
			}
			if NodeUtil_IsNormalOrOptChainGet(calledFn) {
				// 2) Removing "this" context
				//    When calling prefix.myFn(), myFn receives "prefix" as its this context.
				//    Calling "(0, prefix.myFn)()" removes the this context, so this is useful code.
				return true
			}
		}
		return expr != comma.Left().Node() && NodeUtil_IsExpressionResultUsed(parent)
	case rhino.Token_FOR:
		// Only an expression whose result is in the condition part of the
		// expression is used.
		return ast.RoleIn(expr, parent) == ast.RoleForCond
	default:
	}
	return true
}

// NodeUtil_IsExecutedExactlyOnce: n is the expression to check. Returns
// whether the expression is unconditionally executed only once in the
// containing execution scope.
//
// Java: NodeUtil.isExecutedExactlyOnce(Node)
func NodeUtil_IsExecutedExactlyOnce(n rhino.Node) bool {
	// Go-only: language-adapter — Java's do-while; `continue inspect` jumps
	// to the loop condition, which ends each iteration below.
inspect:
	for {
		parent := n.GetParent()
		switch parent.GetToken() {
		case rhino.Token_IF:
			if ast.MustIf(parent).Cond().Node() != n {
				return false
			}
		case rhino.Token_HOOK:
			if ast.MustHook(parent).Cond().Node() != n {
				return false
			}
		case rhino.Token_AND, rhino.Token_OR, rhino.Token_COALESCE:
			if ast.MustBinary(parent).Left().Node() != n {
				return false
			}
		case rhino.Token_ASSIGN_OR, rhino.Token_ASSIGN_AND, rhino.Token_ASSIGN_COALESCE:
			if ast.MustAssignOp(parent).Target().Node() != n {
				return false
			}
		case rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF:
			if ast.MustEnhancedFor(parent).Collection().Node() != n {
				return false
			}
		case rhino.Token_FOR:
			if ast.RoleIn(n, parent) != ast.RoleForInit {
				return false
			}
		case rhino.Token_TRY:
			// Consider all code under a try/catch to be conditionally executed.
			if fin, ok := ast.MustTry(parent).Finally(); !ok || fin.Node() != n {
				return false
			}
		case rhino.Token_WHILE, rhino.Token_DO, rhino.Token_CASE, rhino.Token_DEFAULT_CASE:
			return false
		case rhino.Token_SCRIPT, rhino.Token_FUNCTION:
			// Done, we've reached the scope root.
			break inspect
		default:
		}
		n = n.GetParent()
		if n.IsNil() {
			break
		}
	}
	return true
}
