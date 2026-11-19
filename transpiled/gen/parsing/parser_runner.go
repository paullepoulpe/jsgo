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
// Port of com/google/javascript/jscomp/parsing/ParserRunner.java (generated mechanically by j2go).

package parsing

import (
	"sync"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.parsing.ParserRunner
type ParserRunner struct {
	idHash_ atomic.Int32
}

func (this_ *ParserRunner) BindSelf_(s any) {
	_ = s
}
func (this_ *ParserRunner) IsJavaNil_() bool { return this_ == nil }
func (*ParserRunner) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.ParserRunner", "ParserRunner"
}
func (this_ *ParserRunner) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ParserRunner) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.ParserRunner", "com.google.javascript.jscomp.parsing.ParserRunner", func(v_ any) bool {
		return j2gort.InstanceOf[*ParserRunner](v_)
	})
}

var staticSyncMu_ParserRunner sync.Mutex
var static_ParserRunner_annotationNames j2gort.Ext_com_google_common_collect_ImmutableSet

func StaticGet_ParserRunner_AnnotationNames() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_ParserRunner()
	return static_ParserRunner_annotationNames
}
func StaticSet_ParserRunner_AnnotationNames(v j2gort.Ext_com_google_common_collect_ImmutableSet) {
	clinit_ParserRunner()
	static_ParserRunner_annotationNames = v
}

var static_ParserRunner_suppressionNames j2gort.Ext_com_google_common_collect_ImmutableSet

func StaticGet_ParserRunner_SuppressionNames() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_ParserRunner()
	return static_ParserRunner_suppressionNames
}
func StaticSet_ParserRunner_SuppressionNames(v j2gort.Ext_com_google_common_collect_ImmutableSet) {
	clinit_ParserRunner()
	static_ParserRunner_suppressionNames = v
}

var static_ParserRunner_reservedVars j2gort.Ext_com_google_common_collect_ImmutableSet

func StaticGet_ParserRunner_ReservedVars() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_ParserRunner()
	return static_ParserRunner_reservedVars
}
func StaticSet_ParserRunner_ReservedVars(v j2gort.Ext_com_google_common_collect_ImmutableSet) {
	clinit_ParserRunner()
	static_ParserRunner_reservedVars = v
}

var static_ParserRunner_closurePrimitiveNames j2gort.Ext_com_google_common_collect_ImmutableSet

func StaticGet_ParserRunner_ClosurePrimitiveNames() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_ParserRunner()
	return static_ParserRunner_closurePrimitiveNames
}
func StaticSet_ParserRunner_ClosurePrimitiveNames(v j2gort.Ext_com_google_common_collect_ImmutableSet) {
	clinit_ParserRunner()
	static_ParserRunner_closurePrimitiveNames = v
}

var clinitDone_ParserRunner uint32
var clinitState_ParserRunner j2gort.ClinitState

func clinit_ParserRunner() {
	if atomic.LoadUint32(&clinitDone_ParserRunner) != 0 {
		return
	}
	clinitSlow_ParserRunner()
}

func clinitSlow_ParserRunner() {
	if !j2gort.EnterClinit(&clinitState_ParserRunner, &clinitDone_ParserRunner) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ParserRunner, &clinitDone_ParserRunner)
	static_ParserRunner_annotationNames = nil
	static_ParserRunner_suppressionNames = nil
	static_ParserRunner_reservedVars = nil
	static_ParserRunner_closurePrimitiveNames = nil
}

func (this_ *ParserRunner) initFields_ParserRunner() {
	_ = this_
}

// Java constructor: ParserRunner()
func NewParserRunner() *ParserRunner {
	clinit_ParserRunner()
	this_ := &ParserRunner{}
	InitParserRunner(this_)
	return this_
}

func InitParserRunner(this_ *ParserRunner) {
	this_.initFields_ParserRunner()
}

