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
// Port of com/google/javascript/jscomp/Es6RenameReferences.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.Es6RenameReferences
type Es6RenameReferences struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	renameTable j2gort.Ext_com_google_common_collect_Table `gson:"renameTable"`
}

func (this_ *Es6RenameReferences) BindSelf_(s any) {
	_ = s
}
func (this_ *Es6RenameReferences) IsJavaNil_() bool { return this_ == nil }
func (*Es6RenameReferences) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Es6RenameReferences", "Es6RenameReferences"
}
func (this_ *Es6RenameReferences) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*Es6RenameReferences) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Es6RenameReferences", "com.google.javascript.jscomp.Es6RenameReferences", func(v_ any) bool {
		return j2gort.InstanceOf[*Es6RenameReferences](v_)
	})
}
func (*Es6RenameReferences) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {}
func (*Es6RenameReferences) Implements_jscomp_NodeTraversal_Callback()                  {}

func (this_ *Es6RenameReferences) initFields_Es6RenameReferences() {
	_ = this_
}

// Java constructor: Es6RenameReferences(com.google.common.collect.Table<com.google.javascript.rhino.Node,java.lang.String,java.lang.String> renameTable)
func NewEs6RenameReferences(renameTable j2gort.Ext_com_google_common_collect_Table) *Es6RenameReferences {
	this_ := &Es6RenameReferences{}
	InitEs6RenameReferences(this_, renameTable)
	return this_
}

func InitEs6RenameReferences(this_ *Es6RenameReferences, renameTable j2gort.Ext_com_google_common_collect_Table) {
	_ = renameTable
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_Es6RenameReferences()
	this_.renameTable = renameTable
}

// Java method: Es6RenameReferences.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *Es6RenameReferences) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if NodeUtil_IsReferenceName(n) {
		this_.priv_Es6RenameReferences_renameReference_2556(t, n, n.GetString())
	}
}

// Java method: Es6RenameReferences.renameReference(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, java.lang.String oldName)
func (this_ *Es6RenameReferences) priv_Es6RenameReferences_renameReference_2556(t *NodeTraversal, n pkg_rhino.Node, oldName string) {
	_ = t
	_ = n
	_ = oldName
	var current *Scope = t.GetScope()
	_ = current
	for current != nil {
		var newName string = j2gort.CastString(this_.renameTable.Get(j2gort.BoxNode(current.GetRootNode_0()), j2gort.BoxString(oldName)))
		_ = newName
		if newName != j2gort.NullString {
			n.SetString(newName)
			t.ReportCodeChange_0()
			return
		} else if current.HasOwnSlot(oldName) {
			return
		} else {
			current = j2gort.CastPtr[Scope](current.GetParent())
		}
	}
}
func (*Es6RenameReferences) AlwaysShouldTraverse_()      {}
func (*Es6RenameReferences) CanReleaseSyntacticScopes_() {}
