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
// Port of com/google/javascript/jscomp/PruningAnalysis.java (generated mechanically by j2go).

package jscomp

import (
	"fmt"
	"sync/atomic"

	pkg_graph "github.com/paullepoulpe/jsgo/transpiled/gen/graph"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_pruning_analysis_1 j2gort.Ext_java_util_function_Function
var lambdaSingleton_pruning_analysis_2 j2gort.Ext_java_util_function_Function
var lambdaSingleton_pruning_analysis_3 j2gort.Ext_java_util_function_ToIntFunction

func init() {
	lambdaSingleton_pruning_analysis_1 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any { return j2gort.AsIface[j2gort.Ext_java_util_Map_Entry](mrp0_).GetKey() })
	lambdaSingleton_pruning_analysis_2 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any { return j2gort.AsIface[j2gort.Ext_java_util_Map_Entry](mrp0_).GetValue() })
	lambdaSingleton_pruning_analysis_3 = j2gort.NewLambda_Ext_java_util_function_ToIntFunction(func(lp_6 any) int32 {
		entry := j2gort.AsIface[j2gort.Ext_java_util_Map_Entry](lp_6)
		_ = entry
		return j2gort.UnboxInt(j2gort.CastPtr[j2gort.BoxInt](entry.GetValue()))
	})
}

// Java class: com.google.javascript.jscomp.PruningAnalysis
type PruningAnalysis struct {
	sorter      *SortedDependencies                               `gson:"sorter"`
	entryPoints j2gort.Ext_com_google_common_collect_ImmutableSet `gson:"entryPoints"`
	idHash_     atomic.Int32
}

func (this_ *PruningAnalysis) BindSelf_(s any) {
	_ = s
}
func (this_ *PruningAnalysis) IsJavaNil_() bool { return this_ == nil }
func (*PruningAnalysis) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.PruningAnalysis", "PruningAnalysis"
}
func (this_ *PruningAnalysis) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*PruningAnalysis) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.PruningAnalysis", "com.google.javascript.jscomp.PruningAnalysis", func(v_ any) bool {
		return j2gort.InstanceOf[*PruningAnalysis](v_)
	})
}

func StaticGet_PruningAnalysis_BOTTLENECK_THRESHOLD() int32 { return int32(25) }
func StaticGet_PruningAnalysis_VIRTUAL_ROOT() string        { return "ENTRY POINTS ROOT" }
func (this_ *PruningAnalysis) initFields_PruningAnalysis() {
	_ = this_
}

// Java constructor: PruningAnalysis(com.google.javascript.jscomp.deps.SortedDependencies<com.google.javascript.jscomp.CompilerInput> sorter, java.util.Collection<com.google.javascript.jscomp.CompilerInput> entryPoints)
func NewPruningAnalysis(sorter *SortedDependencies, entryPoints j2gort.Ext_java_util_Collection) *PruningAnalysis {
	this_ := &PruningAnalysis{}
	InitPruningAnalysis(this_, sorter, entryPoints)
	return this_
}

func InitPruningAnalysis(this_ *PruningAnalysis, sorter *SortedDependencies, entryPoints j2gort.Ext_java_util_Collection) {
	_ = sorter
	_ = entryPoints
	this_.initFields_PruningAnalysis()
	this_.sorter = sorter
	this_.entryPoints = j2gort.Ext_com_google_common_collect_ImmutableSet_CopyOf_1_Ext_java_util_Collection(entryPoints)
}

// Java static method: PruningAnalysis.create(com.google.javascript.jscomp.deps.SortedDependencies<com.google.javascript.jscomp.CompilerInput> sorter, java.util.Collection<com.google.javascript.jscomp.CompilerInput> entryPoints)
func PruningAnalysis_Create(sorter *SortedDependencies, entryPoints j2gort.Ext_java_util_Collection) *PruningAnalysis {
	_ = sorter
	_ = entryPoints
	return NewPruningAnalysis(sorter, entryPoints)
}

// Java method: PruningAnalysis.analyze()
func (this_ *PruningAnalysis) Analyze() *PruningAnalysis_Result {
	var total j2gort.Ext_com_google_common_collect_ImmutableMap = this_.priv_PruningAnalysis_calculateTotalBlame_4585()
	_ = total
	var bottlenecks j2gort.Ext_com_google_common_collect_ImmutableMap = this_.priv_PruningAnalysis_calculateBottlenecks_4586()
	_ = bottlenecks
	return NewPruningAnalysis_Result(total, bottlenecks)
}

