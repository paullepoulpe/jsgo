package jscomp

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

func newR96NormalizedCompiler(t *testing.T, externsCode string, srcs []string, threads int, fastMode bool) *Compiler {
	t.Helper()
	c := NewCompiler()
	opts := NewCompilerOptions()
	opts.SetPrettyPrint(true)
	opts.SetNumParallelThreads(threads)
	if fastMode {
		opts.SetFastMode(true)
	}
	externs := []*SourceFile{SourceFile_FromCode("externs.js", externsCode)}
	inputs := make([]*SourceFile, len(srcs))
	for i, code := range srcs {
		inputs[i] = SourceFile_FromCode(fmt.Sprintf("in_%02d.js", i), code)
	}
	c.Init(externs, inputs, opts)
	c.ParseInputs()
	if c.HasErrors() {
		t.Fatalf("parse errors: %v", c.GetErrors())
	}
	NewNormalize(c, false).Process(c.GetExternsRoot(), c.GetJsRoot())
	c.SetLifeCycleStage(AbstractCompiler_LifeCycleStage_NORMALIZED)
	NewGatherGetterAndSetterProperties(c).Process(c.GetExternsRoot(), c.GetJsRoot())
	return c
}

func TestRound96InlineObjectLiteralsOptChainJavaParity(t *testing.T) {
	externs := "function use(x) {}\n"
	src := `
function f() {
  var term = {a: 1, b: 2};
  var invoked = {fn: function() { return 4; }};
  use(term?.a + term.b);
  use(invoked?.fn());
}
use(f());
`
	c := newR96NormalizedCompiler(t, externs, []string{src}, 1, false)
	NewInlineObjectLiterals(c, c.GetUniqueNameIdSupplier()).Process(c.GetExternsRoot(), c.GetJsRoot())
	NewAstValidator(c).ValidateRoot(c.GetRoot())
	got := c.ToSource()
	if !strings.Contains(got, "JSCompiler_object_inline_a_") {
		t.Fatalf("expected terminal optchain object literal to be inlined per Java v20261028 parity, got:\n%s", got)
	}
	if !strings.Contains(got, "invoked?.fn()") {
		t.Fatalf("expected invocation target invoked?.fn() NOT to be inlined, got:\n%s", got)
	}
}

