package jscomp

import (
	"bytes"
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

func runCheckControlFlowCompileWithThreads(
	t *testing.T,
	externs []*SourceFile,
	inputs []*SourceFile,
	numThreads int,
) (errs []string, warns []string, report string) {
	t.Helper()

	opts := NewCompilerOptions()
	opts.SetChecksOnly(true)
	opts.SetCheckTypes(true)
	opts.SetClosurePass(true)
	opts.SetNumParallelThreads(numThreads)
	opts.SetWarningLevel(DiagnosticGroups_CHECK_USELESS_CODE, CheckLevel_WARNING)
	opts.SetWarningLevel(DiagnosticGroups_MISSING_RETURN, CheckLevel_WARNING)

	var reportBuf bytes.Buffer
	comp := NewCompilerWithOutStream(jutil.NewPrintStream(&reportBuf))
	comp.DisableThreads()
	comp.Compile(externs, inputs, opts)

	for _, e := range comp.GetErrors() {
		errs = append(errs, e.String())
	}
	for _, w := range comp.GetWarnings() {
		warns = append(warns, w.String())
	}
	return errs, warns, reportBuf.String()
}

func runCheckUnreachableDirectWithThreads(
	t *testing.T,
	inputs []*SourceFile,
	numThreads int,
) (errs []string, warns []string, report string) {
	t.Helper()

	opts := NewCompilerOptions()
	opts.SetChecksOnly(true)
	opts.SetNumParallelThreads(numThreads)
	opts.SetWarningLevel(DiagnosticGroups_CHECK_USELESS_CODE, CheckLevel_WARNING)
	opts.SetWarningLevel(DiagnosticGroups_MISSING_RETURN, CheckLevel_WARNING)

	var reportBuf bytes.Buffer
	comp := NewCompilerWithOutStream(jutil.NewPrintStream(&reportBuf))
	comp.DisableThreads()
	comp.Init([]*SourceFile{}, inputs, opts)
	root := comp.ParseInputs()
	if root.IsNil() || comp.HasErrors() {
		for _, e := range comp.GetErrors() {
			errs = append(errs, e.String())
		}
		return errs, nil, reportBuf.String()
	}

	externsRoot := root.GetFirstChild()
	jsRoot := root.GetSecondChild()
	defaultPass_checkControlFlow.Create(comp).Process(externsRoot, jsRoot)

	for _, e := range comp.GetErrors() {
		errs = append(errs, e.String())
	}
	for _, w := range comp.GetWarnings() {
		warns = append(warns, w.String())
	}
	comp.GetErrorManager().GenerateReport()
	return errs, warns, reportBuf.String()
}

func TestCheckUnreachableCodeParallelParity(t *testing.T) {
	t.Run("MultiScriptUnreachableMissingReturnAndSuppressions", func(t *testing.T) {
		externs := []*SourceFile{
			SourceFile_FromCode("externs.js", strings.Join([]string{
				"/** @fileoverview @externs */",
				"/** @constructor @template T */ function Promise(resolver) {}",
				"/** @interface @template Y, R, N */ function Generator() {}",
				"/** @constructor @param {string=} msg */ function Error(msg) {}",
				"/** @param {...*} args */ function sideEffect(args) {}",
			}, "\n")),
		}

		var inputs []*SourceFile
		for i := range 24 {
			name := fmt.Sprintf("script_%02d.js", i)
			var lines []string
			switch i % 6 {
			case 0:
				// Straight-line script with straight-line and branching functions.
				lines = append(lines,
					fmt.Sprintf("var straight_%d = %d;", i, i),
					fmt.Sprintf("/** @param {boolean} b @return {number} */ function missingRet_%d(b) {", i),
					"  if (b) { return 1; }",
					"}",
					fmt.Sprintf("function unreachAfterRet_%d() {", i),
					"  return 42;",
					"  sideEffect('dead after return');",
					"  sideEffect('second dead stmt');",
					"}",
				)
			case 1:
				// Suppressed vs unsuppressed unreachable code across scripts.
				lines = append(lines,
					fmt.Sprintf("/** @suppress {uselessCode} */ function suppressedFn_%d() {", i),
					"  return 1;",
					"  sideEffect('suppressed dead');",
					"}",
					fmt.Sprintf("function unsuppressedFn_%d(x) {", i),
					"  if (false) { sideEffect('dead if false'); }",
					"  while (false) { sideEffect('dead while false'); }",
					"  switch (x) { case 1: throw new Error('boom'); sideEffect('dead case'); }",
					"}",
				)
			case 2:
				// Top-level unreachable statement followed by reachable code.
				lines = append(lines,
					"if (false) { sideEffect('top-level dead branch'); }",
					fmt.Sprintf("try { throw new Error('caught'); sideEffect('dead in try'); } catch (e) { sideEffect(e); }"),
					fmt.Sprintf("var topAlive_%d = 1;", i),
				)
			case 3:
				// Classes, methods, and class static blocks with unreachable code + missing return.
				lines = append(lines,
					fmt.Sprintf("class Box_%d {", i),
					"  static {",
					"    if (false) { sideEffect('dead static block'); }",
					"  }",
					"  /** @param {boolean} b @return {string} */",
					"  compute(b) {",
					"    if (b) { return 'ok'; throw new Error('dead in method'); }",
					"  }",
					"}",
				)
			case 4:
				// Async + generator functions sharing CFGs between CheckUnreachableCode and CheckMissingReturn.
				lines = append(lines,
					fmt.Sprintf("/** @param {boolean} b @return {!Promise<number>} */ async function asyncFn_%d(b) {", i),
					"  if (b) { return 7; sideEffect('dead async'); }",
					"}",
					fmt.Sprintf("/** @param {boolean} b @return {!Generator<number, string, *>} */ function* genFn_%d(b) {", i),
					"  yield 1;",
					"  if (b) { return 'done'; sideEffect('dead gen'); }",
					"}",
				)
			case 5:
				// ES6 module script with unreachable code.
				lines = append(lines,
					fmt.Sprintf("export const modVal_%d = %d;", i, i),
					fmt.Sprintf("export function modUnreach_%d() {", i),
					"  throw new Error('stop');",
					"  sideEffect('dead in es6 module fn');",
					"}",
					"if (false) { sideEffect('dead in es6 module top-level'); }",
				)
			}
			inputs = append(inputs, SourceFile_FromCode(name, strings.Join(lines, "\n")))
		}

		errs1, warns1, report1 := runCheckControlFlowCompileWithThreads(t, externs, inputs, 1)
		errs4, warns4, report4 := runCheckControlFlowCompileWithThreads(t, externs, inputs, 4)

		if len(errs1) != 0 || len(errs4) != 0 {
			t.Fatalf("expected 0 errors, got threads=1: %v, threads=4: %v", errs1, errs4)
		}
		if len(warns1) == 0 {
			t.Fatalf("expected non-empty warnings")
		}
		if !slices.Equal(warns1, warns4) {
			t.Fatalf("warnings mismatch:\nthreads=1 (%d):\n%s\nthreads=4 (%d):\n%s",
				len(warns1), strings.Join(warns1, "\n"),
				len(warns4), strings.Join(warns4, "\n"))
		}
		if report1 != report4 {
			t.Fatalf("report mismatch:\nthreads=1:\n%s\nthreads=4:\n%s", report1, report4)
		}

		hasUnreach := false
		hasMissingRet := false
		for _, w := range warns4 {
			if strings.Contains(w, "JSC_UNREACHABLE_CODE") {
				hasUnreach = true
			}
			if strings.Contains(w, "JSC_MISSING_RETURN_STATEMENT") {
				hasMissingRet = true
			}
		}
		if !hasUnreach || !hasMissingRet {
			t.Fatalf("expected both JSC_UNREACHABLE_CODE (%v) and JSC_MISSING_RETURN_STATEMENT (%v)", hasUnreach, hasMissingRet)
		}
	})

	t.Run("CrossScriptTopLevelUnconditionalThrow", func(t *testing.T) {
		inputs := []*SourceFile{
			SourceFile_FromCode("s0.js", "var a = 1;"),
			SourceFile_FromCode("s1.js", "var beforeThrow = 2;\nthrow new Error('stop all following scripts');"),
			SourceFile_FromCode("s2.js", "var unreachScript = 3;\nif (false) { var nestedDead = 4; }"),
			SourceFile_FromCode("s3.js", "var reachableAfterRecompute = 5;"),
		}

		errs1, warns1, report1 := runCheckUnreachableDirectWithThreads(t, inputs, 1)
		errs4, warns4, report4 := runCheckUnreachableDirectWithThreads(t, inputs, 4)

		if !slices.Equal(errs1, errs4) {
			t.Fatalf("errors mismatch:\nthreads=1:\n%v\nthreads=4:\n%v", errs1, errs4)
		}
		if !slices.Equal(warns1, warns4) {
			t.Fatalf("warnings mismatch:\nthreads=1:\n%s\nthreads=4:\n%s",
				strings.Join(warns1, "\n"), strings.Join(warns4, "\n"))
		}
		if report1 != report4 {
			t.Fatalf("report mismatch:\nthreads=1:\n%s\nthreads=4:\n%s", report1, report4)
		}
		if len(warns4) == 0 {
			t.Fatalf("expected unreachable-script warning after top-level unconditional throw")
		}
	})

	t.Run("CrossScriptModuleTopLevelThrow", func(t *testing.T) {
		inputs := []*SourceFile{
			SourceFile_FromCode("m0.js", "export const x = 1;\nthrow new Error('module halt');\nconst deadAfterThrow = 2;"),
			SourceFile_FromCode("m1.js", "export const y = 3;"),
			SourceFile_FromCode("m2.js", "export const z = 4;"),
		}

		errs1, warns1, report1 := runCheckUnreachableDirectWithThreads(t, inputs, 1)
		errs4, warns4, report4 := runCheckUnreachableDirectWithThreads(t, inputs, 4)

		if !slices.Equal(errs1, errs4) {
			t.Fatalf("errors mismatch:\nthreads=1:\n%v\nthreads=4:\n%v", errs1, errs4)
		}
		if !slices.Equal(warns1, warns4) {
			t.Fatalf("warnings mismatch:\nthreads=1:\n%s\nthreads=4:\n%s",
				strings.Join(warns1, "\n"), strings.Join(warns4, "\n"))
		}
		if report1 != report4 {
			t.Fatalf("report mismatch:\nthreads=1:\n%s\nthreads=4:\n%s", report1, report4)
		}
	})
}
