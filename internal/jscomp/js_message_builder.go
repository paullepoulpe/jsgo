// Port of com/google/javascript/jscomp/JsMessage.java @ closure-compiler v20261021

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
	"encoding/binary"
	"maps"
	"slices"
	"unicode/utf8"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// JsMessage_IsLowerCamelCaseWithNumericSuffixes returns whether a string is
// nonempty, begins with a lowercase letter, and contains only digits and
// underscores after the first underscore (`\A[a-z][a-zA-Z0-9]*[_0-9]*\z`).
//
// Java: JsMessage.isLowerCamelCaseWithNumericSuffixes(String)
func JsMessage_IsLowerCamelCaseWithNumericSuffixes(input string) bool {
	if len(input) == 0 || input[0] < 'a' || input[0] > 'z' {
		return false
	}
	i := 1
	for i < len(input) {
		c := input[i]
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') {
			i++
		} else {
			break
		}
	}
	for ; i < len(input); i++ {
		c := input[i]
		if c != '_' && (c < '0' || c > '9') {
			return false
		}
	}
	return true
}

// JsMessage_ToLowerCamelCaseWithNumericSuffixes converts the given string
// from upper-underscore case to lower-camel case, preserving numeric
// suffixes. For example: "NAME" -> "name", "A4_LETTER" -> "a4Letter",
// "START_SPAN_1_23" -> "startSpan_1_23".
//
// Java: JsMessage.toLowerCamelCaseWithNumericSuffixes(String)
func JsMessage_ToLowerCamelCaseWithNumericSuffixes(input string) string {
	// Determine where the numeric suffixes begin by scanning backwards.
	// Surrogate code units are never Character.isDigit, so BMP UTF-8 decoding
	// finds the exact same suffix boundary as UTF-16 code-unit scanning with
	// zero intermediate slice allocations.
	suffixStart := len(input)
	if jutil.IsASCII(input) {
		for suffixStart > 0 {
			numberStart := suffixStart
			for numberStart > 0 && input[numberStart-1] >= '0' && input[numberStart-1] <= '9' {
				numberStart--
			}
			if numberStart > 0 && numberStart < suffixStart && input[numberStart-1] == '_' {
				suffixStart = numberStart - 1
			} else {
				break
			}
		}
	} else {
		for suffixStart > 0 {
			numberStart := suffixStart
			var lastRune rune
			var lastSize int
			for numberStart > 0 {
				r, size := utf8.DecodeLastRuneInString(input[:numberStart])
				if r < 0x10000 && (r != utf8.RuneError || size != 1) && jutil.CharacterIsDigit(uint16(r)) {
					numberStart -= size
				} else {
					lastRune = r
					lastSize = size
					break
				}
			}
			if numberStart > 0 && numberStart < suffixStart && lastSize == 1 && lastRune == '_' {
				suffixStart = numberStart - 1
			} else {
				break
			}
		}
	}

	if suffixStart == len(input) {
		return jutil.CaseFormat_UpperUnderscoreToLowerCamel(input)
	}
	return jutil.CaseFormat_UpperUnderscoreToLowerCamel(input[:suffixStart]) + input[suffixStart:]
}

// JsMessage_IsCanonicalPlaceholderNameFormat reports whether name matches
// `\A[A-Z0-9_]*\z` (the canonical format for placeholder names in XTB/XMB).
//
// Java: JsMessage.isCanonicalPlaceholderNameFormat(String)
func JsMessage_IsCanonicalPlaceholderNameFormat(name string) bool {
	for i := 0; i < len(name); i++ {
		c := name[i]
		if (c < 'A' || c > 'Z') && (c < '0' || c > '9') && c != '_' {
			return false
		}
	}
	return true
}