func TestRound96InlineObjectLiteralsBloomSkipAndTransitiveHoisted(t *testing.T) {
	externs := "function use(x, y) {}\nfunction escape(x) {}\n"
	var sb strings.Builder
	sb.WriteString("(function() {\n")
	sb.WriteString("  var inlinable = {alpha: 10, beta: 20};\n")
	sb.WriteString("  var escapedViaHoisted = {gamma: 30, delta: 40};\n")
	sb.WriteString("  var orderSensitive = {k: 1};\n")
	sb.WriteString("  var dupName = {p: 99};\n")
	for i := 0; i < 60; i++ {
		fmt.Fprintf(&sb, "  function mod_%d(a, b) { var loc = a + b + %d; use(loc, %d); return loc * 2; }\n", i, i, i+1)
	}
	sb.WriteString("  var runOrder = function() { hop1(); };\n")
	sb.WriteString("  runOrder();\n")
	sb.WriteString("  orderSensitive = {k: 2};\n")
	sb.WriteString("  escapeTrigger();\n")
	sb.WriteString("  function dupName() { return 7; }\n")
	sb.WriteString("  function hop1() { hop2(); }\n")
	sb.WriteString("  function hop2() { use(orderSensitive.k, 0); }\n")
	sb.WriteString("  function escapeTrigger() { escapeInner(); }\n")
	sb.WriteString("  function escapeInner() { escape(escapedViaHoisted); }\n")
	for i := 0; i < 60; i++ {
		fmt.Fprintf(&sb, "  use(mod_%d(%d, %d), 0);\n", i, i, i+2)
	}
	sb.WriteString("  use(inlinable.alpha, inlinable.beta);\n")
	sb.WriteString("  use(escapedViaHoisted.gamma, orderSensitive.k);\n")
	sb.WriteString("  use(dupName.p, 0);\n")
	sb.WriteString("})();\n")
	// Add two top-level hoisted candidate functions to exercise filterActionableCandFns at n >= 2.
	sb.WriteString("function candHoistedA() { var oa = {u: 1, v: 2}; return oa.u + oa.v; }\n")
	sb.WriteString("function candHoistedB() { var ob = {u: 3, v: 4}; escape(ob); return ob.u; }\n")
	sb.WriteString("use(candHoistedA(), candHoistedB());\n")

	bundleSrc := sb.String()

	runMode := func(threads int, fastMode bool) (string, int) {
		c := newR96NormalizedCompiler(t, externs, []string{bundleSrc}, threads, fastMode)
		idSup := c.GetUniqueNameIdSupplier()
		iol := NewInlineObjectLiterals(c, idSup)
		iol.Process(c.GetExternsRoot(), c.GetJsRoot())
		NewAstValidator(c).ValidateRoot(c.GetRoot())
		iol.blooms.mu.RLock()
		cachedCount := len(iol.blooms.entries)
		iol.blooms.mu.RUnlock()
		return c.ToSource(), cachedCount
	}

	out1, cnt1 := runMode(1, false)
	out4, cnt4 := runMode(4, false)
	if out1 != out4 {
		t.Fatalf("1-thread vs 4-thread mismatch in default mode:\nout1:\n%s\nout4:\n%s", out1, out4)
	}
	if cnt1 < 60 || cnt4 < 60 {
		t.Fatalf("expected bloom cache to populate >= 60 inner functions, got cnt1=%d cnt4=%d", cnt1, cnt4)
	}
	if !strings.Contains(out1, "JSCompiler_object_inline_alpha_") {
		t.Fatalf("expected inlinable object literal to be split, got:\n%s", out1)
	}
	if !strings.Contains(out1, "escape(escapedViaHoisted)") {
		t.Fatalf("expected escapedViaHoisted (via transitive hoisted call) NOT to be split, got:\n%s", out1)
	}
	if !strings.Contains(out1, "function dupName()") || !strings.Contains(out1, "dupName.p") {
		t.Fatalf("expected dupName colliding with hoisted function declaration NOT to be split, got:\n%s", out1)
	}

	fast1, _ := runMode(1, true)
	fast4, _ := runMode(4, true)
	if fast1 != fast4 {
		t.Fatalf("1-thread vs 4-thread mismatch in fastMode:\nfast1:\n%s\nfast4:\n%s", fast1, fast4)
	}
}

