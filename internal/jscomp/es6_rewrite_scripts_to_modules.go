// Port of com/google/javascript/jscomp/Es6RewriteScriptsToModules.java @ closure-compiler v20261015

/*
 * Copyright 2017 The Closure Compiler Authors.
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

// Es6RewriteScriptsToModules rewrites a script which was imported as a
// module into an ES6 module.
type Es6RewriteScriptsToModules struct {
	NodeTraversal_AbstractPreOrderCallbackBase

	compiler AbstractCompiler
}

// NewEs6RewriteScriptsToModules creates a new Es6RewriteModules instance
// which can be used to rewrite ES6 modules to a concatenable form.
//
// Java: Es6RewriteScriptsToModules(AbstractCompiler)
func NewEs6RewriteScriptsToModules(compiler AbstractCompiler) *Es6RewriteScriptsToModules {
	e := &Es6RewriteScriptsToModules{}
	e.compiler = compiler
	return e
}

// ForceToEs6Module forces rewriting of a script into an ES6 module, such as
// for imported files that contain no "import" or "export" statements.
//
// Java: Es6RewriteScriptsToModules.forceToEs6Module(Node)
func (e *Es6RewriteScriptsToModules) ForceToEs6Module(root rhino.Node) {
	if Es6RewriteModules_IsEs6ModuleRoot(root) {
		return
	}
	moduleNode := rhino.NewNode(rhino.Token_MODULE_BODY).Srcref(root)
	moduleNode.AddChildrenToBack(root.RemoveChildren())
	root.AddChildToBack(moduleNode)
	e.compiler.ReportChangeToChangeScope(root)
}

// Java: Es6RewriteScriptsToModules.process(Node, Node)
func (e *Es6RewriteScriptsToModules) Process(externs rhino.Node, root rhino.Node) {
	for file := root.GetFirstChild(); !file.IsNil(); file = file.GetNext() {
		jutil.CheckArgument(file.IsScript())
		NodeTraversal_Traverse(e.compiler, file, e)
	}
}

// Java: Es6RewriteScriptsToModules.shouldTraverse(NodeTraversal, Node, Node)
func (e *Es6RewriteScriptsToModules) ShouldTraverse(nodeTraversal *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	if n.IsScript() {
		moduleType := e.compiler.GetInput(n.GetInputId()).GetJsModuleType()

		if moduleType == CompilerInput_ModuleType_IMPORTED_SCRIPT {
			e.ForceToEs6Module(n)
		}
	}

	return false
}
