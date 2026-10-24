// Port of com/google/javascript/jscomp/ReplaceMessagesTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/ReplaceMessagesTest.java)

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
	"iter"
	"strconv"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ReplaceMessagesTest is a test which checks that replacer works correctly.
//
// Java: public final class ReplaceMessagesTest extends CompilerTestCase
type ReplaceMessagesTest struct {
	CompilerTestCaseBase

	// Messages returned from fake bundle, keyed by `JsMessage.id`.
	messages *jutil.JavaHashMap[string, *JsMessage]
	// If `true` report errors for messages that are not found in the bundle.
	strictReplacement bool
	// If `true` pass TEST_ID_GENERATOR in to ReplaceMessages via the fake bundle, so it will be
	// used to calculate the message IDs from the meaning and parts instead of just using the message
	// key as its id.
	useTestIdGenerator bool
	testMode           ReplaceMessagesTest_TestMode
}

// Java: ReplaceMessagesTest()
func NewReplaceMessagesTest() *ReplaceMessagesTest {
	c := &ReplaceMessagesTest{testMode: ReplaceMessagesTest_TestMode_FULL_REPLACE}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestReplaceMessagesTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewReplaceMessagesTest() })
}

// Go-only: language-adapter — the anonymous JsMessage.IdGenerator of
// ReplaceMessagesTest.TEST_ID_GENERATOR.
type replaceMessagesTest_testIdGenerator struct{}

// Java: ReplaceMessagesTest.TEST_ID_GENERATOR.generateId(String, List<Part>)
func (replaceMessagesTest_testIdGenerator) GenerateId(meaning *string, messageParts []JsMessage_Part) string {
	var idBuilder strings.Builder
	if meaning == nil {
		idBuilder.WriteString("null")
	} else {
		idBuilder.WriteString(*meaning)
	}
	idBuilder.WriteByte('_')
	for _, messagePart := range messageParts {
		if messagePart.IsPlaceholder() {
			idBuilder.WriteByte('p')
			idBuilder.WriteString(strconv.Itoa(jutil.StringLength(messagePart.GetCanonicalPlaceholderName())))
		} else {
			idBuilder.WriteByte('s')
			idBuilder.WriteString(strconv.Itoa(jutil.StringLength(messagePart.GetString())))
		}
	}

	return idBuilder.String()
}

// ReplaceMessagesTest_TEST_ID_GENERATOR generates IDs of the form
// `MEANING_PARTCOUNT[PARTCOUNT...]`:
// PARTCOUNT = 'sN' for a string part with N == string length;
// PARTCOUNT = 'pN' for a placeholder with N == length of the canonical placeholder name.
//
// Java: ReplaceMessagesTest.TEST_ID_GENERATOR
var ReplaceMessagesTest_TEST_ID_GENERATOR JsMessage_IdGenerator = replaceMessagesTest_testIdGenerator{}

// ReplaceMessagesTest_TestMode indicates which part of the replacement we're currently testing.
//
// Java: enum ReplaceMessagesTest.TestMode
type ReplaceMessagesTest_TestMode int

const (
	// Test full replacement from `goog.getMsg()` to final message values.
	ReplaceMessagesTest_TestMode_FULL_REPLACE ReplaceMessagesTest_TestMode = iota
	// Test replacement of `goog.getMsg()` with the protected function call form.
	// e.g.
	// ```javascript
	// var MSG_G =
	//     __jscomp_define_msg__(
	//         {
	//           "alt_id": null,
	//           "key":    "MSG_G",
	//           "msg_text": "${$amount}"
	//         },
	//         {amount:x + ""});`
	// ```
	ReplaceMessagesTest_TestMode_PROTECT_MSGS
	// Test replacement of the protected function call form with the final message values.
	ReplaceMessagesTest_TestMode_REPLACE_PROTECTED_MSGS
)

// ReplaceMessagesTest_SimpleMessageBundle is the fake bundle backed by the test's messages.
//
// Java: private class ReplaceMessagesTest.SimpleMessageBundle implements MessageBundle
type ReplaceMessagesTest_SimpleMessageBundle struct {
	outer *ReplaceMessagesTest
}

// Java: ReplaceMessagesTest.SimpleMessageBundle()
func NewReplaceMessagesTest_SimpleMessageBundle(outer *ReplaceMessagesTest) *ReplaceMessagesTest_SimpleMessageBundle {
	return &ReplaceMessagesTest_SimpleMessageBundle{outer: outer}
}

// Java: ReplaceMessagesTest.SimpleMessageBundle.getMessage(String)
func (b *ReplaceMessagesTest_SimpleMessageBundle) GetMessage(id string) *JsMessage {
	return b.outer.messages.Get(id)
}

// Java: ReplaceMessagesTest.SimpleMessageBundle.getAllMessages()
func (b *ReplaceMessagesTest_SimpleMessageBundle) GetAllMessages() iter.Seq[*JsMessage] {
	return b.outer.messages.Values()
}

// Java: ReplaceMessagesTest.SimpleMessageBundle.idGenerator()
func (b *ReplaceMessagesTest_SimpleMessageBundle) IdGenerator() JsMessage_IdGenerator {
	if b.outer.useTestIdGenerator {
		return ReplaceMessagesTest_TEST_ID_GENERATOR
	}
	return nil
}

// Java: ReplaceMessagesTest.getProcessor(Compiler)
func (c *ReplaceMessagesTest) GetProcessor(compiler *Compiler) CompilerPass {
	replaceMessages := NewReplaceMessages(compiler, NewReplaceMessagesTest_SimpleMessageBundle(c), c.strictReplacement)
	switch c.testMode {
	case ReplaceMessagesTest_TestMode_FULL_REPLACE:
		return replaceMessages.GetFullReplacementPass()
	case ReplaceMessagesTest_TestMode_PROTECT_MSGS:
		return replaceMessages.GetMsgProtectionPass()
	case ReplaceMessagesTest_TestMode_REPLACE_PROTECTED_MSGS:
		return replaceMessages.GetReplacementCompletionPass()
	}
	// Go-only: unreachable, the Java switch expression is exhaustive.
	panic(jutil.NewIllegalStateException())
}

// Java: ReplaceMessagesTest.multiPhaseTest(String, String, String)
func (c *ReplaceMessagesTest) MultiPhaseTest(originalJs string, protectedJs string, expectedJs string) {
	// The PROTECT_MSGS mode needs to add externs for the protection functions.
	c.AllowExternsChanges()
	c.testMode = ReplaceMessagesTest_TestMode_FULL_REPLACE
	c.Test(originalJs, expectedJs)
	c.testMode = ReplaceMessagesTest_TestMode_PROTECT_MSGS
	c.Test(originalJs, protectedJs)
	c.testMode = ReplaceMessagesTest_TestMode_REPLACE_PROTECTED_MSGS
	c.Test(protectedJs, expectedJs)
}

// Java: ReplaceMessagesTest.multiPhaseTestWarning(String, String, String, DiagnosticType)
func (c *ReplaceMessagesTest) MultiPhaseTestWarning(originalJs string, protectedJs string, expectedJs string, diagnosticType *DiagnosticType) {
	// The PROTECT_MSGS mode needs to add externs for the protection functions.
	c.AllowExternsChanges()
	c.testMode = ReplaceMessagesTest_TestMode_FULL_REPLACE
	c.TestWarning(originalJs, diagnosticType)
	c.testMode = ReplaceMessagesTest_TestMode_PROTECT_MSGS
	c.TestWarning(originalJs, diagnosticType)
	c.testMode = ReplaceMessagesTest_TestMode_REPLACE_PROTECTED_MSGS
	c.Test(protectedJs, expectedJs)
}

// Java: ReplaceMessagesTest.multiPhaseTestPreLookupError(String, DiagnosticType)
func (c *ReplaceMessagesTest) MultiPhaseTestPreLookupError(originalJs string, diagnosticType *DiagnosticType) {
	// The PROTECT_MSGS mode needs to add externs for the protection functions.
	c.AllowExternsChanges()
	c.testMode = ReplaceMessagesTest_TestMode_FULL_REPLACE
	c.TestError(originalJs, diagnosticType)
	c.testMode = ReplaceMessagesTest_TestMode_PROTECT_MSGS
	c.TestError(originalJs, diagnosticType)
}

// Java: ReplaceMessagesTest.multiPhaseTestPostLookupError(String, String, DiagnosticType)
func (c *ReplaceMessagesTest) MultiPhaseTestPostLookupError(originalJs string, protectedJs string, diagnosticType *DiagnosticType) {
	// The PROTECT_MSGS mode needs to add externs for the protection functions.
	c.AllowExternsChanges()
	c.testMode = ReplaceMessagesTest_TestMode_FULL_REPLACE
	c.TestError(originalJs, diagnosticType)
	c.testMode = ReplaceMessagesTest_TestMode_PROTECT_MSGS
	c.Test(originalJs, protectedJs)
	c.testMode = ReplaceMessagesTest_TestMode_REPLACE_PROTECTED_MSGS
	c.TestError(protectedJs, diagnosticType)
}

// Java: ReplaceMessagesTest.multiPhaseTestPostLookupError(String, String, DiagnosticType, String)
func (c *ReplaceMessagesTest) MultiPhaseTestPostLookupErrorWithDescription(originalJs string, protectedJs string, diagnosticType *DiagnosticType, description string) {
	// The PROTECT_MSGS mode needs to add externs for the protection functions.
	c.AllowExternsChanges()
	c.testMode = ReplaceMessagesTest_TestMode_FULL_REPLACE
	c.TestError(originalJs, diagnosticType, description)
	c.testMode = ReplaceMessagesTest_TestMode_PROTECT_MSGS
	c.Test(originalJs, protectedJs)
	c.testMode = ReplaceMessagesTest_TestMode_REPLACE_PROTECTED_MSGS
	c.TestError(protectedJs, diagnosticType, description)
}

// Java: ReplaceMessagesTest.setUp()
func (c *ReplaceMessagesTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.messages = jutil.NewJavaHashMap[string, *JsMessage]()
	c.strictReplacement = false
	c.useTestIdGenerator = false
	c.EnableTypeCheck()
	c.ReplaceTypesWithColors()
	c.EnableTypeInfoValidation()
	c.IgnoreWarnings(DiagnosticGroups_MISSING_PROPERTIES)
}

// Java: ReplaceMessagesTest.testProtectedMessagesAndFallbackAreRemovable()
func (c *ReplaceMessagesTest) TestProtectedMessagesAndFallbackAreRemovable() {
	// Externs for the protection functions will be added.
	c.AllowExternsChanges()
	c.testMode = ReplaceMessagesTest_TestMode_PROTECT_MSGS
	c.Test(
		`/** @desc d */
var MSG_A = goog.getMsg('asdf');
/** @desc d */
var MSG_B = goog.getMsg('qwerty');
var x = goog.getMsgWithFallback(MSG_A, MSG_B);
`,
		`/** @desc d */
var MSG_A =
    __jscomp_define_msg__(
        {
          "key":    "MSG_A",
          "msg_text":"asdf",
        });
/** @desc d */
var MSG_B =
    __jscomp_define_msg__(
        {
          "key":    "MSG_B",
          "msg_text":"qwerty",
      });
var x = __jscomp_msg_fallback__("MSG_A", MSG_A, "MSG_B", MSG_B);
`)
	// It's important that all of the protective function calls be marked as having no side effects,
	// so they will be removed during optimizations if they are unused.
	jsRoot := c.GetLastCompiler().GetJsRoot()
	script := jsRoot.GetFirstChild()
	msgAVarNode := script.GetFirstChild()

	msgANameNode := msgAVarNode.GetOnlyChild()
	AssertNode(msgANameNode).IsName("MSG_A")
	msgACallNode := msgANameNode.GetOnlyChild()
	AssertThat(msgACallNode.GetSideEffectFlags()).IsEqualTo(rhino.Node_SideEffectFlags_NO_SIDE_EFFECTS)

	msgBVarNode := msgAVarNode.GetNext()
	msgBNameNode := msgBVarNode.GetOnlyChild()
	AssertNode(msgBNameNode).IsName("MSG_B")
	msgBCallNode := msgANameNode.GetOnlyChild()
	AssertThat(msgBCallNode.GetSideEffectFlags()).IsEqualTo(rhino.Node_SideEffectFlags_NO_SIDE_EFFECTS)

	xVarNode := msgBVarNode.GetNext()
	xNameNode := xVarNode.GetOnlyChild()
	AssertNode(xNameNode).IsName("x")
	xCallNode := xNameNode.GetOnlyChild()
	AssertThat(xCallNode.GetSideEffectFlags()).IsEqualTo(rhino.Node_SideEffectFlags_NO_SIDE_EFFECTS)
}

// Java: ReplaceMessagesTest.testReplaceSimpleMessage()
func (c *ReplaceMessagesTest) TestReplaceSimpleMessage() {
	c.RegisterMessage(c.GetTestMessageBuilder("MSG_A").AppendStringPart("Hi\nthere").Build())

	c.MultiPhaseTest(
		`/** @desc d */
var MSG_A = goog.getMsg('asdf');
`,
		`/**
 * @desc d
 */
var MSG_A =
    __jscomp_define_msg__(
        {
          "key":    "MSG_A",
          "msg_text":"asdf",
        });
`,
		`/** @desc d */
var MSG_A='Hi\nthere'
`)
}

