/*
 * Copyright 2015 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/deps/JsFileFullParser.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parsing "github.com/paullepoulpe/jsgo/transpiled/gen/parsing"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.deps.JsFileFullParser
type JsFileFullParser struct {
	idHash_ atomic.Int32
}

func (this_ *JsFileFullParser) BindSelf_(s any) {
	_ = s
}
func (this_ *JsFileFullParser) IsJavaNil_() bool { return this_ == nil }
func (*JsFileFullParser) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.deps.JsFileFullParser", "JsFileFullParser"
}
func (this_ *JsFileFullParser) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*JsFileFullParser) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.deps.JsFileFullParser", "com.google.javascript.jscomp.deps.JsFileFullParser", func(v_ any) bool {
		return j2gort.InstanceOf[*JsFileFullParser](v_)
	})
}

func (this_ *JsFileFullParser) initFields_JsFileFullParser() {
	_ = this_
}

// Java constructor: JsFileFullParser()
func NewJsFileFullParser() *JsFileFullParser {
	this_ := &JsFileFullParser{}
	InitJsFileFullParser(this_)
	return this_
}

func InitJsFileFullParser(this_ *JsFileFullParser) {
	this_.initFields_JsFileFullParser()
}

// Java static method: JsFileFullParser.parse(java.lang.String code, java.lang.String filename, com.google.javascript.jscomp.deps.JsFileFullParser.Reporter reporter)
func JsFileFullParser_Parse(code string, filename string, reporter JsFileFullParser_Reporter) *JsFileFullParser_FileInfo {
	_ = code
	_ = filename
	_ = reporter
	var val_1 JsFileFullParser_DelegatingReporter
	InitJsFileFullParser_DelegatingReporter(&val_1, reporter)
	var errorReporter *JsFileFullParser_DelegatingReporter = &val_1
	_ = errorReporter
	var compiler *Compiler = NewCompiler_1_ErrorManager(NewJsFileFullParser__Anon_106(errorReporter))
	_ = compiler
	var source *SourceFile = SourceFile_FromCode_2(filename, code)
	_ = source
	clinit_CompilerOptions()
	var val_2 CompilerOptions
	InitCompilerOptions(&val_2)
	var options *CompilerOptions = &val_2
	_ = options
	options.SetChecksOnly(true)
	compiler.Init(j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()), j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.ImmListOfPtr(source)), options)
	var config pkg_parsing.Config = pkg_parsing.ParserRunner_CreateConfig_6(pkg_parsing.StaticGet_Config_LanguageMode_ES_NEXT(), pkg_parsing.StaticGet_Config_JsDocParsing_INCLUDE_DESCRIPTIONS_NO_WHITESPACE(), pkg_parsing.StaticGet_Config_RunMode_STOP_AFTER_ERROR(), j2gort.AsIface[j2gort.Ext_java_util_Set](j2gort.Ext_com_google_common_collect_ImmutableSet_Of_0()), true, pkg_parsing.StaticGet_Config_StrictMode_SLOPPY())
	_ = config
	var val_3 JsFileFullParser_FileInfo
	InitJsFileFullParser_FileInfo(&val_3)
	var info *JsFileFullParser_FileInfo = &val_3
	_ = info
	var parsed *pkg_parsing.ParserRunner_ParseResult = pkg_parsing.ParserRunner_Parse(pkg_rhino.ToIface_StaticSourceFile(source), code, config, errorReporter)
	_ = parsed
	if errorReporter.seenError {
		return info
	}
	{
		r_4 := info
		r_4.parsed = parsed
	}
	parsed.FieldGet_ast().SetInputId(pkg_rhino.NewInputId(filename))
	var version string = parsed.FieldGet_features().Version()
	_ = version
	if !(j2gort.CheckNotNullString(version) == "es3") {
		j2gort.MultimapPutStr(info.loadFlags, "lang", j2gort.BoxString(version))
	}
	clinit_GatherModuleMetadata()
	var val_5 GatherModuleMetadata
	InitGatherModuleMetadata(&val_5, compiler, false, StaticGet_ModuleLoader_ResolutionMode_BROWSER())
	var gatherModuleMetadata *GatherModuleMetadata = &val_5
	_ = gatherModuleMetadata
	gatherModuleMetadata.Process_2(pkg_rhino.NewNode_1(pkg_rhino.StaticGet_Token_ROOT()), parsed.FieldGet_ast())
	compiler.GenerateReport()
	if compiler.GetModuleMetadataMap().GetModulesByPath().Size() != int32(1) {
		return info
	}
	var module ModuleMetadataMap_ModuleMetadata = j2gort.AsIface[ModuleMetadataMap_ModuleMetadata](j2gort.Ext_com_google_common_collect_Iterables_GetOnlyElement_1(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableCollection](compiler.GetModuleMetadataMap().GetModulesByPath().Values())))
	_ = module
	base_module_ := module.GetModuleMetadataMap_ModuleMetadataBase()
	if base_module_.IsEs6Module() {
		j2gort.MultimapPutStr(info.loadFlags, "module", "es6")
	} else if base_module_.IsGoogModule() {
		j2gort.MultimapPutStr(info.loadFlags, "module", "goog")
	}
	switch (module.ModuleType()).Ordinal() {
	case int32(4) /* ModuleType.GOOG_PROVIDE */ :
		{
			r_6 := info
			r_6.moduleType = StaticGet_JsFileFullParser_FileInfo_ModuleType_GOOG_PROVIDE()
		}
	case int32(1) /* ModuleType.GOOG_MODULE */ :
		{
			r_7 := info
			r_7.moduleType = StaticGet_JsFileFullParser_FileInfo_ModuleType_GOOG_MODULE()
		}
	case int32(2) /* ModuleType.LEGACY_GOOG_MODULE */ :
		{
			r_8 := info
			r_8.moduleType = StaticGet_JsFileFullParser_FileInfo_ModuleType_GOOG_MODULE()
		}
		{
			r_9 := info
			r_9.isLegacyNamespace = true
		}
	case int32(0) /* ModuleType.ES6_MODULE */ :
		{
			r_10 := info
			r_10.moduleType = StaticGet_JsFileFullParser_FileInfo_ModuleType_ES_MODULE()
		}
	case int32(3) /* ModuleType.COMMON_JS */, int32(5) /* ModuleType.SCRIPT */ :
		{
			r_11 := info
			r_11.moduleType = StaticGet_JsFileFullParser_FileInfo_ModuleType_UNKNOWN()
		}
	}
	{
		r_12 := info
		r_12.goog = module.UsesClosure()
	}
	JsFileFullParser_RecordModuleMetadata(info, module)
	return info
}

