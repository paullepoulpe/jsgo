// Port of com/google/javascript/jscomp/AbstractPeepholeOptimization.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/AbstractPeepholeTranspilation.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/OptimizeLetAndConstPeephole.java @ closure-compiler v20261021

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

/*
 * Copyright 2026 The Closure Compiler Authors.
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
	"math/big"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// AbstractPeepholeOptimization is an abstract class whose implementations
// run peephole optimizations: optimizations that look at a small section of
// code and either remove that code (if it is not needed) or replaces it with
// smaller code.
//
// # How to write a subclass (PORTING.md §4)
//
// Every peephole pass (PeepholeFoldConstants, PeepholeRemoveDeadCode,
// StatementFusion, the J2cl rewriters, ...) is a struct that embeds
// AbstractPeepholeOptimizationBase (by value) and implements OptimizeSubtree:
//
//	type PeepholeFoo struct {
//		AbstractPeepholeOptimizationBase
//		late bool // the subclass's own fields
//	}
//
//	// Java: PeepholeFoo(boolean)
//	func NewPeepholeFoo(late bool) *PeepholeFoo {
//		return &PeepholeFoo{late: late}
//	}
//
//	// Java: PeepholeFoo.optimizeSubtree(Node)
//	func (p *PeepholeFoo) OptimizeSubtree(subtree rhino.Node) rhino.Node { ... }
//
// The subclass calls the protected helpers on itself (p.MayHaveSideEffects(n),
// p.ReportChangeToEnclosingScope(n), p.GetSideEffectFreeBooleanValue(n), ...);
// they are promoted from the embedded base. A subclass that overrides
// beginTraversal (the J2cl rewriters do) defines its own BeginTraversal and
// calls p.AbstractPeepholeOptimizationBase.BeginTraversal(compiler) for
// super.beginTraversal(compiler). Callers that hold the Java type
// AbstractPeepholeOptimization (PeepholeOptimizationsPass, DefaultPassConfig)
// use this interface, so the override is the one that runs.
//
// The compiler and the AstAnalyzer are deliberately private to the base, as
// in Java ("Intentionally not exposed to subclasses").
type AbstractPeepholeOptimization interface {
	// OptimizeSubtree, given a node to optimize and a traversal, optimizes
	// the node. Subclasses should override to provide their own peephole
	// optimization.
	//
	// subtree is the subtree that will be optimized. Returns the new version
	// of the subtree (or null if the subtree or one of its parents was
	// removed from the AST). If the subtree has not changed, this method
	// must return subtree.
	//
	// Java: AbstractPeepholeOptimization.optimizeSubtree(Node)
	OptimizeSubtree(subtree rhino.Node) rhino.Node

	// Java: AbstractPeepholeOptimization.beginTraversal(AbstractCompiler)
	BeginTraversal(compiler AbstractCompiler)

	// Java: AbstractPeepholeOptimization.endTraversal()
	EndTraversal()

	// GetAbstractPeepholeOptimizationBase gives access to the base fields
	// through the interface (PORTING.md §4).
	GetAbstractPeepholeOptimizationBase() *AbstractPeepholeOptimizationBase

	peepholeConfig() peepholeOptConfig
	cloneForWorker() AbstractPeepholeOptimization
}

type peepholeOptConfig struct {
	tokens                ast.TokenSet
	isStandard            bool
	tracksFeatures        bool
	isJ2clCallPass        bool
	isOptimizeLetAndConst bool
}

func peepholeTokenSet(tokens ...rhino.Token) ast.TokenSet {
	var s ast.TokenSet
	for _, t := range tokens {
		s.Add(t)
	}
	return s
}

func peepholeTokenSetWith(base ast.TokenSet, tokens ...rhino.Token) ast.TokenSet {
	for _, t := range tokens {
		base.Add(t)
	}
	return base
}

var peepholeCallOnlyTokens = peepholeTokenSet(rhino.Token_CALL)

// AbstractPeepholeOptimizationBase holds the fields and the base
// implementations of AbstractPeepholeOptimization; every subclass embeds it.
type AbstractPeepholeOptimizationBase struct {
	// Intentionally not exposed to subclasses
	compiler AbstractCompiler
	// Intentionally not exposed to subclasses
	astAnalyzer *AstAnalyzer

	// newFeatures are new parser features added in some OptimizeSubtree call.
	//
	// PeepholeOptimizationsPass uses this to call
	// NodeUtil#addFeaturesToScript in closer to O(1) time, as this class
	// doesn't know what script node it's currently in, and would need to
	// walk the AST.
	newFeatures *parser.FeatureSet
}

// GetAbstractPeepholeOptimizationBase returns the embedded base.
func (b *AbstractPeepholeOptimizationBase) GetAbstractPeepholeOptimizationBase() *AbstractPeepholeOptimizationBase {
	return b
}

func (b *AbstractPeepholeOptimizationBase) peepholeConfig() peepholeOptConfig {
	return peepholeOptConfig{tracksFeatures: true}
}

func (b *AbstractPeepholeOptimizationBase) cloneForWorker() AbstractPeepholeOptimization {
	return nil
}

// Report is a helper method for reporting an error to the compiler when
// applying a peephole optimization.
//
// diagnostic is the error type; n is the node for which the error should be
// reported.
//
// Java: AbstractPeepholeOptimization.report(DiagnosticType, Node)
func (b *AbstractPeepholeOptimizationBase) Report(diagnostic *DiagnosticType, n rhino.Node) {
	error_ := JSError_MakeWithNode(n, diagnostic, n.String())
	b.compiler.Report(error_)
}

// AreNodesEqualForInlining reports whether the nodes are equal for the
// purpose of inlining. If type aware optimizations are on, type equality is
// checked.
//
// Java: AbstractPeepholeOptimization.areNodesEqualForInlining(Node, Node)
func (b *AbstractPeepholeOptimizationBase) AreNodesEqualForInlining(n1 rhino.Node, n2 rhino.Node) bool {
	/* Our implementation delegates to the compiler. We provide this
	 * method because we don't want to expose Compiler to PeepholeOptimizations.
	 */
	jutil.CheckNotNull(b.compiler)
	return b.compiler.AreNodesEqualForInlining(n1, n2)
}

