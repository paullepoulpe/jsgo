// Port of com/google/javascript/jscomp/CodePrinter.java @ closure-compiler v20261021

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

package jscomp

import (
	"math"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/sourcemap"
)

// CodePrinter_PrettyCodePrinter_INDENT is PrettyCodePrinter.INDENT.
const CodePrinter_PrettyCodePrinter_INDENT = "  "

type prettyNumberSource struct {
	code  string
	ascii bool
	units []uint16
}

const initialBufferCap = 512

type codePrinterInstance struct {
	cg                  CodeGenerator
	prettySrc           prettyNumberSource
	codeBuf             [initialBufferCap]byte
	openMappingsBuf     [64]int
	completeMappingsBuf [16]int
}

type codePrinterScratch struct {
	inst *codePrinterInstance
}

func (c *CodeConsumer) initPrinter(
	mode codeSinkMode,
	lineBreak bool,
	lineLengthThreshold int,
	createSrcMap bool,
	sourceMapDetailLevel SourceMap_DetailLevel,
	licenseTracker CodePrinter_LicenseTracker,
	codeBuf []byte,
	openBuf []int,
	completeBuf []int,
	prettySrc *prettyNumberSource,
) {
	keepCode := c.code
	if cap(keepCode) == 0 {
		keepCode = codeBuf
	} else {
		keepCode = keepCode[:0]
	}
	thresh := lineLengthThreshold
	if thresh <= 0 {
		thresh = math.MaxInt32
	}
	var openMappings, completeMappings []int
	if createSrcMap {
		openMappings = openBuf
		completeMappings = completeBuf
	}
	*c = CodeConsumer{
		mode:                 mode,
		continueProcessing:   true,
		lineBreak:            lineBreak,
		createSrcMap:         createSrcMap,
		code:                 keepCode,
		lineLengthThreshold:  thresh,
		sourceMapDetailLevel: sourceMapDetailLevel,
		licenseTracker:       licenseTracker,
		openMappings:         openMappings,
		completeMappings:     completeMappings,
		prettySrc:            prettySrc,
	}
}

func (inst *codePrinterInstance) resetAfterPrint() {
	var keepCode []byte
	if cap(inst.cg.cc.code) <= 1<<20 {
		keepCode = inst.cg.cc.code[:0]
	}
	var keepScratch []uint16
	if cap(inst.cg.scratchUnits) <= 65536 {
		keepScratch = inst.cg.scratchUnits[:0]
	}
	inst.cg = CodeGenerator{
		cc:           CodeConsumer{code: keepCode},
		scratchUnits: keepScratch,
	}
	inst.prettySrc = prettyNumberSource{}
}

func (c *CodeConsumer) adjustSourceMappings(code string) {
	if !c.createSrcMap || len(c.allMappings) == 0 || code == "" {
		return
	}
	lineEndPos := strings.IndexByte(code, '\n')
	if lineEndPos == -1 {
		return
	}
	var stackBuf [64]int
	var lineLengths []int
	if c.lineIndex > len(stackBuf) {
		lineLengths = make([]int, 0, c.lineIndex)
	} else {
		lineLengths = stackBuf[:0]
	}
	asciiKnown := !c.hasNonASCII && len(code) == len(c.code)
	lineStartPos := 0
	for lineEndPos > -1 {
		if asciiKnown {
			lineLengths = append(lineLengths, lineEndPos-lineStartPos)
		} else {
			lineLengths = append(lineLengths, jutil.StringLength(code[lineStartPos:lineEndPos]))
		}
		lineStartPos = lineEndPos + 1
		lineEndPos = strings.IndexByte(code[lineStartPos:], '\n')
		if lineEndPos > -1 {
			lineEndPos += lineStartPos
		}
	}
	for i := range c.allMappings {
		c.allMappings[i].end = adjustEndPosition(lineLengths, c.allMappings[i].end)
	}
}

// ReportLineCut reports to the code consumer that the given line has been
// cut at the given position, i.e. a \n has been inserted there. All mappings
// in the source maps after that position will be renormalized as needed.
//
// Java: CodePrinter.MappedCodePrinter.reportLineCut(int, int)
func (c *CodeConsumer) ReportLineCut(lineIndex int, charIndex int) {
	if c.createSrcMap {
		for i := c.cutIdx; i < len(c.allMappings); i++ {
			mapping := &c.allMappings[i]
			mapping.start = convertPositionAfterLineCut(mapping.start, lineIndex, charIndex)
			if mapping.end.GetLine() >= 0 {
				mapping.end = convertPositionAfterLineCut(mapping.end, lineIndex, charIndex)
			}
		}
		c.cutIdx = len(c.allMappings)

		for _, idx := range c.completeMappings {
			c.allMappings[idx].end = convertPositionAfterLineCut(c.allMappings[idx].end, lineIndex, charIndex)
		}
		c.completeMappings = c.completeMappings[:0]
		c.cutStackDepth = len(c.openMappings)
	}
}

