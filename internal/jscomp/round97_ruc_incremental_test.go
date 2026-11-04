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

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

func setupNormalizedTestCompiler(t *testing.T, externsCode string, scripts []string, threads int) *Compiler {
	t.Helper()
	compiler := NewCompiler()
	opts := NewCompilerOptions()
	opts.SetLanguageIn(CompilerOptions_LanguageMode_ECMASCRIPT_2020)
	opts.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT_2020)
	if threads > 1 {
		opts.SetNumParallelThreads(threads)
	}
	inputs := make([]*SourceFile, len(scripts))
	for i, code := range scripts {
		inputs[i] = SourceFile_FromCode(fmt.Sprintf("in_%d.js", i), code)
	}
	compiler.Init(
		[]*SourceFile{SourceFile_FromCode("externs.js", externsCode)},
		inputs,
		opts,
	)
	compiler.ParseInputs()
	if len(compiler.GetErrors()) > 0 {
		t.Fatalf("parse errors: %v", compiler.GetErrors())
	}
	NewNormalize(compiler, false).Process(compiler.GetExternsRoot(), compiler.GetJsRoot())
	NewGatherExternProperties(compiler, GatherExternProperties_Mode_OPTIMIZE).Process(compiler.GetExternsRoot(), compiler.GetJsRoot())
	compiler.SetAccessorSummary(&AccessorSummary{})
	return compiler
}

func advancedRucOpts() RemoveUnusedCodeOptions {
	return RemoveUnusedCodeOptions{
		RemoveLocalVars:                               true,
		RemoveGlobals:                                 true,
		RemoveUnusedPrototypeProperties:               true,
		RemoveUnusedThisProperties:                    true,
		RemoveUnusedObjectDefinePropertiesDefinitions: true,
		RemoveUnusedPolyfills:                         true,
	}
}

func deleteGlobalCallTarget(t *testing.T, compiler *Compiler, fnName string) {
	t.Helper()
	var walk func(n rhino.Node) bool
	walk = func(n rhino.Node) bool {
		if n.IsCall() && !n.GetFirstChild().IsNil() && n.GetFirstChild().IsName() && n.GetFirstChild().GetString() == fnName {
			stmt := n.GetParent()
			if stmt.IsExprResult() {
				NodeUtil_DeleteNode(stmt, compiler)
				return true
			}
			compiler.ReportChangeToEnclosingScope(n)
			n.ReplaceWith(rhino.IR_Number(0).Srcref(n))
			return true
		}
		for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
			if walk(c) {
				return true
			}
		}
		return false
	}
	if !walk(compiler.GetJsRoot()) {
		t.Fatalf("call to %q not found in AST", fnName)
	}
}

func renderClean(compiler *Compiler) string {
	return compiler.ToSource()
}

