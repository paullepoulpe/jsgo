// Port of com/google/debugging/sourcemap/SourceMapGenerator.java @ closure-compiler v20260415

package sourcemap

import "io"

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
	AddSourcesContent(source string, content string)

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
