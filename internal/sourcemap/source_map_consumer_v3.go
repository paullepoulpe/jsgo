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
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// SourceMapConsumerV3_UNMAPPED marks an absent source file or name id.
const SourceMapConsumerV3_UNMAPPED = -1

// SourceMapConsumerV3 is a class for parsing version 3 of the SourceMap
// format, as produced by the Closure Compiler, etc.
// https://github.com/google/closure-compiler/wiki/Source-Maps
type SourceMapConsumerV3 struct {
	sources        []*string
	sourcesContent []*string
	names          []*string
	file           *string
	lineCount      int

	mappings *SourceMapConsumerV3_Mappings

	// (originalFile path, original line) ==> target mappings
	reverseSourceMapping map[sourceMapConsumerV3_lineKey][]*Mapping_OriginalMapping
	sourceRoot           *string
	extensions           *jutil.LinkedHashMap[string, any]
}

// sourceMapConsumerV3_lineKey is a (nullable Java String, line) map key.
//
// Go-only: language-adapter — reverseSourceMapping is keyed by
// sources[id], which is null for a JSON null in "sources"; a Java HashMap
// keeps null as a key distinct from every string.
type sourceMapConsumerV3_lineKey struct {
	name   string
	line   int
	isNull bool
}

func sourceMapConsumerV3_keyOf(s *string, line int) sourceMapConsumerV3_lineKey {
	if s == nil {
		return sourceMapConsumerV3_lineKey{line: line, isNull: true}
	}
	return sourceMapConsumerV3_lineKey{name: *s, line: line}
}

// NewSourceMapConsumerV3 constructs an empty consumer.
//
// Java: SourceMapConsumerV3()
func NewSourceMapConsumerV3() *SourceMapConsumerV3 {
	return &SourceMapConsumerV3{extensions: jutil.NewLinkedHashMap[string, any]()}
}

// SourceMapConsumerV3_DefaultSourceMapSupplier supplies no maps.
type SourceMapConsumerV3_DefaultSourceMapSupplier struct{}

// Java: SourceMapConsumerV3.DefaultSourceMapSupplier.getSourceMap(String)
func (*SourceMapConsumerV3_DefaultSourceMapSupplier) GetSourceMap(url string) *string {
	return nil
}

// Parse parses the given contents containing a source map. It panics with
// a *SourceMapParseException.
//
// Java: SourceMapConsumerV3.parse(String)
func (c *SourceMapConsumerV3) Parse(contents string) {
	sourceMapObject := SourceMapObjectParser_Parse(contents)
	c.ParseSourceMapObject(sourceMapObject, nil)
}

// ParseSourceMapObject parses the given contents containing a source map.
// It panics with a *SourceMapParseException.
//
// Java: SourceMapConsumerV3.parse(SourceMapObject, SourceMapSupplier)
func (c *SourceMapConsumerV3) ParseSourceMapObject(sourceMapObject *SourceMapObject, sectionSupplier SourceMapSupplier) {
	if sourceMapObject.GetVersion() != 3 {
		panic(NewSourceMapParseException("Unknown version: " + jutil.ToJavaString(sourceMapObject.GetVersion())))
	}

	c.file = sourceMapObject.GetFile()
	if c.file != nil && *c.file == "" {
		panic(NewSourceMapParseException("File entry is empty"))
	}

	if sourceMapObject.GetSections() != nil {
		// Looks like a index map, try to parse it that way.
		c.ParseMetaMap(sourceMapObject, sectionSupplier)
		return
	}

	c.lineCount = sourceMapObject.GetLineCount()
	c.sourceRoot = sourceMapObject.GetSourceRoot()
	c.sources = sourceMapObject.GetSources()
	c.sourcesContent = sourceMapObject.GetSourcesContent()
	c.names = sourceMapObject.GetNames()

	if c.sources == nil {
		// sources.length
		panic(jutil.NewNullPointerException())
	}
	useCompactMappings :=
		(len(c.sources) < 65535) && (c.names == nil || len(c.names) < 65535)

	var estimatedEntries int32
	if c.lineCount >= 0 {
		estimatedEntries = int32(c.lineCount) * 2
	} else {
		estimatedEntries = 1000
	}
	if sourceMapObject.GetMappings() != nil {
		// Estimate the number of entries in the string encoding. A entry has 5
		// values and a separator, for a minimum of 6 characters.
		//
		// It is better to overestimate than underestimate the number of
		// entries - otherwise, we will have to resize the array many times.
		// XJS base chunk, for example, have ~1000 lines but the map has
		// ~266,000 entries. If we start small (e.g. 2000), we need multiple
		// allocations/copies as it grows: 2k -> 4k -> ... -> 512k. Because
		// 266K lies between 256K and 512K, it tops off at 512K and the final
		// trim() will still trigger an Arrays.copyOf to scale back down to
		// 266K! If we start with a bigger estimate, we might still need to
		// trim down at the end, but we save potentially many array copies.
		estimatedEntries = max(estimatedEntries, int32(jutil.StringLength(*sourceMapObject.GetMappings())/6))
	}

	// The value type of each extension is the native JSON type (e.g.
	// JsonObject, or JSONObject when compiled with GWT).
	if sourceMapObject.GetExtensions() == nil {
		// extensions.putAll(null)
		panic(jutil.NewNullPointerException())
	}
	for k, v := range sourceMapObject.GetExtensions().All() {
		c.extensions.Put(k, v)
	}
	c.mappings =
		NewSourceMapConsumerV3_MappingBuilder(
			c, sourceMapObject.GetMappings(), c.lineCount, useCompactMappings, estimatedEntries).
			Build()
}

