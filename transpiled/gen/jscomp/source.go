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
// Port of com/google/javascript/jscomp/bundle/Source.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_source_1 Source_Transformer

func init() {
	lambdaSingleton_source_1 = NewLambda_Source_Transformer(func(x Source) Source {
		_ = x
		return x
	})
}

// Java polymorphic class interface: com.google.javascript.jscomp.bundle.Source
type Source interface {
	Implements_jscomp_Source()
	IsJavaNil_() bool
	GetSourceBase() *SourceBase
	Path() j2gort.Ext_java_nio_file_Path
	SourceMap() string
	SourceUrl() string
	SourceMappingUrl() string
	Runtimes() j2gort.Ext_com_google_common_collect_ImmutableSet
	LoadFlags() j2gort.Ext_com_google_common_collect_ImmutableMap
	EstimatedSize() int32
	Code() string
	OriginalCode() string
	ToBuilder() Source_Builder
	CodeSupplier() Source_Lazy
	OriginalCodeSupplier() Source_Lazy
}

func ToIface_Source[T interface {
	Source
	comparable
}](v T) Source {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_Source[T Source](v T) Source {
	var vi Source = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.jscomp.bundle.Source
type SourceBase struct {
	self_   Source
	idHash_ atomic.Int32
}

func (this_ *SourceBase) BindSelf_(s any) {
	_ = s
	this_.self_ = s.(Source)
}
func (this_ *SourceBase) IsJavaNil_() bool { return this_ == nil }
func (*SourceBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.bundle.Source", "Source"
}
func (this_ *SourceBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.bundle.Source", "com.google.javascript.jscomp.bundle.Source", func(v_ any) bool {
		return j2gort.InstanceOf[Source](v_)
	})
}
func (b *SourceBase) GetSourceBase() *SourceBase { return b }
func (*SourceBase) Implements_jscomp_Source()    {}

var static_Source_DEV_NULL j2gort.Ext_java_nio_file_Path

func StaticGet_Source_DEV_NULL() j2gort.Ext_java_nio_file_Path {
	clinit_Source()
	return static_Source_DEV_NULL
}
func StaticSet_Source_DEV_NULL(v j2gort.Ext_java_nio_file_Path) {
	clinit_Source()
	static_Source_DEV_NULL = v
}

var clinitDone_Source uint32
var clinitState_Source j2gort.ClinitState

func clinit_Source() {
	if atomic.LoadUint32(&clinitDone_Source) != 0 {
		return
	}
	clinitSlow_Source()
}

func clinitSlow_Source() {
	if !j2gort.EnterClinit(&clinitState_Source, &clinitDone_Source) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_Source, &clinitDone_Source)
	static_Source_DEV_NULL = j2gort.Ext_java_nio_file_Path_Of_2("/dev/null", []string{})
}

func (this_ *SourceBase) initFields_Source() {
	_ = this_
}

func InitSource(this_ *SourceBase) {
	this_.initFields_Source()
}

// Java method: Source.path()
func (this_ *SourceBase) Path() j2gort.Ext_java_nio_file_Path {
	panic("abstract or native method invoked")
}

// Java method: Source.sourceMap()
func (this_ *SourceBase) SourceMap() string {
	panic("abstract or native method invoked")
}

// Java method: Source.sourceUrl()
func (this_ *SourceBase) SourceUrl() string {
	panic("abstract or native method invoked")
}

// Java method: Source.sourceMappingUrl()
func (this_ *SourceBase) SourceMappingUrl() string {
	panic("abstract or native method invoked")
}

// Java method: Source.runtimes()
func (this_ *SourceBase) Runtimes() j2gort.Ext_com_google_common_collect_ImmutableSet {
	panic("abstract or native method invoked")
}

// Java method: Source.loadFlags()
func (this_ *SourceBase) LoadFlags() j2gort.Ext_com_google_common_collect_ImmutableMap {
	panic("abstract or native method invoked")
}

