/*
 * Copyright 2004 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/CodePrinter.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	pkg_base "github.com/paullepoulpe/jsgo/transpiled/gen/base"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	pkg_sourcemap "github.com/paullepoulpe/jsgo/transpiled/gen/sourcemap"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javaio"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.CodePrinter
type CodePrinter struct {
	idHash_ atomic.Int32
}

func (this_ *CodePrinter) BindSelf_(s any) {
	_ = s
}
func (this_ *CodePrinter) IsJavaNil_() bool { return this_ == nil }
func (*CodePrinter) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CodePrinter", "CodePrinter"
}
func (this_ *CodePrinter) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CodePrinter) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CodePrinter", "com.google.javascript.jscomp.CodePrinter", func(v_ any) bool {
		return j2gort.InstanceOf[*CodePrinter](v_)
	})
}

func (this_ *CodePrinter) initFields_CodePrinter() {
	_ = this_
}

// Java constructor: CodePrinter()
func NewCodePrinter() *CodePrinter {
	this_ := &CodePrinter{}
	InitCodePrinter(this_)
	return this_
}

func InitCodePrinter(this_ *CodePrinter) {
	this_.initFields_CodePrinter()
}

// Java static method: CodePrinter.toSource(com.google.javascript.rhino.Node root, com.google.javascript.jscomp.CodePrinter.Format outputFormat, com.google.javascript.jscomp.CompilerOptions options, com.google.javascript.jscomp.CodePrinter.@org.jspecify.annotations.Nullable LicenseTracker licenseTracker, boolean tagAsTypeSummary, boolean tagAsStrict, boolean lineBreak, com.google.javascript.jscomp.CodePrinter.Builder.CodeGeneratorFactory codeGeneratorFactory)
func CodePrinter_ToSource(root pkg_rhino.Node, outputFormat CodePrinter_Format, options *CompilerOptions, licenseTracker CodePrinter_LicenseTracker, tagAsTypeSummary bool, tagAsStrict bool, lineBreak bool, codeGeneratorFactory CodePrinter_Builder_CodeGeneratorFactory) *CodePrinter_SourceAndMappings {
	_ = root
	_ = outputFormat
	_ = options
	_ = licenseTracker
	_ = tagAsTypeSummary
	_ = tagAsStrict
	_ = lineBreak
	_ = codeGeneratorFactory
	if !(options.GetSourceMapDetailLevel() != nil) {
		j2gort.CheckState(false)
	}
	var mcp CodePrinter_MappedCodePrinter
	_ = mcp
	if outputFormat == StaticGet_CodePrinter_Format_COMPACT() {
		mcp = NewCodePrinter_CompactCodePrinter(lineBreak, options.GetLineLengthThreshold(), options.ShouldGatherSourceMapInfo(), options.GetSourceMapDetailLevel(), licenseTracker)
	} else {
		mcp = NewCodePrinter_PrettyCodePrinter(options.GetLineLengthThreshold(), options.ShouldGatherSourceMapInfo(), options.GetSourceMapDetailLevel(), licenseTracker)
	}
	var cg CodeGenerator = codeGeneratorFactory.GetCodeGenerator(outputFormat, j2gort.AsIface[CodeConsumer](mcp))
	_ = cg
	if tagAsTypeSummary {
		cg.TagAsTypeSummary()
	}
	if tagAsStrict {
		cg.TagAsStrict()
	}
	cg.Add_1_Node(root)
	mcp.EndFile()
	var val_1 CodePrinter_SourceAndMappings
	InitCodePrinter_SourceAndMappings(&val_1)
	var result *CodePrinter_SourceAndMappings = &val_1
	_ = result
	{
		r_2 := result
		r_2.source = mcp.GetCode()
	}
	if options.ShouldGatherSourceMapInfo() {
		{
			r_3 := result
			r_3.mappings = mcp.GetSourceMappings(result.source)
		}
	}
	return result
}

// Java polymorphic class interface: com.google.javascript.jscomp.CodePrinter.MappedCodePrinter
type CodePrinter_MappedCodePrinter interface {
	Implements_jscomp_CodePrinter_MappedCodePrinter()
	IsJavaNil_() bool
	GetCodePrinter_MappedCodePrinterBase() *CodePrinter_MappedCodePrinterBase
	CodeConsumer
	GetSourceMappings(string) j2gort.Ext_java_util_List
	ReportLineCut(int32, int32)
	GetCode() string
	GetCurrentCharIndex() int32
	GetCurrentLineIndex() int32
}

func ToIface_CodePrinter_MappedCodePrinter[T interface {
	CodePrinter_MappedCodePrinter
	comparable
}](v T) CodePrinter_MappedCodePrinter {
	var z T
	if v == z {
		return nil
	}
	return v
}

// Java class: com.google.javascript.jscomp.CodePrinter.MappedCodePrinter
type CodePrinter_MappedCodePrinterBase struct {
	CodeConsumerBase
	mappings                    j2gort.Deque[*SourceMap_Mapping]    `gson:"mappings"`
	allMappings                 []*SourceMap_Mapping                `gson:"allMappings"`
	completeMappings            []*SourceMap_Mapping                `gson:"completeMappings"`
	sourceMapDetailLevel        *SourceMap_DetailLevel              `gson:"sourceMapDetailLevel"`
	licenseTracker              CodePrinter_LicenseTracker          `gson:"licenseTracker"`
	code                        *j2gort.Ext_java_lang_StringBuilder `gson:"code"`
	firstCandidateMappingForCut int32                               `gson:"firstCandidateMappingForCut"`
	lineLengthThreshold         int32                               `gson:"lineLengthThreshold"`
	lineLength                  int32                               `gson:"lineLength"`
	lineIndex                   int32                               `gson:"lineIndex"`
	createSrcMap                bool                                `gson:"createSrcMap"`
}

func (this_ *CodePrinter_MappedCodePrinterBase) BindSelf_(s any) {
	_ = s
	this_.CodeConsumerBase.BindSelf_(s)
}
func (this_ *CodePrinter_MappedCodePrinterBase) IsJavaNil_() bool { return this_ == nil }
func (*CodePrinter_MappedCodePrinterBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CodePrinter$MappedCodePrinter", "MappedCodePrinter"
}
func (this_ *CodePrinter_MappedCodePrinterBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.CodeConsumerBase.IdentityHash_()
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CodePrinter.MappedCodePrinter", "com.google.javascript.jscomp.CodePrinter$MappedCodePrinter", func(v_ any) bool {
		return j2gort.InstanceOf[CodePrinter_MappedCodePrinter](v_)
	})
}
func (b *CodePrinter_MappedCodePrinterBase) GetCodePrinter_MappedCodePrinterBase() *CodePrinter_MappedCodePrinterBase {
	return b
}
func (*CodePrinter_MappedCodePrinterBase) Implements_jscomp_CodePrinter_MappedCodePrinter() {}
func (*CodePrinter_MappedCodePrinterBase) Implements_jscomp_CodeConsumer()                  {}

func (this_ *CodePrinter_MappedCodePrinterBase) initFields_CodePrinter_MappedCodePrinter() {
	_ = this_
	this_.mappings = j2gort.Deque[*SourceMap_Mapping]{}
	this_.allMappings = []*SourceMap_Mapping{}
	this_.completeMappings = []*SourceMap_Mapping{}

	this_.code = j2gort.NewExt_java_lang_StringBuilder_1_int32(int32(1024))

}

func InitCodePrinter_MappedCodePrinter(this_ *CodePrinter_MappedCodePrinterBase, lineLengthThreshold int32, createSrcMap bool, sourceMapDetailLevel *SourceMap_DetailLevel, licenseTracker CodePrinter_LicenseTracker) {
	_ = lineLengthThreshold
	_ = createSrcMap
	_ = sourceMapDetailLevel
	_ = licenseTracker
	InitCodeConsumer(&this_.CodeConsumerBase)
	this_.initFields_CodePrinter_MappedCodePrinter()
	if !(sourceMapDetailLevel != nil) {
		j2gort.CheckState(false)
	}
	if lineLengthThreshold <= int32(0) {
		this_.lineLengthThreshold = int32(2147483647)
	} else {
		this_.lineLengthThreshold = lineLengthThreshold
	}
	this_.createSrcMap = createSrcMap
	this_.sourceMapDetailLevel = sourceMapDetailLevel
	this_.licenseTracker = licenseTracker
}

// Java method: MappedCodePrinter.append(java.lang.String str)
func (this_ *CodePrinter_MappedCodePrinterBase) Append(str string) {
	_ = str
	this_.code.Append_1_string(str)
	this_.lineLength = (this_.lineLength + j2gort.JStrLen(str))
}

// Java method: MappedCodePrinter.trackLicenses(com.google.javascript.rhino.Node node)
func (this_ *CodePrinter_MappedCodePrinterBase) TrackLicenses(node pkg_rhino.Node) {
	_ = node
	if this_.licenseTracker != nil && !this_.licenseTracker.IsJavaNil_() {
		this_.licenseTracker.TrackLicensesForNode(node)
	}
}

// Java method: MappedCodePrinter.startSourceMapping(com.google.javascript.rhino.Node node)
func (this_ *CodePrinter_MappedCodePrinterBase) StartSourceMapping(node pkg_rhino.Node) {
	_ = node
	if !(this_.sourceMapDetailLevel != nil) {
		j2gort.CheckState(false)
	}
	if !(node != 0) {
		j2gort.CheckState(false)
	}
	if ((this_.createSrcMap && (node.GetSourceFileName() != j2gort.NullString)) && (node.GetLineno() > int32(0))) && j2gort.PredTestNode(this_.sourceMapDetailLevel, uint32(node)) {
		var line int32 = this_.GetCurrentLineIndex()
		_ = line
		var index int32 = this_.GetCurrentCharIndex()
		_ = index
		if !(line >= int32(0)) {
			j2gort.CheckState(false)
		}
		var val_4 SourceMap_Mapping
		InitSourceMap_Mapping(&val_4)
		var mapping *SourceMap_Mapping = &val_4
		_ = mapping
		{
			r_5 := mapping
			r_5.node = node
		}
		{
			r_6 := mapping
			r_6.start = pkg_sourcemap.NewFilePosition(line, index)
		}
		this_.mappings.Push(mapping)
		this_.allMappings = append(this_.allMappings, mapping)
	}
}

// Java method: MappedCodePrinter.endSourceMapping(com.google.javascript.rhino.Node node)
func (this_ *CodePrinter_MappedCodePrinterBase) EndSourceMapping(node pkg_rhino.Node) {
	_ = node
	if (this_.createSrcMap && !(this_.mappings.IsEmpty())) && (this_.mappings.Peek().node == node) {
		var mapping *SourceMap_Mapping = this_.mappings.Pop()
		_ = mapping
		var line int32 = this_.GetCurrentLineIndex()
		_ = line
		var index int32 = this_.GetCurrentCharIndex()
		_ = index
		if !(line >= int32(0)) {
			j2gort.CheckState(false)
		}
		{
			r_7 := mapping
			r_7.end = pkg_sourcemap.NewFilePosition(line, index)
		}
		this_.completeMappings = append(this_.completeMappings, mapping)
	}
}

// Java method: MappedCodePrinter.getSourceMappings(java.lang.String code)
func (this_ *CodePrinter_MappedCodePrinterBase) GetSourceMappings(code string) j2gort.Ext_java_util_List {
	_ = code
	if !(this_.createSrcMap) {
		return nil
	}
	var lineLengths j2gort.Ext_com_google_common_collect_ImmutableList = CodePrinter_MappedCodePrinter_ComputeLineLengths(code)
	_ = lineLengths
	var fixedMappings []*SourceMap_Mapping = j2gort.MakeNonNilSlice[*SourceMap_Mapping](int32(len(this_.allMappings)))
	_ = fixedMappings
	{
		var i int32 = int32(0)
		_ = i
		var size int32 = int32(len(this_.allMappings))
		_ = size
		for ; i < size; i = i + 1 {
			var mapping *SourceMap_Mapping = this_.allMappings[i]
			_ = mapping
			var adjusted *SourceMap_Mapping = NewSourceMap_Mapping()
			_ = adjusted
			{
				r_8 := adjusted
				r_8.node = mapping.node
			}
			{
				r_9 := adjusted
				r_9.start = mapping.start
			}
			{
				r_10 := adjusted
				r_10.end = CodePrinter_MappedCodePrinter_AdjustEndPosition(j2gort.AsIface[j2gort.Ext_java_util_List](lineLengths), mapping.end)
			}
			fixedMappings = append(fixedMappings, adjusted)
		}
	}
	return j2gort.WrapTypedArrayList(fixedMappings)
}

// Java method: MappedCodePrinter.reportLineCut(int lineIndex, int charIndex)
func (this_ *CodePrinter_MappedCodePrinterBase) ReportLineCut(lineIndex int32, charIndex int32) {
	_ = lineIndex
	_ = charIndex
	if this_.createSrcMap {
		var mappingCount int32 = int32(len(this_.allMappings))
		_ = mappingCount
		{
			var i int32 = this_.firstCandidateMappingForCut
			_ = i
			for ; i < mappingCount; i = i + 1 {
				var mapping *SourceMap_Mapping = this_.allMappings[i]
				_ = mapping
				{
					r_11 := mapping
					r_11.start = CodePrinter_MappedCodePrinter_ConvertPositionAfterLineCut(mapping.start, lineIndex, charIndex)
				}
			}
		}
		this_.firstCandidateMappingForCut = mappingCount
		{
			var i_2 int32 = int32(0)
			_ = i_2
			var size int32 = int32(len(this_.completeMappings))
			_ = size
			for ; i_2 < size; i_2 = i_2 + 1 {
				var mapping_2 *SourceMap_Mapping = this_.completeMappings[i_2]
				_ = mapping_2
				{
					r_12 := mapping_2
					r_12.end = CodePrinter_MappedCodePrinter_ConvertPositionAfterLineCut(mapping_2.end, lineIndex, charIndex)
				}
			}
		}
		clear(this_.completeMappings)
		this_.completeMappings = this_.completeMappings[:0]
	}
}

// Java static method: MappedCodePrinter.convertPositionAfterLineCut(com.google.debugging.sourcemap.FilePosition position, int lineIndex, int characterPosition)
func CodePrinter_MappedCodePrinter_ConvertPositionAfterLineCut(position *pkg_sourcemap.FilePosition, lineIndex int32, characterPosition int32) *pkg_sourcemap.FilePosition {
	_ = position
	_ = lineIndex
	_ = characterPosition
	var originalLine int32 = position.GetLine()
	_ = originalLine
	var originalChar int32 = position.GetColumn()
	_ = originalChar
	if (originalLine == lineIndex) && (originalChar >= characterPosition) {
		return pkg_sourcemap.NewFilePosition((originalLine + int32(1)), (originalChar - characterPosition))
	} else {
		return position
	}
}

// Java method: MappedCodePrinter.getCode()
func (this_ *CodePrinter_MappedCodePrinterBase) GetCode() string {
	return this_.code.String()
}

// Java method: MappedCodePrinter.getLastChar()
func (this_ *CodePrinter_MappedCodePrinterBase) GetLastChar() uint16 {
	if this_.code.Length() > int32(0) {
		return this_.code.CharAt((this_.code.Length() - int32(1)))
	} else {
		return uint16(0)
	}
}

// Java method: MappedCodePrinter.getCurrentCharIndex()
func (this_ *CodePrinter_MappedCodePrinterBase) GetCurrentCharIndex() int32 {
	return this_.lineLength
}

// Java method: MappedCodePrinter.getCurrentLineIndex()
func (this_ *CodePrinter_MappedCodePrinterBase) GetCurrentLineIndex() int32 {
	return this_.lineIndex
}

// Java static method: MappedCodePrinter.computeLineLengths(java.lang.String code)
func CodePrinter_MappedCodePrinter_ComputeLineLengths(code string) j2gort.Ext_com_google_common_collect_ImmutableList {
	_ = code
	var builder []int32 = []int32{}
	_ = builder
	var lineStartPos int32 = int32(0)
	_ = lineStartPos
	var lineEndPos int32 = j2gort.JStrIndexOfCp(code, int32(uint16(10)))
	_ = lineEndPos
	for lineEndPos > int32(-1) {
		builder = append(builder, lineEndPos-lineStartPos)
		lineStartPos = (lineEndPos + int32(1))
		lineEndPos = j2gort.JStrIndexOfCpFrom(code, int32(uint16(10)), lineStartPos)
	}
	return j2gort.ImmListFromTypedSlice(builder)
}

// Java static method: MappedCodePrinter.adjustEndPosition(java.util.List<java.lang.Integer> lineLengths, com.google.debugging.sourcemap.FilePosition endPosition)
func CodePrinter_MappedCodePrinter_AdjustEndPosition(lineLengths j2gort.Ext_java_util_List, endPosition *pkg_sourcemap.FilePosition) *pkg_sourcemap.FilePosition {
	_ = lineLengths
	_ = endPosition
	var line int32 = endPosition.GetLine()
	_ = line
	if line >= lineLengths.Size() {
		return endPosition
	}
	j2gort.CheckStateFmt((endPosition.GetColumn() <= j2gort.UnboxInt(j2gort.CastPtr[j2gort.BoxInt](lineLengths.Get(line)))), "End position %s points to a column larger than line length %s", j2gort.PtrToAny(endPosition), j2gort.PtrToAny(j2gort.CastPtr[j2gort.BoxInt](lineLengths.Get(line))))
	if endPosition.GetColumn() == j2gort.UnboxInt(j2gort.CastPtr[j2gort.BoxInt](lineLengths.Get(line))) {
		return pkg_sourcemap.NewFilePosition((line + int32(1)), int32(0))
	}
	return endPosition
}

// Java class: com.google.javascript.jscomp.CodePrinter.PrettyCodePrinter
type CodePrinter_PrettyCodePrinter struct {
	CodePrinter_MappedCodePrinterBase
	indent int32 `gson:"indent"`
}

func (this_ *CodePrinter_PrettyCodePrinter) BindSelf_(s any) {
	_ = s
	this_.CodePrinter_MappedCodePrinterBase.BindSelf_(s)
}
func (this_ *CodePrinter_PrettyCodePrinter) IsJavaNil_() bool { return this_ == nil }
func (*CodePrinter_PrettyCodePrinter) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CodePrinter$PrettyCodePrinter", "PrettyCodePrinter"
}
func (this_ *CodePrinter_PrettyCodePrinter) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.CodePrinter_MappedCodePrinterBase.IdentityHash_()
}
func (*CodePrinter_PrettyCodePrinter) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CodePrinter.PrettyCodePrinter", "com.google.javascript.jscomp.CodePrinter$PrettyCodePrinter", func(v_ any) bool {
		return j2gort.InstanceOf[*CodePrinter_PrettyCodePrinter](v_)
	})
}
func (*CodePrinter_PrettyCodePrinter) Implements_jscomp_CodePrinter_MappedCodePrinter() {}
func (*CodePrinter_PrettyCodePrinter) Implements_jscomp_CodeConsumer()                  {}

func StaticGet_CodePrinter_PrettyCodePrinter_INDENT() string { return "  " }
func (this_ *CodePrinter_PrettyCodePrinter) initFields_CodePrinter_PrettyCodePrinter() {
	_ = this_

}

// Java constructor: PrettyCodePrinter(int lineLengthThreshold, boolean createSourceMap, com.google.javascript.jscomp.SourceMap.DetailLevel sourceMapDetailLevel, com.google.javascript.jscomp.CodePrinter.@org.jspecify.annotations.Nullable LicenseTracker licenseTracker)
func NewCodePrinter_PrettyCodePrinter(lineLengthThreshold int32, createSourceMap bool, sourceMapDetailLevel *SourceMap_DetailLevel, licenseTracker CodePrinter_LicenseTracker) *CodePrinter_PrettyCodePrinter {
	this_ := &CodePrinter_PrettyCodePrinter{}
	this_.BindSelf_(this_)
	InitCodePrinter_PrettyCodePrinter(this_, lineLengthThreshold, createSourceMap, sourceMapDetailLevel, licenseTracker)
	return this_
}

func InitCodePrinter_PrettyCodePrinter(this_ *CodePrinter_PrettyCodePrinter, lineLengthThreshold int32, createSourceMap bool, sourceMapDetailLevel *SourceMap_DetailLevel, licenseTracker CodePrinter_LicenseTracker) {
	_ = lineLengthThreshold
	_ = createSourceMap
	_ = sourceMapDetailLevel
	_ = licenseTracker
	InitCodePrinter_MappedCodePrinter(&this_.CodePrinter_MappedCodePrinterBase, lineLengthThreshold, createSourceMap, sourceMapDetailLevel, licenseTracker)
	this_.initFields_CodePrinter_PrettyCodePrinter()
}

// Java method: PrettyCodePrinter.append(java.lang.String str)
func (this_ *CodePrinter_PrettyCodePrinter) Append(str string) {
	_ = str
	if (this_.GetCodePrinter_MappedCodePrinterBase().lineLength == int32(0)) && !(this_.IsInTemplateLiteral()) {
		{
			var i int32 = int32(0)
			_ = i
			for ; i < this_.indent; i = i + 1 {
				this_.GetCodePrinter_MappedCodePrinterBase().code.Append_1_string("  ")
				this_.GetCodePrinter_MappedCodePrinterBase().lineLength = (this_.GetCodePrinter_MappedCodePrinterBase().lineLength + j2gort.JStrLen("  "))
			}
		}
	}
	this_.CodePrinter_MappedCodePrinterBase.Append(str)
}

// Java method: PrettyCodePrinter.addNumber(double x, com.google.javascript.rhino.Node n)
func (this_ *CodePrinter_PrettyCodePrinter) AddNumber(x float64, n pkg_rhino.Node) {
	_ = x
	_ = n
	if !(pkg_base.JSCompDoubles_IsPositive(x)) {
		j2gort.CheckStateMsg(false, j2gort.ValueOfDouble(x))
	}
	var numberFromSource string = CodePrinter_PrettyCodePrinter_GetNumberFromSource(n)
	_ = numberFromSource
	if numberFromSource == j2gort.NullString {
		this_.CodePrinter_MappedCodePrinterBase.AddNumber(x, n)
		return
	}
	var d float64 = float64(0)
	_ = d
	var didRet_13 bool
	func() {
		defer func() {
			if rec_14 := recover(); rec_14 != nil {
				jex_15 := javalang.RecoverException(rec_14)
				_ = jex_15
				if j2gort.InstanceOf[javalang.INumberFormatException](jex_15) {
					e := j2gort.AsIface[javalang.INumberFormatException](jex_15)
					_ = e
					this_.CodePrinter_MappedCodePrinterBase.AddNumber(x, n)
					didRet_13 = true
					return
				} else {
					panic(rec_14)
				}
			}
		}()
		d = j2gort.ParseDouble(numberFromSource)
	}()
	if didRet_13 {
		return
	}
	if x != d {
		this_.CodePrinter_MappedCodePrinterBase.AddNumber(x, n)
		return
	}
	this_.AddConstant(numberFromSource)
}

// Java method: PrettyCodePrinter.startNewLine()
func (this_ *CodePrinter_PrettyCodePrinter) StartNewLine() {
	if (this_.GetCodePrinter_MappedCodePrinterBase().lineLength <= int32(0)) && !(this_.IsInTemplateLiteral()) {
		return
	}
	this_.GetCodePrinter_MappedCodePrinterBase().code.Append_1_uint16(uint16(10))
	this_.GetCodePrinter_MappedCodePrinterBase().lineIndex = this_.GetCodePrinter_MappedCodePrinterBase().lineIndex + 1
	this_.GetCodePrinter_MappedCodePrinterBase().lineLength = int32(0)
}

// Java method: PrettyCodePrinter.maybeLineBreak()
func (this_ *CodePrinter_PrettyCodePrinter) MaybeLineBreak() {
	this_.MaybeCutLine()
}

// Java method: PrettyCodePrinter.maybeCutLine()
func (this_ *CodePrinter_PrettyCodePrinter) MaybeCutLine() {
	if this_.GetCodePrinter_MappedCodePrinterBase().lineLength > this_.GetCodePrinter_MappedCodePrinterBase().lineLengthThreshold {
		this_.StartNewLine()
	}
}

// Java method: PrettyCodePrinter.endLine()
func (this_ *CodePrinter_PrettyCodePrinter) EndLine() {
	this_.StartNewLine()
}

// Java method: PrettyCodePrinter.appendBlockStart()
func (this_ *CodePrinter_PrettyCodePrinter) AppendBlockStart() {
	this_.MaybeInsertSpace()
	this_.Add_1_string("{")
	this_.indent = this_.indent + 1
}

// Java method: PrettyCodePrinter.appendBlockEnd()
func (this_ *CodePrinter_PrettyCodePrinter) AppendBlockEnd() {
	this_.MaybeEndStatement()
	this_.EndLine()
	this_.indent = this_.indent - 1
	this_.Add_1_string("}")
}

// Java method: PrettyCodePrinter.listSeparator()
func (this_ *CodePrinter_PrettyCodePrinter) ListSeparator() {
	this_.Add_1_string(", ")
	this_.MaybeLineBreak()
}

// Java method: PrettyCodePrinter.optionalListSeparator()
func (this_ *CodePrinter_PrettyCodePrinter) OptionalListSeparator() {
	this_.Add_1_string(",")
	this_.MaybeLineBreak()
}

// Java method: PrettyCodePrinter.endFunction(boolean statementContext)
func (this_ *CodePrinter_PrettyCodePrinter) EndFunction(statementContext bool) {
	_ = statementContext
	this_.CodePrinter_MappedCodePrinterBase.EndFunction(statementContext)
	if statementContext {
		this_.StartNewLine()
	}
}

// Java method: PrettyCodePrinter.beginCaseBody()
func (this_ *CodePrinter_PrettyCodePrinter) BeginCaseBody() {
	this_.CodePrinter_MappedCodePrinterBase.BeginCaseBody()
	this_.indent = this_.indent + 1
	this_.EndLine()
}

// Java method: PrettyCodePrinter.endCaseBody()
func (this_ *CodePrinter_PrettyCodePrinter) EndCaseBody() {
	this_.CodePrinter_MappedCodePrinterBase.EndCaseBody()
	this_.indent = this_.indent - 1
}

// Java method: PrettyCodePrinter.appendOp(java.lang.String op, boolean binOp)
func (this_ *CodePrinter_PrettyCodePrinter) AppendOp(op string, binOp bool) {
	_ = op
	_ = binOp
	if ((int32(this_.GetLastChar()) != int32(uint16(32))) && binOp) && (int32(j2gort.JStrCharAt(op, int32(0))) != int32(uint16(44))) {
		this_.Add_1_string(" ")
	}
	this_.Add_1_string(op)
	if binOp {
		this_.Add_1_string(" ")
	}
}

// Java method: PrettyCodePrinter.shouldPreserveExtras(com.google.javascript.rhino.Node n)
func (this_ *CodePrinter_PrettyCodePrinter) ShouldPreserveExtras(n pkg_rhino.Node) bool {
	_ = n
	if (!(n.IsBlock()) || !(n.IsAddedBlock())) || !(n.HasParent()) {
		return true
	}
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	var isElse bool = ((parent.IsIf() && parent.HasXChildren(int32(3))) && (n == parent.GetLastChild()))
	_ = isElse
	var onlyChildIsIf bool = (n.HasOneChild() && n.GetFirstChild().IsIf())
	_ = onlyChildIsIf
	if isElse && onlyChildIsIf {
		return false
	}
	return true
}

// Java method: PrettyCodePrinter.maybeInsertSpace()
func (this_ *CodePrinter_PrettyCodePrinter) MaybeInsertSpace() {
	if (int32(this_.GetLastChar()) != int32(uint16(32))) && (int32(this_.GetLastChar()) != int32(uint16(10))) {
		this_.Add_1_string(" ")
	}
}

// Java static method: PrettyCodePrinter.getTryForCatch(com.google.javascript.rhino.Node n)
func CodePrinter_PrettyCodePrinter_GetTryForCatch(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	return n.GetGrandparent()
}

// Java method: PrettyCodePrinter.breakAfterBlockFor(com.google.javascript.rhino.Node n, boolean isStatementContext)
func (this_ *CodePrinter_PrettyCodePrinter) BreakAfterBlockFor(n pkg_rhino.Node, isStatementContext bool) bool {
	_ = n
	_ = isStatementContext
	if !(n.IsBlock()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(n))
	}
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	var type_ pkg_rhino.Token = parent.GetToken()
	_ = type_
	switch (type_).Ordinal() {
	case int32(79) /* Token.DO */ :
		return false
	case int32(73) /* Token.FUNCTION */ :
		return false
	case int32(48) /* Token.TRY */ :
		return (n != parent.GetFirstChild())
	case int32(86) /* Token.CATCH */ :
		return !(NodeUtil_HasFinally(CodePrinter_PrettyCodePrinter_GetTryForCatch(parent)))
	case int32(74) /* Token.IF */ :
		return (n == parent.GetLastChild())
	default:
		return true
	}
	panic("unreachable switch expression")
}

