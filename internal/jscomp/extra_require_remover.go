// Port of com/google/javascript/jscomp/ExtraRequireRemover.java @ closure-compiler v20261021

/*
 * Copyright 2024 The Closure Compiler Authors.
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

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ExtraRequireRemover walks the AST looking for usages of qualified names,
// and 'goog.require's of those names. Then, reconciles the two lists, and
// removes any unnecessary require statements.
//
//	               Input js                    |     Generated i.js (with ExtraRequireRemover)
//	--------------------------------------------------------------------------------------------
//	    const Foo1 = goog.require('Foo');      |       const Foo1 = goog.require('Foo');
//	    const Bar1 = goog.require('Bar');      |
//	    /** @type {!Foo1} */                   |        /** @type {!Foo1} *\
//	    let foo;                               |        let foo;
//	    let a = A1();                          |        let a = A1();
//
// Note how `const Bar1 = goog.require('Bar');` is removed bc it is not used
// in the .i.js file.
type ExtraRequireRemover struct {
	compiler AbstractCompiler

	// Keys are the local name of a required namespace. Values are the
	// goog.require CALL node.
	requires *jutil.LinkedHashMap[string, rhino.Node]

	// Adding an entry to `usages` indicates that the name (either a fully
	// qualified or local name) is used and can be required. Note that since
	// `usages` are name-based and not scoped, any usage that shadows an
	// unused require in that file will cause the extra require not to be
	// removed.
	usages *jutil.LinkedHashSet[string]
}

var extraRequireRemover_DEFAULT_EXTRA_NAMESPACES = jutil.NewLinkedHashSet(
	"goog.testing.asserts", "goog.testing.jsunit", "goog.labs.testing.Environment")

// Java: ExtraRequireRemover(AbstractCompiler)
func NewExtraRequireRemover(compiler AbstractCompiler) *ExtraRequireRemover {
	e := &ExtraRequireRemover{}
	e.requires = jutil.NewLinkedHashMap[string, rhino.Node]()
	e.usages = jutil.NewLinkedHashSet[string]()
	e.compiler = compiler
	return e
}

// Java: ExtraRequireRemover.process(Node, Node)
func (e *ExtraRequireRemover) Process(externs rhino.Node, root rhino.Node) {
	e.Reset()
	NodeTraversal_Traverse(e.compiler, root, e)
}

// ExtractNamespace extracts the namespace from matching goog.require,
// goog.requireType, and goog.forwardDeclare nodes.
//
// Ex: `goog.require('foo.bar');` will return `foo.bar`
//
// Java: ExtraRequireRemover.extractNamespace(Node, QualifiedName...)
func (e *ExtraRequireRemover) ExtractNamespace(call rhino.Node, primitiveNames ...rhino.QualifiedName) *string {
	callee := call.GetFirstChild()
	if !callee.IsGetProp() {
		return nil
	}
	for _, primitiveName := range primitiveNames {
		if primitiveName.Matches(callee) {
			target := callee.GetNext()
			if !target.IsNil() && target.IsStringLit() {
				s := target.GetString()
				return &s
			}
		}
	}
	return nil
}

var extraRequireRemover_GOOG_REQUIRE = rhino.QualifiedName_Of("goog.require")
var extraRequireRemover_GOOG_REQUIRE_TYPE = rhino.QualifiedName_Of("goog.requireType")
var extraRequireRemover_GOOG_FORWARD_DECLARE = rhino.QualifiedName_Of("goog.forwardDeclare")

// ExtractNamespaceIfRequire extracts namespace from goog.require and
// goog.requireType nodes
//
// Java: ExtraRequireRemover.extractNamespaceIfRequire(Node)
func (e *ExtraRequireRemover) ExtractNamespaceIfRequire(call rhino.Node) *string {
	return e.ExtractNamespace(call, extraRequireRemover_GOOG_REQUIRE, extraRequireRemover_GOOG_REQUIRE_TYPE, extraRequireRemover_GOOG_FORWARD_DECLARE)
}

// Java: ExtraRequireRemover.shouldTraverse(NodeTraversal, Node, Node)
func (e *ExtraRequireRemover) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	// Skip ES modules. We could support them in the future if desired.
	return !n.IsModuleBody() || n.GetParent().GetBooleanProp(rhino.Node_GOOG_MODULE)
}

// Visit visits each corresponding node and tracks all goog.requires as well
// as their usages
//
// Java: ExtraRequireRemover.visit(NodeTraversal, Node, Node)
func (e *ExtraRequireRemover) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	e.MaybeAddJsDocUsages(n)
	switch n.GetToken() {
	case rhino.Token_NAME:
		if !NodeUtil_IsLValue(n) && !parent.IsGetProp() {
			e.VisitQualifiedName(n)
		}
	case rhino.Token_GETPROP:
		// If parent is a GETPROP, they will handle all the usages.
		if !parent.IsGetProp() && n.IsQualifiedName() {
			e.VisitQualifiedName(n)
		}
	case rhino.Token_CALL:
		e.VisitCallNode(n, parent)
	case rhino.Token_SCRIPT:
		e.RemoveExtraRequiresInScript()
		e.Reset()
	default:
	}
}

// Java: ExtraRequireRemover.reset()
func (e *ExtraRequireRemover) Reset() {
	e.usages.Clear()
	e.requires.Clear()
}

// RemoveExtraRequiresInScript: for every goog.require, check that there is
// a usage and remove the import if the goog.require is not found in
// `usages`
//
// Java: ExtraRequireRemover.removeExtraRequiresInScript()
func (e *ExtraRequireRemover) RemoveExtraRequiresInScript() {
	for require, call := range e.requires.All() {
		e.RemoveExtraRequire(call, require)
	}
}

// RemoveExtraRequire reconciles the two lists (`usages` and `requires`) and
// removes requires that are used (found in `usages`). There are different
// cases in which we determine to remove an unnecessary goog.require:
//
//  1. Extra requires with local-level suppressions should not be removed.
//     (@suppress {extraRequire})
//  2. Remove expression statements. ex: const x = goog.require('foo');
//  3. Remove aliased goog.require imports: const bar =
//     goog.require('foo.bar');
//  4. Remove destructured goog.require imports if the identifier(s) are
//     unused: const {x, unused} = goog.require('foo.bar');
//
// Java: ExtraRequireRemover.removeExtraRequire(Node, String)
func (e *ExtraRequireRemover) RemoveExtraRequire(call rhino.Node, require string) {
	if e.usages.Contains(require) {
		return
	}

	// If a goog.require contains a /** @suppress {extraRequire} */ local-level suppression
	// (suppression right above the import line), we will NOT prune the goog.require. This
	// suppression indicates that the goog.require is actually needed even though it does not appear
	// to be referenced in the file, so don't prune it.
	jsDoc := NodeUtil_GetBestJSDocInfo(call)
	if jsDoc != nil && jsDoc.HasSuppression("extraRequire") {
		return
	}

	// Workaround for tsickle bug affecting JSPB imports.
	// TODO: b/520110716 - Remove this once the bug is fixed.
	if ExtraRequireRemover_IsTsickleProtoCall(call) {
		return
	}

	// Case 1: remove if this is an expression statement.
	// ex: goog.require('foo');
	if call.GetParent().IsExprResult() {
		// all unaliased imports in goog.module are not pruned regardless of if they are used or
		// not
		if !call.GetGrandparent().IsModuleBody() {
			e.compiler.ReportChangeToEnclosingScope(call.GetParent())
			call.GetParent().Detach()
		}
		return
	}

	// Case 2: this is a 'default' import esque thing remove the entire statement.
	// ex: Remove `const bar = goog.require('foo.bar');`
	if call.GetParent().IsName() {
		nameDeclaration := call.GetGrandparent()
		jutil.CheckState(NodeUtil_IsNameDeclaration(call.GetGrandparent()))
		e.compiler.ReportChangeToEnclosingScope(nameDeclaration.GetParent())
		nameDeclaration.Detach()
		return
	}

	// Case 3: this is part of a destructuring import.
	//   const {unused} = goog.require('completely.unused');
	//   const {x, unused} = goog.require('foo.bar');
	// Prune destructured imports that are partially unused by pruning the
	// unused identifier(s). If all the identifiers are unused, prune the entire
	// import.
	jutil.CheckStateMsg1(call.GetParent().IsStringKey(), "%s", call.GetParent())
	objectPattern := call.GetGrandparent()

	// the 'const {unused} = goog.require('completely.unused');' case
	if objectPattern.HasOneChild() {
		nameDeclaration := objectPattern.GetGrandparent()
		jutil.CheckStateMsg1(NodeUtil_IsNameDeclaration(nameDeclaration), "%s", nameDeclaration)
		e.compiler.ReportChangeToEnclosingScope(nameDeclaration.GetParent())
		nameDeclaration.Detach()
		return
	}

	// 'const {x, unused} = goog.require('foo.bar');' in this case just remove 'unused'.
	e.compiler.ReportChangeToEnclosingScope(call.GetGrandparent())
	call.GetParent().Detach()
}

