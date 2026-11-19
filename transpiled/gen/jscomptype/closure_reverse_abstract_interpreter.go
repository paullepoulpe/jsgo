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
// Port of com/google/javascript/jscomp/type/ClosureReverseAbstractInterpreter.java (generated mechanically by j2go).

package jscomptype

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.type.ClosureReverseAbstractInterpreter
type ClosureReverseAbstractInterpreter struct {
	ChainableReverseAbstractInterpreterBase
	restrictToObjectVisitor    pkg_rhino.Visitor                                 `gson:"restrictToObjectVisitor"`
	restrictToNotObjectVisitor pkg_rhino.Visitor                                 `gson:"restrictToNotObjectVisitor"`
	restricters                j2gort.Ext_com_google_common_collect_ImmutableMap `gson:"restricters"`
}

func (this_ *ClosureReverseAbstractInterpreter) BindSelf_(s any) {
	_ = s
	this_.ChainableReverseAbstractInterpreterBase.BindSelf_(s)
}
func (this_ *ClosureReverseAbstractInterpreter) IsJavaNil_() bool { return this_ == nil }
func (*ClosureReverseAbstractInterpreter) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.type.ClosureReverseAbstractInterpreter", "ClosureReverseAbstractInterpreter"
}
func (this_ *ClosureReverseAbstractInterpreter) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.ChainableReverseAbstractInterpreterBase.IdentityHash_()
}
func (*ClosureReverseAbstractInterpreter) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.type.ClosureReverseAbstractInterpreter", "com.google.javascript.jscomp.type.ClosureReverseAbstractInterpreter", func(v_ any) bool {
		return j2gort.InstanceOf[*ClosureReverseAbstractInterpreter](v_)
	})
}
func (*ClosureReverseAbstractInterpreter) Implements_jscomptype_ChainableReverseAbstractInterpreter() {
}
func (*ClosureReverseAbstractInterpreter) Implements_jscomptype_ReverseAbstractInterpreter() {}

func (this_ *ClosureReverseAbstractInterpreter) initFields_ClosureReverseAbstractInterpreter() {
	_ = this_
	this_.restrictToObjectVisitor = NewClosureReverseAbstractInterpreter__Anon_131(this_, this_)
	this_.restrictToNotObjectVisitor = NewClosureReverseAbstractInterpreter__Anon_132(this_, this_)
}

// Java constructor: ClosureReverseAbstractInterpreter(com.google.javascript.rhino.jstype.JSTypeRegistry typeRegistry)
func NewClosureReverseAbstractInterpreter(typeRegistry *pkg_rhino.JSTypeRegistry) *ClosureReverseAbstractInterpreter {
	this_ := &ClosureReverseAbstractInterpreter{}
	this_.BindSelf_(this_)
	InitClosureReverseAbstractInterpreter(this_, typeRegistry)
	return this_
}

func InitClosureReverseAbstractInterpreter(this_ *ClosureReverseAbstractInterpreter, typeRegistry *pkg_rhino.JSTypeRegistry) {
	_ = typeRegistry
	InitChainableReverseAbstractInterpreter(&this_.ChainableReverseAbstractInterpreterBase, typeRegistry)
	this_.initFields_ClosureReverseAbstractInterpreter()
	this_.restricters = j2gort.Ext_com_google_common_collect_ImmutableMap_Of_2("isObject", j2gort.NewLambda_Ext_com_google_common_base_Function(func(lp_1 any) any {
		p := j2gort.CastPtr[ClosureReverseAbstractInterpreter_TypeRestriction](lp_1)
		_ = p
		if p.type_ == nil || p.type_.IsJavaNil_() {
			if p.outcome {
				return this_.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_OBJECT_TYPE())
			} else {
				return nil
			}
		}
		var visitor pkg_rhino.Visitor
		_ = visitor
		if p.outcome {
			visitor = this_.restrictToObjectVisitor
		} else {
			visitor = this_.restrictToNotObjectVisitor
		}
		return j2gort.AsIface[pkg_rhino.JSType](p.type_.Visit_1(visitor))
	}))
}

// Java method: ClosureReverseAbstractInterpreter.getPreciserScopeKnowingConditionOutcome(com.google.javascript.rhino.Node condition, com.google.javascript.jscomp.type.FlowScope blindScope, com.google.javascript.rhino.Outcome outcome)
func (this_ *ClosureReverseAbstractInterpreter) GetPreciserScopeKnowingConditionOutcome(condition pkg_rhino.Node, blindScope FlowScope, outcome *pkg_rhino.Outcome) FlowScope {
	_ = condition
	_ = blindScope
	_ = outcome
	if condition.IsCall() && condition.HasTwoChildren() {
		var callee pkg_rhino.Node = condition.GetFirstChild()
		_ = callee
		var param pkg_rhino.Node = condition.GetLastChild()
		_ = param
		if callee.IsGetProp() && param.IsQualifiedName() {
			var paramType pkg_rhino.JSType = this_.GetTypeIfRefinable(param, blindScope)
			_ = paramType
			var receiver pkg_rhino.Node = callee.GetFirstChild()
			_ = receiver
			if receiver.IsName() && func() bool {
				sr_2 := "goog"
				se_3 := receiver.GetString()
				return j2gort.CheckNotNullString(sr_2) == se_3
			}() {
				var restricter j2gort.Ext_com_google_common_base_Function = j2gort.AsIface[j2gort.Ext_com_google_common_base_Function](j2gort.MapGetStr(this_.restricters, callee.GetString()))
				_ = restricter
				if restricter != nil {
					return this_.priv_ClosureReverseAbstractInterpreter_restrictParameter_7849(param, paramType, blindScope, restricter, outcome.IsTruthy())
				}
			}
		}
	}
	return this_.NextPreciserScopeKnowingConditionOutcome(condition, blindScope, outcome)
}

