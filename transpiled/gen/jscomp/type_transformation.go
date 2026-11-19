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
// Port of com/google/javascript/jscomp/TypeTransformation.java (generated mechanically by j2go).

package jscomp

import (
	"strconv"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parsing "github.com/paullepoulpe/jsgo/transpiled/gen/parsing"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.TypeTransformation
type TypeTransformation struct {
	compiler AbstractCompiler           `gson:"compiler"`
	registry *pkg_rhino.JSTypeRegistry  `gson:"registry"`
	typeEnv  pkg_rhino.StaticTypedScope `gson:"typeEnv"`
	idHash_  atomic.Int32
}

func (this_ *TypeTransformation) BindSelf_(s any) {
	_ = s
}
func (this_ *TypeTransformation) IsJavaNil_() bool { return this_ == nil }
func (*TypeTransformation) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.TypeTransformation", "TypeTransformation"
}
func (this_ *TypeTransformation) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*TypeTransformation) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.TypeTransformation", "com.google.javascript.jscomp.TypeTransformation", func(v_ any) bool {
		return j2gort.InstanceOf[*TypeTransformation](v_)
	})
}

func StaticGet_TypeTransformation_VIRTUAL_FILE() string { return "<TypeTransformation.java>" }

var static_TypeTransformation_UNKNOWN_TYPEVAR *DiagnosticType

func StaticGet_TypeTransformation_UNKNOWN_TYPEVAR() *DiagnosticType {
	clinit_TypeTransformation()
	return static_TypeTransformation_UNKNOWN_TYPEVAR
}
func StaticSet_TypeTransformation_UNKNOWN_TYPEVAR(v *DiagnosticType) {
	clinit_TypeTransformation()
	static_TypeTransformation_UNKNOWN_TYPEVAR = v
}

var static_TypeTransformation_UNKNOWN_STRVAR *DiagnosticType

func StaticGet_TypeTransformation_UNKNOWN_STRVAR() *DiagnosticType {
	clinit_TypeTransformation()
	return static_TypeTransformation_UNKNOWN_STRVAR
}
func StaticSet_TypeTransformation_UNKNOWN_STRVAR(v *DiagnosticType) {
	clinit_TypeTransformation()
	static_TypeTransformation_UNKNOWN_STRVAR = v
}

var static_TypeTransformation_UNKNOWN_TYPENAME *DiagnosticType

func StaticGet_TypeTransformation_UNKNOWN_TYPENAME() *DiagnosticType {
	clinit_TypeTransformation()
	return static_TypeTransformation_UNKNOWN_TYPENAME
}
func StaticSet_TypeTransformation_UNKNOWN_TYPENAME(v *DiagnosticType) {
	clinit_TypeTransformation()
	static_TypeTransformation_UNKNOWN_TYPENAME = v
}

var static_TypeTransformation_BASETYPE_INVALID *DiagnosticType

func StaticGet_TypeTransformation_BASETYPE_INVALID() *DiagnosticType {
	clinit_TypeTransformation()
	return static_TypeTransformation_BASETYPE_INVALID
}
func StaticSet_TypeTransformation_BASETYPE_INVALID(v *DiagnosticType) {
	clinit_TypeTransformation()
	static_TypeTransformation_BASETYPE_INVALID = v
}

var static_TypeTransformation_TEMPTYPE_INVALID *DiagnosticType

func StaticGet_TypeTransformation_TEMPTYPE_INVALID() *DiagnosticType {
	clinit_TypeTransformation()
	return static_TypeTransformation_TEMPTYPE_INVALID
}
func StaticSet_TypeTransformation_TEMPTYPE_INVALID(v *DiagnosticType) {
	clinit_TypeTransformation()
	static_TypeTransformation_TEMPTYPE_INVALID = v
}

var static_TypeTransformation_INDEX_OUTOFBOUNDS *DiagnosticType

func StaticGet_TypeTransformation_INDEX_OUTOFBOUNDS() *DiagnosticType {
	clinit_TypeTransformation()
	return static_TypeTransformation_INDEX_OUTOFBOUNDS
}
func StaticSet_TypeTransformation_INDEX_OUTOFBOUNDS(v *DiagnosticType) {
	clinit_TypeTransformation()
	static_TypeTransformation_INDEX_OUTOFBOUNDS = v
}

var static_TypeTransformation_DUPLICATE_VARIABLE *DiagnosticType

func StaticGet_TypeTransformation_DUPLICATE_VARIABLE() *DiagnosticType {
	clinit_TypeTransformation()
	return static_TypeTransformation_DUPLICATE_VARIABLE
}
func StaticSet_TypeTransformation_DUPLICATE_VARIABLE(v *DiagnosticType) {
	clinit_TypeTransformation()
	static_TypeTransformation_DUPLICATE_VARIABLE = v
}

var static_TypeTransformation_UNKNOWN_NAMEVAR *DiagnosticType

func StaticGet_TypeTransformation_UNKNOWN_NAMEVAR() *DiagnosticType {
	clinit_TypeTransformation()
	return static_TypeTransformation_UNKNOWN_NAMEVAR
}
func StaticSet_TypeTransformation_UNKNOWN_NAMEVAR(v *DiagnosticType) {
	clinit_TypeTransformation()
	static_TypeTransformation_UNKNOWN_NAMEVAR = v
}

var static_TypeTransformation_RECTYPE_INVALID *DiagnosticType

func StaticGet_TypeTransformation_RECTYPE_INVALID() *DiagnosticType {
	clinit_TypeTransformation()
	return static_TypeTransformation_RECTYPE_INVALID
}
func StaticSet_TypeTransformation_RECTYPE_INVALID(v *DiagnosticType) {
	clinit_TypeTransformation()
	static_TypeTransformation_RECTYPE_INVALID = v
}

var static_TypeTransformation_MAPRECORD_BODY_INVALID *DiagnosticType

func StaticGet_TypeTransformation_MAPRECORD_BODY_INVALID() *DiagnosticType {
	clinit_TypeTransformation()
	return static_TypeTransformation_MAPRECORD_BODY_INVALID
}
func StaticSet_TypeTransformation_MAPRECORD_BODY_INVALID(v *DiagnosticType) {
	clinit_TypeTransformation()
	static_TypeTransformation_MAPRECORD_BODY_INVALID = v
}

var static_TypeTransformation_VAR_UNDEFINED *DiagnosticType

func StaticGet_TypeTransformation_VAR_UNDEFINED() *DiagnosticType {
	clinit_TypeTransformation()
	return static_TypeTransformation_VAR_UNDEFINED
}
func StaticSet_TypeTransformation_VAR_UNDEFINED(v *DiagnosticType) {
	clinit_TypeTransformation()
	static_TypeTransformation_VAR_UNDEFINED = v
}

var static_TypeTransformation_INVALID_CTOR *DiagnosticType

func StaticGet_TypeTransformation_INVALID_CTOR() *DiagnosticType {
	clinit_TypeTransformation()
	return static_TypeTransformation_INVALID_CTOR
}
func StaticSet_TypeTransformation_INVALID_CTOR(v *DiagnosticType) {
	clinit_TypeTransformation()
	static_TypeTransformation_INVALID_CTOR = v
}

var static_TypeTransformation_RECPARAM_INVALID *DiagnosticType

func StaticGet_TypeTransformation_RECPARAM_INVALID() *DiagnosticType {
	clinit_TypeTransformation()
	return static_TypeTransformation_RECPARAM_INVALID
}
func StaticSet_TypeTransformation_RECPARAM_INVALID(v *DiagnosticType) {
	clinit_TypeTransformation()
	static_TypeTransformation_RECPARAM_INVALID = v
}

var static_TypeTransformation_PROPTYPE_INVALID *DiagnosticType

func StaticGet_TypeTransformation_PROPTYPE_INVALID() *DiagnosticType {
	clinit_TypeTransformation()
	return static_TypeTransformation_PROPTYPE_INVALID
}
func StaticSet_TypeTransformation_PROPTYPE_INVALID(v *DiagnosticType) {
	clinit_TypeTransformation()
	static_TypeTransformation_PROPTYPE_INVALID = v
}

var clinitDone_TypeTransformation uint32
var clinitState_TypeTransformation j2gort.ClinitState

func clinit_TypeTransformation() {
	if atomic.LoadUint32(&clinitDone_TypeTransformation) != 0 {
		return
	}
	clinitSlow_TypeTransformation()
}

