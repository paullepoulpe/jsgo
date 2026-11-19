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
// Port of com/google/javascript/jscomp/type/ChainableReverseAbstractInterpreter.java (generated mechanically by j2go).

package jscomptype

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java polymorphic class interface: com.google.javascript.jscomp.type.ChainableReverseAbstractInterpreter
type ChainableReverseAbstractInterpreter interface {
	Implements_jscomptype_ChainableReverseAbstractInterpreter()
	IsJavaNil_() bool
	GetChainableReverseAbstractInterpreterBase() *ChainableReverseAbstractInterpreterBase
	ReverseAbstractInterpreter
	Append(ChainableReverseAbstractInterpreter) ChainableReverseAbstractInterpreter
	GetFirst() ChainableReverseAbstractInterpreter
	FirstPreciserScopeKnowingConditionOutcome(pkg_rhino.Node, FlowScope, *pkg_rhino.Outcome) FlowScope
	NextPreciserScopeKnowingConditionOutcome(pkg_rhino.Node, FlowScope, *pkg_rhino.Outcome) FlowScope
	GetTypeIfRefinable(pkg_rhino.Node, FlowScope) pkg_rhino.JSType
	DeclareNameInScope(FlowScope, pkg_rhino.Node, pkg_rhino.JSType) FlowScope
	GetRestrictedByTypeOfResult(pkg_rhino.JSType, string, bool) pkg_rhino.JSType
	GetNativeType(pkg_rhino.JSTypeNative) pkg_rhino.JSType
}

func ToIface_ChainableReverseAbstractInterpreter[T interface {
	ChainableReverseAbstractInterpreter
	comparable
}](v T) ChainableReverseAbstractInterpreter {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_ChainableReverseAbstractInterpreter[T ChainableReverseAbstractInterpreter](v T) ChainableReverseAbstractInterpreter {
	var vi ChainableReverseAbstractInterpreter = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.jscomp.type.ChainableReverseAbstractInterpreter
type ChainableReverseAbstractInterpreterBase struct {
	self_                     ChainableReverseAbstractInterpreter
	typeRegistry              *pkg_rhino.JSTypeRegistry           `gson:"typeRegistry"`
	firstLink                 ChainableReverseAbstractInterpreter `gson:"firstLink"`
	nextLink                  ChainableReverseAbstractInterpreter `gson:"nextLink"`
	restrictToArrayVisitor    pkg_rhino.Visitor                   `gson:"restrictToArrayVisitor"`
	restrictToNotArrayVisitor pkg_rhino.Visitor                   `gson:"restrictToNotArrayVisitor"`
	idHash_                   atomic.Int32
}

func (this_ *ChainableReverseAbstractInterpreterBase) BindSelf_(s any) {
	_ = s
	this_.self_ = s.(ChainableReverseAbstractInterpreter)
}
func (this_ *ChainableReverseAbstractInterpreterBase) IsJavaNil_() bool { return this_ == nil }
func (*ChainableReverseAbstractInterpreterBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.type.ChainableReverseAbstractInterpreter", "ChainableReverseAbstractInterpreter"
}
func (this_ *ChainableReverseAbstractInterpreterBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.type.ChainableReverseAbstractInterpreter", "com.google.javascript.jscomp.type.ChainableReverseAbstractInterpreter", func(v_ any) bool {
		return j2gort.InstanceOf[ChainableReverseAbstractInterpreter](v_)
	})
}
func (b *ChainableReverseAbstractInterpreterBase) GetChainableReverseAbstractInterpreterBase() *ChainableReverseAbstractInterpreterBase {
	return b
}
func (*ChainableReverseAbstractInterpreterBase) Implements_jscomptype_ChainableReverseAbstractInterpreter() {
}
func (*ChainableReverseAbstractInterpreterBase) Implements_jscomptype_ReverseAbstractInterpreter() {}

func (this_ *ChainableReverseAbstractInterpreterBase) initFields_ChainableReverseAbstractInterpreter() {
	_ = this_
	this_.restrictToArrayVisitor = NewChainableReverseAbstractInterpreter__Anon_129(this_.self_, this_.self_)
	this_.restrictToNotArrayVisitor = NewChainableReverseAbstractInterpreter__Anon_130(this_.self_, this_.self_)
}

func InitChainableReverseAbstractInterpreter(this_ *ChainableReverseAbstractInterpreterBase, typeRegistry *pkg_rhino.JSTypeRegistry) {
	_ = typeRegistry
	this_.initFields_ChainableReverseAbstractInterpreter()
	this_.typeRegistry = typeRegistry
	this_.firstLink = this_.self_
	this_.nextLink = nil
}

// Java method: ChainableReverseAbstractInterpreter.append(com.google.javascript.jscomp.type.ChainableReverseAbstractInterpreter lastLink)
func (this_ *ChainableReverseAbstractInterpreterBase) Append(lastLink ChainableReverseAbstractInterpreter) ChainableReverseAbstractInterpreter {
	_ = lastLink
	if !(j2gort.IfaceEquals[ChainableReverseAbstractInterpreter](lastLink.GetChainableReverseAbstractInterpreterBase().nextLink, nil)) {
		j2gort.CheckArgument(false)
	}
	this_.nextLink = lastLink
	{
		r_1 := lastLink
		r_1.GetChainableReverseAbstractInterpreterBase().firstLink = this_.firstLink
	}
	return lastLink
}

// Java method: ChainableReverseAbstractInterpreter.getFirst()
func (this_ *ChainableReverseAbstractInterpreterBase) GetFirst() ChainableReverseAbstractInterpreter {
	return this_.firstLink
}

