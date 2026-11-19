/*
 * Copyright 2013 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/MinimizedCondition.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.MinimizedCondition
type MinimizedCondition struct {
	positive *MinimizedCondition_MeasuredNode `gson:"positive"`
	negative *MinimizedCondition_MeasuredNode `gson:"negative"`
	idHash_  atomic.Int32
}

func (this_ *MinimizedCondition) BindSelf_(s any) {
	_ = s
}
func (this_ *MinimizedCondition) IsJavaNil_() bool { return this_ == nil }
func (*MinimizedCondition) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.MinimizedCondition", "MinimizedCondition"
}
func (this_ *MinimizedCondition) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*MinimizedCondition) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.MinimizedCondition", "com.google.javascript.jscomp.MinimizedCondition", func(v_ any) bool {
		return j2gort.InstanceOf[*MinimizedCondition](v_)
	})
}

var static_MinimizedCondition_UNOPTIMIZED_NEG *MinimizedCondition_MeasuredNode

func StaticGet_MinimizedCondition_UNOPTIMIZED_NEG() *MinimizedCondition_MeasuredNode {
	clinit_MinimizedCondition()
	return static_MinimizedCondition_UNOPTIMIZED_NEG
}
func StaticSet_MinimizedCondition_UNOPTIMIZED_NEG(v *MinimizedCondition_MeasuredNode) {
	clinit_MinimizedCondition()
	static_MinimizedCondition_UNOPTIMIZED_NEG = v
}

var clinitDone_MinimizedCondition uint32
var clinitState_MinimizedCondition j2gort.ClinitState

func clinit_MinimizedCondition() {
	if atomic.LoadUint32(&clinitDone_MinimizedCondition) != 0 {
		return
	}
	clinitSlow_MinimizedCondition()
}

func clinitSlow_MinimizedCondition() {
	if !j2gort.EnterClinit(&clinitState_MinimizedCondition, &clinitDone_MinimizedCondition) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_MinimizedCondition, &clinitDone_MinimizedCondition)
	static_MinimizedCondition_UNOPTIMIZED_NEG = NewMinimizedCondition_MeasuredNode(pkg_rhino.Node(0), int32(-1), nil, nil, nil, int32(2147483647), true)
}

func (this_ *MinimizedCondition) initFields_MinimizedCondition() {
	_ = this_
}

// Java constructor: MinimizedCondition(com.google.javascript.jscomp.MinimizedCondition.MeasuredNode p, com.google.javascript.jscomp.MinimizedCondition.MeasuredNode n)
func NewMinimizedCondition(p *MinimizedCondition_MeasuredNode, n *MinimizedCondition_MeasuredNode) *MinimizedCondition {
	clinit_MinimizedCondition()
	this_ := &MinimizedCondition{}
	InitMinimizedCondition(this_, p, n)
	return this_
}
func allocPooled_MinimizedCondition(slabMeasured_ *j2gort.PooledSlab[MinimizedCondition_MeasuredNode], slabCond_ *j2gort.PooledSlab[MinimizedCondition], p *MinimizedCondition_MeasuredNode, n *MinimizedCondition_MeasuredNode) *MinimizedCondition {
	clinit_MinimizedCondition()
	this_ := slabCond_.Alloc()

	initMinimizedConditionWithSlab_(slabMeasured_, this_, p, n)
	return this_
}

func InitMinimizedCondition(this_ *MinimizedCondition, p *MinimizedCondition_MeasuredNode, n *MinimizedCondition_MeasuredNode) {
	_ = p
	_ = n
	this_.initFields_MinimizedCondition()
	this_.positive = p
	this_.negative = n.priv_MeasuredNode_change_3877()
}
func initMinimizedConditionWithSlab_(slab_ *j2gort.PooledSlab[MinimizedCondition_MeasuredNode], this_ *MinimizedCondition, p *MinimizedCondition_MeasuredNode, n *MinimizedCondition_MeasuredNode) {
	_ = p
	_ = n
	this_.initFields_MinimizedCondition()
	this_.positive = p
	this_.negative = n.changeWithSlab_(slab_)
}

// Java static method: MinimizedCondition.isMinimizable(com.google.javascript.rhino.Node n)
func MinimizedCondition_IsMinimizable(n pkg_rhino.Node) bool {
	clinit_MinimizedCondition()
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(19) /* Token.NOT */, int32(69) /* Token.AND */, int32(68) /* Token.OR */, int32(67) /* Token.HOOK */, int32(50) /* Token.COMMA */ :
		return true
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java static method: MinimizedCondition.fromConditionNode(com.google.javascript.rhino.Node n)
func MinimizedCondition_FromConditionNode(n pkg_rhino.Node) *MinimizedCondition {
	clinit_MinimizedCondition()
	_ = n
	if !(n.HasParent()) {
		j2gort.CheckState(false)
	}
	switch n.TokenSwitchOrd_() {
	case int32(19) /* Token.NOT */, int32(69) /* Token.AND */, int32(68) /* Token.OR */, int32(67) /* Token.HOOK */, int32(50) /* Token.COMMA */ :
		return MinimizedCondition_ComputeMinimizedCondition(n)
	default:
		return MinimizedCondition_Unoptimized(n)
	}
	panic("unreachable switch expression")
}
func minimizedConditionFromConditionNodeWithSlab_(slabMeasured_ *j2gort.PooledSlab[MinimizedCondition_MeasuredNode], slabCond_ *j2gort.PooledSlab[MinimizedCondition], n pkg_rhino.Node) *MinimizedCondition {
	slabMeasured_.Reset()
	slabCond_.Reset()
	clinit_MinimizedCondition()
	_ = n
	if !(n.HasParent()) {
		j2gort.CheckState(false)
	}
	switch n.TokenSwitchOrd_() {
	case int32(19), int32(69), int32(68), int32(67), int32(50):
		return minimizedConditionComputeWithSlab_(slabMeasured_, slabCond_, n)
	default:
		return minimizedConditionUnoptimizedWithSlab_(slabMeasured_, slabCond_, n)
	}
	panic("unreachable switch expression")
}