// IsASTNormalized reports whether the current AST is normalized (e.g. has
// the Normalize pass been run and has the Denormalize pass not yet been
// run?).
//
// Java: AbstractPeepholeOptimization.isASTNormalized()
func (b *AbstractPeepholeOptimizationBase) IsASTNormalized() bool {
	jutil.CheckNotNull(b.compiler)

	return b.compiler.GetLifeCycleStage().IsNormalized()
}

// BeginTraversal informs the optimization that a traversal will begin.
//
// Java: AbstractPeepholeOptimization.beginTraversal(AbstractCompiler)
func (b *AbstractPeepholeOptimizationBase) BeginTraversal(compiler AbstractCompiler) {
	b.compiler = jutil.CheckNotNull(compiler)
	b.astAnalyzer = compiler.GetAstAnalyzer()
	if b.newFeatures == nil {
		b.newFeatures = parser.FeatureSet_BARE_MINIMUM
	}
}

// EndTraversal informs the optimization that a traversal has ended.
//
// This class cannot be used after this method is called, unless
// BeginTraversal is called again.
//
// Java: AbstractPeepholeOptimization.endTraversal()
func (b *AbstractPeepholeOptimizationBase) EndTraversal() {
	jutil.CheckStateMsg1(
		b.newFeatures == nil || b.newFeatures.IsEmpty(),
		"Expected getNewFeatures() to be empty in endTraversal() but found %s",
		b.newFeatures)
	b.compiler = nil
	b.astAnalyzer = nil
}

// MayEffectMutableState returns whether the node may create new mutable
// state, or change existing state.
//
// Java: AbstractPeepholeOptimization.mayEffectMutableState(Node)
func (b *AbstractPeepholeOptimizationBase) MayEffectMutableState(n rhino.Node) bool {
	return b.astAnalyzer.MayEffectMutableState(n)
}

// MayHaveSideEffects returns whether the node may have side effects when
// executed.
//
// Java: AbstractPeepholeOptimization.mayHaveSideEffects(Node)
func (b *AbstractPeepholeOptimizationBase) MayHaveSideEffects(n rhino.Node) bool {
	return b.astAnalyzer.MayHaveSideEffects(n)
}

