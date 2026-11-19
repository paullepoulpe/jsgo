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
// Port of com/google/javascript/jscomp/NodeIterators.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.NodeIterators
type NodeIterators struct {
	idHash_ atomic.Int32
}

func (this_ *NodeIterators) BindSelf_(s any) {
	_ = s
}
func (this_ *NodeIterators) IsJavaNil_() bool { return this_ == nil }
func (*NodeIterators) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.NodeIterators", "NodeIterators"
}
func (this_ *NodeIterators) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*NodeIterators) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.NodeIterators", "com.google.javascript.jscomp.NodeIterators", func(v_ any) bool {
		return j2gort.InstanceOf[*NodeIterators](v_)
	})
}

func (this_ *NodeIterators) initFields_NodeIterators() {
	_ = this_
}

// Java constructor: NodeIterators()
func NewNodeIterators() *NodeIterators {
	this_ := &NodeIterators{}
	InitNodeIterators(this_)
	return this_
}

func InitNodeIterators(this_ *NodeIterators) {
	this_.initFields_NodeIterators()
}

// Java class: com.google.javascript.jscomp.NodeIterators.FunctionlessLocalScope
type NodeIterators_FunctionlessLocalScope struct {
	ancestors j2gort.Ext_java_util_Deque `gson:"ancestors"`
	idHash_   atomic.Int32
}

func (this_ *NodeIterators_FunctionlessLocalScope) BindSelf_(s any) {
	_ = s
}
func (this_ *NodeIterators_FunctionlessLocalScope) IsJavaNil_() bool { return this_ == nil }
func (*NodeIterators_FunctionlessLocalScope) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.NodeIterators$FunctionlessLocalScope", "FunctionlessLocalScope"
}
func (this_ *NodeIterators_FunctionlessLocalScope) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*NodeIterators_FunctionlessLocalScope) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.NodeIterators.FunctionlessLocalScope", "com.google.javascript.jscomp.NodeIterators$FunctionlessLocalScope", func(v_ any) bool {
		return j2gort.InstanceOf[*NodeIterators_FunctionlessLocalScope](v_)
	})
}
func (*NodeIterators_FunctionlessLocalScope) Implements_j2gort_Ext_java_util_Iterator() {}

func (this_ *NodeIterators_FunctionlessLocalScope) initFields_NodeIterators_FunctionlessLocalScope() {
	_ = this_
	this_.ancestors = j2gort.NewExt_java_util_ArrayDeque()
}

// Java constructor: FunctionlessLocalScope(com.google.javascript.rhino.Node[] ancestors)
func NewNodeIterators_FunctionlessLocalScope(ancestors []pkg_rhino.Node) *NodeIterators_FunctionlessLocalScope {
	this_ := &NodeIterators_FunctionlessLocalScope{}
	InitNodeIterators_FunctionlessLocalScope(this_, ancestors)
	return this_
}

func InitNodeIterators_FunctionlessLocalScope(this_ *NodeIterators_FunctionlessLocalScope, ancestors []pkg_rhino.Node) {
	_ = ancestors
	this_.initFields_NodeIterators_FunctionlessLocalScope()
	if !(int32(len(ancestors)) > int32(0)) {
		j2gort.CheckArgument(false)
	}
	for _, el_1 := range j2gort.CheckNotNullSlice(ancestors) {
		n := el_1
		_ = n
		if n.IsFunction() {
			break
		}
		this_.ancestors.AddFirst(j2gort.BoxNode(n))
	}
}

// Java method: FunctionlessLocalScope.hasNext()
func (this_ *NodeIterators_FunctionlessLocalScope) HasNext() bool {
	return !((this_.ancestors.Size() == int32(1)) && (pkg_rhino.Node(j2gort.DequePeekLastNode(this_.ancestors)).GetNext() == 0))
}

// Java method: FunctionlessLocalScope.next()
func (this_ *NodeIterators_FunctionlessLocalScope) Next() any {
	var current pkg_rhino.Node = pkg_rhino.Node(j2gort.DequeRemoveLastNode(this_.ancestors))
	_ = current
	if current.GetNext() == 0 {
		current = pkg_rhino.Node(j2gort.DequePeekLastNode(this_.ancestors))
		if current.IsFunction() {
			return j2gort.BoxNode(j2gort.CastNode[pkg_rhino.Node](this_.Next()))
		}
	} else {
		current = current.GetNext()
		j2gort.DequeAddLastNode(this_.ancestors, uint32(current))
		if current.IsFunction() {
			return j2gort.BoxNode(j2gort.CastNode[pkg_rhino.Node](this_.Next()))
		}
		for current.HasChildren() {
			current = current.GetFirstChild()
			j2gort.DequeAddLastNode(this_.ancestors, uint32(current))
			if current.IsFunction() {
				return j2gort.BoxNode(j2gort.CastNode[pkg_rhino.Node](this_.Next()))
			}
		}
	}
	return j2gort.BoxNode(current)
}

