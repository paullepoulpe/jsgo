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
// Port of com/google/javascript/rhino/JSDocInfo.java (generated mechanically by j2go).

package rhino

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_js_doc_info_1 JSDocInfo_TypeTransform

func init() {
	lambdaSingleton_js_doc_info_1 = NewLambda_JSDocInfo_TypeTransform(func(mrp0_ JSTypeExpression) JSTypeExpression { return mrp0_.Copy() })
}

// Java class: com.google.javascript.rhino.JSDocInfo
type JSDocInfo struct {
	type0                JSTypeExpression `gson:"type0"`
	propertyValue0       any              `gson:"propertyValue0"`
	propertyValue1OrRest any              `gson:"propertyValue1OrRest"`
	propertyBits         int64            `gson:"propertyBits"`
	propertyKeysBitset   int64            `gson:"propertyKeysBitset"`
	_jsDocHandle         atomic.Uint32
	idHash_              atomic.Int32
}

func (this_ *JSDocInfo) BindSelf_(s any) {
	_ = s
}
func (this_ *JSDocInfo) IsJavaNil_() bool { return this_ == nil }
func (*JSDocInfo) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.JSDocInfo", "JSDocInfo"
}
func (this_ *JSDocInfo) JSDocHandlePtr_() *atomic.Uint32 { return &this_._jsDocHandle }
func (this_ *JSDocInfo) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*JSDocInfo) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.JSDocInfo", "com.google.javascript.rhino.JSDocInfo", func(v_ any) bool {
		return j2gort.InstanceOf[*JSDocInfo](v_)
	})
}

func StaticGet_JSDocInfo_SerialVersionUID() int64 { return int64(1) }

var static_JSDocInfo_VISIBILITY JSDocInfo_Property

func StaticGet_JSDocInfo_VISIBILITY() JSDocInfo_Property {
	clinit_JSDocInfo()
	return static_JSDocInfo_VISIBILITY
}
func StaticSet_JSDocInfo_VISIBILITY(v JSDocInfo_Property) {
	clinit_JSDocInfo()
	static_JSDocInfo_VISIBILITY = v
}

var static_JSDocInfo_TYPE *JSDocInfo_TypeProperty

func StaticGet_JSDocInfo_TYPE() *JSDocInfo_TypeProperty {
	clinit_JSDocInfo()
	return static_JSDocInfo_TYPE
}
func StaticSet_JSDocInfo_TYPE(v *JSDocInfo_TypeProperty) {
	clinit_JSDocInfo()
	static_JSDocInfo_TYPE = v
}

var static_JSDocInfo_RETURN_TYPE *JSDocInfo_TypeProperty

func StaticGet_JSDocInfo_RETURN_TYPE() *JSDocInfo_TypeProperty {
	clinit_JSDocInfo()
	return static_JSDocInfo_RETURN_TYPE
}
func StaticSet_JSDocInfo_RETURN_TYPE(v *JSDocInfo_TypeProperty) {
	clinit_JSDocInfo()
	static_JSDocInfo_RETURN_TYPE = v
}

var static_JSDocInfo_ENUM_PARAMETER_TYPE *JSDocInfo_TypeProperty

func StaticGet_JSDocInfo_ENUM_PARAMETER_TYPE() *JSDocInfo_TypeProperty {
	clinit_JSDocInfo()
	return static_JSDocInfo_ENUM_PARAMETER_TYPE
}
func StaticSet_JSDocInfo_ENUM_PARAMETER_TYPE(v *JSDocInfo_TypeProperty) {
	clinit_JSDocInfo()
	static_JSDocInfo_ENUM_PARAMETER_TYPE = v
}

var static_JSDocInfo_TYPEDEF_TYPE *JSDocInfo_TypeProperty

func StaticGet_JSDocInfo_TYPEDEF_TYPE() *JSDocInfo_TypeProperty {
	clinit_JSDocInfo()
	return static_JSDocInfo_TYPEDEF_TYPE
}
func StaticSet_JSDocInfo_TYPEDEF_TYPE(v *JSDocInfo_TypeProperty) {
	clinit_JSDocInfo()
	static_JSDocInfo_TYPEDEF_TYPE = v
}

var static_JSDocInfo_THIS_TYPE *JSDocInfo_TypeProperty

func StaticGet_JSDocInfo_THIS_TYPE() *JSDocInfo_TypeProperty {
	clinit_JSDocInfo()
	return static_JSDocInfo_THIS_TYPE
}
func StaticSet_JSDocInfo_THIS_TYPE(v *JSDocInfo_TypeProperty) {
	clinit_JSDocInfo()
	static_JSDocInfo_THIS_TYPE = v
}

var static_JSDocInfo_ORIGINAL_COMMENT_POSITION JSDocInfo_Property

func StaticGet_JSDocInfo_ORIGINAL_COMMENT_POSITION() JSDocInfo_Property {
	clinit_JSDocInfo()
	return static_JSDocInfo_ORIGINAL_COMMENT_POSITION
}
func StaticSet_JSDocInfo_ORIGINAL_COMMENT_POSITION(v JSDocInfo_Property) {
	clinit_JSDocInfo()
	static_JSDocInfo_ORIGINAL_COMMENT_POSITION = v
}

var static_JSDocInfo_ID_GENERATOR JSDocInfo_Property

func StaticGet_JSDocInfo_ID_GENERATOR() JSDocInfo_Property {
	clinit_JSDocInfo()
	return static_JSDocInfo_ID_GENERATOR
}
func StaticSet_JSDocInfo_ID_GENERATOR(v JSDocInfo_Property) {
	clinit_JSDocInfo()
	static_JSDocInfo_ID_GENERATOR = v
}

var static_JSDocInfo_BASE_TYPE *JSDocInfo_TypeProperty

func StaticGet_JSDocInfo_BASE_TYPE() *JSDocInfo_TypeProperty {
	clinit_JSDocInfo()
	return static_JSDocInfo_BASE_TYPE
}
func StaticSet_JSDocInfo_BASE_TYPE(v *JSDocInfo_TypeProperty) {
	clinit_JSDocInfo()
	static_JSDocInfo_BASE_TYPE = v
}

var static_JSDocInfo_EXTENDED_INTERFACES *JSDocInfo_TypeListProperty

func StaticGet_JSDocInfo_EXTENDED_INTERFACES() *JSDocInfo_TypeListProperty {
	clinit_JSDocInfo()
	return static_JSDocInfo_EXTENDED_INTERFACES
}
func StaticSet_JSDocInfo_EXTENDED_INTERFACES(v *JSDocInfo_TypeListProperty) {
	clinit_JSDocInfo()
	static_JSDocInfo_EXTENDED_INTERFACES = v
}

var static_JSDocInfo_IMPLEMENTED_INTERFACES *JSDocInfo_TypeListProperty

func StaticGet_JSDocInfo_IMPLEMENTED_INTERFACES() *JSDocInfo_TypeListProperty {
	clinit_JSDocInfo()
	return static_JSDocInfo_IMPLEMENTED_INTERFACES
}
func StaticSet_JSDocInfo_IMPLEMENTED_INTERFACES(v *JSDocInfo_TypeListProperty) {
	clinit_JSDocInfo()
	static_JSDocInfo_IMPLEMENTED_INTERFACES = v
}

var static_JSDocInfo_PARAMETERS *JSDocInfo_TypeMapProperty

func StaticGet_JSDocInfo_PARAMETERS() *JSDocInfo_TypeMapProperty {
	clinit_JSDocInfo()
	return static_JSDocInfo_PARAMETERS
}
func StaticSet_JSDocInfo_PARAMETERS(v *JSDocInfo_TypeMapProperty) {
	clinit_JSDocInfo()
	static_JSDocInfo_PARAMETERS = v
}

var static_JSDocInfo_THROWS_ANNOTATIONS JSDocInfo_Property

func StaticGet_JSDocInfo_THROWS_ANNOTATIONS() JSDocInfo_Property {
	clinit_JSDocInfo()
	return static_JSDocInfo_THROWS_ANNOTATIONS
}
func StaticSet_JSDocInfo_THROWS_ANNOTATIONS(v JSDocInfo_Property) {
	clinit_JSDocInfo()
	static_JSDocInfo_THROWS_ANNOTATIONS = v
}

var static_JSDocInfo_TEMPLATE_TYPE_NAMES *JSDocInfo_TypeMapProperty

func StaticGet_JSDocInfo_TEMPLATE_TYPE_NAMES() *JSDocInfo_TypeMapProperty {
	clinit_JSDocInfo()
	return static_JSDocInfo_TEMPLATE_TYPE_NAMES
}
func StaticSet_JSDocInfo_TEMPLATE_TYPE_NAMES(v *JSDocInfo_TypeMapProperty) {
	clinit_JSDocInfo()
	static_JSDocInfo_TEMPLATE_TYPE_NAMES = v
}

var static_JSDocInfo_TYPE_TRANSFORMATIONS *JSDocInfo_NodeMapProperty

func StaticGet_JSDocInfo_TYPE_TRANSFORMATIONS() *JSDocInfo_NodeMapProperty {
	clinit_JSDocInfo()
	return static_JSDocInfo_TYPE_TRANSFORMATIONS
}
func StaticSet_JSDocInfo_TYPE_TRANSFORMATIONS(v *JSDocInfo_NodeMapProperty) {
	clinit_JSDocInfo()
	static_JSDocInfo_TYPE_TRANSFORMATIONS = v
}

var static_JSDocInfo_DESCRIPTION JSDocInfo_Property

func StaticGet_JSDocInfo_DESCRIPTION() JSDocInfo_Property {
	clinit_JSDocInfo()
	return static_JSDocInfo_DESCRIPTION
}
func StaticSet_JSDocInfo_DESCRIPTION(v JSDocInfo_Property) {
	clinit_JSDocInfo()
	static_JSDocInfo_DESCRIPTION = v
}

var static_JSDocInfo_MEANING JSDocInfo_Property

func StaticGet_JSDocInfo_MEANING() JSDocInfo_Property {
	clinit_JSDocInfo()
	return static_JSDocInfo_MEANING
}
func StaticSet_JSDocInfo_MEANING(v JSDocInfo_Property) {
	clinit_JSDocInfo()
	static_JSDocInfo_MEANING = v
}

var static_JSDocInfo_ALTERNATE_MESSAGE_ID JSDocInfo_Property

func StaticGet_JSDocInfo_ALTERNATE_MESSAGE_ID() JSDocInfo_Property {
	clinit_JSDocInfo()
	return static_JSDocInfo_ALTERNATE_MESSAGE_ID
}
func StaticSet_JSDocInfo_ALTERNATE_MESSAGE_ID(v JSDocInfo_Property) {
	clinit_JSDocInfo()
	static_JSDocInfo_ALTERNATE_MESSAGE_ID = v
}

var static_JSDocInfo_DEPRECATION_REASON JSDocInfo_Property

func StaticGet_JSDocInfo_DEPRECATION_REASON() JSDocInfo_Property {
	clinit_JSDocInfo()
	return static_JSDocInfo_DEPRECATION_REASON
}
func StaticSet_JSDocInfo_DEPRECATION_REASON(v JSDocInfo_Property) {
	clinit_JSDocInfo()
	static_JSDocInfo_DEPRECATION_REASON = v
}

var static_JSDocInfo_LICENSE JSDocInfo_Property

func StaticGet_JSDocInfo_LICENSE() JSDocInfo_Property {
	clinit_JSDocInfo()
	return static_JSDocInfo_LICENSE
}
func StaticSet_JSDocInfo_LICENSE(v JSDocInfo_Property) {
	clinit_JSDocInfo()
	static_JSDocInfo_LICENSE = v
}

var static_JSDocInfo_SUPPRESSIONS JSDocInfo_Property

func StaticGet_JSDocInfo_SUPPRESSIONS() JSDocInfo_Property {
	clinit_JSDocInfo()
	return static_JSDocInfo_SUPPRESSIONS
}
func StaticSet_JSDocInfo_SUPPRESSIONS(v JSDocInfo_Property) {
	clinit_JSDocInfo()
	static_JSDocInfo_SUPPRESSIONS = v
}

var static_JSDocInfo_MODIFIES JSDocInfo_Property

func StaticGet_JSDocInfo_MODIFIES() JSDocInfo_Property {
	clinit_JSDocInfo()
	return static_JSDocInfo_MODIFIES
}
func StaticSet_JSDocInfo_MODIFIES(v JSDocInfo_Property) {
	clinit_JSDocInfo()
	static_JSDocInfo_MODIFIES = v
}

var static_JSDocInfo_LENDS_NAME *JSDocInfo_TypeProperty

func StaticGet_JSDocInfo_LENDS_NAME() *JSDocInfo_TypeProperty {
	clinit_JSDocInfo()
	return static_JSDocInfo_LENDS_NAME
}
func StaticSet_JSDocInfo_LENDS_NAME(v *JSDocInfo_TypeProperty) {
	clinit_JSDocInfo()
	static_JSDocInfo_LENDS_NAME = v
}

var static_JSDocInfo_CLOSURE_PRIMITIVE_ID JSDocInfo_Property

func StaticGet_JSDocInfo_CLOSURE_PRIMITIVE_ID() JSDocInfo_Property {
	clinit_JSDocInfo()
	return static_JSDocInfo_CLOSURE_PRIMITIVE_ID
}
func StaticSet_JSDocInfo_CLOSURE_PRIMITIVE_ID(v JSDocInfo_Property) {
	clinit_JSDocInfo()
	static_JSDocInfo_CLOSURE_PRIMITIVE_ID = v
}

var static_JSDocInfo_SOURCE_COMMENT JSDocInfo_Property

func StaticGet_JSDocInfo_SOURCE_COMMENT() JSDocInfo_Property {
	clinit_JSDocInfo()
	return static_JSDocInfo_SOURCE_COMMENT
}
func StaticSet_JSDocInfo_SOURCE_COMMENT(v JSDocInfo_Property) {
	clinit_JSDocInfo()
	static_JSDocInfo_SOURCE_COMMENT = v
}

var static_JSDocInfo_MARKERS JSDocInfo_Property

func StaticGet_JSDocInfo_MARKERS() JSDocInfo_Property {
	clinit_JSDocInfo()
	return static_JSDocInfo_MARKERS
}
func StaticSet_JSDocInfo_MARKERS(v JSDocInfo_Property) {
	clinit_JSDocInfo()
	static_JSDocInfo_MARKERS = v
}

var static_JSDocInfo_PARAMETER_DESCRIPTIONS *JSDocInfo_StringMapProperty

func StaticGet_JSDocInfo_PARAMETER_DESCRIPTIONS() *JSDocInfo_StringMapProperty {
	clinit_JSDocInfo()
	return static_JSDocInfo_PARAMETER_DESCRIPTIONS
}
func StaticSet_JSDocInfo_PARAMETER_DESCRIPTIONS(v *JSDocInfo_StringMapProperty) {
	clinit_JSDocInfo()
	static_JSDocInfo_PARAMETER_DESCRIPTIONS = v
}

var static_JSDocInfo_BLOCK_DESCRIPTION JSDocInfo_Property

func StaticGet_JSDocInfo_BLOCK_DESCRIPTION() JSDocInfo_Property {
	clinit_JSDocInfo()
	return static_JSDocInfo_BLOCK_DESCRIPTION
}
func StaticSet_JSDocInfo_BLOCK_DESCRIPTION(v JSDocInfo_Property) {
	clinit_JSDocInfo()
	static_JSDocInfo_BLOCK_DESCRIPTION = v
}

var static_JSDocInfo_FILEOVERVIEW_DESCRIPTION JSDocInfo_Property

func StaticGet_JSDocInfo_FILEOVERVIEW_DESCRIPTION() JSDocInfo_Property {
	clinit_JSDocInfo()
	return static_JSDocInfo_FILEOVERVIEW_DESCRIPTION
}
func StaticSet_JSDocInfo_FILEOVERVIEW_DESCRIPTION(v JSDocInfo_Property) {
	clinit_JSDocInfo()
	static_JSDocInfo_FILEOVERVIEW_DESCRIPTION = v
}

var static_JSDocInfo_RETURN_DESCRIPTION JSDocInfo_Property

func StaticGet_JSDocInfo_RETURN_DESCRIPTION() JSDocInfo_Property {
	clinit_JSDocInfo()
	return static_JSDocInfo_RETURN_DESCRIPTION
}
func StaticSet_JSDocInfo_RETURN_DESCRIPTION(v JSDocInfo_Property) {
	clinit_JSDocInfo()
	static_JSDocInfo_RETURN_DESCRIPTION = v
}

var static_JSDocInfo_ENHANCED_NAMESPACE JSDocInfo_Property

func StaticGet_JSDocInfo_ENHANCED_NAMESPACE() JSDocInfo_Property {
	clinit_JSDocInfo()
	return static_JSDocInfo_ENHANCED_NAMESPACE
}
func StaticSet_JSDocInfo_ENHANCED_NAMESPACE(v JSDocInfo_Property) {
	clinit_JSDocInfo()
	static_JSDocInfo_ENHANCED_NAMESPACE = v
}

var static_JSDocInfo_MODS JSDocInfo_Property

func StaticGet_JSDocInfo_MODS() JSDocInfo_Property {
	clinit_JSDocInfo()
	return static_JSDocInfo_MODS
}
func StaticSet_JSDocInfo_MODS(v JSDocInfo_Property) {
	clinit_JSDocInfo()
	static_JSDocInfo_MODS = v
}

var static_JSDocInfo_TS_TYPES JSDocInfo_Property

func StaticGet_JSDocInfo_TS_TYPES() JSDocInfo_Property {
	clinit_JSDocInfo()
	return static_JSDocInfo_TS_TYPES
}
func StaticSet_JSDocInfo_TS_TYPES(v JSDocInfo_Property) {
	clinit_JSDocInfo()
	static_JSDocInfo_TS_TYPES = v
}

var static_JSDocInfo_CLOSURE_UNAWARE_CONFIG_VALUE JSDocInfo_Property

func StaticGet_JSDocInfo_CLOSURE_UNAWARE_CONFIG_VALUE() JSDocInfo_Property {
	clinit_JSDocInfo()
	return static_JSDocInfo_CLOSURE_UNAWARE_CONFIG_VALUE
}
func StaticSet_JSDocInfo_CLOSURE_UNAWARE_CONFIG_VALUE(v JSDocInfo_Property) {
	clinit_JSDocInfo()
	static_JSDocInfo_CLOSURE_UNAWARE_CONFIG_VALUE = v
}

var static_JSDocInfo_AUTHORS JSDocInfo_Property

func StaticGet_JSDocInfo_AUTHORS() JSDocInfo_Property {
	clinit_JSDocInfo()
	return static_JSDocInfo_AUTHORS
}
func StaticSet_JSDocInfo_AUTHORS(v JSDocInfo_Property) {
	clinit_JSDocInfo()
	static_JSDocInfo_AUTHORS = v
}

var static_JSDocInfo_SEES JSDocInfo_Property

func StaticGet_JSDocInfo_SEES() JSDocInfo_Property {
	clinit_JSDocInfo()
	return static_JSDocInfo_SEES
}
func StaticSet_JSDocInfo_SEES(v JSDocInfo_Property) {
	clinit_JSDocInfo()
	static_JSDocInfo_SEES = v
}

var static_JSDocInfo_SINGLE_TYPE_EXPR_PROPERTIES_MASK int64

func StaticGet_JSDocInfo_SINGLE_TYPE_EXPR_PROPERTIES_MASK() int64 {
	clinit_JSDocInfo()
	return static_JSDocInfo_SINGLE_TYPE_EXPR_PROPERTIES_MASK
}
func StaticSet_JSDocInfo_SINGLE_TYPE_EXPR_PROPERTIES_MASK(v int64) {
	clinit_JSDocInfo()
	static_JSDocInfo_SINGLE_TYPE_EXPR_PROPERTIES_MASK = v
}

var static_JSDocInfo_TYPE_EXPR_PROPERTIES_MASK int64

func StaticGet_JSDocInfo_TYPE_EXPR_PROPERTIES_MASK() int64 {
	clinit_JSDocInfo()
	return static_JSDocInfo_TYPE_EXPR_PROPERTIES_MASK
}
func StaticSet_JSDocInfo_TYPE_EXPR_PROPERTIES_MASK(v int64) {
	clinit_JSDocInfo()
	static_JSDocInfo_TYPE_EXPR_PROPERTIES_MASK = v
}

var static_JSDocInfo_EQUIVALENCE_IGNORED_BITS int64

func StaticGet_JSDocInfo_EQUIVALENCE_IGNORED_BITS() int64 {
	clinit_JSDocInfo()
	return static_JSDocInfo_EQUIVALENCE_IGNORED_BITS
}
func StaticSet_JSDocInfo_EQUIVALENCE_IGNORED_BITS(v int64) {
	clinit_JSDocInfo()
	static_JSDocInfo_EQUIVALENCE_IGNORED_BITS = v
}

var clinitDone_JSDocInfo uint32
var clinitState_JSDocInfo j2gort.ClinitState

func clinit_JSDocInfo() {
	if atomic.LoadUint32(&clinitDone_JSDocInfo) != 0 {
		return
	}
	clinitSlow_JSDocInfo()
}

func clinitSlow_JSDocInfo() {
	if !j2gort.EnterClinit(&clinitState_JSDocInfo, &clinitDone_JSDocInfo) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_JSDocInfo, &clinitDone_JSDocInfo)
	static_JSDocInfo_VISIBILITY = NewJSDocInfo_Property("visibility")
	static_JSDocInfo_TYPE = NewJSDocInfo_TypeProperty("type")
	static_JSDocInfo_RETURN_TYPE = NewJSDocInfo_TypeProperty("returnType")
	static_JSDocInfo_ENUM_PARAMETER_TYPE = NewJSDocInfo_TypeProperty("enumParameterType")
	static_JSDocInfo_TYPEDEF_TYPE = NewJSDocInfo_TypeProperty("typedefType")
	static_JSDocInfo_THIS_TYPE = NewJSDocInfo_TypeProperty("thisType")
	static_JSDocInfo_ORIGINAL_COMMENT_POSITION = NewJSDocInfo_Property("originalCommentPosition")
	static_JSDocInfo_ID_GENERATOR = NewJSDocInfo_Property("idGenerator")
	static_JSDocInfo_BASE_TYPE = NewJSDocInfo_TypeProperty("baseType")
	static_JSDocInfo_EXTENDED_INTERFACES = NewJSDocInfo_TypeListProperty("extendedInterfaces")
	static_JSDocInfo_IMPLEMENTED_INTERFACES = NewJSDocInfo_TypeListProperty("extendedInterfaces")
	static_JSDocInfo_PARAMETERS = NewJSDocInfo_TypeMapProperty("parameters")
	static_JSDocInfo_THROWS_ANNOTATIONS = NewJSDocInfo_Property("throws")
	static_JSDocInfo_TEMPLATE_TYPE_NAMES = NewJSDocInfo_TypeMapProperty("templateTypeNames")
	static_JSDocInfo_TYPE_TRANSFORMATIONS = NewJSDocInfo_NodeMapProperty("typeTransformations")
	static_JSDocInfo_DESCRIPTION = NewJSDocInfo_Property("description")
	static_JSDocInfo_MEANING = NewJSDocInfo_Property("meaning")
	static_JSDocInfo_ALTERNATE_MESSAGE_ID = NewJSDocInfo_Property("alternateMessageId")
	static_JSDocInfo_DEPRECATION_REASON = NewJSDocInfo_Property("deprecationReason")
	static_JSDocInfo_LICENSE = NewJSDocInfo_Property("license")
	static_JSDocInfo_SUPPRESSIONS = NewJSDocInfo_Property("suppressions")
	static_JSDocInfo_MODIFIES = NewJSDocInfo_Property("modifies")
	static_JSDocInfo_LENDS_NAME = NewJSDocInfo_TypeProperty("lendsName")
	static_JSDocInfo_CLOSURE_PRIMITIVE_ID = NewJSDocInfo_Property("closurePrimitiveId")
	static_JSDocInfo_SOURCE_COMMENT = NewJSDocInfo_Property("sourceComment")
	static_JSDocInfo_MARKERS = NewJSDocInfo_MarkerListProperty("markers")
	static_JSDocInfo_PARAMETER_DESCRIPTIONS = NewJSDocInfo_StringMapProperty("parameterDescriptions")
	static_JSDocInfo_BLOCK_DESCRIPTION = NewJSDocInfo_Property("blockDescription")
	static_JSDocInfo_FILEOVERVIEW_DESCRIPTION = NewJSDocInfo_Property("fileoverviewDescription")
	static_JSDocInfo_RETURN_DESCRIPTION = NewJSDocInfo_Property("returnDescription")
	static_JSDocInfo_ENHANCED_NAMESPACE = NewJSDocInfo_Property("enhance")
	static_JSDocInfo_MODS = NewJSDocInfo_Property("mods")
	static_JSDocInfo_TS_TYPES = NewJSDocInfo_Property("tsType")
	static_JSDocInfo_CLOSURE_UNAWARE_CONFIG_VALUE = NewJSDocInfo_Property("perFileClosureUnawareMode")
	static_JSDocInfo_AUTHORS = NewJSDocInfo_Property("authors")
	static_JSDocInfo_SEES = NewJSDocInfo_Property("sees")
	static_JSDocInfo_SINGLE_TYPE_EXPR_PROPERTIES_MASK = ((((((StaticGet_JSDocInfo_TYPE().GetJSDocInfo_PropertyBase().mask | StaticGet_JSDocInfo_RETURN_TYPE().GetJSDocInfo_PropertyBase().mask) | StaticGet_JSDocInfo_ENUM_PARAMETER_TYPE().GetJSDocInfo_PropertyBase().mask) | StaticGet_JSDocInfo_TYPEDEF_TYPE().GetJSDocInfo_PropertyBase().mask) | StaticGet_JSDocInfo_THIS_TYPE().GetJSDocInfo_PropertyBase().mask) | StaticGet_JSDocInfo_BASE_TYPE().GetJSDocInfo_PropertyBase().mask) | StaticGet_JSDocInfo_LENDS_NAME().GetJSDocInfo_PropertyBase().mask)
	static_JSDocInfo_TYPE_EXPR_PROPERTIES_MASK = ((((StaticGet_JSDocInfo_SINGLE_TYPE_EXPR_PROPERTIES_MASK() | StaticGet_JSDocInfo_EXTENDED_INTERFACES().GetJSDocInfo_PropertyBase().mask) | StaticGet_JSDocInfo_IMPLEMENTED_INTERFACES().GetJSDocInfo_PropertyBase().mask) | StaticGet_JSDocInfo_PARAMETERS().GetJSDocInfo_PropertyBase().mask) | StaticGet_JSDocInfo_TEMPLATE_TYPE_NAMES().GetJSDocInfo_PropertyBase().mask)
	static_JSDocInfo_EQUIVALENCE_IGNORED_BITS = (StaticGet_JSDocInfo_Bit_INCLUDE_DOCUMENTATION().mask | StaticGet_JSDocInfo_Bit_INLINE_TYPE().mask)
}

func (this_ *JSDocInfo) initFields_JSDocInfo() {
	_ = this_
}

// Java constructor: JSDocInfo(com.google.javascript.rhino.JSDocInfo.Builder builder)
func NewJSDocInfo(builder *JSDocInfo_Builder) *JSDocInfo {
	clinit_JSDocInfo()
	this_ := &JSDocInfo{}
	InitJSDocInfo(this_, builder)
	return this_
}

func InitJSDocInfo(this_ *JSDocInfo, builder *JSDocInfo_Builder) {
	_ = builder
	this_.initFields_JSDocInfo()
	this_.propertyBits = builder.bits
	var singleKeys int64 = builder.singleTypeKeys
	_ = singleKeys
	var propKeys int64 = builder.propKeys
	_ = propKeys
	if (singleKeys == int64(0)) && (propKeys == int64(0)) {
		this_.propertyKeysBitset = int64(0)
		this_.type0 = (JSTypeExpression{})
		this_.propertyValue0 = nil
		this_.propertyValue1OrRest = nil
		return
	}
	if singleKeys != int64(0) {
		this_.type0 = builder.typeVal0
	} else {
		this_.type0 = (JSTypeExpression{})
	}
	var extraSingleKeys int64 = (singleKeys & (singleKeys - int64(1)))
	_ = extraSingleKeys
	var propValues []any = builder.propValues
	_ = propValues
	if (propValues == nil) && (extraSingleKeys == int64(0)) {
		var keys int64 = int64(0)
		_ = keys
		var count int32 = int32(0)
		_ = count
		var kept0 any = nil
		_ = kept0
		var kept1 any = nil
		_ = kept1
		var kept2 any = nil
		_ = kept2
		var kept3 any = nil
		_ = kept3
		var rem int64 = propKeys
		_ = rem
		var inlineIdx int32 = int32(0)
		_ = inlineIdx
		for rem != int64(0) {
			var bit int32 = j2gort.BoxLong_NumberOfTrailingZeros(rem)
			_ = bit
			rem = (rem & (rem - int64(1)))
			var prop JSDocInfo_Property = StaticGet_JSDocInfo_Property_Values()[bit]
			_ = prop
			var value any = builder.GetInlineValAt(func() int32 { old_1 := inlineIdx; nw_2 := old_1 + 1; inlineIdx = nw_2; return old_1 }())
			_ = value
			if !(prop.IsDefault(value)) {
				if !(!(j2gort.InstanceOf[[]any](value)) && (value != nil)) {
					j2gort.CheckState(false)
				}
				keys = (keys | prop.GetJSDocInfo_PropertyBase().mask)
				switch count {
				case int32(0):
					kept0 = value
				case int32(1):
					kept1 = value
				case int32(2):
					kept2 = value
				default:
					kept3 = value
				}
				count = count + 1
			}
		}
		this_.propertyKeysBitset = (keys | singleKeys)
		this_.propertyValue0 = kept0
		if count <= int32(1) {
			this_.propertyValue1OrRest = nil
		} else if count == int32(2) {
			this_.propertyValue1OrRest = kept1
		} else if count == int32(3) {
			this_.propertyValue1OrRest = []any{kept1, kept2}
		} else {
			this_.propertyValue1OrRest = []any{kept1, kept2, kept3}
		}
		return
	}
	var validObjKeys int64 = int64(0)
	_ = validObjKeys
	if propKeys != int64(0) {
		var rem_2 int64 = propKeys
		_ = rem_2
		var inlineIdx_2 int32 = int32(0)
		_ = inlineIdx_2
		for rem_2 != int64(0) {
			var bit_2 int32 = j2gort.BoxLong_NumberOfTrailingZeros(rem_2)
			_ = bit_2
			rem_2 = (rem_2 & (rem_2 - int64(1)))
			var prop_2 JSDocInfo_Property = StaticGet_JSDocInfo_Property_Values()[bit_2]
			_ = prop_2
			var value_2 any
			_ = value_2
			if propValues != nil {
				value_2 = propValues[bit_2]
			} else {
				value_2 = builder.GetInlineValAt(func() int32 { old_3 := inlineIdx_2; nw_4 := old_3 + 1; inlineIdx_2 = nw_4; return old_3 }())
			}
			if !(prop_2.IsDefault(value_2)) {
				if !(!(j2gort.InstanceOf[[]any](value_2)) && (value_2 != nil)) {
					j2gort.CheckState(false)
				}
				validObjKeys = (validObjKeys | prop_2.GetJSDocInfo_PropertyBase().mask)
			}
		}
	}
	var nonType0Keys int64 = (validObjKeys | extraSingleKeys)
	_ = nonType0Keys
	this_.propertyKeysBitset = (validObjKeys | singleKeys)
	var count_2 int32 = j2gort.BoxLong_BitCount(nonType0Keys)
	_ = count_2
	if count_2 == int32(0) {
		this_.propertyValue0 = nil
		this_.propertyValue1OrRest = nil
		return
	}
	var v0 any = nil
	_ = v0
	var v1 any = nil
	_ = v1
	var rest []any
	_ = rest
	if count_2 >= int32(3) {
		rest = make([]any, (count_2 - int32(1)))
	} else {
		rest = nil
	}
	var rem_3 int64 = nonType0Keys
	_ = rem_3
	var idx int32 = int32(0)
	_ = idx
	for rem_3 != int64(0) {
		var bit_3 int32 = j2gort.BoxLong_NumberOfTrailingZeros(rem_3)
		_ = bit_3
		rem_3 = (rem_3 & (rem_3 - int64(1)))
		var bitMask int64 = javalang.Shl64(int64(1), int(bit_3))
		_ = bitMask
		var value_3 any = nil
		_ = value_3
		if (bitMask & extraSingleKeys) != int64(0) {
			value_3 = j2gort.BoxNullableVal(builder.GetTypeValAt(j2gort.BoxLong_BitCount((singleKeys & (bitMask - int64(1))))))
		} else if propValues != nil {
			value_3 = propValues[bit_3]
		} else {
			value_3 = builder.GetInlineValAt(j2gort.BoxLong_BitCount((propKeys & (bitMask - int64(1)))))
		}
		if idx == int32(0) {
			v0 = value_3
		} else if rest != nil {
			{
				a_5 := rest
				i_6 := (idx - int32(1))
				a_5[i_6] = value_3
			}
		} else {
			v1 = value_3
		}
		idx = idx + 1
	}
	this_.propertyValue0 = v0
	if rest != nil {
		this_.propertyValue1OrRest = j2gort.SliceToAny(rest)
	} else {
		this_.propertyValue1OrRest = v1
	}
}

// Java method: JSDocInfo.checkBit(com.google.javascript.rhino.JSDocInfo.Bit bit)
func (this_ *JSDocInfo) priv_JSDocInfo_checkBit_7975(bit *JSDocInfo_Bit) bool {
	_ = bit
	return ((this_.propertyBits & bit.mask) != int64(int32(0)))
}

// Java static method: JSDocInfo.builder()
func JSDocInfo_S_Builder() *JSDocInfo_Builder {
	clinit_JSDocInfo()
	return NewJSDocInfo_Builder()
}

// Java method: JSDocInfo.toBuilder()
func (this_ *JSDocInfo) ToBuilder_0() *JSDocInfo_Builder {
	return this_.priv_JSDocInfo_toBuilder_7977(nil)
}

// Java method: JSDocInfo.toBuilder(com.google.javascript.rhino.JSDocInfo.@org.jspecify.annotations.Nullable TypeTransform transform)
func (this_ *JSDocInfo) priv_JSDocInfo_toBuilder_7977(transform JSDocInfo_TypeTransform) *JSDocInfo_Builder {
	_ = transform
	var val_7 JSDocInfo_Builder
	InitJSDocInfo_Builder(&val_7)
	var builder *JSDocInfo_Builder = &val_7
	_ = builder
	{
		r_8 := builder
		r_8.bits = (this_.propertyBits & ^(StaticGet_JSDocInfo_Bit_INLINE_TYPE().mask))
	}
	{
		r_9 := builder
		r_9.populated = true
	}
	var keys int64 = this_.propertyKeysBitset
	_ = keys
	if keys == int64(0) {
		return builder
	}
	var singleKeys int64 = (keys & StaticGet_JSDocInfo_SINGLE_TYPE_EXPR_PROPERTIES_MASK())
	_ = singleKeys
	if singleKeys != int64(0) {
		var typeIdx int32 = int32(0)
		_ = typeIdx
		var remSingle int64 = singleKeys
		_ = remSingle
		var actualSingleKeys int64 = int64(0)
		_ = actualSingleKeys
		for remSingle != int64(0) {
			var low int32 = j2gort.BoxLong_NumberOfTrailingZeros(remSingle)
			_ = low
			remSingle = (remSingle & (remSingle - int64(1)))
			var expr JSTypeExpression = this_.priv_JSDocInfo_getSingleTypeForBit_8107(low)
			_ = expr
			var cloned JSTypeExpression
			_ = cloned
			if (expr != (JSTypeExpression{})) && (transform != nil && !transform.IsJavaNil_()) {
				cloned = transform.Apply(expr)
			} else {
				cloned = expr
			}
			if cloned != (JSTypeExpression{}) {
				builder.SetTypeValAt(func() int32 { old_10 := typeIdx; nw_11 := old_10 + 1; typeIdx = nw_11; return old_10 }(), cloned)
				actualSingleKeys = (actualSingleKeys | (javalang.Shl64(int64(1), int(low))))
			}
		}
		{
			r_12 := builder
			r_12.singleTypeKeys = actualSingleKeys
		}
	}
	var objKeys int64 = (keys & ^(StaticGet_JSDocInfo_SINGLE_TYPE_EXPR_PROPERTIES_MASK()))
	_ = objKeys
	if objKeys == int64(0) {
		return builder
	}
	if j2gort.BoxLong_BitCount(objKeys) <= int32(4) {
		var actualKeys int64 = int64(0)
		_ = actualKeys
		var dstIndex int32 = int32(0)
		_ = dstIndex
		for objKeys != int64(0) {
			var low_2 int32 = j2gort.BoxLong_NumberOfTrailingZeros(objKeys)
			_ = low_2
			objKeys = (objKeys & (objKeys - int64(1)))
			var prop JSDocInfo_Property = StaticGet_JSDocInfo_Property_Values()[low_2]
			_ = prop
			var cloned_2 any = prop.Clone(this_.priv_JSDocInfo_getNonType0ValueForBit_8108(low_2), transform)
			_ = cloned_2
			if cloned_2 != nil {
				builder.SetInlineValAt(func() int32 { old_13 := dstIndex; nw_14 := old_13 + 1; dstIndex = nw_14; return old_13 }(), cloned_2)
				actualKeys = (actualKeys | (javalang.Shl64(int64(1), int(low_2))))
			}
		}
		{
			r_15 := builder
			r_15.propKeys = actualKeys
		}
		return builder
	}
	var vals []any = make([]any, StaticGet_JSDocInfo_Property_BitCounter())
	_ = vals
	var actualKeys_2 int64 = int64(0)
	_ = actualKeys_2
	for objKeys != int64(0) {
		var low_3 int32 = j2gort.BoxLong_NumberOfTrailingZeros(objKeys)
		_ = low_3
		objKeys = (objKeys & (objKeys - int64(1)))
		var prop_2 JSDocInfo_Property = StaticGet_JSDocInfo_Property_Values()[low_3]
		_ = prop_2
		var cloned_3 any = prop_2.Clone(this_.priv_JSDocInfo_getNonType0ValueForBit_8108(low_3), transform)
		_ = cloned_3
		if cloned_3 != nil {
			{
				a_16 := vals
				i_17 := low_3
				a_16[i_17] = cloned_3
			}
			actualKeys_2 = (actualKeys_2 | (javalang.Shl64(int64(1), int(low_3))))
		}
	}
	{
		r_18 := builder
		r_18.propKeys = actualKeys_2
	}
	{
		r_19 := builder
		r_19.propValues = vals
	}
	return builder
}

// Java method: JSDocInfo.clone()
func (this_ *JSDocInfo) Clone_0() *JSDocInfo {
	return this_.Clone_1(false)
}

// Java method: JSDocInfo.cloneAndReplaceTypeNames(java.util.Set<java.lang.String> names)
func (this_ *JSDocInfo) CloneAndReplaceTypeNames(names j2gort.Ext_java_util_Set) *JSDocInfo {
	_ = names
	return this_.priv_JSDocInfo_toBuilder_7977(NewLambda_JSDocInfo_TypeTransform(func(type_ JSTypeExpression) JSTypeExpression {
		_ = type_
		return type_.ReplaceNamesWithUnknownType(names)
	})).Build_0()
}

