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
// Port of com/google/javascript/rhino/jstype/JSTypeResolver.java (generated mechanically by j2go).

package rhino

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.rhino.jstype.JSTypeResolver
type JSTypeResolver struct {
	registry        *JSTypeRegistry      `gson:"registry"`
	captureStack    j2gort.Deque[JSType] `gson:"captureStack"`
	resolutionQueue j2gort.Deque[JSType] `gson:"resolutionQueue"`
	idHash_         atomic.Int32
	state           JSTypeResolver_State `gson:"state"`
}

func (this_ *JSTypeResolver) BindSelf_(s any) {
	_ = s
}
func (this_ *JSTypeResolver) IsJavaNil_() bool { return this_ == nil }
func (*JSTypeResolver) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.JSTypeResolver", "JSTypeResolver"
}
func (this_ *JSTypeResolver) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*JSTypeResolver) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.JSTypeResolver", "com.google.javascript.rhino.jstype.JSTypeResolver", func(v_ any) bool {
		return j2gort.InstanceOf[*JSTypeResolver](v_)
	})
}

func (this_ *JSTypeResolver) initFields_JSTypeResolver() {
	_ = this_
	this_.captureStack = j2gort.Deque[JSType]{}
	this_.resolutionQueue = j2gort.Deque[JSType]{}
	this_.state = StaticGet_JSTypeResolver_State_CLOSED()
}

// Java constructor: JSTypeResolver(com.google.javascript.rhino.jstype.JSTypeRegistry registry)
func NewJSTypeResolver(registry *JSTypeRegistry) *JSTypeResolver {
	this_ := &JSTypeResolver{}
	InitJSTypeResolver(this_, registry)
	return this_
}

func InitJSTypeResolver(this_ *JSTypeResolver, registry *JSTypeRegistry) {
	_ = registry
	this_.initFields_JSTypeResolver()
	this_.registry = registry
}

// Java static method: JSTypeResolver.create(com.google.javascript.rhino.jstype.JSTypeRegistry registry)
func JSTypeResolver_Create(registry *JSTypeRegistry) *JSTypeResolver {
	_ = registry
	if registry.GetResolver() != nil {
		return registry.GetResolver()
	}
	return NewJSTypeResolver(registry)
}

// Java method: JSTypeResolver.addUnresolved(com.google.javascript.rhino.jstype.JSType captured)
func (this_ *JSTypeResolver) AddUnresolved(captured JSType) {
	_ = captured
	if !(!(captured.GetJSTypeBase().IsResolved())) {
		j2gort.CheckState(false)
	}
	this_.captureStack.AddLast(captured)
}

// Java method: JSTypeResolver.resolveIfClosed(com.google.javascript.rhino.jstype.JSType captured, com.google.javascript.rhino.jstype.JSTypeClass caller)
func (this_ *JSTypeResolver) ResolveIfClosed(captured JSType, caller JSTypeClass) {
	_ = captured
	_ = caller
	if !(caller.IsTypeOf(captured)) {
		return
	}
	var expected JSType = this_.captureStack.RemoveLast()
	_ = expected
	if !(captured == expected || ((captured == nil || captured.IsJavaNil_()) && (expected == nil || expected.IsJavaNil_()))) {
		j2gort.CheckStateFmt(false, "Captured %s; Expected %s", captured, expected)
	}
	switch (this_.state).Ordinal() {
	case int32(0) /* State.CLOSED */, int32(2) /* State.CLOSING */ :
		this_.priv_JSTypeResolver_doResolve_9287(captured)
	case int32(1) /* State.OPEN */ :
		this_.resolutionQueue.AddLast(captured)
	}
}

// Java method: JSTypeResolver.openForDefinition()
func (this_ *JSTypeResolver) OpenForDefinition() *JSTypeResolver_Closer {
	{
		var iifeCond_1 bool
		{

			re_1 := this_.state
			ae_2 := StaticGet_JSTypeResolver_State_CLOSED()
			if re_1 == 0 {
				j2gort.ThrowNPE()
			}
			iifeCond_1 = re_1 == ae_2
		}
		if !(iifeCond_1) {
			j2gort.CheckState(false)
		}
	}
	if !(this_.captureStack.IsEmpty()) {
		j2gort.CheckState(false)
	}
	this_.state = StaticGet_JSTypeResolver_State_OPEN()
	return NewJSTypeResolver_Closer(this_)
}

