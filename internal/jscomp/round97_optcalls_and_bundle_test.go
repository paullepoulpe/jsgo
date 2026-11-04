/*
 * Copyright 2026 The Closure Compiler Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package jscomp

import (
	"fmt"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

func compileRound97Advanced(t *testing.T, externs string, scripts []*SourceFile, threads int) (string, []*JSError, []*JSError) {
	t.Helper()
	opts := NewCompilerOptions()
	CompilationLevel_ADVANCED_OPTIMIZATIONS.SetOptionsForCompilationLevel(opts)
	opts.SetNumParallelThreads(threads)
	c := NewCompiler()
	res := c.Compile(
		[]*SourceFile{SourceFile_FromCode("externs.js", externs)},
		scripts,
		opts,
	)
	if !res.Success {
		t.Fatalf("advanced compilation failed: %v", res.Errors)
	}
	return c.ToSource(), res.Errors, res.Warnings
}

func TestRound97_B1_GlobalRefcountsAndGoogShadowing(t *testing.T) {
	externs := CompilerTestCase_DEFAULT_EXTERNS + `
/** @param {...*} var_args */
function alert(var_args) {}
var extShadowed;
`
	s1 := `
goog.provide('my.ns.leaf');
my.ns.leaf = 123;
/** @suppress {duplicate} */
var extShadowed = function(x) { alert(x); return x; };
function takeGlobals(a, b, unused) {
  alert(a, b);
}
`
	s2 := `
takeGlobals(my.ns, extShadowed, 999);
takeGlobals(my.ns, extShadowed, 888);
alert(extShadowed(42));
`
	inputs := func() []*SourceFile {
		return []*SourceFile{
			SourceFile_FromCode("s1.js", s1),
			SourceFile_FromCode("s2.js", s2),
		}
	}

	withPropIndex(t, "0", false, false)
	want, _, _ := compileRound97Advanced(t, externs, inputs(), 1)

	withPropIndex(t, "", true, true)
	got, _, _ := compileRound97Advanced(t, externs, inputs(), 4)
	if got != want {
		t.Fatalf("B1 output mismatch:\ngot:  %s\nwant: %s", got, want)
	}

	// Verify lazy GetGlobalScope() and active B1 refcounts directly on OptimizeCalls.
	withPropIndex(t, "", true, true)
	c := NewCompiler()
	opts := NewCompilerOptions()
	c.Init(
		[]*SourceFile{SourceFile_FromCode("externs.js", externs)},
		inputs(),
		opts,
	)
	c.ParseInputs()
	c.SetLifeCycleStage(AbstractCompiler_LifeCycleStage_NORMALIZED)

	var captured *OptimizeCalls_ReferenceMap
	spy := callGraphPassFunc(func(_, _ rhino.Node, rm *OptimizeCalls_ReferenceMap) {
		captured = rm
	})
	NewOptimizeCalls(c, false, spy).Process(c.externsRoot, c.jsRoot)
	if captured == nil {
		t.Fatal("expected captured ReferenceMap")
	}
	if captured.propIndex == nil {
		t.Fatal("expected B1 global refcounts (captured.propIndex != nil) on real CompilerInputs")
	}
	if captured.globalScope != nil {
		t.Fatal("expected globalScope not to be eagerly materialized when B1 refcounts are active")
	}
	extID := uint32(strpool.Intern("extShadowed"))
	if !captured.HasGlobalVarByID(extID) {
		t.Fatal("expected HasGlobalVarByID(extShadowed) == true")
	}
	if captured.propIndex.isGlobalNonExtern(extID) {
		t.Fatal("expected isGlobalNonExtern(extShadowed) == false due to extern shadowing")
	}
	myID := uint32(strpool.Intern("my"))
	if !captured.HasGlobalVarByID(myID) {
		t.Fatal("expected HasGlobalVarByID(my) == true for implicit goog.provide root")
	}
	// Lazily materialize GetGlobalScope() on demand and check parity.
	gs := captured.GetGlobalScope()
	if gs == nil || gs.GetVar("my") == nil || gs.GetVar("extShadowed") == nil {
		t.Fatal("expected lazy GetGlobalScope() to materialize full global scope on demand")
	}

	// Synthetic script without CompilerInput falls back to CreateScopeOfScope cleanly.
	c.dropPropUseIndex()
	synthScript := rhino.IR_Script()
	synthScript.SetInputId(rhino.NewInputId("missing_synthetic.js"))
	c.jsRoot.AddChildToBack(synthScript)
	c.ReportChangeToChangeScope(synthScript)
	var capturedFallback *OptimizeCalls_ReferenceMap
	NewOptimizeCalls(c, false, callGraphPassFunc(func(_, _ rhino.Node, rm *OptimizeCalls_ReferenceMap) {
		capturedFallback = rm
	})).Process(c.externsRoot, c.jsRoot)
	if capturedFallback.propIndex != nil {
		t.Fatal("expected fallback to CreateScopeOfScope when synthetic script lacks CompilerInput")
	}
	if capturedFallback.GetGlobalScope() == nil {
		t.Fatal("expected non-nil global scope in synthetic fallback mode")
	}
}

func TestRound97_B2_CSRAndAccessorPreFilterParity(t *testing.T) {
	externs := CompilerTestCase_DEFAULT_EXTERNS + "function alert(x) {}\n"
	code := `
