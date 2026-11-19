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
// Port of com/google/javascript/jscomp/JSChunkGraph.java (generated mechanically by j2go).

package jscomp

import (
	"fmt"
	"sync/atomic"

	pkg_base "github.com/paullepoulpe/jsgo/transpiled/gen/base"
	pkg_graph "github.com/paullepoulpe/jsgo/transpiled/gen/graph"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_js_chunk_graph_1 j2gort.Ext_com_google_common_base_Function
var lambdaSingleton_js_chunk_graph_2 j2gort.Ext_java_util_function_Function

func init() {
	lambdaSingleton_js_chunk_graph_1 = j2gort.NewLambda_Ext_com_google_common_base_Function(func(mrp0_ any) any { return j2gort.CastPtr[JSChunk](mrp0_).GetInputs() })
	lambdaSingleton_js_chunk_graph_2 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any { return j2gort.CastPtr[JSChunk](mrp0_).GetAllDependencies() })
}

// Java class: com.google.javascript.jscomp.JSChunkGraph
type JSChunkGraph struct {
	chunks                 []*JSChunk                                                          `gson:"chunks"`
	selfPlusTransitiveDeps [][]int64                                                           `gson:"selfPlusTransitiveDeps"`
	scratchCandidateWords  []int64                                                             `gson:"scratchCandidateWords"`
	subtreeSize            []int32                                                             `gson:"subtreeSize"`
	chunksByDepth          j2gort.Ext_java_util_List                                           `gson:"chunksByDepth"`
	dependencyMap          *pkg_base.LinkedIdentityHashMap[*JSChunk, j2gort.Ext_java_util_Set] `gson:"-"`
	idHash_                atomic.Int32
}

func (this_ *JSChunkGraph) BindSelf_(s any) {
	_ = s
}
func (this_ *JSChunkGraph) IsJavaNil_() bool { return this_ == nil }
func (*JSChunkGraph) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.JSChunkGraph", "JSChunkGraph"
}
func (this_ *JSChunkGraph) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*JSChunkGraph) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.JSChunkGraph", "com.google.javascript.jscomp.JSChunkGraph", func(v_ any) bool {
		return j2gort.InstanceOf[*JSChunkGraph](v_)
	})
}

var static_JSChunkGraph_WEAK_FILE_REACHABLE_FROM_ENTRY_POINT_ERROR *DiagnosticType

func StaticGet_JSChunkGraph_WEAK_FILE_REACHABLE_FROM_ENTRY_POINT_ERROR() *DiagnosticType {
	clinit_JSChunkGraph()
	return static_JSChunkGraph_WEAK_FILE_REACHABLE_FROM_ENTRY_POINT_ERROR
}
func StaticSet_JSChunkGraph_WEAK_FILE_REACHABLE_FROM_ENTRY_POINT_ERROR(v *DiagnosticType) {
	clinit_JSChunkGraph()
	static_JSChunkGraph_WEAK_FILE_REACHABLE_FROM_ENTRY_POINT_ERROR = v
}

var static_JSChunkGraph_EXPLICIT_WEAK_ENTRY_POINT_ERROR *DiagnosticType

func StaticGet_JSChunkGraph_EXPLICIT_WEAK_ENTRY_POINT_ERROR() *DiagnosticType {
	clinit_JSChunkGraph()
	return static_JSChunkGraph_EXPLICIT_WEAK_ENTRY_POINT_ERROR
}
func StaticSet_JSChunkGraph_EXPLICIT_WEAK_ENTRY_POINT_ERROR(v *DiagnosticType) {
	clinit_JSChunkGraph()
	static_JSChunkGraph_EXPLICIT_WEAK_ENTRY_POINT_ERROR = v
}

var static_JSChunkGraph_IMPLICIT_WEAK_ENTRY_POINT_ERROR *DiagnosticType

func StaticGet_JSChunkGraph_IMPLICIT_WEAK_ENTRY_POINT_ERROR() *DiagnosticType {
	clinit_JSChunkGraph()
	return static_JSChunkGraph_IMPLICIT_WEAK_ENTRY_POINT_ERROR
}
func StaticSet_JSChunkGraph_IMPLICIT_WEAK_ENTRY_POINT_ERROR(v *DiagnosticType) {
	clinit_JSChunkGraph()
	static_JSChunkGraph_IMPLICIT_WEAK_ENTRY_POINT_ERROR = v
}

var clinitDone_JSChunkGraph uint32
var clinitState_JSChunkGraph j2gort.ClinitState

func clinit_JSChunkGraph() {
	if atomic.LoadUint32(&clinitDone_JSChunkGraph) != 0 {
		return
	}
	clinitSlow_JSChunkGraph()
}

func clinitSlow_JSChunkGraph() {
	if !j2gort.EnterClinit(&clinitState_JSChunkGraph, &clinitDone_JSChunkGraph) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_JSChunkGraph, &clinitDone_JSChunkGraph)
	static_JSChunkGraph_WEAK_FILE_REACHABLE_FROM_ENTRY_POINT_ERROR = DiagnosticType_Error("JSC_WEAK_FILE_REACHABLE_FROM_ENTRY_POINT_ERROR", "File strongly reachable from an entry point must not be weak: {0}")
	static_JSChunkGraph_EXPLICIT_WEAK_ENTRY_POINT_ERROR = DiagnosticType_Error("JSC_EXPLICIT_WEAK_ENTRY_POINT_ERROR", "Explicit entry point input must not be weak: {0}")
	static_JSChunkGraph_IMPLICIT_WEAK_ENTRY_POINT_ERROR = DiagnosticType_Warning("JSC_IMPLICIT_WEAK_ENTRY_POINT_ERROR", "Implicit entry point input should not be weak: {0}")
}

func (this_ *JSChunkGraph) initFields_JSChunkGraph() {
	_ = this_
	this_.dependencyMap = pkg_base.NewLinkedIdentityHashMap[*JSChunk, j2gort.Ext_java_util_Set]()
}

// Java constructor: JSChunkGraph(com.google.javascript.jscomp.JSChunk[] chunksInDepOrder)
func NewJSChunkGraph_1_JSChunk(chunksInDepOrder []*JSChunk) *JSChunkGraph {
	clinit_JSChunkGraph()
	this_ := &JSChunkGraph{}
	InitJSChunkGraph_1_JSChunk(this_, chunksInDepOrder)
	return this_
}

func InitJSChunkGraph_1_JSChunk(this_ *JSChunkGraph, chunksInDepOrder []*JSChunk) {
	_ = chunksInDepOrder
	var iifeArg_1 []any
iifeSwitch_1:
	switch {
	default:

		arr_1 := chunksInDepOrder
		if arr_1 == nil {
			iifeArg_1 = nil
			break iifeSwitch_1
		}
		out_2 := make([]any, len(arr_1))
		for i_, el_3 := range arr_1 {
			out_2[i_] = j2gort.PtrToAny(el_3)
		}
		iifeArg_1 = out_2
		break iifeSwitch_1
	}
	InitJSChunkGraph_1_Ext_java_util_List(this_, j2gort.Ext_java_util_Arrays_AsList(iifeArg_1))
}

// Java constructor: JSChunkGraph(java.util.List<com.google.javascript.jscomp.JSChunk> chunksInDepOrder)
func NewJSChunkGraph_1_Ext_java_util_List(chunksInDepOrder j2gort.Ext_java_util_List) *JSChunkGraph {
	clinit_JSChunkGraph()
	this_ := &JSChunkGraph{}
	InitJSChunkGraph_1_Ext_java_util_List(this_, chunksInDepOrder)
	return this_
}

func InitJSChunkGraph_1_Ext_java_util_List(this_ *JSChunkGraph, chunksInDepOrder j2gort.Ext_java_util_List) {
	_ = chunksInDepOrder
	this_.initFields_JSChunkGraph()
	if !(!(chunksInDepOrder.IsEmpty())) {
		j2gort.CheckState(false)
	}
	chunksInDepOrder = this_.priv_JSChunkGraph_makeWeakChunk_3530(chunksInDepOrder)
	this_.chunks = make([]*JSChunk, chunksInDepOrder.Size())
	{
		var chunkIndex int32 = int32(0)
		_ = chunkIndex
		for ; chunkIndex < int32(len(this_.chunks)); chunkIndex = chunkIndex + 1 {
			var chunk *JSChunk = j2gort.CastPtr[JSChunk](chunksInDepOrder.Get(chunkIndex))
			_ = chunk
			if !(chunk.GetIndex() == int32(-1)) {
				j2gort.CheckStateFmt(false, "Chunk index already set: %s", j2gort.PtrToAny(chunk))
			}
			chunk.SetIndex(chunkIndex)
			{
				a_4 := this_.chunks
				i_5 := chunkIndex
				a_4[i_5] = chunk
			}
		}
	}
	this_.chunksByDepth = this_.priv_JSChunkGraph_initChunksByDepth_3529()
	this_.selfPlusTransitiveDeps = this_.priv_JSChunkGraph_initTransitiveDepsWords_3531()
	this_.scratchCandidateWords = make([]int64, javalang.Ushr32((int32(len(this_.chunks))+int32(63)), int(int32(6))))
	this_.subtreeSize = this_.priv_JSChunkGraph_initSubtreeSize_3532()
	JSChunkGraph_MoveMarkedWeakSources(this_.GetChunkByName("$weak$"), this_.GetAllInputs())
}

