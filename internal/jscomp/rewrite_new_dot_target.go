// Port of com/google/javascript/jscomp/RewriteNewDotTarget.java @ closure-compiler v20261021

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

// RewriteNewDotTarget transpiles away `new.target`.
type RewriteNewDotTarget struct {
	compiler   AbstractCompiler
	astFactory *AstFactory
}

var rewriteNewDotTargetTokens = func() ast.TokenSet {
	var s ast.TokenSet
	s.Add(rhino.Token_NEW_TARGET)
	return s
}()

// Java: RewriteNewDotTarget(AbstractCompiler)
func NewRewriteNewDotTarget(compiler AbstractCompiler) *RewriteNewDotTarget {
	return &RewriteNewDotTarget{
		compiler:   compiler,
		astFactory: compiler.CreateAstFactory(),
	}
}

func (*RewriteNewDotTarget) HandledTokens() ast.TokenSet {
	return rewriteNewDotTargetTokens
}

func (*RewriteNewDotTarget) GetAdditionalFeaturesToRunOn() *parser.FeatureSet {
	return parser.FeatureSet_BARE_MINIMUM
}

// Java: RewriteNewDotTarget.getTranspiledAwayFeatures()
func (r *RewriteNewDotTarget) GetTranspiledAwayFeatures() *parser.FeatureSet {
	return parser.FeatureSet_BARE_MINIMUM.With(parser.FeatureSet_Feature_NEW_TARGET)
}

// Java: RewriteNewDotTarget.transpileSubtree(Node)
func (r *RewriteNewDotTarget) TranspileSubtree(n rhino.Node) rhino.Node {
	if n.GetToken() != rhino.Token_NEW_TARGET {
		// there's nothing to rewrite
		return n
	}

	enclosingNonArrowFunction := NodeUtil_GetEnclosingNonArrowFunction(n)
	if !enclosingNonArrowFunction.IsNil() && NodeUtil_IsEs6Constructor(enclosingNonArrowFunction) {
		// Within an ES6 class constructor that we're about to transpile.
		// `new.target` -> `this.constructor`
		enclosingClass := enclosingNonArrowFunction.GetParent().GetGrandparent()
		replacement := r.astFactory.
			CreateGetProp(
				r.astFactory.CreateThisForEs6Class(enclosingClass), "constructor", AstFactory_Type_(n)).
			SrcrefTree(n)
		n.ReplaceWith(replacement)
		r.compiler.ReportChangeToEnclosingScope(replacement)
		return replacement
	}
	// Getting new.target correct in functions other than transpiled ES6 class constructors
	// requires determining whether the function was called with `new` or not, which is more
	// hassle than its worth. There's no good reason to use `new.target` in such places
	// anyway.
	TranspilationUtil_CannotConvertYet(r.compiler, n, "new.target")
	return n
}
