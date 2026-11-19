/*
 * Copyright 2009 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/CreateSyntheticBlocks.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.CreateSyntheticBlocks
type CreateSyntheticBlocks struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler        AbstractCompiler                `gson:"compiler"`
	startMarkerName string                          `gson:"startMarkerName"`
	endMarkerName   string                          `gson:"endMarkerName"`
	markerStack     j2gort.Deque[pkg_rhino.Node]    `gson:"markerStack"`
	validMarkers    []*CreateSyntheticBlocks_Marker `gson:"validMarkers"`
}

func (this_ *CreateSyntheticBlocks) BindSelf_(s any) {
	_ = s
}
func (this_ *CreateSyntheticBlocks) IsJavaNil_() bool { return this_ == nil }
func (*CreateSyntheticBlocks) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CreateSyntheticBlocks", "CreateSyntheticBlocks"
}
func (this_ *CreateSyntheticBlocks) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*CreateSyntheticBlocks) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CreateSyntheticBlocks", "com.google.javascript.jscomp.CreateSyntheticBlocks", func(v_ any) bool {
		return j2gort.InstanceOf[*CreateSyntheticBlocks](v_)
	})
}
func (*CreateSyntheticBlocks) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {}
func (*CreateSyntheticBlocks) Implements_jscomp_CompilerPass()                            {}
func (*CreateSyntheticBlocks) Implements_jscomp_NodeTraversal_Callback()                  {}

var static_CreateSyntheticBlocks_UNMATCHED_START_MARKER *DiagnosticType

func StaticGet_CreateSyntheticBlocks_UNMATCHED_START_MARKER() *DiagnosticType {
	clinit_CreateSyntheticBlocks()
	return static_CreateSyntheticBlocks_UNMATCHED_START_MARKER
}
func StaticSet_CreateSyntheticBlocks_UNMATCHED_START_MARKER(v *DiagnosticType) {
	clinit_CreateSyntheticBlocks()
	static_CreateSyntheticBlocks_UNMATCHED_START_MARKER = v
}

var static_CreateSyntheticBlocks_UNMATCHED_END_MARKER *DiagnosticType

func StaticGet_CreateSyntheticBlocks_UNMATCHED_END_MARKER() *DiagnosticType {
	clinit_CreateSyntheticBlocks()
	return static_CreateSyntheticBlocks_UNMATCHED_END_MARKER
}
func StaticSet_CreateSyntheticBlocks_UNMATCHED_END_MARKER(v *DiagnosticType) {
	clinit_CreateSyntheticBlocks()
	static_CreateSyntheticBlocks_UNMATCHED_END_MARKER = v
}

var static_CreateSyntheticBlocks_INVALID_MARKER_USAGE *DiagnosticType

func StaticGet_CreateSyntheticBlocks_INVALID_MARKER_USAGE() *DiagnosticType {
	clinit_CreateSyntheticBlocks()
	return static_CreateSyntheticBlocks_INVALID_MARKER_USAGE
}
func StaticSet_CreateSyntheticBlocks_INVALID_MARKER_USAGE(v *DiagnosticType) {
	clinit_CreateSyntheticBlocks()
	static_CreateSyntheticBlocks_INVALID_MARKER_USAGE = v
}

var clinitDone_CreateSyntheticBlocks uint32
var clinitState_CreateSyntheticBlocks j2gort.ClinitState

func clinit_CreateSyntheticBlocks() {
	if atomic.LoadUint32(&clinitDone_CreateSyntheticBlocks) != 0 {
		return
	}
	clinitSlow_CreateSyntheticBlocks()
}

func clinitSlow_CreateSyntheticBlocks() {
	if !j2gort.EnterClinit(&clinitState_CreateSyntheticBlocks, &clinitDone_CreateSyntheticBlocks) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_CreateSyntheticBlocks, &clinitDone_CreateSyntheticBlocks)
	static_CreateSyntheticBlocks_UNMATCHED_START_MARKER = DiagnosticType_Error("JSC_UNMATCHED_START_MARKER", "Unmatched {0}")
	static_CreateSyntheticBlocks_UNMATCHED_END_MARKER = DiagnosticType_Error("JSC_UNMATCHED_END_MARKER", "Unmatched {1} - {0} not in the same block")
	static_CreateSyntheticBlocks_INVALID_MARKER_USAGE = DiagnosticType_Error("JSC_INVALID_MARKER_USAGE", "Marker {0} can only be used in a simple call expression")
}

func (this_ *CreateSyntheticBlocks) initFields_CreateSyntheticBlocks() {
	_ = this_
	this_.startMarkerName = j2gort.NullString
	this_.endMarkerName = j2gort.NullString
	this_.markerStack = j2gort.Deque[pkg_rhino.Node]{}
	this_.validMarkers = []*CreateSyntheticBlocks_Marker{}
}

// Java constructor: CreateSyntheticBlocks(com.google.javascript.jscomp.AbstractCompiler compiler, java.lang.String startMarkerName, java.lang.String endMarkerName)
func NewCreateSyntheticBlocks(compiler AbstractCompiler, startMarkerName string, endMarkerName string) *CreateSyntheticBlocks {
	clinit_CreateSyntheticBlocks()
	this_ := &CreateSyntheticBlocks{}
	InitCreateSyntheticBlocks(this_, compiler, startMarkerName, endMarkerName)
	return this_
}

func InitCreateSyntheticBlocks(this_ *CreateSyntheticBlocks, compiler AbstractCompiler, startMarkerName string, endMarkerName string) {
	_ = compiler
	_ = startMarkerName
	_ = endMarkerName
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_CreateSyntheticBlocks()
	this_.compiler = compiler
	this_.startMarkerName = startMarkerName
	this_.endMarkerName = endMarkerName
}

// Java method: CreateSyntheticBlocks.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *CreateSyntheticBlocks) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_Traverse(this_.compiler, root, this_)
	{
		cur_1 := this_.markerStack.Iter()
		for cur_1.HasNext() {
			node := cur_1.Next()
			_ = node
			this_.compiler.Report_1(JSError_Make_3(node, StaticGet_CreateSyntheticBlocks_UNMATCHED_START_MARKER(), []string{this_.startMarkerName}))
		}
	}
	for _, marker := range this_.validMarkers {
		_ = marker
		this_.priv_CreateSyntheticBlocks_addBlocks_2231(marker)
	}

}

// Java method: CreateSyntheticBlocks.rewriteFunctionDeclaration(com.google.javascript.rhino.Node n)
func (this_ *CreateSyntheticBlocks) priv_CreateSyntheticBlocks_rewriteFunctionDeclaration_2229(n pkg_rhino.Node) {
	_ = n
	var oldNameNode pkg_rhino.Node = n.GetFirstChild()
	_ = oldNameNode
	var fnNameNode pkg_rhino.Node = oldNameNode.CloneNode_0()
	_ = fnNameNode
	var var_ pkg_rhino.Node = pkg_rhino.IR_Var_1(fnNameNode).Srcref(n)
	_ = var_
	oldNameNode.SetString("")
	this_.compiler.ReportChangeToEnclosingScope(oldNameNode)
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	n.Detach()
	parent.AddChildToFront(var_)
	this_.compiler.ReportChangeToEnclosingScope(var_)
	fnNameNode.AddChildToFront(n)
}

// Java method: CreateSyntheticBlocks.rewriteFunctionDeclarationsInBlock(com.google.javascript.rhino.Node block)
func (this_ *CreateSyntheticBlocks) priv_CreateSyntheticBlocks_rewriteFunctionDeclarationsInBlock_2230(block pkg_rhino.Node) {
	_ = block
	if !(block.IsBlock()) {
		j2gort.CheckState(false)
	}
	var next pkg_rhino.Node = pkg_rhino.Node(0)
	_ = next
	{
		var current pkg_rhino.Node = block.GetFirstChild()
		_ = current
		for ; current != 0; current = next {
			next = current.GetNext()
			if NodeUtil_IsFunctionDeclaration(current) {
				this_.priv_CreateSyntheticBlocks_rewriteFunctionDeclaration_2229(current)
			}
		}
	}
}

// Java method: CreateSyntheticBlocks.addBlocks(com.google.javascript.jscomp.CreateSyntheticBlocks.Marker marker)
func (this_ *CreateSyntheticBlocks) priv_CreateSyntheticBlocks_addBlocks_2231(marker *CreateSyntheticBlocks_Marker) {
	_ = marker
	var outerBlock pkg_rhino.Node = pkg_rhino.IR_Block_0().Srcref(marker.startMarker)
	_ = outerBlock
	outerBlock.SetIsSyntheticBlock(true)
	outerBlock.InsertBefore(marker.startMarker)
	var innerBlock pkg_rhino.Node = pkg_rhino.IR_Block_0().Srcref(marker.startMarker)
	_ = innerBlock
	innerBlock.SetIsSyntheticBlock(true)
	this_.priv_CreateSyntheticBlocks_moveSiblingExclusive_2232(innerBlock, marker.startMarker, marker.endMarker)
	outerBlock.AddChildToBack(outerBlock.GetNext().Detach())
	outerBlock.AddChildToBack(innerBlock)
	outerBlock.AddChildToBack(outerBlock.GetNext().Detach())
	this_.priv_CreateSyntheticBlocks_rewriteFunctionDeclarationsInBlock_2230(innerBlock)
	if !(outerBlock.GetParent().IsBlock() || outerBlock.GetParent().IsScript()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(outerBlock.GetParent()))
	}
	if !(innerBlock.GetParent().IsBlock()) {
		j2gort.CheckState(false)
	}
	this_.compiler.ReportChangeToEnclosingScope(outerBlock)
}

// Java method: CreateSyntheticBlocks.moveSiblingExclusive(com.google.javascript.rhino.Node dest, com.google.javascript.rhino.Node start, com.google.javascript.rhino.Node end)
func (this_ *CreateSyntheticBlocks) priv_CreateSyntheticBlocks_moveSiblingExclusive_2232(dest pkg_rhino.Node, start pkg_rhino.Node, end pkg_rhino.Node) {
	_ = dest
	_ = start
	_ = end
	j2gort.CheckNotNullNode(start)
	j2gort.CheckNotNullNode(end)
	for start.GetNext() != end {
		var child pkg_rhino.Node = start.GetNext().Detach()
		_ = child
		dest.AddChildToBack(child)
	}
}

// Java method: CreateSyntheticBlocks.visit(com.google.javascript.jscomp.NodeTraversal unused, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *CreateSyntheticBlocks) Visit_3(unused *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = unused
	_ = n
	_ = parent
	if !(n.IsCall()) || !(n.GetFirstChild().IsName()) {
		return
	}
	var callTarget pkg_rhino.Node = n.GetFirstChild()
	_ = callTarget
	var callName string = callTarget.GetString()
	_ = callName
	if j2gort.CheckNotNullString(this_.startMarkerName) == callName {
		if !(parent.IsExprResult()) {
			this_.compiler.Report_1(JSError_Make_3(n, StaticGet_CreateSyntheticBlocks_INVALID_MARKER_USAGE(), []string{this_.startMarkerName}))
			return
		}
		this_.markerStack.Push(j2gort.CheckNotNullNode(parent))
		return
	}
	if !(j2gort.CheckNotNullString(this_.endMarkerName) == callName) {
		return
	}
	var endMarkerNode pkg_rhino.Node = parent
	_ = endMarkerNode
	if !(endMarkerNode.IsExprResult()) {
		this_.compiler.Report_1(JSError_Make_3(n, StaticGet_CreateSyntheticBlocks_INVALID_MARKER_USAGE(), []string{this_.endMarkerName}))
		return
	}
	if this_.markerStack.IsEmpty() {
		this_.compiler.Report_1(JSError_Make_3(n, StaticGet_CreateSyntheticBlocks_UNMATCHED_END_MARKER(), []string{this_.startMarkerName, this_.endMarkerName}))
		return
	}
	var startMarkerNode pkg_rhino.Node = this_.markerStack.Pop()
	_ = startMarkerNode
	if endMarkerNode.GetParent() != startMarkerNode.GetParent() {
		this_.compiler.Report_1(JSError_Make_3(n, StaticGet_CreateSyntheticBlocks_UNMATCHED_END_MARKER(), []string{this_.startMarkerName, this_.endMarkerName}))
		return
	}
	this_.validMarkers = append(this_.validMarkers, NewCreateSyntheticBlocks_Marker(startMarkerNode, endMarkerNode))
}

// Java class: com.google.javascript.jscomp.CreateSyntheticBlocks.Marker
type CreateSyntheticBlocks_Marker struct {
	startMarker pkg_rhino.Node `gson:"startMarker"`
	endMarker   pkg_rhino.Node `gson:"endMarker"`
	idHash_     atomic.Int32
}

func (this_ *CreateSyntheticBlocks_Marker) BindSelf_(s any) {
	_ = s
}
func (this_ *CreateSyntheticBlocks_Marker) IsJavaNil_() bool { return this_ == nil }
func (*CreateSyntheticBlocks_Marker) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CreateSyntheticBlocks$Marker", "Marker"
}
func (this_ *CreateSyntheticBlocks_Marker) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CreateSyntheticBlocks_Marker) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CreateSyntheticBlocks.Marker", "com.google.javascript.jscomp.CreateSyntheticBlocks$Marker", func(v_ any) bool {
		return j2gort.InstanceOf[*CreateSyntheticBlocks_Marker](v_)
	})
}

func (this_ *CreateSyntheticBlocks_Marker) initFields_CreateSyntheticBlocks_Marker() {
	_ = this_
}

// Java constructor: Marker(com.google.javascript.rhino.Node startMarker, com.google.javascript.rhino.Node endMarker)
func NewCreateSyntheticBlocks_Marker(startMarker pkg_rhino.Node, endMarker pkg_rhino.Node) *CreateSyntheticBlocks_Marker {
	this_ := &CreateSyntheticBlocks_Marker{}
	InitCreateSyntheticBlocks_Marker(this_, startMarker, endMarker)
	return this_
}

func InitCreateSyntheticBlocks_Marker(this_ *CreateSyntheticBlocks_Marker, startMarker pkg_rhino.Node, endMarker pkg_rhino.Node) {
	_ = startMarker
	_ = endMarker
	this_.initFields_CreateSyntheticBlocks_Marker()
	this_.startMarker = startMarker
	this_.endMarker = endMarker
}
func (*CreateSyntheticBlocks) AlwaysShouldTraverse_()      {}
func (*CreateSyntheticBlocks) CanReleaseSyntacticScopes_() {}
