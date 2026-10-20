// Port of com/google/javascript/jscomp/RemoveUnusedCode.java @ closure-compiler v20261015

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

package jscomp

import (
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/paullepoulpe/jsgo/internal/diagnostic"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

const (
	rucSlabMinChunk = 64
	rucSlabMaxChunk = 4096
)

// rucSlab is a reusable multi-chunk slab allocator that retains its backing
// chunk slices across passes and resets in O(touched chunks) time.
type rucSlab[T any] struct {
	chunks   [][]T
	chunkIdx int
	offset   int
}

func (s *rucSlab[T]) New() *T {
	if s.chunkIdx < len(s.chunks) {
		chunk := s.chunks[s.chunkIdx]
		if s.offset < len(chunk) {
			p := &chunk[s.offset]
			s.offset++
			return p
		}
		s.chunkIdx++
		s.offset = 0
		if s.chunkIdx < len(s.chunks) {
			chunk = s.chunks[s.chunkIdx]
			p := &chunk[0]
			s.offset = 1
			return p
		}
	}
	size := rucSlabMinChunk
	if n := len(s.chunks); n > 0 {
		prev := len(s.chunks[n-1])
		size = prev * 2
		if size > rucSlabMaxChunk {
			size = rucSlabMaxChunk
		}
	}
	chunk := make([]T, size)
	s.chunks = append(s.chunks, chunk)
	s.chunkIdx = len(s.chunks) - 1
	s.offset = 1
	return &chunk[0]
}

func (s *rucSlab[T]) reset() {
	if len(s.chunks) == 0 {
		return
	}
	limit := s.chunkIdx
	if limit >= len(s.chunks) {
		limit = len(s.chunks) - 1
	}
	const maxRetainedChunks = 40
	for i := 0; i < limit && i < maxRetainedChunks; i++ {
		clear(s.chunks[i])
	}
	if s.offset > 0 && limit < len(s.chunks) && limit < maxRetainedChunks {
		clear(s.chunks[limit][:s.offset])
	}
	if len(s.chunks) > maxRetainedChunks {
		clear(s.chunks[maxRetainedChunks:])
		s.chunks = s.chunks[:maxRetainedChunks]
	}
	s.chunkIdx = 0
	s.offset = 0
}

type propRemovableList struct {
	head *RemoveUnusedCode_RemovableBase
	tail *RemoveUnusedCode_RemovableBase
}

type rucVarCacheEntry struct {
	v    *Var
	info RemoveUnusedCode_VarInfo
}

func rucVarCacheIdx(v *Var) uint32 {
	return (uint32(uintptr(unsafe.Pointer(v))>>4) * 0x9e3779b1) >> 21
}

type removeUnusedCodeArena struct {
	worklist     []RemoveUnusedCode_Continuation
	worklistHead int
	worklistSize int

	varBuckets  []*RemoveUnusedCode_RealVarInfo
	varInfos    []*RemoveUnusedCode_RealVarInfo
	lastVar     *Var
	lastVarInfo RemoveUnusedCode_VarInfo
	varPtrCache [2048]rucVarCacheEntry

	pinnedPropertyNames map[string]struct{}
	lastPinnedProp      string
	propRemovables      map[string]propRemovableList
	propRemovableKeys   []string

	allFunctionParamScopes []*Scope

	realVarInfoSlab      rucSlab[RemoveUnusedCode_RealVarInfo]
	assignSlab           rucSlab[RemoveUnusedCode_Assign]
	nameDeclStmtSlab     rucSlab[RemoveUnusedCode_NameDeclarationStatement]
	indirectAssignSlab   rucSlab[RemoveUnusedCode_IndirectAssign]
	anonProtoAssignSlab  rucSlab[RemoveUnusedCode_AnonymousPrototypeNamedPropertyAssign]
	instanceofNameSlab   rucSlab[RemoveUnusedCode_InstanceofName]
	funcDeclSlab         rucSlab[RemoveUnusedCode_FunctionDeclaration]
	classSetupCallSlab   rucSlab[RemoveUnusedCode_ClassSetupCall]
	unusedReadRefSlab    rucSlab[RemoveUnusedCode_UnusedReadReference]
	classOrProtoPropSlab rucSlab[RemoveUnusedCode_ClassOrPrototypeNamedProperty]
	classDeclSlab        rucSlab[RemoveUnusedCode_ClassDeclaration]
	vanillaForDeclSlab   rucSlab[RemoveUnusedCode_VanillaForNameDeclaration]
}

func (a *removeUnusedCodeArena) reset() {
	if a.worklistSize > 0 {
		clear(a.worklist)
	}
	a.worklistHead = 0
	a.worklistSize = 0

	a.lastVar = nil
	a.lastVarInfo = nil
	clear(a.varPtrCache[:])
	if len(a.varBuckets) > 0 {
		if len(a.varInfos)*8 < len(a.varBuckets) {
			mask := uint32(len(a.varBuckets) - 1)
			for _, info := range a.varInfos {
				a.varBuckets[info.hash&mask] = nil
			}
		} else {
			clear(a.varBuckets)
		}
	}
	clear(a.varInfos)
	a.varInfos = a.varInfos[:0]

	clear(a.pinnedPropertyNames)
	a.lastPinnedProp = ""
	clear(a.propRemovables)
	clear(a.propRemovableKeys)
	a.propRemovableKeys = a.propRemovableKeys[:0]

	clear(a.allFunctionParamScopes)
	a.allFunctionParamScopes = a.allFunctionParamScopes[:0]

	a.realVarInfoSlab.reset()
	a.assignSlab.reset()
	a.nameDeclStmtSlab.reset()
	a.indirectAssignSlab.reset()
	a.anonProtoAssignSlab.reset()
	a.instanceofNameSlab.reset()
	a.funcDeclSlab.reset()
	a.classSetupCallSlab.reset()
	a.unusedReadRefSlab.reset()
	a.classOrProtoPropSlab.reset()
	a.classDeclSlab.reset()
	a.vanillaForDeclSlab.reset()
}

var cachedRUCArena atomic.Pointer[removeUnusedCodeArena]

var removeUnusedCodeArenaPool = sync.Pool{
	New: func() any {
		return &removeUnusedCodeArena{
			worklist:               make([]RemoveUnusedCode_Continuation, 512),
			varBuckets:             make([]*RemoveUnusedCode_RealVarInfo, 16384),
			varInfos:               make([]*RemoveUnusedCode_RealVarInfo, 0, 2048),
			pinnedPropertyNames:    make(map[string]struct{}, 256),
			propRemovables:         make(map[string]propRemovableList, 256),
			propRemovableKeys:      make([]string, 0, 256),
			allFunctionParamScopes: make([]*Scope, 0, 512),
		}
	},
}

func acquireRUCArena() *removeUnusedCodeArena {
	if a := cachedRUCArena.Swap(nil); a != nil {
		return a
	}
	return removeUnusedCodeArenaPool.Get().(*removeUnusedCodeArena)
}

func releaseRUCArena(a *removeUnusedCodeArena) {
	if cachedRUCArena.CompareAndSwap(nil, a) {
		return
	}
	removeUnusedCodeArenaPool.Put(a)
}

// RemoveUnusedCode is garbage collection for variable and function
// definitions. Basically performs a mark-and-sweep type algorithm over the
// JavaScript parse tree.
//
// For each scope: (1) Scan the variable/function declarations at that scope.
// (2) Traverse the scope for references, marking all referenced variables.
// Unlike other compiler passes, this is a pre-order traversal, not a
// post-order traversal. (3) If the traversal encounters an assign without
// other side-effects, create a continuation. Continue the continuation iff
// the assigned variable is referenced. (4) When the traversal completes,
// remove all unreferenced variables.
//
// If it makes it easier, you can think of the continuations of the traversal
// as a reference graph. Each continuation represents a set of edges, where
// the source node is a known variable, and the destination nodes are lazily
// evaluated when the continuation is executed.
//
// This algorithm is similar to the algorithm used by SmartNameRemoval.
// SmartNameRemoval maintains an explicit graph of dependencies between
// global symbols. However, SmartNameRemoval cannot handle non-trivial edges
// in the reference graph ("A is referenced iff both B and C are
// referenced"), or local variables. SmartNameRemoval is also substantially
// more complicated because it tries to handle namespaces (which is largely
// unnecessary in the presence of CollapseProperties.
//
// This pass also uses a more complex analysis of assignments, where an
// assignment to a variable or a property of that variable does not
// necessarily count as a reference to that variable, unless we can prove
// that it modifies external state. This is similar to
// FlowSensitiveInlineVariables, except that it works for variables used
// across scopes.
//
// Multiple datastructures are used to accumulate nodes, some of which are
// later removed. Since some nodes encompass a subtree of nodes, the removal
// can sometimes pre-remove other nodes which are also referenced in these
// datastructures for later removal. Attempting double-removal violates
// scope change notification constraints so there is a desire to excise
// already-removed subtree nodes from these datastructures. But not all of
// the datastructures are conducive to flexible removal and the ones that are
// conducive don't necessarily track all flavors of nodes. So instead of
// updating datastructures on the fly a pre-check is performed to skip
// already-removed nodes right before the moment an attempt to remove them
// would otherwise be made.
type RemoveUnusedCode struct {
	compiler    AbstractCompiler
	astAnalyzer *AstAnalyzer

	codingConvention CodingConvention

	removeLocalVars bool
	removeGlobals   bool

	preserveFunctionExpressionNames bool

	externProperties *jutil.LinkedHashSet[string]
	accessorSummary  *AccessorSummary
	// Non-nil only when the exported implicit-property list differs from its default.
	customImplicitProperties map[string]struct{}

	// Single value to use for all vars for which we cannot remove anything at
	// all.
	canonicalUnremovableVarInfo    RemoveUnusedCode_VarInfo
	canonicalUnremovableVarInfoVal RemoveUnusedCode_CanonicalUnremovableVarInfo

	// Stores the names of all "leaf" properties that are polyfilled, to avoid
	// unnecessary qualified name matching and searches for all the other
	// properties. This includes global names such as "Promise" and "Map",
	// static methods on global names such as "Array.from" and "Math.fround",
	// and instance properties such as "String.prototype.repeat" and
	// "Promise.prototype.finally".
	polyfills *jutil.HashMultimap[string, RemoveUnusedCode_PolyfillInfo]

	guardedUsages         *jutil.LinkedHashSet[*rhino.Node]
	guardedUsagesComputed bool
	root                  *rhino.Node

	polyfillsFromTable *PolyfillUsageFinder_Polyfills

	scopeCreator *SyntacticScopeCreator

	removeUnusedPrototypeProperties               bool
	removeUnusedThisProperties                    bool
	removeUnusedObjectDefinePropertiesDefinitions bool
	removeUnusedPolyfills                         bool
	assumeGettersArePure                          bool

	// Allocated & cleaned up by process()
	removalLog     diagnostic.LogFile
	unremovableLog diagnostic.LogFile

	// Per-pass pooled arena holding worklist, varBuckets, varInfos, property maps, and slabs.
	removeUnusedCodeArena
}

// Properties that are implicitly used as part of the JS language.
var RemoveUnusedCode_IMPLICITLY_USED_PROPERTIES = []string{"length", "toString", "valueOf", "constructor", "prototype"}

func isImplicitlyUsedProperty(name string) bool {
	switch name {
	case "length", "toString", "valueOf", "constructor", "prototype":
		return true
	default:
		return false
	}
}

// Java: RemoveUnusedCode(Builder)
func NewRemoveUnusedCode(builder *RemoveUnusedCode_Builder) *RemoveUnusedCode {
	sc := NewSyntacticScopeCreator(builder.compiler)
	sc.EnableArenaPool()
	r := &RemoveUnusedCode{
		compiler:                        builder.compiler,
		astAnalyzer:                     builder.compiler.GetAstAnalyzer(),
		codingConvention:                builder.compiler.GetCodingConvention(),
		scopeCreator:                    sc,
		removeLocalVars:                 builder.removeLocalVars,
		removeGlobals:                   builder.removeGlobals,
		preserveFunctionExpressionNames: builder.preserveFunctionExpressionNames,
		removeUnusedPrototypeProperties: builder.removeUnusedPrototypeProperties,
		removeUnusedThisProperties:      builder.removeUnusedThisProperties,
		removeUnusedObjectDefinePropertiesDefinitions: builder.removeUnusedObjectDefinePropertiesDefinitions,
		removeUnusedPolyfills:                         builder.removeUnusedPolyfills,
		polyfillsFromTable:                            defaultPolyfillsTable(),
		assumeGettersArePure:                          builder.assumeGettersArePure,
	}

	// Preserve the constructor's snapshot of the exported mutable list. Keep the
	// allocation-free switch for the stock list, including later global changes.
	defaults := [...]string{"length", "toString", "valueOf", "constructor", "prototype"}
	stock := len(RemoveUnusedCode_IMPLICITLY_USED_PROPERTIES) == len(defaults)
	if stock {
		for i, name := range defaults {
			if RemoveUnusedCode_IMPLICITLY_USED_PROPERTIES[i] != name {
				stock = false
				break
			}
		}
	}
	if !stock {
		r.customImplicitProperties = make(map[string]struct{}, len(RemoveUnusedCode_IMPLICITLY_USED_PROPERTIES))
		for _, name := range RemoveUnusedCode_IMPLICITLY_USED_PROPERTIES {
			r.customImplicitProperties[name] = struct{}{}
		}
	}

	// All Vars that are completely unremovable will share this VarInfo instance.
	r.canonicalUnremovableVarInfoVal.outer = r
	r.canonicalUnremovableVarInfo = &r.canonicalUnremovableVarInfoVal
	return r
}

// RemoveUnusedCode_Builder is the nested class Builder.
type RemoveUnusedCode_Builder struct {
	compiler AbstractCompiler

	removeLocalVars                               bool
	removeGlobals                                 bool
	preserveFunctionExpressionNames               bool
	removeUnusedPrototypeProperties               bool
	removeUnusedThisProperties                    bool
	removeUnusedObjectDefinePropertiesDefinitions bool
	removeUnusedPolyfills                         bool
	assumeGettersArePure                          bool
}

// Java: RemoveUnusedCode.Builder(AbstractCompiler)
func NewRemoveUnusedCode_Builder(compiler AbstractCompiler) *RemoveUnusedCode_Builder {
	return &RemoveUnusedCode_Builder{compiler: compiler}
}

// Java: RemoveUnusedCode.Builder.removeLocalVars(boolean)
func (b *RemoveUnusedCode_Builder) RemoveLocalVars(value bool) *RemoveUnusedCode_Builder {
	b.removeLocalVars = value
	return b
}

// Java: RemoveUnusedCode.Builder.removeGlobals(boolean)
func (b *RemoveUnusedCode_Builder) RemoveGlobals(value bool) *RemoveUnusedCode_Builder {
	b.removeGlobals = value
	return b
}

// Java: RemoveUnusedCode.Builder.preserveFunctionExpressionNames(boolean)
func (b *RemoveUnusedCode_Builder) PreserveFunctionExpressionNames(value bool) *RemoveUnusedCode_Builder {
	b.preserveFunctionExpressionNames = value
	return b
}

// Java: RemoveUnusedCode.Builder.removeUnusedPrototypeProperties(boolean)
func (b *RemoveUnusedCode_Builder) RemoveUnusedPrototypeProperties(value bool) *RemoveUnusedCode_Builder {
	b.removeUnusedPrototypeProperties = value
	return b
}

// Java: RemoveUnusedCode.Builder.removeUnusedThisProperties(boolean)
func (b *RemoveUnusedCode_Builder) RemoveUnusedThisProperties(value bool) *RemoveUnusedCode_Builder {
	b.removeUnusedThisProperties = value
	return b
}

// Java: RemoveUnusedCode.Builder.removeUnusedObjectDefinePropertiesDefinitions(boolean)
func (b *RemoveUnusedCode_Builder) RemoveUnusedObjectDefinePropertiesDefinitions(value bool) *RemoveUnusedCode_Builder {
	b.removeUnusedObjectDefinePropertiesDefinitions = value
	return b
}

// Java: RemoveUnusedCode.Builder.removeUnusedPolyfills(boolean)
func (b *RemoveUnusedCode_Builder) RemoveUnusedPolyfills(value bool) *RemoveUnusedCode_Builder {
	b.removeUnusedPolyfills = value
	return b
}

// Java: RemoveUnusedCode.Builder.assumeGettersArePure(boolean)
func (b *RemoveUnusedCode_Builder) AssumeGettersArePure(value bool) *RemoveUnusedCode_Builder {
	b.assumeGettersArePure = value
	return b
}

// Java: RemoveUnusedCode.Builder.build()
func (b *RemoveUnusedCode_Builder) Build() *RemoveUnusedCode {
	return NewRemoveUnusedCode(b)
}

// RemoveUnusedCode_RemovalLogRecord supplies the string needed for an entry
// in the removal log.
type RemoveUnusedCode_RemovalLogRecord struct {
	kind                 string
	nameSupplier         func() string
	functionNameSupplier func() string
}

// Get returns a log entry string.
//
// Each entry is one tab-separated line of the form:
//
//	KIND NAME [FUNCTION_NAME]
//
// See specific methods below for details.
//
// Java: RemoveUnusedCode.RemovalLogRecord.get()
func (r *RemoveUnusedCode_RemovalLogRecord) Get() string {
	return strings.Join([]string{r.kind, r.nameSupplier(), r.functionNameSupplier()}, "\t")
}

// Java: RemoveUnusedCode.RemovalLogRecord(String, Supplier<String>, Supplier<String>)
func NewRemoveUnusedCode_RemovalLogRecordWithFunctionName(
	kind string, nameSupplier func() string, functionNameSupplier func() string) *RemoveUnusedCode_RemovalLogRecord {
	return &RemoveUnusedCode_RemovalLogRecord{
		kind:                 kind, // checkNotNull: a Go string is never null.
		nameSupplier:         jutil.CheckNotNull(nameSupplier),
		functionNameSupplier: jutil.CheckNotNull(functionNameSupplier),
	}
}

// Java: RemoveUnusedCode.RemovalLogRecord(String, Supplier<String>)
func NewRemoveUnusedCode_RemovalLogRecord(kind string, nameSupplier func() string) *RemoveUnusedCode_RemovalLogRecord {
	// No function name
	return NewRemoveUnusedCode_RemovalLogRecordWithFunctionName(kind, nameSupplier, func() string { return "" })
}

// Java: RemoveUnusedCode.RemovalLogRecord.forProperty(String)
func RemoveUnusedCode_RemovalLogRecord_ForProperty(propName string) *RemoveUnusedCode_RemovalLogRecord {
	return NewRemoveUnusedCode_RemovalLogRecord("prop", func() string { return propName })
}

// Java: RemoveUnusedCode.RemovalLogRecord.forVar(Var)
func RemoveUnusedCode_RemovalLogRecord_ForVar(var_ *Var) *RemoveUnusedCode_RemovalLogRecord {
	return NewRemoveUnusedCode_RemovalLogRecord("var", var_.GetName)
}

// Java: RemoveUnusedCode.RemovalLogRecord.forPolyfill(PolyfillInfo)
func RemoveUnusedCode_RemovalLogRecord_ForPolyfill(polyfillInfo RemoveUnusedCode_PolyfillInfo) *RemoveUnusedCode_RemovalLogRecord {
	return NewRemoveUnusedCode_RemovalLogRecord("poly", polyfillInfo.GetName)
}

// RemoveUnusedCode_RemovalLogRecord_ForNamedArg records removal of a named
// function parameter.
//
// nameNode is the parameter's NAME node; argList the function's PARAM_LIST
// node.
//
// Java: RemoveUnusedCode.RemovalLogRecord.forNamedArg(Node, Node)
func RemoveUnusedCode_RemovalLogRecord_ForNamedArg(nameNode *rhino.Node, argList *rhino.Node) *RemoveUnusedCode_RemovalLogRecord {
	return NewRemoveUnusedCode_RemovalLogRecordWithFunctionName(
		"arg", nameNode.GetString, RemoveUnusedCode_RemovalLogRecord_GetLoggableFunctionNameSupplier(argList))
}

// RemoveUnusedCode_RemovalLogRecord_ForDestructuringArg records removal of a
// destructuring function parameter.
//
// argList is the function's PARAM_LIST node.
//
// Java: RemoveUnusedCode.RemovalLogRecord.forDestructuringArg(Node)
func RemoveUnusedCode_RemovalLogRecord_ForDestructuringArg(argList *rhino.Node) *RemoveUnusedCode_RemovalLogRecord {
	return NewRemoveUnusedCode_RemovalLogRecordWithFunctionName(
		"arg", func() string { return "<pattern>" }, RemoveUnusedCode_RemovalLogRecord_GetLoggableFunctionNameSupplier(argList))
}

// RemoveUnusedCode_RemovalLogRecord_ForMarkingNamedArg records that a named
// parameter is marked as unused for possible removal by OptimizeParameters.
//
// nameNode is the parameter's NAME node; argList the function's PARAM_LIST
// node.
//
// Java: RemoveUnusedCode.RemovalLogRecord.forMarkingNamedArg(Node, Node)
func RemoveUnusedCode_RemovalLogRecord_ForMarkingNamedArg(nameNode *rhino.Node, argList *rhino.Node) *RemoveUnusedCode_RemovalLogRecord {
	return NewRemoveUnusedCode_RemovalLogRecordWithFunctionName(
		"argmark", nameNode.GetString, RemoveUnusedCode_RemovalLogRecord_GetLoggableFunctionNameSupplier(argList))
}

// RemoveUnusedCode_RemovalLogRecord_GetLoggableFunctionNameSupplier returns
// a supplier for the FUNCTION_NAME field of an argument removal log entry.
//
// If no good name can be found, then "<anonymous>" will be supplied.
//
// argList is the function's PARAM_LIST node.
//
// Java: RemoveUnusedCode.RemovalLogRecord.getLoggableFunctionNameSupplier(Node)
func RemoveUnusedCode_RemovalLogRecord_GetLoggableFunctionNameSupplier(argList *rhino.Node) func() string {
	return func() string {
		functionName := NodeUtil_GetNearestFunctionName(jutil.CheckNotNull(argList).GetParent())
		if functionName == nil {
			return "<anonymous>"
		}
		return *functionName
	}
}

// Process traverses the root, removing all unused variables. Multiple
// traversals may occur to ensure all unused variables are removed.
//
// Java: RemoveUnusedCode.process(Node, Node)
func (r *RemoveUnusedCode) Process(externs *rhino.Node, root *rhino.Node) {
	jutil.CheckState(r.compiler.GetLifeCycleStage().IsNormalized())
	r.externProperties = r.compiler.GetExternProperties()
	r.accessorSummary = r.compiler.GetAccessorSummary()

	arena := acquireRUCArena()
	r.removeUnusedCodeArena = *arena
	defer func() {
		if r.scopeCreator != nil {
			r.scopeCreator.ReleaseArena()
		}
		*arena = r.removeUnusedCodeArena
		arena.reset()
		r.removeUnusedCodeArena = removeUnusedCodeArena{}
		releaseRUCArena(arena)
	}()

	removalLogFile :=
		r.compiler.CreateOrReopenIndexedLog(reflect.TypeFor[RemoveUnusedCode](), "removals.log")
	defer removalLogFile.Close()
	keepLogFile :=
		r.compiler.CreateOrReopenIndexedLog(reflect.TypeFor[RemoveUnusedCode](), "unremovable.log")
	defer keepLogFile.Close()
	defer func() {
		r.removalLog = nil
		r.unremovableLog = nil
	}()
	r.removalLog = removalLogFile // avoid passing the log file through a bunch of methods
	r.unremovableLog = keepLogFile
	r.TraverseAndRemoveUnusedReferences(root)
}

func (r *RemoveUnusedCode) enqueueContinuation(c RemoveUnusedCode_Continuation) {
	if r.worklistSize == len(r.worklist) {
		r.growWorklist()
	}
	idx := r.worklistHead + r.worklistSize
	if idx >= len(r.worklist) {
		idx -= len(r.worklist)
	}
	r.worklist[idx] = c
	r.worklistSize++
}

func (r *RemoveUnusedCode) growWorklist() {
	n := len(r.worklist) * 2
	if n < 256 {
		n = 256
	}
	nb := make([]RemoveUnusedCode_Continuation, n)
	if r.worklistSize > 0 {
		firstPart := len(r.worklist) - r.worklistHead
		if firstPart >= r.worklistSize {
			copy(nb, r.worklist[r.worklistHead:r.worklistHead+r.worklistSize])
		} else {
			copy(nb, r.worklist[r.worklistHead:])
			copy(nb[firstPart:], r.worklist[:r.worklistSize-firstPart])
		}
	}
	r.worklist = nb
	r.worklistHead = 0
}

func (r *RemoveUnusedCode) popContinuation() RemoveUnusedCode_Continuation {
	c := r.worklist[r.worklistHead]
	r.worklist[r.worklistHead] = RemoveUnusedCode_Continuation{}
	r.worklistHead++
	if r.worklistHead == len(r.worklist) {
		r.worklistHead = 0
	}
	r.worklistSize--
	return c
}

// TraverseAndRemoveUnusedReferences traverses a node recursively. Call this
// once per pass.
//
// Java: RemoveUnusedCode.traverseAndRemoveUnusedReferences(Node)
func (r *RemoveUnusedCode) TraverseAndRemoveUnusedReferences(root *rhino.Node) {
	if r.pinnedPropertyNames == nil {
		r.pinnedPropertyNames = make(map[string]struct{}, 64)
		r.propRemovables = make(map[string]propRemovableList, 64)
	}
	// Create scope from parent of root node, which also has externs as a child, so we'll
	// have extern definitions in scope.
	scope := r.scopeCreator.CreateScopeOfScope(root.GetParent(), nil)
	if !scope.HasSlot(NodeUtil_JSC_PROPERTY_NAME_FN) {
		// TODO(b/70730762): Passes that add references to this should ensure it is declared.
		// NOTE: null input makes this an extern var.
		scope.Declare(
			NodeUtil_JSC_PROPERTY_NAME_FN /* no declaration node */, nil /* no input */, nil)
	}

	r.root = root
	r.guardedUsagesComputed = false

	r.enqueueContinuation(NewRemoveUnusedCode_Continuation(r, root, scope))
	for r.worklistSize > 0 {
		continuation := r.popContinuation() // worklist.remove()
		continuation.Apply(r)
	}

	r.RemoveUnreferencedVarsAndPolyfills()
	r.RemoveIndependentlyRemovableProperties()
	for _, fparamScope := range r.allFunctionParamScopes {
		r.RemoveUnreferencedFunctionArgs(fparamScope)
	}
}

// Java: RemoveUnusedCode.storePolyfill(PolyfillUsage)
func (r *RemoveUnusedCode) StorePolyfill(polyfillUsage *PolyfillUsageFinder_PolyfillUsage) {
	if r.guardedUsages == nil {
		r.guardedUsages = jutil.NewLinkedHashSet[*rhino.Node]()
	}
	r.guardedUsages.Add(polyfillUsage.Node())
}

// Java: RemoveUnusedCode.removeIndependentlyRemovableProperties()
func (r *RemoveUnusedCode) RemoveIndependentlyRemovableProperties() {
	for _, propName := range r.propRemovableKeys {
		list, ok := r.propRemovables[propName]
		if !ok || list.head == nil {
			continue
		}
		if r.removalLog.IsLogging() {
			r.removalLog.LogSupplier(RemoveUnusedCode_RemovalLogRecord_ForProperty(propName).Get)
		}
		for cur := list.head; cur != nil; {
			next := cur.propNext
			cur.propNext = nil
			cur.inPropList = false
			cur.self.Remove(r.compiler)
			cur = next
		}
	}
}

func hasBlockScopedDeclarations(n *rhino.Node) bool {
	for child := n.GetFirstChild(); child != nil; child = child.GetNext() {
		if hasBlockScopedDeclarationsChild(child) {
			return true
		}
	}
	return false
}

func hasBlockScopedDeclarationsChild(n *rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_LET, rhino.Token_CONST, rhino.Token_CATCH, rhino.Token_IMPORT:
		return true
	case rhino.Token_EXPORT:
		first := n.GetFirstChild()
		return first != nil && hasBlockScopedDeclarationsChild(first)
	case rhino.Token_FUNCTION:
		return !NodeUtil_IsFunctionExpression(n) && n.GetFirstChild().GetString() != ""
	case rhino.Token_CLASS:
		return !NodeUtil_IsClassExpression(n) && n.GetFirstChild().GetString() != ""
	case rhino.Token_BLOCK:
		if scanVarsBlockCreatesScope(n) {
			return false
		}
		return hasBlockScopedDeclarations(n)
	case rhino.Token_ROOT, rhino.Token_DO, rhino.Token_WHILE, rhino.Token_WITH, rhino.Token_IF,
		rhino.Token_LABEL, rhino.Token_TRY, rhino.Token_SWITCH, rhino.Token_CASE, rhino.Token_DEFAULT_CASE:
		return hasBlockScopedDeclarations(n)
	default:
		return false
	}
}

