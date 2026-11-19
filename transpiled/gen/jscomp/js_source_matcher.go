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
// Port of com/google/javascript/refactoring/JsSourceMatcher.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.refactoring.JsSourceMatcher
type JsSourceMatcher struct {
	matcher *TemplateAstMatcher `gson:"matcher"`
	idHash_ atomic.Int32
}

func (this_ *JsSourceMatcher) BindSelf_(s any) {
	_ = s
}
func (this_ *JsSourceMatcher) IsJavaNil_() bool { return this_ == nil }
func (*JsSourceMatcher) JavaClassInfo_() (string, string) { return "com.google.javascript.refactoring.JsSourceMatcher", "JsSourceMatcher" }
func (this_ *JsSourceMatcher) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*JsSourceMatcher) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.refactoring.JsSourceMatcher", "com.google.javascript.refactoring.JsSourceMatcher", func(v_ any) bool {
		return j2gort.InstanceOf[*JsSourceMatcher](v_)
	})
}
func (*JsSourceMatcher) Implements_jscomp_Matcher() {}

func (this_ *JsSourceMatcher) initFields_JsSourceMatcher() {
	_ = this_
}

// Java constructor: JsSourceMatcher(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.rhino.Node templateNode, com.google.javascript.jscomp.TypeMatchingStrategy typeMatchingStrategy)
func NewJsSourceMatcher(compiler AbstractCompiler, templateNode pkg_rhino.Node, typeMatchingStrategy TypeMatchingStrategy) *JsSourceMatcher {
	this_ := &JsSourceMatcher{}
	InitJsSourceMatcher(this_, compiler, templateNode, typeMatchingStrategy)
	return this_
}

func InitJsSourceMatcher(this_ *JsSourceMatcher, compiler AbstractCompiler, templateNode pkg_rhino.Node, typeMatchingStrategy TypeMatchingStrategy) {
	_ = compiler
	_ = templateNode
	_ = typeMatchingStrategy
	this_.initFields_JsSourceMatcher()
	this_.matcher = NewTemplateAstMatcher(compiler, templateNode, typeMatchingStrategy)
}

// Java method: JsSourceMatcher.matches(com.google.javascript.rhino.Node n, com.google.javascript.refactoring.NodeMetadata metadata)
func (this_ *JsSourceMatcher) Matches(n pkg_rhino.Node, metadata *NodeMetadata) bool {
	_ = n
	_ = metadata
	return this_.matcher.Matches(n)
}

// Java method: JsSourceMatcher.getTemplateNodeToMatchMap()
func (this_ *JsSourceMatcher) GetTemplateNodeToMatchMap() j2gort.Ext_java_util_Map {
	return this_.matcher.GetTemplateNodeToMatchMap()
}