// Java method: PruningAnalysis.calculateTotalBlame()
func (this_ *PruningAnalysis) priv_PruningAnalysis_calculateTotalBlame_4585() j2gort.Ext_com_google_common_collect_ImmutableMap {
	var reachabilityPerEntryPoint j2gort.Ext_com_google_common_collect_ImmutableMap_Builder = j2gort.Ext_com_google_common_collect_ImmutableMap_S_Builder()
	_ = reachabilityPerEntryPoint
	{
		cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.entryPoints))
		for cur_1.HasNext() {
			rawElem_2 := cur_1.Next()
			provider := j2gort.CastPtr[CompilerInput](rawElem_2)
			_ = provider
			var transitiveClosure j2gort.Ext_com_google_common_collect_ImmutableSet = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet_Builder](j2gort.Ext_com_google_common_collect_ImmutableSet_S_Builder().AddAll_1_Ext_java_lang_Iterable(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.sorter.GetStrongDependenciesOf(j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.ImmListOfPtr(provider)), false)))).AddAll_1_Ext_java_lang_Iterable(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.sorter.GetSortedWeakDependenciesOf(j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.ImmListOfPtr(provider)))))).Build())
			_ = transitiveClosure
			reachabilityPerEntryPoint.Put_2(j2gort.BoxString(provider.GetName()), j2gort.ValueOfInt(transitiveClosure.Size()))
		}
	}
	return j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap](reachabilityPerEntryPoint.BuildOrThrow())
}

// Java method: PruningAnalysis.calculateBottlenecks()
func (this_ *PruningAnalysis) priv_PruningAnalysis_calculateBottlenecks_4586() j2gort.Ext_com_google_common_collect_ImmutableMap {
	var graph pkg_graph.DiGraph[string, string] = this_.priv_PruningAnalysis_buildGraph_4587()
	_ = graph
	var dominatorTree *pkg_graph.DominatorTree[string] = pkg_graph.DominatorTree_Compute[string, string](graph, "ENTRY POINTS ROOT")
	_ = dominatorTree
	var bottleneckBlame j2gort.OrderedMap[string, *j2gort.BoxInt] = j2gort.NewOrderedMap[string, *j2gort.BoxInt]()
	_ = bottleneckBlame
	{
		cur_3 := j2gort.NewMapEntryCursor(dominatorTree.GetAllSubtreeSizes())
		for cur_3.NextEntry() {
			var node string = cur_3.KeyStr()
			_ = node
			var retainedCount int32 = j2gort.UnboxInt(j2gort.CastPtr[j2gort.BoxInt](cur_3.GetValue()))
			_ = retainedCount
			{
				var iifeCond_1 bool
				{

					sr_4 := node
					se_5 := "ENTRY POINTS ROOT"
					iifeCond_1 = j2gort.CheckNotNullString(sr_4) == se_5
				}
				if !(iifeCond_1) && (retainedCount > int32(1)) {
					bottleneckBlame.Put(node, j2gort.ValueOfInt(retainedCount))
				}
			}
		}
	}
	return j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap](j2gort.AsIface[j2gort.Ext_java_util_Set](bottleneckBlame.ToJavaLinkedHashMap().EntrySet()).Stream().Sorted_1(j2gort.Ext_java_util_Comparator_ComparingInt(lambdaSingleton_pruning_analysis_3).Reversed()).Limit(int64(int32(25))).Collect_1(j2gort.Ext_com_google_common_collect_ImmutableMap_ToImmutableMap_2(lambdaSingleton_pruning_analysis_1, lambdaSingleton_pruning_analysis_2)))
}

// Java method: PruningAnalysis.buildGraph()
func (this_ *PruningAnalysis) priv_PruningAnalysis_buildGraph_4587() pkg_graph.DiGraph[string, string] {
	var graph pkg_graph.DiGraph[string, string] = j2gort.AsIface[pkg_graph.DiGraph[string, string]](pkg_graph.LinkedDirectedGraph_Create[string, string]())
	_ = graph
	graph.CreateNode("ENTRY POINTS ROOT")
	{
		cur_7 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.sorter.GetSortedList()))
		for cur_7.HasNext() {
			rawElem_8 := cur_7.Next()
			input := j2gort.CastPtr[CompilerInput](rawElem_8)
			_ = input
			graph.CreateNode(input.GetName())
		}
	}
	{
		cur_9 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.entryPoints))
		for cur_9.HasNext() {
			rawElem_10 := cur_9.Next()
			input_2 := j2gort.CastPtr[CompilerInput](rawElem_10)
			_ = input_2
			graph.Connect_3_any_any_any("ENTRY POINTS ROOT", j2gort.NullString, input_2.GetName())
		}
	}
	{
		cur_11 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.sorter.GetSortedList()))
		for cur_11.HasNext() {
			rawElem_12 := cur_11.Next()
			input_3 := j2gort.CastPtr[CompilerInput](rawElem_12)
			_ = input_3
			var name string = input_3.GetName()
			_ = name
			var depSymbols j2gort.LocalOrderedSet[string] = j2gort.NewLocalOrderedSet[string]()
			_ = depSymbols
			{
				cur_13 := j2gort.NewIterCursor(j2gort.CheckNotNull(input_3.GetRequires()))
				for cur_13.HasNext() {
					rawElem_14 := cur_13.Next()
					req := j2gort.CastPtr[DependencyInfo_Require](rawElem_14)
					_ = req
					depSymbols.Add(req.GetSymbol())
				}
			}
			j2gort.LocalOrderedSetAddAllStr(&depSymbols, j2gort.AsIface[j2gort.Ext_java_util_Collection](input_3.GetTypeRequires()))
			{
				cur_15 := depSymbols.Iter()
				for cur_15.HasNext() {
					symbol := cur_15.
						Next()
					_ = symbol
					func() {
						defer func() {
							if rec_16 := recover(); rec_16 != nil {
								jex_17 := javalang.RecoverException(rec_16)
								_ = jex_17
								if j2gort.InstanceOf[*SortedDependencies_MissingProvideException](jex_17) {
									e := j2gort.AsIface[*SortedDependencies_MissingProvideException](jex_17)
									_ = e
								} else {
									panic(rec_16)
								}
							}
						}()
						var dep *CompilerInput = j2gort.CastPtr[CompilerInput](this_.sorter.GetInputProviding(symbol))
						_ = dep
						if !(j2gort.CheckNotNullPtr(dep) == input_3) {
							graph.ConnectIfNotFound(name, j2gort.NullString, dep.GetName())
						}
					}()
				}
			}
		}
	}
	return graph
}

