// Port of com/google/javascript/jscomp/serialization/AutoValue_TypedAstDeserializer_DeserializedAst.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.serialization.AutoValue_TypedAstDeserializer_DeserializedAst
type AutoValue_TypedAstDeserializer_DeserializedAst struct {
	TypedAstDeserializer_DeserializedAstBase
	filesystem       j2gort.Ext_java_util_concurrent_ConcurrentMap     `gson:"filesystem"`
	colorRegistry    j2gort.Ext_com_google_common_base_Optional        `gson:"colorRegistry"`
	externProperties j2gort.Ext_com_google_common_collect_ImmutableSet `gson:"externProperties"`
	runtimeLibraries j2gort.Ext_com_google_common_collect_ImmutableSet `gson:"runtimeLibraries"`
}

func (this_ *AutoValue_TypedAstDeserializer_DeserializedAst) BindSelf_(s any) {
	_ = s
}
func (this_ *AutoValue_TypedAstDeserializer_DeserializedAst) IsJavaNil_() bool { return this_ == nil }
func (*AutoValue_TypedAstDeserializer_DeserializedAst) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.serialization.AutoValue_TypedAstDeserializer_DeserializedAst", "AutoValue_TypedAstDeserializer_DeserializedAst"
}
func (this_ *AutoValue_TypedAstDeserializer_DeserializedAst) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.TypedAstDeserializer_DeserializedAstBase.IdentityHash_()
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.serialization.AutoValue_TypedAstDeserializer_DeserializedAst", "com.google.javascript.jscomp.serialization.AutoValue_TypedAstDeserializer_DeserializedAst", func(v_ any) bool {
		return j2gort.InstanceOf[*AutoValue_TypedAstDeserializer_DeserializedAst](v_)
	})
}
func (*AutoValue_TypedAstDeserializer_DeserializedAst) Implements_jscomp_TypedAstDeserializer_DeserializedAst() {
}

func (this_ *AutoValue_TypedAstDeserializer_DeserializedAst) initFields_AutoValue_TypedAstDeserializer_DeserializedAst() {
	_ = this_
}

// Java constructor: AutoValue_TypedAstDeserializer_DeserializedAst(java.util.concurrent.ConcurrentMap<com.google.javascript.jscomp.SourceFile,java.util.function.Supplier<com.google.javascript.rhino.Node>> filesystem, com.google.common.base.Optional<com.google.javascript.jscomp.colors.ColorRegistry> colorRegistry, com.google.common.collect.@org.jspecify.annotations.Nullable ImmutableSet<java.lang.String> externProperties, com.google.common.collect.ImmutableSet<java.lang.String> runtimeLibraries)
func NewAutoValue_TypedAstDeserializer_DeserializedAst(filesystem j2gort.Ext_java_util_concurrent_ConcurrentMap, colorRegistry j2gort.Ext_com_google_common_base_Optional, externProperties j2gort.Ext_com_google_common_collect_ImmutableSet, runtimeLibraries j2gort.Ext_com_google_common_collect_ImmutableSet) *AutoValue_TypedAstDeserializer_DeserializedAst {
	this_ := &AutoValue_TypedAstDeserializer_DeserializedAst{}
	InitAutoValue_TypedAstDeserializer_DeserializedAst(this_, filesystem, colorRegistry, externProperties, runtimeLibraries)
	return this_
}

func InitAutoValue_TypedAstDeserializer_DeserializedAst(this_ *AutoValue_TypedAstDeserializer_DeserializedAst, filesystem j2gort.Ext_java_util_concurrent_ConcurrentMap, colorRegistry j2gort.Ext_com_google_common_base_Optional, externProperties j2gort.Ext_com_google_common_collect_ImmutableSet, runtimeLibraries j2gort.Ext_com_google_common_collect_ImmutableSet) {
	_ = filesystem
	_ = colorRegistry
	_ = externProperties
	_ = runtimeLibraries
	InitTypedAstDeserializer_DeserializedAst(&this_.TypedAstDeserializer_DeserializedAstBase)
	this_.initFields_AutoValue_TypedAstDeserializer_DeserializedAst()
	if filesystem == nil {
		panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", "Null filesystem")))
	}
	this_.filesystem = filesystem
	if colorRegistry == nil {
		panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", "Null colorRegistry")))
	}
	this_.colorRegistry = colorRegistry
	this_.externProperties = externProperties
	if runtimeLibraries == nil {
		panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", "Null runtimeLibraries")))
	}
	this_.runtimeLibraries = runtimeLibraries
}

