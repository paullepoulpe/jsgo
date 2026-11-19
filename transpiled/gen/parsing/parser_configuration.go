/*
 * Copyright 2020 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/parsing/ParserConfiguration.java (generated mechanically by j2go).

package parsing

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.ParserConfiguration
type ParserConfiguration struct {
	idHash_ atomic.Int32
}

func (this_ *ParserConfiguration) BindSelf_(s any) {
	_ = s
}
func (this_ *ParserConfiguration) IsJavaNil_() bool { return this_ == nil }
func (*ParserConfiguration) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.parsing.ParserConfiguration", "ParserConfiguration" }
func (this_ *ParserConfiguration) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ParserConfiguration) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.ParserConfiguration", "com.google.javascript.jscomp.parsing.ParserConfiguration", func(v_ any) bool {
		return j2gort.InstanceOf[*ParserConfiguration](v_)
	})
}

func (this_ *ParserConfiguration) initFields_ParserConfiguration() {
	_ = this_
}

// Java constructor: ParserConfiguration()
func NewParserConfiguration() *ParserConfiguration {
	this_ := &ParserConfiguration{}
	InitParserConfiguration(this_)
	return this_
}

func InitParserConfiguration(this_ *ParserConfiguration) {
	this_.initFields_ParserConfiguration()
}

// Java static method: ParserConfiguration.getString(java.lang.String name)
func ParserConfiguration_GetString(name string) string {
	_ = name
	return j2gort.Ext_java_util_ResourceBundle_GetBundle_1("com.google.javascript.jscomp.parsing.ParserConfig").GetString(name)
}