// GetSideEffectFreeNumberValue returns the number value of the node if it
// has one and it cannot have side effects.
//
// Java: AbstractPeepholeOptimization.getSideEffectFreeNumberValue(Node)
func (b *AbstractPeepholeOptimizationBase) GetSideEffectFreeNumberValue(n rhino.Node) (float64, bool) {
	// Calculating the number value, if any, is likely to be faster than calculating side effects,
	// and there are only a very few cases where we can compute a number value, but there could
	// also be side effects. e.g. `void doSomething()` has value NaN, regardless of the behavior
	// of `doSomething()`
	v, ok := nodeUtil_doGetNumberValue(n, true)
	if !ok || b.astAnalyzer.MayHaveSideEffects(n) {
		return 0, false
	}
	return v, true
}

// Java: AbstractPeepholeOptimization.getSideEffectFreeNumberValueNoConversion(Node)
func (b *AbstractPeepholeOptimizationBase) GetSideEffectFreeNumberValueNoConversion(n rhino.Node) (float64, bool) {
	v, ok := nodeUtil_doGetNumberValue(n, false)
	if !ok || b.astAnalyzer.MayHaveSideEffects(n) {
		return 0, false
	}
	return v, true
}

// GetSideEffectFreeBigIntValue returns the bigint value of the node if it
// has one and it cannot have side effects.
//
// Returns null otherwise.
//
// Java: AbstractPeepholeOptimization.getSideEffectFreeBigIntValue(Node)
func (b *AbstractPeepholeOptimizationBase) GetSideEffectFreeBigIntValue(n rhino.Node) *big.Int {
	value := NodeUtil_GetBigIntValue(n)
	// Calculating the bigint value, if any, is likely to be faster than calculating side effects,
	// and there are only a very few cases where we can compute a bigint value, but there could
	// also be side effects. e.g. `void doSomething()` has value NaN, regardless of the behavior
	// of `doSomething()`
	if value != nil && b.astAnalyzer.MayHaveSideEffects(n) {
		value = nil
	}
	return value
}

// GetSideEffectFreeStringValue gets the value of a node as a String if it can
// be converted and has no side effects.
//
// This method effectively emulates the String() JavaScript cast function
// when possible and the node has no side effects.
//
// Java: AbstractPeepholeOptimization.getSideEffectFreeStringValue(Node)
func (b *AbstractPeepholeOptimizationBase) GetSideEffectFreeStringValue(n rhino.Node) (string, bool) {
	value, ok := nodeUtil_getStringValue(n)
	// Calculating the string value, if any, is likely to be faster than calculating side effects,
	// and there are only a very few cases where we can compute a string value, but there could
	// also be side effects. e.g. `void doSomething()` has value 'undefined', regardless of the
	// behavior of `doSomething()`
	if !ok || b.astAnalyzer.MayHaveSideEffects(n) {
		return "", false
	}
	return value, true
}

// GetSideEffectFreeBooleanValue calculates the known boolean value for a
// node if possible and if it has no side effects.
//
// Returns Tri#UNKNOWN if the node has side effects or its value cannot be
// statically determined.
//
// Java: AbstractPeepholeOptimization.getSideEffectFreeBooleanValue(Node)
func (b *AbstractPeepholeOptimizationBase) GetSideEffectFreeBooleanValue(n rhino.Node) base.Tri {
	value := NodeUtil_GetBooleanValue(n)
	// Calculating the boolean value, if any, is likely to be faster than calculating side effects,
	// and there are only a very few cases where we can compute a boolean value, but there could
	// also be side effects. e.g. `void doSomething()` has value `false`, regardless of the
	// behavior of `doSomething()`
	if value != base.Tri_UNKNOWN && b.astAnalyzer.MayHaveSideEffects(n) {
		value = base.Tri_UNKNOWN
	}
	return value
}

// NodeTypeMayHaveSideEffects returns true if the current node's type implies
// side effects.
//
// This is a non-recursive version of the may have side effects check; used
// to check wherever the current node's type is one of the reason's why a
// subtree has side effects.
//
// Java: AbstractPeepholeOptimization.nodeTypeMayHaveSideEffects(Node)
func (b *AbstractPeepholeOptimizationBase) NodeTypeMayHaveSideEffects(n rhino.Node) bool {
	return b.astAnalyzer.NodeTypeMayHaveSideEffects(n)
}

