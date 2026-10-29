// Port of com/google/debugging/sourcemap/SourceMapConsumerV3.java @ closure-compiler v20261021

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
	"slices"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// SourceMapConsumerV3_MappingBuilder decodes the "mappings" field.
type SourceMapConsumerV3_MappingBuilder struct {
	// Go-only: language-adapter — the enclosing SourceMapConsumerV3.
	outer *SourceMapConsumerV3

	content *SourceMapConsumerV3_StringCharIterator

	// Data for the Mappings we're building.
	//
	// Go-only: performance — the Java int[] of entries is the delta-varint
	// stream of SourceMapConsumerV3_Mappings; prevInLine holds the previous
	// decoded entry of the current line that deltas are encoded against.
	//
	// flatLength tracks the length the Java int[] would have (it starts at
	// the estimate and doubles): a zero estimate makes the first entry fail
	// like Java's out-of-bounds store.
	flatLength int
	prevInLine sourceMapConsumerV3_entry
	data       []byte
	lineStart  []uint32
	lineCount  int
	entrySize  int

	// Parsing state
	line              int
	previousCol       int32
	previousSrcId     int32
	previousSrcLine   int32
	previousSrcColumn int32
	previousNameId    int32
}

// SourceMapConsumerV3_MappingBuilder_MAX_ENTRY_VALUES is the most values an
// entry has.
const SourceMapConsumerV3_MappingBuilder_MAX_ENTRY_VALUES = 5

// NewSourceMapConsumerV3_MappingBuilder constructs a builder over lineMap.
//
// Java: SourceMapConsumerV3.MappingBuilder(String, int, boolean, int)
func NewSourceMapConsumerV3_MappingBuilder(
	outer *SourceMapConsumerV3, lineMap *string, lineCount int, useCompactMappings bool, estimatedEntries int32) *SourceMapConsumerV3_MappingBuilder {
	b, err := tryNewMappingBuilder(outer, makeOptionalString(lineMap), lineCount, useCompactMappings, estimatedEntries)
	if err != nil {
		panic(err)
	}
	return b
}

func tryNewMappingBuilder(
	outer *SourceMapConsumerV3, lineMap optionalString, lineCount int, useCompactMappings bool, estimatedEntries int32) (*SourceMapConsumerV3_MappingBuilder, jutil.Throwable) {
	if !lineMap.present {
		return nil, jutil.NewNullPointerException()
	}
	b := &SourceMapConsumerV3_MappingBuilder{outer: outer}
	b.content = NewSourceMapConsumerV3_StringCharIterator(lineMap.text)
	b.lineCount = lineCount
	if useCompactMappings {
		b.entrySize = 4
	} else {
		b.entrySize = 5
	}
	b.flatLength = int(estimatedEntries) * b.entrySize
	// An entry's varints take at most as many bytes as its Base64 VLQ
	// characters, so the "mappings" text bounds the stream.
	b.data = make([]byte, 0, len(lineMap.text))
	estimatedLines := 1000
	if lineCount >= 0 {
		estimatedLines = lineCount
	}
	b.lineStart = make([]uint32, 0, jutil.IntAdd(estimatedLines, 1))
	return b, nil
}

// Java: SourceMapConsumerV3.MappingBuilder.build()
func (b *SourceMapConsumerV3_MappingBuilder) Build() *SourceMapConsumerV3_Mappings {
	m, err := b.tryBuild()
	if err != nil {
		panic(err)
	}
	return m
}

