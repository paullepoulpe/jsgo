// Port of com/google/javascript/jscomp/RewritePrivateClassProperties.java @ closure-compiler v20260415

/*
 * Copyright 2025 The Closure Compiler Authors.
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

// RewritePrivateClassProperties transpiles away usages of private properties
// in ES6 classes.
type RewritePrivateClassProperties struct {
	AbstractPeepholeTranspilationBase
	compiler AbstractCompiler
}

// Java: RewritePrivateClassProperties(AbstractCompiler)
func NewRewritePrivateClassProperties(compiler AbstractCompiler) *RewritePrivateClassProperties {
	r := &RewritePrivateClassProperties{}
	r.compiler = compiler
	return r
}

// Java: RewritePrivateClassProperties.getTranspiledAwayFeatures()
func (r *RewritePrivateClassProperties) GetTranspiledAwayFeatures() *parser.FeatureSet {
	return parser.FeatureSet_BARE_MINIMUM.With(parser.FeatureSet_Feature_PRIVATE_CLASS_PROPERTIES)
}

// Java: RewritePrivateClassProperties.transpileSubtree(Node)
func (r *RewritePrivateClassProperties) TranspileSubtree(n *rhino.Node) *rhino.Node {
	if n.IsPrivateIdentifier() {
		TranspilationUtil_CannotConvertYet(r.compiler, n, "private class properties")
	}
	return n
}