// Java static method: JsFileFullParser.recordModuleMetadata(com.google.javascript.jscomp.deps.JsFileFullParser.FileInfo info, com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleMetadata module)
func JsFileFullParser_RecordModuleMetadata(info *JsFileFullParser_FileInfo, module ModuleMetadataMap_ModuleMetadata) {
	_ = info
	_ = module
	info.importedModules.AddAll_1(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](module.Es6ImportSpecifiers().ElementSet()))
	if module.UsesClosure() {
		info.provides.AddAll_1(j2gort.AsIface[j2gort.Ext_java_util_Collection](module.GoogNamespaces()))
		info.requires.AddAll_1(j2gort.AsIface[j2gort.Ext_java_util_Collection](module.StronglyRequiredGoogNamespaces()))
		info.maybeRequires.AddAll_1(j2gort.AsIface[j2gort.Ext_java_util_Collection](module.MaybeRequiredGoogNamespaces()))
		info.dynamicRequires.AddAll_1(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](module.DynamicallyRequiredGoogNamespaces().ElementSet()))
		info.readToggles.AddAll_1(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](module.ReadToggles().ElementSet()))
		info.typeRequires.AddAll_1(j2gort.AsIface[j2gort.Ext_java_util_Collection](module.WeaklyRequiredGoogNamespaces()))
		{
			r_13 := info
			r_13.testonly = module.IsTestOnly()
		}
	}
	{
		cur_14 := j2gort.NewIterCursor(j2gort.CheckNotNull(module.NestedModules()))
		for cur_14.HasNext() {
			rawElem_15 := cur_14.Next()
			nested := j2gort.AsIface[ModuleMetadataMap_ModuleMetadata](rawElem_15)
			_ = nested
			JsFileFullParser_RecordModuleMetadata(info, nested)
		}
	}
}