// Java method: JSTypeResolver.resolveAll()
func (this_ *JSTypeResolver) priv_JSTypeResolver_resolveAll_9285() {
	{
		var iifeCond_1 bool
		{

			re_3 := this_.state
			ae_4 := StaticGet_JSTypeResolver_State_OPEN()
			if re_3 == 0 {
				j2gort.ThrowNPE()
			}
			iifeCond_1 = re_3 == ae_4
		}
		if !(iifeCond_1) {
			j2gort.CheckState(false)
		}
	}
	if !(this_.captureStack.IsEmpty()) {
		j2gort.CheckState(false)
	}
	this_.state = StaticGet_JSTypeResolver_State_CLOSING()
	for !(this_.resolutionQueue.IsEmpty()) {
		this_.priv_JSTypeResolver_doResolve_9287(this_.resolutionQueue.RemoveFirst())
	}
	this_.resolutionQueue = j2gort.Deque[JSType]{}
	var globalThis PrototypeObjectType = j2gort.AsIface[PrototypeObjectType](this_.registry.GetNativeType(StaticGet_JSTypeNative_GLOBAL_THIS()))
	_ = globalThis
	var windowType JSType = this_.registry.GetGlobalType("Window")
	_ = windowType
	base_globalThis_ := globalThis.GetPrototypeObjectTypeBase()
	if base_globalThis_.IsUnknownType() {
		var windowObjType ObjectType = ObjectType_Cast(windowType)
		_ = windowObjType
		if windowObjType != nil && !windowObjType.IsJavaNil_() {
			base_globalThis_.SetImplicitPrototype(windowObjType)
		} else {
			base_globalThis_.SetImplicitPrototype(this_.registry.GetNativeObjectType(StaticGet_JSTypeNative_OBJECT_TYPE()))
		}
	}
	this_.state = StaticGet_JSTypeResolver_State_CLOSED()
}

// Java method: JSTypeResolver.isClosed()
func (this_ *JSTypeResolver) IsClosed() bool {
	return (this_.state == StaticGet_JSTypeResolver_State_CLOSED())
}

// Java method: JSTypeResolver.doResolve(com.google.javascript.rhino.jstype.JSType type)
func (this_ *JSTypeResolver) priv_JSTypeResolver_doResolve_9287(type_ JSType) {
	_ = type_
	type_.GetJSTypeBase().Resolve(this_.registry.GetErrorReporter())
}

// Java method: JSTypeResolver.assertLegalToResolveTypes()
func (this_ *JSTypeResolver) AssertLegalToResolveTypes() {
	{
		var iifeCond_1 bool
		{

			re_5 := this_.state
			ae_6 := StaticGet_JSTypeResolver_State_OPEN()
			if re_5 == 0 {
				j2gort.ThrowNPE()
			}
			iifeCond_1 = re_5 == ae_6
		}
		if !(!(iifeCond_1)) {
			j2gort.CheckStateStr(false, "Types cannot be resolved while the registry is open")
		}
	}
}

// Java class: com.google.javascript.rhino.jstype.JSTypeResolver.Closer
type JSTypeResolver_Closer struct {
	outer_  *JSTypeResolver
	idHash_ atomic.Int32
	hasRun  bool `gson:"hasRun"`
}

func (this_ *JSTypeResolver_Closer) BindSelf_(s any) {
	_ = s
}
func (this_ *JSTypeResolver_Closer) IsJavaNil_() bool { return this_ == nil }
func (*JSTypeResolver_Closer) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.JSTypeResolver$Closer", "Closer"
}
func (this_ *JSTypeResolver_Closer) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*JSTypeResolver_Closer) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.JSTypeResolver.Closer", "com.google.javascript.rhino.jstype.JSTypeResolver$Closer", func(v_ any) bool {
		return j2gort.InstanceOf[*JSTypeResolver_Closer](v_)
	})
}
func (this_ *JSTypeResolver_Closer) SetOuter_(o *JSTypeResolver)              { this_.outer_ = o }
func (this_ *JSTypeResolver_Closer) GetOuter_() *JSTypeResolver               { return this_.outer_ }
func (*JSTypeResolver_Closer) Implements_j2gort_Ext_java_lang_AutoCloseable() {}