// Java method: ChainableReverseAbstractInterpreter.firstPreciserScopeKnowingConditionOutcome(com.google.javascript.rhino.Node condition, com.google.javascript.jscomp.type.FlowScope blindScope, com.google.javascript.rhino.Outcome outcome)
func (this_ *ChainableReverseAbstractInterpreterBase) FirstPreciserScopeKnowingConditionOutcome(condition pkg_rhino.Node, blindScope FlowScope, outcome *pkg_rhino.Outcome) FlowScope {
	_ = condition
	_ = blindScope
	_ = outcome
	return this_.firstLink.GetPreciserScopeKnowingConditionOutcome(condition, blindScope, outcome)
}

// Java method: ChainableReverseAbstractInterpreter.nextPreciserScopeKnowingConditionOutcome(com.google.javascript.rhino.Node condition, com.google.javascript.jscomp.type.FlowScope blindScope, com.google.javascript.rhino.Outcome outcome)
func (this_ *ChainableReverseAbstractInterpreterBase) NextPreciserScopeKnowingConditionOutcome(condition pkg_rhino.Node, blindScope FlowScope, outcome *pkg_rhino.Outcome) FlowScope {
	_ = condition
	_ = blindScope
	_ = outcome
	if this_.nextLink != nil && !this_.nextLink.IsJavaNil_() {
		return this_.nextLink.GetPreciserScopeKnowingConditionOutcome(condition, blindScope, outcome)
	} else {
		return blindScope
	}
}

// Java method: ChainableReverseAbstractInterpreter.getTypeIfRefinable(com.google.javascript.rhino.Node node, com.google.javascript.jscomp.type.FlowScope scope)
func (this_ *ChainableReverseAbstractInterpreterBase) GetTypeIfRefinable(node pkg_rhino.Node, scope FlowScope) pkg_rhino.JSType {
	_ = node
	_ = scope
	switch node.TokenSwitchOrd_() {
	case int32(32) /* Token.NAME */ :
		var nameVar pkg_rhino.StaticTypedSlot = j2gort.AsIface[pkg_rhino.StaticTypedSlot](scope.GetSlot(node.GetString()))
		_ = nameVar
		if nameVar != nil && !nameVar.IsJavaNil_() {
			var nameVarType pkg_rhino.JSType = nameVar.GetType()
			_ = nameVarType
			if nameVarType == nil || nameVarType.IsJavaNil_() {
				nameVarType = node.GetJSType()
			}
			return nameVarType
		}
		return nil
	case int32(26) /* Token.GETPROP */ :
		var qualifiedName string = node.GetQualifiedName()
		_ = qualifiedName
		if qualifiedName == j2gort.NullString {
			return nil
		}
		var propVar pkg_rhino.StaticTypedSlot = j2gort.AsIface[pkg_rhino.StaticTypedSlot](scope.GetSlot(qualifiedName))
		_ = propVar
		var propVarType pkg_rhino.JSType = nil
		_ = propVarType
		if propVar != nil && !propVar.IsJavaNil_() {
			propVarType = propVar.GetType()
		}
		if propVarType == nil || propVarType.IsJavaNil_() {
			propVarType = node.GetJSType()
		}
		if propVarType == nil || propVarType.IsJavaNil_() {
			propVarType = this_.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_UNKNOWN_TYPE())
		}
		return propVarType
	default:
	}
	return nil
}

// Java method: ChainableReverseAbstractInterpreter.declareNameInScope(com.google.javascript.jscomp.type.FlowScope scope, com.google.javascript.rhino.Node node, com.google.javascript.rhino.jstype.JSType type)
func (this_ *ChainableReverseAbstractInterpreterBase) DeclareNameInScope(scope FlowScope, node pkg_rhino.Node, type_ pkg_rhino.JSType) FlowScope {
	_ = scope
	_ = node
	_ = type_
	switch node.TokenSwitchOrd_() {
	case int32(32) /* Token.NAME */ :
		return j2gort.AsIface[FlowScope](scope.InferSlotType(node.GetString(), type_))
	case int32(26) /* Token.GETPROP */ :
		var qualifiedName string = node.GetQualifiedName()
		_ = qualifiedName
		j2gort.CheckNotNullString(qualifiedName)
		var origType pkg_rhino.JSType = node.GetJSType()
		_ = origType
		if origType == nil || origType.IsJavaNil_() {
			origType = this_.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_UNKNOWN_TYPE())
		} else {
			origType = origType
		}
		return j2gort.AsIface[FlowScope](scope.InferQualifiedSlot(node, qualifiedName, origType, type_, false))
	case int32(37) /* Token.THIS */ :
		return scope
	default:
		panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", "Node cannot be refined. \n"+j2gort.StrValString(node.ToStringTree()))))
	}
	panic("unreachable switch expression")
}

// Java method: ChainableReverseAbstractInterpreter.getRestrictedByTypeOfResult(com.google.javascript.rhino.jstype.JSType type, java.lang.String value, boolean resultEqualsValue)
func (this_ *ChainableReverseAbstractInterpreterBase) GetRestrictedByTypeOfResult(type_ pkg_rhino.JSType, value string, resultEqualsValue bool) pkg_rhino.JSType {
	_ = type_
	_ = value
	_ = resultEqualsValue
	if type_ == nil || type_.IsJavaNil_() {
		if resultEqualsValue {
			var result pkg_rhino.JSType = this_.priv_ChainableReverseAbstractInterpreter_getNativeTypeForTypeOf_7827(value)
			_ = result
			if result == nil || result.IsJavaNil_() {
				return this_.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_CHECKED_UNKNOWN_TYPE())
			} else {
				return result
			}
		} else {
			return nil
		}
	}
	return j2gort.AsIface[pkg_rhino.JSType](type_.Visit_1(NewChainableReverseAbstractInterpreter_RestrictByOneTypeOfResultVisitor(this_.self_, value, resultEqualsValue)))
}