// Java method: JSDocInfo.clone(boolean cloneTypeNodes)
func (this_ *JSDocInfo) Clone_1(cloneTypeNodes bool) *JSDocInfo {
	_ = cloneTypeNodes
	if cloneTypeNodes {
		return this_.priv_JSDocInfo_toBuilder_7977(lambdaSingleton_js_doc_info_1).Build_0()
	} else {
		return this_.ToBuilder_0().Build_0()
	}
}

// Java static method: JSDocInfo.areEquivalent(com.google.javascript.rhino.JSDocInfo jsDoc1, com.google.javascript.rhino.JSDocInfo jsDoc2)
func JSDocInfo_AreEquivalent(jsDoc1 *JSDocInfo, jsDoc2 *JSDocInfo) bool {
	clinit_JSDocInfo()
	_ = jsDoc1
	_ = jsDoc2
	if (jsDoc1 == nil) != (jsDoc2 == nil) {
		return false
	} else if jsDoc1 == nil {
		return true
	}
	if ((jsDoc1.propertyBits ^ jsDoc2.propertyBits) & ^(StaticGet_JSDocInfo_EQUIVALENCE_IGNORED_BITS())) != int64(int32(0)) {
		return false
	}
	if jsDoc1.propertyKeysBitset != jsDoc2.propertyKeysBitset {
		return false
	}
	var singleBits int64 = (jsDoc1.propertyKeysBitset & StaticGet_JSDocInfo_SINGLE_TYPE_EXPR_PROPERTIES_MASK())
	_ = singleBits
	for singleBits != int64(0) {
		var low int32 = j2gort.BoxLong_NumberOfTrailingZeros(singleBits)
		_ = low
		singleBits = (singleBits & (singleBits - int64(1)))
		var a JSTypeExpression = jsDoc1.priv_JSDocInfo_getSingleTypeForBit_8107(low)
		_ = a
		var b JSTypeExpression = jsDoc2.priv_JSDocInfo_getSingleTypeForBit_8107(low)
		_ = b
		if ((a == (JSTypeExpression{})) != (b == (JSTypeExpression{}))) || ((a != (JSTypeExpression{})) && !(a.IsEquivalentTo(b))) {
			return false
		}
	}
	var objBits int64 = (jsDoc1.propertyKeysBitset & ^(StaticGet_JSDocInfo_SINGLE_TYPE_EXPR_PROPERTIES_MASK()))
	_ = objBits
	for objBits != int64(0) {
		var low_2 int32 = j2gort.BoxLong_NumberOfTrailingZeros(objBits)
		_ = low_2
		objBits = (objBits & (objBits - int64(1)))
		var prop JSDocInfo_Property = StaticGet_JSDocInfo_Property_Values()[low_2]
		_ = prop
		var a_2 any = jsDoc1.priv_JSDocInfo_getNonType0ValueForBit_8108(low_2)
		_ = a_2
		var b_2 any = jsDoc2.priv_JSDocInfo_getNonType0ValueForBit_8108(low_2)
		_ = b_2
		if ((a_2 == nil) != (b_2 == nil)) || ((a_2 != nil) && !(prop.EqualValues(a_2, b_2))) {
			return false
		}
	}
	return true
}

// Java method: JSDocInfo.isDocumentationIncluded()
func (this_ *JSDocInfo) IsDocumentationIncluded() bool {
	return this_.priv_JSDocInfo_checkBit_7975(StaticGet_JSDocInfo_Bit_INCLUDE_DOCUMENTATION())
}

// Java method: JSDocInfo.isAnyIdGenerator()
func (this_ *JSDocInfo) IsAnyIdGenerator() bool {
	return ((this_.propertyKeysBitset & StaticGet_JSDocInfo_ID_GENERATOR().GetJSDocInfo_PropertyBase().mask) != int64(int32(0)))
}

// Java method: JSDocInfo.isConsistentIdGenerator()
func (this_ *JSDocInfo) IsConsistentIdGenerator() bool {
	return (j2gort.Cast[JSDocInfo_IdGenerator](StaticGet_JSDocInfo_ID_GENERATOR().Get(this_)) == StaticGet_JSDocInfo_IdGenerator_CONSISTENT())
}

// Java method: JSDocInfo.isStableIdGenerator()
func (this_ *JSDocInfo) IsStableIdGenerator() bool {
	return (j2gort.Cast[JSDocInfo_IdGenerator](StaticGet_JSDocInfo_ID_GENERATOR().Get(this_)) == StaticGet_JSDocInfo_IdGenerator_STABLE())
}

// Java method: JSDocInfo.isXidGenerator()
func (this_ *JSDocInfo) IsXidGenerator() bool {
	return (j2gort.Cast[JSDocInfo_IdGenerator](StaticGet_JSDocInfo_ID_GENERATOR().Get(this_)) == StaticGet_JSDocInfo_IdGenerator_XID())
}

// Java method: JSDocInfo.isMappedIdGenerator()
func (this_ *JSDocInfo) IsMappedIdGenerator() bool {
	return (j2gort.Cast[JSDocInfo_IdGenerator](StaticGet_JSDocInfo_ID_GENERATOR().Get(this_)) == StaticGet_JSDocInfo_IdGenerator_MAPPED())
}

// Java method: JSDocInfo.isIdGenerator()
func (this_ *JSDocInfo) IsIdGenerator() bool {
	return (j2gort.Cast[JSDocInfo_IdGenerator](StaticGet_JSDocInfo_ID_GENERATOR().Get(this_)) == StaticGet_JSDocInfo_IdGenerator_UNIQUE())
}

// Java method: JSDocInfo.isConstant()
func (this_ *JSDocInfo) IsConstant() bool {
	return (((this_.propertyBits & ((StaticGet_JSDocInfo_Bit_CONST().mask | StaticGet_JSDocInfo_Bit_DEFINE().mask) | StaticGet_JSDocInfo_Bit_FINAL().mask)) != int64(int32(0))) || ((this_.propertyKeysBitset & StaticGet_JSDocInfo_DESCRIPTION().GetJSDocInfo_PropertyBase().mask) != int64(int32(0))))
}

// Java method: JSDocInfo.hasConstAnnotation()
func (this_ *JSDocInfo) HasConstAnnotation() bool {
	return this_.priv_JSDocInfo_checkBit_7975(StaticGet_JSDocInfo_Bit_CONST())
}

// Java method: JSDocInfo.isFinal()
func (this_ *JSDocInfo) IsFinal() bool {
	return this_.priv_JSDocInfo_checkBit_7975(StaticGet_JSDocInfo_Bit_FINAL())
}

// Java method: JSDocInfo.isConstructor()
func (this_ *JSDocInfo) IsConstructor() bool {
	return this_.priv_JSDocInfo_checkBit_7975(StaticGet_JSDocInfo_Bit_CONSTRUCTOR())
}

// Java method: JSDocInfo.isAbstract()
func (this_ *JSDocInfo) IsAbstract() bool {
	return this_.priv_JSDocInfo_checkBit_7975(StaticGet_JSDocInfo_Bit_ABSTRACT())
}

// Java method: JSDocInfo.usesImplicitMatch()
func (this_ *JSDocInfo) UsesImplicitMatch() bool {
	return this_.priv_JSDocInfo_checkBit_7975(StaticGet_JSDocInfo_Bit_RECORD())
}

// Java method: JSDocInfo.makesUnrestricted()
func (this_ *JSDocInfo) MakesUnrestricted() bool {
	return this_.priv_JSDocInfo_checkBit_7975(StaticGet_JSDocInfo_Bit_UNRESTRICTED())
}

// Java method: JSDocInfo.makesStructs()
func (this_ *JSDocInfo) MakesStructs() bool {
	return this_.priv_JSDocInfo_checkBit_7975(StaticGet_JSDocInfo_Bit_STRUCT())
}

// Java method: JSDocInfo.makesDicts()
func (this_ *JSDocInfo) MakesDicts() bool {
	return this_.priv_JSDocInfo_checkBit_7975(StaticGet_JSDocInfo_Bit_DICT())
}

// Java method: JSDocInfo.isDefine()
func (this_ *JSDocInfo) IsDefine() bool {
	return this_.priv_JSDocInfo_checkBit_7975(StaticGet_JSDocInfo_Bit_DEFINE())
}

// Java method: JSDocInfo.isOverride()
func (this_ *JSDocInfo) IsOverride() bool {
	return this_.priv_JSDocInfo_checkBit_7975(StaticGet_JSDocInfo_Bit_OVERRIDE())
}

// Java method: JSDocInfo.isDeprecated()
func (this_ *JSDocInfo) IsDeprecated() bool {
	return this_.priv_JSDocInfo_checkBit_7975(StaticGet_JSDocInfo_Bit_DEPRECATED())
}

// Java method: JSDocInfo.isInterface()
func (this_ *JSDocInfo) IsInterface() bool {
	return ((this_.propertyBits & (StaticGet_JSDocInfo_Bit_INTERFACE().mask | StaticGet_JSDocInfo_Bit_RECORD().mask)) != int64(int32(0)))
}

// Java method: JSDocInfo.isConstructorOrInterface()
func (this_ *JSDocInfo) IsConstructorOrInterface() bool {
	return ((this_.propertyBits & ((StaticGet_JSDocInfo_Bit_CONSTRUCTOR().mask | StaticGet_JSDocInfo_Bit_INTERFACE().mask) | StaticGet_JSDocInfo_Bit_RECORD().mask)) != int64(int32(0)))
}

// Java method: JSDocInfo.isExport()
func (this_ *JSDocInfo) IsExport() bool {
	return this_.priv_JSDocInfo_checkBit_7975(StaticGet_JSDocInfo_Bit_EXPORT())
}

// Java method: JSDocInfo.isImplicitCast()
func (this_ *JSDocInfo) IsImplicitCast() bool {
	return this_.priv_JSDocInfo_checkBit_7975(StaticGet_JSDocInfo_Bit_IMPLICITCAST())
}

// Java method: JSDocInfo.isNoSideEffects()
func (this_ *JSDocInfo) IsNoSideEffects() bool {
	return this_.priv_JSDocInfo_checkBit_7975(StaticGet_JSDocInfo_Bit_NOSIDEEFFECTS())
}

// Java method: JSDocInfo.isExterns()
func (this_ *JSDocInfo) IsExterns() bool {
	return this_.priv_JSDocInfo_checkBit_7975(StaticGet_JSDocInfo_Bit_EXTERNS())
}

// Java method: JSDocInfo.isNoCoverage()
func (this_ *JSDocInfo) IsNoCoverage() bool {
	return this_.priv_JSDocInfo_checkBit_7975(StaticGet_JSDocInfo_Bit_NOCOVERAGE())
}

// Java method: JSDocInfo.isTypeSummary()
func (this_ *JSDocInfo) IsTypeSummary() bool {
	return this_.priv_JSDocInfo_checkBit_7975(StaticGet_JSDocInfo_Bit_TYPE_SUMMARY())
}

// Java method: JSDocInfo.isNoCompile()
func (this_ *JSDocInfo) IsNoCompile() bool {
	return this_.priv_JSDocInfo_checkBit_7975(StaticGet_JSDocInfo_Bit_NOCOMPILE())
}

// Java method: JSDocInfo.isNoDts()
func (this_ *JSDocInfo) IsNoDts() bool {
	return this_.priv_JSDocInfo_checkBit_7975(StaticGet_JSDocInfo_Bit_NODTS())
}

// Java method: JSDocInfo.isLegacyNamedTsExport()
func (this_ *JSDocInfo) IsLegacyNamedTsExport() bool {
	return this_.priv_JSDocInfo_checkBit_7975(StaticGet_JSDocInfo_Bit_LEGACY_NAMED_TS_EXPORT())
}

// Java method: JSDocInfo.isNoCollapse()
func (this_ *JSDocInfo) IsNoCollapse() bool {
	return this_.priv_JSDocInfo_checkBit_7975(StaticGet_JSDocInfo_Bit_NOCOLLAPSE())
}

// Java method: JSDocInfo.isNoInline()
func (this_ *JSDocInfo) IsNoInline() bool {
	return this_.priv_JSDocInfo_checkBit_7975(StaticGet_JSDocInfo_Bit_NOINLINE())
}

// Java method: JSDocInfo.isRequireInlining()
func (this_ *JSDocInfo) IsRequireInlining() bool {
	return this_.priv_JSDocInfo_checkBit_7975(StaticGet_JSDocInfo_Bit_REQUIRE_INLINING())
}

// Java method: JSDocInfo.isEncourageInlining()
func (this_ *JSDocInfo) IsEncourageInlining() bool {
	return this_.priv_JSDocInfo_checkBit_7975(StaticGet_JSDocInfo_Bit_ENCOURAGE_INLINING())
}

// Java method: JSDocInfo.isCollapsibleOrBreakMyCode()
func (this_ *JSDocInfo) IsCollapsibleOrBreakMyCode() bool {
	return this_.priv_JSDocInfo_checkBit_7975(StaticGet_JSDocInfo_Bit_COLLAPSIBLE_OR_BREAK_MY_CODE())
}

// Java method: JSDocInfo.isPureOrBreakMyCode()
func (this_ *JSDocInfo) IsPureOrBreakMyCode() bool {
	return this_.priv_JSDocInfo_checkBit_7975(StaticGet_JSDocInfo_Bit_PURE_OR_BREAK_MY_CODE())
}

// Java method: JSDocInfo.isProvideGoog()
func (this_ *JSDocInfo) IsProvideGoog() bool {
	return this_.priv_JSDocInfo_checkBit_7975(StaticGet_JSDocInfo_Bit_PROVIDE_GOOG())
}

// Java method: JSDocInfo.isProvideAlreadyProvided()
func (this_ *JSDocInfo) IsProvideAlreadyProvided() bool {
	return this_.priv_JSDocInfo_checkBit_7975(StaticGet_JSDocInfo_Bit_PROVIDE_ALREADY_PROVIDED())
}

// Java method: JSDocInfo.containsDeclarationExcludingTypelessConst()
func (this_ *JSDocInfo) ContainsDeclarationExcludingTypelessConst() bool {
	return (((this_.propertyBits & ((((((((StaticGet_JSDocInfo_Bit_CONSTRUCTOR().mask | StaticGet_JSDocInfo_Bit_DEFINE().mask) | StaticGet_JSDocInfo_Bit_OVERRIDE().mask) | StaticGet_JSDocInfo_Bit_EXPORT().mask) | StaticGet_JSDocInfo_Bit_DEPRECATED().mask) | StaticGet_JSDocInfo_Bit_INTERFACE().mask) | StaticGet_JSDocInfo_Bit_IMPLICITCAST().mask) | StaticGet_JSDocInfo_Bit_NOSIDEEFFECTS().mask) | StaticGet_JSDocInfo_Bit_RECORD().mask)) != int64(int32(0))) || ((this_.propertyKeysBitset & ((((((((StaticGet_JSDocInfo_TYPE().GetJSDocInfo_PropertyBase().mask | StaticGet_JSDocInfo_RETURN_TYPE().GetJSDocInfo_PropertyBase().mask) | StaticGet_JSDocInfo_ENUM_PARAMETER_TYPE().GetJSDocInfo_PropertyBase().mask) | StaticGet_JSDocInfo_TYPEDEF_TYPE().GetJSDocInfo_PropertyBase().mask) | StaticGet_JSDocInfo_THIS_TYPE().GetJSDocInfo_PropertyBase().mask) | StaticGet_JSDocInfo_PARAMETERS().GetJSDocInfo_PropertyBase().mask) | StaticGet_JSDocInfo_IMPLEMENTED_INTERFACES().GetJSDocInfo_PropertyBase().mask) | StaticGet_JSDocInfo_BASE_TYPE().GetJSDocInfo_PropertyBase().mask) | StaticGet_JSDocInfo_VISIBILITY().GetJSDocInfo_PropertyBase().mask)) != int64(int32(0))))
}

// Java method: JSDocInfo.containsTypeDeclaration()
func (this_ *JSDocInfo) ContainsTypeDeclaration() bool {
	return ((this_.propertyKeysBitset & ((((((StaticGet_JSDocInfo_TYPE().GetJSDocInfo_PropertyBase().mask | StaticGet_JSDocInfo_RETURN_TYPE().GetJSDocInfo_PropertyBase().mask) | StaticGet_JSDocInfo_ENUM_PARAMETER_TYPE().GetJSDocInfo_PropertyBase().mask) | StaticGet_JSDocInfo_TYPEDEF_TYPE().GetJSDocInfo_PropertyBase().mask) | StaticGet_JSDocInfo_THIS_TYPE().GetJSDocInfo_PropertyBase().mask) | StaticGet_JSDocInfo_PARAMETERS().GetJSDocInfo_PropertyBase().mask) | StaticGet_JSDocInfo_BASE_TYPE().GetJSDocInfo_PropertyBase().mask)) != int64(int32(0)))
}

// Java method: JSDocInfo.containsDeclaration()
func (this_ *JSDocInfo) ContainsDeclaration() bool {
	return (this_.ContainsDeclarationExcludingTypelessConst() || this_.priv_JSDocInfo_checkBit_7975(StaticGet_JSDocInfo_Bit_CONST()))
}

// Java method: JSDocInfo.containsFunctionDeclaration()
func (this_ *JSDocInfo) ContainsFunctionDeclaration() bool {
	if this_.priv_JSDocInfo_checkBit_7975(StaticGet_JSDocInfo_Bit_CONSTRUCTOR()) || ((this_.propertyKeysBitset & ((StaticGet_JSDocInfo_RETURN_TYPE().GetJSDocInfo_PropertyBase().mask | StaticGet_JSDocInfo_THIS_TYPE().GetJSDocInfo_PropertyBase().mask) | StaticGet_JSDocInfo_PARAMETERS().GetJSDocInfo_PropertyBase().mask)) != int64(int32(0))) {
		return true
	}
	var type_ JSTypeExpression = j2gort.Cast[JSTypeExpression](StaticGet_JSDocInfo_TYPE().Get(this_))
	_ = type_
	if type_ != (JSTypeExpression{}) {
		return type_.GetRoot().IsFunction()
	}
	return this_.priv_JSDocInfo_checkBit_7975(StaticGet_JSDocInfo_Bit_NOSIDEEFFECTS())
}

// Java method: JSDocInfo.containsTypeDefinition()
func (this_ *JSDocInfo) ContainsTypeDefinition() bool {
	return (((this_.propertyBits & (StaticGet_JSDocInfo_Bit_CONSTRUCTOR().mask | StaticGet_JSDocInfo_Bit_INTERFACE().mask)) != int64(int32(0))) || ((this_.propertyKeysBitset & (StaticGet_JSDocInfo_ENUM_PARAMETER_TYPE().GetJSDocInfo_PropertyBase().mask | StaticGet_JSDocInfo_TYPEDEF_TYPE().GetJSDocInfo_PropertyBase().mask)) != int64(int32(0))))
}

// Java method: JSDocInfo.isAtSignCodePresent()
func (this_ *JSDocInfo) IsAtSignCodePresent() bool {
	var entireComment string = this_.GetOriginalCommentString()
	_ = entireComment
	if entireComment == j2gort.NullString {
		return false
	} else {
		return (j2gort.JStrIndexOf(entireComment, "@code") >= 0)
	}
}

// Java method: JSDocInfo.getVisibility()
func (this_ *JSDocInfo) GetVisibility() JSDocInfo_Visibility {
	var visibility JSDocInfo_Visibility = j2gort.Cast[JSDocInfo_Visibility](StaticGet_JSDocInfo_VISIBILITY().Get(this_))
	_ = visibility
	if visibility != 0 {
		return visibility
	} else {
		return StaticGet_JSDocInfo_Visibility_INHERITED()
	}
}

// Java method: JSDocInfo.getParameterType(java.lang.String parameter)
func (this_ *JSDocInfo) GetParameterType(parameter string) JSTypeExpression {
	_ = parameter
	var params *JSDocInfo_TypeMap = j2gort.CastPtr[JSDocInfo_TypeMap](StaticGet_JSDocInfo_PARAMETERS().Get(this_))
	_ = params
	if params != nil {
		return params.Get(parameter)
	} else {
		return (JSTypeExpression{})
	}
}

// Java method: JSDocInfo.hasParameter(java.lang.String parameter)
func (this_ *JSDocInfo) HasParameter(parameter string) bool {
	_ = parameter
	var params *JSDocInfo_TypeMap = j2gort.CastPtr[JSDocInfo_TypeMap](StaticGet_JSDocInfo_PARAMETERS().Get(this_))
	_ = params
	return ((params != nil) && params.ContainsKey(parameter))
}

// Java method: JSDocInfo.hasParameterType(java.lang.String parameter)
func (this_ *JSDocInfo) HasParameterType(parameter string) bool {
	_ = parameter
	return (this_.GetParameterType(parameter) != (JSTypeExpression{}))
}

// Java method: JSDocInfo.getParameterNames()
func (this_ *JSDocInfo) GetParameterNames() j2gort.Ext_java_util_Set {
	var params *JSDocInfo_TypeMap = j2gort.CastPtr[JSDocInfo_TypeMap](StaticGet_JSDocInfo_PARAMETERS().Get(this_))
	_ = params
	if params != nil {
		return j2gort.AsIface[j2gort.Ext_java_util_Set](params.KeysAsImmutableSet())
	} else {
		return j2gort.AsIface[j2gort.Ext_java_util_Set](j2gort.Ext_com_google_common_collect_ImmutableSet_Of_0())
	}
}

// Java method: JSDocInfo.getParameterNameAt(int index)
func (this_ *JSDocInfo) GetParameterNameAt(index int32) string {
	_ = index
	var params *JSDocInfo_TypeMap = j2gort.CastPtr[JSDocInfo_TypeMap](StaticGet_JSDocInfo_PARAMETERS().Get(this_))
	_ = params
	if (params == nil) || (index >= params.Size()) {
		return j2gort.NullString
	}
	return params.KeyAt(index)
}

// Java method: JSDocInfo.getParameterCount()
func (this_ *JSDocInfo) GetParameterCount() int32 {
	var params *JSDocInfo_TypeMap = j2gort.CastPtr[JSDocInfo_TypeMap](StaticGet_JSDocInfo_PARAMETERS().Get(this_))
	_ = params
	if params != nil {
		return params.Size()
	} else {
		return int32(0)
	}
}

// Java method: JSDocInfo.getThrowsAnnotations()
func (this_ *JSDocInfo) GetThrowsAnnotations() j2gort.Ext_java_util_List {
	var annotations j2gort.Ext_java_util_List = j2gort.AsIface[j2gort.Ext_java_util_List](StaticGet_JSDocInfo_THROWS_ANNOTATIONS().Get(this_))
	_ = annotations
	if annotations != nil {
		return annotations
	} else {
		return j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.Ext_com_google_common_collect_ImmutableList_Of_0())
	}
}

// Java method: JSDocInfo.hasEnumParameterType()
func (this_ *JSDocInfo) HasEnumParameterType() bool {
	return ((this_.propertyKeysBitset & StaticGet_JSDocInfo_ENUM_PARAMETER_TYPE().GetJSDocInfo_PropertyBase().mask) != int64(0))
}

// Java method: JSDocInfo.hasTypedefType()
func (this_ *JSDocInfo) HasTypedefType() bool {
	return ((this_.propertyKeysBitset & StaticGet_JSDocInfo_TYPEDEF_TYPE().GetJSDocInfo_PropertyBase().mask) != int64(0))
}

// Java method: JSDocInfo.hasReturnType()
func (this_ *JSDocInfo) HasReturnType() bool {
	return ((this_.propertyKeysBitset & StaticGet_JSDocInfo_RETURN_TYPE().GetJSDocInfo_PropertyBase().mask) != int64(0))
}

// Java method: JSDocInfo.hasType()
func (this_ *JSDocInfo) HasType() bool {
	return ((this_.propertyKeysBitset & StaticGet_JSDocInfo_TYPE().GetJSDocInfo_PropertyBase().mask) != int64(0))
}

// Java method: JSDocInfo.hasTypeInformation()
func (this_ *JSDocInfo) HasTypeInformation() bool {
	return ((this_.propertyKeysBitset & (((StaticGet_JSDocInfo_TYPE().GetJSDocInfo_PropertyBase().mask | StaticGet_JSDocInfo_TYPEDEF_TYPE().GetJSDocInfo_PropertyBase().mask) | StaticGet_JSDocInfo_ENUM_PARAMETER_TYPE().GetJSDocInfo_PropertyBase().mask) | StaticGet_JSDocInfo_RETURN_TYPE().GetJSDocInfo_PropertyBase().mask)) != int64(int32(0)))
}

// Java method: JSDocInfo.isInlineType()
func (this_ *JSDocInfo) IsInlineType() bool {
	return this_.priv_JSDocInfo_checkBit_7975(StaticGet_JSDocInfo_Bit_INLINE_TYPE())
}

// Java method: JSDocInfo.getReturnType()
func (this_ *JSDocInfo) GetReturnType() JSTypeExpression {
	return StaticGet_JSDocInfo_RETURN_TYPE().GetTyped(this_)
}

// Java method: JSDocInfo.getEnumParameterType()
func (this_ *JSDocInfo) GetEnumParameterType() JSTypeExpression {
	return StaticGet_JSDocInfo_ENUM_PARAMETER_TYPE().GetTyped(this_)
}

// Java method: JSDocInfo.getTypedefType()
func (this_ *JSDocInfo) GetTypedefType() JSTypeExpression {
	return StaticGet_JSDocInfo_TYPEDEF_TYPE().GetTyped(this_)
}

// Java method: JSDocInfo.getType()
func (this_ *JSDocInfo) GetType() JSTypeExpression {
	return StaticGet_JSDocInfo_TYPE().GetTyped(this_)
}

// Java method: JSDocInfo.getThisType()
func (this_ *JSDocInfo) GetThisType() JSTypeExpression {
	return StaticGet_JSDocInfo_THIS_TYPE().GetTyped(this_)
}

// Java method: JSDocInfo.hasThisType()
func (this_ *JSDocInfo) HasThisType() bool {
	return ((this_.propertyKeysBitset & StaticGet_JSDocInfo_THIS_TYPE().GetJSDocInfo_PropertyBase().mask) != int64(0))
}

// Java method: JSDocInfo.getBaseType()
func (this_ *JSDocInfo) GetBaseType() JSTypeExpression {
	return StaticGet_JSDocInfo_BASE_TYPE().GetTyped(this_)
}

// Java method: JSDocInfo.getDescription()
func (this_ *JSDocInfo) GetDescription() string {
	return j2gort.CastString(StaticGet_JSDocInfo_DESCRIPTION().Get(this_))
}

// Java method: JSDocInfo.getTsTypes()
func (this_ *JSDocInfo) GetTsTypes() j2gort.Ext_com_google_common_collect_ImmutableList {
	var types j2gort.Ext_java_util_List = j2gort.AsIface[j2gort.Ext_java_util_List](StaticGet_JSDocInfo_TS_TYPES().Get(this_))
	_ = types
	if types != nil {
		return j2gort.Ext_com_google_common_collect_ImmutableList_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](types))
	} else {
		return j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()
	}
}

// Java method: JSDocInfo.getMeaning()
func (this_ *JSDocInfo) GetMeaning() string {
	return j2gort.CastString(StaticGet_JSDocInfo_MEANING().Get(this_))
}

// Java method: JSDocInfo.getAlternateMessageId()
func (this_ *JSDocInfo) GetAlternateMessageId() string {
	return j2gort.CastString(StaticGet_JSDocInfo_ALTERNATE_MESSAGE_ID().Get(this_))
}

// Java method: JSDocInfo.getLendsName()
func (this_ *JSDocInfo) GetLendsName() JSTypeExpression {
	return StaticGet_JSDocInfo_LENDS_NAME().GetTyped(this_)
}

// Java method: JSDocInfo.hasLendsName()
func (this_ *JSDocInfo) HasLendsName() bool {
	return ((this_.propertyKeysBitset & StaticGet_JSDocInfo_LENDS_NAME().GetJSDocInfo_PropertyBase().mask) != int64(0))
}

// Java method: JSDocInfo.getClosurePrimitiveId()
func (this_ *JSDocInfo) GetClosurePrimitiveId() string {
	return j2gort.CastString(StaticGet_JSDocInfo_CLOSURE_PRIMITIVE_ID().Get(this_))
}

// Java method: JSDocInfo.hasClosurePrimitiveId()
func (this_ *JSDocInfo) HasClosurePrimitiveId() bool {
	return (j2gort.CastString(StaticGet_JSDocInfo_CLOSURE_PRIMITIVE_ID().Get(this_)) != j2gort.NullString)
}

// Java method: JSDocInfo.isNgInject()
func (this_ *JSDocInfo) IsNgInject() bool {
	return this_.priv_JSDocInfo_checkBit_7975(StaticGet_JSDocInfo_Bit_NG_INJECT())
}

// Java method: JSDocInfo.isWizaction()
func (this_ *JSDocInfo) IsWizaction() bool {
	return this_.priv_JSDocInfo_checkBit_7975(StaticGet_JSDocInfo_Bit_WIZ_ACTION())
}

// Java method: JSDocInfo.isWizcallback()
func (this_ *JSDocInfo) IsWizcallback() bool {
	return this_.priv_JSDocInfo_checkBit_7975(StaticGet_JSDocInfo_Bit_WIZ_CALLBACK())
}

// Java method: JSDocInfo.isPolymerBehavior()
func (this_ *JSDocInfo) IsPolymerBehavior() bool {
	return this_.priv_JSDocInfo_checkBit_7975(StaticGet_JSDocInfo_Bit_POLYMER_BEHAVIOR())
}

// Java method: JSDocInfo.isPolymer()
func (this_ *JSDocInfo) IsPolymer() bool {
	return this_.priv_JSDocInfo_checkBit_7975(StaticGet_JSDocInfo_Bit_POLYMER())
}

// Java method: JSDocInfo.isCustomElement()
func (this_ *JSDocInfo) IsCustomElement() bool {
	return this_.priv_JSDocInfo_checkBit_7975(StaticGet_JSDocInfo_Bit_CUSTOM_ELEMENT())
}

// Java method: JSDocInfo.isMixinClass()
func (this_ *JSDocInfo) IsMixinClass() bool {
	return this_.priv_JSDocInfo_checkBit_7975(StaticGet_JSDocInfo_Bit_MIXIN_CLASS())
}

// Java method: JSDocInfo.isMixinFunction()
func (this_ *JSDocInfo) IsMixinFunction() bool {
	return this_.priv_JSDocInfo_checkBit_7975(StaticGet_JSDocInfo_Bit_MIXIN_FUNCTION())
}

// Java method: JSDocInfo.isSassGeneratedCssTs()
func (this_ *JSDocInfo) IsSassGeneratedCssTs() bool {
	return this_.priv_JSDocInfo_checkBit_7975(StaticGet_JSDocInfo_Bit_SASS_GENERATED_CSS_TS())
}

// Java method: JSDocInfo.isClosureUnawareCode()
func (this_ *JSDocInfo) IsClosureUnawareCode() bool {
	return (j2gort.Cast[JSDocInfo_PerFileClosureUnawareMode](StaticGet_JSDocInfo_CLOSURE_UNAWARE_CONFIG_VALUE().Get(this_)) != 0)
}

// Java method: JSDocInfo.getPerFileClosureUnawareMode()
func (this_ *JSDocInfo) GetPerFileClosureUnawareMode() JSDocInfo_PerFileClosureUnawareMode {
	return j2gort.Cast[JSDocInfo_PerFileClosureUnawareMode](StaticGet_JSDocInfo_CLOSURE_UNAWARE_CONFIG_VALUE().Get(this_))
}

// Java method: JSDocInfo.isUsedViaDotConstructor()
func (this_ *JSDocInfo) IsUsedViaDotConstructor() bool {
	return this_.priv_JSDocInfo_checkBit_7975(StaticGet_JSDocInfo_Bit_USED_VIA_DOT_CONSTRUCTOR())
}

// Java method: JSDocInfo.getLicense()
func (this_ *JSDocInfo) GetLicense() string {
	return j2gort.CastString(StaticGet_JSDocInfo_LICENSE().Get(this_))
}

// Java method: JSDocInfo.toString()
func (this_ *JSDocInfo) String() string {
	return this_.ToStringVerbose()
}

// Java method: JSDocInfo.toStringVerbose()
func (this_ *JSDocInfo) ToStringVerbose() string {
	var helper *j2gort.Ext_com_google_common_base_MoreObjects_ToStringHelper = j2gort.Ext_com_google_common_base_MoreObjects_S_ToStringHelper_1_any(this_).Add_2_string_any("bitset", func() any {
		if this_.propertyBits == int64(int32(0)) {
			return nil
		}
		return javalang.LongToHexString(this_.propertyBits)
	}())
	_ = helper
	for _, el_20 := range j2gort.CheckNotNullSlice(JSDocInfo_Bit_Values()) {
		b := el_20
		_ = b
		if this_.priv_JSDocInfo_checkBit_7975(b) {
			helper.Add_2_string_bool(("bit:" + j2gort.StrValString(b.name)), true)
		}
	}
	var bits int64 = this_.propertyKeysBitset
	_ = bits
	for bits != int64(0) {
		var low int32 = j2gort.BoxLong_NumberOfTrailingZeros(bits)
		_ = low
		bits = (bits & (bits - int64(1)))
		helper = helper.Add_2_string_any(StaticGet_JSDocInfo_Property_Values()[low].GetJSDocInfo_PropertyBase().name, this_.priv_JSDocInfo_getPropertyValueForBit_8109(low))
	}
	return helper.OmitNullValues().String()
}

// Java method: JSDocInfo.hasBaseType()
func (this_ *JSDocInfo) HasBaseType() bool {
	return ((this_.propertyKeysBitset & StaticGet_JSDocInfo_BASE_TYPE().GetJSDocInfo_PropertyBase().mask) != int64(0))
}

// Java method: JSDocInfo.getImplementedInterfaces()
func (this_ *JSDocInfo) GetImplementedInterfaces() j2gort.Ext_java_util_List {
	return StaticGet_JSDocInfo_IMPLEMENTED_INTERFACES().GetUnmodifiable(this_)
}

// Java method: JSDocInfo.getImplementedInterfaceCount()
func (this_ *JSDocInfo) GetImplementedInterfaceCount() int32 {
	var list *JSDocInfo_TypeList = j2gort.CastPtr[JSDocInfo_TypeList](StaticGet_JSDocInfo_IMPLEMENTED_INTERFACES().Get(this_))
	_ = list
	if list != nil {
		return list.Size()
	} else {
		return int32(0)
	}
}

// Java method: JSDocInfo.getImplementedInterfaceAt(int index)
func (this_ *JSDocInfo) GetImplementedInterfaceAt(index int32) JSTypeExpression {
	_ = index
	return j2gort.CastPtr[JSDocInfo_TypeList](StaticGet_JSDocInfo_IMPLEMENTED_INTERFACES().Get(this_)).Get(index)
}

// Java method: JSDocInfo.getExtendedInterfaces()
func (this_ *JSDocInfo) GetExtendedInterfaces() j2gort.Ext_java_util_List {
	return StaticGet_JSDocInfo_EXTENDED_INTERFACES().GetUnmodifiable(this_)
}

// Java method: JSDocInfo.getExtendedInterfacesCount()
func (this_ *JSDocInfo) GetExtendedInterfacesCount() int32 {
	var list *JSDocInfo_TypeList = j2gort.CastPtr[JSDocInfo_TypeList](StaticGet_JSDocInfo_EXTENDED_INTERFACES().Get(this_))
	_ = list
	if list != nil {
		return list.Size()
	} else {
		return int32(0)
	}
}

// Java method: JSDocInfo.getExtendedInterfaceAt(int index)
func (this_ *JSDocInfo) GetExtendedInterfaceAt(index int32) JSTypeExpression {
	_ = index
	return j2gort.CastPtr[JSDocInfo_TypeList](StaticGet_JSDocInfo_EXTENDED_INTERFACES().Get(this_)).Get(index)
}

// Java method: JSDocInfo.getDeprecationReason()
func (this_ *JSDocInfo) GetDeprecationReason() string {
	return j2gort.CastString(StaticGet_JSDocInfo_DEPRECATION_REASON().Get(this_))
}

// Java method: JSDocInfo.getSuppressions()
func (this_ *JSDocInfo) GetSuppressions() j2gort.Ext_java_util_Set {
	var suppressions j2gort.Ext_com_google_common_collect_ImmutableMap = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap](StaticGet_JSDocInfo_SUPPRESSIONS().Get(this_))
	_ = suppressions
	if (suppressions == nil) || suppressions.IsEmpty() {
		return j2gort.AsIface[j2gort.Ext_java_util_Set](j2gort.Ext_com_google_common_collect_ImmutableSet_Of_0())
	}
	if suppressions.Size() == int32(1) {
		return j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](j2gort.Ext_com_google_common_collect_Iterables_GetOnlyElement_1(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](suppressions.KeySet_0())))
	}
	var names j2gort.Ext_com_google_common_collect_ImmutableSet_Builder = j2gort.Ext_com_google_common_collect_ImmutableSet_S_Builder()
	_ = names
	{
		cur_21 := j2gort.NewMapKeyCursor(suppressions)
		for cur_21.HasNext() {
			rawElem_22 := cur_21.Next()
			nameSet := j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](rawElem_22)
			_ = nameSet
			names.AddAll_1_Ext_java_lang_Iterable(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](nameSet))
		}
	}
	return j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](names.Build())
}

// Java method: JSDocInfo.getSuppressionsAndTheirDescription()
func (this_ *JSDocInfo) GetSuppressionsAndTheirDescription() j2gort.Ext_com_google_common_collect_ImmutableMap {
	var suppressions j2gort.Ext_com_google_common_collect_ImmutableMap = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap](StaticGet_JSDocInfo_SUPPRESSIONS().Get(this_))
	_ = suppressions
	if suppressions != nil {
		return suppressions
	} else {
		return j2gort.Ext_com_google_common_collect_ImmutableMap_Of_0()
	}
}

// Java method: JSDocInfo.getModifies()
func (this_ *JSDocInfo) GetModifies() j2gort.Ext_java_util_Set {
	var modifies j2gort.Ext_com_google_common_collect_ImmutableSet = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](StaticGet_JSDocInfo_MODIFIES().Get(this_))
	_ = modifies
	if modifies != nil {
		return j2gort.AsIface[j2gort.Ext_java_util_Set](modifies)
	} else {
		return j2gort.AsIface[j2gort.Ext_java_util_Set](j2gort.Ext_com_google_common_collect_ImmutableSet_Of_0())
	}
}

// Java method: JSDocInfo.hasDescriptionForParameter(java.lang.String name)
func (this_ *JSDocInfo) HasDescriptionForParameter(name string) bool {
	_ = name
	var params *JSDocInfo_StringMap = j2gort.CastPtr[JSDocInfo_StringMap](StaticGet_JSDocInfo_PARAMETER_DESCRIPTIONS().Get(this_))
	_ = params
	return ((params != nil) && params.ContainsKey(name))
}

// Java method: JSDocInfo.getDescriptionForParameter(java.lang.String name)
func (this_ *JSDocInfo) GetDescriptionForParameter(name string) string {
	_ = name
	var params *JSDocInfo_StringMap = j2gort.CastPtr[JSDocInfo_StringMap](StaticGet_JSDocInfo_PARAMETER_DESCRIPTIONS().Get(this_))
	_ = params
	if params != nil {
		return params.Get(name)
	} else {
		return j2gort.NullString
	}
}

