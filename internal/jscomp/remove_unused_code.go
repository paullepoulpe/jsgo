// Port of com/google/javascript/jscomp/RemoveUnusedCode.java @ closure-compiler v20261028

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
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

// Go-only: the traversal methods and the Removables read the shapes they
// handle through the typed view (internal/ast): the parts of assignments,
// declarations, loops, classes, functions, parameters and patterns, and
// the x.prototype and this.x forms. The generic subtree walks, the
// removals' tree edits and the parallel bookkeeping stay on rhino.Node.

type rucPropEntry struct {
	removables      rucRemovableQueue
	polyHead        *rucPolyfillInfo
	polyTail        *rucPolyfillInfo
	propID          strpool.ID
	checkedUnpinned bool
	inRemOrder      bool
}

// removeUnusedCodeState is the state of one run of the pass; Process drops
// it when the run ends.
type removeUnusedCodeState struct {
	worklist fifoQueue[RemoveUnusedCode_Continuation]

	// Java: varInfoMap (insertion-ordered Map<Var, VarInfo>, where Var equality is
	// (scope root, name)). varByKey indexes VarInfos by scopedVarKey and
	// varInfos preserves insertion order.
	varByKey map[scopedVarKey]*RemoveUnusedCode_RealVarInfo
	varInfos []*RemoveUnusedCode_RealVarInfo

	trackProps   bool
	pinnedProps  smallBitSet
	polyBits     smallBitSet
	propIdxByID  map[strpool.ID]uint32
	propEntries  []rucPropEntry
	propRemOrder []uint32

	// Stores all polyfills seen during traversal, indexed by leaf property ID in
	// propEntries (gated by polyBits) and ordered by polyfillOrder.
	polyfillOrder         []*rucPolyfillInfo
	guardedUsages         map[rhino.Node]bool
	guardedUsagesComputed bool
	// Go-only: the kept guarded usages, when the run uses them
	// (remove_unused_code_incremental.go).
	guardedCache            *rucGuardedCache
	externs                 rhino.Node
	root                    rhino.Node
	suppressPolyfillRemoval bool

	// Go-only: per-run function summary recording and replay state
	// (remove_unused_code_summary.go).
	fnCache               *rucFnSummaryCache
	dirtySubtreeRoots     map[rhino.Node]struct{}
	recState              *rucRecorderState
	activeGlobalRec       *rucRecordingSummary
	activeLocalRec        *rucRecordingSummary
	deferredSummaryFns    []rucDeferredFn
	replayVars            []*Var
	savedPolyBits         smallBitSet
	replayedFns           map[rhino.Node]struct{}
	replayedCount         int
	deferredCount         int
	polyEpoch             uint32
	propPinEpoch          uint32
	recMode               rucRecordMode
	summaryReplayDisabled bool

	allFunctionParamScopes []*Scope

	extraContSlab []rucExtraCont

	// Go-only: GetVarInfo may read and write Var.rucInfo (not on parallel
	// workers, which share the global scope's vars, nor while logging).
	cacheVarInfos bool

	// Go-only: performance — one run's removables, var infos, and polyfill infos
	// are created together and recycled per script in parallel mode.
	removables   jutil.ReuseSlab[rucRemovable]
	varInfoSlab  jutil.ReuseSlab[RemoveUnusedCode_RealVarInfo]
	polyfillSlab jutil.ReuseSlab[rucPolyfillInfo]
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

	externProperties *ExternPropSet
	accessorSummary  *AccessorSummary

	// Single value to use for all vars for which we cannot remove anything at
	// all.
	canonicalUnremovableVarInfo    *RemoveUnusedCode_RealVarInfo
	canonicalUnremovableVarInfoVal RemoveUnusedCode_RealVarInfo

	polyfillsFromTable  *PolyfillUsageFinder_Polyfills
	polyfillUsageFinder *PolyfillUsageFinder

	scopeCreator *SyntacticScopeCreator

	removeUnusedPrototypeProperties               bool
	removeUnusedThisProperties                    bool
	removeUnusedObjectDefinePropertiesDefinitions bool
	removeUnusedPolyfills                         bool
	assumeGettersArePure                          bool
	disambiguateProperties                        bool

	// Allocated & cleaned up by process()
	removalLog     diagnostic.LogFile
	unremovableLog diagnostic.LogFile

	removeUnusedCodeState
}

