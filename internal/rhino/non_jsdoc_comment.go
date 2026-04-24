// Port of com/google/javascript/rhino/NonJSDocComment.java
// @ closure-compiler v20260415
//
// The Java class uses com.google.javascript.jscomp.parsing.parser.util.SourcePosition
// (a parser-level position with offset/line/column), not the rhino-level
// SourcePosition<T>. In Go we use a flat CommentPosition struct for the
// same fields rather than importing a not-yet-ported parser package.

package rhino

// CommentPosition records offset, line, and column for a comment boundary.
// Mirrors com.google.javascript.jscomp.parsing.parser.util.SourcePosition
// for the purpose of NonJSDocComment.
type CommentPosition struct {
	Offset int
	Line   int
	Column int
}

// NonJSDocComment holds minimal information about a non-JSDoc comment's
// source location and contents. Mirrors
// com.google.javascript.rhino.NonJSDocComment.
type NonJSDocComment struct {
	startPosition     CommentPosition
	endPosition       CommentPosition
	contents          string
	endsAsLineComment bool
	isInline          bool
}

// NewNonJSDocComment constructs a NonJSDocComment.
func NewNonJSDocComment(start, end CommentPosition, s string) *NonJSDocComment {
	return &NonJSDocComment{
		startPosition: start,
		endPosition:   end,
		contents:      s,
	}
}

// GetCommentString returns the comment text, or "" if nil.
func (c *NonJSDocComment) GetCommentString() string {
	if c.contents == "" {
		return ""
	}
	return c.contents
}

// GetStartPosition returns the start position of the comment.
func (c *NonJSDocComment) GetStartPosition() CommentPosition {
	return c.startPosition
}

// GetEndPosition returns the end position of the comment.
func (c *NonJSDocComment) GetEndPosition() CommentPosition {
	return c.endPosition
}

// SetEndsAsLineComment records whether the comment ends as a line comment.
func (c *NonJSDocComment) SetEndsAsLineComment(v bool) {
	c.endsAsLineComment = v
}

// SetIsInline records whether the comment is an inline comment.
func (c *NonJSDocComment) SetIsInline(v bool) {
	c.isInline = v
}

// IsEndingAsLineComment reports whether this comment ends as a line comment.
func (c *NonJSDocComment) IsEndingAsLineComment() bool {
	return c.endsAsLineComment
}

// IsInline reports whether this is an inline comment.
func (c *NonJSDocComment) IsInline() bool {
	return c.isInline
}

// String returns a diagnostic representation.
func (c *NonJSDocComment) String() string {
	return "NonJSDocComment : " + c.GetCommentString()
}