// Java method: JSDocInfo.getAuthors()
func (this_ *JSDocInfo) GetAuthors() j2gort.Ext_java_util_List {
	return j2gort.AsIface[j2gort.Ext_java_util_List](StaticGet_JSDocInfo_AUTHORS().Get(this_))
}

// Java method: JSDocInfo.getReferences()
func (this_ *JSDocInfo) GetReferences() j2gort.Ext_java_util_List {
	return j2gort.AsIface[j2gort.Ext_java_util_List](StaticGet_JSDocInfo_SEES().Get(this_))
}

// Java method: JSDocInfo.getReturnDescription()
func (this_ *JSDocInfo) GetReturnDescription() string {
	return j2gort.CastString(StaticGet_JSDocInfo_RETURN_DESCRIPTION().Get(this_))
}

// Java method: JSDocInfo.getBlockDescription()
func (this_ *JSDocInfo) GetBlockDescription() string {
	return j2gort.CastString(StaticGet_JSDocInfo_BLOCK_DESCRIPTION().Get(this_))
}

// Java method: JSDocInfo.hasFileOverview()
func (this_ *JSDocInfo) HasFileOverview() bool {
	return this_.priv_JSDocInfo_checkBit_7975(StaticGet_JSDocInfo_Bit_FILEOVERVIEW())
}

// Java method: JSDocInfo.getFileOverview()
func (this_ *JSDocInfo) GetFileOverview() string {
	return j2gort.CastString(StaticGet_JSDocInfo_FILEOVERVIEW_DESCRIPTION().Get(this_))
}

// Java method: JSDocInfo.hasEnhance()
func (this_ *JSDocInfo) HasEnhance() bool {
	return this_.priv_JSDocInfo_checkBit_7975(StaticGet_JSDocInfo_Bit_ENHANCED_NAMESPACE())
}

// Java method: JSDocInfo.getEnhance()
func (this_ *JSDocInfo) GetEnhance() string {
	return j2gort.CastString(StaticGet_JSDocInfo_ENHANCED_NAMESPACE().Get(this_))
}

// Java method: JSDocInfo.hasMods()
func (this_ *JSDocInfo) HasMods() bool {
	return (j2gort.CastString(StaticGet_JSDocInfo_MODS().Get(this_)) != j2gort.NullString)
}

// Java method: JSDocInfo.getMods()
func (this_ *JSDocInfo) GetMods() string {
	return j2gort.CastString(StaticGet_JSDocInfo_MODS().Get(this_))
}

// Java method: JSDocInfo.getMarkers()
func (this_ *JSDocInfo) GetMarkers() j2gort.Ext_java_util_Collection {
	var markers j2gort.Ext_java_util_ArrayList = j2gort.AsIface[j2gort.Ext_java_util_ArrayList](StaticGet_JSDocInfo_MARKERS().Get(this_))
	_ = markers
	if markers != nil {
		return j2gort.AsIface[j2gort.Ext_java_util_Collection](markers)
	} else {
		return j2gort.AsIface[j2gort.Ext_java_util_Collection](j2gort.Ext_com_google_common_collect_ImmutableList_Of_0())
	}
}

// Java method: JSDocInfo.getTemplateTypeNames()
func (this_ *JSDocInfo) GetTemplateTypeNames() j2gort.Ext_com_google_common_collect_ImmutableList {
	var map_ *JSDocInfo_TypeMap = j2gort.CastPtr[JSDocInfo_TypeMap](StaticGet_JSDocInfo_TEMPLATE_TYPE_NAMES().Get(this_))
	_ = map_
	if map_ != nil {
		return map_.KeysAsImmutableList()
	} else {
		return j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()
	}
}

// Java method: JSDocInfo.getTemplateTypes()
func (this_ *JSDocInfo) GetTemplateTypes() j2gort.Ext_com_google_common_collect_ImmutableMap {
	var map_ *JSDocInfo_TypeMap = j2gort.CastPtr[JSDocInfo_TypeMap](StaticGet_JSDocInfo_TEMPLATE_TYPE_NAMES().Get(this_))
	_ = map_
	if map_ != nil {
		return map_.ToImmutableMap()
	} else {
		return j2gort.Ext_com_google_common_collect_ImmutableMap_Of_0()
	}
}

// Java method: JSDocInfo.getTypeTransformations()
func (this_ *JSDocInfo) GetTypeTransformations() j2gort.Ext_com_google_common_collect_ImmutableMap {
	var map_ *JSDocInfo_NodeMap = j2gort.CastPtr[JSDocInfo_NodeMap](StaticGet_JSDocInfo_TYPE_TRANSFORMATIONS().Get(this_))
	_ = map_
	if map_ != nil {
		return map_.ToImmutableMap()
	} else {
		return j2gort.Ext_com_google_common_collect_ImmutableMap_Of_0()
	}
}

// Java method: JSDocInfo.hasTypeNodes()
func (this_ *JSDocInfo) HasTypeNodes() bool {
	var typeBits int64 = (this_.propertyKeysBitset & StaticGet_JSDocInfo_TYPE_EXPR_PROPERTIES_MASK())
	_ = typeBits
	if typeBits == int64(0) {
		return false
	}
	if (typeBits & StaticGet_JSDocInfo_SINGLE_TYPE_EXPR_PROPERTIES_MASK()) != int64(0) {
		return true
	}
	for typeBits != int64(0) {
		var low int32 = j2gort.BoxLong_NumberOfTrailingZeros(typeBits)
		_ = low
		typeBits = (typeBits & (typeBits - int64(1)))
		var prop JSDocInfo_Property = StaticGet_JSDocInfo_Property_Values()[low]
		_ = prop
		if prop.HasAnyTypeExpression(this_.priv_JSDocInfo_getNonType0ValueForBit_8108(low)) {
			return true
		}
	}
	return false
}

// Java method: JSDocInfo.getTypeExpressions()
func (this_ *JSDocInfo) GetTypeExpressions() j2gort.Ext_java_util_Collection {
	var typeBits int64 = (this_.propertyKeysBitset & StaticGet_JSDocInfo_TYPE_EXPR_PROPERTIES_MASK())
	_ = typeBits
	if typeBits == int64(0) {
		return j2gort.AsIface[j2gort.Ext_java_util_Collection](j2gort.Ext_com_google_common_collect_ImmutableList_Of_0())
	}
	if ((typeBits & (typeBits - int64(1))) == int64(0)) && ((typeBits & StaticGet_JSDocInfo_SINGLE_TYPE_EXPR_PROPERTIES_MASK()) != int64(0)) {
		if this_.type0 != (JSTypeExpression{}) {
			return j2gort.AsIface[j2gort.Ext_java_util_Collection](j2gort.Ext_com_google_common_collect_ImmutableList_Of_1(j2gort.BoxNullableVal(this_.type0)))
		} else {
			return j2gort.AsIface[j2gort.Ext_java_util_Collection](j2gort.Ext_com_google_common_collect_ImmutableList_Of_0())
		}
	}
	var builder *j2gort.Ext_com_google_common_collect_ImmutableList_Builder = nil
	_ = builder
	for typeBits != int64(0) {
		var low int32 = j2gort.BoxLong_NumberOfTrailingZeros(typeBits)
		_ = low
		typeBits = (typeBits & (typeBits - int64(1)))
		if ((javalang.Shl64(int64(1), int(low))) & StaticGet_JSDocInfo_SINGLE_TYPE_EXPR_PROPERTIES_MASK()) != int64(0) {
			var expr JSTypeExpression = this_.priv_JSDocInfo_getSingleTypeForBit_8107(low)
			_ = expr
			if expr != (JSTypeExpression{}) {
				if builder == nil {
					builder = (&j2gort.Ext_com_google_common_collect_ImmutableList_Builder{})
				}
				builder.AddVoid_(j2gort.BoxNullableVal(expr))
			}
		} else {
			var prop JSDocInfo_Property = StaticGet_JSDocInfo_Property_Values()[low]
			_ = prop
			var value any = this_.priv_JSDocInfo_getNonType0ValueForBit_8108(low)
			_ = value
			if prop.HasAnyTypeExpression(value) {
				if builder == nil {
					builder = (&j2gort.Ext_com_google_common_collect_ImmutableList_Builder{})
				}
				prop.CollectTypeExpressions(value, builder)
			}
		}
	}
	if builder != nil {
		return j2gort.AsIface[j2gort.Ext_java_util_Collection](builder.BuildImmList_())
	} else {
		return j2gort.AsIface[j2gort.Ext_java_util_Collection](j2gort.Ext_com_google_common_collect_ImmutableList_Of_0())
	}
}

// Java method: JSDocInfo.getTypeNodeCount()
func (this_ *JSDocInfo) GetTypeNodeCount() int32 {
	var typeBits int64 = (this_.propertyKeysBitset & StaticGet_JSDocInfo_TYPE_EXPR_PROPERTIES_MASK())
	_ = typeBits
	if typeBits == int64(0) {
		return int32(0)
	}
	if ((typeBits & (typeBits - int64(1))) == int64(0)) && ((typeBits & StaticGet_JSDocInfo_SINGLE_TYPE_EXPR_PROPERTIES_MASK()) != int64(0)) {
		if this_.type0 != (JSTypeExpression{}) {
			return int32(1)
		} else {
			return int32(0)
		}
	}
	var count int32 = int32(0)
	_ = count
	for typeBits != int64(0) {
		var low int32 = j2gort.BoxLong_NumberOfTrailingZeros(typeBits)
		_ = low
		typeBits = (typeBits & (typeBits - int64(1)))
		if ((javalang.Shl64(int64(1), int(low))) & StaticGet_JSDocInfo_SINGLE_TYPE_EXPR_PROPERTIES_MASK()) != int64(0) {
			if this_.priv_JSDocInfo_getSingleTypeForBit_8107(low) != (JSTypeExpression{}) {
				count = count + 1
			}
		} else {
			var prop JSDocInfo_Property = StaticGet_JSDocInfo_Property_Values()[low]
			_ = prop
			count = (count + prop.CountTypeNodes(this_.priv_JSDocInfo_getNonType0ValueForBit_8108(low)))
		}
	}
	return count
}

// Java method: JSDocInfo.getTypeNodeAt(int index)
func (this_ *JSDocInfo) GetTypeNodeAt(index int32) Node {
	_ = index
	if index < int32(0) {
		panic(j2gort.AsIface[javalang.IIndexOutOfBoundsException](j2gort.ConstructShimThrowable("java.lang.IndexOutOfBoundsException", j2gort.ValueOfInt(index))))
	}
	var typeBits int64 = (this_.propertyKeysBitset & StaticGet_JSDocInfo_TYPE_EXPR_PROPERTIES_MASK())
	_ = typeBits
	if ((typeBits & (typeBits - int64(1))) == int64(0)) && ((typeBits & StaticGet_JSDocInfo_SINGLE_TYPE_EXPR_PROPERTIES_MASK()) != int64(0)) {
		if (index == int32(0)) && (this_.type0 != (JSTypeExpression{})) {
			return this_.type0.GetRoot()
		}
		panic(j2gort.AsIface[javalang.IIndexOutOfBoundsException](j2gort.ConstructShimThrowable("java.lang.IndexOutOfBoundsException", j2gort.ValueOfInt(index))))
	}
	var remaining int32 = index
	_ = remaining
	for typeBits != int64(0) {
		var low int32 = j2gort.BoxLong_NumberOfTrailingZeros(typeBits)
		_ = low
		typeBits = (typeBits & (typeBits - int64(1)))
		if ((javalang.Shl64(int64(1), int(low))) & StaticGet_JSDocInfo_SINGLE_TYPE_EXPR_PROPERTIES_MASK()) != int64(0) {
			var expr JSTypeExpression = this_.priv_JSDocInfo_getSingleTypeForBit_8107(low)
			_ = expr
			if expr != (JSTypeExpression{}) {
				if remaining == int32(0) {
					return expr.GetRoot()
				}
				remaining = remaining - 1
			}
		} else {
			var prop JSDocInfo_Property = StaticGet_JSDocInfo_Property_Values()[low]
			_ = prop
			var val any = this_.priv_JSDocInfo_getNonType0ValueForBit_8108(low)
			_ = val
			var subCount int32 = prop.CountTypeNodes(val)
			_ = subCount
			if remaining < subCount {
				return prop.TypeNodeAt(val, remaining)
			}
			remaining = (remaining - subCount)
		}
	}
	panic(j2gort.AsIface[javalang.IIndexOutOfBoundsException](j2gort.ConstructShimThrowable("java.lang.IndexOutOfBoundsException", j2gort.ValueOfInt(index))))
}

// Java method: JSDocInfo.getTypeNodes()
func (this_ *JSDocInfo) GetTypeNodes() j2gort.Ext_java_util_Collection {
	var typeBits int64 = (this_.propertyKeysBitset & StaticGet_JSDocInfo_TYPE_EXPR_PROPERTIES_MASK())
	_ = typeBits
	if typeBits == int64(0) {
		return j2gort.AsIface[j2gort.Ext_java_util_Collection](j2gort.Ext_com_google_common_collect_ImmutableList_Of_0())
	}
	if ((typeBits & (typeBits - int64(1))) == int64(0)) && ((typeBits & StaticGet_JSDocInfo_SINGLE_TYPE_EXPR_PROPERTIES_MASK()) != int64(0)) {
		if this_.type0 != (JSTypeExpression{}) {
			return j2gort.AsIface[j2gort.Ext_java_util_Collection](j2gort.Ext_com_google_common_collect_ImmutableList_Of_1(j2gort.BoxNode(this_.type0.GetRoot())))
		} else {
			return j2gort.AsIface[j2gort.Ext_java_util_Collection](j2gort.Ext_com_google_common_collect_ImmutableList_Of_0())
		}
	}
	var nodes j2gort.Ext_java_util_ArrayList = nil
	_ = nodes
	for typeBits != int64(0) {
		var low int32 = j2gort.BoxLong_NumberOfTrailingZeros(typeBits)
		_ = low
		typeBits = (typeBits & (typeBits - int64(1)))
		if ((javalang.Shl64(int64(1), int(low))) & StaticGet_JSDocInfo_SINGLE_TYPE_EXPR_PROPERTIES_MASK()) != int64(0) {
			var expr JSTypeExpression = this_.priv_JSDocInfo_getSingleTypeForBit_8107(low)
			_ = expr
			if expr != (JSTypeExpression{}) {
				if nodes == nil {
					nodes = j2gort.NewExt_java_util_ArrayList_1_int32(int32(4))
				}
				nodes.Add_1(j2gort.BoxNode(expr.GetRoot()))
			}
		} else {
			var prop JSDocInfo_Property = StaticGet_JSDocInfo_Property_Values()[low]
			_ = prop
			nodes = prop.CollectTypeNodes(this_.priv_JSDocInfo_getNonType0ValueForBit_8108(low), nodes)
		}
	}
	if nodes != nil {
		return j2gort.AsIface[j2gort.Ext_java_util_Collection](nodes)
	} else {
		return j2gort.AsIface[j2gort.Ext_java_util_Collection](j2gort.Ext_com_google_common_collect_ImmutableList_Of_0())
	}
}

// Java method: JSDocInfo.hasModifies()
func (this_ *JSDocInfo) HasModifies() bool {
	return (j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](StaticGet_JSDocInfo_MODIFIES().Get(this_)) != nil)
}

// Java method: JSDocInfo.getOriginalCommentString()
func (this_ *JSDocInfo) GetOriginalCommentString() string {
	return j2gort.CastString(StaticGet_JSDocInfo_SOURCE_COMMENT().Get(this_))
}

// Java method: JSDocInfo.getOriginalCommentPosition()
func (this_ *JSDocInfo) GetOriginalCommentPosition() int32 {
	var pos *j2gort.BoxInt = j2gort.CastPtr[j2gort.BoxInt](StaticGet_JSDocInfo_ORIGINAL_COMMENT_POSITION().Get(this_))
	_ = pos
	if pos != nil {
		return j2gort.UnboxInt(pos)
	} else {
		return int32(0)
	}
}

// Java method: JSDocInfo.modifiesThis()
func (this_ *JSDocInfo) ModifiesThis() bool {
	return j2gort.SetContainsStr(this_.GetModifies(), "this")
}

// Java method: JSDocInfo.hasSideEffectsArgumentsAnnotation()
func (this_ *JSDocInfo) HasSideEffectsArgumentsAnnotation() bool {
	var modifies j2gort.Ext_java_util_Set = this_.GetModifies()
	_ = modifies
	return ((modifies.Size() > int32(1)) || ((modifies.Size() == int32(1)) && !(j2gort.SetContainsStr(modifies, "this"))))
}

// Java method: JSDocInfo.getLogTypeInCompiler()
func (this_ *JSDocInfo) GetLogTypeInCompiler() bool {
	return this_.priv_JSDocInfo_checkBit_7975(StaticGet_JSDocInfo_Bit_LOG_TYPE_IN_COMPILER())
}

// Java method: JSDocInfo.getSingleTypeForBit(int bit)
func (this_ *JSDocInfo) priv_JSDocInfo_getSingleTypeForBit_8107(bit int32) JSTypeExpression {
	_ = bit
	var mask int64 = javalang.Shl64(int64(1), int(bit))
	_ = mask
	var typeBits int64 = (this_.propertyKeysBitset & StaticGet_JSDocInfo_SINGLE_TYPE_EXPR_PROPERTIES_MASK())
	_ = typeBits
	if (typeBits & -(typeBits)) == mask {
		return this_.type0
	}
	var nonType0Keys int64 = (this_.propertyKeysBitset ^ (typeBits & -(typeBits)))
	_ = nonType0Keys
	return j2gort.Cast[JSTypeExpression](this_.priv_JSDocInfo_getNonType0ValueByIndex_8110(j2gort.BoxLong_BitCount((nonType0Keys & (mask - int64(1))))))
}

// Java method: JSDocInfo.getNonType0ValueForBit(int bit)
func (this_ *JSDocInfo) priv_JSDocInfo_getNonType0ValueForBit_8108(bit int32) any {
	_ = bit
	var mask int64 = javalang.Shl64(int64(1), int(bit))
	_ = mask
	var typeBits int64 = (this_.propertyKeysBitset & StaticGet_JSDocInfo_SINGLE_TYPE_EXPR_PROPERTIES_MASK())
	_ = typeBits
	var nonType0Keys int64 = (this_.propertyKeysBitset ^ (typeBits & -(typeBits)))
	_ = nonType0Keys
	return this_.priv_JSDocInfo_getNonType0ValueByIndex_8110(j2gort.BoxLong_BitCount((nonType0Keys & (mask - int64(1)))))
}

// Java method: JSDocInfo.getPropertyValueForBit(int bit)
func (this_ *JSDocInfo) priv_JSDocInfo_getPropertyValueForBit_8109(bit int32) any {
	_ = bit
	if ((javalang.Shl64(int64(1), int(bit))) & StaticGet_JSDocInfo_SINGLE_TYPE_EXPR_PROPERTIES_MASK()) != int64(0) {
		return j2gort.BoxNullableVal(this_.priv_JSDocInfo_getSingleTypeForBit_8107(bit))
	}
	return this_.priv_JSDocInfo_getNonType0ValueForBit_8108(bit)
}

// Java method: JSDocInfo.getNonType0ValueByIndex(int index)
func (this_ *JSDocInfo) priv_JSDocInfo_getNonType0ValueByIndex_8110(index int32) any {
	_ = index
	if this_.propertyValue0 == nil {
		panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", "no property value")))
	}
	if index == int32(0) {
		return this_.propertyValue0
	}
	var rest []any
	_ = rest
	var ok_23 bool
	rest, ok_23 = j2gort.TryCast[[]any](this_.propertyValue1OrRest)
	if ok_23 {
		return rest[(index - int32(1))]
	}
	if (index == int32(1)) && (this_.propertyValue1OrRest != nil) {
		return this_.propertyValue1OrRest
	}
	panic(j2gort.AsIface[javalang.IArrayIndexOutOfBoundsException](j2gort.ConstructShimThrowable("java.lang.ArrayIndexOutOfBoundsException", j2gort.ValueOfInt(index))))
}

// Java method: JSDocInfo.getPropertyValueByIndex(int index)
func (this_ *JSDocInfo) GetPropertyValueByIndex(index int32) any {
	_ = index
	var rem int64 = this_.propertyKeysBitset
	_ = rem
	{
		var i int32 = int32(0)
		_ = i
		for ; (i < index) && (rem != int64(0)); i = i + 1 {
			rem = (rem & (rem - int64(1)))
		}
	}
	if rem == int64(0) {
		panic(j2gort.AsIface[javalang.IArrayIndexOutOfBoundsException](j2gort.ConstructShimThrowable("java.lang.ArrayIndexOutOfBoundsException", j2gort.ValueOfInt(index))))
	}
	return this_.priv_JSDocInfo_getPropertyValueForBit_8109(j2gort.BoxLong_NumberOfTrailingZeros(rem))
}

// Java static method: JSDocInfo.internSet(java.util.Set<java.lang.String> strings)
func JSDocInfo_InternSet(strings j2gort.Ext_java_util_Set) j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_JSDocInfo()
	_ = strings
	if strings.IsEmpty() {
		return j2gort.Ext_com_google_common_collect_ImmutableSet_Of_0()
	}
	if strings.Size() == int32(1) {
		return j2gort.ImmSetOfStr1(RhinoStringPool_AddOrGet(j2gort.CastString(j2gort.Ext_com_google_common_collect_Iterables_GetOnlyElement_1(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](strings)))))
	}
	var interned j2gort.Ext_com_google_common_collect_ImmutableSet_Builder = j2gort.Ext_com_google_common_collect_ImmutableSet_BuilderWithExpectedSize(strings.Size())
	_ = interned
	{
		cur_24 := j2gort.NewIterCursor(j2gort.CheckNotNull(strings))
		for cur_24.HasNext() {
			s := cur_24.NextStr()
			_ = s
			j2gort.ImmSetBuilderAddStrVoid(interned, RhinoStringPool_AddOrGet(s))
		}
	}
	return j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](interned.Build())
}

// Java polymorphic class interface: com.google.javascript.rhino.JSDocInfo.Property
type JSDocInfo_Property interface {
	Implements_rhino_JSDocInfo_Property()
	IsJavaNil_() bool
	GetJSDocInfo_PropertyBase() *JSDocInfo_PropertyBase
	Get(*JSDocInfo) any
	Clone(any, JSDocInfo_TypeTransform) any
	IsDefault(any) bool
	EqualValues(any, any) bool
	HasAnyTypeExpression(any) bool
	CollectTypeExpressions(any, *j2gort.Ext_com_google_common_collect_ImmutableList_Builder)
	CollectTypeNodes(any, j2gort.Ext_java_util_ArrayList) j2gort.Ext_java_util_ArrayList
	CountTypeNodes(any) int32
	TypeNodeAt(any, int32) Node
	CompareTo(any) int32
}

func ToIface_JSDocInfo_Property[T interface {
	JSDocInfo_Property
	comparable
}](v T) JSDocInfo_Property {
	var z T
	if v == z {
		return nil
	}
	return v
}

// Java class: com.google.javascript.rhino.JSDocInfo.Property
type JSDocInfo_PropertyBase struct {
	self_   JSDocInfo_Property
	name    string `gson:"name"`
	mask    int64  `gson:"mask"`
	bit     int32  `gson:"bit"`
	idHash_ atomic.Int32
}

func (this_ *JSDocInfo_PropertyBase) BindSelf_(s any) {
	_ = s
	this_.self_ = s.(JSDocInfo_Property)
}
func (this_ *JSDocInfo_PropertyBase) IsJavaNil_() bool { return this_ == nil }
func (*JSDocInfo_PropertyBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.JSDocInfo$Property", "Property"
}
func (this_ *JSDocInfo_PropertyBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.JSDocInfo.Property", "com.google.javascript.rhino.JSDocInfo$Property", func(v_ any) bool {
		return j2gort.InstanceOf[JSDocInfo_Property](v_)
	})
}
func (b *JSDocInfo_PropertyBase) GetJSDocInfo_PropertyBase() *JSDocInfo_PropertyBase { return b }
func (*JSDocInfo_PropertyBase) Implements_rhino_JSDocInfo_Property()                 {}

var static_JSDocInfo_Property_bitCounter int32

func StaticGet_JSDocInfo_Property_BitCounter() int32 {
	clinit_JSDocInfo_Property()
	return static_JSDocInfo_Property_bitCounter
}
func StaticSet_JSDocInfo_Property_BitCounter(v int32) {
	clinit_JSDocInfo_Property()
	static_JSDocInfo_Property_bitCounter = v
}

var static_JSDocInfo_Property_values []JSDocInfo_Property

func StaticGet_JSDocInfo_Property_Values() []JSDocInfo_Property {
	clinit_JSDocInfo_Property()
	return static_JSDocInfo_Property_values
}
func StaticSet_JSDocInfo_Property_Values(v []JSDocInfo_Property) {
	clinit_JSDocInfo_Property()
	static_JSDocInfo_Property_values = v
}

var clinitDone_JSDocInfo_Property uint32
var clinitState_JSDocInfo_Property j2gort.ClinitState

func clinit_JSDocInfo_Property() {
	if atomic.LoadUint32(&clinitDone_JSDocInfo_Property) != 0 {
		return
	}
	clinitSlow_JSDocInfo_Property()
}

func clinitSlow_JSDocInfo_Property() {
	if !j2gort.EnterClinit(&clinitState_JSDocInfo_Property, &clinitDone_JSDocInfo_Property) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_JSDocInfo_Property, &clinitDone_JSDocInfo_Property)
	static_JSDocInfo_Property_bitCounter = int32(0)
	static_JSDocInfo_Property_values = make([]JSDocInfo_Property, int32(64))
}

func (this_ *JSDocInfo_PropertyBase) initFields_JSDocInfo_Property() {
	_ = this_
	this_.name = j2gort.NullString
}

// Java constructor: Property(java.lang.String name)
func NewJSDocInfo_Property(name string) JSDocInfo_Property {
	clinit_JSDocInfo_Property()
	this_ := &JSDocInfo_PropertyBase{}
	this_.BindSelf_(this_)
	InitJSDocInfo_Property(this_, name)
	return this_
}

func InitJSDocInfo_Property(this_ *JSDocInfo_PropertyBase, name string) {
	_ = name
	this_.initFields_JSDocInfo_Property()
	this_.name = name
	this_.bit = func() int32 {
		old_25 := StaticGet_JSDocInfo_Property_BitCounter()
		nw_26 := old_25 + 1
		StaticSet_JSDocInfo_Property_BitCounter(nw_26)
		return old_25
	}()
	this_.mask = javalang.Shl64(int64(1), int(this_.bit))
	if this_.bit > int32(63) {
		panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError", "Too many Properties")))
	}
	{
		a_27 := StaticGet_JSDocInfo_Property_Values()
		i_28 := this_.bit
		a_27[i_28] = this_.self_
	}
}

// Java method: Property.get(com.google.javascript.rhino.JSDocInfo info)
func (this_ *JSDocInfo_PropertyBase) Get(info *JSDocInfo) any {
	_ = info
	if (info.propertyKeysBitset & this_.mask) == int64(0) {
		return nil
	}
	return info.priv_JSDocInfo_getPropertyValueForBit_8109(this_.bit)
}

// Java method: Property.clone(T arg, com.google.javascript.rhino.JSDocInfo.@org.jspecify.annotations.Nullable TypeTransform unusedTransform)
func (this_ *JSDocInfo_PropertyBase) Clone(arg_29 any, unusedTransform JSDocInfo_TypeTransform) any {
	arg := arg_29
	_ = arg
	_ = unusedTransform
	return arg
}

// Java method: Property.isDefault(T value)
func (this_ *JSDocInfo_PropertyBase) IsDefault(arg_30 any) bool {
	value := arg_30
	_ = value
	return ((((value == nil) || j2gort.RefEquals(value, StaticGet_JSDocInfo_Visibility_INHERITED())) || (j2gort.InstanceOf[j2gort.Ext_java_util_Collection](value) && (j2gort.AsIface[j2gort.Ext_java_util_Collection](value)).IsEmpty())) || (j2gort.InstanceOf[j2gort.Ext_java_util_Map](value) && (j2gort.AsIface[j2gort.Ext_java_util_Map](value)).IsEmpty()))
}

// Java method: Property.equalValues(T left, T right)
func (this_ *JSDocInfo_PropertyBase) EqualValues(arg_31 any, arg_32 any) bool {
	left := arg_31
	_ = left
	right := arg_32
	_ = right
	return javalang.ObjectsEquals(j2gort.CheckNotNull(left), right)
}

// Java method: Property.hasAnyTypeExpression(T unusedValue)
func (this_ *JSDocInfo_PropertyBase) HasAnyTypeExpression(arg_33 any) bool {
	unusedValue := arg_33
	_ = unusedValue
	return false
}

// Java method: Property.collectTypeExpressions(T unusedValue, com.google.common.collect.ImmutableList.Builder<com.google.javascript.rhino.JSTypeExpression> unusedSink)
func (this_ *JSDocInfo_PropertyBase) CollectTypeExpressions(arg_34 any, unusedSink *j2gort.Ext_com_google_common_collect_ImmutableList_Builder) {
	unusedValue := arg_34
	_ = unusedValue
	_ = unusedSink
}

// Java method: Property.collectTypeNodes(T unusedValue, java.util.@org.jspecify.annotations.Nullable ArrayList<com.google.javascript.rhino.Node> sink)
func (this_ *JSDocInfo_PropertyBase) CollectTypeNodes(arg_35 any, sink j2gort.Ext_java_util_ArrayList) j2gort.Ext_java_util_ArrayList {
	unusedValue := arg_35
	_ = unusedValue
	_ = sink
	return sink
}

// Java method: Property.countTypeNodes(T unusedValue)
func (this_ *JSDocInfo_PropertyBase) CountTypeNodes(arg_36 any) int32 {
	unusedValue := arg_36
	_ = unusedValue
	return int32(0)
}

// Java method: Property.typeNodeAt(T unusedValue, int unusedOffset)
func (this_ *JSDocInfo_PropertyBase) TypeNodeAt(arg_37 any, unusedOffset int32) Node {
	unusedValue := arg_37
	_ = unusedValue
	_ = unusedOffset
	return Node(0)
}

// Java method: Property.compareTo(com.google.javascript.rhino.JSDocInfo.Property<java.lang.Object> that)
func (this_ *JSDocInfo_PropertyBase) CompareTo(arg_38 any) int32 {
	that := j2gort.AsIface[JSDocInfo_Property](arg_38)
	_ = that
	return (this_.bit - that.GetJSDocInfo_PropertyBase().bit)
}

// Java class: com.google.javascript.rhino.JSDocInfo.MarkerListProperty
type JSDocInfo_MarkerListProperty struct {
	JSDocInfo_PropertyBase
}

func (this_ *JSDocInfo_MarkerListProperty) BindSelf_(s any) {
	_ = s
	this_.JSDocInfo_PropertyBase.BindSelf_(s)
}
func (this_ *JSDocInfo_MarkerListProperty) IsJavaNil_() bool { return this_ == nil }
func (*JSDocInfo_MarkerListProperty) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.JSDocInfo$MarkerListProperty", "MarkerListProperty"
}
func (this_ *JSDocInfo_MarkerListProperty) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.JSDocInfo_PropertyBase.IdentityHash_()
}
func (*JSDocInfo_MarkerListProperty) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.JSDocInfo.MarkerListProperty", "com.google.javascript.rhino.JSDocInfo$MarkerListProperty", func(v_ any) bool {
		return j2gort.InstanceOf[*JSDocInfo_MarkerListProperty](v_)
	})
}
func (*JSDocInfo_MarkerListProperty) Implements_rhino_JSDocInfo_Property() {}

func (this_ *JSDocInfo_MarkerListProperty) initFields_JSDocInfo_MarkerListProperty() {
	_ = this_
}

// Java constructor: MarkerListProperty(java.lang.String name)
func NewJSDocInfo_MarkerListProperty(name string) *JSDocInfo_MarkerListProperty {
	this_ := &JSDocInfo_MarkerListProperty{}
	this_.BindSelf_(this_)
	InitJSDocInfo_MarkerListProperty(this_, name)
	return this_
}

func InitJSDocInfo_MarkerListProperty(this_ *JSDocInfo_MarkerListProperty, name string) {
	_ = name
	InitJSDocInfo_Property(&this_.JSDocInfo_PropertyBase, name)
	this_.initFields_JSDocInfo_MarkerListProperty()
}

// Java method: MarkerListProperty.equalValues(java.util.ArrayList<com.google.javascript.rhino.JSDocInfo.Marker> a, java.util.ArrayList<com.google.javascript.rhino.JSDocInfo.Marker> b)
func (this_ *JSDocInfo_MarkerListProperty) EqualValues(arg_39 any, arg_40 any) bool {
	a := j2gort.AsIface[j2gort.Ext_java_util_ArrayList](arg_39)
	_ = a
	b := j2gort.AsIface[j2gort.Ext_java_util_ArrayList](arg_40)
	_ = b
	if a.Size() != b.Size() {
		return false
	}
	{
		var i int32 = int32(0)
		_ = i
		for ; i < a.Size(); i = i + 1 {
			var m1 *JSDocInfo_Marker = j2gort.CastPtr[JSDocInfo_Marker](a.Get(i))
			_ = m1
			var m2 *JSDocInfo_Marker = j2gort.CastPtr[JSDocInfo_Marker](b.Get(i))
			_ = m2
			if ((m1 == nil) != (m2 == nil)) || ((m1 != nil) && !(m1.priv_Marker_isEquivalentTo_8169(m2))) {
				return false
			}
		}
	}
	return true
}

// Java class: com.google.javascript.rhino.JSDocInfo.TypeProperty
type JSDocInfo_TypeProperty struct {
	JSDocInfo_PropertyBase
}

func (this_ *JSDocInfo_TypeProperty) BindSelf_(s any) {
	_ = s
	this_.JSDocInfo_PropertyBase.BindSelf_(s)
}
func (this_ *JSDocInfo_TypeProperty) IsJavaNil_() bool { return this_ == nil }
func (*JSDocInfo_TypeProperty) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.JSDocInfo$TypeProperty", "TypeProperty"
}
func (this_ *JSDocInfo_TypeProperty) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.JSDocInfo_PropertyBase.IdentityHash_()
}
func (*JSDocInfo_TypeProperty) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.JSDocInfo.TypeProperty", "com.google.javascript.rhino.JSDocInfo$TypeProperty", func(v_ any) bool {
		return j2gort.InstanceOf[*JSDocInfo_TypeProperty](v_)
	})
}
func (*JSDocInfo_TypeProperty) Implements_rhino_JSDocInfo_Property() {}

func (this_ *JSDocInfo_TypeProperty) initFields_JSDocInfo_TypeProperty() {
	_ = this_
}

// Java constructor: TypeProperty(java.lang.String name)
func NewJSDocInfo_TypeProperty(name string) *JSDocInfo_TypeProperty {
	this_ := &JSDocInfo_TypeProperty{}
	this_.BindSelf_(this_)
	InitJSDocInfo_TypeProperty(this_, name)
	return this_
}

func InitJSDocInfo_TypeProperty(this_ *JSDocInfo_TypeProperty, name string) {
	_ = name
	InitJSDocInfo_Property(&this_.JSDocInfo_PropertyBase, name)
	this_.initFields_JSDocInfo_TypeProperty()
}

// Java method: TypeProperty.getTyped(com.google.javascript.rhino.JSDocInfo info)
func (this_ *JSDocInfo_TypeProperty) GetTyped(info *JSDocInfo) JSTypeExpression {
	_ = info
	if (info.propertyKeysBitset & this_.GetJSDocInfo_PropertyBase().mask) == int64(0) {
		return (JSTypeExpression{})
	}
	return info.priv_JSDocInfo_getSingleTypeForBit_8107(this_.GetJSDocInfo_PropertyBase().bit)
}

// Java method: TypeProperty.get(com.google.javascript.rhino.JSDocInfo info)
func (this_ *JSDocInfo_TypeProperty) Get(info *JSDocInfo) any {
	_ = info
	return j2gort.BoxNullableVal(this_.GetTyped(info))
}

// Java method: TypeProperty.clone(com.google.javascript.rhino.JSTypeExpression arg, com.google.javascript.rhino.JSDocInfo.TypeTransform transform)
func (this_ *JSDocInfo_TypeProperty) Clone(arg_41 any, transform JSDocInfo_TypeTransform) any {
	arg := j2gort.Cast[JSTypeExpression](arg_41)
	_ = arg
	_ = transform
	if (arg != (JSTypeExpression{})) && (transform != nil && !transform.IsJavaNil_()) {
		return j2gort.BoxNullableVal(transform.Apply(arg))
	} else {
		return j2gort.BoxNullableVal(arg)
	}
}

// Java method: TypeProperty.equalValues(com.google.javascript.rhino.JSTypeExpression left, com.google.javascript.rhino.JSTypeExpression right)
func (this_ *JSDocInfo_TypeProperty) EqualValues(arg_42 any, arg_43 any) bool {
	left := j2gort.Cast[JSTypeExpression](arg_42)
	_ = left
	right := j2gort.Cast[JSTypeExpression](arg_43)
	_ = right
	return left.IsEquivalentTo(right)
}

// Java method: TypeProperty.hasAnyTypeExpression(com.google.javascript.rhino.JSTypeExpression type)
func (this_ *JSDocInfo_TypeProperty) HasAnyTypeExpression(arg_44 any) bool {
	type_ := j2gort.Cast[JSTypeExpression](arg_44)
	_ = type_
	return (type_ != (JSTypeExpression{}))
}

// Java method: TypeProperty.collectTypeExpressions(com.google.javascript.rhino.JSTypeExpression type, com.google.common.collect.ImmutableList.Builder<com.google.javascript.rhino.JSTypeExpression> sink)
func (this_ *JSDocInfo_TypeProperty) CollectTypeExpressions(arg_45 any, sink *j2gort.Ext_com_google_common_collect_ImmutableList_Builder) {
	type_ := j2gort.Cast[JSTypeExpression](arg_45)
	_ = type_
	_ = sink
	if type_ != (JSTypeExpression{}) {
		sink.AddVoid_(j2gort.BoxNullableVal(type_))
	}
}

// Java method: TypeProperty.collectTypeNodes(com.google.javascript.rhino.JSTypeExpression type, java.util.@org.jspecify.annotations.Nullable ArrayList<com.google.javascript.rhino.Node> sink)
func (this_ *JSDocInfo_TypeProperty) CollectTypeNodes(arg_46 any, sink j2gort.Ext_java_util_ArrayList) j2gort.Ext_java_util_ArrayList {
	type_ := j2gort.Cast[JSTypeExpression](arg_46)
	_ = type_
	_ = sink
	if type_ != (JSTypeExpression{}) {
		if sink == nil {
			sink = j2gort.NewExt_java_util_ArrayList_1_int32(int32(4))
		}
		sink.Add_1(j2gort.BoxNode(type_.GetRoot()))
	}
	return sink
}

// Java method: TypeProperty.countTypeNodes(com.google.javascript.rhino.JSTypeExpression type)
func (this_ *JSDocInfo_TypeProperty) CountTypeNodes(arg_47 any) int32 {
	type_ := j2gort.Cast[JSTypeExpression](arg_47)
	_ = type_
	if type_ != (JSTypeExpression{}) {
		return int32(1)
	} else {
		return int32(0)
	}
}

