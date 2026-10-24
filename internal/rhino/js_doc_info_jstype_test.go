// Port of com/google/javascript/rhino/JSDocInfoTest.java @ closure-compiler v20261021

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

package rhino_test

// The JSDocInfoTest cases that resolve type expressions through the
// JSTypeRegistry. Go-only: language-adapter — they use rhinotesting, which
// imports rhino, so they live in the external test package; the other
// JSDocInfoTest cases are in js_doc_info_test.go.

import (
	"testing"

	. "github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/rhino/rhinotesting"
)

type jsDocInfoTest struct {
	rhinotesting.Truth
	errorReporter *rhinotesting.TestErrorReporter
	registry      *JSTypeRegistry
}

func newJsDocInfoTest(t *testing.T) *jsDocInfoTest {
	s := &jsDocInfoTest{Truth: rhinotesting.NewTruth(t)}
	s.errorReporter = rhinotesting.NewTestErrorReporter()
	s.registry = NewJSTypeRegistry(s.errorReporter)
	t.Cleanup(s.validateWarningsAndErrors)
	return s
}

func (s *jsDocInfoTest) getNativeType(typeId JSTypeNative) JSType {
	return s.registry.GetNativeType(typeId)
}

func (s *jsDocInfoTest) validateWarningsAndErrors() {
	s.errorReporter.VerifyHasEncounteredAllWarningsAndErrors(s.T)
}

func TestJSDocInfo_SetType(t *testing.T) {
	s := newJsDocInfoTest(t)
	builder := JSDocInfo_Builder_()
	builder.RecordType(jsDocInfoTest_fromString("string"))
	info := builder.Build()

	s.AssertThat(info.GetBaseType()).IsNull()
	s.AssertThat(info.GetDescription()).IsNull()
	s.AssertThat(info.GetEnumParameterType()).IsNull()
	s.AssertThat(info.GetParameterCount()).IsEqualTo(0)
	s.AssertThat(info.GetReturnType()).IsNull()
	s.AssertType(s.resolve(info.GetType())).IsEqualTo(s.getNativeType(JSTypeNative_STRING_TYPE))
	s.AssertThat(info.GetVisibility()).IsEqualTo(JSDocInfo_Visibility_INHERITED)
	s.AssertThat(info.HasType()).IsTrue()
	s.AssertThat(info.IsConstant()).IsFalse()
	s.AssertThat(info.IsConstructor()).IsFalse()
}

func TestJSDocInfo_SetTypeAndVisibility(t *testing.T) {
	s := newJsDocInfoTest(t)
	builder := JSDocInfo_Builder_()
	builder.RecordType(jsDocInfoTest_fromString("string"))
	builder.RecordVisibility(JSDocInfo_Visibility_PROTECTED)
	info := builder.Build()

	s.AssertThat(info.GetBaseType()).IsNull()
	s.AssertThat(info.GetDescription()).IsNull()
	s.AssertThat(info.GetEnumParameterType()).IsNull()
	s.AssertThat(info.GetParameterCount()).IsEqualTo(0)
	s.AssertThat(info.GetReturnType()).IsNull()
	s.AssertType(s.resolve(info.GetType())).IsEqualTo(s.getNativeType(JSTypeNative_STRING_TYPE))
	s.AssertThat(info.GetVisibility()).IsEqualTo(JSDocInfo_Visibility_PROTECTED)
	s.AssertThat(info.HasType()).IsTrue()
	s.AssertThat(info.IsConstant()).IsFalse()
	s.AssertThat(info.IsConstructor()).IsFalse()
}

func TestJSDocInfo_SetReturnType(t *testing.T) {
	s := newJsDocInfoTest(t)
	builder := JSDocInfo_Builder_()
	builder.RecordReturnType(jsDocInfoTest_fromString("string"))
	info := builder.Build()

	s.AssertThat(info.GetBaseType()).IsNull()
	s.AssertThat(info.GetDescription()).IsNull()
	s.AssertThat(info.GetEnumParameterType()).IsNull()
	s.AssertThat(info.GetParameterCount()).IsEqualTo(0)
	s.AssertType(s.resolve(info.GetReturnType())).IsEqualTo(s.getNativeType(JSTypeNative_STRING_TYPE))
	s.AssertThat(info.GetType()).IsNull()
	s.AssertThat(info.GetVisibility()).IsEqualTo(JSDocInfo_Visibility_INHERITED)
	s.AssertThat(info.HasType()).IsFalse()
	s.AssertThat(info.IsConstant()).IsFalse()
	s.AssertThat(info.IsConstructor()).IsFalse()
}

