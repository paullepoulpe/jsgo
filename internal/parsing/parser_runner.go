// Port of com/google/javascript/jscomp/parsing/ParserRunner.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/ParserConfiguration.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/ParsingUtil.java @ closure-compiler v20261021

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

/*
 * Copyright 2020 The Closure Compiler Authors.
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

/*
 * Copyright 2019 The Closure Compiler Authors.
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

package parsing

import (
	_ "embed"
	"maps"
	"slices"
	"strings"
	"sync"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

var (
	parserRunner_defaultExtraAnnotations map[string]struct{}
	parserRunner_defaultExtraOrder       []string
	parserRunner_defaultAnnotations      map[string]Annotation
	parserRunner_suppressionNames        *jutil.LinkedHashSet[string]
	parserRunner_suppressionSet          map[string]struct{}
	parserRunner_suppressionOrder        []string
	parserRunner_reservedVars            *jutil.LinkedHashSet[string]
	parserRunner_closurePrimitiveSet     map[string]struct{}
	parserRunner_closurePrimitiveOrder   []string

	parserRunner_initOnce sync.Once
)

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
	extra := parserRunner_defaultExtraAnnotations
	extraOrder := parserRunner_defaultExtraOrder
	var annotations map[string]Annotation
	if extraAnnotationNames == nil || extraAnnotationNames.IsEmpty() {
		annotations = parserRunner_defaultAnnotations
	} else {
		var custom map[string]struct{}
		var customOrder []string
		if len(extra) > 0 {
			custom = maps.Clone(extra)
			customOrder = slices.Clone(extraOrder)
		}
		for name := range extraAnnotationNames.All() {
			if name == "" {
				continue
			}
			if _, recognized := Annotation_recognizedAnnotations[name]; !recognized {
				if custom == nil {
					custom = make(map[string]struct{})
				}
				if _, dup := custom[name]; !dup {
					custom[name] = struct{}{}
					customOrder = append(customOrder, name)
				}
			}
		}
		extra = custom
		extraOrder = customOrder
		if len(extra) == 0 {
			annotations = Annotation_recognizedAnnotations
		}
	}
	return &Config{
		languageMode:          languageMode,
		strictMode:            strictMode,
		jsDocParsingMode:      jsdocParsingMode,
		runMode:               runMode,
		annotations:           annotations,
		extraAnnotations:      extra,
		extraAnnotationOrder:  extraOrder,
		suppressionNames:      parserRunner_suppressionSet,
		suppressionOrder:      parserRunner_suppressionOrder,
		closurePrimitiveNames: parserRunner_closurePrimitiveSet,
		closurePrimitiveOrder: parserRunner_closurePrimitiveOrder,
		parseInlineSourceMaps: parseInlineSourceMaps,
	}
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
	parserRunner_initOnce.Do(func() {
		annotationNames := ParserRunner_ExtractList(ParserConfiguration_GetString("jsdoc.annotations"))
		parserRunner_suppressionNames = ParserRunner_ExtractList(ParserConfiguration_GetString("jsdoc.suppressions"))
		closurePrimitiveNames := ParserRunner_ExtractList(ParserConfiguration_GetString("jsdoc.primitives"))
		parserRunner_reservedVars = ParserRunner_ExtractList(ParserConfiguration_GetString("compiler.reserved.vars"))
		parserRunner_suppressionSet, parserRunner_suppressionOrder = config_stringSetAndOrderCopyOf(parserRunner_suppressionNames.All())
		parserRunner_closurePrimitiveSet, parserRunner_closurePrimitiveOrder = config_stringSetAndOrderCopyOf(closurePrimitiveNames.All())
		for name := range annotationNames.All() {
			if name == "" {
				continue
			}
			if _, recognized := Annotation_recognizedAnnotations[name]; !recognized {
				if parserRunner_defaultExtraAnnotations == nil {
					parserRunner_defaultExtraAnnotations = make(map[string]struct{})
				}
				if _, dup := parserRunner_defaultExtraAnnotations[name]; !dup {
					parserRunner_defaultExtraAnnotations[name] = struct{}{}
					parserRunner_defaultExtraOrder = append(parserRunner_defaultExtraOrder, name)
				}
			}
		}
		parserRunner_defaultAnnotations = Config_BuildAnnotations(annotationNames.All())
	})
}

// Java: ParserRunner.extractList(String)
func ParserRunner_ExtractList(configProp string) *jutil.LinkedHashSet[string] {
	parts := strings.Split(configProp, ",")
	for i, part := range parts {
		u := jutil.StringToUTF16(part)
		first, last := 0, len(u)-1
		for first < len(u) && jutil.CharMatcher_Whitespace_Matches(u[first]) {
			first++
		}
		for last > first && jutil.CharMatcher_Whitespace_Matches(u[last]) {
			last--
		}
		parts[i] = jutil.StringFromUTF16(u[first : last+1])
	}
	return jutil.NewLinkedHashSet(parts...)
}

// Java: ParserRunner.parse(StaticSourceFile, String, Config, ErrorReporter)
func ParserRunner_Parse(
	sourceFile rhino.StaticSourceFile,
	sourceString string,
	config *Config,
	errorReporter rhino.ErrorReporter) *ParserRunner_ParseResult {
	return ParserRunner_ParseWithAllocator(sourceFile, sourceString, config, errorReporter, nil)
}

// ParserRunner_ParseWithAllocator is ParserRunner_Parse taking the AST's
// node ids from alloc (see rhino.NodeAllocator; nil: the store's counter).
func ParserRunner_ParseWithAllocator(
	sourceFile rhino.StaticSourceFile,
	sourceString string,
	config *Config,
	errorReporter rhino.ErrorReporter,
	alloc *rhino.NodeAllocator) *ParserRunner_ParseResult {
	sourceName := sourceFile.GetName()
	defer func() {
		if r := recover(); r != nil {
			panic(jutil.NewRuntimeExceptionWithMessageAndCause("Exception parsing \""+sourceName+"\"", jutil.AsThrowable(r)))
		}
	}()
	file := parser.NewSourceFile(sourceName, sourceString)
	keepGoing := config.RunMode() == Config_RunMode_KEEP_GOING
	es6ErrorReporter := NewParserRunner_Es6ErrorReporter(errorReporter, keepGoing)
	es6ErrorReporter.sourceName = sourceName
	es6config := ParserRunner_NewParserConfig(config)
	p := parser.NewParser(es6config, es6ErrorReporter, file)
	tree := p.ParseProgram()
	var root rhino.Node
	comments := []parser.Comment{}
	topLevelStatementRanges := []parser.SourceRange{}
	features := p.GetFeatures()
	if tree != nil && (!es6ErrorReporter.HadError() || keepGoing) {
		factory := irFactoryTransformTree(
			tree, sourceFile, config, errorReporter, file, alloc)
		root = factory.GetResultNode()
		features = features.Union(factory.GetFeatures())
		root.PutProp(rhino.Node_FEATURE_SET, features)

		if config.JsDocParsingMode().ShouldParseDescriptions() {
			comments = p.GetComments()
			topLevelStatementRanges = make([]parser.SourceRange, len(tree.SourceElements))
			for i, element := range tree.SourceElements {
				topLevelStatementRanges[i] = element.GetParseTreeBase().Location
			}
		}
	}
	return NewParserRunner_ParseResultWithTopLevelStatementRanges(
		root, comments, topLevelStatementRanges, features, p.GetSourceMapURL())
}

// Java: ParserRunner.newParserConfig(Config)
func ParserRunner_NewParserConfig(config *Config) *parser.Parser_Config {
	var mode parser.Parser_Config_Mode
	switch config.LanguageMode() {
	case Config_LanguageMode_ECMASCRIPT3:
		mode = parser.Parser_Config_Mode_ES3
	case Config_LanguageMode_ECMASCRIPT5:
		mode = parser.Parser_Config_Mode_ES5
	case Config_LanguageMode_ECMASCRIPT_2015, Config_LanguageMode_ECMASCRIPT_2016:
		mode = parser.Parser_Config_Mode_ES6_OR_ES7
	case Config_LanguageMode_ECMASCRIPT_2017,
		Config_LanguageMode_ECMASCRIPT_2018,
		Config_LanguageMode_ECMASCRIPT_2019,
		Config_LanguageMode_ECMASCRIPT_2020,
		Config_LanguageMode_ECMASCRIPT_2021,
		Config_LanguageMode_ECMASCRIPT_2022,
		Config_LanguageMode_ES_NEXT,
		Config_LanguageMode_UNSTABLE,
		Config_LanguageMode_UNSUPPORTED:
		mode = parser.Parser_Config_Mode_ES8_OR_GREATER
	default:
		panic(jutil.NewIllegalStateException())
	}
	return parser.NewParser_ConfigWithMode(mode, config.StrictMode().IsStrict())
}

// ParserRunner_Es6ErrorReporter forwards the parser's errors and warnings to
// a rhino ErrorReporter.
type ParserRunner_Es6ErrorReporter struct {
	reporter        rhino.ErrorReporter
	errorSeen       bool
	reportAllErrors bool
	sourceName      string
}

func NewParserRunner_Es6ErrorReporter(reporter rhino.ErrorReporter, reportAllErrors bool) *ParserRunner_Es6ErrorReporter {
	return &ParserRunner_Es6ErrorReporter{
		reporter:        reporter,
		reportAllErrors: reportAllErrors,
	}
}

func (r *ParserRunner_Es6ErrorReporter) ReportError(location parser.LinePosition, message string) {
	if r.reportAllErrors || !r.errorSeen {
		r.errorSeen = true
		r.reporter.Error(message, r.sourceName, int(location.Line)+1, int(location.Column))
	}
}

func (r *ParserRunner_Es6ErrorReporter) ReportWarning(location parser.LinePosition, message string) {
	r.reporter.Warning(message, r.sourceName, int(location.Line)+1, int(location.Column))
}

func (r *ParserRunner_Es6ErrorReporter) HadError() bool {
	return r.errorSeen
}

// ParserRunner_ParseResult holds results of parsing.
type ParserRunner_ParseResult struct {
	Ast                     rhino.Node
	Comments                []parser.Comment
	TopLevelStatementRanges []parser.SourceRange
	Features                *parser.FeatureSet
	SourceMapURL            *string
}

func NewParserRunner_ParseResult(ast rhino.Node, comments []parser.Comment, features *parser.FeatureSet, sourceMapURL *string) *ParserRunner_ParseResult {
	return NewParserRunner_ParseResultWithTopLevelStatementRanges(ast, comments, []parser.SourceRange{}, features, sourceMapURL)
}

func NewParserRunner_ParseResultWithTopLevelStatementRanges(ast rhino.Node, comments []parser.Comment, topLevelStatementRanges []parser.SourceRange, features *parser.FeatureSet, sourceMapURL *string) *ParserRunner_ParseResult {
	return &ParserRunner_ParseResult{
		Ast:                     ast,
		Comments:                comments,
		TopLevelStatementRanges: topLevelStatementRanges,
		Features:                features,
		SourceMapURL:            sourceMapURL,
	}
}

//go:embed ParserConfig.properties
var parserConfiguration_properties string

var (
	parserConfiguration_bundleOnce sync.Once
	parserConfiguration_bundle     map[string]string
)

// ParserConfiguration_GetString reads a key from ParserConfig.properties.
func ParserConfiguration_GetString(name string) string {
	parserConfiguration_bundleOnce.Do(func() {
		parserConfiguration_bundle = propertiesLoad(parserConfiguration_properties)
	})
	value, ok := parserConfiguration_bundle[name]
	if !ok {
		panic(jutil.NewIllegalStateExceptionMsg("missing ParserConfig.properties resource: " + name))
	}
	return value
}

func propertiesLoad(text string) map[string]string {
	in := jutil.StringToUTF16(text)
	result := make(map[string]string)
	pos := 0
lines:
	for pos < len(in) {
		start, end, next := propertiesNaturalLine(in, pos)
		pos = next
		i := propertiesSkipSpace(in, start, end)
		if i == end || in[i] == '#' || in[i] == '!' {
			continue
		}
		var logical []uint16
		for {
			if !propertiesContinues(in, i, end) {
				logical = append(logical, in[i:end]...)
				break
			}
			logical = append(logical, in[i:end-1]...)
			if pos >= len(in) {
				if len(logical) == 0 && next-end == 2 {
					continue lines
				}
				break
			}
			if len(logical) == 0 {
				continue lines
			}
			start, end, next = propertiesNaturalLine(in, pos)
			pos = next
			i = propertiesSkipSpace(in, start, end)
			if i == end {
				break
			}
		}
		key, value := propertiesSplitLine(logical)
		result[jutil.StringFromUTF16(propertiesUnescape(key))] = jutil.StringFromUTF16(propertiesUnescape(value))
	}
	return result
}

func propertiesNaturalLine(in []uint16, pos int) (start, end, next int) {
	end = pos
	for end < len(in) && in[end] != '\n' && in[end] != '\r' {
		end++
	}
	next = end
	if next < len(in) {
		if in[next] == '\r' && next+1 < len(in) && in[next+1] == '\n' {
			next++
		}
		next++
	}
	return pos, end, next
}

func propertiesIsSpace(c uint16) bool { return c == ' ' || c == '\t' || c == '\f' }

func propertiesSkipSpace(in []uint16, i, end int) int {
	for i < end && propertiesIsSpace(in[i]) {
		i++
	}
	return i
}

func propertiesContinues(in []uint16, i, end int) bool {
	n := 0
	for j := end - 1; j >= i && in[j] == '\\'; j-- {
		n++
	}
	return n%2 == 1
}

func propertiesSplitLine(line []uint16) (key, value []uint16) {
	k := 0
	for k < len(line) {
		c := line[k]
		if c == '\\' {
			k += 2
			continue
		}
		if c == '=' || c == ':' || propertiesIsSpace(c) {
			break
		}
		k++
	}
	if k > len(line) {
		k = len(line)
	}
	v := propertiesSkipSpace(line, k, len(line))
	if v < len(line) && (line[v] == '=' || line[v] == ':') {
		v = propertiesSkipSpace(line, v+1, len(line))
	}
	return line[:k], line[v:]
}

func propertiesUnescape(s []uint16) []uint16 {
	out := make([]uint16, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c != '\\' {
			out = append(out, c)
			continue
		}
		i++
		if i == len(s) {
			break
		}
		switch c = s[i]; c {
		case 't':
			out = append(out, '\t')
		case 'n':
			out = append(out, '\n')
		case 'r':
			out = append(out, '\r')
		case 'f':
			out = append(out, '\f')
		case 'u':
			var v uint16
			for k := 1; k <= 4; k++ {
				d := -1
				if i+k < len(s) {
					d = propertiesHexDigit(s[i+k])
				}
				if d < 0 {
					panic(jutil.NewIllegalArgumentExceptionMsg(`Malformed \uxxxx encoding.`))
				}
				v = v<<4 | uint16(d)
			}
			out = append(out, v)
			i += 4
		default:
			out = append(out, c)
		}
	}
	return out
}

func propertiesHexDigit(c uint16) int {
	switch {
	case c >= '0' && c <= '9':
		return int(c - '0')
	case c >= 'a' && c <= 'f':
		return int(c-'a') + 10
	case c >= 'A' && c <= 'F':
		return int(c-'A') + 10
	}
	return -1
}

// ParsingUtil_GetParamOrPatternNames calls cb with all NAMEs declared in a
// PARAM_LIST or destructuring pattern.
func ParsingUtil_GetParamOrPatternNames(n rhino.Node, cb func(rhino.Node)) {
	rhino.CheckNotNull(n)
	switch n.GetToken() {
	case rhino.Token_EMPTY, rhino.Token_GETELEM, rhino.Token_GETPROP:
		// An empty node or assignment target means nothing declared here.
	case rhino.Token_NAME:
		cb(n)
	case rhino.Token_ITER_REST, rhino.Token_OBJECT_REST, rhino.Token_STRING_KEY, rhino.Token_DEFAULT_VALUE:
		ParsingUtil_GetParamOrPatternNames(n.GetFirstChild(), cb)
	case rhino.Token_COMPUTED_PROP:
		ParsingUtil_GetParamOrPatternNames(n.GetSecondChild(), cb)
	case rhino.Token_OBJECT_PATTERN, rhino.Token_ARRAY_PATTERN, rhino.Token_PARAM_LIST:
		for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
			ParsingUtil_GetParamOrPatternNames(c, cb)
		}
	default:
		panic(jutil.NewIllegalStateExceptionMsg("Unexpected parameter structure"))
	}
}
