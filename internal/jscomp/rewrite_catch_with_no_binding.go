// Port of com/google/javascript/jscomp/RewriteCatchWithNoBinding.java @ closure-compiler v20261015

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

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// RewriteCatchWithNoBinding transpiles catch statements with no bindings by
// adding an unused binding.
//
//	try {
//	  stuff();
//	} catch {
//	  onError();
//	}
//
// Becomes
//
//	try {
//	  stuff();
//	} catch ($jscomp$unused$catch) {
//	  onError();
//	}
type RewriteCatchWithNoBinding struct {
	AbstractPeepholeTranspilationBase
	compiler         AbstractCompiler
	astFactory       *AstFactory
	uniqueIdSupplier *UniqueIdSupplier
}

// Java: RewriteCatchWithNoBinding.BINDING_NAME
const RewriteCatchWithNoBinding_BINDING_NAME = "$jscomp$unused$catch$"

// Java: RewriteCatchWithNoBinding(AbstractCompiler)
func NewRewriteCatchWithNoBinding(compiler AbstractCompiler) *RewriteCatchWithNoBinding {
	r := &RewriteCatchWithNoBinding{}
	r.compiler = compiler
	r.astFactory = compiler.CreateAstFactory()
	r.uniqueIdSupplier = compiler.GetUniqueIdSupplier()
	return r
}

// Java: RewriteCatchWithNoBinding.getTranspiledAwayFeatures()
func (r *RewriteCatchWithNoBinding) GetTranspiledAwayFeatures() *parser.FeatureSet {
	return parser.FeatureSet_BARE_MINIMUM.With(parser.FeatureSet_Feature_OPTIONAL_CATCH_BINDING)
}

// Java: RewriteCatchWithNoBinding.transpileSubtree(Node)
func (r *RewriteCatchWithNoBinding) TranspileSubtree(n *rhino.Node) *rhino.Node {
	if !n.IsCatch() || !n.GetFirstChild().IsEmpty() {
		return n
	}

	name := r.astFactory.CreateNameWithUnknownType(
		RewriteCatchWithNoBinding_BINDING_NAME +
			r.uniqueIdSupplier.GetUniqueId(r.compiler.GetInput(NodeUtil_GetInputId(n))))
	n.GetFirstChild().ReplaceWith(name.SrcrefTree(n.GetFirstChild()))
	r.compiler.ReportChangeToEnclosingScope(name)
	return n
}
