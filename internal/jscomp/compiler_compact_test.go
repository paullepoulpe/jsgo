package jscomp

import (
	"io"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parsing"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Go-only: AST compaction (compiler_compact.go) with holders that only some
// options create.

// compileForCompactTest compiles input in ADVANCED mode with type checking,
// the tracer on and JSDoc documentation parsed, so that the performance
// tracker holds the roots across the compaction after typesToColors (it
// crashed with --tracer_mode before the tracker's roots were forwarded), and
// returns the output and the number of compactions. Documentation markers
// hold nodes only until typesToColors simplifies the JSDoc, before any
// compaction.
func compileForCompactTest(t *testing.T, input string, threads int) (string, int) {
	t.Helper()
	options := NewCompilerOptions()
	CompilationLevel_ADVANCED_OPTIMIZATIONS.SetOptionsForCompilationLevel(options)
	options.SetCheckTypes(true)
	options.SetNumParallelThreads(threads)
	options.SetTracerMode(CompilerOptions_TracerMode_ALL)
	options.SetParseJsDocDocumentation(parsing.Config_JsDocParsing_INCLUDE_ALL_COMMENTS)
	c := NewCompilerWithOutStream(jutil.NewPrintStream(io.Discard))
	result := c.Compile(
		[]*SourceFile{SourceFile_FromCode("externs.js", propIndexLoopExterns)},
		[]*SourceFile{
			SourceFile_FromCode("a.js", "/** @param {number} x the value\n @return {number} */ function twice(x) { return x * 2; }\nalert(twice);\n"),
			SourceFile_FromCode("in.js", input),
		},
		options)
	if !result.Success {
		t.Fatalf("compilation failed: %v", result.Errors)
	}
	return c.ToSource(), c.astCompact.count
}

// TestCompactionWithTracerAndDocumentation: with the tracer on and JSDoc
// documentation parsed, the compaction runs and the output is the same with
// and without it.
func TestCompactionWithTracerAndDocumentation(t *testing.T) {
	input := propIndexManyScopes() + "alert(twice(3));\n"
	for _, threads := range []int{1, 4} {
		saved := astCompactEnv
		astCompactEnv = "0"
		want, n := compileForCompactTest(t, input, threads)
		astCompactEnv = saved
		if n != 0 {
			t.Fatalf("%d compactions with compaction off", n)
		}
		got, m := compileForCompactTest(t, input, threads)
		if m < 1 {
			t.Fatalf("%d threads: no compaction after typesToColors", threads)
		}
		if got != want {
			t.Errorf("%d threads: output with compaction:\n%s\nwithout:\n%s", threads, got, want)
		}
	}
}

func TestCompactionClearsPropIndexUnderDevModeEveryPassAndPropIndexVerify(t *testing.T) {
	savedEnv, savedForce, savedVerify := propIndexEnv, propIndexForce, propIndexVerify
	propIndexEnv, propIndexForce, propIndexVerify = "1", true, true
	defer func() {
		propIndexEnv, propIndexForce, propIndexVerify = savedEnv, savedForce, savedVerify
	}()

	input := propIndexManyScopes() + "alert(twice(3));\n"
	for _, threads := range []int{1, 4} {
		options := NewCompilerOptions()
		CompilationLevel_ADVANCED_OPTIMIZATIONS.SetOptionsForCompilationLevel(options)
		options.SetCheckTypes(true)
		options.SetDevMode(CompilerOptions_DevMode_EVERY_PASS)
		options.SetNumParallelThreads(threads)
		c := NewCompilerWithOutStream(jutil.NewPrintStream(io.Discard))
		result := c.Compile(
			[]*SourceFile{SourceFile_FromCode("externs.js", propIndexLoopExterns)},
			[]*SourceFile{
				SourceFile_FromCode("a.js", "function twice(x) { return x * 2; }\nalert(twice);\n"),
				SourceFile_FromCode("in.js", input),
			},
			options)
		if !result.Success {
			t.Fatalf("threads=%d: compilation failed under DevMode_EVERY_PASS + propIndexVerify: %v", threads, result.Errors)
		}
		if c.astCompact.count < 1 {
			t.Fatalf("threads=%d: expected >= 1 AST compaction, got %d", threads, c.astCompact.count)
		}
	}
}

func TestCompactionResetsGlobalNamespaceVarNameHash(t *testing.T) {
	rhino.ResetNodeStore()
	c := NewCompilerWithOutStream(jutil.NewPrintStream(io.Discard))
	options := NewCompilerOptions()
	// Allocate dead nodes BEFORE parsing so SlideWithThreads shifts live scope root IDs down!
	for range 512 {
		_ = rhino.IR_Name("deadBeforeParse")
	}
	c.Init(
		[]*SourceFile{SourceFile_FromCode("externs.js", "function alert(x) {}\n")},
		[]*SourceFile{SourceFile_FromCode("in.js", "var ns = {a: 1}; function f(localParam) { var localVar = ns.a + localParam; return localVar; } alert(f(2));\n")},
		options,
	)
	c.ParseForCompilation()
	if c.HasErrors() {
		t.Fatalf("parse errors: %v", c.GetErrors())
	}

	gn := NewGlobalNamespaceWithRoots(c, c.externsRoot, c.jsRoot)
	_ = gn.GetSlot("ns")
	c.transpilationNamespace = gn

	// Prime v.nameHash caches across all scopes referenced by GlobalNamespace BEFORE compaction.
	var vars []*Var
	seenVars := make(map[*Var]struct{})
	for i := uint32(0); i < gn.refCount; i++ {
		ref := &gn.refPages[i>>gnRefPageShift][i&gnRefPageMask]
		for s := ref.scope; s != nil; s = s.GetParent() {
			for _, v := range s.varList {
				if v != nil {
					if _, dup := seenVars[v]; !dup {
						seenVars[v] = struct{}{}
						_ = v.HashCode()
						vars = append(vars, v)
					}
				}
			}
		}
	}
	if len(vars) == 0 {
		t.Fatal("expected non-empty Var set in GlobalNamespace scopes")
	}

	externs, root := c.externsRoot, c.jsRoot
	c.compactAST("testGlobalNamespaceVarHash", true, &externs, &root, nil)
	if c.astCompact.count != 1 {
		t.Fatalf("expected 1 forced compaction, got %d", c.astCompact.count)
	}

	for _, v := range vars {
		gotHash := v.HashCode()
		wantHash := ScopedName_Of(v.GetName(), v.GetScopeRoot()).HashCode()
		if gotHash != wantHash {
			t.Fatalf("Var %q (scopeRoot=%v): HashCode()=%d after compaction, want ScopedName.HashCode()=%d (stale v.nameHash was not reset)",
				v.GetName(), v.GetScopeRoot(), gotHash, wantHash)
		}
	}
}