func TestJSDocInfo_SetReturnTypeAndBaseType(t *testing.T) {
	s := newJsDocInfoTest(t)
	builder := JSDocInfo_Builder_()
	builder.RecordBaseType(
		NewJSTypeExpression(NewNodeWithChild(Token_BANG, Node_NewString("Number")), ""))
	builder.RecordReturnType(jsDocInfoTest_fromString("string"))
	info := builder.Build()

	s.AssertType(s.resolve(info.GetBaseType())).IsEqualTo(s.getNativeType(JSTypeNative_NUMBER_OBJECT_TYPE))
	s.AssertThat(info.GetDescription()).IsNull()
	s.AssertThat(info.GetEnumParameterType()).IsNull()
	s.AssertThat(info.GetParameterCount()).IsEqualTo(0)
	s.AssertType(s.resolve(info.GetReturnType())).IsEqualTo(s.getNativeType(JSTypeNative_STRING_TYPE))
	s.AssertThat(info.GetType()).IsNull()
	s.AssertThat(info.GetVisibility()).IsEqualTo(JSDocInfo_Visibility_INHERITED)
	s.AssertThat(info.HasType()).IsFalse()
	s.AssertThat(info.IsConstant()).IsFalse()
	s.AssertThat(info.IsConstructor()).IsFalse()
}

func TestJSDocInfo_SetEnumParameterType(t *testing.T) {
	s := newJsDocInfoTest(t)
	builder := JSDocInfo_Builder_()
	builder.RecordEnumParameterType(jsDocInfoTest_fromString("string"))
	info := builder.Build()

	s.AssertThat(info.GetBaseType()).IsNull()
	s.AssertThat(info.GetDescription()).IsNull()
	s.AssertType(s.resolve(info.GetEnumParameterType())).IsEqualTo(s.getNativeType(JSTypeNative_STRING_TYPE))
	s.AssertThat(info.GetParameterCount()).IsEqualTo(0)
	s.AssertThat(info.GetReturnType()).IsNull()
	s.AssertThat(info.GetType()).IsNull()
	s.AssertThat(info.GetVisibility()).IsEqualTo(JSDocInfo_Visibility_INHERITED)
	s.AssertThat(info.HasType()).IsFalse()
	s.AssertThat(info.IsConstant()).IsFalse()
	s.AssertThat(info.IsConstructor()).IsFalse()
}

func TestJSDocInfo_MultipleSetType(t *testing.T) {
	s := newJsDocInfoTest(t)
	builder := JSDocInfo_Builder_()
	builder.RecordType(jsDocInfoTest_fromString("number"))

	s.AssertThat(builder.RecordReturnType(jsDocInfoTest_fromString("boolean"))).IsFalse()
	s.AssertThat(builder.RecordEnumParameterType(jsDocInfoTest_fromString("string"))).IsFalse()
	s.AssertThat(builder.RecordTypedef(jsDocInfoTest_fromString("string"))).IsFalse()
	info := builder.Build()

	s.AssertType(s.resolve(info.GetType())).IsEqualTo(s.getNativeType(JSTypeNative_NUMBER_TYPE))
	s.AssertThat(info.GetReturnType()).IsNull()
	s.AssertThat(info.GetEnumParameterType()).IsNull()
	s.AssertThat(info.GetTypedefType()).IsNull()
	s.AssertThat(info.HasType()).IsTrue()
}

func TestJSDocInfo_MultipleSetType2(t *testing.T) {
	s := newJsDocInfoTest(t)
	builder := JSDocInfo_Builder_()

	builder.RecordReturnType(jsDocInfoTest_fromString("boolean"))

	s.AssertThat(builder.RecordType(jsDocInfoTest_fromString("number"))).IsFalse()
	s.AssertThat(builder.RecordEnumParameterType(jsDocInfoTest_fromString("string"))).IsFalse()
	s.AssertThat(builder.RecordTypedef(jsDocInfoTest_fromString("string"))).IsFalse()
	info := builder.Build()

	s.AssertType(s.resolve(info.GetReturnType())).IsEqualTo(s.getNativeType(JSTypeNative_BOOLEAN_TYPE))
	s.AssertThat(info.GetEnumParameterType()).IsNull()
	s.AssertThat(info.GetType()).IsNull()
	s.AssertThat(info.GetTypedefType()).IsNull()
	s.AssertThat(info.HasType()).IsFalse()
}

func TestJSDocInfo_MultipleSetType3(t *testing.T) {
	s := newJsDocInfoTest(t)
	builder := JSDocInfo_Builder_()
	builder.RecordEnumParameterType(jsDocInfoTest_fromString("boolean"))

	s.AssertThat(builder.RecordType(jsDocInfoTest_fromString("number"))).IsFalse()
	s.AssertThat(builder.RecordReturnType(jsDocInfoTest_fromString("string"))).IsFalse()
	s.AssertThat(builder.RecordTypedef(jsDocInfoTest_fromString("string"))).IsFalse()
	info := builder.Build()

	s.AssertThat(info.GetType()).IsNull()
	s.AssertThat(info.GetTypedefType()).IsNull()
	s.AssertThat(info.GetReturnType()).IsNull()
	s.AssertType(s.resolve(info.GetEnumParameterType())).IsEqualTo(s.getNativeType(JSTypeNative_BOOLEAN_TYPE))
}