class Base {
  constructor(a, unusedBase) { alert(a); }
  beforeAcc(x, u1) { return x + 1; }
  get acc() { return 10; }
  set acc(v) { alert(v); }
  afterAcc(y, u2) { return y + 2; }
}
class Sub extends Base {
  constructor(a) { super(a, 999); }
}
const inst = new Sub(5);
inst.beforeAcc(1, 77);
alert(inst.acc);
inst.acc = 20;
inst?.afterAcc(2, 88);
`
	withPropIndex(t, "0", false, false)
	want, _, _ := compileRound97Advanced(t, externs, []*SourceFile{SourceFile_FromCode("in.js", code)}, 1)

	withPropIndex(t, "", true, true)
	got, _, _ := compileRound97Advanced(t, externs, []*SourceFile{SourceFile_FromCode("in.js", code)}, 2)
	if got != want {
		t.Fatalf("B2 CSR/accessor parity mismatch:\ngot:  %s\nwant: %s", got, want)
	}

	// Also verify CSR removeByID and AddReference map-materialization fallback.
	c := NewCompiler()
	c.Init(
		[]*SourceFile{SourceFile_FromCode("externs.js", externs)},
		[]*SourceFile{SourceFile_FromCode("in.js", code)},
		NewCompilerOptions(),
	)
	c.ParseInputs()
	c.SetLifeCycleStage(AbstractCompiler_LifeCycleStage_NORMALIZED)
	NewGatherGetterAndSetterProperties(c).Process(c.externsRoot, c.jsRoot)

	NewOptimizeCalls(c, false, callGraphPassFunc(func(_, _ rhino.Node, rm *OptimizeCalls_ReferenceMap) {
		if rm.props.byID != nil {
			t.Fatal("expected CSR-backed props table before mutation")
		}
		accID := strpool.Intern("acc")
		for id := range rm.GetPropReferencesByID() {
			if strpool.ID(id) == accID {
				t.Fatal("expected accessor property 'acc' to be pre-filtered from CSR table")
			}
		}
		origSize := rm.props.Size()
		beforeID := strpool.Intern("beforeAcc")
		rm.props.removeByID(beforeID)
		if rm.props.Size() != origSize-1 {
			t.Fatalf("expected Size() %d after removeByID, got %d", origSize-1, rm.props.Size())
		}
		dummy := rhino.IR_Name("extra")
		rm.AddReference(&rm.props, "extraProp", dummy)
		if rm.props.byID == nil {
			t.Fatal("expected AddReference to materialize byID map cleanly")
		}
		if rm.props.Size() != origSize {
			t.Fatalf("expected Size() %d after AddReference, got %d", origSize, rm.props.Size())
		}
	})).Process(c.externsRoot, c.jsRoot)
}

func TestRound97_B3_ArgumentsArrowsHasValueReturnAndMayOptMemo(t *testing.T) {
	externs := CompilerTestCase_DEFAULT_EXTERNS + "function alert(x) {}\n"
	code := `