// Java method: Source.estimatedSize()
func (this_ *SourceBase) EstimatedSize() int32 {
	panic("abstract or native method invoked")
}

// Java method: Source.code()
func (this_ *SourceBase) Code() string {
	return j2gort.CastString(this_.self_.CodeSupplier().Get())
}

// Java method: Source.originalCode()
func (this_ *SourceBase) OriginalCode() string {
	return j2gort.CastString(this_.self_.OriginalCodeSupplier().Get())
}

// Java method: Source.toBuilder()
func (this_ *SourceBase) ToBuilder() Source_Builder {
	panic("abstract or native method invoked")
}

// Java static method: Source.builder()
func Source_S_Builder() Source_Builder {
	clinit_Source()
	return NewAutoValue_Source_Builder().SetPath(StaticGet_Source_DEV_NULL()).GetSource_BuilderBase().SetCode_1_string("").SetOriginalCodeSupplier(nil).SetSourceMap("").SetSourceUrl("").SetSourceMappingUrl("").SetRuntimes(j2gort.Ext_com_google_common_collect_ImmutableSet_Of_0()).SetLoadFlags(j2gort.Ext_com_google_common_collect_ImmutableMap_Of_0()).SetEstimatedSize(int32(0))
}

// Java method: Source.codeSupplier()
func (this_ *SourceBase) CodeSupplier() Source_Lazy {
	panic("abstract or native method invoked")
}

// Java method: Source.originalCodeSupplier()
func (this_ *SourceBase) OriginalCodeSupplier() Source_Lazy {
	panic("abstract or native method invoked")
}

// Java polymorphic class interface: com.google.javascript.jscomp.bundle.Source.Builder
type Source_Builder interface {
	Implements_jscomp_Source_Builder()
	IsJavaNil_() bool
	GetSource_BuilderBase() *Source_BuilderBase
	SetPath(j2gort.Ext_java_nio_file_Path) Source_Builder
	SetSourceMap(string) Source_Builder
	SetSourceUrl(string) Source_Builder
	SetSourceMappingUrl(string) Source_Builder
	SetRuntimes(j2gort.Ext_com_google_common_collect_ImmutableSet) Source_Builder
	SetLoadFlags(j2gort.Ext_com_google_common_collect_ImmutableMap) Source_Builder
	SetEstimatedSize(int32) Source_Builder
	SetCode_1_Ext_com_google_common_base_Supplier(j2gort.Ext_com_google_common_base_Supplier) Source_Builder
	SetCode_1_string(string) Source_Builder
	SetOriginalCode(string) Source_Builder
	AddRuntime([]string) Source_Builder
	SetDependencyInfo(DependencyInfo) Source_Builder
	Build() Source
	SetCodeSupplier(Source_Lazy) Source_Builder
	SetOriginalCodeSupplier(Source_Lazy) Source_Builder
	Runtimes() j2gort.Ext_com_google_common_collect_ImmutableSet
	CodeSupplier() Source_Lazy
	AutoBuild() Source
	OriginalCodeSupplier() Source_Lazy
}

func ToIface_Source_Builder[T interface {
	Source_Builder
	comparable
}](v T) Source_Builder {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_Source_Builder[T Source_Builder](v T) Source_Builder {
	var vi Source_Builder = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.jscomp.bundle.Source.Builder
type Source_BuilderBase struct {
	self_   Source_Builder
	idHash_ atomic.Int32
}

func (this_ *Source_BuilderBase) BindSelf_(s any) {
	_ = s
	this_.self_ = s.(Source_Builder)
}
func (this_ *Source_BuilderBase) IsJavaNil_() bool { return this_ == nil }
func (*Source_BuilderBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.bundle.Source$Builder", "Builder"
}
func (this_ *Source_BuilderBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.bundle.Source.Builder", "com.google.javascript.jscomp.bundle.Source$Builder", func(v_ any) bool {
		return j2gort.InstanceOf[Source_Builder](v_)
	})
}
func (b *Source_BuilderBase) GetSource_BuilderBase() *Source_BuilderBase { return b }
func (*Source_BuilderBase) Implements_jscomp_Source_Builder()            {}