// IsEcmaScript5OrGreater returns whether the output language is ECMAScript 5
// or later. Workarounds for quirks in browsers that do not support ES5 can
// be ignored when this is true.
//
// Java: AbstractPeepholeOptimization.isEcmaScript5OrGreater()
func (b *AbstractPeepholeOptimizationBase) IsEcmaScript5OrGreater() bool {
	return b.compiler != nil &&
		b.compiler.GetOptions().GetOutputFeatureSet().Contains(parser.FeatureSet_ES5)
}

// GetCodingConvention returns the current coding convention.
//
// Java: AbstractPeepholeOptimization.getCodingConvention()
func (b *AbstractPeepholeOptimizationBase) GetCodingConvention() CodingConvention {
	// Note: this assumes a thread safe coding convention object.
	return b.compiler.GetCodingConvention()
}

// Java: AbstractPeepholeOptimization.reportChangeToEnclosingScope(Node)
func (b *AbstractPeepholeOptimizationBase) ReportChangeToEnclosingScope(n rhino.Node) {
	b.compiler.ReportChangeToEnclosingScope(n)
}

// DeleteNode calls NodeUtil#deleteNode(Node, AbstractCompiler).
//
// Java: AbstractPeepholeOptimization.deleteNode(Node)
func (b *AbstractPeepholeOptimizationBase) DeleteNode(property rhino.Node) {
	jutil.CheckNotNull(b.compiler)
	NodeUtil_DeleteNode(property, b.compiler)
}

// MarkFunctionsDeleted calls NodeUtil#markFunctionsDeleted(Node,
// AbstractCompiler).
//
// Java: AbstractPeepholeOptimization.markFunctionsDeleted(Node)
func (b *AbstractPeepholeOptimizationBase) MarkFunctionsDeleted(function rhino.Node) {
	jutil.CheckNotNull(b.compiler)
	NodeUtil_MarkFunctionsDeleted(function, b.compiler)
}

// MarkNewScopesChanged calls NodeUtil#markNewScopesChanged(Node,
// AbstractCompiler).
//
// Java: AbstractPeepholeOptimization.markNewScopesChanged(Node)
func (b *AbstractPeepholeOptimizationBase) MarkNewScopesChanged(n rhino.Node) {
	jutil.CheckNotNull(b.compiler)
	NodeUtil_MarkNewScopesChanged(n, b.compiler)
}

// AddFeatureToEnclosingScript calls NodeUtil#addFeatureToScript with the
// script currently being visited.
//
// Java: AbstractPeepholeOptimization.addFeatureToEnclosingScript(Feature)
func (b *AbstractPeepholeOptimizationBase) AddFeatureToEnclosingScript(feature parser.FeatureSet_Feature) {
	// NOTE: we could implement this as
	//  protected final void addFeatureToEnclosingScript(Node node, Feature feature) {
	//     NodeUtil.addFeatureToScript(NodeUtil.getEnclosingScript(node), ...
	// This is expected to be behaviorally equivalent.
	// However, walking the AST to get the enclosing script is O(depth of the AST) per call. With
	// the  current implementation, we let PeepholeOptimizationsPass to find the script in what's
	// usually O(1) time, or worst case it only calls getEnclosingScript() once per NodeTraversal
	// and then caches the result.
	if b.newFeatures == nil {
		b.newFeatures = parser.FeatureSet_BARE_MINIMUM
	}
	b.newFeatures = b.newFeatures.With(feature)
}

