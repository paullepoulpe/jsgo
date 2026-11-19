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
// Port of com/google/javascript/rhino/jstype/NamedType.java (generated mechanically by j2go).

package rhino

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.rhino.jstype.NamedType
type NamedType struct {
	ProxyObjectTypeBase
	reference             string                                      `gson:"reference"`
	sourceName            string                                      `gson:"sourceName"`
	resolutionScope       StaticTypedScope                            `gson:"resolutionScope"`
	validator             j2gort.Ext_com_google_common_base_Predicate `gson:"-"`
	propertyContinuations []*NamedType_PropertyContinuation           `gson:"-"`
	templateTypes         []JSType                                    `gson:"templateTypes"`
	namedResolveResult    JSType                                      `gson:"namedResolveResult"`
	lineno                int32                                       `gson:"lineno"`
	charno                int32                                       `gson:"charno"`
	resolutionKind        NamedType_ResolutionKind                    `gson:"resolutionKind"`
	restrictByNull        bool                                        `gson:"restrictByNull"`
}

func (this_ *NamedType) BindSelf_(s any) {
	_ = s
	this_.ProxyObjectTypeBase.BindSelf_(s)
}
func (this_ *NamedType) IsJavaNil_() bool { return this_ == nil }
func (*NamedType) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.NamedType", "NamedType"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.NamedType", "com.google.javascript.rhino.jstype.NamedType", func(v_ any) bool {
		return j2gort.InstanceOf[*NamedType](v_)
	})
}
func (*NamedType) Implements_rhino_ProxyObjectType() {}
func (*NamedType) Implements_rhino_ObjectType()      {}
func (*NamedType) Implements_rhino_JSType()          {}

var static_NamedType_TYPE_CLASS JSTypeClass

func StaticGet_NamedType_TYPE_CLASS() JSTypeClass {
	clinit_NamedType()
	return static_NamedType_TYPE_CLASS
}
func StaticSet_NamedType_TYPE_CLASS(v JSTypeClass) {
	clinit_NamedType()
	static_NamedType_TYPE_CLASS = v
}

var clinitDone_NamedType uint32
var clinitState_NamedType j2gort.ClinitState

func clinit_NamedType() {
	if atomic.LoadUint32(&clinitDone_NamedType) != 0 {
		return
	}
	clinitSlow_NamedType()
}

func clinitSlow_NamedType() {
	if !j2gort.EnterClinit(&clinitState_NamedType, &clinitDone_NamedType) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_NamedType, &clinitDone_NamedType)
	static_NamedType_TYPE_CLASS = StaticGet_JSTypeClass_NAMED()
}

func (this_ *NamedType) initFields_NamedType() {
	_ = this_
	this_.reference = j2gort.NullString
	this_.sourceName = j2gort.NullString

}

// Java constructor: NamedType(com.google.javascript.rhino.jstype.NamedType.Builder builder)
func NewNamedType(builder *NamedType_Builder) *NamedType {
	this_ := &NamedType{}
	this_.BindSelf_(this_)
	InitNamedType(this_, builder)
	return this_
}

func InitNamedType(this_ *NamedType, builder *NamedType_Builder) {
	_ = builder
	InitProxyObjectType_4(&this_.ProxyObjectTypeBase, builder.registry, builder.referencedType, nil, int32(1024))
	this_.initFields_NamedType()
	j2gort.CheckNotNullString(builder.referenceName)
	j2gort.CheckNotNullEnum(builder.resolutionKind)
	j2gort.CheckNotNull(j2gort.ImmListFromTypedSlice(builder.templateTypes))
	{
		var iifeCond_1 bool
		{

			re_1 := builder.resolutionKind
			ae_2 := StaticGet_NamedType_ResolutionKind_TYPEOF()
			if re_1 == 0 {
				j2gort.ThrowNPE()
			}
			iifeCond_1 = re_1 == ae_2
		}
		if iifeCond_1 {
			if !(javalang.StringStartsWith(j2gort.CheckNotNullString(builder.referenceName), "typeof ")) {
				j2gort.CheckState(false)
			}
		}
	}
	this_.restrictByNull = builder.restrictByNull
	this_.resolutionScope = builder.scope
	this_.reference = builder.referenceName
	this_.sourceName = builder.sourceName
	this_.lineno = builder.lineno
	this_.charno = builder.charno
	this_.templateTypes = builder.templateTypes
	this_.resolutionKind = builder.resolutionKind
	this_.GetJSTypeBase().registry.GetResolver().ResolveIfClosed(ToIface_JSType(this_), StaticGet_NamedType_TYPE_CLASS())
}

