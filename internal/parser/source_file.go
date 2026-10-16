// Port of com/google/javascript/jscomp/parsing/parser/SourceFile.java @ closure-compiler v20260415

package parser

import "github.com/paullepoulpe/jsgo/internal/jutil"

// SourceFile is a source file.
//
// Immutable.
type SourceFile struct {
	Name     string
	Contents string

	// Go-only: language-adapter — Contents as UTF-16 code units (Java chars),
	// converted once, so that the scanner and every offset work in Java char
	// indices. Read-only.
	ContentsUTF16 []uint16
}

// NewSourceFile constructs a SourceFile.
//
// Java: SourceFile(String, String)
func NewSourceFile(name, contents string) *SourceFile {
	s := &SourceFile{}
	s.Name = name
	s.Contents = contents
	s.ContentsUTF16 = jutil.StringToUTF16(contents)
	return s
}
