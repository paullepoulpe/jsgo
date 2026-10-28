// Port of com/google/javascript/rhino/JSDocInfo.java @ closure-compiler v20261021

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
 *   Bob Jervis
 *   Google Inc.
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */

package rhino

import (
	"math/bits"
	"slices"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// JSDocInfo is JSDoc information describing JavaScript code. JSDoc is
// represented as a unified object with fields for each JSDoc annotation,
// even though some combinations are incorrect. For instance, if a JSDoc
// describes an enum, it cannot have information about a return type. This
// implementation takes advantage of such incompatibilities to reuse fields
// for multiple purposes, reducing memory consumption.
//
// Constructing JSDocInfo objects is simplified by JSDocInfo.Builder which
// provides early incompatibility detection.
//
// Go-only: language-adapter — Java's Property<T> values are `any`; each
// typed accessor converts with jsdocGet. The property value types are:
// String → string, Integer → int, List<String> → []string,
// ArrayList<JSTypeExpression> → []*JSTypeExpression,
// LinkedHashMap<String, X> → *jutil.LinkedHashMap[string, X],
// ImmutableSet<String> → *jutil.LinkedHashSet[string], and
// ImmutableMap<ImmutableSet<String>, String> →
// *jutil.ObjectLinkedHashMap[*jutil.LinkedHashSet[string], string].
type JSDocInfo struct {
	propertyBits       int64
	propertyKeysBitset int64

	// This value will be `null` if there is no property, a literal value if
	// there is only one, or an `[]any` (`Object[]`) if there are multiple.
	//
	// This is meant to elide eager allocation of single-valued or empty
	// lists.
	propertyValues any
}

// jsdocGet is the unchecked cast `(T) value` of a property value (null is
// the zero T).
func jsdocGet[T any](value any) T {
	return node_cast[T](value)
}

// JSDocParam is an ordered (Name, Type) entry for @param or @template.
type JSDocParam struct {
	Name string
	Type *JSTypeExpression
}

// jsDocKV is an ordered (key, val) string pair for parameter descriptions.
type jsDocKV struct {
	key string
	val string
}

// JSDocSuppression represents a single @suppress annotation's deduplicated
// warning list and optional trailing description.
type JSDocSuppression struct {
	Warnings    []string
	Description string
}

// FormatJSDocSuppressions formats suppressions identically to Java's
// ImmutableMap<ImmutableSet<String>, String>.toString().
func FormatJSDocSuppressions(s []JSDocSuppression) string {
	var b strings.Builder
	b.WriteByte('{')
	for i, sup := range s {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteByte('[')
		for j, w := range sup.Warnings {
			if j > 0 {
				b.WriteString(", ")
			}
			b.WriteString(w)
		}
		b.WriteString("]=")
		b.WriteString(sup.Description)
	}
	b.WriteByte('}')
	return b.String()
}

func sameWarningSet(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for _, w := range a {
		if !slices.Contains(b, w) {
			return false
		}
	}
	return true
}

func formatJsDocParams(s []JSDocParam) string {
	var b strings.Builder
	b.WriteByte('{')
	for i, p := range s {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(p.Name)
		b.WriteByte('=')
		if p.Type == nil {
			b.WriteString("null")
		} else {
			b.WriteString(p.Type.String())
		}
	}
	b.WriteByte('}')
	return b.String()
}

func formatJsDocKVs(s []jsDocKV) string {
	var b strings.Builder
	b.WriteByte('{')
	for i, kv := range s {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(kv.key)
		b.WriteByte('=')
		b.WriteString(kv.val)
	}
	b.WriteByte('}')
	return b.String()
}

func jsDocParamsToMap(s []JSDocParam) *jutil.LinkedHashMap[string, *JSTypeExpression] {
	m := jutil.NewLinkedHashMapWithCapacity[string, *JSTypeExpression](len(s))
	for _, p := range s {
		m.Put(p.Name, p.Type)
	}
	return m
}

// JSDocInfo_Property is a JSDoc property kind.
//
// Go-only: language-adapter — Java's Property<T> and its subclasses
// (MarkerListProperty, TypeProperty, TypeListProperty, TypeMapProperty,
// NodeMapProperty) are one struct; kind selects the subclass's overrides.
// T is erased to `any`.
type JSDocInfo_Property struct {
	name string
	bit  int
	mask int64
	kind jsDocPropKind
}

// jsDocPropKind is the Property subclass.
type jsDocPropKind uint8

const (
	jsDocPropKindPlain        jsDocPropKind = iota // scalar comparable (enum, int, string)
	jsDocPropKindStringList                        // []string
	jsDocPropKindMarkerList                        // []JSDocInfo_Marker
	jsDocPropKindType                              // *JSTypeExpression
	jsDocPropKindTypeList                          // []*JSTypeExpression
	jsDocPropKindTypeMap                           // []JSDocParam
	jsDocPropKindStringMap                         // []jsDocKV
	jsDocPropKindSuppressions                      // []JSDocSuppression
	jsDocPropKindStringSet                         // *jutil.LinkedHashSet[string]
	jsDocPropKindNodeMap                           // *jutil.LinkedHashMap[string, Node]
)

// Java: JSDocInfo.Property.bitCounter
var JSDocInfo_Property_bitCounter = 0

// Java: JSDocInfo.Property.values
var JSDocInfo_Property_values = make([]*JSDocInfo_Property, 64)

// jsDocInfo_typePropertiesMask has the bits of the properties holding type
// expressions (the type, type list and type map kinds).
var jsDocInfo_typePropertiesMask int64

// Java: JSDocInfo.Property(String), and the subclass constructors
func newJSDocInfo_Property(name string, kind jsDocPropKind) *JSDocInfo_Property {
	p := &JSDocInfo_Property{name: name, bit: JSDocInfo_Property_bitCounter, kind: kind}
	JSDocInfo_Property_bitCounter++
	p.mask = int64(1) << p.bit
	if p.bit > 63 {
		panic(jutil.NewAssertionError("Too many Properties"))
	}
	JSDocInfo_Property_values[p.bit] = p
	switch kind {
	case jsDocPropKindType, jsDocPropKindTypeList, jsDocPropKindTypeMap:
		jsDocInfo_typePropertiesMask |= p.mask
	}
	return p
}

func (p *JSDocInfo_Property) getRaw(info *JSDocInfo) any {
	if (info.propertyKeysBitset & p.mask) == 0 {
		return nil
	}
	return info.GetPropertyValueByIndex(bits.OnesCount64(uint64(info.propertyKeysBitset & (p.mask - 1))))
}

// Get returns the property value of info, or null; the cast to T is unsafe
// but guaranteed by builder.
//
// Java: JSDocInfo.Property.get(JSDocInfo)
func (p *JSDocInfo_Property) Get(info *JSDocInfo) any {
	v := p.getRaw(info)
	if v == nil {
		return nil
	}
	if p.kind == jsDocPropKindTypeMap {
		return jsDocParamsToMap(jsdocGet[[]JSDocParam](v))
	}
	return v
}

func isNilPropValue(kind jsDocPropKind, v any) bool {
	if v == nil {
		return true
	}
	switch kind {
	case jsDocPropKindType:
		return jsdocGet[*JSTypeExpression](v) == nil
	case jsDocPropKindStringSet:
		return jsdocGet[*jutil.LinkedHashSet[string]](v) == nil
	case jsDocPropKindNodeMap:
		return jsdocGet[*jutil.LinkedHashMap[string, Node]](v) == nil
	}
	return false
}

// Clone copies a value for JSDocInfo.toBuilder. NOTE: documentation props
// are never cloned.
//
// Java: JSDocInfo.Property.clone(T, TypeTransform), and its overrides
func (p *JSDocInfo_Property) Clone(arg any, transform JSDocInfo_TypeTransform) any {
	switch p.kind {
	case jsDocPropKindType:
		e := jsdocGet[*JSTypeExpression](arg)
		if e != nil && transform != nil {
			return transform(e)
		}
		return arg
	case jsDocPropKindTypeList:
		out := slices.Clone(jsdocGet[[]*JSTypeExpression](arg))
		if transform != nil {
			for i, elem := range out {
				if elem != nil {
					out[i] = transform(elem)
				}
			}
		}
		return out
	case jsDocPropKindStringList:
		return slices.Clip(jsdocGet[[]string](arg))
	case jsDocPropKindMarkerList:
		return slices.Clip(jsdocGet[[]JSDocInfo_Marker](arg))
	case jsDocPropKindTypeMap:
		out := slices.Clone(jsdocGet[[]JSDocParam](arg))
		if transform != nil {
			for i := range out {
				if out[i].Type != nil {
					out[i].Type = transform(out[i].Type)
				}
			}
		}
		return out
	case jsDocPropKindStringMap:
		return slices.Clip(jsdocGet[[]jsDocKV](arg))
	case jsDocPropKindSuppressions:
		return slices.Clip(jsdocGet[[]JSDocSuppression](arg))
	case jsDocPropKindNodeMap:
		return jutil.NewLinkedHashMapFrom(jsdocGet[*jutil.LinkedHashMap[string, Node]](arg))
	}
	return arg
}

// Java: JSDocInfo.Property.isDefault(T)
func (p *JSDocInfo_Property) IsDefault(value any) bool {
	if value == nil {
		return true
	}
	switch p.kind {
	case jsDocPropKindPlain:
		if v, ok := value.(JSDocInfo_Visibility); ok {
			return v == JSDocInfo_Visibility_INHERITED
		}
		return false
	case jsDocPropKindStringList:
		return len(jsdocGet[[]string](value)) == 0
	case jsDocPropKindMarkerList:
		return len(jsdocGet[[]JSDocInfo_Marker](value)) == 0
	case jsDocPropKindType:
		return jsdocGet[*JSTypeExpression](value) == nil
	case jsDocPropKindTypeList:
		return len(jsdocGet[[]*JSTypeExpression](value)) == 0
	case jsDocPropKindTypeMap:
		return len(jsdocGet[[]JSDocParam](value)) == 0
	case jsDocPropKindStringMap:
		return len(jsdocGet[[]jsDocKV](value)) == 0
	case jsDocPropKindSuppressions:
		return len(jsdocGet[[]JSDocSuppression](value)) == 0
	case jsDocPropKindStringSet:
		s := jsdocGet[*jutil.LinkedHashSet[string]](value)
		return s == nil || s.IsEmpty()
	case jsDocPropKindNodeMap:
		m := jsdocGet[*jutil.LinkedHashMap[string, Node]](value)
		return m == nil || m.IsEmpty()
	}
	return false
}

// Java: JSDocInfo.Property.equalValues(T, T), and its overrides
func (p *JSDocInfo_Property) EqualValues(left any, right any) bool {
	switch p.kind {
	case jsDocPropKindPlain:
		return left == right
	case jsDocPropKindStringList:
		return slices.Equal(jsdocGet[[]string](left), jsdocGet[[]string](right))
	case jsDocPropKindMarkerList:
		a := jsdocGet[[]JSDocInfo_Marker](left)
		b := jsdocGet[[]JSDocInfo_Marker](right)
		if len(a) != len(b) {
			return false
		}
		for i := range a {
			if !(&a[i]).IsEquivalentTo(&b[i]) {
				return false
			}
		}
		return true
	case jsDocPropKindType:
		return jsdocGet[*JSTypeExpression](left).IsEquivalentTo(jsdocGet[*JSTypeExpression](right))
	case jsDocPropKindTypeList:
		l := jsdocGet[[]*JSTypeExpression](left)
		r := jsdocGet[[]*JSTypeExpression](right)
		if len(l) != len(r) {
			return false
		}
		for i, le := range l {
			if !le.IsEquivalentTo(r[i]) {
				return false
			}
		}
		return true
	case jsDocPropKindTypeMap:
		l := jsdocGet[[]JSDocParam](left)
		r := jsdocGet[[]JSDocParam](right)
		if len(l) != len(r) {
			return false
		}
		for _, lp := range l {
			found := false
			for _, rp := range r {
				if rp.Name == lp.Name {
					if !jsDocInfo_TypeMapProperty_AreEquivalent(lp.Type, rp.Type) {
						return false
					}
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
		return true
	case jsDocPropKindStringMap:
		l := jsdocGet[[]jsDocKV](left)
		r := jsdocGet[[]jsDocKV](right)
		if len(l) != len(r) {
			return false
		}
		for _, lkv := range l {
			found := false
			for _, rkv := range r {
				if rkv.key == lkv.key {
					if rkv.val != lkv.val {
						return false
					}
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
		return true
	case jsDocPropKindSuppressions:
		l := jsdocGet[[]JSDocSuppression](left)
		r := jsdocGet[[]JSDocSuppression](right)
		if len(l) != len(r) {
			return false
		}
		for _, le := range l {
			found := false
			for _, re := range r {
				if sameWarningSet(le.Warnings, re.Warnings) {
					if le.Description != re.Description {
						return false
					}
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
		return true
	case jsDocPropKindStringSet:
		return jsdocGet[*jutil.LinkedHashSet[string]](left).Equals(jsdocGet[*jutil.LinkedHashSet[string]](right))
	case jsDocPropKindNodeMap:
		l := jsdocGet[*jutil.LinkedHashMap[string, Node]](left)
		r := jsdocGet[*jutil.LinkedHashMap[string, Node]](right)
		if l.Size() != r.Size() {
			return false
		}
		for k, v := range l.All() {
			rightValue := r.Get(k)
			if rightValue.IsNil() || !v.IsEquivalentTo(rightValue) {
				return false
			}
		}
		return true
	}
	return left == right
}

// Java: JSDocInfo.TypeMapProperty.areEquivalent(JSTypeExpression, JSTypeExpression)
func jsDocInfo_TypeMapProperty_AreEquivalent(expr1 *JSTypeExpression, expr2 *JSTypeExpression) bool {
	return expr1 == expr2 || (expr1 != nil && expr1.IsEquivalentTo(expr2))
}

// CompareTo orders properties by bit. NOTE: These need to be comparable so
// that we can iterate over them to build the list.
//
// Java: JSDocInfo.Property.compareTo(Property<Object>)
func (p *JSDocInfo_Property) CompareTo(that *JSDocInfo_Property) int {
	return p.bit - that.bit
}

// GetUnmodifiable returns the elements of a type list property (or an empty
// list).
//
// Go-only: language-adapter — Java returns an unmodifiable view of the
// list; this returns its current elements.
//
// Java: JSDocInfo.TypeListProperty.getUnmodifiable(JSDocInfo)
func (p *JSDocInfo_Property) GetUnmodifiable(info *JSDocInfo) []*JSTypeExpression {
	if value := jsdocGet[[]*JSTypeExpression](p.getRaw(info)); value != nil {
		return value
	}
	return []*JSTypeExpression{}
}

// JSDocInfo_Bit enumerates the JSDoc boolean annotations.
//
// Java: JSDocInfo.Bit
type JSDocInfo_Bit int32

const (
	// Whether the type annotation was inlined.
	JSDocInfo_Bit_INLINE_TYPE JSDocInfo_Bit = iota
	// Whether to include documentation.
	JSDocInfo_Bit_INCLUDE_DOCUMENTATION

	// The following are all straightforward annotations.
	JSDocInfo_Bit_CONST
	JSDocInfo_Bit_CONSTRUCTOR
	JSDocInfo_Bit_DEFINE
	JSDocInfo_Bit_TYPE_SUMMARY
	JSDocInfo_Bit_FINAL
	JSDocInfo_Bit_OVERRIDE

	JSDocInfo_Bit_DEPRECATED
	JSDocInfo_Bit_INTERFACE
	JSDocInfo_Bit_EXPORT
	JSDocInfo_Bit_ENHANCED_NAMESPACE
	JSDocInfo_Bit_NOINLINE
	JSDocInfo_Bit_FILEOVERVIEW
	JSDocInfo_Bit_IMPLICITCAST
	JSDocInfo_Bit_NOSIDEEFFECTS
	JSDocInfo_Bit_EXTERNS
	JSDocInfo_Bit_NOCOMPILE
	JSDocInfo_Bit_NODTS
	JSDocInfo_Bit_UNRESTRICTED
	JSDocInfo_Bit_USED_VIA_DOT_CONSTRUCTOR
	JSDocInfo_Bit_STRUCT
	JSDocInfo_Bit_DICT
	JSDocInfo_Bit_NOCOLLAPSE
	JSDocInfo_Bit_RECORD
	JSDocInfo_Bit_ABSTRACT
	JSDocInfo_Bit_PURE_OR_BREAK_MY_CODE
	JSDocInfo_Bit_COLLAPSIBLE_OR_BREAK_MY_CODE
	JSDocInfo_Bit_NOCOVERAGE
	JSDocInfo_Bit_REQUIRE_INLINING
	JSDocInfo_Bit_ENCOURAGE_INLINING

	JSDocInfo_Bit_NG_INJECT
	JSDocInfo_Bit_WIZ_ACTION
	JSDocInfo_Bit_POLYMER_BEHAVIOR
	JSDocInfo_Bit_POLYMER
	JSDocInfo_Bit_CUSTOM_ELEMENT
	JSDocInfo_Bit_MIXIN_CLASS
	JSDocInfo_Bit_MIXIN_FUNCTION
	JSDocInfo_Bit_SASS_GENERATED_CSS_TS
	JSDocInfo_Bit_CLOSURE_UNAWARE_CODE

	// `@provideGoog` only appears in base.js
	JSDocInfo_Bit_PROVIDE_GOOG
	JSDocInfo_Bit_PROVIDE_ALREADY_PROVIDED

	JSDocInfo_Bit_WIZ_CALLBACK
	JSDocInfo_Bit_LOG_TYPE_IN_COMPILER
	JSDocInfo_Bit_CLOSURE_UNAWARE_CONFIG
)

var jsdocInfo_Bit_names = [...]string{
	"INLINE_TYPE", "INCLUDE_DOCUMENTATION", "CONST", "CONSTRUCTOR", "DEFINE", "TYPE_SUMMARY",
	"FINAL", "OVERRIDE", "DEPRECATED", "INTERFACE", "EXPORT", "ENHANCED_NAMESPACE", "NOINLINE",
	"FILEOVERVIEW", "IMPLICITCAST", "NOSIDEEFFECTS", "EXTERNS", "NOCOMPILE", "NODTS",
	"UNRESTRICTED", "USED_VIA_DOT_CONSTRUCTOR", "STRUCT", "DICT", "NOCOLLAPSE", "RECORD",
	"ABSTRACT", "PURE_OR_BREAK_MY_CODE", "COLLAPSIBLE_OR_BREAK_MY_CODE", "NOCOVERAGE",
	"REQUIRE_INLINING", "ENCOURAGE_INLINING", "NG_INJECT", "WIZ_ACTION", "POLYMER_BEHAVIOR",
	"POLYMER", "CUSTOM_ELEMENT", "MIXIN_CLASS", "MIXIN_FUNCTION", "SASS_GENERATED_CSS_TS",
	"CLOSURE_UNAWARE_CODE", "PROVIDE_GOOG", "PROVIDE_ALREADY_PROVIDED", "WIZ_CALLBACK",
	"LOG_TYPE_IN_COMPILER", "CLOSURE_UNAWARE_CONFIG",
}

// JSDocInfo_Bit_values is Bit.values().
func JSDocInfo_Bit_values() []JSDocInfo_Bit {
	out := make([]JSDocInfo_Bit, len(jsdocInfo_Bit_names))
	for i := range out {
		out[i] = JSDocInfo_Bit(i)
	}
	return out
}

// Name is Enum.name().
func (b JSDocInfo_Bit) Name() string { return jsdocInfo_Bit_names[b] }

// String is Enum.toString().
func (b JSDocInfo_Bit) String() string { return jsdocInfo_Bit_names[b] }

// The `name` field of each Bit, computed at init without slice allocation.
var jsdocInfo_Bit_name = func() [len(jsdocInfo_Bit_names)]string {
	var names [len(jsdocInfo_Bit_names)]string
	if len(jsdocInfo_Bit_names) > 64 {
		panic(jutil.NewAssertionError("Too many Bits"))
	}
	for i, raw := range jsdocInfo_Bit_names {
		names[i] = jutil.CaseFormat_UpperUnderscoreToLowerCamel(raw)
	}
	return names
}()

// Java: JSDocInfo.Bit.name
func (b JSDocInfo_Bit) BitName() string { return jsdocInfo_Bit_name[b] }

// Java: JSDocInfo.Bit.mask
func (b JSDocInfo_Bit) Mask() int64 { return int64(1) << b }

// JSDocInfo_TypeTransform transforms a type expression.
//
// Go-only: language-adapter — a functional interface, so a Go func type.
//
// Java: JSDocInfo.TypeTransform.apply(JSTypeExpression)
type JSDocInfo_TypeTransform func(arg *JSTypeExpression) *JSTypeExpression

func newJSDocInfoFromBuilder(bits_ int64, propMask int64, entries []any) *JSDocInfo {
	var keys int64 = 0
	var valsBuf [8]any
	values := valsBuf[:0]
	mask := propMask
	for i := 0; mask != 0; i++ {
		low := bits.TrailingZeros64(uint64(mask))
		propBit := int64(1) << low
		mask &^= propBit
		prop := JSDocInfo_Property_values[low]
		value := entries[i]
		if !prop.IsDefault(value) {
			keys |= propBit
			_, isArray := value.([]any)
			jutil.CheckState(!isArray && !isNilPropValue(prop.kind, value))
			values = append(values, value)
		}
	}
	return &JSDocInfo{
		propertyBits:       bits_,
		propertyKeysBitset: keys,
		propertyValues:     JSDocInfo_PackPropertyValues(values),
	}
}

// Java: JSDocInfo.checkBit(Bit)
func (i *JSDocInfo) CheckBit(bit JSDocInfo_Bit) bool {
	return (i.propertyBits & bit.Mask()) != 0
}

// JSDocInfo_Visibility enumerates visibility categories. The ordinal can be
// used as a numerical indicator of privacy, where 0 is the most private.
// This means that compareTo can be used to determine if a visibility is more
// permissive than another.
//
// Java: JSDocInfo.Visibility
type JSDocInfo_Visibility int32

const (
	JSDocInfo_Visibility_PRIVATE JSDocInfo_Visibility = iota
	JSDocInfo_Visibility_PACKAGE
	JSDocInfo_Visibility_PROTECTED
	JSDocInfo_Visibility_PUBLIC

	// If visibility is not specified, we just assume that visibility is
	// inherited from the super class. This should never be included in the
	// map of an actual built JSDocInfo (though it may be in the builder), but
	// will be returned by default if none is specified.
	JSDocInfo_Visibility_INHERITED
)

var jsdocInfo_Visibility_names = [...]string{"PRIVATE", "PACKAGE", "PROTECTED", "PUBLIC", "INHERITED"}

// String is Enum.name().
func (v JSDocInfo_Visibility) String() string { return jsdocInfo_Visibility_names[v] }

// Ordinal is Enum.ordinal().
func (v JSDocInfo_Visibility) Ordinal() int { return int(v) }

// JSDocInfo_IdGenerator enumerates the @idGenerator kinds.
//
// Java: JSDocInfo.IdGenerator
type JSDocInfo_IdGenerator int32

const (
	JSDocInfo_IdGenerator_XID JSDocInfo_IdGenerator = iota
	JSDocInfo_IdGenerator_CONSISTENT
	JSDocInfo_IdGenerator_UNIQUE
	JSDocInfo_IdGenerator_STABLE
	JSDocInfo_IdGenerator_MAPPED
)

var jsdocInfo_IdGenerator_names = [...]string{"XID", "CONSISTENT", "UNIQUE", "STABLE", "MAPPED"}

// String is Enum.name().
func (g JSDocInfo_IdGenerator) String() string { return jsdocInfo_IdGenerator_names[g] }

// JSDocInfo_PerFileClosureUnawareMode configures whether any @closureUnaware
// blocks in a file will go through any compiler optimization passes or
// transpilation.
//
// SIMPLE: the file goes through a predefined set of safe optimizations &
// transpilation
//
// WHITESPACE: the file goes through no optimizations or transpilation,
// although code printing does strip unnecessary whitespace and comments.
//
// UNSPECIFIED: indicates the annotation was just `@closureUnaware` with no
// mode specified. Currently this is equivalent to SIMPLE.
//
// Java: JSDocInfo.PerFileClosureUnawareMode
type JSDocInfo_PerFileClosureUnawareMode int32

const (
	JSDocInfo_PerFileClosureUnawareMode_UNSPECIFIED JSDocInfo_PerFileClosureUnawareMode = iota
	JSDocInfo_PerFileClosureUnawareMode_SIMPLE
	JSDocInfo_PerFileClosureUnawareMode_WHITESPACE
)

var jsdocInfo_PerFileClosureUnawareMode_names = [...]string{"UNSPECIFIED", "SIMPLE", "WHITESPACE"}

// JSDocInfo_PerFileClosureUnawareMode_values is
// PerFileClosureUnawareMode.values().
func JSDocInfo_PerFileClosureUnawareMode_values() []JSDocInfo_PerFileClosureUnawareMode {
	return []JSDocInfo_PerFileClosureUnawareMode{
		JSDocInfo_PerFileClosureUnawareMode_UNSPECIFIED, JSDocInfo_PerFileClosureUnawareMode_SIMPLE,
		JSDocInfo_PerFileClosureUnawareMode_WHITESPACE,
	}
}

// String is Enum.name().
func (m JSDocInfo_PerFileClosureUnawareMode) String() string {
	return jsdocInfo_PerFileClosureUnawareMode_names[m]
}

// The properties, in Java declaration order (which assigns their bits).
var (
	JSDocInfo_VISIBILITY                = newJSDocInfo_Property("visibility", jsDocPropKindPlain)
	JSDocInfo_TYPE                      = newJSDocInfo_Property("type", jsDocPropKindType)
	JSDocInfo_RETURN_TYPE               = newJSDocInfo_Property("returnType", jsDocPropKindType)
	JSDocInfo_ENUM_PARAMETER_TYPE       = newJSDocInfo_Property("enumParameterType", jsDocPropKindType)
	JSDocInfo_TYPEDEF_TYPE              = newJSDocInfo_Property("typedefType", jsDocPropKindType)
	JSDocInfo_THIS_TYPE                 = newJSDocInfo_Property("thisType", jsDocPropKindType)
	JSDocInfo_ORIGINAL_COMMENT_POSITION = newJSDocInfo_Property("originalCommentPosition", jsDocPropKindPlain)
	JSDocInfo_ID_GENERATOR              = newJSDocInfo_Property("idGenerator", jsDocPropKindPlain)
	JSDocInfo_BASE_TYPE                 = newJSDocInfo_Property("baseType", jsDocPropKindType)
	JSDocInfo_EXTENDED_INTERFACES       = newJSDocInfo_Property("extendedInterfaces", jsDocPropKindTypeList)
	JSDocInfo_IMPLEMENTED_INTERFACES    = newJSDocInfo_Property("extendedInterfaces", jsDocPropKindTypeList)
	JSDocInfo_PARAMETERS                = newJSDocInfo_Property("parameters", jsDocPropKindTypeMap)
	JSDocInfo_THROWS_ANNOTATIONS        = newJSDocInfo_Property("throws", jsDocPropKindStringList)
	JSDocInfo_TEMPLATE_TYPE_NAMES       = newJSDocInfo_Property("templateTypeNames", jsDocPropKindTypeMap)
	JSDocInfo_TYPE_TRANSFORMATIONS      = newJSDocInfo_Property("typeTransformations", jsDocPropKindNodeMap)

	JSDocInfo_DESCRIPTION          = newJSDocInfo_Property("description", jsDocPropKindPlain)
	JSDocInfo_MEANING              = newJSDocInfo_Property("meaning", jsDocPropKindPlain)
	JSDocInfo_ALTERNATE_MESSAGE_ID = newJSDocInfo_Property("alternateMessageId", jsDocPropKindPlain)
	JSDocInfo_DEPRECATION_REASON   = newJSDocInfo_Property("deprecationReason", jsDocPropKindPlain)
	JSDocInfo_LICENSE              = newJSDocInfo_Property("license", jsDocPropKindPlain)

	// For example, `@suppress {first, second} Some description applying to
	// the set (of both first and second).`
	// For example, `@suppress {first} Some description applying to only to
	// first.`
	JSDocInfo_SUPPRESSIONS         = newJSDocInfo_Property("suppressions", jsDocPropKindSuppressions)
	JSDocInfo_MODIFIES             = newJSDocInfo_Property("modifies", jsDocPropKindStringSet)
	JSDocInfo_LENDS_NAME           = newJSDocInfo_Property("lendsName", jsDocPropKindType)
	JSDocInfo_CLOSURE_PRIMITIVE_ID = newJSDocInfo_Property("closurePrimitiveId", jsDocPropKindPlain)

	// NOTE: The following properties are "documentation properties", which do
	// _not_ get deeply copied when JSDocInfo is cloned (i.e. none of these
	// Properties override Property#clone).
	JSDocInfo_SOURCE_COMMENT               = newJSDocInfo_Property("sourceComment", jsDocPropKindPlain)
	JSDocInfo_MARKERS                      = newJSDocInfo_Property("markers", jsDocPropKindMarkerList)
	JSDocInfo_PARAMETER_DESCRIPTIONS       = newJSDocInfo_Property("parameterDescriptions", jsDocPropKindStringMap)
	JSDocInfo_BLOCK_DESCRIPTION            = newJSDocInfo_Property("blockDescription", jsDocPropKindPlain)
	JSDocInfo_FILEOVERVIEW_DESCRIPTION     = newJSDocInfo_Property("fileoverviewDescription", jsDocPropKindPlain)
	JSDocInfo_RETURN_DESCRIPTION           = newJSDocInfo_Property("returnDescription", jsDocPropKindPlain)
	JSDocInfo_ENHANCED_NAMESPACE           = newJSDocInfo_Property("enhance", jsDocPropKindPlain)
	JSDocInfo_MODS                         = newJSDocInfo_Property("mods", jsDocPropKindPlain)
	JSDocInfo_TS_TYPES                     = newJSDocInfo_Property("tsType", jsDocPropKindStringList)
	JSDocInfo_CLOSURE_UNAWARE_CONFIG_VALUE = newJSDocInfo_Property("perFileClosureUnawareMode", jsDocPropKindPlain)

	JSDocInfo_AUTHORS = newJSDocInfo_Property("authors", jsDocPropKindStringList)
	JSDocInfo_SEES    = newJSDocInfo_Property("sees", jsDocPropKindStringList)
)

// JSDocInfo_Builder_ creates a new JSDocInfo.Builder object.
//
// Go-only: language-adapter — named with a trailing underscore because the
// nested type JSDocInfo_Builder has the same Go name.
//
// Java: JSDocInfo.builder()
func JSDocInfo_Builder_() *JSDocInfo_Builder {
	return NewJSDocInfo_Builder()
}

// Java: JSDocInfo.toBuilder()
func (i *JSDocInfo) ToBuilder() *JSDocInfo_Builder {
	return i.ToBuilderWithTransform(nil)
}

// Java: JSDocInfo.toBuilder(TypeTransform)
func (i *JSDocInfo) ToBuilderWithTransform(transform JSDocInfo_TypeTransform) *JSDocInfo_Builder {
	builder := NewJSDocInfo_Builder()
	// inline type isn't copied...?
	builder.bits = i.propertyBits &^ JSDocInfo_Bit_INLINE_TYPE.Mask()
	builder.populated = true
	bits_ := i.propertyKeysBitset
	builder.propMask = bits_
	for index := 0; bits_ != 0; index++ {
		low := bits.TrailingZeros64(uint64(bits_))
		bits_ &^= int64(1) << low
		prop := JSDocInfo_Property_values[low]
		builder.entries = append(builder.entries, prop.Clone(i.GetPropertyValueByIndex(index), transform))
	}
	return builder
}

// Java: JSDocInfo.clone()
func (i *JSDocInfo) Clone() *JSDocInfo {
	return i.CloneWithTypeNodes(false)
}

// CloneAndReplaceTypeNames clones this JSDoc but replaces the given names in
// any type related annotation with unknown type.
//
// Java: JSDocInfo.cloneAndReplaceTypeNames(Set<String>)
func (i *JSDocInfo) CloneAndReplaceTypeNames(names *jutil.LinkedHashSet[string]) *JSDocInfo {
	return i.ToBuilderWithTransform(func(type_ *JSTypeExpression) *JSTypeExpression {
		return type_.ReplaceNamesWithUnknownType(names)
	}).Build()
}

// Java: JSDocInfo.clone(boolean)
func (i *JSDocInfo) CloneWithTypeNodes(cloneTypeNodes bool) *JSDocInfo {
	if cloneTypeNodes {
		return i.ToBuilderWithTransform((*JSTypeExpression).Copy).Build()
	}
	return i.ToBuilder().Build()
}

// Java: JSDocInfo.areEquivalent(JSDocInfo, JSDocInfo)
func JSDocInfo_AreEquivalent(jsDoc1 *JSDocInfo, jsDoc2 *JSDocInfo) bool {
	if (jsDoc1 == nil) != (jsDoc2 == nil) {
		return false
	} else if jsDoc1 == nil {
		return true
	}
	if ((jsDoc1.propertyBits ^ jsDoc2.propertyBits) &^ JSDocInfo_EQUIVALENCE_IGNORED_BITS) != 0 {
		return false
	}
	if jsDoc1.propertyKeysBitset != jsDoc2.propertyKeysBitset {
		return false
	}
	bits_ := jsDoc1.propertyKeysBitset
	index := 0
	for bits_ > 0 {
		low := bits.TrailingZeros64(uint64(bits_))
		bits_ &^= int64(1) << low
		prop := JSDocInfo_Property_values[low]
		a := jsDoc1.GetPropertyValueByIndex(index)
		b := jsDoc2.GetPropertyValueByIndex(index)
		index++
		if !prop.EqualValues(a, b) {
			return false
		}
	}
	return true
}

// NOTE: includeDocumentation and inlineType are not part of equivalence.
//
// Java: JSDocInfo.EQUIVALENCE_IGNORED_BITS
var JSDocInfo_EQUIVALENCE_IGNORED_BITS = JSDocInfo_Bit_INCLUDE_DOCUMENTATION.Mask() | JSDocInfo_Bit_INLINE_TYPE.Mask()

// Java: JSDocInfo.isDocumentationIncluded()
func (i *JSDocInfo) IsDocumentationIncluded() bool {
	return i.CheckBit(JSDocInfo_Bit_INCLUDE_DOCUMENTATION)
}

// IsAnyIdGenerator returns whether any `@idGenerator` or `idGenerator {?}`
// annotation is present.
//
// Java: JSDocInfo.isAnyIdGenerator()
func (i *JSDocInfo) IsAnyIdGenerator() bool {
	return (i.propertyKeysBitset & JSDocInfo_ID_GENERATOR.mask) != 0
}

func (i *JSDocInfo) idGeneratorIs(g JSDocInfo_IdGenerator) bool {
	v := JSDocInfo_ID_GENERATOR.getRaw(i)
	return v != nil && jsdocGet[JSDocInfo_IdGenerator](v) == g
}

// IsConsistentIdGenerator returns whether the `@idGenerator {consistent}` is
// present on this JSDocInfo.
//
// Java: JSDocInfo.isConsistentIdGenerator()
func (i *JSDocInfo) IsConsistentIdGenerator() bool {
	return i.idGeneratorIs(JSDocInfo_IdGenerator_CONSISTENT)
}

// IsStableIdGenerator returns whether the `@idGenerator {stable}` is present
// on this JSDocInfo.
//
// Java: JSDocInfo.isStableIdGenerator()
func (i *JSDocInfo) IsStableIdGenerator() bool {
	return i.idGeneratorIs(JSDocInfo_IdGenerator_STABLE)
}

// IsXidGenerator returns whether the `@idGenerator {xid}` is present on this
// JSDocInfo.
//
// Java: JSDocInfo.isXidGenerator()
func (i *JSDocInfo) IsXidGenerator() bool {
	return i.idGeneratorIs(JSDocInfo_IdGenerator_XID)
}

// IsMappedIdGenerator returns whether the `@idGenerator {mapped}` is present
// on this JSDocInfo.
//
// Java: JSDocInfo.isMappedIdGenerator()
func (i *JSDocInfo) IsMappedIdGenerator() bool {
	return i.idGeneratorIs(JSDocInfo_IdGenerator_MAPPED)
}

// IsIdGenerator returns whether the `@idGenerator` is present on this
// JSDocInfo.
//
// Java: JSDocInfo.isIdGenerator()
func (i *JSDocInfo) IsIdGenerator() bool {
	return i.idGeneratorIs(JSDocInfo_IdGenerator_UNIQUE)
}

// IsConstant returns whether this JSDocInfo implies that annotated value is
// constant.
//
// Java: JSDocInfo.isConstant()
func (i *JSDocInfo) IsConstant() bool {
	// @desc is used with goog.getMsg to define messages to be translated, and
	// thus must be @const in order for translation to work correctly.
	return (i.propertyBits&(JSDocInfo_Bit_CONST.Mask()|JSDocInfo_Bit_DEFINE.Mask()|JSDocInfo_Bit_FINAL.Mask())) != 0 ||
		(i.propertyKeysBitset&JSDocInfo_DESCRIPTION.mask) != 0
}

// HasConstAnnotation returns whether the `@const` annotation is present on
// this JSDocInfo.
//
// Java: JSDocInfo.hasConstAnnotation()
func (i *JSDocInfo) HasConstAnnotation() bool {
	return i.CheckBit(JSDocInfo_Bit_CONST)
}

// IsFinal returns whether the `@final` annotation is present on this
// JSDocInfo.
//
// Java: JSDocInfo.isFinal()
func (i *JSDocInfo) IsFinal() bool {
	return i.CheckBit(JSDocInfo_Bit_FINAL)
}

// IsConstructor returns whether the `@constructor` annotation is present on
// this JSDocInfo.
//
// Java: JSDocInfo.isConstructor()
func (i *JSDocInfo) IsConstructor() bool {
	return i.CheckBit(JSDocInfo_Bit_CONSTRUCTOR)
}

// IsAbstract returns whether the `@abstract` annotation is present on this
// JSDocInfo.
//
// Java: JSDocInfo.isAbstract()
func (i *JSDocInfo) IsAbstract() bool {
	return i.CheckBit(JSDocInfo_Bit_ABSTRACT)
}

// UsesImplicitMatch returns whether the `@record` annotation is present on
// this JSDocInfo.
//
// Java: JSDocInfo.usesImplicitMatch()
func (i *JSDocInfo) UsesImplicitMatch() bool {
	return i.CheckBit(JSDocInfo_Bit_RECORD)
}

// MakesUnrestricted returns whether the `@unrestricted` annotation is present
// on this JSDocInfo.
//
// Java: JSDocInfo.makesUnrestricted()
func (i *JSDocInfo) MakesUnrestricted() bool {
	return i.CheckBit(JSDocInfo_Bit_UNRESTRICTED)
}

// MakesStructs returns whether the `@struct` annotation is present on this
// JSDocInfo.
//
// Java: JSDocInfo.makesStructs()
func (i *JSDocInfo) MakesStructs() bool {
	return i.CheckBit(JSDocInfo_Bit_STRUCT)
}

// MakesDicts returns whether the `@dict` annotation is present on this
// JSDocInfo.
//
// Java: JSDocInfo.makesDicts()
func (i *JSDocInfo) MakesDicts() bool {
	return i.CheckBit(JSDocInfo_Bit_DICT)
}

// IsDefine returns whether the `@define` annotation is present on this
// JSDocInfo. If this annotation is present, then the #getType() method will
// retrieve the define type.
//
// Java: JSDocInfo.isDefine()
func (i *JSDocInfo) IsDefine() bool {
	return i.CheckBit(JSDocInfo_Bit_DEFINE)
}

// IsOverride returns whether the `@override` annotation is present on this
// JSDocInfo.
//
// Java: JSDocInfo.isOverride()
func (i *JSDocInfo) IsOverride() bool {
	return i.CheckBit(JSDocInfo_Bit_OVERRIDE)
}

// IsDeprecated returns whether the `@deprecated` annotation is present on
// this JSDocInfo.
//
// Java: JSDocInfo.isDeprecated()
func (i *JSDocInfo) IsDeprecated() bool {
	return i.CheckBit(JSDocInfo_Bit_DEPRECATED)
}

// IsInterface returns whether the `@interface` annotation is present on this
// JSDocInfo.
//
// Java: JSDocInfo.isInterface()
func (i *JSDocInfo) IsInterface() bool {
	return (i.propertyBits & (JSDocInfo_Bit_INTERFACE.Mask() | JSDocInfo_Bit_RECORD.Mask())) != 0
}

// Java: JSDocInfo.isConstructorOrInterface()
func (i *JSDocInfo) IsConstructorOrInterface() bool {
	return (i.propertyBits & (JSDocInfo_Bit_CONSTRUCTOR.Mask() | JSDocInfo_Bit_INTERFACE.Mask() | JSDocInfo_Bit_RECORD.Mask())) != 0
}

// IsExport returns whether the `@export` annotation is present on this
// JSDocInfo.
//
// Java: JSDocInfo.isExport()
func (i *JSDocInfo) IsExport() bool {
	return i.CheckBit(JSDocInfo_Bit_EXPORT)
}

// IsImplicitCast returns whether the `@implicitCast` annotation is present on
// this JSDocInfo.
//
// Java: JSDocInfo.isImplicitCast()
func (i *JSDocInfo) IsImplicitCast() bool {
	return i.CheckBit(JSDocInfo_Bit_IMPLICITCAST)
}

// IsNoSideEffects returns whether the `@nosideeffects` annotation is present
// on this JSDocInfo.
//
// Java: JSDocInfo.isNoSideEffects()
func (i *JSDocInfo) IsNoSideEffects() bool {
	return i.CheckBit(JSDocInfo_Bit_NOSIDEEFFECTS)
}

// IsExterns returns whether the `@externs` annotation is present on this
// JSDocInfo.
//
// Java: JSDocInfo.isExterns()
func (i *JSDocInfo) IsExterns() bool {
	return i.CheckBit(JSDocInfo_Bit_EXTERNS)
}

// IsNoCoverage returns whether the `@nocoverage` annotation is present on
// this JSDocInfo.
//
// Java: JSDocInfo.isNoCoverage()
func (i *JSDocInfo) IsNoCoverage() bool {
	return i.CheckBit(JSDocInfo_Bit_NOCOVERAGE)
}

// IsTypeSummary returns whether the `@typeSummary` annotation is present on
// this JSDocInfo.
//
// Java: JSDocInfo.isTypeSummary()
func (i *JSDocInfo) IsTypeSummary() bool {
	return i.CheckBit(JSDocInfo_Bit_TYPE_SUMMARY)
}

// IsNoCompile returns whether the `@nocompile` annotation is present on this
// JSDocInfo.
//
// Java: JSDocInfo.isNoCompile()
func (i *JSDocInfo) IsNoCompile() bool {
	return i.CheckBit(JSDocInfo_Bit_NOCOMPILE)
}

// IsNoDts returns whether the `@nodts` annotation is present on this
// JSDocInfo.
//
// Java: JSDocInfo.isNoDts()
func (i *JSDocInfo) IsNoDts() bool {
	return i.CheckBit(JSDocInfo_Bit_NODTS)
}

// IsNoCollapse returns whether the `@nocollapse` annotation is present on
// this JSDocInfo.
//
// Java: JSDocInfo.isNoCollapse()
func (i *JSDocInfo) IsNoCollapse() bool {
	return i.CheckBit(JSDocInfo_Bit_NOCOLLAPSE)
}

// IsNoInline returns whether the `@noinline` annotation is present on this
// JSDocInfo.
//
// Java: JSDocInfo.isNoInline()
func (i *JSDocInfo) IsNoInline() bool {
	return i.CheckBit(JSDocInfo_Bit_NOINLINE)
}

// IsRequireInlining returns whether the `@requireInlining` annotation is
// present on this JSDocInfo.
//
// Java: JSDocInfo.isRequireInlining()
func (i *JSDocInfo) IsRequireInlining() bool {
	return i.CheckBit(JSDocInfo_Bit_REQUIRE_INLINING)
}

// IsEncourageInlining returns whether the `@encourageInlining` annotation is
// present on this JSDocInfo.
//
// Java: JSDocInfo.isEncourageInlining()
func (i *JSDocInfo) IsEncourageInlining() bool {
	return i.CheckBit(JSDocInfo_Bit_ENCOURAGE_INLINING)
}

// IsCollapsibleOrBreakMyCode returns whether the `@collapsibleOrBreakMyCode`
// annotation is present on this JSDocInfo.
//
// Java: JSDocInfo.isCollapsibleOrBreakMyCode()
func (i *JSDocInfo) IsCollapsibleOrBreakMyCode() bool {
	return i.CheckBit(JSDocInfo_Bit_COLLAPSIBLE_OR_BREAK_MY_CODE)
}

// IsPureOrBreakMyCode returns whether the `@pureOrBreakMyCode` annotation is
// present on this JSDocInfo.
//
// Java: JSDocInfo.isPureOrBreakMyCode()
func (i *JSDocInfo) IsPureOrBreakMyCode() bool {
	return i.CheckBit(JSDocInfo_Bit_PURE_OR_BREAK_MY_CODE)
}

// IsProvideGoog returns whether the `@provideGoog` annotation is present on
// this JSDocInfo.
//
// Java: JSDocInfo.isProvideGoog()
func (i *JSDocInfo) IsProvideGoog() bool {
	return i.CheckBit(JSDocInfo_Bit_PROVIDE_GOOG)
}

// IsProvideAlreadyProvided returns whether the `@provideAlreadyProvided`
// annotation is present on this JSDocInfo.
//
// Java: JSDocInfo.isProvideAlreadyProvided()
func (i *JSDocInfo) IsProvideAlreadyProvided() bool {
	return i.CheckBit(JSDocInfo_Bit_PROVIDE_ALREADY_PROVIDED)
}

// ContainsDeclarationExcludingTypelessConst returns whether there is a
// declaration present on this JSDocInfo.
//
// Does not consider `@const` (without a following type) to indicate a
// declaration. Whether you want this method, or the`containsDeclaration`
// that includes const, depends on whether you want to consider
// `/** @const * / a.b.c = 0` a declaration or not.
//
// Java: JSDocInfo.containsDeclarationExcludingTypelessConst()
func (i *JSDocInfo) ContainsDeclarationExcludingTypelessConst() bool {
	return (i.propertyBits&
		(JSDocInfo_Bit_CONSTRUCTOR.Mask()|
			JSDocInfo_Bit_DEFINE.Mask()|
			JSDocInfo_Bit_OVERRIDE.Mask()|
			JSDocInfo_Bit_EXPORT.Mask()|
			JSDocInfo_Bit_DEPRECATED.Mask()|
			JSDocInfo_Bit_INTERFACE.Mask()|
			JSDocInfo_Bit_IMPLICITCAST.Mask()|
			JSDocInfo_Bit_NOSIDEEFFECTS.Mask()|
			JSDocInfo_Bit_RECORD.Mask())) != 0 ||
		(i.propertyKeysBitset&
			(JSDocInfo_TYPE.mask|
				JSDocInfo_RETURN_TYPE.mask|
				JSDocInfo_ENUM_PARAMETER_TYPE.mask|
				JSDocInfo_TYPEDEF_TYPE.mask|
				JSDocInfo_THIS_TYPE.mask|
				JSDocInfo_PARAMETERS.mask|
				JSDocInfo_IMPLEMENTED_INTERFACES.mask|
				JSDocInfo_BASE_TYPE.mask|
				JSDocInfo_VISIBILITY.mask)) != 0
}

// ContainsTypeDeclaration returns whether this JSDoc contains a type
// declaration such as `/** @type {string}`.
//
// Java: JSDocInfo.containsTypeDeclaration()
func (i *JSDocInfo) ContainsTypeDeclaration() bool {
	return (i.propertyKeysBitset &
		(JSDocInfo_TYPE.mask |
			JSDocInfo_RETURN_TYPE.mask |
			JSDocInfo_ENUM_PARAMETER_TYPE.mask |
			JSDocInfo_TYPEDEF_TYPE.mask |
			JSDocInfo_THIS_TYPE.mask |
			JSDocInfo_PARAMETERS.mask |
			JSDocInfo_BASE_TYPE.mask)) != 0
}

// ContainsDeclaration returns whether there is a declaration present on this
// JSDocInfo, including a typeless @const like `/** @const * / a.b.c = 0`.
//
// Java: JSDocInfo.containsDeclaration()
func (i *JSDocInfo) ContainsDeclaration() bool {
	return i.ContainsDeclarationExcludingTypelessConst() || i.CheckBit(JSDocInfo_Bit_CONST)
}

// ContainsFunctionDeclaration returns whether there is a declaration of a
// callable type.
//
// Deprecated: This method is quite heuristic, looking for @type annotations
// that start with "function". Other methods like containsDeclaration() and
// containsTypeDefinition are generally preferred.
//
// Java: JSDocInfo.containsFunctionDeclaration()
func (i *JSDocInfo) ContainsFunctionDeclaration() bool {
	if i.CheckBit(JSDocInfo_Bit_CONSTRUCTOR) ||
		(i.propertyKeysBitset&(JSDocInfo_RETURN_TYPE.mask|JSDocInfo_THIS_TYPE.mask|JSDocInfo_PARAMETERS.mask)) != 0 {
		return true
	}
	type_ := jsdocGet[*JSTypeExpression](JSDocInfo_TYPE.getRaw(i))
	if type_ != nil {
		return type_.GetRoot().IsFunction()
	}
	return i.CheckBit(JSDocInfo_Bit_NOSIDEEFFECTS)
}

// ContainsTypeDefinition is for jsdocs that create new types. Not to be
// confused with jsdocs that declare the type of a variable or property.
//
// Java: JSDocInfo.containsTypeDefinition()
func (i *JSDocInfo) ContainsTypeDefinition() bool {
	return (i.propertyBits&(JSDocInfo_Bit_CONSTRUCTOR.Mask()|JSDocInfo_Bit_INTERFACE.Mask())) != 0 ||
		(i.propertyKeysBitset&(JSDocInfo_ENUM_PARAMETER_TYPE.mask|JSDocInfo_TYPEDEF_TYPE.mask)) != 0
}

// IsAtSignCodePresent returns whether the `@code` is present within this
// JSDocInfo.
//
// Java: JSDocInfo.isAtSignCodePresent()
func (i *JSDocInfo) IsAtSignCodePresent() bool {
	entireComment, ok := i.GetOriginalCommentStringValue()
	return ok && strings.Contains(entireComment, "@code")
}

// GetVisibility gets the visibility specified by `@private`, `@protected` or
// `@public` annotation. If no visibility is specified, visibility is
// inherited from the base class.
//
// Java: JSDocInfo.getVisibility()
func (i *JSDocInfo) GetVisibility() JSDocInfo_Visibility {
	visibility := JSDocInfo_VISIBILITY.getRaw(i)
	if visibility != nil {
		return jsdocGet[JSDocInfo_Visibility](visibility)
	}
	return JSDocInfo_Visibility_INHERITED
}

// GetParameterEntries returns the ordered @param (Name, Type) entries without allocating.
func (i *JSDocInfo) GetParameterEntries() []JSDocParam {
	return jsdocGet[[]JSDocParam](JSDocInfo_PARAMETERS.getRaw(i))
}

// GetParameterType gets the type of a given named parameter; null if this
// parameter is not defined or has a null type.
//
// Java: JSDocInfo.getParameterType(String)
func (i *JSDocInfo) GetParameterType(parameter string) *JSTypeExpression {
	for _, p := range i.GetParameterEntries() {
		if p.Name == parameter {
			return p.Type
		}
	}
	return nil
}

// HasParameter returns whether the parameter is defined.
//
// Java: JSDocInfo.hasParameter(String)
func (i *JSDocInfo) HasParameter(parameter string) bool {
	for _, p := range i.GetParameterEntries() {
		if p.Name == parameter {
			return true
		}
	}
	return false
}

// HasParameterType returns whether the parameter has an attached type:
// true if the parameter has an attached type, false if the parameter has no
// attached type or does not exist.
//
// Java: JSDocInfo.hasParameterType(String)
func (i *JSDocInfo) HasParameterType(parameter string) bool {
	return i.GetParameterType(parameter) != nil
}

// GetParameterNames returns the set of names of the defined parameters. The
// iteration order of the returned set is the order in which parameters are
// defined in the JSDoc, rather than the order in which the function declares
// them. The returned set is immutable.
//
// Java: JSDocInfo.getParameterNames()
func (i *JSDocInfo) GetParameterNames() *jutil.LinkedHashSet[string] {
	params := i.GetParameterEntries()
	if len(params) > 0 {
		set := jutil.NewLinkedHashSet[string]()
		for _, p := range params {
			set.Add(p.Name)
		}
		return set
	}
	return jsDocInfo_noSuppressions
}

// GetParameterNameAt returns the nth name in the defined parameters. The
// iteration order is the order in which parameters are defined in the
// JSDoc, rather than the order in which the function declares them.
//
// Java: JSDocInfo.getParameterNameAt(int)
func (i *JSDocInfo) GetParameterNameAt(index int) *string {
	name, ok := i.GetParameterNameValueAt(index)
	if !ok {
		return nil
	}
	return &name
}

// GetParameterNameValueAt returns the nth defined parameter name without
// allocating a *string pointer.
func (i *JSDocInfo) GetParameterNameValueAt(index int) (string, bool) {
	params := i.GetParameterEntries()
	if index < 0 || index >= len(params) {
		return "", false
	}
	return params[index].Name, true
}

// GetParameterTypeAt returns the nth defined parameter type expression in JSDoc
// definition order without allocating an iterator closure.
func (i *JSDocInfo) GetParameterTypeAt(index int) *JSTypeExpression {
	params := i.GetParameterEntries()
	if index < 0 || index >= len(params) {
		return nil
	}
	return params[index].Type
}

// GetParameterCount gets the number of parameters defined.
//
// Java: JSDocInfo.getParameterCount()
func (i *JSDocInfo) GetParameterCount() int {
	return len(i.GetParameterEntries())
}

// HasThrowsAnnotations reports whether any @throws annotations are present.
func (i *JSDocInfo) HasThrowsAnnotations() bool {
	return (i.propertyKeysBitset & JSDocInfo_THROWS_ANNOTATIONS.mask) != 0
}

// GetThrowsAnnotations returns the list of thrown types and descriptions as
// text.
//
// Java: JSDocInfo.getThrowsAnnotations()
func (i *JSDocInfo) GetThrowsAnnotations() []string {
	if annotations := jsdocGet[[]string](JSDocInfo_THROWS_ANNOTATIONS.getRaw(i)); annotations != nil {
		return annotations
	}
	return []string{}
}

// HasEnumParameterType returns whether an enum parameter type, specified
// using the `@enum` annotation, is present on this JSDoc.
//
// Java: JSDocInfo.hasEnumParameterType()
func (i *JSDocInfo) HasEnumParameterType() bool {
	return (i.propertyKeysBitset & JSDocInfo_ENUM_PARAMETER_TYPE.mask) != 0
}

// HasTypedefType returns whether a typedef parameter type, specified using
// the `@typedef` annotation, is present on this JSDoc.
//
// Java: JSDocInfo.hasTypedefType()
func (i *JSDocInfo) HasTypedefType() bool {
	return (i.propertyKeysBitset & JSDocInfo_TYPEDEF_TYPE.mask) != 0
}

// HasReturnType returns whether this JSDocInfo contains a type for `@return`
// annotation.
//
// Java: JSDocInfo.hasReturnType()
func (i *JSDocInfo) HasReturnType() bool {
	return (i.propertyKeysBitset & JSDocInfo_RETURN_TYPE.mask) != 0
}

// HasType returns whether a type, specified using the `@type` annotation, is
// present on this JSDoc.
//
// Java: JSDocInfo.hasType()
func (i *JSDocInfo) HasType() bool {
	return (i.propertyKeysBitset & JSDocInfo_TYPE.mask) != 0
}

// Java: JSDocInfo.hasTypeInformation()
func (i *JSDocInfo) HasTypeInformation() bool {
	return (i.propertyKeysBitset &
		(JSDocInfo_TYPE.mask | JSDocInfo_TYPEDEF_TYPE.mask | JSDocInfo_ENUM_PARAMETER_TYPE.mask | JSDocInfo_RETURN_TYPE.mask)) != 0
}

// IsInlineType returns whether the type annotation was inlined.
//
// Java: JSDocInfo.isInlineType()
func (i *JSDocInfo) IsInlineType() bool {
	return i.CheckBit(JSDocInfo_Bit_INLINE_TYPE)
}

// GetReturnType gets the return type specified by the `@return` annotation.
//
// Java: JSDocInfo.getReturnType()
func (i *JSDocInfo) GetReturnType() *JSTypeExpression {
	return jsdocGet[*JSTypeExpression](JSDocInfo_RETURN_TYPE.getRaw(i))
}

// GetEnumParameterType gets the enum parameter type specified by the `@enum`
// annotation.
//
// Java: JSDocInfo.getEnumParameterType()
func (i *JSDocInfo) GetEnumParameterType() *JSTypeExpression {
	return jsdocGet[*JSTypeExpression](JSDocInfo_ENUM_PARAMETER_TYPE.getRaw(i))
}

// GetTypedefType gets the typedef type specified by the `@type` annotation.
//
// Java: JSDocInfo.getTypedefType()
func (i *JSDocInfo) GetTypedefType() *JSTypeExpression {
	return jsdocGet[*JSTypeExpression](JSDocInfo_TYPEDEF_TYPE.getRaw(i))
}

// GetType gets the type specified by the `@type` annotation.
//
// Java: JSDocInfo.getType()
func (i *JSDocInfo) GetType() *JSTypeExpression {
	return jsdocGet[*JSTypeExpression](JSDocInfo_TYPE.getRaw(i))
}

// GetThisType gets the type specified by the `@this` annotation.
//
// Java: JSDocInfo.getThisType()
func (i *JSDocInfo) GetThisType() *JSTypeExpression {
	return jsdocGet[*JSTypeExpression](JSDocInfo_THIS_TYPE.getRaw(i))
}

// HasThisType returns whether this JSDocInfo contains a type for `@this`
// annotation.
//
// Java: JSDocInfo.hasThisType()
func (i *JSDocInfo) HasThisType() bool {
	return (i.propertyKeysBitset & JSDocInfo_THIS_TYPE.mask) != 0
}

// GetBaseType gets the base type specified by the `@extends` annotation.
//
// Java: JSDocInfo.getBaseType()
func (i *JSDocInfo) GetBaseType() *JSTypeExpression {
	return jsdocGet[*JSTypeExpression](JSDocInfo_BASE_TYPE.getRaw(i))
}

func (i *JSDocInfo) getStringProp(p *JSDocInfo_Property) (string, bool) {
	v := p.getRaw(i)
	if v == nil {
		return "", false
	}
	return jsdocGet[string](v), true
}

func (i *JSDocInfo) getStringPtrProp(p *JSDocInfo_Property) *string {
	s, ok := i.getStringProp(p)
	if !ok {
		return nil
	}
	return &s
}

// GetDescription gets the description specified by the `@desc` annotation.
//
// Java: JSDocInfo.getDescription()
func (i *JSDocInfo) GetDescription() *string {
	return i.getStringPtrProp(JSDocInfo_DESCRIPTION)
}

// GetDescriptionString returns the `@desc` description without allocating a *string.
func (i *JSDocInfo) GetDescriptionString() (string, bool) {
	return i.getStringProp(JSDocInfo_DESCRIPTION)
}

// GetTsTypes gets the ts type declarations specified by the `@tsType`
// annotations.
//
// Java: JSDocInfo.getTsTypes()
func (i *JSDocInfo) GetTsTypes() []string {
	if types := jsdocGet[[]string](JSDocInfo_TS_TYPES.getRaw(i)); types != nil {
		return slices.Clone(types)
	}
	return []string{}
}

// GetMeaning gets the meaning specified by the `@meaning` annotation.
//
// In localization systems, two messages with the same content but different
// "meanings" may be translated differently. By default, we use the name of
// the variable that the message is initialized to as the "meaning" of the
// message.
//
// But some code generators (like Closure Templates) inject their own meaning
// with the jsdoc `@meaning` annotation.
//
// Java: JSDocInfo.getMeaning()
func (i *JSDocInfo) GetMeaning() *string {
	return i.getStringPtrProp(JSDocInfo_MEANING)
}

// GetMeaningString returns the `@meaning` string without allocating a *string.
func (i *JSDocInfo) GetMeaningString() (string, bool) {
	return i.getStringProp(JSDocInfo_MEANING)
}

// GetAlternateMessageId gets the alternate message ID specified by the
// `@alternateMessageId` annotation.
//
// In localization systems, if we migrate from one message ID algorithm to
// another, we can specify the old one via `@alternateMessageId`. This allows
// the product to use the previous translation while waiting for the new one
// to be translated.
//
// Some code generators (like Closure Templates) inject this.
//
// Java: JSDocInfo.getAlternateMessageId()
func (i *JSDocInfo) GetAlternateMessageId() *string {
	return i.getStringPtrProp(JSDocInfo_ALTERNATE_MESSAGE_ID)
}

// GetAlternateMessageIdString returns the `@alternateMessageId` string without allocating a *string.
func (i *JSDocInfo) GetAlternateMessageIdString() (string, bool) {
	return i.getStringProp(JSDocInfo_ALTERNATE_MESSAGE_ID)
}

// GetLendsName gets the name we're lending to in a `@lends` annotation.
//
// In many reflection APIs, you pass an anonymous object to a function, and
// that function mixes the anonymous object into another object. The `@lends`
// annotation allows the type system to track those property assignments.
//
// Java: JSDocInfo.getLendsName()
func (i *JSDocInfo) GetLendsName() *JSTypeExpression {
	return jsdocGet[*JSTypeExpression](JSDocInfo_LENDS_NAME.getRaw(i))
}

// Java: JSDocInfo.hasLendsName()
func (i *JSDocInfo) HasLendsName() bool {
	return (i.propertyKeysBitset & JSDocInfo_LENDS_NAME.mask) != 0
}

// GetClosurePrimitiveId returns the `@closurePrimitive {id}` identifier.
//
// Java: JSDocInfo.getClosurePrimitiveId()
func (i *JSDocInfo) GetClosurePrimitiveId() *string {
	return i.getStringPtrProp(JSDocInfo_CLOSURE_PRIMITIVE_ID)
}

// GetClosurePrimitiveIdString returns the `@closurePrimitive` id without allocating a *string.
func (i *JSDocInfo) GetClosurePrimitiveIdString() (string, bool) {
	return i.getStringProp(JSDocInfo_CLOSURE_PRIMITIVE_ID)
}

// HasClosurePrimitiveId reports whether this JSDoc is annotated with
// `@closurePrimitive`.
//
// Java: JSDocInfo.hasClosurePrimitiveId()
func (i *JSDocInfo) HasClosurePrimitiveId() bool {
	return (i.propertyKeysBitset & JSDocInfo_CLOSURE_PRIMITIVE_ID.mask) != 0
}

// IsNgInject returns whether JSDoc is annotated with `@ngInject` annotation.
//
// Java: JSDocInfo.isNgInject()
func (i *JSDocInfo) IsNgInject() bool {
	return i.CheckBit(JSDocInfo_Bit_NG_INJECT)
}

// IsWizaction returns whether JSDoc is annotated with `@wizaction`
// annotation.
//
// Java: JSDocInfo.isWizaction()
func (i *JSDocInfo) IsWizaction() bool {
	return i.CheckBit(JSDocInfo_Bit_WIZ_ACTION)
}

// IsWizcallback returns if JSDoc is annotated with `@wizcallback` annotation.
//
// Java: JSDocInfo.isWizcallback()
func (i *JSDocInfo) IsWizcallback() bool {
	return i.CheckBit(JSDocInfo_Bit_WIZ_CALLBACK)
}

// IsPolymerBehavior returns whether JSDoc is annotated with
// `@polymerBehavior` annotation.
//
// Java: JSDocInfo.isPolymerBehavior()
func (i *JSDocInfo) IsPolymerBehavior() bool {
	return i.CheckBit(JSDocInfo_Bit_POLYMER_BEHAVIOR)
}

// IsPolymer returns whether JSDoc is annotated with `@polymer` annotation.
//
// Java: JSDocInfo.isPolymer()
func (i *JSDocInfo) IsPolymer() bool {
	return i.CheckBit(JSDocInfo_Bit_POLYMER)
}

// IsCustomElement returns whether JSDoc is annotated with `@customElement`
// annotation.
//
// Java: JSDocInfo.isCustomElement()
func (i *JSDocInfo) IsCustomElement() bool {
	return i.CheckBit(JSDocInfo_Bit_CUSTOM_ELEMENT)
}

// IsMixinClass returns whether JSDoc is annotated with `@mixinClass`
// annotation.
//
// Java: JSDocInfo.isMixinClass()
func (i *JSDocInfo) IsMixinClass() bool {
	return i.CheckBit(JSDocInfo_Bit_MIXIN_CLASS)
}

// IsMixinFunction returns whether JSDoc is annotated with `@mixinFunction`
// annotation.
//
// Java: JSDocInfo.isMixinFunction()
func (i *JSDocInfo) IsMixinFunction() bool {
	return i.CheckBit(JSDocInfo_Bit_MIXIN_FUNCTION)
}

// Java: JSDocInfo.isSassGeneratedCssTs()
func (i *JSDocInfo) IsSassGeneratedCssTs() bool {
	return i.CheckBit(JSDocInfo_Bit_SASS_GENERATED_CSS_TS)
}

// IsClosureUnawareCode returns whether JSDoc is annotated with the
// `@closureUnawareCode` annotation. Generally, prefer calling
// Node#isClosureUnawareCode() instead of referencing this jsdoc, as this
// annotation is NOT serialized to TypedAST (because an equivalent bit is
// instead set in the SourceFile's TypedAST representation).
//
// Java: JSDocInfo.isClosureUnawareCode()
func (i *JSDocInfo) IsClosureUnawareCode() bool {
	return (i.propertyKeysBitset & JSDocInfo_CLOSURE_UNAWARE_CONFIG_VALUE.mask) != 0
}

// GetPerFileClosureUnawareMode returns the value of the
// `@perFileClosureUnawareMode` annotation, or null if not present.
//
// Java: JSDocInfo.getPerFileClosureUnawareMode()
func (i *JSDocInfo) GetPerFileClosureUnawareMode() *JSDocInfo_PerFileClosureUnawareMode {
	m, ok := i.GetPerFileClosureUnawareModeValue()
	if !ok {
		return nil
	}
	return &m
}

// GetPerFileClosureUnawareModeValue returns the `@perFileClosureUnawareMode`
// value without allocating a pointer.
func (i *JSDocInfo) GetPerFileClosureUnawareModeValue() (JSDocInfo_PerFileClosureUnawareMode, bool) {
	v := JSDocInfo_CLOSURE_UNAWARE_CONFIG_VALUE.getRaw(i)
	if v == nil {
		return 0, false
	}
	return jsdocGet[JSDocInfo_PerFileClosureUnawareMode](v), true
}

// IsUsedViaDotConstructor returns whether JSDoc is annotated with the
// `@usedViaDotConstructor` annotation.
//
// Java: JSDocInfo.isUsedViaDotConstructor()
func (i *JSDocInfo) IsUsedViaDotConstructor() bool {
	return i.CheckBit(JSDocInfo_Bit_USED_VIA_DOT_CONSTRUCTOR)
}

// GetLicense gets the description specified by the `@license` annotation.
//
// Java: JSDocInfo.getLicense()
func (i *JSDocInfo) GetLicense() *string {
	return i.getStringPtrProp(JSDocInfo_LICENSE)
}

// GetLicenseString returns the `@license` text without allocating a *string.
func (i *JSDocInfo) GetLicenseString() (string, bool) {
	return i.getStringProp(JSDocInfo_LICENSE)
}

// Java: JSDocInfo.toString()
func (i *JSDocInfo) String() string {
	return i.ToStringVerbose()
}

// Java: JSDocInfo.toStringVerbose()
func (i *JSDocInfo) ToStringVerbose() string {
	var bitset any
	if i.propertyBits != 0 {
		bitset = jutil.LongToHexString(i.propertyBits)
	}
	helper := jutil.MoreObjects_ToStringHelper("JSDocInfo").Add("bitset", bitset)

	for pb := uint64(i.propertyBits); pb != 0; {
		low := bits.TrailingZeros64(pb)
		pb &^= uint64(1) << low
		if low < len(jsdocInfo_Bit_name) {
			helper.AddBool("bit:"+JSDocInfo_Bit(low).BitName(), true)
		}
	}

	bits_ := i.propertyKeysBitset
	index := 0
	for bits_ > 0 {
		low := bits.TrailingZeros64(uint64(bits_))
		bits_ &^= int64(1) << low
		prop := JSDocInfo_Property_values[low]
		val := i.GetPropertyValueByIndex(index)
		switch prop.kind {
		case jsDocPropKindStringList:
			if s := jsdocGet[[]string](val); s != nil {
				val = jutil.CollectionToString(s)
			}
		case jsDocPropKindMarkerList:
			if s := jsdocGet[[]JSDocInfo_Marker](val); s != nil {
				var b strings.Builder
				b.WriteByte('[')
				for idx := range s {
					if idx > 0 {
						b.WriteString(", ")
					}
					b.WriteString((&s[idx]).String())
				}
				b.WriteByte(']')
				val = b.String()
			}
		case jsDocPropKindTypeList:
			if s := jsdocGet[[]*JSTypeExpression](val); s != nil {
				val = jutil.CollectionToString(s)
			}
		case jsDocPropKindTypeMap:
			if s := jsdocGet[[]JSDocParam](val); s != nil {
				val = formatJsDocParams(s)
			}
		case jsDocPropKindStringMap:
			if s := jsdocGet[[]jsDocKV](val); s != nil {
				val = formatJsDocKVs(s)
			}
		case jsDocPropKindSuppressions:
			if s := jsdocGet[[]JSDocSuppression](val); s != nil {
				val = FormatJSDocSuppressions(s)
			}
		}
		helper = helper.Add(prop.name, val)
		index++
	}

	return helper.OmitNullValues().String()
}

// HasBaseType returns whether this JSDocInfo contains a type for `@extends`
// annotation.
//
// Java: JSDocInfo.hasBaseType()
func (i *JSDocInfo) HasBaseType() bool {
	return (i.propertyKeysBitset & JSDocInfo_BASE_TYPE.mask) != 0
}

// GetImplementedInterfaces returns the types specified by the `@implements`
// annotation: an immutable list of JSTypeExpression objects that can be
// resolved to types.
//
// Java: JSDocInfo.getImplementedInterfaces()
func (i *JSDocInfo) GetImplementedInterfaces() []*JSTypeExpression {
	return JSDocInfo_IMPLEMENTED_INTERFACES.GetUnmodifiable(i)
}

// GetImplementedInterfaceCount gets the number of interfaces specified by
// the `@implements` annotation.
//
// Java: JSDocInfo.getImplementedInterfaceCount()
func (i *JSDocInfo) GetImplementedInterfaceCount() int {
	return len(jsdocGet[[]*JSTypeExpression](JSDocInfo_IMPLEMENTED_INTERFACES.getRaw(i)))
}

// GetExtendedInterfaces returns the interfaces extended by an interface: an
// immutable list of JSTypeExpression objects that can be resolved to types.
//
// Java: JSDocInfo.getExtendedInterfaces()
func (i *JSDocInfo) GetExtendedInterfaces() []*JSTypeExpression {
	return JSDocInfo_EXTENDED_INTERFACES.GetUnmodifiable(i)
}

// GetExtendedInterfacesCount gets the number of extended interfaces
// specified.
//
// Java: JSDocInfo.getExtendedInterfacesCount()
func (i *JSDocInfo) GetExtendedInterfacesCount() int {
	return len(jsdocGet[[]*JSTypeExpression](JSDocInfo_EXTENDED_INTERFACES.getRaw(i)))
}

// GetDeprecationReason returns the deprecation reason or null if none
// specified.
//
// Java: JSDocInfo.getDeprecationReason()
func (i *JSDocInfo) GetDeprecationReason() *string {
	return i.getStringPtrProp(JSDocInfo_DEPRECATION_REASON)
}

// GetDeprecationReasonString returns the `@deprecated` reason without allocating a *string.
func (i *JSDocInfo) GetDeprecationReasonString() (string, bool) {
	return i.getStringProp(JSDocInfo_DEPRECATION_REASON)
}

func (i *JSDocInfo) suppressions() []JSDocSuppression {
	return jsdocGet[[]JSDocSuppression](JSDocInfo_SUPPRESSIONS.getRaw(i))
}

// GetSuppressions returns the set of suppressed warnings.
//
// Java: JSDocInfo.getSuppressions()
func (i *JSDocInfo) GetSuppressions() *jutil.LinkedHashSet[string] {
	suppressions := i.suppressions()
	if len(suppressions) == 0 {
		// Go-only: performance — ImmutableSet.of() is a shared singleton;
		// callers only read the result, so share one empty set instead of
		// allocating one per call (this runs for every JSDoc, externs
		// included).
		return jsDocInfo_noSuppressions
	}
	names := jutil.NewLinkedHashSet[string]()
	for _, sup := range suppressions {
		for _, name := range sup.Warnings {
			names.Add(name)
		}
	}
	return names
}

// HasAnySuppressions reports whether i has at least one @suppress warning name
// without allocating a merged set when multiple @suppress tags are present.
func (i *JSDocInfo) HasAnySuppressions() bool {
	for _, sup := range i.suppressions() {
		if len(sup.Warnings) > 0 {
			return true
		}
	}
	return false
}

// HasSuppression reports whether warning name is in i's @suppress set without
// allocating a merged set when multiple @suppress tags are present.
func (i *JSDocInfo) HasSuppression(name string) bool {
	for _, sup := range i.suppressions() {
		if slices.Contains(sup.Warnings, name) {
			return true
		}
	}
	return false
}

// ForEachSuppression visits each suppressed warning name on i (may visit a
// duplicate name across multiple @suppress tags if the caller is idempotent).
func (i *JSDocInfo) ForEachSuppression(visit func(string) bool) {
	for _, sup := range i.suppressions() {
		for _, name := range sup.Warnings {
			if !visit(name) {
				return
			}
		}
	}
}

// jsDocInfo_noSuppressions is GetSuppressions' result for a JSDocInfo
// without suppressions (ImmutableSet.of()). Never modified.
var jsDocInfo_noSuppressions = jutil.NewLinkedHashSet[string]()

// GetSuppressionsAndTheirDescription returns the slice of @suppress warning
// groups and their corresponding descriptions.
//
// Java: JSDocInfo.getSuppressionsAndTheirDescription()
func (i *JSDocInfo) GetSuppressionsAndTheirDescription() []JSDocSuppression {
	if suppressions := i.suppressions(); suppressions != nil {
		return suppressions
	}
	return []JSDocSuppression{}
}

// GetModifies returns the set of sideeffect notations.
//
// Java: JSDocInfo.getModifies()
func (i *JSDocInfo) GetModifies() *jutil.LinkedHashSet[string] {
	modifies := jsdocGet[*jutil.LinkedHashSet[string]](JSDocInfo_MODIFIES.getRaw(i))
	if modifies != nil {
		return modifies
	}
	return jsDocInfo_noSuppressions
}

func (i *JSDocInfo) parameterDescriptions() []jsDocKV {
	return jsdocGet[[]jsDocKV](JSDocInfo_PARAMETER_DESCRIPTIONS.getRaw(i))
}

// HasDescriptionForParameter returns whether a description exists for the
// parameter with the specified name.
//
// Java: JSDocInfo.hasDescriptionForParameter(String)
func (i *JSDocInfo) HasDescriptionForParameter(name string) bool {
	_, ok := i.GetDescriptionForParameterString(name)
	return ok
}

// GetDescriptionForParameter returns the description for the parameter with
// the given name, if its exists.
//
// Java: JSDocInfo.getDescriptionForParameter(String)
func (i *JSDocInfo) GetDescriptionForParameter(name string) *string {
	d, ok := i.GetDescriptionForParameterString(name)
	if !ok {
		return nil
	}
	return &d
}

// GetDescriptionForParameterString returns the description for the parameter
// with the given name without allocating a *string.
func (i *JSDocInfo) GetDescriptionForParameterString(name string) (string, bool) {
	for _, kv := range i.parameterDescriptions() {
		if kv.key == name {
			return kv.val, true
		}
	}
	return "", false
}

// GetAuthors returns the list of authors or null if none.
//
// Java: JSDocInfo.getAuthors()
func (i *JSDocInfo) GetAuthors() []string {
	return jsdocGet[[]string](JSDocInfo_AUTHORS.getRaw(i))
}

// GetReferences returns the list of references or null if none.
//
// Java: JSDocInfo.getReferences()
func (i *JSDocInfo) GetReferences() []string {
	return jsdocGet[[]string](JSDocInfo_SEES.getRaw(i))
}

// GetReturnDescription returns the description of the returned object or
// null if none specified.
//
// Java: JSDocInfo.getReturnDescription()
func (i *JSDocInfo) GetReturnDescription() *string {
	return i.getStringPtrProp(JSDocInfo_RETURN_DESCRIPTION)
}

// GetReturnDescriptionString returns the `@return` description without allocating a *string.
func (i *JSDocInfo) GetReturnDescriptionString() (string, bool) {
	return i.getStringProp(JSDocInfo_RETURN_DESCRIPTION)
}

// GetBlockDescription returns the block-level description or null if none
// specified.
//
// Java: JSDocInfo.getBlockDescription()
func (i *JSDocInfo) GetBlockDescription() *string {
	return i.getStringPtrProp(JSDocInfo_BLOCK_DESCRIPTION)
}

// GetBlockDescriptionString returns the block-level description without allocating a *string.
func (i *JSDocInfo) GetBlockDescriptionString() (string, bool) {
	return i.getStringProp(JSDocInfo_BLOCK_DESCRIPTION)
}

// HasFileOverview returns whether this has a fileoverview flag.
//
// Java: JSDocInfo.hasFileOverview()
func (i *JSDocInfo) HasFileOverview() bool {
	return i.CheckBit(JSDocInfo_Bit_FILEOVERVIEW)
}

// GetFileOverview returns the file overview or null if none specified.
//
// Java: JSDocInfo.getFileOverview()
func (i *JSDocInfo) GetFileOverview() *string {
	return i.getStringPtrProp(JSDocInfo_FILEOVERVIEW_DESCRIPTION)
}

// GetFileOverviewString returns the `@fileoverview` description without allocating a *string.
func (i *JSDocInfo) GetFileOverviewString() (string, bool) {
	return i.getStringProp(JSDocInfo_FILEOVERVIEW_DESCRIPTION)
}

// HasEnhance returns whether this has an enhanced namespace.
//
// Java: JSDocInfo.hasEnhance()
func (i *JSDocInfo) HasEnhance() bool {
	return i.CheckBit(JSDocInfo_Bit_ENHANCED_NAMESPACE)
}

// GetEnhance returns the enhanced namespace or null if none is specified.
//
// Java: JSDocInfo.getEnhance()
func (i *JSDocInfo) GetEnhance() *string {
	return i.getStringPtrProp(JSDocInfo_ENHANCED_NAMESPACE)
}

// GetEnhanceString returns the `@enhance` namespace without allocating a *string.
func (i *JSDocInfo) GetEnhanceString() (string, bool) {
	return i.getStringProp(JSDocInfo_ENHANCED_NAMESPACE)
}

// HasMods returns whether this has a modded namespace.
//
// Java: JSDocInfo.hasMods()
func (i *JSDocInfo) HasMods() bool {
	return (i.propertyKeysBitset & JSDocInfo_MODS.mask) != 0
}

// GetMods returns the modded namespace or null if none is specified.
//
// Java: JSDocInfo.getMods()
func (i *JSDocInfo) GetMods() *string {
	return i.getStringPtrProp(JSDocInfo_MODS)
}

// GetModsString returns the `@mods` namespace without allocating a *string.
func (i *JSDocInfo) GetModsString() (string, bool) {
	return i.getStringProp(JSDocInfo_MODS)
}

// GetMarkers gets the list of all markers for the documentation in this
// JSDoc.
//
// Java: JSDocInfo.getMarkers()
func (i *JSDocInfo) GetMarkers() []JSDocInfo_Marker {
	if markers := jsdocGet[[]JSDocInfo_Marker](JSDocInfo_MARKERS.getRaw(i)); markers != nil {
		return markers
	}
	return []JSDocInfo_Marker{}
}

// GetTemplateTypeEntries returns the ordered @template (Name, Type) entries without allocating.
func (i *JSDocInfo) GetTemplateTypeEntries() []JSDocParam {
	if i == nil {
		return nil
	}
	return jsdocGet[[]JSDocParam](JSDocInfo_TEMPLATE_TYPE_NAMES.getRaw(i))
}

// GetTemplateTypeNames gets the @template type names.
//
// Excludes @template types from TTL; get those with
// getTypeTransformations().
//
// Java: JSDocInfo.getTemplateTypeNames()
func (i *JSDocInfo) GetTemplateTypeNames() []string {
	params := i.GetTemplateTypeEntries()
	if len(params) == 0 {
		return []string{}
	}
	out := make([]string, len(params))
	for idx, p := range params {
		out[idx] = p.Name
	}
	return out
}

// HasTemplateTypeName reports whether name is declared as a @template type name.
func (i *JSDocInfo) HasTemplateTypeName(name string) bool {
	for _, p := range i.GetTemplateTypeEntries() {
		if p.Name == name {
			return true
		}
	}
	return false
}

var (
	jsDocInfo_noTemplateTypes       = jutil.NewLinkedHashMap[string, *JSTypeExpression]()
	jsDocInfo_noTypeTransformations = jutil.NewLinkedHashMap[string, Node]()
)

// Java: JSDocInfo.getTemplateTypes()
func (i *JSDocInfo) GetTemplateTypes() *jutil.LinkedHashMap[string, *JSTypeExpression] {
	params := i.GetTemplateTypeEntries()
	if len(params) != 0 {
		return jsDocParamsToMap(params)
	}
	return jsDocInfo_noTemplateTypes
}

// GetTypeTransformations gets the type transformations.
//
// Java: JSDocInfo.getTypeTransformations()
func (i *JSDocInfo) GetTypeTransformations() *jutil.LinkedHashMap[string, Node] {
	m := jsdocGet[*jutil.LinkedHashMap[string, Node]](JSDocInfo_TYPE_TRANSFORMATIONS.getRaw(i))
	if m != nil {
		return m
	}
	return jsDocInfo_noTypeTransformations
}

// Java: JSDocInfo.hasModifies()
func (i *JSDocInfo) HasModifies() bool {
	return (i.propertyKeysBitset & JSDocInfo_MODIFIES.mask) != 0
}

// GetOriginalCommentString returns the original JSDoc comment string.
// Returns null unless parseJsDocDocumentation is enabled via the
// ParserConfig.
//
// Java: JSDocInfo.getOriginalCommentString()
func (i *JSDocInfo) GetOriginalCommentString() *string {
	return i.getStringPtrProp(JSDocInfo_SOURCE_COMMENT)
}

// GetOriginalCommentStringValue returns the original JSDoc comment string without allocating a *string.
func (i *JSDocInfo) GetOriginalCommentStringValue() (string, bool) {
	return i.getStringProp(JSDocInfo_SOURCE_COMMENT)
}

// Java: JSDocInfo.getOriginalCommentPosition()
func (i *JSDocInfo) GetOriginalCommentPosition() int {
	pos := JSDocInfo_ORIGINAL_COMMENT_POSITION.getRaw(i)
	if pos != nil {
		return jsdocGet[int](pos)
	}
	return 0
}

// ModifiesThis gets the value of the @modifies{this} annotation stored in
// the doc info.
//
// Java: JSDocInfo.modifiesThis()
func (i *JSDocInfo) ModifiesThis() bool {
	modifies := jsdocGet[*jutil.LinkedHashSet[string]](JSDocInfo_MODIFIES.getRaw(i))
	return modifies != nil && modifies.Contains("this")
}

// HasSideEffectsArgumentsAnnotation returns whether the @modifies annotation
// includes "arguments" or any named parameters.
//
// Java: JSDocInfo.hasSideEffectsArgumentsAnnotation()
func (i *JSDocInfo) HasSideEffectsArgumentsAnnotation() bool {
	modifies := jsdocGet[*jutil.LinkedHashSet[string]](JSDocInfo_MODIFIES.getRaw(i))
	if modifies == nil {
		return false
	}
	// TODO(johnlenz): if we start tracking parameters individually this should
	// simply be a check for "arguments".
	return modifies.Size() > 1 || (modifies.Size() == 1 && !modifies.Contains("this"))
}

// GetLogTypeInCompiler returns whether we should log the type of values
// passed to this function.
//
// Java: JSDocInfo.getLogTypeInCompiler()
func (i *JSDocInfo) GetLogTypeInCompiler() bool {
	return i.CheckBit(JSDocInfo_Bit_LOG_TYPE_IN_COMPILER)
}

// Java: JSDocInfo.getPropertyValueByIndex(int)
func (i *JSDocInfo) GetPropertyValueByIndex(index int) any {
	if i.propertyValues == nil {
		panic(jutil.NewIllegalArgumentExceptionMsg("no property value"))
	}

	if (i.propertyKeysBitset & (i.propertyKeysBitset - 1)) != 0 {
		vals := i.propertyValues.([]any)
		if index < 0 || index >= len(vals) {
			panic(jutil.NewArrayIndexOutOfBoundsExceptionForLength(index, len(vals)))
		}
		return vals[index]
	}

	if index != 0 {
		panic(jutil.NewArrayIndexOutOfBoundsException(index))
	}

	return i.propertyValues
}

// Java: JSDocInfo.packPropertyValues(List<Object>)
func JSDocInfo_PackPropertyValues(values []any) any {
	if len(values) == 0 {
		return nil
	}
	if len(values) == 1 {
		return values[0]
	}
	out := make([]any, len(values))
	copy(out, values)
	return out
}
