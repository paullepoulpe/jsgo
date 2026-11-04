package jscomp

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"sync/atomic"
	"time"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

// Go-only: OptimizeCalls' reference map read from the property-use index.
//
// OptimizeCalls.Process used to build its ReferenceMap by walking the whole
// program on every run (optimizeCallsRefCollector.traverseNode), although
// only a few percent of the program changes between two runs of the main
// loop. The index's records (prop_use_index.go) also list, per change scope
// and in the order the walk meets them, what the walk files (piListOC):
//
//   - every NAME, with its id (piOCName). Which of them are global, non-extern
//     names is decided when the list is read, from incremental global refcounts
//     (or the run's global scope on synthetic inputs without CompilerInputs);
//   - every property reference: GETPROP, OPTCHAIN_GETPROP, unquoted
//     STRING_KEY and class or object members (piOCProp). Extern properties and
//     getter/setter accessor names are filtered via bitset before filing;
//   - two kinds of nodes to replay, because what they file depends on more
//     than their change scope or can fail: a SUPER (visitSuper files a
//     super() call under the enclosing class's extends clause, which lies in
//     another change scope) and a call of a property rename function (its
//     first argument's string).
//
// A read files the items into dense per-id CSR arrays on reusable scratch
// buffers, preserving exact first-appearance ID order and per-ID node order
// without map inserts or per-symbol slab allocations.
//
// JSGO_PROPINDEX_OC=0 makes OptimizeCalls walk again. With
// JSGO_PROPINDEX_VERIFY=1 every read also walks and panics if the maps
// differ.

var (
	propIndexOC = os.Getenv("JSGO_PROPINDEX_OC") != "0"
	// propIndexOCReads counts OptimizeCalls' reads (tests).
	propIndexOCReads atomic.Int64
)

// propUseIndexForOptimizeCalls returns the index when OptimizeCalls may read
// its reference map from it, or nil when it must walk.
func propUseIndexForOptimizeCalls(o *OptimizeCalls, externs, root rhino.Node) *propUseIndex {
	if !propIndexOC || o.considerExterns {
		return nil
	}
	ix := propUseIndexFor(o.compiler, externs, root)
	if ix != nil {
		propIndexOCReads.Add(1)
	}
	return ix
}

// ocName records a NAME, as maybeAddNameReferenceNode filters it before
// resolving it, and flags non-self `arguments` references on the scope.
func (r *piRecorder) ocName(n rhino.Node) {
	if id := n.GetStringID(); id > strpool.EmptyID {
		if id == implicitArgumentsID && !(r.cur.root.IsFunction() && n == r.cur.root.GetFirstChild()) {
			r.cur.flags |= piScopeFlagHasArgumentsName
		}
		r.lists[piListOC] = append(r.lists[piListOC], piItem{n: n, prop: id, bits: piOCName})
	}
}

// The verdicts of a read on a name id the global scope may declare (fallback
// mode when synthetic scripts lack CompilerInputs).
const (
	ocVerdictUnknown uint8 = iota // not looked up yet
	ocVerdictKeep                 // a global, non-extern name
	ocVerdictDrop                 // not one
)

// Clean-symbol cache masks per OptimizeCalls sub-pass.
const (
	ocCleanReturns uint8 = 1 << iota
	ocCleanParams
)

type ocCleanEntry struct {
	off   int32
	len   int32
	flags uint8
}

type ocCleanTable struct {
	entries map[uint32]ocCleanEntry
	nodes   []rhino.Node
	dead    int
}

func (c *ocCleanTable) reset() {
	clear(c.entries)
	c.nodes = c.nodes[:0]
	c.dead = 0
}

func (c *ocCleanTable) invalidate(id uint32) {
	if e, ok := c.entries[id]; ok {
		c.dead += int(e.len)
		delete(c.entries, id)
	}
}

func (c *ocCleanTable) isClean(id uint32, mask uint8, refs []rhino.Node) bool {
	e, ok := c.entries[id]
	if !ok || e.flags&mask == 0 || int(e.len) != len(refs) {
		return false
	}
	return slices.Equal(refs, c.nodes[e.off:e.off+e.len])
}