// Java method: PrettyCodePrinter.endStatement(boolean needsSemicolon, boolean hasTrailingCommentOnSameLine)
func (this_ *CodePrinter_PrettyCodePrinter) EndStatement_2(needsSemicolon bool, hasTrailingCommentOnSameLine bool) {
	_ = needsSemicolon
	_ = hasTrailingCommentOnSameLine
	this_.Add_1_string(";")
	if !(hasTrailingCommentOnSameLine) {
		this_.EndLine()
	}
	this_.GetCodeConsumerBase().statementNeedsEnded = false
}

// Java method: PrettyCodePrinter.endFile()
func (this_ *CodePrinter_PrettyCodePrinter) EndFile() {
	this_.MaybeEndStatement()
}

// Java static method: PrettyCodePrinter.getNumberFromSource(com.google.javascript.rhino.Node n)
func CodePrinter_PrettyCodePrinter_GetNumberFromSource(n pkg_rhino.Node) string {
	_ = n
	if !(n.IsNumber()) {
		return j2gort.NullString
	}
	var staticSrc pkg_rhino.StaticSourceFile = NodeUtil_GetSourceFile(n)
	_ = staticSrc
	var src *SourceFile
	_ = src
	var ok_16 bool
	src, ok_16 = j2gort.TryCast[*SourceFile](staticSrc)
	if !ok_16 {
		return j2gort.NullString
	}
	var srcCode string = j2gort.NullString
	_ = srcCode
	var didRet_17 bool
	var retVal_18 string
	func() {
		defer func() {
			if rec_19 := recover(); rec_19 != nil {
				jex_20 := javalang.RecoverException(rec_19)
				_ = jex_20
				if j2gort.InstanceOf[javaio.IIOException](jex_20) {
					e := j2gort.AsIface[javaio.IIOException](jex_20)
					_ = e
					retVal_18 = j2gort.NullString
					didRet_17 = true
					return
				} else {
					panic(rec_19)
				}
			}
		}()
		if src.IsStubSourceFileForAlreadyProvidedInput() {
			retVal_18 = j2gort.NullString
			didRet_17 = true
			return
		}
		srcCode = src.GetCode()
	}()
	if didRet_17 {
		return retVal_18
	}
	var offset int32 = int32(0)
	_ = offset
	var didRet_21 bool
	var retVal_22 string
	func() {
		defer func() {
			if rec_23 := recover(); rec_23 != nil {
				jex_24 := javalang.RecoverException(rec_23)
				_ = jex_24
				if j2gort.InstanceOf[javalang.IIllegalArgumentException](jex_24) {
					e_2 := j2gort.AsIface[javalang.IIllegalArgumentException](jex_24)
					_ = e_2
					retVal_22 = j2gort.NullString
					didRet_21 = true
					return
				} else {
					panic(rec_23)
				}
			}
		}()
		offset = n.GetSourceOffset()
	}()
	if didRet_21 {
		return retVal_22
	}
	var endOffset int32 = (offset + n.GetLength())
	_ = endOffset
	if (offset < int32(0)) || (endOffset > j2gort.JStrLen(srcCode)) {
		return j2gort.NullString
	}
	return j2gort.JStrSub(srcCode, offset, endOffset)
}

