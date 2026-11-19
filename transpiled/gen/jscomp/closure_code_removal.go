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
// Port of com/google/javascript/jscomp/ClosureCodeRemoval.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.ClosureCodeRemoval
type ClosureCodeRemoval struct {
	compiler                      AbstractCompiler                          `gson:"compiler"`
	abstractMethodAssignmentNodes []*ClosureCodeRemoval_RemovableAssignment `gson:"abstractMethodAssignmentNodes"`
	abstractMemberFunctionNodes   []pkg_rhino.Node                          `gson:"abstractMemberFunctionNodes"`
	assertionCalls                []pkg_rhino.Node                          `gson:"assertionCalls"`
	idHash_                       atomic.Int32
	removeAbstractMethods         bool `gson:"removeAbstractMethods"`
	removeAssertionCalls          bool `gson:"removeAssertionCalls"`
}

func (this_ *ClosureCodeRemoval) BindSelf_(s any) {
	_ = s
}
func (this_ *ClosureCodeRemoval) IsJavaNil_() bool { return this_ == nil }
func (*ClosureCodeRemoval) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ClosureCodeRemoval", "ClosureCodeRemoval"
}
func (this_ *ClosureCodeRemoval) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ClosureCodeRemoval) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ClosureCodeRemoval", "com.google.javascript.jscomp.ClosureCodeRemoval", func(v_ any) bool {
		return j2gort.InstanceOf[*ClosureCodeRemoval](v_)
	})
}
func (*ClosureCodeRemoval) Implements_jscomp_CompilerPass() {}

func StaticGet_ClosureCodeRemoval_ABSTRACT_METHOD_NAME() string { return "goog.abstractMethod" }
func (this_ *ClosureCodeRemoval) initFields_ClosureCodeRemoval() {
	_ = this_
	this_.abstractMethodAssignmentNodes = []*ClosureCodeRemoval_RemovableAssignment{}
	this_.abstractMemberFunctionNodes = []pkg_rhino.Node{}
	this_.assertionCalls = []pkg_rhino.Node{}
}

// Java constructor: ClosureCodeRemoval(com.google.javascript.jscomp.AbstractCompiler compiler, boolean removeAbstractMethods, boolean removeAssertionCalls)
func NewClosureCodeRemoval(compiler AbstractCompiler, removeAbstractMethods bool, removeAssertionCalls bool) *ClosureCodeRemoval {
	this_ := &ClosureCodeRemoval{}
	InitClosureCodeRemoval(this_, compiler, removeAbstractMethods, removeAssertionCalls)
	return this_
}

func InitClosureCodeRemoval(this_ *ClosureCodeRemoval, compiler AbstractCompiler, removeAbstractMethods bool, removeAssertionCalls bool) {
	_ = compiler
	_ = removeAbstractMethods
	_ = removeAssertionCalls
	this_.initFields_ClosureCodeRemoval()
	this_.compiler = compiler
	this_.removeAbstractMethods = removeAbstractMethods
	this_.removeAssertionCalls = removeAssertionCalls
}

// Java method: ClosureCodeRemoval.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *ClosureCodeRemoval) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	var passes j2gort.Ext_java_util_List = j2gort.NewExt_java_util_ArrayList()
	_ = passes
	if this_.removeAbstractMethods {
		passes.Add_1(NewClosureCodeRemoval_FindAbstractMethods(this_))
	}
	if this_.removeAssertionCalls {
		passes.Add_1(NewClosureCodeRemoval_FindAssertionCalls(this_))
	}
	CombinedCompilerPass_Traverse(this_.compiler, root, passes)
	for _, assignment := range this_.abstractMethodAssignmentNodes {
		_ = assignment
		assignment.Remove()
	}

	for _, memberFunction := range this_.abstractMemberFunctionNodes {
		_ = memberFunction
		this_.compiler.ReportFunctionDeleted(memberFunction.GetFirstChild())
		var parent pkg_rhino.Node = memberFunction.GetParent()
		_ = parent
		memberFunction.Detach()
		this_.compiler.ReportChangeToEnclosingScope(parent)
	}

	for _, call := range this_.assertionCalls {
		_ = call
		this_.compiler.ReportChangeToEnclosingScope(call)
		var parent_2 pkg_rhino.Node = call.GetParent()
		_ = parent_2
		if parent_2.IsExprResult() {
			parent_2.Detach()
			NodeUtil_MarkFunctionsDeleted_2_Node_AbstractCompiler(parent_2, this_.compiler)
		} else {
			var firstArg pkg_rhino.Node = call.GetSecondChild()
			_ = firstArg
			if firstArg == 0 {
				call.ReplaceWith(NodeUtil_NewUndefinedNode(call))
			} else {
				var replacement pkg_rhino.Node = firstArg.Detach()
				_ = replacement
				replacement.CopyTypeFrom(call)
				call.ReplaceWith(replacement)
			}
			NodeUtil_MarkFunctionsDeleted_2_Node_AbstractCompiler(call, this_.compiler)
		}
	}

}