// Java static method: NamedType.nominalHashCode(com.google.javascript.rhino.jstype.ObjectType type)
func NamedType_NominalHashCode(type_ ObjectType) int32 {
	_ = type_
	if !(type_.GetObjectTypeBase().HasReferenceName()) {
		j2gort.CheckState(false)
	}
	var name string = j2gort.CheckNotNullString(type_.GetReferenceName())
	_ = name
	return javalang.StringHashCode(j2gort.CheckNotNullString(name))
}

// Java method: NamedType.getTypeClass()
func (this_ *NamedType) GetTypeClass() JSTypeClass {
	return StaticGet_NamedType_TYPE_CLASS()
}

// Java method: NamedType.getBangType()
func (this_ *NamedType) GetBangType() JSType {
	if this_.restrictByNull {
		return this_
	} else if this_.IsResolved() {
		if this_.IsNoResolvedType() || this_.IsUnknownType() {
			return this_
		} else {
			return this_.GetReferencedType().RestrictByNotNullOrUndefined()
		}
	}
	return ToIface_JSType(this_.ToBuilder().priv_Builder_setRestrictByNull_9312(true).Build())
}

// Java method: NamedType.getTemplateTypes()
func (this_ *NamedType) GetTemplateTypes() j2gort.Ext_com_google_common_collect_ImmutableList {
	return j2gort.ImmListFromTypedSlice(this_.templateTypes)
}

// Java method: NamedType.defineProperty(java.lang.String propertyName, com.google.javascript.rhino.jstype.JSType type, boolean inferred, com.google.javascript.rhino.Node propertyNode)
func (this_ *NamedType) DefineProperty_4_string_JSType_bool_Node(propertyName string, type_ JSType, inferred bool, propertyNode Node) bool {
	_ = propertyName
	_ = type_
	_ = inferred
	_ = propertyNode
	if !(this_.IsResolved()) {
		if this_.propertyContinuations == nil {
			this_.propertyContinuations = []*NamedType_PropertyContinuation{}
		}
		this_.propertyContinuations = append(this_.propertyContinuations, NewNamedType_PropertyContinuation(NewProperty_StringKey(propertyName), type_, inferred, propertyNode))
		return true
	} else {
		return this_.ProxyObjectTypeBase.DefineProperty_4_string_JSType_bool_Node(propertyName, type_, inferred, propertyNode)
	}
}

// Java method: NamedType.defineProperty(com.google.javascript.rhino.jstype.Property.Key propertyName, com.google.javascript.rhino.jstype.JSType type, boolean inferred, com.google.javascript.rhino.Node propertyNode)
func (this_ *NamedType) DefineProperty_4_Property_Key_JSType_bool_Node(propertyName Property_Key, type_ JSType, inferred bool, propertyNode Node) bool {
	_ = propertyName
	_ = type_
	_ = inferred
	_ = propertyNode
	if !(this_.IsResolved()) {
		if this_.propertyContinuations == nil {
			this_.propertyContinuations = []*NamedType_PropertyContinuation{}
		}
		this_.propertyContinuations = append(this_.propertyContinuations, NewNamedType_PropertyContinuation(propertyName, type_, inferred, propertyNode))
		return true
	} else {
		return this_.ProxyObjectTypeBase.DefineProperty_4_Property_Key_JSType_bool_Node(propertyName, type_, inferred, propertyNode)
	}
}

// Java method: NamedType.finishPropertyContinuations()
func (this_ *NamedType) priv_NamedType_finishPropertyContinuations_9293() {
	var referencedObjType ObjectType = this_.GetReferencedObjTypeInternal()
	_ = referencedObjType
	if ((referencedObjType != nil && !referencedObjType.IsJavaNil_()) && !(referencedObjType.GetObjectTypeBase().IsUnknownType())) && (this_.propertyContinuations != nil) {
		for _, c := range this_.propertyContinuations {
			_ = c
			c.Commit(this_)
		}

	}
	this_.propertyContinuations = nil
}

// Java method: NamedType.getReferencedType()
func (this_ *NamedType) GetReferencedType() JSType {
	return this_.GetReferencedTypeInternal()
}

// Java method: NamedType.getReferenceName()
func (this_ *NamedType) GetReferenceName() string {
	return this_.reference
}

