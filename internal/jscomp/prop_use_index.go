package jscomp

import (
	"fmt"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

// Go-only: an index of property uses and definitions, kept current across the
// optimization loop.
//
// InlineSimpleMethods and InlineProperties run once per iteration of the main
// optimization loop, and each used to walk the whole program twice per run (a
// gathering walk, then a rewriting walk), although the nodes they act on are a
// small part of the tree and only a few percent of it changes between runs.
// The index holds, for every change scope (a script or a function, the unit
// the ChangeTracker reports), what those walks act on, in the order a
// whole-program NodeTraversal meets it; the passes replay their callbacks from
// it instead of walking. DeadPropertyAssignmentElimination reads its change
// scopes, in the traversal's order, instead of walking to find them.
//
// What it holds. A change scope's record has one list per reader, each in the
// scope's NodeTraversal post-order with nested functions left out (see the
// piList* and pi* constants): the GETPROPs; the CALLs with a GETPROP callee;
// what GatherSignatures and GatherCandidates find in the scope, recorded by
// running their own decision code on it when it depends on the scope alone,
// or else the node to run it on. It also says where each nested function's
// traversal starts in each list (piKid). A whole-program traversal is the
// scripts of each root in order, with each nested function's record spliced
// in at its place, so a replay meets the nodes in exactly the order of the
// walk it replaces. The records are built with a NodeTraversal per change
// scope, so they follow NodeTraversal's own order (classes, computed members)
// by construction. What a record holds depends only on its scope's subtree;
// what else the passes read (definitions' bodies, colors, JSDoc, ancestors
// outside the scope) is read when they replay.
//
// How it stays current. Each read first takes the change scopes reported to
// the ChangeTracker since the previous read (a timeline mark, as
// RemoveUnusedCode's guarded polyfill usages do), and indexes again those
// that are still in the tree, with the parent of a new function; a nested
// function met for the first time is indexed too. Records of deleted or
// detached scopes are dropped. Everything is rebuilt on the first read, after
// an AST compaction (node ids change) and when the roots change, and dropped
// when the loop ends: the index is used only inside a PhaseOptimizer loop,
// where every change is made by a loop pass and reported (ChangeVerifier
// checks the reports in the tests; its notion of a scope's content, the scope
// without its nested functions' insides, is the index's).
//
// JSGO_PROPINDEX=0 turns it off, JSGO_PROPINDEX_DPAE=0 for
// DeadPropertyAssignmentElimination alone. JSGO_PROPINDEX_VERIFY=1 rebuilds it
// from scratch after every loop pass and before every read and panics on a
// difference, naming the pass that changed the tree without reporting it.
// JSGO_PROPINDEX_STATS=1 prints what it does.

var (
	propIndexEnv    = os.Getenv("JSGO_PROPINDEX")
	propIndexVerify = os.Getenv("JSGO_PROPINDEX_VERIFY") == "1"
	// propIndexForce makes the passes use the index outside a loop too
	// (tests: every run builds it afresh).
	propIndexForce = false
	// propIndexReads counts the reads (tests).
	propIndexReads atomic.Int64
)

const propIndexTimelineKey = "__propUseIndex__"

// The bits of an index item: what it is to its reader.
const (
	// piUse: a GETPROP (InlineProperties' ReplaceCandidates).
	piUse uint8 = 1 << iota
	// piCall: a CALL whose callee is a GETPROP; prop is the callee's
	// (InlineSimpleMethods' InlineTrivialAccessors).
	piCall
	// piISM: a node InlineSimpleMethods' GatherSignatures acts on in a way
	// that depends on more than its change scope (replayed).
	piISM
	// piIP: a definition InlineProperties' GatherCandidates evaluates.
	piIP
	// piHoist: the traversal was in the global hoist scope at the node
	// (NodeTraversal.inGlobalHoistScope, read by GatherCandidates).
	piHoist
	// piISMEv: what GatherSignatures found, recorded: markNonInlineable(prop)
	// when n is nil, addSignature(prop, n) when n is a function.
	piISMEv
	// piIPInv: GatherCandidates invalidates prop here.
	piIPInv
)

var piTokens = func() ast.TokenSet {
	var s ast.TokenSet
	for _, t := range []rhino.Token{
		rhino.Token_GETPROP, rhino.Token_GETELEM, rhino.Token_OBJECTLIT, rhino.Token_CLASS_MEMBERS,
		rhino.Token_CALL, rhino.Token_STRING_KEY, rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF,
		rhino.Token_MEMBER_FUNCTION_DEF, rhino.Token_MEMBER_FIELD_DEF,
	} {
		s.Add(t)
	}
	return s
}()

type piItem struct {
	n    rhino.Node
	prop uint32
	bits uint8
}

// A record keeps its items in one list per reader, so that a reader scans
// only its own.
const (
	piListUses  = iota // GETPROPs (piUse)
	piListCalls        // CALLs with a GETPROP callee (piCall)
	piListISM          // GatherSignatures: findings (piISMEv) and nodes to replay (piISM)
	piListIP           // GatherCandidates: invalidations (piIPInv) and definitions (piIP)
	piLists
)

// piKid is a nested function whose traversal starts before items[l][at[l]]
// in each list l.
type piKid struct {
	at [piLists]int32
	s  *piScope
}

type piScope struct {
	root  rhino.Node
	items [piLists][]piItem
	kids  []piKid
	gen   uint32 // the update that last indexed it
	mark  uint32 // the last sweep that reached it
}

type propUseIndex struct {
	c             *Compiler
	externs, root rhino.Node
	compacts      int
	scopes        map[rhino.Node]*piScope
	gen           uint32
	sweeps        uint32
	live          int // records reached by the last sweep

	conv  CodingConvention
	cur   *piScope
	queue []*piScope

	// Statistics (JSGO_PROPINDEX_STATS=1).
	stats propIndexStats
}

type propIndexStats struct {
	builds, updates, reindexed, newScopes, dropped, items int
}

var propIndexStatsOn = os.Getenv("JSGO_PROPINDEX_STATS") == "1"

func propIndexEnabled() bool {
	return propIndexEnv != "0"
}

// propUseIndexFor returns the index, brought up to date, when a pass running
// on externs and root may read it, or nil when the pass must walk.
func propUseIndexFor(compiler AbstractCompiler, externs, root rhino.Node) *propUseIndex {
	if !propIndexEnabled() {
		return nil
	}
	c := compiler
	if c == nil || c.parallelRecorder != nil || c.changeTracker == nil {
		return nil
	}
	if externs.IsNil() || root.IsNil() || c.externAndJsRoot.IsNil() ||
		externs.GetParent() != c.externAndJsRoot || root.GetParent() != c.externAndJsRoot {
		return nil
	}
	if c.phaseOptimizer == nil || !c.phaseOptimizer.inLoop {
		if !propIndexForce {
			return nil
		}
		// Outside a loop nothing says what changed since the last read.
		c.dropPropUseIndex()
	}
	ix := c.refreshPropUseIndex(externs, root)
	if propIndexVerify {
		ix.verify("before " + c.currentPassName())
	}
	propIndexReads.Add(1)
	return ix
}

var propIndexDPAE = os.Getenv("JSGO_PROPINDEX_DPAE") != "0"

// propUseIndexForDPAE is propUseIndexFor for DeadPropertyAssignmentElimination,
// which walks the program root.
func propUseIndexForDPAE(compiler AbstractCompiler, externs, root rhino.Node) *propUseIndex {
	if !propIndexDPAE || root != compiler.GetJsRoot() {
		return nil
	}
	return propUseIndexFor(compiler, externs, root)
}

func (c *Compiler) currentPassName() string {
	if p := c.phaseOptimizer; p != nil && p.currentPass != nil {
		return p.currentPass.name
	}
	return "?"
}

// refreshPropUseIndex brings the compiler's index up to date, or builds it.
func (c *Compiler) refreshPropUseIndex(externs, root rhino.Node) *propUseIndex {
	changed := c.changeTracker.GetChangedScopeNodesForPass(propIndexTimelineKey)
	ix := c.propIndex
	if ix == nil || changed == nil || ix.externs != externs || ix.root != root || ix.compacts != c.astCompact.count {
		if ix != nil && propIndexStatsOn {
			ix.printStats()
		}
		ix = newPropUseIndex(c, externs, root)
		ix.build()
		c.propIndex = ix
		return ix
	}
	ix.update(changed)
	return ix
}

// dropPropUseIndex forgets the index (the loop that used it ended).
func (c *Compiler) dropPropUseIndex() {
	if c.propIndex == nil {
		return
	}
	if propIndexStatsOn {
		c.propIndex.printStats()
	}
	c.propIndex = nil
	if c.changeTracker != nil {
		c.changeTracker.changeTimeline.Unmark(propIndexTimelineKey)
	}
}

// propIndexAfterPass checks the index against a rebuild after a loop pass
// (JSGO_PROPINDEX_VERIFY=1).
func (c *Compiler) propIndexAfterPass(name string) {
	ix := c.propIndex
	if ix == nil || c.changeTracker == nil {
		return
	}
	ix = c.refreshPropUseIndex(ix.externs, ix.root)
	ix.verify("after " + name)
}

func newPropUseIndex(c *Compiler, externs, root rhino.Node) *propUseIndex {
	return &propUseIndex{
		c:        c,
		externs:  externs,
		root:     root,
		compacts: c.astCompact.count,
		scopes:   make(map[rhino.Node]*piScope, 1024),
		conv:     c.GetCodingConvention(),
	}
}

func (ix *propUseIndex) printStats() {
	s := ix.stats
	var byBit [8]int
	total := 0
	for _, top := range [2]rhino.Node{ix.externs, ix.root} {
		for l := 0; l < piLists; l++ {
			ix.each(top, l, nil, func(_ *piScope, it piItem) {
				total++
				for b := 0; b < 8; b++ {
					if it.bits&(1<<b) != 0 {
						byBit[b]++
					}
				}
			})
		}
	}
	fmt.Fprintf(os.Stderr, "[propindex] builds %d updates %d reindexed %d new %d dropped %d scopes %d items %d (use %d call %d ism %d ip %d hoist %d ismEv %d ipInv %d)\n",
		s.builds, s.updates, s.reindexed, s.newScopes, s.dropped, len(ix.scopes), total, byBit[0], byBit[1], byBit[2], byBit[3], byBit[4], byBit[5], byBit[6])
}

// build indexes every script of both roots.
func (ix *propUseIndex) build() {
	ix.gen++
	ix.stats.builds++
	for _, top := range [2]rhino.Node{ix.externs, ix.root} {
		for script := top.GetFirstChild(); !script.IsNil(); script = script.GetNext() {
			ix.index(ix.scopeFor(script))
		}
	}
	ix.drain()
	ix.live = len(ix.scopes)
}

// update indexes again the change scopes reported since the last read.
func (ix *propUseIndex) update(changed []rhino.Node) {
	ix.gen++
	ix.stats.updates++
	if propIndexStatsOn {
		t0, re0, it0 := time.Now(), ix.stats.reindexed, ix.stats.items
		defer func() {
			fmt.Fprintf(os.Stderr, "[propindex] update before %s: %d changed, %d scopes / %d items indexed, %v\n",
				ix.c.currentPassName(), len(changed), ix.stats.reindexed-re0, ix.stats.items-it0, time.Since(t0).Round(10*time.Microsecond))
		}()
	}
	for _, n := range changed {
		if !ix.attached(n) {
			if _, ok := ix.scopes[n]; ok {
				delete(ix.scopes, n)
				ix.stats.dropped++
			}
			continue
		}
		_, known := ix.scopes[n]
		ix.index(ix.scopeFor(n))
		// A new function: its parent scope lists where it starts. The
		// ChangeVerifier does not ask for the parent's report when the new
		// function replaced another one in place.
		for !known && n.IsFunction() {
			n = ChangeTracker_GetEnclosingChangeScopeRoot(n.GetParent())
			_, known = ix.scopes[n]
			ix.index(ix.scopeFor(n))
		}
	}
	// A script that was never reported (none is added in the loops today).
	for _, top := range [2]rhino.Node{ix.externs, ix.root} {
		for script := top.GetFirstChild(); !script.IsNil(); script = script.GetNext() {
			if _, ok := ix.scopes[script]; !ok {
				ix.index(ix.scopeFor(script))
			}
		}
	}
	ix.drain()
	if len(ix.scopes) > 2*ix.live+4096 {
		ix.sweep()
	}
}

func (ix *propUseIndex) scopeFor(n rhino.Node) *piScope {
	s := ix.scopes[n]
	if s == nil {
		s = &piScope{root: n}
		ix.scopes[n] = s
	}
	return s
}

// attached reports whether the change scope root n is still in the tree.
func (ix *propUseIndex) attached(n rhino.Node) bool {
	if n.IsDeleted() {
		return false
	}
	top := ix.c.externAndJsRoot
	for p := n.GetParent(); !p.IsNil(); p = p.GetParent() {
		if p == top {
			return true
		}
	}
	return false
}

// index lists the items of s (once per update).
func (ix *propUseIndex) index(s *piScope) {
	if s.gen == ix.gen {
		return
	}
	s.gen = ix.gen
	for l := range s.items {
		s.items[l] = s.items[l][:0]
	}
	s.kids = s.kids[:0]
	ix.cur = s
	ix.stats.reindexed++
	NodeTraversal_Traverse(ix.c, s.root, (*piRecorder)(ix))
	ix.cur = nil
	for _, list := range s.items {
		ix.stats.items += len(list)
	}
}

// drain indexes the nested functions met for the first time.
func (ix *propUseIndex) drain() {
	for len(ix.queue) > 0 {
		s := ix.queue[len(ix.queue)-1]
		ix.queue = ix.queue[:len(ix.queue)-1]
		ix.index(s)
	}
}

// sweep drops the records no longer reachable from the roots.
func (ix *propUseIndex) sweep() {
	ix.sweeps++
	m := ix.sweeps
	var reach func(s *piScope)
	reach = func(s *piScope) {
		s.mark = m
		for _, k := range s.kids {
			reach(k.s)
		}
	}
	for _, top := range [2]rhino.Node{ix.externs, ix.root} {
		for script := top.GetFirstChild(); !script.IsNil(); script = script.GetNext() {
			if s := ix.scopes[script]; s != nil {
				reach(s)
			}
		}
	}
	live := 0
	// order-independent: deletes the records the roots do not reach
	for n, s := range ix.scopes {
		if s.mark != m {
			delete(ix.scopes, n)
			ix.stats.dropped++
		} else {
			live++
		}
	}
	ix.live = live
}

// piRecorder is the NodeTraversal callback that indexes one change scope.
type piRecorder propUseIndex

func (r *piRecorder) visitTokens() *ast.TokenSet { return &piTokens }

func (r *piRecorder) shouldTraverseOnlyScriptAndFunction() bool { return true }

func (r *piRecorder) ShouldTraverse(_ *NodeTraversal, n, _ rhino.Node) bool {
	s := r.cur
	if n == s.root || !n.IsFunction() {
		return true
	}
	ix := (*propUseIndex)(r)
	kid := ix.scopes[n]
	if kid == nil {
		kid = &piScope{root: n}
		ix.scopes[n] = kid
		ix.queue = append(ix.queue, kid)
		ix.stats.newScopes++
	}
	k := piKid{s: kid}
	for l, list := range s.items {
		k.at[l] = int32(len(list))
	}
	s.kids = append(s.kids, k)
	return false
}

func (r *piRecorder) Visit(t *NodeTraversal, n, parent rhino.Node) {
	items := &r.cur.items
	switch n.GetToken() {
	case rhino.Token_GETPROP:
		prop := n.GetStringID()
		items[piListUses] = append(items[piListUses], piItem{n: n, prop: prop, bits: piUse})
		switch _, kind := inlineClassifyGather(r.c, n, parent); kind {
		case inlineGatherInvalidate:
			items[piListIP] = append(items[piListIP], piItem{n: n, prop: prop, bits: piIPInv})
		case inlineGatherGetpropDef:
			bits := piIP
			if t.InGlobalHoistScope() {
				bits |= piHoist
			}
			items[piListIP] = append(items[piListIP], piItem{n: n, prop: prop, bits: bits})
		}
		r.ismGatherGet(n, parent)
	case rhino.Token_GETELEM:
		r.ismGatherGet(n, parent)
	case rhino.Token_OBJECTLIT, rhino.Token_CLASS_MEMBERS:
		if ismMemberKeysKnown(n) {
			ismGatherMemberContainer(r, n)
		} else {
			items[piListISM] = append(items[piListISM], piItem{n: n, bits: piISM})
		}
	case rhino.Token_STRING_KEY, rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF,
		rhino.Token_MEMBER_FUNCTION_DEF, rhino.Token_MEMBER_FIELD_DEF:
		switch prop, kind := inlineClassifyGather(r.c, n, parent); kind {
		case inlineGatherInvalidate:
			items[piListIP] = append(items[piListIP], piItem{n: n, prop: prop, bits: piIPInv})
		case inlineGatherFieldDef:
			items[piListIP] = append(items[piListIP], piItem{n: n, prop: prop, bits: piIP})
		}
	case rhino.Token_CALL:
		call := ast.MustCall(n)
		if callee := call.Callee().Node(); callee.IsGetProp() {
			items[piListCalls] = append(items[piListCalls], piItem{n: n, prop: callee.GetStringID(), bits: piCall})
		}
		ismGatherCall(r, r.conv, call)
		if prop, kind := inlineClassifyGather(r.c, n, parent); kind == inlineGatherInvalidate {
			items[piListIP] = append(items[piListIP], piItem{n: n, prop: prop, bits: piIPInv})
		}
	}
}

// ismGatherGet records what GatherSignatures.visit finds at the GETPROP or
// GETELEM n, or, when that depends on more than n's change scope (the parent
// of an arrow function whose body is `x.prototype.y`: processPrototypeParent
// reads the parent's grandparent), n itself, to be replayed.
func (r *piRecorder) ismGatherGet(n, parent rhino.Node) {
	prop, ok := ismStaticPropertyID(n)
	if !ok {
		return
	}
	switch parent.GetToken() {
	case rhino.Token_ASSIGN:
	case rhino.Token_GETPROP, rhino.Token_GETELEM:
		if prop != ismPrototypeID {
			return
		}
		if parent.GetParent() == r.cur.root {
			r.cur.items[piListISM] = append(r.cur.items[piListISM], piItem{n: n, prop: prop, bits: piISM})
			return
		}
	default:
		return
	}
	ismGatherNormalGet(r, n)
}

// The index records what GatherSignatures finds, to be replayed
// (ismSignatureSink).
func (r *piRecorder) addPossibleSignature(propID uint32, node rhino.Node) {
	fn := rhino.NilNode
	if f, ok := ast.AsFunction(node); ok { // as addPossibleSignature
		fn = f.Node()
	}
	r.cur.items[piListISM] = append(r.cur.items[piListISM], piItem{n: fn, prop: propID, bits: piISMEv})
}

func (r *piRecorder) markNonInlineable(propID uint32) {
	r.cur.items[piListISM] = append(r.cur.items[piListISM], piItem{prop: propID, bits: piISMEv})
}

// each calls fn for every item of list under top (the externs or the
// program root), in the order of a NodeTraversal of top. With a non-nil props,
// only items whose prop is in it.
func (ix *propUseIndex) each(top rhino.Node, list int, props *inlineBitSet, fn func(s *piScope, it piItem)) {
	var walk func(s *piScope)
	walk = func(s *piScope) {
		kids := s.kids
		k := 0
		for i, it := range s.items[list] {
			for k < len(kids) && int(kids[k].at[list]) == i {
				walk(kids[k].s)
				k++
			}
			if props == nil || props.Get(int(it.prop)) {
				fn(s, it)
			}
		}
		for ; k < len(kids); k++ {
			walk(kids[k].s)
		}
	}
	for script := top.GetFirstChild(); !script.IsNil(); script = script.GetNext() {
		walk(ix.scopes[script])
	}
}

// eachScope calls fn for every change scope under top, in the order a
// NodeTraversal of top enters them (pre-order).
func (ix *propUseIndex) eachScope(top rhino.Node, fn func(s *piScope)) {
	var walk func(s *piScope)
	walk = func(s *piScope) {
		fn(s)
		for _, k := range s.kids {
			walk(k.s)
		}
	}
	for script := top.GetFirstChild(); !script.IsNil(); script = script.GetNext() {
		walk(ix.scopes[script])
	}
}

// enclosingFunction is NodeTraversal.getEnclosingFunction at a node of s.
func (s *piScope) enclosingFunction() rhino.Node {
	if s.root.IsFunction() {
		return s.root
	}
	return rhino.NilNode
}

// verify compares the index with a fresh build and panics on a difference.
func (ix *propUseIndex) verify(when string) {
	fresh := newPropUseIndex(ix.c, ix.externs, ix.root)
	fresh.build()
	type entry struct {
		it    piItem
		owner rhino.Node
	}
	flat := func(x *propUseIndex) []entry {
		var out []entry
		for _, top := range [2]rhino.Node{x.externs, x.root} {
			for l := 0; l < piLists; l++ {
				x.each(top, l, nil, func(s *piScope, it piItem) {
					out = append(out, entry{it, s.root})
				})
			}
		}
		return out
	}
	got, want := flat(ix), flat(fresh)
	for i := 0; i < len(got) || i < len(want); i++ {
		var g, w entry
		if i < len(got) {
			g = got[i]
		}
		if i < len(want) {
			w = want[i]
		}
		if g != w {
			var b strings.Builder
			fmt.Fprintf(&b, "property use index differs from a rebuild %s (item %d of %d kept, %d rebuilt)\n", when, i, len(got), len(want))
			fmt.Fprintf(&b, "  kept:    %s\n", piDescribe(g.it, g.owner))
			fmt.Fprintf(&b, "  rebuilt: %s\n", piDescribe(w.it, w.owner))
			panic(b.String())
		}
	}
}

func piDescribe(it piItem, owner rhino.Node) string {
	if it.n.IsNil() {
		return "(none)"
	}
	return fmt.Sprintf("%s %q bits %05b at %s, in %s", it.n.GetToken(), strpool.String(strpool.ID(it.prop)), it.bits,
		piPosition(it.n), piPosition(owner))
}

func piPosition(n rhino.Node) string {
	if n.IsNil() {
		return "nil"
	}
	file := "?"
	if f := n.GetSourceFileName(); f != nil {
		file = *f
	}
	desc := n.GetToken().String()
	if n.IsFunction() {
		if name := NodeUtil_GetName(n); name != nil {
			desc += " " + *name
		}
	}
	return fmt.Sprintf("%s %s:%d:%d (deleted %v, change time %d)", desc, file, n.GetLineno(), n.GetCharno(), n.IsDeleted(), n.GetChangeTime())
}

// piReplay reports a panic of a replayed callback as the traversal it
// replaces does (NodeTraversal.throwUnexpectedException at the visited node):
//
//	rp := piReplay{compiler: c}
//	defer rp.done()
//	... rp.cur = node; callback(node) ...
type piReplay struct {
	compiler AbstractCompiler
	cur      rhino.Node
}

func (rp *piReplay) done() {
	if r := recover(); r != nil {
		t := acquireNodeTraversalDirect(rp.compiler, NodeTraversal_AbstractPostOrderCallbackInterface(func(*NodeTraversal, rhino.Node, rhino.Node) {}), nil, false)
		t.currentNode = rp.cur
		t.currentScript = NodeUtil_GetEnclosingScript(rp.cur)
		t.ThrowUnexpectedException(jutil.AsThrowable(r))
	}
}