func convertPositionAfterLineCut(
	position sourcemap.FilePosition, lineIndex int, characterPosition int) sourcemap.FilePosition {
	originalLine := position.GetLine()
	originalChar := position.GetColumn()
	if originalLine == lineIndex && originalChar >= characterPosition {
		return sourcemap.MakeFilePosition(originalLine+1, originalChar-characterPosition)
	}
	return position
}

// Java: CodePrinter.MappedCodePrinter.getCode()
func (c *CodeConsumer) GetCode() string {
	if len(c.code) == 0 {
		return ""
	}
	return string(c.code)
}

// Java: CodePrinter.MappedCodePrinter.getCurrentCharIndex()
func (c *CodeConsumer) GetCurrentCharIndex() int {
	return c.lineLength
}

// Java: CodePrinter.MappedCodePrinter.getCurrentLineIndex()
func (c *CodeConsumer) GetCurrentLineIndex() int {
	return c.lineIndex
}

func adjustEndPosition(
	lineLengths []int, endPosition sourcemap.FilePosition) sourcemap.FilePosition {
	line := endPosition.GetLine()
	if line >= len(lineLengths) {
		return endPosition
	}
	col := endPosition.GetColumn()
	if col > lineLengths[line] {
		codePrinter_panicBadEndPosition(endPosition, lineLengths[line])
	}
	if col == lineLengths[line] {
		return sourcemap.MakeFilePosition(line+1, 0)
	}
	return endPosition
}

//go:noinline
func codePrinter_panicBadEndPosition(endPosition sourcemap.FilePosition, lineLen int) {
	jutil.CheckStateMsg2(
		false,
		"End position %s points to a column larger than line length %s",
		&endPosition,
		lineLen)
}

// getTryForCatch returns the TRY node for the specified CATCH node.
//
// Java: CodePrinter.PrettyCodePrinter.getTryForCatch(Node)
func getTryForCatch(n rhino.Node) rhino.Node {
	return n.GetGrandparent()
}

func codePrinter_parseDouble(s string) (float64, bool) {
	return jutil.TryParseDouble(s)
}

// Java: CodePrinter.PrettyCodePrinter.getNumberFromSource(Node)
func (c *CodeConsumer) GetNumberFromSource(n rhino.Node) (string, bool) {
	if !n.IsNumber() {
		return "", false
	}

	staticSrc := NodeUtil_GetSourceFile(n)
	src, isSourceFile := staticSrc.(*SourceFile)
	if !isSourceFile || src == nil || src.IsStubSourceFileForAlreadyProvidedInput() {
		return "", false
	}
	srcCode, err := src.getCode()
	if err != nil {
		return "", false
	}

	offset, ok := codePrinter_getSourceOffset(n)
	if !ok {
		return "", false
	}
	if c.prettySrc == nil {
		c.prettySrc = &prettyNumberSource{}
	}
	endOffset := offset + n.GetLength()
	if offset < 0 || endOffset > c.prettySrc.length(srcCode) {
		return "", false
	}

	return c.prettySrc.substring(offset, endOffset), true
}

func (p *prettyNumberSource) length(srcCode string) int {
	if srcCode != p.code {
		p.code = srcCode
		p.ascii = jutil.IsASCII(srcCode)
		if p.ascii {
			p.units = nil
		} else {
			p.units = jutil.StringToUTF16(srcCode)
		}
	}
	if p.ascii {
		return len(p.code)
	}
	return len(p.units)
}

func (p *prettyNumberSource) substring(begin int, end int) string {
	if begin > end {
		panic(jutil.NewStringIndexOutOfBoundsExceptionForRange(begin, end, p.length(p.code)))
	}
	if p.ascii {
		return p.code[begin:end]
	}
	return jutil.StringFromUTF16(p.units[begin:end])
}

func codePrinter_getSourceOffset(n rhino.Node) (int, bool) {
	src, ok := n.GetStaticSourceFile().(*SourceFile)
	if !ok || src == nil {
		return 0, false
	}
	lineno := n.GetLineno()
	if lineno < 1 {
		return 0, false
	}
	src.FindLineOffsets()
	if lineno > len(src.lineOffsets) {
		return 0, false
	}
	offset := int(int32(src.lineOffsets[lineno-1]) + int32(n.GetCharno()))
	if offset < 0 {
		return 0, false
	}
	return offset, true
}

