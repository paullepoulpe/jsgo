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

func simpleRucOptsWithPolyfills() RemoveUnusedCodeOptions {
	return RemoveUnusedCodeOptions{
		RemoveLocalVars:       true,
		RemoveGlobals:         false,
		RemoveUnusedPolyfills: true,
	}
}

func buildRound98TypeScriptStyleBundle(helperCount, checkerCount int) []string {
	var umd strings.Builder
	umd.WriteString("(function(root, factory) {\n")
	umd.WriteString("  if (typeof exports === 'object') {\n")
	umd.WriteString("    factory(exports);\n")
	umd.WriteString("  } else {\n")
	umd.WriteString("    factory((root['umdLib'] = {}));\n")
	umd.WriteString("  }\n")
	umd.WriteString("})(this, function(exports) {\n")
	umd.WriteString("  var umdShared = 11;\n")
	umd.WriteString("  var umdDead = 22;\n")
	for i := range 20 {
		fmt.Fprintf(&umd, "  function umdKeep_%d(x) { return x + umdShared + %d; }\n", i, i)
		fmt.Fprintf(&umd, "  function umdDrop_%d() { return umdDead + %d; }\n", i, i)
	}
	umd.WriteString("  var sum = 0;\n")
	for i := range 20 {
		fmt.Fprintf(&umd, "  sum += umdKeep_%d(%d);\n", i, i)
	}
	umd.WriteString("  exports['sum'] = sum;\n")
	umd.WriteString("});\n")

	var tsBundle strings.Builder
	tsBundle.WriteString("var ts;\n")
	tsBundle.WriteString("(function(ts) {\n")
	tsBundle.WriteString("  var sharedFlags = 100;\n")
	tsBundle.WriteString("  var onlyPinnedBy0 = 31415;\n")
	tsBundle.WriteString("  var deadNamespaceLocal = 999;\n")
	tsBundle.WriteString("  var mutatedOnlyDead = 0;\n")
	tsBundle.WriteString("  function writeDead(v) { mutatedOnlyDead = v; return sharedFlags; }\n")
	tsBundle.WriteString("  function keepArrowArgs(usedArg, preservedUnusedArg) {\n")
	tsBundle.WriteString("    var countArgs = () => arguments.length;\n")
	tsBundle.WriteString("    return usedArg + countArgs() + sharedFlags;\n")
	tsBundle.WriteString("  }\n")
	tsBundle.WriteString("  function createTypeChecker(host) {\n")
	tsBundle.WriteString("    var checkerState = host + sharedFlags;\n")
	tsBundle.WriteString("    var deadCheckerVar = 888;\n")
	for i := range checkerCount {
		fmt.Fprintf(&tsBundle, "    function checkNode_%d(n) { return n + checkerState + %d; }\n", i, i)
		fmt.Fprintf(&tsBundle, "    function deadCheck_%d() { return deadCheckerVar + %d; }\n", i, i)
	}
	tsBundle.WriteString("    var total = 0;\n")
	for i := range checkerCount {
		fmt.Fprintf(&tsBundle, "    total += checkNode_%d(%d);\n", i, i)
	}
	tsBundle.WriteString("    return total;\n")
	tsBundle.WriteString("  }\n")
	fmt.Fprintf(&tsBundle, "  function helper_0(v) { return v + onlyPinnedBy0 + sharedFlags; }\n")
	for i := 1; i < helperCount; i++ {
		fmt.Fprintf(&tsBundle, "  function helper_%d(v) { return v + sharedFlags + %d; }\n", i, i)
		fmt.Fprintf(&tsBundle, "  function unusedHelper_%d() { return deadNamespaceLocal + %d; }\n", i, i)
	}
	tsBundle.WriteString("  function removableEntry() { return helper_0(1); }\n")
	tsBundle.WriteString("  var acc = writeDead(5) + keepArrowArgs(3, 4) + createTypeChecker(7);\n")
	for i := 1; i < helperCount; i++ {
		fmt.Fprintf(&tsBundle, "  acc += helper_%d(%d);\n", i, i)
	}
	tsBundle.WriteString("  removableEntry();\n")
	tsBundle.WriteString("  ts['version'] = '5.8';\n")
	tsBundle.WriteString("  ts['run'] = acc;\n")
	tsBundle.WriteString("})(ts || (ts = {}));\n")
	tsBundle.WriteString("use(ts);\n")

	return []string{
		"var globalBanner = 'keep_global_in_simple'; (() => { var a = 1; var deadA = 2; use(a); })();\n",
		umd.String(),
		tsBundle.String(),
		"(() => { var b = 3; var deadB = 4; use(b); })();\n",
		"(() => { var c = 5; var deadC = 6; use(c); })();\n",
	}
}

