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

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

func TestRound96_ParallelDAEArrowInDefaultParam(t *testing.T) {
	var srcBuilder strings.Builder
	var wantSnippets []string
	const fnCount = 24
	for i := range fnCount {
		fmt.Fprintf(&srcBuilder,
			"function f%d(p = () => { let y = 1; y = 2 + arguments.length; return y; }) { let x = 1; x = 2; return x + p(); }\n",
			i,
		)
		// Outer function f_i contains an inner arrow function in its default parameter,
		// so outer f_i is NOT a leaf candidate; only the inner arrow function is optimized.
		wantSnippets = append(wantSnippets, fmt.Sprintf("function f%d(", i))
	}

	compiler := NewCompiler()
	opts := NewCompilerOptions()
	opts.SetLanguage(CompilerOptions_LanguageMode_ECMASCRIPT_NEXT)
	opts.SetNumParallelThreads(4)
	compiler.Init(
		[]*SourceFile{SourceFile_FromCode("externs.js", "")},
		[]*SourceFile{SourceFile_FromCode("input.js", srcBuilder.String())},
		opts,
	)
	compiler.ParseInputs()
	if errs := compiler.GetErrors(); len(errs) != 0 {
		t.Fatalf("unexpected parse errors: %v", errs)
	}

	externs := compiler.GetExternsRoot()
	jsRoot := compiler.GetJsRoot()
	NewNormalize(compiler, false).Process(externs, jsRoot)
	NewDeadAssignmentsElimination(compiler).Process(externs, jsRoot)

	got := compiler.ToSource()
	for _, snippet := range wantSnippets {
		if !strings.Contains(got, snippet) {
			t.Fatalf("missing %q in output:\n%s", snippet, got)
		}
	}
	// Every inner arrow function had `let y = 1; y = 2 + arguments.length;`:
	// after Normalize (`let y$jscomp$... = 1; y$jscomp$... = 2 + ...`), DAE must
	// drop the dead `= 1` initializer into an expression statement `1;`.
	if gotCount := strings.Count(got, ";1;y"); gotCount != fnCount {
		t.Fatalf("expected %d eliminated inner arrow assignments (`;1;y`), got %d in:\n%s", fnCount, gotCount, got)
	}
	// Outer functions have inner functions in PARAM_LIST, so `let x$... = 1` must remain untouched.
	if gotOuterCount := strings.Count(got, "=1;x"); gotOuterCount != fnCount {
		t.Fatalf("expected %d preserved outer initializers (`=1;x`), got %d in:\n%s", fnCount, gotOuterCount, got)
	}
}

func TestRound96_ParallelPropUseIndexPrewarm(t *testing.T) {
	var srcBuilder strings.Builder
	const fnCount = 80 // > piParallelMin (64) so round 2 runs in parallel
	for i := range fnCount {
		fmt.Fprintf(&srcBuilder,
			"function f%d(obj) { obj.p%d = new Ctor(); return obj.q%d(); }\n",
			i, i, i,
		)
	}

	compiler := NewCompiler()
	opts := NewCompilerOptions()
	opts.SetLanguage(CompilerOptions_LanguageMode_ECMASCRIPT_NEXT)
	opts.SetNumParallelThreads(4)
	compiler.Init(
		[]*SourceFile{SourceFile_FromCode("externs.js", "/** @constructor */ function Ctor() {}")},
		[]*SourceFile{SourceFile_FromCode("input.js", srcBuilder.String())},
		opts,
	)
	compiler.ParseInputs()
	if errs := compiler.GetErrors(); len(errs) != 0 {
		t.Fatalf("unexpected parse errors: %v", errs)
	}

	ix := newPropUseIndex(compiler, compiler.GetExternsRoot(), compiler.GetJsRoot())
	ix.build()
	if ix.live < fnCount {
		t.Fatalf("expected at least %d indexed scopes, got %d", fnCount, ix.live)
	}
}

func TestRound96_DAEDefaultParamOnlyDeadAssignment(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return newRound96DAETestCase() })
}

type round96DAETestCase struct {
	CompilerTestCaseBase
}

func newRound96DAETestCase() *round96DAETestCase {
	c := &round96DAETestCase{}
	c.InitCompilerTestCaseBaseWithExterns(c, "var extern;")
	return c
}

func (c *round96DAETestCase) JUnitBefore() {
	c.CompilerTestCaseBase.JUnitBefore()
	c.EnableNormalize()
}

func (c *round96DAETestCase) GetProcessor(compiler *Compiler) CompilerPass {
	return CompilerPassFunc(func(externs rhino.Node, js rhino.Node) {
		NewDeadAssignmentsElimination(compiler).Process(externs, js)
	})
}

