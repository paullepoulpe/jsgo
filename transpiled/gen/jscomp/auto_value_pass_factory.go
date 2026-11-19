// Port of com/google/javascript/jscomp/AutoValue_PassFactory.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.AutoValue_PassFactory
type AutoValue_PassFactory struct {
	PassFactoryBase
	name                string                                 `gson:"name"`
	condition           j2gort.Ext_java_util_function_Function `gson:"condition"`
	preconditionCheck   j2gort.Ext_java_util_function_Function `gson:"preconditionCheck"`
	internalFactory     j2gort.Ext_java_util_function_Function `gson:"internalFactory"`
	runInFixedPointLoop bool                                   `gson:"runInFixedPointLoop"`
}

func (this_ *AutoValue_PassFactory) BindSelf_(s any) {
	_ = s
	this_.PassFactoryBase.BindSelf_(s)
}
func (this_ *AutoValue_PassFactory) IsJavaNil_() bool { return this_ == nil }
func (*AutoValue_PassFactory) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.AutoValue_PassFactory", "AutoValue_PassFactory"
}
func (this_ *AutoValue_PassFactory) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.PassFactoryBase.IdentityHash_()
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.AutoValue_PassFactory", "com.google.javascript.jscomp.AutoValue_PassFactory", func(v_ any) bool {
		return j2gort.InstanceOf[*AutoValue_PassFactory](v_)
	})
}
func (*AutoValue_PassFactory) Implements_jscomp_PassFactory() {}

func (this_ *AutoValue_PassFactory) initFields_AutoValue_PassFactory() {
	_ = this_
	this_.name = j2gort.NullString
}

// Java constructor: AutoValue_PassFactory(java.lang.String name, java.util.function.Function<com.google.javascript.jscomp.CompilerOptions,java.lang.Boolean> condition, java.util.function.Function<com.google.javascript.jscomp.CompilerOptions,com.google.javascript.jscomp.PassFactory.PreconditionResult> preconditionCheck, boolean runInFixedPointLoop, java.util.function.Function<com.google.javascript.jscomp.AbstractCompiler,? extends com.google.javascript.jscomp.CompilerPass> internalFactory)
func NewAutoValue_PassFactory(name string, condition j2gort.Ext_java_util_function_Function, preconditionCheck j2gort.Ext_java_util_function_Function, runInFixedPointLoop bool, internalFactory j2gort.Ext_java_util_function_Function) *AutoValue_PassFactory {
	this_ := &AutoValue_PassFactory{}
	this_.BindSelf_(this_)
	InitAutoValue_PassFactory(this_, name, condition, preconditionCheck, runInFixedPointLoop, internalFactory)
	return this_
}

func InitAutoValue_PassFactory(this_ *AutoValue_PassFactory, name string, condition j2gort.Ext_java_util_function_Function, preconditionCheck j2gort.Ext_java_util_function_Function, runInFixedPointLoop bool, internalFactory j2gort.Ext_java_util_function_Function) {
	_ = name
	_ = condition
	_ = preconditionCheck
	_ = runInFixedPointLoop
	_ = internalFactory
	InitPassFactory(&this_.PassFactoryBase)
	this_.initFields_AutoValue_PassFactory()
	this_.name = name
	this_.condition = condition
	this_.preconditionCheck = preconditionCheck
	this_.runInFixedPointLoop = runInFixedPointLoop
	this_.internalFactory = internalFactory
}

// Java method: AutoValue_PassFactory.getName()
func (this_ *AutoValue_PassFactory) GetName() string {
	return this_.name
}

// Java method: AutoValue_PassFactory.getCondition()
func (this_ *AutoValue_PassFactory) GetCondition() j2gort.Ext_java_util_function_Function {
	return this_.condition
}

// Java method: AutoValue_PassFactory.getPreconditionCheck()
func (this_ *AutoValue_PassFactory) GetPreconditionCheck() j2gort.Ext_java_util_function_Function {
	return this_.preconditionCheck
}

// Java method: AutoValue_PassFactory.isRunInFixedPointLoop()
func (this_ *AutoValue_PassFactory) IsRunInFixedPointLoop() bool {
	return this_.runInFixedPointLoop
}

// Java method: AutoValue_PassFactory.getInternalFactory()
func (this_ *AutoValue_PassFactory) GetInternalFactory() j2gort.Ext_java_util_function_Function {
	return this_.internalFactory
}