// Java method: ChainableReverseAbstractInterpreter.getNativeType(com.google.javascript.rhino.jstype.JSTypeNative typeId)
func (this_ *ChainableReverseAbstractInterpreterBase) GetNativeType(typeId pkg_rhino.JSTypeNative) pkg_rhino.JSType {
	_ = typeId
	return this_.typeRegistry.GetNativeType(typeId)
}

// Java method: ChainableReverseAbstractInterpreter.getNativeTypeForTypeOf(java.lang.String value)
func (this_ *ChainableReverseAbstractInterpreterBase) priv_ChainableReverseAbstractInterpreter_getNativeTypeForTypeOf_7827(value string) pkg_rhino.JSType {
	_ = value
	switch j2gort.CheckNotNullString(value) {
	case "number":
		return this_.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_NUMBER_TYPE())
	case "boolean":
		return this_.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_BOOLEAN_TYPE())
	case "string":
		return this_.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_STRING_TYPE())
	case "symbol":
		return this_.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_SYMBOL_TYPE())
	case "undefined":
		return this_.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_VOID_TYPE())
	case "object":
		return this_.typeRegistry.CreateUnionType_2(this_.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_OBJECT_TYPE()), this_.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_NULL_TYPE()))
	case "function":
		return this_.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_FUNCTION_TYPE())
	default:
		return nil
	}
	panic("unreachable switch expression")
}

// Java polymorphic class interface: com.google.javascript.jscomp.type.ChainableReverseAbstractInterpreter.RestrictByTypeOfResultVisitor
type ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitor interface {
	Implements_jscomptype_ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitor()
	IsJavaNil_() bool
	GetChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitorBase() *ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitorBase
	pkg_rhino.Visitor
	CaseTopType(pkg_rhino.JSType) pkg_rhino.JSType
}

func ToIface_ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitor[T interface {
	ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitor
	comparable
}](v T) ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitor {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitor[T ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitor](v T) ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitor {
	var vi ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitor = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.jscomp.type.ChainableReverseAbstractInterpreter.RestrictByTypeOfResultVisitor
type ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitorBase struct {
	self_   ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitor
	outer_  ChainableReverseAbstractInterpreter
	idHash_ atomic.Int32
}

func (this_ *ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitorBase) BindSelf_(s any) {
	_ = s
	this_.self_ = s.(ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitor)
}
func (this_ *ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitorBase) IsJavaNil_() bool {
	return this_ == nil
}
func (*ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitorBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.type.ChainableReverseAbstractInterpreter$RestrictByTypeOfResultVisitor", "RestrictByTypeOfResultVisitor"
}
func (this_ *ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitorBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.type.ChainableReverseAbstractInterpreter.RestrictByTypeOfResultVisitor", "com.google.javascript.jscomp.type.ChainableReverseAbstractInterpreter$RestrictByTypeOfResultVisitor", func(v_ any) bool {
		return j2gort.InstanceOf[ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitor](v_)
	})
}
func (b *ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitorBase) GetChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitorBase() *ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitorBase {
	return b
}
func (this_ *ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitorBase) SetOuter_(o ChainableReverseAbstractInterpreter) {
	this_.outer_ = o
}
func (this_ *ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitorBase) GetOuter_() ChainableReverseAbstractInterpreter {
	return this_.outer_
}
func (*ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitorBase) Implements_jscomptype_ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitor() {
}
func (*ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitorBase) Implements_rhino_Visitor() {
}

func (this_ *ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitorBase) initFields_ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitor() {
	_ = this_
}

func InitChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitor(this_ *ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitorBase) {
	this_.initFields_ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitor()
}

// Java method: RestrictByTypeOfResultVisitor.caseTopType(com.google.javascript.rhino.jstype.JSType topType)
func (this_ *ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitorBase) CaseTopType(topType pkg_rhino.JSType) pkg_rhino.JSType {
	_ = topType
	panic("abstract or native method invoked")
}

// Java method: RestrictByTypeOfResultVisitor.caseAllType()
func (this_ *ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitorBase) CaseAllType() any {
	return this_.self_.CaseTopType(this_.outer_.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_ALL_TYPE()))
}

// Java method: RestrictByTypeOfResultVisitor.caseUnknownType()
func (this_ *ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitorBase) CaseUnknownType() any {
	return this_.self_.CaseTopType(this_.outer_.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_CHECKED_UNKNOWN_TYPE()))
}

// Java method: RestrictByTypeOfResultVisitor.caseUnionType(com.google.javascript.rhino.jstype.UnionType type)
func (this_ *ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitorBase) CaseUnionType(type_ *pkg_rhino.UnionType) any {
	_ = type_
	var restricted pkg_rhino.JSType = nil
	_ = restricted
	{
		var i int32 = int32(0)
		_ = i
		var n int32 = type_.GetAlternatesCount()
		_ = n
		for ; i < n; i = i + 1 {
			var restrictedAlternate pkg_rhino.JSType = j2gort.AsIface[pkg_rhino.JSType](type_.GetAlternateAt(i).Visit_1(pkg_rhino.Upcast_Visitor(this_.self_)))
			_ = restrictedAlternate
			if restrictedAlternate != nil && !restrictedAlternate.IsJavaNil_() {
				if restricted == nil || restricted.IsJavaNil_() {
					restricted = restrictedAlternate
				} else {
					restricted = restrictedAlternate.GetLeastSupertype(restricted)
				}
			}
		}
	}
	return restricted
}