// Java method: MinimizedCondition.getMinimized(com.google.javascript.jscomp.MinimizedCondition.MinimizationStyle style)
func (this_ *MinimizedCondition) GetMinimized(style MinimizedCondition_MinimizationStyle) *MinimizedCondition_MeasuredNode {
	_ = style
	if ((style == StaticGet_MinimizedCondition_MinimizationStyle_PREFER_UNNEGATED()) || this_.positive.node.IsNot()) || (this_.positive.length <= this_.negative.length) {
		return this_.positive
	} else {
		return this_.negative.priv_MeasuredNode_addNot_3876()
	}
}
func (this_ *MinimizedCondition) getMinimizedWithSlab_(slab_ *j2gort.PooledSlab[MinimizedCondition_MeasuredNode], style MinimizedCondition_MinimizationStyle) *MinimizedCondition_MeasuredNode {
	_ = style
	if ((style == StaticGet_MinimizedCondition_MinimizationStyle_PREFER_UNNEGATED()) || this_.positive.node.IsNot()) || (this_.positive.length <= this_.negative.length) {
		return this_.positive
	} else {
		return this_.negative.addNotWithSlab_(slab_)
	}
}

// Java static method: MinimizedCondition.unoptimized(com.google.javascript.rhino.Node n)
func MinimizedCondition_Unoptimized(n pkg_rhino.Node) *MinimizedCondition {
	clinit_MinimizedCondition()
	_ = n
	j2gort.CheckNotNullNode(n.GetParent())
	clinit_MinimizedCondition_MeasuredNode()
	var val_1 MinimizedCondition_MeasuredNode
	InitMinimizedCondition_MeasuredNode(&val_1, n, int32(-1), nil, nil, nil, int32(0), false)
	var pos *MinimizedCondition_MeasuredNode = &val_1
	_ = pos
	return NewMinimizedCondition(pos, StaticGet_MinimizedCondition_UNOPTIMIZED_NEG())
}
func minimizedConditionUnoptimizedWithSlab_(slabMeasured_ *j2gort.PooledSlab[MinimizedCondition_MeasuredNode], slabCond_ *j2gort.PooledSlab[MinimizedCondition], n pkg_rhino.Node) *MinimizedCondition {
	clinit_MinimizedCondition()
	_ = n
	j2gort.CheckNotNullNode(n.GetParent())
	clinit_MinimizedCondition_MeasuredNode()
	val_1_slabPtr_ := allocPooled_MinimizedCondition_MeasuredNode(slabMeasured_, n, int32(-1), nil, nil, nil, int32(0), false)
	var pos *MinimizedCondition_MeasuredNode = val_1_slabPtr_
	_ = pos
	return allocPooled_MinimizedCondition(slabMeasured_, slabCond_, pos, StaticGet_MinimizedCondition_UNOPTIMIZED_NEG())
}

// Java static method: MinimizedCondition.pickBest(com.google.javascript.jscomp.MinimizedCondition.MeasuredNode a, com.google.javascript.jscomp.MinimizedCondition.MeasuredNode b)
func MinimizedCondition_PickBest(a *MinimizedCondition_MeasuredNode, b *MinimizedCondition_MeasuredNode) *MinimizedCondition_MeasuredNode {
	clinit_MinimizedCondition()
	_ = a
	_ = b
	if a.length == b.length {
		if b.IsChanged() {
			return a
		} else {
			return b
		}
	}
	if a.length < b.length {
		return a
	} else {
		return b
	}
}

