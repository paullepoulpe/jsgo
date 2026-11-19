// Port of com/google/javascript/jscomp/parsing/AutoValue_Config.java (generated mechanically by j2go).

package parsing

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.parsing.AutoValue_Config
type AutoValue_Config struct {
	ConfigBase
	languageMode          *Config_LanguageMode                              `gson:"languageMode"`
	annotations           j2gort.Ext_com_google_common_collect_ImmutableMap `gson:"annotations"`
	suppressionNames      j2gort.Ext_com_google_common_collect_ImmutableSet `gson:"suppressionNames"`
	closurePrimitiveNames j2gort.Ext_com_google_common_collect_ImmutableSet `gson:"closurePrimitiveNames"`
	strictMode            Config_StrictMode                                 `gson:"strictMode"`
	jsDocParsingMode      Config_JsDocParsing                               `gson:"jsDocParsingMode"`
	runMode               Config_RunMode                                    `gson:"runMode"`
	parseInlineSourceMaps bool                                              `gson:"parseInlineSourceMaps"`
}

func (this_ *AutoValue_Config) BindSelf_(s any) {
	_ = s
	this_.ConfigBase.BindSelf_(s)
}
func (this_ *AutoValue_Config) IsJavaNil_() bool { return this_ == nil }
func (*AutoValue_Config) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.AutoValue_Config", "AutoValue_Config"
}
func (this_ *AutoValue_Config) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.ConfigBase.IdentityHash_()
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.AutoValue_Config", "com.google.javascript.jscomp.parsing.AutoValue_Config", func(v_ any) bool {
		return j2gort.InstanceOf[*AutoValue_Config](v_)
	})
}
func (*AutoValue_Config) Implements_parsing_Config() {}

func (this_ *AutoValue_Config) initFields_AutoValue_Config() {
	_ = this_
}

// Java constructor: AutoValue_Config(com.google.javascript.jscomp.parsing.Config.LanguageMode languageMode, com.google.javascript.jscomp.parsing.Config.StrictMode strictMode, com.google.javascript.jscomp.parsing.Config.JsDocParsing jsDocParsingMode, com.google.javascript.jscomp.parsing.Config.RunMode runMode, com.google.common.collect.ImmutableMap<java.lang.String,com.google.javascript.jscomp.parsing.Annotation> annotations, com.google.common.collect.ImmutableSet<java.lang.String> suppressionNames, com.google.common.collect.ImmutableSet<java.lang.String> closurePrimitiveNames, boolean parseInlineSourceMaps)
func NewAutoValue_Config(languageMode *Config_LanguageMode, strictMode Config_StrictMode, jsDocParsingMode Config_JsDocParsing, runMode Config_RunMode, annotations j2gort.Ext_com_google_common_collect_ImmutableMap, suppressionNames j2gort.Ext_com_google_common_collect_ImmutableSet, closurePrimitiveNames j2gort.Ext_com_google_common_collect_ImmutableSet, parseInlineSourceMaps bool) *AutoValue_Config {
	this_ := &AutoValue_Config{}
	this_.BindSelf_(this_)
	InitAutoValue_Config(this_, languageMode, strictMode, jsDocParsingMode, runMode, annotations, suppressionNames, closurePrimitiveNames, parseInlineSourceMaps)
	return this_
}

func InitAutoValue_Config(this_ *AutoValue_Config, languageMode *Config_LanguageMode, strictMode Config_StrictMode, jsDocParsingMode Config_JsDocParsing, runMode Config_RunMode, annotations j2gort.Ext_com_google_common_collect_ImmutableMap, suppressionNames j2gort.Ext_com_google_common_collect_ImmutableSet, closurePrimitiveNames j2gort.Ext_com_google_common_collect_ImmutableSet, parseInlineSourceMaps bool) {
	_ = languageMode
	_ = strictMode
	_ = jsDocParsingMode
	_ = runMode
	_ = annotations
	_ = suppressionNames
	_ = closurePrimitiveNames
	_ = parseInlineSourceMaps
	InitConfig(&this_.ConfigBase)
	this_.initFields_AutoValue_Config()
	this_.languageMode = languageMode
	this_.strictMode = strictMode
	this_.jsDocParsingMode = jsDocParsingMode
	this_.runMode = runMode
	this_.annotations = annotations
	this_.suppressionNames = suppressionNames
	this_.closurePrimitiveNames = closurePrimitiveNames
	this_.parseInlineSourceMaps = parseInlineSourceMaps
}

// Java method: AutoValue_Config.languageMode()
func (this_ *AutoValue_Config) LanguageMode() *Config_LanguageMode {
	return this_.languageMode
}