// Java method: RestrictByTypeOfResultVisitor.caseNoType(com.google.javascript.rhino.jstype.FunctionType type)
func (this_ *ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitorBase) CaseNoType(type_ *pkg_rhino.FunctionType) any {
	_ = type_
	return j2gort.PtrToAny(type_)
}

// Java method: RestrictByTypeOfResultVisitor.caseEnumElementType(com.google.javascript.rhino.jstype.EnumElementType enumElementType)
func (this_ *ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitorBase) CaseEnumElementType(enumElementType *pkg_rhino.EnumElementType) any {
	_ = enumElementType
	var type_ pkg_rhino.JSType = j2gort.AsIface[pkg_rhino.JSType](enumElementType.GetPrimitiveType().Visit_1(pkg_rhino.Upcast_Visitor(this_.self_)))
	_ = type_
	if (type_ != nil && !type_.IsJavaNil_()) && enumElementType.GetPrimitiveType().GetJSTypeBase().EqualsJSType(type_) {
		return j2gort.PtrToAny(enumElementType)
	} else {
		return type_
	}
}

// Java method: RestrictByTypeOfResultVisitor.caseTemplatizedType(com.google.javascript.rhino.jstype.TemplatizedType type)
func (this_ *ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitorBase) CaseTemplatizedType(type_ *pkg_rhino.TemplatizedType) any {
	_ = type_
	return j2gort.AsIface[pkg_rhino.JSType](this_.self_.CaseObjectType(pkg_rhino.ToIface_ObjectType(type_)))
}

// Java method: RestrictByTypeOfResultVisitor.caseTemplateType(com.google.javascript.rhino.jstype.TemplateType templateType)
func (this_ *ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitorBase) CaseTemplateType(templateType *pkg_rhino.TemplateType) any {
	_ = templateType
	return j2gort.AsIface[pkg_rhino.JSType](this_.self_.CaseObjectType(pkg_rhino.ToIface_ObjectType(templateType)))
}

// Java method: RestrictByTypeOfResultVisitor.caseNamedType(com.google.javascript.rhino.jstype.NamedType type)
func (this_ *ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitorBase) CaseNamedType(type_ *pkg_rhino.NamedType) any {
	_ = type_
	return j2gort.AsIface[pkg_rhino.JSType](this_.CaseProxyObjectType(pkg_rhino.ToIface_ProxyObjectType(type_)))
}

// Java method: RestrictByTypeOfResultVisitor.caseProxyObjectType(com.google.javascript.rhino.jstype.ProxyObjectType type)
func (this_ *ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitorBase) CaseProxyObjectType(type_ pkg_rhino.ProxyObjectType) any {
	_ = type_
	return j2gort.AsIface[pkg_rhino.JSType](type_.GetProxyObjectTypeBase().VisitReferenceType(pkg_rhino.Upcast_Visitor(this_.self_)))
}

// Java polymorphic class interface: com.google.javascript.jscomp.type.ChainableReverseAbstractInterpreter.RestrictByTrueTypeOfResultVisitor
type ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitor interface {
	Implements_jscomptype_ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitor()
	IsJavaNil_() bool
	GetChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitorBase() *ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitorBase
	ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitor
}

func ToIface_ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitor[T interface {
	ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitor
	comparable
}](v T) ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitor {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitor[T ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitor](v T) ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitor {
	var vi ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitor = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.jscomp.type.ChainableReverseAbstractInterpreter.RestrictByTrueTypeOfResultVisitor
type ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitorBase struct {
	ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitorBase
	outer_ ChainableReverseAbstractInterpreter
}

func (this_ *ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitorBase) BindSelf_(s any) {
	_ = s
	this_.ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitorBase.BindSelf_(s)
}
func (this_ *ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitorBase) IsJavaNil_() bool {
	return this_ == nil
}
func (*ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitorBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.type.ChainableReverseAbstractInterpreter$RestrictByTrueTypeOfResultVisitor", "RestrictByTrueTypeOfResultVisitor"
}
func (this_ *ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitorBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitorBase.IdentityHash_()
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.type.ChainableReverseAbstractInterpreter.RestrictByTrueTypeOfResultVisitor", "com.google.javascript.jscomp.type.ChainableReverseAbstractInterpreter$RestrictByTrueTypeOfResultVisitor", func(v_ any) bool {
		return j2gort.InstanceOf[ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitor](v_)
	})
}
func (b *ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitorBase) GetChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitorBase() *ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitorBase {
	return b
}
func (this_ *ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitorBase) SetOuter_(o ChainableReverseAbstractInterpreter) {
	this_.outer_ = o
}
func (this_ *ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitorBase) GetOuter_() ChainableReverseAbstractInterpreter {
	return this_.outer_
}
func (*ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitorBase) Implements_jscomptype_ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitor() {
}
func (*ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitorBase) Implements_jscomptype_ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitor() {
}
func (*ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitorBase) Implements_rhino_Visitor() {
}

func (this_ *ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitorBase) initFields_ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitor() {
	_ = this_
}

func InitChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitor(this_ *ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitorBase) {
	this_.ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitorBase.outer_ = this_.outer_
	InitChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitor(&this_.ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitorBase)
	this_.initFields_ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitor()
}

// Java method: RestrictByTrueTypeOfResultVisitor.caseNoObjectType()
func (this_ *ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitorBase) CaseNoObjectType() any {
	return nil
}

// Java method: RestrictByTrueTypeOfResultVisitor.caseBooleanType()
func (this_ *ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitorBase) CaseBooleanType() any {
	return nil
}