// Java class: com.google.javascript.jscomp.CodePrinter.CompactCodePrinter
type CodePrinter_CompactCodePrinter struct {
	CodePrinter_MappedCodePrinterBase
	lineStartPosition      int32 `gson:"lineStartPosition"`
	preferredBreakPosition int32 `gson:"preferredBreakPosition"`
	lineBreak              bool  `gson:"lineBreak"`
}

func (this_ *CodePrinter_CompactCodePrinter) BindSelf_(s any) {
	_ = s
	this_.CodePrinter_MappedCodePrinterBase.BindSelf_(s)
}
func (this_ *CodePrinter_CompactCodePrinter) IsJavaNil_() bool { return this_ == nil }
func (*CodePrinter_CompactCodePrinter) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CodePrinter$CompactCodePrinter", "CompactCodePrinter"
}
func (this_ *CodePrinter_CompactCodePrinter) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.CodePrinter_MappedCodePrinterBase.IdentityHash_()
}
func (*CodePrinter_CompactCodePrinter) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CodePrinter.CompactCodePrinter", "com.google.javascript.jscomp.CodePrinter$CompactCodePrinter", func(v_ any) bool {
		return j2gort.InstanceOf[*CodePrinter_CompactCodePrinter](v_)
	})
}
func (*CodePrinter_CompactCodePrinter) Implements_jscomp_CodePrinter_MappedCodePrinter() {}
func (*CodePrinter_CompactCodePrinter) Implements_jscomp_CodeConsumer()                  {}

