package jscomp

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Parity of the parallel paths of checkAstValidity, varCheckValidity and
// checkTypeImportCodeReferences with their sequential traversals, including
// the violation / panic paths.

const parallelValidityExterns = `/** @externs */
/** @type {!Object} */ var window;
/** @param {...*} args */ function sink(args) {}
`

func parallelValidityInputs() []string {
	var inputs []string
	for i := 0; i < 8; i++ {
		inputs = append(inputs, fmt.Sprintf(`
function f%[1]d(a, b) { var s = 0; for (var i = 0; i < a; i++) { s += b[i] * %[1]d; } return s; }
window['f%[1]d'] = f%[1]d; sink(f%[1]d(3, [1, 2, 3]));
`, i))
	}
	return inputs
}

// compiledForValidity compiles inputs in ADVANCED mode and returns the
// compiler, its AST still in place, set to run later passes on threads.
func compiledForValidity(t *testing.T, threads int) *Compiler {
	t.Helper()
	externs := []*SourceFile{
		SourceFile_FromCode("extern_0.js", parallelValidityExterns),
		SourceFile_FromCode("extern_1.js", "/** @externs */ /** @type {number} */ var extern_one;"),
	}
	var inputs []*SourceFile
	for i, code := range parallelValidityInputs() {
		inputs = append(inputs, SourceFile_FromCode(fmt.Sprintf("input_%d.js", i), code))
	}
	opts := NewCompilerOptions()
	CompilationLevel_ADVANCED_OPTIMIZATIONS.SetOptionsForCompilationLevel(opts)
	opts.SetNumParallelThreads(threads)
	var buf bytes.Buffer
	comp := NewCompilerWithOutStream(jutil.NewPrintStream(&buf))
	comp.DisableThreads()
	res := comp.Compile(externs, inputs, opts)
	if len(res.Errors) != 0 {
		t.Fatalf("compile errors: %v", res.Errors)
	}
	return comp
}

// noGrammarCheck turns off the AST grammar check that the package's tests
// turn on (ast_validator_test.go): the AST validator validates scripts on
// workers only without it.
func noGrammarCheck(t *testing.T) {
	saved := astValidatorChecksGrammar
	astValidatorChecksGrammar = false
	t.Cleanup(func() { astValidatorChecksGrammar = saved })
}

// runPass runs pass and returns what it panicked with (as text) and the
// errors reported meanwhile.
func runPass(comp *Compiler, pass CompilerPass) (panicked string, errs []string) {
	before := len(comp.GetErrors())
	func() {
		defer func() {
			if r := recover(); r != nil {
				panicked = fmt.Sprint(jutil.AsThrowable(r))
			}
		}()
		pass.Process(comp.GetExternsRoot(), comp.GetJsRoot())
	}()
	for _, e := range comp.GetErrors()[before:] {
		errs = append(errs, e.String())
	}
	return panicked, errs
}

// nthScript returns the n-th script of root.
func nthScript(root rhino.Node, n int) rhino.Node {
	s := root.GetFirstChild()
	for ; n > 0; n-- {
		s = s.GetNext()
	}
	return s
}

// corruptions break the AST of some scripts. Each returns the scripts it
// touched, for the test's message.
var validityCorruptions = []struct {
	name string
	// breakTree changes the tree of comp; the AST validator and VarCheck
	// in validity mode must fail on the first broken script in order.
	breakTree func(comp *Compiler)
}{
	{"none", func(*Compiler) {}},
	{"root scripts 2 and 5", func(comp *Compiler) {
		for _, n := range []int{5, 2} {
			s := nthScript(comp.GetJsRoot(), n)
			s.AddChildToFront(rhino.IR_ExprResult(rhino.IR_Name(fmt.Sprintf("undeclared_%d", n))))
			// A NAME with a child is invalid for the AST validator.
			bad := rhino.IR_Name(fmt.Sprintf("bad_%d", n))
			bad.AddChildToBack(rhino.IR_Number(1))
			s.AddChildToBack(rhino.IR_ExprResult(bad))
		}
	}},
	{"extern script 1 and root script 3", func(comp *Compiler) {
		s := nthScript(comp.GetJsRoot(), 3)
		s.AddChildToFront(rhino.IR_ExprResult(rhino.IR_Name("undeclared_root")))
		bad := rhino.IR_Name("bad_root")
		bad.AddChildToBack(rhino.IR_Number(1))
		s.AddChildToBack(rhino.IR_ExprResult(bad))
		e := nthScript(comp.GetExternsRoot(), 1)
		e.AddChildToBack(rhino.IR_ExprResult(rhino.IR_Name("undeclared_extern")))
		bad = rhino.IR_Name("bad_extern")
		bad.AddChildToBack(rhino.IR_Number(2))
		e.AddChildToBack(rhino.IR_ExprResult(bad))
	}},
}