// Java method: JSChunkGraph.initChunksByDepth()
func (this_ *JSChunkGraph) priv_JSChunkGraph_initChunksByDepth_3529() j2gort.Ext_java_util_List {
	var tmpChunksByDepth []j2gort.Ext_java_util_List = []j2gort.Ext_java_util_List{}
	_ = tmpChunksByDepth
	{
		var chunkIndex int32 = int32(0)
		_ = chunkIndex
		for ; chunkIndex < int32(len(this_.chunks)); chunkIndex = chunkIndex + 1 {
			var chunk *JSChunk = this_.chunks[chunkIndex]
			_ = chunk
			if !(chunk.GetDepth() == int32(-1)) {
				j2gort.CheckStateFmt(false, "Chunk depth already set: %s", j2gort.PtrToAny(chunk))
			}
			var depth int32 = int32(0)
			_ = depth
			{
				cur_6 := j2gort.NewIterCursor(j2gort.CheckNotNull(chunk.GetDependencies()))
				for cur_6.HasNext() {
					rawElem_7 := cur_6.Next()
					dep := j2gort.CastPtr[JSChunk](rawElem_7)
					_ = dep
					var depDepth int32 = dep.GetDepth()
					_ = depDepth
					if depDepth < int32(0) {
						panic(NewJSChunkGraph_ChunkDependenceException(j2gort.Ext_java_lang_String_Format_2("Chunks not in dependency order: %s preceded %s", []any{j2gort.BoxString(chunk.GetName()), j2gort.BoxString(dep.GetName())}), chunk, dep))
					}
					depth = max(depth, (depDepth + int32(1)))
				}
			}
			chunk.SetDepth(depth)
			if depth == int32(len(tmpChunksByDepth)) {
				tmpChunksByDepth = append(tmpChunksByDepth, j2gort.NewExt_java_util_ArrayList())
			}
			tmpChunksByDepth[depth].Add_1(j2gort.PtrToAny(chunk))
		}
	}
	return j2gort.WrapTypedArrayList(tmpChunksByDepth)
}

// Java method: JSChunkGraph.makeWeakChunk(java.util.List<com.google.javascript.jscomp.JSChunk> chunksInDepOrder)
func (this_ *JSChunkGraph) priv_JSChunkGraph_makeWeakChunk_3530(chunksInDepOrder j2gort.Ext_java_util_List) j2gort.Ext_java_util_List {
	_ = chunksInDepOrder
	var hasWeakChunk bool = false
	_ = hasWeakChunk
	{
		cur_8 := j2gort.NewIterCursor(j2gort.CheckNotNull(chunksInDepOrder))
		for cur_8.HasNext() {
			rawElem_9 := cur_8.Next()
			chunk := j2gort.CastPtr[JSChunk](rawElem_9)
			_ = chunk
			{
				var iifeCond_1 bool
				{

					sr_10 := chunk.GetName()
					se_11 := "$weak$"
					iifeCond_1 = j2gort.CheckNotNullString(sr_10) == se_11
				}
				if iifeCond_1 {
					hasWeakChunk = true
					var allOtherChunks j2gort.Ext_java_util_Set = j2gort.NewExt_java_util_LinkedHashSet_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](chunksInDepOrder))
					_ = allOtherChunks
					allOtherChunks.Remove_1_any(j2gort.PtrToAny(chunk))
					if !(chunk.GetAllDependencies().ContainsAll(allOtherChunks)) {
						j2gort.CheckStateStr(false, "A weak chunk already exists but it does not depend on every other chunk.")
					}
					if !(chunk.GetAllDependencies().Size() == allOtherChunks.Size()) {
						j2gort.CheckStateStr(false, "The weak chunk cannot have extra dependencies.")
					}
					break
				}
			}
		}
	}
	if hasWeakChunk {
		var misplacedWeakFiles []string = []string{}
		_ = misplacedWeakFiles
		var misplacedStrongFiles []string = []string{}
		_ = misplacedStrongFiles
		{
			cur_12 := j2gort.NewIterCursor(j2gort.CheckNotNull(chunksInDepOrder))
			for cur_12.HasNext() {
				rawElem_13 := cur_12.Next()
				chunk_2 := j2gort.CastPtr[JSChunk](rawElem_13)
				_ = chunk_2
				var isWeakChunk bool
				{
					sr_14 := chunk_2.GetName()
					se_15 := "$weak$"
					isWeakChunk = j2gort.CheckNotNullString(sr_14) == se_15
				}

				_ = isWeakChunk
				{
					cur_16 := j2gort.NewIterCursor(j2gort.CheckNotNull(chunk_2.GetInputs()))
					for cur_16.HasNext() {
						rawElem_17 := cur_16.Next()
						input := j2gort.CastPtr[CompilerInput](rawElem_17)
						_ = input
						if isWeakChunk && !(input.GetSourceFile().IsWeak()) {
							misplacedStrongFiles = append(misplacedStrongFiles, input.GetSourceFile().GetName())
						} else if !(isWeakChunk) && input.GetSourceFile().IsWeak() {
							misplacedWeakFiles = append(misplacedWeakFiles, j2gort.StrValString(input.GetSourceFile().GetName())+" (in chunk "+j2gort.StrValString(chunk_2.GetName())+")")
						}
					}
				}
			}
		}
		if !(len(misplacedStrongFiles) == 0 && len(misplacedWeakFiles) == 0) {
			var sb j2gort.LocalStringBuilder
			sb.InitStr("A weak chunk exists but some sources are misplaced.")

			if len(misplacedStrongFiles) != 0 {
				sb.AppendStr("\nFound these strong sources in the weak chunk:\n  ").AppendStr(j2gort.Ext_com_google_common_base_Joiner_On_1_string("\n  ").Join_1_Ext_java_lang_Iterable(j2gort.SliceReadOnlyIterable(misplacedStrongFiles)))
			}
			if len(misplacedWeakFiles) != 0 {
				sb.AppendStr("\nFound these weak sources in other chunks:\n  ").AppendStr(j2gort.Ext_com_google_common_base_Joiner_On_1_string("\n  ").Join_1_Ext_java_lang_Iterable(j2gort.SliceReadOnlyIterable(misplacedWeakFiles)))
			}
			panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", j2gort.BoxString(sb.String()))))
		}
	} else {
		var val_18 JSChunk
		InitJSChunk(&val_18, "$weak$")
		var weakChunk *JSChunk = &val_18
		_ = weakChunk
		{
			cur_19 := j2gort.NewIterCursor(j2gort.CheckNotNull(chunksInDepOrder))
			for cur_19.HasNext() {
				rawElem_20 := cur_19.Next()
				chunk_3 := j2gort.CastPtr[JSChunk](rawElem_20)
				_ = chunk_3
				weakChunk.AddDependency(chunk_3)
			}
		}
		chunksInDepOrder = j2gort.NewExt_java_util_ArrayList_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](chunksInDepOrder))
		chunksInDepOrder.Add_1(weakChunk)
	}
	return chunksInDepOrder
}

// Java method: JSChunkGraph.initTransitiveDepsWords()
func (this_ *JSChunkGraph) priv_JSChunkGraph_initTransitiveDepsWords_3531() [][]int64 {
	var array [][]int64 = make([][]int64, int32(len(this_.chunks)))
	_ = array
	{
		var chunkIndex int32 = int32(0)
		_ = chunkIndex
		for ; chunkIndex < int32(len(this_.chunks)); chunkIndex = chunkIndex + 1 {
			var chunk *JSChunk = this_.chunks[chunkIndex]
			_ = chunk
			var wordLen int32 = ((javalang.Ushr32(chunkIndex, int(int32(6)))) + int32(1))
			_ = wordLen
			var deps []int64 = make([]int64, wordLen)
			_ = deps
			{
				a_21 := array
				i_22 := chunkIndex
				a_21[i_22] = deps
			}
			a_23 := deps
			i_24 := javalang.Ushr32(chunkIndex, int(int32(6)))
			a_23[i_24] = (a_23[i_24] | javalang.Shl64(int64(1), int(chunkIndex)))
			{
				cur_25 := j2gort.NewIterCursor(j2gort.CheckNotNull(chunk.GetDependencies()))
				for cur_25.HasNext() {
					rawElem_26 := cur_25.Next()
					dep := j2gort.CastPtr[JSChunk](rawElem_26)
					_ = dep
					var depWords []int64 = array[dep.GetIndex()]
					_ = depWords
					{
						var w int32 = int32(0)
						_ = w
						var len_ int32 = int32(len(depWords))
						_ = len_
						for ; w < len_; w = w + 1 {
							a_27 := deps
							i_28 := w
							a_27[i_28] = (a_27[i_28] | depWords[w])
						}
					}
				}
			}
		}
	}
	return array
}