// TraverseNode traverses everything in the current scope and marks
// variables that are referenced.
//
// During traversal, we identify subtrees that will only be referenced if
// their enclosing variables are referenced. Instead of traversing those
// subtrees, we create a continuation for them, and traverse them lazily.
//
// Java: RemoveUnusedCode.traverseNode(Node, Scope)
func (r *RemoveUnusedCode) TraverseNode(n *rhino.Node, scope *Scope) {
	type_ := n.GetToken()
	switch type_ {
	case rhino.Token_CATCH:
		r.TraverseCatch(n, scope)
	case rhino.Token_FUNCTION:
		{
			var varInfo RemoveUnusedCode_VarInfo
			// If this function is a removable var, then create a continuation
			// for it instead of traversing immediately.
			if NodeUtil_IsFunctionDeclaration(n) {
				varInfo = r.TraverseNameNode(n.GetFirstChild(), scope)
				functionDeclaration :=
					NewRemoveUnusedCode_RemovableBuilder(r).
						AddContinuation(NewRemoveUnusedCode_Continuation(r, n, scope)).
						BuildFunctionDeclaration(n)
				varInfo.AddRemovable(functionDeclaration)
				if n.GetParent().IsExport() {
					varInfo.SetIsExplicitlyNotRemovable(func() string { return "exported class" })
				}
			} else {
				r.TraverseFunction(n, scope)
			}
		}
	case rhino.Token_ASSIGN:
		r.TraverseAssign(n, scope)
	case rhino.Token_ASSIGN_BITOR,
		rhino.Token_ASSIGN_BITXOR,
		rhino.Token_ASSIGN_BITAND,
		rhino.Token_ASSIGN_LSH,
		rhino.Token_ASSIGN_RSH,
		rhino.Token_ASSIGN_URSH,
		rhino.Token_ASSIGN_ADD,
		rhino.Token_ASSIGN_SUB,
		rhino.Token_ASSIGN_MUL,
		rhino.Token_ASSIGN_EXPONENT,
		rhino.Token_ASSIGN_DIV,
		rhino.Token_ASSIGN_MOD:
		r.TraverseCompoundAssign(n, scope)
	case rhino.Token_INC, rhino.Token_DEC:
		r.TraverseIncrementOrDecrementOp(n, scope)
	case rhino.Token_CALL, rhino.Token_OPTCHAIN_CALL:
		r.TraverseCall(n, scope)
	case rhino.Token_SWITCH_BODY, rhino.Token_BLOCK:
		// This case if for if there are let and const variables in block scopes.
		// Otherwise other variables will be hoisted up into the global scope and already be
		// handled.
		if NodeUtil_CreatesBlockScope(n) &&
			(NodeUtil_IsFunctionBlock(n) || NodeUtil_IsClassStaticBlock(n) || hasBlockScopedDeclarations(n)) {
			r.TraverseChildren(n, r.scopeCreator.CreateScopeOfScope(n, scope))
		} else {
			r.TraverseChildren(n, scope)
		}
	case rhino.Token_MODULE_BODY:
		r.TraverseChildren(n, r.scopeCreator.CreateScopeOfScope(n, scope))
	case rhino.Token_CLASS:
		r.TraverseClass(n, scope)
	case rhino.Token_CLASS_MEMBERS:
		r.TraverseClassMembers(n, scope)
	case rhino.Token_ARRAY_PATTERN, rhino.Token_PARAM_LIST:
		r.TraverseIndirectAssignmentList(n, scope)
	case rhino.Token_OBJECT_PATTERN:
		r.TraverseObjectPattern(n, scope)
	case rhino.Token_OBJECTLIT:
		r.TraverseObjectLiteral(n, scope)
	case rhino.Token_FOR:
		r.TraverseVanillaFor(n, scope)
	case rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF:
		r.TraverseEnhancedFor(n, scope)
	case rhino.Token_LET, rhino.Token_CONST, rhino.Token_VAR:
		// for-loop cases are handled by custom traversal methods.
		jutil.CheckState(NodeUtil_IsStatement(n))
		r.TraverseDeclarationStatement(n, scope)
	case rhino.Token_INSTANCEOF:
		r.TraverseInstanceof(n, scope)
	case rhino.Token_NAME:
		// The only cases that should reach this point are parameter declarations and references
		// to names. The name node does not have children in these cases.
		jutil.CheckState(!n.HasChildren())
		parent := n.GetParent()
		// the parameter declaration is not a read of the name
		if !parent.IsParamList() {
			// var|let|const name;
			// are handled at a higher level.
			jutil.CheckState(!NodeUtil_IsNameDeclaration(parent))
			// function name() {}
			// class name() {}
			// handled at a higher level
			jutil.CheckState(!((parent.IsFunction() || parent.IsClass()) && parent.GetFirstChild() == n))
			varInfo := r.TraverseNameNode(n, scope)
			if varInfo.IsRemovable() {
				var reason func() string
				if r.unremovableLog.IsLogging() {
					reason = func() string { return "reference found: " + n.GetLocation() }
				}
				varInfo.SetIsExplicitlyNotRemovable(reason)
			}
		}
	case rhino.Token_GETPROP, rhino.Token_OPTCHAIN_GETPROP:
		r.TraverseNormalOrOptChainGetProp(n, scope)
	default:
		r.TraverseChildren(n, scope)
	}
}

// Java: RemoveUnusedCode.traverseInstanceof(Node, Scope)
func (r *RemoveUnusedCode) TraverseInstanceof(instanceofNode *rhino.Node, scope *Scope) {
	jutil.CheckArgumentMsg1(instanceofNode.IsInstanceOf(), "%s", instanceofNode)
	lhs := instanceofNode.GetFirstChild()
	rhs := lhs.GetNext()
	r.TraverseNode(lhs, scope)
	if rhs.IsName() {
		varInfo := r.TraverseNameNode(rhs, scope)
		builder := NewRemoveUnusedCode_RemovableBuilder(r)
		varInfo.AddRemovable(builder.BuildInstanceofName(instanceofNode))
	} else {
		r.TraverseNode(rhs, scope)
	}
}

// TraverseNormalOrOptChainGetProp traverses `expr.prop` or `expr?.prop`.
//
// Note that this method is called only for RHS nodes. Property references
// that are being assigned to are handled by the logic traversing their
// parent (e.g. ASSIGN) node.
//
// The primary purpose of this method is to make sure the property reference
// is correctly recorded.
//
// Java: RemoveUnusedCode.traverseNormalOrOptChainGetProp(Node, Scope)
func (r *RemoveUnusedCode) TraverseNormalOrOptChainGetProp(getProp *rhino.Node, scope *Scope) {
	jutil.CheckStateMsg1(NodeUtil_IsNormalOrOptChainGetProp(getProp), "%s", getProp)
	objectNode := getProp.GetFirstChild()
	propertyName := getProp.GetString()

	if r.polyfills != nil && !r.polyfills.IsEmpty() && r.polyfills.ContainsKey(propertyName) {
		for info := range r.polyfills.Get(propertyName).All() {
			if info.GetRemoveUnusedCode_PolyfillInfoBase().isRemovable {
				info.ConsiderPossibleReference(getProp)
			}
		}
	}

	if NodeUtil_IsExpressionResultUsed(getProp) ||
		r.ConsiderForAccessorSideEffects(getProp, AccessorSummary_PropertyAccessKind_GETTER_ONLY) {
		// must record as reference to the property and continue traversal.
		r.MarkPropertyNameAsPinned(propertyName)
		r.TraverseNode(objectNode, scope)
	} else if objectNode.IsThis() {
		// This is probably the declaration of a class field in a constructor.
		// /** @private {number} */
		// this.propName;
		// We don't want to consider this a real usage that should prevent removal.
		builder := NewRemoveUnusedCode_RemovableBuilder(r).SetIsThisDotPropertyReference(true)
		r.ConsiderForIndependentRemoval(builder.BuildUnusedReadReference(getProp, getProp))
	} else if RemoveUnusedCode_IsDotPrototype(objectNode) {
		// (objExpression).prototype.propName;
		builder := NewRemoveUnusedCode_RemovableBuilder(r).SetIsPrototypeDotPropertyReference(true)
		objExpression := objectNode.GetFirstChild()
		if objExpression.IsName() {
			// name.prototype.propName;
			varInfo := r.TraverseNameNode(objExpression, scope)
			varInfo.AddRemovable(builder.BuildUnusedReadReference(getProp, getProp))
		} else {
			// (objExpression).prototype.propName;
			if r.MayHaveSideEffects(objExpression) {
				r.TraverseNode(objExpression, scope)
			} else {
				builder.AddContinuation(NewRemoveUnusedCode_Continuation(r, objExpression, scope))
			}
			r.ConsiderForIndependentRemoval(builder.BuildUnusedReadReference(getProp, getProp))
		}
	} else {
		// TODO(bradfordcsmith): add removal of `varName.propName;`
		r.MarkPropertyNameAsPinned(propertyName)
		r.TraverseNode(objectNode, scope)
	}
}

// TraverseIncrementOrDecrementOp.
//
// TODO(b/137380742): Combine with `traverseCompoundAssign`.
//
// Java: RemoveUnusedCode.traverseIncrementOrDecrementOp(Node, Scope)
func (r *RemoveUnusedCode) TraverseIncrementOrDecrementOp(incOrDecOp *rhino.Node, scope *Scope) {
	jutil.CheckArgumentMsg1(incOrDecOp.IsInc() || incOrDecOp.IsDec(), "%s", incOrDecOp)
	arg := incOrDecOp.GetOnlyChild()
	if NodeUtil_IsExpressionResultUsed(incOrDecOp) {
		// If expression result is used, then this expression is definitely not removable.
		r.TraverseNode(arg, scope)
	} else if arg.IsGetProp() {
		getPropObj := arg.GetFirstChild()

		if r.ConsiderForAccessorSideEffects(arg, AccessorSummary_PropertyAccessKind_GETTER_AND_SETTER) {
			r.TraverseNode(getPropObj, scope) // Don't re-traverse the GETPROP as a read.
		} else if getPropObj.IsThis() {
			// this.propName++
			builder := NewRemoveUnusedCode_RemovableBuilder(r).SetIsThisDotPropertyReference(true)
			r.ConsiderForIndependentRemoval(builder.BuildIncOrDepOp(incOrDecOp, arg, nil))
		} else if RemoveUnusedCode_IsDotPrototype(getPropObj) {
			// someExpression.prototype.propName++
			exprObj := getPropObj.GetFirstChild()
			builder := NewRemoveUnusedCode_RemovableBuilder(r).SetIsPrototypeDotPropertyReference(true)
			if exprObj.IsName() {
				// varName.prototype.propName++
				varInfo := r.TraverseNameNode(exprObj, scope)
				varInfo.AddRemovable(builder.BuildIncOrDepOp(incOrDecOp, arg, nil))
			} else {
				// (someExpression).prototype.propName++
				var toPreserve *rhino.Node
				if r.MayHaveSideEffects(exprObj) {
					toPreserve = exprObj
					r.TraverseNode(exprObj, scope)
				} else {
					builder.AddContinuation(NewRemoveUnusedCode_Continuation(r, exprObj, scope))
				}
				r.ConsiderForIndependentRemoval(builder.BuildIncOrDepOp(incOrDecOp, arg, toPreserve))
			}
		} else {
			// someExpression.propName++ is not removable except in the cases covered above
			r.TraverseNode(arg, scope)
		}
	} else {
		// TODO(bradfordcsmith): varName++ should be removable if varName is otherwise unused
		r.TraverseNode(arg, scope)
	}
}

// TraverseCompoundAssign.
//
// TODO(b/137380742): Combine with `traverseIncrementOrDecrement`.
//
// Java: RemoveUnusedCode.traverseCompoundAssign(Node, Scope)
func (r *RemoveUnusedCode) TraverseCompoundAssign(compoundAssignNode *rhino.Node, scope *Scope) {
	// We'll allow removal of compound assignment to a `this` property as long as the result of the
	// assignment is unused.
	// e.g. `this.prop += 3;`

	// NOTE: Some history here, as there were questions about "why is 'this' special".  The "remove
	// unused properties" is not a general property removal algorithm.  It only removes unreferenced
	// properties that are part of class definitions.  "SomeClass.prop += 3" and
	// "SomeClass.prototype.prop += 3" could so be candidates but they aren't considered here.

	targetNode := compoundAssignNode.GetFirstChild()
	valueNode := compoundAssignNode.GetLastChild()
	if targetNode.IsGetProp() {
		if r.ConsiderForAccessorSideEffects(targetNode, AccessorSummary_PropertyAccessKind_GETTER_AND_SETTER) {
			r.TraverseNode(targetNode.GetFirstChild(), scope) // Don't re-traverse the GETPROP as a read.
			r.TraverseNode(valueNode, scope)
		} else if targetNode.GetFirstChild().IsThis() &&
			!NodeUtil_IsExpressionResultUsed(compoundAssignNode) {
			builder := NewRemoveUnusedCode_RemovableBuilder(r).SetIsThisDotPropertyReference(true)
			r.TraverseRemovableAssignValue(valueNode, builder, scope)
			r.ConsiderForIndependentRemoval(
				builder.BuildNamedPropertyAssign(compoundAssignNode, targetNode))
		} else {
			r.TraverseNode(targetNode, scope)
			r.TraverseNode(valueNode, scope)
		}
	} else {
		r.TraverseNode(targetNode, scope)
		r.TraverseNode(valueNode, scope)
	}
}

// Java: RemoveUnusedCode.traverseNameNode(Node, Scope)
func (r *RemoveUnusedCode) TraverseNameNode(n *rhino.Node, scope *Scope) RemoveUnusedCode_VarInfo {
	name := n.GetString()
	if r.polyfills != nil && !r.polyfills.IsEmpty() && r.polyfills.ContainsKey(name) {
		for info := range r.polyfills.Get(name).All() {
			if info.GetRemoveUnusedCode_PolyfillInfoBase().isRemovable {
				info.ConsiderPossibleReference(n)
			}
		}
	}

	return r.TraverseVar(jutil.CheckNotNullMsg1(scope.GetVar(name), "%s", n))
}

// Java: RemoveUnusedCode.traverseCall(Node, Scope)
func (r *RemoveUnusedCode) TraverseCall(callNode *rhino.Node, scope *Scope) {
	callee := callNode.GetFirstChild()

	if r.codingConvention.IsPropertyRenameFunction(callee) {
		propertyNameNode := callee.GetNext()
		if propertyNameNode != nil && propertyNameNode.IsStringLit() {
			r.MarkPropertyNameAsPinned(propertyNameNode.GetString())
		}
		r.TraverseChildren(callNode, scope)
	} else if NodeUtil_IsObjectDefinePropertiesDefinition(callNode) {
		// TODO(bradfordcsmith): Should also handle Object.create() and Object.defineProperty().
		r.TraverseObjectDefinePropertiesCall(callNode, scope)
	} else if r.removeUnusedPolyfills && RemoveUnusedCode_IsJscompPolyfill(callee) {
		firstArg := callee.GetNext()
		polyfillName := firstArg.GetString()
		if strings.HasSuffix(callee.GetString(), "polyfillTypedArrayMethod") {
			polyfillName = "TypedArray.prototype." + polyfillName
		}
		info := r.CreatePolyfillInfo(callNode, scope, polyfillName)
		if r.polyfills == nil {
			r.polyfills = jutil.HashMultimap_Create[string, RemoveUnusedCode_PolyfillInfo]()
		}
		r.polyfills.Put(info.GetRemoveUnusedCode_PolyfillInfoBase().key, info)
		// Only traverse the callee (to mark it as used).  The arguments may be traversed later.
		r.TraverseNode(callNode.GetFirstChild(), scope)
	} else if NodeUtil_IsGoogWeakUsageCall(callNode) &&
		callNode.HasTwoChildren() &&
		callNode.GetSecondChild().IsName() {
		// goog.weakUsage() should have exactly one argument, and it should be either a name or a
		// qualified name (this is checked in ProcessClosurePrimitives.java).
		//
		// If it is a qualified name, then we do not attempt to remove it at this time (see condition
		// above). We rely on CollapseProperties to turn qualified names into simple names where
		// possible.

		// Mark this call as removable if the var is not otherwise referenced.
		varInfo := r.TraverseNameNode(callNode.GetSecondChild(), scope)
		builder := NewRemoveUnusedCode_RemovableBuilder(r)
		varInfo.AddRemovable(builder.BuildWeakUsageCall(callNode))

		// We need to traverse the goog.weakUsage function itself (to mark it as used, in case our
		// usage of it is not removed).
		r.TraverseNode(callNode.GetFirstChild(), scope)
	} else {
		parent := callNode.GetParent()
		var classVarName *string
		classDefiningCall := false

		// A call that is a statement unto itself or the left side of a comma expression might be
		// a call to a known method for doing class setup
		// e.g. $jscomp.inherits(Class, BaseClass) or goog.addSingletonGetter(Class)
		// Such methods never have meaningful return values, so we won't look for them in other
		// contexts
		if parent.IsExprResult() || (parent.IsComma() && parent.GetFirstChild() == callNode) {
			subclassRelationship :=
				r.codingConvention.GetClassesDefinedByCall(callNode)
			if subclassRelationship != nil {
				// e.g. goog.inherits(DerivedClass, BaseClass);
				// NOTE: DerivedClass and BaseClass must be QNames. Otherwise getClassesDefinedByCall()
				// will return null.
				subclassName := subclassRelationship.subclassName
				classVarName = &subclassName
				classDefiningCall = true
			} else {
				// Look for calls to addSingletonGetter calls.
				classVarName = r.codingConvention.GetSingletonGetterClassName(callNode)
			}
		}

		var classVar *Var
		if classVarName != nil && NodeUtil_IsValidSimpleName(*classVarName) {
			classVar = jutil.CheckNotNullMsg1(scope.GetVar(*classVarName), "%s", *classVarName)
		}

		if classVar == nil || !classVar.IsGlobal() {
			// The call we are traversing does not modify a class definition,
			// or the class is not specified with a simple variable name,
			// or the variable name is not global.
			// TODO(bradfordcsmith): It would be more correct to check whether the class name
			// references a known constructor and expand to allow QNames.
			r.TraverseChildren(callNode, scope)
		} else {
			builder := NewRemoveUnusedCode_RemovableBuilder(r)
			for child := callNode.GetFirstChild(); child != nil; child = child.GetNext() {
				builder.AddContinuation(NewRemoveUnusedCode_Continuation(r, child, scope))
			}
			r.TraverseVar(classVar).
				AddRemovable(builder.BuildClassSetupCallWithClassDefiningCall(callNode, classDefiningCall))
		}
	}
}

// RemoveUnusedCode_IsJscompPolyfill checks whether this is a recognizable
// call to $jscomp.polyfill.
//
// Java: RemoveUnusedCode.isJscompPolyfill(Node)
func RemoveUnusedCode_IsJscompPolyfill(n *rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_NAME:
		{
			// Need to work correctly after CollapseProperties.
			name := n.GetString()
			return (name == "$jscomp$polyfill" ||
				name == "$jscomp$patch" ||
				name == "$jscomp$polyfillTypedArrayMethod") &&
				n.GetNext().IsStringLit()
		}
	case rhino.Token_GETPROP:
		{
			// Need to work correctly without CollapseProperties.
			propertyName := n.GetString()
			return (propertyName == "polyfill" ||
				propertyName == "patch" ||
				propertyName == "polyfillTypedArrayMethod") &&
				n.GetFirstChild().IsName() &&
				n.GetFirstChild().GetString() == "$jscomp" &&
				n.GetNext().IsStringLit()
		}
	default:
		return false
	}
}

