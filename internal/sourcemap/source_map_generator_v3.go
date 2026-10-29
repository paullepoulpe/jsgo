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
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// SourceMapGeneratorV3 collects information mapping the generated
// (compiled) source back to its original source for debugging purposes.
//
// Source Map Revision 3 Proposal:
// https://docs.google.com/document/d/1U1RGAehQwRypUTovF1KRlpiOFze0b-_2gc6fAH0KY0k/edit?usp=sharing
type SourceMapGeneratorV3 struct {
	// Go-only: goroutines AppendTo may use (source_map_generator_v3_parallel.go).
	parallelWorkers int

	// A map of source names to source file contents
	sourceFileContentMap map[string]optionalString

	// The mappings, in the order they were added.
	//
	// Go-only: performance — stored delta-encoded (see
	// sourceMapGeneratorV3_mappingStore) rather than as Mapping objects: a
	// big build keeps millions of them alive through the whole output.
	mappings sourceMapGeneratorV3_mappingStore

	// The position that the current source map is offset in the buffer being
	// used to generated the compiled source file.
	offsetPosition FilePosition

	// The position that the current source map is offset in the generated
	// the compiled source file by the addition of a an output wrapper prefix.
	prefixPosition FilePosition

	// A list of extensions to be added to sourcemap. The value is a object to
	// permit single values, like strings or numbers, and JsonObject or
	// JsonArray objects.
	extensions []SourceMapExtension

	// The source root path for relocating source fails or avoid duplicate
	// values on the source entry.
	sourceRootPath string

	addedSourceFiles       []string
	addedSourceFileIndex   map[string]int32
	sourceIdxToFileId      []int32
	orderedSourceFiles     []string
	lastAddedSourceFile    string
	lastAddedSourceFileIdx int32
	addedSymbolNames       []string
	addedSymbolIndex       map[string]int32
	symbolIdxToNameId      []int32
	orderedSymbolNames     []string
	lastAddedSymbolName    string
	lastAddedSymbolNameIdx int32
	addedSourceFilesBuf    [16]string
	addedSymbolNamesBuf    [16]string
}

// SourceMapGeneratorV3_ExtensionMergeAction provides the merging strategy
// when an extension conflict appears because of merging two source maps on
// method MergeMapSectionWithMergeAction.
type SourceMapGeneratorV3_ExtensionMergeAction interface {
	// Merge returns the merged value between two extensions with the same
	// name when merging two source maps.
	//
	// extensionKey is the extension name in conflict; currentValue is the
	// extension value in the current source map; newValue is the extension
	// value in the input source map.
	//
	// Java: SourceMapGeneratorV3.ExtensionMergeAction.merge(String, Object, Object)
	Merge(extensionKey string, currentValue any, newValue any) any
}

// SourceMapGeneratorV3_UNMAPPED is the id of a mapping not (yet) used.
const SourceMapGeneratorV3_UNMAPPED = -1

// NewSourceMapGeneratorV3 constructs an empty generator.
//
// Java: SourceMapGeneratorV3()
func NewSourceMapGeneratorV3() *SourceMapGeneratorV3 {
	g := &SourceMapGeneratorV3{}
	g.addedSourceFiles = g.addedSourceFilesBuf[:0:len(g.addedSourceFilesBuf)]
	g.lastAddedSourceFileIdx = -1
	g.addedSymbolNames = g.addedSymbolNamesBuf[:0:len(g.addedSymbolNamesBuf)]
	g.lastAddedSymbolNameIdx = -1
	g.offsetPosition = MakeFilePosition(0, 0)
	g.prefixPosition = MakeFilePosition(0, 0)
	return g
}

