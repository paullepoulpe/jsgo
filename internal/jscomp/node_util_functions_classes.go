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
	"iter"
	"slices"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

func NodeUtil_GetNameString(n rhino.Node) (string, bool) {
	nameNode := NodeUtil_GetNameNode(n)
	if nameNode.IsNil() {
		return "", false
	}
	return nameNode.GetQualifiedNameString()
}

// NodeUtil_GetNameNode gets the node of a function or class's name. This
// method recognizes five forms:
//
//   - class name {...}
//   - var name = class {...}
//   - qualified.name = class {...}
//   - var name2 = class name1 {...}
//   - qualified.name2 = class name1 {...}
//
// In two last cases with named function expressions, the second name is
// returned (the variable or qualified name).
//
// n is a function or class node; returns the node best representing the
// class's name.
//
// Java: NodeUtil.getNameNode(Node)
func NodeUtil_GetNameNode(n rhino.Node) rhino.Node {
	jutil.CheckStateMsg1(n.IsFunction() || n.IsClass(), "%s", n)
	parent := n.GetParent()
	switch parent.GetToken() {
	case rhino.Token_NAME:
		// var name = function() ...
		// var name2 = function name1() ...
		return parent
	case rhino.Token_ASSIGN:
		// qualified.name = function() ...
		// qualified.name2 = function name1() ...
		target := ast.MustAssign(parent).Target().Node()
		if target.IsQualifiedName() {
			return target
		}
		return rhino.NilNode
	default:
		// function name() ...
		//   or
		// class Name ...
		if fn, ok := ast.AsFunction(n); ok {
			if name, ok := fn.Name(); ok {
				return name.Node()
			}
			return rhino.NilNode
		}
		if name, ok := ast.MustClass(n).Name(); ok && name.Ident() != "" {
			return name.Node()
		}
		return rhino.NilNode
	}
}

// NodeUtil_RemoveName sets the given function/class node to an empty name.
//
// Java: NodeUtil.removeName(Node)
func NodeUtil_RemoveName(n rhino.Node) {
	jutil.CheckState(n.IsFunction() || n.IsClass())
	originalName := n.GetFirstChild()
	var emptyName rhino.Node
	if n.IsFunction() {
		emptyName = rhino.IR_Name("")
	} else {
		emptyName = rhino.IR_Empty()
	}
	originalName.ReplaceWith(emptyName.Srcref(originalName))
}

// NodeUtil_GetNearestFunctionName gets the function's name. This method
// recognizes the forms:
//
//   - {'name': function() ...}
//   - {name: function() ...}
//   - function name() ...
//   - var name = function() ...
//   - var obj = {name() {} ...}
//   - qualified.name = function() ...
//   - var name2 = function name1() ...
//   - qualified.name2 = function name1() ...
//
// n is a node whose type is Token#FUNCTION; returns the function's name, or
// null if it has no name.
//
// Java: NodeUtil.getNearestFunctionName(Node)
func NodeUtil_GetNearestFunctionName(n rhino.Node) *string {
	if !n.IsFunction() {
		return nil
	}

	name := NodeUtil_GetName(n)
	if name != nil {
		return name
	}

	// Check for the form { 'x' : function() { }} and {x() {}}
	parent := n.GetParent()
	switch parent.GetToken() {
	// Return the name of the literal's key.
	case rhino.Token_MEMBER_FUNCTION_DEF, rhino.Token_SETTER_DEF, rhino.Token_GETTER_DEF, rhino.Token_STRING_KEY:
		return jutil.StringPtr(parent.GetString())
	case rhino.Token_NUMBER:
		return NodeUtil_GetStringValue(parent)
	default:
		return nil
	}
}

// Java: NodeUtil.getClassMembers(Node)
func NodeUtil_GetClassMembers(n rhino.Node) rhino.Node {
	jutil.CheckArgument(n.IsClass())
	return ast.MustClass(n).Body().Node()
}

// Java: NodeUtil.getEs6ClassConstructorMemberFunctionDef(Node)
func NodeUtil_GetEs6ClassConstructorMemberFunctionDef(classNode rhino.Node) rhino.Node {
	jutil.CheckArgumentMsg1(classNode.IsClass(), "%s", classNode)
	cls := ast.MustClass(classNode)
	rhino.CheckNotNullMsg1(cls.Body().Node(), "%s", classNode)
	if ctor, ok := cls.Constructor(); ok {
		return ctor.Node()
	}
	return rhino.NilNode
}

// NodeUtil_GetFirstGetterMatchingKey returns the first getter in the class
// members that matches the key.
//
// Java: NodeUtil.getFirstGetterMatchingKey(Node, String)
func NodeUtil_GetFirstGetterMatchingKey(n rhino.Node, keyName string) rhino.Node {
	jutil.CheckStateMsg1(n.IsClassMembers() || n.IsObjectLit(), "%s", n)
	for keyNode := n.GetFirstChild(); !keyNode.IsNil(); keyNode = keyNode.GetNext() {
		if keyNode.IsGetterDef() && keyNode.GetString() == keyName {
			return keyNode
		}
	}
	return rhino.NilNode
}

// NodeUtil_ReferencesOwnReceiver returns true if this function references
// its receiver object.
//
// We define this function in terms of "receiver" rather than specific
// syntax, however internally we search for `this` and `super`.
//
// Arrow functions return false. They do not have their own receiver, but
// rather capture the receiver of the enclosing scope.
//
// Java: NodeUtil.referencesOwnReceiver(Node)
func NodeUtil_ReferencesOwnReceiver(fn rhino.Node) bool {
	jutil.CheckState(fn.IsFunction())
	if fn.IsArrowFunction() {
		return false
	}

	return NodeUtil_ReferencesEnclosingReceiver(NodeUtil_GetFunctionParameters(fn)) ||
		NodeUtil_ReferencesEnclosingReceiver(NodeUtil_GetFunctionBody(fn))
}

// NodeUtil_ReferencesEnclosingReceiver returns true if this subtree
// references the receiver object from its enclosing function scope.
//
// We define this function in terms of "receiver" rather than specific
// syntax, however internally we search for `this` and `super`.
//
// Arrow functions may return true. They capture the receiver of the
// enclosing scope, rather than having their own.
//
// Java: NodeUtil.referencesEnclosingReceiver(Node)
func NodeUtil_ReferencesEnclosingReceiver(n rhino.Node) bool {
	return NodeUtil_Has(n, func(c rhino.Node) bool { return c.IsThis() || c.IsSuper() },
		NodeUtil_MATCH_ANYTHING_BUT_NON_ARROW_FUNCTION)
}

// NodeUtil_ReferencesSuper returns true if the current scope contains
// references to the 'super' keyword. Note that if there are classes declared
// inside the current class, super calls which reference those classes are
// not reported.
//
// Java: NodeUtil.referencesSuper(Node)
func NodeUtil_ReferencesSuper(n rhino.Node) bool {
	curr := n.GetFirstChild()
	for !curr.IsNil() {
		if NodeUtil_Has(curr, (rhino.Node).IsSuper, func(node rhino.Node) bool { return !node.IsClass() }) {
			return true
		}
		curr = curr.GetNext()
	}
	return false
}

