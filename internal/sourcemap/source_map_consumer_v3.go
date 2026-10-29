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
	"strconv"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// SourceMapConsumerV3_UNMAPPED marks an absent source file or name id.
const SourceMapConsumerV3_UNMAPPED = -1

// SourceMapConsumerV3 is a class for parsing version 3 of the SourceMap
// format, as produced by the Closure Compiler, etc.
// https://github.com/google/closure-compiler/wiki/Source-Maps
type SourceMapConsumerV3 struct {
	sources        []optionalString
	sourcesContent []optionalString
	names          []optionalString
	file           optionalString
	lineCount      int

	mappings *SourceMapConsumerV3_Mappings

	// (originalFile path, original line) ==> target mappings
	reverseSourceMapping map[sourceMapConsumerV3_lineKey][]*Mapping_OriginalMapping
	sourceRoot           optionalString
	extensions           []SourceMapExtension
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

func sourceMapConsumerV3_keyOf(s optionalString, line int) sourceMapConsumerV3_lineKey {
	if !s.present {
		return sourceMapConsumerV3_lineKey{line: line, isNull: true}
	}
	return sourceMapConsumerV3_lineKey{name: s.text, line: line}
}

// NewSourceMapConsumerV3 constructs an empty consumer.
//
// Java: SourceMapConsumerV3()
func NewSourceMapConsumerV3() *SourceMapConsumerV3 {
	return &SourceMapConsumerV3{}
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
	if err := c.tryParseThrowable(contents); err != nil {
		panic(err)
	}
}

// TryParse parses the given contents containing a source map, returning a
// *SourceMapParseException instead of panicking on malformed source maps.
func (c *SourceMapConsumerV3) TryParse(contents string) *SourceMapParseException {
	err := c.tryParseThrowable(contents)
	if err == nil {
		return nil
	}
	if sm, ok := err.(*SourceMapParseException); ok {
		return sm
	}
	return NewSourceMapParseException(err.GetMessage())
}

func (c *SourceMapConsumerV3) tryParseThrowable(contents string) jutil.Throwable {
	sourceMapObject, err := parseSourceMapObject(contents)
	if err != nil {
		return err
	}
	return c.tryParseSourceMapObject(sourceMapObject, nil)
}

// ParseSourceMapObject parses the given contents containing a source map.
// It panics with a *SourceMapParseException.
//
// Java: SourceMapConsumerV3.parse(SourceMapObject, SourceMapSupplier)
func (c *SourceMapConsumerV3) ParseSourceMapObject(sourceMapObject *SourceMapObject, sectionSupplier SourceMapSupplier) {
	if err := c.tryParseSourceMapObject(sourceMapObject, sectionSupplier); err != nil {
		panic(err)
	}
}

func (c *SourceMapConsumerV3) tryParseSourceMapObject(sourceMapObject *SourceMapObject, sectionSupplier SourceMapSupplier) jutil.Throwable {
	if sourceMapObject == nil {
		return jutil.NewNullPointerException()
	}
	if sourceMapObject.version != 3 {
		return NewSourceMapParseException("Unknown version: " + strconv.Itoa(sourceMapObject.version))
	}

	c.file = sourceMapObject.file
	if c.file.present && c.file.text == "" {
		return NewSourceMapParseException("File entry is empty")
	}

	if sourceMapObject.sections != nil {
		// Looks like a index map, try to parse it that way.
		return c.tryParseMetaMap(sourceMapObject, sectionSupplier)
	}

	c.lineCount = sourceMapObject.lineCount
	c.sourceRoot = sourceMapObject.sourceRoot
	c.sources = sourceMapObject.sources
	c.sourcesContent = sourceMapObject.sourcesContent
	c.names = sourceMapObject.names

	if c.sources == nil {
		// sources.length
		return jutil.NewNullPointerException()
	}
	useCompactMappings :=
		(len(c.sources) < 65535) && (c.names == nil || len(c.names) < 65535)

	var estimatedEntries int32
	if c.lineCount >= 0 {
		estimatedEntries = int32(c.lineCount) * 2
	} else {
		estimatedEntries = 1000
	}
	if sourceMapObject.mappings.present {
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
		estimatedEntries = max(estimatedEntries, int32(jutil.StringLength(sourceMapObject.mappings.text)/6))
	}

	// The value type of each extension is the native JSON type (e.g.
	// JsonObject, or JSONObject when compiled with GWT).
	for _, ext := range sourceMapObject.extensions {
		c.extensions = putExtension(c.extensions, ext.Key, ext.Value)
	}
	mb, err := tryNewMappingBuilder(
		c, sourceMapObject.mappings, c.lineCount, useCompactMappings, estimatedEntries)
	if err != nil {
		return err
	}
	mappings, err := mb.tryBuild()
	if err != nil {
		return err
	}
	c.mappings = mappings
	return nil
}

// ParseMetaMap parses an index map by merging its sections directly into this consumer.
//
// Java: SourceMapConsumerV3.parseMetaMap(SourceMapObject, SourceMapSupplier)
func (c *SourceMapConsumerV3) ParseMetaMap(sourceMapObject *SourceMapObject, sectionSupplier SourceMapSupplier) {
	if err := c.tryParseMetaMap(sourceMapObject, sectionSupplier); err != nil {
		panic(err)
	}
}

type metaSegment struct {
	sourceIdx int32
	nameIdx   int32
	origLine  int32
	origCol   int32
	startLine int32
	startCol  int32
	endLine   int32
	endCol    int32
}

func (c *SourceMapConsumerV3) tryParseMetaMap(sourceMapObject *SourceMapObject, sectionSupplier SourceMapSupplier) jutil.Throwable {
	if sourceMapObject == nil {
		return jutil.NewNullPointerException()
	}
	if sectionSupplier == nil {
		sectionSupplier = &SourceMapConsumerV3_DefaultSourceMapSupplier{}
	}

	if sourceMapObject.lineCount >= 0 ||
		sourceMapObject.mappings.present ||
		len(sourceMapObject.sources) > 0 ||
		sourceMapObject.names != nil {
		return NewSourceMapParseException("Invalid map format")
	}

	var segments []metaSegment
	var addedSourcesBuf [16]string
	addedSources := addedSourcesBuf[:0]
	var addedSourceIdx map[string]int32
	var lastAddedSource string
	lastAddedSourceIdx := int32(-1)

	var addedNamesBuf [16]string
	addedNames := addedNamesBuf[:0]
	var addedNameIdx map[string]int32
	var lastAddedName string
	lastAddedNameIdx := int32(-1)

	var lastStartLine, lastStartCol int32

	addSegment := func(
		offLine, offCol int32,
		srcName optionalString, symName optionalString,
		origLine, origCol, startLine, startCol, endLine, endCol int32,
	) jutil.Throwable {
		if !srcName.present || origLine < 0 {
			return nil
		}
		if offLine != 0 || offCol != 0 {
			startOffCol := offCol
			endOffCol := offCol
			if startLine > 0 {
				startOffCol = 0
			}
			if endLine > 0 {
				endOffCol = 0
			}
			startLine += offLine
			startCol += startOffCol
			endLine += offLine
			endCol += endOffCol
		}
		if len(segments) > 0 {
			if !(startLine > lastStartLine || (startLine == lastStartLine && startCol >= lastStartCol)) {
				return jutil.NewIllegalStateExceptionMsg(jutil.LenientFormat(
					"Incorrect source mappings order, previous : (%s,%s)\nnew : (%s,%s)",
					int(lastStartLine), int(lastStartCol), int(startLine), int(startCol)))
			}
		}
		lastStartLine = startLine
		lastStartCol = startCol

		sourceIdx := internSourceMapString(&addedSources, &addedSourceIdx, &lastAddedSource, &lastAddedSourceIdx, srcName.text)
		nameIdx := int32(-1)
		if symName.present {
			nameIdx = internSourceMapString(&addedNames, &addedNameIdx, &lastAddedName, &lastAddedNameIdx, symName.text)
		}
		segments = append(segments, metaSegment{
			sourceIdx: sourceIdx,
			nameIdx:   nameIdx,
			origLine:  origLine,
			origCol:   origCol,
			startLine: startLine,
			startCol:  startCol,
			endLine:   endLine,
			endCol:    endCol,
		})
		return nil
	}

	for _, section := range sourceMapObject.sections {
		if section == nil {
			return jutil.NewNullPointerException()
		}
		var mapText string
		var hasMap bool
		if section.type_ == SourceMapSection_SectionType_URL {
			if !section.value.present {
				return jutil.NewNullPointerException()
			}
			supplied, supplierErr := callSupplierGetSourceMap(sectionSupplier, section.value.text)
			if supplierErr != nil {
				return supplierErr
			}
			if supplied != nil {
				mapText = *supplied
				hasMap = true
			}
		} else if section.value.present {
			mapText = section.value.text
			hasMap = true
		}
		if !hasMap {
			valStr := "null"
			if section.value.present {
				valStr = section.value.text
			}
			return NewSourceMapParseException("Unable to retrieve: " + valStr)
		}
		if section.line < 0 || section.column < 0 {
			return jutil.NewIllegalStateException()
		}
		sectionConsumer := NewSourceMapConsumerV3()
		if err := sectionConsumer.tryParseThrowable(mapText); err != nil {
			return err
		}

		offLine := int32(section.line)
		offCol := int32(section.column)
		pending := false
		var pendingSrc, pendingSym optionalString
		var pendingOrigLine, pendingOrigCol, pendingStartLine, pendingStartCol int32
		lineCount := sectionConsumer.mappings.GetParsedLineCount()
		for i := 0; i < lineCount; i++ {
			if sectionConsumer.mappings.lineIsEmpty(i) {
				continue
			}
			entries := sectionConsumer.mappings.line(i)
			for j := range entries {
				entry := &entries[j]
				if pending {
					if err := addSegment(
						offLine, offCol, pendingSrc, pendingSym,
						pendingOrigLine, pendingOrigCol,
						pendingStartLine, pendingStartCol, int32(i), entry.genCol,
					); err != nil {
						return err
					}
					pending = false
				}
				if entry.srcId != SourceMapConsumerV3_UNMAPPED {
					pending = true
					pendingSrc = sectionConsumer.sources[entry.srcId]
					if entry.nameId != SourceMapConsumerV3_UNMAPPED {
						pendingSym = sectionConsumer.names[entry.nameId]
					} else {
						pendingSym = optionalString{}
					}
					pendingOrigLine = entry.srcLine
					pendingOrigCol = entry.srcCol
					pendingStartLine = int32(i)
					pendingStartCol = entry.genCol
				}
			}
		}
		if pending {
			if err := addSegment(
				offLine, offCol, pendingSrc, pendingSym,
				pendingOrigLine, pendingOrigCol,
				pendingStartLine, pendingStartCol, pendingStartLine, pendingStartCol+1,
			); err != nil {
				return err
			}
		}
	}

	maxEndLine := 0
	if err := traverseMetaSegments(segments, func(m *metaSegment, line, col, nextLine, nextCol int) jutil.Throwable {
		if m != nil && int(m.endLine) > maxEndLine {
			maxEndLine = int(m.endLine)
		}
		return nil
	}); err != nil {
		return err
	}
	maxLine := jutil.IntAdd(maxEndLine, 1)

	sourceIdxToId := make([]int32, len(addedSources))
	for i := range sourceIdxToId {
		sourceIdxToId[i] = -1
	}
	orderedSources := make([]optionalString, 0, len(addedSources))

	nameIdxToId := make([]int32, len(addedNames))
	for i := range nameIdxToId {
		nameIdxToId[i] = -1
	}
	orderedNames := make([]optionalString, 0, len(addedNames))

	data := make([]byte, 0, len(segments)*4)
	estLines := 1
	if maxLine > 0 {
		estLines = maxLine
	}
	lineStart := make([]uint32, 1, jutil.IntAdd(estLines, 1))
	linesCount := 0
	var prevInLine sourceMapConsumerV3_entry

	closeLine := func() {
		linesCount++
		lineStart = append(lineStart, uint32(len(data)))
		prevInLine = sourceMapConsumerV3_entry{}
	}

	if err := traverseMetaSegments(segments, func(m *metaSegment, line, col, nextLine, nextCol int) jutil.Throwable {
		if line != nextLine || col != nextCol {
			if line < maxLine {
				e := sourceMapConsumerV3_entry{
					genCol:  int32(col),
					srcLine: SourceMapConsumerV3_UNMAPPED,
					srcCol:  SourceMapConsumerV3_UNMAPPED,
					srcId:   SourceMapConsumerV3_UNMAPPED,
					nameId:  SourceMapConsumerV3_UNMAPPED,
				}
				entryValues := 1
				if m != nil {
					srcId := sourceIdxToId[m.sourceIdx]
					if srcId < 0 {
						srcId = int32(len(orderedSources))
						orderedSources = append(orderedSources, optionalString{text: addedSources[m.sourceIdx], present: true})
						sourceIdxToId[m.sourceIdx] = srcId
					}
					e.srcId = srcId
					e.srcLine = m.origLine
					e.srcCol = m.origCol
					entryValues = 4
					if m.nameIdx >= 0 {
						nameId := nameIdxToId[m.nameIdx]
						if nameId < 0 {
							nameId = int32(len(orderedNames))
							orderedNames = append(orderedNames, optionalString{text: addedNames[m.nameIdx], present: true})
							nameIdxToId[m.nameIdx] = nameId
						}
						e.nameId = nameId
						entryValues = 5
					}
				}
				data = encodeMappingEntry(data, &prevInLine, entryValues, e)
			} else if m != nil {
				return jutil.NewIllegalStateException()
			}
		}
		for i := line; i <= nextLine && i < maxLine; i++ {
			if i == nextLine {
				break
			}
			closeLine()
		}
		return nil
	}); err != nil {
		return err
	}
	closeLine()

	c.file = sourceMapObject.file
	if c.file.present && c.file.text == "" {
		return NewSourceMapParseException("File entry is empty")
	}
	c.lineCount = maxLine
	c.sourceRoot = optionalString{}
	c.sources = orderedSources
	c.sourcesContent = nil
	c.names = orderedNames
	entrySize := 5
	if len(orderedSources) < 65535 && len(orderedNames) < 65535 {
		entrySize = 4
	}
	c.mappings = newSourceMapConsumerV3_Mappings(data, lineStart, linesCount, entrySize)
	return nil
}

func traverseMetaSegments(segments []metaSegment, visit func(m *metaSegment, line, col, nextLine, nextCol int) jutil.Throwable) jutil.Throwable {
	var stackBuf [16]metaSegment
	stack := stackBuf[:0]
	tLine, tCol := 0, 0

	emit := func(m *metaSegment, nextLine, nextCol int) jutil.Throwable {
		if !(tLine <= nextLine && (tLine < nextLine || tCol < nextCol)) {
			return jutil.NewIllegalStateException()
		}
		if err := visit(m, tLine, tCol, nextLine, nextCol); err != nil {
			return err
		}
		tLine = nextLine
		tCol = nextCol
		return nil
	}

	maybeVisit := func(m *metaSegment) jutil.Throwable {
		nextLine := int(m.endLine)
		nextCol := int(m.endCol)
		if tLine < nextLine || (tLine == nextLine && tCol < nextCol) {
			return emit(m, nextLine, nextCol)
		}
		return nil
	}

	for i := range segments {
		m := &segments[i]
		for len(stack) > 0 {
			top := &stack[len(stack)-1]
			if (top.endLine == m.startLine && top.endCol >= m.startCol) || top.endLine > m.startLine {
				break
			}
			closing := *top
			stack = stack[:len(stack)-1]
			if err := maybeVisit(&closing); err != nil {
				return err
			}
		}
		var parent *metaSegment
		if len(stack) > 0 {
			parent = &stack[len(stack)-1]
		}
		nextLine := int(m.startLine)
		nextCol := int(m.startCol)
		if !(tLine < nextLine || tCol <= nextCol) {
			return jutil.NewIllegalStateException()
		}
		if tLine < nextLine || (tLine == nextLine && tCol < nextCol) {
			if err := emit(parent, nextLine, nextCol); err != nil {
				return err
			}
		}
		stack = append(stack, *m)
	}
	for len(stack) > 0 {
		closing := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if err := maybeVisit(&closing); err != nil {
			return err
		}
	}
	return nil
}

func callSupplierGetSourceMap(supplier SourceMapSupplier, url string) (supplied *string, parseErr jutil.Throwable) {
	defer func() {
		if r := recover(); r != nil {
			if ex, ok := r.(jutil.IIOException); ok {
				parseErr = NewSourceMapParseException("IO exception: " + ex.Error())
				return
			}
			if t, ok := r.(jutil.Throwable); ok {
				parseErr = t
				return
			}
			panic(r)
		}
	}()
	return supplier.GetSourceMap(url), nil
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
	return optionalStringsToPtrs(c.sources)
}

// GetOriginalSourcesContent returns the sources' contents, or nil.
//
// Java: SourceMapConsumerV3.getOriginalSourcesContent()
func (c *SourceMapConsumerV3) GetOriginalSourcesContent() []*string {
	if c.sourcesContent == nil {
		return nil
	}
	return optionalStringsToPtrs(c.sourcesContent)
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
	return optionalStringsToPtrs(c.names)
}

// Java: SourceMapConsumerV3.getFile()
func (c *SourceMapConsumerV3) GetFile() *string {
	return c.file.ptr()
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
	return c.sourceRoot.ptr()
}

// GetExtensions returns all extensions and their values (which can be any
// json value) in insertion order.
//
// Java: SourceMapConsumerV3.getExtensions()
func (c *SourceMapConsumerV3) GetExtensions() []SourceMapExtension {
	return c.extensions
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
	origSrc := &c.sources[sourceFileId]
	if !origSrc.present {
		panic(jutil.NewNullPointerException())
	}
	origFile := origSrc.text
	lineNum := int(entry.srcLine + 1)
	colPos := int(entry.srcCol + 1)
	bitField := int32(mapping_originalFileBit | mapping_lineNumberBit | mapping_columnPositionBit | mapping_precisionBit)
	var identifier string
	var identifierPtr *string
	nameId := entry.nameId
	if nameId != SourceMapConsumerV3_UNMAPPED {
		origName := &c.names[nameId]
		if !origName.present {
			panic(jutil.NewNullPointerException())
		}
		identifier = origName.text
		bitField |= mapping_identifierBit
		if identifier != "" {
			identifierPtr = &origName.text
		}
	}
	return Mapping_OriginalMapping{
		bitField0_:      bitField,
		precision_:      precision,
		originalFile_:   origFile,
		lineNumber_:     lineNum,
		columnPosition_: colPos,
		identifier_:     identifier,
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
				mapping := &Mapping_OriginalMapping{
					bitField0_:      mapping_lineNumberBit | mapping_columnPositionBit,
					lineNumber_:     targetLine,
					columnPosition_: int(entry.genCol),
				}
				c.reverseSourceMapping[key] = append(c.reverseSourceMapping[key], mapping)
			}
		}
	}
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
	if gv, ok := visitor.(*SourceMapGeneratorV3_ConsumerEntryVisitor); ok {
		c.visitMappingsToGenerator(gv.outer)
		return
	}

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
					sourceName = c.sources[sourceFileId].ptr()
					nameId := entry.nameId
					if nameId != SourceMapConsumerV3_UNMAPPED {
						symbolName = c.names[nameId].ptr()
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

func (c *SourceMapConsumerV3) visitMappingsToGenerator(g *SourceMapGeneratorV3) {
	pending := false
	var sourceName optionalString
	var symbolName *string
	var sourceStartPosition FilePosition
	var startPosition FilePosition

	for i := 0; i < c.mappings.GetParsedLineCount(); i++ {
		if !c.mappings.lineIsEmpty(i) {
			entries := c.mappings.line(i)
			for j := range entries {
				entry := &entries[j]
				if pending {
					if sourceName.present {
						endPosition := MakeFilePosition(i, int(entry.genCol))
						g.AddMappingPositions(sourceName.text, symbolName, sourceStartPosition, startPosition, endPosition)
					}
					pending = false
				}

				sourceFileId := entry.srcId
				if sourceFileId != SourceMapConsumerV3_UNMAPPED {
					pending = true
					sourceName = c.sources[sourceFileId]
					nameId := entry.nameId
					if nameId != SourceMapConsumerV3_UNMAPPED {
						symbolName = c.names[nameId].ptr()
					} else {
						symbolName = nil
					}
					sourceStartPosition = MakeFilePosition(int(entry.srcLine), int(entry.srcCol))
					startPosition = MakeFilePosition(i, int(entry.genCol))
				}
			}
		}
	}
	if pending && sourceName.present {
		endPosition := MakeFilePosition(startPosition.GetLine(), jutil.IntAdd(startPosition.GetColumn(), 1))
		g.AddMappingPositions(sourceName.text, symbolName, sourceStartPosition, startPosition, endPosition)
	}
}