// Java: SourceMapGeneratorV3.reset()
func (g *SourceMapGeneratorV3) Reset() {
	// Do not reset sourceFileContentMap
	g.mappings.reset()
	clear(g.addedSourceFiles)
	g.addedSourceFiles = g.addedSourceFiles[:0]
	clear(g.addedSourceFileIndex)
	g.sourceIdxToFileId = g.sourceIdxToFileId[:0]
	clear(g.orderedSourceFiles)
	g.orderedSourceFiles = g.orderedSourceFiles[:0]
	g.lastAddedSourceFile = ""
	g.lastAddedSourceFileIdx = -1
	clear(g.addedSymbolNames)
	g.addedSymbolNames = g.addedSymbolNames[:0]
	clear(g.addedSymbolIndex)
	g.symbolIdxToNameId = g.symbolIdxToNameId[:0]
	clear(g.orderedSymbolNames)
	g.orderedSymbolNames = g.orderedSymbolNames[:0]
	g.lastAddedSymbolName = ""
	g.lastAddedSymbolNameIdx = -1
	g.offsetPosition = MakeFilePosition(0, 0)
	g.prefixPosition = MakeFilePosition(0, 0)
}

// Validate sets whether to perform (potentially costly) validation on the
// generated source map.
//
// Java: SourceMapGeneratorV3.validate(boolean)
func (g *SourceMapGeneratorV3) Validate(validate bool) {
	// Nothing currently.
}

// SetWrapperPrefix sets the prefix used for wrapping the generated source
// file before it is written. This ensures that the source map is adjusted
// for the change in character offsets.
//
// Java: SourceMapGeneratorV3.setWrapperPrefix(String)
func (g *SourceMapGeneratorV3) SetWrapperPrefix(prefix string) {
	// Determine the current line and character position (in UTF-16 code units).
	prefixLine := 0
	var prefixIndex int
	if lastNL := strings.LastIndexByte(prefix, '\n'); lastNL >= 0 {
		prefixLine = strings.Count(prefix[:lastNL+1], "\n")
		prefixIndex = jutil.StringLength(prefix[lastNL+1:])
	} else {
		prefixIndex = jutil.StringLength(prefix)
	}
	g.prefixPosition = MakeFilePosition(prefixLine, prefixIndex)
}

// SetStartingPosition sets the source code that exists in the buffer for
// which the generated code is being generated. This ensures that the source
// map accurately reflects the fact that the source is being appended to an
// existing buffer and as such, does not start at line 0, position 0 but
// rather some other line and position.
//
// Java: SourceMapGeneratorV3.setStartingPosition(int, int)
func (g *SourceMapGeneratorV3) SetStartingPosition(offsetLine, offsetIndex int) {
	jutil.CheckState(offsetLine >= 0)
	jutil.CheckState(offsetIndex >= 0)
	g.offsetPosition = MakeFilePosition(offsetLine, offsetIndex)
}

// AddMapping adds a mapping for the given node. Mappings must be added in
// order.
//
// Java: SourceMapGeneratorV3.addMapping(String, String, FilePosition, FilePosition, FilePosition)
func (g *SourceMapGeneratorV3) AddMapping(
	sourceName *string, symbolName *string,
	sourceStartPosition *FilePosition,
	startPosition *FilePosition, endPosition *FilePosition) {
	if sourceName == nil || sourceStartPosition.GetLine() < 0 {
		return
	}
	g.AddMappingPositions(*sourceName, symbolName, *sourceStartPosition, *startPosition, *endPosition)
}