function directArgs(a) { alert(arguments[0]); return arguments[0]; }
function arrowArgs(a) { const f = () => arguments[0]; alert(f()); return f(); }
function deepArrowArgs(a) { const g = () => () => arguments.length; alert(g()()); return g()(); }
function nestedRegularFnArgs(a, dropMe) {
  /** @param {...*} var_args */
  function inner(var_args) { return arguments[0]; }
  alert(a);
  return a + inner(9);
}
function noValueReturn(x) {
  if (x > 0) { alert(x); return; }
  function nestedVal() { return 42; }
  alert(nestedVal());
}
function withValueReturn(x) {
  return alert(x);
}
alert(directArgs(10));
alert(directArgs(11));
alert(arrowArgs(20));
alert(arrowArgs(21));
alert(deepArrowArgs(30));
alert(deepArrowArgs(31));
alert(nestedRegularFnArgs(40, 999));
alert(nestedRegularFnArgs(41, 888));
noValueReturn(1);
noValueReturn(2);
withValueReturn(3);
withValueReturn(4);
`
	withPropIndex(t, "0", false, false)
	want, _, _ := compileRound97Advanced(t, externs, []*SourceFile{SourceFile_FromCode("in.js", code)}, 1)

	withPropIndex(t, "", true, true)
	got, _, _ := compileRound97Advanced(t, externs, []*SourceFile{SourceFile_FromCode("in.js", code)}, 2)
	if got != want {
		t.Fatalf("B3 output mismatch:\ngot:  %s\nwant: %s", got, want)
	}
	if strings.Contains(got, "999") {
		t.Fatalf("expected unused arg 999 of nestedRegularFnArgs to be eliminated, got: %s", got)
	}
	for _, kept := range []string{"10", "20", "30"} {
		if !strings.Contains(got, kept) {
			t.Fatalf("expected arg %s passed to arguments-referencing function to be preserved, got: %s", kept, got)
		}
	}

	// Verify function named `arguments` does not self-flag its own name node,
	// and verify MayBeOptimizableNameByID invalidates when CodingConvention changes.
	c := NewCompiler()
	c.Init(
		[]*SourceFile{SourceFile_FromCode("externs.js", externs)},
		[]*SourceFile{SourceFile_FromCode("in.js", "function f() {} f();")},
		NewCompilerOptions(),
	)
	c.ParseInputs()
	c.SetLifeCycleStage(AbstractCompiler_LifeCycleStage_NORMALIZED)

	fnNamedArguments := rhino.IR_Function(
		rhino.IR_Name("arguments"),
		rhino.IR_ParamList(rhino.IR_Name("p")),
		rhino.IR_Block(),
	)
	c.jsRoot.GetFirstChild().AddChildToBack(fnNamedArguments)
	c.ReportChangeToChangeScope(c.jsRoot.GetFirstChild())
	c.changeTracker.MarkNewScopesChanged(fnNamedArguments)

	ix := propUseIndexFor(c, c.externsRoot, c.jsRoot)
	if ix == nil {
		t.Fatal("expected non-nil propUseIndex")
	}
	if ix.doesFunctionRefOwnArguments(fnNamedArguments) {
		t.Fatal("function named 'arguments' must not self-flag piScopeFlagHasArgumentsName")
	}

	privID := uint32(strpool.Intern("_exportedUnderscore"))
	c.defaultCodingConvention = NewClosureCodingConvention()
	first := OptimizeCalls_MayBeOptimizableNameByID(c, privID)
	c.defaultCodingConvention = NewGoogleCodingConvention()
	second := OptimizeCalls_MayBeOptimizableNameByID(c, privID)
	wantGoogle := !NewGoogleCodingConvention().IsExported("_exportedUnderscore")
	if second != wantGoogle || first == second {
		t.Fatalf("MayBeOptimizableNameByID convention invalidation failed: closure=%v google=%v wantGoogle=%v", first, second, wantGoogle)
	}
}

func TestRound97_B4_IncrementalCleanSkippingAndCallSiteChangeTracking(t *testing.T) {
	withPropIndex(t, "", false, true)
	externs := CompilerTestCase_DEFAULT_EXTERNS + "function alert(x) {} var sideEffectGlobal;\n"
	s1 := `
