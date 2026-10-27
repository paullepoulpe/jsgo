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

// NodeUtil_IsNameDeclaration: is this node a name declaration?
//
// Returns true if n is VAR, LET or CONST.
//
// Java: NodeUtil.isNameDeclaration(Node)
func NodeUtil_IsNameDeclaration(n rhino.Node) bool {
	_, ok := ast.AsNameDeclaration(n)
	return ok
}

// NodeUtil_IsDestructuringDeclaration returns true if n is a VAR, LET or
// CONST that contains a destructuring pattern.
//
// Java: NodeUtil.isDestructuringDeclaration(Node)
func NodeUtil_IsDestructuringDeclaration(n rhino.Node) bool {
	if decl, ok := ast.AsNameDeclaration(n); ok {
		for d := range decl.Declarators() {
			if _, ok := d.AsDestructuringLhs(); ok {
				return true
			}
		}
	}
	return false
}

// NodeUtil_GetAssignedValue: for an assignment or variable declaration get
// the assigned value.
//
// Returns the value node representing the new value.
//
// Java: NodeUtil.getAssignedValue(Node)
func NodeUtil_GetAssignedValue(n rhino.Node) rhino.Node {
	jutil.CheckStateMsg1(n.IsName() || n.IsGetProp(), "%s", n)
	parent := n.GetParent()
	if NodeUtil_IsNameDeclaration(parent) {
		if name, ok := ast.AsName(n); ok {
			init, _ := name.Init()
			return init.Node()
		}
		return n.GetFirstChild()
	} else if ast.RoleIn(n, parent) == ast.RoleAssignTarget {
		return ast.MustAssign(parent).Value().Node()
	} else {
		return rhino.NilNode
	}
}

// NodeUtil_IsExprAssign: is this node an assignment expression statement?
//
// Returns true if n is EXPR_RESULT and n's first child is ASSIGN.
//
// Java: NodeUtil.isExprAssign(Node)
func NodeUtil_IsExprAssign(n rhino.Node) bool {
	exprStmt, ok := ast.AsExprStmt(n)
	return ok && exprStmt.Expr().Node().IsAssign()
}

// NodeUtil_IsExprCall: is this node a call expression statement?
//
// Returns true if n is EXPR_RESULT and n's first child is CALL.
//
// Java: NodeUtil.isExprCall(Node)
func NodeUtil_IsExprCall(n rhino.Node) bool {
	exprStmt, ok := ast.AsExprStmt(n)
	return ok && exprStmt.Expr().Node().IsCall()
}

// Java: NodeUtil.isNonArrowFunction(Node)
func NodeUtil_IsNonArrowFunction(n rhino.Node) bool {
	fn, ok := ast.AsFunction(n)
	return ok && !fn.IsArrow()
}

// NodeUtil_IsNonlocalModuleExportName returns whether the given name in an
// import or export spec is not defined within the module, but is an exported
// name from this or another module.
//
// Examples include `nonlocal` in:
//
//   - export {a as nonlocal};
//   - import {nonlocal} from './foo.js';
//   - import {nonlocal as a} from './foo.js';
//   - export {nonlocal as a} from './foo.js';
//   - export {a as nonlocal} from './foo.js';
//
// n is a NAME node.
//
// Java: NodeUtil.isNonlocalModuleExportName(Node)
func NodeUtil_IsNonlocalModuleExportName(n rhino.Node) bool {
	jutil.CheckArgumentMsg1(n.IsName(), "%s", n)
	parent := n.GetParent()
	switch ast.RoleIn(n, parent) {
	case ast.RoleImportSpecImported, ast.RoleExportSpecExported:
		// import {nonlocal as x} from './foo.js'
		// export {local as nonlocal};
		return true
	case ast.RoleExportSpecLocal:
		// export {nonlocal as b} from './foo.js';
		return NodeUtil_IsExportFrom(parent.GetGrandparent())
	default:
		return false
	}
}