// Java class: com.google.javascript.jscomp.ClosureCodeRemoval.RemovableAssignment
type ClosureCodeRemoval_RemovableAssignment struct {
	outer_          *ClosureCodeRemoval
	assignAncestors []pkg_rhino.Node `gson:"assignAncestors"`
	node            pkg_rhino.Node   `gson:"node"`
	parent          pkg_rhino.Node   `gson:"parent"`
	lastAncestor    pkg_rhino.Node   `gson:"lastAncestor"`
	idHash_         atomic.Int32
}

func (this_ *ClosureCodeRemoval_RemovableAssignment) BindSelf_(s any) {
	_ = s
}
func (this_ *ClosureCodeRemoval_RemovableAssignment) IsJavaNil_() bool { return this_ == nil }
func (*ClosureCodeRemoval_RemovableAssignment) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ClosureCodeRemoval$RemovableAssignment", "RemovableAssignment"
}
func (this_ *ClosureCodeRemoval_RemovableAssignment) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ClosureCodeRemoval_RemovableAssignment) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ClosureCodeRemoval.RemovableAssignment", "com.google.javascript.jscomp.ClosureCodeRemoval$RemovableAssignment", func(v_ any) bool {
		return j2gort.InstanceOf[*ClosureCodeRemoval_RemovableAssignment](v_)
	})
}
func (this_ *ClosureCodeRemoval_RemovableAssignment) SetOuter_(o *ClosureCodeRemoval) {
	this_.outer_ = o
}
func (this_ *ClosureCodeRemoval_RemovableAssignment) GetOuter_() *ClosureCodeRemoval {
	return this_.outer_
}

func (this_ *ClosureCodeRemoval_RemovableAssignment) initFields_ClosureCodeRemoval_RemovableAssignment() {
	_ = this_
	this_.assignAncestors = []pkg_rhino.Node{}
}

// Java constructor: RemovableAssignment(com.google.javascript.rhino.Node nameNode, com.google.javascript.rhino.Node assignNode)
func NewClosureCodeRemoval_RemovableAssignment(outer_ *ClosureCodeRemoval, nameNode pkg_rhino.Node, assignNode pkg_rhino.Node) *ClosureCodeRemoval_RemovableAssignment {
	this_ := &ClosureCodeRemoval_RemovableAssignment{}
	this_.outer_ = outer_
	InitClosureCodeRemoval_RemovableAssignment(this_, nameNode, assignNode)
	return this_
}

func InitClosureCodeRemoval_RemovableAssignment(this_ *ClosureCodeRemoval_RemovableAssignment, nameNode pkg_rhino.Node, assignNode pkg_rhino.Node) {
	_ = nameNode
	_ = assignNode
	this_.initFields_ClosureCodeRemoval_RemovableAssignment()
	this_.node = nameNode
	this_.parent = assignNode
	var ancestor pkg_rhino.Node = assignNode
	_ = ancestor
	for doFirst_5 := true; doFirst_5 || (ancestor.IsAssign() && ancestor.GetFirstChild().IsQualifiedName()); doFirst_5 = false {
		ancestor = ancestor.GetParent()
		this_.assignAncestors = append(this_.assignAncestors, ancestor)
	}
	this_.lastAncestor = ancestor.GetParent()
}

