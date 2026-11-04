// Port of com/google/javascript/jscomp/AngularPass.java @ closure-compiler v20261028

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

	"github.com/paullepoulpe/jsgo/internal/ast"
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
	injectables []angularInjectable
}

type angularInjectable struct {
	name         string
	node         rhino.Node
	functionNode ast.Function
	target       rhino.Node
}

// Java: AngularPass(AbstractCompiler)
func NewAngularPass(compiler AbstractCompiler) *AngularPass {
	return &AngularPass{
		compiler: compiler,
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
func (p *AngularPass) Process(externs rhino.Node, root rhino.Node) {
	// Traverses AST looking for nodes annotated with @ngInject.
	NodeTraversal_Traverse(p.compiler, root, p)
	// iterates through annotated nodes adding $inject property to elements.
	for i := range p.injectables {
		entry := &p.injectables[i]
		dependencies := p.createDependenciesList(entry.functionNode)
		// skips entry if it does not have any dependencies.
		if len(dependencies) == 0 {
			continue
		}
		dependenciesArray := rhino.IR_ArraylitOfIterable(slices.Values(dependencies))
		// creates `something.$inject = ['param1', 'param2']` node.
		assign := rhino.IR_Assign(
			rhino.IR_Getelem(NodeUtil_NewQName(p.compiler, entry.name), rhino.IR_String(AngularPass_INJECT_PROPERTY_NAME)),
			dependenciesArray)
		statement := rhino.IR_ExprResult(assign).SrcrefTree(entry.node)
		statement.SetOriginalName(entry.name)
		// Set the visibility of the newly created property.
		newPropertyDoc := rhino.JSDocInfo_Builder_()
		newPropertyDoc.RecordVisibility(rhino.JSDocInfo_Visibility_PUBLIC)
		assign.SetJSDocInfo(newPropertyDoc.Build())

		// adds `something.$inject = [...]` node after the annotated node or the following
		// goog.inherits call.
		insertionPoint := entry.target
		for next := insertionPoint.GetNext(); !next.IsNil(); next = insertionPoint.GetNext() {
			exprStmt, ok := ast.AsExprStmt(next)
			if !ok || !exprStmt.Expr().Node().IsCall() ||
				p.compiler.GetCodingConvention().GetClassesDefinedByCall(exprStmt.Expr().Node()) == nil {
				break
			}
			insertionPoint = next
		}

		statement.InsertAfter(insertionPoint)
		p.compiler.ReportChangeToEnclosingScope(statement)
	}
}

// createDependenciesList: given a FUNCTION node returns array of STRING
// nodes representing function parameters.
//
// Java: AngularPass.createDependenciesList(Node)
func (p *AngularPass) createDependenciesList(fn ast.Function) []rhino.Node {
	var names []rhino.Node
	for param := range fn.Params().Params() {
		if name, ok := param.AsName(); ok {
			names = append(names, rhino.IR_String(name.Ident()).Srcref(param.Node()))
		} else if _, ok := param.AsPattern(); ok {
			p.compiler.Report(JSError_MakeWithNode(param.Node(), AngularPass_INJECTED_FUNCTION_HAS_DESTRUCTURED_PARAM))
			return nil
		} else if _, ok := param.AsDefaultValue(); ok {
			p.compiler.Report(JSError_MakeWithNode(param.Node(), AngularPass_INJECTED_FUNCTION_HAS_DEFAULT_VALUE))
			return nil
		}
	}
	return names
}

// Java: AngularPass.visit(NodeTraversal, Node, Node)
func (p *AngularPass) Visit(unused *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if docInfo := n.GetJSDocInfo(); docInfo != nil && docInfo.IsNgInject() {
		p.addNode(n)
	}
}

// addNode adds node to the list of injectables.
//
// Java: AngularPass.addNode(Node)
func (p *AngularPass) addNode(n rhino.Node) {
	var injectAfter rhino.Node
	var fnNode rhino.Node
	var name string
	var hasName bool

	switch n.GetToken() {
	case rhino.Token_ASSIGN:
		assign := ast.MustAssign(n)
		lhs := assign.Target().Node()
		if !lhs.IsQualifiedName() {
			p.compiler.Report(JSError_MakeWithNode(n, AngularPass_INJECTED_FUNCTION_ON_NON_QNAME))
			return
		}
		name, hasName = lhs.GetQualifiedNameString()
		fnNode = unwrapAssignChain(assign.Value().Node())
		injectAfter = n.GetParent()
	case rhino.Token_FUNCTION:
		name, hasName = NodeUtil_GetName(n)
		fnNode, injectAfter = n, n
		if parent := n.GetParent(); parent.IsAssign() && parent.GetJSDocInfo() != nil && parent.GetJSDocInfo().IsNgInject() {
			// Handled when visiting the assign.
			return
		}
	case rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST:
		if nameNode, ok := ast.MustNameDeclaration(n).MustFirstDeclarator().AsName(); ok {
			name, hasName = nameNode.Ident(), true
			if init, hasInit := nameNode.Init(); hasInit {
				fnNode = unwrapAssignChain(init.Node())
			}
		}
		injectAfter = n
	case rhino.Token_MEMBER_FUNCTION_DEF:
		m := ast.MustMemberFunctionDef(n)
		if parent := n.GetParent(); parent.IsClassMembers() {
			classNode := parent.GetParent()
			if assign, ok := ast.AsAssign(classNode.GetParent()); ok &&
				assign.Value().Is(classNode) && !assign.Target().Node().IsQualifiedName() {
				p.compiler.Report(JSError_MakeWithNode(n, AngularPass_INJECTED_FUNCTION_ON_NON_QNAME))
				return
			}
			if m.IsEs6Constructor() {
				name, hasName = NodeUtil_GetName(classNode)
			} else if className, ok := NodeUtil_GetName(classNode); ok {
				midPart := ".prototype."
				if m.IsStatic() {
					midPart = "."
				}
				name, hasName = className+midPart+m.Key(), true
			}
			fnNode = m.Function().Node()
			injectAfter = angularClassInjectionTarget(classNode)
		}
	}

	fn, ok := ast.AsFunction(fnNode)
	if !ok {
		p.compiler.Report(JSError_MakeWithNode(n, AngularPass_INJECT_NON_FUNCTION_ERROR))
		return
	}
	if !injectAfter.IsNil() && !injectAfter.GetParent().IsNil() && injectAfter.GetParent().IsExport() {
		// handle `export class Foo {` or `export function(`
		injectAfter = injectAfter.GetParent()
	}
	if !hasName || injectAfter.IsNil() || injectAfter.GetParent().IsNil() || !NodeUtil_IsStatementBlock(injectAfter.GetParent()) {
		p.compiler.Report(JSError_MakeWithNode(n, AngularPass_INJECT_IN_NON_GLOBAL_OR_BLOCK_ERROR))
		return
	}
	p.injectables = append(p.injectables, angularInjectable{
		name:         name,
		node:         n,
		functionNode: fn,
		target:       injectAfter,
	})
}

func angularClassInjectionTarget(classNode rhino.Node) rhino.Node {
	cur := classNode
	for {
		assign, ok := ast.AsAssign(cur.GetParent())
		if !ok || !assign.Value().Is(cur) {
			break
		}
		cur = assign.Node()
	}
	parent := cur.GetParent()
	if parent.IsName() && NodeUtil_IsNameDeclaration(parent.GetParent()) {
		return parent.GetParent()
	}
	if cur != classNode && parent.IsExprResult() {
		return parent
	}
	return cur
}

func unwrapAssignChain(n rhino.Node) rhino.Node {
	for inner, ok := ast.AsAssign(n); ok; inner, ok = ast.AsAssign(n) {
		n = inner.Value().Node()
	}
	return n
}