func deleteFirstMatchingCallStmt(t *testing.T, compiler *Compiler, calleeSubstr string) {
	t.Helper()
	var walk func(n rhino.Node) bool
	walk = func(n rhino.Node) bool {
		if n.IsCall() && !n.GetFirstChild().IsNil() && n.GetFirstChild().IsName() &&
			strings.Contains(n.GetFirstChild().GetString(), calleeSubstr) {
			stmt := n.GetParent()
			if stmt.IsExprResult() {
				NodeUtil_DeleteNode(stmt, compiler)
				return true
			}
		}
		for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
			if walk(c) {
				return true
			}
		}
		return false
	}
	if !walk(compiler.GetJsRoot()) {
		t.Fatalf("call statement matching %q not found", calleeSubstr)
	}
}

func TestRound98_SimpleModePolyfillFreeUMDBundleParallelAndLocalReplay(t *testing.T) {
	externs := "var exports; var window = {}; function use(x) {}\n"
	scripts := buildRound98TypeScriptStyleBundle(50, 20)
	opts := simpleRucOptsWithPolyfills()

	cPar := setupNormalizedTestCompiler(t, externs, scripts, 4)
	cSeq := setupNormalizedTestCompiler(t, externs, scripts, 1)
	cNoIncr := setupNormalizedTestCompiler(t, externs, scripts, 1)

	// Iteration 1: zero $jscomp.polyfill definitions exist, so SIMPLE mode unlocks
	// parallel local-script execution and rucRecLocal nested-function summary recording!
	NewRemoveUnusedCode(cPar, opts).Process(cPar.GetExternsRoot(), cPar.GetJsRoot())
	NewRemoveUnusedCode(cSeq, opts).Process(cSeq.GetExternsRoot(), cSeq.GetJsRoot())

	if cPar.rucGuarded == nil || cPar.rucGuarded.fnCache == nil {
		t.Fatalf("expected rucGuarded.fnCache after iter 1")
	}
	if cPar.rucGuarded.polyDefs.polyScripts != 0 {
		t.Fatalf("expected 0 polyfill scripts, got %d", cPar.rucGuarded.polyDefs.polyScripts)
	}
	if cPar.rucGuarded.fnCache.recMode != rucRecLocal {
		t.Fatalf("expected rucRecLocal when zero polyfills exist, got %v", cPar.rucGuarded.fnCache.recMode)
	}
	if !cPar.rucGuarded.fnCache.hadChangesDuring {
		t.Fatalf("expected iter 1 to record AST removals (hadChangesDuring=true)")
	}

	// Iteration 2: convergence pass (0 removals) — skips clean top-level IIFE/factory units.
	NewRemoveUnusedCode(cPar, opts).Process(cPar.GetExternsRoot(), cPar.GetJsRoot())
	NewRemoveUnusedCode(cSeq, opts).Process(cSeq.GetExternsRoot(), cSeq.GetJsRoot())
	if cPar.rucGuarded.fnCache.lastReplayed < 5 {
		t.Fatalf("expected >= 5 top-level functions replayed on iter 2, got %d", cPar.rucGuarded.fnCache.lastReplayed)
	}
	if cPar.rucGuarded.fnCache.hadChangesDuring {
		t.Fatalf("expected 0 changes on convergence iter 2")
	}

	// Iteration 3: zero AST changes since convergence -> O(1) idempotent fast-exit!
	NewRemoveUnusedCode(cPar, opts).Process(cPar.GetExternsRoot(), cPar.GetJsRoot())
	NewRemoveUnusedCode(cSeq, opts).Process(cSeq.GetExternsRoot(), cSeq.GetJsRoot())
	if !cPar.rucGuarded.fnCache.lastIdempotentSkip || cPar.rucGuarded.fnCache.idempotentSkipCount != 1 {
		t.Fatalf("expected idempotent fast-exit on iter 3, got skip=%v count=%d",
			cPar.rucGuarded.fnCache.lastIdempotentSkip, cPar.rucGuarded.fnCache.idempotentSkipCount)
	}

	// Iteration 4: delete `removableEntry();` inside the giant `ts` UMD namespace IIFE!
	// Only the outer `ts` IIFE is dirty; all ~50+ nested helper_i + createTypeChecker (and its 20 inner fns)
	// + all other scripts remain clean and replay their freePinIDs in O(1), while removableEntry, helper_0,
	// and onlyPinnedBy0 cascade to dead!
	deleteFirstMatchingCallStmt(t, cPar, "removableEntry")
	deleteFirstMatchingCallStmt(t, cSeq, "removableEntry")
	deleteFirstMatchingCallStmt(t, cNoIncr, "removableEntry")

	NewRemoveUnusedCode(cPar, opts).Process(cPar.GetExternsRoot(), cPar.GetJsRoot())
	NewRemoveUnusedCode(cSeq, opts).Process(cSeq.GetExternsRoot(), cSeq.GetJsRoot())

	if cPar.rucGuarded.fnCache.lastIdempotentSkip {
		t.Fatalf("expected iter 4 to wake up after AST mutation")
	}
	// Expect ~49 helper_i + createTypeChecker + keepArrowArgs + other scripts' IIFEs to replay!
	if cPar.rucGuarded.fnCache.lastReplayed < 55 {
		t.Fatalf("expected >= 55 nested/sibling functions to replay in parallel SIMPLE mode, got %d",
			cPar.rucGuarded.fnCache.lastReplayed)
	}
	if cSeq.rucGuarded.fnCache.lastReplayed < 55 {
		t.Fatalf("expected >= 55 nested/sibling functions to replay in sequential SIMPLE mode, got %d",
			cSeq.rucGuarded.fnCache.lastReplayed)
	}

	prevEnv := rucIncrEnv
	rucIncrEnv = "0"
	NewRemoveUnusedCode(cNoIncr, opts).Process(cNoIncr.GetExternsRoot(), cNoIncr.GetJsRoot())
	NewRemoveUnusedCode(cNoIncr, opts).Process(cNoIncr.GetExternsRoot(), cNoIncr.GetJsRoot())
	rucIncrEnv = prevEnv

	gotPar := renderClean(cPar)
	gotSeq := renderClean(cSeq)
	gotNoIncr := renderClean(cNoIncr)
	if gotPar != gotSeq || gotPar != gotNoIncr {
		t.Fatalf("SIMPLE UMD parity mismatch:\npar:\n%s\nseq:\n%s\nnoIncr:\n%s", gotPar, gotSeq, gotNoIncr)
	}

	for _, deadSym := range []string{
		"umdDrop_0", "umdDead", "deadNamespaceLocal", "mutatedOnlyDead", "deadCheckerVar",
		"unusedHelper_1", "removableEntry", "helper_0", "onlyPinnedBy0",
	} {
		if strings.Contains(gotPar, deadSym) {
			t.Fatalf("expected %q to be removed, got:\n%s", deadSym, gotPar)
		}
	}
	for _, keptSym := range []string{
		"globalBanner", "preservedUnusedArg", "createTypeChecker", "checkNode_0", "helper_1", "umdKeep_0",
	} {
		if !strings.Contains(gotPar, keptSym) {
			t.Fatalf("expected %q to be preserved in SIMPLE mode, got:\n%s", keptSym, gotPar)
		}
	}

	// Full end-to-end CompilationLevel_SIMPLE_OPTIMIZATIONS parity check (4 threads vs 1 thread vs JSGO_RUC_INCR=0).
	compileSimple := func(threads int, incr string) string {
		prev := rucIncrEnv
		rucIncrEnv = incr
		defer func() { rucIncrEnv = prev }()

		comp := NewCompiler()
		o := NewCompilerOptions()
		CompilationLevel_SIMPLE_OPTIMIZATIONS.SetOptionsForCompilationLevel(o)
		o.SetLanguageIn(CompilerOptions_LanguageMode_ECMASCRIPT_2020)
		o.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT_2020)
		if threads > 1 {
			o.SetNumParallelThreads(threads)
		}
		ins := make([]*SourceFile, len(scripts))
		for i, s := range scripts {
			ins[i] = SourceFile_FromCode(fmt.Sprintf("bundle_%d.js", i), s)
		}
		res := comp.Compile(
			[]*SourceFile{SourceFile_FromCode("externs.js", externs)},
			ins,
			o,
		)
		if !res.Success {
			t.Fatalf("SIMPLE_OPTIMIZATIONS compile failed (threads=%d incr=%s): %v", threads, incr, comp.GetErrors())
		}
		return comp.ToSource()
	}

	outPar := compileSimple(4, "1")
	outSeq := compileSimple(1, "1")
	outRef := compileSimple(1, "0")
	if outPar != outSeq || outPar != outRef {
		t.Fatalf("full SIMPLE_OPTIMIZATIONS output mismatch (threads=4 vs threads=1 vs JSGO_RUC_INCR=0)")
	}
}