func (c *ocCleanTable) markClean(id uint32, mask uint8, refs []rhino.Node) {
	if c.entries == nil {
		c.entries = make(map[uint32]ocCleanEntry, 256)
	}
	if e, ok := c.entries[id]; ok {
		if int(e.len) == len(refs) && slices.Equal(refs, c.nodes[e.off:e.off+e.len]) {
			e.flags |= mask
			c.entries[id] = e
			return
		}
		if len(refs) <= int(e.len) && e.len > 0 {
			copy(c.nodes[e.off:int(e.off)+len(refs)], refs)
			c.dead += int(e.len) - len(refs)
			e.len = int32(len(refs))
			e.flags = mask
			c.entries[id] = e
			return
		}
		c.dead += int(e.len)
		delete(c.entries, id)
	}
	if c.dead > 4096 && c.dead > len(c.nodes)/2 {
		c.compact()
	}
	off := int32(len(c.nodes))
	c.nodes = append(c.nodes, refs...)
	c.entries[id] = ocCleanEntry{
		off:   off,
		len:   int32(len(refs)),
		flags: mask,
	}
}

func (c *ocCleanTable) compact() {
	dst := make([]rhino.Node, 0, max(len(c.nodes)-c.dead+64, 64))
	// order-independent: each live entry gets a private contiguous slice of its identical refs
	for id, e := range c.entries {
		if e.flags == 0 {
			delete(c.entries, id)
			continue
		}
		newOff := int32(len(dst))
		dst = append(dst, c.nodes[e.off:e.off+e.len]...)
		e.off = newOff
		c.entries[id] = e
	}
	c.nodes = dst
	c.dead = 0
}

// ocFillTable is one table of the map (names or properties) as a read fills
// it: per-id lists in dense CSR arrays reused across runs.
type ocFillTable struct {
	slot      []int32      // id -> 1 + its index in ids, 0 when not met in this read
	ids       []strpool.ID // the ids in first-appearance order
	count     []int32      // per index: its nodes
	slotDirty []bool       // per index: whether any ref came from a dirty scope
	refs      []ocRef      // the nodes filed, in traversal order

	flatIDs     []strpool.ID
	flatOffsets []int32
	flatNodes   []rhino.Node
	flatDirty   []bool
}

type ocRef struct {
	slot int32
	n    rhino.Node
}

func (t *ocFillTable) add(id strpool.ID, n rhino.Node, dirty bool) {
	if int(id) >= len(t.slot) {
		t.slot = slices.Grow(t.slot, int(id)+1-len(t.slot)+len(t.slot)/2)
		t.slot = t.slot[:cap(t.slot)]
	}
	s := t.slot[id]
	if s == 0 {
		t.ids = append(t.ids, id)
		t.count = append(t.count, 0)
		t.slotDirty = append(t.slotDirty, dirty)
		s = int32(len(t.ids))
		t.slot[id] = s
	} else if dirty {
		t.slotDirty[s-1] = true
	}
	t.count[s-1]++
	t.refs = append(t.refs, ocRef{slot: s - 1, n: n})
}