// Java method: RestrictByTrueTypeOfResultVisitor.caseFunctionType(com.google.javascript.rhino.jstype.FunctionType type)
func (this_ *ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitorBase) CaseFunctionType(type_ *pkg_rhino.FunctionType) any {
	_ = type_
	return nil
}

// Java method: RestrictByTrueTypeOfResultVisitor.caseNullType()
func (this_ *ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitorBase) CaseNullType() any {
	return nil
}

// Java method: RestrictByTrueTypeOfResultVisitor.caseNumberType()
func (this_ *ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitorBase) CaseNumberType() any {
	return nil
}

// Java method: RestrictByTrueTypeOfResultVisitor.caseBigIntType()
func (this_ *ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitorBase) CaseBigIntType() any {
	return nil
}

// Java method: RestrictByTrueTypeOfResultVisitor.caseObjectType(com.google.javascript.rhino.jstype.ObjectType type)
func (this_ *ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitorBase) CaseObjectType(type_ pkg_rhino.ObjectType) any {
	_ = type_
	return nil
}

// Java method: RestrictByTrueTypeOfResultVisitor.caseStringType()
func (this_ *ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitorBase) CaseStringType() any {
	return nil
}

// Java method: RestrictByTrueTypeOfResultVisitor.caseSymbolType(com.google.javascript.rhino.jstype.SymbolType type)
func (this_ *ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitorBase) CaseSymbolType(type_ pkg_rhino.SymbolType) any {
	_ = type_
	return nil
}

// Java method: RestrictByTrueTypeOfResultVisitor.caseVoidType()
func (this_ *ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitorBase) CaseVoidType() any {
	return nil
}

// Java polymorphic class interface: com.google.javascript.jscomp.type.ChainableReverseAbstractInterpreter.RestrictByFalseTypeOfResultVisitor
type ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitor interface {
	Implements_jscomptype_ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitor()
	IsJavaNil_() bool
	GetChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitorBase() *ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitorBase
	ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitor
}

func ToIface_ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitor[T interface {
	ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitor
	comparable
}](v T) ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitor {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitor[T ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitor](v T) ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitor {
	var vi ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitor = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.jscomp.type.ChainableReverseAbstractInterpreter.RestrictByFalseTypeOfResultVisitor
type ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitorBase struct {
	ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitorBase
	outer_ ChainableReverseAbstractInterpreter
}

func (this_ *ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitorBase) BindSelf_(s any) {
	_ = s
	this_.ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitorBase.BindSelf_(s)
}
func (this_ *ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitorBase) IsJavaNil_() bool {
	return this_ == nil
}
func (*ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitorBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.type.ChainableReverseAbstractInterpreter$RestrictByFalseTypeOfResultVisitor", "RestrictByFalseTypeOfResultVisitor"
}
func (this_ *ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitorBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitorBase.IdentityHash_()
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.type.ChainableReverseAbstractInterpreter.RestrictByFalseTypeOfResultVisitor", "com.google.javascript.jscomp.type.ChainableReverseAbstractInterpreter$RestrictByFalseTypeOfResultVisitor", func(v_ any) bool {
		return j2gort.InstanceOf[ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitor](v_)
	})
}
func (b *ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitorBase) GetChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitorBase() *ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitorBase {
	return b
}
func (this_ *ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitorBase) SetOuter_(o ChainableReverseAbstractInterpreter) {
	this_.outer_ = o
}
func (this_ *ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitorBase) GetOuter_() ChainableReverseAbstractInterpreter {
	return this_.outer_
}
func (*ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitorBase) Implements_jscomptype_ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitor() {
}
func (*ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitorBase) Implements_jscomptype_ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitor() {
}
func (*ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitorBase) Implements_rhino_Visitor() {
}

func (this_ *ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitorBase) initFields_ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitor() {
	_ = this_
}

func InitChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitor(this_ *ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitorBase) {
	this_.ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitorBase.outer_ = this_.outer_
	InitChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitor(&this_.ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitorBase)
	this_.initFields_ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitor()
}

// Java method: RestrictByFalseTypeOfResultVisitor.caseTopType(com.google.javascript.rhino.jstype.JSType topType)
func (this_ *ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitorBase) CaseTopType(topType pkg_rhino.JSType) pkg_rhino.JSType {
	_ = topType
	return topType
}

// Java method: RestrictByFalseTypeOfResultVisitor.caseNoObjectType()
func (this_ *ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitorBase) CaseNoObjectType() any {
	return this_.outer_.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_NO_OBJECT_TYPE())
}

// Java method: RestrictByFalseTypeOfResultVisitor.caseBooleanType()
func (this_ *ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitorBase) CaseBooleanType() any {
	return this_.outer_.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_BOOLEAN_TYPE())
}

// Java method: RestrictByFalseTypeOfResultVisitor.caseFunctionType(com.google.javascript.rhino.jstype.FunctionType type)
func (this_ *ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitorBase) CaseFunctionType(type_ *pkg_rhino.FunctionType) any {
	_ = type_
	return j2gort.PtrToAny(type_)
}

// Java method: RestrictByFalseTypeOfResultVisitor.caseNullType()
func (this_ *ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitorBase) CaseNullType() any {
	return this_.outer_.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_NULL_TYPE())
}

// Java method: RestrictByFalseTypeOfResultVisitor.caseNumberType()
func (this_ *ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitorBase) CaseNumberType() any {
	return this_.outer_.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_NUMBER_TYPE())
}

// Java method: RestrictByFalseTypeOfResultVisitor.caseBigIntType()
func (this_ *ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitorBase) CaseBigIntType() any {
	return this_.outer_.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_BIGINT_TYPE())
}

// Java method: RestrictByFalseTypeOfResultVisitor.caseObjectType(com.google.javascript.rhino.jstype.ObjectType type)
func (this_ *ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitorBase) CaseObjectType(type_ pkg_rhino.ObjectType) any {
	_ = type_
	return type_
}

