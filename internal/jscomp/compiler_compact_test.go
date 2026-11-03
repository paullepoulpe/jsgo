package jscomp

import (
	"io"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parsing"
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
