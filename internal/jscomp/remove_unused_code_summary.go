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

package jscomp

import (
	"fmt"
	"slices"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

const rucSummaryTimelineKey = "__removeUnusedCode_fnSummary__"

type rucRecordMode uint8

const (
	rucRecNone rucRecordMode = iota
	// rucRecGlobal records and replays top-level functions (parentScope.IsGlobal())
	// when removing globals, properties, or polyfills.
	rucRecGlobal
	// rucRecLocal records and replays both top-level and nested functions when
	// only local variables are being removed (canParallelizeLocalScripts()).
	rucRecLocal
)

type rucFnSummaryFlags uint8

const (
	rucSumFlagValid rucFnSummaryFlags = 1 << iota
	rucSumFlagIsGlobalParent
	rucSumFlagIsGetOrSetKey
	rucSumFlagNamedFnExpr
	rucSumFlagExprResultUsed
)

const rucContextFlagMask = rucSumFlagIsGetOrSetKey | rucSumFlagNamedFnExpr | rucSumFlagExprResultUsed

// rucFnSummary is a 36-byte pointer-free summary of a single FUNCTION subtree.
// All variable/property IDs and AST node handles live in contiguous arenas on
// rucFnSummaryCache so Go's GC never scans rucFnSummary map buckets.
type rucFnSummary struct {
	varOff   uint32
	varLen   uint32
	propOff  uint32
	propLen  uint32
	indepOff uint32
	indepLen uint32
	polyOff  uint32
	polyLen  uint32
	flags    rucFnSummaryFlags
}

func (s rucFnSummary) isValid() bool {
	return s.flags&rucSumFlagValid != 0
}

func (s rucFnSummary) isGlobalParent() bool {
	return s.flags&rucSumFlagIsGlobalParent != 0
}

func (s rucFnSummary) totalElements() uint32 {
	return s.varLen + s.propLen + s.indepLen + s.polyLen
}

type rucDeferredFn struct {
	fn  rhino.Node
	sum rucFnSummary
}

func computeFnContextFlags(fn rhino.Node) rucFnSummaryFlags {
	var flags rucFnSummaryFlags
	if NodeUtil_IsGetOrSetKey(fn.GetParent()) {
		flags |= rucSumFlagIsGetOrSetKey
	}
	if _, hasName := ast.MustFunction(fn).Name(); hasName && NodeUtil_IsFunctionExpression(fn) {
		flags |= rucSumFlagNamedFnExpr
		if NodeUtil_IsExpressionResultUsed(fn) {
			flags |= rucSumFlagExprResultUsed
		}
	}
	return flags
}

func (s rucFnSummary) matchesContext(fn rhino.Node) bool {
	return (s.flags & rucContextFlagMask) == computeFnContextFlags(fn)
}

// rucFnSummaryCache persists clean function summaries across RemoveUnusedCode
// iterations within the normalized unobfuscated optimization stage.
type rucFnSummaryCache struct {
	root               rhino.Node
	externsRoot        rhino.Node
	externsScriptCount int
	jsScriptCount      int
	compacts           int
	recMode            rucRecordMode
	opts               RemoveUnusedCodeOptions
	hasAstAnalyzer     bool
	astAnalyzerVal     AstAnalyzer
	externProps        *ExternPropSet
	accessorSummary    *AccessorSummary
	polyfills          *PolyfillUsageFinder_Polyfills
	colorRegistry      *colors.ColorRegistry

	recordedPolyBits smallBitSet

	byFn           map[rhino.Node]rucFnSummary
	varIDArena     []uint32
	propIDArena    []uint32
	indepPropArena []uint32
	polyNodeArena  []rhino.Node
	deadElements   uint32

	completed        bool
	hadChangesDuring bool

	// Diagnostic counters updated at the end of each run for unit tests.
	lastReplayed        int
	lastDeferred        int
	lastIdempotentSkip  bool
	idempotentSkipCount int
}

func (fc *rucFnSummaryCache) clearAll() {
	clear(fc.byFn)
	fc.varIDArena = fc.varIDArena[:0]
	fc.propIDArena = fc.propIDArena[:0]
	fc.indepPropArena = fc.indepPropArena[:0]
	fc.polyNodeArena = fc.polyNodeArena[:0]
	fc.deadElements = 0
	fc.recordedPolyBits.ClearAll()
}

func (fc *rucFnSummaryCache) maybeCompactArenas() {
	liveTotal := len(fc.varIDArena) + len(fc.propIDArena) + len(fc.indepPropArena) + len(fc.polyNodeArena)
	if fc.deadElements <= 1024 || int(fc.deadElements)*2 < liveTotal {
		return
	}
	newVars := make([]uint32, 0, len(fc.varIDArena))
	newProps := make([]uint32, 0, len(fc.propIDArena))
	newIndeps := make([]uint32, 0, len(fc.indepPropArena))
	newPolys := make([]rhino.Node, 0, len(fc.polyNodeArena))

	// order-independent: each function entry receives its own disjoint subslice with identical contents
	for fn, s := range fc.byFn {
		if fn.IsDeleted() {
			delete(fc.byFn, fn)
			continue
		}
		var ns rucFnSummary
		ns.flags = s.flags
		if s.varLen > 0 {
			ns.varOff = uint32(len(newVars))
			ns.varLen = s.varLen
			newVars = append(newVars, fc.varIDArena[s.varOff:s.varOff+s.varLen]...)
		}
		if s.propLen > 0 {
			ns.propOff = uint32(len(newProps))
			ns.propLen = s.propLen
			newProps = append(newProps, fc.propIDArena[s.propOff:s.propOff+s.propLen]...)
		}
		if s.indepLen > 0 {
			ns.indepOff = uint32(len(newIndeps))
			ns.indepLen = s.indepLen
			newIndeps = append(newIndeps, fc.indepPropArena[s.indepOff:s.indepOff+s.indepLen]...)
		}
		if s.polyLen > 0 {
			ns.polyOff = uint32(len(newPolys))
			ns.polyLen = s.polyLen
			newPolys = append(newPolys, fc.polyNodeArena[s.polyOff:s.polyOff+s.polyLen]...)
		}
		fc.byFn[fn] = ns
	}
	fc.varIDArena = newVars
	fc.propIDArena = newProps
	fc.indepPropArena = newIndeps
	fc.polyNodeArena = newPolys
	fc.deadElements = 0
}

type rucRecordingSummary struct {
	fn             rhino.Node
	parentRec      *rucRecordingSummary
	parentDepth    int32
	startPolyEpoch uint32
	flags          rucFnSummaryFlags
	ineligible     bool

	varIDs       []uint32
	varSet       map[uint32]struct{}
	propIDs      []uint32
	propSet      map[uint32]struct{}
	indepPropIDs []uint32
	indepPropSet map[uint32]struct{}
	polyNodes    []rhino.Node
	polySet      map[rhino.Node]struct{}
}

func appendUniqueUint32(slice *[]uint32, set *map[uint32]struct{}, v uint32) {
	s := *slice
	if len(s) < 8 {
		for _, cur := range s {
			if cur == v {
				return
			}
		}
		*slice = append(s, v)
		return
	}
	m := *set
	if m == nil {
		m = make(map[uint32]struct{}, 16)
		for _, cur := range s {
			m[cur] = struct{}{}
		}
		*set = m
	}
	if _, exists := m[v]; exists {
		return
	}
	m[v] = struct{}{}
	*slice = append(s, v)
}

func (rec *rucRecordingSummary) addVarID(id uint32) {
	if !rec.ineligible && id != 0 {
		appendUniqueUint32(&rec.varIDs, &rec.varSet, id)
	}
}

func (rec *rucRecordingSummary) addPropID(id uint32) {
	if !rec.ineligible && id != 0 {
		appendUniqueUint32(&rec.propIDs, &rec.propSet, id)
	}
}

func (rec *rucRecordingSummary) addIndepPropID(id uint32) {
	if !rec.ineligible && id != 0 {
		appendUniqueUint32(&rec.indepPropIDs, &rec.indepPropSet, id)
	}
}

func (rec *rucRecordingSummary) addPolyNode(n rhino.Node) {
	if rec.ineligible || n.IsNil() {
		return
	}
	s := rec.polyNodes
	if len(s) < 8 {
		for _, cur := range s {
			if cur == n {
				return
			}
		}
		rec.polyNodes = append(s, n)
		return
	}
	m := rec.polySet
	if m == nil {
		m = make(map[rhino.Node]struct{}, 16)
		for _, cur := range s {
			m[cur] = struct{}{}
		}
		rec.polySet = m
	}
	if _, exists := m[n]; exists {
		return
	}
	m[n] = struct{}{}
	rec.polyNodes = append(s, n)
}

type rucRecorderState struct {
	byFn  map[rhino.Node]*rucRecordingSummary
	order []*rucRecordingSummary
}

func newRucRecorderState() *rucRecorderState {
	return &rucRecorderState{
		byFn: make(map[rhino.Node]*rucRecordingSummary, 64),
	}
}

func (r *RemoveUnusedCode) optionsSnapshot() RemoveUnusedCodeOptions {
	return RemoveUnusedCodeOptions{
		RemoveLocalVars:                               r.removeLocalVars,
		RemoveGlobals:                                 r.removeGlobals,
		PreserveFunctionExpressionNames:               r.preserveFunctionExpressionNames,
		RemoveUnusedPrototypeProperties:               r.removeUnusedPrototypeProperties,
		RemoveUnusedThisProperties:                    r.removeUnusedThisProperties,
		RemoveUnusedObjectDefinePropertiesDefinitions: r.removeUnusedObjectDefinePropertiesDefinitions,
		RemoveUnusedPolyfills:                         r.removesPolyfills(),
		AssumeGettersArePure:                          r.assumeGettersArePure,
		DisambiguateProperties:                        r.disambiguateProperties,
	}
}

func (r *RemoveUnusedCode) resolveExternsRoot(root rhino.Node) rhino.Node {
	if !r.externs.IsNil() {
		return r.externs
	}
	if parent := root.GetParent(); !parent.IsNil() {
		if first := parent.GetFirstChild(); first != root {
			return first
		}
	}
	return rhino.NilNode
}

// beginSummaryRun prepares incremental function summary state for root and
// reports whether the entire pass can exit immediately because zero AST scopes
// changed since the end of a compatible, zero-change previous run.
func (r *RemoveUnusedCode) beginSummaryRun(root rhino.Node) bool {
	r.root = root
	if r.removalLog.IsLogging() || r.unremovableLog.IsLogging() {
		return false
	}
	c := r.incrCompiler()
	if c == nil {
		return false
	}
	if r.accessorSummary == nil {
		r.accessorSummary = c.GetAccessorSummary()
	}
	if r.canParallelizeLocalScripts() {
		r.recMode = rucRecLocal
	} else {
		r.recMode = rucRecGlobal
	}

	changedBefore := c.changeTracker.GetChangedScopeNodesForPass(rucSummaryTimelineKey)
	g := c.rucGuarded
	if g == nil {
		g = &rucGuardedCache{}
		c.rucGuarded = g
	}

	externsRoot := r.resolveExternsRoot(root)
	opts := r.optionsSnapshot()
	hasAnalyzer := r.astAnalyzer != nil
	var analyzerVal AstAnalyzer
	if hasAnalyzer {
		analyzerVal = *r.astAnalyzer
	}
	var colorReg *colors.ColorRegistry
	if r.polyfillUsageFinder != nil {
		colorReg = r.polyfillUsageFinder.registry
	}

	fc := g.fnCache
	compatible := fc != nil &&
		changedBefore != nil &&
		fc.root == root &&
		fc.externsRoot == externsRoot &&
		fc.compacts == c.astCompact.count &&
		fc.recMode == r.recMode &&
		fc.opts == opts &&
		fc.hasAstAnalyzer == hasAnalyzer &&
		(!hasAnalyzer || fc.astAnalyzerVal == analyzerVal) &&
		fc.externProps == r.externProperties &&
		fc.accessorSummary == r.accessorSummary &&
		fc.polyfills == r.polyfillsFromTable &&
		fc.colorRegistry == colorReg

	if compatible && fc.completed && !fc.hadChangesDuring && len(changedBefore) == 0 &&
		fc.externsScriptCount == countNodeChildren(externsRoot) &&
		fc.jsScriptCount == countNodeChildren(root) {
		fc.lastReplayed = 0
		fc.lastDeferred = 0
		fc.lastIdempotentSkip = true
		fc.idempotentSkipCount++
		g.fnCache = fc
		if rucIncrVerify {
			r.verifyIdempotentSkip(root)
		}
		return true
	}

	if compatible && len(changedBefore) > 0 && !externsRoot.IsNil() {
		for _, n := range changedBefore {
			if !rucScriptUnder(n, externsRoot).IsNil() {
				compatible = false
				break
			}
		}
	}

	var dirty map[rhino.Node]struct{}
	if !compatible {
		fc = &rucFnSummaryCache{
			root:            root,
			externsRoot:     externsRoot,
			compacts:        c.astCompact.count,
			recMode:         r.recMode,
			opts:            opts,
			hasAstAnalyzer:  hasAnalyzer,
			astAnalyzerVal:  analyzerVal,
			externProps:     r.externProperties,
			accessorSummary: r.accessorSummary,
			polyfills:       r.polyfillsFromTable,
			colorRegistry:   colorReg,
			byFn:            make(map[rhino.Node]rucFnSummary, 64),
		}
	} else if len(changedBefore) > 0 {
		dirty = make(map[rhino.Node]struct{}, len(changedBefore)*2)
		for _, n := range changedBefore {
			if n.IsFunction() {
				if old, ok := fc.byFn[n]; ok {
					fc.deadElements += old.totalElements()
					delete(fc.byFn, n)
				}
			}
			for cur := n; !cur.IsNil(); cur = cur.GetParent() {
				if cur.IsFunction() || cur.IsScript() {
					dirty[cur] = struct{}{}
					if cur.IsFunction() {
						if old, ok := fc.byFn[cur]; ok {
							fc.deadElements += old.totalElements()
							delete(fc.byFn, cur)
						}
					}
				}
			}
		}
	}

	// Detach fnCache during execution so any panic leaves g.fnCache == nil.
	fc.completed = false
	g.fnCache = nil
	r.fnCache = fc
	r.dirtySubtreeRoots = dirty
	r.recState = newRucRecorderState()
	return false
}

func (r *RemoveUnusedCode) finishSummaryRun(extraRecs []*rucRecorderState) {
	fc := r.fnCache
	if fc == nil {
		return
	}
	c := r.incrCompiler()
	if c == nil {
		return
	}
	if r.accessorSummary == nil {
		r.accessorSummary = c.GetAccessorSummary()
	}

	changedDuring := c.changeTracker.GetChangedScopeNodesForPass(rucSummaryTimelineKey)
	var changedFns map[rhino.Node]struct{}
	if len(changedDuring) > 0 {
		changedFns = make(map[rhino.Node]struct{}, len(changedDuring)*2)
		for _, n := range changedDuring {
			if n.IsFunction() {
				changedFns[n] = struct{}{}
				if old, ok := fc.byFn[n]; ok {
					fc.deadElements += old.totalElements()
					delete(fc.byFn, n)
				}
			}
			for cur := n; !cur.IsNil(); cur = cur.GetParent() {
				if cur.IsFunction() {
					changedFns[cur] = struct{}{}
					if old, ok := fc.byFn[cur]; ok {
						fc.deadElements += old.totalElements()
						delete(fc.byFn, cur)
					}
				}
			}
		}
	}

	if r.summaryReplayDisabled {
		fc.clearAll()
	}

	commitRecorder := func(rs *rucRecorderState) {
		if rs == nil {
			return
		}
		for _, rec := range rs.order {
			fn := rec.fn
			if rec.ineligible || fn.IsDeleted() || rec.startPolyEpoch != r.polyEpoch {
				if old, ok := fc.byFn[fn]; ok {
					fc.deadElements += old.totalElements()
					delete(fc.byFn, fn)
				}
				continue
			}
			if _, modified := changedFns[fn]; modified {
				continue
			}
			if RemoveUnusedCode_AlreadyRemoved(fn) {
				continue
			}
			if old, ok := fc.byFn[fn]; ok {
				fc.deadElements += old.totalElements()
			}
			var s rucFnSummary
			s.flags = rec.flags | rucSumFlagValid
			if n := len(rec.varIDs); n > 0 {
				s.varOff = uint32(len(fc.varIDArena))
				s.varLen = uint32(n)
				fc.varIDArena = append(fc.varIDArena, rec.varIDs...)
			}
			if n := len(rec.propIDs); n > 0 {
				s.propOff = uint32(len(fc.propIDArena))
				s.propLen = uint32(n)
				fc.propIDArena = append(fc.propIDArena, rec.propIDs...)
			}
			if n := len(rec.indepPropIDs); n > 0 {
				s.indepOff = uint32(len(fc.indepPropArena))
				s.indepLen = uint32(n)
				fc.indepPropArena = append(fc.indepPropArena, rec.indepPropIDs...)
			}
			if n := len(rec.polyNodes); n > 0 {
				s.polyOff = uint32(len(fc.polyNodeArena))
				s.polyLen = uint32(n)
				fc.polyNodeArena = append(fc.polyNodeArena, rec.polyNodes...)
			}
			fc.byFn[fn] = s
		}
	}

	commitRecorder(r.recState)
	for _, rs := range extraRecs {
		commitRecorder(rs)
	}

	fc.recordedPolyBits.w0 = r.savedPolyBits.w0
	fc.recordedPolyBits.w1 = r.savedPolyBits.w1
	fc.recordedPolyBits.ext = slices.Clone(r.savedPolyBits.ext)
	fc.externsScriptCount = countNodeChildren(fc.externsRoot)
	fc.jsScriptCount = countNodeChildren(fc.root)
	fc.hadChangesDuring = len(changedDuring) > 0
	fc.completed = true
	fc.lastReplayed = r.replayedCount
	fc.lastDeferred = r.deferredCount
	fc.lastIdempotentSkip = false
	fc.maybeCompactArenas()

	g := c.rucGuarded
	if g == nil {
		g = &rucGuardedCache{}
		c.rucGuarded = g
	}
	g.fnCache = fc
}

func (r *RemoveUnusedCode) isDirtySubtree(n rhino.Node) bool {
	if r.dirtySubtreeRoots == nil {
		return false
	}
	_, dirty := r.dirtySubtreeRoots[n]
	return dirty
}

func (r *RemoveUnusedCode) resolveGlobalRec(scope *Scope) *rucRecordingSummary {
	if scope == nil || scope.depth < 1 || r.recState == nil {
		return nil
	}
	s := scope
	for s.depth > 1 {
		s = s.parent
	}
	if !s.rootNode.IsFunction() {
		return nil
	}
	return r.recState.byFn[s.rootNode]
}

func (r *RemoveUnusedCode) resolveLocalRec(scope *Scope) *rucRecordingSummary {
	if scope == nil || scope.depth < 1 || r.recState == nil {
		return nil
	}
	for s := scope; s != nil && s.depth >= 1; s = s.parent {
		if s.rootNode.IsFunction() {
			return r.recState.byFn[s.rootNode]
		}
	}
	return nil
}

func (r *RemoveUnusedCode) beginFnRecord(fn rhino.Node, parentScope *Scope) *rucRecordingSummary {
	if r.recState == nil {
		return nil
	}
	flags := computeFnContextFlags(fn)
	if parentScope.IsGlobal() {
		flags |= rucSumFlagIsGlobalParent
	}
	rec := &rucRecordingSummary{
		fn:             fn,
		parentDepth:    parentScope.depth,
		startPolyEpoch: r.polyEpoch,
		flags:          flags,
	}
	if r.recMode == rucRecLocal {
		rec.parentRec = r.activeLocalRec
	}
	r.recState.byFn[fn] = rec
	r.recState.order = append(r.recState.order, rec)
	return rec
}

func (r *RemoveUnusedCode) onVarPinned(var_ *Var) {
	if var_ == nil {
		return
	}
	switch r.recMode {
	case rucRecGlobal:
		if r.activeGlobalRec != nil && var_.IsGlobal() {
			r.activeGlobalRec.addVarID(var_.nameID)
		}
	case rucRecLocal:
		vDepth := var_.scope.depth
		if vDepth > 0 {
			for rec := r.activeLocalRec; rec != nil && vDepth <= rec.parentDepth; rec = rec.parentRec {
				rec.addVarID(var_.nameID)
			}
		}
	}
}

func (r *RemoveUnusedCode) onVarMutated(var_ *Var) {
	if var_ == nil {
		return
	}
	switch r.recMode {
	case rucRecGlobal:
		if r.activeGlobalRec != nil && var_.IsGlobal() {
			r.activeGlobalRec.ineligible = true
		}
	case rucRecLocal:
		vDepth := var_.scope.depth
		if vDepth > 0 {
			for rec := r.activeLocalRec; rec != nil && vDepth <= rec.parentDepth; rec = rec.parentRec {
				rec.ineligible = true
			}
		}
	}
}

func (r *RemoveUnusedCode) onPolyfillDefined(keyID strpool.ID) {
	if r.activeGlobalRec != nil {
		r.activeGlobalRec.ineligible = true
	}
	if !r.polyBits.Get(int(keyID)) {
		r.polyEpoch++
		if r.fnCache != nil && len(r.fnCache.byFn) > 0 && !r.fnCache.recordedPolyBits.Get(int(keyID)) {
			r.summaryReplayDisabled = true
		}
	}
}

func (r *RemoveUnusedCode) isPropPinnedOrExported(id uint32) bool {
	if !r.trackProps || r.pinnedProps.Get(int(id)) {
		return true
	}
	idx := r.ensurePropEntry(strpool.ID(id))
	if !r.propEntries[idx].checkedUnpinned {
		if r.externProperties.HasID(strpool.ID(id)) ||
			r.isExportedByConvention(strpool.String(strpool.ID(id)), false) {
			r.pinnedProps.Set(int(id))
			return true
		}
		r.propEntries[idx].checkedUnpinned = true
	}
	return false
}

func (r *RemoveUnusedCode) allIndepPropsPinned(fc *rucFnSummaryCache, s rucFnSummary) bool {
	if s.indepLen == 0 {
		return true
	}
	for _, id := range fc.indepPropArena[s.indepOff : s.indepOff+s.indepLen] {
		if !r.isPropPinnedOrExported(id) {
			return false
		}
	}
	return true
}

func (r *RemoveUnusedCode) resolveReplayVars(parentScope *Scope, vars []uint32) bool {
	r.replayVars = r.replayVars[:0]
	for _, id := range vars {
		v := parentScope.GetVarByID(id)
		if v == nil {
			clear(r.replayVars)
			r.replayVars = r.replayVars[:0]
			return false
		}
		r.replayVars = append(r.replayVars, v)
	}
	return true
}

func (r *RemoveUnusedCode) applyGlobalSummary(fc *rucFnSummaryCache, fn rhino.Node, s rucFnSummary, parentScope *Scope) bool {
	vars := fc.varIDArena[s.varOff : s.varOff+s.varLen]
	if !r.resolveReplayVars(parentScope, vars) {
		return false
	}
	prevGlobal := r.activeGlobalRec
	r.activeGlobalRec = nil
	for _, v := range r.replayVars {
		r.TraverseVar(v).pin(rucPinReferenceFound)
	}
	clear(r.replayVars)
	r.replayVars = r.replayVars[:0]
	for _, propID := range fc.propIDArena[s.propOff : s.propOff+s.propLen] {
		r.markPropertyIDPinned(propID)
	}
	for _, n := range fc.polyNodeArena[s.polyOff : s.polyOff+s.polyLen] {
		keyID := n.GetStringID()
		if r.polyBits.Get(int(keyID)) {
			r.considerPolyfillRef(keyID, n)
		}
	}
	r.activeGlobalRec = prevGlobal
	r.markFnReplayed(fn)
	return true
}

func (r *RemoveUnusedCode) tryReplayDeferredSummary(d rucDeferredFn, globalScope *Scope) bool {
	if r.summaryReplayDisabled || r.fnCache == nil || d.fn.IsDeleted() {
		return false
	}
	fc := r.fnCache
	s := d.sum
	if !r.allIndepPropsPinned(fc, s) {
		return false
	}
	return r.applyGlobalSummary(fc, d.fn, s, globalScope)
}

func (r *RemoveUnusedCode) tryReplayOrDeferFunction(fn rhino.Node, parentScope *Scope, allowDefer bool) bool {
	if r.summaryReplayDisabled || r.fnCache == nil || fn.IsDeleted() || r.isDirtySubtree(fn) {
		return false
	}
	fc := r.fnCache
	s, ok := fc.byFn[fn]
	if !ok || !s.isValid() {
		return false
	}
	isGlobalParent := parentScope.IsGlobal()
	if s.isGlobalParent() != isGlobalParent || !s.matchesContext(fn) {
		return false
	}

	if r.recMode == rucRecGlobal {
		if !isGlobalParent {
			return false
		}
		if s.polyLen > 0 {
			for _, n := range fc.polyNodeArena[s.polyOff : s.polyOff+s.polyLen] {
				if n.IsDeleted() {
					return false
				}
				keyID := n.GetStringID()
				if fc.recordedPolyBits.Get(int(keyID)) && !r.polyBits.Get(int(keyID)) {
					// Polyfill definition has not been visited yet on this run.
					return false
				}
			}
		}
		if !r.allIndepPropsPinned(fc, s) {
			if allowDefer {
				r.deferredSummaryFns = append(r.deferredSummaryFns, rucDeferredFn{fn: fn, sum: s})
				r.deferredCount++
				return true
			}
			return false
		}
		return r.applyGlobalSummary(fc, fn, s, parentScope)
	}

	if r.recMode == rucRecLocal {
		if !isGlobalParent && s.varLen > 0 {
			vars := fc.varIDArena[s.varOff : s.varOff+s.varLen]
			if !r.resolveReplayVars(parentScope, vars) {
				return false
			}
			// Keep r.activeLocalRec active so pins on outer scopes propagate to
			// any enclosing function currently being recorded.
			for _, v := range r.replayVars {
				r.TraverseVar(v).pin(rucPinReferenceFound)
			}
			clear(r.replayVars)
			r.replayVars = r.replayVars[:0]
		}
		r.markFnReplayed(fn)
		return true
	}

	return false
}

func (r *RemoveUnusedCode) markFnReplayed(fn rhino.Node) {
	r.replayedCount++
	if rucIncrVerify {
		if r.replayedFns == nil {
			r.replayedFns = make(map[rhino.Node]struct{}, 64)
		}
		r.replayedFns[fn] = struct{}{}
	}
}

func (r *RemoveUnusedCode) drainWorklistAndDeferred(globalScope *Scope) {
	deferHead := 0
	lastCheckedPinEpoch := ^uint32(0)
	for {
		for !r.worklist.isEmpty() {
			continuation := r.worklist.pop()
			continuation.Apply(r)
		}
		if deferHead >= len(r.deferredSummaryFns) {
			r.deferredSummaryFns = r.deferredSummaryFns[:0]
			return
		}
		if r.summaryReplayDisabled {
			fns := r.deferredSummaryFns[deferHead:]
			r.deferredSummaryFns = nil
			deferHead = 0
			for _, d := range fns {
				r.traverseFunctionFull(d.fn, globalScope)
			}
			continue
		}
		if r.propPinEpoch != lastCheckedPinEpoch {
			lastCheckedPinEpoch = r.propPinEpoch
			progressed := false
			writeIdx := deferHead
			for i := deferHead; i < len(r.deferredSummaryFns); i++ {
				d := r.deferredSummaryFns[i]
				if r.tryReplayDeferredSummary(d, globalScope) {
					progressed = true
				} else {
					r.deferredSummaryFns[writeIdx] = d
					writeIdx++
				}
			}
			r.deferredSummaryFns = r.deferredSummaryFns[:writeIdx]
			if progressed || deferHead >= len(r.deferredSummaryFns) {
				continue
			}
		}
		d := r.deferredSummaryFns[deferHead]
		deferHead++
		r.traverseFunctionFull(d.fn, globalScope)
	}
}

func (r *RemoveUnusedCode) canSkipGlobalLocalFunctionUnit(u rhino.Node) bool {
	if !u.IsFunction() || r.summaryReplayDisabled || r.fnCache == nil || u.IsDeleted() || r.isDirtySubtree(u) {
		return false
	}
	s, ok := r.fnCache.byFn[u]
	return ok && s.isValid() && s.isGlobalParent() && s.matchesContext(u)
}

func (r *RemoveUnusedCode) traverseAndRemoveLocalScriptsParallel(root rhino.Node) bool {
	if !r.canParallelizeLocalScripts() || root.GetParent().IsNil() {
		return false
	}
	if r.compiler == nil || r.compiler.GetOptions() == nil || r.compiler.GetOptions().GetNumParallelThreads() <= 1 {
		return false
	}
	scripts, ok := collectRootScripts(root)
	if !ok || len(scripts) == 0 {
		return false
	}

	var units []rhino.Node
	for _, script := range scripts {
		mark := len(units)
		if (!script.GetFirstChild().IsNil() && script.GetFirstChild().IsModuleBody()) ||
			!collectGlobalLocalScopeUnits(script, &units) {
			units = append(units[:mark], script)
		}
	}

	workUnits := make([]rhino.Node, 0, len(units))
	for _, u := range units {
		if r.canSkipGlobalLocalFunctionUnit(u) {
			r.markFnReplayed(u)
			continue
		}
		workUnits = append(workUnits, u)
	}

	if len(workUnits) == 0 {
		if rucIncrVerify && r.replayedCount > 0 {
			gs := r.createGlobalScope(root)
			r.verifySummaryReplay(root, gs)
		}
		r.finishSummaryRun(nil)
		return true
	}

	numWorkers := parallelThreadCount(r.compiler, len(workUnits))
	if numWorkers <= 1 {
		r.cacheVarInfos = r.compiler.parallelRecorder == nil
		if r.varByKey == nil {
			r.varByKey = make(map[scopedVarKey]*RemoveUnusedCode_RealVarInfo, 256)
		}
		gs := r.createGlobalScope(root)
		for _, u := range workUnits {
			r.traverseLocalUnit(u, gs)
		}
		r.finishSummaryRun(nil)
		return true
	}

	globalScope := r.createGlobalScope(root)
	freezeGlobalScopeForParallelReads(globalScope)
	prewarmCompilerForParallelReads(r.compiler)
	restore := enableParallelTypeRegistryReads(r.compiler)
	defer restore()

	workerRecs := make([]*rucRecorderState, numWorkers)
	workerReplayed := make([]int, numWorkers)
	runParallelWork(r.compiler, numWorkers, len(workUnits), func(wIdx int) (func(int, *recordingCompiler), func()) {
		var wr *RemoveUnusedCode
		return func(idx int, rc *recordingCompiler) {
				if wr == nil {
					wr = r.newLocalScriptWorker(rc)
					wr.root = root
					wr.fnCache = r.fnCache
					wr.dirtySubtreeRoots = r.dirtySubtreeRoots
					wr.summaryReplayDisabled = r.summaryReplayDisabled
					if r.recMode == rucRecLocal {
						wr.recMode = rucRecLocal
						wr.recState = newRucRecorderState()
						workerRecs[wIdx] = wr.recState
					}
				} else {
					wr.resetLocalScriptWorker()
				}
				wr.traverseRootAndRemoveVars(workUnits[idx], globalScope)
			}, func() {
				if wr != nil {
					workerReplayed[wIdx] = wr.replayedCount
				}
			}
	}, false)

	for _, n := range workerReplayed {
		r.replayedCount += n
	}
	r.finishSummaryRun(workerRecs)
	return true
}

func (r *RemoveUnusedCode) traverseLocalUnit(unit rhino.Node, globalScope *Scope) {
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
	r.traverseRootAndRemoveVars(unit, globalScope)
}

func (r *RemoveUnusedCode) wasEnclosedByReplayedFn(n rhino.Node) bool {
	if len(r.replayedFns) == 0 {
		return false
	}
	for cur := n; !cur.IsNil(); cur = cur.GetParent() {
		if _, ok := r.replayedFns[cur]; ok {
			return true
		}
	}
	return false
}

func (r *RemoveUnusedCode) verifySummaryReplay(node rhino.Node, scope *Scope) {
	if !rucIncrVerify || r.replayedCount == 0 {
		return
	}
	opts := r.optionsSnapshot()
	opts.RemoveUnusedPolyfills = r.removeUnusedPolyfills
	ref := NewRemoveUnusedCode(r.compiler, opts)
	ref.externProperties = r.externProperties
	ref.accessorSummary = r.accessorSummary
	ref.cacheVarInfos = false
	ref.summaryReplayDisabled = true
	ref.recMode = rucRecNone
	ref.root = r.root
	ref.trackProps = ref.removesAnyProperties()
	if ref.trackProps {
		for _, id := range rucImplicitPropIDs {
			ref.pinnedProps.Set(int(id))
		}
	}
	ref.varByKey = make(map[scopedVarKey]*RemoveUnusedCode_RealVarInfo, 256)

	refScope := scope
	if r.recMode == rucRecGlobal {
		refScope = ref.createGlobalScope(r.root)
	}
	ref.worklist.push(RemoveUnusedCode_Continuation{node: node, scope: refScope})
	for !ref.worklist.isEmpty() {
		cont := ref.worklist.pop()
		cont.Apply(ref)
	}

	if r.recMode == rucRecGlobal {
		for _, vi := range ref.varInfos {
			if !vi.var_.IsGlobal() {
				continue
			}
			gotVar := scope.GetVarByID(vi.var_.nameID)
			if gotVar == nil {
				panic(fmt.Sprintf("RemoveUnusedCode verify: missing global var %s", vi.var_.GetName()))
			}
			gotRem := r.GetVarInfo(gotVar).IsRemovable()
			if gotRem != vi.IsRemovable() {
				panic(fmt.Sprintf("RemoveUnusedCode verify: global var %s removable=%v, full=%v", vi.var_.GetName(), gotRem, vi.IsRemovable()))
			}
		}
		if r.trackProps {
			for _, idx := range ref.propRemOrder {
				entry := &ref.propEntries[idx]
				if !entry.removables.isEmpty() && r.isPropPinnedOrExported(uint32(entry.propID)) {
					panic(fmt.Sprintf("RemoveUnusedCode verify: prop %s pinned in replay but unreferenced in full", strpool.String(entry.propID)))
				}
			}
			for _, idx := range r.propRemOrder {
				entry := &r.propEntries[idx]
				if !entry.removables.isEmpty() && ref.isPropPinnedOrExported(uint32(entry.propID)) {
					panic(fmt.Sprintf("RemoveUnusedCode verify: prop %s unreferenced in replay but pinned in full", strpool.String(entry.propID)))
				}
			}
		}
		if len(ref.polyfillOrder) != len(r.polyfillOrder) {
			panic(fmt.Sprintf("RemoveUnusedCode verify: polyfill count replay=%d full=%d", len(r.polyfillOrder), len(ref.polyfillOrder)))
		}
		for i := range ref.polyfillOrder {
			if ref.polyfillOrder[i].isRemovable != r.polyfillOrder[i].isRemovable {
				panic(fmt.Sprintf("RemoveUnusedCode verify: polyfill %s removable=%v full=%v",
					ref.polyfillOrder[i].getName(), r.polyfillOrder[i].isRemovable, ref.polyfillOrder[i].isRemovable))
			}
		}
		if ref.removeGlobals {
			for _, fparamScope := range ref.allFunctionParamScopes {
				if _, replayed := r.replayedFns[fparamScope.GetRootNode()]; replayed {
					ref.assertParamScopeNoPendingChanges(fparamScope, "replayed global fn")
				}
			}
		}
		return
	}

	if len(ref.polyfillOrder) != 0 {
		panic(fmt.Sprintf("RemoveUnusedCode verify: rucRecLocal found %d polyfill defs", len(ref.polyfillOrder)))
	}
	for _, vi := range ref.varInfos {
		if r.wasEnclosedByReplayedFn(vi.var_.scope.rootNode) {
			if !vi.IsRemovable() {
				continue
			}
			for cur := vi.removables.head; cur != nil; cur = cur.nextInQueue {
				if cur.remKind == remKindIndirectAssign &&
					cur.indirectAssignRoot().GetParent().IsParamList() &&
					!cur.indirectAssignRoot().IsDefaultValue() {
					continue
				}
				panic(fmt.Sprintf("RemoveUnusedCode verify: skipped fn contains removable var %s", vi.var_.GetName()))
			}
			toRemove := vi.var_.GetNameNode().GetParent()
			if fn, isFn := ast.AsFunction(toRemove); isFn && NodeUtil_IsFunctionExpression(fn.Node()) && !r.preserveFunctionExpressionNames {
				panic(fmt.Sprintf("RemoveUnusedCode verify: skipped fn has removable bleeding name %s", vi.var_.GetName()))
			}
			continue
		}
		rVi := r.varByKey[scopedVarKeyOf(vi.var_)]
		if rVi == nil || rVi.IsRemovable() != vi.IsRemovable() {
			panic(fmt.Sprintf("RemoveUnusedCode verify: local var %s removable mismatch (got=%v, full=%v)",
				vi.var_.GetName(), rVi != nil && rVi.IsRemovable(), vi.IsRemovable()))
		}
	}
}