// Java method: JSChunkGraph.initSubtreeSize()
func (this_ *JSChunkGraph) priv_JSChunkGraph_initSubtreeSize_3532() []int32 {
	var subtreeSize []int32 = make([]int32, int32(len(this_.chunks)))
	_ = subtreeSize
	{
		var dependentIndex int32 = int32(0)
		_ = dependentIndex
		for ; dependentIndex < int32(len(this_.chunks)); dependentIndex = dependentIndex + 1 {
			var dependencies []int64 = this_.selfPlusTransitiveDeps[dependentIndex]
			_ = dependencies
			{
				var w int32 = (int32(len(dependencies)) - int32(1))
				_ = w
				for ; w >= int32(0); w = w - 1 {
					var word int64 = dependencies[w]
					_ = word
					for word != int64(0) {
						var bit int32 = j2gort.BoxLong_NumberOfTrailingZeros(word)
						_ = bit
						a_29 := subtreeSize
						i_30 := ((javalang.Shl32(w, int(int32(6)))) + bit)
						a_29[i_30] = (a_29[i_30] + int32(1))
						word = (word & (word - int64(1)))
					}
				}
			}
		}
	}
	return subtreeSize
}

// Java method: JSChunkGraph.getAllInputs()
func (this_ *JSChunkGraph) GetAllInputs() j2gort.Ext_java_lang_Iterable {
	var iifeArg_2 []any
iifeSwitch_2:
	switch {
	default:

		arr_31 := this_.chunks
		if arr_31 == nil {
			iifeArg_2 = nil
			break iifeSwitch_2
		}
		out_32 := make([]any, len(arr_31))
		for i_, el_33 := range arr_31 {
			out_32[i_] = j2gort.PtrToAny(el_33)
		}
		iifeArg_2 = out_32
		break iifeSwitch_2
	}
	return j2gort.Ext_com_google_common_collect_Iterables_Concat_1_Ext_java_lang_Iterable(j2gort.Ext_com_google_common_collect_Iterables_Transform(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](j2gort.Ext_java_util_Arrays_AsList(iifeArg_2)), lambdaSingleton_js_chunk_graph_1))
}

// Java method: JSChunkGraph.getInputCount()
func (this_ *JSChunkGraph) GetInputCount() int32 {
	var count int32 = int32(0)
	_ = count
	for _, el_34 := range j2gort.CheckNotNullSlice(this_.chunks) {
		chunk := el_34
		_ = chunk
		count = (count + chunk.GetInputCount())
	}
	return count
}

// Java method: JSChunkGraph.getAllChunks()
func (this_ *JSChunkGraph) GetAllChunks() j2gort.Ext_java_lang_Iterable {
	var iifeArg_3 []any
iifeSwitch_3:
	switch {
	default:

		arr_35 := this_.chunks
		if arr_35 == nil {
			iifeArg_3 = nil
			break iifeSwitch_3
		}
		out_36 := make([]any, len(arr_35))
		for i_, el_37 := range arr_35 {
			out_36[i_] = j2gort.PtrToAny(el_37)
		}
		iifeArg_3 = out_36
		break iifeSwitch_3
	}
	return j2gort.AsIface[j2gort.Ext_java_lang_Iterable](j2gort.Ext_java_util_Arrays_AsList(iifeArg_3))
}

// Java method: JSChunkGraph.getIntegralChunkArrayForSerialization()
func (this_ *JSChunkGraph) GetIntegralChunkArrayForSerialization() []*JSChunk {
	return this_.chunks
}

// Java method: JSChunkGraph.getChunkByName(java.lang.String name)
func (this_ *JSChunkGraph) GetChunkByName(name string) *JSChunk {
	_ = name
	for _, el_38 := range j2gort.CheckNotNullSlice(this_.chunks) {
		m := el_38
		_ = m
		if j2gort.CheckNotNullString(m.GetName()) == name {
			return m
		}
	}
	return nil
}

// Java method: JSChunkGraph.getChunksByName()
func (this_ *JSChunkGraph) GetChunksByName() j2gort.Ext_java_util_Map {
	var result j2gort.Ext_java_util_Map = j2gort.NewExt_java_util_LinkedHashMap()
	_ = result
	for _, el_39 := range j2gort.CheckNotNullSlice(this_.chunks) {
		m := el_39
		_ = m
		j2gort.MapPutStr(result, m.GetName(), j2gort.PtrToAny(m))
	}
	return result
}

// Java method: JSChunkGraph.getChunkCount()
func (this_ *JSChunkGraph) GetChunkCount() int32 {
	return int32(len(this_.chunks))
}

// Java method: JSChunkGraph.getRootChunk()
func (this_ *JSChunkGraph) GetRootChunk() *JSChunk {
	return j2gort.CastPtr[JSChunk](j2gort.Ext_com_google_common_collect_Iterables_GetOnlyElement_1(j2gort.AsIface[j2gort.Ext_java_util_List](this_.chunksByDepth.Get(int32(0)))))
}

// Java method: JSChunkGraph.toJson()
func (this_ *JSChunkGraph) ToJson() *j2gort.Ext_com_google_gson_JsonArray {
	var chunks *j2gort.Ext_com_google_gson_JsonArray = j2gort.NewExt_com_google_gson_JsonArray()
	_ = chunks
	{
		cur_40 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.GetAllChunks()))
		for cur_40.HasNext() {
			rawElem_41 := cur_40.Next()
			chunk := j2gort.CastPtr[JSChunk](rawElem_41)
			_ = chunk
			var node *j2gort.Ext_com_google_gson_JsonObject = j2gort.NewExt_com_google_gson_JsonObject()
			_ = node
			node.Add("name", j2gort.NewExt_com_google_gson_JsonPrimitive_1_string(chunk.GetName()))
			var deps *j2gort.Ext_com_google_gson_JsonArray = j2gort.NewExt_com_google_gson_JsonArray()
			_ = deps
			node.Add("dependencies", deps)
			{
				cur_42 := j2gort.NewIterCursor(j2gort.CheckNotNull(chunk.GetDependencies()))
				for cur_42.HasNext() {
					rawElem_43 := cur_42.Next()
					m := j2gort.CastPtr[JSChunk](rawElem_43)
					_ = m
					deps.Add_1_Ext_com_google_gson_JsonElement(j2gort.NewExt_com_google_gson_JsonPrimitive_1_string(m.GetName()))
				}
			}
			var transitiveDeps *j2gort.Ext_com_google_gson_JsonArray = j2gort.NewExt_com_google_gson_JsonArray()
			_ = transitiveDeps
			node.Add("transitive-dependencies", transitiveDeps)
			{
				cur_44 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.GetTransitiveDepsDeepestFirst(chunk)))
				for cur_44.HasNext() {
					rawElem_45 := cur_44.Next()
					m_2 := j2gort.CastPtr[JSChunk](rawElem_45)
					_ = m_2
					transitiveDeps.Add_1_Ext_com_google_gson_JsonElement(j2gort.NewExt_com_google_gson_JsonPrimitive_1_string(m_2.GetName()))
				}
			}
			var inputs *j2gort.Ext_com_google_gson_JsonArray = j2gort.NewExt_com_google_gson_JsonArray()
			_ = inputs
			node.Add("inputs", inputs)
			{
				cur_46 := j2gort.NewIterCursor(j2gort.CheckNotNull(chunk.GetInputs()))
				for cur_46.HasNext() {
					rawElem_47 := cur_46.Next()
					input := j2gort.CastPtr[CompilerInput](rawElem_47)
					_ = input
					inputs.Add_1_Ext_com_google_gson_JsonElement(j2gort.NewExt_com_google_gson_JsonPrimitive_1_string(input.GetSourceFile().GetName()))
				}
			}
			chunks.Add_1_Ext_com_google_gson_JsonElement(node)
		}
	}
	return chunks
}