// Java method: NamedType.appendTo(com.google.javascript.rhino.jstype.TypeStringBuilder sb)
func (this_ *NamedType) AppendTo(sb *TypeStringBuilder) {
	_ = sb
	var type_ JSType = this_.GetReferencedType()
	_ = type_
	if !(this_.IsResolved()) || type_.GetJSTypeBase().IsNoResolvedType() {
		sb.Append_1_string(this_.GetReferenceName())
	} else {
		sb.Append_1_JSType(type_)
	}
}

// Java method: NamedType.toMaybeNamedType()
func (this_ *NamedType) ToMaybeNamedType() *NamedType {
	return this_
}

// Java method: NamedType.isNominalType()
func (this_ *NamedType) IsNominalType() bool {
	if this_.IsResolved() {
		return this_.ProxyObjectTypeBase.IsNominalType()
	} else {
		return true
	}
}

// Java method: NamedType.recursionUnsafeHashCode()
func (this_ *NamedType) RecursionUnsafeHashCode() int32 {
	if this_.IsSuccessfullyResolved() {
		return this_.ProxyObjectTypeBase.RecursionUnsafeHashCode()
	} else {
		return NamedType_NominalHashCode(this_)
	}
}

// Java method: NamedType.getCachedResolveResult()
func (this_ *NamedType) GetCachedResolveResult() JSType {
	if !(this_.IsResolved()) {
		return nil
	}
	if this_.namedResolveResult != nil && !this_.namedResolveResult.IsJavaNil_() {
		return this_.namedResolveResult
	} else {
		return this_
	}
}

// Java method: NamedType.setCachedResolveResult(com.google.javascript.rhino.jstype.JSType result)
func (this_ *NamedType) SetCachedResolveResult(result JSType) {
	_ = result
	j2gort.CheckNotNullIface(result)
	if result == this_ || ((result == nil || result.IsJavaNil_()) && (this_ == nil || this_.IsJavaNil_())) {
		this_.namedResolveResult = nil
	} else {
		this_.namedResolveResult = result
	}
	this_.MarkResolved()
}

// Java method: NamedType.clearCachedResolveResult()
func (this_ *NamedType) ClearCachedResolveResult() {
	this_.namedResolveResult = nil
	this_.ProxyObjectTypeBase.ClearCachedResolveResult()
}

// Java method: NamedType.resolveInternal(com.google.javascript.rhino.ErrorReporter reporter)
func (this_ *NamedType) ResolveInternal(reporter ErrorReporter) JSType {
	_ = reporter
	var resolvedTypeArgs j2gort.Ext_com_google_common_collect_ImmutableList = nil
	_ = resolvedTypeArgs
	if len(this_.templateTypes) == 0 {
		resolvedTypeArgs = j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()
	} else {
		var argsBuilder *j2gort.Ext_com_google_common_collect_ImmutableList_Builder = j2gort.InitImmListBuilderCap(&j2gort.Ext_com_google_common_collect_ImmutableList_Builder{}, int32(len(this_.templateTypes)))
		_ = argsBuilder
		{
			var i int32 = int32(0)
			_ = i
			var n int32 = int32(len(this_.templateTypes))
			_ = n
			for ; i < n; i = i + 1 {
				argsBuilder.AddVoid_(this_.templateTypes[i].GetJSTypeBase().Resolve(reporter))
			}
		}
		resolvedTypeArgs = argsBuilder.BuildImmList_()
	}
	{
		var iifeCond_1 bool
		{

			re_5 := this_.resolutionKind
			ae_6 := StaticGet_NamedType_ResolutionKind_NONE()
			if re_5 == 0 {
				j2gort.ThrowNPE()
			}
			iifeCond_1 = re_5 == ae_6
		}
		if iifeCond_1 {
			return this_.ProxyObjectTypeBase.ResolveInternal(reporter)
		}
	}
	if !(this_.GetReferencedType().GetJSTypeBase().IsUnknownType()) {
		j2gort.CheckStateStr(false, "NamedTypes given a referenced type pre-resolution should have ResolutionKind.NONE")
	}
	if this_.resolutionScope == nil || this_.resolutionScope.IsJavaNil_() {
		return this_
	}
	var unused bool = (this_.priv_NamedType_resolveTypeof_9296(reporter) || this_.priv_NamedType_resolveViaRegistry_9295(reporter))
	_ = unused
	this_.ProxyObjectTypeBase.ResolveInternal(reporter)
	if this_.DetectInheritanceCycle() {
		this_.priv_NamedType_handleTypeCycle_9298(reporter)
	}
	this_.priv_NamedType_finishPropertyContinuations_9293()
	var result JSType = this_.GetReferencedType()
	_ = result
	if this_.IsSuccessfullyResolved() {
		this_.resolutionScope = nil
		var resultAsObject ObjectType = result.GetJSTypeBase().ToMaybeObjectType()
		_ = resultAsObject
		if resultAsObject == nil || resultAsObject.IsJavaNil_() {
			return result
		}
		if resolvedTypeArgs.IsEmpty() || !(resultAsObject.GetObjectTypeBase().IsRawTypeOfTemplatizedType()) {
			return result
		}
		var numKeys int32 = result.GetTemplateParamCount()
		_ = numKeys
		if numKeys < resolvedTypeArgs.Size() {
			resolvedTypeArgs = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableList](resolvedTypeArgs.SubList(int32(0), numKeys))
		}
		result = ToIface_JSType(this_.GetJSTypeBase().registry.CreateTemplatizedType_2_ObjectType_Ext_com_google_common_collect_ImmutableList(resultAsObject, resolvedTypeArgs))
		this_.SetReferencedType(result.GetJSTypeBase().Resolve(reporter))
	}
	return result
}

