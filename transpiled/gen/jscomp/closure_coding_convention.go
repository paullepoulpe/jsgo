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
// Port of com/google/javascript/jscomp/ClosureCodingConvention.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.ClosureCodingConvention
type ClosureCodingConvention struct {
	CodingConventions_ProxyBase
}

func (this_ *ClosureCodingConvention) BindSelf_(s any) {
	_ = s
	this_.CodingConventions_ProxyBase.BindSelf_(s)
}
func (this_ *ClosureCodingConvention) IsJavaNil_() bool { return this_ == nil }
func (*ClosureCodingConvention) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ClosureCodingConvention", "ClosureCodingConvention"
}
func (this_ *ClosureCodingConvention) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.CodingConventions_ProxyBase.IdentityHash_()
}
func (*ClosureCodingConvention) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ClosureCodingConvention", "com.google.javascript.jscomp.ClosureCodingConvention", func(v_ any) bool {
		return j2gort.InstanceOf[*ClosureCodingConvention](v_)
	})
}
func (*ClosureCodingConvention) Implements_jscomp_CodingConventions_Proxy() {}
func (*ClosureCodingConvention) Implements_jscomp_CodingConvention()        {}

func StaticGet_ClosureCodingConvention_SerialVersionUID() int64 { return int64(1) }

var static_ClosureCodingConvention_OBJECTLIT_EXPECTED *DiagnosticType

func StaticGet_ClosureCodingConvention_OBJECTLIT_EXPECTED() *DiagnosticType {
	clinit_ClosureCodingConvention()
	return static_ClosureCodingConvention_OBJECTLIT_EXPECTED
}
func StaticSet_ClosureCodingConvention_OBJECTLIT_EXPECTED(v *DiagnosticType) {
	clinit_ClosureCodingConvention()
	static_ClosureCodingConvention_OBJECTLIT_EXPECTED = v
}

var static_ClosureCodingConvention_GOOG_FORWARDDECLARE pkg_rhino.QualifiedName

func StaticGet_ClosureCodingConvention_GOOG_FORWARDDECLARE() pkg_rhino.QualifiedName {
	clinit_ClosureCodingConvention()
	return static_ClosureCodingConvention_GOOG_FORWARDDECLARE
}
func StaticSet_ClosureCodingConvention_GOOG_FORWARDDECLARE(v pkg_rhino.QualifiedName) {
	clinit_ClosureCodingConvention()
	static_ClosureCodingConvention_GOOG_FORWARDDECLARE = v
}

var static_ClosureCodingConvention_GOOG_ADDSINGLETONGETTER pkg_rhino.QualifiedName

func StaticGet_ClosureCodingConvention_GOOG_ADDSINGLETONGETTER() pkg_rhino.QualifiedName {
	clinit_ClosureCodingConvention()
	return static_ClosureCodingConvention_GOOG_ADDSINGLETONGETTER
}
func StaticSet_ClosureCodingConvention_GOOG_ADDSINGLETONGETTER(v pkg_rhino.QualifiedName) {
	clinit_ClosureCodingConvention()
	static_ClosureCodingConvention_GOOG_ADDSINGLETONGETTER = v
}

var static_ClosureCodingConvention_GOOG_ADDSINGLETONGETTER_MANGLED pkg_rhino.QualifiedName

func StaticGet_ClosureCodingConvention_GOOG_ADDSINGLETONGETTER_MANGLED() pkg_rhino.QualifiedName {
	clinit_ClosureCodingConvention()
	return static_ClosureCodingConvention_GOOG_ADDSINGLETONGETTER_MANGLED
}
func StaticSet_ClosureCodingConvention_GOOG_ADDSINGLETONGETTER_MANGLED(v pkg_rhino.QualifiedName) {
	clinit_ClosureCodingConvention()
	static_ClosureCodingConvention_GOOG_ADDSINGLETONGETTER_MANGLED = v
}

var static_ClosureCodingConvention_GOOG_REFLECT_OBJECTPROPERTY pkg_rhino.QualifiedName

