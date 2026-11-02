// Port of test/com/google/javascript/jscomp/WarningsGuardTest.java @ closure-compiler v20261028

/*
 * Copyright 2008 The Closure Compiler Authors.
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
	"math"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Java: TypeCheck.DETERMINISTIC_TEST, CheckAccessControls.DEPRECATED_NAME,
// BAR_WARNING and the DiagnosticGroups the tests use.
var (
	jsaWgDETERMINISTIC_TEST *DiagnosticType
	jsaWgDEPRECATED_NAME    *DiagnosticType

	jsaWgACCESS_CONTROLS      *DiagnosticGroup
	jsaWgVISIBILITY           *DiagnosticGroup
	jsaWgDEPRECATED           *DiagnosticGroup
	jsaWgCHECK_TYPES          *DiagnosticGroup
	jsaWgMESSAGE_DESCRIPTIONS *DiagnosticGroup

	jsaWgBAR_WARNING = DiagnosticType_Warning("BAR", "Bar description")

	jsaWgVisibilityOff     WarningsGuard
	jsaWgVisibilityWarning WarningsGuard
)

// Go-only: language-adapter — the DiagnosticGroups static fields are assigned
// from diagnostic_groups.go's init() (Java static initializer order), so the
// test's static fields that read them are assigned from an init() that runs
// after it.
func init() {
	jsaWgDETERMINISTIC_TEST = TypeCheck_DETERMINISTIC_TEST
	jsaWgDEPRECATED_NAME = CheckAccessControls_DEPRECATED_NAME

	jsaWgACCESS_CONTROLS = DiagnosticGroups_ACCESS_CONTROLS
	jsaWgVISIBILITY = DiagnosticGroups_VISIBILITY
	jsaWgDEPRECATED = DiagnosticGroups_DEPRECATED
	jsaWgCHECK_TYPES = DiagnosticGroups_CHECK_TYPES
	jsaWgMESSAGE_DESCRIPTIONS = DiagnosticGroups_MESSAGE_DESCRIPTIONS

	jsaWgVisibilityOff = NewDiagnosticGroupWarningsGuard(DiagnosticGroups_ACCESS_CONTROLS, CheckLevel_OFF)
	jsaWgVisibilityWarning = NewDiagnosticGroupWarningsGuard(DiagnosticGroups_ACCESS_CONTROLS, CheckLevel_WARNING)
}

func jsaWgLevelIs(t *testing.T, guard WarningsGuard, err *JSError, want CheckLevel, ctx ...any) {
	t.Helper()
	got, ok := guard.Level(err)
	if !ok || got != want {
		t.Errorf("%sgot (%v, %t), want (%v, true)", jsaCtx(ctx), got, ok, want)
	}
}

func jsaWgNoLevel(t *testing.T, guard WarningsGuard, err *JSError, ctx ...any) {
	t.Helper()
	if got, ok := guard.Level(err); ok {
		t.Errorf("%sgot (%v, true), want no level", jsaCtx(ctx), got)
	}
}

func TestWarningsGuard_ShowByPathGuard_Restrict(t *testing.T) {
	includeGuard := NewShowByPathWarningsGuardWithShowType("/foo/", ShowByPathWarningsGuard_ShowType_INCLUDE)

	jsaWgNoLevel(t, includeGuard, jsaWgMakeErrorLevel(jsaStr("asasasd/foo/hello.js"), CheckLevel_WARNING))
	jsaWgNoLevel(t, includeGuard, jsaWgMakeErrorLevel(jsaStr("asasasd/foo/hello.js"), CheckLevel_ERROR))
	jsaWgLevelIs(t, includeGuard, jsaWgMakeErrorLevel(jsaStr("asasasd/hello.js"), CheckLevel_WARNING), CheckLevel_OFF)
	jsaWgLevelIs(t, includeGuard, jsaWgMakeErrorLevel(jsaStr("asasasd/hello.js"), CheckLevel_OFF), CheckLevel_OFF)
	jsaWgNoLevel(t, includeGuard, jsaWgMakeErrorLevel(jsaStr("asasasd/hello.js"), CheckLevel_ERROR))
	jsaWgNoLevel(t, includeGuard, jsaWgMakeError(nil))
	jsaWgNoLevel(t, includeGuard, jsaWgMakeErrorLevel(nil, CheckLevel_WARNING))

	jsaEq(t, includeGuard.MustRunChecks(jsaWgDEPRECATED), base.Tri_UNKNOWN)
}

func TestWarningsGuard_ShowByPathGuard_Suppress(t *testing.T) {
	excludeGuard := NewShowByPathWarningsGuardWithPathsAndShowType([]string{"/foo/", "/bar/"}, ShowByPathWarningsGuard_ShowType_EXCLUDE)

	jsaWgLevelIs(t, excludeGuard, jsaWgMakeErrorLevel(jsaStr("asasasd/foo/hello.js"), CheckLevel_WARNING), CheckLevel_OFF)
	jsaWgLevelIs(t, excludeGuard, jsaWgMakeErrorLevel(jsaStr("asasasd/foo/bar/hello.js"), CheckLevel_WARNING), CheckLevel_OFF)
	jsaWgLevelIs(t, excludeGuard, jsaWgMakeErrorLevel(jsaStr("asasasd/bar/hello.js"), CheckLevel_WARNING), CheckLevel_OFF)
	jsaWgLevelIs(t, excludeGuard, jsaWgMakeErrorLevel(jsaStr("asasasd/foo/bar/hello.js"), CheckLevel_WARNING), CheckLevel_OFF)
	jsaWgNoLevel(t, excludeGuard, jsaWgMakeErrorLevel(jsaStr("asasasd/foo/hello.js"), CheckLevel_ERROR))
	jsaWgNoLevel(t, excludeGuard, jsaWgMakeErrorLevel(jsaStr("asasasd/hello.js"), CheckLevel_WARNING))
	jsaWgNoLevel(t, excludeGuard, jsaWgMakeErrorLevel(jsaStr("asasasd/hello.js"), CheckLevel_OFF))
	jsaWgNoLevel(t, excludeGuard, jsaWgMakeErrorLevel(jsaStr("asasasd/hello.js"), CheckLevel_ERROR))
	jsaWgNoLevel(t, excludeGuard, jsaWgMakeError(nil))
	jsaWgNoLevel(t, excludeGuard, jsaWgMakeErrorLevel(nil, CheckLevel_WARNING))

	jsaEq(t, excludeGuard.MustRunChecks(jsaWgDEPRECATED), base.Tri_UNKNOWN)
}

func TestWarningsGuard_StrictGuard(t *testing.T) {
	guard := NewStrictWarningsGuard()

	jsaWgLevelIs(t, guard, jsaWgMakeErrorLevel(jsaStr("foo/hello.js"), CheckLevel_WARNING), CheckLevel_ERROR)
	jsaWgNoLevel(t, guard, jsaWgMakeErrorLevel(jsaStr("foo/hello.js"), CheckLevel_OFF))
	jsaWgLevelIs(t, guard, jsaWgMakeErrorLevel(jsaStr("bar.js"), CheckLevel_ERROR), CheckLevel_ERROR)

	jsaEq(t, guard.MustRunChecks(jsaWgDEPRECATED), base.Tri_UNKNOWN)
}

func TestWarningsGuard_ByPathGuard(t *testing.T) {
	strictGuard := ByPathWarningsGuard_ForPath([]string{"/foo/"}, CheckLevel_ERROR)

	jsaWgLevelIs(t, strictGuard, jsaWgMakeErrorLevel(jsaStr("asasasd/foo/hello.js"), CheckLevel_WARNING), CheckLevel_ERROR)
	jsaWgNoLevel(t, strictGuard, jsaWgMakeErrorLevel(jsaStr("asasasd/foo/hello.js"), CheckLevel_ERROR))
	jsaWgNoLevel(t, strictGuard, jsaWgMakeErrorLevel(jsaStr("asasasd/hello.js"), CheckLevel_WARNING))
	jsaWgNoLevel(t, strictGuard, jsaWgMakeErrorLevel(jsaStr("asasasd/hello.js"), CheckLevel_OFF))
	jsaWgNoLevel(t, strictGuard, jsaWgMakeErrorLevel(jsaStr("asasasd/hello.js"), CheckLevel_ERROR))
	jsaWgNoLevel(t, strictGuard, jsaWgMakeError(nil))
	jsaWgNoLevel(t, strictGuard, jsaWgMakeErrorLevel(nil, CheckLevel_WARNING))

	jsaEq(t, strictGuard.MustRunChecks(jsaWgDEPRECATED), base.Tri_UNKNOWN)
}

// jsaWgFuncGuard is an anonymous WarningsGuard subclass.
type jsaWgFuncGuard struct {
	WarningsGuardBase
	level func(error *JSError) (CheckLevel, bool)
}

func (g *jsaWgFuncGuard) Level(error *JSError) (CheckLevel, bool) { return g.level(error) }

func TestWarningsGuard_ComposeGuard(t *testing.T) {
	g1 := &jsaWgFuncGuard{level: func(error *JSError) (CheckLevel, bool) {
		if *error.SourceName() == "123456" {
			return CheckLevel_ERROR, true
		}
		return 0, false
	}}

	g2 := &jsaWgFuncGuard{level: func(error *JSError) (CheckLevel, bool) {
		if error.Lineno() == 12 {
			return CheckLevel_WARNING, true
		}
		return 0, false
	}}

	guard := NewComposeWarningsGuard(g1, g2)

	jsaWgNoLevel(t, guard, jsaWgMakeError(jsaStr("aaa")))
	jsaWgNoLevel(t, guard, jsaWgMakeError(jsaStr("12345")))
	jsaWgLevelIs(t, guard, jsaWgMakeError(jsaStr("123456")), CheckLevel_ERROR)
	jsaWgLevelIs(t, guard, jsaWgMakeErrorLine("12345", 12), CheckLevel_WARNING)
	jsaWgNoLevel(t, guard, jsaWgMakeErrorLine("12345", 13))

	jsaEq(t, guard.MustRunChecks(jsaWgDEPRECATED), base.Tri_UNKNOWN)
}

func TestWarningsGuard_ComposeGuard2(t *testing.T) {
	pathGuard := NewShowByPathWarningsGuard("/foo/")
	strictGuard := NewStrictWarningsGuard()

	guard := NewComposeWarningsGuard(strictGuard, pathGuard)
	jsaWgLevelIs(t, guard, jsaWgMakeErrorLevel(jsaStr("asasasd/hello.js"), CheckLevel_WARNING), CheckLevel_OFF)
	jsaWgLevelIs(t, guard, jsaWgMakeErrorLevel(jsaStr("asasasd/foo/hello.js"), CheckLevel_WARNING), CheckLevel_ERROR)

	// Try again with the guards reversed.
	guard = NewComposeWarningsGuard(pathGuard, strictGuard)
	jsaWgLevelIs(t, guard, jsaWgMakeErrorLevel(jsaStr("asasasd/hello.js"), CheckLevel_WARNING), CheckLevel_OFF)
	jsaWgLevelIs(t, guard, jsaWgMakeErrorLevel(jsaStr("asasasd/foo/hello.js"), CheckLevel_WARNING), CheckLevel_ERROR)
}

func TestWarningsGuard_ComposeGuard3(t *testing.T) {
	// Confirm that explicit diagnostic groups override the promotion of
	// warnings to errors done by StrictWarningGuard.

	typeGuard := NewDiagnosticGroupWarningsGuard(jsaWgDEPRECATED, CheckLevel_WARNING)
	strictGuard := NewStrictWarningsGuard()

	guard := NewComposeWarningsGuard(strictGuard, typeGuard)
	jsaWgLevelIs(t, guard, JSError_MakeWithSourceLocation(jsaStr("example.js"), 1, 0, jsaWgDEPRECATED_NAME), CheckLevel_WARNING)

	// Ordering applied doesn't matter, do it again reversed.
	guard = NewComposeWarningsGuard(typeGuard, strictGuard)
	jsaWgLevelIs(t, guard, JSError_MakeWithSourceLocation(jsaStr("example.js"), 1, 0, jsaWgDEPRECATED_NAME), CheckLevel_WARNING)
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

	jsaWgLevelIs(t, guard, jsaWgMakeErrorType("foo", jsaWgDETERMINISTIC_TEST), CheckLevel_ERROR)

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
	m := jutil.NewLinkedHashMap[string, *DiagnosticGroup]()
	m.Put("deprecated", NewDiagnosticGroup(jsaWgBAR_WARNING))
	compiler := NewCompiler()
	var guard WarningsGuard = NewSuppressDocWarningsGuard(compiler, m)

	code := compiler.ParseTestCode(
		"/** @suppress {deprecated} */ function f() { a; }\n" +
			"function g() { b; }\n")
	jsaWgNoLevel(t, guard, JSError_MakeWithNode(code, jsaWgBAR_WARNING))
	jsaWgLevelIs(t, guard, JSError_MakeWithNode(jsaWgFindNameNode(code, "a"), jsaWgBAR_WARNING), CheckLevel_OFF)
	jsaWgNoLevel(t, guard, JSError_MakeWithNode(jsaWgFindNameNode(code, "b"), jsaWgBAR_WARNING))
}

