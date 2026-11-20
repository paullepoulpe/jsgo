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
// Port of com/google/javascript/jscomp/RenameLabels.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.RenameLabels
type RenameLabels struct {
	compiler     AbstractCompiler                       `gson:"compiler"`
	nameSupplier j2gort.Ext_java_util_function_Supplier `gson:"nameSupplier"`
	idHash_      atomic.Int32
	removeUnused bool `gson:"removeUnused"`
	markChanges  bool `gson:"markChanges"`
}

func (this_ *RenameLabels) BindSelf_(s any) {
	_ = s
}
func (this_ *RenameLabels) IsJavaNil_() bool { return this_ == nil }
func (*RenameLabels) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RenameLabels", "RenameLabels"
}
func (this_ *RenameLabels) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*RenameLabels) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RenameLabels", "com.google.javascript.jscomp.RenameLabels", func(v_ any) bool {
		return j2gort.InstanceOf[*RenameLabels](v_)
	})
}
func (*RenameLabels) Implements_jscomp_CompilerPass() {}

func (this_ *RenameLabels) initFields_RenameLabels() {
	_ = this_
}

// Java constructor: RenameLabels(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewRenameLabels_1(compiler AbstractCompiler) *RenameLabels {
	this_ := &RenameLabels{}
	InitRenameLabels_1(this_, compiler)
	return this_
}

func InitRenameLabels_1(this_ *RenameLabels, compiler AbstractCompiler) {
	_ = compiler
	InitRenameLabels_4(this_, compiler, NewRenameLabels_DefaultNameSupplier(), true, true)
}

// Java constructor: RenameLabels(com.google.javascript.jscomp.AbstractCompiler compiler, java.util.function.Supplier<java.lang.String> supplier, boolean removeUnused, boolean markChanges)
func NewRenameLabels_4(compiler AbstractCompiler, supplier j2gort.Ext_java_util_function_Supplier, removeUnused bool, markChanges bool) *RenameLabels {
	this_ := &RenameLabels{}
	InitRenameLabels_4(this_, compiler, supplier, removeUnused, markChanges)
	return this_
}

func InitRenameLabels_4(this_ *RenameLabels, compiler AbstractCompiler, supplier j2gort.Ext_java_util_function_Supplier, removeUnused bool, markChanges bool) {
	_ = compiler
	_ = supplier
	_ = removeUnused
	_ = markChanges
	this_.initFields_RenameLabels()
	this_.compiler = compiler
	this_.nameSupplier = supplier
	this_.removeUnused = removeUnused
	this_.markChanges = markChanges
}

// Java method: RenameLabels.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *RenameLabels) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_Traverse(this_.compiler, root, NewRenameLabels_ProcessLabels(this_, this_.markChanges))
}

// Java class: com.google.javascript.jscomp.RenameLabels.DefaultNameSupplier
type RenameLabels_DefaultNameSupplier struct {
	nameGenerator NameGenerator `gson:"nameGenerator"`
	idHash_       atomic.Int32
}

func (this_ *RenameLabels_DefaultNameSupplier) BindSelf_(s any) {
	_ = s
}
func (this_ *RenameLabels_DefaultNameSupplier) IsJavaNil_() bool { return this_ == nil }
func (*RenameLabels_DefaultNameSupplier) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RenameLabels$DefaultNameSupplier", "DefaultNameSupplier"
}
func (this_ *RenameLabels_DefaultNameSupplier) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*RenameLabels_DefaultNameSupplier) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RenameLabels.DefaultNameSupplier", "com.google.javascript.jscomp.RenameLabels$DefaultNameSupplier", func(v_ any) bool {
		return j2gort.InstanceOf[*RenameLabels_DefaultNameSupplier](v_)
	})
}
func (*RenameLabels_DefaultNameSupplier) Implements_j2gort_Ext_java_util_function_Supplier() {}

func (this_ *RenameLabels_DefaultNameSupplier) initFields_RenameLabels_DefaultNameSupplier() {
	_ = this_
}

