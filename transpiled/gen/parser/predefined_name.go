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
// Port of com/google/javascript/jscomp/parsing/parser/PredefinedName.java (generated mechanically by j2go).

package parser

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.PredefinedName
type PredefinedName struct {
	idHash_ atomic.Int32
}

func (this_ *PredefinedName) BindSelf_(s any) {
	_ = s
}
func (this_ *PredefinedName) IsJavaNil_() bool { return this_ == nil }
func (*PredefinedName) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.parsing.parser.PredefinedName", "PredefinedName" }
func (this_ *PredefinedName) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*PredefinedName) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.PredefinedName", "com.google.javascript.jscomp.parsing.parser.PredefinedName", func(v_ any) bool {
		return j2gort.InstanceOf[*PredefinedName](v_)
	})
}

func StaticGet_PredefinedName_AS() string { return "as" }
func StaticGet_PredefinedName_FROM() string { return "from" }
func StaticGet_PredefinedName_GET() string { return "get" }
func StaticGet_PredefinedName_OF() string { return "of" }
func StaticGet_PredefinedName_SET() string { return "set" }
func (this_ *PredefinedName) initFields_PredefinedName() {
	_ = this_
}

// Java constructor: PredefinedName()
func NewPredefinedName() *PredefinedName {
	this_ := &PredefinedName{}
	InitPredefinedName(this_)
	return this_
}

func InitPredefinedName(this_ *PredefinedName) {
	this_.initFields_PredefinedName()
}