func TestJSDocInfo_SetTypedefType(t *testing.T) {
	s := newJsDocInfoTest(t)
	builder := JSDocInfo_Builder_()
	builder.RecordTypedef(jsDocInfoTest_fromString("boolean"))
	info := builder.Build()

	s.AssertType(s.resolve(info.GetTypedefType())).IsEqualTo(s.getNativeType(JSTypeNative_BOOLEAN_TYPE))
	s.AssertThat(info.HasTypedefType()).IsTrue()
	s.AssertThat(info.HasType()).IsFalse()
	s.AssertThat(info.HasEnumParameterType()).IsFalse()
	s.AssertThat(info.HasReturnType()).IsFalse()
}

func TestJSDocInfo_Clone(t *testing.T) {
	s := newJsDocInfoTest(t)
	builder := JSDocInfo_Builder_()
	builder.RecordDescription("The source info")
	builder.RecordConstancy()
	builder.RecordConstructor()
	builder.RecordNoCollapse()
	builder.RecordBaseType(
		NewJSTypeExpression(NewNodeWithChild(Token_BANG, Node_NewString("Number")), ""))
	builder.RecordReturnType(jsDocInfoTest_fromString("string"))
	info := builder.Build()

	cloned := info.Clone()

	s.AssertType(s.resolve(cloned.GetBaseType())).IsEqualTo(s.getNativeType(JSTypeNative_NUMBER_OBJECT_TYPE))
	s.AssertThat(cloned.GetDescription()).IsEqualTo("The source info")
	s.AssertType(s.resolve(cloned.GetReturnType())).IsEqualTo(s.getNativeType(JSTypeNative_STRING_TYPE))
	s.AssertThat(cloned.IsConstant()).IsTrue()
	s.AssertThat(cloned.IsConstructor()).IsTrue()
	s.AssertThat(cloned.IsNoCollapse()).IsTrue()

	// TODO - cannot change these things!
	// cloned.recordDescription("The cloned info");
	// cloned.setHidden(false);
	// cloned.recordBaseType(fromString("string"));

	// assertType(resolve(cloned.getBaseType())).isEqualTo(getNativeType(STRING_TYPE));
	// assertThat(cloned.getDescription()).isEqualTo("The cloned info");
	// assertThat(cloned.isNoCollapse()).isFalse();

	// // Original info should be unchanged.
	// assertType(resolve(info.getBaseType()))
	//     .isEqualTo(getNativeType(NUMBER_OBJECT_TYPE));
	// assertThat(info.getDescription()).isEqualTo("The source info");
	// assertType(resolve(info.getReturnType())).isEqualTo(getNativeType(STRING_TYPE));
	// assertThat(info.isConstant()).isTrue();
	// assertThat(info.isConstructor()).isTrue();
	// assertThat(info.isNoCollapse()).isTrue();
}

func TestJSDocInfo_CloneTypeExpressions1(t *testing.T) {
	s := newJsDocInfoTest(t)
	builder := JSDocInfo_Builder_()
	builder.RecordDescription("The source info")
	builder.RecordConstancy()
	builder.RecordConstructor()
	builder.RecordBaseType(
		NewJSTypeExpression(NewNodeWithChild(Token_BANG, Node_NewString("Number")), ""))
	builder.RecordReturnType(jsDocInfoTest_fromString("string"))
	builder.RecordParameter("a", jsDocInfoTest_fromString("string"))
	info := builder.Build()

	cloned := info.CloneWithTypeNodes(true)

	s.AssertThat(cloned.GetBaseType().GetRoot()).IsNotSameInstanceAs(info.GetBaseType().GetRoot())
	s.AssertType(s.resolve(cloned.GetBaseType())).IsEqualTo(s.getNativeType(JSTypeNative_NUMBER_OBJECT_TYPE))
	s.AssertThat(cloned.GetDescription()).IsEqualTo("The source info")
	s.AssertThat(cloned.GetReturnType().GetRoot()).
		IsNotSameInstanceAs(info.GetReturnType().GetRoot())
	s.AssertType(s.resolve(cloned.GetReturnType())).IsEqualTo(s.getNativeType(JSTypeNative_STRING_TYPE))
	s.AssertThat(cloned.GetParameterType("a").GetRoot()).
		IsNotSameInstanceAs(info.GetParameterType("a").GetRoot())
	s.AssertType(s.resolve(cloned.GetParameterType("a"))).IsEqualTo(s.getNativeType(JSTypeNative_STRING_TYPE))
}

// Gets the type expression for a simple type name.
//
// Java: JSDocInfoTest.fromString(String)
func jsDocInfoTest_fromString(s string) *JSTypeExpression {
	return NewJSTypeExpression(Node_NewString(s), "")
}

// Java: JSDocInfoTest.resolve(JSTypeExpression, String...)
func (s *jsDocInfoTest) resolve(n *JSTypeExpression, warnings ...string) JSType {
	s.errorReporter.ExpectAllWarnings(warnings...)
	return n.Evaluate(nil, s.registry)
}
