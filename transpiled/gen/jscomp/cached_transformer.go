/*
 * Copyright 2017 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/bundle/CachedTransformer.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.bundle.CachedTransformer
type CachedTransformer struct {
	cache j2gort.Ext_com_google_common_cache_LoadingCache `gson:"cache"`
	idHash_ atomic.Int32
}

func (this_ *CachedTransformer) BindSelf_(s any) {
	_ = s
}
func (this_ *CachedTransformer) IsJavaNil_() bool { return this_ == nil }
func (*CachedTransformer) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.bundle.CachedTransformer", "CachedTransformer" }
func (this_ *CachedTransformer) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CachedTransformer) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.bundle.CachedTransformer", "com.google.javascript.jscomp.bundle.CachedTransformer", func(v_ any) bool {
		return j2gort.InstanceOf[*CachedTransformer](v_)
	})
}
func (*CachedTransformer) Implements_jscomp_Source_Transformer() {}

func StaticGet_CachedTransformer_DEFAULT_CACHE_SPEC() string { return "maximumSize=10000" }
func (this_ *CachedTransformer) initFields_CachedTransformer() {
	_ = this_
}

// Java constructor: CachedTransformer(com.google.javascript.jscomp.bundle.Source.Transformer delegate, com.google.common.cache.CacheBuilder<? super com.google.javascript.jscomp.bundle.Source,? super com.google.javascript.jscomp.bundle.Source> builder)
func NewCachedTransformer_2_Source_Transformer_Ext_com_google_common_cache_CacheBuilder(delegate Source_Transformer, builder *j2gort.Ext_com_google_common_cache_CacheBuilder) *CachedTransformer {
	this_ := &CachedTransformer{}
	InitCachedTransformer_2_Source_Transformer_Ext_com_google_common_cache_CacheBuilder(this_, delegate, builder)
	return this_
}

func InitCachedTransformer_2_Source_Transformer_Ext_com_google_common_cache_CacheBuilder(this_ *CachedTransformer, delegate Source_Transformer, builder *j2gort.Ext_com_google_common_cache_CacheBuilder) {
	_ = delegate
	_ = builder
	this_.initFields_CachedTransformer()
	this_.cache = builder.Build_1(j2gort.Ext_com_google_common_cache_CacheLoader_From_1_Ext_com_google_common_base_Function(func() j2gort.Ext_com_google_common_base_Function { bound_1 := delegate; return j2gort.NewLambda_Ext_com_google_common_base_Function(func(mrp0_ any) any { return bound_1.Transform(j2gort.AsIface[Source](mrp0_)) }) }()))
}

// Java constructor: CachedTransformer(com.google.javascript.jscomp.bundle.Source.Transformer delegate, java.lang.String spec)
func NewCachedTransformer_2_Source_Transformer_string(delegate Source_Transformer, spec string) *CachedTransformer {
	this_ := &CachedTransformer{}
	InitCachedTransformer_2_Source_Transformer_string(this_, delegate, spec)
	return this_
}

func InitCachedTransformer_2_Source_Transformer_string(this_ *CachedTransformer, delegate Source_Transformer, spec string) {
	_ = delegate
	_ = spec
	InitCachedTransformer_2_Source_Transformer_Ext_com_google_common_cache_CacheBuilder(this_, delegate, j2gort.Ext_com_google_common_cache_CacheBuilder_From_1_string(spec))
}

// Java constructor: CachedTransformer(com.google.javascript.jscomp.bundle.Source.Transformer delegate)
func NewCachedTransformer_1(delegate Source_Transformer) *CachedTransformer {
	this_ := &CachedTransformer{}
	InitCachedTransformer_1(this_, delegate)
	return this_
}

func InitCachedTransformer_1(this_ *CachedTransformer, delegate Source_Transformer) {
	_ = delegate
	InitCachedTransformer_2_Source_Transformer_string(this_, delegate, "maximumSize=10000")
}

// Java method: CachedTransformer.transform(com.google.javascript.jscomp.bundle.Source input)
func (this_ *CachedTransformer) Transform(input Source) Source {
	_ = input
	return j2gort.AsIface[Source](this_.cache.GetUnchecked(input))
}

func (this_ *CachedTransformer) AsFunction() j2gort.Ext_java_util_function_Function {
	return IfaceDefault_Source_Transformer_AsFunction(this_)
}

func (this_ *CachedTransformer) AndThen(p0_ Source_Transformer) Source_Transformer {
	return IfaceDefault_Source_Transformer_AndThen(this_, p0_)
}

func (this_ *CachedTransformer) Compose(p0_ Source_Transformer) Source_Transformer {
	return IfaceDefault_Source_Transformer_Compose(this_, p0_)
}

