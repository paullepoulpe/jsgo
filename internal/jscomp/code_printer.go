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
	"iter"
	"math"
	"strings"
	"sync"

	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/sourcemap"
)

// CodePrinter prints out JS code in either pretty format or compact format.
//
// See CodeGenerator.
type CodePrinter struct {
	// Go-only: language-adapter — distinct identity (Java gives every `new`
	// one; a zero-size struct would not).
	_ byte
}

// There are two separate CodeConsumers, one for pretty-printing and
// another for compact printing.

// There are two implementations because the CompactCodePrinter
// potentially has a very different implementation to the pretty
// version.

// CodePrinter_MappedCodePrinter is the abstract base of the two code
// consumers.
//
// Go-only: language-adapter — an abstract class with subclasses: this
// interface declares its methods beyond CodeConsumer's, and
// CodePrinter_MappedCodePrinterBase holds the fields and implementations.
type CodePrinter_MappedCodePrinter interface {
	CodeConsumer
	GetSourceMappings(code string) []*SourceMap_Mapping
	ReportLineCut(lineIndex int, charIndex int)
	GetCode() string
	GetCurrentCharIndex() int
	GetCurrentLineIndex() int

	GetCodePrinter_MappedCodePrinterBase() *CodePrinter_MappedCodePrinterBase
}

// CodePrinter_MappedCodePrinterBase holds MappedCodePrinter's fields and
// implementations.
type CodePrinter_MappedCodePrinterBase struct {
	CodeConsumerBase

	// Go-only: language-adapter — the ArrayDeque used as a stack
	// (push/peek/pop) is a slice whose top is its last element. nil is null.
	mappings []*SourceMap_Mapping
	// The ordered list of finalized mappings that started before the last line
	// break and ended since the last line break. See #reportLineCut.
	completeMappings     []*SourceMap_Mapping // nil is null
	mappingChunks        [][]SourceMap_Mapping
	mappingSlab          []SourceMap_Mapping
	mappingChunksBuf     [16][]SourceMap_Mapping
	mappingBuf           [32]SourceMap_Mapping
	mappingsBuf          [64]*SourceMap_Mapping
	completeMappingsBuf  [16]*SourceMap_Mapping
	cutChunkIdx          int
	cutElemIdx           int
	cutStackDepth        int
	createSrcMap         bool
	sourceMapDetailLevel SourceMap_DetailLevel
	licenseTracker       CodePrinter_LicenseTracker // nil is null
	codeASCIIBuf         [128]byte
	codeASCII            []byte
	codeUTF16            *jutil.StringBuilder
	lineLengthThreshold  int
	lineLength           int
	lineIndex            int
	cg                   CodeGeneratorBase
}

const codePrinter_mappingSlabSize = 512

var codePrinter_mappingSlabPool = sync.Pool{
	New: func() any {
		return new([codePrinter_mappingSlabSize]SourceMap_Mapping)
	},
}

func (m *CodePrinter_MappedCodePrinterBase) allocMapping() *SourceMap_Mapping {
	if len(m.mappingSlab) == cap(m.mappingSlab) {
		if len(m.mappingSlab) > 0 {
			m.mappingChunks = append(m.mappingChunks, m.mappingSlab)
		}
		arr := codePrinter_mappingSlabPool.Get().(*[codePrinter_mappingSlabSize]SourceMap_Mapping)
		m.mappingSlab = arr[:0:codePrinter_mappingSlabSize]
	}
	i := len(m.mappingSlab)
	m.mappingSlab = m.mappingSlab[:i+1]
	return &m.mappingSlab[i]
}

func (m *CodePrinter_MappedCodePrinterBase) releaseMappingSlabs() {
	for i, chunk := range m.mappingChunks {
		if cap(chunk) == codePrinter_mappingSlabSize {
			clear(chunk[:cap(chunk)])
			codePrinter_mappingSlabPool.Put((*[codePrinter_mappingSlabSize]SourceMap_Mapping)(chunk[:cap(chunk)]))
		}
		m.mappingChunks[i] = nil
	}
	m.mappingChunks = m.mappingChunks[:0]
	if cap(m.mappingSlab) == codePrinter_mappingSlabSize {
		clear(m.mappingSlab[:cap(m.mappingSlab)])
		codePrinter_mappingSlabPool.Put((*[codePrinter_mappingSlabSize]SourceMap_Mapping)(m.mappingSlab[:cap(m.mappingSlab)]))
	}
	m.mappingSlab = nil
}

// GetCodePrinter_MappedCodePrinterBase gives access to the base fields
// through the interface.
func (m *CodePrinter_MappedCodePrinterBase) GetCodePrinter_MappedCodePrinterBase() *CodePrinter_MappedCodePrinterBase {
	return m
}

// InitCodePrinter_MappedCodePrinterBase is the MappedCodePrinter(int,
// boolean, SourceMap.DetailLevel, LicenseTracker) constructor; self is the
// object being constructed.
//
// Java: CodePrinter.MappedCodePrinter(int, boolean, SourceMap.DetailLevel, LicenseTracker)
func (m *CodePrinter_MappedCodePrinterBase) InitCodePrinter_MappedCodePrinterBase(
	self CodePrinter_MappedCodePrinter,
	lineLengthThreshold int,
	createSrcMap bool,
	sourceMapDetailLevel SourceMap_DetailLevel,
	licenseTracker CodePrinter_LicenseTracker) {
	m.InitCodeConsumerBase(self)
	m.cutChunkIdx = 0
	m.cutElemIdx = 0
	m.cutStackDepth = 0
	if cap(m.codeASCII) == 0 {
		m.codeASCII = m.codeASCIIBuf[:0:len(m.codeASCIIBuf)]
	} else {
		m.codeASCII = m.codeASCII[:0]
	}
	m.codeUTF16 = nil
	m.lineLength = 0
	m.lineIndex = 0
	// checkState(sourceMapDetailLevel != null)
	//
	// Go-only: language-adapter — an enum value cannot be null in Go, so the
	// check always holds.
	if lineLengthThreshold <= 0 {
		m.lineLengthThreshold = math.MaxInt32
	} else {
		m.lineLengthThreshold = lineLengthThreshold
	}
	m.createSrcMap = createSrcMap
	m.sourceMapDetailLevel = sourceMapDetailLevel
	if createSrcMap {
		m.mappingChunks = m.mappingChunksBuf[:0:len(m.mappingChunksBuf)]
		m.mappingSlab = m.mappingBuf[:0:len(m.mappingBuf)]
		m.mappings = m.mappingsBuf[:0:len(m.mappingsBuf)]
		m.completeMappings = m.completeMappingsBuf[:0:len(m.completeMappingsBuf)]
	} else {
		m.mappingChunks = nil
		m.mappingSlab = nil
		m.mappings = nil
		m.completeMappings = nil
	}
	m.licenseTracker = licenseTracker
}

func (m *CodePrinter_MappedCodePrinterBase) resetForPool() {
	if m.createSrcMap {
		m.releaseMappingSlabs()
		clear(m.mappings)
		clear(m.mappingsBuf[:])
		m.mappings = nil
		clear(m.completeMappings)
		clear(m.completeMappingsBuf[:])
		m.completeMappings = nil
	}
	if cap(m.codeASCII) > 1<<20 {
		m.codeASCII = nil
	} else {
		m.codeASCII = m.codeASCII[:0]
	}
	m.codeUTF16 = nil
	m.licenseTracker = nil
	m.self = nil
	var keepScratch []uint16
	if cap(m.cg.scratchUnits) <= 65536 {
		keepScratch = m.cg.scratchUnits[:0]
	}
	m.cg = CodeGeneratorBase{scratchUnits: keepScratch}
}