// Java static method: ParserRunner.createConfig(com.google.javascript.jscomp.parsing.Config.LanguageMode languageMode, java.util.Set<java.lang.String> extraAnnotationNames, com.google.javascript.jscomp.parsing.Config.StrictMode strictMode)
func ParserRunner_CreateConfig_3(languageMode *Config_LanguageMode, extraAnnotationNames j2gort.Ext_java_util_Set, strictMode Config_StrictMode) Config {
	clinit_ParserRunner()
	_ = languageMode
	_ = extraAnnotationNames
	_ = strictMode
	return ParserRunner_CreateConfig_6(languageMode, StaticGet_Config_JsDocParsing_TYPES_ONLY(), StaticGet_Config_RunMode_STOP_AFTER_ERROR(), extraAnnotationNames, true, strictMode)
}

// Java static method: ParserRunner.createConfig(com.google.javascript.jscomp.parsing.Config.LanguageMode languageMode, com.google.javascript.jscomp.parsing.Config.JsDocParsing jsdocParsingMode, com.google.javascript.jscomp.parsing.Config.RunMode runMode, java.util.Set<java.lang.String> extraAnnotationNames, boolean parseInlineSourceMaps, com.google.javascript.jscomp.parsing.Config.StrictMode strictMode)
func ParserRunner_CreateConfig_6(languageMode *Config_LanguageMode, jsdocParsingMode Config_JsDocParsing, runMode Config_RunMode, extraAnnotationNames j2gort.Ext_java_util_Set, parseInlineSourceMaps bool, strictMode Config_StrictMode) Config {
	clinit_ParserRunner()
	_ = languageMode
	_ = jsdocParsingMode
	_ = runMode
	_ = extraAnnotationNames
	_ = parseInlineSourceMaps
	_ = strictMode
	ParserRunner_InitResourceConfig()
	var effectiveAnnotationNames j2gort.Ext_java_util_Set = nil
	_ = effectiveAnnotationNames
	if extraAnnotationNames == nil {
		effectiveAnnotationNames = j2gort.AsIface[j2gort.Ext_java_util_Set](StaticGet_ParserRunner_AnnotationNames())
	} else {
		effectiveAnnotationNames = j2gort.NewExt_java_util_LinkedHashSet_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](StaticGet_ParserRunner_AnnotationNames()))
		effectiveAnnotationNames.AddAll_1(j2gort.AsIface[j2gort.Ext_java_util_Collection](extraAnnotationNames))
	}
	return Config_S_Builder().GetConfig_BuilderBase().SetExtraAnnotationNames(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](effectiveAnnotationNames)).SetJsDocParsingMode(jsdocParsingMode).SetRunMode(runMode).SetSuppressionNames(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](StaticGet_ParserRunner_SuppressionNames())).SetClosurePrimitiveNames(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](StaticGet_ParserRunner_ClosurePrimitiveNames())).SetLanguageMode(languageMode).SetParseInlineSourceMaps(parseInlineSourceMaps).SetStrictMode(strictMode).Build()
}

// Java static method: ParserRunner.getReservedVars()
func ParserRunner_GetReservedVars() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_ParserRunner()
	ParserRunner_InitResourceConfig()
	return StaticGet_ParserRunner_ReservedVars()
}

// Java static method: ParserRunner.getSuppressionNames()
func ParserRunner_GetSuppressionNames() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_ParserRunner()
	ParserRunner_InitResourceConfig()
	return StaticGet_ParserRunner_SuppressionNames()
}

// Java static method: ParserRunner.initResourceConfig()
func ParserRunner_InitResourceConfig() {
	clinit_ParserRunner()
	staticSyncMu_ParserRunner.Lock()
	defer staticSyncMu_ParserRunner.Unlock()
	if StaticGet_ParserRunner_AnnotationNames() != nil {
		return
	}
	StaticSet_ParserRunner_AnnotationNames(ParserRunner_ExtractList(ParserConfiguration_GetString("jsdoc.annotations")))
	StaticSet_ParserRunner_SuppressionNames(ParserRunner_ExtractList(ParserConfiguration_GetString("jsdoc.suppressions")))
	StaticSet_ParserRunner_ClosurePrimitiveNames(ParserRunner_ExtractList(ParserConfiguration_GetString("jsdoc.primitives")))
	StaticSet_ParserRunner_ReservedVars(ParserRunner_ExtractList(ParserConfiguration_GetString("compiler.reserved.vars")))
}