function keepBoth(a, b) { alert(a + b); return a + b; }
function cleanVoid(a) { alert(a); }
function targetFn(a, b) { alert(a); return b; }
`
	s2 := `
alert(keepBoth(1, 2));
keepBoth(3, 4, 777);
keepBoth(5, void 0);
var usedRes = targetFn(1, 2);
alert(usedRes);
targetFn(1, 9);
`
	s3 := `
cleanVoid(10);
cleanVoid(20);
`
	c := NewCompiler()
	opts := NewCompilerOptions()
	c.Init(
		[]*SourceFile{SourceFile_FromCode("externs.js", externs)},
		[]*SourceFile{
			SourceFile_FromCode("s1.js", s1),
			SourceFile_FromCode("s2.js", s2),
			SourceFile_FromCode("s3.js", s3),
		},
		opts,
	)
	c.ParseInputs()
	c.SetLifeCycleStage(AbstractCompiler_LifeCycleStage_NORMALIZED)
	c.phaseOptimizer = &PhaseOptimizer{inLoop: true}
	defer c.dropPropUseIndex()

	oc := NewOptimizeCalls(c, false, NewOptimizeReturns(c), NewOptimizeParameters(c))
	// Run 1: extra 777 and trailing void 0 on keepBoth should be removed; cleanVoid (whose calls
	// are in untouched s3) marked clean for both passes; targetFn inlines a=1, keeps return b.
	oc.Process(c.externsRoot, c.jsRoot)
	src1 := c.ToSource()
	if strings.Contains(src1, "777") || strings.Contains(src1, "void 0") {
		t.Fatalf("expected extra arg 777 and trailing void 0 to be removed in Run 1, got:\n%s", src1)
	}

	ix := c.propIndex
	if ix == nil {
		t.Fatal("expected active propUseIndex")
	}
	cleanVoidID := uint32(strpool.Intern("cleanVoid"))
	if ix.oc.namesClean.entries[cleanVoidID].flags&(ocCleanReturns|ocCleanParams) != (ocCleanReturns | ocCleanParams) {
		t.Fatalf("expected cleanVoid to be cached clean for both passes after Run 1, got flags=%08b",
			ix.oc.namesClean.entries[cleanVoidID].flags)
	}

	// Run 2 with zero AST mutations: cleanVoid is skipped cleanly via CSR + clean cache, and keepBoth becomes clean too.
	oc.Process(c.externsRoot, c.jsRoot)
	if c.propIndex != ix {
		t.Fatal("expected propUseIndex to be updated incrementally in-loop, not rebuilt")
	}
	if c.ToSource() != src1 {
		t.Fatal("expected idempotent Run 2 with zero changes")
	}
	keepBothID := uint32(strpool.Intern("keepBoth"))
	if ix.oc.namesClean.entries[keepBothID].flags&(ocCleanReturns|ocCleanParams) != (ocCleanReturns | ocCleanParams) {
		t.Fatalf("expected keepBoth to be cached clean for both passes after Run 2, got flags=%08b",
			ix.oc.namesClean.entries[keepBothID].flags)
	}

	// Now mutate s2 so targetFn(b)'s return is no longer used AND both calls pass b=2,
	// and report change to s2's SCRIPT scope. Run 3 must invalidate targetFn and optimize both return and param!
	script2 := c.jsRoot.GetSecondChild()
	NodeUtil_MarkFunctionsDeleted(script2, c)
	replacement := c.ParseSyntheticCode("s2_mut.js", "targetFn(2); targetFn(2); alert(keepBoth(1, 2));")
	script2.RemoveChildren()
	for ch := replacement.RemoveFirstChild(); !ch.IsNil(); ch = replacement.RemoveFirstChild() {
		script2.AddChildToBack(ch)
	}
	c.ReportChangeToChangeScope(script2)

	oc.Process(c.externsRoot, c.jsRoot)
	src3 := c.ToSource()
	if strings.Contains(src3, "return b") || strings.Contains(src3, "targetFn(2)") {
		t.Fatalf("expected targetFn return value and constant param b=2 to be optimized after invalidation, got:\n%s", src3)
	}

	// Verify optimizeCallSite reports change to enclosing scope when callArgNode.IsNil()
	// (e.g. a 0-arg call `fn()` whose side-effect flags are updated when moving a mutable default/fixed arg).
	c2 := NewCompiler()
	c2.Init(
		[]*SourceFile{SourceFile_FromCode("externs.js", externs)},
		[]*SourceFile{SourceFile_FromCode("in.js", "function fn(x) { return x; } fn();")},
		NewCompilerOptions(),
	)
	c2.ParseInputs()
	c2.SetLifeCycleStage(AbstractCompiler_LifeCycleStage_NORMALIZED)
	script := c2.jsRoot.GetFirstChild()
	exprRes := script.GetSecondChild()
	callNode := exprRes.GetFirstChild()
	// Pure call flag initially (no global state mutation), with MayMutateArguments set:
	callNode.SetSideEffectFlagsObject(rhino.NewNode_SideEffectFlags().ClearAllFlags().SetMutatesArguments())
	op := NewOptimizeParameters(c2)
	// Drain timeline marker before calling optimizeCallSite directly with a non-immutable arg ([]) at index 0
	// while callNode itself has 0 actual argument nodes (callArgNode.IsNil()).
	_ = c2.GetChangeTracker().GetChangedScopeNodesForPass("r97Check")
	op.optimizeCallSite([]optimizeParametersParam{{
		arg:          rhino.IR_Arraylit(),
		shouldRemove: true,
	}}, callNode.GetFirstChild())
	changed := c2.GetChangeTracker().GetChangedScopeNodesForPass("r97Check")
	if len(changed) != 1 || changed[0] != script {
		t.Fatalf("expected optimizeCallSite to report change on enclosing script when callArgNode.IsNil(), got: %v", changed)
	}
}

type r97RawDiagCollector struct {
	items []string
}

func (r *r97RawDiagCollector) Report(level CheckLevel, err *JSError) {
	r.items = append(r.items, fmt.Sprintf("%s:%s", level, err.String()))
}

func buildNestedPeepholeBundleJS(fnCount int) string {
	var b strings.Builder
	b.WriteString("var sink = [1 << 0.5];\n")
	b.WriteString("(function(ns) {\n")
	b.WriteString("  var iifeDiag = 2 >> 0.5;\n")
	for i := 0; i < fnCount; i++ {
		fmt.Fprintf(&b, `
  class C%[1]d extends (function() { return Object; })() {
    ["m" + (%[1]d + 1)](x) {
      var a = (%[1]d + 100) + (%[1]d << 0.25);
      if (false) {
        var deadFn = function() { return (999 + %[1]d) + (3 >>> 0.5); };
      }
      for (var k in x) {
        var inner = (y) => {
          var z = (() => (y + %[1]d) + (%[2]d.5 << 1))();
          return z + (%[1]d >> 0.75);
        };
        sink.push(inner(a), k, iifeDiag);
      }
      return a;
    }
  }
  ns["C%[1]d"] = C%[1]d;
