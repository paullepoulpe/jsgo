package jscomp

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

func compileVarRefCheckWithThreads(
	t *testing.T,
	externs []*SourceFile,
	inputs []*SourceFile,
	numThreads int,
) (errs []string, warns []string, report string) {
	t.Helper()

	opts := NewCompilerOptions()
	opts.SetChecksOnly(true)
	opts.SetClosurePass(true)
	opts.SetWarningLevel(DiagnosticGroups_CHECK_VARIABLES, CheckLevel_WARNING)
	opts.SetWarningLevel(DiagnosticGroups_UNUSED_LOCAL_VARIABLE, CheckLevel_WARNING)
	opts.SetNumParallelThreads(numThreads)

	var reportBuf bytes.Buffer
	comp := NewCompilerWithOutStream(jutil.NewPrintStream(&reportBuf))
	comp.DisableThreads()
	res := comp.Compile(externs, inputs, opts)
	for _, e := range res.Errors {
		errs = append(errs, e.String())
	}
	for _, w := range res.Warnings {
		warns = append(warns, w.String())
	}
	return errs, warns, reportBuf.String()
}

func TestVariableReferenceCheckParallel_CrossScriptHoistedAndRedeclarations(t *testing.T) {
	externs := []*SourceFile{
		SourceFile_FromCode("externs_0.js", `/** @externs */
function sink(x) {}
`),
	}

	inputs := []*SourceFile{
		SourceFile_FromCode("script_0.js", `
// Call global hoisted function declared in script_2.js before declaration:
sink(globalHoistedFn(1));

// Early reference to global var declared in script_1.js:
sink(earlyGlobalVar);

var dupGlobalVar = 10;
const globalConst = 5;
`),
		SourceFile_FromCode("script_1.js", `
var earlyGlobalVar = 20;

// Redeclaration of global var from script_0.js:
var dupGlobalVar = 30;

// Suppressed duplicate global declaration:
/** @suppress {duplicate} */
var dupGlobalVar = 40;

// Reassignment of global const from script_0.js:
globalConst = 6;
`),
		SourceFile_FromCode("script_2.js", `
function globalHoistedFn(a) {
  // Local early reference and unused local inside a jumped-to global hoisted function:
  sink(localEarly);
  var localEarly = a + 1;
  var unusedInHoisted = 99;
  return localEarly;
}

function defaultParamCheck(a = b, b = 1) {
  return a + b;
}
sink(defaultParamCheck(2, 3));
`),
		SourceFile_FromCode("mod_0.js", `goog.module('test.mod0');
const earlyUse = exports;
exports = {earlyUse};
`),
	}

	for i := 0; i < 8; i++ {
		inputs = append(inputs, SourceFile_FromCode(
			fmt.Sprintf("worker_local_%d.js", i),
			fmt.Sprintf(`
function localWorkerFn_%[1]d(x) {
  var dupLocal = x;
  var dupLocal = x + 1;
  var unreadLocal = 42;
  return dupLocal;
}
sink(localWorkerFn_%[1]d(%[1]d));
`, i)))
	}

	seqErrs, seqWarns, seqReport := compileVarRefCheckWithThreads(t, externs, inputs, 1)
	parErrs, parWarns, parReport := compileVarRefCheckWithThreads(t, externs, inputs, 4)

	if strings.Join(seqErrs, "\n") != strings.Join(parErrs, "\n") {
		t.Fatalf("errors mismatch:\n--- seq (threads=1) ---\n%s\n--- par (threads=4) ---\n%s",
			strings.Join(seqErrs, "\n"), strings.Join(parErrs, "\n"))
	}
	if strings.Join(seqWarns, "\n") != strings.Join(parWarns, "\n") {
		t.Fatalf("warnings mismatch:\n--- seq (threads=1) ---\n%s\n--- par (threads=4) ---\n%s",
			strings.Join(seqWarns, "\n"), strings.Join(parWarns, "\n"))
	}
	if seqReport != parReport {
		t.Fatalf("formatted report mismatch:\n--- seq (threads=1) ---\n%s\n--- par (threads=4) ---\n%s",
			seqReport, parReport)
	}

	all := strings.Join(append(seqErrs, seqWarns...), "\n")
	for _, want := range []string{
		"JSC_REFERENCE_BEFORE_DECLARE",
		"JSC_VAR_MULTIPLY_DECLARED_ERROR",
		"JSC_REASSIGNED_CONSTANT",
		"JSC_REDECLARED_VARIABLE",
		"JSC_UNUSED_LOCAL_ASSIGNMENT",
	} {
		if !strings.Contains(all, want) {
			t.Fatalf("expected diagnostic %s in output, got:\n%s", want, all)
		}
	}
}