// Java: ReplaceMessagesTest.testReplaceExternalMessage()
func (c *ReplaceMessagesTest) TestReplaceExternalMessage() {
	c.RegisterMessage(c.GetTestMessageBuilder("12345").AppendStringPart("Saluton!").Build())

	c.MultiPhaseTest(
		`/** @desc d */
var MSG_EXTERNAL_12345 = goog.getMsg('Hello!');
`,
		`/**
 * @desc d
 */
var MSG_EXTERNAL_12345 =
    __jscomp_define_msg__(
        {
          "key":    "MSG_EXTERNAL_12345",
          "msg_text":"Hello!",
        });
`,
		`/** @desc d */
var MSG_EXTERNAL_12345='Saluton!'
`)
}

// Java: ReplaceMessagesTest.testReplaceIcuTemplateMessageWithBundleAndJsPlaceholders()
func (c *ReplaceMessagesTest) TestReplaceIcuTemplateMessageWithBundleAndJsPlaceholders() {
	// This unit test contains an ICU template with placeholders ("{EMAIL}"). We cannot treat this
	// message as a single string part, because it has multiple parts. Otherwise, we will generate
	// the wrong message id and this unit test will fail.
	c.useTestIdGenerator = true
	c.strictReplacement = true

	meaning := "MSG_SHOW_EMAIL"
	originalStringPart := JsMessage_StringPart_Create("Email: ")
	originalPlaceholerPart := JsMessage_PlaceholderReference_CreateForCanonicalName("EMAIL")
	expectedMessageId :=
		ReplaceMessagesTest_TEST_ID_GENERATOR.GenerateId(
			&meaning, []JsMessage_Part{originalStringPart, originalPlaceholerPart})

	// Create and register the translation we expect to find in the message bundle
	showEmailTranslatedMsg :=
		NewJsMessage_Builder().
			SetKey(meaning).
			SetMeaning(&meaning).
			AppendStringPart("Retpoŝtadreso: "). // translated string.
			AppendPart(originalPlaceholerPart).  // placeholder is the same as original.
			SetId(expectedMessageId).            // message ID was calculated from the original.
			Build()
	c.RegisterMessage(showEmailTranslatedMsg)

	c.MultiPhaseTest(
		`const {declareIcuTemplate} = goog.require('goog.i18n.messages');

const MSG_SHOW_EMAIL =
    declareIcuTemplate(
        'Email: {EMAIL}',
        {
          description: 'Labeled email address',
// The example text is dropped, since it is only used for XMB extraction.
// However, it does cause the JsMessage read from the JS code to have a placeholder
// in it.
          example: {
            'EMAIL': 'me@foo.com'
           }
        });
`,
		`const {declareIcuTemplate} = goog.require('goog.i18n.messages');

const MSG_SHOW_EMAIL =
    __jscomp_define_msg__(
        {
          "key":    "MSG_SHOW_EMAIL",
          "icu_placeholder_names": ["EMAIL"],
          "msg_text": "Email: {EMAIL}",
          "isIcuTemplate": ""
        });
`,
		`const {declareIcuTemplate} = goog.require('goog.i18n.messages');

const MSG_SHOW_EMAIL = 'Retpoŝtadreso: {EMAIL}';
`)
}

// Java: ReplaceMessagesTest.testReplaceIcuTemplateMessageWithoutJsPlaceholders()
func (c *ReplaceMessagesTest) TestReplaceIcuTemplateMessageWithoutJsPlaceholders() {
	// Message in the bundle has a placeholder and is NOT in ICU selector format.
	//
	// (i.e. it does not start with "{WORD,").
	//
	// Here we want to make sure that messages created with declareIcuTemplate()
	// get treated as ICU messages even without that distinguishing feature.
	c.RegisterMessage(
		c.GetTestMessageBuilder("MSG_SHOW_EMAIL").
			AppendStringPart("Retpoŝtadreso: ").
			AppendCanonicalPlaceholderReference("EMAIL").
			Build())

	c.MultiPhaseTest(
		`const {declareIcuTemplate} = goog.require('goog.i18n.messages');

// Note that no placeholder information is specified here, so the JsMessage as it is
// read from the JS code will have no placeholders.
// In this test case we've put a placeholder in the bundle above, but if the bundle
// were created based on this code, it would not have a placeholder.
// This situation could occur with an externally-produced message bundle.
const MSG_SHOW_EMAIL = declareIcuTemplate(
    'Email: {EMAIL}', { description: 'Labeled email address' });
`,
		`const {declareIcuTemplate} = goog.require('goog.i18n.messages');

const MSG_SHOW_EMAIL =
    __jscomp_define_msg__(
        {
          "key":    "MSG_SHOW_EMAIL",
          "msg_text": "Email: {EMAIL}",
          "isIcuTemplate": ""
        });
`,
		`const {declareIcuTemplate} = goog.require('goog.i18n.messages');

const MSG_SHOW_EMAIL = 'Retpoŝtadreso: {EMAIL}';
`)
}

// Java: ReplaceMessagesTest.testReplaceIcuTemplateMessageWithJsPlaceholders()
func (c *ReplaceMessagesTest) TestReplaceIcuTemplateMessageWithJsPlaceholders() {
	// Make sure ICU messages with multiple parts are handled correctly.
	// We cannot treat this ICU message as a single string part, because it has two placeholders
	// (EMAIL1 and EMAIL2) that cause the message to be split into multiple parts.
	c.RegisterMessage(
		c.GetTestMessageBuilder("MSG_SHOW_EMAIL").
			AppendStringPart("Retpoŝtadreso: ").
			AppendCanonicalPlaceholderReference("EMAIL_1").
			AppendStringPart(" aŭ ").
			AppendCanonicalPlaceholderReference("EMAIL_2").
			Build())

	c.MultiPhaseTest(
		`const {declareIcuTemplate} = goog.require('goog.i18n.messages');

const MSG_SHOW_EMAIL =
    declareIcuTemplate(
        'Email Options: {EMAIL_1} or {EMAIL_2}',
        {
          description: 'Labeled email address',
          example: {
            'EMAIL_1': 'me1@foo.com',
            'EMAIL_2': 'me2@foo.com'
           }
        });
`,
		`const {declareIcuTemplate} = goog.require('goog.i18n.messages');

const MSG_SHOW_EMAIL =
    __jscomp_define_msg__(
        {
          "key":    "MSG_SHOW_EMAIL",
          "icu_placeholder_names": ["EMAIL_1", "EMAIL_2"],
          "msg_text": "Email Options: {EMAIL_1} or {EMAIL_2}",
          "isIcuTemplate": ""
        });
`,
		`const {declareIcuTemplate} = goog.require('goog.i18n.messages');

const MSG_SHOW_EMAIL = 'Retpoŝtadreso: {EMAIL_1} aŭ {EMAIL_2}';
`)
}

// Java: ReplaceMessagesTest.testMissingIcuTemplateMessage()
func (c *ReplaceMessagesTest) TestMissingIcuTemplateMessage() {
	// We don't registerMessage() here, so there are no messages in the bundle used by this test.

	c.MultiPhaseTest(
		"const {declareIcuTemplate} = goog.require('goog.i18n.messages');\n\nconst MSG_SHOW_EMAIL =\n    declareIcuTemplate(\n        'Email: {EMAIL}',\n        {\n          description: 'Labeled email address',\n// The example text is dropped, since it is only used for XMB extraction.\n// However, it does cause the JsMessage read from the JS code to have a placeholder\n// in it.\n// We add this placeholder in the \"icu_placeholder_names\" field to keep track of how the\n// message has multiple parts, which is necessary for the message ID to be generated\n// correctly.\n// The purpose of this test is to:\n// 1. make sure the message template is properly put back together.\n// 2. make sure the \"icu_placeholder_names\" field is populated with `EMAIL`\n          example: {\n            'EMAIL': 'me@foo.com'\n           }\n        });\n",
		`const {declareIcuTemplate} = goog.require('goog.i18n.messages');

const MSG_SHOW_EMAIL =
    __jscomp_define_msg__(
        {
          "key":    "MSG_SHOW_EMAIL",
          "icu_placeholder_names": ["EMAIL"],
          "msg_text": "Email: {EMAIL}",
          "isIcuTemplate": ""
        });
`,
		`const {declareIcuTemplate} = goog.require('goog.i18n.messages');

const MSG_SHOW_EMAIL = 'Email: {EMAIL}';
`)
}

// Java: ReplaceMessagesTest.testReplaceExternalIcuSelectorMessageWithPlaceholders()
func (c *ReplaceMessagesTest) TestReplaceExternalIcuSelectorMessageWithPlaceholders() {
	// Message in the bundle is in ICU selector format with has placeholders with explicit
	// placeholders.
	// The JS code treats the message as a simple string without placeholders.
	// The compiler should join the placeholder names together with the string parts in order to
	// get the runtime string value.
	c.RegisterMessage(
		c.GetTestMessageBuilder("123456").
			AppendStringPart("{USER_GENDER,select,female{Saluton ").
			AppendCanonicalPlaceholderReference("USER_IDENTIFIER").
			AppendStringPart(".}male{Saluton ").
			AppendCanonicalPlaceholderReference("USER_IDENTIFIER").
			AppendStringPart(".}other{Saluton ").
			AppendCanonicalPlaceholderReference("USER_IDENTIFIER").
			AppendStringPart(".}}").
			Build())

	c.MultiPhaseTest(
		`/** @desc ICU gender-sensitive greeting */
// Message in the JS code does not define placeholders for the compiler.
const MSG_EXTERNAL_123456 = goog.getMsg(
    '{USER_GENDER,select,' +
    'female{Hello {USER_IDENTIFIER}.}' +
    'male{Hello {USER_IDENTIFIER}.}' +
    'other{Hello {USER_IDENTIFIER}.}}');
`,
		`/** @desc ICU gender-sensitive greeting */
const MSG_EXTERNAL_123456 =
    __jscomp_define_msg__(
        {
          "key":    "MSG_EXTERNAL_123456",
          "msg_text":
    '{USER_GENDER,select,female{Hello {USER_IDENTIFIER}.}male{Hello {USER_IDENTIFIER}.}other{Hello {USER_IDENTIFIER}.}}',
        });
`,
		`/** @desc ICU gender-sensitive greeting */
const MSG_EXTERNAL_123456 =
    '{USER_GENDER,select,female{Saluton {USER_IDENTIFIER}.}male{Saluton {USER_IDENTIFIER}.}other{Saluton {USER_IDENTIFIER}.}}';
`)
}

// Java: ReplaceMessagesTest.testReplaceSimpleMessageDefinedWithAdd()
func (c *ReplaceMessagesTest) TestReplaceSimpleMessageDefinedWithAdd() {
	c.RegisterMessage(c.GetTestMessageBuilder("MSG_A").AppendStringPart("Hi\nthere").Build())

	c.MultiPhaseTest(
		`/** @desc d */
var MSG_A = goog.getMsg('abcd' + 'efgh');
`,
		`/**
 * @desc d
 */
var MSG_A =
    __jscomp_define_msg__(
        {
          "key":"MSG_A",
          "msg_text":"abcdefgh",
        });
`,
		`/** @desc d */
var MSG_A='Hi\nthere'
`)
}

// Java: ReplaceMessagesTest.testMissingAlternateMessage()
func (c *ReplaceMessagesTest) TestMissingAlternateMessage() {
	c.MultiPhaseTest(
		`/**
 * @desc d
 * @alternateMessageId 1984
 */
var MSG_A = goog.getMsg('asdf');
`,
		`/**
 * @desc d
 * @alternateMessageId 1984
 */
var MSG_A =
    __jscomp_define_msg__(
        {
          "key":"MSG_A",
          "alt_id":"1984",
          "msg_text":"asdf",
        });
`,
		`/**
 * @desc d
 * @alternateMessageId 1984
 */
var MSG_A='asdf'
`)
}

// Java: ReplaceMessagesTest.testAlternateMessageWithMismatchedParts()
func (c *ReplaceMessagesTest) TestAlternateMessageWithMismatchedParts() {
	c.RegisterMessage(
		c.GetTestMessageBuilder("1984").
			SetDesc(jutil.StringPtr("B desc")).
			SetMeaning(jutil.StringPtr("B meaning")).
			AppendStringPart("Hello!").
			AppendStringPart(" Welcome!").
			Build())

	c.MultiPhaseTest(
		`/**
 * @desc d
 * @alternateMessageId 1984
 */
var MSG_A = goog.getMsg('asdf');
`,
		`/**
 * @desc d
 * @alternateMessageId 1984
 */
var MSG_A =
    __jscomp_define_msg__(
        {
          "key":"MSG_A",
          "alt_id":"1984",
          "msg_text":"asdf",
        });
`,
		`/**
 * @desc d
 * @alternateMessageId 1984
 */
var MSG_A = 'Hello! Welcome!';
`)
}