func TestWarningsGuard_SuppressGuard2(t *testing.T) {
	m := jutil.NewLinkedHashMap[string, *DiagnosticGroup]()
	m.Put("deprecated", NewDiagnosticGroup(jsaWgBAR_WARNING))
	compiler := NewCompiler()
	var guard WarningsGuard = NewSuppressDocWarningsGuard(compiler, m)

	code := compiler.ParseTestCode(
		"/** @fileoverview \n" +
			" * @suppress {deprecated} */ function f() { a; }\n" +
			"function g() { b; }\n")
	jsaWgLevelIs(t, guard, JSError_MakeWithNode(jsaWgFindNameNode(code, "a"), jsaWgBAR_WARNING), CheckLevel_OFF)
	jsaWgLevelIs(t, guard, JSError_MakeWithNode(jsaWgFindNameNode(code, "b"), jsaWgBAR_WARNING), CheckLevel_OFF)
}

func TestWarningsGuard_SuppressGuard3(t *testing.T) {
	m := jutil.NewLinkedHashMap[string, *DiagnosticGroup]()
	m.Put("deprecated", NewDiagnosticGroup(jsaWgBAR_WARNING))
	compiler := NewCompiler()
	var guard WarningsGuard = NewSuppressDocWarningsGuard(compiler, m)

	code := compiler.ParseTestCode("/** @suppress {deprecated} */ var f = function() { a; }")
	jsaWgLevelIs(t, guard, JSError_MakeWithNode(jsaWgFindNameNode(code, "a"), jsaWgBAR_WARNING), CheckLevel_OFF)
}

