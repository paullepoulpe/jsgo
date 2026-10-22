// Port of com/google/javascript/jscomp/deps/SortedDependencies.java @ closure-compiler v20261015

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

package jscomp

import (
	"slices"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// SortedDependencies is a sorted list of inputs following the ES6 module
// ordering spec.
//
// Orders such that each input always comes after its dependencies. Circular
// references are allowed by emitting the current input in the moment before
// a loop would complete.
//
// The resulting order is not the same as the user-provided order but is
// influenced by it since it may take more than one graph traversal to
// account for all provided inputs and the graph traversals start with the
// first user provided input and continue from there.
//
// Also exposes other information about the inputs, like which inputs do not
// provide symbols.
//
// Inputs are compared with equals/hashCode, as in the Java collections.
type SortedDependencies[InputT DependencyInfo] struct {
	userOrderedInputs          *jutil.ArrayList[InputT]
	importOrderedInputs        *jutil.ArrayList[InputT]
	nonExportingInputs         *jutil.LinkedHashMap[string, InputT]
	exportingInputBySymbolName *jutil.LinkedHashMap[string, InputT]

	uniqueInputs            []InputT
	ptrToID                 map[any]int32
	objToID                 *jutil.ObjectLinkedHashMap[InputT, int32]
	exportingIDBySymbolName map[string]int32
	userOrderedIDs          []int32
	importOrderedIDs        []int32
	importOrderPos          []int32
	strongDeps              [][]int32
	weakDeps                [][]int32
	visited                 []uint32
	visitGen                uint32
	worklist                []int32
	weakWorklist            []int32
	hasDuplicateUserInputs  bool
}

// Java: SortedDependencies(List<InputT>)
func NewSortedDependencies[InputT DependencyInfo](userOrderedInputs []InputT) *SortedDependencies[InputT] {
	d := &SortedDependencies[InputT]{
		userOrderedInputs:          jutil.NewArrayList[InputT](),
		importOrderedInputs:        jutil.NewArrayList[InputT](),
		nonExportingInputs:         jutil.NewLinkedHashMap[string, InputT](),
		exportingInputBySymbolName: jutil.NewLinkedHashMap[string, InputT](),
	}
	d.userOrderedInputs.AddAll(userOrderedInputs...)
	d.ProcessInputs()
	return d
}

func (d *SortedDependencies[InputT]) lookupID(input InputT) (int32, bool) {
	if id, ok := d.ptrToID[input]; ok {
		return id, true
	}
	if d.objToID != nil {
		return d.objToID.Lookup(input)
	}
	return -1, false
}

func (d *SortedDependencies[InputT]) nextGen() uint32 {
	d.visitGen++
	if d.visitGen == 0 {
		clear(d.visited)
		d.visitGen = 1
	}
	return d.visitGen
}

func (d *SortedDependencies[InputT]) collectStrongDeps(rootInputs []InputT) uint32 {
	gen := d.nextGen()
	worklist := d.worklist[:0]
	for _, root := range rootInputs {
		if id, ok := d.lookupID(root); ok {
			if d.visited[id] != gen {
				d.visited[id] = gen
				worklist = append(worklist, id)
			}
		} else if !jutil.IsNil(root) {
			for _, req := range root.GetRequires() {
				if depID, ok := d.exportingIDBySymbolName[req.GetSymbol()]; ok && d.visited[depID] != gen {
					d.visited[depID] = gen
					worklist = append(worklist, depID)
				}
			}
		}
	}
	for i := 0; i < len(worklist); i++ {
		u := worklist[i]
		for _, depID := range d.strongDeps[u] {
			if d.visited[depID] != gen {
				d.visited[depID] = gen
				worklist = append(worklist, depID)
			}
		}
	}
	d.worklist = worklist
	return gen
}

// GetStrongDependenciesOf gets all the strong dependencies of the given
// roots. The inputs must be returned in a stable order. In other words, if A
// comes before B, and A does not transitively depend on B, then A must also
// come before B in the returned list.
//
// sorted: if true, get them in topologically sorted order. If false, get
// them in the original order they were passed to the compiler.
//
// Java: SortedDependencies.getStrongDependenciesOf(List<InputT>, boolean)
func (d *SortedDependencies[InputT]) GetStrongDependenciesOf(rootInputs []InputT, sorted bool) []InputT {
	if len(rootInputs) == 0 {
		return nil
	}
	gen := d.collectStrongDeps(rootInputs)
	worklist := d.worklist
	if len(worklist) == 0 {
		return nil
	}

	if sorted {
		if len(worklist) == 1 {
			return []InputT{d.uniqueInputs[worklist[0]]}
		}
		impSlice := d.importOrderedInputs.Slice()
		if len(worklist)*16 < len(d.importOrderedIDs) {
			for i, id := range worklist {
				worklist[i] = d.importOrderPos[id]
			}
			slices.Sort(worklist)
			builder := make([]InputT, len(worklist))
			for i, pos := range worklist {
				builder[i] = impSlice[pos]
			}
			return builder
		}
		builder := make([]InputT, 0, len(worklist))
		for i, id := range d.importOrderedIDs {
			if d.visited[id] == gen {
				builder = append(builder, impSlice[i])
			}
		}
		return builder
	}

	if !d.hasDuplicateUserInputs {
		if len(worklist) == 1 {
			return []InputT{d.uniqueInputs[worklist[0]]}
		}
		if len(worklist)*16 < len(d.uniqueInputs) {
			slices.Sort(worklist)
			builder := make([]InputT, len(worklist))
			for i, id := range worklist {
				builder[i] = d.uniqueInputs[id]
			}
			return builder
		}
		builder := make([]InputT, 0, len(worklist))
		for id, input := range d.uniqueInputs {
			if d.visited[id] == gen {
				builder = append(builder, input)
			}
		}
		return builder
	}

	userSlice := d.userOrderedInputs.Slice()
	builder := make([]InputT, 0, len(worklist))
	for i, id := range d.userOrderedIDs {
		if d.visited[id] == gen {
			builder = append(builder, userSlice[i])
		}
	}
	return builder
}

// GetInputProviding returns the input that gives us the given symbol. It
// panics with *SortedDependencies_MissingProvideException if there is no
// input for this symbol.
//
// Java: SortedDependencies.getInputProviding(String)
func (d *SortedDependencies[InputT]) GetInputProviding(symbolName string) InputT {
	input := d.MaybeGetInputProviding(symbolName)
	if !jutil.IsNil(input) {
		return input
	}
	panic(NewSortedDependencies_MissingProvideException(symbolName))
}

// Java: SortedDependencies.getInputsWithoutProvides()
func (d *SortedDependencies[InputT]) GetInputsWithoutProvides() []InputT {
	return slices.Collect(d.nonExportingInputs.Values())
}

// GetSortedStrongDependenciesOf gets all the strong dependencies of the
// given roots, in topologically sorted order.
//
// Java: SortedDependencies.getSortedStrongDependenciesOf(List<InputT>)
func (d *SortedDependencies[InputT]) GetSortedStrongDependenciesOf(roots []InputT) []InputT {
	return d.GetStrongDependenciesOf(roots, true)
}

// GetSortedWeakDependenciesOf gets all the weak dependencies of the given
// roots. The inputs must be returned in stable order. In other words, if A
// comes before B, and A does not transitively depend on B, then A must also
// come before B in the returned list.
//
// The weak dependencies are those that are only reachable via type requires
// from the roots. Note that if a root weakly requires another input, then
// all of its transitive dependencies (strong or weak) that are not strongly
// reachable from the roots will be included. e.g. if A weakly requires B,
// and B strongly requires C, and A is the sole root, then this will return
// B and C. However, if we add D as a root, and D strongly requires C, then
// this will only return B.
//
// Root inputs will never be in the returned list as they are all considered
// strong.
//
// Java: SortedDependencies.getSortedWeakDependenciesOf(List<InputT>)
func (d *SortedDependencies[InputT]) GetSortedWeakDependenciesOf(rootInputs []InputT) []InputT {
	if len(rootInputs) == 0 {
		return nil
	}
	strongGen := d.collectStrongDeps(rootInputs)
	strongList := d.worklist

	weakGen := d.nextGen()
	weakWorklist := d.weakWorklist[:0]
	for _, u := range strongList {
		for _, depID := range d.weakDeps[u] {
			if d.visited[depID] != strongGen && d.visited[depID] != weakGen {
				d.visited[depID] = weakGen
				weakWorklist = append(weakWorklist, depID)
			}
		}
	}
	for _, root := range rootInputs {
		if _, ok := d.lookupID(root); !ok && !jutil.IsNil(root) {
			for _, symbolName := range root.GetTypeRequires() {
				if depID, ok := d.exportingIDBySymbolName[symbolName]; ok && d.visited[depID] != strongGen && d.visited[depID] != weakGen {
					d.visited[depID] = weakGen
					weakWorklist = append(weakWorklist, depID)
				}
			}
		}
	}
	for i := 0; i < len(weakWorklist); i++ {
		u := weakWorklist[i]
		for _, depID := range d.strongDeps[u] {
			if d.visited[depID] != strongGen && d.visited[depID] != weakGen {
				d.visited[depID] = weakGen
				weakWorklist = append(weakWorklist, depID)
			}
		}
		for _, depID := range d.weakDeps[u] {
			if d.visited[depID] != strongGen && d.visited[depID] != weakGen {
				d.visited[depID] = weakGen
				weakWorklist = append(weakWorklist, depID)
			}
		}
	}
	d.weakWorklist = weakWorklist
	if len(weakWorklist) == 0 {
		return nil
	}
	if len(weakWorklist) == 1 {
		return []InputT{d.uniqueInputs[weakWorklist[0]]}
	}
	impSlice := d.importOrderedInputs.Slice()
	if len(weakWorklist)*16 < len(d.importOrderedIDs) {
		for i, id := range weakWorklist {
			weakWorklist[i] = d.importOrderPos[id]
		}
		slices.Sort(weakWorklist)
		builder := make([]InputT, len(weakWorklist))
		for i, pos := range weakWorklist {
			builder[i] = impSlice[pos]
		}
		return builder
	}

	builder := make([]InputT, 0, len(weakWorklist))
	for i, id := range d.importOrderedIDs {
		if d.visited[id] == weakGen {
			builder = append(builder, impSlice[i])
		}
	}
	return builder
}

// GetSortedList returns the inputs in import order (an unmodifiable view in
// Java; callers must not modify it).
//
// Java: SortedDependencies.getSortedList()
func (d *SortedDependencies[InputT]) GetSortedList() []InputT {
	return d.importOrderedInputs.Slice()
}

// MaybeGetInputProviding returns the input that gives us the given symbol,
// or null (the zero InputT).
//
// Java: SortedDependencies.maybeGetInputProviding(String)
func (d *SortedDependencies[InputT]) MaybeGetInputProviding(symbol string) InputT {
	if v, ok := d.exportingInputBySymbolName.Lookup(symbol); ok {
		return v
	}
	return d.nonExportingInputs.Get(ModuleNames_FileToModuleName(symbol))
}

func (d *SortedDependencies[InputT]) ProcessInputs() {
	userSlice := d.userOrderedInputs.Slice()
	n := len(userSlice)
	d.uniqueInputs = make([]InputT, 0, n)
	d.ptrToID = make(map[any]int32, n)
	var zero InputT
	if _, isEqualer := any(zero).(jutil.Equaler); isEqualer {
		d.objToID = jutil.NewObjectLinkedHashMap[InputT, int32]()
	}
	d.userOrderedIDs = make([]int32, n)
	d.exportingIDBySymbolName = make(map[string]int32, n)

	// Index.
	for i, userOrderedInput := range userSlice {
		id, ok := d.ptrToID[userOrderedInput]
		if !ok && d.objToID != nil {
			id, ok = d.objToID.Lookup(userOrderedInput)
			if ok {
				d.ptrToID[userOrderedInput] = id
			}
		}
		if !ok {
			id = int32(len(d.uniqueInputs))
			d.uniqueInputs = append(d.uniqueInputs, userOrderedInput)
			d.ptrToID[userOrderedInput] = id
			if d.objToID != nil {
				d.objToID.Put(userOrderedInput, id)
			}
		} else {
			d.hasDuplicateUserInputs = true
		}
		d.userOrderedIDs[i] = id

		provides := userOrderedInput.GetProvides()
		var firstProvide *string
		if len(provides) > 0 {
			firstProvide = &provides[0]
		}
		if firstProvide == nil ||
			// "module$" indicates the provide is generated from the path. If
			// this is the only thing the module provides and it is not an ES6
			// module then it is just a script and doesn't export anything.
			(len(provides) == 1 &&
				strings.HasPrefix(*firstProvide, "module$") &&
				// ES6 modules should always be considered as exporting
				// something.
				!userOrderedInput.IsEs6Module()) {
			d.nonExportingInputs.Put(ModuleNames_FileToModuleName(userOrderedInput.GetName()), userOrderedInput)
		}
		for _, providedSymbolName := range provides {
			d.exportingInputBySymbolName.Put(providedSymbolName, userOrderedInput)
			d.exportingIDBySymbolName[providedSymbolName] = id
		}
	}

	uCount := len(d.uniqueInputs)
	d.strongDeps = make([][]int32, uCount)
	d.weakDeps = make([][]int32, uCount)
	d.visited = make([]uint32, uCount)

	for i, userOrderedInput := range userSlice {
		u := d.userOrderedIDs[i]
		reqs := userOrderedInput.GetRequires()
		if len(reqs) > 0 {
			deps := d.strongDeps[u]
			for _, req := range reqs {
				if depID, ok := d.exportingIDBySymbolName[req.GetSymbol()]; ok && !slices.Contains(deps, depID) {
					deps = append(deps, depID)
				}
			}
			d.strongDeps[u] = deps
		}
		typeReqs := userOrderedInput.GetTypeRequires()
		if len(typeReqs) > 0 {
			wdeps := d.weakDeps[u]
			for _, symbolName := range typeReqs {
				if depID, ok := d.exportingIDBySymbolName[symbolName]; ok && !slices.Contains(wdeps, depID) {
					wdeps = append(wdeps, depID)
				}
			}
			d.weakDeps[u] = wdeps
		}
	}

	// Order.
	// For each input, traverse in user-provided order.
	completed := make([]bool, uCount)
	d.importOrderedIDs = make([]int32, 0, uCount)
	d.importOrderPos = make([]int32, uCount)
	var orderInput func(u int32)
	orderInput = func(u int32) {
		if completed[u] {
			return
		}
		completed[u] = true
		for _, depID := range d.strongDeps[u] {
			orderInput(depID)
		}
		// Emit an input after its imports have been emitted.
		d.importOrderPos[u] = int32(len(d.importOrderedIDs))
		d.importOrderedIDs = append(d.importOrderedIDs, u)
		d.importOrderedInputs.Add(d.uniqueInputs[u])
	}
	for _, u := range d.userOrderedIDs {
		orderInput(u)
	}
}

// SortedDependencies_MissingProvideException is
// SortedDependencies.MissingProvideException, a checked exception.
type SortedDependencies_MissingProvideException struct{ jutil.JavaException }

func (*SortedDependencies_MissingProvideException) missingProvideException() {}

// ISortedDependencies_MissingProvideException matches
// SortedDependencies_MissingProvideException.
type ISortedDependencies_MissingProvideException interface {
	jutil.IException
	missingProvideException()
}

// Java: SortedDependencies.MissingProvideException(String)
func NewSortedDependencies_MissingProvideException(provide string) *SortedDependencies_MissingProvideException {
	e := &SortedDependencies_MissingProvideException{}
	jutil.InitThrowable(e, "com.google.javascript.jscomp.deps.SortedDependencies$MissingProvideException", &provide, nil)
	return e
}

// Java: SortedDependencies.MissingProvideException(String, Exception)
func NewSortedDependencies_MissingProvideExceptionWithCause(provide string, e jutil.Throwable) *SortedDependencies_MissingProvideException {
	ex := &SortedDependencies_MissingProvideException{}
	jutil.InitThrowable(ex, "com.google.javascript.jscomp.deps.SortedDependencies$MissingProvideException", &provide, e)
	return ex
}