// NodeUtil_IsDeclaration: is the node a var, const, let, function, or class
// declaration? See NodeUtil_IsFunctionDeclaration, NodeUtil_IsClassDeclaration,
// and NodeUtil_IsNameDeclaration.
//
// Java: NodeUtil.isDeclaration(Node)
func NodeUtil_IsDeclaration(n rhino.Node) bool {
	return NodeUtil_IsNameDeclaration(n) || NodeUtil_IsFunctionDeclaration(n) || NodeUtil_IsClassDeclaration(n)
}

// NodeUtil_IsNamedExportsLiteral: whether this is an assignment to 'exports'
// that creates named exports.
//
//   - exports = {a, b}; // named export, returns true.
//   - exports = 0; // namespace export, returns false.
//   - exports = {a: 0, b}; // namespace export, returns false.
//
// Java: NodeUtil.isNamedExportsLiteral(Node)
func NodeUtil_IsNamedExportsLiteral(objectLiteral rhino.Node) bool {
	obj, ok := ast.AsObjectLit(objectLiteral)
	if !ok || !obj.HasMembers() {
		return false
	}
	for member := range obj.Members() {
		key, ok := member.AsStringKey()
		if !ok || key.IsQuoted() || !key.Value().Node().IsName() {
			return false
		}
	}
	return true
}

// NodeUtil_IsLValue determines whether this node is used as an L-value.
// Notice that sometimes names are used as both L-values and R-values.
//
// We treat "var x;" and "let x;" as an L-value because it's syntactically
// similar to "var x = undefined", even though it's technically not an
// L-value. But it kind of makes sense if you treat it as "assignment to
// 'undefined' at the top of the scope".
//
// Returns true if n is an L-value.
//
// Java: NodeUtil.isLValue(Node)
func NodeUtil_IsLValue(n rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_NAME, rhino.Token_GETPROP, rhino.Token_GETELEM:
	default:
		return false
	}

	switch ast.RoleOf(n) {
	case ast.RoleImportSpecLocal,
		ast.RoleNameDeclarationDeclarators,
		ast.RoleIterRestTarget,
		ast.RoleObjectRestTarget,
		ast.RoleParamListParams,
		ast.RoleImportDefault,
		ast.RoleIncDecTarget,
		ast.RoleCatchParam,
		ast.RoleClassName,
		ast.RoleFunctionName,
		ast.RoleDefaultValueTarget,
		ast.RoleForInit,
		ast.RoleForInTarget,
		ast.RoleForOfTarget,
		ast.RoleForAwaitOfTarget,
		ast.RoleAssignTarget,
		ast.RoleAssignOpTarget:
		return true
	case ast.RoleArrayPatternElements, ast.RoleStringKeyValue, ast.RoleComputedPropValue:
		return NodeUtil_IsLhsByDestructuring(n)
	default:
		return false
	}
}

// NodeUtil_IsDeclarationLValue determines whether this node is used as an
// L-value that is a declaration.
//
// `x = 5;` is an L-value but does not declare a variable.
//
// Java: NodeUtil.isDeclarationLValue(Node)
func NodeUtil_IsDeclarationLValue(n rhino.Node) bool {
	isLValue := NodeUtil_IsLValue(n)

	if !isLValue {
		return false
	}

	parent := n.GetParent()

	switch parent.GetToken() {
	case rhino.Token_IMPORT_SPEC, rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST, rhino.Token_PARAM_LIST,
		rhino.Token_IMPORT, rhino.Token_CATCH, rhino.Token_CLASS, rhino.Token_FUNCTION:
		return true
	case rhino.Token_STRING_KEY,
		rhino.Token_OBJECT_PATTERN,
		rhino.Token_ARRAY_PATTERN,
		rhino.Token_COMPUTED_PROP,
		rhino.Token_DEFAULT_VALUE,
		rhino.Token_ITER_REST,
		rhino.Token_OBJECT_REST:
		rootTargetParent := NodeUtil_GetRootTarget(n).GetParent()
		if rootTargetParent.IsRest() || rootTargetParent.IsDefaultValue() {
			return rootTargetParent.GetParent().IsParamList()
		} else if rootTargetParent.IsDestructuringLhs() {
			return NodeUtil_IsNameDeclaration(rootTargetParent.GetParent())
		}
		return rootTargetParent.IsParamList() || rootTargetParent.IsCatch()
	default:
		return false
	}
}

