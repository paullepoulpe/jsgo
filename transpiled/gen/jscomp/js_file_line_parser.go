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
// Port of com/google/javascript/jscomp/deps/JsFileLineParser.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javaio"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java polymorphic class interface: com.google.javascript.jscomp.deps.JsFileLineParser
type JsFileLineParser interface {
	Implements_jscomp_JsFileLineParser()
	IsJavaNil_() bool
	GetJsFileLineParserBase() *JsFileLineParserBase
	SetShortcutMode(bool)
	DidParseSucceed() bool
	DoParse(string, j2gort.Ext_java_io_Reader)
	ParseLine(string) bool
	ParseJsDocCommentLine(string) bool
	ParseJsString(string) string
	ParseJsStringArray(string) j2gort.Ext_java_util_List
	ParseJsStringMap(string) j2gort.Ext_java_util_Map
}

func ToIface_JsFileLineParser[T interface {
	JsFileLineParser
	comparable
}](v T) JsFileLineParser {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_JsFileLineParser[T JsFileLineParser](v T) JsFileLineParser {
	var vi JsFileLineParser = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.jscomp.deps.JsFileLineParser
type JsFileLineParserBase struct {
	self_          JsFileLineParser
	shortcutMode   bool                                `gson:"shortcutMode"`
	valueMatcher   *j2gort.Ext_java_util_regex_Matcher `gson:"valueMatcher"`
	filePath       string                              `gson:"filePath"`
	lineNum        int32                               `gson:"lineNum"`
	errorManager   ErrorManager                        `gson:"errorManager"`
	parseSucceeded bool                                `gson:"parseSucceeded"`
	idHash_        atomic.Int32
}

func (this_ *JsFileLineParserBase) BindSelf_(s any) {
	_ = s
	this_.self_ = s.(JsFileLineParser)
}
func (this_ *JsFileLineParserBase) IsJavaNil_() bool { return this_ == nil }
func (*JsFileLineParserBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.deps.JsFileLineParser", "JsFileLineParser"
}
func (this_ *JsFileLineParserBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.deps.JsFileLineParser", "com.google.javascript.jscomp.deps.JsFileLineParser", func(v_ any) bool {
		return j2gort.InstanceOf[JsFileLineParser](v_)
	})
}
func (b *JsFileLineParserBase) GetJsFileLineParserBase() *JsFileLineParserBase { return b }
func (*JsFileLineParserBase) Implements_jscomp_JsFileLineParser()              {}

var static_JsFileLineParser_PARSE_WARNING *DiagnosticType

func StaticGet_JsFileLineParser_PARSE_WARNING() *DiagnosticType {
	clinit_JsFileLineParser()
	return static_JsFileLineParser_PARSE_WARNING
}
func StaticSet_JsFileLineParser_PARSE_WARNING(v *DiagnosticType) {
	clinit_JsFileLineParser()
	static_JsFileLineParser_PARSE_WARNING = v
}

var static_JsFileLineParser_PARSE_ERROR *DiagnosticType

func StaticGet_JsFileLineParser_PARSE_ERROR() *DiagnosticType {
	clinit_JsFileLineParser()
	return static_JsFileLineParser_PARSE_ERROR
}
func StaticSet_JsFileLineParser_PARSE_ERROR(v *DiagnosticType) {
	clinit_JsFileLineParser()
	static_JsFileLineParser_PARSE_ERROR = v
}

var static_JsFileLineParser_STRING_LITERAL_PATTERN *j2gort.Ext_java_util_regex_Pattern

func StaticGet_JsFileLineParser_STRING_LITERAL_PATTERN() *j2gort.Ext_java_util_regex_Pattern {
	clinit_JsFileLineParser()
	return static_JsFileLineParser_STRING_LITERAL_PATTERN
}
func StaticSet_JsFileLineParser_STRING_LITERAL_PATTERN(v *j2gort.Ext_java_util_regex_Pattern) {
	clinit_JsFileLineParser()
	static_JsFileLineParser_STRING_LITERAL_PATTERN = v
}

var clinitDone_JsFileLineParser uint32
var clinitState_JsFileLineParser j2gort.ClinitState

func clinit_JsFileLineParser() {
	if atomic.LoadUint32(&clinitDone_JsFileLineParser) != 0 {
		return
	}
	clinitSlow_JsFileLineParser()
}

