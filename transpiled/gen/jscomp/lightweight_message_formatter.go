/*
 * Copyright 2007 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/LightweightMessageFormatter.java (generated mechanically by j2go).

package jscomp

import (
	"strconv"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.LightweightMessageFormatter
type LightweightMessageFormatter struct {
	AbstractMessageFormatterBase
	defaultFormat   *SourceExcerptProvider_SourceExcerpt `gson:"defaultFormat"`
	includeLocation bool                                 `gson:"includeLocation"`
	includeLevel    bool                                 `gson:"includeLevel"`
}

func (this_ *LightweightMessageFormatter) BindSelf_(s any) {
	_ = s
}
func (this_ *LightweightMessageFormatter) IsJavaNil_() bool { return this_ == nil }
func (*LightweightMessageFormatter) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.LightweightMessageFormatter", "LightweightMessageFormatter"
}
func (this_ *LightweightMessageFormatter) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.AbstractMessageFormatterBase.IdentityHash_()
}
func (*LightweightMessageFormatter) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.LightweightMessageFormatter", "com.google.javascript.jscomp.LightweightMessageFormatter", func(v_ any) bool {
		return j2gort.InstanceOf[*LightweightMessageFormatter](v_)
	})
}
func (*LightweightMessageFormatter) Implements_jscomp_AbstractMessageFormatter() {}
func (*LightweightMessageFormatter) Implements_jscomp_MessageFormatter()         {}

var static_LightweightMessageFormatter_excerptFormatter SourceExcerptProvider_ExcerptFormatter

func StaticGet_LightweightMessageFormatter_ExcerptFormatter() SourceExcerptProvider_ExcerptFormatter {
	clinit_LightweightMessageFormatter()
	return static_LightweightMessageFormatter_excerptFormatter
}
func StaticSet_LightweightMessageFormatter_ExcerptFormatter(v SourceExcerptProvider_ExcerptFormatter) {
	clinit_LightweightMessageFormatter()
	static_LightweightMessageFormatter_excerptFormatter = v
}
func StaticGet_LightweightMessageFormatter_MAX_MULTILINE_ERROR_LENGTH() int32 { return int32(4) }

var clinitDone_LightweightMessageFormatter uint32
var clinitState_LightweightMessageFormatter j2gort.ClinitState

func clinit_LightweightMessageFormatter() {
	if atomic.LoadUint32(&clinitDone_LightweightMessageFormatter) != 0 {
		return
	}
	clinitSlow_LightweightMessageFormatter()
}

func clinitSlow_LightweightMessageFormatter() {
	if !j2gort.EnterClinit(&clinitState_LightweightMessageFormatter, &clinitDone_LightweightMessageFormatter) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_LightweightMessageFormatter, &clinitDone_LightweightMessageFormatter)
	static_LightweightMessageFormatter_excerptFormatter = NewLightweightMessageFormatter_LineNumberingFormatter()
}

func (this_ *LightweightMessageFormatter) initFields_LightweightMessageFormatter() {
	_ = this_
	this_.includeLocation = true
	this_.includeLevel = true
}

// Java constructor: LightweightMessageFormatter()
func NewLightweightMessageFormatter() *LightweightMessageFormatter {
	clinit_LightweightMessageFormatter()
	this_ := &LightweightMessageFormatter{}
	InitLightweightMessageFormatter(this_)
	return this_
}

func InitLightweightMessageFormatter(this_ *LightweightMessageFormatter) {
	InitAbstractMessageFormatter(&this_.AbstractMessageFormatterBase, nil)
	this_.initFields_LightweightMessageFormatter()
	this_.defaultFormat = StaticGet_SourceExcerptProvider_SourceExcerpt_LINE()
}

// Java constructor: LightweightMessageFormatter(com.google.javascript.jscomp.SourceExcerptProvider source)
func NewLightweightMessageFormatter_1(source SourceExcerptProvider) *LightweightMessageFormatter {
	clinit_LightweightMessageFormatter()
	this_ := &LightweightMessageFormatter{}
	InitLightweightMessageFormatter_1(this_, source)
	return this_
}

func InitLightweightMessageFormatter_1(this_ *LightweightMessageFormatter, source SourceExcerptProvider) {
	_ = source
	InitLightweightMessageFormatter_2(this_, source, StaticGet_SourceExcerptProvider_SourceExcerpt_LINE())
}

// Java constructor: LightweightMessageFormatter(com.google.javascript.jscomp.SourceExcerptProvider source, com.google.javascript.jscomp.SourceExcerptProvider.SourceExcerpt defaultFormat)
func NewLightweightMessageFormatter_2(source SourceExcerptProvider, defaultFormat *SourceExcerptProvider_SourceExcerpt) *LightweightMessageFormatter {
	clinit_LightweightMessageFormatter()
	this_ := &LightweightMessageFormatter{}
	InitLightweightMessageFormatter_2(this_, source, defaultFormat)
	return this_
}

func InitLightweightMessageFormatter_2(this_ *LightweightMessageFormatter, source SourceExcerptProvider, defaultFormat *SourceExcerptProvider_SourceExcerpt) {
	_ = source
	_ = defaultFormat
	InitAbstractMessageFormatter(&this_.AbstractMessageFormatterBase, source)
	this_.initFields_LightweightMessageFormatter()
	j2gort.CheckNotNullIface(source)
	this_.defaultFormat = defaultFormat
}

// Java static method: LightweightMessageFormatter.withoutSource()
func LightweightMessageFormatter_WithoutSource() *LightweightMessageFormatter {
	clinit_LightweightMessageFormatter()
	return NewLightweightMessageFormatter()
}

// Java method: LightweightMessageFormatter.setIncludeLocation(boolean includeLocation)
func (this_ *LightweightMessageFormatter) SetIncludeLocation(includeLocation bool) *LightweightMessageFormatter {
	_ = includeLocation
	this_.includeLocation = includeLocation
	return this_
}

// Java method: LightweightMessageFormatter.setIncludeLevel(boolean includeLevel)
func (this_ *LightweightMessageFormatter) SetIncludeLevel(includeLevel bool) *LightweightMessageFormatter {
	_ = includeLevel
	this_.includeLevel = includeLevel
	return this_
}

// Java method: LightweightMessageFormatter.formatError(com.google.javascript.jscomp.JSError error)
func (this_ *LightweightMessageFormatter) FormatError(error_ *JSError) string {
	_ = error_
	return this_.priv_LightweightMessageFormatter_format_3723(error_, false)
}

// Java method: LightweightMessageFormatter.formatWarning(com.google.javascript.jscomp.JSError warning)
func (this_ *LightweightMessageFormatter) FormatWarning(warning *JSError) string {
	_ = warning
	return this_.priv_LightweightMessageFormatter_format_3723(warning, true)
}

// Java method: LightweightMessageFormatter.format(com.google.javascript.jscomp.JSError error, boolean warning)
func (this_ *LightweightMessageFormatter) priv_LightweightMessageFormatter_format_3723(error_ *JSError, warning bool) string {
	_ = error_
	_ = warning
	var source SourceExcerptProvider = this_.GetSource()
	_ = source
	var sourceName string = error_.SourceName()
	_ = sourceName
	var lineNumber int32 = error_.GetLineNumber()
	_ = lineNumber
	var charno int32 = error_.Charno()
	_ = charno
	var b *j2gort.Ext_java_lang_StringBuilder = j2gort.NewExt_java_lang_StringBuilder()
	_ = b
	var boldLine *j2gort.Ext_java_lang_StringBuilder = j2gort.NewExt_java_lang_StringBuilder()
	_ = boldLine
	var mapping *j2gort.Ext_com_google_debugging_sourcemap_proto_Mapping_OriginalMapping
	_ = mapping
	if source == nil || source.IsJavaNil_() {
		mapping = nil
	} else {
		mapping = source.GetSourceMapping(error_.SourceName(), error_.GetLineNumber(), error_.Charno())
	}
	if this_.includeLocation {
		if mapping != nil {
			LightweightMessageFormatter_AppendPosition(b, sourceName, lineNumber, charno)
			sourceName = mapping.GetOriginalFile()
			lineNumber = mapping.GetLineNumber()
			charno = mapping.GetColumnPosition()
			b.Append_1_string("\nOriginally at:\n")
		}
		LightweightMessageFormatter_AppendPosition(boldLine, sourceName, lineNumber, charno)
	}
	if this_.includeLevel {
		boldLine.Append_1_string(this_.GetLevelName(func() *CheckLevel {
			if warning {
				return StaticGet_CheckLevel_WARNING()
			}
			return StaticGet_CheckLevel_ERROR()
		}()))
		boldLine.Append_1_string(" - [")
		boldLine.Append_1_string(error_.Type().key)
		boldLine.Append_1_string("] ")
	}
	boldLine.Append_1_string(error_.Description())
	b.Append_1_string(this_.MaybeEmbolden(boldLine.String()))
	b.Append_1_uint16(uint16(10))
	var sourceExcerptWithPosition string = this_.priv_LightweightMessageFormatter_getExcerptWithPosition_3725(error_, sourceName, lineNumber, charno, func() *SourceExcerptProvider_SourceExcerpt {
		if mapping != nil {
			return StaticGet_SourceExcerptProvider_SourceExcerpt_LINE()
		}
		return this_.defaultFormat
	}())
	_ = sourceExcerptWithPosition
	if sourceExcerptWithPosition != j2gort.NullString {
		b.Append_1_string(sourceExcerptWithPosition)
	}
	return b.String()
}

// Java method: LightweightMessageFormatter.getExcerptWithPosition(com.google.javascript.jscomp.JSError error)
func (this_ *LightweightMessageFormatter) GetExcerptWithPosition_1(error_ *JSError) string {
	_ = error_
	return this_.priv_LightweightMessageFormatter_getExcerptWithPosition_3726(error_, error_.SourceName(), error_.GetLineNumber(), error_.Charno(), error_.Length(), this_.defaultFormat)
}

// Java method: LightweightMessageFormatter.getExcerptWithPosition(com.google.javascript.jscomp.JSError error, java.lang.String sourceName, int lineNumber, int charno, com.google.javascript.jscomp.SourceExcerptProvider.SourceExcerpt format)
func (this_ *LightweightMessageFormatter) priv_LightweightMessageFormatter_getExcerptWithPosition_3725(error_ *JSError, sourceName string, lineNumber int32, charno int32, format *SourceExcerptProvider_SourceExcerpt) string {
	_ = error_
	_ = sourceName
	_ = lineNumber
	_ = charno
	_ = format
	var nodeLength int32 = error_.Length()
	_ = nodeLength
	return this_.priv_LightweightMessageFormatter_getExcerptWithPosition_3726(error_, sourceName, lineNumber, charno, nodeLength, format)
}

// Java method: LightweightMessageFormatter.getExcerptWithPosition(com.google.javascript.jscomp.JSError error, java.lang.String sourceName, int lineNumber, int charno, int length, com.google.javascript.jscomp.SourceExcerptProvider.SourceExcerpt format)
func (this_ *LightweightMessageFormatter) priv_LightweightMessageFormatter_getExcerptWithPosition_3726(error_ *JSError, sourceName string, lineNumber int32, charno int32, length int32, format *SourceExcerptProvider_SourceExcerpt) string {
	_ = error_
	_ = sourceName
	_ = lineNumber
	_ = charno
	_ = length
	_ = format
	var b *j2gort.Ext_java_lang_StringBuilder = j2gort.NewExt_java_lang_StringBuilder()
	_ = b
	var source SourceExcerptProvider = this_.GetSource()
	_ = source
	var cookedLength int32
	_ = cookedLength
	if (charno >= int32(0)) && (length >= int32(0)) {
		cookedLength = (charno + length)
	} else {
		cookedLength = int32(-1)
	}
	var sourceExcerpt string
	_ = sourceExcerpt
	if source == nil || source.IsJavaNil_() {
		sourceExcerpt = j2gort.NullString
	} else {
		sourceExcerpt = format.Get_5(source, sourceName, lineNumber, cookedLength, StaticGet_LightweightMessageFormatter_ExcerptFormatter())
	}
	if sourceExcerpt == j2gort.NullString {
		return ""
	}
	var success bool = this_.priv_LightweightMessageFormatter_padSourceExcerpt_3727(error_, sourceExcerpt, charno, length, format, b)
	_ = success
	if success {
		return b.String()
	}
	var message string = "Source excerpt could not be formatted. This may indicate a bug in the compiler."
	_ = message
	return message
}

// Java method: LightweightMessageFormatter.padSourceExcerpt(com.google.javascript.jscomp.JSError error, java.lang.String sourceExcerpt, int charno, int length, com.google.javascript.jscomp.SourceExcerptProvider.SourceExcerpt format, java.lang.StringBuilder b)
func (this_ *LightweightMessageFormatter) priv_LightweightMessageFormatter_padSourceExcerpt_3727(error_ *JSError, sourceExcerpt string, charno int32, length int32, format *SourceExcerptProvider_SourceExcerpt, b *j2gort.Ext_java_lang_StringBuilder) bool {
	_ = error_
	_ = sourceExcerpt
	_ = charno
	_ = length
	_ = format
	_ = b
	if j2gort.RefEquals(format, StaticGet_SourceExcerptProvider_SourceExcerpt_FULL()) {
		if charno >= int32(0) {
			var success bool = this_.priv_LightweightMessageFormatter_padMultipleLines_3729(charno, sourceExcerpt, b, error_.Node())
			_ = success
			if !(success) {
				return false
			}
		} else {
			b.Append_1_string(sourceExcerpt)
			b.Append_1_uint16(uint16(10))
		}
	} else {
		b.Append_1_string(sourceExcerpt)
		b.Append_1_uint16(uint16(10))
		if (j2gort.RefEquals(format, StaticGet_SourceExcerptProvider_SourceExcerpt_LINE()) && (int32(0) <= charno)) && (charno <= j2gort.JStrLen(sourceExcerpt)) {
			var success_2 bool = this_.priv_LightweightMessageFormatter_padLine_3728(charno, sourceExcerpt, b, length, error_.Node())
			_ = success_2
			if !(success_2) {
				return false
			}
		}
	}
	return true
}

// Java static method: LightweightMessageFormatter.appendPosition(java.lang.StringBuilder b, java.lang.String sourceName, int lineNumber, int charno)
func LightweightMessageFormatter_AppendPosition(b *j2gort.Ext_java_lang_StringBuilder, sourceName string, lineNumber int32, charno int32) {
	clinit_LightweightMessageFormatter()
	_ = b
	_ = sourceName
	_ = lineNumber
	_ = charno
	if sourceName != j2gort.NullString {
		b.Append_1_string(sourceName)
		if lineNumber > int32(0) {
			j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](b.Append_1_uint16(uint16(58))).Append_1_int32(lineNumber)
			if charno >= int32(0) {
				j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](b.Append_1_uint16(uint16(58))).Append_1_int32(charno)
			}
		}
		b.Append_1_string(": ")
	}
}

// Java method: LightweightMessageFormatter.padLine(int charno, java.lang.String sourceExcerpt, java.lang.StringBuilder b, int errLength, com.google.javascript.rhino.Node errorNode)
func (this_ *LightweightMessageFormatter) priv_LightweightMessageFormatter_padLine_3728(charno int32, sourceExcerpt string, b *j2gort.Ext_java_lang_StringBuilder, errLength int32, errorNode pkg_rhino.Node) bool {
	_ = charno
	_ = sourceExcerpt
	_ = b
	_ = errLength
	_ = errorNode
	if charno > j2gort.JStrLen(sourceExcerpt) {
		return false
	}
	{
		var i int32 = int32(0)
		_ = i
		for ; i < charno; i = i + 1 {
			var c uint16 = j2gort.JStrCharAt(sourceExcerpt, i)
			_ = c
			if pkg_rhino.TokenUtil_IsWhitespace(int32(c)) {
				b.Append_1_uint16(c)
			} else {
				b.Append_1_uint16(uint16(32))
			}
		}
	}
	if errorNode == 0 {
		b.Append_1_string("^")
	} else {
		var length int32 = max(int32(1), min(errLength, (j2gort.JStrLen(sourceExcerpt)-charno)))
		_ = length
		{
			var i_2 int32 = int32(0)
			_ = i_2
			for ; i_2 < length; i_2 = i_2 + 1 {
				b.Append_1_string("^")
			}
		}
	}
	b.Append_1_string("\n")
	return true
}

// Java method: LightweightMessageFormatter.padMultipleLines(int startCharno, java.lang.String sourceExcerpt, java.lang.StringBuilder b, com.google.javascript.rhino.Node errorNode)
func (this_ *LightweightMessageFormatter) priv_LightweightMessageFormatter_padMultipleLines_3729(startCharno int32, sourceExcerpt string, b *j2gort.Ext_java_lang_StringBuilder, errorNode pkg_rhino.Node) bool {
	_ = startCharno
	_ = sourceExcerpt
	_ = b
	_ = errorNode
	if errorNode == 0 {
		b.Append_1_string(sourceExcerpt)
		b.Append_1_string("\n")
		var charWithLineNumberOffset int32 = ((startCharno + j2gort.JStrIndexOfCp(sourceExcerpt, int32(uint16(124)))) + int32(2))
		_ = charWithLineNumberOffset
		return this_.priv_LightweightMessageFormatter_padLine_3728(charWithLineNumberOffset, sourceExcerpt, b, int32(-1), pkg_rhino.Node(0))
	}
	var lines j2gort.Ext_java_util_List = j2gort.Ext_com_google_common_base_Splitter_On_1_uint16(uint16(10)).SplitToList(j2gort.BoxString(sourceExcerpt))
	_ = lines
	var requiresTruncation bool = (lines.Size() > int32(4))
	_ = requiresTruncation
	var truncationStart int32 = (int32(4) / int32(2))
	_ = truncationStart
	var truncationEnd int32 = (lines.Size() - (int32(4) / int32(2)))
	_ = truncationEnd
	var remainingLength int32 = errorNode.GetLength()
	_ = remainingLength
	var charno int32 = startCharno
	_ = charno
	{
		var i int32 = int32(0)
		_ = i
		for ; i < lines.Size(); i = i + 1 {
			var line string = j2gort.ListGetStr(lines, i)
			_ = line
			if requiresTruncation && (i == truncationStart) {
				b.Append_1_string("...\n")
			}
			var shouldPrintLine bool = ((!(requiresTruncation) || (i < truncationStart)) || (i >= truncationEnd))
			_ = shouldPrintLine
			var charWithLineNumberOffset_2 int32 = ((charno + j2gort.JStrIndexOfCp(line, int32(uint16(124)))) + int32(2))
			_ = charWithLineNumberOffset_2
			if shouldPrintLine {
				b.Append_1_string(line)
				b.Append_1_string("\n")
				var success bool = this_.priv_LightweightMessageFormatter_padLine_3728(charWithLineNumberOffset_2, line, b, remainingLength, errorNode)
				_ = success
				if !(success) {
					return false
				}
			}
			remainingLength = (remainingLength - ((j2gort.JStrLen(line) + int32(1)) - charWithLineNumberOffset_2))
			charno = int32(0)
		}
	}
	return true
}

// Java class: com.google.javascript.jscomp.LightweightMessageFormatter.LineNumberingFormatter
type LightweightMessageFormatter_LineNumberingFormatter struct {
	idHash_ atomic.Int32
}

func (this_ *LightweightMessageFormatter_LineNumberingFormatter) BindSelf_(s any) {
	_ = s
}
func (this_ *LightweightMessageFormatter_LineNumberingFormatter) IsJavaNil_() bool {
	return this_ == nil
}
func (*LightweightMessageFormatter_LineNumberingFormatter) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.LightweightMessageFormatter$LineNumberingFormatter", "LineNumberingFormatter"
}
func (this_ *LightweightMessageFormatter_LineNumberingFormatter) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*LightweightMessageFormatter_LineNumberingFormatter) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.LightweightMessageFormatter.LineNumberingFormatter", "com.google.javascript.jscomp.LightweightMessageFormatter$LineNumberingFormatter", func(v_ any) bool {
		return j2gort.InstanceOf[*LightweightMessageFormatter_LineNumberingFormatter](v_)
	})
}
func (*LightweightMessageFormatter_LineNumberingFormatter) Implements_jscomp_SourceExcerptProvider_ExcerptFormatter() {
}

func (this_ *LightweightMessageFormatter_LineNumberingFormatter) initFields_LightweightMessageFormatter_LineNumberingFormatter() {
	_ = this_
}

// Java constructor: LineNumberingFormatter()
func NewLightweightMessageFormatter_LineNumberingFormatter() *LightweightMessageFormatter_LineNumberingFormatter {
	this_ := &LightweightMessageFormatter_LineNumberingFormatter{}
	InitLightweightMessageFormatter_LineNumberingFormatter(this_)
	return this_
}

func InitLightweightMessageFormatter_LineNumberingFormatter(this_ *LightweightMessageFormatter_LineNumberingFormatter) {
	this_.initFields_LightweightMessageFormatter_LineNumberingFormatter()
}

// Java method: LineNumberingFormatter.formatLine(java.lang.String line, int lineNumber)
func (this_ *LightweightMessageFormatter_LineNumberingFormatter) FormatLine(line string, lineNumber int32) string {
	_ = line
	_ = lineNumber
	return line
}

// Java method: LineNumberingFormatter.formatRegion(com.google.javascript.jscomp.@org.jspecify.annotations.Nullable Region region)
func (this_ *LightweightMessageFormatter_LineNumberingFormatter) FormatRegion(region Region) string {
	_ = region
	if region == nil || region.IsJavaNil_() {
		return j2gort.NullString
	}
	var code string = region.GetSourceExcerpt()
	_ = code
	if j2gort.CheckNotNullString(code) == "" {
		return j2gort.NullString
	}
	var numberLength int32 = j2gort.JStrLen(strconv.FormatInt(int64(region.GetEndingLineNumber()), 10))
	_ = numberLength
	var builder *j2gort.Ext_java_lang_StringBuilder = j2gort.NewExt_java_lang_StringBuilder_1_int32((j2gort.JStrLen(code) * int32(2)))
	_ = builder
	var start int32 = int32(0)
	_ = start
	var end int32 = j2gort.JStrIndexOfCpFrom(code, int32(uint16(10)), start)
	_ = end
	var lineNumber int32 = region.GetBeginningLineNumber()
	_ = lineNumber
	for start >= int32(0) {
		var line string = j2gort.NullString
		_ = line
		if end < int32(0) {
			line = j2gort.JStrSubFrom(code, start)
			if j2gort.CheckNotNullString(line) == "" {
				return builder.Substring_2(int32(0), (builder.Length() - int32(1)))
			}
		} else {
			line = j2gort.JStrSub(code, start, end)
		}
		builder.Append_1_string("  ")
		var spaces int32 = (numberLength - j2gort.JStrLen(strconv.FormatInt(int64(lineNumber), 10)))
		_ = spaces
		builder.Repeat_2_any_int32(" ", spaces)
		builder.Append_1_int32(lineNumber)
		builder.Append_1_string("| ")
		if end < int32(0) {
			builder.Append_1_string(line)
			start = int32(-1)
		} else {
			builder.Append_1_string(line)
			builder.Append_1_uint16(uint16(10))
			start = (end + int32(1))
			end = j2gort.JStrIndexOfCpFrom(code, int32(uint16(10)), start)
			lineNumber = lineNumber + 1
		}
	}
	return builder.String()
}