func StaticGet_ClosureCodingConvention_GOOG_REFLECT_OBJECTPROPERTY() pkg_rhino.QualifiedName {
	clinit_ClosureCodingConvention()
	return static_ClosureCodingConvention_GOOG_REFLECT_OBJECTPROPERTY
}
func StaticSet_ClosureCodingConvention_GOOG_REFLECT_OBJECTPROPERTY(v pkg_rhino.QualifiedName) {
	clinit_ClosureCodingConvention()
	static_ClosureCodingConvention_GOOG_REFLECT_OBJECTPROPERTY = v
}

var static_ClosureCodingConvention_GOOG_REFLECT_OBJECTPROPERTY_MANGLED pkg_rhino.QualifiedName

func StaticGet_ClosureCodingConvention_GOOG_REFLECT_OBJECTPROPERTY_MANGLED() pkg_rhino.QualifiedName {
	clinit_ClosureCodingConvention()
	return static_ClosureCodingConvention_GOOG_REFLECT_OBJECTPROPERTY_MANGLED
}
func StaticSet_ClosureCodingConvention_GOOG_REFLECT_OBJECTPROPERTY_MANGLED(v pkg_rhino.QualifiedName) {
	clinit_ClosureCodingConvention()
	static_ClosureCodingConvention_GOOG_REFLECT_OBJECTPROPERTY_MANGLED = v
}

var static_ClosureCodingConvention_GOOG_REFLECT_OBJECT pkg_rhino.QualifiedName

func StaticGet_ClosureCodingConvention_GOOG_REFLECT_OBJECT() pkg_rhino.QualifiedName {
	clinit_ClosureCodingConvention()
	return static_ClosureCodingConvention_GOOG_REFLECT_OBJECT
}
func StaticSet_ClosureCodingConvention_GOOG_REFLECT_OBJECT(v pkg_rhino.QualifiedName) {
	clinit_ClosureCodingConvention()
	static_ClosureCodingConvention_GOOG_REFLECT_OBJECT = v
}

var static_ClosureCodingConvention_JSCOMP_REFLECTOBJECT pkg_rhino.QualifiedName

func StaticGet_ClosureCodingConvention_JSCOMP_REFLECTOBJECT() pkg_rhino.QualifiedName {
	clinit_ClosureCodingConvention()
	return static_ClosureCodingConvention_JSCOMP_REFLECTOBJECT
}
func StaticSet_ClosureCodingConvention_JSCOMP_REFLECTOBJECT(v pkg_rhino.QualifiedName) {
	clinit_ClosureCodingConvention()
	static_ClosureCodingConvention_JSCOMP_REFLECTOBJECT = v
}

var static_ClosureCodingConvention_GOOG_BIND pkg_rhino.QualifiedName

func StaticGet_ClosureCodingConvention_GOOG_BIND() pkg_rhino.QualifiedName {
	clinit_ClosureCodingConvention()
	return static_ClosureCodingConvention_GOOG_BIND
}
func StaticSet_ClosureCodingConvention_GOOG_BIND(v pkg_rhino.QualifiedName) {
	clinit_ClosureCodingConvention()
	static_ClosureCodingConvention_GOOG_BIND = v
}

var static_ClosureCodingConvention_GOOG_PARTIAL pkg_rhino.QualifiedName

func StaticGet_ClosureCodingConvention_GOOG_PARTIAL() pkg_rhino.QualifiedName {
	clinit_ClosureCodingConvention()
	return static_ClosureCodingConvention_GOOG_PARTIAL
}
func StaticSet_ClosureCodingConvention_GOOG_PARTIAL(v pkg_rhino.QualifiedName) {
	clinit_ClosureCodingConvention()
	static_ClosureCodingConvention_GOOG_PARTIAL = v
}

var static_ClosureCodingConvention_googCacheReflect pkg_rhino.Node

