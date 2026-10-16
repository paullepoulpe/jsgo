// Port of com/google/javascript/rhino/NonJSDocComment.java @ closure-compiler v20260415

package rhino

import "github.com/paullepoulpe/jsgo/internal/parser"

// NonJSDocComment is a minimal class holding information about a nonJSDoc
// comment's source location and contents.
type NonJSDocComment struct {
	startPosition     *parser.SourcePosition
	endPosition       *parser.SourcePosition
	contents          *string
	endsAsLineComment bool
	isInline          bool
}

// Java: NonJSDocComment(SourcePosition, SourcePosition, String)
func NewNonJSDocComment(start *parser.SourcePosition, end *parser.SourcePosition, s *string) *NonJSDocComment {
	c := &NonJSDocComment{}
	c.startPosition = start
	c.endPosition = end
	c.contents = s
	c.endsAsLineComment = false
	c.isInline = false
	return c
}

// Java: NonJSDocComment.getCommentString()
func (c *NonJSDocComment) GetCommentString() string {
	if c.contents == nil {
		return ""
	}
	return *c.contents
}

// Java: NonJSDocComment.getStartPosition()
func (c *NonJSDocComment) GetStartPosition() *parser.SourcePosition {
	return c.startPosition
}

// Java: NonJSDocComment.getEndPosition()
func (c *NonJSDocComment) GetEndPosition() *parser.SourcePosition {
	return c.endPosition
}

// Java: NonJSDocComment.setEndsAsLineComment(boolean)
func (c *NonJSDocComment) SetEndsAsLineComment(endsAsLineComment bool) {
	c.endsAsLineComment = endsAsLineComment
}

// Java: NonJSDocComment.setIsInline(boolean)
func (c *NonJSDocComment) SetIsInline(isInline bool) {
	c.isInline = isInline
}

// Java: NonJSDocComment.isEndingAsLineComment()
func (c *NonJSDocComment) IsEndingAsLineComment() bool {
	return c.endsAsLineComment
}

// Java: NonJSDocComment.isInline()
func (c *NonJSDocComment) IsInline() bool {
	return c.isInline
}

// Java: NonJSDocComment.toString()
func (c *NonJSDocComment) String() string {
	return "NonJSDocComment : " + c.GetCommentString()
}
