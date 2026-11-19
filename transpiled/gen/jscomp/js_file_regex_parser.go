/*
 * Copyright 2008 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/deps/JsFileRegexParser.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.deps.JsFileRegexParser
type JsFileRegexParser struct {
	JsFileLineParserBase
	googMatcher               *j2gort.Ext_java_util_regex_Matcher `gson:"googMatcher"`
	es6Matcher                *j2gort.Ext_java_util_regex_Matcher `gson:"es6Matcher"`
	googMatcherBuffer         string                              `gson:"googMatcherBuffer"`
	provides                  j2gort.Ext_java_util_List           `gson:"provides"`
	requires                  j2gort.Ext_java_util_List           `gson:"requires"`
	typeRequires              j2gort.Ext_java_util_List           `gson:"typeRequires"`
	fileHasProvidesOrRequires bool                                `gson:"fileHasProvidesOrRequires"`
	hasExternsAnnotation      bool                                `gson:"hasExternsAnnotation"`
	hasNoCompileAnnotation    bool                                `gson:"hasNoCompileAnnotation"`
	loader                    *ModuleLoader                       `gson:"loader"`
	file                      *ModuleLoader_ModulePath            `gson:"file"`
	moduleType                JsFileRegexParser_ModuleType        `gson:"moduleType"`
	seenLoadModule            bool                                `gson:"seenLoadModule"`
	includeGoogBase           bool                                `gson:"includeGoogBase"`
}

func (this_ *JsFileRegexParser) BindSelf_(s any) {
	_ = s
	this_.JsFileLineParserBase.BindSelf_(s)
}
func (this_ *JsFileRegexParser) IsJavaNil_() bool { return this_ == nil }
func (*JsFileRegexParser) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.deps.JsFileRegexParser", "JsFileRegexParser"
}
func (this_ *JsFileRegexParser) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.JsFileLineParserBase.IdentityHash_()
}
func (*JsFileRegexParser) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.deps.JsFileRegexParser", "com.google.javascript.jscomp.deps.JsFileRegexParser", func(v_ any) bool {
		return j2gort.InstanceOf[*JsFileRegexParser](v_)
	})
}
func (*JsFileRegexParser) Implements_jscomp_JsFileLineParser() {}

var static_JsFileRegexParser_logger *j2gort.Ext_java_util_logging_Logger

func StaticGet_JsFileRegexParser_Logger() *j2gort.Ext_java_util_logging_Logger {
	clinit_JsFileRegexParser()
	return static_JsFileRegexParser_logger
}
func StaticSet_JsFileRegexParser_Logger(v *j2gort.Ext_java_util_logging_Logger) {
	clinit_JsFileRegexParser()
	static_JsFileRegexParser_logger = v
}

var static_JsFileRegexParser_GOOG_PROVIDE_REQUIRE_PATTERN *j2gort.Ext_java_util_regex_Pattern

func StaticGet_JsFileRegexParser_GOOG_PROVIDE_REQUIRE_PATTERN() *j2gort.Ext_java_util_regex_Pattern {
	clinit_JsFileRegexParser()
	return static_JsFileRegexParser_GOOG_PROVIDE_REQUIRE_PATTERN
}
func StaticSet_JsFileRegexParser_GOOG_PROVIDE_REQUIRE_PATTERN(v *j2gort.Ext_java_util_regex_Pattern) {
	clinit_JsFileRegexParser()
	static_JsFileRegexParser_GOOG_PROVIDE_REQUIRE_PATTERN = v
}

var static_JsFileRegexParser_ES6_MODULE_PATTERN *j2gort.Ext_java_util_regex_Pattern

func StaticGet_JsFileRegexParser_ES6_MODULE_PATTERN() *j2gort.Ext_java_util_regex_Pattern {
	clinit_JsFileRegexParser()
	return static_JsFileRegexParser_ES6_MODULE_PATTERN
}
func StaticSet_JsFileRegexParser_ES6_MODULE_PATTERN(v *j2gort.Ext_java_util_regex_Pattern) {
	clinit_JsFileRegexParser()
	static_JsFileRegexParser_ES6_MODULE_PATTERN = v
}

var static_JsFileRegexParser_ES6_EXPORT_PATTERN *j2gort.Ext_java_util_regex_Pattern

func StaticGet_JsFileRegexParser_ES6_EXPORT_PATTERN() *j2gort.Ext_java_util_regex_Pattern {
	clinit_JsFileRegexParser()
	return static_JsFileRegexParser_ES6_EXPORT_PATTERN
}
func StaticSet_JsFileRegexParser_ES6_EXPORT_PATTERN(v *j2gort.Ext_java_util_regex_Pattern) {
	clinit_JsFileRegexParser()
	static_JsFileRegexParser_ES6_EXPORT_PATTERN = v
}
func StaticGet_JsFileRegexParser_PROVIDES_GOOG_COMMENT() string { return "@provideGoog" }
func StaticGet_JsFileRegexParser_EXTERNS_COMMENT() string       { return "@externs" }
func StaticGet_JsFileRegexParser_NOCOMPILE_COMMENT() string     { return "@nocompile" }
func StaticGet_JsFileRegexParser_BUNDLED_GOOG_MODULE_START() string {
	return "goog.loadModule(function("
}

var clinitDone_JsFileRegexParser uint32
var clinitState_JsFileRegexParser j2gort.ClinitState

func clinit_JsFileRegexParser() {
	if atomic.LoadUint32(&clinitDone_JsFileRegexParser) != 0 {
		return
	}
	clinitSlow_JsFileRegexParser()
}

func clinitSlow_JsFileRegexParser() {
	if !j2gort.EnterClinit(&clinitState_JsFileRegexParser, &clinitDone_JsFileRegexParser) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_JsFileRegexParser, &clinitDone_JsFileRegexParser)
	static_JsFileRegexParser_logger = j2gort.Ext_java_util_logging_Logger_GetLogger_1(j2gort.CheckNotNullPtr(j2gort.ClassLiteral("com.google.javascript.jscomp.deps.JsFileRegexParser", "JsFileRegexParser")).GetName())
	static_JsFileRegexParser_GOOG_PROVIDE_REQUIRE_PATTERN = j2gort.Ext_java_util_regex_Pattern_Compile_1("(?:^|;)(?:[\\p{L}\\p{Nl}\\p{Nd}$_,:{}\\s]+=)?\\s*goog\\.(?<func>provide|module|require|requireType|forwardDeclare|addDependency|declareModuleId)\\s*\\(\\s*(?<args>.*?)\\s*\\)")
	static_JsFileRegexParser_ES6_MODULE_PATTERN = j2gort.Ext_java_util_regex_Pattern_Compile_1("^(?:import|export)\\b\\s*(?:[a-zA-Z0-9$_*,{}\\s]+\\bfrom\\s*|)(?:['\"]([^'\"]+)['\"])?\\s*;")
	static_JsFileRegexParser_ES6_EXPORT_PATTERN = j2gort.Ext_java_util_regex_Pattern_Compile_1("^export\\b")
}

func (this_ *JsFileRegexParser) initFields_JsFileRegexParser() {
	_ = this_
	this_.googMatcher = StaticGet_JsFileRegexParser_GOOG_PROVIDE_REQUIRE_PATTERN().Matcher("")
	this_.es6Matcher = StaticGet_JsFileRegexParser_ES6_MODULE_PATTERN().Matcher("")

	this_.loader = StaticGet_ModuleLoader_EMPTY()

}

// Java constructor: JsFileRegexParser(com.google.javascript.jscomp.ErrorManager errorManager)
func NewJsFileRegexParser(errorManager ErrorManager) *JsFileRegexParser {
	clinit_JsFileRegexParser()
	this_ := &JsFileRegexParser{}
	this_.BindSelf_(this_)
	InitJsFileRegexParser(this_, errorManager)
	return this_
}

func InitJsFileRegexParser(this_ *JsFileRegexParser, errorManager ErrorManager) {
	_ = errorManager
	InitJsFileLineParser(&this_.JsFileLineParserBase, errorManager)
	this_.initFields_JsFileRegexParser()
}

// Java method: JsFileRegexParser.setIncludeGoogBase(boolean include)
func (this_ *JsFileRegexParser) SetIncludeGoogBase(include bool) *JsFileRegexParser {
	_ = include
	if !(JsFileRegexParser_IsSupported()) {
		j2gort.CheckState(false)
	}
	this_.includeGoogBase = include
	return this_
}

// Java method: JsFileRegexParser.setModuleLoader(com.google.javascript.jscomp.deps.ModuleLoader loader)
func (this_ *JsFileRegexParser) SetModuleLoader(loader *ModuleLoader) *JsFileRegexParser {
	_ = loader
	this_.loader = loader
	return this_
}

// Java method: JsFileRegexParser.parseFile(java.lang.String filePath, java.lang.String closureRelativePath, java.lang.String fileContents)
func (this_ *JsFileRegexParser) ParseFile(filePath string, closureRelativePath string, fileContents string) DependencyInfo {
	_ = filePath
	_ = closureRelativePath
	_ = fileContents
	return this_.priv_JsFileRegexParser_parseReader_6100(filePath, closureRelativePath, j2gort.NewExt_java_io_StringReader(fileContents))
}

// Java method: JsFileRegexParser.parseReader(java.lang.String filePath, java.lang.String closureRelativePath, java.io.Reader fileContents)
func (this_ *JsFileRegexParser) priv_JsFileRegexParser_parseReader_6100(filePath string, closureRelativePath string, fileContents j2gort.Ext_java_io_Reader) DependencyInfo {
	_ = filePath
	_ = closureRelativePath
	_ = fileContents
	this_.provides = j2gort.NewExt_java_util_ArrayList()
	this_.requires = j2gort.NewExt_java_util_ArrayList()
	this_.typeRequires = j2gort.NewExt_java_util_ArrayList()
	this_.fileHasProvidesOrRequires = false
	this_.hasExternsAnnotation = false
	this_.hasNoCompileAnnotation = false
	this_.file = this_.loader.Resolve(filePath)
	this_.moduleType = StaticGet_JsFileRegexParser_ModuleType_NON_MODULE()
	this_.seenLoadModule = false
	this_.googMatcherBuffer = ""
	if StaticGet_JsFileRegexParser_Logger().IsLoggable(j2gort.Ext_java_util_logging_Level_F_FINE) {
		StaticGet_JsFileRegexParser_Logger().Fine_1_string(("Parsing Source: " + j2gort.StrValString(filePath)))
	}
	this_.DoParse(filePath, fileContents)
	if this_.moduleType == StaticGet_JsFileRegexParser_ModuleType_ES6_MODULE() {
		this_.provides.Add_1(j2gort.BoxString(this_.file.ToModuleName()))
	}
	var loadFlags j2gort.Ext_java_util_Map = j2gort.NewExt_java_util_LinkedHashMap()
	_ = loadFlags
	switch (this_.moduleType).Ordinal() {
	case int32(1) /* ModuleType.GOOG_MODULE */ :
		j2gort.MapPutStr(loadFlags, "module", "goog")
	case int32(3) /* ModuleType.ES6_MODULE */ :
		j2gort.MapPutStr(loadFlags, "module", "es6")
	default:
	}
	var dependencyInfo DependencyInfo = j2gort.AsIface[DependencyInfo](SimpleDependencyInfo_S_Builder(closureRelativePath, filePath).SetProvides_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.provides)).SetRequires_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.requires)).SetTypeRequires_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.typeRequires)).SetLoadFlags(loadFlags).SetHasExternsAnnotation(this_.hasExternsAnnotation).SetHasNoCompileAnnotation(this_.hasNoCompileAnnotation).Build())
	_ = dependencyInfo
	if StaticGet_JsFileRegexParser_Logger().IsLoggable(j2gort.Ext_java_util_logging_Level_F_FINE) {
		StaticGet_JsFileRegexParser_Logger().Fine_1_string(("DepInfo: " + j2gort.StrVal(dependencyInfo)))
	}
	return dependencyInfo
}