// Java: NodeUtil.isLhsOfAssign(Node)
func NodeUtil_IsLhsOfAssign(n rhino.Node) bool {
	return ast.RoleOf(n) == ast.RoleAssignTarget
}

// Java: NodeUtil.isImportedName(Node)
func NodeUtil_IsImportedName(n rhino.Node) bool {
	parent := n.GetParent()
	return parent.IsImport() || ast.RoleIn(n, parent) == ast.RoleImportSpecLocal
}

// NodeUtil_GetDeclaringParent returns the node that is effectively declaring
// the given target.
//
// Examples:
//
//	const a = 1; // getDeclaringParent(a) returns CONST
//	let {[expression]: [x = 3]} = obj; // getDeclaringParent(x) returns LET
//	function foo({a, b}) {}; // getDeclaringParent(a) returns PARAM_LIST
//	function foo(a = 1) {}; // getDeclaringParent(a) returns PARAM_LIST
//	function foo({a, b} = obj) {}; // getDeclaringParent(a) returns PARAM_LIST
//	function foo(...a) {}; // getDeclaringParent(a) returns PARAM_LIST
//	function foo() {}; // gotRootTarget(foo) returns FUNCTION
//	class foo {}; // gotRootTarget(foo) returns CLASS
//	import foo from './foo'; // getDeclaringParent(foo) returns IMPORT
//	import {foo} from './foo'; // getDeclaringParent(foo) returns IMPORT
//	import {foo as bar} from './foo'; // getDeclaringParent(bar) returns IMPORT
//	} catch (err) { // getDeclaringParent(err) returns CATCH
//
// targetNode is a NAME, OBJECT_PATTERN, or ARRAY_PATTERN; returns a node of
// type LET, CONST, VAR, FUNCTION, CLASS, PARAM_LIST, CATCH, or IMPORT. Throws
// IllegalStateException if targetNode is not actually used as a declaration
// target.
//
// Java: NodeUtil.getDeclaringParent(Node)
func NodeUtil_GetDeclaringParent(targetNode rhino.Node) rhino.Node {
	rootTarget := NodeUtil_GetRootTarget(targetNode)
	parent := rootTarget.GetParent()
	if parent.IsRest() || parent.IsDefaultValue() {
		// e.g. `function foo(targetNode1 = default, ...targetNode2) {}`
		parent = parent.GetParent()
		jutil.CheckStateMsg1(parent.IsParamList(), "%s", parent)
	} else if parent.IsDestructuringLhs() {
		// e.g. `let [a, b] = something;` targetNode is `[a, b]`
		parent = parent.GetParent()
		jutil.CheckStateMsg1(NodeUtil_IsNameDeclaration(parent), "%s", parent)
	} else if parent.IsClass() || parent.IsFunction() {
		// e.g. `function targetNode() {}`
		// e.g. `class targetNode {}`
		jutil.CheckStateMsg1(targetNode == parent.GetFirstChild(), "%s", targetNode)
	} else if parent.IsImportSpec() {
		// e.g. `import {foo as targetNode} from './foo';
		jutil.CheckStateMsg1(targetNode == parent.GetSecondChild(), "%s", targetNode)
		// import -> import_specs -> import_spec
		// we want import
		parent = parent.GetGrandparent()
		jutil.CheckStateMsg1(parent.IsImport(), "%s", parent)
	} else {
		// e.g. `function foo(targetNode) {};`
		// e.g. `let targetNode = something;`
		// e.g. `import targetNode from './foo';
		// e.g. `} catch (foo) {`
		jutil.CheckStateMsg1(
			parent.IsParamList() ||
				NodeUtil_IsNameDeclaration(parent) ||
				parent.IsImport() ||
				parent.IsCatch(),
			"%s",
			parent)
	}
	return parent
}

