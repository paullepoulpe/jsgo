package jscomp

import (
	"os"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

func parallelFSIVEnabled() bool { return os.Getenv("JSGO_EXPERIMENT_PARALLEL_FSIV") == "1" }

// FSIV rewrites only function-local definitions and uses. Scripts can therefore
// run independently, provided their scope chains share the ONE canonical global
// scope of the real combined ROOT. Copying that scope or using SCRIPT as its root
// would change global binding identity (including escaped-variable decisions).
//
// Every job retains original traversal/analysis/rewrite order inside its script:
// an outer function can move a nested function before that nested function is
// visited. There is no speculative function analysis or AST cloning here.
// This is an opt-in CLI experiment. On an unexpected internal exception later
// jobs can already have rewritten their scripts; the earliest exception and its
// diagnostic prefix are deterministic, but partial in-memory ASTs after failure
// are not an API-equivalence guarantee. Custom observers always run serially.
type parallelFSIVJob struct {
	script   *rhino.Node
	compiler parallelFSIVCompiler
}

type parallelFSIVEvent struct {
	changeScope *rhino.Node
	diagnostic  *JSError
}

// Deliberately leave the embedded interface nil. Only the audited read surface
// is available; a future pass dependency cannot accidentally mutate the shared
// compiler via a promoted method. Compiler writes are published after the join.
type parallelFSIVCompiler struct {
	AbstractCompiler
	owner      *Compiler
	options    *CompilerOptions
	convention CodingConvention
	analyzer   *AstAnalyzer
	events     []parallelFSIVEvent
}

func (c *parallelFSIVCompiler) GetOptions() *CompilerOptions { return c.options }
func (c *parallelFSIVCompiler) GetInput(id *rhino.InputId) *CompilerInput {
	return c.owner.GetInput(id)
}
func (c *parallelFSIVCompiler) GetCodingConvention() CodingConvention {
	return c.convention
}
func (c *parallelFSIVCompiler) GetAstAnalyzer() *AstAnalyzer { return c.analyzer }
func (c *parallelFSIVCompiler) GetSourceLine(name *string, line int) *string {
	// SourceFile publishes immutable source/line-offset snapshots atomically.
	return c.owner.GetSourceLine(name, line)
}
func (c *parallelFSIVCompiler) ThrowInternalError(message string, cause jutil.Throwable) {
	c.owner.ThrowInternalError(message, cause) // constructs/throws, no compiler writes
}
func (c *parallelFSIVCompiler) Report(diagnostic *JSError) {
	c.events = append(c.events, parallelFSIVEvent{diagnostic: diagnostic})
}
func (c *parallelFSIVCompiler) ReportChangeToEnclosingScope(n *rhino.Node) {
	// Resolve NOW: FSIV reports an assignment parent before detaching it. Walking
	// its ancestors at publication time would lose the original change scope.
	c.ReportChangeToChangeScope(c.owner.changeTracker.getChangeScopeForNode(n))
}
func (c *parallelFSIVCompiler) ReportChangeToChangeScope(scope *rhino.Node) {
	jutil.CheckState(scope.IsScript() || scope.IsFunction())
	c.events = append(c.events, parallelFSIVEvent{changeScope: scope})
}

func parallelFSIVConventionSafe(convention CodingConvention) bool {
	switch c := convention.(type) {
	case *CodingConventions_DefaultCodingConvention:
		return c != nil
	case *ClosureCodingConvention:
		return c != nil && parallelFSIVConventionSafe(c.nextConvention)
	case *GoogleCodingConvention:
		return c != nil && parallelFSIVConventionSafe(c.nextConvention)
	default:
		return false
	}
}

func parallelFSIVErrorManagerSafe(manager ErrorManager) bool {
	switch m := manager.(type) {
	case *ThreadSafeDelegatingErrorManager:
		return m != nil && parallelFSIVErrorManagerSafe(m.delegated)
	case *LoggerErrorManager, *PrintStreamErrorManager, *SortingErrorManagerBase:
		return !jutil.IsNil(manager)
	default:
		return false
	}
}

func parallelFSIVObserversSafe(c *Compiler) bool {
	o := c.GetOptions()
	if _, ok := SyntacticScopeCreator_DEFAULT_REDECLARATION_HANDLER.(*SyntacticScopeCreator_DefaultRedeclarationHandler); !ok {
		return false
	}
	if o.GetErrorHandler() != nil || o.GetDebugLogDirectory() != nil ||
		o.GetCheckDeterminism() || o.GetTracerMode() != CompilerOptions_TracerMode_OFF ||
		!parallelFSIVConventionSafe(c.GetCodingConvention()) || !parallelFSIVErrorManagerSafe(c.errorManager) {
		return false
	}
	for _, handler := range c.changeTracker.codeChangeHandlers.Slice() {
		switch handler.(type) {
		case *RecentChange, *PhaseOptimizer_ScopedChangeHandler:
			// These consume only ordered change stamps, never the live AST. FSIV
			// does not query change stamps while it is running.
		default:
			return false
		}
	}
	return true
}

func parallelFSIVHasRawTypes(n *rhino.Node) bool {
	if n.GetJSType() != nil {
		return true
	}
	for child := n.GetFirstChild(); child != nil; child = child.GetNext() {
		if parallelFSIVHasRawTypes(child) {
			return true
		}
	}
	return false
}

func (f *FlowSensitiveInlineVariables) prepareParallelScripts(externs, root *rhino.Node) ([]parallelFSIVJob, *Scope) {
	c, ok := f.compiler.(*Compiler)
	if !ok || c == nil || c.GetOptions().GetNumParallelThreads() < 2 ||
		externs == nil || root == nil || externs.GetParent() == nil ||
		root.GetParent() != externs.GetParent() || !root.IsRoot() ||
		!externs.GetParent().IsRoot() || !parallelFSIVObserversSafe(c) {
		return nil, nil
	}
	if first := root.GetFirstChild(); first == nil || first.GetNext() == nil {
		return nil, nil
	}
	jobs := []parallelFSIVJob{}
	analyzer := c.GetAstAnalyzer() // Snapshot settings and shared readonly summaries serially.
	for script := root.GetFirstChild(); script != nil; script = script.GetNext() {
		if !script.IsScript() || script.GetInputId() == nil {
			return nil, nil
		}
		input := c.GetInput(script.GetInputId())
		if input == nil || input.IsExtern() {
			return nil, nil
		}
		// Raw JSType equality can resolve/memoize shared types. Color equality
		// against StandardColors_STRING (AstAnalyzer's sole color comparison)
		// has empty comparison sets, hence never calls lazy Color.HashCode.
		if analyzer.useTypesForLocalOptimization && parallelFSIVHasRawTypes(script) {
			return nil, nil
		}
		copy := *analyzer
		jobs = append(jobs, parallelFSIVJob{script: script, compiler: parallelFSIVCompiler{
			owner: c, options: c.GetOptions(), convention: c.GetCodingConvention(), analyzer: &copy,
		}})
	}
	if len(jobs) < 2 {
		return nil, nil
	}
	// This is exactly the global root and scanner used by TraverseRoots. Its
	// declarations/name nodes remain untouched; global implicit variables cannot
	// exist at a ROOT. Var's only lazy hash is atomic. Descendant scopes and their
	// implicit slots are created privately by each script's normal scope creator.
	global := NewSyntacticScopeCreator(c).CreateScopeOfScope(externs.GetParent(), nil)
	return jobs, global
}

func (f *FlowSensitiveInlineVariables) processScriptsParallel(externs, root *rhino.Node) bool {
	jobs, global := f.prepareParallelScripts(externs, root)
	if global == nil {
		return false
	}
	failures := parallelPrebuild(len(jobs), f.compiler.GetOptions().GetNumParallelThreads(), func(index int) {
		job := &jobs[index]
		pass := NewFlowSensitiveInlineVariables(&job.compiler)
		pass.nativeFactsEnabled = f.nativeFactsEnabled
		NodeTraversal_Builder_().SetCompiler(&job.compiler).SetCallback(pass).
			TraverseWithScope(job.script, global)
	})
	publishParallelFSIVJobs(f.compiler, jobs, failures)
	return true
}

func publishParallelFSIVJobs(compiler AbstractCompiler, jobs []parallelFSIVJob, failures []jutil.Throwable) {
	// Reproduce serial change stamps/timeline insertion and diagnostics. If an
	// unexpected internal error occurs, other jobs may finish, but only the
	// original input-order prefix (including events before that error) publishes.
	for index := range jobs {
		for _, event := range jobs[index].compiler.events {
			if event.diagnostic != nil {
				compiler.Report(event.diagnostic)
			} else {
				compiler.ReportChangeToChangeScope(event.changeScope)
			}
		}
		if failures[index] != nil {
			panic(failures[index])
		}
	}
}
