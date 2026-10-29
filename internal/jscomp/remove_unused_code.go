// Port of com/google/javascript/jscomp/RemoveUnusedCode.java @ closure-compiler v20261021

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
	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/diagnostic"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Go-only: the traversal methods and the Removables read the shapes they
// handle through the typed view (internal/ast): the parts of assignments,
// declarations, loops, classes, functions, parameters and patterns, and
// the x.prototype and this.x forms. The generic subtree walks, the
// removals' tree edits and the parallel bookkeeping stay on rhino.Node.

type rucPropEntry struct {
	remHead         *RemoveUnusedCode_RemovableBase
	remTail         *RemoveUnusedCode_RemovableBase
	polyHead        *rucPolyfillInfo
	polyTail        *rucPolyfillInfo
	propID          uint32
	checkedUnpinned bool
	inRemOrder      bool
}

// removeUnusedCodeState is the state of one run of the pass; Process drops
// it when the run ends.
type removeUnusedCodeState struct {
	worklist     []RemoveUnusedCode_Continuation
	worklistHead int
	worklistSize int

	// Java: varInfoMap (LinkedHashMap<Var, VarInfo>, where Var equality is
	// (scope root, name)). varByKey indexes VarInfos by scopedVarKey and
	// varInfos preserves insertion order.
	varByKey map[scopedVarKey]*RemoveUnusedCode_RealVarInfo
	varInfos []*RemoveUnusedCode_RealVarInfo

	trackProps   bool
	pinnedProps  inlineBitSet
	polyBits     inlineBitSet
	propIdxByID  map[uint32]uint32
	propEntries  []rucPropEntry
	propRemOrder []uint32

	// Stores all polyfills seen during traversal, indexed by leaf property ID in
	// propEntries (gated by polyBits) and ordered by polyfillOrder.
	polyfillOrder         []*rucPolyfillInfo
	guardedUsages         map[rhino.Node]bool
	guardedUsagesComputed bool
	root                  rhino.Node

	allFunctionParamScopes []*Scope

	extraContSlab []RemoveUnusedCode_Continuation

	// Go-only: performance — one run's removables, var infos, and polyfill infos
	// are created together and die together with this state (see jutil.Slab).
	removables   jutil.Slab[RemoveUnusedCode_RemovableBase]
	varInfoSlab  jutil.Slab[RemoveUnusedCode_RealVarInfo]
	polyfillSlab jutil.Slab[rucPolyfillInfo]
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

	// Single value to use for all vars for which we cannot remove anything at
	// all.
	canonicalUnremovableVarInfo    *RemoveUnusedCode_RealVarInfo
	canonicalUnremovableVarInfoVal RemoveUnusedCode_RealVarInfo

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

	removeUnusedCodeState
}

var (
	// Properties that are implicitly used as part of the JS language.
	rucPropLengthID      = rhino.RhinoStringPool_InternID("length")
	rucPropToStringID    = rhino.RhinoStringPool_InternID("toString")
	rucPropValueOfID     = rhino.RhinoStringPool_InternID("valueOf")
	rucPropConstructorID = rhino.RhinoStringPool_InternID("constructor")
	rucPropPrototypeID   = rhino.RhinoStringPool_InternID("prototype")

	rucImplicitPropIDs = [5]uint32{
		rucPropLengthID,
		rucPropToStringID,
		rucPropValueOfID,
		rucPropConstructorID,
		rucPropPrototypeID,
	}

	rucCalleeInheritsID                 = rhino.RhinoStringPool_InternID("inherits")
	rucCalleeMixinID                    = rhino.RhinoStringPool_InternID("mixin")
	rucCalleeDefinePropertiesID         = rhino.RhinoStringPool_InternID("defineProperties")
	rucCalleeObjectPropertyID           = rhino.RhinoStringPool_InternID("objectProperty")
	rucCalleeAddSingletonGetterID       = rhino.RhinoStringPool_InternID("addSingletonGetter")
	rucCalleePolyfillID                 = rhino.RhinoStringPool_InternID("polyfill")
	rucCalleePatchID                    = rhino.RhinoStringPool_InternID("patch")
	rucCalleePolyfillTypedArrayMethodID = rhino.RhinoStringPool_InternID("polyfillTypedArrayMethod")
	rucCalleeWeakUsageID                = rhino.RhinoStringPool_InternID("weakUsage")

	rucCalleeJscPropNameFnID                  = rhino.RhinoStringPool_InternID(NodeUtil_JSC_PROPERTY_NAME_FN)
	rucCalleeGoogReflectObjectPropertyID      = rhino.RhinoStringPool_InternID("goog$reflect$objectProperty")
	rucCalleeJscompPolyfillID                 = rhino.RhinoStringPool_InternID("$jscomp$polyfill")
	rucCalleeJscompPatchID                    = rhino.RhinoStringPool_InternID("$jscomp$patch")
	rucCalleeJscompPolyfillTypedArrayMethodID = rhino.RhinoStringPool_InternID("$jscomp$polyfillTypedArrayMethod")
	rucCalleeGoogWeakUsageID                  = rhino.RhinoStringPool_InternID("goog$weakUsage")
	rucCalleeGoogAddSingletonGetterID         = rhino.RhinoStringPool_InternID("goog$addSingletonGetter")
	rucJscompIdentID                          = rhino.RhinoStringPool_InternID("$jscomp")
)