// NodeUtil_GetRootTarget returns the outermost target enclosing the given
// assignment target.
//
// Returns targetNode itself if there is no enclosing target.
//
// Examples:
//
//	const a = 1; // getRootTarget(a) returns a
//	let {[expression]: [x = 3]} = obj; // getRootTarget(x) returns {[expression]: [x = 3]}
//	{a = 1} = obj; // getRootTarget(a) returns {a = 1}
//	{[expression]: [x = 3]} = obj; // getRootTarget(x) returns {[expression]: [x = 3]}
//	function foo({a, b}) {}; // getRootTarget(a) returns {a, b}
//	function foo(a = 1) {}; // getRootTarget(a) returns a
//	function foo({a, b} = obj) {}; // getRootTarget(a) returns a
//	function foo(...a) {}; // getRootTarget(a) returns a
//	function foo() {}; // gotRootTarget(foo) returns foo
//	class foo {}; // gotRootTarget(foo) returns foo
//	import foo from './foo'; // getRootTarget(foo) returns foo
//	import {foo} from './foo'; // getRootTarget(foo) returns foo
//	import {foo as bar} from './foo'; // getRootTarget(bar) returns bar
//
// Throws IllegalStateException if targetNode is not actually used as a
// target.
//
// Java: NodeUtil.getRootTarget(Node)
func NodeUtil_GetRootTarget(targetNode rhino.Node) rhino.Node {
	enclosingTarget := targetNode
	for nextTarget := NodeUtil_GetEnclosingTarget(enclosingTarget); !nextTarget.IsNil(); nextTarget = NodeUtil_GetEnclosingTarget(enclosingTarget) {
		enclosingTarget = nextTarget
	}
	return enclosingTarget
}

// NodeUtil_GetEnclosingTarget returns the immediately enclosing target node
// for a given target node, or null if none found.
//
// See NodeUtil_GetRootTarget for examples.
//
// Java: NodeUtil.getEnclosingTarget(Node)
func NodeUtil_GetEnclosingTarget(targetNode rhino.Node) rhino.Node {
	jutil.CheckStateMsg1(rhino.CheckNotNull(targetNode).IsValidAssignmentTarget(), "%s", targetNode)
	parent := rhino.CheckNotNullMsg1(targetNode.GetParent(), "%s", targetNode)
	targetIsFirstChild := parent.GetFirstChild() == targetNode
	if parent.IsDefaultValue() || parent.IsRest() {
		// in `([something = targetNode] = x)` targetNode isn't actually acting
		// as a target.
		jutil.CheckStateMsg1(targetIsFirstChild, "%s", parent)
		// The DEFAULT_VALUE or REST occupies the place where the assignment target it contains would
		// otherwise be in the AST, so pretend it is the target for the logic below.
		targetNode = parent
		parent = rhino.CheckNotNull(targetNode.GetParent())
		targetIsFirstChild = targetNode == parent.GetFirstChild()
	}
	switch parent.GetToken() {
	case rhino.Token_ARRAY_PATTERN:
		// e.g. ([targetNode] = something)
		return parent

	case rhino.Token_OBJECT_PATTERN:
		// e.g. ({...rest} = something);
		return parent

	case rhino.Token_COMPUTED_PROP:
		// e.g. ({[expression]: targetNode} = something)
		// e.g. ({[expression]: targetNode = default} = something)
		// make sure the effective target (targetNode or DEFAULT_VALUE containing it)
		// isn't the expression part
		jutil.CheckStateMsg1(!targetIsFirstChild, "%s", parent)
		grandparent := rhino.CheckNotNullMsg1(parent.GetParent(), "%s", parent)
		jutil.CheckStateMsg1(grandparent.IsObjectPattern(), "%s", grandparent)
		return grandparent
	case rhino.Token_STRING_KEY:
		// e.g. ({parent: targetNode} = something)
		grandparent := rhino.CheckNotNullMsg1(parent.GetParent(), "%s", parent)
		jutil.CheckStateMsg1(grandparent.IsObjectPattern(), "%s", grandparent)
		return grandparent

	case rhino.Token_PARAM_LIST, rhino.Token_LET, rhino.Token_CONST, rhino.Token_VAR:
		// e.g. `function foo(targetNode) {}`
		// or non-destructured declarations
		// e.g. `let targetNode = 3;`
		return rhino.NilNode

	case rhino.Token_FUNCTION, rhino.Token_CLASS:
		// e.g. `function targetNode() {}`
		// e.g. `class targetNode {}`
		jutil.CheckStateMsg1(targetIsFirstChild, "%s", targetNode)
		return rhino.NilNode

	case rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF:
		// e.g. `for ({length} in obj) {}` // targetNode is `{length}`
		// e.g. `for ({length} of obj) {}` // targetNode is `{length}`
		jutil.CheckStateMsg1(targetIsFirstChild, "%s", targetNode)
		return rhino.NilNode

	case rhino.Token_DESTRUCTURING_LHS:
		// destructured declarations
		// e.g. `let [a] = 3`; // targetNode is `[a]`
		jutil.CheckStateMsg1(targetIsFirstChild, "%s", targetNode)
		return rhino.NilNode

	case rhino.Token_IMPORT:
		// e.g. `import targetNode from './foo/bar';`
		return rhino.NilNode

	case rhino.Token_IMPORT_SPEC:
		// e.g. `import {bar as targetNode} from './foo/bar';`
		// e.g. `import {targetNode} from './foo/bar';` // AST will have {targetNode as targetNode}
		jutil.CheckStateMsg1(!targetIsFirstChild, "%s", parent)
		return rhino.NilNode

	case rhino.Token_CATCH:
		// e.g. `try {} catch (foo) {}`
		return rhino.NilNode

	default:
		// e.g. targetNode = something
		jutil.CheckStateMsg1(NodeUtil_IsAssignmentOp(parent) && targetIsFirstChild, "%s", parent)
		return rhino.NilNode
	}
}

