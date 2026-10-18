package jscomp

import (
	"fmt"
	"slices"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

func newParallelFSIVTestCompiler(t *testing.T, workers int, programs ...string) *Compiler {
	t.Helper()
	options := NewCompilerOptions()
	options.SetLanguageIn(CompilerOptions_LanguageMode_ECMASCRIPT_NEXT)
	options.SetNumParallelThreads(workers)
	options.SetUseTypesForLocalOptimization(true)
	sources := make([]*SourceFile, len(programs))
	for i, program := range programs {
		sources[i] = SourceFile_FromCode(fmt.Sprintf("script%d.js", i), program)
	}
	compiler := NewCompiler()
	compiler.Init([]*SourceFile{SourceFile_FromCode("externs.js", `var use; var mutate; var outside;`)}, sources, options)
	compiler.ParseInputs()
	if compiler.GetErrorCount() != 0 {
		t.Fatal("fixture failed parsing")
	}
	Normalize_CreateNormalizeForOptimizations(compiler).Process(compiler.GetExternsRoot(), compiler.GetJsRoot())
	NewGatherGetterAndSetterProperties(compiler).Process(compiler.GetExternsRoot(), compiler.GetJsRoot())
	NewPureFunctionIdentifier_Driver(compiler).Process(compiler.GetExternsRoot(), compiler.GetJsRoot())
	if compiler.GetErrorCount() != 0 {
		t.Fatal("fixture failed normalization or side-effect setup")
	}
	compiler.changeTracker.GetChangedScopeNodesForPass("parallel-fsiv-test")
	return compiler
}

func parallelFSIVChangeSnapshot(c *Compiler) []string {
	out := []string{fmt.Sprint(c.changeTracker.GetChangeStamp())}
	var walk func(*rhino.Node)
	walk = func(n *rhino.Node) {
		if n.IsFunction() || n.IsScript() {
			out = append(out, fmt.Sprintf("%s:%d", n.String(), n.GetChangeTime()))
		}
		for child := n.GetFirstChild(); child != nil; child = child.GetNext() {
			walk(child)
		}
	}
	walk(c.GetRoot())
	for _, node := range c.changeTracker.GetChangedScopeNodesForPass("parallel-fsiv-test") {
		out = append(out, node.String())
	}
	return out
}

func TestParallelFSIVDifferential(t *testing.T) {
	fixtures := [][]string{
		{`var shared = 1; function f(a) { var x; x = a + 1; return x + shared; }`,
			`function g(a) { var x = a + 2; return x + shared; }`},
		{`function f(a) { var fn = function(b) { var x = a + b; return x; }; return fn; }`,
			`function g(a) { var fn = () => { let x = a + 2; return x; }; return fn; }`},
		{`function f(a) { var x = a; if (a) { x = 2; } return x; }`,
			`function g(a) { var x = a; while (a--) { use(x); x = a; } return x; }`},
		{`function f(a) { const x = a + 1; return x; }`,
			`function g(a) { var x; label: x = a + 1; return x; }`},
		{`function f(a) { var x = a; try { mutate(); } catch(e) { let y = e; use(y); } finally { use(a); } return x; }`,
			`function g(a) { var x = a; use(() => x); var y = a + 1; return y; }`},
		{`function f(a) { let x = a; [x] = outside; return x; }`,
			`function g(a) { let {x, y = a} = outside; var z = y; return z; }`},
		{`function f(a) { var x = a + 1; arguments[0] = 2; return x; }`,
			`function g(a) { var x = a + 1; use(arguments); return x; }`},
		{`function f(a) { var x = a + 1; outside?.(a); return x; }`,
			`function g(a) { var x = a + 1; return outside?.[x]; }`},
		{`let shared = 2; function f(a = shared) { let x = a; return x; }`,
			`function g(a = shared) { class C { field = () => a; method() { const x = a + 1; return x; } } var x = new C; return x; }`},
		{`function f(a) { var _exported = a + 1; return _exported; }`,
			`function g(a) { var shared = a + 1; function h() { var y = shared; return y; } var x = h; return x; }`},
		{`export function f(a) { var x = a + 1; return x; }`,
			`import {f} from './script0.js'; export function g(a) { let x = f(a); return x; }`},
	}
	for index, programs := range fixtures {
		for _, native := range []string{"", "1"} {
			t.Run(fmt.Sprintf("case%d/native=%s", index, native), func(t *testing.T) {
				t.Setenv("JSGO_EXPERIMENT_NATIVE_FACTS", native)
				want := newParallelFSIVTestCompiler(t, 1, programs...)
				serial := NewFlowSensitiveInlineVariables(want)
				serial.nativeFactsEnabled = nativeFunctionFactsEnabled()
				initialStamp := want.changeTracker.GetChangeStamp()
				for range 3 {
					serial.processSerial(want.GetExternsRoot(), want.GetJsRoot())
				}
				wantChanges := parallelFSIVChangeSnapshot(want)
				if index == 0 && want.changeTracker.GetChangeStamp() == initialStamp {
					t.Fatal("fixture performed no actual inlining")
				}
				for _, workers := range []int{2, 4} {
					got := newParallelFSIVTestCompiler(t, workers, programs...)
					pass := NewFlowSensitiveInlineVariables(got)
					pass.nativeFactsEnabled = nativeFunctionFactsEnabled()
					for range 3 {
						if !pass.processScriptsParallel(got.GetExternsRoot(), got.GetJsRoot()) {
							t.Fatal("fixture fell back instead of exercising parallel FSIV")
						}
					}
					if got.ToSource() != want.ToSource() || !got.GetRoot().IsEquivalentTo(want.GetRoot()) {
						t.Fatalf("workers=%d changed output:\n%s\nwant:\n%s", workers, got.ToSource(), want.ToSource())
					}
					if changes := parallelFSIVChangeSnapshot(got); !slices.Equal(changes, wantChanges) {
						t.Fatalf("workers=%d changed ordered change reports: %v; want %v", workers, changes, wantChanges)
					}
					if got.GetErrorCount() != 0 || got.GetWarningCount() != want.GetWarningCount() {
						t.Fatal("diagnostics changed")
					}
				}
			})
		}
	}
}

type parallelFSIVScopeProbe struct {
	NodeTraversal_AbstractScopedCallbackBase
	global *Scope
	found  []*Var
	bad    bool
}

func (p *parallelFSIVScopeProbe) Visit(t *NodeTraversal, n, parent *rhino.Node) {
	if !n.IsName() || n.GetString() != "shared" {
		return
	}
	var_ := t.GetScope().GetVar("shared")
	p.found = append(p.found, var_)
	p.bad = p.bad || var_ != p.global.GetVar("shared")
}

func TestParallelFSIVCanonicalGlobalScope(t *testing.T) {
	c := newParallelFSIVTestCompiler(t, 2,
		`let shared = 1; function f(a) { let x = shared + a; return () => x + shared; }`,
		`function g(a) { let x = shared + a; return () => x + shared; }`)
	jobs, global := NewFlowSensitiveInlineVariables(c).prepareParallelScripts(c.GetExternsRoot(), c.GetJsRoot())
	if global == nil || len(jobs) != 2 || global.GetRootNode() != c.GetRoot() {
		t.Fatal("did not prepare canonical common ROOT scope")
	}
	probes := make([]parallelFSIVScopeProbe, len(jobs))
	failures := parallelPrebuild(len(jobs), 2, func(i int) {
		probes[i].global = global
		NodeTraversal_Builder_().SetCompiler(&jobs[i].compiler).SetCallback(&probes[i]).TraverseWithScope(jobs[i].script, global)
	})
	for i, probe := range probes {
		if failures[i] != nil || len(probe.found) < 2 || probe.bad {
			t.Fatalf("job %d lost canonical global Var identity", i)
		}
	}
	if global.implicitVars != nil {
		t.Fatal("worker unexpectedly lazily initialized global implicit vars")
	}
}

type parallelFSIVCustomObserver struct{ changed bool }

func (o *parallelFSIVCustomObserver) ReportChange() { o.changed = true }

func TestParallelFSIVUnsupportedObserversAndRawTypesFallBack(t *testing.T) {
	for _, unsupported := range []string{"observer", "rawtype", "debug", "singleworker", "customconvention"} {
		t.Run(unsupported, func(t *testing.T) {
			c := newParallelFSIVTestCompiler(t, 2,
				`function f(a) { let x = a + 1; return x; }`, `function g(a) { let y = a + 1; return y; }`)
			switch unsupported {
			case "observer":
				c.changeTracker.AddChangeHandler(&parallelFSIVCustomObserver{})
			case "rawtype":
				c.GetJsRoot().GetFirstChild().GetFirstChild().SetJSType(c.GetTypeRegistry().GetNativeType(rhino.JSTypeNative_STRING_TYPE))
			case "debug":
				c.GetOptions().SetCheckDeterminism(true)
			case "singleworker":
				c.GetOptions().SetNumParallelThreads(1)
			case "customconvention":
				c.GetOptions().SetCodingConvention(&parallelFSIVUnknownConvention{CodingConvention: c.GetCodingConvention()})
			}
			before := c.ToSource()
			if NewFlowSensitiveInlineVariables(c).processScriptsParallel(c.GetExternsRoot(), c.GetJsRoot()) {
				t.Fatal("unsupported setup ran in parallel")
			}
			if c.ToSource() != before {
				t.Fatal("fallback mutated AST before serial path")
			}
		})
	}
}

type parallelFSIVUnknownConvention struct{ CodingConvention }

func TestParallelFSIVColorEqualityAvoidsSharedLazyHash(t *testing.T) {
	// These are the exact shared color sets consulted by AstAnalyzer. Empty
	// primitive sets make equality against STRING read-only, even for a distinct
	// Color instance with the same id (and initially uncomputed hash).
	for _, color := range []*colors.Color{colors.StandardColors_STRING, colors.StandardColors_NUMBER, colors.StandardColors_UNKNOWN} {
		if !color.GetPrototypes().IsEmpty() || !color.GetInstanceColors().IsEmpty() ||
			!color.GetOwnProperties().IsEmpty() || !color.GetUnionElements().IsEmpty() {
			t.Fatal("primitive-color equality audit requires empty sets")
		}
	}
	shared := colors.Color_SingleBuilder().SetId(colors.StandardColors_STRING.GetId()).
		SetBoxId(colors.StandardColors_STRING.GetBoxId()).BuildInternal()
	programs := []string{
		`function f(a) { var x = a.replace(/x/, 'y'); return x; }`,
		`function g(a) { var x = a.split('x'); return x; }`,
	}
	withColors := func(workers int) *Compiler {
		c := newParallelFSIVTestCompiler(t, workers, programs...)
		var paint func(*rhino.Node)
		paint = func(n *rhino.Node) {
			if n.IsName() && n.GetString() == "a" {
				n.SetColor(shared)
			}
			for child := n.GetFirstChild(); child != nil; child = child.GetNext() {
				paint(child)
			}
		}
		paint(c.GetJsRoot())
		return c
	}
	want := withColors(1)
	NewFlowSensitiveInlineVariables(want).processSerial(want.GetExternsRoot(), want.GetJsRoot())
	got := withColors(2)
	if !NewFlowSensitiveInlineVariables(got).processScriptsParallel(got.GetExternsRoot(), got.GetJsRoot()) || got.ToSource() != want.ToSource() {
		t.Fatal("color-bearing scripts did not run equivalently in parallel")
	}
}

type parallelFSIVPublicationRecorder struct {
	AbstractCompiler
	events []any
}

func (r *parallelFSIVPublicationRecorder) ReportChangeToChangeScope(n *rhino.Node) {
	r.events = append(r.events, n)
}
func (r *parallelFSIVPublicationRecorder) Report(e *JSError) { r.events = append(r.events, e) }

func TestParallelFSIVPublicationPreservesOrderAndEarliestFailure(t *testing.T) {
	script := rhino.IR_Script()
	firstDiagnostic := JSError_Make(DiagnosticType_Warning("JSC_PARALLEL_FSIV_TEST_FIRST", "first"))
	secondDiagnostic := JSError_Make(DiagnosticType_Warning("JSC_PARALLEL_FSIV_TEST_SECOND", "second"))
	early := jutil.NewIllegalStateExceptionMsg("earlier script failed")
	late := jutil.NewIllegalStateExceptionMsg("later script failed")
	jobs := []parallelFSIVJob{
		{compiler: parallelFSIVCompiler{events: []parallelFSIVEvent{{changeScope: script}, {diagnostic: firstDiagnostic}}}},
		{compiler: parallelFSIVCompiler{events: []parallelFSIVEvent{{diagnostic: secondDiagnostic}, {changeScope: script}}}},
		{compiler: parallelFSIVCompiler{events: []parallelFSIVEvent{{changeScope: script}}}},
	}
	recorder := &parallelFSIVPublicationRecorder{}
	var failure jutil.Throwable
	jutil.Try(func() {
		publishParallelFSIVJobs(recorder, jobs, []jutil.Throwable{nil, early, late})
	}, func(t jutil.Throwable) { failure = t })
	if failure != early || !slices.Equal(recorder.events, []any{script, firstDiagnostic, secondDiagnostic, script}) {
		t.Fatalf("publication lost input-order prefix or earliest failure: %v, %v", recorder.events, failure)
	}
}

func TestParallelFSIVChangeScopeResolvedBeforeDetach(t *testing.T) {
	statement := rhino.IR_ExprResult(rhino.IR_Number(1))
	inner := rhino.IR_Function(rhino.IR_Name("inner"), rhino.IR_ParamList(), rhino.IR_BlockOfNodes(statement))
	outer := rhino.IR_Function(rhino.IR_Name("outer"), rhino.IR_ParamList(), rhino.IR_BlockOfNodes(inner))
	script := rhino.IR_ScriptOfNodes(outer)
	owner := NewCompiler()
	facade := parallelFSIVCompiler{owner: owner}
	facade.ReportChangeToEnclosingScope(statement)
	facade.ReportChangeToEnclosingScope(inner) // Function itself is excluded, just as in ChangeTracker.
	facade.ReportChangeToEnclosingScope(script)
	statement.Detach()
	want := []*rhino.Node{inner, outer, script}
	for index, event := range facade.events {
		if event.changeScope != want[index] {
			t.Fatal("buffered report resolved the wrong original enclosing scope")
		}
	}
}
