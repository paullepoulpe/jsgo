// Port of com/google/javascript/rhino/StaticSourceFile.java @ closure-compiler v20261015

/*
 *
 * ***** BEGIN LICENSE BLOCK *****
 * Version: MPL 1.1
 *
 * The contents of this file are subject to the Mozilla Public License Version
 * 1.1 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 * http://www.mozilla.org/MPL/
 *
 * Software distributed under the License is distributed on an "AS IS" basis,
 * WITHOUT WARRANTY OF ANY KIND, either express or implied. See the License
 * for the specific language governing rights and limitations under the
 * License.
 *
 * The Original Code is Rhino code, released
 * May 6, 1999.
 *
 * The Initial Developer of the Original Code is
 * Netscape Communications Corporation.
 * Portions created by the Initial Developer are Copyright (C) 1997-1999
 * the Initial Developer. All Rights Reserved.
 *
 * Contributor(s):
 *   Nick Santos
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */

package rhino

import "strings"

// StaticSourceFile contains information about a compiler input.
//
// Java default methods are declared in the interface (implementations may
// override them, as in Java) and their default bodies are the
// StaticSourceFile_X functions, which implementations call to inherit them.
type StaticSourceFile interface {
	// GetName is the name of the file. Must be unique across all files in the
	// compilation.
	GetName() string

	// GetKind is the source kind.
	GetKind() StaticSourceFile_SourceKind

	// IsStrong reports whether the source kind is STRONG.
	IsStrong() bool

	// IsWeak reports whether the source kind is WEAK.
	IsWeak() bool

	// IsExtern reports whether the source kind is EXTERN.
	IsExtern() bool

	// IsNonCode reports whether the source kind is NON_CODE.
	IsNonCode() bool

	IsTypeScriptSource() bool

	// IsClosureUnawareCode reports whether this source file contains code
	// that should be optimized using a minimal set of known-safe
	// optimizations. In particular, code in this StaticSourceFile is unlikely
	// to be well-typed (or have types annotated at all), and is likely written
	// by someone unaware of the limitations of this compiler.
	IsClosureUnawareCode() bool

	// GetLineOffset returns the offset of the given line number relative to
	// the file start. Line number should be 1-based.
	//
	// If the source file doesn't have line information, it should return
	// Integer.MIN_VALUE. The negative offsets will make it more obvious what
	// happened.
	//
	// Throws IllegalArgumentException if lineno is less than 1 or greater
	// than the number of lines in the source.
	GetLineOffset(lineNumber int) int

	// GetLineOfOffset gets the 1-based line number of the given source offset.
	// The behavior is undefined if this offset does not exist in the source
	// file.
	GetLineOfOffset(offset int) int

	// GetColumnOfOffset gets the 0-based column number of the given source
	// offset. The behavior is undefined if this offset does not exist in the
	// source file.
	GetColumnOfOffset(offset int) int
}

// StaticSourceFileFast is an optional Go-only fast-path interface implemented
// by SourceFile and SimpleSourceFile so Node.GetSourceFileName avoids *string
// heap boxing.
type StaticSourceFileFast interface {
	GetNamePtr() *string
}

// StaticSourceFile_SourceKind enumerates source kinds.
//
// Java: StaticSourceFile.SourceKind
type StaticSourceFile_SourceKind int32

const (
	// A file whose contents are necessary both for type checking and emitting
	// code.
	StaticSourceFile_SourceKind_STRONG StaticSourceFile_SourceKind = iota
	// A file whose contents are necessary for type checking only.
	StaticSourceFile_SourceKind_WEAK
	// A file whose contents are extern declarations.
	StaticSourceFile_SourceKind_EXTERN
	// A file that is not part of the compilation unit and is not strong, weak,
	// or an extern
	StaticSourceFile_SourceKind_NON_CODE
)

var staticSourceFile_SourceKind_names = [...]string{"STRONG", "WEAK", "EXTERN", "NON_CODE"}

// StaticSourceFile_SourceKind_values is SourceKind.values().
func StaticSourceFile_SourceKind_values() []StaticSourceFile_SourceKind {
	return []StaticSourceFile_SourceKind{
		StaticSourceFile_SourceKind_STRONG, StaticSourceFile_SourceKind_WEAK,
		StaticSourceFile_SourceKind_EXTERN, StaticSourceFile_SourceKind_NON_CODE,
	}
}

// String is Enum.name().
func (k StaticSourceFile_SourceKind) String() string { return staticSourceFile_SourceKind_names[k] }

// Java: StaticSourceFile.isStrong() (default method)
func StaticSourceFile_IsStrong(f StaticSourceFile) bool {
	return f.GetKind() == StaticSourceFile_SourceKind_STRONG
}

// Java: StaticSourceFile.isWeak() (default method)
func StaticSourceFile_IsWeak(f StaticSourceFile) bool {
	return f.GetKind() == StaticSourceFile_SourceKind_WEAK
}

// Java: StaticSourceFile.isExtern() (default method)
func StaticSourceFile_IsExtern(f StaticSourceFile) bool {
	return f.GetKind() == StaticSourceFile_SourceKind_EXTERN
}

// Java: StaticSourceFile.isNonCode() (default method)
func StaticSourceFile_IsNonCode(f StaticSourceFile) bool {
	return f.GetKind() == StaticSourceFile_SourceKind_NON_CODE
}

// Java: StaticSourceFile.isTypeScriptSource() (default method)
func StaticSourceFile_IsTypeScriptSource(f StaticSourceFile) bool {
	sourceName := f.GetName()
	return strings.HasSuffix(sourceName, ".closure.js") || strings.HasSuffix(sourceName, ".tsx.cl.js")
}