// NodeUtil_IsNormalOrOptChainGet: is this a GETPROP, OPTCHAIN_GETPROP,
// GETELEM, or OPTCHAIN_GETELEM?
//
// Java: NodeUtil.isNormalOrOptChainGet(Node)
func NodeUtil_IsNormalOrOptChainGet(n rhino.Node) bool {
	return NodeUtil_IsNormalGet(n) || NodeUtil_IsOptChainGet(n)
}

// NodeUtil_IsNormalOrOptChainGetProp: is this a GETPROP or OPTCHAIN_GETPROP?
//
// Java: NodeUtil.isNormalOrOptChainGetProp(Node)
func NodeUtil_IsNormalOrOptChainGetProp(n rhino.Node) bool {
	return n.IsGetProp() || n.IsOptChainGetProp()
}

// NodeUtil_IsNormalOrOptChainCall: is this a CALL or OPTCHAIN_CALL?
//
// Java: NodeUtil.isNormalOrOptChainCall(Node)
func NodeUtil_IsNormalOrOptChainCall(n rhino.Node) bool {
	return n.IsCall() || n.IsOptChainCall()
}

// NodeUtil_IsNormalGet: is this a GETPROP or GETELEM node?
//
// Java: NodeUtil.isNormalGet(Node)
func NodeUtil_IsNormalGet(n rhino.Node) bool {
	return n.IsGetProp() || n.IsGetElem()
}

// NodeUtil_IsOptChainGet: is this an OPTCHAIN_GETPROP or OPTCHAIN_GETELEM
// node?
//
// Java: NodeUtil.isOptChainGet(Node)
func NodeUtil_IsOptChainGet(n rhino.Node) bool {
	return n.IsOptChainGetProp() || n.IsOptChainGetElem()
}

// NodeUtil_IsOptChainNode: is this a OPTCHAIN_GETPROP, OPTCHAIN_GETELEM,
// OPTCHAIN_CALL node?
//
// Java: NodeUtil.isOptChainNode(Node)
func NodeUtil_IsOptChainNode(n rhino.Node) bool {
	return n.IsOptChainGetProp() || n.IsOptChainGetElem() || n.IsOptChainCall()
}

// NodeUtil_GetStartOfOptChainSegment finds the start of the optional chain.
// E.g Find the `a?. ...` node in `a?.b.c.d` given any other node.
//
// n is a node in an optional chain; returns the start of the optional chain
// that n is part of.
//
// Java: NodeUtil.getStartOfOptChainSegment(Node)
func NodeUtil_GetStartOfOptChainSegment(n rhino.Node) rhino.Node {
	jutil.CheckStateMsg1(NodeUtil_IsOptChainNode(n), "%s", n)
	if n.IsOptionalChainStart() {
		return n
	}
	return NodeUtil_GetStartOfOptChainSegment(n.GetFirstChild())
}

// NodeUtil_GetEndOfOptChainSegment finds the end of an optional chain
// segment.
//
// Each `?.` ends one segment and starts another.
//
// Examples
//
//	a?.b.c.d   // end is the node with children `a?.b.c` and `d`
//	a?.b.c?.d  // given a?.b end is the node with children `a?.b` and `c`
//	a?.b.c?.d  // given a?.b.c end is the node with children `a?.b.c` and `d`
//	(a?.b.c).d // given a?.b end is the node with children `a?.b` and `c`
//
// n is a node in an optional chain; returns the end of the optional chain
// that n is part of.
//
// Java: NodeUtil.getEndOfOptChainSegment(Node)
func NodeUtil_GetEndOfOptChainSegment(n rhino.Node) rhino.Node {
	jutil.CheckStateMsg1(NodeUtil_IsOptChainNode(n), "%s", n)
	if NodeUtil_IsEndOfOptChainSegment(n) {
		return n
	} else {
		return NodeUtil_GetEndOfOptChainSegment(n.GetParent())
	}
}

// NodeUtil_IsEndOfFullOptChain: is this node the final node of a full
// optional chain?
//
// e.g. for `a?.b.c?.d` this method returns true only for the Node with
// children `a?.b.c` and `d`. That node is the end of the whole chain, and
// also represents the whole chain in the AST.
//
// Java: NodeUtil.isEndOfFullOptChain(Node)
func NodeUtil_IsEndOfFullOptChain(n rhino.Node) bool {
	if NodeUtil_IsOptChainNode(n) {
		parent := n.GetParent()
		// the chain continues if this is the first child of another optional chain node
		return !(NodeUtil_IsOptChainNode(parent) && n.IsFirstChildOf(parent))
	} else {
		// not even an optional chain node
		return false
	}
}

// NodeUtil_IsEndOfOptChainSegment: is this node the end of an optional chain
// segment?
//
// Each `?.` begins a new segment and ends the previous one, if any. The end
// of the whole chain is also the end of its final segment.
//
// Java: NodeUtil.isEndOfOptChainSegment(Node)
func NodeUtil_IsEndOfOptChainSegment(n rhino.Node) bool {
	if !NodeUtil_IsOptChainNode(n) {
		return false
	} else {
		parent := n.GetParent()
		// Check for null so this method will work for a disconnected Node.
		if !parent.IsNil() && n.IsFirstChildOf(parent) && NodeUtil_IsOptChainNode(parent) {
			// The parent is a continuation of this optional chain.
			// If it starts a new segment, then this node is the end of a segment
			return parent.IsOptionalChainStart()
		} else {
			// The parent doesn't continue this node's optional chain, though it might be part of a
			// different one.
			// e.g. in `a?.(x?.y.z)` the parent of `x?.y.z` is part of a different optional chain.
			return true
		}
	}
}

// NodeUtil_ConvertToNonOptionalChainSegment: given the end of an optional
// chain segment changes all nodes from the end down to the start into
// non-optional nodes. e.g `({a})?.a.b.c.d()?.x.y.z` gets converted to
// `({a}).a.b.c.d()?.x.y.z` when the passed in endOfOptChainSegment is
// `({a})?.a.b.c.d()`.
//
// Java: NodeUtil.convertToNonOptionalChainSegment(Node)
func NodeUtil_ConvertToNonOptionalChainSegment(endOfOptChainSegment rhino.Node) {
	start := NodeUtil_GetStartOfOptChainSegment(endOfOptChainSegment)
	NodeUtil_ConvertToNonOptionalChainSegmentDownTo(endOfOptChainSegment, start.GetFirstChild())
}