// Java static method: MinimizedCondition.computeMinimizedCondition(com.google.javascript.rhino.Node n)
func MinimizedCondition_ComputeMinimizedCondition(n pkg_rhino.Node) *MinimizedCondition {
	clinit_MinimizedCondition()
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(19) /* Token.NOT */ :
		var subtree *MinimizedCondition = MinimizedCondition_ComputeMinimizedCondition(n.GetFirstChild())
		_ = subtree
		var positive *MinimizedCondition_MeasuredNode = MinimizedCondition_PickBest(MinimizedCondition_MeasuredNode_AddNode_2(n, subtree.positive), subtree.negative)
		_ = positive
		var negative *MinimizedCondition_MeasuredNode = MinimizedCondition_PickBest(subtree.negative.priv_MeasuredNode_negate_3874(), subtree.positive)
		_ = negative
		return NewMinimizedCondition(positive, negative)
	case int32(69) /* Token.AND */, int32(68) /* Token.OR */ :
		var complementNode pkg_rhino.Node = pkg_rhino.NewNode_1(func() pkg_rhino.Token {
			if n.IsAnd() {
				return pkg_rhino.StaticGet_Token_OR()
			}
			return pkg_rhino.StaticGet_Token_AND()
		}()).Srcref(n)
		_ = complementNode
		var leftSubtree *MinimizedCondition = MinimizedCondition_ComputeMinimizedCondition(n.GetFirstChild())
		_ = leftSubtree
		var rightSubtree *MinimizedCondition = MinimizedCondition_ComputeMinimizedCondition(n.GetLastChild())
		_ = rightSubtree
		var posTree *MinimizedCondition_MeasuredNode = MinimizedCondition_MeasuredNode_AddNode_3(n, leftSubtree.positive, rightSubtree.positive)
		_ = posTree
		var compTree *MinimizedCondition_MeasuredNode = MinimizedCondition_MeasuredNode_AddNode_3(complementNode, leftSubtree.negative, rightSubtree.negative)
		_ = compTree
		var positive_2 *MinimizedCondition_MeasuredNode = MinimizedCondition_PickBest(posTree, compTree.priv_MeasuredNode_negate_3874())
		_ = positive_2
		var negative_2 *MinimizedCondition_MeasuredNode = MinimizedCondition_PickBest(posTree.priv_MeasuredNode_negate_3874(), compTree.priv_MeasuredNode_change_3877())
		_ = negative_2
		return NewMinimizedCondition(positive_2, negative_2)
	case int32(67) /* Token.HOOK */ :
		var cond pkg_rhino.Node = n.GetFirstChild()
		_ = cond
		var thenNode pkg_rhino.Node = cond.GetNext()
		_ = thenNode
		var elseNode pkg_rhino.Node = thenNode.GetNext()
		_ = elseNode
		var thenSubtree *MinimizedCondition = MinimizedCondition_ComputeMinimizedCondition(thenNode)
		_ = thenSubtree
		var elseSubtree *MinimizedCondition = MinimizedCondition_ComputeMinimizedCondition(elseNode)
		_ = elseSubtree
		var condNode *MinimizedCondition_MeasuredNode = MinimizedCondition_MeasuredNode_ForNode(cond)
		_ = condNode
		var positive_3 *MinimizedCondition_MeasuredNode = MinimizedCondition_MeasuredNode_AddNode_4(n, condNode, thenSubtree.positive, elseSubtree.positive)
		_ = positive_3
		var negative_3 *MinimizedCondition_MeasuredNode = MinimizedCondition_MeasuredNode_AddNode_4(n, condNode, thenSubtree.negative, elseSubtree.negative)
		_ = negative_3
		return NewMinimizedCondition(positive_3, negative_3)
	case int32(50) /* Token.COMMA */ :
		var lhs pkg_rhino.Node = n.GetFirstChild()
		_ = lhs
		var rhsSubtree *MinimizedCondition = MinimizedCondition_ComputeMinimizedCondition(lhs.GetNext())
		_ = rhsSubtree
		var lhsNode *MinimizedCondition_MeasuredNode = MinimizedCondition_MeasuredNode_ForNode(lhs)
		_ = lhsNode
		var positive_4 *MinimizedCondition_MeasuredNode = MinimizedCondition_MeasuredNode_AddNode_3(n, lhsNode, rhsSubtree.positive)
		_ = positive_4
		var negative_4 *MinimizedCondition_MeasuredNode = MinimizedCondition_MeasuredNode_AddNode_3(n, lhsNode, rhsSubtree.negative)
		_ = negative_4
		return NewMinimizedCondition(positive_4, negative_4)
	default:
		var pos *MinimizedCondition_MeasuredNode = MinimizedCondition_MeasuredNode_ForNode(n)
		_ = pos
		var neg *MinimizedCondition_MeasuredNode = pos.priv_MeasuredNode_negate_3874()
		_ = neg
		return NewMinimizedCondition(pos, neg)
	}
	panic("unreachable switch expression")
}
func minimizedConditionComputeWithSlab_(slabMeasured_ *j2gort.PooledSlab[MinimizedCondition_MeasuredNode], slabCond_ *j2gort.PooledSlab[MinimizedCondition], n pkg_rhino.Node) *MinimizedCondition {
	clinit_MinimizedCondition()
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(19):
		var subtree *MinimizedCondition = minimizedConditionComputeWithSlab_(slabMeasured_, slabCond_, n.GetFirstChild())
		_ = subtree
		var positive *MinimizedCondition_MeasuredNode = MinimizedCondition_PickBest(minimizedConditionMeasuredNodeAddNode2WithSlab_(slabMeasured_, n, subtree.positive), subtree.negative)
		_ = positive
		var negative *MinimizedCondition_MeasuredNode = MinimizedCondition_PickBest(subtree.negative.negateWithSlab_(slabMeasured_), subtree.positive)
		_ = negative
		return allocPooled_MinimizedCondition(slabMeasured_, slabCond_, positive, negative)
	case int32(69), int32(68):
		var complementNode pkg_rhino.Node = pkg_rhino.NewNode_1(func() pkg_rhino.Token {
			if n.IsAnd() {
				return pkg_rhino.StaticGet_Token_OR()
			}
			return pkg_rhino.StaticGet_Token_AND()
		}()).Srcref(n)
		_ = complementNode
		var leftSubtree *MinimizedCondition = minimizedConditionComputeWithSlab_(slabMeasured_, slabCond_, n.GetFirstChild())
		_ = leftSubtree
		var rightSubtree *MinimizedCondition = minimizedConditionComputeWithSlab_(slabMeasured_, slabCond_, n.GetLastChild())
		_ = rightSubtree
		var posTree *MinimizedCondition_MeasuredNode = minimizedConditionMeasuredNodeAddNode3WithSlab_(slabMeasured_, n, leftSubtree.positive, rightSubtree.positive)
		_ = posTree
		var compTree *MinimizedCondition_MeasuredNode = minimizedConditionMeasuredNodeAddNode3WithSlab_(slabMeasured_, complementNode, leftSubtree.negative, rightSubtree.negative)
		_ = compTree
		var positive_2 *MinimizedCondition_MeasuredNode = MinimizedCondition_PickBest(posTree, compTree.negateWithSlab_(slabMeasured_))
		_ = positive_2
		var negative_2 *MinimizedCondition_MeasuredNode = MinimizedCondition_PickBest(posTree.negateWithSlab_(slabMeasured_), compTree.changeWithSlab_(slabMeasured_))
		_ = negative_2
		return allocPooled_MinimizedCondition(slabMeasured_, slabCond_, positive_2, negative_2)
	case int32(67):
		var cond pkg_rhino.Node = n.GetFirstChild()
		_ = cond
		var thenNode pkg_rhino.Node = cond.GetNext()
		_ = thenNode
		var elseNode pkg_rhino.Node = thenNode.GetNext()
		_ = elseNode
		var thenSubtree *MinimizedCondition = minimizedConditionComputeWithSlab_(slabMeasured_, slabCond_, thenNode)
		_ = thenSubtree
		var elseSubtree *MinimizedCondition = minimizedConditionComputeWithSlab_(slabMeasured_, slabCond_, elseNode)
		_ = elseSubtree
		var condNode *MinimizedCondition_MeasuredNode = minimizedConditionMeasuredNodeForNodeWithSlab_(slabMeasured_, cond)
		_ = condNode
		var positive_3 *MinimizedCondition_MeasuredNode = minimizedConditionMeasuredNodeAddNode4WithSlab_(slabMeasured_, n, condNode, thenSubtree.positive, elseSubtree.positive)
		_ = positive_3
		var negative_3 *MinimizedCondition_MeasuredNode = minimizedConditionMeasuredNodeAddNode4WithSlab_(slabMeasured_, n, condNode, thenSubtree.negative, elseSubtree.negative)
		_ = negative_3
		return allocPooled_MinimizedCondition(slabMeasured_, slabCond_, positive_3, negative_3)
	case int32(50):
		var lhs pkg_rhino.Node = n.GetFirstChild()
		_ = lhs
		var rhsSubtree *MinimizedCondition = minimizedConditionComputeWithSlab_(slabMeasured_, slabCond_, lhs.GetNext())
		_ = rhsSubtree
		var lhsNode *MinimizedCondition_MeasuredNode = minimizedConditionMeasuredNodeForNodeWithSlab_(slabMeasured_, lhs)
		_ = lhsNode
		var positive_4 *MinimizedCondition_MeasuredNode = minimizedConditionMeasuredNodeAddNode3WithSlab_(slabMeasured_, n, lhsNode, rhsSubtree.positive)
		_ = positive_4
		var negative_4 *MinimizedCondition_MeasuredNode = minimizedConditionMeasuredNodeAddNode3WithSlab_(slabMeasured_, n, lhsNode, rhsSubtree.negative)
		_ = negative_4
		return allocPooled_MinimizedCondition(slabMeasured_, slabCond_, positive_4, negative_4)
	default:
		var pos *MinimizedCondition_MeasuredNode = minimizedConditionMeasuredNodeForNodeWithSlab_(slabMeasured_, n)
		_ = pos
		var neg *MinimizedCondition_MeasuredNode = pos.negateWithSlab_(slabMeasured_)
		_ = neg
		return allocPooled_MinimizedCondition(slabMeasured_, slabCond_, pos, neg)
	}
	panic("unreachable switch expression")
}

