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
 *   Michael Zhou
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */
// Port of com/google/javascript/rhino/TypeDeclarationsIR.java (generated mechanically by j2go).

package rhino

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.rhino.TypeDeclarationsIR
type TypeDeclarationsIR struct {
	idHash_ atomic.Int32
}

func (this_ *TypeDeclarationsIR) BindSelf_(s any) {
	_ = s
}
func (this_ *TypeDeclarationsIR) IsJavaNil_() bool { return this_ == nil }
func (*TypeDeclarationsIR) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.TypeDeclarationsIR", "TypeDeclarationsIR"
}
func (this_ *TypeDeclarationsIR) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*TypeDeclarationsIR) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.TypeDeclarationsIR", "com.google.javascript.rhino.TypeDeclarationsIR", func(v_ any) bool {
		return j2gort.InstanceOf[*TypeDeclarationsIR](v_)
	})
}

func (this_ *TypeDeclarationsIR) initFields_TypeDeclarationsIR() {
	_ = this_
}

// Java constructor: TypeDeclarationsIR()
func NewTypeDeclarationsIR() *TypeDeclarationsIR {
	this_ := &TypeDeclarationsIR{}
	InitTypeDeclarationsIR(this_)
	return this_
}

func InitTypeDeclarationsIR(this_ *TypeDeclarationsIR) {
	this_.initFields_TypeDeclarationsIR()
}

// Java static method: TypeDeclarationsIR.stringType()
func TypeDeclarationsIR_StringType() Node {
	return NewNode_1(StaticGet_Token_STRING_TYPE())
}

// Java static method: TypeDeclarationsIR.numberType()
func TypeDeclarationsIR_NumberType() Node {
	return NewNode_1(StaticGet_Token_NUMBER_TYPE())
}

// Java static method: TypeDeclarationsIR.booleanType()
func TypeDeclarationsIR_BooleanType() Node {
	return NewNode_1(StaticGet_Token_BOOLEAN_TYPE())
}

// Java static method: TypeDeclarationsIR.anyType()
func TypeDeclarationsIR_AnyType() Node {
	return NewNode_1(StaticGet_Token_ANY_TYPE())
}

// Java static method: TypeDeclarationsIR.voidType()
func TypeDeclarationsIR_VoidType() Node {
	return NewNode_1(StaticGet_Token_VOID_TYPE())
}

// Java static method: TypeDeclarationsIR.undefinedType()
func TypeDeclarationsIR_UndefinedType() Node {
	return NewNode_1(StaticGet_Token_UNDEFINED_TYPE())
}

// Java static method: TypeDeclarationsIR.namedType(java.lang.String typeName)
func TypeDeclarationsIR_NamedType_1_string(typeName string) Node {
	_ = typeName
	return TypeDeclarationsIR_NamedType_1_Ext_java_lang_Iterable(j2gort.SplitterSplitStr(j2gort.Ext_com_google_common_base_Splitter_On_1_uint16(uint16(46)), typeName))
}

// Java static method: TypeDeclarationsIR.namedType(java.lang.Iterable<java.lang.String> segments)
func TypeDeclarationsIR_NamedType_1_Ext_java_lang_Iterable(segments j2gort.Ext_java_lang_Iterable) Node {
	_ = segments
	var segmentsIt j2gort.Ext_java_util_Iterator = j2gort.AsIface[j2gort.Ext_java_util_Iterator](segments.Iterator())
	_ = segmentsIt
	var node Node = IR_Name(j2gort.CastString(segmentsIt.Next()))
	_ = node
	for segmentsIt.HasNext() {
		node = IR_Getprop_2(node, j2gort.CastString(segmentsIt.Next()))
	}
	return NewNode_2(StaticGet_Token_NAMED_TYPE(), node)
}

// Java static method: TypeDeclarationsIR.recordType(java.util.LinkedHashMap<java.lang.String,com.google.javascript.rhino.Node> properties)
func TypeDeclarationsIR_RecordType(properties *j2gort.Ext_java_util_LinkedHashMap) Node {
	_ = properties
	var node Node = NewNode_1(StaticGet_Token_RECORD_TYPE())
	_ = node
	{
		cur_1 := j2gort.NewMapEntryCursor(properties)
		for cur_1.NextEntry() {
			var stringKey Node = IR_StringKey_1(cur_1.KeyStr())
			_ = stringKey
			node.AddChildToBack(stringKey)
			if j2gort.CastNode[Node](cur_1.GetValue()) != 0 {
				stringKey.AddChildToFront(j2gort.CastNode[Node](cur_1.GetValue()))
			}
		}
	}
	return node
}

// Java static method: TypeDeclarationsIR.maybeAddType(com.google.javascript.rhino.Node node, com.google.javascript.rhino.Node type)
func TypeDeclarationsIR_MaybeAddType(node Node, type_ Node) Node {
	_ = node
	_ = type_
	if type_ != 0 {
		node.SetDeclaredTypeExpression(type_)
	}
	return node
}