// AddMappingPositions adds a mapping using value FilePositions without heap allocation.
func (g *SourceMapGeneratorV3) AddMappingPositions(
	sourceName string, symbolName *string,
	sourceStartPosition FilePosition,
	startPosition FilePosition, endPosition FilePosition) {

	// Don't bother if there is not sufficient information to be useful.
	if sourceStartPosition.GetLine() < 0 {
		return
	}

	adjustedStart := startPosition
	adjustedEnd := endPosition

	if g.offsetPosition.GetLine() != 0 ||
		g.offsetPosition.GetColumn() != 0 {
		// If the mapping is found on the first line, we need to offset its
		// character position by the number of characters found on the *last*
		// line of the source file to which the code is being generated.
		offsetLine := g.offsetPosition.GetLine()
		startOffsetPosition := g.offsetPosition.GetColumn()
		endOffsetPosition := g.offsetPosition.GetColumn()

		if startPosition.GetLine() > 0 {
			startOffsetPosition = 0
		}

		if endPosition.GetLine() > 0 {
			endOffsetPosition = 0
		}

		adjustedStart = MakeFilePosition(
			jutil.IntAdd(startPosition.GetLine(), offsetLine),
			jutil.IntAdd(startPosition.GetColumn(), startOffsetPosition))

		adjustedEnd = MakeFilePosition(
			jutil.IntAdd(endPosition.GetLine(), offsetLine),
			jutil.IntAdd(endPosition.GetColumn(), endOffsetPosition))
	}

	// Validate the mappings are in a proper order.
	if g.mappings.n > 0 {
		lastLine := g.mappings.last.startLine
		lastColumn := g.mappings.last.startCol
		nextLine := int32(adjustedStart.GetLine())
		nextColumn := int32(adjustedStart.GetColumn())
		// Go-only: language-adapter — the message arguments are boxed only on
		// failure (PORTING.md §5); jutil.CheckStateMsg throws exactly as
		// Preconditions.checkState(boolean, String, Object...).
		if !(nextLine > lastLine || (nextLine == lastLine && nextColumn >= lastColumn)) {
			jutil.CheckStateMsg(false,
				"Incorrect source mappings order, previous : (%s,%s)\n"+
					"new : (%s,%s)",
				int(lastLine), int(lastColumn), int(nextLine), int(nextColumn))
		}
	}

	sourceFileIdx := internSourceMapString(
		&g.addedSourceFiles, &g.addedSourceFileIndex,
		&g.lastAddedSourceFile, &g.lastAddedSourceFileIdx,
		sourceName)

	originalNameIdx := int32(-1)
	if symbolName != nil {
		originalNameIdx = internSourceMapString(
			&g.addedSymbolNames, &g.addedSymbolIndex,
			&g.lastAddedSymbolName, &g.lastAddedSymbolNameIdx,
			*symbolName)
	}

	// Create the new mapping.
	mapping := SourceMapGeneratorV3_Mapping{
		sourceFileIdx:   sourceFileIdx,
		originalLine:    int32(sourceStartPosition.GetLine()),
		originalCol:     int32(sourceStartPosition.GetColumn()),
		originalNameIdx: originalNameIdx,
		startLine:       int32(adjustedStart.GetLine()),
		startCol:        int32(adjustedStart.GetColumn()),
		endLine:         int32(adjustedEnd.GetLine()),
		endCol:          int32(adjustedEnd.GetColumn()),
	}
	g.mappings.add(&mapping)
}

func internSourceMapString(
	items *[]string,
	index *map[string]int32,
	last *string,
	lastIdx *int32,
	s string,
) int32 {
	if *lastIdx >= 0 && *last == s {
		return *lastIdx
	}
	var idx int32
	if len(*items) < 16 {
		found := false
		for i, item := range *items {
			if item == s {
				idx = int32(i)
				found = true
				break
			}
		}
		if !found {
			idx = int32(len(*items))
			*items = append(*items, s)
		}
	} else {
		if len(*index) == 0 {
			if *index == nil {
				*index = make(map[string]int32, 256)
			}
			for i, item := range *items {
				(*index)[item] = int32(i)
			}
		}
		if existing, ok := (*index)[s]; ok {
			idx = existing
		} else {
			idx = int32(len(*items))
			*items = append(*items, s)
			(*index)[s] = idx
		}
	}
	*last = s
	*lastIdx = idx
	return idx
}

// Java: SourceMapGeneratorV3.addSourcesContent(String, String)
func (g *SourceMapGeneratorV3) AddSourcesContent(source string, content *string) {
	if g.sourceFileContentMap == nil {
		g.sourceFileContentMap = make(map[string]optionalString)
	}
	g.sourceFileContentMap[source] = makeOptionalString(content)
}

// ShareConfigFrom shares read-only configuration (sourcesContent map,
// sourceRootPath, and extensions) from a template generator so parallel
// per-chunk generators produce identical output without re-populating
// sourcesContent per chunk.
func (g *SourceMapGeneratorV3) ShareConfigFrom(other *SourceMapGeneratorV3) {
	if other == nil {
		return
	}
	g.sourceFileContentMap = other.sourceFileContentMap
	g.sourceRootPath = other.sourceRootPath
	if len(other.extensions) > 0 {
		g.extensions = slices.Clone(other.extensions)
	}
}