// VisitRequire: localName is the name that should be used in this file.
//
//	Require style                        | localName
//	-------------------------------------|----------
//	goog.require('foo.bar');             | foo.bar
//	var bar = goog.require('foo.bar');   | bar
//	var {qux} = goog.require('foo.bar'); | qux
//	import {qux} from 'foo.bar';         | qux
//
// namespace is the namespace being imported. May be identical to localName
// in the case of goog.require('foo.bar');
//
// Java: ExtraRequireRemover.visitRequire(String, String, Node)
func (e *ExtraRequireRemover) VisitRequire(localName string, namespace string, node rhino.Node) {
	if extraRequireRemover_DEFAULT_EXTRA_NAMESPACES.Contains(namespace) {
		return
	}
	e.requires.PutIfAbsent(localName, node)
}

// VisitGoogRequire visits local names and adds the name to `usages`.
//
// Ex: `const A = goog.require('foo');` add `A` to `usages`
//
// If it is a destructured import, add each identifier to `usages`.
//
// Ex: `const {Foo, Bar, Baz} = goog.require('foo');` add `Foo`, `Bar`, and
// `Baz` to `usages`
//
// Java: ExtraRequireRemover.visitGoogRequire(String, Node, Node)
func (e *ExtraRequireRemover) VisitGoogRequire(namespace string, googRequireCall rhino.Node, parent rhino.Node) {
	// local names
	if parent.IsName() {
		e.VisitRequire(parent.GetString(), namespace, googRequireCall)
	} else if parent.IsDestructuringLhs() && parent.GetFirstChild().IsObjectPattern() {
		// destructured imports
		if parent.GetFirstChild().HasChildren() {
			// loop through identifiers in destructured import
			for stringKey := parent.GetFirstFirstChild(); !stringKey.IsNil(); stringKey = stringKey.GetNext() {
				importName := stringKey.GetFirstChild()
				e.VisitRequire(importName.GetString(), namespace, importName)
			}
		} else {
			e.VisitRequire(namespace, namespace, googRequireCall)
		}
	} else {
		e.VisitRequire(namespace, namespace, googRequireCall)
	}
}

