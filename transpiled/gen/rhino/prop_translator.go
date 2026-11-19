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
 *   Google Inc.
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */
// Port of com/google/javascript/rhino/PropTranslator.java (generated mechanically by j2go).

package rhino

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.rhino.PropTranslator
type PropTranslator struct {
	idHash_ atomic.Int32
}

func (this_ *PropTranslator) BindSelf_(s any) {
	_ = s
}
func (this_ *PropTranslator) IsJavaNil_() bool { return this_ == nil }
func (*PropTranslator) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.PropTranslator", "PropTranslator"
}
func (this_ *PropTranslator) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*PropTranslator) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.PropTranslator", "com.google.javascript.rhino.PropTranslator", func(v_ any) bool {
		return j2gort.InstanceOf[*PropTranslator](v_)
	})
}

var static_PropTranslator_protoToRhinoProp []Node_Prop

func StaticGet_PropTranslator_ProtoToRhinoProp() []Node_Prop {
	clinit_PropTranslator()
	return static_PropTranslator_protoToRhinoProp
}
func StaticSet_PropTranslator_ProtoToRhinoProp(v []Node_Prop) {
	clinit_PropTranslator()
	static_PropTranslator_protoToRhinoProp = v
}

var static_PropTranslator_rhinoToProtoProp []*j2gort.Ext_com_google_javascript_jscomp_serialization_NodeProperty

func StaticGet_PropTranslator_RhinoToProtoProp() []*j2gort.Ext_com_google_javascript_jscomp_serialization_NodeProperty {
	clinit_PropTranslator()
	return static_PropTranslator_rhinoToProtoProp
}
func StaticSet_PropTranslator_RhinoToProtoProp(v []*j2gort.Ext_com_google_javascript_jscomp_serialization_NodeProperty) {
	clinit_PropTranslator()
	static_PropTranslator_rhinoToProtoProp = v
}

var clinitDone_PropTranslator uint32
var clinitState_PropTranslator j2gort.ClinitState

func clinit_PropTranslator() {
	if atomic.LoadUint32(&clinitDone_PropTranslator) != 0 {
		return
	}
	clinitSlow_PropTranslator()
}

func clinitSlow_PropTranslator() {
	if !j2gort.EnterClinit(&clinitState_PropTranslator, &clinitDone_PropTranslator) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_PropTranslator, &clinitDone_PropTranslator)
	static_PropTranslator_protoToRhinoProp = make([]Node_Prop, int32(len(j2gort.Ext_com_google_javascript_jscomp_serialization_NodeProperty_Values())))
	static_PropTranslator_rhinoToProtoProp = make([]*j2gort.Ext_com_google_javascript_jscomp_serialization_NodeProperty, int32(len(Node_Prop_Values())))
	PropTranslator_SetProps()
	PropTranslator_CheckUnexpectedNullProtoProps()
}

func (this_ *PropTranslator) initFields_PropTranslator() {
	_ = this_
}

// Java constructor: PropTranslator()
func NewPropTranslator() *PropTranslator {
	clinit_PropTranslator()
	this_ := &PropTranslator{}
	InitPropTranslator(this_)
	return this_
}

func InitPropTranslator(this_ *PropTranslator) {
	this_.initFields_PropTranslator()
	panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError")))
}

// Java static method: PropTranslator.serialize(com.google.javascript.rhino.Node.Prop x)
func PropTranslator_Serialize(x Node_Prop) *j2gort.Ext_com_google_javascript_jscomp_serialization_NodeProperty {
	clinit_PropTranslator()
	_ = x
	return StaticGet_PropTranslator_RhinoToProtoProp()[(x).Ordinal()]
}

// Java static method: PropTranslator.deserialize(com.google.javascript.jscomp.serialization.NodeProperty x)
func PropTranslator_Deserialize(x *j2gort.Ext_com_google_javascript_jscomp_serialization_NodeProperty) Node_Prop {
	clinit_PropTranslator()
	_ = x
	return StaticGet_PropTranslator_ProtoToRhinoProp()[j2gort.EnumOrdinal(x)]
}