// Java method: TypeProperty.typeNodeAt(com.google.javascript.rhino.JSTypeExpression type, int offset)
func (this_ *JSDocInfo_TypeProperty) TypeNodeAt(arg_48 any, offset int32) Node {
	type_ := j2gort.Cast[JSTypeExpression](arg_48)
	_ = type_
	_ = offset
	if (offset == int32(0)) && (type_ != (JSTypeExpression{})) {
		return type_.GetRoot()
	} else {
		return Node(0)
	}
}

// Java class: com.google.javascript.rhino.JSDocInfo.TypeList
type JSDocInfo_TypeList struct {
	values  []JSTypeExpression `gson:"values"`
	size    int32              `gson:"size"`
	idHash_ atomic.Int32
}

func (this_ *JSDocInfo_TypeList) BindSelf_(s any) {
	_ = s
}
func (this_ *JSDocInfo_TypeList) IsJavaNil_() bool { return this_ == nil }
func (*JSDocInfo_TypeList) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.JSDocInfo$TypeList", "TypeList"
}
func (this_ *JSDocInfo_TypeList) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*JSDocInfo_TypeList) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.JSDocInfo.TypeList", "com.google.javascript.rhino.JSDocInfo$TypeList", func(v_ any) bool {
		return j2gort.InstanceOf[*JSDocInfo_TypeList](v_)
	})
}

func (this_ *JSDocInfo_TypeList) initFields_JSDocInfo_TypeList() {
	_ = this_
}

// Java constructor: TypeList(int initialCapacity)
func NewJSDocInfo_TypeList(initialCapacity int32) *JSDocInfo_TypeList {
	this_ := &JSDocInfo_TypeList{}
	InitJSDocInfo_TypeList(this_, initialCapacity)
	return this_
}

func InitJSDocInfo_TypeList(this_ *JSDocInfo_TypeList, initialCapacity int32) {
	_ = initialCapacity
	this_.initFields_JSDocInfo_TypeList()
	this_.values = make([]JSTypeExpression, max(int32(2), initialCapacity))
	this_.size = int32(0)
}

// Java method: TypeList.size()
func (this_ *JSDocInfo_TypeList) Size() int32 {
	return this_.size
}

// Java method: TypeList.isEmpty()
func (this_ *JSDocInfo_TypeList) IsEmpty() bool {
	return (this_.size == int32(0))
}

// Java method: TypeList.get(int index)
func (this_ *JSDocInfo_TypeList) Get(index int32) JSTypeExpression {
	_ = index
	if (index < int32(0)) || (index >= this_.size) {
		panic(j2gort.AsIface[javalang.IIndexOutOfBoundsException](j2gort.ConstructShimThrowable("java.lang.IndexOutOfBoundsException", j2gort.ValueOfInt(index))))
	}
	return this_.values[index]
}

// Java method: TypeList.addUnique(com.google.javascript.rhino.JSTypeExpression elem)
func (this_ *JSDocInfo_TypeList) AddUnique(elem JSTypeExpression) bool {
	_ = elem
	{
		var i int32 = int32(0)
		_ = i
		for ; i < this_.size; i = i + 1 {
			if elem.IsEquivalentTo(this_.values[i]) {
				return false
			}
		}
	}
	if this_.size == int32(len(this_.values)) {
		this_.values = func() []JSTypeExpression {
			arrSrc_49 := this_.values
			arrLen_50 := (this_.size * int32(2))
			j2gort.CheckNotNullSlice(arrSrc_49)
			if arrLen_50 < 0 {
				panic(j2gort.ConstructShimThrowable("java.lang.NegativeArraySizeException"))
			}
			arrOut_51 := make([]JSTypeExpression, arrLen_50)
			copy(arrOut_51, arrSrc_49)
			return arrOut_51
		}()
	}
	{
		a_52 := this_.values
		i_53 := func() int32 { old_54 := this_.size; nw_55 := old_54 + 1; this_.size = nw_55; return old_54 }()
		a_52[i_53] = elem
	}
	return true
}

// Java method: TypeList.toUnmodifiableList()
func (this_ *JSDocInfo_TypeList) ToUnmodifiableList() j2gort.Ext_java_util_List {
	if this_.size == int32(0) {
		return j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.Ext_com_google_common_collect_ImmutableList_Of_0())
	}
	if this_.size == int32(1) {
		return j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.Ext_com_google_common_collect_ImmutableList_Of_1(j2gort.BoxNullableVal(this_.values[int32(0)])))
	}
	if this_.size == int32(2) {
		return j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.Ext_com_google_common_collect_ImmutableList_Of_2(j2gort.BoxNullableVal(this_.values[int32(0)]), j2gort.BoxNullableVal(this_.values[int32(1)])))
	}
	var builder *j2gort.Ext_com_google_common_collect_ImmutableList_Builder = j2gort.InitImmListBuilderCap(&j2gort.Ext_com_google_common_collect_ImmutableList_Builder{}, this_.size)
	_ = builder
	{
		var i int32 = int32(0)
		_ = i
		for ; i < this_.size; i = i + 1 {
			builder.AddVoid_(j2gort.BoxNullableVal(this_.values[i]))
		}
	}
	return j2gort.AsIface[j2gort.Ext_java_util_List](builder.BuildImmList_())
}

// Java method: TypeList.toString()
func (this_ *JSDocInfo_TypeList) String() string {
	if this_.size == int32(0) {
		return "[]"
	}
	var sb *j2gort.Ext_java_lang_StringBuilder = j2gort.NewExt_java_lang_StringBuilder()
	_ = sb
	sb.Append_1_uint16(uint16(91))
	{
		var i int32 = int32(0)
		_ = i
		for ; i < this_.size; i = i + 1 {
			if i > int32(0) {
				sb.Append_1_string(", ")
			}
			sb.Append_1_any_9796(j2gort.BoxNullableVal(this_.values[i]))
		}
	}
	sb.Append_1_uint16(uint16(93))
	return sb.String()
}

// Java class: com.google.javascript.rhino.JSDocInfo.TypeListProperty
type JSDocInfo_TypeListProperty struct {
	JSDocInfo_PropertyBase
}

func (this_ *JSDocInfo_TypeListProperty) BindSelf_(s any) {
	_ = s
	this_.JSDocInfo_PropertyBase.BindSelf_(s)
}
func (this_ *JSDocInfo_TypeListProperty) IsJavaNil_() bool { return this_ == nil }
func (*JSDocInfo_TypeListProperty) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.JSDocInfo$TypeListProperty", "TypeListProperty"
}
func (this_ *JSDocInfo_TypeListProperty) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.JSDocInfo_PropertyBase.IdentityHash_()
}
func (*JSDocInfo_TypeListProperty) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.JSDocInfo.TypeListProperty", "com.google.javascript.rhino.JSDocInfo$TypeListProperty", func(v_ any) bool {
		return j2gort.InstanceOf[*JSDocInfo_TypeListProperty](v_)
	})
}
func (*JSDocInfo_TypeListProperty) Implements_rhino_JSDocInfo_Property() {}

func (this_ *JSDocInfo_TypeListProperty) initFields_JSDocInfo_TypeListProperty() {
	_ = this_
}

// Java constructor: TypeListProperty(java.lang.String name)
func NewJSDocInfo_TypeListProperty(name string) *JSDocInfo_TypeListProperty {
	this_ := &JSDocInfo_TypeListProperty{}
	this_.BindSelf_(this_)
	InitJSDocInfo_TypeListProperty(this_, name)
	return this_
}

func InitJSDocInfo_TypeListProperty(this_ *JSDocInfo_TypeListProperty, name string) {
	_ = name
	InitJSDocInfo_Property(&this_.JSDocInfo_PropertyBase, name)
	this_.initFields_JSDocInfo_TypeListProperty()
}

// Java method: TypeListProperty.getUnmodifiable(com.google.javascript.rhino.JSDocInfo info)
func (this_ *JSDocInfo_TypeListProperty) GetUnmodifiable(info *JSDocInfo) j2gort.Ext_java_util_List {
	_ = info
	var value *JSDocInfo_TypeList = j2gort.CastPtr[JSDocInfo_TypeList](this_.Get(info))
	_ = value
	if value != nil {
		return value.ToUnmodifiableList()
	} else {
		return j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.Ext_com_google_common_collect_ImmutableList_Of_0())
	}
}

// Java method: TypeListProperty.isDefault(com.google.javascript.rhino.JSDocInfo.TypeList value)
func (this_ *JSDocInfo_TypeListProperty) IsDefault(arg_56 any) bool {
	value := j2gort.CastPtr[JSDocInfo_TypeList](arg_56)
	_ = value
	return ((value == nil) || value.IsEmpty())
}

// Java method: TypeListProperty.clone(com.google.javascript.rhino.JSDocInfo.TypeList arg, com.google.javascript.rhino.JSDocInfo.TypeTransform transform)
func (this_ *JSDocInfo_TypeListProperty) Clone(arg_57 any, transform JSDocInfo_TypeTransform) any {
	arg := j2gort.CastPtr[JSDocInfo_TypeList](arg_57)
	_ = arg
	_ = transform
	var n int32 = arg.size
	_ = n
	var val_58 JSDocInfo_TypeList
	InitJSDocInfo_TypeList(&val_58, n)
	var out *JSDocInfo_TypeList = &val_58
	_ = out
	{
		r_59 := out
		r_59.size = n
	}
	{
		var i int32 = int32(0)
		_ = i
		for ; i < n; i = i + 1 {
			var elem JSTypeExpression = arg.values[i]
			_ = elem
			{
				a_60 := out.values
				i_61 := i
				a_60[i_61] = func() JSTypeExpression {
					if (elem != (JSTypeExpression{})) && (transform != nil && !transform.IsJavaNil_()) {
						return transform.Apply(elem)
					}
					return elem
				}()
			}
		}
	}
	return out
}

// Java method: TypeListProperty.equalValues(com.google.javascript.rhino.JSDocInfo.TypeList left, com.google.javascript.rhino.JSDocInfo.TypeList right)
func (this_ *JSDocInfo_TypeListProperty) EqualValues(arg_62 any, arg_63 any) bool {
	left := j2gort.CastPtr[JSDocInfo_TypeList](arg_62)
	_ = left
	right := j2gort.CastPtr[JSDocInfo_TypeList](arg_63)
	_ = right
	var n int32 = left.size
	_ = n
	if n != right.size {
		return false
	}
	{
		var i int32 = int32(0)
		_ = i
		for ; i < n; i = i + 1 {
			if !(left.values[i].IsEquivalentTo(right.values[i])) {
				return false
			}
		}
	}
	return true
}

// Java method: TypeListProperty.hasAnyTypeExpression(com.google.javascript.rhino.JSDocInfo.TypeList types)
func (this_ *JSDocInfo_TypeListProperty) HasAnyTypeExpression(arg_64 any) bool {
	types := j2gort.CastPtr[JSDocInfo_TypeList](arg_64)
	_ = types
	{
		var i int32 = int32(0)
		_ = i
		var n int32 = types.size
		_ = n
		for ; i < n; i = i + 1 {
			if types.values[i] != (JSTypeExpression{}) {
				return true
			}
		}
	}
	return false
}

// Java method: TypeListProperty.collectTypeExpressions(com.google.javascript.rhino.JSDocInfo.TypeList types, com.google.common.collect.ImmutableList.Builder<com.google.javascript.rhino.JSTypeExpression> sink)
func (this_ *JSDocInfo_TypeListProperty) CollectTypeExpressions(arg_65 any, sink *j2gort.Ext_com_google_common_collect_ImmutableList_Builder) {
	types := j2gort.CastPtr[JSDocInfo_TypeList](arg_65)
	_ = types
	_ = sink
	{
		var i int32 = int32(0)
		_ = i
		var n int32 = types.size
		_ = n
		for ; i < n; i = i + 1 {
			var type_ JSTypeExpression = types.values[i]
			_ = type_
			if type_ != (JSTypeExpression{}) {
				sink.AddVoid_(j2gort.BoxNullableVal(type_))
			}
		}
	}
}

// Java method: TypeListProperty.collectTypeNodes(com.google.javascript.rhino.JSDocInfo.TypeList types, java.util.@org.jspecify.annotations.Nullable ArrayList<com.google.javascript.rhino.Node> sink)
func (this_ *JSDocInfo_TypeListProperty) CollectTypeNodes(arg_66 any, sink j2gort.Ext_java_util_ArrayList) j2gort.Ext_java_util_ArrayList {
	types := j2gort.CastPtr[JSDocInfo_TypeList](arg_66)
	_ = types
	_ = sink
	var n int32 = types.size
	_ = n
	{
		var i int32 = int32(0)
		_ = i
		for ; i < n; i = i + 1 {
			var type_ JSTypeExpression = types.values[i]
			_ = type_
			if type_ != (JSTypeExpression{}) {
				if sink == nil {
					sink = j2gort.NewExt_java_util_ArrayList_1_int32(max(int32(4), n))
				}
				sink.Add_1(j2gort.BoxNode(type_.GetRoot()))
			}
		}
	}
	return sink
}

// Java method: TypeListProperty.countTypeNodes(com.google.javascript.rhino.JSDocInfo.TypeList types)
func (this_ *JSDocInfo_TypeListProperty) CountTypeNodes(arg_67 any) int32 {
	types := j2gort.CastPtr[JSDocInfo_TypeList](arg_67)
	_ = types
	var count int32 = int32(0)
	_ = count
	{
		var i int32 = int32(0)
		_ = i
		var n int32 = types.size
		_ = n
		for ; i < n; i = i + 1 {
			if types.values[i] != (JSTypeExpression{}) {
				count = count + 1
			}
		}
	}
	return count
}

// Java method: TypeListProperty.typeNodeAt(com.google.javascript.rhino.JSDocInfo.TypeList types, int offset)
func (this_ *JSDocInfo_TypeListProperty) TypeNodeAt(arg_68 any, offset int32) Node {
	types := j2gort.CastPtr[JSDocInfo_TypeList](arg_68)
	_ = types
	_ = offset
	var rem int32 = offset
	_ = rem
	{
		var i int32 = int32(0)
		_ = i
		var n int32 = types.size
		_ = n
		for ; i < n; i = i + 1 {
			var type_ JSTypeExpression = types.values[i]
			_ = type_
			if type_ != (JSTypeExpression{}) {
				if rem == int32(0) {
					return type_.GetRoot()
				}
				rem = rem - 1
			}
		}
	}
	return Node(0)
}

// Java class: com.google.javascript.rhino.JSDocInfo.TypeMap
type JSDocInfo_TypeMap struct {
	keys    []string           `gson:"keys"`
	values  []JSTypeExpression `gson:"values"`
	idHash_ atomic.Int32
}

func (this_ *JSDocInfo_TypeMap) BindSelf_(s any) {
	_ = s
}
func (this_ *JSDocInfo_TypeMap) IsJavaNil_() bool { return this_ == nil }
func (*JSDocInfo_TypeMap) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.JSDocInfo$TypeMap", "TypeMap"
}
func (this_ *JSDocInfo_TypeMap) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*JSDocInfo_TypeMap) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.JSDocInfo.TypeMap", "com.google.javascript.rhino.JSDocInfo$TypeMap", func(v_ any) bool {
		return j2gort.InstanceOf[*JSDocInfo_TypeMap](v_)
	})
}

func (this_ *JSDocInfo_TypeMap) initFields_JSDocInfo_TypeMap() {
	_ = this_
}

// Java constructor: TypeMap(int initialCapacity)
func NewJSDocInfo_TypeMap(initialCapacity int32) *JSDocInfo_TypeMap {
	this_ := &JSDocInfo_TypeMap{}
	InitJSDocInfo_TypeMap(this_, initialCapacity)
	return this_
}

func InitJSDocInfo_TypeMap(this_ *JSDocInfo_TypeMap, initialCapacity int32) {
	_ = initialCapacity
	this_.initFields_JSDocInfo_TypeMap()
	this_.keys = j2gort.MakeNonNilSlice[string](initialCapacity)
	this_.values = j2gort.MakeNonNilSlice[JSTypeExpression](initialCapacity)
}

// Java method: TypeMap.size()
func (this_ *JSDocInfo_TypeMap) Size() int32 {
	return int32(len(this_.keys))
}

// Java method: TypeMap.isEmpty()
func (this_ *JSDocInfo_TypeMap) IsEmpty() bool {
	return len(this_.keys) == 0
}

// Java method: TypeMap.indexOfKey(java.lang.String key)
func (this_ *JSDocInfo_TypeMap) IndexOfKey(key string) int32 {
	_ = key
	{
		var i int32 = int32(0)
		_ = i
		var n int32 = int32(len(this_.keys))
		_ = n
		for ; i < n; i = i + 1 {
			if j2gort.CheckNotNullString(this_.keys[i]) == key {
				return i
			}
		}
	}
	return int32(-1)
}

// Java method: TypeMap.containsKey(java.lang.String key)
func (this_ *JSDocInfo_TypeMap) ContainsKey(key string) bool {
	_ = key
	return (this_.IndexOfKey(key) >= int32(0))
}

// Java method: TypeMap.get(java.lang.String key)
func (this_ *JSDocInfo_TypeMap) Get(key string) JSTypeExpression {
	_ = key
	var idx int32 = this_.IndexOfKey(key)
	_ = idx
	if idx >= int32(0) {
		return this_.values[idx]
	} else {
		return (JSTypeExpression{})
	}
}

// Java method: TypeMap.keyAt(int index)
func (this_ *JSDocInfo_TypeMap) KeyAt(index int32) string {
	_ = index
	return this_.keys[index]
}

// Java method: TypeMap.putIfAbsent(java.lang.String key, com.google.javascript.rhino.@org.jspecify.annotations.Nullable JSTypeExpression value)
func (this_ *JSDocInfo_TypeMap) PutIfAbsent(key string, value JSTypeExpression) bool {
	_ = key
	_ = value
	var idx int32 = this_.IndexOfKey(key)
	_ = idx
	if idx < int32(0) {
		this_.keys = append(this_.keys, key)
		this_.values = append(this_.values, value)
		return true
	}
	if this_.values[idx] == (JSTypeExpression{}) {
		this_.values[idx] = value
		return true
	}
	return false
}

// Java method: TypeMap.keysAsImmutableSet()
func (this_ *JSDocInfo_TypeMap) KeysAsImmutableSet() j2gort.Ext_com_google_common_collect_ImmutableSet {
	var n int32 = int32(len(this_.keys))
	_ = n
	if n == int32(0) {
		return j2gort.Ext_com_google_common_collect_ImmutableSet_Of_0()
	}
	if n == int32(1) {
		var k0 string = this_.keys[int32(0)]
		_ = k0
		return j2gort.ImmSetOfStr1(k0)
	}
	var builder j2gort.Ext_com_google_common_collect_ImmutableSet_Builder = j2gort.Ext_com_google_common_collect_ImmutableSet_BuilderWithExpectedSize(n)
	_ = builder
	{
		var i int32 = int32(0)
		_ = i
		for ; i < n; i = i + 1 {
			var k string = this_.keys[i]
			_ = k
			j2gort.ImmSetBuilderAddStrVoid(builder, k)
		}
	}
	return j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](builder.Build())
}

// Java method: TypeMap.keysAsImmutableList()
func (this_ *JSDocInfo_TypeMap) KeysAsImmutableList() j2gort.Ext_com_google_common_collect_ImmutableList {
	return j2gort.ImmListFromSlice(this_.keys)
}

// Java method: TypeMap.toImmutableMap()
func (this_ *JSDocInfo_TypeMap) ToImmutableMap() j2gort.Ext_com_google_common_collect_ImmutableMap {
	var n int32 = int32(len(this_.keys))
	_ = n
	if n == int32(0) {
		return j2gort.Ext_com_google_common_collect_ImmutableMap_Of_0()
	}
	if n == int32(1) {
		var k0 string = this_.keys[int32(0)]
		_ = k0
		var v0 JSTypeExpression = this_.values[int32(0)]
		_ = v0
		return j2gort.Ext_com_google_common_collect_ImmutableMap_Of_2(j2gort.BoxString(k0), j2gort.BoxNullableVal(v0))
	}
	var builder j2gort.Ext_com_google_common_collect_ImmutableMap_Builder = j2gort.Ext_com_google_common_collect_ImmutableMap_BuilderWithExpectedSize(n)
	_ = builder
	{
		var i int32 = int32(0)
		_ = i
		for ; i < n; i = i + 1 {
			var k string = this_.keys[i]
			_ = k
			var v JSTypeExpression = this_.values[i]
			_ = v
			builder.Put_2(j2gort.BoxString(k), j2gort.BoxNullableVal(v))
		}
	}
	return j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap](builder.BuildOrThrow())
}

// Java method: TypeMap.toString()
func (this_ *JSDocInfo_TypeMap) String() string {
	var n int32 = int32(len(this_.keys))
	_ = n
	if n == int32(0) {
		return "{}"
	}
	var sb *j2gort.Ext_java_lang_StringBuilder = j2gort.NewExt_java_lang_StringBuilder()
	_ = sb
	sb.Append_1_uint16(uint16(123))
	{
		var i int32 = int32(0)
		_ = i
		for ; i < n; i = i + 1 {
			if i > int32(0) {
				sb.Append_1_string(", ")
			}
			var k string = this_.keys[i]
			_ = k
			var v JSTypeExpression = this_.values[i]
			_ = v
			j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](sb.Append_1_string(k)).Append_1_uint16(uint16(61))).Append_1_any_9796(j2gort.BoxNullableVal(v))
		}
	}
	sb.Append_1_uint16(uint16(125))
	return sb.String()
}

// Java class: com.google.javascript.rhino.JSDocInfo.TypeMapProperty
type JSDocInfo_TypeMapProperty struct {
	JSDocInfo_PropertyBase
}

func (this_ *JSDocInfo_TypeMapProperty) BindSelf_(s any) {
	_ = s
	this_.JSDocInfo_PropertyBase.BindSelf_(s)
}
func (this_ *JSDocInfo_TypeMapProperty) IsJavaNil_() bool { return this_ == nil }
func (*JSDocInfo_TypeMapProperty) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.JSDocInfo$TypeMapProperty", "TypeMapProperty"
}
func (this_ *JSDocInfo_TypeMapProperty) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.JSDocInfo_PropertyBase.IdentityHash_()
}
func (*JSDocInfo_TypeMapProperty) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.JSDocInfo.TypeMapProperty", "com.google.javascript.rhino.JSDocInfo$TypeMapProperty", func(v_ any) bool {
		return j2gort.InstanceOf[*JSDocInfo_TypeMapProperty](v_)
	})
}
func (*JSDocInfo_TypeMapProperty) Implements_rhino_JSDocInfo_Property() {}

func (this_ *JSDocInfo_TypeMapProperty) initFields_JSDocInfo_TypeMapProperty() {
	_ = this_
}

// Java constructor: TypeMapProperty(java.lang.String name)
func NewJSDocInfo_TypeMapProperty(name string) *JSDocInfo_TypeMapProperty {
	this_ := &JSDocInfo_TypeMapProperty{}
	this_.BindSelf_(this_)
	InitJSDocInfo_TypeMapProperty(this_, name)
	return this_
}

func InitJSDocInfo_TypeMapProperty(this_ *JSDocInfo_TypeMapProperty, name string) {
	_ = name
	InitJSDocInfo_Property(&this_.JSDocInfo_PropertyBase, name)
	this_.initFields_JSDocInfo_TypeMapProperty()
}

// Java method: TypeMapProperty.isDefault(com.google.javascript.rhino.JSDocInfo.TypeMap value)
func (this_ *JSDocInfo_TypeMapProperty) IsDefault(arg_69 any) bool {
	value := j2gort.CastPtr[JSDocInfo_TypeMap](arg_69)
	_ = value
	return ((value == nil) || value.IsEmpty())
}

// Java method: TypeMapProperty.clone(com.google.javascript.rhino.JSDocInfo.TypeMap arg, com.google.javascript.rhino.JSDocInfo.TypeTransform transform)
func (this_ *JSDocInfo_TypeMapProperty) Clone(arg_70 any, transform JSDocInfo_TypeTransform) any {
	arg := j2gort.CastPtr[JSDocInfo_TypeMap](arg_70)
	_ = arg
	_ = transform
	var n int32 = arg.Size()
	_ = n
	var val_71 JSDocInfo_TypeMap
	InitJSDocInfo_TypeMap(&val_71, n)
	var out *JSDocInfo_TypeMap = &val_71
	_ = out
	{
		var i int32 = int32(0)
		_ = i
		for ; i < n; i = i + 1 {
			var k string = arg.keys[i]
			_ = k
			var elem JSTypeExpression = arg.values[i]
			_ = elem
			if (elem != (JSTypeExpression{})) && (transform != nil && !transform.IsJavaNil_()) {
				elem = transform.Apply(elem)
			}
			out.keys = append(out.keys, k)
			out.values = append(out.values, elem)
		}
	}
	return out
}

// Java method: TypeMapProperty.equalValues(com.google.javascript.rhino.JSDocInfo.TypeMap left, com.google.javascript.rhino.JSDocInfo.TypeMap right)
func (this_ *JSDocInfo_TypeMapProperty) EqualValues(arg_72 any, arg_73 any) bool {
	left := j2gort.CastPtr[JSDocInfo_TypeMap](arg_72)
	_ = left
	right := j2gort.CastPtr[JSDocInfo_TypeMap](arg_73)
	_ = right
	var n int32 = left.Size()
	_ = n
	if n != right.Size() {
		return false
	}
	{
		var i int32 = int32(0)
		_ = i
		for ; i < n; i = i + 1 {
			var rightIdx int32 = right.IndexOfKey(left.keys[i])
			_ = rightIdx
			if (rightIdx < int32(0)) || !(this_.priv_TypeMapProperty_areEquivalent_8140(left.values[i], right.values[rightIdx])) {
				return false
			}
		}
	}
	return true
}

// Java method: TypeMapProperty.areEquivalent(com.google.javascript.rhino.JSTypeExpression expr1, com.google.javascript.rhino.JSTypeExpression expr2)
func (this_ *JSDocInfo_TypeMapProperty) priv_TypeMapProperty_areEquivalent_8140(expr1 JSTypeExpression, expr2 JSTypeExpression) bool {
	_ = expr1
	_ = expr2
	return ((expr1 == expr2) || ((expr1 != (JSTypeExpression{})) && expr1.IsEquivalentTo(expr2)))
}

// Java method: TypeMapProperty.hasAnyTypeExpression(com.google.javascript.rhino.JSDocInfo.TypeMap map)
func (this_ *JSDocInfo_TypeMapProperty) HasAnyTypeExpression(arg_74 any) bool {
	map_ := j2gort.CastPtr[JSDocInfo_TypeMap](arg_74)
	_ = map_
	{
		var i int32 = int32(0)
		_ = i
		var n int32 = int32(len(map_.values))
		_ = n
		for ; i < n; i = i + 1 {
			if map_.values[i] != (JSTypeExpression{}) {
				return true
			}
		}
	}
	return false
}

// Java method: TypeMapProperty.collectTypeExpressions(com.google.javascript.rhino.JSDocInfo.TypeMap map, com.google.common.collect.ImmutableList.Builder<com.google.javascript.rhino.JSTypeExpression> sink)
func (this_ *JSDocInfo_TypeMapProperty) CollectTypeExpressions(arg_75 any, sink *j2gort.Ext_com_google_common_collect_ImmutableList_Builder) {
	map_ := j2gort.CastPtr[JSDocInfo_TypeMap](arg_75)
	_ = map_
	_ = sink
	{
		var i int32 = int32(0)
		_ = i
		var n int32 = int32(len(map_.values))
		_ = n
		for ; i < n; i = i + 1 {
			var type_ JSTypeExpression = map_.values[i]
			_ = type_
			if type_ != (JSTypeExpression{}) {
				sink.AddVoid_(j2gort.BoxNullableVal(type_))
			}
		}
	}
}

// Java method: TypeMapProperty.collectTypeNodes(com.google.javascript.rhino.JSDocInfo.TypeMap map, java.util.@org.jspecify.annotations.Nullable ArrayList<com.google.javascript.rhino.Node> sink)
func (this_ *JSDocInfo_TypeMapProperty) CollectTypeNodes(arg_76 any, sink j2gort.Ext_java_util_ArrayList) j2gort.Ext_java_util_ArrayList {
	map_ := j2gort.CastPtr[JSDocInfo_TypeMap](arg_76)
	_ = map_
	_ = sink
	var n int32 = int32(len(map_.values))
	_ = n
	{
		var i int32 = int32(0)
		_ = i
		for ; i < n; i = i + 1 {
			var type_ JSTypeExpression = map_.values[i]
			_ = type_
			if type_ != (JSTypeExpression{}) {
				if sink == nil {
					sink = j2gort.NewExt_java_util_ArrayList_1_int32(max(int32(4), n))
				}
				sink.Add_1(j2gort.BoxNode(type_.GetRoot()))
			}
		}
	}
	return sink
}

// Java method: TypeMapProperty.countTypeNodes(com.google.javascript.rhino.JSDocInfo.TypeMap map)
func (this_ *JSDocInfo_TypeMapProperty) CountTypeNodes(arg_77 any) int32 {
	map_ := j2gort.CastPtr[JSDocInfo_TypeMap](arg_77)
	_ = map_
	var count int32 = int32(0)
	_ = count
	{
		var i int32 = int32(0)
		_ = i
		var n int32 = int32(len(map_.values))
		_ = n
		for ; i < n; i = i + 1 {
			if map_.values[i] != (JSTypeExpression{}) {
				count = count + 1
			}
		}
	}
	return count
}

// Java method: TypeMapProperty.typeNodeAt(com.google.javascript.rhino.JSDocInfo.TypeMap map, int offset)
func (this_ *JSDocInfo_TypeMapProperty) TypeNodeAt(arg_78 any, offset int32) Node {
	map_ := j2gort.CastPtr[JSDocInfo_TypeMap](arg_78)
	_ = map_
	_ = offset
	var rem int32 = offset
	_ = rem
	{
		var i int32 = int32(0)
		_ = i
		var n int32 = int32(len(map_.values))
		_ = n
		for ; i < n; i = i + 1 {
			var type_ JSTypeExpression = map_.values[i]
			_ = type_
			if type_ != (JSTypeExpression{}) {
				if rem == int32(0) {
					return type_.GetRoot()
				}
				rem = rem - 1
			}
		}
	}
	return Node(0)
}

// Java class: com.google.javascript.rhino.JSDocInfo.StringMap
type JSDocInfo_StringMap struct {
	keys    []string `gson:"keys"`
	values  []string `gson:"values"`
	idHash_ atomic.Int32
}

func (this_ *JSDocInfo_StringMap) BindSelf_(s any) {
	_ = s
}
func (this_ *JSDocInfo_StringMap) IsJavaNil_() bool { return this_ == nil }
func (*JSDocInfo_StringMap) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.JSDocInfo$StringMap", "StringMap"
}
func (this_ *JSDocInfo_StringMap) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*JSDocInfo_StringMap) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.JSDocInfo.StringMap", "com.google.javascript.rhino.JSDocInfo$StringMap", func(v_ any) bool {
		return j2gort.InstanceOf[*JSDocInfo_StringMap](v_)
	})
}

func (this_ *JSDocInfo_StringMap) initFields_JSDocInfo_StringMap() {
	_ = this_
	this_.keys = make([]string, 0, int32(4))
	this_.values = make([]string, 0, int32(4))
}

// Java constructor: StringMap()
func NewJSDocInfo_StringMap() *JSDocInfo_StringMap {
	this_ := &JSDocInfo_StringMap{}
	InitJSDocInfo_StringMap(this_)
	return this_
}

func InitJSDocInfo_StringMap(this_ *JSDocInfo_StringMap) {
	this_.initFields_JSDocInfo_StringMap()
}

// Java method: StringMap.size()
func (this_ *JSDocInfo_StringMap) Size() int32 {
	return int32(len(this_.keys))
}

// Java method: StringMap.isEmpty()
func (this_ *JSDocInfo_StringMap) IsEmpty() bool {
	return len(this_.keys) == 0
}

// Java method: StringMap.indexOfKey(java.lang.String key)
func (this_ *JSDocInfo_StringMap) IndexOfKey(key string) int32 {
	_ = key
	{
		var i int32 = int32(0)
		_ = i
		var n int32 = int32(len(this_.keys))
		_ = n
		for ; i < n; i = i + 1 {
			if j2gort.CheckNotNullString(this_.keys[i]) == key {
				return i
			}
		}
	}
	return int32(-1)
}

// Java method: StringMap.containsKey(java.lang.String key)
func (this_ *JSDocInfo_StringMap) ContainsKey(key string) bool {
	_ = key
	return (this_.IndexOfKey(key) >= int32(0))
}

// Java method: StringMap.get(java.lang.String key)
func (this_ *JSDocInfo_StringMap) Get(key string) string {
	_ = key
	var idx int32 = this_.IndexOfKey(key)
	_ = idx
	if idx >= int32(0) {
		return this_.values[idx]
	} else {
		return j2gort.NullString
	}
}

// Java method: StringMap.putIfAbsent(java.lang.String key, java.lang.@org.jspecify.annotations.Nullable String value)
func (this_ *JSDocInfo_StringMap) PutIfAbsent(key string, value string) bool {
	_ = key
	_ = value
	var idx int32 = this_.IndexOfKey(key)
	_ = idx
	if idx < int32(0) {
		this_.keys = append(this_.keys, key)
		this_.values = append(this_.values, value)
		return true
	}
	if this_.values[idx] == j2gort.NullString {
		this_.values[idx] = value
		return true
	}
	return false
}

// Java method: StringMap.toString()
func (this_ *JSDocInfo_StringMap) String() string {
	var n int32 = int32(len(this_.keys))
	_ = n
	if n == int32(0) {
		return "{}"
	}
	var sb j2gort.LocalStringBuilder

	sb.AppendChar(uint16(123))
	{
		var i int32 = int32(0)
		_ = i
		for ; i < n; i = i + 1 {
			if i > int32(0) {
				sb.AppendStr(", ")
			}
			var k string = this_.keys[i]
			_ = k
			var v string = this_.values[i]
			_ = v
			sb.AppendStr(k).AppendChar(uint16(61)).AppendStr(v)
		}
	}
	sb.AppendChar(uint16(125))
	return sb.String()
}

// Java class: com.google.javascript.rhino.JSDocInfo.StringMapProperty
type JSDocInfo_StringMapProperty struct {
	JSDocInfo_PropertyBase
}

func (this_ *JSDocInfo_StringMapProperty) BindSelf_(s any) {
	_ = s
	this_.JSDocInfo_PropertyBase.BindSelf_(s)
}
func (this_ *JSDocInfo_StringMapProperty) IsJavaNil_() bool { return this_ == nil }
func (*JSDocInfo_StringMapProperty) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.JSDocInfo$StringMapProperty", "StringMapProperty"
}
func (this_ *JSDocInfo_StringMapProperty) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.JSDocInfo_PropertyBase.IdentityHash_()
}
func (*JSDocInfo_StringMapProperty) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.JSDocInfo.StringMapProperty", "com.google.javascript.rhino.JSDocInfo$StringMapProperty", func(v_ any) bool {
		return j2gort.InstanceOf[*JSDocInfo_StringMapProperty](v_)
	})
}
func (*JSDocInfo_StringMapProperty) Implements_rhino_JSDocInfo_Property() {}

func (this_ *JSDocInfo_StringMapProperty) initFields_JSDocInfo_StringMapProperty() {
	_ = this_
}

// Java constructor: StringMapProperty(java.lang.String name)
func NewJSDocInfo_StringMapProperty(name string) *JSDocInfo_StringMapProperty {
	this_ := &JSDocInfo_StringMapProperty{}
	this_.BindSelf_(this_)
	InitJSDocInfo_StringMapProperty(this_, name)
	return this_
}

func InitJSDocInfo_StringMapProperty(this_ *JSDocInfo_StringMapProperty, name string) {
	_ = name
	InitJSDocInfo_Property(&this_.JSDocInfo_PropertyBase, name)
	this_.initFields_JSDocInfo_StringMapProperty()
}

// Java method: StringMapProperty.isDefault(com.google.javascript.rhino.JSDocInfo.StringMap value)
func (this_ *JSDocInfo_StringMapProperty) IsDefault(arg_79 any) bool {
	value := j2gort.CastPtr[JSDocInfo_StringMap](arg_79)
	_ = value
	return ((value == nil) || value.IsEmpty())
}

// Java method: StringMapProperty.equalValues(com.google.javascript.rhino.JSDocInfo.StringMap left, com.google.javascript.rhino.JSDocInfo.StringMap right)
func (this_ *JSDocInfo_StringMapProperty) EqualValues(arg_80 any, arg_81 any) bool {
	left := j2gort.CastPtr[JSDocInfo_StringMap](arg_80)
	_ = left
	right := j2gort.CastPtr[JSDocInfo_StringMap](arg_81)
	_ = right
	var n int32 = left.Size()
	_ = n
	if n != right.Size() {
		return false
	}
	{
		var i int32 = int32(0)
		_ = i
		for ; i < n; i = i + 1 {
			var rightIdx int32 = right.IndexOfKey(left.keys[i])
			_ = rightIdx
			if (rightIdx < int32(0)) || !(left.values[i] == right.values[rightIdx]) {
				return false
			}
		}
	}
	return true
}

// Java class: com.google.javascript.rhino.JSDocInfo.NodeMap
type JSDocInfo_NodeMap struct {
	keys    []string `gson:"keys"`
	values  []Node   `gson:"values"`
	idHash_ atomic.Int32
}

func (this_ *JSDocInfo_NodeMap) BindSelf_(s any) {
	_ = s
}
func (this_ *JSDocInfo_NodeMap) IsJavaNil_() bool { return this_ == nil }
func (*JSDocInfo_NodeMap) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.JSDocInfo$NodeMap", "NodeMap"
}
func (this_ *JSDocInfo_NodeMap) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*JSDocInfo_NodeMap) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.JSDocInfo.NodeMap", "com.google.javascript.rhino.JSDocInfo$NodeMap", func(v_ any) bool {
		return j2gort.InstanceOf[*JSDocInfo_NodeMap](v_)
	})
}

func (this_ *JSDocInfo_NodeMap) initFields_JSDocInfo_NodeMap() {
	_ = this_
}

// Java constructor: NodeMap(int initialCapacity)
func NewJSDocInfo_NodeMap(initialCapacity int32) *JSDocInfo_NodeMap {
	this_ := &JSDocInfo_NodeMap{}
	InitJSDocInfo_NodeMap(this_, initialCapacity)
	return this_
}

func InitJSDocInfo_NodeMap(this_ *JSDocInfo_NodeMap, initialCapacity int32) {
	_ = initialCapacity
	this_.initFields_JSDocInfo_NodeMap()
	this_.keys = j2gort.MakeNonNilSlice[string](initialCapacity)
	this_.values = j2gort.MakeNonNilSlice[Node](initialCapacity)
}

// Java method: NodeMap.size()
func (this_ *JSDocInfo_NodeMap) Size() int32 {
	return int32(len(this_.keys))
}

// Java method: NodeMap.isEmpty()
func (this_ *JSDocInfo_NodeMap) IsEmpty() bool {
	return len(this_.keys) == 0
}