// Java static method: ParserRunner.extractList(java.lang.String configProp)
func ParserRunner_ExtractList(configProp string) j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_ParserRunner()
	_ = configProp
	return j2gort.Ext_com_google_common_collect_ImmutableSet_CopyOf_1_Ext_java_lang_Iterable(j2gort.SplitterSplitStr(j2gort.Ext_com_google_common_base_Splitter_On_1_uint16(uint16(44)).TrimResults_0(), configProp))
}

// Java static method: ParserRunner.parse(com.google.javascript.rhino.StaticSourceFile sourceFile, java.lang.String sourceString, com.google.javascript.jscomp.parsing.Config config, com.google.javascript.rhino.ErrorReporter errorReporter)
func ParserRunner_Parse(sourceFile pkg_rhino.StaticSourceFile, sourceString string, config Config, errorReporter pkg_rhino.ErrorReporter) *ParserRunner_ParseResult {
	clinit_ParserRunner()
	_ = sourceFile
	_ = sourceString
	_ = config
	_ = errorReporter
	var sourceName string = sourceFile.GetName()
	_ = sourceName
	var didRet_1 bool
	var retVal_2 *ParserRunner_ParseResult
	func() {
		defer func() {
			if rec_3 := recover(); rec_3 != nil {
				jex_4 := javalang.RecoverException(rec_3)
				_ = jex_4
				if j2gort.InstanceOf[javalang.Throwable](jex_4) {
					t := j2gort.AsIface[javalang.Throwable](jex_4)
					_ = t
					panic(j2gort.AsIface[javalang.IRuntimeException](j2gort.ConstructShimThrowable("java.lang.RuntimeException", "Exception parsing \""+j2gort.StrValString(sourceName)+"\"", t)))
				} else {
					panic(rec_3)
				}
			}
		}()
		var file *pkg_parser.SourceFile = pkg_parser.NewSourceFile(sourceName, sourceString)
		_ = file
		var keepGoing bool = (config.RunMode() == StaticGet_Config_RunMode_KEEP_GOING())
		_ = keepGoing
		var es6ErrorReporter *ParserRunner_Es6ErrorReporter = NewParserRunner_Es6ErrorReporter(errorReporter, file.FieldGet_name(), keepGoing)
		_ = es6ErrorReporter
		var es6config *pkg_parser.Parser_Config = ParserRunner_NewParserConfig(config)
		_ = es6config
		var p *pkg_parser.Parser = acquirePooledParser_(es6config, es6ErrorReporter, file)
		_ = p
		defer releasePooledParser_(p)
		var tree *pkg_parser.ProgramTree = p.ParseProgram()
		_ = tree
		var root pkg_rhino.Node = pkg_rhino.Node(0)
		_ = root
		var comments j2gort.Ext_java_util_List = j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.Ext_com_google_common_collect_ImmutableList_Of_0())
		_ = comments
		var topLevelStatementRanges j2gort.Ext_com_google_common_collect_ImmutableList = j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()
		_ = topLevelStatementRanges
		var features *pkg_parser.FeatureSet = p.GetFeatures()
		_ = features
		if (tree != nil) && (!(es6ErrorReporter.HadError()) || keepGoing) {
			var factory *IRFactory = IRFactory_TransformTree(tree, sourceFile, config, errorReporter, file)
			_ = factory
			root = factory.GetResultNode()
			features = features.Union(factory.GetFeatures())
			root.SetFeatureSet(j2gort.PtrToAny(features))
			if config.JsDocParsingMode().ShouldParseDescriptions() {
				comments = p.GetComments()
				var builder *j2gort.Ext_com_google_common_collect_ImmutableList_Builder = j2gort.InitImmListBuilderCap(&j2gort.Ext_com_google_common_collect_ImmutableList_Builder{}, tree.FieldGet_sourceElements().Size())
				_ = builder
				{
					cur_5 := j2gort.NewIterCursor(j2gort.CheckNotNull(tree.FieldGet_sourceElements()))
					for cur_5.HasNext() {
						rawElem_6 := cur_5.Next()
						element := j2gort.AsIface[pkg_parser.ParseTree](rawElem_6)
						_ = element
						builder.AddVoid_(element.GetParseTreeBase().FieldGet_location())
					}
				}
				topLevelStatementRanges = builder.BuildImmList_()
			}
		}
		retVal_2 = NewParserRunner_ParseResult_5(root, comments, topLevelStatementRanges, features, p.GetSourceMapURL())
		didRet_1 = true
		return
	}()
	if didRet_1 {
		return retVal_2
	}
	panic("unreachable")
}