// TraverseObjectDefinePropertiesCall traverses
// `Object.defineProperties(someObject, propertyDefinitions);`.
//
// Java: RemoveUnusedCode.traverseObjectDefinePropertiesCall(Node, Scope)
func (r *RemoveUnusedCode) TraverseObjectDefinePropertiesCall(callNode *rhino.Node, scope *Scope) {
	// First child is Object.defineProperties or some equivalent of it.
	callee := callNode.GetFirstChild()
	targetObject := callNode.GetSecondChild()
	propertyDefinitions := targetObject.GetNext()

	if (targetObject.IsName() || r.IsNameDotPrototype(targetObject)) &&
		!NodeUtil_IsExpressionResultUsed(callNode) {
		// NOTE: Object.defineProperties() returns its first argument, so if its return value is used
		// that counts as a use of the targetObject.
		var nameNode *rhino.Node
		if targetObject.IsName() {
			nameNode = targetObject
		} else {
			nameNode = targetObject.GetFirstChild()
		}
		varInfo := r.TraverseNameNode(nameNode, scope)
		builder := NewRemoveUnusedCode_RemovableBuilder(r)
		// TODO(bradfordcsmith): Is it really necessary to traverse the callee
		// (aka. Object.defineProperties)?
		builder.AddContinuation(NewRemoveUnusedCode_Continuation(r, callee, scope))
		if r.MayHaveSideEffects(propertyDefinitions) {
			r.TraverseNode(propertyDefinitions, scope)
		} else {
			builder.AddContinuation(NewRemoveUnusedCode_Continuation(r, propertyDefinitions, scope))
		}
		varInfo.AddRemovable(builder.BuildClassSetupCall(callNode))
	} else {
		// TODO(bradfordcsmith): Is it really necessary to traverse the callee
		// (aka. Object.defineProperties)?
		r.TraverseNode(callee, scope)
		r.TraverseNode(targetObject, scope)
		r.TraverseNode(propertyDefinitions, scope)
	}
}

// TraverseObjectDefinePropertiesLiteral traverses the object literal passed
// as the second argument to `Object.defineProperties()`.
//
// Java: RemoveUnusedCode.traverseObjectDefinePropertiesLiteral(Node, Scope)
func (r *RemoveUnusedCode) TraverseObjectDefinePropertiesLiteral(propertyDefinitions *rhino.Node, scope *Scope) {
	for property := propertyDefinitions.GetFirstChild(); property != nil; property = property.GetNext() {
		if property.IsQuotedStringKey() {
			// Quoted property name counts as a reference to the property and protects it from removal.
			r.MarkPropertyNameAsPinned(property.GetString())
			r.TraverseNode(property.GetOnlyChild(), scope)
		} else if property.IsStringKey() {
			definition := property.GetOnlyChild()
			if r.MayHaveSideEffects(definition) {
				r.TraverseNode(definition, scope)
			} else {
				r.ConsiderForIndependentRemoval(
					NewRemoveUnusedCode_RemovableBuilder(r).
						AddContinuation(NewRemoveUnusedCode_Continuation(r, definition, scope)).
						BuildObjectDefinePropertiesDefinition(property))
			}
		} else {
			// TODO(bradfordcsmith): Maybe report error for anything other than a computed property,
			// since getters, setters, and methods don't make much sense in this context.
			r.TraverseNode(property, scope)
		}
	}
}

// Java: RemoveUnusedCode.getVarForNameNode(Node, Scope)
func (r *RemoveUnusedCode) GetVarForNameNode(nameNode *rhino.Node, scope *Scope) *Var {
	return jutil.CheckNotNullMsg1(scope.GetVar(nameNode.GetString()), "%s", nameNode)
}

// Java: RemoveUnusedCode.traverseObjectLiteral(Node, Scope)
func (r *RemoveUnusedCode) TraverseObjectLiteral(objectLiteral *rhino.Node, scope *Scope) {
	jutil.CheckArgumentMsg1(objectLiteral.IsObjectLit(), "%s", objectLiteral)
	// Is this an object literal that is assigned directly to a 'prototype' property?
	if r.IsAssignmentToPrototype(objectLiteral.GetParent()) {
		r.TraversePrototypeLiteral(objectLiteral, scope)
	} else if r.IsObjectDefinePropertiesSecondArgument(objectLiteral) {
		// TODO(bradfordcsmith): Consider restricting special handling of the properties literal to
		// cases where the target object is a known class, prototype, or this.
		r.TraverseObjectDefinePropertiesLiteral(objectLiteral, scope)
	} else {
		r.TraverseNonPrototypeObjectLiteral(objectLiteral, scope)
	}
}

// Java: RemoveUnusedCode.isObjectDefinePropertiesSecondArgument(Node)
func (r *RemoveUnusedCode) IsObjectDefinePropertiesSecondArgument(n *rhino.Node) bool {
	parent := n.GetParent()
	return NodeUtil_IsObjectDefinePropertiesDefinition(parent) && parent.GetLastChild() == n
}

// Java: RemoveUnusedCode.traverseNonPrototypeObjectLiteral(Node, Scope)
func (r *RemoveUnusedCode) TraverseNonPrototypeObjectLiteral(objectLiteral *rhino.Node, scope *Scope) {
	for propertyNode := objectLiteral.GetFirstChild(); propertyNode != nil; propertyNode = propertyNode.GetNext() {
		if propertyNode.IsStringKey() {
			// A property name in an object literal counts as a reference,
			// because of some reflection patterns.
			// Note that we are intentionally treating both quoted and unquoted keys as
			// references.
			r.MarkPropertyNameAsPinned(propertyNode.GetString())
			r.TraverseNode(propertyNode.GetFirstChild(), scope)
		} else {
			r.TraverseNode(propertyNode, scope)
		}
	}
}

// Java: RemoveUnusedCode.traversePrototypeLiteral(Node, Scope)
func (r *RemoveUnusedCode) TraversePrototypeLiteral(objectLiteral *rhino.Node, scope *Scope) {
	for propertyNode := objectLiteral.GetFirstChild(); propertyNode != nil; propertyNode = propertyNode.GetNext() {
		if propertyNode.IsComputedProp() || propertyNode.IsQuotedStringKey() {
			r.TraverseChildren(propertyNode, scope)
		} else {
			valueNode := propertyNode.GetOnlyChild()
			if r.MayHaveSideEffects(valueNode) {
				// TODO(bradfordcsmith): Ideally we should preserve the side-effect without keeping the
				// property itself alive.
				r.TraverseNode(valueNode, scope)
			} else {
				// If we've come this far, we already know we're keeping the prototype literal itself,
				// but we may be able to remove unreferenced properties in it.
				r.ConsiderForIndependentRemoval(
					NewRemoveUnusedCode_RemovableBuilder(r).
						AddContinuation(NewRemoveUnusedCode_Continuation(r, valueNode, scope)).
						BuildClassOrPrototypeNamedProperty(propertyNode))
			}
		}
	}
}

// Java: RemoveUnusedCode.isAssignmentToPrototype(Node)
func (r *RemoveUnusedCode) IsAssignmentToPrototype(n *rhino.Node) bool {
	return n.IsAssign() && RemoveUnusedCode_IsDotPrototype(n.GetFirstChild())
}

// RemoveUnusedCode_IsDotPrototype: true for `someExpression.prototype`.
//
// Java: RemoveUnusedCode.isDotPrototype(Node)
func RemoveUnusedCode_IsDotPrototype(n *rhino.Node) bool {
	return NodeUtil_IsNormalOrOptChainGetProp(n) && n.GetString() == "prototype"
}

// Java: RemoveUnusedCode.traverseCatch(Node, Scope)
func (r *RemoveUnusedCode) TraverseCatch(catchNode *rhino.Node, scope *Scope) {
	exceptionNameNode := catchNode.GetFirstChild()
	block := exceptionNameNode.GetNext()
	if exceptionNameNode.IsName() {
		// exceptionNameNode can be an empty node if not using a binding in 2019.
		exceptionVarInfo := r.TraverseNameNode(exceptionNameNode, scope)
		exceptionVarInfo.SetIsExplicitlyNotRemovable(func() string { return "catch variable" })
	} else {
		r.TraverseNode(exceptionNameNode, scope)
	}
	r.TraverseNode(block, scope)
}

// Java: RemoveUnusedCode.traverseEnhancedFor(Node, Scope)
func (r *RemoveUnusedCode) TraverseEnhancedFor(enhancedFor *rhino.Node, scope *Scope) {
	forScope := scope
	if hasBlockScopedDeclarations(enhancedFor) {
		forScope = r.scopeCreator.CreateScopeOfScope(enhancedFor, scope)
	}
	// for (iterationTarget in|of collection) body;
	iterationTarget := enhancedFor.GetFirstChild()
	collection := iterationTarget.GetNext()
	body := collection.GetNext()
	if iterationTarget.IsName() {
		// using previously-declared loop variable. e.g.
		// `for (varName of collection) {}`
		varInfo := r.TraverseNameNode(iterationTarget, forScope)
		varInfo.SetIsExplicitlyNotRemovable(func() string { return "for-of or for-in loop variable" })
	} else if NodeUtil_IsNameDeclaration(iterationTarget) {
		// loop has const/var/let declaration
		declNode := iterationTarget.GetOnlyChild()
		if declNode.IsDestructuringLhs() {
			// e.g.
			// `for (const [a, b] of pairList) {}`
			// destructuring is handled at a lower level
			// Note that destructuring assignments are always considered to set an unknown value
			// equivalent to what we set for the var name case above and below.
			// It isn't necessary to set the variable names as not removable, though, because the
			// thing that isn't removable is the destructuring pattern itself, which we never remove.
			// TODO(bradfordcsmith): The need to explain all the above shows this should be reworked.
			r.TraverseNode(declNode, forScope)
		} else {
			// e.g.
			// `for (const varName of collection) {}`
			jutil.CheckState(declNode.IsName())
			jutil.CheckState(!declNode.HasChildren())
			// We can never remove the loop variable of a for-in or for-of loop, because it's
			// essential to loop syntax.
			varInfo := r.TraverseNameNode(declNode, forScope)
			varInfo.SetIsExplicitlyNotRemovable(func() string { return "for-of or for-in loop variable" })
		}
	} else {
		// using some general LHS value e.g.
		// `for ([a, b] of collection) {}` destructuring with existing vars
		// `for (a.x of collection) {}` using a property as the loop var
		// TODO(bradfordcsmith): This should be considered a write if it's a property reference.
		r.TraverseNode(iterationTarget, forScope)
	}
	r.TraverseNode(collection, forScope)
	r.TraverseNode(body, forScope)
}

// Java: RemoveUnusedCode.traverseVanillaFor(Node, Scope)
func (r *RemoveUnusedCode) TraverseVanillaFor(forNode *rhino.Node, scope *Scope) {
	forScope := scope
	if hasBlockScopedDeclarations(forNode) {
		forScope = r.scopeCreator.CreateScopeOfScope(forNode, scope)
	}
	initialization := forNode.GetFirstChild()
	condition := initialization.GetNext()
	update := condition.GetNext()
	block := update.GetNext()
	if NodeUtil_IsNameDeclaration(initialization) {
		r.TraverseVanillaForNameDeclarations(initialization, forScope)
	} else {
		r.TraverseNode(initialization, forScope)
	}
	r.TraverseNode(condition, forScope)
	r.TraverseNode(update, forScope)
	r.TraverseNode(block, forScope)
}

// Java: RemoveUnusedCode.traverseVanillaForNameDeclarations(Node, Scope)
func (r *RemoveUnusedCode) TraverseVanillaForNameDeclarations(nameDeclaration *rhino.Node, scope *Scope) {
	for child := nameDeclaration.GetFirstChild(); child != nil; child = child.GetNext() {
		if !child.IsName() {
			// TODO(bradfordcsmith): Customize handling of destructuring
			r.TraverseNode(child, scope)
		} else {
			nameNode := child
			valueNode := child.GetFirstChild()
			varInfo := r.TraverseNameNode(nameNode, scope)
			if valueNode == nil {
				varInfo.AddRemovable(NewRemoveUnusedCode_RemovableBuilder(r).BuildVanillaForNameDeclaration(nameNode))
			} else if r.MayHaveSideEffects(valueNode) {
				// TODO(bradfordcsmith): Actually allow for removing the variable while keeping the
				// valueNode for its side-effects.
				varInfo.SetIsExplicitlyNotRemovable(
					func() string { return "for-loop variable initialization has side-effects" })
				r.TraverseNode(valueNode, scope)
			} else {
				vanillaForNameDeclaration :=
					NewRemoveUnusedCode_RemovableBuilder(r).
						AddContinuation(NewRemoveUnusedCode_Continuation(r, valueNode, scope)).
						BuildVanillaForNameDeclaration(nameNode)
				varInfo.AddRemovable(vanillaForNameDeclaration)
			}
		}
	}
}

// Java: RemoveUnusedCode.traverseDeclarationStatement(Node, Scope)
func (r *RemoveUnusedCode) TraverseDeclarationStatement(declarationStatement *rhino.Node, scope *Scope) {
	// Normalization should ensure that declaration statements always have just one child.
	nameNode := declarationStatement.GetOnlyChild()
	if !nameNode.IsName() {
		// Destructuring declarations are handled elsewhere.
		r.TraverseNode(nameNode, scope)
	} else {
		valueNode := nameNode.GetFirstChild()
		varInfo := r.TraverseNameNode(nameNode, scope)
		builder := NewRemoveUnusedCode_RemovableBuilder(r)
		if valueNode == nil {
			varInfo.AddRemovable(builder.BuildNameDeclarationStatement(declarationStatement))
		} else {
			if r.MayHaveSideEffects(valueNode) {
				r.TraverseNode(valueNode, scope)
			} else {
				builder.AddContinuation(NewRemoveUnusedCode_Continuation(r, valueNode, scope))
			}
			removable :=
				builder.BuildNameDeclarationStatement(declarationStatement)
			varInfo.AddRemovable(removable)
		}
	}
}

// Java: RemoveUnusedCode.traverseAssign(Node, Scope)
func (r *RemoveUnusedCode) TraverseAssign(assignNode *rhino.Node, scope *Scope) {
	jutil.CheckState(NodeUtil_IsAssignmentOp(assignNode))

	lhs := assignNode.GetFirstChild()
	valueNode := assignNode.GetLastChild()
	if lhs.IsName() {
		// varName = something
		varInfo := r.TraverseNameNode(lhs, scope)
		builder := NewRemoveUnusedCode_RemovableBuilder(r)
		r.TraverseRemovableAssignValue(valueNode, builder, scope)
		varInfo.AddRemovable(builder.BuildVariableAssign(assignNode, varInfo))
	} else if lhs.IsGetElem() {
		getElemObj := lhs.GetFirstChild()
		getElemKey := lhs.GetLastChild()
		var varNameNode *rhino.Node
		if getElemObj.IsName() {
			varNameNode = getElemObj
		} else if r.IsNameDotPrototype(getElemObj) {
			varNameNode = getElemObj.GetFirstChild()
		} else {
			varNameNode = nil
		}

		if varNameNode != nil {
			// varName[someExpression] = someValue
			// OR
			// varName.prototype[someExpression] = someValue
			varInfo := r.TraverseNameNode(varNameNode, scope)
			builder := NewRemoveUnusedCode_RemovableBuilder(r)
			if r.MayHaveSideEffects(getElemKey) {
				r.TraverseNode(getElemKey, scope)
			} else {
				builder.AddContinuation(NewRemoveUnusedCode_Continuation(r, getElemKey, scope))
			}
			r.TraverseRemovableAssignValue(valueNode, builder, scope)
			varInfo.AddRemovable(builder.BuildComputedPropertyAssign(assignNode, getElemKey, varInfo))
		} else {
			r.TraverseNode(getElemObj, scope)
			r.TraverseNode(getElemKey, scope)
			r.TraverseNode(valueNode, scope)
		}
	} else if lhs.IsGetProp() {
		getPropLhs := lhs.GetFirstChild()
		// Assignments `Foo.prototype.bar = function() {`
		isDotPrototypeLhs := RemoveUnusedCode_IsDotPrototype(getPropLhs)
		isPrototypeMethodDef := isDotPrototypeLhs && valueNode.IsFunction()

		if !isPrototypeMethodDef &&
			r.ConsiderForAccessorSideEffects(lhs, AccessorSummary_PropertyAccessKind_SETTER_ONLY) {
			// And the possible side-effects mean we can't do any removal. We don't use the
			// `AstAnalyzer` because we only want to consider side-effect from the assignment, not the
			// entire l-value subtree.
			// Assume prototype method assignments never trigger setters, matching ES class semantics
			r.TraverseNode(getPropLhs, scope) // Don't re-traverse the GETPROP as a read.
			r.TraverseNode(valueNode, scope)
		} else if getPropLhs.IsName() {
			// varName.propertyName = someValue
			varInfo := r.TraverseNameNode(getPropLhs, scope)
			builder := NewRemoveUnusedCode_RemovableBuilder(r)
			r.TraverseRemovableAssignValue(valueNode, builder, scope)
			varInfo.AddRemovable(builder.BuildNamedPropertyAssignWithVarInfo(assignNode, lhs, varInfo))
		} else if isDotPrototypeLhs {
			// objExpression.prototype.propertyName = someValue
			objExpression := getPropLhs.GetFirstChild()
			builder := NewRemoveUnusedCode_RemovableBuilder(r).SetIsPrototypeDotPropertyReference(true)
			r.TraverseRemovableAssignValue(valueNode, builder, scope)
			if objExpression.IsName() {
				// varName.prototype.propertyName = someValue
				varInfo := r.TraverseNameNode(getPropLhs.GetFirstChild(), scope)
				varInfo.AddRemovable(builder.BuildNamedPropertyAssignWithVarInfo(assignNode, lhs, varInfo))
			} else {
				// (someExpression).prototype.propertyName = someValue
				if r.MayHaveSideEffects(objExpression) {
					r.TraverseNode(objExpression, scope)
				} else {
					builder.AddContinuation(NewRemoveUnusedCode_Continuation(r, objExpression, scope))
				}
				r.ConsiderForIndependentRemoval(
					builder.BuildAnonymousPrototypeNamedPropertyAssign(assignNode, lhs.GetString()))
			}
		} else if getPropLhs.IsThis() {
			// this.propertyName = someValue
			builder := NewRemoveUnusedCode_RemovableBuilder(r).SetIsThisDotPropertyReference(true)
			r.TraverseRemovableAssignValue(valueNode, builder, scope)
			r.ConsiderForIndependentRemoval(builder.BuildNamedPropertyAssign(assignNode, lhs))
		} else {
			r.TraverseNode(lhs, scope)
			r.TraverseNode(valueNode, scope)
		}
	} else {
		// no other cases are removable
		r.TraverseNode(lhs, scope)
		r.TraverseNode(valueNode, scope)
	}
}

// Java: RemoveUnusedCode.traverseRemovableAssignValue(Node, RemovableBuilder, Scope)
func (r *RemoveUnusedCode) TraverseRemovableAssignValue(valueNode *rhino.Node, builder *RemoveUnusedCode_RemovableBuilder, scope *Scope) {
	if r.MayHaveSideEffects(valueNode) || NodeUtil_IsExpressionResultUsed(valueNode.GetParent()) {
		r.TraverseNode(valueNode, scope)
	} else {
		builder.AddContinuation(NewRemoveUnusedCode_Continuation(r, valueNode, scope))
	}
}

// Java: RemoveUnusedCode.isNameDotPrototype(Node)
func (r *RemoveUnusedCode) IsNameDotPrototype(n *rhino.Node) bool {
	return n.IsGetProp() && n.GetFirstChild().IsName() && n.GetString() == "prototype"
}

// Java: RemoveUnusedCode.traverseObjectPattern(Node, Scope)
func (r *RemoveUnusedCode) TraverseObjectPattern(pattern *rhino.Node, scope *Scope) {
	jutil.CheckStateMsg1(pattern.IsObjectPattern(), "%s", pattern)

	for elem := pattern.GetFirstChild(); elem != nil; elem = elem.GetNext() {
		switch elem.GetToken() {
		case rhino.Token_COMPUTED_PROP:
			r.TraverseIndirectAssignment(elem, elem.GetSecondChild(), scope)
		case rhino.Token_STRING_KEY:
			if !elem.IsQuotedStringKey() {
				r.MarkPropertyNameAsPinned(elem.GetString())
			}
			r.TraverseIndirectAssignment(elem, elem.GetOnlyChild(), scope)
		case rhino.Token_ITER_REST, rhino.Token_OBJECT_REST:
			// Recall that the rest target can be any l-value expression
			r.TraverseIndirectAssignment(elem, elem.GetOnlyChild(), scope)
		default:
			panic(jutil.NewIllegalStateExceptionMsg(
				"Unexpected child of " + pattern.GetToken().String() + ": " + elem.ToStringTree()))
		}
	}
}

// Java: RemoveUnusedCode.traverseIndirectAssignmentList(Node, Scope)
func (r *RemoveUnusedCode) TraverseIndirectAssignmentList(list *rhino.Node, scope *Scope) {
	jutil.CheckStateMsg1(list.IsArrayPattern() || list.IsParamList(), "%s", list)

	for elem := list.GetFirstChild(); elem != nil; elem = elem.GetNext() {
		switch elem.GetToken() {
		case rhino.Token_EMPTY:
		case rhino.Token_ARRAY_PATTERN, rhino.Token_DEFAULT_VALUE, rhino.Token_GETELEM, rhino.Token_GETPROP, rhino.Token_NAME, rhino.Token_OBJECT_PATTERN:
			r.TraverseIndirectAssignment(elem, elem, scope)
		case rhino.Token_ITER_REST, rhino.Token_OBJECT_REST:
			r.TraverseIndirectAssignment(elem, elem.GetOnlyChild(), scope)
		default:
			panic(jutil.NewIllegalStateExceptionMsg(
				"Unexpected child of " + list.GetToken().String() + ": " + elem.ToStringTree()))
		}
	}
}

// TraverseIndirectAssignment traverses an AST structure representing an
// assignment operation for which the target and value are far apart.
//
// Examples include destructurings and function parameters.
//
// root is the root of the assignment subtree; target the l-value expression
// being assigned to.
//
// Java: RemoveUnusedCode.traverseIndirectAssignment(Node, Node, Scope)
func (r *RemoveUnusedCode) TraverseIndirectAssignment(root *rhino.Node, target *rhino.Node, scope *Scope) {
	rootParent := root.GetParent()
	jutil.CheckArgumentMsg1(rootParent.IsDestructuringPattern() || rootParent.IsParamList(), "%s", rootParent)

	// Flatten out the case where the target is a default value. We always have to consider it.
	if target.IsDefaultValue() {
		target = target.GetFirstChild()
	}

	if target.IsGetProp() {
		r.ConsiderForAccessorSideEffects(target, AccessorSummary_PropertyAccessKind_SETTER_ONLY)
	}

	builder :=
		NewRemoveUnusedCode_RemovableBuilder(r).AddContinuation(NewRemoveUnusedCode_Continuation(r, root, scope))

	if r.MayHaveSideEffects(root) {
		// If anywhere in the assignment subtree has side-effects, it means that even if the target is
		// removable the subtree is not.
		r.TraverseNode(root, scope)
		// TODO(bradfordcsmith): Preserve side effects without preventing removal of variables and
		// properties. We could probably do this by subbing in an empty object pattern.
	} else if target.IsName() {
		varInfo := r.TraverseNameNode(target, scope)
		varInfo.AddRemovable(builder.BuildIndirectAssign(root, target))
	} else if r.IsNameDotPrototype(target) || RemoveUnusedCode_IsThisDotProperty(target) {
		r.ConsiderForIndependentRemoval(builder.BuildIndirectAssign(root, target))
	} else {
		// TODO(bradfordcsmith): Handle property assignments also
		// e.g. `({a: foo.bar, b: foo.baz}) = {a: 1, b: 2}`
		r.TraverseNode(root, scope)
	}
}

// Java: RemoveUnusedCode.traverseChildren(Node, Scope)
func (r *RemoveUnusedCode) TraverseChildren(n *rhino.Node, scope *Scope) {
	for c := n.GetFirstChild(); c != nil; c = c.GetNext() {
		r.TraverseNode(c, scope)
	}
}