func (this_ *Source_BuilderBase) initFields_Source_Builder() {
	_ = this_
}

func InitSource_Builder(this_ *Source_BuilderBase) {
	this_.initFields_Source_Builder()
}

// Java method: Builder.setPath(java.nio.file.Path path)
func (this_ *Source_BuilderBase) SetPath(path j2gort.Ext_java_nio_file_Path) Source_Builder {
	_ = path
	panic("abstract or native method invoked")
}

// Java method: Builder.setSourceMap(java.lang.String sourceMap)
func (this_ *Source_BuilderBase) SetSourceMap(sourceMap string) Source_Builder {
	_ = sourceMap
	panic("abstract or native method invoked")
}

// Java method: Builder.setSourceUrl(java.lang.String sourceUrl)
func (this_ *Source_BuilderBase) SetSourceUrl(sourceUrl string) Source_Builder {
	_ = sourceUrl
	panic("abstract or native method invoked")
}

// Java method: Builder.setSourceMappingUrl(java.lang.String sourceMappingUrl)
func (this_ *Source_BuilderBase) SetSourceMappingUrl(sourceMappingUrl string) Source_Builder {
	_ = sourceMappingUrl
	panic("abstract or native method invoked")
}

// Java method: Builder.setRuntimes(com.google.common.collect.ImmutableSet<java.lang.String> runtimes)
func (this_ *Source_BuilderBase) SetRuntimes(runtimes j2gort.Ext_com_google_common_collect_ImmutableSet) Source_Builder {
	_ = runtimes
	panic("abstract or native method invoked")
}

// Java method: Builder.setLoadFlags(com.google.common.collect.ImmutableMap<java.lang.String,java.lang.String> flags)
func (this_ *Source_BuilderBase) SetLoadFlags(flags j2gort.Ext_com_google_common_collect_ImmutableMap) Source_Builder {
	_ = flags
	panic("abstract or native method invoked")
}

// Java method: Builder.setEstimatedSize(int estimatedSize)
func (this_ *Source_BuilderBase) SetEstimatedSize(estimatedSize int32) Source_Builder {
	_ = estimatedSize
	panic("abstract or native method invoked")
}

// Java method: Builder.setCode(com.google.common.base.Supplier<java.lang.String> code)
func (this_ *Source_BuilderBase) SetCode_1_Ext_com_google_common_base_Supplier(code j2gort.Ext_com_google_common_base_Supplier) Source_Builder {
	_ = code
	return this_.self_.SetCodeSupplier(Source_Lazy_Memoize(code))
}

// Java method: Builder.setCode(java.lang.String code)
func (this_ *Source_BuilderBase) SetCode_1_string(code string) Source_Builder {
	_ = code
	return this_.self_.SetCodeSupplier(Source_Lazy_OfInstance(j2gort.BoxString(code)))
}

// Java method: Builder.setOriginalCode(java.lang.String code)
func (this_ *Source_BuilderBase) SetOriginalCode(code string) Source_Builder {
	_ = code
	return this_.self_.SetOriginalCodeSupplier(Source_Lazy_OfInstance(j2gort.BoxString(code)))
}

// Java method: Builder.addRuntime(java.lang.String[] runtimes)
func (this_ *Source_BuilderBase) AddRuntime(runtimes []string) Source_Builder {
	_ = runtimes
	return this_.self_.SetRuntimes(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet_Builder](j2gort.Ext_com_google_common_collect_ImmutableSet_S_Builder().AddAll_1_Ext_java_lang_Iterable(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.self_.Runtimes()))).Add_1_slice_any(func() []any {
		arr_1 := runtimes
		if arr_1 == nil {
			return nil
		}
		out_2 := make([]any, len(arr_1))
		for i_, el_3 := range arr_1 {
			out_2[i_] = j2gort.BoxString(el_3)
		}
		return out_2
	}())).Build()))
}

