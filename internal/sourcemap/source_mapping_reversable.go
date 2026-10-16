// Port of com/google/debugging/sourcemap/SourceMappingReversable.java @ closure-compiler v20260415

package sourcemap

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