// Java method: ClosureReverseAbstractInterpreter.restrictParameter(com.google.javascript.rhino.Node parameter, com.google.javascript.rhino.jstype.JSType type, com.google.javascript.jscomp.type.FlowScope blindScope, com.google.common.base.Function<com.google.javascript.jscomp.type.ClosureReverseAbstractInterpreter.TypeRestriction,com.google.javascript.rhino.jstype.JSType> restriction, boolean outcome)
func (this_ *ClosureReverseAbstractInterpreter) priv_ClosureReverseAbstractInterpreter_restrictParameter_7849(parameter pkg_rhino.Node, type_ pkg_rhino.JSType, blindScope FlowScope, restriction j2gort.Ext_com_google_common_base_Function, outcome bool) FlowScope {
	_ = parameter
	_ = type_
	_ = blindScope
	_ = restriction
	_ = outcome
	type_ = j2gort.AsIface[pkg_rhino.JSType](restriction.Apply(NewClosureReverseAbstractInterpreter_TypeRestriction(type_, outcome)))
	if type_ != nil && !type_.IsJavaNil_() {
		return this_.DeclareNameInScope(blindScope, parameter, type_)
	} else {
		return blindScope
	}
}

// Java class: com.google.javascript.jscomp.type.ClosureReverseAbstractInterpreter.TypeRestriction
type ClosureReverseAbstractInterpreter_TypeRestriction struct {
	type_   pkg_rhino.JSType `gson:"type"`
	idHash_ atomic.Int32
	outcome bool `gson:"outcome"`
}

func (this_ *ClosureReverseAbstractInterpreter_TypeRestriction) BindSelf_(s any) {
	_ = s
}
func (this_ *ClosureReverseAbstractInterpreter_TypeRestriction) IsJavaNil_() bool {
	return this_ == nil
}
func (*ClosureReverseAbstractInterpreter_TypeRestriction) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.type.ClosureReverseAbstractInterpreter$TypeRestriction", "TypeRestriction"
}
func (this_ *ClosureReverseAbstractInterpreter_TypeRestriction) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ClosureReverseAbstractInterpreter_TypeRestriction) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.type.ClosureReverseAbstractInterpreter.TypeRestriction", "com.google.javascript.jscomp.type.ClosureReverseAbstractInterpreter$TypeRestriction", func(v_ any) bool {
		return j2gort.InstanceOf[*ClosureReverseAbstractInterpreter_TypeRestriction](v_)
	})
}

func (this_ *ClosureReverseAbstractInterpreter_TypeRestriction) initFields_ClosureReverseAbstractInterpreter_TypeRestriction() {
	_ = this_
}

// Java constructor: TypeRestriction(com.google.javascript.rhino.jstype.JSType type, boolean outcome)
func NewClosureReverseAbstractInterpreter_TypeRestriction(type_ pkg_rhino.JSType, outcome bool) *ClosureReverseAbstractInterpreter_TypeRestriction {
	this_ := &ClosureReverseAbstractInterpreter_TypeRestriction{}
	InitClosureReverseAbstractInterpreter_TypeRestriction(this_, type_, outcome)
	return this_
}

func InitClosureReverseAbstractInterpreter_TypeRestriction(this_ *ClosureReverseAbstractInterpreter_TypeRestriction, type_ pkg_rhino.JSType, outcome bool) {
	_ = type_
	_ = outcome
	this_.initFields_ClosureReverseAbstractInterpreter_TypeRestriction()
	this_.type_ = type_
	this_.outcome = outcome
}

// Java class:
type ClosureReverseAbstractInterpreter__Anon_131 struct {
	ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitorBase
	outer_ *ClosureReverseAbstractInterpreter
}