// NodeUtil_IsLhsByDestructuring returns true if the node is a lhs value of a
// destructuring assignment.
//
// For example, x in `var [x] = [1];`, `var [...x] = [1];`, and
// `var {a: x} = {a: 1}` or a.b in `([a.b] = [1]);` or
// `({key: a.b} = {key: 1});`
//
// Java: NodeUtil.isLhsByDestructuring(Node)
func NodeUtil_IsLhsByDestructuring(n rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_NAME, rhino.Token_GETPROP, rhino.Token_GETELEM:
		return NodeUtil_IsLhsByDestructuringHelper(n)
	default:
		return false
	}
}

// NodeUtil_IsLhsByDestructuringHelper returns true if the given node is
// either an LHS node in a destructuring pattern or if one of its descendants
// contains an LHS node in a destructuring pattern. For example, in
// `var {a: b = 3}}`, this returns true given the NAME b or the DEFAULT_VALUE
// node containing b.
//
// Java: NodeUtil.isLhsByDestructuringHelper(Node)
func NodeUtil_IsLhsByDestructuringHelper(n rhino.Node) bool {
	parent := n.GetParent()
	switch ast.RoleIn(n, parent) {
	// `b` in `var [b] = ...`
	// `b` in `var [...b] = ...`
	case ast.RoleArrayPatternElements, ast.RoleIterRestTarget, ast.RoleObjectRestTarget:
		return true
	// the "b" in "var {[k]: b} = ..." or "var {a: b} = ..."
	case ast.RoleComputedPropValue, ast.RoleStringKeyValue:
		return parent.GetParent().IsObjectPattern()
	case ast.RoleDefaultValueTarget:
		// The first child of a DEFAULT_VALUE is a potential LHS; the second child is the value.
		return NodeUtil_IsLhsByDestructuringHelper(parent)
	default:
		return false
	}
}