func StaticGet_ClosureCodingConvention_GoogCacheReflect() pkg_rhino.Node {
	clinit_ClosureCodingConvention()
	return static_ClosureCodingConvention_googCacheReflect
}
func StaticSet_ClosureCodingConvention_GoogCacheReflect(v pkg_rhino.Node) {
	clinit_ClosureCodingConvention()
	static_ClosureCodingConvention_googCacheReflect = v
}

var clinitDone_ClosureCodingConvention uint32
var clinitState_ClosureCodingConvention j2gort.ClinitState

func clinit_ClosureCodingConvention() {
	if atomic.LoadUint32(&clinitDone_ClosureCodingConvention) != 0 {
		return
	}
	clinitSlow_ClosureCodingConvention()
}

func clinitSlow_ClosureCodingConvention() {
	if !j2gort.EnterClinit(&clinitState_ClosureCodingConvention, &clinitDone_ClosureCodingConvention) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ClosureCodingConvention, &clinitDone_ClosureCodingConvention)
	static_ClosureCodingConvention_OBJECTLIT_EXPECTED = DiagnosticType_Warning("JSC_REFLECT_OBJECTLIT_EXPECTED", "Object literal expected as second argument")
	static_ClosureCodingConvention_GOOG_FORWARDDECLARE = pkg_rhino.QualifiedName_Of("goog.forwardDeclare")
	static_ClosureCodingConvention_GOOG_ADDSINGLETONGETTER = pkg_rhino.QualifiedName_Of("goog.addSingletonGetter")
	static_ClosureCodingConvention_GOOG_ADDSINGLETONGETTER_MANGLED = pkg_rhino.QualifiedName_Of("goog$addSingletonGetter")
	static_ClosureCodingConvention_GOOG_REFLECT_OBJECTPROPERTY = pkg_rhino.QualifiedName_Of("goog.reflect.objectProperty")
	static_ClosureCodingConvention_GOOG_REFLECT_OBJECTPROPERTY_MANGLED = pkg_rhino.QualifiedName_Of("goog$reflect$objectProperty")
	static_ClosureCodingConvention_GOOG_REFLECT_OBJECT = pkg_rhino.QualifiedName_Of("goog.reflect.object")
	static_ClosureCodingConvention_JSCOMP_REFLECTOBJECT = pkg_rhino.QualifiedName_Of("$jscomp.reflectObject")
	static_ClosureCodingConvention_GOOG_BIND = pkg_rhino.QualifiedName_Of("goog.bind")
	static_ClosureCodingConvention_GOOG_PARTIAL = pkg_rhino.QualifiedName_Of("goog.partial")
	static_ClosureCodingConvention_googCacheReflect = pkg_rhino.IR_Getprop_3(pkg_rhino.IR_Name("goog"), "reflect", []string{"cache"})
}

func (this_ *ClosureCodingConvention) initFields_ClosureCodingConvention() {
	_ = this_
}

// Java constructor: ClosureCodingConvention()
func NewClosureCodingConvention() *ClosureCodingConvention {
	clinit_ClosureCodingConvention()
	this_ := &ClosureCodingConvention{}
	this_.BindSelf_(this_)
	InitClosureCodingConvention(this_)
	return this_
}

func InitClosureCodingConvention(this_ *ClosureCodingConvention) {
	InitClosureCodingConvention_1(this_, CodingConventions_GetDefault())
}

// Java constructor: ClosureCodingConvention(com.google.javascript.jscomp.CodingConvention wrapped)
func NewClosureCodingConvention_1(wrapped CodingConvention) *ClosureCodingConvention {
	clinit_ClosureCodingConvention()
	this_ := &ClosureCodingConvention{}
	this_.BindSelf_(this_)
	InitClosureCodingConvention_1(this_, wrapped)
	return this_
}

func InitClosureCodingConvention_1(this_ *ClosureCodingConvention, wrapped CodingConvention) {
	_ = wrapped
	InitCodingConventions_Proxy(&this_.CodingConventions_ProxyBase, wrapped)
	this_.initFields_ClosureCodingConvention()
}