`, i, i+1)
	}
	b.WriteString("})(sink);\n")
	return b.String()
}

func runSinglePeepholePassForTest(t *testing.T, scripts []*SourceFile, threads int) (string, []string, []parser.FeatureSet) {
	t.Helper()
	var raw r97RawDiagCollector
	opts := NewCompilerOptions()
	opts.SetNumParallelThreads(threads)
	opts.SetErrorHandler(&raw)
	c := NewCompiler()
	c.Init(
		[]*SourceFile{SourceFile_FromCode("externs.js", CompilerTestCase_DEFAULT_EXTERNS)},
		scripts,
		opts,
	)
	c.ParseInputs()
	pass := NewPeepholeOptimizationsPass(
		c,
		"peepholeR97",
		NewPeepholeRemoveDeadCode(),
		NewPeepholeMinimizeConditions(true),
		NewPeepholeSubstituteAlternateSyntax(true),
		NewPeepholeReplaceKnownMethods(true, false),
		NewPeepholeFoldConstants(true, false),
	)
	pass.Process(c.externsRoot, c.jsRoot)

	var feats []parser.FeatureSet
	for sc := c.jsRoot.GetFirstChild(); !sc.IsNil(); sc = sc.GetNext() {
		feats = append(feats, NodeUtil_GetFeatureSetOfScript(sc))
	}
	return c.ToSource(), raw.items, feats
}

func TestRound97_PeepholeWaveParallelParityAndRace(t *testing.T) {
	bundle := buildNestedPeepholeBundleJS(40)
	makeScripts := func() []*SourceFile {
		return []*SourceFile{
			SourceFile_FromCode("bundle.js", bundle),
			SourceFile_FromCode("tail.js", "var tailVal = (1 + 2) + (7 << 0.5);"),
		}
	}

	wantSrc, wantDiags, wantFeats := runSinglePeepholePassForTest(t, makeScripts(), 1)
	gotSrc, gotDiags, gotFeats := runSinglePeepholePassForTest(t, makeScripts(), 4)

	if gotSrc != wantSrc {
		t.Fatalf("wave-parallel Peephole source mismatch:\ngot:\n%s\nwant:\n%s", gotSrc, wantSrc)
	}
	if len(wantDiags) < 100 {
		t.Fatalf("expected >=100 raw Peephole diagnostics across depths, got %d", len(wantDiags))
	}
	if len(gotDiags) != len(wantDiags) {
		t.Fatalf("diagnostic count mismatch: got %d, want %d", len(gotDiags), len(wantDiags))
	}
	for i := range wantDiags {
		if gotDiags[i] != wantDiags[i] {
			t.Fatalf("raw arrival-order diagnostic[%d] mismatch:\ngot:  %s\nwant: %s", i, gotDiags[i], wantDiags[i])
		}
	}
	for i := range wantFeats {
		if gotFeats[i] != wantFeats[i] {
			t.Fatalf("script[%d] FeatureSet mismatch: got %v, want %v", i, gotFeats[i], wantFeats[i])
		}
	}

	// Test 1 for processChangedScopesParallel: 3-wave live ancestor+descendant FUNCTION scopes
	// interleaved in Timeline order, verifying exact 0..n-1 raw diagnostic order, FeatureSet, and AST parity.
	buildLiveIncremental := func(threads int) (string, []string, parser.FeatureSet) {
		var sb strings.Builder
		sb.WriteString("var out = [];\n")
		for i := 0; i < 12; i++ {
			fmt.Fprintf(&sb, `
