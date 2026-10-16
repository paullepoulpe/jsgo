// Port of com/google/javascript/jscomp/parsing/parser/util/SourceRange.java @ closure-compiler v20260415

package parser

import "github.com/paullepoulpe/jsgo/internal/jutil"

// SourceRange is a range of positions in a source string.
// Start is inclusive.
// End is exclusive.
type SourceRange struct {
	Start *SourcePosition
	End   *SourcePosition
}

// NewSourceRange constructs a SourceRange.
//
// Java: SourceRange(SourcePosition, SourcePosition)
func NewSourceRange(start, end *SourcePosition) *SourceRange {
	s := &SourceRange{}
	s.Start = start
	s.End = end
	return s
}

// String renders the range as "<start - end>".
//
// Java: SourceRange.toString()
func (s *SourceRange) String() string {
	return jutil.StringFormat("<%s - %s>", s.Start, s.End)
}