// Java method: NodeMap.indexOfKey(java.lang.String key)
func (this_ *JSDocInfo_NodeMap) IndexOfKey(key string) int32 {
	_ = key
	{
		var i int32 = int32(0)
		_ = i
		var n int32 = int32(len(this_.keys))
		_ = n
		for ; i < n; i = i + 1 {
			if j2gort.CheckNotNullString(this_.keys[i]) == key {
				return i
			}
		}
	}
	return int32(-1)
}

// Java method: NodeMap.containsKey(java.lang.String key)
func (this_ *JSDocInfo_NodeMap) ContainsKey(key string) bool {
	_ = key
	return (this_.IndexOfKey(key) >= int32(0))
}

// Java method: NodeMap.get(java.lang.String key)
func (this_ *JSDocInfo_NodeMap) Get(key string) Node {
	_ = key
	var idx int32 = this_.IndexOfKey(key)
	_ = idx
	if idx >= int32(0) {
		return this_.values[idx]
	} else {
		return Node(0)
	}
}

// Java method: NodeMap.putIfAbsent(java.lang.String key, com.google.javascript.rhino.Node value)
func (this_ *JSDocInfo_NodeMap) PutIfAbsent(key string, value Node) bool {
	_ = key
	_ = value
	var idx int32 = this_.IndexOfKey(key)
	_ = idx
	if idx < int32(0) {
		this_.keys = append(this_.keys, key)
		this_.values = append(this_.values, value)
		return true
	}
	if this_.values[idx] == 0 {
		this_.values[idx] = value
		return true
	}
	return false
}

// Java method: NodeMap.toImmutableMap()
func (this_ *JSDocInfo_NodeMap) ToImmutableMap() j2gort.Ext_com_google_common_collect_ImmutableMap {
	var n int32 = int32(len(this_.keys))
	_ = n
	if n == int32(0) {
		return j2gort.Ext_com_google_common_collect_ImmutableMap_Of_0()
	}
	if n == int32(1) {
		var k0 string = this_.keys[int32(0)]
		_ = k0
		var v0 Node = this_.values[int32(0)]
		_ = v0
		return j2gort.Ext_com_google_common_collect_ImmutableMap_Of_2(j2gort.BoxString(k0), j2gort.BoxNode(v0))
	}
	var builder j2gort.Ext_com_google_common_collect_ImmutableMap_Builder = j2gort.Ext_com_google_common_collect_ImmutableMap_BuilderWithExpectedSize(n)
	_ = builder
	{
		var i int32 = int32(0)
		_ = i
		for ; i < n; i = i + 1 {
			var k string = this_.keys[i]
			_ = k
			var v Node = this_.values[i]
			_ = v
			builder.Put_2(j2gort.BoxString(k), j2gort.BoxNode(v))
		}
	}
	return j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap](builder.BuildOrThrow())
}

// Java method: NodeMap.toString()
func (this_ *JSDocInfo_NodeMap) String() string {
	var n int32 = int32(len(this_.keys))
	_ = n
	if n == int32(0) {
		return "{}"
	}
	var sb *j2gort.Ext_java_lang_StringBuilder = j2gort.NewExt_java_lang_StringBuilder()
	_ = sb
	sb.Append_1_uint16(uint16(123))
	{
		var i int32 = int32(0)
		_ = i
		for ; i < n; i = i + 1 {
			if i > int32(0) {
				sb.Append_1_string(", ")
			}
			var k string = this_.keys[i]
			_ = k
			var v Node = this_.values[i]
			_ = v
			j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](sb.Append_1_string(k)).Append_1_uint16(uint16(61))).Append_1_any_9796(j2gort.BoxNode(v))
		}
	}
	sb.Append_1_uint16(uint16(125))
	return sb.String()
}

// Java class: com.google.javascript.rhino.JSDocInfo.NodeMapProperty
type JSDocInfo_NodeMapProperty struct {
	JSDocInfo_PropertyBase
}

func (this_ *JSDocInfo_NodeMapProperty) BindSelf_(s any) {
	_ = s
	this_.JSDocInfo_PropertyBase.BindSelf_(s)
}
func (this_ *JSDocInfo_NodeMapProperty) IsJavaNil_() bool { return this_ == nil }
func (*JSDocInfo_NodeMapProperty) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.JSDocInfo$NodeMapProperty", "NodeMapProperty"
}
func (this_ *JSDocInfo_NodeMapProperty) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.JSDocInfo_PropertyBase.IdentityHash_()
}
func (*JSDocInfo_NodeMapProperty) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.JSDocInfo.NodeMapProperty", "com.google.javascript.rhino.JSDocInfo$NodeMapProperty", func(v_ any) bool {
		return j2gort.InstanceOf[*JSDocInfo_NodeMapProperty](v_)
	})
}
func (*JSDocInfo_NodeMapProperty) Implements_rhino_JSDocInfo_Property() {}

func (this_ *JSDocInfo_NodeMapProperty) initFields_JSDocInfo_NodeMapProperty() {
	_ = this_
}

// Java constructor: NodeMapProperty(java.lang.String name)
func NewJSDocInfo_NodeMapProperty(name string) *JSDocInfo_NodeMapProperty {
	this_ := &JSDocInfo_NodeMapProperty{}
	this_.BindSelf_(this_)
	InitJSDocInfo_NodeMapProperty(this_, name)
	return this_
}

func InitJSDocInfo_NodeMapProperty(this_ *JSDocInfo_NodeMapProperty, name string) {
	_ = name
	InitJSDocInfo_Property(&this_.JSDocInfo_PropertyBase, name)
	this_.initFields_JSDocInfo_NodeMapProperty()
}

// Java method: NodeMapProperty.isDefault(com.google.javascript.rhino.JSDocInfo.NodeMap value)
func (this_ *JSDocInfo_NodeMapProperty) IsDefault(arg_82 any) bool {
	value := j2gort.CastPtr[JSDocInfo_NodeMap](arg_82)
	_ = value
	return ((value == nil) || value.IsEmpty())
}

// Java method: NodeMapProperty.clone(com.google.javascript.rhino.JSDocInfo.NodeMap arg, com.google.javascript.rhino.JSDocInfo.TypeTransform transform)
func (this_ *JSDocInfo_NodeMapProperty) Clone(arg_83 any, transform JSDocInfo_TypeTransform) any {
	arg := j2gort.CastPtr[JSDocInfo_NodeMap](arg_83)
	_ = arg
	_ = transform
	var n int32 = arg.Size()
	_ = n
	var val_84 JSDocInfo_NodeMap
	InitJSDocInfo_NodeMap(&val_84, n)
	var out *JSDocInfo_NodeMap = &val_84
	_ = out
	{
		var i int32 = int32(0)
		_ = i
		for ; i < n; i = i + 1 {
			var k string = arg.keys[i]
			_ = k
			var v Node = arg.values[i]
			_ = v
			out.keys = append(out.keys, k)
			out.values = append(out.values, v)
		}
	}
	return out
}

// Java method: NodeMapProperty.equalValues(com.google.javascript.rhino.JSDocInfo.NodeMap left, com.google.javascript.rhino.JSDocInfo.NodeMap right)
func (this_ *JSDocInfo_NodeMapProperty) EqualValues(arg_85 any, arg_86 any) bool {
	left := j2gort.CastPtr[JSDocInfo_NodeMap](arg_85)
	_ = left
	right := j2gort.CastPtr[JSDocInfo_NodeMap](arg_86)
	_ = right
	var n int32 = left.Size()
	_ = n
	if n != right.Size() {
		return false
	}
	{
		var i int32 = int32(0)
		_ = i
		for ; i < n; i = i + 1 {
			var rightValue Node = right.Get(left.keys[i])
			_ = rightValue
			if (rightValue == 0) || !(left.values[i].IsEquivalentTo_1(rightValue)) {
				return false
			}
		}
	}
	return true
}

// Java enum: com.google.javascript.rhino.JSDocInfo.Bit
type JSDocInfo_Bit struct {
	j2gort.EnumBase
	name string
	mask int64
}

func (this_ *JSDocInfo_Bit) GetEnumBase() *j2gort.EnumBase {
	if this_ == nil {
		return nil
	}
	return &this_.EnumBase
}
func (this_ *JSDocInfo_Bit) EnumOrdinal_() int32 {
	if this_ == nil {
		return -1
	}
	return this_.EnumBase.Ordinal_
}
func (*JSDocInfo_Bit) Implements_j2gort_Ext_java_lang_Enum()               {}
func (*JSDocInfo_Bit) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ *JSDocInfo_Bit) initFields_JSDocInfo_Bit() {
	_ = this_
	this_.name = j2gort.NullString
}

func initEnumJSDocInfo_Bit(this_ *JSDocInfo_Bit, enumName_ string, enumOrd_ int32) *JSDocInfo_Bit {
	this_.EnumBase.Name_ = enumName_
	this_.EnumBase.Ordinal_ = enumOrd_
	this_.initFields_JSDocInfo_Bit()
	if j2gort.EnumOrdinal(this_) > int32(63) {
		panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError", "Too many Bits")))
	}
	this_.name = j2gort.Ext_com_google_common_base_CaseFormat_F_UPPER_UNDERSCORE.To(j2gort.Ext_com_google_common_base_CaseFormat_F_LOWER_CAMEL, j2gort.EnumName(this_))
	this_.mask = javalang.Shl64(int64(1), int(j2gort.EnumOrdinal(this_)))
	return this_
}

var enum_JSDocInfo_Bit_INLINE_TYPE *JSDocInfo_Bit

func StaticGet_JSDocInfo_Bit_INLINE_TYPE() *JSDocInfo_Bit {
	clinit_JSDocInfo_Bit()
	return enum_JSDocInfo_Bit_INLINE_TYPE
}

var enum_JSDocInfo_Bit_INCLUDE_DOCUMENTATION *JSDocInfo_Bit

func StaticGet_JSDocInfo_Bit_INCLUDE_DOCUMENTATION() *JSDocInfo_Bit {
	clinit_JSDocInfo_Bit()
	return enum_JSDocInfo_Bit_INCLUDE_DOCUMENTATION
}

var enum_JSDocInfo_Bit_CONST *JSDocInfo_Bit

func StaticGet_JSDocInfo_Bit_CONST() *JSDocInfo_Bit {
	clinit_JSDocInfo_Bit()
	return enum_JSDocInfo_Bit_CONST
}

var enum_JSDocInfo_Bit_CONSTRUCTOR *JSDocInfo_Bit

func StaticGet_JSDocInfo_Bit_CONSTRUCTOR() *JSDocInfo_Bit {
	clinit_JSDocInfo_Bit()
	return enum_JSDocInfo_Bit_CONSTRUCTOR
}

var enum_JSDocInfo_Bit_DEFINE *JSDocInfo_Bit

func StaticGet_JSDocInfo_Bit_DEFINE() *JSDocInfo_Bit {
	clinit_JSDocInfo_Bit()
	return enum_JSDocInfo_Bit_DEFINE
}

var enum_JSDocInfo_Bit_TYPE_SUMMARY *JSDocInfo_Bit

func StaticGet_JSDocInfo_Bit_TYPE_SUMMARY() *JSDocInfo_Bit {
	clinit_JSDocInfo_Bit()
	return enum_JSDocInfo_Bit_TYPE_SUMMARY
}

var enum_JSDocInfo_Bit_FINAL *JSDocInfo_Bit

func StaticGet_JSDocInfo_Bit_FINAL() *JSDocInfo_Bit {
	clinit_JSDocInfo_Bit()
	return enum_JSDocInfo_Bit_FINAL
}

var enum_JSDocInfo_Bit_OVERRIDE *JSDocInfo_Bit

func StaticGet_JSDocInfo_Bit_OVERRIDE() *JSDocInfo_Bit {
	clinit_JSDocInfo_Bit()
	return enum_JSDocInfo_Bit_OVERRIDE
}

var enum_JSDocInfo_Bit_DEPRECATED *JSDocInfo_Bit

func StaticGet_JSDocInfo_Bit_DEPRECATED() *JSDocInfo_Bit {
	clinit_JSDocInfo_Bit()
	return enum_JSDocInfo_Bit_DEPRECATED
}

var enum_JSDocInfo_Bit_INTERFACE *JSDocInfo_Bit

func StaticGet_JSDocInfo_Bit_INTERFACE() *JSDocInfo_Bit {
	clinit_JSDocInfo_Bit()
	return enum_JSDocInfo_Bit_INTERFACE
}

var enum_JSDocInfo_Bit_EXPORT *JSDocInfo_Bit

func StaticGet_JSDocInfo_Bit_EXPORT() *JSDocInfo_Bit {
	clinit_JSDocInfo_Bit()
	return enum_JSDocInfo_Bit_EXPORT
}

var enum_JSDocInfo_Bit_ENHANCED_NAMESPACE *JSDocInfo_Bit

func StaticGet_JSDocInfo_Bit_ENHANCED_NAMESPACE() *JSDocInfo_Bit {
	clinit_JSDocInfo_Bit()
	return enum_JSDocInfo_Bit_ENHANCED_NAMESPACE
}

var enum_JSDocInfo_Bit_NOINLINE *JSDocInfo_Bit

func StaticGet_JSDocInfo_Bit_NOINLINE() *JSDocInfo_Bit {
	clinit_JSDocInfo_Bit()
	return enum_JSDocInfo_Bit_NOINLINE
}

var enum_JSDocInfo_Bit_FILEOVERVIEW *JSDocInfo_Bit

func StaticGet_JSDocInfo_Bit_FILEOVERVIEW() *JSDocInfo_Bit {
	clinit_JSDocInfo_Bit()
	return enum_JSDocInfo_Bit_FILEOVERVIEW
}

var enum_JSDocInfo_Bit_IMPLICITCAST *JSDocInfo_Bit

func StaticGet_JSDocInfo_Bit_IMPLICITCAST() *JSDocInfo_Bit {
	clinit_JSDocInfo_Bit()
	return enum_JSDocInfo_Bit_IMPLICITCAST
}

var enum_JSDocInfo_Bit_NOSIDEEFFECTS *JSDocInfo_Bit

func StaticGet_JSDocInfo_Bit_NOSIDEEFFECTS() *JSDocInfo_Bit {
	clinit_JSDocInfo_Bit()
	return enum_JSDocInfo_Bit_NOSIDEEFFECTS
}

var enum_JSDocInfo_Bit_EXTERNS *JSDocInfo_Bit

func StaticGet_JSDocInfo_Bit_EXTERNS() *JSDocInfo_Bit {
	clinit_JSDocInfo_Bit()
	return enum_JSDocInfo_Bit_EXTERNS
}

var enum_JSDocInfo_Bit_NOCOMPILE *JSDocInfo_Bit

func StaticGet_JSDocInfo_Bit_NOCOMPILE() *JSDocInfo_Bit {
	clinit_JSDocInfo_Bit()
	return enum_JSDocInfo_Bit_NOCOMPILE
}

var enum_JSDocInfo_Bit_NODTS *JSDocInfo_Bit

func StaticGet_JSDocInfo_Bit_NODTS() *JSDocInfo_Bit {
	clinit_JSDocInfo_Bit()
	return enum_JSDocInfo_Bit_NODTS
}

var enum_JSDocInfo_Bit_UNRESTRICTED *JSDocInfo_Bit

func StaticGet_JSDocInfo_Bit_UNRESTRICTED() *JSDocInfo_Bit {
	clinit_JSDocInfo_Bit()
	return enum_JSDocInfo_Bit_UNRESTRICTED
}

var enum_JSDocInfo_Bit_USED_VIA_DOT_CONSTRUCTOR *JSDocInfo_Bit

func StaticGet_JSDocInfo_Bit_USED_VIA_DOT_CONSTRUCTOR() *JSDocInfo_Bit {
	clinit_JSDocInfo_Bit()
	return enum_JSDocInfo_Bit_USED_VIA_DOT_CONSTRUCTOR
}

var enum_JSDocInfo_Bit_STRUCT *JSDocInfo_Bit

func StaticGet_JSDocInfo_Bit_STRUCT() *JSDocInfo_Bit {
	clinit_JSDocInfo_Bit()
	return enum_JSDocInfo_Bit_STRUCT
}

var enum_JSDocInfo_Bit_DICT *JSDocInfo_Bit

func StaticGet_JSDocInfo_Bit_DICT() *JSDocInfo_Bit {
	clinit_JSDocInfo_Bit()
	return enum_JSDocInfo_Bit_DICT
}

var enum_JSDocInfo_Bit_NOCOLLAPSE *JSDocInfo_Bit

func StaticGet_JSDocInfo_Bit_NOCOLLAPSE() *JSDocInfo_Bit {
	clinit_JSDocInfo_Bit()
	return enum_JSDocInfo_Bit_NOCOLLAPSE
}

var enum_JSDocInfo_Bit_RECORD *JSDocInfo_Bit

func StaticGet_JSDocInfo_Bit_RECORD() *JSDocInfo_Bit {
	clinit_JSDocInfo_Bit()
	return enum_JSDocInfo_Bit_RECORD
}

var enum_JSDocInfo_Bit_ABSTRACT *JSDocInfo_Bit

func StaticGet_JSDocInfo_Bit_ABSTRACT() *JSDocInfo_Bit {
	clinit_JSDocInfo_Bit()
	return enum_JSDocInfo_Bit_ABSTRACT
}

var enum_JSDocInfo_Bit_PURE_OR_BREAK_MY_CODE *JSDocInfo_Bit

func StaticGet_JSDocInfo_Bit_PURE_OR_BREAK_MY_CODE() *JSDocInfo_Bit {
	clinit_JSDocInfo_Bit()
	return enum_JSDocInfo_Bit_PURE_OR_BREAK_MY_CODE
}

var enum_JSDocInfo_Bit_COLLAPSIBLE_OR_BREAK_MY_CODE *JSDocInfo_Bit

func StaticGet_JSDocInfo_Bit_COLLAPSIBLE_OR_BREAK_MY_CODE() *JSDocInfo_Bit {
	clinit_JSDocInfo_Bit()
	return enum_JSDocInfo_Bit_COLLAPSIBLE_OR_BREAK_MY_CODE
}

var enum_JSDocInfo_Bit_NOCOVERAGE *JSDocInfo_Bit

func StaticGet_JSDocInfo_Bit_NOCOVERAGE() *JSDocInfo_Bit {
	clinit_JSDocInfo_Bit()
	return enum_JSDocInfo_Bit_NOCOVERAGE
}

var enum_JSDocInfo_Bit_REQUIRE_INLINING *JSDocInfo_Bit

func StaticGet_JSDocInfo_Bit_REQUIRE_INLINING() *JSDocInfo_Bit {
	clinit_JSDocInfo_Bit()
	return enum_JSDocInfo_Bit_REQUIRE_INLINING
}

var enum_JSDocInfo_Bit_ENCOURAGE_INLINING *JSDocInfo_Bit

func StaticGet_JSDocInfo_Bit_ENCOURAGE_INLINING() *JSDocInfo_Bit {
	clinit_JSDocInfo_Bit()
	return enum_JSDocInfo_Bit_ENCOURAGE_INLINING
}

var enum_JSDocInfo_Bit_NG_INJECT *JSDocInfo_Bit

func StaticGet_JSDocInfo_Bit_NG_INJECT() *JSDocInfo_Bit {
	clinit_JSDocInfo_Bit()
	return enum_JSDocInfo_Bit_NG_INJECT
}

var enum_JSDocInfo_Bit_WIZ_ACTION *JSDocInfo_Bit

func StaticGet_JSDocInfo_Bit_WIZ_ACTION() *JSDocInfo_Bit {
	clinit_JSDocInfo_Bit()
	return enum_JSDocInfo_Bit_WIZ_ACTION
}

var enum_JSDocInfo_Bit_POLYMER_BEHAVIOR *JSDocInfo_Bit

func StaticGet_JSDocInfo_Bit_POLYMER_BEHAVIOR() *JSDocInfo_Bit {
	clinit_JSDocInfo_Bit()
	return enum_JSDocInfo_Bit_POLYMER_BEHAVIOR
}

var enum_JSDocInfo_Bit_POLYMER *JSDocInfo_Bit

func StaticGet_JSDocInfo_Bit_POLYMER() *JSDocInfo_Bit {
	clinit_JSDocInfo_Bit()
	return enum_JSDocInfo_Bit_POLYMER
}

var enum_JSDocInfo_Bit_CUSTOM_ELEMENT *JSDocInfo_Bit

func StaticGet_JSDocInfo_Bit_CUSTOM_ELEMENT() *JSDocInfo_Bit {
	clinit_JSDocInfo_Bit()
	return enum_JSDocInfo_Bit_CUSTOM_ELEMENT
}

var enum_JSDocInfo_Bit_MIXIN_CLASS *JSDocInfo_Bit

func StaticGet_JSDocInfo_Bit_MIXIN_CLASS() *JSDocInfo_Bit {
	clinit_JSDocInfo_Bit()
	return enum_JSDocInfo_Bit_MIXIN_CLASS
}

var enum_JSDocInfo_Bit_MIXIN_FUNCTION *JSDocInfo_Bit

func StaticGet_JSDocInfo_Bit_MIXIN_FUNCTION() *JSDocInfo_Bit {
	clinit_JSDocInfo_Bit()
	return enum_JSDocInfo_Bit_MIXIN_FUNCTION
}

var enum_JSDocInfo_Bit_SASS_GENERATED_CSS_TS *JSDocInfo_Bit

func StaticGet_JSDocInfo_Bit_SASS_GENERATED_CSS_TS() *JSDocInfo_Bit {
	clinit_JSDocInfo_Bit()
	return enum_JSDocInfo_Bit_SASS_GENERATED_CSS_TS
}

var enum_JSDocInfo_Bit_CLOSURE_UNAWARE_CODE *JSDocInfo_Bit

func StaticGet_JSDocInfo_Bit_CLOSURE_UNAWARE_CODE() *JSDocInfo_Bit {
	clinit_JSDocInfo_Bit()
	return enum_JSDocInfo_Bit_CLOSURE_UNAWARE_CODE
}

var enum_JSDocInfo_Bit_PROVIDE_GOOG *JSDocInfo_Bit

func StaticGet_JSDocInfo_Bit_PROVIDE_GOOG() *JSDocInfo_Bit {
	clinit_JSDocInfo_Bit()
	return enum_JSDocInfo_Bit_PROVIDE_GOOG
}

var enum_JSDocInfo_Bit_PROVIDE_ALREADY_PROVIDED *JSDocInfo_Bit

func StaticGet_JSDocInfo_Bit_PROVIDE_ALREADY_PROVIDED() *JSDocInfo_Bit {
	clinit_JSDocInfo_Bit()
	return enum_JSDocInfo_Bit_PROVIDE_ALREADY_PROVIDED
}

var enum_JSDocInfo_Bit_WIZ_CALLBACK *JSDocInfo_Bit

func StaticGet_JSDocInfo_Bit_WIZ_CALLBACK() *JSDocInfo_Bit {
	clinit_JSDocInfo_Bit()
	return enum_JSDocInfo_Bit_WIZ_CALLBACK
}

var enum_JSDocInfo_Bit_LOG_TYPE_IN_COMPILER *JSDocInfo_Bit

func StaticGet_JSDocInfo_Bit_LOG_TYPE_IN_COMPILER() *JSDocInfo_Bit {
	clinit_JSDocInfo_Bit()
	return enum_JSDocInfo_Bit_LOG_TYPE_IN_COMPILER
}

var enum_JSDocInfo_Bit_CLOSURE_UNAWARE_CONFIG *JSDocInfo_Bit

func StaticGet_JSDocInfo_Bit_CLOSURE_UNAWARE_CONFIG() *JSDocInfo_Bit {
	clinit_JSDocInfo_Bit()
	return enum_JSDocInfo_Bit_CLOSURE_UNAWARE_CONFIG
}

var enum_JSDocInfo_Bit_LEGACY_NAMED_TS_EXPORT *JSDocInfo_Bit

func StaticGet_JSDocInfo_Bit_LEGACY_NAMED_TS_EXPORT() *JSDocInfo_Bit {
	clinit_JSDocInfo_Bit()
	return enum_JSDocInfo_Bit_LEGACY_NAMED_TS_EXPORT
}

var clinitDone_JSDocInfo_Bit uint32
var clinitState_JSDocInfo_Bit j2gort.ClinitState

func clinit_JSDocInfo_Bit() {
	if atomic.LoadUint32(&clinitDone_JSDocInfo_Bit) != 0 {
		return
	}
	clinitSlow_JSDocInfo_Bit()
}

func clinitSlow_JSDocInfo_Bit() {
	if !j2gort.EnterClinit(&clinitState_JSDocInfo_Bit, &clinitDone_JSDocInfo_Bit) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_JSDocInfo_Bit, &clinitDone_JSDocInfo_Bit)
	enum_JSDocInfo_Bit_INLINE_TYPE = &JSDocInfo_Bit{}
	initEnumJSDocInfo_Bit(enum_JSDocInfo_Bit_INLINE_TYPE, "INLINE_TYPE", int32(0))
	enum_JSDocInfo_Bit_INCLUDE_DOCUMENTATION = &JSDocInfo_Bit{}
	initEnumJSDocInfo_Bit(enum_JSDocInfo_Bit_INCLUDE_DOCUMENTATION, "INCLUDE_DOCUMENTATION", int32(1))
	enum_JSDocInfo_Bit_CONST = &JSDocInfo_Bit{}
	initEnumJSDocInfo_Bit(enum_JSDocInfo_Bit_CONST, "CONST", int32(2))
	enum_JSDocInfo_Bit_CONSTRUCTOR = &JSDocInfo_Bit{}
	initEnumJSDocInfo_Bit(enum_JSDocInfo_Bit_CONSTRUCTOR, "CONSTRUCTOR", int32(3))
	enum_JSDocInfo_Bit_DEFINE = &JSDocInfo_Bit{}
	initEnumJSDocInfo_Bit(enum_JSDocInfo_Bit_DEFINE, "DEFINE", int32(4))
	enum_JSDocInfo_Bit_TYPE_SUMMARY = &JSDocInfo_Bit{}
	initEnumJSDocInfo_Bit(enum_JSDocInfo_Bit_TYPE_SUMMARY, "TYPE_SUMMARY", int32(5))
	enum_JSDocInfo_Bit_FINAL = &JSDocInfo_Bit{}
	initEnumJSDocInfo_Bit(enum_JSDocInfo_Bit_FINAL, "FINAL", int32(6))
	enum_JSDocInfo_Bit_OVERRIDE = &JSDocInfo_Bit{}
	initEnumJSDocInfo_Bit(enum_JSDocInfo_Bit_OVERRIDE, "OVERRIDE", int32(7))
	enum_JSDocInfo_Bit_DEPRECATED = &JSDocInfo_Bit{}
	initEnumJSDocInfo_Bit(enum_JSDocInfo_Bit_DEPRECATED, "DEPRECATED", int32(8))
	enum_JSDocInfo_Bit_INTERFACE = &JSDocInfo_Bit{}
	initEnumJSDocInfo_Bit(enum_JSDocInfo_Bit_INTERFACE, "INTERFACE", int32(9))
	enum_JSDocInfo_Bit_EXPORT = &JSDocInfo_Bit{}
	initEnumJSDocInfo_Bit(enum_JSDocInfo_Bit_EXPORT, "EXPORT", int32(10))
	enum_JSDocInfo_Bit_ENHANCED_NAMESPACE = &JSDocInfo_Bit{}
	initEnumJSDocInfo_Bit(enum_JSDocInfo_Bit_ENHANCED_NAMESPACE, "ENHANCED_NAMESPACE", int32(11))
	enum_JSDocInfo_Bit_NOINLINE = &JSDocInfo_Bit{}
	initEnumJSDocInfo_Bit(enum_JSDocInfo_Bit_NOINLINE, "NOINLINE", int32(12))
	enum_JSDocInfo_Bit_FILEOVERVIEW = &JSDocInfo_Bit{}
	initEnumJSDocInfo_Bit(enum_JSDocInfo_Bit_FILEOVERVIEW, "FILEOVERVIEW", int32(13))
	enum_JSDocInfo_Bit_IMPLICITCAST = &JSDocInfo_Bit{}
	initEnumJSDocInfo_Bit(enum_JSDocInfo_Bit_IMPLICITCAST, "IMPLICITCAST", int32(14))
	enum_JSDocInfo_Bit_NOSIDEEFFECTS = &JSDocInfo_Bit{}
	initEnumJSDocInfo_Bit(enum_JSDocInfo_Bit_NOSIDEEFFECTS, "NOSIDEEFFECTS", int32(15))
	enum_JSDocInfo_Bit_EXTERNS = &JSDocInfo_Bit{}
	initEnumJSDocInfo_Bit(enum_JSDocInfo_Bit_EXTERNS, "EXTERNS", int32(16))
	enum_JSDocInfo_Bit_NOCOMPILE = &JSDocInfo_Bit{}
	initEnumJSDocInfo_Bit(enum_JSDocInfo_Bit_NOCOMPILE, "NOCOMPILE", int32(17))
	enum_JSDocInfo_Bit_NODTS = &JSDocInfo_Bit{}
	initEnumJSDocInfo_Bit(enum_JSDocInfo_Bit_NODTS, "NODTS", int32(18))
	enum_JSDocInfo_Bit_UNRESTRICTED = &JSDocInfo_Bit{}
	initEnumJSDocInfo_Bit(enum_JSDocInfo_Bit_UNRESTRICTED, "UNRESTRICTED", int32(19))
	enum_JSDocInfo_Bit_USED_VIA_DOT_CONSTRUCTOR = &JSDocInfo_Bit{}
	initEnumJSDocInfo_Bit(enum_JSDocInfo_Bit_USED_VIA_DOT_CONSTRUCTOR, "USED_VIA_DOT_CONSTRUCTOR", int32(20))
	enum_JSDocInfo_Bit_STRUCT = &JSDocInfo_Bit{}
	initEnumJSDocInfo_Bit(enum_JSDocInfo_Bit_STRUCT, "STRUCT", int32(21))
	enum_JSDocInfo_Bit_DICT = &JSDocInfo_Bit{}
	initEnumJSDocInfo_Bit(enum_JSDocInfo_Bit_DICT, "DICT", int32(22))
	enum_JSDocInfo_Bit_NOCOLLAPSE = &JSDocInfo_Bit{}
	initEnumJSDocInfo_Bit(enum_JSDocInfo_Bit_NOCOLLAPSE, "NOCOLLAPSE", int32(23))
	enum_JSDocInfo_Bit_RECORD = &JSDocInfo_Bit{}
	initEnumJSDocInfo_Bit(enum_JSDocInfo_Bit_RECORD, "RECORD", int32(24))
	enum_JSDocInfo_Bit_ABSTRACT = &JSDocInfo_Bit{}
	initEnumJSDocInfo_Bit(enum_JSDocInfo_Bit_ABSTRACT, "ABSTRACT", int32(25))
	enum_JSDocInfo_Bit_PURE_OR_BREAK_MY_CODE = &JSDocInfo_Bit{}
	initEnumJSDocInfo_Bit(enum_JSDocInfo_Bit_PURE_OR_BREAK_MY_CODE, "PURE_OR_BREAK_MY_CODE", int32(26))
	enum_JSDocInfo_Bit_COLLAPSIBLE_OR_BREAK_MY_CODE = &JSDocInfo_Bit{}
	initEnumJSDocInfo_Bit(enum_JSDocInfo_Bit_COLLAPSIBLE_OR_BREAK_MY_CODE, "COLLAPSIBLE_OR_BREAK_MY_CODE", int32(27))
	enum_JSDocInfo_Bit_NOCOVERAGE = &JSDocInfo_Bit{}
	initEnumJSDocInfo_Bit(enum_JSDocInfo_Bit_NOCOVERAGE, "NOCOVERAGE", int32(28))
	enum_JSDocInfo_Bit_REQUIRE_INLINING = &JSDocInfo_Bit{}
	initEnumJSDocInfo_Bit(enum_JSDocInfo_Bit_REQUIRE_INLINING, "REQUIRE_INLINING", int32(29))
	enum_JSDocInfo_Bit_ENCOURAGE_INLINING = &JSDocInfo_Bit{}
	initEnumJSDocInfo_Bit(enum_JSDocInfo_Bit_ENCOURAGE_INLINING, "ENCOURAGE_INLINING", int32(30))
	enum_JSDocInfo_Bit_NG_INJECT = &JSDocInfo_Bit{}
	initEnumJSDocInfo_Bit(enum_JSDocInfo_Bit_NG_INJECT, "NG_INJECT", int32(31))
	enum_JSDocInfo_Bit_WIZ_ACTION = &JSDocInfo_Bit{}
	initEnumJSDocInfo_Bit(enum_JSDocInfo_Bit_WIZ_ACTION, "WIZ_ACTION", int32(32))
	enum_JSDocInfo_Bit_POLYMER_BEHAVIOR = &JSDocInfo_Bit{}
	initEnumJSDocInfo_Bit(enum_JSDocInfo_Bit_POLYMER_BEHAVIOR, "POLYMER_BEHAVIOR", int32(33))
	enum_JSDocInfo_Bit_POLYMER = &JSDocInfo_Bit{}
	initEnumJSDocInfo_Bit(enum_JSDocInfo_Bit_POLYMER, "POLYMER", int32(34))
	enum_JSDocInfo_Bit_CUSTOM_ELEMENT = &JSDocInfo_Bit{}
	initEnumJSDocInfo_Bit(enum_JSDocInfo_Bit_CUSTOM_ELEMENT, "CUSTOM_ELEMENT", int32(35))
	enum_JSDocInfo_Bit_MIXIN_CLASS = &JSDocInfo_Bit{}
	initEnumJSDocInfo_Bit(enum_JSDocInfo_Bit_MIXIN_CLASS, "MIXIN_CLASS", int32(36))
	enum_JSDocInfo_Bit_MIXIN_FUNCTION = &JSDocInfo_Bit{}
	initEnumJSDocInfo_Bit(enum_JSDocInfo_Bit_MIXIN_FUNCTION, "MIXIN_FUNCTION", int32(37))
	enum_JSDocInfo_Bit_SASS_GENERATED_CSS_TS = &JSDocInfo_Bit{}
	initEnumJSDocInfo_Bit(enum_JSDocInfo_Bit_SASS_GENERATED_CSS_TS, "SASS_GENERATED_CSS_TS", int32(38))
	enum_JSDocInfo_Bit_CLOSURE_UNAWARE_CODE = &JSDocInfo_Bit{}
	initEnumJSDocInfo_Bit(enum_JSDocInfo_Bit_CLOSURE_UNAWARE_CODE, "CLOSURE_UNAWARE_CODE", int32(39))
	enum_JSDocInfo_Bit_PROVIDE_GOOG = &JSDocInfo_Bit{}
	initEnumJSDocInfo_Bit(enum_JSDocInfo_Bit_PROVIDE_GOOG, "PROVIDE_GOOG", int32(40))
	enum_JSDocInfo_Bit_PROVIDE_ALREADY_PROVIDED = &JSDocInfo_Bit{}
	initEnumJSDocInfo_Bit(enum_JSDocInfo_Bit_PROVIDE_ALREADY_PROVIDED, "PROVIDE_ALREADY_PROVIDED", int32(41))
	enum_JSDocInfo_Bit_WIZ_CALLBACK = &JSDocInfo_Bit{}
	initEnumJSDocInfo_Bit(enum_JSDocInfo_Bit_WIZ_CALLBACK, "WIZ_CALLBACK", int32(42))
	enum_JSDocInfo_Bit_LOG_TYPE_IN_COMPILER = &JSDocInfo_Bit{}
	initEnumJSDocInfo_Bit(enum_JSDocInfo_Bit_LOG_TYPE_IN_COMPILER, "LOG_TYPE_IN_COMPILER", int32(43))
	enum_JSDocInfo_Bit_CLOSURE_UNAWARE_CONFIG = &JSDocInfo_Bit{}
	initEnumJSDocInfo_Bit(enum_JSDocInfo_Bit_CLOSURE_UNAWARE_CONFIG, "CLOSURE_UNAWARE_CONFIG", int32(44))
	enum_JSDocInfo_Bit_LEGACY_NAMED_TS_EXPORT = &JSDocInfo_Bit{}
	initEnumJSDocInfo_Bit(enum_JSDocInfo_Bit_LEGACY_NAMED_TS_EXPORT, "LEGACY_NAMED_TS_EXPORT", int32(45))
	enumValues_JSDocInfo_Bit = []*JSDocInfo_Bit{enum_JSDocInfo_Bit_INLINE_TYPE, enum_JSDocInfo_Bit_INCLUDE_DOCUMENTATION, enum_JSDocInfo_Bit_CONST, enum_JSDocInfo_Bit_CONSTRUCTOR, enum_JSDocInfo_Bit_DEFINE, enum_JSDocInfo_Bit_TYPE_SUMMARY, enum_JSDocInfo_Bit_FINAL, enum_JSDocInfo_Bit_OVERRIDE, enum_JSDocInfo_Bit_DEPRECATED, enum_JSDocInfo_Bit_INTERFACE, enum_JSDocInfo_Bit_EXPORT, enum_JSDocInfo_Bit_ENHANCED_NAMESPACE, enum_JSDocInfo_Bit_NOINLINE, enum_JSDocInfo_Bit_FILEOVERVIEW, enum_JSDocInfo_Bit_IMPLICITCAST, enum_JSDocInfo_Bit_NOSIDEEFFECTS, enum_JSDocInfo_Bit_EXTERNS, enum_JSDocInfo_Bit_NOCOMPILE, enum_JSDocInfo_Bit_NODTS, enum_JSDocInfo_Bit_UNRESTRICTED, enum_JSDocInfo_Bit_USED_VIA_DOT_CONSTRUCTOR, enum_JSDocInfo_Bit_STRUCT, enum_JSDocInfo_Bit_DICT, enum_JSDocInfo_Bit_NOCOLLAPSE, enum_JSDocInfo_Bit_RECORD, enum_JSDocInfo_Bit_ABSTRACT, enum_JSDocInfo_Bit_PURE_OR_BREAK_MY_CODE, enum_JSDocInfo_Bit_COLLAPSIBLE_OR_BREAK_MY_CODE, enum_JSDocInfo_Bit_NOCOVERAGE, enum_JSDocInfo_Bit_REQUIRE_INLINING, enum_JSDocInfo_Bit_ENCOURAGE_INLINING, enum_JSDocInfo_Bit_NG_INJECT, enum_JSDocInfo_Bit_WIZ_ACTION, enum_JSDocInfo_Bit_POLYMER_BEHAVIOR, enum_JSDocInfo_Bit_POLYMER, enum_JSDocInfo_Bit_CUSTOM_ELEMENT, enum_JSDocInfo_Bit_MIXIN_CLASS, enum_JSDocInfo_Bit_MIXIN_FUNCTION, enum_JSDocInfo_Bit_SASS_GENERATED_CSS_TS, enum_JSDocInfo_Bit_CLOSURE_UNAWARE_CODE, enum_JSDocInfo_Bit_PROVIDE_GOOG, enum_JSDocInfo_Bit_PROVIDE_ALREADY_PROVIDED, enum_JSDocInfo_Bit_WIZ_CALLBACK, enum_JSDocInfo_Bit_LOG_TYPE_IN_COMPILER, enum_JSDocInfo_Bit_CLOSURE_UNAWARE_CONFIG, enum_JSDocInfo_Bit_LEGACY_NAMED_TS_EXPORT}
	j2gort.RegisterEnumConstants("com.google.javascript.rhino.JSDocInfo.Bit", []any{enum_JSDocInfo_Bit_INLINE_TYPE, enum_JSDocInfo_Bit_INCLUDE_DOCUMENTATION, enum_JSDocInfo_Bit_CONST, enum_JSDocInfo_Bit_CONSTRUCTOR, enum_JSDocInfo_Bit_DEFINE, enum_JSDocInfo_Bit_TYPE_SUMMARY, enum_JSDocInfo_Bit_FINAL, enum_JSDocInfo_Bit_OVERRIDE, enum_JSDocInfo_Bit_DEPRECATED, enum_JSDocInfo_Bit_INTERFACE, enum_JSDocInfo_Bit_EXPORT, enum_JSDocInfo_Bit_ENHANCED_NAMESPACE, enum_JSDocInfo_Bit_NOINLINE, enum_JSDocInfo_Bit_FILEOVERVIEW, enum_JSDocInfo_Bit_IMPLICITCAST, enum_JSDocInfo_Bit_NOSIDEEFFECTS, enum_JSDocInfo_Bit_EXTERNS, enum_JSDocInfo_Bit_NOCOMPILE, enum_JSDocInfo_Bit_NODTS, enum_JSDocInfo_Bit_UNRESTRICTED, enum_JSDocInfo_Bit_USED_VIA_DOT_CONSTRUCTOR, enum_JSDocInfo_Bit_STRUCT, enum_JSDocInfo_Bit_DICT, enum_JSDocInfo_Bit_NOCOLLAPSE, enum_JSDocInfo_Bit_RECORD, enum_JSDocInfo_Bit_ABSTRACT, enum_JSDocInfo_Bit_PURE_OR_BREAK_MY_CODE, enum_JSDocInfo_Bit_COLLAPSIBLE_OR_BREAK_MY_CODE, enum_JSDocInfo_Bit_NOCOVERAGE, enum_JSDocInfo_Bit_REQUIRE_INLINING, enum_JSDocInfo_Bit_ENCOURAGE_INLINING, enum_JSDocInfo_Bit_NG_INJECT, enum_JSDocInfo_Bit_WIZ_ACTION, enum_JSDocInfo_Bit_POLYMER_BEHAVIOR, enum_JSDocInfo_Bit_POLYMER, enum_JSDocInfo_Bit_CUSTOM_ELEMENT, enum_JSDocInfo_Bit_MIXIN_CLASS, enum_JSDocInfo_Bit_MIXIN_FUNCTION, enum_JSDocInfo_Bit_SASS_GENERATED_CSS_TS, enum_JSDocInfo_Bit_CLOSURE_UNAWARE_CODE, enum_JSDocInfo_Bit_PROVIDE_GOOG, enum_JSDocInfo_Bit_PROVIDE_ALREADY_PROVIDED, enum_JSDocInfo_Bit_WIZ_CALLBACK, enum_JSDocInfo_Bit_LOG_TYPE_IN_COMPILER, enum_JSDocInfo_Bit_CLOSURE_UNAWARE_CONFIG, enum_JSDocInfo_Bit_LEGACY_NAMED_TS_EXPORT})
	j2gort.RegisterEnumConstants("com.google.javascript.rhino.JSDocInfo$Bit", []any{enum_JSDocInfo_Bit_INLINE_TYPE, enum_JSDocInfo_Bit_INCLUDE_DOCUMENTATION, enum_JSDocInfo_Bit_CONST, enum_JSDocInfo_Bit_CONSTRUCTOR, enum_JSDocInfo_Bit_DEFINE, enum_JSDocInfo_Bit_TYPE_SUMMARY, enum_JSDocInfo_Bit_FINAL, enum_JSDocInfo_Bit_OVERRIDE, enum_JSDocInfo_Bit_DEPRECATED, enum_JSDocInfo_Bit_INTERFACE, enum_JSDocInfo_Bit_EXPORT, enum_JSDocInfo_Bit_ENHANCED_NAMESPACE, enum_JSDocInfo_Bit_NOINLINE, enum_JSDocInfo_Bit_FILEOVERVIEW, enum_JSDocInfo_Bit_IMPLICITCAST, enum_JSDocInfo_Bit_NOSIDEEFFECTS, enum_JSDocInfo_Bit_EXTERNS, enum_JSDocInfo_Bit_NOCOMPILE, enum_JSDocInfo_Bit_NODTS, enum_JSDocInfo_Bit_UNRESTRICTED, enum_JSDocInfo_Bit_USED_VIA_DOT_CONSTRUCTOR, enum_JSDocInfo_Bit_STRUCT, enum_JSDocInfo_Bit_DICT, enum_JSDocInfo_Bit_NOCOLLAPSE, enum_JSDocInfo_Bit_RECORD, enum_JSDocInfo_Bit_ABSTRACT, enum_JSDocInfo_Bit_PURE_OR_BREAK_MY_CODE, enum_JSDocInfo_Bit_COLLAPSIBLE_OR_BREAK_MY_CODE, enum_JSDocInfo_Bit_NOCOVERAGE, enum_JSDocInfo_Bit_REQUIRE_INLINING, enum_JSDocInfo_Bit_ENCOURAGE_INLINING, enum_JSDocInfo_Bit_NG_INJECT, enum_JSDocInfo_Bit_WIZ_ACTION, enum_JSDocInfo_Bit_POLYMER_BEHAVIOR, enum_JSDocInfo_Bit_POLYMER, enum_JSDocInfo_Bit_CUSTOM_ELEMENT, enum_JSDocInfo_Bit_MIXIN_CLASS, enum_JSDocInfo_Bit_MIXIN_FUNCTION, enum_JSDocInfo_Bit_SASS_GENERATED_CSS_TS, enum_JSDocInfo_Bit_CLOSURE_UNAWARE_CODE, enum_JSDocInfo_Bit_PROVIDE_GOOG, enum_JSDocInfo_Bit_PROVIDE_ALREADY_PROVIDED, enum_JSDocInfo_Bit_WIZ_CALLBACK, enum_JSDocInfo_Bit_LOG_TYPE_IN_COMPILER, enum_JSDocInfo_Bit_CLOSURE_UNAWARE_CONFIG, enum_JSDocInfo_Bit_LEGACY_NAMED_TS_EXPORT})
}