// Java class: com.google.javascript.jscomp.deps.JsFileFullParser.FileInfo
type JsFileFullParser_FileInfo struct {
	goog              bool                                          `gson:"goog"`
	testonly          bool                                          `gson:"testonly"`
	moduleType        JsFileFullParser_FileInfo_ModuleType          `gson:"moduleType"`
	isLegacyNamespace bool                                          `gson:"isLegacyNamespace"`
	importedModules   j2gort.Ext_java_util_Set                      `gson:"importedModules"`
	dynamicRequires   j2gort.Ext_java_util_Set                      `gson:"dynamicRequires"`
	readToggles       j2gort.Ext_java_util_Set                      `gson:"readToggles"`
	provides          j2gort.Ext_com_google_common_collect_Multiset `gson:"provides"`
	requires          j2gort.Ext_com_google_common_collect_Multiset `gson:"requires"`
	typeRequires      j2gort.Ext_com_google_common_collect_Multiset `gson:"typeRequires"`
	maybeRequires     j2gort.Ext_com_google_common_collect_Multiset `gson:"maybeRequires"`
	loadFlags         j2gort.Ext_com_google_common_collect_Multimap `gson:"loadFlags"`
	parsed            *pkg_parsing.ParserRunner_ParseResult         `gson:"parsed"`
	idHash_           atomic.Int32
}

func (this_ *JsFileFullParser_FileInfo) BindSelf_(s any) {
	_ = s
}
func (this_ *JsFileFullParser_FileInfo) IsJavaNil_() bool { return this_ == nil }
func (*JsFileFullParser_FileInfo) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.deps.JsFileFullParser$FileInfo", "FileInfo"
}
func (this_ *JsFileFullParser_FileInfo) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*JsFileFullParser_FileInfo) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.deps.JsFileFullParser.FileInfo", "com.google.javascript.jscomp.deps.JsFileFullParser$FileInfo", func(v_ any) bool {
		return j2gort.InstanceOf[*JsFileFullParser_FileInfo](v_)
	})
}

func (this_ *JsFileFullParser_FileInfo) initFields_JsFileFullParser_FileInfo() {
	_ = this_

	this_.moduleType = StaticGet_JsFileFullParser_FileInfo_ModuleType_UNKNOWN()

	this_.importedModules = j2gort.NewExt_java_util_LinkedHashSet()
	this_.dynamicRequires = j2gort.NewExt_java_util_LinkedHashSet()
	this_.readToggles = j2gort.NewExt_java_util_LinkedHashSet()
	this_.provides = j2gort.AsIface[j2gort.Ext_com_google_common_collect_Multiset](j2gort.Ext_com_google_common_collect_TreeMultiset_Create_0())
	this_.requires = j2gort.AsIface[j2gort.Ext_com_google_common_collect_Multiset](j2gort.Ext_com_google_common_collect_TreeMultiset_Create_0())
	this_.typeRequires = j2gort.AsIface[j2gort.Ext_com_google_common_collect_Multiset](j2gort.Ext_com_google_common_collect_TreeMultiset_Create_0())
	this_.maybeRequires = j2gort.AsIface[j2gort.Ext_com_google_common_collect_Multiset](j2gort.Ext_com_google_common_collect_TreeMultiset_Create_0())
	this_.loadFlags = j2gort.AsIface[j2gort.Ext_com_google_common_collect_Multimap](j2gort.Ext_com_google_common_collect_TreeMultimap_Create_0())

}