// Java method: Builder.setDependencyInfo(com.google.javascript.jscomp.deps.DependencyInfo info)
func (this_ *Source_BuilderBase) SetDependencyInfo(info DependencyInfo) Source_Builder {
	_ = info
	return this_.self_.SetLoadFlags(info.GetLoadFlags())
}

// Java method: Builder.build()
func (this_ *Source_BuilderBase) Build() Source {
	if j2gort.IfaceEquals[Source_Lazy](this_.self_.OriginalCodeSupplier(), nil) {
		this_.self_.SetOriginalCodeSupplier(this_.self_.CodeSupplier())
	}
	return this_.self_.AutoBuild()
}

// Java method: Builder.setCodeSupplier(com.google.javascript.jscomp.bundle.Source.Lazy<java.lang.String> code)
func (this_ *Source_BuilderBase) SetCodeSupplier(code Source_Lazy) Source_Builder {
	_ = code
	panic("abstract or native method invoked")
}

// Java method: Builder.setOriginalCodeSupplier(com.google.javascript.jscomp.bundle.Source.@org.jspecify.annotations.Nullable Lazy<java.lang.String> code)
func (this_ *Source_BuilderBase) SetOriginalCodeSupplier(code Source_Lazy) Source_Builder {
	_ = code
	panic("abstract or native method invoked")
}

// Java method: Builder.runtimes()
func (this_ *Source_BuilderBase) Runtimes() j2gort.Ext_com_google_common_collect_ImmutableSet {
	panic("abstract or native method invoked")
}

// Java method: Builder.codeSupplier()
func (this_ *Source_BuilderBase) CodeSupplier() Source_Lazy {
	panic("abstract or native method invoked")
}

// Java method: Builder.autoBuild()
func (this_ *Source_BuilderBase) AutoBuild() Source {
	panic("abstract or native method invoked")
}

// Java method: Builder.originalCodeSupplier()
func (this_ *Source_BuilderBase) OriginalCodeSupplier() Source_Lazy {
	panic("abstract or native method invoked")
}

// Java interface: com.google.javascript.jscomp.bundle.Source.Transformer
type Source_Transformer interface {
	Implements_jscomp_Source_Transformer()
	IsJavaNil_() bool
	Transform(Source) Source
	AsFunction() j2gort.Ext_java_util_function_Function
	AndThen(Source_Transformer) Source_Transformer
	Compose(Source_Transformer) Source_Transformer
}

func ToIface_Source_Transformer[T interface {
	Source_Transformer
	comparable
}](v T) Source_Transformer {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_Source_Transformer[T Source_Transformer](v T) Source_Transformer {
	var vi Source_Transformer = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.bundle.Source.Transformer", "com.google.javascript.jscomp.bundle.Source$Transformer", func(v_ any) bool {
		return j2gort.InstanceOf[Source_Transformer](v_)
	})
}

// Java static method: Transformer.of(java.util.function.Function<com.google.javascript.jscomp.bundle.Source,com.google.javascript.jscomp.bundle.Source> function)
func Source_Transformer_Of(function j2gort.Ext_java_util_function_Function) Source_Transformer {
	_ = function
	{

		bound_4 := function
		return NewLambda_Source_Transformer(func(mrp0_ Source) Source { return j2gort.AsIface[Source](bound_4.Apply(mrp0_)) })
	}

}

// Java static method: Transformer.identity()
func Source_Transformer_Identity() Source_Transformer {
	return lambdaSingleton_source_1
}

// Java default interface method: Transformer.asFunction()
func IfaceDefault_Source_Transformer_AsFunction(this_ Source_Transformer) j2gort.Ext_java_util_function_Function {
	_ = this_
	{

		bound_5 := this_
		return j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any { return bound_5.Transform(j2gort.AsIface[Source](mrp0_)) })
	}

}

