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

type round93PeepholeDeadCodeCase struct {
	CompilerTestCaseBase
}

func newRound93PeepholeDeadCodeCase() *round93PeepholeDeadCodeCase {
	c := &round93PeepholeDeadCodeCase{}
	c.InitCompilerTestCaseBaseWithExterns(c, "function sideEffect(x) {}\nfunction g() {}\nfunction upd(x) {}\n")
	return c
}

func (c *round93PeepholeDeadCodeCase) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.EnableNormalize()
	c.EnableNormalizeExpectedOutput()
	c.EnableComputeSideEffects()
	// ValidateAstChangeMarking (ChangeVerifier) is enabled by default on CompilerTestCaseBase.
}

func (c *round93PeepholeDeadCodeCase) GetNumRepetitions() int {
	return 2
}

func (c *round93PeepholeDeadCodeCase) GetProcessor(compiler *Compiler) CompilerPass {
	return CompilerPassFunc(func(externs, root rhino.Node) {
		NewPeepholeOptimizationsPass(compiler, "round93PeepholeDeadCode", NewPeepholeRemoveDeadCode()).Process(externs, root)
	})
}

func TestRound93Slice3_PeepholeRemoveDeadCodeMarkFunctionsDeleted(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return newRound93PeepholeDeadCodeCase() })
}

func (c *round93PeepholeDeadCodeCase) TestFoldEmptyTryWithFunctionInCatch() {
	c.Test("try {} catch (e) { function f() {} }", "")
	c.Test("try {} catch (e) { var x = 1; function f() { return x; } }", "var x;")
	c.Test("try {} catch (e) { function f() {} } finally { g(); }", "g();")
}

func (c *round93PeepholeDeadCodeCase) TestFoldDeadForWithFunctionInBodyAndCond() {
	c.Test("for (; false; ) { function f() {} }", "")
	c.Test("lbl: for (; false; ) { function f() {} }", "")
	c.Test("for (; (function h() {}, false); ) { function f() {} }", "")
	c.Test(
		"for (; (sideEffect(function kept() {}), false); upd(function dropped() {})) { function f() {} }",
		"sideEffect(function kept() {});",
	)
}

type round93ConcretizeStaticCase struct {
	CompilerTestCaseBase
}

func newRound93ConcretizeStaticCase() *round93ConcretizeStaticCase {
	c := &round93ConcretizeStaticCase{}
	c.InitCompilerTestCaseBaseWithExterns(c, "var $jscomp = {}; $jscomp.inherits = function(a, b, c) {}; function alert(x) {}\n")
	return c
}

func (c *round93ConcretizeStaticCase) GetProcessor(compiler *Compiler) CompilerPass {
	return NewConcretizeStaticInheritanceForInlining(compiler)
}

func TestRound93Slice3_ConcretizeStaticInheritanceArity(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return newRound93ConcretizeStaticCase() })
}

func (c *round93ConcretizeStaticCase) TestOneArgAndThreeArgInheritsDoNotCrashOrCopyWrongly() {
	// 0-arg and 1-arg $jscomp.inherits must not panic on NilNode.
	c.TestSame(
		"/** @constructor */\n" +
			"function Example() {}\n" +
			"Example.staticMethod = function() { alert(1); };\n" +
			"/** @constructor @extends {Example} */\n" +
			"function Subclass() {}\n" +
			"$jscomp.inherits();\n" +
			"$jscomp.inherits(Subclass);\n",
	)

	// 3-arg $jscomp.inherits(Subclass, Example, Extra) must be ignored rather than
	// misreading Extra as superclass and Example as subclass.
	c.TestSame(
		"/** @constructor */\n" +
			"function Example() {}\n" +
			"Example.staticMethod = function() { alert(1); };\n" +
			"/** @constructor @extends {Example} */\n" +
			"function Subclass() {}\n" +
			"/** @constructor */\n" +
			"function Extra() {}\n" +
			"Extra.otherStatic = function() { alert(2); };\n" +
			"$jscomp.inherits(Subclass, Example, Extra);\n",
	)

	// Standard 2-arg $jscomp.inherits still copies static members accurately.
	c.Test(
		"/** @constructor */\n"+
			"function Example() {}\n"+
			"Example.staticMethod = function() { alert(1); };\n"+
			"/** @constructor @extends {Example} */\n"+
			"function Subclass() {}\n"+
			"$jscomp.inherits(Subclass, Example);\n",
		"/** @constructor */\n"+
			"function Example() {}\n"+
			"Example.staticMethod = function() { alert(1); };\n"+
			"/** @constructor @extends {Example} */\n"+
			"function Subclass() {}\n"+
			"$jscomp.inherits(Subclass, Example);\n"+
			"Subclass.staticMethod = Example.staticMethod;\n",
	)
}

