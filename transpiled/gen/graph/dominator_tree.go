/*
 * Copyright 2026 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/graph/DominatorTree.java (generated mechanically by j2go).

package graph

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.graph.DominatorTree
type DominatorTree[N comparable] struct {
	idoms        j2gort.Ext_com_google_common_collect_ImmutableMap `gson:"idoms"`
	subtreeSizes j2gort.Ext_com_google_common_collect_ImmutableMap `gson:"subtreeSizes"`
	idHash_      atomic.Int32
}

func (this_ *DominatorTree[N]) BindSelf_(s any) {
	_ = s
}
func (this_ *DominatorTree[N]) IsJavaNil_() bool { return this_ == nil }
func (*DominatorTree[N]) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.graph.DominatorTree", "DominatorTree"
}
func (this_ *DominatorTree[N]) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*DominatorTree[N]) HasIdentityEquals_()                 {}
func (*DominatorTree[N]) IsGenericInst_graph_DominatorTree_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.graph.DominatorTree", "com.google.javascript.jscomp.graph.DominatorTree", func(v_ any) bool {
		m_, ok_ := v_.(interface {
			IsGenericInst_graph_DominatorTree_()
			IsJavaNil_() bool
		})
		return ok_ && !m_.IsJavaNil_()
	})
}

func (this_ *DominatorTree[N]) initFields_DominatorTree() {
	_ = this_
}

// Java constructor: DominatorTree(java.util.Map<N,N> idoms, java.util.Map<N,java.lang.Integer> subtreeSizes)
func NewDominatorTree[N comparable](idoms j2gort.Ext_java_util_Map, subtreeSizes j2gort.Ext_java_util_Map) *DominatorTree[N] {
	this_ := &DominatorTree[N]{}
	InitDominatorTree[N](this_, idoms, subtreeSizes)
	return this_
}

func InitDominatorTree[N comparable](this_ *DominatorTree[N], idoms j2gort.Ext_java_util_Map, subtreeSizes j2gort.Ext_java_util_Map) {
	_ = idoms
	_ = subtreeSizes
	this_.initFields_DominatorTree()
	this_.idoms = j2gort.Ext_com_google_common_collect_ImmutableMap_CopyOf_1_Ext_java_util_Map(idoms)
	this_.subtreeSizes = j2gort.Ext_com_google_common_collect_ImmutableMap_CopyOf_1_Ext_java_util_Map(subtreeSizes)
}

// Java static method: DominatorTree.compute(com.google.javascript.jscomp.graph.DiGraph<N,E> graph, N entry)
func DominatorTree_Compute[N comparable, E any](graph DiGraph[N, E], entry N) *DominatorTree[N] {
	_ = graph
	_ = entry
	j2gort.CheckNotNullIface(j2gort.AsIface[DiGraph_DiGraphNode[N, E]](graph.GetNode(entry)))
	var postOrder j2gort.Ext_java_util_ArrayList = j2gort.NewExt_java_util_ArrayList()
	_ = postOrder
	var postOrderIndex j2gort.OrderedMap[N, int32] = j2gort.NewOrderedMap[N, int32]()
	_ = postOrderIndex
	DominatorTree_BuildPostOrder[N, E](graph, entry, postOrder, j2gort.NewExt_java_util_LinkedHashMap())
	var numNodes int32 = postOrder.Size()
	_ = numNodes
	{
		var i int32 = int32(0)
		_ = i
		for ; i < numNodes; i = i + 1 {
			postOrderIndex.Put(func() N {
				tv_1 := postOrder.Get(i)
				if tv_1 == nil {
					{

						if _, ok_ := any((*N)(nil)).(*string); ok_ {
							return any(j2gort.NullString).(N)
						}
						var z_ N
						return z_
					}

				}
				return j2gort.Cast[N](tv_1)
			}(), i)
		}
	}
	var idomIndexes []int32 = make([]int32, numNodes)
	_ = idomIndexes
	j2gort.ArraysFillRange(idomIndexes, int32(0), numNodes, int32(-1))
	var startNodeIndex int32
iifeSwitch_6:
	switch {
	default:
		if v_, ok_ := postOrderIndex.GetOk(entry); ok_ {
			startNodeIndex = v_
			break iifeSwitch_6
		}
		j2gort.ThrowNPE()
		startNodeIndex = 0
		break iifeSwitch_6
	}

	_ = startNodeIndex
	{
		a_2 := idomIndexes
		i_3 := startNodeIndex
		a_2[i_3] = startNodeIndex
	}
	var changed bool = true
	_ = changed
	for changed {
		changed = false
		{
			var r int32 = (numNodes - int32(1))
			_ = r
			for ; r >= int32(0); r = r - 1 {
				var node N
			iifeSwitch_1:
				switch {
				default:
					tv_4 := postOrder.Get(r)
					if tv_4 == nil {
						{

							if _, ok_ := any((*N)(nil)).(*string); ok_ {
								node = any(j2gort.NullString).(N)
								break iifeSwitch_1
							}
							var z_ N
							node = z_
							break iifeSwitch_1
						}

					}
					node = j2gort.Cast[N](tv_4)
					break iifeSwitch_1
				}

				_ = node
				if func() bool {
					rk_5 := node
					ak_6 := entry
					if _, ok_ := any((*N)(nil)).(*string); ok_ {
						_ = j2gort.CheckNotNullString(any(rk_5).(string))
						return rk_5 == ak_6
					}
					var z_ N
					if rk_5 == z_ {
						j2gort.ThrowNPE()
					}
					if rk_5 == ak_6 {
						return true
					}
					return javalang.ObjectsEquals(j2gort.ToAny(rk_5), j2gort.ToAny(ak_6))
				}() {
					continue
				}
				var nodeIdx int32
			iifeSwitch_5:
				switch {
				default:
					if v_, ok_ := postOrderIndex.GetOk(node); ok_ {
						nodeIdx = v_
						break iifeSwitch_5
					}
					j2gort.ThrowNPE()
					nodeIdx = 0
					break iifeSwitch_5
				}

				_ = nodeIdx
				var newIdomIdx int32 = int32(-1)
				_ = newIdomIdx
				{
					cur_7 := j2gort.NewIterCursor(j2gort.CheckNotNull(graph.GetDirectedPredNodes_1_any(node)))
					for cur_7.HasNext() {
						rawElem_8 := cur_7.Next()
						pred := j2gort.AsIface[DiGraph_DiGraphNode[N, E]](rawElem_8)
						_ = pred
						var predValue N = pred.GetValue()
						_ = predValue
						predIdx, predIdxOk_ := postOrderIndex.GetOk(predValue)
						_ = predIdxOk_
						_ = predIdx
						if (predIdxOk_) && (idomIndexes[func() int32 {
							if !predIdxOk_ {
								j2gort.ThrowNPE()
							}
							return predIdx
						}()] != int32(-1)) {
							{

								if !predIdxOk_ {
									j2gort.ThrowNPE()
								}
								newIdomIdx = predIdx
							}

							break
						}
					}
				}
				if newIdomIdx != int32(-1) {
					{
						cur_9 := j2gort.NewIterCursor(j2gort.CheckNotNull(graph.GetDirectedPredNodes_1_any(node)))
						for cur_9.HasNext() {
							rawElem_10 := cur_9.Next()
							pred_2 := j2gort.AsIface[DiGraph_DiGraphNode[N, E]](rawElem_10)
							_ = pred_2
							var predValue_2 N = pred_2.GetValue()
							_ = predValue_2
							predIdx_2, predIdx_2Ok_ := postOrderIndex.GetOk(predValue_2)
							_ = predIdx_2Ok_
							_ = predIdx_2
							if (predIdx_2Ok_) && (idomIndexes[func() int32 {
								if !predIdx_2Ok_ {
									j2gort.ThrowNPE()
								}
								return predIdx_2
							}()] != int32(-1)) {
								var iifeArg_1 int32
								{

									if !predIdx_2Ok_ {
										j2gort.ThrowNPE()
									}
									iifeArg_1 = predIdx_2
								}
								newIdomIdx = DominatorTree_Intersect(idomIndexes, iifeArg_1,

									newIdomIdx)
							}
						}
					}
					if idomIndexes[nodeIdx] != newIdomIdx {
						{
							a_11 := idomIndexes
							i_12 := nodeIdx
							a_11[i_12] = newIdomIdx
						}
						changed = true
					}
				}
			}
		}
	}
	var idomsMap *j2gort.Ext_java_util_LinkedHashMap = j2gort.NewExt_java_util_LinkedHashMap()
	_ = idomsMap
	var children j2gort.Ext_com_google_common_collect_ListMultimap = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ListMultimap](j2gort.Ext_com_google_common_collect_ArrayListMultimap_Create_0())
	_ = children
	{
		var i_2 int32 = int32(0)
		_ = i_2
		for ; i_2 < numNodes; i_2 = i_2 + 1 {
			var node_2 N
		iifeSwitch_3:
			switch {
			default:
				tv_13 := postOrder.Get(i_2)
				if tv_13 == nil {
					{

						if _, ok_ := any((*N)(nil)).(*string); ok_ {
							node_2 = any(j2gort.NullString).(N)
							break iifeSwitch_3
						}
						var z_ N
						node_2 = z_
						break iifeSwitch_3
					}

				}
				node_2 = j2gort.Cast[N](tv_13)
				break iifeSwitch_3
			}

			_ = node_2
			var idomIdx int32 = idomIndexes[i_2]
			_ = idomIdx
			if (idomIdx != int32(-1)) && (i_2 != startNodeIndex) {
				var idom N
			iifeSwitch_2:
				switch {
				default:
					tv_14 := postOrder.Get(idomIdx)
					if tv_14 == nil {
						{

							if _, ok_ := any((*N)(nil)).(*string); ok_ {
								idom = any(j2gort.NullString).(N)
								break iifeSwitch_2
							}
							var z_ N
							idom = z_
							break iifeSwitch_2
						}

					}
					idom = j2gort.Cast[N](tv_14)
					break iifeSwitch_2
				}

				_ = idom
				func() N {
					tv_15 := idomsMap.Put_2(j2gort.ToAny(node_2), j2gort.ToAny(idom))
					if tv_15 == nil {
						{

							if _, ok_ := any((*N)(nil)).(*string); ok_ {
								return any(j2gort.NullString).(N)
							}
							var z_ N
							return z_
						}

					}
					return j2gort.Cast[N](tv_15)
				}()
				children.Put(j2gort.ToAny(idom), j2gort.ToAny(node_2))
			}
		}
	}
	var sizesMap *j2gort.Ext_java_util_LinkedHashMap = j2gort.NewExt_java_util_LinkedHashMap()
	_ = sizesMap
	DominatorTree_ComputeSubtreeSizes[N](entry, children, sizesMap)
	return NewDominatorTree[N](idomsMap, sizesMap)
}

// Java static method: DominatorTree.intersect(int[] idomIndexes, int i, int j)
func DominatorTree_Intersect(idomIndexes []int32, i int32, j int32) int32 {
	_ = idomIndexes
	_ = i
	_ = j
	var finger1 int32 = i
	_ = finger1
	var finger2 int32 = j
	_ = finger2
	for finger1 != finger2 {
		for finger1 < finger2 {
			finger1 = idomIndexes[finger1]
		}
		for finger2 < finger1 {
			finger2 = idomIndexes[finger2]
		}
	}
	return finger1
}

// Java static method: DominatorTree.buildPostOrder(com.google.javascript.jscomp.graph.DiGraph<N,E> graph, N node, java.util.ArrayList<N> postOrder, java.util.LinkedHashMap<N,java.lang.Boolean> visited)
func DominatorTree_BuildPostOrder[N comparable, E any](graph DiGraph[N, E], node N, postOrder j2gort.Ext_java_util_ArrayList, visited *j2gort.Ext_java_util_LinkedHashMap) {
	_ = graph
	_ = node
	_ = postOrder
	_ = visited
	visited.Put_2(j2gort.ToAny(node), j2gort.ValueOfBoolean(true))
	{
		cur_16 := j2gort.NewIterCursor(j2gort.CheckNotNull(graph.GetDirectedSuccNodes_1_any(node)))
		for cur_16.HasNext() {
			rawElem_17 := cur_16.Next()
			succ := j2gort.AsIface[DiGraph_DiGraphNode[N, E]](rawElem_17)
			_ = succ
			var succValue N = succ.GetValue()
			_ = succValue
			if !(visited.ContainsKey(j2gort.ToAny(succValue))) {
				DominatorTree_BuildPostOrder[N, E](graph, succValue, postOrder, visited)
			}
		}
	}
	postOrder.Add_1(j2gort.ToAny(node))
}

// Java static method: DominatorTree.computeSubtreeSizes(N node, com.google.common.collect.ListMultimap<N,N> children, java.util.LinkedHashMap<N,java.lang.Integer> sizes)
func DominatorTree_ComputeSubtreeSizes[N comparable](node N, children j2gort.Ext_com_google_common_collect_ListMultimap, sizes *j2gort.Ext_java_util_LinkedHashMap) int32 {
	_ = node
	_ = children
	_ = sizes
	var size int32 = int32(1)
	_ = size
	var nodeChildren j2gort.Ext_java_util_List = j2gort.AsIface[j2gort.Ext_java_util_List](children.Get(j2gort.ToAny(node)))
	_ = nodeChildren
	if nodeChildren != nil {
		{
			cur_18 := j2gort.NewIterCursor(j2gort.CheckNotNull(nodeChildren))
			for cur_18.HasNext() {
				rawElem_19 := cur_18.Next()
				var child N
			iifeSwitch_4:
				switch {
				default:

					tv_20 := rawElem_19
					if tv_20 == nil {
						{

							if _, ok_ := any((*N)(nil)).(*string); ok_ {
								child = any(j2gort.NullString).(N)
								break iifeSwitch_4
							}
							var z_ N
							child = z_
							break iifeSwitch_4
						}

					}
					child = j2gort.Cast[N](tv_20)
					break iifeSwitch_4
				}

				_ = child
				size = (size + DominatorTree_ComputeSubtreeSizes[N](child, children, sizes))
			}
		}
	}
	sizes.Put_2(j2gort.ToAny(node), j2gort.ValueOfInt(size))
	return size
}

// Java method: DominatorTree.getImmediateDominator(N node)
func (this_ *DominatorTree[N]) GetImmediateDominator(node N) N {
	_ = node
	{

		tv_21 := this_.idoms.Get(j2gort.ToAny(node))
		if tv_21 == nil {
			{

				if _, ok_ := any((*N)(nil)).(*string); ok_ {
					return any(j2gort.NullString).(N)
				}
				var z_ N
				return z_
			}

		}
		return j2gort.Cast[N](tv_21)
	}

}

// Java method: DominatorTree.getDominatedSubtreeSize(N node)
func (this_ *DominatorTree[N]) GetDominatedSubtreeSize(node N) int32 {
	_ = node
	return j2gort.UnboxInt(j2gort.CastPtr[j2gort.BoxInt](this_.subtreeSizes.GetOrDefault(j2gort.ToAny(node), j2gort.ValueOfInt(int32(0)))))
}

// Java method: DominatorTree.getAllSubtreeSizes()
func (this_ *DominatorTree[N]) GetAllSubtreeSizes() j2gort.Ext_com_google_common_collect_ImmutableMap {
	return this_.subtreeSizes
}