func TestRound98_SimpleModeWithPolyfillsPreservesRemovalAndUpgradesMode(t *testing.T) {
	externs := `
var $jscomp = {};
$jscomp.polyfill = function(target, polyfill, fromLang, toLang) {};
$jscomp.patch = function(target, polyfill) {};
var Array = function() {};
Array.prototype.flat = function() {};
var String = function() {};
String.prototype.startsWith = function(s) {};
var Promise = function() {};
Promise.allSettled = function(p) {};
var window = {};
function use(x) {}
`

	t.Run("MixedUsedUnusedAndPatchedPolyfillsInSimpleMode", func(t *testing.T) {
		scripts := []string{
			`
$jscomp.polyfill('Array.prototype.flat', function(orig) { return orig; }, 'es6', 'es5');
$jscomp.polyfill('String.prototype.startsWith', function(orig) { return orig; }, 'es6', 'es5');
$jscomp.patch('Promise.allSettled', function(orig) { return orig; });
`,
			`
(function() {
  var deadLocal = 123;
  if (Array.prototype.flat) {
    use([].flat());
  }
  if (Promise.allSettled) {
    use(Promise.allSettled([]));
  }
})();
`,
			`
(function() {
  var anotherDead = 456;
  use('hello'.startsWith('he'));
})();
`,
		}

		opts := simpleRucOptsWithPolyfills()
		cPar := setupNormalizedTestCompiler(t, externs, scripts, 4)
		cRef := setupNormalizedTestCompiler(t, externs, scripts, 1)

		NewRemoveUnusedCode(cPar, opts).Process(cPar.GetExternsRoot(), cPar.GetJsRoot())
		if cPar.rucGuarded.polyDefs.polyScripts != 1 {
			t.Fatalf("expected 1 script with polyfill defs, got %d", cPar.rucGuarded.polyDefs.polyScripts)
		}
		if cPar.rucGuarded.fnCache.recMode != rucRecGlobal {
			t.Fatalf("expected rucRecGlobal when polyfills exist, got %v", cPar.rucGuarded.fnCache.recMode)
		}

		prev := rucIncrEnv
		rucIncrEnv = "0"
		NewRemoveUnusedCode(cRef, opts).Process(cRef.GetExternsRoot(), cRef.GetJsRoot())
		rucIncrEnv = prev

		got := renderClean(cPar)
		want := renderClean(cRef)
		if got != want {
			t.Fatalf("SIMPLE polyfill removal mismatch:\ngot:\n%s\nwant:\n%s", got, want)
		}
		if strings.Contains(got, "Array.prototype.flat',") || strings.Contains(got, "Array.prototype.flat\",") {
			t.Fatalf("expected guarded-only Array.prototype.flat polyfill to be removed, got:\n%s", got)
		}
		if !strings.Contains(got, "String.prototype.startsWith") {
			t.Fatalf("expected unguarded String.prototype.startsWith polyfill to be kept, got:\n%s", got)
		}
		if !strings.Contains(got, "$jscomp.patch") {
			t.Fatalf("expected $jscomp.patch(Promise.allSettled) to be kept, got:\n%s", got)
		}
	})

	t.Run("RemovingLastPolyfillUpgradesNextRunToLocalParallelMode", func(t *testing.T) {
		var umd strings.Builder
		umd.WriteString("(function(exports) {\n")
		umd.WriteString("  var shared = 7;\n")
		umd.WriteString("  var onlyByDrop = 99;\n")
		for i := range 32 {
			fmt.Fprintf(&umd, "  function fn_%d(x) { return x + shared + %d; }\n", i, i)
		}
		umd.WriteString("  function dropLater() { return onlyByDrop; }\n")
		umd.WriteString("  var s = 0;\n")
		for i := range 32 {
			fmt.Fprintf(&umd, "  s += fn_%d(%d);\n", i, i)
		}
		umd.WriteString("  dropLater();\n")
		umd.WriteString("  exports['s'] = s;\n")
		umd.WriteString("})(window);\n")

		scripts := []string{
			"$jscomp.polyfill('Array.prototype.flat', function(orig) { return orig; }, 'es6', 'es5');\n",
			umd.String(),
		}

		opts := simpleRucOptsWithPolyfills()
		c := setupNormalizedTestCompiler(t, externs, scripts, 4)

		// Iter 1: polyfill exists -> rucRecGlobal removes the unused polyfill!
		NewRemoveUnusedCode(c, opts).Process(c.GetExternsRoot(), c.GetJsRoot())
		if c.rucGuarded.fnCache.recMode != rucRecGlobal {
			t.Fatalf("expected iter 1 recMode=rucRecGlobal, got %v", c.rucGuarded.fnCache.recMode)
		}

		// Iter 2: Script 0 changed when its polyfill was deleted -> hasAnyPolyfillDef rescans Script 0,
		// sees 0 remaining polyfills, and upgrades recMode from rucRecGlobal -> rucRecLocal!
		NewRemoveUnusedCode(c, opts).Process(c.GetExternsRoot(), c.GetJsRoot())
		if c.rucGuarded.polyDefs.polyScripts != 0 {
			t.Fatalf("expected polyScripts=0 on iter 2 after polyfill removal, got %d", c.rucGuarded.polyDefs.polyScripts)
		}
		if c.rucGuarded.fnCache.recMode != rucRecLocal {
			t.Fatalf("expected iter 2 to upgrade to rucRecLocal, got %v", c.rucGuarded.fnCache.recMode)
		}

		// Iter 3: zero changes -> idempotent fast-exit!
		NewRemoveUnusedCode(c, opts).Process(c.GetExternsRoot(), c.GetJsRoot())
		if !c.rucGuarded.fnCache.lastIdempotentSkip {
			t.Fatalf("expected iter 3 idempotent fast-exit")
		}

		// Iter 4: delete dropLater() inside UMD -> wakes up and replays all 32 inner fn_i summaries!
		deleteFirstMatchingCallStmt(t, c, "dropLater")
		NewRemoveUnusedCode(c, opts).Process(c.GetExternsRoot(), c.GetJsRoot())
		if c.rucGuarded.fnCache.lastReplayed < 32 {
			t.Fatalf("expected >= 32 inner functions to replay after mode upgrade, got %d", c.rucGuarded.fnCache.lastReplayed)
		}

		cRef := setupNormalizedTestCompiler(t, externs, scripts, 1)
		deleteFirstMatchingCallStmt(t, cRef, "dropLater")
		prev := rucIncrEnv
		rucIncrEnv = "0"
		NewRemoveUnusedCode(cRef, opts).Process(cRef.GetExternsRoot(), cRef.GetJsRoot())
		NewRemoveUnusedCode(cRef, opts).Process(cRef.GetExternsRoot(), cRef.GetJsRoot())
		rucIncrEnv = prev

		got := renderClean(c)
		want := renderClean(cRef)
		if got != want {
			t.Fatalf("post-upgrade parity mismatch:\ngot:\n%s\nwant:\n%s", got, want)
		}
		if strings.Contains(got, "$jscomp.polyfill") || strings.Contains(got, "dropLater") || strings.Contains(got, "onlyByDrop") {
			t.Fatalf("expected polyfill, dropLater, and onlyByDrop to be removed, got:\n%s", got)
		}
	})

	t.Run("TypedArrayAndCollapsedPolyfillsAndDetachedScript", func(t *testing.T) {
		typedExterns := externs + `
$jscomp.polyfillTypedArrayMethod = function(name, polyfill, fromLang, toLang) {};
function $jscomp$polyfill(target, polyfill, fromLang, toLang) {}
`
		scripts := []string{
			"$jscomp.polyfillTypedArrayMethod('at', function(orig) { return orig; }, 'es_2022', 'es5');\n",
			"$jscomp$polyfill('Array.prototype.flat', function(orig) { return orig; }, 'es6', 'es5');\n",
			"(() => { var live = 1; var dead = 2; use(live); })();\n",
		}
		opts := simpleRucOptsWithPolyfills()
		c := setupNormalizedTestCompiler(t, typedExterns, scripts, 4)

		NewRemoveUnusedCode(c, opts).Process(c.GetExternsRoot(), c.GetJsRoot())
		if c.rucGuarded.fnCache.recMode != rucRecGlobal {
			t.Fatalf("expected rucRecGlobal when typed-array/collapsed polyfills exist")
		}
		got := renderClean(c)
		if strings.Contains(got, "polyfillTypedArrayMethod") || strings.Contains(got, "$jscomp$polyfill") || strings.Contains(got, "dead") {
			t.Fatalf("expected unused typed-array and collapsed polyfills to be removed, got:\n%s", got)
		}

		// Also verify detaching a script containing a polyfill wakes up both polyDefs and beginSummaryRun.
		scriptsKeep := []string{
			"$jscomp.patch('Promise.allSettled', function(orig) { return orig; }); if (Promise.allSettled) { use(1); }\n",
			"(() => { var live = 9; use(live); })();\n",
		}
		c2 := setupNormalizedTestCompiler(t, typedExterns, scriptsKeep, 4)
		NewRemoveUnusedCode(c2, opts).Process(c2.GetExternsRoot(), c2.GetJsRoot())
		NewRemoveUnusedCode(c2, opts).Process(c2.GetExternsRoot(), c2.GetJsRoot())
		if c2.rucGuarded.polyDefs.polyScripts != 1 || c2.rucGuarded.fnCache.recMode != rucRecGlobal {
			t.Fatalf("expected polyScripts=1 and rucRecGlobal before script detach")
		}
		// Detach Script 0 directly from JsRoot without touching Script 1.
		c2.GetJsRoot().GetFirstChild().Detach()
		NewRemoveUnusedCode(c2, opts).Process(c2.GetExternsRoot(), c2.GetJsRoot())
		if c2.rucGuarded.fnCache.lastIdempotentSkip {
			t.Fatalf("detaching a script must not idempotent-skip")
		}
		if c2.rucGuarded.polyDefs.polyScripts != 0 || c2.rucGuarded.fnCache.recMode != rucRecLocal {
			t.Fatalf("expected detached polyfill script to drop polyScripts to 0 and upgrade to rucRecLocal, got polyScripts=%d recMode=%v",
				c2.rucGuarded.polyDefs.polyScripts, c2.rucGuarded.fnCache.recMode)
		}
	})
}

