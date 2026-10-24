// Port of com/google/javascript/jscomp/XtbMessageBundleTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/XtbMessageBundleTest.java)

/*
 * Copyright 2006 The Closure Compiler Authors.
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
	"bytes"
	"slices"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// XtbMessageBundleTest tests XtbMessageBundle.
type XtbMessageBundleTest struct {
	// Go-only: test-helper — the Go test running the current test method.
	t *testing.T
}

const xtbMessageBundleTest_PROJECT_ID = "TestProject"

// Java: XtbMessageBundleTest()
func NewXtbMessageBundleTest() *XtbMessageBundleTest {
	return &XtbMessageBundleTest{}
}

func TestXtbMessageBundleTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewXtbMessageBundleTest() })
}

// JUnitSetT gives the instance the Go test running its test method.
func (c *XtbMessageBundleTest) JUnitSetT(t *testing.T) { c.t = t }

// JUnitBefore: the class has no @Before method.
func (c *XtbMessageBundleTest) JUnitBefore() {}

// JUnitAfter: the class has no @After method.
func (c *XtbMessageBundleTest) JUnitAfter() {}

// Java: XtbMessageBundleTest.testXtbBundle()
func (c *XtbMessageBundleTest) TestXtbBundle() {
	xtb :=
		"<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" +
			"<!DOCTYPE translationbundle SYSTEM \"translationbundle.dtd\">\n" +
			"<translationbundle lang=\"zh-HK\">\n" +
			"<translation id=\"7639678437384034548\">descargar</translation>\n" +
			"<translation id=\"2398375912250604550\">Se han\n" +
			"ignorado <ph name=\"NUM\"/> conversaciones.</translation>\n" +
			"<translation id=\"6323937743550839320\"><ph name=\"P_START\"/>Si, puede <ph name=\"LINK_START_1_3\"/>hacer clic<ph name=\"LINK_END_1_3\"/> para utilizar.<ph name=\"P_END\"/><ph name=\"P_START\"/>Esperamos poder ampliar.<ph name=\"P_END\"/></translation>\n" +
			"<translation id=\"3945720239421293834\"></translation>\n" +
			"</translationbundle>"
	stream := bytes.NewReader([]byte(xtb))
	bundle := NewXtbMessageBundle(stream, jutil.StringPtr(xtbMessageBundleTest_PROJECT_ID))

	message := bundle.GetMessage("7639678437384034548")
	AssertThat(message.AsJsMessageString()).IsEqualTo("descargar")

	message = bundle.GetMessage("2398375912250604550")
	AssertThat(message.AsJsMessageString()).IsEqualTo("Se han\nignorado {$num} conversaciones.")

	message = bundle.GetMessage("6323937743550839320")
	AssertThat(message.AsJsMessageString()).
		IsEqualTo(
			"{$pStart}Si, puede {$linkStart_1_3}hacer clic{$linkEnd_1_3} para utilizar.{$pEnd}{$pStart}Esperamos poder ampliar.{$pEnd}")

	message = bundle.GetMessage("3945720239421293834")
	AssertThat(message.AsJsMessageString()).IsEmpty()
	AssertThat(message.GetParts()).IsNotEmpty()
}

// Java: XtbMessageBundleTest.testXtbBundle_genderedMessageVariantsWithoutPlaceholders()
func (c *XtbMessageBundleTest) TestXtbBundle_genderedMessageVariantsWithoutPlaceholders() {
	xtb :=
		"<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" +
			"<!DOCTYPE translationbundle SYSTEM \"translationbundle.dtd\">\n" +
			"<translationbundle lang=\"es_ES\">\n" +
			"<translation id=\"7639678437384034548\">\n" +
			"<branch variants=\"variants { grammatical_gender_variant { grammatical_gender_case:\n" +
			"MASCULINE } }\">Bienvenido!</branch><branch variants=\"variants { grammatical_gender_variant { grammatical_gender_case:\n" +
			"FEMININE } }\">Bienvenida!</branch><branch variants=\"variants { grammatical_gender_variant { grammatical_gender_case:\n" +
			"NEUTER } }\">Te damos la bienvenida!</branch><branch variants=\"variants { grammatical_gender_variant { grammatical_gender_case:\n" +
			"OTHER } }\">Te damos la bienvenida! - OTHER</branch></translation>\n" +
			"</translationbundle>\n"
	stream := bytes.NewReader([]byte(xtb))
	bundle := NewXtbMessageBundle(stream, jutil.StringPtr(xtbMessageBundleTest_PROJECT_ID))

	message := bundle.GetMessage("7639678437384034548")
	// Testing invalid calls to asJsMessageString() as it should call
	// asJsMessageString(GrammaticalGenderCase)
	AssertThrows[jutil.IUnsupportedOperationException](func() { message.AsJsMessageString() })

	AssertThat(message.AsJsMessageStringForGender(JsMessage_GrammaticalGenderCase_MASCULINE)).
		IsEqualTo("Bienvenido!")
	AssertThat(message.AsJsMessageStringForGender(JsMessage_GrammaticalGenderCase_FEMININE)).
		IsEqualTo("Bienvenida!")
	AssertThat(message.AsJsMessageStringForGender(JsMessage_GrammaticalGenderCase_NEUTER)).
		IsEqualTo("Te damos la bienvenida!")
	AssertThat(message.AsJsMessageStringForGender(JsMessage_GrammaticalGenderCase_OTHER)).
		IsEqualTo("Te damos la bienvenida! - OTHER")

	AssertThat(message.GetGenderedMessageParts(JsMessage_GrammaticalGenderCase_MASCULINE)).
		HasSize(1)
	AssertThat(
		message.
			GetGenderedMessageParts(JsMessage_GrammaticalGenderCase_MASCULINE)[0].
			GetString()).
		IsEqualTo("Bienvenido!")
	AssertThat(
		message.
			GetGenderedMessageParts(JsMessage_GrammaticalGenderCase_FEMININE)[0].
			GetString()).
		IsEqualTo("Bienvenida!")
	AssertThat(
		message.
			GetGenderedMessageParts(JsMessage_GrammaticalGenderCase_NEUTER)[0].
			GetString()).
		IsEqualTo("Te damos la bienvenida!")
	AssertThat(
		message.
			GetGenderedMessageParts(JsMessage_GrammaticalGenderCase_OTHER)[0].
			GetString()).
		IsEqualTo("Te damos la bienvenida! - OTHER")
}

// Java: XtbMessageBundleTest.testXtbBundle_genderedMessageVariantsWithPlaceholders()
func (c *XtbMessageBundleTest) TestXtbBundle_genderedMessageVariantsWithPlaceholders() {
	xtb :=
		"<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" +
			"<!DOCTYPE translationbundle SYSTEM \"translationbundle.dtd\">\n" +
			"<translationbundle lang=\"es_ES\">\n" +
			"<translation id=\"7639678437384034548\">\n" +
			"<branch variants=\"variants { grammatical_gender_variant { grammatical_gender_case:\n" +
			"MASCULINE } }\">Bienvenido, <ph name=\"NAME\"/></branch><branch variants=\"variants { grammatical_gender_variant { grammatical_gender_case:\n" +
			"FEMININE } }\">Bienvenida, <ph name=\"NAME\"/></branch><branch variants=\"variants { grammatical_gender_variant { grammatical_gender_case:\n" +
			"NEUTER } }\">Te damos la bienvenida, <ph name=\"NAME\"/></branch><branch variants=\"variants { grammatical_gender_variant { grammatical_gender_case:\n" +
			"OTHER } }\">Te damos la bienvenida - OTHER, <ph name=\"NAME\"/></branch></translation>\n" +
			"\n" +
			"</translationbundle>\n"
	stream := bytes.NewReader([]byte(xtb))
	bundle := NewXtbMessageBundle(stream, jutil.StringPtr(xtbMessageBundleTest_PROJECT_ID))

	message := bundle.GetMessage("7639678437384034548")
	// Testing invalid calls to asJsMessageString() as it should call
	// asJsMessageString(GrammaticalGenderCase)
	AssertThrows[jutil.IUnsupportedOperationException](func() { message.AsJsMessageString() })
	// Testing calls to asJsMessageString(GrammaticalGenderCase)
	AssertThrows[jutil.IIllegalArgumentException](
		func() { message.AsJsMessageStringForGender(JsMessage_GrammaticalGenderCase_ValueOf("FOO")) })

	// Testing valid calls to asJsMessageString()
	AssertThat(message.AsJsMessageStringForGender(JsMessage_GrammaticalGenderCase_MASCULINE)).
		IsEqualTo("Bienvenido, {$name}")
	AssertThat(message.AsJsMessageStringForGender(JsMessage_GrammaticalGenderCase_FEMININE)).
		IsEqualTo("Bienvenida, {$name}")
	AssertThat(message.AsJsMessageStringForGender(JsMessage_GrammaticalGenderCase_NEUTER)).
		IsEqualTo("Te damos la bienvenida, {$name}")
	AssertThat(message.AsJsMessageStringForGender(JsMessage_GrammaticalGenderCase_OTHER)).
		IsEqualTo("Te damos la bienvenida - OTHER, {$name}")

	masculineParts :=
		message.GetGenderedMessageParts(JsMessage_GrammaticalGenderCase_MASCULINE)
	AssertThat(masculineParts).HasSize(2)
	AssertThat(masculineParts[0].GetString()).IsEqualTo("Bienvenido, ")
	AssertThat(masculineParts[1].GetJsPlaceholderName()).IsEqualTo("name")

	feminineParts :=
		message.GetGenderedMessageParts(JsMessage_GrammaticalGenderCase_FEMININE)
	AssertThat(feminineParts).HasSize(2)
	AssertThat(feminineParts[0].GetString()).IsEqualTo("Bienvenida, ")
	AssertThat(feminineParts[1].GetJsPlaceholderName()).IsEqualTo("name")

	neuterParts :=
		message.GetGenderedMessageParts(JsMessage_GrammaticalGenderCase_NEUTER)
	AssertThat(neuterParts).HasSize(2)
	AssertThat(neuterParts[0].GetString()).IsEqualTo("Te damos la bienvenida, ")
	AssertThat(neuterParts[1].GetJsPlaceholderName()).IsEqualTo("name")

	otherParts :=
		message.GetGenderedMessageParts(JsMessage_GrammaticalGenderCase_OTHER)
	AssertThat(otherParts).HasSize(2)
	AssertThat(otherParts[0].GetString()).IsEqualTo("Te damos la bienvenida - OTHER, ")
	AssertThat(otherParts[1].GetJsPlaceholderName()).IsEqualTo("name")
}

// Java: XtbMessageBundleTest.testXtbBundleWithIcu_genderedMessageVariantsWithPlaceholders()
func (c *XtbMessageBundleTest) TestXtbBundleWithIcu_genderedMessageVariantsWithPlaceholders() {
	xtb :=
		"<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" +
			"<!DOCTYPE translationbundle SYSTEM \"translationbundle.dtd\">\n" +
			"<translationbundle lang=\"es_ES\">\n" +
			"<translation id=\"7639678437384034548\">\n" +
			"<branch variants=\"variants { grammatical_gender_variant { grammatical_gender_case:\n" +
			"MASCULINE } }\">Bienvenido, <ph name=\"NAME\"/></branch><branch variants=\"variants { grammatical_gender_variant { grammatical_gender_case:\n" +
			"FEMININE } }\">Bienvenida, <ph name=\"NAME\"/></branch><branch variants=\"variants { grammatical_gender_variant { grammatical_gender_case:\n" +
			"NEUTER } }\">Te damos la bienvenida, <ph name=\"NAME\"/></branch><branch variants=\"variants { grammatical_gender_variant { grammatical_gender_case:\n" +
			"OTHER } }\">Te damos la bienvenida - OTHER, <ph name=\"NAME\"/></branch></translation>\n" +
			"\n" +
			"</translationbundle>\n"
	stream := bytes.NewReader([]byte(xtb))
	bundle := NewXtbMessageBundle(stream, jutil.StringPtr(xtbMessageBundleTest_PROJECT_ID))

	message := bundle.GetMessage("7639678437384034548")

	AssertThat(message.IsEmpty()).IsFalse()
	AssertThrows[jutil.IUnsupportedOperationException](func() { message.AsIcuMessageString() })

	AssertThat(message.AsIcuMessageStringForGender(JsMessage_GrammaticalGenderCase_MASCULINE)).
		IsEqualTo("Bienvenido, {NAME}")
	AssertThat(message.AsIcuMessageStringForGender(JsMessage_GrammaticalGenderCase_FEMININE)).
		IsEqualTo("Bienvenida, {NAME}")
	AssertThat(message.AsIcuMessageStringForGender(JsMessage_GrammaticalGenderCase_NEUTER)).
		IsEqualTo("Te damos la bienvenida, {NAME}")
	AssertThat(message.AsIcuMessageStringForGender(JsMessage_GrammaticalGenderCase_OTHER)).
		IsEqualTo("Te damos la bienvenida - OTHER, {NAME}")

	masculineParts :=
		message.GetGenderedMessageParts(JsMessage_GrammaticalGenderCase_MASCULINE)
	AssertThat(masculineParts).HasSize(2)
	AssertThat(masculineParts[0].GetString()).IsEqualTo("Bienvenido, ")
	AssertThat(masculineParts[1].GetCanonicalPlaceholderName()).IsEqualTo("NAME")

	feminineParts :=
		message.GetGenderedMessageParts(JsMessage_GrammaticalGenderCase_FEMININE)
	AssertThat(feminineParts).HasSize(2)
	AssertThat(feminineParts[0].GetString()).IsEqualTo("Bienvenida, ")
	AssertThat(feminineParts[1].GetCanonicalPlaceholderName()).IsEqualTo("NAME")

	neuterParts :=
		message.GetGenderedMessageParts(JsMessage_GrammaticalGenderCase_NEUTER)
	AssertThat(neuterParts).HasSize(2)
	AssertThat(neuterParts[0].GetString()).IsEqualTo("Te damos la bienvenida, ")
	AssertThat(neuterParts[1].GetCanonicalPlaceholderName()).IsEqualTo("NAME")

	otherParts :=
		message.GetGenderedMessageParts(JsMessage_GrammaticalGenderCase_OTHER)
	AssertThat(otherParts).HasSize(2)
	AssertThat(otherParts[0].GetString()).IsEqualTo("Te damos la bienvenida - OTHER, ")
	AssertThat(otherParts[1].GetCanonicalPlaceholderName()).IsEqualTo("NAME")
}

// Java: XtbMessageBundleTest.testXtbBundle_genderedMessageVariantsWithInvalidGenderCase_throwsException()
func (c *XtbMessageBundleTest) TestXtbBundle_genderedMessageVariantsWithInvalidGenderCase_throwsException() {
	xtb :=
		"<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" +
			"<!DOCTYPE translationbundle SYSTEM \"translationbundle.dtd\">\n" +
			"<translationbundle lang=\"es_ES\">\n" +
			"<translation id=\"7639678437384034548\">\n" +
			"<branch variants=\"variants { grammatical_gender_variant { grammatical_gender_case:\n" +
			"  MASCULINE } }\">Bienvenido!</branch><branch variants=\"variants { grammatical_gender_variant { grammatical_gender_case:\n" +
			"  FEMININE } }\">Bienvenida!</branch><branch variants=\"variants { grammatical_gender_variant { grammatical_gender_case:\n" +
			"  FOO } }\">Foo bar!</branch><branch variants=\"variants { grammatical_gender_variant { grammatical_gender_case:\n" +
			"  OTHER } }\">Te damos la bienvenida! - OTHER</branch></translation>\n" +
			"\n" +
			"</translationbundle>\n"
	stream := bytes.NewReader([]byte(xtb))

	e :=
		AssertThrows[jutil.IIllegalStateException](func() { NewXtbMessageBundle(stream, jutil.StringPtr(xtbMessageBundleTest_PROJECT_ID)) })
	AssertThat(e.GetMessage()).
		Contains(
			"Gender case must be one of the following: MASCULINE, FEMININE, NEUTER, or OTHER.")
}

/**
 * When using EXTERNAL messages with plurals/selects, the XTB files may contain a mix of ICU style
 * placeholders (i.e. {@code {foo}}) and regular placeholders (i.e. {@code <ph name="foo"/>}).
 * However, JsMessage and the Closure Library runtime don't expect to see regular placeholders, so
 * they must be rewritten.
 */