func TestWarningsGuard_SuppressGuard_strictMissingPropertyOnUnionTypes(t *testing.T) {
	compiler := NewCompiler()
	var guard WarningsGuard = NewSuppressDocWarningsGuard(compiler, DiagnosticGroups_GetRegisteredGroups())

	code := compiler.ParseTestCode(
		"class C {}\n" +
			"class D{}\n" +
			"/** @type {(C|D)} */\n" +
			"let obj;\n" +
			"/** @suppress {strictMissingProperties} */\n" +
			"obj.prop\n")
	jsaWgLevelIs(t,
		guard,
		JSError_MakeWithNode(
			jsaWgFindGetPropNode(code, "prop"), TypeCheck_ILLEGAL_PROPERTY_CREATION_ON_UNION_TYPE),
		CheckLevel_OFF)
}

func TestWarningsGuard_SuppressGuard4(t *testing.T) {
	m := jutil.NewLinkedHashMap[string, *DiagnosticGroup]()
	m.Put("deprecated", NewDiagnosticGroup(jsaWgBAR_WARNING))
	compiler := NewCompiler()
	var guard WarningsGuard = NewSuppressDocWarningsGuard(compiler, m)

	code := compiler.ParseTestCode(
		"var goog = {};\n" +
			"/** @suppress {deprecated} */ goog.f = function() { a; }\n")
	jsaWgLevelIs(t, guard, JSError_MakeWithNode(jsaWgFindNameNode(code, "a"), jsaWgBAR_WARNING), CheckLevel_OFF)
}