// TraverseClass handles a class that is not the RHS child of an assignment
// or a variable declaration initializer.
//
// Java: RemoveUnusedCode.traverseClass(Node, Scope)
func (r *RemoveUnusedCode) TraverseClass(classNode *rhino.Node, scope *Scope) {
	jutil.CheckArgument(classNode.IsClass())
	if NodeUtil_IsClassDeclaration(classNode) {
		r.TraverseClassDeclaration(classNode, scope)
	} else {
		r.TraverseClassExpression(classNode, scope)
	}
}

// Java: RemoveUnusedCode.traverseClassDeclaration(Node, Scope)
func (r *RemoveUnusedCode) TraverseClassDeclaration(classNode *rhino.Node, scope *Scope) {
	jutil.CheckArgument(classNode.IsClass())
	classNameNode := classNode.GetFirstChild()
	baseClassExpression := classNameNode.GetNext()
	classBodyNode := baseClassExpression.GetNext()
	// Class declarations never declare the class name in the inner class scope (only class
	// expressions do), so creating an inner scope here would always produce an empty scope.
	classScope := scope

	varInfo := r.TraverseNameNode(classNameNode, scope)
	if classNode.GetParent().IsExport() {
		// Cannot remove an exported class.
		varInfo.SetIsExplicitlyNotRemovable(func() string { return "exported class" })
		r.TraverseNode(baseClassExpression, scope)
		// Use traverseChildren() here, because we should not consider any properties on the exported
		// class to be removable.
		r.TraverseChildren(classBodyNode, classScope)
	} else if r.MayHaveSideEffects(baseClassExpression) {
		// TODO(bradfordcsmith): implement removal without losing side-effects for this case
		varInfo.SetIsExplicitlyNotRemovable(func() string { return "base class expression has side-effects" })
		r.TraverseNode(baseClassExpression, scope)
		r.TraverseClassMembers(classBodyNode, classScope)
	} else if r.MayHaveSideEffects(classBodyNode) {
		varInfo.SetIsExplicitlyNotRemovable(func() string { return "class body has side-effects" })
		r.TraverseNode(baseClassExpression, scope)
		r.TraverseClassMembers(classBodyNode, classScope)
	} else {
		builder :=
			NewRemoveUnusedCode_RemovableBuilder(r).
				AddContinuation(NewRemoveUnusedCode_Continuation(r, baseClassExpression, classScope)).
				AddContinuation(NewRemoveUnusedCode_Continuation(r, classBodyNode, classScope))
		varInfo.AddRemovable(builder.BuildClassDeclaration(classNode))
	}
}

// Java: RemoveUnusedCode.traverseClassExpression(Node, Scope)
func (r *RemoveUnusedCode) TraverseClassExpression(classNode *rhino.Node, scope *Scope) {
	jutil.CheckArgument(classNode.IsClass())
	classNameNode := classNode.GetFirstChild()
	baseClassExpression := classNameNode.GetNext()
	classBodyNode := baseClassExpression.GetNext()
	classScope := scope

	if classNameNode.IsName() {
		classScope = r.scopeCreator.CreateScopeOfScope(classNode, scope)
		// We may be able to remove the name node if nothing ends up referring to it.
		varInfo := r.TraverseNameNode(classNameNode, classScope)
		// The class is non-local, because it is accessible by unknown code outside
		// of the scope where InnerName is defined.
		// e.g. `use(class InnerName {})`
		varInfo.SetHasNonLocalOrNonLiteralValue()
		varInfo.AddRemovable(NewRemoveUnusedCode_RemovableBuilder(r).BuildNamedClassExpression(classNode))
	}
	// If we're traversing the class expression, we've already decided we cannot remove it.
	r.TraverseNode(baseClassExpression, scope)
	r.TraverseClassMembers(classBodyNode, classScope)
}

// Java: RemoveUnusedCode.traverseClassMembers(Node, Scope)
func (r *RemoveUnusedCode) TraverseClassMembers(node *rhino.Node, scope *Scope) {
	jutil.CheckArgumentMsg1(node.IsClassMembers(), "%s", node)
	if !r.removeUnusedPrototypeProperties {
		r.TraverseChildren(node, scope)
		return
	}

	for member := node.GetFirstChild(); member != nil; member = member.GetNext() {
		switch member.GetToken() {
		case rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF, rhino.Token_MEMBER_FUNCTION_DEF:
			// If we get as far as traversing the members of a class, we've already decided that
			// we cannot remove the class itself, so just consider individual members for removal.
			r.ConsiderForIndependentRemoval(
				NewRemoveUnusedCode_RemovableBuilder(r).
					AddContinuation(NewRemoveUnusedCode_Continuation(r, member, scope)).
					BuildClassOrPrototypeNamedProperty(member))
		case rhino.Token_MEMBER_FIELD_DEF:
			// TODO(bradfordcsmith): currently if the RHS of a field has side effects, we do not
			// remove any part of the field. The proper behavior of class C { x = alert(); }
			// would be to remove x, leaving class C { constructor() { alert(); } }
			// but currently we aren't removing anything.
			if !member.HasChildren() || !r.MayHaveSideEffects(member.GetFirstChild()) {
				r.ConsiderForIndependentRemoval(
					NewRemoveUnusedCode_RemovableBuilder(r).
						AddContinuation(NewRemoveUnusedCode_Continuation(r, member, scope)).
						BuildClassOrPrototypeNamedProperty(member))
			}
			// TODO: b/354704593 - remove the entire class when it is not referenced
			if member.HasChildren() {
				r.TraverseChildren(member, scope)
			}
		case rhino.Token_COMPUTED_PROP, rhino.Token_COMPUTED_FIELD_DEF:
			r.TraverseChildren(member, scope)
		case rhino.Token_BLOCK:
			r.TraverseChildren(member, r.scopeCreator.CreateScopeOfScope(member, scope))
		default:
			panic(jutil.NewIllegalStateExceptionMsg(
				"Unexpected child of CLASS_MEMBERS: " + member.ToStringTree()))
		}
	}
}

// TraverseFunction traverses a function
//
// ES6 scopes of a function include the parameter scope and the body scope of
// the function.
//
// Note that CATCH blocks also create a new scope, but only for the catch
// variable. Declarations within the block actually belong to the enclosing
// scope. Because we don't remove catch variables, there's no need to treat
// CATCH blocks differently like we do functions.
//
// Java: RemoveUnusedCode.traverseFunction(Node, Scope)
func (r *RemoveUnusedCode) TraverseFunction(function *rhino.Node, parentScope *Scope) {
	jutil.CheckStateMsg1(function.HasXChildren(3), "%s", function)
	jutil.CheckStateMsg1(function.IsFunction(), "%s", function)

	paramlist := NodeUtil_GetFunctionParameters(function)
	body := function.GetLastChild()
	jutil.CheckStateMsg1(body.GetNext() == nil && body.IsBlock(), "%s", body)

	// Checking the parameters
	fparamScope := r.scopeCreator.CreateScopeOfScope(function, parentScope)

	// Checking the function body
	fbodyScope := r.scopeCreator.CreateScopeOfScope(body, fparamScope)

	nameNode := function.GetFirstChild()
	if nameNode.GetString() != "" {
		// var x = function funcName() {};
		// make sure funcName gets into the varInfoMap so it will be considered for removal.
		varInfo := r.TraverseNameNode(nameNode, fparamScope)
		if NodeUtil_IsExpressionResultUsed(function) {
			// var f = function g() {};
			// The f is an alias for g, so g escapes from the scope where it is defined.
			varInfo.SetHasNonLocalOrNonLiteralValue()
		}
	}

	r.TraverseNode(paramlist, fparamScope)
	r.TraverseChildren(body, fbodyScope)

	r.allFunctionParamScopes = append(r.allFunctionParamScopes, fparamScope)
}

// Java: RemoveUnusedCode.canRemoveParameters(Node)
func (r *RemoveUnusedCode) CanRemoveParameters(parameterList *rhino.Node) bool {
	jutil.CheckState(parameterList.IsParamList())
	function := parameterList.GetParent()
	return r.removeGlobals && !NodeUtil_IsGetOrSetKey(function.GetParent())
}

// RemoveUnreferencedFunctionArgs removes unreferenced arguments from a
// function declaration and when possible the function's callSites.
//
// fparamScope is the function parameter scope.
//
// Java: RemoveUnusedCode.removeUnreferencedFunctionArgs(Scope)
func (r *RemoveUnusedCode) RemoveUnreferencedFunctionArgs(fparamScope *Scope) {
	// Notice that removing unreferenced function args breaks
	// Function.prototype.length. In advanced mode, we don't really care
	// about this: we consider "length" the equivalent of reflecting on
	// the function's lexical source.
	//
	// Rather than create a new option for this, we assume that if the user
	// is removing globals, then it's OK to remove unused function args.
	//
	// See http://blickly.github.io/closure-compiler-issues/#253
	if !r.removeGlobals {
		return
	}

	function := fparamScope.GetRootNode()
	jutil.CheckState(function.IsFunction())
	if NodeUtil_IsGetOrSetKey(function.GetParent()) {
		// The parameters object literal setters can not be removed.
		return
	}

	argList := NodeUtil_GetFunctionParameters(function)
	// Strip as many unreferenced args off the end of the function declaration as possible.
	r.MaybeRemoveUnusedTrailingParameters(argList, fparamScope)

	// Mark any remaining unused parameters are unused to OptimizeParameters can try to remove
	// them.
	r.MarkUnusedParameters(argList, fparamScope)
}

func (r *RemoveUnusedCode) isPropertyPinned(propertyName string) bool {
	if propertyName == r.lastPinnedProp && propertyName != "" {
		return true
	}
	if r.customImplicitProperties == nil {
		if isImplicitlyUsedProperty(propertyName) {
			r.lastPinnedProp = propertyName
			return true
		}
	} else if _, ok := r.customImplicitProperties[propertyName]; ok {
		r.lastPinnedProp = propertyName
		return true
	}
	if _, ok := r.pinnedPropertyNames[propertyName]; ok {
		r.lastPinnedProp = propertyName
		return true
	}
	if r.externProperties != nil && r.externProperties.Contains(propertyName) {
		r.pinnedPropertyNames[propertyName] = struct{}{}
		r.lastPinnedProp = propertyName
		return true
	}
	return false
}

// Java: RemoveUnusedCode.markPropertyNameAsPinned(String)
func (r *RemoveUnusedCode) MarkPropertyNameAsPinned(propertyName string) {
	if r.isPropertyPinned(propertyName) {
		return
	}
	r.pinnedPropertyNames[propertyName] = struct{}{}
	r.lastPinnedProp = propertyName
	// Continue traversal of all of the property name's values and no longer consider them for
	// removal.
	if list, ok := r.propRemovables[propertyName]; ok {
		delete(r.propRemovables, propertyName)
		for cur := list.head; cur != nil; {
			next := cur.propNext
			cur.propNext = nil
			cur.inPropList = false
			cur.self.ApplyContinuations()
			cur = next
		}
	}
}

// Java: RemoveUnusedCode.considerForIndependentRemoval(Removable)
func (r *RemoveUnusedCode) ConsiderForIndependentRemoval(removable RemoveUnusedCode_Removable) {
	if removable.IsNamedProperty() {
		propertyName := removable.GetPropertyName()

		if r.isPropertyPinned(propertyName) ||
			r.codingConvention.IsExportedWithLocal(propertyName /* local= */, false) {
			// Referenced or exported, so not removable.
			removable.ApplyContinuations()
		} else if r.IsIndependentlyRemovable(removable) {
			// Store for possible removal later.
			base := removable.GetRemoveUnusedCode_RemovableBase()
			if !base.inPropList {
				base.inPropList = true
				base.propNext = nil
				list, exists := r.propRemovables[propertyName]
				if !exists {
					r.propRemovableKeys = append(r.propRemovableKeys, propertyName)
					r.propRemovables[propertyName] = propRemovableList{head: base, tail: base}
				} else {
					list.tail.propNext = base
					list.tail = base
					r.propRemovables[propertyName] = list
				}
			}
		} else {
			removable.ApplyContinuations()
			// This assignment counts as a reference, since we won't be removing it.
			// This is necessary in order to preserve getters and setters for the property.
			r.MarkPropertyNameAsPinned(propertyName)
		}
	} else {
		removable.ApplyContinuations()
	}
}

// ConsiderForAccessorSideEffects returns whether or not accessor
// side-effect are a possibility.
//
// Java: RemoveUnusedCode.considerForAccessorSideEffects(Node, PropertyAccessKind)
func (r *RemoveUnusedCode) ConsiderForAccessorSideEffects(getprop *rhino.Node, usage AccessorSummary_PropertyAccessKind) bool {
	// Other node types may make sense in the future.
	jutil.CheckStateMsg1(NodeUtil_IsNormalOrOptChainGetProp(getprop), "%s", getprop)
	propName := getprop.GetString()
	summary := r.accessorSummary
	if summary == nil {
		summary = r.compiler.GetAccessorSummary()
	}
	recorded := summary.GetKind(propName)
	if (recorded.HasGetter() && usage.HasGetter() && !r.assumeGettersArePure) ||
		(recorded.HasSetter() && usage.HasSetter()) {
		r.MarkPropertyNameAsPinned(propName)
		return true
	}

	return false
}

// Java: RemoveUnusedCode.isIndependentlyRemovable(Removable)
func (r *RemoveUnusedCode) IsIndependentlyRemovable(removable RemoveUnusedCode_Removable) bool {
	if removable.IsPrototypeProperty() {
		// `foo.prototype.prop = something;`
		// `class C { prop() {} }`
		return r.removeUnusedPrototypeProperties
	} else if removable.IsObjectDefinePropertiesDefinition() {
		// `Object.defineProperties({ prop: {...}});`
		return r.removeUnusedObjectDefinePropertiesDefinitions
	} else if removable.IsThisDotPropertyReference() {
		// `this.prop = something;`
		return r.removeUnusedThisProperties
	} else if removable.IsStaticProperty() {
		// `class Foo { static prop() {} }`
		// `Foo.otherStaticProp = value;`
		// TODO(b/139319709): removeUnusedThisProperties has ended up covering more than it was
		// originally intended to cover for arbitrary reasons.
		return r.removeUnusedThisProperties
	} else {
		return false
	}
}

// MarkUnusedParameters marks any remaining unused parameters as being unused
// so it can be used elsewhere.
//
// paramList is the list of function's parameters.
//
// Java: RemoveUnusedCode.markUnusedParameters(Node, Scope)
func (r *RemoveUnusedCode) MarkUnusedParameters(paramList *rhino.Node, fparamScope *Scope) {
	jutil.CheckArgumentMsg1(paramList.IsParamList(), "%s", paramList)

	for param := paramList.GetFirstChild(); param != nil; param = param.GetNext() {
		if param.IsUnusedParameter() {
			// already marked
			continue
		}

		paramNameNode := RemoveUnusedCode_NameOfParam(param)
		if paramNameNode == nil {
			// destructuring pattern parameters don't have a name that applies to the whole parameter
			// TODO(bradfordcsmith): We could mark this if we determined that all vars created by
			// the pattern are unused.
			continue
		}

		varInfo := r.TraverseNameNode(paramNameNode, fparamScope)
		if varInfo.IsRemovable() {
			param.SetUnusedParameter(true)
			r.compiler.ReportChangeToEnclosingScope(paramList)
			if r.removalLog.IsLogging() {
				r.removalLog.LogSupplier(RemoveUnusedCode_RemovalLogRecord_ForMarkingNamedArg(paramNameNode, paramList).Get)
			}
		}
	}
}

// MaybeRemoveUnusedTrailingParameters strips as many unreferenced args off
// the end of the function declaration as possible. We start from the end of
// the function declaration because removing parameters from the middle of
// the param list could mess up the interpretation of parameters being sent
// over by any function calls.
//
// argList is the list of function's arguments.
//
// Java: RemoveUnusedCode.maybeRemoveUnusedTrailingParameters(Node, Scope)
func (r *RemoveUnusedCode) MaybeRemoveUnusedTrailingParameters(argList *rhino.Node, fparamScope *Scope) {
	jutil.CheckArgumentMsg1(argList.IsParamList(), "%s", argList)
	for {
		lastArg := argList.GetLastChild()
		if lastArg == nil {
			break
		}
		argNode := lastArg
		if lastArg.IsDefaultValue() {
			argNode = lastArg.GetFirstChild()
			if r.MayHaveSideEffects(lastArg.GetLastChild()) {
				break
			}
		}

		if argNode.IsRest() {
			argNode = argNode.GetFirstChild()
		}

		if argNode.IsDestructuringPattern() {
			if argNode.HasChildren() {
				// TODO(johnlenz): handle the case where there are no assignments.
				break
			} else {
				// Remove empty destructuring patterns and their associated object literal assignment
				// if it exists and if the right hand side does not have side effects. Note, a
				// destructuring pattern with a "leftover" property key as in {a:{}} is not considered
				// empty in this case!
				NodeUtil_DeleteNode(lastArg, r.compiler)
				if r.removalLog.IsLogging() {
					r.removalLog.LogSupplier(RemoveUnusedCode_RemovalLogRecord_ForDestructuringArg(argList).Get)
				}
				continue
			}
		}

		varInfo := r.GetVarInfo(r.GetVarForNameNode(argNode, fparamScope))
		if varInfo.IsRemovable() {
			NodeUtil_DeleteNode(lastArg, r.compiler)
			if r.removalLog.IsLogging() {
				r.removalLog.LogSupplier(RemoveUnusedCode_RemovalLogRecord_ForNamedArg(argNode, argList).Get)
			}
		} else {
			break
		}
	}
}

// TraverseVar handles a variable reference seen during traversal and
// returns a VarInfo object appropriate for the given Var.
//
// This is a wrapper for getVarInfo that handles additional logic needed when
// we're getting the VarInfo during traversal.
//
// Java: RemoveUnusedCode.traverseVar(Var)
func (r *RemoveUnusedCode) TraverseVar(var_ *Var) RemoveUnusedCode_VarInfo {
	if var_ == nil {
		panic(jutil.NewNullPointerException())
	}
	if var_ == r.lastVar {
		return r.lastVarInfo
	}
	if slot := &r.varPtrCache[rucVarCacheIdx(var_)]; slot.v == var_ {
		r.lastVar = var_
		r.lastVarInfo = slot.info
		return slot.info
	}
	if r.removeLocalVars && var_.IsArguments() {
		// If we are considering removing local variables, that includes parameters.
		// If `arguments` is used in a function we must consider all parameters to be referenced.
		varScope, _ := var_.GetScope().(*Scope)
		functionScope := varScope.GetClosestHoistScope()
		paramList := NodeUtil_GetFunctionParameters(functionScope.GetRootNode())
		for param := paramList.GetFirstChild(); param != nil; param = param.GetNext() {
			lValue := RemoveUnusedCode_NameOfParam(param)
			if lValue == nil {
				continue
			}

			r.GetVarInfo(r.GetVarForNameNode(lValue, functionScope)).
				SetIsExplicitlyNotRemovable(func() string { return "parameter in function using arguments" })
		}
		// `arguments` is never removable.
		return r.canonicalUnremovableVarInfo
	} else {
		return r.GetVarInfo(var_)
	}
}

// RemoveUnusedCode_NameOfParam returns the NAME node associated with a
// function parameter (the child of a PARAM_LIST), or null if there is no
// single name.
//
// Java: RemoveUnusedCode.nameOfParam(Node)
func RemoveUnusedCode_NameOfParam(param *rhino.Node) *rhino.Node {
	switch param.GetToken() {
	case rhino.Token_NAME:
		return param
	case rhino.Token_DEFAULT_VALUE:
		return RemoveUnusedCode_NameOfParam(param.GetFirstChild())
	case rhino.Token_ITER_REST:
		return RemoveUnusedCode_NameOfParam(param.GetOnlyChild())
	case rhino.Token_ARRAY_PATTERN, rhino.Token_OBJECT_PATTERN:
		return nil
	default:
		panic(jutil.NewIllegalStateExceptionMsg(
			"Unexpected child of PARAM_LIST: " + param.ToStringTree()))
	}
}

func hashVar(v *Var) uint32 {
	h := uint32(v.HashCode())
	h ^= h >> 16
	h *= 0x85ebca6b
	h ^= h >> 13
	return h
}

func (r *RemoveUnusedCode) putVarInfo(v *RemoveUnusedCode_RealVarInfo) {
	r.varInfos = append(r.varInfos, v)
	if len(r.varInfos)*4 > len(r.varBuckets)*3 && len(r.varBuckets) < 262144 {
		newCap := len(r.varBuckets) * 2
		newBuckets := make([]*RemoveUnusedCode_RealVarInfo, newCap)
		mask := uint32(newCap - 1)
		for _, info := range r.varInfos {
			idx := info.hash & mask
			info.hashNext = newBuckets[idx]
			newBuckets[idx] = info
		}
		r.varBuckets = newBuckets
		return
	}
	idx := v.hash & uint32(len(r.varBuckets)-1)
	v.hashNext = r.varBuckets[idx]
	r.varBuckets[idx] = v
}

// GetVarInfo gets the right VarInfo object to use for the given Var.
//
// This method is responsible for managing the entries in varInfoMap.
//
// Note: Several Vars may share the same VarInfo when they should be treated
// the same way.
//
// Java: RemoveUnusedCode.getVarInfo(Var)
func (r *RemoveUnusedCode) GetVarInfo(var_ *Var) RemoveUnusedCode_VarInfo {
	if var_ == nil {
		panic(jutil.NewNullPointerException())
	}
	if var_.IsExtern() {
		if r.unremovableLog.IsLogging() {
			name := var_.GetName()
			r.unremovableLog.LogSupplier(func() string { return name + ": extern" })
		}
		return r.canonicalUnremovableVarInfo
	}
	isGlobal := var_.IsGlobal()
	name := var_.GetName()
	// Coding conventions are extensible and may classify each lookup differently.
	// Preserve their callback before arguments handling or private cache hits.
	if r.codingConvention.IsExportedWithLocal(name /* local= */, !isGlobal) {
		if r.unremovableLog.IsLogging() {
			r.unremovableLog.LogSupplier(func() string { return name + ": exported by convention" })
		}
		return r.canonicalUnremovableVarInfo
	}
	if var_.IsArguments() {
		// No point in logging that we cannot remove "arguments"
		return r.canonicalUnremovableVarInfo
	}
	if var_ == r.lastVar {
		return r.lastVarInfo
	}
	cacheIdx := rucVarCacheIdx(var_)
	if slot := &r.varPtrCache[cacheIdx]; slot.v == var_ {
		r.lastVar = var_
		r.lastVarInfo = slot.info
		return slot.info
	}
	scopeRoot := var_.GetScopeRoot()
	h := hashVar(var_)
	if len(r.varBuckets) == 0 {
		r.varBuckets = make([]*RemoveUnusedCode_RealVarInfo, 16384)
	}
	mask := uint32(len(r.varBuckets) - 1)
	for cur := r.varBuckets[h&mask]; cur != nil; cur = cur.hashNext {
		if cur.var_ == var_ || (cur.hash == h && cur.scopeRoot == scopeRoot && cur.varName == name) {
			r.lastVar = var_
			r.lastVarInfo = cur
			r.varPtrCache[cacheIdx] = rucVarCacheEntry{v: var_, info: cur}
			return cur
		}
	}

	varInfo := NewRemoveUnusedCode_RealVarInfo(r, var_, scopeRoot, name, h)
	if var_.GetParentNode().IsParamList() {
		varInfo.SetHasNonLocalOrNonLiteralValue()
	}
	// Cannot use canonicalUnremovableVarInfo for the 2 non-removable cases below, because each
	// varInfo needs to track what value is assigned to it for the purpose of correctly allowing
	// or preventing removal of properties set on it.
	if !r.removeGlobals && isGlobal {
		varInfo.SetIsExplicitlyNotRemovable(func() string { return "not removing globals" })
	} else if !r.removeLocalVars && !isGlobal {
		varInfo.SetIsExplicitlyNotRemovable(func() string { return "not removing locals" })
	}
	r.putVarInfo(varInfo)
	r.lastVar = var_
	r.lastVarInfo = varInfo
	r.varPtrCache[cacheIdx] = rucVarCacheEntry{v: var_, info: varInfo}
	return varInfo
}

