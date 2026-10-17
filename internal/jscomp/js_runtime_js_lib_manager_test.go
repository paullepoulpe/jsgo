// Port of com/google/javascript/jscomp/js/RuntimeJsLibManagerTest.java @ closure-compiler v20260415

package jscomp

import (
	"slices"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// stubResourceLoader is RuntimeJsLibManagerTest.StubResourceLoader.
func stubResourceLoader(t *testing.T, resources map[string]string) RuntimeJsLibManager_ResourceProvider {
	return func(resourceName string, unused string) *rhino.Node {
		code, ok := resources[resourceName]
		if !ok {
			return nil
		}
		compiler := NewCompiler()
		return compiler.ParseSyntheticCode("lib.js", code)
	}
}

// runtimeJsLibManagerTestJs is RuntimeJsLibManagerTest.js(String).
func runtimeJsLibManagerTestJs(t *testing.T, code string) *rhino.Node {
	compiler := NewCompiler()
	return compiler.ParseSyntheticCode("lib.js", code)
}

var runtimeJsLibManagerTest_JSCOMP_INHERITS = rhino.IR_Getprop(rhino.IR_Name("$jscomp"), "inherits")

func scriptSupplier(script *rhino.Node) func() *rhino.Node {
	return func() *rhino.Node { return script }
}

func freshScript() *rhino.Node { return rhino.IR_Script() }

func TestRuntimeJsLibManager_NoOpModeDoesNotRecordOrInject(t *testing.T) {
	script := rhino.IR_Script()
	manager := RuntimeJsLibManager_Create(
		RuntimeJsLibManager_RuntimeLibraryMode_NO_OP,
		stubResourceLoader(t, map[string]string{"lib": "0"}),
		NewChangeTracker(),
		scriptSupplier(script))

	manager.EnsureLibraryInjected("lib" /* force= */, false)

	assertEqual(t, "injected", len(manager.GetInjectedLibraries()), 0)
}

func TestRuntimeJsLibManager_NoOpModeDoesRecordAndInjectWithForceInjection(t *testing.T) {
	script := rhino.IR_Script()
	manager := RuntimeJsLibManager_Create(
		RuntimeJsLibManager_RuntimeLibraryMode_NO_OP,
		stubResourceLoader(t, map[string]string{"lib": "0"}),
		NewChangeTracker(),
		scriptSupplier(script))

	manager.EnsureLibraryInjected("lib" /* force= */, true)

	assertSliceEqual(t, "injected", manager.GetInjectedLibraries(), []string{"lib"})
}

func TestRuntimeJsLibManager_RecordOnlyModeDoesNotInjectButRecordsLibraries(t *testing.T) {
	manager := RuntimeJsLibManager_Create(
		RuntimeJsLibManager_RuntimeLibraryMode_RECORD_ONLY,
		stubResourceLoader(t, map[string]string{"lib": "0"}),
		NewChangeTracker(),
		freshScript)

	injected := manager.EnsureLibraryInjected("lib" /* force= */, false)

	if injected != nil {
		t.Errorf("injected %v, want null", injected)
	}
	assertSliceEqual(t, "injected", manager.GetInjectedLibraries(), []string{"lib"})
}

func TestRuntimeJsLibManager_RecordOnlyModeDoesRecordAndInjectWithForceInjection(t *testing.T) {
	manager := RuntimeJsLibManager_Create(
		RuntimeJsLibManager_RuntimeLibraryMode_RECORD_ONLY,
		stubResourceLoader(t, map[string]string{"lib": "0"}),
		NewChangeTracker(),
		freshScript)

	manager.EnsureLibraryInjected("lib" /* force= */, true)

	assertSliceEqual(t, "injected", manager.GetInjectedLibraries(), []string{"lib"})
}

func TestRuntimeJsLibManager_InjectModeDoesRecordAndInjectAndReturnsLastStatementInLibrary(t *testing.T) {
	script := rhino.IR_Script()
	manager := RuntimeJsLibManager_Create(
		RuntimeJsLibManager_RuntimeLibraryMode_INJECT,
		stubResourceLoader(t, map[string]string{"lib": "function foo() {}\n0;\n"}),
		NewChangeTracker(),
		scriptSupplier(script))

	injected := manager.EnsureLibraryInjected("lib" /* force= */, false)

	assertSliceEqual(t, "injected", manager.GetInjectedLibraries(), []string{"lib"})
	assertNodeEquivalent(t, script, runtimeJsLibManagerTestJs(t, "function foo() {} 0;"))
	assertEqual(t, "second child", script.GetSecondChild(), injected)
}

func TestRuntimeJsLibManager_DoesNotInjectSameLibraryTwice(t *testing.T) {
	script := rhino.IR_Script()
	manager := RuntimeJsLibManager_Create(
		RuntimeJsLibManager_RuntimeLibraryMode_INJECT,
		stubResourceLoader(t, map[string]string{"lib": "0"}),
		NewChangeTracker(),
		scriptSupplier(script))

	injected1 := manager.EnsureLibraryInjected("lib" /* force= */, false)
	injected2 := manager.EnsureLibraryInjected("lib" /* force= */, false)

	assertEqual(t, "same", injected1, injected2)
	assertEqual(t, "one child", script.HasOneChild(), true)
}

func TestRuntimeJsLibManager_InjectsBeforeExistingCode(t *testing.T) {
	existing := rhino.IR_Block()
	script := rhino.IR_ScriptOfNodes(existing)

	manager := RuntimeJsLibManager_Create(
		RuntimeJsLibManager_RuntimeLibraryMode_INJECT,
		stubResourceLoader(t, map[string]string{"lib": "0"}),
		NewChangeTracker(),
		scriptSupplier(script))

	injected := manager.EnsureLibraryInjected("lib" /* force= */, false)

	assertNodeEquivalent(t, script.GetFirstChild(), injected)
	assertEqual(t, "second child", script.GetSecondChild(), existing)
}

func TestRuntimeJsLibManager_InjectsAfterPreviouslyInjectedLib(t *testing.T) {
	script := rhino.IR_Script()
	manager := RuntimeJsLibManager_Create(
		RuntimeJsLibManager_RuntimeLibraryMode_INJECT,
		stubResourceLoader(t, map[string]string{"lib0": "0", "lib1": "1"}),
		NewChangeTracker(),
		scriptSupplier(script))

	injected0 := manager.EnsureLibraryInjected("lib0" /* force= */, false)
	injected1 := manager.EnsureLibraryInjected("lib1" /* force= */, false)

	assertEqual(t, "first child", script.GetFirstChild(), injected0)
	assertEqual(t, "second child", script.GetSecondChild(), injected1)
	assertNodeEquivalent(t, script, runtimeJsLibManagerTestJs(t, "0; 1;"))
}

func TestRuntimeJsLibManager_InjectsRequiredLibsOnlyOnce(t *testing.T) {
	script := rhino.IR_Script()
	// Create a dependency tree:
	//       (base)
	//      /      \
	//  (childA)   (childB)
	resources := stubResourceLoader(t, map[string]string{
		"base.js":   "print('in BASE');",
		"childA.js": "'require base.js';\nprint('in CHILD_A');\n",
		"childB.js": "'require base.js';\nprint('in CHILD_B');\n",
	})
	manager := RuntimeJsLibManager_Create(
		RuntimeJsLibManager_RuntimeLibraryMode_INJECT, resources, NewChangeTracker(), scriptSupplier(script))

	childA := manager.EnsureLibraryInjected("childA.js" /* force= */, false)
	childB := manager.EnsureLibraryInjected("childB.js" /* force= */, false)

	// containsExactly without inOrder(): the libraries are recorded as they
	// are reached (childA.js before the base.js it requires).
	assertSliceEqual(t, "injected", slices.Sorted(slices.Values(manager.GetInjectedLibraries())),
		[]string{"base.js", "childA.js", "childB.js"})
	assertEqual(t, "second child", script.GetSecondChild(), childA)
	assertEqual(t, "last child", script.GetLastChild(), childB)
	assertNodeEquivalent(t, script, runtimeJsLibManagerTestJs(t,
		"print('in BASE');\nprint('in CHILD_A');\nprint('in CHILD_B');\n"))
}

func TestRuntimeJsLibManager_InjectModeDoesNotPassAlongUseStrictDirective(t *testing.T) {
	script := rhino.IR_Script()
	manager := RuntimeJsLibManager_Create(
		RuntimeJsLibManager_RuntimeLibraryMode_INJECT,
		stubResourceLoader(t, map[string]string{"lib": "'use strict';\n0;\n"}),
		NewChangeTracker(),
		scriptSupplier(script))

	injected := manager.EnsureLibraryInjected("lib" /* force= */, false)

	assertNodeEquivalent(t, injected, runtimeJsLibManagerTestJs(t, "0;").GetOnlyChild())
	assertEqual(t, "one child", script.HasOneChild(), true)
	assertSliceEqual(t, "injected", manager.GetInjectedLibraries(), []string{"lib"})
}

func TestRuntimeJsLibManager_InjectingMissingLibThrowsException(t *testing.T) {
	manager := RuntimeJsLibManager_Create(
		RuntimeJsLibManager_RuntimeLibraryMode_INJECT,
		stubResourceLoader(t, map[string]string{}),
		NewChangeTracker(),
		freshScript)

	// Java: NullPointerException, from dereferencing the null AST. In Go
	// that is a nil dereference, which jutil reports as a RuntimeException.
	assertThrows[jutil.IRuntimeException](t, func() {
		manager.EnsureLibraryInjected("lib" /* force= */, false)
	})
}

func TestRuntimeJsLibManager_InjectLibForFieldAddsCorrespondingLibrary(t *testing.T) {
	manager := RuntimeJsLibManager_Create(
		RuntimeJsLibManager_RuntimeLibraryMode_RECORD_ONLY,
		stubResourceLoader(t, map[string]string{}),
		NewChangeTracker(),
		freshScript)

	manager.InjectLibForField("$jscomp.inherits")

	assertSliceEqual(t, "injected", manager.GetInjectedLibraries(), []string{"es6/util/inherits"})
}

func TestRuntimeJsLibManager_InjectLibForFieldCrashesIfNonQualifiedName(t *testing.T) {
	manager := RuntimeJsLibManager_Create(
		RuntimeJsLibManager_RuntimeLibraryMode_RECORD_ONLY,
		stubResourceLoader(t, map[string]string{}),
		NewChangeTracker(),
		freshScript)

	assertThrows[jutil.IIllegalArgumentException](t, func() { manager.InjectLibForField("3 + 4") })
}

func TestRuntimeJsLibManager_InjectLibForFieldCrashesIfPassedNonJscompName(t *testing.T) {
	manager := RuntimeJsLibManager_Create(
		RuntimeJsLibManager_RuntimeLibraryMode_RECORD_ONLY,
		stubResourceLoader(t, map[string]string{}),
		NewChangeTracker(),
		freshScript)

	assertThrows[jutil.IIllegalArgumentException](t, func() { manager.InjectLibForField("foobar") })
}

func TestRuntimeJsLibManager_InjectLibForFieldCrashesIfPassedJustJscomp(t *testing.T) {
	manager := RuntimeJsLibManager_Create(
		RuntimeJsLibManager_RuntimeLibraryMode_RECORD_ONLY,
		stubResourceLoader(t, map[string]string{}),
		NewChangeTracker(),
		freshScript)

	assertThrows[jutil.IIllegalArgumentException](t, func() { manager.InjectLibForField("$jscomp") })
}

func TestRuntimeJsLibManager_InjectLibForFieldCrashesIfPassedMissingProperty(t *testing.T) {
	manager := RuntimeJsLibManager_Create(
		RuntimeJsLibManager_RuntimeLibraryMode_RECORD_ONLY,
		stubResourceLoader(t, map[string]string{}),
		NewChangeTracker(),
		freshScript)

	ex := assertThrows[jutil.INullPointerException](t, func() { manager.InjectLibForField("$jscomp.doesNotExist") })
	if !strings.Contains(ex.GetMessage(), "Cannot find definition of $jscomp.doesNotExist") {
		t.Errorf("message %q", ex.GetMessage())
	}
}

func TestRuntimeJsLibManager_InjectLibForFieldThenAssertInjectedAndGetQualifiedNamePasses(t *testing.T) {
	manager := RuntimeJsLibManager_Create(
		RuntimeJsLibManager_RuntimeLibraryMode_RECORD_ONLY,
		stubResourceLoader(t, map[string]string{}),
		NewChangeTracker(),
		freshScript)

	manager.InjectLibForField("$jscomp.inherits")
	field := manager.GetJsLibField("$jscomp.inherits")

	assertEqual(t, "qualified name", field.AssertInjected().QualifiedName(), "$jscomp.inherits")
}

func TestRuntimeJsLibManager_NoinjectLibForFieldThenAssertInjectedRecordAndValidateFieldsModeFails(t *testing.T) {
	manager := RuntimeJsLibManager_Create(
		RuntimeJsLibManager_RuntimeLibraryMode_RECORD_AND_VALIDATE_FIELDS,
		stubResourceLoader(t, map[string]string{}),
		NewChangeTracker(),
		freshScript)

	field := manager.GetJsLibField("$jscomp.inherits")

	assertThrows[jutil.IIllegalStateException](t, func() { field.AssertInjected() })
}

func TestRuntimeJsLibManager_NoinjectLibForFieldThenAssertInjectedInjectModeFails(t *testing.T) {
	manager := RuntimeJsLibManager_Create(
		RuntimeJsLibManager_RuntimeLibraryMode_RECORD_AND_VALIDATE_FIELDS,
		stubResourceLoader(t, map[string]string{}),
		NewChangeTracker(),
		freshScript)

	field := manager.GetJsLibField("$jscomp.inherits")

	assertThrows[jutil.IIllegalStateException](t, func() { field.AssertInjected() })
}

func TestRuntimeJsLibManager_NoinjectLibForFieldThenAssertInjectedAndGetQualifiedNameRecordOnlyModeSucceeds(t *testing.T) {
	manager := RuntimeJsLibManager_Create(
		RuntimeJsLibManager_RuntimeLibraryMode_RECORD_ONLY,
		stubResourceLoader(t, map[string]string{}),
		NewChangeTracker(),
		freshScript)

	field := manager.GetJsLibField("$jscomp.inherits")

	_ = field.AssertInjected().QualifiedName()
}

func TestRuntimeJsLibManager_NoinjectLibForFieldThenAssertInjectedAndGetQualifiedNameNoOpModeSucceeds(t *testing.T) {
	manager := RuntimeJsLibManager_Create(
		RuntimeJsLibManager_RuntimeLibraryMode_NO_OP,
		stubResourceLoader(t, map[string]string{}),
		NewChangeTracker(),
		freshScript)

	field := manager.GetJsLibField("$jscomp.inherits")

	_ = field.AssertInjected().QualifiedName()
}

func TestRuntimeJsLibManager_FieldWithoutInjectionAllowsCallingMatches(t *testing.T) {
	manager := RuntimeJsLibManager_Create(
		RuntimeJsLibManager_RuntimeLibraryMode_RECORD_AND_VALIDATE_FIELDS,
		stubResourceLoader(t, map[string]string{}),
		NewChangeTracker(),
		freshScript)

	field := manager.GetJsLibField("$jscomp.inherits")

	assertEqual(t, "foo", field.Matches(rhino.IR_Name("foo")), false)
	assertEqual(t, "$jscomp.inherits", field.Matches(runtimeJsLibManagerTest_JSCOMP_INHERITS), true)
}

func TestRuntimeJsLibManager_ExternFieldNamesModeDoesNotInjectButRecordsLibraries(t *testing.T) {
	manager := RuntimeJsLibManager_Create(
		RuntimeJsLibManager_RuntimeLibraryMode_EXTERN_FIELD_NAMES,
		stubResourceLoader(t, map[string]string{"lib": "0"}),
		NewChangeTracker(),
		freshScript)

	injected := manager.EnsureLibraryInjected("lib" /* force= */, false)

	if injected != nil {
		t.Errorf("injected %v, want null", injected)
	}
	assertSliceEqual(t, "injected", manager.GetInjectedLibraries(), []string{"lib"})
}

func TestRuntimeJsLibManager_ExternFieldNamesModeDoesRecordAndInjectWithForceInjection(t *testing.T) {
	script := rhino.IR_Script()
	manager := RuntimeJsLibManager_Create(
		RuntimeJsLibManager_RuntimeLibraryMode_EXTERN_FIELD_NAMES,
		stubResourceLoader(t, map[string]string{"lib": "0"}),
		NewChangeTracker(),
		scriptSupplier(script))

	manager.EnsureLibraryInjected("lib" /* force= */, true)

	assertSliceEqual(t, "injected", manager.GetInjectedLibraries(), []string{"lib"})
	assertNodeEquivalent(t, script, runtimeJsLibManagerTestJs(t, "0;"))
}

func TestRuntimeJsLibManager_ExternFieldNamesModeInjectLibForFieldRecordsAndAddsExternDecl(t *testing.T) {
	externs := rhino.IR_Script()
	externs.SetStaticSourceFile(SourceFile_FromCodeWithKind("ex.js", "", rhino.StaticSourceFile_SourceKind_EXTERN))
	manager := RuntimeJsLibManager_Create(
		RuntimeJsLibManager_RuntimeLibraryMode_EXTERN_FIELD_NAMES,
		stubResourceLoader(t, map[string]string{}),
		NewChangeTracker(),
		scriptSupplier(externs))

	manager.InjectLibForField("$jscomp.inherits")

	assertSliceEqual(t, "injected", manager.GetInjectedLibraries(), []string{"es6/util/inherits"})
	externedFields := manager.GetExternedFields()
	assertEqual(t, "externed fields", len(externedFields), 1)
	assertEqual(t, "uncompiledName", externedFields[0].UncompiledName(), "$jscomp.inherits")
	assertEqual(t, "qualifiedName", externedFields[0].QualifiedName(), "$jscomp_inherits")
	assertNodeEquivalent(t, externs, runtimeJsLibManagerTestJs(t, "var $jscomp_inherits;"))
	assertEqual(t, "IS_CONSTANT_NAME",
		externs.GetOnlyChild().GetOnlyChild().GetBooleanProp(rhino.Node_IS_CONSTANT_NAME), true)
}

func TestRuntimeJsLibManager_ExternFieldNamesModeInjectLibForFieldTwiceDoesNotAddDuplicateExterns(t *testing.T) {
	externs := rhino.IR_Script()
	externs.SetStaticSourceFile(SourceFile_FromCodeWithKind("ex.js", "", rhino.StaticSourceFile_SourceKind_EXTERN))
	manager := RuntimeJsLibManager_Create(
		RuntimeJsLibManager_RuntimeLibraryMode_EXTERN_FIELD_NAMES,
		stubResourceLoader(t, map[string]string{}),
		NewChangeTracker(),
		scriptSupplier(externs))

	manager.InjectLibForField("$jscomp.inherits")
	manager.InjectLibForField("$jscomp.inherits")

	assertNodeEquivalent(t, externs, runtimeJsLibManagerTestJs(t, "var $jscomp_inherits;"))
}

func TestRuntimeJsLibManager_ExternFieldNamesModeNoinjectLibForFieldThenAssertInjectedFails(t *testing.T) {
	manager := RuntimeJsLibManager_Create(
		RuntimeJsLibManager_RuntimeLibraryMode_EXTERN_FIELD_NAMES,
		stubResourceLoader(t, map[string]string{}),
		NewChangeTracker(),
		freshScript)

	field := manager.GetJsLibField("$jscomp.inherits")

	assertThrows[jutil.IIllegalStateException](t, func() { field.AssertInjected() })
}
