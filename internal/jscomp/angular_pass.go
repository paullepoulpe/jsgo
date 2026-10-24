// Port of com/google/javascript/jscomp/AngularPass.java @ closure-compiler v20261021

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
	"slices"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// AngularPass is the compiler pass for AngularJS-specific needs. Generates
// `$inject` properties for functions (class constructors, wrappers, etc)
// annotated with `@ngInject`. Without this pass, AngularJS will not work
// properly if variable renaming is enabled, because the function arguments
// will be renamed.
//
// See http://docs.angularjs.org/tutorial/step_05#a-note-on-minification
//
// For example, the following code:
//
//	/** @ngInject * /
//	function Controller(dependency1, dependency2) {
//	// do something
//	}
//
// will be transformed into:
//
//	function Controller(dependency1, dependency2) {
//	// do something
//	}
//	Controller.$inject = ['dependency1', 'dependency2'];
//
// This pass also supports assignments of function expressions to variables
// like:
//
//	/** @ngInject * /
//	var filter = function(a, b) {};
//
//	var ns = {};
//	/** @ngInject * /
//	ns.method = function(a,b,c) {};
//
//	/** @ngInject * /
//	var shorthand = ns.method2 = function(a,b,c,) {}
//
// Java: class AngularPass extends AbstractPostOrderCallback implements CompilerPass
type AngularPass struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	compiler AbstractCompiler

	// injectables are the nodes annotated with @ngInject
	injectables *jutil.ArrayList[*AngularPass_NodeContext]
}

// Java: AngularPass(AbstractCompiler)
func NewAngularPass(compiler AbstractCompiler) *AngularPass {
	return &AngularPass{
		compiler:    compiler,
		injectables: jutil.NewArrayList[*AngularPass_NodeContext](),
	}
}

const AngularPass_INJECT_PROPERTY_NAME = "$inject"

var AngularPass_INJECT_IN_NON_GLOBAL_OR_BLOCK_ERROR = DiagnosticType_Error("JSC_INJECT_IN_NON_GLOBAL_OR_BLOCK_ERROR",
	"@ngInject only applies to functions defined in blocks or "+
		"global scope.")

var AngularPass_INJECT_NON_FUNCTION_ERROR = DiagnosticType_Error("JSC_INJECT_NON_FUNCTION_ERROR",
	"@ngInject can only be used when defining a function or "+
		"assigning a function expression.")

var AngularPass_INJECTED_FUNCTION_HAS_DESTRUCTURED_PARAM = DiagnosticType_Error("JSC_INJECTED_FUNCTION_HAS_DESTRUCTURED_PARAM",
	"@ngInject cannot be used on functions containing "+
		"destructured parameter.")

var AngularPass_INJECTED_FUNCTION_HAS_DEFAULT_VALUE = DiagnosticType_Error("JSC_INJECTED_FUNCTION_HAS_DEFAULT_VALUE",
	"@ngInject cannot be used on functions containing default value.")

var AngularPass_INJECTED_FUNCTION_ON_NON_QNAME = DiagnosticType_Error("JSC_INJECTED_FUNCTION_ON_NON_QNAME",
	"@ngInject can only be used on qualified names.")

// Java: AngularPass.process(Node, Node)
func (p *AngularPass) Process(externs *rhino.Node, root *rhino.Node) {
	// Traverses AST looking for nodes annotated with @ngInject.
	NodeTraversal_Traverse(p.compiler, root, p)
	// iterates through annotated nodes adding $inject property to elements.
	for entry := range p.injectables.All() {
		name := entry.GetName()
		fn := entry.GetFunctionNode()
		dependencies := p.createDependenciesList(fn)
		// skips entry if it does not have any dependencies.
		if dependencies.IsEmpty() {
			continue
		}
		dependenciesArray := rhino.IR_ArraylitOfIterable(slices.Values(dependencies.Slice()))
		// creates `something.$inject = ['param1', 'param2']` node.
		statement :=
			rhino.IR_ExprResult(
				rhino.IR_Assign(
					rhino.IR_Getelem(NodeUtil_NewQName(p.compiler, name), rhino.IR_String(AngularPass_INJECT_PROPERTY_NAME)),
					dependenciesArray))
		statement.SrcrefTree(entry.GetNode())
		statement.SetOriginalName(&name)
		// Set the visibility of the newly created property.
		newPropertyDoc := rhino.JSDocInfo_Builder_()
		newPropertyDoc.RecordVisibility(rhino.JSDocInfo_Visibility_PUBLIC)
		statement.GetFirstChild().SetJSDocInfo(newPropertyDoc.Build())

		// adds `something.$inject = [...]` node after the annotated node or the following
		// goog.inherits call.
		insertionPoint := entry.GetTarget()
		next := insertionPoint.GetNext()
		for next != nil &&
			NodeUtil_IsExprCall(next) &&
			p.compiler.GetCodingConvention().GetClassesDefinedByCall(
				next.GetFirstChild()) != nil {
			insertionPoint = next
			next = insertionPoint.GetNext()
		}

		statement.InsertAfter(insertionPoint)
		p.compiler.ReportChangeToEnclosingScope(statement)
	}
}

