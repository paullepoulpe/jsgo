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

	"github.com/paullepoulpe/jsgo/internal/jscomp"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ScriptMetadata is a summary of a script for use during fixes.
//
// Instances are mutable so that changes to a script caused by one fix cen be
// reflected during subseuqent fixes.
type ScriptMetadata struct {
	script                 rhino.Node
	supportsRequireAliases bool
	requireAliases         map[string]string
	namesInUse             map[string]bool
}

// Java: ScriptMetadata.create(Node, NodeMetadata)
func ScriptMetadata_CreateFromNodeMetadata(script rhino.Node, metadata *NodeMetadata) *ScriptMetadata {
	return ScriptMetadata_Create(script, metadata.GetCompiler())
}

// Java: ScriptMetadata.create(Node, AbstractCompiler)
func ScriptMetadata_Create(script rhino.Node, compiler jscomp.AbstractCompiler) *ScriptMetadata {
	jutil.CheckArgument(script.IsScript())

	toFill := newScriptMetadata(script)
	jscomp.NodeTraversal_Traverse(compiler, script, newScriptMetadata_FillerCallback(toFill))

	// Even if the file was not found to be a module, the existence of other aliases means aliases
	// are supported.
	if len(toFill.requireAliases) != 0 {
		toFill.supportsRequireAliases = true
	}
	return toFill
}

// Java: ScriptMetadata(Node)
func newScriptMetadata(script rhino.Node) *ScriptMetadata {
	return &ScriptMetadata{
		script:         script,
		requireAliases: make(map[string]string),
		namesInUse:     make(map[string]bool),
	}
}

// Java: ScriptMetadata.getScript()
func (m *ScriptMetadata) GetScript() rhino.Node {
	return m.script
}

// SupportsRequireAliases: When requires are added to this script, should
// they be given local aliases?
//
// Java: ScriptMetadata.supportsRequireAliases()
func (m *ScriptMetadata) SupportsRequireAliases() bool {
	return m.supportsRequireAliases
}

// UsesName: Is `name` in use within this script?
//
// This includes both executable and JsDoc references.
//
// Java: ScriptMetadata.usesName(String)
func (m *ScriptMetadata) UsesName(name string) bool {
	return m.namesInUse[name]
}

// GetAlias returns the local alias for `namespace`.
//
// If `namespace` is not required or has no alias, returns `null`. In
// particular, destructured aliases will not be available.
//
// Java: ScriptMetadata.getAlias(String)
func (m *ScriptMetadata) GetAlias(namespace string) *string {
	alias, ok := m.requireAliases[namespace]
	if !ok {
		return nil
	}
	return &alias
}

// AddAlias records that `namespace` has ben required with local alias
// `alias`.
//
// As a corallary, any alias is also a name in use within this script.
//
// Java: ScriptMetadata.addAlias(String, String)
func (m *ScriptMetadata) AddAlias(namespace string, alias string) {
	jutil.CheckState(m.supportsRequireAliases)
	m.namesInUse[alias] = true
	m.requireAliases[namespace] = alias
}

var scriptMetadata_FillerCallback_GOOG_MODULE = rhino.QualifiedName_Of("goog.module")
var scriptMetadata_FillerCallback_GOOG_REQUIRE = rhino.QualifiedName_Of("goog.require")
var scriptMetadata_FillerCallback_GOOG_REQUIRETYPE = rhino.QualifiedName_Of("goog.requireType")
var scriptMetadata_FillerCallback_GOOG_FORWARDDECLARE = rhino.QualifiedName_Of("goog.forwardDeclare")

// Java: private static class ScriptMetadata.FillerCallback extends AbstractPostOrderCallback
type scriptMetadata_FillerCallback struct {
	jscomp.NodeTraversal_AbstractPostOrderCallbackBase

	toFill *ScriptMetadata
}

// Java: ScriptMetadata.FillerCallback(ScriptMetadata)
func newScriptMetadata_FillerCallback(toFill *ScriptMetadata) *scriptMetadata_FillerCallback {
	return &scriptMetadata_FillerCallback{toFill: toFill}
}

// Java: ScriptMetadata.FillerCallback.visit(NodeTraversal, Node, Node)
func (c *scriptMetadata_FillerCallback) Visit(nodeTraversal *jscomp.NodeTraversal, n rhino.Node, parent rhino.Node) {
	c.updateSupportsRequireAlias(n)
	c.maybeCollectRequirelikeAlias(n)
	c.maybeCollectNames(n)
}

// Java: ScriptMetadata.FillerCallback.updateSupportsRequireAlias(Node)
func (c *scriptMetadata_FillerCallback) updateSupportsRequireAlias(n rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_MODULE_BODY:
	case rhino.Token_CALL:
		if n.GetChildCount() != 2 ||
			!scriptMetadata_FillerCallback_GOOG_MODULE.Matches(n.GetFirstChild()) ||
			!n.GetSecondChild().IsStringLit() {
			return
		}
	default:
		return
	}

	c.toFill.supportsRequireAliases = true
}

// Java: ScriptMetadata.FillerCallback.maybeCollectRequirelikeAlias(Node)
func (c *scriptMetadata_FillerCallback) maybeCollectRequirelikeAlias(n rhino.Node) {
	if !jscomp.NodeUtil_IsNameDeclaration(n) {
		return
	}

	aliasNode := n.GetFirstChild()

	// TODO(b/139953612): respect destructured goog.requires
	if aliasNode.IsDestructuringLhs() {
		return
	}

	requirelikeCall := aliasNode.GetFirstChild()
	if !c.isRequirelikeCall(requirelikeCall) {
		return
	}

	var alias string
	if originalName := aliasNode.GetOriginalName(); originalName != nil {
		alias = *originalName
	} else {
		alias = aliasNode.GetString()
	}

	namespace := requirelikeCall.GetSecondChild().GetString()
	c.toFill.requireAliases[namespace] = alias
}

// Java: ScriptMetadata.FillerCallback.isRequirelikeCall(Node)
func (c *scriptMetadata_FillerCallback) isRequirelikeCall(call rhino.Node) bool {
	if call.IsNil() ||
		!call.IsCall() ||
		call.GetChildCount() != 2 ||
		!call.GetSecondChild().IsStringLit() {
		return false
	}

	callee := call.GetFirstChild()
	return scriptMetadata_FillerCallback_GOOG_REQUIRE.Matches(callee) ||
		scriptMetadata_FillerCallback_GOOG_REQUIRETYPE.Matches(callee) ||
		scriptMetadata_FillerCallback_GOOG_FORWARDDECLARE.Matches(callee)
}

// Java: ScriptMetadata.FillerCallback.maybeCollectNames(Node)
func (c *scriptMetadata_FillerCallback) maybeCollectNames(n rhino.Node) {
	if n.IsName() {
		c.toFill.namesInUse[n.GetString()] = true
	}

	jsdoc := n.GetJSDocInfo()
	if jsdoc != nil {
		for _, expr := range jsdoc.GetTypeExpressions() {
			for type_ := range expr.GetAllTypeNames().All() {
				c.toFill.namesInUse[scriptMetadata_firstDotSplit(type_)] = true
			}
		}
	}
}

// scriptMetadata_firstDotSplit is DOT_SPLITTER.split(type).iterator().next()
// with DOT_SPLITTER = Splitter.on('.'): the text before the first '.'.
func scriptMetadata_firstDotSplit(s string) string {
	before, _, _ := strings.Cut(s, ".")
	return before
}