func clinitSlow_JsFileLineParser() {
	if !j2gort.EnterClinit(&clinitState_JsFileLineParser, &clinitDone_JsFileLineParser) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_JsFileLineParser, &clinitDone_JsFileLineParser)
	static_JsFileLineParser_PARSE_WARNING = DiagnosticType_Warning("DEPS_PARSE_WARNING", "{0}\n{1}")
	static_JsFileLineParser_PARSE_ERROR = DiagnosticType_Error("DEPS_PARSE_ERROR", "{0}\n{1}")
	static_JsFileLineParser_STRING_LITERAL_PATTERN = j2gort.Ext_java_util_regex_Pattern_Compile_1("\\s*(?:'((?:\\\\.|[^'\\\\])*)'|\"((?:\\\\.|[^\"\\\\])*)\")\\s*")
}

func (this_ *JsFileLineParserBase) initFields_JsFileLineParser() {
	_ = this_

	this_.valueMatcher = StaticGet_JsFileLineParser_STRING_LITERAL_PATTERN().Matcher("")
	this_.filePath = j2gort.NullString
}

func InitJsFileLineParser(this_ *JsFileLineParserBase, errorManager ErrorManager) {
	_ = errorManager
	this_.initFields_JsFileLineParser()
	this_.errorManager = errorManager
}

// Java method: JsFileLineParser.setShortcutMode(boolean mode)
func (this_ *JsFileLineParserBase) SetShortcutMode(mode bool) {
	_ = mode
	this_.shortcutMode = mode
}

// Java method: JsFileLineParser.didParseSucceed()
func (this_ *JsFileLineParserBase) DidParseSucceed() bool {
	return this_.parseSucceeded
}