// Java constructor: DefaultNameSupplier()
func NewRenameLabels_DefaultNameSupplier() *RenameLabels_DefaultNameSupplier {
	this_ := &RenameLabels_DefaultNameSupplier{}
	InitRenameLabels_DefaultNameSupplier(this_)
	return this_
}

func InitRenameLabels_DefaultNameSupplier(this_ *RenameLabels_DefaultNameSupplier) {
	this_.initFields_RenameLabels_DefaultNameSupplier()
	this_.nameGenerator = NewDefaultNameGenerator_3(j2gort.NewExt_java_util_LinkedHashSet(), "", j2gort.AsIface[j2gort.Ext_java_util_Set](j2gort.Ext_com_google_common_collect_ImmutableSet_Of_0()))
}

// Java method: DefaultNameSupplier.get()
func (this_ *RenameLabels_DefaultNameSupplier) Get() any {
	return j2gort.BoxString(this_.nameGenerator.GenerateNextName())
}

// Java class: com.google.javascript.jscomp.RenameLabels.ProcessLabels
type RenameLabels_ProcessLabels struct {
	outer_         *RenameLabels
	namespaceStack j2gort.Deque[*RenameLabels_LabelNamespace] `gson:"namespaceStack"`
	names          []string                                   `gson:"names"`
	idHash_        atomic.Int32
	markChanges    bool `gson:"markChanges"`
}

func (this_ *RenameLabels_ProcessLabels) BindSelf_(s any) {
	_ = s
}
func (this_ *RenameLabels_ProcessLabels) IsJavaNil_() bool { return this_ == nil }
func (*RenameLabels_ProcessLabels) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RenameLabels$ProcessLabels", "ProcessLabels"
}
func (this_ *RenameLabels_ProcessLabels) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*RenameLabels_ProcessLabels) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RenameLabels.ProcessLabels", "com.google.javascript.jscomp.RenameLabels$ProcessLabels", func(v_ any) bool {
		return j2gort.InstanceOf[*RenameLabels_ProcessLabels](v_)
	})
}
func (this_ *RenameLabels_ProcessLabels) SetOuter_(o *RenameLabels)                 { this_.outer_ = o }
func (this_ *RenameLabels_ProcessLabels) GetOuter_() *RenameLabels                  { return this_.outer_ }
func (*RenameLabels_ProcessLabels) Implements_jscomp_NodeTraversal_ScopedCallback() {}
func (*RenameLabels_ProcessLabels) Implements_jscomp_NodeTraversal_Callback()       {}

func (this_ *RenameLabels_ProcessLabels) initFields_RenameLabels_ProcessLabels() {
	_ = this_
	this_.namespaceStack = j2gort.Deque[*RenameLabels_LabelNamespace]{}
	this_.names = []string{}
}

// Java constructor: ProcessLabels(boolean markChanges)
func NewRenameLabels_ProcessLabels(outer_ *RenameLabels, markChanges bool) *RenameLabels_ProcessLabels {
	this_ := &RenameLabels_ProcessLabels{}
	this_.outer_ = outer_
	InitRenameLabels_ProcessLabels(this_, markChanges)
	return this_
}

func InitRenameLabels_ProcessLabels(this_ *RenameLabels_ProcessLabels, markChanges bool) {
	_ = markChanges
	this_.initFields_RenameLabels_ProcessLabels()
	this_.markChanges = markChanges
	this_.namespaceStack.Push(NewRenameLabels_LabelNamespace())
}

// Java method: ProcessLabels.enterScope(com.google.javascript.jscomp.NodeTraversal nodeTraversal)
func (this_ *RenameLabels_ProcessLabels) EnterScope(nodeTraversal *NodeTraversal) {
	_ = nodeTraversal
	if nodeTraversal.GetScopeRoot().IsFunction() {
		this_.namespaceStack.Push(NewRenameLabels_LabelNamespace())
	}
}

