package jscomp

import (
	"fmt"
	"testing"
	"time"

	"github.com/paullepoulpe/jsgo/internal/jutil"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

func TestFunctionAnalysisRevisionIncludesAncestorsAndPosition(t *testing.T) {
	inner := rhino.IR_Function(rhino.IR_Name("inner"), rhino.IR_ParamList(), rhino.IR_Block())
	outer := rhino.IR_Function(rhino.IR_Name("outer"), rhino.IR_ParamList(), rhino.IR_BlockOfNodes(inner))
	script := rhino.IR_ScriptOfNodes(outer)
	revision := snapshotFunctionAnalysisRevision(inner)
	if !revision.valid() {
		t.Fatal("fresh snapshot invalid")
	}
	outer.SetChangeTime(1)
	if revision.valid() {
		t.Fatal("outer rewrite did not invalidate inner summary")
	}
	revision = snapshotFunctionAnalysisRevision(inner)
	script.SetChangeTime(2)
	if revision.valid() {
		t.Fatal("script rewrite did not invalidate nested summary")
	}
	tracker := NewChangeTracker()
	tracker.ReportChangeToChangeScope(outer)
	revision = snapshotFunctionAnalysisRevision(inner)
	tracker.ReportChangeToChangeScope(outer)
	if revision.valid() {
		t.Fatal("a repeated outer change report reused the old revision")
	}
	revision = snapshotFunctionAnalysisRevision(inner)
	inner.Detach()
	if revision.valid() {
		t.Fatal("detaching a function did not invalidate its summary")
	}
}

func TestImmutableFunctionAnalysisMatchesOriginalCoalescing(t *testing.T) {
	t.Setenv("JSGO_EXPERIMENT_FUNCTION_ANALYSIS", "")
	programs := []string{
		`function f(a) { var x = a + 1; use(x); var y = a + 2; return y; }
		 function g(b) { var q = b * 2; use(q); var r = b * 3; return r; }`,
		`function outer(a) { var x = 1; use(x); var y = 2; use(y);
		 function inner(b) { var z = b + 1; use(z); var q = a + 2; return q; } return inner(y); }`,
		`function f(a, b) { var x = a + b; use(x); var y = a - b; return y; }
		 function g(a, b) { arguments[0] = 4; var x = a + 1; use(x); var y = b + 1; return y; }`,
		`function f(a) { try { var x = a + 1; use(x); throw a; }
		 catch (e) { let q = e + 1; use(q); } finally { use(a); }
		 var y = a + 2; return y; }`,
		`function f(a) { var x = 1; while (a--) { let y = x + a; use(y); }
		 var z = 2; return z; }
		 function g(a) { var x = 1; for (var k in a) { use(k, x); } var y = 2; return y; }`,
		`function f(a = 1) { let x = a + 1; use(x); let y = a + 2; return y; }
		 function g({a, b}) { let x = a + b; use(x); let y = a - b; return y; }`,
		`function f(a) { const inner = function(b) { let x = b + a; use(x); let y = b - a; return y; };
		 let x = 1; use(x); let y = 2; return inner(y); }`,
		`function f(a) { var x = function named(b) { var q = b + 1; use(q); var r = b + 2; return r; };
		 use(x); var y = function(c) { var q = c + 1; use(q); var r = c + 2; return r; }; use(y); return a; }`,
		`function f(a) { class C { field = () => a;
		 m(b = a) { let x = b + 1; use(x); let y = b + 2; return y; } }
		 let x = 1; use(x); let y = 2; return new C().m(y); }`,
		`function outer(a, b = a) { function named(c = b) { let x = c + 1; use(x); let y = c + 2; return y; }
		 try { throw a; } catch (e) { let x = 1; use(x); let y = 2; use(y); return () => named(e); } }`,
		`function f(a) { let x = 1; use(x); { let y = 2; use(y); }
		 class C { m(b) { let x = b + 1; use(x); let y = b + 2; return y; } }
		 let z = a + 3; return new C().m(z); }`,
	}
	for index, program := range programs {
		for _, pseudo := range []bool{false, true} {
			t.Run(fmt.Sprintf("case%d/pseudo=%v", index, pseudo), func(t *testing.T) {
				render := func(workers int, experimental bool) (string, *coalesceFunctionAnalysisScheduler) {
					options := NewCompilerOptions()
					options.SetNumParallelThreads(workers)
					options.SetLanguageIn(CompilerOptions_LanguageMode_ECMASCRIPT_NEXT)
					compiler := NewCompiler()
					compiler.Init(nil, []*SourceFile{SourceFile_FromCode("functions.js", program)}, options)
					compiler.ParseInputs()
					root := compiler.GetJsRoot()
					Normalize_CreateNormalizeForOptimizations(compiler).Process(compiler.GetExternsRoot(), root)
					pass := NewCoalesceVariableNames(compiler, pseudo)
					var scheduler *coalesceFunctionAnalysisScheduler
					if experimental {
						scheduler = newCoalesceFunctionAnalysisScheduler(pass, root)
						pass.functionAnalyses = scheduler
						NodeTraversal_Builder_().SetCompiler(compiler).SetCallback(pass).
							SetScopeCreator(pass.scopeCreator).Traverse(root)
						compiler.SetLifeCycleStage(AbstractCompiler_LifeCycleStage_RAW)
					} else {
						pass.Process(compiler.GetExternsRoot(), root)
					}
					if compiler.GetErrorCount() != 0 {
						t.Fatalf("unexpected compile error in case %d", index)
					}
					return compiler.ToSource(), scheduler
				}
				want, _ := render(1, false)
				for _, workers := range []int{1, 2, 4} {
					got, scheduler := render(workers, true)
					if got != want {
						t.Fatalf("workers=%d changed coalescing:\n%s\nwant:\n%s", workers, got, want)
					}
					if scheduler.computed == 0 || scheduler.reused == 0 {
						t.Fatalf("workers=%d did not consume actual worker summaries: computed=%d reused=%d", workers, scheduler.computed, scheduler.reused)
					}
					if len(scheduler.results) > workers {
						t.Fatalf("retained %d private results with only %d workers", len(scheduler.results), workers)
					}
				}
			})
		}
	}
}

func TestFunctionAnalysisDiscardsStaleResultBeforePublication(t *testing.T) {
	function := rhino.IR_Function(rhino.IR_Name("f"), rhino.IR_ParamList(), rhino.IR_Block())
	rhino.IR_ScriptOfNodes(function)
	result := &coalesceFunctionAnalysisResult{}
	scheduler := &coalesceFunctionAnalysisScheduler{
		jobs:   []coalesceFunctionAnalysisJob{{revision: snapshotFunctionAnalysisRevision(function)}},
		byRoot: map[*rhino.Node]int{function: 0}, windowEnd: 1,
		results: []*coalesceFunctionAnalysisResult{result},
	}
	function.SetChangeTime(1)
	if got := scheduler.take(function); got != nil || scheduler.invalidated != 1 || scheduler.reused != 0 {
		t.Fatalf("stale result published: result=%v invalidated=%d reused=%d", got, scheduler.invalidated, scheduler.reused)
	}
	if scheduler.results[0] != nil {
		t.Fatal("stale private result retained")
	}
}

func TestFunctionAnalysisRealJobsOverlapWithinWindow(t *testing.T) {
	options := NewCompilerOptions()
	options.SetNumParallelThreads(2)
	compiler := NewCompiler()
	compiler.Init(nil, []*SourceFile{SourceFile_FromCode("overlap.js", `
	 function f(a) { var x = a+1; use(x); var y = a+2; return y; }
	 function g(a) { var x = a+1; use(x); var y = a+2; return y; }
	 function h(a) { var x = a+1; use(x); var y = a+2; return y; }
	 function i(a) { var x = a+1; use(x); var y = a+2; return y; }
	`)}, options)
	compiler.ParseInputs()
	Normalize_CreateNormalizeForOptimizations(compiler).Process(compiler.GetExternsRoot(), compiler.GetJsRoot())
	pass := NewCoalesceVariableNames(compiler, false)
	scheduler := newCoalesceFunctionAnalysisScheduler(pass, compiler.GetJsRoot())
	entered := make(chan struct{}, 4)
	release := make(chan struct{})
	released := false
	defer func() {
		if !released {
			close(release)
		}
	}()
	scheduler.analyze = func(job *coalesceFunctionAnalysisJob, escaped *jutil.ObjectLinkedHashSet[*Var],
		compiler *functionAnalysisReadOnlyCompiler) *coalesceFunctionAnalysisResult {
		entered <- struct{}{}
		<-release
		return analyzeCoalesceFunctionSnapshot(job, escaped, compiler)
	}
	done := make(chan any, 1)
	go func() {
		defer func() { done <- recover() }()
		scheduler.computeWindow(0)
	}()
	for range 2 {
		select {
		case <-entered:
		case failure := <-done:
			t.Fatalf("window ended before two real analysis jobs overlapped: %v", failure)
		case <-time.After(5 * time.Second):
			t.Fatal("analysis jobs did not overlap")
		}
	}
	select {
	case <-entered:
		t.Fatal("more analysis jobs started than the worker/window budget")
	default:
	}
	close(release)
	released = true
	select {
	case failure := <-done:
		if failure != nil {
			t.Fatalf("analysis window failed: %v", failure)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("analysis jobs did not finish")
	}
	if len(scheduler.results) != 2 || scheduler.results[0] == nil || scheduler.results[1] == nil {
		t.Fatal("real CFG/liveness/coloring results were not produced")
	}
}