// JsMessage_Builder contains functionality for creating JS messages.
type JsMessage_Builder struct {
	key         string
	meaning     string
	desc        string
	id          string
	alternateId string
	sourceName  string

	parts                              []JsMessage_Part
	genderedMessageMap                 map[JsMessage_GrammaticalGenderCase][]JsMessage_Part
	genderedMessageOrder               []JsMessage_GrammaticalGenderCase
	jsPlaceholderNames                 map[string]struct{}
	jsPlaceholderNamesOrder            []string
	canonicalPlaceholderNames          map[string]struct{}
	canonicalPlaceholderNamesOrder     []string
	placeholderNameToExampleMap        map[string]string
	placeholderNameToExampleOrder      []string
	placeholderNameToOriginalCodeMap   map[string]string
	placeholderNameToOriginalCodeOrder []string

	hasKey         bool
	hasMeaning     bool
	hasDesc        bool
	hasId          bool
	hasAlternateId bool
	hasSourceName  bool
	isAnonymous    bool
	isExternal     bool
}

// NewJsMessage_Builder is new JsMessage.Builder().
func NewJsMessage_Builder() *JsMessage_Builder {
	return &JsMessage_Builder{}
}

func (b *JsMessage_Builder) addJsPlaceholderName(name string) {
	if b.jsPlaceholderNames == nil {
		b.jsPlaceholderNames = make(map[string]struct{})
	}
	if _, exists := b.jsPlaceholderNames[name]; !exists {
		b.jsPlaceholderNames[name] = struct{}{}
		b.jsPlaceholderNamesOrder = append(b.jsPlaceholderNamesOrder, name)
	}
}

func (b *JsMessage_Builder) addCanonicalPlaceholderName(name string) {
	if b.canonicalPlaceholderNames == nil {
		b.canonicalPlaceholderNames = make(map[string]struct{})
	}
	if _, exists := b.canonicalPlaceholderNames[name]; !exists {
		b.canonicalPlaceholderNames[name] = struct{}{}
		b.canonicalPlaceholderNamesOrder = append(b.canonicalPlaceholderNamesOrder, name)
	}
}

// GetKey gets the message's key (e.g. "MSG_HELLO").
//
// Java: JsMessage.Builder.getKey()
func (b *JsMessage_Builder) GetKey() *string {
	if !b.hasKey {
		return nil
	}
	return &b.key
}

// SetKey sets a key that uniquely identifies this message (e.g. "MSG_HELLO").
//
// Java: JsMessage.Builder.setKey(String)
func (b *JsMessage_Builder) SetKey(key string) *JsMessage_Builder {
	b.key = key
	b.hasKey = true
	return b
}

// SetSourceName sets the message's sourceName.
//
// Java: JsMessage.Builder.setSourceName(String)
func (b *JsMessage_Builder) SetSourceName(sourceName *string) *JsMessage_Builder {
	if sourceName == nil {
		b.sourceName = ""
		b.hasSourceName = false
	} else {
		b.sourceName = *sourceName
		b.hasSourceName = true
	}
	return b
}

// SetSourceNameString sets a non-null sourceName without pointer boxing.
func (b *JsMessage_Builder) SetSourceNameString(sourceName string) *JsMessage_Builder {
	b.sourceName = sourceName
	b.hasSourceName = true
	return b
}

// Java: JsMessage.Builder.appendPart(Part)
func (b *JsMessage_Builder) AppendPart(part JsMessage_Part) *JsMessage_Builder {
	b.parts = append(b.parts, part)
	if part.IsPlaceholder() {
		b.addJsPlaceholderName(part.GetJsPlaceholderName())
		b.addCanonicalPlaceholderName(part.GetCanonicalPlaceholderName())
	}
	return b
}

// Java: JsMessage.Builder.appendParts(List<Part>)
func (b *JsMessage_Builder) AppendParts(parts []JsMessage_Part) *JsMessage_Builder {
	for _, part := range parts {
		b.AppendPart(part)
	}
	return b
}

