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
	"bytes"
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type round95DefinesRunResult struct {
	source      string
	errors      []string
	warnings    []string
	defineNames []string
}

func runRound95ProcessDefines(
	t *testing.T,
	mode ProcessDefines_Mode,
	externsSrc string,
	scripts []string,
	overrides []processDefinesTest_override,
	injectFullNamespace bool,
	numThreads int,
) round95DefinesRunResult {
	t.Helper()
	compiler := NewCompiler()
	opts := NewCompilerOptions()
	opts.SetLanguageIn(CompilerOptions_LanguageMode_ECMASCRIPT_NEXT)
	opts.SetNumParallelThreads(numThreads)
	externs := []*SourceFile{
		SourceFile_FromCode("externs.js", CompilerTestCase_DEFAULT_EXTERNS+"var goog = {}; goog.define = function(a, b) {}; goog.provide = function(a) {};\n"+externsSrc),
	}
	inputs := make([]*SourceFile, len(scripts))
	for i, code := range scripts {
		inputs[i] = SourceFile_FromCode(fmt.Sprintf("input_%02d.js", i), code)
	}
	compiler.Init(externs, inputs, opts)
	compiler.ParseInputs()
	if compiler.HasErrors() {
		t.Fatalf("unexpected parse errors: %v", formatJSErrors(compiler.GetErrors()))
	}

	var replacements OrderedStringMap[rhino.Node]
	for _, ov := range overrides {
		replacements.Put(ov.name, ov.value.CloneTree())
	}
	pdOpts := ProcessDefinesOptions{
		Replacements:            replacements,
		Mode:                    mode,
		RecognizeClosureDefines: true,
	}
	if injectFullNamespace {
		ns := NewGlobalNamespaceWithRoots(compiler, compiler.GetExternsRoot(), compiler.GetJsRoot())
		pdOpts.NamespaceSupplier = func() *GlobalNamespace { return ns }
	}
	NewProcessDefines(compiler, pdOpts).Process(compiler.GetExternsRoot(), compiler.GetJsRoot())

	return round95DefinesRunResult{
		source:      compiler.ToSource(),
		errors:      formatJSErrors(compiler.GetErrors()),
		warnings:    formatJSErrors(compiler.GetWarnings()),
		defineNames: compiler.GetDefineNames().ToSlice(),
	}
}

func formatJSErrors(errs []*JSError) []string {
	out := make([]string, len(errs))
	for i, e := range errs {
		out[i] = fmt.Sprintf("%s:%d:%d:%s", e.Type().Key(), e.GetLineNumber(), e.GetCharno(), e.Description())
	}
	return out
}

func assertRound95DefinesParity(
	t *testing.T,
	externsSrc string,
	scripts []string,
	overrides []processDefinesTest_override,
) {
	t.Helper()
	modes := []ProcessDefines_Mode{
		ProcessDefines_Mode_CHECK,
		ProcessDefines_Mode_OPTIMIZE,
		ProcessDefines_Mode_CHECK_AND_OPTIMIZE,
	}
	for _, mode := range modes {
		want := runRound95ProcessDefines(t, mode, externsSrc, scripts, overrides, true, 1)
		for _, threads := range []int{1, 4} {
			got := runRound95ProcessDefines(t, mode, externsSrc, scripts, overrides, false, threads)
			if got.source != want.source {
				t.Fatalf("mode=%s threads=%d source mismatch:\ngot:\n%s\nwant:\n%s", mode, threads, got.source, want.source)
			}
			if !slices.Equal(got.errors, want.errors) {
				t.Fatalf("mode=%s threads=%d errors mismatch:\ngot:  %v\nwant: %v", mode, threads, got.errors, want.errors)
			}
			if !slices.Equal(got.warnings, want.warnings) {
				t.Fatalf("mode=%s threads=%d warnings mismatch:\ngot:  %v\nwant: %v", mode, threads, got.warnings, want.warnings)
			}
			if !slices.Equal(got.defineNames, want.defineNames) {
				t.Fatalf("mode=%s threads=%d defineNames mismatch:\ngot:  %v\nwant: %v", mode, threads, got.defineNames, want.defineNames)
			}
		}
	}
}

func makeNoiseScripts(count int) []string {
	scripts := make([]string, count)
	for i := range count {
		scripts[i] = fmt.Sprintf(
			"goog.provide('noise.mod%d');\n"+
				"const TOP_UNUSED_%d = %d;\n"+
				"function work_%d(obj) {\n"+
				"  const id = obj.id;\n"+
				"  const name = obj.name;\n"+
				"  const value = obj.value;\n"+
				"  return id + name + value;\n"+
				"}\n",
			i, i, i, i,
		)
	}
	return scripts
}

