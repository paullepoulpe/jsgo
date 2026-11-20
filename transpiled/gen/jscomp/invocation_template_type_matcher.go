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
// Port of com/google/javascript/jscomp/InvocationTemplateTypeMatcher.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	pkg_base "github.com/paullepoulpe/jsgo/transpiled/gen/base"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

var lambdaSingleton_invocation_template_type_matcher_1 j2gort.Ext_java_util_function_BiFunction

func init() {
	lambdaSingleton_invocation_template_type_matcher_1 = j2gort.NewLambda_Ext_java_util_function_BiFunction(func(mrp0_ any, mrp1_ any) any {
		return j2gort.AsIface[pkg_rhino.JSType](mrp0_).GetLeastSupertype(j2gort.AsIface[pkg_rhino.JSType](mrp1_))
	})
}

// Java class: com.google.javascript.jscomp.InvocationTemplateTypeMatcher
type InvocationTemplateTypeMatcher struct {
	matchedKey0   *pkg_rhino.TemplateType                                                    `gson:"matchedKey0"`
	matchedVal0   pkg_rhino.JSType                                                           `gson:"matchedVal0"`
	matchedTypes  *pkg_base.LinkedIdentityHashMap[*pkg_rhino.TemplateType, pkg_rhino.JSType] `gson:"matchedTypes"`
	seenTypes     *pkg_base.LinkedIdentityHashSet[pkg_rhino.JSType]                          `gson:"seenTypes"`
	registry      *pkg_rhino.JSTypeRegistry                                                  `gson:"registry"`
	calleeType    *pkg_rhino.FunctionType                                                    `gson:"calleeType"`
	localThisType pkg_rhino.JSType                                                           `gson:"localThisType"`
	unknownType   pkg_rhino.JSType                                                           `gson:"unknownType"`
	invocation    pkg_rhino.Node                                                             `gson:"invocation"`
	idHash_       atomic.Int32
}

func (this_ *InvocationTemplateTypeMatcher) BindSelf_(s any) {
	_ = s
}
func (this_ *InvocationTemplateTypeMatcher) IsJavaNil_() bool { return this_ == nil }
func (*InvocationTemplateTypeMatcher) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.InvocationTemplateTypeMatcher", "InvocationTemplateTypeMatcher"
}
func (this_ *InvocationTemplateTypeMatcher) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*InvocationTemplateTypeMatcher) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.InvocationTemplateTypeMatcher", "com.google.javascript.jscomp.InvocationTemplateTypeMatcher", func(v_ any) bool {
		return j2gort.InstanceOf[*InvocationTemplateTypeMatcher](v_)
	})
}

func (this_ *InvocationTemplateTypeMatcher) initFields_InvocationTemplateTypeMatcher() {
	_ = this_

}

// Java constructor: InvocationTemplateTypeMatcher(com.google.javascript.rhino.jstype.JSTypeRegistry registry, com.google.javascript.rhino.jstype.FunctionType calleeType, com.google.javascript.rhino.jstype.JSType localThisType, com.google.javascript.rhino.Node invocation)
func NewInvocationTemplateTypeMatcher(registry *pkg_rhino.JSTypeRegistry, calleeType *pkg_rhino.FunctionType, localThisType pkg_rhino.JSType, invocation pkg_rhino.Node) *InvocationTemplateTypeMatcher {
	this_ := &InvocationTemplateTypeMatcher{}
	InitInvocationTemplateTypeMatcher(this_, registry, calleeType, localThisType, invocation)
	return this_
}

func InitInvocationTemplateTypeMatcher(this_ *InvocationTemplateTypeMatcher, registry *pkg_rhino.JSTypeRegistry, calleeType *pkg_rhino.FunctionType, localThisType pkg_rhino.JSType, invocation pkg_rhino.Node) {
	_ = registry
	_ = calleeType
	_ = localThisType
	_ = invocation
	this_.initFields_InvocationTemplateTypeMatcher()
	this_.registry = registry
	this_.calleeType = calleeType
	this_.localThisType = localThisType
	this_.invocation = invocation
	this_.unknownType = registry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_UNKNOWN_TYPE())
}