// Java method: JsFileRegexParser.setModuleType(com.google.javascript.jscomp.deps.JsFileRegexParser.ModuleType type)
func (this_ *JsFileRegexParser) priv_JsFileRegexParser_setModuleType_6101(type_ JsFileRegexParser_ModuleType) bool {
	_ = type_
	var provide bool = ((type_ == StaticGet_JsFileRegexParser_ModuleType_GOOG_PROVIDE()) || (this_.moduleType == StaticGet_JsFileRegexParser_ModuleType_GOOG_PROVIDE()))
	_ = provide
	var es6Module bool = ((type_ == StaticGet_JsFileRegexParser_ModuleType_ES6_MODULE()) || (this_.moduleType == StaticGet_JsFileRegexParser_ModuleType_ES6_MODULE()))
	_ = es6Module
	var googModule bool = ((type_ == StaticGet_JsFileRegexParser_ModuleType_GOOG_MODULE()) || (this_.moduleType == StaticGet_JsFileRegexParser_ModuleType_GOOG_MODULE()))
	_ = googModule
	if (googModule && provide) && this_.seenLoadModule {
		this_.moduleType = StaticGet_JsFileRegexParser_ModuleType_GOOG_PROVIDE()
		return true
	}
	var provideGoogModuleConflict bool = ((googModule && provide) && !(this_.seenLoadModule))
	_ = provideGoogModuleConflict
	var provideEs6ModuleConflict bool = (es6Module && provide)
	_ = provideEs6ModuleConflict
	var googEs6ModuleConflict bool = ((googModule || this_.seenLoadModule) && es6Module)
	_ = googEs6ModuleConflict
	if (provideGoogModuleConflict || provideEs6ModuleConflict) || googEs6ModuleConflict {
		this_.GetJsFileLineParserBase().errorManager.Report_2(StaticGet_CheckLevel_WARNING(), JSError_Make_5(this_.file.String(), int32(-1), int32(-1), StaticGet_ModuleLoader_MODULE_CONFLICT(), []string{this_.file.String()}))
		return false
	}
	this_.moduleType = type_
	return true
}

