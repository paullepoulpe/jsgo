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
// Port of com/google/javascript/jscomp/parsing/parser/SourceFile.java (generated mechanically by j2go).

package parser

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.SourceFile
type SourceFile struct {
	name string `gson:"name"`
	contents string `gson:"contents"`
	idHash_ atomic.Int32
}

func (this_ *SourceFile) BindSelf_(s any) {
	_ = s
}
func (this_ *SourceFile) IsJavaNil_() bool { return this_ == nil }
func (*SourceFile) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.parsing.parser.SourceFile", "SourceFile" }
func (this_ *SourceFile) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*SourceFile) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.SourceFile", "com.google.javascript.jscomp.parsing.parser.SourceFile", func(v_ any) bool {
		return j2gort.InstanceOf[*SourceFile](v_)
	})
}
func (this_ *SourceFile) FieldGet_name() string { return this_.name }
func (this_ *SourceFile) FieldGet_contents() string { return this_.contents }

func (this_ *SourceFile) initFields_SourceFile() {
	_ = this_
	this_.name = j2gort.NullString
	this_.contents = j2gort.NullString
}

// Java constructor: SourceFile(java.lang.String name, java.lang.String contents)
func NewSourceFile(name string, contents string) *SourceFile {
	this_ := &SourceFile{}
	InitSourceFile(this_, name, contents)
	return this_
}

func InitSourceFile(this_ *SourceFile, name string, contents string) {
	_ = name
	_ = contents
	this_.initFields_SourceFile()
	this_.name = name
	this_.contents = contents
}