// Java method: JSChunkGraph.dependsOn(com.google.javascript.jscomp.JSChunk src, com.google.javascript.jscomp.JSChunk m)
func (this_ *JSChunkGraph) DependsOn(src *JSChunk, m *JSChunk) bool {
	_ = src
	_ = m
	if src == m {
		return false
	}
	var mIdx int32 = m.GetIndex()
	_ = mIdx
	var deps []int64 = this_.selfPlusTransitiveDeps[src.GetIndex()]
	_ = deps
	var w int32 = javalang.Ushr32(mIdx, int(int32(6)))
	_ = w
	return ((w < int32(len(deps))) && ((deps[w] & (javalang.Shl64(int64(1), int(mIdx)))) != int64(0)))
}

// Java method: JSChunkGraph.getSmallestCoveringSubtree(com.google.javascript.jscomp.JSChunk parentTree, java.util.BitSet dependentChunks)
func (this_ *JSChunkGraph) GetSmallestCoveringSubtree(parentTree *JSChunk, dependentChunks *j2gort.Ext_java_util_BitSet) *JSChunk {
	_ = parentTree
	_ = dependentChunks
	var firstDependent int32 = dependentChunks.NextSetBit(int32(0))
	_ = firstDependent
	if !(firstDependent >= int32(0)) {
		j2gort.CheckState(false)
	}
	var firstDeps []int64 = this_.selfPlusTransitiveDeps[firstDependent]
	_ = firstDeps
	var wordCount int32 = int32(len(firstDeps))
	_ = wordCount
	var candWords []int64 = this_.scratchCandidateWords
	_ = candWords
	func() {
		acSrc_48 := firstDeps
		acSP_49 := int32(0)
		acDst_50 := candWords
		acDP_51 := int32(0)
		acLen_52 := wordCount
		j2gort.CheckNotNullSlice(acSrc_48)
		j2gort.CheckNotNullSlice(acDst_50)
		if acSP_49 < 0 || acDP_51 < 0 || acLen_52 < 0 || int64(acSP_49)+int64(acLen_52) > int64(len(acSrc_48)) || int64(acDP_51)+int64(acLen_52) > int64(len(acDst_50)) {
			panic(j2gort.ConstructShimThrowable("java.lang.ArrayIndexOutOfBoundsException"))
		}
		copy(acDst_50[int(acDP_51):int(acDP_51)+int(acLen_52)], acSrc_48[int(acSP_49):int(acSP_49)+int(acLen_52)])
	}()
	{
		var dependentIndex int32 = dependentChunks.NextSetBit((firstDependent + int32(1)))
		_ = dependentIndex
		for ; dependentIndex >= int32(0); dependentIndex = dependentChunks.NextSetBit((dependentIndex + int32(1))) {
			var depWords []int64 = this_.selfPlusTransitiveDeps[dependentIndex]
			_ = depWords
			{
				var w int32 = int32(0)
				_ = w
				for ; w < wordCount; w = w + 1 {
					a_53 := candWords
					i_54 := w
					a_53[i_54] = (a_53[i_54] & depWords[w])
				}
			}
		}
	}
	var hasAnyCandidate bool = false
	_ = hasAnyCandidate
	{
		var w_2 int32 = int32(0)
		_ = w_2
		for ; w_2 < wordCount; w_2 = w_2 + 1 {
			if candWords[w_2] != int64(0) {
				hasAnyCandidate = true
				break
			}
		}
	}
	if !(hasAnyCandidate) {
		j2gort.CheckStateFmt(false, "No common dependency found for %s", j2gort.PtrToAny(dependentChunks))
	}
	var parentTreeIndex int32 = parentTree.GetIndex()
	_ = parentTreeIndex
	var parentWord int32 = javalang.Ushr32(parentTreeIndex, int(int32(6)))
	_ = parentWord
	var parentBit int64 = javalang.Shl64(int64(1), int(parentTreeIndex))
	_ = parentBit
	var bestCandidateIndex int32 = parentTreeIndex
	_ = bestCandidateIndex
	var bestSubtreeSize int32 = this_.subtreeSize[bestCandidateIndex]
	_ = bestSubtreeSize
	{
		var w_3 int32 = (wordCount - int32(1))
		_ = w_3
		for ; w_3 >= int32(0); w_3 = w_3 - 1 {
			var word int64 = int64(0)
			_ = word
			for (func() int64 { var asgn_55 int64 = candWords[w_3]; word = asgn_55; return asgn_55 }()) != int64(0) {
				var bitInWord int32 = (int32(63) - j2gort.BoxLong_NumberOfLeadingZeros(word))
				_ = bitInWord
				var candidateIndex int32 = ((javalang.Shl32(w_3, int(int32(6)))) + bitInWord)
				_ = candidateIndex
				var candDeps []int64 = this_.selfPlusTransitiveDeps[candidateIndex]
				_ = candDeps
				if (parentWord < int32(len(candDeps))) && ((candDeps[parentWord] & parentBit) != int64(0)) {
					{
						var k int32 = int32(0)
						_ = k
						var len_ int32 = int32(len(candDeps))
						_ = len_
						for ; k < len_; k = k + 1 {
							a_56 := candWords
							i_57 := k
							a_56[i_57] = (a_56[i_57] & ^(candDeps[k]))
						}
					}
					var candSize int32 = this_.subtreeSize[candidateIndex]
					_ = candSize
					if candSize < bestSubtreeSize {
						bestCandidateIndex = candidateIndex
						bestSubtreeSize = candSize
					}
				} else {
					{
						a_58 := candWords
						i_59 := w_3
						a_58[i_59] = (word & ^(javalang.Shl64(int64(1), int(bitInWord))))
					}
				}
			}
		}
	}
	return this_.chunks[bestCandidateIndex]
}

// Java method: JSChunkGraph.getDeepestCommonDependency(com.google.javascript.jscomp.JSChunk m1, com.google.javascript.jscomp.JSChunk m2)
func (this_ *JSChunkGraph) GetDeepestCommonDependency(m1 *JSChunk, m2 *JSChunk) *JSChunk {
	_ = m1
	_ = m2
	var m1Depth int32 = m1.GetDepth()
	_ = m1Depth
	var m2Depth int32 = m2.GetDepth()
	_ = m2Depth
	{
		var depth int32 = (min(m1Depth, m2Depth) - int32(1))
		_ = depth
		for ; depth >= int32(0); depth = depth - 1 {
			var chunksAtDepth j2gort.Ext_java_util_List = j2gort.AsIface[j2gort.Ext_java_util_List](this_.chunksByDepth.Get(depth))
			_ = chunksAtDepth
			{
				var i int32 = (chunksAtDepth.Size() - int32(1))
				_ = i
				for ; i >= int32(0); i = i - 1 {
					var m *JSChunk = j2gort.CastPtr[JSChunk](chunksAtDepth.Get(i))
					_ = m
					if this_.DependsOn(m1, m) && this_.DependsOn(m2, m) {
						return m
					}
				}
			}
		}
	}
	return nil
}

// Java method: JSChunkGraph.getDeepestCommonDependencyInclusive(com.google.javascript.jscomp.JSChunk m1, com.google.javascript.jscomp.JSChunk m2)
func (this_ *JSChunkGraph) GetDeepestCommonDependencyInclusive_2(m1 *JSChunk, m2 *JSChunk) *JSChunk {
	_ = m1
	_ = m2
	if (m2 == m1) || this_.DependsOn(m2, m1) {
		return m1
	} else if this_.DependsOn(m1, m2) {
		return m2
	}
	return this_.GetDeepestCommonDependency(m1, m2)
}

// Java method: JSChunkGraph.getDeepestCommonDependencyInclusive(java.util.Collection<com.google.javascript.jscomp.JSChunk> chunks)
func (this_ *JSChunkGraph) GetDeepestCommonDependencyInclusive_1(chunks j2gort.Ext_java_util_Collection) *JSChunk {
	_ = chunks
	var iter j2gort.Ext_java_util_Iterator = j2gort.AsIface[j2gort.Ext_java_util_Iterator](chunks.Iterator())
	_ = iter
	var dep *JSChunk = j2gort.CastPtr[JSChunk](iter.Next())
	_ = dep
	for iter.HasNext() {
		dep = this_.GetDeepestCommonDependencyInclusive_2(dep, j2gort.CastPtr[JSChunk](iter.Next()))
	}
	return dep
}