// Java enum: com.google.javascript.jscomp.MinimizedCondition.MinimizationStyle
type MinimizedCondition_MinimizationStyle uint8

const (
	enum_MinimizedCondition_MinimizationStyle_PREFER_UNNEGATED  MinimizedCondition_MinimizationStyle = 1
	enum_MinimizedCondition_MinimizationStyle_ALLOW_LEADING_NOT MinimizedCondition_MinimizationStyle = 2
)

var enumBaseTable_MinimizedCondition_MinimizationStyle = [...]j2gort.EnumBase{
	{Name_: "PREFER_UNNEGATED", Ordinal_: 0},
	{Name_: "ALLOW_LEADING_NOT", Ordinal_: 1},
}
var enumValues_MinimizedCondition_MinimizationStyle = []MinimizedCondition_MinimizationStyle{1, 2}

func StaticGet_MinimizedCondition_MinimizationStyle_PREFER_UNNEGATED() MinimizedCondition_MinimizationStyle {
	return enum_MinimizedCondition_MinimizationStyle_PREFER_UNNEGATED
}
func StaticGet_MinimizedCondition_MinimizationStyle_ALLOW_LEADING_NOT() MinimizedCondition_MinimizationStyle {
	return enum_MinimizedCondition_MinimizationStyle_ALLOW_LEADING_NOT
}

func MinimizedCondition_MinimizationStyle_Values() []MinimizedCondition_MinimizationStyle {
	return j2gort.CloneSlice(enumValues_MinimizedCondition_MinimizationStyle)
}

func MinimizedCondition_MinimizationStyle_ValueOf(name string) MinimizedCondition_MinimizationStyle {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_MinimizedCondition_MinimizationStyle {
		if enumBaseTable_MinimizedCondition_MinimizationStyle[i_].Name_ == name {
			return MinimizedCondition_MinimizationStyle(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.MinimizedCondition.MinimizationStyle." + name))
}

func (this_ MinimizedCondition_MinimizationStyle) IsJavaNil_() bool                              { return this_ == 0 }
func (MinimizedCondition_MinimizationStyle) HasIdentityEquals_()                                 {}
func (MinimizedCondition_MinimizationStyle) Implements_j2gort_Ext_java_lang_Enum()               {}
func (MinimizedCondition_MinimizationStyle) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ MinimizedCondition_MinimizationStyle) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (MinimizedCondition_MinimizationStyle) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.MinimizedCondition$MinimizationStyle", "MinimizationStyle"
}
func (this_ MinimizedCondition_MinimizationStyle) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_MinimizedCondition_MinimizationStyle[this_-1]
}
func (this_ MinimizedCondition_MinimizationStyle) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ MinimizedCondition_MinimizationStyle) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ MinimizedCondition_MinimizationStyle) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_MinimizedCondition_MinimizationStyle[this_-1].Name_
}
func (this_ MinimizedCondition_MinimizationStyle) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_MinimizedCondition_MinimizationStyle[this_-1].Name_
}
func (this_ MinimizedCondition_MinimizationStyle) Equals(other_ any) bool {
	o_, ok_ := other_.(MinimizedCondition_MinimizationStyle)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ MinimizedCondition_MinimizationStyle) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_MinimizedCondition_MinimizationStyle[this_-1].HashCode()
}
func (this_ MinimizedCondition_MinimizationStyle) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_MinimizedCondition_MinimizationStyle[this_-1].CompareTo(other_)
}
func (this_ MinimizedCondition_MinimizationStyle) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_MinimizedCondition_MinimizationStyle[this_-1].Name_)
}
func init() {
	vals_ := []any{MinimizedCondition_MinimizationStyle(1), MinimizedCondition_MinimizationStyle(2)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.MinimizedCondition.MinimizationStyle", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.MinimizedCondition$MinimizationStyle", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.MinimizedCondition.MinimizationStyle", "com.google.javascript.jscomp.MinimizedCondition$MinimizationStyle", func(v_ any) bool {
		return j2gort.InstanceOf[MinimizedCondition_MinimizationStyle](v_)
	})
}

