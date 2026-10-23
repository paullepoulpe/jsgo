// Port of com/google/javascript/rhino/JSDocInfo.java @ closure-compiler v20261015

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
	"iter"
	"math/bits"
	"strings"
	"sync/atomic"
	"unsafe"

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
// String → string, Integer → int, List<String> → *jutil.ArrayList[string],
// ArrayList<JSTypeExpression> → *jutil.ArrayList[*JSTypeExpression],
// LinkedHashMap<String, X> → *jutil.LinkedHashMap[string, X],
// ImmutableSet<String> → *jutil.LinkedHashSet[string], and
// ImmutableMap<ImmutableSet<String>, String> →
// *jutil.ObjectLinkedHashMap[*jutil.LinkedHashSet[string], string].
type JSDocInfo struct {
	propertyBits       int64
	propertyKeysBitset int64

	// This value will be `null` if there is no property, a literal value if
	// there is only one, or an `Object[]` if there are multiple.
	//
	// This is meant to elide eager allocation of single-valued or empty
	// lists.
	//
	// Go-only: performance — the Object[] is a *any to the first element of
	// an exact-size array (its length is the popcount of
	// propertyKeysBitset), so no slice header is boxed.
	propertyValues any
}

// Java: JSDocInfo.serialVersionUID
const JSDocInfo_serialVersionUID int64 = 1

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
	jsDocPropKindPlain      jsDocPropKind = iota // Property
	jsDocPropKindMarkerList                      // MarkerListProperty: ArrayList<Marker>
	jsDocPropKindType                            // TypeProperty: JSTypeExpression
	jsDocPropKindTypeList                        // TypeListProperty: ArrayList<JSTypeExpression>
	jsDocPropKindTypeMap                         // TypeMapProperty: LinkedHashMap<String, JSTypeExpression>
	jsDocPropKindNodeMap                         // NodeMapProperty: LinkedHashMap<String, Node>
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
	p.mask = jutil.Shl64(1, int64(p.bit))
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

// Get returns the property value of info, or null; the cast to T is unsafe
// but guaranteed by builder.
//
// Java: JSDocInfo.Property.get(JSDocInfo)
func (p *JSDocInfo_Property) Get(info *JSDocInfo) any {
	if (info.propertyKeysBitset & p.mask) == 0 {
		return nil
	}
	return info.GetPropertyValueByIndex(bits.OnesCount64(uint64(info.propertyKeysBitset & (p.mask - 1))))
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
		out := jutil.NewArrayListFrom(jsdocGet[*jutil.ArrayList[*JSTypeExpression]](arg).Slice())
		if transform != nil {
			for i := 0; i < out.Size(); i++ {
				elem := out.Get(i)
				if elem != nil {
					out.Set(i, transform(elem))
				}
			}
		}
		return out
	case jsDocPropKindTypeMap:
		out := jutil.NewLinkedHashMapFrom(jsdocGet[*jutil.LinkedHashMap[string, *JSTypeExpression]](arg))
		if transform != nil {
			for k, elem := range out.All() {
				if elem != nil {
					out.Put(k, transform(elem))
				}
			}
		}
		return out
	case jsDocPropKindNodeMap:
		return jutil.NewLinkedHashMapFrom(jsdocGet[*jutil.LinkedHashMap[string, *Node]](arg))
	}
	return arg
}

// Java: JSDocInfo.Property.isDefault(T)
func (p *JSDocInfo_Property) IsDefault(value any) bool {
	if jutil.IsNil(value) {
		return true
	}
	if v, ok := value.(JSDocInfo_Visibility); ok && v == JSDocInfo_Visibility_INHERITED {
		return true
	}
	if c, ok := value.(jutil.JavaCollection); ok && c.IsEmpty() {
		return true
	}
	if m, ok := value.(jutil.JavaMap); ok && m.IsEmpty() {
		return true
	}
	return false
}

// Java: JSDocInfo.Property.equalValues(T, T), and its overrides
func (p *JSDocInfo_Property) EqualValues(left any, right any) bool {
	switch p.kind {
	case jsDocPropKindMarkerList:
		a := jsdocGet[*jutil.ArrayList[*JSDocInfo_Marker]](left)
		b := jsdocGet[*jutil.ArrayList[*JSDocInfo_Marker]](right)
		if a.Size() != b.Size() {
			return false
		}
		for i := 0; i < a.Size(); i++ {
			m1 := a.Get(i)
			m2 := b.Get(i)
			if (m1 == nil) != (m2 == nil) || (m1 != nil && !m1.IsEquivalentTo(m2)) {
				return false
			}
		}
		return true
	case jsDocPropKindType:
		return jsdocGet[*JSTypeExpression](left).IsEquivalentTo(jsdocGet[*JSTypeExpression](right))
	case jsDocPropKindTypeList:
		l := jsdocGet[*jutil.ArrayList[*JSTypeExpression]](left)
		r := jsdocGet[*jutil.ArrayList[*JSTypeExpression]](right)
		if l.Size() != r.Size() {
			return false
		}
		for i := 0; i < l.Size(); i++ {
			if !l.Get(i).IsEquivalentTo(r.Get(i)) {
				return false
			}
		}
		return true
	case jsDocPropKindTypeMap:
		l := jsdocGet[*jutil.LinkedHashMap[string, *JSTypeExpression]](left)
		r := jsdocGet[*jutil.LinkedHashMap[string, *JSTypeExpression]](right)
		leftKeys := jutil.NewLinkedHashSet(l.KeySlice()...)
		rightKeys := jutil.NewLinkedHashSet(r.KeySlice()...)
		if !leftKeys.Equals(rightKeys) {
			return false
		}
		for key := range leftKeys.All() {
			if !jsDocInfo_TypeMapProperty_AreEquivalent(l.Get(key), r.Get(key)) {
				return false
			}
		}
		return true
	case jsDocPropKindNodeMap:
		l := jsdocGet[*jutil.LinkedHashMap[string, *Node]](left)
		r := jsdocGet[*jutil.LinkedHashMap[string, *Node]](right)
		if l.Size() != r.Size() {
			return false
		}
		for k, v := range l.All() {
			rightValue := r.Get(k)
			if rightValue == nil || !v.IsEquivalentTo(rightValue) {
				return false
			}
		}
		return true
	}
	return jutil.ObjectsEquals(left, right)
}