// Java method: ClosureCodingConvention.applySubclassRelationship(com.google.javascript.rhino.NominalTypeBuilder parent, com.google.javascript.rhino.NominalTypeBuilder child, com.google.javascript.jscomp.CodingConvention.SubclassType type)
func (this_ *ClosureCodingConvention) ApplySubclassRelationship(parent *pkg_rhino.NominalTypeBuilder, child *pkg_rhino.NominalTypeBuilder, type_ CodingConvention_SubclassType) {
	_ = parent
	_ = child
	_ = type_
	this_.CodingConventions_ProxyBase.ApplySubclassRelationship(parent, child, type_)
	if type_ == StaticGet_CodingConvention_SubclassType_INHERITS() {
		var childCtor *pkg_rhino.FunctionType = child.Constructor()
		_ = childCtor
		child.DeclareConstructorProperty("superClass_", pkg_rhino.Upcast_JSType(parent.PrototypeOrInstance()), childCtor.GetSource())
		var qmarkCtor *pkg_rhino.FunctionType = childCtor.ForgetParameterAndReturnTypes()
		_ = qmarkCtor
		child.DeclarePrototypeProperty("constructor", pkg_rhino.ToIface_JSType(qmarkCtor), childCtor.GetSource())
	}
}

// Java method: ClosureCodingConvention.getClassesDefinedByCall(com.google.javascript.rhino.Node callNode)
func (this_ *ClosureCodingConvention) GetClassesDefinedByCall(callNode pkg_rhino.Node) *CodingConvention_SubclassRelationship {
	_ = callNode
	var relationship *CodingConvention_SubclassRelationship = this_.CodingConventions_ProxyBase.GetClassesDefinedByCall(callNode)
	_ = relationship
	if relationship != nil {
		return relationship
	}
	var callName pkg_rhino.Node = callNode.GetFirstChild()
	_ = callName
	var type_ CodingConvention_SubclassType = ClosureCodingConvention_TypeofClassDefiningName(callName)
	_ = type_
	if type_ != 0 {
		if callNode.GetChildCount() < int32(3) {
			return nil
		}
		var subclass pkg_rhino.Node = callName.GetNext()
		_ = subclass
		var superclass pkg_rhino.Node = subclass.GetNext()
		_ = superclass
		if type_ == StaticGet_CodingConvention_SubclassType_MIXIN() {
			if ClosureCodingConvention_EndsWithPrototype(superclass) {
				superclass = superclass.GetFirstChild()
			}
			if ClosureCodingConvention_EndsWithPrototype(subclass) {
				subclass = subclass.GetFirstChild()
			}
		}
		if ((subclass != 0) && subclass.IsUnscopedQualifiedName()) && superclass.IsUnscopedQualifiedName() {
			return NewCodingConvention_SubclassRelationship(type_, subclass, superclass)
		}
	}
	return nil
}

// Java static method: ClosureCodingConvention.typeofClassDefiningName(com.google.javascript.rhino.Node callName)
func ClosureCodingConvention_TypeofClassDefiningName(callName pkg_rhino.Node) CodingConvention_SubclassType {
	clinit_ClosureCodingConvention()
	_ = callName
	var methodName string = j2gort.NullString
	_ = methodName
	if callName.IsGetProp() {
		methodName = callName.GetString()
	} else if callName.IsName() {
		var name string = callName.GetString()
		_ = name
		var dollarIndex int32 = j2gort.JStrLastIndexOfCp(name, int32(uint16(36)))
		_ = dollarIndex
		if dollarIndex != int32(-1) {
			methodName = j2gort.JStrSubFrom(name, (dollarIndex + int32(1)))
		}
	}
	if methodName != j2gort.NullString {
		if j2gort.CheckNotNullString(methodName) == "inherits" {
			return StaticGet_CodingConvention_SubclassType_INHERITS()
		} else if j2gort.CheckNotNullString(methodName) == "mixin" {
			return StaticGet_CodingConvention_SubclassType_MIXIN()
		}
	}
	return CodingConvention_SubclassType(0)
}

