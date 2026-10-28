// Port of com/google/javascript/jscomp/RewriteCatchWithNoBinding.java @ closure-compiler v20261021

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
	"github.com/paullepoulpe/jsgo/internal/ast"
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
	compiler         AbstractCompiler
	astFactory       *AstFactory
	uniqueIdSupplier *UniqueIdSupplier
}

var rewriteCatchWithNoBindingTokens = func() ast.TokenSet {
	var s ast.TokenSet
	s.Add(rhino.Token_CATCH)
	return s
}()

// Java: RewriteCatchWithNoBinding.BINDING_NAME
const RewriteCatchWithNoBinding_BINDING_NAME = "$jscomp$unused$catch$"

// Java: RewriteCatchWithNoBinding(AbstractCompiler)
func NewRewriteCatchWithNoBinding(compiler AbstractCompiler) *RewriteCatchWithNoBinding {
	return &RewriteCatchWithNoBinding{
		compiler:         compiler,
		astFactory:       compiler.CreateAstFactory(),
		uniqueIdSupplier: compiler.GetUniqueIdSupplier(),
	}
}

func (*RewriteCatchWithNoBinding) HandledTokens() ast.TokenSet {
	return rewriteCatchWithNoBindingTokens
}

func (*RewriteCatchWithNoBinding) GetAdditionalFeaturesToRunOn() *parser.FeatureSet {
	return parser.FeatureSet_BARE_MINIMUM
}

// Java: RewriteCatchWithNoBinding.getTranspiledAwayFeatures()
func (r *RewriteCatchWithNoBinding) GetTranspiledAwayFeatures() *parser.FeatureSet {
	return parser.FeatureSet_BARE_MINIMUM.With(parser.FeatureSet_Feature_OPTIONAL_CATCH_BINDING)
}

// Java: RewriteCatchWithNoBinding.transpileSubtree(Node)
func (r *RewriteCatchWithNoBinding) TranspileSubtree(n rhino.Node) rhino.Node {
	catchNode, ok := ast.AsCatch(n)
	if !ok {
		return n
	}
	if _, hasParam := catchNode.Param(); hasParam {
		return n
	}

	placeholder := n.GetFirstChild()
	name := r.astFactory.CreateNameWithUnknownType(
		RewriteCatchWithNoBinding_BINDING_NAME +
			r.uniqueIdSupplier.GetUniqueId(r.compiler.GetInput(NodeUtil_GetInputId(n)))).
		SrcrefTree(placeholder)
	catchNode.SetParam(ast.MustName(name).AsCatchParam())
	r.compiler.ReportChangeToEnclosingScope(name)
	return n
}
