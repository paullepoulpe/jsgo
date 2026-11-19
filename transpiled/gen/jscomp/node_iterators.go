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
	ancestors     []pkg_rhino.Node `gson:"ancestors"`
	ancestorCount int32            `gson:"ancestorCount"`
	idHash_       atomic.Int32
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
}

// Java constructor: FunctionlessLocalScope(com.google.javascript.rhino.Node a0, com.google.javascript.rhino.Node a1, com.google.javascript.rhino.Node a2)
func NewNodeIterators_FunctionlessLocalScope_3(a0 pkg_rhino.Node, a1 pkg_rhino.Node, a2 pkg_rhino.Node) *NodeIterators_FunctionlessLocalScope {
	this_ := &NodeIterators_FunctionlessLocalScope{}
	InitNodeIterators_FunctionlessLocalScope_3(this_, a0, a1, a2)
	return this_
}

func InitNodeIterators_FunctionlessLocalScope_3(this_ *NodeIterators_FunctionlessLocalScope, a0 pkg_rhino.Node, a1 pkg_rhino.Node, a2 pkg_rhino.Node) {
	_ = a0
	_ = a1
	_ = a2
	this_.initFields_NodeIterators_FunctionlessLocalScope()
	this_.ancestors = make([]pkg_rhino.Node, int32(16))
	var valid int32 = int32(0)
	_ = valid
	if !(a0.IsFunction()) {
		valid = int32(1)
		if !(a1.IsFunction()) {
			valid = int32(2)
			if !(a2.IsFunction()) {
				valid = int32(3)
			}
		}
	}
	this_.ancestorCount = valid
	if valid == int32(3) {
		{
			a_1 := this_.ancestors
			i_2 := int32(0)
			a_1[i_2] = a2
		}
		{
			a_3 := this_.ancestors
			i_4 := int32(1)
			a_3[i_4] = a1
		}
		{
			a_5 := this_.ancestors
			i_6 := int32(2)
			a_5[i_6] = a0
		}
	} else if valid == int32(2) {
		{
			a_7 := this_.ancestors
			i_8 := int32(0)
			a_7[i_8] = a1
		}
		{
			a_9 := this_.ancestors
			i_10 := int32(1)
			a_9[i_10] = a0
		}
	} else if valid == int32(1) {
		{
			a_11 := this_.ancestors
			i_12 := int32(0)
			a_11[i_12] = a0
		}
	}
}
func initReuse_NodeIterators_FunctionlessLocalScope_3(this_ *NodeIterators_FunctionlessLocalScope, a0 pkg_rhino.Node, a1 pkg_rhino.Node, a2 pkg_rhino.Node) {
	_ = a0
	_ = a1
	_ = a2
	this_.initFields_NodeIterators_FunctionlessLocalScope()
	if len(this_.ancestors) < 16 || len(this_.ancestors) > 64 {
		this_.ancestors = make([]pkg_rhino.Node, 16)
	} else {
		clear(this_.ancestors)
	}
	this_.idHash_.Store(0)
	var valid int32 = int32(0)
	_ = valid
	if !(a0.IsFunction()) {
		valid = int32(1)
		if !(a1.IsFunction()) {
			valid = int32(2)
			if !(a2.IsFunction()) {
				valid = int32(3)
			}
		}
	}
	this_.ancestorCount = valid
	if valid == int32(3) {
		{
			a_1 := this_.ancestors
			i_2 := int32(0)
			a_1[i_2] = a2
		}
		{
			a_3 := this_.ancestors
			i_4 := int32(1)
			a_3[i_4] = a1
		}
		{
			a_5 := this_.ancestors
			i_6 := int32(2)
			a_5[i_6] = a0
		}
	} else if valid == int32(2) {
		{
			a_7 := this_.ancestors
			i_8 := int32(0)
			a_7[i_8] = a1
		}
		{
			a_9 := this_.ancestors
			i_10 := int32(1)
			a_9[i_10] = a0
		}
	} else if valid == int32(1) {
		{
			a_11 := this_.ancestors
			i_12 := int32(0)
			a_11[i_12] = a0
		}
	}
}

