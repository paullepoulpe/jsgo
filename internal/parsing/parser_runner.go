// Port of com/google/javascript/jscomp/parsing/ParserRunner.java @ closure-compiler v20260415

package parsing

import (
	"sync"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ParserRunner is the parser runner.
//
// Java: final class ParserRunner
type ParserRunner struct{}

var (
	parserRunner_annotationNames       *jutil.LinkedHashSet[string]
	parserRunner_suppressionNames      *jutil.LinkedHashSet[string]
	parserRunner_reservedVars          *jutil.LinkedHashSet[string]
	parserRunner_closurePrimitiveNames *jutil.LinkedHashSet[string]

	// Go-only: language-adapter — the monitor of the static synchronized
	// initResourceConfig().
	parserRunner_initResourceConfigMu sync.Mutex
)

// Should never need to instantiate class of static methods.
//
// Java: ParserRunner()
func NewParserRunner() *ParserRunner { return &ParserRunner{} }

// Java: ParserRunner.createConfig(LanguageMode, Set<String>, StrictMode)
func ParserRunner_CreateConfig(
	languageMode Config_LanguageMode, extraAnnotationNames *jutil.LinkedHashSet[string], strictMode Config_StrictMode) *Config {
	return ParserRunner_CreateConfigWithModes(
		languageMode,
		Config_JsDocParsing_TYPES_ONLY,
		Config_RunMode_STOP_AFTER_ERROR,
		extraAnnotationNames,
		true,
		strictMode)
}

// Java: ParserRunner.createConfig(LanguageMode, JsDocParsing, RunMode, Set<String>, boolean, StrictMode)
func ParserRunner_CreateConfigWithModes(
	languageMode Config_LanguageMode,
	jsdocParsingMode Config_JsDocParsing,
	runMode Config_RunMode,
	extraAnnotationNames *jutil.LinkedHashSet[string],
	parseInlineSourceMaps bool,
	strictMode Config_StrictMode) *Config {

	ParserRunner_InitResourceConfig()
	var effectiveAnnotationNames *jutil.LinkedHashSet[string]
	if extraAnnotationNames == nil {
		effectiveAnnotationNames = parserRunner_annotationNames
	} else {
		effectiveAnnotationNames = jutil.NewLinkedHashSet(parserRunner_annotationNames.ToSlice()...)
		effectiveAnnotationNames.AddAll(extraAnnotationNames.ToSlice()...)
	}
	return Config_Builder_().
		SetExtraAnnotationNames(effectiveAnnotationNames.All()).
		SetJsDocParsingMode(jsdocParsingMode).
		SetRunMode(runMode).
		SetSuppressionNames(parserRunner_suppressionNames.All()).
		SetClosurePrimitiveNames(parserRunner_closurePrimitiveNames.All()).
		SetLanguageMode(languageMode).
		SetParseInlineSourceMaps(parseInlineSourceMaps).
		SetStrictMode(strictMode).
		Build()
}

// Java: ParserRunner.getReservedVars()
func ParserRunner_GetReservedVars() *jutil.LinkedHashSet[string] {
	ParserRunner_InitResourceConfig()
	return parserRunner_reservedVars
}

// Java: ParserRunner.getSuppressionNames()
func ParserRunner_GetSuppressionNames() *jutil.LinkedHashSet[string] {
	ParserRunner_InitResourceConfig()
	return parserRunner_suppressionNames
}

// Java: ParserRunner.initResourceConfig()
func ParserRunner_InitResourceConfig() {
	parserRunner_initResourceConfigMu.Lock()
	defer parserRunner_initResourceConfigMu.Unlock()
	if parserRunner_annotationNames != nil {
		return
	}

	parserRunner_annotationNames = ParserRunner_ExtractList(ParserConfiguration_GetString("jsdoc.annotations"))
	parserRunner_suppressionNames = ParserRunner_ExtractList(ParserConfiguration_GetString("jsdoc.suppressions"))
	parserRunner_closurePrimitiveNames = ParserRunner_ExtractList(ParserConfiguration_GetString("jsdoc.primitives"))
	parserRunner_reservedVars = ParserRunner_ExtractList(ParserConfiguration_GetString("compiler.reserved.vars"))
}

// Java: ParserRunner.extractList(String)
func ParserRunner_ExtractList(configProp string) *jutil.LinkedHashSet[string] {
	return jutil.NewLinkedHashSet(jutil.Splitter_OnCharTrimResults(configProp, ',')...)
}

// Java: ParserRunner.parse(StaticSourceFile, String, Config, ErrorReporter)
func ParserRunner_Parse(
	sourceFile rhino.StaticSourceFile,
	sourceString string,
	config *Config,
	errorReporter rhino.ErrorReporter) *ParserRunner_ParseResult {
	// TODO(johnlenz): unify "SourceFile", "Es6ErrorReporter" and "Config"

	sourceName := sourceFile.GetName()
	var result *ParserRunner_ParseResult
	jutil.Try(func() {
		file := parser.NewSourceFile(sourceName, sourceString)
		keepGoing := config.RunMode() == Config_RunMode_KEEP_GOING
		es6ErrorReporter := NewParserRunner_Es6ErrorReporter(errorReporter, keepGoing)
		es6config := ParserRunner_NewParserConfig(config)
		p := parser.NewParser(es6config, es6ErrorReporter, file)
		tree := p.ParseProgram()
		var root *rhino.Node
		comments := []*parser.Comment{}
		features := p.GetFeatures()
		if tree != nil && (!es6ErrorReporter.HadError() || keepGoing) {
			factory := IRFactory_TransformTree(tree, sourceFile, config, errorReporter, file)
			root = factory.GetResultNode()
			features = features.Union(factory.GetFeatures())
			root.PutProp(rhino.Node_FEATURE_SET, features)

			if config.JsDocParsingMode().ShouldParseDescriptions() {
				comments = p.GetComments()
			}
		}
		result = NewParserRunner_ParseResult(root, comments, features, p.GetSourceMapURL())
	}, func(t jutil.Throwable) {
		panic(jutil.NewRuntimeExceptionWithMessageAndCause("Exception parsing \""+sourceName+"\"", t))
	})
	return result
}

// Java: ParserRunner.newParserConfig(Config)
func ParserRunner_NewParserConfig(config *Config) *parser.Parser_Config {
	languageMode := config.LanguageMode()
	isStrictMode := config.StrictMode().IsStrict()
	var parserConfigLanguageMode parser.Parser_Config_Mode
	switch languageMode {
	case Config_LanguageMode_ECMASCRIPT3:
		parserConfigLanguageMode = parser.Parser_Config_Mode_ES3
	case Config_LanguageMode_ECMASCRIPT5:
		parserConfigLanguageMode = parser.Parser_Config_Mode_ES5
	case Config_LanguageMode_ECMASCRIPT_2015, Config_LanguageMode_ECMASCRIPT_2016:
		parserConfigLanguageMode = parser.Parser_Config_Mode_ES6_OR_ES7
	case Config_LanguageMode_ECMASCRIPT_2017,
		Config_LanguageMode_ECMASCRIPT_2018,
		Config_LanguageMode_ECMASCRIPT_2019,
		Config_LanguageMode_ECMASCRIPT_2020,
		Config_LanguageMode_ECMASCRIPT_2021,
		Config_LanguageMode_ECMASCRIPT_2022,
		Config_LanguageMode_ES_NEXT,
		Config_LanguageMode_UNSTABLE,
		Config_LanguageMode_UNSUPPORTED:
		parserConfigLanguageMode = parser.Parser_Config_Mode_ES8_OR_GREATER
	default:
		// Go-only: language-adapter — an exhaustive Java switch expression
		// throws MatchException for an enum constant it does not know.
		panic(jutil.NewIllegalStateException())
	}
	return parser.NewParser_ConfigWithMode(jutil.CheckNotNull(parserConfigLanguageMode), isStrictMode)
}

// ParserRunner_Es6ErrorReporter forwards the parser's errors and warnings to
// a rhino ErrorReporter.
//
// Java: class ParserRunner.Es6ErrorReporter extends parser.util.ErrorReporter
type ParserRunner_Es6ErrorReporter struct {
	parser.ErrorReporterBase

	reporter        rhino.ErrorReporter
	errorSeen       bool
	reportAllErrors bool
}

// Java: ParserRunner.Es6ErrorReporter(ErrorReporter, boolean)
func NewParserRunner_Es6ErrorReporter(reporter rhino.ErrorReporter, reportAllErrors bool) *ParserRunner_Es6ErrorReporter {
	r := &ParserRunner_Es6ErrorReporter{}
	r.InitErrorReporterBase(r)
	r.reporter = reporter
	r.reportAllErrors = reportAllErrors
	return r
}

// Java: ParserRunner.Es6ErrorReporter.reportError(SourcePosition, String)
func (r *ParserRunner_Es6ErrorReporter) ReportError(location *parser.SourcePosition, message string) {
	// In normal usage, only the first parse error should be reported, but
	// sometimes it is useful to keep going.
	if r.reportAllErrors || !r.errorSeen {
		r.errorSeen = true
		r.reporter.Error(message, location.Source.Name, location.Line+1, location.Column)
	}
}

// Java: ParserRunner.Es6ErrorReporter.reportWarning(SourcePosition, String)
func (r *ParserRunner_Es6ErrorReporter) ReportWarning(location *parser.SourcePosition, message string) {
	r.reporter.Warning(message, location.Source.Name, location.Line+1, location.Column)
}

// ParserRunner_ParseResult holds results of parsing.
//
// Java: class ParserRunner.ParseResult
type ParserRunner_ParseResult struct {
	Ast      *rhino.Node
	Comments []*parser.Comment
	Features *parser.FeatureSet
	// Go-only: language-adapter — nil is Java null.
	SourceMapURL *string
}

// Java: ParserRunner.ParseResult(Node, List<Comment>, FeatureSet, String)
func NewParserRunner_ParseResult(ast *rhino.Node, comments []*parser.Comment, features *parser.FeatureSet, sourceMapURL *string) *ParserRunner_ParseResult {
	r := &ParserRunner_ParseResult{}
	r.Ast = ast
	r.Comments = comments
	r.Features = features
	r.SourceMapURL = sourceMapURL
	return r
}