// Java method: ClosureCodingConvention.isSuperClassReference(java.lang.String propertyName)
func (this_ *ClosureCodingConvention) IsSuperClassReference(propertyName string) bool {
	_ = propertyName
	return (("superClass_" == propertyName) || this_.CodingConventions_ProxyBase.IsSuperClassReference(propertyName))
}

// Java static method: ClosureCodingConvention.endsWithPrototype(com.google.javascript.rhino.Node qualifiedName)
func ClosureCodingConvention_EndsWithPrototype(qualifiedName pkg_rhino.Node) bool {
	clinit_ClosureCodingConvention()
	_ = qualifiedName
	return (qualifiedName.IsGetProp() && (j2gort.CheckNotNullString(qualifiedName.GetString()) == "prototype"))
}

// Java method: ClosureCodingConvention.extractIsModuleFile(com.google.javascript.rhino.Node node, com.google.javascript.rhino.Node parent)
func (this_ *ClosureCodingConvention) ExtractIsModuleFile(node pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = node
	_ = parent
	var namespace string = ClosureCodingConvention_ExtractClassNameIfGoog(node, parent, "goog.module")
	_ = namespace
	return (namespace != j2gort.NullString)
}

// Java method: ClosureCodingConvention.extractClassNameIfProvide(com.google.javascript.rhino.Node node, com.google.javascript.rhino.Node parent)
func (this_ *ClosureCodingConvention) ExtractClassNameIfProvide(node pkg_rhino.Node, parent pkg_rhino.Node) string {
	_ = node
	_ = parent
	var namespace string = ClosureCodingConvention_ExtractClassNameIfGoog(node, parent, "goog.provide")
	_ = namespace
	if namespace == j2gort.NullString {
		namespace = ClosureCodingConvention_ExtractClassNameIfGoog(node, parent, "goog.module")
	}
	return namespace
}

// Java method: ClosureCodingConvention.extractClassNameIfRequire(com.google.javascript.rhino.Node node, com.google.javascript.rhino.Node parent)
func (this_ *ClosureCodingConvention) ExtractClassNameIfRequire(node pkg_rhino.Node, parent pkg_rhino.Node) string {
	_ = node
	_ = parent
	return ClosureCodingConvention_ExtractClassNameIfGoog(node, parent, "goog.require")
}

// Java static method: ClosureCodingConvention.extractClassNameIfGoog(com.google.javascript.rhino.Node node, com.google.javascript.rhino.Node parent, java.lang.String functionName)
func ClosureCodingConvention_ExtractClassNameIfGoog(node pkg_rhino.Node, parent pkg_rhino.Node, functionName string) string {
	clinit_ClosureCodingConvention()
	_ = node
	_ = parent
	_ = functionName
	var className string = j2gort.NullString
	_ = className
	if NodeUtil_IsExprCall(parent) {
		var callee pkg_rhino.Node = node.GetFirstChild()
		_ = callee
		if ((callee != 0) && callee.IsGetProp()) && callee.MatchesQualifiedName_1_string(functionName) {
			var target pkg_rhino.Node = callee.GetNext()
			_ = target
			if (target != 0) && target.IsStringLit() {
				className = target.GetString()
			}
		}
	}
	return className
}

// Java method: ClosureCodingConvention.getExportPropertyFunction()
func (this_ *ClosureCodingConvention) GetExportPropertyFunction() string {
	return "goog.exportProperty"
}

// Java method: ClosureCodingConvention.getExportSymbolFunction()
func (this_ *ClosureCodingConvention) GetExportSymbolFunction() string {
	return "goog.exportSymbol"
}