var extraRequireRemover_GOOG_MODULE_GET = rhino.QualifiedName_Of("goog.module.get")

// VisitCallNode extracts the namespace and check if the qualified name
// matches the callee of 'goog.require' , 'goog.requireType', or
// 'goog.forwardDeclare'. If so, we will visit those corresponding nodes and
// add those usages.
//
// Ex: Qualified name is 'goog.module.get('foo')' and callee is
// 'goog.require('foo')'
//
// Java: ExtraRequireRemover.visitCallNode(Node, Node)
func (e *ExtraRequireRemover) VisitCallNode(call rhino.Node, parent rhino.Node) {
	required := e.ExtractNamespaceIfRequire(call)
	if required != nil {
		e.VisitGoogRequire(*required, call, parent)
		return
	}

	callee := call.GetFirstChild()
	if extraRequireRemover_GOOG_MODULE_GET.Matches(callee) && call.GetSecondChild().IsStringLit() {
		e.usages.Add(call.GetSecondChild().GetString())
	}
}

// AddUsagesOfAllPrefixes: for qualified name node, we want to add all
// usages for each prefix of the name, because one of them might be
// required.
//
// For 'foo.bar.baz.qux' add usages for 'foo.bar.baz.qux', 'foo.bar.baz',
// 'foo.bar', and 'foo' because any of those might be a require that we need
// to include.
//
// Java: ExtraRequireRemover.addUsagesOfAllPrefixes(QualifiedName)
func (e *ExtraRequireRemover) AddUsagesOfAllPrefixes(qualifiedName rhino.QualifiedName) {
	e.usages.Add(qualifiedName.Join())

	for !qualifiedName.IsSimple() {
		qualifiedName = qualifiedName.GetOwner()
		e.usages.Add(qualifiedName.Join())
	}
}

// Java: ExtraRequireRemover.visitQualifiedName(Node)
func (e *ExtraRequireRemover) VisitQualifiedName(n rhino.Node) {
	jutil.CheckArgumentMsg1(n.IsQualifiedName(), "%s", n)
	e.AddUsagesOfAllPrefixes(n.GetQualifiedNameObject())
}

// MaybeAddJsDocUsages: if the node has corresponding jsdoc info then go
// through its type names and track usages
//
// Java: ExtraRequireRemover.maybeAddJsDocUsages(Node)
func (e *ExtraRequireRemover) MaybeAddJsDocUsages(n rhino.Node) {
	info := n.GetJSDocInfo()
	if info == nil {
		return
	}

	for it := info.TypeExprIter(); it.Next(); {
		for typeName := range it.Expr().GetAllTypeNames().All() {
			e.AddUsagesOfAllPrefixes(rhino.QualifiedName_Of(typeName))
		}
	}
}

// Java: ExtraRequireRemover.isTsickleProtoCall(Node)
func ExtraRequireRemover_IsTsickleProtoCall(call rhino.Node) bool {
	if !call.IsCall() || !strings.HasSuffix(*call.GetSourceFileName(), ".closure.js") {
		return false
	}
	moduleId := call.GetSecondChild()
	return !moduleId.IsNil() && moduleId.IsStringLit() && strings.HasPrefix(moduleId.GetString(), "proto.")
}
