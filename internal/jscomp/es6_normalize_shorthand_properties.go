// Port of com/google/javascript/jscomp/Es6NormalizeShorthandProperties.java @ closure-compiler v20261015

/*
 * Copyright 2017 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Es6NormalizeShorthandProperties normalizes shorthand object properties.
// This should be one of the first things done when transpiling from ES6 down
// to ES5, as it allows all the following checks and transpilations to not
// care about shorthand and destructured assignments.
type Es6NormalizeShorthandProperties struct {
	AbstractPeepholeTranspilationBase
}

// Java: Es6NormalizeShorthandProperties(AbstractCompiler)
func NewEs6NormalizeShorthandProperties(compiler AbstractCompiler) *Es6NormalizeShorthandProperties {
	return &Es6NormalizeShorthandProperties{}
}

// Java: Es6NormalizeShorthandProperties.transpileSubtree(Node)
func (e *Es6NormalizeShorthandProperties) TranspileSubtree(subtree rhino.Node) rhino.Node {
	if subtree.IsStringKey() {
		subtree.SetShorthandProperty(false)
	}
	return subtree
}

// Java: Es6NormalizeShorthandProperties.getTranspiledAwayFeatures()
func (e *Es6NormalizeShorthandProperties) GetTranspiledAwayFeatures() *parser.FeatureSet {
	return parser.FeatureSet_BARE_MINIMUM.With(parser.FeatureSet_Feature_SHORTHAND_OBJECT_PROPERTIES)
}