// Java method: NamedType.resolveViaRegistry(com.google.javascript.rhino.ErrorReporter reporter)
func (this_ *NamedType) priv_NamedType_resolveViaRegistry_9295(reporter ErrorReporter) bool {
	_ = reporter
	var type_ JSType = this_.GetJSTypeBase().registry.GetType_2(Upcast_StaticScope(this_.resolutionScope), this_.reference)
	_ = type_
	if type_ == nil || type_.IsJavaNil_() {
		this_.priv_NamedType_handleUnresolvedType_9301(reporter)
		return false
	}
	this_.priv_NamedType_setReferencedAndResolvedType_9297(type_, reporter)
	return true
}

// Java method: NamedType.resolveTypeof(com.google.javascript.rhino.ErrorReporter reporter)
func (this_ *NamedType) priv_NamedType_resolveTypeof_9296(reporter ErrorReporter) bool {
	_ = reporter
	{
		var iifeCond_1 bool
		{

			re_7 := this_.resolutionKind
			ae_8 := StaticGet_NamedType_ResolutionKind_TYPEOF()
			if re_7 == 0 {
				j2gort.ThrowNPE()
			}
			iifeCond_1 = re_7 == ae_8
		}
		if !(iifeCond_1) {
			return false
		}
	}
	var scopeName string = j2gort.JStrSubFrom(this_.reference, j2gort.JStrLen("typeof "))
	_ = scopeName
	var type_ JSType = this_.resolutionScope.LookupQualifiedName_1_string(scopeName)
	_ = type_
	if (type_ == nil || type_.IsJavaNil_()) || type_.GetJSTypeBase().IsUnknownType() {
		if this_.GetJSTypeBase().registry.IsForwardDeclaredType(scopeName) {
			this_.SetReferencedType(NewNoResolvedType(this_.GetJSTypeBase().registry, this_.GetReferenceName(), this_.GetTemplateTypes()))
			if this_.validator != nil {
				var unused bool = this_.validator.Apply(this_.GetReferencedType())
				_ = unused
			}
		} else {
			this_.Warning(reporter, "Missing type for `typeof` value. The value must be declared and const.")
			this_.priv_NamedType_setReferencedAndResolvedType_9297(this_.GetJSTypeBase().registry.GetNativeType(StaticGet_JSTypeNative_UNKNOWN_TYPE()), reporter)
		}
	} else {
		if type_.GetJSTypeBase().IsLiteralObject() {
			var objlit JSType = type_
			_ = objlit
			type_ = ToIface_JSType(NamedType_S_Builder(this_.GetJSTypeBase().registry, this_.GetReferenceName()).SetResolutionKind(StaticGet_NamedType_ResolutionKind_NONE()).SetReferencedType(objlit).Build())
		}
		this_.priv_NamedType_setReferencedAndResolvedType_9297(type_, reporter)
	}
	return true
}

// Java method: NamedType.setReferencedAndResolvedType(com.google.javascript.rhino.jstype.JSType type, com.google.javascript.rhino.ErrorReporter reporter)
func (this_ *NamedType) priv_NamedType_setReferencedAndResolvedType_9297(type_ JSType, reporter ErrorReporter) {
	_ = type_
	_ = reporter
	if this_.restrictByNull {
		type_ = type_.RestrictByNotNullOrUndefined()
	}
	if this_.validator != nil {
		var unused bool = this_.validator.Apply(type_)
		_ = unused
	}
	this_.SetReferencedType(type_)
	this_.priv_NamedType_checkEnumElementCycle_9299(reporter)
	this_.priv_NamedType_checkProtoCycle_9300(reporter)
}

