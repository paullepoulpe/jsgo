/*
 * Copyright 2008 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/LinkedFlowScope.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_jscomptype "github.com/paullepoulpe/jsgo/transpiled/gen/jscomptype"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_linked_flow_scope_1 j2gort.Ext_java_util_function_BiPredicate
var lambdaSingleton_linked_flow_scope_2 j2gort.Ext_java_util_function_BiPredicate

func init() {
	lambdaSingleton_linked_flow_scope_1 = j2gort.NewLambda_Ext_java_util_function_BiPredicate(func(mrp0_ any, mrp1_ any) bool {
		return LinkedFlowScope_EqualScopes(j2gort.CastPtr[LinkedFlowScope_OverlayScope](mrp0_), j2gort.CastPtr[LinkedFlowScope_OverlayScope](mrp1_))
	})
	lambdaSingleton_linked_flow_scope_2 = j2gort.NewLambda_Ext_java_util_function_BiPredicate(func(mrp0_ any, mrp1_ any) bool {
		return LinkedFlowScope_EqualSlots(j2gort.AsIface[pkg_rhino.StaticTypedSlot](mrp0_), j2gort.AsIface[pkg_rhino.StaticTypedSlot](mrp1_))
	})
}

// Java class: com.google.javascript.jscomp.LinkedFlowScope
type LinkedFlowScope struct {
	pool               *LinkedFlowScope_Pool                                           `gson:"pool"`
	inputProvider      CompilerInputProvider                                           `gson:"inputProvider"`
	scopes             *pkg_rhino.HamtPMap[*TypedScope, *LinkedFlowScope_OverlayScope] `gson:"scopes"`
	singleScopeKey     *TypedScope                                                     `gson:"singleScopeKey"`
	singleScopeVal     *LinkedFlowScope_OverlayScope                                   `gson:"singleScopeVal"`
	functionScope      *TypedScope                                                     `gson:"functionScope"`
	syntacticScope     *TypedScope                                                     `gson:"syntacticScope"`
	unindexedRootBloom int64                                                           `gson:"unindexedRootBloom"`
}

func (this_ *LinkedFlowScope) BindSelf_(s any) {
	_ = s
}
func (this_ *LinkedFlowScope) IsJavaNil_() bool { return this_ == nil }
func (*LinkedFlowScope) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.LinkedFlowScope", "LinkedFlowScope"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.LinkedFlowScope", "com.google.javascript.jscomp.LinkedFlowScope", func(v_ any) bool {
		return j2gort.InstanceOf[*LinkedFlowScope](v_)
	})
}
func (*LinkedFlowScope) Implements_jscomptype_FlowScope()   {}
func (*LinkedFlowScope) Implements_rhino_StaticTypedScope() {}
func (*LinkedFlowScope) Implements_graph_LatticeElement()   {}
func (*LinkedFlowScope) Implements_rhino_StaticScope()      {}
func (*LinkedFlowScope) Implements_graph_Annotation()       {}

var static_LinkedFlowScope_EQUAL_SCOPES_PREDICATE j2gort.Ext_java_util_function_BiPredicate

func StaticGet_LinkedFlowScope_EQUAL_SCOPES_PREDICATE() j2gort.Ext_java_util_function_BiPredicate {
	clinit_LinkedFlowScope()
	return static_LinkedFlowScope_EQUAL_SCOPES_PREDICATE
}
func StaticSet_LinkedFlowScope_EQUAL_SCOPES_PREDICATE(v j2gort.Ext_java_util_function_BiPredicate) {
	clinit_LinkedFlowScope()
	static_LinkedFlowScope_EQUAL_SCOPES_PREDICATE = v
}

var static_LinkedFlowScope_EQUAL_SLOTS_PREDICATE j2gort.Ext_java_util_function_BiPredicate

func StaticGet_LinkedFlowScope_EQUAL_SLOTS_PREDICATE() j2gort.Ext_java_util_function_BiPredicate {
	clinit_LinkedFlowScope()
	return static_LinkedFlowScope_EQUAL_SLOTS_PREDICATE
}
func StaticSet_LinkedFlowScope_EQUAL_SLOTS_PREDICATE(v j2gort.Ext_java_util_function_BiPredicate) {
	clinit_LinkedFlowScope()
	static_LinkedFlowScope_EQUAL_SLOTS_PREDICATE = v
}

var static_LinkedFlowScope_EMPTY_SCOPES *pkg_rhino.HamtPMap[*TypedScope, *LinkedFlowScope_OverlayScope]

func StaticGet_LinkedFlowScope_EMPTY_SCOPES() *pkg_rhino.HamtPMap[*TypedScope, *LinkedFlowScope_OverlayScope] {
	clinit_LinkedFlowScope()
	return static_LinkedFlowScope_EMPTY_SCOPES
}
func StaticSet_LinkedFlowScope_EMPTY_SCOPES(v *pkg_rhino.HamtPMap[*TypedScope, *LinkedFlowScope_OverlayScope]) {
	clinit_LinkedFlowScope()
	static_LinkedFlowScope_EMPTY_SCOPES = v
}

var static_LinkedFlowScope_EMPTY_SLOTS *pkg_rhino.HamtPMap[string, *LinkedFlowScope_OverlaySlot]

func StaticGet_LinkedFlowScope_EMPTY_SLOTS() *pkg_rhino.HamtPMap[string, *LinkedFlowScope_OverlaySlot] {
	clinit_LinkedFlowScope()
	return static_LinkedFlowScope_EMPTY_SLOTS
}
func StaticSet_LinkedFlowScope_EMPTY_SLOTS(v *pkg_rhino.HamtPMap[string, *LinkedFlowScope_OverlaySlot]) {
	clinit_LinkedFlowScope()
	static_LinkedFlowScope_EMPTY_SLOTS = v
}

var clinitDone_LinkedFlowScope uint32
var clinitState_LinkedFlowScope j2gort.ClinitState

func clinit_LinkedFlowScope() {
	if atomic.LoadUint32(&clinitDone_LinkedFlowScope) != 0 {
		return
	}
	clinitSlow_LinkedFlowScope()
}

func clinitSlow_LinkedFlowScope() {
	if !j2gort.EnterClinit(&clinitState_LinkedFlowScope, &clinitDone_LinkedFlowScope) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_LinkedFlowScope, &clinitDone_LinkedFlowScope)
	static_LinkedFlowScope_EQUAL_SCOPES_PREDICATE = lambdaSingleton_linked_flow_scope_1
	static_LinkedFlowScope_EQUAL_SLOTS_PREDICATE = lambdaSingleton_linked_flow_scope_2
	static_LinkedFlowScope_EMPTY_SCOPES = pkg_rhino.HamtPMap_Empty[*TypedScope, *LinkedFlowScope_OverlayScope]()
	static_LinkedFlowScope_EMPTY_SLOTS = pkg_rhino.HamtPMap_Empty[string, *LinkedFlowScope_OverlaySlot]()
}

func (this_ *LinkedFlowScope) initFields_LinkedFlowScope() {
	_ = this_
}

// Java constructor: LinkedFlowScope()
func NewLinkedFlowScope() *LinkedFlowScope {
	clinit_LinkedFlowScope()
	this_ := &LinkedFlowScope{}
	InitLinkedFlowScope(this_)
	return this_
}

func InitLinkedFlowScope(this_ *LinkedFlowScope) {
	this_.initFields_LinkedFlowScope()
}

// Java static method: LinkedFlowScope.revHash(java.lang.String s)
func LinkedFlowScope_RevHash(s string) int32 {
	clinit_LinkedFlowScope()
	_ = s
	return j2gort.BoxInt_Reverse(javalang.StringHashCode(j2gort.CheckNotNullString(s)))
}

// Java static method: LinkedFlowScope.compareUnsigned(int left, int right)
func LinkedFlowScope_CompareUnsigned(left int32, right int32) int32 {
	clinit_LinkedFlowScope()
	_ = left
	_ = right
	var diff int32 = ((javalang.Ushr32(left, int(int32(2)))) - (javalang.Ushr32(right, int(int32(2)))))
	_ = diff
	if diff != int32(0) {
		return diff
	} else {
		return ((left & int32(3)) - (right & int32(3)))
	}
}

// Java static method: LinkedFlowScope.createFlowScope(com.google.javascript.jscomp.LinkedFlowScope.@org.jspecify.annotations.Nullable Pool pool, com.google.javascript.jscomp.CompilerInputProvider inputProvider, com.google.javascript.rhino.@org.jspecify.annotations.Nullable HamtPMap<com.google.javascript.jscomp.TypedScope,com.google.javascript.jscomp.LinkedFlowScope.OverlayScope> scopes, com.google.javascript.jscomp.@org.jspecify.annotations.Nullable TypedScope singleScopeKey, com.google.javascript.jscomp.LinkedFlowScope.@org.jspecify.annotations.Nullable OverlayScope singleScopeVal, com.google.javascript.jscomp.TypedScope syntacticScope, com.google.javascript.jscomp.TypedScope functionScope, long unindexedRootBloom)
func LinkedFlowScope_CreateFlowScope(pool *LinkedFlowScope_Pool, inputProvider CompilerInputProvider, scopes *pkg_rhino.HamtPMap[*TypedScope, *LinkedFlowScope_OverlayScope], singleScopeKey *TypedScope, singleScopeVal *LinkedFlowScope_OverlayScope, syntacticScope *TypedScope, functionScope *TypedScope, unindexedRootBloom int64) *LinkedFlowScope {
	clinit_LinkedFlowScope()
	_ = pool
	_ = inputProvider
	_ = scopes
	_ = singleScopeKey
	_ = singleScopeVal
	_ = syntacticScope
	_ = functionScope
	_ = unindexedRootBloom
	if pool != nil {
		return pool.AcquireFlowScope(inputProvider, scopes, singleScopeKey, singleScopeVal, syntacticScope, functionScope, unindexedRootBloom)
	}
	clinit_LinkedFlowScope()
	var val_1 LinkedFlowScope
	InitLinkedFlowScope(&val_1)
	var s *LinkedFlowScope = &val_1
	_ = s
	s.priv_LinkedFlowScope_init_3736(nil, inputProvider, scopes, singleScopeKey, singleScopeVal, syntacticScope, functionScope, unindexedRootBloom)
	return s
}

// Java static method: LinkedFlowScope.createFlowScopeFromMap(com.google.javascript.jscomp.LinkedFlowScope.@org.jspecify.annotations.Nullable Pool pool, com.google.javascript.jscomp.CompilerInputProvider inputProvider, com.google.javascript.rhino.@org.jspecify.annotations.Nullable HamtPMap<com.google.javascript.jscomp.TypedScope,com.google.javascript.jscomp.LinkedFlowScope.OverlayScope> map, com.google.javascript.jscomp.TypedScope syntacticScope, com.google.javascript.jscomp.TypedScope functionScope, long unindexedRootBloom)
func LinkedFlowScope_CreateFlowScopeFromMap(pool *LinkedFlowScope_Pool, inputProvider CompilerInputProvider, map_ *pkg_rhino.HamtPMap[*TypedScope, *LinkedFlowScope_OverlayScope], syntacticScope *TypedScope, functionScope *TypedScope, unindexedRootBloom int64) *LinkedFlowScope {
	clinit_LinkedFlowScope()
	_ = pool
	_ = inputProvider
	_ = map_
	_ = syntacticScope
	_ = functionScope
	_ = unindexedRootBloom
	if (map_ == nil) || map_.IsEmpty() {
		return LinkedFlowScope_CreateFlowScope(pool, inputProvider, nil, nil, nil, syntacticScope, functionScope, int64(0))
	}
	if map_.IsSingleEntry() {
		return LinkedFlowScope_CreateFlowScope(pool, inputProvider, nil, map_.RootKey(), map_.RootValue(), syntacticScope, functionScope, unindexedRootBloom)
	}
	return LinkedFlowScope_CreateFlowScope(pool, inputProvider, map_, nil, nil, syntacticScope, functionScope, unindexedRootBloom)
}

// Java static method: LinkedFlowScope.createOverlayScope(com.google.javascript.jscomp.LinkedFlowScope.@org.jspecify.annotations.Nullable Pool pool, com.google.javascript.jscomp.TypedScope scope, com.google.javascript.rhino.@org.jspecify.annotations.Nullable HamtPMap<java.lang.String,com.google.javascript.jscomp.LinkedFlowScope.OverlaySlot> slots, com.google.javascript.jscomp.LinkedFlowScope.@org.jspecify.annotations.Nullable OverlaySlot slot0, com.google.javascript.jscomp.LinkedFlowScope.@org.jspecify.annotations.Nullable OverlaySlot slot1, boolean hasQualifiedSlots)
func LinkedFlowScope_CreateOverlayScope(pool *LinkedFlowScope_Pool, scope *TypedScope, slots *pkg_rhino.HamtPMap[string, *LinkedFlowScope_OverlaySlot], slot0 *LinkedFlowScope_OverlaySlot, slot1 *LinkedFlowScope_OverlaySlot, hasQualifiedSlots bool) *LinkedFlowScope_OverlayScope {
	clinit_LinkedFlowScope()
	_ = pool
	_ = scope
	_ = slots
	_ = slot0
	_ = slot1
	_ = hasQualifiedSlots
	if pool != nil {
		return pool.AcquireOverlayScope(scope, slots, slot0, slot1, hasQualifiedSlots)
	}
	var val_2 LinkedFlowScope_OverlayScope
	InitLinkedFlowScope_OverlayScope(&val_2)
	var s *LinkedFlowScope_OverlayScope = &val_2
	_ = s
	s.Init(nil, scope, slots, slot0, slot1, hasQualifiedSlots)
	return s
}

// Java static method: LinkedFlowScope.createOverlayScopeFromMap(com.google.javascript.jscomp.LinkedFlowScope.@org.jspecify.annotations.Nullable Pool pool, com.google.javascript.jscomp.TypedScope scope, com.google.javascript.rhino.@org.jspecify.annotations.Nullable HamtPMap<java.lang.String,com.google.javascript.jscomp.LinkedFlowScope.OverlaySlot> map, boolean hasQualifiedSlots)
func LinkedFlowScope_CreateOverlayScopeFromMap(pool *LinkedFlowScope_Pool, scope *TypedScope, map_ *pkg_rhino.HamtPMap[string, *LinkedFlowScope_OverlaySlot], hasQualifiedSlots bool) *LinkedFlowScope_OverlayScope {
	clinit_LinkedFlowScope()
	_ = pool
	_ = scope
	_ = map_
	_ = hasQualifiedSlots
	if (map_ == nil) || map_.IsEmpty() {
		return LinkedFlowScope_CreateOverlayScope(pool, scope, nil, nil, nil, false)
	}
	if map_.IsSingleEntry() {
		return LinkedFlowScope_CreateOverlayScope(pool, scope, nil, map_.RootValue(), nil, hasQualifiedSlots)
	}
	return LinkedFlowScope_CreateOverlayScope(pool, scope, map_, nil, nil, hasQualifiedSlots)
}

// Java static method: LinkedFlowScope.createOverlaySlot(com.google.javascript.jscomp.LinkedFlowScope.@org.jspecify.annotations.Nullable Pool pool, java.lang.String name, com.google.javascript.rhino.jstype.JSType type)
func LinkedFlowScope_CreateOverlaySlot(pool *LinkedFlowScope_Pool, name string, type_ pkg_rhino.JSType) *LinkedFlowScope_OverlaySlot {
	clinit_LinkedFlowScope()
	_ = pool
	_ = name
	_ = type_
	if pool != nil {
		return pool.AcquireOverlaySlot(name, type_)
	}
	var val_3 LinkedFlowScope_OverlaySlot
	InitLinkedFlowScope_OverlaySlot(&val_3)
	var s *LinkedFlowScope_OverlaySlot = &val_3
	_ = s
	s.Init(name, type_)
	return s
}

// Java method: LinkedFlowScope.init(com.google.javascript.jscomp.LinkedFlowScope.@org.jspecify.annotations.Nullable Pool pool, com.google.javascript.jscomp.CompilerInputProvider inputProvider, com.google.javascript.rhino.@org.jspecify.annotations.Nullable HamtPMap<com.google.javascript.jscomp.TypedScope,com.google.javascript.jscomp.LinkedFlowScope.OverlayScope> scopes, com.google.javascript.jscomp.@org.jspecify.annotations.Nullable TypedScope singleScopeKey, com.google.javascript.jscomp.LinkedFlowScope.@org.jspecify.annotations.Nullable OverlayScope singleScopeVal, com.google.javascript.jscomp.TypedScope syntacticScope, com.google.javascript.jscomp.TypedScope functionScope, long unindexedRootBloom)
func (this_ *LinkedFlowScope) priv_LinkedFlowScope_init_3736(pool *LinkedFlowScope_Pool, inputProvider CompilerInputProvider, scopes *pkg_rhino.HamtPMap[*TypedScope, *LinkedFlowScope_OverlayScope], singleScopeKey *TypedScope, singleScopeVal *LinkedFlowScope_OverlayScope, syntacticScope *TypedScope, functionScope *TypedScope, unindexedRootBloom int64) {
	_ = pool
	_ = inputProvider
	_ = scopes
	_ = singleScopeKey
	_ = singleScopeVal
	_ = syntacticScope
	_ = functionScope
	_ = unindexedRootBloom
	this_.pool = pool
	this_.inputProvider = inputProvider
	this_.scopes = scopes
	this_.singleScopeKey = singleScopeKey
	this_.singleScopeVal = singleScopeVal
	this_.syntacticScope = syntacticScope
	this_.functionScope = functionScope
	this_.unindexedRootBloom = unindexedRootBloom
}

// Java method: LinkedFlowScope.clear()
func (this_ *LinkedFlowScope) priv_LinkedFlowScope_clear_3737() {
	this_.pool = nil
	this_.inputProvider = nil
	this_.scopes = nil
	this_.singleScopeKey = nil
	this_.singleScopeVal = nil
	this_.syntacticScope = nil
	this_.functionScope = nil
	this_.unindexedRootBloom = int64(0)
}

// Java method: LinkedFlowScope.hasNoScopes()
func (this_ *LinkedFlowScope) priv_LinkedFlowScope_hasNoScopes_3738() bool {
	return ((this_.scopes == nil) && (this_.singleScopeKey == nil))
}

// Java method: LinkedFlowScope.sameScopes(com.google.javascript.jscomp.LinkedFlowScope other)
func (this_ *LinkedFlowScope) priv_LinkedFlowScope_sameScopes_3739(other *LinkedFlowScope) bool {
	_ = other
	return (((this_.scopes == other.scopes) && (this_.singleScopeKey == other.singleScopeKey)) && (this_.singleScopeVal == other.singleScopeVal))
}

// Java method: LinkedFlowScope.getScopeOverlay(com.google.javascript.jscomp.TypedScope key)
func (this_ *LinkedFlowScope) priv_LinkedFlowScope_getScopeOverlay_3740(key *TypedScope) *LinkedFlowScope_OverlayScope {
	_ = key
	if this_.scopes != nil {
		return this_.scopes.Get_1(key)
	}
	if this_.singleScopeKey == key {
		return this_.singleScopeVal
	} else {
		return nil
	}
}

// Java method: LinkedFlowScope.trimsAwaySingleKey(com.google.javascript.jscomp.TypedScope scope, com.google.javascript.jscomp.TypedScope keyToCheck)
func (this_ *LinkedFlowScope) priv_LinkedFlowScope_trimsAwaySingleKey_3741(scope *TypedScope, keyToCheck *TypedScope) bool {
	_ = scope
	_ = keyToCheck
	var thisScope *TypedScope = this_.syntacticScope
	_ = thisScope
	var thatScope *TypedScope = scope
	_ = thatScope
	var thisDepth int32 = thisScope.GetDepth()
	_ = thisDepth
	var thatDepth int32 = thatScope.GetDepth()
	_ = thatDepth
	for thatDepth > thisDepth {
		thatScope = j2gort.CastPtr[TypedScope](thatScope.GetParent())
		thatDepth = thatDepth - 1
	}
	for thisDepth > thatDepth {
		if thisScope == keyToCheck {
			return true
		}
		thisScope = j2gort.CastPtr[TypedScope](thisScope.GetParent())
		thisDepth = thisDepth - 1
	}
	for ((thisScope != thatScope) && (thisScope != nil)) && (thatScope != nil) {
		if thisScope == keyToCheck {
			return true
		}
		thisScope = j2gort.CastPtr[TypedScope](thisScope.GetParent())
		thatScope = j2gort.CastPtr[TypedScope](thatScope.GetParent())
	}
	return false
}

// Java static method: LinkedFlowScope.trimScopesMap(com.google.javascript.rhino.HamtPMap<com.google.javascript.jscomp.TypedScope,com.google.javascript.jscomp.LinkedFlowScope.OverlayScope> startMap, com.google.javascript.jscomp.TypedScope fromScope, com.google.javascript.jscomp.TypedScope toScope)
func LinkedFlowScope_TrimScopesMap(startMap *pkg_rhino.HamtPMap[*TypedScope, *LinkedFlowScope_OverlayScope], fromScope *TypedScope, toScope *TypedScope) *pkg_rhino.HamtPMap[*TypedScope, *LinkedFlowScope_OverlayScope] {
	clinit_LinkedFlowScope()
	_ = startMap
	_ = fromScope
	_ = toScope
	var thisScope *TypedScope = fromScope
	_ = thisScope
	var thatScope *TypedScope = toScope
	_ = thatScope
	var thisDepth int32 = thisScope.GetDepth()
	_ = thisDepth
	var thatDepth int32 = thatScope.GetDepth()
	_ = thatDepth
	var result *pkg_rhino.HamtPMap[*TypedScope, *LinkedFlowScope_OverlayScope] = startMap
	_ = result
	for thatDepth > thisDepth {
		thatScope = j2gort.CastPtr[TypedScope](thatScope.GetParent())
		thatDepth = thatDepth - 1
	}
	for (thisDepth > thatDepth) && (result != nil) {
		result = result.MinusOrNull_1(thisScope)
		thisScope = j2gort.CastPtr[TypedScope](thisScope.GetParent())
		thisDepth = thisDepth - 1
	}
	for (((thisScope != thatScope) && (thisScope != nil)) && (thatScope != nil)) && (result != nil) {
		result = result.MinusOrNull_1(thisScope)
		thisScope = j2gort.CastPtr[TypedScope](thisScope.GetParent())
		thatScope = j2gort.CastPtr[TypedScope](thatScope.GetParent())
	}
	return result
}

// Java method: LinkedFlowScope.flowsFromBottom()
func (this_ *LinkedFlowScope) priv_LinkedFlowScope_flowsFromBottom_3742() bool {
	return this_.functionScope.IsBottom()
}

// Java static method: LinkedFlowScope.createEntryLattice(com.google.javascript.jscomp.CompilerInputProvider inputProvider, com.google.javascript.jscomp.TypedScope scope)
func LinkedFlowScope_CreateEntryLattice_2(inputProvider CompilerInputProvider, scope *TypedScope) *LinkedFlowScope {
	clinit_LinkedFlowScope()
	_ = inputProvider
	_ = scope
	return LinkedFlowScope_CreateFlowScope(nil, inputProvider, nil, nil, nil, scope, scope, int64(0))
}

// Java static method: LinkedFlowScope.createEntryLattice(com.google.javascript.jscomp.CompilerInputProvider inputProvider, com.google.javascript.jscomp.TypedScope scope, com.google.javascript.jscomp.LinkedFlowScope.@org.jspecify.annotations.Nullable Pool pool)
func LinkedFlowScope_CreateEntryLattice_3(inputProvider CompilerInputProvider, scope *TypedScope, pool *LinkedFlowScope_Pool) *LinkedFlowScope {
	clinit_LinkedFlowScope()
	_ = inputProvider
	_ = scope
	_ = pool
	return LinkedFlowScope_CreateFlowScope(pool, inputProvider, nil, nil, nil, scope, scope, int64(0))
}

// Java static method: LinkedFlowScope.hashPrefixBeforeDot(java.lang.String symbol, int dot)
func LinkedFlowScope_HashPrefixBeforeDot(symbol string, dot int32) int32 {
	clinit_LinkedFlowScope()
	_ = symbol
	_ = dot
	var h int32 = int32(0)
	_ = h
	{
		var i int32 = int32(0)
		_ = i
		for ; i < dot; i = i + 1 {
			h = ((int32(31) * h) + int32(j2gort.JStrCharAt(symbol, i)))
		}
	}
	return h
}

// Java method: LinkedFlowScope.inferSlotType(java.lang.String symbol, com.google.javascript.rhino.jstype.JSType type)
func (this_ *LinkedFlowScope) InferSlotType(symbol string, type_ pkg_rhino.JSType) any {
	_ = symbol
	_ = type_
	var rootVar *TypedVar = this_.syntacticScope.LookupTypedVar(LinkedFlowScope_GetRootOfQualifiedName(symbol))
	_ = rootVar
	var targetScope *TypedScope
	_ = targetScope
	if rootVar != nil {
		targetScope = j2gort.CastPtr[TypedScope](rootVar.GetScope())
	} else {
		targetScope = nil
	}
	if targetScope == nil {
		targetScope = this_.functionScope
	}
	var existingScope *LinkedFlowScope_OverlayScope = this_.priv_LinkedFlowScope_getScopeOverlay_3740(targetScope)
	_ = existingScope
	var newScope *LinkedFlowScope_OverlayScope = nil
	_ = newScope
	var dot int32 = j2gort.JStrIndexOfCp(symbol, int32(uint16(46)))
	_ = dot
	if existingScope == nil {
		var slot *LinkedFlowScope_OverlaySlot = LinkedFlowScope_CreateOverlaySlot(this_.pool, symbol, type_)
		_ = slot
		newScope = LinkedFlowScope_CreateOverlayScope(this_.pool, targetScope, nil, slot, nil, (dot >= int32(0)))
	} else {
		newScope = existingScope.Infer(symbol, type_, (dot >= int32(0)))
		if newScope == existingScope {
			return this_
		}
	}
	var nextBloom int64 = this_.unindexedRootBloom
	_ = nextBloom
	if !(newScope.IsEmpty()) && (dot >= int32(0)) {
		nextBloom = (nextBloom | (javalang.Shl64(int64(1), int((LinkedFlowScope_HashPrefixBeforeDot(symbol, dot) & int32(63))))))
	}
	if this_.scopes == nil {
		if this_.singleScopeKey == nil {
			if newScope.IsEmpty() {
				return this_
			}
			return j2gort.PtrToAny(LinkedFlowScope_CreateFlowScope(this_.pool, this_.inputProvider, nil, targetScope, newScope, this_.syntacticScope, this_.functionScope, nextBloom))
		}
		if this_.singleScopeKey == targetScope {
			if newScope.IsEmpty() {
				return j2gort.PtrToAny(LinkedFlowScope_CreateFlowScope(this_.pool, this_.inputProvider, nil, nil, nil, this_.syntacticScope, this_.functionScope, int64(0)))
			}
			return j2gort.PtrToAny(LinkedFlowScope_CreateFlowScope(this_.pool, this_.inputProvider, nil, targetScope, newScope, this_.syntacticScope, this_.functionScope, nextBloom))
		}
		if newScope.IsEmpty() {
			return this_
		}
		var twoScopes *pkg_rhino.HamtPMap[*TypedScope, *LinkedFlowScope_OverlayScope] = pkg_rhino.HamtPMap_OfTwo[*TypedScope, *LinkedFlowScope_OverlayScope](this_.singleScopeKey, this_.singleScopeVal, targetScope, newScope)
		_ = twoScopes
		return j2gort.PtrToAny(LinkedFlowScope_CreateFlowScope(this_.pool, this_.inputProvider, twoScopes, nil, nil, this_.syntacticScope, this_.functionScope, nextBloom))
	}
	var newScopes *pkg_rhino.HamtPMap[*TypedScope, *LinkedFlowScope_OverlayScope]
	_ = newScopes
	if !(newScope.IsEmpty()) {
		newScopes = j2gort.CastPtr[pkg_rhino.HamtPMap[*TypedScope, *LinkedFlowScope_OverlayScope]](this_.scopes.Plus_2(targetScope, newScope))
	} else {
		newScopes = this_.scopes.MinusOrNull_1(targetScope)
	}
	if newScopes == this_.scopes {
		return this_
	}
	return j2gort.PtrToAny(LinkedFlowScope_CreateFlowScopeFromMap(this_.pool, this_.inputProvider, newScopes, this_.syntacticScope, this_.functionScope, nextBloom))
}

// Java method: LinkedFlowScope.inferQualifiedSlot(com.google.javascript.rhino.Node node, java.lang.String symbol, com.google.javascript.rhino.jstype.JSType bottomType, com.google.javascript.rhino.jstype.JSType inferredType, boolean declared)
func (this_ *LinkedFlowScope) InferQualifiedSlot(node pkg_rhino.Node, symbol string, bottomType pkg_rhino.JSType, inferredType pkg_rhino.JSType, declared bool) any {
	_ = node
	_ = symbol
	_ = bottomType
	_ = inferredType
	_ = declared
	if this_.functionScope.IsGlobal() {
		return this_
	}
	var v *TypedVar = this_.syntacticScope.LookupTypedVar(symbol)
	_ = v
	if (v == nil) && !(this_.functionScope.IsBottom()) {
		var dot int32 = j2gort.JStrIndexOfCp(symbol, int32(uint16(46)))
		_ = dot
		var rootVar *TypedVar
		_ = rootVar
		if dot < int32(0) {
			rootVar = nil
		} else {
			rootVar = this_.syntacticScope.LookupTypedVar(j2gort.JStrSub(symbol, int32(0), dot))
		}
		var rootScope *TypedScope
		_ = rootScope
		if rootVar != nil {
			rootScope = j2gort.CastPtr[TypedScope](rootVar.GetScope())
		} else {
			rootScope = this_.syntacticScope.GetGlobalTypedScope()
		}
		v = rootScope.Declare(symbol, node, bottomType, this_.inputProvider.GetInput(NodeUtil_GetInputId(node)), !(declared))
	}
	var declaredType pkg_rhino.JSType
	_ = declaredType
	if v != nil {
		declaredType = v.GetType()
	} else {
		declaredType = nil
	}
	if v != nil {
		if !(v.IsTypeInferred()) {
			if (((declaredType == nil || declaredType.IsJavaNil_()) || !(inferredType.GetJSTypeBase().IsSubtypeOf_1(declaredType))) || declaredType.GetJSTypeBase().IsSubtypeOf_1(inferredType)) || inferredType.GetJSTypeBase().EqualsJSType(declaredType) {
				return this_
			}
		} else if (declaredType != nil && !declaredType.IsJavaNil_()) && !(inferredType.GetJSTypeBase().IsSubtypeOf_1(declaredType)) {
			v.SetType(v.GetType().GetLeastSupertype(inferredType))
		}
	}
	return j2gort.PtrToAny(j2gort.CastPtr[LinkedFlowScope](this_.InferSlotType(symbol, inferredType)))
}

// Java method: LinkedFlowScope.getTypeOfThis()
func (this_ *LinkedFlowScope) GetTypeOfThis() pkg_rhino.JSType {
	return this_.syntacticScope.GetTypeOfThis()
}

// Java method: LinkedFlowScope.getRootNode()
func (this_ *LinkedFlowScope) GetRootNode_0() pkg_rhino.Node {
	return this_.syntacticScope.GetRootNode_0()
}

// Java method: LinkedFlowScope.getParentScope()
func (this_ *LinkedFlowScope) GetParentScope() any {
	panic(j2gort.AsIface[javalang.IUnsupportedOperationException](j2gort.ConstructShimThrowable("java.lang.UnsupportedOperationException")))
}

// Java method: LinkedFlowScope.getSlot(java.lang.String name)
func (this_ *LinkedFlowScope) GetSlot(name string) any {
	_ = name
	var var_ *TypedVar = this_.syntacticScope.LookupTypedVar(name)
	_ = var_
	if var_ != nil {
		var scope *LinkedFlowScope_OverlayScope = this_.priv_LinkedFlowScope_getScopeOverlay_3740(j2gort.CastPtr[TypedScope](var_.GetScope()))
		_ = scope
		if scope != nil {
			return scope.GetOverlayOrFallback(name, var_)
		} else {
			return pkg_rhino.ToIface_StaticTypedSlot(var_)
		}
	}
	var dot int32 = j2gort.JStrIndexOfCp(name, int32(uint16(46)))
	_ = dot
	var scope_2 *LinkedFlowScope_OverlayScope = nil
	_ = scope_2
	if dot < int32(0) {
		scope_2 = this_.priv_LinkedFlowScope_getScopeOverlay_3740(this_.functionScope)
	} else {
		var rootVar *TypedVar = this_.syntacticScope.LookupTypedVar(j2gort.JStrSub(name, int32(0), dot))
		_ = rootVar
		scope_2 = this_.priv_LinkedFlowScope_getScopeOverlay_3740(func() *TypedScope {
			if rootVar != nil {
				return j2gort.CastPtr[TypedScope](rootVar.GetScope())
			}
			return this_.functionScope
		}())
	}
	if scope_2 != nil {
		return scope_2.GetSlot(name)
	} else {
		return nil
	}
}

// Java method: LinkedFlowScope.hasGlobalFunctionScope()
func (this_ *LinkedFlowScope) HasGlobalFunctionScope() bool {
	return this_.functionScope.IsGlobal()
}

// Java method: LinkedFlowScope.mayHaveQualifiedSlotForRoot(java.lang.String rootName)
func (this_ *LinkedFlowScope) MayHaveQualifiedSlotForRoot(rootName string) bool {
	_ = rootName
	if (this_.unindexedRootBloom != int64(0)) && ((this_.unindexedRootBloom & (javalang.Shl64(int64(1), int((javalang.StringHashCode(j2gort.CheckNotNullString(rootName)) & int32(63)))))) != int64(0)) {
		return true
	}
	return (this_.syntacticScope.MayHaveQualifiedVarForRoot(rootName) || ((this_.functionScope != this_.syntacticScope) && this_.functionScope.MayHaveQualifiedVarForRoot(rootName)))
}

// Java method: LinkedFlowScope.mayHaveQualifiedSlotForNode(com.google.javascript.rhino.Node n)
func (this_ *LinkedFlowScope) MayHaveQualifiedSlotForNode(n pkg_rhino.Node) bool {
	_ = n
	if !(n.IsGetProp()) {
		return n.IsQualifiedName()
	}
	var hashes int64 = n.QualifiedNameHashes()
	_ = hashes
	if hashes == int64(0) {
		return false
	}
	var mixedHash int32 = int32((javalang.Ushr64(hashes, int(int32(32)))))
	_ = mixedHash
	if this_.syntacticScope.MayHaveQualifiedVarWithHash(mixedHash) || ((this_.functionScope != this_.syntacticScope) && this_.functionScope.MayHaveQualifiedVarWithHash(mixedHash)) {
		return true
	}
	if this_.unindexedRootBloom == int64(0) {
		return false
	}
	var rawKeyHash int32 = int32(hashes)
	_ = rawKeyHash
	if this_.scopes != nil {
		return LinkedFlowScope_ContainsOverlaySlotWithHash(this_.scopes, rawKeyHash)
	}
	var single *LinkedFlowScope_OverlayScope = this_.singleScopeVal
	_ = single
	return (((single != nil) && single.hasQualifiedSlots) && single.ContainsSlotWithHash(rawKeyHash))
}

// Java static method: LinkedFlowScope.containsOverlaySlotWithHash(com.google.javascript.rhino.HamtPMap<com.google.javascript.jscomp.TypedScope,com.google.javascript.jscomp.LinkedFlowScope.OverlayScope> map, int rawKeyHash)
func LinkedFlowScope_ContainsOverlaySlotWithHash(map_ *pkg_rhino.HamtPMap[*TypedScope, *LinkedFlowScope_OverlayScope], rawKeyHash int32) bool {
	clinit_LinkedFlowScope()
	_ = map_
	_ = rawKeyHash
	var overlay *LinkedFlowScope_OverlayScope = map_.RootValue()
	_ = overlay
	if ((overlay != nil) && overlay.hasQualifiedSlots) && overlay.ContainsSlotWithHash(rawKeyHash) {
		return true
	}
	{
		var i int32 = int32(0)
		_ = i
		var len_ int32 = map_.ChildCount()
		_ = len_
		for ; i < len_; i = i + 1 {
			if LinkedFlowScope_ContainsOverlaySlotWithHash(map_.ChildAt(i), rawKeyHash) {
				return true
			}
		}
	}
	return false
}

// Java static method: LinkedFlowScope.getRootOfQualifiedName(java.lang.String name)
func LinkedFlowScope_GetRootOfQualifiedName(name string) string {
	clinit_LinkedFlowScope()
	_ = name
	var index int32 = j2gort.JStrIndexOfCp(name, int32(uint16(46)))
	_ = index
	if index < int32(0) {
		return name
	} else {
		return j2gort.JStrSub(name, int32(0), index)
	}
}

// Java method: LinkedFlowScope.getOverlayScopeForName(java.lang.String name)
func (this_ *LinkedFlowScope) priv_LinkedFlowScope_getOverlayScopeForName_3749(name string) *LinkedFlowScope_OverlayScope {
	_ = name
	var rootVar *TypedVar = this_.syntacticScope.LookupTypedVar(LinkedFlowScope_GetRootOfQualifiedName(name))
	_ = rootVar
	var scope *TypedScope
	_ = scope
	if rootVar != nil {
		scope = j2gort.CastPtr[TypedScope](rootVar.GetScope())
	} else {
		scope = nil
	}
	if scope != nil {
		scope = scope
	} else {
		scope = this_.functionScope
	}
	return this_.priv_LinkedFlowScope_getScopeOverlay_3740(scope)
}

// Java method: LinkedFlowScope.getOwnSlot(java.lang.String name)
func (this_ *LinkedFlowScope) GetOwnSlot(name string) any {
	_ = name
	panic(j2gort.AsIface[javalang.IUnsupportedOperationException](j2gort.ConstructShimThrowable("java.lang.UnsupportedOperationException")))
}

// Java method: LinkedFlowScope.withSyntacticScope(com.google.javascript.rhino.jstype.StaticTypedScope scope)
func (this_ *LinkedFlowScope) WithSyntacticScope(scope pkg_rhino.StaticTypedScope) pkg_jscomptype.FlowScope {
	_ = scope
	var typedScope *TypedScope = j2gort.CastPtr[TypedScope](scope)
	_ = typedScope
	if scope == this_.syntacticScope || ((scope == nil || scope.IsJavaNil_()) && (this_.syntacticScope == nil || this_.syntacticScope.IsJavaNil_())) {
		return this_
	}
	if this_.priv_LinkedFlowScope_hasNoScopes_3738() {
		return pkg_jscomptype.ToIface_FlowScope(LinkedFlowScope_CreateFlowScope(this_.pool, this_.inputProvider, nil, nil, nil, typedScope, this_.functionScope, int64(0)))
	}
	if this_.scopes == nil {
		if this_.priv_LinkedFlowScope_trimsAwaySingleKey_3741(typedScope, this_.singleScopeKey) {
			return pkg_jscomptype.ToIface_FlowScope(LinkedFlowScope_CreateFlowScope(this_.pool, this_.inputProvider, nil, nil, nil, typedScope, this_.functionScope, int64(0)))
		}
		return pkg_jscomptype.ToIface_FlowScope(LinkedFlowScope_CreateFlowScope(this_.pool, this_.inputProvider, nil, this_.singleScopeKey, this_.singleScopeVal, typedScope, this_.functionScope, this_.unindexedRootBloom))
	}
	var trimmed *pkg_rhino.HamtPMap[*TypedScope, *LinkedFlowScope_OverlayScope] = LinkedFlowScope_TrimScopesMap(this_.scopes, this_.syntacticScope, typedScope)
	_ = trimmed
	return pkg_jscomptype.ToIface_FlowScope(LinkedFlowScope_CreateFlowScopeFromMap(this_.pool, this_.inputProvider, trimmed, typedScope, this_.functionScope, this_.unindexedRootBloom))
}

// Java method: LinkedFlowScope.getDeclarationScope()
func (this_ *LinkedFlowScope) GetDeclarationScope() any {
	return j2gort.PtrToAny(this_.syntacticScope)
}

// Java static method: LinkedFlowScope.getCommonParentDeclarationScope(com.google.javascript.jscomp.LinkedFlowScope left, com.google.javascript.jscomp.LinkedFlowScope right)
func LinkedFlowScope_GetCommonParentDeclarationScope(left *LinkedFlowScope, right *LinkedFlowScope) *TypedScope {
	clinit_LinkedFlowScope()
	_ = left
	_ = right
	if left.priv_LinkedFlowScope_flowsFromBottom_3742() {
		return right.syntacticScope
	} else if right.priv_LinkedFlowScope_flowsFromBottom_3742() {
		return left.syntacticScope
	}
	return j2gort.CastPtr[TypedScope](left.syntacticScope.GetCommonParent(ToIface_AbstractScope(right.syntacticScope)))
}

// Java method: LinkedFlowScope.equals(java.lang.Object other)
func (this_ *LinkedFlowScope) Equals(arg_4 any) bool {
	other := arg_4
	_ = other
	var that *LinkedFlowScope
	_ = that
	var ok_5 bool
	that, ok_5 = j2gort.TryCast[*LinkedFlowScope](other)
	if !ok_5 {
		return false
	}
	if this_.functionScope != that.functionScope {
		return false
	}
	if (this_.scopes == nil) && (that.scopes == nil) {
		if this_.singleScopeKey != that.singleScopeKey {
			return false
		}
		return ((this_.singleScopeKey == nil) || LinkedFlowScope_EqualScopes(this_.singleScopeVal, that.singleScopeVal))
	}
	if (this_.scopes == nil) || (that.scopes == nil) {
		return false
	}
	return this_.scopes.Equivalent(pkg_rhino.ToIface_PMap[*TypedScope, *LinkedFlowScope_OverlayScope](that.scopes), StaticGet_LinkedFlowScope_EQUAL_SCOPES_PREDICATE())
}

// Java static method: LinkedFlowScope.equalScopes(com.google.javascript.jscomp.LinkedFlowScope.OverlayScope left, com.google.javascript.jscomp.LinkedFlowScope.OverlayScope right)
func LinkedFlowScope_EqualScopes(left *LinkedFlowScope_OverlayScope, right *LinkedFlowScope_OverlayScope) bool {
	clinit_LinkedFlowScope()
	_ = left
	_ = right
	if left == right {
		return true
	}
	if (left.slots == nil) && (right.slots == nil) {
		var l0 *LinkedFlowScope_OverlaySlot = left.slot0
		_ = l0
		var r0 *LinkedFlowScope_OverlaySlot = right.slot0
		_ = r0
		if (l0 == nil) || (r0 == nil) {
			return (l0 == r0)
		}
		var l1 *LinkedFlowScope_OverlaySlot = left.slot1
		_ = l1
		var r1 *LinkedFlowScope_OverlaySlot = right.slot1
		_ = r1
		if (l1 == nil) != (r1 == nil) {
			return false
		}
		{
			var iifeCond_1 bool
			{

				sr_6 := l0.name
				se_7 := r0.name
				iifeCond_1 = j2gort.CheckNotNullString(sr_6) == se_7
			}
			if !(iifeCond_1) || !(LinkedFlowScope_EqualSlots(pkg_rhino.ToIface_StaticTypedSlot(l0), pkg_rhino.ToIface_StaticTypedSlot(r0))) {
				return false
			}
		}
		return ((l1 == nil) || (func() bool { sr_8 := l1.name; se_9 := r1.name; return j2gort.CheckNotNullString(sr_8) == se_9 }() && LinkedFlowScope_EqualSlots(pkg_rhino.ToIface_StaticTypedSlot(l1), pkg_rhino.ToIface_StaticTypedSlot(r1))))
	}
	if (left.slots == nil) || (right.slots == nil) {
		return false
	}
	return left.slots.Equivalent(pkg_rhino.ToIface_PMap[string, *LinkedFlowScope_OverlaySlot](right.slots), StaticGet_LinkedFlowScope_EQUAL_SLOTS_PREDICATE())
}

// Java static method: LinkedFlowScope.equalSlots(com.google.javascript.rhino.jstype.StaticTypedSlot slotA, com.google.javascript.rhino.jstype.StaticTypedSlot slotB)
func LinkedFlowScope_EqualSlots(slotA pkg_rhino.StaticTypedSlot, slotB pkg_rhino.StaticTypedSlot) bool {
	clinit_LinkedFlowScope()
	_ = slotA
	_ = slotB
	return ((slotA == slotB || ((slotA == nil || slotA.IsJavaNil_()) && (slotB == nil || slotB.IsJavaNil_()))) || !(slotA.GetType().GetJSTypeBase().DiffersFrom(slotB.GetType())))
}

// Java method: LinkedFlowScope.hashCode()
func (this_ *LinkedFlowScope) HashCode() int32 {
	panic(j2gort.AsIface[javalang.IUnsupportedOperationException](j2gort.ConstructShimThrowable("java.lang.UnsupportedOperationException")))
}

func (this_ *LinkedFlowScope) LookupQualifiedName_1_QualifiedName(p0_ pkg_rhino.QualifiedName) pkg_rhino.JSType {
	return pkg_rhino.IfaceDefault_StaticTypedScope_LookupQualifiedName_1_QualifiedName(this_, p0_)
}

func (this_ *LinkedFlowScope) LookupQualifiedName_1_string(p0_ string) pkg_rhino.JSType {
	return pkg_rhino.IfaceDefault_StaticTypedScope_LookupQualifiedName_1_string(this_, p0_)
}

func (this_ *LinkedFlowScope) LookupQualifiedNameNode(p0_ pkg_rhino.Node) pkg_rhino.JSType {
	return pkg_rhino.IfaceDefault_StaticTypedScope_LookupQualifiedNameNode(this_, p0_)
}

func (this_ *LinkedFlowScope) GetTopmostScopeOfEventualDeclaration(p0_ string) pkg_rhino.StaticScope {
	return pkg_rhino.IfaceDefault_StaticScope_GetTopmostScopeOfEventualDeclaration(this_, p0_)
}

// Java class: com.google.javascript.jscomp.LinkedFlowScope.Pool
type LinkedFlowScope_Pool struct {
	flowScopes        []*LinkedFlowScope              `gson:"flowScopes"`
	overlayScopes     []*LinkedFlowScope_OverlayScope `gson:"overlayScopes"`
	overlaySlots      []*LinkedFlowScope_OverlaySlot  `gson:"overlaySlots"`
	flowScopeCount    int32                           `gson:"flowScopeCount"`
	flowScopeAlloc    int32                           `gson:"flowScopeAlloc"`
	overlayScopeCount int32                           `gson:"overlayScopeCount"`
	overlayScopeAlloc int32                           `gson:"overlayScopeAlloc"`
	overlaySlotCount  int32                           `gson:"overlaySlotCount"`
	overlaySlotAlloc  int32                           `gson:"overlaySlotAlloc"`
	idHash_           atomic.Int32
}

func (this_ *LinkedFlowScope_Pool) BindSelf_(s any) {
	_ = s
}
func (this_ *LinkedFlowScope_Pool) IsJavaNil_() bool { return this_ == nil }
func (*LinkedFlowScope_Pool) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.LinkedFlowScope$Pool", "Pool"
}
func (this_ *LinkedFlowScope_Pool) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*LinkedFlowScope_Pool) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.LinkedFlowScope.Pool", "com.google.javascript.jscomp.LinkedFlowScope$Pool", func(v_ any) bool {
		return j2gort.InstanceOf[*LinkedFlowScope_Pool](v_)
	})
}

func StaticGet_LinkedFlowScope_Pool_MAX_POOLED() int32 { return int32(4096) }
func (this_ *LinkedFlowScope_Pool) initFields_LinkedFlowScope_Pool() {
	_ = this_
	this_.flowScopes = make([]*LinkedFlowScope, int32(64))
	this_.overlayScopes = make([]*LinkedFlowScope_OverlayScope, int32(64))
	this_.overlaySlots = make([]*LinkedFlowScope_OverlaySlot, int32(64))
}

// Java constructor: Pool()
func NewLinkedFlowScope_Pool() *LinkedFlowScope_Pool {
	this_ := &LinkedFlowScope_Pool{}
	InitLinkedFlowScope_Pool(this_)
	return this_
}

func InitLinkedFlowScope_Pool(this_ *LinkedFlowScope_Pool) {
	this_.initFields_LinkedFlowScope_Pool()
}

// Java method: Pool.acquireFlowScope(com.google.javascript.jscomp.CompilerInputProvider inputProvider, com.google.javascript.rhino.@org.jspecify.annotations.Nullable HamtPMap<com.google.javascript.jscomp.TypedScope,com.google.javascript.jscomp.LinkedFlowScope.OverlayScope> scopes, com.google.javascript.jscomp.@org.jspecify.annotations.Nullable TypedScope singleScopeKey, com.google.javascript.jscomp.LinkedFlowScope.@org.jspecify.annotations.Nullable OverlayScope singleScopeVal, com.google.javascript.jscomp.TypedScope syntacticScope, com.google.javascript.jscomp.TypedScope functionScope, long unindexedRootBloom)
func (this_ *LinkedFlowScope_Pool) AcquireFlowScope(inputProvider CompilerInputProvider, scopes *pkg_rhino.HamtPMap[*TypedScope, *LinkedFlowScope_OverlayScope], singleScopeKey *TypedScope, singleScopeVal *LinkedFlowScope_OverlayScope, syntacticScope *TypedScope, functionScope *TypedScope, unindexedRootBloom int64) *LinkedFlowScope {
	_ = inputProvider
	_ = scopes
	_ = singleScopeKey
	_ = singleScopeVal
	_ = syntacticScope
	_ = functionScope
	_ = unindexedRootBloom
	var s *LinkedFlowScope = nil
	_ = s
	if this_.flowScopeCount < this_.flowScopeAlloc {
		s = this_.flowScopes[func() int32 {
			old_10 := this_.flowScopeCount
			nw_11 := old_10 + 1
			this_.flowScopeCount = nw_11
			return old_10
		}()]
	} else if this_.flowScopeCount < int32(4096) {
		if this_.flowScopeCount == int32(len(this_.flowScopes)) {
			this_.flowScopes = func() []*LinkedFlowScope {
				arrSrc_12 := this_.flowScopes
				arrLen_13 := min(int32(4096), (int32(len(this_.flowScopes)) * int32(2)))
				j2gort.CheckNotNullSlice(arrSrc_12)
				if arrLen_13 < 0 {
					panic(j2gort.ConstructShimThrowable("java.lang.NegativeArraySizeException"))
				}
				arrOut_14 := make([]*LinkedFlowScope, arrLen_13)
				copy(arrOut_14, arrSrc_12)
				return arrOut_14
			}()
		}
		s = NewLinkedFlowScope()
		{
			a_15 := this_.flowScopes
			i_16 := func() int32 {
				old_17 := this_.flowScopeCount
				nw_18 := old_17 + 1
				this_.flowScopeCount = nw_18
				return old_17
			}()
			a_15[i_16] = s
		}
		this_.flowScopeAlloc = this_.flowScopeCount
	} else {
		s = NewLinkedFlowScope()
	}
	s.priv_LinkedFlowScope_init_3736(this_, inputProvider, scopes, singleScopeKey, singleScopeVal, syntacticScope, functionScope, unindexedRootBloom)
	return s
}

// Java method: Pool.acquireOverlayScope(com.google.javascript.jscomp.TypedScope scope, com.google.javascript.rhino.@org.jspecify.annotations.Nullable HamtPMap<java.lang.String,com.google.javascript.jscomp.LinkedFlowScope.OverlaySlot> slots, com.google.javascript.jscomp.LinkedFlowScope.@org.jspecify.annotations.Nullable OverlaySlot slot0, com.google.javascript.jscomp.LinkedFlowScope.@org.jspecify.annotations.Nullable OverlaySlot slot1, boolean hasQualifiedSlots)
func (this_ *LinkedFlowScope_Pool) AcquireOverlayScope(scope *TypedScope, slots *pkg_rhino.HamtPMap[string, *LinkedFlowScope_OverlaySlot], slot0 *LinkedFlowScope_OverlaySlot, slot1 *LinkedFlowScope_OverlaySlot, hasQualifiedSlots bool) *LinkedFlowScope_OverlayScope {
	_ = scope
	_ = slots
	_ = slot0
	_ = slot1
	_ = hasQualifiedSlots
	var s *LinkedFlowScope_OverlayScope = nil
	_ = s
	if this_.overlayScopeCount < this_.overlayScopeAlloc {
		s = this_.overlayScopes[func() int32 {
			old_19 := this_.overlayScopeCount
			nw_20 := old_19 + 1
			this_.overlayScopeCount = nw_20
			return old_19
		}()]
	} else if this_.overlayScopeCount < int32(4096) {
		if this_.overlayScopeCount == int32(len(this_.overlayScopes)) {
			this_.overlayScopes = func() []*LinkedFlowScope_OverlayScope {
				arrSrc_21 := this_.overlayScopes
				arrLen_22 := min(int32(4096), (int32(len(this_.overlayScopes)) * int32(2)))
				j2gort.CheckNotNullSlice(arrSrc_21)
				if arrLen_22 < 0 {
					panic(j2gort.ConstructShimThrowable("java.lang.NegativeArraySizeException"))
				}
				arrOut_23 := make([]*LinkedFlowScope_OverlayScope, arrLen_22)
				copy(arrOut_23, arrSrc_21)
				return arrOut_23
			}()
		}
		s = NewLinkedFlowScope_OverlayScope()
		{
			a_24 := this_.overlayScopes
			i_25 := func() int32 {
				old_26 := this_.overlayScopeCount
				nw_27 := old_26 + 1
				this_.overlayScopeCount = nw_27
				return old_26
			}()
			a_24[i_25] = s
		}
		this_.overlayScopeAlloc = this_.overlayScopeCount
	} else {
		s = NewLinkedFlowScope_OverlayScope()
	}
	s.Init(this_, scope, slots, slot0, slot1, hasQualifiedSlots)
	return s
}

// Java method: Pool.acquireOverlaySlot(java.lang.String name, com.google.javascript.rhino.jstype.JSType type)
func (this_ *LinkedFlowScope_Pool) AcquireOverlaySlot(name string, type_ pkg_rhino.JSType) *LinkedFlowScope_OverlaySlot {
	_ = name
	_ = type_
	var s *LinkedFlowScope_OverlaySlot = nil
	_ = s
	if this_.overlaySlotCount < this_.overlaySlotAlloc {
		s = this_.overlaySlots[func() int32 {
			old_28 := this_.overlaySlotCount
			nw_29 := old_28 + 1
			this_.overlaySlotCount = nw_29
			return old_28
		}()]
	} else if this_.overlaySlotCount < int32(4096) {
		if this_.overlaySlotCount == int32(len(this_.overlaySlots)) {
			this_.overlaySlots = func() []*LinkedFlowScope_OverlaySlot {
				arrSrc_30 := this_.overlaySlots
				arrLen_31 := min(int32(4096), (int32(len(this_.overlaySlots)) * int32(2)))
				j2gort.CheckNotNullSlice(arrSrc_30)
				if arrLen_31 < 0 {
					panic(j2gort.ConstructShimThrowable("java.lang.NegativeArraySizeException"))
				}
				arrOut_32 := make([]*LinkedFlowScope_OverlaySlot, arrLen_31)
				copy(arrOut_32, arrSrc_30)
				return arrOut_32
			}()
		}
		s = NewLinkedFlowScope_OverlaySlot()
		{
			a_33 := this_.overlaySlots
			i_34 := func() int32 {
				old_35 := this_.overlaySlotCount
				nw_36 := old_35 + 1
				this_.overlaySlotCount = nw_36
				return old_35
			}()
			a_33[i_34] = s
		}
		this_.overlaySlotAlloc = this_.overlaySlotCount
	} else {
		s = NewLinkedFlowScope_OverlaySlot()
	}
	s.Init(name, type_)
	return s
}

// Java method: Pool.reset()
func (this_ *LinkedFlowScope_Pool) Reset() {
	var fsLim int32 = min(this_.flowScopeCount, this_.flowScopeAlloc)
	_ = fsLim
	{
		var i int32 = int32(0)
		_ = i
		for ; i < fsLim; i = i + 1 {
			this_.flowScopes[i].priv_LinkedFlowScope_clear_3737()
		}
	}
	this_.flowScopeCount = int32(0)
	var osLim int32 = min(this_.overlayScopeCount, this_.overlayScopeAlloc)
	_ = osLim
	{
		var i_2 int32 = int32(0)
		_ = i_2
		for ; i_2 < osLim; i_2 = i_2 + 1 {
			this_.overlayScopes[i_2].Clear()
		}
	}
	this_.overlayScopeCount = int32(0)
	var slotLim int32 = min(this_.overlaySlotCount, this_.overlaySlotAlloc)
	_ = slotLim
	{
		var i_3 int32 = int32(0)
		_ = i_3
		for ; i_3 < slotLim; i_3 = i_3 + 1 {
			this_.overlaySlots[i_3].Clear()
		}
	}
	this_.overlaySlotCount = int32(0)
}

// Java class: com.google.javascript.jscomp.LinkedFlowScope.FlowScopeJoinOp
type LinkedFlowScope_FlowScopeJoinOp struct {
	result          *LinkedFlowScope                 `gson:"result"`
	inputProvider   CompilerInputProvider            `gson:"inputProvider"`
	pool            *LinkedFlowScope_Pool            `gson:"pool"`
	scopeReconciler *LinkedFlowScope_ScopeReconciler `gson:"scopeReconciler"`
	idHash_         atomic.Int32
	inUse           bool `gson:"inUse"`
}

func (this_ *LinkedFlowScope_FlowScopeJoinOp) BindSelf_(s any) {
	_ = s
}
func (this_ *LinkedFlowScope_FlowScopeJoinOp) IsJavaNil_() bool { return this_ == nil }
func (*LinkedFlowScope_FlowScopeJoinOp) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.LinkedFlowScope$FlowScopeJoinOp", "FlowScopeJoinOp"
}
func (this_ *LinkedFlowScope_FlowScopeJoinOp) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*LinkedFlowScope_FlowScopeJoinOp) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.LinkedFlowScope.FlowScopeJoinOp", "com.google.javascript.jscomp.LinkedFlowScope$FlowScopeJoinOp", func(v_ any) bool {
		return j2gort.InstanceOf[*LinkedFlowScope_FlowScopeJoinOp](v_)
	})
}
func (*LinkedFlowScope_FlowScopeJoinOp) Implements_jscomp_DataFlowAnalysis_FlowJoiner() {}

func (this_ *LinkedFlowScope_FlowScopeJoinOp) initFields_LinkedFlowScope_FlowScopeJoinOp() {
	_ = this_

}

// Java constructor: FlowScopeJoinOp(com.google.javascript.jscomp.CompilerInputProvider inputProvider, com.google.javascript.jscomp.LinkedFlowScope.@org.jspecify.annotations.Nullable Pool pool)
func NewLinkedFlowScope_FlowScopeJoinOp_2(inputProvider CompilerInputProvider, pool *LinkedFlowScope_Pool) *LinkedFlowScope_FlowScopeJoinOp {
	this_ := &LinkedFlowScope_FlowScopeJoinOp{}
	InitLinkedFlowScope_FlowScopeJoinOp_2(this_, inputProvider, pool)
	return this_
}

func InitLinkedFlowScope_FlowScopeJoinOp_2(this_ *LinkedFlowScope_FlowScopeJoinOp, inputProvider CompilerInputProvider, pool *LinkedFlowScope_Pool) {
	_ = inputProvider
	_ = pool
	this_.initFields_LinkedFlowScope_FlowScopeJoinOp()
	this_.inputProvider = inputProvider
	this_.pool = pool
	this_.scopeReconciler = NewLinkedFlowScope_ScopeReconciler(pool)
}

// Java constructor: FlowScopeJoinOp(com.google.javascript.jscomp.CompilerInputProvider inputProvider)
func NewLinkedFlowScope_FlowScopeJoinOp_1(inputProvider CompilerInputProvider) *LinkedFlowScope_FlowScopeJoinOp {
	this_ := &LinkedFlowScope_FlowScopeJoinOp{}
	InitLinkedFlowScope_FlowScopeJoinOp_1(this_, inputProvider)
	return this_
}

func InitLinkedFlowScope_FlowScopeJoinOp_1(this_ *LinkedFlowScope_FlowScopeJoinOp, inputProvider CompilerInputProvider) {
	_ = inputProvider
	InitLinkedFlowScope_FlowScopeJoinOp_2(this_, inputProvider, nil)
}

// Java method: FlowScopeJoinOp.tryAcquire()
func (this_ *LinkedFlowScope_FlowScopeJoinOp) TryAcquire() bool {
	if this_.inUse {
		return false
	}
	this_.inUse = true
	this_.result = nil
	return true
}

// Java method: FlowScopeJoinOp.joinFlow(com.google.javascript.jscomp.type.FlowScope input)
func (this_ *LinkedFlowScope_FlowScopeJoinOp) JoinFlow(input pkg_jscomptype.FlowScope) {
	_ = input
	var linkedInput *LinkedFlowScope = j2gort.CastPtr[LinkedFlowScope](input)
	_ = linkedInput
	if this_.result == nil {
		this_.result = linkedInput
		return
	} else if (this_.result == linkedInput) || (this_.result.priv_LinkedFlowScope_sameScopes_3739(linkedInput) && (this_.result.functionScope == linkedInput.functionScope)) {
		return
	} else if (this_.result.priv_LinkedFlowScope_flowsFromBottom_3742() && this_.result.priv_LinkedFlowScope_hasNoScopes_3738()) && !(linkedInput.priv_LinkedFlowScope_flowsFromBottom_3742()) {
		this_.result = linkedInput
		return
	} else if (linkedInput.priv_LinkedFlowScope_flowsFromBottom_3742() && linkedInput.priv_LinkedFlowScope_hasNoScopes_3738()) && !(this_.result.priv_LinkedFlowScope_flowsFromBottom_3742()) {
		return
	}
	var common *TypedScope = LinkedFlowScope_GetCommonParentDeclarationScope(this_.result, linkedInput)
	_ = common
	var activePool *LinkedFlowScope_Pool
	_ = activePool
	if this_.pool != nil {
		activePool = this_.pool
	} else if this_.result.pool != nil {
		activePool = this_.result.pool
	} else {
		activePool = linkedInput.pool
	}
	var aMap *pkg_rhino.HamtPMap[*TypedScope, *LinkedFlowScope_OverlayScope] = nil
	_ = aMap
	var aKey *TypedScope = nil
	_ = aKey
	var aVal *LinkedFlowScope_OverlayScope = nil
	_ = aVal
	if !(this_.result.priv_LinkedFlowScope_hasNoScopes_3738()) {
		if this_.result.scopes == nil {
			if (this_.result.syntacticScope == common) || !(this_.result.priv_LinkedFlowScope_trimsAwaySingleKey_3741(common, this_.result.singleScopeKey)) {
				aKey = this_.result.singleScopeKey
				aVal = this_.result.singleScopeVal
			}
		} else {
			var m *pkg_rhino.HamtPMap[*TypedScope, *LinkedFlowScope_OverlayScope]
			_ = m
			if this_.result.syntacticScope == common {
				m = this_.result.scopes
			} else {
				m = LinkedFlowScope_TrimScopesMap(this_.result.scopes, this_.result.syntacticScope, common)
			}
			if (m != nil) && !(m.IsEmpty()) {
				if m.IsSingleEntry() {
					aKey = m.RootKey()
					aVal = m.RootValue()
				} else {
					aMap = m
				}
			}
		}
	}
	var bMap *pkg_rhino.HamtPMap[*TypedScope, *LinkedFlowScope_OverlayScope] = nil
	_ = bMap
	var bKey *TypedScope = nil
	_ = bKey
	var bVal *LinkedFlowScope_OverlayScope = nil
	_ = bVal
	if !(linkedInput.priv_LinkedFlowScope_hasNoScopes_3738()) {
		if linkedInput.scopes == nil {
			if (linkedInput.syntacticScope == common) || !(linkedInput.priv_LinkedFlowScope_trimsAwaySingleKey_3741(common, linkedInput.singleScopeKey)) {
				bKey = linkedInput.singleScopeKey
				bVal = linkedInput.singleScopeVal
			}
		} else {
			var m_2 *pkg_rhino.HamtPMap[*TypedScope, *LinkedFlowScope_OverlayScope]
			_ = m_2
			if linkedInput.syntacticScope == common {
				m_2 = linkedInput.scopes
			} else {
				m_2 = LinkedFlowScope_TrimScopesMap(linkedInput.scopes, linkedInput.syntacticScope, common)
			}
			if (m_2 != nil) && !(m_2.IsEmpty()) {
				if m_2.IsSingleEntry() {
					bKey = m_2.RootKey()
					bVal = m_2.RootValue()
				} else {
					bMap = m_2
				}
			}
		}
	}
	var jMap *pkg_rhino.HamtPMap[*TypedScope, *LinkedFlowScope_OverlayScope] = nil
	_ = jMap
	var jKey *TypedScope = nil
	_ = jKey
	var jVal *LinkedFlowScope_OverlayScope = nil
	_ = jVal
	if ((aMap == nil) && (bMap == nil)) && (((aKey == nil) || (bKey == nil)) || (aKey == bKey)) {
		if (aKey != nil) || (bKey != nil) {
			this_.scopeReconciler.Reset(activePool, this_.result, linkedInput)
			if bKey == nil {
				jKey = aKey
				jVal = this_.scopeReconciler.Merge(aKey, aVal, nil)
			} else if aKey == nil {
				jKey = bKey
				jVal = this_.scopeReconciler.Merge(bKey, nil, bVal)
			} else {
				jKey = aKey
				if aVal == bVal {
					jVal = aVal
				} else {
					jVal = this_.scopeReconciler.Merge(aKey, aVal, bVal)
				}
			}
			this_.scopeReconciler.ClearRefs()
		}
	} else {
		var leftMap *pkg_rhino.HamtPMap[*TypedScope, *LinkedFlowScope_OverlayScope]
		_ = leftMap
		if aMap != nil {
			leftMap = aMap
		} else if aKey != nil {
			leftMap = pkg_rhino.HamtPMap_Singleton[*TypedScope, *LinkedFlowScope_OverlayScope](aKey, aVal)
		} else {
			leftMap = StaticGet_LinkedFlowScope_EMPTY_SCOPES()
		}
		var rightMap *pkg_rhino.HamtPMap[*TypedScope, *LinkedFlowScope_OverlayScope]
		_ = rightMap
		if bMap != nil {
			rightMap = bMap
		} else if bKey != nil {
			rightMap = pkg_rhino.HamtPMap_Singleton[*TypedScope, *LinkedFlowScope_OverlayScope](bKey, bVal)
		} else {
			rightMap = StaticGet_LinkedFlowScope_EMPTY_SCOPES()
		}
		this_.scopeReconciler.Reset(activePool, this_.result, linkedInput)
		var rec *pkg_rhino.HamtPMap[*TypedScope, *LinkedFlowScope_OverlayScope] = leftMap.ReconcileOrNull(rightMap, pkg_rhino.ToIface_PMap_Reconciler[*TypedScope, *LinkedFlowScope_OverlayScope](this_.scopeReconciler))
		_ = rec
		this_.scopeReconciler.ClearRefs()
		if (rec != nil) && !(rec.IsEmpty()) {
			if rec.IsSingleEntry() {
				jKey = rec.RootKey()
				jVal = rec.RootValue()
			} else {
				jMap = rec
			}
		}
	}
	var nextFnScope *TypedScope
	_ = nextFnScope
	if this_.result.priv_LinkedFlowScope_flowsFromBottom_3742() {
		nextFnScope = linkedInput.functionScope
	} else {
		nextFnScope = this_.result.functionScope
	}
	var nextBloom int64
	_ = nextBloom
	if (jMap == nil) && (jKey == nil) {
		nextBloom = int64(0)
	} else {
		nextBloom = (this_.result.unindexedRootBloom | linkedInput.unindexedRootBloom)
	}
	if (((((jMap == this_.result.scopes) && (jKey == this_.result.singleScopeKey)) && (jVal == this_.result.singleScopeVal)) && (common == this_.result.syntacticScope)) && (nextFnScope == this_.result.functionScope)) && (nextBloom == this_.result.unindexedRootBloom) {
		return
	}
	if (((((jMap == linkedInput.scopes) && (jKey == linkedInput.singleScopeKey)) && (jVal == linkedInput.singleScopeVal)) && (common == linkedInput.syntacticScope)) && (nextFnScope == linkedInput.functionScope)) && (nextBloom == linkedInput.unindexedRootBloom) {
		this_.result = linkedInput
		return
	}
	this_.result = LinkedFlowScope_CreateFlowScope(activePool, this_.inputProvider, jMap, jKey, jVal, common, nextFnScope, nextBloom)
}

// Java method: FlowScopeJoinOp.finish()
func (this_ *LinkedFlowScope_FlowScopeJoinOp) Finish() pkg_jscomptype.FlowScope {
	var res *LinkedFlowScope = this_.result
	_ = res
	this_.result = nil
	this_.inUse = false
	return pkg_jscomptype.ToIface_FlowScope(res)
}

// Java class: com.google.javascript.jscomp.LinkedFlowScope.SlotReconciler
type LinkedFlowScope_SlotReconciler struct {
	pool        *LinkedFlowScope_Pool `gson:"pool"`
	typedScopeA *TypedScope           `gson:"typedScopeA"`
	typedScopeB *TypedScope           `gson:"typedScopeB"`
	idHash_     atomic.Int32
}

func (this_ *LinkedFlowScope_SlotReconciler) BindSelf_(s any) {
	_ = s
}
func (this_ *LinkedFlowScope_SlotReconciler) IsJavaNil_() bool { return this_ == nil }
func (*LinkedFlowScope_SlotReconciler) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.LinkedFlowScope$SlotReconciler", "SlotReconciler"
}
func (this_ *LinkedFlowScope_SlotReconciler) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*LinkedFlowScope_SlotReconciler) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.LinkedFlowScope.SlotReconciler", "com.google.javascript.jscomp.LinkedFlowScope$SlotReconciler", func(v_ any) bool {
		return j2gort.InstanceOf[*LinkedFlowScope_SlotReconciler](v_)
	})
}
func (*LinkedFlowScope_SlotReconciler) Implements_rhino_PMap_Reconciler() {}

func (this_ *LinkedFlowScope_SlotReconciler) initFields_LinkedFlowScope_SlotReconciler() {
	_ = this_
}

// Java constructor: SlotReconciler()
func NewLinkedFlowScope_SlotReconciler() *LinkedFlowScope_SlotReconciler {
	this_ := &LinkedFlowScope_SlotReconciler{}
	InitLinkedFlowScope_SlotReconciler(this_)
	return this_
}

func InitLinkedFlowScope_SlotReconciler(this_ *LinkedFlowScope_SlotReconciler) {
	this_.initFields_LinkedFlowScope_SlotReconciler()
}

// Java method: SlotReconciler.merge(java.lang.String slotKey, com.google.javascript.jscomp.LinkedFlowScope.@org.jspecify.annotations.Nullable OverlaySlot slotA, com.google.javascript.jscomp.LinkedFlowScope.@org.jspecify.annotations.Nullable OverlaySlot slotB)
func (this_ *LinkedFlowScope_SlotReconciler) Merge(slotKey string, slotA *LinkedFlowScope_OverlaySlot, slotB *LinkedFlowScope_OverlaySlot) *LinkedFlowScope_OverlaySlot {
	_ = slotKey
	_ = slotA
	_ = slotB
	var name string
	_ = name
	if slotA != nil {
		name = slotA.GetName()
	} else {
		name = slotB.GetName()
	}
	if (slotB == nil) || j2gort.IfaceEquals[pkg_rhino.JSType](slotB.GetType(), nil) {
		var fnSlot *TypedVar
		_ = fnSlot
		if this_.typedScopeB != nil {
			fnSlot = this_.typedScopeB.LookupTypedVar(name)
		} else {
			fnSlot = nil
		}
		var fnSlotType pkg_rhino.JSType
		_ = fnSlotType
		if fnSlot != nil {
			fnSlotType = fnSlot.GetType()
		} else {
			fnSlotType = nil
		}
		if (fnSlotType == nil || fnSlotType.IsJavaNil_()) || j2gort.IfaceEquals[pkg_rhino.JSType](fnSlotType, slotA.GetType()) {
			return slotA
		} else {
			var joinedType pkg_rhino.JSType = slotA.GetType().GetLeastSupertype(fnSlotType)
			_ = joinedType
			if j2gort.IfaceEquals[pkg_rhino.JSType](joinedType, slotA.GetType()) {
				return slotA
			} else {
				return LinkedFlowScope_CreateOverlaySlot(this_.pool, name, joinedType)
			}
		}
	} else if (slotA == nil) || j2gort.IfaceEquals[pkg_rhino.JSType](slotA.GetType(), nil) {
		var fnSlot_2 *TypedVar
		_ = fnSlot_2
		if this_.typedScopeA != nil {
			fnSlot_2 = this_.typedScopeA.LookupTypedVar(name)
		} else {
			fnSlot_2 = nil
		}
		var fnSlotType_2 pkg_rhino.JSType
		_ = fnSlotType_2
		if fnSlot_2 != nil {
			fnSlotType_2 = fnSlot_2.GetType()
		} else {
			fnSlotType_2 = nil
		}
		if (fnSlotType_2 == nil || fnSlotType_2.IsJavaNil_()) || j2gort.IfaceEquals[pkg_rhino.JSType](fnSlotType_2, slotB.GetType()) {
			return slotB
		} else {
			var joinedType_2 pkg_rhino.JSType = slotB.GetType().GetLeastSupertype(fnSlotType_2)
			_ = joinedType_2
			if j2gort.IfaceEquals[pkg_rhino.JSType](joinedType_2, slotB.GetType()) {
				return slotB
			} else {
				return LinkedFlowScope_CreateOverlaySlot(this_.pool, name, joinedType_2)
			}
		}
	}
	if j2gort.IfaceEquals[pkg_rhino.JSType](slotA.GetType(), slotB.GetType()) {
		return slotA
	}
	var joinedType_3 pkg_rhino.JSType = slotA.GetType().GetLeastSupertype(slotB.GetType())
	_ = joinedType_3
	if j2gort.IfaceEquals[pkg_rhino.JSType](joinedType_3, slotA.GetType()) {
		return slotA
	} else {
		return LinkedFlowScope_CreateOverlaySlot(this_.pool, name, joinedType_3)
	}
}

// Java class: com.google.javascript.jscomp.LinkedFlowScope.ScopeReconciler
type LinkedFlowScope_ScopeReconciler struct {
	pool           *LinkedFlowScope_Pool           `gson:"pool"`
	linkedA        *LinkedFlowScope                `gson:"linkedA"`
	linkedB        *LinkedFlowScope                `gson:"linkedB"`
	slotReconciler *LinkedFlowScope_SlotReconciler `gson:"slotReconciler"`
	idHash_        atomic.Int32
}

func (this_ *LinkedFlowScope_ScopeReconciler) BindSelf_(s any) {
	_ = s
}
func (this_ *LinkedFlowScope_ScopeReconciler) IsJavaNil_() bool { return this_ == nil }
func (*LinkedFlowScope_ScopeReconciler) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.LinkedFlowScope$ScopeReconciler", "ScopeReconciler"
}
func (this_ *LinkedFlowScope_ScopeReconciler) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*LinkedFlowScope_ScopeReconciler) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.LinkedFlowScope.ScopeReconciler", "com.google.javascript.jscomp.LinkedFlowScope$ScopeReconciler", func(v_ any) bool {
		return j2gort.InstanceOf[*LinkedFlowScope_ScopeReconciler](v_)
	})
}
func (*LinkedFlowScope_ScopeReconciler) Implements_rhino_PMap_Reconciler() {}

func (this_ *LinkedFlowScope_ScopeReconciler) initFields_LinkedFlowScope_ScopeReconciler() {
	_ = this_
	this_.slotReconciler = NewLinkedFlowScope_SlotReconciler()
}

// Java constructor: ScopeReconciler(com.google.javascript.jscomp.LinkedFlowScope.@org.jspecify.annotations.Nullable Pool pool)
func NewLinkedFlowScope_ScopeReconciler(pool *LinkedFlowScope_Pool) *LinkedFlowScope_ScopeReconciler {
	this_ := &LinkedFlowScope_ScopeReconciler{}
	InitLinkedFlowScope_ScopeReconciler(this_, pool)
	return this_
}

func InitLinkedFlowScope_ScopeReconciler(this_ *LinkedFlowScope_ScopeReconciler, pool *LinkedFlowScope_Pool) {
	_ = pool
	this_.initFields_LinkedFlowScope_ScopeReconciler()
	this_.pool = pool
}

// Java method: ScopeReconciler.reset(com.google.javascript.jscomp.LinkedFlowScope.@org.jspecify.annotations.Nullable Pool pool, com.google.javascript.jscomp.LinkedFlowScope linkedA, com.google.javascript.jscomp.LinkedFlowScope linkedB)
func (this_ *LinkedFlowScope_ScopeReconciler) Reset(pool *LinkedFlowScope_Pool, linkedA *LinkedFlowScope, linkedB *LinkedFlowScope) {
	_ = pool
	_ = linkedA
	_ = linkedB
	this_.pool = pool
	this_.linkedA = linkedA
	this_.linkedB = linkedB
	{
		r_37 := this_.slotReconciler
		r_37.pool = pool
	}
}

// Java method: ScopeReconciler.clearRefs()
func (this_ *LinkedFlowScope_ScopeReconciler) ClearRefs() {
	this_.linkedA = nil
	this_.linkedB = nil
	{
		r_38 := this_.slotReconciler
		r_38.typedScopeA = nil
	}
	{
		r_39 := this_.slotReconciler
		r_39.typedScopeB = nil
	}
}

// Java method: ScopeReconciler.merge(com.google.javascript.jscomp.TypedScope scopeKey, com.google.javascript.jscomp.LinkedFlowScope.@org.jspecify.annotations.Nullable OverlayScope scopeA, com.google.javascript.jscomp.LinkedFlowScope.@org.jspecify.annotations.Nullable OverlayScope scopeB)
func (this_ *LinkedFlowScope_ScopeReconciler) Merge(scopeKey *TypedScope, scopeA *LinkedFlowScope_OverlayScope, scopeB *LinkedFlowScope_OverlayScope) *LinkedFlowScope_OverlayScope {
	_ = scopeKey
	_ = scopeA
	_ = scopeB
	var typedScopeA *TypedScope
	_ = typedScopeA
	if this_.linkedA.priv_LinkedFlowScope_flowsFromBottom_3742() {
		typedScopeA = nil
	} else if scopeA != nil {
		typedScopeA = scopeA.scope
	} else {
		typedScopeA = scopeB.scope
	}
	var typedScopeB *TypedScope
	_ = typedScopeB
	if this_.linkedB.priv_LinkedFlowScope_flowsFromBottom_3742() {
		typedScopeB = nil
	} else if scopeB != nil {
		typedScopeB = scopeB.scope
	} else {
		typedScopeB = scopeA.scope
	}
	var bestScope *TypedScope
	_ = bestScope
	if typedScopeA != nil {
		bestScope = typedScopeA
	} else {
		bestScope = typedScopeB
	}
	if bestScope != nil {
		bestScope = bestScope
	} else if scopeA != nil {
		bestScope = scopeA.scope
	} else {
		bestScope = scopeB.scope
	}
	{
		r_40 := this_.slotReconciler
		r_40.typedScopeA = typedScopeA
	}
	{
		r_41 := this_.slotReconciler
		r_41.typedScopeB = typedScopeB
	}
	if ((scopeA == nil) || (scopeA.slots == nil)) && ((scopeB == nil) || (scopeB.slots == nil)) {
		var a0 *LinkedFlowScope_OverlaySlot
		_ = a0
		if scopeA != nil {
			a0 = scopeA.slot0
		} else {
			a0 = nil
		}
		var a1 *LinkedFlowScope_OverlaySlot
		_ = a1
		if scopeA != nil {
			a1 = scopeA.slot1
		} else {
			a1 = nil
		}
		var b0 *LinkedFlowScope_OverlaySlot
		_ = b0
		if scopeB != nil {
			b0 = scopeB.slot0
		} else {
			b0 = nil
		}
		var b1 *LinkedFlowScope_OverlaySlot
		_ = b1
		if scopeB != nil {
			b1 = scopeB.slot1
		} else {
			b1 = nil
		}
		var aLen int32
		_ = aLen
		if a0 == nil {
			aLen = int32(0)
		} else {
			aLen = (func() int32 {
				if a1 == nil {
					return int32(1)
				}
				return int32(2)
			}())
		}
		var bLen int32
		_ = bLen
		if b0 == nil {
			bLen = int32(0)
		} else {
			bLen = (func() int32 {
				if b1 == nil {
					return int32(1)
				}
				return int32(2)
			}())
		}
		var k0a *LinkedFlowScope_OverlaySlot = nil
		_ = k0a
		var k0b *LinkedFlowScope_OverlaySlot = nil
		_ = k0b
		var k1a *LinkedFlowScope_OverlaySlot = nil
		_ = k1a
		var k1b *LinkedFlowScope_OverlaySlot = nil
		_ = k1b
		var outCount int32 = int32(0)
		_ = outCount
		var fallbackToMap bool = false
		_ = fallbackToMap
		var ia int32 = int32(0)
		_ = ia
		var ib int32 = int32(0)
		_ = ib
		for (ia < aLen) || (ib < bLen) {
			var sa *LinkedFlowScope_OverlaySlot
			_ = sa
			if ia == int32(0) {
				sa = a0
			} else {
				sa = a1
			}
			var sb *LinkedFlowScope_OverlaySlot
			_ = sb
			if ib == int32(0) {
				sb = b0
			} else {
				sb = b1
			}
			var pickA *LinkedFlowScope_OverlaySlot = nil
			_ = pickA
			var pickB *LinkedFlowScope_OverlaySlot = nil
			_ = pickB
			if ib >= bLen {
				pickA = sa
				ia = ia + 1
			} else if ia >= aLen {
				pickB = sb
				ib = ib + 1
			} else {
				var iifeCond_1 bool
				{

					sr_42 := sa.name
					se_43 := sb.name
					iifeCond_1 = j2gort.CheckNotNullString(sr_42) == se_43
				}
				if iifeCond_1 {
					pickA = sa
					pickB = sb
					ia = ia + 1
					ib = ib + 1
				} else {
					var cmp int32 = LinkedFlowScope_CompareUnsigned(LinkedFlowScope_RevHash(sa.name), LinkedFlowScope_RevHash(sb.name))
					_ = cmp
					if cmp == int32(0) {
						fallbackToMap = true
						break
					} else if cmp < int32(0) {
						pickA = sa
						ia = ia + 1
					} else {
						pickB = sb
						ib = ib + 1
					}
				}
			}
			if outCount == int32(0) {
				k0a = pickA
				k0b = pickB
				outCount = int32(1)
			} else if outCount == int32(1) {
				k1a = pickA
				k1b = pickB
				outCount = int32(2)
			} else {
				fallbackToMap = true
				break
			}
		}
		if !(fallbackToMap) {
			var m0 *LinkedFlowScope_OverlaySlot = nil
			_ = m0
			var m1 *LinkedFlowScope_OverlaySlot = nil
			_ = m1
			if outCount >= int32(1) {
				if (k0a != nil) && (k0a == k0b) {
					m0 = k0a
				} else {
					m0 = this_.slotReconciler.Merge(func() string {
						if k0a != nil {
							return k0a.name
						}
						return k0b.name
					}(), k0a, k0b)
				}
				if outCount == int32(2) {
					if (k1a != nil) && (k1a == k1b) {
						m1 = k1a
					} else {
						m1 = this_.slotReconciler.Merge(func() string {
							if k1a != nil {
								return k1a.name
							}
							return k1b.name
						}(), k1a, k1b)
					}
				}
			}
			if (((scopeA != nil) && (scopeA.scope == bestScope)) && (m0 == a0)) && (m1 == a1) {
				return scopeA
			}
			if (((scopeB != nil) && (scopeB.scope == bestScope)) && (m0 == b0)) && (m1 == b1) {
				return scopeB
			}
			var hasQual bool = ((m0 != nil) && (((scopeA != nil) && scopeA.hasQualifiedSlots) || ((scopeB != nil) && scopeB.hasQualifiedSlots)))
			_ = hasQual
			return LinkedFlowScope_CreateOverlayScope(this_.pool, bestScope, nil, m0, m1, hasQual)
		}
	}
	var slotsA *pkg_rhino.HamtPMap[string, *LinkedFlowScope_OverlaySlot]
	_ = slotsA
	if scopeA != nil {
		slotsA = scopeA.ToSlotsMap()
	} else {
		slotsA = StaticGet_LinkedFlowScope_EMPTY_SLOTS()
	}
	var slotsB *pkg_rhino.HamtPMap[string, *LinkedFlowScope_OverlaySlot]
	_ = slotsB
	if scopeB != nil {
		slotsB = scopeB.ToSlotsMap()
	} else {
		slotsB = StaticGet_LinkedFlowScope_EMPTY_SLOTS()
	}
	var reconciledSlots *pkg_rhino.HamtPMap[string, *LinkedFlowScope_OverlaySlot] = slotsA.ReconcileOrNull(slotsB, pkg_rhino.ToIface_PMap_Reconciler[string, *LinkedFlowScope_OverlaySlot](this_.slotReconciler))
	_ = reconciledSlots
	if (((scopeA != nil) && (scopeA.scope == bestScope)) && (scopeA.slots != nil)) && (reconciledSlots == scopeA.slots) {
		return scopeA
	}
	if (((scopeB != nil) && (scopeB.scope == bestScope)) && (scopeB.slots != nil)) && (reconciledSlots == scopeB.slots) {
		return scopeB
	}
	var hasQual_2 bool = (((reconciledSlots != nil) && !(reconciledSlots.IsEmpty())) && (((scopeA != nil) && scopeA.hasQualifiedSlots) || ((scopeB != nil) && scopeB.hasQualifiedSlots)))
	_ = hasQual_2
	return LinkedFlowScope_CreateOverlayScopeFromMap(this_.pool, bestScope, reconciledSlots, hasQual_2)
}

// Java class: com.google.javascript.jscomp.LinkedFlowScope.OverlayScope
type LinkedFlowScope_OverlayScope struct {
	pool              *LinkedFlowScope_Pool                                     `gson:"pool"`
	scope             *TypedScope                                               `gson:"scope"`
	slots             *pkg_rhino.HamtPMap[string, *LinkedFlowScope_OverlaySlot] `gson:"slots"`
	slot0             *LinkedFlowScope_OverlaySlot                              `gson:"slot0"`
	slot1             *LinkedFlowScope_OverlaySlot                              `gson:"slot1"`
	idHash_           atomic.Int32
	hasQualifiedSlots bool `gson:"hasQualifiedSlots"`
}

func (this_ *LinkedFlowScope_OverlayScope) BindSelf_(s any) {
	_ = s
}
func (this_ *LinkedFlowScope_OverlayScope) IsJavaNil_() bool { return this_ == nil }
func (*LinkedFlowScope_OverlayScope) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.LinkedFlowScope$OverlayScope", "OverlayScope"
}
func (this_ *LinkedFlowScope_OverlayScope) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*LinkedFlowScope_OverlayScope) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.LinkedFlowScope.OverlayScope", "com.google.javascript.jscomp.LinkedFlowScope$OverlayScope", func(v_ any) bool {
		return j2gort.InstanceOf[*LinkedFlowScope_OverlayScope](v_)
	})
}

func (this_ *LinkedFlowScope_OverlayScope) initFields_LinkedFlowScope_OverlayScope() {
	_ = this_
}

// Java constructor: OverlayScope()
func NewLinkedFlowScope_OverlayScope() *LinkedFlowScope_OverlayScope {
	this_ := &LinkedFlowScope_OverlayScope{}
	InitLinkedFlowScope_OverlayScope(this_)
	return this_
}

func InitLinkedFlowScope_OverlayScope(this_ *LinkedFlowScope_OverlayScope) {
	this_.initFields_LinkedFlowScope_OverlayScope()
}

// Java method: OverlayScope.init(com.google.javascript.jscomp.LinkedFlowScope.@org.jspecify.annotations.Nullable Pool pool, com.google.javascript.jscomp.TypedScope scope, com.google.javascript.rhino.@org.jspecify.annotations.Nullable HamtPMap<java.lang.String,com.google.javascript.jscomp.LinkedFlowScope.OverlaySlot> slots, com.google.javascript.jscomp.LinkedFlowScope.@org.jspecify.annotations.Nullable OverlaySlot slot0, com.google.javascript.jscomp.LinkedFlowScope.@org.jspecify.annotations.Nullable OverlaySlot slot1, boolean hasQualifiedSlots)
func (this_ *LinkedFlowScope_OverlayScope) Init(pool *LinkedFlowScope_Pool, scope *TypedScope, slots *pkg_rhino.HamtPMap[string, *LinkedFlowScope_OverlaySlot], slot0 *LinkedFlowScope_OverlaySlot, slot1 *LinkedFlowScope_OverlaySlot, hasQualifiedSlots bool) {
	_ = pool
	_ = scope
	_ = slots
	_ = slot0
	_ = slot1
	_ = hasQualifiedSlots
	this_.pool = pool
	this_.scope = j2gort.CheckNotNullPtr(scope)
	this_.slots = slots
	this_.slot0 = slot0
	this_.slot1 = slot1
	this_.hasQualifiedSlots = hasQualifiedSlots
}

// Java method: OverlayScope.clear()
func (this_ *LinkedFlowScope_OverlayScope) Clear() {
	this_.pool = nil
	this_.scope = nil
	this_.slots = nil
	this_.slot0 = nil
	this_.slot1 = nil
	this_.hasQualifiedSlots = false
}

// Java method: OverlayScope.isEmpty()
func (this_ *LinkedFlowScope_OverlayScope) IsEmpty() bool {
	return ((this_.slots == nil) && (this_.slot0 == nil))
}

// Java method: OverlayScope.containsSlotWithHash(int rawKeyHash)
func (this_ *LinkedFlowScope_OverlayScope) ContainsSlotWithHash(rawKeyHash int32) bool {
	_ = rawKeyHash
	if this_.slots != nil {
		return this_.slots.ContainsHash(rawKeyHash)
	}
	var s0 *LinkedFlowScope_OverlaySlot = this_.slot0
	_ = s0
	if s0 == nil {
		return false
	}
	if javalang.StringHashCode(j2gort.CheckNotNullString(s0.name)) == rawKeyHash {
		return true
	}
	var s1 *LinkedFlowScope_OverlaySlot = this_.slot1
	_ = s1
	return ((s1 != nil) && (javalang.StringHashCode(j2gort.CheckNotNullString(s1.name)) == rawKeyHash))
}

// Java method: OverlayScope.toSlotsMap()
func (this_ *LinkedFlowScope_OverlayScope) ToSlotsMap() *pkg_rhino.HamtPMap[string, *LinkedFlowScope_OverlaySlot] {
	if this_.slots != nil {
		return this_.slots
	}
	var s0 *LinkedFlowScope_OverlaySlot = this_.slot0
	_ = s0
	if s0 == nil {
		return StaticGet_LinkedFlowScope_EMPTY_SLOTS()
	}
	var s1 *LinkedFlowScope_OverlaySlot = this_.slot1
	_ = s1
	if s1 == nil {
		return pkg_rhino.HamtPMap_Singleton[string, *LinkedFlowScope_OverlaySlot](s0.name, s0)
	}
	return pkg_rhino.HamtPMap_OfTwo[string, *LinkedFlowScope_OverlaySlot](s0.name, s0, s1.name, s1)
}

// Java method: OverlayScope.infer(java.lang.String name, com.google.javascript.rhino.jstype.JSType type, boolean symbolIsQualified)
func (this_ *LinkedFlowScope_OverlayScope) Infer(name string, type_ pkg_rhino.JSType, symbolIsQualified bool) *LinkedFlowScope_OverlayScope {
	_ = name
	_ = type_
	_ = symbolIsQualified
	if this_.slots == nil {
		var s0 *LinkedFlowScope_OverlaySlot = this_.slot0
		_ = s0
		if s0 == nil {
			var nextSlot *LinkedFlowScope_OverlaySlot = LinkedFlowScope_CreateOverlaySlot(this_.pool, name, type_)
			_ = nextSlot
			return LinkedFlowScope_CreateOverlayScope(this_.pool, this_.scope, nil, nextSlot, nil, (this_.hasQualifiedSlots || symbolIsQualified))
		}
		if j2gort.CheckNotNullString(s0.name) == name {
			if type_ == s0.type_ || ((type_ == nil || type_.IsJavaNil_()) && (s0.type_ == nil || s0.type_.IsJavaNil_())) {
				return this_
			}
			var nextSlot_2 *LinkedFlowScope_OverlaySlot = LinkedFlowScope_CreateOverlaySlot(this_.pool, name, type_)
			_ = nextSlot_2
			return LinkedFlowScope_CreateOverlayScope(this_.pool, this_.scope, nil, nextSlot_2, this_.slot1, (this_.hasQualifiedSlots || symbolIsQualified))
		}
		var s1 *LinkedFlowScope_OverlaySlot = this_.slot1
		_ = s1
		if s1 == nil {
			var nextSlot_3 *LinkedFlowScope_OverlaySlot = LinkedFlowScope_CreateOverlaySlot(this_.pool, name, type_)
			_ = nextSlot_3
			var nextQual bool = (this_.hasQualifiedSlots || symbolIsQualified)
			_ = nextQual
			var cmp int32 = LinkedFlowScope_CompareUnsigned(LinkedFlowScope_RevHash(s0.name), LinkedFlowScope_RevHash(name))
			_ = cmp
			if cmp < int32(0) {
				return LinkedFlowScope_CreateOverlayScope(this_.pool, this_.scope, nil, s0, nextSlot_3, nextQual)
			} else if cmp > int32(0) {
				return LinkedFlowScope_CreateOverlayScope(this_.pool, this_.scope, nil, nextSlot_3, s0, nextQual)
			} else {
				var map_ *pkg_rhino.HamtPMap[string, *LinkedFlowScope_OverlaySlot] = pkg_rhino.HamtPMap_OfTwo[string, *LinkedFlowScope_OverlaySlot](s0.name, s0, name, nextSlot_3)
				_ = map_
				return LinkedFlowScope_CreateOverlayScope(this_.pool, this_.scope, map_, nil, nil, nextQual)
			}
		}
		if j2gort.CheckNotNullString(s1.name) == name {
			if type_ == s1.type_ || ((type_ == nil || type_.IsJavaNil_()) && (s1.type_ == nil || s1.type_.IsJavaNil_())) {
				return this_
			}
			var nextSlot_4 *LinkedFlowScope_OverlaySlot = LinkedFlowScope_CreateOverlaySlot(this_.pool, name, type_)
			_ = nextSlot_4
			return LinkedFlowScope_CreateOverlayScope(this_.pool, this_.scope, nil, s0, nextSlot_4, (this_.hasQualifiedSlots || symbolIsQualified))
		}
		var nextSlot_5 *LinkedFlowScope_OverlaySlot = LinkedFlowScope_CreateOverlaySlot(this_.pool, name, type_)
		_ = nextSlot_5
		var nextQual_2 bool = (this_.hasQualifiedSlots || symbolIsQualified)
		_ = nextQual_2
		var map__2 *pkg_rhino.HamtPMap[string, *LinkedFlowScope_OverlaySlot] = j2gort.CastPtr[pkg_rhino.HamtPMap[string, *LinkedFlowScope_OverlaySlot]](pkg_rhino.HamtPMap_OfTwo[string, *LinkedFlowScope_OverlaySlot](s0.name, s0, s1.name, s1).Plus_2(name, nextSlot_5))
		_ = map__2
		return LinkedFlowScope_CreateOverlayScope(this_.pool, this_.scope, map__2, nil, nil, nextQual_2)
	}
	var slot *LinkedFlowScope_OverlaySlot = this_.slots.Get_1(name)
	_ = slot
	if (slot != nil) && (type_ == slot.type_ || ((type_ == nil || type_.IsJavaNil_()) && (slot.type_ == nil || slot.type_.IsJavaNil_()))) {
		return this_
	}
	var nextHasQualified bool = (this_.hasQualifiedSlots || symbolIsQualified)
	_ = nextHasQualified
	var nextSlot_6 *LinkedFlowScope_OverlaySlot = LinkedFlowScope_CreateOverlaySlot(this_.pool, name, type_)
	_ = nextSlot_6
	return LinkedFlowScope_CreateOverlayScope(this_.pool, this_.scope, j2gort.CastPtr[pkg_rhino.HamtPMap[string, *LinkedFlowScope_OverlaySlot]](this_.slots.Plus_2(name, nextSlot_6)), nil, nil, nextHasQualified)
}

// Java method: OverlayScope.getOverlayOrFallback(java.lang.String name, com.google.javascript.jscomp.@org.jspecify.annotations.Nullable TypedVar fallbackVar)
func (this_ *LinkedFlowScope_OverlayScope) GetOverlayOrFallback(name string, fallbackVar *TypedVar) pkg_rhino.StaticTypedSlot {
	_ = name
	_ = fallbackVar
	if this_.slots != nil {
		var slot *LinkedFlowScope_OverlaySlot = this_.slots.Get_1(name)
		_ = slot
		if slot != nil {
			return pkg_rhino.ToIface_StaticTypedSlot(slot)
		} else {
			return pkg_rhino.ToIface_StaticTypedSlot(fallbackVar)
		}
	}
	var s0 *LinkedFlowScope_OverlaySlot = this_.slot0
	_ = s0
	if s0 != nil {
		if j2gort.CheckNotNullString(s0.name) == name {
			return pkg_rhino.ToIface_StaticTypedSlot(s0)
		}
		var s1 *LinkedFlowScope_OverlaySlot = this_.slot1
		_ = s1
		if (s1 != nil) && (j2gort.CheckNotNullString(s1.name) == name) {
			return pkg_rhino.ToIface_StaticTypedSlot(s1)
		}
	}
	return pkg_rhino.ToIface_StaticTypedSlot(fallbackVar)
}

// Java method: OverlayScope.getSlot(java.lang.String name)
func (this_ *LinkedFlowScope_OverlayScope) GetSlot(name string) pkg_rhino.StaticTypedSlot {
	_ = name
	if this_.slots != nil {
		var slot *LinkedFlowScope_OverlaySlot = this_.slots.Get_1(name)
		_ = slot
		if slot != nil {
			return pkg_rhino.ToIface_StaticTypedSlot(slot)
		} else {
			return pkg_rhino.ToIface_StaticTypedSlot(this_.scope.LookupTypedVar(name))
		}
	}
	var s0 *LinkedFlowScope_OverlaySlot = this_.slot0
	_ = s0
	if s0 != nil {
		if j2gort.CheckNotNullString(s0.name) == name {
			return pkg_rhino.ToIface_StaticTypedSlot(s0)
		}
		var s1 *LinkedFlowScope_OverlaySlot = this_.slot1
		_ = s1
		if (s1 != nil) && (j2gort.CheckNotNullString(s1.name) == name) {
			return pkg_rhino.ToIface_StaticTypedSlot(s1)
		}
	}
	return pkg_rhino.ToIface_StaticTypedSlot(this_.scope.LookupTypedVar(name))
}

// Java class: com.google.javascript.jscomp.LinkedFlowScope.OverlaySlot
type LinkedFlowScope_OverlaySlot struct {
	name    string           `gson:"name"`
	type_   pkg_rhino.JSType `gson:"type"`
	idHash_ atomic.Int32
}

func (this_ *LinkedFlowScope_OverlaySlot) BindSelf_(s any) {
	_ = s
}
func (this_ *LinkedFlowScope_OverlaySlot) IsJavaNil_() bool { return this_ == nil }
func (*LinkedFlowScope_OverlaySlot) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.LinkedFlowScope$OverlaySlot", "OverlaySlot"
}
func (this_ *LinkedFlowScope_OverlaySlot) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*LinkedFlowScope_OverlaySlot) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.LinkedFlowScope.OverlaySlot", "com.google.javascript.jscomp.LinkedFlowScope$OverlaySlot", func(v_ any) bool {
		return j2gort.InstanceOf[*LinkedFlowScope_OverlaySlot](v_)
	})
}
func (*LinkedFlowScope_OverlaySlot) Implements_rhino_StaticTypedSlot() {}
func (*LinkedFlowScope_OverlaySlot) Implements_rhino_StaticSlot()      {}

func (this_ *LinkedFlowScope_OverlaySlot) initFields_LinkedFlowScope_OverlaySlot() {
	_ = this_
	this_.name = j2gort.NullString
}

// Java constructor: OverlaySlot()
func NewLinkedFlowScope_OverlaySlot() *LinkedFlowScope_OverlaySlot {
	this_ := &LinkedFlowScope_OverlaySlot{}
	InitLinkedFlowScope_OverlaySlot(this_)
	return this_
}

func InitLinkedFlowScope_OverlaySlot(this_ *LinkedFlowScope_OverlaySlot) {
	this_.initFields_LinkedFlowScope_OverlaySlot()
}

// Java method: OverlaySlot.init(java.lang.String name, com.google.javascript.rhino.jstype.JSType type)
func (this_ *LinkedFlowScope_OverlaySlot) Init(name string, type_ pkg_rhino.JSType) {
	_ = name
	_ = type_
	this_.name = name
	this_.type_ = type_
}

// Java method: OverlaySlot.clear()
func (this_ *LinkedFlowScope_OverlaySlot) Clear() {
	this_.name = j2gort.NullString
	this_.type_ = nil
}

// Java method: OverlaySlot.getName()
func (this_ *LinkedFlowScope_OverlaySlot) GetName() string {
	return this_.name
}

// Java method: OverlaySlot.getType()
func (this_ *LinkedFlowScope_OverlaySlot) GetType() pkg_rhino.JSType {
	return this_.type_
}

// Java method: OverlaySlot.isTypeInferred()
func (this_ *LinkedFlowScope_OverlaySlot) IsTypeInferred() bool {
	return true
}

// Java method: OverlaySlot.getDeclaration()
func (this_ *LinkedFlowScope_OverlaySlot) GetDeclaration() any {
	return nil
}

// Java method: OverlaySlot.getJSDocInfo()
func (this_ *LinkedFlowScope_OverlaySlot) GetJSDocInfo() *pkg_rhino.JSDocInfo {
	return nil
}

// Java method: OverlaySlot.getScope()
func (this_ *LinkedFlowScope_OverlaySlot) GetScope() any {
	panic(j2gort.AsIface[javalang.IUnsupportedOperationException](j2gort.ConstructShimThrowable("java.lang.UnsupportedOperationException")))
}