function outer%[1]d(x) {
  var c = 10 + %[1]d + (%[1]d << 0.5);
  function mid%[1]d(y) {
    var m = (20 + %[1]d) + (%[1]d >>> 0.25);
    function inner%[1]d(z) {
      for (var k in z) {
        out.push(k);
      }
      return (30 + %[1]d) + (%[1]d >> 0.75);
    }
    return y + m + inner%[1]d(x);
  }
  return mid%[1]d(c);
}
`, i)
		}
		var raw r97RawDiagCollector
		opts := NewCompilerOptions()
		opts.SetNumParallelThreads(threads)
		opts.SetErrorHandler(&raw)
		c := NewCompiler()
		c.Init(
			[]*SourceFile{SourceFile_FromCode("externs.js", CompilerTestCase_DEFAULT_EXTERNS)},
			[]*SourceFile{SourceFile_FromCode("inc.js", sb.String())},
			opts,
		)
		c.ParseInputs()
		pass := NewPeepholeOptimizationsPass(
			c,
			"peepholeIncLiveR97",
			NewPeepholeRemoveDeadCode(),
			NewPeepholeFoldConstants(true, false),
		)
		// Mark pass once so next Process takes the incremental changedScopeNodes path!
		_ = c.GetChangeTracker().GetChangedScopeNodesForPass("peepholeIncLiveR97")
		// Report scopes in interleaved order: [outer_i, inner_i, mid_i] so Timeline order != wave order!
		for fn := c.jsRoot.GetFirstChild().GetFirstChild(); !fn.IsNil(); fn = fn.GetNext() {
			if !fn.IsFunction() {
				continue
			}
			var mids, inners []rhino.Node
			var walk func(n rhino.Node, depth int)
			walk = func(n rhino.Node, depth int) {
				for ch := n.GetFirstChild(); !ch.IsNil(); ch = ch.GetNext() {
					if ch.IsFunction() {
						if depth == 0 {
							mids = append(mids, ch)
						} else {
							inners = append(inners, ch)
						}
						walk(ch, depth+1)
					} else {
						walk(ch, depth)
					}
				}
			}
			walk(fn, 0)
			c.ReportChangeToChangeScope(fn)
			for _, in := range inners {
				c.ReportChangeToChangeScope(in)
			}
			for _, mid := range mids {
				c.ReportChangeToChangeScope(mid)
			}
		}
		pass.Process(c.externsRoot, c.jsRoot)
		return c.ToSource(), raw.items, NodeUtil_GetFeatureSetOfScript(c.jsRoot.GetFirstChild())
	}

	incWantSrc, incWantDiags, incWantFeat := buildLiveIncremental(1)
	incGotSrc, incGotDiags, incGotFeat := buildLiveIncremental(4)
	if incGotSrc != incWantSrc {
		t.Fatalf("incremental wave Peephole source mismatch:\ngot:\n%s\nwant:\n%s", incGotSrc, incWantSrc)
	}
	if len(incWantDiags) < 36 {
		t.Fatalf("expected >=36 diagnostics across 3 waves, got %d", len(incWantDiags))
	}
	if len(incGotDiags) != len(incWantDiags) {
		t.Fatalf("incremental diagnostic count mismatch: got %d, want %d", len(incGotDiags), len(incWantDiags))
	}
	for i := range incWantDiags {
		if incGotDiags[i] != incWantDiags[i] {
			t.Fatalf("incremental diagnostic[%d] Timeline order mismatch:\ngot:  %s\nwant: %s", i, incGotDiags[i], incWantDiags[i])
		}
	}
	if incGotFeat != incWantFeat {
		t.Fatalf("incremental script FeatureSet mismatch: got %v, want %v", incGotFeat, incWantFeat)
	}

	// Test 2 for processChangedScopesParallel: ancestor scope in wave 0 detaches/deletes a dead
	// descendant FUNCTION scope that was also in changedScopeNodes; wave 1 must skip the detached
	// descendant cleanly without panicking or mutating detached subtrees.
	var sbDetached strings.Builder
	sbDetached.WriteString("var out = [];\n")
	for i := 0; i < 12; i++ {
		fmt.Fprintf(&sbDetached, `