// appendGenderedPart appends to genderedMessageMap.get(genderCase), which must
// exist (Java throws NullPointerException on the null list).
func (b *JsMessage_Builder) appendGenderedPart(genderCase JsMessage_GrammaticalGenderCase, part JsMessage_Part) {
	parts, ok := b.genderedMessageMap[genderCase]
	if !ok {
		panic(jutil.NewNullPointerException())
	}
	b.genderedMessageMap[genderCase] = append(parts, part)
}

// AppendJsPlaceholderReference appends a placeholder reference to the message.
//
// Java: JsMessage.Builder.appendJsPlaceholderReference(String)
func (b *JsMessage_Builder) AppendJsPlaceholderReference(name string) *JsMessage_Builder {
	b.parts = append(b.parts, JsMessage_PlaceholderReference_CreateForJsName(name))
	b.addJsPlaceholderName(name)
	return b
}

// AppendJsPlaceholderReferenceForGender appends a placeholder reference to the
// corresponding gendered message.
//
// Java: JsMessage.Builder.appendJsPlaceholderReference(GrammaticalGenderCase, String)
func (b *JsMessage_Builder) AppendJsPlaceholderReferenceForGender(genderCase JsMessage_GrammaticalGenderCase, name string) *JsMessage_Builder {
	b.appendGenderedPart(genderCase, JsMessage_PlaceholderReference_CreateForJsName(name))
	b.addJsPlaceholderName(name)
	return b
}

// AppendCanonicalPlaceholderReference appends a canonical placeholder
// reference to the message.
//
// Java: JsMessage.Builder.appendCanonicalPlaceholderReference(String)
func (b *JsMessage_Builder) AppendCanonicalPlaceholderReference(name string) *JsMessage_Builder {
	placeholder := JsMessage_PlaceholderReference_CreateForCanonicalName(name)
	b.parts = append(b.parts, placeholder)
	b.addJsPlaceholderName(placeholder.GetJsPlaceholderName())
	return b
}

// AppendCanonicalPlaceholderReferenceForGender appends a canonical placeholder
// reference to the corresponding gendered message.
//
// Java: JsMessage.Builder.appendCanonicalPlaceholderReference(GrammaticalGenderCase, String)
func (b *JsMessage_Builder) AppendCanonicalPlaceholderReferenceForGender(
	genderCase JsMessage_GrammaticalGenderCase, name string) *JsMessage_Builder {
	placeholder := JsMessage_PlaceholderReference_CreateForCanonicalName(name)
	b.appendGenderedPart(genderCase, placeholder)
	b.addJsPlaceholderName(placeholder.GetJsPlaceholderName())
	return b
}

// AppendStringPart appends a translatable string literal to the message.
//
// Java: JsMessage.Builder.appendStringPart(String)
func (b *JsMessage_Builder) AppendStringPart(part string) *JsMessage_Builder {
	b.parts = append(b.parts, JsMessage_StringPart_Create(part))
	return b
}

// AppendStringPartForGender appends a translatable string literal to the
// corresponding gendered message.
//
// Java: JsMessage.Builder.appendStringPart(GrammaticalGenderCase, String)
func (b *JsMessage_Builder) AppendStringPartForGender(genderCase JsMessage_GrammaticalGenderCase, part string) *JsMessage_Builder {
	b.appendGenderedPart(genderCase, JsMessage_StringPart_Create(part))
	return b
}

// AddGenderedMessageKey adds a gendered message key to genderedMessageParts map.
//
// Java: JsMessage.Builder.addGenderedMessageKey(GrammaticalGenderCase)
func (b *JsMessage_Builder) AddGenderedMessageKey(part JsMessage_GrammaticalGenderCase) *JsMessage_Builder {
	if b.genderedMessageMap == nil {
		b.genderedMessageMap = make(map[JsMessage_GrammaticalGenderCase][]JsMessage_Part, 4)
	}
	if _, exists := b.genderedMessageMap[part]; !exists {
		b.genderedMessageOrder = append(b.genderedMessageOrder, part)
	}
	b.genderedMessageMap[part] = nil
	return b
}