// RemoveUnusedCodeOptions configures a RemoveUnusedCode pass.
type RemoveUnusedCodeOptions struct {
	RemoveLocalVars                               bool
	RemoveGlobals                                 bool
	PreserveFunctionExpressionNames               bool
	RemoveUnusedPrototypeProperties               bool
	RemoveUnusedThisProperties                    bool
	RemoveUnusedObjectDefinePropertiesDefinitions bool
	RemoveUnusedPolyfills                         bool
	AssumeGettersArePure                          bool
}

// Java: RemoveUnusedCode(Builder)
func NewRemoveUnusedCode(compiler AbstractCompiler, opts RemoveUnusedCodeOptions) *RemoveUnusedCode {
	sc := NewSyntacticScopeCreator(compiler)
	r := &RemoveUnusedCode{
		compiler:                        compiler,
		astAnalyzer:                     compiler.GetAstAnalyzer(),
		codingConvention:                compiler.GetCodingConvention(),
		scopeCreator:                    sc,
		removeLocalVars:                 opts.RemoveLocalVars,
		removeGlobals:                   opts.RemoveGlobals,
		preserveFunctionExpressionNames: opts.PreserveFunctionExpressionNames,
		removeUnusedPrototypeProperties: opts.RemoveUnusedPrototypeProperties,
		removeUnusedThisProperties:      opts.RemoveUnusedThisProperties,
		removeUnusedObjectDefinePropertiesDefinitions: opts.RemoveUnusedObjectDefinePropertiesDefinitions,
		removeUnusedPolyfills:                         opts.RemoveUnusedPolyfills,
		polyfillsFromTable:                            defaultPolyfillsTable(),
		assumeGettersArePure:                          opts.AssumeGettersArePure,
	}

	// All Vars that are completely unremovable will share this VarInfo instance.
	r.canonicalUnremovableVarInfoVal.outer = r
	r.canonicalUnremovableVarInfoVal.flags = rucVarFlagCanonical
	r.canonicalUnremovableVarInfo = &r.canonicalUnremovableVarInfoVal
	return r
}

func (r *RemoveUnusedCode) logRemoval(kind, name, fnName string) {
	r.removalLog.LogString(kind + "\t" + name + "\t" + fnName)
}

func loggableFunctionName(argList rhino.Node) string {
	functionName := NodeUtil_GetNearestFunctionName(rhino.CheckNotNull(argList).GetParent())
	if functionName == nil {
		return "<anonymous>"
	}
	return *functionName
}

func (r *RemoveUnusedCode) removesAnyProperties() bool {
	return r.removeUnusedPrototypeProperties ||
		r.removeUnusedThisProperties ||
		r.removeUnusedObjectDefinePropertiesDefinitions
}

func (r *RemoveUnusedCode) ensurePropEntry(id uint32) uint32 {
	if idx1 := r.propIdxByID[id]; idx1 != 0 {
		return idx1 - 1
	}
	if r.propIdxByID == nil {
		r.propIdxByID = make(map[uint32]uint32, 64)
	}
	idx := uint32(len(r.propEntries))
	r.propEntries = append(r.propEntries, rucPropEntry{propID: id})
	r.propIdxByID[id] = idx + 1
	return idx
}

func (r *RemoveUnusedCode) considerPolyfillRef(keyID uint32, n rhino.Node) {
	idx1 := r.propIdxByID[keyID]
	if idx1 == 0 {
		return
	}
	for info := r.propEntries[idx1-1].polyHead; info != nil; info = info.nextByKey {
		if info.isRemovable {
			info.considerPossibleReference(n)
		}
	}
}