func init() {
	j2gort.RegisterEnumConstantsLoader("com.google.javascript.rhino.JSDocInfo.Bit", "com.google.javascript.rhino.JSDocInfo$Bit", func() {
		clinit_JSDocInfo_Bit()
	})
}

var enumValues_JSDocInfo_Bit []*JSDocInfo_Bit

func JSDocInfo_Bit_Values() []*JSDocInfo_Bit {
	clinit_JSDocInfo_Bit()
	return j2gort.CloneSlice(enumValues_JSDocInfo_Bit)
}

func JSDocInfo_Bit_ValueOf(name string) *JSDocInfo_Bit {
	clinit_JSDocInfo_Bit()
	for _, c := range enumValues_JSDocInfo_Bit {
		if c.Name() == name {
			return c
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.rhino.JSDocInfo.Bit." + name))
}

func (this_ *JSDocInfo_Bit) IsJavaNil_() bool { return this_ == nil }
func (*JSDocInfo_Bit) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.JSDocInfo$Bit", "Bit"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.JSDocInfo.Bit", "com.google.javascript.rhino.JSDocInfo$Bit", func(v_ any) bool {
		return j2gort.InstanceOf[*JSDocInfo_Bit](v_)
	})
}

// Java interface: com.google.javascript.rhino.JSDocInfo.TypeTransform
type JSDocInfo_TypeTransform interface {
	Implements_rhino_JSDocInfo_TypeTransform()
	IsJavaNil_() bool
	Apply(JSTypeExpression) JSTypeExpression
}

func ToIface_JSDocInfo_TypeTransform[T interface {
	JSDocInfo_TypeTransform
	comparable
}](v T) JSDocInfo_TypeTransform {
	var z T
	if v == z {
		return nil
	}
	return v
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.JSDocInfo.TypeTransform", "com.google.javascript.rhino.JSDocInfo$TypeTransform", func(v_ any) bool {
		return j2gort.InstanceOf[JSDocInfo_TypeTransform](v_)
	})
}

type lambda_JSDocInfo_TypeTransform struct {
	fn func(JSTypeExpression) JSTypeExpression
}

func (l *lambda_JSDocInfo_TypeTransform) IsJavaNil_() bool { return l == nil }
func NewLambda_JSDocInfo_TypeTransform(fn func(JSTypeExpression) JSTypeExpression) JSDocInfo_TypeTransform {
	if fn == nil {
		return nil
	}
	return &lambda_JSDocInfo_TypeTransform{fn: fn}
}

func (*lambda_JSDocInfo_TypeTransform) Implements_rhino_JSDocInfo_TypeTransform() {}
func (l *lambda_JSDocInfo_TypeTransform) Apply(p0_ JSTypeExpression) JSTypeExpression {
	return l.fn(p0_)
}

// Java enum: com.google.javascript.rhino.JSDocInfo.Visibility
type JSDocInfo_Visibility uint8

const (
	enum_JSDocInfo_Visibility_PRIVATE   JSDocInfo_Visibility = 1
	enum_JSDocInfo_Visibility_PACKAGE   JSDocInfo_Visibility = 2
	enum_JSDocInfo_Visibility_PROTECTED JSDocInfo_Visibility = 3
	enum_JSDocInfo_Visibility_PUBLIC    JSDocInfo_Visibility = 4
	enum_JSDocInfo_Visibility_INHERITED JSDocInfo_Visibility = 5
)

var enumBaseTable_JSDocInfo_Visibility = [...]j2gort.EnumBase{
	{Name_: "PRIVATE", Ordinal_: 0},
	{Name_: "PACKAGE", Ordinal_: 1},
	{Name_: "PROTECTED", Ordinal_: 2},
	{Name_: "PUBLIC", Ordinal_: 3},
	{Name_: "INHERITED", Ordinal_: 4},
}
var enumValues_JSDocInfo_Visibility = []JSDocInfo_Visibility{1, 2, 3, 4, 5}

func StaticGet_JSDocInfo_Visibility_PRIVATE() JSDocInfo_Visibility {
	return enum_JSDocInfo_Visibility_PRIVATE
}
func StaticGet_JSDocInfo_Visibility_PACKAGE() JSDocInfo_Visibility {
	return enum_JSDocInfo_Visibility_PACKAGE
}
func StaticGet_JSDocInfo_Visibility_PROTECTED() JSDocInfo_Visibility {
	return enum_JSDocInfo_Visibility_PROTECTED
}
func StaticGet_JSDocInfo_Visibility_PUBLIC() JSDocInfo_Visibility {
	return enum_JSDocInfo_Visibility_PUBLIC
}
func StaticGet_JSDocInfo_Visibility_INHERITED() JSDocInfo_Visibility {
	return enum_JSDocInfo_Visibility_INHERITED
}

func JSDocInfo_Visibility_Values() []JSDocInfo_Visibility {
	return j2gort.CloneSlice(enumValues_JSDocInfo_Visibility)
}

func JSDocInfo_Visibility_ValueOf(name string) JSDocInfo_Visibility {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_JSDocInfo_Visibility {
		if enumBaseTable_JSDocInfo_Visibility[i_].Name_ == name {
			return JSDocInfo_Visibility(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.rhino.JSDocInfo.Visibility." + name))
}

func (this_ JSDocInfo_Visibility) IsJavaNil_() bool                              { return this_ == 0 }
func (JSDocInfo_Visibility) HasIdentityEquals_()                                 {}
func (JSDocInfo_Visibility) Implements_j2gort_Ext_java_lang_Enum()               {}
func (JSDocInfo_Visibility) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ JSDocInfo_Visibility) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (JSDocInfo_Visibility) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.JSDocInfo$Visibility", "Visibility"
}
func (this_ JSDocInfo_Visibility) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_JSDocInfo_Visibility[this_-1]
}
func (this_ JSDocInfo_Visibility) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ JSDocInfo_Visibility) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ JSDocInfo_Visibility) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_JSDocInfo_Visibility[this_-1].Name_
}
func (this_ JSDocInfo_Visibility) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_JSDocInfo_Visibility[this_-1].Name_
}
func (this_ JSDocInfo_Visibility) Equals(other_ any) bool {
	o_, ok_ := other_.(JSDocInfo_Visibility)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ JSDocInfo_Visibility) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_JSDocInfo_Visibility[this_-1].HashCode()
}
func (this_ JSDocInfo_Visibility) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_JSDocInfo_Visibility[this_-1].CompareTo(other_)
}
func (this_ JSDocInfo_Visibility) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_JSDocInfo_Visibility[this_-1].Name_)
}
func init() {
	vals_ := []any{JSDocInfo_Visibility(1), JSDocInfo_Visibility(2), JSDocInfo_Visibility(3), JSDocInfo_Visibility(4), JSDocInfo_Visibility(5)}
	j2gort.RegisterEnumConstants("com.google.javascript.rhino.JSDocInfo.Visibility", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.rhino.JSDocInfo$Visibility", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.JSDocInfo.Visibility", "com.google.javascript.rhino.JSDocInfo$Visibility", func(v_ any) bool {
		return j2gort.InstanceOf[JSDocInfo_Visibility](v_)
	})
}

// Java enum: com.google.javascript.rhino.JSDocInfo.IdGenerator
type JSDocInfo_IdGenerator uint8

const (
	enum_JSDocInfo_IdGenerator_XID        JSDocInfo_IdGenerator = 1
	enum_JSDocInfo_IdGenerator_CONSISTENT JSDocInfo_IdGenerator = 2
	enum_JSDocInfo_IdGenerator_UNIQUE     JSDocInfo_IdGenerator = 3
	enum_JSDocInfo_IdGenerator_STABLE     JSDocInfo_IdGenerator = 4
	enum_JSDocInfo_IdGenerator_MAPPED     JSDocInfo_IdGenerator = 5
)

var enumBaseTable_JSDocInfo_IdGenerator = [...]j2gort.EnumBase{
	{Name_: "XID", Ordinal_: 0},
	{Name_: "CONSISTENT", Ordinal_: 1},
	{Name_: "UNIQUE", Ordinal_: 2},
	{Name_: "STABLE", Ordinal_: 3},
	{Name_: "MAPPED", Ordinal_: 4},
}
var enumValues_JSDocInfo_IdGenerator = []JSDocInfo_IdGenerator{1, 2, 3, 4, 5}

func StaticGet_JSDocInfo_IdGenerator_XID() JSDocInfo_IdGenerator {
	return enum_JSDocInfo_IdGenerator_XID
}
func StaticGet_JSDocInfo_IdGenerator_CONSISTENT() JSDocInfo_IdGenerator {
	return enum_JSDocInfo_IdGenerator_CONSISTENT
}
func StaticGet_JSDocInfo_IdGenerator_UNIQUE() JSDocInfo_IdGenerator {
	return enum_JSDocInfo_IdGenerator_UNIQUE
}
func StaticGet_JSDocInfo_IdGenerator_STABLE() JSDocInfo_IdGenerator {
	return enum_JSDocInfo_IdGenerator_STABLE
}
func StaticGet_JSDocInfo_IdGenerator_MAPPED() JSDocInfo_IdGenerator {
	return enum_JSDocInfo_IdGenerator_MAPPED
}

func JSDocInfo_IdGenerator_Values() []JSDocInfo_IdGenerator {
	return j2gort.CloneSlice(enumValues_JSDocInfo_IdGenerator)
}

func JSDocInfo_IdGenerator_ValueOf(name string) JSDocInfo_IdGenerator {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_JSDocInfo_IdGenerator {
		if enumBaseTable_JSDocInfo_IdGenerator[i_].Name_ == name {
			return JSDocInfo_IdGenerator(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.rhino.JSDocInfo.IdGenerator." + name))
}

func (this_ JSDocInfo_IdGenerator) IsJavaNil_() bool                              { return this_ == 0 }
func (JSDocInfo_IdGenerator) HasIdentityEquals_()                                 {}
func (JSDocInfo_IdGenerator) Implements_j2gort_Ext_java_lang_Enum()               {}
func (JSDocInfo_IdGenerator) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ JSDocInfo_IdGenerator) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (JSDocInfo_IdGenerator) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.JSDocInfo$IdGenerator", "IdGenerator"
}
func (this_ JSDocInfo_IdGenerator) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_JSDocInfo_IdGenerator[this_-1]
}
func (this_ JSDocInfo_IdGenerator) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ JSDocInfo_IdGenerator) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ JSDocInfo_IdGenerator) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_JSDocInfo_IdGenerator[this_-1].Name_
}
func (this_ JSDocInfo_IdGenerator) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_JSDocInfo_IdGenerator[this_-1].Name_
}
func (this_ JSDocInfo_IdGenerator) Equals(other_ any) bool {
	o_, ok_ := other_.(JSDocInfo_IdGenerator)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ JSDocInfo_IdGenerator) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_JSDocInfo_IdGenerator[this_-1].HashCode()
}
func (this_ JSDocInfo_IdGenerator) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_JSDocInfo_IdGenerator[this_-1].CompareTo(other_)
}
func (this_ JSDocInfo_IdGenerator) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_JSDocInfo_IdGenerator[this_-1].Name_)
}
func init() {
	vals_ := []any{JSDocInfo_IdGenerator(1), JSDocInfo_IdGenerator(2), JSDocInfo_IdGenerator(3), JSDocInfo_IdGenerator(4), JSDocInfo_IdGenerator(5)}
	j2gort.RegisterEnumConstants("com.google.javascript.rhino.JSDocInfo.IdGenerator", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.rhino.JSDocInfo$IdGenerator", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.JSDocInfo.IdGenerator", "com.google.javascript.rhino.JSDocInfo$IdGenerator", func(v_ any) bool {
		return j2gort.InstanceOf[JSDocInfo_IdGenerator](v_)
	})
}

// Java enum: com.google.javascript.rhino.JSDocInfo.PerFileClosureUnawareMode
type JSDocInfo_PerFileClosureUnawareMode uint8

const (
	enum_JSDocInfo_PerFileClosureUnawareMode_UNSPECIFIED JSDocInfo_PerFileClosureUnawareMode = 1
	enum_JSDocInfo_PerFileClosureUnawareMode_SIMPLE      JSDocInfo_PerFileClosureUnawareMode = 2
	enum_JSDocInfo_PerFileClosureUnawareMode_WHITESPACE  JSDocInfo_PerFileClosureUnawareMode = 3
)

var enumBaseTable_JSDocInfo_PerFileClosureUnawareMode = [...]j2gort.EnumBase{
	{Name_: "UNSPECIFIED", Ordinal_: 0},
	{Name_: "SIMPLE", Ordinal_: 1},
	{Name_: "WHITESPACE", Ordinal_: 2},
}
var enumValues_JSDocInfo_PerFileClosureUnawareMode = []JSDocInfo_PerFileClosureUnawareMode{1, 2, 3}

func StaticGet_JSDocInfo_PerFileClosureUnawareMode_UNSPECIFIED() JSDocInfo_PerFileClosureUnawareMode {
	return enum_JSDocInfo_PerFileClosureUnawareMode_UNSPECIFIED
}
func StaticGet_JSDocInfo_PerFileClosureUnawareMode_SIMPLE() JSDocInfo_PerFileClosureUnawareMode {
	return enum_JSDocInfo_PerFileClosureUnawareMode_SIMPLE
}
func StaticGet_JSDocInfo_PerFileClosureUnawareMode_WHITESPACE() JSDocInfo_PerFileClosureUnawareMode {
	return enum_JSDocInfo_PerFileClosureUnawareMode_WHITESPACE
}

func JSDocInfo_PerFileClosureUnawareMode_Values() []JSDocInfo_PerFileClosureUnawareMode {
	return j2gort.CloneSlice(enumValues_JSDocInfo_PerFileClosureUnawareMode)
}

func JSDocInfo_PerFileClosureUnawareMode_ValueOf(name string) JSDocInfo_PerFileClosureUnawareMode {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_JSDocInfo_PerFileClosureUnawareMode {
		if enumBaseTable_JSDocInfo_PerFileClosureUnawareMode[i_].Name_ == name {
			return JSDocInfo_PerFileClosureUnawareMode(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.rhino.JSDocInfo.PerFileClosureUnawareMode." + name))
}

func (this_ JSDocInfo_PerFileClosureUnawareMode) IsJavaNil_() bool                              { return this_ == 0 }
func (JSDocInfo_PerFileClosureUnawareMode) HasIdentityEquals_()                                 {}
func (JSDocInfo_PerFileClosureUnawareMode) Implements_j2gort_Ext_java_lang_Enum()               {}
func (JSDocInfo_PerFileClosureUnawareMode) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ JSDocInfo_PerFileClosureUnawareMode) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (JSDocInfo_PerFileClosureUnawareMode) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.JSDocInfo$PerFileClosureUnawareMode", "PerFileClosureUnawareMode"
}
func (this_ JSDocInfo_PerFileClosureUnawareMode) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_JSDocInfo_PerFileClosureUnawareMode[this_-1]
}
func (this_ JSDocInfo_PerFileClosureUnawareMode) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ JSDocInfo_PerFileClosureUnawareMode) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ JSDocInfo_PerFileClosureUnawareMode) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_JSDocInfo_PerFileClosureUnawareMode[this_-1].Name_
}
func (this_ JSDocInfo_PerFileClosureUnawareMode) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_JSDocInfo_PerFileClosureUnawareMode[this_-1].Name_
}
func (this_ JSDocInfo_PerFileClosureUnawareMode) Equals(other_ any) bool {
	o_, ok_ := other_.(JSDocInfo_PerFileClosureUnawareMode)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ JSDocInfo_PerFileClosureUnawareMode) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_JSDocInfo_PerFileClosureUnawareMode[this_-1].HashCode()
}
func (this_ JSDocInfo_PerFileClosureUnawareMode) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_JSDocInfo_PerFileClosureUnawareMode[this_-1].CompareTo(other_)
}
func (this_ JSDocInfo_PerFileClosureUnawareMode) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_JSDocInfo_PerFileClosureUnawareMode[this_-1].Name_)
}
func init() {
	vals_ := []any{JSDocInfo_PerFileClosureUnawareMode(1), JSDocInfo_PerFileClosureUnawareMode(2), JSDocInfo_PerFileClosureUnawareMode(3)}
	j2gort.RegisterEnumConstants("com.google.javascript.rhino.JSDocInfo.PerFileClosureUnawareMode", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.rhino.JSDocInfo$PerFileClosureUnawareMode", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.JSDocInfo.PerFileClosureUnawareMode", "com.google.javascript.rhino.JSDocInfo$PerFileClosureUnawareMode", func(v_ any) bool {
		return j2gort.InstanceOf[JSDocInfo_PerFileClosureUnawareMode](v_)
	})
}

// Java polymorphic class interface: com.google.javascript.rhino.JSDocInfo.ComparableSourcePosition
type JSDocInfo_ComparableSourcePosition interface {
	Implements_rhino_JSDocInfo_ComparableSourcePosition()
	IsJavaNil_() bool
	GetJSDocInfo_ComparableSourcePositionBase() *JSDocInfo_ComparableSourcePositionBase
	SourcePosition
	IsEquivalentTo(SourcePosition) bool
}

func ToIface_JSDocInfo_ComparableSourcePosition[T interface {
	JSDocInfo_ComparableSourcePosition
	comparable
}](v T) JSDocInfo_ComparableSourcePosition {
	var z T
	if v == z {
		return nil
	}
	return v
}

// Java class: com.google.javascript.rhino.JSDocInfo.ComparableSourcePosition
type JSDocInfo_ComparableSourcePositionBase struct {
	SourcePositionBase
}

func (this_ *JSDocInfo_ComparableSourcePositionBase) BindSelf_(s any) {
	_ = s
}
func (this_ *JSDocInfo_ComparableSourcePositionBase) IsJavaNil_() bool { return this_ == nil }
func (*JSDocInfo_ComparableSourcePositionBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.JSDocInfo$ComparableSourcePosition", "ComparableSourcePosition"
}
func (this_ *JSDocInfo_ComparableSourcePositionBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.SourcePositionBase.IdentityHash_()
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.JSDocInfo.ComparableSourcePosition", "com.google.javascript.rhino.JSDocInfo$ComparableSourcePosition", func(v_ any) bool {
		return j2gort.InstanceOf[JSDocInfo_ComparableSourcePosition](v_)
	})
}
func (b *JSDocInfo_ComparableSourcePositionBase) GetJSDocInfo_ComparableSourcePositionBase() *JSDocInfo_ComparableSourcePositionBase {
	return b
}
func (*JSDocInfo_ComparableSourcePositionBase) Implements_rhino_JSDocInfo_ComparableSourcePosition() {
}
func (*JSDocInfo_ComparableSourcePositionBase) Implements_rhino_SourcePosition() {}

func (this_ *JSDocInfo_ComparableSourcePositionBase) initFields_JSDocInfo_ComparableSourcePosition() {
	_ = this_
}

func InitJSDocInfo_ComparableSourcePosition(this_ *JSDocInfo_ComparableSourcePositionBase) {
	InitSourcePosition(&this_.SourcePositionBase)
	this_.initFields_JSDocInfo_ComparableSourcePosition()
}

// Java method: ComparableSourcePosition.isEquivalentTo(com.google.javascript.rhino.SourcePosition<T> that)
func (this_ *JSDocInfo_ComparableSourcePositionBase) IsEquivalentTo(that SourcePosition) bool {
	_ = that
	panic("abstract or native method invoked")
}

// Java polymorphic class interface: com.google.javascript.rhino.JSDocInfo.StringPosition
type JSDocInfo_StringPosition interface {
	Implements_rhino_JSDocInfo_StringPosition()
	IsJavaNil_() bool
	GetJSDocInfo_StringPositionBase() *JSDocInfo_StringPositionBase
	JSDocInfo_ComparableSourcePosition
}

func ToIface_JSDocInfo_StringPosition[T interface {
	JSDocInfo_StringPosition
	comparable
}](v T) JSDocInfo_StringPosition {
	var z T
	if v == z {
		return nil
	}
	return v
}

// Java class: com.google.javascript.rhino.JSDocInfo.StringPosition
type JSDocInfo_StringPositionBase struct {
	JSDocInfo_ComparableSourcePositionBase
}

func (this_ *JSDocInfo_StringPositionBase) BindSelf_(s any) {
	_ = s
}
func (this_ *JSDocInfo_StringPositionBase) IsJavaNil_() bool { return this_ == nil }
func (*JSDocInfo_StringPositionBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.JSDocInfo$StringPosition", "StringPosition"
}
func (this_ *JSDocInfo_StringPositionBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.JSDocInfo_ComparableSourcePositionBase.IdentityHash_()
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.JSDocInfo.StringPosition", "com.google.javascript.rhino.JSDocInfo$StringPosition", func(v_ any) bool {
		return j2gort.InstanceOf[JSDocInfo_StringPosition](v_)
	})
}
func (b *JSDocInfo_StringPositionBase) GetJSDocInfo_StringPositionBase() *JSDocInfo_StringPositionBase {
	return b
}
func (*JSDocInfo_StringPositionBase) Implements_rhino_JSDocInfo_StringPosition()           {}
func (*JSDocInfo_StringPositionBase) Implements_rhino_JSDocInfo_ComparableSourcePosition() {}
func (*JSDocInfo_StringPositionBase) Implements_rhino_SourcePosition()                     {}

func (this_ *JSDocInfo_StringPositionBase) initFields_JSDocInfo_StringPosition() {
	_ = this_
}

// Java constructor: StringPosition()
func NewJSDocInfo_StringPosition() JSDocInfo_StringPosition {
	this_ := &JSDocInfo_StringPositionBase{}
	InitJSDocInfo_StringPosition(this_)
	return this_
}

func InitJSDocInfo_StringPosition(this_ *JSDocInfo_StringPositionBase) {
	InitJSDocInfo_ComparableSourcePosition(&this_.JSDocInfo_ComparableSourcePositionBase)
	this_.initFields_JSDocInfo_StringPosition()
}

// Java method: StringPosition.isEquivalentTo(com.google.javascript.rhino.SourcePosition<java.lang.String> that)
func (this_ *JSDocInfo_StringPositionBase) IsEquivalentTo(that SourcePosition) bool {
	_ = that
	return (((that != nil && !that.IsJavaNil_()) && this_.IsSamePositionAs(that)) && (j2gort.CastString(this_.GetItem()) == j2gort.CastString(that.GetItem())))
}

// Java class: com.google.javascript.rhino.JSDocInfo.TrimmedStringPosition
type JSDocInfo_TrimmedStringPosition struct {
	JSDocInfo_StringPositionBase
}

func (this_ *JSDocInfo_TrimmedStringPosition) BindSelf_(s any) {
	_ = s
}
func (this_ *JSDocInfo_TrimmedStringPosition) IsJavaNil_() bool { return this_ == nil }
func (*JSDocInfo_TrimmedStringPosition) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.JSDocInfo$TrimmedStringPosition", "TrimmedStringPosition"
}
func (this_ *JSDocInfo_TrimmedStringPosition) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.JSDocInfo_StringPositionBase.IdentityHash_()
}
func (*JSDocInfo_TrimmedStringPosition) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.JSDocInfo.TrimmedStringPosition", "com.google.javascript.rhino.JSDocInfo$TrimmedStringPosition", func(v_ any) bool {
		return j2gort.InstanceOf[*JSDocInfo_TrimmedStringPosition](v_)
	})
}
func (*JSDocInfo_TrimmedStringPosition) Implements_rhino_JSDocInfo_StringPosition()           {}
func (*JSDocInfo_TrimmedStringPosition) Implements_rhino_JSDocInfo_ComparableSourcePosition() {}
func (*JSDocInfo_TrimmedStringPosition) Implements_rhino_SourcePosition()                     {}

func (this_ *JSDocInfo_TrimmedStringPosition) initFields_JSDocInfo_TrimmedStringPosition() {
	_ = this_
}

// Java constructor: TrimmedStringPosition()
func NewJSDocInfo_TrimmedStringPosition() *JSDocInfo_TrimmedStringPosition {
	this_ := &JSDocInfo_TrimmedStringPosition{}
	InitJSDocInfo_TrimmedStringPosition(this_)
	return this_
}

func InitJSDocInfo_TrimmedStringPosition(this_ *JSDocInfo_TrimmedStringPosition) {
	InitJSDocInfo_StringPosition(&this_.JSDocInfo_StringPositionBase)
	this_.initFields_JSDocInfo_TrimmedStringPosition()
}

// Java method: TrimmedStringPosition.setItem(java.lang.String item)
func (this_ *JSDocInfo_TrimmedStringPosition) SetItem(arg_87 any) {
	item := j2gort.CastString(arg_87)
	_ = item
	if !((int32(j2gort.JStrCharAt(item, int32(0))) != int32(uint16(32))) && (int32(j2gort.JStrCharAt(item, (j2gort.JStrLen(item)-int32(1)))) != int32(uint16(32)))) {
		j2gort.CheckArgumentStr(false, "String has leading or trailing whitespace")
	}
	this_.JSDocInfo_StringPositionBase.SetItem(j2gort.BoxString(item))
}

// Java class: com.google.javascript.rhino.JSDocInfo.NamePosition
type JSDocInfo_NamePosition struct {
	JSDocInfo_ComparableSourcePositionBase
}

func (this_ *JSDocInfo_NamePosition) BindSelf_(s any) {
	_ = s
}
func (this_ *JSDocInfo_NamePosition) IsJavaNil_() bool { return this_ == nil }
func (*JSDocInfo_NamePosition) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.JSDocInfo$NamePosition", "NamePosition"
}
func (this_ *JSDocInfo_NamePosition) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.JSDocInfo_ComparableSourcePositionBase.IdentityHash_()
}
func (*JSDocInfo_NamePosition) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.JSDocInfo.NamePosition", "com.google.javascript.rhino.JSDocInfo$NamePosition", func(v_ any) bool {
		return j2gort.InstanceOf[*JSDocInfo_NamePosition](v_)
	})
}
func (*JSDocInfo_NamePosition) Implements_rhino_JSDocInfo_ComparableSourcePosition() {}
func (*JSDocInfo_NamePosition) Implements_rhino_SourcePosition()                     {}

func (this_ *JSDocInfo_NamePosition) initFields_JSDocInfo_NamePosition() {
	_ = this_
}

// Java constructor: NamePosition()
func NewJSDocInfo_NamePosition() *JSDocInfo_NamePosition {
	this_ := &JSDocInfo_NamePosition{}
	InitJSDocInfo_NamePosition(this_)
	return this_
}

func InitJSDocInfo_NamePosition(this_ *JSDocInfo_NamePosition) {
	InitJSDocInfo_ComparableSourcePosition(&this_.JSDocInfo_ComparableSourcePositionBase)
	this_.initFields_JSDocInfo_NamePosition()
}

// Java method: NamePosition.isEquivalentTo(com.google.javascript.rhino.SourcePosition<com.google.javascript.rhino.Node> that)
func (this_ *JSDocInfo_NamePosition) IsEquivalentTo(that SourcePosition) bool {
	_ = that
	if ((that == nil || that.IsJavaNil_()) || !(this_.IsSamePositionAs(that))) || ((j2gort.CastNode[Node](this_.GetItem()) == 0) != (j2gort.CastNode[Node](that.GetItem()) == 0)) {
		return false
	}
	return ((j2gort.CastNode[Node](this_.GetItem()) == 0) || j2gort.CastNode[Node](this_.GetItem()).IsEquivalentTo_1(j2gort.CastNode[Node](that.GetItem())))
}

// Java class: com.google.javascript.rhino.JSDocInfo.TypePosition
type JSDocInfo_TypePosition struct {
	JSDocInfo_ComparableSourcePositionBase
	brackets bool `gson:"brackets"`
}

func (this_ *JSDocInfo_TypePosition) BindSelf_(s any) {
	_ = s
}
func (this_ *JSDocInfo_TypePosition) IsJavaNil_() bool { return this_ == nil }
func (*JSDocInfo_TypePosition) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.JSDocInfo$TypePosition", "TypePosition"
}
func (this_ *JSDocInfo_TypePosition) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.JSDocInfo_ComparableSourcePositionBase.IdentityHash_()
}
func (*JSDocInfo_TypePosition) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.JSDocInfo.TypePosition", "com.google.javascript.rhino.JSDocInfo$TypePosition", func(v_ any) bool {
		return j2gort.InstanceOf[*JSDocInfo_TypePosition](v_)
	})
}
func (*JSDocInfo_TypePosition) Implements_rhino_JSDocInfo_ComparableSourcePosition() {}
func (*JSDocInfo_TypePosition) Implements_rhino_SourcePosition()                     {}

func (this_ *JSDocInfo_TypePosition) initFields_JSDocInfo_TypePosition() {
	_ = this_

}

// Java constructor: TypePosition()
func NewJSDocInfo_TypePosition() *JSDocInfo_TypePosition {
	this_ := &JSDocInfo_TypePosition{}
	InitJSDocInfo_TypePosition(this_)
	return this_
}

func InitJSDocInfo_TypePosition(this_ *JSDocInfo_TypePosition) {
	InitJSDocInfo_ComparableSourcePosition(&this_.JSDocInfo_ComparableSourcePositionBase)
	this_.initFields_JSDocInfo_TypePosition()
}

// Java method: TypePosition.hasBrackets()
func (this_ *JSDocInfo_TypePosition) HasBrackets() bool {
	return this_.brackets
}

// Java method: TypePosition.setHasBrackets(boolean newVal)
func (this_ *JSDocInfo_TypePosition) SetHasBrackets(newVal bool) {
	_ = newVal
	this_.brackets = newVal
}

// Java method: TypePosition.isEquivalentTo(com.google.javascript.rhino.SourcePosition<com.google.javascript.rhino.Node> that)
func (this_ *JSDocInfo_TypePosition) IsEquivalentTo(that SourcePosition) bool {
	_ = that
	var typePosition *JSDocInfo_TypePosition
	_ = typePosition
	var ok_89 bool
	typePosition, ok_89 = j2gort.TryCast[*JSDocInfo_TypePosition](that)
	ok_88 := !ok_89
	if ok_88 || (!(this_.IsSamePositionAs(that))) || (this_.brackets != typePosition.brackets) || ((j2gort.CastNode[Node](this_.GetItem()) == 0) != (j2gort.CastNode[Node](that.GetItem()) == 0)) {
		return false
	}
	return ((j2gort.CastNode[Node](this_.GetItem()) == 0) || j2gort.CastNode[Node](this_.GetItem()).IsEquivalentTo_1(j2gort.CastNode[Node](that.GetItem())))
}

// Java class: com.google.javascript.rhino.JSDocInfo.Marker
type JSDocInfo_Marker struct {
	annotation  *JSDocInfo_TrimmedStringPosition `gson:"annotation"`
	nameNode    *JSDocInfo_NamePosition          `gson:"nameNode"`
	description JSDocInfo_StringPosition         `gson:"description"`
	type_       *JSDocInfo_TypePosition          `gson:"type"`
	idHash_     atomic.Int32
}

func (this_ *JSDocInfo_Marker) BindSelf_(s any) {
	_ = s
}
func (this_ *JSDocInfo_Marker) IsJavaNil_() bool { return this_ == nil }
func (*JSDocInfo_Marker) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.JSDocInfo$Marker", "Marker"
}
func (this_ *JSDocInfo_Marker) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*JSDocInfo_Marker) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.JSDocInfo.Marker", "com.google.javascript.rhino.JSDocInfo$Marker", func(v_ any) bool {
		return j2gort.InstanceOf[*JSDocInfo_Marker](v_)
	})
}

func (this_ *JSDocInfo_Marker) initFields_JSDocInfo_Marker() {
	_ = this_
}

// Java constructor: Marker()
func NewJSDocInfo_Marker() *JSDocInfo_Marker {
	this_ := &JSDocInfo_Marker{}
	InitJSDocInfo_Marker(this_)
	return this_
}

func InitJSDocInfo_Marker(this_ *JSDocInfo_Marker) {
	this_.initFields_JSDocInfo_Marker()
}

// Java method: Marker.getAnnotation()
func (this_ *JSDocInfo_Marker) GetAnnotation() JSDocInfo_StringPosition {
	return ToIface_JSDocInfo_StringPosition(this_.annotation)
}

// Java method: Marker.setAnnotation(com.google.javascript.rhino.JSDocInfo.TrimmedStringPosition p)
func (this_ *JSDocInfo_Marker) SetAnnotation(p *JSDocInfo_TrimmedStringPosition) {
	_ = p
	this_.annotation = p
}

// Java method: Marker.getNameNode()
func (this_ *JSDocInfo_Marker) GetNameNode() *JSDocInfo_NamePosition {
	return this_.nameNode
}

// Java method: Marker.setNameNode(com.google.javascript.rhino.JSDocInfo.NamePosition p)
func (this_ *JSDocInfo_Marker) SetNameNode(p *JSDocInfo_NamePosition) {
	_ = p
	this_.nameNode = p
}

// Java method: Marker.getDescription()
func (this_ *JSDocInfo_Marker) GetDescription() JSDocInfo_StringPosition {
	return this_.description
}

// Java method: Marker.setDescription(com.google.javascript.rhino.JSDocInfo.StringPosition p)
func (this_ *JSDocInfo_Marker) SetDescription(p JSDocInfo_StringPosition) {
	_ = p
	this_.description = p
}