// Java class: com.google.javascript.jscomp.MinimizedCondition.MeasuredNode
type MinimizedCondition_MeasuredNode struct {
	c0         *MinimizedCondition_MeasuredNode `gson:"c0"`
	c1         *MinimizedCondition_MeasuredNode `gson:"c1"`
	c2         *MinimizedCondition_MeasuredNode `gson:"c2"`
	node       pkg_rhino.Node                   `gson:"node"`
	length     int32                            `gson:"length"`
	childCount int32                            `gson:"childCount"`
	idHash_    atomic.Int32
	changed    bool `gson:"changed"`
}

func (this_ *MinimizedCondition_MeasuredNode) BindSelf_(s any) {
	_ = s
}
func (this_ *MinimizedCondition_MeasuredNode) IsJavaNil_() bool { return this_ == nil }
func (*MinimizedCondition_MeasuredNode) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.MinimizedCondition$MeasuredNode", "MeasuredNode"
}
func (this_ *MinimizedCondition_MeasuredNode) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*MinimizedCondition_MeasuredNode) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.MinimizedCondition.MeasuredNode", "com.google.javascript.jscomp.MinimizedCondition$MeasuredNode", func(v_ any) bool {
		return j2gort.InstanceOf[*MinimizedCondition_MeasuredNode](v_)
	})
}

var static_MinimizedCondition_MeasuredNode_NOT_PRECEDENCE int32

func StaticGet_MinimizedCondition_MeasuredNode_NOT_PRECEDENCE() int32 {
	clinit_MinimizedCondition_MeasuredNode()
	return static_MinimizedCondition_MeasuredNode_NOT_PRECEDENCE
}
func StaticSet_MinimizedCondition_MeasuredNode_NOT_PRECEDENCE(v int32) {
	clinit_MinimizedCondition_MeasuredNode()
	static_MinimizedCondition_MeasuredNode_NOT_PRECEDENCE = v
}

var clinitDone_MinimizedCondition_MeasuredNode uint32
var clinitState_MinimizedCondition_MeasuredNode j2gort.ClinitState

func clinit_MinimizedCondition_MeasuredNode() {
	if atomic.LoadUint32(&clinitDone_MinimizedCondition_MeasuredNode) != 0 {
		return
	}
	clinitSlow_MinimizedCondition_MeasuredNode()
}

func clinitSlow_MinimizedCondition_MeasuredNode() {
	if !j2gort.EnterClinit(&clinitState_MinimizedCondition_MeasuredNode, &clinitDone_MinimizedCondition_MeasuredNode) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_MinimizedCondition_MeasuredNode, &clinitDone_MinimizedCondition_MeasuredNode)
	static_MinimizedCondition_MeasuredNode_NOT_PRECEDENCE = NodeUtil_Precedence(pkg_rhino.StaticGet_Token_NOT())
}

func (this_ *MinimizedCondition_MeasuredNode) initFields_MinimizedCondition_MeasuredNode() {
	_ = this_
}

// Java constructor: MeasuredNode(com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node n, int childCount, com.google.javascript.jscomp.MinimizedCondition.@org.jspecify.annotations.Nullable MeasuredNode c0, com.google.javascript.jscomp.MinimizedCondition.@org.jspecify.annotations.Nullable MeasuredNode c1, com.google.javascript.jscomp.MinimizedCondition.@org.jspecify.annotations.Nullable MeasuredNode c2, int len, boolean ch)
func NewMinimizedCondition_MeasuredNode(n pkg_rhino.Node, childCount int32, c0 *MinimizedCondition_MeasuredNode, c1 *MinimizedCondition_MeasuredNode, c2 *MinimizedCondition_MeasuredNode, len_ int32, ch bool) *MinimizedCondition_MeasuredNode {
	clinit_MinimizedCondition_MeasuredNode()
	this_ := &MinimizedCondition_MeasuredNode{}
	InitMinimizedCondition_MeasuredNode(this_, n, childCount, c0, c1, c2, len_, ch)
	return this_
}
func allocPooled_MinimizedCondition_MeasuredNode(slab_ *j2gort.PooledSlab[MinimizedCondition_MeasuredNode], n pkg_rhino.Node, childCount int32, c0 *MinimizedCondition_MeasuredNode, c1 *MinimizedCondition_MeasuredNode, c2 *MinimizedCondition_MeasuredNode, len_ int32, ch bool) *MinimizedCondition_MeasuredNode {
	clinit_MinimizedCondition_MeasuredNode()
	this_ := slab_.Alloc()
	InitMinimizedCondition_MeasuredNode(this_, n, childCount, c0, c1, c2, len_, ch)
	return this_
}

func InitMinimizedCondition_MeasuredNode(this_ *MinimizedCondition_MeasuredNode, n pkg_rhino.Node, childCount int32, c0 *MinimizedCondition_MeasuredNode, c1 *MinimizedCondition_MeasuredNode, c2 *MinimizedCondition_MeasuredNode, len_ int32, ch bool) {
	_ = n
	_ = childCount
	_ = c0
	_ = c1
	_ = c2
	_ = len_
	_ = ch
	this_.initFields_MinimizedCondition_MeasuredNode()
	this_.node = n
	this_.childCount = childCount
	this_.c0 = c0
	this_.c1 = c1
	this_.c2 = c2
	this_.length = len_
	this_.changed = ch
}

// Java method: MeasuredNode.isChanged()
func (this_ *MinimizedCondition_MeasuredNode) IsChanged() bool {
	return this_.changed
}

// Java method: MeasuredNode.isNot()
func (this_ *MinimizedCondition_MeasuredNode) IsNot() bool {
	return this_.node.IsNot()
}