func TestRound98_IdempotentZeroChangeFastExitAndWakeup(t *testing.T) {
	externs := "var window = {}; function use(x) {}\n"

	t.Run("SimpleModeFastExitAndMutationWakeup", func(t *testing.T) {
		scripts := []string{
			`
(function() {
  var keptLocal = 10;
  var initialDead = 20;
  var lateDead = 30;
  function readLate() { return lateDead; }
  function triggerLate() { use(readLate()); }
  use(keptLocal);
  triggerLate();
})();
`,
			`
(function() {
  var otherKept = 40;
  var otherDead = 50;
  use(otherKept);
})();
`,
		}
		opts := simpleRucOptsWithPolyfills()
		c := setupNormalizedTestCompiler(t, externs, scripts, 4)

		// Iter 1: removes initialDead and otherDead (hadChangesDuring=true).
		NewRemoveUnusedCode(c, opts).Process(c.GetExternsRoot(), c.GetJsRoot())
		if c.rucGuarded.fnCache.lastIdempotentSkip || !c.rucGuarded.fnCache.hadChangesDuring {
			t.Fatalf("unexpected iter 1 state: skip=%v changed=%v",
				c.rucGuarded.fnCache.lastIdempotentSkip, c.rucGuarded.fnCache.hadChangesDuring)
		}

		// Iter 2: convergence run (0 changes).
		NewRemoveUnusedCode(c, opts).Process(c.GetExternsRoot(), c.GetJsRoot())
		if c.rucGuarded.fnCache.lastIdempotentSkip || c.rucGuarded.fnCache.hadChangesDuring {
			t.Fatalf("unexpected iter 2 state: skip=%v changed=%v",
				c.rucGuarded.fnCache.lastIdempotentSkip, c.rucGuarded.fnCache.hadChangesDuring)
		}

		// Simulate PhaseOptimizer advancing ChangeStamp between passes with zero code edits.
		c.GetChangeTracker().IncrementChangeStamp()

		// Iter 3: O(1) idempotent fast-exit!
		NewRemoveUnusedCode(c, opts).Process(c.GetExternsRoot(), c.GetJsRoot())
		if !c.rucGuarded.fnCache.lastIdempotentSkip || c.rucGuarded.fnCache.idempotentSkipCount != 1 {
			t.Fatalf("expected iter 3 idempotent skip=true count=1, got skip=%v count=%d",
				c.rucGuarded.fnCache.lastIdempotentSkip, c.rucGuarded.fnCache.idempotentSkipCount)
		}

		// Mutate Script 0 by deleting `triggerLate();` -> wakes up RUC and removes triggerLate, readLate, lateDead!
		deleteFirstMatchingCallStmt(t, c, "triggerLate")
		NewRemoveUnusedCode(c, opts).Process(c.GetExternsRoot(), c.GetJsRoot())
		if c.rucGuarded.fnCache.lastIdempotentSkip || !c.rucGuarded.fnCache.hadChangesDuring {
			t.Fatalf("expected iter 4 wakeup with removals, got skip=%v changed=%v",
				c.rucGuarded.fnCache.lastIdempotentSkip, c.rucGuarded.fnCache.hadChangesDuring)
		}

		// Iter 5 converges, Iter 6 fast-exits again!
		NewRemoveUnusedCode(c, opts).Process(c.GetExternsRoot(), c.GetJsRoot())
		NewRemoveUnusedCode(c, opts).Process(c.GetExternsRoot(), c.GetJsRoot())
		if !c.rucGuarded.fnCache.lastIdempotentSkip || c.rucGuarded.fnCache.idempotentSkipCount != 2 {
			t.Fatalf("expected iter 6 idempotent skip count=2, got skip=%v count=%d",
				c.rucGuarded.fnCache.lastIdempotentSkip, c.rucGuarded.fnCache.idempotentSkipCount)
		}

		cRef := setupNormalizedTestCompiler(t, externs, scripts, 1)
		deleteFirstMatchingCallStmt(t, cRef, "triggerLate")
		prev := rucIncrEnv
		rucIncrEnv = "0"
		NewRemoveUnusedCode(cRef, opts).Process(cRef.GetExternsRoot(), cRef.GetJsRoot())
		rucIncrEnv = prev

		got := renderClean(c)
		want := renderClean(cRef)
		if got != want {
			t.Fatalf("SIMPLE wakeup mismatch:\ngot:\n%s\nwant:\n%s", got, want)
		}
		if strings.Contains(got, "lateDead") || strings.Contains(got, "readLate") || strings.Contains(got, "triggerLate") {
			t.Fatalf("expected newly dead local vars to be removed, got:\n%s", got)
		}
	})

	t.Run("AdvancedModeFastExitWakeupAndInstanceofCascade", func(t *testing.T) {
		// Notice `use(instCandidate instanceof UnusedCls)`:
		// Iter 1 pins instCandidate during traversal, then removes UnusedCls and folds
		// `instCandidate instanceof UnusedCls` -> `false` (hadChangesDuring=true).
		// Iter 2 MUST NOT fast-exit so it can remove instCandidate!
		src := `
var instCandidate = 123;
class UnusedCls {}
use(instCandidate instanceof UnusedCls);

function LateGlobalHelper() { return 55; }
/** @constructor */
function Box() {
  this.lateProp_ = () => LateGlobalHelper();
}
function invokeBox(b) { return use(b.lateProp_()); }
function keepEntry() { return use(9); }

window['b'] = new Box();
invokeBox(window['b']);
keepEntry();
`
		opts := advancedRucOpts()
		c := setupNormalizedTestCompiler(t, externs, []string{src}, 4)

		// Iter 1: removes UnusedCls, rewrites instanceof -> false (instCandidate still present).
		NewRemoveUnusedCode(c, opts).Process(c.GetExternsRoot(), c.GetJsRoot())
		if !c.rucGuarded.fnCache.hadChangesDuring || c.rucGuarded.fnCache.lastIdempotentSkip {
			t.Fatalf("expected iter 1 changes without skip")
		}

		// Iter 2: MUST run (not skip!) because iter 1 hadChangesDuring=true -> removes instCandidate!
		NewRemoveUnusedCode(c, opts).Process(c.GetExternsRoot(), c.GetJsRoot())
		if c.rucGuarded.fnCache.lastIdempotentSkip {
			t.Fatalf("iter 2 must not skip when iter 1 modified AST")
		}
		if strings.Contains(renderClean(c), "instCandidate") {
			t.Fatalf("expected instCandidate to be removed on iter 2, got:\n%s", renderClean(c))
		}

		// Iter 3: converges with 0 changes.
		NewRemoveUnusedCode(c, opts).Process(c.GetExternsRoot(), c.GetJsRoot())
		if c.rucGuarded.fnCache.hadChangesDuring || c.rucGuarded.fnCache.lastIdempotentSkip {
			t.Fatalf("expected iter 3 clean convergence without skip")
		}

		// Iter 4: O(1) idempotent fast-exit!
		c.GetChangeTracker().IncrementChangeStamp()
		NewRemoveUnusedCode(c, opts).Process(c.GetExternsRoot(), c.GetJsRoot())
		if !c.rucGuarded.fnCache.lastIdempotentSkip || c.rucGuarded.fnCache.idempotentSkipCount != 1 {
			t.Fatalf("expected iter 4 idempotent skip, got skip=%v count=%d",
				c.rucGuarded.fnCache.lastIdempotentSkip, c.rucGuarded.fnCache.idempotentSkipCount)
		}

		// Mutate scope before Iter 5: delete `invokeBox(window['b'])` -> wakes up and removes
		// invokeBox, Box.this.lateProp_, and LateGlobalHelper!
		deleteGlobalCallTarget(t, c, "invokeBox")
		NewRemoveUnusedCode(c, opts).Process(c.GetExternsRoot(), c.GetJsRoot())
		if c.rucGuarded.fnCache.lastIdempotentSkip {
			t.Fatalf("expected iter 5 wakeup after deleting invokeBox")
		}

		// Iter 6 converges, Iter 7 fast-exits again!
		NewRemoveUnusedCode(c, opts).Process(c.GetExternsRoot(), c.GetJsRoot())
		NewRemoveUnusedCode(c, opts).Process(c.GetExternsRoot(), c.GetJsRoot())
		if !c.rucGuarded.fnCache.lastIdempotentSkip || c.rucGuarded.fnCache.idempotentSkipCount != 2 {
			t.Fatalf("expected iter 7 idempotent skip count=2, got skip=%v count=%d",
				c.rucGuarded.fnCache.lastIdempotentSkip, c.rucGuarded.fnCache.idempotentSkipCount)
		}

		cRef := setupNormalizedTestCompiler(t, externs, []string{src}, 1)
		deleteGlobalCallTarget(t, cRef, "invokeBox")
		prev := rucIncrEnv
		rucIncrEnv = "0"
		NewRemoveUnusedCode(cRef, opts).Process(cRef.GetExternsRoot(), cRef.GetJsRoot())
		NewRemoveUnusedCode(cRef, opts).Process(cRef.GetExternsRoot(), cRef.GetJsRoot())
		rucIncrEnv = prev

		got := renderClean(c)
		want := renderClean(cRef)
		if got != want {
			t.Fatalf("ADVANCED wakeup mismatch:\ngot:\n%s\nwant:\n%s", got, want)
		}
		for _, deadSym := range []string{"UnusedCls", "instCandidate", "invokeBox", "lateProp_", "LateGlobalHelper"} {
			if strings.Contains(got, deadSym) {
				t.Fatalf("expected %q to be removed, got:\n%s", deadSym, got)
			}
		}
	})
}