// flush scatters the collected references into reusable CSR arrays on t,
// invalidates dirty slots in clean, and resets the scratch accumulators.
func (t *ocFillTable) flush(data *optimizeCalls_refTable, clean *ocCleanTable) {
	nIDs := len(t.ids)
	nRefs := len(t.refs)

	t.flatIDs = append(t.flatIDs[:0], t.ids...)
	t.flatDirty = append(t.flatDirty[:0], t.slotDirty...)
	if cap(t.flatOffsets) < nIDs+1 {
		t.flatOffsets = make([]int32, nIDs+1)
	} else {
		t.flatOffsets = t.flatOffsets[:nIDs+1]
	}
	if cap(t.flatNodes) < nRefs {
		t.flatNodes = make([]rhino.Node, nRefs)
	} else {
		t.flatNodes = t.flatNodes[:nRefs]
	}

	t.flatOffsets[0] = 0
	end := int32(0)
	for i, c := range t.count {
		t.count[i] = end
		end += c
		t.flatOffsets[i+1] = end
	}
	for _, r := range t.refs {
		t.flatNodes[t.count[r.slot]] = r.n
		t.count[r.slot]++
	}
	for i, id := range t.ids {
		if t.slotDirty[i] && clean != nil {
			clean.invalidate(uint32(id))
		}
		t.slot[id] = 0
	}
	t.ids = t.ids[:0]
	t.count = t.count[:0]
	t.slotDirty = t.slotDirty[:0]
	t.refs = t.refs[:0]

	data.byID = nil
	data.order = nil
	data.flatIDs = t.flatIDs
	data.flatOffsets = t.flatOffsets
	data.flatNodes = t.flatNodes
	data.flatDirty = t.flatDirty
	data.size = nIDs
}

// ocFillState is what a read of the index needs besides the records, kept
// across reads and dropped with propUseIndex in dropPropUseIndex.
type ocFillState struct {
	names, props ocFillTable
	skipPropBits smallBitSet

	ix           *propUseIndex
	useRefcounts bool
	maybe        []uint64     // fallback global scope bitset
	verdict      []uint8      // fallback name id -> ocVerdict*
	marked       []strpool.ID // ids set in maybe or verdict
	global       *Scope
	lookupAll    bool

	namesClean           ocCleanTable
	propsClean           ocCleanTable
	cleanReady           bool
	lastGlobalAddedEpoch uint32
	lastSideEffectCfg    inlineSideEffectConfig
	lastConv             CodingConvention
	liveMutatedScopes    map[rhino.Node]struct{}

	// Statistics (JSGO_PROPINDEX_STATS=1).
	reads int
	fill  time.Duration
}

func (st *ocFillState) mark(id strpool.ID) {
	w := int(id >> 6)
	if w >= len(st.maybe) {
		st.maybe = slices.Grow(st.maybe, w+1-len(st.maybe)+len(st.maybe)/2)
		st.maybe = st.maybe[:cap(st.maybe)]
	}
	if int(id) >= len(st.verdict) {
		st.verdict = slices.Grow(st.verdict, int(id)+1-len(st.verdict)+len(st.verdict)/2)
		st.verdict = st.verdict[:cap(st.verdict)]
	}
	if bit := uint64(1) << (id & 63); st.maybe[w]&bit == 0 {
		st.maybe[w] |= bit
		st.marked = append(st.marked, id)
	}
}

func (st *ocFillState) begin(ix *propUseIndex, gs *Scope) {
	st.ix = ix
	if gs == nil {
		st.useRefcounts = true
		return
	}
	st.useRefcounts = false
	st.global = gs
	st.lookupAll = gs.lookupParent != nil
	for _, v := range gs.varList {
		st.mark(strpool.ID(v.nameID))
	}
	for _, id := range [...]uint32{implicitArgumentsID, implicitThisID, implicitExportsID, implicitSuperID} {
		st.mark(strpool.ID(id))
	}
}

func (st *ocFillState) end() {
	if !st.useRefcounts {
		for _, id := range st.marked {
			st.maybe[id>>6] = 0
			st.verdict[id] = ocVerdictUnknown
		}
		st.marked = st.marked[:0]
		st.global = nil
	}
	st.skipPropBits.ClearAll()
}

// isGlobalName is maybeAddNameReferenceNode's test with considerExterns
// false: the run's global scope declares id, and not in the externs.
func (st *ocFillState) isGlobalName(id strpool.ID) bool {
	if st.useRefcounts {
		return st.ix.isGlobalNonExtern(uint32(id))
	}
	if w := int(id >> 6); w >= len(st.maybe) || st.maybe[w]&(uint64(1)<<(id&63)) == 0 {
		if !st.lookupAll {
			return false
		}
		st.mark(id)
	}
	switch st.verdict[id] {
	case ocVerdictKeep:
		return true
	case ocVerdictDrop:
		return false
	}
	keep := false
	if v := st.global.GetVarByID(uint32(id)); v != nil && !v.IsExtern() {
		keep = true
	}
	if keep {
		st.verdict[id] = ocVerdictKeep
	} else {
		st.verdict[id] = ocVerdictDrop
	}
	return keep
}