// Java method: NamedType.handleTypeCycle(com.google.javascript.rhino.ErrorReporter reporter)
func (this_ *NamedType) priv_NamedType_handleTypeCycle_9298(reporter ErrorReporter) {
	_ = reporter
	this_.SetReferencedType(Upcast_JSType(this_.GetJSTypeBase().registry.GetNativeObjectType(StaticGet_JSTypeNative_UNKNOWN_TYPE())))
	this_.Warning(reporter, ("Cycle detected in inheritance chain of type " + j2gort.StrValString(this_.reference)))
}

// Java method: NamedType.checkEnumElementCycle(com.google.javascript.rhino.ErrorReporter reporter)
func (this_ *NamedType) priv_NamedType_checkEnumElementCycle_9299(reporter ErrorReporter) {
	_ = reporter
	var referencedType JSType = this_.GetReferencedType()
	_ = referencedType
	var enumElementType *EnumElementType
	_ = enumElementType
	var ok_9 bool
	enumElementType, ok_9 = j2gort.TryCast[*EnumElementType](referencedType)
	if ok_9 && (j2gort.IfaceEquals[JSType](this_, enumElementType.GetPrimitiveType())) {
		this_.priv_NamedType_handleTypeCycle_9298(reporter)
	}
}

// Java method: NamedType.checkProtoCycle(com.google.javascript.rhino.ErrorReporter reporter)
func (this_ *NamedType) priv_NamedType_checkProtoCycle_9300(reporter ErrorReporter) {
	_ = reporter
	var referencedType JSType = this_.GetReferencedType()
	_ = referencedType
	if referencedType == this_ || ((referencedType == nil || referencedType.IsJavaNil_()) && (this_ == nil || this_.IsJavaNil_())) {
		this_.priv_NamedType_handleTypeCycle_9298(reporter)
	}
}

// Java method: NamedType.handleUnresolvedType(com.google.javascript.rhino.ErrorReporter reporter)
func (this_ *NamedType) priv_NamedType_handleUnresolvedType_9301(reporter ErrorReporter) {
	_ = reporter
	var isForwardDeclared bool = this_.GetJSTypeBase().registry.IsForwardDeclaredType(this_.reference)
	_ = isForwardDeclared
	if !(isForwardDeclared) {
		var msg string = ("Bad type annotation. Unknown type " + j2gort.StrValString(this_.reference))
		_ = msg
		var root string
		_ = root
		if j2gort.JStrIndexOf(this_.reference, ".") >= 0 {
			root = j2gort.JStrSub(this_.reference, int32(0), j2gort.JStrIndexOf(this_.reference, "."))
		} else {
			root = this_.reference
		}
		if this_.priv_NamedType_localVariableShadowsGlobalNamespace_9302(root) {
			msg = (j2gort.StrValString(msg) + "\nIt's possible that a local variable called '" + j2gort.StrValString(root) + "' is shadowing the intended global namespace.")
		} else if !(j2gort.IfaceEquals[StaticTypedSlot](j2gort.AsIface[StaticTypedSlot](this_.resolutionScope.GetSlot(root)), nil)) {
			msg = (j2gort.StrValString(msg) + "\nIt's possible that '" + j2gort.StrValString(this_.reference) + "' refers to a value, not a type.")
		}
		this_.Warning(reporter, msg)
	} else {
		this_.SetReferencedType(NewNoResolvedType(this_.GetJSTypeBase().registry, this_.GetReferenceName(), this_.GetTemplateTypes()))
		if this_.validator != nil {
			var unused bool = this_.validator.Apply(this_.GetReferencedType())
			_ = unused
		}
	}
}

// Java method: NamedType.localVariableShadowsGlobalNamespace(java.lang.String root)
func (this_ *NamedType) priv_NamedType_localVariableShadowsGlobalNamespace_9302(root string) bool {
	_ = root
	var rootVar StaticSlot = Upcast_StaticSlot(j2gort.AsIface[StaticTypedSlot](this_.resolutionScope.GetSlot(root)))
	_ = rootVar
	if rootVar != nil && !rootVar.IsJavaNil_() {
		j2gort.CheckNotNullIface(j2gort.AsIface[StaticScope](rootVar.GetScope()))
		var parent StaticScope = j2gort.AsIface[StaticScope](j2gort.AsIface[StaticScope](rootVar.GetScope()).GetParentScope())
		_ = parent
		if parent != nil && !parent.IsJavaNil_() {
			var globalVar StaticSlot = j2gort.AsIface[StaticSlot](parent.GetSlot(root))
			_ = globalVar
			return (globalVar != nil && !globalVar.IsJavaNil_())
		}
	}
	return false
}