// Java method: InvocationTemplateTypeMatcher.runMatch()
func (this_ *InvocationTemplateTypeMatcher) RunMatch() {
	if this_.calleeType.GetTemplateTypeMap().IsEmpty() {
		return
	}
	var target pkg_rhino.Node = this_.invocation.GetFirstChild()
	_ = target
	if NodeUtil_IsNormalGet(target) {
		var obj pkg_rhino.Node = target.GetFirstChild()
		_ = obj
		var typeOfThisRequiredByTheFunction pkg_rhino.JSType = this_.calleeType.GetTypeOfThis()
		_ = typeOfThisRequiredByTheFunction
		var typeOfThisProvidedByTheCall pkg_rhino.JSType
		_ = typeOfThisProvidedByTheCall
		if obj.IsSuper() {
			typeOfThisProvidedByTheCall = this_.localThisType
		} else {
			typeOfThisProvidedByTheCall = this_.priv_InvocationTemplateTypeMatcher_getTypeOrUnknown_3459(obj)
		}
		typeOfThisProvidedByTheCall = typeOfThisProvidedByTheCall.RestrictByNotNullOrUndefined()
		this_.priv_InvocationTemplateTypeMatcher_matchTemplateTypesRecursive_3454(typeOfThisRequiredByTheFunction, typeOfThisProvidedByTheCall)
	}
	var calleeParameters []pkg_rhino.FunctionType_Parameter = this_.calleeType.GetParameterArray()
	_ = calleeParameters
	if this_.invocation.IsTaggedTemplateLit() {
		if int32(len(calleeParameters)) == int32(0) {
			return
		}
		this_.priv_InvocationTemplateTypeMatcher_matchTemplateTypesRecursive_3454(calleeParameters[int32(0)].GetJSType(), this_.registry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_I_TEMPLATE_ARRAY_TYPE()))
		this_.priv_InvocationTemplateTypeMatcher_matchTemplateTypesFromNodes_3455(calleeParameters, int32(1), true)
	} else if this_.invocation.HasMoreThanOneChild() {
		this_.priv_InvocationTemplateTypeMatcher_matchTemplateTypesFromNodes_3455(calleeParameters, int32(0), false)
	}
}

// Java method: InvocationTemplateTypeMatcher.lookupMatched(com.google.javascript.rhino.jstype.TemplateType key)
func (this_ *InvocationTemplateTypeMatcher) LookupMatched(key *pkg_rhino.TemplateType) pkg_rhino.JSType {
	_ = key
	if this_.matchedTypes != nil {
		return this_.matchedTypes.Get(key)
	}
	if this_.matchedKey0 == key {
		return this_.matchedVal0
	} else {
		return nil
	}
}

// Java method: InvocationTemplateTypeMatcher.match()
func (this_ *InvocationTemplateTypeMatcher) Match() j2gort.Ext_com_google_common_collect_ImmutableMap {
	this_.RunMatch()
	return this_.priv_InvocationTemplateTypeMatcher_copyMatchedTypes_3458()
}

