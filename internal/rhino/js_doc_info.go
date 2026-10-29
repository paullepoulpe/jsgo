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
// List<Marker> → []JSDocInfo_Marker,
// ArrayList<JSTypeExpression> → []*JSTypeExpression,
// LinkedHashMap<String, JSTypeExpression> → []JSDocParam,
// LinkedHashMap<String, String> → []jsDocKV,
// ImmutableMap<ImmutableSet<String>, String> → []JSDocSuppression,
// ImmutableSet<String> → []string (deduplicated), and
// LinkedHashMap<String, Node> → []JSDocTypeTransform.
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
	jsDocPropKindStringSet                         // []string (deduplicated)
	jsDocPropKindNodeMap                           // []JSDocTypeTransform
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
	switch p.kind {
	case jsDocPropKindTypeMap:
		return jsDocParamsToMap(jsdocGet[[]JSDocParam](v))
	case jsDocPropKindStringSet:
		return jsDocStringSetFromSlice(jsdocGet[[]string](v))
	case jsDocPropKindNodeMap:
		return jsDocTypeTransformsToMap(jsdocGet[[]JSDocTypeTransform](v))
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
		return jsdocGet[[]string](v) == nil
	case jsDocPropKindNodeMap:
		return jsdocGet[[]JSDocTypeTransform](v) == nil
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
	case jsDocPropKindStringList, jsDocPropKindStringSet:
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
		return slices.Clone(jsdocGet[[]JSDocTypeTransform](arg))
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
	case jsDocPropKindStringList, jsDocPropKindStringSet:
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
	case jsDocPropKindNodeMap:
		return len(jsdocGet[[]JSDocTypeTransform](value)) == 0
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
		return sameWarningSet(jsdocGet[[]string](left), jsdocGet[[]string](right))
	case jsDocPropKindNodeMap:
		l := jsdocGet[[]JSDocTypeTransform](left)
		r := jsdocGet[[]JSDocTypeTransform](right)
		if len(l) != len(r) {
			return false
		}
		for _, lt := range l {
			found := false
			for _, rt := range r {
				if rt.Name == lt.Name {
					if rt.TypeNode.IsNil() || !lt.TypeNode.IsEquivalentTo(rt.TypeNode) {
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
	}
	return left == right
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

func newJSDocInfoFromBuilder(bits_ int64, propMask int64, entries []any) *JSDocInfo {
	if propMask == 0 {
		return &JSDocInfo{propertyBits: bits_}
	}
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
	if n := bits.OnesCount64(uint64(bits_)); n > cap(builder.entries) {
		builder.entries = make([]any, 0, n)
	}
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
		case jsDocPropKindStringList, jsDocPropKindStringSet:
			if s := jsdocGet[[]string](val); s != nil {
				val = "[" + strings.Join(s, ", ") + "]"
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
				var b strings.Builder
				b.WriteByte('[')
				for idx, elem := range s {
					if idx > 0 {
						b.WriteString(", ")
					}
					if elem == nil {
						b.WriteString("null")
					} else {
						b.WriteString(elem.String())
					}
				}
				b.WriteByte(']')
				val = b.String()
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
		case jsDocPropKindNodeMap:
			if s := jsdocGet[[]JSDocTypeTransform](val); s != nil {
				val = formatJsDocTypeTransforms(s)
			}
		}
		helper = helper.Add(prop.name, val)
		index++
	}

	return helper.OmitNullValues().String()
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