// Java method: RestrictByFalseTypeOfResultVisitor.caseStringType()
func (this_ *ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitorBase) CaseStringType() any {
	return this_.outer_.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_STRING_TYPE())
}

// Java method: RestrictByFalseTypeOfResultVisitor.caseSymbolType(com.google.javascript.rhino.jstype.SymbolType type)
func (this_ *ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitorBase) CaseSymbolType(type_ pkg_rhino.SymbolType) any {
	_ = type_
	return type_
}

// Java method: RestrictByFalseTypeOfResultVisitor.caseVoidType()
func (this_ *ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitorBase) CaseVoidType() any {
	return this_.outer_.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_VOID_TYPE())
}

// Java class: com.google.javascript.jscomp.type.ChainableReverseAbstractInterpreter.RestrictByOneTypeOfResultVisitor
type ChainableReverseAbstractInterpreter_RestrictByOneTypeOfResultVisitor struct {
	ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitorBase
	outer_            ChainableReverseAbstractInterpreter
	value             string `gson:"value"`
	resultEqualsValue bool   `gson:"resultEqualsValue"`
}

func (this_ *ChainableReverseAbstractInterpreter_RestrictByOneTypeOfResultVisitor) BindSelf_(s any) {
	_ = s
	this_.ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitorBase.BindSelf_(s)
}
func (this_ *ChainableReverseAbstractInterpreter_RestrictByOneTypeOfResultVisitor) IsJavaNil_() bool {
	return this_ == nil
}
func (*ChainableReverseAbstractInterpreter_RestrictByOneTypeOfResultVisitor) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.type.ChainableReverseAbstractInterpreter$RestrictByOneTypeOfResultVisitor", "RestrictByOneTypeOfResultVisitor"
}
func (this_ *ChainableReverseAbstractInterpreter_RestrictByOneTypeOfResultVisitor) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitorBase.IdentityHash_()
}
func (*ChainableReverseAbstractInterpreter_RestrictByOneTypeOfResultVisitor) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.type.ChainableReverseAbstractInterpreter.RestrictByOneTypeOfResultVisitor", "com.google.javascript.jscomp.type.ChainableReverseAbstractInterpreter$RestrictByOneTypeOfResultVisitor", func(v_ any) bool {
		return j2gort.InstanceOf[*ChainableReverseAbstractInterpreter_RestrictByOneTypeOfResultVisitor](v_)
	})
}
func (this_ *ChainableReverseAbstractInterpreter_RestrictByOneTypeOfResultVisitor) SetOuter_(o ChainableReverseAbstractInterpreter) {
	this_.outer_ = o
}
func (this_ *ChainableReverseAbstractInterpreter_RestrictByOneTypeOfResultVisitor) GetOuter_() ChainableReverseAbstractInterpreter {
	return this_.outer_
}
func (*ChainableReverseAbstractInterpreter_RestrictByOneTypeOfResultVisitor) Implements_jscomptype_ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitor() {
}
func (*ChainableReverseAbstractInterpreter_RestrictByOneTypeOfResultVisitor) Implements_rhino_Visitor() {
}

func (this_ *ChainableReverseAbstractInterpreter_RestrictByOneTypeOfResultVisitor) initFields_ChainableReverseAbstractInterpreter_RestrictByOneTypeOfResultVisitor() {
	_ = this_
	this_.value = j2gort.NullString
}

// Java constructor: RestrictByOneTypeOfResultVisitor(java.lang.String value, boolean resultEqualsValue)
func NewChainableReverseAbstractInterpreter_RestrictByOneTypeOfResultVisitor(outer_ ChainableReverseAbstractInterpreter, value string, resultEqualsValue bool) *ChainableReverseAbstractInterpreter_RestrictByOneTypeOfResultVisitor {
	this_ := &ChainableReverseAbstractInterpreter_RestrictByOneTypeOfResultVisitor{}
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	InitChainableReverseAbstractInterpreter_RestrictByOneTypeOfResultVisitor(this_, value, resultEqualsValue)
	return this_
}

func InitChainableReverseAbstractInterpreter_RestrictByOneTypeOfResultVisitor(this_ *ChainableReverseAbstractInterpreter_RestrictByOneTypeOfResultVisitor, value string, resultEqualsValue bool) {
	_ = value
	_ = resultEqualsValue
	this_.ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitorBase.outer_ = this_.outer_
	InitChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitor(&this_.ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitorBase)
	this_.initFields_ChainableReverseAbstractInterpreter_RestrictByOneTypeOfResultVisitor()
	this_.value = value
	this_.resultEqualsValue = resultEqualsValue
}

// Java method: RestrictByOneTypeOfResultVisitor.matchesExpectation(java.lang.String result)
func (this_ *ChainableReverseAbstractInterpreter_RestrictByOneTypeOfResultVisitor) priv_RestrictByOneTypeOfResultVisitor_matchesExpectation_7848(result string) bool {
	_ = result
	return (func() bool { sr_2 := result; se_3 := this_.value; return j2gort.CheckNotNullString(sr_2) == se_3 }() == this_.resultEqualsValue)
}

// Java method: RestrictByOneTypeOfResultVisitor.caseTopType(com.google.javascript.rhino.jstype.JSType topType)
func (this_ *ChainableReverseAbstractInterpreter_RestrictByOneTypeOfResultVisitor) CaseTopType(topType pkg_rhino.JSType) pkg_rhino.JSType {
	_ = topType
	var result pkg_rhino.JSType = topType
	_ = result
	if this_.resultEqualsValue {
		var typeByName pkg_rhino.JSType = this_.outer_.GetChainableReverseAbstractInterpreterBase().priv_ChainableReverseAbstractInterpreter_getNativeTypeForTypeOf_7827(this_.value)
		_ = typeByName
		if typeByName != nil && !typeByName.IsJavaNil_() {
			result = typeByName
		}
	}
	return result
}