// Java method: ProcessLabels.exitScope(com.google.javascript.jscomp.NodeTraversal nodeTraversal)
func (this_ *RenameLabels_ProcessLabels) ExitScope(nodeTraversal *NodeTraversal) {
	_ = nodeTraversal
	if nodeTraversal.GetScopeRoot().IsFunction() {
		this_.namespaceStack.Pop()
	}
}

// Java method: ProcessLabels.shouldTraverse(com.google.javascript.jscomp.NodeTraversal nodeTraversal, com.google.javascript.rhino.Node node, com.google.javascript.rhino.Node parent)
func (this_ *RenameLabels_ProcessLabels) ShouldTraverse_3(nodeTraversal *NodeTraversal, node pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = nodeTraversal
	_ = node
	_ = parent
	if node.IsLabel() {
		var current *RenameLabels_LabelNamespace = this_.namespaceStack.Peek()
		_ = current
		var currentDepth int32 = (current.renameMap.Size() + int32(1))
		_ = currentDepth
		var name string = node.GetFirstChild().GetString()
		_ = name
		var val_1 RenameLabels_LabelInfo
		InitRenameLabels_LabelInfo(&val_1, currentDepth)
		var li *RenameLabels_LabelInfo = &val_1
		_ = li
		if !(!(current.renameMap.ContainsKey(name))) {
			j2gort.CheckState(false)
		}
		current.renameMap.Put(name, li)
		if int32(len(this_.names)) < currentDepth {
			this_.names = append(this_.names, j2gort.CastString(this_.outer_.nameSupplier.Get()))
		}
	}
	return true
}

// Java method: ProcessLabels.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node node, com.google.javascript.rhino.Node parent)
func (this_ *RenameLabels_ProcessLabels) Visit_3(t *NodeTraversal, node pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = node
	_ = parent
	switch node.TokenSwitchOrd_() {
	case int32(92) /* Token.LABEL */ :
		this_.priv_ProcessLabels_visitLabel_4831(t, node)
	case int32(82) /* Token.BREAK */, int32(83) /* Token.CONTINUE */ :
		this_.priv_ProcessLabels_visitBreakOrContinue_4830(t, node)
	default:
	}
}

// Java method: ProcessLabels.visitBreakOrContinue(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node node)
func (this_ *RenameLabels_ProcessLabels) priv_ProcessLabels_visitBreakOrContinue_4830(t *NodeTraversal, node pkg_rhino.Node) {
	_ = t
	_ = node
	var nameNode pkg_rhino.Node = node.GetFirstChild()
	_ = nameNode
	if nameNode != 0 {
		var name string = nameNode.GetString()
		_ = name
		if !(!(j2gort.CheckNotNullString(name) == "")) {
			j2gort.CheckState(false)
		}
		var li *RenameLabels_LabelInfo = this_.GetLabelInfo(name)
		_ = li
		if li != nil {
			var newName string = this_.GetNameForId(li.id)
			_ = newName
			{
				r_2 := li
				r_2.referenced = true
			}
			if !(j2gort.CheckNotNullString(name) == newName) {
				nameNode.SetString(newName)
				if this_.markChanges {
					t.ReportCodeChange_0()
				}
			}
		}
	}
}

// Java method: ProcessLabels.visitLabel(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node node)
func (this_ *RenameLabels_ProcessLabels) priv_ProcessLabels_visitLabel_4831(t *NodeTraversal, node pkg_rhino.Node) {
	_ = t
	_ = node
	var nameNode pkg_rhino.Node = node.GetFirstChild()
	_ = nameNode
	if !(nameNode != 0) {
		j2gort.CheckState(false)
	}
	var name string = nameNode.GetString()
	_ = name
	var li *RenameLabels_LabelInfo = this_.GetLabelInfo(name)
	_ = li
	if li.referenced || !(this_.outer_.removeUnused) {
		var newName string = this_.GetNameForId(li.id)
		_ = newName
		if !(j2gort.CheckNotNullString(name) == newName) {
			nameNode.SetString(newName)
			if this_.markChanges {
				t.ReportCodeChange_0()
			}
		}
	} else {
		var newChild pkg_rhino.Node = node.GetLastChild()
		_ = newChild
		newChild.Detach()
		node.ReplaceWith(newChild)
		if newChild.IsBlock() {
			NodeUtil_TryMergeBlock(newChild, false)
		}
		if this_.markChanges {
			t.ReportCodeChange_0()
		}
	}
	this_.namespaceStack.Peek().renameMap.Remove(name)
}