// Java: JsMessage.Builder.setPlaceholderNameToExampleMap(Map<String, String>)
func (b *JsMessage_Builder) SetPlaceholderNameToExampleMap(m map[string]string) *JsMessage_Builder {
	if len(m) == 0 {
		b.placeholderNameToExampleMap = nil
		b.placeholderNameToExampleOrder = nil
		return b
	}
	order := make([]string, 0, len(m))
	for k := range m { // order-independent: keys are sorted immediately after collection
		order = append(order, k)
	}
	slices.Sort(order)
	b.placeholderNameToExampleMap = maps.Clone(m)
	b.placeholderNameToExampleOrder = order
	return b
}

func (b *JsMessage_Builder) setPlaceholderNameToExampleOrdered(m map[string]string, order []string) *JsMessage_Builder {
	if len(m) == 0 {
		b.placeholderNameToExampleMap = nil
		b.placeholderNameToExampleOrder = nil
		return b
	}
	b.placeholderNameToExampleMap = maps.Clone(m)
	b.placeholderNameToExampleOrder = slices.Clone(order)
	return b
}

// Java: JsMessage.Builder.setPlaceholderNameToOriginalCodeMap(Map<String, String>)
func (b *JsMessage_Builder) SetPlaceholderNameToOriginalCodeMap(m map[string]string) *JsMessage_Builder {
	if len(m) == 0 {
		b.placeholderNameToOriginalCodeMap = nil
		b.placeholderNameToOriginalCodeOrder = nil
		return b
	}
	order := make([]string, 0, len(m))
	for k := range m { // order-independent: keys are sorted immediately after collection
		order = append(order, k)
	}
	slices.Sort(order)
	b.placeholderNameToOriginalCodeMap = maps.Clone(m)
	b.placeholderNameToOriginalCodeOrder = order
	return b
}

func (b *JsMessage_Builder) setPlaceholderNameToOriginalCodeOrdered(m map[string]string, order []string) *JsMessage_Builder {
	if len(m) == 0 {
		b.placeholderNameToOriginalCodeMap = nil
		b.placeholderNameToOriginalCodeOrder = nil
		return b
	}
	b.placeholderNameToOriginalCodeMap = maps.Clone(m)
	b.placeholderNameToOriginalCodeOrder = slices.Clone(order)
	return b
}

// SetDesc sets the description of the message, which helps translators.
//
// Java: JsMessage.Builder.setDesc(String)
func (b *JsMessage_Builder) SetDesc(desc *string) *JsMessage_Builder {
	if desc == nil {
		b.desc = ""
		b.hasDesc = false
	} else {
		b.desc = *desc
		b.hasDesc = true
	}
	return b
}

// SetDescString sets a non-null description without pointer boxing.
func (b *JsMessage_Builder) SetDescString(desc string) *JsMessage_Builder {
	b.desc = desc
	b.hasDesc = true
	return b
}

// SetDescOpt sets or clears the description without pointer boxing.
func (b *JsMessage_Builder) SetDescOpt(desc string, ok bool) *JsMessage_Builder {
	if !ok {
		b.desc = ""
		b.hasDesc = false
	} else {
		b.desc = desc
		b.hasDesc = true
	}
	return b
}

// SetMeaning sets the programmer-specified meaning of this message, which
// forces this message to translate differently.
//
// Java: JsMessage.Builder.setMeaning(String)
func (b *JsMessage_Builder) SetMeaning(meaning *string) *JsMessage_Builder {
	if meaning == nil {
		b.meaning = ""
		b.hasMeaning = false
	} else {
		b.meaning = *meaning
		b.hasMeaning = true
	}
	return b
}