func TestRound97_AdvancedConstructorIndepPropBecomesDeadAndReplays(t *testing.T) {
	externs := "var window = {}; function use(x) {}\n"

	t.Run("IndepPropBecomesDeadDeletesPropAndGlobalHelper", func(t *testing.T) {
		src := `
function GlobalHelper() { return 42; }
/** @constructor */
function Widget() {
  this.prop_ = () => GlobalHelper();
}
function keepReader(w) { return use(w.prop_()); }
function keepAlive() { return use(1); }
window['w'] = new Widget();
keepReader(window['w']);
keepAlive();
`
		c := setupNormalizedTestCompiler(t, externs, []string{src}, 1)
		opts := advancedRucOpts()

		// Iter 1: Widget has zero AST changes and is cached with indepPropIDs=[prop_], varIDs=[GlobalHelper].
		NewRemoveUnusedCode(c, opts).Process(c.GetExternsRoot(), c.GetJsRoot())
		if c.rucGuarded == nil || c.rucGuarded.fnCache == nil {
			t.Fatalf("expected fnCache to be populated after iter 1")
		}

		// Delete the call keepReader(window['w']) in global script without touching Widget or GlobalHelper.
		deleteGlobalCallTarget(t, c, "keepReader")

		// Iter 2: Widget is deferred because prop_ is unpinned, then traversed in full.
		NewRemoveUnusedCode(c, opts).Process(c.GetExternsRoot(), c.GetJsRoot())
		if c.rucGuarded == nil || c.rucGuarded.fnCache == nil {
			t.Fatalf("expected fnCache after iter 2")
		}
		if c.rucGuarded.fnCache.lastDeferred == 0 {
			t.Fatalf("expected Widget summary to be deferred on iter 2")
		}
		if c.rucGuarded.fnCache.lastReplayed == 0 {
			t.Fatalf("expected keepAlive summary to replay on iter 2")
		}

		// Compare against a single fresh non-incremental run on the same post-edit AST.
		cRef := setupNormalizedTestCompiler(t, externs, []string{src}, 1)
		deleteGlobalCallTarget(t, cRef, "keepReader")
		cRef.rucGuarded = nil
		prevEnv := rucIncrEnv
		rucIncrEnv = "0"
		NewRemoveUnusedCode(cRef, opts).Process(cRef.GetExternsRoot(), cRef.GetJsRoot())
		rucIncrEnv = prevEnv

		got := renderClean(c)
		want := renderClean(cRef)
		if got != want {
			t.Fatalf("AST mismatch after deferred indepProp removal:\ngot:\n%s\nwant:\n%s", got, want)
		}
		if strings.Contains(got, "GlobalHelper") || strings.Contains(got, "prop_") {
			t.Fatalf("expected GlobalHelper and prop_ to be removed, got:\n%s", got)
		}
	})

	t.Run("IndepPropDeferredThenPinnedReplaysSummary", func(t *testing.T) {
		// Widget appears BEFORE lateReader in worklist order, so Widget is deferred first,
		// then lateReader pins prop_, and Widget replays its summary in O(1) without full traversal!
		src := `
function GlobalHelper() { return 99; }
/** @constructor */
function Widget() {
  this.prop_ = GlobalHelper();
}
function lateReader(w) { return use(w.prop_); }
function dropMe() { return use(7); }
window['make'] = function() {
  var w = new Widget();
  lateReader(w);
  return w;
};
dropMe();
`
		c := setupNormalizedTestCompiler(t, externs, []string{src}, 1)
		opts := advancedRucOpts()

		NewRemoveUnusedCode(c, opts).Process(c.GetExternsRoot(), c.GetJsRoot())
		deleteGlobalCallTarget(t, c, "dropMe")
		NewRemoveUnusedCode(c, opts).Process(c.GetExternsRoot(), c.GetJsRoot())

		fc := c.rucGuarded.fnCache
		if fc.lastDeferred == 0 || fc.lastReplayed < 2 {
			t.Fatalf("expected deferred+replayed summary on iter 2, got deferred=%d replayed=%d", fc.lastDeferred, fc.lastReplayed)
		}
		got := renderClean(c)
		if !strings.Contains(got, "GlobalHelper") || !strings.Contains(got, "prop_") || strings.Contains(got, "dropMe") {
			t.Fatalf("unexpected output:\n%s", got)
		}
	})
}