func (b *SourceMapConsumerV3_MappingBuilder) tryBuild() (*SourceMapConsumerV3_Mappings, jutil.Throwable) {
	var temp [SourceMapConsumerV3_MappingBuilder_MAX_ENTRY_VALUES]int32

	entriesCount := 0
	linesCount := 0

	b.SetLineStart(0, 0)

	for b.content.HasNext() {
		// ';' denotes a new line.
		if b.TryConsumeToken(';') {
			// The line is complete, store the result
			linesCount++
			b.SetLineStart(linesCount, len(b.data))
			b.line++
			b.previousCol = 0
			b.prevInLine = sourceMapConsumerV3_entry{}
		} else {
			// grab the next entry for the current line.
			entryValues := 0
			for !b.EntryComplete() {
				v, err := b.tryNextValue()
				if err != nil {
					return nil, err
				}
				if entryValues >= SourceMapConsumerV3_MappingBuilder_MAX_ENTRY_VALUES {
					return nil, jutil.NewArrayIndexOutOfBoundsExceptionForLength(
						entryValues, SourceMapConsumerV3_MappingBuilder_MAX_ENTRY_VALUES)
				}
				temp[entryValues] = v
				entryValues++
			}

			e, err := b.tryDecodeEntry(temp[:], entryValues)
			if err != nil {
				return nil, err
			}
			if err := b.tryValidateEntry(e, b.line); err != nil {
				return nil, err
			}

			b.previousCol = e.genCol
			if entryValues >= 4 {
				b.previousSrcId = e.srcId
				b.previousSrcLine = e.srcLine
				b.previousSrcColumn = e.srcCol
			}
			if entryValues == 5 {
				b.previousNameId = e.nameId
			}

			b.appendEntry(entryValues, e)
			entriesCount++

			// Consume the separating token, if there is one.
			b.TryConsumeToken(',')
		}
	}

	// Some source map generator (e.g.UglifyJS) generates lines without a
	// trailing line separator. So add the rest of the content.
	if len(b.data) > int(b.GetLineStart(linesCount)) {
		linesCount++
		b.SetLineStart(linesCount, len(b.data))
	}

	b.Trim(entriesCount, linesCount)

	return newSourceMapConsumerV3_Mappings(b.data, b.lineStart, b.lineCount, b.entrySize), nil
}

// appendEntry encodes e onto the stream as deltas against the previous entry
// of the line. See SourceMapConsumerV3_Mappings.
func (b *SourceMapConsumerV3_MappingBuilder) appendEntry(entryValues int, e sourceMapConsumerV3_entry) {
	kind := uint32(0)
	switch entryValues {
	case 4:
		kind = 1
	case 5:
		kind = 2
	}
	b.data = sourceMapConsumerV3_appendVarint(b.data, uint64(sourceMapConsumerV3_zigzag(e.genCol-b.prevInLine.genCol))<<2|uint64(kind))
	b.prevInLine.genCol = e.genCol
	if kind == 0 {
		return
	}
	prev := &b.prevInLine
	b.data = sourceMapConsumerV3_appendVarint(b.data, uint64(sourceMapConsumerV3_zigzag(e.srcId-prev.srcId)))
	b.data = sourceMapConsumerV3_appendVarint(b.data, uint64(sourceMapConsumerV3_zigzag(e.srcLine-prev.srcLine)))
	b.data = sourceMapConsumerV3_appendVarint(b.data, uint64(sourceMapConsumerV3_zigzag(e.srcCol-prev.srcCol)))
	if kind == 2 {
		b.data = sourceMapConsumerV3_appendVarint(b.data, uint64(sourceMapConsumerV3_zigzag(e.nameId-prev.nameId)))
	} else {
		e.nameId = prev.nameId
	}
	b.prevInLine = e
}

// Java: SourceMapConsumerV3.MappingBuilder.getLineStart(int)
func (b *SourceMapConsumerV3_MappingBuilder) GetLineStart(line int) uint32 {
	return b.lineStart[line]
}

// SetLineStart records where line's entries start in the stream (always
// the next line, in order).
//
// Java: SourceMapConsumerV3.MappingBuilder.setLineStart(int, int)
func (b *SourceMapConsumerV3_MappingBuilder) SetLineStart(line, dataOffset int) {
	jutil.CheckState(line == len(b.lineStart))
	b.lineStart = append(b.lineStart, uint32(dataOffset))
}

func maskCompactId(id int32) int32 {
	if id == -1 || (id&0xFFFF) == 0xFFFF {
		return SourceMapConsumerV3_UNMAPPED
	}
	return id & 0xFFFF
}