// Java constructor: FunctionlessLocalScope(com.google.javascript.rhino.Node a0, com.google.javascript.rhino.Node a1, com.google.javascript.rhino.Node a2, com.google.javascript.rhino.Node a3)
func NewNodeIterators_FunctionlessLocalScope_4(a0 pkg_rhino.Node, a1 pkg_rhino.Node, a2 pkg_rhino.Node, a3 pkg_rhino.Node) *NodeIterators_FunctionlessLocalScope {
	this_ := &NodeIterators_FunctionlessLocalScope{}
	InitNodeIterators_FunctionlessLocalScope_4(this_, a0, a1, a2, a3)
	return this_
}

func InitNodeIterators_FunctionlessLocalScope_4(this_ *NodeIterators_FunctionlessLocalScope, a0 pkg_rhino.Node, a1 pkg_rhino.Node, a2 pkg_rhino.Node, a3 pkg_rhino.Node) {
	_ = a0
	_ = a1
	_ = a2
	_ = a3
	this_.initFields_NodeIterators_FunctionlessLocalScope()
	this_.ancestors = make([]pkg_rhino.Node, int32(16))
	var valid int32 = int32(0)
	_ = valid
	if !(a0.IsFunction()) {
		valid = int32(1)
		if !(a1.IsFunction()) {
			valid = int32(2)
			if !(a2.IsFunction()) {
				valid = int32(3)
				if !(a3.IsFunction()) {
					valid = int32(4)
				}
			}
		}
	}
	this_.ancestorCount = valid
	if valid == int32(4) {
		{
			a_13 := this_.ancestors
			i_14 := int32(0)
			a_13[i_14] = a3
		}
		{
			a_15 := this_.ancestors
			i_16 := int32(1)
			a_15[i_16] = a2
		}
		{
			a_17 := this_.ancestors
			i_18 := int32(2)
			a_17[i_18] = a1
		}
		{
			a_19 := this_.ancestors
			i_20 := int32(3)
			a_19[i_20] = a0
		}
	} else if valid == int32(3) {
		{
			a_21 := this_.ancestors
			i_22 := int32(0)
			a_21[i_22] = a2
		}
		{
			a_23 := this_.ancestors
			i_24 := int32(1)
			a_23[i_24] = a1
		}
		{
			a_25 := this_.ancestors
			i_26 := int32(2)
			a_25[i_26] = a0
		}
	} else if valid == int32(2) {
		{
			a_27 := this_.ancestors
			i_28 := int32(0)
			a_27[i_28] = a1
		}
		{
			a_29 := this_.ancestors
			i_30 := int32(1)
			a_29[i_30] = a0
		}
	} else if valid == int32(1) {
		{
			a_31 := this_.ancestors
			i_32 := int32(0)
			a_31[i_32] = a0
		}
	}
}

// Java constructor: FunctionlessLocalScope(com.google.javascript.rhino.Node[] ancestors)
func NewNodeIterators_FunctionlessLocalScope_1(ancestors []pkg_rhino.Node) *NodeIterators_FunctionlessLocalScope {
	this_ := &NodeIterators_FunctionlessLocalScope{}
	InitNodeIterators_FunctionlessLocalScope_1(this_, ancestors)
	return this_
}

func InitNodeIterators_FunctionlessLocalScope_1(this_ *NodeIterators_FunctionlessLocalScope, ancestors []pkg_rhino.Node) {
	_ = ancestors
	this_.initFields_NodeIterators_FunctionlessLocalScope()
	if !(int32(len(ancestors)) > int32(0)) {
		j2gort.CheckArgument(false)
	}
	var valid int32 = int32(0)
	_ = valid
	for (valid < int32(len(ancestors))) && !(ancestors[valid].IsFunction()) {
		valid = valid + 1
	}
	this_.ancestors = make([]pkg_rhino.Node, max(int32(16), valid))
	this_.ancestorCount = valid
	{
		var i int32 = int32(0)
		_ = i
		for ; i < valid; i = i + 1 {
			{
				a_33 := this_.ancestors
				i_34 := ((valid - int32(1)) - i)
				a_33[i_34] = ancestors[i]
			}
		}
	}
}