// Java method: JSChunkGraph.getTransitiveDepsDeepestFirst(com.google.javascript.jscomp.JSChunk m)
func (this_ *JSChunkGraph) GetTransitiveDepsDeepestFirst(m *JSChunk) j2gort.Ext_java_util_List {
	_ = m
	return StaticGet_JSChunkGraph_InverseDepthComparator_INSTANCE().SortedCopy(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.priv_JSChunkGraph_getTransitiveDeps_3548(m)))
}

// Java method: JSChunkGraph.getTransitiveDeps(com.google.javascript.jscomp.JSChunk m)
func (this_ *JSChunkGraph) priv_JSChunkGraph_getTransitiveDeps_3548(m *JSChunk) j2gort.Ext_java_util_Set {
	_ = m
	return this_.dependencyMap.ComputeIfAbsent(m, lambdaSingleton_js_chunk_graph_2)
}

// Java static method: JSChunkGraph.moveMarkedWeakSources(com.google.javascript.jscomp.JSChunk weakChunk, java.lang.Iterable<com.google.javascript.jscomp.CompilerInput> inputs)
func JSChunkGraph_MoveMarkedWeakSources(weakChunk *JSChunk, inputs j2gort.Ext_java_lang_Iterable) {
	clinit_JSChunkGraph()
	_ = weakChunk
	_ = inputs
	j2gort.CheckNotNullPtr(weakChunk)
	var allInputs []*CompilerInput = j2gort.ImmListCopyToSlice[*CompilerInput](inputs)
	_ = allInputs
	{

		for _, i := range allInputs {

			_ = i
			if i.GetSourceFile().IsWeak() {
				var existingChunk *JSChunk = i.GetChunk()
				_ = existingChunk
				if existingChunk == weakChunk {
					continue
				}
				if existingChunk != nil {
					existingChunk.Remove(i)
				}
				weakChunk.Add_1_CompilerInput(i)
			}
		}
	}
}

// Java method: JSChunkGraph.manageDependencies(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.DependencyOptions dependencyOptions)
func (this_ *JSChunkGraph) ManageDependencies(compiler AbstractCompiler, dependencyOptions *DependencyOptions) *JSChunkGraph_DependencyManagementResult {
	_ = compiler
	_ = dependencyOptions
	var originalInputs []*CompilerInput = j2gort.ImmListCopyToSlice[*CompilerInput](this_.GetAllInputs())
	_ = originalInputs
	var val_62 SortedDependencies
	InitSortedDependencies(&val_62, j2gort.ImmListFromTypedSlice(originalInputs))
	var sorter *SortedDependencies = &val_62
	_ = sorter
	var entryPointInputs j2gort.OrderedSet[*CompilerInput] = this_.priv_JSChunkGraph_createEntryPointInputs_3551(compiler, dependencyOptions, this_.GetAllInputs(), sorter)
	_ = entryPointInputs
	var inputsByProvide j2gort.OrderedMap[string, *j2gort.Ext_java_util_LinkedHashSet] = j2gort.NewOrderedMap[string, *j2gort.Ext_java_util_LinkedHashSet]()
	_ = inputsByProvide
	{

		for _, input := range originalInputs {

			_ = input
			{
				cur_65 := j2gort.NewIterCursor(j2gort.CheckNotNull(input.GetKnownProvides()))
				for cur_65.HasNext() {
					provide := cur_65.NextStr()
					_ = provide
					var set *j2gort.Ext_java_util_LinkedHashSet = inputsByProvide.Get(provide)
					_ = set
					if set == nil {
						set = j2gort.NewExt_java_util_LinkedHashSet()
						inputsByProvide.Put(provide, set)
					}
					set.Add_1(j2gort.PtrToAny(input))
				}
			}
			var chunkName string = input.GetPath().ToModuleName()
			_ = chunkName
			var chunkSet *j2gort.Ext_java_util_LinkedHashSet = inputsByProvide.Get(chunkName)
			_ = chunkSet
			if chunkSet == nil {
				chunkSet = j2gort.NewExt_java_util_LinkedHashSet()
				inputsByProvide.Put(chunkName, chunkSet)
			}
			chunkSet.Add_1(j2gort.PtrToAny(input))
		}
	}
	{

		for _, input_2 := range originalInputs {

			_ = input_2
			{
				cur_68 := j2gort.NewIterCursor(j2gort.CheckNotNull(input_2.GetDynamicRequires()))
				for cur_68.HasNext() {
					require := cur_68.NextStr()
					_ = require
					var provided *j2gort.Ext_java_util_LinkedHashSet = inputsByProvide.Get(require)
					_ = provided
					if provided != nil {
						{
							cur_69 := j2gort.NewIterCursor(j2gort.CheckNotNullPtr(provided))
							for cur_69.HasNext() {
								rawElem_70 := cur_69.Next()
								p := j2gort.CastPtr[CompilerInput](rawElem_70)
								_ = p
								entryPointInputs.Add(p)
							}
						}
					}
				}
			}
		}
	}
	{

		for _, input_3 := range originalInputs {

			_ = input_3
			{
				cur_73 := j2gort.NewIterCursor(j2gort.CheckNotNull(input_3.GetRequireDynamicImports()))
				for cur_73.HasNext() {
					require_2 := cur_73.NextStr()
					_ = require_2
					var provided_2 *j2gort.Ext_java_util_LinkedHashSet = inputsByProvide.Get(require_2)
					_ = provided_2
					if provided_2 != nil {
						{
							cur_74 := j2gort.NewIterCursor(j2gort.CheckNotNullPtr(provided_2))
							for cur_74.HasNext() {
								rawElem_75 := cur_74.Next()
								p_2 := j2gort.CastPtr[CompilerInput](rawElem_75)
								_ = p_2
								entryPointInputs.Add(p_2)
							}
						}
					}
				}
			}
		}
	}
	var absoluteOrder j2gort.Ext_com_google_common_collect_ImmutableList = sorter.GetStrongDependenciesOf(j2gort.ImmListFromTypedSlice(originalInputs), dependencyOptions.ShouldSort())
	_ = absoluteOrder
	var entryPointInputsPerChunk j2gort.Ext_com_google_common_collect_ListMultimap = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ListMultimap](j2gort.Ext_com_google_common_collect_LinkedListMultimap_Create_0())
	_ = entryPointInputsPerChunk
	{
		cur_76 := entryPointInputs.Iter()
		for cur_76.HasNext() {

			input_4 := cur_76.
				Next()
			_ = input_4
			var chunk *JSChunk = input_4.GetChunk()
			_ = chunk
			j2gort.CheckNotNullPtr(chunk)
			entryPointInputsPerChunk.Put(j2gort.PtrToAny(chunk), j2gort.PtrToAny(input_4))
		}
	}
	{
		cur_78 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.GetAllChunks()))
		for cur_78.HasNext() {
			rawElem_79 := cur_78.Next()
			chunk_2 := j2gort.CastPtr[JSChunk](rawElem_79)
			_ = chunk_2
			chunk_2.RemoveAll()
		}
	}
	var orderedInputs j2gort.Ext_java_util_List = j2gort.NewExt_java_util_ArrayList()
	_ = orderedInputs
	var reachedInputs j2gort.LocalOrderedSet[*CompilerInput] = j2gort.NewLocalOrderedSet[*CompilerInput]()
	_ = reachedInputs
	for _, el_80 := range j2gort.CheckNotNullSlice(this_.chunks) {
		chunk_3 := el_80
		_ = chunk_3
		var transitiveClosure j2gort.Ext_java_util_List = nil
		_ = transitiveClosure
		if dependencyOptions.ShouldSort() && dependencyOptions.ShouldPrune() {
			transitiveClosure = j2gort.NewExt_java_util_ArrayList()
			var inputsNotYetReached j2gort.OrderedSet[*CompilerInput] = j2gort.OrderedSetFromCol[*CompilerInput](j2gort.ImmListFromTypedSlice(originalInputs))
			_ = inputsNotYetReached
			{
				cur_81 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.AsIface[j2gort.Ext_java_util_List](entryPointInputsPerChunk.Get(j2gort.PtrToAny(chunk_3)))))
				for cur_81.HasNext() {
					rawElem_82 := cur_81.Next()
					entryPoint := j2gort.CastPtr[CompilerInput](rawElem_82)
					_ = entryPoint
					transitiveClosure.AddAll_1(j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.priv_JSChunkGraph_getDepthFirstDependenciesOf_3550(entryPoint, &inputsNotYetReached, &inputsByProvide)))
				}
			}
			{
				cur_83 := j2gort.NewIterCursor(j2gort.CheckNotNull(transitiveClosure))
				for cur_83.HasNext() {
					rawElem_84 := cur_83.Next()
					orderedInput := j2gort.CastPtr[CompilerInput](rawElem_84)
					_ = orderedInput
					if reachedInputs.Add(orderedInput) {
						orderedInputs.Add_1(j2gort.PtrToAny(orderedInput))
					}
				}
			}
		} else {
			transitiveClosure = j2gort.AsIface[j2gort.Ext_java_util_List](sorter.GetStrongDependenciesOf(j2gort.AsIface[j2gort.Ext_java_util_List](entryPointInputsPerChunk.Get(j2gort.PtrToAny(chunk_3))), dependencyOptions.ShouldSort()))
		}
		{
			cur_85 := j2gort.NewIterCursor(j2gort.CheckNotNull(transitiveClosure))
			for cur_85.HasNext() {
				rawElem_86 := cur_85.Next()
				input_5 := j2gort.CastPtr[CompilerInput](rawElem_86)
				_ = input_5
				if (dependencyOptions.ShouldPrune() && input_5.GetSourceFile().IsWeak()) && !(entryPointInputs.Contains(input_5)) {
					compiler.Report_1(JSError_Make_2(StaticGet_JSChunkGraph_WEAK_FILE_REACHABLE_FROM_ENTRY_POINT_ERROR(), []string{input_5.GetSourceFile().GetName()}))
				}
				var oldChunk *JSChunk = input_5.GetChunk()
				_ = oldChunk
				if oldChunk == nil {
					input_5.SetChunk(chunk_3)
				} else {
					input_5.SetChunk(nil)
					input_5.SetChunk(this_.GetDeepestCommonDependencyInclusive_2(oldChunk, chunk_3))
				}
			}
		}
	}
	if !(dependencyOptions.ShouldSort() && dependencyOptions.ShouldPrune()) || entryPointInputsPerChunk.IsEmpty() {
		orderedInputs = j2gort.AsIface[j2gort.Ext_java_util_List](absoluteOrder)
	}
	var weakChunk *JSChunk = this_.GetChunkByName("$weak$")
	_ = weakChunk
	j2gort.CheckNotNullPtr(weakChunk)
	if dependencyOptions.ShouldPrune() {
		var weakInputs []*CompilerInput = j2gort.ImmListToSliceNullable[*CompilerInput](sorter.GetSortedWeakDependenciesOf(orderedInputs))
		_ = weakInputs
		{

			for _, i := range weakInputs {

				_ = i
				if !(i.GetChunk() == nil) {
					j2gort.CheckState(false)
				}
				i.GetSourceFile().SetKind(pkg_rhino.StaticGet_StaticSourceFile_SourceKind_WEAK())
				i.SetChunk(weakChunk)
				weakChunk.Add_1_CompilerInput(i)
			}
		}
	} else {
		JSChunkGraph_MoveMarkedWeakSources(weakChunk, j2gort.ImmListFromTypedSlice(originalInputs))
	}
	{
		cur_89 := j2gort.NewIterCursor(j2gort.CheckNotNull(orderedInputs))
		for cur_89.HasNext() {
			rawElem_90 := cur_89.Next()
			input_6 := j2gort.CastPtr[CompilerInput](rawElem_90)
			_ = input_6
			var chunk_4 *JSChunk = input_6.GetChunk()
			_ = chunk_4
			if (chunk_4 != nil) && (chunk_4.GetByName(input_6.GetName()) == nil) {
				chunk_4.Add_1_CompilerInput(input_6)
			}
		}
	}
	return NewJSChunkGraph_DependencyManagementResult(j2gort.Ext_com_google_common_collect_ImmutableList_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](orderedInputs)), sorter, entryPointInputs.ToImmutableSet())
}