// Java static method: TypeDeclarationsIR.functionType(com.google.javascript.rhino.Node returnType, java.util.LinkedHashMap<java.lang.String,com.google.javascript.rhino.Node> requiredParams, java.util.LinkedHashMap<java.lang.String,com.google.javascript.rhino.Node> optionalParams, java.lang.String restName, com.google.javascript.rhino.Node restType)
func TypeDeclarationsIR_FunctionType(returnType Node, requiredParams *j2gort.Ext_java_util_LinkedHashMap, optionalParams *j2gort.Ext_java_util_LinkedHashMap, restName string, restType Node) Node {
	_ = returnType
	_ = requiredParams
	_ = optionalParams
	_ = restName
	_ = restType
	var node Node = NewNode_2(StaticGet_Token_FUNCTION_TYPE(), returnType)
	_ = node
	j2gort.CheckNotNullPtr(requiredParams)
	j2gort.CheckNotNullPtr(optionalParams)
	{
		cur_2 := j2gort.NewMapEntryCursor(requiredParams)
		for cur_2.NextEntry() {
			var name Node = IR_Name(cur_2.KeyStr())
			_ = name
			node.AddChildToBack(TypeDeclarationsIR_MaybeAddType(name, j2gort.CastNode[Node](cur_2.GetValue())))
		}
	}
	{
		cur_3 := j2gort.NewMapEntryCursor(optionalParams)
		for cur_3.NextEntry() {
			var name_2 Node = IR_Name(cur_3.KeyStr())
			_ = name_2
			name_2.PutBooleanProp(StaticGet_Node_OPT_ES6_TYPED(), true)
			node.AddChildToBack(TypeDeclarationsIR_MaybeAddType(name_2, j2gort.CastNode[Node](cur_3.GetValue())))
		}
	}
	if restName != j2gort.NullString {
		var rest Node = NewNode_2(StaticGet_Token_ITER_REST(), IR_Name(restName))
		_ = rest
		node.AddChildToBack(TypeDeclarationsIR_MaybeAddType(rest, restType))
	}
	return node
}

// Java static method: TypeDeclarationsIR.parameterizedType(com.google.javascript.rhino.Node baseType, java.lang.Iterable<com.google.javascript.rhino.Node> typeParameters)
func TypeDeclarationsIR_ParameterizedType(baseType Node, typeParameters j2gort.Ext_java_lang_Iterable) Node {
	_ = baseType
	_ = typeParameters
	if j2gort.Ext_com_google_common_collect_Iterables_IsEmpty(typeParameters) {
		return baseType
	}
	var node Node = NewNode_2(StaticGet_Token_PARAMETERIZED_TYPE(), baseType)
	_ = node
	{
		cur_4 := j2gort.NewIterCursor(j2gort.CheckNotNull(typeParameters))
		for cur_4.HasNext() {
			typeParameter := Node(cur_4.NextNode())
			_ = typeParameter
			node.AddChildToBack(typeParameter)
		}
	}
	return node
}

// Java static method: TypeDeclarationsIR.arrayType(com.google.javascript.rhino.Node elementType)
func TypeDeclarationsIR_ArrayType(elementType Node) Node {
	_ = elementType
	return NewNode_2(StaticGet_Token_ARRAY_TYPE(), elementType)
}

// Java static method: TypeDeclarationsIR.unionType(java.lang.Iterable<com.google.javascript.rhino.Node> options)
func TypeDeclarationsIR_UnionType_1_Ext_java_lang_Iterable(options j2gort.Ext_java_lang_Iterable) Node {
	_ = options
	if !(!(j2gort.Ext_com_google_common_collect_Iterables_IsEmpty(options))) {
		j2gort.CheckArgumentStr(false, "union must have at least one option")
	}
	var node Node = NewNode_1(StaticGet_Token_UNION_TYPE())
	_ = node
	{
		cur_5 := j2gort.NewIterCursor(j2gort.CheckNotNull(options))
		for cur_5.HasNext() {
			option := Node(cur_5.NextNode())
			_ = option
			node.AddChildToBack(option)
		}
	}
	return node
}

// Java static method: TypeDeclarationsIR.unionType(com.google.javascript.rhino.Node[] options)
func TypeDeclarationsIR_UnionType_1_Node(options []Node) Node {
	_ = options
	var iifeArg_1 []any
iifeSwitch_1:
	switch {
	default:

		arr_6 := options
		if arr_6 == nil {
			iifeArg_1 = nil
			break iifeSwitch_1
		}
		out_7 := make([]any, len(arr_6))
		for i_, el_8 := range arr_6 {
			out_7[i_] = j2gort.BoxNode(el_8)
		}
		iifeArg_1 = out_7
		break iifeSwitch_1
	}
	return TypeDeclarationsIR_UnionType_1_Ext_java_lang_Iterable(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](j2gort.Ext_java_util_Arrays_AsList(iifeArg_1)))
}

// Java static method: TypeDeclarationsIR.optionalParameter(com.google.javascript.rhino.Node parameterType)
func TypeDeclarationsIR_OptionalParameter(parameterType Node) Node {
	_ = parameterType
	return NewNode_2(StaticGet_Token_OPTIONAL_PARAMETER(), parameterType)
}
