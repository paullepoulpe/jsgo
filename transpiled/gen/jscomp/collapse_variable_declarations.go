/*
 * Copyright 2006 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/CollapseVariableDeclarations.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.CollapseVariableDeclarations
type CollapseVariableDeclarations struct {
	compiler        AbstractCompiler                         `gson:"compiler"`
	collapses       []*CollapseVariableDeclarations_Collapse `gson:"collapses"`
	nodesToCollapse j2gort.OrderedSet[pkg_rhino.Node]        `gson:"nodesToCollapse"`
	idHash_         atomic.Int32
}

func (this_ *CollapseVariableDeclarations) BindSelf_(s any) {
	_ = s
}
func (this_ *CollapseVariableDeclarations) IsJavaNil_() bool { return this_ == nil }
func (*CollapseVariableDeclarations) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CollapseVariableDeclarations", "CollapseVariableDeclarations"
}
func (this_ *CollapseVariableDeclarations) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CollapseVariableDeclarations) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CollapseVariableDeclarations", "com.google.javascript.jscomp.CollapseVariableDeclarations", func(v_ any) bool {
		return j2gort.InstanceOf[*CollapseVariableDeclarations](v_)
	})
}
func (*CollapseVariableDeclarations) Implements_jscomp_CompilerPass() {}

func (this_ *CollapseVariableDeclarations) initFields_CollapseVariableDeclarations() {
	_ = this_
	this_.collapses = []*CollapseVariableDeclarations_Collapse{}
	this_.nodesToCollapse = j2gort.NewOrderedSet[pkg_rhino.Node]()
}

// Java constructor: CollapseVariableDeclarations(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewCollapseVariableDeclarations(compiler AbstractCompiler) *CollapseVariableDeclarations {
	this_ := &CollapseVariableDeclarations{}
	InitCollapseVariableDeclarations(this_, compiler)
	return this_
}

func InitCollapseVariableDeclarations(this_ *CollapseVariableDeclarations, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_CollapseVariableDeclarations()
	if !(!(compiler.GetLifeCycleStage().IsNormalized())) {
		j2gort.CheckState(false)
	}
	this_.compiler = compiler
}

// Java method: CollapseVariableDeclarations.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *CollapseVariableDeclarations) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	clear(this_.collapses)
	this_.collapses = this_.collapses[:0]
	this_.nodesToCollapse.Clear()
	NodeTraversal_Traverse(this_.compiler, root, NewCollapseVariableDeclarations_GatherCollapses(this_))
	if len(this_.collapses) != 0 {
		this_.priv_CollapseVariableDeclarations_applyCollapses_1396()
	}
}

// Java method: CollapseVariableDeclarations.applyCollapses()
func (this_ *CollapseVariableDeclarations) priv_CollapseVariableDeclarations_applyCollapses_1396() {
	for _, collapse := range this_.collapses {
		_ = collapse
		var var_ pkg_rhino.Node = collapse.startNode
		_ = var_
		this_.compiler.ReportChangeToEnclosingScope(var_)
		for (var_.GetNext() != 0) && (var_.GetNext().RawToken_() == var_.RawToken_()) {
			var next pkg_rhino.Node = var_.GetNext().Detach()
			_ = next
			var_.AddChildrenToBack(next.RemoveChildren())
		}
	}

}

// Java class: com.google.javascript.jscomp.CollapseVariableDeclarations.Collapse
type CollapseVariableDeclarations_Collapse struct {
	startNode pkg_rhino.Node `gson:"startNode"`
	idHash_   atomic.Int32
}

func (this_ *CollapseVariableDeclarations_Collapse) BindSelf_(s any) {
	_ = s
}
func (this_ *CollapseVariableDeclarations_Collapse) IsJavaNil_() bool { return this_ == nil }
func (*CollapseVariableDeclarations_Collapse) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CollapseVariableDeclarations$Collapse", "Collapse"
}
func (this_ *CollapseVariableDeclarations_Collapse) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CollapseVariableDeclarations_Collapse) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CollapseVariableDeclarations.Collapse", "com.google.javascript.jscomp.CollapseVariableDeclarations$Collapse", func(v_ any) bool {
		return j2gort.InstanceOf[*CollapseVariableDeclarations_Collapse](v_)
	})
}

func (this_ *CollapseVariableDeclarations_Collapse) initFields_CollapseVariableDeclarations_Collapse() {
	_ = this_
}

// Java constructor: Collapse(com.google.javascript.rhino.Node startNode)
func NewCollapseVariableDeclarations_Collapse(startNode pkg_rhino.Node) *CollapseVariableDeclarations_Collapse {
	this_ := &CollapseVariableDeclarations_Collapse{}
	InitCollapseVariableDeclarations_Collapse(this_, startNode)
	return this_
}

func InitCollapseVariableDeclarations_Collapse(this_ *CollapseVariableDeclarations_Collapse, startNode pkg_rhino.Node) {
	_ = startNode
	this_.initFields_CollapseVariableDeclarations_Collapse()
	this_.startNode = startNode
}

// Java class: com.google.javascript.jscomp.CollapseVariableDeclarations.GatherCollapses
type CollapseVariableDeclarations_GatherCollapses struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer_ *CollapseVariableDeclarations
}

func (this_ *CollapseVariableDeclarations_GatherCollapses) BindSelf_(s any) {
	_ = s
}
func (this_ *CollapseVariableDeclarations_GatherCollapses) IsJavaNil_() bool { return this_ == nil }
func (*CollapseVariableDeclarations_GatherCollapses) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CollapseVariableDeclarations$GatherCollapses", "GatherCollapses"
}
func (this_ *CollapseVariableDeclarations_GatherCollapses) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*CollapseVariableDeclarations_GatherCollapses) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CollapseVariableDeclarations.GatherCollapses", "com.google.javascript.jscomp.CollapseVariableDeclarations$GatherCollapses", func(v_ any) bool {
		return j2gort.InstanceOf[*CollapseVariableDeclarations_GatherCollapses](v_)
	})
}
func (this_ *CollapseVariableDeclarations_GatherCollapses) SetOuter_(o *CollapseVariableDeclarations) {
	this_.outer_ = o
}
func (this_ *CollapseVariableDeclarations_GatherCollapses) GetOuter_() *CollapseVariableDeclarations {
	return this_.outer_
}
func (*CollapseVariableDeclarations_GatherCollapses) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {
}
func (*CollapseVariableDeclarations_GatherCollapses) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *CollapseVariableDeclarations_GatherCollapses) initFields_CollapseVariableDeclarations_GatherCollapses() {
	_ = this_
}

// Java constructor: GatherCollapses()
func NewCollapseVariableDeclarations_GatherCollapses(outer_ *CollapseVariableDeclarations) *CollapseVariableDeclarations_GatherCollapses {
	this_ := &CollapseVariableDeclarations_GatherCollapses{}
	this_.outer_ = outer_
	InitCollapseVariableDeclarations_GatherCollapses(this_)
	return this_
}

func InitCollapseVariableDeclarations_GatherCollapses(this_ *CollapseVariableDeclarations_GatherCollapses) {
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_CollapseVariableDeclarations_GatherCollapses()
}

// Java method: GatherCollapses.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *CollapseVariableDeclarations_GatherCollapses) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if this_.outer_.nodesToCollapse.Contains(n) {
		return
	}
	if !(NodeUtil_IsNameDeclaration(n)) {
		return
	}
	if parent.IsIf() {
		return
	}
	var varNode pkg_rhino.Node = n
	_ = varNode
	var nType pkg_rhino.Token = n.GetToken()
	_ = nType
	n = n.GetNext()
	var hasNodesToCollapse bool = false
	_ = hasNodesToCollapse
	for (n != 0) && (nType == n.GetToken()) {
		this_.outer_.nodesToCollapse.Add(n)
		hasNodesToCollapse = true
		n = n.GetNext()
	}
	if hasNodesToCollapse {
		this_.outer_.nodesToCollapse.Add(varNode)
		this_.outer_.collapses = append(this_.outer_.collapses, NewCollapseVariableDeclarations_Collapse(varNode))
	}
}
func (*CollapseVariableDeclarations_GatherCollapses) AlwaysShouldTraverse_()      {}
func (*CollapseVariableDeclarations_GatherCollapses) CanReleaseSyntacticScopes_() {}