func TestRound93Slice3_BranchCoverageLabeledLoopAtEnd(t *testing.T) {
	cases := []struct {
		name string
		code string
	}{
		{
			name: "labeled_while_at_script_end",
			code: "function body() {}\nvar cond = false;\nlbl: while (cond) { body(); }\n",
		},
		{
			name: "labeled_while_at_function_end",
			code: "function run(cond, body) {\n  lbl: while (cond) { body(); }\n}\n",
		},
		{
			name: "nested_labels_multiline_at_function_end",
			code: "function run(cond, body) {\n  outer: inner: while (cond) {\n    if (cond) break outer;\n    body();\n  }\n}\n",
		},
		{
			name: "nested_labels_same_line_at_function_end",
			code: "function run(cond, body) {\n  outer: inner: while (cond) { if (cond) break outer; body(); }\n}\n",
		},
		{
			name: "labeled_while_before_another_label",
			code: "function run(cond, body, done) {\n  lbl1: while (cond) { body(); }\n  lbl2: { done(); }\n}\n",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			compiler := NewCompiler()
			opts := NewCompilerOptions()
			opts.SetLanguageIn(CompilerOptions_LanguageMode_ECMASCRIPT_2020)
			opts.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT_2020)
			opts.SetInstrumentForCoverageOption(CompilerOptions_InstrumentOption_BRANCH_ONLY)
			opts.SetInstrumentForCoverageOnly(true)

			res := compiler.Compile(
				[]*SourceFile{SourceFile_FromCode("externs.js", "")},
				[]*SourceFile{SourceFile_FromCode("input.js", tc.code)},
				opts,
			)
			if len(res.Errors) > 0 {
				t.Fatalf("unexpected compile errors: %v", res.Errors)
			}
			// Explicitly run AstValidator to ensure no 3-child Token_LABEL exists.
			NewAstValidator(compiler).ValidateRoot(compiler.GetRoot())
			out := compiler.ToSource()
			if !strings.Contains(out, "JSCompiler_lcov_branch_data_") {
				t.Fatalf("expected branch coverage instrumentation in output, got:\n%s", out)
			}
		})
	}
}

type round93DPAECase struct {
	CompilerTestCaseBase
}

func newRound93DPAECase() *round93DPAECase {
	c := &round93DPAECase{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func (c *round93DPAECase) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.EnableGatherExternProperties()
}

func (c *round93DPAECase) GetProcessor(compiler *Compiler) CompilerPass {
	return NewDeadPropertyAssignmentElimination(compiler)
}

func TestRound93Slice3_DeadPropertyAssignmentShortCircuitLogicalAssign(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return newRound93DPAECase() })
}

func (c *round93DPAECase) TestShortCircuitLogicalAssignmentsDoNotKillPriorWrite() {
	c.TestSame("function f(a, z) { z.w = 0; a.b &&= (z.w = 1); return z.w; }")
	c.TestSame("function f(a, z) { z.w = 0; a.b ||= (z.w = 1); return z.w; }")
	c.TestSame("function f(a, z) { z.w = 0; a.b ??= (z.w = 1); return z.w; }")
}

func buildStage3ParityInputs(numScripts int) ([]*SourceFile, []*JSChunk) {
	externs := []*SourceFile{
		SourceFile_FromCode("externs.js", "/** @externs */\nvar extObj = {keepMe: 1};\nfunction sink(x) {}\n"),
	}
	chunk0 := NewJSChunk("chunk0")
	chunk1 := NewJSChunk("chunk1")
	chunk1.AddDependency(chunk0)

	for i := 0; i < numScripts; i++ {
		var licenseHeader string
		if i%3 == 0 {
			licenseHeader = fmt.Sprintf("/** @license Shared License %d */\n", i%6)
		}
		// Mix ASCII, non-ASCII UTF-16 BMP + astral plane characters, properties, blocks, and long lines.
		code := fmt.Sprintf(
			"%s(function(s) {\n"+
				"  let msg = 'héllo-🚀-%d-' + '%s';\n"+
				"  { let msg = %d; s({alphaProp: msg, betaProp: extObj.keepMe, 'quotedProp': %d}); }\n"+
				"  s(msg);\n"+
				"})(sink);\n",
			licenseHeader,
			i,
			strings.Repeat("x", 90),
			i*7,
			i,
		)
		sf := SourceFile_FromCode(fmt.Sprintf("file_%02d.js", i), code)
		if i < numScripts/2 {
			chunk0.Add(NewCompilerInput(sf))
		} else {
			chunk1.Add(NewCompilerInput(sf))
		}
	}
	return externs, []*JSChunk{chunk0, chunk1}
}