func (ix *propUseIndex) prepareOCDirtyStamp() uint32 {
	st := &ix.oc
	changed := ix.c.changeTracker.GetChangedScopeNodesForPass(ocTimelineKey)
	curCfg := currentInlineSideEffectConfig(ix.c)
	curConv := ix.c.GetCodingConvention()
	if changed == nil || !st.cleanReady ||
		st.lastGlobalAddedEpoch != ix.globalAddedEpoch ||
		st.lastSideEffectCfg != curCfg ||
		st.lastConv != curConv {
		st.namesClean.reset()
		st.propsClean.reset()
		st.cleanReady = true
		st.lastGlobalAddedEpoch = ix.globalAddedEpoch
		st.lastSideEffectCfg = curCfg
		st.lastConv = curConv
	}
	if len(changed) == 0 {
		return 0
	}
	ix.dirtyGen++
	if ix.dirtyGen == 0 {
		ix.dirtyGen = 1
	}
	gen := ix.dirtyGen
	for _, n := range changed {
		for cur := n; !cur.IsNil(); {
			if s := ix.scopes[cur]; s != nil {
				if s.dirtyStamp == gen {
					break
				}
				s.dirtyStamp = gen
			}
			p := cur.GetParent()
			if p.IsNil() {
				break
			}
			cur = ChangeTracker_GetEnclosingChangeScopeRoot(p)
		}
	}
	return gen
}

func (st *ocFillState) beginPassLoop(c *Compiler) {
	if c == nil || c.changeTracker == nil {
		return
	}
	clear(st.liveMutatedScopes)
	_ = c.changeTracker.GetChangedScopeNodesForPass(ocIntraTimelineKey)
}

func (st *ocFillState) syncLiveMutatedScopes(c *Compiler) {
	if c == nil || c.changeTracker == nil {
		return
	}
	changed := c.changeTracker.GetChangedScopeNodesForPass(ocIntraTimelineKey)
	if len(changed) == 0 {
		return
	}
	if st.liveMutatedScopes == nil {
		st.liveMutatedScopes = make(map[rhino.Node]struct{}, max(len(changed)*2, 16))
	}
	for _, n := range changed {
		for cur := n; !cur.IsNil(); {
			if _, seen := st.liveMutatedScopes[cur]; seen {
				break
			}
			st.liveMutatedScopes[cur] = struct{}{}
			p := cur.GetParent()
			if p.IsNil() {
				break
			}
			cur = ChangeTracker_GetEnclosingChangeScopeRoot(p)
		}
	}
}

func (st *ocFillState) endPassLoop(c *Compiler) {
	clear(st.liveMutatedScopes)
	if c != nil && c.changeTracker != nil {
		c.changeTracker.changeTimeline.Unmark(ocIntraTimelineKey)
	}
}

func (st *ocFillState) refsTouchLiveMutatedScopes(refs []rhino.Node) bool {
	if len(st.liveMutatedScopes) == 0 {
		return false
	}
	for _, n := range refs {
		sc := ChangeTracker_GetEnclosingChangeScopeRoot(n)
		if !sc.IsNil() {
			if _, hit := st.liveMutatedScopes[sc]; hit {
				return true
			}
		}
	}
	return false
}

func (st *ocFillState) canSkipClean(isProp bool, slotDirty bool, id uint32, refs []rhino.Node, mask uint8) bool {
	if slotDirty || !st.cleanReady {
		return false
	}
	tbl := &st.namesClean
	if isProp {
		tbl = &st.propsClean
	}
	if !tbl.isClean(id, mask, refs) {
		return false
	}
	return !st.refsTouchLiveMutatedScopes(refs)
}