// Java method: ClosureCodingConvention.identifyTypeDeclarationCall(com.google.javascript.rhino.Node n)
func (this_ *ClosureCodingConvention) IdentifyTypeDeclarationCall(n pkg_rhino.Node) j2gort.Ext_java_util_List {
	_ = n
	var callName pkg_rhino.Node = n.GetFirstChild()
	_ = callName
	if StaticGet_ClosureCodingConvention_GOOG_FORWARDDECLARE().Matches(callName) && n.HasTwoChildren() {
		var typeDeclaration pkg_rhino.Node = n.GetSecondChild()
		_ = typeDeclaration
		if typeDeclaration.IsStringLit() {
			return j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.ImmListOfStr1(typeDeclaration.GetString()))
		}
	}
	return this_.CodingConventions_ProxyBase.IdentifyTypeDeclarationCall(n)
}

// Java method: ClosureCodingConvention.getAbstractMethodName()
func (this_ *ClosureCodingConvention) GetAbstractMethodName() string {
	return "goog.abstractMethod"
}

// Java method: ClosureCodingConvention.getSingletonGetterClassName(com.google.javascript.rhino.Node callNode)
func (this_ *ClosureCodingConvention) GetSingletonGetterClassName(callNode pkg_rhino.Node) string {
	_ = callNode
	var callArg pkg_rhino.Node = callNode.GetFirstChild()
	_ = callArg
	if callNode.HasTwoChildren() && (StaticGet_ClosureCodingConvention_GOOG_ADDSINGLETONGETTER().Matches(callArg) || StaticGet_ClosureCodingConvention_GOOG_ADDSINGLETONGETTER_MANGLED().Matches(callArg)) {
		return callArg.GetNext().GetQualifiedName()
	}
	return this_.CodingConventions_ProxyBase.GetSingletonGetterClassName(callNode)
}

// Java method: ClosureCodingConvention.applySingletonGetter(com.google.javascript.rhino.NominalTypeBuilder classType, com.google.javascript.rhino.jstype.FunctionType getterType)
func (this_ *ClosureCodingConvention) ApplySingletonGetter(classType *pkg_rhino.NominalTypeBuilder, getterType *pkg_rhino.FunctionType) {
	_ = classType
	_ = getterType
	var defSite pkg_rhino.Node = classType.Constructor().GetSource()
	_ = defSite
	classType.DeclareConstructorProperty("getInstance", pkg_rhino.ToIface_JSType(getterType), defSite)
	classType.DeclareConstructorProperty("instance_", pkg_rhino.Upcast_JSType(classType.Instance()), defSite)
}

// Java method: ClosureCodingConvention.isPropertyTestFunction(com.google.javascript.rhino.Node call)
func (this_ *ClosureCodingConvention) IsPropertyTestFunction(call pkg_rhino.Node) bool {
	_ = call
	if !(call.IsCall()) {
		j2gort.CheckArgument(false)
	}
	var target pkg_rhino.Node = call.GetFirstChild()
	_ = target
	if target.IsGetProp() {
		var src pkg_rhino.Node = target.GetFirstChild()
		_ = src
		var prop string = target.GetString()
		_ = prop
		if (src.IsName() && (j2gort.CheckNotNullString(src.GetString()) == "goog")) && ((j2gort.CheckNotNullString(prop) == "isArrayLike") || (j2gort.CheckNotNullString(prop) == "isObject")) {
			return true
		}
	}
	return this_.CodingConventions_ProxyBase.IsPropertyTestFunction(call)
}

// Java method: ClosureCodingConvention.isPropertyRenameFunction(com.google.javascript.rhino.Node nameNode)
func (this_ *ClosureCodingConvention) IsPropertyRenameFunction(nameNode pkg_rhino.Node) bool {
	_ = nameNode
	if this_.CodingConventions_ProxyBase.IsPropertyRenameFunction(nameNode) {
		return true
	}
	return (StaticGet_ClosureCodingConvention_GOOG_REFLECT_OBJECTPROPERTY().Matches(nameNode) || StaticGet_ClosureCodingConvention_GOOG_REFLECT_OBJECTPROPERTY_MANGLED().Matches(nameNode))
}

