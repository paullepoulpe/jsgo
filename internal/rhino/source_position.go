// Port of com/google/javascript/rhino/SourcePosition.java
// @ closure-compiler v20260415
//
// Java's SourcePosition<T> is a generic abstract class. In Go we use a
// concrete struct with an any item field, preserving the full public API.

package rhino

// SourcePosition represents a position in some piece of source code, with an
// associated item found at that position. Mirrors
// com.google.javascript.rhino.SourcePosition<T>.
type SourcePosition struct {
	item       any
	startLineno int
	startCharno int
	endLineno   int
	endCharno   int
}

// SetItem sets the item that this source position references.
func (s *SourcePosition) SetItem(item any) {
	s.item = item
}

// SetPositionInformation sets the position information contained in this
// source position. Mirrors SourcePosition.setPositionInformation.
// Panics if the start/end positions are inconsistent.
func (s *SourcePosition) SetPositionInformation(startLineno, startCharno, endLineno, endCharno int) {
	if startLineno > endLineno {
		panic("Recorded bad position information\nstart-line: " + itoa(startLineno) + "\nend-line: " + itoa(endLineno))
	} else if startLineno == endLineno && startCharno >= endCharno {
		panic("Recorded bad position information\nline: " + itoa(startLineno) + "\nstart-char: " + itoa(startCharno) + "\nend-char: " + itoa(endCharno))
	}
	s.startLineno = startLineno
	s.startCharno = startCharno
	s.endLineno = endLineno
	s.endCharno = endCharno
}

// GetItem returns the item found at this source position.
func (s *SourcePosition) GetItem() any {
	return s.item
}

// GetStartLine returns the starting line number of this position.
func (s *SourcePosition) GetStartLine() int {
	return s.startLineno
}

// GetPositionOnStartLine returns the character position on the starting line.
func (s *SourcePosition) GetPositionOnStartLine() int {
	return s.startCharno
}

// GetEndLine returns the ending line number of this position.
func (s *SourcePosition) GetEndLine() int {
	return s.endLineno
}

// GetPositionOnEndLine returns the character position on the ending line.
func (s *SourcePosition) GetPositionOnEndLine() int {
	return s.endCharno
}

// IsSamePositionAs reports whether this position has the same line/char
// extents as that. Mirrors SourcePosition.isSamePositionAs.
func (s *SourcePosition) IsSamePositionAs(that *SourcePosition) bool {
	return s.startLineno == that.startLineno &&
		s.startCharno == that.startCharno &&
		s.endLineno == that.endLineno &&
		s.endCharno == that.endCharno
}