// Java: AbstractPeepholeOptimization.clearNewFeatures()
func (b *AbstractPeepholeOptimizationBase) ClearNewFeatures() {
	b.newFeatures = parser.FeatureSet_BARE_MINIMUM
}

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
	// Most passes return FeatureSet_BARE_MINIMUM, except passes that need to
	// run on files containing features they don't transpile away (e.g.
	// ReportUntranspilableFeatures triggering on REGEXP_SYNTAX).
	//
	// Java: AbstractPeepholeTranspilation.getAdditionalFeaturesToRunOn()
	GetAdditionalFeaturesToRunOn() *parser.FeatureSet

	// HandledTokens returns the exact set of AST node tokens inspected by
	// TranspileSubtree.
	HandledTokens() ast.TokenSet

	// TranspileSubtree transpiles the given node. Subclasses should
	// override to do their own peephole rewriting.
	//
	// Returns the new version of the subtree (or null if the subtree was
	// removed from the AST). If the subtree has not changed, this method
	// must return subtree.
	//
	// Java: AbstractPeepholeTranspilation.transpileSubtree(Node)
	TranspileSubtree(subtree rhino.Node) rhino.Node
}

// OptimizeLetAndConstPeephole is a late-stage compiler pass that lowers
// `const`/`let` to `var` where safe, and `const` to `let`, improving keyword
// homogeneity for better compression.
type OptimizeLetAndConstPeephole struct {
	AbstractPeepholeOptimizationBase
	assumeOutputIsWrapped bool
}

// Java: OptimizeLetAndConstPeephole(boolean)
func NewOptimizeLetAndConstPeephole(assumeOutputIsWrapped bool) *OptimizeLetAndConstPeephole {
	return &OptimizeLetAndConstPeephole{assumeOutputIsWrapped: assumeOutputIsWrapped}
}

var optimizeLetAndConstTokens = peepholeTokenSet(
	rhino.Token_LET,
	rhino.Token_CONST,
)

func (p *OptimizeLetAndConstPeephole) peepholeConfig() peepholeOptConfig {
	return peepholeOptConfig{
		tokens:                optimizeLetAndConstTokens,
		isStandard:            true,
		tracksFeatures:        true,
		isOptimizeLetAndConst: true,
	}
}

func (p *OptimizeLetAndConstPeephole) cloneForWorker() AbstractPeepholeOptimization {
	return NewOptimizeLetAndConstPeephole(p.assumeOutputIsWrapped)
}

// Java: OptimizeLetAndConstPeephole.optimizeSubtree(Node)
func (p *OptimizeLetAndConstPeephole) OptimizeSubtree(subtree rhino.Node) rhino.Node {
	switch subtree.GetToken() {
	case rhino.Token_LET:
		return p.ProcessLet(subtree)
	case rhino.Token_CONST:
		return p.ProcessConst(subtree)
	default:
		return subtree
	}
}

// Java: OptimizeLetAndConstPeephole.processLet(Node)
func (p *OptimizeLetAndConstPeephole) ProcessLet(n rhino.Node) rhino.Node {
	if !p.assumeOutputIsWrapped && n.GetParent().IsScript() {
		return n
	}

	if p.IsDirectlyInHoistScope(n) {
		n.SetToken(rhino.Token_VAR)
		p.ReportChangeToEnclosingScope(n)
	}
	return n
}

// Java: OptimizeLetAndConstPeephole.processConst(Node)
func (p *OptimizeLetAndConstPeephole) ProcessConst(n rhino.Node) rhino.Node {
	if !p.assumeOutputIsWrapped && n.GetParent().IsScript() {
		return n
	}

	if p.IsDirectlyInHoistScope(n) {
		n.SetToken(rhino.Token_VAR)
	} else {
		n.SetToken(rhino.Token_LET)
		p.AddFeatureToEnclosingScript(parser.FeatureSet_Feature_LET_DECLARATIONS)
	}
	p.ReportChangeToEnclosingScope(n)
	return n
}

// Java: OptimizeLetAndConstPeephole.isDirectlyInHoistScope(Node)
func (p *OptimizeLetAndConstPeephole) IsDirectlyInHoistScope(n rhino.Node) bool {
	parent := n.GetParent()
	for !parent.IsNil() && parent.IsLabel() {
		parent = parent.GetParent()
	}
	if parent.IsNil() {
		return false
	}

	if parent.IsScript() || parent.IsModuleBody() {
		return true
	}

	if parent.IsBlock() {
		grandparent := parent.GetParent()
		if !grandparent.IsNil() && (grandparent.IsFunction() || grandparent.IsClassMembers()) {
			return true
		}
	}

	return false
}
