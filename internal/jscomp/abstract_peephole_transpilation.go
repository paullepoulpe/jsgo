// Port of com/google/javascript/jscomp/AbstractPeepholeTranspilation.java @ closure-compiler v20261015

/*
 * Copyright 2010 The Closure Compiler Authors.
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

// AbstractPeepholeTranspilation is an abstract class whose implementations
// run peephole transpilations: transpilations that look at a small section
// of code with TranspileSubtree and transpile it. The set of transpilations
// done must be small, independent and local rewritings that don't transform
// the AST structure outside of the provided subtree.
type AbstractPeepholeTranspilation interface {
	// GetTranspiledAwayFeatures represents the set of features that get
	// transpiled away by this pass.
	//
	// Java: AbstractPeepholeTranspilation.getTranspiledAwayFeatures()
	GetTranspiledAwayFeatures() *parser.FeatureSet

	// GetAdditionalFeaturesToRunOn is the set of additional features, if
	// any, beyond those transpiled away which if present in a SCRIPT should
	// trigger this pass.
	//
	// Java: AbstractPeepholeTranspilation.getAdditionalFeaturesToRunOn()
	GetAdditionalFeaturesToRunOn() *parser.FeatureSet

	// TranspileSubtree transpiles the given node. Subclasses should
	// override to do their own peephole rewriting.
	//
	// Returns the new version of the subtree (or null if the subtree was
	// removed from the AST). If the subtree has not changed, this method
	// must return subtree.
	//
	// Java: AbstractPeepholeTranspilation.transpileSubtree(Node)
	TranspileSubtree(subtree *rhino.Node) *rhino.Node
}

// AbstractPeepholeTranspilationBase holds the base implementations; every
// subclass embeds it.
type AbstractPeepholeTranspilationBase struct{}

// GetAdditionalFeaturesToRunOn is the set of additional features, if any,
// beyond those transpiled away which if present in a SCRIPT should trigger
// this pass.
//
// Most passes shouldn't override this, but it is available for the rare
// pass that needs to run on files containing features that it doesn't
// transpile away.
//
// For example, we trigger the peephole pass ReportUntranspilableFeatures
// pass for a SCRIPT which uses the ES3 feature Feature.REGEXP_SYNTAX, but
// we only transpile away certain unsupported RegExp flags (e.g.
// REGEXP_LOOKBEHIND) in that pass.
//
// Another example, the Es6NormalizeClasses, although not a peephole
// transpilation, must trigger on Feature.CLASS, but it does not transpile it
// away.
//
// Similarly, the pass Es6ConvertSuper, although not a peephole
// transpilation, runs on Feature.SUPER, but does not transpile it away as
// calls to `super()` are not transpiled by it.
//
// Java: AbstractPeepholeTranspilation.getAdditionalFeaturesToRunOn()
func (b *AbstractPeepholeTranspilationBase) GetAdditionalFeaturesToRunOn() *parser.FeatureSet {
	return parser.FeatureSet_BARE_MINIMUM
}
