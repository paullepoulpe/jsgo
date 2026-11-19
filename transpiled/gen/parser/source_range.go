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
// Port of com/google/javascript/jscomp/parsing/parser/util/SourceRange.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.util.SourceRange
type SourceRange struct {
	start SourcePosition `gson:"start"`
	end SourcePosition `gson:"end"`
}

func (this_ SourceRange) BindSelf_(s any) {
	_ = s
}
func (this_ SourceRange) IsJavaNil_() bool { return false }
func (SourceRange) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.parsing.parser.util.SourceRange", "SourceRange" }
func (SourceRange) HasIdentityEquals_() {}
func (this_ SourceRange) IdentityHash_() int32 {
	var h_ int32 = 1
	h_ = 31*h_ + (this_.start).IdentityHash_()
	h_ = 31*h_ + (this_.end).IdentityHash_()
	return h_
}
func (this_ SourceRange) FastValHash_() (int32, bool, bool) { return this_.IdentityHash_(), true, true }
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.util.SourceRange", "com.google.javascript.jscomp.parsing.parser.util.SourceRange", func(v_ any) bool {
		return j2gort.InstanceOf[SourceRange](v_)
	})
}
func (this_ SourceRange) FieldGet_start() SourcePosition { return this_.start }
func (this_ SourceRange) FieldGet_end() SourcePosition { return this_.end }

func (this_ *SourceRange) initFields_SourceRange() {
	_ = this_
}

// Java constructor: SourceRange(com.google.javascript.jscomp.parsing.parser.util.SourcePosition start, com.google.javascript.jscomp.parsing.parser.util.SourcePosition end)
func NewSourceRange(start SourcePosition, end SourcePosition) SourceRange {
	this_ := SourceRange{}
	InitSourceRange(&this_, start, end)
	return this_
}

func InitSourceRange(this_ *SourceRange, start SourcePosition, end SourcePosition) {
	_ = start
	_ = end
	this_.initFields_SourceRange()
	this_.start = start
	this_.end = end
}

// Java method: SourceRange.toString()
func (this_ SourceRange) String() string {
	return j2gort.Ext_java_lang_String_Format_2("<%s - %s>", []any{this_.start, this_.end})
}