// RemoveUnreferencedVarsAndPolyfills removes any vars in the scope that were
// not referenced. Removes any assignments to those variables as well.
//
// Java: RemoveUnusedCode.removeUnreferencedVarsAndPolyfills()
func (r *RemoveUnusedCode) RemoveUnreferencedVarsAndPolyfills() {
	for _, varInfo := range r.varInfos {
		if !varInfo.IsRemovable() {
			continue
		}
		var_ := varInfo.var_

		if r.removalLog.IsLogging() {
			r.removalLog.LogSupplier(RemoveUnusedCode_RemovalLogRecord_ForVar(var_).Get)
		}
		// Regardless of what happens to the original declaration,
		// we need to remove all assigns, because they may contain references
		// to other unreferenced variables.
		varInfo.RemoveAllRemovables()

		nameNode := var_.GetNameNode()
		toRemove := nameNode.GetParent()
		if toRemove == nil || RemoveUnusedCode_AlreadyRemoved(toRemove) {
			// assignedVarInfo.removeAllRemovables () already removed it
		} else if NodeUtil_IsFunctionExpression(toRemove) {
			// TODO(bradfordcsmith): Add a Removable for this case.
			if !r.preserveFunctionExpressionNames {
				fnNameNode := toRemove.GetFirstChild()
				r.compiler.ReportChangeToEnclosingScope(fnNameNode)
				fnNameNode.SetString("")
			}
		} else {
			// Removables are not created for theses cases.
			// function foo(unused1 = someSideEffectingValue, ...unused2) {}
			// removeUnreferencedFunctionArgs() is responsible for removing these.
			// TODO(bradfordcsmith): handle parameter declarations with removables
			jutil.CheckStateMsg1(
				toRemove.IsParamList() ||
					(toRemove.GetParent().IsParamList() &&
						(toRemove.IsDefaultValue() || toRemove.IsRest())),
				"unremoved code: %s",
				toRemove)
		}
	}

	// Iterator<PolyfillInfo> iter = polyfills.values().iterator(); ... iter.remove();
	//
	// Go-only: language-adapter — the removable polyfills are collected in
	// iteration order and removed from the multimap after the iteration
	// (the JavaHashMap iterators are fail-fast); removing a polyfill's code
	// does not touch the multimap, so the order of effects is unchanged.
	if r.polyfills != nil && !r.polyfills.IsEmpty() {
		var removed []RemoveUnusedCode_PolyfillInfo
		for polyfill := range r.polyfills.Values() {
			if polyfill.GetRemoveUnusedCode_PolyfillInfoBase().isRemovable {
				if r.removalLog.IsLogging() {
					r.removalLog.LogSupplier(RemoveUnusedCode_RemovalLogRecord_ForPolyfill(polyfill).Get)
				}
				polyfill.GetRemoveUnusedCode_PolyfillInfoBase().removable.Remove(r.compiler)
				removed = append(removed, polyfill)
			}
		}
		for _, polyfill := range removed {
			r.polyfills.Remove(polyfill.GetRemoveUnusedCode_PolyfillInfoBase().key, polyfill)
		}
	}
}

// RemoveUnusedCode_Continuation: our progress in a traversal can be
// expressed completely as the current node and scope. The continuation lets
// us save that information so that we can continue the traversal later.
type RemoveUnusedCode_Continuation struct {
	node  *rhino.Node
	scope *Scope
}

// Java: RemoveUnusedCode.Continuation(Node, Scope)
func NewRemoveUnusedCode_Continuation(_ *RemoveUnusedCode, node *rhino.Node, scope *Scope) RemoveUnusedCode_Continuation {
	return RemoveUnusedCode_Continuation{node: node, scope: scope}
}

// Java: RemoveUnusedCode.Continuation.apply()
func (c RemoveUnusedCode_Continuation) Apply(outer *RemoveUnusedCode) {
	if c.node.IsFunction() {
		// Calling traverseNode here would create infinite recursion for a function declaration
		outer.TraverseFunction(c.node, c.scope)
	} else {
		outer.TraverseNode(c.node, c.scope)
	}
}

// RemoveUnusedCode_Removable represents a portion of the AST that can be
// removed.
type RemoveUnusedCode_Removable interface {
	// Java: RemoveUnusedCode.Removable.getPropertyName()
	GetPropertyName() string

	// RemoveInternal removes the associated nodes from the AST.
	//
	// Java: RemoveUnusedCode.Removable.removeInternal(AbstractCompiler)
	RemoveInternal(compiler AbstractCompiler)

	// Java: RemoveUnusedCode.Removable.remove(AbstractCompiler)
	Remove(compiler AbstractCompiler)

	// Java: RemoveUnusedCode.Removable.applyContinuations()
	ApplyContinuations()

	// Java: RemoveUnusedCode.Removable.isVariableAssignment()
	IsVariableAssignment() bool

	// Java: RemoveUnusedCode.Removable.isNamedProperty()
	IsNamedProperty() bool

	// Java: RemoveUnusedCode.Removable.isNamedPropertyAssignment()
	IsNamedPropertyAssignment() bool

	// Java: RemoveUnusedCode.Removable.isAssignedValueLocal()
	IsAssignedValueLocal() bool

	// Java: RemoveUnusedCode.Removable.getLocalAssignedValue()
	GetLocalAssignedValue() *rhino.Node

	// Java: RemoveUnusedCode.Removable.isPrototypeAssignment()
	IsPrototypeAssignment() bool

	// Java: RemoveUnusedCode.Removable.isPrototypeDotPropertyReference()
	IsPrototypeDotPropertyReference() bool

	// Java: RemoveUnusedCode.Removable.isClassOrPrototypeNamedProperty()
	IsClassOrPrototypeNamedProperty() bool

	// Java: RemoveUnusedCode.Removable.isPrototypeProperty()
	IsPrototypeProperty() bool

	// Java: RemoveUnusedCode.Removable.isThisDotPropertyReference()
	IsThisDotPropertyReference() bool

	// Java: RemoveUnusedCode.Removable.isObjectDefinePropertiesDefinition()
	IsObjectDefinePropertiesDefinition() bool

	// Java: RemoveUnusedCode.Removable.isStaticProperty()
	IsStaticProperty() bool

	// Java: RemoveUnusedCode.Removable.preventsRemovalOfVariableWithNonLocalValueOrPrototype()
	PreventsRemovalOfVariableWithNonLocalValueOrPrototype() bool

	GetRemoveUnusedCode_RemovableBase() *RemoveUnusedCode_RemovableBase
}

// RemoveUnusedCode_RemovableBase holds Removable's fields and base
// implementations.
type RemoveUnusedCode_RemovableBase struct {
	outer    *RemoveUnusedCode
	self     RemoveUnusedCode_Removable
	propNext *RemoveUnusedCode_RemovableBase

	cont0      RemoveUnusedCode_Continuation
	cont1      RemoveUnusedCode_Continuation
	extraConts []RemoveUnusedCode_Continuation

	// If this object represents an assignment of a value to a property. This
	// is the name of the property.
	propertyName string

	// If this object represents a variable declaration or assignment of a
	// value, this is the node representing where the value is being stored.
	// e.g. the LHS of an assignment.
	targetNode *rhino.Node

	numConts                        uint32
	hasPropertyName                 bool
	isPrototypeDotPropertyReference bool
	isThisDotPropertyReference      bool
	continuationsAreApplied         bool
	isRemoved                       bool
	inPropList                      bool
}

// InitRemoveUnusedCode_RemovableBase is the Removable(Node, RemovableBuilder)
// constructor; self is the removable being constructed.
//
// Java: RemoveUnusedCode.Removable(Node, RemovableBuilder)
func (b *RemoveUnusedCode_RemovableBase) InitRemoveUnusedCode_RemovableBase(self RemoveUnusedCode_Removable, outer *RemoveUnusedCode, targetNode *rhino.Node, builder *RemoveUnusedCode_RemovableBuilder) {
	b.self = self
	b.outer = outer
	b.continuationsAreApplied = false
	b.isRemoved = false
	b.numConts = builder.numConts
	b.cont0 = builder.cont0
	b.cont1 = builder.cont1
	b.extraConts = builder.extraConts
	b.propertyName = builder.propertyName
	b.hasPropertyName = builder.hasPropertyName
	b.isPrototypeDotPropertyReference = builder.isPrototypeDotPropertyReference
	b.isThisDotPropertyReference = builder.isThisDotPropertyReference
	b.targetNode = targetNode
}

// GetRemoveUnusedCode_RemovableBase gives access to the base through the
// interface.
func (b *RemoveUnusedCode_RemovableBase) GetRemoveUnusedCode_RemovableBase() *RemoveUnusedCode_RemovableBase {
	return b
}

// Java: RemoveUnusedCode.Removable.getPropertyName()
func (b *RemoveUnusedCode_RemovableBase) GetPropertyName() string {
	if !b.hasPropertyName {
		panic(jutil.NewNullPointerException())
	}
	return b.propertyName
}

// Remove removes the associated nodes from the AST, unless they've already
// been removed.
//
// Java: RemoveUnusedCode.Removable.remove(AbstractCompiler)
func (b *RemoveUnusedCode_RemovableBase) Remove(compiler AbstractCompiler) {
	if !b.isRemoved {
		b.isRemoved = true
		b.self.RemoveInternal(compiler)
	}
}

// Java: RemoveUnusedCode.Removable.applyContinuations()
func (b *RemoveUnusedCode_RemovableBase) ApplyContinuations() {
	if !b.continuationsAreApplied {
		b.continuationsAreApplied = true
		// Enqueue the continuation for processing.
		// Don't invoke the continuation immediately, because that can lead to concurrent
		// modification of data structures.
		if b.numConts >= 1 {
			b.outer.enqueueContinuation(b.cont0)
			b.cont0 = RemoveUnusedCode_Continuation{}
		}
		if b.numConts >= 2 {
			b.outer.enqueueContinuation(b.cont1)
			b.cont1 = RemoveUnusedCode_Continuation{}
		}
		for _, c := range b.extraConts {
			b.outer.enqueueContinuation(c)
		}
		clear(b.extraConts)
		b.extraConts = nil
		b.numConts = 0
	}
}

// IsVariableAssignment: true if this object represents assignment to a
// variable.
//
// Java: RemoveUnusedCode.Removable.isVariableAssignment()
func (b *RemoveUnusedCode_RemovableBase) IsVariableAssignment() bool {
	return false
}

// IsNamedProperty: true if this object represents a named property, either
// assignment or declaration.
//
// Java: RemoveUnusedCode.Removable.isNamedProperty()
func (b *RemoveUnusedCode_RemovableBase) IsNamedProperty() bool {
	return b.hasPropertyName
}

// IsNamedPropertyAssignment: true if this object represents assignment to a
// named property.
//
// This does not include class or object literal member declarations.
//
// Java: RemoveUnusedCode.Removable.isNamedPropertyAssignment()
func (b *RemoveUnusedCode_RemovableBase) IsNamedPropertyAssignment() bool {
	return false
}

// Java: RemoveUnusedCode.Removable.isAssignedValueLocal()
func (b *RemoveUnusedCode_RemovableBase) IsAssignedValueLocal() bool {
	return false // assume non-local by default
}

// GetLocalAssignedValue returns the Node representing the local value that
// is being assigned or `null` if the value is non-local or cannot be
// determined.
//
// Java: RemoveUnusedCode.Removable.getLocalAssignedValue()
func (b *RemoveUnusedCode_RemovableBase) GetLocalAssignedValue() *rhino.Node {
	return nil
}

// IsPrototypeAssignment: is this a direct assignment to `varName.prototype`?
//
// Java: RemoveUnusedCode.Removable.isPrototypeAssignment()
func (b *RemoveUnusedCode_RemovableBase) IsPrototypeAssignment() bool {
	if !b.self.IsNamedPropertyAssignment() {
		return false
	}
	if !b.hasPropertyName {
		// propertyName.equals(...) on a null field.
		panic(jutil.NewNullPointerException())
	}
	return b.propertyName == "prototype"
}

// IsPrototypeDotPropertyReference: is this an assignment to a property on a
// prototype object?
//
// Java: RemoveUnusedCode.Removable.isPrototypeDotPropertyReference()
func (b *RemoveUnusedCode_RemovableBase) IsPrototypeDotPropertyReference() bool {
	return b.isPrototypeDotPropertyReference
}

// Java: RemoveUnusedCode.Removable.isClassOrPrototypeNamedProperty()
func (b *RemoveUnusedCode_RemovableBase) IsClassOrPrototypeNamedProperty() bool {
	return false
}

// Java: RemoveUnusedCode.Removable.isPrototypeProperty()
func (b *RemoveUnusedCode_RemovableBase) IsPrototypeProperty() bool {
	return b.self.IsPrototypeDotPropertyReference() || b.self.IsClassOrPrototypeNamedProperty()
}

// Java: RemoveUnusedCode.Removable.isThisDotPropertyReference()
func (b *RemoveUnusedCode_RemovableBase) IsThisDotPropertyReference() bool {
	return b.isThisDotPropertyReference
}

// Java: RemoveUnusedCode.Removable.isObjectDefinePropertiesDefinition()
func (b *RemoveUnusedCode_RemovableBase) IsObjectDefinePropertiesDefinition() bool {
	return false
}

// IsStaticProperty.
//
// TODO(b/134610338): Combine this method with `isPrototypeProperty`.
//
// Java: RemoveUnusedCode.Removable.isStaticProperty()
func (b *RemoveUnusedCode_RemovableBase) IsStaticProperty() bool {
	return false
}

// PreventsRemovalOfVariableWithNonLocalValueOrPrototype: would a nonlocal
// or nonliteral value prevent removal of a variable associated with this
// Removable?
//
// True if the nature of this removable is such that a variable associated
// with it must not be removed if its value or its prototype is not a local,
// literal value.
//
// e.g. When X or X.prototype is nonlocal and / or nonliteral we don't know
// whether it is safe to remove code like this.
//
//	X.propName = something; // Don't know the effect of setting X.propName
//	use(something instanceof X); // can't be certain there are no instances of X
//
// Java: RemoveUnusedCode.Removable.preventsRemovalOfVariableWithNonLocalValueOrPrototype()
func (b *RemoveUnusedCode_RemovableBase) PreventsRemovalOfVariableWithNonLocalValueOrPrototype() bool {
	return false
}

// RemoveUnusedCode_RemovableBuilder is the inner class RemovableBuilder.
type RemoveUnusedCode_RemovableBuilder struct {
	outer      *RemoveUnusedCode
	cont0      RemoveUnusedCode_Continuation
	cont1      RemoveUnusedCode_Continuation
	extraConts []RemoveUnusedCode_Continuation
	numConts   uint32

	propertyName                    string
	hasPropertyName                 bool
	isPrototypeDotPropertyReference bool
	isThisDotPropertyReference      bool
}

// Java: RemoveUnusedCode.RemovableBuilder()
func NewRemoveUnusedCode_RemovableBuilder(outer *RemoveUnusedCode) *RemoveUnusedCode_RemovableBuilder {
	return &RemoveUnusedCode_RemovableBuilder{outer: outer}
}

// Java: RemoveUnusedCode.RemovableBuilder.addContinuation(Continuation)
func (b *RemoveUnusedCode_RemovableBuilder) AddContinuation(continuation RemoveUnusedCode_Continuation) *RemoveUnusedCode_RemovableBuilder {
	switch b.numConts {
	case 0:
		b.cont0 = continuation
	case 1:
		b.cont1 = continuation
	default:
		b.extraConts = append(b.extraConts, continuation)
	}
	b.numConts++
	return b
}

// Java: RemoveUnusedCode.RemovableBuilder.setIsPrototypeDotPropertyReference(boolean)
func (b *RemoveUnusedCode_RemovableBuilder) SetIsPrototypeDotPropertyReference(value bool) *RemoveUnusedCode_RemovableBuilder {
	b.isPrototypeDotPropertyReference = value
	return b
}

// Java: RemoveUnusedCode.RemovableBuilder.setIsThisDotPropertyReference(boolean)
func (b *RemoveUnusedCode_RemovableBuilder) SetIsThisDotPropertyReference(value bool) *RemoveUnusedCode_RemovableBuilder {
	b.isThisDotPropertyReference = value
	return b
}

// setPropertyName is `this.propertyName = name`.
func (b *RemoveUnusedCode_RemovableBuilder) setPropertyName(name string) {
	b.propertyName = name
	b.hasPropertyName = true
}

// Java: RemoveUnusedCode.RemovableBuilder.buildIndirectAssign(Node, Node)
func (b *RemoveUnusedCode_RemovableBuilder) BuildIndirectAssign(root *rhino.Node, targetNode *rhino.Node) *RemoveUnusedCode_IndirectAssign {
	return NewRemoveUnusedCode_IndirectAssign(b.outer, b, root, targetNode)
}

// Java: RemoveUnusedCode.RemovableBuilder.buildPolyfill(Node)
func (b *RemoveUnusedCode_RemovableBuilder) BuildPolyfill(polyfillNode *rhino.Node) *RemoveUnusedCode_Polyfill {
	return NewRemoveUnusedCode_Polyfill(b.outer, b, polyfillNode)
}

// Java: RemoveUnusedCode.RemovableBuilder.buildClassDeclaration(Node)
func (b *RemoveUnusedCode_RemovableBuilder) BuildClassDeclaration(classNode *rhino.Node) *RemoveUnusedCode_ClassDeclaration {
	return NewRemoveUnusedCode_ClassDeclaration(b.outer, b, classNode)
}

// Java: RemoveUnusedCode.RemovableBuilder.buildNamedClassExpression(Node)
func (b *RemoveUnusedCode_RemovableBuilder) BuildNamedClassExpression(classNode *rhino.Node) *RemoveUnusedCode_NamedClassExpression {
	return NewRemoveUnusedCode_NamedClassExpression(b.outer, b, classNode)
}

// Java: RemoveUnusedCode.RemovableBuilder.buildClassOrPrototypeNamedProperty(Node)
func (b *RemoveUnusedCode_RemovableBuilder) BuildClassOrPrototypeNamedProperty(propertyNode *rhino.Node) *RemoveUnusedCode_ClassOrPrototypeNamedProperty {
	jutil.CheckArgumentMsg1(
		propertyNode.IsMemberFunctionDef() ||
			propertyNode.IsMemberFieldDef() ||
			NodeUtil_IsGetOrSetKey(propertyNode) ||
			(propertyNode.IsStringKey() && !propertyNode.IsQuotedStringKey()),
		"%s", propertyNode)
	b.setPropertyName(propertyNode.GetString())
	return NewRemoveUnusedCode_ClassOrPrototypeNamedProperty(b.outer, b, propertyNode)
}

// Java: RemoveUnusedCode.RemovableBuilder.buildObjectDefinePropertiesDefinition(Node)
func (b *RemoveUnusedCode_RemovableBuilder) BuildObjectDefinePropertiesDefinition(propertyNode *rhino.Node) *RemoveUnusedCode_ObjectDefinePropertiesDefinition {
	b.setPropertyName(propertyNode.GetString())
	return NewRemoveUnusedCode_ObjectDefinePropertiesDefinition(b.outer, b, propertyNode)
}

// Java: RemoveUnusedCode.RemovableBuilder.buildFunctionDeclaration(Node)
func (b *RemoveUnusedCode_RemovableBuilder) BuildFunctionDeclaration(functionNode *rhino.Node) *RemoveUnusedCode_FunctionDeclaration {
	return NewRemoveUnusedCode_FunctionDeclaration(b.outer, b, functionNode)
}

// Java: RemoveUnusedCode.RemovableBuilder.buildNameDeclarationStatement(Node)
func (b *RemoveUnusedCode_RemovableBuilder) BuildNameDeclarationStatement(declarationStatement *rhino.Node) *RemoveUnusedCode_NameDeclarationStatement {
	return NewRemoveUnusedCode_NameDeclarationStatement(b.outer, b, declarationStatement)
}

// Java: RemoveUnusedCode.RemovableBuilder.buildNamedPropertyAssign(Node, Node)
func (b *RemoveUnusedCode_RemovableBuilder) BuildNamedPropertyAssign(assignNode *rhino.Node, propertyNode *rhino.Node) *RemoveUnusedCode_Assign {
	return b.BuildNamedPropertyAssignWithVarInfo(assignNode, propertyNode, nil)
}

// Java: RemoveUnusedCode.RemovableBuilder.buildNamedPropertyAssign(Node, Node, VarInfo)
func (b *RemoveUnusedCode_RemovableBuilder) BuildNamedPropertyAssignWithVarInfo(assignNode *rhino.Node, propertyNode *rhino.Node, varInfo RemoveUnusedCode_VarInfo) *RemoveUnusedCode_Assign {
	b.setPropertyName(propertyNode.GetString())
	return NewRemoveUnusedCode_Assign(b.outer, b, assignNode, RemoveUnusedCode_Kind_NAMED_PROPERTY, propertyNode, varInfo)
}

// Java: RemoveUnusedCode.RemovableBuilder.buildComputedPropertyAssign(Node, Node, VarInfo)
func (b *RemoveUnusedCode_RemovableBuilder) BuildComputedPropertyAssign(assignNode *rhino.Node, propertyNode *rhino.Node, varInfo RemoveUnusedCode_VarInfo) *RemoveUnusedCode_Assign {
	return NewRemoveUnusedCode_Assign(b.outer, b, assignNode, RemoveUnusedCode_Kind_COMPUTED_PROPERTY, propertyNode, varInfo)
}

// Java: RemoveUnusedCode.RemovableBuilder.buildVariableAssign(Node, VarInfo)
func (b *RemoveUnusedCode_RemovableBuilder) BuildVariableAssign(assignNode *rhino.Node, varInfo RemoveUnusedCode_VarInfo) *RemoveUnusedCode_Assign {
	return NewRemoveUnusedCode_Assign(b.outer, b, assignNode, RemoveUnusedCode_Kind_VARIABLE /* propertyNode= */, nil, varInfo)
}

// Java: RemoveUnusedCode.RemovableBuilder.buildClassSetupCall(Node)
func (b *RemoveUnusedCode_RemovableBuilder) BuildClassSetupCall(callNode *rhino.Node) *RemoveUnusedCode_ClassSetupCall {
	return b.BuildClassSetupCallWithClassDefiningCall(callNode /* classDefiningCall= */, false)
}

// Java: RemoveUnusedCode.RemovableBuilder.buildClassSetupCall(Node, boolean)
func (b *RemoveUnusedCode_RemovableBuilder) BuildClassSetupCallWithClassDefiningCall(callNode *rhino.Node, classDefiningCall bool) *RemoveUnusedCode_ClassSetupCall {
	return NewRemoveUnusedCode_ClassSetupCall(b.outer, b, callNode, classDefiningCall)
}

// Java: RemoveUnusedCode.RemovableBuilder.buildVanillaForNameDeclaration(Node)
func (b *RemoveUnusedCode_RemovableBuilder) BuildVanillaForNameDeclaration(nameNode *rhino.Node) *RemoveUnusedCode_VanillaForNameDeclaration {
	return NewRemoveUnusedCode_VanillaForNameDeclaration(b.outer, b, nameNode)
}

// Java: RemoveUnusedCode.RemovableBuilder.buildAnonymousPrototypeNamedPropertyAssign(Node, String)
func (b *RemoveUnusedCode_RemovableBuilder) BuildAnonymousPrototypeNamedPropertyAssign(
	assignNode *rhino.Node, propertyName string) *RemoveUnusedCode_AnonymousPrototypeNamedPropertyAssign {
	b.setPropertyName(propertyName)
	return NewRemoveUnusedCode_AnonymousPrototypeNamedPropertyAssign(b.outer, b, assignNode)
}

// Java: RemoveUnusedCode.RemovableBuilder.buildIncOrDepOp(Node, Node, Node)
func (b *RemoveUnusedCode_RemovableBuilder) BuildIncOrDepOp(incOrDecOp *rhino.Node, propertyNode *rhino.Node, toPreseve *rhino.Node) *RemoveUnusedCode_IncOrDecOp {
	b.setPropertyName(propertyNode.GetString())
	return NewRemoveUnusedCode_IncOrDecOp(b.outer, b, incOrDecOp, toPreseve)
}