func TestWarningsGuard_SuppressGuard5(t *testing.T) {
	m := jutil.NewLinkedHashMap[string, *DiagnosticGroup]()
	m.Put("deprecated", NewDiagnosticGroup(jsaWgBAR_WARNING))
	compiler := NewCompiler()
	var guard WarningsGuard = NewSuppressDocWarningsGuard(compiler, m)

	code := compiler.ParseTestCode(
		"var goog = {};\n" +
			"goog.f = function() { /** @suppress {deprecated} */ (a); }\n")

	jsaWgLevelIs(t, guard, JSError_MakeWithNode(jsaWgFindNameNode(code, "a"), jsaWgBAR_WARNING), CheckLevel_OFF)
}

func TestWarningsGuard_SuppressGuard6(t *testing.T) {
	m := jutil.NewLinkedHashMap[string, *DiagnosticGroup]()
	m.Put("deprecated", NewDiagnosticGroup(jsaWgBAR_WARNING))
	compiler := NewCompiler()
	var guard WarningsGuard = NewSuppressDocWarningsGuard(compiler, m)

	code := compiler.ParseTestCode("/** @fileoverview @suppress {deprecated} */\n console.log(a);")

	jsaWgLevelIs(t, guard, JSError_MakeWithNode(jsaWgFindNameNode(code, "a"), jsaWgBAR_WARNING), CheckLevel_OFF)
}