func clinitSlow_TypeTransformation() {
	if !j2gort.EnterClinit(&clinitState_TypeTransformation, &clinitDone_TypeTransformation) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_TypeTransformation, &clinitDone_TypeTransformation)
	static_TypeTransformation_UNKNOWN_TYPEVAR = DiagnosticType_Warning("TYPEVAR_UNDEFINED", "Reference to an unknown type variable {0}")
	static_TypeTransformation_UNKNOWN_STRVAR = DiagnosticType_Warning("UNKNOWN_STRVAR", "Reference to an unknown string variable {0}")
	static_TypeTransformation_UNKNOWN_TYPENAME = DiagnosticType_Warning("TYPENAME_UNDEFINED", "Reference to an unknown type name {0}")
	static_TypeTransformation_BASETYPE_INVALID = DiagnosticType_Warning("BASETYPE_INVALID", "The type {0} cannot be templatized")
	static_TypeTransformation_TEMPTYPE_INVALID = DiagnosticType_Warning("TEMPTYPE_INVALID", "Expected templatized type in {0} found {1}")
	static_TypeTransformation_INDEX_OUTOFBOUNDS = DiagnosticType_Warning("INDEX_OUTOFBOUNDS", "Index out of bounds in templateTypeOf: expected a number less than {0}, found {1}")
	static_TypeTransformation_DUPLICATE_VARIABLE = DiagnosticType_Warning("DUPLICATE_VARIABLE", "The variable {0} is already defined")
	static_TypeTransformation_UNKNOWN_NAMEVAR = DiagnosticType_Warning("UNKNOWN_NAMEVAR", "Reference to an unknown name variable {0}")
	static_TypeTransformation_RECTYPE_INVALID = DiagnosticType_Warning("RECTYPE_INVALID", "The first parameter of a maprecord must be a record type, found {0}")
	static_TypeTransformation_MAPRECORD_BODY_INVALID = DiagnosticType_Warning("MAPRECORD_BODY_INVALID", "The body of a maprecord function must evaluate to a record type or a no type, found {0}")
	static_TypeTransformation_VAR_UNDEFINED = DiagnosticType_Warning("VAR_UNDEFINED", "Variable {0} is undefined in the scope")
	static_TypeTransformation_INVALID_CTOR = DiagnosticType_Warning("INVALID_CTOR", "Expected a constructor type, found {0}")
	static_TypeTransformation_RECPARAM_INVALID = DiagnosticType_Warning("RECPARAM_INVALID", "Expected a record type, found {0}")
	static_TypeTransformation_PROPTYPE_INVALID = DiagnosticType_Warning("PROPTYPE_INVALID", "Expected object type, found {0}")
}

func (this_ *TypeTransformation) initFields_TypeTransformation() {
	_ = this_
}

// Java constructor: TypeTransformation(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.rhino.jstype.StaticTypedScope typeEnv)
func NewTypeTransformation(compiler AbstractCompiler, typeEnv pkg_rhino.StaticTypedScope) *TypeTransformation {
	clinit_TypeTransformation()
	this_ := &TypeTransformation{}
	InitTypeTransformation(this_, compiler, typeEnv)
	return this_
}

func InitTypeTransformation(this_ *TypeTransformation, compiler AbstractCompiler, typeEnv pkg_rhino.StaticTypedScope) {
	_ = compiler
	_ = typeEnv
	this_.initFields_TypeTransformation()
	this_.compiler = compiler
	this_.registry = compiler.GetTypeRegistry()
	this_.typeEnv = typeEnv
}

// Java method: TypeTransformation.isTypeVar(com.google.javascript.rhino.Node n)
func (this_ *TypeTransformation) priv_TypeTransformation_isTypeVar_5593(n pkg_rhino.Node) bool {
	_ = n
	return n.IsName()
}

// Java method: TypeTransformation.isTypeName(com.google.javascript.rhino.Node n)
func (this_ *TypeTransformation) priv_TypeTransformation_isTypeName_5594(n pkg_rhino.Node) bool {
	_ = n
	return n.IsStringLit()
}

// Java method: TypeTransformation.isBooleanOperation(com.google.javascript.rhino.Node n)
func (this_ *TypeTransformation) priv_TypeTransformation_isBooleanOperation_5595(n pkg_rhino.Node) bool {
	_ = n
	return ((n.IsAnd() || n.IsOr()) || n.IsNot())
}

// Java method: TypeTransformation.nameToKeyword(java.lang.String s)
func (this_ *TypeTransformation) priv_TypeTransformation_nameToKeyword_5596(s string) *pkg_parsing.TypeTransformationParser_Keywords {
	_ = s
	return pkg_parsing.TypeTransformationParser_Keywords_ValueOf(j2gort.Ext_com_google_common_base_Ascii_ToUpperCase_1_string(s))
}

// Java method: TypeTransformation.getType(java.lang.String typeName)
func (this_ *TypeTransformation) priv_TypeTransformation_getType_5597(typeName string) pkg_rhino.JSType {
	_ = typeName
	var type_ pkg_rhino.JSType = this_.registry.GetType_2(pkg_rhino.Upcast_StaticScope(this_.typeEnv), typeName)
	_ = type_
	if type_ != nil && !type_.IsJavaNil_() {
		return type_
	}
	var slot pkg_rhino.StaticTypedSlot = j2gort.AsIface[pkg_rhino.StaticTypedSlot](this_.typeEnv.GetSlot(typeName))
	_ = slot
	if slot != nil && !slot.IsJavaNil_() {
		type_ = slot.GetType()
	} else {
		type_ = nil
	}
	if type_ != nil && !type_.IsJavaNil_() {
		if type_.GetJSTypeBase().IsConstructor() || type_.GetJSTypeBase().IsInterface() {
			return pkg_rhino.Upcast_JSType(type_.ToMaybeFunctionType().GetInstanceType().GetObjectTypeBase().GetRawType())
		}
		if type_.GetJSTypeBase().IsEnumElementType() {
			return type_.GetJSTypeBase().GetEnumeratedTypeOfEnumElement()
		}
		return type_
	}
	var jsdoc *pkg_rhino.JSDocInfo
	_ = jsdoc
	if slot == nil || slot.IsJavaNil_() {
		jsdoc = nil
	} else {
		jsdoc = slot.GetJSDocInfo()
	}
	if (jsdoc != nil) && jsdoc.HasTypedefType() {
		return this_.registry.EvaluateTypeExpression(jsdoc.GetTypedefType(), this_.typeEnv)
	}
	return nil
}

// Java method: TypeTransformation.getUnknownType()
func (this_ *TypeTransformation) priv_TypeTransformation_getUnknownType_5598() pkg_rhino.JSType {
	return pkg_rhino.Upcast_JSType(this_.registry.GetNativeObjectType(pkg_rhino.StaticGet_JSTypeNative_UNKNOWN_TYPE()))
}

// Java method: TypeTransformation.getNoType()
func (this_ *TypeTransformation) priv_TypeTransformation_getNoType_5599() pkg_rhino.JSType {
	return pkg_rhino.Upcast_JSType(this_.registry.GetNativeObjectType(pkg_rhino.StaticGet_JSTypeNative_NO_TYPE()))
}

// Java method: TypeTransformation.getAllType()
func (this_ *TypeTransformation) priv_TypeTransformation_getAllType_5600() pkg_rhino.JSType {
	return this_.registry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_ALL_TYPE())
}

// Java method: TypeTransformation.getObjectType()
func (this_ *TypeTransformation) priv_TypeTransformation_getObjectType_5601() pkg_rhino.JSType {
	return this_.registry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_OBJECT_TYPE())
}

// Java method: TypeTransformation.createUnionType(com.google.javascript.rhino.jstype.JSType[] variants)
func (this_ *TypeTransformation) priv_TypeTransformation_createUnionType_5602(variants []pkg_rhino.JSType) pkg_rhino.JSType {
	_ = variants
	return this_.registry.CreateUnionType_1_JSType(variants)
}

// Java method: TypeTransformation.createRecordType(com.google.common.collect.ImmutableMap<java.lang.String,com.google.javascript.rhino.jstype.JSType> props)
func (this_ *TypeTransformation) priv_TypeTransformation_createRecordType_5603(props j2gort.Ext_com_google_common_collect_ImmutableMap) pkg_rhino.JSType {
	_ = props
	return this_.registry.CreateRecordType(j2gort.AsIface[j2gort.Ext_java_util_Map](props))
}

// Java method: TypeTransformation.reportWarning(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.DiagnosticType msg, java.lang.String[] param)
func (this_ *TypeTransformation) priv_TypeTransformation_reportWarning_5604(n pkg_rhino.Node, msg *DiagnosticType, param []string) {
	_ = n
	_ = msg
	_ = param
	this_.compiler.Report_1(JSError_Make_3(n, msg, param))
}

// Java method: TypeTransformation.addNewEntry(com.google.common.collect.ImmutableMap<java.lang.String,T> map, java.lang.String name, T type)
func (this_ *TypeTransformation) priv_TypeTransformation_addNewEntry_5605(map_ j2gort.Ext_com_google_common_collect_ImmutableMap, name string, type_ any) j2gort.Ext_com_google_common_collect_ImmutableMap {
	_ = map_
	_ = name
	_ = type_
	return j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.NewExt_com_google_common_collect_ImmutableMap_Builder().PutAll_1_Ext_java_util_Map(j2gort.AsIface[j2gort.Ext_java_util_Map](map_))).Put_2(j2gort.BoxString(name), type_)).BuildOrThrow())
}