// tryDecodeEntry resolves the VLQ deltas in vals against the builder's running
// state, applying compact 16-bit ID masking when entrySize == 4.
func (b *SourceMapConsumerV3_MappingBuilder) tryDecodeEntry(vals []int32, entryValues int) (sourceMapConsumerV3_entry, jutil.Throwable) {
	if b.entrySize > b.flatLength {
		b.flatLength *= 2
		if b.entrySize > b.flatLength {
			return sourceMapConsumerV3_entry{}, jutil.NewArrayIndexOutOfBoundsExceptionForLength(0, b.flatLength)
		}
	}
	e := sourceMapConsumerV3_entry{
		srcLine: SourceMapConsumerV3_UNMAPPED,
		srcCol:  SourceMapConsumerV3_UNMAPPED,
		srcId:   SourceMapConsumerV3_UNMAPPED,
		nameId:  SourceMapConsumerV3_UNMAPPED,
	}
	switch entryValues {
	case 1:
		e.genCol = vals[0] + b.previousCol
	case 4:
		e.genCol = vals[0] + b.previousCol
		e.srcId = vals[1] + b.previousSrcId
		e.srcLine = vals[2] + b.previousSrcLine
		e.srcCol = vals[3] + b.previousSrcColumn
		if b.entrySize == 4 {
			e.srcId = maskCompactId(e.srcId)
		}
	case 5:
		e.genCol = vals[0] + b.previousCol
		e.srcId = vals[1] + b.previousSrcId
		e.srcLine = vals[2] + b.previousSrcLine
		e.srcCol = vals[3] + b.previousSrcColumn
		e.nameId = vals[4] + b.previousNameId
		if b.entrySize == 4 {
			e.srcId = maskCompactId(e.srcId)
			e.nameId = maskCompactId(e.nameId)
		}
	default:
		return e, NewSourceMapParseException(
			"Unexpected number of values for entry:" + jutil.ToJavaString(entryValues))
	}
	return e, nil
}

// Java: SourceMapConsumerV3.MappingBuilder.validateEntry(int, int)
func (b *SourceMapConsumerV3_MappingBuilder) tryValidateEntry(e sourceMapConsumerV3_entry, line int) jutil.Throwable {
	if !(b.lineCount < 0 || line < b.lineCount) {
		return jutil.NewIllegalStateExceptionMsg(jutil.LenientFormat("line=%s, lineCount=%s", line, b.lineCount))
	}
	if !(e.srcId == SourceMapConsumerV3_UNMAPPED || int(e.srcId) < len(b.outer.sources)) {
		return jutil.NewIllegalStateException()
	}
	if e.nameId != SourceMapConsumerV3_UNMAPPED && b.outer.names == nil {
		// names.length
		return jutil.NewNullPointerException()
	}
	if !(e.nameId == SourceMapConsumerV3_UNMAPPED || int(e.nameId) < len(b.outer.names)) {
		return jutil.NewIllegalStateException()
	}
	return nil
}

// Java: SourceMapConsumerV3.MappingBuilder.trim(int, int)
func (b *SourceMapConsumerV3_MappingBuilder) Trim(entriesCount int, linesCount int) {
	_ = entriesCount
	if len(b.lineStart) > linesCount+1 {
		b.lineStart = b.lineStart[:linesCount+1]
	}
	b.lineCount = linesCount
	// Give the over-estimated capacity back.
	if cap(b.data)-len(b.data) > len(b.data)/8+64 {
		b.data = slices.Clone(b.data)
	}
	if cap(b.lineStart)-len(b.lineStart) > len(b.lineStart)/8+64 {
		b.lineStart = slices.Clone(b.lineStart)
	}
}

// Java: SourceMapConsumerV3.MappingBuilder.tryConsumeToken(char)
func (b *SourceMapConsumerV3_MappingBuilder) TryConsumeToken(token uint16) bool {
	if b.content.HasNext() && b.content.Peek() == token {
		// consume the comma
		b.content.Next()
		return true
	}
	return false
}

// Java: SourceMapConsumerV3.MappingBuilder.entryComplete()
func (b *SourceMapConsumerV3_MappingBuilder) EntryComplete() bool {
	if !b.content.HasNext() {
		return true
	}

	c := b.content.Peek()
	return c == ';' || c == ','
}

// Java: SourceMapConsumerV3.MappingBuilder.nextValue()
func (b *SourceMapConsumerV3_MappingBuilder) NextValue() int32 {
	v, err := b.tryNextValue()
	if err != nil {
		panic(err)
	}
	return v
}

func (b *SourceMapConsumerV3_MappingBuilder) tryNextValue() (int32, jutil.Throwable) {
	var result int32
	var shift int32
	in := b.content
	for {
		if in.current >= in.length {
			return 0, jutil.NewStringIndexOutOfBoundsExceptionForIndex(in.current, in.length)
		}
		c := in.Next()
		if c >= 256 {
			return 0, jutil.NewArrayIndexOutOfBoundsExceptionForLength(int(c), 256)
		}
		digit := int32(Base64_BASE64_DECODE_MAP[c])
		continuation := (digit & Base64VLQ_VLQ_CONTINUATION_BIT) != 0
		digit &= Base64VLQ_VLQ_BASE_MASK
		result = result + jutil.Shl32(digit, int64(shift))
		shift = shift + Base64VLQ_VLQ_BASE_SHIFT
		if !continuation {
			break
		}
	}
	return Base64VLQ_FromVLQSigned(result), nil
}

