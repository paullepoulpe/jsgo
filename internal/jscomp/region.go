// Port of com/google/javascript/jscomp/Region.java @ closure-compiler v20260415

package jscomp

// Region is a source code region.
type Region interface {
	// GetSourceExcerpt gets the source region.
	GetSourceExcerpt() string

	// GetBeginningLineNumber gets the beginning line number.
	GetBeginningLineNumber() int

	// GetEndingLineNumber gets the ending line number.
	GetEndingLineNumber() int
}
