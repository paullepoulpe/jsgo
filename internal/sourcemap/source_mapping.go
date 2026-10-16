// Port of com/google/debugging/sourcemap/SourceMapping.java @ closure-compiler v20260415

package sourcemap

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
