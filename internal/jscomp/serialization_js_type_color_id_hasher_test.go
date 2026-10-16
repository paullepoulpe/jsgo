// Port of com/google/javascript/jscomp/serialization/JSTypeColorIdHasherTest.java @ closure-compiler v20260415

package jscomp

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/rhino/rhinotesting"
)

// The Java test is parameterized on rewriteClosureModules (true, false).
// Every case but returnsNativeId_forNativeTypes is a CompilerTestCase that
// type-checks source code and hashes the types of labeled expressions;
// Compiler, CompilerTestCase and NodeTraversal are not ported yet, so those
// cases are skipped. (JSTypeColorIdHasher is checked against Java on a whole
// type universe by TestJstypeDepsOracle meanwhile.)

func TestJSTypeColorIdHasher_HaveDifferentIds_ctorInstancePrototype(t *testing.T) {
	t.Skip("needs Compiler, CompilerTestCase, NodeTraversal")
}

func TestJSTypeColorIdHasher_HaveDifferentIds_qnamesWithSameEnding(t *testing.T) {
	t.Skip("needs Compiler, CompilerTestCase, NodeTraversal")
}

func TestJSTypeColorIdHasher_HaveDifferentIds_classWithSameName_inDifferentGoogModules(t *testing.T) {
	t.Skip("needs Compiler, CompilerTestCase, NodeTraversal")
}

func TestJSTypeColorIdHasher_HaveDifferentIds_classWithSameName_inDifferentGoogModules_inLoadModuleFile(t *testing.T) {
	t.Skip("needs Compiler, CompilerTestCase, NodeTraversal")
}

func TestJSTypeColorIdHasher_HaveDifferentIds_classWithSameName_insideAndOutsideGoogModule(t *testing.T) {
	t.Skip("needs Compiler, CompilerTestCase, NodeTraversal")
}

func TestJSTypeColorIdHasher_HaveDifferentIds_anonymousTypes_withDifferentPropertyNames(t *testing.T) {
	t.Skip("needs Compiler, CompilerTestCase, NodeTraversal")
}

func TestJSTypeColorIdHasher_HaveSameIds_anonymousTypes_withSamePropertyNames_andDiffererntPropertyTypes(t *testing.T) {
	t.Skip("needs Compiler, CompilerTestCase, NodeTraversal")
}

func TestJSTypeColorIdHasher_HaveSameIds_nominalTypes_fromDifferentScopes_fromSameGoogModule(t *testing.T) {
	t.Skip("needs Compiler, CompilerTestCase, NodeTraversal")
}

func TestJSTypeColorIdHasher_HaveRecompileStableIds_nominalTypes(t *testing.T) {
	t.Skip("needs Compiler, CompilerTestCase, NodeTraversal")
}

func TestJSTypeColorIdHasher_HaveRecompileStableIds_anonymousTypes(t *testing.T) {
	t.Skip("needs Compiler, CompilerTestCase, NodeTraversal")
}

func TestJSTypeColorIdHasher_HaveHeaderDerivableIds_nominalTypes(t *testing.T) {
	t.Skip("needs Compiler, CompilerTestCase, NodeTraversal")
}

func TestJSTypeColorIdHasher_HaveHeaderDerivableIds_anonymousTypes(t *testing.T) {
	t.Skip("needs Compiler, CompilerTestCase, NodeTraversal")
}

func TestJSTypeColorIdHasher_ReturnsNativeId_forNativeTypes(t *testing.T) {
	for _, rewriteClosureModules := range []bool{true, false} {
		_ = rewriteClosureModules // the parameter only affects CompilerTestCase setup
		truth := rhinotesting.NewTruth(t)
		registry := rhino.NewJSTypeRegistry(rhino.ErrorReporter_ALWAYS_THROWS_INSTANCE)
		hasher := NewJSTypeColorIdHasher(registry)

		for n, i := range JSTypeColorIdHasher_NATIVE_TYPE_TO_ID.All() {
			typ := registry.GetNativeObjectType(n)
			truth.AssertThat(hasher.HashObjectType(typ)).IsEqualTo(i)
		}
	}
}