// ParseMetaMap parses an index map: its sections are merged into a new
// generator, whose map is then parsed by this consumer.
//
// Java: SourceMapConsumerV3.parseMetaMap(SourceMapObject, SourceMapSupplier)
func (c *SourceMapConsumerV3) ParseMetaMap(sourceMapObject *SourceMapObject, sectionSupplier SourceMapSupplier) {
	if sectionSupplier == nil {
		sectionSupplier = &SourceMapConsumerV3_DefaultSourceMapSupplier{}
	}

	jutil.Try(func() {
		if sourceMapObject.GetLineCount() >= 0 ||
			sourceMapObject.GetMappings() != nil ||
			(sourceMapObject.GetSources() != nil && len(sourceMapObject.GetSources()) > 0) ||
			sourceMapObject.GetNames() != nil {
			panic(NewSourceMapParseException("Invalid map format"))
		}

		// Build up a new source map in a new generator using the mappings of
		// this metamap. The new map will be rendered to JSON and then parsed
		// using this consumer.
		generator := NewSourceMapGeneratorV3()
		for _, section := range sourceMapObject.GetSections() {
			mapSectionContents := section.GetSectionValue()
			if section.GetSectionType() == SourceMapSection_SectionType_URL {
				mapSectionContents = sectionSupplier.GetSourceMap(*section.GetSectionValue())
			}
			if mapSectionContents == nil {
				panic(NewSourceMapParseException("Unable to retrieve: " + jutil.ToJavaString(section.GetSectionValue())))
			}
			generator.MergeMapSection(section.GetLine(), section.GetColumn(), *mapSectionContents)
		}

		var sb strings.Builder
		generator.AppendTo(&sb, sourceMapObject.GetFile())
		c.Parse(sb.String())

	}, func(ex jutil.IIOException) {
		panic(NewSourceMapParseException("IO exception: " + ex.Error()))
	})
}

// GetMappingForLine returns the original mapping for the line number and
// column position (both 1-based), or nil.
//
// Java: SourceMapConsumerV3.getMappingForLine(int, int)
func (c *SourceMapConsumerV3) GetMappingForLine(lineNumber, column int) *Mapping_OriginalMapping {
	m, _, ok := c.GetMappingForLineValue(lineNumber, column)
	if !ok {
		return nil
	}
	return &m
}

// GetMappingForLineValue returns the original mapping for the line number and
// column position (both 1-based) by value without heap allocation, along with
// a non-empty identifier pointer (if any) and a boolean indicating whether a
// mapping was found.
func (c *SourceMapConsumerV3) GetMappingForLineValue(lineNumber, column int) (Mapping_OriginalMapping, *string, bool) {
	return c.GetMappingForLineValueCached(lineNumber, column, nil)
}