// Java method: InvocationTemplateTypeMatcher.matchTemplateTypesRecursive(com.google.javascript.rhino.jstype.JSType paramType, com.google.javascript.rhino.jstype.JSType argType)
func (this_ *InvocationTemplateTypeMatcher) priv_InvocationTemplateTypeMatcher_matchTemplateTypesRecursive_3454(paramType pkg_rhino.JSType, argType pkg_rhino.JSType) {
	_ = paramType
	_ = argType
	base_paramType_ := paramType.GetJSTypeBase()
	if base_paramType_.IsTemplateType() {
		this_.priv_InvocationTemplateTypeMatcher_recordTemplateMatch_3457(paramType.ToMaybeTemplateType(), argType)
		return
	}
	if base_paramType_.IsUnionType() {
		var union *pkg_rhino.UnionType = paramType.ToMaybeUnionType()
		_ = union
		{
			var i int32 = int32(0)
			_ = i
			var len_ int32 = union.GetAlternatesCount()
			_ = len_
			for ; i < len_; i = i + 1 {
				this_.priv_InvocationTemplateTypeMatcher_matchTemplateTypesRecursive_3454(union.GetAlternateAt(i), argType)
			}
		}
		return
	} else if argType.GetJSTypeBase().IsUnionType() {
		var union_2 *pkg_rhino.UnionType = argType.ToMaybeUnionType()
		_ = union_2
		{
			var i_2 int32 = int32(0)
			_ = i_2
			var len__2 int32 = union_2.GetAlternatesCount()
			_ = len__2
			for ; i_2 < len__2; i_2 = i_2 + 1 {
				this_.priv_InvocationTemplateTypeMatcher_matchTemplateTypesRecursive_3454(paramType, union_2.GetAlternateAt(i_2))
			}
		}
		return
	}
	if base_paramType_.IsFunctionType() {
		var paramFunctionType *pkg_rhino.FunctionType = paramType.ToMaybeFunctionType()
		_ = paramFunctionType
		var argFunctionType *pkg_rhino.FunctionType = argType.RestrictByNotNullOrUndefined().CollapseUnion().ToMaybeFunctionType()
		_ = argFunctionType
		if (argFunctionType != nil) && argFunctionType.IsSubtype_1(paramType) {
			this_.priv_InvocationTemplateTypeMatcher_matchTemplateTypesRecursive_3454(paramFunctionType.GetTypeOfThis(), argFunctionType.GetTypeOfThis())
			this_.priv_InvocationTemplateTypeMatcher_matchTemplateTypesRecursive_3454(paramFunctionType.GetReturnType(), argFunctionType.GetReturnType())
			this_.priv_InvocationTemplateTypeMatcher_matchTemplateTypesFromParameters_3456(paramFunctionType.GetParameterArray(), argFunctionType.GetParameterArray())
		}
	} else if base_paramType_.IsRecordType() && !(paramType.IsNominalType()) {
		if this_.seenTypes == nil {
			this_.seenTypes = pkg_base.NewLinkedIdentityHashSet[pkg_rhino.JSType]()
		}
		if this_.seenTypes.Add(paramType) {
			var paramRecordType pkg_rhino.ObjectType = base_paramType_.ToObjectType()
			_ = paramRecordType
			var argObjectType pkg_rhino.ObjectType = argType.RestrictByNotNullOrUndefined().GetJSTypeBase().ToObjectType()
			_ = argObjectType
			if ((argObjectType != nil && !argObjectType.IsJavaNil_()) && !(argObjectType.GetObjectTypeBase().IsUnknownType())) && !(argObjectType.GetObjectTypeBase().IsEmptyType()) {
				var names j2gort.Ext_com_google_common_collect_ImmutableSortedSet = paramRecordType.GetObjectTypeBase().GetPropertyNames()
				_ = names
				{
					cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(names))
					for cur_1.HasNext() {
						name := cur_1.NextStr()
						_ = name
						if paramRecordType.GetObjectTypeBase().HasOwnProperty_1_string(name) && argObjectType.GetObjectTypeBase().HasProperty_1_string(name) {
							this_.priv_InvocationTemplateTypeMatcher_matchTemplateTypesRecursive_3454(paramRecordType.GetObjectTypeBase().GetPropertyType_1_string(name), argObjectType.GetObjectTypeBase().GetPropertyType_1_string(name))
						}
					}
				}
			}
			this_.seenTypes.Remove(paramType)
		}
	} else if base_paramType_.IsTemplatizedType() {
		var templatizedParamType *pkg_rhino.TemplatizedType = paramType.ToMaybeTemplatizedType()
		_ = templatizedParamType
		var keyCount int32 = j2gort.ImmSliceSize(j2gort.CheckNotNullPtr(templatizedParamType).Slice_GetTemplateTypes())
		_ = keyCount
		if keyCount == int32(0) {
			return
		}
		var argObjectType_2 pkg_rhino.JSType = argType.RestrictByNotNullOrUndefined().CollapseUnion()
		_ = argObjectType_2
		var paramTypeMap *pkg_rhino.TemplateTypeMap = paramType.GetTemplateTypeMap()
		_ = paramTypeMap
		var paramMapSize int32 = paramTypeMap.Size()
		_ = paramMapSize
		var argTypeMap *pkg_rhino.TemplateTypeMap = argObjectType_2.GetTemplateTypeMap()
		_ = argTypeMap
		{
			var index int32 = (paramMapSize - keyCount)
			_ = index
			for ; index < paramMapSize; index = index + 1 {
				var key *pkg_rhino.TemplateType = paramTypeMap.GetTemplateKeyAt(index)
				_ = key
				this_.priv_InvocationTemplateTypeMatcher_matchTemplateTypesRecursive_3454(paramTypeMap.GetResolvedTemplateType(key), argTypeMap.GetResolvedTemplateType(key))
			}
		}
	}
}