// NodeUtil_ConvertToNonOptionalChainSegmentDownTo: given the end of an
// optional chain segment, changes all nodes from endOfOptChainSegment down
// to (but excluding) stopNode into non-optional nodes.
//
// Java: NodeUtil.convertToNonOptionalChainSegmentDownTo(Node, Node)
func NodeUtil_ConvertToNonOptionalChainSegmentDownTo(endOfOptChainSegment rhino.Node, stopNode rhino.Node) {
	jutil.CheckArgumentMsg1(NodeUtil_IsEndOfOptChainSegment(endOfOptChainSegment), "%s", endOfOptChainSegment)
	var segmentNodes []rhino.Node
	segmentNode := endOfOptChainSegment
	for !segmentNode.IsNil() && segmentNode != stopNode {
		if NodeUtil_IsOptChainNode(segmentNode) {
			segmentNodes = append(segmentNodes, segmentNode)
		}
		segmentNode = segmentNode.GetFirstChild()
	}
	for _, n := range segmentNodes {
		n.SetIsOptionalChainStart(false)
		n.SetToken(NodeUtil_GetNonOptChainToken(n.GetToken()))
	}
}

// Java: NodeUtil.getNonOptChainToken(Token)
func NodeUtil_GetNonOptChainToken(optChainToken rhino.Token) rhino.Token {
	switch optChainToken {
	case rhino.Token_OPTCHAIN_CALL:
		return rhino.Token_CALL
	case rhino.Token_OPTCHAIN_GETELEM:
		return rhino.Token_GETELEM
	case rhino.Token_OPTCHAIN_GETPROP:
		return rhino.Token_GETPROP
	default:
		panic(jutil.NewIllegalStateExceptionMsg("Should be an OPTCHAIN token: " + optChainToken.String()))
	}
}

// NodeUtil_IsBlockScopedDeclaration: is this node a block-scoped
// declaration? Checks for let, const, class, or block-scoped function
// declarations.
//
// Returns true if n is a block-scoped declaration statement.
//
// Java: NodeUtil.isBlockScopedDeclaration(Node)
func NodeUtil_IsBlockScopedDeclaration(n rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_LET, rhino.Token_CONST, rhino.Token_CATCH:
		return true
	case rhino.Token_CLASS:
		return NodeUtil_IsClassDeclaration(n)
	case rhino.Token_FUNCTION:
		return NodeUtil_IsBlockScopedFunctionDeclaration(n)
	default:
		return false
	}
}

// Java: NodeUtil.isEnhancedFor(Node)
func NodeUtil_IsEnhancedFor(n rhino.Node) bool {
	return n.IsForOf() || n.IsForAwaitOf() || n.IsForIn()
}

// NodeUtil_IsCallOrNew returns whether the call is a NEW or CALL node.
//
// Java: NodeUtil.isCallOrNew(Node)
func NodeUtil_IsCallOrNew(node rhino.Node) bool {
	return node.IsCall() || node.IsNew() || node.IsOptChainCall()
}

// NodeUtil_GetFunctionBody returns a BLOCK node for the given FUNCTION node.
//
// Java: NodeUtil.getFunctionBody(Node)
func NodeUtil_GetFunctionBody(fn rhino.Node) rhino.Node {
	jutil.CheckArgumentMsg1(fn.IsFunction(), "%s", fn)
	return ast.MustFunction(fn).Body().Node()
}

// NodeUtil_GetCallTargetResolvingIndirectCalls returns the call target for a
// call expression, resolving an indirected call (`(0, foo)()`) if present.
//
// Java: NodeUtil.getCallTargetResolvingIndirectCalls(Node)
func NodeUtil_GetCallTargetResolvingIndirectCalls(call rhino.Node) rhino.Node {
	jutil.CheckArgumentMsg1(
		call.IsCall() || call.IsNew() || call.IsTaggedTemplateLit(),
		"must be call, new or tagged-template-literal expression, got %s",
		call)
	target := ast.MustInvocation(call).Callee().Node()
	if target.IsComma() {
		right := ast.MustBinary(target).Right().Node()
		if right.IsQualifiedName() {
			return right
		}
	}
	return target
}

// NodeUtil_IsFunctionDeclaration: is this node a function declaration? A
// function declaration is a function that has a name that is added to the
// current scope (i.e. a function that is not part of a expression; see
// NodeUtil_IsFunctionExpression).
//
// Java: NodeUtil.isFunctionDeclaration(Node)
func NodeUtil_IsFunctionDeclaration(n rhino.Node) bool {
	// Note: There is currently one case where an unnamed function has a declaration parent.
	// `export default function() {...}`
	// In this case we consider the function to be an expression.
	return n.IsFunction() && NodeUtil_IsDeclarationParent(n.GetParent()) && NodeUtil_IsNamedFunction(n)
}

// NodeUtil_IsMethodDeclaration: is this node a class or object literal
// member function?
//
// examples:
//
//	class C {
//	  f() {}
//	  get x() { return this.x_; }
//	  set x(v) { this.x_ = v; }
//	  [someExpr]() {}
//	}
//	obj = {
//	  f() {}
//	  get x() { return this.x_; }
//	  set x(v) { this.x_ = v; }
//	  [someExpr]() {}
//	}
//
// Java: NodeUtil.isMethodDeclaration(Node)
func NodeUtil_IsMethodDeclaration(n rhino.Node) bool {
	if !n.IsFunction() {
		return false
	}
	parent := n.GetParent()
	switch ast.RoleIn(n, parent) {
	case ast.RoleGetterDefFunction, ast.RoleSetterDefFunction, ast.RoleMemberFunctionDefFunction:
		return true
	case ast.RoleComputedPropValue:
		comp := ast.MustComputedProp(parent)
		return comp.IsMethod() || comp.IsGetter() || comp.IsSetter()
	default:
		return false
	}
}

// NodeUtil_IsClassDeclaration: is this a class declaration.
//
// Java: NodeUtil.isClassDeclaration(Node)
func NodeUtil_IsClassDeclaration(n rhino.Node) bool {
	return n.IsClass() && NodeUtil_IsDeclarationParent(n.GetParent()) && NodeUtil_IsNamedClass(n)
}

// NodeUtil_IsHoistedFunctionDeclaration: is this node a hoisted function
// declaration? A function declaration in the scope root is hoisted to the
// top of the scope. See NodeUtil_IsFunctionDeclaration.
//
// Java: NodeUtil.isHoistedFunctionDeclaration(Node)
func NodeUtil_IsHoistedFunctionDeclaration(n rhino.Node) bool {
	if NodeUtil_IsFunctionDeclaration(n) {
		parent := n.GetParent()
		return parent.IsScript() ||
			parent.IsModuleBody() ||
			parent.GetParent().IsFunction() ||
			parent.IsExport()
	}
	return false
}

// Java: NodeUtil.isBlockScopedFunctionDeclaration(Node)
func NodeUtil_IsBlockScopedFunctionDeclaration(n rhino.Node) bool {
	if !NodeUtil_IsFunctionDeclaration(n) {
		return false
	}
	current := n.GetParent()
	for !current.IsNil() {
		switch current.GetToken() {
		case rhino.Token_BLOCK:
			return !current.GetParent().IsFunction()
		case rhino.Token_FUNCTION, rhino.Token_SCRIPT, rhino.Token_DECLARE, rhino.Token_EXPORT,
			rhino.Token_MODULE_BODY:
			return false
		default:
			jutil.CheckStateMsg1(current.IsLabel(), "%s", current)
			current = current.GetParent()
		}
	}
	return false
}