// CodePrinter_LicenseTracker: License Trackers are responsible for ensuring
// that any licensing information attached to nodes is retained in the final
// output of JSCompiler.
type CodePrinter_LicenseTracker interface {
	// Java: CodePrinter.LicenseTracker.trackLicensesForNode(Node)
	TrackLicensesForNode(node rhino.Node)

	// Java: CodePrinter.LicenseTracker.emitLicenses()
	EmitLicenses() []string
}

// CodePrinter_Builder builds the source code of an AST.
type CodePrinter_Builder struct {
	root             rhino.Node
	options          *CompilerOptions
	lineBreak        bool
	prettyPrint      bool
	outputTypes      bool
	tagAsTypeSummary bool
	tagAsStrict      bool
	licenseTracker   CodePrinter_LicenseTracker // nil is null
	registry         *rhino.JSTypeRegistry      // may be null unless using Format.TYPED
}

// NewCodePrinter_Builder sets the root node from which to generate the
// source code.
//
// Java: CodePrinter.Builder(Node)
func NewCodePrinter_Builder(node rhino.Node) *CodePrinter_Builder {
	return &CodePrinter_Builder{root: node}
}

func (b *CodePrinter_Builder) getOptions() *CompilerOptions {
	if b.options == nil {
		b.options = NewCompilerOptions()
	}
	return b.options
}

// SetCompilerOptions sets the output options from compiler options.
//
// Java: CodePrinter.Builder.setCompilerOptions(CompilerOptions)
func (b *CodePrinter_Builder) SetCompilerOptions(options *CompilerOptions) *CodePrinter_Builder {
	b.options = options
	b.prettyPrint = options.IsPrettyPrint()
	b.lineBreak = options.ShouldAddLineBreak()
	return b
}

// Java: CodePrinter.Builder.setTypeRegistry(JSTypeRegistry)
func (b *CodePrinter_Builder) SetTypeRegistry(registry *rhino.JSTypeRegistry) *CodePrinter_Builder {
	b.registry = registry
	return b
}

// SetPrettyPrint sets whether pretty printing should be used.
//
// Java: CodePrinter.Builder.setPrettyPrint(boolean)
func (b *CodePrinter_Builder) SetPrettyPrint(prettyPrint bool) *CodePrinter_Builder {
	b.prettyPrint = prettyPrint
	return b
}

// SetLineBreak sets whether line breaking should be done automatically.
//
// Java: CodePrinter.Builder.setLineBreak(boolean)
func (b *CodePrinter_Builder) SetLineBreak(lineBreak bool) *CodePrinter_Builder {
	b.lineBreak = lineBreak
	return b
}

// SetOutputTypes sets whether to output closure-style type annotations.
//
// Java: CodePrinter.Builder.setOutputTypes(boolean)
func (b *CodePrinter_Builder) SetOutputTypes(outputTypes bool) *CodePrinter_Builder {
	b.outputTypes = outputTypes
	return b
}

// SetLicenseTracker sets the license tracker to use when printing.
//
// Java: CodePrinter.Builder.setLicenseTracker(LicenseTracker)
func (b *CodePrinter_Builder) SetLicenseTracker(licenseTracker CodePrinter_LicenseTracker) *CodePrinter_Builder {
	b.licenseTracker = licenseTracker
	return b
}

// SetTagAsTypeSummary sets whether the output should be tagged as an .i.js
// file.
//
// Java: CodePrinter.Builder.setTagAsTypeSummary(boolean)
func (b *CodePrinter_Builder) SetTagAsTypeSummary(tagAsTypeSummary bool) *CodePrinter_Builder {
	b.tagAsTypeSummary = tagAsTypeSummary
	return b
}

// SetTagAsStrict sets whether the output should be tags as ECMASCRIPT 5
// Strict.
//
// Java: CodePrinter.Builder.setTagAsStrict(boolean)
func (b *CodePrinter_Builder) SetTagAsStrict(tagAsStrict bool) *CodePrinter_Builder {
	b.tagAsStrict = tagAsStrict
	return b
}

// Build generates the source code and returns it.
//
// Java: CodePrinter.Builder.build()
func (b *CodePrinter_Builder) Build() string {
	return b.buildSourceAndMappingsInto(nil, false, nil)
}

// Java: CodePrinter.Builder.buildWithSourceMappings()
func (b *CodePrinter_Builder) BuildWithSourceMappings() *CodePrinter_SourceAndMappings {
	sm := &CodePrinter_SourceAndMappings{}
	b.buildSourceAndMappingsInto(sm, true, nil)
	return sm
}

