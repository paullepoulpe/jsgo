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
	"strings"
	"sync"
	"unicode/utf16"
	"unicode/utf8"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

var (
	parserRunner_defaultExtraAnnotations *jutil.LinkedHashSet[string]
	parserRunner_defaultAnnotations      map[string]Annotation
	parserRunner_suppressionNames        *jutil.LinkedHashSet[string]
	parserRunner_reservedVars            *jutil.LinkedHashSet[string]
	parserRunner_closurePrimitives       *jutil.LinkedHashSet[string]
	parserRunner_initOnce                sync.Once
)

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
	var annotations map[string]Annotation
	if extraAnnotationNames == nil || extraAnnotationNames.IsEmpty() {
		annotations = parserRunner_defaultAnnotations
	} else {
		var custom *jutil.LinkedHashSet[string]
		if !extra.IsEmpty() {
			custom = jutil.NewLinkedHashSetWithCapacity[string](extra.Size() + extraAnnotationNames.Size())
			for name := range extra.All() {
				custom.AddNew(name)
			}
		}
		for name := range extraAnnotationNames.All() {
			if name == "" {
				continue
			}
			if _, recognized := Annotation_recognizedAnnotations[name]; !recognized {
				if custom == nil {
					custom = jutil.NewLinkedHashSet[string]()
				}
				custom.Add(name)
			}
		}
		extra = custom
		if extra.IsEmpty() {
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
		suppressionNames:      parserRunner_suppressionNames,
		closurePrimitiveNames: parserRunner_closurePrimitives,
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
		annotations := parserRunner_extractList(ParserConfiguration_GetString("jsdoc.annotations"))
		parserRunner_suppressionNames = parserRunner_extractList(ParserConfiguration_GetString("jsdoc.suppressions"))
		parserRunner_closurePrimitives = parserRunner_extractList(ParserConfiguration_GetString("jsdoc.primitives"))
		parserRunner_reservedVars = parserRunner_extractList(ParserConfiguration_GetString("compiler.reserved.vars"))
		for name := range annotations.All() {
			if name == "" {
				continue
			}
			if _, recognized := Annotation_recognizedAnnotations[name]; !recognized {
				if parserRunner_defaultExtraAnnotations == nil {
					parserRunner_defaultExtraAnnotations = jutil.NewLinkedHashSet[string]()
				}
				parserRunner_defaultExtraAnnotations.Add(name)
			}
		}
		parserRunner_defaultAnnotations = Config_BuildAnnotations(annotations.All())
	})
}

func parserRunner_extractList(configProp string) *jutil.LinkedHashSet[string] {
	set := jutil.NewLinkedHashSetWithCapacity[string](strings.Count(configProp, ",") + 1)
	for part := range strings.SplitSeq(configProp, ",") {
		set.Add(parserRunner_trimWhitespace(part))
	}
	return set
}