// NodeUtil_VisitLhsNodesInNode retrieves lhs nodes declared or assigned in a
// given assigning parent node.
//
// An assigning parent node is one that assigns a value to one or more LHS
// nodes.
//
// Java: NodeUtil.visitLhsNodesInNode(Node, Consumer<Node>)
func NodeUtil_VisitLhsNodesInNode(assigningParent rhino.Node, consumer func(rhino.Node)) {
	jutil.CheckArgumentMsg1(
		NodeUtil_IsNameDeclaration(assigningParent) || // e.g. `var [a,b,c];`
			assigningParent.IsParamList() || // e.g. `function foo([a,b,c]) {}`
			NodeUtil_IsAssignmentOp(assigningParent) || // e.g. `[a,b,c] = [4,5,6]`
			assigningParent.IsCatch() ||
			assigningParent.IsDestructuringLhs() ||
			assigningParent.IsDefaultValue() ||
			assigningParent.IsImport() ||
			// enhanced for loops assign to loop variables
			NodeUtil_IsEnhancedFor(assigningParent),
		"%s",
		assigningParent)
	NodeUtil_GetLhsNodesHelper(assigningParent, consumer)
}

// NodeUtil_VisitLhsNodesInDestructuringPattern retrieves lhs nodes declared
// or assigned in a given destructuring pattern node.
//
// Java: NodeUtil.visitLhsNodesInDestructuringPattern(Node, Consumer<Node>)
func NodeUtil_VisitLhsNodesInDestructuringPattern(
	destructuringPattern rhino.Node, consumer func(rhino.Node)) {
	jutil.CheckArgumentMsg(
		destructuringPattern.IsDestructuringPattern(), "Must be a destructuring pattern node.")
	NodeUtil_GetLhsNodesHelper(destructuringPattern, consumer)
}

// NodeUtil_IsExportFrom returns whether an EXPORT node has a from clause.
//
// Java: NodeUtil.isExportFrom(Node)
func NodeUtil_IsExportFrom(n rhino.Node) bool {
	jutil.CheckArgument(n.IsExport())
	_, ok := ast.MustExport(n).Source()
	return ok
}

// Java: NodeUtil.isConstantVar(Node, Scope)
func NodeUtil_IsConstantVar(node rhino.Node, scope *Scope) bool {
	if NodeUtil_IsConstantName(node) {
		return true
	}

	if !node.IsName() || scope == nil {
		return false
	}

	var_ := scope.GetVar(node.GetString())
	return var_ != nil && (var_.IsDeclaredOrInferredConst() || var_.IsConst())
}

// NodeUtil_IsConstantName determines whether a variable is constant:
//
//  1. In Normalize, any name that matches the CodingConvention#isConstant(String)
//     is annotated with an IS_CONSTANT_NAME property.
//
// node is a NAME or STRING node. Returns true if a name node represents a
// constant variable.
//
// TODO(dimvar): this method and the next two do similar but not quite
// identical things. Clean up
//
// Java: NodeUtil.isConstantName(Node)
func NodeUtil_IsConstantName(node rhino.Node) bool {
	return node.GetBooleanProp(rhino.Node_IS_CONSTANT_NAME)
}

// NodeUtil_IsConstantByConvention returns whether the given name is constant
// by coding convention.
//
// Deprecated: we want to delete the constant by convention logic - see
// http://b/135755127
//
// Java: NodeUtil.isConstantByConvention(CodingConvention, Node)
func NodeUtil_IsConstantByConvention(convention CodingConvention, node rhino.Node) bool {
	if NodeUtil_IsNormalOrOptChainGetProp(node) {
		return convention.IsConstantKey(node.GetString())
	} else if NodeUtil_MayBeObjectLitKey(node) {
		return convention.IsConstantKey(node.GetString())
	} else if node.IsName() {
		return convention.IsConstant(node.GetString())
	}
	return false
}