function outer%[1]d(x) {
  var c = 10 + %[1]d;
  if (false) {
    function deadChild%[1]d(obj) { for (var k in obj) { out.push(k); } return 1 + 2; }
    out.push(deadChild%[1]d(x));
  }
  function liveChild%[1]d(y) { return y + (20 + %[1]d); }
  return liveChild%[1]d(c);
}
`, i)
	}
	optsDet := NewCompilerOptions()
	optsDet.SetNumParallelThreads(4)
	cDet := NewCompiler()
	cDet.Init(
		[]*SourceFile{SourceFile_FromCode("externs.js", CompilerTestCase_DEFAULT_EXTERNS)},
		[]*SourceFile{SourceFile_FromCode("det.js", sbDetached.String())},
		optsDet,
	)
	cDet.ParseInputs()
	passDet := NewPeepholeOptimizationsPass(
		cDet,
		"peepholeIncDetR97",
		NewPeepholeRemoveDeadCode(),
		NewPeepholeFoldConstants(true, false),
	)
	_ = cDet.GetChangeTracker().GetChangedScopeNodesForPass("peepholeIncDetR97")
	for fn := cDet.jsRoot.GetFirstChild().GetFirstChild(); !fn.IsNil(); fn = fn.GetNext() {
		if fn.IsFunction() {
			cDet.ReportChangeToChangeScope(fn)
			var reportKids func(n rhino.Node)
			reportKids = func(n rhino.Node) {
				for ch := n.GetFirstChild(); !ch.IsNil(); ch = ch.GetNext() {
					if ch.IsFunction() {
						cDet.ReportChangeToChangeScope(ch)
					}
					reportKids(ch)
				}
			}
			reportKids(fn)
		}
	}
	passDet.Process(cDet.externsRoot, cDet.jsRoot)
	detSrc := cDet.ToSource()
	if strings.Contains(detSrc, "deadChild") || strings.Contains(detSrc, "20 + ") {
		t.Fatalf("expected deadChild to be removed in wave 0 and liveChild folded in wave 1, got:\n%s", detSrc)
	}
}

func TestRound97_VariableReferenceCheck2ScriptBundleParallel(t *testing.T) {
	runVRC := func(threads int) []string {
		opts := NewCompilerOptions()
		opts.SetNumParallelThreads(threads)
		opts.SetWarningLevel(DiagnosticGroups_CHECK_VARIABLES, CheckLevel_WARNING)
		c := NewCompiler()
		res := c.Compile(
			[]*SourceFile{SourceFile_FromCode("externs.js", CompilerTestCase_DEFAULT_EXTERNS)},
			[]*SourceFile{
				SourceFile_FromCode("s1.js", "function f1() { if (true) { var x = 1; var x = 2; } return x; } f1();"),
				SourceFile_FromCode("s2.js", "function f2() { if (true) { var y = 3; var y = 4; } return y; } f2();"),
				SourceFile_FromCode("s3.js", "function f3() { var z = 5; var z = 6; return z; } f3();"),
			},
			opts,
		)
		var out []string
		for _, e := range res.Errors {
			out = append(out, e.String())
		}
		for _, w := range res.Warnings {
			out = append(out, w.String())
		}
		return out
	}

	want := runVRC(1)
	got := runVRC(4)
	if len(want) == 0 {
		t.Fatal("expected VariableReferenceCheck warnings on duplicate var declarations")
	}
	if len(got) != len(want) {
		t.Fatalf("VRC diagnostic count mismatch: got %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("VRC diagnostic[%d] mismatch:\ngot:  %s\nwant: %s", i, got[i], want[i])
		}
	}
}

func TestRound97_InlineFunctionsPhase4PruningParity(t *testing.T) {
	externs := CompilerTestCase_DEFAULT_EXTERNS + "function alert(x) {}\n"
	var sb strings.Builder
	// 50 untouched large functions with no inlineable calls (pruned in Phase 4).
	for i := 0; i < 50; i++ {
		fmt.Fprintf(&sb, `
