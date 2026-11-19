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
// Port of com/google/javascript/jscomp/parsing/parser/IdentifierToken.java (generated mechanically by j2go).

package parser

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.IdentifierToken
type IdentifierToken struct {
	idHash_ atomic.Int32
}

func (this_ *IdentifierToken) BindSelf_(s any) {
	_ = s
}
func (this_ *IdentifierToken) IsJavaNil_() bool { return this_ == nil }
func (*IdentifierToken) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.parsing.parser.IdentifierToken", "IdentifierToken" }
func (this_ *IdentifierToken) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*IdentifierToken) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.IdentifierToken", "com.google.javascript.jscomp.parsing.parser.IdentifierToken", func(v_ any) bool {
		return j2gort.InstanceOf[*IdentifierToken](v_)
	})
}

func (this_ *IdentifierToken) initFields_IdentifierToken() {
	_ = this_
}

// Java constructor: IdentifierToken()
func NewIdentifierToken() *IdentifierToken {
	this_ := &IdentifierToken{}
	InitIdentifierToken(this_)
	return this_
}

func InitIdentifierToken(this_ *IdentifierToken) {
	this_.initFields_IdentifierToken()
}

