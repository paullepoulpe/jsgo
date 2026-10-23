// Port of com/google/javascript/jscomp/lint/CheckExtraRequires.java @ closure-compiler v20261015

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
	requires *jutil.LinkedHashMap[string, *rhino.Node]

	// Adding an entry to usages indicates that the name (either a fully qualified or local name)
	// is used and can be required.  Note that since usages are name-based and not scoped, any usage
	// that shadows an unused require in that file will cause the extra require warning to be missed.
	usages *jutil.LinkedHashSet[string]

	// This is only relevant for the standalone CheckExtraRequires run. This is
	// used to restrict the linter rule only for the modules listed in this set
	requiresToRemove *jutil.LinkedHashSet[string]
}

var CheckExtraRequires_EXTRA_REQUIRE_WARNING = DiagnosticType_Disabled(
	"JSC_EXTRA_REQUIRE_WARNING",
	"extra require: ''{0}'' is never referenced in this file")

// TODO(b/130215517): This should eventually be removed and exceptions suppressed
var checkExtraRequires_DEFAULT_EXTRA_NAMESPACES = jutil.NewLinkedHashSet(
	"goog.testing.asserts", "goog.testing.jsunit")

// NewCheckExtraRequires is new CheckExtraRequires(compiler, requiresToRemove).
//
// requiresToRemove: providing a non-null set to this parameter will result
// in only removing the goog.requires that are in this set. If this is null,
// it will attempt to remove all the unnecessary requires.
//
// Java: CheckExtraRequires(AbstractCompiler, ImmutableSet<String>)
func NewCheckExtraRequires(
	compiler AbstractCompiler, requiresToRemove *jutil.LinkedHashSet[string]) *CheckExtraRequires {
	c := &CheckExtraRequires{
		requires: jutil.NewLinkedHashMap[string, *rhino.Node](),
		usages:   jutil.NewLinkedHashSet[string](),
	}
	c.compiler = compiler
	c.requiresToRemove = requiresToRemove
	return c
}

// Java: CheckExtraRequires.process(Node, Node)
func (c *CheckExtraRequires) Process(externs *rhino.Node, root *rhino.Node) {
	c.Reset()
	NodeTraversal_Traverse(c.compiler, root, c)
}

// Java: CheckExtraRequires.extractNamespace(Node, String...)
func (c *CheckExtraRequires) ExtractNamespace(call *rhino.Node, primitiveNames ...string) *string {
	callee := call.GetFirstChild()
	if !callee.IsGetProp() {
		return nil
	}
	for _, primitiveName := range primitiveNames {
		if callee.MatchesQualifiedName(primitiveName) {
			target := callee.GetNext()
			if target != nil && target.IsStringLit() {
				return jutil.StringPtr(target.GetString())
			}
		}
	}
	return nil
}

// Java: CheckExtraRequires.extractNamespaceIfRequire(Node)
func (c *CheckExtraRequires) ExtractNamespaceIfRequire(call *rhino.Node) *string {
	return c.ExtractNamespace(call, "goog.require", "goog.requireType")
}

// Java: CheckExtraRequires.extractNamespaceIfForwardDeclare(Node)
func (c *CheckExtraRequires) ExtractNamespaceIfForwardDeclare(call *rhino.Node) *string {
	return c.ExtractNamespace(call, "goog.forwardDeclare")
}

// Java: CheckExtraRequires.visit(NodeTraversal, Node, Node)
func (c *CheckExtraRequires) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
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
		c.VisitCallNode(n, parent)
	case rhino.Token_SCRIPT:
		c.VisitScriptNode()
		c.Reset()
	case rhino.Token_IMPORT:
		c.VisitImportNode(n)
	default:
	}
}

// Java: CheckExtraRequires.reset()
func (c *CheckExtraRequires) Reset() {
	c.usages.Clear()
	c.requires.Clear()
}

// Java: CheckExtraRequires.visitScriptNode()
func (c *CheckExtraRequires) VisitScriptNode() {
	// For every goog.require, check that there is a usage and warn if there is not.
	for require, call := range c.requires.All() {
		if !c.usages.Contains(require) &&
			(c.requiresToRemove == nil || c.requiresToRemove.Contains(require)) {
			c.ReportExtraRequireWarning(call, require)
		}
	}
}