// Java default interface method: Transformer.andThen(com.google.javascript.jscomp.bundle.Source.Transformer after)
func IfaceDefault_Source_Transformer_AndThen(this_ Source_Transformer, after Source_Transformer) Source_Transformer {
	_ = this_
	_ = after
	var before Source_Transformer = this_
	_ = before
	return NewLambda_Source_Transformer(func(x Source) Source {
		_ = x
		return after.Transform(before.Transform(x))
	})
}

// Java default interface method: Transformer.compose(com.google.javascript.jscomp.bundle.Source.Transformer before)
func IfaceDefault_Source_Transformer_Compose(this_ Source_Transformer, before Source_Transformer) Source_Transformer {
	_ = this_
	_ = before
	var after Source_Transformer = this_
	_ = after
	return NewLambda_Source_Transformer(func(x Source) Source {
		_ = x
		return after.Transform(before.Transform(x))
	})
}

type lambda_Source_Transformer struct {
	fn func(Source) Source
}

func (l *lambda_Source_Transformer) IsJavaNil_() bool { return l == nil }
func NewLambda_Source_Transformer(fn func(Source) Source) Source_Transformer {
	if fn == nil {
		return nil
	}
	return &lambda_Source_Transformer{fn: fn}
}

func (*lambda_Source_Transformer) Implements_jscomp_Source_Transformer() {}
func (l *lambda_Source_Transformer) Transform(p0_ Source) Source {
	return l.fn(p0_)
}

func (this_ *lambda_Source_Transformer) AsFunction() j2gort.Ext_java_util_function_Function {
	return IfaceDefault_Source_Transformer_AsFunction(this_)
}

func (this_ *lambda_Source_Transformer) AndThen(p0_ Source_Transformer) Source_Transformer {
	return IfaceDefault_Source_Transformer_AndThen(this_, p0_)
}

func (this_ *lambda_Source_Transformer) Compose(p0_ Source_Transformer) Source_Transformer {
	return IfaceDefault_Source_Transformer_Compose(this_, p0_)
}

// Java polymorphic class interface: com.google.javascript.jscomp.bundle.Source.Lazy
type Source_Lazy interface {
	Implements_jscomp_Source_Lazy()
	IsJavaNil_() bool
	GetSource_LazyBase() *Source_LazyBase
	j2gort.Ext_com_google_common_base_Supplier
	Equals(any) bool
	HashCode() int32
}