// Java method: MeasuredNode.withoutNot()
func (this_ *MinimizedCondition_MeasuredNode) WithoutNot() *MinimizedCondition_MeasuredNode {
	if !(this_.IsNot()) {
		j2gort.CheckState(false)
	}
	var first *MinimizedCondition_MeasuredNode
	_ = first
	if this_.childCount >= int32(0) {
		first = this_.c0
	} else {
		first = MinimizedCondition_MeasuredNode_ForNode(this_.node.GetFirstChild())
	}
	return first.priv_MeasuredNode_change_3877()
}
func (this_ *MinimizedCondition_MeasuredNode) withoutNotWithSlab_(slab_ *j2gort.PooledSlab[MinimizedCondition_MeasuredNode]) *MinimizedCondition_MeasuredNode {
	if !(this_.IsNot()) {
		j2gort.CheckState(false)
	}
	var first *MinimizedCondition_MeasuredNode
	_ = first
	if this_.childCount >= int32(0) {
		first = this_.c0
	} else {
		first = minimizedConditionMeasuredNodeForNodeWithSlab_(slab_, this_.node.GetFirstChild())
	}
	return first.changeWithSlab_(slab_)
}

// Java method: MeasuredNode.negate()
func (this_ *MinimizedCondition_MeasuredNode) priv_MeasuredNode_negate_3874() *MinimizedCondition_MeasuredNode {
	switch this_.node.TokenSwitchOrd_() {
	case int32(4) /* Token.EQ */ :
		return this_.priv_MeasuredNode_updateToken_3875(pkg_rhino.StaticGet_Token_NE())
	case int32(5) /* Token.NE */ :
		return this_.priv_MeasuredNode_updateToken_3875(pkg_rhino.StaticGet_Token_EQ())
	case int32(40) /* Token.SHEQ */ :
		return this_.priv_MeasuredNode_updateToken_3875(pkg_rhino.StaticGet_Token_SHNE())
	case int32(41) /* Token.SHNE */ :
		return this_.priv_MeasuredNode_updateToken_3875(pkg_rhino.StaticGet_Token_SHEQ())
	case int32(19) /* Token.NOT */ :
		return this_.WithoutNot()
	default:
		return this_.priv_MeasuredNode_addNot_3876()
	}
	panic("unreachable switch expression")
}
func (this_ *MinimizedCondition_MeasuredNode) negateWithSlab_(slab_ *j2gort.PooledSlab[MinimizedCondition_MeasuredNode]) *MinimizedCondition_MeasuredNode {
	switch this_.node.TokenSwitchOrd_() {
	case int32(4):
		return this_.updateTokenWithSlab_(slab_, pkg_rhino.StaticGet_Token_NE())
	case int32(5):
		return this_.updateTokenWithSlab_(slab_, pkg_rhino.StaticGet_Token_EQ())
	case int32(40):
		return this_.updateTokenWithSlab_(slab_, pkg_rhino.StaticGet_Token_SHNE())
	case int32(41):
		return this_.updateTokenWithSlab_(slab_, pkg_rhino.StaticGet_Token_SHEQ())
	case int32(19):
		return this_.withoutNotWithSlab_(slab_)
	default:
		return this_.addNotWithSlab_(slab_)
	}
	panic("unreachable switch expression")
}

// Java method: MeasuredNode.updateToken(com.google.javascript.rhino.Token token)
func (this_ *MinimizedCondition_MeasuredNode) priv_MeasuredNode_updateToken_3875(token pkg_rhino.Token) *MinimizedCondition_MeasuredNode {
	_ = token
	var newNode pkg_rhino.Node = pkg_rhino.NewNode_1(token).Srcref(this_.node)
	_ = newNode
	if this_.childCount >= int32(0) {
		return NewMinimizedCondition_MeasuredNode(newNode, this_.childCount, this_.c0, this_.c1, this_.c2, this_.length, true)
	}
	var first pkg_rhino.Node = this_.node.GetFirstChild()
	_ = first
	if first == 0 {
		return NewMinimizedCondition_MeasuredNode(newNode, int32(-1), nil, nil, nil, this_.length, true)
	}
	var second pkg_rhino.Node = first.GetNext()
	_ = second
	if !((second != 0) && (second.GetNext() == 0)) {
		j2gort.CheckState(false)
	}
	return NewMinimizedCondition_MeasuredNode(newNode, int32(2), MinimizedCondition_MeasuredNode_ForNode(first), MinimizedCondition_MeasuredNode_ForNode(second), nil, this_.length, true)
}
func (this_ *MinimizedCondition_MeasuredNode) updateTokenWithSlab_(slab_ *j2gort.PooledSlab[MinimizedCondition_MeasuredNode], token pkg_rhino.Token) *MinimizedCondition_MeasuredNode {
	_ = token
	var newNode pkg_rhino.Node = pkg_rhino.NewNode_1(token).Srcref(this_.node)
	_ = newNode
	if this_.childCount >= int32(0) {
		return allocPooled_MinimizedCondition_MeasuredNode(slab_, newNode, this_.childCount, this_.c0, this_.c1, this_.c2, this_.length, true)
	}
	var first pkg_rhino.Node = this_.node.GetFirstChild()
	_ = first
	if first == 0 {
		return allocPooled_MinimizedCondition_MeasuredNode(slab_, newNode, int32(-1), nil, nil, nil, this_.length, true)
	}
	var second pkg_rhino.Node = first.GetNext()
	_ = second
	if !((second != 0) && (second.GetNext() == 0)) {
		j2gort.CheckState(false)
	}
	return allocPooled_MinimizedCondition_MeasuredNode(slab_, newNode, int32(2), minimizedConditionMeasuredNodeForNodeWithSlab_(slab_, first), minimizedConditionMeasuredNodeForNodeWithSlab_(slab_, second), nil, this_.length, true)
}