func (this_ *CodePrinter_CompactCodePrinter) initFields_CodePrinter_CompactCodePrinter() {
	_ = this_

}

// Java constructor: CompactCodePrinter(boolean lineBreak, int lineLengthThreshold, boolean createSrcMap, com.google.javascript.jscomp.SourceMap.DetailLevel sourceMapDetailLevel, com.google.javascript.jscomp.CodePrinter.LicenseTracker licenseTracker)
func NewCodePrinter_CompactCodePrinter(lineBreak bool, lineLengthThreshold int32, createSrcMap bool, sourceMapDetailLevel *SourceMap_DetailLevel, licenseTracker CodePrinter_LicenseTracker) *CodePrinter_CompactCodePrinter {
	this_ := &CodePrinter_CompactCodePrinter{}
	this_.BindSelf_(this_)
	InitCodePrinter_CompactCodePrinter(this_, lineBreak, lineLengthThreshold, createSrcMap, sourceMapDetailLevel, licenseTracker)
	return this_
}

func InitCodePrinter_CompactCodePrinter(this_ *CodePrinter_CompactCodePrinter, lineBreak bool, lineLengthThreshold int32, createSrcMap bool, sourceMapDetailLevel *SourceMap_DetailLevel, licenseTracker CodePrinter_LicenseTracker) {
	_ = lineBreak
	_ = lineLengthThreshold
	_ = createSrcMap
	_ = sourceMapDetailLevel
	_ = licenseTracker
	InitCodePrinter_MappedCodePrinter(&this_.CodePrinter_MappedCodePrinterBase, lineLengthThreshold, createSrcMap, sourceMapDetailLevel, licenseTracker)
	this_.initFields_CodePrinter_CompactCodePrinter()
	this_.lineBreak = lineBreak
}