// Java method: FunctionlessLocalScope.pushAncestor(com.google.javascript.rhino.Node n)
func (this_ *NodeIterators_FunctionlessLocalScope) priv_FunctionlessLocalScope_pushAncestor_3928(n pkg_rhino.Node) {
	_ = n
	var c int32 = this_.ancestorCount
	_ = c
	var arr []pkg_rhino.Node = this_.ancestors
	_ = arr
	if c == int32(len(arr)) {
		{
			arrSrc_35 := arr
			arrLen_36 := (c * int32(2))
			j2gort.CheckNotNullSlice(arrSrc_35)
			if arrLen_36 < 0 {
				panic(j2gort.ConstructShimThrowable("java.lang.NegativeArraySizeException"))
			}
			arrOut_37 := make([]pkg_rhino.Node, arrLen_36)
			copy(arrOut_37, arrSrc_35)
			arr = arrOut_37
		}
		this_.ancestors = arr
	}
	{
		a_38 := arr
		i_39 := c
		a_38[i_39] = n
	}
	this_.ancestorCount = (c + int32(1))
}

// Java method: FunctionlessLocalScope.hasNext()
func (this_ *NodeIterators_FunctionlessLocalScope) HasNext() bool {
	var c int32 = this_.ancestorCount
	_ = c
	return !((c == int32(1)) && (this_.ancestors[int32(0)].GetNext() == 0))
}

// Java method: FunctionlessLocalScope.nextNode()
func (this_ *NodeIterators_FunctionlessLocalScope) NextNode() pkg_rhino.Node {
	for {
		var topIdx int32
		{
			old_40 := this_.ancestorCount
			nw_41 := old_40 - 1
			this_.ancestorCount = nw_41
			topIdx = nw_41
		}
		_ = topIdx
		var current pkg_rhino.Node = this_.ancestors[topIdx]
		_ = current
		{
			a_42 := this_.ancestors
			i_43 := topIdx
			a_42[i_43] = pkg_rhino.Node(0)
		}
		var nextSibling pkg_rhino.Node = current.GetNext()
		_ = nextSibling
		if nextSibling == 0 {
			current = this_.ancestors[(topIdx - int32(1))]
			if current.IsFunction() {
				continue
			}
			return current
		}
		current = nextSibling
		this_.priv_FunctionlessLocalScope_pushAncestor_3928(current)
		if current.IsFunction() {
			continue
		}
		var hitFunction bool = false
		_ = hitFunction
		for current.HasChildren() {
			current = current.GetFirstChild()
			this_.priv_FunctionlessLocalScope_pushAncestor_3928(current)
			if current.IsFunction() {
				hitFunction = true
				break
			}
		}
		if hitFunction {
			continue
		}
		return current
	}
	panic("unreachable")
}

// Java method: FunctionlessLocalScope.next()
func (this_ *NodeIterators_FunctionlessLocalScope) Next() any {
	return j2gort.BoxNode(this_.NextNode())
}

// Java method: FunctionlessLocalScope.remove()
func (this_ *NodeIterators_FunctionlessLocalScope) Remove() {
	panic(j2gort.AsIface[javalang.IUnsupportedOperationException](j2gort.ConstructShimThrowable("java.lang.UnsupportedOperationException", "Not implemented")))
}

// Java method: FunctionlessLocalScope.current()
func (this_ *NodeIterators_FunctionlessLocalScope) Current() pkg_rhino.Node {
	return this_.ancestors[(this_.ancestorCount - int32(1))]
}

// Java method: FunctionlessLocalScope.currentParent()
func (this_ *NodeIterators_FunctionlessLocalScope) CurrentParent() pkg_rhino.Node {
	var c int32 = this_.ancestorCount
	_ = c
	if c >= int32(2) {
		return this_.ancestors[(c - int32(1))].GetParent()
	} else {
		return pkg_rhino.Node(0)
	}
}

