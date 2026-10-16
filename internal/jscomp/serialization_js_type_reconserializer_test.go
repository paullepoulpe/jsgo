// Port of com/google/javascript/jscomp/serialization/JSTypeReconserializerTest.java @ closure-compiler v20260415

package jscomp

import "testing"

// Every JSTypeReconserializerTest case is a CompilerTestCase: it type-checks
// source code with Compiler and serializes the JSTypes on the AST. Compiler,
// CompilerTestCase and NodeTraversal are not ported yet, so the cases are
// skipped. (JSTypeReconserializer is checked against Java on a whole type
// universe by TestJstypeSerializationOracle meanwhile.)

func TestJSTypeReconserializer_NativeTypesAreNotSerialized(t *testing.T) {
	t.Skip("needs Compiler, CompilerTestCase, NodeTraversal")
}

func TestJSTypeReconserializer_ObjectTypeSerializationFormat(t *testing.T) {
	t.Skip("needs Compiler, CompilerTestCase, NodeTraversal")
}

func TestJSTypeReconserializer_ObjectTypeSerializationFormat_canAvoidSerializingProperty(t *testing.T) {
	t.Skip("needs Compiler, CompilerTestCase, NodeTraversal")
}

func TestJSTypeReconserializer_DisambiguationEdges_pointFromInstanceToPrototype(t *testing.T) {
	t.Skip("needs Compiler, CompilerTestCase, NodeTraversal")
}

func TestJSTypeReconserializer_DisambiguationEdges_pointFromInterfaceToPrototype(t *testing.T) {
	t.Skip("needs Compiler, CompilerTestCase, NodeTraversal")
}

func TestJSTypeReconserializer_DisambiguationEdges_pointFromInstanceToInterface(t *testing.T) {
	t.Skip("needs Compiler, CompilerTestCase, NodeTraversal")
}

func TestJSTypeReconserializer_DisambiguationEdges_pointFromSubctorToSuperctor(t *testing.T) {
	t.Skip("needs Compiler, CompilerTestCase, NodeTraversal")
}

func TestJSTypeReconserializer_DisambiguationEdges_dont_pointFromPrototypeToInstance(t *testing.T) {
	t.Skip("needs Compiler, CompilerTestCase, NodeTraversal")
}

func TestJSTypeReconserializer_MultipleAnonymousTypesDoesntCrash(t *testing.T) {
	t.Skip("needs Compiler, CompilerTestCase, NodeTraversal")
}

func TestJSTypeReconserializer_Union_proxyLikeTypes_dontPreventUnionCollapsing(t *testing.T) {
	t.Skip("needs Compiler, CompilerTestCase, NodeTraversal")
}

func TestJSTypeReconserializer_SerializedInstanceTypeHasClassName(t *testing.T) {
	t.Skip("needs Compiler, CompilerTestCase, NodeTraversal")
}

func TestJSTypeReconserializer_EnumElementSerializesToPrimitive(t *testing.T) {
	t.Skip("needs Compiler, CompilerTestCase, NodeTraversal")
}

func TestJSTypeReconserializer_BoxableScalarTypesKeepOriginalName(t *testing.T) {
	t.Skip("needs Compiler, CompilerTestCase, NodeTraversal")
}

func TestJSTypeReconserializer_MarksClosureAssertions(t *testing.T) {
	t.Skip("needs Compiler, CompilerTestCase, NodeTraversal")
}

func TestJSTypeReconserializer_SerializeObjectLiteralTypesAsInvalidating(t *testing.T) {
	t.Skip("needs Compiler, CompilerTestCase, NodeTraversal")
}

func TestJSTypeReconserializer_SerializesTemplateTypesAsUnknown(t *testing.T) {
	t.Skip("needs Compiler, CompilerTestCase, NodeTraversal")
}

func TestJSTypeReconserializer_Reconcile_literalFunctions_toNativeType(t *testing.T) {
	t.Skip("needs Compiler, CompilerTestCase, NodeTraversal")
}

func TestJSTypeReconserializer_Discards_differentTemplatizationsOfSameRawType(t *testing.T) {
	t.Skip("needs Compiler, CompilerTestCase, NodeTraversal")
}

func TestJSTypeReconserializer_DoesntCrashOnFunctionWithNonObjectThisType(t *testing.T) {
	t.Skip("needs Compiler, CompilerTestCase, NodeTraversal")
}

func TestJSTypeReconserializer_SupertypeListForInterface_includesIllegalExtendedTypes(t *testing.T) {
	t.Skip("needs Compiler, CompilerTestCase, NodeTraversal")
}

func TestJSTypeReconserializer_DoesntSerializeNullOrUndefinedUnion(t *testing.T) {
	t.Skip("needs Compiler, CompilerTestCase, NodeTraversal")
}

func TestJSTypeReconserializer_DoesntSerializeForwardDeclaredType(t *testing.T) {
	t.Skip("needs Compiler, CompilerTestCase, NodeTraversal")
}

func TestJSTypeReconserializer_DoesntSerializeMissingType(t *testing.T) {
	t.Skip("needs Compiler, CompilerTestCase, NodeTraversal")
}

func TestJSTypeReconserializer_SerializesSourceOfSyntheticCodeDoesntCrash(t *testing.T) {
	t.Skip("needs Compiler, CompilerTestCase, NodeTraversal")
}

func TestJSTypeReconserializer_Reconcile_setInstanceColors(t *testing.T) {
	t.Skip("needs Compiler, CompilerTestCase, NodeTraversal")
}

func TestJSTypeReconserializer_Reconcile_addPrototypes_noActualTestcase(t *testing.T) {
	t.Skip("needs Compiler, CompilerTestCase, NodeTraversal")
}

func TestJSTypeReconserializer_Reconcile_setOwnProperties(t *testing.T) {
	t.Skip("needs Compiler, CompilerTestCase, NodeTraversal")
}

func TestJSTypeReconserializer_Reconcile_setClosureAssert(t *testing.T) {
	t.Skip("needs Compiler, CompilerTestCase, NodeTraversal")
}

func TestJSTypeReconserializer_Reconcile_setConstructor_noActualTestcase(t *testing.T) {
	t.Skip("needs Compiler, CompilerTestCase, NodeTraversal")
}

func TestJSTypeReconserializer_Reconcile_setInvalidating(t *testing.T) {
	t.Skip("needs Compiler, CompilerTestCase, NodeTraversal")
}

func TestJSTypeReconserializer_Reconcile_setPropertiesKeepOriginalName_noActualTestcase(t *testing.T) {
	t.Skip("needs Compiler, CompilerTestCase, NodeTraversal")
}

func TestJSTypeReconserializer_DebugInfo_mismatchLocations(t *testing.T) {
	t.Skip("needs Compiler, CompilerTestCase, NodeTraversal")
}