// GetMappingForLineValueCached is like GetMappingForLineValue, but uses a
// caller-owned SourceMapLineCache when lc != nil so multiple goroutines can
// query a shared SourceMapConsumerV3 concurrently without data races.
func (c *SourceMapConsumerV3) GetMappingForLineValueCached(lineNumber, column int, lc *SourceMapLineCache) (Mapping_OriginalMapping, *string, bool) {
	// Normalize the line and column numbers to 0.
	ln := int32(lineNumber) - 1
	col := int32(column) - 1

	if ln < 0 || int(ln) >= c.mappings.GetParsedLineCount() {
		return Mapping_OriginalMapping{}, nil, false
	}

	entries := c.mappings.lineCached(int(ln), lc)

	// If the line is empty return the previous mapping.
	if len(entries) == 0 {
		return c.getPreviousMappingValueCached(int(ln), lc)
	}

	if entries[0].genCol > col {
		return c.getPreviousMappingValueCached(int(ln), lc)
	}

	index := sourceMapConsumerV3_search(entries, col, 0, len(entries)-1)
	// The message reports the index in the Java int[] of entries.
	jutil.CheckStateMsg1(index >= 0, "unexpected:%s", index*c.mappings.entrySize)
	return c.getOriginalMappingForEntryValue(&entries[index], Mapping_OriginalMapping_Precision_EXACT)
}

// GetOriginalSources returns the sources (elements may be nil).
//
// Java: SourceMapConsumerV3.getOriginalSources()
func (c *SourceMapConsumerV3) GetOriginalSources() []*string {
	if c.sources == nil {
		// Arrays.asList(null)
		panic(jutil.NewNullPointerException())
	}
	return c.sources
}

// GetOriginalSourcesContent returns the sources' contents, or nil.
//
// Java: SourceMapConsumerV3.getOriginalSourcesContent()
func (c *SourceMapConsumerV3) GetOriginalSourcesContent() []*string {
	if c.sourcesContent == nil {
		return nil
	}
	return c.sourcesContent
}

// DropSourcesContent forgets the sources' contents, once the caller knows
// it will not ask for them (an output map without sources content).
//
// Go-only: performance — the contents are a copy of the original sources.
func (c *SourceMapConsumerV3) DropSourcesContent() {
	c.sourcesContent = nil
}

// GetOriginalNames returns the names (elements may be nil).
//
// Java: SourceMapConsumerV3.getOriginalNames()
func (c *SourceMapConsumerV3) GetOriginalNames() []*string {
	if c.names == nil {
		// Arrays.asList(null)
		panic(jutil.NewNullPointerException())
	}
	return c.names
}

// Java: SourceMapConsumerV3.getFile()
func (c *SourceMapConsumerV3) GetFile() *string {
	return c.file
}

// Java: SourceMapConsumerV3.getLineCount()
func (c *SourceMapConsumerV3) GetLineCount() int {
	return c.lineCount
}

// GetReverseMapping returns the target mappings of a source line.
//
// Java: SourceMapConsumerV3.getReverseMapping(String, int, int)
func (c *SourceMapConsumerV3) GetReverseMapping(originalFile string, line, column int) []*Mapping_OriginalMapping {
	// TODO(user): This implementation currently does not make use of the
	// column parameter.

	// Synchronization needs to be handled by callers.
	if c.reverseSourceMapping == nil {
		c.CreateReverseMapping()
	}

	mappings := c.reverseSourceMapping[sourceMapConsumerV3_lineKey{name: originalFile, line: line}]
	if mappings == nil {
		return []*Mapping_OriginalMapping{}
	}
	return mappings
}

// Java: SourceMapConsumerV3.getSourceRoot()
func (c *SourceMapConsumerV3) GetSourceRoot() *string {
	return c.sourceRoot
}

// GetExtensions returns all extensions and their values (which can be any
// json value) in a Map object.
//
// Java: SourceMapConsumerV3.getExtensions()
func (c *SourceMapConsumerV3) GetExtensions() *jutil.LinkedHashMap[string, any] {
	return c.extensions
}

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
	b := &SourceMapConsumerV3_MappingBuilder{outer: outer}
	b.content = NewSourceMapConsumerV3_StringCharIterator(*lineMap)
	b.lineCount = lineCount
	if useCompactMappings {
		b.entrySize = 4
	} else {
		b.entrySize = 5
	}
	b.flatLength = int(estimatedEntries) * b.entrySize
	// An entry's varints take at most as many bytes as its Base64 VLQ
	// characters, so the "mappings" text bounds the stream.
	b.data = make([]byte, 0, len(*lineMap))
	estimatedLines := 1000
	if lineCount >= 0 {
		estimatedLines = lineCount
	}
	b.lineStart = make([]uint32, 0, jutil.IntAdd(estimatedLines, 1))
	return b
}

