// Port of com/google/javascript/rhino/JSDocInfoTest.java @ closure-compiler v20261028

/*
 *
 * ***** BEGIN LICENSE BLOCK *****
 * Version: MPL 1.1
 *
 * The contents of this file are subject to the Mozilla Public License Version
 * 1.1 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 * http://www.mozilla.org/MPL/
 *
 * Software distributed under the License is distributed on an "AS IS" basis,
 * WITHOUT WARRANTY OF ANY KIND, either express or implied. See the License
 * for the specific language governing rights and limitations under the
 * License.
 *
 * The Original Code is Rhino code, released
 * May 6, 1999.
 *
 * The Initial Developer of the Original Code is
 * Netscape Communications Corporation.
 * Portions created by the Initial Developer are Copyright (C) 1997-1999
 * the Initial Developer. All Rights Reserved.
 *
 * Contributor(s):
 *   Nick Santos
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */

package rhino

// The Java test class owns a JSTypeRegistry (with a TestErrorReporter
// verified after each test) for its resolve(JSTypeExpression) helper. The
// tests that resolve type expressions are in js_doc_info_jstype_test.go
// (package rhino_test, as they use rhinotesting); the ones here never touch
// the registry, so the @After check holds trivially.

import (
	"slices"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

func assertNullTypeExpr(t *testing.T, what string, e *JSTypeExpression) {
	t.Helper()
	if e != nil {
		t.Errorf("%s = %v, want null", what, e)
	}
}

func assertNotContains(t *testing.T, what string, set *jutil.OrderedSet[string], s string) {
	t.Helper()
	if set.Contains(s) {
		t.Errorf("%s %v contains %q", what, set, s)
	}
}

func assertContains(t *testing.T, what string, set *jutil.OrderedSet[string], s string) {
	t.Helper()
	if !set.Contains(s) {
		t.Errorf("%s %v does not contain %q", what, set, s)
	}
}

// Tests the assigned ordinal of the elements of the JSDocInfo.Visibility
// enum.
func TestJSDocInfo_VisibilityOrdinal(t *testing.T) {
	for want, v := range []JSDocInfo_Visibility{
		JSDocInfo_Visibility_PRIVATE,
		JSDocInfo_Visibility_PACKAGE,
		JSDocInfo_Visibility_PROTECTED,
		JSDocInfo_Visibility_PUBLIC,
	} {
		if got := v.Ordinal(); got != want {
			t.Errorf("%v.ordinal() = %d, want %d", v, got, want)
		}
	}
}

// Tests that all module local names get correctly removed from a
// JSTypeExpression
func TestJSDocInfo_RemovesModuleLocalNames(t *testing.T) {
	jsTypeExpression := createSampleTypeExpression()
	mockedModuleLocals := jutil.NewOrderedSet("Item", "AnotherItem")

	newExpr := jsTypeExpression.ReplaceNamesWithUnknownType(mockedModuleLocals)
	replacedNames := newExpr.GetAllTypeNames()
	assertNotContains(t, "replacedNames", replacedNames, "Item")
	assertContains(t, "replacedNames", replacedNames, "string")
	assertContains(t, "replacedNames", replacedNames, "boolean")
	assertNotContains(t, "replacedNames", replacedNames, "AnotherItem")
}

func TestJSDocInfo_SetConstant(t *testing.T) {
	builder := JSDocInfo_Builder_()
	builder.RecordConstancy()
	info := builder.Build()

	assertFalse(t, "hasType", info.HasType())
	assertTrue(t, "isConstant", info.IsConstant())
	assertFalse(t, "isConstructor", info.IsConstructor())
	assertFalse(t, "isDefine", info.IsDefine())
}

func TestJSDocInfo_SetConstructor(t *testing.T) {
	builder := JSDocInfo_Builder_()
	builder.RecordConstructor()
	info := builder.Build()

	assertFalse(t, "isConstant", info.IsConstant())
	assertTrue(t, "isConstructor", info.IsConstructor())
	assertFalse(t, "isDefine", info.IsDefine())
}

func TestJSDocInfo_SetDefine(t *testing.T) {
	builder := JSDocInfo_Builder_()
	builder.RecordDefineType(jsdocFromString("string"))
	info := builder.Build()

	assertTrue(t, "isConstant", info.IsConstant())
	assertFalse(t, "isConstructor", info.IsConstructor())
	assertTrue(t, "isDefine", info.IsDefine())
}

func TestJSDocInfo_SetTypeSummary(t *testing.T) {
	builder := JSDocInfo_Builder_()
	builder.RecordTypeSummary()
	info := builder.Build()

	assertTrue(t, "isTypeSummary", info.IsTypeSummary())
}

func TestJSDocInfo_SetOverride(t *testing.T) {
	builder := JSDocInfo_Builder_()
	builder.RecordOverride()
	info := builder.Build()

	assertFalse(t, "isDeprecated", info.IsDeprecated())
	assertTrue(t, "isOverride", info.IsOverride())
}

func TestJSDocInfo_SetExport(t *testing.T) {
	builder := JSDocInfo_Builder_()
	builder.RecordExport()
	info := builder.Build()

	assertTrue(t, "isExport", info.IsExport())
}

func TestJSDocInfo_SetPolymerBehavior(t *testing.T) {
	assertFalse(t, "EMPTY.isPolymerBehavior", jsdocTestEMPTY.IsPolymerBehavior())

	builder := JSDocInfo_Builder_()
	builder.RecordPolymerBehavior()
	info := builder.Build()
	assertTrue(t, "isPolymerBehavior", info.IsPolymerBehavior())
}

func TestJSDocInfo_SetPolymer(t *testing.T) {
	assertFalse(t, "EMPTY.isPolymer", jsdocTestEMPTY.IsPolymer())

	builder := JSDocInfo_Builder_()
	builder.RecordPolymer()
	info := builder.Build()
	assertTrue(t, "isPolymer", info.IsPolymer())
}

func TestJSDocInfo_SetCustomElement(t *testing.T) {
	assertFalse(t, "EMPTY.isCustomElement", jsdocTestEMPTY.IsCustomElement())

	builder := JSDocInfo_Builder_()
	builder.RecordCustomElement()
	info := builder.Build()
	assertTrue(t, "isCustomElement", info.IsCustomElement())
}

func TestJSDocInfo_SetMixinClass(t *testing.T) {
	assertFalse(t, "EMPTY.isMixinClass", jsdocTestEMPTY.IsMixinClass())

	builder := JSDocInfo_Builder_()
	builder.RecordMixinClass()
	info := builder.Build()
	assertTrue(t, "isMixinClass", info.IsMixinClass())
}

func TestJSDocInfo_SetMixinFunction(t *testing.T) {
	assertFalse(t, "EMPTY.isMixinFunction", jsdocTestEMPTY.IsMixinFunction())

	builder := JSDocInfo_Builder_()
	builder.RecordMixinFunction()
	info := builder.Build()
	assertTrue(t, "isMixinFunction", info.IsMixinFunction())
}

func TestJSDocInfo_SetNoAlias(t *testing.T) {
	assertFalse(t, "EMPTY.isDeprecated", jsdocTestEMPTY.IsDeprecated())
	assertFalse(t, "EMPTY.isOverride", jsdocTestEMPTY.IsOverride())
}

func TestJSDocInfo_SetDeprecated(t *testing.T) {
	builder := JSDocInfo_Builder_()
	builder.RecordDeprecated()
	info := builder.Build()

	assertFalse(t, "isOverride", info.IsOverride())
	assertTrue(t, "isDeprecated", info.IsDeprecated())
}

func TestJSDocInfo_MultipleSetFlags1(t *testing.T) {
	builder := JSDocInfo_Builder_()
	builder.RecordConstancy()
	builder.RecordConstructor()
	builder.RecordNoCollapse()
	info := builder.Build()

	assertFalse(t, "hasType", info.HasType())
	assertTrue(t, "isConstant", info.IsConstant())
	assertTrue(t, "isConstructor", info.IsConstructor())
	assertFalse(t, "isDefine", info.IsDefine())
	assertTrue(t, "isNoCollapse", info.IsNoCollapse())

	builder = info.ToBuilder()
	builder.RecordMutable()
	info = builder.Build()

	assertFalse(t, "isConstant", info.IsConstant())
	assertFalse(t, "isDefine", info.IsDefine())
	assertTrue(t, "isNoCollapse", info.IsNoCollapse())
}

func TestJSDocInfo_DescriptionContainsAtSignCode(t *testing.T) {
	builder := JSDocInfo_Builder_()
	builder.ParseDocumentation()
	builder.RecordOriginalCommentString("Blah blah {@code blah blah} blah blah.")
	info := builder.Build()

	assertTrue(t, "isAtSignCodePresent", info.IsAtSignCodePresent())
}

func TestJSDocInfo_DescriptionDoesNotContainAtSignCode(t *testing.T) {
	builder := JSDocInfo_Builder_()
	builder.ParseDocumentation()
	builder.RecordOriginalCommentString("Blah blah `blah blah` blah blah.")
	info := builder.Build()

	assertFalse(t, "isAtSignCodePresent", info.IsAtSignCodePresent())
}

func TestJSDocInfo_CloneTypeExpressions2(t *testing.T) {
	builder := JSDocInfo_Builder_()
	builder.RecordParameter("a", nil)
	info := builder.Build()

	cloned := info.CloneWithTypeNodes(true)

	assertNullTypeExpr(t, "cloned.getParameterType(a)", cloned.GetParameterType("a"))
}

// Test names in @param get replaced
func TestJSDocInfo_JSDocInfoCloneAndReplaceNames_Params(t *testing.T) {
	builder := JSDocInfo_Builder_()
	jsTypeExpression := createSampleTypeExpression()
	builder.RecordParameter("a", jsTypeExpression)
	info := builder.Build()

	mockedModuleLocals := jutil.NewOrderedSet("Item", "AnotherItem")

	cloned := info.CloneAndReplaceTypeNames(mockedModuleLocals)

	if got := cloned.GetParameterCount(); got != 1 {
		t.Errorf("parameter count = %d, want 1", got)
	}
	if got, ok := cloned.GetParameterNameValueAt(0); !ok || got != "a" {
		t.Errorf("getParameterNameAt(0) = (%q, %v), want (\"a\", true)", got, ok)
	}
	// not same reference after cloning
	assertTrue(t, "cloned param type is a new object", cloned.GetParameterType("a") != jsTypeExpression)
	assertNotContains(t, "param type names", cloned.GetParameterType("a").GetAllTypeNames(), "Item")
	assertNotContains(t, "param type names", cloned.GetParameterType("a").GetAllTypeNames(), "AnotherItem")
}

// Test names in @type get replaced
//
//	          BANG                                        BANG
//	           |                                           |
//	          Item                                         ?
//	        /  |   \                                    /  |   \
//	   QMARK QMARK QMARK          -------->        QMARK QMARK QMARK
//	      /    |    \                                 /    |     \
//	string boolean  AnotherItem                 string   boolean  ?
func TestJSDocInfo_JSDocInfoCloneAndReplaceNames_Type(t *testing.T) {
	builder := JSDocInfo_Builder_()
	jsTypeExpression := createSampleTypeExpression()
	builder.RecordType(jsTypeExpression)
	info := builder.Build()

	mockedModuleLocals := jutil.NewOrderedSet("Item", "AnotherItem")

	cloned := info.CloneAndReplaceTypeNames(mockedModuleLocals)

	assertTrue(t, "cloned type is a new object", cloned.GetType() != jsTypeExpression)
	assertNotContains(t, "type names", cloned.GetType().GetAllTypeNames(), "Item")
	assertNotContains(t, "type names", cloned.GetType().GetAllTypeNames(), "AnotherItem")
}

// Test names in @type get replaced
//
//	          Item                                       ?
//	        /  |   \                                  /  |  \
//	   QMARK QMARK QMARK       --------->        QMARK QMARK QMARK
//	      /    |    \                               /    |    \
//	string boolean  AnotherItem               string   boolean  ?
func TestJSDocInfo_JSDocInfoCloneAndReplaceNames_Type_RootReplacement(t *testing.T) {
	builder := JSDocInfo_Builder_()
	jsTypeExpression := createSampleTypeExpression_rootReplacement()
	builder.RecordType(jsTypeExpression)
	info := builder.Build()

	mockedModuleLocals := jutil.NewOrderedSet("Item", "AnotherItem")

	cloned := info.CloneAndReplaceTypeNames(mockedModuleLocals)

	assertTrue(t, "cloned type is a new object", cloned.GetType() != jsTypeExpression)
	if got := cloned.GetType().GetRoot().GetToken(); got != Token_QMARK {
		t.Errorf("root token = %v, want QMARK", got)
	}
	assertNotContains(t, "type names", cloned.GetType().GetAllTypeNames(), "Item")
	assertNotContains(t, "type names", cloned.GetType().GetAllTypeNames(), "AnotherItem")
}

func TestJSDocInfo_JSDocInfoCloneAndReplaceNames_Type_SingleNode(t *testing.T) {
	builder := JSDocInfo_Builder_()
	root := Node_NewString("Item")
	jsTypeExpression := NewJSTypeExpression(root, "")
	builder.RecordType(jsTypeExpression)
	info := builder.Build()

	mockedModuleLocals := jutil.NewOrderedSet("Item", "AnotherItem")

	cloned := info.CloneAndReplaceTypeNames(mockedModuleLocals)

	assertTrue(t, "cloned type is a new object", cloned.GetType() != jsTypeExpression)
	if got := cloned.GetType().GetRoot().GetToken(); got != Token_QMARK {
		t.Errorf("root token = %v, want QMARK", got)
	}
	assertNotContains(t, "type names", cloned.GetType().GetAllTypeNames(), "Item")
	assertNotContains(t, "type names", cloned.GetType().GetAllTypeNames(), "AnotherItem")
}

func TestJSDocInfo_SetFileOverviewWithDocumentationOff(t *testing.T) {
	builder := JSDocInfo_Builder_()
	builder.RecordFileOverview("hi bob")
	info := builder.BuildWithAlways(true)
	if s, ok := info.GetFileOverviewString(); ok {
		t.Errorf("getFileOverview = %q, want unset", s)
	}
}

func TestJSDocInfo_SetFileOverviewWithDocumentationOn(t *testing.T) {
	builder := JSDocInfo_Builder_()
	builder.ParseDocumentation()
	builder.RecordFileOverview("hi bob")
	info := builder.Build()
	if s, ok := info.GetFileOverviewString(); !ok || s != "hi bob" {
		t.Errorf("getFileOverview = (%q, %v), want (\"hi bob\", true)", s, ok)
	}
}

func TestJSDocInfo_SetSuppressions(t *testing.T) {
	builder := JSDocInfo_Builder_()
	builder.ParseDocumentation()
	builder.RecordSuppressions([]string{"sam", "bob"})
	builder.RecordSuppression("fred")
	info := builder.Build()
	if want := jutil.NewOrderedSet("bob", "sam", "fred"); !info.GetSuppressions().Equals(want) {
		t.Errorf("getSuppressions() = %v, want %v", info.GetSuppressions(), want)
	}
}

func TestJSDocInfo_SetModifies(t *testing.T) {
	builder := JSDocInfo_Builder_()
	builder.ParseDocumentation()
	builder.RecordModifies([]string{"this"})
	info := builder.Build()
	if want := jutil.NewOrderedSet("this"); !info.GetModifies().Equals(want) {
		t.Errorf("getModifies() = %v, want %v", info.GetModifies(), want)
	}
	if got := info.GetModifiesSlice(); !slices.Equal(got, []string{"this"}) {
		t.Errorf("GetModifiesSlice() = %v, want [this]", got)
	}

	builder = JSDocInfo_Builder_()
	builder.ParseDocumentation()
	builder.RecordModifies([]string{"arguments"})
	info = builder.Build()
	if want := jutil.NewOrderedSet("arguments"); !info.GetModifies().Equals(want) {
		t.Errorf("getModifies() = %v, want %v", info.GetModifies(), want)
	}
	if got := info.GetModifiesSlice(); !slices.Equal(got, []string{"arguments"}) {
		t.Errorf("GetModifiesSlice() = %v, want [arguments]", got)
	}

	builder = JSDocInfo_Builder_()
	builder.ParseDocumentation()
	if !builder.RecordModifies([]string{"this", "arguments", "this"}) {
		t.Fatal("RecordModifies failed")
	}
	info = builder.Build()
	if got := info.GetModifiesSlice(); !slices.Equal(got, []string{"this", "arguments"}) {
		t.Errorf("GetModifiesSlice() = %v, want [this arguments]", got)
	}
}

func TestJSDocInfo_AddSingleTemplateTypeName(t *testing.T) {
	builder := JSDocInfo_Builder_()
	builder.ParseDocumentation()
	assertTrue(t, "recordTemplateTypeName(T)", builder.RecordTemplateTypeName("T"))
	info := builder.Build()

	typeNames := []string{"T"}
	if got := info.GetTemplateTypeNames(); !slices.Equal(got, typeNames) {
		t.Errorf("getTemplateTypeNames() = %q, want %q", got, typeNames)
	}
}

func TestJSDocInfo_AddMultipleTemplateTypeName(t *testing.T) {
	builder := JSDocInfo_Builder_()
	builder.ParseDocumentation()
	typeNames := []string{"T", "R"}
	builder.RecordTemplateTypeName("T")
	builder.RecordTemplateTypeName("R")
	info := builder.Build()
	if got := info.GetTemplateTypeNames(); !slices.Equal(got, typeNames) {
		t.Errorf("getTemplateTypeNames() = %q, want %q", got, typeNames)
	}
}

func TestJSDocInfo_FailToAddTemplateTypeName(t *testing.T) {
	builder := JSDocInfo_Builder_()
	builder.ParseDocumentation()
	builder.RecordTemplateTypeName("T")
	assertFalse(t, "second recordTemplateTypeName(T)", builder.RecordTemplateTypeName("T"))
}

func TestJSDocInfo_GetThrowsDescription(t *testing.T) {
	builder := JSDocInfo_Builder_()
	builder.ParseDocumentation()

	// Set a description so that builder is initialized.
	builder.RecordDescription("Lorem")
	builder.RecordThrowsAnnotation("{Error} Because it does.")
	builder.RecordThrowsAnnotation("{not a type}")
	info := builder.Build()
	want := []string{"{Error} Because it does.", "{not a type}"}
	if got := info.GetThrowsAnnotations(); !slices.Equal(got, want) {
		t.Errorf("getThrowsAnnotations() = %q, want %q", got, want)
	}
}

// https://github.com/google/closure-compiler/issues/2328
func TestJSDocInfo_GetTypeNodes_ExcludesNull(t *testing.T) {
	builder := JSDocInfo_Builder_()

	// no way to add to implemented interfaces
	assertFalse(t, "recordImplementedInterface(null)", builder.RecordImplementedInterface(nil))
	info := builder.BuildWithAlways(true)

	nodes := info.GetTypeNodes()
	if len(nodes) != 0 {
		t.Errorf("getTypeNodes() = %v, want empty", nodes)
	}
}

func TestJSDocInfo_GetTypeNodes_IncludesTemplateTypeBounds(t *testing.T) {
	// Given
	builder := JSDocInfo_Builder_()

	// When
	builder.RecordTemplateTypeNameWithBound("A", nil) // Uses `?` as the bounding expression by default.
	builder.RecordTemplateTypeNameWithBound("B", jsdocFromString("Foo"))
	builder.RecordTemplateTypeNameWithBound("C", jsdocFromString("Bar"))
	info := builder.Build()

	// Then
	var upperBoundRoots []Node
	for e := range info.GetTemplateTypes().Values() {
		upperBoundRoots = append(upperBoundRoots, e.GetRoot())
	}
	nodes := info.GetTypeNodes()

	if len(nodes) != 3 {
		t.Fatalf("getTypeNodes() has %d nodes, want 3", len(nodes))
	}
	ok := len(nodes) == len(upperBoundRoots)
	for _, r := range upperBoundRoots {
		ok = ok && slices.Contains(nodes, r)
	}
	if !ok {
		t.Errorf("getTypeNodes() = %v, want exactly %v", nodes, upperBoundRoots)
	}
}

func TestJSDocInfo_ContainsDeclaration_Implements(t *testing.T) {
	builder := JSDocInfo_Builder_()
	builder.RecordVisibility(JSDocInfo_Visibility_INHERITED)
	builder.RecordImplementedInterface(jsdocFromString("MyInterface"))
	info := builder.Build()

	if got := info.GetImplementedInterfaceCount(); got != 1 {
		t.Errorf("getImplementedInterfaceCount() = %d, want 1", got)
	}
	assertTrue(t, "containsDeclaration", info.ContainsDeclaration())
}

func TestJSDocInfo_ContainsDeclaration_Extends(t *testing.T) {
	builder := JSDocInfo_Builder_()
	builder.RecordVisibility(JSDocInfo_Visibility_INHERITED)
	builder.RecordBaseType(jsdocFromString("MyBaseClass"))
	info := builder.Build()

	assertTrue(t, "hasBaseType", info.HasBaseType())
	assertTrue(t, "containsDeclaration", info.ContainsDeclaration())
}

func TestJSDocInfo_ClosurePrimitiveId_AffectsEquality(t *testing.T) {
	builder := JSDocInfo_Builder_()
	builder.RecordClosurePrimitiveId("asserts.fail")
	assertsFailJSDoc := builder.Build()

	assertFalse(t, "areEquivalent(assertsFail, EMPTY)", JSDocInfo_AreEquivalent(assertsFailJSDoc, jsdocTestEMPTY))

	builder = JSDocInfo_Builder_()
	builder.RecordClosurePrimitiveId("asserts.fail")
	assertTrue(t, "areEquivalent(assertsFail, assertsFail')", JSDocInfo_AreEquivalent(assertsFailJSDoc, builder.Build()))
}

func assertPerFileClosureUnawareMode(t *testing.T, info *JSDocInfo, want JSDocInfo_PerFileClosureUnawareMode) {
	t.Helper()
	if got := info.GetPerFileClosureUnawareMode(); got == nil || *got != want {
		t.Errorf("getPerFileClosureUnawareMode() = %v, want %v", got, want)
	}
}

func TestJSDocInfo_RecordClosureUnaware_PerFileMode(t *testing.T) {
	builder := JSDocInfo_Builder_()
	builder.RecordClosureUnawareCodeWithMode(JSDocInfo_PerFileClosureUnawareMode_WHITESPACE)
	info := builder.Build()

	assertTrue(t, "isClosureUnawareCode", info.IsClosureUnawareCode())
	assertPerFileClosureUnawareMode(t, info, JSDocInfo_PerFileClosureUnawareMode_WHITESPACE)
}

func TestJSDocInfo_RecordClosureUnaware_PerFileMode_Duplicate(t *testing.T) {
	builder := JSDocInfo_Builder_()
	builder.RecordClosureUnawareCodeWithMode(JSDocInfo_PerFileClosureUnawareMode_WHITESPACE)
	builder.RecordClosureUnawareCodeWithMode(JSDocInfo_PerFileClosureUnawareMode_SIMPLE)
	info := builder.Build()

	assertTrue(t, "isClosureUnawareCode", info.IsClosureUnawareCode())
	assertPerFileClosureUnawareMode(t, info, JSDocInfo_PerFileClosureUnawareMode_WHITESPACE)
}

func TestJSDocInfo_NoClosureUnawarePerFileMode_DefaultsToUnspecified(t *testing.T) {
	builder := JSDocInfo_Builder_()
	builder.RecordClosureUnawareCode()
	info := builder.Build()
	assertTrue(t, "isClosureUnawareCode", info.IsClosureUnawareCode())
	assertPerFileClosureUnawareMode(t, info, JSDocInfo_PerFileClosureUnawareMode_UNSPECIFIED)
}

// jsdocFromString is JSDocInfoTest.fromString: gets the type expression for
// a simple type name.
func jsdocFromString(s string) *JSTypeExpression {
	return NewJSTypeExpression(Node_NewString(s), "")
}

// Generates a sample type expression tree
//
//	          BANG
//	           |
//	          Item
//	        /  |   \
//	   QMARK QMARK QMARK
//	      /    |    \
//	string boolean  AnotherItem
func createSampleTypeExpression() *JSTypeExpression {
	root := NewNodeWithChild(Token_BANG, Node_NewString("Item"))
	child1 := NewNodeWithChild(Token_QMARK, Node_NewString("string"))
	child2 := NewNodeWithChild(Token_QMARK, Node_NewString("boolean"))
	child3 := NewNodeWithChild(Token_QMARK, Node_NewString("AnotherItem"))
	root.AddChildToBack(child1)
	root.AddChildToBack(child2)
	root.AddChildToBack(child3)
	return NewJSTypeExpression(root, "")
}

// Generates a sample type expression tree with root node as a typename
//
//	          Item
//	        /  |   \
//	   QMARK QMARK QMARK
//	      /    |    \
//	string boolean  AnotherItem
func createSampleTypeExpression_rootReplacement() *JSTypeExpression {
	root := Node_NewString("Item")
	child1 := NewNodeWithChild(Token_QMARK, Node_NewString("string"))
	child2 := NewNodeWithChild(Token_QMARK, Node_NewString("boolean"))
	child3 := NewNodeWithChild(Token_QMARK, Node_NewString("AnotherItem"))
	root.AddChildToBack(child1)
	root.AddChildToBack(child2)
	root.AddChildToBack(child3)
	return NewJSTypeExpression(root, "")
}

var jsdocTestEMPTY = JSDocInfo_Builder_().BuildWithAlways(true)

func TestJSDocInfoPropertyCloneCopiesBackingSlices(t *testing.T) {
	strList := []string{"a", "b"}
	clonedStrList := JSDocInfo_AUTHORS.Clone(strList, nil).([]string)
	clonedStrList[0] = "mutated"
	if strList[0] != "a" {
		t.Fatalf("mutating cloned string list corrupted original: %v", strList)
	}

	markers := []JSDocInfo_Marker{{}}
	clonedMarkers := JSDocInfo_MARKERS.Clone(markers, nil).([]JSDocInfo_Marker)
	if &clonedMarkers[0] == &markers[0] {
		t.Fatal("expected JSDocInfo_MARKERS.Clone to allocate a distinct backing array")
	}

	kvs := []jsDocKV{{key: "k", val: "v1"}}
	clonedKVs := JSDocInfo_PARAMETER_DESCRIPTIONS.Clone(kvs, nil).([]jsDocKV)
	clonedKVs[0].val = "v2"
	if kvs[0].val != "v1" {
		t.Fatalf("mutating cloned jsDocKV slice corrupted original: %+v", kvs)
	}

	supps := []JSDocSuppression{{Warnings: []string{"checkTypes"}, Description: "orig"}}
	clonedSupps := JSDocInfo_SUPPRESSIONS.Clone(supps, nil).([]JSDocSuppression)
	clonedSupps[0].Description = "changed"
	if supps[0].Description != "orig" {
		t.Fatalf("mutating cloned JSDocSuppression slice corrupted original: %+v", supps)
	}
}

func TestJSDocInfoTypeTransformEqualValuesNilTypeNode(t *testing.T) {
	bNil1 := JSDocInfo_Builder_()
	bNil1.RecordTypeTransformation("T", NilNode)
	infoNil1 := bNil1.Build()

	bNil2 := JSDocInfo_Builder_()
	bNil2.RecordTypeTransformation("T", NilNode)
	infoNil2 := bNil2.Build()

	bNonNil := JSDocInfo_Builder_()
	bNonNil.RecordTypeTransformation("T", Node_NewString("string"))
	infoNonNil := bNonNil.Build()

	if !JSDocInfo_AreEquivalent(infoNil1, infoNil1) {
		t.Fatal("expected JSDocInfo with NilNode TypeTransformation to equal itself")
	}
	if !JSDocInfo_AreEquivalent(infoNil1, infoNil2) || !JSDocInfo_AreEquivalent(infoNil2, infoNil1) {
		t.Fatal("expected two JSDocInfos with NilNode TypeTransformation for T to be equal")
	}
	if JSDocInfo_AreEquivalent(infoNil1, infoNonNil) || JSDocInfo_AreEquivalent(infoNonNil, infoNil1) {
		t.Fatal("expected NilNode and non-nil TypeTransformation JSDocInfos to be unequal without panicking")
	}
}

func TestNodeAndJSTypeZeroOverheadStringAccessors(t *testing.T) {
	n := IR_Name("x")
	if name, ok := n.SourceFileName(); ok || name != "" {
		t.Fatalf("expected unset SourceFileName, got (%q, %v)", name, ok)
	}
	if orig, ok := n.OriginalName(); ok || orig != "" {
		t.Fatalf("expected unset OriginalName, got (%q, %v)", orig, ok)
	}
	n.SetSourceFileForTesting("input.js")
	n.SetOriginalNameString("origX")
	if name, ok := n.SourceFileName(); !ok || name != "input.js" {
		t.Fatalf("SourceFileName() = (%q, %v), want (\"input.js\", true)", name, ok)
	}
	if orig, ok := n.OriginalName(); !ok || orig != "origX" {
		t.Fatalf("OriginalName() = (%q, %v), want (\"origX\", true)", orig, ok)
	}

	cooked := "a\nb"
	sub := Node_NewTemplateLitString(&cooked, `a\nb`)
	if s, ok := sub.CookedString(); !ok || s != "a\nb" {
		t.Fatalf("CookedString() = (%q, %v), want (\"a\\nb\", true)", s, ok)
	}
	invalidSub := Node_NewTemplateLitString(nil, `\u{`)
	if s, ok := invalidSub.CookedString(); ok || s != "" {
		t.Fatalf("invalid CookedString() = (%q, %v), want (\"\", false)", s, ok)
	}

	reg := NewJSTypeRegistry(nil)
	if d, ok := reg.GetNativeType(JSTypeNative_NUMBER_TYPE).DisplayName(); !ok || d != "number" {
		t.Fatalf("NUMBER_TYPE.DisplayName() = (%q, %v), want (\"number\", true)", d, ok)
	}
	if r, ok := reg.GetNativeType(JSTypeNative_NUMBER_TYPE).ReferenceName(); ok || r != "" {
		t.Fatalf("NUMBER_TYPE.ReferenceName() = (%q, %v), want (\"\", false)", r, ok)
	}
	if d, ok := reg.GetNativeType(JSTypeNative_NO_TYPE).DisplayName(); ok || d != "" {
		t.Fatalf("NO_TYPE.DisplayName() = (%q, %v), want (\"\", false)", d, ok)
	}
	if r, ok := reg.GetNativeObjectType(JSTypeNative_NO_OBJECT_TYPE).ReferenceName(); ok || r != "" {
		t.Fatalf("NO_OBJECT_TYPE.ReferenceName() = (%q, %v), want (\"\", false)", r, ok)
	}
	if r, ok := reg.GetNativeObjectType(JSTypeNative_ARRAY_TYPE).ReferenceName(); !ok || r != "Array" {
		t.Fatalf("ARRAY_TYPE.ReferenceName() = (%q, %v), want (\"Array\", true)", r, ok)
	}
	tmpl := reg.CreateTemplateType("T")
	if r, ok := tmpl.ReferenceName(); !ok || r != "T" {
		t.Fatalf("TemplateType.ReferenceName() = (%q, %v), want (\"T\", true)", r, ok)
	}
}