// Java static method: ParserRunner.newParserConfig(com.google.javascript.jscomp.parsing.Config config)
func ParserRunner_NewParserConfig(config Config) *pkg_parser.Parser_Config {
	clinit_ParserRunner()
	_ = config
	var languageMode *Config_LanguageMode = config.LanguageMode()
	_ = languageMode
	var isStrictMode bool = config.StrictMode().IsStrict()
	_ = isStrictMode
	var parserConfigLanguageMode pkg_parser.Parser_Config_Mode = pkg_parser.Parser_Config_Mode(0)
	_ = parserConfigLanguageMode
	switch j2gort.CheckNotNullPtr(languageMode).Ordinal() {
	case int32(0) /* LanguageMode.ECMASCRIPT3 */ :
		parserConfigLanguageMode = pkg_parser.StaticGet_Parser_Config_Mode_ES3()
	case int32(1) /* LanguageMode.ECMASCRIPT5 */ :
		parserConfigLanguageMode = pkg_parser.StaticGet_Parser_Config_Mode_ES5()
	case int32(2) /* LanguageMode.ECMASCRIPT_2015 */, int32(3) /* LanguageMode.ECMASCRIPT_2016 */ :
		parserConfigLanguageMode = pkg_parser.StaticGet_Parser_Config_Mode_ES6_OR_ES7()
	case int32(4) /* LanguageMode.ECMASCRIPT_2017 */, int32(5) /* LanguageMode.ECMASCRIPT_2018 */, int32(6) /* LanguageMode.ECMASCRIPT_2019 */, int32(7) /* LanguageMode.ECMASCRIPT_2020 */, int32(8) /* LanguageMode.ECMASCRIPT_2021 */, int32(9) /* LanguageMode.ECMASCRIPT_2022 */, int32(10) /* LanguageMode.ES_NEXT */, int32(11) /* LanguageMode.UNSTABLE */, int32(12) /* LanguageMode.UNSUPPORTED */ :
		parserConfigLanguageMode = pkg_parser.StaticGet_Parser_Config_Mode_ES8_OR_GREATER()
	}
	return pkg_parser.NewParser_Config_2(j2gort.CheckNotNullEnum(parserConfigLanguageMode), isStrictMode)
}

// Java class: com.google.javascript.jscomp.parsing.ParserRunner.Es6ErrorReporter
type ParserRunner_Es6ErrorReporter struct {
	pkg_parser.ErrorReporterBase
	reporter        pkg_rhino.ErrorReporter `gson:"reporter"`
	sourceName      string                  `gson:"sourceName"`
	errorSeen       bool                    `gson:"errorSeen"`
	reportAllErrors bool                    `gson:"reportAllErrors"`
}

func (this_ *ParserRunner_Es6ErrorReporter) BindSelf_(s any) {
	_ = s
	this_.ErrorReporterBase.BindSelf_(s)
}
func (this_ *ParserRunner_Es6ErrorReporter) IsJavaNil_() bool { return this_ == nil }
func (*ParserRunner_Es6ErrorReporter) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.ParserRunner$Es6ErrorReporter", "Es6ErrorReporter"
}
func (this_ *ParserRunner_Es6ErrorReporter) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.ErrorReporterBase.IdentityHash_()
}
func (*ParserRunner_Es6ErrorReporter) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.ParserRunner.Es6ErrorReporter", "com.google.javascript.jscomp.parsing.ParserRunner$Es6ErrorReporter", func(v_ any) bool {
		return j2gort.InstanceOf[*ParserRunner_Es6ErrorReporter](v_)
	})
}
func (*ParserRunner_Es6ErrorReporter) Implements_parser_ErrorReporter() {}