// sourceMapConsumerV3_search performs a binary search on the line's
// entries to find a section that covers the target column.
//
// Java: SourceMapConsumerV3.search(int, int, int)
func sourceMapConsumerV3_search(entries []sourceMapConsumerV3_entry, target int32, start, end int) int {
	for {
		mid := (end-start)/2 + start
		compare := entries[mid].genCol - target
		if compare == 0 {
			return mid
		} else if compare < 0 {
			// it is in the upper half
			start = mid + 1
			if start > end {
				return end
			}
		} else {
			// it is in the lower half
			end = mid - 1
			if end < start {
				return end
			}
		}
	}
}

// SourceMapConsumerV3_StringCharIterator is a implementation of the
// Base64VLQ CharIterator used for decoding the mappings encoded in the JSON
// string.
type SourceMapConsumerV3_StringCharIterator struct {
	// Go-only: language-adapter — the Java String as UTF-16 code units, so
	// that charAt is O(1); when content is pure ASCII (all valid Base64VLQ
	// strings), ascii holds it directly without a []uint16 allocation.
	ascii   string
	content []uint16
	length  int
	current int
}

// NewSourceMapConsumerV3_StringCharIterator constructs an iterator over
// content.
//
// Java: SourceMapConsumerV3.StringCharIterator(String)
func NewSourceMapConsumerV3_StringCharIterator(content string) *SourceMapConsumerV3_StringCharIterator {
	if jutil.IsASCII(content) {
		return &SourceMapConsumerV3_StringCharIterator{
			ascii:  content,
			length: len(content),
		}
	}
	u := jutil.StringToUTF16(content)
	return &SourceMapConsumerV3_StringCharIterator{
		content: u,
		length:  len(u),
	}
}

// Java: SourceMapConsumerV3.StringCharIterator.next()
func (it *SourceMapConsumerV3_StringCharIterator) Next() uint16 {
	i := it.current
	it.current++
	if it.content != nil {
		return it.content[i]
	}
	return uint16(it.ascii[i])
}

// Java: SourceMapConsumerV3.StringCharIterator.peek()
func (it *SourceMapConsumerV3_StringCharIterator) Peek() uint16 {
	if it.content != nil {
		return it.content[it.current]
	}
	return uint16(it.ascii[it.current])
}

// Java: SourceMapConsumerV3.StringCharIterator.hasNext()
func (it *SourceMapConsumerV3_StringCharIterator) HasNext() bool {
	return it.current < it.length
}

// SourceMapConsumerV3_Mappings is an immutable container for mapping
// entries.
//
// Go-only: performance — Java keeps every entry as 4 or 5 ints (16-20
// bytes; 1.4M entries on a bundle of mapped libraries). Here the entries
// are a byte stream of varints, line by line: an entry is
// zigzag(generatedColumn delta) << 2 | kind, with kind 0 for an unmapped
// entry, 1 for a mapped one and 2 for a mapped one with a name, followed
// for kinds 1 and 2 by the zigzag deltas of sourceFileId, sourceLine and
// sourceColumn, and for kind 2 by the nameId delta. Deltas are against the
// previous entry of the same line (zero at the start of a line), so a line
// decodes on its own; a lookup decodes its line into a cache that stays
// valid while the lookups go on in the same line, which is how the code
// printer asks.
type SourceMapConsumerV3_Mappings struct {
	data []byte
	// lineStart[i] is the offset of the first entry for line i in data. The
	// length of lineStart is lineCount + 1, where lineStart[lineCount] is
	// the length of data.
	lineStart []uint32
	lineCount int
	entrySize int

	cachedLine int
	cache      []sourceMapConsumerV3_entry
}

// sourceMapConsumerV3_entry is one decoded entry: the values of the Java
// int[] entry, with SourceMapConsumerV3_UNMAPPED for absent ones.
type sourceMapConsumerV3_entry struct {
	genCol  int32
	srcLine int32
	srcCol  int32
	srcId   int32
	nameId  int32
}

func newSourceMapConsumerV3_Mappings(data []byte, lineStart []uint32, lineCount int, entrySize int) *SourceMapConsumerV3_Mappings {
	return &SourceMapConsumerV3_Mappings{data: data, lineStart: lineStart, lineCount: lineCount, entrySize: entrySize, cachedLine: -1}
}