// Java method: CompactCodePrinter.startNewLine()
func (this_ *CodePrinter_CompactCodePrinter) StartNewLine() {
	if (this_.GetCodePrinter_MappedCodePrinterBase().lineLength <= int32(0)) && !(this_.IsInTemplateLiteral()) {
		return
	}
	this_.GetCodePrinter_MappedCodePrinterBase().code.Append_1_uint16(uint16(10))
	this_.GetCodePrinter_MappedCodePrinterBase().lineLength = int32(0)
	this_.GetCodePrinter_MappedCodePrinterBase().lineIndex = this_.GetCodePrinter_MappedCodePrinterBase().lineIndex + 1
	this_.lineStartPosition = this_.GetCodePrinter_MappedCodePrinterBase().code.Length()
}

// Java method: CompactCodePrinter.maybeLineBreak()
func (this_ *CodePrinter_CompactCodePrinter) MaybeLineBreak() {
	if this_.lineBreak {
		if this_.GetCodeConsumerBase().sawFunction {
			this_.StartNewLine()
			this_.GetCodeConsumerBase().sawFunction = false
		}
	}
	var len_ int32 = this_.GetCodePrinter_MappedCodePrinterBase().code.Length()
	_ = len_
	if this_.preferredBreakPosition == (len_ - int32(1)) {
		var ch uint16 = this_.GetCodePrinter_MappedCodePrinterBase().code.CharAt((len_ - int32(1)))
		_ = ch
		if int32(ch) == int32(uint16(59)) {
			this_.preferredBreakPosition = len_
		}
	}
	this_.MaybeCutLine()
}

