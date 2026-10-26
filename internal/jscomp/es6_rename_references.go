// Port of com/google/javascript/jscomp/Es6RenameReferences.java @ closure-compiler v20261021

/*
 * Copyright 2015 The Closure Compiler Authors.
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

type Es6ScopeRenameKey struct {
	Root rhino.Node
	Name string
}

// Es6RenameReferences renames references in code.
type Es6RenameReferences struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	// Map from the root node of a scope + the original variable name to a
	// new name for the variable.
	//
	// An entry in this map means we should rename a variable that was
	// originally declared in the scope with the given root and had the
	// original variable name to the new name.
	renameTable map[Es6ScopeRenameKey]string
}

// Java: Es6RenameReferences(Table<Node, String, String>)
func NewEs6RenameReferences(renameTable map[Es6ScopeRenameKey]string) *Es6RenameReferences {
	return &Es6RenameReferences{renameTable: renameTable}
}

// Java: Es6RenameReferences.visit(NodeTraversal, Node, Node)
func (e *Es6RenameReferences) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if NodeUtil_IsReferenceName(n) {
		e.RenameReference(t, n, n.GetString())
	}
}

// RenameReference: t is the NodeTraversal; oldName is the root name of a
// qualified name.
//
// Java: Es6RenameReferences.renameReference(NodeTraversal, Node, String)
func (e *Es6RenameReferences) RenameReference(t *NodeTraversal, n rhino.Node, oldName string) {
	current := t.GetScope()
	// You should be wondering:
	//
	// Why are we searching up the stack of scopes here instead of just looking up the Var for
	// oldName and getting the scope from that?
	//
	// The answer is:
	//
	// When Es6RewriteBlockScopedDeclaration uses this class, it has already modified the scopes so
	// that `oldName` is no longer declared in its original scope.
	//
	// I don't know if the other users of this class are also modifying the scopes in advance.
	// - bradfordcsmith@google.com
	for current != nil {
		newName, ok := e.renameTable[Es6ScopeRenameKey{Root: current.GetRootNode(), Name: oldName}]
		if ok {
			n.SetString(newName)
			t.ReportCodeChange()
			return
		} else if current.HasOwnSlot(oldName) {
			return
		} else {
			current = current.GetParent()
		}
	}
}
