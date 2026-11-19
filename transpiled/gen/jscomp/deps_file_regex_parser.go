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
// Port of com/google/javascript/jscomp/deps/DepsFileRegexParser.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

var lambdaSingleton_deps_file_regex_parser_1 j2gort.Ext_java_util_function_Function

func init() {
	lambdaSingleton_deps_file_regex_parser_1 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return j2gort.PtrToAny(DependencyInfo_Require_ParsedFromDeps(j2gort.CastString(mrp0_)))
	})
}

// Java class: com.google.javascript.jscomp.deps.DepsFileRegexParser
type DepsFileRegexParser struct {
	JsFileLineParserBase
	depMatcher     *j2gort.Ext_java_util_regex_Matcher        `gson:"depMatcher"`
	depArgsMatch   *j2gort.Ext_java_util_regex_Matcher        `gson:"depArgsMatch"`
	depInfos       j2gort.Ext_java_util_List                  `gson:"depInfos"`
	pathTranslator j2gort.Ext_com_google_common_base_Function `gson:"pathTranslator"`
}

func (this_ *DepsFileRegexParser) BindSelf_(s any) {
	_ = s
	this_.JsFileLineParserBase.BindSelf_(s)
}
func (this_ *DepsFileRegexParser) IsJavaNil_() bool { return this_ == nil }
func (*DepsFileRegexParser) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.deps.DepsFileRegexParser", "DepsFileRegexParser"
}
func (this_ *DepsFileRegexParser) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.JsFileLineParserBase.IdentityHash_()
}
func (*DepsFileRegexParser) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.deps.DepsFileRegexParser", "com.google.javascript.jscomp.deps.DepsFileRegexParser", func(v_ any) bool {
		return j2gort.InstanceOf[*DepsFileRegexParser](v_)
	})
}
func (*DepsFileRegexParser) Implements_jscomp_JsFileLineParser() {}

var static_DepsFileRegexParser_logger *j2gort.Ext_java_util_logging_Logger

func StaticGet_DepsFileRegexParser_Logger() *j2gort.Ext_java_util_logging_Logger {
	clinit_DepsFileRegexParser()
	return static_DepsFileRegexParser_logger
}
func StaticSet_DepsFileRegexParser_Logger(v *j2gort.Ext_java_util_logging_Logger) {
	clinit_DepsFileRegexParser()
	static_DepsFileRegexParser_logger = v
}

var clinitDone_DepsFileRegexParser uint32
var clinitState_DepsFileRegexParser j2gort.ClinitState

func clinit_DepsFileRegexParser() {
	if atomic.LoadUint32(&clinitDone_DepsFileRegexParser) != 0 {
		return
	}
	clinitSlow_DepsFileRegexParser()
}

func clinitSlow_DepsFileRegexParser() {
	if !j2gort.EnterClinit(&clinitState_DepsFileRegexParser, &clinitDone_DepsFileRegexParser) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_DepsFileRegexParser, &clinitDone_DepsFileRegexParser)
	static_DepsFileRegexParser_logger = j2gort.Ext_java_util_logging_Logger_GetLogger_1(j2gort.CheckNotNullPtr(j2gort.ClassLiteral("com.google.javascript.jscomp.deps.DepsFileRegexParser", "DepsFileRegexParser")).GetName())
}

func (this_ *DepsFileRegexParser) initFields_DepsFileRegexParser() {
	_ = this_
	this_.depMatcher = j2gort.Ext_java_util_regex_Pattern_Compile_1("\\s*goog.addDependency\\((.*)\\);?\\s*").Matcher("")
	this_.depArgsMatch = j2gort.Ext_java_util_regex_Pattern_Compile_1("\\s*([^,]*), (\\[[^\\]]*\\]), (\\[[^\\]]*\\])(?:, (true|false|\\{[^{}]*\\}))?\\s*").Matcher("")
}

// Java constructor: DepsFileRegexParser(com.google.javascript.jscomp.ErrorManager errorManager)
func NewDepsFileRegexParser_1(errorManager ErrorManager) *DepsFileRegexParser {
	clinit_DepsFileRegexParser()
	this_ := &DepsFileRegexParser{}
	this_.BindSelf_(this_)
	InitDepsFileRegexParser_1(this_, errorManager)
	return this_
}

func InitDepsFileRegexParser_1(this_ *DepsFileRegexParser, errorManager ErrorManager) {
	_ = errorManager
	InitDepsFileRegexParser_2(this_, j2gort.Ext_com_google_common_base_Functions_Identity(), errorManager)
}

// Java constructor: DepsFileRegexParser(com.google.common.base.Function<java.lang.String,java.lang.String> pathTranslator, com.google.javascript.jscomp.ErrorManager errorManager)
func NewDepsFileRegexParser_2(pathTranslator j2gort.Ext_com_google_common_base_Function, errorManager ErrorManager) *DepsFileRegexParser {
	clinit_DepsFileRegexParser()
	this_ := &DepsFileRegexParser{}
	this_.BindSelf_(this_)
	InitDepsFileRegexParser_2(this_, pathTranslator, errorManager)
	return this_
}

func InitDepsFileRegexParser_2(this_ *DepsFileRegexParser, pathTranslator j2gort.Ext_com_google_common_base_Function, errorManager ErrorManager) {
	_ = pathTranslator
	_ = errorManager
	InitJsFileLineParser(&this_.JsFileLineParserBase, errorManager)
	this_.initFields_DepsFileRegexParser()
	this_.pathTranslator = pathTranslator
}