// Java method: CompactCodePrinter.maybeCutLine()
func (this_ *CodePrinter_CompactCodePrinter) MaybeCutLine() {
	if this_.GetCodePrinter_MappedCodePrinterBase().lineLength <= this_.GetCodePrinter_MappedCodePrinterBase().lineLengthThreshold {
		return
	}
	if (this_.preferredBreakPosition > this_.lineStartPosition) && (this_.preferredBreakPosition < (this_.lineStartPosition + this_.GetCodePrinter_MappedCodePrinterBase().lineLength)) {
		this_.GetCodePrinter_MappedCodePrinterBase().code.Insert_2_int32_uint16(this_.preferredBreakPosition, uint16(10))
		this_.ReportLineCut(this_.GetCodePrinter_MappedCodePrinterBase().lineIndex, (this_.preferredBreakPosition - this_.lineStartPosition))
		this_.GetCodePrinter_MappedCodePrinterBase().lineIndex = this_.GetCodePrinter_MappedCodePrinterBase().lineIndex + 1
		this_.GetCodePrinter_MappedCodePrinterBase().lineLength = (this_.GetCodePrinter_MappedCodePrinterBase().lineLength - (this_.preferredBreakPosition - this_.lineStartPosition))
		this_.lineStartPosition = (this_.preferredBreakPosition + int32(1))
	} else {
		this_.StartNewLine()
	}
}

// Java method: CompactCodePrinter.notePreferredLineBreak()
func (this_ *CodePrinter_CompactCodePrinter) NotePreferredLineBreak() {
	this_.preferredBreakPosition = this_.GetCodePrinter_MappedCodePrinterBase().code.Length()
}

// Java interface: com.google.javascript.jscomp.CodePrinter.LicenseTracker
type CodePrinter_LicenseTracker interface {
	Implements_jscomp_CodePrinter_LicenseTracker()
	IsJavaNil_() bool
	TrackLicensesForNode(pkg_rhino.Node)
	EmitLicenses() j2gort.Ext_com_google_common_collect_ImmutableSet
}

func ToIface_CodePrinter_LicenseTracker[T interface {
	CodePrinter_LicenseTracker
	comparable
}](v T) CodePrinter_LicenseTracker {
	var z T
	if v == z {
		return nil
	}
	return v
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CodePrinter.LicenseTracker", "com.google.javascript.jscomp.CodePrinter$LicenseTracker", func(v_ any) bool {
		return j2gort.InstanceOf[CodePrinter_LicenseTracker](v_)
	})
}

// Java class: com.google.javascript.jscomp.CodePrinter.Builder
type CodePrinter_Builder struct {
	options              *CompilerOptions                         `gson:"options"`
	licenseTracker       CodePrinter_LicenseTracker               `gson:"licenseTracker"`
	registry             *pkg_rhino.JSTypeRegistry                `gson:"registry"`
	codeGeneratorFactory CodePrinter_Builder_CodeGeneratorFactory `gson:"codeGeneratorFactory"`
	root                 pkg_rhino.Node                           `gson:"root"`
	idHash_              atomic.Int32
	lineBreak            bool `gson:"lineBreak"`
	prettyPrint          bool `gson:"prettyPrint"`
	outputTypes          bool `gson:"outputTypes"`
	tagAsTypeSummary     bool `gson:"tagAsTypeSummary"`
	tagAsStrict          bool `gson:"tagAsStrict"`
}

