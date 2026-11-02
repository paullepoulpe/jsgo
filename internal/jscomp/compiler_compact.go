package jscomp

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Go-only: AST compaction (rhino/node_compact.go).
//
// Nodes link with ids, so dead nodes are reclaimed only by compacting the
// node store: keeping its live nodes, under new ids, which compaction gives
// them in pre-order so that walks read the records nearly in address order.
// This file decides when (always after typesToColors, so that the
// optimizations run on that layout; then whenever the nodes allocated since
// the last compaction exceed a quarter of the live nodes, checked at the end
// of every optimization-loop iteration and after every top-level
// optimization pass, and at the end of the optimization stage) and keeps the
// forwarding registry: the explicit list
// of every rhino.Node the Compiler holds outside the tree across a pass
// boundary (docs/memory-research/node-holder-census.md, section 3.1).
// Everything else that holds nodes is pass-local and dead between passes.
//
// Compaction never runs before typesToColors (the type graph holds hundreds
// of thousands of nodes), nor when the typed scopes are kept
// (PreservesDetailedSourceInfo), nor in a NestedCompilerRunner's shadow
// compiler (its tree is spliced into, and held by, the outer compilation).
//
// JSGO_AST_COMPACT=0|1 overrides the default (on), JSGO_AST_COMPACT_LOG=1
// logs each compaction to stderr.
//
// Layout: JSGO_AST_LAYOUT=slide slides the live nodes down in id order
// instead of renumbering them in pre-order (rhino.Compactor.Permute), and
// skips the compaction after typesToColors when less than 10% of the nodes
// are dead, as before the pre-order layout.
// JSGO_AST_LAYOUT_AT=<point,...> also forces a compaction at those points
// (pass names, "loopL.I" for iteration I of loop L) whatever the dead
// fraction.

type astCompaction struct {
	enabled     bool
	nested      bool // a NestedCompilerRunner compiler: shares nodes with the outer compilation, never compacts
	allowed     bool // inside the optimization stage
	preorder    bool // renumber the live nodes in pre-order instead of sliding them
	nodesAtLast int  // nodes in the store right after the last compaction
	liveAtLast  int
	count       int
}

var (
	astCompactEnv  = os.Getenv("JSGO_AST_COMPACT")
	astCompactLog  = os.Getenv("JSGO_AST_COMPACT_LOG") == "1"
	astLayoutEnv   = os.Getenv("JSGO_AST_LAYOUT")
	astLayoutForce = splitNonEmpty(os.Getenv("JSGO_AST_LAYOUT_AT"))
)

func splitNonEmpty(s string) []string {
	var out []string
	for _, x := range strings.Split(s, ",") {
		if x = strings.TrimSpace(x); x != "" {
			out = append(out, x)
		}
	}
	return out
}

// astLayoutForced reports whether a compaction is forced at point.
func astLayoutForced(point string) bool {
	for _, x := range astLayoutForce {
		if x == point {
			return true
		}
	}
	return false
}

// astCompactMinNodes is the allocation volume below which a compaction is
// never worth its walk.
const astCompactMinNodes = 64 * 1024

func (c *Compiler) initASTCompaction() {
	on := true
	switch astCompactEnv {
	case "0":
		on = false
	case "1":
		on = true
	}
	c.astCompact = astCompaction{
		enabled:  on && !c.astCompact.nested && !c.options.PreservesDetailedSourceInfo(),
		nested:   c.astCompact.nested,
		preorder: astLayoutEnv != "slide",
	}
}

func (c *Compiler) beginASTCompaction() {
	a := &c.astCompact
	a.allowed = true
	if a.nodesAtLast == 0 && c.nodeStore != nil {
		n := c.nodeStore.NodeCount()
		a.nodesAtLast = n
		a.liveAtLast = n
	}
}

// maybeCompactAST compacts the AST after a pass when the policy says so.
// externs and root are the caller's local root pointers (forwarded in
// place); extra forwards holders the caller keeps in locals.
func (c *Compiler) maybeCompactAST(after string, externs, root *rhino.Node, extra func(fwd func(rhino.Node) rhino.Node)) {
	a := &c.astCompact
	if !a.enabled || !a.allowed || c.externAndJsRoot.IsNil() {
		return
	}
	if !c.GetLifeCycleStage().HasColorAndSimplifiedJSDoc() || c.typeRegistry != nil {
		return
	}
	force := a.preorder && (after == "typesToColors" || astLayoutForced(after))
	if after != "typesToColors" && !force {
		allocated := c.nodeStore.NodeCount() - a.nodesAtLast
		if allocated < max(a.liveAtLast/4, astCompactMinNodes) {
			return
		}
	}
	c.compactAST(after, force, externs, root, extra)
}

