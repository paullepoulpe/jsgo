/*
 * Copyright 2007 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/CodingConvention.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_coding_convention_1 j2gort.Ext_java_util_function_Function

func init() {
	lambdaSingleton_coding_convention_1 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return j2gort.AsIface[CodingConvention_AssertionFunctionSpec](mrp0_).GetCodingConvention_AssertionFunctionSpecBase().priv_AssertionFunctionSpec_getId_1387()
	})
}

// Java interface: com.google.javascript.jscomp.CodingConvention
type CodingConvention interface {
	Implements_jscomp_CodingConvention()
	IsJavaNil_() bool
	IsConstant(string) bool
	IsConstantKey(string) bool
	IsValidEnumKey(string) bool
	IsOptionalParameter(pkg_rhino.Node) bool
	IsVarArgsParameter(pkg_rhino.Node) bool
	IsFunctionCallThatAlwaysThrows(pkg_rhino.Node) bool
	IsExported_2(string, bool) bool
	IsExported_1(string) bool
	BlockRenamingForProperty(string) bool
	GetPackageName(pkg_rhino.StaticSourceFile) string
	GetClassesDefinedByCall(pkg_rhino.Node) *CodingConvention_SubclassRelationship
	IsClassFactoryCall(pkg_rhino.Node) bool
	IsSuperClassReference(string) bool
	ExtractIsModuleFile(pkg_rhino.Node, pkg_rhino.Node) bool
	ExtractClassNameIfProvide(pkg_rhino.Node, pkg_rhino.Node) string
	ExtractClassNameIfRequire(pkg_rhino.Node, pkg_rhino.Node) string
	GetExportPropertyFunction() string
	GetExportSymbolFunction() string
	IdentifyTypeDeclarationCall(pkg_rhino.Node) j2gort.Ext_java_util_List
	ApplySubclassRelationship(*pkg_rhino.NominalTypeBuilder, *pkg_rhino.NominalTypeBuilder, CodingConvention_SubclassType)
	GetAbstractMethodName() string
	GetSingletonGetterClassName(pkg_rhino.Node) string
	ApplySingletonGetter(*pkg_rhino.NominalTypeBuilder, *pkg_rhino.FunctionType)
	DescribeFunctionBind(pkg_rhino.Node, bool) *CodingConvention_Bind
	DescribeCachingCall(pkg_rhino.Node) *CodingConvention_Cache
	IsPropertyTestFunction(pkg_rhino.Node) bool
	IsPrototypeAlias(pkg_rhino.Node) bool
	IsPropertyRenameFunction(pkg_rhino.Node) bool
	GetObjectLiteralCast(pkg_rhino.Node) *CodingConvention_ObjectLiteralCast
	GetAssertionFunctions() any
}

func ToIface_CodingConvention[T interface {
	CodingConvention
	comparable
}](v T) CodingConvention {
	var z T
	if v == z {
		return nil
	}
	return v
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CodingConvention", "com.google.javascript.jscomp.CodingConvention", func(v_ any) bool {
		return j2gort.InstanceOf[CodingConvention](v_)
	})
}

// Java default interface method: CodingConvention.isExported(java.lang.String name)
func IfaceDefault_CodingConvention_IsExported_1(this_ CodingConvention, name string) bool {
	_ = this_
	_ = name
	return (this_.IsExported_2(name, true) || this_.IsExported_2(name, false))
}

// Java default interface method: CodingConvention.blockRenamingForProperty(java.lang.String name)
func IfaceDefault_CodingConvention_BlockRenamingForProperty(this_ CodingConvention, name string) bool {
	_ = this_
	_ = name
	return false
}

// Java class: com.google.javascript.jscomp.CodingConvention.Bind
type CodingConvention_Bind struct {
	target     pkg_rhino.Node `gson:"target"`
	thisValue  pkg_rhino.Node `gson:"thisValue"`
	parameters pkg_rhino.Node `gson:"parameters"`
	idHash_    atomic.Int32
}

func (this_ *CodingConvention_Bind) BindSelf_(s any) {
	_ = s
}
func (this_ *CodingConvention_Bind) IsJavaNil_() bool { return this_ == nil }
func (*CodingConvention_Bind) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CodingConvention$Bind", "Bind"
}
func (this_ *CodingConvention_Bind) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CodingConvention_Bind) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CodingConvention.Bind", "com.google.javascript.jscomp.CodingConvention$Bind", func(v_ any) bool {
		return j2gort.InstanceOf[*CodingConvention_Bind](v_)
	})
}

func (this_ *CodingConvention_Bind) initFields_CodingConvention_Bind() {
	_ = this_
}

// Java constructor: Bind(com.google.javascript.rhino.Node target, com.google.javascript.rhino.Node thisValue, com.google.javascript.rhino.Node parameters)
func NewCodingConvention_Bind(target pkg_rhino.Node, thisValue pkg_rhino.Node, parameters pkg_rhino.Node) *CodingConvention_Bind {
	this_ := &CodingConvention_Bind{}
	InitCodingConvention_Bind(this_, target, thisValue, parameters)
	return this_
}

func InitCodingConvention_Bind(this_ *CodingConvention_Bind, target pkg_rhino.Node, thisValue pkg_rhino.Node, parameters pkg_rhino.Node) {
	_ = target
	_ = thisValue
	_ = parameters
	this_.initFields_CodingConvention_Bind()
	this_.target = target
	this_.thisValue = thisValue
	this_.parameters = parameters
}

// Java method: Bind.getBoundParameterCount()
func (this_ *CodingConvention_Bind) GetBoundParameterCount() int32 {
	if this_.parameters == 0 {
		return int32(0)
	}
	var paramParent pkg_rhino.Node = this_.parameters.GetParent()
	_ = paramParent
	return (paramParent.GetChildCount() - paramParent.GetIndexOfChild(this_.parameters))
}

// Java class: com.google.javascript.jscomp.CodingConvention.Cache
type CodingConvention_Cache struct {
	cacheObj pkg_rhino.Node `gson:"cacheObj"`
	key      pkg_rhino.Node `gson:"key"`
	valueFn  pkg_rhino.Node `gson:"valueFn"`
	keyFn    pkg_rhino.Node `gson:"keyFn"`
	idHash_  atomic.Int32
}

func (this_ *CodingConvention_Cache) BindSelf_(s any) {
	_ = s
}
func (this_ *CodingConvention_Cache) IsJavaNil_() bool { return this_ == nil }
func (*CodingConvention_Cache) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CodingConvention$Cache", "Cache"
}
func (this_ *CodingConvention_Cache) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CodingConvention_Cache) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CodingConvention.Cache", "com.google.javascript.jscomp.CodingConvention$Cache", func(v_ any) bool {
		return j2gort.InstanceOf[*CodingConvention_Cache](v_)
	})
}

func (this_ *CodingConvention_Cache) initFields_CodingConvention_Cache() {
	_ = this_
}

// Java constructor: Cache(com.google.javascript.rhino.Node cacheObj, com.google.javascript.rhino.Node key, com.google.javascript.rhino.Node valueFn, com.google.javascript.rhino.Node keyFn)
func NewCodingConvention_Cache(cacheObj pkg_rhino.Node, key pkg_rhino.Node, valueFn pkg_rhino.Node, keyFn pkg_rhino.Node) *CodingConvention_Cache {
	this_ := &CodingConvention_Cache{}
	InitCodingConvention_Cache(this_, cacheObj, key, valueFn, keyFn)
	return this_
}

func InitCodingConvention_Cache(this_ *CodingConvention_Cache, cacheObj pkg_rhino.Node, key pkg_rhino.Node, valueFn pkg_rhino.Node, keyFn pkg_rhino.Node) {
	_ = cacheObj
	_ = key
	_ = valueFn
	_ = keyFn
	this_.initFields_CodingConvention_Cache()
	this_.cacheObj = cacheObj
	this_.key = key
	this_.valueFn = valueFn
	this_.keyFn = keyFn
}

// Java enum: com.google.javascript.jscomp.CodingConvention.SubclassType
type CodingConvention_SubclassType uint8

const (
	enum_CodingConvention_SubclassType_INHERITS CodingConvention_SubclassType = 1
	enum_CodingConvention_SubclassType_MIXIN    CodingConvention_SubclassType = 2
)

var enumBaseTable_CodingConvention_SubclassType = [...]j2gort.EnumBase{
	{Name_: "INHERITS", Ordinal_: 0},
	{Name_: "MIXIN", Ordinal_: 1},
}
var enumValues_CodingConvention_SubclassType = []CodingConvention_SubclassType{1, 2}

func StaticGet_CodingConvention_SubclassType_INHERITS() CodingConvention_SubclassType {
	return enum_CodingConvention_SubclassType_INHERITS
}
func StaticGet_CodingConvention_SubclassType_MIXIN() CodingConvention_SubclassType {
	return enum_CodingConvention_SubclassType_MIXIN
}

func CodingConvention_SubclassType_Values() []CodingConvention_SubclassType {
	return j2gort.CloneSlice(enumValues_CodingConvention_SubclassType)
}

func CodingConvention_SubclassType_ValueOf(name string) CodingConvention_SubclassType {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_CodingConvention_SubclassType {
		if enumBaseTable_CodingConvention_SubclassType[i_].Name_ == name {
			return CodingConvention_SubclassType(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.CodingConvention.SubclassType." + name))
}

func (this_ CodingConvention_SubclassType) IsJavaNil_() bool                              { return this_ == 0 }
func (CodingConvention_SubclassType) HasIdentityEquals_()                                 {}
func (CodingConvention_SubclassType) Implements_j2gort_Ext_java_lang_Enum()               {}
func (CodingConvention_SubclassType) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ CodingConvention_SubclassType) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (CodingConvention_SubclassType) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CodingConvention$SubclassType", "SubclassType"
}
func (this_ CodingConvention_SubclassType) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_CodingConvention_SubclassType[this_-1]
}
func (this_ CodingConvention_SubclassType) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ CodingConvention_SubclassType) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ CodingConvention_SubclassType) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CodingConvention_SubclassType[this_-1].Name_
}
func (this_ CodingConvention_SubclassType) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_CodingConvention_SubclassType[this_-1].Name_
}
func (this_ CodingConvention_SubclassType) Equals(other_ any) bool {
	o_, ok_ := other_.(CodingConvention_SubclassType)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ CodingConvention_SubclassType) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_CodingConvention_SubclassType[this_-1].HashCode()
}
func (this_ CodingConvention_SubclassType) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CodingConvention_SubclassType[this_-1].CompareTo(other_)
}
func (this_ CodingConvention_SubclassType) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_CodingConvention_SubclassType[this_-1].Name_)
}
func init() {
	vals_ := []any{CodingConvention_SubclassType(1), CodingConvention_SubclassType(2)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.CodingConvention.SubclassType", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.CodingConvention$SubclassType", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CodingConvention.SubclassType", "com.google.javascript.jscomp.CodingConvention$SubclassType", func(v_ any) bool {
		return j2gort.InstanceOf[CodingConvention_SubclassType](v_)
	})
}

// Java class: com.google.javascript.jscomp.CodingConvention.SubclassRelationship
type CodingConvention_SubclassRelationship struct {
	subclassName   string `gson:"subclassName"`
	superclassName string `gson:"superclassName"`
	idHash_        atomic.Int32
	type_          CodingConvention_SubclassType `gson:"type"`
}

func (this_ *CodingConvention_SubclassRelationship) BindSelf_(s any) {
	_ = s
}
func (this_ *CodingConvention_SubclassRelationship) IsJavaNil_() bool { return this_ == nil }
func (*CodingConvention_SubclassRelationship) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CodingConvention$SubclassRelationship", "SubclassRelationship"
}
func (this_ *CodingConvention_SubclassRelationship) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CodingConvention_SubclassRelationship) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CodingConvention.SubclassRelationship", "com.google.javascript.jscomp.CodingConvention$SubclassRelationship", func(v_ any) bool {
		return j2gort.InstanceOf[*CodingConvention_SubclassRelationship](v_)
	})
}

func (this_ *CodingConvention_SubclassRelationship) initFields_CodingConvention_SubclassRelationship() {
	_ = this_
	this_.subclassName = j2gort.NullString
	this_.superclassName = j2gort.NullString
}

// Java constructor: SubclassRelationship(com.google.javascript.jscomp.CodingConvention.SubclassType type, com.google.javascript.rhino.Node subclassNode, com.google.javascript.rhino.Node superclassNode)
func NewCodingConvention_SubclassRelationship(type_ CodingConvention_SubclassType, subclassNode pkg_rhino.Node, superclassNode pkg_rhino.Node) *CodingConvention_SubclassRelationship {
	this_ := &CodingConvention_SubclassRelationship{}
	InitCodingConvention_SubclassRelationship(this_, type_, subclassNode, superclassNode)
	return this_
}

func InitCodingConvention_SubclassRelationship(this_ *CodingConvention_SubclassRelationship, type_ CodingConvention_SubclassType, subclassNode pkg_rhino.Node, superclassNode pkg_rhino.Node) {
	_ = type_
	_ = subclassNode
	_ = superclassNode
	this_.initFields_CodingConvention_SubclassRelationship()
	if !(subclassNode.IsQualifiedName()) {
		j2gort.CheckArgumentFmt(false, "Expected qualified name, found: %s", j2gort.BoxNode(subclassNode))
	}
	if !(superclassNode.IsQualifiedName()) {
		j2gort.CheckArgumentFmt(false, "Expected qualified name, found: %s", j2gort.BoxNode(superclassNode))
	}
	this_.type_ = type_
	this_.subclassName = subclassNode.GetQualifiedName()
	this_.superclassName = superclassNode.GetQualifiedName()
}

// Java class: com.google.javascript.jscomp.CodingConvention.ObjectLiteralCast
type CodingConvention_ObjectLiteralCast struct {
	typeName       string          `gson:"typeName"`
	diagnosticType *DiagnosticType `gson:"diagnosticType"`
	objectNode     pkg_rhino.Node  `gson:"objectNode"`
	idHash_        atomic.Int32
}

func (this_ *CodingConvention_ObjectLiteralCast) BindSelf_(s any) {
	_ = s
}
func (this_ *CodingConvention_ObjectLiteralCast) IsJavaNil_() bool { return this_ == nil }
func (*CodingConvention_ObjectLiteralCast) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CodingConvention$ObjectLiteralCast", "ObjectLiteralCast"
}
func (this_ *CodingConvention_ObjectLiteralCast) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CodingConvention_ObjectLiteralCast) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CodingConvention.ObjectLiteralCast", "com.google.javascript.jscomp.CodingConvention$ObjectLiteralCast", func(v_ any) bool {
		return j2gort.InstanceOf[*CodingConvention_ObjectLiteralCast](v_)
	})
}

func (this_ *CodingConvention_ObjectLiteralCast) initFields_CodingConvention_ObjectLiteralCast() {
	_ = this_
	this_.typeName = j2gort.NullString
}

// Java constructor: ObjectLiteralCast(java.lang.String typeName, com.google.javascript.rhino.Node objectNode, com.google.javascript.jscomp.DiagnosticType diagnosticType)
func NewCodingConvention_ObjectLiteralCast(typeName string, objectNode pkg_rhino.Node, diagnosticType *DiagnosticType) *CodingConvention_ObjectLiteralCast {
	this_ := &CodingConvention_ObjectLiteralCast{}
	InitCodingConvention_ObjectLiteralCast(this_, typeName, objectNode, diagnosticType)
	return this_
}

func InitCodingConvention_ObjectLiteralCast(this_ *CodingConvention_ObjectLiteralCast, typeName string, objectNode pkg_rhino.Node, diagnosticType *DiagnosticType) {
	_ = typeName
	_ = objectNode
	_ = diagnosticType
	this_.initFields_CodingConvention_ObjectLiteralCast()
	this_.typeName = typeName
	this_.objectNode = objectNode
	this_.diagnosticType = diagnosticType
}

// Java polymorphic class interface: com.google.javascript.jscomp.CodingConvention.AssertionFunctionSpec
type CodingConvention_AssertionFunctionSpec interface {
	Implements_jscomp_CodingConvention_AssertionFunctionSpec()
	IsJavaNil_() bool
	GetCodingConvention_AssertionFunctionSpecBase() *CodingConvention_AssertionFunctionSpecBase
	GetFunctionName() string
	GetClosurePrimitive() pkg_rhino.ClosurePrimitive
	GetAssertionKind() CodingConvention_AssertionFunctionSpec_AssertionKind
	GetParamIndex() int32
	GetAssertedArg(pkg_rhino.Node) pkg_rhino.Node
}

func ToIface_CodingConvention_AssertionFunctionSpec[T interface {
	CodingConvention_AssertionFunctionSpec
	comparable
}](v T) CodingConvention_AssertionFunctionSpec {
	var z T
	if v == z {
		return nil
	}
	return v
}

// Java class: com.google.javascript.jscomp.CodingConvention.AssertionFunctionSpec
type CodingConvention_AssertionFunctionSpecBase struct {
	self_   CodingConvention_AssertionFunctionSpec
	idHash_ atomic.Int32
}

func (this_ *CodingConvention_AssertionFunctionSpecBase) BindSelf_(s any) {
	_ = s
	this_.self_ = s.(CodingConvention_AssertionFunctionSpec)
}
func (this_ *CodingConvention_AssertionFunctionSpecBase) IsJavaNil_() bool { return this_ == nil }
func (*CodingConvention_AssertionFunctionSpecBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CodingConvention$AssertionFunctionSpec", "AssertionFunctionSpec"
}
func (this_ *CodingConvention_AssertionFunctionSpecBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CodingConvention.AssertionFunctionSpec", "com.google.javascript.jscomp.CodingConvention$AssertionFunctionSpec", func(v_ any) bool {
		return j2gort.InstanceOf[CodingConvention_AssertionFunctionSpec](v_)
	})
}
func (b *CodingConvention_AssertionFunctionSpecBase) GetCodingConvention_AssertionFunctionSpecBase() *CodingConvention_AssertionFunctionSpecBase {
	return b
}
func (*CodingConvention_AssertionFunctionSpecBase) Implements_jscomp_CodingConvention_AssertionFunctionSpec() {
}

func (this_ *CodingConvention_AssertionFunctionSpecBase) initFields_CodingConvention_AssertionFunctionSpec() {
	_ = this_
}

func InitCodingConvention_AssertionFunctionSpec(this_ *CodingConvention_AssertionFunctionSpecBase) {
	this_.initFields_CodingConvention_AssertionFunctionSpec()
}

// Java method: AssertionFunctionSpec.getFunctionName()
func (this_ *CodingConvention_AssertionFunctionSpecBase) GetFunctionName() string {
	panic("abstract or native method invoked")
}

// Java method: AssertionFunctionSpec.getClosurePrimitive()
func (this_ *CodingConvention_AssertionFunctionSpecBase) GetClosurePrimitive() pkg_rhino.ClosurePrimitive {
	panic("abstract or native method invoked")
}

// Java method: AssertionFunctionSpec.getAssertionKind()
func (this_ *CodingConvention_AssertionFunctionSpecBase) GetAssertionKind() CodingConvention_AssertionFunctionSpec_AssertionKind {
	panic("abstract or native method invoked")
}

// Java method: AssertionFunctionSpec.getParamIndex()
func (this_ *CodingConvention_AssertionFunctionSpecBase) GetParamIndex() int32 {
	panic("abstract or native method invoked")
}

// Java static method: AssertionFunctionSpec.builder()
func CodingConvention_AssertionFunctionSpec_S_Builder() CodingConvention_AssertionFunctionSpec_Builder {
	return NewAutoValue_CodingConvention_AssertionFunctionSpec_Builder().SetParamIndex(int32(0))
}

// Java static method: AssertionFunctionSpec.forTruthy()
func CodingConvention_AssertionFunctionSpec_ForTruthy() CodingConvention_AssertionFunctionSpec_Builder {
	return CodingConvention_AssertionFunctionSpec_S_Builder().SetAssertionKind(StaticGet_CodingConvention_AssertionFunctionSpec_AssertionKind_TRUTHY())
}

// Java static method: AssertionFunctionSpec.forMatchesReturn()
func CodingConvention_AssertionFunctionSpec_ForMatchesReturn() CodingConvention_AssertionFunctionSpec_Builder {
	return CodingConvention_AssertionFunctionSpec_S_Builder().SetAssertionKind(StaticGet_CodingConvention_AssertionFunctionSpec_AssertionKind_MATCHES_RETURN_TYPE())
}

// Java method: AssertionFunctionSpec.getId()
func (this_ *CodingConvention_AssertionFunctionSpecBase) priv_AssertionFunctionSpec_getId_1387() any {
	if this_.self_.GetClosurePrimitive() != 0 {
		return j2gort.BoxNullableVal(this_.self_.GetClosurePrimitive())
	} else {
		return j2gort.BoxString(this_.self_.GetFunctionName())
	}
}

// Java method: AssertionFunctionSpec.getAssertedArg(com.google.javascript.rhino.Node firstArg)
func (this_ *CodingConvention_AssertionFunctionSpecBase) GetAssertedArg(firstArg pkg_rhino.Node) pkg_rhino.Node {
	_ = firstArg
	{
		var i int32 = int32(0)
		_ = i
		for ; i < this_.self_.GetParamIndex(); i = i + 1 {
			if firstArg == 0 {
				return pkg_rhino.Node(0)
			}
			firstArg = firstArg.GetNext()
		}
	}
	return firstArg
}

// Java enum: com.google.javascript.jscomp.CodingConvention.AssertionFunctionSpec.AssertionKind
type CodingConvention_AssertionFunctionSpec_AssertionKind uint8

const (
	enum_CodingConvention_AssertionFunctionSpec_AssertionKind_TRUTHY              CodingConvention_AssertionFunctionSpec_AssertionKind = 1
	enum_CodingConvention_AssertionFunctionSpec_AssertionKind_MATCHES_RETURN_TYPE CodingConvention_AssertionFunctionSpec_AssertionKind = 2
)

var enumBaseTable_CodingConvention_AssertionFunctionSpec_AssertionKind = [...]j2gort.EnumBase{
	{Name_: "TRUTHY", Ordinal_: 0},
	{Name_: "MATCHES_RETURN_TYPE", Ordinal_: 1},
}
var enumValues_CodingConvention_AssertionFunctionSpec_AssertionKind = []CodingConvention_AssertionFunctionSpec_AssertionKind{1, 2}

func StaticGet_CodingConvention_AssertionFunctionSpec_AssertionKind_TRUTHY() CodingConvention_AssertionFunctionSpec_AssertionKind {
	return enum_CodingConvention_AssertionFunctionSpec_AssertionKind_TRUTHY
}
func StaticGet_CodingConvention_AssertionFunctionSpec_AssertionKind_MATCHES_RETURN_TYPE() CodingConvention_AssertionFunctionSpec_AssertionKind {
	return enum_CodingConvention_AssertionFunctionSpec_AssertionKind_MATCHES_RETURN_TYPE
}

func CodingConvention_AssertionFunctionSpec_AssertionKind_Values() []CodingConvention_AssertionFunctionSpec_AssertionKind {
	return j2gort.CloneSlice(enumValues_CodingConvention_AssertionFunctionSpec_AssertionKind)
}

func CodingConvention_AssertionFunctionSpec_AssertionKind_ValueOf(name string) CodingConvention_AssertionFunctionSpec_AssertionKind {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_CodingConvention_AssertionFunctionSpec_AssertionKind {
		if enumBaseTable_CodingConvention_AssertionFunctionSpec_AssertionKind[i_].Name_ == name {
			return CodingConvention_AssertionFunctionSpec_AssertionKind(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.CodingConvention.AssertionFunctionSpec.AssertionKind." + name))
}

func (this_ CodingConvention_AssertionFunctionSpec_AssertionKind) IsJavaNil_() bool {
	return this_ == 0
}
func (CodingConvention_AssertionFunctionSpec_AssertionKind) HasIdentityEquals_()                   {}
func (CodingConvention_AssertionFunctionSpec_AssertionKind) Implements_j2gort_Ext_java_lang_Enum() {}
func (CodingConvention_AssertionFunctionSpec_AssertionKind) Implements_j2gort_Ext_java_lang_constant_Constable() {
}
func (this_ CodingConvention_AssertionFunctionSpec_AssertionKind) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (CodingConvention_AssertionFunctionSpec_AssertionKind) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CodingConvention$AssertionFunctionSpec$AssertionKind", "AssertionKind"
}
func (this_ CodingConvention_AssertionFunctionSpec_AssertionKind) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_CodingConvention_AssertionFunctionSpec_AssertionKind[this_-1]
}
func (this_ CodingConvention_AssertionFunctionSpec_AssertionKind) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ CodingConvention_AssertionFunctionSpec_AssertionKind) EnumOrdinal_() int32 {
	return int32(this_) - 1
}
func (this_ CodingConvention_AssertionFunctionSpec_AssertionKind) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CodingConvention_AssertionFunctionSpec_AssertionKind[this_-1].Name_
}
func (this_ CodingConvention_AssertionFunctionSpec_AssertionKind) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_CodingConvention_AssertionFunctionSpec_AssertionKind[this_-1].Name_
}
func (this_ CodingConvention_AssertionFunctionSpec_AssertionKind) Equals(other_ any) bool {
	o_, ok_ := other_.(CodingConvention_AssertionFunctionSpec_AssertionKind)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ CodingConvention_AssertionFunctionSpec_AssertionKind) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_CodingConvention_AssertionFunctionSpec_AssertionKind[this_-1].HashCode()
}
func (this_ CodingConvention_AssertionFunctionSpec_AssertionKind) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CodingConvention_AssertionFunctionSpec_AssertionKind[this_-1].CompareTo(other_)
}
func (this_ CodingConvention_AssertionFunctionSpec_AssertionKind) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_CodingConvention_AssertionFunctionSpec_AssertionKind[this_-1].Name_)
}
func init() {
	vals_ := []any{CodingConvention_AssertionFunctionSpec_AssertionKind(1), CodingConvention_AssertionFunctionSpec_AssertionKind(2)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.CodingConvention.AssertionFunctionSpec.AssertionKind", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.CodingConvention$AssertionFunctionSpec$AssertionKind", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CodingConvention.AssertionFunctionSpec.AssertionKind", "com.google.javascript.jscomp.CodingConvention$AssertionFunctionSpec$AssertionKind", func(v_ any) bool {
		return j2gort.InstanceOf[CodingConvention_AssertionFunctionSpec_AssertionKind](v_)
	})
}

// Java polymorphic class interface: com.google.javascript.jscomp.CodingConvention.AssertionFunctionSpec.Builder
type CodingConvention_AssertionFunctionSpec_Builder interface {
	Implements_jscomp_CodingConvention_AssertionFunctionSpec_Builder()
	IsJavaNil_() bool
	GetCodingConvention_AssertionFunctionSpec_BuilderBase() *CodingConvention_AssertionFunctionSpec_BuilderBase
	SetFunctionName(string) CodingConvention_AssertionFunctionSpec_Builder
	SetClosurePrimitive(pkg_rhino.ClosurePrimitive) CodingConvention_AssertionFunctionSpec_Builder
	SetParamIndex(int32) CodingConvention_AssertionFunctionSpec_Builder
	SetAssertionKind(CodingConvention_AssertionFunctionSpec_AssertionKind) CodingConvention_AssertionFunctionSpec_Builder
	AutoBuild() CodingConvention_AssertionFunctionSpec
	Build() CodingConvention_AssertionFunctionSpec
}

func ToIface_CodingConvention_AssertionFunctionSpec_Builder[T interface {
	CodingConvention_AssertionFunctionSpec_Builder
	comparable
}](v T) CodingConvention_AssertionFunctionSpec_Builder {
	var z T
	if v == z {
		return nil
	}
	return v
}

// Java class: com.google.javascript.jscomp.CodingConvention.AssertionFunctionSpec.Builder
type CodingConvention_AssertionFunctionSpec_BuilderBase struct {
	self_   CodingConvention_AssertionFunctionSpec_Builder
	idHash_ atomic.Int32
}

func (this_ *CodingConvention_AssertionFunctionSpec_BuilderBase) BindSelf_(s any) {
	_ = s
	this_.self_ = s.(CodingConvention_AssertionFunctionSpec_Builder)
}
func (this_ *CodingConvention_AssertionFunctionSpec_BuilderBase) IsJavaNil_() bool {
	return this_ == nil
}
func (*CodingConvention_AssertionFunctionSpec_BuilderBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CodingConvention$AssertionFunctionSpec$Builder", "Builder"
}
func (this_ *CodingConvention_AssertionFunctionSpec_BuilderBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CodingConvention.AssertionFunctionSpec.Builder", "com.google.javascript.jscomp.CodingConvention$AssertionFunctionSpec$Builder", func(v_ any) bool {
		return j2gort.InstanceOf[CodingConvention_AssertionFunctionSpec_Builder](v_)
	})
}
func (b *CodingConvention_AssertionFunctionSpec_BuilderBase) GetCodingConvention_AssertionFunctionSpec_BuilderBase() *CodingConvention_AssertionFunctionSpec_BuilderBase {
	return b
}
func (*CodingConvention_AssertionFunctionSpec_BuilderBase) Implements_jscomp_CodingConvention_AssertionFunctionSpec_Builder() {
}

func (this_ *CodingConvention_AssertionFunctionSpec_BuilderBase) initFields_CodingConvention_AssertionFunctionSpec_Builder() {
	_ = this_
}

func InitCodingConvention_AssertionFunctionSpec_Builder(this_ *CodingConvention_AssertionFunctionSpec_BuilderBase) {
	this_.initFields_CodingConvention_AssertionFunctionSpec_Builder()
}

// Java method: Builder.setFunctionName(java.lang.String name)
func (this_ *CodingConvention_AssertionFunctionSpec_BuilderBase) SetFunctionName(name string) CodingConvention_AssertionFunctionSpec_Builder {
	_ = name
	panic("abstract or native method invoked")
}

// Java method: Builder.setClosurePrimitive(com.google.javascript.rhino.ClosurePrimitive primitive)
func (this_ *CodingConvention_AssertionFunctionSpec_BuilderBase) SetClosurePrimitive(primitive pkg_rhino.ClosurePrimitive) CodingConvention_AssertionFunctionSpec_Builder {
	_ = primitive
	panic("abstract or native method invoked")
}

// Java method: Builder.setParamIndex(int paramIndex)
func (this_ *CodingConvention_AssertionFunctionSpec_BuilderBase) SetParamIndex(paramIndex int32) CodingConvention_AssertionFunctionSpec_Builder {
	_ = paramIndex
	panic("abstract or native method invoked")
}

// Java method: Builder.setAssertionKind(com.google.javascript.jscomp.CodingConvention.AssertionFunctionSpec.AssertionKind kind)
func (this_ *CodingConvention_AssertionFunctionSpec_BuilderBase) SetAssertionKind(kind CodingConvention_AssertionFunctionSpec_AssertionKind) CodingConvention_AssertionFunctionSpec_Builder {
	_ = kind
	panic("abstract or native method invoked")
}

// Java method: Builder.autoBuild()
func (this_ *CodingConvention_AssertionFunctionSpec_BuilderBase) AutoBuild() CodingConvention_AssertionFunctionSpec {
	panic("abstract or native method invoked")
}

// Java method: Builder.build()
func (this_ *CodingConvention_AssertionFunctionSpec_BuilderBase) Build() CodingConvention_AssertionFunctionSpec {
	var spec CodingConvention_AssertionFunctionSpec = this_.self_.AutoBuild()
	_ = spec
	if !((spec.GetFunctionName() != j2gort.NullString) || (spec.GetClosurePrimitive() != 0)) {
		j2gort.CheckStateStr(false, "Must provide a function name or ClosurePrimitive for each spec")
	}
	return spec
}

// Java class: com.google.javascript.jscomp.CodingConvention.AssertionFunctionLookup
type CodingConvention_AssertionFunctionLookup struct {
	internal            j2gort.Ext_com_google_common_collect_ImmutableMap `gson:"internal"`
	qnameHashes         []int32                                           `gson:"qnameHashes"`
	idHash_             atomic.Int32
	hasSimpleNamedSpecs bool `gson:"hasSimpleNamedSpecs"`
}

func (this_ *CodingConvention_AssertionFunctionLookup) BindSelf_(s any) {
	_ = s
}
func (this_ *CodingConvention_AssertionFunctionLookup) IsJavaNil_() bool { return this_ == nil }
func (*CodingConvention_AssertionFunctionLookup) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CodingConvention$AssertionFunctionLookup", "AssertionFunctionLookup"
}
func (this_ *CodingConvention_AssertionFunctionLookup) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CodingConvention_AssertionFunctionLookup) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CodingConvention.AssertionFunctionLookup", "com.google.javascript.jscomp.CodingConvention$AssertionFunctionLookup", func(v_ any) bool {
		return j2gort.InstanceOf[*CodingConvention_AssertionFunctionLookup](v_)
	})
}

func (this_ *CodingConvention_AssertionFunctionLookup) initFields_CodingConvention_AssertionFunctionLookup() {
	_ = this_
}

// Java constructor: AssertionFunctionLookup(com.google.common.collect.ImmutableMap<java.lang.Object,com.google.javascript.jscomp.CodingConvention.AssertionFunctionSpec> internal)
func NewCodingConvention_AssertionFunctionLookup(internal j2gort.Ext_com_google_common_collect_ImmutableMap) *CodingConvention_AssertionFunctionLookup {
	this_ := &CodingConvention_AssertionFunctionLookup{}
	InitCodingConvention_AssertionFunctionLookup(this_, internal)
	return this_
}

func InitCodingConvention_AssertionFunctionLookup(this_ *CodingConvention_AssertionFunctionLookup, internal j2gort.Ext_com_google_common_collect_ImmutableMap) {
	_ = internal
	this_.initFields_CodingConvention_AssertionFunctionLookup()
	this_.internal = internal
	var hasSimple bool = false
	_ = hasSimple
	var dottedCount int32 = int32(0)
	_ = dottedCount
	{
		cur_1 := j2gort.NewMapKeyCursor(internal)
		for cur_1.HasNext() {
			rawElem_2 := cur_1.Next()
			id := rawElem_2
			_ = id
			var name string
			_ = name
			var ok_3 bool
			name, ok_3 = j2gort.TryCastString(id)
			if ok_3 {
				if j2gort.JStrIndexOfCp(name, int32(uint16(46))) >= int32(0) {
					dottedCount = dottedCount + 1
				} else {
					hasSimple = true
				}
			}
		}
	}
	this_.hasSimpleNamedSpecs = hasSimple
	if dottedCount == int32(0) {
		this_.qnameHashes = nil
	} else {
		var cap_ int32 = int32(8)
		_ = cap_
		for cap_ < (dottedCount * int32(2)) {
			cap_ = javalang.Shl32(cap_, int(int32(1)))
		}
		var table []int32 = make([]int32, cap_)
		_ = table
		var mask int32 = (cap_ - int32(1))
		_ = mask
		{
			cur_4 := j2gort.NewMapKeyCursor(internal)
			for cur_4.HasNext() {
				rawElem_5 := cur_4.Next()
				id_2 := rawElem_5
				_ = id_2
				var name_2 string
				_ = name_2
				var ok_6 bool
				name_2, ok_6 = j2gort.TryCastString(id_2)
				if ok_6 && (j2gort.JStrIndexOfCp(name_2, int32(uint16(46))) >= int32(0)) {
					var h int32 = pkg_rhino.Node_HashQualifiedNameString(name_2)
					_ = h
					if h != int32(0) {
						var slot int32 = (h & mask)
						_ = slot
						for (table[slot] != int32(0)) && (table[slot] != h) {
							slot = ((slot + int32(1)) & mask)
						}
						{
							a_7 := table
							i_8 := slot
							a_7[i_8] = h
						}
					}
				}
			}
		}
		this_.qnameHashes = table
	}
}

// Java static method: AssertionFunctionLookup.containsHash(int[] table, int h)
func CodingConvention_AssertionFunctionLookup_ContainsHash(table []int32, h int32) bool {
	_ = table
	_ = h
	var mask int32 = (int32(len(table)) - int32(1))
	_ = mask
	var slot int32 = (h & mask)
	_ = slot
	for {
		var existing int32 = table[slot]
		_ = existing
		if existing == int32(0) {
			return false
		}
		if existing == h {
			return true
		}
		slot = ((slot + int32(1)) & mask)
	}
	panic("unreachable")
}

// Java static method: AssertionFunctionLookup.of(java.util.Collection<com.google.javascript.jscomp.CodingConvention.AssertionFunctionSpec> specs)
func CodingConvention_AssertionFunctionLookup_Of(specs j2gort.Ext_java_util_Collection) *CodingConvention_AssertionFunctionLookup {
	_ = specs
	var idToSpecMap j2gort.Ext_com_google_common_collect_ImmutableMap = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap](specs.Stream().Collect_1(j2gort.Ext_com_google_common_collect_ImmutableMap_ToImmutableMap_2(lambdaSingleton_coding_convention_1, j2gort.Ext_java_util_function_Function_Identity())))
	_ = idToSpecMap
	return NewCodingConvention_AssertionFunctionLookup(idToSpecMap)
}

// Java method: AssertionFunctionLookup.lookupByCallee(com.google.javascript.rhino.Node callee)
func (this_ *CodingConvention_AssertionFunctionLookup) LookupByCallee(callee pkg_rhino.Node) CodingConvention_AssertionFunctionSpec {
	_ = callee
	var fnType *pkg_rhino.FunctionType = pkg_rhino.JSType_ToMaybeFunctionType(callee.GetJSType())
	_ = fnType
	if (fnType != nil) && (fnType.GetClosurePrimitive() != 0) {
		var spec CodingConvention_AssertionFunctionSpec = j2gort.AsIface[CodingConvention_AssertionFunctionSpec](this_.internal.Get(j2gort.BoxNullableVal(fnType.GetClosurePrimitive())))
		_ = spec
		if spec != nil && !spec.IsJavaNil_() {
			return spec
		}
	}
	if callee.IsGetProp() {
		var table []int32 = this_.qnameHashes
		_ = table
		if table == nil {
			return nil
		}
		var h int32 = callee.HashQualifiedName()
		_ = h
		if (h == int32(0)) || !(CodingConvention_AssertionFunctionLookup_ContainsHash(table, h)) {
			return nil
		}
		return j2gort.AsIface[CodingConvention_AssertionFunctionSpec](j2gort.MapGetStr(this_.internal, callee.GetQualifiedName()))
	}
	if callee.IsName() {
		var name string = callee.GetString()
		_ = name
		if j2gort.CheckNotNullString(name) == "" {
			return nil
		}
		if j2gort.JStrIndexOfCp(name, int32(uint16(46))) >= int32(0) {
			if (this_.qnameHashes != nil) && CodingConvention_AssertionFunctionLookup_ContainsHash(this_.qnameHashes, pkg_rhino.Node_HashQualifiedNameString(name)) {
				return j2gort.AsIface[CodingConvention_AssertionFunctionSpec](j2gort.MapGetStr(this_.internal, name))
			} else {
				return nil
			}
		}
		if this_.hasSimpleNamedSpecs {
			return j2gort.AsIface[CodingConvention_AssertionFunctionSpec](j2gort.MapGetStr(this_.internal, name))
		} else {
			return nil
		}
	}
	if this_.hasSimpleNamedSpecs && callee.IsQualifiedName() {
		return j2gort.AsIface[CodingConvention_AssertionFunctionSpec](j2gort.MapGetStr(this_.internal, callee.GetQualifiedName()))
	}
	return nil
}
