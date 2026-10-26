// Port of com/google/javascript/jscomp/Es6RenameTypeReferences.java @ closure-compiler v20261021

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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Es6RenameTypeReferences renames references in JSDoc.
type Es6RenameTypeReferences struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	// Map from the root node of a scope + the original variable name to a
	// new name for the variable.
	//
	// An entry in this map means we should rename a variable that was
	// originally declared in the scope with the given root and had the
	// original variable name to the new name.
	renameTable map[Es6ScopeRenameKey]string
}

// Java: Es6RenameTypeReferences(Table<Node, String, String>)
func NewEs6RenameTypeReferences(renameTable map[Es6ScopeRenameKey]string) *Es6RenameTypeReferences {
	return &Es6RenameTypeReferences{renameTable: renameTable}
}

// Java: Es6RenameTypeReferences.visit(NodeTraversal, Node, Node)
func (e *Es6RenameTypeReferences) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	// Remove this branch after module rewriting is always done
	// after type checking.
	info := n.GetJSDocInfo()
	if info != nil {
		for it := info.TypeExprIter(); it.Next(); {
			e.RenameTypeNodeRecursive(t, it.Node())
		}
	}
}

// Java: Es6RenameTypeReferences.renameTypeNodeRecursive(NodeTraversal, Node)
func (e *Es6RenameTypeReferences) RenameTypeNodeRecursive(t *NodeTraversal, n rhino.Node) {
	if n.IsStringLit() {
		e.RenameTypeReference(t, n)
	}

	for child := n.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
		e.RenameTypeNodeRecursive(t, child)
	}
}

// Java: Es6RenameTypeReferences.renameTypeReference(NodeTraversal, Node)
func (e *Es6RenameTypeReferences) RenameTypeReference(t *NodeTraversal, n rhino.Node) {
	jutil.CheckState(n.IsStringLit())
	fullName := n.GetString()
	var rootName string
	var rest *string
	endPos := jutil.StringIndexOfChar(fullName, '.')
	if endPos == -1 {
		rootName = fullName
		rest = nil
	} else {
		rootName = jutil.StringSubstring(fullName, 0, endPos)
		r := jutil.StringSubstringFrom(fullName, endPos)
		rest = &r
	}
	e.RenameReference(t, n, rootName, rest)
}

// RenameReference: t is the NodeTraversal; oldName is the root name of a
// qualified name; rest is the rest of the qualified name or null.
//
// Java: Es6RenameTypeReferences.renameReference(NodeTraversal, Node, String, String)
func (e *Es6RenameTypeReferences) RenameReference(t *NodeTraversal, n rhino.Node, oldName string, rest *string) {
	current := t.GetScope()

	// You should be wondering:
	//
	// Why are we searching up the stack of scopes here instead of just looking up the Var for
	// oldName and getting the scope from that?
	//
	// The answer is:  The Var no longer exists.
	//
	// Es6RewriteModules (the only client for this class) actually deletes the module-level
	// declaration for the old name but still uses the scope rooted at the MODULE_BODY (for
	// goog.module() files), SCRIPT, or function body (for wrapped goog.loadModule() goog.module()s)
	// in the rename map.
	//
	// If we're in some local scope where there's a shadowing local variable, we don't want to do
	// the renaming.
	//
	// Otherwise, we need to use the MODULE_BODY or SCRIPT or goog.loadModule body scope to look up
	// the new name.
	// - bradfordcsmith@google.com
	for current != nil {
		currentRootNode := current.GetRootNode()
		newName, ok := e.renameTable[Es6ScopeRenameKey{Root: currentRootNode, Name: oldName}]
		if ok {
			jutil.CheckStateMsg1(
				currentRootNode.IsModuleBody() ||
					currentRootNode.IsScript() ||
					NodeUtil_IsFunctionBlock(currentRootNode),
				"Not a MODULE_BODY or SCRIPT or goog.loadModule root: %s",
				currentRootNode)
			var newFullName string
			if rest == nil {
				newFullName = newName
			} else {
				newFullName = newName + *rest
			}
			n.SetString(newFullName)
			return
		} else if current.HasOwnSlot(oldName) {
			// This is a reference to some local variable whose name shadows the module-scope variable
			// we actually want to rename.
			return
		} else {
			current = current.GetParent()
		}
	}
}