func (c *round96DAETestCase) TestDefaultParamDeadAssignments() {
	// Dead assignment in default parameter killed by unconditional body assignment.
	c.Test(
		"function f(a, b = (a = 1, 2)) { a = 3; return a + b; }",
		"function f(a, b = (1, 2)) { a = 3; return a + b; }",
	)
	// Function whose ONLY assignment in the entire function is in PARAM_LIST.
	c.Test(
		"function g(a, b = (a = 1, 2)) { return b; }",
		"function g(a, b = (1, 2)) { return b; }",
	)
	// Assignment in default parameter that IS live into body must be kept.
	c.TestSame("function h1(a, b = (a = 1, 2)) { return a + b; }")
	// Closure in default parameter reading outer parameter must prevent eliminating outer assignment.
	c.TestSame("function h3(a = 1, p = () => { return a; }) { a = 2; return p(); }")
	c.Test(
		"function h4(a = 1, p = () => a) { a = 2; return p(); }",
		"function h4(a = 1, p = () => { return a; }) { a = 2; return p(); }",
	)
}

func TestRound96_NodeUtilLValueNilAndRootGuards(t *testing.T) {
	if got := NodeUtil_GetBestLValue(rhino.NilNode); !got.IsNil() {
		t.Fatalf("NodeUtil_GetBestLValue(NilNode) = %v, want NilNode", got)
	}
	if got := NodeUtil_GetRValueOfLValue(rhino.NilNode); !got.IsNil() {
		t.Fatalf("NodeUtil_GetRValueOfLValue(NilNode) = %v, want NilNode", got)
	}

	compiler := NewCompiler()
	script := compiler.ParseTestCode("var x = 1; function fn() {} class Cls {}")
	// Detach declarations so their parent is NilNode (root node case).
	for c := script.GetFirstChild(); !c.IsNil(); {
		next := c.GetNext()
		c.Detach()
		if got := NodeUtil_GetRValueOfLValue(c); !got.IsNil() {
			t.Fatalf("NodeUtil_GetRValueOfLValue(root %v) = %v, want NilNode", c.GetToken(), got)
		}
		c = next
	}
	if got := NodeUtil_GetRValueOfLValue(script); !got.IsNil() {
		t.Fatalf("NodeUtil_GetRValueOfLValue(script) = %v, want NilNode", got)
	}
	// Upstream Java NodeUtil.getBestLValue throws NullPointerException on parentless
	// SCRIPT nodes (verified by TestNodeUtilSweepGolden).
	func() {
		defer func() {
			r := recover()
			if r == nil {
				t.Fatalf("expected NullPointerException from NodeUtil_GetBestLValue(script)")
			}
			if _, ok := jutil.AsThrowable(r).(jutil.INullPointerException); !ok {
				t.Fatalf("expected INullPointerException, got %T (%v)", r, r)
			}
		}()
		_ = NodeUtil_GetBestLValue(script)
	}()
}

