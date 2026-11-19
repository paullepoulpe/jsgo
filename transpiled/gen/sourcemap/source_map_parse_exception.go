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
// Port of com/google/debugging/sourcemap/SourceMapParseException.java (generated mechanically by j2go).

package sourcemap

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.debugging.sourcemap.SourceMapParseException
type SourceMapParseException struct {
	javalang.JavaException
	idHash_ atomic.Int32
}

func (this_ *SourceMapParseException) BindSelf_(s any) {
	_ = s
}
func (this_ *SourceMapParseException) IsJavaNil_() bool { return this_ == nil }
func (*SourceMapParseException) JavaClassInfo_() (string, string) { return "com.google.debugging.sourcemap.SourceMapParseException", "SourceMapParseException" }
func (this_ *SourceMapParseException) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.debugging.sourcemap.SourceMapParseException", "com.google.debugging.sourcemap.SourceMapParseException", func(v_ any) bool {
		return j2gort.InstanceOf[*SourceMapParseException](v_)
	})
}

func (this_ *SourceMapParseException) initFields_SourceMapParseException() {
	_ = this_
}

// Java constructor: SourceMapParseException(java.lang.String message)
func NewSourceMapParseException(message string) *SourceMapParseException {
	this_ := &SourceMapParseException{}
	InitSourceMapParseException(this_, message)
	return this_
}

func InitSourceMapParseException(this_ *SourceMapParseException, message string) {
	_ = message
	j2gort.InitUserThrowable(this_, "com.google.debugging.sourcemap.SourceMapParseException", j2gort.BoxString(message))
	this_.initFields_SourceMapParseException()
}