// Java: NodeUtil.isFunctionBlock(Node)
func NodeUtil_IsFunctionBlock(n rhino.Node) bool {
	return n.IsBlock() && n.HasParent() && n.GetParent().IsFunction()
}

// Java: NodeUtil.isClassStaticBlock(Node)
func NodeUtil_IsClassStaticBlock(n rhino.Node) bool {
	return n.IsBlock() && n.HasParent() && n.GetParent().IsClassMembers()
}

// NodeUtil_IsNonStaticClassFieldInitializer returns whether n is the
// initializer expression of a non-static class field (e.g. a in
// `class { x = a; }` or `class { [k] = a; }`).
//
// Unlike static field initializers and computed field keys, which are
// evaluated eagerly when the class definition is evaluated, non-static field
// initializers are evaluated once per new call (like an inner function
// body).
//
// Java: NodeUtil.isNonStaticClassFieldInitializer(Node)
func NodeUtil_IsNonStaticClassFieldInitializer(n rhino.Node) bool {
	parent := n.GetParent()
	switch ast.RoleIn(n, parent) {
	case ast.RoleMemberFieldDefValue, ast.RoleComputedFieldDefValue:
		return !ast.MustFieldDef(parent).IsStatic()
	default:
		return false
	}
}

// NodeUtil_IsFunctionOrNonStaticClassFieldInitializer returns whether n is a
// function or a non-static class field initializer: both are code whose
// evaluation is deferred until later (and possibly repeated), and which may
// capture variables from the enclosing scope.
//
// Java: NodeUtil.isFunctionOrNonStaticClassFieldInitializer(Node)
func NodeUtil_IsFunctionOrNonStaticClassFieldInitializer(n rhino.Node) bool {
	return n.IsFunction() || NodeUtil_IsNonStaticClassFieldInitializer(n)
}

// NodeUtil_IsFunctionExpression: is a FUNCTION node a function expression?
//
// A function expression is a function that:
//
//   - has either no name or a name that is not added to the current scope
//   - AND can be manipulated as an expression (assigned to variables, passed
//     to functions, etc.) i.e. It is not a method declaration on a class or
//     object literal.
//
// Some examples of function expressions:
//
//	(function () {})
//	(function f() {})()
//	[ function f() {} ]
//	var f = function f() {};
//	for (function f() {};;) {}
//	export default function() {}
//	() => 1
//
// Some examples of functions that are not expressions:
//
//	function f() {}
//	if (x); else function f() {}
//	for (;;) { function f() {} }
//	export default function f() {}
//	({
//	  f() {},
//	  set x(v) {},
//	  get x() {},
//	  [expr]() {}
//	})
//	class {
//	  f() {}
//	  set x(v) {}
//	  get x() {}
//	  [expr]() {}
//	}
//
// Returns whether n is a function used within an expression.
//
// Java: NodeUtil.isFunctionExpression(Node)
func NodeUtil_IsFunctionExpression(n rhino.Node) bool {
	return n.IsFunction() && !NodeUtil_IsFunctionDeclaration(n) && !NodeUtil_IsMethodDeclaration(n)
}

// NodeUtil_IsNamedFunctionExpression returns whether the node is both a
// function expression and the function is named.
//
// Java: NodeUtil.isNamedFunctionExpression(Node)
func NodeUtil_IsNamedFunctionExpression(n rhino.Node) bool {
	if !NodeUtil_IsFunctionExpression(n) {
		return false
	}
	_, hasName := ast.MustFunction(n).Name()
	return hasName
}

// NodeUtil_IsClassExpression: see NodeUtil_IsFunctionExpression. Returns
// whether n is a class used within an expression.
//
// Java: NodeUtil.isClassExpression(Node)
func NodeUtil_IsClassExpression(n rhino.Node) bool {
	return n.IsClass() && (!NodeUtil_IsNamedClass(n) || !NodeUtil_IsDeclarationParent(n.GetParent()))
}

// NodeUtil_IsNamedClassExpression returns whether the node is both a class
// expression and the class is named.
//
// Java: NodeUtil.isNamedClassExpression(Node)
func NodeUtil_IsNamedClassExpression(n rhino.Node) bool {
	if !NodeUtil_IsClassExpression(n) {
		return false
	}
	_, hasName := ast.MustClass(n).Name()
	return hasName
}

// NodeUtil_IsNamedFunction returns whether n is a function with a nonempty
// name. Differs from NodeUtil_IsFunctionDeclaration because the name might
// in a function expression and not be added to the current scope.
//
// Some named functions include
//
//	(function f() {})();
//	export default function f() {};
//	function f() {};
//	var f = function f() {};
//
// Java: NodeUtil.isNamedFunction(Node)
func NodeUtil_IsNamedFunction(n rhino.Node) bool {
	fn, ok := ast.AsFunction(n)
	if !ok {
		return false
	}
	_, hasName := fn.Name()
	return hasName
}

// NodeUtil_IsNamedClass: see NodeUtil_IsNamedFunction. Returns whether n is
// a named class.
//
// Java: NodeUtil.isNamedClass(Node)
func NodeUtil_IsNamedClass(n rhino.Node) bool {
	cls, ok := ast.AsClass(n)
	if !ok {
		return false
	}
	name, hasName := cls.Name()
	return hasName && name.Ident() != ""
}

// NodeUtil_IsBleedingFunctionName returns whether this is a bleeding
// function (an anonymous named function that bleeds into the inner scope).
//
// Java: NodeUtil.isBleedingFunctionName(Node)
func NodeUtil_IsBleedingFunctionName(n rhino.Node) bool {
	if !n.IsName() || n.GetString() == "" {
		return false
	}
	parent := n.GetParent()
	return NodeUtil_IsFunctionExpression(parent) && ast.RoleIn(n, parent) == ast.RoleFunctionName
}

// NodeUtil_IsEmptyFunctionExpression determines if a node is a function
// expression that has an empty body.
//
// Returns whether the given node is a function expression that is empty.
//
// Java: NodeUtil.isEmptyFunctionExpression(Node)
func NodeUtil_IsEmptyFunctionExpression(node rhino.Node) bool {
	return NodeUtil_IsFunctionExpression(node) && NodeUtil_IsEmptyBlock(ast.MustFunction(node).Body().Node())
}

// NodeUtil_DoesFunctionReferenceOwnArgumentsObject returns whether a
// function has a reference to its own "arguments" object.
//
// Java: NodeUtil.doesFunctionReferenceOwnArgumentsObject(Node)
func NodeUtil_DoesFunctionReferenceOwnArgumentsObject(fn rhino.Node) bool {
	jutil.CheckArgument(fn.IsFunction())
	f := ast.MustFunction(fn)
	if f.IsArrow() {
		return false
	}
	return NodeUtil_ReferencesArgumentsHelper(f.Params().Node()) ||
		NodeUtil_ReferencesArgumentsHelper(f.Body().Node())
}

