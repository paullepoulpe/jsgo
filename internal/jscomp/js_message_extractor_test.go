// Port of test/com/google/javascript/jscomp/JsMessageExtractorTest.java @ closure-compiler v20261021

/*
 * Copyright 2004 The Closure Compiler Authors.
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
	"strconv"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// jsMessageExtractorTest_TEST_ID_GENERATOR generates IDs of the form
// `MEANING_PARTCOUNT[PARTCOUNT...]`:
// PARTCOUNT = 'sN' for a string part with N == string length;
// PARTCOUNT = 'pN' for a placeholder with N == length of the canonical
// placeholder name.
var jsMessageExtractorTest_TEST_ID_GENERATOR JsMessage_IdGenerator = jsMessageExtractorTest_idGenerator{}

type jsMessageExtractorTest_idGenerator struct{}

func (jsMessageExtractorTest_idGenerator) GenerateId(meaning *string, messageParts []JsMessage_Part) string {
	idBuilder := jutil.NewStringBuilder()
	idBuilder.Append(jutil.ToJavaString(meaning)).AppendChar('_')
	for _, messagePart := range messageParts {
		if messagePart.IsPlaceholder() {
			idBuilder.AppendChar('p').Append(strconv.Itoa(jutil.StringLength(messagePart.GetCanonicalPlaceholderName())))
		} else {
			idBuilder.AppendChar('s').Append(strconv.Itoa(jutil.StringLength(messagePart.GetString())))
		}
	}

	return idBuilder.String()
}

// Java: JsMessageExtractorTest.extractMessages(JsMessage.IdGenerator, String)
func jsMessageExtractorTest_extractMessages(idGenerator JsMessage_IdGenerator, sourceCode string) []*JsMessage {
	return NewJsMessageExtractor(idGenerator).
		ExtractMessages(SourceFile_FromCode("testcode", sourceCode))
}

// Java: JsMessageExtractorTest.extractMessage(JsMessage.IdGenerator, String)
func jsMessageExtractorTest_extractMessage(t *testing.T, idGenerator JsMessage_IdGenerator, sourceCode string) *JsMessage {
	t.Helper()
	messages := jsMessageExtractorTest_extractMessages( /* idGenerator= */ idGenerator, sourceCode)
	if len(messages) != 1 {
		t.Fatalf("%d messages, want 1", len(messages))
	}
	return messages[0]
}

func jsMessageExtractorTest_assertThrows(t *testing.T, f func()) jutil.IRuntimeException {
	t.Helper()
	var caught jutil.IRuntimeException
	jutil.Try(f, func(e jutil.IRuntimeException) { caught = e })
	if caught == nil {
		t.Fatalf("expected a RuntimeException")
	}
	return caught
}

func jsMessageExtractorTest_assertContains(t *testing.T, s string, sub string) {
	t.Helper()
	if !strings.Contains(s, sub) {
		t.Errorf("%q does not contain %q", s, sub)
	}
}

func TestJsMessageExtractor_SyntaxError1(t *testing.T) {
	e := jsMessageExtractorTest_assertThrows(t, func() { jsMessageExtractorTest_extractMessage(t, nil, "if (true) {}}") })

	jsMessageExtractorTest_assertContains(t, e.GetMessage(), "JSCompiler errors\n")
	jsMessageExtractorTest_assertContains(t, e.GetMessage(), "testcode:1:13: ERROR - [JSC_PARSE_ERROR] Parse error")
	jsMessageExtractorTest_assertContains(t, e.GetMessage(), "if (true) {}}\n")
}

func TestJsMessageExtractor_SyntaxError2(t *testing.T) {
	e := jsMessageExtractorTest_assertThrows(t, func() { jsMessageExtractorTest_extractMessage(t, nil, "\nif (true) {}}") })

	jsMessageExtractorTest_assertContains(t, e.GetMessage(), "JSCompiler errors\n")
	jsMessageExtractorTest_assertContains(t, e.GetMessage(), "testcode:2:13: ERROR - [JSC_PARSE_ERROR] Parse error")
	jsMessageExtractorTest_assertContains(t, e.GetMessage(), "if (true) {}}\n")
}

func TestJsMessageExtractor_ExtractNewStyleMessage1(t *testing.T) {
	// A simple message with no description.
	jsMessageExtractorTest_assertEquals(t,
		NewJsMessage_Builder().
			SetKey("MSG_SILLY").
			SetId("MSG_SILLY").
			AppendStringPart("silly test message").
			Build(),
		jsMessageExtractorTest_extractMessage(t, nil, "var MSG_SILLY = goog.getMsg('silly test message');"))
}

