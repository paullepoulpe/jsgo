/*
 * Copyright 2014 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/GatherExternProperties.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.GatherExternProperties
type GatherExternProperties struct {
	compiler         AbstractCompiler          `gson:"compiler"`
	externProperties j2gort.OrderedSet[string] `gson:"externProperties"`
	idHash_          atomic.Int32
	mode             GatherExternProperties_Mode `gson:"mode"`
}

func (this_ *GatherExternProperties) BindSelf_(s any) {
	_ = s
}
func (this_ *GatherExternProperties) IsJavaNil_() bool { return this_ == nil }
func (*GatherExternProperties) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.GatherExternProperties", "GatherExternProperties"
}
func (this_ *GatherExternProperties) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*GatherExternProperties) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.GatherExternProperties", "com.google.javascript.jscomp.GatherExternProperties", func(v_ any) bool {
		return j2gort.InstanceOf[*GatherExternProperties](v_)
	})
}
func (*GatherExternProperties) Implements_jscomp_NodeTraversal_Callback() {}
func (*GatherExternProperties) Implements_jscomp_CompilerPass()           {}

func (this_ *GatherExternProperties) initFields_GatherExternProperties() {
	_ = this_
	this_.externProperties = j2gort.NewOrderedSet[string]()
}

// Java constructor: GatherExternProperties(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.GatherExternProperties.Mode mode)
func NewGatherExternProperties(compiler AbstractCompiler, mode GatherExternProperties_Mode) *GatherExternProperties {
	this_ := &GatherExternProperties{}
	InitGatherExternProperties(this_, compiler, mode)
	return this_
}

func InitGatherExternProperties(this_ *GatherExternProperties, compiler AbstractCompiler, mode GatherExternProperties_Mode) {
	_ = compiler
	_ = mode
	this_.initFields_GatherExternProperties()
	this_.compiler = compiler
	this_.mode = mode
	if compiler.GetExternProperties() != nil {
		j2gort.OrderedSetAddAllStr(&this_.externProperties, j2gort.AsIface[j2gort.Ext_java_util_Collection](compiler.GetExternProperties()))
	}
}

// Java method: GatherExternProperties.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *GatherExternProperties) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_Traverse(this_.compiler, externs, this_)
	this_.compiler.SetExternProperties(this_.externProperties.ToImmutableSet())
}

// Java method: GatherExternProperties.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *GatherExternProperties) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	return (!(n.IsScript()) || !(NodeUtil_IsFromTypeSummary(n)))
}

// Java method: GatherExternProperties.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *GatherExternProperties) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if (this_.mode).FieldGet_optimize() {
		switch n.TokenSwitchOrd_() {
		case int32(26) /* Token.GETPROP */ :
			this_.externProperties.Add(n.GetString())
		case int32(100) /* Token.STRING_KEY */ :
			if parent.IsObjectLit() {
				this_.externProperties.Add(n.GetString())
			}
		case int32(107) /* Token.MEMBER_FUNCTION_DEF */, int32(108) /* Token.MEMBER_FIELD_DEF */, int32(95) /* Token.GETTER_DEF */, int32(96) /* Token.SETTER_DEF */ :
			this_.externProperties.Add(n.GetString())
		default:
		}
	}
	var jsDocInfo *pkg_rhino.JSDocInfo = n.GetJSDocInfo()
	_ = jsDocInfo
	if (jsDocInfo != nil) && (this_.mode).FieldGet_check() {
		this_.priv_GatherExternProperties_gatherPropertiesFromJSDocInfo_3003(jsDocInfo)
	}
}

// Java method: GatherExternProperties.gatherPropertiesFromJSDocInfo(com.google.javascript.rhino.JSDocInfo jsDocInfo)
func (this_ *GatherExternProperties) priv_GatherExternProperties_gatherPropertiesFromJSDocInfo_3003(jsDocInfo *pkg_rhino.JSDocInfo) {
	_ = jsDocInfo
	{
		var i int32 = int32(0)
		_ = i
		var count int32 = jsDocInfo.GetTypeNodeCount()
		_ = count
		for ; i < count; i = i + 1 {
			this_.priv_GatherExternProperties_gatherPropertiesFromJsTypeExpressionNode_3004(jsDocInfo.GetTypeNodeAt(i))
		}
	}
}

// Java method: GatherExternProperties.gatherPropertiesFromJsTypeExpressionNode(com.google.javascript.rhino.Node jsTypeExpressionNode)
func (this_ *GatherExternProperties) priv_GatherExternProperties_gatherPropertiesFromJsTypeExpressionNode_3004(jsTypeExpressionNode pkg_rhino.Node) {
	_ = jsTypeExpressionNode
	switch jsTypeExpressionNode.TokenSwitchOrd_() {
	case int32(161) /* Token.LB */ :
		this_.priv_GatherExternProperties_gatherPropertiesFromJsDocRecordType_3005(jsTypeExpressionNode)
	default:
		{
			var child pkg_rhino.Node = jsTypeExpressionNode.GetFirstChild()
			_ = child
			for ; child != 0; child = child.GetNext() {
				this_.priv_GatherExternProperties_gatherPropertiesFromJsTypeExpressionNode_3004(child)
			}
		}
	}
}

