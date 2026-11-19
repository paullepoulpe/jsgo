/*
 * Copyright 2019 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/parsing/parser/StringLiteralToken.java (generated mechanically by j2go).

package parser

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.StringLiteralToken
type StringLiteralToken struct {
	idHash_ atomic.Int32
}

func (this_ *StringLiteralToken) BindSelf_(s any) {
	_ = s
}
func (this_ *StringLiteralToken) IsJavaNil_() bool { return this_ == nil }
func (*StringLiteralToken) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.parsing.parser.StringLiteralToken", "StringLiteralToken" }
func (this_ *StringLiteralToken) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*StringLiteralToken) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.StringLiteralToken", "com.google.javascript.jscomp.parsing.parser.StringLiteralToken", func(v_ any) bool {
		return j2gort.InstanceOf[*StringLiteralToken](v_)
	})
}

func (this_ *StringLiteralToken) initFields_StringLiteralToken() {
	_ = this_
}

// Java constructor: StringLiteralToken()
func NewStringLiteralToken() *StringLiteralToken {
	this_ := &StringLiteralToken{}
	InitStringLiteralToken(this_)
	return this_
}

func InitStringLiteralToken(this_ *StringLiteralToken) {
	this_.initFields_StringLiteralToken()
}

