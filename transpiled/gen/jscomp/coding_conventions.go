/*
 * Copyright 2011 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/CodingConventions.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.CodingConventions
type CodingConventions struct {
	idHash_ atomic.Int32
}

func (this_ *CodingConventions) BindSelf_(s any) {
	_ = s
}
func (this_ *CodingConventions) IsJavaNil_() bool { return this_ == nil }
func (*CodingConventions) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CodingConventions", "CodingConventions"
}
func (this_ *CodingConventions) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CodingConventions) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CodingConventions", "com.google.javascript.jscomp.CodingConventions", func(v_ any) bool {
		return j2gort.InstanceOf[*CodingConventions](v_)
	})
}

func (this_ *CodingConventions) initFields_CodingConventions() {
	_ = this_
}

// Java constructor: CodingConventions()
func NewCodingConventions() *CodingConventions {
	this_ := &CodingConventions{}
	InitCodingConventions(this_)
	return this_
}

func InitCodingConventions(this_ *CodingConventions) {
	this_.initFields_CodingConventions()
}

// Java static method: CodingConventions.getDefault()
func CodingConventions_GetDefault() CodingConvention {
	return NewCodingConventions_DefaultCodingConvention()
}

// Java static method: CodingConventions.defaultIsFunctionCallThatAlwaysThrows(com.google.javascript.rhino.Node n, java.lang.String alwaysThrowsFunctionName)
func CodingConventions_DefaultIsFunctionCallThatAlwaysThrows(n pkg_rhino.Node, alwaysThrowsFunctionName string) bool {
	_ = n
	_ = alwaysThrowsFunctionName
	if n.IsExprResult() {
		if !(n.GetFirstChild().IsCall()) {
			return false
		}
	} else if !(n.IsCall()) {
		return false
	}
	if n.IsExprResult() {
		n = n.GetFirstChild()
	}
	return n.GetFirstChild().MatchesQualifiedName_1_string(alwaysThrowsFunctionName)
}

// Java polymorphic class interface: com.google.javascript.jscomp.CodingConventions.Proxy
type CodingConventions_Proxy interface {
	Implements_jscomp_CodingConventions_Proxy()
	IsJavaNil_() bool
	GetCodingConventions_ProxyBase() *CodingConventions_ProxyBase
	CodingConvention
}

func ToIface_CodingConventions_Proxy[T interface {
	CodingConventions_Proxy
	comparable
}](v T) CodingConventions_Proxy {
	var z T
	if v == z {
		return nil
	}
	return v
}

// Java class: com.google.javascript.jscomp.CodingConventions.Proxy
type CodingConventions_ProxyBase struct {
	self_          CodingConventions_Proxy
	nextConvention CodingConvention `gson:"nextConvention"`
	idHash_        atomic.Int32
}

func (this_ *CodingConventions_ProxyBase) BindSelf_(s any) {
	_ = s
	this_.self_ = s.(CodingConventions_Proxy)
}
func (this_ *CodingConventions_ProxyBase) IsJavaNil_() bool { return this_ == nil }
func (*CodingConventions_ProxyBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CodingConventions$Proxy", "Proxy"
}
func (this_ *CodingConventions_ProxyBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CodingConventions.Proxy", "com.google.javascript.jscomp.CodingConventions$Proxy", func(v_ any) bool {
		return j2gort.InstanceOf[CodingConventions_Proxy](v_)
	})
}
func (b *CodingConventions_ProxyBase) GetCodingConventions_ProxyBase() *CodingConventions_ProxyBase {
	return b
}
func (*CodingConventions_ProxyBase) Implements_jscomp_CodingConventions_Proxy() {}
func (*CodingConventions_ProxyBase) Implements_jscomp_CodingConvention()        {}

func (this_ *CodingConventions_ProxyBase) initFields_CodingConventions_Proxy() {
	_ = this_
}

// Java constructor: Proxy(com.google.javascript.jscomp.CodingConvention convention)
func NewCodingConventions_Proxy(convention CodingConvention) CodingConventions_Proxy {
	this_ := &CodingConventions_ProxyBase{}
	this_.BindSelf_(this_)
	InitCodingConventions_Proxy(this_, convention)
	return this_
}

func InitCodingConventions_Proxy(this_ *CodingConventions_ProxyBase, convention CodingConvention) {
	_ = convention
	this_.initFields_CodingConventions_Proxy()
	this_.nextConvention = convention
}

// Java method: Proxy.isConstant(java.lang.String variableName)
func (this_ *CodingConventions_ProxyBase) IsConstant(variableName string) bool {
	_ = variableName
	return this_.nextConvention.IsConstant(variableName)
}

// Java method: Proxy.isConstantKey(java.lang.String keyName)
func (this_ *CodingConventions_ProxyBase) IsConstantKey(keyName string) bool {
	_ = keyName
	return this_.nextConvention.IsConstantKey(keyName)
}

// Java method: Proxy.isValidEnumKey(java.lang.String key)
func (this_ *CodingConventions_ProxyBase) IsValidEnumKey(key string) bool {
	_ = key
	return this_.nextConvention.IsValidEnumKey(key)
}

// Java method: Proxy.isOptionalParameter(com.google.javascript.rhino.Node parameter)
func (this_ *CodingConventions_ProxyBase) IsOptionalParameter(parameter pkg_rhino.Node) bool {
	_ = parameter
	return this_.nextConvention.IsOptionalParameter(parameter)
}

// Java method: Proxy.isVarArgsParameter(com.google.javascript.rhino.Node parameter)
func (this_ *CodingConventions_ProxyBase) IsVarArgsParameter(parameter pkg_rhino.Node) bool {
	_ = parameter
	return this_.nextConvention.IsVarArgsParameter(parameter)
}

// Java method: Proxy.isFunctionCallThatAlwaysThrows(com.google.javascript.rhino.Node n)
func (this_ *CodingConventions_ProxyBase) IsFunctionCallThatAlwaysThrows(n pkg_rhino.Node) bool {
	_ = n
	return this_.nextConvention.IsFunctionCallThatAlwaysThrows(n)
}

// Java method: Proxy.isExported(java.lang.String name, boolean local)
func (this_ *CodingConventions_ProxyBase) IsExported_2(name string, local bool) bool {
	_ = name
	_ = local
	return this_.nextConvention.IsExported_2(name, local)
}

// Java method: Proxy.isExported(java.lang.String name)
func (this_ *CodingConventions_ProxyBase) IsExported_1(name string) bool {
	_ = name
	return IfaceDefault_CodingConvention_IsExported_1(j2gort.AsIface[CodingConvention](this_.self_), name)
}

// Java method: Proxy.getPackageName(com.google.javascript.rhino.StaticSourceFile source)
func (this_ *CodingConventions_ProxyBase) GetPackageName(source pkg_rhino.StaticSourceFile) string {
	_ = source
	return this_.nextConvention.GetPackageName(source)
}

// Java method: Proxy.getClassesDefinedByCall(com.google.javascript.rhino.Node callNode)
func (this_ *CodingConventions_ProxyBase) GetClassesDefinedByCall(callNode pkg_rhino.Node) *CodingConvention_SubclassRelationship {
	_ = callNode
	return this_.nextConvention.GetClassesDefinedByCall(callNode)
}

// Java method: Proxy.isClassFactoryCall(com.google.javascript.rhino.Node callNode)
func (this_ *CodingConventions_ProxyBase) IsClassFactoryCall(callNode pkg_rhino.Node) bool {
	_ = callNode
	return this_.nextConvention.IsClassFactoryCall(callNode)
}

// Java method: Proxy.isSuperClassReference(java.lang.String propertyName)
func (this_ *CodingConventions_ProxyBase) IsSuperClassReference(propertyName string) bool {
	_ = propertyName
	return this_.nextConvention.IsSuperClassReference(propertyName)
}

// Java method: Proxy.extractIsModuleFile(com.google.javascript.rhino.Node node, com.google.javascript.rhino.Node parent)
func (this_ *CodingConventions_ProxyBase) ExtractIsModuleFile(node pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = node
	_ = parent
	return this_.nextConvention.ExtractIsModuleFile(node, parent)
}

// Java method: Proxy.extractClassNameIfProvide(com.google.javascript.rhino.Node node, com.google.javascript.rhino.Node parent)
func (this_ *CodingConventions_ProxyBase) ExtractClassNameIfProvide(node pkg_rhino.Node, parent pkg_rhino.Node) string {
	_ = node
	_ = parent
	return this_.nextConvention.ExtractClassNameIfProvide(node, parent)
}

// Java method: Proxy.extractClassNameIfRequire(com.google.javascript.rhino.Node node, com.google.javascript.rhino.Node parent)
func (this_ *CodingConventions_ProxyBase) ExtractClassNameIfRequire(node pkg_rhino.Node, parent pkg_rhino.Node) string {
	_ = node
	_ = parent
	return this_.nextConvention.ExtractClassNameIfRequire(node, parent)
}

// Java method: Proxy.getExportPropertyFunction()
func (this_ *CodingConventions_ProxyBase) GetExportPropertyFunction() string {
	return this_.nextConvention.GetExportPropertyFunction()
}

// Java method: Proxy.getExportSymbolFunction()
func (this_ *CodingConventions_ProxyBase) GetExportSymbolFunction() string {
	return this_.nextConvention.GetExportSymbolFunction()
}

// Java method: Proxy.identifyTypeDeclarationCall(com.google.javascript.rhino.Node n)
func (this_ *CodingConventions_ProxyBase) IdentifyTypeDeclarationCall(n pkg_rhino.Node) j2gort.Ext_java_util_List {
	_ = n
	return this_.nextConvention.IdentifyTypeDeclarationCall(n)
}

// Java method: Proxy.applySubclassRelationship(com.google.javascript.rhino.NominalTypeBuilder parent, com.google.javascript.rhino.NominalTypeBuilder child, com.google.javascript.jscomp.CodingConvention.SubclassType type)
func (this_ *CodingConventions_ProxyBase) ApplySubclassRelationship(parent *pkg_rhino.NominalTypeBuilder, child *pkg_rhino.NominalTypeBuilder, type_ CodingConvention_SubclassType) {
	_ = parent
	_ = child
	_ = type_
	this_.nextConvention.ApplySubclassRelationship(parent, child, type_)
}

// Java method: Proxy.getAbstractMethodName()
func (this_ *CodingConventions_ProxyBase) GetAbstractMethodName() string {
	return this_.nextConvention.GetAbstractMethodName()
}

// Java method: Proxy.getSingletonGetterClassName(com.google.javascript.rhino.Node callNode)
func (this_ *CodingConventions_ProxyBase) GetSingletonGetterClassName(callNode pkg_rhino.Node) string {
	_ = callNode
	return this_.nextConvention.GetSingletonGetterClassName(callNode)
}

// Java method: Proxy.applySingletonGetter(com.google.javascript.rhino.NominalTypeBuilder classType, com.google.javascript.rhino.jstype.FunctionType getterType)
func (this_ *CodingConventions_ProxyBase) ApplySingletonGetter(classType *pkg_rhino.NominalTypeBuilder, getterType *pkg_rhino.FunctionType) {
	_ = classType
	_ = getterType
	this_.nextConvention.ApplySingletonGetter(classType, getterType)
}

// Java method: Proxy.getAssertionFunctions()
func (this_ *CodingConventions_ProxyBase) GetAssertionFunctions() any {
	return j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.nextConvention.GetAssertionFunctions())
}

// Java method: Proxy.describeFunctionBind(com.google.javascript.rhino.Node n, boolean checkTypes)
func (this_ *CodingConventions_ProxyBase) DescribeFunctionBind(n pkg_rhino.Node, checkTypes bool) *CodingConvention_Bind {
	_ = n
	_ = checkTypes
	return this_.nextConvention.DescribeFunctionBind(n, checkTypes)
}

// Java method: Proxy.describeCachingCall(com.google.javascript.rhino.Node node)
func (this_ *CodingConventions_ProxyBase) DescribeCachingCall(node pkg_rhino.Node) *CodingConvention_Cache {
	_ = node
	return this_.nextConvention.DescribeCachingCall(node)
}

// Java method: Proxy.isPropertyTestFunction(com.google.javascript.rhino.Node call)
func (this_ *CodingConventions_ProxyBase) IsPropertyTestFunction(call pkg_rhino.Node) bool {
	_ = call
	return this_.nextConvention.IsPropertyTestFunction(call)
}

// Java method: Proxy.isPropertyRenameFunction(com.google.javascript.rhino.Node nameNode)
func (this_ *CodingConventions_ProxyBase) IsPropertyRenameFunction(nameNode pkg_rhino.Node) bool {
	_ = nameNode
	return this_.nextConvention.IsPropertyRenameFunction(nameNode)
}

// Java method: Proxy.isPrototypeAlias(com.google.javascript.rhino.Node getProp)
func (this_ *CodingConventions_ProxyBase) IsPrototypeAlias(getProp pkg_rhino.Node) bool {
	_ = getProp
	return false
}

// Java method: Proxy.getObjectLiteralCast(com.google.javascript.rhino.Node callNode)
func (this_ *CodingConventions_ProxyBase) GetObjectLiteralCast(callNode pkg_rhino.Node) *CodingConvention_ObjectLiteralCast {
	_ = callNode
	return this_.nextConvention.GetObjectLiteralCast(callNode)
}

func (this_ *CodingConventions_ProxyBase) BlockRenamingForProperty(p0_ string) bool {
	return IfaceDefault_CodingConvention_BlockRenamingForProperty(this_.self_, p0_)
}

// Java class: com.google.javascript.jscomp.CodingConventions.DefaultCodingConvention
type CodingConventions_DefaultCodingConvention struct {
	idHash_ atomic.Int32
}

func (this_ *CodingConventions_DefaultCodingConvention) BindSelf_(s any) {
	_ = s
}
func (this_ *CodingConventions_DefaultCodingConvention) IsJavaNil_() bool { return this_ == nil }
func (*CodingConventions_DefaultCodingConvention) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CodingConventions$DefaultCodingConvention", "DefaultCodingConvention"
}
func (this_ *CodingConventions_DefaultCodingConvention) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CodingConventions_DefaultCodingConvention) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CodingConventions.DefaultCodingConvention", "com.google.javascript.jscomp.CodingConventions$DefaultCodingConvention", func(v_ any) bool {
		return j2gort.InstanceOf[*CodingConventions_DefaultCodingConvention](v_)
	})
}
func (*CodingConventions_DefaultCodingConvention) Implements_jscomp_CodingConvention() {}

func StaticGet_CodingConventions_DefaultCodingConvention_SerialVersionUID() int64 { return int64(1) }

var static_CodingConventions_DefaultCodingConvention_JSCOMP_INHERITS pkg_rhino.QualifiedName

func StaticGet_CodingConventions_DefaultCodingConvention_JSCOMP_INHERITS() pkg_rhino.QualifiedName {
	clinit_CodingConventions_DefaultCodingConvention()
	return static_CodingConventions_DefaultCodingConvention_JSCOMP_INHERITS
}
func StaticSet_CodingConventions_DefaultCodingConvention_JSCOMP_INHERITS(v pkg_rhino.QualifiedName) {
	clinit_CodingConventions_DefaultCodingConvention()
	static_CodingConventions_DefaultCodingConvention_JSCOMP_INHERITS = v
}

var static_CodingConventions_DefaultCodingConvention_FUNCTION_PROTOTYPE_BIND_CALL pkg_rhino.QualifiedName

func StaticGet_CodingConventions_DefaultCodingConvention_FUNCTION_PROTOTYPE_BIND_CALL() pkg_rhino.QualifiedName {
	clinit_CodingConventions_DefaultCodingConvention()
	return static_CodingConventions_DefaultCodingConvention_FUNCTION_PROTOTYPE_BIND_CALL
}
func StaticSet_CodingConventions_DefaultCodingConvention_FUNCTION_PROTOTYPE_BIND_CALL(v pkg_rhino.QualifiedName) {
	clinit_CodingConventions_DefaultCodingConvention()
	static_CodingConventions_DefaultCodingConvention_FUNCTION_PROTOTYPE_BIND_CALL = v
}

var clinitDone_CodingConventions_DefaultCodingConvention uint32
var clinitState_CodingConventions_DefaultCodingConvention j2gort.ClinitState

func clinit_CodingConventions_DefaultCodingConvention() {
	if atomic.LoadUint32(&clinitDone_CodingConventions_DefaultCodingConvention) != 0 {
		return
	}
	clinitSlow_CodingConventions_DefaultCodingConvention()
}

func clinitSlow_CodingConventions_DefaultCodingConvention() {
	if !j2gort.EnterClinit(&clinitState_CodingConventions_DefaultCodingConvention, &clinitDone_CodingConventions_DefaultCodingConvention) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_CodingConventions_DefaultCodingConvention, &clinitDone_CodingConventions_DefaultCodingConvention)
	static_CodingConventions_DefaultCodingConvention_JSCOMP_INHERITS = pkg_rhino.QualifiedName_Of("$jscomp.inherits")
	static_CodingConventions_DefaultCodingConvention_FUNCTION_PROTOTYPE_BIND_CALL = pkg_rhino.QualifiedName_Of("Function.prototype.bind.call")
}

func (this_ *CodingConventions_DefaultCodingConvention) initFields_CodingConventions_DefaultCodingConvention() {
	_ = this_
}

// Java constructor: DefaultCodingConvention()
func NewCodingConventions_DefaultCodingConvention() *CodingConventions_DefaultCodingConvention {
	this_ := &CodingConventions_DefaultCodingConvention{}
	InitCodingConventions_DefaultCodingConvention(this_)
	return this_
}

func InitCodingConventions_DefaultCodingConvention(this_ *CodingConventions_DefaultCodingConvention) {
	this_.initFields_CodingConventions_DefaultCodingConvention()
}

// Java method: DefaultCodingConvention.isConstant(java.lang.String variableName)
func (this_ *CodingConventions_DefaultCodingConvention) IsConstant(variableName string) bool {
	_ = variableName
	return false
}

// Java method: DefaultCodingConvention.isConstantKey(java.lang.String variableName)
func (this_ *CodingConventions_DefaultCodingConvention) IsConstantKey(variableName string) bool {
	_ = variableName
	return false
}

// Java method: DefaultCodingConvention.isValidEnumKey(java.lang.String key)
func (this_ *CodingConventions_DefaultCodingConvention) IsValidEnumKey(key string) bool {
	_ = key
	return ((key != j2gort.NullString) && (j2gort.JStrLen(key) > int32(0)))
}

// Java method: DefaultCodingConvention.isOptionalParameter(com.google.javascript.rhino.Node parameter)
func (this_ *CodingConventions_DefaultCodingConvention) IsOptionalParameter(parameter pkg_rhino.Node) bool {
	_ = parameter
	return false
}

// Java method: DefaultCodingConvention.isVarArgsParameter(com.google.javascript.rhino.Node parameter)
func (this_ *CodingConventions_DefaultCodingConvention) IsVarArgsParameter(parameter pkg_rhino.Node) bool {
	_ = parameter
	return parameter.IsRest()
}

// Java method: DefaultCodingConvention.isFunctionCallThatAlwaysThrows(com.google.javascript.rhino.Node n)
func (this_ *CodingConventions_DefaultCodingConvention) IsFunctionCallThatAlwaysThrows(n pkg_rhino.Node) bool {
	_ = n
	if NodeUtil_IsExprCall(n) {
		var fnType *pkg_rhino.FunctionType = pkg_rhino.JSType_ToMaybeFunctionType(n.GetFirstFirstChild().GetJSType())
		_ = fnType
		return ((fnType != nil) && (pkg_rhino.StaticGet_ClosurePrimitive_ASSERTS_FAIL() == fnType.GetClosurePrimitive()))
	}
	return false
}

// Java method: DefaultCodingConvention.getPackageName(com.google.javascript.rhino.StaticSourceFile source)
func (this_ *CodingConventions_DefaultCodingConvention) GetPackageName(source pkg_rhino.StaticSourceFile) string {
	_ = source
	var name string = source.GetName()
	_ = name
	var lastSlash int32 = j2gort.JStrLastIndexOfCp(name, int32(uint16(47)))
	_ = lastSlash
	if lastSlash == int32(-1) {
		return ""
	} else {
		return j2gort.JStrSub(name, int32(0), lastSlash)
	}
}

// Java method: DefaultCodingConvention.isExported(java.lang.String name, boolean local)
func (this_ *CodingConventions_DefaultCodingConvention) IsExported_2(name string, local bool) bool {
	_ = name
	_ = local
	return (local && javalang.StringStartsWith(j2gort.CheckNotNullString(name), "$super"))
}

// Java method: DefaultCodingConvention.isExported(java.lang.String name)
func (this_ *CodingConventions_DefaultCodingConvention) IsExported_1(name string) bool {
	_ = name
	return IfaceDefault_CodingConvention_IsExported_1(this_, name)
}

// Java method: DefaultCodingConvention.getClassesDefinedByCall(com.google.javascript.rhino.Node callNode)
func (this_ *CodingConventions_DefaultCodingConvention) GetClassesDefinedByCall(callNode pkg_rhino.Node) *CodingConvention_SubclassRelationship {
	_ = callNode
	var callName pkg_rhino.Node = callNode.GetFirstChild()
	_ = callName
	if (StaticGet_CodingConventions_DefaultCodingConvention_JSCOMP_INHERITS().Matches(callName) || callName.MatchesName_1_string("$jscomp$inherits")) && callNode.HasXChildren(int32(3)) {
		var subclass pkg_rhino.Node = callName.GetNext()
		_ = subclass
		var superclass pkg_rhino.Node = subclass.GetNext()
		_ = superclass
		if subclass.IsQualifiedName() && superclass.IsQualifiedName() {
			return NewCodingConvention_SubclassRelationship(StaticGet_CodingConvention_SubclassType_INHERITS(), subclass, superclass)
		}
	}
	return nil
}

// Java method: DefaultCodingConvention.isClassFactoryCall(com.google.javascript.rhino.Node callNode)
func (this_ *CodingConventions_DefaultCodingConvention) IsClassFactoryCall(callNode pkg_rhino.Node) bool {
	_ = callNode
	return false
}

// Java method: DefaultCodingConvention.isSuperClassReference(java.lang.String propertyName)
func (this_ *CodingConventions_DefaultCodingConvention) IsSuperClassReference(propertyName string) bool {
	_ = propertyName
	return false
}

// Java method: DefaultCodingConvention.extractIsModuleFile(com.google.javascript.rhino.Node node, com.google.javascript.rhino.Node parent)
func (this_ *CodingConventions_DefaultCodingConvention) ExtractIsModuleFile(node pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = node
	_ = parent
	var message string = "only implemented in ClosureCodingConvention"
	_ = message
	panic(j2gort.AsIface[javalang.IUnsupportedOperationException](j2gort.ConstructShimThrowable("java.lang.UnsupportedOperationException", j2gort.BoxString(message))))
}

// Java method: DefaultCodingConvention.extractClassNameIfProvide(com.google.javascript.rhino.Node node, com.google.javascript.rhino.Node parent)
func (this_ *CodingConventions_DefaultCodingConvention) ExtractClassNameIfProvide(node pkg_rhino.Node, parent pkg_rhino.Node) string {
	_ = node
	_ = parent
	var message string = "only implemented in ClosureCodingConvention"
	_ = message
	panic(j2gort.AsIface[javalang.IUnsupportedOperationException](j2gort.ConstructShimThrowable("java.lang.UnsupportedOperationException", j2gort.BoxString(message))))
}

// Java method: DefaultCodingConvention.extractClassNameIfRequire(com.google.javascript.rhino.Node node, com.google.javascript.rhino.Node parent)
func (this_ *CodingConventions_DefaultCodingConvention) ExtractClassNameIfRequire(node pkg_rhino.Node, parent pkg_rhino.Node) string {
	_ = node
	_ = parent
	var message string = "only implemented in ClosureCodingConvention"
	_ = message
	panic(j2gort.AsIface[javalang.IUnsupportedOperationException](j2gort.ConstructShimThrowable("java.lang.UnsupportedOperationException", j2gort.BoxString(message))))
}

// Java method: DefaultCodingConvention.getExportPropertyFunction()
func (this_ *CodingConventions_DefaultCodingConvention) GetExportPropertyFunction() string {
	return j2gort.NullString
}

// Java method: DefaultCodingConvention.getExportSymbolFunction()
func (this_ *CodingConventions_DefaultCodingConvention) GetExportSymbolFunction() string {
	return j2gort.NullString
}

// Java method: DefaultCodingConvention.identifyTypeDeclarationCall(com.google.javascript.rhino.Node n)
func (this_ *CodingConventions_DefaultCodingConvention) IdentifyTypeDeclarationCall(n pkg_rhino.Node) j2gort.Ext_java_util_List {
	_ = n
	return j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.Ext_com_google_common_collect_ImmutableList_Of_0())
}

// Java method: DefaultCodingConvention.applySubclassRelationship(com.google.javascript.rhino.NominalTypeBuilder parent, com.google.javascript.rhino.NominalTypeBuilder child, com.google.javascript.jscomp.CodingConvention.SubclassType type)
func (this_ *CodingConventions_DefaultCodingConvention) ApplySubclassRelationship(parent *pkg_rhino.NominalTypeBuilder, child *pkg_rhino.NominalTypeBuilder, type_ CodingConvention_SubclassType) {
	_ = parent
	_ = child
	_ = type_
}

// Java method: DefaultCodingConvention.getAbstractMethodName()
func (this_ *CodingConventions_DefaultCodingConvention) GetAbstractMethodName() string {
	return j2gort.NullString
}

// Java method: DefaultCodingConvention.getSingletonGetterClassName(com.google.javascript.rhino.Node callNode)
func (this_ *CodingConventions_DefaultCodingConvention) GetSingletonGetterClassName(callNode pkg_rhino.Node) string {
	_ = callNode
	return j2gort.NullString
}

// Java method: DefaultCodingConvention.applySingletonGetter(com.google.javascript.rhino.NominalTypeBuilder classType, com.google.javascript.rhino.jstype.FunctionType getterType)
func (this_ *CodingConventions_DefaultCodingConvention) ApplySingletonGetter(classType *pkg_rhino.NominalTypeBuilder, getterType *pkg_rhino.FunctionType) {
	_ = classType
	_ = getterType
}

// Java method: DefaultCodingConvention.isPropertyTestFunction(com.google.javascript.rhino.Node call)
func (this_ *CodingConventions_DefaultCodingConvention) IsPropertyTestFunction(call pkg_rhino.Node) bool {
	_ = call
	var target pkg_rhino.Node = call.GetFirstChild()
	_ = target
	if target.IsGetProp() {
		var src pkg_rhino.Node = target.GetFirstChild()
		_ = src
		var prop string = target.GetString()
		_ = prop
		if (src.IsName() && (j2gort.CheckNotNullString(src.GetString()) == "Array")) && (j2gort.CheckNotNullString(prop) == "isArray") {
			return true
		}
	}
	return false
}

// Java method: DefaultCodingConvention.isPropertyRenameFunction(com.google.javascript.rhino.Node nameNode)
func (this_ *CodingConventions_DefaultCodingConvention) IsPropertyRenameFunction(nameNode pkg_rhino.Node) bool {
	_ = nameNode
	return nameNode.MatchesName_1_string("JSCompiler_renameProperty")
}

// Java method: DefaultCodingConvention.isPrototypeAlias(com.google.javascript.rhino.Node getProp)
func (this_ *CodingConventions_DefaultCodingConvention) IsPrototypeAlias(getProp pkg_rhino.Node) bool {
	_ = getProp
	return false
}

// Java method: DefaultCodingConvention.getObjectLiteralCast(com.google.javascript.rhino.Node callNode)
func (this_ *CodingConventions_DefaultCodingConvention) GetObjectLiteralCast(callNode pkg_rhino.Node) *CodingConvention_ObjectLiteralCast {
	_ = callNode
	return nil
}

// Java method: DefaultCodingConvention.getAssertionFunctions()
func (this_ *CodingConventions_DefaultCodingConvention) GetAssertionFunctions() any {
	return j2gort.Ext_com_google_common_collect_ImmutableSet_Of_2(CodingConvention_AssertionFunctionSpec_ForTruthy().SetClosurePrimitive(pkg_rhino.StaticGet_ClosurePrimitive_ASSERTS_TRUTHY()).Build(), CodingConvention_AssertionFunctionSpec_ForMatchesReturn().SetClosurePrimitive(pkg_rhino.StaticGet_ClosurePrimitive_ASSERTS_MATCHES_RETURN()).Build())
}

// Java method: DefaultCodingConvention.describeFunctionBind(com.google.javascript.rhino.Node n, boolean checkTypes)
func (this_ *CodingConventions_DefaultCodingConvention) DescribeFunctionBind(n pkg_rhino.Node, checkTypes bool) *CodingConvention_Bind {
	_ = n
	_ = checkTypes
	if !(n.IsCall()) {
		return nil
	}
	var callTarget pkg_rhino.Node = n.GetFirstChild()
	_ = callTarget
	if callTarget.IsQualifiedName() {
		if StaticGet_CodingConventions_DefaultCodingConvention_FUNCTION_PROTOTYPE_BIND_CALL().Matches(callTarget) {
			var fn pkg_rhino.Node = callTarget.GetNext()
			_ = fn
			if fn == 0 {
				return nil
			}
			var thisValue pkg_rhino.Node = CodingConventions_DefaultCodingConvention_SafeNext(fn)
			_ = thisValue
			var parameters pkg_rhino.Node = CodingConventions_DefaultCodingConvention_SafeNext(thisValue)
			_ = parameters
			return NewCodingConvention_Bind(fn, thisValue, parameters)
		}
	}
	if callTarget.IsGetProp() && (j2gort.CheckNotNullString(callTarget.GetString()) == "bind") {
		var maybeFn pkg_rhino.Node = callTarget.GetFirstChild()
		_ = maybeFn
		var maybeFnType pkg_rhino.JSType = maybeFn.GetJSType()
		_ = maybeFnType
		var fnType *pkg_rhino.FunctionType = nil
		_ = fnType
		if checkTypes && (maybeFnType != nil && !maybeFnType.IsJavaNil_()) {
			fnType = maybeFnType.RestrictByNotNullOrUndefined().ToMaybeFunctionType()
		}
		if (fnType != nil) || maybeFn.IsFunction() {
			var thisValue_2 pkg_rhino.Node = callTarget.GetNext()
			_ = thisValue_2
			var parameters_2 pkg_rhino.Node = CodingConventions_DefaultCodingConvention_SafeNext(thisValue_2)
			_ = parameters_2
			return NewCodingConvention_Bind(maybeFn, thisValue_2, parameters_2)
		}
	}
	return nil
}

// Java method: DefaultCodingConvention.describeCachingCall(com.google.javascript.rhino.Node node)
func (this_ *CodingConventions_DefaultCodingConvention) DescribeCachingCall(node pkg_rhino.Node) *CodingConvention_Cache {
	_ = node
	return nil
}

// Java static method: DefaultCodingConvention.safeNext(com.google.javascript.rhino.Node n)
func CodingConventions_DefaultCodingConvention_SafeNext(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	if n != 0 {
		return n.GetNext()
	}
	return pkg_rhino.Node(0)
}

func (this_ *CodingConventions_DefaultCodingConvention) BlockRenamingForProperty(p0_ string) bool {
	return IfaceDefault_CodingConvention_BlockRenamingForProperty(this_, p0_)
}