func TestRound95Defines_LocalConstDoesNotPolluteCandidates(t *testing.T) {
	scripts := makeNoiseScripts(12)
	scripts[5] = "/** @define {boolean} */ var ONLY_DEFINE = true;\n" + scripts[5]

	compiler := NewCompiler()
	opts := NewCompilerOptions()
	opts.SetLanguageIn(CompilerOptions_LanguageMode_ECMASCRIPT_NEXT)
	opts.SetNumParallelThreads(4)
	inputs := make([]*SourceFile, len(scripts))
	for i, s := range scripts {
		inputs[i] = SourceFile_FromCode(fmt.Sprintf("s%d.js", i), s)
	}
	compiler.Init([]*SourceFile{SourceFile_FromCode("ext.js", "")}, inputs, opts)
	compiler.ParseInputs()

	ns := NewGlobalNamespaceWithRoots(compiler, compiler.GetExternsRoot(), compiler.GetJsRoot())
	ns.forProcessDefines = true
	ns.processDefinesCheck = true
	ns.EnsureGenerated()

	if len(ns.defineCandidateLeafIDs) != 1 {
		t.Fatalf("expected exactly 1 candidate leaf ID (ONLY_DEFINE), got %d", len(ns.defineCandidateLeafIDs))
	}
	if ns.LookupName("ONLY_DEFINE") == nil {
		t.Fatalf("expected ONLY_DEFINE in namespace")
	}
	if ns.LookupName("TOP_UNUSED_0") != nil {
		t.Fatalf("unexpected TOP_UNUSED_0 in primary define namespace when includeGlobalConsts=false")
	}
}

func TestRound95Defines_ZeroDefinesFastExit(t *testing.T) {
	scripts := makeNoiseScripts(14)
	scripts[2] = "var CLOSURE_DEFINES = {'COMPILED': true, 'EXTRA_UNKNOWN': 1};\n" + scripts[2]
	overrides := []processDefinesTest_override{
		{name: "CLI_UNKNOWN", value: rhino.NewNode(rhino.Token_TRUE)},
	}
	assertRound95DefinesParity(t, "", scripts, overrides)
}

func TestRound95Defines_ChainedGlobalConstsAndDefines(t *testing.T) {
	scripts := makeNoiseScripts(14)
	scripts[0] = "const BASE_NUM = 10;\n/** @const */ var MULT_NUM = 3;\n"
	scripts[3] = "var cfg = {};\n/** @const */ cfg.OFFSET = BASE_NUM + 5;\n"
	scripts[6] = "/** @define {boolean} */ var DEF_ENABLED = goog.define('DEF_ENABLED', true);\n"
	scripts[9] = "/** @const */ var DEF_DISABLED = !DEF_ENABLED;\n"
	scripts[12] = "/** @define {number} */ const COMPUTED_DEF = cfg.OFFSET * MULT_NUM;\n" +
		"/** @define {boolean} */ var MIRROR_DEF = DEF_DISABLED;\n"

	overrides := []processDefinesTest_override{
		{name: "DEF_ENABLED", value: rhino.NewNode(rhino.Token_FALSE)},
		{name: "COMPUTED_DEF", value: rhino.IR_Number(99)},
	}
	assertRound95DefinesParity(t, "", scripts, overrides)
}

func TestRound95Defines_ReassignedConstAcrossScriptsRejected(t *testing.T) {
	scripts := makeNoiseScripts(12)
	scripts[1] = "/** @const */ var MAYBE_CONST = 7;\n"
	scripts[7] = "function mutateLater() { MAYBE_CONST = 8; }\n"
	scripts[10] = "/** @define {number} */ var BAD_FROM_MUTATED = MAYBE_CONST + 1;\n"

	assertRound95DefinesParity(t, "", scripts, nil)
	res := runRound95ProcessDefines(t, ProcessDefines_Mode_CHECK, "", scripts, nil, false, 4)
	if len(res.errors) != 1 || !strings.HasPrefix(res.errors[0], "JSC_INVALID_DEFINE_VALUE:") {
		t.Fatalf("expected JSC_INVALID_DEFINE_VALUE, got %v", res.errors)
	}
}

func TestRound95Defines_NonConstDefineReassignmentsAcrossScripts(t *testing.T) {
	scripts := makeNoiseScripts(14)
	scripts[0] = "var ns = {};\n/** @define {boolean} */ ns.FLAG = true;\n/** @define {number} */ var COUNT = 1;\n"
	scripts[5] = "function clobberFn() { COUNT = 2; }\n"
	scripts[11] = "ns = {FLAG: false};\n"

	assertRound95DefinesParity(t, "", scripts, nil)
	res := runRound95ProcessDefines(t, ProcessDefines_Mode_CHECK_AND_OPTIMIZE, "", scripts, nil, false, 4)
	if len(res.errors) != 2 {
		t.Fatalf("expected 2 JSC_NON_CONST_DEFINE errors, got %v", res.errors)
	}
}