func jsMessageExtractorTest_map(kv ...string) *jutil.LinkedHashMap[string, string] {
	m := jutil.NewLinkedHashMap[string, string]()
	for i := 0; i < len(kv); i += 2 {
		m.Put(kv[i], kv[i+1])
	}
	return m
}

func TestJsMessageExtractor_OriginalCodeAndExampleMaps(t *testing.T) {
	// A message with placeholders and original code annotations.
	jsMessageExtractorTest_assertEquals(t,
		NewJsMessage_Builder().
			SetKey("MSG_WELCOME").
			SetId("MSG_WELCOME").
			AppendStringPart("Hi ").
			AppendJsPlaceholderReference("interpolation_0").
			AppendStringPart("! Welcome to ").
			AppendJsPlaceholderReference("interpolation_1").
			AppendStringPart(".").
			SetPlaceholderNameToOriginalCodeMap(
				jsMessageExtractorTest_map(
					"interpolation_0", "foo.getUserName()",
					"interpolation_1", "bar.getProductName()")).
			SetPlaceholderNameToExampleMap(
				jsMessageExtractorTest_map(
					"interpolation_0", "Ginny Weasley",
					"interpolation_1", "Google Muggle Finder")).
			SetDesc(jutil.StringPtr("The welcome message.")).
			Build(),
		jsMessageExtractorTest_extractMessage(t, nil,
			"/** @desc The welcome message. */\n"+
				"var MSG_WELCOME = goog.getMsg(\n"+
				"    'Hi {$interpolation_0}! Welcome to {$interpolation_1}.',\n"+
				"    {\n"+
				"        'interpolation_0': 'magic-string-0',\n"+
				"        'interpolation_1': 'magic-string-1',\n"+
				"    },\n"+
				"    {\n"+
				"        original_code: {\n"+
				"            'interpolation_0': 'foo.getUserName()',\n"+
				"            'interpolation_1': 'bar.getProductName()',\n"+
				"        },\n"+
				"        example: {\n"+
				"            'interpolation_0': 'Ginny Weasley',\n"+
				"            'interpolation_1': 'Google Muggle Finder',\n"+
				"        },\n"+
				"    },\n"+
				");\n"))
}

func TestJsMessageExtractor_OriginalCodeAndExampleMapsForDeclareIcuTemplate(t *testing.T) {
	// A message with placeholders and original code annotations.
	jsMessageExtractorTest_assertEquals(t,
		NewJsMessage_Builder().
			SetKey("MSG_WELCOME").
			AppendStringPart("Hi ").                                // "s3" in the ID
			AppendCanonicalPlaceholderReference("INTERPOLATION_0"). // "p15" in the ID
			AppendStringPart("! Welcome to ").                      // "s13" in the ID
			AppendCanonicalPlaceholderReference("INTERPOLATION_1"). // "p15" in the ID
			AppendStringPart(".").                                  // "s1" in the ID
			SetPlaceholderNameToOriginalCodeMap(
				jsMessageExtractorTest_map(
					"INTERPOLATION_0", "foo.getUserName()",
					"INTERPOLATION_1", "bar.getProductName()")).
			SetPlaceholderNameToExampleMap(
				jsMessageExtractorTest_map(
					"INTERPOLATION_0", "Ginny Weasley",
					"INTERPOLATION_1", "Google Muggle Finder")).
			SetDesc(jutil.StringPtr("The welcome message.")).
			SetId("MSG_WELCOME_s3p15s13p15s1").
			Build(),
		jsMessageExtractorTest_extractMessage(t,
			jsMessageExtractorTest_TEST_ID_GENERATOR,
			"var MSG_WELCOME = declareIcuTemplate(\n"+
				"    'Hi {INTERPOLATION_0}! Welcome to {INTERPOLATION_1}.',\n"+
				"    {\n"+
				"        description: 'The welcome message.',\n"+
				"        example: {\n"+
				"            'INTERPOLATION_0': 'Ginny Weasley',\n"+
				"            'INTERPOLATION_1': 'Google Muggle Finder',\n"+
				"        },\n"+
				"    },\n"+
				");\n"))
}