// Java method: Marker.getType()
func (this_ *JSDocInfo_Marker) GetType() *JSDocInfo_TypePosition {
	return this_.type_
}

// Java method: Marker.setType(com.google.javascript.rhino.JSDocInfo.TypePosition p)
func (this_ *JSDocInfo_Marker) SetType(p *JSDocInfo_TypePosition) {
	_ = p
	this_.type_ = p
}

// Java method: Marker.isEquivalentTo(com.google.javascript.rhino.JSDocInfo.Marker that)
func (this_ *JSDocInfo_Marker) priv_Marker_isEquivalentTo_8169(that *JSDocInfo_Marker) bool {
	_ = that
	return (((JSDocInfo_Marker_AreEquivalent(ToIface_JSDocInfo_ComparableSourcePosition(this_.annotation), ToIface_JSDocInfo_ComparableSourcePosition(that.annotation)) && JSDocInfo_Marker_AreEquivalent(ToIface_JSDocInfo_ComparableSourcePosition(this_.nameNode), ToIface_JSDocInfo_ComparableSourcePosition(that.nameNode))) && JSDocInfo_Marker_AreEquivalent(j2gort.AsIface[JSDocInfo_ComparableSourcePosition](this_.description), j2gort.AsIface[JSDocInfo_ComparableSourcePosition](that.description))) && JSDocInfo_Marker_AreEquivalent(ToIface_JSDocInfo_ComparableSourcePosition(this_.type_), ToIface_JSDocInfo_ComparableSourcePosition(that.type_)))
}

// Java static method: Marker.areEquivalent(com.google.javascript.rhino.JSDocInfo.@org.jspecify.annotations.Nullable ComparableSourcePosition<T> p1, com.google.javascript.rhino.JSDocInfo.@org.jspecify.annotations.Nullable ComparableSourcePosition<T> p2)
func JSDocInfo_Marker_AreEquivalent(p1 JSDocInfo_ComparableSourcePosition, p2 JSDocInfo_ComparableSourcePosition) bool {
	_ = p1
	_ = p2
	return (((p1 == nil || p1.IsJavaNil_()) == (p2 == nil || p2.IsJavaNil_())) && ((p1 == nil || p1.IsJavaNil_()) || p1.IsEquivalentTo(j2gort.AsIface[SourcePosition](p2))))
}

// Java class: com.google.javascript.rhino.JSDocInfo.Builder
type JSDocInfo_Builder struct {
	propVal0       any                        `gson:"propVal0"`
	propVal1       any                        `gson:"propVal1"`
	propVal2       any                        `gson:"propVal2"`
	propVal3       any                        `gson:"propVal3"`
	propValues     []any                      `gson:"propValues"`
	typeVal0       JSTypeExpression           `gson:"typeVal0"`
	typeVal1       JSTypeExpression           `gson:"typeVal1"`
	typeVal2       JSTypeExpression           `gson:"typeVal2"`
	typeVal3       JSTypeExpression           `gson:"typeVal3"`
	currentMarker  *JSDocInfo_Marker          `gson:"currentMarker"`
	licenseTexts   *j2gort.OrderedSet[string] `gson:"licenseTexts"`
	propKeys       int64                      `gson:"propKeys"`
	singleTypeKeys int64                      `gson:"singleTypeKeys"`
	bits           int64                      `gson:"bits"`
	idHash_        atomic.Int32
	populated      bool `gson:"populated"`
}

func (this_ *JSDocInfo_Builder) BindSelf_(s any) {
	_ = s
}
func (this_ *JSDocInfo_Builder) IsJavaNil_() bool { return this_ == nil }
func (*JSDocInfo_Builder) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.JSDocInfo$Builder", "Builder"
}
func (this_ *JSDocInfo_Builder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*JSDocInfo_Builder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.JSDocInfo.Builder", "com.google.javascript.rhino.JSDocInfo$Builder", func(v_ any) bool {
		return j2gort.InstanceOf[*JSDocInfo_Builder](v_)
	})
}

func (this_ *JSDocInfo_Builder) initFields_JSDocInfo_Builder() {
	_ = this_

}

// Java constructor: Builder()
func NewJSDocInfo_Builder() *JSDocInfo_Builder {
	this_ := &JSDocInfo_Builder{}
	InitJSDocInfo_Builder(this_)
	return this_
}

func InitJSDocInfo_Builder(this_ *JSDocInfo_Builder) {
	this_.initFields_JSDocInfo_Builder()
}

// Java static method: Builder.copyFrom(com.google.javascript.rhino.JSDocInfo info)
func JSDocInfo_Builder_CopyFrom(info *JSDocInfo) *JSDocInfo_Builder {
	_ = info
	return info.ToBuilder_0()
}

// Java static method: Builder.maybeCopyFrom(com.google.javascript.rhino.@org.jspecify.annotations.Nullable JSDocInfo info)
func JSDocInfo_Builder_MaybeCopyFrom(info *JSDocInfo) *JSDocInfo_Builder {
	_ = info
	if info != nil {
		return info.ToBuilder_0()
	} else {
		return JSDocInfo_S_Builder().ParseDocumentation()
	}
}

// Java static method: Builder.maybeCopyFromWithNewType(com.google.javascript.rhino.JSDocInfo info, com.google.javascript.rhino.JSTypeExpression typeExpression)
func JSDocInfo_Builder_MaybeCopyFromWithNewType(info *JSDocInfo, typeExpression JSTypeExpression) *JSDocInfo_Builder {
	_ = info
	_ = typeExpression
	if info == nil {
		return JSDocInfo_S_Builder().ParseDocumentation().SetType(typeExpression)
	}
	return info.ToBuilder_0().SetType(typeExpression)
}

// Java static method: Builder.copyFromWithNewType(com.google.javascript.rhino.JSDocInfo info, com.google.javascript.rhino.JSTypeExpression typeExpression)
func JSDocInfo_Builder_CopyFromWithNewType(info *JSDocInfo, typeExpression JSTypeExpression) *JSDocInfo_Builder {
	_ = info
	_ = typeExpression
	return info.ToBuilder_0().SetType(typeExpression)
}

// Java static method: Builder.maybeCopyFromAndReplaceNames(com.google.javascript.rhino.JSDocInfo info, java.util.Set<java.lang.String> moduleLocalNamesToReplace)
func JSDocInfo_Builder_MaybeCopyFromAndReplaceNames(info *JSDocInfo, moduleLocalNamesToReplace j2gort.Ext_java_util_Set) *JSDocInfo_Builder {
	_ = info
	_ = moduleLocalNamesToReplace
	if info != nil {
		return JSDocInfo_Builder_CopyFromAndReplaceNames(info, moduleLocalNamesToReplace)
	} else {
		return JSDocInfo_S_Builder().ParseDocumentation()
	}
}

// Java static method: Builder.copyFromAndReplaceNames(com.google.javascript.rhino.JSDocInfo info, java.util.Set<java.lang.String> oldNames)
func JSDocInfo_Builder_CopyFromAndReplaceNames(info *JSDocInfo, oldNames j2gort.Ext_java_util_Set) *JSDocInfo_Builder {
	_ = info
	_ = oldNames
	return info.CloneAndReplaceTypeNames(oldNames).ToBuilder_0()
}

// Java method: Builder.parseDocumentation()
func (this_ *JSDocInfo_Builder) ParseDocumentation() *JSDocInfo_Builder {
	this_.priv_Builder_setBit_8323(StaticGet_JSDocInfo_Bit_INCLUDE_DOCUMENTATION(), true)
	return this_
}

// Java method: Builder.shouldParseDocumentation()
func (this_ *JSDocInfo_Builder) ShouldParseDocumentation() bool {
	return this_.priv_Builder_checkBit_8322(StaticGet_JSDocInfo_Bit_INCLUDE_DOCUMENTATION())
}

// Java method: Builder.recordOriginalCommentString(java.lang.String sourceComment)
func (this_ *JSDocInfo_Builder) RecordOriginalCommentString(sourceComment string) {
	_ = sourceComment
	if this_.ShouldParseDocumentation() {
		this_.populated = true
		this_.priv_Builder_setProp_8314(StaticGet_JSDocInfo_SOURCE_COMMENT(), j2gort.BoxString(sourceComment))
	}
}

// Java method: Builder.recordOriginalCommentPosition(int position)
func (this_ *JSDocInfo_Builder) RecordOriginalCommentPosition(position int32) {
	_ = position
	if this_.ShouldParseDocumentation() {
		this_.populated = true
		this_.priv_Builder_setProp_8314(StaticGet_JSDocInfo_ORIGINAL_COMMENT_POSITION(), j2gort.ValueOfInt(position))
	}
}

// Java method: Builder.isPopulatedWithFileOverview()
func (this_ *JSDocInfo_Builder) IsPopulatedWithFileOverview() bool {
	return (this_.populated && (((this_.bits & (((((StaticGet_JSDocInfo_Bit_FILEOVERVIEW().mask | StaticGet_JSDocInfo_Bit_EXTERNS().mask) | StaticGet_JSDocInfo_Bit_NOCOMPILE().mask) | StaticGet_JSDocInfo_Bit_NOCOVERAGE().mask) | StaticGet_JSDocInfo_Bit_TYPE_SUMMARY().mask) | StaticGet_JSDocInfo_Bit_ENHANCED_NAMESPACE().mask)) != int64(int32(0))) || this_.IsModsRecorded()))
}

// Java method: Builder.isDescriptionRecorded()
func (this_ *JSDocInfo_Builder) IsDescriptionRecorded() bool {
	return (j2gort.CastString(this_.priv_Builder_getProp_8315(StaticGet_JSDocInfo_DESCRIPTION())) != j2gort.NullString)
}

// Java method: Builder.isNoInline()
func (this_ *JSDocInfo_Builder) IsNoInline() bool {
	return this_.priv_Builder_checkBit_8322(StaticGet_JSDocInfo_Bit_NOINLINE())
}

// Java method: Builder.isRequireInlining()
func (this_ *JSDocInfo_Builder) IsRequireInlining() bool {
	return this_.priv_Builder_checkBit_8322(StaticGet_JSDocInfo_Bit_REQUIRE_INLINING())
}

// Java method: Builder.isEncourageInlining()
func (this_ *JSDocInfo_Builder) IsEncourageInlining() bool {
	return this_.priv_Builder_checkBit_8322(StaticGet_JSDocInfo_Bit_ENCOURAGE_INLINING())
}

// Java method: Builder.build()
func (this_ *JSDocInfo_Builder) Build_0() *JSDocInfo {
	return this_.Build_1(false)
}

// Java method: Builder.build(boolean always)
func (this_ *JSDocInfo_Builder) Build_1(always bool) *JSDocInfo {
	_ = always
	if this_.populated || always {
		clinit_JSDocInfo()
		var val_90 JSDocInfo
		InitJSDocInfo(&val_90, this_)
		var info *JSDocInfo = &val_90
		_ = info
		this_.populated = false
		return info
	}
	return nil
}

// Java method: Builder.buildAndReset()
func (this_ *JSDocInfo_Builder) BuildAndReset() *JSDocInfo {
	var info *JSDocInfo = this_.Build_0()
	_ = info
	var keepDocBit int64 = (this_.bits & StaticGet_JSDocInfo_Bit_INCLUDE_DOCUMENTATION().mask)
	_ = keepDocBit
	this_.Reset()
	this_.bits = keepDocBit
	return info
}

// Java method: Builder.reset()
func (this_ *JSDocInfo_Builder) Reset() {
	this_.bits = int64(0)
	if this_.propValues != nil {
		var rem int64 = this_.propKeys
		_ = rem
		for rem != int64(0) {
			var bit int32 = j2gort.BoxLong_NumberOfTrailingZeros(rem)
			_ = bit
			rem = (rem & (rem - int64(1)))
			{
				a_91 := this_.propValues
				i_92 := bit
				a_91[i_92] = nil
			}
		}
	} else if this_.propKeys != int64(0) {
		this_.propVal0 = nil
		this_.propVal1 = nil
		this_.propVal2 = nil
		this_.propVal3 = nil
	}
	this_.propKeys = int64(0)
	if this_.singleTypeKeys != int64(0) {
		this_.typeVal0 = (JSTypeExpression{})
		this_.typeVal1 = (JSTypeExpression{})
		this_.typeVal2 = (JSTypeExpression{})
		this_.typeVal3 = (JSTypeExpression{})
		this_.singleTypeKeys = int64(0)
	}
	this_.populated = false
	this_.currentMarker = nil
	if this_.licenseTexts != nil {
		this_.licenseTexts.Clear()
	}
}

// Java method: Builder.markAnnotation(java.lang.String annotation, int lineno, int charno)
func (this_ *JSDocInfo_Builder) MarkAnnotation(annotation string, lineno int32, charno int32) {
	_ = annotation
	_ = lineno
	_ = charno
	var marker *JSDocInfo_Marker = this_.priv_Builder_addMarker_8184()
	_ = marker
	if marker != nil {
		var position *JSDocInfo_TrimmedStringPosition = NewJSDocInfo_TrimmedStringPosition()
		_ = position
		position.SetItem(j2gort.BoxString(annotation))
		position.SetPositionInformation(lineno, charno, lineno, (charno + j2gort.JStrLen(annotation)))
		marker.SetAnnotation(position)
		this_.populated = true
	}
	this_.currentMarker = marker
}

// Java method: Builder.addMarker()
func (this_ *JSDocInfo_Builder) priv_Builder_addMarker_8184() *JSDocInfo_Marker {
	if this_.ShouldParseDocumentation() {
		var markers j2gort.Ext_java_util_ArrayList = j2gort.AsIface[j2gort.Ext_java_util_ArrayList](this_.priv_Builder_getProp_8315(StaticGet_JSDocInfo_MARKERS()))
		_ = markers
		if markers == nil {
			this_.priv_Builder_setProp_8314(StaticGet_JSDocInfo_MARKERS(), func() j2gort.Ext_java_util_ArrayList {
				var asgn_93 j2gort.Ext_java_util_ArrayList = j2gort.NewExt_java_util_ArrayList()
				markers = asgn_93
				return asgn_93
			}())
		}
		var val_94 JSDocInfo_Marker
		InitJSDocInfo_Marker(&val_94)
		var marker *JSDocInfo_Marker = &val_94
		_ = marker
		markers.Add_1(marker)
		return marker
	}
	return nil
}

// Java method: Builder.markText(java.lang.String text, int startLineno, int startCharno, int endLineno, int endCharno)
func (this_ *JSDocInfo_Builder) MarkText(text string, startLineno int32, startCharno int32, endLineno int32, endCharno int32) {
	_ = text
	_ = startLineno
	_ = startCharno
	_ = endLineno
	_ = endCharno
	if this_.currentMarker != nil {
		var position JSDocInfo_StringPosition = NewJSDocInfo_StringPosition()
		_ = position
		position.SetItem(j2gort.BoxString(text))
		position.SetPositionInformation(startLineno, startCharno, endLineno, endCharno)
		this_.currentMarker.SetDescription(position)
	}
}

// Java method: Builder.markTypeNode(com.google.javascript.rhino.Node typeNode, int lineno, int startCharno, int endLineno, int endCharno, boolean hasLC)
func (this_ *JSDocInfo_Builder) MarkTypeNode(typeNode Node, lineno int32, startCharno int32, endLineno int32, endCharno int32, hasLC bool) {
	_ = typeNode
	_ = lineno
	_ = startCharno
	_ = endLineno
	_ = endCharno
	_ = hasLC
	if this_.currentMarker != nil {
		var position *JSDocInfo_TypePosition = NewJSDocInfo_TypePosition()
		_ = position
		position.SetItem(j2gort.BoxNode(typeNode))
		position.SetHasBrackets(hasLC)
		position.SetPositionInformation(lineno, startCharno, endLineno, endCharno)
		this_.currentMarker.SetType(position)
	}
}

// Java method: Builder.markName(java.lang.String name, com.google.javascript.rhino.Node templateNode, int lineno, int charno)
func (this_ *JSDocInfo_Builder) MarkName(name string, templateNode Node, lineno int32, charno int32) {
	_ = name
	_ = templateNode
	_ = lineno
	_ = charno
	if this_.currentMarker != nil {
		var position *JSDocInfo_TrimmedStringPosition = NewJSDocInfo_TrimmedStringPosition()
		_ = position
		position.SetItem(j2gort.BoxString(name))
		position.SetPositionInformation(lineno, charno, lineno, (charno + j2gort.JStrLen(name)))
		var nodePos *JSDocInfo_NamePosition = NewJSDocInfo_NamePosition()
		_ = nodePos
		var node Node = Node_NewString_2(StaticGet_Token_NAME(), name).SetLinenoCharno(lineno, charno)
		_ = node
		node.SetLength(j2gort.JStrLen(name))
		if templateNode != 0 {
			node.SetStaticSourceFileFrom(templateNode)
		}
		nodePos.SetItem(j2gort.BoxNode(node))
		nodePos.SetPositionInformation(lineno, charno, lineno, (charno + j2gort.JStrLen(name)))
		this_.currentMarker.SetNameNode(nodePos)
	}
}

// Java method: Builder.recordBlockDescription(java.lang.String description)
func (this_ *JSDocInfo_Builder) RecordBlockDescription(description string) bool {
	_ = description
	this_.populated = true
	if !(this_.ShouldParseDocumentation()) {
		return true
	}
	return this_.priv_Builder_populateProp_8321(StaticGet_JSDocInfo_BLOCK_DESCRIPTION(), j2gort.BoxString(description))
}

// Java method: Builder.recordVisibility(com.google.javascript.rhino.JSDocInfo.Visibility visibility)
func (this_ *JSDocInfo_Builder) RecordVisibility(visibility JSDocInfo_Visibility) bool {
	_ = visibility
	if j2gort.Cast[JSDocInfo_Visibility](this_.priv_Builder_getProp_8315(StaticGet_JSDocInfo_VISIBILITY())) == 0 {
		this_.populated = true
		this_.priv_Builder_setProp_8314(StaticGet_JSDocInfo_VISIBILITY(), j2gort.BoxNullableVal(visibility))
		return true
	}
	return false
}

// Java method: Builder.overwriteVisibility(com.google.javascript.rhino.JSDocInfo.Visibility visibility)
func (this_ *JSDocInfo_Builder) OverwriteVisibility(visibility JSDocInfo_Visibility) {
	_ = visibility
	this_.populated = true
	this_.priv_Builder_setProp_8314(StaticGet_JSDocInfo_VISIBILITY(), j2gort.BoxNullableVal(visibility))
}

// Java method: Builder.recordParameter(java.lang.String parameterName, com.google.javascript.rhino.JSTypeExpression type)
func (this_ *JSDocInfo_Builder) RecordParameter(parameterName string, type_ JSTypeExpression) bool {
	_ = parameterName
	_ = type_
	return (!(this_.priv_Builder_hasAnySingletonTypeTags_8304()) && this_.priv_Builder_populateTypePropEntry_8318(StaticGet_JSDocInfo_PARAMETERS(), RhinoStringPool_AddOrGet(parameterName), type_))
}

// Java method: Builder.recordParameterDescription(java.lang.String parameterName, java.lang.String description)
func (this_ *JSDocInfo_Builder) RecordParameterDescription(parameterName string, description string) bool {
	_ = parameterName
	_ = description
	if !(this_.ShouldParseDocumentation()) {
		return true
	}
	return this_.priv_Builder_populateStringPropEntry_8319(StaticGet_JSDocInfo_PARAMETER_DESCRIPTIONS(), RhinoStringPool_AddOrGet(parameterName), description)
}

// Java method: Builder.recordTemplateTypeName(java.lang.String name)
func (this_ *JSDocInfo_Builder) RecordTemplateTypeName_1(name string) bool {
	_ = name
	return this_.RecordTemplateTypeName_2(name, (JSTypeExpression{}))
}

// Java method: Builder.recordTemplateTypeName(java.lang.String name, com.google.javascript.rhino.@org.jspecify.annotations.Nullable JSTypeExpression bound)
func (this_ *JSDocInfo_Builder) RecordTemplateTypeName_2(name string, bound JSTypeExpression) bool {
	_ = name
	_ = bound
	if bound == (JSTypeExpression{}) {
		bound = StaticGet_JSTypeExpression_IMPLICIT_TEMPLATE_BOUND()
	}
	var transformations *JSDocInfo_NodeMap = j2gort.CastPtr[JSDocInfo_NodeMap](this_.priv_Builder_getProp_8315(ToIface_JSDocInfo_Property(StaticGet_JSDocInfo_TYPE_TRANSFORMATIONS())))
	_ = transformations
	if ((transformations != nil) && transformations.ContainsKey(name)) || ((this_.singleTypeKeys & StaticGet_JSDocInfo_TYPEDEF_TYPE().GetJSDocInfo_PropertyBase().mask) != int64(0)) {
		return false
	}
	return this_.priv_Builder_populateTypePropEntry_8318(StaticGet_JSDocInfo_TEMPLATE_TYPE_NAMES(), RhinoStringPool_AddOrGet(name), bound)
}

// Java method: Builder.recordTypeTransformation(java.lang.String name, com.google.javascript.rhino.Node expr)
func (this_ *JSDocInfo_Builder) RecordTypeTransformation(name string, expr Node) bool {
	_ = name
	_ = expr
	var names *JSDocInfo_TypeMap = j2gort.CastPtr[JSDocInfo_TypeMap](this_.priv_Builder_getProp_8315(ToIface_JSDocInfo_Property(StaticGet_JSDocInfo_TEMPLATE_TYPE_NAMES())))
	_ = names
	if (names != nil) && names.ContainsKey(name) {
		return false
	}
	return this_.priv_Builder_populateNodePropEntry_8320(StaticGet_JSDocInfo_TYPE_TRANSFORMATIONS(), RhinoStringPool_AddOrGet(name), expr)
}

// Java method: Builder.recordThrowsAnnotation(java.lang.String annotation)
func (this_ *JSDocInfo_Builder) RecordThrowsAnnotation(annotation string) bool {
	_ = annotation
	this_.populated = true
	if this_.priv_Builder_checkBit_8322(StaticGet_JSDocInfo_Bit_NOSIDEEFFECTS()) {
		return false
	}
	if !(this_.priv_Builder_hasAnySingletonTypeTags_8304()) {
		var throwsAnnotations j2gort.Ext_java_util_List = this_.priv_Builder_getOrCreateStringList_8317(StaticGet_JSDocInfo_THROWS_ANNOTATIONS())
		_ = throwsAnnotations
		if this_.ShouldParseDocumentation() {
			throwsAnnotations.Add_1(j2gort.BoxString(annotation))
		} else if throwsAnnotations.IsEmpty() {
			throwsAnnotations.Add_1("")
		}
	}
	return true
}

// Java method: Builder.recordAuthor(java.lang.String author)
func (this_ *JSDocInfo_Builder) RecordAuthor(author string) bool {
	_ = author
	this_.populated = true
	if this_.ShouldParseDocumentation() {
		this_.priv_Builder_getOrCreateStringList_8317(StaticGet_JSDocInfo_AUTHORS()).Add_1(j2gort.BoxString(author))
	}
	return true
}

// Java method: Builder.recordReference(java.lang.String reference)
func (this_ *JSDocInfo_Builder) RecordReference(reference string) bool {
	_ = reference
	this_.populated = true
	if this_.ShouldParseDocumentation() {
		this_.priv_Builder_getOrCreateStringList_8317(StaticGet_JSDocInfo_SEES()).Add_1(j2gort.BoxString(reference))
	}
	return true
}

// Java method: Builder.recordConsistentIdGenerator()
func (this_ *JSDocInfo_Builder) RecordConsistentIdGenerator() bool {
	return this_.priv_Builder_populateProp_8321(StaticGet_JSDocInfo_ID_GENERATOR(), j2gort.BoxNullableVal(StaticGet_JSDocInfo_IdGenerator_CONSISTENT()))
}

// Java method: Builder.recordStableIdGenerator()
func (this_ *JSDocInfo_Builder) RecordStableIdGenerator() bool {
	return this_.priv_Builder_populateProp_8321(StaticGet_JSDocInfo_ID_GENERATOR(), j2gort.BoxNullableVal(StaticGet_JSDocInfo_IdGenerator_STABLE()))
}

// Java method: Builder.recordXidGenerator()
func (this_ *JSDocInfo_Builder) RecordXidGenerator() bool {
	return this_.priv_Builder_populateProp_8321(StaticGet_JSDocInfo_ID_GENERATOR(), j2gort.BoxNullableVal(StaticGet_JSDocInfo_IdGenerator_XID()))
}

// Java method: Builder.recordMappedIdGenerator()
func (this_ *JSDocInfo_Builder) RecordMappedIdGenerator() bool {
	return this_.priv_Builder_populateProp_8321(StaticGet_JSDocInfo_ID_GENERATOR(), j2gort.BoxNullableVal(StaticGet_JSDocInfo_IdGenerator_MAPPED()))
}

// Java method: Builder.recordIdGenerator()
func (this_ *JSDocInfo_Builder) RecordIdGenerator() bool {
	return this_.priv_Builder_populateProp_8321(StaticGet_JSDocInfo_ID_GENERATOR(), j2gort.BoxNullableVal(StaticGet_JSDocInfo_IdGenerator_UNIQUE()))
}

// Java method: Builder.recordDeprecationReason(java.lang.String reason)
func (this_ *JSDocInfo_Builder) RecordDeprecationReason(reason string) bool {
	_ = reason
	return this_.priv_Builder_populateProp_8321(StaticGet_JSDocInfo_DEPRECATION_REASON(), j2gort.BoxString(reason))
}

// Java method: Builder.isDeprecationReasonRecorded()
func (this_ *JSDocInfo_Builder) IsDeprecationReasonRecorded() bool {
	return (j2gort.CastString(this_.priv_Builder_getProp_8315(StaticGet_JSDocInfo_DEPRECATION_REASON())) != j2gort.NullString)
}

// Java method: Builder.recordSuppressions(com.google.common.collect.ImmutableSet<java.lang.String> suppressions, java.lang.String description)
func (this_ *JSDocInfo_Builder) RecordSuppressions_2(suppressions j2gort.Ext_com_google_common_collect_ImmutableSet, description string) {
	_ = suppressions
	_ = description
	this_.populated = true
	var current j2gort.Ext_com_google_common_collect_ImmutableMap = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap](this_.priv_Builder_getProp_8315(StaticGet_JSDocInfo_SUPPRESSIONS()))
	_ = current
	if current == nil {
		this_.priv_Builder_setProp_8314(StaticGet_JSDocInfo_SUPPRESSIONS(), j2gort.Ext_com_google_common_collect_ImmutableMap_Of_2(JSDocInfo_InternSet(j2gort.AsIface[j2gort.Ext_java_util_Set](suppressions)), j2gort.BoxString(description)))
		return
	}
	if current.ContainsKey(suppressions) {
		return
	}
	var mapBuilder j2gort.Ext_com_google_common_collect_ImmutableMap_Builder = j2gort.Ext_com_google_common_collect_ImmutableMap_BuilderWithExpectedSize((current.Size() + int32(1)))
	_ = mapBuilder
	mapBuilder.PutAll_1_Ext_java_util_Map(j2gort.AsIface[j2gort.Ext_java_util_Map](current))
	mapBuilder.Put_2(JSDocInfo_InternSet(j2gort.AsIface[j2gort.Ext_java_util_Set](suppressions)), j2gort.BoxString(description))
	this_.priv_Builder_setProp_8314(StaticGet_JSDocInfo_SUPPRESSIONS(), j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap](mapBuilder.BuildOrThrow()))
}

// Java method: Builder.recordSuppressions(java.util.Set<java.lang.String> suppressions)
func (this_ *JSDocInfo_Builder) RecordSuppressions_1(suppressions j2gort.Ext_java_util_Set) {
	_ = suppressions
	this_.RecordSuppressions_2(j2gort.Ext_com_google_common_collect_ImmutableSet_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](suppressions)), "")
}

// Java method: Builder.recordSuppression(java.lang.String suppression)
func (this_ *JSDocInfo_Builder) RecordSuppression(suppression string) {
	_ = suppression
	this_.RecordSuppressions_2(j2gort.ImmSetOfStr1(suppression), "")
}

// Java method: Builder.recordModifies(java.util.Set<java.lang.String> modifies)
func (this_ *JSDocInfo_Builder) RecordModifies(modifies j2gort.Ext_java_util_Set) bool {
	_ = modifies
	return (!(this_.priv_Builder_hasAnySingletonSideEffectTags_8310()) && this_.priv_Builder_populateProp_8321(StaticGet_JSDocInfo_MODIFIES(), JSDocInfo_InternSet(modifies)))
}

// Java method: Builder.recordType(com.google.javascript.rhino.JSTypeExpression type)
func (this_ *JSDocInfo_Builder) RecordType(type_ JSTypeExpression) bool {
	_ = type_
	return (((type_ != (JSTypeExpression{})) && !(this_.priv_Builder_hasAnyTypeRelatedTags_8300())) && this_.priv_Builder_populateTypeProp_8309(StaticGet_JSDocInfo_TYPE(), type_))
}

// Java method: Builder.recordInlineType()
func (this_ *JSDocInfo_Builder) RecordInlineType() {
	this_.priv_Builder_populateBit_8324(StaticGet_JSDocInfo_Bit_INLINE_TYPE(), true)
}

// Java method: Builder.recordTypedef(com.google.javascript.rhino.JSTypeExpression type)
func (this_ *JSDocInfo_Builder) RecordTypedef(type_ JSTypeExpression) bool {
	_ = type_
	return ((((type_ != (JSTypeExpression{})) && !(this_.priv_Builder_hasAnyTypeRelatedTags_8300())) && (j2gort.CastPtr[JSDocInfo_TypeMap](this_.priv_Builder_getProp_8315(ToIface_JSDocInfo_Property(StaticGet_JSDocInfo_TEMPLATE_TYPE_NAMES()))) == nil)) && this_.priv_Builder_populateTypeProp_8309(StaticGet_JSDocInfo_TYPEDEF_TYPE(), type_))
}

// Java method: Builder.recordReturnType(com.google.javascript.rhino.JSTypeExpression type)
func (this_ *JSDocInfo_Builder) RecordReturnType(type_ JSTypeExpression) bool {
	_ = type_
	return (((type_ != (JSTypeExpression{})) && !(this_.priv_Builder_hasAnySingletonTypeTags_8304())) && this_.priv_Builder_populateTypeProp_8309(StaticGet_JSDocInfo_RETURN_TYPE(), type_))
}

// Java method: Builder.recordReturnDescription(java.lang.String description)
func (this_ *JSDocInfo_Builder) RecordReturnDescription(description string) bool {
	_ = description
	if !(this_.ShouldParseDocumentation()) {
		return true
	}
	return this_.priv_Builder_populateProp_8321(StaticGet_JSDocInfo_RETURN_DESCRIPTION(), j2gort.BoxString(description))
}

// Java method: Builder.recordDefineType(com.google.javascript.rhino.JSTypeExpression type)
func (this_ *JSDocInfo_Builder) RecordDefineType(type_ JSTypeExpression) bool {
	_ = type_
	if (((type_ != (JSTypeExpression{})) && !(this_.priv_Builder_checkBit_8322(StaticGet_JSDocInfo_Bit_CONST()))) && !(this_.priv_Builder_checkBit_8322(StaticGet_JSDocInfo_Bit_DEFINE()))) && this_.RecordType(type_) {
		return this_.priv_Builder_populateBit_8324(StaticGet_JSDocInfo_Bit_DEFINE(), true)
	}
	return false
}

// Java method: Builder.recordEnumParameterType(com.google.javascript.rhino.JSTypeExpression type)
func (this_ *JSDocInfo_Builder) RecordEnumParameterType(type_ JSTypeExpression) bool {
	_ = type_
	if (type_ != (JSTypeExpression{})) && !(this_.priv_Builder_hasAnyTypeRelatedTags_8300()) {
		this_.priv_Builder_setTypeProp_8307(StaticGet_JSDocInfo_ENUM_PARAMETER_TYPE(), type_)
		this_.populated = true
		return true
	}
	return false
}

// Java method: Builder.recordThisType(com.google.javascript.rhino.JSTypeExpression type)
func (this_ *JSDocInfo_Builder) RecordThisType(type_ JSTypeExpression) bool {
	_ = type_
	return (((type_ != (JSTypeExpression{})) && !(this_.priv_Builder_hasAnySingletonTypeTags_8304())) && this_.priv_Builder_populateTypeProp_8309(StaticGet_JSDocInfo_THIS_TYPE(), type_))
}

// Java method: Builder.recordBaseType(com.google.javascript.rhino.JSTypeExpression type)
func (this_ *JSDocInfo_Builder) RecordBaseType(type_ JSTypeExpression) bool {
	_ = type_
	return (((type_ != (JSTypeExpression{})) && !(this_.priv_Builder_hasAnySingletonTypeTags_8304())) && this_.priv_Builder_populateTypeProp_8309(StaticGet_JSDocInfo_BASE_TYPE(), type_))
}

// Java method: Builder.changeBaseType(com.google.javascript.rhino.JSTypeExpression type)
func (this_ *JSDocInfo_Builder) ChangeBaseType(type_ JSTypeExpression) bool {
	_ = type_
	if (type_ != (JSTypeExpression{})) && !(this_.priv_Builder_hasAnySingletonTypeTags_8304()) {
		this_.priv_Builder_setTypeProp_8307(StaticGet_JSDocInfo_BASE_TYPE(), type_)
		this_.populated = true
		return true
	}
	return false
}

// Java method: Builder.recordConstancy()
func (this_ *JSDocInfo_Builder) RecordConstancy() bool {
	return this_.priv_Builder_populateBit_8324(StaticGet_JSDocInfo_Bit_CONST(), true)
}

// Java method: Builder.recordMutable()
func (this_ *JSDocInfo_Builder) RecordMutable() bool {
	return this_.priv_Builder_populateBit_8324(StaticGet_JSDocInfo_Bit_CONST(), false)
}

// Java method: Builder.recordFinality()
func (this_ *JSDocInfo_Builder) RecordFinality() bool {
	return this_.priv_Builder_populateBit_8324(StaticGet_JSDocInfo_Bit_FINAL(), true)
}

// Java method: Builder.recordDescription(java.lang.String description)
func (this_ *JSDocInfo_Builder) RecordDescription(description string) bool {
	_ = description
	return this_.priv_Builder_populateProp_8321(StaticGet_JSDocInfo_DESCRIPTION(), j2gort.BoxString(description))
}

// Java method: Builder.recordTsType(java.lang.String tsType)
func (this_ *JSDocInfo_Builder) RecordTsType(tsType string) {
	_ = tsType
	this_.populated = true
	this_.priv_Builder_getOrCreateStringList_8317(StaticGet_JSDocInfo_TS_TYPES()).Add_1(j2gort.BoxString(tsType))
}

// Java method: Builder.recordMeaning(java.lang.String meaning)
func (this_ *JSDocInfo_Builder) RecordMeaning(meaning string) bool {
	_ = meaning
	return this_.priv_Builder_populateProp_8321(StaticGet_JSDocInfo_MEANING(), j2gort.BoxString(meaning))
}

// Java method: Builder.recordAlternateMessageId(java.lang.String alternateMessageId)
func (this_ *JSDocInfo_Builder) RecordAlternateMessageId(alternateMessageId string) bool {
	_ = alternateMessageId
	return this_.priv_Builder_populateProp_8321(StaticGet_JSDocInfo_ALTERNATE_MESSAGE_ID(), j2gort.BoxString(alternateMessageId))
}

// Java method: Builder.recordClosurePrimitiveId(java.lang.String closurePrimitiveId)
func (this_ *JSDocInfo_Builder) RecordClosurePrimitiveId(closurePrimitiveId string) bool {
	_ = closurePrimitiveId
	return this_.priv_Builder_populateProp_8321(StaticGet_JSDocInfo_CLOSURE_PRIMITIVE_ID(), j2gort.BoxString(RhinoStringPool_AddOrGet(closurePrimitiveId)))
}

// Java method: Builder.recordFileOverview(java.lang.String description)
func (this_ *JSDocInfo_Builder) RecordFileOverview(description string) bool {
	_ = description
	this_.priv_Builder_setBit_8323(StaticGet_JSDocInfo_Bit_FILEOVERVIEW(), true)
	this_.populated = true
	if !(this_.ShouldParseDocumentation()) {
		return true
	}
	return this_.priv_Builder_populateProp_8321(StaticGet_JSDocInfo_FILEOVERVIEW_DESCRIPTION(), j2gort.BoxString(description))
}

// Java method: Builder.recordEnhance(java.lang.String namespace)
func (this_ *JSDocInfo_Builder) RecordEnhance(namespace string) bool {
	_ = namespace
	this_.priv_Builder_setBit_8323(StaticGet_JSDocInfo_Bit_ENHANCED_NAMESPACE(), true)
	this_.populated = true
	return this_.priv_Builder_populateProp_8321(StaticGet_JSDocInfo_ENHANCED_NAMESPACE(), j2gort.BoxString(namespace))
}

// Java method: Builder.isModsRecorded()
func (this_ *JSDocInfo_Builder) IsModsRecorded() bool {
	return (j2gort.CastString(this_.priv_Builder_getProp_8315(StaticGet_JSDocInfo_MODS())) != j2gort.NullString)
}

// Java method: Builder.recordMods(java.lang.String namespace)
func (this_ *JSDocInfo_Builder) RecordMods(namespace string) bool {
	_ = namespace
	return this_.priv_Builder_populateProp_8321(StaticGet_JSDocInfo_MODS(), j2gort.BoxString(RhinoStringPool_AddOrGet(namespace)))
}

// Java method: Builder.recordLicense(java.lang.String license)
func (this_ *JSDocInfo_Builder) RecordLicense(license string) bool {
	_ = license
	this_.priv_Builder_setProp_8314(StaticGet_JSDocInfo_LICENSE(), j2gort.BoxString(RhinoStringPool_AddOrGet(license)))
	this_.populated = true
	return true
}

// Java method: Builder.addLicense(java.lang.String license)
func (this_ *JSDocInfo_Builder) AddLicense(license string) bool {
	_ = license
	if this_.licenseTexts == nil {
		this_.licenseTexts = &j2gort.OrderedSet[string]{}
	}
	if !(this_.licenseTexts.Add(license)) {
		return false
	}
	var txt string = j2gort.CastString(this_.priv_Builder_getProp_8315(StaticGet_JSDocInfo_LICENSE()))
	_ = txt
	return this_.RecordLicense(RhinoStringPool_AddOrGet((j2gort.StrValString(j2gort.Ext_com_google_common_base_Strings_NullToEmpty(txt)) + j2gort.StrValString(license))))
}

// Java method: Builder.recordNoCompile()
func (this_ *JSDocInfo_Builder) RecordNoCompile() bool {
	return this_.priv_Builder_populateBit_8324(StaticGet_JSDocInfo_Bit_NOCOMPILE(), true)
}

// Java method: Builder.recordNoDts()
func (this_ *JSDocInfo_Builder) RecordNoDts() bool {
	return this_.priv_Builder_populateBit_8324(StaticGet_JSDocInfo_Bit_NODTS(), true)
}

// Java method: Builder.recordLegacyNamedTsExport()
func (this_ *JSDocInfo_Builder) RecordLegacyNamedTsExport() bool {
	return this_.priv_Builder_populateBit_8324(StaticGet_JSDocInfo_Bit_LEGACY_NAMED_TS_EXPORT(), true)
}