// Java method: InvocationTemplateTypeMatcher.matchTemplateTypesFromNodes(com.google.javascript.rhino.jstype.FunctionType.Parameter[] declParams, int declStartIdx, boolean isTaggedTemplate)
func (this_ *InvocationTemplateTypeMatcher) priv_InvocationTemplateTypeMatcher_matchTemplateTypesFromNodes_3455(declParams []pkg_rhino.FunctionType_Parameter, declStartIdx int32, isTaggedTemplate bool) {
	_ = declParams
	_ = declStartIdx
	_ = isTaggedTemplate
	var declLen int32 = int32(len(declParams))
	_ = declLen
	var declIdx int32 = declStartIdx
	_ = declIdx
	var argCursor pkg_rhino.Node = NodeUtil_FirstInvocationArgCursor(this_.invocation, isTaggedTemplate)
	_ = argCursor
	for (declIdx < declLen) && (argCursor != 0) {
		var declParam pkg_rhino.FunctionType_Parameter = declParams[func() int32 { old_2 := declIdx; nw_3 := old_2 + 1; declIdx = nw_3; return old_2 }()]
		_ = declParam
		var arg pkg_rhino.Node = NodeUtil_InvocationArgFromCursor(argCursor, isTaggedTemplate)
		_ = arg
		argCursor = NodeUtil_NextInvocationArgCursor(argCursor, isTaggedTemplate)
		var declParamType pkg_rhino.JSType = declParam.GetJSType()
		_ = declParamType
		this_.priv_InvocationTemplateTypeMatcher_matchTemplateTypesRecursive_3454(declParamType, this_.priv_InvocationTemplateTypeMatcher_getTypeOrUnknown_3459(arg))
		if declParam.IsVariadic() {
			for argCursor != 0 {
				var restArg pkg_rhino.Node = NodeUtil_InvocationArgFromCursor(argCursor, isTaggedTemplate)
				_ = restArg
				argCursor = NodeUtil_NextInvocationArgCursor(argCursor, isTaggedTemplate)
				this_.priv_InvocationTemplateTypeMatcher_matchTemplateTypesRecursive_3454(declParamType, this_.priv_InvocationTemplateTypeMatcher_getTypeOrUnknown_3459(restArg))
			}
		}
	}
}

// Java method: InvocationTemplateTypeMatcher.matchTemplateTypesFromParameters(com.google.javascript.rhino.jstype.FunctionType.Parameter[] declParams, com.google.javascript.rhino.jstype.FunctionType.Parameter[] callParams)
func (this_ *InvocationTemplateTypeMatcher) priv_InvocationTemplateTypeMatcher_matchTemplateTypesFromParameters_3456(declParams []pkg_rhino.FunctionType_Parameter, callParams []pkg_rhino.FunctionType_Parameter) {
	_ = declParams
	_ = callParams
	var declLen int32 = int32(len(declParams))
	_ = declLen
	var callLen int32 = int32(len(callParams))
	_ = callLen
	var declIdx int32 = int32(0)
	_ = declIdx
	var callIdx int32 = int32(0)
	_ = callIdx
	for (declIdx < declLen) && (callIdx < callLen) {
		var declParam pkg_rhino.FunctionType_Parameter = declParams[func() int32 { old_4 := declIdx; nw_5 := old_4 + 1; declIdx = nw_5; return old_4 }()]
		_ = declParam
		var declParamType pkg_rhino.JSType = declParam.GetJSType()
		_ = declParamType
		this_.priv_InvocationTemplateTypeMatcher_matchTemplateTypesRecursive_3454(declParamType, callParams[func() int32 { old_6 := callIdx; nw_7 := old_6 + 1; callIdx = nw_7; return old_6 }()].GetJSType())
		if declParam.IsVariadic() {
			for callIdx < callLen {
				this_.priv_InvocationTemplateTypeMatcher_matchTemplateTypesRecursive_3454(declParamType, callParams[func() int32 { old_8 := callIdx; nw_9 := old_8 + 1; callIdx = nw_9; return old_8 }()].GetJSType())
			}
		}
	}
}