// Java constructor: FileInfo()
func NewJsFileFullParser_FileInfo() *JsFileFullParser_FileInfo {
	this_ := &JsFileFullParser_FileInfo{}
	InitJsFileFullParser_FileInfo(this_)
	return this_
}

func InitJsFileFullParser_FileInfo(this_ *JsFileFullParser_FileInfo) {
	this_.initFields_JsFileFullParser_FileInfo()
}

// Java enum: com.google.javascript.jscomp.deps.JsFileFullParser.FileInfo.ModuleType
type JsFileFullParser_FileInfo_ModuleType uint8

const (
	enum_JsFileFullParser_FileInfo_ModuleType_UNKNOWN      JsFileFullParser_FileInfo_ModuleType = 1
	enum_JsFileFullParser_FileInfo_ModuleType_GOOG_PROVIDE JsFileFullParser_FileInfo_ModuleType = 2
	enum_JsFileFullParser_FileInfo_ModuleType_GOOG_MODULE  JsFileFullParser_FileInfo_ModuleType = 3
	enum_JsFileFullParser_FileInfo_ModuleType_ES_MODULE    JsFileFullParser_FileInfo_ModuleType = 4
)

var enumBaseTable_JsFileFullParser_FileInfo_ModuleType = [...]j2gort.EnumBase{
	{Name_: "UNKNOWN", Ordinal_: 0},
	{Name_: "GOOG_PROVIDE", Ordinal_: 1},
	{Name_: "GOOG_MODULE", Ordinal_: 2},
	{Name_: "ES_MODULE", Ordinal_: 3},
}
var enumValues_JsFileFullParser_FileInfo_ModuleType = []JsFileFullParser_FileInfo_ModuleType{1, 2, 3, 4}

func StaticGet_JsFileFullParser_FileInfo_ModuleType_UNKNOWN() JsFileFullParser_FileInfo_ModuleType {
	return enum_JsFileFullParser_FileInfo_ModuleType_UNKNOWN
}
func StaticGet_JsFileFullParser_FileInfo_ModuleType_GOOG_PROVIDE() JsFileFullParser_FileInfo_ModuleType {
	return enum_JsFileFullParser_FileInfo_ModuleType_GOOG_PROVIDE
}
func StaticGet_JsFileFullParser_FileInfo_ModuleType_GOOG_MODULE() JsFileFullParser_FileInfo_ModuleType {
	return enum_JsFileFullParser_FileInfo_ModuleType_GOOG_MODULE
}
func StaticGet_JsFileFullParser_FileInfo_ModuleType_ES_MODULE() JsFileFullParser_FileInfo_ModuleType {
	return enum_JsFileFullParser_FileInfo_ModuleType_ES_MODULE
}

func JsFileFullParser_FileInfo_ModuleType_Values() []JsFileFullParser_FileInfo_ModuleType {
	return j2gort.CloneSlice(enumValues_JsFileFullParser_FileInfo_ModuleType)
}