func parserRunner_trimWhitespace(s string) string {
	first, last := 0, len(s)
	for first < last {
		if b := s[first]; b < 0x80 {
			if !jutil.CharMatcher_Whitespace_Matches(uint16(b)) {
				break
			}
			first++
			continue
		}
		r, w := utf8.DecodeRuneInString(s[first:last])
		if r > 0xFFFF || !jutil.CharMatcher_Whitespace_Matches(uint16(r)) {
			break
		}
		first += w
	}
	for last > first {
		if b := s[last-1]; b < 0x80 {
			if !jutil.CharMatcher_Whitespace_Matches(uint16(b)) {
				break
			}
			last--
			continue
		}
		r, w := utf8.DecodeLastRuneInString(s[first:last])
		if r > 0xFFFF || !jutil.CharMatcher_Whitespace_Matches(uint16(r)) {
			break
		}
		last -= w
	}
	return s[first:last]
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
	es6ErrorReporter := ParserRunner_Es6ErrorReporter{
		reporter:        errorReporter,
		reportAllErrors: keepGoing,
		sourceName:      sourceName,
	}
	es6config := ParserRunner_NewParserConfig(config)
	p := parser.NewParser(es6config, &es6ErrorReporter, file)
	tree := p.ParseProgram()
	var root rhino.Node
	var comments []parser.Comment
	var topLevelStatementRanges []parser.SourceRange
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
func ParserRunner_NewParserConfig(config *Config) parser.Parser_Config {
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
	result := make(map[string]string)
	pos := 0
lines:
	for pos < len(text) {
		start, end, next := propertiesNaturalLine(text, pos)
		pos = next
		i := propertiesSkipSpace(text, start, end)
		if i == end || text[i] == '#' || text[i] == '!' {
			continue
		}
		var logical string
		if !propertiesContinues(text, i, end) {
			logical = text[i:end]
		} else {
			var sb strings.Builder
			for {
				if !propertiesContinues(text, i, end) {
					sb.WriteString(text[i:end])
					break
				}
				sb.WriteString(text[i : end-1])
				if pos >= len(text) {
					if sb.Len() == 0 && next-end == 2 {
						continue lines
					}
					break
				}
				if sb.Len() == 0 {
					continue lines
				}
				start, end, next = propertiesNaturalLine(text, pos)
				pos = next
				i = propertiesSkipSpace(text, start, end)
				if i == end {
					break
				}
			}
			logical = sb.String()
		}
		key, value := propertiesSplitLine(logical)
		result[propertiesUnescape(key)] = propertiesUnescape(value)
	}
	return result
}

func propertiesNaturalLine(in string, pos int) (start, end, next int) {
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

func propertiesIsSpace(c byte) bool { return c == ' ' || c == '\t' || c == '\f' }

func propertiesSkipSpace(in string, i, end int) int {
	for i < end && propertiesIsSpace(in[i]) {
		i++
	}
	return i
}

func propertiesContinues(in string, i, end int) bool {
	n := 0
	for j := end - 1; j >= i && in[j] == '\\'; j-- {
		n++
	}
	return n%2 == 1
}

func propertiesSplitLine(line string) (key, value string) {
	k := 0
	for k < len(line) {
		c := line[k]
		if c == '\\' {
			k++
			if k < len(line) {
				_, w := utf8.DecodeRuneInString(line[k:])
				k += w
			}
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

func propertiesUnescape(s string) string {
	if strings.IndexByte(s, '\\') < 0 {
		return s
	}
	var out strings.Builder
	out.Grow(len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c != '\\' {
			out.WriteByte(c)
			continue
		}
		i++
		if i == len(s) {
			break
		}
		switch c = s[i]; c {
		case 't':
			out.WriteByte('\t')
		case 'n':
			out.WriteByte('\n')
		case 'r':
			out.WriteByte('\r')
		case 'f':
			out.WriteByte('\f')
		case 'u':
			v := propertiesParseHex4(s, i+1)
			i += 4
			if v >= 0xD800 && v <= 0xDBFF && i+6 < len(s) && s[i+1] == '\\' && s[i+2] == 'u' {
				v2 := propertiesParseHex4(s, i+3)
				if v2 >= 0xDC00 && v2 <= 0xDFFF {
					out.WriteRune(utf16.DecodeRune(rune(v), rune(v2)))
					i += 6
					continue
				}
			}
			if v >= 0xD800 && v <= 0xDFFF {
				out.WriteRune(utf8.RuneError)
			} else {
				out.WriteRune(rune(v))
			}
		default:
			if c < 0x80 {
				out.WriteByte(c)
			} else {
				r, w := utf8.DecodeRuneInString(s[i:])
				out.WriteRune(r)
				i += w - 1
			}
		}
	}
	return out.String()
}

func propertiesParseHex4(s string, pos int) uint16 {
	var v uint16
	for k := 0; k < 4; k++ {
		d := -1
		if pos+k < len(s) {
			d = propertiesHexDigit(s[pos+k])
		}
		if d < 0 {
			panic(jutil.NewIllegalArgumentExceptionMsg(`Malformed \uxxxx encoding.`))
		}
		v = v<<4 | uint16(d)
	}
	return v
}

func propertiesHexDigit(c byte) int {
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