// SourceMapGeneratorV3_ConsumerEntryVisitor adds the mappings a consumer
// visits to the enclosing generator.
type SourceMapGeneratorV3_ConsumerEntryVisitor struct {
	// Go-only: language-adapter — the enclosing SourceMapGeneratorV3.
	outer *SourceMapGeneratorV3
}

// Java: SourceMapGeneratorV3.ConsumerEntryVisitor.visit(String, String, FilePosition, FilePosition, FilePosition)
func (v *SourceMapGeneratorV3_ConsumerEntryVisitor) Visit(
	sourceName *string, symbolName *string,
	sourceStartPosition *FilePosition,
	startPosition *FilePosition, endPosition *FilePosition) {
	v.outer.AddMapping(sourceName, symbolName,
		sourceStartPosition, startPosition, endPosition)
}

// MergeMapSection merges current mapping with mapSectionContents
// considering the offset (line, column). Any extension in the map section
// will be ignored. It panics with a *SourceMapParseException.
//
// Java: SourceMapGeneratorV3.mergeMapSection(int, int, String)
func (g *SourceMapGeneratorV3) MergeMapSection(line, column int, mapSectionContents string) {
	g.SetStartingPosition(line, column)
	section := NewSourceMapConsumerV3()
	section.Parse(mapSectionContents)
	section.VisitMappings(&SourceMapGeneratorV3_ConsumerEntryVisitor{outer: g})
}

// MergeMapSectionWithMergeAction works like MergeMapSection, except that
// extensions from the mapSectionContents are merged to the top level source
// map. For conflicts a mergeAction is performed.
//
// Java: SourceMapGeneratorV3.mergeMapSection(int, int, String, ExtensionMergeAction)
func (g *SourceMapGeneratorV3) MergeMapSectionWithMergeAction(
	line, column int, mapSectionContents string, mergeAction SourceMapGeneratorV3_ExtensionMergeAction) {
	g.SetStartingPosition(line, column)
	section := NewSourceMapConsumerV3()
	section.Parse(mapSectionContents)
	section.VisitMappings(&SourceMapGeneratorV3_ConsumerEntryVisitor{outer: g})
	for _, ext := range section.GetExtensions() {
		if cur, ok := lookupExtension(g.extensions, ext.Key); ok {
			g.extensions = putExtension(g.extensions, ext.Key,
				mergeAction.Merge(ext.Key, cur, ext.Value))
		} else {
			g.extensions = putExtension(g.extensions, ext.Key, ext.Value)
		}
	}
}

// AppendTo writes out the source map in the following format (line numbers
// are for reference only and are not part of the format):
//
//	{                                        // line 1
//	  version: 3,                            // line 2
//	  file: "out.js",                        // line 3
//	  lineCount: 2,                          // line 4
//	  sourceRoot: "",                        // line 5
//	  sources: ["foo.js", "bar.js"],         // line 6
//	  sourcesContent: ["var foo", "var bar"],// line 7
//	  names: ["src", "maps", "are", "fun"],  // line 8
//	  mappings: "a;;abcde,abcd,a;"           // line 9
//	  x_org_extension: value                 // line 10
//	}                                        // line 11
//
// Line 1 is the entire file as a single JSON object. Line 2 is the file
// version (always the first entry in the object). Line 3 [optional] is the
// name of the file that this source map is associated with. Line 4
// [optional] is the number of lines represented in the source map. Line 5
// [optional] is an optional source root, useful for relocating source files
// on a server or removing repeated prefix values in the "sources" entry.
// Line 6 is a list of sources used by the "mappings" entry relative to the
// sourceRoot. Line 7 is an optional list of the full content of the source
// files. Line 8 is a list of symbol names used by the "mapping" entry; this
// list may be incomplete. Line 9 is the mappings field. Line 10 is any
// custom field (extension).
//
// Java: SourceMapGeneratorV3.appendTo(Appendable, String)
func (g *SourceMapGeneratorV3) AppendTo(out io.StringWriter, name *string) {
	if g.parallelWorkers > 1 {
		g.appendToParallel(out, name)
		return
	}
	maxLine := jutil.IntAdd(g.PrepMappings(), 1)

	// Add the header fields.
	appendString(out, "{\n")
	SourceMapGeneratorV3_AppendFirstField(out, "version", "3")
	if name != nil {
		SourceMapGeneratorV3_AppendField(out, "file", SourceMapGeneratorV3_EscapeString(*name))
	}
	SourceMapGeneratorV3_AppendField(out, "lineCount", strconv.Itoa(maxLine))

	//optional source root
	if g.sourceRootPath != "" {
		SourceMapGeneratorV3_AppendField(out, "sourceRoot", SourceMapGeneratorV3_EscapeString(g.sourceRootPath))
	}

	// Add the mappings themselves.
	SourceMapGeneratorV3_AppendFieldStart(out, "mappings")
	// out.append("[");
	NewSourceMapGeneratorV3_LineMapper(g, out, maxLine).AppendLineMappings()

	// out.append("]");
	SourceMapGeneratorV3_AppendFieldEnd(out)
	g.appendToTail(out)
}

