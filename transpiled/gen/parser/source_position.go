/*
 * Copyright 2011 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/parsing/parser/util/SourcePosition.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.util.SourcePosition
type SourcePosition struct {
	offset int32 `gson:"offset"`
	line int32 `gson:"line"`
	column int32 `gson:"column"`
}

func (this_ SourcePosition) BindSelf_(s any) {
	_ = s
}
func (this_ SourcePosition) IsJavaNil_() bool { return false }
func (SourcePosition) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.parsing.parser.util.SourcePosition", "SourcePosition" }
func (SourcePosition) HasIdentityEquals_() {}
func (this_ SourcePosition) IdentityHash_() int32 {
	var h_ int32 = 1
	h_ = 31*h_ + (this_.offset)
	h_ = 31*h_ + (this_.line)
	h_ = 31*h_ + (this_.column)
	return h_
}
func (this_ SourcePosition) FastValHash_() (int32, bool, bool) { return this_.IdentityHash_(), true, true }
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.util.SourcePosition", "com.google.javascript.jscomp.parsing.parser.util.SourcePosition", func(v_ any) bool {
		return j2gort.InstanceOf[SourcePosition](v_)
	})
}
func (this_ SourcePosition) FieldGet_offset() int32 { return this_.offset }
func (this_ SourcePosition) FieldGet_line() int32 { return this_.line }
func (this_ SourcePosition) FieldGet_column() int32 { return this_.column }

func (this_ *SourcePosition) initFields_SourcePosition() {
	_ = this_
}

// Java constructor: SourcePosition(com.google.javascript.jscomp.parsing.parser.@org.jspecify.annotations.Nullable SourceFile source, int offset, int line, int column)
func NewSourcePosition(source *SourceFile, offset int32, line int32, column int32) SourcePosition {
	this_ := SourcePosition{}
	InitSourcePosition(&this_, source, offset, line, column)
	return this_
}

func InitSourcePosition(this_ *SourcePosition, source *SourceFile, offset int32, line int32, column int32) {
	_ = source
	_ = offset
	_ = line
	_ = column
	this_.initFields_SourcePosition()
	this_.offset = offset
	this_.line = line
	this_.column = column
}

// Java method: SourcePosition.getOffset()
func (this_ SourcePosition) GetOffset() int32 {
	return this_.offset
}

// Java method: SourcePosition.toString()
func (this_ SourcePosition) String() string {
	return ("(" + j2gort.StrValInt32((this_.line + int32(1))) + ", " + j2gort.StrValInt32((this_.column + int32(1))) + ")")
}

