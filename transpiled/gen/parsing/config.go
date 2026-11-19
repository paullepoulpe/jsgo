/*
 * Copyright 2009 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/parsing/Config.java (generated mechanically by j2go).

package parsing

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java polymorphic class interface: com.google.javascript.jscomp.parsing.Config
type Config interface {
	Implements_parsing_Config()
	IsJavaNil_() bool
	GetConfigBase() *ConfigBase
	LanguageMode() *Config_LanguageMode
	StrictMode() Config_StrictMode
	JsDocParsingMode() Config_JsDocParsing
	RunMode() Config_RunMode
	Annotations() j2gort.Ext_com_google_common_collect_ImmutableMap
	SuppressionNames() j2gort.Ext_com_google_common_collect_ImmutableSet
	ClosurePrimitiveNames() j2gort.Ext_com_google_common_collect_ImmutableSet
	ParseInlineSourceMaps() bool
	AnnotationNames() j2gort.Ext_com_google_common_collect_ImmutableSet
	ToBuilder() Config_Builder
}

func ToIface_Config[T interface {
	Config
	comparable
}](v T) Config {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_Config[T Config](v T) Config {
	var vi Config = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.jscomp.parsing.Config
type ConfigBase struct {
	self_   Config
	idHash_ atomic.Int32
}

func (this_ *ConfigBase) BindSelf_(s any) {
	_ = s
	this_.self_ = s.(Config)
}
func (this_ *ConfigBase) IsJavaNil_() bool { return this_ == nil }
func (*ConfigBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.Config", "Config"
}
func (this_ *ConfigBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.Config", "com.google.javascript.jscomp.parsing.Config", func(v_ any) bool {
		return j2gort.InstanceOf[Config](v_)
	})
}
func (b *ConfigBase) GetConfigBase() *ConfigBase { return b }
func (*ConfigBase) Implements_parsing_Config()   {}

func (this_ *ConfigBase) initFields_Config() {
	_ = this_
}

func InitConfig(this_ *ConfigBase) {
	this_.initFields_Config()
}

// Java method: Config.languageMode()
func (this_ *ConfigBase) LanguageMode() *Config_LanguageMode {
	panic("abstract or native method invoked")
}

// Java method: Config.strictMode()
func (this_ *ConfigBase) StrictMode() Config_StrictMode {
	panic("abstract or native method invoked")
}

// Java method: Config.jsDocParsingMode()
func (this_ *ConfigBase) JsDocParsingMode() Config_JsDocParsing {
	panic("abstract or native method invoked")
}

// Java method: Config.runMode()
func (this_ *ConfigBase) RunMode() Config_RunMode {
	panic("abstract or native method invoked")
}

// Java method: Config.annotations()
func (this_ *ConfigBase) Annotations() j2gort.Ext_com_google_common_collect_ImmutableMap {
	panic("abstract or native method invoked")
}

// Java method: Config.suppressionNames()
func (this_ *ConfigBase) SuppressionNames() j2gort.Ext_com_google_common_collect_ImmutableSet {
	panic("abstract or native method invoked")
}

// Java method: Config.closurePrimitiveNames()
func (this_ *ConfigBase) ClosurePrimitiveNames() j2gort.Ext_com_google_common_collect_ImmutableSet {
	panic("abstract or native method invoked")
}

// Java method: Config.parseInlineSourceMaps()
func (this_ *ConfigBase) ParseInlineSourceMaps() bool {
	panic("abstract or native method invoked")
}

// Java method: Config.annotationNames()
func (this_ *ConfigBase) AnnotationNames() j2gort.Ext_com_google_common_collect_ImmutableSet {
	return j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](this_.self_.Annotations().KeySet_0())
}

// Java method: Config.toBuilder()
func (this_ *ConfigBase) ToBuilder() Config_Builder {
	panic("abstract or native method invoked")
}

// Java static method: Config.builder()
func Config_S_Builder() Config_Builder {
	return NewAutoValue_Config_Builder().SetLanguageMode(StaticGet_Config_LanguageMode_UNSUPPORTED()).SetStrictMode(StaticGet_Config_StrictMode_STRICT()).SetJsDocParsingMode(StaticGet_Config_JsDocParsing_TYPES_ONLY()).SetRunMode(StaticGet_Config_RunMode_STOP_AFTER_ERROR()).GetConfig_BuilderBase().SetExtraAnnotationNames(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](j2gort.Ext_com_google_common_collect_ImmutableSet_Of_0())).SetSuppressionNames(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](j2gort.Ext_com_google_common_collect_ImmutableSet_Of_0())).SetClosurePrimitiveNames(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](j2gort.Ext_com_google_common_collect_ImmutableSet_Of_0())).SetParseInlineSourceMaps(false)
}

// Java static method: Config.buildAnnotations(java.lang.Iterable<java.lang.String> allowlist)
func Config_BuildAnnotations(allowlist j2gort.Ext_java_lang_Iterable) j2gort.Ext_com_google_common_collect_ImmutableMap {
	_ = allowlist
	var annotationsBuilder j2gort.Ext_com_google_common_collect_ImmutableMap_Builder = j2gort.Ext_com_google_common_collect_ImmutableMap_S_Builder()
	_ = annotationsBuilder
	annotationsBuilder.PutAll_1_Ext_java_util_Map(j2gort.AsIface[j2gort.Ext_java_util_Map](StaticGet_Annotation_RecognizedAnnotations()))
	{
		cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(allowlist))
		for cur_1.HasNext() {
			unrecognizedAnnotation := cur_1.NextStr()
			_ = unrecognizedAnnotation
			if !(j2gort.CheckNotNullString(unrecognizedAnnotation) == "") && !(j2gort.MapContainsKeyStr(StaticGet_Annotation_RecognizedAnnotations(), unrecognizedAnnotation)) {
				annotationsBuilder.Put_2(j2gort.BoxString(unrecognizedAnnotation), j2gort.BoxNullableVal(StaticGet_Annotation_NOT_IMPLEMENTED()))
			}
		}
	}
	return j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap](annotationsBuilder.BuildOrThrow())
}

// Java enum: com.google.javascript.jscomp.parsing.Config.StrictMode
type Config_StrictMode uint8

const (
	enum_Config_StrictMode_STRICT Config_StrictMode = 1
	enum_Config_StrictMode_SLOPPY Config_StrictMode = 2
)

var enumBaseTable_Config_StrictMode = [...]j2gort.EnumBase{
	{Name_: "STRICT", Ordinal_: 0},
	{Name_: "SLOPPY", Ordinal_: 1},
}
var enumValues_Config_StrictMode = []Config_StrictMode{1, 2}

func StaticGet_Config_StrictMode_STRICT() Config_StrictMode {
	return enum_Config_StrictMode_STRICT
}
func StaticGet_Config_StrictMode_SLOPPY() Config_StrictMode {
	return enum_Config_StrictMode_SLOPPY
}

func Config_StrictMode_Values() []Config_StrictMode {
	return j2gort.CloneSlice(enumValues_Config_StrictMode)
}

func Config_StrictMode_ValueOf(name string) Config_StrictMode {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_Config_StrictMode {
		if enumBaseTable_Config_StrictMode[i_].Name_ == name {
			return Config_StrictMode(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.parsing.Config.StrictMode." + name))
}

// Java method: StrictMode.isStrict()
func (this_ Config_StrictMode) IsStrict() bool {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return (this_ == StaticGet_Config_StrictMode_STRICT())
}

func (this_ Config_StrictMode) IsJavaNil_() bool                              { return this_ == 0 }
func (Config_StrictMode) HasIdentityEquals_()                                 {}
func (Config_StrictMode) Implements_j2gort_Ext_java_lang_Enum()               {}
func (Config_StrictMode) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ Config_StrictMode) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (Config_StrictMode) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.Config$StrictMode", "StrictMode"
}
func (this_ Config_StrictMode) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_Config_StrictMode[this_-1]
}
func (this_ Config_StrictMode) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ Config_StrictMode) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ Config_StrictMode) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_Config_StrictMode[this_-1].Name_
}
func (this_ Config_StrictMode) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_Config_StrictMode[this_-1].Name_
}
func (this_ Config_StrictMode) Equals(other_ any) bool {
	o_, ok_ := other_.(Config_StrictMode)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ Config_StrictMode) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_Config_StrictMode[this_-1].HashCode()
}
func (this_ Config_StrictMode) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_Config_StrictMode[this_-1].CompareTo(other_)
}
func (this_ Config_StrictMode) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_Config_StrictMode[this_-1].Name_)
}
func init() {
	vals_ := []any{Config_StrictMode(1), Config_StrictMode(2)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.parsing.Config.StrictMode", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.parsing.Config$StrictMode", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.Config.StrictMode", "com.google.javascript.jscomp.parsing.Config$StrictMode", func(v_ any) bool {
		return j2gort.InstanceOf[Config_StrictMode](v_)
	})
}

// Java enum: com.google.javascript.jscomp.parsing.Config.LanguageMode
type Config_LanguageMode struct {
	j2gort.EnumBase
	featureSet *pkg_parser.FeatureSet
}

func (this_ *Config_LanguageMode) GetEnumBase() *j2gort.EnumBase {
	if this_ == nil {
		return nil
	}
	return &this_.EnumBase
}
func (this_ *Config_LanguageMode) EnumOrdinal_() int32 {
	if this_ == nil {
		return -1
	}
	return this_.EnumBase.Ordinal_
}
func (*Config_LanguageMode) Implements_j2gort_Ext_java_lang_Enum()               {}
func (*Config_LanguageMode) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ *Config_LanguageMode) initFields_Config_LanguageMode() {
	_ = this_
}

func initEnumConfig_LanguageMode(this_ *Config_LanguageMode, enumName_ string, enumOrd_ int32, featureSet *pkg_parser.FeatureSet) *Config_LanguageMode {
	this_.EnumBase.Name_ = enumName_
	this_.EnumBase.Ordinal_ = enumOrd_
	_ = featureSet
	this_.initFields_Config_LanguageMode()
	this_.featureSet = featureSet
	return this_
}

var enum_Config_LanguageMode_ECMASCRIPT3 *Config_LanguageMode

func StaticGet_Config_LanguageMode_ECMASCRIPT3() *Config_LanguageMode {
	clinit_Config_LanguageMode()
	return enum_Config_LanguageMode_ECMASCRIPT3
}

var enum_Config_LanguageMode_ECMASCRIPT5 *Config_LanguageMode

func StaticGet_Config_LanguageMode_ECMASCRIPT5() *Config_LanguageMode {
	clinit_Config_LanguageMode()
	return enum_Config_LanguageMode_ECMASCRIPT5
}

var enum_Config_LanguageMode_ECMASCRIPT_2015 *Config_LanguageMode

func StaticGet_Config_LanguageMode_ECMASCRIPT_2015() *Config_LanguageMode {
	clinit_Config_LanguageMode()
	return enum_Config_LanguageMode_ECMASCRIPT_2015
}

var enum_Config_LanguageMode_ECMASCRIPT_2016 *Config_LanguageMode

func StaticGet_Config_LanguageMode_ECMASCRIPT_2016() *Config_LanguageMode {
	clinit_Config_LanguageMode()
	return enum_Config_LanguageMode_ECMASCRIPT_2016
}

var enum_Config_LanguageMode_ECMASCRIPT_2017 *Config_LanguageMode

func StaticGet_Config_LanguageMode_ECMASCRIPT_2017() *Config_LanguageMode {
	clinit_Config_LanguageMode()
	return enum_Config_LanguageMode_ECMASCRIPT_2017
}

var enum_Config_LanguageMode_ECMASCRIPT_2018 *Config_LanguageMode

func StaticGet_Config_LanguageMode_ECMASCRIPT_2018() *Config_LanguageMode {
	clinit_Config_LanguageMode()
	return enum_Config_LanguageMode_ECMASCRIPT_2018
}

var enum_Config_LanguageMode_ECMASCRIPT_2019 *Config_LanguageMode

func StaticGet_Config_LanguageMode_ECMASCRIPT_2019() *Config_LanguageMode {
	clinit_Config_LanguageMode()
	return enum_Config_LanguageMode_ECMASCRIPT_2019
}

var enum_Config_LanguageMode_ECMASCRIPT_2020 *Config_LanguageMode

func StaticGet_Config_LanguageMode_ECMASCRIPT_2020() *Config_LanguageMode {
	clinit_Config_LanguageMode()
	return enum_Config_LanguageMode_ECMASCRIPT_2020
}

var enum_Config_LanguageMode_ECMASCRIPT_2021 *Config_LanguageMode

func StaticGet_Config_LanguageMode_ECMASCRIPT_2021() *Config_LanguageMode {
	clinit_Config_LanguageMode()
	return enum_Config_LanguageMode_ECMASCRIPT_2021
}

var enum_Config_LanguageMode_ECMASCRIPT_2022 *Config_LanguageMode

func StaticGet_Config_LanguageMode_ECMASCRIPT_2022() *Config_LanguageMode {
	clinit_Config_LanguageMode()
	return enum_Config_LanguageMode_ECMASCRIPT_2022
}

var enum_Config_LanguageMode_ES_NEXT *Config_LanguageMode

func StaticGet_Config_LanguageMode_ES_NEXT() *Config_LanguageMode {
	clinit_Config_LanguageMode()
	return enum_Config_LanguageMode_ES_NEXT
}

var enum_Config_LanguageMode_UNSTABLE *Config_LanguageMode

func StaticGet_Config_LanguageMode_UNSTABLE() *Config_LanguageMode {
	clinit_Config_LanguageMode()
	return enum_Config_LanguageMode_UNSTABLE
}

var enum_Config_LanguageMode_UNSUPPORTED *Config_LanguageMode

func StaticGet_Config_LanguageMode_UNSUPPORTED() *Config_LanguageMode {
	clinit_Config_LanguageMode()
	return enum_Config_LanguageMode_UNSUPPORTED
}

var clinitDone_Config_LanguageMode uint32
var clinitState_Config_LanguageMode j2gort.ClinitState

func clinit_Config_LanguageMode() {
	if atomic.LoadUint32(&clinitDone_Config_LanguageMode) != 0 {
		return
	}
	clinitSlow_Config_LanguageMode()
}

func clinitSlow_Config_LanguageMode() {
	if !j2gort.EnterClinit(&clinitState_Config_LanguageMode, &clinitDone_Config_LanguageMode) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_Config_LanguageMode, &clinitDone_Config_LanguageMode)
	enum_Config_LanguageMode_ECMASCRIPT3 = &Config_LanguageMode{}
	initEnumConfig_LanguageMode(enum_Config_LanguageMode_ECMASCRIPT3, "ECMASCRIPT3", int32(0), pkg_parser.StaticGet_FeatureSet_ES3())
	enum_Config_LanguageMode_ECMASCRIPT5 = &Config_LanguageMode{}
	initEnumConfig_LanguageMode(enum_Config_LanguageMode_ECMASCRIPT5, "ECMASCRIPT5", int32(1), pkg_parser.StaticGet_FeatureSet_ES5())
	enum_Config_LanguageMode_ECMASCRIPT_2015 = &Config_LanguageMode{}
	initEnumConfig_LanguageMode(enum_Config_LanguageMode_ECMASCRIPT_2015, "ECMASCRIPT_2015", int32(2), pkg_parser.StaticGet_FeatureSet_ES2015_MODULES())
	enum_Config_LanguageMode_ECMASCRIPT_2016 = &Config_LanguageMode{}
	initEnumConfig_LanguageMode(enum_Config_LanguageMode_ECMASCRIPT_2016, "ECMASCRIPT_2016", int32(3), pkg_parser.StaticGet_FeatureSet_ES2016_MODULES())
	enum_Config_LanguageMode_ECMASCRIPT_2017 = &Config_LanguageMode{}
	initEnumConfig_LanguageMode(enum_Config_LanguageMode_ECMASCRIPT_2017, "ECMASCRIPT_2017", int32(4), pkg_parser.StaticGet_FeatureSet_ES2017_MODULES())
	enum_Config_LanguageMode_ECMASCRIPT_2018 = &Config_LanguageMode{}
	initEnumConfig_LanguageMode(enum_Config_LanguageMode_ECMASCRIPT_2018, "ECMASCRIPT_2018", int32(5), pkg_parser.StaticGet_FeatureSet_ES2018_MODULES())
	enum_Config_LanguageMode_ECMASCRIPT_2019 = &Config_LanguageMode{}
	initEnumConfig_LanguageMode(enum_Config_LanguageMode_ECMASCRIPT_2019, "ECMASCRIPT_2019", int32(6), pkg_parser.StaticGet_FeatureSet_ES2019_MODULES())
	enum_Config_LanguageMode_ECMASCRIPT_2020 = &Config_LanguageMode{}
	initEnumConfig_LanguageMode(enum_Config_LanguageMode_ECMASCRIPT_2020, "ECMASCRIPT_2020", int32(7), pkg_parser.StaticGet_FeatureSet_ES2020_MODULES())
	enum_Config_LanguageMode_ECMASCRIPT_2021 = &Config_LanguageMode{}
	initEnumConfig_LanguageMode(enum_Config_LanguageMode_ECMASCRIPT_2021, "ECMASCRIPT_2021", int32(8), pkg_parser.StaticGet_FeatureSet_ES2021_MODULES())
	enum_Config_LanguageMode_ECMASCRIPT_2022 = &Config_LanguageMode{}
	initEnumConfig_LanguageMode(enum_Config_LanguageMode_ECMASCRIPT_2022, "ECMASCRIPT_2022", int32(9), pkg_parser.StaticGet_FeatureSet_ES2022_MODULES())
	enum_Config_LanguageMode_ES_NEXT = &Config_LanguageMode{}
	initEnumConfig_LanguageMode(enum_Config_LanguageMode_ES_NEXT, "ES_NEXT", int32(10), pkg_parser.StaticGet_FeatureSet_ES_NEXT())
	enum_Config_LanguageMode_UNSTABLE = &Config_LanguageMode{}
	initEnumConfig_LanguageMode(enum_Config_LanguageMode_UNSTABLE, "UNSTABLE", int32(11), pkg_parser.StaticGet_FeatureSet_ES_UNSTABLE())
	enum_Config_LanguageMode_UNSUPPORTED = &Config_LanguageMode{}
	initEnumConfig_LanguageMode(enum_Config_LanguageMode_UNSUPPORTED, "UNSUPPORTED", int32(12), pkg_parser.StaticGet_FeatureSet_ES_UNSUPPORTED())
	enumValues_Config_LanguageMode = []*Config_LanguageMode{enum_Config_LanguageMode_ECMASCRIPT3, enum_Config_LanguageMode_ECMASCRIPT5, enum_Config_LanguageMode_ECMASCRIPT_2015, enum_Config_LanguageMode_ECMASCRIPT_2016, enum_Config_LanguageMode_ECMASCRIPT_2017, enum_Config_LanguageMode_ECMASCRIPT_2018, enum_Config_LanguageMode_ECMASCRIPT_2019, enum_Config_LanguageMode_ECMASCRIPT_2020, enum_Config_LanguageMode_ECMASCRIPT_2021, enum_Config_LanguageMode_ECMASCRIPT_2022, enum_Config_LanguageMode_ES_NEXT, enum_Config_LanguageMode_UNSTABLE, enum_Config_LanguageMode_UNSUPPORTED}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.parsing.Config.LanguageMode", []any{enum_Config_LanguageMode_ECMASCRIPT3, enum_Config_LanguageMode_ECMASCRIPT5, enum_Config_LanguageMode_ECMASCRIPT_2015, enum_Config_LanguageMode_ECMASCRIPT_2016, enum_Config_LanguageMode_ECMASCRIPT_2017, enum_Config_LanguageMode_ECMASCRIPT_2018, enum_Config_LanguageMode_ECMASCRIPT_2019, enum_Config_LanguageMode_ECMASCRIPT_2020, enum_Config_LanguageMode_ECMASCRIPT_2021, enum_Config_LanguageMode_ECMASCRIPT_2022, enum_Config_LanguageMode_ES_NEXT, enum_Config_LanguageMode_UNSTABLE, enum_Config_LanguageMode_UNSUPPORTED})
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.parsing.Config$LanguageMode", []any{enum_Config_LanguageMode_ECMASCRIPT3, enum_Config_LanguageMode_ECMASCRIPT5, enum_Config_LanguageMode_ECMASCRIPT_2015, enum_Config_LanguageMode_ECMASCRIPT_2016, enum_Config_LanguageMode_ECMASCRIPT_2017, enum_Config_LanguageMode_ECMASCRIPT_2018, enum_Config_LanguageMode_ECMASCRIPT_2019, enum_Config_LanguageMode_ECMASCRIPT_2020, enum_Config_LanguageMode_ECMASCRIPT_2021, enum_Config_LanguageMode_ECMASCRIPT_2022, enum_Config_LanguageMode_ES_NEXT, enum_Config_LanguageMode_UNSTABLE, enum_Config_LanguageMode_UNSUPPORTED})
}

func init() {
	j2gort.RegisterEnumConstantsLoader("com.google.javascript.jscomp.parsing.Config.LanguageMode", "com.google.javascript.jscomp.parsing.Config$LanguageMode", func() {
		clinit_Config_LanguageMode()
	})
}

var enumValues_Config_LanguageMode []*Config_LanguageMode

func Config_LanguageMode_Values() []*Config_LanguageMode {
	clinit_Config_LanguageMode()
	return j2gort.CloneSlice(enumValues_Config_LanguageMode)
}

func Config_LanguageMode_ValueOf(name string) *Config_LanguageMode {
	clinit_Config_LanguageMode()
	for _, c := range enumValues_Config_LanguageMode {
		if c.Name() == name {
			return c
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.parsing.Config.LanguageMode." + name))
}

// Java static method: LanguageMode.minimumRequiredFor(com.google.javascript.jscomp.parsing.parser.FeatureSet.Feature feature)
func Config_LanguageMode_MinimumRequiredFor(feature *pkg_parser.FeatureSet_Feature) *Config_LanguageMode {
	clinit_Config_LanguageMode()
	_ = feature
	for _, el_2 := range j2gort.CheckNotNullSlice(Config_LanguageMode_Values()) {
		mode := el_2
		_ = mode
		if mode.featureSet.Has(feature) {
			return mode
		}
	}
	panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "No input language mode supports feature: "+j2gort.StrVal(feature))))
}

// Java static method: LanguageMode.minimumRequiredForSet(com.google.javascript.jscomp.parsing.parser.FeatureSet featureSet)
func Config_LanguageMode_MinimumRequiredForSet(featureSet *pkg_parser.FeatureSet) *Config_LanguageMode {
	clinit_Config_LanguageMode()
	_ = featureSet
	for _, el_3 := range j2gort.CheckNotNullSlice(Config_LanguageMode_Values()) {
		mode := el_3
		_ = mode
		if mode.featureSet.Contains_1_FeatureSet(featureSet) {
			return mode
		}
	}
	panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "No input language mode supports feature set: "+j2gort.StrVal(featureSet))))
}

// Java static method: LanguageMode.latestEcmaScript()
func Config_LanguageMode_LatestEcmaScript() *Config_LanguageMode {
	clinit_Config_LanguageMode()
	return StaticGet_Config_LanguageMode_ECMASCRIPT_2021()
}

func (this_ *Config_LanguageMode) IsJavaNil_() bool { return this_ == nil }
func (*Config_LanguageMode) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.Config$LanguageMode", "LanguageMode"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.Config.LanguageMode", "com.google.javascript.jscomp.parsing.Config$LanguageMode", func(v_ any) bool {
		return j2gort.InstanceOf[*Config_LanguageMode](v_)
	})
}

// Java enum: com.google.javascript.jscomp.parsing.Config.JsDocParsing
type Config_JsDocParsing uint8

const (
	enum_Config_JsDocParsing_TYPES_ONLY                           Config_JsDocParsing = 1
	enum_Config_JsDocParsing_INCLUDE_DESCRIPTIONS_NO_WHITESPACE   Config_JsDocParsing = 2
	enum_Config_JsDocParsing_INCLUDE_DESCRIPTIONS_WITH_WHITESPACE Config_JsDocParsing = 3
	enum_Config_JsDocParsing_INCLUDE_ALL_COMMENTS                 Config_JsDocParsing = 4
	enum_Config_JsDocParsing_LICENSE_COMMENTS_ONLY                Config_JsDocParsing = 5
)

var enumBaseTable_Config_JsDocParsing = [...]j2gort.EnumBase{
	{Name_: "TYPES_ONLY", Ordinal_: 0},
	{Name_: "INCLUDE_DESCRIPTIONS_NO_WHITESPACE", Ordinal_: 1},
	{Name_: "INCLUDE_DESCRIPTIONS_WITH_WHITESPACE", Ordinal_: 2},
	{Name_: "INCLUDE_ALL_COMMENTS", Ordinal_: 3},
	{Name_: "LICENSE_COMMENTS_ONLY", Ordinal_: 4},
}
var enumValues_Config_JsDocParsing = []Config_JsDocParsing{1, 2, 3, 4, 5}

func StaticGet_Config_JsDocParsing_TYPES_ONLY() Config_JsDocParsing {
	return enum_Config_JsDocParsing_TYPES_ONLY
}
func StaticGet_Config_JsDocParsing_INCLUDE_DESCRIPTIONS_NO_WHITESPACE() Config_JsDocParsing {
	return enum_Config_JsDocParsing_INCLUDE_DESCRIPTIONS_NO_WHITESPACE
}
func StaticGet_Config_JsDocParsing_INCLUDE_DESCRIPTIONS_WITH_WHITESPACE() Config_JsDocParsing {
	return enum_Config_JsDocParsing_INCLUDE_DESCRIPTIONS_WITH_WHITESPACE
}
func StaticGet_Config_JsDocParsing_INCLUDE_ALL_COMMENTS() Config_JsDocParsing {
	return enum_Config_JsDocParsing_INCLUDE_ALL_COMMENTS
}
func StaticGet_Config_JsDocParsing_LICENSE_COMMENTS_ONLY() Config_JsDocParsing {
	return enum_Config_JsDocParsing_LICENSE_COMMENTS_ONLY
}

func Config_JsDocParsing_Values() []Config_JsDocParsing {
	return j2gort.CloneSlice(enumValues_Config_JsDocParsing)
}

func Config_JsDocParsing_ValueOf(name string) Config_JsDocParsing {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_Config_JsDocParsing {
		if enumBaseTable_Config_JsDocParsing[i_].Name_ == name {
			return Config_JsDocParsing(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.parsing.Config.JsDocParsing." + name))
}

// Java method: JsDocParsing.shouldParseDescriptions()
func (this_ Config_JsDocParsing) ShouldParseDescriptions() bool {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return (this_ != StaticGet_Config_JsDocParsing_TYPES_ONLY())
}

// Java method: JsDocParsing.shouldPreserveWhitespace()
func (this_ Config_JsDocParsing) ShouldPreserveWhitespace() bool {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return (((this_ == StaticGet_Config_JsDocParsing_INCLUDE_DESCRIPTIONS_WITH_WHITESPACE()) || (this_ == StaticGet_Config_JsDocParsing_INCLUDE_ALL_COMMENTS())) || (this_ == StaticGet_Config_JsDocParsing_LICENSE_COMMENTS_ONLY()))
}

func (this_ Config_JsDocParsing) IsJavaNil_() bool                              { return this_ == 0 }
func (Config_JsDocParsing) HasIdentityEquals_()                                 {}
func (Config_JsDocParsing) Implements_j2gort_Ext_java_lang_Enum()               {}
func (Config_JsDocParsing) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ Config_JsDocParsing) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (Config_JsDocParsing) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.Config$JsDocParsing", "JsDocParsing"
}
func (this_ Config_JsDocParsing) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_Config_JsDocParsing[this_-1]
}
func (this_ Config_JsDocParsing) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ Config_JsDocParsing) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ Config_JsDocParsing) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_Config_JsDocParsing[this_-1].Name_
}
func (this_ Config_JsDocParsing) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_Config_JsDocParsing[this_-1].Name_
}
func (this_ Config_JsDocParsing) Equals(other_ any) bool {
	o_, ok_ := other_.(Config_JsDocParsing)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ Config_JsDocParsing) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_Config_JsDocParsing[this_-1].HashCode()
}
func (this_ Config_JsDocParsing) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_Config_JsDocParsing[this_-1].CompareTo(other_)
}
func (this_ Config_JsDocParsing) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_Config_JsDocParsing[this_-1].Name_)
}
func init() {
	vals_ := []any{Config_JsDocParsing(1), Config_JsDocParsing(2), Config_JsDocParsing(3), Config_JsDocParsing(4), Config_JsDocParsing(5)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.parsing.Config.JsDocParsing", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.parsing.Config$JsDocParsing", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.Config.JsDocParsing", "com.google.javascript.jscomp.parsing.Config$JsDocParsing", func(v_ any) bool {
		return j2gort.InstanceOf[Config_JsDocParsing](v_)
	})
}

// Java enum: com.google.javascript.jscomp.parsing.Config.RunMode
type Config_RunMode uint8

const (
	enum_Config_RunMode_STOP_AFTER_ERROR Config_RunMode = 1
	enum_Config_RunMode_KEEP_GOING       Config_RunMode = 2
)

var enumBaseTable_Config_RunMode = [...]j2gort.EnumBase{
	{Name_: "STOP_AFTER_ERROR", Ordinal_: 0},
	{Name_: "KEEP_GOING", Ordinal_: 1},
}
var enumValues_Config_RunMode = []Config_RunMode{1, 2}

func StaticGet_Config_RunMode_STOP_AFTER_ERROR() Config_RunMode {
	return enum_Config_RunMode_STOP_AFTER_ERROR
}
func StaticGet_Config_RunMode_KEEP_GOING() Config_RunMode {
	return enum_Config_RunMode_KEEP_GOING
}

func Config_RunMode_Values() []Config_RunMode {
	return j2gort.CloneSlice(enumValues_Config_RunMode)
}

func Config_RunMode_ValueOf(name string) Config_RunMode {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_Config_RunMode {
		if enumBaseTable_Config_RunMode[i_].Name_ == name {
			return Config_RunMode(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.parsing.Config.RunMode." + name))
}

func (this_ Config_RunMode) IsJavaNil_() bool                              { return this_ == 0 }
func (Config_RunMode) HasIdentityEquals_()                                 {}
func (Config_RunMode) Implements_j2gort_Ext_java_lang_Enum()               {}
func (Config_RunMode) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ Config_RunMode) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (Config_RunMode) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.Config$RunMode", "RunMode"
}
func (this_ Config_RunMode) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_Config_RunMode[this_-1]
}
func (this_ Config_RunMode) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ Config_RunMode) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ Config_RunMode) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_Config_RunMode[this_-1].Name_
}
func (this_ Config_RunMode) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_Config_RunMode[this_-1].Name_
}
func (this_ Config_RunMode) Equals(other_ any) bool {
	o_, ok_ := other_.(Config_RunMode)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ Config_RunMode) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_Config_RunMode[this_-1].HashCode()
}
func (this_ Config_RunMode) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_Config_RunMode[this_-1].CompareTo(other_)
}
func (this_ Config_RunMode) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_Config_RunMode[this_-1].Name_)
}
func init() {
	vals_ := []any{Config_RunMode(1), Config_RunMode(2)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.parsing.Config.RunMode", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.parsing.Config$RunMode", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.Config.RunMode", "com.google.javascript.jscomp.parsing.Config$RunMode", func(v_ any) bool {
		return j2gort.InstanceOf[Config_RunMode](v_)
	})
}

// Java polymorphic class interface: com.google.javascript.jscomp.parsing.Config.Builder
type Config_Builder interface {
	Implements_parsing_Config_Builder()
	IsJavaNil_() bool
	GetConfig_BuilderBase() *Config_BuilderBase
	SetLanguageMode(*Config_LanguageMode) Config_Builder
	SetStrictMode(Config_StrictMode) Config_Builder
	SetJsDocParsingMode(Config_JsDocParsing) Config_Builder
	SetRunMode(Config_RunMode) Config_Builder
	SetParseInlineSourceMaps(bool) Config_Builder
	SetSuppressionNames(j2gort.Ext_java_lang_Iterable) Config_Builder
	SetClosurePrimitiveNames(j2gort.Ext_java_lang_Iterable) Config_Builder
	SetExtraAnnotationNames(j2gort.Ext_java_lang_Iterable) Config_Builder
	Build() Config
	SetAnnotations(j2gort.Ext_com_google_common_collect_ImmutableMap) Config_Builder
}

func ToIface_Config_Builder[T interface {
	Config_Builder
	comparable
}](v T) Config_Builder {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_Config_Builder[T Config_Builder](v T) Config_Builder {
	var vi Config_Builder = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.jscomp.parsing.Config.Builder
type Config_BuilderBase struct {
	self_   Config_Builder
	idHash_ atomic.Int32
}

func (this_ *Config_BuilderBase) BindSelf_(s any) {
	_ = s
	this_.self_ = s.(Config_Builder)
}
func (this_ *Config_BuilderBase) IsJavaNil_() bool { return this_ == nil }
func (*Config_BuilderBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.Config$Builder", "Builder"
}
func (this_ *Config_BuilderBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.Config.Builder", "com.google.javascript.jscomp.parsing.Config$Builder", func(v_ any) bool {
		return j2gort.InstanceOf[Config_Builder](v_)
	})
}
func (b *Config_BuilderBase) GetConfig_BuilderBase() *Config_BuilderBase { return b }
func (*Config_BuilderBase) Implements_parsing_Config_Builder()           {}

func (this_ *Config_BuilderBase) initFields_Config_Builder() {
	_ = this_
}

func InitConfig_Builder(this_ *Config_BuilderBase) {
	this_.initFields_Config_Builder()
}

// Java method: Builder.setLanguageMode(com.google.javascript.jscomp.parsing.Config.LanguageMode mode)
func (this_ *Config_BuilderBase) SetLanguageMode(mode *Config_LanguageMode) Config_Builder {
	_ = mode
	panic("abstract or native method invoked")
}

// Java method: Builder.setStrictMode(com.google.javascript.jscomp.parsing.Config.StrictMode mode)
func (this_ *Config_BuilderBase) SetStrictMode(mode Config_StrictMode) Config_Builder {
	_ = mode
	panic("abstract or native method invoked")
}

// Java method: Builder.setJsDocParsingMode(com.google.javascript.jscomp.parsing.Config.JsDocParsing mode)
func (this_ *Config_BuilderBase) SetJsDocParsingMode(mode Config_JsDocParsing) Config_Builder {
	_ = mode
	panic("abstract or native method invoked")
}

// Java method: Builder.setRunMode(com.google.javascript.jscomp.parsing.Config.RunMode mode)
func (this_ *Config_BuilderBase) SetRunMode(mode Config_RunMode) Config_Builder {
	_ = mode
	panic("abstract or native method invoked")
}

// Java method: Builder.setParseInlineSourceMaps(boolean parseInlineSourceMaps)
func (this_ *Config_BuilderBase) SetParseInlineSourceMaps(parseInlineSourceMaps bool) Config_Builder {
	_ = parseInlineSourceMaps
	panic("abstract or native method invoked")
}

// Java method: Builder.setSuppressionNames(java.lang.Iterable<java.lang.String> names)
func (this_ *Config_BuilderBase) SetSuppressionNames(names j2gort.Ext_java_lang_Iterable) Config_Builder {
	_ = names
	panic("abstract or native method invoked")
}

// Java method: Builder.setClosurePrimitiveNames(java.lang.Iterable<java.lang.String> names)
func (this_ *Config_BuilderBase) SetClosurePrimitiveNames(names j2gort.Ext_java_lang_Iterable) Config_Builder {
	_ = names
	panic("abstract or native method invoked")
}

// Java method: Builder.setExtraAnnotationNames(java.lang.Iterable<java.lang.String> names)
func (this_ *Config_BuilderBase) SetExtraAnnotationNames(names j2gort.Ext_java_lang_Iterable) Config_Builder {
	_ = names
	return this_.self_.SetAnnotations(Config_BuildAnnotations(names))
}

// Java method: Builder.build()
func (this_ *Config_BuilderBase) Build() Config {
	panic("abstract or native method invoked")
}

// Java method: Builder.setAnnotations(com.google.common.collect.ImmutableMap<java.lang.String,com.google.javascript.jscomp.parsing.Annotation> names)
func (this_ *Config_BuilderBase) SetAnnotations(names j2gort.Ext_com_google_common_collect_ImmutableMap) Config_Builder {
	_ = names
	panic("abstract or native method invoked")
}