func TestJsMessageExtractor_ExtractNewStyleMessage2(t *testing.T) {
	// A message with placeholders and meta data.
	jsMessageExtractorTest_assertEquals(t,
		NewJsMessage_Builder().
			SetKey("MSG_WELCOME").
			SetId("MSG_WELCOME").
			AppendStringPart("Hi ").
			AppendJsPlaceholderReference("userName").
			AppendStringPart("! Welcome to ").
			AppendJsPlaceholderReference("product").
			AppendStringPart(".").
			SetDesc(jutil.StringPtr("The welcome message.")).
			Build(),
		jsMessageExtractorTest_extractMessage(t, nil,
			"/**\n"+
				" * @desc The welcome\n"+
				" *   message.\n"+
				" */\n"+
				"var MSG_WELCOME = goog.getMsg(\n"+
				"    'Hi {$userName}! Welcome to {$product}.',\n"+
				"    {userName: someUserName, product: getProductName()});\n"))
}

func TestJsMessageExtractor_DuplicateUnnamedVariables(t *testing.T) {
	// Make sure that duplicate unnamed variables don't get swallowed when
	// using a Google-specific ID generator.
	msgs := jsMessageExtractorTest_extractMessages(nil,
		"function a() {\n"+
			"  var MSG_UNNAMED_2 = goog.getMsg('foo');\n"+
			"}\n"+
			"function b() {\n"+
			"  var MSG_UNNAMED_2 = goog.getMsg('bar');\n"+
			"}\n")

	if len(msgs) != 2 {
		t.Fatalf("%d messages, want 2", len(msgs))
	}
	if msgs[0].AsJsMessageString() != "foo" || msgs[1].AsJsMessageString() != "bar" {
		t.Errorf("messages %q, %q", msgs[0].AsJsMessageString(), msgs[1].AsJsMessageString())
	}
}

func TestJsMessageExtractor_MeaningAnnotation(t *testing.T) {
	msgs := jsMessageExtractorTest_extractMessages(nil,
		"var MSG_UNNAMED_1 = goog.getMsg('foo');\n"+
			"var MSG_UNNAMED_2 = goog.getMsg('foo');\n")
	if len(msgs) != 2 {
		t.Fatalf("%d messages, want 2", len(msgs))
	}
	if msgs[0].GetId() != msgs[1].GetId() {
		t.Errorf("ids %q and %q differ", msgs[0].GetId(), msgs[1].GetId())
	}
	jsMessageExtractorTest_assertEquals(t, msgs[0], msgs[1])

	msgs = jsMessageExtractorTest_extractMessages(nil,
		"var MSG_UNNAMED_1 = goog.getMsg('foo');\n"+
			"/** @meaning bar */ var MSG_UNNAMED_2 = goog.getMsg('foo');\n")
	if len(msgs) != 2 {
		t.Fatalf("%d messages, want 2", len(msgs))
	}
	if msgs[0].GetId() == msgs[1].GetId() {
		t.Errorf("ids are both %q", msgs[0].GetId())
	}
}

// Java: JsMessageExtractorTest.assertEquals(JsMessage, JsMessage)
func jsMessageExtractorTest_assertEquals(t *testing.T, expected *JsMessage, actual *JsMessage) {
	t.Helper()
	if actual.GetId() != expected.GetId() {
		t.Errorf("id = %q, want %q", actual.GetId(), expected.GetId())
	}
	if actual.GetKey() != expected.GetKey() {
		t.Errorf("key = %q, want %q", actual.GetKey(), expected.GetKey())
	}
	ap, ep := actual.GetParts(), expected.GetParts()
	same := len(ap) == len(ep)
	for i := 0; same && i < len(ap); i++ {
		same = ap[i] == ep[i]
	}
	if !same {
		t.Errorf("parts = %v, want %v", ap, ep)
	}
	if !jutil.ObjectsEquals(actual.JsPlaceholderNames(), expected.JsPlaceholderNames()) {
		t.Errorf("jsPlaceholderNames = %v, want %v", actual.JsPlaceholderNames(), expected.JsPlaceholderNames())
	}
	if !jutil.ObjectsEquals(actual.GetDesc(), expected.GetDesc()) {
		t.Errorf("desc = %v, want %v", jutil.ToJavaString(actual.GetDesc()), jutil.ToJavaString(expected.GetDesc()))
	}
	if !jutil.ObjectsEquals(actual.GetMeaning(), expected.GetMeaning()) {
		t.Errorf("meaning = %v, want %v", jutil.ToJavaString(actual.GetMeaning()), jutil.ToJavaString(expected.GetMeaning()))
	}
}