var (
	// Properties that are implicitly used as part of the JS language.
	rucPropPrototypeID = uint32(idPrototype)

	rucImplicitPropIDs = [5]uint32{
		uint32(idLength),
		uint32(idToString),
		uint32(idValueOf),
		uint32(idConstructor),
		uint32(idPrototype),
	}

	rucCalleeInheritsID                 = uint32(strpool.Intern("inherits"))
	rucCalleeMixinID                    = uint32(strpool.Intern("mixin"))
	rucCalleeDefinePropertiesID         = uint32(idDefineProperties)
	rucCalleeObjectPropertyID           = uint32(strpool.Intern("objectProperty"))
	rucCalleeAddSingletonGetterID       = uint32(strpool.Intern("addSingletonGetter"))
	rucCalleePolyfillID                 = uint32(strpool.Intern("polyfill"))
	rucCalleePatchID                    = uint32(strpool.Intern("patch"))
	rucCalleePolyfillTypedArrayMethodID = uint32(strpool.Intern("polyfillTypedArrayMethod"))
	rucCalleeWeakUsageID                = uint32(idWeakUsage)

	rucCalleeJscPropNameFnID                  = uint32(idJscPropertyNameFn)
	rucCalleeGoogReflectObjectPropertyID      = uint32(strpool.Intern("goog$reflect$objectProperty"))
	rucCalleeJscompPolyfillID                 = uint32(strpool.Intern("$jscomp$polyfill"))
	rucCalleeJscompPatchID                    = uint32(strpool.Intern("$jscomp$patch"))
	rucCalleeJscompPolyfillTypedArrayMethodID = uint32(strpool.Intern("$jscomp$polyfillTypedArrayMethod"))
	rucCalleeGoogWeakUsageID                  = uint32(idGoogWeakUsage)
	rucCalleeGoogAddSingletonGetterID         = uint32(strpool.Intern("goog$addSingletonGetter"))
	rucJscompIdentID                          = uint32(strpool.Intern("$jscomp"))
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
	DisambiguateProperties                        bool
}

// Java: RemoveUnusedCode(Builder)
func NewRemoveUnusedCode(compiler AbstractCompiler, opts RemoveUnusedCodeOptions) *RemoveUnusedCode {
	sc := NewSyntacticScopeCreator(compiler)
	polyfillsFromTable := defaultPolyfillsTable()
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
		disambiguateProperties:                        opts.DisambiguateProperties,
		polyfillsFromTable:                            polyfillsFromTable,
		polyfillUsageFinder:                           NewPolyfillUsageFinder(compiler, polyfillsFromTable, opts.DisambiguateProperties),
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
	if functionName, ok := NodeUtil_NearestFunctionName(rhino.CheckNotNull(argList).GetParent()); ok {
		return functionName
	}
	return "<anonymous>"
}

func (r *RemoveUnusedCode) removesAnyProperties() bool {
	return r.removeUnusedPrototypeProperties ||
		r.removeUnusedThisProperties ||
		r.removeUnusedObjectDefinePropertiesDefinitions
}

func (r *RemoveUnusedCode) removesPolyfills() bool {
	return r.removeUnusedPolyfills && !r.suppressPolyfillRemoval
}

func (r *RemoveUnusedCode) ensurePropEntry(id strpool.ID) uint32 {
	if idx1 := r.propIdxByID[id]; idx1 != 0 {
		return idx1 - 1
	}
	if r.propIdxByID == nil {
		r.propIdxByID = make(map[strpool.ID]uint32, 64)
	}
	idx := uint32(len(r.propEntries))
	r.propEntries = append(r.propEntries, rucPropEntry{propID: id})
	r.propIdxByID[id] = idx + 1
	return idx
}

