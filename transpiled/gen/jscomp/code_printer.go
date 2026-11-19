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

// Java static method: CodePrinter.toSource(com.google.javascript.rhino.Node root, com.google.javascript.jscomp.CodePrinter.Format outputFormat, com.google.javascript.jscomp.CompilerOptions options, boolean gatherSourceMapInfo, com.google.javascript.jscomp.CodePrinter.@org.jspecify.annotations.Nullable LicenseTracker licenseTracker, boolean tagAsTypeSummary, boolean tagAsStrict, boolean lineBreak, com.google.javascript.jscomp.CodePrinter.Builder.CodeGeneratorFactory codeGeneratorFactory)
func CodePrinter_ToSource(root pkg_rhino.Node, outputFormat CodePrinter_Format, options *CompilerOptions, gatherSourceMapInfo bool, licenseTracker CodePrinter_LicenseTracker, tagAsTypeSummary bool, tagAsStrict bool, lineBreak bool, codeGeneratorFactory CodePrinter_Builder_CodeGeneratorFactory) *CodePrinter_SourceAndMappings {
	_ = root
	_ = outputFormat
	_ = options
	_ = gatherSourceMapInfo
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
		mcp = NewCodePrinter_CompactCodePrinter(lineBreak, options.GetLineLengthThreshold(), gatherSourceMapInfo, options.GetSourceMapDetailLevel(), licenseTracker)
	} else {
		mcp = NewCodePrinter_PrettyCodePrinter(options.GetLineLengthThreshold(), gatherSourceMapInfo, options.GetSourceMapDetailLevel(), licenseTracker)
	}
	var cg CodeGenerator = codeGeneratorFactory.GetCodeGenerator(outputFormat, Upcast_CodeConsumer(mcp))
	_ = cg
	if tagAsTypeSummary {
		cg.GetCodeGeneratorBase().TagAsTypeSummary()
	}
	if tagAsStrict {
		cg.GetCodeGeneratorBase().TagAsStrict()
	}
	cg.GetCodeGeneratorBase().Add_1_Node(root)
	mcp.EndFile()
	var val_1 CodePrinter_SourceAndMappings
	InitCodePrinter_SourceAndMappings(&val_1)
	var result *CodePrinter_SourceAndMappings = &val_1
	_ = result
	{
		r_2 := result
		r_2.source = mcp.GetCodePrinter_MappedCodePrinterBase().GetCode()
	}
	if gatherSourceMapInfo {
		mcp.GetCodePrinter_MappedCodePrinterBase().FinishSourceMappings(result.source, result)
	}
	return result
}