// Java: SourceMapConsumerV3.MappingBuilder.build()
func (b *SourceMapConsumerV3_MappingBuilder) Build() *SourceMapConsumerV3_Mappings {
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
				temp[entryValues] = b.NextValue()
				entryValues++
			}

			e := b.decodeEntry(temp[:], entryValues)
			b.validateEntry(e, b.line)

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

	return newSourceMapConsumerV3_Mappings(b.data, b.lineStart, b.lineCount, b.entrySize)
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

// decodeEntry resolves the VLQ deltas in vals against the builder's running
// state, applying compact 16-bit ID masking when entrySize == 4.
func (b *SourceMapConsumerV3_MappingBuilder) decodeEntry(vals []int32, entryValues int) sourceMapConsumerV3_entry {
	if b.entrySize > b.flatLength {
		b.flatLength *= 2
		if b.entrySize > b.flatLength {
			panic(jutil.NewArrayIndexOutOfBoundsExceptionForLength(0, b.flatLength))
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
		panic(NewSourceMapParseException(
			"Unexpected number of values for entry:" + jutil.ToJavaString(entryValues)))
	}
	return e
}

// Java: SourceMapConsumerV3.MappingBuilder.validateEntry(int, int)
func (b *SourceMapConsumerV3_MappingBuilder) validateEntry(e sourceMapConsumerV3_entry, line int) {
	jutil.CheckStateMsg2(
		b.lineCount < 0 || line < b.lineCount, "line=%s, lineCount=%s", line, b.lineCount)
	jutil.CheckState(e.srcId == SourceMapConsumerV3_UNMAPPED || int(e.srcId) < len(b.outer.sources))
	if e.nameId != SourceMapConsumerV3_UNMAPPED && b.outer.names == nil {
		// names.length
		panic(jutil.NewNullPointerException())
	}
	jutil.CheckState(e.nameId == SourceMapConsumerV3_UNMAPPED || int(e.nameId) < len(b.outer.names))
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
	var result int32
	var shift int32
	in := b.content
	for {
		c := in.Next()
		digit := Base64_FromBase64(c)
		continuation := (digit & Base64VLQ_VLQ_CONTINUATION_BIT) != 0
		digit &= Base64VLQ_VLQ_BASE_MASK
		result = result + jutil.Shl32(digit, int64(shift))
		shift = shift + Base64VLQ_VLQ_BASE_SHIFT
		if !continuation {
			break
		}
	}
	return Base64VLQ_FromVLQSigned(result)
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

// Java: SourceMapConsumerV3.getPreviousMapping(int)
func (c *SourceMapConsumerV3) GetPreviousMapping(lineNumber int) *Mapping_OriginalMapping {
	m, _, ok := c.getPreviousMappingValue(lineNumber)
	if !ok {
		return nil
	}
	return &m
}

func (c *SourceMapConsumerV3) getPreviousMappingValue(lineNumber int) (Mapping_OriginalMapping, *string, bool) {
	return c.getPreviousMappingValueCached(lineNumber, nil)
}

func (c *SourceMapConsumerV3) getPreviousMappingValueCached(lineNumber int, lc *SourceMapLineCache) (Mapping_OriginalMapping, *string, bool) {
	for {
		if lineNumber == 0 {
			return Mapping_OriginalMapping{}, nil, false
		}
		lineNumber--
		if !c.mappings.lineIsEmpty(lineNumber) {
			break
		}
	}
	entries := c.mappings.lineCached(lineNumber, lc)
	return c.getOriginalMappingForEntryValue(&entries[len(entries)-1], Mapping_OriginalMapping_Precision_APPROXIMATE_LINE)
}

// GetOriginalMappingForEntry creates an "OriginalMapping" object for the
// given entry object.
//
// Java: SourceMapConsumerV3.getOriginalMappingForEntry(int, OriginalMapping.Precision)
func (c *SourceMapConsumerV3) GetOriginalMappingForEntry(entry *sourceMapConsumerV3_entry, precision Mapping_OriginalMapping_Precision) *Mapping_OriginalMapping {
	m, _, ok := c.getOriginalMappingForEntryValue(entry, precision)
	if !ok {
		return nil
	}
	return &m
}

func (c *SourceMapConsumerV3) getOriginalMappingForEntryValue(entry *sourceMapConsumerV3_entry, precision Mapping_OriginalMapping_Precision) (Mapping_OriginalMapping, *string, bool) {
	sourceFileId := entry.srcId
	if sourceFileId == SourceMapConsumerV3_UNMAPPED {
		return Mapping_OriginalMapping{}, nil, false
	}
	// Adjust the line/column here to be start at 1.
	origFile := *c.sources[sourceFileId]
	lineNum := int(entry.srcLine + 1)
	colPos := int(entry.srcCol + 1)
	prec := int32(precision.GetNumber())
	bitField := int32(mapping_originalFileBit | mapping_lineNumberBit | mapping_columnPositionBit | mapping_precisionBit)
	var identifier string
	var identifierPtr *string
	nameId := entry.nameId
	if nameId != SourceMapConsumerV3_UNMAPPED {
		namePtr := c.names[nameId]
		identifier = *namePtr
		bitField |= mapping_identifierBit
		if identifier != "" {
			identifierPtr = namePtr
		}
	}
	return Mapping_OriginalMapping{
		bitField0_:      bitField,
		originalFile_:   origFile,
		lineNumber_:     lineNum,
		columnPosition_: colPos,
		identifier_:     identifier,
		precision_:      prec,
	}, identifierPtr, true
}

// CreateReverseMapping reverses the source map; the created mapping will
// allow us to quickly go from a source file and line number to a collection
// of target OriginalMappings.
//
// Java: SourceMapConsumerV3.createReverseMapping()
func (c *SourceMapConsumerV3) CreateReverseMapping() {
	c.reverseSourceMapping = make(map[sourceMapConsumerV3_lineKey][]*Mapping_OriginalMapping)

	for targetLine := 0; targetLine < c.mappings.GetParsedLineCount(); targetLine++ {
		entries := c.mappings.line(targetLine)
		for i := range entries {
			entry := &entries[i]
			sourceFileId := entry.srcId
			sourceLine := entry.srcLine
			if sourceFileId != SourceMapConsumerV3_UNMAPPED && sourceLine != SourceMapConsumerV3_UNMAPPED {
				key := sourceMapConsumerV3_keyOf(c.sources[sourceFileId], int(sourceLine))
				mapping := Mapping_OriginalMapping_NewBuilder().
					SetLineNumber(targetLine).
					SetColumnPosition(int(entry.genCol)).
					Build()
				c.reverseSourceMapping[key] = append(c.reverseSourceMapping[key], mapping)
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

// SourceMapConsumerV3_EntryVisitor receives the mappings of a parsed map.
type SourceMapConsumerV3_EntryVisitor interface {
	// Java: SourceMapConsumerV3.EntryVisitor.visit(String, String, FilePosition, FilePosition, FilePosition)
	Visit(sourceName *string,
		symbolName *string,
		sourceStartPosition *FilePosition,
		startPosition *FilePosition,
		endPosition *FilePosition)
}

// VisitMappings visits every mapped segment, in order.
//
// Java: SourceMapConsumerV3.visitMappings(EntryVisitor)
func (c *SourceMapConsumerV3) VisitMappings(visitor SourceMapConsumerV3_EntryVisitor) {
	pending := false
	var sourceName *string
	var symbolName *string
	var sourceStartPosition *FilePosition
	var startPosition *FilePosition

	for i := 0; i < c.mappings.GetParsedLineCount(); i++ {
		if !c.mappings.lineIsEmpty(i) {
			entries := c.mappings.line(i)
			for j := range entries {
				entry := &entries[j]
				if pending {
					endPosition := NewFilePosition(i, int(entry.genCol))
					visitor.Visit(
						sourceName,
						symbolName,
						sourceStartPosition,
						startPosition,
						endPosition)
					pending = false
				}

				sourceFileId := entry.srcId
				if sourceFileId != SourceMapConsumerV3_UNMAPPED {
					pending = true
					sourceName = c.sources[sourceFileId]
					nameId := entry.nameId
					if nameId != SourceMapConsumerV3_UNMAPPED {
						symbolName = c.names[nameId]
					} else {
						symbolName = nil
					}
					sourceStartPosition =
						NewFilePosition(int(entry.srcLine), int(entry.srcCol))
					startPosition = NewFilePosition(i, int(entry.genCol))
				}
			}
		}
	}
	// Complete pending entry if any.
	if pending {
		// Given that this is the last entry and we don't know how much of the
		// generated file left after that entry - make it of length 1.
		endPosition :=
			NewFilePosition(startPosition.GetLine(), jutil.IntAdd(startPosition.GetColumn(), 1))
		visitor.Visit(sourceName, symbolName, sourceStartPosition, startPosition, endPosition)
	}
}