// Java method: NamedType.setValidator(com.google.common.base.Predicate<com.google.javascript.rhino.jstype.JSType> validator)
func (this_ *NamedType) SetValidator(validator j2gort.Ext_com_google_common_base_Predicate) bool {
	_ = validator
	if this_.IsResolved() {
		return this_.ProxyObjectTypeBase.SetValidator(validator)
	} else {
		this_.validator = validator
		return true
	}
}

// Java method: NamedType.warning(com.google.javascript.rhino.ErrorReporter reporter, java.lang.String message)
func (this_ *NamedType) Warning(reporter ErrorReporter, message string) {
	_ = reporter
	_ = message
	reporter.Warning(message, this_.sourceName, this_.lineno, this_.charno)
}

// Java method: NamedType.isObject()
func (this_ *NamedType) IsObject() bool {
	if this_.IsEnumElementType() {
		return this_.ToMaybeEnumElementType().IsObject()
	}
	return this_.ProxyObjectTypeBase.IsObject()
}

// Java method: NamedType.visit(com.google.javascript.rhino.jstype.Visitor<T> visitor)
func (this_ *NamedType) Visit_1(visitor Visitor) any {
	_ = visitor
	return visitor.CaseNamedType(this_)
}

// Java static method: NamedType.builder(com.google.javascript.rhino.jstype.JSTypeRegistry registry, java.lang.String reference)
func NamedType_S_Builder(registry *JSTypeRegistry, reference string) *NamedType_Builder {
	_ = registry
	_ = reference
	return NewNamedType_Builder(registry, reference)
}

// Java method: NamedType.toBuilder()
func (this_ *NamedType) ToBuilder() *NamedType_Builder {
	if !(!(this_.IsResolved())) {
		j2gort.CheckStateStr(false, "Only call toBuilder on unresolved NamedTypes")
	}
	return NewNamedType_Builder(this_.GetJSTypeBase().registry, this_.reference).SetScope(this_.resolutionScope).SetResolutionKind(this_.resolutionKind).SetErrorReportingLocation(this_.sourceName, this_.lineno, this_.charno).SetTemplateTypes(j2gort.ImmListFromTypedSlice(this_.templateTypes)).SetReferencedType(this_.GetReferencedType()).priv_Builder_setRestrictByNull_9312(this_.restrictByNull)
}

// Java class: com.google.javascript.rhino.jstype.NamedType.PropertyContinuation
type NamedType_PropertyContinuation struct {
	propertyName Property_Key `gson:"propertyName"`
	type_        JSType       `gson:"type"`
	propertyNode Node         `gson:"propertyNode"`
	idHash_      atomic.Int32
	inferred     bool `gson:"inferred"`
}

func (this_ *NamedType_PropertyContinuation) BindSelf_(s any) {
	_ = s
}
func (this_ *NamedType_PropertyContinuation) IsJavaNil_() bool { return this_ == nil }
func (*NamedType_PropertyContinuation) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.NamedType$PropertyContinuation", "PropertyContinuation"
}
func (this_ *NamedType_PropertyContinuation) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*NamedType_PropertyContinuation) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.NamedType.PropertyContinuation", "com.google.javascript.rhino.jstype.NamedType$PropertyContinuation", func(v_ any) bool {
		return j2gort.InstanceOf[*NamedType_PropertyContinuation](v_)
	})
}

func (this_ *NamedType_PropertyContinuation) initFields_NamedType_PropertyContinuation() {
	_ = this_
}

// Java constructor: PropertyContinuation(com.google.javascript.rhino.jstype.Property.Key propertyName, com.google.javascript.rhino.jstype.JSType type, boolean inferred, com.google.javascript.rhino.Node propertyNode)
func NewNamedType_PropertyContinuation(propertyName Property_Key, type_ JSType, inferred bool, propertyNode Node) *NamedType_PropertyContinuation {
	this_ := &NamedType_PropertyContinuation{}
	InitNamedType_PropertyContinuation(this_, propertyName, type_, inferred, propertyNode)
	return this_
}

func InitNamedType_PropertyContinuation(this_ *NamedType_PropertyContinuation, propertyName Property_Key, type_ JSType, inferred bool, propertyNode Node) {
	_ = propertyName
	_ = type_
	_ = inferred
	_ = propertyNode
	this_.initFields_NamedType_PropertyContinuation()
	this_.propertyName = propertyName
	this_.type_ = type_
	this_.inferred = inferred
	this_.propertyNode = propertyNode
}

