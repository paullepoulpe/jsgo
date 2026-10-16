// Port of com/google/javascript/refactoring/CodeReplacement.java (with the AutoValue-generated
// AutoValue_CodeReplacement.java) @ closure-compiler v20260415

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

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// CodeReplacement is a class that represents a set of changes to make to
// the code.
//
// Go-only: language-adapter — an @AutoValue: the abstract class and the
// generated AutoValue_CodeReplacement (fields, equals, hashCode, toString)
// are one struct (PORTING.md §4).
type CodeReplacement struct {
	startPosition int
	length        int
	newContent    string
	sortKey       string
}

// Java: CodeReplacement.create(int, int, String)
func CodeReplacement_Create(startPosition int, length int, newContent string) *CodeReplacement {
	return CodeReplacement_CreateWithSortKey(startPosition, length, newContent, "")
}

// Java: CodeReplacement.create(int, int, String, String)
func CodeReplacement_CreateWithSortKey(startPosition int, length int, newContent string, sortKey string) *CodeReplacement {
	return newAutoValue_CodeReplacement(startPosition, length, newContent, sortKey)
}

// Java: AutoValue_CodeReplacement(int, int, String, String)
func newAutoValue_CodeReplacement(startPosition int, length int, newContent string, sortKey string) *CodeReplacement {
	return &CodeReplacement{
		startPosition: startPosition,
		length:        length,
		newContent:    newContent,
		sortKey:       sortKey,
	}
}

// GetStartPosition returns the start position within the file that the
// modification should be applied starting at.
//
// Java: CodeReplacement.getStartPosition()
func (r *CodeReplacement) GetStartPosition() int {
	return r.startPosition
}

// GetLength returns how many characters the new content should replace in
// the original content.
//
// Java: CodeReplacement.getLength()
func (r *CodeReplacement) GetLength() int {
	return r.length
}

// GetEndPosition returns the end position within the file that the
// modification should be applied starting at.
//
// Java: CodeReplacement.getEndPosition()
func (r *CodeReplacement) GetEndPosition() int {
	return r.GetStartPosition() + r.GetLength()
}

// GetNewContent returns the new content that should be inserted into the
// file.
//
// Java: CodeReplacement.getNewContent()
func (r *CodeReplacement) GetNewContent() string {
	return r.newContent
}

// GetSortKey returns an additional String key that can be used to sort
// replacements using lexical ordering.
//
// Java: CodeReplacement.getSortKey()
func (r *CodeReplacement) GetSortKey() string {
	return r.sortKey
}

// Java: CodeReplacement.compareTo(CodeReplacement)
func (r *CodeReplacement) CompareTo(x *CodeReplacement) int {
	return codeReplacement_APPLICATION_ORDER(r, x)
}

// Remember to compare every field so that it's consistent with equals.
//
// Java: CodeReplacement.APPLICATION_ORDER
func codeReplacement_APPLICATION_ORDER(a, b *CodeReplacement) int {
	if c := codeReplacement_compareInt(a.GetStartPosition(), b.GetStartPosition()); c != 0 {
		return c
	}
	if c := codeReplacement_compareInt(a.GetLength(), b.GetLength()); c != 0 {
		return c
	}
	if c := jutil.StringCompareTo(a.GetSortKey(), b.GetSortKey()); c != 0 {
		return c
	}
	return jutil.StringCompareTo(a.GetNewContent(), b.GetNewContent())
}

// codeReplacement_compareInt is Integer.compare(int, int).
func codeReplacement_compareInt(x, y int) int {
	if x < y {
		return -1
	} else if x == y {
		return 0
	}
	return 1
}

// Java: AutoValue_CodeReplacement.toString()
func (r *CodeReplacement) String() string {
	return "CodeReplacement{" +
		"startPosition=" + jutil.ToJavaString(r.startPosition) + ", " +
		"length=" + jutil.ToJavaString(r.length) + ", " +
		"newContent=" + r.newContent + ", " +
		"sortKey=" + r.sortKey +
		"}"
}

// Java: AutoValue_CodeReplacement.equals(Object)
func (r *CodeReplacement) Equals(o any) bool {
	that, ok := o.(*CodeReplacement)
	if !ok || that == nil {
		return false
	}
	if that == r {
		return true
	}
	return r.startPosition == that.GetStartPosition() &&
		r.length == that.GetLength() &&
		r.newContent == that.GetNewContent() &&
		r.sortKey == that.GetSortKey()
}

// Java: AutoValue_CodeReplacement.hashCode()
func (r *CodeReplacement) HashCode() int32 {
	var h int32 = 1
	h *= 1000003
	h ^= int32(r.startPosition)
	h *= 1000003
	h ^= int32(r.length)
	h *= 1000003
	h ^= jutil.StringHashCode(r.newContent)
	h *= 1000003
	h ^= jutil.StringHashCode(r.sortKey)
	return h
}