// Process traverses the root, removing all unused variables. Multiple
// traversals may occur to ensure all unused variables are removed.
//
// Java: RemoveUnusedCode.process(Node, Node)
func (r *RemoveUnusedCode) Process(externs rhino.Node, root rhino.Node) {
	if !r.compiler.GetLifeCycleStage().IsNormalized() {
		jutil.CheckState(false)
	}
	r.externProperties = r.compiler.GetExternProperties()
	r.accessorSummary = r.compiler.GetAccessorSummary()

	defer func() {
		r.removeUnusedCodeState = removeUnusedCodeState{}
	}()

	removalLogFile :=
		r.compiler.CreateOrReopenIndexedLog("RemoveUnusedCode", "removals.log")
	defer removalLogFile.Close()
	keepLogFile :=
		r.compiler.CreateOrReopenIndexedLog("RemoveUnusedCode", "unremovable.log")
	defer keepLogFile.Close()
	defer func() {
		r.removalLog = diagnostic.LogFile{}
		r.unremovableLog = diagnostic.LogFile{}
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

func (r *RemoveUnusedCode) createGlobalScope(root rhino.Node) *Scope {
	// Create scope from parent of root node, which also has externs as a child, so we'll
	// have extern definitions in scope.
	scope := r.scopeCreator.CreateScopeOfScope(root.GetParent(), nil)
	if !scope.HasSlotByID(rucCalleeJscPropNameFnID) {
		// TODO(b/70730762): Passes that add references to this should ensure it is declared.
		// NOTE: null input makes this an extern var.
		scope.Declare(
			NodeUtil_JSC_PROPERTY_NAME_FN /* no declaration node */, rhino.NilNode /* no input */, nil)
	}
	return scope
}

func (r *RemoveUnusedCode) traverseRootAndRemoveVars(node rhino.Node, scope *Scope) {
	r.root = node
	r.guardedUsagesComputed = false

	r.enqueueContinuation(RemoveUnusedCode_Continuation{node: node, scope: scope})
	for r.worklistSize > 0 {
		continuation := r.popContinuation() // worklist.remove()
		continuation.Apply(r)
	}

	r.RemoveUnreferencedVarsAndPolyfills()
}

// TraverseAndRemoveUnusedReferences traverses a node recursively. Call this
// once per pass.
//
// Java: RemoveUnusedCode.traverseAndRemoveUnusedReferences(Node)
func (r *RemoveUnusedCode) TraverseAndRemoveUnusedReferences(root rhino.Node) {
	if r.traverseAndRemoveLocalScriptsParallel(root) {
		return
	}
	r.trackProps = r.removesAnyProperties()
	if r.trackProps {
		for _, id := range rucImplicitPropIDs {
			r.pinnedProps.Set(int(id))
		}
	}
	scope := r.createGlobalScope(root)
	if r.varByKey == nil {
		r.varByKey = make(map[scopedVarKey]*RemoveUnusedCode_RealVarInfo, 1024)
	}

	r.traverseRootAndRemoveVars(root, scope)
	r.RemoveIndependentlyRemovableProperties()
	if r.removeGlobals && !r.removeUnreferencedFunctionArgsParallel() {
		for _, fparamScope := range r.allFunctionParamScopes {
			r.removeUnreferencedFunctionArgsWithCompiler(r.compiler, fparamScope, false)
		}
	}
}

func (r *RemoveUnusedCode) canParallelizeLocalScripts() bool {
	return !r.removeGlobals &&
		!r.removeUnusedPrototypeProperties &&
		!r.removeUnusedThisProperties &&
		!r.removeUnusedObjectDefinePropertiesDefinitions &&
		!r.removeUnusedPolyfills &&
		!r.removalLog.IsLogging() &&
		!r.unremovableLog.IsLogging()
}

func (r *RemoveUnusedCode) traverseAndRemoveLocalScriptsParallel(root rhino.Node) bool {
	if !r.canParallelizeLocalScripts() || root.GetParent().IsNil() {
		return false
	}
	scripts, ok := collectRootScripts(root)
	if !ok {
		return false
	}
	numWorkers := parallelThreadCount(r.compiler, len(scripts))
	if numWorkers <= 1 {
		return false
	}

	prewarmCompilerForParallelReads(r.compiler)
	restoreReg := enableParallelTypeRegistryReads(r.compiler)
	defer restoreReg()

	globalScope := r.createGlobalScope(root)
	freezeGlobalScopeForParallelReads(globalScope)

	runParallelWork(r.compiler, numWorkers, len(scripts), func(_ int) (func(int, *recordingCompiler), func()) {
		var wr *RemoveUnusedCode
		runItem := func(idx int, rc *recordingCompiler) {
			if wr == nil || wr.compiler != rc {
				wr = &RemoveUnusedCode{
					compiler:                        rc,
					astAnalyzer:                     r.astAnalyzer,
					codingConvention:                r.codingConvention,
					scopeCreator:                    NewSyntacticScopeCreator(rc),
					removeLocalVars:                 r.removeLocalVars,
					removeGlobals:                   false,
					preserveFunctionExpressionNames: r.preserveFunctionExpressionNames,
					externProperties:                r.externProperties,
					accessorSummary:                 r.accessorSummary,
					polyfillsFromTable:              r.polyfillsFromTable,
					assumeGettersArePure:            r.assumeGettersArePure,
					removalLog:                      r.removalLog,
					unremovableLog:                  r.unremovableLog,
				}
				wr.canonicalUnremovableVarInfoVal.outer = wr
				wr.canonicalUnremovableVarInfoVal.flags = rucVarFlagCanonical
				wr.canonicalUnremovableVarInfo = &wr.canonicalUnremovableVarInfoVal
				wr.varByKey = make(map[scopedVarKey]*RemoveUnusedCode_RealVarInfo, 128)
			} else {
				clear(wr.varByKey)
				clear(wr.varInfos)
				wr.varInfos = wr.varInfos[:0]
				wr.worklistHead = 0
				wr.worklistSize = 0
				wr.extraContSlab = wr.extraContSlab[:0]
			}
			wr.traverseRootAndRemoveVars(scripts[idx], globalScope)
		}
		return runItem, nil
	}, false)
	return true
}

// trailingParamBinding returns what the parameter p binds: the name or
// pattern, through a default value or a rest. removable is false when p has
// a default value with side effects, which keeps it.
func (r *RemoveUnusedCode) trailingParamBinding(p ast.Param) (binding rhino.Node, removable bool) {
	binding = p.Node()
	if dv, ok := p.AsDefaultValue(); ok {
		if r.MayHaveSideEffects(dv.Value().Node()) {
			return rhino.NilNode, false
		}
		binding = dv.Target().Node()
	}
	if rest, ok := ast.AsIterRest(binding); ok {
		binding = rest.Target().Node()
	}
	return binding, true
}

func (r *RemoveUnusedCode) isParamVarRemovable(nameNode rhino.Node, fparamScope *Scope, readOnly bool) bool {
	if !readOnly && r.unremovableLog.IsLogging() {
		return r.GetVarInfo(r.GetVarForNameNode(nameNode, fparamScope)).IsRemovable()
	}
	v := fparamScope.GetVarByNode(nameNode)
	if v == nil {
		jutil.CheckNotNullMsg1(v, "%s", nameNode)
	}
	isGlobal := v.IsGlobal()
	if v.IsExtern() || r.isExportedByConvention(v.GetName(), !isGlobal) || v.IsArguments() {
		return false
	}
	if vi := r.varByKey[scopedVarKeyOf(v)]; vi != nil {
		return vi.IsRemovable()
	}
	if !r.removeGlobals && isGlobal {
		return false
	}
	return r.removeLocalVars || isGlobal
}

// removeUnreferencedFunctionArgsWithCompiler removes unreferenced arguments from a
// function declaration and when possible the function's callSites.
//
// Notice that removing unreferenced function args breaks
// Function.prototype.length. In advanced mode, we don't really care
// about this: we consider "length" the equivalent of reflecting on
// the function's lexical source.
//
// Rather than create a new option for this, we assume that if the user
// is removing globals, then it's OK to remove unused function args.
//
// See http://blickly.github.io/closure-compiler-issues/#253
//
// Java: RemoveUnusedCode.removeUnreferencedFunctionArgs(Scope)
func (r *RemoveUnusedCode) removeUnreferencedFunctionArgsWithCompiler(
	compiler AbstractCompiler,
	fparamScope *Scope,
	readOnly bool,
) {
	function, ok := ast.AsFunction(fparamScope.GetRootNode())
	if !ok {
		jutil.CheckState(false)
	}
	if NodeUtil_IsGetOrSetKey(function.Node().GetParent()) {
		// The parameters object literal setters can not be removed.
		return
	}
	argList := function.Params().Node()
	// Strip as many unreferenced args off the end of the function declaration as possible.
	r.maybeRemoveUnusedTrailingParameters(compiler, argList, fparamScope, readOnly)
	// Mark any remaining unused parameters are unused to OptimizeParameters can try to remove
	// them.
	r.markUnusedParameters(compiler, argList, fparamScope, readOnly)
}

type rucParamScriptGroup struct {
	script  rhino.Node
	indices []int
}

func (r *RemoveUnusedCode) removeUnreferencedFunctionArgsParallel() bool {
	if !r.removeGlobals ||
		len(r.allFunctionParamScopes) < 2 ||
		r.removalLog.IsLogging() ||
		r.unremovableLog.IsLogging() {
		return false
	}
	if r.compiler == nil || r.compiler.GetOptions() == nil || r.compiler.GetOptions().GetNumParallelThreads() <= 1 {
		return false
	}

	groupByScript := make(map[rhino.Node]int, 64)
	var groups []rucParamScriptGroup
	var detachedIndices []int
	var anyScope *Scope
	for idx, fparamScope := range r.allFunctionParamScopes {
		if anyScope == nil && fparamScope != nil {
			anyScope = fparamScope
		}
		script := NodeUtil_GetEnclosingScript(fparamScope.GetRootNode())
		if script.IsNil() {
			detachedIndices = append(detachedIndices, idx)
			continue
		}
		gIdx, exists := groupByScript[script]
		if !exists {
			gIdx = len(groups)
			groupByScript[script] = gIdx
			groups = append(groups, rucParamScriptGroup{script: script})
		}
		groups[gIdx].indices = append(groups[gIdx].indices, idx)
	}

	numWorkers := parallelThreadCount(r.compiler, len(groups))
	if numWorkers <= 1 {
		return false
	}

	prewarmCompilerForParallelReads(r.compiler)
	restoreReg := enableParallelTypeRegistryReads(r.compiler)
	defer restoreReg()
	if anyScope != nil {
		freezeGlobalScopeForParallelReads(anyScope.GetGlobalScope())
	}

	eventsByIdx := make([][]recordedCompilerEvent, len(r.allFunctionParamScopes))
	runParallelWork(r.compiler, numWorkers, len(groups), func(_ int) (func(int, *recordingCompiler), func()) {
		runItem := func(gIdx int, rc *recordingCompiler) {
			g := &groups[gIdx]
			for _, origIdx := range g.indices {
				resetRecordingCompiler(rc, r.compiler)
				r.removeUnreferencedFunctionArgsWithCompiler(rc, r.allFunctionParamScopes[origIdx], true)
				eventsByIdx[origIdx] = rc.takeEvents()
			}
		}
		return runItem, nil
	}, false)

	if len(detachedIndices) > 0 {
		var mainRC recordingCompiler
		initRecordingCompiler(&mainRC, r.compiler)
		for _, origIdx := range detachedIndices {
			resetRecordingCompiler(&mainRC, r.compiler)
			r.removeUnreferencedFunctionArgsWithCompiler(&mainRC, r.allFunctionParamScopes[origIdx], true)
			eventsByIdx[origIdx] = mainRC.takeEvents()
		}
	}

	for _, evs := range eventsByIdx {
		if len(evs) > 0 {
			replayCompilerEvents(r.compiler, evs)
		}
	}
	return true
}

// Java: RemoveUnusedCode.storePolyfill(PolyfillUsage)
func (r *RemoveUnusedCode) StorePolyfill(polyfillUsage *PolyfillUsageFinder_PolyfillUsage) {
	if r.guardedUsages == nil {
		r.guardedUsages = make(map[rhino.Node]bool)
	}
	r.guardedUsages[polyfillUsage.Node()] = true
}

// Java: RemoveUnusedCode.removeIndependentlyRemovableProperties()
func (r *RemoveUnusedCode) RemoveIndependentlyRemovableProperties() {
	for _, idx := range r.propRemOrder {
		cur := r.propEntries[idx].remHead
		if cur == nil {
			continue
		}
		r.propEntries[idx].remHead = nil
		r.propEntries[idx].remTail = nil
		if r.removalLog.IsLogging() {
			r.logRemoval("prop", rhino.RhinoStringPool_StringFromID(r.propEntries[idx].propID), "")
		}
		for cur != nil {
			next := cur.propNext
			cur.propNext = nil
			cur.flags &^= remFlagInPropList
			cur.Remove(r)
			cur = next
		}
	}
}

// Java: RemoveUnusedCode.getVarForNameNode(Node, Scope)
func (r *RemoveUnusedCode) GetVarForNameNode(nameNode rhino.Node, scope *Scope) *Var {
	v := scope.GetVarByNode(nameNode)
	if v == nil {
		jutil.CheckNotNullMsg1(v, "%s", nameNode)
	}
	return v
}

// Java: RemoveUnusedCode.canRemoveParameters(Node)
func (r *RemoveUnusedCode) CanRemoveParameters(parameterList rhino.Node) bool {
	if !parameterList.IsParamList() {
		jutil.CheckState(false)
	}
	function := parameterList.GetParent()
	return r.removeGlobals && !NodeUtil_IsGetOrSetKey(function.GetParent())
}

func (r *RemoveUnusedCode) markPropertyIDPinned(id uint32) {
	if !r.trackProps || r.pinnedProps.Get(int(id)) {
		return
	}
	r.pinnedProps.Set(int(id))
	// Continue traversal of all of the property name's values and no longer consider them for
	// removal.
	if idx1 := r.propIdxByID[id]; idx1 != 0 {
		entry := &r.propEntries[idx1-1]
		cur := entry.remHead
		entry.remHead = nil
		entry.remTail = nil
		for cur != nil {
			next := cur.propNext
			cur.propNext = nil
			cur.flags &^= remFlagInPropList
			cur.ApplyContinuations(r)
			cur = next
		}
	}
}

func (r *RemoveUnusedCode) isExportedByConvention(name string, local bool) bool {
	if len(name) == 0 || (name[0] != '_' && name[0] != '$') {
		return false
	}
	return r.codingConvention.IsExportedWithLocal(name, local)
}

// Java: RemoveUnusedCode.considerForIndependentRemoval(Removable)
func (r *RemoveUnusedCode) ConsiderForIndependentRemoval(removable *RemoveUnusedCode_RemovableBase) {
	propID := removable.propertyNameID
	if !r.trackProps || propID == 0 {
		removable.ApplyContinuations(r)
		return
	}
	if r.pinnedProps.Get(int(propID)) {
		// Referenced or exported, so not removable.
		removable.ApplyContinuations(r)
		return
	}
	idx := r.ensurePropEntry(propID)
	if !r.propEntries[idx].checkedUnpinned {
		name := rhino.RhinoStringPool_StringFromID(propID)
		if (r.externProperties != nil && r.externProperties.Contains(name)) ||
			r.isExportedByConvention(name /* local= */, false) {
			r.pinnedProps.Set(int(propID))
			removable.ApplyContinuations(r)
			return
		}
		r.propEntries[idx].checkedUnpinned = true
	}
	if r.IsIndependentlyRemovable(removable) {
		// Store for possible removal later.
		if removable.flags&remFlagInPropList == 0 {
			removable.flags |= remFlagInPropList
			removable.propNext = nil
			entry := &r.propEntries[idx]
			if entry.remHead == nil {
				entry.remHead = removable
				entry.remTail = removable
				if !entry.inRemOrder {
					entry.inRemOrder = true
					r.propRemOrder = append(r.propRemOrder, idx)
				}
			} else {
				entry.remTail.propNext = removable
				entry.remTail = removable
			}
		}
	} else {
		removable.ApplyContinuations(r)
		// This assignment counts as a reference, since we won't be removing it.
		// This is necessary in order to preserve getters and setters for the property.
		r.markPropertyIDPinned(propID)
	}
}

// ConsiderForAccessorSideEffects returns whether or not accessor
// side-effect are a possibility.
//
// Java: RemoveUnusedCode.considerForAccessorSideEffects(Node, PropertyAccessKind)
func (r *RemoveUnusedCode) ConsiderForAccessorSideEffects(getprop rhino.Node, usage AccessorSummary_PropertyAccessKind) bool {
	// Other node types may make sense in the future.
	if !NodeUtil_IsNormalOrOptChainGetProp(getprop) {
		jutil.CheckStateMsg1(false, "%s", getprop)
	}
	if r.assumeGettersArePure && !usage.HasSetter() {
		return false
	}
	summary := r.accessorSummary
	if summary == nil {
		summary = r.compiler.GetAccessorSummary()
		r.accessorSummary = summary
	}
	propID := getprop.GetStringID()
	recorded := summary.GetKindByID(propID)
	if (recorded.HasGetter() && usage.HasGetter() && !r.assumeGettersArePure) ||
		(recorded.HasSetter() && usage.HasSetter()) {
		r.markPropertyIDPinned(propID)
		return true
	}

	return false
}

// Java: RemoveUnusedCode.isIndependentlyRemovable(Removable)
func (r *RemoveUnusedCode) IsIndependentlyRemovable(removable *RemoveUnusedCode_RemovableBase) bool {
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
	}
	return false
}

// markUnusedParameters marks any remaining unused parameters as being unused
// so it can be used elsewhere.
//
// paramList is the list of function's parameters.
//
// Java: RemoveUnusedCode.markUnusedParameters(Node, Scope)
func (r *RemoveUnusedCode) markUnusedParameters(compiler AbstractCompiler, paramList rhino.Node, fparamScope *Scope, readOnly bool) {
	params, ok := ast.AsParamList(paramList)
	if !ok {
		jutil.CheckArgumentMsg1(false, "%s", paramList)
	}

	for p := range params.Params() {
		param := p.Node()
		if param.IsUnusedParameter() {
			// already marked
			continue
		}

		paramNameNode := RemoveUnusedCode_NameOfParam(param)
		if paramNameNode.IsNil() {
			// destructuring pattern parameters don't have a name that applies to the whole parameter
			// TODO(bradfordcsmith): We could mark this if we determined that all vars created by
			// the pattern are unused.
			continue
		}

		if r.isParamVarRemovable(paramNameNode, fparamScope, readOnly) {
			param.SetUnusedParameter(true)
			compiler.ReportChangeToEnclosingScope(paramList)
			if !readOnly && r.removalLog.IsLogging() {
				r.logRemoval("argmark", paramNameNode.GetString(), loggableFunctionName(paramList))
			}
		}
	}
}

// maybeRemoveUnusedTrailingParameters strips as many unreferenced args off
// the end of the function declaration as possible. We start from the end of
// the function declaration because removing parameters from the middle of
// the param list could mess up the interpretation of parameters being sent
// over by any function calls.
//
// argList is the list of function's arguments.
//
// Java: RemoveUnusedCode.maybeRemoveUnusedTrailingParameters(Node, Scope)
func (r *RemoveUnusedCode) maybeRemoveUnusedTrailingParameters(compiler AbstractCompiler, argList rhino.Node, fparamScope *Scope, readOnly bool) {
	params, ok := ast.AsParamList(argList)
	if !ok {
		jutil.CheckArgumentMsg1(false, "%s", argList)
	}
	for {
		last, ok := params.LastParam()
		if !ok {
			break
		}
		lastArg := last.Node()
		argNode, removable := r.trailingParamBinding(last)
		if !removable {
			break
		}

		if argNode.IsDestructuringPattern() {
			if argNode.HasChildren() {
				// TODO(johnlenz): handle the case where there are no assignments.
				break
			}
			// Remove empty destructuring patterns and their associated object literal assignment
			// if it exists and if the right hand side does not have side effects. Note, a
			// destructuring pattern with a "leftover" property key as in {a:{}} is not considered
			// empty in this case!
			NodeUtil_DeleteNode(lastArg, compiler)
			if !readOnly && r.removalLog.IsLogging() {
				r.logRemoval("arg", "<pattern>", loggableFunctionName(argList))
			}
			continue
		}

		if r.isParamVarRemovable(argNode, fparamScope, readOnly) {
			NodeUtil_DeleteNode(lastArg, compiler)
			if !readOnly && r.removalLog.IsLogging() {
				r.logRemoval("arg", argNode.GetString(), loggableFunctionName(argList))
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
func (r *RemoveUnusedCode) TraverseVar(var_ *Var) *RemoveUnusedCode_RealVarInfo {
	if var_ == nil {
		panic(jutil.NewNullPointerException())
	}
	if r.removeLocalVars && var_.IsArguments() {
		// If we are considering removing local variables, that includes parameters.
		// If `arguments` is used in a function we must consider all parameters to be referenced.
		varScope, _ := var_.GetScope().(*Scope)
		functionScope := varScope.GetClosestHoistScope()
		for param := range ast.MustFunction(functionScope.GetRootNode()).Params().Params() {
			lValue := RemoveUnusedCode_NameOfParam(param.Node())
			if lValue.IsNil() {
				continue
			}

			r.GetVarInfo(r.GetVarForNameNode(lValue, functionScope)).
				pin(rucPinParamInFnUsingArguments)
		}
		// `arguments` is never removable.
		return r.canonicalUnremovableVarInfo
	}
	return r.GetVarInfo(var_)
}

// RemoveUnusedCode_NameOfParam returns the NAME node associated with a
// function parameter (the child of a PARAM_LIST), or null if there is no
// single name.
//
// Java: RemoveUnusedCode.nameOfParam(Node)
func RemoveUnusedCode_NameOfParam(param rhino.Node) rhino.Node {
	switch param.GetToken() {
	case rhino.Token_NAME:
		return param
	case rhino.Token_DEFAULT_VALUE:
		return RemoveUnusedCode_NameOfParam(ast.MustDefaultValue(param).Target().Node())
	case rhino.Token_ITER_REST:
		return RemoveUnusedCode_NameOfParam(ast.MustIterRest(param).Target().Node())
	case rhino.Token_ARRAY_PATTERN, rhino.Token_OBJECT_PATTERN:
		return rhino.NilNode
	default:
		panic(jutil.NewIllegalStateExceptionMsg(
			"Unexpected child of PARAM_LIST: " + param.ToStringTree()))
	}
}

// GetVarInfo gets the right VarInfo object to use for the given Var.
//
// This method is responsible for managing the entries in varInfoMap.
//
// Note: Several Vars may share the same VarInfo when they should be treated
// the same way.
//
// Java: RemoveUnusedCode.getVarInfo(Var)
func (r *RemoveUnusedCode) GetVarInfo(var_ *Var) *RemoveUnusedCode_RealVarInfo {
	if var_ == nil {
		panic(jutil.NewNullPointerException())
	}
	isGlobal := var_.IsGlobal()
	if var_.IsExtern() {
		if r.unremovableLog.IsLogging() {
			r.unremovableLog.LogString(var_.GetName() + ": extern")
		}
		return r.canonicalUnremovableVarInfo
	} else if name := var_.GetName(); r.isExportedByConvention(name /* local= */, !isGlobal) {
		if r.unremovableLog.IsLogging() {
			r.unremovableLog.LogString(name + ": exported by convention")
		}
		return r.canonicalUnremovableVarInfo
	} else if var_.IsArguments() {
		// No point in logging that we cannot remove "arguments"
		return r.canonicalUnremovableVarInfo
	}
	key := scopedVarKeyOf(var_)
	if v := r.varByKey[key]; v != nil {
		return v
	}
	varInfo := NewRemoveUnusedCode_RealVarInfo(r, var_)
	if var_.GetParentNode().IsParamList() {
		varInfo.SetHasNonLocalOrNonLiteralValue()
	}
	// Cannot use canonicalUnremovableVarInfo for the 2 non-removable cases below, because each
	// varInfo needs to track what value is assigned to it for the purpose of correctly allowing
	// or preventing removal of properties set on it.
	if !r.removeGlobals && isGlobal {
		varInfo.pin(rucPinNotRemovingGlobals)
	} else if !r.removeLocalVars && !isGlobal {
		varInfo.pin(rucPinNotRemovingLocals)
	}
	r.varByKey[key] = varInfo
	r.varInfos = append(r.varInfos, varInfo)
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
			r.logRemoval("var", var_.GetName(), "")
		}
		// Regardless of what happens to the original declaration,
		// we need to remove all assigns, because they may contain references
		// to other unreferenced variables.
		varInfo.RemoveAllRemovables()

		nameNode := var_.GetNameNode()
		toRemove := nameNode.GetParent()
		if toRemove.IsNil() || RemoveUnusedCode_AlreadyRemoved(toRemove) {
			// assignedVarInfo.removeAllRemovables () already removed it
		} else if NodeUtil_IsFunctionExpression(toRemove) {
			// TODO(bradfordcsmith): Add a Removable for this case.
			if !r.preserveFunctionExpressionNames {
				r.compiler.ReportChangeToEnclosingScope(nameNode)
				ast.MustFunction(toRemove).RemoveName()
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

	if len(r.polyfillOrder) > 0 {
		for _, polyfill := range r.polyfillOrder {
			if polyfill.isRemovable {
				if r.removalLog.IsLogging() {
					r.logRemoval("poly", polyfill.getName(), "")
				}
				polyfill.isRemovable = false
				polyfill.removable.Remove(r)
			}
		}
		r.polyBits.ClearAll()
	}
}
