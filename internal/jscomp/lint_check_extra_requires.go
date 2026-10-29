// Port of com/google/javascript/jscomp/lint/CheckExtraRequires.java @ closure-compiler v20261021

/*
 * Copyright 2008 The Closure Compiler Authors.
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

// CheckExtraRequires walks the AST looking for usages of qualified names, and
// 'goog.require's of those names. Then, reconciles the two lists, and reports
// warning for any unnecessary require statements.
type CheckExtraRequires struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler AbstractCompiler

	// Keys are the local name of a required namespace. Values are the goog.require CALL node.
	// Accompanied by requiresOrder to preserve insertion order for deterministic warning output.
	requires      map[string]rhino.Node
	requiresOrder []string

	// Adding an entry to usages indicates that the name (either a fully qualified or local name)
	// is used and can be required. Note that since usages are name-based and not scoped, any usage
	// that shadows an unused require in that file will cause the extra require warning to be missed.
	usages map[string]bool

	// This is only relevant for the standalone CheckExtraRequires run. This is
	// used to restrict the linter rule only for the modules listed in this set
	requiresToRemove *jutil.LinkedHashSet[string]
}

var CheckExtraRequires_EXTRA_REQUIRE_WARNING = DiagnosticType_Disabled(
	"JSC_EXTRA_REQUIRE_WARNING",
	"extra require: ''{0}'' is never referenced in this file")

// TODO(b/130215517): This should eventually be removed and exceptions suppressed
func checkExtraRequires_isDefaultExtraNamespace(ns string) bool {
	return ns == "goog.testing.asserts" || ns == "goog.testing.jsunit"
}

// NewCheckExtraRequires is new CheckExtraRequires(compiler, requiresToRemove).
//
// requiresToRemove: providing a non-null set to this parameter will result
// in only removing the goog.requires that are in this set. If this is null,
// it will attempt to remove all the unnecessary requires.
//
// Java: CheckExtraRequires(AbstractCompiler, ImmutableSet<String>)
func NewCheckExtraRequires(
	compiler AbstractCompiler, requiresToRemove *jutil.LinkedHashSet[string]) *CheckExtraRequires {
	return &CheckExtraRequires{
		compiler:         compiler,
		requires:         make(map[string]rhino.Node),
		usages:           make(map[string]bool),
		requiresToRemove: requiresToRemove,
	}
}

func (c *CheckExtraRequires) cloneForParallel(rc *Compiler) NodeTraversal_Callback {
	return NewCheckExtraRequires(rc, c.requiresToRemove)
}

// Java: CheckExtraRequires.process(Node, Node)
func (c *CheckExtraRequires) Process(externs rhino.Node, root rhino.Node) {
	c.Reset()
	NodeTraversal_Traverse(c.compiler, root, c)
}

// Java: CheckExtraRequires.extractNamespace(Node, String...)
func (c *CheckExtraRequires) ExtractNamespace(call ast.Call, primitiveNames ...string) (string, bool) {
	callee := call.Callee().Node()
	if !callee.IsGetProp() {
		return "", false
	}
	firstArg, hasArg := call.FirstArg()
	if !hasArg {
		return "", false
	}
	strLit, isStr := ast.AsStringLit(firstArg.Node())
	if !isStr {
		return "", false
	}
	for _, primitiveName := range primitiveNames {
		if callee.MatchesQualifiedName(primitiveName) {
			return strLit.Value(), true
		}
	}
	return "", false
}

// Java: CheckExtraRequires.extractNamespaceIfRequire(Node)
func (c *CheckExtraRequires) ExtractNamespaceIfRequire(call ast.Call) (string, bool) {
	return c.ExtractNamespace(call, "goog.require", "goog.requireType")
}

// Java: CheckExtraRequires.extractNamespaceIfForwardDeclare(Node)
func (c *CheckExtraRequires) ExtractNamespaceIfForwardDeclare(call ast.Call) (string, bool) {
	return c.ExtractNamespace(call, "goog.forwardDeclare")
}

// Java: CheckExtraRequires.visit(NodeTraversal, Node, Node)
func (c *CheckExtraRequires) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	c.MaybeAddJsDocUsages(n)
	switch n.GetToken() {
	case rhino.Token_NAME:
		if !NodeUtil_IsLValue(n) && !parent.IsGetProp() && !parent.IsImportSpec() {
			c.VisitQualifiedName(n)
		}
	case rhino.Token_GETPROP:
		// If parent is a GETPROP, they will handle all the usages.
		if !parent.IsGetProp() && n.IsQualifiedName() {
			c.VisitQualifiedName(n)
		}
	case rhino.Token_CALL:
		c.VisitCallNode(ast.MustCall(n), parent)
	case rhino.Token_SCRIPT:
		c.VisitScriptNode()
		c.Reset()
	case rhino.Token_IMPORT:
		c.VisitImportNode(ast.MustImport(n))
	default:
	}
}

// Java: CheckExtraRequires.reset()
func (c *CheckExtraRequires) Reset() {
	clear(c.usages)
	clear(c.requires)
	c.requiresOrder = c.requiresOrder[:0]
}

// Java: CheckExtraRequires.visitScriptNode()
func (c *CheckExtraRequires) VisitScriptNode() {
	// For every goog.require, check that there is a usage and warn if there is not.
	for _, require := range c.requiresOrder {
		if !c.usages[require] &&
			(c.requiresToRemove == nil || c.requiresToRemove.Contains(require)) {
			c.ReportExtraRequireWarning(c.requires[require], require)
		}
	}
}

// Java: CheckExtraRequires.reportExtraRequireWarning(Node, String)
func (c *CheckExtraRequires) ReportExtraRequireWarning(call rhino.Node, require string) {
	if checkExtraRequires_isDefaultExtraNamespace(require) {
		return
	}
	jsDoc := NodeUtil_GetBestJSDocInfo(call)
	if jsDoc != nil && jsDoc.HasSuppression("extraRequire") {
		// There is a @suppress {extraRequire} on the call node or its enclosing statement.
		// This is one of the acceptable places for a @suppress, per
		// https://github.com/google/closure-compiler/wiki/@suppress-annotations
		return
	}
	c.compiler.Report(JSError_MakeWithNode(call, CheckExtraRequires_EXTRA_REQUIRE_WARNING, require))
}

// VisitRequire records a require; localName is the name that should be used
// in this file.
//
//	Require style                        | localName
//	-------------------------------------|----------
//	goog.require('foo.bar');             | foo.bar
//	var bar = goog.require('foo.bar');   | bar
//	var {qux} = goog.require('foo.bar'); | qux
//	import {qux} from 'foo.bar';         | qux
//
// Java: CheckExtraRequires.visitRequire(String, Node)
func (c *CheckExtraRequires) VisitRequire(localName string, node rhino.Node) {
	if _, exists := c.requires[localName]; !exists {
		c.requires[localName] = node
		c.requiresOrder = append(c.requiresOrder, localName)
	}
}

// Java: CheckExtraRequires.visitImportNode(Node)
func (c *CheckExtraRequires) VisitImportNode(imp ast.Import) {
	if defaultImport, ok := imp.Default(); ok {
		c.VisitRequire(defaultImport.Ident(), imp.Node())
	}
	if bindings, ok := imp.Bindings(); ok {
		if specs, ok := bindings.AsImportSpecs(); ok {
			for spec := range specs.Specs() {
				c.VisitRequire(spec.Local().Ident(), imp.Node())
			}
		}
	}
}

// Java: CheckExtraRequires.visitForwardDeclare(String, Node, Node)
func (c *CheckExtraRequires) VisitForwardDeclare(namespace string, forwardDeclareCall rhino.Node, parent rhino.Node) {
	c.VisitGoogRequire(namespace, forwardDeclareCall, parent)
}

// Java: CheckExtraRequires.visitGoogRequire(String, Node, Node)
func (c *CheckExtraRequires) VisitGoogRequire(namespace string, googRequireCall rhino.Node, parent rhino.Node) {
	if parentName, ok := ast.AsName(parent); ok {
		c.VisitRequire(parentName.Ident(), googRequireCall)
	} else if lhs, ok := ast.AsDestructuringLhs(parent); ok {
		if objPat, ok := lhs.Pattern().AsObjectPattern(); ok {
			if _, hasProps := objPat.FirstProp(); hasProps {
				for prop := range objPat.Props() {
					importName := prop.Node().GetFirstChild()
					if nameNode, ok := ast.AsName(importName); ok {
						c.VisitRequire(nameNode.Ident(), nameNode.Node())
					}
				}
			} else {
				c.VisitRequire(namespace, googRequireCall)
			}
			return
		}
		c.VisitRequire(namespace, googRequireCall)
	} else {
		c.VisitRequire(namespace, googRequireCall)
	}
}

var checkExtraRequires_GOOG_MODULE_GET = rhino.QualifiedName_Of("goog.module.get")

// Java: CheckExtraRequires.visitCallNode(Node, Node)
func (c *CheckExtraRequires) VisitCallNode(call ast.Call, parent rhino.Node) {
	if required, ok := c.ExtractNamespaceIfRequire(call); ok {
		c.VisitGoogRequire(required, call.Node(), parent)
		return
	}
	if declare, ok := c.ExtractNamespaceIfForwardDeclare(call); ok {
		c.VisitForwardDeclare(declare, call.Node(), parent)
		return
	}
	if checkExtraRequires_GOOG_MODULE_GET.Matches(call.Callee().Node()) {
		if firstArg, hasArg := call.FirstArg(); hasArg {
			if strLit, ok := ast.AsStringLit(firstArg.Node()); ok {
				c.usages[strLit.Value()] = true
			}
		}
	}
}

// Java: CheckExtraRequires.addUsagesOfAllPrefixes(QualifiedName)
func (c *CheckExtraRequires) AddUsagesOfAllPrefixes(qualifiedName rhino.QualifiedName) {
	c.usages[qualifiedName.Join()] = true
	// For "foo.bar.baz.qux" add usages for "foo.bar.baz.qux", "foo.bar.baz",
	// "foo.bar", and "foo" because any of those might be a require that
	// we need to include.
	for !qualifiedName.IsSimple() {
		qualifiedName = qualifiedName.GetOwner()
		c.usages[qualifiedName.Join()] = true
	}
}

// Java: CheckExtraRequires.visitQualifiedName(Node)
func (c *CheckExtraRequires) VisitQualifiedName(n rhino.Node) {
	jutil.CheckArgumentMsg1(n.IsQualifiedName(), "%s", n)
	c.AddUsagesOfAllPrefixes(n.GetQualifiedNameObject())
}

// Java: CheckExtraRequires.maybeAddJsDocUsages(Node)
func (c *CheckExtraRequires) MaybeAddJsDocUsages(n rhino.Node) {
	info := n.GetJSDocInfo()
	if info == nil {
		return
	}
	for it := info.TypeExprIter(); it.Next(); {
		for typeName := range it.Expr().GetAllTypeNames().All() {
			c.AddUsagesOfAllPrefixes(rhino.QualifiedName_Of(typeName))
		}
	}
}