func TestRound97_AdvancedFIFONestedContinuationsAndSharedGlobalPin(t *testing.T) {
	externs := "var window = {}; function use(x) {}\n"
	src := `
function SharedGlobal() { return 1; }
function ArrowGlobal(x) { return x + 2; }
function DeclGlobal(y) { return y + 3; }
function LocalContGlobal() { return 4; }
function LocalClassGlobal() { return 5; }

function firstCaller() {
  return use(SharedGlobal());
}

function secondCaller(arr) {
  function innerDecl(v) { return DeclGlobal(v); }
  var loc = LocalContGlobal();
  class LocalCls {
    run() { return LocalClassGlobal(); }
  }
  var mapped = arr.map((item) => ArrowGlobal(item) + SharedGlobal());
  return use(innerDecl(loc) + new LocalCls().run() + mapped.length);
}

firstCaller();
window['entry'] = secondCaller;
`
	c := setupNormalizedTestCompiler(t, externs, []string{src}, 1)
	opts := advancedRucOpts()

	// Iter 1: firstCaller pins SharedGlobal BEFORE secondCaller runs; secondCaller must still record SharedGlobal
	// and all FIFO nested continuations (arrow, inner fn decl, local var cont, local class method).
	NewRemoveUnusedCode(c, opts).Process(c.GetExternsRoot(), c.GetJsRoot())

	// Delete firstCaller() call so secondCaller's summary replay is now the SOLE pin of SharedGlobal!
	deleteGlobalCallTarget(t, c, "firstCaller")
	NewRemoveUnusedCode(c, opts).Process(c.GetExternsRoot(), c.GetJsRoot())

	if c.rucGuarded.fnCache.lastReplayed == 0 {
		t.Fatalf("expected secondCaller to replay from summary")
	}

	cRef := setupNormalizedTestCompiler(t, externs, []string{src}, 1)
	deleteGlobalCallTarget(t, cRef, "firstCaller")
	prevEnv := rucIncrEnv
	rucIncrEnv = "0"
	NewRemoveUnusedCode(cRef, opts).Process(cRef.GetExternsRoot(), cRef.GetJsRoot())
	rucIncrEnv = prevEnv

	got := renderClean(c)
	want := renderClean(cRef)
	if got != want {
		t.Fatalf("FIFO/shared pin replay mismatch:\ngot:\n%s\nwant:\n%s", got, want)
	}
	for _, sym := range []string{"SharedGlobal", "ArrowGlobal", "DeclGlobal", "LocalContGlobal", "LocalClassGlobal"} {
		if !strings.Contains(got, sym) {
			t.Fatalf("expected %s to remain alive, got:\n%s", sym, got)
		}
	}
	if strings.Contains(got, "firstCaller") {
		t.Fatalf("expected firstCaller to be removed, got:\n%s", got)
	}
}

func TestRound97_NamedFunctionExpressionResultUsedFlip(t *testing.T) {
	externs := "function use(x) {}\n"
	src := `use(function namedFn() { namedFn.tag = 1; });`

	c := setupNormalizedTestCompiler(t, externs, []string{src}, 1)
	opts := advancedRucOpts()
	opts.PreserveFunctionExpressionNames = false

	// Iter 1: expression result of namedFn IS used (argument to use(...)), so namedFn escapes and stays pinned.
	NewRemoveUnusedCode(c, opts).Process(c.GetExternsRoot(), c.GetJsRoot())

	// Rewrite `use(function namedFn() { ... });` -> `(function namedFn() { ... }, 0);` WITHOUT mutating inside namedFn!
	script := c.GetJsRoot().GetFirstChild()
	exprRes := script.GetFirstChild()
	useCall := exprRes.GetFirstChild()
	fnNode := useCall.GetSecondChild().Detach()
	commaExpr := rhino.IR_Comma(fnNode, rhino.IR_Number(0)).Srcref(useCall)
	c.ReportChangeToEnclosingScope(useCall)
	useCall.ReplaceWith(commaExpr)

	// Iter 2: namedFn's subtree was untouched, but exprResultUsed flipped true -> false!
	// Summary cache must reject the stale context and remove both namedFn.tag = 1 and bleeding name namedFn.
	NewRemoveUnusedCode(c, opts).Process(c.GetExternsRoot(), c.GetJsRoot())

	got := renderClean(c)
	if strings.Contains(got, "namedFn") || strings.Contains(got, "tag") {
		t.Fatalf("expected namedFn name and namedFn.tag = 1 to be removed after exprResultUsed flip, got:\n%s", got)
	}
}