// compactAST keeps only the live nodes of the store (rhino/node_compact.go):
// it marks what every registered holder refers to, renumbers the live nodes
// in pre-order (or slides them down), and forwards the holders to the
// nodes' new ids. Unless force is set, it does nothing when less than 10%
// of the nodes are dead.
func (c *Compiler) compactAST(after string, force bool, externs, root *rhino.Node, extra func(fwd func(rhino.Node) rhino.Node)) {
	a := &c.astCompact
	t0 := time.Now()
	jutil.CheckStateMsg(c.typedScopeCreator == nil && c.topScope == nil, "AST compaction with typed scopes alive")
	k := rhino.NewCompactor()
	if a.preorder {
		k.RecordPreorder()
	}
	c.forwardASTHolders(k.Mark, externs, root, extra)
	tMarked := time.Now()
	total := k.Total()
	live := k.CountMarked()
	if dead := total - live; !force && (dead <= 0 || int64(dead)*10 < int64(total)) {
		k.PruneDeadAux()
		k.Finish()
		a.nodesAtLast = c.nodeStore.NodeCount()
		a.liveAtLast = live
		if astCompactLog {
			fmt.Fprintf(os.Stderr, "[ast-compact skipped after %s] kept %d of %d nodes (<10%% dead), %v\n", after, live, total, time.Since(t0).Round(time.Millisecond))
		}
		return
	}
	threads := 1
	if c.options != nil && c.options.GetNumParallelThreads() > 1 {
		threads = c.options.GetNumParallelThreads()
	}
	if a.preorder {
		k.Permute(threads)
	} else {
		k.SlideWithThreads(threads)
	}
	tMoved := time.Now()
	c.forwardASTHolders(k.Forward, externs, root, extra)
	k.Finish()

	a.count++
	a.nodesAtLast = c.nodeStore.NodeCount()
	a.liveAtLast = k.Live
	if astCompactLog {
		layout := "slide"
		if a.preorder {
			layout = fmt.Sprintf("pre-order, %d new blocks, %d reused, copy %v, tables %v", k.NewBlocks, k.ReusedBlocks,
				k.CopyTime.Round(time.Millisecond), k.TablesTime.Round(time.Millisecond))
		}
		fmt.Fprintf(os.Stderr, "[ast-compact #%d after %s] %s, %s, %v (mark %v, move %v, forward %v)\n", a.count, after, k, layout,
			time.Since(t0).Round(time.Millisecond), tMarked.Sub(t0).Round(time.Millisecond), tMoved.Sub(tMarked).Round(time.Millisecond), time.Since(tMoved).Round(time.Millisecond))
	}
}

// forwardASTHolders replaces every node the compiler holds outside the tree
// across a pass boundary (and the caller's locals externs, root and those
// extra forwards) with fwd(node). compactAST calls it twice: to mark the
// nodes, then to forward them.
func (c *Compiler) forwardASTHolders(fwd func(rhino.Node) rhino.Node, externs, root *rhino.Node, extra func(fwd func(rhino.Node) rhino.Node)) {
	// 1. Roots.
	c.externAndJsRoot = fwd(c.externAndJsRoot)
	c.externsRoot = fwd(c.externsRoot)
	c.jsRoot = fwd(c.jsRoot)
	if !externs.IsNil() {
		*externs = fwd(*externs)
	}
	if !root.IsNil() {
		*root = fwd(*root)
	}
	for name, n := range c.scriptNodeByFilename { // order-independent: updating existing keys in place without adding or removing keys
		c.scriptNodeByFilename[name] = fwd(n)
	}
	forwardInput := func(in *CompilerInput) {
		if in == nil {
			return
		}
		jutil.CheckState(in.typedScope == nil)
		if in.ast != nil {
			in.ast.root = fwd(in.ast.root)
		}
	}
	for _, in := range c.inputsById { // order-independent: k.Mark sets bitset bits and k.Forward maps precomputed old->new node IDs
		forwardInput(in)
	}
	if c.chunkGraph != nil {
		for _, in := range c.chunkGraph.GetAllInputs() {
			forwardInput(in)
		}
	}
	for _, in := range c.externs {
		forwardInput(in)
	}
	if c.phaseOptimizer != nil {
		c.phaseOptimizer.jsRoot = fwd(c.phaseOptimizer.jsRoot)
	}

	// 2. Change tracking.
	if c.changeTracker != nil && c.changeTracker.changeTimeline != nil {
		c.changeTracker.changeTimeline.ForwardValues(fwd)
	}

	// 3. Modules.
	if c.moduleMetadataMap != nil {
		var md func(m *ModuleMetadataMap_ModuleMetadata)
		md = func(m *ModuleMetadataMap_ModuleMetadata) {
			m.rootNode = fwd(m.rootNode)
			for _, nm := range m.nestedModules {
				md(nm)
			}
		}
		for m := range c.moduleMetadataMap.GetAllModuleMetadata().All() {
			md(m)
		}
	}
	if c.moduleMap != nil {
		seen := make(map[*Module]struct{})
		exp := func(e *Export) {
			if e != nil {
				e.exportNode = fwd(e.exportNode)
				e.nameNode = fwd(e.nameNode)
			}
		}
		bind := func(b *Binding) {
			if b != nil {
				b.sourceNode = fwd(b.sourceNode)
				exp(b.originatingExport)
			}
		}
		mod := func(m *Module) {
			if m == nil {
				return
			}
			if _, done := seen[m]; done {
				return
			}
			seen[m] = struct{}{}
			if m.namespace != nil {
				for _, b := range m.namespace.All() {
					bind(b)
				}
			}
			if m.boundNames != nil {
				for _, b := range m.boundNames.All() {
					bind(b)
				}
			}
			if m.localNameToLocalExport != nil {
				for _, e := range m.localNameToLocalExport.All() {
					exp(e)
				}
			}
		}
		for _, m := range c.moduleMap.GetModulesByPath().All() {
			mod(m)
		}
		for _, m := range c.moduleMap.GetModulesByClosureNamespace().All() {
			mod(m)
		}
	}

	// 4. Diagnostics.
	if c.errorManager != nil {
		for _, e := range c.errorManager.GetErrors() {
			e.node = fwd(e.node)
		}
		for _, e := range c.errorManager.GetWarnings() {
			e.node = fwd(e.node)
		}
	}

	// 5. Runtime library injection.
	if c.runtimeJsLibManager != nil {
		c.runtimeJsLibManager.lastInjectedLibrary = fwd(c.runtimeJsLibManager.lastInjectedLibrary)
	}

	// 6. Transpilation namespace.
	if gn, ok := c.transpilationNamespace.(*GlobalNamespace); ok && gn != nil {
		gn.forwardNodes(fwd)
	}

	// 7. Caller locals.
	if extra != nil {
		extra(fwd)
	}
}