// Java method: JSChunkGraph.getDepthFirstDependenciesOf(com.google.javascript.jscomp.CompilerInput rootInput, java.util.LinkedHashSet<com.google.javascript.jscomp.CompilerInput> unreachedInputs, java.util.LinkedHashMap<java.lang.String,java.util.LinkedHashSet<com.google.javascript.jscomp.CompilerInput>> inputsByProvide)
func (this_ *JSChunkGraph) priv_JSChunkGraph_getDepthFirstDependenciesOf_3550(rootInput *CompilerInput, unreachedInputs *j2gort.OrderedSet[*CompilerInput], inputsByProvide *j2gort.OrderedMap[string, *j2gort.Ext_java_util_LinkedHashSet]) j2gort.Ext_java_util_List {
	_ = rootInput
	_ = unreachedInputs
	_ = inputsByProvide
	var orderedInputs []*CompilerInput = []*CompilerInput{}
	_ = orderedInputs
	if !(unreachedInputs.Remove(rootInput)) {
		return j2gort.WrapTypedArrayList(orderedInputs)
	}
	var requires []*DependencyInfo_Require = j2gort.ImmListToSliceNullable[*DependencyInfo_Require](rootInput.GetRequires())
	_ = requires
	{
		var r int32 = int32(0)
		_ = r
		var rLen int32 = int32(len(requires))
		_ = rLen
		for ; r < rLen; r = r + 1 {
			var providers *j2gort.Ext_java_util_LinkedHashSet = inputsByProvide.Get(requires[r].GetSymbol())
			_ = providers
			if providers != nil {
				{
					cur_91 := j2gort.NewIterCursor(j2gort.CheckNotNullPtr(providers))
					for cur_91.HasNext() {
						rawElem_92 := cur_91.Next()
						input := j2gort.CastPtr[CompilerInput](rawElem_92)
						_ = input
						if j2gort.CheckNotNullPtr(unreachedInputs).Contains(input) {
							orderedInputs = j2gort.AppendCollection[*CompilerInput](orderedInputs, j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.priv_JSChunkGraph_getDepthFirstDependenciesOf_3550(input, unreachedInputs, inputsByProvide)))
						}
					}
				}
			}
		}
	}
	orderedInputs = append(orderedInputs, rootInput)
	return j2gort.WrapTypedArrayList(orderedInputs)
}

// Java method: JSChunkGraph.createEntryPointInputs(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.DependencyOptions dependencyOptions, java.lang.Iterable<com.google.javascript.jscomp.CompilerInput> inputs, com.google.javascript.jscomp.deps.SortedDependencies<com.google.javascript.jscomp.CompilerInput> sorter)
func (this_ *JSChunkGraph) priv_JSChunkGraph_createEntryPointInputs_3551(compiler AbstractCompiler, dependencyOptions *DependencyOptions, inputs j2gort.Ext_java_lang_Iterable, sorter *SortedDependencies) j2gort.OrderedSet[*CompilerInput] {
	_ = compiler
	_ = dependencyOptions
	_ = inputs
	_ = sorter
	var entryPointInputs j2gort.OrderedSet[*CompilerInput] = j2gort.NewOrderedSet[*CompilerInput]()
	_ = entryPointInputs
	var chunksByName j2gort.Ext_java_util_Map = this_.GetChunksByName()
	_ = chunksByName
	if dependencyOptions.ShouldPrune() {
		var baseJs *CompilerInput = j2gort.CastPtr[CompilerInput](sorter.MaybeGetInputProviding("goog"))
		_ = baseJs
		if baseJs != nil {
			entryPointInputs.Add(baseJs)
		}
		if !(dependencyOptions.ShouldDropMoochers()) {
			{
				cur_93 := j2gort.NewIterCursor(j2gort.CheckNotNull(sorter.GetInputsWithoutProvides()))
				for cur_93.HasNext() {
					rawElem_94 := cur_93.Next()
					entryPointInput := j2gort.CastPtr[CompilerInput](rawElem_94)
					_ = entryPointInput
					if entryPointInput.GetSourceFile().IsWeak() {
						compiler.Report_1(JSError_Make_2(StaticGet_JSChunkGraph_IMPLICIT_WEAK_ENTRY_POINT_ERROR(), []string{entryPointInput.GetSourceFile().GetName()}))
					} else {
						entryPointInputs.Add(entryPointInput)
					}
				}
			}
		}
		{

			for _, entryPoint := range j2gort.CheckNotNullSlice(j2gort.CheckNotNullPtr(dependencyOptions).
				Slice_EntryPoints()) {

				_ = entryPoint
				var entryPointInput_2 *CompilerInput = nil
				_ = entryPointInput_2
				func() {
					defer func() {
						if rec_97 := recover(); rec_97 != nil {
							jex_98 := javalang.RecoverException(rec_97)
							_ = jex_98
							if j2gort.InstanceOf[*SortedDependencies_MissingProvideException](jex_98) {
								e := j2gort.AsIface[*SortedDependencies_MissingProvideException](jex_98)
								_ = e
								panic(NewSortedDependencies_MissingProvideException_2(entryPoint.Name(), j2gort.AsIface[javalang.IException](e)))
							} else {
								panic(rec_97)
							}
						}
					}()
					{
						var iifeCond_1 bool
						{

							sr_99 := entryPoint.ClosureNamespace()
							se_100 := entryPoint.ModuleName()
							iifeCond_1 = j2gort.CheckNotNullString(sr_99) == se_100
						}
						if iifeCond_1 {
							entryPointInput_2 = j2gort.CastPtr[CompilerInput](sorter.MaybeGetInputProviding(entryPoint.ClosureNamespace()))
							if entryPointInput_2 == nil {
								entryPointInput_2 = j2gort.CastPtr[CompilerInput](sorter.GetInputProviding(entryPoint.Name()))
							}
						} else {
							var chunk *JSChunk = j2gort.CastPtr[JSChunk](j2gort.MapGetStr(chunksByName, entryPoint.ModuleName()))
							_ = chunk
							if chunk == nil {
								panic(NewJSChunkGraph_MissingChunkException(entryPoint.ModuleName()))
							} else {
								entryPointInput_2 = j2gort.CastPtr[CompilerInput](sorter.GetInputProviding(entryPoint.ClosureNamespace()))
								entryPointInput_2.OverrideChunk(chunk)
							}
						}
					}
				}()
				if entryPointInput_2.GetSourceFile().IsWeak() {
					compiler.Report_1(JSError_Make_2(StaticGet_JSChunkGraph_EXPLICIT_WEAK_ENTRY_POINT_ERROR(), []string{entryPointInput_2.GetSourceFile().GetName()}))
				} else {
					entryPointInputs.Add(entryPointInput_2)
				}
			}
		}
	} else {
		{
			cur_101 := j2gort.NewIterCursor(j2gort.CheckNotNull(inputs))
			for cur_101.HasNext() {
				rawElem_102 := cur_101.Next()
				input := j2gort.CastPtr[CompilerInput](rawElem_102)
				_ = input
				entryPointInputs.Add(input)
			}
		}
	}
	return entryPointInputs
}