// Java method: RestrictByOneTypeOfResultVisitor.caseNoObjectType()
func (this_ *ChainableReverseAbstractInterpreter_RestrictByOneTypeOfResultVisitor) CaseNoObjectType() any {
	if ((j2gort.CheckNotNullString(this_.value) == "object") || (j2gort.CheckNotNullString(this_.value) == "function")) == this_.resultEqualsValue {
		return this_.outer_.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_NO_OBJECT_TYPE())
	} else {
		return nil
	}
}

// Java method: RestrictByOneTypeOfResultVisitor.caseBooleanType()
func (this_ *ChainableReverseAbstractInterpreter_RestrictByOneTypeOfResultVisitor) CaseBooleanType() any {
	if this_.priv_RestrictByOneTypeOfResultVisitor_matchesExpectation_7848("boolean") {
		return this_.outer_.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_BOOLEAN_TYPE())
	} else {
		return nil
	}
}

// Java method: RestrictByOneTypeOfResultVisitor.caseFunctionType(com.google.javascript.rhino.jstype.FunctionType type)
func (this_ *ChainableReverseAbstractInterpreter_RestrictByOneTypeOfResultVisitor) CaseFunctionType(type_ *pkg_rhino.FunctionType) any {
	_ = type_
	if this_.priv_RestrictByOneTypeOfResultVisitor_matchesExpectation_7848("function") {
		return pkg_rhino.ToIface_JSType(type_)
	} else {
		return nil
	}
}

// Java method: RestrictByOneTypeOfResultVisitor.caseNullType()
func (this_ *ChainableReverseAbstractInterpreter_RestrictByOneTypeOfResultVisitor) CaseNullType() any {
	if this_.priv_RestrictByOneTypeOfResultVisitor_matchesExpectation_7848("object") {
		return this_.outer_.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_NULL_TYPE())
	} else {
		return nil
	}
}

// Java method: RestrictByOneTypeOfResultVisitor.caseNumberType()
func (this_ *ChainableReverseAbstractInterpreter_RestrictByOneTypeOfResultVisitor) CaseNumberType() any {
	if this_.priv_RestrictByOneTypeOfResultVisitor_matchesExpectation_7848("number") {
		return this_.outer_.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_NUMBER_TYPE())
	} else {
		return nil
	}
}

// Java method: RestrictByOneTypeOfResultVisitor.caseBigIntType()
func (this_ *ChainableReverseAbstractInterpreter_RestrictByOneTypeOfResultVisitor) CaseBigIntType() any {
	if this_.priv_RestrictByOneTypeOfResultVisitor_matchesExpectation_7848("bigint") {
		return this_.outer_.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_BIGINT_TYPE())
	} else {
		return nil
	}
}

// Java method: RestrictByOneTypeOfResultVisitor.caseObjectType(com.google.javascript.rhino.jstype.ObjectType type)
func (this_ *ChainableReverseAbstractInterpreter_RestrictByOneTypeOfResultVisitor) CaseObjectType(type_ pkg_rhino.ObjectType) any {
	_ = type_
	if j2gort.CheckNotNullString(this_.value) == "function" {
		var ctorType pkg_rhino.JSType = this_.outer_.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_FUNCTION_TYPE())
		_ = ctorType
		if this_.resultEqualsValue {
			return ctorType.GetJSTypeBase().GetGreatestSubtype(pkg_rhino.Upcast_JSType(type_))
		} else {
			if type_.GetObjectTypeBase().IsSubtypeOf_1(ctorType) {
				return nil
			} else {
				return pkg_rhino.Upcast_JSType(type_)
			}
		}
	}
	if this_.priv_RestrictByOneTypeOfResultVisitor_matchesExpectation_7848("object") {
		return pkg_rhino.Upcast_JSType(type_)
	} else {
		return nil
	}
}

// Java method: RestrictByOneTypeOfResultVisitor.caseStringType()
func (this_ *ChainableReverseAbstractInterpreter_RestrictByOneTypeOfResultVisitor) CaseStringType() any {
	if this_.priv_RestrictByOneTypeOfResultVisitor_matchesExpectation_7848("string") {
		return this_.outer_.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_STRING_TYPE())
	} else {
		return nil
	}
}

// Java method: RestrictByOneTypeOfResultVisitor.caseSymbolType(com.google.javascript.rhino.jstype.SymbolType type)
func (this_ *ChainableReverseAbstractInterpreter_RestrictByOneTypeOfResultVisitor) CaseSymbolType(type_ pkg_rhino.SymbolType) any {
	_ = type_
	if this_.priv_RestrictByOneTypeOfResultVisitor_matchesExpectation_7848("symbol") {
		return pkg_rhino.Upcast_JSType(type_)
	} else {
		return nil
	}
}

// Java method: RestrictByOneTypeOfResultVisitor.caseVoidType()
func (this_ *ChainableReverseAbstractInterpreter_RestrictByOneTypeOfResultVisitor) CaseVoidType() any {
	if this_.priv_RestrictByOneTypeOfResultVisitor_matchesExpectation_7848("undefined") {
		return this_.outer_.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_VOID_TYPE())
	} else {
		return nil
	}
}

// Java class:
type ChainableReverseAbstractInterpreter__Anon_129 struct {
	ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitorBase
	outer_ ChainableReverseAbstractInterpreter
}