// Java method: ClosureCodingConvention.isFunctionCallThatAlwaysThrows(com.google.javascript.rhino.Node n)
func (this_ *ClosureCodingConvention) IsFunctionCallThatAlwaysThrows(n pkg_rhino.Node) bool {
	_ = n
	return (this_.CodingConventions_ProxyBase.IsFunctionCallThatAlwaysThrows(n) || CodingConventions_DefaultIsFunctionCallThatAlwaysThrows(n, "goog.asserts.fail"))
}

// Java method: ClosureCodingConvention.getObjectLiteralCast(com.google.javascript.rhino.Node callNode)
func (this_ *ClosureCodingConvention) GetObjectLiteralCast(callNode pkg_rhino.Node) *CodingConvention_ObjectLiteralCast {
	_ = callNode
	if !(callNode.IsCall()) {
		j2gort.CheckArgumentFmt(false, "Expected call node but found %s", j2gort.BoxNode(callNode))
	}
	var proxyCast *CodingConvention_ObjectLiteralCast = this_.CodingConventions_ProxyBase.GetObjectLiteralCast(callNode)
	_ = proxyCast
	if proxyCast != nil {
		return proxyCast
	}
	var callName pkg_rhino.Node = callNode.GetFirstChild()
	_ = callName
	if !(StaticGet_ClosureCodingConvention_GOOG_REFLECT_OBJECT().Matches(callName) || StaticGet_ClosureCodingConvention_JSCOMP_REFLECTOBJECT().Matches(callName)) || !(callNode.HasXChildren(int32(3))) {
		return nil
	}
	var typeNode pkg_rhino.Node = callName.GetNext()
	_ = typeNode
	if !(typeNode.IsQualifiedName()) {
		return nil
	}
	var objectNode pkg_rhino.Node = typeNode.GetNext()
	_ = objectNode
	if !(objectNode.IsObjectLit()) {
		return NewCodingConvention_ObjectLiteralCast(j2gort.NullString, pkg_rhino.Node(0), StaticGet_ClosureCodingConvention_OBJECTLIT_EXPECTED())
	}
	return NewCodingConvention_ObjectLiteralCast(typeNode.GetQualifiedName(), typeNode.GetNext(), nil)
}

// Java method: ClosureCodingConvention.getAssertionFunctions()
func (this_ *ClosureCodingConvention) GetAssertionFunctions() any {
	return j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet_Builder](j2gort.Ext_com_google_common_collect_ImmutableSet_S_Builder().AddAll_1_Ext_java_lang_Iterable(j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.CodingConventions_ProxyBase.GetAssertionFunctions()))).Add_1_slice_any([]any{CodingConvention_AssertionFunctionSpec_ForTruthy().SetFunctionName("goog.asserts.assert").GetCodingConvention_AssertionFunctionSpec_BuilderBase().Build(), ClosureCodingConvention_CreateGoogAssertOnReturn("Array"), ClosureCodingConvention_CreateGoogAssertOnReturn("Boolean"), ClosureCodingConvention_CreateGoogAssertOnReturn("Element"), ClosureCodingConvention_CreateGoogAssertOnReturn("Function"), ClosureCodingConvention_CreateGoogAssertOnReturn("Instanceof"), ClosureCodingConvention_CreateGoogAssertOnReturn("Number"), ClosureCodingConvention_CreateGoogAssertOnReturn("Object"), ClosureCodingConvention_CreateGoogAssertOnReturn("String")})).Build())
}

// Java static method: ClosureCodingConvention.createGoogAssertOnReturn(java.lang.String assertedTypeName)
func ClosureCodingConvention_CreateGoogAssertOnReturn(assertedTypeName string) CodingConvention_AssertionFunctionSpec {
	clinit_ClosureCodingConvention()
	_ = assertedTypeName
	return CodingConvention_AssertionFunctionSpec_ForMatchesReturn().SetFunctionName(("goog.asserts.assert" + j2gort.StrValString(assertedTypeName))).GetCodingConvention_AssertionFunctionSpec_BuilderBase(

	// Java method: ClosureCodingConvention.describeFunctionBind(com.google.javascript.rhino.Node n, boolean checkTypes)
	).Build()
}