func TestWarningsGuard_SuppressGuard7(t *testing.T) {
	m := jutil.NewLinkedHashMap[string, *DiagnosticGroup]()
	m.Put("deprecated", NewDiagnosticGroup(jsaWgBAR_WARNING))
	compiler := NewCompiler()
	var guard WarningsGuard = NewSuppressDocWarningsGuard(compiler, m)

	code := compiler.ParseTestCode("console.log(/** @suppress {deprecated} */ (a));")

	// We don't care about @suppress annotations within nested expressions
	jsaWgNoLevel(t, guard, JSError_MakeWithNode(jsaWgFindNameNode(code, "a"), jsaWgBAR_WARNING))
}

func TestWarningsGuard_SuppressDocGuard_appliesSuppressionsOnComputedPropMethod_toPropNameExpression(t *testing.T) {
	m := jutil.NewLinkedHashMap[string, *DiagnosticGroup]()
	m.Put("deprecated", NewDiagnosticGroup(jsaWgBAR_WARNING))

	compiler := NewCompiler()
	var guard WarningsGuard = NewSuppressDocWarningsGuard(compiler, m)

	code := compiler.ParseTestCode(
		"class Foo {\n" +
			"  /** @suppress {deprecated} */\n" +
			"  [a]() { }\n" +
			"}\n")

	jsaWgLevelIs(t, guard, JSError_MakeWithNode(jsaWgFindNameNode(code, "a"), jsaWgBAR_WARNING), CheckLevel_OFF)
}

func TestWarningsGuard_SuppressGuard_onCompoundAssignment(t *testing.T) {
	compiler := NewCompiler()
	m := jutil.NewLinkedHashMap[string, *DiagnosticGroup]()
	m.Put("deprecated", NewDiagnosticGroup(jsaWgBAR_WARNING))
	var guard WarningsGuard = NewSuppressDocWarningsGuard(compiler, m)

	code := compiler.ParseTestCode(
		"var goog = {};\n" +
			"/** @suppress {deprecated} */ goog.f += a\n")
	jsaWgLevelIs(t, guard, JSError_MakeWithNode(jsaWgFindNameNode(code, "a"), jsaWgBAR_WARNING), CheckLevel_OFF)
}