// NodeUtil_ReferencesArgumentsHelper returns whether any child is a
// reference to the "arguments" object of the root. Effectively, this
// includes arrow method bodies (which don't have their own) and excludes
// other functions which shadow the "arguments" value with their own.
//
// Java: NodeUtil.referencesArgumentsHelper(Node)
func NodeUtil_ReferencesArgumentsHelper(node rhino.Node) bool {
	if node.IsName() && node.GetString() == "arguments" {
		return true
	}

	if NodeUtil_IsNonArrowFunction(node) {
		return false
	}

	for c := node.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		if NodeUtil_ReferencesArgumentsHelper(c) {
			return true
		}
	}

	return false
}

// NodeUtil_IsObjectCallMethod returns whether node is a call to methodName:
//
//	a.f(...)
//	a?.f(...)
//	a['f'](...)
//
// Java: NodeUtil.isObjectCallMethod(Node, String)
func NodeUtil_IsObjectCallMethod(callNode rhino.Node, methodName string) bool {
	if call, ok := ast.AsNormalOrOptChainCall(callNode); ok {
		if get, ok := ast.AsNormalOrOptChainGet(call.Callee().Node()); ok {
			propName, ok := get.StaticPropertyName()
			return ok && propName == methodName
		}
	}
	return false
}

// NodeUtil_IsFunctionObjectCall returns whether the callNode represents an
// expression in the form of: x.call(...) x['call'](...)
//
// Java: NodeUtil.isFunctionObjectCall(Node)
func NodeUtil_IsFunctionObjectCall(callNode rhino.Node) bool {
	return NodeUtil_IsObjectCallMethod(callNode, "call")
}

// NodeUtil_IsFunctionObjectApply returns whether the callNode represents an
// expression in the form of: x.apply(...) x['apply'](...)
//
// Java: NodeUtil.isFunctionObjectApply(Node)
func NodeUtil_IsFunctionObjectApply(callNode rhino.Node) bool {
	return NodeUtil_IsObjectCallMethod(callNode, "apply")
}

// NodeUtil_IsNameDeclOrSimpleAssignLhs determines whether this node is
// strictly on the left hand side of an assign or var initialization.
// Notably, this does not include all L-values, only statements where the
// node is used only as an L-value.
//
// n is the node, parent the parent of the node. Returns true if n is the
// left hand of an assign.
//
// Java: NodeUtil.isNameDeclOrSimpleAssignLhs(Node, Node)
func NodeUtil_IsNameDeclOrSimpleAssignLhs(n rhino.Node, parent rhino.Node) bool {
	return ast.RoleIn(n, parent) == ast.RoleAssignTarget || NodeUtil_IsNameDeclaration(parent)
}

// NodeUtil_MayBeObjectLitKey determines whether a node represents a possible
// object literal key (e.g. key1 in {key1: value1, key2: value2}). Computed
// properties are excluded here (see b/111621528). This method does not check
// whether the node is actually in an object literal! it also returns true
// for object pattern keys, and member functions/getters in ES6 classes.
//
// Java: NodeUtil.mayBeObjectLitKey(Node)
//
// TODO(b/189993301): should fields be added to this method?
func NodeUtil_MayBeObjectLitKey(node rhino.Node) bool {
	switch node.GetToken() {
	case rhino.Token_STRING_KEY, rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF, rhino.Token_MEMBER_FUNCTION_DEF:
		return true
	default:
		return false
	}
}

// NodeUtil_IsObjectLitKey determines whether a node represents an object
// literal key (e.g. key1 in {key1: value1, key2: value2}) and is in an
// object literal. Computed properties are excluded here (see b/111621528).
//
// Java: NodeUtil.isObjectLitKey(Node)
func NodeUtil_IsObjectLitKey(node rhino.Node) bool {
	return node.GetParent().IsObjectLit() && NodeUtil_MayBeObjectLitKey(node)
}

// NodeUtil_GetObjectOrClassLitKeyName gets the name of an object literal
// key.
//
// Java: NodeUtil.getObjectOrClassLitKeyName(Node)
func NodeUtil_GetObjectOrClassLitKeyName(key rhino.Node) string {
	keyNode := NodeUtil_GetObjectOrClassLitKeyNode(key)
	if !keyNode.IsNil() {
		return keyNode.GetString()
	}
	panic(jutil.NewIllegalStateExceptionMsg("Unexpected node type: " + key.String()))
}

// NodeUtil_GetObjectOrClassLitKeyNode gets the Node that defines the name of
// an object literal key.
//
// Java: NodeUtil.getObjectOrClassLitKeyNode(Node)
func NodeUtil_GetObjectOrClassLitKeyNode(key rhino.Node) rhino.Node {
	switch key.GetToken() {
	case rhino.Token_STRING_KEY, rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF, rhino.Token_MEMBER_FUNCTION_DEF,
		rhino.Token_MEMBER_FIELD_DEF:
		return key
	case rhino.Token_COMPUTED_PROP, rhino.Token_COMPUTED_FIELD_DEF:
		if key.GetFirstChild().
			IsStringLit() { // TODO(b/189993301): may be an issue with non string lits
			return key.GetFirstChild()
		}
		return rhino.NilNode
	default:
	}
	panic(jutil.NewIllegalStateExceptionMsg("Unexpected node type: " + key.String()))
}

// NodeUtil_IsGetOrSetKey determines whether a node represents an object
// literal get or set key (e.g. key1 in {get key1() {}, set key2(a){}).
//
// Java: NodeUtil.isGetOrSetKey(Node)
func NodeUtil_IsGetOrSetKey(node rhino.Node) bool {
	if _, ok := ast.AsGetterOrSetterDef(node); ok {
		return true
	}
	if comp, ok := ast.AsComputedProp(node); ok {
		return comp.IsGetter() || comp.IsSetter()
	}
	return false
}

// NodeUtil_OpToStr converts an operator's token value (see Token) to a
// string representation.
//
// operator is the operator's token value to convert; returns the string
// representation or null if the token value is not an operator.
//
// Java: NodeUtil.opToStr(Token)
func NodeUtil_OpToStr(operator rhino.Token) *string {
	// Go-only: language-adapter — Java returns interned string constants; the
	// results of the switch below are computed once per token so that the
	// *string result does not allocate on every call (the code printer calls
	// this per operator).
	return nodeUtil_opToStrTable[operator]
}

// NodeUtil_IsObjectDefinePropertiesDefinition returns true if the node is a
// definition with Object.defineProperties.
//
// Java: NodeUtil.isObjectDefinePropertiesDefinition(Node)
func NodeUtil_IsObjectDefinePropertiesDefinition(n rhino.Node) bool {
	// We intentionally don't check optional Object.defineProperties?.() because we expect it to be
	// defined or polyfill-ed, and to avoid changing the current invariant that this method returns
	// true only for non-optional `Object.defineProperties()`.
	call, ok := ast.AsCall(n)
	if !ok || !n.HasXChildren(3) {
		return false
	}
	getprop, ok := ast.AsGetProp(call.Callee().Node())
	return ok &&
		getprop.Property() == "defineProperties" &&
		NodeUtil_IsKnownGlobalObjectReference(getprop.Object().Node())
}

