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
//     names is decided when the list is read, from the run's global scope:
//     nothing resolved is stored;
//   - every property reference: GETPROP, OPTCHAIN_GETPROP, unquoted
//     STRING_KEY and class or object members (piOCProp). The extern
//     properties are filtered when the list is read;
//   - two kinds of nodes to replay, because what they file depends on more
//     than their change scope or can fail: a SUPER (visitSuper files a
//     super() call under the enclosing class's extends clause, which lies in
//     another change scope) and a call of a property rename function (its
//     first argument's string).
//
// The walk is a NodeTraversal post-order with its class order, done directly
// (traverseNode); the index's records are built by a NodeTraversal per change
// scope, so splicing them gives the walk's order. Extern scripts are skipped
// by the walk when considerExterns is false, the only case read from the
// index, so their records hold no such items.
//
// A read files the items into dense per-id lists, then moves each list into
// the map once, in first-appearance order: the map, its id order and each
// id's nodes are the walk's. The passes are unchanged.
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
// resolving it.
func (r *piRecorder) ocName(n rhino.Node) {
	if id := n.GetStringID(); id > rhino.RhinoStringPool_EmptyID {
		r.lists[piListOC] = append(r.lists[piListOC], piItem{n: n, prop: id, bits: piOCName})
	}
}

// The verdicts of a read on a name id the global scope may declare.
const (
	ocVerdictUnknown uint8 = iota // not looked up yet
	ocVerdictKeep                 // a global, non-extern name
	ocVerdictDrop                 // not one
)

// ocFillTable is one table of the map (names or properties) as a read fills
// it: per-id lists in dense arrays, moved into the map once per id.
type ocFillTable struct {
	slot  []int32  // id -> 1 + its index in ids, 0 when not met in this read
	ids   []uint32 // the ids in first-appearance order
	count []int32  // per index: its nodes
	refs  []ocRef  // the nodes filed, in traversal order
}

type ocRef struct {
	slot int32
	n    rhino.Node
}

func (t *ocFillTable) add(id uint32, n rhino.Node) {
	if int(id) >= len(t.slot) {
		t.slot = slices.Grow(t.slot, int(id)+1-len(t.slot)+len(t.slot)/2)
		t.slot = t.slot[:cap(t.slot)]
	}
	s := t.slot[id]
	if s == 0 {
		t.ids = append(t.ids, id)
		t.count = append(t.count, 0)
		s = int32(len(t.ids))
		t.slot[id] = s
	}
	t.count[s-1]++
	t.refs = append(t.refs, ocRef{slot: s - 1, n: n})
}

// flush moves the lists into data, in first-appearance order, and resets the
// table for the next read.
func (t *ocFillTable) flush(m *OptimizeCalls_ReferenceMap, data *optimizeCalls_refTable) {
	data.byID = make(map[uint32]*optimizeCalls_refList, len(t.ids))
	data.order = append(make([]uint32, 0, len(t.ids)), t.ids...)
	data.size = len(t.ids)
	nodes := make([]rhino.Node, len(t.refs))
	// count becomes each list's end.
	end := int32(0)
	for i, c := range t.count {
		end += c
		t.count[i] = end - c // its start, advanced to its end by the scatter
	}
	for _, r := range t.refs {
		nodes[t.count[r.slot]] = r.n
		t.count[r.slot]++
	}
	start := int32(0)
	for i, id := range t.ids {
		end := t.count[i]
		refs := m.refLists.New()
		refs.nodes = nodes[start:end:end]
		data.byID[id] = refs
		t.slot[id] = 0
		start = end
	}
	t.ids = t.ids[:0]
	t.count = t.count[:0]
	t.refs = t.refs[:0]
}

// ocFillState is what a read of the index needs besides the records, kept
// across reads.
type ocFillState struct {
	names, props ocFillTable
	maybe        []uint64 // the name ids the global scope may declare (a bitset)
	verdict      []uint8  // name id -> ocVerdict*, for the ids in maybe
	marked       []uint32 // the ids set in maybe or verdict
	global       *Scope
	lookupAll    bool // the global scope has a parent: look every name up

	// Statistics (JSGO_PROPINDEX_STATS=1).
	reads int
	fill  time.Duration
}

