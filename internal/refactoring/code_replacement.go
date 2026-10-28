// Port of com/google/javascript/refactoring/CodeReplacement.java (with the AutoValue-generated
// AutoValue_CodeReplacement.java) @ closure-compiler v20261021

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

package refactoring

import (
	"cmp"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// CodeReplacement represents a set of changes to make to the code.
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
	return &CodeReplacement{
		startPosition: startPosition,
		length:        length,
		newContent:    newContent,
		sortKey:       sortKey,
	}
}

// Java: CodeReplacement.getStartPosition()
func (r *CodeReplacement) GetStartPosition() int { return r.startPosition }

// Java: CodeReplacement.getLength()
func (r *CodeReplacement) GetLength() int { return r.length }

// Java: CodeReplacement.getEndPosition()
func (r *CodeReplacement) GetEndPosition() int { return r.startPosition + r.length }

// Java: CodeReplacement.getNewContent()
func (r *CodeReplacement) GetNewContent() string { return r.newContent }

// Java: CodeReplacement.getSortKey()
func (r *CodeReplacement) GetSortKey() string { return r.sortKey }

// Java: CodeReplacement.compareTo(CodeReplacement)
func (r *CodeReplacement) CompareTo(x *CodeReplacement) int {
	if c := cmp.Compare(r.startPosition, x.startPosition); c != 0 {
		return c
	}
	if c := cmp.Compare(r.length, x.length); c != 0 {
		return c
	}
	if c := jutil.StringCompareTo(r.sortKey, x.sortKey); c != 0 {
		return c
	}
	return jutil.StringCompareTo(r.newContent, x.newContent)
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
	return r.startPosition == that.startPosition &&
		r.length == that.length &&
		r.newContent == that.newContent &&
		r.sortKey == that.sortKey
}

// Java: AutoValue_CodeReplacement.hashCode()
func (r *CodeReplacement) HashCode() int32 {
	var h int32 = 1
	h = h*1000003 ^ int32(r.startPosition)
	h = h*1000003 ^ int32(r.length)
	h = h*1000003 ^ jutil.StringHashCode(r.newContent)
	h = h*1000003 ^ jutil.StringHashCode(r.sortKey)
	return h
}