func TestParallelPasses_ValidityPassesParity(t *testing.T) {
	noGrammarCheck(t)
	for _, tc := range validityCorruptions {
		t.Run(tc.name, func(t *testing.T) {
			type result struct{ astPanic, varPanic, varErrs string }
			var got [2]result
			for i, threads := range []int{1, 4} {
				comp := compiledForValidity(t, threads)
				tc.breakTree(comp)
				p, errs := runPass(comp, NewAstValidator(comp))
				got[i].astPanic = p + strings.Join(errs, "\n")
				p, errs = runPass(comp, NewVarCheckWithValidityCheck(comp, true))
				got[i].varPanic, got[i].varErrs = p, strings.Join(errs, "\n")
			}
			if got[0] != got[1] {
				t.Fatalf("threads=1 and threads=4 differ:\n--- 1 ---\n%+v\n--- 4 ---\n%+v", got[0], got[1])
			}
			broken := tc.name != "none"
			if broken != (got[0].astPanic != "") || broken != (got[0].varPanic != "") {
				t.Fatalf("want violations iff the tree is broken, got %+v", got[0])
			}
			if broken {
				first := map[string]string{"root scripts 2 and 5": "2", "extern script 1 and root script 3": "extern"}[tc.name]
				if !strings.Contains(got[0].astPanic, "bad_"+first) || !strings.Contains(got[0].varPanic, "undeclared_"+first) ||
					!strings.Contains(got[0].varErrs, "undeclared_"+first) {
					t.Fatalf("want the first broken script (%s) to fail, got %+v", first, got[0])
				}
			}
		})
	}
}

// A caller's violation handler sees every violation in order (the
// validator does not use workers for it).
func TestParallelPasses_AstValidatorCustomHandlerSequential(t *testing.T) {
	noGrammarCheck(t)
	var seen [2][]string
	for i, threads := range []int{1, 4} {
		comp := compiledForValidity(t, threads)
		validityCorruptions[1].breakTree(comp)
		v := NewAstValidatorWithViolationFunc(comp, func(msg string, n rhino.Node) {
			seen[i] = append(seen[i], msg+" "+n.String())
		}, false, false)
		v.Process(comp.GetExternsRoot(), comp.GetJsRoot())
	}
	if len(seen[0]) < 2 || strings.Join(seen[0], "\n") != strings.Join(seen[1], "\n") {
		t.Fatalf("handler calls differ or too few:\n%v\n%v", seen[0], seen[1])
	}
}

// checkTypeImportCodeReferences reports code references to
// goog.requireType'd names in several scripts in the same order at 1 and 4
// threads, and does nothing without goog.requireType.
func TestParallelPasses_CheckTypeImportCodeReferencesParity(t *testing.T) {
	base := `/** @fileoverview */
var goog = {};
goog.provide = function(ns) {};
goog.module = function(ns) {};
goog.require = function(ns) {};
goog.requireType = function(ns) {};
goog.module.declareLegacyNamespace = function() {};
`
	inputs := []string{base, "goog.module('a.A'); class A {} exports = A;"}
	for i := 0; i < 6; i++ {
		inputs = append(inputs, fmt.Sprintf(`goog.module('m.m%[1]d');
const A = goog.requireType('a.A');
/** @param {!A} x */ function use%[1]d(x) { return x; }
%[2]s
exports = {use%[1]d};`, i, map[bool]string{true: "const bad = new A();", false: ""}[i%2 == 1]))
	}
	assertParallelParity(t, nil, inputs, nil)
	_, errs, _, _ := compileWithParallelThreads(t, nil, inputs, 4, nil)
	n := 0
	for _, e := range errs {
		if strings.HasPrefix(e, "JSC_TYPE_IMPORT_CODE_REFERENCE") {
			n++
		}
	}
	if n != 3 {
		t.Fatalf("want 3 JSC_TYPE_IMPORT_CODE_REFERENCE, got %v", errs)
	}
	noRequireType := []string{base, "goog.module('a.A'); class A {} exports = A;", "goog.module('b'); const A = goog.require('a.A'); new A();"}
	assertParallelParity(t, nil, noRequireType, nil)
}