func TestRound98_DeferredIndepPropsEpochAndParamScopeInvariants(t *testing.T) {
	externs := "var window = {}; function use(x) {}\n"
	src := `
function HelperA() { return 1; }
function HelperB() { return 2; }
function HelperC() { return 3; }

/** @constructor */
function MultiWidget(unusedBeforeKept, keptParam, trailingUnused1, trailingUnused2) {
  this.deadPropA_ = () => HelperA();
  this.deadPropB_ = () => HelperB();
  this.livePropC_ = HelperC() + keptParam;
}

function chainPinA(w) { return use(w.deadPropA_()); }
function chainPinB(w) { return use(w.deadPropB_()); }
function keepLiveC(w) { return use(w.livePropC_); }

window['mk'] = function(x) {
  var w = new MultiWidget(10, x, 20, 30);
  keepLiveC(w);
  return w;
};
chainPinA(new MultiWidget(1, 2, 3, 4));
chainPinB(new MultiWidget(5, 6, 7, 8));
`
	opts := advancedRucOpts()
	c := setupNormalizedTestCompiler(t, externs, []string{src}, 4)

	// Iter 1 removes trailingUnused1, trailingUnused2 and marks unusedBeforeKept as IS_UNUSED_PARAMETER.
	NewRemoveUnusedCode(c, opts).Process(c.GetExternsRoot(), c.GetJsRoot())
	// Iter 2 caches MultiWidget cleanly (now that its params were already trimmed on Iter 1).
	NewRemoveUnusedCode(c, opts).Process(c.GetExternsRoot(), c.GetJsRoot())

	// Delete chainPinA and chainPinB so deadPropA_ and deadPropB_ become unreferenced while MultiWidget is clean.
	deleteGlobalCallTarget(t, c, "chainPinA")
	deleteGlobalCallTarget(t, c, "chainPinB")
	NewRemoveUnusedCode(c, opts).Process(c.GetExternsRoot(), c.GetJsRoot())

	cRef := setupNormalizedTestCompiler(t, externs, []string{src}, 1)
	deleteGlobalCallTarget(t, cRef, "chainPinA")
	deleteGlobalCallTarget(t, cRef, "chainPinB")
	prev := rucIncrEnv
	rucIncrEnv = "0"
	NewRemoveUnusedCode(cRef, opts).Process(cRef.GetExternsRoot(), cRef.GetJsRoot())
	NewRemoveUnusedCode(cRef, opts).Process(cRef.GetExternsRoot(), cRef.GetJsRoot())
	rucIncrEnv = prev

	got := renderClean(c)
	want := renderClean(cRef)
	if got != want {
		t.Fatalf("deferred indepProp + param scope mismatch:\ngot:\n%s\nwant:\n%s", got, want)
	}
	for _, deadSym := range []string{"HelperA", "HelperB", "deadPropA_", "deadPropB_", "trailingUnused1", "trailingUnused2"} {
		if strings.Contains(got, deadSym) {
			t.Fatalf("expected %q to be removed, got:\n%s", deadSym, got)
		}
	}
	for _, keptSym := range []string{"HelperC", "livePropC_", "unusedBeforeKept", "keptParam"} {
		if !strings.Contains(got, keptSym) {
			t.Fatalf("expected %q to remain, got:\n%s", keptSym, got)
		}
	}
}