// Java: XtbMessageBundleTest.testXtbBundle_mixedPlaceholders()
func (c *XtbMessageBundleTest) TestXtbBundle_mixedPlaceholders() {
	xtbWithMixedPlaceholders :=
		"<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" +
			"<!DOCTYPE translationbundle SYSTEM \"translationbundle.dtd\">\n" +
			"<translationbundle lang=\"ru_RU\">\n" +
			"<translation id=\"123456\">{USER_GENDER,select,female{Hello <ph name=\"USER_IDENTIFIER\"/>.}male{Hello <ph name=\"USER_IDENTIFIER\"/>.}other{Hello <ph name=\"USER_IDENTIFIER\"/>.}}</translation>\n" +
			"<translation id=\"123457\"><ph name=\"START_PARAGRAPH\"/>p1<ph name=\"END_PARAGRAPH\"/><ph name=\"START_PARAGRAPH\"/>p1<ph name=\"END_PARAGRAPH\"/></translation></translationbundle>"
	stream := bytes.NewReader([]byte(xtbWithMixedPlaceholders))
	bundle := NewXtbMessageBundle(stream, jutil.StringPtr(xtbMessageBundleTest_PROJECT_ID))

	AssertThat(slices.Collect(bundle.GetAllMessages())).HasSize(2)
	icuMsg := bundle.GetMessage("123456")
	AssertThat(icuMsg.AsIcuMessageString()).
		IsEqualTo(
			"{USER_GENDER,select,female{Hello {USER_IDENTIFIER}.}male{Hello {USER_IDENTIFIER}.}other{Hello {USER_IDENTIFIER}.}}")
	// For an ICU selector formatted message, XtbMessageBundle automatically converts all the
	// placeholders into normal strings.
	icuMsgParts := icuMsg.GetParts()
	AssertThat(icuMsgParts).HasSize(7)
	AssertThat(icuMsgParts[0].GetString()).IsEqualTo("{USER_GENDER,select,female{Hello ")
	AssertThat(icuMsgParts[1].GetCanonicalPlaceholderName()).IsEqualTo("USER_IDENTIFIER")
	AssertThat(icuMsgParts[2].GetString()).IsEqualTo(".}male{Hello ")
	AssertThat(icuMsgParts[3].GetCanonicalPlaceholderName()).IsEqualTo("USER_IDENTIFIER")
	AssertThat(icuMsgParts[4].GetString()).IsEqualTo(".}other{Hello ")
	AssertThat(icuMsgParts[5].GetCanonicalPlaceholderName()).IsEqualTo("USER_IDENTIFIER")
	AssertThat(icuMsgParts[6].GetString()).IsEqualTo(".}}")

	// Previous ICU message should not to affect next message
	normalMsg := bundle.GetMessage("123457")
	AssertThat(normalMsg.AsJsMessageString()).
		IsEqualTo("{$startParagraph}p1{$endParagraph}{$startParagraph}p1{$endParagraph}")
	normalMsgParts := normalMsg.GetParts()
	// For a normal message the placeholders are not turned into strings
	AssertThat(normalMsgParts).HasSize(6)
	AssertThat(normalMsgParts[0].GetJsPlaceholderName()).IsEqualTo("startParagraph")
	AssertThat(normalMsgParts[1].GetString()).IsEqualTo("p1")
	AssertThat(normalMsgParts[2].GetJsPlaceholderName()).IsEqualTo("endParagraph")
	AssertThat(normalMsgParts[3].GetJsPlaceholderName()).IsEqualTo("startParagraph")
	AssertThat(normalMsgParts[4].GetString()).IsEqualTo("p1")
	AssertThat(normalMsgParts[5].GetJsPlaceholderName()).IsEqualTo("endParagraph")
}