func (st *ocFillState) markClean(isProp bool, id uint32, refs []rhino.Node, mask uint8) {
	if !st.cleanReady {
		return
	}
	tbl := &st.namesClean
	if isProp {
		tbl = &st.propsClean
	}
	tbl.markClean(id, mask, refs)
}

// fillOptimizeCalls populates m from the property-use index with
// considerExterns false and pre-flush accessor filtering.
func (ix *propUseIndex) fillOptimizeCalls(o *OptimizeCalls, m *OptimizeCalls_ReferenceMap, fallbackGlobal *Scope) {
	st := &ix.oc
	var t0 time.Time
	if propIndexStatsOn {
		t0 = time.Now()
	}

	dirtyGen := ix.prepareOCDirtyStamp()
	st.skipPropBits.ClearAll()
	if externProps := o.compiler.GetExternPropSet(); externProps != nil {
		for _, id := range externProps.IDs() {
			st.skipPropBits.Set(int(id))
		}
	}
	if summary := o.compiler.GetAccessorSummary(); summary != nil {
		for _, id := range summary.AccessorIDs() {
			st.skipPropBits.Set(int(id))
		}
	}

	st.begin(ix, fallbackGlobal)
	f := ocFiller{ix: ix, st: st, dirtyGen: dirtyGen}
	for _, top := range [2]rhino.Node{ix.externs, ix.root} {
		for script := top.GetFirstChild(); !script.IsNil(); script = script.GetNext() {
			if script.IsFromExterns() {
				continue // as traverseNode
			}
			f.walk(ix.scopes[script], false)
		}
	}
	st.end()
	if propIndexStatsOn {
		st.reads++
	}
	st.names.flush(&m.names, &st.namesClean)
	st.props.flush(&m.props, &st.propsClean)
	if propIndexStatsOn {
		d := time.Since(t0)
		st.fill += d
		fmt.Fprintf(os.Stderr, "[propindex] optimizeCalls read %d: %d items, %d names, %d properties, %v\n",
			st.reads, f.items, m.names.Size(), m.props.Size(), d.Round(10*time.Microsecond))
	}
}

type ocFiller struct {
	ix       *propUseIndex
	st       *ocFillState
	dirtyGen uint32
	items    int
}

// walk files the items of s and of its nested functions' records, spliced in
// where the traversal enters them.
func (f *ocFiller) walk(s *piScope, parentDirty bool) {
	st := f.st
	scopeDirty := f.dirtyGen != 0 && s.dirtyStamp == f.dirtyGen
	fnSelfName := rhino.NilNode
	if parentDirty && !scopeDirty && s.root.IsFunction() {
		fnSelfName = s.root.GetFirstChild()
	}
	kids := s.kids
	k := 0
	items := s.items[piListOC]
	f.items += len(items)
	for i := range items {
		for k < len(kids) && int(kids[k].at[piListOC]) == i {
			f.walk(kids[k].s, scopeDirty)
			k++
		}
		it := &items[i]
		switch it.bits {
		case piOCName:
			if id := strpool.ID(it.prop); st.isGlobalName(id) {
				st.names.add(id, it.n, scopeDirty || it.n == fnSelfName)
			}
		case piOCProp:
			if !st.skipPropBits.Get(int(it.prop)) {
				st.props.add(strpool.ID(it.prop), it.n, scopeDirty)
			}
		case piOCSuper:
			// visitSuper
			target, isName := optimizeCallsSuperTarget(it.n)
			if target.IsNil() {
				break
			}
			superDirty := scopeDirty
			if !superDirty && f.dirtyGen != 0 {
				if cls := NodeUtil_GetEnclosingClass(it.n); !cls.IsNil() {
					if cs := f.ix.scopes[ChangeTracker_GetEnclosingChangeScopeRoot(cls)]; cs != nil && cs.dirtyStamp == f.dirtyGen {
						superDirty = true
					}
				}
			}
			if isName {
				// maybeAddNameReferenceNode(target, superNode)
				if id := nodeStringID(target); id > strpool.EmptyID && st.isGlobalName(id) {
					st.names.add(id, it.n, superDirty)
				}
			} else if id := nodeStringID(target); !st.skipPropBits.Get(int(id)) {
				st.props.add(id, it.n, superDirty)
			}
		case piOCRename:
			// traverseNode's CALL case
			firstArg, _ := ast.MustCall(it.n).FirstArg()
			if id := nodeStringID(firstArg.Node()); !st.skipPropBits.Get(int(id)) {
				st.props.add(id, it.n, scopeDirty)
			}
		}
	}
	for ; k < len(kids); k++ {
		f.walk(kids[k].s, scopeDirty)
	}
}