func (b *CodePrinter_Builder) buildSourceAndMappingsInto(
	sm *CodePrinter_SourceAndMappings,
	wantMappings bool,
	scratch *codePrinterScratch) string {
	if b.root.IsNil() {
		panic(jutil.NewIllegalStateExceptionMsg("Cannot build without root node being specified"))
	}
	options := b.getOptions()
	outputFormat := CodePrinter_Format_FromOptions(options, b.outputTypes, b.prettyPrint)
	return b.buildSourceAndMappingsWithFormatInto(options, outputFormat, sm, wantMappings, scratch)
}

func (b *CodePrinter_Builder) buildSourceAndMappingsWithFormatInto(
	options *CompilerOptions,
	outputFormat CodePrinter_Format,
	sm *CodePrinter_SourceAndMappings,
	wantMappings bool,
	scratch *codePrinterScratch) string {
	gatherSourceMap := wantMappings && options.ShouldGatherSourceMapInfo()

	var inst *codePrinterInstance
	if scratch != nil {
		if scratch.inst == nil {
			scratch.inst = &codePrinterInstance{}
		}
		inst = scratch.inst
	} else {
		var localInst codePrinterInstance
		inst = &localInst
	}

	mode := sinkPretty
	var prettySrc *prettyNumberSource
	if outputFormat == CodePrinter_Format_COMPACT {
		mode = sinkCompact
	} else {
		prettySrc = &inst.prettySrc
	}

	cg := &inst.cg
	cc := &cg.cc
	cc.initPrinter(
		mode,
		b.lineBreak,
		options.GetLineLengthThreshold(),
		gatherSourceMap,
		options.GetSourceMapDetailLevel(),
		b.licenseTracker,
		inst.codeBuf[:0:len(inst.codeBuf)],
		inst.openMappingsBuf[:0:len(inst.openMappingsBuf)],
		inst.completeMappingsBuf[:0:len(inst.completeMappingsBuf)],
		prettySrc,
	)

	if sm != nil && gatherSourceMap {
		if cap(sm.mappings) == 0 {
			cc.allMappings = make([]SourceMap_Mapping, 0, 64)
		} else {
			cc.allMappings = sm.mappings[:0]
		}
	}

	if outputFormat == CodePrinter_Format_TYPED {
		cg.initTypedWithOptions(options, b.registry)
	} else {
		cg.initWithOptions(options)
	}

	if b.tagAsTypeSummary {
		cg.TagAsTypeSummary()
	}
	if b.tagAsStrict {
		cg.TagAsStrict()
	}

	cg.AddNode(b.root)
	cc.EndFile()

	source := cc.GetCode()
	if sm != nil {
		sm.source = source
		if !cc.hasNonASCII {
			sm.utf16Len = len(source)
			sm.lineIndex = cc.lineIndex
			sm.lineLength = cc.lineLength
			sm.hasMetrics = true
		} else {
			sm.hasMetrics = false
		}
		if gatherSourceMap {
			cc.adjustSourceMappings(source)
			sm.mappings = cc.allMappings
		} else {
			sm.mappings = sm.mappings[:0]
		}
	}
	if scratch != nil {
		inst.resetAfterPrint()
	}
	return source
}

// CodePrinter_Format specifies a format for code generation.
//
// Java: CodePrinter.Format
type CodePrinter_Format int32

const (
	CodePrinter_Format_COMPACT CodePrinter_Format = iota
	CodePrinter_Format_PRETTY
	CodePrinter_Format_TYPED
)

var codePrinter_Format_names = [...]string{"COMPACT", "PRETTY", "TYPED"}

// String is Format.name().
func (f CodePrinter_Format) String() string { return codePrinter_Format_names[f] }

// Java: CodePrinter.Format.fromOptions(CompilerOptions, boolean, boolean)
func CodePrinter_Format_FromOptions(options *CompilerOptions, outputTypes bool, prettyPrint bool) CodePrinter_Format {
	if outputTypes {
		return CodePrinter_Format_TYPED
	}
	if prettyPrint || options.GetOutputFeatureSet().ContainsFeature(parser.FeatureSet_Feature_TYPE_ANNOTATION) {
		return CodePrinter_Format_PRETTY
	}
	return CodePrinter_Format_COMPACT
}

// CodePrinter_SourceAndMappings bundles together the source and generated
// SourceMap Mappings for that source.
type CodePrinter_SourceAndMappings struct {
	source     string
	mappings   []SourceMap_Mapping
	utf16Len   int
	lineIndex  int
	lineLength int
	hasMetrics bool
}

// Java: CodePrinter.SourceAndMappings()
func NewCodePrinter_SourceAndMappings() *CodePrinter_SourceAndMappings {
	return &CodePrinter_SourceAndMappings{}
}