// Java: ReplaceMessagesTest.testAlternateIcuSelectorMessageWithPlaceholders()
func (c *ReplaceMessagesTest) TestAlternateIcuSelectorMessageWithPlaceholders() {
	// Message in the bundle is in ICU selector format with has placeholders with explicit
	// placeholders.
	// The JS code treats the message as a simple string without placeholders.
	// The compiler should join the placeholder names together with the string parts in order to
	// get the runtime string value.
	// Note that we are not putting a translation for the actual message from the JS Code into the
	// bundle here. Instead, we are providing the alternate message.
	c.RegisterMessage(
		c.GetTestMessageBuilder("1984").
			AppendStringPart("{USER_GENDER,select,female{Saluton ").
			AppendCanonicalPlaceholderReference("USER_IDENTIFIER").
			AppendStringPart(".}male{Saluton ").
			AppendCanonicalPlaceholderReference("USER_IDENTIFIER").
			AppendStringPart(".}other{Saluton ").
			AppendCanonicalPlaceholderReference("USER_IDENTIFIER").
			AppendStringPart(".}}").
			Build())

	c.MultiPhaseTest(
		`/**
 * @desc ICU gender-sensitive greeting
 * @alternateMessageId 1984
 */
// Message in the JS code does not define placeholders for the compiler.
const MSG_ICU_SELECT = goog.getMsg(
    '{USER_GENDER,select,' +
    'female{Hello {USER_IDENTIFIER}.}' +
    'male{Hello {USER_IDENTIFIER}.}' +
    'other{Hello {USER_IDENTIFIER}.}}');
`,
		`/**
 * @desc ICU gender-sensitive greeting
 * @alternateMessageId 1984
 */
const MSG_ICU_SELECT =
    __jscomp_define_msg__(
        {
          "key":    "MSG_ICU_SELECT",
          "alt_id": "1984",
          "msg_text":
    '{USER_GENDER,select,female{Hello {USER_IDENTIFIER}.}male{Hello {USER_IDENTIFIER}.}other{Hello {USER_IDENTIFIER}.}}',
        });
`,
		`/**
 * @desc ICU gender-sensitive greeting
 * @alternateMessageId 1984
 */
const MSG_ICU_SELECT =
    '{USER_GENDER,select,female{Saluton {USER_IDENTIFIER}.}male{Saluton {USER_IDENTIFIER}.}other{Saluton {USER_IDENTIFIER}.}}';
`)
}

// Java: ReplaceMessagesTest.getTestMessageBuilder(String)
func (c *ReplaceMessagesTest) GetTestMessageBuilder(keyAndId string) *JsMessage_Builder {
	return NewJsMessage_Builder().SetKey(keyAndId).SetId(keyAndId)
}

// Java: ReplaceMessagesTest.testAlternateMessageWithMismatchedPlaceholders()
func (c *ReplaceMessagesTest) TestAlternateMessageWithMismatchedPlaceholders() {
	c.RegisterMessage(
		c.GetTestMessageBuilder("1984").
			SetDesc(jutil.StringPtr("B desc")).
			SetMeaning(jutil.StringPtr("B meaning")).
			AppendStringPart("Hello, ").
			AppendJsPlaceholderReference("firstName").
			AppendStringPart("!").
			Build())

	c.MultiPhaseTestPostLookupError(
		`/**
 * @desc B desc
 * @meaning B meaning
 * @alternateMessageId 1984
 */
var MSG_A = goog.getMsg('Hello, {$name}!', {name: name});
`,
		`/**
 * @desc B desc
 * @meaning B meaning
 * @alternateMessageId 1984
 */
var MSG_A =
    __jscomp_define_msg__(
        {
          "key":    "MSG_A",
          "alt_id": "1984",
          "meaning":"B meaning",
          "msg_text":"Hello, {$name}!"
        },
        {'name': name});
`,
		ReplaceMessages_INVALID_ALTERNATE_MESSAGE_PLACEHOLDERS)
}

// Java: ReplaceMessagesTest.testReplaceValidAlternateMessage()
func (c *ReplaceMessagesTest) TestReplaceValidAlternateMessage() {
	c.RegisterMessage(c.GetTestMessageBuilder("1984").AppendStringPart("Howdy\npardner").Build())

	c.MultiPhaseTest(
		`/**
 * @desc B desc
 * @alternateMessageId 1984
 */
var MSG_A = goog.getMsg('asdf');
`,
		`/**
 * @desc B desc
 * @alternateMessageId 1984
 */
var MSG_A =
    __jscomp_define_msg__(
        {
          "key":"MSG_A",
          "alt_id":"1984",
          "msg_text":"asdf",
        });
`,
		`/**
 * @desc B desc
 * @alternateMessageId 1984
 */
var MSG_A='Howdy\npardner'
`)
}

// Java: ReplaceMessagesTest.testIgnoreUnnecessaryAlternateMessage()
func (c *ReplaceMessagesTest) TestIgnoreUnnecessaryAlternateMessage() {
	c.RegisterMessage(c.GetTestMessageBuilder("1984").AppendStringPart("Howdy\npardner").Build())
	c.RegisterMessage(
		c.GetTestMessageBuilder("MSG_A").
			SetDesc(jutil.StringPtr("Greeting.")).
			SetAlternateId(jutil.StringPtr("1984")).
			AppendStringPart("Hi\nthere").
			Build())

	c.MultiPhaseTest(
		`/**
 * @desc d
 * @alternateMessageId 1984
 */
var MSG_A = goog.getMsg('asdf');
`,
		`/**
 * @desc d
 * @alternateMessageId 1984
 */
var MSG_A =
    __jscomp_define_msg__(
        {
          "key":"MSG_A",
          "alt_id":"1984",
          "msg_text":"asdf",
        });
`,
		`/**
 * @desc d
 * @alternateMessageId 1984
 */
var MSG_A = 'Hi\nthere';
`)
}

// Java: ReplaceMessagesTest.testAlternateTrumpsFallback()
func (c *ReplaceMessagesTest) TestAlternateTrumpsFallback() {
	c.RegisterMessage(c.GetTestMessageBuilder("1984").AppendStringPart("Howdy\npardner").Build())

	c.RegisterMessage(c.GetTestMessageBuilder("MSG_B").AppendStringPart("Good\nmorrow, sir").Build())

	c.MultiPhaseTest(
		`/**
 * @desc d
 * @alternateMessageId 1984
*/
var MSG_A = goog.getMsg('asdf');
/**
 * @desc d
*/
var MSG_B = goog.getMsg('ghjk');
var x = goog.getMsgWithFallback(MSG_A, MSG_B);
`,
		`/**
 * @desc d
 * @alternateMessageId 1984
 */
var MSG_A =
    __jscomp_define_msg__(
        {
          "key":"MSG_A",
          "alt_id":"1984",
          "msg_text":"asdf",
        });
/**
 * @desc d
 */
var MSG_B =
    __jscomp_define_msg__(
        {
          "key":"MSG_B",
          "msg_text":"ghjk",
        });
var x = __jscomp_msg_fallback__("MSG_A", MSG_A, "MSG_B", MSG_B);
`,
		`/**
    @desc d
    @alternateMessageId 1984
*/
var MSG_A = 'Howdy\npardner';
/**
    @desc d
*/
var MSG_B = 'Good\nmorrow, sir';
var x = MSG_A;
`)
}

// Java: ReplaceMessagesTest.testFallbackWithAlternate()
func (c *ReplaceMessagesTest) TestFallbackWithAlternate() {
	c.RegisterMessage(c.GetTestMessageBuilder("1984").AppendStringPart("Howdy\npardner").Build())

	c.MultiPhaseTest(
		`/**
    @desc d
*/
var MSG_A = goog.getMsg('asdf');
/**
    @desc d
    @alternateMessageId 1984
*/
var MSG_B = goog.getMsg('ghjk');
var x = goog.getMsgWithFallback(MSG_A, MSG_B);
`,
		`/**
 * @desc d
 */
var MSG_A =
    __jscomp_define_msg__(
        {
          "key":"MSG_A",
          "msg_text":"asdf",
        });
/**
 * @desc d
    @alternateMessageId 1984     */
var MSG_B =
    __jscomp_define_msg__(
        {
          "key":"MSG_B",
          "alt_id":"1984",
          "msg_text":"ghjk",
        });
var x = __jscomp_msg_fallback__("MSG_A", MSG_A, "MSG_B", MSG_B);
`,
		`/**
    @desc d
*/
var MSG_A = 'asdf';
/**
    @desc d
    @alternateMessageId 1984
*/
var MSG_B = 'Howdy\npardner';
var x = MSG_B;
`)
}

// Java: ReplaceMessagesTest.testNameReplacement()
func (c *ReplaceMessagesTest) TestNameReplacement() {
	c.RegisterMessage(
		c.GetTestMessageBuilder("MSG_B").
			AppendStringPart("One ").
			AppendJsPlaceholderReference("measly").
			AppendStringPart(" ph").
			Build())

	c.MultiPhaseTest(
		`/** @desc d */
var MSG_B=goog.getMsg('asdf {$measly}', {measly: x});
`,
		`/**
 * @desc d
 */
var MSG_B =
    __jscomp_define_msg__(
        {
          "key":"MSG_B",
          "msg_text":"asdf {$measly}",
        }, {'measly': x});
`,
		`/** @desc d */
var MSG_B = 'One ' + x + ' ph';
`)
}

// Java: ReplaceMessagesTest.testNameReplacementWithFullOptionsBag()
func (c *ReplaceMessagesTest) TestNameReplacementWithFullOptionsBag() {
	c.RegisterMessage(
		c.GetTestMessageBuilder("MSG_B").
			AppendStringPart("One ").
			AppendJsPlaceholderReference("measly").
			AppendStringPart(" ph").
			Build())

	c.MultiPhaseTest(
		`/** @desc d */
var MSG_B =
    goog.getMsg(
        'asdf {$measly}',
        {measly: x},
        {
// use all allowed options
          html: true,
          unescapeHtmlEntities: true,
// original_code and example get dropped, because they're only used
// when generating the XMB file.
          original_code: {
            'measly': 'getMeasley()'
          },
          example: {
            'measly': 'very little'
          },
        });
`,
		`/**
 * @desc d
 */
var MSG_B =
    __jscomp_define_msg__(
        {
          "key":"MSG_B",
          "msg_text":"asdf {$measly}",
          "escapeLessThan":"",
          "unescapeHtmlEntities":""
        },
        {'measly': x});
`,
		`/** @desc d */
var MSG_B = 'One ' + x + ' ph';
`)
}

// Java: ReplaceMessagesTest.testGetPropReplacement()
func (c *ReplaceMessagesTest) TestGetPropReplacement() {
	c.RegisterMessage(c.GetTestMessageBuilder("MSG_C").AppendJsPlaceholderReference("amount").Build())

	c.MultiPhaseTest(
		`/** @desc d */
var MSG_C = goog.getMsg('${$amount}', {amount: a.b.amount});
`,
		`/**
 * @desc d
 */
var MSG_C =
    __jscomp_define_msg__(
        {
          "key":"MSG_C",
          "msg_text":"${$amount}",
        }, {'amount': a.b.amount});
`,
		`/** @desc d */
var MSG_C=a.b.amount
`)
}

// Java: ReplaceMessagesTest.testFunctionCallReplacement()
func (c *ReplaceMessagesTest) TestFunctionCallReplacement() {
	c.RegisterMessage(c.GetTestMessageBuilder("MSG_D").AppendJsPlaceholderReference("amount").Build())

	c.MultiPhaseTest(
		`/** @desc d */
var MSG_D = goog.getMsg('${$amount}', {amount: getAmt()});
`,
		`/**
 * @desc d
 */
var MSG_D =
    __jscomp_define_msg__(
        {
          "key":"MSG_D",
          "msg_text":"${$amount}",
        }, {'amount': getAmt()});
`,
		`/** @desc d */
var MSG_D=getAmt()
`)
}

// Java: ReplaceMessagesTest.testMethodCallReplacement()
func (c *ReplaceMessagesTest) TestMethodCallReplacement() {
	c.RegisterMessage(c.GetTestMessageBuilder("MSG_E").AppendJsPlaceholderReference("amount").Build())

	c.MultiPhaseTest(
		`/** @desc d */
var MSG_E = goog.getMsg('${$amount}', {amount: obj.getAmt()});
`,
		`/**
 * @desc d
 */
var MSG_E =
    __jscomp_define_msg__(
        {
          "key":"MSG_E",
          "msg_text":"${$amount}",
        }, {'amount': obj.getAmt()});
`,
		`/** @desc d */
var MSG_E=obj.getAmt()
`)
}

// TODO(b/538170233): ReplaceMessages preserving side effects in placeholders
//
// Java: ReplaceMessagesTest.testMethodCallReplacementEmptyMessage()
func (c *ReplaceMessagesTest) TestMethodCallReplacementEmptyMessage() {
	c.RegisterMessage(c.GetTestMessageBuilder("MSG_M").Build())

	c.MultiPhaseTest(
		`/** @desc d */
var MSG_M = goog.getMsg('${$amount}', {amount: obj.getAmt()});
`,
		`/**
 * @desc d
 */
var MSG_M =
    __jscomp_define_msg__(
        {
          "key":"MSG_M",
          "msg_text":"${$amount}",
        }, {'amount': obj.getAmt()});
`,
		`/** @desc d */
var MSG_M = "";
`)
}