// Java method: TypeTransformation.getFunctionParameter(com.google.javascript.rhino.Node n, int i)
func (this_ *TypeTransformation) priv_TypeTransformation_getFunctionParameter_5606(n pkg_rhino.Node, i int32) string {
	_ = n
	_ = i
	if !(n.IsFunction()) {
		j2gort.CheckArgumentFmt(false, "Expected a function node, found %s", j2gort.BoxNode(n))
	}
	return n.GetSecondChild().GetChildAtIndex(i).GetString()
}

// Java method: TypeTransformation.getCallName(com.google.javascript.rhino.Node n)
func (this_ *TypeTransformation) priv_TypeTransformation_getCallName_5607(n pkg_rhino.Node) string {
	_ = n
	if !(n.IsCall()) {
		j2gort.CheckArgumentFmt(false, "Expected a call node, found %s", j2gort.BoxNode(n))
	}
	return n.GetFirstChild().GetString()
}

// Java method: TypeTransformation.getCallArgument(com.google.javascript.rhino.Node n, int i)
func (this_ *TypeTransformation) priv_TypeTransformation_getCallArgument_5608(n pkg_rhino.Node, i int32) pkg_rhino.Node {
	_ = n
	_ = i
	if !(n.IsCall()) {
		j2gort.CheckArgumentFmt(false, "Expected a call node, found %s", j2gort.BoxNode(n))
	}
	return n.GetChildAtIndex((i + int32(1)))
}

// Java method: TypeTransformation.getCallParamCount(com.google.javascript.rhino.Node n)
func (this_ *TypeTransformation) priv_TypeTransformation_getCallParamCount_5609(n pkg_rhino.Node) int32 {
	_ = n
	if !(n.IsCall()) {
		j2gort.CheckArgumentFmt(false, "Expected a call node, found %s", j2gort.BoxNode(n))
	}
	return (n.GetChildCount() - int32(1))
}

// Java method: TypeTransformation.getCallParams(com.google.javascript.rhino.Node n)
func (this_ *TypeTransformation) priv_TypeTransformation_getCallParams_5610(n pkg_rhino.Node) j2gort.Ext_com_google_common_collect_ImmutableList {
	_ = n
	if !(n.IsCall()) {
		j2gort.CheckArgumentFmt(false, "Expected a call node, found %s", j2gort.BoxNode(n))
	}
	var builder []pkg_rhino.Node = []pkg_rhino.Node{}
	_ = builder
	{
		var i int32 = int32(0)
		_ = i
		for ; i < this_.priv_TypeTransformation_getCallParamCount_5609(n); i = i + 1 {
			builder = append(builder, j2gort.CheckNotNullElem(this_.priv_TypeTransformation_getCallArgument_5608(n, i)))
		}
	}
	return j2gort.ImmListFromTypedSlice(builder)
}

// Java method: TypeTransformation.getComputedPropValue(com.google.javascript.rhino.Node n)
func (this_ *TypeTransformation) priv_TypeTransformation_getComputedPropValue_5611(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	if !(n.IsComputedProp()) {
		j2gort.CheckArgumentFmt(false, "Expected a computed property node, found %s", j2gort.BoxNode(n))
	}
	return n.GetSecondChild()
}

// Java method: TypeTransformation.getComputedPropName(com.google.javascript.rhino.Node n)
func (this_ *TypeTransformation) priv_TypeTransformation_getComputedPropName_5612(n pkg_rhino.Node) string {
	_ = n
	if !(n.IsComputedProp()) {
		j2gort.CheckArgumentFmt(false, "Expected a computed property node, found %s", j2gort.BoxNode(n))
	}
	return n.GetFirstChild().GetString()
}

// Java method: TypeTransformation.eval(com.google.javascript.rhino.Node ttlAst, com.google.common.collect.ImmutableMap<java.lang.String,com.google.javascript.rhino.jstype.JSType> typeVars)
func (this_ *TypeTransformation) Eval_2(ttlAst pkg_rhino.Node, typeVars j2gort.Ext_com_google_common_collect_ImmutableMap) pkg_rhino.JSType {
	_ = ttlAst
	_ = typeVars
	return this_.Eval_3(ttlAst, typeVars, j2gort.Ext_com_google_common_collect_ImmutableMap_Of_0())
}

// Java method: TypeTransformation.eval(com.google.javascript.rhino.Node ttlAst, com.google.common.collect.ImmutableMap<java.lang.String,com.google.javascript.rhino.jstype.JSType> typeVars, com.google.common.collect.ImmutableMap<java.lang.String,java.lang.String> nameVars)
func (this_ *TypeTransformation) Eval_3(ttlAst pkg_rhino.Node, typeVars j2gort.Ext_com_google_common_collect_ImmutableMap, nameVars j2gort.Ext_com_google_common_collect_ImmutableMap) pkg_rhino.JSType {
	_ = ttlAst
	_ = typeVars
	_ = nameVars
	var result pkg_rhino.JSType = this_.priv_TypeTransformation_evalInternal_5615(ttlAst, NewTypeTransformation_NameResolver(typeVars, nameVars))
	_ = result
	if result.GetJSTypeBase().IsEmptyType() {
		return this_.priv_TypeTransformation_getUnknownType_5598()
	} else {
		return result
	}
}

// Java method: TypeTransformation.evalInternal(com.google.javascript.rhino.Node ttlAst, com.google.javascript.jscomp.TypeTransformation.NameResolver nameResolver)
func (this_ *TypeTransformation) priv_TypeTransformation_evalInternal_5615(ttlAst pkg_rhino.Node, nameResolver *TypeTransformation_NameResolver) pkg_rhino.JSType {
	_ = ttlAst
	_ = nameResolver
	if this_.priv_TypeTransformation_isTypeName_5594(ttlAst) {
		return this_.priv_TypeTransformation_evalTypeName_5618(ttlAst)
	}
	if this_.priv_TypeTransformation_isTypeVar_5593(ttlAst) {
		return this_.priv_TypeTransformation_evalTypeVar_5620(ttlAst, nameResolver)
	}
	var name string = this_.priv_TypeTransformation_getCallName_5607(ttlAst)
	_ = name
	var keyword *pkg_parsing.TypeTransformationParser_Keywords = this_.priv_TypeTransformation_nameToKeyword_5596(name)
	_ = keyword
	switch (keyword.FieldGet_kind()).Ordinal() {
	case int32(0) /* OperationKind.TYPE_CONSTRUCTOR */ :
		return this_.priv_TypeTransformation_evalTypeExpression_5617(ttlAst, nameResolver)
	case int32(1) /* OperationKind.OPERATION */ :
		return this_.priv_TypeTransformation_evalOperationExpression_5616(ttlAst, nameResolver)
	default:
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Could not evaluate the type transformation expression")))
	}
	panic("unreachable switch expression")
}

// Java method: TypeTransformation.evalOperationExpression(com.google.javascript.rhino.Node ttlAst, com.google.javascript.jscomp.TypeTransformation.NameResolver nameResolver)
func (this_ *TypeTransformation) priv_TypeTransformation_evalOperationExpression_5616(ttlAst pkg_rhino.Node, nameResolver *TypeTransformation_NameResolver) pkg_rhino.JSType {
	_ = ttlAst
	_ = nameResolver
	var name string = this_.priv_TypeTransformation_getCallName_5607(ttlAst)
	_ = name
	var keyword *pkg_parsing.TypeTransformationParser_Keywords = this_.priv_TypeTransformation_nameToKeyword_5596(name)
	_ = keyword
	switch j2gort.CheckNotNullPtr(keyword).Ordinal() {
	case int32(1) /* Keywords.COND */ :
		return this_.priv_TypeTransformation_evalConditional_5630(ttlAst, nameResolver)
	case int32(9) /* Keywords.MAPUNION */ :
		return this_.priv_TypeTransformation_evalMapunion_5631(ttlAst, nameResolver)
	case int32(10) /* Keywords.MAPRECORD */ :
		return this_.priv_TypeTransformation_evalMaprecord_5639(ttlAst, nameResolver)
	case int32(21) /* Keywords.TYPEOFVAR */ :
		return this_.priv_TypeTransformation_evalTypeOfVar_5640(ttlAst)
	case int32(8) /* Keywords.INSTANCEOF */ :
		return this_.priv_TypeTransformation_evalInstanceOf_5641(ttlAst, nameResolver)
	case int32(12) /* Keywords.PRINTTYPE */ :
		return this_.priv_TypeTransformation_evalPrintType_5643(ttlAst, nameResolver)
	case int32(13) /* Keywords.PROPTYPE */ :
		return this_.priv_TypeTransformation_evalPropType_5644(ttlAst, nameResolver)
	default:
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Invalid type transformation operation")))
	}
	panic("unreachable switch expression")
}

