// Port of com/google/debugging/sourcemap/SourceMapFormat.java @ closure-compiler v20260415

package sourcemap

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

// SourceMapFormat_values is SourceMapFormat.values().
func SourceMapFormat_values() []SourceMapFormat {
	return []SourceMapFormat{SourceMapFormat_DEFAULT, SourceMapFormat_V3}
}

// String is Enum.name().
func (f SourceMapFormat) String() string { return sourceMapFormat_names[f] }