func (this_ *JSTypeResolver_Closer) initFields_JSTypeResolver_Closer() {
	_ = this_

}

// Java constructor: Closer()
func NewJSTypeResolver_Closer(outer_ *JSTypeResolver) *JSTypeResolver_Closer {
	this_ := &JSTypeResolver_Closer{}
	this_.outer_ = outer_
	InitJSTypeResolver_Closer(this_)
	return this_
}

func InitJSTypeResolver_Closer(this_ *JSTypeResolver_Closer) {
	this_.initFields_JSTypeResolver_Closer()
}

// Java method: Closer.close()
func (this_ *JSTypeResolver_Closer) Close_0() {
	if !(!(this_.hasRun)) {
		j2gort.CheckState(false)
	}
	this_.hasRun = true
	this_.outer_.priv_JSTypeResolver_resolveAll_9285()
}

// Java enum: com.google.javascript.rhino.jstype.JSTypeResolver.State
type JSTypeResolver_State uint8

const (
	enum_JSTypeResolver_State_CLOSED  JSTypeResolver_State = 1
	enum_JSTypeResolver_State_OPEN    JSTypeResolver_State = 2
	enum_JSTypeResolver_State_CLOSING JSTypeResolver_State = 3
)

var enumBaseTable_JSTypeResolver_State = [...]j2gort.EnumBase{
	{Name_: "CLOSED", Ordinal_: 0},
	{Name_: "OPEN", Ordinal_: 1},
	{Name_: "CLOSING", Ordinal_: 2},
}
var enumValues_JSTypeResolver_State = []JSTypeResolver_State{1, 2, 3}

func StaticGet_JSTypeResolver_State_CLOSED() JSTypeResolver_State {
	return enum_JSTypeResolver_State_CLOSED
}
func StaticGet_JSTypeResolver_State_OPEN() JSTypeResolver_State {
	return enum_JSTypeResolver_State_OPEN
}
func StaticGet_JSTypeResolver_State_CLOSING() JSTypeResolver_State {
	return enum_JSTypeResolver_State_CLOSING
}

func JSTypeResolver_State_Values() []JSTypeResolver_State {
	return j2gort.CloneSlice(enumValues_JSTypeResolver_State)
}

func JSTypeResolver_State_ValueOf(name string) JSTypeResolver_State {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_JSTypeResolver_State {
		if enumBaseTable_JSTypeResolver_State[i_].Name_ == name {
			return JSTypeResolver_State(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.rhino.jstype.JSTypeResolver.State." + name))
}

func (this_ JSTypeResolver_State) IsJavaNil_() bool                              { return this_ == 0 }
func (JSTypeResolver_State) HasIdentityEquals_()                                 {}
func (JSTypeResolver_State) Implements_j2gort_Ext_java_lang_Enum()               {}
func (JSTypeResolver_State) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ JSTypeResolver_State) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (JSTypeResolver_State) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.JSTypeResolver$State", "State"
}
func (this_ JSTypeResolver_State) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_JSTypeResolver_State[this_-1]
}
func (this_ JSTypeResolver_State) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ JSTypeResolver_State) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ JSTypeResolver_State) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_JSTypeResolver_State[this_-1].Name_
}
func (this_ JSTypeResolver_State) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_JSTypeResolver_State[this_-1].Name_
}
func (this_ JSTypeResolver_State) Equals(other_ any) bool {
	o_, ok_ := other_.(JSTypeResolver_State)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ JSTypeResolver_State) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_JSTypeResolver_State[this_-1].HashCode()
}
func (this_ JSTypeResolver_State) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_JSTypeResolver_State[this_-1].CompareTo(other_)
}
func (this_ JSTypeResolver_State) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_JSTypeResolver_State[this_-1].Name_)
}
func init() {
	vals_ := []any{JSTypeResolver_State(1), JSTypeResolver_State(2), JSTypeResolver_State(3)}
	j2gort.RegisterEnumConstants("com.google.javascript.rhino.jstype.JSTypeResolver.State", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.rhino.jstype.JSTypeResolver$State", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.JSTypeResolver.State", "com.google.javascript.rhino.jstype.JSTypeResolver$State", func(v_ any) bool {
		return j2gort.InstanceOf[JSTypeResolver_State](v_)
	})
}