// verifyOptimizeCallsMap panics if got, read from the index, differs from
// want, built by the walk: the ids of each table in first-appearance order,
// and each id's nodes.
func verifyOptimizeCallsMap(c *Compiler, got, want *OptimizeCalls_ReferenceMap) {
	cmp := func(kind string, g, w *optimizeCalls_refTable) {
		var gi, wi []strpool.ID
		var gr, wr [][]rhino.Node
		for id, refs := range g.allByID() {
			gi = append(gi, strpool.ID(id))
			gr = append(gr, refs)
		}
		for id, refs := range w.allByID() {
			wi = append(wi, strpool.ID(id))
			wr = append(wr, refs)
		}
		fail := func(format string, args ...any) {
			var b strings.Builder
			fmt.Fprintf(&b, "optimizeCalls' reference map read from the property use index differs from the walk's in %s: %s ", c.currentPassName(), kind)
			fmt.Fprintf(&b, format, args...)
			panic(b.String())
		}
		for i := 0; i < len(gi) || i < len(wi); i++ {
			if i >= len(gi) || i >= len(wi) || gi[i] != wi[i] {
				name := func(ids []strpool.ID) string {
					if i < len(ids) {
						return strpool.String(ids[i])
					}
					return "(none)"
				}
				fail("id %d of %d read, %d walked: %q read, %q walked", i, len(gi), len(wi), name(gi), name(wi))
			}
		}
		for i, id := range gi {
			gn, wn := gr[i], wr[i]
			if !slices.Equal(gn, wn) {
				j := 0
				for j < len(gn) && j < len(wn) && gn[j] == wn[j] {
					j++
				}
				at := func(ns []rhino.Node) string {
					if j < len(ns) {
						return piPosition(ns[j])
					}
					return "(none)"
				}
				fail("%q: node %d of %d read, %d walked: %s read, %s walked", strpool.String(id), j, len(gn), len(wn), at(gn), at(wn))
			}
		}
	}
	cmp("names", &got.names, &want.names)
	cmp("properties", &got.props, &want.props)
}

func (ix *propUseIndex) verifyGlobalScopeParity(gs *Scope) {
	maxID := max(len(ix.globalNonExternBits)*64, len(ix.globalAnyBits)*64)
	// order-independent: computes maximum ID across map keys for exhaustive parity check
	for id := range ix.globalCounts {
		if int(id)+1 > maxID {
			maxID = int(id) + 1
		}
	}
	for _, v := range gs.varList {
		if int(v.nameID)+1 > maxID {
			maxID = int(v.nameID) + 1
		}
	}
	for i := 1; i < maxID; i++ {
		id := uint32(i)
		v := gs.GetVarByID(id)
		wantHas := v != nil
		gotHas := ix.hasGlobalVarByID(id)
		if gotHas != wantHas {
			panic(fmt.Sprintf("propUseIndex.hasGlobalVarByID(%q) = %v, want %v",
				strpool.String(strpool.ID(id)), gotHas, wantHas))
		}
		wantNonExt := v != nil && !v.IsExtern()
		gotNonExt := ix.isGlobalNonExtern(id)
		if gotNonExt != wantNonExt {
			panic(fmt.Sprintf("propUseIndex.isGlobalNonExtern(%q) = %v, want %v",
				strpool.String(strpool.ID(id)), gotNonExt, wantNonExt))
		}
	}
}