// forwardNodes forwards every node the namespace holds (roots, Name
// declarations and refs, the refs' scopes, and the two node-keyed caches).
func (g *GlobalNamespace) forwardNodes(fwd func(rhino.Node) rhino.Node) {
	g.root = fwd(g.root)
	g.externsRoot = fwd(g.externsRoot)
	g.globalRoot = fwd(g.globalRoot)
	if g.spreadSiblingCache != nil {
		m := make(map[rhino.Node]bool, len(g.spreadSiblingCache))
		// order-independent: rebuilds the map with forwarded keys.
		for n, v := range g.spreadSiblingCache {
			m[fwd(n)] = v
		}
		g.spreadSiblingCache = m
	}
	if g.defineCandidateScripts != nil {
		m := make(map[rhino.Node]struct{}, len(g.defineCandidateScripts))
		// order-independent: rebuilds the set with forwarded keys.
		for n := range g.defineCandidateScripts {
			m[fwd(n)] = struct{}{}
		}
		g.defineCandidateScripts = m
	}
	scopes := make(map[*Scope]struct{})
	var scope func(s *Scope)
	scope = func(s *Scope) {
		for ; s != nil; s = s.parent {
			if _, done := scopes[s]; done {
				return
			}
			scopes[s] = struct{}{}
			s.rootNode = fwd(s.rootNode)
			for _, v := range s.varList {
				if v != nil {
					v.nameNode = fwd(v.nameNode)
				}
			}
		}
	}
	for i := uint32(0); i < g.refCount; i++ {
		r := &g.refPages[i>>gnRefPageShift][i&gnRefPageMask]
		r.node = fwd(r.node)
		scope(r.scope)
	}
	for i := range g.multiRefs {
		mr := &g.multiRefs[i]
		if mr.refsIndex != nil {
			m := make(map[rhino.Node]int, len(mr.refsIndex))
			for idx, id := range mr.refIDs {
				if id != 0 {
					r := g.refByID(id)
					if !r.node.IsNil() {
						m[r.node] = idx
					}
				}
			}
			mr.refsIndex = m
		}
	}
}

// ForwardValues replaces every value in the timeline through f (the index
// map is rebuilt, as it hashes by value).
func (tl *Timeline[T]) ForwardValues(f func(T) T) {
	index := make(map[T]uint32, len(tl.eventsByValue))
	// Event ids and their order are unchanged; only values are rewritten.
	// order-independent: rebuilds the index with forwarded values.
	for value, id := range tl.eventsByValue {
		nv := f(value)
		tl.event(id).value = nv
		index[nv] = id
	}
	tl.eventsByValue = index
}