func TestRound96_ScopedAliasesNonFunctionArgNilEnclosingFunction(t *testing.T) {
	compiler := NewCompiler()
	opts := NewCompilerOptions()
	compiler.Init(
		[]*SourceFile{SourceFile_FromCode("externs.js", "var goog = {}; goog.scope = function(fn) {};")},
		[]*SourceFile{SourceFile_FromCode("input.js", "goog.scope([function() { { function inner() {} } }, 1]);")},
		opts,
	)
	compiler.ParseInputs()
	NewScopedAliases(compiler, ScopedAliasesOptions{}).Process(compiler.GetExternsRoot(), compiler.GetJsRoot())

	errs := compiler.GetErrors()
	if len(errs) == 0 {
		t.Fatalf("expected GOOG_SCOPE_HAS_BAD_PARAMETERS error, got none")
	}
	found := false
	for _, err := range errs {
		if err.GetType() == ScopedAliases_GOOG_SCOPE_HAS_BAD_PARAMETERS {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected GOOG_SCOPE_HAS_BAD_PARAMETERS in %v", errs)
	}
}

func TestRound96_PeepholeMinimizeConditionsDanglingElse(t *testing.T) {
	// Direct unit check on PeepholeMinimizeConditions_ConsumesDanglingElse for
	// FOR_OF, FOR_AWAIT_OF, and LABEL.
	checkConsumes := func(code string, want bool) {
		t.Helper()
		compiler := NewCompiler()
		opts := NewCompilerOptions()
		opts.SetLanguage(CompilerOptions_LanguageMode_ECMASCRIPT_NEXT)
		compiler.InitOptions(opts)
		script := compiler.ParseTestCode("async function wrap() { " + code + " }")
		if errs := compiler.GetErrors(); len(errs) != 0 {
			t.Fatalf("unexpected parse errors for %q: %v", code, errs)
		}
		fn := ast.MustFunction(script.GetFirstChild())
		stmt := fn.Body().MustBlock().MustFirstStatement()
		if got := PeepholeMinimizeConditions_ConsumesDanglingElse(stmt); got != want {
			t.Fatalf("ConsumesDanglingElse(%q) = %v, want %v", code, got, want)
		}
	}

	checkConsumes("for (const x of y) if (b) c();", true)
	checkConsumes("for await (const x of y) if (b) c();", true)
	checkConsumes("lbl: if (b) c();", true)
	checkConsumes("for (const x of y) if (b) c(); else d();", false)
	checkConsumes("lbl: if (b) c(); else d();", false)

	// End-to-end PeepholeMinimizeConditions + CodeGenerator checks:
	// 1. Braces must be preserved so the outer `else d()` is not captured by the inner `if (b) c()`.
	// 2. `if (!a) ... else <dangling-else consumer>` must NOT swap branches.
	runMinimize := func(src string) string {
		t.Helper()
		compiler := NewCompiler()
		opts := NewCompilerOptions()
		opts.SetLanguage(CompilerOptions_LanguageMode_ECMASCRIPT_NEXT)
		compiler.Init(
			[]*SourceFile{SourceFile_FromCode("externs.js", "")},
			[]*SourceFile{SourceFile_FromCode("input.js", src)},
			opts,
		)
		compiler.ParseInputs()
		if errs := compiler.GetErrors(); len(errs) != 0 {
			t.Fatalf("unexpected parse errors for %q: %v", src, errs)
		}
		pass := NewPeepholeOptimizationsPass(compiler, "minimize", NewPeepholeMinimizeConditions(true))
		pass.Process(compiler.GetExternsRoot(), compiler.GetJsRoot())
		out := compiler.ToSource()

		// Re-parse the printed output and verify the outer IF still has its ELSE branch!
		compiler2 := NewCompiler()
		compiler2.InitOptions(opts)
		reparsed := compiler2.ParseTestCode(out)
		if errs := compiler2.GetErrors(); len(errs) != 0 {
			t.Fatalf("reparse errors for %q: %v", out, errs)
		}
		topStmt := reparsed.GetFirstChild()
		if fn, ok := ast.AsFunction(topStmt); ok {
			topStmt = fn.Body().MustBlock().MustFirstStatement().Node()
		}
		outerIf, ok := ast.AsIf(topStmt)
		if !ok {
			t.Fatalf("expected outer IF after minimizing %q, got %s (%s)", src, topStmt.GetToken(), out)
		}
		if _, hasElse := outerIf.Else(); !hasElse {
			t.Fatalf("outer IF lost its ELSE branch after minimizing %q -> %q", src, out)
		}
		return out
	}

	// When inner `if` cannot be folded to `&&` (e.g. `break` / `return`), braces
	// around the dangling-else consumer MUST be preserved!
	if got := runMinimize("if (a) { for (const x of y) if (b) break; } else d();"); !strings.Contains(got, "{") {
		t.Fatalf("expected braces preserved around dangling-else in for-of, got %q", got)
	}
	if got := runMinimize("if (a) { lbl: if (b) break lbl; } else d();"); !strings.Contains(got, "{") {
		t.Fatalf("expected braces preserved around dangling-else in label, got %q", got)
	}
	if got := runMinimize("async function f() { if (a) { for await (const x of y) if (b) break; } else d(); }"); !strings.Contains(got, "}else d()") {
		t.Fatalf("expected braces preserved before else d() in for-await-of, got %q", got)
	}
	// Also verify expression-bodied inner `if (b) c();` keeps the outer `else d()` intact after re-parse.
	_ = runMinimize("if (a) { for (const x of y) if (b) c(); } else d();")
	_ = runMinimize("if (a) { lbl: if (b) c(); } else d();")

	// Verify `!x` condition is NOT swapped when the else branch consumes a dangling else via FOR_OF, FOR_AWAIT_OF, or LABEL.
	RunJUnit4(t, func() JUnitTestCase { return newRound96PMCTestCase() })
}

type round96PMCTestCase struct {
	CompilerTestCaseBase
}

func newRound96PMCTestCase() *round96PMCTestCase {
	c := &round96PMCTestCase{}
	c.InitCompilerTestCaseBaseWithExterns(c, CompilerTestCase_DEFAULT_EXTERNS)
	return c
}

func (c *round96PMCTestCase) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.DisableTypeCheck()
}

func (c *round96PMCTestCase) GetProcessor(compiler *Compiler) CompilerPass {
	peepholePass := NewPeepholeOptimizationsPass(compiler, c.GetName(), NewPeepholeMinimizeConditions(true))
	peepholePass.SetRetraverseOnChange(false)
	return peepholePass
}

func (c *round96PMCTestCase) TestNoSwapWithForOfAndLabelDanglingElse() {
	c.TestSame("if (!x) { foo(); bar(); } else for (const a of b) if (y) break;")
	c.TestSame("async function g() { if (!x) { foo(); bar(); } else for await (const a of b) if (y) break; }")
	c.TestSame("if (!x) { foo(); bar(); } else lbl: if (y) break lbl;")
}