// Java method: TypeTransformation.evalTypeExpression(com.google.javascript.rhino.Node ttlAst, com.google.javascript.jscomp.TypeTransformation.NameResolver nameResolver)
func (this_ *TypeTransformation) priv_TypeTransformation_evalTypeExpression_5617(ttlAst pkg_rhino.Node, nameResolver *TypeTransformation_NameResolver) pkg_rhino.JSType {
	_ = ttlAst
	_ = nameResolver
	var name string = this_.priv_TypeTransformation_getCallName_5607(ttlAst)
	_ = name
	var keyword *pkg_parsing.TypeTransformationParser_Keywords = this_.priv_TypeTransformation_nameToKeyword_5596(name)
	_ = keyword
	switch j2gort.CheckNotNullPtr(keyword).Ordinal() {
	case int32(19) /* Keywords.TYPE */ :
		return this_.priv_TypeTransformation_evalTemplatizedType_5619(ttlAst, nameResolver)
	case int32(22) /* Keywords.UNION */ :
		return this_.priv_TypeTransformation_evalUnionType_5621(ttlAst, nameResolver)
	case int32(11) /* Keywords.NONE */ :
		return this_.priv_TypeTransformation_getNoType_5599()
	case int32(0) /* Keywords.ALL */ :
		return this_.priv_TypeTransformation_getAllType_5600()
	case int32(23) /* Keywords.UNKNOWN */ :
		return this_.priv_TypeTransformation_getUnknownType_5598()
	case int32(14) /* Keywords.RAWTYPEOF */ :
		return this_.priv_TypeTransformation_evalRawTypeOf_5632(ttlAst, nameResolver)
	case int32(18) /* Keywords.TEMPLATETYPEOF */ :
		return this_.priv_TypeTransformation_evalTemplateTypeOf_5633(ttlAst, nameResolver)
	case int32(17) /* Keywords.RECORD */ :
		return this_.priv_TypeTransformation_evalRecordType_5636(ttlAst, nameResolver)
	case int32(20) /* Keywords.TYPEEXPR */ :
		return this_.priv_TypeTransformation_evalNativeTypeExpr_5642(ttlAst)
	default:
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Invalid type expression")))
	}
	panic("unreachable switch expression")
}

// Java method: TypeTransformation.evalTypeName(com.google.javascript.rhino.Node ttlAst)
func (this_ *TypeTransformation) priv_TypeTransformation_evalTypeName_5618(ttlAst pkg_rhino.Node) pkg_rhino.JSType {
	_ = ttlAst
	var typeName string = ttlAst.GetString()
	_ = typeName
	var resultingType pkg_rhino.JSType = this_.priv_TypeTransformation_getType_5597(typeName)
	_ = resultingType
	if resultingType == nil || resultingType.IsJavaNil_() {
		this_.priv_TypeTransformation_reportWarning_5604(ttlAst, StaticGet_TypeTransformation_UNKNOWN_TYPENAME(), []string{typeName})
		return this_.priv_TypeTransformation_getUnknownType_5598()
	}
	return resultingType
}

// Java method: TypeTransformation.evalTemplatizedType(com.google.javascript.rhino.Node ttlAst, com.google.javascript.jscomp.TypeTransformation.NameResolver nameResolver)
func (this_ *TypeTransformation) priv_TypeTransformation_evalTemplatizedType_5619(ttlAst pkg_rhino.Node, nameResolver *TypeTransformation_NameResolver) pkg_rhino.JSType {
	_ = ttlAst
	_ = nameResolver
	var params []pkg_rhino.Node = j2gort.ImmListToNodeSliceNullable[pkg_rhino.Node](this_.priv_TypeTransformation_getCallParams_5610(ttlAst))
	_ = params
	var firstParam pkg_rhino.JSType = this_.priv_TypeTransformation_evalInternal_5615(params[int32(0)], nameResolver)
	_ = firstParam
	base_firstParam_ := firstParam.GetJSTypeBase()
	if !(base_firstParam_.IsRawTypeOfTemplatizedType()) {
		this_.priv_TypeTransformation_reportWarning_5604(ttlAst, StaticGet_TypeTransformation_BASETYPE_INVALID(), []string{base_firstParam_.String()})
		return this_.priv_TypeTransformation_getUnknownType_5598()
	}
	var templatizedTypes []pkg_rhino.JSType = make([]pkg_rhino.JSType, (int32(len(params)) - int32(1)))
	_ = templatizedTypes
	func() {
		arrMut_1 := templatizedTypes
		var arrBox_2 []any
		if arrMut_1 != nil {
			arrBox_2 = make([]any, len(arrMut_1))
			for i_, el_3 := range arrMut_1 {
				arrBox_2[i_] = el_3
			}
		}
		j2gort.Ext_java_util_Arrays_SetAll_2_slice_any_Ext_java_util_function_IntFunction(arrBox_2, j2gort.NewLambda_Ext_java_util_function_IntFunction(func(i int32) any {
			_ = i
			return this_.priv_TypeTransformation_evalInternal_5615(params[(i+int32(1))], nameResolver)
		}))
		for i_, el_3 := range arrBox_2 {
			arrMut_1[i_] = j2gort.AsIface[pkg_rhino.JSType](el_3)
		}
	}()
	var baseType pkg_rhino.ObjectType = base_firstParam_.ToMaybeObjectType()
	_ = baseType
	return pkg_rhino.ToIface_JSType(this_.registry.CreateTemplatizedType_2_ObjectType_JSType(baseType, templatizedTypes))
}

// Java method: TypeTransformation.evalTypeVar(com.google.javascript.rhino.Node ttlAst, com.google.javascript.jscomp.TypeTransformation.NameResolver nameResolver)
func (this_ *TypeTransformation) priv_TypeTransformation_evalTypeVar_5620(ttlAst pkg_rhino.Node, nameResolver *TypeTransformation_NameResolver) pkg_rhino.JSType {
	_ = ttlAst
	_ = nameResolver
	var typeVar string = ttlAst.GetString()
	_ = typeVar
	var resultingType pkg_rhino.JSType = j2gort.AsIface[pkg_rhino.JSType](j2gort.MapGetStr(nameResolver.typeVars, typeVar))
	_ = resultingType
	if resultingType == nil || resultingType.IsJavaNil_() {
		this_.priv_TypeTransformation_reportWarning_5604(ttlAst, StaticGet_TypeTransformation_UNKNOWN_TYPEVAR(), []string{typeVar})
		return this_.priv_TypeTransformation_getUnknownType_5598()
	}
	return resultingType
}

// Java method: TypeTransformation.evalUnionType(com.google.javascript.rhino.Node ttlAst, com.google.javascript.jscomp.TypeTransformation.NameResolver nameResolver)
func (this_ *TypeTransformation) priv_TypeTransformation_evalUnionType_5621(ttlAst pkg_rhino.Node, nameResolver *TypeTransformation_NameResolver) pkg_rhino.JSType {
	_ = ttlAst
	_ = nameResolver
	var params []pkg_rhino.Node = j2gort.ImmListToNodeSliceNullable[pkg_rhino.Node](this_.priv_TypeTransformation_getCallParams_5610(ttlAst))
	_ = params
	var paramCount int32 = int32(len(params))
	_ = paramCount
	var basicTypes []pkg_rhino.JSType = make([]pkg_rhino.JSType, paramCount)
	_ = basicTypes
	{
		var i int32 = int32(0)
		_ = i
		for ; i < paramCount; i = i + 1 {
			{
				a_4 := basicTypes
				i_5 := i
				a_4[i_5] = this_.priv_TypeTransformation_evalInternal_5615(params[i], nameResolver)
			}
		}
	}
	return this_.priv_TypeTransformation_createUnionType_5602(basicTypes)
}

// Java method: TypeTransformation.evalTypeParams(com.google.javascript.rhino.Node ttlAst, com.google.javascript.jscomp.TypeTransformation.NameResolver nameResolver)
func (this_ *TypeTransformation) priv_TypeTransformation_evalTypeParams_5622(ttlAst pkg_rhino.Node, nameResolver *TypeTransformation_NameResolver) []pkg_rhino.JSType {
	_ = ttlAst
	_ = nameResolver
	var params []pkg_rhino.Node = j2gort.ImmListToNodeSliceNullable[pkg_rhino.Node](this_.priv_TypeTransformation_getCallParams_5610(ttlAst))
	_ = params
	var paramCount int32 = int32(len(params))
	_ = paramCount
	var result []pkg_rhino.JSType = make([]pkg_rhino.JSType, paramCount)
	_ = result
	{
		var i int32 = int32(0)
		_ = i
		for ; i < paramCount; i = i + 1 {
			{
				a_6 := result
				i_7 := i
				a_6[i_7] = this_.priv_TypeTransformation_evalInternal_5615(params[i], nameResolver)
			}
		}
	}
	return result
}

// Java method: TypeTransformation.evalString(com.google.javascript.rhino.Node ttlAst, com.google.javascript.jscomp.TypeTransformation.NameResolver nameResolver)
func (this_ *TypeTransformation) priv_TypeTransformation_evalString_5623(ttlAst pkg_rhino.Node, nameResolver *TypeTransformation_NameResolver) string {
	_ = ttlAst
	_ = nameResolver
	if ttlAst.IsName() {
		if !(j2gort.MapContainsKeyStr(nameResolver.nameVars, ttlAst.GetString())) {
			this_.priv_TypeTransformation_reportWarning_5604(ttlAst, StaticGet_TypeTransformation_UNKNOWN_STRVAR(), []string{ttlAst.GetString()})
			return ""
		}
		return j2gort.CastString(j2gort.MapGetStr(nameResolver.nameVars, ttlAst.GetString()))
	}
	return ttlAst.GetString()
}