func (r *RemoveUnusedCode) considerPolyfillRef(keyID uint32, n rhino.Node) {
	if r.activeGlobalRec != nil {
		r.activeGlobalRec.addPolyNode(n)
	}
	idx1 := r.propIdxByID[strpool.ID(keyID)]
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
	jutil.CheckState(r.compiler.GetLifeCycleStage().IsNormalized())
	r.externProperties = r.compiler.GetExternPropSet()
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
	r.externs = externs
	r.TraverseAndRemoveUnusedReferences(root)
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
	if r.root.IsNil() {
		r.root = node
	}
	r.guardedUsagesComputed = false

	r.worklist.push(RemoveUnusedCode_Continuation{node: node, scope: scope})
	r.drainWorklistAndDeferred(scope)

	if rucIncrVerify && r.replayedCount > 0 {
		r.verifySummaryReplay(node, scope)
	}
	if r.recMode != rucRecNone {
		r.savedPolyBits.w0 = r.polyBits.w0
		r.savedPolyBits.w1 = r.polyBits.w1
		if len(r.polyBits.ext) > 0 {
			r.savedPolyBits.ext = append(r.savedPolyBits.ext[:0], r.polyBits.ext...)
		} else {
			r.savedPolyBits.ext = r.savedPolyBits.ext[:0]
		}
	}

	r.RemoveUnreferencedVarsAndPolyfills()
}

// TraverseAndRemoveUnusedReferences traverses a node recursively. Call this
// once per pass.
//
// Java: RemoveUnusedCode.traverseAndRemoveUnusedReferences(Node)
func (r *RemoveUnusedCode) TraverseAndRemoveUnusedReferences(root rhino.Node) {
	r.root = root
	r.initPolyfillSuppression(root)
	if r.beginSummaryRun(root) {
		return
	}
	handled, scope := r.traverseAndRemoveLocalScriptsParallel(root)
	if handled {
		return
	}
	r.cacheVarInfos = r.compiler != nil && r.compiler.parallelRecorder == nil &&
		!r.removalLog.IsLogging() && !r.unremovableLog.IsLogging()
	r.trackProps = r.removesAnyProperties()
	if r.trackProps {
		for _, id := range rucImplicitPropIDs {
			r.pinnedProps.Set(int(id))
		}
	}
	if scope == nil {
		scope = r.createGlobalScope(root)
	}
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
	r.finishSummaryRun(nil)
}

func (r *RemoveUnusedCode) canParallelizeLocalScripts() bool {
	return !r.removeGlobals &&
		!r.removeUnusedPrototypeProperties &&
		!r.removeUnusedThisProperties &&
		!r.removeUnusedObjectDefinePropertiesDefinitions &&
		!r.removesPolyfills() &&
		!r.removalLog.IsLogging() &&
		!r.unremovableLog.IsLogging()
}

func (r *RemoveUnusedCode) newLocalScriptWorker(rc *Compiler) *RemoveUnusedCode {
	wr := &RemoveUnusedCode{
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
		polyfillUsageFinder:             r.polyfillUsageFinder,
		assumeGettersArePure:            r.assumeGettersArePure,
		disambiguateProperties:          r.disambiguateProperties,
		removalLog:                      r.removalLog,
		unremovableLog:                  r.unremovableLog,
	}
	wr.varByKey = make(map[scopedVarKey]*RemoveUnusedCode_RealVarInfo, 128)
	wr.canonicalUnremovableVarInfoVal.outer = wr
	wr.canonicalUnremovableVarInfoVal.flags = rucVarFlagCanonical
	wr.canonicalUnremovableVarInfo = &wr.canonicalUnremovableVarInfoVal
	return wr
}

func (r *RemoveUnusedCode) resetLocalScriptWorker() {
	clear(r.varByKey)
	clear(r.varInfos)
	r.varInfos = r.varInfos[:0]
	r.worklist.clear()
	clear(r.extraContSlab)
	r.extraContSlab = r.extraContSlab[:0]
	clear(r.replayVars)
	r.replayVars = r.replayVars[:0]
	r.removables.Reset(nil)
	r.varInfoSlab.Reset(nil)
	r.polyfillSlab.Reset(nil)
	r.activeLocalRec = nil
	r.activeGlobalRec = nil
}

type rucLocalUnitStatus uint8

const (
	rucUnitsOK rucLocalUnitStatus = iota
	rucUnitsWholeScript
	rucUnitsFallbackSerial
)