// appendToTail writes the fields after "mappings" (Go-only: shared with
// appendToParallel).
func (g *SourceMapGeneratorV3) appendToTail(out io.StringWriter) {
	// Files names
	SourceMapGeneratorV3_AppendFieldStart(out, "sources")
	appendString(out, "[")
	g.AddSourceNameMap(out)
	appendString(out, "]")
	SourceMapGeneratorV3_AppendFieldEnd(out)

	// Sources contents
	g.AddSourcesContentMap(out)

	// Identifier names
	SourceMapGeneratorV3_AppendFieldStart(out, "names")
	appendString(out, "[")
	g.AddSymbolNameMap(out)
	appendString(out, "]")
	SourceMapGeneratorV3_AppendFieldEnd(out)

	// Extensions, only if there is any
	for _, ext := range g.extensions {
		objValue := ext.Value
		var value string
		if s, ok := sourceMapGeneratorV3_asJavaString(objValue); ok {
			value = SourceMapGeneratorV3_EscapeString(s) // escapes native String
		} else {
			if sourceMapGeneratorV3_isNilExtension(objValue) {
				// objValue.toString()
				panic(jutil.NewNullPointerException())
			}
			value = jutil.ToJavaString(objValue)
		}
		SourceMapGeneratorV3_AppendField(out, ext.Key, value)
	}

	appendString(out, "\n}\n")
}

func sourceMapGeneratorV3_isNilExtension(v any) bool {
	switch x := v.(type) {
	case nil:
		return true
	case *string:
		return x == nil
	case *jutil.JsonObject:
		return x == nil
	case *jutil.JsonArray:
		return x == nil
	case *jutil.JsonPrimitive:
		return x == nil
	case *jutil.JsonNull:
		return x == nil
	default:
		return false
	}
}

// sourceMapGeneratorV3_asJavaString is `objValue instanceof String string`:
// a Java String is a Go string, or a non-nil *string.
func sourceMapGeneratorV3_asJavaString(objValue any) (string, bool) {
	switch v := objValue.(type) {
	case string:
		return v, true
	case *string:
		if v != nil {
			return *v, true
		}
	}
	return "", false
}

// SetSourceRoot sets a prefix to be added to the beginning of each
// sourceName passed to AddMapping. Debuggers expect (prefix + sourceName) to
// be a URL for loading the source code.
//
// path is the URL prefix to save in the sourcemap file. (Not validated.)
//
// Java: SourceMapGeneratorV3.setSourceRoot(String)
func (g *SourceMapGeneratorV3) SetSourceRoot(path string) {
	g.sourceRootPath = path
}

// AddExtension adds field extensions to the json source map. The value is
// allowed to be any value accepted by json, eg. string, JsonObject,
// JsonArray, etc.
//
// Extensions must follow the format x_orgranization_field (based on V3
// proposal), otherwise a SourceMapParseException is thrown (panicked).
//
// Java: SourceMapGeneratorV3.addExtension(String, Object)
func (g *SourceMapGeneratorV3) AddExtension(name string, object any) {
	if !strings.HasPrefix(name, "x_") {
		panic(NewSourceMapParseException("Extension '" + name +
			"' must start with 'x_'"))
	}
	g.extensions = putExtension(g.extensions, name, object)
}

