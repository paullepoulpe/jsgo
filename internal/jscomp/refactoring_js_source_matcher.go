// Port of com/google/javascript/refactoring/JsSourceMatcher.java @ closure-compiler v20261021

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

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// JsSourceMatcher is a Matcher that can take arbitrary JS source code and use
// it as a template to find matches in other source code.
//
// Java: public final class JsSourceMatcher implements Matcher
type JsSourceMatcher struct {
	matcher *TemplateAstMatcher
}

// NewJsSourceMatcher constructs this matcher with a Function node that serves
// as the template to match all other nodes against. The body of the function
// will be used to match against.
//
// Java: JsSourceMatcher(AbstractCompiler, Node, TypeMatchingStrategy)
func NewJsSourceMatcher(
	compiler AbstractCompiler, templateNode rhino.Node, typeMatchingStrategy TypeMatchingStrategy) *JsSourceMatcher {
	m := &JsSourceMatcher{}
	m.matcher = NewTemplateAstMatcher(compiler, templateNode, typeMatchingStrategy)
	return m
}

// Java: JsSourceMatcher.matches(Node, NodeMetadata)
func (m *JsSourceMatcher) Matches(n rhino.Node, metadata *NodeMetadata) bool {
	return m.matcher.Matches(n)
}

// GetTemplateNodeToMatchMap returns a map from named template node strings to
// Nodes that were the equivalent matches from the last matched template.
//
// Java: JsSourceMatcher.getTemplateNodeToMatchMap()
func (m *JsSourceMatcher) GetTemplateNodeToMatchMap() *jutil.LinkedHashMap[string, rhino.Node] {
	return m.matcher.GetTemplateNodeToMatchMap()
}