func returnPooledPrinter(pooledCompact *CodePrinter_CompactCodePrinter, pooledPretty *CodePrinter_PrettyCodePrinter) {
	if pooledCompact != nil {
		pooledCompact.resetForPool()
		codePrinter_compactPool.Put(pooledCompact)
	} else if pooledPretty != nil {
		pooledPretty.resetForPool()
		pooledPretty.numberSourceCode = ""
		pooledPretty.numberSourceUnits = nil
		codePrinter_prettyPool.Put(pooledPretty)
	}
}

func (m *CodePrinter_MappedCodePrinterBase) codeGrowASCII(n int) {
	if need := len(m.codeASCII) + n; need > cap(m.codeASCII) {
		newCap := max(2*cap(m.codeASCII), max(need, 1024))
		buf := make([]byte, len(m.codeASCII), newCap)
		copy(buf, m.codeASCII)
		m.codeASCII = buf
	}
}

func (m *CodePrinter_MappedCodePrinterBase) codeAppend(str string) int {
	if m.codeUTF16 == nil {
		if jutil.IsASCII(str) {
			m.codeGrowASCII(len(str))
			m.codeASCII = append(m.codeASCII, str...)
			return len(str)
		}
		m.codeUTF16 = jutil.NewStringBuilderWithCapacity(max(len(m.codeASCII)+len(str), 1024))
		if len(m.codeASCII) > 0 {
			m.codeUTF16.Append(string(m.codeASCII))
		}
		m.codeASCII = m.codeASCII[:0]
	}
	before := m.codeUTF16.Length()
	m.codeUTF16.Append(str)
	return m.codeUTF16.Length() - before
}

func (m *CodePrinter_MappedCodePrinterBase) codeAppendASCII(str string) {
	if m.codeUTF16 == nil {
		m.codeGrowASCII(len(str))
		m.codeASCII = append(m.codeASCII, str...)
		return
	}
	m.codeUTF16.Append(str)
}

func (m *CodePrinter_MappedCodePrinterBase) codeAppendByte(c byte) {
	if m.codeUTF16 == nil {
		m.codeGrowASCII(1)
		m.codeASCII = append(m.codeASCII, c)
		return
	}
	m.codeUTF16.AppendChar(uint16(c))
}

func (m *CodePrinter_MappedCodePrinterBase) codeLength() int {
	if m.codeUTF16 == nil {
		return len(m.codeASCII)
	}
	return m.codeUTF16.Length()
}

func (m *CodePrinter_MappedCodePrinterBase) codeCharAt(i int) uint16 {
	if m.codeUTF16 == nil {
		if i < 0 || i >= len(m.codeASCII) {
			panic(jutil.NewStringIndexOutOfBoundsException(i, len(m.codeASCII)))
		}
		return uint16(m.codeASCII[i])
	}
	return m.codeUTF16.CharAt(i)
}

func (m *CodePrinter_MappedCodePrinterBase) codeInsertByte(offset int, c byte) {
	if m.codeUTF16 == nil {
		if offset < 0 || offset > len(m.codeASCII) {
			panic(jutil.NewStringIndexOutOfBoundsException(offset, len(m.codeASCII)))
		}
		n := len(m.codeASCII)
		m.codeGrowASCII(1)
		m.codeASCII = append(m.codeASCII, 0)
		copy(m.codeASCII[offset+1:], m.codeASCII[offset:n])
		m.codeASCII[offset] = c
		return
	}
	m.codeUTF16.InsertChar(offset, uint16(c))
}

// Append appends a string to the code, keeping track of the current line
// length.
//
// Java: CodePrinter.MappedCodePrinter.append(String)
func (m *CodePrinter_MappedCodePrinterBase) Append(str string) {
	// Go-only: language-adapter — str.length() is the number of UTF-16
	// units the append added, which codeAppend counts.
	m.lineLength += m.codeAppend(str)
}

// Java: CodePrinter.MappedCodePrinter.trackLicenses(Node)
func (m *CodePrinter_MappedCodePrinterBase) TrackLicenses(node *rhino.Node) {
	if m.licenseTracker != nil {
		m.licenseTracker.TrackLicensesForNode(node)
	}
}

// StartSourceMapping starts the source mapping for the given node at the
// current position.
//
// Java: CodePrinter.MappedCodePrinter.startSourceMapping(Node)
func (m *CodePrinter_MappedCodePrinterBase) StartSourceMapping(node *rhino.Node) {
	// checkState(sourceMapDetailLevel != null): see the constructor.
	jutil.CheckState(node != nil)
	// Go-only: language-adapter — `node.getSourceFileName() != null` is
	// `node.getStaticSourceFile() != null` (a Go file name is never null);
	// it avoids allocating the *string for every node.
	if m.createSrcMap &&
		node.GetLineno() > 0 &&
		(m.sourceMapDetailLevel == SourceMap_DetailLevel_ALL || m.sourceMapDetailLevel.Apply(node)) &&
		node.GetStaticSourceFile() != nil {
		line := m.GetCurrentLineIndex()
		index := m.GetCurrentCharIndex()
		jutil.CheckState(line >= 0)
		mapping := m.allocMapping()
		mapping.node = node
		mapping.start = sourcemap.MakeFilePosition(line, index)
		mapping.end = sourcemap.MakeFilePosition(-1, -1)
		m.mappings = append(m.mappings, mapping)
	}
}

// EndSourceMapping finishes the source mapping for the given node at the
// current position.
//
// Java: CodePrinter.MappedCodePrinter.endSourceMapping(Node)
func (m *CodePrinter_MappedCodePrinterBase) EndSourceMapping(node *rhino.Node) {
	if m.createSrcMap && len(m.mappings) != 0 && m.mappings[len(m.mappings)-1].node == node {
		mapping := m.mappings[len(m.mappings)-1]
		m.mappings[len(m.mappings)-1] = nil
		m.mappings = m.mappings[:len(m.mappings)-1]
		line := m.GetCurrentLineIndex()
		index := m.GetCurrentCharIndex()
		jutil.CheckState(line >= 0)
		mapping.end = sourcemap.MakeFilePosition(line, index)
		if len(m.mappings) < m.cutStackDepth {
			m.cutStackDepth = len(m.mappings)
			m.completeMappings = append(m.completeMappings, mapping)
		}
	}
}

func (m *CodePrinter_MappedCodePrinterBase) adjustSourceMappings(code string) {
	if !m.createSrcMap || (len(m.mappingChunks) == 0 && len(m.mappingSlab) == 0) || code == "" {
		return
	}
	lineEndPos := strings.IndexByte(code, '\n')
	if lineEndPos == -1 {
		return
	}
	var stackBuf [64]int
	var lineLengths []int
	if m.lineIndex > len(stackBuf) {
		lineLengths = make([]int, 0, m.lineIndex)
	} else {
		lineLengths = stackBuf[:0]
	}
	asciiKnown := m.codeUTF16 == nil && len(code) == len(m.codeASCII)
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
	for _, chunk := range m.mappingChunks {
		for i := range chunk {
			chunk[i].end = codePrinter_MappedCodePrinter_adjustEndPositionVal(lineLengths, chunk[i].end)
		}
	}
	for i := range m.mappingSlab {
		m.mappingSlab[i].end = codePrinter_MappedCodePrinter_adjustEndPositionVal(lineLengths, m.mappingSlab[i].end)
	}
}