func (st *ocFillState) mark(id uint32) {
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

// begin marks the names the global scope may declare: its vars and the
// implicit vars. GetVarByID finds nothing else in a scope without a parent
// (Scope.ownVar finds only a var of varList).
func (st *ocFillState) begin(gs *Scope) {
	st.global = gs
	st.lookupAll = gs.lookupParent != nil
	for _, v := range gs.varList {
		st.mark(v.nameID)
	}
	for _, id := range [...]uint32{implicitArgumentsID, implicitThisID, implicitExportsID, implicitSuperID} {
		st.mark(id)
	}
}

func (st *ocFillState) end() {
	for _, id := range st.marked {
		st.maybe[id>>6] = 0
		st.verdict[id] = ocVerdictUnknown
	}
	st.marked = st.marked[:0]
	st.global = nil
}

// isGlobalName is maybeAddNameReferenceNode's test with considerExterns
// false: the run's global scope declares id, and not in the externs.
func (st *ocFillState) isGlobalName(id uint32) bool {
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
	if v := st.global.GetVarByID(id); v != nil && !v.IsExtern() {
		keep = true
	}
	if keep {
		st.verdict[id] = ocVerdictKeep
	} else {
		st.verdict[id] = ocVerdictDrop
	}
	return keep
}

// fillOptimizeCalls adds to cb's map what cb.traverseNode(externs) and
// cb.traverseNode(root) would, with considerExterns false.
func (ix *propUseIndex) fillOptimizeCalls(cb *optimizeCallsRefCollector) {
	st := &ix.oc
	var t0 time.Time
	if propIndexStatsOn {
		t0 = time.Now()
	}
	st.begin(cb.globalScope)
	f := ocFiller{st: st, cb: cb}
	for _, top := range [2]rhino.Node{ix.externs, ix.root} {
		for script := top.GetFirstChild(); !script.IsNil(); script = script.GetNext() {
			if script.IsFromExterns() {
				continue // as traverseNode
			}
			f.walk(ix.scopes[script])
		}
	}
	st.end()
	if propIndexStatsOn {
		st.reads++
	}
	m := cb.references
	st.names.flush(m, &m.names)
	st.props.flush(m, &m.props)
	if propIndexStatsOn {
		d := time.Since(t0)
		st.fill += d
		fmt.Fprintf(os.Stderr, "[propindex] optimizeCalls read %d: %d items, %d names, %d properties, %v\n",
			st.reads, f.items, m.names.Size(), m.props.Size(), d.Round(10*time.Microsecond))
	}
}

type ocFiller struct {
	st    *ocFillState
	cb    *optimizeCallsRefCollector
	items int
}

// walk files the items of s and of its nested functions' records, spliced in
// where the traversal enters them.
func (f *ocFiller) walk(s *piScope) {
	st, cb := f.st, f.cb
	kids := s.kids
	k := 0
	items := s.items[piListOC]
	f.items += len(items)
	for i := range items {
		for k < len(kids) && int(kids[k].at[piListOC]) == i {
			f.walk(kids[k].s)
			k++
		}
		it := &items[i]
		switch it.bits {
		case piOCName:
			if st.isGlobalName(it.prop) {
				st.names.add(it.prop, it.n)
			}
		case piOCProp:
			if !cb.externPropBits.Get(int(it.prop)) {
				st.props.add(it.prop, it.n)
			}
		case piOCSuper:
			// visitSuper
			target, isName := optimizeCallsSuperTarget(it.n)
			if target.IsNil() {
				break
			}
			if isName {
				// maybeAddNameReferenceNode(target, superNode)
				if id := target.GetStringID(); id > rhino.RhinoStringPool_EmptyID && st.isGlobalName(id) {
					st.names.add(id, it.n)
				}
			} else if id := target.GetStringID(); !cb.externPropBits.Get(int(id)) {
				st.props.add(id, it.n)
			}
		case piOCRename:
			// traverseNode's CALL case
			firstArg, _ := ast.MustCall(it.n).FirstArg()
			if id := firstArg.Node().GetStringID(); !cb.externPropBits.Get(int(id)) {
				st.props.add(id, it.n)
			}
		}
	}
	for ; k < len(kids); k++ {
		f.walk(kids[k].s)
	}
}

// verifyOptimizeCallsMap panics if got, read from the index, differs from
// want, built by the walk: the ids of each table in first-appearance order,
// and each id's nodes.
func verifyOptimizeCallsMap(c *Compiler, got, want *OptimizeCalls_ReferenceMap) {
	cmp := func(kind string, g, w *optimizeCalls_refTable) {
		var gi, wi []uint32
		for id := range g.allByID() {
			gi = append(gi, id)
		}
		for id := range w.allByID() {
			wi = append(wi, id)
		}
		fail := func(format string, args ...any) {
			var b strings.Builder
			fmt.Fprintf(&b, "optimizeCalls' reference map read from the property use index differs from the walk's in %s: %s ", c.currentPassName(), kind)
			fmt.Fprintf(&b, format, args...)
			panic(b.String())
		}
		for i := 0; i < len(gi) || i < len(wi); i++ {
			if i >= len(gi) || i >= len(wi) || gi[i] != wi[i] {
				name := func(ids []uint32) string {
					if i < len(ids) {
						return strpool.String(strpool.ID(ids[i]))
					}
					return "(none)"
				}
				fail("id %d of %d read, %d walked: %q read, %q walked", i, len(gi), len(wi), name(gi), name(wi))
			}
		}
		for _, id := range gi {
			gn, wn := g.byID[id].nodes, w.byID[id].nodes
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
				fail("%q: node %d of %d read, %d walked: %s read, %s walked", strpool.String(strpool.ID(id)), j, len(gn), len(wn), at(gn), at(wn))
			}
		}
	}
	cmp("names", &got.names, &want.names)
	cmp("properties", &got.props, &want.props)
}
