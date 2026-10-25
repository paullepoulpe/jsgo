// Port of com/google/javascript/jscomp/parsing/parser/trees/Comment.java @ closure-compiler v20261021

/*
 * Copyright 2014 The Closure Compiler Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package parser

// Comment is a placeholder class.
type Comment struct {
	Value    string
	Location LineRange // Go-only: with lines and columns (see SourcePosition)
	Type     Comment_Type
}

// Comment_Type is Comment.Type.
type Comment_Type int32

const (
	// /* comment */
	Comment_Type_BLOCK Comment_Type = iota

	// // comment
	Comment_Type_LINE

	// /** comment */
	Comment_Type_JSDOC

	// #!/usr/bin/node
	// Only valid at the start of a file.
	Comment_Type_SHEBANG

	// /*! comment */
	Comment_Type_IMPORTANT
)

var comment_Type_names = [...]string{"BLOCK", "LINE", "JSDOC", "SHEBANG", "IMPORTANT"}

// Comment_Type_values is Comment.Type.values().
//
// Java: Comment.Type.values()
func Comment_Type_values() []Comment_Type {
	return []Comment_Type{Comment_Type_BLOCK, Comment_Type_LINE, Comment_Type_JSDOC, Comment_Type_SHEBANG, Comment_Type_IMPORTANT}
}

// Name is Enum.name().
//
// Java: Comment.Type.name()
func (t Comment_Type) Name() string { return comment_Type_names[t] }

// String is Enum.toString().
//
// Java: Comment.Type.toString()
func (t Comment_Type) String() string { return t.Name() }

// NewComment constructs a Comment.
//
// Java: Comment(String, SourceRange, Type)
func NewComment(value string, location LineRange, type_ Comment_Type) *Comment {
	c := &Comment{}
	c.Value = value
	c.Location = location
	c.Type = type_
	return c
}

// IsJsDoc reports whether this is a JSDoc comment.
//
// Java: Comment.isJsDoc()
func (c *Comment) IsJsDoc() bool {
	return c.Type == Comment_Type_JSDOC
}

// GetAbsolutePosition returns the start offset.
//
// Java: Comment.getAbsolutePosition()
func (c *Comment) GetAbsolutePosition() int {
	return int(c.Location.Start.Offset)
}

// GetLength returns the length in UTF-16 code units.
//
// Java: Comment.getLength()
func (c *Comment) GetLength() int {
	return int(c.Location.End.Offset - c.Location.Start.Offset)
}