// Java method: PropertyContinuation.commit(com.google.javascript.rhino.jstype.ObjectType target)
func (this_ *NamedType_PropertyContinuation) Commit(target ObjectType) {
	_ = target
	target.DefineProperty_4_Property_Key_JSType_bool_Node(this_.propertyName, this_.type_, this_.inferred, this_.propertyNode)
}

// Java enum: com.google.javascript.rhino.jstype.NamedType.ResolutionKind
type NamedType_ResolutionKind uint8

const (
	enum_NamedType_ResolutionKind_NONE      NamedType_ResolutionKind = 1
	enum_NamedType_ResolutionKind_TYPE_NAME NamedType_ResolutionKind = 2
	enum_NamedType_ResolutionKind_TYPEOF    NamedType_ResolutionKind = 3
)

var enumBaseTable_NamedType_ResolutionKind = [...]j2gort.EnumBase{
	{Name_: "NONE", Ordinal_: 0},
	{Name_: "TYPE_NAME", Ordinal_: 1},
	{Name_: "TYPEOF", Ordinal_: 2},
}
var enumValues_NamedType_ResolutionKind = []NamedType_ResolutionKind{1, 2, 3}

func StaticGet_NamedType_ResolutionKind_NONE() NamedType_ResolutionKind {
	return enum_NamedType_ResolutionKind_NONE
}
func StaticGet_NamedType_ResolutionKind_TYPE_NAME() NamedType_ResolutionKind {
	return enum_NamedType_ResolutionKind_TYPE_NAME
}
func StaticGet_NamedType_ResolutionKind_TYPEOF() NamedType_ResolutionKind {
	return enum_NamedType_ResolutionKind_TYPEOF
}

func NamedType_ResolutionKind_Values() []NamedType_ResolutionKind {
	return j2gort.CloneSlice(enumValues_NamedType_ResolutionKind)
}

func NamedType_ResolutionKind_ValueOf(name string) NamedType_ResolutionKind {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_NamedType_ResolutionKind {
		if enumBaseTable_NamedType_ResolutionKind[i_].Name_ == name {
			return NamedType_ResolutionKind(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.rhino.jstype.NamedType.ResolutionKind." + name))
}

func (this_ NamedType_ResolutionKind) IsJavaNil_() bool                              { return this_ == 0 }
func (NamedType_ResolutionKind) HasIdentityEquals_()                                 {}
func (NamedType_ResolutionKind) Implements_j2gort_Ext_java_lang_Enum()               {}
func (NamedType_ResolutionKind) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ NamedType_ResolutionKind) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (NamedType_ResolutionKind) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.NamedType$ResolutionKind", "ResolutionKind"
}
func (this_ NamedType_ResolutionKind) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_NamedType_ResolutionKind[this_-1]
}
func (this_ NamedType_ResolutionKind) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ NamedType_ResolutionKind) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ NamedType_ResolutionKind) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_NamedType_ResolutionKind[this_-1].Name_
}
func (this_ NamedType_ResolutionKind) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_NamedType_ResolutionKind[this_-1].Name_
}
func (this_ NamedType_ResolutionKind) Equals(other_ any) bool {
	o_, ok_ := other_.(NamedType_ResolutionKind)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ NamedType_ResolutionKind) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_NamedType_ResolutionKind[this_-1].HashCode()
}
func (this_ NamedType_ResolutionKind) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_NamedType_ResolutionKind[this_-1].CompareTo(other_)
}
func (this_ NamedType_ResolutionKind) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_NamedType_ResolutionKind[this_-1].Name_)
}
func init() {
	vals_ := []any{NamedType_ResolutionKind(1), NamedType_ResolutionKind(2), NamedType_ResolutionKind(3)}
	j2gort.RegisterEnumConstants("com.google.javascript.rhino.jstype.NamedType.ResolutionKind", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.rhino.jstype.NamedType$ResolutionKind", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.NamedType.ResolutionKind", "com.google.javascript.rhino.jstype.NamedType$ResolutionKind", func(v_ any) bool {
		return j2gort.InstanceOf[NamedType_ResolutionKind](v_)
	})
}

// Java class: com.google.javascript.rhino.jstype.NamedType.Builder
type NamedType_Builder struct {
	registry       *JSTypeRegistry  `gson:"registry"`
	referenceName  string           `gson:"referenceName"`
	scope          StaticTypedScope `gson:"scope"`
	sourceName     string           `gson:"sourceName"`
	referencedType JSType           `gson:"referencedType"`
	templateTypes  []JSType         `gson:"templateTypes"`
	lineno         int32            `gson:"lineno"`
	charno         int32            `gson:"charno"`
	idHash_        atomic.Int32
	resolutionKind NamedType_ResolutionKind `gson:"resolutionKind"`
	restrictByNull bool                     `gson:"restrictByNull"`
}

