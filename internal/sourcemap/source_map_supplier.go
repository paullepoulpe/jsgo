// Port of com/google/debugging/sourcemap/SourceMapSupplier.java @ closure-compiler v20260415

package sourcemap

// SourceMapSupplier is a class for mapping source map names to the actual
// contents. Used when parsing index maps.
type SourceMapSupplier interface {
	// GetSourceMap returns the contents of the map associated with the URL
	// (nil for Java null). It may panic with a jutil IOException.
	//
	// Java: SourceMapSupplier.getSourceMap(String)
	GetSourceMap(url string) *string
}