func JsFileFullParser_FileInfo_ModuleType_ValueOf(name string) JsFileFullParser_FileInfo_ModuleType {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_JsFileFullParser_FileInfo_ModuleType {
		if enumBaseTable_JsFileFullParser_FileInfo_ModuleType[i_].Name_ == name {
			return JsFileFullParser_FileInfo_ModuleType(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.deps.JsFileFullParser.FileInfo.ModuleType." + name))
}

func (this_ JsFileFullParser_FileInfo_ModuleType) IsJavaNil_() bool                              { return this_ == 0 }
func (JsFileFullParser_FileInfo_ModuleType) HasIdentityEquals_()                                 {}
func (JsFileFullParser_FileInfo_ModuleType) Implements_j2gort_Ext_java_lang_Enum()               {}
func (JsFileFullParser_FileInfo_ModuleType) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ JsFileFullParser_FileInfo_ModuleType) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (JsFileFullParser_FileInfo_ModuleType) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.deps.JsFileFullParser$FileInfo$ModuleType", "ModuleType"
}
func (this_ JsFileFullParser_FileInfo_ModuleType) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_JsFileFullParser_FileInfo_ModuleType[this_-1]
}
func (this_ JsFileFullParser_FileInfo_ModuleType) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ JsFileFullParser_FileInfo_ModuleType) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ JsFileFullParser_FileInfo_ModuleType) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_JsFileFullParser_FileInfo_ModuleType[this_-1].Name_
}
func (this_ JsFileFullParser_FileInfo_ModuleType) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_JsFileFullParser_FileInfo_ModuleType[this_-1].Name_
}
func (this_ JsFileFullParser_FileInfo_ModuleType) Equals(other_ any) bool {
	o_, ok_ := other_.(JsFileFullParser_FileInfo_ModuleType)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ JsFileFullParser_FileInfo_ModuleType) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_JsFileFullParser_FileInfo_ModuleType[this_-1].HashCode()
}
func (this_ JsFileFullParser_FileInfo_ModuleType) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_JsFileFullParser_FileInfo_ModuleType[this_-1].CompareTo(other_)
}
func (this_ JsFileFullParser_FileInfo_ModuleType) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_JsFileFullParser_FileInfo_ModuleType[this_-1].Name_)
}
func init() {
	vals_ := []any{JsFileFullParser_FileInfo_ModuleType(1), JsFileFullParser_FileInfo_ModuleType(2), JsFileFullParser_FileInfo_ModuleType(3), JsFileFullParser_FileInfo_ModuleType(4)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.deps.JsFileFullParser.FileInfo.ModuleType", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.deps.JsFileFullParser$FileInfo$ModuleType", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.deps.JsFileFullParser.FileInfo.ModuleType", "com.google.javascript.jscomp.deps.JsFileFullParser$FileInfo$ModuleType", func(v_ any) bool {
		return j2gort.InstanceOf[JsFileFullParser_FileInfo_ModuleType](v_)
	})
}

// Java interface: com.google.javascript.jscomp.deps.JsFileFullParser.Reporter
type JsFileFullParser_Reporter interface {
	Implements_jscomp_JsFileFullParser_Reporter()
	IsJavaNil_() bool
	Report(bool, string, string, int32, int32)
}

func ToIface_JsFileFullParser_Reporter[T interface {
	JsFileFullParser_Reporter
	comparable
}](v T) JsFileFullParser_Reporter {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_JsFileFullParser_Reporter[T JsFileFullParser_Reporter](v T) JsFileFullParser_Reporter {
	var vi JsFileFullParser_Reporter = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.deps.JsFileFullParser.Reporter", "com.google.javascript.jscomp.deps.JsFileFullParser$Reporter", func(v_ any) bool {
		return j2gort.InstanceOf[JsFileFullParser_Reporter](v_)
	})
}

// Java class: com.google.javascript.jscomp.deps.JsFileFullParser.DelegatingReporter
type JsFileFullParser_DelegatingReporter struct {
	delegate  JsFileFullParser_Reporter `gson:"delegate"`
	seenError bool                      `gson:"seenError"`
	idHash_   atomic.Int32
}

func (this_ *JsFileFullParser_DelegatingReporter) BindSelf_(s any) {
	_ = s
}
func (this_ *JsFileFullParser_DelegatingReporter) IsJavaNil_() bool { return this_ == nil }
func (*JsFileFullParser_DelegatingReporter) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.deps.JsFileFullParser$DelegatingReporter", "DelegatingReporter"
}
func (this_ *JsFileFullParser_DelegatingReporter) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*JsFileFullParser_DelegatingReporter) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.deps.JsFileFullParser.DelegatingReporter", "com.google.javascript.jscomp.deps.JsFileFullParser$DelegatingReporter", func(v_ any) bool {
		return j2gort.InstanceOf[*JsFileFullParser_DelegatingReporter](v_)
	})
}
func (*JsFileFullParser_DelegatingReporter) Implements_rhino_ErrorReporter() {}

func (this_ *JsFileFullParser_DelegatingReporter) initFields_JsFileFullParser_DelegatingReporter() {
	_ = this_

}

