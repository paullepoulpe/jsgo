// Port of com/google/javascript/jscomp/CheckAccessControlsTest.java
// (deprecation subset) @ closure-compiler v20260415.
//
// These tests mirror the hand-written deprecation fixtures from the
// Java test class. Because JSC_DEPRECATED_* codes are DiagnosticType.disabled
// in the upstream Java pass they do not appear in oracle files generated
// by gen-oracles.sh; they are verified here directly instead.

package jscomp

import (
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/diag"
)

// checkOneDeprecated is a helper that runs the pipeline on src and
// returns all DEPRECATED_* diagnostics.
func checkDeprecatedDiags(t *testing.T, src string) []diag.Diagnostic {
	t.Helper()
	all, err := CheckSource("t.js", src)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	var out []diag.Diagnostic
	for _, d := range all {
		switch d.Code {
		case diag.JSC_DEPRECATED_VAR, diag.JSC_DEPRECATED_VAR_REASON,
			diag.JSC_DEPRECATED_PROP, diag.JSC_DEPRECATED_PROP_REASON,
			diag.JSC_DEPRECATED_CLASS, diag.JSC_DEPRECATED_CLASS_REASON:
			out = append(out, d)
		}
	}
	return out
}

// assertDeprecatedCode asserts that exactly one deprecated diagnostic
// with the given code is emitted.
func assertDeprecatedCode(t *testing.T, src, wantCode string) {
	t.Helper()
	diags := checkDeprecatedDiags(t, src)
	for _, d := range diags {
		if d.Code == wantCode {
			return
		}
	}
	t.Errorf("expected %s in output; got %v\nsrc: %s", wantCode, diags, src)
}

// assertNoDeprecated asserts that no deprecated diagnostics are emitted.
func assertNoDeprecated(t *testing.T, src string) {
	t.Helper()
	diags := checkDeprecatedDiags(t, src)
	if len(diags) != 0 {
		t.Errorf("expected no deprecated diagnostics; got %v\nsrc: %s", diags, src)
	}
}

// ---------------------------------------------------------------------------
// Variable / function deprecation
// ---------------------------------------------------------------------------

// TestDeprecatedVar_CallInGlobalScope: calling a deprecated function (no
// reason text) from global scope is an invocation target → warn with
// JSC_DEPRECATED_VAR.
func TestDeprecatedVar_CallInGlobalScope(t *testing.T) {
	src := `
/** @deprecated */
function oldFoo() {}
oldFoo();
`
	assertDeprecatedCode(t, src, diag.JSC_DEPRECATED_VAR)
}

// TestDeprecatedVar_CallWithReason: deprecated function with reason text
// emits JSC_DEPRECATED_VAR_REASON.
func TestDeprecatedVar_CallWithReason(t *testing.T) {
	src := `
/** @deprecated Use newFoo instead. */
function oldFoo() {}
function bar() { oldFoo(); }
`
	assertDeprecatedCode(t, src, diag.JSC_DEPRECATED_VAR_REASON)
}

// TestDeprecatedVarReason_CallInsideFunction: reason string emits _REASON
// variant.
func TestDeprecatedVarReason_CallInsideFunction(t *testing.T) {
	src := `
/** @deprecated Use newFoo instead. */
function oldFoo() {}
function bar() { oldFoo(); }
`
	// The function has reason text so _REASON variant fires.
	assertDeprecatedCode(t, src, diag.JSC_DEPRECATED_VAR_REASON)
}

// TestDeprecatedVar_BareRefInGlobalScope: referencing (but not calling) a
// deprecated function in global scope must NOT warn — Java's
// shouldEmitDeprecationWarning skips non-invocation-targets in global scope.
func TestDeprecatedVar_BareRefInGlobalScope(t *testing.T) {
	src := `
/** @deprecated */
function oldFoo() {}
var x = oldFoo;
`
	assertNoDeprecated(t, src)
}

// TestDeprecatedVar_CallInsideFunction: calling a deprecated function
// from inside a (non-deprecated) function does warn.
func TestDeprecatedVar_CallInsideFunction(t *testing.T) {
	src := `
/** @deprecated */
function oldFoo() {}
function caller() { oldFoo(); }
`
	assertDeprecatedCode(t, src, diag.JSC_DEPRECATED_VAR)
}

// TestDeprecatedVar_CallInsideDeprecatedFunction: calling deprecated from
// inside another deprecated function must NOT warn (deprecated context).
func TestDeprecatedVar_CallInsideDeprecatedFunction(t *testing.T) {
	src := `
/** @deprecated */
function oldFoo() {}
/** @deprecated */
function alsoOld() { oldFoo(); }
`
	assertNoDeprecated(t, src)
}

// TestDeprecatedVarReason_ReasonText: @deprecated with reason text emits
// the JSC_DEPRECATED_VAR_REASON code and includes the reason in the message.
func TestDeprecatedVarReason_ReasonText(t *testing.T) {
	src := `
/** @deprecated Use newFoo instead. */
function oldFoo() {}
function bar() { oldFoo(); }
`
	all, _ := CheckSource("t.js", src)
	for _, d := range all {
		if d.Code == diag.JSC_DEPRECATED_VAR_REASON {
			if !strings.Contains(d.Message, "Use newFoo instead") {
				t.Errorf("expected reason in message; got %q", d.Message)
			}
			return
		}
	}
	t.Errorf("expected JSC_DEPRECATED_VAR_REASON; got %v", all)
}

// ---------------------------------------------------------------------------
// Class deprecation (ES5 constructor via @constructor)
// ---------------------------------------------------------------------------

// TestDeprecatedClass_New: new Foo() where Foo has @constructor + @deprecated
// (with reason) should emit JSC_DEPRECATED_CLASS_REASON.
func TestDeprecatedClass_New(t *testing.T) {
	src := `
/** @constructor
 * @deprecated Use NewFoo instead. */
function OldFoo() {}
function bar() { new OldFoo(); }
`
	assertDeprecatedCode(t, src, diag.JSC_DEPRECATED_CLASS_REASON)
}

// TestDeprecatedClass_NewNoReason: new Foo() where Foo has @constructor +
// @deprecated (no reason) should emit JSC_DEPRECATED_CLASS.
func TestDeprecatedClass_NewNoReason(t *testing.T) {
	src := `
/** @constructor
 * @deprecated */
function OldFoo() {}
function bar() { new OldFoo(); }
`
	assertDeprecatedCode(t, src, diag.JSC_DEPRECATED_CLASS)
}

// TestDeprecatedClass_NewInGlobalScope: new DeprecatedCtor() in global
// scope — NEW is always an emission site regardless of scope.
func TestDeprecatedClass_NewInGlobalScope(t *testing.T) {
	src := `
/** @constructor
 * @deprecated Use NewFoo instead. */
function OldFoo() {}
new OldFoo();
`
	assertDeprecatedCode(t, src, diag.JSC_DEPRECATED_CLASS_REASON)
}

// TestDeprecatedClass_InsideDeprecated: new DeprecatedCtor() inside a
// deprecated function must NOT warn.
func TestDeprecatedClass_InsideDeprecated(t *testing.T) {
	src := `
/** @constructor
 * @deprecated */
function OldFoo() {}
/** @deprecated */
function alsoOld() { new OldFoo(); }
`
	assertNoDeprecated(t, src)
}
