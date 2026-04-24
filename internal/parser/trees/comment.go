// Port of com/google/javascript/jscomp/parsing/parser/trees/Comment.java
// @ closure-compiler v20260415

package trees

// CommentType identifies the kind of a Comment.
type CommentType int

const (
	// CommentTypeBlock is a /* comment */.
	CommentTypeBlock CommentType = iota
	// CommentTypeLine is a // comment.
	CommentTypeLine
	// CommentTypeJsDoc is a /** comment */.
	CommentTypeJsDoc
	// CommentTypeShebang is a #!/usr/bin/node line; only valid at the start of a file.
	CommentTypeShebang
	// CommentTypeImportant is a /*! comment */.
	CommentTypeImportant
)

// Comment holds a source comment and its position.
type Comment struct {
	Value    string
	Location SourceRange
	Type     CommentType
}

func NewComment(value string, location SourceRange, commentType CommentType) *Comment {
	return &Comment{
		Value:    value,
		Location: location,
		Type:     commentType,
	}
}

func (c *Comment) IsJsDoc() bool {
	return c.Type == CommentTypeJsDoc
}

func (c *Comment) GetAbsolutePosition() int {
	return c.Location.Start.Offset
}

func (c *Comment) GetLength() int {
	return c.Location.End.Offset - c.Location.Start.Offset
}