// Java method: DepsFileRegexParser.parseFile(java.lang.String filePath)
func (this_ *DepsFileRegexParser) ParseFile_1(filePath string) j2gort.Ext_java_util_List {
	_ = filePath
	return this_.ParseFileReader(filePath, j2gort.AsIface[j2gort.Ext_java_io_Reader](j2gort.Ext_com_google_common_io_Files_NewReader(j2gort.NewExt_java_io_File_1_string(filePath), j2gort.Ext_java_nio_charset_StandardCharsets_F_UTF_8)))
}

// Java method: DepsFileRegexParser.parseFile(java.lang.String filePath, java.lang.String fileContents)
func (this_ *DepsFileRegexParser) ParseFile_2(filePath string, fileContents string) j2gort.Ext_java_util_List {
	_ = filePath
	_ = fileContents
	return this_.ParseFileReader(filePath, j2gort.NewExt_java_io_StringReader(fileContents))
}

// Java method: DepsFileRegexParser.parseFileReader(java.lang.String filePath, java.io.Reader reader)
func (this_ *DepsFileRegexParser) ParseFileReader(filePath string, reader j2gort.Ext_java_io_Reader) j2gort.Ext_java_util_List {
	_ = filePath
	_ = reader
	this_.depInfos = j2gort.NewExt_java_util_ArrayList()
	if StaticGet_DepsFileRegexParser_Logger().IsLoggable(j2gort.Ext_java_util_logging_Level_F_FINE) {
		StaticGet_DepsFileRegexParser_Logger().Fine_1_string(("Parsing Dep: " + j2gort.StrValString(filePath)))
	}
	this_.DoParse(filePath, reader)
	return this_.depInfos
}

// Java method: DepsFileRegexParser.parseLine(java.lang.String line)
func (this_ *DepsFileRegexParser) ParseLine(line string) bool {
	_ = line
	var hasDependencies bool = false
	_ = hasDependencies
	if j2gort.JStrIndexOf(line, "addDependency") >= 0 {
		this_.depMatcher.Reset_1(j2gort.BoxString(line))
		if this_.depMatcher.Matches() {
			hasDependencies = true
			var addDependencyParams string = this_.depMatcher.Group_1_int32(int32(1))
			_ = addDependencyParams
			this_.depArgsMatch.Reset_1(j2gort.BoxString(addDependencyParams))
			if !(this_.depArgsMatch.Matches()) {
				panic(NewJsFileLineParser_ParseException(("Invalid arguments to goog.addDependency(). Found: " + j2gort.StrValString(addDependencyParams)), true))
			}
			var relativePath string = this_.ParseJsString(this_.depArgsMatch.Group_1_int32(int32(1)))
			_ = relativePath
			var path string = j2gort.CastString(this_.pathTranslator.Apply(j2gort.BoxString(relativePath)))
			_ = path
			var provides j2gort.Ext_java_util_List = this_.ParseJsStringArray(this_.depArgsMatch.Group_1_int32(int32(2)))
			_ = provides
			var loadFlags j2gort.Ext_java_util_Map = this_.priv_DepsFileRegexParser_parseLoadFlags_6083(this_.depArgsMatch.Group_1_int32(int32(4)))
			_ = loadFlags
			{
				var iifeCond_1 bool
				{

					sr_1 := "es6"
					se_2 := j2gort.CastString(j2gort.MapGetStr(loadFlags, "module"))
					iifeCond_1 = j2gort.CheckNotNullString(sr_1) == se_2
				}
				if iifeCond_1 {
					provides.Add_1(j2gort.BoxString(relativePath))
				}
			}
			var depInfo DependencyInfo = Upcast_DependencyInfo(SimpleDependencyInfo_S_Builder(path, this_.GetJsFileLineParserBase().filePath).SetProvides_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](provides)).SetRequires_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableList](this_.ParseJsStringArray(this_.depArgsMatch.Group_1_int32(int32(3))).Stream().Map(lambdaSingleton_deps_file_regex_parser_1).Collect_1(j2gort.Ext_com_google_common_collect_ImmutableList_ToImmutableList()))).SetLoadFlags(loadFlags).Build())
			_ = depInfo
			if StaticGet_DepsFileRegexParser_Logger().IsLoggable(j2gort.Ext_java_util_logging_Level_F_FINE) {
				StaticGet_DepsFileRegexParser_Logger().Fine_1_string(("Found dep: " + j2gort.StrVal(depInfo)))
			}
			this_.depInfos.Add_1(depInfo)
		}
	}
	return ((!(this_.GetJsFileLineParserBase().shortcutMode) || hasDependencies) || j2gort.Ext_com_google_common_base_CharMatcher_S_Whitespace().MatchesAllOf(j2gort.BoxString(line)))
}

// Java method: DepsFileRegexParser.parseLoadFlags(java.lang.String loadFlags)
func (this_ *DepsFileRegexParser) priv_DepsFileRegexParser_parseLoadFlags_6083(loadFlags string) j2gort.Ext_java_util_Map {
	_ = loadFlags
	if (loadFlags == j2gort.NullString) || (j2gort.CheckNotNullString(loadFlags) == "false") {
		return j2gort.AsIface[j2gort.Ext_java_util_Map](j2gort.Ext_com_google_common_collect_ImmutableMap_Of_0())
	} else if j2gort.CheckNotNullString(loadFlags) == "true" {
		return j2gort.AsIface[j2gort.Ext_java_util_Map](j2gort.Ext_com_google_common_collect_ImmutableMap_Of_2("module", "goog"))
	} else {
		return this_.ParseJsStringMap(loadFlags)
	}
}
