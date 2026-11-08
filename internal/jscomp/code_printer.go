// Port of com/google/javascript/jscomp/CodePrinter.java @ closure-compiler v20261104

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
	"bytes"
	"math"
	"unicode/utf8"

	"github.com/paullepoulpe/jsgo/internal/ast"
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
	mappingsBuf         []SourceMap_Mapping
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
	if cap(inst.mappingsBuf) > 16384 {
		inst.mappingsBuf = nil
	} else {
		inst.mappingsBuf = inst.mappingsBuf[:0]
	}
	inst.cg = CodeGenerator{
		cc:           CodeConsumer{code: keepCode},
		scratchUnits: keepScratch,
	}
	inst.prettySrc = prettyNumberSource{}
}

func asciiPrefixLenBytes(b []byte) int {
	i := 0
	for ; i+8 <= len(b); i += 8 {
		if (b[i] | b[i+1] | b[i+2] | b[i+3] | b[i+4] | b[i+5] | b[i+6] | b[i+7]) >= 0x80 {
			break
		}
	}
	for ; i < len(b); i++ {
		if b[i] >= 0x80 {
			return i
		}
	}
	return len(b)
}

func decodeWTF8Bytes(b []byte) (r rune, size int) {
	c := b[0]
	if c < 0x80 {
		return rune(c), 1
	}
	// Generalized UTF-8 surrogate: ED A0..BF xx
	if c == 0xED && len(b) >= 3 && b[1] >= 0xA0 && b[1] <= 0xBF && b[2]&0xC0 == 0x80 {
		return rune(c&0x0F)<<12 | rune(b[1]&0x3F)<<6 | rune(b[2]&0x3F), 3
	}
	return utf8.DecodeRune(b)
}

func utf16LenBytes(b []byte) int {
	prefix := asciiPrefixLenBytes(b)
	if prefix == len(b) {
		return len(b)
	}
	n := prefix
	for i := prefix; i < len(b); {
		r, size := decodeWTF8Bytes(b[i:])
		i += size
		if r >= 0x10000 {
			n += 2
		} else {
			n++
		}
	}
	return n
}

func (c *CodeConsumer) adjustSourceMappingsBytes(code []byte) {
	if !c.createSrcMap || len(c.allMappings) == 0 || len(code) == 0 {
		return
	}
	lineEndPos := bytes.IndexByte(code, '\n')
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
			lineLengths = append(lineLengths, utf16LenBytes(code[lineStartPos:lineEndPos]))
		}
		lineStartPos = lineEndPos + 1
		lineEndPos = bytes.IndexByte(code[lineStartPos:], '\n')
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
		return sourcemap.NewFilePosition(originalLine+1, originalChar-characterPosition)
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
		return sourcemap.NewFilePosition(line+1, 0)
	}
	return endPosition
}

//go:noinline
func codePrinter_panicBadEndPosition(endPosition sourcemap.FilePosition, lineLen int) {
	panic(jutil.NewIllegalStateExceptionMsg(
		jutil.LenientFormat(
			"End position %s points to a column larger than line length %s",
			&endPosition,
			lineLen)))
}