// Java method: RemovableAssignment.remove()
func (this_ *ClosureCodeRemoval_RemovableAssignment) Remove() {
	var rhs pkg_rhino.Node = this_.node.GetNext()
	_ = rhs
	var last pkg_rhino.Node = this_.parent
	_ = last
	for _, ancestor := range this_.assignAncestors {
		_ = ancestor
		if ancestor.IsExprResult() {
			ancestor.Detach()
			NodeUtil_MarkFunctionsDeleted_2_Node_AbstractCompiler(ancestor, this_.outer_.compiler)
		} else {
			rhs.Detach()
			last.ReplaceWith(rhs)
		}
		last = ancestor
	}

	this_.outer_.compiler.ReportChangeToEnclosingScope(this_.lastAncestor)
}

// Java class: com.google.javascript.jscomp.ClosureCodeRemoval.FindAbstractMethods
type ClosureCodeRemoval_FindAbstractMethods struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer_ *ClosureCodeRemoval
}

func (this_ *ClosureCodeRemoval_FindAbstractMethods) BindSelf_(s any) {
	_ = s
}
func (this_ *ClosureCodeRemoval_FindAbstractMethods) IsJavaNil_() bool { return this_ == nil }
func (*ClosureCodeRemoval_FindAbstractMethods) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ClosureCodeRemoval$FindAbstractMethods", "FindAbstractMethods"
}
func (this_ *ClosureCodeRemoval_FindAbstractMethods) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*ClosureCodeRemoval_FindAbstractMethods) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ClosureCodeRemoval.FindAbstractMethods", "com.google.javascript.jscomp.ClosureCodeRemoval$FindAbstractMethods", func(v_ any) bool {
		return j2gort.InstanceOf[*ClosureCodeRemoval_FindAbstractMethods](v_)
	})
}
func (this_ *ClosureCodeRemoval_FindAbstractMethods) SetOuter_(o *ClosureCodeRemoval) {
	this_.outer_ = o
}
func (this_ *ClosureCodeRemoval_FindAbstractMethods) GetOuter_() *ClosureCodeRemoval {
	return this_.outer_
}
func (*ClosureCodeRemoval_FindAbstractMethods) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {
}
func (*ClosureCodeRemoval_FindAbstractMethods) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *ClosureCodeRemoval_FindAbstractMethods) initFields_ClosureCodeRemoval_FindAbstractMethods() {
	_ = this_
}

// Java constructor: FindAbstractMethods()
func NewClosureCodeRemoval_FindAbstractMethods(outer_ *ClosureCodeRemoval) *ClosureCodeRemoval_FindAbstractMethods {
	this_ := &ClosureCodeRemoval_FindAbstractMethods{}
	this_.outer_ = outer_
	InitClosureCodeRemoval_FindAbstractMethods(this_)
	return this_
}

func InitClosureCodeRemoval_FindAbstractMethods(this_ *ClosureCodeRemoval_FindAbstractMethods) {
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_ClosureCodeRemoval_FindAbstractMethods()
}

// Java method: FindAbstractMethods.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *ClosureCodeRemoval_FindAbstractMethods) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if n.IsAssign() {
		var nameNode pkg_rhino.Node = n.GetFirstChild()
		_ = nameNode
		var valueNode pkg_rhino.Node = n.GetLastChild()
		_ = valueNode
		if (nameNode.IsQualifiedName() && valueNode.IsQualifiedName()) && valueNode.MatchesQualifiedName_1_string("goog.abstractMethod") {
			this_.outer_.abstractMethodAssignmentNodes = append(this_.outer_.abstractMethodAssignmentNodes, NewClosureCodeRemoval_RemovableAssignment(this_.outer_, n.GetFirstChild(), n))
		} else if (((n.GetJSDocInfo() != nil) && n.GetJSDocInfo().IsAbstract()) && NodeUtil_IsEmptyFunctionExpression(valueNode)) && !(n.GetJSDocInfo().IsConstructor()) {
			this_.outer_.abstractMethodAssignmentNodes = append(this_.outer_.abstractMethodAssignmentNodes, NewClosureCodeRemoval_RemovableAssignment(this_.outer_, n.GetFirstChild(), n))
		}
	} else if n.IsMemberFunctionDef() && parent.IsClassMembers() {
		if (n.GetJSDocInfo() != nil) && n.GetJSDocInfo().IsAbstract() {
			this_.outer_.abstractMemberFunctionNodes = append(this_.outer_.abstractMemberFunctionNodes, n)
		}
	}
}