// Java: JSDocInfo.TypeMapProperty.areEquivalent(JSTypeExpression, JSTypeExpression)
func jsDocInfo_TypeMapProperty_AreEquivalent(expr1 *JSTypeExpression, expr2 *JSTypeExpression) bool {
	return jutil.ObjectsEquals(expr1, expr2) || (expr1 != nil && expr1.IsEquivalentTo(expr2))
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
	value := jsdocGet[*jutil.ArrayList[*JSTypeExpression]](p.Get(info))
	if value != nil {
		return value.Slice()
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

// The `name` and `mask` fields of each Bit, computed by the Bit
// constructor.
var jsdocInfo_Bit_name, jsdocInfo_Bit_mask = func() ([]string, []int64) {
	var names []string
	var masks []int64
	for _, b := range JSDocInfo_Bit_values() {
		// Java: JSDocInfo.Bit()
		if int(b) > 63 {
			panic(jutil.NewAssertionError("Too many Bits"))
		}
		names = append(names, jutil.CaseFormat_UpperUnderscoreToLowerCamel(b.Name()))
		masks = append(masks, jutil.Shl64(1, int64(b)))
	}
	return names, masks
}()

// Java: JSDocInfo.Bit.name
func (b JSDocInfo_Bit) BitName() string { return jsdocInfo_Bit_name[b] }

// Java: JSDocInfo.Bit.mask
func (b JSDocInfo_Bit) Mask() int64 { return jsdocInfo_Bit_mask[b] }

// JSDocInfo_TypeTransform transforms a type expression.
//
// Go-only: language-adapter — a functional interface, so a Go func type.
//
// Java: JSDocInfo.TypeTransform.apply(JSTypeExpression)
type JSDocInfo_TypeTransform func(arg *JSTypeExpression) *JSTypeExpression

// Java: JSDocInfo(long, TreeMap<Property<?>, Object>)
func newJSDocInfo(bits_ int64, props *jutil.TreeMap[*JSDocInfo_Property, any]) *JSDocInfo {
	var keys int64 = 0
	var valsBuf [8]any
	values := valsBuf[:0]
	for prop, value := range props.All() {
		if !prop.IsDefault(value) {
			keys |= prop.mask

			// Ensure that another codepath did not inadvertently retrieve a
			// reference to the array or a `null` value.
			_, isArray := value.([]any)
			jutil.CheckState(!isArray && !jutil.IsNil(value))
			values = append(values, value)
		}
	}
	return &JSDocInfo{
		propertyBits:       bits_,
		propertyKeysBitset: keys,
		propertyValues:     JSDocInfo_PackPropertyValues(values),
	}
}

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
			jutil.CheckState(!isArray && !jutil.IsNil(value))
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

func jsdocInfo_NewPropMap() *jutil.TreeMap[*JSDocInfo_Property, any] {
	return jutil.NewTreeMap[*JSDocInfo_Property, any](func(a, b *JSDocInfo_Property) int { return a.CompareTo(b) })
}

// Java: JSDocInfo.asMap()
func (i *JSDocInfo) AsMap() *jutil.TreeMap[*JSDocInfo_Property, any] {
	m := jsdocInfo_NewPropMap()
	bits_ := i.propertyKeysBitset
	index := 0
	for bits_ > 0 {
		low := bits.TrailingZeros64(uint64(bits_))
		bits_ &^= jutil.Shl64(1, int64(low))
		m.Put(JSDocInfo_Property_values[low], i.GetPropertyValueByIndex(index))
		index++
	}
	return m
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

// JSDocInfo_Visibility_values is Visibility.values().
func JSDocInfo_Visibility_values() []JSDocInfo_Visibility {
	return []JSDocInfo_Visibility{
		JSDocInfo_Visibility_PRIVATE, JSDocInfo_Visibility_PACKAGE, JSDocInfo_Visibility_PROTECTED,
		JSDocInfo_Visibility_PUBLIC, JSDocInfo_Visibility_INHERITED,
	}
}

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
	JSDocInfo_THROWS_ANNOTATIONS        = newJSDocInfo_Property("throws", jsDocPropKindPlain)
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
	JSDocInfo_SUPPRESSIONS         = newJSDocInfo_Property("suppressions", jsDocPropKindPlain)
	JSDocInfo_MODIFIES             = newJSDocInfo_Property("modifies", jsDocPropKindPlain)
	JSDocInfo_LENDS_NAME           = newJSDocInfo_Property("lendsName", jsDocPropKindType)
	JSDocInfo_CLOSURE_PRIMITIVE_ID = newJSDocInfo_Property("closurePrimitiveId", jsDocPropKindPlain)

	// NOTE: The following properties are "documentation properties", which do
	// _not_ get deeply copied when JSDocInfo is cloned (i.e. none of these
	// Properties override Property#clone).
	JSDocInfo_SOURCE_COMMENT               = newJSDocInfo_Property("sourceComment", jsDocPropKindPlain)
	JSDocInfo_MARKERS                      = newJSDocInfo_Property("markers", jsDocPropKindMarkerList)
	JSDocInfo_PARAMETER_DESCRIPTIONS       = newJSDocInfo_Property("parameterDescriptions", jsDocPropKindPlain)
	JSDocInfo_BLOCK_DESCRIPTION            = newJSDocInfo_Property("blockDescription", jsDocPropKindPlain)
	JSDocInfo_FILEOVERVIEW_DESCRIPTION     = newJSDocInfo_Property("fileoverviewDescription", jsDocPropKindPlain)
	JSDocInfo_RETURN_DESCRIPTION           = newJSDocInfo_Property("returnDescription", jsDocPropKindPlain)
	JSDocInfo_ENHANCED_NAMESPACE           = newJSDocInfo_Property("enhance", jsDocPropKindPlain)
	JSDocInfo_MODS                         = newJSDocInfo_Property("mods", jsDocPropKindPlain)
	JSDocInfo_TS_TYPES                     = newJSDocInfo_Property("tsType", jsDocPropKindPlain)
	JSDocInfo_CLOSURE_UNAWARE_CONFIG_VALUE = newJSDocInfo_Property("perFileClosureUnawareMode", jsDocPropKindPlain)

	JSDocInfo_AUTHORS = newJSDocInfo_Property("authors", jsDocPropKindPlain)
	JSDocInfo_SEES    = newJSDocInfo_Property("sees", jsDocPropKindPlain)
)

// JSDocInfo_ComparableSourcePosition is a SourcePosition that can be
// compared for equivalence.
//
// Java: JSDocInfo.ComparableSourcePosition<T>
type JSDocInfo_ComparableSourcePosition[T any] interface {
	SourcePosition[T]
	IsEquivalentTo(that SourcePosition[T]) bool
}

// JSDocInfo_StringPosition is a piece of information (found in a marker)
// which contains a position with a string.
//
// Go-only: language-adapter — StringPosition has a subclass
// (TrimmedStringPosition), so it is an interface plus
// JSDocInfo_StringPositionBase; the item type String is Go string (a
// position's item is always set right after construction).
type JSDocInfo_StringPosition interface {
	JSDocInfo_ComparableSourcePosition[string]
	GetJSDocInfo_StringPositionBase() *JSDocInfo_StringPositionBase
}

// JSDocInfo_StringPositionBase is `new StringPosition()`.
type JSDocInfo_StringPositionBase struct {
	SourcePositionBase[string]
}

// NewJSDocInfo_StringPosition is the implicit StringPosition constructor.
//
// Java: JSDocInfo.StringPosition()
func NewJSDocInfo_StringPosition() *JSDocInfo_StringPositionBase {
	return &JSDocInfo_StringPositionBase{}
}

// GetJSDocInfo_StringPositionBase gives access to the base through the
// interface.
func (p *JSDocInfo_StringPositionBase) GetJSDocInfo_StringPositionBase() *JSDocInfo_StringPositionBase {
	return p
}

// Java: JSDocInfo.StringPosition.isEquivalentTo(SourcePosition<String>)
func (p *JSDocInfo_StringPositionBase) IsEquivalentTo(that SourcePosition[string]) bool {
	return !jutil.IsNil(that) && p.IsSamePositionAs(that) && jutil.ObjectsEquals(p.GetItem(), that.GetItem())
}

// JSDocInfo_TrimmedStringPosition is a piece of information (found in a
// marker) which contains a position with a string that has no leading or
// trailing whitespace.
type JSDocInfo_TrimmedStringPosition struct {
	JSDocInfo_StringPositionBase
}

// Java: JSDocInfo.TrimmedStringPosition()
func NewJSDocInfo_TrimmedStringPosition() *JSDocInfo_TrimmedStringPosition {
	return &JSDocInfo_TrimmedStringPosition{}
}

// Java: JSDocInfo.TrimmedStringPosition.setItem(String)
func (p *JSDocInfo_TrimmedStringPosition) SetItem(item string) {
	jutil.CheckArgumentMsg(
		jutil.StringCharAt(item, 0) != ' ' && jutil.StringCharAt(item, jutil.StringLength(item)-1) != ' ',
		"String has leading or trailing whitespace")
	p.JSDocInfo_StringPositionBase.SetItem(item)
}

// JSDocInfo_NamePosition is a piece of information (found in a marker) which
// contains a position with a name node.
type JSDocInfo_NamePosition struct {
	SourcePositionBase[*Node]
}

// Java: JSDocInfo.NamePosition()
func NewJSDocInfo_NamePosition() *JSDocInfo_NamePosition {
	return &JSDocInfo_NamePosition{}
}

// Java: JSDocInfo.NamePosition.isEquivalentTo(SourcePosition<Node>)
func (p *JSDocInfo_NamePosition) IsEquivalentTo(that SourcePosition[*Node]) bool {
	if jutil.IsNil(that) ||
		!p.IsSamePositionAs(that) ||
		(p.GetItem() == nil) != (that.GetItem() == nil) {
		return false
	}
	return p.GetItem() == nil || p.GetItem().IsEquivalentTo(that.GetItem())
}

// JSDocInfo_TypePosition is a piece of information (found in a marker) which
// contains a position with a type expression syntax tree.
type JSDocInfo_TypePosition struct {
	SourcePositionBase[*Node]
	brackets bool
}

// Java: JSDocInfo.TypePosition()
func NewJSDocInfo_TypePosition() *JSDocInfo_TypePosition {
	return &JSDocInfo_TypePosition{brackets: false}
}

// HasBrackets returns whether the type has curly braces around it.
//
// Java: JSDocInfo.TypePosition.hasBrackets()
func (p *JSDocInfo_TypePosition) HasBrackets() bool {
	return p.brackets
}

// Java: JSDocInfo.TypePosition.setHasBrackets(boolean)
func (p *JSDocInfo_TypePosition) SetHasBrackets(newVal bool) {
	p.brackets = newVal
}

// Java: JSDocInfo.TypePosition.isEquivalentTo(SourcePosition<Node>)
func (p *JSDocInfo_TypePosition) IsEquivalentTo(that SourcePosition[*Node]) bool {
	typePosition, ok := that.(*JSDocInfo_TypePosition)
	if !ok || typePosition == nil ||
		!p.IsSamePositionAs(that) ||
		p.brackets != typePosition.brackets ||
		(p.GetItem() == nil) != (that.GetItem() == nil) {
		return false
	}
	return p.GetItem() == nil || p.GetItem().IsEquivalentTo(that.GetItem())
}

// JSDocInfo_Marker defines a class for containing the parsing information
// for this JSDocInfo. For each annotation found in the JsDoc, a marker will
// be created indicating the annotation itself, the name of the annotation
// (if any; for example, a @param has a name, but a @return does not), the
// textual description found on that annotation and, if applicable, the type
// declaration. All this information is only collected if documentation
// collection is turned on.
type JSDocInfo_Marker struct {
	annotation  *JSDocInfo_TrimmedStringPosition
	nameNode    *JSDocInfo_NamePosition
	description JSDocInfo_StringPosition
	type_       *JSDocInfo_TypePosition
}

// Java: JSDocInfo.Marker()
func NewJSDocInfo_Marker() *JSDocInfo_Marker { return &JSDocInfo_Marker{} }

// GetAnnotation gets the position information for the annotation name.
// (e.g., "param")
//
// Java: JSDocInfo.Marker.getAnnotation()
func (m *JSDocInfo_Marker) GetAnnotation() JSDocInfo_StringPosition {
	if m.annotation == nil {
		return nil
	}
	return m.annotation
}

// Java: JSDocInfo.Marker.setAnnotation(TrimmedStringPosition)
func (m *JSDocInfo_Marker) SetAnnotation(p *JSDocInfo_TrimmedStringPosition) {
	m.annotation = p
}

// GetNameNode gets the position information for the name found in a @param
// tag.
//
// Java: JSDocInfo.Marker.getNameNode()
func (m *JSDocInfo_Marker) GetNameNode() *JSDocInfo_NamePosition {
	return m.nameNode
}

// Java: JSDocInfo.Marker.setNameNode(NamePosition)
func (m *JSDocInfo_Marker) SetNameNode(p *JSDocInfo_NamePosition) {
	m.nameNode = p
}

// GetDescription gets the position information for the description found in
// a block tag.
//
// Java: JSDocInfo.Marker.getDescription()
func (m *JSDocInfo_Marker) GetDescription() JSDocInfo_StringPosition {
	return m.description
}

// Java: JSDocInfo.Marker.setDescription(StringPosition)
func (m *JSDocInfo_Marker) SetDescription(p JSDocInfo_StringPosition) {
	if jutil.IsNil(p) {
		m.description = nil
	} else {
		m.description = p
	}
}

// GetType gets the position information for the type expression found in
// some block tags, like "@param" and "@return".
//
// Java: JSDocInfo.Marker.getType()
func (m *JSDocInfo_Marker) GetType() *JSDocInfo_TypePosition {
	return m.type_
}

// Java: JSDocInfo.Marker.setType(TypePosition)
func (m *JSDocInfo_Marker) SetType(p *JSDocInfo_TypePosition) {
	m.type_ = p
}

// Java: JSDocInfo.Marker.isEquivalentTo(Marker)
func (m *JSDocInfo_Marker) IsEquivalentTo(that *JSDocInfo_Marker) bool {
	return JSDocInfo_Marker_AreEquivalent[string](m.annotation, that.annotation) &&
		JSDocInfo_Marker_AreEquivalent[*Node](m.nameNode, that.nameNode) &&
		JSDocInfo_Marker_AreEquivalent[string](m.description, that.description) &&
		JSDocInfo_Marker_AreEquivalent[*Node](m.type_, that.type_)
}

// JSDocInfo_Marker_AreEquivalent compares two nullable positions.
//
// Go-only: language-adapter — a nil pointer converted to the interface
// parameter is Java null (jutil.IsNil).
//
// Java: JSDocInfo.Marker.areEquivalent(ComparableSourcePosition<T>, ComparableSourcePosition<T>)
func JSDocInfo_Marker_AreEquivalent[T any](p1 JSDocInfo_ComparableSourcePosition[T], p2 JSDocInfo_ComparableSourcePosition[T]) bool {
	p1Null, p2Null := jutil.IsNil(p1), jutil.IsNil(p2)
	if p2Null {
		p2 = nil
	}
	return p1Null == p2Null && (p1Null || p1.IsEquivalentTo(p2))
}

// String is java.lang.Object.toString() (Marker does not override it).
//
// Go-only: language-adapter — Java prints an identity hash code, which
// differs from run to run; this one is derived from the object's address.
func (m *JSDocInfo_Marker) String() string {
	return jutil.ObjectToString("com.google.javascript.rhino.JSDocInfo$Marker", m)
}

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
		bits_ &^= jutil.Shl64(1, int64(low))
		prop := JSDocInfo_Property_values[low]
		a := jsDoc1.GetPropertyValueByIndex(index)
		b := jsDoc2.GetPropertyValueByIndex(index)
		index++
		if jutil.IsNil(a) != jutil.IsNil(b) || (!jutil.IsNil(a) && !prop.EqualValues(a, b)) {
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
	v := JSDocInfo_ID_GENERATOR.Get(i)
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
	type_ := jsdocGet[*JSTypeExpression](JSDocInfo_TYPE.Get(i))
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
	entireComment := i.GetOriginalCommentString()
	if entireComment == nil {
		return false
	}
	return strings.Contains(*entireComment, "@code")
}

// GetVisibility gets the visibility specified by `@private`, `@protected` or
// `@public` annotation. If no visibility is specified, visibility is
// inherited from the base class.
//
// Java: JSDocInfo.getVisibility()
func (i *JSDocInfo) GetVisibility() JSDocInfo_Visibility {
	visibility := JSDocInfo_VISIBILITY.Get(i)
	if visibility != nil {
		return jsdocGet[JSDocInfo_Visibility](visibility)
	}
	return JSDocInfo_Visibility_INHERITED
}

func (i *JSDocInfo) parameters() *jutil.LinkedHashMap[string, *JSTypeExpression] {
	return jsdocGet[*jutil.LinkedHashMap[string, *JSTypeExpression]](JSDocInfo_PARAMETERS.Get(i))
}

// GetParameterType gets the type of a given named parameter; null if this
// parameter is not defined or has a null type.
//
// Java: JSDocInfo.getParameterType(String)
func (i *JSDocInfo) GetParameterType(parameter string) *JSTypeExpression {
	params := i.parameters()
	if params != nil {
		return params.Get(parameter)
	}
	return nil
}

// HasParameter returns whether the parameter is defined.
//
// Java: JSDocInfo.hasParameter(String)
func (i *JSDocInfo) HasParameter(parameter string) bool {
	params := i.parameters()
	return params != nil && params.ContainsKey(parameter)
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
	params := i.parameters()
	if params != nil && !params.IsEmpty() {
		set := jutil.NewLinkedHashSet[string]()
		for name := range params.Keys() {
			set.Add(name)
		}
		return set
	}
	return jsDocInfo_noSuppressions
}

// ParameterNames iterates the defined parameter names in JSDoc definition order without allocating a set.
func (i *JSDocInfo) ParameterNames() iter.Seq[string] {
	params := i.parameters()
	if params == nil {
		return func(func(string) bool) {}
	}
	return params.Keys()
}

// ParameterTypes iterates the defined parameter type expressions in JSDoc definition order without allocating a set.
func (i *JSDocInfo) ParameterTypes() iter.Seq[*JSTypeExpression] {
	params := i.parameters()
	if params == nil {
		return func(func(*JSTypeExpression) bool) {}
	}
	return params.Values()
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
	params := i.parameters()
	if params == nil {
		return "", false
	}
	return params.KeyAt(index)
}

// GetParameterTypeAt returns the nth defined parameter type expression in JSDoc
// definition order without allocating an iterator closure.
func (i *JSDocInfo) GetParameterTypeAt(index int) *JSTypeExpression {
	params := i.parameters()
	if params == nil {
		return nil
	}
	t, _ := params.ValueAt(index)
	return t
}

// GetParameterCount gets the number of parameters defined.
//
// Java: JSDocInfo.getParameterCount()
func (i *JSDocInfo) GetParameterCount() int {
	params := i.parameters()
	if params != nil {
		return params.Size()
	}
	return 0
}

// GetThrowsAnnotations returns the list of thrown types and descriptions as
// text.
//
// Java: JSDocInfo.getThrowsAnnotations()
func (i *JSDocInfo) GetThrowsAnnotations() []string {
	annotations := jsdocGet[*jutil.ArrayList[string]](JSDocInfo_THROWS_ANNOTATIONS.Get(i))
	if annotations != nil {
		return annotations.Slice()
	}
	return []string{}
}

// HasEnumParameterType returns whether an enum parameter type, specified
// using the `@enum` annotation, is present on this JSDoc.
//
// Java: JSDocInfo.hasEnumParameterType()
func (i *JSDocInfo) HasEnumParameterType() bool {
	return JSDocInfo_ENUM_PARAMETER_TYPE.Get(i) != nil
}

// HasTypedefType returns whether a typedef parameter type, specified using
// the `@typedef` annotation, is present on this JSDoc.
//
// Java: JSDocInfo.hasTypedefType()
func (i *JSDocInfo) HasTypedefType() bool {
	return JSDocInfo_TYPEDEF_TYPE.Get(i) != nil
}

// HasReturnType returns whether this JSDocInfo contains a type for `@return`
// annotation.
//
// Java: JSDocInfo.hasReturnType()
func (i *JSDocInfo) HasReturnType() bool {
	return JSDocInfo_RETURN_TYPE.Get(i) != nil
}

// HasType returns whether a type, specified using the `@type` annotation, is
// present on this JSDoc.
//
// Java: JSDocInfo.hasType()
func (i *JSDocInfo) HasType() bool {
	return JSDocInfo_TYPE.Get(i) != nil
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
	return jsdocGet[*JSTypeExpression](JSDocInfo_RETURN_TYPE.Get(i))
}

// GetEnumParameterType gets the enum parameter type specified by the `@enum`
// annotation.
//
// Java: JSDocInfo.getEnumParameterType()
func (i *JSDocInfo) GetEnumParameterType() *JSTypeExpression {
	return jsdocGet[*JSTypeExpression](JSDocInfo_ENUM_PARAMETER_TYPE.Get(i))
}

// GetTypedefType gets the typedef type specified by the `@type` annotation.
//
// Java: JSDocInfo.getTypedefType()
func (i *JSDocInfo) GetTypedefType() *JSTypeExpression {
	return jsdocGet[*JSTypeExpression](JSDocInfo_TYPEDEF_TYPE.Get(i))
}

// GetType gets the type specified by the `@type` annotation.
//
// Java: JSDocInfo.getType()
func (i *JSDocInfo) GetType() *JSTypeExpression {
	return jsdocGet[*JSTypeExpression](JSDocInfo_TYPE.Get(i))
}

// GetThisType gets the type specified by the `@this` annotation.
//
// Java: JSDocInfo.getThisType()
func (i *JSDocInfo) GetThisType() *JSTypeExpression {
	return jsdocGet[*JSTypeExpression](JSDocInfo_THIS_TYPE.Get(i))
}

// HasThisType returns whether this JSDocInfo contains a type for `@this`
// annotation.
//
// Java: JSDocInfo.hasThisType()
func (i *JSDocInfo) HasThisType() bool {
	return JSDocInfo_THIS_TYPE.Get(i) != nil
}

// GetBaseType gets the base type specified by the `@extends` annotation.
//
// Java: JSDocInfo.getBaseType()
func (i *JSDocInfo) GetBaseType() *JSTypeExpression {
	return jsdocGet[*JSTypeExpression](JSDocInfo_BASE_TYPE.Get(i))
}

func jsdocString(v any) *string {
	if v == nil {
		return nil
	}
	s := jsdocGet[string](v)
	return &s
}

// GetDescription gets the description specified by the `@desc` annotation.
//
// Java: JSDocInfo.getDescription()
func (i *JSDocInfo) GetDescription() *string {
	return jsdocString(JSDocInfo_DESCRIPTION.Get(i))
}

// GetTsTypes gets the ts type declarations specified by the `@tsType`
// annotations.
//
// Java: JSDocInfo.getTsTypes()
func (i *JSDocInfo) GetTsTypes() []string {
	types := jsdocGet[*jutil.ArrayList[string]](JSDocInfo_TS_TYPES.Get(i))
	if types != nil {
		return types.ToSlice()
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
	return jsdocString(JSDocInfo_MEANING.Get(i))
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
	return jsdocString(JSDocInfo_ALTERNATE_MESSAGE_ID.Get(i))
}

// GetLendsName gets the name we're lending to in a `@lends` annotation.
//
// In many reflection APIs, you pass an anonymous object to a function, and
// that function mixes the anonymous object into another object. The `@lends`
// annotation allows the type system to track those property assignments.
//
// Java: JSDocInfo.getLendsName()
func (i *JSDocInfo) GetLendsName() *JSTypeExpression {
	return jsdocGet[*JSTypeExpression](JSDocInfo_LENDS_NAME.Get(i))
}

// Java: JSDocInfo.hasLendsName()
func (i *JSDocInfo) HasLendsName() bool {
	return i.GetLendsName() != nil
}

// GetClosurePrimitiveId returns the `@closurePrimitive {id}` identifier.
//
// Java: JSDocInfo.getClosurePrimitiveId()
func (i *JSDocInfo) GetClosurePrimitiveId() *string {
	return jsdocString(JSDocInfo_CLOSURE_PRIMITIVE_ID.Get(i))
}

// HasClosurePrimitiveId reports whether this JSDoc is annotated with
// `@closurePrimitive`.
//
// Java: JSDocInfo.hasClosurePrimitiveId()
func (i *JSDocInfo) HasClosurePrimitiveId() bool {
	return JSDocInfo_CLOSURE_PRIMITIVE_ID.Get(i) != nil
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
	return JSDocInfo_CLOSURE_UNAWARE_CONFIG_VALUE.Get(i) != nil
}

// GetPerFileClosureUnawareMode returns the value of the
// `@perFileClosureUnawareMode` annotation, or null if not present.
//
// Java: JSDocInfo.getPerFileClosureUnawareMode()
func (i *JSDocInfo) GetPerFileClosureUnawareMode() *JSDocInfo_PerFileClosureUnawareMode {
	v := JSDocInfo_CLOSURE_UNAWARE_CONFIG_VALUE.Get(i)
	if v == nil {
		return nil
	}
	m := jsdocGet[JSDocInfo_PerFileClosureUnawareMode](v)
	return &m
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
	return jsdocString(JSDocInfo_LICENSE.Get(i))
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

	for _, b := range JSDocInfo_Bit_values() {
		if i.CheckBit(b) {
			helper.AddBool("bit:"+b.BitName(), true)
		}
	}

	bits_ := i.propertyKeysBitset
	index := 0
	for bits_ > 0 {
		low := bits.TrailingZeros64(uint64(bits_))
		bits_ &^= jutil.Shl64(1, int64(low))
		helper = helper.Add(JSDocInfo_Property_values[low].name, i.GetPropertyValueByIndex(index))
		index++
	}

	return helper.OmitNullValues().String()
}

// HasBaseType returns whether this JSDocInfo contains a type for `@extends`
// annotation.
//
// Java: JSDocInfo.hasBaseType()
func (i *JSDocInfo) HasBaseType() bool {
	return i.GetBaseType() != nil
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
	list := jsdocGet[*jutil.ArrayList[*JSTypeExpression]](JSDocInfo_IMPLEMENTED_INTERFACES.Get(i))
	if list != nil {
		return list.Size()
	}
	return 0
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
	list := jsdocGet[*jutil.ArrayList[*JSTypeExpression]](JSDocInfo_EXTENDED_INTERFACES.Get(i))
	if list != nil {
		return list.Size()
	}
	return 0
}

// GetDeprecationReason returns the deprecation reason or null if none
// specified.
//
// Java: JSDocInfo.getDeprecationReason()
func (i *JSDocInfo) GetDeprecationReason() *string {
	return jsdocString(JSDocInfo_DEPRECATION_REASON.Get(i))
}

func (i *JSDocInfo) suppressions() *jutil.ObjectLinkedHashMap[*jutil.LinkedHashSet[string], string] {
	return jsdocGet[*jutil.ObjectLinkedHashMap[*jutil.LinkedHashSet[string], string]](JSDocInfo_SUPPRESSIONS.Get(i))
}

// GetSuppressions returns the set of suppressed warnings.
//
// Java: JSDocInfo.getSuppressions()
func (i *JSDocInfo) GetSuppressions() *jutil.LinkedHashSet[string] {
	suppressions := i.suppressions()
	if suppressions == nil {
		// Go-only: performance — ImmutableSet.of() is a shared singleton;
		// callers only read the result, so share one empty set instead of
		// allocating one per call (this runs for every JSDoc, externs
		// included).
		return jsDocInfo_noSuppressions
	}
	if suppressions.Size() == 1 {
		if nameSet, ok := suppressions.FirstKey(); ok {
			return nameSet
		}
	}
	names := jutil.NewLinkedHashSet[string]()
	for nameSet := range suppressions.Keys() {
		for name := range nameSet.All() {
			names.Add(name)
		}
	}
	// Go-only: performance — ImmutableSet.copyOf(names) would be a copy
	// with the same elements in the same order; names is not shared, so it
	// is returned as is.
	return names
}

// jsDocInfo_noSuppressions is GetSuppressions' result for a JSDocInfo
// without suppressions (ImmutableSet.of()). Never modified.
var jsDocInfo_noSuppressions = jutil.NewLinkedHashSet[string]()

// GetSuppressionsAndTheirDescription returns a map containing key=set of
// suppressions, and value=the corresponding description for the set.
//
// Java: JSDocInfo.getSuppressionsAndTheirDescription()
func (i *JSDocInfo) GetSuppressionsAndTheirDescription() *jutil.ObjectLinkedHashMap[*jutil.LinkedHashSet[string], string] {
	suppressions := i.suppressions()
	if suppressions != nil {
		return suppressions
	}
	return jutil.NewObjectLinkedHashMap[*jutil.LinkedHashSet[string], string]()
}

// GetModifies returns the set of sideeffect notations.
//
// Java: JSDocInfo.getModifies()
func (i *JSDocInfo) GetModifies() *jutil.LinkedHashSet[string] {
	modifies := jsdocGet[*jutil.LinkedHashSet[string]](JSDocInfo_MODIFIES.Get(i))
	if modifies != nil {
		return modifies
	}
	return jsDocInfo_noSuppressions
}

func (i *JSDocInfo) parameterDescriptions() *jutil.LinkedHashMap[string, string] {
	return jsdocGet[*jutil.LinkedHashMap[string, string]](JSDocInfo_PARAMETER_DESCRIPTIONS.Get(i))
}

// HasDescriptionForParameter returns whether a description exists for the
// parameter with the specified name.
//
// Java: JSDocInfo.hasDescriptionForParameter(String)
func (i *JSDocInfo) HasDescriptionForParameter(name string) bool {
	params := i.parameterDescriptions()
	return params != nil && params.ContainsKey(name)
}

// GetDescriptionForParameter returns the description for the parameter with
// the given name, if its exists.
//
// Java: JSDocInfo.getDescriptionForParameter(String)
func (i *JSDocInfo) GetDescriptionForParameter(name string) *string {
	params := i.parameterDescriptions()
	if params != nil {
		if d, ok := params.Lookup(name); ok {
			return &d
		}
	}
	return nil
}

// GetAuthors returns the list of authors or null if none.
//
// Java: JSDocInfo.getAuthors()
func (i *JSDocInfo) GetAuthors() *jutil.ArrayList[string] {
	return jsdocGet[*jutil.ArrayList[string]](JSDocInfo_AUTHORS.Get(i))
}

// GetReferences returns the list of references or null if none.
//
// Java: JSDocInfo.getReferences()
func (i *JSDocInfo) GetReferences() *jutil.ArrayList[string] {
	return jsdocGet[*jutil.ArrayList[string]](JSDocInfo_SEES.Get(i))
}

// GetReturnDescription returns the description of the returned object or
// null if none specified.
//
// Java: JSDocInfo.getReturnDescription()
func (i *JSDocInfo) GetReturnDescription() *string {
	return jsdocString(JSDocInfo_RETURN_DESCRIPTION.Get(i))
}

// GetBlockDescription returns the block-level description or null if none
// specified.
//
// Java: JSDocInfo.getBlockDescription()
func (i *JSDocInfo) GetBlockDescription() *string {
	return jsdocString(JSDocInfo_BLOCK_DESCRIPTION.Get(i))
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
	return jsdocString(JSDocInfo_FILEOVERVIEW_DESCRIPTION.Get(i))
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
	return jsdocString(JSDocInfo_ENHANCED_NAMESPACE.Get(i))
}

// HasMods returns whether this has a modded namespace.
//
// Java: JSDocInfo.hasMods()
func (i *JSDocInfo) HasMods() bool {
	return JSDocInfo_MODS.Get(i) != nil
}

// GetMods returns the modded namespace or null if none is specified.
//
// Java: JSDocInfo.getMods()
func (i *JSDocInfo) GetMods() *string {
	return jsdocString(JSDocInfo_MODS.Get(i))
}

// GetMarkers gets the list of all markers for the documentation in this
// JSDoc.
//
// Java: JSDocInfo.getMarkers()
func (i *JSDocInfo) GetMarkers() []*JSDocInfo_Marker {
	markers := jsdocGet[*jutil.ArrayList[*JSDocInfo_Marker]](JSDocInfo_MARKERS.Get(i))
	if markers != nil {
		return markers.Slice()
	}
	return []*JSDocInfo_Marker{}
}

func (i *JSDocInfo) templateTypeNames() *jutil.LinkedHashMap[string, *JSTypeExpression] {
	return jsdocGet[*jutil.LinkedHashMap[string, *JSTypeExpression]](JSDocInfo_TEMPLATE_TYPE_NAMES.Get(i))
}

// GetTemplateTypeNames gets the @template type names.
//
// Excludes @template types from TTL; get those with
// getTypeTransformations().
//
// Java: JSDocInfo.getTemplateTypeNames()
func (i *JSDocInfo) GetTemplateTypeNames() []string {
	m := i.templateTypeNames()
	if m != nil {
		return m.KeySlice()
	}
	return []string{}
}

var (
	jsDocInfo_noTemplateTypes       = jutil.NewLinkedHashMap[string, *JSTypeExpression]()
	jsDocInfo_noTypeTransformations = jutil.NewLinkedHashMap[string, *Node]()
)

// Java: JSDocInfo.getTemplateTypes()
func (i *JSDocInfo) GetTemplateTypes() *jutil.LinkedHashMap[string, *JSTypeExpression] {
	m := i.templateTypeNames()
	if m != nil {
		return m
	}
	return jsDocInfo_noTemplateTypes
}

// GetTypeTransformations gets the type transformations.
//
// Java: JSDocInfo.getTypeTransformations()
func (i *JSDocInfo) GetTypeTransformations() *jutil.LinkedHashMap[string, *Node] {
	m := jsdocGet[*jutil.LinkedHashMap[string, *Node]](JSDocInfo_TYPE_TRANSFORMATIONS.Get(i))
	if m != nil {
		return m
	}
	return jsDocInfo_noTypeTransformations
}

// GetTypeExpressions returns a collection of all JSTypeExpressions that are
// a part of this JSDocInfo.
//
// This includes:
//
//   - base type
//   - @extends
//   - @implements
//   - @lend
//   - @param
//   - @return
//   - @template
//   - @this
//   - @throws
//   - @type
//
// Any future type specific JSDoc should make sure to add the appropriate
// nodes here.
//
// JSDocTypeExprIterator is a stack-allocated, zero-closure cursor over all
// non-nil JSTypeExpressions in a JSDocInfo.
type JSDocTypeExprIterator struct {
	info     *JSDocInfo
	typeBits uint64
	curList  *jutil.ArrayList[*JSTypeExpression]
	curMap   *jutil.LinkedHashMap[string, *JSTypeExpression]
	subIdx   int
	cur      *JSTypeExpression
}

// TypeExprIter returns a stack-allocated JSDocTypeExprIterator for i.
func (i *JSDocInfo) TypeExprIter() JSDocTypeExprIterator {
	if i == nil {
		return JSDocTypeExprIterator{}
	}
	return JSDocTypeExprIterator{
		info:     i,
		typeBits: uint64(i.propertyKeysBitset & jsDocInfo_typePropertiesMask),
	}
}

// Expr returns the current JSTypeExpression after Next returns true.
func (it *JSDocTypeExprIterator) Expr() *JSTypeExpression { return it.cur }

// Node returns the root Node of the current JSTypeExpression after Next returns true.
func (it *JSDocTypeExprIterator) Node() *Node { return it.cur.GetRoot() }

// Next advances to the next non-nil JSTypeExpression and reports whether one exists.
func (it *JSDocTypeExprIterator) Next() bool {
	for {
		if it.curList != nil {
			n := it.curList.Size()
			for it.subIdx < n {
				e := it.curList.Get(it.subIdx)
				it.subIdx++
				if e != nil {
					it.cur = e
					return true
				}
			}
			it.curList = nil
		}
		if it.curMap != nil {
			for {
				v, nextSlot, ok := it.curMap.NextSlotValue(it.subIdx)
				if !ok {
					break
				}
				it.subIdx = nextSlot
				if v != nil {
					it.cur = v
					return true
				}
			}
			it.curMap = nil
		}
		if it.typeBits == 0 {
			it.cur = nil
			return false
		}
		low := bits.TrailingZeros64(it.typeBits)
		propMask := uint64(1) << low
		it.typeBits &^= propMask
		index := bits.OnesCount64(uint64(it.info.propertyKeysBitset) & (propMask - 1))
		val := it.info.GetPropertyValueByIndex(index)
		switch JSDocInfo_Property_values[low].kind {
		case jsDocPropKindType:
			if e := jsdocGet[*JSTypeExpression](val); e != nil {
				it.cur = e
				return true
			}
		case jsDocPropKindTypeList:
			it.curList = jsdocGet[*jutil.ArrayList[*JSTypeExpression]](val)
			it.subIdx = 0
		case jsDocPropKindTypeMap:
			it.curMap = jsdocGet[*jutil.LinkedHashMap[string, *JSTypeExpression]](val)
			it.subIdx = 0
		}
	}
}

// TypeExpressions iterates over all non-nil JSTypeExpressions in this JSDocInfo
// without allocating intermediate slices.
func (i *JSDocInfo) TypeExpressions() iter.Seq[*JSTypeExpression] {
	return func(yield func(*JSTypeExpression) bool) {
		for it := i.TypeExprIter(); it.Next(); {
			if !yield(it.Expr()) {
				return
			}
		}
	}
}

// TypeNodes iterates over all type AST root nodes in this JSDocInfo without
// allocating intermediate slices.
func (i *JSDocInfo) TypeNodes() iter.Seq[*Node] {
	return func(yield func(*Node) bool) {
		for it := i.TypeExprIter(); it.Next(); {
			if !yield(it.Node()) {
				return
			}
		}
	}
}

// HasTypeNodes reports whether this JSDocInfo contains at least one type node,
// without allocating slices or closures.
func (i *JSDocInfo) HasTypeNodes() bool {
	it := i.TypeExprIter()
	return it.Next()
}

// Java: JSDocInfo.getTypeExpressions()
func (i *JSDocInfo) GetTypeExpressions() []*JSTypeExpression {
	builder := []*JSTypeExpression{}
	for it := i.TypeExprIter(); it.Next(); {
		builder = append(builder, it.Expr())
	}
	return builder
}

// GetTypeNodes returns a collection of all type nodes that are a part of
// this JSDocInfo.
//
// This includes:
//
//   - @extends
//   - @implements
//   - @lend
//   - @param
//   - @return
//   - @template
//   - @this
//   - @throws
//   - @type
//
// Any future type specific JSDoc should make sure to add the appropriate
// nodes here.
//
// Java: JSDocInfo.getTypeNodes()
func (i *JSDocInfo) GetTypeNodes() []*Node {
	nodes := []*Node{}

	for it := i.TypeExprIter(); it.Next(); {
		nodes = append(nodes, it.Node())
	}

	return nodes
}

// Java: JSDocInfo.hasModifies()
func (i *JSDocInfo) HasModifies() bool {
	return JSDocInfo_MODIFIES.Get(i) != nil
}

// GetOriginalCommentString returns the original JSDoc comment string.
// Returns null unless parseJsDocDocumentation is enabled via the
// ParserConfig.
//
// Java: JSDocInfo.getOriginalCommentString()
func (i *JSDocInfo) GetOriginalCommentString() *string {
	return jsdocString(JSDocInfo_SOURCE_COMMENT.Get(i))
}

// Java: JSDocInfo.getOriginalCommentPosition()
func (i *JSDocInfo) GetOriginalCommentPosition() int {
	pos := JSDocInfo_ORIGINAL_COMMENT_POSITION.Get(i)
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
	modifies := jsdocGet[*jutil.LinkedHashSet[string]](JSDocInfo_MODIFIES.Get(i))
	return modifies != nil && modifies.Contains("this")
}

// HasSideEffectsArgumentsAnnotation returns whether the @modifies annotation
// includes "arguments" or any named parameters.
//
// Java: JSDocInfo.hasSideEffectsArgumentsAnnotation()
func (i *JSDocInfo) HasSideEffectsArgumentsAnnotation() bool {
	modifies := jsdocGet[*jutil.LinkedHashSet[string]](JSDocInfo_MODIFIES.Get(i))
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
		n := bits.OnesCount64(uint64(i.propertyKeysBitset))
		if index < 0 || index >= n {
			panic(jutil.NewArrayIndexOutOfBoundsExceptionForLength(index, n))
		}
		return unsafe.Slice(i.propertyValues.(*any), n)[index]
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
	return &out[0]
}

// JSDocInfo_Builder is a builder for JSDocInfo objects. This builder is
// required because JSDocInfo instances have immutable structure. It provides
// early incompatibility detection among properties stored on the JSDocInfo
// object being created.
type JSDocInfo_Builder struct {
	propMask   int64
	entries    []any
	entriesBuf [4]any
	bits       int64

	// whether the current JSDocInfo has valuable information
	populated bool
	// the current marker, if any.
	currentMarker *JSDocInfo_Marker
	// the set of unique license texts
	licenseTexts *jutil.LinkedHashSet[string]
}

// NewJSDocInfo_Builder is the implicit Builder constructor.
//
// Java: JSDocInfo.Builder()
func NewJSDocInfo_Builder() *JSDocInfo_Builder {
	b := &JSDocInfo_Builder{}
	b.entries = b.entriesBuf[:0]
	return b
}

// Reset resets the builder for reuse without allocating.
func (b *JSDocInfo_Builder) Reset() {
	clear(b.entries)
	b.entries = b.entriesBuf[:0]
	b.propMask = 0
	b.bits = 0
	b.populated = false
	b.currentMarker = nil
	b.licenseTexts = nil
}

// Java: JSDocInfo.Builder.copyFrom(JSDocInfo)
func JSDocInfo_Builder_CopyFrom(info *JSDocInfo) *JSDocInfo_Builder {
	return info.ToBuilder()
}

// Java: JSDocInfo.Builder.maybeCopyFrom(JSDocInfo)
func JSDocInfo_Builder_MaybeCopyFrom(info *JSDocInfo) *JSDocInfo_Builder {
	if info != nil {
		return info.ToBuilder()
	}
	return JSDocInfo_Builder_().ParseDocumentation()
}

// JSDocInfo_Builder_MaybeCopyFromWithNewType returns a JSDocInfo.Builder
// that contains a copy of the given JSDocInfo in which only the `@type`
// field of the JSDocInfo is replaced with the given typeExpression. This is
// done to prevent generating code in the client module which references
// local variables from another module.
//
// Java: JSDocInfo.Builder.maybeCopyFromWithNewType(JSDocInfo, JSTypeExpression)
func JSDocInfo_Builder_MaybeCopyFromWithNewType(info *JSDocInfo, typeExpression *JSTypeExpression) *JSDocInfo_Builder {
	if info == nil {
		return JSDocInfo_Builder_().ParseDocumentation().SetType(typeExpression)
	}
	return info.ToBuilder().SetType(typeExpression)
}

// Java: JSDocInfo.Builder.copyFromWithNewType(JSDocInfo, JSTypeExpression)
func JSDocInfo_Builder_CopyFromWithNewType(info *JSDocInfo, typeExpression *JSTypeExpression) *JSDocInfo_Builder {
	return info.ToBuilder().SetType(typeExpression)
}

// JSDocInfo_Builder_MaybeCopyFromAndReplaceNames returns a JSDocInfo.Builder
// that contains a JSDoc in which all module local types (which may be inside
// `@param`, `@type` or `@returns` are replaced with unknown. This is done to
// prevent generating code in the client module which references local
// variables from another module.
//
// Java: JSDocInfo.Builder.maybeCopyFromAndReplaceNames(JSDocInfo, Set<String>)
func JSDocInfo_Builder_MaybeCopyFromAndReplaceNames(info *JSDocInfo, moduleLocalNamesToReplace *jutil.LinkedHashSet[string]) *JSDocInfo_Builder {
	if info != nil {
		return JSDocInfo_Builder_CopyFromAndReplaceNames(info, moduleLocalNamesToReplace)
	}
	return JSDocInfo_Builder_().ParseDocumentation()
}

// Java: JSDocInfo.Builder.copyFromAndReplaceNames(JSDocInfo, Set<String>)
func JSDocInfo_Builder_CopyFromAndReplaceNames(info *JSDocInfo, oldNames *jutil.LinkedHashSet[string]) *JSDocInfo_Builder {
	return info.CloneAndReplaceTypeNames(oldNames).ToBuilder() // TODO - populated
}

// ParseDocumentation configures the builder to parse documentation. This
// should be called immediately after instantiating the builder if
// documentation should be included, since it enables various operations to
// do work that would otherwise be no-ops.
//
// Java: JSDocInfo.Builder.parseDocumentation()
func (b *JSDocInfo_Builder) ParseDocumentation() *JSDocInfo_Builder {
	b.SetBit(JSDocInfo_Bit_INCLUDE_DOCUMENTATION, true)
	return b
}

// Java: JSDocInfo.Builder.shouldParseDocumentation()
func (b *JSDocInfo_Builder) ShouldParseDocumentation() bool {
	return b.CheckBit(JSDocInfo_Bit_INCLUDE_DOCUMENTATION)
}

// RecordOriginalCommentString sets the original JSDoc comment string. This
// is a no-op if the builder isn't configured to record documentation.
//
// Java: JSDocInfo.Builder.recordOriginalCommentString(String)
func (b *JSDocInfo_Builder) RecordOriginalCommentString(sourceComment string) {
	if b.ShouldParseDocumentation() {
		b.populated = true
		b.SetProp(JSDocInfo_SOURCE_COMMENT, sourceComment)
	}
}

// RecordOriginalCommentPosition sets the position of original JSDoc
// comment.
//
// Java: JSDocInfo.Builder.recordOriginalCommentPosition(int)
func (b *JSDocInfo_Builder) RecordOriginalCommentPosition(position int) {
	if b.ShouldParseDocumentation() {
		b.populated = true
		b.SetProp(JSDocInfo_ORIGINAL_COMMENT_POSITION, position)
	}
}

// IsPopulatedWithFileOverview returns whether this builder is populated with
// information that can be used to build a JSDocInfo object that has a
// fileoverview tag.
//
// Java: JSDocInfo.Builder.isPopulatedWithFileOverview()
func (b *JSDocInfo_Builder) IsPopulatedWithFileOverview() bool {
	return b.populated &&
		((b.bits&
			(JSDocInfo_Bit_FILEOVERVIEW.Mask()|
				JSDocInfo_Bit_EXTERNS.Mask()|
				JSDocInfo_Bit_NOCOMPILE.Mask()|
				JSDocInfo_Bit_NOCOVERAGE.Mask()|
				JSDocInfo_Bit_TYPE_SUMMARY.Mask()|
				JSDocInfo_Bit_ENHANCED_NAMESPACE.Mask())) != 0 ||
			b.IsModsRecorded())
}

// IsDescriptionRecorded returns whether this builder recorded a
// description.
//
// Java: JSDocInfo.Builder.isDescriptionRecorded()
func (b *JSDocInfo_Builder) IsDescriptionRecorded() bool {
	return b.GetProp(JSDocInfo_DESCRIPTION) != nil
}

// IsNoInline returns whether the `@noinline` annotation is present on this
// JSDocInfo.
//
// Java: JSDocInfo.Builder.isNoInline()
func (b *JSDocInfo_Builder) IsNoInline() bool {
	return b.CheckBit(JSDocInfo_Bit_NOINLINE)
}

// IsRequireInlining returns whether the `@requireInlining` annotation is
// present on this JSDocInfo.
//
// Java: JSDocInfo.Builder.isRequireInlining()
func (b *JSDocInfo_Builder) IsRequireInlining() bool {
	return b.CheckBit(JSDocInfo_Bit_REQUIRE_INLINING)
}

// IsEncourageInlining returns whether the `@encourageInlining` annotation is
// present on this JSDocInfo.
//
// Java: JSDocInfo.Builder.isEncourageInlining()
func (b *JSDocInfo_Builder) IsEncourageInlining() bool {
	return b.CheckBit(JSDocInfo_Bit_ENCOURAGE_INLINING)
}

// Build builds a JSDocInfo object based on the populated information and
// returns it. If no value was populated, this method simply returns null.
//
// Java: JSDocInfo.Builder.build()
func (b *JSDocInfo_Builder) Build() *JSDocInfo {
	return b.BuildWithAlways( /* always= */ false)
}

// BuildWithAlways builds a JSDocInfo object based on the populated
// information and returns it. If no value was populated and `always` is
// false, returns null. If `always` is true, returns a default JSDocInfo.
//
// Java: JSDocInfo.Builder.build(boolean)
func (b *JSDocInfo_Builder) BuildWithAlways(always bool) *JSDocInfo {
	if b.populated || always {
		info := newJSDocInfoFromBuilder(b.bits, b.propMask, b.entries)
		b.populated = false
		return info
	}
	return nil
}

// BuildAndReset builds a JSDocInfo object based on the populated information
// and returns it. Once this method is called, the builder can be reused to
// build another JSDocInfo object. If no value was populated, this method
// simply returns null.
//
// Java: JSDocInfo.Builder.buildAndReset()
func (b *JSDocInfo_Builder) BuildAndReset() *JSDocInfo {
	info := b.Build()
	b.bits &= JSDocInfo_Bit_INCLUDE_DOCUMENTATION.Mask() // only keep this one flag
	clear(b.entries)
	b.entries = b.entries[:0]
	b.propMask = 0
	b.populated = false
	return info
}

// MarkAnnotation adds a marker to the current JSDocInfo and populates the
// marker with the annotation information.
//
// Java: JSDocInfo.Builder.markAnnotation(String, int, int)
func (b *JSDocInfo_Builder) MarkAnnotation(annotation string, lineno int, charno int) {
	marker := b.AddMarker()
	if marker != nil {
		position := NewJSDocInfo_TrimmedStringPosition()
		position.SetItem(annotation)
		position.SetPositionInformation(lineno, charno, lineno, charno+jutil.StringLength(annotation))
		marker.SetAnnotation(position)
		b.populated = true
	}
	b.currentMarker = marker
}

// Java: JSDocInfo.Builder.addMarker()
func (b *JSDocInfo_Builder) AddMarker() *JSDocInfo_Marker {
	if b.ShouldParseDocumentation() {
		markers := jsdocGet[*jutil.ArrayList[*JSDocInfo_Marker]](b.GetProp(JSDocInfo_MARKERS))
		if markers == nil {
			markers = jutil.NewArrayList[*JSDocInfo_Marker]()
			b.SetProp(JSDocInfo_MARKERS, markers)
		}
		marker := NewJSDocInfo_Marker()
		markers.Add(marker)
		return marker
	}
	return nil
}

// MarkText adds a textual block to the current marker.
//
// Java: JSDocInfo.Builder.markText(String, int, int, int, int)
func (b *JSDocInfo_Builder) MarkText(text string, startLineno int, startCharno int, endLineno int, endCharno int) {
	if b.currentMarker != nil {
		position := NewJSDocInfo_StringPosition()
		position.SetItem(text)
		position.SetPositionInformation(startLineno, startCharno, endLineno, endCharno)
		b.currentMarker.SetDescription(position)
	}
}

// MarkTypeNode adds a type declaration to the current marker.
//
// Java: JSDocInfo.Builder.markTypeNode(Node, int, int, int, int, boolean)
func (b *JSDocInfo_Builder) MarkTypeNode(typeNode *Node, lineno int, startCharno int, endLineno int, endCharno int, hasLC bool) {
	if b.currentMarker != nil {
		position := NewJSDocInfo_TypePosition()
		position.SetItem(typeNode)
		position.SetHasBrackets(hasLC)
		position.SetPositionInformation(lineno, startCharno, endLineno, endCharno)
		b.currentMarker.SetType(position)
	}
}

// MarkName adds a name declaration to the current marker.
//
// Java: JSDocInfo.Builder.markName(String, Node, int, int)
func (b *JSDocInfo_Builder) MarkName(name string, templateNode *Node, lineno int, charno int) {
	if b.currentMarker != nil {
		// Record the name as both a SourcePosition<String> and a
		// SourcePosition<Node>. The <String> form is deprecated, because <Node>
		// is more consistent with how other name references are handled (see
		// #markTypeNode)
		//
		// TODO(nicksantos): Remove all uses of the Name position and replace
		// them with the NameNode position.
		position := NewJSDocInfo_TrimmedStringPosition()
		position.SetItem(name)
		position.SetPositionInformation(lineno, charno, lineno, charno+jutil.StringLength(name))

		nodePos := NewJSDocInfo_NamePosition()
		node := Node_NewStringWithToken(Token_NAME, name).SetLinenoCharno(lineno, charno)
		node.SetLength(jutil.StringLength(name))
		if templateNode != nil {
			node.SetStaticSourceFileFrom(templateNode)
		}
		nodePos.SetItem(node)
		nodePos.SetPositionInformation(lineno, charno, lineno, charno+jutil.StringLength(name))
		b.currentMarker.SetNameNode(nodePos)
	}
}

// RecordBlockDescription records a block-level description; returns true
// if the description was recorded.
//
// Java: JSDocInfo.Builder.recordBlockDescription(String)
func (b *JSDocInfo_Builder) RecordBlockDescription(description string) bool {
	b.populated = true
	if !b.ShouldParseDocumentation() {
		return true
	}
	return b.PopulateProp(JSDocInfo_BLOCK_DESCRIPTION, description)
}

// RecordVisibility records a visibility; returns true if the visibility was
// recorded and false if it was already defined.
//
// Java: JSDocInfo.Builder.recordVisibility(Visibility)
func (b *JSDocInfo_Builder) RecordVisibility(visibility JSDocInfo_Visibility) bool {
	if b.GetProp(JSDocInfo_VISIBILITY) == nil {
		b.populated = true
		b.SetProp(JSDocInfo_VISIBILITY, visibility)
		return true
	}
	return false
}

// Java: JSDocInfo.Builder.overwriteVisibility(Visibility)
func (b *JSDocInfo_Builder) OverwriteVisibility(visibility JSDocInfo_Visibility) {
	b.populated = true
	b.SetProp(JSDocInfo_VISIBILITY, visibility)
}

// RecordParameter records a typed parameter; returns true if the typed
// parameter was recorded and false if a parameter with the same name was
// already defined.
//
// Java: JSDocInfo.Builder.recordParameter(String, JSTypeExpression)
func (b *JSDocInfo_Builder) RecordParameter(parameterName string, type_ *JSTypeExpression) bool {
	return !b.HasAnySingletonTypeTags() &&
		b.populateTypeMapEntry(JSDocInfo_PARAMETERS, RhinoStringPool_AddOrGet(parameterName), type_)
}

// RecordParameterDescription records a parameter's description; returns
// true if the parameter's description was recorded and false if a parameter
// with the same name was already defined.
//
// Java: JSDocInfo.Builder.recordParameterDescription(String, String)
func (b *JSDocInfo_Builder) RecordParameterDescription(parameterName string, description string) bool {
	if !b.ShouldParseDocumentation() {
		return true
	}
	return JSDocInfo_Builder_PopulatePropEntry(
		b, JSDocInfo_PARAMETER_DESCRIPTIONS, RhinoStringPool_AddOrGet(parameterName), description)
}

// RecordTemplateTypeName records a template type name; returns true if the
// template type name was recorded and false if the input template type name
// was already defined.
//
// Java: JSDocInfo.Builder.recordTemplateTypeName(String)
func (b *JSDocInfo_Builder) RecordTemplateTypeName(name string) bool {
	return b.RecordTemplateTypeNameWithBound(name, nil)
}

// Java: JSDocInfo.Builder.recordTemplateTypeName(String, JSTypeExpression)
func (b *JSDocInfo_Builder) RecordTemplateTypeNameWithBound(name string, bound *JSTypeExpression) bool {
	if bound == nil {
		bound = JSTypeExpression_IMPLICIT_TEMPLATE_BOUND
	}
	transformations := jsdocGet[*jutil.LinkedHashMap[string, *Node]](b.GetProp(JSDocInfo_TYPE_TRANSFORMATIONS))
	if (transformations != nil && transformations.ContainsKey(name)) ||
		b.hasPropKey(JSDocInfo_TYPEDEF_TYPE) {
		return false
	}
	return b.populateTypeMapEntry(b_propTemplateTypeNames(), RhinoStringPool_AddOrGet(name), bound)
}

func b_propTemplateTypeNames() *JSDocInfo_Property { return JSDocInfo_TEMPLATE_TYPE_NAMES }

func (b *JSDocInfo_Builder) populateTypeMapEntry(prop *JSDocInfo_Property, key string, value *JSTypeExpression) bool {
	raw := b.GetProp(prop)
	var m *jutil.LinkedHashMap[string, *JSTypeExpression]
	if raw == nil {
		m = jutil.NewLinkedHashMap[string, *JSTypeExpression]()
		b.SetProp(prop, m)
	} else {
		m = jsdocGet[*jutil.LinkedHashMap[string, *JSTypeExpression]](raw)
	}
	if m.PutIfAbsent(key, value) == nil {
		b.populated = true
		return true
	}
	return false
}

// RecordTypeTransformation records a type transformation expression
// together with its template type name.
//
// Java: JSDocInfo.Builder.recordTypeTransformation(String, Node)
func (b *JSDocInfo_Builder) RecordTypeTransformation(name string, expr *Node) bool {
	names := jsdocGet[*jutil.LinkedHashMap[string, *JSTypeExpression]](b.GetProp(JSDocInfo_TEMPLATE_TYPE_NAMES))
	if names != nil && names.ContainsKey(name) {
		return false
	}
	return JSDocInfo_Builder_PopulatePropEntry(b, JSDocInfo_TYPE_TRANSFORMATIONS, RhinoStringPool_AddOrGet(name), expr)
}

func jsdocNewStringList() any { return jutil.NewArrayList[string]() }

// RecordThrowsAnnotation records a throw annotation description; returns
// true if the type's description was recorded. The description of a throw
// annotation is the text including the type.
//
// Java: JSDocInfo.Builder.recordThrowsAnnotation(String)
func (b *JSDocInfo_Builder) RecordThrowsAnnotation(annotation string) bool {
	b.populated = true
	if b.CheckBit(JSDocInfo_Bit_NOSIDEEFFECTS) {
		return false
	}
	// TODO(user): Does it make sense to check for singleton tags here?
	// Note that if the @throws annotation appears before a singleton tag like
	// @type, the throws annotation is preserved, but if it appears after the
	// singleton tag, it gets dropped.
	if !b.HasAnySingletonTypeTags() {
		throwsAnnotations := jsdocGet[*jutil.ArrayList[string]](b.GetPropWithDefault(JSDocInfo_THROWS_ANNOTATIONS, jsdocNewStringList))
		if b.ShouldParseDocumentation() {
			throwsAnnotations.Add(annotation)
		} else if throwsAnnotations.IsEmpty() {
			// Add at least one annotation so that PureFunctionIdentifier knows
			// about the side effect.
			throwsAnnotations.Add("")
		}
	}
	return true
}

// RecordAuthor adds an author to the current information.
//
// Java: JSDocInfo.Builder.recordAuthor(String)
func (b *JSDocInfo_Builder) RecordAuthor(author string) bool {
	b.populated = true
	if b.ShouldParseDocumentation() {
		jsdocGet[*jutil.ArrayList[string]](b.GetPropWithDefault(JSDocInfo_AUTHORS, jsdocNewStringList)).Add(author)
	}
	// NOTE: this could be removed, since it's always true.
	return true
}

// RecordReference adds a reference ("@see") to the current information.
//
// Java: JSDocInfo.Builder.recordReference(String)
func (b *JSDocInfo_Builder) RecordReference(reference string) bool {
	b.populated = true
	if b.ShouldParseDocumentation() {
		jsdocGet[*jutil.ArrayList[string]](b.GetPropWithDefault(JSDocInfo_SEES, jsdocNewStringList)).Add(reference)
	}
	// NOTE: this could be removed, since it's always true.
	return true
}

// RecordConsistentIdGenerator records that the JSDocInfo being built should
// have its JSDocInfo#isConsistentIdGenerator() flag set to true; returns
// true if the consistentIdGenerator flag was recorded and false if it was
// already recorded.
//
// Java: JSDocInfo.Builder.recordConsistentIdGenerator()
func (b *JSDocInfo_Builder) RecordConsistentIdGenerator() bool {
	return b.PopulateProp(JSDocInfo_ID_GENERATOR, JSDocInfo_IdGenerator_CONSISTENT)
}

// RecordStableIdGenerator records that the JSDocInfo being built should have
// its JSDocInfo#isStableIdGenerator() flag set to true; returns true if the
// stableIdGenerator flag was recorded and false if it was already recorded.
//
// Java: JSDocInfo.Builder.recordStableIdGenerator()
func (b *JSDocInfo_Builder) RecordStableIdGenerator() bool {
	return b.PopulateProp(JSDocInfo_ID_GENERATOR, JSDocInfo_IdGenerator_STABLE)
}

// RecordXidGenerator records that the JSDocInfo being built should have its
// JSDocInfo#isXidGenerator() flag set to true; returns true if the
// isXidGenerator flag was recorded and false if it was already recorded.
//
// Java: JSDocInfo.Builder.recordXidGenerator()
func (b *JSDocInfo_Builder) RecordXidGenerator() bool {
	return b.PopulateProp(JSDocInfo_ID_GENERATOR, JSDocInfo_IdGenerator_XID)
}

// RecordMappedIdGenerator records that the JSDocInfo being built should have
// its JSDocInfo#isStableIdGenerator() flag set to true; returns true if the
// stableIdGenerator flag was recorded and false if it was already recorded.
//
// Java: JSDocInfo.Builder.recordMappedIdGenerator()
func (b *JSDocInfo_Builder) RecordMappedIdGenerator() bool {
	return b.PopulateProp(JSDocInfo_ID_GENERATOR, JSDocInfo_IdGenerator_MAPPED)
}

// RecordIdGenerator records that the JSDocInfo being built should have its
// JSDocInfo#isIdGenerator() flag set to true; returns true if the
// idGenerator flag was recorded and false if it was already recorded.
//
// Java: JSDocInfo.Builder.recordIdGenerator()
func (b *JSDocInfo_Builder) RecordIdGenerator() bool {
	return b.PopulateProp(JSDocInfo_ID_GENERATOR, JSDocInfo_IdGenerator_UNIQUE)
}

// RecordDeprecationReason records the deprecation reason.
//
// Java: JSDocInfo.Builder.recordDeprecationReason(String)
func (b *JSDocInfo_Builder) RecordDeprecationReason(reason string) bool {
	return b.PopulateProp(JSDocInfo_DEPRECATION_REASON, reason)
}

// IsDeprecationReasonRecorded returns whether a deprecation reason has been
// recorded.
//
// Java: JSDocInfo.Builder.isDeprecationReasonRecorded()
func (b *JSDocInfo_Builder) IsDeprecationReasonRecorded() bool {
	return b.GetProp(JSDocInfo_DEPRECATION_REASON) != nil
}

// Java: JSDocInfo.Builder.recordSuppressions(ImmutableSet<String>, String)
const suppressionInternCacheSize = 2048

type internedSuppressionEntry struct {
	names []string
	set   *jutil.LinkedHashSet[string]
	mp    *jutil.ObjectLinkedHashMap[*jutil.LinkedHashSet[string], string]
}

var suppressionInternCache [suppressionInternCacheSize]atomic.Pointer[internedSuppressionEntry]

func internSuppressionSlice(rawNames []string) *internedSuppressionEntry {
	var uniqBuf [8]string
	var uniq []string
	if len(rawNames) <= 1 {
		uniq = rawNames
	} else {
		if len(rawNames) <= len(uniqBuf) {
			uniq = uniqBuf[:0]
		} else {
			uniq = make([]string, 0, len(rawNames))
		}
		for _, s := range rawNames {
			dup := false
			for _, u := range uniq {
				if u == s {
					dup = true
					break
				}
			}
			if !dup {
				uniq = append(uniq, s)
			}
		}
	}

	h := uint32(2166136261)
	for _, s := range uniq {
		h = (h ^ uint32(jutil.StringHashCode(s))) * 16777619
	}
	idx := h & (suppressionInternCacheSize - 1)
	if cached := suppressionInternCache[idx].Load(); cached != nil && len(cached.names) == len(uniq) {
		match := true
		for i, u := range uniq {
			if cached.names[i] != u {
				match = false
				break
			}
		}
		if match {
			return cached
		}
	}

	namesCopy := make([]string, len(uniq))
	set := jutil.NewLinkedHashSet[string]()
	for i, u := range uniq {
		interned := RhinoStringPool_AddOrGet(u)
		namesCopy[i] = interned
		set.Add(interned)
	}
	mp := jutil.NewObjectLinkedHashMap[*jutil.LinkedHashSet[string], string]()
	mp.Put(set, "")
	entry := &internedSuppressionEntry{names: namesCopy, set: set, mp: mp}
	suppressionInternCache[idx].Store(entry)
	return entry
}

// RecordSuppressionSliceWithDescription records a parsed @suppress slice without
// allocating temporary LinkedHashSets or ObjectLinkedHashMaps on the common path.
func (b *JSDocInfo_Builder) RecordSuppressionSliceWithDescription(rawNames []string, description string) {
	b.populated = true
	entry := internSuppressionSlice(rawNames)
	current := jsdocGet[*jutil.ObjectLinkedHashMap[*jutil.LinkedHashSet[string], string]](b.GetProp(JSDocInfo_SUPPRESSIONS))
	if current == nil && description == "" {
		b.SetProp(JSDocInfo_SUPPRESSIONS, entry.mp)
		return
	}
	if current != nil && current.ContainsKey(entry.set) {
		return
	}
	mapBuilder := jutil.NewObjectLinkedHashMap[*jutil.LinkedHashSet[string], string]()
	if current != nil {
		for k, v := range current.All() {
			mapBuilder.Put(k, v)
		}
	}
	mapBuilder.Put(entry.set, description)
	b.SetProp(JSDocInfo_SUPPRESSIONS, mapBuilder)
}

func (b *JSDocInfo_Builder) RecordSuppressionsWithDescription(suppressions *jutil.LinkedHashSet[string], description string) {
	if suppressions != nil && suppressions.Size() <= 8 {
		var buf [8]string
		names := buf[:0]
		for s := range suppressions.All() {
			names = append(names, s)
		}
		b.RecordSuppressionSliceWithDescription(names, description)
		return
	}
	b.populated = true
	mapBuilder := jutil.NewObjectLinkedHashMap[*jutil.LinkedHashSet[string], string]()
	current := jsdocGet[*jutil.ObjectLinkedHashMap[*jutil.LinkedHashSet[string], string]](b.GetProp(JSDocInfo_SUPPRESSIONS))
	if current != nil {
		if current.ContainsKey(suppressions) {
			// Exact @suppress warning exists already. Return without recording.
			return
		}
		for k, v := range current.All() {
			mapBuilder.Put(k, v)
		}
	}
	key := JSDocInfo_InternSet(suppressions)
	if mapBuilder.ContainsKey(key) {
		// Java: ImmutableMap.Builder.buildOrThrow() rejects duplicate keys.
		panic(jutil.NewIllegalArgumentExceptionMsg("Multiple entries with same key: " + key.String()))
	}
	mapBuilder.Put(key, description)
	suppressionsMap := mapBuilder
	b.SetProp(JSDocInfo_SUPPRESSIONS, suppressionsMap)
}

// RecordSuppressions records the list of suppressed warnings, possibly
// adding to the set of already configured warnings.
//
// Java: JSDocInfo.Builder.recordSuppressions(Set<String>)
func (b *JSDocInfo_Builder) RecordSuppressions(suppressions *jutil.LinkedHashSet[string]) {
	jutil.CheckNotNull(suppressions)
	b.RecordSuppressionsWithDescription(suppressions, "")
}

// Java: JSDocInfo.Builder.recordSuppression(String)
func (b *JSDocInfo_Builder) RecordSuppression(suppression string) {
	var single [1]string
	single[0] = suppression
	b.RecordSuppressionSliceWithDescription(single[:], "")
}

// RecordModifies records the list of modifies warnings.
//
// Java: JSDocInfo.Builder.recordModifies(Set<String>)
func (b *JSDocInfo_Builder) RecordModifies(modifies *jutil.LinkedHashSet[string]) bool {
	return !b.HasAnySingletonSideEffectTags() && b.PopulateProp(JSDocInfo_MODIFIES, JSDocInfo_InternSet(modifies))
}

// RecordType records a type; returns true if the type was recorded and false
// if it is invalid or was already defined.
//
// Java: JSDocInfo.Builder.recordType(JSTypeExpression)
func (b *JSDocInfo_Builder) RecordType(type_ *JSTypeExpression) bool {
	return type_ != nil && !b.HasAnyTypeRelatedTags() && b.PopulateProp(JSDocInfo_TYPE, type_)
}

// Java: JSDocInfo.Builder.recordInlineType()
func (b *JSDocInfo_Builder) RecordInlineType() {
	b.PopulateBit(JSDocInfo_Bit_INLINE_TYPE, true)
}

// RecordTypedef records that the JSDocInfo being built should be populated
// with a `typedef`'d type.
//
// Java: JSDocInfo.Builder.recordTypedef(JSTypeExpression)
func (b *JSDocInfo_Builder) RecordTypedef(type_ *JSTypeExpression) bool {
	return type_ != nil &&
		!b.HasAnyTypeRelatedTags() &&
		b.GetProp(JSDocInfo_TEMPLATE_TYPE_NAMES) == nil &&
		b.PopulateProp(JSDocInfo_TYPEDEF_TYPE, type_)
}

// RecordReturnType records a return type; returns true if the return type
// was recorded and false if it is invalid or was already defined.
//
// Java: JSDocInfo.Builder.recordReturnType(JSTypeExpression)
func (b *JSDocInfo_Builder) RecordReturnType(type_ *JSTypeExpression) bool {
	return type_ != nil && !b.HasAnySingletonTypeTags() && b.PopulateProp(JSDocInfo_RETURN_TYPE, type_)
}

// RecordReturnDescription records a return description; returns true if the
// return description was recorded and false if it is invalid or was already
// defined.
//
// Java: JSDocInfo.Builder.recordReturnDescription(String)
func (b *JSDocInfo_Builder) RecordReturnDescription(description string) bool {
	if !b.ShouldParseDocumentation() {
		return true
	}
	return b.PopulateProp(JSDocInfo_RETURN_DESCRIPTION, description)
}

// RecordDefineType records the type of a define.
//
// 'Define' values are special constants that may be manipulated by the
// compiler. They are designed to mimic the #define command in the C
// preprocessor.
//
// Java: JSDocInfo.Builder.recordDefineType(JSTypeExpression)
func (b *JSDocInfo_Builder) RecordDefineType(type_ *JSTypeExpression) bool {
	if type_ != nil && !b.CheckBit(JSDocInfo_Bit_CONST) && !b.CheckBit(JSDocInfo_Bit_DEFINE) && b.RecordType(type_) {
		return b.PopulateBit(JSDocInfo_Bit_DEFINE, true)
	}
	return false
}

// RecordEnumParameterType records a parameter type to an enum; returns true
// if the enum's parameter type was recorded and false if it was invalid or
// already defined.
//
// Java: JSDocInfo.Builder.recordEnumParameterType(JSTypeExpression)
func (b *JSDocInfo_Builder) RecordEnumParameterType(type_ *JSTypeExpression) bool {
	if type_ != nil && !b.HasAnyTypeRelatedTags() {
		b.SetProp(JSDocInfo_ENUM_PARAMETER_TYPE, type_)
		b.populated = true
		return true
	}
	return false
}

// TODO(tbreisacher): Disallow nullable types here. If someone writes "@this
// {Foo}" in their JS we automatically treat it as though they'd written
// "@this {!Foo}". But, if the type node is created in the compiler (e.g. in
// the WizPass) we should explicitly add the '!'

// RecordThisType records a type for `@this` annotation; returns true if the
// type was recorded and false if it is invalid or if it collided with
// `@enum` or `@type` annotations.
//
// Java: JSDocInfo.Builder.recordThisType(JSTypeExpression)
func (b *JSDocInfo_Builder) RecordThisType(type_ *JSTypeExpression) bool {
	return type_ != nil && !b.HasAnySingletonTypeTags() && b.PopulateProp(JSDocInfo_THIS_TYPE, type_)
}

// RecordBaseType records a base type; returns true if the base type was
// recorded and false if it was already defined.
//
// Java: JSDocInfo.Builder.recordBaseType(JSTypeExpression)
func (b *JSDocInfo_Builder) RecordBaseType(type_ *JSTypeExpression) bool {
	return type_ != nil && !b.HasAnySingletonTypeTags() && b.PopulateProp(JSDocInfo_BASE_TYPE, type_)
}

// ChangeBaseType changes a base type, even if one has already been set on
// currentInfo; returns true if the base type was changed successfully.
//
// Java: JSDocInfo.Builder.changeBaseType(JSTypeExpression)
func (b *JSDocInfo_Builder) ChangeBaseType(type_ *JSTypeExpression) bool {
	if type_ != nil && !b.HasAnySingletonTypeTags() {
		b.SetProp(JSDocInfo_BASE_TYPE, type_)
		b.populated = true
		return true
	}
	return false
}

// RecordConstancy records that the JSDocInfo being built should have its
// JSDocInfo#isConstant() flag set to true; returns true if the constancy was
// recorded and false if it was already defined.
//
// Java: JSDocInfo.Builder.recordConstancy()
func (b *JSDocInfo_Builder) RecordConstancy() bool {
	return b.PopulateBit(JSDocInfo_Bit_CONST, true)
}

// RecordMutable records that the JSDocInfo being built should have its
// JSDocInfo#isConstant() flag set to false; returns true if the mutability
// was recorded and false if it was already defined.
//
// Java: JSDocInfo.Builder.recordMutable()
func (b *JSDocInfo_Builder) RecordMutable() bool {
	return b.PopulateBit(JSDocInfo_Bit_CONST, false)
}

// RecordFinality records that the JSDocInfo being built should have its
// JSDocInfo#isFinal() flag set to true; returns true if the finality was
// recorded and false if it was already defined.
//
// Java: JSDocInfo.Builder.recordFinality()
func (b *JSDocInfo_Builder) RecordFinality() bool {
	return b.PopulateBit(JSDocInfo_Bit_FINAL, true)
}

// RecordDescription records a description giving context for translation
// (i18n); returns true if the description was recorded and false if the
// description was invalid or was already defined.
//
// Java: JSDocInfo.Builder.recordDescription(String)
func (b *JSDocInfo_Builder) RecordDescription(description string) bool {
	return b.PopulateProp(JSDocInfo_DESCRIPTION, description)
}

// RecordTsType records a tsType giving context for .d.ts generation.
//
// Java: JSDocInfo.Builder.recordTsType(String)
func (b *JSDocInfo_Builder) RecordTsType(tsType string) {
	b.populated = true
	jsdocGet[*jutil.ArrayList[string]](b.GetPropWithDefault(JSDocInfo_TS_TYPES, jsdocNewStringList)).Add(tsType)
}

// RecordMeaning records a meaning giving context for translation (i18n).
// Different meanings will result in different translations; returns true if
// the meaning was successfully updated.
//
// Java: JSDocInfo.Builder.recordMeaning(String)
func (b *JSDocInfo_Builder) RecordMeaning(meaning string) bool {
	return b.PopulateProp(JSDocInfo_MEANING, meaning)
}

// RecordAlternateMessageId records an ID for an alternate message to be used
// if this message is not yet translated; returns true if the alternate
// message ID was successfully updated.
//
// Java: JSDocInfo.Builder.recordAlternateMessageId(String)
func (b *JSDocInfo_Builder) RecordAlternateMessageId(alternateMessageId string) bool {
	return b.PopulateProp(JSDocInfo_ALTERNATE_MESSAGE_ID, alternateMessageId)
}

// RecordClosurePrimitiveId records an identifier for a Closure Primitive.
// function; returns true if the id was successfully updated.
//
// Java: JSDocInfo.Builder.recordClosurePrimitiveId(String)
func (b *JSDocInfo_Builder) RecordClosurePrimitiveId(closurePrimitiveId string) bool {
	return b.PopulateProp(JSDocInfo_CLOSURE_PRIMITIVE_ID, RhinoStringPool_AddOrGet(closurePrimitiveId))
}

// RecordFileOverview records a fileoverview description; returns true if the
// description was recorded and false if the description was invalid or was
// already defined.
//
// Java: JSDocInfo.Builder.recordFileOverview(String)
func (b *JSDocInfo_Builder) RecordFileOverview(description string) bool {
	b.SetBit(JSDocInfo_Bit_FILEOVERVIEW, true)
	b.populated = true
	if !b.ShouldParseDocumentation() {
		return true
	}
	return b.PopulateProp(JSDocInfo_FILEOVERVIEW_DESCRIPTION, description)
}

// RecordEnhance records enhanced namespace; returns true if the enhanced
// namespace was recorded.
//
// Java: JSDocInfo.Builder.recordEnhance(String)
func (b *JSDocInfo_Builder) RecordEnhance(namespace string) bool {
	b.SetBit(JSDocInfo_Bit_ENHANCED_NAMESPACE, true)
	b.populated = true
	return b.PopulateProp(JSDocInfo_ENHANCED_NAMESPACE, namespace)
}

// IsModsRecorded returns whether this builder recorded a modded namespace.
//
// Java: JSDocInfo.Builder.isModsRecorded()
func (b *JSDocInfo_Builder) IsModsRecorded() bool {
	return b.GetProp(JSDocInfo_MODS) != nil
}

// RecordMods records modded namespace; returns true if the modded namespace
// was recorded.
//
// Java: JSDocInfo.Builder.recordMods(String)
func (b *JSDocInfo_Builder) RecordMods(namespace string) bool {
	return b.PopulateProp(JSDocInfo_MODS, RhinoStringPool_AddOrGet(namespace))
}

// Java: JSDocInfo.Builder.recordLicense(String)
func (b *JSDocInfo_Builder) RecordLicense(license string) bool {
	b.SetProp(JSDocInfo_LICENSE, RhinoStringPool_AddOrGet(license))
	b.populated = true
	return true
}

// Java: JSDocInfo.Builder.addLicense(String)
func (b *JSDocInfo_Builder) AddLicense(license string) bool {
	// The vast majority of JSDoc doesn't have @license so it make sense to be
	// lazy about building the LinkedHashSet.
	if b.licenseTexts == nil {
		// The LinkedHashSet is only used to remove duplicates, it is never read
		// beyond the add, so LinkedHashSet is not required.
		b.licenseTexts = jutil.NewLinkedHashSet[string]()
	}

	if !b.licenseTexts.Add(license) {
		return false
	}

	txt := jsdocGet[string](b.GetProp(JSDocInfo_LICENSE)) // Strings.nullToEmpty(txt)
	return b.RecordLicense(RhinoStringPool_AddOrGet(txt + license))
}

// RecordNoCompile records that the JSDocInfo being built should have its
// JSDocInfo#isNoCompile() flag set to `true`. Returns `true` if the no
// compile flag was recorded and `false` if it was already recorded
//
// Java: JSDocInfo.Builder.recordNoCompile()
func (b *JSDocInfo_Builder) RecordNoCompile() bool {
	return b.PopulateBit(JSDocInfo_Bit_NOCOMPILE, true)
}

// RecordNoDts records that the JSDocInfo being built should have its
// JSDocInfo#isNoDtsOutput() flag set to `true`. Returns `true` if the no
// compile flag was recorded and `false` if it was already recorded
//
// Java: JSDocInfo.Builder.recordNoDts()
func (b *JSDocInfo_Builder) RecordNoDts() bool {
	return b.PopulateBit(JSDocInfo_Bit_NODTS, true)
}

// RecordNoCollapse records that the JSDocInfo being built should have its
// JSDocInfo#isNoCollapse() flag set to `true`. Returns `true` if the no
// collapse flag was recorded and `false` if it was already recorded
//
// Java: JSDocInfo.Builder.recordNoCollapse()
func (b *JSDocInfo_Builder) RecordNoCollapse() bool {
	return b.PopulateBit(JSDocInfo_Bit_NOCOLLAPSE, true)
}

// RecordNoInline records that the JSDocInfo being built should have its
// JSDocInfo#isNoInline() flag set to `true`. Returns `true` if the no inline
// flag was recorded and `false` if it was already recorded
//
// Java: JSDocInfo.Builder.recordNoInline()
func (b *JSDocInfo_Builder) RecordNoInline() bool {
	return b.PopulateBit(JSDocInfo_Bit_NOINLINE, true)
}

// RecordEncourageInlining records that the JSDocInfo being built should have
// its JSDocInfo#isEncourageInlining() flag set to `true`. Returns `true` if
// the encourageInlining flag was recorded and `false` if it was already
// recorded
//
// Java: JSDocInfo.Builder.recordEncourageInlining()
func (b *JSDocInfo_Builder) RecordEncourageInlining() bool {
	return b.PopulateBit(JSDocInfo_Bit_ENCOURAGE_INLINING, true)
}

// RecordRequireInlining records that the JSDocInfo being built should have
// its JSDocInfo#isRequireInlining() flag set to `true`. Returns `true` if the
// requireInlining flag was recorded and `false` if it was already recorded
//
// Java: JSDocInfo.Builder.recordRequireInlining()
func (b *JSDocInfo_Builder) RecordRequireInlining() bool {
	return b.PopulateBit(JSDocInfo_Bit_REQUIRE_INLINING, true)
}

// RecordPureOrBreakMyCode records that the JSDocInfo being built should have
// its JSDocInfo#isPureOrBreakMyCode() flag set to `true`. Returns `true` if
// the no pureOrBreakMyCode flag was recorded and `false` if it was already
// recorded
//
// Java: JSDocInfo.Builder.recordPureOrBreakMyCode()
func (b *JSDocInfo_Builder) RecordPureOrBreakMyCode() bool {
	return b.PopulateBit(JSDocInfo_Bit_PURE_OR_BREAK_MY_CODE, true)
}

// RecordCollapsibleOrBreakMyCode records that the JSDocInfo being built
// should have its JSDocInfo#isCollapsibleOrBreakMyCode() flag set to `true`.
// Returns `true` if the no collapsibleOrBreakMyCode flag was recorded and
// `false` if it was already recorded
//
// Java: JSDocInfo.Builder.recordCollapsibleOrBreakMyCode()
func (b *JSDocInfo_Builder) RecordCollapsibleOrBreakMyCode() bool {
	return b.PopulateBit(JSDocInfo_Bit_COLLAPSIBLE_OR_BREAK_MY_CODE, true)
}

// RecordConstructor records that the JSDocInfo being built should have its
// JSDocInfo#isConstructor() flag set to `true`. Returns `true` if the
// constructor was recorded and `false` if it was already defined or it was
// incompatible with the existing flags
//
// Java: JSDocInfo.Builder.recordConstructor()
func (b *JSDocInfo_Builder) RecordConstructor() bool {
	return !b.HasAnySingletonTypeTags() &&
		!b.IsConstructorOrInterface() &&
		b.PopulateBit(JSDocInfo_Bit_CONSTRUCTOR, true)
}

// RecordImplicitMatch records that the JSDocInfo being built should have its
// JSDocInfo#usesImplicitMatch() flag set to `true`. Returns `true` if the
// `@record` tag was recorded and `false` if it was already defined or it was
// incompatible with the existing flags
//
// Java: JSDocInfo.Builder.recordImplicitMatch()
func (b *JSDocInfo_Builder) RecordImplicitMatch() bool {
	return !b.HasAnySingletonTypeTags() &&
		!b.IsConstructorOrInterface() &&
		b.PopulateBit(JSDocInfo_Bit_RECORD, true) &&
		b.PopulateBit(JSDocInfo_Bit_INTERFACE, true)
}

// Java: JSDocInfo.Builder.recordProvideGoog()
func (b *JSDocInfo_Builder) RecordProvideGoog() bool {
	return b.PopulateBit(JSDocInfo_Bit_PROVIDE_GOOG, true)
}

// Java: JSDocInfo.Builder.recordProvideAlreadyProvided()
func (b *JSDocInfo_Builder) RecordProvideAlreadyProvided() bool {
	return b.PopulateBit(JSDocInfo_Bit_PROVIDE_ALREADY_PROVIDED, true)
}

// IsConstructorRecorded reports whether the JSDocInfo being built will have
// its JSDocInfo#isConstructor() flag set to true.
//
// Java: JSDocInfo.Builder.isConstructorRecorded()
func (b *JSDocInfo_Builder) IsConstructorRecorded() bool {
	return b.CheckBit(JSDocInfo_Bit_CONSTRUCTOR)
}

// RecordUnrestricted records that the JSDocInfo being built should have its
// JSDocInfo#makesUnrestricted() flag set to true; returns true if annotation
// was recorded and false if it was already defined or it was incompatible
// with the existing flags.
//
// Java: JSDocInfo.Builder.recordUnrestricted()
func (b *JSDocInfo_Builder) RecordUnrestricted() bool {
	return !b.HasAnySingletonTypeTags() &&
		((b.bits & (JSDocInfo_Bit_INTERFACE.Mask() | JSDocInfo_Bit_DICT.Mask() | JSDocInfo_Bit_STRUCT.Mask())) == 0) &&
		b.PopulateBit(JSDocInfo_Bit_UNRESTRICTED, true)
}

// Java: JSDocInfo.Builder.isUnrestrictedRecorded()
func (b *JSDocInfo_Builder) IsUnrestrictedRecorded() bool {
	return b.CheckBit(JSDocInfo_Bit_UNRESTRICTED)
}

// RecordAbstract records that the JSDocInfo being built should have its
// JSDocInfo#isAbstract() flag set to true; returns true if the flag was
// recorded and false if it was already defined or it was incompatible with
// the existing flags.
//
// Java: JSDocInfo.Builder.recordAbstract()
func (b *JSDocInfo_Builder) RecordAbstract() bool {
	visibility := b.GetProp(JSDocInfo_VISIBILITY)
	isPrivate := visibility != nil && jsdocGet[JSDocInfo_Visibility](visibility) == JSDocInfo_Visibility_PRIVATE
	return !b.HasAnySingletonTypeTags() &&
		((b.bits & (JSDocInfo_Bit_INTERFACE.Mask() | JSDocInfo_Bit_FINAL.Mask())) == 0) &&
		!isPrivate &&
		b.PopulateBit(JSDocInfo_Bit_ABSTRACT, true)
}

// RecordStruct records that the JSDocInfo being built should have its
// JSDocInfo#makesStructs() flag set to true; returns true if the struct was
// recorded and false if it was already defined or it was incompatible with
// the existing flags.
//
// Java: JSDocInfo.Builder.recordStruct()
func (b *JSDocInfo_Builder) RecordStruct() bool {
	return !b.HasAnySingletonTypeTags() &&
		((b.bits & (JSDocInfo_Bit_DICT.Mask() | JSDocInfo_Bit_UNRESTRICTED.Mask())) == 0) &&
		b.PopulateBit(JSDocInfo_Bit_STRUCT, true)
}

// Java: JSDocInfo.Builder.isStructRecorded()
func (b *JSDocInfo_Builder) IsStructRecorded() bool {
	return b.CheckBit(JSDocInfo_Bit_STRUCT)
}

// RecordDict records that the JSDocInfo being built should have its
// JSDocInfo#makesDicts() flag set to true; returns true if the dict was
// recorded and false if it was already defined or it was incompatible with
// the existing flags.
//
// Java: JSDocInfo.Builder.recordDict()
func (b *JSDocInfo_Builder) RecordDict() bool {
	return !b.HasAnySingletonTypeTags() &&
		((b.bits & (JSDocInfo_Bit_STRUCT.Mask() | JSDocInfo_Bit_UNRESTRICTED.Mask())) == 0) &&
		b.PopulateBit(JSDocInfo_Bit_DICT, true)
}

// Java: JSDocInfo.Builder.isDictRecorded()
func (b *JSDocInfo_Builder) IsDictRecorded() bool {
	return b.CheckBit(JSDocInfo_Bit_DICT)
}

// RecordOverride records that the JSDocInfo being built should have its
// JSDocInfo#isOverride() flag set to true.
//
// Java: JSDocInfo.Builder.recordOverride()
func (b *JSDocInfo_Builder) RecordOverride() bool {
	return b.PopulateBit(JSDocInfo_Bit_OVERRIDE, true)
}

// RecordDeprecated records that the JSDocInfo being built should have its
// JSDocInfo#isDeprecated() flag set to true.
//
// Java: JSDocInfo.Builder.recordDeprecated()
func (b *JSDocInfo_Builder) RecordDeprecated() bool {
	return b.PopulateBit(JSDocInfo_Bit_DEPRECATED, true)
}

// RecordInterface records that the JSDocInfo being built should have its
// JSDocInfo#isInterface() flag set to true; returns true if the flag was
// recorded and false if it was already defined or it was incompatible with
// the existing flags.
//
// Java: JSDocInfo.Builder.recordInterface()
func (b *JSDocInfo_Builder) RecordInterface() bool {
	return !b.HasAnySingletonTypeTags() &&
		((b.bits & (JSDocInfo_Bit_CONSTRUCTOR.Mask() | JSDocInfo_Bit_ABSTRACT.Mask())) == 0) &&
		b.PopulateBit(JSDocInfo_Bit_INTERFACE, true)
}

// RecordExport records that the JSDocInfo being built should have its
// JSDocInfo#isExport() flag set to true.
//
// Java: JSDocInfo.Builder.recordExport()
func (b *JSDocInfo_Builder) RecordExport() bool {
	return b.PopulateBit(JSDocInfo_Bit_EXPORT, true)
}

// RemoveExport records that the JSDocInfo being built should have its
// JSDocInfo#isExport() flag set to false.
//
// Java: JSDocInfo.Builder.removeExport()
func (b *JSDocInfo_Builder) RemoveExport() bool {
	return b.PopulateBit(JSDocInfo_Bit_EXPORT, false)
}

// RecordImplicitCast records that the JSDocInfo being built should have its
// JSDocInfo#isImplicitCast() flag set to true.
//
// Java: JSDocInfo.Builder.recordImplicitCast()
func (b *JSDocInfo_Builder) RecordImplicitCast() bool {
	return b.PopulateBit(JSDocInfo_Bit_IMPLICITCAST, true)
}

// RecordNoSideEffects records that the JSDocInfo being built should have its
// JSDocInfo#isNoSideEffects() flag set to true.
//
// Java: JSDocInfo.Builder.recordNoSideEffects()
func (b *JSDocInfo_Builder) RecordNoSideEffects() bool {
	return !b.HasAnySingletonSideEffectTags() &&
		b.IsPropEmpty(JSDocInfo_THROWS_ANNOTATIONS) &&
		b.PopulateBit(JSDocInfo_Bit_NOSIDEEFFECTS, true)
}

// IsNoSideEffectsRecorded returns whether current JSDoc is annotated with
// @nosideeffects.
//
// Java: JSDocInfo.Builder.isNoSideEffectsRecorded()
func (b *JSDocInfo_Builder) IsNoSideEffectsRecorded() bool {
	return b.CheckBit(JSDocInfo_Bit_NOSIDEEFFECTS)
}

// IsModifiesRecorded returns whether current JSDoc is annotated with
// @modifies.
//
// Java: JSDocInfo.Builder.isModifiesRecorded()
func (b *JSDocInfo_Builder) IsModifiesRecorded() bool {
	return !b.IsPropEmpty(JSDocInfo_MODIFIES)
}

// IsThrowsRecorded returns whether current JSDoc is annotated with @throws.
//
// Java: JSDocInfo.Builder.isThrowsRecorded()
func (b *JSDocInfo_Builder) IsThrowsRecorded() bool {
	return !b.IsPropEmpty(JSDocInfo_THROWS_ANNOTATIONS)
}

// RecordExterns records that the JSDocInfo being built should have its
// JSDocInfo#isExterns() flag set to true.
//
// Java: JSDocInfo.Builder.recordExterns()
func (b *JSDocInfo_Builder) RecordExterns() bool {
	return b.PopulateBit(JSDocInfo_Bit_EXTERNS, true)
}

// RecordNoCoverage records that the JSDocInfo being built should have its
// JSDocInfo#isNoCoverage() flag set to true.
//
// Java: JSDocInfo.Builder.recordNoCoverage()
func (b *JSDocInfo_Builder) RecordNoCoverage() bool {
	return b.PopulateBit(JSDocInfo_Bit_NOCOVERAGE, true)
}

// RecordTypeSummary records that the JSDocInfo being built should have its
// JSDocInfo#isTypeSummary() flag set to true.
//
// Java: JSDocInfo.Builder.recordTypeSummary()
func (b *JSDocInfo_Builder) RecordTypeSummary() bool {
	return b.PopulateBit(JSDocInfo_Bit_TYPE_SUMMARY, true)
}

// IsInterfaceRecorded reports whether the JSDocInfo being built will have
// its JSDocInfo#isInterface() flag set to true.
//
// Java: JSDocInfo.Builder.isInterfaceRecorded()
func (b *JSDocInfo_Builder) IsInterfaceRecorded() bool {
	return b.CheckBit(JSDocInfo_Bit_INTERFACE)
}

// HasParameter returns whether a parameter of the given name has already
// been recorded.
//
// Java: JSDocInfo.Builder.hasParameter(String)
func (b *JSDocInfo_Builder) HasParameter(name string) bool {
	params := jsdocGet[*jutil.LinkedHashMap[string, *JSTypeExpression]](b.GetProp(JSDocInfo_PARAMETERS))
	return params != nil && params.ContainsKey(name)
}

// RecordImplementedInterface records an implemented interface.
//
// Java: JSDocInfo.Builder.recordImplementedInterface(JSTypeExpression)
func (b *JSDocInfo_Builder) RecordImplementedInterface(interfaceType *JSTypeExpression) bool {
	return interfaceType != nil && b.AddUnique(JSDocInfo_IMPLEMENTED_INTERFACES, interfaceType)
}

// RecordExtendedInterface records an extended interface type.
//
// Java: JSDocInfo.Builder.recordExtendedInterface(JSTypeExpression)
func (b *JSDocInfo_Builder) RecordExtendedInterface(interfaceType *JSTypeExpression) bool {
	return interfaceType != nil && b.AddUnique(JSDocInfo_EXTENDED_INTERFACES, interfaceType)
}

// Java: JSDocInfo.Builder.addUnique(Property<ArrayList<JSTypeExpression>>, JSTypeExpression)
func (b *JSDocInfo_Builder) AddUnique(prop *JSDocInfo_Property, elem *JSTypeExpression) bool {
	list := jsdocGet[*jutil.ArrayList[*JSTypeExpression]](b.GetPropWithDefault(prop, func() any {
		return jutil.NewArrayList[*JSTypeExpression]()
	}))
	for e := range list.All() {
		if elem.IsEquivalentTo(e) {
			return false
		}
	}
	list.Add(elem)
	b.populated = true
	return true
}

// RecordLends records that we're lending to another name.
//
// Java: JSDocInfo.Builder.recordLends(JSTypeExpression)
func (b *JSDocInfo_Builder) RecordLends(name *JSTypeExpression) bool {
	return !b.HasAnyTypeRelatedTags() && b.PopulateProp(JSDocInfo_LENDS_NAME, name)
}

// IsNgInjectRecorded returns whether current JSDoc is annotated with
// `@ngInject`.
//
// Java: JSDocInfo.Builder.isNgInjectRecorded()
func (b *JSDocInfo_Builder) IsNgInjectRecorded() bool {
	return b.CheckBit(JSDocInfo_Bit_NG_INJECT)
}

// RecordNgInject records that we'd like to add `$inject` property inferred
// from parameters.
//
// Java: JSDocInfo.Builder.recordNgInject(boolean)
func (b *JSDocInfo_Builder) RecordNgInject(ngInject bool) bool {
	return b.PopulateBit(JSDocInfo_Bit_NG_INJECT, true)
}

// IsWizactionRecorded returns whether current JSDoc is annotated with
// `@wizaction`.
//
// Java: JSDocInfo.Builder.isWizactionRecorded()
func (b *JSDocInfo_Builder) IsWizactionRecorded() bool {
	return b.CheckBit(JSDocInfo_Bit_WIZ_ACTION)
}

// RecordWizaction records that this method is to be exposed as a wizaction.
//
// Java: JSDocInfo.Builder.recordWizaction()
func (b *JSDocInfo_Builder) RecordWizaction() bool {
	return b.PopulateBit(JSDocInfo_Bit_WIZ_ACTION, true)
}

// IsWizcallbackRecorded returns if current JSDoc is annotated with
// `@wizcallback`.
//
// Java: JSDocInfo.Builder.isWizcallbackRecorded()
func (b *JSDocInfo_Builder) IsWizcallbackRecorded() bool {
	return b.CheckBit(JSDocInfo_Bit_WIZ_CALLBACK)
}

// RecordWizcallback records that this method is to be exposed as a
// wizcallback.
//
// Java: JSDocInfo.Builder.recordWizcallback()
func (b *JSDocInfo_Builder) RecordWizcallback() bool {
	return b.PopulateBit(JSDocInfo_Bit_WIZ_CALLBACK, true)
}

// IsPolymerBehaviorRecorded returns whether current JSDoc is annotated with
// `@polymerBehavior`.
//
// Java: JSDocInfo.Builder.isPolymerBehaviorRecorded()
func (b *JSDocInfo_Builder) IsPolymerBehaviorRecorded() bool {
	return b.CheckBit(JSDocInfo_Bit_POLYMER_BEHAVIOR)
}

// RecordPolymerBehavior records that this method is to be exposed as a
// polymerBehavior.
//
// Java: JSDocInfo.Builder.recordPolymerBehavior()
func (b *JSDocInfo_Builder) RecordPolymerBehavior() bool {
	return b.PopulateBit(JSDocInfo_Bit_POLYMER_BEHAVIOR, true)
}

// IsPolymerRecorded returns whether current JSDoc is annotated with
// `@polymer`.
//
// Java: JSDocInfo.Builder.isPolymerRecorded()
func (b *JSDocInfo_Builder) IsPolymerRecorded() bool {
	return b.CheckBit(JSDocInfo_Bit_POLYMER)
}

// RecordPolymer records that this method is to be exposed as a polymer
// element.
//
// Java: JSDocInfo.Builder.recordPolymer()
func (b *JSDocInfo_Builder) RecordPolymer() bool {
	return b.PopulateBit(JSDocInfo_Bit_POLYMER, true)
}

// IsCustomElementRecorded returns whether current JSDoc is annotated with
// `@customElement`.
//
// Java: JSDocInfo.Builder.isCustomElementRecorded()
func (b *JSDocInfo_Builder) IsCustomElementRecorded() bool {
	return b.CheckBit(JSDocInfo_Bit_CUSTOM_ELEMENT)
}

// RecordCustomElement records that this method is to be exposed as a
// customElement.
//
// Java: JSDocInfo.Builder.recordCustomElement()
func (b *JSDocInfo_Builder) RecordCustomElement() bool {
	return b.PopulateBit(JSDocInfo_Bit_CUSTOM_ELEMENT, true)
}

// IsMixinClassRecorded returns whether current JSDoc is annotated with
// `@mixinClass`.
//
// Java: JSDocInfo.Builder.isMixinClassRecorded()
func (b *JSDocInfo_Builder) IsMixinClassRecorded() bool {
	return b.CheckBit(JSDocInfo_Bit_MIXIN_CLASS)
}

// RecordMixinClass records that this method is to be exposed as a mixinClass.
//
// Java: JSDocInfo.Builder.recordMixinClass()
func (b *JSDocInfo_Builder) RecordMixinClass() bool {
	return b.PopulateBit(JSDocInfo_Bit_MIXIN_CLASS, true)
}

// IsMixinFunctionRecorded returns whether current JSDoc is annotated with
// `@mixinFunction`.
//
// Java: JSDocInfo.Builder.isMixinFunctionRecorded()
func (b *JSDocInfo_Builder) IsMixinFunctionRecorded() bool {
	return b.CheckBit(JSDocInfo_Bit_MIXIN_FUNCTION)
}

// RecordMixinFunction records that this method is to be exposed as a
// mixinFunction.
//
// Java: JSDocInfo.Builder.recordMixinFunction()
func (b *JSDocInfo_Builder) RecordMixinFunction() bool {
	return b.PopulateBit(JSDocInfo_Bit_MIXIN_FUNCTION, true)
}

// IsSassGeneratedCssTsRecorded returns whether current JSDoc is annotated
// with `@sassGeneratedCssTs`.
//
// Java: JSDocInfo.Builder.isSassGeneratedCssTsRecorded()
func (b *JSDocInfo_Builder) IsSassGeneratedCssTsRecorded() bool {
	return b.CheckBit(JSDocInfo_Bit_SASS_GENERATED_CSS_TS)
}

// RecordSassGeneratedCssTs records that this script was generated by Sass
// from a.css.ts file.
//
// Java: JSDocInfo.Builder.recordSassGeneratedCssTs()
func (b *JSDocInfo_Builder) RecordSassGeneratedCssTs() {
	b.PopulateBit(JSDocInfo_Bit_SASS_GENERATED_CSS_TS, true)
}

// LogTypeInCompiler returns whether we should log the type of values passed
// to this function.
//
// Java: JSDocInfo.Builder.logTypeInCompiler()
func (b *JSDocInfo_Builder) LogTypeInCompiler() bool {
	return b.CheckBit(JSDocInfo_Bit_LOG_TYPE_IN_COMPILER)
}

// RecordLogTypeInCompiler records that the types of values passed to this
// method should be logged in the compiler.
//
// Java: JSDocInfo.Builder.recordLogTypeInCompiler()
func (b *JSDocInfo_Builder) RecordLogTypeInCompiler() bool {
	return b.PopulateBit(JSDocInfo_Bit_LOG_TYPE_IN_COMPILER, true)
}

// IsClosureUnawareCode returns whether JSDoc is annotated with the
// `@closureUnaware` annotation.
//
// Java: JSDocInfo.Builder.isClosureUnawareCode()
func (b *JSDocInfo_Builder) IsClosureUnawareCode() bool {
	return b.GetProp(JSDocInfo_CLOSURE_UNAWARE_CONFIG_VALUE) != nil
}

// RecordClosureUnawareCode records that this JSDoc was annotated with the
// `@closureUnaware` annotation.
//
// Java: JSDocInfo.Builder.recordClosureUnawareCode()
func (b *JSDocInfo_Builder) RecordClosureUnawareCode() bool {
	return b.PopulateProp(JSDocInfo_CLOSURE_UNAWARE_CONFIG_VALUE, JSDocInfo_PerFileClosureUnawareMode_UNSPECIFIED)
}

// RecordClosureUnawareCodeWithMode records that this JSDoc was annotated
// with the `@closureUnaware` annotation.
//
// Java: JSDocInfo.Builder.recordClosureUnawareCode(PerFileClosureUnawareMode)
func (b *JSDocInfo_Builder) RecordClosureUnawareCodeWithMode(mode JSDocInfo_PerFileClosureUnawareMode) bool {
	return b.PopulateProp(JSDocInfo_CLOSURE_UNAWARE_CONFIG_VALUE, mode)
}

// RemoveClosureUnawareCode removes the `@closureUnaware` annotation from
// this JSDoc, returning true the annotation was present before removal.
//
// Java: JSDocInfo.Builder.removeClosureUnawareCode()
func (b *JSDocInfo_Builder) RemoveClosureUnawareCode() bool {
	existing := b.removeProp(JSDocInfo_CLOSURE_UNAWARE_CONFIG_VALUE)
	return existing != nil
}

// RecordUsedViaDotConstructor records that this JSDoc was annotated with the
// `@usedViaDotConstructor` annotation.
//
// Java: JSDocInfo.Builder.recordUsedViaDotConstructor()
func (b *JSDocInfo_Builder) RecordUsedViaDotConstructor() bool {
	return b.PopulateBit(JSDocInfo_Bit_USED_VIA_DOT_CONSTRUCTOR, true)
}

// SetType force-sets the type, while recordType backs off. This is useful
// for (e.g.) copyFromWithNewType.
//
// TODO(sdh): this is a new method - consider removing it in favor of
// recordType? The main difference is that this force-sets the type, while
// recordType backs off.
//
// Java: JSDocInfo.Builder.setType(JSTypeExpression)
func (b *JSDocInfo_Builder) SetType(type_ *JSTypeExpression) *JSDocInfo_Builder {
	b.removeProp(JSDocInfo_RETURN_TYPE)
	b.removeProp(JSDocInfo_ENUM_PARAMETER_TYPE)
	b.removeProp(JSDocInfo_TYPEDEF_TYPE)
	b.SetProp(JSDocInfo_TYPE, type_)
	return b
}

// HasAnyTypeRelatedTags reports whether the current doc info has other type
// tags, like `@param` or `@return` or `@type` or etc.
//
// Java: JSDocInfo.Builder.hasAnyTypeRelatedTags()
func (b *JSDocInfo_Builder) HasAnyTypeRelatedTags() bool {
	return (b.bits&(JSDocInfo_Bit_CONSTRUCTOR.Mask()|JSDocInfo_Bit_INTERFACE.Mask()|JSDocInfo_Bit_ABSTRACT.Mask())) != 0 ||
		b.HasAnyParameters() ||
		b.GetProp(JSDocInfo_RETURN_TYPE) != nil ||
		b.GetProp(JSDocInfo_BASE_TYPE) != nil ||
		!b.IsPropEmpty(JSDocInfo_EXTENDED_INTERFACES) ||
		b.GetProp(JSDocInfo_LENDS_NAME) != nil ||
		b.GetProp(JSDocInfo_THIS_TYPE) != nil ||
		b.HasAnySingletonTypeTags()
}

// Java: JSDocInfo.Builder.hasAnyParameters()
func (b *JSDocInfo_Builder) HasAnyParameters() bool {
	params := b.GetProp(JSDocInfo_PARAMETERS)
	return params != nil && !jsdocGet[jutil.JavaMap](params).IsEmpty()
}

// Java: JSDocInfo.Builder.isPropEmpty(Property<? extends Collection<?>>)
func (b *JSDocInfo_Builder) IsPropEmpty(prop *JSDocInfo_Property) bool {
	c := b.GetProp(prop)
	return c == nil || jsdocGet[jutil.JavaCollection](c).IsEmpty()
}

// HasAnySingletonTypeTags reports whether the current doc info has any of
// the singleton type tags that may not appear with other type tags, like
// `@type` or `@typedef`.
//
// Java: JSDocInfo.Builder.hasAnySingletonTypeTags()
func (b *JSDocInfo_Builder) HasAnySingletonTypeTags() bool {
	return b.GetProp(JSDocInfo_TYPE) != nil ||
		b.GetProp(JSDocInfo_TYPEDEF_TYPE) != nil ||
		b.GetProp(JSDocInfo_ENUM_PARAMETER_TYPE) != nil
}

// HasAnySingletonSideEffectTags reports whether the current doc info has
// any of the singleton type tags that may not appear with other type tags,
// like `@type` or `@typedef`.
//
// Java: JSDocInfo.Builder.hasAnySingletonSideEffectTags()
func (b *JSDocInfo_Builder) HasAnySingletonSideEffectTags() bool {
	return b.CheckBit(JSDocInfo_Bit_NOSIDEEFFECTS) || !b.IsPropEmpty(JSDocInfo_MODIFIES)
}

// Java: JSDocInfo.Builder.isConstructorOrInterface()
func (b *JSDocInfo_Builder) IsConstructorOrInterface() bool {
	return (b.bits & (JSDocInfo_Bit_CONSTRUCTOR.Mask() | JSDocInfo_Bit_INTERFACE.Mask())) != 0
}

func (b *JSDocInfo_Builder) hasPropKey(prop *JSDocInfo_Property) bool {
	return (b.propMask & prop.mask) != 0
}

func (b *JSDocInfo_Builder) removeProp(prop *JSDocInfo_Property) any {
	if (b.propMask & prop.mask) == 0 {
		return nil
	}
	idx := bits.OnesCount64(uint64(b.propMask & (prop.mask - 1)))
	old := b.entries[idx]
	copy(b.entries[idx:], b.entries[idx+1:])
	b.entries[len(b.entries)-1] = nil
	b.entries = b.entries[:len(b.entries)-1]
	b.propMask &^= prop.mask
	return old
}

// Java: JSDocInfo.Builder.setProp(Property<T>, T)
func (b *JSDocInfo_Builder) SetProp(prop *JSDocInfo_Property, value any) {
	idx := bits.OnesCount64(uint64(b.propMask & (prop.mask - 1)))
	if (b.propMask & prop.mask) != 0 {
		b.entries[idx] = value
		return
	}
	b.propMask |= prop.mask
	if b.entries == nil {
		b.entries = b.entriesBuf[:0]
	}
	b.entries = append(b.entries, nil)
	copy(b.entries[idx+1:], b.entries[idx:])
	b.entries[idx] = value
}

// Java: JSDocInfo.Builder.getProp(Property<T>)
func (b *JSDocInfo_Builder) GetProp(prop *JSDocInfo_Property) any {
	if (b.propMask & prop.mask) == 0 {
		return nil
	}
	idx := bits.OnesCount64(uint64(b.propMask & (prop.mask - 1)))
	return b.entries[idx]
}

// Java: JSDocInfo.Builder.getPropWithDefault(Property<T>, Supplier<T>)
func (b *JSDocInfo_Builder) GetPropWithDefault(prop *JSDocInfo_Property, supplier func() any) any {
	value := b.GetProp(prop)
	if value == nil {
		value = supplier()
		b.SetProp(prop, value)
	}
	return value
}

// JSDocInfo_Builder_PutPropEntry puts key → value into the map property
// unless the key already has a (non-null) value.
//
// Go-only: language-adapter — a generic method, so a package function.
//
// Java: JSDocInfo.Builder.putPropEntry(Property<LinkedHashMap<K, V>>, K, V)
func JSDocInfo_Builder_PutPropEntry[V any](b *JSDocInfo_Builder, prop *JSDocInfo_Property, key string, value V) bool {
	m := jsdocGet[*jutil.LinkedHashMap[string, V]](b.GetPropWithDefault(prop, func() any {
		return jutil.NewLinkedHashMap[string, V]()
	}))
	return jutil.IsNil(any(m.PutIfAbsent(key, value)))
}

// JSDocInfo_Builder_PopulatePropEntry is putPropEntry, marking the builder
// populated if the entry was put.
//
// Go-only: language-adapter — a generic method, so a package function.
//
// Java: JSDocInfo.Builder.populatePropEntry(Property<LinkedHashMap<K, V>>, K, V)
func JSDocInfo_Builder_PopulatePropEntry[V any](b *JSDocInfo_Builder, prop *JSDocInfo_Property, key string, value V) bool {
	if JSDocInfo_Builder_PutPropEntry(b, prop, key, value) {
		b.populated = true
		return true
	}
	return false
}

// Java: JSDocInfo.Builder.populateProp(Property<T>, T)
func (b *JSDocInfo_Builder) PopulateProp(prop *JSDocInfo_Property, value any) bool {
	b.populated = true
	idx := bits.OnesCount64(uint64(b.propMask & (prop.mask - 1)))
	if (b.propMask & prop.mask) != 0 {
		if jutil.IsNil(b.entries[idx]) {
			b.entries[idx] = value
			return true
		}
		return false
	}
	b.propMask |= prop.mask
	if b.entries == nil {
		b.entries = b.entriesBuf[:0]
	}
	b.entries = append(b.entries, nil)
	copy(b.entries[idx+1:], b.entries[idx:])
	b.entries[idx] = value
	return true
}

// Java: JSDocInfo.Builder.checkBit(Bit)
func (b *JSDocInfo_Builder) CheckBit(bit JSDocInfo_Bit) bool {
	return (b.bits & bit.Mask()) != 0
}

// Java: JSDocInfo.Builder.setBit(Bit, boolean)
func (b *JSDocInfo_Builder) SetBit(bit JSDocInfo_Bit, value bool) {
	if value {
		b.bits |= bit.Mask()
	} else {
		b.bits &^= bit.Mask()
	}
}

// Java: JSDocInfo.Builder.populateBit(Bit, boolean)
func (b *JSDocInfo_Builder) PopulateBit(bit JSDocInfo_Bit, value bool) bool {
	if b.CheckBit(bit) != value {
		b.SetBit(bit, value)
		b.populated = true
		return b.populated
	}
	return false
}

// Java: JSDocInfo.internSet(Set<String>)
func JSDocInfo_InternSet(strings_ *jutil.LinkedHashSet[string]) *jutil.LinkedHashSet[string] {
	interned := jutil.NewLinkedHashSet[string]()
	for s := range strings_.All() {
		interned.Add(RhinoStringPool_AddOrGet(s))
	}
	return interned
}
