// Port of com/google/debugging/sourcemap/SourceMapGeneratorV3.java @ closure-compiler v20261021

/*
 * Copyright 2011 The Closure Compiler Authors.
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

package sourcemap

import (
	"encoding/binary"
	"io"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// sourceMapGeneratorV3_mappingStore holds a sequence of mappings, each one
// encoded as the zigzag varint differences of its fields from those of the
// previous mapping (about 10 bytes instead of 32).
//
// Go-only: performance.
type sourceMapGeneratorV3_mappingStore struct {
	blocks [][]byte // full blocks, then cur
	cur    []byte
	n      int
	last   SourceMapGeneratorV3_Mapping // the last mapping added
}

const (
	sourceMapGeneratorV3_mappingBlockSize  = 64 << 10
	sourceMapGeneratorV3_maxEncodedMapping = 8 * binary.MaxVarintLen64
)

func (s *sourceMapGeneratorV3_mappingStore) reset() {
	*s = sourceMapGeneratorV3_mappingStore{}
}

func (s *sourceMapGeneratorV3_mappingStore) add(m *SourceMapGeneratorV3_Mapping) {
	if cap(s.cur)-len(s.cur) < sourceMapGeneratorV3_maxEncodedMapping {
		if s.cur != nil {
			s.blocks = append(s.blocks, s.cur)
		}
		s.cur = make([]byte, 0, sourceMapGeneratorV3_mappingBlockSize)
	}
	p := &s.last
	b := s.cur
	b = binary.AppendVarint(b, int64(m.startLine)-int64(p.startLine))
	b = binary.AppendVarint(b, int64(m.startCol)-int64(p.startCol))
	b = binary.AppendVarint(b, int64(m.endLine)-int64(m.startLine))
	b = binary.AppendVarint(b, int64(m.endCol)-int64(m.startCol))
	b = binary.AppendVarint(b, int64(m.sourceFileIdx)-int64(p.sourceFileIdx))
	b = binary.AppendVarint(b, int64(m.originalLine)-int64(p.originalLine))
	b = binary.AppendVarint(b, int64(m.originalCol)-int64(p.originalCol))
	b = binary.AppendVarint(b, int64(m.originalNameIdx)-int64(p.originalNameIdx))
	s.cur = b
	s.last = *m
	s.n++
}

// sourceMapGeneratorV3_mappingReader decodes a mappingStore in order.
type sourceMapGeneratorV3_mappingReader struct {
	store *sourceMapGeneratorV3_mappingStore
	block int
	buf   []byte
	prev  SourceMapGeneratorV3_Mapping
}

func (r *sourceMapGeneratorV3_mappingReader) varint() int32 {
	v, n := binary.Varint(r.buf)
	r.buf = r.buf[n:]
	return int32(v)
}

// next decodes the next mapping into m, and reports whether there was one.
func (r *sourceMapGeneratorV3_mappingReader) next(m *SourceMapGeneratorV3_Mapping) bool {
	for len(r.buf) == 0 {
		switch {
		case r.block < len(r.store.blocks):
			r.buf = r.store.blocks[r.block]
		case r.block == len(r.store.blocks):
			r.buf = r.store.cur
		default:
			return false
		}
		r.block++
	}
	p := &r.prev
	m.startLine = p.startLine + r.varint()
	m.startCol = p.startCol + r.varint()
	m.endLine = m.startLine + r.varint()
	m.endCol = m.startCol + r.varint()
	m.sourceFileIdx = p.sourceFileIdx + r.varint()
	m.originalLine = p.originalLine + r.varint()
	m.originalCol = p.originalCol + r.varint()
	m.originalNameIdx = p.originalNameIdx + r.varint()
	*p = *m
	return true
}

// SourceMapGeneratorV3_Mapping is a mapping from a given position in an
// input source file to a given position in the generated code.
type SourceMapGeneratorV3_Mapping struct {
	// Index into SourceMapGeneratorV3.addedSymbolNames, or -1 if none.
	originalNameIdx int32

	// Index into SourceMapGeneratorV3.addedSourceFiles.
	sourceFileIdx int32

	// The position of the code in the input source file.
	originalLine int32
	originalCol  int32

	// The starting position of the code in the generated source file which
	// this mapping represents. Indexed by 0.
	startLine int32
	startCol  int32

	// The ending position of the code in the generated source file which
	// this mapping represents. Indexed by 0.
	endLine int32
	endCol  int32
}

// NewSourceMapGeneratorV3_Mapping constructs an unused mapping.
//
// Java: SourceMapGeneratorV3.Mapping()
func NewSourceMapGeneratorV3_Mapping() *SourceMapGeneratorV3_Mapping {
	return &SourceMapGeneratorV3_Mapping{originalNameIdx: -1}
}

// SourceMapGeneratorV3_UsedMappingCheck marks any visited mapping as
// "used".
type SourceMapGeneratorV3_UsedMappingCheck struct {
	maxLine int
}

// Java: SourceMapGeneratorV3.UsedMappingCheck.visit(Mapping, int, int, int, int)
func (c *SourceMapGeneratorV3_UsedMappingCheck) Visit(m *SourceMapGeneratorV3_Mapping, line, col, nextLine, nextCol int) {
	if m != nil {
		if endLine := int(m.endLine); endLine > c.maxLine {
			c.maxLine = endLine
		}
	}
}

// SourceMapGeneratorV3_MappingVisitor visits code segments.
type SourceMapGeneratorV3_MappingVisitor interface {
	// Visit is called with the mapping for the current code segment (nil if
	// the segment is unmapped), the starting line and column for this code
	// segment, and the ending line and column.
	//
	// Java: SourceMapGeneratorV3.MappingVisitor.visit(Mapping, int, int, int, int)
	Visit(m *SourceMapGeneratorV3_Mapping, line, col, endLine, endCol int)
}

// SourceMapGeneratorV3_MappingTraversal walks the mappings and visits each
// segment of the mappings, unmapped segments are visited with a null
// mapping, unused mapping are not visited.
type SourceMapGeneratorV3_MappingTraversal struct {
	// Go-only: language-adapter — the enclosing SourceMapGeneratorV3.
	outer *SourceMapGeneratorV3

	// The last line and column written
	line int
	col  int

	// Go-only: performance — the decoded mapping being visited, the one
	// being closed, and the stack of open ones (fields, so that taking
	// their addresses does not allocate per mapping).
	cur, closing SourceMapGeneratorV3_Mapping
	stack        []SourceMapGeneratorV3_Mapping
}

// Java: SourceMapGeneratorV3.MappingTraversal()
func NewSourceMapGeneratorV3_MappingTraversal(outer *SourceMapGeneratorV3) *SourceMapGeneratorV3_MappingTraversal {
	return &SourceMapGeneratorV3_MappingTraversal{outer: outer}
}

// Traverse appends the line mapping entries.
//
// Java: SourceMapGeneratorV3.MappingTraversal.traverse(MappingVisitor)
func (t *SourceMapGeneratorV3_MappingTraversal) Traverse(v SourceMapGeneratorV3_MappingVisitor) {
	// The mapping list is ordered as a pre-order traversal. The mapping
	// positions give us enough information to rebuild the stack and this
	// allows the building of the source map in O(n) time.
	t.stack = t.stack[:0]
	r := sourceMapGeneratorV3_mappingReader{store: &t.outer.mappings}
	for r.next(&t.cur) {
		m := &t.cur
		// Find the closest ancestor of the current mapping:
		// An overlapping mapping is an ancestor of the current mapping, any
		// non-overlapping mappings are siblings (or cousins) and must be
		// closed in the reverse order of when they encountered.
		for len(t.stack) > 0 && !t.IsOverlapped(&t.stack[len(t.stack)-1], m) {
			t.closing = t.stack[len(t.stack)-1]
			t.stack = t.stack[:len(t.stack)-1]
			t.MaybeVisit(v, &t.closing)
		}

		// Any gaps between the current line position and the start of the
		// current mapping belong to the parent.
		var parent *SourceMapGeneratorV3_Mapping
		if len(t.stack) > 0 {
			parent = &t.stack[len(t.stack)-1]
		}
		t.MaybeVisitParent(v, parent, m)

		t.stack = append(t.stack, *m)
	}

	// There are no more children to be had, simply close the remaining
	// mappings in the reverse order of when they encountered.
	for len(t.stack) > 0 {
		t.closing = t.stack[len(t.stack)-1]
		t.stack = t.stack[:len(t.stack)-1]
		t.MaybeVisit(v, &t.closing)
	}
}

// GetAdjustedLine returns the line adjusted for the prefix position.
//
// Java: SourceMapGeneratorV3.MappingTraversal.getAdjustedLine(FilePosition)
func (t *SourceMapGeneratorV3_MappingTraversal) GetAdjustedLine(p FilePosition) int {
	return jutil.IntAdd(p.GetLine(), t.outer.prefixPosition.GetLine())
}

func (t *SourceMapGeneratorV3_MappingTraversal) getAdjustedLine(rawLine int32) int {
	return int(rawLine + int32(t.outer.prefixPosition.GetLine()))
}

// GetAdjustedCol returns the column adjusted for the prefix position.
//
// Java: SourceMapGeneratorV3.MappingTraversal.getAdjustedCol(FilePosition)
func (t *SourceMapGeneratorV3_MappingTraversal) GetAdjustedCol(p FilePosition) int {
	rawLine := p.GetLine()
	rawCol := p.GetColumn()
	// Only the first line needs the character position adjusted.
	if rawLine != 0 {
		return rawCol
	}
	return jutil.IntAdd(rawCol, t.outer.prefixPosition.GetColumn())
}

func (t *SourceMapGeneratorV3_MappingTraversal) getAdjustedCol(rawLine, rawCol int32) int {
	if rawLine != 0 {
		return int(rawCol)
	}
	return int(rawCol + int32(t.outer.prefixPosition.GetColumn()))
}

// IsOverlapped returns whether m1 ends before m2 starts.
//
// Java: SourceMapGeneratorV3.MappingTraversal.isOverlapped(Mapping, Mapping)
func (t *SourceMapGeneratorV3_MappingTraversal) IsOverlapped(m1, m2 *SourceMapGeneratorV3_Mapping) bool {
	// No need to use adjusted values here, relative positions are sufficient.
	l1 := m1.endLine
	l2 := m2.startLine
	c1 := m1.endCol
	c2 := m2.startCol

	return (l1 == l2 && c1 >= c2) || l1 > l2
}

// MaybeVisit writes any needed entries from the current position to the end
// of the provided mapping.
//
// Java: SourceMapGeneratorV3.MappingTraversal.maybeVisit(MappingVisitor, Mapping)
func (t *SourceMapGeneratorV3_MappingTraversal) MaybeVisit(v SourceMapGeneratorV3_MappingVisitor, m *SourceMapGeneratorV3_Mapping) {
	nextLine := t.getAdjustedLine(m.endLine)
	nextCol := t.getAdjustedCol(m.endLine, m.endCol)
	// If this anything remaining in this mapping beyond the current line and
	// column position, write it out now.
	if t.line < nextLine || (t.line == nextLine && t.col < nextCol) {
		t.Visit(v, m, nextLine, nextCol)
	}
}

// MaybeVisitParent writes any needed entries to complete the provided
// mapping.
//
// Java: SourceMapGeneratorV3.MappingTraversal.maybeVisitParent(MappingVisitor, Mapping, Mapping)
func (t *SourceMapGeneratorV3_MappingTraversal) MaybeVisitParent(v SourceMapGeneratorV3_MappingVisitor, parent, m *SourceMapGeneratorV3_Mapping) {
	nextLine := t.getAdjustedLine(m.startLine)
	nextCol := t.getAdjustedCol(m.startLine, m.startCol)
	// If the previous value is null, no mapping exists.
	jutil.CheckState(t.line < nextLine || t.col <= nextCol)
	if t.line < nextLine || (t.line == nextLine && t.col < nextCol) {
		t.Visit(v, parent, nextLine, nextCol)
	}
}

// Visit writes any entries needed between the current position the next
// position and update the current position.
//
// Java: SourceMapGeneratorV3.MappingTraversal.visit(MappingVisitor, Mapping, int, int)
func (t *SourceMapGeneratorV3_MappingTraversal) Visit(v SourceMapGeneratorV3_MappingVisitor, m *SourceMapGeneratorV3_Mapping,
	nextLine, nextCol int) {
	jutil.CheckState(t.line <= nextLine)
	jutil.CheckState(t.line < nextLine || t.col < nextCol)

	if t.line == nextLine && t.col == nextCol {
		// Nothing to do.
		panic(jutil.NewIllegalStateException())
	}

	v.Visit(m, t.line, t.col, nextLine, nextCol)

	t.line = nextLine
	t.col = nextCol
}

// SourceMapGeneratorV3_LineMapper writes the "mappings" field.
type SourceMapGeneratorV3_LineMapper struct {
	// Go-only: language-adapter — the enclosing SourceMapGeneratorV3.
	outer *SourceMapGeneratorV3

	// The destination.
	out     io.StringWriter
	maxLine int // TODO(johnlenz): This shouldn't be necessary to track.

	previousLine   int
	previousColumn int

	// Previous values used for storing relative ids.
	previousSourceFileId int
	previousSourceLine   int
	previousSourceColumn int
	previousNameId       int

	buf       [4096]byte
	bufLen    int
	buffering bool

	// Go-only: the id tables were filled by prepMappingsParallel and are
	// only read.
	frozenIds bool
}

// Java: SourceMapGeneratorV3.LineMapper(Appendable, int)
func NewSourceMapGeneratorV3_LineMapper(outer *SourceMapGeneratorV3, out io.StringWriter, maxLine int) *SourceMapGeneratorV3_LineMapper {
	l := &SourceMapGeneratorV3_LineMapper{outer: outer, previousLine: -1, previousColumn: 0}
	l.out = out
	l.maxLine = maxLine
	return l
}

func (l *SourceMapGeneratorV3_LineMapper) flush() {
	if l.bufLen > 0 {
		appendBytes(l.out, l.buf[:l.bufLen])
		l.bufLen = 0
	}
}

func (l *SourceMapGeneratorV3_LineMapper) writeByte(b byte) {
	if !l.buffering {
		appendByte(l.out, b)
		return
	}
	if l.bufLen >= len(l.buf) {
		l.flush()
	}
	l.buf[l.bufLen] = b
	l.bufLen++
}

func (l *SourceMapGeneratorV3_LineMapper) writeVLQ(value int32) {
	if !l.buffering {
		Base64VLQ_Encode(l.out, value)
		return
	}
	if l.bufLen+8 > len(l.buf) {
		l.flush()
	}
	value = Base64VLQ_ToVLQSigned(value)
	for {
		digit := value & Base64VLQ_VLQ_BASE_MASK
		value = jutil.Ushr32(value, Base64VLQ_VLQ_BASE_SHIFT)
		if value > 0 {
			digit |= Base64VLQ_VLQ_CONTINUATION_BIT
		}
		l.buf[l.bufLen] = byte(Base64_ToBase64(digit))
		l.bufLen++
		if !(value > 0) {
			break
		}
	}
}

// Visit writes out the appropriate line mapping as each segment is visited.
//
// Java: SourceMapGeneratorV3.LineMapper.visit(Mapping, int, int, int, int)
func (l *SourceMapGeneratorV3_LineMapper) Visit(m *SourceMapGeneratorV3_Mapping, line, col, nextLine, nextCol int) {
	if l.previousLine != line {
		l.previousColumn = 0
	}

	if line != nextLine || col != nextCol {
		// TODO(johnlenz): For some reason, we have mappings beyond the max
		// line. So far they're just null mappings and we can ignore them.
		// (If they're non-null, we assert-fail.)
		if line < l.maxLine {
			if l.previousLine == line { // not the first entry for the line
				l.writeByte(',')
			}
			l.WriteEntry(m, col)
			l.previousLine = line
			l.previousColumn = col
		} else {
			jutil.CheckState(m == nil)
		}
	}

	for i := line; i <= nextLine && i < l.maxLine; i++ {
		if i == nextLine {
			break
		}

		l.CloseLine(false)
		l.OpenLine(false)
	}
}

// WriteEntry writes an entry for the given column (of the generated text)
// and associated mapping. The values are stored as relative to the last
// seen values for each field and encoded as Base64VLQs.
//
// Java: SourceMapGeneratorV3.LineMapper.writeEntry(Mapping, int)
func (l *SourceMapGeneratorV3_LineMapper) WriteEntry(m *SourceMapGeneratorV3_Mapping, column int) {
	// The relative generated column number
	l.writeVLQ(int32(column - l.previousColumn))
	l.previousColumn = column
	if m != nil {
		// The relative source file id
		sourceId := l.outer.sourceIdxToFileId[m.sourceFileIdx]
		if sourceId < 0 && l.frozenIds {
			panic(jutil.NewIllegalStateExceptionMsg("source id not assigned in phase 1"))
		}
		if sourceId < 0 {
			sourceId = int32(len(l.outer.orderedSourceFiles))
			l.outer.orderedSourceFiles = append(l.outer.orderedSourceFiles, l.outer.addedSourceFiles[m.sourceFileIdx])
			l.outer.sourceIdxToFileId[m.sourceFileIdx] = sourceId
		}
		l.writeVLQ(sourceId - int32(l.previousSourceFileId))
		l.previousSourceFileId = int(sourceId)

		// The relative source file line and column
		srcline := int(m.originalLine)
		srcColumn := int(m.originalCol)
		l.writeVLQ(int32(srcline - l.previousSourceLine))
		l.previousSourceLine = srcline

		l.writeVLQ(int32(srcColumn - l.previousSourceColumn))
		l.previousSourceColumn = srcColumn

		if m.originalNameIdx >= 0 {
			// The relative id for the associated symbol name
			nameId := l.outer.symbolIdxToNameId[m.originalNameIdx]
			if nameId < 0 && l.frozenIds {
				panic(jutil.NewIllegalStateExceptionMsg("name id not assigned in phase 1"))
			}
			if nameId < 0 {
				nameId = int32(len(l.outer.orderedSymbolNames))
				l.outer.orderedSymbolNames = append(l.outer.orderedSymbolNames, l.outer.addedSymbolNames[m.originalNameIdx])
				l.outer.symbolIdxToNameId[m.originalNameIdx] = nameId
			}
			l.writeVLQ(nameId - int32(l.previousNameId))
			l.previousNameId = int(nameId)
		}
	}
}

// AppendLineMappings appends the line mapping entries.
//
// Java: SourceMapGeneratorV3.LineMapper.appendLineMappings()
func (l *SourceMapGeneratorV3_LineMapper) AppendLineMappings() {
	nSources := len(l.outer.addedSourceFiles)
	if cap(l.outer.sourceIdxToFileId) < nSources {
		l.outer.sourceIdxToFileId = make([]int32, nSources)
	} else {
		l.outer.sourceIdxToFileId = l.outer.sourceIdxToFileId[:nSources]
	}
	for i := range l.outer.sourceIdxToFileId {
		l.outer.sourceIdxToFileId[i] = -1
	}
	if cap(l.outer.orderedSourceFiles) < nSources {
		l.outer.orderedSourceFiles = make([]string, 0, nSources)
	} else {
		clear(l.outer.orderedSourceFiles)
		l.outer.orderedSourceFiles = l.outer.orderedSourceFiles[:0]
	}
	nSymbols := len(l.outer.addedSymbolNames)
	if cap(l.outer.symbolIdxToNameId) < nSymbols {
		l.outer.symbolIdxToNameId = make([]int32, nSymbols)
	} else {
		l.outer.symbolIdxToNameId = l.outer.symbolIdxToNameId[:nSymbols]
	}
	for i := range l.outer.symbolIdxToNameId {
		l.outer.symbolIdxToNameId[i] = -1
	}
	if cap(l.outer.orderedSymbolNames) == 0 && nSymbols > 0 {
		l.outer.orderedSymbolNames = make([]string, 0, min(nSymbols, 1024))
	} else {
		clear(l.outer.orderedSymbolNames)
		l.outer.orderedSymbolNames = l.outer.orderedSymbolNames[:0]
	}

	l.buffering = true
	// Start the first line.
	l.OpenLine(true)

	t := SourceMapGeneratorV3_MappingTraversal{outer: l.outer}
	t.Traverse(l)

	// And close the final line.
	l.CloseLine(true)
	l.flush()
	l.buffering = false
}

// OpenLine begins the entry for a new line.
//
// Java: SourceMapGeneratorV3.LineMapper.openLine(boolean)
func (l *SourceMapGeneratorV3_LineMapper) OpenLine(firstEntry bool) {
	if firstEntry {
		l.writeByte('"')
	}
}

// CloseLine ends the entry for a line.
//
// Java: SourceMapGeneratorV3.LineMapper.closeLine(boolean)
func (l *SourceMapGeneratorV3_LineMapper) CloseLine(finalEntry bool) {
	l.writeByte(';')
	if finalEntry {
		l.writeByte('"')
	}
}