// Java method: FunctionlessLocalScope.remove()
func (this_ *NodeIterators_FunctionlessLocalScope) Remove() {
	panic(j2gort.AsIface[javalang.IUnsupportedOperationException](j2gort.ConstructShimThrowable("java.lang.UnsupportedOperationException", "Not implemented")))
}

// Java method: FunctionlessLocalScope.current()
func (this_ *NodeIterators_FunctionlessLocalScope) Current() pkg_rhino.Node {
	return pkg_rhino.Node(j2gort.DequePeekLastNode(this_.ancestors))
}

// Java method: FunctionlessLocalScope.currentParent()
func (this_ *NodeIterators_FunctionlessLocalScope) CurrentParent() pkg_rhino.Node {
	if this_.ancestors.Size() >= int32(2) {
		return this_.Current().GetParent()
	} else {
		return pkg_rhino.Node(0)
	}
}

// Java method: FunctionlessLocalScope.currentAncestors()
func (this_ *NodeIterators_FunctionlessLocalScope) CurrentAncestors() j2gort.Ext_java_util_List {
	var list j2gort.Ext_java_util_List = j2gort.NewExt_java_util_ArrayList_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.ancestors))
	_ = list
	j2gort.Ext_java_util_Collections_Reverse(list)
	return list
}

// Java class: com.google.javascript.jscomp.NodeIterators.LocalVarMotion
type NodeIterators_LocalVarMotion struct {
	compiler            AbstractCompiler                      `gson:"compiler"`
	iterator            *NodeIterators_FunctionlessLocalScope `gson:"iterator"`
	varName             string                                `gson:"varName"`
	lookAhead           pkg_rhino.Node                        `gson:"lookAhead"`
	idHash_             atomic.Int32
	valueHasSideEffects bool `gson:"valueHasSideEffects"`
}

func (this_ *NodeIterators_LocalVarMotion) BindSelf_(s any) {
	_ = s
}
func (this_ *NodeIterators_LocalVarMotion) IsJavaNil_() bool { return this_ == nil }
func (*NodeIterators_LocalVarMotion) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.NodeIterators$LocalVarMotion", "LocalVarMotion"
}
func (this_ *NodeIterators_LocalVarMotion) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*NodeIterators_LocalVarMotion) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.NodeIterators.LocalVarMotion", "com.google.javascript.jscomp.NodeIterators$LocalVarMotion", func(v_ any) bool {
		return j2gort.InstanceOf[*NodeIterators_LocalVarMotion](v_)
	})
}
func (*NodeIterators_LocalVarMotion) Implements_j2gort_Ext_java_util_Iterator() {}

func (this_ *NodeIterators_LocalVarMotion) initFields_NodeIterators_LocalVarMotion() {
	_ = this_
	this_.varName = j2gort.NullString
}

// Java constructor: LocalVarMotion(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.rhino.Node nameNode, com.google.javascript.jscomp.NodeIterators.FunctionlessLocalScope iterator)
func NewNodeIterators_LocalVarMotion(compiler AbstractCompiler, nameNode pkg_rhino.Node, iterator *NodeIterators_FunctionlessLocalScope) *NodeIterators_LocalVarMotion {
	this_ := &NodeIterators_LocalVarMotion{}
	InitNodeIterators_LocalVarMotion(this_, compiler, nameNode, iterator)
	return this_
}

func InitNodeIterators_LocalVarMotion(this_ *NodeIterators_LocalVarMotion, compiler AbstractCompiler, nameNode pkg_rhino.Node, iterator *NodeIterators_FunctionlessLocalScope) {
	_ = compiler
	_ = nameNode
	_ = iterator
	this_.initFields_NodeIterators_LocalVarMotion()
	if !(nameNode.IsName()) {
		j2gort.CheckArgument(false)
	}
	var valueNode pkg_rhino.Node = NodeUtil_GetAssignedValue(nameNode)
	_ = valueNode
	this_.compiler = j2gort.CheckNotNullIface(compiler)
	this_.varName = nameNode.GetString()
	this_.valueHasSideEffects = ((valueNode != 0) && compiler.GetAstAnalyzer().MayHaveSideEffects(valueNode))
	this_.iterator = iterator
	this_.priv_LocalVarMotion_advanceLookAhead_3930(true)
}

