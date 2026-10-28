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
	"strings"

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
	pinnedProps  jutil.BitSet
	polyBits     jutil.BitSet
	propIdxByID  map[uint32]uint32
	propEntries  []rucPropEntry
	propRemOrder []uint32

	// Stores all polyfills seen during traversal, indexed by leaf property ID in
	// propEntries (gated by polyBits) and ordered by polyfillOrder.
	polyfillOrder         []*rucPolyfillInfo
	guardedUsages         *jutil.LinkedHashSet[rhino.Node]
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

// Java: RemoveUnusedCode(Builder)
func NewRemoveUnusedCode(builder *RemoveUnusedCode_Builder) *RemoveUnusedCode {
	sc := NewSyntacticScopeCreator(builder.compiler)
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

	// All Vars that are completely unremovable will share this VarInfo instance.
	r.canonicalUnremovableVarInfoVal.outer = r
	r.canonicalUnremovableVarInfoVal.flags = rucVarFlagCanonical
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
	// Create scope from parent of root node, which also has externs as a child, so we'll
	// have extern definitions in scope.
	scope := r.scopeCreator.CreateScopeOfScope(root.GetParent(), nil)
	if r.varByKey == nil {
		r.varByKey = make(map[scopedVarKey]*RemoveUnusedCode_RealVarInfo, 1024)
	}
	if !scope.HasSlot(NodeUtil_JSC_PROPERTY_NAME_FN) {
		// TODO(b/70730762): Passes that add references to this should ensure it is declared.
		// NOTE: null input makes this an extern var.
		scope.Declare(
			NodeUtil_JSC_PROPERTY_NAME_FN /* no declaration node */, rhino.NilNode /* no input */, nil)
	}

	r.root = root
	r.guardedUsagesComputed = false

	r.enqueueContinuation(RemoveUnusedCode_Continuation{node: root, scope: scope})
	for r.worklistSize > 0 {
		continuation := r.popContinuation() // worklist.remove()
		continuation.Apply(r)
	}

	r.RemoveUnreferencedVarsAndPolyfills()
	r.RemoveIndependentlyRemovableProperties()
	if !r.removeUnreferencedFunctionArgsParallel() {
		for _, fparamScope := range r.allFunctionParamScopes {
			r.RemoveUnreferencedFunctionArgs(fparamScope)
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

	globalScope := r.scopeCreator.CreateScopeOfScope(root.GetParent(), nil)
	if !globalScope.HasSlot(NodeUtil_JSC_PROPERTY_NAME_FN) {
		globalScope.Declare(NodeUtil_JSC_PROPERTY_NAME_FN, rhino.NilNode, nil)
	}
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
			script := scripts[idx]
			wr.root = script
			wr.guardedUsagesComputed = false
			wr.enqueueContinuation(RemoveUnusedCode_Continuation{node: script, scope: globalScope})
			for wr.worklistSize > 0 {
				cont := wr.popContinuation()
				cont.Apply(wr)
			}
			wr.RemoveUnreferencedVarsAndPolyfills()
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

func (r *RemoveUnusedCode) isParamVarRemovableReadOnly(nameNode rhino.Node, fparamScope *Scope) bool {
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

func (r *RemoveUnusedCode) removeUnreferencedFunctionArgsWithCompiler(
	compiler AbstractCompiler,
	fparamScope *Scope,
) {
	function, ok := ast.AsFunction(fparamScope.GetRootNode())
	if !ok {
		jutil.CheckState(false)
	}
	if NodeUtil_IsGetOrSetKey(function.Node().GetParent()) {
		return
	}
	params := function.Params()
	argList := params.Node()
	for {
		lastArg, ok := params.LastParam()
		if !ok {
			break
		}
		argNode, removable := r.trailingParamBinding(lastArg)
		if !removable {
			break
		}
		if argNode.IsDestructuringPattern() {
			if argNode.HasChildren() {
				break
			}
			NodeUtil_DeleteNode(lastArg.Node(), compiler)
			continue
		}
		if r.isParamVarRemovableReadOnly(argNode, fparamScope) {
			NodeUtil_DeleteNode(lastArg.Node(), compiler)
		} else {
			break
		}
	}

	for param := range params.Params() {
		if param.Node().IsUnusedParameter() {
			continue
		}
		paramNameNode := RemoveUnusedCode_NameOfParam(param.Node())
		if paramNameNode.IsNil() {
			continue
		}
		if r.isParamVarRemovableReadOnly(paramNameNode, fparamScope) {
			param.Node().SetUnusedParameter(true)
			compiler.ReportChangeToEnclosingScope(argList)
		}
	}
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
				initRecordingCompiler(rc, r.compiler)
				r.removeUnreferencedFunctionArgsWithCompiler(rc, r.allFunctionParamScopes[origIdx])
				eventsByIdx[origIdx] = rc.takeEvents()
			}
		}
		return runItem, nil
	}, false)

	if len(detachedIndices) > 0 {
		var mainRC recordingCompiler
		for _, origIdx := range detachedIndices {
			initRecordingCompiler(&mainRC, r.compiler)
			r.removeUnreferencedFunctionArgsWithCompiler(&mainRC, r.allFunctionParamScopes[origIdx])
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
		r.guardedUsages = jutil.NewLinkedHashSet[rhino.Node]()
	}
	r.guardedUsages.Add(polyfillUsage.Node())
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

func (r *RemoveUnusedCode) traverseNameReference(n rhino.Node, scope *Scope) {
	// The only cases that should reach this point are parameter declarations and references
	// to names. The name node does not have children in these cases.
	if n.HasChildren() {
		jutil.CheckState(false)
	}
	// (A switch on the parent's token: this runs for every name, and
	// RoleOf would read the parent's children too.)
	switch n.GetParent().GetToken() {
	case rhino.Token_PARAM_LIST:
		// the parameter declaration is not a read of the name
		return
	case rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST:
		// var|let|const name;
		// are handled at a higher level.
		jutil.CheckState(false)
	case rhino.Token_FUNCTION, rhino.Token_CLASS:
		// function name() {}
		// class name() {}
		// are handled at a higher level.
		if role := ast.RoleOf(n); role == ast.RoleFunctionName || role == ast.RoleClassName {
			jutil.CheckState(false)
		}
	}
	varInfo := r.TraverseNameNode(n, scope)
	if varInfo.IsRemovable() {
		varInfo.pinAt(rucPinReferenceFound, n)
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
func (r *RemoveUnusedCode) TraverseNode(n rhino.Node, scope *Scope) {
	if n.GetFirstChild().IsNil() {
		if n.GetToken() == rhino.Token_NAME {
			r.traverseNameReference(n, scope)
		}
		return
	}
	switch n.GetToken() {
	case rhino.Token_CATCH:
		r.TraverseCatch(n, scope)
	case rhino.Token_FUNCTION:
		// If this function is a removable var, then create a continuation
		// for it instead of traversing immediately.
		if NodeUtil_IsFunctionDeclaration(n) {
			name, _ := ast.MustFunction(n).Name() // a declaration has one
			varInfo := r.TraverseNameNode(name.Node(), scope)
			varInfo.AddRemovable(r.newFunctionDeclRemovable(n, scope))
			if n.GetParent().IsExport() {
				varInfo.pin(rucPinExportedClass)
			}
		} else {
			r.TraverseFunction(n, scope)
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
	case rhino.Token_SWITCH_BODY:
		if nodeUtilBlockScopeMayHaveDecls(n, n) {
			r.TraverseChildren(n, r.scopeCreator.CreateScopeOfScope(n, scope))
		} else {
			r.TraverseChildren(n, scope)
		}
	case rhino.Token_BLOCK:
		// This case if for if there are let and const variables in block scopes.
		// Otherwise other variables will be hoisted up into the global scope and already be
		// handled.
		if scanVarsBlockCreatesScope(n) &&
			(NodeUtil_IsFunctionBlock(n) || NodeUtil_IsClassStaticBlock(n) || syntacticScopeHasBlockDecls(n)) {
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
		if !NodeUtil_IsStatement(n) {
			jutil.CheckState(false)
		}
		r.TraverseDeclarationStatement(n, scope)
	case rhino.Token_INSTANCEOF:
		r.TraverseInstanceof(n, scope)
	case rhino.Token_NAME:
		r.traverseNameReference(n, scope)
	case rhino.Token_GETPROP, rhino.Token_OPTCHAIN_GETPROP:
		r.TraverseNormalOrOptChainGetProp(n, scope)
	default:
		r.TraverseChildren(n, scope)
	}
}

// Java: RemoveUnusedCode.traverseInstanceof(Node, Scope)
func (r *RemoveUnusedCode) TraverseInstanceof(instanceofNode rhino.Node, scope *Scope) {
	instanceof, ok := ast.AsBinaryOf(instanceofNode, rhino.Token_INSTANCEOF)
	if !ok {
		jutil.CheckArgumentMsg1(false, "%s", instanceofNode)
	}
	r.TraverseNode(instanceof.Left().Node(), scope)
	if rhs, ok := instanceof.Right().AsName(); ok {
		varInfo := r.TraverseNameNode(rhs.Node(), scope)
		varInfo.AddRemovable(r.newInstanceofNameRemovable(instanceofNode))
	} else {
		r.TraverseNode(instanceof.Right().Node(), scope)
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
func (r *RemoveUnusedCode) TraverseNormalOrOptChainGetProp(getProp rhino.Node, scope *Scope) {
	access, ok := ast.AsNormalOrOptChainGetProp(getProp)
	if !ok {
		jutil.CheckStateMsg1(false, "%s", getProp)
	}
	objectNode := access.Object().Node()
	propID := getProp.GetStringID()

	if r.polyBits.Get(int(propID)) {
		r.considerPolyfillRef(propID, getProp)
	}

	if NodeUtil_IsExpressionResultUsed(getProp) ||
		r.ConsiderForAccessorSideEffects(getProp, AccessorSummary_PropertyAccessKind_GETTER_ONLY) {
		// must record as reference to the property and continue traversal.
		r.markPropertyIDPinned(propID)
		r.TraverseNode(objectNode, scope)
	} else if objectNode.IsThis() {
		// This is probably the declaration of a class field in a constructor.
		// /** @private {number} */
		// this.propName;
		// We don't want to consider this a real usage that should prevent removal.
		r.ConsiderForIndependentRemoval(r.newUnusedReadRefRemovable(getProp, remFlagThisDotPropRef))
	} else if owner, ok := removeUnusedCode_prototypeOwner(objectNode); ok {
		// (objExpression).prototype.propName;
		objExpression := owner.Node()
		if objExpression.IsName() {
			// name.prototype.propName;
			varInfo := r.TraverseNameNode(objExpression, scope)
			varInfo.AddRemovable(r.newUnusedReadRefRemovable(getProp, remFlagProtoDotPropRef))
		} else {
			// (objExpression).prototype.propName;
			rem := r.newUnusedReadRefRemovable(getProp, remFlagProtoDotPropRef)
			if r.MayHaveSideEffects(objExpression) {
				r.TraverseNode(objExpression, scope)
			} else {
				rem.addContinuation(r, objExpression, scope)
			}
			r.ConsiderForIndependentRemoval(rem)
		}
	} else {
		// TODO(bradfordcsmith): add removal of `varName.propName;`
		r.markPropertyIDPinned(propID)
		r.TraverseNode(objectNode, scope)
	}
}

// TraverseIncrementOrDecrementOp.
//
// TODO(b/137380742): Combine with `traverseCompoundAssign`.
//
// Java: RemoveUnusedCode.traverseIncrementOrDecrementOp(Node, Scope)
func (r *RemoveUnusedCode) TraverseIncrementOrDecrementOp(incOrDecOp rhino.Node, scope *Scope) {
	incDec, ok := ast.AsIncDec(incOrDecOp)
	jutil.CheckArgumentMsg1(ok, "%s", incOrDecOp)
	arg := incDec.Target().Node()
	if NodeUtil_IsExpressionResultUsed(incOrDecOp) {
		// If expression result is used, then this expression is definitely not removable.
		r.TraverseNode(arg, scope)
	} else if getProp, ok := incDec.Target().AsGetProp(); ok {
		getPropObj := getProp.Object().Node()

		if r.ConsiderForAccessorSideEffects(arg, AccessorSummary_PropertyAccessKind_GETTER_AND_SETTER) {
			r.TraverseNode(getPropObj, scope) // Don't re-traverse the GETPROP as a read.
		} else if getPropObj.IsThis() {
			// this.propName++
			r.ConsiderForIndependentRemoval(r.newIncOrDecOpRemovable(incOrDecOp, arg, rhino.NilNode, remFlagThisDotPropRef))
		} else if owner, ok := removeUnusedCode_prototypeOwner(getPropObj); ok {
			// someExpression.prototype.propName++
			exprObj := owner.Node()
			if exprObj.IsName() {
				// varName.prototype.propName++
				varInfo := r.TraverseNameNode(exprObj, scope)
				varInfo.AddRemovable(r.newIncOrDecOpRemovable(incOrDecOp, arg, rhino.NilNode, remFlagProtoDotPropRef))
			} else {
				// (someExpression).prototype.propName++
				var toPreserve rhino.Node
				var rem *RemoveUnusedCode_RemovableBase
				if r.MayHaveSideEffects(exprObj) {
					toPreserve = exprObj
					r.TraverseNode(exprObj, scope)
					rem = r.newIncOrDecOpRemovable(incOrDecOp, arg, toPreserve, remFlagProtoDotPropRef)
				} else {
					rem = r.newIncOrDecOpRemovable(incOrDecOp, arg, rhino.NilNode, remFlagProtoDotPropRef)
					rem.addContinuation(r, exprObj, scope)
				}
				r.ConsiderForIndependentRemoval(rem)
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
func (r *RemoveUnusedCode) TraverseCompoundAssign(compoundAssignNode rhino.Node, scope *Scope) {
	// We'll allow removal of compound assignment to a `this` property as long as the result of the
	// assignment is unused.
	// e.g. `this.prop += 3;`

	// NOTE: Some history here, as there were questions about "why is 'this' special".  The "remove
	// unused properties" is not a general property removal algorithm.  It only removes unreferenced
	// properties that are part of class definitions.  "SomeClass.prop += 3" and
	// "SomeClass.prototype.prop += 3" could so be candidates but they aren't considered here.

	assign := ast.MustAssignOp(compoundAssignNode)
	targetNode := assign.Target().Node()
	valueNode := assign.Value().Node()
	if getProp, ok := assign.Target().AsGetProp(); ok {
		if r.ConsiderForAccessorSideEffects(targetNode, AccessorSummary_PropertyAccessKind_GETTER_AND_SETTER) {
			r.TraverseNode(getProp.Object().Node(), scope) // Don't re-traverse the GETPROP as a read.
			r.TraverseNode(valueNode, scope)
		} else if getProp.Object().Token() == rhino.Token_THIS &&
			!NodeUtil_IsExpressionResultUsed(compoundAssignNode) {
			rem := r.newNamedPropAssignRemovable(compoundAssignNode, targetNode, nil, remFlagThisDotPropRef)
			r.traverseRemovableAssignValue(rem, valueNode, scope)
			r.ConsiderForIndependentRemoval(rem)
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
func (r *RemoveUnusedCode) TraverseNameNode(n rhino.Node, scope *Scope) *RemoveUnusedCode_RealVarInfo {
	nameID := n.GetStringID()
	if r.polyBits.Get(int(nameID)) {
		r.considerPolyfillRef(nameID, n)
	}
	v := scope.GetVarByID(nameID)
	if v == nil {
		jutil.CheckNotNullMsg1(v, "%s", n)
	}
	return r.TraverseVar(v)
}

func isSpecialCallCandidate(callee rhino.Node) bool {
	switch callee.GetToken() {
	case rhino.Token_GETPROP:
		switch callee.GetStringID() {
		case rucCalleeInheritsID, rucCalleeMixinID, rucCalleeDefinePropertiesID,
			rucCalleeObjectPropertyID, rucCalleeAddSingletonGetterID, rucCalleePolyfillID,
			rucCalleePatchID, rucCalleePolyfillTypedArrayMethodID, rucCalleeWeakUsageID:
			return true
		}
		return false
	case rhino.Token_NAME:
		switch callee.GetStringID() {
		case rucCalleeJscPropNameFnID,
			rucCalleeGoogReflectObjectPropertyID,
			rucCalleeJscompPolyfillID,
			rucCalleeJscompPatchID,
			rucCalleeJscompPolyfillTypedArrayMethodID,
			rucCalleeGoogWeakUsageID,
			rucCalleeGoogAddSingletonGetterID:
			return true
		}
		name := callee.GetString()
		return strings.HasSuffix(name, "$inherits") || strings.HasSuffix(name, "$mixin")
	default:
		return false
	}
}

// Java: RemoveUnusedCode.traverseCall(Node, Scope)
func (r *RemoveUnusedCode) TraverseCall(callNode rhino.Node, scope *Scope) {
	call := ast.MustNormalOrOptChainCall(callNode)
	callee := call.Callee().Node()
	if !isSpecialCallCandidate(callee) {
		r.TraverseChildren(callNode, scope)
		return
	}

	if r.codingConvention.IsPropertyRenameFunction(callee) {
		if firstArg, ok := call.FirstArg(); ok && firstArg.Node().IsStringLit() {
			r.markPropertyIDPinned(firstArg.Node().GetStringID())
		}
		r.TraverseChildren(callNode, scope)
	} else if NodeUtil_IsObjectDefinePropertiesDefinition(callNode) {
		// TODO(bradfordcsmith): Should also handle Object.create() and Object.defineProperty().
		r.TraverseObjectDefinePropertiesCall(callNode, scope)
	} else if r.removeUnusedPolyfills && RemoveUnusedCode_IsJscompPolyfill(callee) {
		firstArg, _ := call.FirstArg() // isJscompPolyfill: a string
		polyfillName := firstArg.MustStringLit().Value()
		if calleeID := callee.GetStringID(); calleeID == rucCalleePolyfillTypedArrayMethodID ||
			calleeID == rucCalleeJscompPolyfillTypedArrayMethodID {
			polyfillName = "TypedArray.prototype." + polyfillName
		}
		info := r.createPolyfillInfo(call, scope, polyfillName)
		r.polyBits.Set(int(info.keyID))
		idx := r.ensurePropEntry(info.keyID)
		entry := &r.propEntries[idx]
		if entry.polyHead == nil {
			entry.polyHead = info
			entry.polyTail = info
		} else {
			entry.polyTail.nextByKey = info
			entry.polyTail = info
		}
		r.polyfillOrder = append(r.polyfillOrder, info)
		// Only traverse the callee (to mark it as used).  The arguments may be traversed later.
		r.TraverseNode(callee, scope)
	} else if weakName, ok := removeUnusedCode_onlyNameArg(call); ok && NodeUtil_IsGoogWeakUsageCall(callNode) {
		// goog.weakUsage() should have exactly one argument, and it should be either a name or a
		// qualified name (this is checked in ProcessClosurePrimitives.java).
		//
		// If it is a qualified name, then we do not attempt to remove it at this time (see condition
		// above). We rely on CollapseProperties to turn qualified names into simple names where
		// possible.

		// Mark this call as removable if the var is not otherwise referenced.
		varInfo := r.TraverseNameNode(weakName.Node(), scope)
		varInfo.AddRemovable(r.newWeakUsageCallRemovable(callNode))

		// We need to traverse the goog.weakUsage function itself (to mark it as used, in case our
		// usage of it is not removed).
		r.TraverseNode(callee, scope)
	} else {
		parent := callNode.GetParent()
		var classVarName *string
		classDefiningCall := false

		// A call that is a statement unto itself or the left side of a comma expression might be
		// a call to a known method for doing class setup
		// e.g. $jscomp.inherits(Class, BaseClass) or goog.addSingletonGetter(Class)
		// Such methods never have meaningful return values, so we won't look for them in other
		// contexts
		if parent.IsExprResult() || (parent.IsComma() && ast.RoleOf(callNode) == ast.RoleBinaryLeft) {
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
			classVar = scope.GetVar(*classVarName)
			if classVar == nil {
				jutil.CheckNotNullMsg1(classVar, "%s", *classVarName)
			}
		}

		if classVar == nil || !classVar.IsGlobal() {
			// The call we are traversing does not modify a class definition,
			// or the class is not specified with a simple variable name,
			// or the variable name is not global.
			// TODO(bradfordcsmith): It would be more correct to check whether the class name
			// references a known constructor and expand to allow QNames.
			r.TraverseChildren(callNode, scope)
		} else {
			rem := r.newClassSetupCallRemovable(callNode, classDefiningCall)
			for child := callNode.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
				rem.addContinuation(r, child, scope)
			}
			r.TraverseVar(classVar).AddRemovable(rem)
		}
	}
}

// removeUnusedCode_onlyNameArg returns the argument of a call that has
// exactly one, a name.
func removeUnusedCode_onlyNameArg(call ast.NormalOrOptChainCall) (ast.Name, bool) {
	arg, ok := call.OnlyArg()
	if !ok {
		return ast.Name{}, false
	}
	return arg.AsName()
}

// RemoveUnusedCode_IsJscompPolyfill checks whether this is a recognizable
// call to $jscomp.polyfill.
//
// Java: RemoveUnusedCode.isJscompPolyfill(Node)
func RemoveUnusedCode_IsJscompPolyfill(n rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_NAME:
		// Need to work correctly after CollapseProperties.
		id := n.GetStringID()
		return (id == rucCalleeJscompPolyfillID ||
			id == rucCalleeJscompPatchID ||
			id == rucCalleeJscompPolyfillTypedArrayMethodID) &&
			n.GetNext().IsStringLit()
	case rhino.Token_GETPROP:
		// Need to work correctly without CollapseProperties.
		id := n.GetStringID()
		if id != rucCalleePolyfillID && id != rucCalleePatchID && id != rucCalleePolyfillTypedArrayMethodID {
			return false
		}
		object, ok := ast.MustGetProp(n).Object().AsName()
		return ok && object.Node().GetStringID() == rucJscompIdentID && n.GetNext().IsStringLit()
	default:
		return false
	}
}

// TraverseObjectDefinePropertiesCall traverses
// `Object.defineProperties(someObject, propertyDefinitions);`.
//
// Java: RemoveUnusedCode.traverseObjectDefinePropertiesCall(Node, Scope)
func (r *RemoveUnusedCode) TraverseObjectDefinePropertiesCall(callNode rhino.Node, scope *Scope) {
	// First child is Object.defineProperties or some equivalent of it.
	call := ast.MustNormalOrOptChainCall(callNode)
	callee := call.Callee().Node()
	target, _ := call.FirstArg() // isObjectDefinePropertiesDefinition: two arguments
	definitions, _ := target.Next()
	targetObject, propertyDefinitions := target.Node(), definitions.Node()

	// The name that the target object is or is the prototype of.
	nameNode, ok := target.AsName()
	if !ok {
		nameNode, ok = removeUnusedCode_nameDotPrototype(targetObject)
	}
	if ok && !NodeUtil_IsExpressionResultUsed(callNode) {
		// NOTE: Object.defineProperties() returns its first argument, so if its return value is used
		// that counts as a use of the targetObject.
		varInfo := r.TraverseNameNode(nameNode.Node(), scope)
		rem := r.newClassSetupCallRemovable(callNode, false)
		// TODO(bradfordcsmith): Is it really necessary to traverse the callee
		// (aka. Object.defineProperties)?
		rem.addContinuation(r, callee, scope)
		if r.MayHaveSideEffects(propertyDefinitions) {
			r.TraverseNode(propertyDefinitions, scope)
		} else {
			rem.addContinuation(r, propertyDefinitions, scope)
		}
		varInfo.AddRemovable(rem)
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
func (r *RemoveUnusedCode) TraverseObjectDefinePropertiesLiteral(propertyDefinitions rhino.Node, scope *Scope) {
	for member := range ast.MustObjectLit(propertyDefinitions).Members() {
		property := member.Node()
		if key, ok := member.AsStringKey(); ok && key.IsQuoted() {
			// Quoted property name counts as a reference to the property and protects it from removal.
			r.markPropertyIDPinned(property.GetStringID())
			r.TraverseNode(key.Value().Node(), scope)
		} else if ok {
			definition := key.Value().Node()
			if r.MayHaveSideEffects(definition) {
				r.TraverseNode(definition, scope)
			} else {
				rem := r.newObjectDefinePropsDefRemovable(property)
				rem.addContinuation(r, definition, scope)
				r.ConsiderForIndependentRemoval(rem)
			}
		} else {
			// TODO(bradfordcsmith): Maybe report error for anything other than a computed property,
			// since getters, setters, and methods don't make much sense in this context.
			r.TraverseNode(property, scope)
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

// Java: RemoveUnusedCode.traverseObjectLiteral(Node, Scope)
func (r *RemoveUnusedCode) TraverseObjectLiteral(objectLiteral rhino.Node, scope *Scope) {
	if !objectLiteral.IsObjectLit() {
		jutil.CheckArgumentMsg1(false, "%s", objectLiteral)
	}
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
func (r *RemoveUnusedCode) IsObjectDefinePropertiesSecondArgument(n rhino.Node) bool {
	call, ok := ast.AsCall(n.GetParent())
	if !ok {
		return false
	}
	last, ok := call.LastArg()
	return ok && last.Is(n) && NodeUtil_IsObjectDefinePropertiesDefinition(call.Node())
}

// Java: RemoveUnusedCode.traverseNonPrototypeObjectLiteral(Node, Scope)
func (r *RemoveUnusedCode) TraverseNonPrototypeObjectLiteral(objectLiteral rhino.Node, scope *Scope) {
	for member := range ast.MustObjectLit(objectLiteral).Members() {
		if key, ok := member.AsStringKey(); ok {
			// A property name in an object literal counts as a reference,
			// because of some reflection patterns.
			// Note that we are intentionally treating both quoted and unquoted keys as
			// references.
			r.markPropertyIDPinned(member.Node().GetStringID())
			r.TraverseNode(key.Value().Node(), scope)
		} else {
			r.TraverseNode(member.Node(), scope)
		}
	}
}

// Java: RemoveUnusedCode.traversePrototypeLiteral(Node, Scope)
func (r *RemoveUnusedCode) TraversePrototypeLiteral(objectLiteral rhino.Node, scope *Scope) {
	for propertyNode := objectLiteral.GetFirstChild(); !propertyNode.IsNil(); propertyNode = propertyNode.GetNext() {
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
				rem := r.newClassOrProtoPropRemovable(propertyNode)
				rem.addContinuation(r, valueNode, scope)
				r.ConsiderForIndependentRemoval(rem)
			}
		}
	}
}

// Java: RemoveUnusedCode.isAssignmentToPrototype(Node)
func (r *RemoveUnusedCode) IsAssignmentToPrototype(n rhino.Node) bool {
	assign, ok := ast.AsAssign(n)
	return ok && RemoveUnusedCode_IsDotPrototype(assign.Target().Node())
}

// RemoveUnusedCode_IsDotPrototype: true for `someExpression.prototype`.
//
// Java: RemoveUnusedCode.isDotPrototype(Node)
func RemoveUnusedCode_IsDotPrototype(n rhino.Node) bool {
	_, ok := removeUnusedCode_prototypeOwner(n)
	return ok
}

// removeUnusedCode_prototypeOwner returns someExpression of n =
// someExpression.prototype (or ?.prototype); ok is false for any other n.
func removeUnusedCode_prototypeOwner(n rhino.Node) (ast.Expr, bool) {
	access, ok := ast.AsNormalOrOptChainGetProp(n)
	if !ok || n.GetStringID() != rucPropPrototypeID {
		return ast.Expr{}, false
	}
	return access.Object(), true
}

// removeUnusedCode_nameDotPrototype returns name of n = name.prototype; ok
// is false for any other n.
func removeUnusedCode_nameDotPrototype(n rhino.Node) (ast.Name, bool) {
	getProp, ok := ast.AsGetProp(n)
	if !ok || n.GetStringID() != rucPropPrototypeID {
		return ast.Name{}, false
	}
	return getProp.Object().AsName()
}

// Java: RemoveUnusedCode.traverseCatch(Node, Scope)
func (r *RemoveUnusedCode) TraverseCatch(catchNode rhino.Node, scope *Scope) {
	catch := ast.MustCatch(catchNode)
	// The binding is missing (an EMPTY node) in catch {} (ES2019).
	if param, ok := catch.Param(); ok {
		if exceptionName, ok := param.AsName(); ok {
			exceptionVarInfo := r.TraverseNameNode(exceptionName.Node(), scope)
			exceptionVarInfo.pin(rucPinCatchVariable)
		} else {
			r.TraverseNode(param.Node(), scope)
		}
	}
	r.TraverseNode(catch.Body().Node(), scope)
}

// Java: RemoveUnusedCode.traverseEnhancedFor(Node, Scope)
func (r *RemoveUnusedCode) TraverseEnhancedFor(enhancedFor rhino.Node, scope *Scope) {
	forScope := scope
	if syntacticScopeHasBlockDecls(enhancedFor) {
		forScope = r.scopeCreator.CreateScopeOfScope(enhancedFor, scope)
	}
	// for (iterationTarget in|of collection) body;
	loop := ast.MustEnhancedFor(enhancedFor)
	iterationTarget := loop.Target().Node()
	if name, ok := loop.Target().AsName(); ok {
		// using previously-declared loop variable. e.g.
		// `for (varName of collection) {}`
		varInfo := r.TraverseNameNode(name.Node(), forScope)
		varInfo.pin(rucPinEnhancedForVar)
	} else if decl, ok := loop.Target().AsNameDeclaration(); ok {
		// loop has const/var/let declaration
		d, ok := decl.OnlyDeclarator()
		jutil.CheckState(ok)
		declNode := d.Node()
		if name, ok := d.AsName(); !ok {
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
			_, hasInit := name.Init()
			jutil.CheckState(!hasInit)
			// We can never remove the loop variable of a for-in or for-of loop, because it's
			// essential to loop syntax.
			varInfo := r.TraverseNameNode(name.Node(), forScope)
			varInfo.pin(rucPinEnhancedForVar)
		}
	} else {
		// using some general LHS value e.g.
		// `for ([a, b] of collection) {}` destructuring with existing vars
		// `for (a.x of collection) {}` using a property as the loop var
		// TODO(bradfordcsmith): This should be considered a write if it's a property reference.
		r.TraverseNode(iterationTarget, forScope)
	}
	r.TraverseNode(loop.Collection().Node(), forScope)
	r.TraverseNode(loop.AsLoop().Body().Node(), forScope)
}

// Java: RemoveUnusedCode.traverseVanillaFor(Node, Scope)
func (r *RemoveUnusedCode) TraverseVanillaFor(forNode rhino.Node, scope *Scope) {
	forScope := scope
	if syntacticScopeHasBlockDecls(forNode) {
		forScope = r.scopeCreator.CreateScopeOfScope(forNode, scope)
	}
	// The missing parts of a for (an EMPTY node each) hold nothing to traverse.
	loop := ast.MustFor(forNode)
	if initialization, ok := loop.Init(); ok {
		if decl, ok := initialization.AsNameDeclaration(); ok {
			r.TraverseVanillaForNameDeclarations(decl.Node(), forScope)
		} else {
			r.TraverseNode(initialization.Node(), forScope)
		}
	}
	if condition, ok := loop.Cond(); ok {
		r.TraverseNode(condition.Node(), forScope)
	}
	if update, ok := loop.Update(); ok {
		r.TraverseNode(update.Node(), forScope)
	}
	r.TraverseNode(loop.Body().Node(), forScope)
}

// Java: RemoveUnusedCode.traverseVanillaForNameDeclarations(Node, Scope)
func (r *RemoveUnusedCode) TraverseVanillaForNameDeclarations(nameDeclaration rhino.Node, scope *Scope) {
	for child := range ast.MustNameDeclaration(nameDeclaration).Declarators() {
		if name, ok := child.AsName(); !ok {
			// TODO(bradfordcsmith): Customize handling of destructuring
			r.TraverseNode(child.Node(), scope)
		} else {
			nameNode := name.Node()
			value, hasValue := name.Init()
			valueNode := value.Node()
			varInfo := r.TraverseNameNode(nameNode, scope)
			if !hasValue {
				varInfo.AddRemovable(r.newVanillaForNameDeclRemovable(nameNode, false))
			} else if r.MayHaveSideEffects(valueNode) {
				// TODO(bradfordcsmith): Actually allow for removing the variable while keeping the
				// valueNode for its side-effects.
				varInfo.pin(rucPinForLoopInitSideEffects)
				r.TraverseNode(valueNode, scope)
			} else {
				rem := r.newVanillaForNameDeclRemovable(nameNode, true)
				rem.addContinuation(r, valueNode, scope)
				varInfo.AddRemovable(rem)
			}
		}
	}
}

// Java: RemoveUnusedCode.traverseDeclarationStatement(Node, Scope)
func (r *RemoveUnusedCode) TraverseDeclarationStatement(declarationStatement rhino.Node, scope *Scope) {
	// Normalization should ensure that declaration statements always have just one child.
	declarator, ok := ast.MustNameDeclaration(declarationStatement).OnlyDeclarator()
	jutil.CheckState(ok)
	if name, ok := declarator.AsName(); !ok {
		// Destructuring declarations are handled elsewhere.
		r.TraverseNode(declarator.Node(), scope)
	} else {
		value, hasValue := name.Init()
		valueNode := value.Node()
		varInfo := r.TraverseNameNode(name.Node(), scope)
		rem := r.newNameDeclStmtRemovable(declarationStatement, hasValue)
		if hasValue {
			if r.MayHaveSideEffects(valueNode) {
				r.TraverseNode(valueNode, scope)
			} else {
				rem.addContinuation(r, valueNode, scope)
			}
		}
		varInfo.AddRemovable(rem)
	}
}

// Java: RemoveUnusedCode.traverseAssign(Node, Scope)
func (r *RemoveUnusedCode) TraverseAssign(assignNode rhino.Node, scope *Scope) {
	assign, ok := ast.AsAssignment(assignNode)
	if !ok {
		jutil.CheckState(false)
	}

	lhs := assign.Target().Node()
	valueNode := assign.Value().Node()
	if _, ok := assign.Target().AsName(); ok {
		// varName = something
		varInfo := r.TraverseNameNode(lhs, scope)
		rem := r.newVarAssignRemovable(assignNode, varInfo)
		r.traverseRemovableAssignValue(rem, valueNode, scope)
		varInfo.AddRemovable(rem)
	} else if getElem, ok := assign.Target().AsGetElem(); ok {
		getElemObj := getElem.Object().Node()
		getElemKey := getElem.Index().Node()
		varName, ok := getElem.Object().AsName()
		if !ok {
			varName, ok = removeUnusedCode_nameDotPrototype(getElemObj)
		}

		if ok {
			// varName[someExpression] = someValue
			// OR
			// varName.prototype[someExpression] = someValue
			varInfo := r.TraverseNameNode(varName.Node(), scope)
			rem := r.newComputedPropAssignRemovable(assignNode, getElemKey, varInfo)
			if r.MayHaveSideEffects(getElemKey) {
				r.TraverseNode(getElemKey, scope)
			} else {
				rem.addContinuation(r, getElemKey, scope)
			}
			r.traverseRemovableAssignValue(rem, valueNode, scope)
			varInfo.AddRemovable(rem)
		} else {
			r.TraverseNode(getElemObj, scope)
			r.TraverseNode(getElemKey, scope)
			r.TraverseNode(valueNode, scope)
		}
	} else if getProp, ok := assign.Target().AsGetProp(); ok {
		getPropLhs := getProp.Object().Node()
		// Assignments `Foo.prototype.bar = function() {`
		owner, isDotPrototypeLhs := removeUnusedCode_prototypeOwner(getPropLhs)
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
			rem := r.newNamedPropAssignRemovable(assignNode, lhs, varInfo, 0)
			r.traverseRemovableAssignValue(rem, valueNode, scope)
			varInfo.AddRemovable(rem)
		} else if isDotPrototypeLhs {
			// objExpression.prototype.propertyName = someValue
			objExpression := owner.Node()
			if objExpression.IsName() {
				// varName.prototype.propertyName = someValue
				rem := r.newNamedPropAssignRemovable(assignNode, lhs, nil, remFlagProtoDotPropRef)
				r.traverseRemovableAssignValue(rem, valueNode, scope)
				varInfo := r.TraverseNameNode(objExpression, scope)
				rem.varInfo = varInfo
				varInfo.AddRemovable(rem)
			} else {
				// (someExpression).prototype.propertyName = someValue
				rem := r.newAnonProtoPropAssignRemovable(assignNode, lhs)
				r.traverseRemovableAssignValue(rem, valueNode, scope)
				if r.MayHaveSideEffects(objExpression) {
					r.TraverseNode(objExpression, scope)
				} else {
					rem.addContinuation(r, objExpression, scope)
				}
				r.ConsiderForIndependentRemoval(rem)
			}
		} else if getPropLhs.IsThis() {
			// this.propertyName = someValue
			rem := r.newNamedPropAssignRemovable(assignNode, lhs, nil, remFlagThisDotPropRef)
			r.traverseRemovableAssignValue(rem, valueNode, scope)
			r.ConsiderForIndependentRemoval(rem)
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

// Java: RemoveUnusedCode.isNameDotPrototype(Node)
func (r *RemoveUnusedCode) IsNameDotPrototype(n rhino.Node) bool {
	_, ok := removeUnusedCode_nameDotPrototype(n)
	return ok
}

// Java: RemoveUnusedCode.traverseObjectPattern(Node, Scope)
func (r *RemoveUnusedCode) TraverseObjectPattern(pattern rhino.Node, scope *Scope) {
	objectPattern, ok := ast.AsObjectPattern(pattern)
	if !ok {
		jutil.CheckStateMsg1(false, "%s", pattern)
	}

	for prop := range objectPattern.Props() {
		elem := prop.Node()
		switch prop.Token() {
		case rhino.Token_COMPUTED_PROP:
			value, _ := prop.MustComputedProp().Value() // always there in a pattern
			r.TraverseIndirectAssignment(elem, value.Node(), scope)
		case rhino.Token_STRING_KEY:
			key := prop.MustStringKey()
			if !key.IsQuoted() {
				r.markPropertyIDPinned(elem.GetStringID())
			}
			r.TraverseIndirectAssignment(elem, key.Value().Node(), scope)
		case rhino.Token_OBJECT_REST:
			// Recall that the rest target can be any l-value expression
			r.TraverseIndirectAssignment(elem, prop.MustObjectRest().Target().Node(), scope)
		default:
			panic(jutil.NewIllegalStateExceptionMsg(
				"Unexpected child of " + pattern.GetToken().String() + ": " + elem.ToStringTree()))
		}
	}
}

// Java: RemoveUnusedCode.traverseIndirectAssignmentList(Node, Scope)
func (r *RemoveUnusedCode) TraverseIndirectAssignmentList(list rhino.Node, scope *Scope) {
	if !list.IsArrayPattern() && !list.IsParamList() {
		jutil.CheckStateMsg1(false, "%s", list)
	}

	for elem := list.GetFirstChild(); !elem.IsNil(); elem = elem.GetNext() {
		switch elem.GetToken() {
		case rhino.Token_EMPTY:
		case rhino.Token_ARRAY_PATTERN, rhino.Token_DEFAULT_VALUE, rhino.Token_GETELEM, rhino.Token_GETPROP, rhino.Token_NAME, rhino.Token_OBJECT_PATTERN:
			r.TraverseIndirectAssignment(elem, elem, scope)
		case rhino.Token_ITER_REST:
			r.TraverseIndirectAssignment(elem, ast.MustIterRest(elem).Target().Node(), scope)
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
func (r *RemoveUnusedCode) TraverseIndirectAssignment(root rhino.Node, target rhino.Node, scope *Scope) {
	rootParent := root.GetParent()
	if !rootParent.IsDestructuringPattern() && !rootParent.IsParamList() {
		jutil.CheckArgumentMsg1(false, "%s", rootParent)
	}

	// Flatten out the case where the target is a default value. We always have to consider it.
	if defaultValue, ok := ast.AsDefaultValue(target); ok {
		target = defaultValue.Target().Node()
	}

	if target.IsGetProp() {
		r.ConsiderForAccessorSideEffects(target, AccessorSummary_PropertyAccessKind_SETTER_ONLY)
	}

	if r.MayHaveSideEffects(root) {
		// If anywhere in the assignment subtree has side-effects, it means that even if the target is
		// removable the subtree is not.
		r.TraverseNode(root, scope)
		// TODO(bradfordcsmith): Preserve side effects without preventing removal of variables and
		// properties. We could probably do this by subbing in an empty object pattern.
	} else if target.IsName() {
		varInfo := r.TraverseNameNode(target, scope)
		varInfo.AddRemovable(r.newIndirectAssignRemovable(root, target, scope))
	} else if r.IsNameDotPrototype(target) || RemoveUnusedCode_IsThisDotProperty(target) {
		r.ConsiderForIndependentRemoval(r.newIndirectAssignRemovable(root, target, scope))
	} else {
		// TODO(bradfordcsmith): Handle property assignments also
		// e.g. `({a: foo.bar, b: foo.baz}) = {a: 1, b: 2}`
		r.TraverseNode(root, scope)
	}
}

// Java: RemoveUnusedCode.traverseChildren(Node, Scope)
func (r *RemoveUnusedCode) TraverseChildren(n rhino.Node, scope *Scope) {
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		if c.GetFirstChild().IsNil() {
			if c.GetToken() == rhino.Token_NAME {
				r.traverseNameReference(c, scope)
			}
			continue
		}
		r.TraverseNode(c, scope)
	}
}

// TraverseClass handles a class that is not the RHS child of an assignment
// or a variable declaration initializer.
//
// Java: RemoveUnusedCode.traverseClass(Node, Scope)
func (r *RemoveUnusedCode) TraverseClass(classNode rhino.Node, scope *Scope) {
	jutil.CheckArgument(classNode.IsClass())
	if NodeUtil_IsClassDeclaration(classNode) {
		r.TraverseClassDeclaration(classNode, scope)
	} else {
		r.TraverseClassExpression(classNode, scope)
	}
}

// Java: RemoveUnusedCode.traverseClassDeclaration(Node, Scope)
func (r *RemoveUnusedCode) TraverseClassDeclaration(classNode rhino.Node, scope *Scope) {
	class, ok := ast.AsClass(classNode)
	jutil.CheckArgument(ok)
	className, _ := class.Name() // a declaration has one
	// A missing extends (an EMPTY node) has nothing to traverse.
	baseClass, hasBaseClass := class.Extends()
	classBodyNode := class.Body().Node()
	// Class declarations never declare the class name in the inner class scope (only class
	// expressions do), so creating an inner scope here would always produce an empty scope.
	classScope := scope

	varInfo := r.TraverseNameNode(className.Node(), scope)
	if classNode.GetParent().IsExport() {
		// Cannot remove an exported class.
		varInfo.pin(rucPinExportedClass)
		if hasBaseClass {
			r.TraverseNode(baseClass.Node(), scope)
		}
		// Use traverseChildren() here, because we should not consider any properties on the exported
		// class to be removable.
		r.TraverseChildren(classBodyNode, classScope)
	} else if hasBaseClass && r.MayHaveSideEffects(baseClass.Node()) {
		// TODO(bradfordcsmith): implement removal without losing side-effects for this case
		varInfo.pin(rucPinBaseClassSideEffects)
		r.TraverseNode(baseClass.Node(), scope)
		r.TraverseClassMembers(classBodyNode, classScope)
	} else if r.MayHaveSideEffects(classBodyNode) {
		varInfo.pin(rucPinClassBodySideEffects)
		if hasBaseClass {
			r.TraverseNode(baseClass.Node(), scope)
		}
		r.TraverseClassMembers(classBodyNode, classScope)
	} else {
		rem := r.newClassDeclRemovable(classNode)
		if hasBaseClass {
			rem.addContinuation(r, baseClass.Node(), classScope)
		}
		rem.addContinuation(r, classBodyNode, classScope)
		varInfo.AddRemovable(rem)
	}
}

// Java: RemoveUnusedCode.traverseClassExpression(Node, Scope)
func (r *RemoveUnusedCode) TraverseClassExpression(classNode rhino.Node, scope *Scope) {
	class, ok := ast.AsClass(classNode)
	jutil.CheckArgument(ok)
	classScope := scope

	if className, ok := class.Name(); ok {
		classScope = r.scopeCreator.CreateScopeOfScope(classNode, scope)
		// We may be able to remove the name node if nothing ends up referring to it.
		varInfo := r.TraverseNameNode(className.Node(), classScope)
		// The class is non-local, because it is accessible by unknown code outside
		// of the scope where InnerName is defined.
		// e.g. `use(class InnerName {})`
		varInfo.SetHasNonLocalOrNonLiteralValue()
		varInfo.AddRemovable(r.newNamedClassExprRemovable(classNode))
	}
	// If we're traversing the class expression, we've already decided we cannot remove it.
	if baseClass, ok := class.Extends(); ok {
		r.TraverseNode(baseClass.Node(), scope)
	}
	r.TraverseClassMembers(class.Body().Node(), classScope)
}

// Java: RemoveUnusedCode.traverseClassMembers(Node, Scope)
func (r *RemoveUnusedCode) TraverseClassMembers(node rhino.Node, scope *Scope) {
	jutil.CheckArgumentMsg1(node.IsClassMembers(), "%s", node)
	if !r.removeUnusedPrototypeProperties {
		r.TraverseChildren(node, scope)
		return
	}

	for member := node.GetFirstChild(); !member.IsNil(); member = member.GetNext() {
		switch member.GetToken() {
		case rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF, rhino.Token_MEMBER_FUNCTION_DEF:
			// If we get as far as traversing the members of a class, we've already decided that
			// we cannot remove the class itself, so just consider individual members for removal.
			rem := r.newClassOrProtoPropRemovable(member)
			rem.addContinuation(r, member, scope)
			r.ConsiderForIndependentRemoval(rem)
		case rhino.Token_MEMBER_FIELD_DEF:
			// TODO(bradfordcsmith): currently if the RHS of a field has side effects, we do not
			// remove any part of the field. The proper behavior of class C { x = alert(); }
			// would be to remove x, leaving class C { constructor() { alert(); } }
			// but currently we aren't removing anything.
			value, hasValue := ast.MustMemberFieldDef(member).Value()
			if !hasValue || !r.MayHaveSideEffects(value.Node()) {
				rem := r.newClassOrProtoPropRemovable(member)
				rem.addContinuation(r, member, scope)
				r.ConsiderForIndependentRemoval(rem)
			}
			// TODO: b/354704593 - remove the entire class when it is not referenced
			if hasValue {
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
func (r *RemoveUnusedCode) TraverseFunction(function rhino.Node, parentScope *Scope) {
	fn, ok := ast.AsFunction(function)
	if !ok {
		jutil.CheckStateMsg1(false, "%s", function)
	}
	body := fn.Body().Node()
	if !body.IsBlock() {
		jutil.CheckStateMsg1(false, "%s", body)
	}
	paramlist := fn.Params()

	// Checking the parameters
	fparamScope := r.scopeCreator.CreateScopeOfScope(function, parentScope)

	// Checking the function body
	fbodyScope := r.scopeCreator.CreateScopeOfScope(body, fparamScope)

	if name, ok := fn.Name(); ok {
		// var x = function funcName() {};
		// make sure funcName gets into the varInfoMap so it will be considered for removal.
		varInfo := r.TraverseNameNode(name.Node(), fparamScope)
		if NodeUtil_IsExpressionResultUsed(function) {
			// var f = function g() {};
			// The f is an alias for g, so g escapes from the scope where it is defined.
			varInfo.SetHasNonLocalOrNonLiteralValue()
		}
	}

	if paramlist.HasParams() {
		r.TraverseIndirectAssignmentList(paramlist.Node(), fparamScope)
		if r.removeGlobals {
			r.allFunctionParamScopes = append(r.allFunctionParamScopes, fparamScope)
		}
	}
	r.TraverseChildren(body, fbodyScope)
}

// Java: RemoveUnusedCode.canRemoveParameters(Node)
func (r *RemoveUnusedCode) CanRemoveParameters(parameterList rhino.Node) bool {
	if !parameterList.IsParamList() {
		jutil.CheckState(false)
	}
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
	if !function.IsFunction() {
		jutil.CheckState(false)
	}
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

// MarkUnusedParameters marks any remaining unused parameters as being unused
// so it can be used elsewhere.
//
// paramList is the list of function's parameters.
//
// Java: RemoveUnusedCode.markUnusedParameters(Node, Scope)
func (r *RemoveUnusedCode) MarkUnusedParameters(paramList rhino.Node, fparamScope *Scope) {
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

		varInfo := r.TraverseNameNode(paramNameNode, fparamScope)
		if varInfo.IsRemovable() {
			param.SetUnusedParameter(true)
			r.compiler.ReportChangeToEnclosingScope(paramList)
			if r.removalLog.IsLogging() {
				r.logRemoval("argmark", paramNameNode.GetString(), loggableFunctionName(paramList))
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
func (r *RemoveUnusedCode) MaybeRemoveUnusedTrailingParameters(argList rhino.Node, fparamScope *Scope) {
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
			NodeUtil_DeleteNode(lastArg, r.compiler)
			if r.removalLog.IsLogging() {
				r.logRemoval("arg", "<pattern>", loggableFunctionName(argList))
			}
			continue
		}

		varInfo := r.GetVarInfo(r.GetVarForNameNode(argNode, fparamScope))
		if varInfo.IsRemovable() {
			NodeUtil_DeleteNode(lastArg, r.compiler)
			if r.removalLog.IsLogging() {
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