// Java method: AutoValue_TypedAstDeserializer_DeserializedAst.getFilesystem()
func (this_ *AutoValue_TypedAstDeserializer_DeserializedAst) GetFilesystem() j2gort.Ext_java_util_concurrent_ConcurrentMap {
	return this_.filesystem
}

// Java method: AutoValue_TypedAstDeserializer_DeserializedAst.getColorRegistry()
func (this_ *AutoValue_TypedAstDeserializer_DeserializedAst) GetColorRegistry() j2gort.Ext_com_google_common_base_Optional {
	return this_.colorRegistry
}

// Java method: AutoValue_TypedAstDeserializer_DeserializedAst.getExternProperties()
func (this_ *AutoValue_TypedAstDeserializer_DeserializedAst) GetExternProperties() j2gort.Ext_com_google_common_collect_ImmutableSet {
	return this_.externProperties
}

// Java method: AutoValue_TypedAstDeserializer_DeserializedAst.getRuntimeLibraries()
func (this_ *AutoValue_TypedAstDeserializer_DeserializedAst) GetRuntimeLibraries() j2gort.Ext_com_google_common_collect_ImmutableSet {
	return this_.runtimeLibraries
}

// Java method: AutoValue_TypedAstDeserializer_DeserializedAst.toString()
func (this_ *AutoValue_TypedAstDeserializer_DeserializedAst) String() string {
	return ("DeserializedAst{filesystem=" + j2gort.StrVal(this_.filesystem) + ", colorRegistry=" + j2gort.StrVal(this_.colorRegistry) + ", externProperties=" + j2gort.StrVal(this_.externProperties) + ", runtimeLibraries=" + j2gort.StrVal(this_.runtimeLibraries) + "}")
}

// Java method: AutoValue_TypedAstDeserializer_DeserializedAst.equals(java.lang.@org.jspecify.annotations.Nullable Object o)
func (this_ *AutoValue_TypedAstDeserializer_DeserializedAst) Equals(arg_1 any) bool {
	o := arg_1
	_ = o
	if j2gort.RefEquals(o, this_) {
		return true
	}
	if j2gort.InstanceOf[TypedAstDeserializer_DeserializedAst](o) {
		var that TypedAstDeserializer_DeserializedAst = j2gort.AsIface[TypedAstDeserializer_DeserializedAst](o)
		_ = that
		return (((this_.filesystem.Equals(that.GetFilesystem()) && this_.colorRegistry.Equals(that.GetColorRegistry())) && (func() bool {
			if this_.externProperties == nil {
				return (that.GetExternProperties() == nil)
			}
			return this_.externProperties.Equals(that.GetExternProperties())
		}())) && this_.runtimeLibraries.Equals(that.GetRuntimeLibraries()))
	}
	return false
}

// Java method: AutoValue_TypedAstDeserializer_DeserializedAst.hashCode()
func (this_ *AutoValue_TypedAstDeserializer_DeserializedAst) HashCode() int32 {
	var h_d int32 = int32(1)
	_ = h_d
	h_d = (h_d * int32(1000003))
	h_d = (h_d ^ this_.filesystem.HashCode())
	h_d = (h_d * int32(1000003))
	h_d = (h_d ^ this_.colorRegistry.HashCode())
	h_d = (h_d * int32(1000003))
	h_d = (h_d ^ func() int32 {
		if this_.externProperties == nil {
			return int32(0)
		}
		return this_.externProperties.HashCode()
	}())
	h_d = (h_d * int32(1000003))
	h_d = (h_d ^ this_.runtimeLibraries.HashCode())
	return h_d
}