func (this_ *CodePrinter_Builder) BindSelf_(s any) {
	_ = s
}
func (this_ *CodePrinter_Builder) IsJavaNil_() bool { return this_ == nil }
func (*CodePrinter_Builder) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CodePrinter$Builder", "Builder"
}
func (this_ *CodePrinter_Builder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CodePrinter_Builder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CodePrinter.Builder", "com.google.javascript.jscomp.CodePrinter$Builder", func(v_ any) bool {
		return j2gort.InstanceOf[*CodePrinter_Builder](v_)
	})
}

func (this_ *CodePrinter_Builder) initFields_CodePrinter_Builder() {
	_ = this_

	this_.codeGeneratorFactory = NewCodePrinter__Anon_10(this_)
}

// Java constructor: Builder(com.google.javascript.rhino.Node node)
func NewCodePrinter_Builder(node pkg_rhino.Node) *CodePrinter_Builder {
	this_ := &CodePrinter_Builder{}
	InitCodePrinter_Builder(this_, node)
	return this_
}

func InitCodePrinter_Builder(this_ *CodePrinter_Builder, node pkg_rhino.Node) {
	_ = node
	this_.initFields_CodePrinter_Builder()
	this_.root = node
}

// Java method: Builder.getOptions()
func (this_ *CodePrinter_Builder) priv_Builder_getOptions_1362() *CompilerOptions {
	if this_.options == nil {
		this_.options = NewCompilerOptions()
	}
	return this_.options
}

// Java method: Builder.setCompilerOptions(com.google.javascript.jscomp.CompilerOptions options)
func (this_ *CodePrinter_Builder) SetCompilerOptions(options *CompilerOptions) *CodePrinter_Builder {
	_ = options
	this_.options = options
	this_.prettyPrint = options.IsPrettyPrint()
	this_.lineBreak = options.ShouldAddLineBreak()
	return this_
}

// Java method: Builder.setTypeRegistry(com.google.javascript.rhino.jstype.JSTypeRegistry registry)
func (this_ *CodePrinter_Builder) SetTypeRegistry(registry *pkg_rhino.JSTypeRegistry) *CodePrinter_Builder {
	_ = registry
	this_.registry = registry
	return this_
}

// Java method: Builder.setPrettyPrint(boolean prettyPrint)
func (this_ *CodePrinter_Builder) SetPrettyPrint(prettyPrint bool) *CodePrinter_Builder {
	_ = prettyPrint
	this_.prettyPrint = prettyPrint
	return this_
}

// Java method: Builder.setLineBreak(boolean lineBreak)
func (this_ *CodePrinter_Builder) SetLineBreak(lineBreak bool) *CodePrinter_Builder {
	_ = lineBreak
	this_.lineBreak = lineBreak
	return this_
}

// Java method: Builder.setOutputTypes(boolean outputTypes)
func (this_ *CodePrinter_Builder) SetOutputTypes(outputTypes bool) *CodePrinter_Builder {
	_ = outputTypes
	this_.outputTypes = outputTypes
	return this_
}

// Java method: Builder.setLicenseTracker(com.google.javascript.jscomp.CodePrinter.LicenseTracker licenseTracker)
func (this_ *CodePrinter_Builder) SetLicenseTracker(licenseTracker CodePrinter_LicenseTracker) *CodePrinter_Builder {
	_ = licenseTracker
	this_.licenseTracker = licenseTracker
	return this_
}

// Java method: Builder.setTagAsTypeSummary(boolean tagAsTypeSummary)
func (this_ *CodePrinter_Builder) SetTagAsTypeSummary(tagAsTypeSummary bool) *CodePrinter_Builder {
	_ = tagAsTypeSummary
	this_.tagAsTypeSummary = tagAsTypeSummary
	return this_
}

// Java method: Builder.setTagAsStrict(boolean tagAsStrict)
func (this_ *CodePrinter_Builder) SetTagAsStrict(tagAsStrict bool) *CodePrinter_Builder {
	_ = tagAsStrict
	this_.tagAsStrict = tagAsStrict
	return this_
}

// Java method: Builder.setCodeGeneratorFactory(com.google.javascript.jscomp.CodePrinter.Builder.CodeGeneratorFactory factory)
func (this_ *CodePrinter_Builder) SetCodeGeneratorFactory(factory CodePrinter_Builder_CodeGeneratorFactory) *CodePrinter_Builder {
	_ = factory
	this_.codeGeneratorFactory = factory
	return this_
}

// Java method: Builder.build()
func (this_ *CodePrinter_Builder) Build() string {
	return this_.BuildWithSourceMappings().source
}

// Java method: Builder.buildWithSourceMappings()
func (this_ *CodePrinter_Builder) BuildWithSourceMappings() *CodePrinter_SourceAndMappings {
	if this_.root == 0 {
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Cannot build without root node being specified")))
	}
	var opts *CompilerOptions = this_.priv_Builder_getOptions_1362()
	_ = opts
	return CodePrinter_ToSource(this_.root, CodePrinter_Format_FromOptions(opts, this_.outputTypes, this_.prettyPrint), opts, this_.licenseTracker, this_.tagAsTypeSummary, this_.tagAsStrict, this_.lineBreak, this_.codeGeneratorFactory)
}

// Java interface: com.google.javascript.jscomp.CodePrinter.Builder.CodeGeneratorFactory
type CodePrinter_Builder_CodeGeneratorFactory interface {
	Implements_jscomp_CodePrinter_Builder_CodeGeneratorFactory()
	IsJavaNil_() bool
	GetCodeGenerator(CodePrinter_Format, CodeConsumer) CodeGenerator
}

func ToIface_CodePrinter_Builder_CodeGeneratorFactory[T interface {
	CodePrinter_Builder_CodeGeneratorFactory
	comparable
}](v T) CodePrinter_Builder_CodeGeneratorFactory {
	var z T
	if v == z {
		return nil
	}
	return v
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CodePrinter.Builder.CodeGeneratorFactory", "com.google.javascript.jscomp.CodePrinter$Builder$CodeGeneratorFactory", func(v_ any) bool {
		return j2gort.InstanceOf[CodePrinter_Builder_CodeGeneratorFactory](v_)
	})
}

// Java enum: com.google.javascript.jscomp.CodePrinter.Format
type CodePrinter_Format uint8

const (
	enum_CodePrinter_Format_COMPACT CodePrinter_Format = 1
	enum_CodePrinter_Format_PRETTY  CodePrinter_Format = 2
	enum_CodePrinter_Format_TYPED   CodePrinter_Format = 3
)