// Java record: com.google.javascript.jscomp.PruningAnalysis.Result
type PruningAnalysis_Result struct {
	entryPointDependencyCount j2gort.Ext_com_google_common_collect_ImmutableMap `gson:"entryPointDependencyCount"`
	bottleneckBlame           j2gort.Ext_com_google_common_collect_ImmutableMap `gson:"bottleneckBlame"`
}

func (this_ *PruningAnalysis_Result) BindSelf_(s any) {
	_ = s
}
func (this_ *PruningAnalysis_Result) IsJavaNil_() bool { return this_ == nil }
func (*PruningAnalysis_Result) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.PruningAnalysis$Result", "Result"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.PruningAnalysis.Result", "com.google.javascript.jscomp.PruningAnalysis$Result", func(v_ any) bool {
		return j2gort.InstanceOf[*PruningAnalysis_Result](v_)
	})
}
func (*PruningAnalysis_Result) Implements_j2gort_Ext_java_lang_Record() {}

func (this_ *PruningAnalysis_Result) initFields_PruningAnalysis_Result() {
	_ = this_
}

// Java constructor: Result(com.google.common.collect.ImmutableMap<java.lang.String,java.lang.Integer> entryPointDependencyCount, com.google.common.collect.ImmutableMap<java.lang.String,java.lang.Integer> bottleneckBlame)
func NewPruningAnalysis_Result(entryPointDependencyCount j2gort.Ext_com_google_common_collect_ImmutableMap, bottleneckBlame j2gort.Ext_com_google_common_collect_ImmutableMap) *PruningAnalysis_Result {
	this_ := &PruningAnalysis_Result{}
	InitPruningAnalysis_Result(this_, entryPointDependencyCount, bottleneckBlame)
	return this_
}

func InitPruningAnalysis_Result(this_ *PruningAnalysis_Result, entryPointDependencyCount j2gort.Ext_com_google_common_collect_ImmutableMap, bottleneckBlame j2gort.Ext_com_google_common_collect_ImmutableMap) {
	_ = entryPointDependencyCount
	_ = bottleneckBlame
	this_.initFields_PruningAnalysis_Result()
	this_.entryPointDependencyCount = entryPointDependencyCount
	this_.bottleneckBlame = bottleneckBlame
}

func (this_ *PruningAnalysis_Result) EntryPointDependencyCount() j2gort.Ext_com_google_common_collect_ImmutableMap {
	return this_.entryPointDependencyCount
}
func (this_ *PruningAnalysis_Result) BottleneckBlame() j2gort.Ext_com_google_common_collect_ImmutableMap {
	return this_.bottleneckBlame
}
func (this_ *PruningAnalysis_Result) Equals(other_ any) bool {
	o_, ok_ := other_.(*PruningAnalysis_Result)
	if !ok_ || o_ == nil {
		return false
	}
	if this_ == o_ {
		return true
	}
	if !javalang.ObjectsEquals(this_.entryPointDependencyCount, o_.entryPointDependencyCount) {
		return false
	}
	if !javalang.ObjectsEquals(this_.bottleneckBlame, o_.bottleneckBlame) {
		return false
	}
	return true
}
func (this_ *PruningAnalysis_Result) HashCode() int32 {
	var h_ int32 = 0
	h_ = 31*h_ + j2gort.ObjectHash(this_.entryPointDependencyCount)
	h_ = 31*h_ + j2gort.ObjectHash(this_.bottleneckBlame)
	return h_
}
func (this_ *PruningAnalysis_Result) String() string {
	_ = fmt.Sprint
	return "Result[" + "entryPointDependencyCount=" + j2gort.StrVal(this_.entryPointDependencyCount) + ", " + "bottleneckBlame=" + j2gort.StrVal(this_.bottleneckBlame) + "]"
}