// NodeUtil_IsObjectCreateWithPropertiesDefinition returns true if the node is
// a 2-argument Object.create(proto, descriptors).
//
// Java: NodeUtil.isObjectCreateWithPropertiesDefinition(Node)
func NodeUtil_IsObjectCreateWithPropertiesDefinition(n rhino.Node) bool {
	call, ok := ast.AsCall(n)
	if !ok || !n.HasXChildren(3) {
		return false
	}
	getprop, ok := ast.AsGetProp(call.Callee().Node())
	return ok &&
		getprop.Property() == "create" &&
		NodeUtil_IsKnownGlobalObjectReference(getprop.Object().Node())
}

// Java: NodeUtil.isKnownGlobalObjectReference(Node)
func NodeUtil_IsKnownGlobalObjectReference(n rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_NAME:
		return n.GetString() == "Object"
	case rhino.Token_GETPROP:
		return NodeUtil_GLOBAL_OBJECT.Matches(n) || NodeUtil_GLOBAL_OBJECT_MANGLED.Matches(n)
	default:
		return false
	}
}

// NodeUtil_IsObjectDefinePropertyDefinition returns true if the node is a
// definition with Object.defineProperty.
//
// Java: NodeUtil.isObjectDefinePropertyDefinition(Node)
func NodeUtil_IsObjectDefinePropertyDefinition(n rhino.Node) bool {
	call, ok := ast.AsCall(n)
	if !ok || !n.HasXChildren(4) {
		return false
	}
	getprop, ok := ast.AsGetProp(call.Callee().Node())
	return ok &&
		getprop.Property() == "defineProperty" &&
		NodeUtil_IsKnownGlobalObjectReference(getprop.Object().Node())
}

// NodeUtil_GetObjectDefinedPropertiesKeys returns a list of STRING_KEY
// properties defined by a Object.defineProperties(o, {...}) call.
//
// Java: NodeUtil.getObjectDefinedPropertiesKeys(Node)
func NodeUtil_GetObjectDefinedPropertiesKeys(definePropertiesCall rhino.Node) []rhino.Node {
	jutil.CheckArgument(NodeUtil_IsObjectDefinePropertiesDefinition(definePropertiesCall))
	var properties []rhino.Node
	if obj, ok := ast.AsObjectLit(definePropertiesCall.GetLastChild()); ok {
		for member := range obj.Members() {
			if key, ok := member.AsStringKey(); ok {
				properties = append(properties, key.Node())
			}
		}
	}
	return properties
}

// NodeUtil_IsPrototypePropertyDeclaration returns true if the node an
// assignment to a prototype property of some constructor.
//
// Java: NodeUtil.isPrototypePropertyDeclaration(Node)
func NodeUtil_IsPrototypePropertyDeclaration(n rhino.Node) bool {
	return NodeUtil_IsExprAssign(n) && NodeUtil_IsPrototypeProperty(n.GetFirstFirstChild())
}

// NodeUtil_IsPrototypeProperty returns whether the node represents a
// qualified prototype property.
//
// Java: NodeUtil.isPrototypeProperty(Node)
func NodeUtil_IsPrototypeProperty(n rhino.Node) bool {
	prop, ok := ast.AsGetProp(n)
	if !ok {
		return false
	}
	recv, ok := ast.AsGetProp(prop.Object().Node())
	return ok && recv.Property() == "prototype"
}

// NodeUtil_IsPrototypeMethod returns whether the node represents a prototype
// method.
//
// Java: NodeUtil.isPrototypeMethod(Node)
func NodeUtil_IsPrototypeMethod(n rhino.Node) bool {
	if !n.IsFunction() {
		return false
	}
	assignNode := n.GetParent()
	if !assignNode.IsAssign() {
		return false
	}
	return NodeUtil_IsPrototypePropertyDeclaration(assignNode.GetParent())
}

// Java: NodeUtil.isPrototypeAssignment(Node)
func NodeUtil_IsPrototypeAssignment(getProp rhino.Node) bool {
	prop, ok := ast.AsGetProp(getProp)
	return ok &&
		ast.RoleOf(getProp) == ast.RoleAssignTarget &&
		prop.Property() == "prototype"
}

// NodeUtil_IsPropertyTest determines whether this node is testing for the
// existence of a property. If true, we will not emit warnings about a
// missing property.
//
// propAccess is the GETPROP or GETELEM being tested.
//
// Java: NodeUtil.isPropertyTest(AbstractCompiler, Node)
func NodeUtil_IsPropertyTest(compiler AbstractCompiler, propAccess rhino.Node) bool {
	parent := propAccess.GetParent()
	switch parent.GetToken() {
	case rhino.Token_CALL:
		return parent.GetFirstChild() != propAccess &&
			compiler.GetCodingConvention().IsPropertyTestFunction(parent)
	case rhino.Token_IF, rhino.Token_WHILE, rhino.Token_DO, rhino.Token_FOR, rhino.Token_FOR_IN:
		return NodeUtil_GetConditionExpression(parent) == propAccess

	case rhino.Token_INSTANCEOF, rhino.Token_TYPEOF, rhino.Token_AND, rhino.Token_OR, rhino.Token_COALESCE,
		rhino.Token_OPTCHAIN_GETPROP:
		return true

	case rhino.Token_NE, rhino.Token_SHNE:
		var other rhino.Node
		if parent.GetFirstChild() == propAccess {
			other = parent.GetSecondChild()
		} else {
			other = parent.GetFirstChild()
		}
		return NodeUtil_IsUndefined(other) || (parent.IsNE() && other.IsNull())
	case rhino.Token_OPTCHAIN_CALL, rhino.Token_OPTCHAIN_GETELEM, rhino.Token_HOOK:
		return parent.GetFirstChild() == propAccess
	case rhino.Token_NOT:
		return parent.GetParent().IsOr() && parent.GetParent().GetFirstChild() == parent
	case rhino.Token_CAST:
		return NodeUtil_IsPropertyTest(compiler, parent)
	default:
		return false
	}
}

// Java: NodeUtil.isPropertyAbsenceTest(Node)
func NodeUtil_IsPropertyAbsenceTest(propAccess rhino.Node) bool {
	parent := propAccess.GetParent()
	switch parent.GetToken() {
	case rhino.Token_EQ, rhino.Token_SHEQ:
		var other rhino.Node
		if parent.GetFirstChild() == propAccess {
			other = parent.GetSecondChild()
		} else {
			other = parent.GetFirstChild()
		}
		return NodeUtil_IsUndefined(other) || (parent.IsEQ() && other.IsNull())
	default:
		return false
	}
}

// NodeUtil_GetPrototypeClassName: qName is a qualified name node
// representing a class prototype, or a property on that prototype, e.g.
// foo.Bar.prototype, or foo.Bar.prototype.toString. Returns the class name
// part of a qualified prototype name, e.g. foo.Bar.
//
// Java: NodeUtil.getPrototypeClassName(Node)
func NodeUtil_GetPrototypeClassName(qName rhino.Node) rhino.Node {
	prop, ok := ast.AsGetProp(qName)
	if !ok {
		return rhino.NilNode
	}
	if prop.Property() == "prototype" {
		return prop.Object().Node()
	}
	if recv, ok := ast.AsGetProp(prop.Object().Node()); ok && recv.Property() == "prototype" {
		return recv.Object().Node()
	}
	return rhino.NilNode
}