// Java static method: LocalVarMotion.forVar(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.rhino.Node name, com.google.javascript.rhino.Node var, com.google.javascript.rhino.Node block)
func NodeIterators_LocalVarMotion_ForVar(compiler AbstractCompiler, name pkg_rhino.Node, var_ pkg_rhino.Node, block pkg_rhino.Node) *NodeIterators_LocalVarMotion {
	_ = compiler
	_ = name
	_ = var_
	_ = block
	if !(NodeUtil_IsNameDeclaration(var_)) {
		j2gort.CheckArgument(false)
	}
	if !(NodeUtil_IsStatement(var_)) {
		j2gort.CheckArgument(false)
	}
	return NewNodeIterators_LocalVarMotion(compiler, name, NewNodeIterators_FunctionlessLocalScope([]pkg_rhino.Node{name, var_, block}))
}

// Java static method: LocalVarMotion.forAssign(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.rhino.Node name, com.google.javascript.rhino.Node assign, com.google.javascript.rhino.Node expr, com.google.javascript.rhino.Node block)
func NodeIterators_LocalVarMotion_ForAssign(compiler AbstractCompiler, name pkg_rhino.Node, assign pkg_rhino.Node, expr pkg_rhino.Node, block pkg_rhino.Node) *NodeIterators_LocalVarMotion {
	_ = compiler
	_ = name
	_ = assign
	_ = expr
	_ = block
	if !(assign.IsAssign()) {
		j2gort.CheckArgument(false)
	}
	if !(expr.IsExprResult()) {
		j2gort.CheckArgument(false)
	}
	return NewNodeIterators_LocalVarMotion(compiler, name, NewNodeIterators_FunctionlessLocalScope([]pkg_rhino.Node{assign, expr, block}))
}

// Java method: LocalVarMotion.hasNext()
func (this_ *NodeIterators_LocalVarMotion) HasNext() bool {
	return (this_.lookAhead != 0)
}

// Java method: LocalVarMotion.next()
func (this_ *NodeIterators_LocalVarMotion) Next() any {
	var next pkg_rhino.Node = this_.lookAhead
	_ = next
	this_.priv_LocalVarMotion_advanceLookAhead_3930(false)
	return j2gort.BoxNode(next)
}

// Java method: LocalVarMotion.remove()
func (this_ *NodeIterators_LocalVarMotion) Remove() {
	panic(j2gort.AsIface[javalang.IUnsupportedOperationException](j2gort.ConstructShimThrowable("java.lang.UnsupportedOperationException", "Not implemented")))
}

// Java method: LocalVarMotion.advanceLookAhead(boolean atStart)
func (this_ *NodeIterators_LocalVarMotion) priv_LocalVarMotion_advanceLookAhead_3930(atStart bool) {
	_ = atStart
	if !(atStart) {
		if this_.lookAhead == 0 {
			return
		}
		var curNode pkg_rhino.Node = this_.iterator.Current()
		_ = curNode
		if curNode.IsName() && func() bool {
			sr_2 := this_.varName
			se_3 := curNode.GetString()
			return j2gort.CheckNotNullString(sr_2) == se_3
		}() {
			this_.lookAhead = pkg_rhino.Node(0)
			return
		}
	}
	if !(this_.iterator.HasNext()) {
		this_.lookAhead = pkg_rhino.Node(0)
		return
	}
	var nextNode pkg_rhino.Node = j2gort.CastNode[pkg_rhino.Node](this_.iterator.Next())
	_ = nextNode
	var nextParent pkg_rhino.Node = this_.iterator.CurrentParent()
	_ = nextParent
	var type_ pkg_rhino.Token = nextNode.GetToken()
	_ = type_
	if this_.valueHasSideEffects {
		var readsState bool = false
		_ = readsState
		if (nextNode.IsName() && !(func() bool {
			sr_4 := this_.varName
			se_5 := nextNode.GetString()
			return j2gort.CheckNotNullString(sr_4) == se_5
		}())) || (((nextNode.IsGetProp() || nextNode.IsGetElem()) || nextNode.IsOptChainGetProp()) || nextNode.IsOptChainGetElem()) {
			if (nextParent == 0) || !(NodeUtil_IsNameDeclOrSimpleAssignLhs(nextNode, nextParent)) {
				readsState = true
			}
		} else if (nextNode.IsCall() || nextNode.IsNew()) || nextNode.IsOptChainCall() {
			readsState = true
		}
		if readsState {
			this_.lookAhead = pkg_rhino.Node(0)
			return
		}
	}
	if (this_.compiler.GetAstAnalyzer().NodeTypeMayHaveSideEffects(nextNode) && (type_ != pkg_rhino.StaticGet_Token_NAME())) || ((type_ == pkg_rhino.StaticGet_Token_NAME()) && nextParent.IsCatch()) {
		this_.lookAhead = pkg_rhino.Node(0)
		return
	}
	this_.lookAhead = nextNode
}