// Java method: JsFileRegexParser.parseJsDocCommentLine(java.lang.String line)
func (this_ *JsFileRegexParser) ParseJsDocCommentLine(line string) bool {
	_ = line
	if this_.includeGoogBase && (j2gort.JStrIndexOf(line, "@provideGoog") >= 0) {
		this_.provides.Add_1("goog")
		return false
	} else if j2gort.JStrIndexOf(line, "@externs") >= 0 {
		this_.hasExternsAnnotation = true
		return false
	} else if j2gort.JStrIndexOf(line, "@nocompile") >= 0 {
		this_.hasNoCompileAnnotation = true
		return false
	}
	return true
}

// Java method: JsFileRegexParser.applyGoogMatcher(java.lang.String line)
func (this_ *JsFileRegexParser) priv_JsFileRegexParser_applyGoogMatcher_6102(line string) bool {
	_ = line
	var lineHasProvidesOrRequires bool = false
	_ = lineHasProvidesOrRequires
	this_.googMatcher.Reset_1(j2gort.BoxString(line))
	for this_.googMatcher.Find_0() {
		lineHasProvidesOrRequires = true
		if this_.includeGoogBase && !(this_.fileHasProvidesOrRequires) {
			this_.fileHasProvidesOrRequires = true
			this_.requires.Add_1(j2gort.PtrToAny(StaticGet_DependencyInfo_Require_BASE()))
		}
		var methodName string = this_.googMatcher.Group_1_string("func")
		_ = methodName
		var firstChar uint16 = j2gort.JStrCharAt(methodName, int32(0))
		_ = firstChar
		var isDeclareModuleNamespace bool = (int32(firstChar) == int32(uint16(100)))
		_ = isDeclareModuleNamespace
		var isModule bool = (!(isDeclareModuleNamespace) && (int32(firstChar) == int32(uint16(109))))
		_ = isModule
		var isProvide bool = (int32(firstChar) == int32(uint16(112)))
		_ = isProvide
		var providesNamespace bool = ((isProvide || isModule) || isDeclareModuleNamespace)
		_ = providesNamespace
		var isRequire bool = ((int32(firstChar) == int32(uint16(114))) || (int32(firstChar) == int32(uint16(102))))
		_ = isRequire
		if isModule && !(this_.seenLoadModule) {
			providesNamespace = this_.priv_JsFileRegexParser_setModuleType_6101(StaticGet_JsFileRegexParser_ModuleType_GOOG_MODULE())
		}
		if isProvide {
			providesNamespace = this_.priv_JsFileRegexParser_setModuleType_6101(StaticGet_JsFileRegexParser_ModuleType_GOOG_PROVIDE())
		}
		if providesNamespace || isRequire {
			var arg string = this_.ParseJsString(this_.googMatcher.Group_1_string("args"))
			_ = arg
			if isRequire {
				if ("requireType" == methodName) || ("forwardDeclare" == methodName) {
					this_.typeRequires.Add_1(j2gort.BoxString(arg))
				} else if !("goog" == arg) {
					var require *DependencyInfo_Require = DependencyInfo_Require_GoogRequireSymbol(arg)
					_ = require
					this_.requires.Add_1(j2gort.PtrToAny(require))
				}
			} else {
				this_.provides.Add_1(j2gort.BoxString(arg))
			}
		}
	}
	return lineHasProvidesOrRequires
}