func (this_ *ParserRunner_Es6ErrorReporter) initFields_ParserRunner_Es6ErrorReporter() {
	_ = this_
	this_.sourceName = j2gort.NullString

}

// Java constructor: Es6ErrorReporter(com.google.javascript.rhino.ErrorReporter reporter, java.lang.String sourceName, boolean reportAllErrors)
func NewParserRunner_Es6ErrorReporter(reporter pkg_rhino.ErrorReporter, sourceName string, reportAllErrors bool) *ParserRunner_Es6ErrorReporter {
	this_ := &ParserRunner_Es6ErrorReporter{}
	this_.BindSelf_(this_)
	InitParserRunner_Es6ErrorReporter(this_, reporter, sourceName, reportAllErrors)
	return this_
}

func InitParserRunner_Es6ErrorReporter(this_ *ParserRunner_Es6ErrorReporter, reporter pkg_rhino.ErrorReporter, sourceName string, reportAllErrors bool) {
	_ = reporter
	_ = sourceName
	_ = reportAllErrors
	pkg_parser.InitErrorReporter(&this_.ErrorReporterBase)
	this_.initFields_ParserRunner_Es6ErrorReporter()
	this_.reporter = reporter
	this_.sourceName = sourceName
	this_.reportAllErrors = reportAllErrors
}

// Java method: Es6ErrorReporter.reportError(com.google.javascript.jscomp.parsing.parser.util.SourcePosition location, java.lang.String message)
func (this_ *ParserRunner_Es6ErrorReporter) ReportError_2(location pkg_parser.SourcePosition, message string) {
	_ = location
	_ = message
	if this_.reportAllErrors || !(this_.errorSeen) {
		this_.errorSeen = true
		this_.reporter.M_Error(message, this_.sourceName, (location.FieldGet_line() + int32(1)), location.FieldGet_column())
	}
}

// Java method: Es6ErrorReporter.reportWarning(com.google.javascript.jscomp.parsing.parser.util.SourcePosition location, java.lang.String message)
func (this_ *ParserRunner_Es6ErrorReporter) ReportWarning_2(location pkg_parser.SourcePosition, message string) {
	_ = location
	_ = message
	this_.reporter.Warning(message, this_.sourceName, (location.FieldGet_line() + int32(1)), location.FieldGet_column())
}

// Java class: com.google.javascript.jscomp.parsing.ParserRunner.ParseResult
type ParserRunner_ParseResult struct {
	comments                j2gort.Ext_java_util_List                          `gson:"comments"`
	topLevelStatementRanges j2gort.Ext_com_google_common_collect_ImmutableList `gson:"topLevelStatementRanges"`
	features                *pkg_parser.FeatureSet                             `gson:"features"`
	sourceMapURL            string                                             `gson:"sourceMapURL"`
	ast                     pkg_rhino.Node                                     `gson:"ast"`
	idHash_                 atomic.Int32
}

func (this_ *ParserRunner_ParseResult) BindSelf_(s any) {
	_ = s
}
func (this_ *ParserRunner_ParseResult) IsJavaNil_() bool { return this_ == nil }
func (*ParserRunner_ParseResult) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.ParserRunner$ParseResult", "ParseResult"
}
func (this_ *ParserRunner_ParseResult) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ParserRunner_ParseResult) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.ParserRunner.ParseResult", "com.google.javascript.jscomp.parsing.ParserRunner$ParseResult", func(v_ any) bool {
		return j2gort.InstanceOf[*ParserRunner_ParseResult](v_)
	})
}
func (this_ *ParserRunner_ParseResult) FieldGet_ast() pkg_rhino.Node { return this_.ast }
func (this_ *ParserRunner_ParseResult) FieldGet_comments() j2gort.Ext_java_util_List {
	return this_.comments
}
func (this_ *ParserRunner_ParseResult) FieldGet_features() *pkg_parser.FeatureSet {
	return this_.features
}
func (this_ *ParserRunner_ParseResult) FieldGet_sourceMapURL() string { return this_.sourceMapURL }

