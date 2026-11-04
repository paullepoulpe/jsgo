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
	"strconv"
	"strings"
	"testing"
)

const round98Externs = `
/** @constructor */ function Function() {}
/** @type {function(*, ...*): *} */
Function.prototype.call = function(thisArg, ...args) {};
/** @nosideeffects */ function pureExt(x) {}
function sink(x) {}
var globalObj = {};
`

func buildRound98GiantUMDSource(prefix string, count int) string {
	var b strings.Builder
	fmt.Fprintf(&b, "var %s_pre = function(x) { return x + 10; };\n", prefix)
	fmt.Fprintf(&b, "sink(%s_pre(1));\n", prefix)
	fmt.Fprintf(&b, "(function(exports) {\n")
	for i := 0; i < count; i++ {
		switch i % 8 {
		case 0:
			// Direct-inlinable helper called inside decomposed expression.
			fmt.Fprintf(&b,
				"  function %s_dir_%d(a, b) { return a + b + %d; }\n",
				prefix, i, i)
		case 1:
			// Block-inlinable helper with local var and side effect.
			fmt.Fprintf(&b,
				"  var %s_blk_%d = function(a, b) { var t = a * 2 + %d; sink(t); return t + b; };\n",
				prefix, i, i)
		case 2:
			// Multi-step chain: calls previous direct and block helpers (requires multiple iterations).
			if i >= 2 {
				fmt.Fprintf(&b,
					"  function %s_chain_%d(x) { return %s_blk_%d(%s_dir_%d(x, 1), %d); }\n",
					prefix, i, prefix, i-1, prefix, i-2, i)
			} else {
				fmt.Fprintf(&b,
					"  function %s_chain_%d(x) { return x + %d; }\n",
					prefix, i, i)
			}
		case 3:
			// Shadows outer UMD `exports` parameter and uses `.call(this, ...)`.
			fmt.Fprintf(&b,
				"  function %s_shad_%d(exports, y) { var z = exports + y + %d; return z * 3; }\n",
				prefix, i, i)
		case 4:
			// Reassigned function (must be disallowed in Phase 2).
			fmt.Fprintf(&b,
				"  var %s_reass_%d = function(x) { return x + %d; };\n",
				prefix, i, i)
		case 5:
			// Const alias + @noinline or multi-caller large function kept by cost model.
			fmt.Fprintf(&b,
				"  function %s_heavy_%d(a, b, c) { sink(a + %d); var q = a * b - c; sink(q + b); sink(q - c); return q + %d; }\n",
				prefix, i, i, i+1)
		case 6:
			// Arrow helper + inner anonymous IIFE at call site.
			fmt.Fprintf(&b,
				"  const %s_arr_%d = (u, v) => u * %d + v;\n",
				prefix, i, i+1)
		case 7:
			// Escaped / unremovable function reference.
			fmt.Fprintf(&b,
				"  function %s_esc_%d(k) { return k - %d; }\n",
				prefix, i, i)
		}
	}

	b.WriteString("  var acc = 0;\n")
	for i := 0; i < count; i++ {
		switch i % 8 {
		case 0:
			// Call in expression requiring decomposition (sink(acc) before call).
			fmt.Fprintf(&b, "  acc += sink(acc) + %s_dir_%d(sink(%d), acc);\n", prefix, i, i)
		case 1:
			// Block-inlined in expression context (forces decomposition + temp vars).
			fmt.Fprintf(&b, "  acc = acc + %s_blk_%d(sink(%d), acc);\n", prefix, i, i)
		case 2:
			fmt.Fprintf(&b, "  acc += %s_chain_%d(acc + %d);\n", prefix, i, i)
		case 3:
			fmt.Fprintf(&b, "  acc += %s_shad_%d.call(this, acc, %d) + (function(w) { var m = w + %d; sink(m); return m; })(acc);\n",
				prefix, i, i, i)
		case 4:
			fmt.Fprintf(&b, "  if (acc > %d) { %s_reass_%d = sink; } acc += %s_reass_%d(acc);\n",
				i, prefix, i, prefix, i)
		case 5:
			fmt.Fprintf(&b, "  acc += %s_heavy_%d(acc, 1, 2) + %s_heavy_%d(3, acc, 4) + %s_heavy_%d(5, 6, acc);\n",
				prefix, i, prefix, i, prefix, i)
		case 6:
			fmt.Fprintf(&b, "  acc += %s_arr_%d(sink(%d), (function(p) { return p + 1; })(acc));\n",
				prefix, i, i)
		case 7:
			fmt.Fprintf(&b, "  sink(%s_esc_%d); acc += %s_esc_%d(acc);\n", prefix, i, prefix, i)
		}
	}
	b.WriteString("  exports.result = acc;\n")
	b.WriteString("})({ get exports() { return globalObj; }, set result(v) { sink(v); } });\n")
	fmt.Fprintf(&b, "var %s_post = function(z) { sink(z); return z + 99; };\n", prefix)
	fmt.Fprintf(&b, "sink(%s_post(globalObj));\n", prefix)
	return b.String()
}