func collectGlobalLocalScopeUnits(n rhino.Node, globalScope *Scope, units *[]rhino.Node) rucLocalUnitStatus {
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		if c.GetFirstChild().IsNil() {
			if c.GetToken() == rhino.Token_NAME && globalScope.GetVarByID(c.GetStringID()) == nil {
				return rucUnitsFallbackSerial
			}
			continue
		}
		switch c.GetToken() {
		case rhino.Token_MODULE_BODY:
			return rucUnitsWholeScript
		case rhino.Token_NAME:
			if globalScope.GetVarByID(c.GetStringID()) == nil {
				return rucUnitsFallbackSerial
			}
			if st := collectGlobalLocalScopeUnits(c, globalScope, units); st != rucUnitsOK {
				return st
			}
		case rhino.Token_FUNCTION, rhino.Token_CLASS,
			rhino.Token_FOR, rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF,
			rhino.Token_SWITCH_BODY:
			*units = append(*units, c)
		case rhino.Token_BLOCK:
			if scanVarsBlockCreatesScope(c) && (NodeUtil_IsClassStaticBlock(c) || syntacticScopeHasBlockDecls(c)) {
				*units = append(*units, c)
			} else if st := collectGlobalLocalScopeUnits(c, globalScope, units); st != rucUnitsOK {
				return st
			}
		default:
			if st := collectGlobalLocalScopeUnits(c, globalScope, units); st != rucUnitsOK {
				return st
			}
		}
	}
	return rucUnitsOK
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
		panic(jutil.NewNullPointerExceptionMsg(jutil.LenientFormat("%s", nameNode)))
	}
	if r.cacheVarInfos {
		if c := v.rucInfo; c != 0 {
			if c == rucVarInfoCanonical {
				return false
			}
			if i := int(c - 1); i < len(r.varInfos) && r.varInfos[i].var_ == v {
				return r.varInfos[i].IsRemovable()
			}
		}
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
	jutil.CheckState(ok)
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

func (r *RemoveUnusedCode) removeUnreferencedFunctionArgsParallel() bool {
	if !r.removeGlobals || r.removalLog.IsLogging() || r.unremovableLog.IsLogging() {
		return false
	}
	scopes := r.allFunctionParamScopes
	return runParallelNodesByScript(r.compiler, len(scopes), parallelNodeTask[struct{}]{
		allowDetached:           true,
		enableTypeRegistryReads: true,
		prepare: func() {
			for _, s := range scopes {
				if s != nil {
					freezeGlobalScopeForParallelReads(s.GetGlobalScope())
					break
				}
			}
		},
		nodeAt: func(idx int) rhino.Node { return scopes[idx].GetRootNode() },
		runItem: func(_ struct{}, idx int, rc *Compiler) {
			r.removeUnreferencedFunctionArgsWithCompiler(rc, scopes[idx], true)
		},
	})
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
		entry := &r.propEntries[idx]
		if entry.removables.isEmpty() {
			continue
		}
		if r.removalLog.IsLogging() {
			r.logRemoval("prop", strpool.String(entry.propID), "")
		}
		entry.removables.drain(func(cur *rucRemovable) {
			cur.flags &^= remFlagInPropList
			cur.Remove(r)
		})
	}
}

// Java: RemoveUnusedCode.getVarForNameNode(Node, Scope)
func (r *RemoveUnusedCode) GetVarForNameNode(nameNode rhino.Node, scope *Scope) *Var {
	v := scope.GetVarByNode(nameNode)
	if v == nil {
		panic(jutil.NewNullPointerExceptionMsg(jutil.LenientFormat("%s", nameNode)))
	}
	return v
}

// Java: RemoveUnusedCode.canRemoveParameters(Node)
func (r *RemoveUnusedCode) CanRemoveParameters(parameterList rhino.Node) bool {
	params, ok := ast.AsParamList(parameterList)
	jutil.CheckState(ok)
	function := params.Node().GetParent()
	return r.removeGlobals && !NodeUtil_IsGetOrSetKey(function.GetParent())
}

func (r *RemoveUnusedCode) markPropertyIDPinned(id uint32) {
	if !r.trackProps {
		return
	}
	if r.activeGlobalRec != nil && id != 0 {
		r.activeGlobalRec.addPropID(id)
	}
	if r.pinnedProps.Get(int(id)) {
		return
	}
	r.pinnedProps.Set(int(id))
	r.propPinEpoch++
	// Continue traversal of all of the property name's values and no longer consider them for
	// removal.
	if idx1 := r.propIdxByID[strpool.ID(id)]; idx1 != 0 {
		r.propEntries[idx1-1].removables.drain(func(cur *rucRemovable) {
			cur.flags &^= remFlagInPropList
			cur.ApplyContinuations(r)
		})
	}
}