func TestRound97_PolyfillGuardedAndUnguardedSummaryReplay(t *testing.T) {
	externs := `
var $jscomp = {};
$jscomp.polyfill = function(target, polyfill, fromLang, toLang) {};
$jscomp.patch = function(target, polyfill) {};
var Array = function() {};
Array.prototype.flat = function() {};
var Promise = function() {};
Promise.allSettled = function(x) {};
var window = {};
function use(x) {}
`
	src := `
$jscomp.polyfill('Array.prototype.flat', function(orig) { return orig; }, 'es6', 'es5');
$jscomp.patch('Promise.allSettled', function(orig) { return orig; });

function guardedCaller(arr) {
  if (Array.prototype.flat) {
    use(arr.flat());
  }
  if (Promise.allSettled) {
    use(Promise.allSettled(arr));
  }
}

function unguardedCaller(arr) {
  return use(arr.flat());
}

window['g'] = guardedCaller;
unguardedCaller([]);
`
	c := setupNormalizedTestCompiler(t, externs, []string{src}, 1)
	opts := advancedRucOpts()

	// Iter 1: unguardedCaller keeps Array.prototype.flat alive; guardedCaller is cached with polyNodes.
	NewRemoveUnusedCode(c, opts).Process(c.GetExternsRoot(), c.GetJsRoot())

	// Delete unguardedCaller([]); on Iter 2 guardedCaller replays its polyNodes!
	// Array.prototype.flat now has ONLY guarded usages -> removed!
	// Promise.allSettled is $jscomp.patch -> guarded usage still keeps it alive!
	deleteGlobalCallTarget(t, c, "unguardedCaller")
	NewRemoveUnusedCode(c, opts).Process(c.GetExternsRoot(), c.GetJsRoot())

	if c.rucGuarded.fnCache.lastReplayed == 0 {
		t.Fatalf("expected guardedCaller to replay from summary")
	}

	cRef := setupNormalizedTestCompiler(t, externs, []string{src}, 1)
	deleteGlobalCallTarget(t, cRef, "unguardedCaller")
	prevEnv := rucIncrEnv
	rucIncrEnv = "0"
	NewRemoveUnusedCode(cRef, opts).Process(cRef.GetExternsRoot(), cRef.GetJsRoot())
	rucIncrEnv = prevEnv

	got := renderClean(c)
	want := renderClean(cRef)
	if got != want {
		t.Fatalf("polyfill replay mismatch:\ngot:\n%s\nwant:\n%s", got, want)
	}
	if strings.Contains(got, "$jscomp.polyfill") {
		t.Fatalf("expected guarded-only $jscomp.polyfill to be removed, got:\n%s", got)
	}
	if !strings.Contains(got, "$jscomp.patch") {
		t.Fatalf("expected $jscomp.patch to remain alive from guarded usage, got:\n%s", got)
	}
}

func TestRound97_DirtyInnerFunctionInvalidatesAncestors(t *testing.T) {
	externs := "var window = {}; function use(x) {}\n"
	src := `
function TargetGlobal() { return 77; }
function outerFn() {
  var outerLocal = 88;
  function midFn() {
    function deepFn() {
      return TargetGlobal() + outerLocal;
    }
    return deepFn();
  }
  return midFn();
}
window['out'] = outerFn;
`
	for _, removeGlobals := range []bool{true, false} {
		t.Run(fmt.Sprintf("removeGlobals=%v", removeGlobals), func(t *testing.T) {
			c := setupNormalizedTestCompiler(t, externs, []string{src}, 1)
			opts := RemoveUnusedCodeOptions{
				RemoveLocalVars: true,
				RemoveGlobals:   removeGlobals,
			}

			// Iter 1 caches outerFn (and midFn/deepFn in local mode).
			NewRemoveUnusedCode(c, opts).Process(c.GetExternsRoot(), c.GetJsRoot())

			// Mutate ONLY inside deepFn (`return TargetGlobal() + outerLocal;` -> `return 0;`),
			// reporting change ONLY on deepFn's return child!
			var mutated bool
			var editDeep func(n rhino.Node)
			editDeep = func(n rhino.Node) {
				if n.IsFunction() && !n.GetFirstChild().IsNil() && strings.HasPrefix(n.GetFirstChild().GetString(), "deepFn") {
					ret := n.GetLastChild().GetFirstChild()
					expr := ret.GetFirstChild()
					c.ReportChangeToEnclosingScope(expr)
					expr.ReplaceWith(rhino.IR_Number(0).Srcref(expr))
					mutated = true
					return
				}
				for ch := n.GetFirstChild(); !ch.IsNil(); ch = ch.GetNext() {
					editDeep(ch)
				}
			}
			editDeep(c.GetJsRoot())
			if !mutated {
				t.Fatalf("failed to locate deepFn")
			}

			// Iter 2 must invalidate midFn and outerFn via ancestor-closed dirtySubtreeRoots!
			NewRemoveUnusedCode(c, opts).Process(c.GetExternsRoot(), c.GetJsRoot())

			got := renderClean(c)
			if strings.Contains(got, "outerLocal") {
				t.Fatalf("expected outerLocal to be removed after deepFn edit, got:\n%s", got)
			}
			if removeGlobals && strings.Contains(got, "TargetGlobal") {
				t.Fatalf("expected TargetGlobal to be removed after deepFn edit, got:\n%s", got)
			}
		})
	}
}

