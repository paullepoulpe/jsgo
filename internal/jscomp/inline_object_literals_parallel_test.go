package jscomp

import (
	"fmt"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

func runIOLPassOnce(t *testing.T, srcs []string, threads int, passes int) (out string, nextID string, changeStamp int32) {
	t.Helper()
	compiler := NewCompiler()
	opts := NewCompilerOptions()
	opts.SetPrettyPrint(true)
	opts.SetNumParallelThreads(threads)
	if opts.isFastMode() {
		t.Fatal("expected default (!fastMode) options")
	}

	externs := []*SourceFile{SourceFile_FromCode("externs.js", "function use(x, y) {}\nfunction escape(x) {}\nvar ns = {};\n")}
	inputs := make([]*SourceFile, len(srcs))
	for i, code := range srcs {
		inputs[i] = SourceFile_FromCode(fmt.Sprintf("s%02d.js", i), code)
	}
	compiler.Init(externs, inputs, opts)
	compiler.ParseInputs()
	if compiler.HasErrors() {
		t.Fatalf("unexpected parse errors: %v", compiler.GetErrors())
	}
	NewNormalize(compiler, false).Process(compiler.GetExternsRoot(), compiler.GetJsRoot())
	idSupplier := compiler.GetUniqueNameIdSupplier()
	for i := 0; i < passes; i++ {
		NewInlineObjectLiterals(compiler, idSupplier).Process(compiler.GetExternsRoot(), compiler.GetJsRoot())
	}
	return compiler.ToSource(), idSupplier(), compiler.GetChangeTracker().GetChangeStamp()
}

func TestInlineObjectLiteralsParallel_DefaultModeByteIdentical(t *testing.T) {
	makeEscapedScript := func(i int) string {
		return fmt.Sprintf(
			"function hoistedEscaped_%d() {\n"+
				"  var esc = {a: %d, b: %d};\n"+
				"  escape(esc);\n"+
				"  return esc.a;\n"+
				"}\n"+
				"use(hoistedEscaped_%d(), 0);\n",
			i, i, i+1, i,
		)
	}

	t.Run("zero_actionable_among_many_hoisted", func(t *testing.T) {
		srcs := make([]string, 12)
		for i := range srcs {
			srcs[i] = makeEscapedScript(i)
		}
		out1, id1, stamp1 := runIOLPassOnce(t, srcs, 1, 2)
		out4, id4, stamp4 := runIOLPassOnce(t, srcs, 4, 2)
		if out1 != out4 || id1 != id4 || stamp1 != stamp4 {
			t.Fatalf("mismatch (1 vs 4 threads):\nid: %s vs %s\nstamp: %d vs %d\nout1:\n%s\nout4:\n%s",
				id1, id4, stamp1, stamp4, out1, out4)
		}
		if id1 != "0" {
			t.Fatalf("expected 0 consumed unique IDs when nothing splits, got nextID=%s", id1)
		}
	})

	t.Run("single_actionable_in_global_hoisted", func(t *testing.T) {
		srcs := make([]string, 12)
		for i := range srcs {
			srcs[i] = makeEscapedScript(i)
		}
		srcs[7] += "function hoistedActionable() { var o = {x: 10, y: 20}; use(o.x, o.y); }\nhoistedActionable();\n"

		out1, id1, stamp1 := runIOLPassOnce(t, srcs, 1, 2)
		out4, id4, stamp4 := runIOLPassOnce(t, srcs, 4, 2)
		if out1 != out4 || id1 != id4 || stamp1 != stamp4 {
			t.Fatalf("mismatch (1 vs 4 threads):\nid: %s vs %s\nstamp: %d vs %d\nout1:\n%s\nout4:\n%s",
				id1, id4, stamp1, stamp4, out1, out4)
		}
		if !strings.Contains(out4, "JSCompiler_object_inline_x_0") || !strings.Contains(out4, "JSCompiler_object_inline_y_1") {
			t.Fatalf("expected inline temp vars 0 and 1 in output:\n%s", out4)
		}
	})

	t.Run("distinct_scripts_local_hoisted_only", func(t *testing.T) {
		srcs := make([]string, 12)
		for i := range srcs {
			srcs[i] = makeEscapedScript(i)
		}
		// Local hoisted functions inside IIFEs in distinct scripts: cannot be jumped into from earlier scripts.
		srcs[2] += "(function() { localH2(); function localH2() { var p = {k: 1}; use(p.k, 0); } })();\n"
		srcs[5] += "(function() { localH5(); function localH5() { var q = {m: 2}; use(q.m, 0); } })();\n"
		srcs[10] += "(function() { localH10(); function localH10() { var r = {n: 3}; use(r.n, 0); } })();\n"

		out1, id1, stamp1 := runIOLPassOnce(t, srcs, 1, 2)
		out4, id4, stamp4 := runIOLPassOnce(t, srcs, 4, 2)
		if out1 != out4 || id1 != id4 || stamp1 != stamp4 {
			t.Fatalf("mismatch (1 vs 4 threads):\nid: %s vs %s\nstamp: %d vs %d\nout1:\n%s\nout4:\n%s",
				id1, id4, stamp1, stamp4, out1, out4)
		}
		if !strings.Contains(out4, "JSCompiler_object_inline_k_0") ||
			!strings.Contains(out4, "JSCompiler_object_inline_m_1") ||
			!strings.Contains(out4, "JSCompiler_object_inline_n_2") {
			t.Fatalf("expected exact sequential IDs 0, 1, 2:\n%s", out4)
		}
	})

	t.Run("cross_script_global_hoisted_reference_order_jump", func(t *testing.T) {
		srcs := make([]string, 12)
		for i := range srcs {
			srcs[i] = makeEscapedScript(i)
		}
		// Script 1 calls globalHoistedIn9() BEFORE defining its own IIFE with an actionable literal!
		// ReferenceCollector jumps into globalHoistedIn9 first (ID 0), then script 1's IIFE gets ID 1.
		srcs[1] = "globalHoistedIn9();\n(function() { var second = {beta: 2}; use(second.beta, 0); })();\n" + srcs[1]
		srcs[9] += "function globalHoistedIn9() { var first = {alpha: 1}; use(first.alpha, 0); }\n"

		out1, id1, stamp1 := runIOLPassOnce(t, srcs, 1, 2)
		out4, id4, stamp4 := runIOLPassOnce(t, srcs, 4, 2)
		if out1 != out4 || id1 != id4 || stamp1 != stamp4 {
			t.Fatalf("mismatch (1 vs 4 threads):\nid: %s vs %s\nstamp: %d vs %d\nout1:\n%s\nout4:\n%s",
				id1, id4, stamp1, stamp4, out1, out4)
		}
		if !strings.Contains(out4, "JSCompiler_object_inline_alpha_0") ||
			!strings.Contains(out4, "JSCompiler_object_inline_beta_1") {
			t.Fatalf("expected cross-script hoisted jump order (alpha_0 then beta_1):\n%s", out4)
		}
	})

	t.Run("same_script_local_hoisted_reference_order_jump", func(t *testing.T) {
		srcs := make([]string, 12)
		for i := range srcs {
			srcs[i] = makeEscapedScript(i)
		}
		srcs[4] += "(function() {\n" +
			"  laterHoisted();\n" +
			"  var f = function() { var bObj = {secondKey: 2}; use(bObj.secondKey, 0); };\n" +
			"  f();\n" +
			"  function laterHoisted() { var aObj = {firstKey: 1}; use(aObj.firstKey, 0); }\n" +
			"})();\n"

		out1, id1, stamp1 := runIOLPassOnce(t, srcs, 1, 2)
		out4, id4, stamp4 := runIOLPassOnce(t, srcs, 4, 2)
		if out1 != out4 || id1 != id4 || stamp1 != stamp4 {
			t.Fatalf("mismatch (1 vs 4 threads):\nid: %s vs %s\nstamp: %d vs %d\nout1:\n%s\nout4:\n%s",
				id1, id4, stamp1, stamp4, out1, out4)
		}
		if !strings.Contains(out4, "JSCompiler_object_inline_firstKey_0") ||
			!strings.Contains(out4, "JSCompiler_object_inline_secondKey_1") {
			t.Fatalf("expected local hoisted jump order (firstKey_0 then secondKey_1):\n%s", out4)
		}
	})
}

func TestInlineObjectLiteralsParallel_IncrementalChangeTracker(t *testing.T) {
	// 1. Nested object literal splits across Run 1 and Run 2, and Run 3 is a 0-change no-op.
	srcs := make([]string, 12)
	for i := range srcs {
		srcs[i] = fmt.Sprintf("(function() { var keep = {a: %d}; escape(keep); })();\n", i)
	}
	srcs[6] = "(function() { var outer = {inner: {leaf: 99}}; use(outer.inner.leaf, 0); })();\n"

	out1, id1, stamp1 := runIOLPassOnce(t, srcs, 1, 3)
	out4, id4, stamp4 := runIOLPassOnce(t, srcs, 4, 3)
	if out1 != out4 || id1 != id4 || stamp1 != stamp4 {
		t.Fatalf("nested incremental mismatch:\nid: %s vs %s\nstamp: %d vs %d\nout1:\n%s\nout4:\n%s",
			id1, id4, stamp1, stamp4, out1, out4)
	}
	if !strings.Contains(out4, "JSCompiler_object_inline_leaf_1") || strings.Contains(out4, "inner:") {
		t.Fatalf("expected nested object literal to fully inline by pass 2:\n%s", out4)
	}

	// 2. Giant outer IIFE where mutating only child function targetFn unlocks both targetFn's local
	// object literal AND the outer IIFE's own object literal while skipping unchanged siblings.
	compiler := NewCompiler()
	opts := NewCompilerOptions()
	opts.SetNumParallelThreads(4)
	var iife strings.Builder
	iife.WriteString("(function(ns) {\n  var outerTable = { entry: 77 };\n")
	for i := 0; i < 16; i++ {
		fmt.Fprintf(&iife, "  function sibling_%d() { var s = { k: %d }; escape(s); return s.k; }\n", i, i)
	}
	iife.WriteString("  var targetFn = function() { var loc = { val: 88 }; escape(loc); escape(outerTable); };\n")
	iife.WriteString("  targetFn();\n})(ns);\n")

	compiler.Init(
		[]*SourceFile{SourceFile_FromCode("externs.js", "function use(x, y) {}\nfunction escape(x) {}\nvar ns = {};\n")},
		[]*SourceFile{SourceFile_FromCode("app.js", iife.String())},
		opts,
	)
	compiler.ParseInputs()
	NewNormalize(compiler, false).Process(compiler.GetExternsRoot(), compiler.GetJsRoot())
	idSupplier := compiler.GetUniqueNameIdSupplier()

	// Run 1: nothing is inlinable because targetFn escapes both `loc` and `outerTable`.
	NewInlineObjectLiterals(compiler, idSupplier).Process(compiler.GetExternsRoot(), compiler.GetJsRoot())
	if strings.Contains(compiler.ToSource(), "JSCompiler_object_inline_") {
		t.Fatalf("expected zero splits on Run 1, got:\n%s", compiler.ToSource())
	}

	// Mutate ONLY targetFn: replace `escape(loc)` -> `use(loc.val, 0)` and `escape(outerTable)` -> `use(outerTable.entry, 0)`,
	// reporting change ONLY to targetFn (NOT the outer IIFE and NOT sibling_0..15).
	var targetFnNode rhino.Node
	var findTargetFn func(n rhino.Node)
	findTargetFn = func(n rhino.Node) {
		if n.IsName() && n.GetString() == "targetFn" && !n.GetFirstChild().IsNil() && n.GetFirstChild().IsFunction() {
			targetFnNode = n.GetFirstChild()
			return
		}
		for c := n.GetFirstChild(); !c.IsNil() && targetFnNode.IsNil(); c = c.GetNext() {
			findTargetFn(c)
		}
	}
	findTargetFn(compiler.GetJsRoot())
	if targetFnNode.IsNil() {
		t.Fatal("failed to locate targetFn")
	}
	replaced := compiler.ParseSyntheticCode("patch",
		"function targetFn() { var loc = { val: 88 }; use(loc.val, 0); use(outerTable.entry, 0); }").GetFirstChild()
	newBody := replaced.GetLastChild().Detach()
	targetFnNode.GetLastChild().ReplaceWith(newBody)
	compiler.ReportChangeToChangeScope(targetFnNode)

	// Run 2 (incremental): must inline BOTH `loc` (in targetFn) AND `outerTable` (in outer IIFE)!
	NewInlineObjectLiterals(compiler, idSupplier).Process(compiler.GetExternsRoot(), compiler.GetJsRoot())
	gotAfterRun2 := compiler.ToSource()
	if !strings.Contains(gotAfterRun2, "JSCompiler_object_inline_val_0") ||
		!strings.Contains(gotAfterRun2, "JSCompiler_object_inline_entry_1") {
		t.Fatalf("expected incremental Run 2 to split both loc and outerTable, got:\n%s", gotAfterRun2)
	}

	// Run 3 (0 changes since Run 2): immediate no-op.
	stampBeforeRun3 := compiler.GetChangeTracker().GetChangeStamp()
	NewInlineObjectLiterals(compiler, idSupplier).Process(compiler.GetExternsRoot(), compiler.GetJsRoot())
	if got := compiler.GetChangeTracker().GetChangeStamp(); got != stampBeforeRun3 {
		t.Fatalf("expected no-op Run 3 to preserve change stamp %d, got %d", stampBeforeRun3, got)
	}
}