func TestRound96InlineObjectLiteralsIncrementalSkipAndNestedInvalidation(t *testing.T) {
	externs := "function use(x) {}\nfunction escape(x) {}\n"
	var sb strings.Builder
	sb.WriteString("(function outerIIFE() {\n")
	sb.WriteString("  var keepUninlinable = {a: 1};\n")
	sb.WriteString("  escape(keepUninlinable);\n")
	for i := 0; i < 32; i++ {
		fmt.Fprintf(&sb, "  function inner_%d(x) { return use(x + %d); }\n", i, i)
		fmt.Fprintf(&sb, "  inner_%d(%d);\n", i, i)
	}
	sb.WriteString("  var nestedHolder = function() { return function() { return 123; }; };\n")
	sb.WriteString("  var unchangedHelper = function() {\n")
	sb.WriteString("    h1();\n")
	sb.WriteString("    function h1() { h2(); }\n")
	sb.WriteString("    function h2() { escape(lateCand); }\n")
	sb.WriteString("  };\n")
	sb.WriteString("  var lateCand = 0;\n")
	sb.WriteString("  var nestedCand = 0;\n")
	sb.WriteString("  use(nestedHolder()());\n")
	sb.WriteString("  unchangedHelper();\n")
	sb.WriteString("  use(lateCand);\n")
	sb.WriteString("  use(nestedCand);\n")
	sb.WriteString("})();\n")

	c := newR96NormalizedCompiler(t, externs, []string{sb.String()}, 1, false)
	idSup := c.GetUniqueNameIdSupplier()

	// Pass 1: full scan populates shared bloom cache across pass instances.
	p1 := NewInlineObjectLiterals(c, idSup)
	p1.Process(c.GetExternsRoot(), c.GetJsRoot())
	if p1.blooms == nil {
		t.Fatal("expected shared rcFnBloomCache on compiler with ChangeTracker")
	}

	// Pass 2: mutate outerIIFE directly without adding any object literal -> incremental scan skips all inner functions.
	outerFn := findNamedFunctionNode(c.GetJsRoot(), "outerIIFE")
	if outerFn.IsNil() {
		t.Fatal("expected outerIIFE function node")
	}
	c.GetChangeTracker().IncrementChangeStamp()
	c.ReportChangeToChangeScope(outerFn)

	p2 := NewInlineObjectLiterals(c, idSup)
	var s2 iolCandidateScanner
	if p2.scanCandidates(&s2, c.GetJsRoot()) {
		// keepUninlinable is still directly inside outerIIFE, so scanIncremental sees outerIIFE's own var
		// while skipping all 32 unchanged inner functions!
		p2.Process(c.GetExternsRoot(), c.GetJsRoot())
	}
	if p2.blooms != p1.blooms {
		t.Fatal("expected fresh InlineObjectLiterals instance on same Compiler to reuse shared rcFnBloomCache")
	}

	// Pass 3: mutate ONLY the deepest arrow/function inside nestedHolder to escape(nestedCand),
	// AND turn `var nestedCand = 0; use(nestedCand);` into `var nestedCand = {z: 9}; use(nestedCand.z);`
	// AND `var lateCand = 0; use(lateCand);` into `var lateCand = {w: 8}; use(lateCand.w);`.
	outerBody := outerFn.GetLastChild()
	var deepRet rhino.Node
	for stmt := outerBody.GetFirstChild(); !stmt.IsNil(); stmt = stmt.GetNext() {
		if stmt.IsVar() {
			lhs := stmt.GetFirstChild()
			switch lhs.GetString() {
			case "nestedHolder":
				outerAnon := lhs.GetFirstChild()
				innerAnon := outerAnon.GetLastChild().GetFirstChild().GetFirstChild()
				deepRet = innerAnon.GetLastChild().GetFirstChild()
			case "nestedCand":
				obj := rhino.IR_Objectlit(rhino.IR_StringKeyWithValue("z", rhino.IR_Number(9)))
				lhs.GetFirstChild().ReplaceWith(obj)
				c.ReportChangeToEnclosingScope(obj)
			case "lateCand":
				obj := rhino.IR_Objectlit(rhino.IR_StringKeyWithValue("w", rhino.IR_Number(8)))
				lhs.GetFirstChild().ReplaceWith(obj)
				c.ReportChangeToEnclosingScope(obj)
			}
		} else if stmt.IsExprResult() && stmt.GetFirstChild().IsCall() {
			call := stmt.GetFirstChild()
			arg := call.GetSecondChild()
			if !arg.IsNil() && arg.IsName() {
				switch arg.GetString() {
				case "nestedCand":
					gp := rhino.IR_Getprop(rhino.IR_Name("nestedCand"), "z")
					arg.ReplaceWith(gp)
					c.ReportChangeToEnclosingScope(gp)
				case "lateCand":
					gp := rhino.IR_Getprop(rhino.IR_Name("lateCand"), "w")
					arg.ReplaceWith(gp)
					c.ReportChangeToEnclosingScope(gp)
				}
			}
		}
	}
	if deepRet.IsNil() {
		t.Fatal("failed to locate deep return inside nestedHolder")
	}
	c.GetChangeTracker().IncrementChangeStamp()
	newRet := rhino.NewNodeWithChild(rhino.Token_RETURN, rhino.IR_Call(rhino.IR_Name("escape"), rhino.IR_Name("nestedCand")))
	deepRet.ReplaceWith(newRet)
	c.ReportChangeToEnclosingScope(newRet)

	p3 := NewInlineObjectLiterals(c, idSup)
	p3.Process(c.GetExternsRoot(), c.GetJsRoot())
	NewAstValidator(c).ValidateRoot(c.GetRoot())
	got := c.ToSource()
	if !strings.Contains(got, "escape(nestedCand)") {
		t.Fatalf("expected nestedCand NOT to be inlined after deep inner function mutation, got:\n%s", got)
	}
	if !strings.Contains(got, "escape(lateCand)") {
		t.Fatalf("expected lateCand NOT to be inlined because unchangedHelper's nested h1->h2 escapes it, got:\n%s", got)
	}
}

