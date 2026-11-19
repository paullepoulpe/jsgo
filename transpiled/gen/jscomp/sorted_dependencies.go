/*
 * Copyright 2010 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/deps/SortedDependencies.java (generated mechanically by j2go).

package jscomp

import (
	"sync"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.deps.SortedDependencies
type SortedDependencies struct {
	userOrderedInputs          []DependencyInfo                          `gson:"userOrderedInputs"`
	importOrderedInputs        []DependencyInfo                          `gson:"importOrderedInputs"`
	sortedListView             j2gort.Ext_java_util_List                 `gson:"sortedListView"`
	nonExportingInputs         j2gort.OrderedMap[string, DependencyInfo] `gson:"nonExportingInputs"`
	exportingInputBySymbolName j2gort.OrderedMap[string, DependencyInfo] `gson:"exportingInputBySymbolName"`
	uniqueInputs               []any                                     `gson:"uniqueInputs"`
	inputHashes                []int32                                   `gson:"inputHashes"`
	inputTable                 []int32                                   `gson:"inputTable"`
	uniqueInputCount           int32                                     `gson:"uniqueInputCount"`
	userInputIds               []int32                                   `gson:"userInputIds"`
	importInputIds             []int32                                   `gson:"importInputIds"`
	scratchStrongInUse         bool                                      `gson:"scratchStrongInUse"`
	scratchIncludedById        []bool                                    `gson:"scratchIncludedById"`
	scratchWorklist            *j2gort.Deque[DependencyInfo]             `gson:"scratchWorklist"`
	scratchWeakInUse           bool                                      `gson:"scratchWeakInUse"`
	scratchStrongById          []bool                                    `gson:"scratchStrongById"`
	scratchWeakById            []bool                                    `gson:"scratchWeakById"`
	scratchWeakWorklist        *j2gort.Deque[DependencyInfo]             `gson:"scratchWeakWorklist"`
	idHash_                    atomic.Int32
}

func (this_ *SortedDependencies) BindSelf_(s any) {
	_ = s
}
func (this_ *SortedDependencies) IsJavaNil_() bool { return this_ == nil }
func (*SortedDependencies) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.deps.SortedDependencies", "SortedDependencies"
}
func (this_ *SortedDependencies) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*SortedDependencies) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.deps.SortedDependencies", "com.google.javascript.jscomp.deps.SortedDependencies", func(v_ any) bool {
		return j2gort.InstanceOf[*SortedDependencies](v_)
	})
}

var staticSyncMu_SortedDependencies sync.Mutex

func (this_ *SortedDependencies) initFields_SortedDependencies() {
	_ = this_
	this_.userOrderedInputs = []DependencyInfo{}
	this_.importOrderedInputs = []DependencyInfo{}
	this_.nonExportingInputs = j2gort.NewOrderedMap[string, DependencyInfo]()
	this_.exportingInputBySymbolName = j2gort.NewOrderedMap[string, DependencyInfo]()
}

// Java constructor: SortedDependencies(java.util.List<InputT> userOrderedInputs)
func NewSortedDependencies(userOrderedInputs j2gort.Ext_java_util_List) *SortedDependencies {
	this_ := &SortedDependencies{}
	InitSortedDependencies(this_, userOrderedInputs)
	return this_
}

func InitSortedDependencies(this_ *SortedDependencies, userOrderedInputs j2gort.Ext_java_util_List) {
	_ = userOrderedInputs
	this_.initFields_SortedDependencies()
	{
		var i int32 = int32(0)
		_ = i
		var size int32 = userOrderedInputs.Size()
		_ = size
		for ; i < size; i = i + 1 {
			this_.userOrderedInputs = append(this_.userOrderedInputs, j2gort.AsIface[DependencyInfo](userOrderedInputs.Get(i)))
		}
	}
	this_.priv_SortedDependencies_processInputs_6162()
}

// Java static method: SortedDependencies.tryAcquireStrongScratch(com.google.javascript.jscomp.deps.SortedDependencies<?> self)
func SortedDependencies_TryAcquireStrongScratch(self *SortedDependencies) bool {
	staticSyncMu_SortedDependencies.Lock()
	defer staticSyncMu_SortedDependencies.Unlock()
	_ = self
	if self.scratchStrongInUse {
		return false
	}
	{
		r_1 := self
		r_1.scratchStrongInUse = true
	}
	return true
}

// Java static method: SortedDependencies.releaseStrongScratch(com.google.javascript.jscomp.deps.SortedDependencies<?> self)
func SortedDependencies_ReleaseStrongScratch(self *SortedDependencies) {
	staticSyncMu_SortedDependencies.Lock()
	defer staticSyncMu_SortedDependencies.Unlock()
	_ = self
	{
		r_2 := self
		r_2.scratchStrongInUse = false
	}
}

// Java static method: SortedDependencies.tryAcquireWeakScratch(com.google.javascript.jscomp.deps.SortedDependencies<?> self)
func SortedDependencies_TryAcquireWeakScratch(self *SortedDependencies) bool {
	staticSyncMu_SortedDependencies.Lock()
	defer staticSyncMu_SortedDependencies.Unlock()
	_ = self
	if self.scratchWeakInUse {
		return false
	}
	{
		r_3 := self
		r_3.scratchWeakInUse = true
	}
	return true
}

// Java static method: SortedDependencies.releaseWeakScratch(com.google.javascript.jscomp.deps.SortedDependencies<?> self)
func SortedDependencies_ReleaseWeakScratch(self *SortedDependencies) {
	staticSyncMu_SortedDependencies.Lock()
	defer staticSyncMu_SortedDependencies.Unlock()
	_ = self
	{
		r_4 := self
		r_4.scratchWeakInUse = false
	}
}

// Java static method: SortedDependencies.mixHash(int h)
func SortedDependencies_MixHash(h int32) int32 {
	_ = h
	h = (h * int32(-1640531527))
	return (h ^ (javalang.Ushr32(h, int(int32(16)))))
}

// Java method: SortedDependencies.registerOrGetInputId(InputT input)
func (this_ *SortedDependencies) priv_SortedDependencies_registerOrGetInputId_6150(input DependencyInfo) int32 {
	_ = input
	var h int32 = javalang.ObjectHashCode(j2gort.CheckNotNull(input))
	_ = h
	var mask int32 = (int32(len(this_.inputTable)) - int32(1))
	_ = mask
	var slot int32 = (SortedDependencies_MixHash(h) & mask)
	_ = slot
	for {
		var entry int32 = this_.inputTable[slot]
		_ = entry
		if entry == int32(0) {
			var id int32 = func() int32 {
				old_5 := this_.uniqueInputCount
				nw_6 := old_5 + 1
				this_.uniqueInputCount = nw_6
				return old_5
			}()
			_ = id
			{
				a_7 := this_.uniqueInputs
				i_8 := id
				a_7[i_8] = input
			}
			{
				a_9 := this_.inputHashes
				i_10 := id
				a_9[i_10] = h
			}
			{
				a_11 := this_.inputTable
				i_12 := slot
				a_11[i_12] = (id + int32(1))
			}
			return id
		}
		var id_2 int32 = (entry - int32(1))
		_ = id_2
		var existing any = this_.uniqueInputs[id_2]
		_ = existing
		if j2gort.RefEquals(existing, input) || ((this_.inputHashes[id_2] == h) && javalang.ObjectsEquals(j2gort.CheckNotNull(input), existing)) {
			return id_2
		}
		slot = ((slot + int32(1)) & mask)
	}
	panic("unreachable")
}

// Java method: SortedDependencies.findInputId(InputT input)
func (this_ *SortedDependencies) priv_SortedDependencies_findInputId_6151(input DependencyInfo) int32 {
	_ = input
	if (this_.uniqueInputCount == int32(0)) || (input == nil || input.IsJavaNil_()) {
		return int32(-1)
	}
	var h int32 = javalang.ObjectHashCode(j2gort.CheckNotNull(input))
	_ = h
	var mask int32 = (int32(len(this_.inputTable)) - int32(1))
	_ = mask
	var slot int32 = (SortedDependencies_MixHash(h) & mask)
	_ = slot
	for {
		var entry int32 = this_.inputTable[slot]
		_ = entry
		if entry == int32(0) {
			return int32(-1)
		}
		var id int32 = (entry - int32(1))
		_ = id
		var existing any = this_.uniqueInputs[id]
		_ = existing
		if j2gort.RefEquals(existing, input) || ((this_.inputHashes[id] == h) && javalang.ObjectsEquals(j2gort.CheckNotNull(input), existing)) {
			return id
		}
		slot = ((slot + int32(1)) & mask)
	}
	panic("unreachable")
}

// Java method: SortedDependencies.getStrongDependenciesOf(java.util.List<InputT> rootInputs, boolean sorted)
func (this_ *SortedDependencies) GetStrongDependenciesOf(rootInputs j2gort.Ext_java_util_List, sorted bool) j2gort.Ext_com_google_common_collect_ImmutableList {
	_ = rootInputs
	_ = sorted
	if rootInputs.IsEmpty() || (this_.uniqueInputCount == int32(0)) {
		return j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()
	}
	var ownsScratch bool = SortedDependencies_TryAcquireStrongScratch(this_)
	_ = ownsScratch
	var includedById []bool = nil
	_ = includedById
	var worklist *j2gort.Deque[DependencyInfo]
	_ = worklist
	if ownsScratch {
		includedById = this_.scratchIncludedById
		if (includedById == nil) || (int32(len(includedById)) < this_.uniqueInputCount) {
			includedById = make([]bool, this_.uniqueInputCount)
			this_.scratchIncludedById = includedById
		}
		worklist = this_.scratchWorklist
		if worklist == nil {
			worklist = &j2gort.Deque[DependencyInfo]{}
			this_.scratchWorklist = worklist
		}
		{
			var i int32 = int32(0)
			_ = i
			var size int32 = rootInputs.Size()
			_ = size
			for ; i < size; i = i + 1 {
				worklist.AddLast(j2gort.AsIface[DependencyInfo](rootInputs.Get(i)))
			}
		}
	} else {
		includedById = make([]bool, this_.uniqueInputCount)
		worklist = j2gort.DequePtrFromCol[DependencyInfo](j2gort.AsIface[j2gort.Ext_java_util_Collection](rootInputs))
	}
	defer func() {
		if ownsScratch {
			j2gort.ArraysFillRange(includedById, int32(0), this_.uniqueInputCount, false)
			worklist.Clear()
			SortedDependencies_ReleaseStrongScratch(this_)
		}
	}()
	{
		var extraIncluded *j2gort.Ext_java_util_LinkedHashSet = nil
		_ = extraIncluded
		for !(worklist.IsEmpty()) {
			var input DependencyInfo = worklist.Pop()
			_ = input
			var id int32 = this_.priv_SortedDependencies_findInputId_6151(input)
			_ = id
			var added bool = false
			_ = added
			if id >= int32(0) {
				added = !(includedById[id])
				{
					a_13 := includedById
					i_14 := id
					a_13[i_14] = true
				}
			} else {
				if extraIncluded == nil {
					extraIncluded = j2gort.NewExt_java_util_LinkedHashSet()
				}
				added = extraIncluded.Add_1(input)
			}
			if added {
				var requires []*DependencyInfo_Require = j2gort.ImmListToSliceNullable[*DependencyInfo_Require](input.GetRequires())
				_ = requires
				{
					var r int32 = int32(0)
					_ = r
					var rLen int32 = int32(len(requires))
					_ = rLen
					for ; r < rLen; r = r + 1 {
						var importedSymbolName DependencyInfo = this_.exportingInputBySymbolName.Get(requires[r].GetSymbol())
						_ = importedSymbolName
						if importedSymbolName != nil && !importedSymbolName.IsJavaNil_() {
							worklist.AddLast(importedSymbolName)
						}
					}
				}
			}
		}
		var builder []DependencyInfo = []DependencyInfo{}
		_ = builder
		if sorted {
			{
				var i_2 int32 = int32(0)
				_ = i_2
				var size_2 int32 = int32(len(this_.importOrderedInputs))
				_ = size_2
				for ; i_2 < size_2; i_2 = i_2 + 1 {
					if includedById[this_.importInputIds[i_2]] {
						builder = append(builder, j2gort.CheckNotNullElem(this_.importOrderedInputs[i_2]))
					}
				}
			}
		} else {
			{
				var i_3 int32 = int32(0)
				_ = i_3
				var size_3 int32 = int32(len(this_.userOrderedInputs))
				_ = size_3
				for ; i_3 < size_3; i_3 = i_3 + 1 {
					if includedById[this_.userInputIds[i_3]] {
						builder = append(builder, j2gort.CheckNotNullElem(this_.userOrderedInputs[i_3]))
					}
				}
			}
		}
		return j2gort.ImmListFromTypedSlice(builder)
	}
	panic("unreachable")
}

// Java method: SortedDependencies.getInputProviding(java.lang.String symbolName)
func (this_ *SortedDependencies) GetInputProviding(symbolName string) DependencyInfo {
	_ = symbolName
	var input DependencyInfo = this_.MaybeGetInputProviding(symbolName)
	_ = input
	if input != nil && !input.IsJavaNil_() {
		return input
	}
	panic(NewSortedDependencies_MissingProvideException_1(symbolName))
}

// Java method: SortedDependencies.getInputsWithoutProvides()
func (this_ *SortedDependencies) GetInputsWithoutProvides() j2gort.Ext_com_google_common_collect_ImmutableList {
	return this_.nonExportingInputs.ValuesToImmutableList()
}

// Java method: SortedDependencies.getSortedStrongDependenciesOf(java.util.List<InputT> roots)
func (this_ *SortedDependencies) GetSortedStrongDependenciesOf(roots j2gort.Ext_java_util_List) j2gort.Ext_com_google_common_collect_ImmutableList {
	_ = roots
	return this_.GetStrongDependenciesOf(roots, true)
}

// Java method: SortedDependencies.getSortedWeakDependenciesOf(java.util.List<InputT> rootInputs)
func (this_ *SortedDependencies) GetSortedWeakDependenciesOf(rootInputs j2gort.Ext_java_util_List) j2gort.Ext_com_google_common_collect_ImmutableList {
	_ = rootInputs
	var strongList []DependencyInfo = j2gort.ImmListToSliceNullable[DependencyInfo](this_.GetSortedStrongDependenciesOf(rootInputs))
	_ = strongList
	if len(strongList) == 0 || (this_.uniqueInputCount == int32(0)) {
		return j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()
	}
	var ownsScratch bool = SortedDependencies_TryAcquireWeakScratch(this_)
	_ = ownsScratch
	var strongById []bool = nil
	_ = strongById
	var weakById []bool = nil
	_ = weakById
	var worklist *j2gort.Deque[DependencyInfo]
	_ = worklist
	if ownsScratch {
		strongById = this_.scratchStrongById
		if (strongById == nil) || (int32(len(strongById)) < this_.uniqueInputCount) {
			strongById = make([]bool, this_.uniqueInputCount)
			this_.scratchStrongById = strongById
		}
		weakById = this_.scratchWeakById
		if (weakById == nil) || (int32(len(weakById)) < this_.uniqueInputCount) {
			weakById = make([]bool, this_.uniqueInputCount)
			this_.scratchWeakById = weakById
		}
		worklist = this_.scratchWeakWorklist
		if worklist == nil {
			worklist = j2gort.NewDequePtrCap[DependencyInfo](int32(len(strongList)))
			this_.scratchWeakWorklist = worklist
		}
	} else {
		strongById = make([]bool, this_.uniqueInputCount)
		weakById = make([]bool, this_.uniqueInputCount)
		worklist = j2gort.NewDequePtrCap[DependencyInfo](int32(len(strongList)))
	}
	defer func() {
		if ownsScratch {
			j2gort.ArraysFillRange(strongById, int32(0), this_.uniqueInputCount, false)
			j2gort.ArraysFillRange(weakById, int32(0), this_.uniqueInputCount, false)
			worklist.Clear()
			SortedDependencies_ReleaseWeakScratch(this_)
		}
	}()
	{
		{
			var i int32 = int32(0)
			_ = i
			var size int32 = int32(len(strongList))
			_ = size
			for ; i < size; i = i + 1 {
				var strongInput DependencyInfo = strongList[i]
				_ = strongInput
				var sId int32 = this_.priv_SortedDependencies_findInputId_6151(strongInput)
				_ = sId
				if !(strongById[sId]) {
					{
						a_15 := strongById
						i_16 := sId
						a_15[i_16] = true
					}
					worklist.AddLast(strongInput)
				}
			}
		}
		for !(worklist.IsEmpty()) {
			var input DependencyInfo = worklist.Pop()
			_ = input
			var id int32 = this_.priv_SortedDependencies_findInputId_6151(input)
			_ = id
			var isStrong bool = strongById[id]
			_ = isStrong
			if !(isStrong) {
				if weakById[id] {
					continue
				}
				{
					a_17 := weakById
					i_18 := id
					a_17[i_18] = true
				}
				this_.priv_SortedDependencies_enqueueWeakRequireEdges_6157(input.GetRequires(), strongById, weakById, worklist)
			}
			this_.priv_SortedDependencies_enqueueWeakEdges_6158(input.GetTypeRequires(), strongById, weakById, worklist)
		}
		var builder []DependencyInfo = []DependencyInfo{}
		_ = builder
		{
			var i_2 int32 = int32(0)
			_ = i_2
			var size_2 int32 = int32(len(this_.importOrderedInputs))
			_ = size_2
			for ; i_2 < size_2; i_2 = i_2 + 1 {
				if weakById[this_.importInputIds[i_2]] {
					builder = append(builder, j2gort.CheckNotNullElem(this_.importOrderedInputs[i_2]))
				}
			}
		}
		return j2gort.ImmListFromTypedSlice(builder)
	}
	panic("unreachable")
}

// Java method: SortedDependencies.enqueueWeakRequireEdges(com.google.common.collect.ImmutableList<com.google.javascript.jscomp.deps.DependencyInfo.Require> requires, boolean[] strongById, boolean[] weakById, java.util.ArrayDeque<InputT> worklist)
func (this_ *SortedDependencies) priv_SortedDependencies_enqueueWeakRequireEdges_6157(requires j2gort.Ext_com_google_common_collect_ImmutableList, strongById []bool, weakById []bool, worklist *j2gort.Deque[DependencyInfo]) {
	_ = requires
	_ = strongById
	_ = weakById
	_ = worklist
	{
		var i int32 = int32(0)
		_ = i
		var len_ int32 = requires.Size()
		_ = len_
		for ; i < len_; i = i + 1 {
			var importedSymbolName DependencyInfo = this_.exportingInputBySymbolName.Get(j2gort.CastPtr[DependencyInfo_Require](requires.Get(i)).GetSymbol())
			_ = importedSymbolName
			if importedSymbolName != nil && !importedSymbolName.IsJavaNil_() {
				var impId int32 = this_.priv_SortedDependencies_findInputId_6151(importedSymbolName)
				_ = impId
				if !(strongById[impId]) && !(weakById[impId]) {
					worklist.AddLast(importedSymbolName)
				}
			}
		}
	}
}

// Java method: SortedDependencies.enqueueWeakEdges(com.google.common.collect.ImmutableList<java.lang.String> symbols, boolean[] strongById, boolean[] weakById, java.util.ArrayDeque<InputT> worklist)
func (this_ *SortedDependencies) priv_SortedDependencies_enqueueWeakEdges_6158(symbols j2gort.Ext_com_google_common_collect_ImmutableList, strongById []bool, weakById []bool, worklist *j2gort.Deque[DependencyInfo]) {
	_ = symbols
	_ = strongById
	_ = weakById
	_ = worklist
	{
		var i int32 = int32(0)
		_ = i
		var len_ int32 = symbols.Size()
		_ = len_
		for ; i < len_; i = i + 1 {
			var importedSymbolName DependencyInfo = this_.exportingInputBySymbolName.Get(j2gort.ListGetStr(symbols, i))
			_ = importedSymbolName
			if importedSymbolName != nil && !importedSymbolName.IsJavaNil_() {
				var impId int32 = this_.priv_SortedDependencies_findInputId_6151(importedSymbolName)
				_ = impId
				if !(strongById[impId]) && !(weakById[impId]) {
					worklist.AddLast(importedSymbolName)
				}
			}
		}
	}
}

// Java method: SortedDependencies.getSortedList()
func (this_ *SortedDependencies) GetSortedList() j2gort.Ext_java_util_List {
	var view j2gort.Ext_java_util_List = this_.sortedListView
	_ = view
	if view == nil {
		var copy_ []DependencyInfo = j2gort.MakeNonNilSlice[DependencyInfo](int32(len(this_.importOrderedInputs)))
		_ = copy_
		{
			var i int32 = int32(0)
			_ = i
			var size int32 = int32(len(this_.importOrderedInputs))
			_ = size
			for ; i < size; i = i + 1 {
				copy_ = append(copy_, this_.importOrderedInputs[i])
			}
		}
		view = j2gort.SliceUnmodifiableList(copy_)
		this_.sortedListView = view
	}
	return view
}

// Java method: SortedDependencies.maybeGetInputProviding(java.lang.String symbol)
func (this_ *SortedDependencies) MaybeGetInputProviding(symbol string) DependencyInfo {
	_ = symbol
	if this_.exportingInputBySymbolName.ContainsKey(symbol) {
		return this_.exportingInputBySymbolName.Get(symbol)
	}
	return this_.nonExportingInputs.Get(ModuleNames_FileToModuleName(symbol))
}

// Java method: SortedDependencies.orderInput(InputT input, int id, boolean[] completedById, java.util.ArrayList<InputT>[] importedById, int[][] importedIdsById, int[] importedCounts)
func (this_ *SortedDependencies) priv_SortedDependencies_orderInput_6161(input DependencyInfo, id int32, completedById []bool, importedById []j2gort.Ext_java_util_ArrayList, importedIdsById [][]int32, importedCounts []int32) {
	_ = input
	_ = id
	_ = completedById
	_ = importedById
	_ = importedIdsById
	_ = importedCounts
	if completedById[id] {
		return
	}
	{
		a_19 := completedById
		i_20 := id
		a_19[i_20] = true
	}
	var depCount int32 = importedCounts[id]
	_ = depCount
	if depCount > int32(0) {
		var deps j2gort.Ext_java_util_ArrayList = importedById[id]
		_ = deps
		var depIds []int32 = importedIdsById[id]
		_ = depIds
		{
			var i int32 = int32(0)
			_ = i
			for ; i < depCount; i = i + 1 {
				this_.priv_SortedDependencies_orderInput_6161(j2gort.AsIface[DependencyInfo](deps.Get(i)), depIds[i], completedById, importedById, importedIdsById, importedCounts)
			}
		}
	}
	this_.importOrderedInputs = append(this_.importOrderedInputs, input)
}

// Java method: SortedDependencies.processInputs()
func (this_ *SortedDependencies) priv_SortedDependencies_processInputs_6162() {
	var inputCount int32 = int32(len(this_.userOrderedInputs))
	_ = inputCount
	if inputCount == int32(0) {
		this_.uniqueInputs = make([]any, int32(0))
		this_.inputHashes = make([]int32, int32(0))
		this_.inputTable = make([]int32, int32(0))
		this_.uniqueInputCount = int32(0)
		this_.userInputIds = make([]int32, int32(0))
		this_.importInputIds = make([]int32, int32(0))
		return
	}
	var tableCap int32 = int32(16)
	_ = tableCap
	for (int64(tableCap) * int64(int32(3))) < (int64(inputCount) * int64(int32(4))) {
		tableCap = javalang.Shl32(tableCap, int(int32(1)))
	}
	this_.uniqueInputs = make([]any, inputCount)
	this_.inputHashes = make([]int32, inputCount)
	this_.inputTable = make([]int32, tableCap)
	this_.userInputIds = make([]int32, inputCount)
	{
		var i int32 = int32(0)
		_ = i
		for ; i < inputCount; i = i + 1 {
			var userOrderedInput DependencyInfo = this_.userOrderedInputs[i]
			_ = userOrderedInput
			{
				a_21 := this_.userInputIds
				i_22 := i
				a_21[i_22] = this_.priv_SortedDependencies_registerOrGetInputId_6150(userOrderedInput)
			}
			var provides []string = j2gort.ImmListToSliceNullable[string](userOrderedInput.GetProvides())
			_ = provides
			var firstProvide string
			_ = firstProvide
			if len(provides) == 0 {
				firstProvide = j2gort.NullString
			} else {
				firstProvide = provides[int32(0)]
			}
			if (firstProvide == j2gort.NullString) || (((int32(len(provides)) == int32(1)) && javalang.StringStartsWith(j2gort.CheckNotNullString(firstProvide), "module$")) && !(userOrderedInput.IsEs6Module())) {
				this_.nonExportingInputs.Put(ModuleNames_FileToModuleName(userOrderedInput.GetName()), userOrderedInput)
			}
			{

				for _, providedSymbolName := range provides {

					_ = providedSymbolName
					this_.exportingInputBySymbolName.Put(providedSymbolName, userOrderedInput)
				}
			}
		}
	}
	var importedById []j2gort.Ext_java_util_ArrayList = make([]j2gort.Ext_java_util_ArrayList, this_.uniqueInputCount)
	_ = importedById
	var importedIdsById [][]int32 = make([][]int32, this_.uniqueInputCount)
	_ = importedIdsById
	var importedCounts []int32 = make([]int32, this_.uniqueInputCount)
	_ = importedCounts
	{
		var i_2 int32 = int32(0)
		_ = i_2
		for ; i_2 < inputCount; i_2 = i_2 + 1 {
			var userOrderedInput_2 DependencyInfo = this_.userOrderedInputs[i_2]
			_ = userOrderedInput_2
			var srcId int32 = this_.userInputIds[i_2]
			_ = srcId
			var requires []*DependencyInfo_Require = j2gort.ImmListToSliceNullable[*DependencyInfo_Require](userOrderedInput_2.GetRequires())
			_ = requires
			{
				var r int32 = int32(0)
				_ = r
				var rLen int32 = int32(len(requires))
				_ = rLen
				for ; r < rLen; r = r + 1 {
					var importedInput DependencyInfo = this_.exportingInputBySymbolName.Get(requires[r].GetSymbol())
					_ = importedInput
					if importedInput != nil && !importedInput.IsJavaNil_() {
						var dstId int32 = this_.priv_SortedDependencies_findInputId_6151(importedInput)
						_ = dstId
						var count int32 = importedCounts[srcId]
						_ = count
						var dstIds []int32 = importedIdsById[srcId]
						_ = dstIds
						var alreadyPresent bool = false
						_ = alreadyPresent
						if dstIds == nil {
							dstIds = make([]int32, int32(4))
							{
								a_24 := importedIdsById
								i_25 := srcId
								a_24[i_25] = dstIds
							}
							{
								a_26 := importedById
								i_27 := srcId
								a_26[i_27] = j2gort.NewExt_java_util_ArrayList_1_int32(int32(4))
							}
						} else {
							{
								var k int32 = int32(0)
								_ = k
								for ; k < count; k = k + 1 {
									if dstIds[k] == dstId {
										alreadyPresent = true
										break
									}
								}
							}
						}
						if !(alreadyPresent) {
							if count == int32(len(dstIds)) {
								dstIds = func() []int32 {
									arrSrc_28 := dstIds
									arrLen_29 := (count * int32(2))
									j2gort.CheckNotNullSlice(arrSrc_28)
									if arrLen_29 < 0 {
										panic(j2gort.ConstructShimThrowable("java.lang.NegativeArraySizeException"))
									}
									arrOut_30 := make([]int32, arrLen_29)
									copy(arrOut_30, arrSrc_28)
									return arrOut_30
								}()
								{
									a_31 := importedIdsById
									i_32 := srcId
									a_31[i_32] = dstIds
								}
							}
							{
								a_33 := dstIds
								i_34 := count
								a_33[i_34] = dstId
							}
							importedById[srcId].Add_1(importedInput)
							{
								a_35 := importedCounts
								i_36 := srcId
								a_35[i_36] = (count + int32(1))
							}
						}
					}
				}
			}
		}
	}
	var completedById []bool = make([]bool, this_.uniqueInputCount)
	_ = completedById
	{
		var i_3 int32 = int32(0)
		_ = i_3
		for ; i_3 < inputCount; i_3 = i_3 + 1 {
			this_.priv_SortedDependencies_orderInput_6161(this_.userOrderedInputs[i_3], this_.userInputIds[i_3], completedById, importedById, importedIdsById, importedCounts)
		}
	}
	var importCount int32 = int32(len(this_.importOrderedInputs))
	_ = importCount
	this_.importInputIds = make([]int32, importCount)
	{
		var i_4 int32 = int32(0)
		_ = i_4
		for ; i_4 < importCount; i_4 = i_4 + 1 {
			{
				a_37 := this_.importInputIds
				i_38 := i_4
				a_37[i_38] = this_.priv_SortedDependencies_findInputId_6151(this_.importOrderedInputs[i_4])
			}
		}
	}
}

// Java class: com.google.javascript.jscomp.deps.SortedDependencies.MissingProvideException
type SortedDependencies_MissingProvideException struct {
	javalang.JavaException
	idHash_ atomic.Int32
}

func (this_ *SortedDependencies_MissingProvideException) BindSelf_(s any) {
	_ = s
}
func (this_ *SortedDependencies_MissingProvideException) IsJavaNil_() bool { return this_ == nil }
func (*SortedDependencies_MissingProvideException) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.deps.SortedDependencies$MissingProvideException", "MissingProvideException"
}
func (this_ *SortedDependencies_MissingProvideException) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.deps.SortedDependencies.MissingProvideException", "com.google.javascript.jscomp.deps.SortedDependencies$MissingProvideException", func(v_ any) bool {
		return j2gort.InstanceOf[*SortedDependencies_MissingProvideException](v_)
	})
}

func (this_ *SortedDependencies_MissingProvideException) initFields_SortedDependencies_MissingProvideException() {
	_ = this_
}

// Java constructor: MissingProvideException(java.lang.String provide)
func NewSortedDependencies_MissingProvideException_1(provide string) *SortedDependencies_MissingProvideException {
	this_ := &SortedDependencies_MissingProvideException{}
	InitSortedDependencies_MissingProvideException_1(this_, provide)
	return this_
}

func InitSortedDependencies_MissingProvideException_1(this_ *SortedDependencies_MissingProvideException, provide string) {
	_ = provide
	j2gort.InitUserThrowable(this_, "com.google.javascript.jscomp.deps.SortedDependencies.MissingProvideException", j2gort.BoxString(provide))
	this_.initFields_SortedDependencies_MissingProvideException()
}

// Java constructor: MissingProvideException(java.lang.String provide, java.lang.Exception e)
func NewSortedDependencies_MissingProvideException_2(provide string, e javalang.IException) *SortedDependencies_MissingProvideException {
	this_ := &SortedDependencies_MissingProvideException{}
	InitSortedDependencies_MissingProvideException_2(this_, provide, e)
	return this_
}

func InitSortedDependencies_MissingProvideException_2(this_ *SortedDependencies_MissingProvideException, provide string, e javalang.IException) {
	_ = provide
	_ = e
	j2gort.InitUserThrowable(this_, "com.google.javascript.jscomp.deps.SortedDependencies.MissingProvideException", j2gort.BoxString(provide), e)
	this_.initFields_SortedDependencies_MissingProvideException()
}