// Java: ReplaceMessagesTest.testHookReplacement()
func (c *ReplaceMessagesTest) TestHookReplacement() {
	c.RegisterMessage(
		c.GetTestMessageBuilder("MSG_F").
			AppendStringPart("#").
			AppendJsPlaceholderReference("amount").
			AppendStringPart(".").
			Build())

	c.MultiPhaseTest(
		"/** @desc d */\n var MSG_F = goog.getMsg('${$amount}', {amount: (a ? b : c)});",
		`/**
 * @desc d
 */
var MSG_F =
    __jscomp_define_msg__(
        {
          "key":"MSG_F",
          "msg_text":"${$amount}",
        }, {'amount': a ? b : c});
`,
		`/** @desc d */
var MSG_F = '#' + (a?b:c) + '.';
`)
}

// Java: ReplaceMessagesTest.testAddReplacement()
func (c *ReplaceMessagesTest) TestAddReplacement() {
	c.RegisterMessage(c.GetTestMessageBuilder("MSG_G").AppendJsPlaceholderReference("amount").Build())

	c.MultiPhaseTest(
		`/** @desc d */
var MSG_G = goog.getMsg('${$amount}', {amount: x + ''});
`,
		`/** @desc d */
var MSG_G =
    __jscomp_define_msg__(
        {
          'key':'MSG_G',
          "msg_text":"${$amount}",
        },
        {'amount': x + ''});
`,
		`/** @desc d */
var MSG_G=x+''
`)
}

// Java: ReplaceMessagesTest.testPlaceholderValueReferencedTwice()
func (c *ReplaceMessagesTest) TestPlaceholderValueReferencedTwice() {
	c.RegisterMessage(
		c.GetTestMessageBuilder("MSG_H").
			AppendJsPlaceholderReference("dick").
			AppendStringPart(", ").
			AppendJsPlaceholderReference("dick").
			AppendStringPart(" and ").
			AppendJsPlaceholderReference("jane").
			Build())

	c.MultiPhaseTest(
		"/** @desc d */\n var MSG_H = goog.getMsg('{$dick}{$jane}', {jane: x, dick: y});",
		`/**
 * @desc d
 */
var MSG_H =
    __jscomp_define_msg__(
        {
          "key":"MSG_H",
          "msg_text":"{$dick}{$jane}",
        }, {'jane': x, 'dick': y});
`,
		`/** @desc d */
var MSG_H = y + ', ' + y + ' and ' + x;
`)
}

// Java: ReplaceMessagesTest.testInvalidMessageStringType()
func (c *ReplaceMessagesTest) TestInvalidMessageStringType() {
	c.MultiPhaseTestPreLookupError(
		`/** @desc d */
const MSG_H = goog.getMsg(10);
`,
		JsMessageVisitor_MESSAGE_TREE_MALFORMED)
}

// Java: ReplaceMessagesTest.testPlaceholderValueDefinedTwice()
func (c *ReplaceMessagesTest) TestPlaceholderValueDefinedTwice() {
	c.MultiPhaseTestPreLookupError(
		`/** @desc d */
const MSG_H = goog.getMsg(
    '{$dick}{$jane}',
    {jane: x, dick: y, jane: x});
`,
		JsMessageVisitor_MESSAGE_TREE_MALFORMED)
}

// Java: ReplaceMessagesTest.testInvalidPlaceholderArgument()
func (c *ReplaceMessagesTest) TestInvalidPlaceholderArgument() {
	c.MultiPhaseTestPreLookupError(
		`/** @desc d */
const MSG_H = goog.getMsg(
    '{$dick}{$jane}',
    'this should be an object literal');
`,
		JsMessageVisitor_MESSAGE_TREE_MALFORMED)
}

// Java: ReplaceMessagesTest.testInvalidOptionsArgumentType()
func (c *ReplaceMessagesTest) TestInvalidOptionsArgumentType() {
	c.MultiPhaseTestPreLookupError(
		`/** @desc d */
const MSG_H = goog.getMsg(
    '{$dick}{$jane}',
    {jane: x, dick: y},
    'should be an object literal');
`,
		JsMessageVisitor_MESSAGE_TREE_MALFORMED)
}

// Java: ReplaceMessagesTest.testComputedKeyInOptions()
func (c *ReplaceMessagesTest) TestComputedKeyInOptions() {
	c.MultiPhaseTestPreLookupError(
		`/** @desc d */
const MSG_H = goog.getMsg(
    '{$dick}{$jane}',
    {jane: x, dick: y},
    {[computedOpt]: true});
`,
		JsMessageVisitor_MESSAGE_TREE_MALFORMED)
}

// Java: ReplaceMessagesTest.testPlaceholderNameInLowerCamelCase()
func (c *ReplaceMessagesTest) TestPlaceholderNameInLowerCamelCase() {
	c.RegisterMessage(
		c.GetTestMessageBuilder("MSG_I").
			AppendStringPart("Sum: $").
			AppendJsPlaceholderReference("amtEarned").
			Build())

	c.MultiPhaseTest(
		`/** @desc d */
var MSG_I = goog.getMsg('${$amtEarned}', {amtEarned: x});
`,
		`/**
 * @desc d
 */
var MSG_I =
    __jscomp_define_msg__(
        {
          "key":"MSG_I",
          "msg_text":"${$amtEarned}",
        }, {'amtEarned': x});
`,
		`/** @desc d */
var MSG_I='Sum: $'+x
`)
}

// Java: ReplaceMessagesTest.testQualifiedMessageName()
func (c *ReplaceMessagesTest) TestQualifiedMessageName() {
	c.RegisterMessage(
		c.GetTestMessageBuilder("MSG_J").
			AppendStringPart("One ").
			AppendJsPlaceholderReference("measly").
			AppendStringPart(" ph").
			Build())

	c.MultiPhaseTest(
		`/** @desc d */
a.b.c.MSG_J = goog.getMsg('asdf {$measly}', {measly: x});
`,
		`/**
 * @desc d
 */
    a.b.c.MSG_J =
    __jscomp_define_msg__(
        {
          "key":   'MSG_J',
          "msg_text":"asdf {$measly}",
        },
        {'measly': x});
`,
		`/** @desc d */
a.b.c.MSG_J = 'One ' + x + ' ph';
`)
}

// Java: ReplaceMessagesTest.testPlaceholderInPlaceholderValue()
func (c *ReplaceMessagesTest) TestPlaceholderInPlaceholderValue() {
	c.RegisterMessage(
		c.GetTestMessageBuilder("MSG_L").
			AppendJsPlaceholderReference("a").
			AppendStringPart(" has ").
			AppendJsPlaceholderReference("b").
			Build())

	c.MultiPhaseTest(
		`/** @desc d */
var MSG_L = goog.getMsg('{$a} has {$b}', {a: '{$b}', b: 1});
`,
		`/**
 * @desc d
 */
var MSG_L =
    __jscomp_define_msg__(
        {
          "key":"MSG_L",
          "msg_text":"{$a} has {$b}"
        }, {'a': "{$b}", 'b': 1});
`,
		`/** @desc d */
var MSG_L = '{$b} has ' + 1;
`)
}

// Java: ReplaceMessagesTest.testSimpleMessageReplacementMissing()
func (c *ReplaceMessagesTest) TestSimpleMessageReplacementMissing() {
	c.MultiPhaseTestWarning(
		`/** @desc d */
var MSG_E = 'd*6a0@z>t';
`, //
		`/** @desc d */
var MSG_E =
    __jscomp_define_msg__({"key":"MSG_E", "msg_text":"d*6a0@z\x3et"});
`,
		`/** @desc d */
var MSG_E = 'd*6a0@z>t'
`,
		JsMessageVisitor_MESSAGE_NOT_INITIALIZED_CORRECTLY)
}

// Java: ReplaceMessagesTest.testSimpleMessageReplacementMissingWithNewStyle()
func (c *ReplaceMessagesTest) TestSimpleMessageReplacementMissingWithNewStyle() {
	c.MultiPhaseTest(
		"/** @desc d */\n var MSG_E = goog.getMsg('missing');",
		`/**
 * @desc d
 */
var MSG_E =
    __jscomp_define_msg__(
        {
          "key":"MSG_E",
          "msg_text":"missing",
        });
`,
		"/** @desc d */\n var MSG_E = 'missing'")
}

// Java: ReplaceMessagesTest.testStrictModeAndMessageReplacementAbsentInBundle()
func (c *ReplaceMessagesTest) TestStrictModeAndMessageReplacementAbsentInBundle() {
	c.strictReplacement = true
	c.MultiPhaseTestPostLookupError(
		`/** @desc d */
var MSG_E = goog.getMsg('Hello');
`,
		`/** @desc d */
var MSG_E = __jscomp_define_msg__({"key":"MSG_E", "msg_text":"Hello"});
`,
		ReplaceMessages_BUNDLE_DOES_NOT_HAVE_THE_MESSAGE)
}

// Java: ReplaceMessagesTest.testStrictModeAndMessageReplacementAbsentInNonEmptyBundle()
func (c *ReplaceMessagesTest) TestStrictModeAndMessageReplacementAbsentInNonEmptyBundle() {
	c.RegisterMessage(
		c.GetTestMessageBuilder("MSG_J").
			AppendStringPart("One ").
			AppendJsPlaceholderReference("measly").
			AppendStringPart(" ph").
			Build())

	c.strictReplacement = true
	c.MultiPhaseTestPostLookupError(
		`/** @desc d */
var MSG_E = goog.getMsg('Hello');
`,
		`/** @desc d */
var MSG_E = __jscomp_define_msg__({"key":"MSG_E", "msg_text":"Hello"});
`,
		ReplaceMessages_BUNDLE_DOES_NOT_HAVE_THE_MESSAGE)
}

// Java: ReplaceMessagesTest.testFunctionReplacementMissing()
func (c *ReplaceMessagesTest) TestFunctionReplacementMissing() {
	c.MultiPhaseTestWarning(
		"var MSG_F = function() {return 'asdf'};", //
		`var MSG_F = function() {
  return __jscomp_define_msg__(
      {
        "key":"MSG_F",
        "msg_text":"asdf"
      },
      {});
};
`,
		"var MSG_F = function() {return'asdf'}",
		JsMessageVisitor_MESSAGE_NOT_INITIALIZED_CORRECTLY)
}

// Java: ReplaceMessagesTest.testFunctionWithParamReplacementMissing()
func (c *ReplaceMessagesTest) TestFunctionWithParamReplacementMissing() {
	c.MultiPhaseTestWarning(
		"var MSG_G = function(measly) { return 'asdf' + measly};",
		`var MSG_G = function(measly) {
    return __jscomp_define_msg__(
        {
          "key":"MSG_G",
          "msg_text":"asdf{$measly}"
        },
        {"measly":measly});
    };
`,
		"var MSG_G = function(measly) { return 'asdf' + measly}",
		JsMessageVisitor_MESSAGE_NOT_INITIALIZED_CORRECTLY)
}

// Java: ReplaceMessagesTest.testPlaceholderNameInLowerUnderscoreCase()
func (c *ReplaceMessagesTest) TestPlaceholderNameInLowerUnderscoreCase() {
	c.MultiPhaseTestPreLookupError(
		"var MSG_J = goog.getMsg('${$amt_earned}', {amt_earned: x});", JsMessageVisitor_MESSAGE_TREE_MALFORMED)
}

// Java: ReplaceMessagesTest.testBadPlaceholderReferenceInReplacement()
func (c *ReplaceMessagesTest) TestBadPlaceholderReferenceInReplacement() {
	c.RegisterMessage(c.GetTestMessageBuilder("MSG_K").AppendJsPlaceholderReference("amount").Build())

	c.MultiPhaseTestPostLookupError(
		`/** @desc d */
var MSG_K = goog.getMsg('Hi {$jane}', {jane: x});
`,
		`/** @desc d */
var MSG_K =
    __jscomp_define_msg__(
        { "key":"MSG_K", "msg_text":"Hi {$jane}" },
        {'jane': x});
`,
		JsMessageVisitor_MESSAGE_TREE_MALFORMED)
}

// Java: ReplaceMessagesTest.testTranslatedMessageWithPlaceholdersForGoogGetMsgWithoutAny()
func (c *ReplaceMessagesTest) TestTranslatedMessageWithPlaceholdersForGoogGetMsgWithoutAny() {
	c.RegisterMessage(
		c.GetTestMessageBuilder("MSG_E").
			AppendStringPart("You have purchased ").
			AppendJsPlaceholderReference("amount").
			AppendStringPart(" items.").
			Build())

	c.MultiPhaseTestPostLookupErrorWithDescription(
		`/** @desc d */
var MSG_E = goog.getMsg('no placeholders');
`,
		`/** @desc d */
var MSG_E =
    __jscomp_define_msg__({"key":"MSG_E", "msg_text":"no placeholders"});
`,
		JsMessageVisitor_MESSAGE_TREE_MALFORMED,
		"Message parse tree malformed. The translated message has placeholders, but the definition"+
			" in the JS code does not.")
}