// Java method: JsFileRegexParser.parseLine(java.lang.String line)
func (this_ *JsFileRegexParser) ParseLine(line string) bool {
	_ = line
	var lineHasProvidesOrRequires bool = false
	_ = lineHasProvidesOrRequires
	if javalang.StringStartsWith(j2gort.CheckNotNullString(line), "goog.loadModule(function(") {
		this_.seenLoadModule = true
	}
	var lineHasProvidesOrRequiresWords bool = ((((((j2gort.JStrIndexOf(line, "provide") >= 0) || (j2gort.JStrIndexOf(line, "require") >= 0)) || (j2gort.JStrIndexOf(line, "forwardDeclare") >= 0)) || (j2gort.JStrIndexOf(line, "module") >= 0)) || (j2gort.JStrIndexOf(line, "addDependency") >= 0)) || (j2gort.JStrIndexOf(line, "declareModuleId") >= 0))
	_ = lineHasProvidesOrRequiresWords
	if !(j2gort.CheckNotNullString(this_.googMatcherBuffer) == "") {
		lineHasProvidesOrRequires = this_.priv_JsFileRegexParser_applyGoogMatcher_6102((j2gort.StrValString(this_.googMatcherBuffer) + j2gort.StrValString(line)))
	}
	if !(lineHasProvidesOrRequires) && lineHasProvidesOrRequiresWords {
		lineHasProvidesOrRequires = this_.priv_JsFileRegexParser_applyGoogMatcher_6102(line)
	}
	if !(lineHasProvidesOrRequires) && lineHasProvidesOrRequiresWords {
		this_.googMatcherBuffer = line
	} else {
		this_.googMatcherBuffer = ""
	}
	if javalang.StringStartsWith(j2gort.CheckNotNullString(line), "import") || javalang.StringStartsWith(j2gort.CheckNotNullString(line), "export") {
		this_.es6Matcher.Reset_1(j2gort.BoxString(line))
		for this_.es6Matcher.Find_0() {
			this_.priv_JsFileRegexParser_setModuleType_6101(StaticGet_JsFileRegexParser_ModuleType_ES6_MODULE())
			lineHasProvidesOrRequires = true
			var arg string = this_.es6Matcher.Group_1_int32(int32(1))
			_ = arg
			if arg != j2gort.NullString {
				if javalang.StringStartsWith(j2gort.CheckNotNullString(arg), "goog:") {
					this_.requires.Add_1(j2gort.PtrToAny(DependencyInfo_Require_GoogRequireSymbol(j2gort.JStrSubFrom(arg, int32(5)))))
				} else {
					var path *ModuleLoader_ModulePath = this_.file.ResolveJsModule(arg, this_.GetJsFileLineParserBase().filePath, this_.GetJsFileLineParserBase().lineNum, this_.es6Matcher.Start_0())
					_ = path
					if path == nil {
						path = this_.file.ResolveModuleAsPath(arg)
					}
					this_.requires.Add_1(j2gort.PtrToAny(DependencyInfo_Require_Es6Import(path.ToModuleName(), arg)))
				}
			}
		}
		if (this_.moduleType != StaticGet_JsFileRegexParser_ModuleType_ES6_MODULE()) && StaticGet_JsFileRegexParser_ES6_EXPORT_PATTERN().Matcher(j2gort.BoxString(line)).LookingAt() {
			this_.priv_JsFileRegexParser_setModuleType_6101(StaticGet_JsFileRegexParser_ModuleType_ES6_MODULE())
		}
	}
	return ((((((!(this_.GetJsFileLineParserBase().shortcutMode) || lineHasProvidesOrRequires) || !(j2gort.CheckNotNullString(this_.googMatcherBuffer) == "")) || j2gort.Ext_com_google_common_base_CharMatcher_S_Whitespace().MatchesAllOf(j2gort.BoxString(line))) || !(j2gort.JStrIndexOf(line, ";") >= 0)) || (j2gort.JStrIndexOf(line, "goog.setTestOnly") >= 0)) || (j2gort.JStrIndexOf(line, "goog.module.declareLegacyNamespace") >= 0))
}