func TestRound96OptimizeParametersSyntacticPrepassAndLoggingParity(t *testing.T) {
	externs := "function use(x) {}\nfunction escape(x) {}\nfunction sideEffect() {}\n"
	src := `
function neverCalled(a, b) { return arguments.length + a; }
function escapedFn(a, b) { return a; }
escape(escapedFn);
escapedFn(1, 2);

function usesArgsAllowedRef(a) { return arguments[0]; }
use(typeof usesArgsAllowedRef);
usesArgsAllowedRef(10);

function usesArgsDisallowedDef(a) { return arguments[0]; }
var alias = usesArgsDisallowedDef;
use(alias);
usesArgsDisallowedDef(20);

function optimizableConst(a, b) { return a + 1; }
use(optimizableConst(42, 99));

var condFn = use(1) ? function(x) { return 1; } : function(x) { return 2; };
use(condFn(77));

class WithSideEffectField {
  f = sideEffect();
  constructor(unused) {}
}
use(new WithSideEffectField(123));

function tagUnused(strings, x) { return x; }
use(tagUnused` + "`hello ${5}`" + `);

function tagUsed(strings, x) { return strings[0] + x; }
use(tagUsed` + "`world ${6}`" + `);
`

	runOptParams := func(withLog bool, threads int) (string, string) {
		c := NewCompiler()
		opts := NewCompilerOptions()
		opts.SetPrettyPrint(true)
		opts.SetNumParallelThreads(threads)
		var logDir string
		if withLog {
			logDir = t.TempDir()
			opts.SetDebugLogDirectory(jutil.Path_Of(logDir))
			opts.SetDebugLogFilter("OptimizeParameters")
		}
		c.Init(
			[]*SourceFile{SourceFile_FromCode("externs.js", externs)},
			[]*SourceFile{SourceFile_FromCode("in_00.js", src)},
			opts,
		)
		c.ParseInputs()
		if c.HasErrors() {
			t.Fatalf("parse errors: %v", c.GetErrors())
		}
		NewNormalize(c, false).Process(c.GetExternsRoot(), c.GetJsRoot())
		c.SetLifeCycleStage(AbstractCompiler_LifeCycleStage_NORMALIZED)
		NewGatherGetterAndSetterProperties(c).Process(c.GetExternsRoot(), c.GetJsRoot())
		c.BeforePass("removeUnusedCode")
		NewRemoveUnusedCode(c, RemoveUnusedCodeOptions{
			RemoveLocalVars: true,
			RemoveGlobals:   true,
		}).Process(c.GetExternsRoot(), c.GetJsRoot())
		c.AfterPass("removeUnusedCode")
		c.BeforePass("optimizeCalls")
		NewOptimizeParameters(c).Process(c.GetExternsRoot(), c.GetJsRoot())
		c.AfterPass("optimizeCalls")
		NewAstValidator(c).ValidateRoot(c.GetRoot())
		var logText string
		if withLog {
			matches, err := filepath.Glob(filepath.Join(logDir, "OptimizeParameters", "*decisions.log"))
			if err != nil || len(matches) == 0 {
				t.Fatalf("expected OptimizeParameters decisions.log in %s (err=%v)", logDir, err)
			}
			b, err := os.ReadFile(matches[0])
			if err != nil {
				t.Fatalf("read decisions.log: %v", err)
			}
			logText = string(b)
		}
		return c.ToSource(), logText
	}

	outFast1, _ := runOptParams(false, 1)
	outFast4, _ := runOptParams(false, 4)
	outLogged, logContent := runOptParams(true, 1)
	if outFast1 != outLogged || outFast4 != outLogged {
		t.Fatalf("OptimizeParameters output mismatch between fast syntactic pre-pass and logging path:\nfast1:\n%s\nfast4:\n%s\nlogged:\n%s",
			outFast1, outFast4, outLogged)
	}
	if !strings.Contains(logContent, "decisions") && len(logContent) == 0 {
		t.Fatal("expected non-empty decisions.log when logging is enabled")
	}
	if strings.Contains(outFast1, "42, 99") {
		t.Fatalf("expected optimizableConst(42, 99) parameters to be optimized, got:\n%s", outFast1)
	}
	if strings.Contains(outFast1, "tagUnused`") {
		t.Fatalf("expected tagUnused tagged template literal to convert to normal call, got:\n%s", outFast1)
	}
	if !strings.Contains(outFast1, "tagUsed`") {
		t.Fatalf("expected tagUsed tagged template literal to remain unconverted, got:\n%s", outFast1)
	}
}