// Java method: AutoValue_Config.strictMode()
func (this_ *AutoValue_Config) StrictMode() Config_StrictMode {
	return this_.strictMode
}

// Java method: AutoValue_Config.jsDocParsingMode()
func (this_ *AutoValue_Config) JsDocParsingMode() Config_JsDocParsing {
	return this_.jsDocParsingMode
}

// Java method: AutoValue_Config.runMode()
func (this_ *AutoValue_Config) RunMode() Config_RunMode {
	return this_.runMode
}

// Java method: AutoValue_Config.annotations()
func (this_ *AutoValue_Config) Annotations() j2gort.Ext_com_google_common_collect_ImmutableMap {
	return this_.annotations
}

// Java method: AutoValue_Config.suppressionNames()
func (this_ *AutoValue_Config) SuppressionNames() j2gort.Ext_com_google_common_collect_ImmutableSet {
	return this_.suppressionNames
}

// Java method: AutoValue_Config.closurePrimitiveNames()
func (this_ *AutoValue_Config) ClosurePrimitiveNames() j2gort.Ext_com_google_common_collect_ImmutableSet {
	return this_.closurePrimitiveNames
}

// Java method: AutoValue_Config.parseInlineSourceMaps()
func (this_ *AutoValue_Config) ParseInlineSourceMaps() bool {
	return this_.parseInlineSourceMaps
}

// Java method: AutoValue_Config.toString()
func (this_ *AutoValue_Config) String() string {
	return ("Config{languageMode=" + j2gort.StrVal(this_.languageMode) + ", strictMode=" + j2gort.StrVal(this_.strictMode) + ", jsDocParsingMode=" + j2gort.StrVal(this_.jsDocParsingMode) + ", runMode=" + j2gort.StrVal(this_.runMode) + ", annotations=" + j2gort.StrVal(this_.annotations) + ", suppressionNames=" + j2gort.StrVal(this_.suppressionNames) + ", closurePrimitiveNames=" + j2gort.StrVal(this_.closurePrimitiveNames) + ", parseInlineSourceMaps=" + j2gort.StrValBool(this_.parseInlineSourceMaps) + "}")
}

// Java method: AutoValue_Config.equals(java.lang.Object o)
func (this_ *AutoValue_Config) Equals(arg_1 any) bool {
	o := arg_1
	_ = o
	if j2gort.RefEquals(o, this_) {
		return true
	}
	if j2gort.InstanceOf[Config](o) {
		var that Config = j2gort.AsIface[Config](o)
		_ = that
		return (((((((j2gort.RefEquals(this_.languageMode, that.LanguageMode()) && func() bool {
			re_2 := this_.strictMode
			ae_3 := that.StrictMode()
			if re_2 == 0 {
				j2gort.ThrowNPE()
			}
			return re_2 == ae_3
		}()) && func() bool {
			re_4 := this_.jsDocParsingMode
			ae_5 := that.JsDocParsingMode()
			if re_4 == 0 {
				j2gort.ThrowNPE()
			}
			return re_4 == ae_5
		}()) && func() bool {
			re_6 := this_.runMode
			ae_7 := that.RunMode()
			if re_6 == 0 {
				j2gort.ThrowNPE()
			}
			return re_6 == ae_7
		}()) && this_.annotations.Equals(that.Annotations())) && this_.suppressionNames.Equals(that.SuppressionNames())) && this_.closurePrimitiveNames.Equals(that.ClosurePrimitiveNames())) && (this_.parseInlineSourceMaps == that.ParseInlineSourceMaps()))
	}
	return false
}

// Java method: AutoValue_Config.hashCode()
func (this_ *AutoValue_Config) HashCode() int32 {
	var h_d int32 = int32(1)
	_ = h_d
	h_d = (h_d * int32(1000003))
	h_d = (h_d ^ this_.languageMode.HashCode())
	h_d = (h_d * int32(1000003))
	h_d = (h_d ^ j2gort.CheckNotNullEnum(this_.strictMode).HashCode())
	h_d = (h_d * int32(1000003))
	h_d = (h_d ^ j2gort.CheckNotNullEnum(this_.jsDocParsingMode).HashCode())
	h_d = (h_d * int32(1000003))
	h_d = (h_d ^ j2gort.CheckNotNullEnum(this_.runMode).HashCode())
	h_d = (h_d * int32(1000003))
	h_d = (h_d ^ this_.annotations.HashCode())
	h_d = (h_d * int32(1000003))
	h_d = (h_d ^ this_.suppressionNames.HashCode())
	h_d = (h_d * int32(1000003))
	h_d = (h_d ^ this_.closurePrimitiveNames.HashCode())
	h_d = (h_d * int32(1000003))
	h_d = (h_d ^ func() int32 {
		if this_.parseInlineSourceMaps {
			return int32(1231)
		}
		return int32(1237)
	}())
	return h_d
}