func runInlineFunctionsLoopForTest(
	t *testing.T,
	scripts []string,
	threads int,
	reach CompilerOptions_Reach,
	maxSize int,
	iters int,
) (perIterCode []string, perIterCounter []int) {
	t.Helper()
	options := NewCompilerOptions()
	if reach == CompilerOptions_Reach_ALL {
		CompilationLevel_ADVANCED_OPTIMIZATIONS.SetOptionsForCompilationLevel(options)
	} else {
		CompilationLevel_SIMPLE_OPTIMIZATIONS.SetOptionsForCompilationLevel(options)
	}
	options.SetNumParallelThreads(threads)

	c := NewCompiler()
	inputs := make([]*SourceFile, len(scripts))
	for i, s := range scripts {
		inputs[i] = SourceFile_FromCode(fmt.Sprintf("script_%d.js", i), s)
	}
	c.Init(
		[]*SourceFile{SourceFile_FromCode("externs.js", round98Externs)},
		inputs,
		options,
	)
	c.ParseInputs()
	if c.HasErrors() {
		t.Fatalf("parse errors (threads=%d): %v", threads, c.GetErrors())
	}
	NewNormalize(c, false).Process(c.externsRoot, c.jsRoot)
	c.SetLifeCycleStage(AbstractCompiler_LifeCycleStage_NORMALIZED)

	nextID := 0
	supplier := func() string {
		v := nextID
		nextID++
		return "r98_" + strconv.Itoa(v)
	}

	perIterCode = make([]string, 0, iters)
	perIterCounter = make([]int, 0, iters)
	for iter := 0; iter < iters; iter++ {
		pass := NewInlineFunctions(c, supplier, reach, false, maxSize)
		pass.Process(c.externsRoot, c.jsRoot)
		if c.HasErrors() {
			t.Fatalf("iter %d errors (threads=%d): %v", iter, threads, c.GetErrors())
		}
		perIterCode = append(perIterCode, c.ToSource())
		perIterCounter = append(perIterCounter, nextID)
	}
	return perIterCode, perIterCounter
}

func compileFullForRound98(
	t *testing.T,
	scripts []string,
	level CompilationLevel,
	threads int,
) string {
	t.Helper()
	opts := NewCompilerOptions()
	level.SetOptionsForCompilationLevel(opts)
	opts.SetLanguageIn(CompilerOptions_LanguageMode_ECMASCRIPT_NEXT)
	opts.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT_2020)
	opts.SetNumParallelThreads(threads)

	c := NewCompiler()
	inputs := make([]*SourceFile, len(scripts))
	for i, s := range scripts {
		inputs[i] = SourceFile_FromCode(fmt.Sprintf("in_%d.js", i), s)
	}
	res := c.Compile(
		[]*SourceFile{SourceFile_FromCode("externs.js", round98Externs)},
		inputs,
		opts,
	)
	if !res.Success {
		t.Fatalf("compile failed (level=%v, threads=%d): %v", level, threads, res.Errors)
	}
	return c.ToSource()
}