func (this_ *ClosureCodingConvention) DescribeFunctionBind(n pkg_rhino.Node, checkTypes bool) *CodingConvention_Bind {
	_ = n
	_ = checkTypes
	if !(n.IsCall()) {
		return nil
	}
	var callTarget pkg_rhino.Node = n.GetFirstChild()
	_ = callTarget
	if callTarget.IsQualifiedName() {
		if StaticGet_ClosureCodingConvention_GOOG_BIND().Matches(callTarget) || callTarget.MatchesName_1_string("goog$bind") {
			var fn pkg_rhino.Node = callTarget.GetNext()
			_ = fn
			if fn == 0 {
				return nil
			}
			var thisValue pkg_rhino.Node = ClosureCodingConvention_SafeNext(fn)
			_ = thisValue
			var parameters pkg_rhino.Node = ClosureCodingConvention_SafeNext(thisValue)
			_ = parameters
			return NewCodingConvention_Bind(fn, thisValue, parameters)
		}
		if StaticGet_ClosureCodingConvention_GOOG_PARTIAL().Matches(callTarget) || callTarget.MatchesName_1_string("goog$partial") {
			var fn_2 pkg_rhino.Node = callTarget.GetNext()
			_ = fn_2
			if fn_2 == 0 {
				return nil
			}
			var thisValue_2 pkg_rhino.Node = pkg_rhino.Node(0)
			_ = thisValue_2
			var parameters_2 pkg_rhino.Node = ClosureCodingConvention_SafeNext(fn_2)
			_ = parameters_2
			return NewCodingConvention_Bind(fn_2, thisValue_2, parameters_2)
		}
	}
	return this_.CodingConventions_ProxyBase.DescribeFunctionBind(n, checkTypes)
}

// Java method: ClosureCodingConvention.describeCachingCall(com.google.javascript.rhino.Node node)
func (this_ *ClosureCodingConvention) DescribeCachingCall(node pkg_rhino.Node) *CodingConvention_Cache {
	_ = node
	if !(node.IsCall()) {
		return nil
	}
	var callTarget pkg_rhino.Node = node.GetFirstChild()
	_ = callTarget
	if this_.priv_ClosureCodingConvention_matchesCacheMethodName_1164(callTarget) {
		var paramCount int32 = (node.GetChildCount() - int32(1))
		_ = paramCount
		if (int32(3) <= paramCount) && (paramCount <= int32(4)) {
			var cacheObj pkg_rhino.Node = callTarget.GetNext()
			_ = cacheObj
			var keyNode pkg_rhino.Node = cacheObj.GetNext()
			_ = keyNode
			var valueFn pkg_rhino.Node = keyNode.GetNext()
			_ = valueFn
			var keyFn pkg_rhino.Node = valueFn.GetNext()
			_ = keyFn
			return NewCodingConvention_Cache(cacheObj, keyNode, valueFn, keyFn)
		}
	}
	return this_.CodingConventions_ProxyBase.DescribeCachingCall(node)
}

// Java method: ClosureCodingConvention.matchesCacheMethodName(com.google.javascript.rhino.Node target)
func (this_ *ClosureCodingConvention) priv_ClosureCodingConvention_matchesCacheMethodName_1164(target pkg_rhino.Node) bool {
	_ = target
	if target.IsGetProp() {
		return target.MatchesQualifiedName_1_Node(StaticGet_ClosureCodingConvention_GoogCacheReflect())
	} else if target.IsName() {
		return (j2gort.CheckNotNullString(target.GetString()) == "goog$reflect$cache")
	}
	return false
}

// Java static method: ClosureCodingConvention.safeNext(com.google.javascript.rhino.Node n)
func ClosureCodingConvention_SafeNext(n pkg_rhino.Node) pkg_rhino.Node {
	clinit_ClosureCodingConvention()
	_ = n
	if n != 0 {
		return n.GetNext()
	}
	return pkg_rhino.Node(0)
}