// Java method: AutoValue_Config.toBuilder()
func (this_ *AutoValue_Config) ToBuilder() Config_Builder {
	return NewAutoValue_Config_Builder_1(this_)
}

// Java class: com.google.javascript.jscomp.parsing.AutoValue_Config.Builder
type AutoValue_Config_Builder struct {
	Config_BuilderBase
	languageMode          *Config_LanguageMode                              `gson:"languageMode"`
	annotations           j2gort.Ext_com_google_common_collect_ImmutableMap `gson:"annotations"`
	suppressionNames      j2gort.Ext_com_google_common_collect_ImmutableSet `gson:"suppressionNames"`
	closurePrimitiveNames j2gort.Ext_com_google_common_collect_ImmutableSet `gson:"closurePrimitiveNames"`
	strictMode            Config_StrictMode                                 `gson:"strictMode"`
	jsDocParsingMode      Config_JsDocParsing                               `gson:"jsDocParsingMode"`
	runMode               Config_RunMode                                    `gson:"runMode"`
	parseInlineSourceMaps bool                                              `gson:"parseInlineSourceMaps"`
	set_d0                int8                                              `gson:"set$0"`
}

func (this_ *AutoValue_Config_Builder) BindSelf_(s any) {
	_ = s
	this_.Config_BuilderBase.BindSelf_(s)
}
func (this_ *AutoValue_Config_Builder) IsJavaNil_() bool { return this_ == nil }
func (*AutoValue_Config_Builder) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.AutoValue_Config$Builder", "Builder"
}
func (this_ *AutoValue_Config_Builder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.Config_BuilderBase.IdentityHash_()
}
func (*AutoValue_Config_Builder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.AutoValue_Config.Builder", "com.google.javascript.jscomp.parsing.AutoValue_Config$Builder", func(v_ any) bool {
		return j2gort.InstanceOf[*AutoValue_Config_Builder](v_)
	})
}
func (*AutoValue_Config_Builder) Implements_parsing_Config_Builder() {}

func (this_ *AutoValue_Config_Builder) initFields_AutoValue_Config_Builder() {
	_ = this_
}

// Java constructor: Builder()
func NewAutoValue_Config_Builder() *AutoValue_Config_Builder {
	this_ := &AutoValue_Config_Builder{}
	this_.BindSelf_(this_)
	InitAutoValue_Config_Builder(this_)
	return this_
}

func InitAutoValue_Config_Builder(this_ *AutoValue_Config_Builder) {
	InitConfig_Builder(&this_.Config_BuilderBase)
	this_.initFields_AutoValue_Config_Builder()
}

// Java constructor: Builder(com.google.javascript.jscomp.parsing.Config source)
func NewAutoValue_Config_Builder_1(source Config) *AutoValue_Config_Builder {
	this_ := &AutoValue_Config_Builder{}
	this_.BindSelf_(this_)
	InitAutoValue_Config_Builder_1(this_, source)
	return this_
}

func InitAutoValue_Config_Builder_1(this_ *AutoValue_Config_Builder, source Config) {
	_ = source
	InitConfig_Builder(&this_.Config_BuilderBase)
	this_.initFields_AutoValue_Config_Builder()
	this_.languageMode = source.LanguageMode()
	this_.strictMode = source.StrictMode()
	this_.jsDocParsingMode = source.JsDocParsingMode()
	this_.runMode = source.RunMode()
	this_.annotations = source.Annotations()
	this_.suppressionNames = source.SuppressionNames()
	this_.closurePrimitiveNames = source.ClosurePrimitiveNames()
	this_.parseInlineSourceMaps = source.ParseInlineSourceMaps()
	this_.set_d0 = int8(1)
}

// Java method: Builder.setLanguageMode(com.google.javascript.jscomp.parsing.Config.LanguageMode languageMode)
func (this_ *AutoValue_Config_Builder) SetLanguageMode(languageMode *Config_LanguageMode) Config_Builder {
	_ = languageMode
	if languageMode == nil {
		panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", "Null languageMode")))
	}
	this_.languageMode = languageMode
	return this_
}

// Java method: Builder.setStrictMode(com.google.javascript.jscomp.parsing.Config.StrictMode strictMode)
func (this_ *AutoValue_Config_Builder) SetStrictMode(strictMode Config_StrictMode) Config_Builder {
	_ = strictMode
	if strictMode == 0 {
		panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", "Null strictMode")))
	}
	this_.strictMode = strictMode
	return this_
}