// Java: ReplaceMessagesTest.testLegacyStyleNoPlaceholdersVarSyntaxConcat()
func (c *ReplaceMessagesTest) TestLegacyStyleNoPlaceholdersVarSyntaxConcat() {
	c.RegisterMessage(c.GetTestMessageBuilder("MSG_A").AppendStringPart("Hi\nthere").Build())
	c.MultiPhaseTestWarning(
		"var MSG_A = 'abc' + 'def';", //
		"var MSG_A = __jscomp_define_msg__({\"key\":\"MSG_A\", \"msg_text\":\"abcdef\"});",
		"var MSG_A = 'Hi\\nthere'",
		JsMessageVisitor_MESSAGE_NOT_INITIALIZED_CORRECTLY)
}

// Java: ReplaceMessagesTest.testLegacyStyleNoPlaceholdersVarSyntax()
func (c *ReplaceMessagesTest) TestLegacyStyleNoPlaceholdersVarSyntax() {
	c.RegisterMessage(c.GetTestMessageBuilder("MSG_A").AppendStringPart("Hi\nthere").Build())
	c.MultiPhaseTestWarning(
		"var MSG_A = 'd*6a0@z>t';", //
		"var MSG_A = __jscomp_define_msg__({\"key\":\"MSG_A\", \"msg_text\":\"d*6a0@z\\x3et\"});",
		"var MSG_A='Hi\\nthere'",
		JsMessageVisitor_MESSAGE_NOT_INITIALIZED_CORRECTLY)
}

// Java: ReplaceMessagesTest.testLegacyStyleNoPlaceholdersFunctionSyntax()
func (c *ReplaceMessagesTest) TestLegacyStyleNoPlaceholdersFunctionSyntax() {
	c.RegisterMessage(c.GetTestMessageBuilder("MSG_B").AppendStringPart("Hi\nthere").Build())
	c.MultiPhaseTestWarning(
		"var MSG_B = function() {return 'asdf'};", //
		`var MSG_B = function() {
    return __jscomp_define_msg__(
        {
          "key":"MSG_B",
          "msg_text":"asdf"
        },
        {});
};
`,
		"var MSG_B=function(){return'Hi\\nthere'}",
		JsMessageVisitor_MESSAGE_NOT_INITIALIZED_CORRECTLY)
}

// Java: ReplaceMessagesTest.testLegacyStyleOnePlaceholder()
func (c *ReplaceMessagesTest) TestLegacyStyleOnePlaceholder() {
	c.RegisterMessage(
		c.GetTestMessageBuilder("MSG_C").
			AppendStringPart("One ").
			AppendJsPlaceholderReference("measly").
			AppendStringPart(" ph").
			Build())
	c.MultiPhaseTestWarning(
		"var MSG_C = function(measly) {return 'asdf' + measly};",
		`var MSG_C = function(measly) {
    return __jscomp_define_msg__(
        {
          "key":"MSG_C",
          "msg_text":"asdf{$measly}"
        },
        {"measly":measly});
};
`,
		"var MSG_C=function(measly){ return 'One ' + measly + ' ph'; }",
		JsMessageVisitor_MESSAGE_NOT_INITIALIZED_CORRECTLY)
}

// Java: ReplaceMessagesTest.testLegacyStyleTwoPlaceholders()
func (c *ReplaceMessagesTest) TestLegacyStyleTwoPlaceholders() {
	c.RegisterMessage(
		c.GetTestMessageBuilder("MSG_D").
			AppendJsPlaceholderReference("dick").
			AppendStringPart(" and ").
			AppendJsPlaceholderReference("jane").
			Build())
	c.MultiPhaseTestWarning(
		"var MSG_D = function(jane, dick) {return jane + dick};", //
		`var MSG_D = function(jane, dick) {
    return __jscomp_define_msg__(
        {
          "key":"MSG_D",
          "msg_text":"{$jane}{$dick}"
        },
        {"jane":jane, "dick":dick});
};
`,
		"var MSG_D = function(jane,dick) { return dick + ' and ' + jane; }",
		JsMessageVisitor_MESSAGE_NOT_INITIALIZED_CORRECTLY)
}

// Java: ReplaceMessagesTest.testLegacyStylePlaceholderNameInLowerCamelCase()
func (c *ReplaceMessagesTest) TestLegacyStylePlaceholderNameInLowerCamelCase() {
	c.RegisterMessage(
		c.GetTestMessageBuilder("MSG_E").
			AppendStringPart("Sum: $").
			AppendJsPlaceholderReference("amtEarned").
			Build())
	c.MultiPhaseTestWarning(
		"var MSG_E = function(amtEarned) {return amtEarned + 'x'};",
		`var MSG_E = function(amtEarned) {
    return __jscomp_define_msg__(
        {
          "key":"MSG_E",
          "msg_text":"{$amtEarned}x"
        },
        {"amtEarned":amtEarned});
};
`,
		"var MSG_E=function(amtEarned){return'Sum: $'+amtEarned}",
		JsMessageVisitor_MESSAGE_NOT_INITIALIZED_CORRECTLY)
}

// Java: ReplaceMessagesTest.testInvalidRhs()
func (c *ReplaceMessagesTest) TestInvalidRhs() {
	// If the RHS of a variable named `MSG_*` is not a function call, just report a warning.
	c.MultiPhaseTestWarning(
		"var MSG_A = 'string value';",
		"var MSG_A = 'string value';",
		"var MSG_A = 'string value';",
		JsMessageVisitor_MESSAGE_NOT_INITIALIZED_CORRECTLY)

	c.MultiPhaseTestWarning(
		"var MSG_A = 15 * 12;",
		"var MSG_A = 15 * 12;",
		"var MSG_A = 15 * 12;",
		JsMessageVisitor_MESSAGE_NOT_INITIALIZED_CORRECTLY)
}

// Java: ReplaceMessagesTest.testTranslatedPlaceHolderMissMatch()
func (c *ReplaceMessagesTest) TestTranslatedPlaceHolderMissMatch() {
	c.RegisterMessage(
		c.GetTestMessageBuilder("MSG_A").
			AppendJsPlaceholderReference("a").
			AppendStringPart("!").
			Build())

	c.MultiPhaseTestPreLookupError("var MSG_A = goog.getMsg('{$a}');", JsMessageVisitor_MESSAGE_TREE_MALFORMED)
}

// Java: ReplaceMessagesTest.testTranslatedBadBooleanOptionValue()
func (c *ReplaceMessagesTest) TestTranslatedBadBooleanOptionValue() {
	c.RegisterMessage(
		c.GetTestMessageBuilder("MSG_A").
			AppendJsPlaceholderReference("a").
			AppendStringPart("!").
			Build())

	c.MultiPhaseTestPreLookupError(
		// used an object when a boolean is required
		"var MSG_A = goog.getMsg('{$a}', {'a': 'something'}, { html: {} });",
		JsMessageVisitor_MESSAGE_TREE_MALFORMED)
	c.MultiPhaseTestPreLookupError(
		// used an object when a boolean is required
		"var MSG_A = goog.getMsg('{$a}', {'a': 'something'}, { unescapeHtmlEntities: {} });",
		JsMessageVisitor_MESSAGE_TREE_MALFORMED)
}

// Java: ReplaceMessagesTest.testTranslatedMisspelledExamples()
func (c *ReplaceMessagesTest) TestTranslatedMisspelledExamples() {
	c.RegisterMessage(
		c.GetTestMessageBuilder("MSG_A").
			AppendJsPlaceholderReference("a").
			AppendStringPart("!").
			Build())

	c.MultiPhaseTestPreLookupError(
		// mistakenly used "examples" instead of "example"
		"var MSG_A = goog.getMsg('{$a}', {'a': 'something'}, { examples: { 'a': 'example a' } });",
		JsMessageVisitor_MESSAGE_TREE_MALFORMED)
}

// Java: ReplaceMessagesTest.testTranslatedMisspelledOriginalCode()
func (c *ReplaceMessagesTest) TestTranslatedMisspelledOriginalCode() {
	c.RegisterMessage(
		c.GetTestMessageBuilder("MSG_A").
			AppendJsPlaceholderReference("a").
			AppendStringPart("!").
			Build())

	c.MultiPhaseTestPreLookupError(
		// mistakenly used "original" instead of "original_code"
		"var MSG_A = goog.getMsg('{$a}', {'a': 'something'}, { original: { 'a': 'code' } });",
		JsMessageVisitor_MESSAGE_TREE_MALFORMED)
}

// Java: ReplaceMessagesTest.testTranslatedExampleWithUnknownPlaceholder()
func (c *ReplaceMessagesTest) TestTranslatedExampleWithUnknownPlaceholder() {
	c.RegisterMessage(
		c.GetTestMessageBuilder("MSG_A").
			AppendJsPlaceholderReference("a").
			AppendStringPart("!").
			Build())

	c.MultiPhaseTestPreLookupError(
		"var MSG_A = goog.getMsg('{$a}', {'a': 'something'}, { example: { 'b': 'example a' } });",
		JsMessageVisitor_MESSAGE_TREE_MALFORMED)
}

// Java: ReplaceMessagesTest.testTranslatedExampleWithNonStringPlaceholderValue()
func (c *ReplaceMessagesTest) TestTranslatedExampleWithNonStringPlaceholderValue() {
	c.RegisterMessage(
		c.GetTestMessageBuilder("MSG_A").
			AppendJsPlaceholderReference("a").
			AppendStringPart("!").
			Build())

	c.MultiPhaseTestPreLookupError(
		"var MSG_A = goog.getMsg('{$a}', {'a': 'something'}, { example: { 'a': 1 } });",
		JsMessageVisitor_MESSAGE_TREE_MALFORMED)
}

// Java: ReplaceMessagesTest.testTranslatedExampleWithBadValue()
func (c *ReplaceMessagesTest) TestTranslatedExampleWithBadValue() {
	c.RegisterMessage(
		c.GetTestMessageBuilder("MSG_A").
			AppendJsPlaceholderReference("a").
			AppendStringPart("!").
			Build())

	c.MultiPhaseTestPreLookupError(
		"var MSG_A = goog.getMsg('{$a}', {'a': 'something'}, { example: 'bad value' });",
		JsMessageVisitor_MESSAGE_TREE_MALFORMED)
}

// Java: ReplaceMessagesTest.testTranslatedExampleWithComputedProperty()
func (c *ReplaceMessagesTest) TestTranslatedExampleWithComputedProperty() {
	c.RegisterMessage(
		c.GetTestMessageBuilder("MSG_A").
			AppendJsPlaceholderReference("a").
			AppendStringPart("!").
			Build())

	c.MultiPhaseTestPreLookupError(
		// computed property is not allowed for examples
		"var MSG_A = goog.getMsg('{$a}', {'a': 'something'}, { example: { ['a']: 'wrong' } });",
		JsMessageVisitor_MESSAGE_TREE_MALFORMED)
}

// Java: ReplaceMessagesTest.testBadFallbackSyntax1()
func (c *ReplaceMessagesTest) TestBadFallbackSyntax1() {
	c.MultiPhaseTestPreLookupError(
		`/** @desc d */

var MSG_A = goog.getMsg('asdf');
var x = goog.getMsgWithFallback(MSG_A);
`,
		JsMessageVisitor_BAD_FALLBACK_SYNTAX)
}

// Java: ReplaceMessagesTest.testBadFallbackSyntax2()
func (c *ReplaceMessagesTest) TestBadFallbackSyntax2() {
	c.MultiPhaseTestPreLookupError(
		"var x = goog.getMsgWithFallback('abc', 'bcd');", JsMessageVisitor_BAD_FALLBACK_SYNTAX)
}

// Java: ReplaceMessagesTest.testBadFallbackSyntax3()
func (c *ReplaceMessagesTest) TestBadFallbackSyntax3() {
	c.MultiPhaseTestPreLookupError(
		`/** @desc d */

var MSG_A = goog.getMsg('asdf');var x = goog.getMsgWithFallback(MSG_A, NOT_A_MESSAGE);
`,
		JsMessageVisitor_BAD_FALLBACK_SYNTAX)
}

// Java: ReplaceMessagesTest.testBadFallbackSyntax4()
func (c *ReplaceMessagesTest) TestBadFallbackSyntax4() {
	c.MultiPhaseTestPreLookupError(
		`/** @desc d */

var MSG_A = goog.getMsg('asdf');var x = goog.getMsgWithFallback(NOT_A_MESSAGE, MSG_A);
`,
		JsMessageVisitor_BAD_FALLBACK_SYNTAX)
}

// Java: ReplaceMessagesTest.testBadFallbackSyntax5()
func (c *ReplaceMessagesTest) TestBadFallbackSyntax5() {
	c.MultiPhaseTestPreLookupError(
		`/** @desc d */

var MSG_A = goog.getMsg('asdf');var x = goog.getMsgWithFallback(MSG_A, MSG_DOES_NOT_EXIST);
`,
		JsMessageVisitor_FALLBACK_ARG_ERROR)
}