// getTryForCatch returns the TRY node for the specified CATCH node.
//
// Java: CodePrinter.PrettyCodePrinter.getTryForCatch(Node)
func getTryForCatch(c ast.Catch) ast.Try {
	return ast.MustTry(c.Node().GetGrandparent())
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

var defaultCodePrintOptions = &CompilerOptions{
	languageIn:          CompilerOptions_LanguageMode_STABLE_IN,
	lineLengthThreshold: CompilerOptions_DEFAULT_LINE_LENGTH_THRESHOLD,
}

// CodePrintConfig configures AST source code generation for PrintCode and
// PrintCodeWithMappings.
//
// Java: CodePrinter.Builder
type CodePrintConfig struct {
	Root             rhino.Node
	Options          *CompilerOptions
	PrettyPrint      bool
	LineBreak        bool
	OutputTypes      bool
	TagAsTypeSummary bool
	TagAsStrict      bool
	LicenseTracker   CodePrinter_LicenseTracker // nil is null
	Registry         *rhino.JSTypeRegistry      // may be nil unless OutputTypes is true
}

// PrintCode generates the source code for cfg.Root and returns it.
//
// Java: CodePrinter.Builder.build()
func PrintCode(cfg CodePrintConfig) string {
	return printCodeInto(cfg, nil, false, nil)
}

// PrintCodeWithMappings generates the source code and source mappings for cfg.Root.
//
// Java: CodePrinter.Builder.buildWithSourceMappings()
func PrintCodeWithMappings(cfg CodePrintConfig) *CodePrinter_SourceAndMappings {
	sm := &CodePrinter_SourceAndMappings{}
	printCodeInto(cfg, sm, true, nil)
	return sm
}

func printCodeInto(
	cfg CodePrintConfig,
	sm *CodePrinter_SourceAndMappings,
	wantMappings bool,
	scratch *codePrinterScratch) string {
	if cfg.Root.IsNil() {
		panic(jutil.NewIllegalStateExceptionMsg("Cannot build without root node being specified"))
	}
	options := cfg.Options
	if options == nil {
		options = defaultCodePrintOptions
	}
	prettyPrint := cfg.PrettyPrint || options.IsPrettyPrint()
	lineBreak := cfg.LineBreak || options.ShouldAddLineBreak()
	outputFormat := CodePrinter_Format_FromOptions(options, cfg.OutputTypes, prettyPrint)
	return printCodeWithFormatInto(cfg, options, outputFormat, lineBreak, sm, wantMappings, scratch)
}

func printCodeWithFormatInto(
	cfg CodePrintConfig,
	options *CompilerOptions,
	outputFormat CodePrinter_Format,
	lineBreak bool,
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
	useSlotBufs := sm != nil && sm.ownCodeBuf
	var savedInstCode []byte
	initCodeBuf := inst.codeBuf[:0:len(inst.codeBuf)]
	if useSlotBufs {
		savedInstCode = cc.code
		cc.code = nil
		if cap(sm.codeBytes) == 0 {
			initCodeBuf = make([]byte, 0, initialBufferCap)
		} else {
			initCodeBuf = sm.codeBytes[:0]
		}
	}
	cc.initPrinter(
		mode,
		lineBreak,
		options.GetLineLengthThreshold(),
		gatherSourceMap,
		options.GetSourceMapDetailLevel(),
		cfg.LicenseTracker,
		initCodeBuf,
		inst.openMappingsBuf[:0:len(inst.openMappingsBuf)],
		inst.completeMappingsBuf[:0:len(inst.completeMappingsBuf)],
		prettySrc,
	)

	if sm != nil && gatherSourceMap {
		if cap(sm.mappings) > 0 {
			cc.allMappings = sm.mappings[:0]
		} else if !useSlotBufs && scratch != nil && cap(inst.mappingsBuf) > 0 {
			cc.allMappings = inst.mappingsBuf[:0]
		} else {
			cc.allMappings = make([]SourceMap_Mapping, 0, 64)
		}
	}

	if outputFormat == CodePrinter_Format_TYPED {
		cg.initTypedWithOptions(options, cfg.Registry)
	} else {
		cg.initWithOptions(options)
	}

	if cfg.TagAsTypeSummary {
		cg.TagAsTypeSummary()
	}
	if cfg.TagAsStrict {
		cg.TagAsStrict()
	}

	cg.AddNode(cfg.Root)
	cc.EndFile()

	var source string
	if sm == nil || !sm.bytesOnly {
		source = cc.GetCode()
	}
	if sm != nil {
		if sm.bytesOnly {
			sm.source = ""
			sm.codeBytes = cc.code
		} else {
			sm.source = source
			sm.codeBytes = nil
		}
		if !cc.hasNonASCII {
			sm.utf16Len = len(cc.code)
			sm.lineIndex = cc.lineIndex
			sm.lineLength = cc.lineLength
			sm.hasMetrics = true
		} else {
			sm.hasMetrics = false
		}
		if gatherSourceMap {
			cc.adjustSourceMappingsBytes(cc.code)
			sm.mappings = cc.allMappings
			if !useSlotBufs && scratch != nil {
				inst.mappingsBuf = cc.allMappings
			}
		} else {
			sm.mappings = sm.mappings[:0]
		}
	}
	if useSlotBufs {
		cc.code = savedInstCode
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
	codeBytes  []byte
	mappings   []SourceMap_Mapping
	utf16Len   int
	lineIndex  int
	lineLength int
	hasMetrics bool
	bytesOnly  bool
	ownCodeBuf bool
}

// Java: CodePrinter.SourceAndMappings()
func NewCodePrinter_SourceAndMappings() *CodePrinter_SourceAndMappings {
	return &CodePrinter_SourceAndMappings{}
}
