// Port of com/google/debugging/sourcemap/FilePosition.java @ closure-compiler v20261021
// Port of com/google/debugging/sourcemap/SourceMapConsumer.java @ closure-compiler v20261021
// Port of com/google/debugging/sourcemap/SourceMapConsumerFactory.java @ closure-compiler v20261021
// Port of com/google/debugging/sourcemap/SourceMapFormat.java @ closure-compiler v20261021
// Port of com/google/debugging/sourcemap/SourceMapGenerator.java @ closure-compiler v20261021
// Port of com/google/debugging/sourcemap/SourceMapGeneratorFactory.java @ closure-compiler v20261021
// Port of com/google/debugging/sourcemap/SourceMapParseException.java @ closure-compiler v20261021
// Port of com/google/debugging/sourcemap/SourceMapSupplier.java @ closure-compiler v20261021
// Port of com/google/debugging/sourcemap/SourceMapping.java @ closure-compiler v20261021
// Port of com/google/debugging/sourcemap/SourceMappingReversable.java @ closure-compiler v20261021

/*
 * Copyright 2009 The Closure Compiler Authors.
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

// Package sourcemap is a 1:1 port of com.google.debugging.sourcemap
// (closure-compiler v20261021): the source map V3 generator
// (SourceMapGeneratorV3), consumer (SourceMapConsumerV3), their JSON
// parsers (SourceMapObjectParser on Gson, and the SourceMapJsonLexer fast
// path), and the Base64 VLQ encoding. proto_mapping.go reconstructs the
// protobuf-generated com.google.debugging.sourcemap.proto.Mapping.
//
// Java Strings that may be null at public boundaries are *string; Java arrays
// that may be null are nil slices (a non-null empty array is a non-nil empty
// slice); a Java Appendable is an io.StringWriter whose write errors become
// IOExceptions. See PORTING.md for the translation rules.
package sourcemap

import (
	"io"
	"strconv"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// optionalString stores a nullable Java String by value without heap-allocating
// a *string wrapper.
type optionalString struct {
	text    string
	present bool
}

func makeOptionalString(p *string) optionalString {
	if p == nil {
		return optionalString{}
	}
	return optionalString{text: *p, present: true}
}

func (s *optionalString) ptr() *string {
	if !s.present {
		return nil
	}
	return &s.text
}

func optionalStringsFromPtrs(xs []*string) []optionalString {
	if xs == nil {
		return nil
	}
	out := make([]optionalString, len(xs))
	for i, p := range xs {
		if p != nil {
			out[i] = optionalString{text: *p, present: true}
		}
	}
	return out
}

func optionalStringsToPtrs(xs []optionalString) []*string {
	if xs == nil {
		return nil
	}
	out := make([]*string, len(xs))
	for i := range xs {
		if xs[i].present {
			out[i] = &xs[i].text
		}
	}
	return out
}

// FilePosition represents a position in a source file.
type FilePosition struct {
	line   int32
	column int32
}

// NewFilePosition constructs a FilePosition.
//
// Java: FilePosition(int, int)
func NewFilePosition(line, column int) *FilePosition {
	return &FilePosition{line: int32(line), column: int32(column)}
}

// MakeFilePosition constructs a FilePosition by value without heap allocation.
func MakeFilePosition(line, column int) FilePosition {
	return FilePosition{line: int32(line), column: int32(column)}
}

// GetLine returns the line number of this position. Note: The v1 and v2
// source maps use a line number with the first line being 1, whereas the v3
// source map corrects this and uses a first line number of 0 to be
// consistent with the column representation.
//
// Java: FilePosition.getLine()
func (p *FilePosition) GetLine() int {
	return int(p.line)
}

// GetColumn returns the character index on the line of this position, with
// the first column being 0.
//
// Java: FilePosition.getColumn()
func (p *FilePosition) GetColumn() int {
	return int(p.column)
}

// String is for debugging only.
//
// Java: FilePosition.toString()
func (p *FilePosition) String() string {
	return "FilePosition{line=" + strconv.Itoa(int(p.line)) + ", column=" + strconv.Itoa(int(p.column)) + "}"
}

// SourceMapping is an interface for provide a way of mapping (line, column)
// positions back to positions in the original (uncompiled) source code.
type SourceMapping interface {
	// GetMappingForLine returns the original mapping for the line number and
	// column position found in the source map. Returns null if none is
	// found.
	//
	// lineNumber is the line number, with the first being '1'; columnIndex
	// is the column index, with the first being '1'.
	//
	// Java: SourceMapping.getMappingForLine(int, int)
	GetMappingForLine(lineNumber, columnIndex int) *Mapping_OriginalMapping
}

// SourceMappingReversable is a SourceMapping that can provide the reverse
// (source → target) source mapping.
type SourceMappingReversable interface {
	SourceMapping

	// GetOriginalSources returns the collection of original sources in this
	// source mapping.
	//
	// Java: SourceMappingReversable.getOriginalSources()
	GetOriginalSources() []*string

	// GetReverseMapping returns, given a source file, line, and column, the
	// reverse mapping (source → target). A collection is returned as in some
	// cases (like a function being inlined), one source line may map to more
	// then one target location. An empty collection is returned if there
	// were no matches.
	//
	// Java: SourceMappingReversable.getReverseMapping(String, int, int)
	GetReverseMapping(originalFile string, line, column int) []*Mapping_OriginalMapping
}

// SourceMapConsumer is a SourceMapping provide that can parse from a raw
// string.
type SourceMapConsumer interface {
	SourceMapping

	// Parse parses the given contents containing a source map to provide
	// initialize a class providing SourceMapping. It panics with a
	// *SourceMapParseException.
	//
	// Java: SourceMapConsumer.parse(String)
	Parse(contents string)
}

// SourceMapSupplier is a class for mapping source map names to the actual
// contents. Used when parsing index maps.
type SourceMapSupplier interface {
	// GetSourceMap returns the contents of the map associated with the URL
	// (nil for Java null). It may panic with a jutil IOException.
	//
	// Java: SourceMapSupplier.getSourceMap(String)
	GetSourceMap(url string) *string
}

// SourceMapFormat is a list of currently support SourceMap format
// revisions.
type SourceMapFormat int32

const (
	// The latest "stable" format
	SourceMapFormat_DEFAULT SourceMapFormat = iota

	// V3: A nice compact format
	SourceMapFormat_V3
)

var sourceMapFormat_names = [...]string{"DEFAULT", "V3"}

// String is Enum.name().
func (f SourceMapFormat) String() string { return sourceMapFormat_names[f] }

// SourceMapGenerator collects information mapping the generated (compiled)
// source back to its original source for debugging purposes.
type SourceMapGenerator interface {
	// AppendTo appends the source map to the given buffer.
	//
	// out is the stream to which the map will be appended; name is the name
	// of the generated source file that this source map represents.
	//
	// Java: SourceMapGenerator.appendTo(Appendable, String)
	AppendTo(out io.StringWriter, name *string)

	// AppendIndexMapTo appends the index source map to the given buffer.
	//
	// out is the stream to which the map will be appended; name is the name
	// of the generated source file that this source map represents; sections
	// is an ordered list of map sections to include in the index.
	//
	// Java: SourceMapGenerator.appendIndexMapTo(Appendable, String, List)
	AppendIndexMapTo(out io.StringWriter, name string, sections []*SourceMapSection)

	// Reset resets the source map for reuse. A reset needs to be called
	// between each generated output file.
	//
	// Java: SourceMapGenerator.reset()
	Reset()

	// AddMapping adds a mapping for the given node. Mappings must be added in
	// order.
	//
	// sourceName is the file name to use in the generate source map to
	// represent this source; symbolName is the symbol name associated with
	// this position in the source map (nil for none); sourceStartPosition is
	// the starting position in the original source for represented range
	// outputStartPosition to outputEndPosition in the generated file;
	// outputStartPosition is the position on the starting line;
	// outputEndPosition is the position on the ending line.
	//
	// Java: SourceMapGenerator.addMapping(String, String, FilePosition, FilePosition, FilePosition)
	AddMapping(sourceName *string, symbolName *string,
		sourceStartPosition *FilePosition,
		outputStartPosition *FilePosition, outputEndPosition *FilePosition)

	// AddSourcesContent adds a sourcesContent entry for the corresponding
	// source.
	//
	// Java: SourceMapGenerator.addSourcesContent(String, String)
	AddSourcesContent(source string, content *string)

	// SetWrapperPrefix sets the prefix used for wrapping the generated source
	// file before it is written. This ensures that the source map is adjusted
	// for the change in character offsets.
	//
	// Java: SourceMapGenerator.setWrapperPrefix(String)
	SetWrapperPrefix(prefix string)

	// SetStartingPosition sets the source code that exists in the buffer for
	// which the generated code is being generated. This ensures that the
	// source map accurately reflects the fact that the source is being
	// appended to an existing buffer and as such, does not start at line 0,
	// position 0 but rather some other line and position.
	//
	// offsetLine is the index of the current line being printed; offsetIndex
	// is the column index of the current character being printed.
	//
	// Java: SourceMapGenerator.setStartingPosition(int, int)
	SetStartingPosition(offsetLine, offsetIndex int)

	// Validate sets whether to perform additional validation on the source
	// map.
	//
	// Java: SourceMapGenerator.validate(boolean)
	Validate(validate bool)
}

// SourceMapGeneratorFactory_GetInstance returns the appropriate source map
// object for the given source map format.
//
// Java: SourceMapGeneratorFactory.getInstance(SourceMapFormat)
func SourceMapGeneratorFactory_GetInstance(format SourceMapFormat) SourceMapGenerator {
	switch format {
	case SourceMapFormat_DEFAULT, SourceMapFormat_V3:
		return NewSourceMapGeneratorV3()
	}
	// Go-only: language-adapter — unreachable for a valid enum value: the
	// Java switch expression is exhaustive.
	panic(jutil.NewIllegalStateException())
}

// SourceMapConsumerFactory_Parse parses the string representing the source
// map file contents. It panics with a *SourceMapParseException.
//
// Java: SourceMapConsumerFactory.parse(String)
func SourceMapConsumerFactory_Parse(contents string) SourceMapping {
	return SourceMapConsumerFactory_ParseWithSupplier(contents, nil)
}

// SourceMapConsumerFactory_ParseWithSupplier parses the string representing
// the source map file contents, using supplier (may be nil) for any
// referenced maps.
//
// Java: SourceMapConsumerFactory.parse(String, SourceMapSupplier)
func SourceMapConsumerFactory_ParseWithSupplier(contents string, supplier SourceMapSupplier) SourceMapping {
	// Version 1, starts with a magic string
	if strings.HasPrefix(contents, "/** Begin line maps. **/") {
		panic(NewSourceMapParseException(
			"This appears to be a V1 SourceMap, which is not supported."))
	} else if strings.HasPrefix(contents, "{") {
		sourceMapObject := SourceMapObjectParser_Parse(contents)

		// Check basic assertions about the format.
		switch sourceMapObject.version {
		case 3:
			consumer := NewSourceMapConsumerV3()
			consumer.ParseSourceMapObject(sourceMapObject, supplier)
			return consumer
		default:
			panic(NewSourceMapParseException(
				"Unknown source map version:" + strconv.Itoa(sourceMapObject.version)))
		}
	}

	panic(NewSourceMapParseException("unable to detect source map format"))
}

