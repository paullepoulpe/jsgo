// Port of com/google/debugging/sourcemap/SourceMapConsumer.java @ closure-compiler v20260415

package sourcemap

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