func TestRound95Defines_MisplacedLocalDefinesWithZeroGlobalDefines(t *testing.T) {
	scripts := makeNoiseScripts(14)
	scripts[1] = "function f() { /** @define {boolean} */ var localDef = true; }\n"
	scripts[4] = "var bag = { /** @define {boolean} */ propDef: false };\n"
	scripts[8] = "var missingDoc = goog.define('MISSING_DOC', 1);\n"
	scripts[10] = "function g() { var x = goog.define('LOCAL_GOOG_DEF', true); }\n"
	scripts[13] = "function h() { var CLOSURE_DEFINES = {'FOO': 1}; }\n"

	assertRound95DefinesParity(t, "", scripts, nil)
	res := runRound95ProcessDefines(t, ProcessDefines_Mode_CHECK, "", scripts, nil, false, 4)
	if len(res.errors) != 5 {
		t.Fatalf("expected 5 invalid-location/goog.define/CLOSURE_DEFINES errors, got %v", res.errors)
	}
}

func TestRound95Defines_GoogProvideAndLegacyModuleBlockOrder(t *testing.T) {
	scripts := makeNoiseScripts(14)
	scripts[0] = "goog.provide('alpha.beta.EARLY_FLAG');\ngoog.provide('pkg.sub');\n"
	scripts[2] = "{\n  goog.provide('legacy.mod');\n  const scopedHelper = 1;\n}\n"
	scripts[4] = "/** @define {boolean} */ var MID_FLAG = true;\n"
	scripts[9] = "/** @define {boolean} */ pkg.sub.ENABLE_X = goog.define('pkg.sub.ENABLE_X', false);\n"
	scripts[11] = "/** @define {number} */ legacy.mod.LIMIT = 42;\n"
	scripts[13] = "/** @define {boolean} */ alpha.beta.EARLY_FLAG = false;\n"

	overrides := []processDefinesTest_override{
		{name: "pkg.sub.ENABLE_X", value: rhino.NewNode(rhino.Token_TRUE)},
		{name: "legacy.mod.LIMIT", value: rhino.IR_Number(100)},
		{name: "alpha.beta.EARLY_FLAG", value: rhino.NewNode(rhino.Token_TRUE)},
	}
	assertRound95DefinesParity(t, "", scripts, overrides)
}

func TestRound95Defines_FunctionAndDestructuringEdgeCases(t *testing.T) {
	scripts := makeNoiseScripts(12)
	scripts[1] = "/** @define {boolean} */ function badFnDefine() { return true; }\n"
	scripts[4] = "var holder = {};\n/** @define {boolean} */ holder.D = true;\n"
	scripts[7] = "const {D: copied} = holder;\n"
	scripts[10] = "var viaInit = (/** @define {number} */ holder.NUM = 5);\n"

	assertRound95DefinesParity(t, "", scripts, nil)
}

func TestRound95Defines_JvmMetricsFormatting(t *testing.T) {
	var empty bytes.Buffer
	JvmMetrics_MaybeWriteJvmMetrics(&empty, base.NoString())
	if empty.Len() != 0 {
		t.Fatalf("expected empty output for NoString, got %q", empty.String())
	}

	var pretty bytes.Buffer
	JvmMetrics_MaybeWriteJvmMetricsOpt(&pretty, base.SomeString("verbose:pretty:all"))
	prettyStr := pretty.String()
	for _, wantSub := range []string{
		"Garbage collection stats\n\tCollector Go GC\n\t\tcollection count   :",
		"Memory usage\n\tHeap\n\t\tavailable         :",
		"\tPool Go heap\n",
		"JIT Stats\n\tGo jit time: 0 ms\n",
	} {
		if !strings.Contains(prettyStr, wantSub) {
			t.Fatalf("missing %q in verbose:pretty:all output:\n%s", wantSub, prettyStr)
		}
	}

	var tabular bytes.Buffer
	JvmMetrics_WriteJvmMetrics(&tabular, "compact:tabular:gc,mem,jit")
	for _, line := range strings.Split(strings.TrimSpace(tabular.String()), "\n") {
		colon := strings.IndexByte(line, ':')
		if colon != jvmTabularColonPos {
			t.Fatalf("expected colon at column %d in %q, got %d", jvmTabularColonPos, line, colon)
		}
	}
}