// Java method: TypeTransformation.evalStringParams(com.google.javascript.rhino.Node ttlAst, com.google.javascript.jscomp.TypeTransformation.NameResolver nameResolver)
func (this_ *TypeTransformation) priv_TypeTransformation_evalStringParams_5624(ttlAst pkg_rhino.Node, nameResolver *TypeTransformation_NameResolver) []string {
	_ = ttlAst
	_ = nameResolver
	var params []pkg_rhino.Node = j2gort.ImmListToNodeSliceNullable[pkg_rhino.Node](this_.priv_TypeTransformation_getCallParams_5610(ttlAst))
	_ = params
	var paramCount int32 = int32(len(params))
	_ = paramCount
	var result []string = j2gort.NewStringArray(paramCount)
	_ = result
	{
		var i int32 = int32(0)
		_ = i
		for ; i < paramCount; i = i + 1 {
			{
				a_8 := result
				i_9 := i
				a_8[i_9] = this_.priv_TypeTransformation_evalString_5623(params[i], nameResolver)
			}
		}
	}
	return result
}

// Java method: TypeTransformation.evalTypePredicate(com.google.javascript.rhino.Node ttlAst, com.google.javascript.jscomp.TypeTransformation.NameResolver nameResolver)
func (this_ *TypeTransformation) priv_TypeTransformation_evalTypePredicate_5625(ttlAst pkg_rhino.Node, nameResolver *TypeTransformation_NameResolver) bool {
	_ = ttlAst
	_ = nameResolver
	var params []pkg_rhino.JSType = this_.priv_TypeTransformation_evalTypeParams_5622(ttlAst, nameResolver)
	_ = params
	var name string = this_.priv_TypeTransformation_getCallName_5607(ttlAst)
	_ = name
	var keyword *pkg_parsing.TypeTransformationParser_Keywords = this_.priv_TypeTransformation_nameToKeyword_5596(name)
	_ = keyword
	var type_ pkg_rhino.JSType = params[int32(0)]
	_ = type_
	switch j2gort.CheckNotNullPtr(keyword).Ordinal() {
	case int32(2) /* Keywords.EQ */ :
		return type_.GetJSTypeBase().EqualsJSType(params[int32(1)])
	case int32(15) /* Keywords.SUB */ :
		return type_.GetJSTypeBase().IsSubtypeOf_1(params[int32(1)])
	case int32(3) /* Keywords.ISCTOR */ :
		return type_.GetJSTypeBase().IsConstructor()
	case int32(6) /* Keywords.ISTEMPLATIZED */ :
		return type_.GetJSTypeBase().IsTemplatizedType()
	case int32(5) /* Keywords.ISRECORD */ :
		return type_.GetJSTypeBase().IsRecordType()
	case int32(7) /* Keywords.ISUNKNOWN */ :
		return type_.GetJSTypeBase().IsUnknownType()
	default:
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Invalid type predicate in the type transformation")))
	}
	panic("unreachable switch expression")
}

// Java method: TypeTransformation.evalStringPredicate(com.google.javascript.rhino.Node ttlAst, com.google.javascript.jscomp.TypeTransformation.NameResolver nameResolver)
func (this_ *TypeTransformation) priv_TypeTransformation_evalStringPredicate_5626(ttlAst pkg_rhino.Node, nameResolver *TypeTransformation_NameResolver) bool {
	_ = ttlAst
	_ = nameResolver
	var params []string = this_.priv_TypeTransformation_evalStringParams_5624(ttlAst, nameResolver)
	_ = params
	{
		var i int32 = int32(0)
		_ = i
		for ; i < int32(len(params)); i = i + 1 {
			if j2gort.CheckNotNullString(params[i]) == "" {
				return false
			}
		}
	}
	var name string = this_.priv_TypeTransformation_getCallName_5607(ttlAst)
	_ = name
	var keyword *pkg_parsing.TypeTransformationParser_Keywords = this_.priv_TypeTransformation_nameToKeyword_5596(name)
	_ = keyword
	switch j2gort.CheckNotNullPtr(keyword).Ordinal() {
	case int32(16) /* Keywords.STREQ */ :
		{

			sr_10 := params[int32(0)]
			se_11 := params[int32(1)]
			return j2gort.CheckNotNullString(sr_10) == se_11
		}

	default:
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Invalid string predicate in the type transformation")))
	}
	panic("unreachable switch expression")
}

// Java method: TypeTransformation.evalTypevarPredicate(com.google.javascript.rhino.Node ttlAst, com.google.javascript.jscomp.TypeTransformation.NameResolver nameResolver)
func (this_ *TypeTransformation) priv_TypeTransformation_evalTypevarPredicate_5627(ttlAst pkg_rhino.Node, nameResolver *TypeTransformation_NameResolver) bool {
	_ = ttlAst
	_ = nameResolver
	var name string = this_.priv_TypeTransformation_getCallName_5607(ttlAst)
	_ = name
	var keyword *pkg_parsing.TypeTransformationParser_Keywords = this_.priv_TypeTransformation_nameToKeyword_5596(name)
	_ = keyword
	switch j2gort.CheckNotNullPtr(keyword).Ordinal() {
	case int32(4) /* Keywords.ISDEFINED */ :
		return j2gort.MapContainsKeyStr(nameResolver.typeVars, this_.priv_TypeTransformation_getCallArgument_5608(ttlAst, int32(0)).GetString())
	default:
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Invalid typevar predicate in the type transformation")))
	}
	panic("unreachable switch expression")
}

// Java method: TypeTransformation.evalBooleanOperation(com.google.javascript.rhino.Node ttlAst, com.google.javascript.jscomp.TypeTransformation.NameResolver nameResolver)
func (this_ *TypeTransformation) priv_TypeTransformation_evalBooleanOperation_5628(ttlAst pkg_rhino.Node, nameResolver *TypeTransformation_NameResolver) bool {
	_ = ttlAst
	_ = nameResolver
	var param0 bool = this_.priv_TypeTransformation_evalBoolean_5629(ttlAst.GetFirstChild(), nameResolver)
	_ = param0
	if ttlAst.IsNot() {
		return !(param0)
	}
	if ttlAst.IsAnd() {
		return (param0 && this_.priv_TypeTransformation_evalBoolean_5629(ttlAst.GetLastChild(), nameResolver))
	}
	if ttlAst.IsOr() {
		return (param0 || this_.priv_TypeTransformation_evalBoolean_5629(ttlAst.GetLastChild(), nameResolver))
	}
	panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Invalid boolean predicate in the type transformation")))
}

// Java method: TypeTransformation.evalBoolean(com.google.javascript.rhino.Node ttlAst, com.google.javascript.jscomp.TypeTransformation.NameResolver nameResolver)
func (this_ *TypeTransformation) priv_TypeTransformation_evalBoolean_5629(ttlAst pkg_rhino.Node, nameResolver *TypeTransformation_NameResolver) bool {
	_ = ttlAst
	_ = nameResolver
	if this_.priv_TypeTransformation_isBooleanOperation_5595(ttlAst) {
		return this_.priv_TypeTransformation_evalBooleanOperation_5628(ttlAst, nameResolver)
	}
	var name string = this_.priv_TypeTransformation_getCallName_5607(ttlAst)
	_ = name
	var keyword *pkg_parsing.TypeTransformationParser_Keywords = this_.priv_TypeTransformation_nameToKeyword_5596(name)
	_ = keyword
	switch (keyword.FieldGet_kind()).Ordinal() {
	case int32(2) /* OperationKind.STRING_PREDICATE */ :
		return this_.priv_TypeTransformation_evalStringPredicate_5626(ttlAst, nameResolver)
	case int32(3) /* OperationKind.TYPE_PREDICATE */ :
		return this_.priv_TypeTransformation_evalTypePredicate_5625(ttlAst, nameResolver)
	case int32(4) /* OperationKind.TYPEVAR_PREDICATE */ :
		return this_.priv_TypeTransformation_evalTypevarPredicate_5627(ttlAst, nameResolver)
	default:
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Invalid boolean predicate in the type transformation")))
	}
	panic("unreachable switch expression")
}

// Java method: TypeTransformation.evalConditional(com.google.javascript.rhino.Node ttlAst, com.google.javascript.jscomp.TypeTransformation.NameResolver nameResolver)
func (this_ *TypeTransformation) priv_TypeTransformation_evalConditional_5630(ttlAst pkg_rhino.Node, nameResolver *TypeTransformation_NameResolver) pkg_rhino.JSType {
	_ = ttlAst
	_ = nameResolver
	var params []pkg_rhino.Node = j2gort.ImmListToNodeSliceNullable[pkg_rhino.Node](this_.priv_TypeTransformation_getCallParams_5610(ttlAst))
	_ = params
	if this_.priv_TypeTransformation_evalBoolean_5629(params[int32(0)], nameResolver) {
		return this_.priv_TypeTransformation_evalInternal_5615(params[int32(1)], nameResolver)
	} else {
		return this_.priv_TypeTransformation_evalInternal_5615(params[int32(2)], nameResolver)
	}
}