// Java method: MeasuredNode.addNot()
func (this_ *MinimizedCondition_MeasuredNode) priv_MeasuredNode_addNot_3876() *MinimizedCondition_MeasuredNode {
	var cost int32 = (this_.length + int32(1))
	_ = cost
	if this_.IsLowerPrecedenceThan(StaticGet_MinimizedCondition_MeasuredNode_NOT_PRECEDENCE()) {
		cost = (cost + int32(2))
	}
	return NewMinimizedCondition_MeasuredNode(pkg_rhino.NewNode_1(pkg_rhino.StaticGet_Token_NOT()).Srcref(this_.node), int32(1), this_, nil, nil, cost, true)
}
func (this_ *MinimizedCondition_MeasuredNode) addNotWithSlab_(slab_ *j2gort.PooledSlab[MinimizedCondition_MeasuredNode]) *MinimizedCondition_MeasuredNode {
	var cost int32 = (this_.length + int32(1))
	_ = cost
	if this_.IsLowerPrecedenceThan(StaticGet_MinimizedCondition_MeasuredNode_NOT_PRECEDENCE()) {
		cost = (cost + int32(2))
	}
	return allocPooled_MinimizedCondition_MeasuredNode(slab_, pkg_rhino.NewNode_1(pkg_rhino.StaticGet_Token_NOT()).Srcref(this_.node), int32(1), this_, nil, nil, cost, true)
}

// Java method: MeasuredNode.change()
func (this_ *MinimizedCondition_MeasuredNode) priv_MeasuredNode_change_3877() *MinimizedCondition_MeasuredNode {
	if this_.IsChanged() {
		return this_
	} else {
		return NewMinimizedCondition_MeasuredNode(this_.node, this_.childCount, this_.c0, this_.c1, this_.c2, this_.length, true)
	}
}
func (this_ *MinimizedCondition_MeasuredNode) changeWithSlab_(slab_ *j2gort.PooledSlab[MinimizedCondition_MeasuredNode]) *MinimizedCondition_MeasuredNode {
	if this_.IsChanged() {
		return this_
	} else {
		return allocPooled_MinimizedCondition_MeasuredNode(slab_, this_.node, this_.childCount, this_.c0, this_.c1, this_.c2, this_.length, true)
	}
}

// Java method: MeasuredNode.isLowerPrecedenceThan(int precedence)
func (this_ *MinimizedCondition_MeasuredNode) IsLowerPrecedenceThan(precedence int32) bool {
	_ = precedence
	return (NodeUtil_Precedence(this_.node.GetToken()) < precedence)
}

// Java static method: MeasuredNode.addNode(com.google.javascript.rhino.Node parent, com.google.javascript.jscomp.MinimizedCondition.MeasuredNode child0)
func MinimizedCondition_MeasuredNode_AddNode_2(parent pkg_rhino.Node, child0 *MinimizedCondition_MeasuredNode) *MinimizedCondition_MeasuredNode {
	clinit_MinimizedCondition_MeasuredNode()
	_ = parent
	_ = child0
	var cost int32 = child0.length
	_ = cost
	if parent.IsNot() {
		cost = cost + 1
	}
	var parentPrecedence int32 = NodeUtil_Precedence(parent.GetToken())
	_ = parentPrecedence
	if child0.IsLowerPrecedenceThan(parentPrecedence) {
		cost = (cost + int32(2))
	}
	return NewMinimizedCondition_MeasuredNode(parent, int32(1), child0, nil, nil, cost, child0.changed)
}
func minimizedConditionMeasuredNodeAddNode2WithSlab_(slab_ *j2gort.PooledSlab[MinimizedCondition_MeasuredNode], parent pkg_rhino.Node, child0 *MinimizedCondition_MeasuredNode) *MinimizedCondition_MeasuredNode {
	clinit_MinimizedCondition_MeasuredNode()
	_ = parent
	_ = child0
	var cost int32 = child0.length
	_ = cost
	if parent.IsNot() {
		cost = cost + 1
	}
	var parentPrecedence int32 = NodeUtil_Precedence(parent.GetToken())
	_ = parentPrecedence
	if child0.IsLowerPrecedenceThan(parentPrecedence) {
		cost = (cost + int32(2))
	}
	return allocPooled_MinimizedCondition_MeasuredNode(slab_, parent, int32(1), child0, nil, nil, cost, child0.changed)
}

// Java static method: MeasuredNode.addNode(com.google.javascript.rhino.Node parent, com.google.javascript.jscomp.MinimizedCondition.MeasuredNode child0, com.google.javascript.jscomp.MinimizedCondition.MeasuredNode child1)
func MinimizedCondition_MeasuredNode_AddNode_3(parent pkg_rhino.Node, child0 *MinimizedCondition_MeasuredNode, child1 *MinimizedCondition_MeasuredNode) *MinimizedCondition_MeasuredNode {
	clinit_MinimizedCondition_MeasuredNode()
	_ = parent
	_ = child0
	_ = child1
	var cost int32 = (child0.length + child1.length)
	_ = cost
	if parent.IsNot() {
		cost = cost + 1
	}
	var parentPrecedence int32 = NodeUtil_Precedence(parent.GetToken())
	_ = parentPrecedence
	if child0.IsLowerPrecedenceThan(parentPrecedence) {
		cost = (cost + int32(2))
	}
	if child1.IsLowerPrecedenceThan(parentPrecedence) {
		cost = (cost + int32(2))
	}
	return NewMinimizedCondition_MeasuredNode(parent, int32(2), child0, child1, nil, cost, (child0.changed || child1.changed))
}
func minimizedConditionMeasuredNodeAddNode3WithSlab_(slab_ *j2gort.PooledSlab[MinimizedCondition_MeasuredNode], parent pkg_rhino.Node, child0 *MinimizedCondition_MeasuredNode, child1 *MinimizedCondition_MeasuredNode) *MinimizedCondition_MeasuredNode {
	clinit_MinimizedCondition_MeasuredNode()
	_ = parent
	_ = child0
	_ = child1
	var cost int32 = (child0.length + child1.length)
	_ = cost
	if parent.IsNot() {
		cost = cost + 1
	}
	var parentPrecedence int32 = NodeUtil_Precedence(parent.GetToken())
	_ = parentPrecedence
	if child0.IsLowerPrecedenceThan(parentPrecedence) {
		cost = (cost + int32(2))
	}
	if child1.IsLowerPrecedenceThan(parentPrecedence) {
		cost = (cost + int32(2))
	}
	return allocPooled_MinimizedCondition_MeasuredNode(slab_, parent, int32(2), child0, child1, nil, cost, (child0.changed || child1.changed))
}