func (m *CodePrinter_MappedCodePrinterBase) iterMappings(yield func(int, *SourceMap_Mapping) bool) {
	if !m.createSrcMap {
		return
	}
	idx := 0
	for _, chunk := range m.mappingChunks {
		for i := range chunk {
			if !yield(idx, &chunk[i]) {
				return
			}
			idx++
		}
	}
	for i := range m.mappingSlab {
		if !yield(idx, &m.mappingSlab[i]) {
			return
		}
		idx++
	}
	m.releaseMappingSlabs()
}

// GetSourceMappings returns a list of sourcemap mappings that were
// generated while printing this code. Only useful if createSrcMap was true
// for this MappedCodePrinter.
//
// Java: CodePrinter.MappedCodePrinter.getSourceMappings(String)
func (m *CodePrinter_MappedCodePrinterBase) GetSourceMappings(code string) []*SourceMap_Mapping {
	if !m.createSrcMap {
		return nil
	}
	m.adjustSourceMappings(code)
	total := len(m.mappingSlab)
	for _, chunk := range m.mappingChunks {
		total += len(chunk)
	}
	out := make([]*SourceMap_Mapping, 0, total)
	for _, chunk := range m.mappingChunks {
		for i := range chunk {
			out = append(out, &chunk[i])
		}
	}
	for i := range m.mappingSlab {
		out = append(out, &m.mappingSlab[i])
	}
	return out
}

// ReportLineCut reports to the code consumer that the given line has been
// cut at the given position, i.e. a \n has been inserted there. All mappings
// in the source maps after that position will be renormalized as needed.
//
// Java: CodePrinter.MappedCodePrinter.reportLineCut(int, int)
func (m *CodePrinter_MappedCodePrinterBase) ReportLineCut(lineIndex int, charIndex int) {
	if m.createSrcMap {
		for m.cutChunkIdx < len(m.mappingChunks) {
			chunk := m.mappingChunks[m.cutChunkIdx]
			for i := m.cutElemIdx; i < len(chunk); i++ {
				mapping := &chunk[i]
				mapping.start = codePrinter_MappedCodePrinter_convertPositionAfterLineCutVal(mapping.start, lineIndex, charIndex)
				if mapping.end.GetLine() >= 0 {
					mapping.end = codePrinter_MappedCodePrinter_convertPositionAfterLineCutVal(mapping.end, lineIndex, charIndex)
				}
			}
			m.cutChunkIdx++
			m.cutElemIdx = 0
		}
		for i := m.cutElemIdx; i < len(m.mappingSlab); i++ {
			mapping := &m.mappingSlab[i]
			mapping.start = codePrinter_MappedCodePrinter_convertPositionAfterLineCutVal(mapping.start, lineIndex, charIndex)
			if mapping.end.GetLine() >= 0 {
				mapping.end = codePrinter_MappedCodePrinter_convertPositionAfterLineCutVal(mapping.end, lineIndex, charIndex)
			}
		}
		m.cutElemIdx = len(m.mappingSlab)

		for _, mapping := range m.completeMappings {
			mapping.end = codePrinter_MappedCodePrinter_convertPositionAfterLineCutVal(mapping.end, lineIndex, charIndex)
		}
		clear(m.completeMappings)
		m.completeMappings = m.completeMappings[:0]
		m.cutStackDepth = len(m.mappings)
	}
}

func codePrinter_MappedCodePrinter_convertPositionAfterLineCutVal(
	position sourcemap.FilePosition, lineIndex int, characterPosition int) sourcemap.FilePosition {
	originalLine := position.GetLine()
	originalChar := position.GetColumn()
	if originalLine == lineIndex && originalChar >= characterPosition {
		return sourcemap.MakeFilePosition(originalLine+1, originalChar-characterPosition)
	}
	return position
}

// CodePrinter_MappedCodePrinter_ConvertPositionAfterLineCut converts the
// given position by normalizing it against the insertion at the given line
// and character position.
//
// position is the existing position before the newline was inserted;
// lineIndex is the index of the line at which the newline was inserted;
// characterPosition is the position on the line at which the newline was
// inserted. Returns the normalized position.
//
// Java: CodePrinter.MappedCodePrinter.convertPositionAfterLineCut(FilePosition, int, int)
func CodePrinter_MappedCodePrinter_ConvertPositionAfterLineCut(
	position *sourcemap.FilePosition, lineIndex int, characterPosition int) *sourcemap.FilePosition {
	originalLine := position.GetLine()
	originalChar := position.GetColumn()
	if originalLine == lineIndex && originalChar >= characterPosition {
		// If the position falls on the line itself, then normalize it
		// if it falls at or after the place the newline was inserted.
		return sourcemap.NewFilePosition(originalLine+1, originalChar-characterPosition)
	} else {
		return position
	}
}

// Java: CodePrinter.MappedCodePrinter.getCode()
func (m *CodePrinter_MappedCodePrinterBase) GetCode() string {
	if m.codeUTF16 != nil {
		return m.codeUTF16.String()
	}
	if len(m.codeASCII) == 0 {
		return ""
	}
	return string(m.codeASCII)
}

// Java: CodePrinter.MappedCodePrinter.getLastChar()
func (m *CodePrinter_MappedCodePrinterBase) GetLastChar() uint16 {
	if m.codeUTF16 != nil {
		if m.codeUTF16.Length() > 0 {
			return m.codeUTF16.CharAt(m.codeUTF16.Length() - 1)
		}
		return 0
	}
	if len(m.codeASCII) > 0 {
		return uint16(m.codeASCII[len(m.codeASCII)-1])
	}
	return 0
}

// GetCurrentCharIndex is final in Java.
//
// Java: CodePrinter.MappedCodePrinter.getCurrentCharIndex()
func (m *CodePrinter_MappedCodePrinterBase) GetCurrentCharIndex() int {
	return m.lineLength
}

// GetCurrentLineIndex is final in Java.
//
// Java: CodePrinter.MappedCodePrinter.getCurrentLineIndex()
func (m *CodePrinter_MappedCodePrinterBase) GetCurrentLineIndex() int {
	return m.lineIndex
}

// CodePrinter_MappedCodePrinter_ComputeLineLengths calculates length of each
// line in compiled code.
//
// Java: CodePrinter.MappedCodePrinter.computeLineLengths(String)
func CodePrinter_MappedCodePrinter_ComputeLineLengths(code string) []int {
	// Go-only: language-adapter — the lengths are in UTF-16 units, as
	// Java's indexOf positions are; '\n' is ASCII, so lines can be cut on
	// its bytes.
	var builder []int
	lineStartPos := 0
	lineEndPos := strings.IndexByte(code, '\n')
	for lineEndPos > -1 {
		builder = append(builder, jutil.StringLength(code[lineStartPos:lineEndPos]))
		// Next line starts where current line ends + 1 to skip "\n"
		// character.
		lineStartPos = lineEndPos + 1
		lineEndPos = strings.IndexByte(code[lineStartPos:], '\n')
		if lineEndPos > -1 {
			lineEndPos += lineStartPos
		}
	}
	return builder
}