// Java constructor: DelegatingReporter(com.google.javascript.jscomp.deps.JsFileFullParser.Reporter delegate)
func NewJsFileFullParser_DelegatingReporter(delegate JsFileFullParser_Reporter) *JsFileFullParser_DelegatingReporter {
	this_ := &JsFileFullParser_DelegatingReporter{}
	InitJsFileFullParser_DelegatingReporter(this_, delegate)
	return this_
}

func InitJsFileFullParser_DelegatingReporter(this_ *JsFileFullParser_DelegatingReporter, delegate JsFileFullParser_Reporter) {
	_ = delegate
	this_.initFields_JsFileFullParser_DelegatingReporter()
	this_.delegate = j2gort.CheckNotNullIface(delegate)
}

// Java method: DelegatingReporter.warning(java.lang.String message, java.lang.String sourceName, int line, int lineOffset)
func (this_ *JsFileFullParser_DelegatingReporter) Warning(message string, sourceName string, line int32, lineOffset int32) {
	_ = message
	_ = sourceName
	_ = line
	_ = lineOffset
	this_.delegate.Report(false, message, sourceName, line, lineOffset)
}

// Java method: DelegatingReporter.error(java.lang.String message, java.lang.String sourceName, int line, int lineOffset)
func (this_ *JsFileFullParser_DelegatingReporter) M_Error(message string, sourceName string, line int32, lineOffset int32) {
	_ = message
	_ = sourceName
	_ = line
	_ = lineOffset
	this_.seenError = true
	this_.delegate.Report(true, message, sourceName, line, lineOffset)
}

// Java class:
type JsFileFullParser__Anon_106 struct {
	BasicErrorManagerBase
	cap_0_errorReporter *JsFileFullParser_DelegatingReporter
}

func (this_ *JsFileFullParser__Anon_106) BindSelf_(s any) {
	_ = s
	this_.BasicErrorManagerBase.BindSelf_(s)
}
func (this_ *JsFileFullParser__Anon_106) IsJavaNil_() bool { return this_ == nil }
func (*JsFileFullParser__Anon_106) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.deps.JsFileFullParser$1", ""
}
func (this_ *JsFileFullParser__Anon_106) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.BasicErrorManagerBase.IdentityHash_()
}
func (*JsFileFullParser__Anon_106) HasIdentityEquals_()                    {}
func (*JsFileFullParser__Anon_106) Implements_jscomp_BasicErrorManager()   {}
func (*JsFileFullParser__Anon_106) Implements_jscomp_SortingErrorManager() {}
func (*JsFileFullParser__Anon_106) Implements_jscomp_ErrorManager()        {}
func (*JsFileFullParser__Anon_106) Implements_jscomp_ErrorHandler()        {}

func (this_ *JsFileFullParser__Anon_106) initFields_JsFileFullParser__Anon_106() {
	_ = this_
}

func NewJsFileFullParser__Anon_106(cap_p_errorReporter *JsFileFullParser_DelegatingReporter) *JsFileFullParser__Anon_106 {
	this_ := &JsFileFullParser__Anon_106{}
	this_.BindSelf_(this_)
	this_.cap_0_errorReporter = cap_p_errorReporter
	InitBasicErrorManager(&this_.BasicErrorManagerBase)
	this_.initFields_JsFileFullParser__Anon_106()
	return this_
}

// Java method: .println(com.google.javascript.jscomp.CheckLevel level, com.google.javascript.jscomp.JSError error)
func (this_ *JsFileFullParser__Anon_106) Println(level *CheckLevel, error_ *JSError) {
	_ = level
	_ = error_
	if level == StaticGet_CheckLevel_ERROR() {
		this_.cap_0_errorReporter.M_Error(error_.Description(), error_.SourceName(), error_.GetLineNumber(), error_.Charno())
	} else if level == StaticGet_CheckLevel_WARNING() {
		this_.cap_0_errorReporter.Warning(error_.Description(), error_.SourceName(), error_.GetLineNumber(), error_.Charno())
	}
}

// Java method: .printSummary()
func (this_ *JsFileFullParser__Anon_106) PrintSummary() {
}
