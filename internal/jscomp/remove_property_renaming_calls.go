// Port of com/google/javascript/jscomp/RemovePropertyRenamingCalls.java @ closure-compiler v20261015

/*
 * Copyright 2023 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// RemovePropertyRenamingCalls replaces all calls to property-renaming
// functions like JSCompiler_renameProperty with the string literal property
// name.
//
// Intended for use when RenameProperties isn't running, since this
// functionality is also baked into RenameProperties.
type RemovePropertyRenamingCalls struct {
	compiler AbstractCompiler
}

// Java: RemovePropertyRenamingCalls(AbstractCompiler)
func NewRemovePropertyRenamingCalls(compiler AbstractCompiler) *RemovePropertyRenamingCalls {
	return &RemovePropertyRenamingCalls{compiler: compiler}
}

// Java: RemovePropertyRenamingCalls.process(Node, Node)
func (r *RemovePropertyRenamingCalls) Process(externs rhino.Node, root rhino.Node) {
	NodeTraversal_Traverse(r.compiler, root, &RemovePropertyRenamingCalls_RemoveCompilerPropertyRenamingCallback{outer: r})
}

// RemovePropertyRenamingCalls_RemoveCompilerPropertyRenamingCallback is the
// inner class RemoveCompilerPropertyRenamingCallback.
type RemovePropertyRenamingCalls_RemoveCompilerPropertyRenamingCallback struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer *RemovePropertyRenamingCalls
}

// Java: RemovePropertyRenamingCalls.RemoveCompilerPropertyRenamingCallback.visit(NodeTraversal, Node, Node)
func (c *RemovePropertyRenamingCalls_RemoveCompilerPropertyRenamingCallback) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if n.IsCall() {
		fn := n.GetFirstChild()
		if c.outer.compiler.GetCodingConvention().IsPropertyRenameFunction(fn) {
			c.ReplacePropertyRenamingCall(n)
		}
	}
}

// ReplacePropertyRenamingCall replaces JSCompiler_renameProperty("bar", obj)
// with "bar".
//
// Java: RemovePropertyRenamingCalls.RemoveCompilerPropertyRenamingCallback.replacePropertyRenamingCall(Node)
func (c *RemovePropertyRenamingCalls_RemoveCompilerPropertyRenamingCallback) ReplacePropertyRenamingCall(callNode rhino.Node) {
	propName := NodeUtil_GetArgumentForCallOrNew(callNode, 0)
	if !propName.IsNil() {
		callNode.ReplaceWith(propName.Detach())
		c.outer.compiler.ReportChangeToEnclosingScope(propName)
		return
	}
}