var enumBaseTable_CodePrinter_Format = [...]j2gort.EnumBase{
	{Name_: "COMPACT", Ordinal_: 0},
	{Name_: "PRETTY", Ordinal_: 1},
	{Name_: "TYPED", Ordinal_: 2},
}
var enumValues_CodePrinter_Format = []CodePrinter_Format{1, 2, 3}

func StaticGet_CodePrinter_Format_COMPACT() CodePrinter_Format {
	return enum_CodePrinter_Format_COMPACT
}
func StaticGet_CodePrinter_Format_PRETTY() CodePrinter_Format {
	return enum_CodePrinter_Format_PRETTY
}
func StaticGet_CodePrinter_Format_TYPED() CodePrinter_Format {
	return enum_CodePrinter_Format_TYPED
}

func CodePrinter_Format_Values() []CodePrinter_Format {
	return j2gort.CloneSlice(enumValues_CodePrinter_Format)
}

func CodePrinter_Format_ValueOf(name string) CodePrinter_Format {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_CodePrinter_Format {
		if enumBaseTable_CodePrinter_Format[i_].Name_ == name {
			return CodePrinter_Format(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.CodePrinter.Format." + name))
}

// Java static method: Format.fromOptions(com.google.javascript.jscomp.CompilerOptions options, boolean outputTypes, boolean prettyPrint)
func CodePrinter_Format_FromOptions(options *CompilerOptions, outputTypes bool, prettyPrint bool) CodePrinter_Format {
	_ = options
	_ = outputTypes
	_ = prettyPrint
	if outputTypes {
		return StaticGet_CodePrinter_Format_TYPED()
	}
	if prettyPrint || options.GetOutputFeatureSet().Contains_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_TYPE_ANNOTATION()) {
		return StaticGet_CodePrinter_Format_PRETTY()
	}
	return StaticGet_CodePrinter_Format_COMPACT()
}

func (this_ CodePrinter_Format) IsJavaNil_() bool                              { return this_ == 0 }
func (CodePrinter_Format) HasIdentityEquals_()                                 {}
func (CodePrinter_Format) Implements_j2gort_Ext_java_lang_Enum()               {}
func (CodePrinter_Format) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ CodePrinter_Format) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (CodePrinter_Format) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CodePrinter$Format", "Format"
}
func (this_ CodePrinter_Format) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_CodePrinter_Format[this_-1]
}
func (this_ CodePrinter_Format) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ CodePrinter_Format) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ CodePrinter_Format) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CodePrinter_Format[this_-1].Name_
}
func (this_ CodePrinter_Format) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_CodePrinter_Format[this_-1].Name_
}
func (this_ CodePrinter_Format) Equals(other_ any) bool {
	o_, ok_ := other_.(CodePrinter_Format)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ CodePrinter_Format) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_CodePrinter_Format[this_-1].HashCode()
}
func (this_ CodePrinter_Format) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CodePrinter_Format[this_-1].CompareTo(other_)
}
func (this_ CodePrinter_Format) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_CodePrinter_Format[this_-1].Name_)
}
func init() {
	vals_ := []any{CodePrinter_Format(1), CodePrinter_Format(2), CodePrinter_Format(3)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.CodePrinter.Format", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.CodePrinter$Format", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CodePrinter.Format", "com.google.javascript.jscomp.CodePrinter$Format", func(v_ any) bool {
		return j2gort.InstanceOf[CodePrinter_Format](v_)
	})
}

// Java class: com.google.javascript.jscomp.CodePrinter.SourceAndMappings
type CodePrinter_SourceAndMappings struct {
	source   string                    `gson:"source"`
	mappings j2gort.Ext_java_util_List `gson:"mappings"`
	idHash_  atomic.Int32
}

func (this_ *CodePrinter_SourceAndMappings) BindSelf_(s any) {
	_ = s
}
func (this_ *CodePrinter_SourceAndMappings) IsJavaNil_() bool { return this_ == nil }
func (*CodePrinter_SourceAndMappings) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CodePrinter$SourceAndMappings", "SourceAndMappings"
}
func (this_ *CodePrinter_SourceAndMappings) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CodePrinter_SourceAndMappings) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CodePrinter.SourceAndMappings", "com.google.javascript.jscomp.CodePrinter$SourceAndMappings", func(v_ any) bool {
		return j2gort.InstanceOf[*CodePrinter_SourceAndMappings](v_)
	})
}

func (this_ *CodePrinter_SourceAndMappings) initFields_CodePrinter_SourceAndMappings() {
	_ = this_
	this_.source = j2gort.NullString
}

// Java constructor: SourceAndMappings()
func NewCodePrinter_SourceAndMappings() *CodePrinter_SourceAndMappings {
	this_ := &CodePrinter_SourceAndMappings{}
	InitCodePrinter_SourceAndMappings(this_)
	return this_
}

func InitCodePrinter_SourceAndMappings(this_ *CodePrinter_SourceAndMappings) {
	this_.initFields_CodePrinter_SourceAndMappings()
}

// Java class:
type CodePrinter__Anon_10 struct {
	outer_  *CodePrinter_Builder
	idHash_ atomic.Int32
}

func (this_ *CodePrinter__Anon_10) BindSelf_(s any) {
	_ = s
}
func (this_ *CodePrinter__Anon_10) IsJavaNil_() bool { return this_ == nil }
func (*CodePrinter__Anon_10) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CodePrinter$Builder$1", ""
}
func (this_ *CodePrinter__Anon_10) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CodePrinter__Anon_10) HasIdentityEquals_()                                         {}
func (this_ *CodePrinter__Anon_10) SetOuter_(o *CodePrinter_Builder)                      { this_.outer_ = o }
func (this_ *CodePrinter__Anon_10) GetOuter_() *CodePrinter_Builder                       { return this_.outer_ }
func (*CodePrinter__Anon_10) Implements_jscomp_CodePrinter_Builder_CodeGeneratorFactory() {}

func (this_ *CodePrinter__Anon_10) initFields_CodePrinter__Anon_10() {
	_ = this_
}

func NewCodePrinter__Anon_10(outer_ *CodePrinter_Builder) *CodePrinter__Anon_10 {
	this_ := &CodePrinter__Anon_10{}
	this_.outer_ = outer_
	this_.initFields_CodePrinter__Anon_10()
	return this_
}

// Java method: .getCodeGenerator(com.google.javascript.jscomp.CodePrinter.Format outputFormat, com.google.javascript.jscomp.CodeConsumer cc)
func (this_ *CodePrinter__Anon_10) GetCodeGenerator(outputFormat CodePrinter_Format, cc CodeConsumer) CodeGenerator {
	_ = outputFormat
	_ = cc
	var opts *CompilerOptions = this_.outer_.priv_Builder_getOptions_1362()
	_ = opts
	if outputFormat == StaticGet_CodePrinter_Format_TYPED() {
		return NewTypedCodeGenerator(cc, opts, this_.outer_.registry)
	} else {
		return NewCodeGenerator_2(cc, opts)
	}
}
