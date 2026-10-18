package jscomp

import (
	"cmp"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// functionAnalysisRevision covers the function's position and every enclosing
// change scope. A serial outer rewrite can invalidate a nested function even
// when that nested function has not received a direct change report.
type functionAnalysisRevision struct {
	root   *rhino.Node
	parent *rhino.Node
	chain  []functionAnalysisScopeRevision
}

type functionAnalysisScopeRevision struct {
	root  *rhino.Node
	stamp int32
}

func snapshotFunctionAnalysisRevision(root *rhino.Node) functionAnalysisRevision {
	revision := functionAnalysisRevision{root: root, parent: root.GetParent()}
	for ancestor := root; ancestor != nil; ancestor = ancestor.GetParent() {
		if ChangeTracker_IsChangeScopeRoot(ancestor) {
			revision.chain = append(revision.chain, functionAnalysisScopeRevision{ancestor, ancestor.GetChangeTime()})
		}
	}
	return revision
}

func (r *functionAnalysisRevision) valid() bool {
	if r.root.IsDeleted() || r.root.GetParent() != r.parent {
		return false
	}
	index := 0
	for ancestor := r.root; ancestor != nil; ancestor = ancestor.GetParent() {
		if !ChangeTracker_IsChangeScopeRoot(ancestor) {
			continue
		}
		if index >= len(r.chain) || r.chain[index].root != ancestor || r.chain[index].stamp != ancestor.GetChangeTime() {
			return false
		}
		index++
	}
	return index == len(r.chain)
}

// A prepared job contains canonical, fully materialized lexical metadata and a
// private escaped set. Workers never touch a ScopeCreator or compiler mutation
// APIs. All AST reads finish before the compiler goroutine resumes rewriting.
type coalesceFunctionAnalysisJob struct {
	revision functionAnalysisRevision
	scope    *Scope
	vars     *NodeUtil_AllVarsDeclaredInFunction
}

type coalesceFunctionAnalysisResult struct {
	liveness *LiveVariablesAnalysis
	coloring GraphColoring[*Var, any]
}

// The analysis strategy receives only a prepared snapshot, private escaped state
// and the read-only compiler facade. It also makes the phase's concurrency and
// failure behavior testable with the real analysis behind a blocking wrapper.
type coalesceFunctionAnalyzer func(*coalesceFunctionAnalysisJob, *jutil.ObjectLinkedHashSet[*Var],
	*functionAnalysisReadOnlyCompiler) *coalesceFunctionAnalysisResult

// Only the small readonly compiler surface required by CFG construction is
// exposed to workers. An unexpected access fails locally and causes the original
// serial path to run at the original traversal position. In particular, workers
// cannot create compiler scopes, report diagnostics or allocate global names.
type functionAnalysisReadOnlyCompiler struct {
	AbstractCompiler
	options *CompilerOptions
}

func (c *functionAnalysisReadOnlyCompiler) GetOptions() *CompilerOptions       { return c.options }
func (c *functionAnalysisReadOnlyCompiler) GetSourceLine(*string, int) *string { return nil }
func (c *functionAnalysisReadOnlyCompiler) ThrowInternalError(message string, cause jutil.Throwable) {
	panic(jutil.NewRuntimeExceptionWithMessageAndCause(message, cause))
}

// coalesceFunctionAnalysisScheduler keeps at most one workers-sized window of
// private CFGs/lattices/colorings. Discovering canonical scopes and escape facts
// is serialized; computing the complete analysis is parallel; consuming results
// and AST mutation remain in the original traversal order.
type coalesceFunctionAnalysisScheduler struct {
	pass       *CoalesceVariableNames
	jobs       []coalesceFunctionAnalysisJob
	byRoot     map[*rhino.Node]int
	windowFrom int
	windowEnd  int
	results    []*coalesceFunctionAnalysisResult
	analyze    coalesceFunctionAnalyzer
	// Counters are caller-owned and support differential coverage of reuse and
	// conservative invalidation without worker-side writes to shared state.
	computed, reused, invalidated int
}

type coalesceFunctionAnalysisCollector struct {
	NodeTraversal_AbstractScopedCallbackBase
	scheduler *coalesceFunctionAnalysisScheduler
}

func (c *coalesceFunctionAnalysisCollector) Visit(*NodeTraversal, *rhino.Node, *rhino.Node) {}
func (c *coalesceFunctionAnalysisCollector) EnterScope(t *NodeTraversal) {
	// Materialize scopes in the caller, including block/nested scopes that
	// ComputeEscaped would otherwise create lazily. Keep only scope headers
	// here; per-function variable maps are built for the active window.
	scope := t.GetScope()
	if !scope.IsFunctionScope() {
		return
	}
	s := c.scheduler
	root := t.GetScopeRoot()
	s.byRoot[root] = len(s.jobs)
	s.jobs = append(s.jobs, coalesceFunctionAnalysisJob{
		revision: snapshotFunctionAnalysisRevision(root), scope: scope,
	})
}

func newCoalesceFunctionAnalysisScheduler(pass *CoalesceVariableNames, root *rhino.Node) *coalesceFunctionAnalysisScheduler {
	s := &coalesceFunctionAnalysisScheduler{pass: pass, byRoot: make(map[*rhino.Node]int), analyze: analyzeCoalesceFunctionSnapshot}
	NodeTraversal_Builder_().SetCompiler(pass.compiler).SetScopeCreator(pass.scopeCreator).
		SetCallback(&coalesceFunctionAnalysisCollector{scheduler: s}).Traverse(root)
	return s
}

func (s *coalesceFunctionAnalysisScheduler) take(root *rhino.Node) *coalesceFunctionAnalysisResult {
	index, ok := s.byRoot[root]
	if !ok {
		return nil
	}
	if index < s.windowFrom || index >= s.windowEnd {
		s.computeWindow(index)
	}
	result := s.results[index-s.windowFrom]
	s.results[index-s.windowFrom] = nil
	if result == nil {
		return nil
	}
	if !s.jobs[index].revision.valid() {
		s.invalidated++
		return nil
	}
	s.reused++
	return result
}

func (s *coalesceFunctionAnalysisScheduler) computeWindow(start int) {
	workers := max(1, s.pass.compiler.GetOptions().GetNumParallelThreads())
	count := min(workers, len(s.jobs)-start)
	s.windowFrom, s.windowEnd = start, start+count
	s.results = make([]*coalesceFunctionAnalysisResult, count)
	escaped := make([]*jutil.ObjectLinkedHashSet[*Var], count)
	for slot := range count {
		job := &s.jobs[start+slot]
		if !job.revision.valid() {
			s.invalidated++
			continue // The original traversal will recompute against its live state.
		}
		job.vars = NodeUtil_GetAllVarsDeclaredInFunction(s.pass.compiler, s.pass.scopeCreator, job.scope)
		if len(job.vars.GetAllVariablesInOrder()) >= LiveVariablesAnalysis_MAX_VARIABLES_TO_ANALYZE {
			job.vars = nil
			continue
		}
		// Escape discovery traverses nested scopes and resolves their names in
		// the caller. Workers subsequently only query explicit parameter names
		// or HasOwnSlot; neither path creates shared implicit scope variables.
		escaped[slot] = jutil.NewObjectLinkedHashSet[*Var]()
		DataFlowAnalysis_ComputeEscaped(job.scope, escaped[slot], s.pass.compiler, s.pass.scopeCreator,
			job.vars.GetAllVariables())
	}
	compiler := &functionAnalysisReadOnlyCompiler{options: s.pass.compiler.GetOptions()}
	failures := parallelPrebuild(count, workers, func(slot int) {
		if escaped[slot] == nil {
			return
		}
		s.results[slot] = s.analyze(&s.jobs[start+slot], escaped[slot], compiler)
	})
	for slot, failure := range failures {
		// Successful results own the variable maps through their liveness
		// instance. Release preparation-only headers for all other functions.
		s.jobs[start+slot].vars = nil
		if failure != nil {
			// Re-run the unchanged serial implementation at consumption, which
			// preserves its diagnostics, exception wrapping and pass ordering.
			s.results[slot] = nil
		}
		if escaped[slot] != nil {
			s.computed++
		}
	}
}

func analyzeCoalesceFunctionSnapshot(job *coalesceFunctionAnalysisJob, escaped *jutil.ObjectLinkedHashSet[*Var],
	compiler *functionAnalysisReadOnlyCompiler) *coalesceFunctionAnalysisResult {
	cfg := ControlFlowAnalysis_Builder_().SetCompiler(compiler).SetCfgRoot(job.revision.root).
		SetIncludeEdgeAnnotations(true).ComputeCfg()
	liveness := newLiveVariablesAnalysisFromSnapshot(cfg, job.scope, nil, job.vars, escaped)
	liveness.AddScopeVariables()
	if parser.FeatureSet_ES3.Contains(compiler.options.GetOutputFeatureSet()) &&
		NodeUtil_GetFunctionParameters(job.revision.root).HasTwoChildren() {
		liveness.MarkAllParametersEscaped()
	}
	liveness.Analyze()
	private := &CoalesceVariableNames{liveness: liveness}
	graph := private.ComputeVariableNamesInterferenceGraph(cfg, liveness.GetEscapedLocals())
	coloring := NewGraphColoring_GreedyGraphColoringWithTieBreaker[*Var, any](graph,
		func(a, b *Var) int {
			return cmp.Compare(liveness.GetVarIndex(a.GetName()), liveness.GetVarIndex(b.GetName()))
		})
	coloring.Color()
	return &coalesceFunctionAnalysisResult{liveness: liveness, coloring: coloring}
}
