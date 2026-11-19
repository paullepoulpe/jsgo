// Port of com/google/javascript/jscomp/AutoBuilder_AstAnalyzer_Options_Builder.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.AutoBuilder_AstAnalyzer_Options_Builder
type AutoBuilder_AstAnalyzer_Options_Builder struct {
	idHash_                      atomic.Int32
	useTypesForLocalOptimization bool `gson:"useTypesForLocalOptimization"`
	assumeGettersArePure         bool `gson:"assumeGettersArePure"`
	hasRegexpGlobalReferences    bool `gson:"hasRegexpGlobalReferences"`
	assumeKnownBuiltinsArePure   bool `gson:"assumeKnownBuiltinsArePure"`
	set_d0                       int8 `gson:"set$0"`
}

func (this_ *AutoBuilder_AstAnalyzer_Options_Builder) BindSelf_(s any) {
	_ = s
}
func (this_ *AutoBuilder_AstAnalyzer_Options_Builder) IsJavaNil_() bool { return this_ == nil }
func (*AutoBuilder_AstAnalyzer_Options_Builder) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.AutoBuilder_AstAnalyzer_Options_Builder", "AutoBuilder_AstAnalyzer_Options_Builder"
}
func (this_ *AutoBuilder_AstAnalyzer_Options_Builder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*AutoBuilder_AstAnalyzer_Options_Builder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.AutoBuilder_AstAnalyzer_Options_Builder", "com.google.javascript.jscomp.AutoBuilder_AstAnalyzer_Options_Builder", func(v_ any) bool {
		return j2gort.InstanceOf[*AutoBuilder_AstAnalyzer_Options_Builder](v_)
	})
}
func (*AutoBuilder_AstAnalyzer_Options_Builder) Implements_jscomp_AstAnalyzer_Options_Builder() {}

func (this_ *AutoBuilder_AstAnalyzer_Options_Builder) initFields_AutoBuilder_AstAnalyzer_Options_Builder() {
	_ = this_
}

// Java constructor: AutoBuilder_AstAnalyzer_Options_Builder()
func NewAutoBuilder_AstAnalyzer_Options_Builder() *AutoBuilder_AstAnalyzer_Options_Builder {
	this_ := &AutoBuilder_AstAnalyzer_Options_Builder{}
	InitAutoBuilder_AstAnalyzer_Options_Builder(this_)
	return this_
}

func InitAutoBuilder_AstAnalyzer_Options_Builder(this_ *AutoBuilder_AstAnalyzer_Options_Builder) {
	this_.initFields_AutoBuilder_AstAnalyzer_Options_Builder()
}

// Java constructor: AutoBuilder_AstAnalyzer_Options_Builder(com.google.javascript.jscomp.AstAnalyzer.Options source)
func NewAutoBuilder_AstAnalyzer_Options_Builder_1(source *AstAnalyzer_Options) *AutoBuilder_AstAnalyzer_Options_Builder {
	this_ := &AutoBuilder_AstAnalyzer_Options_Builder{}
	InitAutoBuilder_AstAnalyzer_Options_Builder_1(this_, source)
	return this_
}

func InitAutoBuilder_AstAnalyzer_Options_Builder_1(this_ *AutoBuilder_AstAnalyzer_Options_Builder, source *AstAnalyzer_Options) {
	_ = source
	this_.initFields_AutoBuilder_AstAnalyzer_Options_Builder()
	this_.useTypesForLocalOptimization = source.UseTypesForLocalOptimization()
	this_.assumeGettersArePure = source.AssumeGettersArePure()
	this_.hasRegexpGlobalReferences = source.HasRegexpGlobalReferences()
	this_.assumeKnownBuiltinsArePure = source.AssumeKnownBuiltinsArePure()
	this_.set_d0 = int8(15)
}

// Java method: AutoBuilder_AstAnalyzer_Options_Builder.setUseTypesForLocalOptimization(boolean useTypesForLocalOptimization)
func (this_ *AutoBuilder_AstAnalyzer_Options_Builder) SetUseTypesForLocalOptimization(useTypesForLocalOptimization bool) AstAnalyzer_Options_Builder {
	_ = useTypesForLocalOptimization
	this_.useTypesForLocalOptimization = useTypesForLocalOptimization
	this_.set_d0 = int8((int32(this_.set_d0) | int32(int8(1))))
	return this_
}

// Java method: AutoBuilder_AstAnalyzer_Options_Builder.setAssumeGettersArePure(boolean assumeGettersArePure)
func (this_ *AutoBuilder_AstAnalyzer_Options_Builder) SetAssumeGettersArePure(assumeGettersArePure bool) AstAnalyzer_Options_Builder {
	_ = assumeGettersArePure
	this_.assumeGettersArePure = assumeGettersArePure
	this_.set_d0 = int8((int32(this_.set_d0) | int32(int8(2))))
	return this_
}

// Java method: AutoBuilder_AstAnalyzer_Options_Builder.setHasRegexpGlobalReferences(boolean hasRegexpGlobalReferences)
func (this_ *AutoBuilder_AstAnalyzer_Options_Builder) SetHasRegexpGlobalReferences(hasRegexpGlobalReferences bool) AstAnalyzer_Options_Builder {
	_ = hasRegexpGlobalReferences
	this_.hasRegexpGlobalReferences = hasRegexpGlobalReferences
	this_.set_d0 = int8((int32(this_.set_d0) | int32(int8(4))))
	return this_
}

// Java method: AutoBuilder_AstAnalyzer_Options_Builder.setAssumeKnownBuiltinsArePure(boolean assumeKnownBuiltinsArePure)
func (this_ *AutoBuilder_AstAnalyzer_Options_Builder) SetAssumeKnownBuiltinsArePure(assumeKnownBuiltinsArePure bool) AstAnalyzer_Options_Builder {
	_ = assumeKnownBuiltinsArePure
	this_.assumeKnownBuiltinsArePure = assumeKnownBuiltinsArePure
	this_.set_d0 = int8((int32(this_.set_d0) | int32(int8(8))))
	return this_
}

// Java method: AutoBuilder_AstAnalyzer_Options_Builder.build()
func (this_ *AutoBuilder_AstAnalyzer_Options_Builder) Build() *AstAnalyzer_Options {
	if int32(this_.set_d0) != int32(15) {
		var missing *j2gort.Ext_java_lang_StringBuilder = j2gort.NewExt_java_lang_StringBuilder()
		_ = missing
		if (int32(this_.set_d0) & int32(1)) == int32(0) {
			missing.Append_1_string(" useTypesForLocalOptimization")
		}
		if (int32(this_.set_d0) & int32(2)) == int32(0) {
			missing.Append_1_string(" assumeGettersArePure")
		}
		if (int32(this_.set_d0) & int32(4)) == int32(0) {
			missing.Append_1_string(" hasRegexpGlobalReferences")
		}
		if (int32(this_.set_d0) & int32(8)) == int32(0) {
			missing.Append_1_string(" assumeKnownBuiltinsArePure")
		}
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Missing required properties:"+j2gort.StrVal(missing))))
	}
	return NewAstAnalyzer_Options(this_.useTypesForLocalOptimization, this_.assumeGettersArePure, this_.hasRegexpGlobalReferences, this_.assumeKnownBuiltinsArePure)
}