func buildMultiScriptIIFEBundle(helperCount int) []string {
	var big strings.Builder
	big.WriteString("((module) => {\n")
	big.WriteString("  var shared = 100;\n")
	big.WriteString("  var deadOuter = 999;\n")
	big.WriteString("  var mutatedUnused = 1;\n")
	big.WriteString("  function mutator(x) { mutatedUnused = x; return shared; }\n")
	big.WriteString("  function withArrowArgs(keptParam, alsoKeptParam) {\n")
	big.WriteString("    const readArgs = () => arguments.length;\n")
	big.WriteString("    return readArgs() + shared;\n")
	big.WriteString("  }\n")
	for i := range helperCount {
		fmt.Fprintf(&big, "  function helper_%d(v) { return v + shared + %d; }\n", i, i)
		fmt.Fprintf(&big, "  function unused_%d() { return deadOuter + %d; }\n", i, i)
	}
	big.WriteString("  var acc = mutator(2) + withArrowArgs(10, 20);\n")
	for i := range helperCount {
		fmt.Fprintf(&big, "  acc += helper_%d(%d);\n", i, i)
	}
	big.WriteString("  module['exports'] = acc;\n")
	big.WriteString("})(window);\n")

	return []string{
		"(() => { var a = 1; var unusedA = 2; use(a); })();\n",
		"(() => { var b = 3; var unusedB = 4; use(b); })();\n",
		big.String(),
		"(() => { var c = 5; var unusedC = 6; use(c); })();\n",
	}
}