// Java: RemoveUnusedCode.RemovableBuilder.buildUnusedReadReference(Node, Node)
func (b *RemoveUnusedCode_RemovableBuilder) BuildUnusedReadReference(referenceNode *rhino.Node, propertyNode *rhino.Node) *RemoveUnusedCode_UnusedReadReference {
	b.setPropertyName(propertyNode.GetString())
	return NewRemoveUnusedCode_UnusedReadReference(b.outer, b, referenceNode)
}

// Java: RemoveUnusedCode.RemovableBuilder.buildInstanceofName(Node)
func (b *RemoveUnusedCode_RemovableBuilder) BuildInstanceofName(instanceofNode *rhino.Node) RemoveUnusedCode_Removable {
	return NewRemoveUnusedCode_InstanceofName(b.outer, b, instanceofNode)
}

// Java: RemoveUnusedCode.RemovableBuilder.buildWeakUsageCall(Node)
func (b *RemoveUnusedCode_RemovableBuilder) BuildWeakUsageCall(weakUsageCall *rhino.Node) RemoveUnusedCode_Removable {
	return NewRemoveUnusedCode_WeakUsageCall(b.outer, b, weakUsageCall)
}

// RemoveUnusedCode_UnusedReadReference represents a read reference whose
// value is not used.
type RemoveUnusedCode_UnusedReadReference struct {
	RemoveUnusedCode_RemovableBase
	referenceNode *rhino.Node
}

// Java: RemoveUnusedCode.UnusedReadReference(RemovableBuilder, Node)
func NewRemoveUnusedCode_UnusedReadReference(outer *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, referenceNode *rhino.Node) *RemoveUnusedCode_UnusedReadReference {
	u := outer.unusedReadRefSlab.New()
	u.InitRemoveUnusedCode_RemovableBase(u, outer /* targetNode= */, nil, builder)
	// TODO(bradfordcsmith): handle `name;` and `name.property;` references
	jutil.CheckStateMsg1(
		RemoveUnusedCode_IsThisDotProperty(referenceNode) || RemoveUnusedCode_IsDotPrototypeDotProperty(referenceNode),
		"%s", referenceNode)
	u.referenceNode = referenceNode
	return u
}

// Java: RemoveUnusedCode.UnusedReadReference.removeInternal(AbstractCompiler)
func (u *RemoveUnusedCode_UnusedReadReference) RemoveInternal(compiler AbstractCompiler) {
	if !RemoveUnusedCode_AlreadyRemoved(u.referenceNode) {
		if RemoveUnusedCode_IsThisDotProperty(u.referenceNode) {
			u.outer.RemoveExpressionCompletely(u.referenceNode)
		} else {
			jutil.CheckStateMsg1(RemoveUnusedCode_IsDotPrototypeDotProperty(u.referenceNode), "%s", u.referenceNode)
			// objExpression.prototype.propertyName
			objExpression := u.referenceNode.GetFirstFirstChild()
			if u.outer.MayHaveSideEffects(objExpression) {
				u.outer.ReplaceNodeWith(u.referenceNode, objExpression.Detach())
			} else {
				u.outer.RemoveExpressionCompletely(u.referenceNode)
			}
		}
	}
}

// Java: RemoveUnusedCode.UnusedReadReference.toString()
func (u *RemoveUnusedCode_UnusedReadReference) String() string {
	return "UnusedReadReference:" + u.referenceNode.String()
}

// RemoveUnusedCode_InstanceofName represents `something instanceof
// varName`.
//
// If `varName` is removed, this expression can be replaced with `false` or
// `(something, false)` to preserve side effects.
type RemoveUnusedCode_InstanceofName struct {
	RemoveUnusedCode_RemovableBase
	instanceofNode *rhino.Node
}

// Java: RemoveUnusedCode.InstanceofName(RemovableBuilder, Node)
func NewRemoveUnusedCode_InstanceofName(outer *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, instanceofNode *rhino.Node) *RemoveUnusedCode_InstanceofName {
	i := outer.instanceofNameSlab.New()
	i.InitRemoveUnusedCode_RemovableBase(i, outer /* targetNode= */, nil, builder)
	jutil.CheckArgumentMsg1(instanceofNode.IsInstanceOf(), "%s", instanceofNode)
	i.instanceofNode = instanceofNode
	return i
}

// Java: RemoveUnusedCode.InstanceofName.removeInternal(AbstractCompiler)
func (i *RemoveUnusedCode_InstanceofName) RemoveInternal(compiler AbstractCompiler) {
	if !RemoveUnusedCode_AlreadyRemoved(i.instanceofNode) {
		lhs := i.instanceofNode.GetFirstChild()
		falseNode := rhino.IR_FalseNode().Srcref(i.instanceofNode)
		if i.outer.MayHaveSideEffects(lhs) {
			i.outer.ReplaceNodeWith(i.instanceofNode, rhino.IR_Comma(lhs.Detach(), falseNode).Srcref(i.instanceofNode))
		} else {
			i.outer.ReplaceNodeWith(i.instanceofNode, falseNode)
		}
	}
}

// Java: RemoveUnusedCode.InstanceofName.preventsRemovalOfVariableWithNonLocalValueOrPrototype()
func (i *RemoveUnusedCode_InstanceofName) PreventsRemovalOfVariableWithNonLocalValueOrPrototype() bool {
	// If we aren't sure where X comes from and what aliases it might have, we cannot be sure
	// there are no instances of it.
	return true
}

// Java: RemoveUnusedCode.InstanceofName.toString()
func (i *RemoveUnusedCode_InstanceofName) String() string {
	return "InstanceofName:" + i.instanceofNode.String()
}

// RemoveUnusedCode_WeakUsageCall represents `goog.weakUsage(varName)`.
//
// If `varName` is removed, this expression can be replaced with
// `undefined`.
type RemoveUnusedCode_WeakUsageCall struct {
	RemoveUnusedCode_RemovableBase
	weakCallNode *rhino.Node
}

// Java: RemoveUnusedCode.WeakUsageCall(RemovableBuilder, Node)
func NewRemoveUnusedCode_WeakUsageCall(outer *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, weakCallNode *rhino.Node) *RemoveUnusedCode_WeakUsageCall {
	w := &RemoveUnusedCode_WeakUsageCall{}
	w.InitRemoveUnusedCode_RemovableBase(w, outer /* targetNode= */, nil, builder)
	jutil.CheckArgumentMsg1(weakCallNode.IsCall(), "%s", weakCallNode)
	w.weakCallNode = weakCallNode
	return w
}

// Java: RemoveUnusedCode.WeakUsageCall.removeInternal(AbstractCompiler)
func (w *RemoveUnusedCode_WeakUsageCall) RemoveInternal(compiler AbstractCompiler) {
	if !RemoveUnusedCode_AlreadyRemoved(w.weakCallNode) {
		undefinedNode := NodeUtil_NewUndefinedNode(w.weakCallNode)
		w.outer.ReplaceNodeWith(w.weakCallNode, undefinedNode)
	}
}

// Java: RemoveUnusedCode.WeakUsageCall.toString()
func (w *RemoveUnusedCode_WeakUsageCall) String() string {
	return "weakUsageCall:" + w.weakCallNode.String()
}

// RemoveUnusedCode_IncOrDecOp represents an increment or decrement
// operation that could be removed.
type RemoveUnusedCode_IncOrDecOp struct {
	RemoveUnusedCode_RemovableBase
	incOrDecNode *rhino.Node
	toPreserve   *rhino.Node
}

// Java: RemoveUnusedCode.IncOrDecOp(RemovableBuilder, Node, Node)
func NewRemoveUnusedCode_IncOrDecOp(outer *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, incOrDecNode *rhino.Node, toPreserve *rhino.Node) *RemoveUnusedCode_IncOrDecOp {
	o := &RemoveUnusedCode_IncOrDecOp{}
	o.InitRemoveUnusedCode_RemovableBase(o, outer, incOrDecNode.GetOnlyChild(), builder)
	jutil.CheckArgumentMsg1(incOrDecNode.IsInc() || incOrDecNode.IsDec(), "%s", incOrDecNode)

	arg := incOrDecNode.GetOnlyChild()
	// TODO(bradfordcsmith): handle `name;` and `name.property;` references
	jutil.CheckStateMsg1(RemoveUnusedCode_IsThisDotProperty(arg) || RemoveUnusedCode_IsDotPrototypeDotProperty(arg), "%s", arg)

	o.incOrDecNode = incOrDecNode
	o.toPreserve = toPreserve
	return o
}

// Java: RemoveUnusedCode.IncOrDecOp.removeInternal(AbstractCompiler)
func (o *RemoveUnusedCode_IncOrDecOp) RemoveInternal(compiler AbstractCompiler) {
	if RemoveUnusedCode_AlreadyRemoved(o.incOrDecNode) {
		return
	}

	arg := o.incOrDecNode.GetOnlyChild()
	jutil.CheckStateMsg1(arg.IsGetProp(), "%s", arg)

	if o.toPreserve == nil {
		o.outer.RemoveExpressionCompletely(o.incOrDecNode)
	} else {
		o.outer.ReplaceNodeWith(o.incOrDecNode, o.toPreserve.Detach())
	}
}

// Java: RemoveUnusedCode.IncOrDecOp.toString()
func (o *RemoveUnusedCode_IncOrDecOp) String() string {
	return "IncOrDecOp:" + o.incOrDecNode.String()
}

// RemoveUnusedCode_IsThisDotProperty: true for `this.propertyName`.
//
// Java: RemoveUnusedCode.isThisDotProperty(Node)
func RemoveUnusedCode_IsThisDotProperty(n *rhino.Node) bool {
	return NodeUtil_IsNormalOrOptChainGetProp(n) && n.GetFirstChild().IsThis()
}

// RemoveUnusedCode_IsDotPrototypeDotProperty: true for
// `(something).prototype.propertyName`.
//
// Java: RemoveUnusedCode.isDotPrototypeDotProperty(Node)
func RemoveUnusedCode_IsDotPrototypeDotProperty(n *rhino.Node) bool {
	return NodeUtil_IsNormalOrOptChainGetProp(n) && RemoveUnusedCode_IsDotPrototype(n.GetFirstChild())
}

// RemoveUnusedCode_IndirectAssign is the inner class IndirectAssign.
type RemoveUnusedCode_IndirectAssign struct {
	RemoveUnusedCode_RemovableBase
	// The subtree which can be removed if the assignment is removable.
	root *rhino.Node
}

// Java: RemoveUnusedCode.IndirectAssign(RemovableBuilder, Node, Node)
func NewRemoveUnusedCode_IndirectAssign(outer *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, root *rhino.Node, targetNode *rhino.Node) *RemoveUnusedCode_IndirectAssign {
	a := outer.indirectAssignSlab.New()
	a.InitRemoveUnusedCode_RemovableBase(a, outer, targetNode, builder)

	rootParent := root.GetParent()
	jutil.CheckStateMsg1(rootParent.IsDestructuringPattern() || rootParent.IsParamList(), "%s", rootParent)
	jutil.CheckStateMsg1(targetNode.IsName() || targetNode.IsGetProp(), "%s", targetNode)

	a.root = root
	return a
}

// Java: RemoveUnusedCode.IndirectAssign.isVariableAssignment()
func (a *RemoveUnusedCode_IndirectAssign) IsVariableAssignment() bool {
	return a.targetNode.IsName()
}

// Java: RemoveUnusedCode.IndirectAssign.isThisDotPropertyReference()
func (a *RemoveUnusedCode_IndirectAssign) IsThisDotPropertyReference() bool {
	return RemoveUnusedCode_IsThisDotProperty(a.targetNode)
}

// Java: RemoveUnusedCode.IndirectAssign.isNamedProperty()
func (a *RemoveUnusedCode_IndirectAssign) IsNamedProperty() bool {
	return a.targetNode.IsGetProp()
}

// Java: RemoveUnusedCode.IndirectAssign.preventsRemovalOfVariableWithNonLocalValueOrPrototype()
func (a *RemoveUnusedCode_IndirectAssign) PreventsRemovalOfVariableWithNonLocalValueOrPrototype() bool {
	if a.targetNode.IsGetProp() {
		getPropLhs := a.targetNode.GetFirstChild()
		// assignment to varName.property or varName.prototype.property
		// cannot be removed unless varName and varName.prototype have literal, local values.
		return getPropLhs.IsName() || a.outer.IsNameDotPrototype(getPropLhs)
	} else {
		return false
	}
}

// Java: RemoveUnusedCode.IndirectAssign.isNamedPropertyAssignment()
func (a *RemoveUnusedCode_IndirectAssign) IsNamedPropertyAssignment() bool {
	return a.targetNode.IsGetProp()
}

// Java: RemoveUnusedCode.IndirectAssign.getPropertyName()
func (a *RemoveUnusedCode_IndirectAssign) GetPropertyName() string {
	jutil.CheckStateMsg1(a.targetNode.IsGetProp(), "%s", a.targetNode)
	return a.targetNode.GetString()
}

// Java: RemoveUnusedCode.IndirectAssign.removeInternal(AbstractCompiler)
func (a *RemoveUnusedCode_IndirectAssign) RemoveInternal(compiler AbstractCompiler) {
	if !RemoveUnusedCode_AlreadyRemoved(a.targetNode) {
		a.RemoveRoot()
	}
}

// Java: RemoveUnusedCode.IndirectAssign.removeRoot()
func (a *RemoveUnusedCode_IndirectAssign) RemoveRoot() {
	rootParent := a.root.GetParent()

	switch rootParent.GetToken() {
	case rhino.Token_ARRAY_PATTERN:
		// [a, root, b] = something;
		// [a, root] = something;
		// Replace root with an empty node to avoid messing up the order of patterns,
		// then clean up trailing empties.
		a.outer.ReplaceNodeWith(a.root, rhino.IR_Empty().Srcref(a.root))
		// We prefer `[a, b]` to `[a, b, , , , ]`
		// So remove any trailing empty nodes.
		for maybeEmpty := rootParent.GetLastChild(); maybeEmpty != nil && maybeEmpty.IsEmpty(); maybeEmpty = rootParent.GetLastChild() {
			maybeEmpty.Detach()
		}
		a.outer.compiler.ReportChangeToEnclosingScope(rootParent)
		// TODO(bradfordcsmith): If the array pattern is now empty, try to remove it entirely.
	case rhino.Token_PARAM_LIST:
		if !a.root.IsDefaultValue() {
			// removeUnreferencedFunctionArgs() is responsible for removal of function parameter
			// positions, so all we can do here is remove the default value.
			// NOTE: traverseRest() avoids creating a removable for a rest parameter.
			// TODO(bradfordcsmith): Handle parameter removal consistently with other removals.
			return
		}

		// function(removableName = removableValue)
		a.outer.compiler.ReportChangeToEnclosingScope(rootParent)
		// preserve the slot in the parameter list
		name := a.root.GetFirstChild()
		jutil.CheckState(name.IsName())
		if a.root == rootParent.GetLastChild() &&
			a.outer.removeGlobals &&
			a.outer.CanRemoveParameters(rootParent) {
			// function(p1, removableName = removableDefault)
			// and we're allowed to remove the parameter entirely
			a.root.Detach()
		} else {
			// function(removableName = removableDefault, otherParam)
			// or removableName is at the end, but cannot be completely removed.
			a.root.ReplaceWith(name.Detach())
		}
		NodeUtil_MarkFunctionsDeleted(a.root, a.outer.compiler)
	case rhino.Token_OBJECT_PATTERN:
		// ({ [propExpression]: root } = something)
		// becomes
		// ({} = something)
		NodeUtil_DeleteNode(a.root, a.outer.compiler)
	default:
		panic(jutil.NewIllegalStateExceptionMsg(
			"Unexpected parent of indirect assignment: " + rootParent.ToStringTree()))
	}
}

// RemoveUnusedCode_Polyfill is a call to $jscomp.polyfill or $jscomp.patch
// that can be removed if it is no longer referenced.
type RemoveUnusedCode_Polyfill struct {
	RemoveUnusedCode_RemovableBase
	polyfillNode *rhino.Node
	// Whether this is a call to $jscomp.patch
	isPatch bool
}

// Java: RemoveUnusedCode.Polyfill(RemovableBuilder, Node)
func NewRemoveUnusedCode_Polyfill(outer *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, polyfillNode *rhino.Node) *RemoveUnusedCode_Polyfill {
	p := &RemoveUnusedCode_Polyfill{}
	p.InitRemoveUnusedCode_RemovableBase(p, outer /* targetNode= */, nil, builder)
	p.polyfillNode = polyfillNode
	p.isPatch = strings.Contains(polyfillNode.GetFirstFirstChild().GetString(), "patch")
	return p
}

// Java: RemoveUnusedCode.Polyfill.removeInternal(AbstractCompiler)
func (p *RemoveUnusedCode_Polyfill) RemoveInternal(compiler AbstractCompiler) {
	NodeUtil_DeleteNode(p.polyfillNode, compiler)
}

// Java: RemoveUnusedCode.Polyfill.toString()
func (p *RemoveUnusedCode_Polyfill) String() string {
	if p.isPatch {
		return "Patch:" + p.polyfillNode.String()
	}
	return "Polyfill:" + p.polyfillNode.String()
}

// RemoveUnusedCode_ClassDeclaration is the inner class ClassDeclaration.
type RemoveUnusedCode_ClassDeclaration struct {
	RemoveUnusedCode_RemovableBase
	classDeclarationNode *rhino.Node
}

// Java: RemoveUnusedCode.ClassDeclaration(RemovableBuilder, Node)
func NewRemoveUnusedCode_ClassDeclaration(outer *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, classDeclarationNode *rhino.Node) *RemoveUnusedCode_ClassDeclaration {
	c := outer.classDeclSlab.New()
	// First child of the CLASS is the NAME node for the name to which the class is being
	// assigned.
	c.InitRemoveUnusedCode_RemovableBase(c, outer, classDeclarationNode.GetFirstChild(), builder)
	c.classDeclarationNode = classDeclarationNode
	return c
}

// Java: RemoveUnusedCode.ClassDeclaration.removeInternal(AbstractCompiler)
func (c *RemoveUnusedCode_ClassDeclaration) RemoveInternal(compiler AbstractCompiler) {
	NodeUtil_DeleteNode(c.classDeclarationNode, compiler)
}

// Java: RemoveUnusedCode.ClassDeclaration.isVariableAssignment()
func (c *RemoveUnusedCode_ClassDeclaration) IsVariableAssignment() bool {
	return true
}

// Java: RemoveUnusedCode.ClassDeclaration.isAssignedValueLocal()
func (c *RemoveUnusedCode_ClassDeclaration) IsAssignedValueLocal() bool {
	return true
}

// Java: RemoveUnusedCode.ClassDeclaration.getLocalAssignedValue()
func (c *RemoveUnusedCode_ClassDeclaration) GetLocalAssignedValue() *rhino.Node {
	return c.classDeclarationNode
}

// Java: RemoveUnusedCode.ClassDeclaration.toString()
func (c *RemoveUnusedCode_ClassDeclaration) String() string {
	return "ClassDeclaration:" + c.classDeclarationNode.String()
}

// RemoveUnusedCode_NamedClassExpression is the inner class
// NamedClassExpression.
type RemoveUnusedCode_NamedClassExpression struct {
	RemoveUnusedCode_RemovableBase
	classNode *rhino.Node
}

// Java: RemoveUnusedCode.NamedClassExpression(RemovableBuilder, Node)
func NewRemoveUnusedCode_NamedClassExpression(outer *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, classNode *rhino.Node) *RemoveUnusedCode_NamedClassExpression {
	c := &RemoveUnusedCode_NamedClassExpression{}
	c.InitRemoveUnusedCode_RemovableBase(c, outer, classNode.GetFirstChild(), builder)
	c.classNode = classNode
	return c
}

// Java: RemoveUnusedCode.NamedClassExpression.removeInternal(AbstractCompiler)
func (c *RemoveUnusedCode_NamedClassExpression) RemoveInternal(compiler AbstractCompiler) {
	if !RemoveUnusedCode_AlreadyRemoved(c.classNode) {
		nameNode := c.classNode.GetFirstChild()
		if !nameNode.IsEmpty() {
			// Just empty the class's name. If the expression is assigned to an unused variable,
			// then the whole class might still be removed as part of that assignment.
			nameNode.ReplaceWith(rhino.IR_Empty().Srcref(nameNode))
			compiler.ReportChangeToEnclosingScope(c.classNode)
		}
	}
}

// Java: RemoveUnusedCode.NamedClassExpression.toString()
func (c *RemoveUnusedCode_NamedClassExpression) String() string {
	return "NamedClassExpression:" + c.classNode.String()
}

// RemoveUnusedCode_ClassOrPrototypeNamedProperty is the inner class
// ClassOrPrototypeNamedProperty.
type RemoveUnusedCode_ClassOrPrototypeNamedProperty struct {
	RemoveUnusedCode_RemovableBase
	propertyNode *rhino.Node
}

// Java: RemoveUnusedCode.ClassOrPrototypeNamedProperty(RemovableBuilder, Node)
func NewRemoveUnusedCode_ClassOrPrototypeNamedProperty(outer *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, propertyNode *rhino.Node) *RemoveUnusedCode_ClassOrPrototypeNamedProperty {
	p := outer.classOrProtoPropSlab.New()
	p.InitRemoveUnusedCode_RemovableBase(p, outer /* targetNode= */, nil, builder)
	p.propertyNode = propertyNode
	return p
}

// Java: RemoveUnusedCode.ClassOrPrototypeNamedProperty.isStaticProperty()
func (p *RemoveUnusedCode_ClassOrPrototypeNamedProperty) IsStaticProperty() bool {
	return p.propertyNode.IsStaticMember()
}

// Java: RemoveUnusedCode.ClassOrPrototypeNamedProperty.isClassOrPrototypeNamedProperty()
func (p *RemoveUnusedCode_ClassOrPrototypeNamedProperty) IsClassOrPrototypeNamedProperty() bool {
	return !p.IsStaticProperty()
}

// Java: RemoveUnusedCode.ClassOrPrototypeNamedProperty.removeInternal(AbstractCompiler)
func (p *RemoveUnusedCode_ClassOrPrototypeNamedProperty) RemoveInternal(compiler AbstractCompiler) {
	NodeUtil_DeleteNode(p.propertyNode, compiler)
}

// Java: RemoveUnusedCode.ClassOrPrototypeNamedProperty.toString()
func (p *RemoveUnusedCode_ClassOrPrototypeNamedProperty) String() string {
	return "ClassOrPrototypeNamedProperty:" + p.propertyNode.String()
}

// RemoveUnusedCode_ObjectDefinePropertiesDefinition represents a single
// property definition in the object literal passed as the second argument
// to e.g. `Object.defineProperties(obj, {p1: {value: 1}, p2: {value: 3}});`.
type RemoveUnusedCode_ObjectDefinePropertiesDefinition struct {
	RemoveUnusedCode_RemovableBase
	propertyNode *rhino.Node
}

// Java: RemoveUnusedCode.ObjectDefinePropertiesDefinition(RemovableBuilder, Node)
func NewRemoveUnusedCode_ObjectDefinePropertiesDefinition(outer *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, propertyNode *rhino.Node) *RemoveUnusedCode_ObjectDefinePropertiesDefinition {
	d := &RemoveUnusedCode_ObjectDefinePropertiesDefinition{}
	d.InitRemoveUnusedCode_RemovableBase(d, outer /* targetNode= */, nil, builder)
	d.propertyNode = propertyNode
	return d
}

// Java: RemoveUnusedCode.ObjectDefinePropertiesDefinition.isObjectDefinePropertiesDefinition()
func (d *RemoveUnusedCode_ObjectDefinePropertiesDefinition) IsObjectDefinePropertiesDefinition() bool {
	return true
}

// Java: RemoveUnusedCode.ObjectDefinePropertiesDefinition.removeInternal(AbstractCompiler)
func (d *RemoveUnusedCode_ObjectDefinePropertiesDefinition) RemoveInternal(compiler AbstractCompiler) {
	NodeUtil_DeleteNode(d.propertyNode, compiler)
}

// RemoveUnusedCode_FunctionDeclaration is the inner class
// FunctionDeclaration.
type RemoveUnusedCode_FunctionDeclaration struct {
	RemoveUnusedCode_RemovableBase
	functionDeclarationNode *rhino.Node
}

