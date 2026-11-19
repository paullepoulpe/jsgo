/*
 * Copyright 2009 The Closure Compiler Authors.
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
// Port of com/google/debugging/sourcemap/FilePosition.java (generated mechanically by j2go).

package sourcemap

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.debugging.sourcemap.FilePosition
type FilePosition struct {
	line    int32 `gson:"line"`
	column  int32 `gson:"column"`
	idHash_ atomic.Int32
}

func (this_ *FilePosition) BindSelf_(s any) {
	_ = s
}
func (this_ *FilePosition) IsJavaNil_() bool { return this_ == nil }
func (*FilePosition) JavaClassInfo_() (string, string) {
	return "com.google.debugging.sourcemap.FilePosition", "FilePosition"
}
func (this_ *FilePosition) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*FilePosition) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.debugging.sourcemap.FilePosition", "com.google.debugging.sourcemap.FilePosition", func(v_ any) bool {
		return j2gort.InstanceOf[*FilePosition](v_)
	})
}

func (this_ *FilePosition) initFields_FilePosition() {
	_ = this_
}

// Java constructor: FilePosition(int line, int column)
func NewFilePosition(line int32, column int32) *FilePosition {
	this_ := &FilePosition{}
	InitFilePosition(this_, line, column)
	return this_
}

func InitFilePosition(this_ *FilePosition, line int32, column int32) {
	_ = line
	_ = column
	this_.initFields_FilePosition()
	this_.line = line
	this_.column = column
}

// Java method: FilePosition.getLine()
func (this_ *FilePosition) GetLine() int32 {
	return this_.line
}

// Java method: FilePosition.getColumn()
func (this_ *FilePosition) GetColumn() int32 {
	return this_.column
}

// Java method: FilePosition.toString()
func (this_ *FilePosition) String() string {
	return j2gort.Ext_com_google_common_base_MoreObjects_S_ToStringHelper_1_any(this_).Add_2_string_int32("line", this_.line).Add_2_string_int32("column", this_.column).String()
}