// Java method: AutoValue_PassFactory.toString()
func (this_ *AutoValue_PassFactory) String() string {
	return ("PassFactory{name=" + j2gort.StrValString(this_.name) + ", condition=" + j2gort.StrVal(this_.condition) + ", preconditionCheck=" + j2gort.StrVal(this_.preconditionCheck) + ", runInFixedPointLoop=" + j2gort.StrValBool(this_.runInFixedPointLoop) + ", internalFactory=" + j2gort.StrVal(this_.internalFactory) + "}")
}

// Java method: AutoValue_PassFactory.equals(java.lang.Object o)
func (this_ *AutoValue_PassFactory) Equals(arg_1 any) bool {
	o := arg_1
	_ = o
	if j2gort.RefEquals(o, this_) {
		return true
	}
	if j2gort.InstanceOf[PassFactory](o) {
		var that PassFactory = j2gort.AsIface[PassFactory](o)
		_ = that
		return ((((func() bool {
			sr_2 := this_.name
			se_3 := that.GetName()
			return j2gort.CheckNotNullString(sr_2) == se_3
		}() && javalang.ObjectsEquals(j2gort.CheckNotNull(this_.condition), that.GetCondition())) && javalang.ObjectsEquals(j2gort.CheckNotNull(this_.preconditionCheck), that.GetPreconditionCheck())) && (this_.runInFixedPointLoop == that.IsRunInFixedPointLoop())) && javalang.ObjectsEquals(j2gort.CheckNotNull(this_.internalFactory), that.GetInternalFactory()))
	}
	return false
}

// Java method: AutoValue_PassFactory.hashCode()
func (this_ *AutoValue_PassFactory) HashCode() int32 {
	var h_d int32 = int32(1)
	_ = h_d
	h_d = (h_d * int32(1000003))
	h_d = (h_d ^ javalang.StringHashCode(j2gort.CheckNotNullString(this_.name)))
	h_d = (h_d * int32(1000003))
	h_d = (h_d ^ javalang.ObjectHashCode(j2gort.CheckNotNull(this_.condition)))
	h_d = (h_d * int32(1000003))
	h_d = (h_d ^ javalang.ObjectHashCode(j2gort.CheckNotNull(this_.preconditionCheck)))
	h_d = (h_d * int32(1000003))
	h_d = (h_d ^ func() int32 {
		if this_.runInFixedPointLoop {
			return int32(1231)
		}
		return int32(1237)
	}())
	h_d = (h_d * int32(1000003))
	h_d = (h_d ^ javalang.ObjectHashCode(j2gort.CheckNotNull(this_.internalFactory)))
	return h_d
}

// Java method: AutoValue_PassFactory.toBuilder()
func (this_ *AutoValue_PassFactory) ToBuilder() PassFactory_Builder {
	return NewAutoValue_PassFactory_Builder_1(this_)
}

// Java class: com.google.javascript.jscomp.AutoValue_PassFactory.Builder
type AutoValue_PassFactory_Builder struct {
	PassFactory_BuilderBase
	name                string                                 `gson:"name"`
	condition           j2gort.Ext_java_util_function_Function `gson:"condition"`
	preconditionCheck   j2gort.Ext_java_util_function_Function `gson:"preconditionCheck"`
	internalFactory     j2gort.Ext_java_util_function_Function `gson:"internalFactory"`
	runInFixedPointLoop bool                                   `gson:"runInFixedPointLoop"`
	set_d0              int8                                   `gson:"set$0"`
}

func (this_ *AutoValue_PassFactory_Builder) BindSelf_(s any) {
	_ = s
	this_.PassFactory_BuilderBase.BindSelf_(s)
}
func (this_ *AutoValue_PassFactory_Builder) IsJavaNil_() bool { return this_ == nil }
func (*AutoValue_PassFactory_Builder) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.AutoValue_PassFactory$Builder", "Builder"
}
func (this_ *AutoValue_PassFactory_Builder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.PassFactory_BuilderBase.IdentityHash_()
}
func (*AutoValue_PassFactory_Builder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.AutoValue_PassFactory.Builder", "com.google.javascript.jscomp.AutoValue_PassFactory$Builder", func(v_ any) bool {
		return j2gort.InstanceOf[*AutoValue_PassFactory_Builder](v_)
	})
}
func (*AutoValue_PassFactory_Builder) Implements_jscomp_PassFactory_Builder() {}

func (this_ *AutoValue_PassFactory_Builder) initFields_AutoValue_PassFactory_Builder() {
	_ = this_
	this_.name = j2gort.NullString
}

// Java constructor: Builder()
func NewAutoValue_PassFactory_Builder() *AutoValue_PassFactory_Builder {
	this_ := &AutoValue_PassFactory_Builder{}
	this_.BindSelf_(this_)
	InitAutoValue_PassFactory_Builder(this_)
	return this_
}

