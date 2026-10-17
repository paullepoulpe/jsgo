// Port of com/google/debugging/sourcemap/SourceMapConsumerV3.java @ closure-compiler v20260415

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

	// originalFile path ==> original line ==> target mappings
	reverseSourceMapping *jutil.LinkedHashMap[sourceMapConsumerV3_fileKey,
		*jutil.LinkedHashMap[int, *jutil.ArrayList[*Mapping_OriginalMapping]]]
	sourceRoot *string
	extensions *jutil.LinkedHashMap[string, any]
}

// sourceMapConsumerV3_fileKey is a nullable Java String as a map key.
//
// Go-only: language-adapter — reverseSourceMapping is keyed by
// sources[id], which is null for a JSON null in "sources"; a Java HashMap
// keeps null as a key distinct from every string.
type sourceMapConsumerV3_fileKey struct {
	name   string
	isNull bool
}

func sourceMapConsumerV3_keyOf(s *string) sourceMapConsumerV3_fileKey {
	if s == nil {
		return sourceMapConsumerV3_fileKey{isNull: true}
	}
	return sourceMapConsumerV3_fileKey{name: *s}
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

	c.mappings = NewSourceMapConsumerV3_Mappings(c, c.lineCount, useCompactMappings, estimatedEntries)

	// The value type of each extension is the native JSON type (e.g.
	// JsonObject, or JSONObject when compiled with GWT).
	if sourceMapObject.GetExtensions() == nil {
		// extensions.putAll(null)
		panic(jutil.NewNullPointerException())
	}
	for k, v := range sourceMapObject.GetExtensions().All() {
		c.extensions.Put(k, v)
	}
	NewSourceMapConsumerV3_MappingBuilder(c, sourceMapObject.GetMappings()).Build()
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
			sourceMapObject.GetSources() != nil ||
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
	// Normalize the line and column numbers to 0.
	ln := int32(lineNumber) - 1
	col := int32(column) - 1

	if ln < 0 || int(ln) >= c.mappings.GetParsedLineCount() {
		return nil
	}

	start := c.mappings.GetLineStart(int(ln))
	end := c.mappings.GetLineStart(int(ln) + 1)

	// If the line is empty return the previous mapping.
	if start == end {
		return c.GetPreviousMapping(int(ln))
	}

	if c.mappings.AtIndex(start).GetGeneratedColumn() > col {
		return c.GetPreviousMapping(int(ln))
	}

	index := c.Search(col, start, end-c.mappings.GetEntrySize())
	jutil.CheckStateMsg1(index >= 0, "unexpected:%s", index)
	return c.GetOriginalMappingForEntry(index, Mapping_OriginalMapping_Precision_EXACT)
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

	sourceLineToCollectionMap := c.reverseSourceMapping.Get(sourceMapConsumerV3_keyOf(&originalFile))

	if sourceLineToCollectionMap == nil {
		return []*Mapping_OriginalMapping{}
	} else {
		mappings := sourceLineToCollectionMap.Get(line)

		if mappings == nil {
			return []*Mapping_OriginalMapping{}
		} else {
			return mappings.Slice()
		}
	}
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

	content           *SourceMapConsumerV3_StringCharIterator
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
// Java: SourceMapConsumerV3.MappingBuilder(String)
func NewSourceMapConsumerV3_MappingBuilder(outer *SourceMapConsumerV3, lineMap *string) *SourceMapConsumerV3_MappingBuilder {
	b := &SourceMapConsumerV3_MappingBuilder{outer: outer}
	b.content = NewSourceMapConsumerV3_StringCharIterator(*lineMap)
	return b
}

// Java: SourceMapConsumerV3.MappingBuilder.build()
func (b *SourceMapConsumerV3_MappingBuilder) Build() {
	var temp [SourceMapConsumerV3_MappingBuilder_MAX_ENTRY_VALUES]int32

	entriesCount := 0
	linesCount := 0

	mappings := b.outer.mappings
	mappings.SetLineStart(0, 0)

	for b.content.HasNext() {
		// ';' denotes a new line.
		if b.TryConsumeToken(';') {
			// The line is complete, store the result
			linesCount++
			mappings.SetLineStart(linesCount, entriesCount)
			b.line++
			b.previousCol = 0
		} else {
			// grab the next entry for the current line.
			entryValues := 0
			for !b.EntryComplete() {
				temp[entryValues] = b.NextValue()
				entryValues++
			}

			mappings.SetEntry(
				entriesCount,
				temp[:],
				entryValues,
				b.previousCol,
				b.previousSrcId,
				b.previousSrcLine,
				b.previousSrcColumn,
				b.previousNameId)
			mappings.ValidateEntry(entriesCount, b.line)

			mappings.AtIndex(entriesCount)
			b.previousCol = mappings.GetGeneratedColumn()
			if entryValues >= 4 {
				b.previousSrcId = mappings.GetSourceFileId()
				b.previousSrcLine = mappings.GetSourceLine()
				b.previousSrcColumn = mappings.GetSourceColumn()
			}
			if entryValues == 5 {
				b.previousNameId = mappings.GetNameId()
			}

			entriesCount += mappings.GetEntrySize()

			// Consume the separating token, if there is one.
			b.TryConsumeToken(',')
		}
	}

	// Some source map generator (e.g.UglifyJS) generates lines without a
	// trailing line separator. So add the rest of the content.
	if entriesCount > mappings.GetLineStart(linesCount) {
		linesCount++
		mappings.SetLineStart(linesCount, entriesCount)
	}

	mappings.Trim(entriesCount, linesCount)
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
	return Base64VLQ_Decode(b.content)
}

// Search performs a binary search on the array to find a section that
// covers the target column.
//
// Java: SourceMapConsumerV3.search(int, int, int)
func (c *SourceMapConsumerV3) Search(target int32, start, end int) int {
	entrySize := c.mappings.GetEntrySize()
	for {
		mid := ((end-start)/(entrySize*2))*entrySize + start
		compare := c.mappings.AtIndex(mid).GetGeneratedColumn() - target
		if compare == 0 {
			return mid
		} else if compare < 0 {
			// it is in the upper half
			start = mid + entrySize
			if start > end {
				return end
			}
		} else {
			// it is in the lower half
			end = mid - entrySize
			if end < start {
				return end
			}
		}
	}
}

// Java: SourceMapConsumerV3.getPreviousMapping(int)
func (c *SourceMapConsumerV3) GetPreviousMapping(lineNumber int) *Mapping_OriginalMapping {
	for {
		if lineNumber == 0 {
			return nil
		}
		lineNumber--
		if c.mappings.GetLineStart(lineNumber) != c.mappings.GetLineStart(lineNumber+1) {
			break
		}
	}
	index := c.mappings.GetLineStart(lineNumber+1) - c.mappings.GetEntrySize()
	return c.GetOriginalMappingForEntry(index, Mapping_OriginalMapping_Precision_APPROXIMATE_LINE)
}

// GetOriginalMappingForEntry creates an "OriginalMapping" object for the
// given entry object.
//
// Java: SourceMapConsumerV3.getOriginalMappingForEntry(int, OriginalMapping.Precision)
func (c *SourceMapConsumerV3) GetOriginalMappingForEntry(index int, precision Mapping_OriginalMapping_Precision) *Mapping_OriginalMapping {
	c.mappings.AtIndex(index)
	sourceFileId := c.mappings.GetSourceFileId()
	if sourceFileId == SourceMapConsumerV3_UNMAPPED {
		return nil
	} else {
		// Adjust the line/column here to be start at 1.
		x := Mapping_OriginalMapping_NewBuilder().
			SetOriginalFile(*c.sources[sourceFileId]).
			SetLineNumber(int(c.mappings.GetSourceLine() + 1)).
			SetColumnPosition(int(c.mappings.GetSourceColumn() + 1)).
			SetPrecision(precision)
		nameId := c.mappings.GetNameId()
		if nameId != SourceMapConsumerV3_UNMAPPED {
			x.SetIdentifier(*c.names[nameId])
		}
		return x.Build()
	}
}

// CreateReverseMapping reverses the source map; the created mapping will
// allow us to quickly go from a source file and line number to a collection
// of target OriginalMappings.
//
// Java: SourceMapConsumerV3.createReverseMapping()
func (c *SourceMapConsumerV3) CreateReverseMapping() {
	c.reverseSourceMapping = jutil.NewLinkedHashMap[sourceMapConsumerV3_fileKey,
		*jutil.LinkedHashMap[int, *jutil.ArrayList[*Mapping_OriginalMapping]]]()

	for targetLine := 0; targetLine < c.mappings.GetParsedLineCount(); targetLine++ {
		start := c.mappings.GetLineStart(targetLine)
		end := c.mappings.GetLineStart(targetLine + 1)

		for i := start; i < end; i += c.mappings.GetEntrySize() {
			c.mappings.AtIndex(i)
			sourceFileId := c.mappings.GetSourceFileId()
			sourceLine := c.mappings.GetSourceLine()
			if sourceFileId != SourceMapConsumerV3_UNMAPPED && sourceLine != SourceMapConsumerV3_UNMAPPED {
				originalFile := c.sources[sourceFileId]

				lineToCollectionMap := c.reverseSourceMapping.ComputeIfAbsent(
					sourceMapConsumerV3_keyOf(originalFile),
					func(k sourceMapConsumerV3_fileKey) *jutil.LinkedHashMap[int, *jutil.ArrayList[*Mapping_OriginalMapping]] {
						return jutil.NewLinkedHashMap[int, *jutil.ArrayList[*Mapping_OriginalMapping]]()
					})

				if !lineToCollectionMap.ContainsKey(int(sourceLine)) {
					lineToCollectionMap.Put(int(sourceLine), jutil.NewArrayListWithCapacity[*Mapping_OriginalMapping](1))
				}

				mappingsLine := lineToCollectionMap.Get(int(sourceLine))

				builder := Mapping_OriginalMapping_NewBuilder().
					SetLineNumber(targetLine).
					SetColumnPosition(int(c.mappings.GetGeneratedColumn()))

				mappingsLine.Add(builder.Build())
			}
		}
	}
}

// SourceMapConsumerV3_StringCharIterator is a implementation of the
// Base64VLQ CharIterator used for decoding the mappings encoded in the JSON
// string.
type SourceMapConsumerV3_StringCharIterator struct {
	// Go-only: language-adapter — the Java String as UTF-16 code units, so
	// that charAt is O(1).
	content []uint16
	length  int
	current int
}

// NewSourceMapConsumerV3_StringCharIterator constructs an iterator over
// content.
//
// Java: SourceMapConsumerV3.StringCharIterator(String)
func NewSourceMapConsumerV3_StringCharIterator(content string) *SourceMapConsumerV3_StringCharIterator {
	it := &SourceMapConsumerV3_StringCharIterator{}
	it.content = jutil.StringToUTF16(content)
	it.length = len(it.content)
	return it
}

// Java: SourceMapConsumerV3.StringCharIterator.next()
func (it *SourceMapConsumerV3_StringCharIterator) Next() uint16 {
	i := it.current
	it.current++
	return it.content[i]
}

// Java: SourceMapConsumerV3.StringCharIterator.peek()
func (it *SourceMapConsumerV3_StringCharIterator) Peek() uint16 {
	return it.content[it.current]
}

// Java: SourceMapConsumerV3.StringCharIterator.hasNext()
func (it *SourceMapConsumerV3_StringCharIterator) HasNext() bool {
	return it.current < it.length
}

// SourceMapConsumerV3_Mappings is a container for mapping entries.
type SourceMapConsumerV3_Mappings struct {
	// Go-only: language-adapter — the enclosing SourceMapConsumerV3.
	outer *SourceMapConsumerV3

	// flatEntries stores all entries sequentially. Each entry is packed into
	// 5 consecutive integers: [generatedCol, sourceLine, sourceCol,
	// sourceFileId, nameId]. For unmapped entries, sourceFileId (and other
	// source-related fields) will be set to UNMAPPED (-1).
	flatEntries []int32
	// lineStart[i] is the index of the first entry for line i in
	// flatEntries. The length of lineStart is lineCount + 1, where
	// lineStart[lineCount] is the total number of entries.
	lineStart []int
	lineCount int
	index     int

	entrySize int
}

// NewSourceMapConsumerV3_Mappings constructs an empty container.
//
// Java: SourceMapConsumerV3.Mappings(int, boolean, int)
func NewSourceMapConsumerV3_Mappings(outer *SourceMapConsumerV3, lineCount int, useCompactMappings bool, estimatedEntries int32) *SourceMapConsumerV3_Mappings {
	m := &SourceMapConsumerV3_Mappings{outer: outer}
	m.lineCount = lineCount
	if useCompactMappings {
		m.entrySize = 4
	} else {
		m.entrySize = 5
	}
	m.flatEntries = make([]int32, estimatedEntries*int32(m.entrySize))
	estimatedLines := 1000
	if lineCount >= 0 {
		estimatedLines = lineCount
	}
	m.lineStart = make([]int, jutil.IntAdd(estimatedLines, 1))
	return m
}

// Java: SourceMapConsumerV3.Mappings.getEntrySize()
func (m *SourceMapConsumerV3_Mappings) GetEntrySize() int {
	return m.entrySize
}

// Java: SourceMapConsumerV3.Mappings.atIndex(int)
func (m *SourceMapConsumerV3_Mappings) AtIndex(index int) *SourceMapConsumerV3_Mappings {
	m.index = index
	return m
}

// Java: SourceMapConsumerV3.Mappings.getGeneratedColumn()
func (m *SourceMapConsumerV3_Mappings) GetGeneratedColumn() int32 {
	return m.flatEntries[m.index]
}

// Java: SourceMapConsumerV3.Mappings.getSourceFileId()
func (m *SourceMapConsumerV3_Mappings) GetSourceFileId() int32 {
	if m.entrySize == 4 {
		val := m.flatEntries[m.index+3] >> 16
		if val == -1 || (val&0xFFFF) == 0xFFFF {
			return SourceMapConsumerV3_UNMAPPED
		}
		return val & 0xFFFF
	}
	return m.flatEntries[m.index+3]
}

// Java: SourceMapConsumerV3.Mappings.getSourceLine()
func (m *SourceMapConsumerV3_Mappings) GetSourceLine() int32 {
	return m.flatEntries[m.index+1]
}

// Java: SourceMapConsumerV3.Mappings.getSourceColumn()
func (m *SourceMapConsumerV3_Mappings) GetSourceColumn() int32 {
	return m.flatEntries[m.index+2]
}

// Java: SourceMapConsumerV3.Mappings.getNameId()
func (m *SourceMapConsumerV3_Mappings) GetNameId() int32 {
	if m.entrySize == 4 {
		val := m.flatEntries[m.index+3] & 0xFFFF
		if val == 0xFFFF {
			return SourceMapConsumerV3_UNMAPPED
		}
		return val
	}
	return m.flatEntries[m.index+4]
}

// Java: SourceMapConsumerV3.Mappings.getParsedLineCount()
func (m *SourceMapConsumerV3_Mappings) GetParsedLineCount() int {
	return m.lineCount
}

// Java: SourceMapConsumerV3.Mappings.getLineStart(int)
func (m *SourceMapConsumerV3_Mappings) GetLineStart(line int) int {
	return m.lineStart[line]
}

// Java: SourceMapConsumerV3.Mappings.setLineStart(int, int)
func (m *SourceMapConsumerV3_Mappings) SetLineStart(line, entriesCount int) {
	if line >= len(m.lineStart) {
		m.lineStart = sourceMapConsumerV3_copyOfInts(m.lineStart, len(m.lineStart)*2)
	}
	m.lineStart[line] = entriesCount
}

// Java: SourceMapConsumerV3.Mappings.setEntry(int, int[], int, int, int, int, int, int)
func (m *SourceMapConsumerV3_Mappings) SetEntry(
	entriesCount int,
	vals []int32,
	entryValues int,
	previousCol int32,
	previousSrcId int32,
	previousSrcLine int32,
	previousSrcColumn int32,
	previousNameId int32) {
	if entriesCount+m.entrySize > len(m.flatEntries) {
		m.flatEntries = sourceMapConsumerV3_copyOfInt32s(m.flatEntries, len(m.flatEntries)*2)
	}
	// `entrySize == 4` indicates `useCompactMappings` is true, meaning the
	// number of unique source files and names are both small enough (under
	// 65535) that their IDs can be bit-packed into a single integer.
	// Otherwise, they each get their own integer.
	if m.entrySize == 4 {
		m.SetEntryCompact(
			entriesCount,
			vals,
			entryValues,
			previousCol,
			previousSrcId,
			previousSrcLine,
			previousSrcColumn,
			previousNameId)
	} else {
		m.SetEntryNormal(
			entriesCount,
			vals,
			entryValues,
			previousCol,
			previousSrcId,
			previousSrcLine,
			previousSrcColumn,
			previousNameId)
	}
}

// Java: SourceMapConsumerV3.Mappings.setEntryCompact(int, int[], int, int, int, int, int, int)
func (m *SourceMapConsumerV3_Mappings) SetEntryCompact(
	entriesCount int,
	vals []int32,
	entryValues int,
	previousCol int32,
	previousSrcId int32,
	previousSrcLine int32,
	previousSrcColumn int32,
	previousNameId int32) {
	// The `vals` array holds the decoded VLQ values in the exact order
	// established by the SourceMap V3 specification: [generatedCol,
	// sourceFileId, sourceLine, sourceCol, nameId]. We reorder these values
	// when storing them in `flatEntries` to minimize branching overhead
	// during lookups, putting `sourceLine` and `sourceColumn` before the IDs.
	switch entryValues {
	case 1:
		m.flatEntries[entriesCount] = vals[0] + previousCol
		m.flatEntries[entriesCount+1] = SourceMapConsumerV3_UNMAPPED
		m.flatEntries[entriesCount+2] = SourceMapConsumerV3_UNMAPPED
		m.flatEntries[entriesCount+3] = SourceMapConsumerV3_UNMAPPED // Both IDs are -1
	case 4:
		m.flatEntries[entriesCount] = vals[0] + previousCol
		m.flatEntries[entriesCount+1] = vals[2] + previousSrcLine
		m.flatEntries[entriesCount+2] = vals[3] + previousSrcColumn
		srcId := vals[1] + previousSrcId
		m.flatEntries[entriesCount+3] = (srcId << 16) | 0xFFFF
	case 5:
		m.flatEntries[entriesCount] = vals[0] + previousCol
		m.flatEntries[entriesCount+1] = vals[2] + previousSrcLine
		m.flatEntries[entriesCount+2] = vals[3] + previousSrcColumn
		srcId := vals[1] + previousSrcId
		nameId := vals[4] + previousNameId
		m.flatEntries[entriesCount+3] = (srcId << 16) | (nameId & 0xFFFF)
	default:
		panic(NewSourceMapParseException(
			"Unexpected number of values for entry:" + jutil.ToJavaString(entryValues)))
	}
}

// Java: SourceMapConsumerV3.Mappings.setEntryNormal(int, int[], int, int, int, int, int, int)
func (m *SourceMapConsumerV3_Mappings) SetEntryNormal(
	entriesCount int,
	vals []int32,
	entryValues int,
	previousCol int32,
	previousSrcId int32,
	previousSrcLine int32,
	previousSrcColumn int32,
	previousNameId int32) {
	// The `vals` array holds the decoded VLQ values in the exact order
	// established by the SourceMap V3 specification: [generatedCol,
	// sourceFileId, sourceLine, sourceCol, nameId]. We reorder these values
	// when storing them in `flatEntries` to minimize branching overhead
	// during lookups, putting `sourceLine` and `sourceColumn` before the IDs.
	switch entryValues {
	case 1:
		m.flatEntries[entriesCount] = vals[0] + previousCol
		m.flatEntries[entriesCount+1] = SourceMapConsumerV3_UNMAPPED
		m.flatEntries[entriesCount+2] = SourceMapConsumerV3_UNMAPPED
		m.flatEntries[entriesCount+3] = SourceMapConsumerV3_UNMAPPED
		m.flatEntries[entriesCount+4] = SourceMapConsumerV3_UNMAPPED
	case 4:
		m.flatEntries[entriesCount] = vals[0] + previousCol
		m.flatEntries[entriesCount+1] = vals[2] + previousSrcLine
		m.flatEntries[entriesCount+2] = vals[3] + previousSrcColumn
		m.flatEntries[entriesCount+3] = vals[1] + previousSrcId
		m.flatEntries[entriesCount+4] = SourceMapConsumerV3_UNMAPPED
	case 5:
		m.flatEntries[entriesCount] = vals[0] + previousCol
		m.flatEntries[entriesCount+1] = vals[2] + previousSrcLine
		m.flatEntries[entriesCount+2] = vals[3] + previousSrcColumn
		m.flatEntries[entriesCount+3] = vals[1] + previousSrcId
		m.flatEntries[entriesCount+4] = vals[4] + previousNameId
	default:
		panic(NewSourceMapParseException(
			"Unexpected number of values for entry:" + jutil.ToJavaString(entryValues)))
	}
}

// Java: SourceMapConsumerV3.Mappings.validateEntry(int, int)
func (m *SourceMapConsumerV3_Mappings) ValidateEntry(entriesCount int, line int) {
	jutil.CheckStateMsg2(
		m.lineCount < 0 || line < m.lineCount, "line=%s, lineCount=%s", line, m.lineCount)
	m.AtIndex(entriesCount)
	sourceFileId := m.GetSourceFileId()
	nameId := m.GetNameId()
	jutil.CheckState(sourceFileId == SourceMapConsumerV3_UNMAPPED || int(sourceFileId) < len(m.outer.sources))
	if nameId != SourceMapConsumerV3_UNMAPPED && m.outer.names == nil {
		// names.length
		panic(jutil.NewNullPointerException())
	}
	jutil.CheckState(nameId == SourceMapConsumerV3_UNMAPPED || int(nameId) < len(m.outer.names))
}

// Java: SourceMapConsumerV3.Mappings.trim(int, int)
func (m *SourceMapConsumerV3_Mappings) Trim(entriesCount int, linesCount int) {
	if len(m.flatEntries) > entriesCount {
		m.flatEntries = sourceMapConsumerV3_copyOfInt32s(m.flatEntries, entriesCount)
	}
	if len(m.lineStart) > linesCount+1 {
		m.lineStart = sourceMapConsumerV3_copyOfInts(m.lineStart, linesCount+1)
	}
	m.lineCount = linesCount
}

// Go-only: language-adapter — Arrays.copyOf(int[], int): a new array of the
// given length, zero-padded.
func sourceMapConsumerV3_copyOfInt32s(a []int32, newLength int) []int32 {
	if newLength <= len(a) {
		return slices.Clone(a[:newLength])
	}
	out := make([]int32, newLength)
	copy(out, a)
	return out
}

func sourceMapConsumerV3_copyOfInts(a []int, newLength int) []int {
	if newLength <= len(a) {
		return slices.Clone(a[:newLength])
	}
	out := make([]int, newLength)
	copy(out, a)
	return out
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
		start := c.mappings.GetLineStart(i)
		end := c.mappings.GetLineStart(i + 1)
		if start != end {
			for j := start; j < end; j += c.mappings.GetEntrySize() {
				c.mappings.AtIndex(j)
				if pending {
					endPosition := NewFilePosition(i, int(c.mappings.GetGeneratedColumn()))
					visitor.Visit(
						sourceName,
						symbolName,
						sourceStartPosition,
						startPosition,
						endPosition)
					pending = false
				}

				sourceFileId := c.mappings.GetSourceFileId()
				if sourceFileId != SourceMapConsumerV3_UNMAPPED {
					pending = true
					sourceName = c.sources[sourceFileId]
					nameId := c.mappings.GetNameId()
					if nameId != SourceMapConsumerV3_UNMAPPED {
						symbolName = c.names[nameId]
					} else {
						symbolName = nil
					}
					sourceStartPosition =
						NewFilePosition(int(c.mappings.GetSourceLine()), int(c.mappings.GetSourceColumn()))
					startPosition = NewFilePosition(i, int(c.mappings.GetGeneratedColumn()))
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