func (this_ *ChainableReverseAbstractInterpreter__Anon_129) BindSelf_(s any) {
	_ = s
	this_.ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitorBase.BindSelf_(s)
}
func (this_ *ChainableReverseAbstractInterpreter__Anon_129) IsJavaNil_() bool { return this_ == nil }
func (*ChainableReverseAbstractInterpreter__Anon_129) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.type.ChainableReverseAbstractInterpreter$1", ""
}
func (this_ *ChainableReverseAbstractInterpreter__Anon_129) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitorBase.IdentityHash_()
}
func (*ChainableReverseAbstractInterpreter__Anon_129) HasIdentityEquals_() {}
func (this_ *ChainableReverseAbstractInterpreter__Anon_129) SetOuter_(o ChainableReverseAbstractInterpreter) {
	this_.outer_ = o
}
func (this_ *ChainableReverseAbstractInterpreter__Anon_129) GetOuter_() ChainableReverseAbstractInterpreter {
	return this_.outer_
}
func (*ChainableReverseAbstractInterpreter__Anon_129) Implements_jscomptype_ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitor() {
}
func (*ChainableReverseAbstractInterpreter__Anon_129) Implements_jscomptype_ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitor() {
}
func (*ChainableReverseAbstractInterpreter__Anon_129) Implements_rhino_Visitor() {}

func (this_ *ChainableReverseAbstractInterpreter__Anon_129) initFields_ChainableReverseAbstractInterpreter__Anon_129() {
	_ = this_
}

func NewChainableReverseAbstractInterpreter__Anon_129(outer_ ChainableReverseAbstractInterpreter, superOuter_ ChainableReverseAbstractInterpreter) *ChainableReverseAbstractInterpreter__Anon_129 {
	this_ := &ChainableReverseAbstractInterpreter__Anon_129{}
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	this_.ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitorBase.outer_ = superOuter_
	InitChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitor(&this_.ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitorBase)
	this_.initFields_ChainableReverseAbstractInterpreter__Anon_129()
	return this_
}

// Java method: .caseTopType(com.google.javascript.rhino.jstype.JSType topType)
func (this_ *ChainableReverseAbstractInterpreter__Anon_129) CaseTopType(topType pkg_rhino.JSType) pkg_rhino.JSType {
	_ = topType
	if topType.GetJSTypeBase().IsAllType() {
		return this_.outer_.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_ARRAY_TYPE())
	} else {
		return topType
	}
}

// Java method: .caseObjectType(com.google.javascript.rhino.jstype.ObjectType type)
func (this_ *ChainableReverseAbstractInterpreter__Anon_129) CaseObjectType(type_ pkg_rhino.ObjectType) any {
	_ = type_
	var arrayType pkg_rhino.JSType = this_.outer_.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_ARRAY_TYPE())
	_ = arrayType
	if arrayType.GetJSTypeBase().IsSubtypeOf_1(pkg_rhino.Upcast_JSType(type_)) {
		return arrayType
	} else {
		return nil
	}
}

// Java class:
type ChainableReverseAbstractInterpreter__Anon_130 struct {
	ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitorBase
	outer_ ChainableReverseAbstractInterpreter
}

func (this_ *ChainableReverseAbstractInterpreter__Anon_130) BindSelf_(s any) {
	_ = s
	this_.ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitorBase.BindSelf_(s)
}
func (this_ *ChainableReverseAbstractInterpreter__Anon_130) IsJavaNil_() bool { return this_ == nil }
func (*ChainableReverseAbstractInterpreter__Anon_130) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.type.ChainableReverseAbstractInterpreter$2", ""
}
func (this_ *ChainableReverseAbstractInterpreter__Anon_130) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitorBase.IdentityHash_()
}
func (*ChainableReverseAbstractInterpreter__Anon_130) HasIdentityEquals_() {}
func (this_ *ChainableReverseAbstractInterpreter__Anon_130) SetOuter_(o ChainableReverseAbstractInterpreter) {
	this_.outer_ = o
}
func (this_ *ChainableReverseAbstractInterpreter__Anon_130) GetOuter_() ChainableReverseAbstractInterpreter {
	return this_.outer_
}
func (*ChainableReverseAbstractInterpreter__Anon_130) Implements_jscomptype_ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitor() {
}
func (*ChainableReverseAbstractInterpreter__Anon_130) Implements_jscomptype_ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitor() {
}
func (*ChainableReverseAbstractInterpreter__Anon_130) Implements_rhino_Visitor() {}

func (this_ *ChainableReverseAbstractInterpreter__Anon_130) initFields_ChainableReverseAbstractInterpreter__Anon_130() {
	_ = this_
}

func NewChainableReverseAbstractInterpreter__Anon_130(outer_ ChainableReverseAbstractInterpreter, superOuter_ ChainableReverseAbstractInterpreter) *ChainableReverseAbstractInterpreter__Anon_130 {
	this_ := &ChainableReverseAbstractInterpreter__Anon_130{}
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	this_.ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitorBase.outer_ = superOuter_
	InitChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitor(&this_.ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitorBase)
	this_.initFields_ChainableReverseAbstractInterpreter__Anon_130()
	return this_
}

// Java method: .caseObjectType(com.google.javascript.rhino.jstype.ObjectType type)
func (this_ *ChainableReverseAbstractInterpreter__Anon_130) CaseObjectType(type_ pkg_rhino.ObjectType) any {
	_ = type_
	if type_.GetObjectTypeBase().IsSubtypeOf_1(this_.outer_.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_ARRAY_TYPE())) {
		return nil
	} else {
		return pkg_rhino.Upcast_JSType(type_)
	}
}