// Java: CheckExtraRequires.reportExtraRequireWarning(Node, String)
func (c *CheckExtraRequires) ReportExtraRequireWarning(call *rhino.Node, require string) {
	if checkExtraRequires_DEFAULT_EXTRA_NAMESPACES.Contains(require) {
		return
	}
	jsDoc := NodeUtil_GetBestJSDocInfo(call)
	if jsDoc != nil && jsDoc.GetSuppressions().Contains("extraRequire") {
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
func (c *CheckExtraRequires) VisitRequire(localName string, node *rhino.Node) {
	c.requires.PutIfAbsent(localName, node)
}

// Java: CheckExtraRequires.visitImportNode(Node)
func (c *CheckExtraRequires) VisitImportNode(importNode *rhino.Node) {
	defaultImport := importNode.GetFirstChild()
	if defaultImport.IsName() {
		c.VisitRequire(defaultImport.GetString(), importNode)
	}
	namedImports := defaultImport.GetNext()
	if namedImports.IsImportSpecs() {
		for importSpec := namedImports.GetFirstChild(); importSpec != nil; importSpec = importSpec.GetNext() {
			c.VisitRequire(importSpec.GetLastChild().GetString(), importNode)
		}
	}
}

// Java: CheckExtraRequires.visitForwardDeclare(String, Node, Node)
func (c *CheckExtraRequires) VisitForwardDeclare(namespace string, forwardDeclareCall *rhino.Node, parent *rhino.Node) {
	c.VisitGoogRequire(namespace, forwardDeclareCall, parent)
}

// Java: CheckExtraRequires.visitGoogRequire(String, Node, Node)
func (c *CheckExtraRequires) VisitGoogRequire(namespace string, googRequireCall *rhino.Node, parent *rhino.Node) {
	if parent.IsName() {
		c.VisitRequire(parent.GetString(), googRequireCall)
	} else if parent.IsDestructuringLhs() && parent.GetFirstChild().IsObjectPattern() {
		if parent.GetFirstChild().HasChildren() {
			for stringKey := parent.GetFirstFirstChild(); stringKey != nil; stringKey = stringKey.GetNext() {
				importName := stringKey.GetFirstChild()
				if !importName.IsName() {
					// invalid reported elsewhere
					continue
				}
				c.VisitRequire(importName.GetString(), importName)
			}
		} else {
			c.VisitRequire(namespace, googRequireCall)
		}
	} else {
		c.VisitRequire(namespace, googRequireCall)
	}
}

var checkExtraRequires_GOOG_MODULE_GET = rhino.QualifiedName_Of("goog.module.get")

// Java: CheckExtraRequires.visitCallNode(Node, Node)
func (c *CheckExtraRequires) VisitCallNode(call *rhino.Node, parent *rhino.Node) {
	required := c.ExtractNamespaceIfRequire(call)
	if required != nil {
		c.VisitGoogRequire(*required, call, parent)
		return
	}
	declare := c.ExtractNamespaceIfForwardDeclare(call)
	if declare != nil {
		c.VisitForwardDeclare(*declare, call, parent)
		return
	}
	callee := call.GetFirstChild()
	if checkExtraRequires_GOOG_MODULE_GET.Matches(callee) && call.GetSecondChild().IsStringLit() {
		c.usages.Add(call.GetSecondChild().GetString())
	}
}

// Java: CheckExtraRequires.addUsagesOfAllPrefixes(QualifiedName)
func (c *CheckExtraRequires) AddUsagesOfAllPrefixes(qualifiedName rhino.QualifiedName) {
	c.usages.Add(qualifiedName.Join())
	// For "foo.bar.baz.qux" add usages for "foo.bar.baz.qux", "foo.bar.baz",
	// "foo.bar", and "foo" because any of those might be a require that
	// we need to include.
	for !qualifiedName.IsSimple() {
		qualifiedName = qualifiedName.GetOwner()
		c.usages.Add(qualifiedName.Join())
	}
}

// Java: CheckExtraRequires.visitQualifiedName(Node)
func (c *CheckExtraRequires) VisitQualifiedName(n *rhino.Node) {
	jutil.CheckArgumentMsg1(n.IsQualifiedName(), "%s", n)
	c.AddUsagesOfAllPrefixes(n.GetQualifiedNameObject())
}

// Java: CheckExtraRequires.maybeAddJsDocUsages(Node)
func (c *CheckExtraRequires) MaybeAddJsDocUsages(n *rhino.Node) {
	info := n.GetJSDocInfo()
	if info == nil {
		return
	}
	for e := range info.TypeExpressions() {
		for typeName := range e.GetAllTypeNames().All() {
			c.AddUsagesOfAllPrefixes(rhino.QualifiedName_Of(typeName))
		}
	}
}