// RemoveExtension removes an extension by name if present.
//
// Java: SourceMapGeneratorV3.removeExtension(String)
func (g *SourceMapGeneratorV3) RemoveExtension(name string) {
	g.extensions = removeExtension(g.extensions, name)
}

// HasExtension checks whether or not the sourcemap has an extension.
//
// Java: SourceMapGeneratorV3.hasExtension(String)
func (g *SourceMapGeneratorV3) HasExtension(name string) bool {
	_, ok := lookupExtension(g.extensions, name)
	return ok
}

// GetExtension returns the value mapped by the specified extension or nil
// if this extension does not exist.
//
// Java: SourceMapGeneratorV3.getExtension(String)
func (g *SourceMapGeneratorV3) GetExtension(name string) any {
	v, _ := lookupExtension(g.extensions, name)
	return v
}

// AddSourceNameMap writes the source name map to 'out'.
//
// Java: SourceMapGeneratorV3.addSourceNameMap(Appendable)
func (g *SourceMapGeneratorV3) AddSourceNameMap(out io.StringWriter) {
	SourceMapGeneratorV3_AddNameSlice(out, g.orderedSourceFiles)
}

// Java: SourceMapGeneratorV3.addSourcesContentMap(Appendable)
func (g *SourceMapGeneratorV3) AddSourcesContentMap(out io.StringWriter) {
	if len(g.sourceFileContentMap) == 0 {
		return
	}
	found := false
	for _, key := range g.orderedSourceFiles {
		if g.sourceFileContentMap[key].present {
			found = true
			break
		}
	}
	if !found {
		return
	}
	SourceMapGeneratorV3_AppendFieldStart(out, "sourcesContent")
	appendString(out, "[")
	for i, key := range g.orderedSourceFiles {
		if i != 0 {
			appendString(out, ",")
		}
		appendString(out, SourceMapGeneratorV3_EscapeString(g.sourceFileContentMap[key].text))
	}
	appendString(out, "]")
	SourceMapGeneratorV3_AppendFieldEnd(out)
}

// AddSymbolNameMap writes the source name map to 'out'.
//
// Java: SourceMapGeneratorV3.addSymbolNameMap(Appendable)
func (g *SourceMapGeneratorV3) AddSymbolNameMap(out io.StringWriter) {
	SourceMapGeneratorV3_AddNameSlice(out, g.orderedSymbolNames)
}

func SourceMapGeneratorV3_AddNameSlice(out io.StringWriter, names []string) {
	var buf [4096]byte
	bufLen := 0
	for i, key := range names {
		if !util_needsEscape(key) && len(key)+3 <= len(buf) {
			if bufLen+len(key)+3 > len(buf) {
				appendBytes(out, buf[:bufLen])
				bufLen = 0
			}
			if i != 0 {
				buf[bufLen] = ','
				bufLen++
			}
			buf[bufLen] = '"'
			bufLen++
			copy(buf[bufLen:], key)
			bufLen += len(key)
			buf[bufLen] = '"'
			bufLen++
		} else {
			if bufLen > 0 {
				appendBytes(out, buf[:bufLen])
				bufLen = 0
			}
			if i != 0 {
				appendString(out, ",")
			}
			appendString(out, SourceMapGeneratorV3_EscapeString(key))
		}
	}
	if bufLen > 0 {
		appendBytes(out, buf[:bufLen])
	}
}

// SourceMapGeneratorV3_EscapeString escapes the given string for JSON.
//
// Java: SourceMapGeneratorV3.escapeString(String)
func SourceMapGeneratorV3_EscapeString(value string) string {
	return Util_EscapeString(value)
}

// Source map field helpers.

// Java: SourceMapGeneratorV3.appendFirstField(Appendable, String, CharSequence)
func SourceMapGeneratorV3_AppendFirstField(out io.StringWriter, name string, value string) {
	SourceMapGeneratorV3_AppendFieldStartWithFirst(out, name, true)
	appendString(out, value)
}