// Java method: JSChunkGraph.toGraphvizGraph()
func (this_ *JSChunkGraph) ToGraphvizGraph() pkg_graph.LinkedDirectedGraph[*JSChunk, string] {
	var graphViz pkg_graph.LinkedDirectedGraph[*JSChunk, string] = pkg_graph.LinkedDirectedGraph_Create[*JSChunk, string]()
	_ = graphViz
	{
		cur_103 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.GetAllChunks()))
		for cur_103.HasNext() {
			rawElem_104 := cur_103.Next()
			chunk := j2gort.CastPtr[JSChunk](rawElem_104)
			_ = chunk
			graphViz.CreateNode(chunk)
			{
				cur_105 := j2gort.NewIterCursor(j2gort.CheckNotNull(chunk.GetDependencies()))
				for cur_105.HasNext() {
					rawElem_106 := cur_105.Next()
					dep := j2gort.CastPtr[JSChunk](rawElem_106)
					_ = dep
					graphViz.CreateNode(dep)
					graphViz.Connect_3_any_any_any(chunk, "->", dep)
				}
			}
		}
	}
	return graphViz
}

// Java static method: JSChunkGraph.depthCompare(com.google.javascript.jscomp.JSChunk m1, com.google.javascript.jscomp.JSChunk m2)
func JSChunkGraph_DepthCompare(m1 *JSChunk, m2 *JSChunk) int32 {
	clinit_JSChunkGraph()
	_ = m1
	_ = m2
	if m1 == m2 {
		return int32(0)
	}
	var d1 int32 = m1.GetDepth()
	_ = d1
	var d2 int32 = m2.GetDepth()
	_ = d2
	if d1 < d2 {
		return int32(-1)
	} else {
		{

			if d2 == d1 {
				return j2gort.JStrCompareTo(m1.GetName(), m2.GetName())
			}
			return int32(1)
		}

	}
}

// Java record: com.google.javascript.jscomp.JSChunkGraph.DependencyManagementResult
type JSChunkGraph_DependencyManagementResult struct {
	orderedInputs    j2gort.Ext_com_google_common_collect_ImmutableList `gson:"orderedInputs"`
	sorter           *SortedDependencies                                `gson:"sorter"`
	entryPointInputs j2gort.Ext_com_google_common_collect_ImmutableSet  `gson:"entryPointInputs"`
}

func (this_ *JSChunkGraph_DependencyManagementResult) BindSelf_(s any) {
	_ = s
}
func (this_ *JSChunkGraph_DependencyManagementResult) IsJavaNil_() bool { return this_ == nil }
func (*JSChunkGraph_DependencyManagementResult) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.JSChunkGraph$DependencyManagementResult", "DependencyManagementResult"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.JSChunkGraph.DependencyManagementResult", "com.google.javascript.jscomp.JSChunkGraph$DependencyManagementResult", func(v_ any) bool {
		return j2gort.InstanceOf[*JSChunkGraph_DependencyManagementResult](v_)
	})
}
func (*JSChunkGraph_DependencyManagementResult) Implements_j2gort_Ext_java_lang_Record() {}

func (this_ *JSChunkGraph_DependencyManagementResult) initFields_JSChunkGraph_DependencyManagementResult() {
	_ = this_
}

// Java constructor: DependencyManagementResult(com.google.common.collect.ImmutableList<com.google.javascript.jscomp.CompilerInput> orderedInputs, com.google.javascript.jscomp.deps.SortedDependencies<com.google.javascript.jscomp.CompilerInput> sorter, com.google.common.collect.ImmutableSet<com.google.javascript.jscomp.CompilerInput> entryPointInputs)
func NewJSChunkGraph_DependencyManagementResult(orderedInputs j2gort.Ext_com_google_common_collect_ImmutableList, sorter *SortedDependencies, entryPointInputs j2gort.Ext_com_google_common_collect_ImmutableSet) *JSChunkGraph_DependencyManagementResult {
	this_ := &JSChunkGraph_DependencyManagementResult{}
	InitJSChunkGraph_DependencyManagementResult(this_, orderedInputs, sorter, entryPointInputs)
	return this_
}

func InitJSChunkGraph_DependencyManagementResult(this_ *JSChunkGraph_DependencyManagementResult, orderedInputs j2gort.Ext_com_google_common_collect_ImmutableList, sorter *SortedDependencies, entryPointInputs j2gort.Ext_com_google_common_collect_ImmutableSet) {
	_ = orderedInputs
	_ = sorter
	_ = entryPointInputs
	this_.initFields_JSChunkGraph_DependencyManagementResult()
	this_.orderedInputs = orderedInputs
	this_.sorter = sorter
	this_.entryPointInputs = entryPointInputs
}

func (this_ *JSChunkGraph_DependencyManagementResult) OrderedInputs() j2gort.Ext_com_google_common_collect_ImmutableList {
	return this_.orderedInputs
}
func (this_ *JSChunkGraph_DependencyManagementResult) Sorter() *SortedDependencies {
	return this_.sorter
}
func (this_ *JSChunkGraph_DependencyManagementResult) EntryPointInputs() j2gort.Ext_com_google_common_collect_ImmutableSet {
	return this_.entryPointInputs
}
func (this_ *JSChunkGraph_DependencyManagementResult) Equals(other_ any) bool {
	o_, ok_ := other_.(*JSChunkGraph_DependencyManagementResult)
	if !ok_ || o_ == nil {
		return false
	}
	if this_ == o_ {
		return true
	}
	if !javalang.ObjectsEquals(this_.orderedInputs, o_.orderedInputs) {
		return false
	}
	if this_.sorter != o_.sorter {
		return false
	}
	if !javalang.ObjectsEquals(this_.entryPointInputs, o_.entryPointInputs) {
		return false
	}
	return true
}
func (this_ *JSChunkGraph_DependencyManagementResult) HashCode() int32 {
	var h_ int32 = 0
	h_ = 31*h_ + j2gort.ObjectHash(this_.orderedInputs)
	{
		hp_107 := this_.sorter
		h_ = 31 * h_
		if hp_107 != nil {
			h_ += j2gort.IdentityHash(hp_107)
		}
	}
	h_ = 31*h_ + j2gort.ObjectHash(this_.entryPointInputs)
	return h_
}
func (this_ *JSChunkGraph_DependencyManagementResult) String() string {
	_ = fmt.Sprint
	return "DependencyManagementResult[" + "orderedInputs=" + j2gort.StrVal(this_.orderedInputs) + ", " + "sorter=" + j2gort.StrVal(j2gort.PtrToAny(this_.sorter)) + ", " + "entryPointInputs=" + j2gort.StrVal(this_.entryPointInputs) + "]"
}