// SetMeaningString sets a non-null meaning without pointer boxing.
func (b *JsMessage_Builder) SetMeaningString(meaning string) *JsMessage_Builder {
	b.meaning = meaning
	b.hasMeaning = true
	return b
}

// SetMeaningOpt sets or clears the meaning without pointer boxing.
func (b *JsMessage_Builder) SetMeaningOpt(meaning string, ok bool) *JsMessage_Builder {
	if !ok {
		b.meaning = ""
		b.hasMeaning = false
	} else {
		b.meaning = meaning
		b.hasMeaning = true
	}
	return b
}

// SetAlternateId sets the alternate message ID, to be used if the primary ID
// is not yet translated.
//
// Java: JsMessage.Builder.setAlternateId(String)
func (b *JsMessage_Builder) SetAlternateId(alternateId *string) *JsMessage_Builder {
	if alternateId == nil {
		b.alternateId = ""
		b.hasAlternateId = false
	} else {
		b.alternateId = *alternateId
		b.hasAlternateId = true
	}
	return b
}

// SetAlternateIdString sets a non-null alternate message ID without pointer boxing.
func (b *JsMessage_Builder) SetAlternateIdString(alternateId string) *JsMessage_Builder {
	b.alternateId = alternateId
	b.hasAlternateId = true
	return b
}

// SetAlternateIdOpt sets or clears the alternate message ID without pointer boxing.
func (b *JsMessage_Builder) SetAlternateIdOpt(alternateId string, ok bool) *JsMessage_Builder {
	if !ok {
		b.alternateId = ""
		b.hasAlternateId = false
	} else {
		b.alternateId = alternateId
		b.hasAlternateId = true
	}
	return b
}

// Java: JsMessage.Builder.setIsAnonymous(boolean)
func (b *JsMessage_Builder) SetIsAnonymous(isAnonymous bool) *JsMessage_Builder {
	b.isAnonymous = isAnonymous
	return b
}

// Java: JsMessage.Builder.setId(String)
func (b *JsMessage_Builder) SetId(id string) *JsMessage_Builder {
	if b.hasId {
		panic(jutil.NewIllegalStateExceptionMsg(
			"id already set to '" + b.id + "': cannot change it to '" + id + "'"))
	}
	b.id = id
	b.hasId = true
	return b
}

// Java: JsMessage.Builder.setIsExternalMsg(boolean)
func (b *JsMessage_Builder) SetIsExternalMsg(isExternalMsg bool) *JsMessage_Builder {
	b.isExternal = isExternalMsg
	return b
}

// HasParts gets whether at least one part has been appended.
//
// Java: JsMessage.Builder.hasParts()
func (b *JsMessage_Builder) HasParts() bool {
	return len(b.parts) > 0
}

// HasGenderedVariants gets whether the message has gendered variants.
//
// Java: JsMessage.Builder.hasGenderedVariants()
func (b *JsMessage_Builder) HasGenderedVariants() bool {
	return len(b.genderedMessageMap) > 0
}

// Java: JsMessage.Builder.getParts()
func (b *JsMessage_Builder) GetParts() []JsMessage_Part {
	return b.parts
}