func TestWarningsGuard_SuppressGuard_onDetachedNode(t *testing.T) {
	compiler := NewCompiler()
	m := jutil.NewLinkedHashMap[string, *DiagnosticGroup]()
	m.Put("deprecated", NewDiagnosticGroup(jsaWgBAR_WARNING))
	var guard WarningsGuard = NewSuppressDocWarningsGuard(compiler, m)

	code := compiler.ParseTestCode("/** @fileoverview @suppress {deprecated} */\n\nvar x;")
	// Create an error that has the right source file, but no parents.
	// This is the state of JSDoc nodes.
	error := jsaWgMakeError(code.GetSourceFileName())
	jsaWgLevelIs(t, guard, error, CheckLevel_OFF)
}

func TestWarningsGuard_ComposeGuardCycle(t *testing.T) {
	guard := NewComposeWarningsGuard(jsaWgVisibilityOff, jsaWgVisibilityWarning)
	guard.AddGuard(guard)
	jsaEq(t, guard.String(), "DiagnosticGroup<visibility>(WARNING), DiagnosticGroup<visibility>(OFF)")
}

// Java: WarningsGuardTest.findNameNode(Node, String)
func jsaWgFindNameNode(root rhino.Node, name string) rhino.Node {
	if root.IsName() && root.GetString() == name {
		return root
	}

	for n := root.GetFirstChild(); !n.IsNil(); n = n.GetNext() {
		result := jsaWgFindNameNode(n, name)
		if !result.IsNil() {
			return result
		}
	}
	return rhino.NilNode
}

// Java: WarningsGuardTest.findGetPropNode(Node, String)
func jsaWgFindGetPropNode(root rhino.Node, name string) rhino.Node {
	if root.IsGetProp() && root.GetString() == name {
		return root
	}

	for n := root.GetFirstChild(); !n.IsNil(); n = n.GetNext() {
		result := jsaWgFindGetPropNode(n, name)
		if !result.IsNil() {
			return result
		}
	}
	return rhino.NilNode
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
	jsaWgLevelIs(t, ts, jsaWgMakeErrorLevel(jsaStr("a.i.js"), CheckLevel_ERROR), CheckLevel_OFF)
	jsaWgNoLevel(t, ts, jsaWgMakeErrorLevel(jsaStr("a.js"), CheckLevel_ERROR))
	jsaWgNoLevel(t, ts, jsaWgMakeErrorType("a.i.js", ClosurePrimitiveErrors_DUPLICATE_NAMESPACE_AND_MODULE))
	jsaEq(t, ts.GetPriority(), 40)

	ps := NewDiagnosticGroupPathSuppressingWarningsGuard(jsaWgDEPRECATED, "gen/")
	jsaWgLevelIs(t, ps, jsaWgMakeErrorType("x/gen/a.js", jsaWgDEPRECATED_NAME), CheckLevel_OFF)
	jsaWgNoLevel(t, ps, jsaWgMakeErrorType("x/a.js", jsaWgDEPRECATED_NAME))
	jsaWgNoLevel(t, ps, jsaWgMakeErrorType("x/gen/a.js", jsaWgBAR_WARNING))
	jsaEq(t, ps.String(), "DiagnosticGroup<deprecated>(gen/)")
	jsaEq(t, ps.GetPriority(), 50)

	cu := NewClosureUnawareCodeWarningsGuard()
	n := rhino.NewNode(rhino.Token_EMPTY)
	jsaWgNoLevel(t, cu, JSError_MakeWithNode(n, RhinoErrorReporter_TYPE_PARSE_ERROR))
	jsaWgNoLevel(t, cu, JSError_Make(RhinoErrorReporter_TYPE_PARSE_ERROR))
	jsaEq(t, cu.GetPriority(), 1)
}