func InitAutoValue_PassFactory_Builder(this_ *AutoValue_PassFactory_Builder) {
	InitPassFactory_Builder(&this_.PassFactory_BuilderBase)
	this_.initFields_AutoValue_PassFactory_Builder()
}

// Java constructor: Builder(com.google.javascript.jscomp.PassFactory source)
func NewAutoValue_PassFactory_Builder_1(source PassFactory) *AutoValue_PassFactory_Builder {
	this_ := &AutoValue_PassFactory_Builder{}
	this_.BindSelf_(this_)
	InitAutoValue_PassFactory_Builder_1(this_, source)
	return this_
}

func InitAutoValue_PassFactory_Builder_1(this_ *AutoValue_PassFactory_Builder, source PassFactory) {
	_ = source
	InitPassFactory_Builder(&this_.PassFactory_BuilderBase)
	this_.initFields_AutoValue_PassFactory_Builder()
	this_.name = source.GetName()
	this_.condition = source.GetCondition()
	this_.preconditionCheck = source.GetPreconditionCheck()
	this_.runInFixedPointLoop = source.IsRunInFixedPointLoop()
	this_.internalFactory = source.GetInternalFactory()
	this_.set_d0 = int8(1)
}

// Java method: Builder.setName(java.lang.String name)
func (this_ *AutoValue_PassFactory_Builder) SetName(name string) PassFactory_Builder {
	_ = name
	if name == j2gort.NullString {
		panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", "Null name")))
	}
	this_.name = name
	return this_
}

// Java method: Builder.setCondition(java.util.function.Function<com.google.javascript.jscomp.CompilerOptions,java.lang.Boolean> condition)
func (this_ *AutoValue_PassFactory_Builder) SetCondition(condition j2gort.Ext_java_util_function_Function) PassFactory_Builder {
	_ = condition
	if condition == nil {
		panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", "Null condition")))
	}
	this_.condition = condition
	return this_
}

// Java method: Builder.setPreconditionCheck(java.util.function.Function<com.google.javascript.jscomp.CompilerOptions,com.google.javascript.jscomp.PassFactory.PreconditionResult> preconditionCheck)
func (this_ *AutoValue_PassFactory_Builder) SetPreconditionCheck(preconditionCheck j2gort.Ext_java_util_function_Function) PassFactory_Builder {
	_ = preconditionCheck
	if preconditionCheck == nil {
		panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", "Null preconditionCheck")))
	}
	this_.preconditionCheck = preconditionCheck
	return this_
}

// Java method: Builder.setRunInFixedPointLoop(boolean runInFixedPointLoop)
func (this_ *AutoValue_PassFactory_Builder) SetRunInFixedPointLoop(runInFixedPointLoop bool) PassFactory_Builder {
	_ = runInFixedPointLoop
	this_.runInFixedPointLoop = runInFixedPointLoop
	this_.set_d0 = int8((int32(this_.set_d0) | int32(int8(1))))
	return this_
}

// Java method: Builder.setInternalFactory(java.util.function.Function<com.google.javascript.jscomp.AbstractCompiler,? extends com.google.javascript.jscomp.CompilerPass> internalFactory)
func (this_ *AutoValue_PassFactory_Builder) SetInternalFactory(internalFactory j2gort.Ext_java_util_function_Function) PassFactory_Builder {
	_ = internalFactory
	if internalFactory == nil {
		panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", "Null internalFactory")))
	}
	this_.internalFactory = internalFactory
	return this_
}

// Java method: Builder.autoBuild()
func (this_ *AutoValue_PassFactory_Builder) AutoBuild() PassFactory {
	if ((((int32(this_.set_d0) != int32(1)) || (this_.name == j2gort.NullString)) || (this_.condition == nil)) || (this_.preconditionCheck == nil)) || (this_.internalFactory == nil) {
		var missing *j2gort.Ext_java_lang_StringBuilder = j2gort.NewExt_java_lang_StringBuilder()
		_ = missing
		if this_.name == j2gort.NullString {
			missing.Append_1_string(" name")
		}
		if this_.condition == nil {
			missing.Append_1_string(" condition")
		}
		if this_.preconditionCheck == nil {
			missing.Append_1_string(" preconditionCheck")
		}
		if (int32(this_.set_d0) & int32(1)) == int32(0) {
			missing.Append_1_string(" runInFixedPointLoop")
		}
		if this_.internalFactory == nil {
			missing.Append_1_string(" internalFactory")
		}
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Missing required properties:"+j2gort.StrVal(missing))))
	}
	return NewAutoValue_PassFactory(this_.name, this_.condition, this_.preconditionCheck, this_.runInFixedPointLoop, this_.internalFactory)
}