func (this_ *ParserRunner_ParseResult) initFields_ParserRunner_ParseResult() {
	_ = this_
	this_.sourceMapURL = j2gort.NullString
}

// Java constructor: ParseResult(com.google.javascript.rhino.Node ast, java.util.List<com.google.javascript.jscomp.parsing.parser.trees.Comment> comments, com.google.javascript.jscomp.parsing.parser.FeatureSet features, java.lang.String sourceMapURL)
func NewParserRunner_ParseResult_4(ast pkg_rhino.Node, comments j2gort.Ext_java_util_List, features *pkg_parser.FeatureSet, sourceMapURL string) *ParserRunner_ParseResult {
	this_ := &ParserRunner_ParseResult{}
	InitParserRunner_ParseResult_4(this_, ast, comments, features, sourceMapURL)
	return this_
}

func InitParserRunner_ParseResult_4(this_ *ParserRunner_ParseResult, ast pkg_rhino.Node, comments j2gort.Ext_java_util_List, features *pkg_parser.FeatureSet, sourceMapURL string) {
	_ = ast
	_ = comments
	_ = features
	_ = sourceMapURL
	InitParserRunner_ParseResult_5(this_, ast, comments, j2gort.Ext_com_google_common_collect_ImmutableList_Of_0(), features, sourceMapURL)
}

// Java constructor: ParseResult(com.google.javascript.rhino.Node ast, java.util.List<com.google.javascript.jscomp.parsing.parser.trees.Comment> comments, com.google.common.collect.ImmutableList<com.google.javascript.jscomp.parsing.parser.util.SourceRange> topLevelStatementRanges, com.google.javascript.jscomp.parsing.parser.FeatureSet features, java.lang.String sourceMapURL)
func NewParserRunner_ParseResult_5(ast pkg_rhino.Node, comments j2gort.Ext_java_util_List, topLevelStatementRanges j2gort.Ext_com_google_common_collect_ImmutableList, features *pkg_parser.FeatureSet, sourceMapURL string) *ParserRunner_ParseResult {
	this_ := &ParserRunner_ParseResult{}
	InitParserRunner_ParseResult_5(this_, ast, comments, topLevelStatementRanges, features, sourceMapURL)
	return this_
}

func InitParserRunner_ParseResult_5(this_ *ParserRunner_ParseResult, ast pkg_rhino.Node, comments j2gort.Ext_java_util_List, topLevelStatementRanges j2gort.Ext_com_google_common_collect_ImmutableList, features *pkg_parser.FeatureSet, sourceMapURL string) {
	_ = ast
	_ = comments
	_ = topLevelStatementRanges
	_ = features
	_ = sourceMapURL
	this_.initFields_ParserRunner_ParseResult()
	this_.ast = ast
	this_.comments = comments
	this_.topLevelStatementRanges = topLevelStatementRanges
	this_.features = features
	this_.sourceMapURL = sourceMapURL
}

var parserPool_ sync.Pool

func acquirePooledParser_(config *pkg_parser.Parser_Config, errorReporter pkg_parser.ErrorReporter, source *pkg_parser.SourceFile) *pkg_parser.Parser {
	if v := parserPool_.Get(); v != nil {
		p := v.(*pkg_parser.Parser)
		p.ResetForReuse_(config, errorReporter, source)
		return p
	}
	return pkg_parser.NewParser(config, errorReporter, source)
}

func releasePooledParser_(p *pkg_parser.Parser) {
	if p == nil {
		return
	}
	p.ClearForPool_()
	parserPool_.Put(p)
}