// Java static method: PropTranslator.setProps()
func PropTranslator_SetProps() {
	clinit_PropTranslator()
	for _, el_1 := range j2gort.CheckNotNullSlice(Node_Prop_Values()) {
		rhinoProp := el_1
		_ = rhinoProp
		var protoProp *j2gort.Ext_com_google_javascript_jscomp_serialization_NodeProperty = PropTranslator_SerializeProp(rhinoProp)
		_ = protoProp
		if protoProp != nil {
			if !(protoProp.GetNumber() < int32(63)) {
				j2gort.CheckStateFmt(false, "enum %s value %s", j2gort.PtrToAny(protoProp), j2gort.ValueOfInt(protoProp.GetNumber()))
			}
			{
				a_2 := StaticGet_PropTranslator_ProtoToRhinoProp()
				i_3 := j2gort.EnumOrdinal(protoProp)
				a_2[i_3] = rhinoProp
			}
			{
				a_4 := StaticGet_PropTranslator_RhinoToProtoProp()
				i_5 := (rhinoProp).Ordinal()
				a_4[i_5] = protoProp
			}
		}
	}
}

// Java static method: PropTranslator.serializeProp(com.google.javascript.rhino.Node.Prop prop)
func PropTranslator_SerializeProp(prop Node_Prop) *j2gort.Ext_com_google_javascript_jscomp_serialization_NodeProperty {
	clinit_PropTranslator()
	_ = prop
	switch (prop).Ordinal() {
	case int32(20) /* Prop.ARROW_FN */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeProperty_F_ARROW_FN
	case int32(21) /* Prop.ASYNC_FN */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeProperty_F_ASYNC_FN
	case int32(19) /* Prop.GENERATOR_FN */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeProperty_F_GENERATOR_FN
	case int32(22) /* Prop.YIELD_ALL */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeProperty_F_YIELD_ALL
	case int32(0) /* Prop.IS_PARENTHESIZED */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeProperty_F_IS_PARENTHESIZED
	case int32(6) /* Prop.SYNTHETIC */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeProperty_F_SYNTHETIC
	case int32(7) /* Prop.ADDED_BLOCK */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeProperty_F_ADDED_BLOCK
	case int32(18) /* Prop.STATIC_MEMBER */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeProperty_F_STATIC_MEMBER
	case int32(26) /* Prop.IS_GENERATOR_MARKER */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeProperty_F_IS_GENERATOR_MARKER
	case int32(27) /* Prop.IS_GENERATOR_SAFE */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeProperty_F_IS_GENERATOR_SAFE
	case int32(34) /* Prop.COLOR_FROM_CAST */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeProperty_F_COLOR_FROM_CAST
	case int32(40) /* Prop.NON_INDEXABLE */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeProperty_F_NON_INDEXABLE
	case int32(44) /* Prop.DELETED */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeProperty_F_DELETED
	case int32(46) /* Prop.IS_UNUSED_PARAMETER */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeProperty_F_IS_UNUSED_PARAMETER
	case int32(48) /* Prop.IS_SHORTHAND_PROPERTY */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeProperty_F_IS_SHORTHAND_PROPERTY
	case int32(51) /* Prop.START_OF_OPT_CHAIN */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeProperty_F_START_OF_OPT_CHAIN
	case int32(52) /* Prop.TRAILING_COMMA */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeProperty_F_TRAILING_COMMA
	case int32(9) /* Prop.IS_CONSTANT_NAME */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeProperty_F_IS_CONSTANT_NAME
	case int32(10) /* Prop.IS_NAMESPACE */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeProperty_F_IS_NAMESPACE
	case int32(12) /* Prop.DIRECT_EVAL */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeProperty_F_DIRECT_EVAL
	case int32(13) /* Prop.FREE_CALL */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeProperty_F_FREE_CALL
	case int32(17) /* Prop.REFLECTED_OBJECT */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeProperty_F_REFLECTED_OBJECT
	case int32(23) /* Prop.EXPORT_DEFAULT */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeProperty_F_EXPORT_DEFAULT
	case int32(24) /* Prop.EXPORT_ALL_FROM */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeProperty_F_EXPORT_ALL_FROM
	case int32(28) /* Prop.COMPUTED_PROP_METHOD */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeProperty_F_COMPUTED_PROP_METHOD
	case int32(29) /* Prop.COMPUTED_PROP_GETTER */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeProperty_F_COMPUTED_PROP_GETTER
	case int32(30) /* Prop.COMPUTED_PROP_SETTER */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeProperty_F_COMPUTED_PROP_SETTER
	case int32(31) /* Prop.COMPUTED_PROP_VARIABLE */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeProperty_F_COMPUTED_PROP_VARIABLE
	case int32(41) /* Prop.GOOG_MODULE */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeProperty_F_GOOG_MODULE
	case int32(45) /* Prop.MODULE_ALIAS */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeProperty_F_MODULE_ALIAS
	case int32(47) /* Prop.MODULE_EXPORT */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeProperty_F_MODULE_EXPORT
	case int32(49) /* Prop.ES6_MODULE */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeProperty_F_ES6_MODULE
	case int32(25) /* Prop.CONSTANT_VAR_FLAGS */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeProperty_F_CONSTANT_VAR_FLAGS
	case int32(53) /* Prop.SYNTHESIZED_UNFULFILLED_NAME_DECLARATION */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeProperty_F_SYNTHESIZED_UNFULFILLED_NAME_DECLARATION
	case int32(55) /* Prop.PRIVATE_IDENTIFIER */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeProperty_F_PRIVATE_IDENTIFIER
	case int32(54) /* Prop.CLOSURE_UNAWARE_SHADOW */, int32(56) /* Prop.IS_IN_CLOSURE_UNAWARE_SUBTREE */, int32(8) /* Prop.SIDE_EFFECT_FLAGS */, int32(32) /* Prop.DECLARED_TYPE_EXPR */, int32(42) /* Prop.FEATURE_SET */, int32(33) /* Prop.TYPE_BEFORE_CAST */, int32(1) /* Prop.NON_JSDOC_COMMENT */, int32(2) /* Prop.TRAILING_NON_JSDOC_COMMENT */, int32(3) /* Prop.JSDOC_INFO */, int32(4) /* Prop.INCRDECR */, int32(5) /* Prop.QUOTED */, int32(11) /* Prop.USE_STRICT */, int32(14) /* Prop.SOURCE_FILE */, int32(15) /* Prop.INPUT_ID */, int32(16) /* Prop.CHANGE_TIME */, int32(35) /* Prop.OPT_ES6_TYPED */, int32(36) /* Prop.GENERIC_TYPE */, int32(37) /* Prop.IMPLEMENTS */, int32(38) /* Prop.CONSTRUCT_SIGNATURE */, int32(39) /* Prop.ACCESS_MODIFIER */, int32(43) /* Prop.IS_TYPESCRIPT_ABSTRACT */, int32(50) /* Prop.TYPEDEF_TYPE */ :
		return nil
	}
	panic("unreachable switch expression")
}