func assertInlineLoopParity(
	t *testing.T,
	label string,
	scripts []string,
	reach CompilerOptions_Reach,
	maxSize int,
	iters int,
) {
	t.Helper()
	wantCode, wantIDs := runInlineFunctionsLoopForTest(t, scripts, 1, reach, maxSize, iters)
	for _, threads := range []int{2, 4, 8} {
		gotCode, gotIDs := runInlineFunctionsLoopForTest(t, scripts, threads, reach, maxSize, iters)
		for i := 0; i < iters; i++ {
			if gotIDs[i] != wantIDs[i] {
				t.Fatalf("%s (threads=%d, iter=%d): safeNameIdSupplier counter = %d, want %d",
					label, threads, i, gotIDs[i], wantIDs[i])
			}
			if gotCode[i] != wantCode[i] {
				t.Fatalf("%s (threads=%d, iter=%d): output mismatch vs threads=1\ngot (len=%d):\n%.600s\nwant (len=%d):\n%.600s",
					label, threads, i, len(gotCode[i]), gotCode[i], len(wantCode[i]), wantCode[i])
			}
		}
	}
}

func TestRound98_SingleAndMultiScriptGiantUMDBundleParity(t *testing.T) {
	giant := buildRound98GiantUMDSource("ts", 320)
	small0 := "function s0(a) { return a + 1; } sink(s0(10));\n"
	small1 := "(function() { function s1(a) { var y = a * 3; sink(y); return y; } sink(s1(5)); })();\n"
	small2 := "var s2 = function(a, b) { return pureExt(a) ? b : a; }; sink(s2(1, 2));\n"

	for _, cacheCfg := range []struct {
		name   string
		env    string
		force  bool
		verify bool
	}{
		{name: "cacheVerify", env: "", force: true, verify: true},
		{name: "cacheFast", env: "", force: true, verify: false},
		{name: "uncached", env: "0", force: false, verify: false},
	} {
		t.Run(cacheCfg.name, func(t *testing.T) {
			withInlineFnCacheConfig(t, cacheCfg.env, cacheCfg.force, cacheCfg.verify)

			// Single-script giant UMD bundle (SIMPLE / LOCALS_ONLY and ADVANCED / ALL).
			assertInlineLoopParity(t, "single-locals", []string{giant},
				CompilerOptions_Reach_LOCAL_ONLY, CompilerOptions_UNLIMITED_FUN_SIZE_AFTER_INLINING, 4)
			assertInlineLoopParity(t, "single-all-maxSize", []string{giant},
				CompilerOptions_Reach_ALL, 150, 4)

			// Multi-script bundle (3 small scripts + 1 giant UMD script, mirroring bundle benchmark).
			multi := []string{small0, small1, small2, giant}
			assertInlineLoopParity(t, "multi-locals", multi,
				CompilerOptions_Reach_LOCAL_ONLY, CompilerOptions_UNLIMITED_FUN_SIZE_AFTER_INLINING, 4)
			assertInlineLoopParity(t, "multi-all", multi,
				CompilerOptions_Reach_ALL, CompilerOptions_UNLIMITED_FUN_SIZE_AFTER_INLINING, 5)
		})
	}

	// End-to-end Compiler.Compile parity (SIMPLE and ADVANCED).
	multi := []string{small0, small1, small2, giant}
	for _, lvl := range []CompilationLevel{
		CompilationLevel_SIMPLE_OPTIMIZATIONS,
		CompilationLevel_ADVANCED_OPTIMIZATIONS,
	} {
		want := compileFullForRound98(t, multi, lvl, 1)
		got := compileFullForRound98(t, multi, lvl, 4)
		if got != want {
			t.Fatalf("full compile mismatch (level=%v) between threads=1 and threads=4", lvl)
		}
	}
}