// NodeUtil_IsConstantDeclaration determines whether the given lvalue is
// declared constant or is a name assigned exactly once.
//
// Note that this intentionally excludes variables that are constant
// according to the coding convention.
//
// node is some lvalue node. Throws IllegalStateException if the given node
// is not an lvalue.
//
// Java: NodeUtil.isConstantDeclaration(JSDocInfo, Node)
func NodeUtil_IsConstantDeclaration(info *rhino.JSDocInfo, node rhino.Node) bool {
	if NodeUtil_IsObjectLitKey(node) ||
		(node.GetParent().IsAssign() && node.IsFirstChildOf(node.GetParent())) ||
		(node.GetParent().IsExprResult() && NodeUtil_IsNormalGet(node)) ||
		node.IsMemberFieldDef() ||
		node.IsComputedFieldDef() ||
		node.IsComputedProp() {
		return info != nil && info.IsConstant()
	}
	jutil.CheckArgumentMsg1(node.IsName(), "%s", node)
	declaringParent := NodeUtil_GetDeclaringParent(node) // throws an error if `node` is not an lvalue
	if declaringParent.IsConst() {
		return true
	} else if info != nil && info.IsConstant() {
		return true
	}
	return node.IsInferredConstantVar()
}

// Java: NodeUtil.functionHasInlineJsdocs(Node)
func NodeUtil_FunctionHasInlineJsdocs(fn rhino.Node) bool {
	f, ok := ast.AsFunction(fn)
	if !ok {
		return false
	}
	// Check inline return annotation
	if f.NameNode().Node().GetJSDocInfo() != nil {
		return true
	}
	// Check inline parameter annotations
	for param := range f.Params().Params() {
		if param.Node().GetJSDocInfo() != nil {
			return true
		}
	}
	return false
}

// NodeUtil_GetBestLValue finds the l-value that the given r-value is being
// assigned to.
//
// Java: NodeUtil.getBestLValue(Node)
func NodeUtil_GetBestLValue(n rhino.Node) rhino.Node {
	parent := n.GetParent()
	if NodeUtil_IsFunctionDeclaration(n) {
		return ast.MustFunction(n).NameNode().Node()
	} else if NodeUtil_IsClassDeclaration(n) {
		return ast.MustClass(n).NameNode()
	} else if n.IsClassMembers() {
		return NodeUtil_GetBestLValue(parent)
	} else if parent.IsName() {
		return parent
	} else if assign, ok := ast.AsAssign(parent); ok {
		return assign.Target().Node()
	} else if NodeUtil_MayBeObjectLitKey(parent) || parent.IsComputedProp() {
		return parent
	} else if role := ast.RoleIn(n, parent); role == ast.RoleHookThen ||
		role == ast.RoleHookElse ||
		parent.IsOr() ||
		parent.IsAnd() ||
		(parent.IsComma() && role == ast.RoleBinaryRight) ||
		role == ast.RoleCastExpr {
		return NodeUtil_GetBestLValue(parent)
	}
	return rhino.NilNode
}

// NodeUtil_GetRValueOfLValue gets the r-value (or initializer) of a node
// returned by getBestLValue.
//
// Java: NodeUtil.getRValueOfLValue(Node)
func NodeUtil_GetRValueOfLValue(n rhino.Node) rhino.Node {
	parent := n.GetParent()
	switch ast.RoleIn(n, parent) {
	case ast.RoleAssignTarget:
		return ast.MustAssign(parent).Value().Node()
	case ast.RoleAssignOpTarget:
		return ast.MustAssignOp(parent).Value().Node()
	case ast.RoleDestructuringLhsPattern:
		if init, ok := ast.MustDestructuringLhs(parent).Init(); ok {
			return init.Node()
		}
		return rhino.NilNode
	case ast.RoleAssignValue, ast.RoleAssignOpValue, ast.RoleDestructuringLhsInit:
		return rhino.NilNode
	}
	switch parent.GetToken() {
	case rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST:
		return n.GetLastChild()
	case rhino.Token_OBJECTLIT:
		jutil.CheckStateMsg1(
			n.IsStringKey() ||
				n.IsComputedProp() ||
				n.IsMemberFunctionDef() ||
				n.IsGetterDef() ||
				n.IsSetterDef(),
			"%s",
			n)
		return n.GetLastChild()
	case rhino.Token_CLASS_MEMBERS:
		jutil.CheckStateMsg1(
			n.IsMemberFunctionDef() ||
				n.IsMemberFieldDef() ||
				n.IsComputedFieldDef() ||
				n.IsGetterDef() ||
				n.IsSetterDef(),
			"%s",
			n)
		return n.GetLastChild()
	case rhino.Token_FUNCTION, rhino.Token_CLASS:
		return parent
	default:
	}
	return rhino.NilNode
}