// Java method: FunctionlessLocalScope.currentAncestors()
func (this_ *NodeIterators_FunctionlessLocalScope) CurrentAncestors() j2gort.Ext_java_util_List {
	var c int32 = this_.ancestorCount
	_ = c
	var list []pkg_rhino.Node = j2gort.MakeNonNilSlice[pkg_rhino.Node](c)
	_ = list
	{
		var i int32 = (c - int32(1))
		_ = i
		for ; i >= int32(0); i = i - 1 {
			list = append(list, this_.ancestors[i])
		}
	}
	return j2gort.WrapTypedArrayList(list)
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
	this_.priv_LocalVarMotion_advanceLookAhead_3937(true)
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
	return NewNodeIterators_LocalVarMotion(compiler, name, NewNodeIterators_FunctionlessLocalScope_3(name, var_, block))
}
func allocReuseVia_NodeIterators_LocalVarMotion_ForVar(lvm_ *NodeIterators_LocalVarMotion, fls_ *NodeIterators_FunctionlessLocalScope, compiler AbstractCompiler, name pkg_rhino.Node, var_ pkg_rhino.Node, block pkg_rhino.Node) *NodeIterators_LocalVarMotion {
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
	initReuse_NodeIterators_FunctionlessLocalScope_3(fls_, name, var_, block)
	var zeroLvm_ NodeIterators_LocalVarMotion
	*lvm_ = zeroLvm_
	InitNodeIterators_LocalVarMotion(lvm_, compiler, name, fls_)
	return lvm_

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
	return NewNodeIterators_LocalVarMotion(compiler, name, NewNodeIterators_FunctionlessLocalScope_3(assign, expr, block))
}
func allocReuseVia_NodeIterators_LocalVarMotion_ForAssign(lvm_ *NodeIterators_LocalVarMotion, fls_ *NodeIterators_FunctionlessLocalScope, compiler AbstractCompiler, name pkg_rhino.Node, assign pkg_rhino.Node, expr pkg_rhino.Node, block pkg_rhino.Node) *NodeIterators_LocalVarMotion {
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
	initReuse_NodeIterators_FunctionlessLocalScope_3(fls_, assign, expr, block)
	var zeroLvm_ NodeIterators_LocalVarMotion
	*lvm_ = zeroLvm_
	InitNodeIterators_LocalVarMotion(lvm_, compiler, name, fls_)
	return lvm_

}

// Java method: LocalVarMotion.hasNext()
func (this_ *NodeIterators_LocalVarMotion) HasNext() bool {
	return (this_.lookAhead != 0)
}

// Java method: LocalVarMotion.nextNode()
func (this_ *NodeIterators_LocalVarMotion) NextNode() pkg_rhino.Node {
	var next pkg_rhino.Node = this_.lookAhead
	_ = next
	this_.priv_LocalVarMotion_advanceLookAhead_3937(false)
	return next
}

// Java method: LocalVarMotion.next()
func (this_ *NodeIterators_LocalVarMotion) Next() any {
	return j2gort.BoxNode(this_.NextNode())
}

// Java method: LocalVarMotion.remove()
func (this_ *NodeIterators_LocalVarMotion) Remove() {
	panic(j2gort.AsIface[javalang.IUnsupportedOperationException](j2gort.ConstructShimThrowable("java.lang.UnsupportedOperationException", "Not implemented")))
}

// Java method: LocalVarMotion.advanceLookAhead(boolean atStart)
func (this_ *NodeIterators_LocalVarMotion) priv_LocalVarMotion_advanceLookAhead_3937(atStart bool) {
	_ = atStart
	if !(atStart) {
		if this_.lookAhead == 0 {
			return
		}
		var curNode pkg_rhino.Node = this_.iterator.Current()
		_ = curNode
		if curNode.IsName() && func() bool {
			sr_44 := this_.varName
			se_45 := curNode.GetString()
			return j2gort.CheckNotNullString(sr_44) == se_45
		}() {
			this_.lookAhead = pkg_rhino.Node(0)
			return
		}
	}
	if !(this_.iterator.HasNext()) {
		this_.lookAhead = pkg_rhino.Node(0)
		return
	}
	var nextNode pkg_rhino.Node = this_.iterator.NextNode()
	_ = nextNode
	var nextParent pkg_rhino.Node = this_.iterator.CurrentParent()
	_ = nextParent
	var type_ pkg_rhino.Token = nextNode.GetToken()
	_ = type_
	if this_.valueHasSideEffects {
		var readsState bool = false
		_ = readsState
		if (nextNode.IsName() && !(func() bool {
			sr_46 := this_.varName
			se_47 := nextNode.GetString()
			return j2gort.CheckNotNullString(sr_46) == se_47
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