// Java: ReplaceMessagesTest.testBadFallbackSyntax6()
func (c *ReplaceMessagesTest) TestBadFallbackSyntax6() {
	c.MultiPhaseTestPreLookupError(
		`/** @desc d */

var MSG_A = goog.getMsg('asdf');var x = goog.getMsgWithFallback(MSG_DOES_NOT_EXIST, MSG_A);
`,
		JsMessageVisitor_FALLBACK_ARG_ERROR)
}

// Java: ReplaceMessagesTest.testUseFallback()
func (c *ReplaceMessagesTest) TestUseFallback() {
	c.RegisterMessage(c.GetTestMessageBuilder("MSG_B").AppendStringPart("translated").Build())
	c.MultiPhaseTest(
		`/** @desc d */
var MSG_A = goog.getMsg('msg A');
/** @desc d */
var MSG_B = goog.getMsg('msg B');
var x = goog.getMsgWithFallback(MSG_A, MSG_B);
`,
		`/**
 * @desc d
 */
var MSG_A =
    __jscomp_define_msg__(
        {
          "key":"MSG_A",
          "msg_text":"msg A",
        });
/**
 * @desc d
 */
var MSG_B =
    __jscomp_define_msg__(
        {
          "key":"MSG_B",
          "msg_text":"msg B",
        });
var x = __jscomp_msg_fallback__("MSG_A", MSG_A, "MSG_B", MSG_B);
`,
		`/** @desc d */
var MSG_A = 'msg A';
/** @desc d */
var MSG_B = 'translated';
var x = MSG_B;
`)
}

// Java: ReplaceMessagesTest.testFallbackEmptyBundle()
func (c *ReplaceMessagesTest) TestFallbackEmptyBundle() {
	c.MultiPhaseTest(
		`/** @desc d */
var MSG_A = goog.getMsg('msg A');
/** @desc d */
var MSG_B = goog.getMsg('msg B');
var x = goog.getMsgWithFallback(MSG_A, MSG_B);
`,
		`/**
 * @desc d
 */
var MSG_A =
    __jscomp_define_msg__(
        {
          "key":"MSG_A",
          "msg_text":"msg A",
        });
/**
 * @desc d
 */
var MSG_B =
    __jscomp_define_msg__(
        {
          "key":"MSG_B",
          "msg_text":"msg B",
        });
var x = __jscomp_msg_fallback__("MSG_A", MSG_A, "MSG_B", MSG_B);
`,
		`/** @desc d */
var MSG_A = 'msg A';
/** @desc d */
var MSG_B = 'msg B';
var x = MSG_A;
`)
}

// Java: ReplaceMessagesTest.testNoUseFallback()
func (c *ReplaceMessagesTest) TestNoUseFallback() {
	c.RegisterMessage(c.GetTestMessageBuilder("MSG_A").AppendStringPart("translated").Build())
	c.MultiPhaseTest(
		`/** @desc d */
var MSG_A = goog.getMsg('msg A');
/** @desc d */
var MSG_B = goog.getMsg('msg B');
var x = goog.getMsgWithFallback(MSG_A, MSG_B);
`,
		`/**
 * @desc d
 */
var MSG_A =
    __jscomp_define_msg__(
        {
          "key":"MSG_A",
          "msg_text":"msg A",
        });
/**
 * @desc d
 */
var MSG_B =
    __jscomp_define_msg__(
        {
          "key":"MSG_B",
          "msg_text":"msg B",
        });
var x = __jscomp_msg_fallback__("MSG_A", MSG_A, "MSG_B", MSG_B);
`,
		`/** @desc d */
var MSG_A = 'translated';
/** @desc d */
var MSG_B = 'msg B';
var x = MSG_A;
`)
}

// Java: ReplaceMessagesTest.testNoUseFallback2()
func (c *ReplaceMessagesTest) TestNoUseFallback2() {
	c.RegisterMessage(c.GetTestMessageBuilder("MSG_C").AppendStringPart("translated").Build())
	c.MultiPhaseTest(
		`/** @desc d */
var MSG_A = goog.getMsg('msg A');
/** @desc d */
var MSG_B = goog.getMsg('msg B');
var x = goog.getMsgWithFallback(MSG_A, MSG_B);
`,
		`/**
 * @desc d
 */
var MSG_A =
    __jscomp_define_msg__(
        {
          "key":"MSG_A",
          "msg_text":"msg A",
        });
/**
 * @desc d
 */
var MSG_B =
    __jscomp_define_msg__(
        {
          "key":"MSG_B",
          "msg_text":"msg B",
        });
var x = __jscomp_msg_fallback__("MSG_A", MSG_A, "MSG_B", MSG_B);
`,
		`/** @desc d */
var MSG_A = 'msg A';
/** @desc d */
var MSG_B = 'msg B';
var x = MSG_A;
`)
}

// Java: ReplaceMessagesTest.testTemplateLiteralSimple()
func (c *ReplaceMessagesTest) TestTemplateLiteralSimple() {
	c.RegisterMessage(c.GetTestMessageBuilder("MSG_A").AppendStringPart("Hi\nthere").Build())

	c.MultiPhaseTest(
		"/** @desc d */\n var MSG_A = goog.getMsg(`asdf`);",
		`/**
 * @desc d
 */
var MSG_A =
    __jscomp_define_msg__(
        {
          "key":"MSG_A",
          "msg_text":"asdf",
        });
`,
		"/** @desc d */\n var MSG_A='Hi\\nthere'")
}

// Java: ReplaceMessagesTest.testTemplateLiteralNameReplacement()
func (c *ReplaceMessagesTest) TestTemplateLiteralNameReplacement() {
	c.RegisterMessage(
		c.GetTestMessageBuilder("MSG_B").
			AppendStringPart("One ").
			AppendJsPlaceholderReference("measly").
			AppendStringPart(" ph").
			Build())

	c.MultiPhaseTest(
		"/** @desc d */\n var MSG_B=goog.getMsg(`asdf {$measly}`, {measly: x});",
		`/**
 * @desc d
 */
var MSG_B =
    __jscomp_define_msg__(
        {
          "key":"MSG_B",
          "msg_text":"asdf {$measly}",
        },
        {'measly': x});
`,
		`/** @desc d */
var MSG_B = 'One ' + x + ' ph';
`)
}

// Java: ReplaceMessagesTest.testTemplateLiteralSubstitutions()
func (c *ReplaceMessagesTest) TestTemplateLiteralSubstitutions() {
	// Only allow template literals that are constant strings
	c.RegisterMessage(c.GetTestMessageBuilder("MSG_C").AppendStringPart("Hi\nthere").Build())

	c.MultiPhaseTestPreLookupError(
		"/** @desc d */\n var MSG_C = goog.getMsg(`asdf ${42}`);",
		JsMessageVisitor_MESSAGE_TREE_MALFORMED)
}

// Java: ReplaceMessagesTest.testReplaceUnescapeHtmlEntitiesMessage()
func (c *ReplaceMessagesTest) TestReplaceUnescapeHtmlEntitiesMessage() {
	c.MultiPhaseTest(
		"/** @desc d */\n var MSG_A = goog.getMsg('A', {}, {unescapeHtmlEntities: true});",
		`/**
 * @desc d
 */
var MSG_A =
    __jscomp_define_msg__(
        {
          "key":"MSG_A",
          "msg_text":"A",
          "unescapeHtmlEntities":""
        },
        {});
`,
		"/** @desc d */\n var MSG_A = 'A';")
	c.MultiPhaseTest(
		`/** @desc d */

var MSG_A = goog.getMsg('User&apos;s &lt; email &amp; address &gt; are &quot;correct&quot;', {}, {unescapeHtmlEntities: true});
`,
		`/**
 * @desc d
 */
var MSG_A =
    __jscomp_define_msg__(
        {
          "key":"MSG_A",
          "msg_text":"User\x26apos;s \x26lt; email \x26amp; address \x26gt; are \x26quot;correct\x26quot;",
          "unescapeHtmlEntities":""
        },
        {});

`,
		"/** @desc d */\n var MSG_A = 'User\\'s < email & address > are \"correct\"';")
	c.MultiPhaseTest(
		`/** @desc d */

var MSG_A = goog.getMsg('&lt; {$startSpan}email &amp; address{$endSpan} &gt;', {'startSpan': '<span title="&lt;info&gt;">', 'endSpan': '</span>'}, {unescapeHtmlEntities: true});
`,
		`/**
 * @desc d
 */
var MSG_A =
    __jscomp_define_msg__(
        {
          "key":"MSG_A",
          "msg_text":"\x26lt; {$startSpan}email \x26amp; address{$endSpan} \x26gt;",
          "unescapeHtmlEntities":""
        },
        {
          "startSpan":'\x3cspan title\x3d"\x26lt;info\x26gt;"\x3e',
          "endSpan":"\x3c/span\x3e"
        });
`,
		`/** @desc d */
var MSG_A =
    '< <span title="&lt;info&gt;">email & address</span> >';
`)
	c.MultiPhaseTest(
		`/** @desc d */

var MSG_A = goog.getMsg('&amp;lt;double &amp;amp; escaping&amp;gt;', {}, {unescapeHtmlEntities: true});
`,
		`/**
* @desc d
*/
var MSG_A =
    __jscomp_define_msg__(
        {
          "key":"MSG_A",
          "msg_text":"\x26amp;lt;double \x26amp;amp; escaping\x26amp;gt;",
          "unescapeHtmlEntities":""
        },
        {});
`,
		"/** @desc d */\n var MSG_A = '&lt;double &amp; escaping&gt;';")
}

// Java: ReplaceMessagesTest.testReplaceUnescapeHtmlEntitiesMessageWithReplacement()
func (c *ReplaceMessagesTest) TestReplaceUnescapeHtmlEntitiesMessageWithReplacement() {
	c.RegisterMessage(
		c.GetTestMessageBuilder("MSG_A").
			AppendStringPart("User").
			AppendStringPart("&").
			AppendStringPart("apos;s &").
			AppendStringPart("lt;").
			AppendStringPart(" email &a").
			AppendStringPart("mp; address &gt").
			AppendStringPart("; are &quot;correct").
			AppendStringPart("&quot;").
			Build())
	c.MultiPhaseTest(
		"/** @desc d */\n var MSG_A = goog.getMsg('A', {}, {unescapeHtmlEntities: true});",
		`/**
 * @desc d
 */
var MSG_A =
    __jscomp_define_msg__(
        {
          "key":"MSG_A",
          "msg_text":"A",
          "unescapeHtmlEntities":""
        },
        {});
`,
		`/** @desc d */
var MSG_A = 'User\'s < email & address > are "correct"';
`)

	c.RegisterMessage(
		c.GetTestMessageBuilder("MSG_B").
			AppendStringPart("User").
			AppendStringPart("&apos;").
			AppendStringPart("s ").
			AppendStringPart("&lt;").
			AppendStringPart(" email ").
			AppendStringPart("&amp;").
			AppendStringPart(" address ").
			AppendStringPart("&gt;").
			AppendStringPart(" are ").
			AppendStringPart("&quot;").
			AppendStringPart("correct").
			AppendStringPart("&quot;").
			Build())
	c.MultiPhaseTest(
		"/** @desc d */\n var MSG_B = goog.getMsg('B', {}, {unescapeHtmlEntities: true});",
		`/**
 * @desc d
 */
var MSG_B =
    __jscomp_define_msg__(
        {
          "key":"MSG_B",
          "msg_text":"B",
          "unescapeHtmlEntities":""
        },
        {});
`,
		"/** @desc d */\n var MSG_B = 'User\\'s < email & address > are \"correct\"';")

	c.RegisterMessage(
		c.GetTestMessageBuilder("MSG_C").
			AppendJsPlaceholderReference("br").
			AppendStringPart("&").
			AppendStringPart("amp;").
			AppendJsPlaceholderReference("x").
			AppendJsPlaceholderReference("y").
			AppendStringPart("&ap").
			AppendJsPlaceholderReference("z").
			AppendStringPart("os;").
			Build())
	c.MultiPhaseTest(
		`/** @desc d */

var MSG_C = goog.getMsg('{$br}{$x}{$y}{$z}', {'br': '<br>', 'x': 'X', 'y': 'Y', 'z': 'Z'}, {unescapeHtmlEntities: true});
`,
		`/**
 * @desc d
 */
var MSG_C =
    __jscomp_define_msg__(
        {
          "key":"MSG_C",
          "msg_text":"{$br}{$x}{$y}{$z}",
          "unescapeHtmlEntities":""
        },
        {
          "br":"\x3cbr\x3e",
          "x":"X",
          "y":"Y",
          "z":"Z"
        });
`,
		`/** @desc d */
var MSG_C = '<br>&XY&apZos;';
`)
}