// NodeUtil_GetBestLValueOwner gets the owner of the given l-value node.
//
// Java: NodeUtil.getBestLValueOwner(Node)
//
// TODO(b/189993301): do I have to add logic for class fields in these LValue
// functions?
func NodeUtil_GetBestLValueOwner(lValue rhino.Node) rhino.Node {
	if lValue.IsNil() || lValue.GetParent().IsNil() {
		return rhino.NilNode
	}
	if NodeUtil_MayBeObjectLitKey(lValue) || lValue.IsComputedProp() {
		return NodeUtil_GetBestLValue(lValue.GetParent())
	} else if NodeUtil_IsNormalGet(lValue) {
		return lValue.GetFirstChild()
	}

	return rhino.NilNode
}

// NodeUtil_GetBestLValueName gets the name of the given l-value node.
//
// Go-only: performance — the body is in nodeUtil_getBestLValueNameString, so
// that this wrapper inlines and a caller that only reads the result keeps
// the string header on its stack.
//
// Java: NodeUtil.getBestLValueName(Node)
func NodeUtil_GetBestLValueName(lValue rhino.Node) *string {
	s, ok := nodeUtil_getBestLValueNameString(lValue)
	if !ok {
		return nil
	}
	return &s
}

// nodeUtil_getBestLValueNameString is the body of getBestLValueName; ok is
// false where Java returns null.
func nodeUtil_getBestLValueNameString(lValue rhino.Node) (string, bool) {
	if lValue.IsNil() || lValue.GetParent().IsNil() {
		return "", false
	}
	if lValue.GetParent().IsClassMembers() && !lValue.IsComputedProp() {
		className := NodeUtil_GetName(lValue.GetGrandparent())
		if className == nil { // Anonymous class
			return "", false
		}
		methodName := lValue.GetString()
		if lValue.IsStaticMember() {
			return *className + "." + methodName, true
		}
		return *className + ".prototype." + methodName, true
	}
	// TODO(sdh): Tighten this to simply require !lValue.isQuotedString()
	// Could get rid of the isJSIdentifier check, but may need to fix depot.
	if NodeUtil_MayBeObjectLitKey(lValue) {
		owner := NodeUtil_GetBestLValue(lValue.GetParent())
		if !owner.IsNil() {
			key := NodeUtil_GetObjectOrClassLitKeyName(lValue)
			if !rhino.TokenStream_IsJSIdentifier(key) {
				return "", false
			}
			ownerName, ok := nodeUtil_getBestLValueNameString(owner)
			if ok {
				return ownerName + "." + key, true
			}
		}
		return "", false
	}
	return lValue.GetQualifiedNameString()
}

// NodeUtil_GetBestLValueRoot gets the root of a qualified name l-value.
//
// Java: NodeUtil.getBestLValueRoot(Node)
func NodeUtil_GetBestLValueRoot(lValue rhino.Node) rhino.Node {
	if lValue.IsNil() {
		return rhino.NilNode
	}
	switch lValue.GetToken() {
	case rhino.Token_STRING_KEY:
		// NOTE: beware of getBestLValue returning null (or be null-permissive?)
		return NodeUtil_GetBestLValueRoot(NodeUtil_GetBestLValue(lValue.GetParent()))
	case rhino.Token_GETPROP, rhino.Token_GETELEM:
		return NodeUtil_GetBestLValueRoot(lValue.GetFirstChild())
	case rhino.Token_THIS, rhino.Token_SUPER, rhino.Token_NAME:
		return lValue
	default:
		return rhino.NilNode
	}
}
