// Port of com/google/javascript/refactoring/ScriptMetadata.java @ closure-compiler v20261021

/*
 * Copyright 2014 The Closure Compiler Authors.
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

package refactoring

import (
	"strings"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jscomp"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ScriptMetadata is a summary of a script for use during fixes.
//
// Instances are mutable so that changes to a script caused by one fix can be
// reflected during subsequent fixes.
type ScriptMetadata struct {
	script                 rhino.Node
	supportsRequireAliases bool
	requireAliases         map[string]string
	namesInUse             map[string]struct{}
}

// Java: ScriptMetadata.create(Node, NodeMetadata)
func ScriptMetadata_CreateFromNodeMetadata(script rhino.Node, metadata *NodeMetadata) *ScriptMetadata {
	return ScriptMetadata_Create(script, metadata.GetCompiler())
}

// Java: ScriptMetadata.create(Node, AbstractCompiler)
func ScriptMetadata_Create(script rhino.Node, compiler jscomp.AbstractCompiler) *ScriptMetadata {
	jutil.CheckArgument(script.IsScript())

	toFill := newScriptMetadata(script)
	jscomp.NodeTraversal_Traverse(compiler, script, &scriptMetadata_FillerCallback{toFill: toFill})

	// Even if the file was not found to be a module, the existence of other aliases means aliases
	// are supported.
	if len(toFill.requireAliases) != 0 {
		toFill.supportsRequireAliases = true
	}
	return toFill
}

func newScriptMetadata(script rhino.Node) *ScriptMetadata {
	return &ScriptMetadata{
		script:         script,
		requireAliases: make(map[string]string),
		namesInUse:     make(map[string]struct{}),
	}
}

// Java: ScriptMetadata.getScript()
func (m *ScriptMetadata) GetScript() rhino.Node {
	return m.script
}

// SupportsRequireAliases reports whether requires added to this script should
// be given local aliases.
//
// Java: ScriptMetadata.supportsRequireAliases()
func (m *ScriptMetadata) SupportsRequireAliases() bool {
	return m.supportsRequireAliases
}

// UsesName reports whether `name` is in use within this script (including both
// executable and JsDoc references).
//
// Java: ScriptMetadata.usesName(String)
func (m *ScriptMetadata) UsesName(name string) bool {
	_, ok := m.namesInUse[name]
	return ok
}

func (m *ScriptMetadata) lookupAlias(namespace string) (string, bool) {
	alias, ok := m.requireAliases[namespace]
	return alias, ok
}

// GetAlias returns the local alias for `namespace`, or nil if not required or
// unaliased.
//
// Java: ScriptMetadata.getAlias(String)
func (m *ScriptMetadata) GetAlias(namespace string) *string {
	alias, ok := m.requireAliases[namespace]
	if !ok {
		return nil
	}
	return &alias
}

// AddAlias records that `namespace` has been required with local alias `alias`.
//
// Java: ScriptMetadata.addAlias(String, String)
func (m *ScriptMetadata) AddAlias(namespace string, alias string) {
	jutil.CheckState(m.supportsRequireAliases)
	m.namesInUse[alias] = struct{}{}
	m.requireAliases[namespace] = alias
}

var (
	scriptMetadata_FillerCallback_GOOG_MODULE         = rhino.QualifiedName_Of("goog.module")
	scriptMetadata_FillerCallback_GOOG_REQUIRE        = rhino.QualifiedName_Of("goog.require")
	scriptMetadata_FillerCallback_GOOG_REQUIRETYPE    = rhino.QualifiedName_Of("goog.requireType")
	scriptMetadata_FillerCallback_GOOG_FORWARDDECLARE = rhino.QualifiedName_Of("goog.forwardDeclare")
)

type scriptMetadata_FillerCallback struct {
	jscomp.NodeTraversal_AbstractPostOrderCallbackBase
	toFill *ScriptMetadata
}

func (c *scriptMetadata_FillerCallback) Visit(nodeTraversal *jscomp.NodeTraversal, n rhino.Node, parent rhino.Node) {
	c.updateSupportsRequireAlias(n)
	c.maybeCollectRequirelikeAlias(n)
	c.maybeCollectNames(n)
}

func (c *scriptMetadata_FillerCallback) updateSupportsRequireAlias(n rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_MODULE_BODY:
	case rhino.Token_CALL:
		call := ast.MustCall(n)
		arg, ok := call.OnlyArg()
		if !ok || !scriptMetadata_FillerCallback_GOOG_MODULE.Matches(call.Callee().Node()) {
			return
		}
		if _, isStr := arg.AsStringLit(); !isStr {
			return
		}
	default:
		return
	}
	c.toFill.supportsRequireAliases = true
}

func (c *scriptMetadata_FillerCallback) maybeCollectRequirelikeAlias(n rhino.Node) {
	decl, ok := ast.AsNameDeclaration(n)
	if !ok {
		return
	}
	firstDecl, ok := decl.FirstDeclarator()
	if !ok {
		return
	}
	// TODO(b/139953612): respect destructured goog.requires
	aliasName, ok := ast.AsName(firstDecl.Node())
	if !ok {
		return
	}
	initExpr, ok := aliasName.Init()
	if !ok {
		return
	}
	namespace, ok := requirelikeNamespace(initExpr.Node())
	if !ok {
		return
	}
	alias := aliasName.Ident()
	if originalName := aliasName.Node().GetOriginalName(); originalName != nil {
		alias = *originalName
	}
	c.toFill.requireAliases[namespace] = alias
}

func requirelikeNamespace(node rhino.Node) (string, bool) {
	call, ok := ast.AsCall(node)
	if !ok {
		return "", false
	}
	arg, ok := call.OnlyArg()
	if !ok {
		return "", false
	}
	strLit, ok := arg.AsStringLit()
	if !ok {
		return "", false
	}
	callee := call.Callee().Node()
	if !scriptMetadata_FillerCallback_GOOG_REQUIRE.Matches(callee) &&
		!scriptMetadata_FillerCallback_GOOG_REQUIRETYPE.Matches(callee) &&
		!scriptMetadata_FillerCallback_GOOG_FORWARDDECLARE.Matches(callee) {
		return "", false
	}
	return strLit.Value(), true
}

func (c *scriptMetadata_FillerCallback) maybeCollectNames(n rhino.Node) {
	if name, ok := ast.AsName(n); ok {
		c.toFill.namesInUse[name.Ident()] = struct{}{}
	}
	if jsdoc := n.GetJSDocInfo(); jsdoc != nil {
		for _, expr := range jsdoc.GetTypeExpressions() {
			for type_ := range expr.GetAllTypeNames().All() {
				before, _, _ := strings.Cut(type_, ".")
				c.toFill.namesInUse[before] = struct{}{}
			}
		}
	}
}
