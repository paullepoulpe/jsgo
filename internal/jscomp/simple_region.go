// Port of com/google/javascript/jscomp/SimpleRegion.java @ closure-compiler v20260415

package jscomp

// SimpleRegion is a simple region.
type SimpleRegion struct {
	beginningLineNumber int
	endingLineNumber    int
	source              string
}

// Java: SimpleRegion(int, int, String)
func NewSimpleRegion(beginningLineNumber int, endingLineNumber int, source string) *SimpleRegion {
	return &SimpleRegion{beginningLineNumber: beginningLineNumber, endingLineNumber: endingLineNumber, source: source}
}

// Java: SimpleRegion.getBeginningLineNumber()
func (r *SimpleRegion) GetBeginningLineNumber() int { return r.beginningLineNumber }

// Java: SimpleRegion.getEndingLineNumber()
func (r *SimpleRegion) GetEndingLineNumber() int { return r.endingLineNumber }

// Java: SimpleRegion.getSourceExcerpt()
func (r *SimpleRegion) GetSourceExcerpt() string { return r.source }