func TestRound98_WrapperIIFEVariantsAndInlinedOuterBlock(t *testing.T) {
	withInlineFnCacheConfig(t, "", true, false)

	buildBody := func(tag string, n int) string {
		var b strings.Builder
		for i := 0; i < n; i++ {
			switch i % 3 {
			case 0:
				fmt.Fprintf(&b, "  function %s_f%d(x) { return x + %d; }\n", tag, i, i)
			case 1:
				fmt.Fprintf(&b, "  function %s_g%d(x) { var v = %s_f%d(x) * 2; sink(v); return v; }\n", tag, i, tag, i-1)
			case 2:
				fmt.Fprintf(&b, "  sink(%s_g%d(sink(%d)));\n", tag, i-1, i)
			}
		}
		return b.String()
	}

	// 1a. Zero-param arrow IIFE calling a top-level block-inlinable helper 150 times:
	// iter 1 unwraps the outer IIFE and block-inlines it into a top-level BLOCK of 150 statements;
	// iter 2 unwraps that top-level BLOCK in both Phase 1 and Phase 2 and inlines the helper calls.
	var zb strings.Builder
	zb.WriteString("function topBlk(x) { var r = x + 7; sink(r); return r; }\n(() => {\n")
	for i := 0; i < 150; i++ {
		if i%10 == 0 {
			fmt.Fprintf(&zb, "  sink(topBlk(sink(%d)));\n", i)
		} else {
			fmt.Fprintf(&zb, "  sink(%d);\n", i)
		}
	}
	zb.WriteString("})();\n")
	zeroParamIIFE := zb.String()
	{
		c := newNormalizedTestCompiler(t, round98Externs, zeroParamIIFE)
		c.GetOptions().SetNumParallelThreads(4)
		inf := NewInlineFunctions(c, func() string { return "id" }, CompilerOptions_Reach_ALL, false,
			CompilerOptions_UNLIMITED_FUN_SIZE_AFTER_INLINING)
		plan1, workers1, ok1 := inf.buildIfShardPlan(c.jsRoot, true)
		if !ok1 || workers1 <= 1 || len(plan1.unwrappedFns) != 1 || len(plan1.units) < 3 {
			t.Fatalf("iter 1 expected unwrapped IIFE shard plan, got ok=%v workers=%d fns=%d units=%d",
				ok1, workers1, len(plan1.unwrappedFns), len(plan1.units))
		}
		inf.Process(c.externsRoot, c.jsRoot)

		plan2P1, workers2P1, ok2P1 := inf.buildIfShardPlan(c.jsRoot, true)
		plan2P2, workers2P2, ok2P2 := inf.buildIfShardPlan(c.jsRoot, false)
		if !ok2P1 || workers2P1 <= 1 || len(plan2P1.unwrappedBlocks) < 1 || len(plan2P1.units) < 3 {
			t.Fatalf("iter 2 Phase 1 expected unwrapped top-level BLOCK shard plan, got ok=%v workers=%d blks=%d units=%d",
				ok2P1, workers2P1, len(plan2P1.unwrappedBlocks), len(plan2P1.units))
		}
		if !ok2P2 || workers2P2 <= 1 || len(plan2P2.unwrappedBlocks) < 1 || len(plan2P2.units) < 3 {
			t.Fatalf("iter 2 Phase 2 expected unwrapped top-level BLOCK shard plan, got ok=%v workers=%d blks=%d units=%d",
				ok2P2, workers2P2, len(plan2P2.unwrappedBlocks), len(plan2P2.units))
		}
	}
	assertInlineLoopParity(t, "zero-param-iife-block-inlined", []string{zeroParamIIFE},
		CompilerOptions_Reach_ALL, CompilerOptions_UNLIMITED_FUN_SIZE_AFTER_INLINING, 3)

	// 1b. Parameterized arrow IIFE with 150 inner functions + >150 call statements:
	// iters 1-2 inline inner functions; iter 3 block-inlines the outer IIFE into a top-level BLOCK;
	// iter 4 shards that top-level BLOCK.
	arrowIIFE := "((module) => {\n  var localSeed = module;\n" + buildBody("arr", 420) + "  sink(localSeed);\n})(globalObj);\n"
	assertInlineLoopParity(t, "arrow-iife-block-inlined", []string{arrowIIFE},
		CompilerOptions_Reach_ALL, CompilerOptions_UNLIMITED_FUN_SIZE_AFTER_INLINING, 5)

	// 2. `.call(this, ...)` wrapper, unary `!function(){}()`, and UMD `(function(r,f){f(r)})(this, function(e){...})`.
	dotCall := "(function(exports, extra) {\n" + buildBody("dc", 150) + "  sink(exports); sink(extra);\n}).call(this, globalObj, sink(1));\n"
	bangIIFE := "!function(exports) {\n" + buildBody("bg", 150) + "  sink(exports);\n}(globalObj);\n"
	umdFactory := "(function(root, factory) { factory(root); })(this, function(exports) {\n" + buildBody("umd", 156) + "  sink(exports);\n});\n"

	for name, src := range map[string]string{
		"dotCall":    dotCall,
		"bangIIFE":   bangIIFE,
		"umdFactory": umdFactory,
	} {
		assertInlineLoopParity(t, name+"-locals", []string{src},
			CompilerOptions_Reach_LOCAL_ONLY, CompilerOptions_UNLIMITED_FUN_SIZE_AFTER_INLINING, 3)
		assertInlineLoopParity(t, name+"-all", []string{src},
			CompilerOptions_Reach_ALL, CompilerOptions_UNLIMITED_FUN_SIZE_AFTER_INLINING, 4)
	}
}