func (r *RemoveUnusedCode) isExportedByConvention(name string, local bool) bool {
	if len(name) == 0 || (name[0] != '_' && name[0] != '$') {
		return false
	}
	return r.codingConvention.IsExportedWithLocal(name, local)
}

// Java: RemoveUnusedCode.considerForIndependentRemoval(Removable)
func (r *RemoveUnusedCode) ConsiderForIndependentRemoval(removable *rucRemovable) {
	propID := removable.propertyNameID
	if !r.trackProps || propID == 0 {
		removable.ApplyContinuations(r)
		return
	}
	indep := r.IsIndependentlyRemovable(removable)
	if r.activeGlobalRec != nil {
		if indep {
			r.activeGlobalRec.addIndepPropID(propID)
		} else {
			r.activeGlobalRec.addPropID(propID)
		}
	}
	if r.pinnedProps.Get(int(propID)) {
		// Referenced or exported, so not removable.
		removable.ApplyContinuations(r)
		return
	}
	idx := r.ensurePropEntry(strpool.ID(propID))
	if !r.propEntries[idx].checkedUnpinned {
		if r.externProperties.HasID(strpool.ID(propID)) ||
			r.isExportedByConvention(strpool.String(strpool.ID(propID)) /* local= */, false) {
			r.pinnedProps.Set(int(propID))
			removable.ApplyContinuations(r)
			return
		}
		r.propEntries[idx].checkedUnpinned = true
	}
	if indep {
		// Store for possible removal later.
		if removable.flags&remFlagInPropList == 0 {
			removable.flags |= remFlagInPropList
			entry := &r.propEntries[idx]
			if !entry.inRemOrder {
				entry.inRemOrder = true
				r.propRemOrder = append(r.propRemOrder, idx)
			}
			entry.removables.push(removable)
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
	access, ok := ast.AsNormalOrOptChainGetProp(getprop)
	jutil.CheckStateMsg1(ok, "%s", getprop)
	if r.assumeGettersArePure && !usage.HasSetter() {
		return false
	}
	summary := r.accessorSummary
	if summary == nil {
		summary = r.compiler.GetAccessorSummary()
		r.accessorSummary = summary
	}
	propID := access.Node().GetStringID()
	recorded := summary.GetKindByID(propID)
	if (recorded.HasGetter() && usage.HasGetter() && !r.assumeGettersArePure) ||
		(recorded.HasSetter() && usage.HasSetter()) {
		r.markPropertyIDPinned(propID)
		return true
	}

	return false
}

// Java: RemoveUnusedCode.isIndependentlyRemovable(Removable)
func (r *RemoveUnusedCode) IsIndependentlyRemovable(removable *rucRemovable) bool {
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
	jutil.CheckArgumentMsg1(ok, "%s", paramList)

	for p := range params.Params() {
		if p.IsUnused() {
			// already marked
			continue
		}

		paramName, hasName := rucParamName(p)
		if !hasName {
			// destructuring pattern parameters don't have a name that applies to the whole parameter
			// TODO(bradfordcsmith): We could mark this if we determined that all vars created by
			// the pattern are unused.
			continue
		}

		if r.isParamVarRemovable(paramName.Node(), fparamScope, readOnly) {
			p.SetUnused(true)
			compiler.ReportChangeToEnclosingScope(paramList)
			if !readOnly && r.removalLog.IsLogging() {
				r.logRemoval("argmark", paramName.Ident(), loggableFunctionName(paramList))
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
	jutil.CheckArgumentMsg1(ok, "%s", argList)
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

		if pat, isPattern := ast.AsPattern(argNode); isPattern {
			if pat.Node().HasChildren() {
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
				r.logRemoval("arg", ast.MustName(argNode).Ident(), loggableFunctionName(argList))
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
		if r.recMode == rucRecLocal {
			r.onVarPinned(var_)
		}
		// If we are considering removing local variables, that includes parameters.
		// If `arguments` is used in a function we must consider all parameters to be referenced.
		functionScope := var_.scope.GetClosestHoistScope()
		for param := range ast.MustFunction(functionScope.GetRootNode()).Params().Params() {
			paramName, ok := rucParamName(param)
			if !ok {
				continue
			}

			r.GetVarInfo(r.GetVarForNameNode(paramName.Node(), functionScope)).
				pin(rucPinParamInFnUsingArguments)
		}
		// `arguments` is never removable.
		return r.canonicalUnremovableVarInfo
	}
	return r.GetVarInfo(var_)
}

func rucParamName(p ast.Param) (ast.Name, bool) {
	n := RemoveUnusedCode_NameOfParam(p.Node())
	if n.IsNil() {
		return ast.Name{}, false
	}
	return ast.MustName(n), true
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
	if !r.trackProps && !r.unremovableLog.IsLogging() &&
		((isGlobal && !r.removeGlobals) || (!isGlobal && !r.removeLocalVars)) {
		return r.canonicalUnremovableVarInfo
	}
	if r.cacheVarInfos {
		if c := var_.rucInfo; c != 0 {
			if c == rucVarInfoCanonical {
				return r.canonicalUnremovableVarInfo
			}
			if i := int(c - 1); i < len(r.varInfos) && r.varInfos[i].var_ == var_ {
				return r.varInfos[i]
			}
		}
	}
	key := scopedVarKeyOf(var_)
	if v := r.varByKey[key]; v != nil {
		return v
	}
	// (These three depend on the var alone: the same answer in every run.)
	if var_.IsExtern() {
		if r.unremovableLog.IsLogging() {
			r.unremovableLog.LogString(var_.GetName() + ": extern")
		}
		return r.cachedCanonicalVarInfo(var_)
	} else if name := var_.GetName(); r.isExportedByConvention(name /* local= */, !isGlobal) {
		if r.unremovableLog.IsLogging() {
			r.unremovableLog.LogString(name + ": exported by convention")
		}
		return r.cachedCanonicalVarInfo(var_)
	} else if var_.IsArguments() {
		// No point in logging that we cannot remove "arguments"
		return r.cachedCanonicalVarInfo(var_)
	}
	varInfo := NewRemoveUnusedCode_RealVarInfo(r, var_)
	if _, isParamList := ast.AsParamList(var_.GetParentNode()); isParamList {
		varInfo.flags |= rucVarFlagHasNonLocalOrNonLit
	}
	// When trackProps is enabled, non-removable vars still need their own VarInfo
	// to track assigned values for static property removal.
	if !r.removeGlobals && isGlobal {
		varInfo.pin(rucPinNotRemovingGlobals)
	} else if !r.removeLocalVars && !isGlobal {
		varInfo.pin(rucPinNotRemovingLocals)
	}
	r.varByKey[key] = varInfo
	r.varInfos = append(r.varInfos, varInfo)
	if r.cacheVarInfos {
		var_.rucInfo = uint32(len(r.varInfos))
	}
	return varInfo
}

// rucVarInfoCanonical is Var.rucInfo for a var whose var info is always the
// canonical unremovable one.
const rucVarInfoCanonical = ^uint32(0)

func (r *RemoveUnusedCode) cachedCanonicalVarInfo(var_ *Var) *RemoveUnusedCode_RealVarInfo {
	if r.cacheVarInfos {
		var_.rucInfo = rucVarInfoCanonical
	}
	return r.canonicalUnremovableVarInfo
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
		} else if fn, isFn := ast.AsFunction(toRemove); isFn && NodeUtil_IsFunctionExpression(fn.Node()) {
			// TODO(bradfordcsmith): Add a Removable for this case.
			if !r.preserveFunctionExpressionNames {
				r.compiler.ReportChangeToEnclosingScope(nameNode)
				fn.RemoveName()
			}
		} else {
			// Removables are not created for theses cases.
			// function foo(unused1 = someSideEffectingValue, ...unused2) {}
			// removeUnreferencedFunctionArgs() is responsible for removing these.
			// TODO(bradfordcsmith): handle parameter declarations with removables
			_, isParamList := ast.AsParamList(toRemove)
			_, parentIsParamList := ast.AsParamList(toRemove.GetParent())
			_, isDefaultVal := ast.AsDefaultValue(toRemove)
			_, isRest := ast.AsRest(toRemove)
			jutil.CheckStateMsg1(
				isParamList ||
					(parentIsParamList && (isDefaultVal || isRest)),
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