func codePrinter_MappedCodePrinter_adjustEndPositionVal(
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

// CodePrinter_MappedCodePrinter_AdjustEndPosition adjusts end position of a
// mapping. End position points to a column *after* the last character that
// is covered by a mapping. And if it's end of the line there are 2
// possibilities: either point to the non-existent character after the last
// char on a line or point to the first character on the next line. In some
// cases we end up with 2 mappings which should have the same end position,
// but they use different styles as described above it leads to invalid
// source maps.
//
// This method adjusts all such end positions, so if it points to the
// non-existing character at the end of line - it is changed to point to the
// first character on the next line.
//
// lineLengths is the list of all line lengths in compiled code; endPosition
// is the end position of a mapping.
//
// Java: CodePrinter.MappedCodePrinter.adjustEndPosition(List<Integer>, FilePosition)
func CodePrinter_MappedCodePrinter_AdjustEndPosition(
	lineLengths []int, endPosition *sourcemap.FilePosition) *sourcemap.FilePosition {
	line := endPosition.GetLine()
	// if position points to non-existing line, return it unmodified
	if line >= len(lineLengths) {
		return endPosition
	}

	jutil.CheckStateMsg2(
		endPosition.GetColumn() <= lineLengths[line],
		"End position %s points to a column larger than line length %s",
		endPosition,
		lineLengths[line])

	// if end position points to the column just after the last character on
	// the line - change it to point the first character on the next line
	if endPosition.GetColumn() == lineLengths[line] {
		return sourcemap.NewFilePosition(line+1, 0)
	}
	return endPosition
}

// CodePrinter_PrettyCodePrinter_INDENT is PrettyCodePrinter.INDENT.
const CodePrinter_PrettyCodePrinter_INDENT = "  "

// CodePrinter_PrettyCodePrinter is the pretty-printing code consumer.
type CodePrinter_PrettyCodePrinter struct {
	CodePrinter_MappedCodePrinterBase

	indent int

	// Go-only: language-adapter — getNumberFromSource takes the length of,
	// and a substring of, the whole source code for every number; a Go
	// string has neither in O(1) by UTF-16 index, so the UTF-16 view of the
	// last source code is kept here (keyed by the immutable string itself).
	numberSourceCode  string
	numberSourceASCII bool
	numberSourceUnits []uint16
}

var codePrinter_prettyPool = sync.Pool{
	New: func() any {
		return &CodePrinter_PrettyCodePrinter{}
	},
}

func initCodePrinter_PrettyCodePrinter(
	p *CodePrinter_PrettyCodePrinter,
	lineLengthThreshold int,
	createSourceMap bool,
	sourceMapDetailLevel SourceMap_DetailLevel,
	licenseTracker CodePrinter_LicenseTracker) {
	p.InitCodePrinter_MappedCodePrinterBase(p, lineLengthThreshold, createSourceMap, sourceMapDetailLevel, licenseTracker)
	p.indent = 0
	p.numberSourceCode = ""
	p.numberSourceASCII = false
	p.numberSourceUnits = nil
}

// NewCodePrinter_PrettyCodePrinter: lineLengthThreshold is the length of a
// line after which we force a newline when possible; createSourceMap is
// whether to generate source map data; sourceMapDetailLevel is a filter to
// control which nodes get mapped into the source map; licenseTracker is a
// license tracking implementation to manage license text emit. The
// CodePrinter will never emit license information directly.
//
// Java: CodePrinter.PrettyCodePrinter(int, boolean, SourceMap.DetailLevel, LicenseTracker)
func NewCodePrinter_PrettyCodePrinter(
	lineLengthThreshold int,
	createSourceMap bool,
	sourceMapDetailLevel SourceMap_DetailLevel,
	licenseTracker CodePrinter_LicenseTracker) *CodePrinter_PrettyCodePrinter {
	p := &CodePrinter_PrettyCodePrinter{}
	initCodePrinter_PrettyCodePrinter(p, lineLengthThreshold, createSourceMap, sourceMapDetailLevel, licenseTracker)
	return p
}

// Append appends an appropriately indented string to the code, keeping
// track of the current line count and line length.
//
// Java: CodePrinter.PrettyCodePrinter.append(String)
func (p *CodePrinter_PrettyCodePrinter) Append(str string) {
	// For pretty printing: indent at the beginning of the line, except
	// template literal lines.
	if p.lineLength == 0 && !p.IsInTemplateLiteral() {
		for i := 0; i < p.indent; i++ {
			p.codeAppendASCII(CodePrinter_PrettyCodePrinter_INDENT)
			p.lineLength += len(CodePrinter_PrettyCodePrinter_INDENT)
		}
	}

	p.CodePrinter_MappedCodePrinterBase.Append(str)
}

// AddNumber attempts to read the number format out of the original source
// location, falling back to the default behavior if we cannot locate it.
//
// Java: CodePrinter.PrettyCodePrinter.addNumber(double, Node)
func (p *CodePrinter_PrettyCodePrinter) AddNumber(x float64, n *rhino.Node) {
	if !base.JSCompDoubles_IsPositive(x) {
		// checkState(isPositive(x), x)
		panic(jutil.NewIllegalStateExceptionMsg(jutil.DoubleToString(x)))
	}

	numberFromSource := p.GetNumberFromSource(n)
	if numberFromSource == nil {
		p.CodePrinter_MappedCodePrinterBase.AddNumber(x, n)
		return
	}

	// The string we extract from the source code is not always a number.
	// Conservatively, we only use it if we can verify that it is as a
	// number with the right value. This excludes some valid constants (hex,
	// etc.) for simplicity.
	d, ok := codePrinter_parseDouble(*numberFromSource)
	if !ok {
		// catch (NumberFormatException e)
		p.CodePrinter_MappedCodePrinterBase.AddNumber(x, n)
		return
	}

	if x != d {
		p.CodePrinter_MappedCodePrinterBase.AddNumber(x, n)
		return
	}

	p.self.AddConstant(*numberFromSource)
}

// codePrinter_parseDouble is `Double.parseDouble(s)` inside
// `try { ... } catch (NumberFormatException e)`; ok is false when it
// threw.
//
// Go-only: language-adapter — a named function with defer/recover instead
// of a jutil.Try closure, as this runs for every number (PORTING.md §5).
func codePrinter_parseDouble(s string) (d float64, ok bool) {
	defer func() {
		if r := recover(); r != nil {
			if _, isNFE := r.(jutil.INumberFormatException); isNFE {
				d, ok = 0, false
				return
			}
			panic(r)
		}
	}()
	return jutil.DoubleParseDouble(s), true
}

// StartNewLine adds a newline to the code, resetting the line length and
// handling indenting for pretty printing.
//
// Java: CodePrinter.PrettyCodePrinter.startNewLine()
func (p *CodePrinter_PrettyCodePrinter) StartNewLine() {
	if p.lineLength <= 0 && !p.IsInTemplateLiteral() {
		return
	}

	p.codeAppendByte('\n')
	p.lineIndex++
	p.lineLength = 0
}

// Java: CodePrinter.PrettyCodePrinter.maybeLineBreak()
func (p *CodePrinter_PrettyCodePrinter) MaybeLineBreak() {
	p.self.MaybeCutLine()
}

// MaybeCutLine may start a new line if the current line is longer than the
// line length threshold.
//
// Java: CodePrinter.PrettyCodePrinter.maybeCutLine()
func (p *CodePrinter_PrettyCodePrinter) MaybeCutLine() {
	if p.lineLength > p.lineLengthThreshold {
		p.self.StartNewLine()
	}
}

// Java: CodePrinter.PrettyCodePrinter.endLine()
func (p *CodePrinter_PrettyCodePrinter) EndLine() {
	p.self.StartNewLine()
}

// Java: CodePrinter.PrettyCodePrinter.appendBlockStart()
func (p *CodePrinter_PrettyCodePrinter) AppendBlockStart() {
	p.self.MaybeInsertSpace()
	p.self.Add("{")
	p.indent++
}

// Java: CodePrinter.PrettyCodePrinter.appendBlockEnd()
func (p *CodePrinter_PrettyCodePrinter) AppendBlockEnd() {
	p.self.MaybeEndStatement()
	p.self.EndLine()
	p.indent--
	p.self.Add("}")
}

// Java: CodePrinter.PrettyCodePrinter.listSeparator()
func (p *CodePrinter_PrettyCodePrinter) ListSeparator() {
	p.self.Add(", ")
	p.self.MaybeLineBreak()
}

// Java: CodePrinter.PrettyCodePrinter.optionalListSeparator()
func (p *CodePrinter_PrettyCodePrinter) OptionalListSeparator() {
	p.self.Add(",")
	p.self.MaybeLineBreak()
}

// Java: CodePrinter.PrettyCodePrinter.endFunction(boolean)
func (p *CodePrinter_PrettyCodePrinter) EndFunction(statementContext bool) {
	p.CodePrinter_MappedCodePrinterBase.EndFunction(statementContext)
	if statementContext {
		p.self.StartNewLine()
	}
}

// Java: CodePrinter.PrettyCodePrinter.beginCaseBody()
func (p *CodePrinter_PrettyCodePrinter) BeginCaseBody() {
	p.CodePrinter_MappedCodePrinterBase.BeginCaseBody()
	p.indent++
	p.self.EndLine()
}

// Java: CodePrinter.PrettyCodePrinter.endCaseBody()
func (p *CodePrinter_PrettyCodePrinter) EndCaseBody() {
	p.CodePrinter_MappedCodePrinterBase.EndCaseBody()
	p.indent--
}

// Java: CodePrinter.PrettyCodePrinter.appendOp(String, boolean)
func (p *CodePrinter_PrettyCodePrinter) AppendOp(op string, binOp bool) {
	if p.self.GetLastChar() != ' ' && binOp && codeConsumer_firstChar(op) != ',' {
		p.self.Add(" ")
	}
	p.self.Add(op)
	if binOp {
		p.self.Add(" ")
	}
}

// ShouldPreserveExtras: if the body of a for loop or the then clause of an
// if statement has a single statement, should it be wrapped in a block?
// And similar.
//
// Java: CodePrinter.PrettyCodePrinter.shouldPreserveExtras(Node)
func (p *CodePrinter_PrettyCodePrinter) ShouldPreserveExtras(n *rhino.Node) bool {
	// When pretty-printing, always place the statement in its own block so
	// it is printed on a separate line. This allows breakpoints to be
	// placed on the statement.
	//
	// The only exception is an added block around an else-if statement.
	// Example code:
	//   if (0) {
	//     0;
	//   } else if (1) {
	//     1;
	//   }
	// Resulting node tree:
	//   IF
	//     NUMBER
	//     BLOCK
	//       EXPR_RESULT
	//         NUMBER
	//     BLOCK [added_block: 1] <-- we don't want to print this block
	//       IF
	//         NUMBER
	//         BLOCK
	//           EXPR_RESULT
	//             NUMBER

	if !n.IsBlock() || !n.IsAddedBlock() || !n.HasParent() {
		return true
	}

	parent := n.GetParent()
	isElse := parent.IsIf() && parent.HasXChildren(3) && n == parent.GetLastChild()
	onlyChildIsIf := n.HasOneChild() && n.GetFirstChild().IsIf()
	if isElse && onlyChildIsIf {
		return false
	}

	return true
}

// Java: CodePrinter.PrettyCodePrinter.maybeInsertSpace()
func (p *CodePrinter_PrettyCodePrinter) MaybeInsertSpace() {
	if p.self.GetLastChar() != ' ' && p.self.GetLastChar() != '\n' {
		p.self.Add(" ")
	}
}

// CodePrinter_PrettyCodePrinter_GetTryForCatch returns the TRY node for the
// specified CATCH node.
//
// Java: CodePrinter.PrettyCodePrinter.getTryForCatch(Node)
func CodePrinter_PrettyCodePrinter_GetTryForCatch(n *rhino.Node) *rhino.Node {
	return n.GetGrandparent()
}

// BreakAfterBlockFor returns whether the a line break should be added after
// the specified BLOCK.
//
// Java: CodePrinter.PrettyCodePrinter.breakAfterBlockFor(Node, boolean)
func (p *CodePrinter_PrettyCodePrinter) BreakAfterBlockFor(n *rhino.Node, isStatementContext bool) bool {
	jutil.CheckStateMsg1(n.IsBlock(), "%s", n)
	parent := n.GetParent()
	type_ := parent.GetToken()
	switch type_ {
	case rhino.Token_DO:
		return false // Don't break before 'while' in DO-WHILE statements.
	case rhino.Token_FUNCTION:
		return false // FUNCTIONs are handled separately, don't break here.
	case rhino.Token_TRY:
		return n != parent.GetFirstChild() // Don't break before catch
	case rhino.Token_CATCH:
		return !NodeUtil_HasFinally(CodePrinter_PrettyCodePrinter_GetTryForCatch(parent)) // Don't break before finally
	case rhino.Token_IF:
		return n == parent.GetLastChild() // Don't break before else
	default:
		return true
	}
}

// Java: CodePrinter.PrettyCodePrinter.endStatement(boolean, boolean)
func (p *CodePrinter_PrettyCodePrinter) EndStatementWithSemiColon(needsSemicolon bool, hasTrailingCommentOnSameLine bool) {
	p.self.Add(";")
	if !hasTrailingCommentOnSameLine {
		p.self.EndLine()
	}
	p.statementNeedsEnded = false
}

// Java: CodePrinter.PrettyCodePrinter.endFile()
func (p *CodePrinter_PrettyCodePrinter) EndFile() {
	p.self.MaybeEndStatement()
}

// GetNumberFromSource is static in Java.
//
// Go-only: language-adapter — a method, to reach the printer's UTF-16 view
// of the source code (see the numberSource fields).
//
// Java: CodePrinter.PrettyCodePrinter.getNumberFromSource(Node)
func (p *CodePrinter_PrettyCodePrinter) GetNumberFromSource(n *rhino.Node) *string {
	if !n.IsNumber() {
		return nil
	}

	staticSrc := NodeUtil_GetSourceFile(n)
	src, isSourceFile := staticSrc.(*SourceFile)
	if !isSourceFile || src == nil {
		return nil
	}

	if src.IsStubSourceFileForAlreadyProvidedInput() {
		// source file is a stub file, so we can not get number from source.
		return nil
	}
	srcCode, ok := codePrinter_getCode(src)
	if !ok {
		// catch (IOException e)
		return nil
	}

	offset, ok := codePrinter_getSourceOffset(n)
	if !ok {
		// catch (IllegalArgumentException e)
		return nil
	}
	endOffset := offset + n.GetLength()
	if offset < 0 || endOffset > p.numberSourceLength(srcCode) {
		return nil
	}

	s := p.numberSourceSubstring(offset, endOffset)
	return &s
}

// numberSourceLength is srcCode.length(), remembering the UTF-16 view of
// srcCode for numberSourceSubstring.
//
// Go-only: language-adapter — see the numberSource fields.
func (p *CodePrinter_PrettyCodePrinter) numberSourceLength(srcCode string) int {
	if srcCode != p.numberSourceCode {
		p.numberSourceCode = srcCode
		p.numberSourceASCII = true
		for i := 0; i < len(srcCode); i++ {
			if srcCode[i] >= 0x80 {
				p.numberSourceASCII = false
				break
			}
		}
		if p.numberSourceASCII {
			p.numberSourceUnits = nil
		} else {
			p.numberSourceUnits = jutil.StringToUTF16(srcCode)
		}
	}
	if p.numberSourceASCII {
		return len(p.numberSourceCode)
	}
	return len(p.numberSourceUnits)
}

// numberSourceSubstring is srcCode.substring(begin, end) on the code last
// passed to numberSourceLength, with 0 <= begin and end <= its length.
//
// Go-only: language-adapter — see the numberSource fields.
func (p *CodePrinter_PrettyCodePrinter) numberSourceSubstring(begin int, end int) string {
	if begin > end {
		panic(jutil.NewStringIndexOutOfBoundsExceptionForRange(begin, end, p.numberSourceLength(p.numberSourceCode)))
	}
	if p.numberSourceASCII {
		return p.numberSourceCode[begin:end]
	}
	return jutil.StringFromUTF16(p.numberSourceUnits[begin:end])
}

// codePrinter_getCode is `src.getCode()` inside
// `try { ... } catch (IOException e)`; ok is false when it threw.
//
// Go-only: language-adapter — a named function with defer/recover instead
// of a jutil.Try closure, as this runs for every number (PORTING.md §5).
func codePrinter_getCode(src *SourceFile) (code string, ok bool) {
	defer func() {
		if r := recover(); r != nil {
			if _, isIOE := r.(jutil.IIOException); isIOE {
				code, ok = "", false
				return
			}
			panic(r)
		}
	}()
	return src.GetCode(), true
}

// codePrinter_getSourceOffset is `n.getSourceOffset()` inside
// `try { ... } catch (IllegalArgumentException e)`; ok is false when it
// threw.
//
// Go-only: language-adapter — as codePrinter_getCode.
func codePrinter_getSourceOffset(n *rhino.Node) (offset int, ok bool) {
	defer func() {
		if r := recover(); r != nil {
			if _, isIAE := r.(jutil.IIllegalArgumentException); isIAE {
				offset, ok = 0, false
				return
			}
			panic(r)
		}
	}()
	return n.GetSourceOffset(), true
}

// CodePrinter_CompactCodePrinter is the compact code consumer.
//
// The CompactCodePrinter tries to emit just enough newlines to stop there
// being lines longer than the threshold. Since the output is going to be
// gzipped, it makes sense to try to make the newlines appear in similar
// contexts so that gzip can encode them for 'free'.
//
// This version tries to break the lines at 'preferred' places, which are
// between the top-level forms. This works because top-level forms tend to
// be more uniform than arbitrary legal contexts. Better compression would
// probably require explicit modeling of the gzip algorithm.
type CodePrinter_CompactCodePrinter struct {
	CodePrinter_MappedCodePrinterBase

	lineBreak bool

	lineStartPosition      int
	preferredBreakPosition int
}

var codePrinter_compactPool = sync.Pool{
	New: func() any {
		return &CodePrinter_CompactCodePrinter{}
	},
}

func initCodePrinter_CompactCodePrinter(
	c *CodePrinter_CompactCodePrinter,
	lineBreak bool,
	lineLengthThreshold int,
	createSrcMap bool,
	sourceMapDetailLevel SourceMap_DetailLevel,
	licenseTracker CodePrinter_LicenseTracker) {
	c.InitCodePrinter_MappedCodePrinterBase(c, lineLengthThreshold, createSrcMap, sourceMapDetailLevel, licenseTracker)
	c.lineStartPosition = 0
	c.preferredBreakPosition = 0
	c.lineBreak = lineBreak
}

// NewCodePrinter_CompactCodePrinter: lineBreak is to break the lines a bit
// more aggressively; lineLengthThreshold is the length of a line after
// which we force a newline when possible; createSrcMap is whether to gather
// source position mapping information when printing; sourceMapDetailLevel
// is a filter to control which nodes get mapped into the source map;
// licenseTracker is a license tracking implementation to manage license
// text emit. The CodePrinter will never emit license information directly -
// it only ever passes nodes to the license tracker to request tracking.
//
// Java: CodePrinter.CompactCodePrinter(boolean, int, boolean, SourceMap.DetailLevel, LicenseTracker)
func NewCodePrinter_CompactCodePrinter(
	lineBreak bool,
	lineLengthThreshold int,
	createSrcMap bool,
	sourceMapDetailLevel SourceMap_DetailLevel,
	licenseTracker CodePrinter_LicenseTracker) *CodePrinter_CompactCodePrinter {
	c := &CodePrinter_CompactCodePrinter{}
	initCodePrinter_CompactCodePrinter(c, lineBreak, lineLengthThreshold, createSrcMap, sourceMapDetailLevel, licenseTracker)
	return c
}

// StartNewLine adds a newline to the code, resetting the line length.
//
// Java: CodePrinter.CompactCodePrinter.startNewLine()
func (c *CodePrinter_CompactCodePrinter) StartNewLine() {
	if c.lineLength <= 0 && !c.IsInTemplateLiteral() {
		return
	}

	c.codeAppendByte('\n')
	c.lineLength = 0
	c.lineIndex++
	c.lineStartPosition = c.codeLength()
}

// Java: CodePrinter.CompactCodePrinter.maybeLineBreak()
func (c *CodePrinter_CompactCodePrinter) MaybeLineBreak() {
	if c.lineBreak {
		if c.sawFunction {
			c.self.StartNewLine()
			c.sawFunction = false
		}
	}

	// Since we are at a legal line break, can we upgrade the
	// preferred break position? We prefer to break after a
	// semicolon rather than before it.
	len_ := c.codeLength()
	if c.preferredBreakPosition == len_-1 {
		ch := c.codeCharAt(len_ - 1)
		if ch == ';' {
			c.preferredBreakPosition = len_
		}
	}
	c.self.MaybeCutLine()
}

// MaybeCutLine may start a new line if the current line is longer than the
// line length threshold.
//
// Java: CodePrinter.CompactCodePrinter.maybeCutLine()
func (c *CodePrinter_CompactCodePrinter) MaybeCutLine() {
	if c.lineLength <= c.lineLengthThreshold {
		return
	}

	// Use the preferred position provided it will break the line.
	if c.preferredBreakPosition > c.lineStartPosition &&
		c.preferredBreakPosition < c.lineStartPosition+c.lineLength {
		// If the preferred break position is on the current line.
		c.codeInsertByte(c.preferredBreakPosition, '\n')
		c.ReportLineCut(c.lineIndex, c.preferredBreakPosition-c.lineStartPosition)
		c.lineIndex++
		c.lineLength -= (c.preferredBreakPosition - c.lineStartPosition)
		c.lineStartPosition = c.preferredBreakPosition + 1 // Jump over the inserted newline.
	} else {
		c.self.StartNewLine()
	}
}

// Java: CodePrinter.CompactCodePrinter.notePreferredLineBreak()
func (c *CodePrinter_CompactCodePrinter) NotePreferredLineBreak() {
	c.preferredBreakPosition = c.codeLength()
}

// CodePrinter_LicenseTracker: License Trackers are responsible for ensuring
// that any licensing information attached to nodes is retained in the final
// output of JSCompiler.
//
// The Code Printer will visit every node being printed and call
// trackLicensesForNode for that node. Later, some other code can call
// emitLicenses to retrieve all relevant licenses. Deciding what is relevant
// varies between implementations of License Trackers - read their javadoc
// to determine how they are intended to be used.
type CodePrinter_LicenseTracker interface {
	// Java: CodePrinter.LicenseTracker.trackLicensesForNode(Node)
	TrackLicensesForNode(node *rhino.Node)

	// Java: CodePrinter.LicenseTracker.emitLicenses()
	EmitLicenses() *jutil.LinkedHashSet[string]
}

// CodePrinter_Builder builds the source code of an AST.
type CodePrinter_Builder struct {
	root                 *rhino.Node
	options              *CompilerOptions
	lineBreak            bool
	prettyPrint          bool
	outputTypes          bool
	tagAsTypeSummary     bool
	tagAsStrict          bool
	licenseTracker       CodePrinter_LicenseTracker // nil is null
	registry             *rhino.JSTypeRegistry      // may be null unless using Format.TYPED
	codeGeneratorFactory CodePrinter_Builder_CodeGeneratorFactory
}

// codePrinter_Builder_defaultCodeGeneratorFactory is the Builder's
// anonymous default CodeGeneratorFactory; it reads the builder's options
// and registry when it is called.
type codePrinter_Builder_defaultCodeGeneratorFactory struct {
	builder *CodePrinter_Builder
}

// Java: CodePrinter.Builder.<anonymous CodeGeneratorFactory>.getCodeGenerator(Format, CodeConsumer)
func (f *codePrinter_Builder_defaultCodeGeneratorFactory) GetCodeGenerator(outputFormat CodePrinter_Format, cc CodeConsumer) CodeGenerator {
	if outputFormat == CodePrinter_Format_TYPED {
		return NewTypedCodeGenerator(cc, f.builder.getOptions(), f.builder.registry)
	}
	return NewCodeGeneratorWithOptions(cc, f.builder.getOptions())
}

// NewCodePrinter_Builder sets the root node from which to generate the
// source code.
//
// Java: CodePrinter.Builder(Node)
func NewCodePrinter_Builder(node *rhino.Node) *CodePrinter_Builder {
	b := &CodePrinter_Builder{}
	// Go-only: performance — Java initializes options to a new
	// CompilerOptions(), which almost every caller (Compiler.toSource, once
	// per printed node in checks such as CheckDuplicateCase) replaces at
	// once with setCompilerOptions. Building one costs ~100 allocations
	// (its DefaultNameGenerator alone builds a 64-entry map), so it is
	// created by getOptions only if it is used before being replaced.
	b.options = nil
	b.outputTypes = false
	b.codeGeneratorFactory = nil
	b.root = node
	return b
}

// getCodeGeneratorFactory returns codeGeneratorFactory, creating the default
// factory only when passing to external callers of CodePrinter_ToSource.
func (b *CodePrinter_Builder) getCodeGeneratorFactory() CodePrinter_Builder_CodeGeneratorFactory {
	if b.codeGeneratorFactory == nil {
		b.codeGeneratorFactory = &codePrinter_Builder_defaultCodeGeneratorFactory{builder: b}
	}
	return b.codeGeneratorFactory
}

// getOptions returns options, made on first use if setCompilerOptions has
// not set them (see NewCodePrinter_Builder).
//
// Go-only: performance — Java reads the field.
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

// SetLicenseTracker sets the license tracker to use when printing. It can
// be null to disable license tracking.
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

// SetCodeGeneratorFactory sets a custom code generator factory to enable
// custom code generation.
//
// Java: CodePrinter.Builder.setCodeGeneratorFactory(CodeGeneratorFactory)
func (b *CodePrinter_Builder) SetCodeGeneratorFactory(factory CodePrinter_Builder_CodeGeneratorFactory) *CodePrinter_Builder {
	b.codeGeneratorFactory = factory
	return b
}

// CodePrinter_Builder_CodeGeneratorFactory creates the code generator of a
// build.
type CodePrinter_Builder_CodeGeneratorFactory interface {
	// Java: CodePrinter.Builder.CodeGeneratorFactory.getCodeGenerator(Format, CodeConsumer)
	GetCodeGenerator(outputFormat CodePrinter_Format, cc CodeConsumer) CodeGenerator
}

// Build generates the source code and returns it.
//
// Java: CodePrinter.Builder.build()
func (b *CodePrinter_Builder) Build() string {
	return b.buildSourceAndMappingsInto(nil, false)
}

// Java: CodePrinter.Builder.buildWithSourceMappings()
func (b *CodePrinter_Builder) BuildWithSourceMappings() *CodePrinter_SourceAndMappings {
	sm := &CodePrinter_SourceAndMappings{}
	b.buildSourceAndMappingsInto(sm, true)
	sm.mappings = sm.iterMappings
	return sm
}

func codePrinter_emptyMappingsSeq(func(int, *SourceMap_Mapping) bool) {}

func (b *CodePrinter_Builder) buildSourceAndMappingsInto(sm *CodePrinter_SourceAndMappings, wantMappings bool) string {
	if b.root == nil {
		panic(jutil.NewIllegalStateExceptionMsg("Cannot build without root node being specified"))
	}
	options := b.getOptions()
	outputFormat := CodePrinter_Format_FromOptions(options, b.outputTypes, b.prettyPrint)
	gatherSourceMap := wantMappings && options.ShouldGatherSourceMapInfo()

	var mcp CodePrinter_MappedCodePrinter
	var pooledCompact *CodePrinter_CompactCodePrinter
	var pooledPretty *CodePrinter_PrettyCodePrinter
	if outputFormat == CodePrinter_Format_COMPACT {
		pooledCompact = codePrinter_compactPool.Get().(*CodePrinter_CompactCodePrinter)
		initCodePrinter_CompactCodePrinter(
			pooledCompact,
			b.lineBreak,
			options.GetLineLengthThreshold(),
			gatherSourceMap,
			options.GetSourceMapDetailLevel(),
			b.licenseTracker)
		mcp = pooledCompact
	} else {
		pooledPretty = codePrinter_prettyPool.Get().(*CodePrinter_PrettyCodePrinter)
		initCodePrinter_PrettyCodePrinter(
			pooledPretty,
			options.GetLineLengthThreshold(),
			gatherSourceMap,
			options.GetSourceMapDetailLevel(),
			b.licenseTracker)
		mcp = pooledPretty
	}

	var cg CodeGenerator
	if b.codeGeneratorFactory != nil {
		cg = b.codeGeneratorFactory.GetCodeGenerator(outputFormat, mcp)
	} else if outputFormat == CodePrinter_Format_TYPED {
		cg = NewTypedCodeGenerator(mcp, options, b.registry)
	} else {
		baseCG := &mcp.GetCodePrinter_MappedCodePrinterBase().cg
		baseCG.InitCodeGeneratorBaseWithOptions(baseCG, mcp, options)
		cg = baseCG
	}

	if b.tagAsTypeSummary {
		cg.TagAsTypeSummary()
	}
	if b.tagAsStrict {
		cg.TagAsStrict()
	}

	cg.AddNode(b.root)
	mcp.EndFile()

	source := mcp.GetCode()
	if sm != nil {
		sm.source = source
		if gatherSourceMap {
			base := mcp.GetCodePrinter_MappedCodePrinterBase()
			base.adjustSourceMappings(source)
			sm.takeMappingsFrom(base)
		} else {
			sm.mappingChunks = nil
			sm.mappingSlab = nil
		}
	}
	returnPooledPrinter(pooledCompact, pooledPretty)
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

// CodePrinter_Format_values is Format.values().
func CodePrinter_Format_values() []CodePrinter_Format {
	return []CodePrinter_Format{CodePrinter_Format_COMPACT, CodePrinter_Format_PRETTY, CodePrinter_Format_TYPED}
}

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
	source           string
	mappings         iter.Seq2[int, *SourceMap_Mapping] // nil is null
	mappingChunks    [][]SourceMap_Mapping
	mappingSlab      []SourceMap_Mapping
	mappingChunksBuf [16][]SourceMap_Mapping
	mappingBuf       [32]SourceMap_Mapping
}

// Java: CodePrinter.SourceAndMappings()
func NewCodePrinter_SourceAndMappings() *CodePrinter_SourceAndMappings {
	return &CodePrinter_SourceAndMappings{}
}

func (s *CodePrinter_SourceAndMappings) takeMappingsFrom(m *CodePrinter_MappedCodePrinterBase) {
	if len(m.mappingChunks) == 0 {
		s.mappingChunks = nil
		n := len(m.mappingSlab)
		if n == 0 {
			s.mappingSlab = nil
		} else {
			copy(s.mappingBuf[:n], m.mappingSlab)
			clear(m.mappingSlab)
			s.mappingSlab = s.mappingBuf[:n]
		}
		m.mappingSlab = nil
		return
	}
	n0 := len(m.mappingChunks[0])
	copy(s.mappingBuf[:n0], m.mappingChunks[0])
	clear(m.mappingChunks[0])
	if len(m.mappingChunks) <= len(s.mappingChunksBuf) {
		nChunks := len(m.mappingChunks)
		copy(s.mappingChunksBuf[:nChunks], m.mappingChunks)
		clear(m.mappingChunks)
		s.mappingChunks = s.mappingChunksBuf[:nChunks]
	} else {
		clear(m.mappingChunksBuf[:])
		s.mappingChunks = m.mappingChunks
	}
	s.mappingChunks[0] = s.mappingBuf[:n0]
	s.mappingSlab = m.mappingSlab
	m.mappingChunks = nil
	m.mappingSlab = nil
}

func (s *CodePrinter_SourceAndMappings) iterMappings(yield func(int, *SourceMap_Mapping) bool) {
	idx := 0
	for _, chunk := range s.mappingChunks {
		for i := range chunk {
			if !yield(idx, &chunk[i]) {
				return
			}
			idx++
		}
	}
	for i := range s.mappingSlab {
		if !yield(idx, &s.mappingSlab[i]) {
			return
		}
		idx++
	}
	s.Release()
}

// Release returns any pooled mapping slabs owned by s to the pool.
func (s *CodePrinter_SourceAndMappings) Release() {
	if s.mappingChunks == nil && s.mappingSlab == nil {
		return
	}
	for i, chunk := range s.mappingChunks {
		if cap(chunk) == codePrinter_mappingSlabSize {
			clear(chunk[:cap(chunk)])
			codePrinter_mappingSlabPool.Put((*[codePrinter_mappingSlabSize]SourceMap_Mapping)(chunk[:cap(chunk)]))
		} else {
			clear(chunk)
		}
		s.mappingChunks[i] = nil
	}
	s.mappingChunks = nil
	if cap(s.mappingSlab) == codePrinter_mappingSlabSize {
		clear(s.mappingSlab[:cap(s.mappingSlab)])
		codePrinter_mappingSlabPool.Put((*[codePrinter_mappingSlabSize]SourceMap_Mapping)(s.mappingSlab[:cap(s.mappingSlab)]))
	} else {
		clear(s.mappingSlab)
	}
	s.mappingSlab = nil
}

// CodePrinter_ToSource converts a tree to JS code.
//
// Java: CodePrinter.toSource(Node, Format, CompilerOptions, LicenseTracker, boolean, boolean, boolean, CodeGeneratorFactory)
func CodePrinter_ToSource(
	root *rhino.Node,
	outputFormat CodePrinter_Format,
	options *CompilerOptions,
	licenseTracker CodePrinter_LicenseTracker,
	tagAsTypeSummary bool,
	tagAsStrict bool,
	lineBreak bool,
	codeGeneratorFactory CodePrinter_Builder_CodeGeneratorFactory) *CodePrinter_SourceAndMappings {
	// checkState(options.getSourceMapDetailLevel() != null)
	//
	// Go-only: language-adapter — an enum value cannot be null in Go, so the
	// check always holds.

	gatherSourceMap := options.ShouldGatherSourceMapInfo()
	var mcp CodePrinter_MappedCodePrinter
	var pooledCompact *CodePrinter_CompactCodePrinter
	var pooledPretty *CodePrinter_PrettyCodePrinter
	if outputFormat == CodePrinter_Format_COMPACT {
		pooledCompact = codePrinter_compactPool.Get().(*CodePrinter_CompactCodePrinter)
		initCodePrinter_CompactCodePrinter(
			pooledCompact,
			lineBreak,
			options.GetLineLengthThreshold(),
			gatherSourceMap,
			options.GetSourceMapDetailLevel(),
			licenseTracker)
		mcp = pooledCompact
	} else {
		pooledPretty = codePrinter_prettyPool.Get().(*CodePrinter_PrettyCodePrinter)
		initCodePrinter_PrettyCodePrinter(
			pooledPretty,
			options.GetLineLengthThreshold(),
			gatherSourceMap,
			options.GetSourceMapDetailLevel(),
			licenseTracker)
		mcp = pooledPretty
	}
	cg := codeGeneratorFactory.GetCodeGenerator(outputFormat, mcp)

	if tagAsTypeSummary {
		cg.TagAsTypeSummary()
	}
	if tagAsStrict {
		cg.TagAsStrict()
	}

	cg.AddNode(root)
	mcp.EndFile()

	result := NewCodePrinter_SourceAndMappings()
	result.source = mcp.GetCode()

	if gatherSourceMap {
		base := mcp.GetCodePrinter_MappedCodePrinterBase()
		base.adjustSourceMappings(result.source)
		result.takeMappingsFrom(base)
		result.mappings = result.iterMappings
	} else {
		result.mappings = codePrinter_emptyMappingsSeq
	}
	returnPooledPrinter(pooledCompact, pooledPretty)

	return result
}

// Java: CodePrinter()
func NewCodePrinter() *CodePrinter {
	return &CodePrinter{}
}