func TestRound98_IIFEDefaultParamsMultipleIIFEsFlatScriptAndESModule(t *testing.T) {
	withInlineFnCacheConfig(t, "", true, true)

	makeStmts := func(tag string, n int) string {
		var b strings.Builder
		for i := 0; i < n; i++ {
			if i%2 == 0 {
				fmt.Fprintf(&b, "function %s_fn_%d(a) { var r = a + %d; sink(r); return r; }\n", tag, i, i)
			} else {
				fmt.Fprintf(&b, "sink(%s_fn_%d(sink(%d)));\n", tag, i-1, i)
			}
		}
		return b.String()
	}

	// Edge case 1: Named IIFE with default parameters referencing globals & candidate functions.
	defaultParamIIFE := `
function topHelper(v) { return v + 42; }
(function namedWrapper(a = topHelper(sink(1)), b = (function(z) { return z * 2; })(3)) {
` + makeStmts("dp", 160) + `
  if (a > 1000) { sink(namedWrapper); }
  sink(a + b);
})(undefined, sink(2));
`
	assertInlineLoopParity(t, "default-param-named-iife-all", []string{defaultParamIIFE},
		CompilerOptions_Reach_ALL, CompilerOptions_UNLIMITED_FUN_SIZE_AFTER_INLINING, 4)
	assertInlineLoopParity(t, "default-param-named-iife-locals", []string{defaultParamIIFE},
		CompilerOptions_Reach_LOCAL_ONLY, CompilerOptions_UNLIMITED_FUN_SIZE_AFTER_INLINING, 3)

	// Edge case 2: Multiple large IIFEs interleaved with top-level statements in one script.
	multiIIFE := "var pre = function(x) { return x - 1; }; sink(pre(9));\n" +
		"(function(e1) {\n" + makeStmts("m1", 150) + "sink(e1);\n})(globalObj);\n" +
		"var mid = function(x) { return x + 2; }; sink(mid(8));\n" +
		"!function(e2) {\n" + makeStmts("m2", 150) + "sink(e2);\n}(globalObj);\n" +
		"var post = function(x) { return x * 3; }; sink(post(7));\n"
	assertInlineLoopParity(t, "multiple-large-iifes", []string{multiIIFE},
		CompilerOptions_Reach_ALL, CompilerOptions_UNLIMITED_FUN_SIZE_AFTER_INLINING, 4)

	// Edge case 3: Large flat top-level SCRIPT (>= 200 direct statements, no wrapper IIFE).
	flatScript := makeStmts("flat", 220)
	assertInlineLoopParity(t, "flat-large-script", []string{flatScript},
		CompilerOptions_Reach_ALL, CompilerOptions_UNLIMITED_FUN_SIZE_AFTER_INLINING, 4)

	// Edge case 4: ES module (MODULE_BODY) script alongside a large UMD script.
	esModule := "export function exportedKeep(x) { return x + 1; }\n" + makeStmts("mod", 180)
	modWant := compileFullForRound98(t, []string{esModule, multiIIFE}, CompilationLevel_SIMPLE_OPTIMIZATIONS, 1)
	modGot := compileFullForRound98(t, []string{esModule, multiIIFE}, CompilationLevel_SIMPLE_OPTIMIZATIONS, 4)
	if modGot != modWant {
		t.Fatalf("ES module + UMD compile mismatch between threads=1 and threads=4:\ngot:\n%.500s\nwant:\n%.500s", modGot, modWant)
	}
}