// Java method: ProcessLabels.getNameForId(int id)
func (this_ *RenameLabels_ProcessLabels) GetNameForId(id int32) string {
	_ = id
	return this_.names[(id - int32(1))]
}

// Java method: ProcessLabels.getLabelInfo(java.lang.String name)
func (this_ *RenameLabels_ProcessLabels) GetLabelInfo(name string) *RenameLabels_LabelInfo {
	_ = name
	return this_.namespaceStack.Peek().renameMap.Get(name)
}

// Java class: com.google.javascript.jscomp.RenameLabels.LabelInfo
type RenameLabels_LabelInfo struct {
	id         int32 `gson:"id"`
	idHash_    atomic.Int32
	referenced bool `gson:"referenced"`
}

func (this_ *RenameLabels_LabelInfo) BindSelf_(s any) {
	_ = s
}
func (this_ *RenameLabels_LabelInfo) IsJavaNil_() bool { return this_ == nil }
func (*RenameLabels_LabelInfo) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RenameLabels$LabelInfo", "LabelInfo"
}
func (this_ *RenameLabels_LabelInfo) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*RenameLabels_LabelInfo) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RenameLabels.LabelInfo", "com.google.javascript.jscomp.RenameLabels$LabelInfo", func(v_ any) bool {
		return j2gort.InstanceOf[*RenameLabels_LabelInfo](v_)
	})
}

func (this_ *RenameLabels_LabelInfo) initFields_RenameLabels_LabelInfo() {
	_ = this_

}

// Java constructor: LabelInfo(int id)
func NewRenameLabels_LabelInfo(id int32) *RenameLabels_LabelInfo {
	this_ := &RenameLabels_LabelInfo{}
	InitRenameLabels_LabelInfo(this_, id)
	return this_
}

func InitRenameLabels_LabelInfo(this_ *RenameLabels_LabelInfo, id int32) {
	_ = id
	this_.initFields_RenameLabels_LabelInfo()
	this_.id = id
}

// Java class: com.google.javascript.jscomp.RenameLabels.LabelNamespace
type RenameLabels_LabelNamespace struct {
	renameMap j2gort.OrderedMap[string, *RenameLabels_LabelInfo] `gson:"renameMap"`
	idHash_   atomic.Int32
}

func (this_ *RenameLabels_LabelNamespace) BindSelf_(s any) {
	_ = s
}
func (this_ *RenameLabels_LabelNamespace) IsJavaNil_() bool { return this_ == nil }
func (*RenameLabels_LabelNamespace) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RenameLabels$LabelNamespace", "LabelNamespace"
}
func (this_ *RenameLabels_LabelNamespace) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*RenameLabels_LabelNamespace) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RenameLabels.LabelNamespace", "com.google.javascript.jscomp.RenameLabels$LabelNamespace", func(v_ any) bool {
		return j2gort.InstanceOf[*RenameLabels_LabelNamespace](v_)
	})
}

func (this_ *RenameLabels_LabelNamespace) initFields_RenameLabels_LabelNamespace() {
	_ = this_
	this_.renameMap = j2gort.NewOrderedMap[string, *RenameLabels_LabelInfo]()
}

// Java constructor: LabelNamespace()
func NewRenameLabels_LabelNamespace() *RenameLabels_LabelNamespace {
	this_ := &RenameLabels_LabelNamespace{}
	InitRenameLabels_LabelNamespace(this_)
	return this_
}

func InitRenameLabels_LabelNamespace(this_ *RenameLabels_LabelNamespace) {
	this_.initFields_RenameLabels_LabelNamespace()
}
func (*RenameLabels_ProcessLabels) CanReleaseSyntacticScopes_() {}