// NodeUtil_GetPrototypePropertyName returns the string property name part of
// a qualified prototype name.
//
// Java: NodeUtil.getPrototypePropertyName(Node)
func NodeUtil_GetPrototypePropertyName(qName rhino.Node) string {
	qNameStr, ok := qName.GetQualifiedNameString()
	if !ok {
		panic(jutil.NewNullPointerException())
	}
	prototypeIdx := jutil.StringLastIndexOf(qNameStr, ".prototype.")
	memberIndex := prototypeIdx + jutil.StringLength(".prototype") + 1
	return jutil.StringSubstringFrom(qNameStr, memberIndex)
}

// NodeUtil_GetFunctionParameters: fnNode is the function. Returns the Node
// containing the Function parameters.
//
// Java: NodeUtil.getFunctionParameters(Node)
func NodeUtil_GetFunctionParameters(fnNode rhino.Node) rhino.Node {
	jutil.CheckArgument(fnNode.IsFunction())
	return ast.MustFunction(fnNode).Params().Node()
}

// NodeUtil_HasNonSimpleParameters returns whether the given function node
// has any non-simple parameters (e.g. default values, rest parameters, or
// destructuring patterns).
//
// fnNode is a FUNCTION node. Returns true if function parameters are
// non-simple; false otherwise.
//
// Java: NodeUtil.hasNonSimpleParameters(Node)
func NodeUtil_HasNonSimpleParameters(fnNode rhino.Node) bool {
	jutil.CheckArgument(fnNode.IsFunction())
	for param := range ast.MustFunction(fnNode).Params().Params() {
		if _, ok := param.AsName(); !ok {
			return true
		}
	}
	return false
}

// NodeUtil_FunctionParametersMayHaveSideEffects returns whether any parameter
// of the given function node may have side effects when evaluated, or may
// throw an exception during parameter pattern matching / destructuring.
//
// fnNode is a FUNCTION node; astAnalyzer is the AstAnalyzer to check side
// effects of default parameter initializers. Returns true if function
// parameters may have side effects or throw; false otherwise.
//
// Java: NodeUtil.functionParametersMayHaveSideEffects(Node, AstAnalyzer)
func NodeUtil_FunctionParametersMayHaveSideEffects(fnNode rhino.Node, astAnalyzer *AstAnalyzer) bool {
	jutil.CheckArgument(fnNode.IsFunction())
	for param := range ast.MustFunction(fnNode).Params().Params() {
		if _, ok := param.AsName(); ok {
			continue
		}
		if rest, ok := param.AsIterRest(); ok {
			if _, ok := rest.Target().AsName(); ok {
				continue
			}
			return true
		}
		if def, ok := param.AsDefaultValue(); ok {
			if _, ok := def.Target().AsName(); !ok {
				return true
			}
			if astAnalyzer.MayHaveSideEffects(def.Value().Node()) {
				return true
			}
			continue
		}
		return true
	}
	return false
}

// NodeUtil_GetArgumentForFunction: given the function, this returns the nth
// argument or null if no such parameter exists.
//
// Java: NodeUtil.getArgumentForFunction(Node, int)
func NodeUtil_GetArgumentForFunction(function rhino.Node, index int) rhino.Node {
	jutil.CheckState(function.IsFunction())
	if index < 0 {
		return NodeUtil_GetNthSibling(function.GetSecondChild().GetFirstChild(), index)
	}
	param, _ := ast.MustFunction(function).Params().ParamAt(index)
	return param.Node()
}

// NodeUtil_GetArgumentForCallOrNew: given the new or call, this returns the
// nth argument of the call or null if no such argument exists.
//
// Java: NodeUtil.getArgumentForCallOrNew(Node, int)
func NodeUtil_GetArgumentForCallOrNew(call rhino.Node, index int) rhino.Node {
	jutil.CheckState(NodeUtil_IsCallOrNew(call))
	if index < 0 {
		return NodeUtil_GetNthSibling(call.GetSecondChild(), index)
	}
	arg, _ := ast.MustCallOrNew(call).ArgAt(index)
	return arg.Node()
}

// NodeUtil_IsInvocationTarget returns whether this is a target of a call or
// new.
//
// Java: NodeUtil.isInvocationTarget(Node)
func NodeUtil_IsInvocationTarget(n rhino.Node) bool {
	switch ast.RoleOf(n) {
	case ast.RoleCallCallee, ast.RoleOptChainCallCallee, ast.RoleNewCallee, ast.RoleTaggedTemplateLitTag:
		return true
	default:
		return false
	}
}

// NodeUtil_IsInvocation returns whether this is a call (including tagged
// template lits) or new.
//
// Java: NodeUtil.isInvocation(Node)
func NodeUtil_IsInvocation(n rhino.Node) bool {
	_, ok := ast.AsInvocation(n)
	return ok
}

// Java: NodeUtil.isCallOrNewArgument(Node)
func NodeUtil_IsCallOrNewArgument(n rhino.Node) bool {
	switch ast.RoleOf(n) {
	case ast.RoleCallArgs, ast.RoleOptChainCallArgs, ast.RoleNewArgs:
		return true
	default:
		return false
	}
}

// Java: NodeUtil.isToStringMethodCall(Node)
func NodeUtil_IsToStringMethodCall(call rhino.Node) bool {
	if get, ok := ast.AsNormalOrOptChainGet(call.GetFirstChild()); ok {
		name, ok := get.StaticPropertyName()
		return ok && name == "toString"
	}
	return false
}

// NodeUtil_IsRhsClass: e.g. `class A` in `/** some */ var x = class A {};`
//
// Java: NodeUtil.isRhsClass(Node)
func NodeUtil_IsRhsClass(n rhino.Node) bool {
	return NodeUtil_IsNamedClassExpression(n) &&
		(n.GetParent().IsName() || n.GetParent().IsAssign())
}

// NodeUtil_IsRhsFunction: e.g. `function A` in
// `/** some */ var x = function A() {};`
//
// Java: NodeUtil.isRhsFunction(Node)
func NodeUtil_IsRhsFunction(n rhino.Node) bool {
	return NodeUtil_IsNamedFunctionExpression(n) &&
		(n.GetParent().IsName() || n.GetParent().IsAssign())
}

// Java: NodeUtil.isConstructor(Node)
func NodeUtil_IsConstructor(fnNode rhino.Node) bool {
	if fnNode.IsNil() || !fnNode.IsFunction() {
		return false
	}

	type_ := fnNode.GetJSType()
	jsDocInfo := NodeUtil_GetBestJSDocInfo(fnNode)
	color := fnNode.GetColor()

	return (type_ != nil && type_.IsConstructor()) ||
		(jsDocInfo != nil && jsDocInfo.IsConstructor()) ||
		(color != nil && color.IsConstructor()) ||
		NodeUtil_IsEs6Constructor(fnNode)
}