func TestRound97_SimpleMultiScriptIIFEBundleAndParallelParity(t *testing.T) {
	externs := "var window = {}; function use(x) {}\n"
	scripts := buildMultiScriptIIFEBundle(64)

	// 1. Multi-iteration direct RemoveUnusedCode in parallel (4 threads) vs sequential (1 thread) vs non-incremental.
	cPar := setupNormalizedTestCompiler(t, externs, scripts, 4)
	cSeq := setupNormalizedTestCompiler(t, externs, scripts, 1)
	cNoIncr := setupNormalizedTestCompiler(t, externs, scripts, 1)

	simpleOpts := RemoveUnusedCodeOptions{RemoveLocalVars: true}

	// Iteration 1: removes unused_* and deadOuter and mutatedUnused across all 4 scripts.
	NewRemoveUnusedCode(cPar, simpleOpts).Process(cPar.GetExternsRoot(), cPar.GetJsRoot())
	NewRemoveUnusedCode(cSeq, simpleOpts).Process(cSeq.GetExternsRoot(), cSeq.GetJsRoot())

	// Iteration 2: completely clean — all 4 top-level IIFE units skip in O(1)!
	NewRemoveUnusedCode(cPar, simpleOpts).Process(cPar.GetExternsRoot(), cPar.GetJsRoot())
	NewRemoveUnusedCode(cSeq, simpleOpts).Process(cSeq.GetExternsRoot(), cSeq.GetJsRoot())
	if cPar.rucGuarded.fnCache.lastReplayed < 4 {
		t.Fatalf("expected >= 4 top-level IIFEs skipped on clean iter 2, got %d", cPar.rucGuarded.fnCache.lastReplayed)
	}

	// Iteration 3: mutate ONE statement in the giant IIFE body (marking ONLY the outer IIFE dirty,
	// while all 64 nested helper_i + withArrowArgs functions replay their freePinIDs!).
	touchOuterIIFE := func(c *Compiler) {
		bigScript := c.GetJsRoot().GetFirstChild().GetNext().GetNext()
		iifeBody := bigScript.GetFirstChild().GetFirstChild().GetFirstChild().GetLastChild()
		lastStmt := iifeBody.GetLastChild()
		c.ReportChangeToEnclosingScope(lastStmt)
	}
	touchOuterIIFE(cPar)
	touchOuterIIFE(cSeq)

	NewRemoveUnusedCode(cPar, simpleOpts).Process(cPar.GetExternsRoot(), cPar.GetJsRoot())
	NewRemoveUnusedCode(cSeq, simpleOpts).Process(cSeq.GetExternsRoot(), cSeq.GetJsRoot())

	// 3 clean top-level IIFEs + >= 64 nested functions inside the giant IIFE should replay!
	if cPar.rucGuarded.fnCache.lastReplayed < 65 {
		t.Fatalf("expected >= 65 replayed functions on iter 3, got %d", cPar.rucGuarded.fnCache.lastReplayed)
	}
	if cSeq.rucGuarded.fnCache.lastReplayed < 65 {
		t.Fatalf("expected >= 65 replayed functions in seq mode on iter 3, got %d", cSeq.rucGuarded.fnCache.lastReplayed)
	}

	prevEnv := rucIncrEnv
	rucIncrEnv = "0"
	NewRemoveUnusedCode(cNoIncr, simpleOpts).Process(cNoIncr.GetExternsRoot(), cNoIncr.GetJsRoot())
	NewRemoveUnusedCode(cNoIncr, simpleOpts).Process(cNoIncr.GetExternsRoot(), cNoIncr.GetJsRoot())
	rucIncrEnv = prevEnv

	gotPar := renderClean(cPar)
	gotSeq := renderClean(cSeq)
	gotNoIncr := renderClean(cNoIncr)
	if gotPar != gotSeq || gotPar != gotNoIncr {
		t.Fatalf("multi-script IIFE bundle mismatch:\npar:\n%s\nseq:\n%s\nnoIncr:\n%s", gotPar, gotSeq, gotNoIncr)
	}
	if !strings.Contains(gotPar, "alsoKeptParam") {
		t.Fatalf("expected arrow arguments capture to preserve outer params, got:\n%s", gotPar)
	}
	if strings.Contains(gotPar, "deadOuter") || strings.Contains(gotPar, "mutatedUnused") || strings.Contains(gotPar, "unused_0") {
		t.Fatalf("expected dead locals to be removed, got:\n%s", gotPar)
	}

	// 2. Full end-to-end SIMPLE_OPTIMIZATIONS and ADVANCED_OPTIMIZATIONS compilation lifecycle check.
	for _, level := range []CompilationLevel{CompilationLevel_SIMPLE_OPTIMIZATIONS, CompilationLevel_ADVANCED_OPTIMIZATIONS} {
		compileWith := func(threads int, incr string) (string, *Compiler) {
			prev := rucIncrEnv
			rucIncrEnv = incr
			defer func() { rucIncrEnv = prev }()

			comp := NewCompiler()
			o := NewCompilerOptions()
			level.SetOptionsForCompilationLevel(o)
			o.SetLanguageIn(CompilerOptions_LanguageMode_ECMASCRIPT_2020)
			o.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT_2020)
			if threads > 1 {
				o.SetNumParallelThreads(threads)
			}
			ins := make([]*SourceFile, len(scripts))
			for i, s := range scripts {
				ins[i] = SourceFile_FromCode(fmt.Sprintf("bundle_%d.js", i), s)
			}
			res := comp.Compile([]*SourceFile{SourceFile_FromCode("externs.js", externs)}, ins, o)
			if len(res.Errors) > 0 {
				t.Fatalf("compile errors (%v, threads=%d, incr=%s): %v", level, threads, incr, res.Errors)
			}
			return comp.ToSource(), comp
		}

		outSeqNoIncr, _ := compileWith(1, "0")
		outSeqIncr, compSeq := compileWith(1, "1")
		outParIncr, compPar := compileWith(4, "1")
		if outSeqIncr != outSeqNoIncr || outParIncr != outSeqNoIncr {
			t.Fatalf("full compile parity mismatch for level=%v:\nseqNoIncr=%s\nseqIncr=%s\nparIncr=%s",
				level, outSeqNoIncr, outSeqIncr, outParIncr)
		}
		if compSeq.rucGuarded != nil || compPar.rucGuarded != nil {
			t.Fatalf("expected compiler.rucGuarded == nil after Compile")
		}
	}
}
