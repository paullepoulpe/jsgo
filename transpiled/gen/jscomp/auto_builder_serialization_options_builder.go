// Port of com/google/javascript/jscomp/serialization/AutoBuilder_SerializationOptions_Builder.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.serialization.AutoBuilder_SerializationOptions_Builder
type AutoBuilder_SerializationOptions_Builder struct {
	SerializationOptions_BuilderBase
	runtimeLibraries j2gort.Ext_com_google_common_collect_ImmutableList `gson:"runtimeLibraries"`
	includeDebugInfo bool                                               `gson:"includeDebugInfo"`
	runValidation    bool                                               `gson:"runValidation"`
	set_d0           int8                                               `gson:"set$0"`
}

func (this_ *AutoBuilder_SerializationOptions_Builder) BindSelf_(s any) {
	_ = s
}
func (this_ *AutoBuilder_SerializationOptions_Builder) IsJavaNil_() bool { return this_ == nil }
func (*AutoBuilder_SerializationOptions_Builder) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.serialization.AutoBuilder_SerializationOptions_Builder", "AutoBuilder_SerializationOptions_Builder"
}
func (this_ *AutoBuilder_SerializationOptions_Builder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.SerializationOptions_BuilderBase.IdentityHash_()
}
func (*AutoBuilder_SerializationOptions_Builder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.serialization.AutoBuilder_SerializationOptions_Builder", "com.google.javascript.jscomp.serialization.AutoBuilder_SerializationOptions_Builder", func(v_ any) bool {
		return j2gort.InstanceOf[*AutoBuilder_SerializationOptions_Builder](v_)
	})
}
func (*AutoBuilder_SerializationOptions_Builder) Implements_jscomp_SerializationOptions_Builder() {}

func (this_ *AutoBuilder_SerializationOptions_Builder) initFields_AutoBuilder_SerializationOptions_Builder() {
	_ = this_
}

// Java constructor: AutoBuilder_SerializationOptions_Builder()
func NewAutoBuilder_SerializationOptions_Builder() *AutoBuilder_SerializationOptions_Builder {
	this_ := &AutoBuilder_SerializationOptions_Builder{}
	InitAutoBuilder_SerializationOptions_Builder(this_)
	return this_
}

func InitAutoBuilder_SerializationOptions_Builder(this_ *AutoBuilder_SerializationOptions_Builder) {
	InitSerializationOptions_Builder(&this_.SerializationOptions_BuilderBase)
	this_.initFields_AutoBuilder_SerializationOptions_Builder()
}

// Java constructor: AutoBuilder_SerializationOptions_Builder(com.google.javascript.jscomp.serialization.SerializationOptions source)
func NewAutoBuilder_SerializationOptions_Builder_1(source *SerializationOptions) *AutoBuilder_SerializationOptions_Builder {
	this_ := &AutoBuilder_SerializationOptions_Builder{}
	InitAutoBuilder_SerializationOptions_Builder_1(this_, source)
	return this_
}

func InitAutoBuilder_SerializationOptions_Builder_1(this_ *AutoBuilder_SerializationOptions_Builder, source *SerializationOptions) {
	_ = source
	InitSerializationOptions_Builder(&this_.SerializationOptions_BuilderBase)
	this_.initFields_AutoBuilder_SerializationOptions_Builder()
	this_.includeDebugInfo = source.IncludeDebugInfo()
	this_.runValidation = source.RunValidation()
	this_.runtimeLibraries = source.RuntimeLibraries()
	this_.set_d0 = int8(3)
}

// Java method: AutoBuilder_SerializationOptions_Builder.setIncludeDebugInfo(boolean includeDebugInfo)
func (this_ *AutoBuilder_SerializationOptions_Builder) SetIncludeDebugInfo(includeDebugInfo bool) SerializationOptions_Builder {
	_ = includeDebugInfo
	this_.includeDebugInfo = includeDebugInfo
	this_.set_d0 = int8((int32(this_.set_d0) | int32(int8(1))))
	return this_
}

// Java method: AutoBuilder_SerializationOptions_Builder.setRunValidation(boolean runValidation)
func (this_ *AutoBuilder_SerializationOptions_Builder) SetRunValidation(runValidation bool) SerializationOptions_Builder {
	_ = runValidation
	this_.runValidation = runValidation
	this_.set_d0 = int8((int32(this_.set_d0) | int32(int8(2))))
	return this_
}

// Java method: AutoBuilder_SerializationOptions_Builder.setRuntimeLibraries(com.google.common.collect.ImmutableList<java.lang.String> runtimeLibraries)
func (this_ *AutoBuilder_SerializationOptions_Builder) SetRuntimeLibraries(runtimeLibraries j2gort.Ext_com_google_common_collect_ImmutableList) SerializationOptions_Builder {
	_ = runtimeLibraries
	if runtimeLibraries == nil {
		panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", "Null runtimeLibraries")))
	}
	this_.runtimeLibraries = runtimeLibraries
	return this_
}

// Java method: AutoBuilder_SerializationOptions_Builder.build()
func (this_ *AutoBuilder_SerializationOptions_Builder) Build() *SerializationOptions {
	if (int32(this_.set_d0) != int32(3)) || (this_.runtimeLibraries == nil) {
		var missing *j2gort.Ext_java_lang_StringBuilder = j2gort.NewExt_java_lang_StringBuilder()
		_ = missing
		if (int32(this_.set_d0) & int32(1)) == int32(0) {
			missing.Append_1_string(" includeDebugInfo")
		}
		if (int32(this_.set_d0) & int32(2)) == int32(0) {
			missing.Append_1_string(" runValidation")
		}
		if this_.runtimeLibraries == nil {
			missing.Append_1_string(" runtimeLibraries")
		}
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Missing required properties:"+j2gort.StrVal(missing))))
	}
	return NewSerializationOptions(this_.includeDebugInfo, this_.runValidation, this_.runtimeLibraries)
}