// Java method: GatherExternProperties.gatherPropertiesFromJsDocRecordType(com.google.javascript.rhino.Node jsDocRecordNode)
func (this_ *GatherExternProperties) priv_GatherExternProperties_gatherPropertiesFromJsDocRecordType_3005(jsDocRecordNode pkg_rhino.Node) {
	_ = jsDocRecordNode
	if !(jsDocRecordNode.RawToken_() == uint8(161)) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(jsDocRecordNode))
	}
	{
		var fieldNode pkg_rhino.Node = jsDocRecordNode.GetFirstChild()
		_ = fieldNode
		for ; fieldNode != 0; fieldNode = fieldNode.GetNext() {
			var fieldNameNode pkg_rhino.Node = pkg_rhino.Node(0)
			_ = fieldNameNode
			var fieldTypeNode pkg_rhino.Node = pkg_rhino.Node(0)
			_ = fieldTypeNode
			if fieldNode.RawToken_() == uint8(163) {
				fieldNameNode = fieldNode.GetFirstChild()
				fieldTypeNode = fieldNameNode.GetNext()
			} else {
				fieldNameNode = fieldNode
				fieldTypeNode = pkg_rhino.Node(0)
			}
			if !(fieldNameNode.IsStringKey()) {
				j2gort.CheckStateMsg(false, j2gort.BoxNode(fieldNameNode))
			}
			var fieldName string = fieldNameNode.GetString()
			_ = fieldName
			if javalang.StringStartsWith(j2gort.CheckNotNullString(fieldName), "'") || javalang.StringStartsWith(j2gort.CheckNotNullString(fieldName), "\"") {
				fieldName = j2gort.JStrSub(fieldName, int32(1), (j2gort.JStrLen(fieldName) - int32(1)))
			}
			this_.externProperties.Add(fieldName)
			if fieldTypeNode != 0 {
				this_.priv_GatherExternProperties_gatherPropertiesFromJsTypeExpressionNode_3004(fieldTypeNode)
			}
		}
	}
}

// Java enum: com.google.javascript.jscomp.GatherExternProperties.Mode
type GatherExternProperties_Mode uint8

const (
	enum_GatherExternProperties_Mode_CHECK              GatherExternProperties_Mode = 1
	enum_GatherExternProperties_Mode_OPTIMIZE           GatherExternProperties_Mode = 2
	enum_GatherExternProperties_Mode_CHECK_AND_OPTIMIZE GatherExternProperties_Mode = 3
)

var enumBaseTable_GatherExternProperties_Mode = [...]j2gort.EnumBase{
	{Name_: "CHECK", Ordinal_: 0},
	{Name_: "OPTIMIZE", Ordinal_: 1},
	{Name_: "CHECK_AND_OPTIMIZE", Ordinal_: 2},
}
var enumValues_GatherExternProperties_Mode = []GatherExternProperties_Mode{1, 2, 3}

var enumField_GatherExternProperties_Mode_check = [...]bool{
	true,
	false,
	true,
}

func (this_ GatherExternProperties_Mode) FieldGet_check() bool {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumField_GatherExternProperties_Mode_check[this_-1]
}

var enumField_GatherExternProperties_Mode_optimize = [...]bool{
	false,
	true,
	true,
}

func (this_ GatherExternProperties_Mode) FieldGet_optimize() bool {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumField_GatherExternProperties_Mode_optimize[this_-1]
}

func StaticGet_GatherExternProperties_Mode_CHECK() GatherExternProperties_Mode {
	return enum_GatherExternProperties_Mode_CHECK
}
func StaticGet_GatherExternProperties_Mode_OPTIMIZE() GatherExternProperties_Mode {
	return enum_GatherExternProperties_Mode_OPTIMIZE
}
func StaticGet_GatherExternProperties_Mode_CHECK_AND_OPTIMIZE() GatherExternProperties_Mode {
	return enum_GatherExternProperties_Mode_CHECK_AND_OPTIMIZE
}

func GatherExternProperties_Mode_Values() []GatherExternProperties_Mode {
	return j2gort.CloneSlice(enumValues_GatherExternProperties_Mode)
}

func GatherExternProperties_Mode_ValueOf(name string) GatherExternProperties_Mode {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_GatherExternProperties_Mode {
		if enumBaseTable_GatherExternProperties_Mode[i_].Name_ == name {
			return GatherExternProperties_Mode(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.GatherExternProperties.Mode." + name))
}

func (this_ GatherExternProperties_Mode) IsJavaNil_() bool                              { return this_ == 0 }
func (GatherExternProperties_Mode) HasIdentityEquals_()                                 {}
func (GatherExternProperties_Mode) Implements_j2gort_Ext_java_lang_Enum()               {}
func (GatherExternProperties_Mode) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ GatherExternProperties_Mode) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (GatherExternProperties_Mode) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.GatherExternProperties$Mode", "Mode"
}
func (this_ GatherExternProperties_Mode) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_GatherExternProperties_Mode[this_-1]
}
func (this_ GatherExternProperties_Mode) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ GatherExternProperties_Mode) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ GatherExternProperties_Mode) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_GatherExternProperties_Mode[this_-1].Name_
}
func (this_ GatherExternProperties_Mode) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_GatherExternProperties_Mode[this_-1].Name_
}
func (this_ GatherExternProperties_Mode) Equals(other_ any) bool {
	o_, ok_ := other_.(GatherExternProperties_Mode)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ GatherExternProperties_Mode) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_GatherExternProperties_Mode[this_-1].HashCode()
}
func (this_ GatherExternProperties_Mode) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_GatherExternProperties_Mode[this_-1].CompareTo(other_)
}
func (this_ GatherExternProperties_Mode) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_GatherExternProperties_Mode[this_-1].Name_)
}
func init() {
	vals_ := []any{GatherExternProperties_Mode(1), GatherExternProperties_Mode(2), GatherExternProperties_Mode(3)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.GatherExternProperties.Mode", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.GatherExternProperties$Mode", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.GatherExternProperties.Mode", "com.google.javascript.jscomp.GatherExternProperties$Mode", func(v_ any) bool {
		return j2gort.InstanceOf[GatherExternProperties_Mode](v_)
	})
}
func (*GatherExternProperties) CanReleaseSyntacticScopes_() {}
