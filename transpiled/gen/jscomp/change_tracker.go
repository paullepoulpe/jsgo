/*
 * Copyright 2025 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/ChangeTracker.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.ChangeTracker
type ChangeTracker struct {
	changeTimeline     *Timeline                 `gson:"changeTimeline"`
	recentChange       *RecentChange             `gson:"recentChange"`
	codeChangeHandlers j2gort.Ext_java_util_List `gson:"codeChangeHandlers"`
	changeStamp        int32                     `gson:"changeStamp"`
	idHash_            atomic.Int32
}

func (this_ *ChangeTracker) BindSelf_(s any) {
	_ = s
}
func (this_ *ChangeTracker) IsJavaNil_() bool { return this_ == nil }
func (*ChangeTracker) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ChangeTracker", "ChangeTracker"
}
func (this_ *ChangeTracker) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ChangeTracker) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ChangeTracker", "com.google.javascript.jscomp.ChangeTracker", func(v_ any) bool {
		return j2gort.InstanceOf[*ChangeTracker](v_)
	})
}

func (this_ *ChangeTracker) initFields_ChangeTracker() {
	_ = this_
	this_.changeStamp = int32(1)
	this_.changeTimeline = NewTimeline()
	this_.recentChange = NewRecentChange()
	this_.codeChangeHandlers = j2gort.NewExt_java_util_ArrayList()
}

// Java constructor: ChangeTracker()
func NewChangeTracker() *ChangeTracker {
	this_ := &ChangeTracker{}
	InitChangeTracker(this_)
	return this_
}

func InitChangeTracker(this_ *ChangeTracker) {
	this_.initFields_ChangeTracker()
}

// Java method: ChangeTracker.addChangeHandler(com.google.javascript.jscomp.CodeChangeHandler handler)
func (this_ *ChangeTracker) AddChangeHandler(handler CodeChangeHandler) {
	_ = handler
	this_.codeChangeHandlers.Add_1(handler)
}

// Java method: ChangeTracker.removeChangeHandler(com.google.javascript.jscomp.CodeChangeHandler handler)
func (this_ *ChangeTracker) RemoveChangeHandler(handler CodeChangeHandler) {
	_ = handler
	this_.codeChangeHandlers.Remove_1_any(handler)
}

// Java method: ChangeTracker.getRecentChange()
func (this_ *ChangeTracker) GetRecentChange() *RecentChange {
	return this_.recentChange
}

// Java method: ChangeTracker.reportChangeToEnclosingScope(com.google.javascript.rhino.Node n)
func (this_ *ChangeTracker) ReportChangeToEnclosingScope(n pkg_rhino.Node) {
	_ = n
	this_.priv_ChangeTracker_recordChange_973(this_.priv_ChangeTracker_getChangeScopeForNode_972(n))
	this_.priv_ChangeTracker_notifyChangeHandlers_971()
}

// Java method: ChangeTracker.reportChangeToChangeScope(com.google.javascript.rhino.Node changeScopeRoot)
func (this_ *ChangeTracker) ReportChangeToChangeScope(changeScopeRoot pkg_rhino.Node) {
	_ = changeScopeRoot
	if !(changeScopeRoot.IsScript() || changeScopeRoot.IsFunction()) {
		j2gort.CheckState(false)
	}
	this_.priv_ChangeTracker_recordChange_973(changeScopeRoot)
	this_.priv_ChangeTracker_notifyChangeHandlers_971()
}

// Java method: ChangeTracker.markNewScopesChanged(com.google.javascript.rhino.Node node)
func (this_ *ChangeTracker) MarkNewScopesChanged(node pkg_rhino.Node) {
	_ = node
	if node.IsFunction() {
		this_.ReportChangeToChangeScope(node)
	}
	{
		var child pkg_rhino.Node = node.GetFirstChild()
		_ = child
		for ; child != 0; child = child.GetNext() {
			this_.MarkNewScopesChanged(child)
		}
	}
}

// Java method: ChangeTracker.reportFunctionDeleted(com.google.javascript.rhino.Node n)
func (this_ *ChangeTracker) ReportFunctionDeleted(n pkg_rhino.Node) {
	_ = n
	if !(n.IsFunction()) {
		j2gort.CheckState(false)
	}
	n.SetDeleted(true)
	this_.changeTimeline.RemoveNode_(n)
}

// Java method: ChangeTracker.getChangedScopeNodesForPass(java.lang.String passName)
func (this_ *ChangeTracker) GetChangedScopeNodesForPass(passName string) j2gort.Ext_java_util_List {
	_ = passName
	var changedScopeNodes j2gort.Ext_java_util_List = this_.changeTimeline.GetSince(passName)
	_ = changedScopeNodes
	this_.changeTimeline.Mark(passName)
	return changedScopeNodes
}

// Java method: ChangeTracker.getChangeStamp()
func (this_ *ChangeTracker) GetChangeStamp() int32 {
	return this_.changeStamp
}

// Java method: ChangeTracker.incrementChangeStamp()
func (this_ *ChangeTracker) IncrementChangeStamp() {
	this_.changeStamp = this_.changeStamp + 1
}

// Java method: ChangeTracker.resetChangeStamp()
func (this_ *ChangeTracker) ResetChangeStamp() {
	this_.changeStamp = int32(1)
}

// Java static method: ChangeTracker.isChangeScopeRoot(com.google.javascript.rhino.Node n)
func ChangeTracker_IsChangeScopeRoot(n pkg_rhino.Node) bool {
	_ = n
	return (n.IsScript() || n.IsFunction())
}

// Java static method: ChangeTracker.getEnclosingChangeScopeRoot(com.google.javascript.rhino.Node n)
func ChangeTracker_GetEnclosingChangeScopeRoot(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	for (n != 0) && !(ChangeTracker_IsChangeScopeRoot(n)) {
		n = n.GetParent()
	}
	return n
}

// Java method: ChangeTracker.notifyChangeHandlers()
func (this_ *ChangeTracker) priv_ChangeTracker_notifyChangeHandlers_971() {
	{
		cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.codeChangeHandlers))
		for cur_1.HasNext() {
			rawElem_2 := cur_1.Next()
			handler := j2gort.AsIface[CodeChangeHandler](rawElem_2)
			_ = handler
			handler.ReportChange()
		}
	}
}

// Java method: ChangeTracker.getChangeScopeForNode(com.google.javascript.rhino.Node n)
func (this_ *ChangeTracker) priv_ChangeTracker_getChangeScopeForNode_972(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	if n.IsScript() {
		return n
	}
	var enclosingScopeNode pkg_rhino.Node = ChangeTracker_GetEnclosingChangeScopeRoot(n.GetParent())
	_ = enclosingScopeNode
	if enclosingScopeNode == 0 {
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "An enclosing scope is required for change reports but node "+j2gort.StrVal(n)+" doesn't have one.")))
	}
	return enclosingScopeNode
}

// Java method: ChangeTracker.recordChange(com.google.javascript.rhino.Node n)
func (this_ *ChangeTracker) priv_ChangeTracker_recordChange_973(n pkg_rhino.Node) {
	_ = n
	if n.IsDeleted() {
		return
	}
	n.SetChangeTime(this_.changeStamp)
	this_.IncrementChangeStamp()
	this_.changeTimeline.AddNode_(n)
}