// Java method: JsFileLineParser.doParse(java.lang.String filePath, java.io.Reader fileContents)
func (this_ *JsFileLineParserBase) DoParse(filePath string, fileContents j2gort.Ext_java_io_Reader) {
	_ = filePath
	_ = fileContents
	this_.filePath = filePath
	this_.parseSucceeded = true
	var lineBuffer *j2gort.Ext_java_io_BufferedReader = j2gort.NewExt_java_io_BufferedReader_1(fileContents)
	_ = lineBuffer
	var line string = j2gort.NullString
	_ = line
	this_.lineNum = int32(0)
	var inMultilineComment bool = false
	_ = inMultilineComment
	var inJsDocComment bool = false
	_ = inJsDocComment
	func() {
		defer func() {
			if rec_1 := recover(); rec_1 != nil {
				jex_2 := javalang.RecoverException(rec_1)
				_ = jex_2
				if j2gort.InstanceOf[javaio.IIOException](jex_2) {
					e := j2gort.AsIface[javaio.IIOException](jex_2)
					_ = e
					this_.errorManager.Report_2(StaticGet_CheckLevel_ERROR(), JSError_Make_5(filePath, int32(0), int32(0), StaticGet_JsFileLineParser_PARSE_ERROR(), []string{("Error reading file: " + j2gort.StrValString(filePath))}))
					this_.parseSucceeded = false
				} else {
					panic(rec_1)
				}
			}
		}()
		for (func() string { var asgn_3 string = lineBuffer.ReadLine(); line = asgn_3; return asgn_3 }()) != j2gort.NullString {
			this_.lineNum = this_.lineNum + 1
			var tryFlow_4 int
			func() {
				defer func() {
					if rec_5 := recover(); rec_5 != nil {
						jex_6 := javalang.RecoverException(rec_5)
						_ = jex_6
						if j2gort.InstanceOf[*JsFileLineParser_ParseException](jex_6) {
							e_2 := j2gort.AsIface[*JsFileLineParser_ParseException](jex_6)
							_ = e_2
							this_.errorManager.Report_2(func() *CheckLevel {
								if e_2.IsFatal() {
									return StaticGet_CheckLevel_ERROR()
								}
								return StaticGet_CheckLevel_WARNING()
							}(), JSError_Make_5(filePath, this_.lineNum, int32(0), func() *DiagnosticType {
								if e_2.IsFatal() {
									return StaticGet_JsFileLineParser_PARSE_ERROR()
								}
								return StaticGet_JsFileLineParser_PARSE_WARNING()
							}(), []string{j2gort.ThrowableGetMessage(e_2), line}))
							this_.parseSucceeded = (this_.parseSucceeded && !(e_2.IsFatal()))
						} else {
							panic(rec_5)
						}
					}
				}()
				var revisedLine string = line
				_ = revisedLine
				var revisedJsDocCommentLine string = ""
				_ = revisedJsDocCommentLine
				if inMultilineComment {
					var endOfComment int32 = j2gort.JStrIndexOf(revisedLine, "*/")
					_ = endOfComment
					if endOfComment != int32(-1) {
						if inJsDocComment {
							revisedJsDocCommentLine = j2gort.JStrSub(revisedLine, int32(0), (endOfComment + int32(2)))
							inJsDocComment = false
						}
						revisedLine = j2gort.JStrSubFrom(revisedLine, (endOfComment + int32(2)))
						inMultilineComment = false
					} else {
						if inJsDocComment {
							revisedJsDocCommentLine = line
						}
						revisedLine = ""
					}
				}
				if !(inMultilineComment) {
					for {
						var startOfLineComment int32 = j2gort.JStrIndexOf(revisedLine, "//")
						_ = startOfLineComment
						var startOfMultilineComment int32 = j2gort.JStrIndexOf(revisedLine, "/*")
						_ = startOfMultilineComment
						if (startOfLineComment != int32(-1)) && ((startOfMultilineComment == int32(-1)) || (startOfLineComment < startOfMultilineComment)) {
							revisedLine = j2gort.JStrSub(revisedLine, int32(0), startOfLineComment)
							break
						} else if startOfMultilineComment != int32(-1) {
							if this_.priv_JsFileLineParser_isCommentQuoted_6098(revisedLine, startOfMultilineComment, uint16(39)) {
								break
							}
							if this_.priv_JsFileLineParser_isCommentQuoted_6098(revisedLine, startOfMultilineComment, uint16(34)) {
								break
							}
							if startOfMultilineComment == j2gort.JStrIndexOf(revisedLine, "/**") {
								inJsDocComment = true
							}
							var endOfMultilineComment int32 = j2gort.JStrIndexOfFrom(revisedLine, "*/", (startOfMultilineComment + int32(2)))
							_ = endOfMultilineComment
							if endOfMultilineComment == int32(-1) {
								if inJsDocComment {
									revisedJsDocCommentLine = j2gort.JStrSubFrom(revisedLine, startOfMultilineComment)
								}
								revisedLine = j2gort.JStrSub(revisedLine, int32(0), startOfMultilineComment)
								inMultilineComment = true
								break
							} else {
								if inJsDocComment {
									var jsDocComment string = j2gort.JStrSub(revisedLine, startOfMultilineComment, (endOfMultilineComment + int32(2)))
									_ = jsDocComment
									if !(this_.self_.ParseJsDocCommentLine(jsDocComment)) && this_.shortcutMode {
										break
									}
									inJsDocComment = false
								}
								revisedLine = (j2gort.StrValString(j2gort.JStrSub(revisedLine, int32(0), startOfMultilineComment)) + j2gort.StrValString(j2gort.JStrSubFrom(revisedLine, (endOfMultilineComment+int32(2)))))
							}
						} else {
							break
						}
					}
				}
				if !(j2gort.CheckNotNullString(revisedJsDocCommentLine) == "") {
					if !(this_.self_.ParseJsDocCommentLine(revisedJsDocCommentLine)) && this_.shortcutMode {
						tryFlow_4 = 1
						return
					}
				}
				if !(j2gort.CheckNotNullString(revisedLine) == "") {
					if !(this_.self_.ParseLine(revisedLine)) && this_.shortcutMode {
						tryFlow_4 = 1
						return
					}
				}
			}()
			if tryFlow_4 == 1 {
				break
			}
		}
	}()
}

// Java method: JsFileLineParser.isCommentQuoted(java.lang.String line, int startOfMultilineComment, char quoteChar)
func (this_ *JsFileLineParserBase) priv_JsFileLineParser_isCommentQuoted_6098(line string, startOfMultilineComment int32, quoteChar uint16) bool {
	_ = line
	_ = startOfMultilineComment
	_ = quoteChar
	var startQuoteIndex int32 = j2gort.JStrIndexOfCp(line, int32(quoteChar))
	_ = startQuoteIndex
	for (startQuoteIndex != int32(-1)) && (startQuoteIndex < startOfMultilineComment) {
		var closingQuoteIndex int32 = j2gort.JStrIndexOfCpFrom(line, int32(quoteChar), (startQuoteIndex + int32(1)))
		_ = closingQuoteIndex
		if closingQuoteIndex == int32(-1) {
			return false
		}
		for this_.priv_JsFileLineParser_isEscaped_6099(line, closingQuoteIndex) {
			closingQuoteIndex = j2gort.JStrIndexOfCpFrom(line, int32(quoteChar), (closingQuoteIndex + int32(1)))
			if closingQuoteIndex == int32(-1) {
				return false
			}
		}
		if closingQuoteIndex > startOfMultilineComment {
			return true
		}
		startQuoteIndex = j2gort.JStrIndexOfCpFrom(line, int32(quoteChar), (closingQuoteIndex + int32(1)))
	}
	return false
}

