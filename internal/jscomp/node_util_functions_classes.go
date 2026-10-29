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
	return func(yield func(rhino.Node) bool) {
		for arg := invocation.GetSecondChild(); !arg.IsNil(); {
			next := arg.GetNext()
			if !yield(arg) {
				return
			}
			arg = next
		}
	}
}

// NodeUtil_GetInvocationArgsCount returns the number of arguments in this
// invocation. For template literals it takes into account the implicit first
// argument of ITemplateArray.
//
// Java: NodeUtil.getInvocationArgsCount(Node)
func NodeUtil_GetInvocationArgsCount(invocation rhino.Node) int {
	if invocation.IsTaggedTemplateLit() {
		templateLit := invocation.GetLastChild()
		jutil.CheckState(templateLit.IsTemplateLit())
		size := 1
		for child := templateLit.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
			if child.IsTemplateLitSub() {
				size++
			}
		}
		return size
	}
	return invocation.GetChildCount() - 1
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
func NewNodeUtil_TemplateArgsIterable(templateLit rhino.Node) NodeUtil_TemplateArgsIterable {
	jutil.CheckState(templateLit.IsTemplateLit())
	return NodeUtil_TemplateArgsIterable{templateLit: templateLit}
}

// All iterates the iterable.
//
// Go-only: language-adapter — Java's for-each over the Iterable.
func (t NodeUtil_TemplateArgsIterable) All() iter.Seq[rhino.Node] {
	templateLit := t.templateLit
	return func(yield func(rhino.Node) bool) {
		for child := templateLit.GetFirstChild(); !child.IsNil(); {
			next := child.GetNext()
			if child.IsTemplateLitSub() {
				if !yield(child.GetFirstChild()) {
					return
				}
			}
			child = next
		}
	}
}

// NodeUtil_IsBlocklessArrowFunctionResult returns whether the node
// represents the return value of a blockless Arrow function.
//
// Java: NodeUtil.isBlocklessArrowFunctionResult(Node)
func NodeUtil_IsBlocklessArrowFunctionResult(n rhino.Node) bool {
	return ast.RoleOf(n) == ast.RoleFunctionBody && !n.IsBlock()
}