// Java method: Builder.setJsDocParsingMode(com.google.javascript.jscomp.parsing.Config.JsDocParsing jsDocParsingMode)
func (this_ *AutoValue_Config_Builder) SetJsDocParsingMode(jsDocParsingMode Config_JsDocParsing) Config_Builder {
	_ = jsDocParsingMode
	if jsDocParsingMode == 0 {
		panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", "Null jsDocParsingMode")))
	}
	this_.jsDocParsingMode = jsDocParsingMode
	return this_
}

// Java method: Builder.setRunMode(com.google.javascript.jscomp.parsing.Config.RunMode runMode)
func (this_ *AutoValue_Config_Builder) SetRunMode(runMode Config_RunMode) Config_Builder {
	_ = runMode
	if runMode == 0 {
		panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", "Null runMode")))
	}
	this_.runMode = runMode
	return this_
}

// Java method: Builder.setAnnotations(com.google.common.collect.ImmutableMap<java.lang.String,com.google.javascript.jscomp.parsing.Annotation> annotations)
func (this_ *AutoValue_Config_Builder) SetAnnotations(annotations j2gort.Ext_com_google_common_collect_ImmutableMap) Config_Builder {
	_ = annotations
	if annotations == nil {
		panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", "Null annotations")))
	}
	this_.annotations = annotations
	return this_
}

// Java method: Builder.setSuppressionNames(java.lang.Iterable<java.lang.String> suppressionNames)
func (this_ *AutoValue_Config_Builder) SetSuppressionNames(suppressionNames j2gort.Ext_java_lang_Iterable) Config_Builder {
	_ = suppressionNames
	this_.suppressionNames = j2gort.Ext_com_google_common_collect_ImmutableSet_CopyOf_1_Ext_java_lang_Iterable(suppressionNames)
	return this_
}

// Java method: Builder.setClosurePrimitiveNames(java.lang.Iterable<java.lang.String> closurePrimitiveNames)
func (this_ *AutoValue_Config_Builder) SetClosurePrimitiveNames(closurePrimitiveNames j2gort.Ext_java_lang_Iterable) Config_Builder {
	_ = closurePrimitiveNames
	this_.closurePrimitiveNames = j2gort.Ext_com_google_common_collect_ImmutableSet_CopyOf_1_Ext_java_lang_Iterable(closurePrimitiveNames)
	return this_
}

// Java method: Builder.setParseInlineSourceMaps(boolean parseInlineSourceMaps)
func (this_ *AutoValue_Config_Builder) SetParseInlineSourceMaps(parseInlineSourceMaps bool) Config_Builder {
	_ = parseInlineSourceMaps
	this_.parseInlineSourceMaps = parseInlineSourceMaps
	this_.set_d0 = int8((int32(this_.set_d0) | int32(int8(1))))
	return this_
}

// Java method: Builder.build()
func (this_ *AutoValue_Config_Builder) Build() Config {
	if (((((((int32(this_.set_d0) != int32(1)) || (this_.languageMode == nil)) || (this_.strictMode == 0)) || (this_.jsDocParsingMode == 0)) || (this_.runMode == 0)) || (this_.annotations == nil)) || (this_.suppressionNames == nil)) || (this_.closurePrimitiveNames == nil) {
		var missing *j2gort.Ext_java_lang_StringBuilder = j2gort.NewExt_java_lang_StringBuilder()
		_ = missing
		if this_.languageMode == nil {
			missing.Append_1_string(" languageMode")
		}
		if this_.strictMode == 0 {
			missing.Append_1_string(" strictMode")
		}
		if this_.jsDocParsingMode == 0 {
			missing.Append_1_string(" jsDocParsingMode")
		}
		if this_.runMode == 0 {
			missing.Append_1_string(" runMode")
		}
		if this_.annotations == nil {
			missing.Append_1_string(" annotations")
		}
		if this_.suppressionNames == nil {
			missing.Append_1_string(" suppressionNames")
		}
		if this_.closurePrimitiveNames == nil {
			missing.Append_1_string(" closurePrimitiveNames")
		}
		if (int32(this_.set_d0) & int32(1)) == int32(0) {
			missing.Append_1_string(" parseInlineSourceMaps")
		}
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Missing required properties:"+j2gort.StrVal(missing))))
	}
	return NewAutoValue_Config(this_.languageMode, this_.strictMode, this_.jsDocParsingMode, this_.runMode, this_.annotations, this_.suppressionNames, this_.closurePrimitiveNames, this_.parseInlineSourceMaps)
}