function heavyUntouched%[1]d(a, b) {
  var r = a * %[1]d + b;
  if (r > 10) { alert(r); } else { alert(r - 1); }
  return r + a - b;
}
alert(heavyUntouched%[1]d(%[1]d, %[1]d + 1));
`, i)
	}
	// Active inline candidates: direct, block, decomposed expression, nested, computed class key, and unreferenced removable.
	sb.WriteString(`
function inlineDirect(x) { return x + 42; }
function inlineInner(y) { return y * 3; }
function inlineBlock(z) { var t = z + 1; alert(t); return t * 2; }
function unreferencedRemovable() { return 77777; }

class Widget {
  [inlineDirect(1)]() {
    var combined = alert(100) || (inlineBlock(inlineInner(5)) + inlineDirect(2));
    return combined;
  }
}
alert(new Widget()[43]());
`)
	code := sb.String()

	withPropIndex(t, "0", false, false)
	want, _, _ := compileRound97Advanced(t, externs, []*SourceFile{
		SourceFile_FromCode("part1.js", "var seed = 1;\n"),
		SourceFile_FromCode("part2.js", code),
	}, 1)

	withPropIndex(t, "", true, true)
	got, _, _ := compileRound97Advanced(t, externs, []*SourceFile{
		SourceFile_FromCode("part1.js", "var seed = 1;\n"),
		SourceFile_FromCode("part2.js", code),
	}, 4)

	if got != want {
		t.Fatalf("InlineFunctions + full pipeline parity mismatch:\ngot:\n%s\nwant:\n%s", got, want)
	}
	if strings.Contains(got, "77777") {
		t.Fatalf("expected unreferencedRemovable to be removed, got:\n%s", got)
	}
}