// createDependenciesList: given a FUNCTION node returns array of STRING
// nodes representing function parameters.
//
// n is the FUNCTION node. Returns STRING nodes.
//
// Java: AngularPass.createDependenciesList(Node)
func (p *AngularPass) createDependenciesList(n *rhino.Node) *jutil.ArrayList[*rhino.Node] {
	jutil.CheckArgument(n.IsFunction())
	params := NodeUtil_GetFunctionParameters(n)
	if params != nil {
		return p.createStringsFromParamList(params)
	}
	return jutil.NewArrayList[*rhino.Node]()
}

// createStringsFromParamList: given a PARAM_LIST node creates an array of
// corresponding STRING nodes.
//
// params is the PARAM_LIST node. Returns array of STRING nodes.
//
// Java: AngularPass.createStringsFromParamList(Node)
func (p *AngularPass) createStringsFromParamList(params *rhino.Node) *jutil.ArrayList[*rhino.Node] {
	param := params.GetFirstChild()
	names := jutil.NewArrayList[*rhino.Node]()
	for param != nil {
		if param.IsName() {
			names.Add(rhino.IR_String(param.GetString()).Srcref(param))
		} else if param.IsDestructuringPattern() {
			p.compiler.Report(JSError_MakeWithNode(param,
				AngularPass_INJECTED_FUNCTION_HAS_DESTRUCTURED_PARAM))
			return jutil.NewArrayList[*rhino.Node]()
		} else if param.IsDefaultValue() {
			p.compiler.Report(JSError_MakeWithNode(param,
				AngularPass_INJECTED_FUNCTION_HAS_DEFAULT_VALUE))
			return jutil.NewArrayList[*rhino.Node]()
		}
		param = param.GetNext()
	}
	return names
}

// Java: AngularPass.visit(NodeTraversal, Node, Node)
func (p *AngularPass) Visit(unused *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	docInfo := n.GetJSDocInfo()
	if docInfo != nil && docInfo.IsNgInject() {
		p.addNode(n)
	}
}

