package jscomp

import (
	"fmt"
	"os"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/strpool"
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

const (
	rucGuardedTimelineKey  = "__removeUnusedCode_guardedUsages__"
	rucPolyDefsTimelineKey = "__removeUnusedCode_polyDefs__"
)

type rucPolyDefsCache struct {
	root        rhino.Node
	compacts    int
	byScript    map[rhino.Node]bool
	polyScripts int
}

// rucGuardedCache holds the guarded polyfill usages of every script of one
// root, per script.
type rucGuardedCache struct {
	root          rhino.Node
	polyfills     *PolyfillUsageFinder_Polyfills
	disambiguate  bool
	colorRegistry *colors.ColorRegistry
	compacts      int
	byScript      map[rhino.Node][]rhino.Node
	// usages is the union of byScript, the set a run consults.
	usages map[rhino.Node]bool

	// Go-only: per-script $jscomp.polyfill definition presence cache.
	polyDefs rucPolyDefsCache

	// Go-only: function-level summary cache across RemoveUnusedCode runs
	// (remove_unused_code_summary.go).
	fnCache *rucFnSummaryCache
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
		r.polyfillUsageFinder.TraverseOnlyGuarded(s, func(u *PolyfillUsageFinder_PolyfillUsage) {
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
			r.polyfillUsageFinder.TraverseOnlyGuarded(r.root, r.StorePolyfill)
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
		r.polyfillUsageFinder.TraverseOnlyGuarded(r.root, r.StorePolyfill)
		return
	}
	changed := c.changeTracker.GetChangedScopeNodesForPass(rucGuardedTimelineKey)
	g := c.rucGuarded
	if g == nil {
		g = &rucGuardedCache{}
		c.rucGuarded = g
	}
	if g.byScript == nil || changed == nil || g.root != r.root || g.compacts != c.astCompact.count || g.polyfills != r.polyfillsFromTable || g.disambiguate != r.disambiguateProperties || g.colorRegistry != r.polyfillUsageFinder.registry {
		g.root = r.root
		g.polyfills = r.polyfillsFromTable
		g.disambiguate = r.disambiguateProperties
		g.colorRegistry = r.polyfillUsageFinder.registry
		g.compacts = c.astCompact.count
		g.byScript = make(map[rhino.Node][]rhino.Node)
		g.usages = make(map[rhino.Node]bool)
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

func (r *RemoveUnusedCode) initPolyfillSuppression(root rhino.Node) {
	if !r.removeUnusedPolyfills || r.removalLog.IsLogging() || r.unremovableLog.IsLogging() {
		return
	}
	r.root = root
	c := r.incrCompiler()
	if c == nil || root != c.GetJsRoot() {
		return
	}
	if !r.hasAnyPolyfillDef(c, root) {
		r.suppressPolyfillRemoval = true
	}
}

func (r *RemoveUnusedCode) hasAnyPolyfillDef(c *Compiler, root rhino.Node) bool {
	changed := c.changeTracker.GetChangedScopeNodesForPass(rucPolyDefsTimelineKey)
	g := c.rucGuarded
	if g == nil {
		g = &rucGuardedCache{}
		c.rucGuarded = g
	}
	p := &g.polyDefs
	if p.byScript == nil || changed == nil || p.root != root || p.compacts != c.astCompact.count {
		p.root = root
		p.compacts = c.astCompact.count
		p.byScript = make(map[rhino.Node]bool)
		p.polyScripts = 0
	} else {
		for _, n := range changed {
			if s := rucScriptUnder(n, root); !s.IsNil() {
				if had, exists := p.byScript[s]; exists {
					if had {
						p.polyScripts--
					}
					delete(p.byScript, s)
				}
			}
		}
		// order-independent: removes detached scripts and subtracts 1 per script that had a polyfill def
		for s, had := range p.byScript {
			if s.GetParent() != root {
				if had {
					p.polyScripts--
				}
				delete(p.byScript, s)
			}
		}
	}

	var toScan []rhino.Node
	for s := root.GetFirstChild(); !s.IsNil(); s = s.GetNext() {
		if _, ok := p.byScript[s]; !ok {
			toScan = append(toScan, s)
		}
	}
	if len(toScan) > 0 {
		results := make([]bool, len(toScan))
		if numWorkers := parallelThreadCount(c, len(toScan)); numWorkers > 1 {
			_ = c.ensureSync()
			if wp := spawnIndexedWorkers(len(toScan), numWorkers, func(_ int, claim func() (int, bool)) {
				for {
					idx, ok := claim()
					if !ok {
						return
					}
					results[idx] = subtreeHasJscompPolyfillDef(toScan[idx])
				}
			}); wp != nil {
				rethrowParallelPanic(c, wp)
			}
		} else {
			for i, s := range toScan {
				results[i] = subtreeHasJscompPolyfillDef(s)
			}
		}
		for i, s := range toScan {
			p.byScript[s] = results[i]
			if results[i] {
				p.polyScripts++
			}
		}
	}

	hasAny := p.polyScripts > 0
	if rucIncrVerify {
		want := subtreeHasJscompPolyfillDef(root)
		if hasAny != want {
			panic(fmt.Sprintf("RemoveUnusedCode: hasAnyPolyfillDef=%v, full scan=%v", hasAny, want))
		}
	}
	return hasAny
}

func subtreeHasJscompPolyfillDef(n rhino.Node) bool {
	tok := n.GetToken()
	if tok == rhino.Token_CALL || tok == rhino.Token_OPTCHAIN_CALL {
		if _, ok := isJscompPolyfillCall(ast.MustNormalOrOptChainCall(n)); ok {
			return true
		}
	}
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		if !c.GetFirstChild().IsNil() && subtreeHasJscompPolyfillDef(c) {
			return true
		}
	}
	return false
}

func countNodeChildren(n rhino.Node) int {
	if n.IsNil() {
		return 0
	}
	count := 0
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		count++
	}
	return count
}

func (r *RemoveUnusedCode) assertParamScopeNoPendingChanges(fparamScope *Scope, label string) {
	fnNode := fparamScope.GetRootNode()
	if NodeUtil_IsGetOrSetKey(fnNode.GetParent()) {
		return
	}
	params := ast.MustFunction(fnNode).Params()
	if last, ok := params.LastParam(); ok {
		if argNode, removable := r.trailingParamBinding(last); removable {
			if pat, isPattern := ast.AsPattern(argNode); isPattern {
				if !pat.Node().HasChildren() {
					panic(fmt.Sprintf("RemoveUnusedCode verify (%s): trailing empty pattern param removable in %v", label, fnNode))
				}
			} else if r.isParamVarRemovable(argNode, fparamScope, true) {
				panic(fmt.Sprintf("RemoveUnusedCode verify (%s): trailing param %s removable in %v", label, ast.MustName(argNode).Ident(), fnNode))
			}
		}
	}
	for p := range params.Params() {
		if p.IsUnused() {
			continue
		}
		paramName, hasName := rucParamName(p)
		if !hasName {
			continue
		}
		if r.isParamVarRemovable(paramName.Node(), fparamScope, true) {
			panic(fmt.Sprintf("RemoveUnusedCode verify (%s): unmarked unused param %s in %v", label, paramName.Ident(), fnNode))
		}
	}
}

func (r *RemoveUnusedCode) verifyIdempotentSkip(root rhino.Node) {
	if !rucIncrVerify {
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
	ref.root = root
	ref.trackProps = ref.removesAnyProperties()
	if ref.trackProps {
		for _, id := range rucImplicitPropIDs {
			ref.pinnedProps.Set(int(id))
		}
	}
	ref.varByKey = make(map[scopedVarKey]*RemoveUnusedCode_RealVarInfo, 256)
	scope := ref.createGlobalScope(root)
	ref.worklist.push(RemoveUnusedCode_Continuation{node: root, scope: scope})
	for !ref.worklist.isEmpty() {
		cont := ref.worklist.pop()
		cont.Apply(ref)
	}

	for _, vi := range ref.varInfos {
		if !vi.IsRemovable() {
			continue
		}
		for cur := vi.removables.head; cur != nil; cur = cur.nextInQueue {
			if cur.remKind == remKindIndirectAssign &&
				cur.indirectAssignRoot().GetParent().IsParamList() &&
				!cur.indirectAssignRoot().IsDefaultValue() {
				continue
			}
			panic(fmt.Sprintf("RemoveUnusedCode verify (idempotent skip): removable var %s (%s)", vi.var_.GetName(), cur.String()))
		}
		toRemove := vi.var_.GetNameNode().GetParent()
		if fn, isFn := ast.AsFunction(toRemove); isFn && NodeUtil_IsFunctionExpression(fn.Node()) && !ref.preserveFunctionExpressionNames {
			panic(fmt.Sprintf("RemoveUnusedCode verify (idempotent skip): removable fn expr name %s", vi.var_.GetName()))
		}
	}
	if ref.trackProps {
		for _, idx := range ref.propRemOrder {
			entry := &ref.propEntries[idx]
			if !entry.removables.isEmpty() {
				panic(fmt.Sprintf("RemoveUnusedCode verify (idempotent skip): removable prop %s", strpool.String(entry.propID)))
			}
		}
	}
	for _, poly := range ref.polyfillOrder {
		if poly.isRemovable {
			panic(fmt.Sprintf("RemoveUnusedCode verify (idempotent skip): removable polyfill %s", poly.getName()))
		}
	}
	if ref.removeGlobals {
		for _, fparamScope := range ref.allFunctionParamScopes {
			ref.assertParamScopeNoPendingChanges(fparamScope, "idempotent skip")
		}
	}
}
