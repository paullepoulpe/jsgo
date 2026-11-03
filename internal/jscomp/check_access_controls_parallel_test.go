package jscomp

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

func compileAccessControlsWithThreads(
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
	opts.SetWarningLevel(DiagnosticGroups_ACCESS_CONTROLS, CheckLevel_WARNING)
	opts.SetWarningLevel(DiagnosticGroups_DEPRECATED, CheckLevel_WARNING)
	opts.SetWarningLevel(DiagnosticGroups_CONSTANT_PROPERTY, CheckLevel_WARNING)
	opts.SetWarningLevel(DiagnosticGroups_VISIBILITY, CheckLevel_WARNING)
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

func assertAccessControlsParallelParity(
	t *testing.T,
	externs []*SourceFile,
	inputs []*SourceFile,
) (errs []string, warns []string) {
	t.Helper()
	seqErrs, seqWarns, seqReport := compileAccessControlsWithThreads(t, externs, inputs, 1)
	parErrs, parWarns, parReport := compileAccessControlsWithThreads(t, externs, inputs, 4)

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
	return seqErrs, seqWarns
}

func TestCheckAccessControlsParallel_CrossScriptConstPropsAndVisibility(t *testing.T) {
	externs := []*SourceFile{
		SourceFile_FromCode("extern_0.js", `/** @externs */
/** @constructor */
function ExternType() {}
/** @const {number} */
ExternType.EXT_CONST = 1;
`),
		SourceFile_FromCode("extern_1.js", `/** @externs */
/** @param {...*} var_args */
function useVal(var_args) {}
`),
	}

	inputs := []*SourceFile{
		SourceFile_FromCode("pkg/base.js", `
/** @const */
var app = {};

/** @constructor */
app.Base = function() {
  /** @private {number} */
  this.secret_ = 10;
  /** @protected {number} */
  this.prot_ = 20;
};

/** @final {number} */
app.Base.prototype.finalProp = 1;

/** @const {number} */
app.Base.STATIC_CONST = 42;

/** @deprecated Use newFn instead. */
app.deprecatedFn = function() {};
`),
		SourceFile_FromCode("pkg/pkg_file.js", `/**
 * @fileoverview
 * @package
 */
app.pkgHelper = function() {
  return 7;
};
`),
		SourceFile_FromCode("other/sub_and_mutations.js", `
class Sub extends app.Base {
  method() {
    // Allowed protected access on subclass instance, disallowed private access:
    useVal(this.prot_, this.secret_);
  }
}

// @final override across scripts:
Sub.prototype.finalProp = 2;

// Cross-script @const reassignment (initialized in pkg/base.js):
app.Base.STATIC_CONST = 99;

// Cross-script extern @const reassignment (initialized in extern_0.js):
ExternType.EXT_CONST = 2;

// Cross-script @const deletion:
delete app.Base.STATIC_CONST;

// Cross-package access to @fileoverview @package symbol, @protected property outside hierarchy, and @deprecated call:
var extBase = new app.Base();
useVal(extBase.prot_, app.pkgHelper(), app.deprecatedFn());
`),
		SourceFile_FromCode("other/extends_scope.js", `
/** @constructor */
function Holder() {}
/** @protected {function(new:app.Base)} */
Holder.ProtCtor = app.Base;

class BadExtends extends Holder.ProtCtor {}
useVal(BadExtends);
`),
	}

	for i := 0; i < 6; i++ {
		inputs = append(inputs, SourceFile_FromCode(
			fmt.Sprintf("other/extra_%d.js", i),
			fmt.Sprintf(`
/** @constructor */
function LocalBox_%[1]d() {
  /** @const {number} */
  this.immutableVal = %[1]d;
}
var box_%[1]d = new LocalBox_%[1]d();
box_%[1]d.immutableVal = %[1]d + 1;
useVal(box_%[1]d.secret_);
`, i)))
	}

	errs, warns := assertAccessControlsParallelParity(t, externs, inputs)
	all := strings.Join(append(errs, warns...), "\n")
	for _, want := range []string{
		"JSC_BAD_PRIVATE_PROPERTY_ACCESS",
		"JSC_FINAL_PROPERTY_OVERRIDDEN",
		"JSC_CONSTANT_PROPERTY_REASSIGNED_VALUE",
		"JSC_CONSTANT_PROPERTY_DELETED",
		"JSC_BAD_PACKAGE_PROPERTY_ACCESS",
		"JSC_DEPRECATED_PROP_REASON",
		"JSC_BAD_PROTECTED_PROPERTY_ACCESS",
	} {
		if !strings.Contains(all, want) {
			t.Fatalf("expected diagnostic %s in output, got:\n%s", want, all)
		}
	}
}