// Java static method: MeasuredNode.addNode(com.google.javascript.rhino.Node parent, com.google.javascript.jscomp.MinimizedCondition.MeasuredNode child0, com.google.javascript.jscomp.MinimizedCondition.MeasuredNode child1, com.google.javascript.jscomp.MinimizedCondition.MeasuredNode child2)
func MinimizedCondition_MeasuredNode_AddNode_4(parent pkg_rhino.Node, child0 *MinimizedCondition_MeasuredNode, child1 *MinimizedCondition_MeasuredNode, child2 *MinimizedCondition_MeasuredNode) *MinimizedCondition_MeasuredNode {
	clinit_MinimizedCondition_MeasuredNode()
	_ = parent
	_ = child0
	_ = child1
	_ = child2
	var cost int32 = ((child0.length + child1.length) + child2.length)
	_ = cost
	if parent.IsNot() {
		cost = cost + 1
	}
	var parentPrecedence int32 = NodeUtil_Precedence(parent.GetToken())
	_ = parentPrecedence
	if child0.IsLowerPrecedenceThan(parentPrecedence) {
		cost = (cost + int32(2))
	}
	if child1.IsLowerPrecedenceThan(parentPrecedence) {
		cost = (cost + int32(2))
	}
	if child2.IsLowerPrecedenceThan(parentPrecedence) {
		cost = (cost + int32(2))
	}
	return NewMinimizedCondition_MeasuredNode(parent, int32(3), child0, child1, child2, cost, ((child0.changed || child1.changed) || child2.changed))
}
func minimizedConditionMeasuredNodeAddNode4WithSlab_(slab_ *j2gort.PooledSlab[MinimizedCondition_MeasuredNode], parent pkg_rhino.Node, child0 *MinimizedCondition_MeasuredNode, child1 *MinimizedCondition_MeasuredNode, child2 *MinimizedCondition_MeasuredNode) *MinimizedCondition_MeasuredNode {
	clinit_MinimizedCondition_MeasuredNode()
	_ = parent
	_ = child0
	_ = child1
	_ = child2
	var cost int32 = ((child0.length + child1.length) + child2.length)
	_ = cost
	if parent.IsNot() {
		cost = cost + 1
	}
	var parentPrecedence int32 = NodeUtil_Precedence(parent.GetToken())
	_ = parentPrecedence
	if child0.IsLowerPrecedenceThan(parentPrecedence) {
		cost = (cost + int32(2))
	}
	if child1.IsLowerPrecedenceThan(parentPrecedence) {
		cost = (cost + int32(2))
	}
	if child2.IsLowerPrecedenceThan(parentPrecedence) {
		cost = (cost + int32(2))
	}
	return allocPooled_MinimizedCondition_MeasuredNode(slab_, parent, int32(3), child0, child1, child2, cost, ((child0.changed || child1.changed) || child2.changed))
}

// Java static method: MeasuredNode.forNode(com.google.javascript.rhino.Node n)
func MinimizedCondition_MeasuredNode_ForNode(n pkg_rhino.Node) *MinimizedCondition_MeasuredNode {
	clinit_MinimizedCondition_MeasuredNode()
	_ = n
	return NewMinimizedCondition_MeasuredNode(n, int32(-1), nil, nil, nil, int32(0), false)
}
func minimizedConditionMeasuredNodeForNodeWithSlab_(slab_ *j2gort.PooledSlab[MinimizedCondition_MeasuredNode], n pkg_rhino.Node) *MinimizedCondition_MeasuredNode {
	clinit_MinimizedCondition_MeasuredNode()
	_ = n
	return allocPooled_MinimizedCondition_MeasuredNode(slab_, n, int32(-1), nil, nil, nil, int32(0), false)
}

// Java method: MeasuredNode.willChange(com.google.javascript.rhino.Node original)
func (this_ *MinimizedCondition_MeasuredNode) WillChange(original pkg_rhino.Node) bool {
	_ = original
	j2gort.CheckNotNullNode(original)
	return ((original != this_.node) || this_.IsChanged())
}

// Java method: MeasuredNode.applyTo(com.google.javascript.rhino.Node original)
func (this_ *MinimizedCondition_MeasuredNode) ApplyTo(original pkg_rhino.Node) pkg_rhino.Node {
	_ = original
	j2gort.CheckNotNullNode(original)
	if !(this_.WillChange(original)) {
		j2gort.CheckState(false)
	}
	var replacement pkg_rhino.Node = this_.BuildReplacement()
	_ = replacement
	if original != replacement {
		this_.priv_MeasuredNode_safeDetach_3881(replacement)
		original.ReplaceWith(replacement)
	}
	return replacement
}

// Java method: MeasuredNode.safeDetach(com.google.javascript.rhino.Node n)
func (this_ *MinimizedCondition_MeasuredNode) priv_MeasuredNode_safeDetach_3881(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	if n.HasParent() {
		return n.Detach()
	} else {
		return n
	}
}

// Java method: MeasuredNode.buildReplacement()
func (this_ *MinimizedCondition_MeasuredNode) BuildReplacement() pkg_rhino.Node {
	if this_.childCount >= int32(0) {
		this_.node.DetachChildren()
		if this_.childCount >= int32(1) {
			this_.node.AddChildToBack(this_.priv_MeasuredNode_safeDetach_3881(this_.c0.BuildReplacement()))
		}
		if this_.childCount >= int32(2) {
			this_.node.AddChildToBack(this_.priv_MeasuredNode_safeDetach_3881(this_.c1.BuildReplacement()))
		}
		if this_.childCount >= int32(3) {
			this_.node.AddChildToBack(this_.priv_MeasuredNode_safeDetach_3881(this_.c2.BuildReplacement()))
		}
	}
	return this_.node
}