// Java class: com.google.javascript.jscomp.ClosureCodeRemoval.FindAssertionCalls
type ClosureCodeRemoval_FindAssertionCalls struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer_         *ClosureCodeRemoval
	assertionNames *CodingConvention_AssertionFunctionLookup `gson:"assertionNames"`
}

func (this_ *ClosureCodeRemoval_FindAssertionCalls) BindSelf_(s any) {
	_ = s
}
func (this_ *ClosureCodeRemoval_FindAssertionCalls) IsJavaNil_() bool { return this_ == nil }
func (*ClosureCodeRemoval_FindAssertionCalls) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ClosureCodeRemoval$FindAssertionCalls", "FindAssertionCalls"
}
func (this_ *ClosureCodeRemoval_FindAssertionCalls) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*ClosureCodeRemoval_FindAssertionCalls) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ClosureCodeRemoval.FindAssertionCalls", "com.google.javascript.jscomp.ClosureCodeRemoval$FindAssertionCalls", func(v_ any) bool {
		return j2gort.InstanceOf[*ClosureCodeRemoval_FindAssertionCalls](v_)
	})
}
func (this_ *ClosureCodeRemoval_FindAssertionCalls) SetOuter_(o *ClosureCodeRemoval) {
	this_.outer_ = o
}
func (this_ *ClosureCodeRemoval_FindAssertionCalls) GetOuter_() *ClosureCodeRemoval {
	return this_.outer_
}
func (*ClosureCodeRemoval_FindAssertionCalls) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {
}
func (*ClosureCodeRemoval_FindAssertionCalls) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *ClosureCodeRemoval_FindAssertionCalls) initFields_ClosureCodeRemoval_FindAssertionCalls() {
	_ = this_
}

// Java constructor: FindAssertionCalls()
func NewClosureCodeRemoval_FindAssertionCalls(outer_ *ClosureCodeRemoval) *ClosureCodeRemoval_FindAssertionCalls {
	this_ := &ClosureCodeRemoval_FindAssertionCalls{}
	this_.outer_ = outer_
	InitClosureCodeRemoval_FindAssertionCalls(this_)
	return this_
}

func InitClosureCodeRemoval_FindAssertionCalls(this_ *ClosureCodeRemoval_FindAssertionCalls) {
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_ClosureCodeRemoval_FindAssertionCalls()
	this_.assertionNames = CodingConvention_AssertionFunctionLookup_Of(j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.outer_.compiler.GetCodingConvention().GetAssertionFunctions()))
}

// Java method: FindAssertionCalls.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *ClosureCodeRemoval_FindAssertionCalls) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if !(n.IsCall()) {
		return
	}
	var callee pkg_rhino.Node = n.GetFirstChild()
	_ = callee
	if !(j2gort.IfaceEquals[CodingConvention_AssertionFunctionSpec](this_.assertionNames.LookupByCallee(callee), nil)) || ((callee.GetColor() != nil) && callee.GetColor().IsClosureAssert()) {
		this_.outer_.assertionCalls = append(this_.outer_.assertionCalls, n)
	}
}
func (*ClosureCodeRemoval_FindAbstractMethods) AlwaysShouldTraverse_()      {}
func (*ClosureCodeRemoval_FindAbstractMethods) CanReleaseSyntacticScopes_() {}
func (*ClosureCodeRemoval_FindAssertionCalls) AlwaysShouldTraverse_()       {}
func (*ClosureCodeRemoval_FindAssertionCalls) CanReleaseSyntacticScopes_()  {}