// Java static method: JsFileRegexParser.isSupported()
func JsFileRegexParser_IsSupported() bool {
	clinit_JsFileRegexParser()
	return true
}

// Java enum: com.google.javascript.jscomp.deps.JsFileRegexParser.ModuleType
type JsFileRegexParser_ModuleType uint8

const (
	enum_JsFileRegexParser_ModuleType_NON_MODULE   JsFileRegexParser_ModuleType = 1
	enum_JsFileRegexParser_ModuleType_GOOG_MODULE  JsFileRegexParser_ModuleType = 2
	enum_JsFileRegexParser_ModuleType_GOOG_PROVIDE JsFileRegexParser_ModuleType = 3
	enum_JsFileRegexParser_ModuleType_ES6_MODULE   JsFileRegexParser_ModuleType = 4
)

var enumBaseTable_JsFileRegexParser_ModuleType = [...]j2gort.EnumBase{
	{Name_: "NON_MODULE", Ordinal_: 0},
	{Name_: "GOOG_MODULE", Ordinal_: 1},
	{Name_: "GOOG_PROVIDE", Ordinal_: 2},
	{Name_: "ES6_MODULE", Ordinal_: 3},
}
var enumValues_JsFileRegexParser_ModuleType = []JsFileRegexParser_ModuleType{1, 2, 3, 4}