type stage3ParityOutputs struct {
	js        string
	srcMap    string
	chunkJS   []string
	chunkMaps []string
	rendered  []RenderedChunk
}

func compileForStage3Parity(t *testing.T, numThreads int, fastMode bool) stage3ParityOutputs {
	t.Helper()
	compiler := NewCompiler()
	externs, chunks := buildStage3ParityInputs(14)
	opts := NewCompilerOptions()
	CompilationLevel_ADVANCED_OPTIMIZATIONS.SetOptionsForCompilationLevel(opts)
	opts.SetLanguageIn(CompilerOptions_LanguageMode_ECMASCRIPT_2020)
	opts.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT_2020)
	opts.SetNumParallelThreads(numThreads)
	opts.SetFastMode(fastMode)
	opts.SetLineLengthThreshold(80)
	opts.SetPrintInputDelimiter(true)
	opts.SetInputDelimiter("// Input %num%: %name%")
	opts.SetSourceMapOutputPath("out.js.map")
	opts.SetSourceMapDetailLevel(SourceMap_DetailLevel_ALL)

	res := compiler.CompileChunks(externs, chunks, opts)
	if len(res.Errors) > 0 {
		t.Fatalf("CompileChunks(threads=%d, fast=%v) errors: %v", numThreads, fastMode, res.Errors)
	}
	js := compiler.ToSource()
	var mapSB strings.Builder
	compiler.GetSourceMap().AppendTo(&mapSB, "out.js")

	lt := NewCompiler_ChunkGraphAwareLicenseTracker(compiler)
	chunkJS := make([]string, len(chunks))
	chunkMaps := make([]string, len(chunks))
	for idx := range chunks {
		compiler.GetSourceMap().Reset()
		lt.SetCurrentChunkContext(chunks[idx])
		chunkJS[idx] = compiler.ToSourceOfLicenseTrackerAndJSChunk(lt, chunks[idx])
		var chunkMapSB strings.Builder
		compiler.GetSourceMap().AppendTo(&chunkMapSB, chunks[idx].GetName()+".js")
		chunkMaps[idx] = chunkMapSB.String()
	}

	rendered, err := compiler.RenderChunksAndSourceMaps(chunks, RenderChunkOptions{
		GenerateSourceMaps: true,
		NumWorkers:         numThreads,
	})
	if err != nil {
		t.Fatalf("RenderChunksAndSourceMaps(workers=%d) error: %v", numThreads, err)
	}

	return stage3ParityOutputs{
		js:        js,
		srcMap:    mapSB.String(),
		chunkJS:   chunkJS,
		chunkMaps: chunkMaps,
		rendered:  rendered,
	}
}

func TestRound93Slice3_Stage3CodegenAndSourceMapParity(t *testing.T) {
	seq := compileForStage3Parity(t, 1, true)
	parFast := compileForStage3Parity(t, 4, true)
	parCompat := compileForStage3Parity(t, 4, false)

	if seq.js != parFast.js {
		t.Fatalf("ToSource mismatch between threads=1 and threads=4 (fast):\nseq len=%d, par len=%d", len(seq.js), len(parFast.js))
	}
	if seq.srcMap != parFast.srcMap {
		t.Fatalf("SourceMap mismatch between threads=1 and threads=4 (fast):\nseq=%s\npar=%s", seq.srcMap, parFast.srcMap)
	}
	if seq.js != parCompat.js {
		t.Fatalf("ToSource mismatch between threads=1 and threads=4 (compat)")
	}
	if seq.srcMap != parCompat.srcMap {
		t.Fatalf("SourceMap mismatch between threads=1 and threads=4 (compat)")
	}

	for idx := range seq.chunkJS {
		if seq.chunkJS[idx] != parFast.chunkJS[idx] {
			t.Fatalf("chunk %d ToSourceOfLicenseTrackerAndJSChunk mismatch (fast):\nseq=%q\npar=%q", idx, seq.chunkJS[idx], parFast.chunkJS[idx])
		}
		if seq.chunkJS[idx] != parCompat.chunkJS[idx] {
			t.Fatalf("chunk %d ToSourceOfLicenseTrackerAndJSChunk mismatch (compat)", idx)
		}
		if seq.chunkMaps[idx] != parFast.chunkMaps[idx] || seq.chunkMaps[idx] != parCompat.chunkMaps[idx] {
			t.Fatalf("chunk %d ToSourceOfLicenseTrackerAndJSChunk sourcemap mismatch", idx)
		}
	}

	if len(seq.rendered) != len(parFast.rendered) {
		t.Fatalf("rendered chunk count mismatch: %d vs %d", len(seq.rendered), len(parFast.rendered))
	}
	for i := range seq.rendered {
		if seq.rendered[i].Code != parFast.rendered[i].Code || seq.rendered[i].Code != parCompat.rendered[i].Code {
			t.Fatalf("chunk %d rendered code mismatch", i)
		}
		if seq.rendered[i].SourceMap != parFast.rendered[i].SourceMap || seq.rendered[i].SourceMap != parCompat.rendered[i].SourceMap {
			t.Fatalf("chunk %d rendered sourcemap mismatch:\nseq=%s\npar=%s", i, seq.rendered[i].SourceMap, parFast.rendered[i].SourceMap)
		}
	}
}