// Java method: JsFileLineParser.isEscaped(java.lang.String line, int closingQuoteIndex)
func (this_ *JsFileLineParserBase) priv_JsFileLineParser_isEscaped_6099(line string, closingQuoteIndex int32) bool {
	_ = line
	_ = closingQuoteIndex
	var escapeCharCount int32 = int32(0)
	_ = escapeCharCount
	var index int32 = (closingQuoteIndex - int32(1))
	_ = index
	for int32(j2gort.JStrCharAt(line, index)) == int32(uint16(92)) {
		index = index - 1
		escapeCharCount = escapeCharCount + 1
	}
	return ((escapeCharCount % int32(2)) == int32(1))
}

// Java method: JsFileLineParser.parseLine(java.lang.String line)
func (this_ *JsFileLineParserBase) ParseLine(line string) bool {
	_ = line
	panic("abstract or native method invoked")
}

// Java method: JsFileLineParser.parseJsDocCommentLine(java.lang.String line)
func (this_ *JsFileLineParserBase) ParseJsDocCommentLine(line string) bool {
	_ = line
	return true
}

// Java method: JsFileLineParser.parseJsString(java.lang.String jsStringLiteral)
func (this_ *JsFileLineParserBase) ParseJsString(jsStringLiteral string) string {
	_ = jsStringLiteral
	this_.valueMatcher.Reset_1(j2gort.BoxString(jsStringLiteral))
	if !(this_.valueMatcher.Matches()) {
		panic(NewJsFileLineParser_ParseException("Syntax error in JS String literal", true))
	}
	if this_.valueMatcher.Group_1_int32(int32(1)) != j2gort.NullString {
		return this_.valueMatcher.Group_1_int32(int32(1))
	} else {
		return this_.valueMatcher.Group_1_int32(int32(2))
	}
}

// Java method: JsFileLineParser.parseJsStringArray(java.lang.String input)
func (this_ *JsFileLineParserBase) ParseJsStringArray(input string) j2gort.Ext_java_util_List {
	_ = input
	var results []string = []string{}
	_ = results
	var indexStart int32 = j2gort.JStrIndexOfCp(input, int32(uint16(91)))
	_ = indexStart
	var indexEnd int32 = j2gort.JStrLastIndexOfCp(input, int32(uint16(93)))
	_ = indexEnd
	if (indexStart == int32(-1)) || (indexEnd == int32(-1)) {
		panic(NewJsFileLineParser_ParseException("Syntax error when parsing JS array", true))
	}
	var innerValues string = j2gort.JStrSub(input, (indexStart + int32(1)), indexEnd)
	_ = innerValues
	if !(j2gort.CheckNotNullString(javalang.StringTrim(j2gort.CheckNotNullString(innerValues))) == "") {
		this_.valueMatcher.Reset_1(j2gort.BoxString(innerValues))
		for {
			if !(this_.valueMatcher.LookingAt()) {
				panic(NewJsFileLineParser_ParseException("Syntax error in JS String literal", true))
			}
			results = append(results, func() string {
				if this_.valueMatcher.Group_1_int32(int32(1)) != j2gort.NullString {
					return this_.valueMatcher.Group_1_int32(int32(1))
				}
				return this_.valueMatcher.Group_1_int32(int32(2))
			}())
			if this_.valueMatcher.HitEnd() {
				break
			}
			if int32(j2gort.JStrCharAt(innerValues, this_.valueMatcher.End_0())) != int32(uint16(44)) {
				panic(NewJsFileLineParser_ParseException("Missing comma in string array", true))
			}
			this_.valueMatcher.Region((this_.valueMatcher.End_0() + int32(1)), this_.valueMatcher.RegionEnd())
		}
	}
	return j2gort.WrapTypedArrayList(results)
}