// Java: ReplaceMessagesTest.testReplaceHtmlMessageWithPlaceholder()
func (c *ReplaceMessagesTest) TestReplaceHtmlMessageWithPlaceholder() {
	c.RegisterMessage(
		c.GetTestMessageBuilder("MSG_A").
			AppendStringPart("Hello <"). // html option changes `<` to `&lt;.
			AppendJsPlaceholderReference("br").
			AppendStringPart("&gt;").
			Build())

	c.MultiPhaseTest(
		"/** @desc d */\n var MSG_A = goog.getMsg('{$br}', {'br': '<br>'}, {html: true});",
		`/**
 * @desc d
 */
var MSG_A =
    __jscomp_define_msg__(
        {
          "key":"MSG_A",
          "msg_text":"{$br}",
          "escapeLessThan":"",
        },
        {"br":"\x3cbr\x3e"});
`,
		`/** @desc d */
var MSG_A = 'Hello &lt;<br>&gt;';
`)

	// Confirm that the default behavior is to leave `<` unchanged
	c.MultiPhaseTest(
		"/** @desc d */\n var MSG_A = goog.getMsg('{$br}', {'br': '<br>'});",
		`/**
 * @desc d
 */
var MSG_A =
    __jscomp_define_msg__(
        {
          "key":"MSG_A",
          "msg_text":"{$br}",
        },
        {"br":"\x3cbr\x3e"});
`,
		`/** @desc d */
var MSG_A = 'Hello <<br>&gt;';
`)
}

// Java: ReplaceMessagesTest.testReplaceGenderedMessagesWithoutPlaceholders()
func (c *ReplaceMessagesTest) TestReplaceGenderedMessagesWithoutPlaceholders() {
	c.RegisterMessage(
		c.GetTestMessageBuilder("MSG_E").
			AddGenderedMessageKey(JsMessage_GrammaticalGenderCase_MASCULINE).
			AppendStringPartForGender(JsMessage_GrammaticalGenderCase_MASCULINE, "Bienvenido!").
			AddGenderedMessageKey(JsMessage_GrammaticalGenderCase_FEMININE).
			AppendStringPartForGender(JsMessage_GrammaticalGenderCase_FEMININE, "Bienvenida!").
			AddGenderedMessageKey(JsMessage_GrammaticalGenderCase_NEUTER).
			AppendStringPartForGender(JsMessage_GrammaticalGenderCase_NEUTER, "Les damos la bienvenida!").
			AddGenderedMessageKey(JsMessage_GrammaticalGenderCase_OTHER).
			AppendStringPartForGender(
				JsMessage_GrammaticalGenderCase_OTHER, "Les damos la bienvenida! - OTHER").
			Build())

	c.MultiPhaseTest(
		"/** @desc d */\nvar MSG_E = goog.getMsg('Welcome!');",
		`/**
 * @desc d
 */
var MSG_E = __jscomp_define_msg__({"key":"MSG_E", "msg_text":"Welcome!"});
`,
		`/** @desc d */
var MSG_E = goog.msgKind.MASCULINE ? 'Bienvenido!' :
            goog.msgKind.FEMININE  ? 'Bienvenida!' :
            goog.msgKind.NEUTER  ? 'Les damos la bienvenida!' :
            'Les damos la bienvenida! - OTHER';
`)
}

// Java: ReplaceMessagesTest.testReplaceGenderedMessagesWithoutPlaceholdersNoNeuter()
func (c *ReplaceMessagesTest) TestReplaceGenderedMessagesWithoutPlaceholdersNoNeuter() {
	c.RegisterMessage(
		c.GetTestMessageBuilder("MSG_E").
			AddGenderedMessageKey(JsMessage_GrammaticalGenderCase_MASCULINE).
			AppendStringPartForGender(JsMessage_GrammaticalGenderCase_MASCULINE, "Bienvenido!").
			AddGenderedMessageKey(JsMessage_GrammaticalGenderCase_FEMININE).
			AppendStringPartForGender(JsMessage_GrammaticalGenderCase_FEMININE, "Bienvenida!").
			AddGenderedMessageKey(JsMessage_GrammaticalGenderCase_OTHER).
			AppendStringPartForGender(JsMessage_GrammaticalGenderCase_OTHER, "Les damos la bienvenida!").
			Build())

	c.MultiPhaseTest(
		"/** @desc d */\nvar MSG_E = goog.getMsg('Welcome!');",
		`/**
 * @desc d
 */
var MSG_E = __jscomp_define_msg__({"key":"MSG_E", "msg_text":"Welcome!"});
`,
		`/** @desc d */
var MSG_E = goog.msgKind.MASCULINE ? 'Bienvenido!' :
            goog.msgKind.FEMININE  ? 'Bienvenida!' :
            'Les damos la bienvenida!';
`)
}

// Java: ReplaceMessagesTest.testReplaceGenderedMessagesWithoutPlaceholdersDifferentOrder()
func (c *ReplaceMessagesTest) TestReplaceGenderedMessagesWithoutPlaceholdersDifferentOrder() {
	c.RegisterMessage(
		c.GetTestMessageBuilder("MSG_E").
			AddGenderedMessageKey(JsMessage_GrammaticalGenderCase_OTHER).
			AppendStringPartForGender(
				JsMessage_GrammaticalGenderCase_OTHER, "Les damos la bienvenida! - OTHER").
			AddGenderedMessageKey(JsMessage_GrammaticalGenderCase_NEUTER).
			AppendStringPartForGender(JsMessage_GrammaticalGenderCase_NEUTER, "Les damos la bienvenida!").
			AddGenderedMessageKey(JsMessage_GrammaticalGenderCase_FEMININE).
			AppendStringPartForGender(JsMessage_GrammaticalGenderCase_FEMININE, "Bienvenida!").
			AddGenderedMessageKey(JsMessage_GrammaticalGenderCase_MASCULINE).
			AppendStringPartForGender(JsMessage_GrammaticalGenderCase_MASCULINE, "Bienvenido!").
			Build())

	c.MultiPhaseTest(
		"/** @desc d */\nvar MSG_E = goog.getMsg('Welcome!');",
		`/**
 * @desc d
 */
var MSG_E = __jscomp_define_msg__({"key":"MSG_E", "msg_text":"Welcome!"});
`,
		`/** @desc d */
var MSG_E = goog.msgKind.NEUTER  ? 'Les damos la bienvenida!' :
            goog.msgKind.FEMININE  ? 'Bienvenida!' :
            goog.msgKind.MASCULINE ? 'Bienvenido!' :
            'Les damos la bienvenida! - OTHER';
`)
}

// Java: ReplaceMessagesTest.testReplaceGenderedMessagesWithoutPlaceholdersDifferentOrder2()
func (c *ReplaceMessagesTest) TestReplaceGenderedMessagesWithoutPlaceholdersDifferentOrder2() {
	c.RegisterMessage(
		c.GetTestMessageBuilder("MSG_E").
			AddGenderedMessageKey(JsMessage_GrammaticalGenderCase_NEUTER).
			AppendStringPartForGender(JsMessage_GrammaticalGenderCase_NEUTER, "Les damos la bienvenida!").
			AddGenderedMessageKey(JsMessage_GrammaticalGenderCase_OTHER).
			AppendStringPartForGender(
				JsMessage_GrammaticalGenderCase_OTHER, "Les damos la bienvenida! - OTHER").
			AddGenderedMessageKey(JsMessage_GrammaticalGenderCase_FEMININE).
			AppendStringPartForGender(JsMessage_GrammaticalGenderCase_FEMININE, "Bienvenida!").
			AddGenderedMessageKey(JsMessage_GrammaticalGenderCase_MASCULINE).
			AppendStringPartForGender(JsMessage_GrammaticalGenderCase_MASCULINE, "Bienvenido!").
			Build())

	c.MultiPhaseTest(
		"/** @desc d */\nvar MSG_E = goog.getMsg('Welcome!');",
		`/**
 * @desc d
 */
var MSG_E = __jscomp_define_msg__({"key":"MSG_E", "msg_text":"Welcome!"});
`,
		`/** @desc d */
var MSG_E = goog.msgKind.NEUTER  ? 'Les damos la bienvenida!' :
            goog.msgKind.FEMININE ? 'Bienvenida!' :
            goog.msgKind.MASCULINE  ? 'Bienvenido!' :
            'Les damos la bienvenida! - OTHER';
`)
}

// Java: ReplaceMessagesTest.testReplaceGenderedMessagesWithoutPlaceholdersDifferentOrder3()
func (c *ReplaceMessagesTest) TestReplaceGenderedMessagesWithoutPlaceholdersDifferentOrder3() {
	c.RegisterMessage(
		c.GetTestMessageBuilder("MSG_E").
			AddGenderedMessageKey(JsMessage_GrammaticalGenderCase_NEUTER).
			AppendStringPartForGender(JsMessage_GrammaticalGenderCase_NEUTER, "Les damos la bienvenida!").
			AddGenderedMessageKey(JsMessage_GrammaticalGenderCase_FEMININE).
			AppendStringPartForGender(JsMessage_GrammaticalGenderCase_FEMININE, "Bienvenida!").
			AddGenderedMessageKey(JsMessage_GrammaticalGenderCase_OTHER).
			AppendStringPartForGender(
				JsMessage_GrammaticalGenderCase_OTHER, "Les damos la bienvenida! - OTHER").
			AddGenderedMessageKey(JsMessage_GrammaticalGenderCase_MASCULINE).
			AppendStringPartForGender(JsMessage_GrammaticalGenderCase_MASCULINE, "Bienvenido!").
			Build())

	c.MultiPhaseTest(
		"/** @desc d */\nvar MSG_E = goog.getMsg('Welcome!');",
		`/**
 * @desc d
 */
var MSG_E = __jscomp_define_msg__({"key":"MSG_E", "msg_text":"Welcome!"});
`,
		`/** @desc d */
var MSG_E = goog.msgKind.NEUTER  ? 'Les damos la bienvenida!' :
            goog.msgKind.FEMININE  ? 'Bienvenida!' :
            goog.msgKind.MASCULINE  ? 'Bienvenido!' :
            'Les damos la bienvenida! - OTHER';
`)
}

// Java: ReplaceMessagesTest.testReplaceGenderedMessageWithPlaceholders()
func (c *ReplaceMessagesTest) TestReplaceGenderedMessageWithPlaceholders() {
	c.RegisterMessage(
		c.GetTestMessageBuilder("MSG_E").
			AddGenderedMessageKey(JsMessage_GrammaticalGenderCase_MASCULINE).
			AppendStringPartForGender(JsMessage_GrammaticalGenderCase_MASCULINE, "Bienvenido ").
			AppendJsPlaceholderReferenceForGender(JsMessage_GrammaticalGenderCase_MASCULINE, "name").
			AppendStringPartForGender(JsMessage_GrammaticalGenderCase_MASCULINE, ", a tu nuevo ").
			AppendJsPlaceholderReferenceForGender(JsMessage_GrammaticalGenderCase_MASCULINE, "device").
			AddGenderedMessageKey(JsMessage_GrammaticalGenderCase_FEMININE).
			AppendStringPartForGender(JsMessage_GrammaticalGenderCase_FEMININE, "Bienvenida ").
			AppendJsPlaceholderReferenceForGender(JsMessage_GrammaticalGenderCase_FEMININE, "name").
			AppendStringPartForGender(JsMessage_GrammaticalGenderCase_FEMININE, ", a tu nuevo ").
			AppendJsPlaceholderReferenceForGender(JsMessage_GrammaticalGenderCase_FEMININE, "device").
			AddGenderedMessageKey(JsMessage_GrammaticalGenderCase_NEUTER).
			AppendStringPartForGender(JsMessage_GrammaticalGenderCase_NEUTER, "Les damos la bienvenida ").
			AppendJsPlaceholderReferenceForGender(JsMessage_GrammaticalGenderCase_NEUTER, "name").
			AppendStringPartForGender(JsMessage_GrammaticalGenderCase_NEUTER, ", a tu nuevo ").
			AppendJsPlaceholderReferenceForGender(JsMessage_GrammaticalGenderCase_NEUTER, "device").
			AddGenderedMessageKey(JsMessage_GrammaticalGenderCase_OTHER).
			AppendStringPartForGender(
				JsMessage_GrammaticalGenderCase_OTHER, "Les damos la bienvenida! - OTHER ").
			AppendJsPlaceholderReferenceForGender(JsMessage_GrammaticalGenderCase_OTHER, "name").
			AppendStringPartForGender(JsMessage_GrammaticalGenderCase_OTHER, ", a tu nuevo ").
			AppendJsPlaceholderReferenceForGender(JsMessage_GrammaticalGenderCase_OTHER, "device").
			Build())

	c.MultiPhaseTest(
		`/** @desc d */
var MSG_E = goog.getMsg(
    'Welcome, {$name}, to your new {$device}.',
    {name: user.getName(), device: user.getDevice()});
`,
		`/**
 * @desc d
 */
 var MSG_E = __jscomp_define_msg__({
  "key": "MSG_E",
  "msg_text": "Welcome, {$name}, to your new {$device}."
}, {
  "name": user.getName(),
  "device": user.getDevice()
});
`,
		`/** @desc d */
var MSG_E = function(namem1146332801$0, devicem1146332801$0) {
  return goog.msgKind.MASCULINE ? 'Bienvenido ' + namem1146332801$0 + ', a tu nuevo ' + devicem1146332801$0 :
         goog.msgKind.FEMININE  ? 'Bienvenida ' + namem1146332801$0 + ', a tu nuevo ' + devicem1146332801$0 :
         goog.msgKind.NEUTER  ? 'Les damos la bienvenida ' + namem1146332801$0 + ', a tu nuevo ' + devicem1146332801$0 :
         'Les damos la bienvenida! - OTHER ' + namem1146332801$0 + ', a tu nuevo ' + devicem1146332801$0;
}(user.getName(), user.getDevice());
`)
}