func (this_ *ClosureReverseAbstractInterpreter__Anon_131) BindSelf_(s any) {
	_ = s
	this_.ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitorBase.BindSelf_(s)
}
func (this_ *ClosureReverseAbstractInterpreter__Anon_131) IsJavaNil_() bool { return this_ == nil }
func (*ClosureReverseAbstractInterpreter__Anon_131) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.type.ClosureReverseAbstractInterpreter$1", ""
}
func (this_ *ClosureReverseAbstractInterpreter__Anon_131) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitorBase.IdentityHash_()
}
func (*ClosureReverseAbstractInterpreter__Anon_131) HasIdentityEquals_() {}
func (this_ *ClosureReverseAbstractInterpreter__Anon_131) SetOuter_(o *ClosureReverseAbstractInterpreter) {
	this_.outer_ = o
}
func (this_ *ClosureReverseAbstractInterpreter__Anon_131) GetOuter_() *ClosureReverseAbstractInterpreter {
	return this_.outer_
}
func (*ClosureReverseAbstractInterpreter__Anon_131) Implements_jscomptype_ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitor() {
}
func (*ClosureReverseAbstractInterpreter__Anon_131) Implements_jscomptype_ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitor() {
}
func (*ClosureReverseAbstractInterpreter__Anon_131) Implements_rhino_Visitor() {}

func (this_ *ClosureReverseAbstractInterpreter__Anon_131) initFields_ClosureReverseAbstractInterpreter__Anon_131() {
	_ = this_
}

func NewClosureReverseAbstractInterpreter__Anon_131(outer_ *ClosureReverseAbstractInterpreter, superOuter_ ChainableReverseAbstractInterpreter) *ClosureReverseAbstractInterpreter__Anon_131 {
	this_ := &ClosureReverseAbstractInterpreter__Anon_131{}
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	this_.ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitorBase.outer_ = superOuter_
	InitChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitor(&this_.ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitorBase)
	this_.initFields_ClosureReverseAbstractInterpreter__Anon_131()
	return this_
}

// Java method: .caseTopType(com.google.javascript.rhino.jstype.JSType topType)
func (this_ *ClosureReverseAbstractInterpreter__Anon_131) CaseTopType(topType pkg_rhino.JSType) pkg_rhino.JSType {
	_ = topType
	return this_.outer_.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_NO_OBJECT_TYPE())
}

// Java method: .caseObjectType(com.google.javascript.rhino.jstype.ObjectType type)
func (this_ *ClosureReverseAbstractInterpreter__Anon_131) CaseObjectType(type_ pkg_rhino.ObjectType) any {
	_ = type_
	return type_
}

// Java method: .caseFunctionType(com.google.javascript.rhino.jstype.FunctionType type)
func (this_ *ClosureReverseAbstractInterpreter__Anon_131) CaseFunctionType(type_ *pkg_rhino.FunctionType) any {
	_ = type_
	return j2gort.PtrToAny(type_)
}

// Java class:
type ClosureReverseAbstractInterpreter__Anon_132 struct {
	ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitorBase
	outer_ *ClosureReverseAbstractInterpreter
}

func (this_ *ClosureReverseAbstractInterpreter__Anon_132) BindSelf_(s any) {
	_ = s
	this_.ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitorBase.BindSelf_(s)
}
func (this_ *ClosureReverseAbstractInterpreter__Anon_132) IsJavaNil_() bool { return this_ == nil }
func (*ClosureReverseAbstractInterpreter__Anon_132) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.type.ClosureReverseAbstractInterpreter$2", ""
}
func (this_ *ClosureReverseAbstractInterpreter__Anon_132) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitorBase.IdentityHash_()
}
func (*ClosureReverseAbstractInterpreter__Anon_132) HasIdentityEquals_() {}
func (this_ *ClosureReverseAbstractInterpreter__Anon_132) SetOuter_(o *ClosureReverseAbstractInterpreter) {
	this_.outer_ = o
}
func (this_ *ClosureReverseAbstractInterpreter__Anon_132) GetOuter_() *ClosureReverseAbstractInterpreter {
	return this_.outer_
}
func (*ClosureReverseAbstractInterpreter__Anon_132) Implements_jscomptype_ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitor() {
}
func (*ClosureReverseAbstractInterpreter__Anon_132) Implements_jscomptype_ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitor() {
}
func (*ClosureReverseAbstractInterpreter__Anon_132) Implements_rhino_Visitor() {}

func (this_ *ClosureReverseAbstractInterpreter__Anon_132) initFields_ClosureReverseAbstractInterpreter__Anon_132() {
	_ = this_
}

func NewClosureReverseAbstractInterpreter__Anon_132(outer_ *ClosureReverseAbstractInterpreter, superOuter_ ChainableReverseAbstractInterpreter) *ClosureReverseAbstractInterpreter__Anon_132 {
	this_ := &ClosureReverseAbstractInterpreter__Anon_132{}
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	this_.ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitorBase.outer_ = superOuter_
	InitChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitor(&this_.ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitorBase)
	this_.initFields_ClosureReverseAbstractInterpreter__Anon_132()
	return this_
}

// Java method: .caseObjectType(com.google.javascript.rhino.jstype.ObjectType type)
func (this_ *ClosureReverseAbstractInterpreter__Anon_132) CaseObjectType(type_ pkg_rhino.ObjectType) any {
	_ = type_
	return nil
}

// Java method: .caseFunctionType(com.google.javascript.rhino.jstype.FunctionType type)
func (this_ *ClosureReverseAbstractInterpreter__Anon_132) CaseFunctionType(type_ *pkg_rhino.FunctionType) any {
	_ = type_
	return nil
}