// addNode adds node to the list of injectables.
//
// n is the node to add.
//
// Java: AngularPass.addNode(Node)
func (p *AngularPass) addNode(n *rhino.Node) {
	var injectAfter *rhino.Node
	var fn *rhino.Node
	var name *string

	switch n.GetToken() {
	case rhino.Token_ASSIGN:
		// handles assignment cases like:
		// a = function() {}
		// a = b = c = function() {}
		if !n.GetFirstChild().IsQualifiedName() {
			p.compiler.Report(JSError_MakeWithNode(n, AngularPass_INJECTED_FUNCTION_ON_NON_QNAME))
			return
		}
		name = n.GetFirstChild().GetQualifiedName()
		// last node of chained assignment.
		fn = n
		for fn.IsAssign() {
			fn = fn.GetLastChild()
		}
		injectAfter = n.GetParent()
	case rhino.Token_FUNCTION:
		// handles function case:
		// function fnName() {}
		name = NodeUtil_GetName(n)
		fn = n
		injectAfter = n
		if n.GetParent().IsAssign() && n.GetParent().GetJSDocInfo().IsNgInject() {
			// This is a function assigned into a symbol, e.g. a regular function
			// declaration in a goog.module or goog.scope.
			// Skip in this traversal, it is handled when visiting the assign.
			return
		}
	case rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST:
		// handles var declaration cases like:
		// var a = function() {}
		// var a = b = function() {}
		name = jutil.StringPtr(n.GetFirstChild().GetString())
		// looks for a function node.
		fn = angularPass_getDeclarationRValue(n)
		injectAfter = n
	case rhino.Token_MEMBER_FUNCTION_DEF:
		// handles class method case:
		// class clName(){
		//   constructor(){}
		//   someMethod(){} <===
		// }
		parent := n.GetParent()
		if parent.IsClassMembers() {
			classNode := parent.GetParent()
			midPart := ".prototype."
			if n.IsStaticMember() {
				midPart = "."
			}
			name = jutil.StringPtr(jutil.ToJavaString(NodeUtil_GetName(classNode)) + midPart + n.GetString())
			if NodeUtil_IsEs6ConstructorMemberFunctionDef(n) {
				name = NodeUtil_GetName(classNode)
			}
			fn = n.GetFirstChild()
			injectAfter = NodeUtil_GetEnclosingStatement(classNode)
		}
	default:
	}

	if fn == nil || !fn.IsFunction() {
		p.compiler.Report(JSError_MakeWithNode(n, AngularPass_INJECT_NON_FUNCTION_ERROR))
		return
	}
	if injectAfter.GetParent().IsExport() {
		// handle `export class Foo {` or `export function(`
		injectAfter = injectAfter.GetParent()
	}
	// report an error if the function declaration did not take place in the root of a statement
	// for example, `fn(/** @inject */ function(x) {});` is forbidden
	if !NodeUtil_IsStatementBlock(injectAfter.GetParent()) {
		p.compiler.Report(JSError_MakeWithNode(n, AngularPass_INJECT_IN_NON_GLOBAL_OR_BLOCK_ERROR))
		return
	}
	// checks that name is present, which must always be the case unless the
	// compiler allowed a syntax error or a dangling anonymous function
	// expression.
	jutil.CheckNotNull(name)
	// registers the node.
	p.injectables.Add(NewAngularPass_NodeContext(*name, n, fn, injectAfter))
}

// angularPass_getDeclarationRValue: given a VAR node (variable declaration)
// returns the node of initial value.
//
//	var x;  // null
//	var y = "value"; // STRING "value" node
//	var z = x = y = function() {}; // FUNCTION node
//
// n is the VAR node. Returns the assigned initial value, or the rightmost
// rvalue of an assignment chain, or null.
//
// Java: AngularPass.getDeclarationRValue(Node)
func angularPass_getDeclarationRValue(n *rhino.Node) *rhino.Node {
	jutil.CheckNotNull(n)
	jutil.CheckArgument(NodeUtil_IsNameDeclaration(n))
	n = n.GetFirstFirstChild()
	if n == nil {
		return nil
	}
	for n.IsAssign() {
		n = n.GetLastChild()
	}
	return n
}

// AngularPass_NodeContext is an annotated node.
//
// Java: static class AngularPass.NodeContext
type AngularPass_NodeContext struct {
	// name is the name of the function/object.
	name string
	// node is the node jsDoc is attached to.
	node *rhino.Node
	// functionNode is the function node
	functionNode *rhino.Node
	// target is the node after which to inject the new code
	target *rhino.Node
}

// Java: AngularPass.NodeContext(String, Node, Node, Node)
func NewAngularPass_NodeContext(name string, node *rhino.Node, functionNode *rhino.Node, target *rhino.Node) *AngularPass_NodeContext {
	return &AngularPass_NodeContext{
		name:         name,
		node:         node,
		functionNode: functionNode,
		target:       target,
	}
}

// GetName returns the name.
//
// Java: AngularPass.NodeContext.getName()
func (c *AngularPass_NodeContext) GetName() string {
	return c.name
}

// GetNode returns the node.
//
// Java: AngularPass.NodeContext.getNode()
func (c *AngularPass_NodeContext) GetNode() *rhino.Node {
	return c.node
}

// GetFunctionNode returns the context.
//
// Java: AngularPass.NodeContext.getFunctionNode()
func (c *AngularPass_NodeContext) GetFunctionNode() *rhino.Node {
	return c.functionNode
}

// GetTarget returns the context.
//
// Java: AngularPass.NodeContext.getTarget()
func (c *AngularPass_NodeContext) GetTarget() *rhino.Node {
	return c.target
}