// Java: ReplaceMessagesTest.testReplaceIcuTemplateGenderedMessageWithBundleAndJsPlaceholders()
func (c *ReplaceMessagesTest) TestReplaceIcuTemplateGenderedMessageWithBundleAndJsPlaceholders() {
	c.RegisterMessage(
		c.GetTestMessageBuilder("MSG_E").
			AddGenderedMessageKey(JsMessage_GrammaticalGenderCase_MASCULINE).
			AppendStringPartForGender(JsMessage_GrammaticalGenderCase_MASCULINE, "Bienvenido ").
			AppendCanonicalPlaceholderReferenceForGender(JsMessage_GrammaticalGenderCase_MASCULINE, "NAME").
			AppendStringPartForGender(JsMessage_GrammaticalGenderCase_MASCULINE, ", a tu nuevo ").
			AppendCanonicalPlaceholderReferenceForGender(JsMessage_GrammaticalGenderCase_MASCULINE, "PHONE").
			AddGenderedMessageKey(JsMessage_GrammaticalGenderCase_FEMININE).
			AppendStringPartForGender(JsMessage_GrammaticalGenderCase_FEMININE, "Bienvenida ").
			AppendCanonicalPlaceholderReferenceForGender(JsMessage_GrammaticalGenderCase_FEMININE, "NAME").
			AppendStringPartForGender(JsMessage_GrammaticalGenderCase_FEMININE, ", a tu nuevo ").
			AppendCanonicalPlaceholderReferenceForGender(JsMessage_GrammaticalGenderCase_FEMININE, "PHONE").
			AddGenderedMessageKey(JsMessage_GrammaticalGenderCase_NEUTER).
			AppendStringPartForGender(JsMessage_GrammaticalGenderCase_NEUTER, "Les damos la bienvenida ").
			AppendCanonicalPlaceholderReferenceForGender(JsMessage_GrammaticalGenderCase_NEUTER, "NAME").
			AppendStringPartForGender(JsMessage_GrammaticalGenderCase_NEUTER, ", a tu nuevo ").
			AppendCanonicalPlaceholderReferenceForGender(JsMessage_GrammaticalGenderCase_NEUTER, "PHONE").
			AddGenderedMessageKey(JsMessage_GrammaticalGenderCase_OTHER).
			AppendStringPartForGender(
				JsMessage_GrammaticalGenderCase_OTHER, "Les damos la bienvenida! - OTHER ").
			AppendCanonicalPlaceholderReferenceForGender(JsMessage_GrammaticalGenderCase_OTHER, "NAME").
			AppendStringPartForGender(JsMessage_GrammaticalGenderCase_OTHER, ", a tu nuevo ").
			AppendCanonicalPlaceholderReferenceForGender(JsMessage_GrammaticalGenderCase_OTHER, "PHONE").
			Build())

	c.MultiPhaseTest(
		`const {declareIcuTemplate} = goog.require('goog.i18n.messages');

var MSG_E =
    declareIcuTemplate(
        'Welcome {NAME}, to your new {PHONE}',
        {
          description: 'Welcome message',
          example: {
            'NAME': 'John Doe',
            'PHONE': 'Pixel'
           }
        });
`,
		`const {declareIcuTemplate} = goog.require('goog.i18n.messages');

var MSG_E =
    __jscomp_define_msg__(
        {
          "key":    "MSG_E",
          "icu_placeholder_names": ["NAME", "PHONE"],
          "msg_text": "Welcome {NAME}, to your new {PHONE}",
          "isIcuTemplate": ""
        });
`,
		`const {declareIcuTemplate} = goog.require('goog.i18n.messages');

var MSG_E = goog.msgKind.MASCULINE ? 'Bienvenido {NAME}, a tu nuevo {PHONE}' :
            goog.msgKind.FEMININE  ? 'Bienvenida {NAME}, a tu nuevo {PHONE}' :
            goog.msgKind.NEUTER  ? 'Les damos la bienvenida {NAME}, a tu nuevo {PHONE}' :
            'Les damos la bienvenida! - OTHER {NAME}, a tu nuevo {PHONE}';
`)
}

// Java: ReplaceMessagesTest.testReplaceGenderedMessageWithIcuPlaceholderAtStart()
func (c *ReplaceMessagesTest) TestReplaceGenderedMessageWithIcuPlaceholderAtStart() {
	c.RegisterMessage(
		c.GetTestMessageBuilder("MSG_StartWithPlaceholder").
			AddGenderedMessageKey(JsMessage_GrammaticalGenderCase_MASCULINE).
			AppendCanonicalPlaceholderReferenceForGender(JsMessage_GrammaticalGenderCase_MASCULINE, "PH").
			AppendStringPartForGender(JsMessage_GrammaticalGenderCase_MASCULINE, " suffix").
			AddGenderedMessageKey(JsMessage_GrammaticalGenderCase_OTHER).
			AppendCanonicalPlaceholderReferenceForGender(JsMessage_GrammaticalGenderCase_OTHER, "PH").
			AppendStringPartForGender(JsMessage_GrammaticalGenderCase_OTHER, " suffix").
			Build())

	// We expect this to fail with NPE before the fix
	c.MultiPhaseTest(
		`/** @desc d */
var MSG_StartWithPlaceholder = goog.getMsg('{PH} suffix');
`,
		`/**
 * @desc d
 */
var MSG_StartWithPlaceholder =
    __jscomp_define_msg__(
        {
          "key": "MSG_StartWithPlaceholder",
          "msg_text": "{PH} suffix"
        });
`,
		`/** @desc d */
var MSG_StartWithPlaceholder = goog.msgKind.MASCULINE ? '{PH} suffix' : '{PH} suffix';
`)
}

// TODO(b/538170233): ReplaceMessages preserving side effects in placeholders
//
// Java: ReplaceMessagesTest.testRepeatedPlaceholderWithSideEffects()
func (c *ReplaceMessagesTest) TestRepeatedPlaceholderWithSideEffects() {
	c.RegisterMessage(
		c.GetTestMessageBuilder("MSG_REPEATED_SE").
			AppendJsPlaceholderReference("t").
			AppendStringPart(" - ").
			AppendJsPlaceholderReference("t").
			Build())

	c.MultiPhaseTest(
		`/** @desc d */
var MSG_REPEATED_SE = goog.getMsg('{$t} - {$t}', {t: nextToken()});
`,
		`/**
 * @desc d
 */
var MSG_REPEATED_SE =
    __jscomp_define_msg__(
        {
          "key": "MSG_REPEATED_SE",
          "msg_text": "{$t} - {$t}"
        },
        {"t": nextToken()});
`,
		`/** @desc d */
var MSG_REPEATED_SE = nextToken() + " - " + nextToken();
`)
}

// TODO(b/538170233): ReplaceMessages preserving side effects in placeholders
//
// Java: ReplaceMessagesTest.testDroppedPlaceholderWithSideEffects()
func (c *ReplaceMessagesTest) TestDroppedPlaceholderWithSideEffects() {
	c.RegisterMessage(
		c.GetTestMessageBuilder("MSG_DROPPED_SE").AppendStringPart("Bonjour").Build())

	c.MultiPhaseTest(
		`/** @desc d */
var MSG_DROPPED_SE = goog.getMsg('Hi {$x}', {x: nextToken()});
`,
		`/**
 * @desc d
 */
var MSG_DROPPED_SE =
    __jscomp_define_msg__(
        {
          "key": "MSG_DROPPED_SE",
          "msg_text": "Hi {$x}"
        },
        {"x": nextToken()});
`,
		`/** @desc d */
var MSG_DROPPED_SE = "Bonjour";
`)
}

// Java: ReplaceMessagesTest.testPlaceholderWithoutSideEffectsNoIIFE()
func (c *ReplaceMessagesTest) TestPlaceholderWithoutSideEffectsNoIIFE() {
	c.RegisterMessage(
		c.GetTestMessageBuilder("MSG_NO_SE").
			AppendStringPart("Hi ").
			AppendJsPlaceholderReference("x").
			Build())

	c.MultiPhaseTest(
		`/** @desc d */
var MSG_NO_SE = goog.getMsg('Hi {$x}', {x: 'User'});
`,
		`/**
 * @desc d
 */
var MSG_NO_SE =
    __jscomp_define_msg__(
        {
          "key": "MSG_NO_SE",
          "msg_text": "Hi {$x}"
        },
        {"x": "User"});
`,
		`/** @desc d */
var MSG_NO_SE = "Hi User";
`)
}

// Java: ReplaceMessagesTest.testRepeatedPlaceholderWithoutSideEffectsNoIIFE()
func (c *ReplaceMessagesTest) TestRepeatedPlaceholderWithoutSideEffectsNoIIFE() {
	c.RegisterMessage(
		c.GetTestMessageBuilder("MSG_REPEATED_NO_SE").
			AppendJsPlaceholderReference("x").
			AppendStringPart(" - ").
			AppendJsPlaceholderReference("x").
			Build())

	c.MultiPhaseTest(
		`/** @desc d */
var MSG_REPEATED_NO_SE = goog.getMsg('{$x} - {$x}', {x: 'User'});
`,
		`/**
 * @desc d
 */
var MSG_REPEATED_NO_SE =
    __jscomp_define_msg__(
        {
          "key": "MSG_REPEATED_NO_SE",
          "msg_text": "{$x} - {$x}"
        },
        {"x": "User"});
`,
		`/** @desc d */
var MSG_REPEATED_NO_SE = "User - User";
`)
}

// Java: ReplaceMessagesTest.registerMessage(JsMessage)
func (c *ReplaceMessagesTest) RegisterMessage(message *JsMessage) {
	c.messages.Put(message.GetId(), message)
}

// Java: ReplaceMessagesTest.testLeakingPlaceholderRegression()
func (c *ReplaceMessagesTest) TestLeakingPlaceholderRegression() {
	// Message A: Gendered, has placeholder "name".
	builderA := NewJsMessage_Builder().SetKey("MSG_A").SetId("MSG_A")
	builderA.AddGenderedMessageKey(JsMessage_GrammaticalGenderCase_MASCULINE)
	builderA.AppendStringPartForGender(JsMessage_GrammaticalGenderCase_MASCULINE, "A male ")
	builderA.AppendJsPlaceholderReferenceForGender(JsMessage_GrammaticalGenderCase_MASCULINE, "name")
	builderA.AddGenderedMessageKey(JsMessage_GrammaticalGenderCase_OTHER)
	builderA.AppendStringPartForGender(JsMessage_GrammaticalGenderCase_OTHER, "A other ")
	builderA.AppendJsPlaceholderReferenceForGender(JsMessage_GrammaticalGenderCase_OTHER, "name")
	c.RegisterMessage(builderA.Build())

	// Message B: Gendered, has placeholder "name" in the translation parts (simulating ICU mixed),
	// but will NOT have it in the goog.getMsg placeholders map.
	builderB := NewJsMessage_Builder().SetKey("MSG_B").SetId("MSG_B")
	builderB.AddGenderedMessageKey(JsMessage_GrammaticalGenderCase_MASCULINE)
	builderB.AppendStringPartForGender(JsMessage_GrammaticalGenderCase_MASCULINE, "B male {name}")
	// This part is a placeholder, but since it's not in the goog.getMsg args, it should be treated
	// as literal
	builderB.AppendJsPlaceholderReferenceForGender(JsMessage_GrammaticalGenderCase_MASCULINE, "name")

	builderB.AddGenderedMessageKey(JsMessage_GrammaticalGenderCase_OTHER)
	builderB.AppendStringPartForGender(JsMessage_GrammaticalGenderCase_OTHER, "B other {name}")
	builderB.AppendJsPlaceholderReferenceForGender(JsMessage_GrammaticalGenderCase_OTHER, "name")

	c.RegisterMessage(builderB.Build())

	// Note: We use test() directly instead of multiPhaseTest because we want to test interaction
	// between two messages in the same pass run, and we are specifically targeting the bug in
	// ReplacementCompletionPass logic where state leaked.
	// However, multiPhaseTest does separate compilations for each phase.
	// If we use FULL_REPLACE mode (default), it runs getFullReplacementPass() which does it all in
	// one go (for legacy reasons).
	// The bug was in ReplaceMessages class structure (field vs local), which is shared if the same
	// instance is used.
	// getProcessor() creates a new instance.
	// But getFullReplacementPass() returns an inner class instance.
	// So if we run FULL_REPLACE, we exercise the logic.

	c.testMode = ReplaceMessagesTest_TestMode_FULL_REPLACE
	c.Test(
		`/** @desc d */
var MSG_A = goog.getMsg('A {$name}', {name: 'User'});
/** @desc d */
var MSG_B = goog.getMsg('B {name}');
`,
		`/** @desc d */
var MSG_A = function(namem1146332801$0) {
  return goog.msgKind.MASCULINE ? "A male " + namem1146332801$0 : "A other " + namem1146332801$0;
}("User");
/** @desc d */
var MSG_B = goog.msgKind.MASCULINE ? 'B male {name}{NAME}' : 'B other {name}{NAME}';
`)
}