// Java: SourceMapConsumerV3.Mappings.getParsedLineCount()
func (m *SourceMapConsumerV3_Mappings) GetParsedLineCount() int {
	return m.lineCount
}

func (m *SourceMapConsumerV3_Mappings) lineIsEmpty(line int) bool {
	return m.lineStart[line] == m.lineStart[line+1]
}

// SourceMapLineCache holds per-caller line-decoding state so multiple
// goroutines can query a shared SourceMapConsumerV3 concurrently without
// mutating SourceMapConsumerV3_Mappings.cachedLine / cache.
type SourceMapLineCache struct {
	mappings   *SourceMapConsumerV3_Mappings
	cachedLine int
	cache      []sourceMapConsumerV3_entry
}

// Reset clears the cached mapping pointer and line index while retaining backing capacity.
func (c *SourceMapLineCache) Reset() {
	c.mappings = nil
	c.cachedLine = -1
	c.cache = c.cache[:0]
}

// line returns the decoded entries of a line; the slice is only valid until
// the next call.
func (m *SourceMapConsumerV3_Mappings) line(line int) []sourceMapConsumerV3_entry {
	if line == m.cachedLine {
		return m.cache
	}
	entries := m.decodeLineInto(line, m.cache[:0])
	m.cache = entries
	m.cachedLine = line
	return entries
}

func (m *SourceMapConsumerV3_Mappings) lineCached(line int, lc *SourceMapLineCache) []sourceMapConsumerV3_entry {
	if lc == nil {
		return m.line(line)
	}
	if lc.mappings == m && line == lc.cachedLine {
		return lc.cache
	}
	entries := m.decodeLineInto(line, lc.cache[:0])
	lc.mappings = m
	lc.cachedLine = line
	lc.cache = entries
	return entries
}

func (m *SourceMapConsumerV3_Mappings) decodeLineInto(line int, entries []sourceMapConsumerV3_entry) []sourceMapConsumerV3_entry {
	data := m.data[m.lineStart[line]:m.lineStart[line+1]]
	var prev sourceMapConsumerV3_entry
	for len(data) > 0 {
		var head uint64
		head, data = sourceMapConsumerV3_readVarint(data)
		e := sourceMapConsumerV3_entry{
			genCol:  prev.genCol + sourceMapConsumerV3_unzigzag(uint32(head>>2)),
			srcLine: SourceMapConsumerV3_UNMAPPED,
			srcCol:  SourceMapConsumerV3_UNMAPPED,
			srcId:   SourceMapConsumerV3_UNMAPPED,
			nameId:  SourceMapConsumerV3_UNMAPPED,
		}
		prev.genCol = e.genCol
		if kind := head & 3; kind != 0 {
			var v uint64
			v, data = sourceMapConsumerV3_readVarint(data)
			e.srcId = prev.srcId + sourceMapConsumerV3_unzigzag(uint32(v))
			v, data = sourceMapConsumerV3_readVarint(data)
			e.srcLine = prev.srcLine + sourceMapConsumerV3_unzigzag(uint32(v))
			v, data = sourceMapConsumerV3_readVarint(data)
			e.srcCol = prev.srcCol + sourceMapConsumerV3_unzigzag(uint32(v))
			if kind == 2 {
				v, data = sourceMapConsumerV3_readVarint(data)
				e.nameId = prev.nameId + sourceMapConsumerV3_unzigzag(uint32(v))
				prev = e
			} else {
				// The name delta base carries over an entry without a name.
				nameId := prev.nameId
				prev = e
				prev.nameId = nameId
			}
		}
		entries = append(entries, e)
	}
	return entries
}

func sourceMapConsumerV3_zigzag(v int32) uint32 {
	return uint32((v << 1) ^ (v >> 31))
}

func sourceMapConsumerV3_unzigzag(u uint32) int32 {
	return int32(u>>1) ^ -int32(u&1)
}

func sourceMapConsumerV3_appendVarint(data []byte, v uint64) []byte {
	for v >= 0x80 {
		data = append(data, byte(v)|0x80)
		v >>= 7
	}
	return append(data, byte(v))
}

func sourceMapConsumerV3_readVarint(data []byte) (uint64, []byte) {
	var v uint64
	for i, b := range data {
		v |= uint64(b&0x7F) << (7 * uint(i))
		if b < 0x80 {
			return v, data[i+1:]
		}
	}
	return v, nil
}