// Java: NodeUtil.isEs6ConstructorMemberFunctionDef(Node)
func NodeUtil_IsEs6ConstructorMemberFunctionDef(memberFunctionDef rhino.Node) bool {
	m, ok := ast.AsMemberFunctionDef(memberFunctionDef)
	return ok && m.IsEs6Constructor()
}

// Java: NodeUtil.isEs6Constructor(Node)
func NodeUtil_IsEs6Constructor(fnNode rhino.Node) bool {
	fn, ok := ast.AsFunction(fnNode)
	return ok && fn.IsEs6Constructor()
}

// Java: NodeUtil.isGetterOrSetter(Node)
func NodeUtil_IsGetterOrSetter(propNode rhino.Node) bool {
	if NodeUtil_IsGetOrSetKey(propNode) {
		return true
	}
	key, ok := ast.AsStringKey(propNode)
	if !ok || !key.Value().Node().IsFunction() {
		return false
	}
	keyName := key.Key()
	return keyName == "get" || keyName == "set"
}

// Java: NodeUtil.isCallTo(Node, String)
func NodeUtil_IsCallTo(n rhino.Node, qualifiedName string) bool {
	return n.IsCall() && n.GetFirstChild().MatchesQualifiedName(qualifiedName)
}

// Java: NodeUtil.isCallTo(Node, QualifiedName)
func NodeUtil_IsCallToQualifiedName(n rhino.Node, qualifiedName rhino.QualifiedName) bool {
	return n.IsCall() && qualifiedName.Matches(n.GetFirstChild())
}

// NodeUtil_IsCallToNode is a faster version of NodeUtil_IsCallTo.
//
// n is the node to check if a call; targetMethod is the prebuilt AST getprop
// node that represents the method to check.
//
// Java: NodeUtil.isCallTo(Node, Node)
func NodeUtil_IsCallToNode(n rhino.Node, targetMethod rhino.Node) bool {
	if !n.IsCall() {
		return false
	}

	return n.GetFirstChild().MatchesQualifiedNameNode(targetMethod)
}

// Java: NodeUtil.getInvocationArgsAsIterable(Node)
func NodeUtil_GetInvocationArgsAsIterable(invocation rhino.Node) iter.Seq[rhino.Node] {
	if invocation.IsTaggedTemplateLit() {
		return NewNodeUtil_TemplateArgsIterable(invocation.GetLastChild()).All()
	}

	jutil.CheckStateMsg1(NodeUtil_IsCallOrNew(invocation), "%s", invocation)
	if invocation.HasOneChild() {
		return slices.Values([]rhino.Node(nil))
	}

	var list []rhino.Node
	for arg := invocation.GetSecondChild(); !arg.IsNil(); arg = arg.GetNext() {
		list = append(list, arg)
	}
	return slices.Values(list)
}

// NodeUtil_GetInvocationArgsCount returns the number of arguments in this
// invocation. For template literals it takes into account the implicit first
// argument of ITemplateArray.
//
// Java: NodeUtil.getInvocationArgsCount(Node)
func NodeUtil_GetInvocationArgsCount(invocation rhino.Node) int {
	if invocation.IsTaggedTemplateLit() {
		args := NewNodeUtil_TemplateArgsIterable(invocation.GetLastChild())
		size := 0
		for range args.All() {
			size++
		}
		return size + 1
	} else {
		return invocation.GetChildCount() - 1
	}
}

// NodeUtil_TemplateArgsIterable represents an iterable of the children of
// templatelit_sub nodes of a template lit node This iterable will skip over
// the String children of the template lit node.
//
// Java: NodeUtil.TemplateArgsIterable
type NodeUtil_TemplateArgsIterable struct {
	templateLit rhino.Node
}

// Java: NodeUtil.TemplateArgsIterable(Node)
func NewNodeUtil_TemplateArgsIterable(templateLit rhino.Node) *NodeUtil_TemplateArgsIterable {
	jutil.CheckState(templateLit.IsTemplateLit())
	return &NodeUtil_TemplateArgsIterable{templateLit: templateLit}
}

// Iterator returns the anonymous AbstractIterator<Node> of the Java class.
//
// Java: NodeUtil.TemplateArgsIterable.iterator()
func (t *NodeUtil_TemplateArgsIterable) Iterator() *nodeUtil_templateArgsIterator {
	return &nodeUtil_templateArgsIterator{nextChild: t.templateLit.GetFirstChild()}
}

// All iterates the iterable.
//
// Go-only: language-adapter — Java's for-each over the Iterable.
func (t *NodeUtil_TemplateArgsIterable) All() iter.Seq[rhino.Node] {
	return func(yield func(rhino.Node) bool) {
		it := t.Iterator()
		for it.HasNext() {
			if !yield(it.Next()) {
				return
			}
		}
	}
}

// nodeUtil_templateArgsIterator is the anonymous AbstractIterator<Node>
// returned by TemplateArgsIterable.iterator().
type nodeUtil_templateArgsIterator struct {
	nextChild rhino.Node

	// Go-only: language-adapter — AbstractIterator's state.
	next    rhino.Node
	hasNext bool
	done    bool
}

// HasNext is AbstractIterator.hasNext().
func (it *nodeUtil_templateArgsIterator) HasNext() bool {
	if it.done {
		return false
	}
	if !it.hasNext {
		next, ok := it.ComputeNext()
		if !ok {
			it.done = true
			return false
		}
		it.next = next
		it.hasNext = true
	}
	return true
}

// Next is AbstractIterator.next().
func (it *nodeUtil_templateArgsIterator) Next() rhino.Node {
	if !it.HasNext() {
		panic(jutil.NewNoSuchElementException())
	}
	it.hasNext = false
	result := it.next
	it.next = rhino.NilNode
	return result
}

// ComputeNext returns the next argument, or ok=false for endOfData().
//
// Java: NodeUtil.TemplateArgsIterable.iterator().computeNext()
func (it *nodeUtil_templateArgsIterator) ComputeNext() (next rhino.Node, ok bool) {
	for !it.nextChild.IsNil() && !it.nextChild.IsTemplateLitSub() {
		it.nextChild = it.nextChild.GetNext()
	}
	if it.nextChild.IsNil() {
		return rhino.NilNode, false
	} else {
		result := it.nextChild.GetFirstChild()
		it.nextChild = it.nextChild.GetNext()
		return result, true
	}
}

// NodeUtil_IsObjLitProperty returns true if the node is a property of an
// object literal.
//
// Java: NodeUtil.isObjLitProperty(Node)
func NodeUtil_IsObjLitProperty(node rhino.Node) bool {
	return node.IsStringKey() ||
		node.IsGetterDef() ||
		node.IsSetterDef() ||
		node.IsMemberFunctionDef() ||
		node.IsComputedProp()
}

// NodeUtil_IsBlocklessArrowFunctionResult returns whether the node
// represents the return value of a blockless Arrow function.
//
// Java: NodeUtil.isBlocklessArrowFunctionResult(Node)
func NodeUtil_IsBlocklessArrowFunctionResult(n rhino.Node) bool {
	return ast.RoleOf(n) == ast.RoleFunctionBody && !n.IsBlock()
}