func TestRound93Slice3_RenamePropertiesExtraReservedAndPrevMap(t *testing.T) {
	compiler := NewCompiler()
	opts := NewCompilerOptions()
	opts.SetLanguageIn(CompilerOptions_LanguageMode_ECMASCRIPT_2020)
	opts.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT_2020)
	compiler.Init(
		[]*SourceFile{SourceFile_FromCode("externs.js", "var ext = {a: 1};\n")},
		[]*SourceFile{SourceFile_FromCode("in.js", "var o = {firstProp: 1, secondProp: 2, 'b': 3}; o.firstProp = o.secondProp;\n")},
		opts,
	)
	compiler.ParseInputs()
	NewGatherExternProperties(compiler, GatherExternProperties_Mode_OPTIMIZE).Process(compiler.GetExternsRoot(), compiler.GetJsRoot())
	compiler.SetLifeCycleStage(AbstractCompiler_LifeCycleStage_NORMALIZED)

	prevMap := newVariableMapFromEntries([]variableMapEntry{
		{orig: "secondProp", newName: "uninternedPrevPropName_r93_unique_xyz"},
	})

	rp := NewRenameProperties(
		compiler,
		false,
		prevMap,
		ASCIICharSet{},
		ASCIICharSet{},
		NewDefaultNameGenerator(),
	)
	rp.Process(compiler.GetExternsRoot(), compiler.GetJsRoot())
	propMap := rp.GetPropertyMap()
	if got := propMap.origToNew["secondProp"]; got != "uninternedPrevPropName_r93_unique_xyz" {
		t.Fatalf("expected secondProp to reuse prevUsedPropertyMap name, got %q", got)
	}
	if got := propMap.origToNew["firstProp"]; got == "a" || got == "b" || got == "" {
		t.Fatalf("expected firstProp to avoid extern 'a' and quoted 'b', got %q", got)
	}
}

func TestRound93Slice3_ContextualRenameInverterAndRucGuardedLifecycle(t *testing.T) {
	compiler := NewCompiler()
	opts := NewCompilerOptions()
	CompilationLevel_SIMPLE_OPTIMIZATIONS.SetOptionsForCompilationLevel(opts)
	opts.SetLanguageIn(CompilerOptions_LanguageMode_ECMASCRIPT_2020)
	opts.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT_2020)

	// Deeply nested scopes (>10) exercise MakeDeclaredNamesUnique ContextualRenamer frame pooling + overflow.
	src := "function outer(a) {\n" +
		"  let x = a;\n" +
		"  { let x = 1; { let x = 2; { let x = 3; { let x = 4; { let x = 5; { let x = 6; { let x = 7; { let x = 8; { let x = 9; { var hoisted = x; a += hoisted; } } } } } } } } } }\n" +
		"  return x + a;\n" +
		"}\n" +
		"window['outer'] = outer;\n"

	res := compiler.Compile(
		[]*SourceFile{SourceFile_FromCode("externs.js", "var window = {};\n")},
		[]*SourceFile{SourceFile_FromCode("in.js", src)},
		opts,
	)
	if len(res.Errors) > 0 {
		t.Fatalf("SIMPLE compile errors: %v", res.Errors)
	}
	if compiler.rucGuarded != nil {
		t.Fatalf("expected compiler.rucGuarded to be nil after compilation")
	}

	// Verify ContextualRenameInverter.Process releases transient arenas on return and can be called twice.
	inv := NewMakeDeclaredNamesUnique_ContextualRenameInverter(compiler)
	inv.Process(compiler.GetExternsRoot(), compiler.GetJsRoot())
	if inv.vars != nil || inv.scopes != nil || inv.declVars != nil || inv.renamableDecls != nil ||
		inv.scopeRefNodes != nil || inv.refLinks != nil || inv.conflictLinks != nil || inv.shadowByID != nil {
		t.Fatalf("expected ContextualRenameInverter transient arenas to be nil after Process")
	}
	inv.Process(compiler.GetExternsRoot(), compiler.GetJsRoot())
}