// Java: RemoveUnusedCode.FunctionDeclaration(RemovableBuilder, Node)
func NewRemoveUnusedCode_FunctionDeclaration(outer *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, functionDeclarationNode *rhino.Node) *RemoveUnusedCode_FunctionDeclaration {
	f := outer.funcDeclSlab.New()
	f.InitRemoveUnusedCode_RemovableBase(f, outer, functionDeclarationNode.GetFirstChild(), builder)
	f.functionDeclarationNode = functionDeclarationNode
	return f
}

// Java: RemoveUnusedCode.FunctionDeclaration.removeInternal(AbstractCompiler)
func (f *RemoveUnusedCode_FunctionDeclaration) RemoveInternal(compiler AbstractCompiler) {
	NodeUtil_DeleteNode(f.functionDeclarationNode, compiler)
}

// Java: RemoveUnusedCode.FunctionDeclaration.isVariableAssignment()
func (f *RemoveUnusedCode_FunctionDeclaration) IsVariableAssignment() bool {
	return true
}

// Java: RemoveUnusedCode.FunctionDeclaration.isAssignedValueLocal()
func (f *RemoveUnusedCode_FunctionDeclaration) IsAssignedValueLocal() bool {
	// The declared function is always created locally.
	return true
}

// Java: RemoveUnusedCode.FunctionDeclaration.getLocalAssignedValue()
func (f *RemoveUnusedCode_FunctionDeclaration) GetLocalAssignedValue() *rhino.Node {
	return f.functionDeclarationNode
}

// Java: RemoveUnusedCode.FunctionDeclaration.toString()
func (f *RemoveUnusedCode_FunctionDeclaration) String() string {
	return "FunctionDeclaration:" + f.functionDeclarationNode.String()
}

// RemoveUnusedCode_NameDeclarationStatement is the inner class
// NameDeclarationStatement.
type RemoveUnusedCode_NameDeclarationStatement struct {
	RemoveUnusedCode_RemovableBase
	declarationStatement *rhino.Node
}

// Java: RemoveUnusedCode.NameDeclarationStatement(RemovableBuilder, Node)
func NewRemoveUnusedCode_NameDeclarationStatement(outer *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, declarationStatement *rhino.Node) *RemoveUnusedCode_NameDeclarationStatement {
	s := outer.nameDeclStmtSlab.New()
	s.InitRemoveUnusedCode_RemovableBase(s, outer, declarationStatement.GetOnlyChild(), builder)
	jutil.CheckArgumentMsg1(NodeUtil_IsNameDeclaration(declarationStatement), "%s", declarationStatement)
	s.declarationStatement = declarationStatement
	return s
}

// Java: RemoveUnusedCode.NameDeclarationStatement.removeInternal(AbstractCompiler)
func (s *RemoveUnusedCode_NameDeclarationStatement) RemoveInternal(compiler AbstractCompiler) {
	nameNode := s.declarationStatement.GetOnlyChild()
	valueNode := nameNode.GetFirstChild()
	if valueNode != nil && s.outer.MayHaveSideEffects(valueNode) {
		compiler.ReportChangeToEnclosingScope(s.declarationStatement)
		valueNode.Detach()
		s.declarationStatement.ReplaceWith(rhino.IR_ExprResult(valueNode).Srcref(valueNode))
	} else {
		NodeUtil_DeleteNode(s.declarationStatement, compiler)
	}
}

// Java: RemoveUnusedCode.NameDeclarationStatement.isVariableAssignment()
func (s *RemoveUnusedCode_NameDeclarationStatement) IsVariableAssignment() bool {
	return true
}

// Java: RemoveUnusedCode.NameDeclarationStatement.isAssignedValueLocal()
func (s *RemoveUnusedCode_NameDeclarationStatement) IsAssignedValueLocal() bool {
	nameNode := s.declarationStatement.GetOnlyChild()
	initialValueNode := nameNode.GetFirstChild()
	if initialValueNode == nil {
		// `var foo;`
		// the "assigned" value is undefined, which should be considered a "local" value,
		// since it is a constant.
		return true
	}
	// Handle `var name = name || defaultValue;`
	valueNode := RemoveUnusedCode_MaybeUnwrapQnameOrDefaultValueNode(nameNode, initialValueNode)
	return NodeUtil_EvaluatesToLocalValue(valueNode)
}

// Java: RemoveUnusedCode.NameDeclarationStatement.getLocalAssignedValue()
func (s *RemoveUnusedCode_NameDeclarationStatement) GetLocalAssignedValue() *rhino.Node {
	nameNode := s.declarationStatement.GetOnlyChild()
	initialValueNode := nameNode.GetFirstChild()
	if initialValueNode == nil {
		// `var foo;` has no node to represent the `undefined` value that is assigned.
		return nil
	}
	// Handle `var name = name || defaultValue;`
	valueNode := RemoveUnusedCode_MaybeUnwrapQnameOrDefaultValueNode(nameNode, initialValueNode)
	if NodeUtil_EvaluatesToLocalValue(valueNode) {
		return valueNode
	}
	return nil
}

// Java: RemoveUnusedCode.NameDeclarationStatement.toString()
func (s *RemoveUnusedCode_NameDeclarationStatement) String() string {
	return "NameDeclStmt:" + s.declarationStatement.String()
}

// RemoveUnusedCode_MaybeUnwrapQnameOrDefaultValueNode: targetNode is the
// node to which a value is being assigned; valueNode the value being
// assigned. If `valueNode` has the form `qualifiedName || defaultValue` and
// `qualifiedName` matches `targetNode`, return `defaultValue`. Otherwise
// return `valueNode`.
//
// Java: RemoveUnusedCode.maybeUnwrapQnameOrDefaultValueNode(Node, Node)
func RemoveUnusedCode_MaybeUnwrapQnameOrDefaultValueNode(targetNode *rhino.Node, valueNode *rhino.Node) *rhino.Node {
	if valueNode.IsOr() && targetNode.IsQualifiedName() {
		lhsOfOr := jutil.CheckNotNull(valueNode.GetFirstChild())
		if lhsOfOr.MatchesQualifiedNameNode(targetNode) {
			// We use `matchesQualifiedName` rather than `isEquivalentTo` to properly handle the
			// variable declaration case where the assigned value is a child of the name we want to
			// match.
			return valueNode.GetLastChild()
		}
	}
	return valueNode
}

// RemoveUnusedCode_Kind is the nested enum Kind.
type RemoveUnusedCode_Kind int32

const (
	// X = something;
	RemoveUnusedCode_Kind_VARIABLE RemoveUnusedCode_Kind = iota
	// X.propertyName = something;
	// X.prototype.propertyName = something;
	RemoveUnusedCode_Kind_NAMED_PROPERTY
	// X[expression] = something;
	// X.prototype[expression] = something;
	RemoveUnusedCode_Kind_COMPUTED_PROPERTY
)

var removeUnusedCode_Kind_names = [...]string{"VARIABLE", "NAMED_PROPERTY", "COMPUTED_PROPERTY"}

// Java: RemoveUnusedCode.Kind.values()
func RemoveUnusedCode_Kind_values() []RemoveUnusedCode_Kind {
	return []RemoveUnusedCode_Kind{
		RemoveUnusedCode_Kind_VARIABLE, RemoveUnusedCode_Kind_NAMED_PROPERTY, RemoveUnusedCode_Kind_COMPUTED_PROPERTY,
	}
}

// Java: RemoveUnusedCode.Kind.name()
func (k RemoveUnusedCode_Kind) String() string { return removeUnusedCode_Kind_names[k] }

// RemoveUnusedCode_Assign is the inner class Assign.
type RemoveUnusedCode_Assign struct {
	RemoveUnusedCode_RemovableBase

	assignNode *rhino.Node
	kind       RemoveUnusedCode_Kind

	// The VarInfo associated with the LHS of the assignment.
	//
	// The VarInfo for X in the following cases. `null` for all others.
	//
	//	X = class {};
	//	X.prop = 1;
	//	X.prototype.prop = 1;
	//	X[prop] = 1;
	varInfo RemoveUnusedCode_VarInfo
}

// Java: RemoveUnusedCode.Assign(RemovableBuilder, Node, Kind, Node, VarInfo)
func NewRemoveUnusedCode_Assign(
	outer *RemoveUnusedCode,
	builder *RemoveUnusedCode_RemovableBuilder,
	assignNode *rhino.Node,
	kind RemoveUnusedCode_Kind,
	propertyNode *rhino.Node,
	varInfo RemoveUnusedCode_VarInfo) *RemoveUnusedCode_Assign {
	a := outer.assignSlab.New()
	a.InitRemoveUnusedCode_RemovableBase(a, outer, assignNode.GetFirstChild(), builder)
	jutil.CheckArgumentMsg1(NodeUtil_IsAssignmentOp(assignNode), "%s", assignNode)
	if kind == RemoveUnusedCode_Kind_VARIABLE {
		jutil.CheckArgumentMsg1(
			propertyNode == nil,
			"got property node for simple variable assignment: %s",
			propertyNode)
		jutil.CheckArgumentMsg1(varInfo != nil, "missing VarInfo for variable assignment: %s", propertyNode)
	} else {
		jutil.CheckArgumentMsg(propertyNode != nil, "missing property node")
		if kind == RemoveUnusedCode_Kind_NAMED_PROPERTY {
			jutil.CheckArgumentMsg1(
				propertyNode.IsGetProp(), "property name is not a GETPROP: %s", propertyNode)
		}
	}
	a.assignNode = assignNode
	a.kind = kind
	a.varInfo = varInfo
	return a
}

// IsVariableAssignment: true for `varName = value` assignments.
//
// Java: RemoveUnusedCode.Assign.isVariableAssignment()
func (a *RemoveUnusedCode_Assign) IsVariableAssignment() bool {
	return a.kind == RemoveUnusedCode_Kind_VARIABLE
}

// Java: RemoveUnusedCode.Assign.isAssignedValueLocal()
func (a *RemoveUnusedCode_Assign) IsAssignedValueLocal() bool {
	return a.GetLocalAssignedValue() != nil
}

// Java: RemoveUnusedCode.Assign.getLocalAssignedValue()
func (a *RemoveUnusedCode_Assign) GetLocalAssignedValue() *rhino.Node {
	if NodeUtil_IsExpressionResultUsed(a.assignNode) {
		// assigned value may escape or be aliased
		return nil
	} else {
		// Handle `qname = qname || defaultValue;`
		valueNode :=
			RemoveUnusedCode_MaybeUnwrapQnameOrDefaultValueNode(
				a.assignNode.GetFirstChild(), a.assignNode.GetLastChild())
		if NodeUtil_EvaluatesToLocalValue(valueNode) {
			return valueNode
		} else {
			return nil
		}
	}
}

// Java: RemoveUnusedCode.Assign.preventsRemovalOfVariableWithNonLocalValueOrPrototype()
func (a *RemoveUnusedCode_Assign) PreventsRemovalOfVariableWithNonLocalValueOrPrototype() bool {
	// If we don't know where the variable comes from or where it may go, then we don't know
	// whether it is safe to remove assignments to properties on it.
	return a.IsNamedPropertyAssignment() || a.IsComputedPropertyAssignment()
}

// IsNamedPropertyAssignment: true for `varName.propName = value` and
// `varName.prototype.propName = value` assignments.
//
// Java: RemoveUnusedCode.Assign.isNamedPropertyAssignment()
func (a *RemoveUnusedCode_Assign) IsNamedPropertyAssignment() bool {
	return a.kind == RemoveUnusedCode_Kind_NAMED_PROPERTY
}

// IsComputedPropertyAssignment: true for `varName[expr] = value` and
// `varName.prototype[expr] = value` assignments.
//
// Java: RemoveUnusedCode.Assign.isComputedPropertyAssignment()
func (a *RemoveUnusedCode_Assign) IsComputedPropertyAssignment() bool {
	return a.kind == RemoveUnusedCode_Kind_COMPUTED_PROPERTY
}

// Java: RemoveUnusedCode.Assign.isStaticProperty()
func (a *RemoveUnusedCode_Assign) IsStaticProperty() bool {
	if a.kind == RemoveUnusedCode_Kind_NAMED_PROPERTY &&
		a.varInfo != nil &&
		a.varInfo.HasFunctionOrClassLiteralValue() {
		// We have either
		// `classOrFunctionVar.prop = something;` which is static
		// or
		// `classOrFunctionVar.prototype.prop = something;` which is not.
		return a.targetNode.GetFirstChild().IsName()
	} else {
		return false
	}
}

// RemoveInternal replaces the current assign with its right hand side.
//
// Java: RemoveUnusedCode.Assign.removeInternal(AbstractCompiler)
func (a *RemoveUnusedCode_Assign) RemoveInternal(compiler AbstractCompiler) {
	if RemoveUnusedCode_AlreadyRemoved(a.assignNode) {
		return
	}
	parent := a.assignNode.GetParent()
	compiler.ReportChangeToEnclosingScope(parent)
	lhs := a.assignNode.GetFirstChild()
	rhs := a.assignNode.GetSecondChild()
	mustPreserveRhs :=
		a.outer.MayHaveSideEffects(rhs) || NodeUtil_IsExpressionResultUsed(a.assignNode)
	mustPreserveGetElmExpr := lhs.IsGetElem() && a.outer.MayHaveSideEffects(lhs.GetLastChild())

	if mustPreserveRhs && mustPreserveGetElmExpr {
		replacement := rhino.IR_Comma(lhs.GetLastChild().Detach(), rhs.Detach()).Srcref(a.assignNode)
		a.outer.ReplaceNodeWith(a.assignNode, replacement)
	} else if mustPreserveGetElmExpr {
		a.outer.ReplaceNodeWith(a.assignNode, lhs.GetLastChild().Detach())
	} else if mustPreserveRhs {
		a.outer.ReplaceNodeWith(a.assignNode, rhs.Detach())
	} else {
		a.outer.RemoveExpressionCompletely(a.assignNode)
	}
}

// Java: RemoveUnusedCode.Assign.toString()
func (a *RemoveUnusedCode_Assign) String() string {
	return "Assign:" + a.assignNode.String()
}

// Java: RemoveUnusedCode.mayHaveSideEffects(Node)
func (r *RemoveUnusedCode) MayHaveSideEffects(node *rhino.Node) bool {
	// check for @pureOrBreakMyCode on call expressions.
	var jsDocInfo *rhino.JSDocInfo
	if node.IsCall() {
		jsDocInfo = node.GetJSDocInfo()
	}
	if jsDocInfo != nil && jsDocInfo.IsPureOrBreakMyCode() {
		return false
	}
	return r.astAnalyzer.MayHaveSideEffects(node)
}

// RemoveUnusedCode_AnonymousPrototypeNamedPropertyAssign represents
// `(someObjectExpression).prototype.propertyName = someValue`.
type RemoveUnusedCode_AnonymousPrototypeNamedPropertyAssign struct {
	RemoveUnusedCode_RemovableBase
	assignNode *rhino.Node
}

// Java: RemoveUnusedCode.AnonymousPrototypeNamedPropertyAssign(RemovableBuilder, Node)
func NewRemoveUnusedCode_AnonymousPrototypeNamedPropertyAssign(outer *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, assignNode *rhino.Node) *RemoveUnusedCode_AnonymousPrototypeNamedPropertyAssign {
	a := outer.anonProtoAssignSlab.New()
	a.InitRemoveUnusedCode_RemovableBase(a, outer, assignNode.GetFirstChild(), builder)
	if !builder.hasPropertyName {
		panic(jutil.NewNullPointerException())
	}
	jutil.CheckArgumentMsg1(assignNode.IsAssign(), "%s", assignNode)
	a.assignNode = assignNode
	return a
}

// Java: RemoveUnusedCode.AnonymousPrototypeNamedPropertyAssign.removeInternal(AbstractCompiler)
func (a *RemoveUnusedCode_AnonymousPrototypeNamedPropertyAssign) RemoveInternal(compiler AbstractCompiler) {
	if RemoveUnusedCode_AlreadyRemoved(a.assignNode) {
		return
	}
	parent := a.assignNode.GetParent()
	compiler.ReportChangeToEnclosingScope(parent)
	lhs := a.assignNode.GetFirstChild()
	rhs := a.assignNode.GetLastChild()

	jutil.CheckStateMsg1(lhs.IsGetProp(), "%s", lhs)
	objDotPrototype := lhs.GetFirstChild()
	jutil.CheckStateMsg1(objDotPrototype.IsGetProp(), "%s", objDotPrototype)
	objExpression := objDotPrototype.GetFirstChild()
	jutil.CheckStateMsg1(objDotPrototype.GetString() == "prototype", "%s", objDotPrototype)

	mustPreserveRhs :=
		a.outer.MayHaveSideEffects(rhs) || NodeUtil_IsExpressionResultUsed(a.assignNode)
	mustPreserveObjExpression := a.outer.MayHaveSideEffects(objExpression)

	if mustPreserveRhs && mustPreserveObjExpression {
		replacement := rhino.IR_Comma(objExpression.Detach(), rhs.Detach()).Srcref(a.assignNode)
		a.outer.ReplaceNodeWith(a.assignNode, replacement)
	} else if mustPreserveObjExpression {
		a.outer.ReplaceNodeWith(a.assignNode, objExpression.Detach())
	} else if mustPreserveRhs {
		a.outer.ReplaceNodeWith(a.assignNode, rhs.Detach())
	} else {
		a.outer.RemoveExpressionCompletely(a.assignNode)
	}
}

// Java: RemoveUnusedCode.AnonymousPrototypeNamedPropertyAssign.isPrototypeProperty()
func (a *RemoveUnusedCode_AnonymousPrototypeNamedPropertyAssign) IsPrototypeProperty() bool {
	return true
}

// Java: RemoveUnusedCode.AnonymousPrototypeNamedPropertyAssign.toString()
func (a *RemoveUnusedCode_AnonymousPrototypeNamedPropertyAssign) String() string {
	return "AnonymousPrototypeNamedPropertyAssign:" + a.assignNode.String()
}

// RemoveUnusedCode_ClassSetupCall represents a call to a class setup method
// such as `goog.inherits()` or `goog.addSingletonGetter()`.
type RemoveUnusedCode_ClassSetupCall struct {
	RemoveUnusedCode_RemovableBase
	callNode          *rhino.Node
	classDefiningCall bool
}

// Java: RemoveUnusedCode.ClassSetupCall(RemovableBuilder, Node, boolean)
func NewRemoveUnusedCode_ClassSetupCall(outer *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, callNode *rhino.Node, classDefiningCall bool) *RemoveUnusedCode_ClassSetupCall {
	c := outer.classSetupCallSlab.New()
	c.InitRemoveUnusedCode_RemovableBase(c, outer /* targetNode= */, nil, builder)
	c.callNode = callNode
	c.classDefiningCall = classDefiningCall
	return c
}

// Java: RemoveUnusedCode.ClassSetupCall.removeInternal(AbstractCompiler)
func (c *RemoveUnusedCode_ClassSetupCall) RemoveInternal(compiler AbstractCompiler) {
	parent := c.callNode.GetParent()

	var replacement *rhino.Node
	// Need to keep call args that have side effects.
	// Easiest thing to do is break apart the call node as we go.
	// First child is the callee (aka. Object.defineProperties or equivalent)
	c.callNode.RemoveFirstChild()
	for arg := c.callNode.GetLastChild(); arg != nil; arg = c.callNode.GetLastChild() {
		arg.Detach()
		// If this is a class defining call, the arguments are well defined and verified so they are
		// always safe to drop.
		if !c.classDefiningCall && c.outer.MayHaveSideEffects(arg) {
			if replacement == nil {
				replacement = arg
			} else {
				replacement = rhino.IR_Comma(arg, replacement).Srcref(c.callNode)
			}
		} else {
			NodeUtil_MarkFunctionsDeleted(arg, compiler)
		}
	}

	// This can be part of an arbitrary expression but the results must be unused.
	if replacement != nil {
		c.outer.ReplaceNodeWith(c.callNode, replacement)
	} else if parent.IsExprResult() {
		NodeUtil_DeleteNode(parent, compiler)
	} else {
		// We have been asked to remove the value inside an expression. This will only happen if
		// we know the result of this sub-expression is otherwise unused (doesn't change the result
		// of the expression when removed).
		if parent.IsComma() || parent.IsAnd() || parent.IsOr() {
			if parent.GetFirstChild() == c.callNode {
				// `(goog.inherits(A, B), something)` -> `something`
				rhs := jutil.CheckNotNull(c.callNode.GetNext())
				compiler.ReportChangeToEnclosingScope(parent)
				parent.ReplaceWith(rhs.Detach())
			} else {
				// `(something, Object.defineProperties(A, B))` -> `something`
				lhs := parent.GetFirstChild()
				compiler.ReportChangeToEnclosingScope(parent)
				parent.ReplaceWith(lhs.Detach())
			}
		} else {
			// `x ? Object.defineProperties(A, B) : something` -> `x ? 0 : something`
			// Leave simplifying arbitrary expressions to the peephole passes.
			compiler.ReportChangeToEnclosingScope(parent)
			c.callNode.ReplaceWith(rhino.IR_Number(0))
		}
	}
}

// Java: RemoveUnusedCode.ClassSetupCall.preventsRemovalOfVariableWithNonLocalValueOrPrototype()
func (c *RemoveUnusedCode_ClassSetupCall) PreventsRemovalOfVariableWithNonLocalValueOrPrototype() bool {
	// If we aren't sure where X comes from and what aliases it might have, we cannot be sure
	// it's safe to remove the class setup for it.
	return true
}

// Java: RemoveUnusedCode.ClassSetupCall.toString()
func (c *RemoveUnusedCode_ClassSetupCall) String() string {
	return "ClassSetupCall:" + c.callNode.String()
}

// Java: RemoveUnusedCode.alreadyRemoved(Node)
func RemoveUnusedCode_AlreadyRemoved(n *rhino.Node) bool {
	parent := n.GetParent()
	if parent == nil {
		return true
	}
	if parent.IsRoot() {
		return false
	}
	return RemoveUnusedCode_AlreadyRemoved(parent)
}

// RemoveUnusedCode_VarInfo tracks whether a variable is removable or not,
// including tracking the Removable objects associated with it.
type RemoveUnusedCode_VarInfo interface {
	// AddRemovable adds a Removable representing code that must be removed
	// if this variable is removed.
	//
	//   - The contents of the Removable could cause the variable to be no
	//     longer safe to remove.
	//   - If the variable is not safe to remove, this method will either
	//     apply the continuations within the `removable` or allow it to be
	//     considered for independent removal.
	//
	// Java: RemoveUnusedCode.VarInfo.addRemovable(Removable)
	AddRemovable(removable RemoveUnusedCode_Removable)

	// IsRemovable: at the current point of execution, does the variable
	// appear safe to remove?
	//
	// Java: RemoveUnusedCode.VarInfo.isRemovable()
	IsRemovable() bool

	// SetIsExplicitlyNotRemovable records that the variable cannot be
	// removed.
	//
	// If the variable was previously considered safe to remove, then this
	// method will examine all of the `Removable` objects associated with this
	// variable and either apply their continuations or consider them for
	// independent removal.
	//
	// Java: RemoveUnusedCode.VarInfo.setIsExplicitlyNotRemovable(Supplier<String>)
	SetIsExplicitlyNotRemovable(reasonSupplier func() string)

	// SetHasNonLocalOrNonLiteralValue records that at least one value
	// assigned to the variable is non-local (comes from or escapes to
	// another scope) and / or non-literal.
	//
	// Java: RemoveUnusedCode.VarInfo.setHasNonLocalOrNonLiteralValue()
	SetHasNonLocalOrNonLiteralValue()

	// HasFunctionOrClassLiteralValue: is at least one value assigned to the
	// variable a function or class literal?
	//
	// Java: RemoveUnusedCode.VarInfo.hasFunctionOrClassLiteralValue()
	HasFunctionOrClassLiteralValue() bool

	// RemoveAllRemovables invokes the `remove()` method on all removables
	// associated with this variable.
	//
	// Does nothing if the variable has been found to be unsafe to remove.
	//
	// Java: RemoveUnusedCode.VarInfo.removeAllRemovables()
	RemoveAllRemovables()
}

// RemoveUnusedCode_CanonicalUnremovableVarInfo represents a variable that we
// know can never be removed regardless of how it is used.
//
// We create just one instance of this class and use it for many variables
// in order to save memory.
type RemoveUnusedCode_CanonicalUnremovableVarInfo struct {
	outer *RemoveUnusedCode
}