// Java class: com.google.javascript.jscomp.JSChunkGraph.InverseDepthComparator
type JSChunkGraph_InverseDepthComparator struct {
	j2gort.Ext_com_google_common_collect_OrderingBase
	idHash_ atomic.Int32
}

func (this_ *JSChunkGraph_InverseDepthComparator) BindSelf_(s any) {
	_ = s
	this_.Ext_com_google_common_collect_OrderingBase.BindExtSelf_(s)
}
func (this_ *JSChunkGraph_InverseDepthComparator) IsJavaNil_() bool { return this_ == nil }
func (*JSChunkGraph_InverseDepthComparator) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.JSChunkGraph$InverseDepthComparator", "InverseDepthComparator"
}
func (this_ *JSChunkGraph_InverseDepthComparator) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.JSChunkGraph.InverseDepthComparator", "com.google.javascript.jscomp.JSChunkGraph$InverseDepthComparator", func(v_ any) bool {
		return j2gort.InstanceOf[*JSChunkGraph_InverseDepthComparator](v_)
	})
}
func (*JSChunkGraph_InverseDepthComparator) Implements_j2gort_Ext_com_google_common_collect_Ordering() {
}
func (*JSChunkGraph_InverseDepthComparator) Implements_j2gort_Ext_java_util_Comparator() {}

var static_JSChunkGraph_InverseDepthComparator_INSTANCE *JSChunkGraph_InverseDepthComparator

func StaticGet_JSChunkGraph_InverseDepthComparator_INSTANCE() *JSChunkGraph_InverseDepthComparator {
	clinit_JSChunkGraph_InverseDepthComparator()
	return static_JSChunkGraph_InverseDepthComparator_INSTANCE
}
func StaticSet_JSChunkGraph_InverseDepthComparator_INSTANCE(v *JSChunkGraph_InverseDepthComparator) {
	clinit_JSChunkGraph_InverseDepthComparator()
	static_JSChunkGraph_InverseDepthComparator_INSTANCE = v
}

var clinitDone_JSChunkGraph_InverseDepthComparator uint32
var clinitState_JSChunkGraph_InverseDepthComparator j2gort.ClinitState

func clinit_JSChunkGraph_InverseDepthComparator() {
	if atomic.LoadUint32(&clinitDone_JSChunkGraph_InverseDepthComparator) != 0 {
		return
	}
	clinitSlow_JSChunkGraph_InverseDepthComparator()
}

func clinitSlow_JSChunkGraph_InverseDepthComparator() {
	if !j2gort.EnterClinit(&clinitState_JSChunkGraph_InverseDepthComparator, &clinitDone_JSChunkGraph_InverseDepthComparator) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_JSChunkGraph_InverseDepthComparator, &clinitDone_JSChunkGraph_InverseDepthComparator)
	static_JSChunkGraph_InverseDepthComparator_INSTANCE = NewJSChunkGraph_InverseDepthComparator()
}

func (this_ *JSChunkGraph_InverseDepthComparator) initFields_JSChunkGraph_InverseDepthComparator() {
	_ = this_
}

// Java constructor: InverseDepthComparator()
func NewJSChunkGraph_InverseDepthComparator() *JSChunkGraph_InverseDepthComparator {
	clinit_JSChunkGraph_InverseDepthComparator()
	this_ := &JSChunkGraph_InverseDepthComparator{}
	this_.BindSelf_(this_)
	InitJSChunkGraph_InverseDepthComparator(this_)
	return this_
}

func InitJSChunkGraph_InverseDepthComparator(this_ *JSChunkGraph_InverseDepthComparator) {
	j2gort.InitExt_com_google_common_collect_Ordering(&this_.Ext_com_google_common_collect_OrderingBase)
	this_.initFields_JSChunkGraph_InverseDepthComparator()
}

// Java method: InverseDepthComparator.compare(com.google.javascript.jscomp.JSChunk m1, com.google.javascript.jscomp.JSChunk m2)
func (this_ *JSChunkGraph_InverseDepthComparator) Compare(arg_107 any, arg_108 any) int32 {
	m1 := j2gort.CastPtr[JSChunk](arg_107)
	_ = m1
	m2 := j2gort.CastPtr[JSChunk](arg_108)
	_ = m2
	return JSChunkGraph_DepthCompare(m2, m1)
}

// Java class: com.google.javascript.jscomp.JSChunkGraph.ChunkDependenceException
type JSChunkGraph_ChunkDependenceException struct {
	javalang.IllegalArgumentException
	chunk          *JSChunk `gson:"chunk"`
	dependentChunk *JSChunk `gson:"dependentChunk"`
	idHash_        atomic.Int32
}

func (this_ *JSChunkGraph_ChunkDependenceException) BindSelf_(s any) {
	_ = s
}
func (this_ *JSChunkGraph_ChunkDependenceException) IsJavaNil_() bool { return this_ == nil }
func (*JSChunkGraph_ChunkDependenceException) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.JSChunkGraph$ChunkDependenceException", "ChunkDependenceException"
}
func (this_ *JSChunkGraph_ChunkDependenceException) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.JSChunkGraph.ChunkDependenceException", "com.google.javascript.jscomp.JSChunkGraph$ChunkDependenceException", func(v_ any) bool {
		return j2gort.InstanceOf[*JSChunkGraph_ChunkDependenceException](v_)
	})
}

func StaticGet_JSChunkGraph_ChunkDependenceException_SerialVersionUID() int64 { return int64(1) }
func (this_ *JSChunkGraph_ChunkDependenceException) initFields_JSChunkGraph_ChunkDependenceException() {
	_ = this_
}

// Java constructor: ChunkDependenceException(java.lang.String message, com.google.javascript.jscomp.JSChunk chunk, com.google.javascript.jscomp.JSChunk dependentChunk)
func NewJSChunkGraph_ChunkDependenceException(message string, chunk *JSChunk, dependentChunk *JSChunk) *JSChunkGraph_ChunkDependenceException {
	this_ := &JSChunkGraph_ChunkDependenceException{}
	InitJSChunkGraph_ChunkDependenceException(this_, message, chunk, dependentChunk)
	return this_
}

func InitJSChunkGraph_ChunkDependenceException(this_ *JSChunkGraph_ChunkDependenceException, message string, chunk *JSChunk, dependentChunk *JSChunk) {
	_ = message
	_ = chunk
	_ = dependentChunk
	j2gort.InitUserThrowable(this_, "com.google.javascript.jscomp.JSChunkGraph.ChunkDependenceException", j2gort.BoxString(message))
	this_.initFields_JSChunkGraph_ChunkDependenceException()
	this_.chunk = chunk
	this_.dependentChunk = dependentChunk
}

// Java method: ChunkDependenceException.getChunk()
func (this_ *JSChunkGraph_ChunkDependenceException) GetChunk() *JSChunk {
	return this_.chunk
}

// Java method: ChunkDependenceException.getDependentChunk()
func (this_ *JSChunkGraph_ChunkDependenceException) GetDependentChunk() *JSChunk {
	return this_.dependentChunk
}

// Java class: com.google.javascript.jscomp.JSChunkGraph.MissingChunkException
type JSChunkGraph_MissingChunkException struct {
	javalang.JavaException
	idHash_ atomic.Int32
}

func (this_ *JSChunkGraph_MissingChunkException) BindSelf_(s any) {
	_ = s
}
func (this_ *JSChunkGraph_MissingChunkException) IsJavaNil_() bool { return this_ == nil }
func (*JSChunkGraph_MissingChunkException) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.JSChunkGraph$MissingChunkException", "MissingChunkException"
}
func (this_ *JSChunkGraph_MissingChunkException) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.JSChunkGraph.MissingChunkException", "com.google.javascript.jscomp.JSChunkGraph$MissingChunkException", func(v_ any) bool {
		return j2gort.InstanceOf[*JSChunkGraph_MissingChunkException](v_)
	})
}

func (this_ *JSChunkGraph_MissingChunkException) initFields_JSChunkGraph_MissingChunkException() {
	_ = this_
}

// Java constructor: MissingChunkException(java.lang.String chunkName)
func NewJSChunkGraph_MissingChunkException(chunkName string) *JSChunkGraph_MissingChunkException {
	this_ := &JSChunkGraph_MissingChunkException{}
	InitJSChunkGraph_MissingChunkException(this_, chunkName)
	return this_
}

func InitJSChunkGraph_MissingChunkException(this_ *JSChunkGraph_MissingChunkException, chunkName string) {
	_ = chunkName
	j2gort.InitUserThrowable(this_, "com.google.javascript.jscomp.JSChunkGraph.MissingChunkException", j2gort.BoxString(chunkName))
	this_.initFields_JSChunkGraph_MissingChunkException()
}
