package jscomp

import (
	"fmt"
	"os"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Go-only: RemoveUnusedCode keeps what it learns between runs.
//
// The pass runs up to a dozen times per compilation (in both code-removal
// loops, and once before and after them), and every run used to rescan the
// whole program even when a handful of functions had changed since the last
// one. What a run computes from the tree is a function of the tree alone, so
// the parts that depend only on code that has not changed can be kept: the
// compiler's ChangeTracker says which change scopes (functions and scripts)
// changed since a named point of its timeline, and every pass reports its
// changes there (PhaseOptimizer's change verifier checks it in the tests;
// peephole optimizations rely on it the same way).
//
// The guarded polyfill usages (PolyfillUsageFinder.traverseOnlyGuarded, a
// whole-program NodeTraversal at the first polyfill each run meets) are
// kept per script: a guard never reaches outside the script it is in (the
// guard state starts empty under ROOT and is keyed by nodes of the script),
// and the run only asks whether a node is a guarded usage. A run traverses
// a script only when it asks about a node in it and the script holds a
// scope changed since it was last traversed.
//
// Everything is dropped when the AST is compacted (node ids change) and
// when the root or the polyfill table differ. JSGO_RUC_INCR=0 turns it off.

var (
	rucIncrEnv = os.Getenv("JSGO_RUC_INCR")
	// JSGO_RUC_INCR_VERIFY=1 recomputes everything from scratch after each
	// incremental computation and panics on a difference (a debugging aid).
	rucIncrVerify = os.Getenv("JSGO_RUC_INCR_VERIFY") == "1"
)

func rucIncrEnabled() bool {
	return rucIncrEnv != "0"
}

const rucGuardedTimelineKey = "__removeUnusedCode_guardedUsages__"

// rucGuardedCache holds the guarded polyfill usages of every script of one
// root, per script.
type rucGuardedCache struct {
	root      rhino.Node
	polyfills *PolyfillUsageFinder_Polyfills
	compacts  int
	byScript  map[rhino.Node][]rhino.Node
	// usages is the union of byScript, the set a run consults.
	usages map[rhino.Node]bool
}

// incrCompiler returns the Compiler whose caches this run may use, or
// nil when the run must compute everything itself.
func (r *RemoveUnusedCode) incrCompiler() *Compiler {
	if !rucIncrEnabled() {
		return nil
	}
	c := r.compiler
	if c == nil || c.parallelRecorder != nil || c.changeTracker == nil {
		return nil
	}
	if r.root.IsNil() || !r.root.IsRoot() || r.root.GetParent().IsNil() {
		return nil
	}
	if !c.GetLifeCycleStage().IsNormalizedUnobfuscated() {
		return nil
	}
	return c
}

// isGuardedUsage reports whether n is a polyfill usage behind a guard.
//
// Go-only: the run used to compute the whole program's guarded usages at
// its first polyfill; it is only asked about the references of polyfills it
// may still remove, so a run that can keep them computes them per script,
// for the scripts it is asked about, at the first question.
func (r *RemoveUnusedCode) isGuardedUsage(n rhino.Node) bool {
	if !r.guardedUsagesComputed {
		r.computeGuardedUsages()
	}
	g := r.guardedCache
	if g == nil {
		return r.guardedUsages[n]
	}
	s := rucScriptUnder(n, r.root)
	if _, ok := g.byScript[s]; !ok && !s.IsNil() {
		var found []rhino.Node
		NewPolyfillUsageFinder(r.compiler, r.polyfillsFromTable).TraverseOnlyGuarded(s, func(u *PolyfillUsageFinder_PolyfillUsage) {
			found = append(found, u.Node())
			g.usages[u.Node()] = true
		})
		if found == nil {
			found = []rhino.Node{}
		}
		g.byScript[s] = found
	}
	guarded := g.usages[n]
	if rucIncrVerify {
		if r.guardedUsages == nil {
			r.guardedUsages = make(map[rhino.Node]bool)
			NewPolyfillUsageFinder(r.compiler, r.polyfillsFromTable).TraverseOnlyGuarded(r.root, r.StorePolyfill)
		}
		if guarded != r.guardedUsages[n] {
			panic(fmt.Sprintf("RemoveUnusedCode: guarded usage of %v: kept %v, full traversal %v", n, guarded, r.guardedUsages[n]))
		}
	}
	return guarded
}

// computeGuardedUsages fills r.guardedUsages with the polyfill usages of
// r.root that are behind a guard, or, when the run may keep them, drops the
// kept usages of the scripts that changed and sets r.guardedCache.
func (r *RemoveUnusedCode) computeGuardedUsages() {
	r.guardedUsagesComputed = true
	c := r.incrCompiler()
	if c == nil {
		NewPolyfillUsageFinder(r.compiler, r.polyfillsFromTable).
			TraverseOnlyGuarded(r.root, r.StorePolyfill)
		return
	}
	changed := c.changeTracker.GetChangedScopeNodesForPass(rucGuardedTimelineKey)
	g := c.rucGuarded
	if g == nil || changed == nil || g.root != r.root || g.compacts != c.astCompact.count || g.polyfills != r.polyfillsFromTable {
		g = &rucGuardedCache{
			root:      r.root,
			polyfills: r.polyfillsFromTable,
			compacts:  c.astCompact.count,
			byScript:  make(map[rhino.Node][]rhino.Node),
			usages:    make(map[rhino.Node]bool),
		}
		c.rucGuarded = g
	} else {
		for _, n := range changed {
			if s := rucScriptUnder(n, r.root); !s.IsNil() {
				g.forget(s)
			}
		}
		// order-independent: drops the entries of scripts no longer under the root
		for s := range g.byScript {
			if s.GetParent() != r.root {
				g.forget(s)
			}
		}
	}
	r.guardedCache = g
}

func (g *rucGuardedCache) forget(script rhino.Node) {
	for _, n := range g.byScript[script] {
		delete(g.usages, n)
	}
	delete(g.byScript, script)
}

// rucScriptUnder returns the child of root that contains n, or NilNode when
// n is no longer under root.
func rucScriptUnder(n, root rhino.Node) rhino.Node {
	for !n.IsNil() {
		p := n.GetParent()
		if p == root {
			return n
		}
		n = p
	}
	return rhino.NilNode
}