// Java method: InvocationTemplateTypeMatcher.recordTemplateMatch(com.google.javascript.rhino.jstype.TemplateType template, com.google.javascript.rhino.jstype.JSType match)
func (this_ *InvocationTemplateTypeMatcher) priv_InvocationTemplateTypeMatcher_recordTemplateMatch_3457(template *pkg_rhino.TemplateType, match pkg_rhino.JSType) {
	_ = template
	_ = match
	if match.GetJSTypeBase().IsUnknownType() {
		return
	}
	if this_.matchedTypes != nil {
		this_.matchedTypes.Merge(template, match, lambdaSingleton_invocation_template_type_matcher_1)
		return
	}
	if this_.matchedKey0 == nil {
		this_.matchedKey0 = template
		this_.matchedVal0 = match
		return
	}
	if this_.matchedKey0 == template {
		this_.matchedVal0 = this_.matchedVal0.GetLeastSupertype(match)
		return
	}
	var map_ *pkg_base.LinkedIdentityHashMap[*pkg_rhino.TemplateType, pkg_rhino.JSType] = pkg_base.NewLinkedIdentityHashMap[*pkg_rhino.TemplateType, pkg_rhino.JSType]()
	_ = map_
	map_.Put(this_.matchedKey0, this_.matchedVal0)
	map_.Put(template, match)
	this_.matchedKey0 = nil
	this_.matchedVal0 = nil
	this_.matchedTypes = map_
}

// Java method: InvocationTemplateTypeMatcher.copyMatchedTypes()
func (this_ *InvocationTemplateTypeMatcher) priv_InvocationTemplateTypeMatcher_copyMatchedTypes_3458() j2gort.Ext_com_google_common_collect_ImmutableMap {
	if this_.matchedTypes != nil {
		var builder j2gort.Ext_com_google_common_collect_ImmutableMap_Builder = j2gort.Ext_com_google_common_collect_ImmutableMap_S_Builder()
		_ = builder
		var iifeArg_1 j2gort.
			Ext_java_util_function_BiConsumer
		{

			bound_10 := builder
			iifeArg_1 = j2gort.NewLambda_Ext_java_util_function_BiConsumer(func(mrp0_ any, mrp1_ any) { bound_10.Put_2(mrp0_, mrp1_) })
		}
		this_.matchedTypes.ForEach(iifeArg_1)
		return j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap](builder.BuildOrThrow())
	}
	if this_.matchedKey0 != nil {
		return j2gort.Ext_com_google_common_collect_ImmutableMap_Of_2(j2gort.PtrToAny(this_.matchedKey0), this_.matchedVal0)
	}
	return j2gort.Ext_com_google_common_collect_ImmutableMap_Of_0()
}

// Java method: InvocationTemplateTypeMatcher.getTypeOrUnknown(com.google.javascript.rhino.Node n)
func (this_ *InvocationTemplateTypeMatcher) priv_InvocationTemplateTypeMatcher_getTypeOrUnknown_3459(n pkg_rhino.Node) pkg_rhino.JSType {
	_ = n
	var type_ pkg_rhino.JSType = n.GetJSType()
	_ = type_
	if type_ == nil || type_.IsJavaNil_() {
		return this_.unknownType
	} else {
		return type_
	}
}