// SourceMapConsumerFactory_ParseFast uses the fast parser to parse the
// source map. This only supports version 3 source maps.
//
// Java: SourceMapConsumerFactory.parseFast(String)
func SourceMapConsumerFactory_ParseFast(contents string) *SourceMapConsumerV3 {
	if strings.HasPrefix(contents, "{") {
		sourceMapObject := SourceMapObjectParser_ParseFast(contents)
		jutil.CheckArgumentMsg1(
			sourceMapObject.version == 3,
			"Unknown source map version:%s",
			sourceMapObject.version)
		consumer := NewSourceMapConsumerV3()
		consumer.ParseSourceMapObject(sourceMapObject, nil)
		return consumer
	}
	panic(NewSourceMapParseException("unable to detect source map format"))
}

// SourceMapParseException is thrown if an invalid or unknown source map is
// encountered. It is a checked exception (java.lang.Exception).
type SourceMapParseException struct{ jutil.JavaException }

func (*SourceMapParseException) sourceMapParseException() {}

// ISourceMapParseException matches SourceMapParseException.
type ISourceMapParseException interface {
	jutil.IException
	sourceMapParseException()
}

// NewSourceMapParseException constructs a SourceMapParseException.
//
// Java: SourceMapParseException(String)
func NewSourceMapParseException(message string) *SourceMapParseException {
	e := &SourceMapParseException{}
	jutil.InitThrowable(e, "com.google.debugging.sourcemap.SourceMapParseException", &message, nil)
	return e
}