// Java: RemoveUnusedCode.CanonicalUnremovableVarInfo()
func NewRemoveUnusedCode_CanonicalUnremovableVarInfo(outer *RemoveUnusedCode) *RemoveUnusedCode_CanonicalUnremovableVarInfo {
	return &RemoveUnusedCode_CanonicalUnremovableVarInfo{outer: outer}
}

// Java: RemoveUnusedCode.CanonicalUnremovableVarInfo.addRemovable(Removable)
func (v *RemoveUnusedCode_CanonicalUnremovableVarInfo) AddRemovable(removable RemoveUnusedCode_Removable) {
	// Immediately pass the argument off for potential independent removal.
	v.outer.ConsiderForIndependentRemoval(removable)
}

// Java: RemoveUnusedCode.CanonicalUnremovableVarInfo.isRemovable()
func (v *RemoveUnusedCode_CanonicalUnremovableVarInfo) IsRemovable() bool {
	return false
}

// Java: RemoveUnusedCode.CanonicalUnremovableVarInfo.setIsExplicitlyNotRemovable(Supplier<String>)
func (v *RemoveUnusedCode_CanonicalUnremovableVarInfo) SetIsExplicitlyNotRemovable(reasonSupplier func() string) {
	// nothing to do
}

// Java: RemoveUnusedCode.CanonicalUnremovableVarInfo.setHasNonLocalOrNonLiteralValue()
func (v *RemoveUnusedCode_CanonicalUnremovableVarInfo) SetHasNonLocalOrNonLiteralValue() {
	// nothing to do
}

// Java: RemoveUnusedCode.CanonicalUnremovableVarInfo.hasFunctionOrClassLiteralValue()
func (v *RemoveUnusedCode_CanonicalUnremovableVarInfo) HasFunctionOrClassLiteralValue() bool {
	// Returning false here will prevent some properties assigned on unremovable variables from
	// being independently removed, but returning `true` would cause incorrect removal of
	// properties.
	return false
}

// Java: RemoveUnusedCode.CanonicalUnremovableVarInfo.removeAllRemovables()
func (v *RemoveUnusedCode_CanonicalUnremovableVarInfo) RemoveAllRemovables() {
	// nothing to do
}

// RemoveUnusedCode_RealVarInfo tracks the removable code and other state
// related to variables we may be able to remove.
type RemoveUnusedCode_RealVarInfo struct {
	outer     *RemoveUnusedCode
	var_      *Var
	scopeRoot *rhino.Node
	hashNext  *RemoveUnusedCode_RealVarInfo
	varName   string

	// Objects that represent variable declarations, assignments, or class
	// setup calls that can be removed.
	//
	// NOTE: Once we realize that we cannot remove the variable, this list
	// will be cleared and no more will be added.
	removables    []RemoveUnusedCode_Removable
	removablesBuf [2]RemoveUnusedCode_Removable

	hash                uint32
	isEntirelyRemovable bool

	// At least one assignment to the variable is a non-local and/or
	// non-literal value.
	hasNonLocalOrNonLiteralValue bool

	// NOTE: We are assuming that if one value assigned to a variable is a
	// class or function literal, than it is very likely that all other
	// values, if any, assigned to the variable are functions or classes. At
	// present this information is used only to decide whether
	// `varName.propName = something` should be considered to be an ES5-style
	// static property. If this assumption is wrong we may end up removing
	// `propName` even though it's not actually a static class property. This
	// seems a reasonable risk, because that removal would only occur if there
	// were no references to `propName` anywhere in the sources or externs.
	//
	// At least one assignment to the variable is a function or class literal.
	hasFunctionOrClassLiteralValue      bool
	requiresLocalLiteralValueForRemoval bool
}

// Java: RemoveUnusedCode.RealVarInfo(String)
func NewRemoveUnusedCode_RealVarInfo(outer *RemoveUnusedCode, var_ *Var, scopeRoot *rhino.Node, varName string, hash uint32) *RemoveUnusedCode_RealVarInfo {
	v := outer.realVarInfoSlab.New()
	v.outer = outer
	v.var_ = var_
	v.scopeRoot = scopeRoot
	v.varName = varName
	v.hash = hash
	v.isEntirelyRemovable = true
	v.removables = v.removablesBuf[:0]
	return v
}

// Java: RemoveUnusedCode.RealVarInfo.addRemovable(Removable)
func (v *RemoveUnusedCode_RealVarInfo) AddRemovable(removable RemoveUnusedCode_Removable) {
	if removable.IsVariableAssignment() {
		// class name {}
		// function name {}
		// let name = something;
		// name = something;
		// let {a} = something;
		if localValue := removable.GetLocalAssignedValue(); localValue != nil {
			if localValue.IsFunction() || localValue.IsClass() {
				v.hasFunctionOrClassLiteralValue = true
			}
		} else if !removable.IsAssignedValueLocal() {
			// Still have to check isAssignedValueLocal() when localValue is nil because of variable
			// declarations without initial values:
			// `var a;` isAssignedValueLocal() == true but getLocalAssignedValue() == null
			v.hasNonLocalOrNonLiteralValue = true
		}
	} else if removable.IsPrototypeAssignment() && !removable.IsAssignedValueLocal() {
		// `name.prototype = someNonLocalValue;`
		v.hasNonLocalOrNonLiteralValue = true
	}
	if removable.PreventsRemovalOfVariableWithNonLocalValueOrPrototype() {
		v.requiresLocalLiteralValueForRemoval = true
	}
	if v.hasNonLocalOrNonLiteralValue && v.requiresLocalLiteralValueForRemoval {
		v.SetIsExplicitlyNotRemovable(
			func() string { return "hasNonLocalOrNonLiteralValue && requiresLocalLiteralValueForRemoval" })
	}

	if v.isEntirelyRemovable {
		// Store for possible removal later.
		v.removables = append(v.removables, removable)
	} else {
		v.outer.ConsiderForIndependentRemoval(removable)
	}
}

// Java: RemoveUnusedCode.RealVarInfo.isRemovable()
func (v *RemoveUnusedCode_RealVarInfo) IsRemovable() bool {
	return v.isEntirelyRemovable
}

// Java: RemoveUnusedCode.RealVarInfo.setIsExplicitlyNotRemovable(Supplier<String>)
func (v *RemoveUnusedCode_RealVarInfo) SetIsExplicitlyNotRemovable(reasonSupplier func() string) {
	if v.isEntirelyRemovable {
		v.isEntirelyRemovable = false
		if v.outer.unremovableLog.IsLogging() {
			v.outer.unremovableLog.LogString(v.varName + ": " + reasonSupplier())
		}
		for _, r := range v.removables {
			v.outer.ConsiderForIndependentRemoval(r)
		}
		clear(v.removables)
		v.removables = v.removables[:0]
	}
}

// Java: RemoveUnusedCode.RealVarInfo.setHasNonLocalOrNonLiteralValue()
func (v *RemoveUnusedCode_RealVarInfo) SetHasNonLocalOrNonLiteralValue() {
	v.hasNonLocalOrNonLiteralValue = true
}

// Java: RemoveUnusedCode.RealVarInfo.hasFunctionOrClassLiteralValue()
func (v *RemoveUnusedCode_RealVarInfo) HasFunctionOrClassLiteralValue() bool {
	return v.hasFunctionOrClassLiteralValue
}

// Java: RemoveUnusedCode.RealVarInfo.removeAllRemovables()
func (v *RemoveUnusedCode_RealVarInfo) RemoveAllRemovables() {
	jutil.CheckState(v.isEntirelyRemovable)
	for _, removable := range v.removables {
		removable.Remove(v.outer.compiler)
	}
	clear(v.removables)
	v.removables = v.removables[:0]
}

// CreatePolyfillInfo makes a new PolyfillInfo, including the correct
// Removable. Parses the name to determine whether this is a global, static,
// or prototype polyfill.
//
// Java: RemoveUnusedCode.createPolyfillInfo(Node, Scope, String)
func (r *RemoveUnusedCode) CreatePolyfillInfo(call *rhino.Node, scope *Scope, name string) RemoveUnusedCode_PolyfillInfo {
	jutil.CheckState(call.GetParent().IsExprResult())
	if !r.guardedUsagesComputed {
		r.guardedUsagesComputed = true
		NewPolyfillUsageFinder(r.compiler, r.polyfillsFromTable).
			TraverseOnlyGuarded(r.root, r.StorePolyfill)
	}
	// Make the removable and polyfill info.  Add continuations for all arguments.
	builder := NewRemoveUnusedCode_RemovableBuilder(r)
	for n := call.GetFirstChild().GetNext(); n != nil; n = n.GetNext() {
		builder.AddContinuation(NewRemoveUnusedCode_Continuation(r, n, scope))
	}
	removable := builder.BuildPolyfill(call.GetParent())
	lastDot := jutil.StringLastIndexOf(name, ".")
	if lastDot < 0 {
		return NewRemoveUnusedCode_GlobalPolyfillInfo(r, removable, name)
	}
	owner := jutil.StringSubstring(name, 0, lastDot)
	prop := jutil.StringSubstringFrom(name, lastDot+1)
	if strings.HasSuffix(owner, RemoveUnusedCode_DOT_PROTOTYPE) {
		owner = jutil.StringSubstring(owner, 0, jutil.StringLength(owner)-jutil.StringLength(RemoveUnusedCode_DOT_PROTOTYPE))
		return NewRemoveUnusedCode_PrototypePropertyPolyfillInfo(r, removable, prop, owner)
	}
	return NewRemoveUnusedCode_StaticPropertyPolyfillInfo(r, removable, prop, owner)
}

const RemoveUnusedCode_DOT_PROTOTYPE = ".prototype"

// RemoveUnusedCode_PolyfillInfo stores information about definitions and
// usages of polyfills.
//
// The polyfill removal strategy is as follows. First, look for all the
// polyfill definitions, whose names are stores as strings passed as the
// first argument to $jscomp.polyfill. Each definition falls into one of
// three categories: (1) global names, such as Map or Promise; (2) static
// properties, such as Array.from or Reflect.get, which must always have
// exactly two name components (TODO(bradfordcsmith): is this assumption
// required?); or (3) prototype properties, such as String.prototype.repeat
// or Promise.prototype.finally, which must always have exactly three name
// components. The definition can be removed once it is found that there are
// no references to it.
//
// References are ignored if they are "guarded". This allows removing, e.g,
// the Promise polyfill if it is only referenced in `if (typeof Promise ===
// 'function') { use(Promise); }`. Note that a guarded reference to a
// polyfill does not guarantee its removal, either. Polyfills may have
// nonguarded references as well.
//
// Determining whether a node is a reference depends on the type of
// polyfill. When type information is available, the type of the expected
// owner (i.e. the global object for global polyfills, the namespace or class
// for static polyfills, or an instance of the owning class (or its implicit
// prototype) for prototype polyfills) is used exclusively to determine this
// with very good accuracy. Types are considered to match if a direct cast
// would be allowed without a warning (i.e. some element of the union is a
// direct subtype or supertype).
//
// When type information is not available (or is too loose) then we fall
// back on a heuristic:
//
//   - globals are referenced by any same-named NAME node or any GETPROP node
//     whose last child has the same string (this allows matching
//     goog.global.Map, but will also match MyOuter.Map).
//   - static properties are referenced by any GETPROP node whose last child
//     is the same as the polyfill's property name and whose owner
//     references the polyfill owner per the above rule.
//   - prototype properties are referenced by any GETPROP node whose last
//     child is the same as the polyfill's property name, regardless of its
//     owner.
//
// Note that this results in both false positives and false negatives in
// untyped code: we may remove polyfills that are actually used (e.g. if
// Array.from is accessed via a subclass as SubArray.from or in a subclass'
// static method as this.from) and we may retain polyfills that are not used
// (e.g. if a user-defined nested class shares the same name as a global
// builtin, as in Foo.Map). For greater consistency we may shift this balance
// in the future to eliminate the possibility of incorrect removals, at the
// cost of more incorrect retentions.
type RemoveUnusedCode_PolyfillInfo interface {
	// Java: RemoveUnusedCode.PolyfillInfo.considerPossibleReference(Node)
	ConsiderPossibleReference(n *rhino.Node)

	// Java: RemoveUnusedCode.PolyfillInfo.getName()
	GetName() string

	// ConsiderPossibleReferenceInternal is the template method to check the
	// node.
	//
	// Java: RemoveUnusedCode.PolyfillInfo.considerPossibleReferenceInternal(Node)
	ConsiderPossibleReferenceInternal(n *rhino.Node)

	GetRemoveUnusedCode_PolyfillInfoBase() *RemoveUnusedCode_PolyfillInfoBase
}

// RemoveUnusedCode_PolyfillInfoBase holds PolyfillInfo's fields and base
// implementations.
type RemoveUnusedCode_PolyfillInfoBase struct {
	outer *RemoveUnusedCode
	self  RemoveUnusedCode_PolyfillInfo

	// The Polyfill instance corresponding to the polyfill's definition.
	removable *RemoveUnusedCode_Polyfill

	// The rightmost component of the polyfill's qualified name (does not
	// contain a dot).
	key string

	// Whether the polyfill is unreferenced and this can be removed safely.
	isRemovable bool
}

// InitRemoveUnusedCode_PolyfillInfoBase is the PolyfillInfo(Polyfill,
// String) constructor; self is the polyfill info being constructed.
//
// Java: RemoveUnusedCode.PolyfillInfo(Polyfill, String)
func (b *RemoveUnusedCode_PolyfillInfoBase) InitRemoveUnusedCode_PolyfillInfoBase(self RemoveUnusedCode_PolyfillInfo, outer *RemoveUnusedCode, removable *RemoveUnusedCode_Polyfill, key string) {
	b.self = self
	b.outer = outer
	b.isRemovable = true
	b.removable = removable
	b.key = key
}

// GetRemoveUnusedCode_PolyfillInfoBase gives access to the base through the
// interface.
func (b *RemoveUnusedCode_PolyfillInfoBase) GetRemoveUnusedCode_PolyfillInfoBase() *RemoveUnusedCode_PolyfillInfoBase {
	return b
}

// ConsiderPossibleReference accepts a NAME or GETPROP node whose (property)
// string matches key and checks whether the node should be considered as a
// possible reference to this polyfill. If so, mark the polyfill as
// referenced and therefore not removable.
//
// Java: RemoveUnusedCode.PolyfillInfo.considerPossibleReference(Node)
func (b *RemoveUnusedCode_PolyfillInfoBase) ConsiderPossibleReference(n *rhino.Node) {
	if b.isRemovable && (b.removable.isPatch || b.outer.guardedUsages == nil || !b.outer.guardedUsages.Contains(n)) {
		b.self.ConsiderPossibleReferenceInternal(n)
		if !b.isRemovable {
			b.removable.ApplyContinuations()
		}
	}
}

// Java: RemoveUnusedCode.PolyfillInfo.getName()
func (b *RemoveUnusedCode_PolyfillInfoBase) GetName() string {
	return b.key
}

// RemoveUnusedCode_GlobalPolyfillInfo is the inner class GlobalPolyfillInfo.
type RemoveUnusedCode_GlobalPolyfillInfo struct {
	RemoveUnusedCode_PolyfillInfoBase
}

// Java: RemoveUnusedCode.GlobalPolyfillInfo(Polyfill, String)
func NewRemoveUnusedCode_GlobalPolyfillInfo(outer *RemoveUnusedCode, removable *RemoveUnusedCode_Polyfill, name string) *RemoveUnusedCode_GlobalPolyfillInfo {
	g := &RemoveUnusedCode_GlobalPolyfillInfo{}
	g.InitRemoveUnusedCode_PolyfillInfoBase(g, outer, removable, name)
	return g
}

// Java: RemoveUnusedCode.GlobalPolyfillInfo.considerPossibleReferenceInternal(Node)
func (g *RemoveUnusedCode_GlobalPolyfillInfo) ConsiderPossibleReferenceInternal(possiblyReferencingNode *rhino.Node) {
	if possiblyReferencingNode.IsName() {
		// A matching NAME node must be a reference (there's no need to check that the referenced
		// Var is global, since local variables have all been renamed by normalization).
		g.isRemovable = false
	} else if NodeUtil_IsNormalOrOptChainGetProp(possiblyReferencingNode) {
		// Assume that the owner is possibly the global `this` and skip removal.
		g.isRemovable = false
	}
}

// RemoveUnusedCode_StaticPropertyPolyfillInfo is the inner class
// StaticPropertyPolyfillInfo.
type RemoveUnusedCode_StaticPropertyPolyfillInfo struct {
	RemoveUnusedCode_PolyfillInfoBase
	// Name of the owning type, used only for debugging.
	polyfillOwnerName string
}

// Java: RemoveUnusedCode.StaticPropertyPolyfillInfo(Polyfill, String, String)
func NewRemoveUnusedCode_StaticPropertyPolyfillInfo(outer *RemoveUnusedCode, removable *RemoveUnusedCode_Polyfill, key string, ownerName string) *RemoveUnusedCode_StaticPropertyPolyfillInfo {
	s := &RemoveUnusedCode_StaticPropertyPolyfillInfo{}
	s.InitRemoveUnusedCode_PolyfillInfoBase(s, outer, removable, key)
	s.polyfillOwnerName = ownerName // checkNotNull: a Go string is never null.
	return s
}

// Java: RemoveUnusedCode.StaticPropertyPolyfillInfo.getName()
func (s *RemoveUnusedCode_StaticPropertyPolyfillInfo) GetName() string {
	return s.polyfillOwnerName + "." + s.key
}

// Java: RemoveUnusedCode.StaticPropertyPolyfillInfo.considerPossibleReferenceInternal(Node)
func (s *RemoveUnusedCode_StaticPropertyPolyfillInfo) ConsiderPossibleReferenceInternal(possiblyReferencingNode *rhino.Node) {
	if NodeUtil_IsNormalOrOptChainGetProp(possiblyReferencingNode) {
		s.isRemovable = false
	}
}

// RemoveUnusedCode_PrototypePropertyPolyfillInfo is the inner class
// PrototypePropertyPolyfillInfo.
type RemoveUnusedCode_PrototypePropertyPolyfillInfo struct {
	RemoveUnusedCode_PolyfillInfoBase
	// Name of the owning type, used only for debugging.
	polyfillOwnerName string
}

// Java: RemoveUnusedCode.PrototypePropertyPolyfillInfo(Polyfill, String, String)
func NewRemoveUnusedCode_PrototypePropertyPolyfillInfo(outer *RemoveUnusedCode, removable *RemoveUnusedCode_Polyfill, key string, polyfillOwnerName string) *RemoveUnusedCode_PrototypePropertyPolyfillInfo {
	p := &RemoveUnusedCode_PrototypePropertyPolyfillInfo{}
	p.InitRemoveUnusedCode_PolyfillInfoBase(p, outer, removable, key)
	p.polyfillOwnerName = polyfillOwnerName // checkNotNull: a Go string is never null.
	return p
}

// Java: RemoveUnusedCode.PrototypePropertyPolyfillInfo.getName()
func (p *RemoveUnusedCode_PrototypePropertyPolyfillInfo) GetName() string {
	return p.polyfillOwnerName + ".prototype." + p.key
}

// Java: RemoveUnusedCode.PrototypePropertyPolyfillInfo.considerPossibleReferenceInternal(Node)
func (p *RemoveUnusedCode_PrototypePropertyPolyfillInfo) ConsiderPossibleReferenceInternal(possiblyReferencingNode *rhino.Node) {
	if NodeUtil_IsNormalOrOptChainGetProp(possiblyReferencingNode) {
		// Prototype properties are simply not removable.
		p.isRemovable = false
	}
}

// RemoveUnusedCode_VanillaForNameDeclaration represents declarations in the
// standard for-loop initialization.
//
// e.g. the `let i = 0` part of `for (let i = 0; i < 10; ++i) {...}`. These
// must be handled differently from declaration statements because:
//
//  1. For-loop declarations may declare more than one variable. The
//     normalization doesn't break them up as it does for declaration
//     statements.
//  2. Removal must be handled differently.
//  3. We don't currently preserve initializers with side effects here.
//     Instead, we just consider such cases non-removable.
type RemoveUnusedCode_VanillaForNameDeclaration struct {
	RemoveUnusedCode_RemovableBase
	nameNode *rhino.Node
}

// Java: RemoveUnusedCode.VanillaForNameDeclaration(RemovableBuilder, Node)
func NewRemoveUnusedCode_VanillaForNameDeclaration(outer *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, nameNode *rhino.Node) *RemoveUnusedCode_VanillaForNameDeclaration {
	d := outer.vanillaForDeclSlab.New()
	d.InitRemoveUnusedCode_RemovableBase(d, outer, nameNode, builder)
	d.nameNode = nameNode
	return d
}

// Java: RemoveUnusedCode.VanillaForNameDeclaration.removeInternal(AbstractCompiler)
func (d *RemoveUnusedCode_VanillaForNameDeclaration) RemoveInternal(compiler AbstractCompiler) {
	declaration := jutil.CheckNotNull(d.nameNode.GetParent())
	compiler.ReportChangeToEnclosingScope(declaration)
	// NOTE: We don't need to preserve the initializer value, because we currently do not remove
	//     for-loop vars whose initializing values have side effects.
	if d.nameNode.GetPrevious() == nil && d.nameNode.GetNext() == nil {
		// only child, so we can remove the whole declaration
		declaration.ReplaceWith(rhino.IR_Empty().Srcref(declaration))
	} else {
		d.nameNode.Detach()
	}
	NodeUtil_MarkFunctionsDeleted(d.nameNode, compiler)
}

// Java: RemoveUnusedCode.VanillaForNameDeclaration.isVariableAssignment()
func (d *RemoveUnusedCode_VanillaForNameDeclaration) IsVariableAssignment() bool {
	return true
}

// Java: RemoveUnusedCode.VanillaForNameDeclaration.isAssignedValueLocal()
func (d *RemoveUnusedCode_VanillaForNameDeclaration) IsAssignedValueLocal() bool {
	initialValueNode := d.nameNode.GetFirstChild()
	if initialValueNode == nil {
		// `var foo;`
		// the "assigned" value is undefined, which should be considered a "local" value,
		// since it is a constant.
		return true
	}
	// Handle `var name = name || defaultValue;`
	valueNode := RemoveUnusedCode_MaybeUnwrapQnameOrDefaultValueNode(d.nameNode, initialValueNode)
	return NodeUtil_EvaluatesToLocalValue(valueNode)
}

// Java: RemoveUnusedCode.VanillaForNameDeclaration.getLocalAssignedValue()
func (d *RemoveUnusedCode_VanillaForNameDeclaration) GetLocalAssignedValue() *rhino.Node {
	initialValueNode := d.nameNode.GetFirstChild()
	if initialValueNode == nil {
		// `var foo;` has no node to represent the `undefined` value that is assigned.
		return nil
	}
	// Handle `var name = name || defaultValue;`
	valueNode := RemoveUnusedCode_MaybeUnwrapQnameOrDefaultValueNode(d.nameNode, initialValueNode)
	if NodeUtil_EvaluatesToLocalValue(valueNode) {
		return valueNode
	}
	return nil
}

// Java: RemoveUnusedCode.removeExpressionCompletely(Node)
func (r *RemoveUnusedCode) RemoveExpressionCompletely(expression *rhino.Node) {
	jutil.CheckStateMsg1(!NodeUtil_IsExpressionResultUsed(expression), "%s", expression)
	parent := expression.GetParent()
	if parent.IsExprResult() {
		NodeUtil_DeleteNode(parent, r.compiler)
	} else if parent.IsComma() {
		// Expression is probably the first child of the comma,
		// but it could be the second if the entire comma expression value is unused.
		otherChild := expression.GetNext()
		if otherChild == nil {
			otherChild = expression.GetPrevious()
		}
		r.ReplaceNodeWith(parent, otherChild.Detach())
	} else {
		// value isn't needed, but we need to keep the AST valid.
		r.ReplaceNodeWith(expression, rhino.IR_Number(0).Srcref(expression))
	}
}

// Java: RemoveUnusedCode.replaceNodeWith(Node, Node)
func (r *RemoveUnusedCode) ReplaceNodeWith(n *rhino.Node, replacement *rhino.Node) {
	r.compiler.ReportChangeToEnclosingScope(n)
	n.ReplaceWith(replacement)
	NodeUtil_MarkFunctionsDeleted(n, r.compiler)
}