// Java method: TypeTransformation.evalMapunion(com.google.javascript.rhino.Node ttlAst, com.google.javascript.jscomp.TypeTransformation.NameResolver nameResolver)
func (this_ *TypeTransformation) priv_TypeTransformation_evalMapunion_5631(ttlAst pkg_rhino.Node, nameResolver *TypeTransformation_NameResolver) pkg_rhino.JSType {
	_ = ttlAst
	_ = nameResolver
	var params []pkg_rhino.Node = j2gort.ImmListToNodeSliceNullable[pkg_rhino.Node](this_.priv_TypeTransformation_getCallParams_5610(ttlAst))
	_ = params
	var unionParam pkg_rhino.Node = params[int32(0)]
	_ = unionParam
	var mapFunction pkg_rhino.Node = params[int32(1)]
	_ = mapFunction
	var paramName string = this_.priv_TypeTransformation_getFunctionParameter_5606(mapFunction, int32(0))
	_ = paramName
	if j2gort.MapContainsKeyStr(nameResolver.typeVars, paramName) {
		this_.priv_TypeTransformation_reportWarning_5604(ttlAst, StaticGet_TypeTransformation_DUPLICATE_VARIABLE(), []string{paramName})
		return this_.priv_TypeTransformation_getUnknownType_5598()
	}
	var mapFunctionBody pkg_rhino.Node = NodeUtil_GetFunctionBody(mapFunction)
	_ = mapFunctionBody
	var unionType pkg_rhino.JSType = this_.priv_TypeTransformation_evalInternal_5615(unionParam, nameResolver)
	_ = unionType
	base_unionType_ := unionType.GetJSTypeBase()
	if !(base_unionType_.IsUnionType()) {
		var val_12 TypeTransformation_NameResolver
		InitTypeTransformation_NameResolver(&val_12, this_.priv_TypeTransformation_addNewEntry_5605(nameResolver.typeVars, paramName, unionType), nameResolver.nameVars)
		var newNameResolver *TypeTransformation_NameResolver = &val_12
		_ = newNameResolver
		return this_.priv_TypeTransformation_evalInternal_5615(mapFunctionBody, newNameResolver)
	}
	var unionElms []pkg_rhino.JSType = j2gort.ImmListCopyToSlice[pkg_rhino.JSType](base_unionType_.GetUnionMembers())
	_ = unionElms
	var unionSize int32 = int32(len(unionElms))
	_ = unionSize
	var newUnionElms []pkg_rhino.JSType = make([]pkg_rhino.JSType, unionSize)
	_ = newUnionElms
	var i int32 = int32(0)
	_ = i
	{

		for _, elm := range unionElms {

			_ = elm
			var newNameResolver_2 *TypeTransformation_NameResolver = NewTypeTransformation_NameResolver(this_.priv_TypeTransformation_addNewEntry_5605(nameResolver.typeVars, paramName, elm), nameResolver.nameVars)
			_ = newNameResolver_2
			{
				a_15 := newUnionElms
				i_16 := i
				a_15[i_16] = this_.priv_TypeTransformation_evalInternal_5615(mapFunctionBody, newNameResolver_2)
			}
			i = i + 1
		}
	}
	return this_.priv_TypeTransformation_createUnionType_5602(newUnionElms)
}

// Java method: TypeTransformation.evalRawTypeOf(com.google.javascript.rhino.Node ttlAst, com.google.javascript.jscomp.TypeTransformation.NameResolver nameResolver)
func (this_ *TypeTransformation) priv_TypeTransformation_evalRawTypeOf_5632(ttlAst pkg_rhino.Node, nameResolver *TypeTransformation_NameResolver) pkg_rhino.JSType {
	_ = ttlAst
	_ = nameResolver
	var params []pkg_rhino.Node = j2gort.ImmListToNodeSliceNullable[pkg_rhino.Node](this_.priv_TypeTransformation_getCallParams_5610(ttlAst))
	_ = params
	var type_ pkg_rhino.JSType = this_.priv_TypeTransformation_evalInternal_5615(params[int32(0)], nameResolver)
	_ = type_
	base_type__ := type_.GetJSTypeBase()
	if !(base_type__.IsTemplatizedType()) {
		this_.priv_TypeTransformation_reportWarning_5604(ttlAst, StaticGet_TypeTransformation_TEMPTYPE_INVALID(), []string{"rawTypeOf", base_type__.String()})
		return this_.priv_TypeTransformation_getUnknownType_5598()
	}
	return pkg_rhino.Upcast_JSType(base_type__.ToMaybeObjectType().GetObjectTypeBase(

	// Java method: TypeTransformation.evalTemplateTypeOf(com.google.javascript.rhino.Node ttlAst, com.google.javascript.jscomp.TypeTransformation.NameResolver nameResolver)
	).GetRawType())
}

func (this_ *TypeTransformation) priv_TypeTransformation_evalTemplateTypeOf_5633(ttlAst pkg_rhino.Node, nameResolver *TypeTransformation_NameResolver) pkg_rhino.JSType {
	_ = ttlAst
	_ = nameResolver
	var params []pkg_rhino.Node = j2gort.ImmListToNodeSliceNullable[pkg_rhino.Node](this_.priv_TypeTransformation_getCallParams_5610(ttlAst))
	_ = params
	var type_ pkg_rhino.JSType = this_.priv_TypeTransformation_evalInternal_5615(params[int32(0)], nameResolver)
	_ = type_
	base_type__ := type_.GetJSTypeBase()
	if !(base_type__.IsTemplatizedType()) {
		this_.priv_TypeTransformation_reportWarning_5604(ttlAst, StaticGet_TypeTransformation_TEMPTYPE_INVALID(), []string{"templateTypeOf", base_type__.String()})
		return this_.priv_TypeTransformation_getUnknownType_5598()
	}
	var index int32 = javalang.D2I(params[int32(1)].GetDouble())
	_ = index
	var templateTypes []pkg_rhino.JSType = j2gort.ImmListToSliceNullable[pkg_rhino.JSType](base_type__.ToMaybeObjectType().GetTemplateTypes())
	_ = templateTypes
	if index >= int32(len(templateTypes)) {
		this_.priv_TypeTransformation_reportWarning_5604(ttlAst, StaticGet_TypeTransformation_INDEX_OUTOFBOUNDS(), []string{strconv.FormatInt(int64(int32(len(templateTypes))), 10), strconv.FormatInt(int64(index), 10)})
		return this_.priv_TypeTransformation_getUnknownType_5598()
	}
	return templateTypes[index]
}

// Java method: TypeTransformation.evalRecord(com.google.javascript.rhino.Node record, com.google.javascript.jscomp.TypeTransformation.NameResolver nameResolver)
func (this_ *TypeTransformation) priv_TypeTransformation_evalRecord_5634(record pkg_rhino.Node, nameResolver *TypeTransformation_NameResolver) pkg_rhino.JSType {
	_ = record
	_ = nameResolver
	var props j2gort.Ext_java_util_Map = j2gort.NewExt_java_util_LinkedHashMap()
	_ = props
	{
		var propNode pkg_rhino.Node = record.GetFirstChild()
		_ = propNode
		for ; propNode != 0; propNode = propNode.GetNext() {
			if propNode.IsComputedProp() {
				var compPropName string = this_.priv_TypeTransformation_getComputedPropName_5612(propNode)
				_ = compPropName
				if !(j2gort.MapContainsKeyStr(nameResolver.nameVars, compPropName)) {
					this_.priv_TypeTransformation_reportWarning_5604(record, StaticGet_TypeTransformation_UNKNOWN_NAMEVAR(), []string{compPropName})
					return this_.priv_TypeTransformation_getUnknownType_5598()
				}
				var propValue pkg_rhino.Node = this_.priv_TypeTransformation_getComputedPropValue_5611(propNode)
				_ = propValue
				var resolvedName string = j2gort.CastString(j2gort.MapGetStr(nameResolver.nameVars, compPropName))
				_ = resolvedName
				var resultingType pkg_rhino.JSType = this_.priv_TypeTransformation_evalInternal_5615(propValue, nameResolver)
				_ = resultingType
				j2gort.MapPutStr(props, resolvedName, resultingType)
			} else {
				var propName string = propNode.GetString()
				_ = propName
				var resultingType_2 pkg_rhino.JSType = this_.priv_TypeTransformation_evalInternal_5615(propNode.GetFirstChild(), nameResolver)
				_ = resultingType_2
				j2gort.MapPutStr(props, propName, resultingType_2)
			}
		}
	}
	return this_.registry.CreateRecordType(props)
}