// Java: SourceMapGeneratorV3.appendField(Appendable, String, CharSequence)
func SourceMapGeneratorV3_AppendField(out io.StringWriter, name string, value string) {
	SourceMapGeneratorV3_AppendFieldStartWithFirst(out, name, false)
	appendString(out, value)
}

// Java: SourceMapGeneratorV3.appendFieldStart(Appendable, String)
func SourceMapGeneratorV3_AppendFieldStart(out io.StringWriter, name string) {
	SourceMapGeneratorV3_AppendFieldStartWithFirst(out, name, false)
}

// Java: SourceMapGeneratorV3.appendFieldStart(Appendable, String, boolean)
func SourceMapGeneratorV3_AppendFieldStartWithFirst(out io.StringWriter, name string, first bool) {
	if !first {
		appendString(out, ",\n")
	}
	appendString(out, "\"")
	appendString(out, name)
	appendString(out, "\"")
	appendString(out, ":")
}

// Java: SourceMapGeneratorV3.appendFieldEnd(Appendable)
func SourceMapGeneratorV3_AppendFieldEnd(out io.StringWriter) {}

// PrepMappings assigns sequential ids to used mappings, and returns the last
// line mapped.
//
// Java: SourceMapGeneratorV3.prepMappings()
func (g *SourceMapGeneratorV3) PrepMappings() int {
	var check SourceMapGeneratorV3_UsedMappingCheck
	t := SourceMapGeneratorV3_MappingTraversal{outer: g}
	t.Traverse(&check)

	// Adjust for the prefix.
	return jutil.IntAdd(check.maxLine, g.prefixPosition.GetLine())
}

// AppendIndexMapTo appends the index source map to the given buffer.
//
// out is the stream to which the map will be appended; name is the name of
// the generated source file that this source map represents; sections is an
// ordered list of map sections to include in the index.
//
// Java: SourceMapGeneratorV3.appendIndexMapTo(Appendable, String, List)
func (g *SourceMapGeneratorV3) AppendIndexMapTo(out io.StringWriter, name string, sections []*SourceMapSection) {
	// Add the header fields.
	appendString(out, "{\n")
	SourceMapGeneratorV3_AppendFirstField(out, "version", "3")
	SourceMapGeneratorV3_AppendField(out, "file", SourceMapGeneratorV3_EscapeString(name))

	// Add the line character maps.
	SourceMapGeneratorV3_AppendFieldStart(out, "sections")
	appendString(out, "[\n")
	first := true
	for _, section := range sections {
		if first {
			first = false
		} else {
			appendString(out, ",\n")
		}
		appendString(out, "{\n")
		SourceMapGeneratorV3_AppendFieldStartWithFirst(out, "offset", true)
		SourceMapGeneratorV3_AppendOffsetValue(out, section.line, section.column)
		if section.type_ == SourceMapSection_SectionType_URL {
			if !section.value.present {
				panic(jutil.NewNullPointerException())
			}
			SourceMapGeneratorV3_AppendField(out, "url", SourceMapGeneratorV3_EscapeString(section.value.text))
		} else if section.type_ == SourceMapSection_SectionType_MAP {
			// Appendable.append(null) appends "null".
			mapVal := "null"
			if section.value.present {
				mapVal = section.value.text
			}
			SourceMapGeneratorV3_AppendField(out, "map", mapVal)
		} else {
			panic(jutil.NewIOException("Unexpected section type"))
		}
		appendString(out, "\n}")
	}

	appendString(out, "\n]")
	SourceMapGeneratorV3_AppendFieldEnd(out)

	appendString(out, "\n}\n")
}

// Java: SourceMapGeneratorV3.appendOffsetValue(Appendable, int, int)
func SourceMapGeneratorV3_AppendOffsetValue(out io.StringWriter, line, column int) {
	appendString(out, "{\n")
	SourceMapGeneratorV3_AppendFirstField(out, "line", strconv.Itoa(line))
	SourceMapGeneratorV3_AppendField(out, "column", strconv.Itoa(column))
	appendString(out, "\n}")
}