/**
 * An ICU message using {@code phex} to describe its variables should result in the same
 * JavaScript code as the same ICU message without those {@code phex} attributes.
 */
// Java: XtbMessageBundleTest.testXtbBundle_icuPluralWithAndWithoutPhex()
func (c *XtbMessageBundleTest) TestXtbBundle_icuPluralWithAndWithoutPhex() {
	xtb :=
		"<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" +
			"<!DOCTYPE translationbundle SYSTEM \"translationbundle.dtd\">\n" +
			"<translationbundle lang=\"de_CH\">\n" +
			// With phex
			"<translation id=\"123456\">{NUM,plural, =1{Setting: {START_STRONG}<ph name=\"UDC_SETTING\"/>{END_STRONG} Products: {START_STRONG}<ph name=\"PRODUCT_LIST\"/>{END_STRONG}.}other{Settings: {START_STRONG}<ph name=\"UDC_SETTING_LIST\"/>{END_STRONG} Products: {START_STRONG}<ph name=\"PRODUCT_LIST\"/>{END_STRONG}.}}</translation>\n" +
			// Without phex
			"<translation id=\"987654\">{NUM,plural, =1{Setting: {START_STRONG}{UDC_SETTING}{END_STRONG} Products: {START_STRONG}{PRODUCT_LIST}{END_STRONG}.}other{Settings: {START_STRONG}{UDC_SETTING_LIST}{END_STRONG} Products: {START_STRONG}{PRODUCT_LIST}{END_STRONG}.}}</translation>\n" +
			"</translationbundle>"
	stream := bytes.NewReader([]byte(xtb))
	bundle := NewXtbMessageBundle(stream, jutil.StringPtr(xtbMessageBundleTest_PROJECT_ID))

	AssertThat(slices.Collect(bundle.GetAllMessages())).HasSize(2)
	AssertThat(bundle.GetMessage("123456").AsIcuMessageString()).
		IsEqualTo(
			"{NUM,plural, =1{Setting: {START_STRONG}{UDC_SETTING}{END_STRONG} Products: {START_STRONG}{PRODUCT_LIST}{END_STRONG}.}other{Settings: {START_STRONG}{UDC_SETTING_LIST}{END_STRONG} Products: {START_STRONG}{PRODUCT_LIST}{END_STRONG}.}}")
	// Both translation entries should result into the same message in JavaScript.
	AssertThat(bundle.GetMessage("987654").AsIcuMessageString()).
		IsEqualTo(bundle.GetMessage("123456").AsIcuMessageString())
}
