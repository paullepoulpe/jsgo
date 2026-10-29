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
	"strings"

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
	if err := c.TryParse(contents); err != nil {
		panic(err)
	}
}

// TryParse parses the given contents containing a source map, returning a
// *SourceMapParseException instead of panicking on malformed source maps.
func (c *SourceMapConsumerV3) TryParse(contents string) *SourceMapParseException {
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

func (c *SourceMapConsumerV3) tryParseSourceMapObject(sourceMapObject *SourceMapObject, sectionSupplier SourceMapSupplier) *SourceMapParseException {
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
	mappings, err := newMappingBuilder(
		c, sourceMapObject.mappings, c.lineCount, useCompactMappings, estimatedEntries).
		tryBuild()
	if err != nil {
		return err
	}
	c.mappings = mappings
	return nil
}

// ParseMetaMap parses an index map: its sections are merged into a new
// generator, whose map is then parsed by this consumer.
//
// Java: SourceMapConsumerV3.parseMetaMap(SourceMapObject, SourceMapSupplier)
func (c *SourceMapConsumerV3) ParseMetaMap(sourceMapObject *SourceMapObject, sectionSupplier SourceMapSupplier) {
	if err := c.tryParseMetaMap(sourceMapObject, sectionSupplier); err != nil {
		panic(err)
	}
}

func (c *SourceMapConsumerV3) tryParseMetaMap(sourceMapObject *SourceMapObject, sectionSupplier SourceMapSupplier) *SourceMapParseException {
	if sectionSupplier == nil {
		sectionSupplier = &SourceMapConsumerV3_DefaultSourceMapSupplier{}
	}

	if sourceMapObject.lineCount >= 0 ||
		sourceMapObject.mappings.present ||
		len(sourceMapObject.sources) > 0 ||
		sourceMapObject.names != nil {
		return NewSourceMapParseException("Invalid map format")
	}

	// Build up a new source map in a new generator using the mappings of
	// this metamap. The new map will be rendered to JSON and then parsed
	// using this consumer.
	generator := NewSourceMapGeneratorV3()
	for _, section := range sourceMapObject.sections {
		var mapText string
		var hasMap bool
		if section.type_ == SourceMapSection_SectionType_URL {
			if !section.value.present {
				panic(jutil.NewNullPointerException())
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
		generator.SetStartingPosition(section.line, section.column)
		sectionConsumer := NewSourceMapConsumerV3()
		if err := sectionConsumer.TryParse(mapText); err != nil {
			return err
		}
		sectionConsumer.VisitMappings(&SourceMapGeneratorV3_ConsumerEntryVisitor{outer: generator})
	}

	var sb strings.Builder
	generator.AppendTo(&sb, sourceMapObject.file.ptr())
	return c.TryParse(sb.String())
}

func callSupplierGetSourceMap(supplier SourceMapSupplier, url string) (supplied *string, parseErr *SourceMapParseException) {
	defer func() {
		if r := recover(); r != nil {
			if ex, ok := r.(jutil.IIOException); ok {
				parseErr = NewSourceMapParseException("IO exception: " + ex.Error())
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
	prec := int32(precision.GetNumber())
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