func StaticGet_JsFileRegexParser_ModuleType_NON_MODULE() JsFileRegexParser_ModuleType {
	return enum_JsFileRegexParser_ModuleType_NON_MODULE
}
func StaticGet_JsFileRegexParser_ModuleType_GOOG_MODULE() JsFileRegexParser_ModuleType {
	return enum_JsFileRegexParser_ModuleType_GOOG_MODULE
}
func StaticGet_JsFileRegexParser_ModuleType_GOOG_PROVIDE() JsFileRegexParser_ModuleType {
	return enum_JsFileRegexParser_ModuleType_GOOG_PROVIDE
}
func StaticGet_JsFileRegexParser_ModuleType_ES6_MODULE() JsFileRegexParser_ModuleType {
	return enum_JsFileRegexParser_ModuleType_ES6_MODULE
}

func JsFileRegexParser_ModuleType_Values() []JsFileRegexParser_ModuleType {
	return j2gort.CloneSlice(enumValues_JsFileRegexParser_ModuleType)
}

func JsFileRegexParser_ModuleType_ValueOf(name string) JsFileRegexParser_ModuleType {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_JsFileRegexParser_ModuleType {
		if enumBaseTable_JsFileRegexParser_ModuleType[i_].Name_ == name {
			return JsFileRegexParser_ModuleType(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.deps.JsFileRegexParser.ModuleType." + name))
}

func (this_ JsFileRegexParser_ModuleType) IsJavaNil_() bool                              { return this_ == 0 }
func (JsFileRegexParser_ModuleType) HasIdentityEquals_()                                 {}
func (JsFileRegexParser_ModuleType) Implements_j2gort_Ext_java_lang_Enum()               {}
func (JsFileRegexParser_ModuleType) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ JsFileRegexParser_ModuleType) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (JsFileRegexParser_ModuleType) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.deps.JsFileRegexParser$ModuleType", "ModuleType"
}
func (this_ JsFileRegexParser_ModuleType) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_JsFileRegexParser_ModuleType[this_-1]
}
func (this_ JsFileRegexParser_ModuleType) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ JsFileRegexParser_ModuleType) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ JsFileRegexParser_ModuleType) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_JsFileRegexParser_ModuleType[this_-1].Name_
}
func (this_ JsFileRegexParser_ModuleType) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_JsFileRegexParser_ModuleType[this_-1].Name_
}
func (this_ JsFileRegexParser_ModuleType) Equals(other_ any) bool {
	o_, ok_ := other_.(JsFileRegexParser_ModuleType)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ JsFileRegexParser_ModuleType) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_JsFileRegexParser_ModuleType[this_-1].HashCode()
}
func (this_ JsFileRegexParser_ModuleType) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_JsFileRegexParser_ModuleType[this_-1].CompareTo(other_)
}
func (this_ JsFileRegexParser_ModuleType) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_JsFileRegexParser_ModuleType[this_-1].Name_)
}
func init() {
	vals_ := []any{JsFileRegexParser_ModuleType(1), JsFileRegexParser_ModuleType(2), JsFileRegexParser_ModuleType(3), JsFileRegexParser_ModuleType(4)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.deps.JsFileRegexParser.ModuleType", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.deps.JsFileRegexParser$ModuleType", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.deps.JsFileRegexParser.ModuleType", "com.google.javascript.jscomp.deps.JsFileRegexParser$ModuleType", func(v_ any) bool {
		return j2gort.InstanceOf[JsFileRegexParser_ModuleType](v_)
	})
}
