// Port of com/google/javascript/jscomp/PeepholeTranspilationsPass.java @ closure-compiler v20260415

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
	"reflect"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// PeepholeTranspilationsPass is a compiler pass to run various peephole
// transpilations (e.g. rewriteCatchBindings, rewriteNewDotTarget,
// reportUntranspilableFeatures, es6NormalizeShorthandProperties, etc).
type PeepholeTranspilationsPass struct {
	compiler AbstractCompiler
	// NOTE: Use a native array rather than a List to avoid creating iterators for every node in the
	// AST.
	peepholeTranspilations []AbstractPeepholeTranspilation

	// The feature set to mark as transpiled away after all peephole transpilations are done.
	featureSetToMarkAsTranspiledAway *parser.FeatureSet

	// The feature set that triggers this pass. That is, this pass will only run on scripts that
	// contain at least one of these features.
	featureSetToRunOn *parser.FeatureSet
}

// Java: PeepholeTranspilationsPass(AbstractCompiler, List<AbstractPeepholeTranspilation>, FeatureSet, FeatureSet)
func newPeepholeTranspilationsPass(
	compiler AbstractCompiler,
	transpilations []AbstractPeepholeTranspilation,
	featureSetToRunOn *parser.FeatureSet,
	featureSetToMarkAsTranspiledAway *parser.FeatureSet) *PeepholeTranspilationsPass {
	p := &PeepholeTranspilationsPass{}
	p.compiler = compiler
	p.peepholeTranspilations = append([]AbstractPeepholeTranspilation(nil), transpilations...)
	p.featureSetToRunOn = featureSetToRunOn
	p.featureSetToMarkAsTranspiledAway = featureSetToMarkAsTranspiledAway
	return p
}

// PeepholeTranspilationsPass_Create creates a peephole optimization pass that
// runs the given optimizations.
//
// Java: PeepholeTranspilationsPass.create(AbstractCompiler, List<AbstractPeepholeTranspilation>)
func PeepholeTranspilationsPass_Create(
	compiler AbstractCompiler, transpilations []AbstractPeepholeTranspilation) *PeepholeTranspilationsPass {
	featureSetToRunOn := parser.FeatureSet_BARE_MINIMUM
	featureSetToMarkAsTranspiledAway := parser.FeatureSet_BARE_MINIMUM
	// initialize the feature sets to the union of all the features that the peephole transpilations
	// need to run on and the features that they transpile away.
	for _, transpilation := range transpilations {
		jutil.CheckStateMsg1(
			!transpilation.
				GetAdditionalFeaturesToRunOn().
				ContainsAtLeastOneOf(transpilation.GetTranspiledAwayFeatures()),
			"Transpilation pass %s has getAdditionalFeaturesToRunOn() that contains features that it"+
				" transpiles away.",
			peepholeTranspilationsPass_simpleName(transpilation))
		featureSetToRunOn = featureSetToRunOn.WithFeatureSet(transpilation.GetTranspiledAwayFeatures())
		// some passes need to run on additional features beyond those they transpile away.
		featureSetToRunOn = featureSetToRunOn.WithFeatureSet(transpilation.GetAdditionalFeaturesToRunOn())
		featureSetToMarkAsTranspiledAway =
			featureSetToMarkAsTranspiledAway.WithFeatureSet(transpilation.GetTranspiledAwayFeatures())
	}
	return newPeepholeTranspilationsPass(
		compiler, transpilations, featureSetToRunOn, featureSetToMarkAsTranspiledAway)
}

// peepholeTranspilationsPass_simpleName is transpilation.getClass().getSimpleName().
//
// Go-only: language-adapter — the Go struct of each transpilation is named
// after its Java class.
func peepholeTranspilationsPass_simpleName(transpilation AbstractPeepholeTranspilation) string {
	t := reflect.TypeOf(transpilation)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	return t.Name()
}

// Java: PeepholeTranspilationsPass.process(Node, Node)
func (p *PeepholeTranspilationsPass) Process(externs *rhino.Node, root *rhino.Node) {
	NodeTraversal_Traverse(p.compiler, root, &PeepholeTranspilationsPass_PeepCallback{outer: p})
	// Update the featureSets with all features that got transpiled.
	TranspilationPasses_MaybeMarkFeaturesAsTranspiledAway(
		p.compiler, root, p.featureSetToMarkAsTranspiledAway)
}

// PeepholeTranspilationsPass_PeepCallback is the traversal callback of the
// pass.
type PeepholeTranspilationsPass_PeepCallback struct {
	outer *PeepholeTranspilationsPass
}

// Java: PeepholeTranspilationsPass.PeepCallback.visit(NodeTraversal, Node, Node)
func (c *PeepholeTranspilationsPass_PeepCallback) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {

	currentNode := n
	for _, transpilation := range c.outer.peepholeTranspilations {
		currentNode = transpilation.TranspileSubtree(currentNode)
		if currentNode == nil {
			// this subtree was removed by the current transpilation, so we don't need to run the rest
			// of the peephole transpilations.
			return
		}
	}
}

// Java: PeepholeTranspilationsPass.PeepCallback.shouldTraverse(NodeTraversal, Node, Node)
func (c *PeepholeTranspilationsPass_PeepCallback) ShouldTraverse(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_SCRIPT:
		// check if the script contains any of the features that we are transpiling.
		scriptFeatures := NodeUtil_GetFeatureSetOfScript(n)
		return scriptFeatures.ContainsAtLeastOneOf(c.outer.featureSetToRunOn)
	default:
		return true
	}
}