// Java polymorphic class interface: com.google.javascript.jscomp.CodePrinter.MappedCodePrinter
type CodePrinter_MappedCodePrinter interface {
	Implements_jscomp_CodePrinter_MappedCodePrinter()
	IsJavaNil_() bool
	GetCodePrinter_MappedCodePrinterBase() *CodePrinter_MappedCodePrinterBase
	CodeConsumer
	FinishSourceMappings(string, *CodePrinter_SourceAndMappings)
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

func Upcast_CodePrinter_MappedCodePrinter[T CodePrinter_MappedCodePrinter](v T) CodePrinter_MappedCodePrinter {
	var vi CodePrinter_MappedCodePrinter = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.jscomp.CodePrinter.MappedCodePrinter
type CodePrinter_MappedCodePrinterBase struct {
	CodeConsumerBase
	mappingNodes                []pkg_rhino.Node                    `gson:"mappingNodes"`
	startLines                  []int32                             `gson:"startLines"`
	startCols                   []int32                             `gson:"startCols"`
	endLines                    []int32                             `gson:"endLines"`
	endCols                     []int32                             `gson:"endCols"`
	openMappingStack            []int32                             `gson:"openMappingStack"`
	completeMappingIds          []int32                             `gson:"completeMappingIds"`
	sourceMapDetailLevel        *SourceMap_DetailLevel              `gson:"sourceMapDetailLevel"`
	licenseTracker              CodePrinter_LicenseTracker          `gson:"licenseTracker"`
	code                        *j2gort.Ext_java_lang_StringBuilder `gson:"code"`
	mappingCount                int32                               `gson:"mappingCount"`
	openMappingDepth            int32                               `gson:"openMappingDepth"`
	completeMappingCount        int32                               `gson:"completeMappingCount"`
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

var static_CodePrinter_MappedCodePrinter_EMPTY_INT_ARRAY []int32

func StaticGet_CodePrinter_MappedCodePrinter_EMPTY_INT_ARRAY() []int32 {
	clinit_CodePrinter_MappedCodePrinter()
	return static_CodePrinter_MappedCodePrinter_EMPTY_INT_ARRAY
}
func StaticSet_CodePrinter_MappedCodePrinter_EMPTY_INT_ARRAY(v []int32) {
	clinit_CodePrinter_MappedCodePrinter()
	static_CodePrinter_MappedCodePrinter_EMPTY_INT_ARRAY = v
}

var clinitDone_CodePrinter_MappedCodePrinter uint32
var clinitState_CodePrinter_MappedCodePrinter j2gort.ClinitState

func clinit_CodePrinter_MappedCodePrinter() {
	if atomic.LoadUint32(&clinitDone_CodePrinter_MappedCodePrinter) != 0 {
		return
	}
	clinitSlow_CodePrinter_MappedCodePrinter()
}

func clinitSlow_CodePrinter_MappedCodePrinter() {
	if !j2gort.EnterClinit(&clinitState_CodePrinter_MappedCodePrinter, &clinitDone_CodePrinter_MappedCodePrinter) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_CodePrinter_MappedCodePrinter, &clinitDone_CodePrinter_MappedCodePrinter)
	static_CodePrinter_MappedCodePrinter_EMPTY_INT_ARRAY = make([]int32, int32(0))
}

func (this_ *CodePrinter_MappedCodePrinterBase) initFields_CodePrinter_MappedCodePrinter() {
	_ = this_

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
		var id int32 = this_.mappingCount
		_ = id
		var nodes []pkg_rhino.Node = this_.mappingNodes
		_ = nodes
		if nodes == nil {
			nodes = make([]pkg_rhino.Node, int32(64))
			this_.mappingNodes = nodes
			this_.startLines = make([]int32, int32(64))
			this_.startCols = make([]int32, int32(64))
			this_.endLines = make([]int32, int32(64))
			this_.endCols = make([]int32, int32(64))
			this_.openMappingStack = make([]int32, int32(32))
			this_.completeMappingIds = make([]int32, int32(64))
		} else if id == int32(len(nodes)) {
			var newCap int32 = (id * int32(2))
			_ = newCap
			nodes = func() []pkg_rhino.Node {
				arrSrc_3 := nodes
				arrLen_4 := newCap
				j2gort.CheckNotNullSlice(arrSrc_3)
				if arrLen_4 < 0 {
					panic(j2gort.ConstructShimThrowable("java.lang.NegativeArraySizeException"))
				}
				arrOut_5 := make([]pkg_rhino.Node, arrLen_4)
				copy(arrOut_5, arrSrc_3)
				return arrOut_5
			}()
			this_.mappingNodes = nodes
			this_.startLines = func() []int32 {
				arrSrc_6 := this_.startLines
				arrLen_7 := newCap
				j2gort.CheckNotNullSlice(arrSrc_6)
				if arrLen_7 < 0 {
					panic(j2gort.ConstructShimThrowable("java.lang.NegativeArraySizeException"))
				}
				arrOut_8 := make([]int32, arrLen_7)
				copy(arrOut_8, arrSrc_6)
				return arrOut_8
			}()
			this_.startCols = func() []int32 {
				arrSrc_9 := this_.startCols
				arrLen_10 := newCap
				j2gort.CheckNotNullSlice(arrSrc_9)
				if arrLen_10 < 0 {
					panic(j2gort.ConstructShimThrowable("java.lang.NegativeArraySizeException"))
				}
				arrOut_11 := make([]int32, arrLen_10)
				copy(arrOut_11, arrSrc_9)
				return arrOut_11
			}()
			this_.endLines = func() []int32 {
				arrSrc_12 := this_.endLines
				arrLen_13 := newCap
				j2gort.CheckNotNullSlice(arrSrc_12)
				if arrLen_13 < 0 {
					panic(j2gort.ConstructShimThrowable("java.lang.NegativeArraySizeException"))
				}
				arrOut_14 := make([]int32, arrLen_13)
				copy(arrOut_14, arrSrc_12)
				return arrOut_14
			}()
			this_.endCols = func() []int32 {
				arrSrc_15 := this_.endCols
				arrLen_16 := newCap
				j2gort.CheckNotNullSlice(arrSrc_15)
				if arrLen_16 < 0 {
					panic(j2gort.ConstructShimThrowable("java.lang.NegativeArraySizeException"))
				}
				arrOut_17 := make([]int32, arrLen_16)
				copy(arrOut_17, arrSrc_15)
				return arrOut_17
			}()
		}
		{
			a_18 := nodes
			i_19 := id
			a_18[i_19] = node
		}
		{
			a_20 := this_.startLines
			i_21 := id
			a_20[i_21] = line
		}
		{
			a_22 := this_.startCols
			i_23 := id
			a_22[i_23] = index
		}
		this_.mappingCount = (id + int32(1))
		var depth int32 = this_.openMappingDepth
		_ = depth
		var stack []int32 = this_.openMappingStack
		_ = stack
		if depth == int32(len(stack)) {
			stack = func() []int32 {
				arrSrc_24 := stack
				arrLen_25 := (depth * int32(2))
				j2gort.CheckNotNullSlice(arrSrc_24)
				if arrLen_25 < 0 {
					panic(j2gort.ConstructShimThrowable("java.lang.NegativeArraySizeException"))
				}
				arrOut_26 := make([]int32, arrLen_25)
				copy(arrOut_26, arrSrc_24)
				return arrOut_26
			}()
			this_.openMappingStack = stack
		}
		{
			a_27 := stack
			i_28 := depth
			a_27[i_28] = id
		}
		this_.openMappingDepth = (depth + int32(1))
	}
}

// Java method: MappedCodePrinter.endSourceMapping(com.google.javascript.rhino.Node node)
func (this_ *CodePrinter_MappedCodePrinterBase) EndSourceMapping(node pkg_rhino.Node) {
	_ = node
	if (this_.createSrcMap && (this_.openMappingDepth > int32(0))) && (this_.mappingNodes[this_.openMappingStack[(this_.openMappingDepth-int32(1))]] == node) {
		var depth int32 = (this_.openMappingDepth - int32(1))
		_ = depth
		this_.openMappingDepth = depth
		var id int32 = this_.openMappingStack[depth]
		_ = id
		var line int32 = this_.GetCurrentLineIndex()
		_ = line
		var index int32 = this_.GetCurrentCharIndex()
		_ = index
		if !(line >= int32(0)) {
			j2gort.CheckState(false)
		}
		{
			a_29 := this_.endLines
			i_30 := id
			a_29[i_30] = line
		}
		{
			a_31 := this_.endCols
			i_32 := id
			a_31[i_32] = index
		}
		var compIdx int32 = this_.completeMappingCount
		_ = compIdx
		var compIds []int32 = this_.completeMappingIds
		_ = compIds
		if compIdx == int32(len(compIds)) {
			compIds = func() []int32 {
				arrSrc_33 := compIds
				arrLen_34 := (compIdx * int32(2))
				j2gort.CheckNotNullSlice(arrSrc_33)
				if arrLen_34 < 0 {
					panic(j2gort.ConstructShimThrowable("java.lang.NegativeArraySizeException"))
				}
				arrOut_35 := make([]int32, arrLen_34)
				copy(arrOut_35, arrSrc_33)
				return arrOut_35
			}()
			this_.completeMappingIds = compIds
		}
		{
			a_36 := compIds
			i_37 := compIdx
			a_36[i_37] = id
		}
		this_.completeMappingCount = (compIdx + int32(1))
	}
}

// Java method: MappedCodePrinter.finishSourceMappings(java.lang.String code, com.google.javascript.jscomp.CodePrinter.SourceAndMappings out)
func (this_ *CodePrinter_MappedCodePrinterBase) FinishSourceMappings(code string, out *CodePrinter_SourceAndMappings) {
	_ = code
	_ = out
	if !(this_.createSrcMap) || (this_.mappingCount == int32(0)) {
		return
	}
	var lineLengths []int32 = CodePrinter_MappedCodePrinter_ComputeLineLengths(code)
	_ = lineLengths
	var lineCount int32 = int32(len(lineLengths))
	_ = lineCount
	var count int32 = this_.mappingCount
	_ = count
	var eLines []int32 = this_.endLines
	_ = eLines
	var eCols []int32 = this_.endCols
	_ = eCols
	{
		var i int32 = int32(0)
		_ = i
		for ; i < count; i = i + 1 {
			var line int32 = eLines[i]
			_ = line
			if line < lineCount {
				var col int32 = eCols[i]
				_ = col
				var lineLen int32 = lineLengths[line]
				_ = lineLen
				if col > lineLen {
					j2gort.CheckStateFmt(false, "End position %s points to a column larger than line length %s", pkg_sourcemap.NewFilePosition(line, col), j2gort.ValueOfInt(lineLen))
				}
				if col == lineLen {
					{
						a_38 := eLines
						i_39 := i
						a_38[i_39] = (line + int32(1))
					}
					{
						a_40 := eCols
						i_41 := i
						a_40[i_41] = int32(0)
					}
				}
			}
		}
	}
	{
		r_42 := out
		r_42.mappingNodes = this_.mappingNodes
	}
	{
		r_43 := out
		r_43.startLines = this_.startLines
	}
	{
		r_44 := out
		r_44.startCols = this_.startCols
	}
	{
		r_45 := out
		r_45.endLines = eLines
	}
	{
		r_46 := out
		r_46.endCols = eCols
	}
	{
		r_47 := out
		r_47.mappingCount = count
	}
}

// Java method: MappedCodePrinter.reportLineCut(int lineIndex, int charIndex)
func (this_ *CodePrinter_MappedCodePrinterBase) ReportLineCut(lineIndex int32, charIndex int32) {
	_ = lineIndex
	_ = charIndex
	if this_.createSrcMap && (this_.mappingCount > int32(0)) {
		var count int32 = this_.mappingCount
		_ = count
		var sLines []int32 = this_.startLines
		_ = sLines
		var sCols []int32 = this_.startCols
		_ = sCols
		{
			var i int32 = this_.firstCandidateMappingForCut
			_ = i
			for ; i < count; i = i + 1 {
				if (sLines[i] == lineIndex) && (sCols[i] >= charIndex) {
					{
						a_48 := sLines
						i_49 := i
						a_48[i_49] = (lineIndex + int32(1))
					}
					a_50 := sCols
					i_51 := i
					a_50[i_51] = (a_50[i_51] - charIndex)
				}
			}
		}
		this_.firstCandidateMappingForCut = count
		var compCount int32 = this_.completeMappingCount
		_ = compCount
		if compCount > int32(0) {
			var compIds []int32 = this_.completeMappingIds
			_ = compIds
			var eLines []int32 = this_.endLines
			_ = eLines
			var eCols []int32 = this_.endCols
			_ = eCols
			{
				var i_2 int32 = int32(0)
				_ = i_2
				for ; i_2 < compCount; i_2 = i_2 + 1 {
					var id int32 = compIds[i_2]
					_ = id
					if (eLines[id] == lineIndex) && (eCols[id] >= charIndex) {
						{
							a_52 := eLines
							i_53 := id
							a_52[i_53] = (lineIndex + int32(1))
						}
						a_54 := eCols
						i_55 := id
						a_54[i_55] = (a_54[i_55] - charIndex)
					}
				}
			}
			this_.completeMappingCount = int32(0)
		}
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
func CodePrinter_MappedCodePrinter_ComputeLineLengths(code string) []int32 {
	clinit_CodePrinter_MappedCodePrinter()
	_ = code
	var lineEndPos int32 = j2gort.JStrIndexOfCp(code, int32(uint16(10)))
	_ = lineEndPos
	if lineEndPos == int32(-1) {
		return StaticGet_CodePrinter_MappedCodePrinter_EMPTY_INT_ARRAY()
	}
	var lengths []int32 = make([]int32, int32(64))
	_ = lengths
	var count int32 = int32(0)
	_ = count
	var lineStartPos int32 = int32(0)
	_ = lineStartPos
	for lineEndPos > int32(-1) {
		if count == int32(len(lengths)) {
			lengths = func() []int32 {
				arrSrc_56 := lengths
				arrLen_57 := (count * int32(2))
				j2gort.CheckNotNullSlice(arrSrc_56)
				if arrLen_57 < 0 {
					panic(j2gort.ConstructShimThrowable("java.lang.NegativeArraySizeException"))
				}
				arrOut_58 := make([]int32, arrLen_57)
				copy(arrOut_58, arrSrc_56)
				return arrOut_58
			}()
		}
		{
			a_59 := lengths
			i_60 := func() int32 { old_61 := count; nw_62 := old_61 + 1; count = nw_62; return old_61 }()
			a_59[i_60] = (lineEndPos - lineStartPos)
		}
		lineStartPos = (lineEndPos + int32(1))
		lineEndPos = j2gort.JStrIndexOfCpFrom(code, int32(uint16(10)), lineStartPos)
	}
	if count == int32(len(lengths)) {
		return lengths
	} else {
		return func() []int32 {
			arrSrc_63 := lengths
			arrLen_64 := count
			j2gort.CheckNotNullSlice(arrSrc_63)
			if arrLen_64 < 0 {
				panic(j2gort.ConstructShimThrowable("java.lang.NegativeArraySizeException"))
			}
			arrOut_65 := make([]int32, arrLen_64)
			copy(arrOut_65, arrSrc_63)
			return arrOut_65
		}()
	}
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
	var didRet_66 bool
	func() {
		defer func() {
			if rec_67 := recover(); rec_67 != nil {
				jex_68 := javalang.RecoverException(rec_67)
				_ = jex_68
				if j2gort.InstanceOf[javalang.INumberFormatException](jex_68) {
					e := j2gort.AsIface[javalang.INumberFormatException](jex_68)
					_ = e
					this_.CodePrinter_MappedCodePrinterBase.AddNumber(x, n)
					didRet_66 = true
					return
				} else {
					panic(rec_67)
				}
			}
		}()
		d = j2gort.ParseDouble(numberFromSource)
	}()
	if didRet_66 {
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
	var ok_69 bool
	src, ok_69 = j2gort.TryCast[*SourceFile](staticSrc)
	if !ok_69 {
		return j2gort.NullString
	}
	var srcCode string = j2gort.NullString
	_ = srcCode
	var didRet_70 bool
	var retVal_71 string
	func() {
		defer func() {
			if rec_72 := recover(); rec_72 != nil {
				jex_73 := javalang.RecoverException(rec_72)
				_ = jex_73
				if j2gort.InstanceOf[javaio.IIOException](jex_73) {
					e := j2gort.AsIface[javaio.IIOException](jex_73)
					_ = e
					retVal_71 = j2gort.NullString
					didRet_70 = true
					return
				} else {
					panic(rec_72)
				}
			}
		}()
		if src.IsStubSourceFileForAlreadyProvidedInput() {
			retVal_71 = j2gort.NullString
			didRet_70 = true
			return
		}
		srcCode = src.GetCode()
	}()
	if didRet_70 {
		return retVal_71
	}
	var offset int32 = int32(0)
	_ = offset
	var didRet_74 bool
	var retVal_75 string
	func() {
		defer func() {
			if rec_76 := recover(); rec_76 != nil {
				jex_77 := javalang.RecoverException(rec_76)
				_ = jex_77
				if j2gort.InstanceOf[javalang.IIllegalArgumentException](jex_77) {
					e_2 := j2gort.AsIface[javalang.IIllegalArgumentException](jex_77)
					_ = e_2
					retVal_75 = j2gort.NullString
					didRet_74 = true
					return
				} else {
					panic(rec_76)
				}
			}
		}()
		offset = n.GetSourceOffset()
	}()
	if didRet_74 {
		return retVal_75
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

func Upcast_CodePrinter_LicenseTracker[T CodePrinter_LicenseTracker](v T) CodePrinter_LicenseTracker {
	var vi CodePrinter_LicenseTracker = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
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
func (this_ *CodePrinter_Builder) priv_Builder_getOptions_1363() *CompilerOptions {
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
	return this_.priv_Builder_buildInternal_1375(false).source
}

// Java method: Builder.buildWithSourceMappings()
func (this_ *CodePrinter_Builder) BuildWithSourceMappings() *CodePrinter_SourceAndMappings {
	return this_.priv_Builder_buildInternal_1375(true)
}

// Java method: Builder.buildInternal(boolean allowSourceMap)
func (this_ *CodePrinter_Builder) priv_Builder_buildInternal_1375(allowSourceMap bool) *CodePrinter_SourceAndMappings {
	_ = allowSourceMap
	if this_.root == 0 {
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Cannot build without root node being specified")))
	}
	var opts *CompilerOptions = this_.priv_Builder_getOptions_1363()
	_ = opts
	return CodePrinter_ToSource(this_.root, CodePrinter_Format_FromOptions(opts, this_.outputTypes, this_.prettyPrint), opts, (allowSourceMap && opts.ShouldGatherSourceMapInfo()), this_.licenseTracker, this_.tagAsTypeSummary, this_.tagAsStrict, this_.lineBreak, this_.codeGeneratorFactory)
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

func Upcast_CodePrinter_Builder_CodeGeneratorFactory[T CodePrinter_Builder_CodeGeneratorFactory](v T) CodePrinter_Builder_CodeGeneratorFactory {
	var vi CodePrinter_Builder_CodeGeneratorFactory = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
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
	source       string           `gson:"source"`
	mappingNodes []pkg_rhino.Node `gson:"mappingNodes"`
	startLines   []int32          `gson:"startLines"`
	startCols    []int32          `gson:"startCols"`
	endLines     []int32          `gson:"endLines"`
	endCols      []int32          `gson:"endCols"`
	mappingCount int32            `gson:"mappingCount"`
	idHash_      atomic.Int32
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
	var opts *CompilerOptions = this_.outer_.priv_Builder_getOptions_1363()
	_ = opts
	if outputFormat == StaticGet_CodePrinter_Format_TYPED() {
		return NewTypedCodeGenerator(cc, opts, this_.outer_.registry)
	} else {
		return NewCodeGenerator_2(cc, opts)
	}
}
