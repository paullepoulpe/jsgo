// Port of com/google/javascript/jscomp/JsMessageTest.java @ closure-compiler v20261015

/*
 * Copyright 2009 The Closure Compiler Authors.
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
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

func TestJsMessage_IsLowerCamelCaseWithNumericSuffixes(t *testing.T) {
	assertEqual(t, "name", JsMessage_IsLowerCamelCaseWithNumericSuffixes("name"), true)
	assertEqual(t, "NAME", JsMessage_IsLowerCamelCaseWithNumericSuffixes("NAME"), false)
	assertEqual(t, "Name", JsMessage_IsLowerCamelCaseWithNumericSuffixes("Name"), false)

	assertEqual(t, "a4Letter", JsMessage_IsLowerCamelCaseWithNumericSuffixes("a4Letter"), true)
	assertEqual(t, "A4_LETTER", JsMessage_IsLowerCamelCaseWithNumericSuffixes("A4_LETTER"), false)

	assertEqual(t, "startSpan_1_23", JsMessage_IsLowerCamelCaseWithNumericSuffixes("startSpan_1_23"), true)
	assertEqual(t, "startSpan_1_23b", JsMessage_IsLowerCamelCaseWithNumericSuffixes("startSpan_1_23b"), false)
	assertEqual(t, "START_SPAN_1_23", JsMessage_IsLowerCamelCaseWithNumericSuffixes("START_SPAN_1_23"), false)

	assertEqual(t, "empty", JsMessage_IsLowerCamelCaseWithNumericSuffixes(""), false)
}

func TestJsMessage_ToLowerCamelCaseWithNumericSuffixes(t *testing.T) {
	assertEqual(t, "NAME", JsMessage_ToLowerCamelCaseWithNumericSuffixes("NAME"), "name")
	assertEqual(t, "A4_LETTER", JsMessage_ToLowerCamelCaseWithNumericSuffixes("A4_LETTER"), "a4Letter")
	assertEqual(t, "START_SPAN_1_23", JsMessage_ToLowerCamelCaseWithNumericSuffixes("START_SPAN_1_23"),
		"startSpan_1_23")
}

func TestJsMessage_IsEmpty(t *testing.T) {
	assertEqual(t, "1", newTestMessageBuilder("MSG_KEY").Build().IsEmpty(), true)
	assertEqual(t, "2", newTestMessageBuilder("MSG_KEY").AppendStringPart("").Build().IsEmpty(), true)
	assertEqual(t, "3",
		newTestMessageBuilder("MSG_KEY").
			AppendStringPart("").
			AppendStringPart("").
			Build().
			IsEmpty(),
		true)
	assertEqual(t, "4",
		newTestMessageBuilder("MSG_KEY").
			AppendStringPart("s").
			AppendStringPart("").
			Build().
			IsEmpty(),
		false)
	assertEqual(t, "5",
		newTestMessageBuilder("MSG_KEY").AppendJsPlaceholderReference("ph").Build().IsEmpty(),
		false)
}

// newTestMessageBuilder returns a new message builder that uses the given
// string as both its key and its ID.
func newTestMessageBuilder(keyAndId string) *JsMessage_Builder {
	return NewJsMessage_Builder().SetKey(keyAndId).SetId(keyAndId)
}

func TestJsMessage_HashValues(t *testing.T) {
	var answerString64 int64 = 0x43ec5d9731515874
	s := "Hello, world"
	assertEqual(t, "Hello, world", JsMessage_Hash_Hash64(&s), answerString64)

	var answerEmpty64 int64 = 0x468d9ea2c42361aa
	e := ""
	assertEqual(t, "empty", JsMessage_Hash_Hash64(&e), answerEmpty64)
}

func TestJsMessage_NoAlternateId(t *testing.T) {
	msg := newTestMessageBuilder("MSG_SOME_KEY").SetDesc(jutil.StringPtr("Hello.")).Build()
	assertEqual(t, "desc", strPtrValue(msg.GetDesc()), "Hello.")
	assertEqual(t, "id", msg.GetId(), "MSG_SOME_KEY")
	if msg.GetAlternateId() != nil {
		t.Errorf("alternate id %q, want null", *msg.GetAlternateId())
	}
}

func TestJsMessage_SelfReferentialAlternateId(t *testing.T) {
	msg := newTestMessageBuilder("MSG_SOME_NAME").
		SetDesc(jutil.StringPtr("Hello.")).
		SetAlternateId(jutil.StringPtr("MSG_SOME_NAME")).
		Build()
	assertEqual(t, "desc", strPtrValue(msg.GetDesc()), "Hello.")
	assertEqual(t, "id", msg.GetId(), "MSG_SOME_NAME")
	if msg.GetAlternateId() != nil {
		t.Errorf("alternate id %q, want null", *msg.GetAlternateId())
	}
}

func TestJsMessage_AlternateId(t *testing.T) {
	msg := NewJsMessage_Builder().
		SetKey("MSG_KEY").
		SetDesc(jutil.StringPtr("Hello.")).
		SetAlternateId(jutil.StringPtr("foo")).
		SetMeaning(jutil.StringPtr("meaning")).
		SetId("meaning").
		AppendJsPlaceholderReference("placeholder0").
		AppendJsPlaceholderReference("placeholder1").
		AppendStringPart("part0").
		AppendStringPart("part1").
		AppendStringPart("part2").
		Build()
	assertEqual(t, "desc", strPtrValue(msg.GetDesc()), "Hello.")
	assertEqual(t, "id", msg.GetId(), "meaning")
	assertEqual(t, "alternate id", strPtrValue(msg.GetAlternateId()), "foo")
}

func TestJsMessage_MessageWithGenderedVariants(t *testing.T) {
	msg := NewJsMessage_Builder().
		SetKey("MSG_KEY").
		SetId("1234567890").
		SetDesc(jutil.StringPtr("A welcome message.")).
		AddGenderedMessageKey(JsMessage_GrammaticalGenderCase_MASCULINE).
		AppendStringPartForGender(JsMessage_GrammaticalGenderCase_MASCULINE, "Beinvenido, ").
		AppendJsPlaceholderReferenceForGender(JsMessage_GrammaticalGenderCase_MASCULINE, "placeholder").
		AddGenderedMessageKey(JsMessage_GrammaticalGenderCase_FEMININE).
		AppendStringPartForGender(JsMessage_GrammaticalGenderCase_FEMININE, "Beinvenida, ").
		AppendJsPlaceholderReferenceForGender(JsMessage_GrammaticalGenderCase_FEMININE, "placeholder").
		AddGenderedMessageKey(JsMessage_GrammaticalGenderCase_NEUTER).
		AppendStringPartForGender(JsMessage_GrammaticalGenderCase_NEUTER, "Te damos la bienvenida, ").
		AppendJsPlaceholderReferenceForGender(JsMessage_GrammaticalGenderCase_NEUTER, "placeholder").
		AddGenderedMessageKey(JsMessage_GrammaticalGenderCase_OTHER).
		AppendStringPartForGender(
			JsMessage_GrammaticalGenderCase_OTHER, "Te damos la bienvenida - OTHER, ").
		AppendJsPlaceholderReferenceForGender(JsMessage_GrammaticalGenderCase_OTHER, "placeholder").
		Build()

	for _, c := range []struct {
		g    JsMessage_GrammaticalGenderCase
		text string
	}{
		{JsMessage_GrammaticalGenderCase_MASCULINE, "Beinvenido, "},
		{JsMessage_GrammaticalGenderCase_FEMININE, "Beinvenida, "},
		{JsMessage_GrammaticalGenderCase_NEUTER, "Te damos la bienvenida, "},
		{JsMessage_GrammaticalGenderCase_OTHER, "Te damos la bienvenida - OTHER, "},
	} {
		parts := msg.GetGenderedMessageParts(c.g)
		assertEqual(t, c.g.String()+" size", len(parts), 2)
		assertEqual(t, c.g.String()+" string", parts[0].GetString(), c.text)
		assertEqual(t, c.g.String()+" placeholder", parts[1].GetJsPlaceholderName(), "placeholder")
	}
}

func TestJsMessage_GetGenderedVariants(t *testing.T) {
	msg := NewJsMessage_Builder().
		SetKey("MSG_KEY").
		SetId("1234567890").
		SetDesc(jutil.StringPtr("A welcome message.")).
		AddGenderedMessageKey(JsMessage_GrammaticalGenderCase_MASCULINE).
		AppendStringPartForGender(JsMessage_GrammaticalGenderCase_MASCULINE, "Bienvenido, ").
		AddGenderedMessageKey(JsMessage_GrammaticalGenderCase_FEMININE).
		AppendStringPartForGender(JsMessage_GrammaticalGenderCase_FEMININE, "Bienvenida, ").
		AddGenderedMessageKey(JsMessage_GrammaticalGenderCase_NEUTER).
		AppendStringPartForGender(JsMessage_GrammaticalGenderCase_NEUTER, "Te damos la bienvenida, ").
		AddGenderedMessageKey(JsMessage_GrammaticalGenderCase_OTHER).
		AppendStringPartForGender(
			JsMessage_GrammaticalGenderCase_OTHER, "Te damos la bienvenida - OTHER, ").
		Build()

	expectedGenderedVariants := []JsMessage_GrammaticalGenderCase{
		JsMessage_GrammaticalGenderCase_MASCULINE,
		JsMessage_GrammaticalGenderCase_FEMININE,
		JsMessage_GrammaticalGenderCase_NEUTER,
		JsMessage_GrammaticalGenderCase_OTHER,
	}

	assertSliceEqual(t, "variants", msg.GetGenderedMessageVariants(), expectedGenderedVariants)
}

func TestJsMessage_GetGenderedVariantsNoNeuter(t *testing.T) {
	msg := NewJsMessage_Builder().
		SetKey("MSG_KEY").
		SetId("1234567890").
		SetDesc(jutil.StringPtr("A welcome message.")).
		AddGenderedMessageKey(JsMessage_GrammaticalGenderCase_MASCULINE).
		AppendStringPartForGender(JsMessage_GrammaticalGenderCase_MASCULINE, "Bienvenido, ").
		AddGenderedMessageKey(JsMessage_GrammaticalGenderCase_FEMININE).
		AppendStringPartForGender(JsMessage_GrammaticalGenderCase_FEMININE, "Bienvenida, ").
		AddGenderedMessageKey(JsMessage_GrammaticalGenderCase_OTHER).
		AppendStringPartForGender(
			JsMessage_GrammaticalGenderCase_OTHER, "Te damos la bienvenida - OTHER, ").
		Build()

	expectedGenderedVariants := []JsMessage_GrammaticalGenderCase{
		JsMessage_GrammaticalGenderCase_MASCULINE,
		JsMessage_GrammaticalGenderCase_FEMININE,
		JsMessage_GrammaticalGenderCase_OTHER,
	}

	assertSliceEqual(t, "variants", msg.GetGenderedMessageVariants(), expectedGenderedVariants)
}

func TestJsMessage_GenderedMessageAsString(t *testing.T) {
	msg := NewJsMessage_Builder().
		SetKey("MSG_KEY").
		SetId("1234567890").
		SetDesc(jutil.StringPtr("A welcome message.")).
		AddGenderedMessageKey(JsMessage_GrammaticalGenderCase_MASCULINE).
		AppendStringPartForGender(JsMessage_GrammaticalGenderCase_MASCULINE, "Beinvenido, ").
		AppendJsPlaceholderReferenceForGender(JsMessage_GrammaticalGenderCase_MASCULINE, "placeholder").
		AddGenderedMessageKey(JsMessage_GrammaticalGenderCase_FEMININE).
		AppendStringPartForGender(JsMessage_GrammaticalGenderCase_FEMININE, "Beinvenida, ").
		AppendJsPlaceholderReferenceForGender(JsMessage_GrammaticalGenderCase_FEMININE, "placeholder").
		AddGenderedMessageKey(JsMessage_GrammaticalGenderCase_NEUTER).
		AppendStringPartForGender(JsMessage_GrammaticalGenderCase_NEUTER, "Te damos la bienvenida, ").
		AppendJsPlaceholderReferenceForGender(JsMessage_GrammaticalGenderCase_NEUTER, "placeholder").
		AddGenderedMessageKey(JsMessage_GrammaticalGenderCase_OTHER).
		AppendStringPartForGender(
			JsMessage_GrammaticalGenderCase_OTHER, "Te damos la bienvenida - OTHER, ").
		AppendJsPlaceholderReferenceForGender(JsMessage_GrammaticalGenderCase_OTHER, "placeholder").
		Build()

	assertEqual(t, "js M", msg.AsJsMessageStringForGender(JsMessage_GrammaticalGenderCase_MASCULINE),
		"Beinvenido, {$placeholder}")
	assertEqual(t, "js F", msg.AsJsMessageStringForGender(JsMessage_GrammaticalGenderCase_FEMININE),
		"Beinvenida, {$placeholder}")
	assertEqual(t, "js N", msg.AsJsMessageStringForGender(JsMessage_GrammaticalGenderCase_NEUTER),
		"Te damos la bienvenida, {$placeholder}")
	assertEqual(t, "js O", msg.AsJsMessageStringForGender(JsMessage_GrammaticalGenderCase_OTHER),
		"Te damos la bienvenida - OTHER, {$placeholder}")

	assertEqual(t, "icu M", msg.AsIcuMessageStringForGender(JsMessage_GrammaticalGenderCase_MASCULINE),
		"Beinvenido, {PLACEHOLDER}")
	assertEqual(t, "icu F", msg.AsIcuMessageStringForGender(JsMessage_GrammaticalGenderCase_FEMININE),
		"Beinvenida, {PLACEHOLDER}")
	assertEqual(t, "icu N", msg.AsIcuMessageStringForGender(JsMessage_GrammaticalGenderCase_NEUTER),
		"Te damos la bienvenida, {PLACEHOLDER}")
	assertEqual(t, "icu O", msg.AsIcuMessageStringForGender(JsMessage_GrammaticalGenderCase_OTHER),
		"Te damos la bienvenida - OTHER, {PLACEHOLDER}")
}

func TestJsMessage_GenderedMessageAsStringInvalidCall(t *testing.T) {
	genderedMsg := NewJsMessage_Builder().
		SetKey("MSG_KEY").
		SetId("1234567890").
		SetDesc(jutil.StringPtr("A welcome message.")).
		AddGenderedMessageKey(JsMessage_GrammaticalGenderCase_MASCULINE).
		AppendStringPartForGender(JsMessage_GrammaticalGenderCase_MASCULINE, "Bienvenido, ").
		AddGenderedMessageKey(JsMessage_GrammaticalGenderCase_FEMININE).
		AppendStringPartForGender(JsMessage_GrammaticalGenderCase_FEMININE, "Bienvenida, ").
		AddGenderedMessageKey(JsMessage_GrammaticalGenderCase_OTHER).
		AppendStringPartForGender(
			JsMessage_GrammaticalGenderCase_OTHER, "Te damos la bienvenida - OTHER, ").
		Build()

	// No grammatical gender argument provided
	assertThrows[jutil.IUnsupportedOperationException](t, func() { genderedMsg.AsJsMessageString() })
	assertThrows[jutil.IUnsupportedOperationException](t, func() { genderedMsg.AsIcuMessageString() })

	// No message for the provided grammatical gender
	assertThrows[jutil.IIllegalArgumentException](t, func() {
		genderedMsg.AsJsMessageStringForGender(JsMessage_GrammaticalGenderCase_NEUTER)
	})
	assertThrows[jutil.IIllegalArgumentException](t, func() {
		genderedMsg.AsIcuMessageStringForGender(JsMessage_GrammaticalGenderCase_NEUTER)
	})
}