// Java static method: PropTranslator.checkUnexpectedNullProtoProps()
func PropTranslator_CheckUnexpectedNullProtoProps() {
	clinit_PropTranslator()
	for _, el_6 := range j2gort.CheckNotNullSlice(j2gort.Ext_com_google_javascript_jscomp_serialization_NodeProperty_Values()) {
		protoProp := el_6
		_ = protoProp
		switch j2gort.CheckNotNullPtr(protoProp).Ordinal() {
		case int32(0) /* NodeProperty.NODE_PROPERTY_UNSPECIFIED */ :
		case int32(36) /* NodeProperty.IS_DECLARED_CONSTANT */, int32(35) /* NodeProperty.IS_INFERRED_CONSTANT */ :
		case int32(44) /* NodeProperty.UNRECOGNIZED */, int32(8) /* NodeProperty.UNUSED_11 */ :
		case int32(38) /* NodeProperty.MUTATES_GLOBAL_STATE */, int32(39) /* NodeProperty.MUTATES_THIS */, int32(40) /* NodeProperty.MUTATES_ARGUMENTS */, int32(41) /* NodeProperty.THROWS */ :
		case int32(42) /* NodeProperty.CLOSURE_UNAWARE_SHADOW */ :
		default:
			if !(StaticGet_PropTranslator_ProtoToRhinoProp()[j2gort.EnumOrdinal(protoProp)] != 0) {
				j2gort.CheckStateFmt(false, "Hit unhandled node property: %s", j2gort.PtrToAny(protoProp))
			}
		}
	}
}