func ToIface_Source_Lazy[T interface {
	Source_Lazy
	comparable
}](v T) Source_Lazy {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_Source_Lazy[T Source_Lazy](v T) Source_Lazy {
	var vi Source_Lazy = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.jscomp.bundle.Source.Lazy
type Source_LazyBase struct {
	self_ Source_Lazy
}

func (this_ *Source_LazyBase) BindSelf_(s any) {
	_ = s
	this_.self_ = s.(Source_Lazy)
}
func (this_ *Source_LazyBase) IsJavaNil_() bool { return this_ == nil }
func (*Source_LazyBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.bundle.Source$Lazy", "Lazy"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.bundle.Source.Lazy", "com.google.javascript.jscomp.bundle.Source$Lazy", func(v_ any) bool {
		return j2gort.InstanceOf[Source_Lazy](v_)
	})
}
func (b *Source_LazyBase) GetSource_LazyBase() *Source_LazyBase                 { return b }
func (*Source_LazyBase) Implements_jscomp_Source_Lazy()                         {}
func (*Source_LazyBase) Implements_j2gort_Ext_com_google_common_base_Supplier() {}
func (*Source_LazyBase) Implements_j2gort_Ext_java_util_function_Supplier()     {}

func (this_ *Source_LazyBase) initFields_Source_Lazy() {
	_ = this_
}

func InitSource_Lazy(this_ *Source_LazyBase) {
	this_.initFields_Source_Lazy()
}

// Java method: Lazy.equals(java.lang.Object other)
func (this_ *Source_LazyBase) Equals(arg_6 any) bool {
	other := arg_6
	_ = other
	return (j2gort.InstanceOf[Source_Lazy](other) && javalang.ObjectsEquals(this_.self_.Get(), (j2gort.AsIface[Source_Lazy](other)).Get()))
}

// Java method: Lazy.hashCode()
func (this_ *Source_LazyBase) HashCode() int32 {
	return j2gort.ObjectHash(this_.self_.Get())
}

// Java static method: Lazy.ofInstance(T instance)
func Source_Lazy_OfInstance(instance any) Source_Lazy {
	_ = instance
	return NewSource__Anon_100(instance)
}

// Java static method: Lazy.memoize(com.google.common.base.Supplier<T> supplier)
func Source_Lazy_Memoize(supplier j2gort.Ext_com_google_common_base_Supplier) Source_Lazy {
	_ = supplier
	var memoized j2gort.Ext_com_google_common_base_Supplier = j2gort.Ext_com_google_common_base_Suppliers_Memoize(supplier)
	_ = memoized
	return NewSource__Anon_101(memoized)
}

// Java class:
type Source__Anon_100 struct {
	Source_LazyBase
	cap_0_instance any
}

func (this_ *Source__Anon_100) BindSelf_(s any) {
	_ = s
	this_.Source_LazyBase.BindSelf_(s)
}
func (this_ *Source__Anon_100) IsJavaNil_() bool { return this_ == nil }
func (*Source__Anon_100) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.bundle.Source$Lazy$1", ""
}
func (*Source__Anon_100) Implements_jscomp_Source_Lazy()                         {}
func (*Source__Anon_100) Implements_j2gort_Ext_com_google_common_base_Supplier() {}
func (*Source__Anon_100) Implements_j2gort_Ext_java_util_function_Supplier()     {}

func (this_ *Source__Anon_100) initFields_Source__Anon_100() {
	_ = this_
}

func NewSource__Anon_100(cap_p_instance any) *Source__Anon_100 {
	this_ := &Source__Anon_100{}
	this_.BindSelf_(this_)
	this_.cap_0_instance = cap_p_instance
	InitSource_Lazy(&this_.Source_LazyBase)
	this_.initFields_Source__Anon_100()
	return this_
}

// Java method: .get()
func (this_ *Source__Anon_100) Get() any {
	return this_.cap_0_instance
}

// Java class:
type Source__Anon_101 struct {
	Source_LazyBase
	cap_0_memoized j2gort.Ext_com_google_common_base_Supplier
}

func (this_ *Source__Anon_101) BindSelf_(s any) {
	_ = s
	this_.Source_LazyBase.BindSelf_(s)
}
func (this_ *Source__Anon_101) IsJavaNil_() bool { return this_ == nil }
func (*Source__Anon_101) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.bundle.Source$Lazy$2", ""
}
func (*Source__Anon_101) Implements_jscomp_Source_Lazy()                         {}
func (*Source__Anon_101) Implements_j2gort_Ext_com_google_common_base_Supplier() {}
func (*Source__Anon_101) Implements_j2gort_Ext_java_util_function_Supplier()     {}

func (this_ *Source__Anon_101) initFields_Source__Anon_101() {
	_ = this_
}

func NewSource__Anon_101(cap_p_memoized j2gort.Ext_com_google_common_base_Supplier) *Source__Anon_101 {
	this_ := &Source__Anon_101{}
	this_.BindSelf_(this_)
	this_.cap_0_memoized = cap_p_memoized
	InitSource_Lazy(&this_.Source_LazyBase)
	this_.initFields_Source__Anon_101()
	return this_
}

// Java method: .get()
func (this_ *Source__Anon_101) Get() any {
	return this_.cap_0_memoized.Get()
}