func (this_ *NamedType_Builder) BindSelf_(s any) {
	_ = s
}
func (this_ *NamedType_Builder) IsJavaNil_() bool { return this_ == nil }
func (*NamedType_Builder) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.NamedType$Builder", "Builder"
}
func (this_ *NamedType_Builder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*NamedType_Builder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.NamedType.Builder", "com.google.javascript.rhino.jstype.NamedType$Builder", func(v_ any) bool {
		return j2gort.InstanceOf[*NamedType_Builder](v_)
	})
}

func (this_ *NamedType_Builder) initFields_NamedType_Builder() {
	_ = this_
	this_.referenceName = j2gort.NullString
	this_.sourceName = j2gort.NullString
	this_.templateTypes = []JSType{}
}

// Java constructor: Builder(com.google.javascript.rhino.jstype.JSTypeRegistry registry, java.lang.String referenceName)
func NewNamedType_Builder(registry *JSTypeRegistry, referenceName string) *NamedType_Builder {
	this_ := &NamedType_Builder{}
	InitNamedType_Builder(this_, registry, referenceName)
	return this_
}

func InitNamedType_Builder(this_ *NamedType_Builder, registry *JSTypeRegistry, referenceName string) {
	_ = registry
	_ = referenceName
	this_.initFields_NamedType_Builder()
	this_.registry = registry
	this_.referenceName = referenceName
	this_.referencedType = registry.GetNativeType(StaticGet_JSTypeNative_UNKNOWN_TYPE())
}

// Java method: Builder.setScope(com.google.javascript.rhino.jstype.StaticTypedScope scope)
func (this_ *NamedType_Builder) SetScope(scope StaticTypedScope) *NamedType_Builder {
	_ = scope
	this_.scope = scope
	return this_
}

// Java method: Builder.setResolutionKind(com.google.javascript.rhino.jstype.NamedType.ResolutionKind resolutionKind)
func (this_ *NamedType_Builder) SetResolutionKind(resolutionKind NamedType_ResolutionKind) *NamedType_Builder {
	_ = resolutionKind
	this_.resolutionKind = resolutionKind
	return this_
}

// Java method: Builder.setErrorReportingLocation(java.lang.String sourceName, int lineno, int charno)
func (this_ *NamedType_Builder) SetErrorReportingLocation(sourceName string, lineno int32, charno int32) *NamedType_Builder {
	_ = sourceName
	_ = lineno
	_ = charno
	this_.sourceName = sourceName
	this_.lineno = lineno
	this_.charno = charno
	return this_
}

// Java method: Builder.setErrorReportingLocationFrom(com.google.javascript.rhino.Node source)
func (this_ *NamedType_Builder) SetErrorReportingLocationFrom(source Node) *NamedType_Builder {
	_ = source
	this_.sourceName = source.GetSourceFileName()
	this_.lineno = source.GetLineno()
	this_.charno = source.GetCharno()
	return this_
}

// Java method: Builder.setTemplateTypes(com.google.common.collect.ImmutableList<com.google.javascript.rhino.jstype.JSType> templateTypes)
func (this_ *NamedType_Builder) SetTemplateTypes(templateTypes j2gort.Ext_com_google_common_collect_ImmutableList) *NamedType_Builder {
	_ = templateTypes
	this_.templateTypes = j2gort.ImmListToSliceNullable[JSType](templateTypes)
	return this_
}

// Java method: Builder.setReferencedType(com.google.javascript.rhino.jstype.JSType referencedType)
func (this_ *NamedType_Builder) SetReferencedType(referencedType JSType) *NamedType_Builder {
	_ = referencedType
	this_.referencedType = referencedType
	return this_
}

// Java method: Builder.setRestrictByNull(boolean restrictByNull)
func (this_ *NamedType_Builder) priv_Builder_setRestrictByNull_9312(restrictByNull bool) *NamedType_Builder {
	_ = restrictByNull
	this_.restrictByNull = restrictByNull
	return this_
}

// Java method: Builder.build()
func (this_ *NamedType_Builder) Build() *NamedType {
	return NewNamedType(this_)
}
func (this_ *NamedType) Slice_GetTemplateTypes() []JSType {
	return this_.templateTypes
}