func TestRound96DevirtualizeMethodsEarlyCallSiteFilter(t *testing.T) {
	externs := "function use(x) {}\n"
	src := `
class Box {
  ok(x) { return x + 1; }
  tearOff(x) { return x + 2; }
  optCall(x) { return x + 3; }
  tagged(s) { return s; }
  constructed(x) { return x; }
}
var b = new Box();
use(b.ok(10));
use(b.tearOff);
use(b.optCall?.(20));
use(b.tagged` + "`hi`" + `);
use(new b.constructed(30));
`
	c := newR96NormalizedCompiler(t, externs, []string{src}, 4, false)
	NewOptimizeCalls(c, false, NewDevirtualizeMethods(c)).Process(c.GetExternsRoot(), c.GetJsRoot())
	NewAstValidator(c).ValidateRoot(c.GetRoot())
	got := c.ToSource()
	if !strings.Contains(got, "JSCompiler_StaticMethods_ok") {
		t.Fatalf("expected Box.prototype.ok to be devirtualized, got:\n%s", got)
	}
	for _, kept := range []string{"b.tearOff", "b.optCall?.(20)", "b.tagged`hi`", "new b.constructed(30)"} {
		if !strings.Contains(got, kept) {
			t.Fatalf("expected %q NOT to be devirtualized, got:\n%s", kept, got)
		}
	}
}

