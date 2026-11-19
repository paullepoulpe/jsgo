/*
 * Copyright 2018 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/Promises.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.Promises
type Promises struct {
	idHash_ atomic.Int32
}

func (this_ *Promises) BindSelf_(s any) {
	_ = s
}
func (this_ *Promises) IsJavaNil_() bool { return this_ == nil }
func (*Promises) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.Promises", "Promises" }
func (this_ *Promises) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Promises) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Promises", "com.google.javascript.jscomp.Promises", func(v_ any) bool {
		return j2gort.InstanceOf[*Promises](v_)
	})
}

func (this_ *Promises) initFields_Promises() {
	_ = this_
}

// Java constructor: Promises()
func NewPromises() *Promises {
	this_ := &Promises{}
	InitPromises(this_)
	return this_
}

func InitPromises(this_ *Promises) {
	this_.initFields_Promises()
}

// Java static method: Promises.getTemplateTypeOfThenable(com.google.javascript.rhino.jstype.JSTypeRegistry registry, com.google.javascript.rhino.jstype.JSType maybeThenable)
func Promises_GetTemplateTypeOfThenable(registry *pkg_rhino.JSTypeRegistry, maybeThenable pkg_rhino.JSType) pkg_rhino.JSType {
	_ = registry
	_ = maybeThenable
	var templateType *pkg_rhino.TemplateType = registry.GetIThenableTemplate()
	_ = templateType
	return maybeThenable.RestrictByNotNullOrUndefined().GetTemplateTypeMap().GetResolvedTemplateType(templateType)
}

// Java static method: Promises.getResolvedType(com.google.javascript.rhino.jstype.JSTypeRegistry registry, com.google.javascript.rhino.jstype.JSType type)
func Promises_GetResolvedType(registry *pkg_rhino.JSTypeRegistry, type_ pkg_rhino.JSType) pkg_rhino.JSType {
	_ = registry
	_ = type_
	if (type_.IsUnknownType()) {
		return type_
	}
	if (type_.IsUnionType()) {
		var union *pkg_rhino.UnionType = type_.ToMaybeUnionType()
		_ = union
		var n int32 = union.GetAlternatesCount()
		_ = n
		var unionTypeBuilder *pkg_rhino.UnionType_Builder = registry.AcquireUnionTypeBuilder()
		_ = unionTypeBuilder
		defer func() {
			registry.ReleaseUnionTypeBuilder(unionTypeBuilder)
		}()
		{
			{
				var i int32 = int32(0)
				_ = i
				for ; (i < n); i = i + 1 {
					unionTypeBuilder.AddAlternate(Promises_GetResolvedType(registry, union.GetAlternateAt(i)))
				}
			}
			return unionTypeBuilder.Build()
		}
	}
	var templates *pkg_rhino.TemplateTypeMap = type_.GetTemplateTypeMap()
	_ = templates
	if (templates.HasTemplateKey_1(registry.GetIThenableTemplate())) {
		return Promises_GetResolvedType(registry, templates.GetResolvedTemplateType(registry.GetIThenableTemplate()))
	}
	if (type_.IsSubtypeOf_1(registry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_THENABLE_TYPE()))) {
		return registry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_UNKNOWN_TYPE())
	}
	return type_
}

// Java static method: Promises.wrapInIThenable(com.google.javascript.rhino.jstype.JSTypeRegistry registry, com.google.javascript.rhino.jstype.JSType maybeThenable)
func Promises_WrapInIThenable(registry *pkg_rhino.JSTypeRegistry, maybeThenable pkg_rhino.JSType) pkg_rhino.JSType {
	_ = registry
	_ = maybeThenable
	var unwrapped pkg_rhino.JSType = Promises_GetResolvedType(registry, maybeThenable)
	_ = unwrapped
	return pkg_rhino.ToIface_JSType(registry.CreateTemplatizedType_2_ObjectType_JSType(registry.GetNativeObjectType(pkg_rhino.StaticGet_JSTypeNative_I_THENABLE_TYPE()), []pkg_rhino.JSType{unwrapped}))
}

// Java static method: Promises.createAsyncReturnableType(com.google.javascript.rhino.jstype.JSTypeRegistry registry, com.google.javascript.rhino.jstype.JSType maybeThenable)
func Promises_CreateAsyncReturnableType(registry *pkg_rhino.JSTypeRegistry, maybeThenable pkg_rhino.JSType) pkg_rhino.JSType {
	_ = registry
	_ = maybeThenable
	var unknownType pkg_rhino.JSType = registry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_UNKNOWN_TYPE())
	_ = unknownType
	var iThenableType pkg_rhino.ObjectType = registry.GetNativeObjectType(pkg_rhino.StaticGet_JSTypeNative_I_THENABLE_TYPE())
	_ = iThenableType
	var iThenableOfUnknownType pkg_rhino.JSType = pkg_rhino.ToIface_JSType(registry.CreateTemplatizedType_2_ObjectType_JSType(iThenableType, []pkg_rhino.JSType{unknownType}))
	_ = iThenableOfUnknownType
	var asyncTemplateBuilder *pkg_rhino.UnionType_Builder = nil
	_ = asyncTemplateBuilder
	var asyncTemplateUnion pkg_rhino.JSType = nil
	_ = asyncTemplateUnion
	var didRet_1 bool
	var retVal_2 pkg_rhino.JSType
	func() {
		defer func() {
			if ((asyncTemplateBuilder != nil)) {
				registry.ReleaseUnionTypeBuilder(asyncTemplateBuilder)
			}
		}()
		if (maybeThenable.IsUnionType()) {
			var union *pkg_rhino.UnionType = maybeThenable.ToMaybeUnionType()
			_ = union
			{
				var i int32 = int32(0)
				_ = i
				var n int32 = union.GetAlternatesCount()
				_ = n
				for ; (i < n); i = i + 1 {
					var t pkg_rhino.JSType = union.GetAlternateAt(i)
					_ = t
					if (t.IsSubtypeOf_1(iThenableOfUnknownType)) {
						if ((asyncTemplateBuilder == nil)) {
							asyncTemplateBuilder = registry.AcquireUnionTypeBuilder()
						}
						asyncTemplateBuilder.AddAlternate(Promises_GetTemplateTypeOfThenable(registry, t))
					}
				}
			}
		} else if (maybeThenable.IsSubtypeOf_1(iThenableOfUnknownType)) {
			asyncTemplateBuilder = registry.AcquireUnionTypeBuilder()
			asyncTemplateBuilder.AddAlternate(Promises_GetTemplateTypeOfThenable(registry, maybeThenable))
		}
		if ((asyncTemplateBuilder == nil)) {
			retVal_2 = unknownType
			didRet_1 = true
			return
		}
		asyncTemplateUnion = asyncTemplateBuilder.Build()
	}()
	if didRet_1 { return retVal_2 }
	return registry.CreateUnionType_2(asyncTemplateUnion, pkg_rhino.ToIface_JSType(registry.CreateTemplatizedType_2_ObjectType_JSType(iThenableType, []pkg_rhino.JSType{asyncTemplateUnion})))
}