// Java method: TypeTransformation.evalRecordParam(com.google.javascript.rhino.Node ttlAst, com.google.javascript.jscomp.TypeTransformation.NameResolver nameResolver)
func (this_ *TypeTransformation) priv_TypeTransformation_evalRecordParam_5635(ttlAst pkg_rhino.Node, nameResolver *TypeTransformation_NameResolver) pkg_rhino.JSType {
	_ = ttlAst
	_ = nameResolver
	if ttlAst.IsObjectLit() {
		return this_.priv_TypeTransformation_evalRecord_5634(ttlAst, nameResolver)
	}
	return this_.priv_TypeTransformation_evalInternal_5615(ttlAst, nameResolver)
}

// Java method: TypeTransformation.evalRecordType(com.google.javascript.rhino.Node ttlAst, com.google.javascript.jscomp.TypeTransformation.NameResolver nameResolver)
func (this_ *TypeTransformation) priv_TypeTransformation_evalRecordType_5636(ttlAst pkg_rhino.Node, nameResolver *TypeTransformation_NameResolver) pkg_rhino.JSType {
	_ = ttlAst
	_ = nameResolver
	var paramCount int32 = this_.priv_TypeTransformation_getCallParamCount_5609(ttlAst)
	_ = paramCount
	var recTypesBuilder []pkg_rhino.ObjectType = []pkg_rhino.ObjectType{}
	_ = recTypesBuilder
	{
		var i int32 = int32(0)
		_ = i
		for ; i < paramCount; i = i + 1 {
			var type_ pkg_rhino.JSType = this_.priv_TypeTransformation_evalRecordParam_5635(this_.priv_TypeTransformation_getCallArgument_5608(ttlAst, i), nameResolver)
			_ = type_
			var objType pkg_rhino.ObjectType = type_.GetJSTypeBase().ToMaybeObjectType()
			_ = objType
			if (objType == nil || objType.IsJavaNil_()) || objType.GetObjectTypeBase().IsUnknownType() {
				this_.priv_TypeTransformation_reportWarning_5604(ttlAst, StaticGet_TypeTransformation_RECPARAM_INVALID(), []string{type_.GetJSTypeBase().String()})
				return this_.priv_TypeTransformation_getUnknownType_5598()
			}
			var recType pkg_rhino.JSType = this_.registry.BuildRecordTypeFromObject(objType)
			_ = recType
			if !(recType.GetJSTypeBase().EqualsJSType(this_.priv_TypeTransformation_getObjectType_5601())) {
				recTypesBuilder = append(recTypesBuilder, j2gort.CheckNotNullElem(recType.GetJSTypeBase().ToMaybeObjectType()))
			}
		}
	}
	return this_.priv_TypeTransformation_joinRecordTypes_5638(j2gort.ImmListFromTypedSlice(recTypesBuilder))
}

// Java method: TypeTransformation.putNewPropInPropertyMap(java.util.Map<java.lang.String,com.google.javascript.rhino.jstype.JSType> props, java.lang.String newPropName, com.google.javascript.rhino.jstype.JSType newPropValue)
func (this_ *TypeTransformation) priv_TypeTransformation_putNewPropInPropertyMap_5637(props *j2gort.OrderedMap[string, pkg_rhino.JSType], newPropName string, newPropValue pkg_rhino.JSType) {
	_ = props
	_ = newPropName
	_ = newPropValue
	if (!(props.ContainsKey(newPropName)) || !(newPropValue.GetJSTypeBase().IsRecordType())) || !(props.Get(newPropName).GetJSTypeBase().IsRecordType()) {
		props.Put(newPropName, newPropValue)
		return
	}
	props.Put(newPropName, this_.priv_TypeTransformation_joinRecordTypes_5638(j2gort.Ext_com_google_common_collect_ImmutableList_Of_2(j2gort.AsIface[pkg_rhino.ObjectType](props.Get(newPropName)), j2gort.AsIface[pkg_rhino.ObjectType](newPropValue))))
}

// Java method: TypeTransformation.joinRecordTypes(com.google.common.collect.ImmutableList<com.google.javascript.rhino.jstype.ObjectType> recTypes)
func (this_ *TypeTransformation) priv_TypeTransformation_joinRecordTypes_5638(recTypes j2gort.Ext_com_google_common_collect_ImmutableList) pkg_rhino.JSType {
	_ = recTypes
	var props j2gort.OrderedMap[string, pkg_rhino.JSType] = j2gort.NewOrderedMap[string, pkg_rhino.JSType]()
	_ = props
	{
		cur_17 := j2gort.NewIterCursor(j2gort.CheckNotNull(recTypes))
		for cur_17.HasNext() {
			rawElem_18 := cur_17.Next()
			recType := j2gort.AsIface[pkg_rhino.ObjectType](rawElem_18)
			_ = recType
			{
				cur_19 := j2gort.NewIterCursor(j2gort.CheckNotNull(recType.GetOwnPropertyNames()))
				for cur_19.HasNext() {
					newPropName := cur_19.NextStr()
					_ = newPropName
					var newPropValue pkg_rhino.JSType = recType.GetObjectTypeBase().GetPropertyType_1_string(newPropName)
					_ = newPropValue
					this_.priv_TypeTransformation_putNewPropInPropertyMap_5637(&props, newPropName, newPropValue)
				}
			}
		}
	}
	return this_.priv_TypeTransformation_createRecordType_5603(props.ToImmutableMap())
}

// Java method: TypeTransformation.evalMaprecord(com.google.javascript.rhino.Node ttlAst, com.google.javascript.jscomp.TypeTransformation.NameResolver nameResolver)
func (this_ *TypeTransformation) priv_TypeTransformation_evalMaprecord_5639(ttlAst pkg_rhino.Node, nameResolver *TypeTransformation_NameResolver) pkg_rhino.JSType {
	_ = ttlAst
	_ = nameResolver
	var recordNode pkg_rhino.Node = ttlAst.GetSecondChild()
	_ = recordNode
	var mapFunction pkg_rhino.Node = ttlAst.GetChildAtIndex(int32(2))
	_ = mapFunction
	var type_ pkg_rhino.JSType = this_.priv_TypeTransformation_evalInternal_5615(recordNode, nameResolver)
	_ = type_

	base_type__ := type_.GetJSTypeBase()
	if base_type__.EqualsJSType(this_.priv_TypeTransformation_getObjectType_5601()) {
		return this_.priv_TypeTransformation_getObjectType_5601()
	}

	if !(base_type__.IsRecordType()) {
		this_.priv_TypeTransformation_reportWarning_5604(recordNode, StaticGet_TypeTransformation_RECTYPE_INVALID(), []string{base_type__.String()})
		return this_.priv_TypeTransformation_getUnknownType_5598()
	}
	var objtype pkg_rhino.ObjectType = base_type__.ToMaybeObjectType()
	_ = objtype
	var paramKey string = this_.priv_TypeTransformation_getFunctionParameter_5606(mapFunction, int32(0))
	_ = paramKey
	var paramValue string = this_.priv_TypeTransformation_getFunctionParameter_5606(mapFunction, int32(1))
	_ = paramValue
	if j2gort.MapContainsKeyStr(nameResolver.nameVars, paramKey) {
		this_.priv_TypeTransformation_reportWarning_5604(ttlAst, StaticGet_TypeTransformation_DUPLICATE_VARIABLE(), []string{paramKey})
		return this_.priv_TypeTransformation_getUnknownType_5598()
	}
	if j2gort.MapContainsKeyStr(nameResolver.typeVars, paramValue) {
		this_.priv_TypeTransformation_reportWarning_5604(ttlAst, StaticGet_TypeTransformation_DUPLICATE_VARIABLE(), []string{paramValue})
		return this_.priv_TypeTransformation_getUnknownType_5598()
	}
	var mapFnBody pkg_rhino.Node = NodeUtil_GetFunctionBody(mapFunction)
	_ = mapFnBody
	var newProps j2gort.OrderedMap[string, pkg_rhino.JSType] = j2gort.NewOrderedMap[string, pkg_rhino.JSType]()
	_ = newProps
	{
		cur_20 := j2gort.NewIterCursor(j2gort.CheckNotNull(objtype.GetOwnPropertyNames()))
		for cur_20.HasNext() {
			propName := cur_20.NextStr()
			_ = propName
			var propValue pkg_rhino.JSType = objtype.GetObjectTypeBase().GetPropertyType_1_string(propName)
			_ = propValue
			var newNameResolver *TypeTransformation_NameResolver = NewTypeTransformation_NameResolver(this_.priv_TypeTransformation_addNewEntry_5605(nameResolver.typeVars, paramValue, propValue), this_.priv_TypeTransformation_addNewEntry_5605(nameResolver.nameVars, paramKey, j2gort.BoxString(propName)))
			_ = newNameResolver
			var body pkg_rhino.JSType = this_.priv_TypeTransformation_evalInternal_5615(mapFnBody, newNameResolver)
			_ = body
			if body.GetJSTypeBase().IsUnknownType() {
				return this_.priv_TypeTransformation_getUnknownType_5598()
			}
			if body.GetJSTypeBase().IsEmptyType() || body.GetJSTypeBase().EqualsJSType(this_.priv_TypeTransformation_getObjectType_5601()) {
				continue
			}
			if !(body.GetJSTypeBase().IsRecordType()) {
				this_.priv_TypeTransformation_reportWarning_5604(ttlAst, StaticGet_TypeTransformation_MAPRECORD_BODY_INVALID(), []string{body.GetJSTypeBase().String()})
				return this_.priv_TypeTransformation_getUnknownType_5598()
			}
			var bodyAsObj pkg_rhino.ObjectType = body.GetJSTypeBase().ToMaybeObjectType()
			_ = bodyAsObj
			{
				cur_21 := j2gort.NewIterCursor(j2gort.CheckNotNull(bodyAsObj.GetOwnPropertyNames()))
				for cur_21.HasNext() {
					newPropName := cur_21.NextStr()
					_ = newPropName
					var newPropValue pkg_rhino.JSType = bodyAsObj.GetObjectTypeBase().GetPropertyType_1_string(newPropName)
					_ = newPropValue
					this_.priv_TypeTransformation_putNewPropInPropertyMap_5637(&newProps, newPropName, newPropValue)
				}
			}
		}
	}
	return this_.priv_TypeTransformation_createRecordType_5603(newProps.ToImmutableMap())
}