// Java: JsMessage.Builder.build()
func (b *JsMessage_Builder) Build() *JsMessage {
	if !b.hasKey {
		panic(jutil.NewNullPointerExceptionMsg("key has not been set"))
	}
	if !b.hasId {
		panic(jutil.NewNullPointerExceptionMsg("id has not been set"))
	}
	jutil.CheckStateMsg(!b.isExternal || !b.isAnonymous, "a message cannot be both anonymous and external")

	// An alternate ID that points to itself is a no-op, so just omit it.
	hasAltId := b.hasAlternateId
	altId := b.alternateId
	if hasAltId && altId == b.id {
		hasAltId = false
		altId = ""
	}

	var genderedMap map[JsMessage_GrammaticalGenderCase][]JsMessage_Part
	var genderedOrder []JsMessage_GrammaticalGenderCase
	if len(b.genderedMessageOrder) > 0 {
		genderedMap = make(map[JsMessage_GrammaticalGenderCase][]JsMessage_Part, len(b.genderedMessageOrder))
		genderedOrder = slices.Clone(b.genderedMessageOrder)
		for _, k := range genderedOrder {
			genderedMap[k] = slices.Clone(b.genderedMessageMap[k])
		}
	}

	var exampleMap map[string]string
	var exampleOrder []string
	if len(b.placeholderNameToExampleOrder) > 0 {
		exampleMap = maps.Clone(b.placeholderNameToExampleMap)
		exampleOrder = slices.Clone(b.placeholderNameToExampleOrder)
	}

	var originalCodeMap map[string]string
	var originalCodeOrder []string
	if len(b.placeholderNameToOriginalCodeOrder) > 0 {
		originalCodeMap = maps.Clone(b.placeholderNameToOriginalCodeMap)
		originalCodeOrder = slices.Clone(b.placeholderNameToOriginalCodeOrder)
	}

	var jsPhNames map[string]struct{}
	var jsPhOrder []string
	if len(b.jsPlaceholderNamesOrder) > 0 {
		jsPhNames = maps.Clone(b.jsPlaceholderNames)
		jsPhOrder = slices.Clone(b.jsPlaceholderNamesOrder)
	}

	var canonicalPhNames map[string]struct{}
	var canonicalPhOrder []string
	if len(b.canonicalPlaceholderNamesOrder) > 0 {
		canonicalPhNames = maps.Clone(b.canonicalPlaceholderNames)
		canonicalPhOrder = slices.Clone(b.canonicalPlaceholderNamesOrder)
	}

	return &JsMessage{
		sourceName:                         b.sourceName,
		hasSourceName:                      b.hasSourceName,
		key:                                b.key,
		isAnonymous:                        b.isAnonymous,
		isExternal:                         b.isExternal,
		id:                                 b.id,
		parts:                              slices.Clone(b.parts),
		genderedMessagesMap:                genderedMap,
		genderedMessagesOrder:              genderedOrder,
		alternateId:                        altId,
		hasAlternateId:                     hasAltId,
		desc:                               b.desc,
		hasDesc:                            b.hasDesc,
		meaning:                            b.meaning,
		hasMeaning:                         b.hasMeaning,
		placeholderNameToExampleMap:        exampleMap,
		placeholderNameToExampleOrder:      exampleOrder,
		placeholderNameToOriginalCodeMap:   originalCodeMap,
		placeholderNameToOriginalCodeOrder: originalCodeOrder,
		jsPlaceholderNames:                 jsPhNames,
		jsPlaceholderNamesOrder:            jsPhOrder,
		canonicalPlaceholderNames:          canonicalPhNames,
		canonicalPlaceholderNamesOrder:     canonicalPhOrder,
	}
}

// JsMessage_Hash contains routines for hashing.
type JsMessage_Hash struct{}

const (
	// Default hash seed (64 bit)
	jsMessage_Hash_SEED64 uint64 = 0x2b992ddfa23249d6 // part of pi, arbitrary

	// Hash constant (64 bit)
	jsMessage_Hash_CONSTANT64 uint64 = 0xe08c1d668b756f82 // part of golden ratio, arbitrary
)

// JsMessage_Hash_Hash64 hashes a string to a 64 bit value. The digits of
// pi are used for the hash seed.
//
// Java: JsMessage.Hash.hash64(String)
func JsMessage_Hash_Hash64(value *string) int64 {
	if value == nil {
		return jsMessage_Hash_hash64Range(nil, 0, 0, jsMessage_Hash_SEED64)
	}
	return jsMessage_Hash_hash64String(*value)
}

