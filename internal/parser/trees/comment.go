// Port of com/google/javascript/jscomp/parsing/parser/trees/Comment.java
// @ closure-compiler v20260415

package trees

// CommentType mirrors parsing/parser/trees/Comment.Type in Java.
type CommentType int

const (
	CommentTypeBlock CommentType = iota
	CommentTypeLine
	CommentTypeJSDoc
	CommentTypeImportant
	CommentTypeShebang
)

// Comment represents a single comment in the parse tree.
// Mirrors parsing/parser/trees/Comment.java.
type Comment struct {
	Value    string
	Location SourceRange
	Type     CommentType
}

// NewComment constructs a Comment.
// Go-only: language-adapter
func NewComment(value string, location SourceRange, kind CommentType) *Comment {
	return &Comment{Value: value, Location: location, Type: kind}
}

// IsJsDoc mirrors Java's Comment.isJsDoc().
func (c *Comment) IsJsDoc() bool { return c.Type == CommentTypeJSDoc }
