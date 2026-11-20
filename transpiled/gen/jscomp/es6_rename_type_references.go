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
// Port of com/google/javascript/jscomp/Es6RenameTypeReferences.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.Es6RenameTypeReferences
type Es6RenameTypeReferences struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	renameTable j2gort.Ext_com_google_common_collect_Table `gson:"renameTable"`
}

func (this_ *Es6RenameTypeReferences) BindSelf_(s any) {
	_ = s
}
func (this_ *Es6RenameTypeReferences) IsJavaNil_() bool { return this_ == nil }
func (*Es6RenameTypeReferences) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Es6RenameTypeReferences", "Es6RenameTypeReferences"
}
func (this_ *Es6RenameTypeReferences) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*Es6RenameTypeReferences) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Es6RenameTypeReferences", "com.google.javascript.jscomp.Es6RenameTypeReferences", func(v_ any) bool {
		return j2gort.InstanceOf[*Es6RenameTypeReferences](v_)
	})
}
func (*Es6RenameTypeReferences) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {}
func (*Es6RenameTypeReferences) Implements_jscomp_NodeTraversal_Callback()                  {}

func (this_ *Es6RenameTypeReferences) initFields_Es6RenameTypeReferences() {
	_ = this_
}

// Java constructor: Es6RenameTypeReferences(com.google.common.collect.Table<com.google.javascript.rhino.Node,java.lang.String,java.lang.String> renameTable)
func NewEs6RenameTypeReferences(renameTable j2gort.Ext_com_google_common_collect_Table) *Es6RenameTypeReferences {
	this_ := &Es6RenameTypeReferences{}
	InitEs6RenameTypeReferences(this_, renameTable)
	return this_
}

func InitEs6RenameTypeReferences(this_ *Es6RenameTypeReferences, renameTable j2gort.Ext_com_google_common_collect_Table) {
	_ = renameTable
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_Es6RenameTypeReferences()
	this_.renameTable = renameTable
}

// Java method: Es6RenameTypeReferences.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *Es6RenameTypeReferences) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	var info *pkg_rhino.JSDocInfo = n.GetJSDocInfo()
	_ = info
	if info != nil {
		{
			var i int32 = int32(0)
			_ = i
			var count int32 = info.GetTypeNodeCount()
			_ = count
			for ; i < count; i = i + 1 {
				this_.priv_Es6RenameTypeReferences_renameTypeNodeRecursive_2574(t, info.GetTypeNodeAt(i))
			}
		}
	}
}

// Java method: Es6RenameTypeReferences.renameTypeNodeRecursive(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *Es6RenameTypeReferences) priv_Es6RenameTypeReferences_renameTypeNodeRecursive_2574(t *NodeTraversal, n pkg_rhino.Node) {
	_ = t
	_ = n
	if n.IsStringLit() {
		this_.priv_Es6RenameTypeReferences_renameTypeReference_2575(t, n)
	}
	{
		var child pkg_rhino.Node = n.GetFirstChild()
		_ = child
		for ; child != 0; child = child.GetNext() {
			this_.priv_Es6RenameTypeReferences_renameTypeNodeRecursive_2574(t, child)
		}
	}
}

// Java method: Es6RenameTypeReferences.renameTypeReference(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *Es6RenameTypeReferences) priv_Es6RenameTypeReferences_renameTypeReference_2575(t *NodeTraversal, n pkg_rhino.Node) {
	_ = t
	_ = n
	if !(n.IsStringLit()) {
		j2gort.CheckState(false)
	}
	var fullName string = n.GetString()
	_ = fullName
	var rootName string = j2gort.NullString
	_ = rootName
	var rest string = j2gort.NullString
	_ = rest
	var endPos int32 = j2gort.JStrIndexOfCp(fullName, int32(uint16(46)))
	_ = endPos
	if endPos == int32(-1) {
		rootName = fullName
		rest = j2gort.NullString
	} else {
		rootName = j2gort.JStrSub(fullName, int32(0), endPos)
		rest = j2gort.JStrSubFrom(fullName, endPos)
	}
	this_.priv_Es6RenameTypeReferences_renameReference_2576(t, n, rootName, rest)
}

// Java method: Es6RenameTypeReferences.renameReference(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, java.lang.String oldName, java.lang.@org.jspecify.annotations.Nullable String rest)
func (this_ *Es6RenameTypeReferences) priv_Es6RenameTypeReferences_renameReference_2576(t *NodeTraversal, n pkg_rhino.Node, oldName string, rest string) {
	_ = t
	_ = n
	_ = oldName
	_ = rest
	var current *Scope = t.GetScope()
	_ = current
	for current != nil {
		var currentRootNode pkg_rhino.Node = current.GetRootNode_0()
		_ = currentRootNode
		var newName string = j2gort.CastString(this_.renameTable.Get(j2gort.BoxNode(currentRootNode), j2gort.BoxString(oldName)))
		_ = newName
		if newName != j2gort.NullString {
			if !((currentRootNode.IsModuleBody() || currentRootNode.IsScript()) || NodeUtil_IsFunctionBlock(currentRootNode)) {
				j2gort.CheckStateFmt(false, "Not a MODULE_BODY or SCRIPT or goog.loadModule root: %s", j2gort.BoxNode(currentRootNode))
			}
			var newFullName string
			_ = newFullName
			if rest == j2gort.NullString {
				newFullName = newName
			} else {
				newFullName = (j2gort.StrValString(newName) + j2gort.StrValString(rest))
			}
			n.SetString(newFullName)
			return
		} else if current.HasOwnSlot(oldName) {
			return
		} else {
			current = j2gort.CastPtr[Scope](current.GetParent())
		}
	}
}
func (*Es6RenameTypeReferences) AlwaysShouldTraverse_()      {}
func (*Es6RenameTypeReferences) CanReleaseSyntacticScopes_() {}