func TestRound96RemoveUnusedCodeSingleScriptParallelAndParamCache(t *testing.T) {
	externs := "function sink(x) {}\n"
	var sb strings.Builder
	sb.WriteString("var globalKept = 1;\n")
	for i := 0; i < 6; i++ {
		fmt.Fprintf(&sb, "(function(usedArg_%d, unusedTrailing_%d) {\n", i, i)
		fmt.Fprintf(&sb, "  var deadLocal_%d = %d;\n", i, i+100)
		fmt.Fprintf(&sb, "  var liveLocal_%d = usedArg_%d + globalKept;\n", i, i)
		fmt.Fprintf(&sb, "  function deadInnerFn_%d() { return %d; }\n", i, i)
		fmt.Fprintf(&sb, "  sink(liveLocal_%d);\n", i)
		fmt.Fprintf(&sb, "})(%d, %d);\n", i, i+1)
	}
	sb.WriteString("function globalFnDecl(a, b) { var deadInGlobalFn = 99; return a; }\n")
	sb.WriteString("sink(function unusedFeName() { var deadFeVar = 1; return 42; }());\n")
	sb.WriteString("sink(new (class UnusedClassExprName { m() { var deadMethodVar = 2; return 7; } })().m());\n")
	sb.WriteString("{ let deadBlockLet = 3; const liveBlockConst = 4; sink(liveBlockConst); }\n")
	sb.WriteString("for (let i = 0, deadLoopLet = 5; i < 1; i++) { let deadBodyLet = 6; sink(i); }\n")
	sb.WriteString("try { var hoistedGlobalInTry = 7; let deadTryLet = 8; sink(1); } catch (err) { let deadCatchLet = 9; sink(err); }\n")

	singleScript := sb.String()

	runLocalRUC := func(threads int) (string, int32) {
		c := newR96NormalizedCompiler(t, externs, []string{singleScript}, threads, false)
		NewRemoveUnusedCode(c, RemoveUnusedCodeOptions{
			RemoveLocalVars:                 true,
			RemoveGlobals:                   false,
			PreserveFunctionExpressionNames: false,
		}).Process(c.GetExternsRoot(), c.GetJsRoot())
		NewAstValidator(c).ValidateRoot(c.GetRoot())
		return c.ToSource(), c.GetChangeTracker().GetChangeStamp()
	}

	out1, stamp1 := runLocalRUC(1)
	out4, stamp4 := runLocalRUC(4)
	if out1 != out4 {
		t.Fatalf("1-thread vs 4-thread single-script RemoveUnusedCode mismatch:\nout1:\n%s\nout4:\n%s", out1, out4)
	}
	if stamp1 == 0 || stamp4 == 0 {
		t.Fatalf("expected non-zero change stamps after removals, got %d vs %d", stamp1, stamp4)
	}
	for _, mustGone := range []string{
		"deadLocal_0", "deadInnerFn_0", "deadInGlobalFn", "unusedFeName", "deadFeVar",
		"UnusedClassExprName", "deadMethodVar", "deadBlockLet", "deadLoopLet", "deadBodyLet", "deadTryLet", "deadCatchLet",
	} {
		if strings.Contains(out4, mustGone) {
			t.Fatalf("expected local symbol %q to be removed in single-script parallel mode, got:\n%s", mustGone, out4)
		}
	}
	for _, mustStay := range []string{"globalKept", "globalFnDecl", "hoistedGlobalInTry"} {
		if !strings.Contains(out4, mustStay) {
			t.Fatalf("expected global symbol %q to remain when RemoveGlobals=false, got:\n%s", mustStay, out4)
		}
	}

	// Also exercise isParamVarRemovable with cacheVarInfos=true (RemoveGlobals=true) across 1 vs 4 threads.
	paramSrc := `
function fnRemovableParams(usedFirst, unusedMid, unusedLast) { return usedFirst + 1; }
function fnExportedParam(usedA, _exportedConvention) { return usedA; }
function fnAllUnused(a, b = 2, ...rest) { return 99; }
sink(fnRemovableParams(1, 2, 3));
sink(fnExportedParam(4, 5));
sink(fnAllUnused(6, 7, 8));
`
	runGlobalRUC := func(threads int) string {
		c := newR96NormalizedCompiler(t, externs, []string{paramSrc}, threads, false)
		NewRemoveUnusedCode(c, RemoveUnusedCodeOptions{
			RemoveLocalVars: true,
			RemoveGlobals:   true,
		}).Process(c.GetExternsRoot(), c.GetJsRoot())
		NewAstValidator(c).ValidateRoot(c.GetRoot())
		return c.ToSource()
	}
	g1 := runGlobalRUC(1)
	g4 := runGlobalRUC(4)
	if g1 != g4 {
		t.Fatalf("global RUC param removal 1-thread vs 4-thread mismatch:\ng1:\n%s\ng4:\n%s", g1, g4)
	}
	if strings.Contains(g1, "unusedLast") || strings.Contains(g1, "rest") {
		t.Fatalf("expected trailing unused params to be removed via cached v.rucInfo path, got:\n%s", g1)
	}
}