// Java method: JsFileLineParser.parseJsStringMap(java.lang.String input)
func (this_ *JsFileLineParserBase) ParseJsStringMap(input string) j2gort.Ext_java_util_Map {
	_ = input
	input = j2gort.Ext_com_google_common_base_CharMatcher_S_Whitespace().TrimFrom(j2gort.BoxString(input))
	JsFileLineParser_Check(((!(j2gort.CheckNotNullString(input) == "") && (int32(j2gort.JStrCharAt(input, int32(0))) == int32(uint16(123)))) && (int32(j2gort.JStrCharAt(input, (j2gort.JStrLen(input)-int32(1)))) == int32(uint16(125)))), "Syntax error when parsing JS object")
	input = javalang.StringTrim(j2gort.CheckNotNullString(j2gort.JStrSub(input, int32(1), (j2gort.JStrLen(input) - int32(1)))))
	var results j2gort.Ext_java_util_Map = j2gort.NewExt_java_util_LinkedHashMap()
	_ = results
	var done bool = (j2gort.CheckNotNullString(input) == "")
	_ = done
	this_.valueMatcher.Reset_1(j2gort.BoxString(input))
	for !(done) {
		JsFileLineParser_Check(this_.valueMatcher.LookingAt(), "Bad key in JS object literal")
		var key string
		_ = key
		if this_.valueMatcher.Group_1_int32(int32(1)) != j2gort.NullString {
			key = this_.valueMatcher.Group_1_int32(int32(1))
		} else {
			key = this_.valueMatcher.Group_1_int32(int32(2))
		}
		JsFileLineParser_Check(!(this_.valueMatcher.HitEnd()), "Missing value in JS object literal")
		JsFileLineParser_Check((int32(j2gort.JStrCharAt(input, this_.valueMatcher.End_0())) == int32(uint16(58))), "Missing colon in JS object literal")
		this_.valueMatcher.Region((this_.valueMatcher.End_0() + int32(1)), this_.valueMatcher.RegionEnd())
		JsFileLineParser_Check(this_.valueMatcher.LookingAt(), "Bad value in JS object literal")
		var val string
		_ = val
		if this_.valueMatcher.Group_1_int32(int32(1)) != j2gort.NullString {
			val = this_.valueMatcher.Group_1_int32(int32(1))
		} else {
			val = this_.valueMatcher.Group_1_int32(int32(2))
		}
		j2gort.MapPutStr(results, key, j2gort.BoxString(val))
		if !(this_.valueMatcher.HitEnd()) {
			JsFileLineParser_Check((int32(j2gort.JStrCharAt(input, this_.valueMatcher.End_0())) == int32(uint16(44))), "Missing comma in JS object literal")
			this_.valueMatcher.Region((this_.valueMatcher.End_0() + int32(1)), this_.valueMatcher.RegionEnd())
		} else {
			done = true
		}
	}
	return results
}

// Java static method: JsFileLineParser.check(boolean condition, java.lang.String message)
func JsFileLineParser_Check(condition bool, message string) {
	clinit_JsFileLineParser()
	_ = condition
	_ = message
	if !(condition) {
		panic(NewJsFileLineParser_ParseException(message, true))
	}
}

// Java class: com.google.javascript.jscomp.deps.JsFileLineParser.ParseException
type JsFileLineParser_ParseException struct {
	javalang.JavaException
	fatal   bool `gson:"fatal"`
	idHash_ atomic.Int32
}

func (this_ *JsFileLineParser_ParseException) BindSelf_(s any) {
	_ = s
}
func (this_ *JsFileLineParser_ParseException) IsJavaNil_() bool { return this_ == nil }
func (*JsFileLineParser_ParseException) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.deps.JsFileLineParser$ParseException", "ParseException"
}
func (this_ *JsFileLineParser_ParseException) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.deps.JsFileLineParser.ParseException", "com.google.javascript.jscomp.deps.JsFileLineParser$ParseException", func(v_ any) bool {
		return j2gort.InstanceOf[*JsFileLineParser_ParseException](v_)
	})
}

func StaticGet_JsFileLineParser_ParseException_SerialVersionUID() int64 { return int64(1) }
func (this_ *JsFileLineParser_ParseException) initFields_JsFileLineParser_ParseException() {
	_ = this_
}

// Java constructor: ParseException(java.lang.String message, boolean fatal)
func NewJsFileLineParser_ParseException(message string, fatal bool) *JsFileLineParser_ParseException {
	this_ := &JsFileLineParser_ParseException{}
	InitJsFileLineParser_ParseException(this_, message, fatal)
	return this_
}

func InitJsFileLineParser_ParseException(this_ *JsFileLineParser_ParseException, message string, fatal bool) {
	_ = message
	_ = fatal
	j2gort.InitUserThrowable(this_, "com.google.javascript.jscomp.deps.JsFileLineParser.ParseException", j2gort.BoxString(message))
	this_.initFields_JsFileLineParser_ParseException()
	this_.fatal = fatal
}

// Java method: ParseException.isFatal()
func (this_ *JsFileLineParser_ParseException) IsFatal() bool {
	return this_.fatal
}