func jsMessage_Hash_hash64String(value string) int64 {
	var buf [128]byte
	b := jutil.AppendStringBytesUTF8(buf[:0], value)
	return jsMessage_Hash_hash64Range(b, 0, len(b), jsMessage_Hash_SEED64)
}

// jsMessage_Hash_hash64Range hashes byte array to a 64 bit value using the
// supplied seed.
//
// Java: JsMessage.Hash.hash64(byte[], int, int, long)
func jsMessage_Hash_hash64Range(value []byte, offset int, length int, seed uint64) int64 {
	a := jsMessage_Hash_CONSTANT64
	b := a
	c := seed
	keylen := length

	for ; keylen >= 24; keylen, offset = keylen-24, offset+24 {
		a += binary.LittleEndian.Uint64(value[offset:])
		b += binary.LittleEndian.Uint64(value[offset+8:])
		c += binary.LittleEndian.Uint64(value[offset+16:])
		a, b, c = jsMessage_Hash_mix64(a, b, c)
	}

	c += uint64(length)
	tail := value[offset : offset+keylen]
	switch keylen {
	case 23:
		c += uint64(tail[22]) << 56
		fallthrough
	case 22:
		c += uint64(tail[21]) << 48
		fallthrough
	case 21:
		c += uint64(tail[20]) << 40
		fallthrough
	case 20:
		c += uint64(tail[19]) << 32
		fallthrough
	case 19:
		c += uint64(tail[18]) << 24
		fallthrough
	case 18:
		c += uint64(tail[17]) << 16
		fallthrough
	case 17:
		c += uint64(tail[16]) << 8
		// the first byte of c is reserved for the length
		fallthrough
	case 16:
		b += binary.LittleEndian.Uint64(tail[8:])
		a += binary.LittleEndian.Uint64(tail)
	case 15:
		b += uint64(tail[14]) << 48
		fallthrough
	case 14:
		b += uint64(tail[13]) << 40
		fallthrough
	case 13:
		b += uint64(tail[12]) << 32
		fallthrough
	case 12:
		b += uint64(tail[11]) << 24
		fallthrough
	case 11:
		b += uint64(tail[10]) << 16
		fallthrough
	case 10:
		b += uint64(tail[9]) << 8
		fallthrough
	case 9:
		b += uint64(tail[8])
		fallthrough
	case 8:
		a += binary.LittleEndian.Uint64(tail)
	case 7:
		a += uint64(tail[6]) << 48
		fallthrough
	case 6:
		a += uint64(tail[5]) << 40
		fallthrough
	case 5:
		a += uint64(tail[4]) << 32
		fallthrough
	case 4:
		a += uint64(tail[3]) << 24
		fallthrough
	case 3:
		a += uint64(tail[2]) << 16
		fallthrough
	case 2:
		a += uint64(tail[1]) << 8
		fallthrough
	case 1:
		a += uint64(tail[0])
		// case 0: nothing left to add
	}
	_, _, c = jsMessage_Hash_mix64(a, b, c)
	return int64(c)
}

// jsMessage_Hash_mix64 mixes 64-bit words a, b, and c.
//
// Java: JsMessage.Hash.mix64(long, long, long)
func jsMessage_Hash_mix64(a, b, c uint64) (uint64, uint64, uint64) {
	a -= b
	a -= c
	a ^= c >> 43
	b -= c
	b -= a
	b ^= a << 9
	c -= a
	c -= b
	c ^= b >> 8
	a -= b
	a -= c
	a ^= c >> 38
	b -= c
	b -= a
	b ^= a << 23
	c -= a
	c -= b
	c ^= b >> 5
	a -= b
	a -= c
	a ^= c >> 35
	b -= c
	b -= a
	b ^= a << 49
	c -= a
	c -= b
	c ^= b >> 11
	a -= b
	a -= c
	a ^= c >> 12
	b -= c
	b -= a
	b ^= a << 18
	c -= a
	c -= b
	c ^= b >> 22
	return a, b, c
}
