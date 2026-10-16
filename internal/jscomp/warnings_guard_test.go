// Port of test/com/google/javascript/jscomp/WarningsGuardTest.java @ closure-compiler v20260415

package jscomp

import (
	"math"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// DiagnosticGroups is not ported in this wave: the groups the tests use are
// replaced by stand-ins with the same names and with distinct diagnostic
// types of the same default levels (the guards only depend on group
// membership).
var (
	jsaWgDETERMINISTIC_TEST = DiagnosticType_Warning("JSC_DETERMINISTIC_TEST", "deterministic test")
	jsaWgDEPRECATED_NAME    = DiagnosticType_Disabled("JSC_DEPRECATED_VAR", "Variable {0} has been deprecated.")

	jsaWgACCESS_CONTROLS = NewDiagnosticGroupWithName(jsaStr("visibility"),
		DiagnosticType_Warning("JSC_BAD_PRIVATE_GLOBAL_ACCESS", "x"),
		DiagnosticType_Warning("JSC_BAD_PRIVATE_PROPERTY_ACCESS", "x"))
	jsaWgVISIBILITY           = jsaWgACCESS_CONTROLS
	jsaWgDEPRECATED           = NewDiagnosticGroupWithName(jsaStr("deprecated"), jsaWgDEPRECATED_NAME)
	jsaWgCHECK_TYPES          = NewDiagnosticGroupWithName(jsaStr("checkTypes"), jsaWgDETERMINISTIC_TEST, DiagnosticType_Warning("JSC_TYPE_MISMATCH", "x"))
	jsaWgMESSAGE_DESCRIPTIONS = NewDiagnosticGroupWithName(jsaStr("msgDescriptions"), DiagnosticType_Disabled("JSC_MSG_HAS_NO_DESCRIPTION", "x"))

	jsaWgBAR_WARNING = DiagnosticType_Warning("BAR", "Bar description")

	jsaWgVisibilityOff     = NewDiagnosticGroupWarningsGuard(jsaWgACCESS_CONTROLS, CheckLevel_OFF)
	jsaWgVisibilityWarning = NewDiagnosticGroupWarningsGuard(jsaWgACCESS_CONTROLS, CheckLevel_WARNING)
)

func jsaWgLevelIs(t *testing.T, got *CheckLevel, want *CheckLevel, ctx ...any) {
	t.Helper()
	if (got == nil) != (want == nil) || (got != nil && *got != *want) {
		t.Errorf("%sgot %v, want %v", jsaCtx(ctx), got, want)
	}
}

var (
	jsaWgERROR   = CheckLevel_ERROR.Ptr()
	jsaWgWARNING = CheckLevel_WARNING.Ptr()
	jsaWgOFF     = CheckLevel_OFF.Ptr()
)

func TestWarningsGuard_ShowByPathGuard_Restrict(t *testing.T) {
	includeGuard := NewShowByPathWarningsGuardWithShowType("/foo/", ShowByPathWarningsGuard_ShowType_INCLUDE)

	jsaWgLevelIs(t, includeGuard.Level(jsaWgMakeErrorLevel(jsaStr("asasasd/foo/hello.js"), CheckLevel_WARNING)), nil)
	jsaWgLevelIs(t, includeGuard.Level(jsaWgMakeErrorLevel(jsaStr("asasasd/foo/hello.js"), CheckLevel_ERROR)), nil)
	jsaWgLevelIs(t, includeGuard.Level(jsaWgMakeErrorLevel(jsaStr("asasasd/hello.js"), CheckLevel_WARNING)), jsaWgOFF)
	jsaWgLevelIs(t, includeGuard.Level(jsaWgMakeErrorLevel(jsaStr("asasasd/hello.js"), CheckLevel_OFF)), jsaWgOFF)
	jsaWgLevelIs(t, includeGuard.Level(jsaWgMakeErrorLevel(jsaStr("asasasd/hello.js"), CheckLevel_ERROR)), nil)
	jsaWgLevelIs(t, includeGuard.Level(jsaWgMakeError(nil)), nil)
	jsaWgLevelIs(t, includeGuard.Level(jsaWgMakeErrorLevel(nil, CheckLevel_WARNING)), nil)

	jsaEq(t, includeGuard.MustRunChecks(jsaWgDEPRECATED), base.Tri_UNKNOWN)
}

func TestWarningsGuard_ShowByPathGuard_Suppress(t *testing.T) {
	excludeGuard := NewShowByPathWarningsGuardWithPathsAndShowType([]string{"/foo/", "/bar/"}, ShowByPathWarningsGuard_ShowType_EXCLUDE)

	jsaWgLevelIs(t, excludeGuard.Level(jsaWgMakeErrorLevel(jsaStr("asasasd/foo/hello.js"), CheckLevel_WARNING)), jsaWgOFF)
	jsaWgLevelIs(t, excludeGuard.Level(jsaWgMakeErrorLevel(jsaStr("asasasd/foo/bar/hello.js"), CheckLevel_WARNING)), jsaWgOFF)
	jsaWgLevelIs(t, excludeGuard.Level(jsaWgMakeErrorLevel(jsaStr("asasasd/bar/hello.js"), CheckLevel_WARNING)), jsaWgOFF)
	jsaWgLevelIs(t, excludeGuard.Level(jsaWgMakeErrorLevel(jsaStr("asasasd/foo/bar/hello.js"), CheckLevel_WARNING)), jsaWgOFF)
	jsaWgLevelIs(t, excludeGuard.Level(jsaWgMakeErrorLevel(jsaStr("asasasd/foo/hello.js"), CheckLevel_ERROR)), nil)
	jsaWgLevelIs(t, excludeGuard.Level(jsaWgMakeErrorLevel(jsaStr("asasasd/hello.js"), CheckLevel_WARNING)), nil)
	jsaWgLevelIs(t, excludeGuard.Level(jsaWgMakeErrorLevel(jsaStr("asasasd/hello.js"), CheckLevel_OFF)), nil)
	jsaWgLevelIs(t, excludeGuard.Level(jsaWgMakeErrorLevel(jsaStr("asasasd/hello.js"), CheckLevel_ERROR)), nil)
	jsaWgLevelIs(t, excludeGuard.Level(jsaWgMakeError(nil)), nil)
	jsaWgLevelIs(t, excludeGuard.Level(jsaWgMakeErrorLevel(nil, CheckLevel_WARNING)), nil)

	jsaEq(t, excludeGuard.MustRunChecks(jsaWgDEPRECATED), base.Tri_UNKNOWN)
}

func TestWarningsGuard_StrictGuard(t *testing.T) {
	guard := NewStrictWarningsGuard()

	jsaWgLevelIs(t, guard.Level(jsaWgMakeErrorLevel(jsaStr("foo/hello.js"), CheckLevel_WARNING)), jsaWgERROR)
	jsaWgLevelIs(t, guard.Level(jsaWgMakeErrorLevel(jsaStr("foo/hello.js"), CheckLevel_OFF)), nil)
	jsaWgLevelIs(t, guard.Level(jsaWgMakeErrorLevel(jsaStr("bar.js"), CheckLevel_ERROR)), jsaWgERROR)

	jsaEq(t, guard.MustRunChecks(jsaWgDEPRECATED), base.Tri_UNKNOWN)
}

func TestWarningsGuard_ByPathGuard(t *testing.T) {
	strictGuard := ByPathWarningsGuard_ForPath([]string{"/foo/"}, CheckLevel_ERROR)

	jsaWgLevelIs(t, strictGuard.Level(jsaWgMakeErrorLevel(jsaStr("asasasd/foo/hello.js"), CheckLevel_WARNING)), jsaWgERROR)
	jsaWgLevelIs(t, strictGuard.Level(jsaWgMakeErrorLevel(jsaStr("asasasd/foo/hello.js"), CheckLevel_ERROR)), nil)
	jsaWgLevelIs(t, strictGuard.Level(jsaWgMakeErrorLevel(jsaStr("asasasd/hello.js"), CheckLevel_WARNING)), nil)
	jsaWgLevelIs(t, strictGuard.Level(jsaWgMakeErrorLevel(jsaStr("asasasd/hello.js"), CheckLevel_OFF)), nil)
	jsaWgLevelIs(t, strictGuard.Level(jsaWgMakeErrorLevel(jsaStr("asasasd/hello.js"), CheckLevel_ERROR)), nil)
	jsaWgLevelIs(t, strictGuard.Level(jsaWgMakeError(nil)), nil)
	jsaWgLevelIs(t, strictGuard.Level(jsaWgMakeErrorLevel(nil, CheckLevel_WARNING)), nil)

	jsaEq(t, strictGuard.MustRunChecks(jsaWgDEPRECATED), base.Tri_UNKNOWN)
}

// jsaWgFuncGuard is an anonymous WarningsGuard subclass.
type jsaWgFuncGuard struct {
	WarningsGuardBase
	level func(error *JSError) *CheckLevel
}

func (g *jsaWgFuncGuard) Level(error *JSError) *CheckLevel { return g.level(error) }

func TestWarningsGuard_ComposeGuard(t *testing.T) {
	g1 := &jsaWgFuncGuard{level: func(error *JSError) *CheckLevel {
		if *error.SourceName() == "123456" {
			return jsaWgERROR
		}
		return nil
	}}

	g2 := &jsaWgFuncGuard{level: func(error *JSError) *CheckLevel {
		if error.GetLineNumber() == 12 {
			return jsaWgWARNING
		}
		return nil
	}}

	guard := NewComposeWarningsGuard(g1, g2)

	jsaWgLevelIs(t, guard.Level(jsaWgMakeError(jsaStr("aaa"))), nil)
	jsaWgLevelIs(t, guard.Level(jsaWgMakeError(jsaStr("12345"))), nil)
	jsaWgLevelIs(t, guard.Level(jsaWgMakeError(jsaStr("123456"))), jsaWgERROR)
	jsaWgLevelIs(t, guard.Level(jsaWgMakeErrorLine("12345", 12)), jsaWgWARNING)
	jsaWgLevelIs(t, guard.Level(jsaWgMakeErrorLine("12345", 13)), nil)

	jsaEq(t, guard.MustRunChecks(jsaWgDEPRECATED), base.Tri_UNKNOWN)
}

func TestWarningsGuard_ComposeGuard2(t *testing.T) {
	pathGuard := NewShowByPathWarningsGuard("/foo/")
	strictGuard := NewStrictWarningsGuard()

	guard := NewComposeWarningsGuard(strictGuard, pathGuard)
	jsaWgLevelIs(t, guard.Level(jsaWgMakeErrorLevel(jsaStr("asasasd/hello.js"), CheckLevel_WARNING)), jsaWgOFF)
	jsaWgLevelIs(t, guard.Level(jsaWgMakeErrorLevel(jsaStr("asasasd/foo/hello.js"), CheckLevel_WARNING)), jsaWgERROR)

	// Try again with the guards reversed.
	guard = NewComposeWarningsGuard(pathGuard, strictGuard)
	jsaWgLevelIs(t, guard.Level(jsaWgMakeErrorLevel(jsaStr("asasasd/hello.js"), CheckLevel_WARNING)), jsaWgOFF)
	jsaWgLevelIs(t, guard.Level(jsaWgMakeErrorLevel(jsaStr("asasasd/foo/hello.js"), CheckLevel_WARNING)), jsaWgERROR)
}

func TestWarningsGuard_ComposeGuard3(t *testing.T) {
	// Confirm that explicit diagnostic groups override the promotion of
	// warnings to errors done by StrictWarningGuard.

	typeGuard := NewDiagnosticGroupWarningsGuard(jsaWgDEPRECATED, CheckLevel_WARNING)
	strictGuard := NewStrictWarningsGuard()

	guard := NewComposeWarningsGuard(strictGuard, typeGuard)
	jsaWgLevelIs(t, guard.Level(JSError_MakeWithSourceLocation(jsaStr("example.js"), 1, 0, jsaWgDEPRECATED_NAME)), jsaWgWARNING)

	// Ordering applied doesn't matter, do it again reversed.
	guard = NewComposeWarningsGuard(typeGuard, strictGuard)
	jsaWgLevelIs(t, guard.Level(JSError_MakeWithSourceLocation(jsaStr("example.js"), 1, 0, jsaWgDEPRECATED_NAME)), jsaWgWARNING)
}

func TestWarningsGuard_ComposeGuardOrdering(t *testing.T) {
	pathGuard1 := NewShowByPathWarningsGuard("/foo/")
	pathGuard2 := NewShowByPathWarningsGuard("/bar/")
	strictGuard := NewStrictWarningsGuard()
	diagnosticGuard1 := NewDiagnosticGroupWarningsGuard(jsaWgACCESS_CONTROLS, CheckLevel_OFF)
	diagnosticGuard2 := NewDiagnosticGroupWarningsGuard(jsaWgACCESS_CONTROLS, CheckLevel_OFF)
	diagnosticGuard3 := NewDiagnosticGroupWarningsGuard(jsaWgACCESS_CONTROLS, CheckLevel_OFF)

	guard := NewComposeWarningsGuard(
		pathGuard1,
		diagnosticGuard1,
		strictGuard,
		diagnosticGuard2,
		pathGuard2,
		diagnosticGuard3)

	guards := guard.GetGuards()
	jsaEq(t, len(guards), 6)

	prevPriority := math.MinInt32
	for _, g := range guards {
		jsaTrue(t, g.GetPriority() >= prevPriority)
		prevPriority = g.GetPriority()
	}
}

func TestWarningsGuard_ComposeGuardOrdering2(t *testing.T) {
	// Ensure that guards added later always override, when two guards
	// have the same priority.
	guardA := NewComposeWarningsGuard()
	guardA.AddGuard(jsaWgVisibilityWarning)
	guardA.AddGuard(jsaWgVisibilityOff)
	guardA.AddGuard(jsaWgVisibilityWarning)

	guardB := NewComposeWarningsGuard()
	guardB.AddGuard(jsaWgVisibilityOff)
	guardB.AddGuard(jsaWgVisibilityWarning)
	guardB.AddGuard(jsaWgVisibilityOff)

	jsaEq(t, guardA.MustRunChecks(jsaWgACCESS_CONTROLS), base.Tri_TRUE)
	jsaEq(t, guardB.MustRunChecks(jsaWgACCESS_CONTROLS), base.Tri_FALSE)
}

func TestWarningsGuard_ComposeGuardOrdering3(t *testing.T) {
	guardA := NewComposeWarningsGuard()
	guardA.AddGuard(NewDiagnosticGroupWarningsGuard(jsaWgACCESS_CONTROLS, CheckLevel_WARNING))
	guardA.AddGuard(NewDiagnosticGroupWarningsGuard(jsaWgACCESS_CONTROLS, CheckLevel_WARNING))
	guardA.AddGuard(NewDiagnosticGroupWarningsGuard(jsaWgACCESS_CONTROLS, CheckLevel_WARNING))
	guardA.AddGuard(NewDiagnosticGroupWarningsGuard(jsaWgACCESS_CONTROLS, CheckLevel_WARNING))
	guardA.AddGuard(jsaWgVisibilityOff)

	jsaEq(t, guardA.MustRunChecks(jsaWgACCESS_CONTROLS), base.Tri_FALSE)
}

func TestWarningsGuard_ComposeGuardOrdering4(t *testing.T) {
	guardA := NewComposeWarningsGuard()
	guardA.AddGuard(jsaWgVisibilityWarning)
	guardA.AddGuard(jsaWgVisibilityWarning)
	guardA.AddGuard(jsaWgVisibilityWarning)
	guardA.AddGuard(jsaWgVisibilityWarning)
	guardA.AddGuard(jsaWgVisibilityWarning)
	guardA.AddGuard(jsaWgVisibilityWarning)
	guardA.AddGuard(jsaWgVisibilityOff)

	jsaEq(t, guardA.MustRunChecks(jsaWgACCESS_CONTROLS), base.Tri_FALSE)
}

func TestWarningsGuard_DiagnosticGuard1(t *testing.T) {
	guard := NewDiagnosticGroupWarningsGuard(jsaWgCHECK_TYPES, CheckLevel_ERROR)

	jsaWgLevelIs(t, guard.Level(jsaWgMakeErrorType("foo", jsaWgDETERMINISTIC_TEST)), jsaWgERROR)

	jsaEq(t, guard.MustRunChecks(jsaWgCHECK_TYPES), base.Tri_TRUE)
	jsaEq(t, guard.MustRunChecks(jsaWgMESSAGE_DESCRIPTIONS), base.Tri_UNKNOWN)
}

func TestWarningsGuard_DiagnosticGuard3(t *testing.T) {
	guard := NewDiagnosticGroupWarningsGuard(jsaWgCHECK_TYPES, CheckLevel_OFF)

	jsaEq(t, guard.MustRunChecks(jsaWgCHECK_TYPES), base.Tri_FALSE)
	jsaEq(t, guard.MustRunChecks(jsaWgMESSAGE_DESCRIPTIONS), base.Tri_UNKNOWN)
}

func TestWarningsGuard_DiagnosticGuard4(t *testing.T) {
	guard := NewDiagnosticGroupWarningsGuard(jsaWgDEPRECATED, CheckLevel_OFF)

	jsaEq(t, guard.MustRunChecks(jsaWgDEPRECATED), base.Tri_FALSE)
	jsaEq(t, guard.MustRunChecks(jsaWgVISIBILITY), base.Tri_UNKNOWN)
	jsaEq(t, guard.MustRunChecks(jsaWgACCESS_CONTROLS), base.Tri_UNKNOWN)
	jsaEq(t, guard.MustRunChecks(jsaWgMESSAGE_DESCRIPTIONS), base.Tri_UNKNOWN)
}

func TestWarningsGuard_SuppressGuard1(t *testing.T) {
	t.Skip("needs SuppressDocWarningsGuard, Compiler")
}
func TestWarningsGuard_SuppressGuard2(t *testing.T) {
	t.Skip("needs SuppressDocWarningsGuard, Compiler")
}
func TestWarningsGuard_SuppressGuard3(t *testing.T) {
	t.Skip("needs SuppressDocWarningsGuard, Compiler")
}
func TestWarningsGuard_SuppressGuard_strictMissingPropertyOnUnionTypes(t *testing.T) {
	t.Skip("needs SuppressDocWarningsGuard, Compiler")
}
func TestWarningsGuard_SuppressGuard4(t *testing.T) {
	t.Skip("needs SuppressDocWarningsGuard, Compiler")
}
func TestWarningsGuard_SuppressGuard5(t *testing.T) {
	t.Skip("needs SuppressDocWarningsGuard, Compiler")
}
func TestWarningsGuard_SuppressGuard6(t *testing.T) {
	t.Skip("needs SuppressDocWarningsGuard, Compiler")
}
func TestWarningsGuard_SuppressGuard7(t *testing.T) {
	t.Skip("needs SuppressDocWarningsGuard, Compiler")
}
func TestWarningsGuard_SuppressDocGuard_appliesSuppressionsOnComputedPropMethod_toPropNameExpression(t *testing.T) {
	t.Skip("needs SuppressDocWarningsGuard, Compiler")
}
func TestWarningsGuard_SuppressGuard_onCompoundAssignment(t *testing.T) {
	t.Skip("needs SuppressDocWarningsGuard, Compiler")
}
func TestWarningsGuard_SuppressGuard_onDetachedNode(t *testing.T) {
	t.Skip("needs SuppressDocWarningsGuard, Compiler")
}

func TestWarningsGuard_ComposeGuardCycle(t *testing.T) {
	guard := NewComposeWarningsGuard(jsaWgVisibilityOff, jsaWgVisibilityWarning)
	guard.AddGuard(guard)
	jsaEq(t, guard.String(), "DiagnosticGroup<visibility>(WARNING), DiagnosticGroup<visibility>(OFF)")
}

func jsaWgMakeError(sourcePath *string) *JSError {
	n := rhino.NewNode(rhino.Token_EMPTY)
	if sourcePath != nil {
		n.SetSourceFileForTesting(*sourcePath)
	}
	return JSError_MakeWithNode(n, jsaWgBAR_WARNING)
}

func jsaWgMakeErrorType(sourcePath string, type_ *DiagnosticType) *JSError {
	n := rhino.NewNode(rhino.Token_EMPTY)
	n.SetSourceFileForTesting(sourcePath)
	return JSError_MakeWithNode(n, type_)
}

func jsaWgMakeErrorLevel(sourcePath *string, level CheckLevel) *JSError {
	n := rhino.NewNode(rhino.Token_EMPTY)
	if sourcePath != nil {
		n.SetSourceFileForTesting(*sourcePath)
	}
	return JSError_MakeWithNode(n, DiagnosticType_Make("FOO", level, "Foo description"))
}

func jsaWgMakeErrorLine(sourcePath string, lineno int) *JSError {
	return JSError_MakeWithSourceLocation(&sourcePath, lineno, -1, jsaWgBAR_WARNING)
}

// Go-only: extra coverage — the other guards of this wave.
func TestWarningsGuard_OtherGuards(t *testing.T) {
	ts := NewCheckTypeSummaryWarningsGuard(CheckLevel_OFF)
	jsaWgLevelIs(t, ts.Level(jsaWgMakeErrorLevel(jsaStr("a.i.js"), CheckLevel_ERROR)), jsaWgOFF)
	jsaWgLevelIs(t, ts.Level(jsaWgMakeErrorLevel(jsaStr("a.js"), CheckLevel_ERROR)), nil)
	jsaWgLevelIs(t, ts.Level(jsaWgMakeErrorType("a.i.js", ClosurePrimitiveErrors_DUPLICATE_NAMESPACE_AND_MODULE)), nil)
	jsaEq(t, ts.GetPriority(), 40)

	ps := NewDiagnosticGroupPathSuppressingWarningsGuard(jsaWgDEPRECATED, "gen/")
	jsaWgLevelIs(t, ps.Level(jsaWgMakeErrorType("x/gen/a.js", jsaWgDEPRECATED_NAME)), jsaWgOFF)
	jsaWgLevelIs(t, ps.Level(jsaWgMakeErrorType("x/a.js", jsaWgDEPRECATED_NAME)), nil)
	jsaWgLevelIs(t, ps.Level(jsaWgMakeErrorType("x/gen/a.js", jsaWgBAR_WARNING)), nil)
	jsaEq(t, ps.String(), "DiagnosticGroup<deprecated>(gen/)")
	jsaEq(t, ps.GetPriority(), 50)

	cu := NewClosureUnawareCodeWarningsGuard()
	n := rhino.NewNode(rhino.Token_EMPTY)
	jsaWgLevelIs(t, cu.Level(JSError_MakeWithNode(n, RhinoErrorReporter_TYPE_PARSE_ERROR)), nil)
	jsaWgLevelIs(t, cu.Level(JSError_Make(RhinoErrorReporter_TYPE_PARSE_ERROR)), nil)
	jsaEq(t, cu.GetPriority(), 1)
}