// Java method: TypeTransformation.evalTypeOfVar(com.google.javascript.rhino.Node ttlAst)
func (this_ *TypeTransformation) priv_TypeTransformation_evalTypeOfVar_5640(ttlAst pkg_rhino.Node) pkg_rhino.JSType {
	_ = ttlAst
	var name string = this_.priv_TypeTransformation_getCallArgument_5608(ttlAst, int32(0)).GetString()
	_ = name
	var slot pkg_rhino.StaticTypedSlot = j2gort.AsIface[pkg_rhino.StaticTypedSlot](this_.typeEnv.GetSlot(name))
	_ = slot
	var type_ pkg_rhino.JSType
	_ = type_
	if slot != nil && !slot.IsJavaNil_() {
		type_ = slot.GetType()
	} else {
		type_ = nil
	}
	if type_ == nil || type_.IsJavaNil_() {
		this_.priv_TypeTransformation_reportWarning_5604(ttlAst, StaticGet_TypeTransformation_VAR_UNDEFINED(), []string{name})
		return this_.priv_TypeTransformation_getUnknownType_5598()
	}
	return type_
}

// Java method: TypeTransformation.evalInstanceOf(com.google.javascript.rhino.Node ttlAst, com.google.javascript.jscomp.TypeTransformation.NameResolver nameResolver)
func (this_ *TypeTransformation) priv_TypeTransformation_evalInstanceOf_5641(ttlAst pkg_rhino.Node, nameResolver *TypeTransformation_NameResolver) pkg_rhino.JSType {
	_ = ttlAst
	_ = nameResolver
	var type_ pkg_rhino.JSType = this_.priv_TypeTransformation_evalInternal_5615(this_.priv_TypeTransformation_getCallArgument_5608(ttlAst, int32(0)), nameResolver)
	_ = type_
	base_type__ := type_.GetJSTypeBase()
	if base_type__.IsUnknownType() || !(base_type__.IsConstructor()) {
		this_.priv_TypeTransformation_reportWarning_5604(ttlAst, StaticGet_TypeTransformation_INVALID_CTOR(), []string{type_.GetDisplayName()})
		return this_.priv_TypeTransformation_getUnknownType_5598()
	}
	return pkg_rhino.Upcast_JSType(type_.ToMaybeFunctionType().GetInstanceType())
}

// Java method: TypeTransformation.evalNativeTypeExpr(com.google.javascript.rhino.Node ttlAst)
func (this_ *TypeTransformation) priv_TypeTransformation_evalNativeTypeExpr_5642(ttlAst pkg_rhino.Node) pkg_rhino.JSType {
	_ = ttlAst
	var expr pkg_rhino.JSTypeExpression = pkg_rhino.NewJSTypeExpression(this_.priv_TypeTransformation_getCallArgument_5608(ttlAst, int32(0)).CloneTree_0(), "<TypeTransformation.java>")
	_ = expr
	return this_.registry.EvaluateTypeExpression(expr, this_.typeEnv)
}

// Java method: TypeTransformation.evalPrintType(com.google.javascript.rhino.Node ttlAst, com.google.javascript.jscomp.TypeTransformation.NameResolver nameResolver)
func (this_ *TypeTransformation) priv_TypeTransformation_evalPrintType_5643(ttlAst pkg_rhino.Node, nameResolver *TypeTransformation_NameResolver) pkg_rhino.JSType {
	_ = ttlAst
	_ = nameResolver
	var type_ pkg_rhino.JSType = this_.priv_TypeTransformation_evalInternal_5615(this_.priv_TypeTransformation_getCallArgument_5608(ttlAst, int32(1)), nameResolver)
	_ = type_
	var msg string = (j2gort.StrValString(this_.priv_TypeTransformation_getCallArgument_5608(ttlAst, int32(0)).GetString()) + j2gort.StrVal(type_))
	_ = msg
	j2gort.Ext_java_lang_System_F_Out.Println_1_string(msg)
	return type_
}

// Java method: TypeTransformation.evalPropType(com.google.javascript.rhino.Node ttlAst, com.google.javascript.jscomp.TypeTransformation.NameResolver nameResolver)
func (this_ *TypeTransformation) priv_TypeTransformation_evalPropType_5644(ttlAst pkg_rhino.Node, nameResolver *TypeTransformation_NameResolver) pkg_rhino.JSType {
	_ = ttlAst
	_ = nameResolver
	var type_ pkg_rhino.JSType = this_.priv_TypeTransformation_evalInternal_5615(this_.priv_TypeTransformation_getCallArgument_5608(ttlAst, int32(1)), nameResolver)
	_ = type_
	base_type__ := type_.GetJSTypeBase()
	var objType pkg_rhino.ObjectType = base_type__.ToMaybeObjectType()
	_ = objType
	if objType == nil || objType.IsJavaNil_() {
		this_.priv_TypeTransformation_reportWarning_5604(ttlAst, StaticGet_TypeTransformation_PROPTYPE_INVALID(), []string{base_type__.String()})
		return this_.priv_TypeTransformation_getUnknownType_5598()
	}
	var propType pkg_rhino.JSType = objType.GetObjectTypeBase().GetPropertyType_1_string(this_.priv_TypeTransformation_getCallArgument_5608(ttlAst, int32(0)).GetString())
	_ = propType
	return j2gort.AsIface[pkg_rhino.JSType](j2gort.Ext_com_google_common_base_MoreObjects_FirstNonNull(propType, this_.priv_TypeTransformation_getUnknownType_5598()))
}

// Java class: com.google.javascript.jscomp.TypeTransformation.NameResolver
type TypeTransformation_NameResolver struct {
	typeVars j2gort.Ext_com_google_common_collect_ImmutableMap `gson:"typeVars"`
	nameVars j2gort.Ext_com_google_common_collect_ImmutableMap `gson:"nameVars"`
	idHash_  atomic.Int32
}

func (this_ *TypeTransformation_NameResolver) BindSelf_(s any) {
	_ = s
}
func (this_ *TypeTransformation_NameResolver) IsJavaNil_() bool { return this_ == nil }
func (*TypeTransformation_NameResolver) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.TypeTransformation$NameResolver", "NameResolver"
}
func (this_ *TypeTransformation_NameResolver) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*TypeTransformation_NameResolver) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.TypeTransformation.NameResolver", "com.google.javascript.jscomp.TypeTransformation$NameResolver", func(v_ any) bool {
		return j2gort.InstanceOf[*TypeTransformation_NameResolver](v_)
	})
}

func (this_ *TypeTransformation_NameResolver) initFields_TypeTransformation_NameResolver() {
	_ = this_
}

// Java constructor: NameResolver(com.google.common.collect.ImmutableMap<java.lang.String,com.google.javascript.rhino.jstype.JSType> typeVars, com.google.common.collect.ImmutableMap<java.lang.String,java.lang.String> nameVars)
func NewTypeTransformation_NameResolver(typeVars j2gort.Ext_com_google_common_collect_ImmutableMap, nameVars j2gort.Ext_com_google_common_collect_ImmutableMap) *TypeTransformation_NameResolver {
	this_ := &TypeTransformation_NameResolver{}
	InitTypeTransformation_NameResolver(this_, typeVars, nameVars)
	return this_
}

func InitTypeTransformation_NameResolver(this_ *TypeTransformation_NameResolver, typeVars j2gort.Ext_com_google_common_collect_ImmutableMap, nameVars j2gort.Ext_com_google_common_collect_ImmutableMap) {
	_ = typeVars
	_ = nameVars
	this_.initFields_TypeTransformation_NameResolver()
	this_.typeVars = typeVars
	this_.nameVars = nameVars
}