// Java method: Builder.recordNoCollapse()
func (this_ *JSDocInfo_Builder) RecordNoCollapse() bool {
	return this_.priv_Builder_populateBit_8324(StaticGet_JSDocInfo_Bit_NOCOLLAPSE(), true)
}

// Java method: Builder.recordNoInline()
func (this_ *JSDocInfo_Builder) RecordNoInline() bool {
	return this_.priv_Builder_populateBit_8324(StaticGet_JSDocInfo_Bit_NOINLINE(), true)
}

// Java method: Builder.recordEncourageInlining()
func (this_ *JSDocInfo_Builder) RecordEncourageInlining() bool {
	return this_.priv_Builder_populateBit_8324(StaticGet_JSDocInfo_Bit_ENCOURAGE_INLINING(), true)
}

// Java method: Builder.recordRequireInlining()
func (this_ *JSDocInfo_Builder) RecordRequireInlining() bool {
	return this_.priv_Builder_populateBit_8324(StaticGet_JSDocInfo_Bit_REQUIRE_INLINING(), true)
}

// Java method: Builder.recordPureOrBreakMyCode()
func (this_ *JSDocInfo_Builder) RecordPureOrBreakMyCode() bool {
	return this_.priv_Builder_populateBit_8324(StaticGet_JSDocInfo_Bit_PURE_OR_BREAK_MY_CODE(), true)
}

// Java method: Builder.recordCollapsibleOrBreakMyCode()
func (this_ *JSDocInfo_Builder) RecordCollapsibleOrBreakMyCode() bool {
	return this_.priv_Builder_populateBit_8324(StaticGet_JSDocInfo_Bit_COLLAPSIBLE_OR_BREAK_MY_CODE(), true)
}

// Java method: Builder.recordConstructor()
func (this_ *JSDocInfo_Builder) RecordConstructor() bool {
	return ((!(this_.priv_Builder_hasAnySingletonTypeTags_8304()) && !(this_.priv_Builder_isConstructorOrInterface_8311())) && this_.priv_Builder_populateBit_8324(StaticGet_JSDocInfo_Bit_CONSTRUCTOR(), true))
}

// Java method: Builder.recordImplicitMatch()
func (this_ *JSDocInfo_Builder) RecordImplicitMatch() bool {
	return (((!(this_.priv_Builder_hasAnySingletonTypeTags_8304()) && !(this_.priv_Builder_isConstructorOrInterface_8311())) && this_.priv_Builder_populateBit_8324(StaticGet_JSDocInfo_Bit_RECORD(), true)) && this_.priv_Builder_populateBit_8324(StaticGet_JSDocInfo_Bit_INTERFACE(), true))
}

// Java method: Builder.recordProvideGoog()
func (this_ *JSDocInfo_Builder) RecordProvideGoog() bool {
	return this_.priv_Builder_populateBit_8324(StaticGet_JSDocInfo_Bit_PROVIDE_GOOG(), true)
}

// Java method: Builder.recordProvideAlreadyProvided()
func (this_ *JSDocInfo_Builder) RecordProvideAlreadyProvided() bool {
	return this_.priv_Builder_populateBit_8324(StaticGet_JSDocInfo_Bit_PROVIDE_ALREADY_PROVIDED(), true)
}

// Java method: Builder.isConstructorRecorded()
func (this_ *JSDocInfo_Builder) IsConstructorRecorded() bool {
	return this_.priv_Builder_checkBit_8322(StaticGet_JSDocInfo_Bit_CONSTRUCTOR())
}

// Java method: Builder.recordUnrestricted()
func (this_ *JSDocInfo_Builder) RecordUnrestricted() bool {
	return ((!(this_.priv_Builder_hasAnySingletonTypeTags_8304()) && ((this_.bits & ((StaticGet_JSDocInfo_Bit_INTERFACE().mask | StaticGet_JSDocInfo_Bit_DICT().mask) | StaticGet_JSDocInfo_Bit_STRUCT().mask)) == int64(int32(0)))) && this_.priv_Builder_populateBit_8324(StaticGet_JSDocInfo_Bit_UNRESTRICTED(), true))
}

// Java method: Builder.isUnrestrictedRecorded()
func (this_ *JSDocInfo_Builder) IsUnrestrictedRecorded() bool {
	return this_.priv_Builder_checkBit_8322(StaticGet_JSDocInfo_Bit_UNRESTRICTED())
}

// Java method: Builder.recordAbstract()
func (this_ *JSDocInfo_Builder) RecordAbstract() bool {
	return (((!(this_.priv_Builder_hasAnySingletonTypeTags_8304()) && ((this_.bits & (StaticGet_JSDocInfo_Bit_INTERFACE().mask | StaticGet_JSDocInfo_Bit_FINAL().mask)) == int64(int32(0)))) && (j2gort.Cast[JSDocInfo_Visibility](this_.priv_Builder_getProp_8315(StaticGet_JSDocInfo_VISIBILITY())) != StaticGet_JSDocInfo_Visibility_PRIVATE())) && this_.priv_Builder_populateBit_8324(StaticGet_JSDocInfo_Bit_ABSTRACT(), true))
}

// Java method: Builder.recordStruct()
func (this_ *JSDocInfo_Builder) RecordStruct() bool {
	return ((!(this_.priv_Builder_hasAnySingletonTypeTags_8304()) && ((this_.bits & (StaticGet_JSDocInfo_Bit_DICT().mask | StaticGet_JSDocInfo_Bit_UNRESTRICTED().mask)) == int64(int32(0)))) && this_.priv_Builder_populateBit_8324(StaticGet_JSDocInfo_Bit_STRUCT(), true))
}

// Java method: Builder.isStructRecorded()
func (this_ *JSDocInfo_Builder) IsStructRecorded() bool {
	return this_.priv_Builder_checkBit_8322(StaticGet_JSDocInfo_Bit_STRUCT())
}

// Java method: Builder.recordDict()
func (this_ *JSDocInfo_Builder) RecordDict() bool {
	return ((!(this_.priv_Builder_hasAnySingletonTypeTags_8304()) && ((this_.bits & (StaticGet_JSDocInfo_Bit_STRUCT().mask | StaticGet_JSDocInfo_Bit_UNRESTRICTED().mask)) == int64(int32(0)))) && this_.priv_Builder_populateBit_8324(StaticGet_JSDocInfo_Bit_DICT(), true))
}

// Java method: Builder.isDictRecorded()
func (this_ *JSDocInfo_Builder) IsDictRecorded() bool {
	return this_.priv_Builder_checkBit_8322(StaticGet_JSDocInfo_Bit_DICT())
}

// Java method: Builder.recordOverride()
func (this_ *JSDocInfo_Builder) RecordOverride() bool {
	return this_.priv_Builder_populateBit_8324(StaticGet_JSDocInfo_Bit_OVERRIDE(), true)
}

// Java method: Builder.recordDeprecated()
func (this_ *JSDocInfo_Builder) RecordDeprecated() bool {
	return this_.priv_Builder_populateBit_8324(StaticGet_JSDocInfo_Bit_DEPRECATED(), true)
}

// Java method: Builder.recordInterface()
func (this_ *JSDocInfo_Builder) RecordInterface() bool {
	return ((!(this_.priv_Builder_hasAnySingletonTypeTags_8304()) && ((this_.bits & (StaticGet_JSDocInfo_Bit_CONSTRUCTOR().mask | StaticGet_JSDocInfo_Bit_ABSTRACT().mask)) == int64(int32(0)))) && this_.priv_Builder_populateBit_8324(StaticGet_JSDocInfo_Bit_INTERFACE(), true))
}

// Java method: Builder.recordExport()
func (this_ *JSDocInfo_Builder) RecordExport() bool {
	return this_.priv_Builder_populateBit_8324(StaticGet_JSDocInfo_Bit_EXPORT(), true)
}

// Java method: Builder.removeExport()
func (this_ *JSDocInfo_Builder) RemoveExport() bool {
	return this_.priv_Builder_populateBit_8324(StaticGet_JSDocInfo_Bit_EXPORT(), false)
}

// Java method: Builder.recordImplicitCast()
func (this_ *JSDocInfo_Builder) RecordImplicitCast() bool {
	return this_.priv_Builder_populateBit_8324(StaticGet_JSDocInfo_Bit_IMPLICITCAST(), true)
}

// Java method: Builder.recordNoSideEffects()
func (this_ *JSDocInfo_Builder) RecordNoSideEffects() bool {
	return ((!(this_.priv_Builder_hasAnySingletonSideEffectTags_8310()) && this_.priv_Builder_isPropEmpty_8303(StaticGet_JSDocInfo_THROWS_ANNOTATIONS())) && this_.priv_Builder_populateBit_8324(StaticGet_JSDocInfo_Bit_NOSIDEEFFECTS(), true))
}

// Java method: Builder.isNoSideEffectsRecorded()
func (this_ *JSDocInfo_Builder) IsNoSideEffectsRecorded() bool {
	return this_.priv_Builder_checkBit_8322(StaticGet_JSDocInfo_Bit_NOSIDEEFFECTS())
}

// Java method: Builder.isModifiesRecorded()
func (this_ *JSDocInfo_Builder) IsModifiesRecorded() bool {
	return !(this_.priv_Builder_isPropEmpty_8303(StaticGet_JSDocInfo_MODIFIES()))
}

// Java method: Builder.isThrowsRecorded()
func (this_ *JSDocInfo_Builder) IsThrowsRecorded() bool {
	return !(this_.priv_Builder_isPropEmpty_8303(StaticGet_JSDocInfo_THROWS_ANNOTATIONS()))
}

// Java method: Builder.recordExterns()
func (this_ *JSDocInfo_Builder) RecordExterns() bool {
	return this_.priv_Builder_populateBit_8324(StaticGet_JSDocInfo_Bit_EXTERNS(), true)
}

// Java method: Builder.recordNoCoverage()
func (this_ *JSDocInfo_Builder) RecordNoCoverage() bool {
	return this_.priv_Builder_populateBit_8324(StaticGet_JSDocInfo_Bit_NOCOVERAGE(), true)
}

// Java method: Builder.recordTypeSummary()
func (this_ *JSDocInfo_Builder) RecordTypeSummary() bool {
	return this_.priv_Builder_populateBit_8324(StaticGet_JSDocInfo_Bit_TYPE_SUMMARY(), true)
}

// Java method: Builder.isInterfaceRecorded()
func (this_ *JSDocInfo_Builder) IsInterfaceRecorded() bool {
	return this_.priv_Builder_checkBit_8322(StaticGet_JSDocInfo_Bit_INTERFACE())
}

// Java method: Builder.hasParameter(java.lang.String name)
func (this_ *JSDocInfo_Builder) HasParameter(name string) bool {
	_ = name
	var params *JSDocInfo_TypeMap = j2gort.CastPtr[JSDocInfo_TypeMap](this_.priv_Builder_getProp_8315(ToIface_JSDocInfo_Property(StaticGet_JSDocInfo_PARAMETERS())))
	_ = params
	return ((params != nil) && params.ContainsKey(name))
}

// Java method: Builder.recordImplementedInterface(com.google.javascript.rhino.JSTypeExpression interfaceType)
func (this_ *JSDocInfo_Builder) RecordImplementedInterface(interfaceType JSTypeExpression) bool {
	_ = interfaceType
	return ((interfaceType != (JSTypeExpression{})) && this_.priv_Builder_addUnique_8272(StaticGet_JSDocInfo_IMPLEMENTED_INTERFACES(), interfaceType))
}

// Java method: Builder.recordExtendedInterface(com.google.javascript.rhino.JSTypeExpression interfaceType)
func (this_ *JSDocInfo_Builder) RecordExtendedInterface(interfaceType JSTypeExpression) bool {
	_ = interfaceType
	return ((interfaceType != (JSTypeExpression{})) && this_.priv_Builder_addUnique_8272(StaticGet_JSDocInfo_EXTENDED_INTERFACES(), interfaceType))
}

// Java method: Builder.addUnique(com.google.javascript.rhino.JSDocInfo.TypeListProperty prop, com.google.javascript.rhino.JSTypeExpression elem)
func (this_ *JSDocInfo_Builder) priv_Builder_addUnique_8272(prop *JSDocInfo_TypeListProperty, elem JSTypeExpression) bool {
	_ = prop
	_ = elem
	var list *JSDocInfo_TypeList = j2gort.CastPtr[JSDocInfo_TypeList](this_.priv_Builder_getProp_8315(ToIface_JSDocInfo_Property(prop)))
	_ = list
	if list == nil {
		list = NewJSDocInfo_TypeList(int32(2))
		this_.priv_Builder_setProp_8314(ToIface_JSDocInfo_Property(prop), j2gort.PtrToAny(list))
	}
	if !(list.AddUnique(elem)) {
		return false
	}
	this_.populated = true
	return true
}

// Java method: Builder.recordLends(com.google.javascript.rhino.JSTypeExpression name)
func (this_ *JSDocInfo_Builder) RecordLends(name JSTypeExpression) bool {
	_ = name
	return (!(this_.priv_Builder_hasAnyTypeRelatedTags_8300()) && this_.priv_Builder_populateTypeProp_8309(StaticGet_JSDocInfo_LENDS_NAME(), name))
}

// Java method: Builder.isNgInjectRecorded()
func (this_ *JSDocInfo_Builder) IsNgInjectRecorded() bool {
	return this_.priv_Builder_checkBit_8322(StaticGet_JSDocInfo_Bit_NG_INJECT())
}

// Java method: Builder.recordNgInject(boolean ngInject)
func (this_ *JSDocInfo_Builder) RecordNgInject(ngInject bool) bool {
	_ = ngInject
	return this_.priv_Builder_populateBit_8324(StaticGet_JSDocInfo_Bit_NG_INJECT(), true)
}

// Java method: Builder.isWizactionRecorded()
func (this_ *JSDocInfo_Builder) IsWizactionRecorded() bool {
	return this_.priv_Builder_checkBit_8322(StaticGet_JSDocInfo_Bit_WIZ_ACTION())
}

// Java method: Builder.recordWizaction()
func (this_ *JSDocInfo_Builder) RecordWizaction() bool {
	return this_.priv_Builder_populateBit_8324(StaticGet_JSDocInfo_Bit_WIZ_ACTION(), true)
}

// Java method: Builder.isWizcallbackRecorded()
func (this_ *JSDocInfo_Builder) IsWizcallbackRecorded() bool {
	return this_.priv_Builder_checkBit_8322(StaticGet_JSDocInfo_Bit_WIZ_CALLBACK())
}

// Java method: Builder.recordWizcallback()
func (this_ *JSDocInfo_Builder) RecordWizcallback() bool {
	return this_.priv_Builder_populateBit_8324(StaticGet_JSDocInfo_Bit_WIZ_CALLBACK(), true)
}

// Java method: Builder.isPolymerBehaviorRecorded()
func (this_ *JSDocInfo_Builder) IsPolymerBehaviorRecorded() bool {
	return this_.priv_Builder_checkBit_8322(StaticGet_JSDocInfo_Bit_POLYMER_BEHAVIOR())
}

// Java method: Builder.recordPolymerBehavior()
func (this_ *JSDocInfo_Builder) RecordPolymerBehavior() bool {
	return this_.priv_Builder_populateBit_8324(StaticGet_JSDocInfo_Bit_POLYMER_BEHAVIOR(), true)
}

// Java method: Builder.isPolymerRecorded()
func (this_ *JSDocInfo_Builder) IsPolymerRecorded() bool {
	return this_.priv_Builder_checkBit_8322(StaticGet_JSDocInfo_Bit_POLYMER())
}

// Java method: Builder.recordPolymer()
func (this_ *JSDocInfo_Builder) RecordPolymer() bool {
	return this_.priv_Builder_populateBit_8324(StaticGet_JSDocInfo_Bit_POLYMER(), true)
}

// Java method: Builder.isCustomElementRecorded()
func (this_ *JSDocInfo_Builder) IsCustomElementRecorded() bool {
	return this_.priv_Builder_checkBit_8322(StaticGet_JSDocInfo_Bit_CUSTOM_ELEMENT())
}

// Java method: Builder.recordCustomElement()
func (this_ *JSDocInfo_Builder) RecordCustomElement() bool {
	return this_.priv_Builder_populateBit_8324(StaticGet_JSDocInfo_Bit_CUSTOM_ELEMENT(), true)
}

// Java method: Builder.isMixinClassRecorded()
func (this_ *JSDocInfo_Builder) IsMixinClassRecorded() bool {
	return this_.priv_Builder_checkBit_8322(StaticGet_JSDocInfo_Bit_MIXIN_CLASS())
}

// Java method: Builder.recordMixinClass()
func (this_ *JSDocInfo_Builder) RecordMixinClass() bool {
	return this_.priv_Builder_populateBit_8324(StaticGet_JSDocInfo_Bit_MIXIN_CLASS(), true)
}

// Java method: Builder.isMixinFunctionRecorded()
func (this_ *JSDocInfo_Builder) IsMixinFunctionRecorded() bool {
	return this_.priv_Builder_checkBit_8322(StaticGet_JSDocInfo_Bit_MIXIN_FUNCTION())
}

// Java method: Builder.recordMixinFunction()
func (this_ *JSDocInfo_Builder) RecordMixinFunction() bool {
	return this_.priv_Builder_populateBit_8324(StaticGet_JSDocInfo_Bit_MIXIN_FUNCTION(), true)
}

// Java method: Builder.isSassGeneratedCssTsRecorded()
func (this_ *JSDocInfo_Builder) IsSassGeneratedCssTsRecorded() bool {
	return this_.priv_Builder_checkBit_8322(StaticGet_JSDocInfo_Bit_SASS_GENERATED_CSS_TS())
}

// Java method: Builder.recordSassGeneratedCssTs()
func (this_ *JSDocInfo_Builder) RecordSassGeneratedCssTs() {
	this_.priv_Builder_populateBit_8324(StaticGet_JSDocInfo_Bit_SASS_GENERATED_CSS_TS(), true)
}

// Java method: Builder.logTypeInCompiler()
func (this_ *JSDocInfo_Builder) LogTypeInCompiler() bool {
	return this_.priv_Builder_checkBit_8322(StaticGet_JSDocInfo_Bit_LOG_TYPE_IN_COMPILER())
}

// Java method: Builder.recordLogTypeInCompiler()
func (this_ *JSDocInfo_Builder) RecordLogTypeInCompiler() bool {
	return this_.priv_Builder_populateBit_8324(StaticGet_JSDocInfo_Bit_LOG_TYPE_IN_COMPILER(), true)
}

// Java method: Builder.isClosureUnawareCode()
func (this_ *JSDocInfo_Builder) IsClosureUnawareCode() bool {
	return (j2gort.Cast[JSDocInfo_PerFileClosureUnawareMode](this_.priv_Builder_getProp_8315(StaticGet_JSDocInfo_CLOSURE_UNAWARE_CONFIG_VALUE())) != 0)
}

// Java method: Builder.recordClosureUnawareCode()
func (this_ *JSDocInfo_Builder) RecordClosureUnawareCode_0() bool {
	return this_.priv_Builder_populateProp_8321(StaticGet_JSDocInfo_CLOSURE_UNAWARE_CONFIG_VALUE(), j2gort.BoxNullableVal(StaticGet_JSDocInfo_PerFileClosureUnawareMode_UNSPECIFIED()))
}

// Java method: Builder.recordClosureUnawareCode(com.google.javascript.rhino.JSDocInfo.PerFileClosureUnawareMode mode)
func (this_ *JSDocInfo_Builder) RecordClosureUnawareCode_1(mode JSDocInfo_PerFileClosureUnawareMode) bool {
	_ = mode
	return this_.priv_Builder_populateProp_8321(StaticGet_JSDocInfo_CLOSURE_UNAWARE_CONFIG_VALUE(), j2gort.BoxNullableVal(mode))
}

// Java method: Builder.removeClosureUnawareCode()
func (this_ *JSDocInfo_Builder) RemoveClosureUnawareCode() bool {
	return (j2gort.Cast[JSDocInfo_PerFileClosureUnawareMode](this_.priv_Builder_removeProp_8316(StaticGet_JSDocInfo_CLOSURE_UNAWARE_CONFIG_VALUE())) != 0)
}

// Java method: Builder.recordUsedViaDotConstructor()
func (this_ *JSDocInfo_Builder) RecordUsedViaDotConstructor() bool {
	return this_.priv_Builder_populateBit_8324(StaticGet_JSDocInfo_Bit_USED_VIA_DOT_CONSTRUCTOR(), true)
}

// Java method: Builder.setType(com.google.javascript.rhino.JSTypeExpression type)
func (this_ *JSDocInfo_Builder) SetType(type_ JSTypeExpression) *JSDocInfo_Builder {
	_ = type_
	this_.priv_Builder_removeTypeProp_8308(StaticGet_JSDocInfo_RETURN_TYPE())
	this_.priv_Builder_removeTypeProp_8308(StaticGet_JSDocInfo_ENUM_PARAMETER_TYPE())
	this_.priv_Builder_removeTypeProp_8308(StaticGet_JSDocInfo_TYPEDEF_TYPE())
	this_.priv_Builder_setTypeProp_8307(StaticGet_JSDocInfo_TYPE(), type_)
	return this_
}

// Java method: Builder.hasAnyTypeRelatedTags()
func (this_ *JSDocInfo_Builder) priv_Builder_hasAnyTypeRelatedTags_8300() bool {
	return ((((this_.singleTypeKeys != int64(0)) || ((this_.bits & ((StaticGet_JSDocInfo_Bit_CONSTRUCTOR().mask | StaticGet_JSDocInfo_Bit_INTERFACE().mask) | StaticGet_JSDocInfo_Bit_ABSTRACT().mask)) != int64(int32(0)))) || this_.priv_Builder_hasAnyParameters_8301()) || !(this_.priv_Builder_isTypeListEmpty_8302(StaticGet_JSDocInfo_EXTENDED_INTERFACES())))
}

// Java method: Builder.hasAnyParameters()
func (this_ *JSDocInfo_Builder) priv_Builder_hasAnyParameters_8301() bool {
	var params *JSDocInfo_TypeMap = j2gort.CastPtr[JSDocInfo_TypeMap](this_.priv_Builder_getProp_8315(ToIface_JSDocInfo_Property(StaticGet_JSDocInfo_PARAMETERS())))
	_ = params
	return ((params != nil) && !(params.IsEmpty()))
}

// Java method: Builder.isTypeListEmpty(com.google.javascript.rhino.JSDocInfo.TypeListProperty prop)
func (this_ *JSDocInfo_Builder) priv_Builder_isTypeListEmpty_8302(prop *JSDocInfo_TypeListProperty) bool {
	_ = prop
	var list *JSDocInfo_TypeList = j2gort.CastPtr[JSDocInfo_TypeList](this_.priv_Builder_getProp_8315(ToIface_JSDocInfo_Property(prop)))
	_ = list
	return ((list == nil) || list.IsEmpty())
}

// Java method: Builder.isPropEmpty(com.google.javascript.rhino.JSDocInfo.Property<? extends java.util.Collection<?>> prop)
func (this_ *JSDocInfo_Builder) priv_Builder_isPropEmpty_8303(prop JSDocInfo_Property) bool {
	_ = prop
	var c j2gort.Ext_java_util_Collection = j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.priv_Builder_getProp_8315(prop))
	_ = c
	return ((c == nil) || c.IsEmpty())
}

// Java method: Builder.hasAnySingletonTypeTags()
func (this_ *JSDocInfo_Builder) priv_Builder_hasAnySingletonTypeTags_8304() bool {
	return ((this_.singleTypeKeys & ((StaticGet_JSDocInfo_TYPE().GetJSDocInfo_PropertyBase().mask | StaticGet_JSDocInfo_TYPEDEF_TYPE().GetJSDocInfo_PropertyBase().mask) | StaticGet_JSDocInfo_ENUM_PARAMETER_TYPE().GetJSDocInfo_PropertyBase().mask)) != int64(0))
}

// Java method: Builder.getTypeValAt(int idx)
func (this_ *JSDocInfo_Builder) GetTypeValAt(idx int32) JSTypeExpression {
	_ = idx
	switch idx {
	case int32(0):
		return this_.typeVal0
	case int32(1):
		return this_.typeVal1
	case int32(2):
		return this_.typeVal2
	default:
		return this_.typeVal3
	}
	panic("unreachable switch expression")
}

// Java method: Builder.setTypeValAt(int idx, com.google.javascript.rhino.@org.jspecify.annotations.Nullable JSTypeExpression v)
func (this_ *JSDocInfo_Builder) SetTypeValAt(idx int32, v JSTypeExpression) {
	_ = idx
	_ = v
	switch idx {
	case int32(0):
		this_.typeVal0 = v
	case int32(1):
		this_.typeVal1 = v
	case int32(2):
		this_.typeVal2 = v
	default:
		this_.typeVal3 = v
	}
}

// Java method: Builder.setTypeProp(com.google.javascript.rhino.JSDocInfo.TypeProperty prop, com.google.javascript.rhino.@org.jspecify.annotations.Nullable JSTypeExpression value)
func (this_ *JSDocInfo_Builder) priv_Builder_setTypeProp_8307(prop *JSDocInfo_TypeProperty, value JSTypeExpression) {
	_ = prop
	_ = value
	if value == (JSTypeExpression{}) {
		this_.priv_Builder_removeTypeProp_8308(prop)
		return
	}
	var idx int32 = j2gort.BoxLong_BitCount((this_.singleTypeKeys & (prop.GetJSDocInfo_PropertyBase().mask - int64(1))))
	_ = idx
	if (this_.singleTypeKeys & prop.GetJSDocInfo_PropertyBase().mask) != int64(0) {
		this_.SetTypeValAt(idx, value)
		return
	}
	var count int32 = j2gort.BoxLong_BitCount(this_.singleTypeKeys)
	_ = count
	if !(count < int32(4)) {
		j2gort.CheckState(false)
	}
	{
		var i int32 = count
		_ = i
		for ; i > idx; i = i - 1 {
			this_.SetTypeValAt(i, this_.GetTypeValAt((i - int32(1))))
		}
	}
	this_.SetTypeValAt(idx, value)
	this_.singleTypeKeys = (this_.singleTypeKeys | prop.GetJSDocInfo_PropertyBase().mask)
}

// Java method: Builder.removeTypeProp(com.google.javascript.rhino.JSDocInfo.TypeProperty prop)
func (this_ *JSDocInfo_Builder) priv_Builder_removeTypeProp_8308(prop *JSDocInfo_TypeProperty) {
	_ = prop
	if (this_.singleTypeKeys & prop.GetJSDocInfo_PropertyBase().mask) == int64(0) {
		return
	}
	var idx int32 = j2gort.BoxLong_BitCount((this_.singleTypeKeys & (prop.GetJSDocInfo_PropertyBase().mask - int64(1))))
	_ = idx
	var count int32 = j2gort.BoxLong_BitCount(this_.singleTypeKeys)
	_ = count
	{
		var i int32 = idx
		_ = i
		for ; i < (count - int32(1)); i = i + 1 {
			this_.SetTypeValAt(i, this_.GetTypeValAt((i + int32(1))))
		}
	}
	this_.SetTypeValAt((count - int32(1)), (JSTypeExpression{}))
	this_.singleTypeKeys = (this_.singleTypeKeys & ^(prop.GetJSDocInfo_PropertyBase().mask))
}

// Java method: Builder.populateTypeProp(com.google.javascript.rhino.JSDocInfo.TypeProperty prop, com.google.javascript.rhino.@org.jspecify.annotations.Nullable JSTypeExpression value)
func (this_ *JSDocInfo_Builder) priv_Builder_populateTypeProp_8309(prop *JSDocInfo_TypeProperty, value JSTypeExpression) bool {
	_ = prop
	_ = value
	this_.populated = true
	if (this_.singleTypeKeys & prop.GetJSDocInfo_PropertyBase().mask) != int64(0) {
		return false
	}
	this_.priv_Builder_setTypeProp_8307(prop, value)
	return true
}

// Java method: Builder.hasAnySingletonSideEffectTags()
func (this_ *JSDocInfo_Builder) priv_Builder_hasAnySingletonSideEffectTags_8310() bool {
	return (this_.priv_Builder_checkBit_8322(StaticGet_JSDocInfo_Bit_NOSIDEEFFECTS()) || !(this_.priv_Builder_isPropEmpty_8303(StaticGet_JSDocInfo_MODIFIES())))
}

// Java method: Builder.isConstructorOrInterface()
func (this_ *JSDocInfo_Builder) priv_Builder_isConstructorOrInterface_8311() bool {
	return ((this_.bits & (StaticGet_JSDocInfo_Bit_CONSTRUCTOR().mask | StaticGet_JSDocInfo_Bit_INTERFACE().mask)) != int64(int32(0)))
}

// Java method: Builder.getInlineValAt(int idx)
func (this_ *JSDocInfo_Builder) GetInlineValAt(idx int32) any {
	_ = idx
	switch idx {
	case int32(0):
		return this_.propVal0
	case int32(1):
		return this_.propVal1
	case int32(2):
		return this_.propVal2
	default:
		return this_.propVal3
	}
	panic("unreachable switch expression")
}

// Java method: Builder.setInlineValAt(int idx, java.lang.@org.jspecify.annotations.Nullable Object v)
func (this_ *JSDocInfo_Builder) SetInlineValAt(idx int32, v any) {
	_ = idx
	_ = v
	switch idx {
	case int32(0):
		this_.propVal0 = v
	case int32(1):
		this_.propVal1 = v
	case int32(2):
		this_.propVal2 = v
	default:
		this_.propVal3 = v
	}
}

// Java method: Builder.setProp(com.google.javascript.rhino.JSDocInfo.Property<T> prop, @org.jspecify.annotations.Nullable T value)
func (this_ *JSDocInfo_Builder) priv_Builder_setProp_8314(prop JSDocInfo_Property, value any) {
	_ = prop
	_ = value
	if value == nil {
		this_.priv_Builder_removeProp_8316(prop)
		return
	}
	if this_.propValues != nil {
		{
			a_95 := this_.propValues
			i_96 := prop.GetJSDocInfo_PropertyBase().bit
			a_95[i_96] = value
		}
		this_.propKeys = (this_.propKeys | prop.GetJSDocInfo_PropertyBase().mask)
		return
	}
	var idx int32 = j2gort.BoxLong_BitCount((this_.propKeys & (prop.GetJSDocInfo_PropertyBase().mask - int64(1))))
	_ = idx
	if (this_.propKeys & prop.GetJSDocInfo_PropertyBase().mask) != int64(0) {
		this_.SetInlineValAt(idx, value)
		return
	}
	var count int32 = j2gort.BoxLong_BitCount(this_.propKeys)
	_ = count
	if count < int32(4) {
		{
			var i int32 = count
			_ = i
			for ; i > idx; i = i - 1 {
				this_.SetInlineValAt(i, this_.GetInlineValAt((i - int32(1))))
			}
		}
		this_.SetInlineValAt(idx, value)
		this_.propKeys = (this_.propKeys | prop.GetJSDocInfo_PropertyBase().mask)
		return
	}
	var vals []any = make([]any, StaticGet_JSDocInfo_Property_BitCounter())
	_ = vals
	var rem int64 = this_.propKeys
	_ = rem
	{
		var i_2 int32 = int32(0)
		_ = i_2
		for ; i_2 < int32(4); i_2 = i_2 + 1 {
			var b int32 = j2gort.BoxLong_NumberOfTrailingZeros(rem)
			_ = b
			rem = (rem & (rem - int64(1)))
			{
				a_97 := vals
				i_98 := b
				a_97[i_98] = this_.GetInlineValAt(i_2)
			}
		}
	}
	this_.propVal0 = nil
	this_.propVal1 = nil
	this_.propVal2 = nil
	this_.propVal3 = nil
	{
		a_99 := vals
		i_100 := prop.GetJSDocInfo_PropertyBase().bit
		a_99[i_100] = value
	}
	this_.propValues = vals
	this_.propKeys = (this_.propKeys | prop.GetJSDocInfo_PropertyBase().mask)
}

// Java method: Builder.getProp(com.google.javascript.rhino.JSDocInfo.Property<T> prop)
func (this_ *JSDocInfo_Builder) priv_Builder_getProp_8315(prop JSDocInfo_Property) any {
	_ = prop
	if (this_.propKeys & prop.GetJSDocInfo_PropertyBase().mask) == int64(0) {
		return nil
	}
	if this_.propValues != nil {
		return this_.propValues[prop.GetJSDocInfo_PropertyBase().bit]
	}
	return this_.GetInlineValAt(j2gort.BoxLong_BitCount((this_.propKeys & (prop.GetJSDocInfo_PropertyBase().mask - int64(1)))))
}

// Java method: Builder.removeProp(com.google.javascript.rhino.JSDocInfo.Property<T> prop)
func (this_ *JSDocInfo_Builder) priv_Builder_removeProp_8316(prop JSDocInfo_Property) any {
	_ = prop
	if (this_.propKeys & prop.GetJSDocInfo_PropertyBase().mask) == int64(0) {
		return nil
	}
	var old any = nil
	_ = old
	if this_.propValues != nil {
		old = this_.propValues[prop.GetJSDocInfo_PropertyBase().bit]
		{
			a_101 := this_.propValues
			i_102 := prop.GetJSDocInfo_PropertyBase().bit
			a_101[i_102] = nil
		}
	} else {
		var idx int32 = j2gort.BoxLong_BitCount((this_.propKeys & (prop.GetJSDocInfo_PropertyBase().mask - int64(1))))
		_ = idx
		var count int32 = j2gort.BoxLong_BitCount(this_.propKeys)
		_ = count
		old = this_.GetInlineValAt(idx)
		{
			var i int32 = idx
			_ = i
			for ; i < (count - int32(1)); i = i + 1 {
				this_.SetInlineValAt(i, this_.GetInlineValAt((i + int32(1))))
			}
		}
		this_.SetInlineValAt((count - int32(1)), nil)
	}
	this_.propKeys = (this_.propKeys & ^(prop.GetJSDocInfo_PropertyBase().mask))
	return old
}

// Java method: Builder.getOrCreateStringList(com.google.javascript.rhino.JSDocInfo.Property<java.util.List<java.lang.String>> prop)
func (this_ *JSDocInfo_Builder) priv_Builder_getOrCreateStringList_8317(prop JSDocInfo_Property) j2gort.Ext_java_util_List {
	_ = prop
	var list j2gort.Ext_java_util_List = j2gort.AsIface[j2gort.Ext_java_util_List](this_.priv_Builder_getProp_8315(prop))
	_ = list
	if list == nil {
		list = j2gort.NewExt_java_util_ArrayList()
		this_.priv_Builder_setProp_8314(prop, list)
	}
	return list
}

// Java method: Builder.populateTypePropEntry(com.google.javascript.rhino.JSDocInfo.TypeMapProperty prop, java.lang.String key, com.google.javascript.rhino.@org.jspecify.annotations.Nullable JSTypeExpression value)
func (this_ *JSDocInfo_Builder) priv_Builder_populateTypePropEntry_8318(prop *JSDocInfo_TypeMapProperty, key string, value JSTypeExpression) bool {
	_ = prop
	_ = key
	_ = value
	var map_ *JSDocInfo_TypeMap = j2gort.CastPtr[JSDocInfo_TypeMap](this_.priv_Builder_getProp_8315(ToIface_JSDocInfo_Property(prop)))
	_ = map_
	if map_ == nil {
		map_ = NewJSDocInfo_TypeMap(int32(4))
		this_.priv_Builder_setProp_8314(ToIface_JSDocInfo_Property(prop), j2gort.PtrToAny(map_))
	}
	return (map_.PutIfAbsent(key, value) && (func() bool { var asgn_103 bool = true; this_.populated = asgn_103; return asgn_103 }()))
}

// Java method: Builder.populateStringPropEntry(com.google.javascript.rhino.JSDocInfo.StringMapProperty prop, java.lang.String key, java.lang.@org.jspecify.annotations.Nullable String value)
func (this_ *JSDocInfo_Builder) priv_Builder_populateStringPropEntry_8319(prop *JSDocInfo_StringMapProperty, key string, value string) bool {
	_ = prop
	_ = key
	_ = value
	var map_ *JSDocInfo_StringMap = j2gort.CastPtr[JSDocInfo_StringMap](this_.priv_Builder_getProp_8315(ToIface_JSDocInfo_Property(prop)))
	_ = map_
	if map_ == nil {
		map_ = NewJSDocInfo_StringMap()
		this_.priv_Builder_setProp_8314(ToIface_JSDocInfo_Property(prop), j2gort.PtrToAny(map_))
	}
	return (map_.PutIfAbsent(key, value) && (func() bool { var asgn_104 bool = true; this_.populated = asgn_104; return asgn_104 }()))
}

// Java method: Builder.populateNodePropEntry(com.google.javascript.rhino.JSDocInfo.NodeMapProperty prop, java.lang.String key, com.google.javascript.rhino.Node value)
func (this_ *JSDocInfo_Builder) priv_Builder_populateNodePropEntry_8320(prop *JSDocInfo_NodeMapProperty, key string, value Node) bool {
	_ = prop
	_ = key
	_ = value
	var map_ *JSDocInfo_NodeMap = j2gort.CastPtr[JSDocInfo_NodeMap](this_.priv_Builder_getProp_8315(ToIface_JSDocInfo_Property(prop)))
	_ = map_
	if map_ == nil {
		map_ = NewJSDocInfo_NodeMap(int32(4))
		this_.priv_Builder_setProp_8314(ToIface_JSDocInfo_Property(prop), j2gort.PtrToAny(map_))
	}
	return (map_.PutIfAbsent(key, value) && (func() bool { var asgn_105 bool = true; this_.populated = asgn_105; return asgn_105 }()))
}

// Java method: Builder.populateProp(com.google.javascript.rhino.JSDocInfo.Property<T> prop, T value)
func (this_ *JSDocInfo_Builder) priv_Builder_populateProp_8321(prop JSDocInfo_Property, value any) bool {
	_ = prop
	_ = value
	this_.populated = true
	if (this_.propKeys & prop.GetJSDocInfo_PropertyBase().mask) != int64(0) {
		return false
	}
	this_.priv_Builder_setProp_8314(prop, value)
	return true
}

// Java method: Builder.checkBit(com.google.javascript.rhino.JSDocInfo.Bit bit)
func (this_ *JSDocInfo_Builder) priv_Builder_checkBit_8322(bit *JSDocInfo_Bit) bool {
	_ = bit
	return ((this_.bits & bit.mask) != int64(int32(0)))
}

// Java method: Builder.setBit(com.google.javascript.rhino.JSDocInfo.Bit bit, boolean value)
func (this_ *JSDocInfo_Builder) priv_Builder_setBit_8323(bit *JSDocInfo_Bit, value bool) {
	_ = bit
	_ = value
	if value {
		this_.bits = (this_.bits | bit.mask)
	} else {
		this_.bits = (this_.bits & ^(bit.mask))
	}
}

// Java method: Builder.populateBit(com.google.javascript.rhino.JSDocInfo.Bit bit, boolean value)
func (this_ *JSDocInfo_Builder) priv_Builder_populateBit_8324(bit *JSDocInfo_Bit, value bool) bool {
	_ = bit
	_ = value
	if this_.priv_Builder_checkBit_8322(bit) != value {
		this_.priv_Builder_setBit_8323(bit, value)
		return func() bool { var asgn_106 bool = true; this_.populated = asgn_106; return asgn_106 }()
	}
	return false
}
