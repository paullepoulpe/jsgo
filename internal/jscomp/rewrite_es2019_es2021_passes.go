// Port of com/google/javascript/jscomp/RewriteCatchWithNoBinding.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/RewriteLogicalAssignmentOperatorsPass.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/RewriteNewDotTarget.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/RewriteNullishCoalesceOperator.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/TranspilationUtil.java @ closure-compiler v20261021

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

/*
 * Copyright 2021 The Closure Compiler Authors.
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
 * Copyright 2020 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
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

// RewriteLogicalAssignmentOperatorsPass replaces the ES2020 `||=`, `&&=`,
// and `??=` operators.
type RewriteLogicalAssignmentOperatorsPass struct {
	compiler                                AbstractCompiler
	rewriteLogicalAssignmentOperatorsHelper *RewriteLogicalAssignmentOperatorsHelper
}

// Java: RewriteLogicalAssignmentOperatorsPass(AbstractCompiler)
func NewRewriteLogicalAssignmentOperatorsPass(compiler AbstractCompiler) *RewriteLogicalAssignmentOperatorsPass {
	r := &RewriteLogicalAssignmentOperatorsPass{}
	r.compiler = compiler
	r.rewriteLogicalAssignmentOperatorsHelper =
		NewRewriteLogicalAssignmentOperatorsHelper(
			compiler, compiler.CreateAstFactory(), compiler.GetUniqueIdSupplier())
	return r
}

// Java: RewriteLogicalAssignmentOperatorsPass.process(Node, Node)
func (r *RewriteLogicalAssignmentOperatorsPass) Process(externs rhino.Node, root rhino.Node) {
	NodeTraversal_Traverse(r.compiler, root, r)
	TranspilationPasses_MaybeMarkFeatureAsTranspiledAway(
		r.compiler, root, parser.FeatureSet_Feature_LOGICAL_ASSIGNMENT)
}

// Java: RewriteLogicalAssignmentOperatorsPass.shouldTraverse(NodeTraversal, Node, Node)
func (r *RewriteLogicalAssignmentOperatorsPass) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	if n.IsScript() {
		scriptFeatures := NodeUtil_GetFeatureSetOfScript(n)
		// runs only if logical assignment present
		return scriptFeatures == nil || scriptFeatures.ContainsFeature(parser.FeatureSet_Feature_LOGICAL_ASSIGNMENT)
	}
	return true
}

// Java: RewriteLogicalAssignmentOperatorsPass.visit(NodeTraversal, Node, Node)
func (r *RewriteLogicalAssignmentOperatorsPass) Visit(t *NodeTraversal, logicalAssignment rhino.Node, parent rhino.Node) {
	switch logicalAssignment.GetToken() {
	case rhino.Token_ASSIGN_OR, rhino.Token_ASSIGN_AND, rhino.Token_ASSIGN_COALESCE:
		r.rewriteLogicalAssignmentOperatorsHelper.VisitLogicalAssignmentOperator(
			t, logicalAssignment)
	default:
	}
}

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

// RewriteNullishCoalesceOperator replaces the ES2020 `??` operator with
// conditional (?:).
type RewriteNullishCoalesceOperator struct {
	compiler         AbstractCompiler
	astFactory       *AstFactory
	uniqueIdSupplier *UniqueIdSupplier
}

// Java: RewriteNullishCoalesceOperator.TEMP_VAR_NAME_PREFIX
const RewriteNullishCoalesceOperator_TEMP_VAR_NAME_PREFIX = "$jscomp$nullish$tmp"

// Java: RewriteNullishCoalesceOperator(AbstractCompiler)
func NewRewriteNullishCoalesceOperator(compiler AbstractCompiler) *RewriteNullishCoalesceOperator {
	r := &RewriteNullishCoalesceOperator{}
	r.compiler = compiler
	r.astFactory = compiler.CreateAstFactory()
	r.uniqueIdSupplier = compiler.GetUniqueIdSupplier()
	return r
}

// Java: RewriteNullishCoalesceOperator.process(Node, Node)
func (r *RewriteNullishCoalesceOperator) Process(externs rhino.Node, root rhino.Node) {
	NodeTraversal_Traverse(r.compiler, root, r)
	TranspilationPasses_MaybeMarkFeatureAsTranspiledAway(r.compiler, root, parser.FeatureSet_Feature_NULL_COALESCE_OP)
}

// Java: RewriteNullishCoalesceOperator.shouldTraverse(NodeTraversal, Node, Node)
func (r *RewriteNullishCoalesceOperator) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	if n.IsScript() {
		scriptFeatures := NodeUtil_GetFeatureSetOfScript(n)
		return scriptFeatures == nil || scriptFeatures.ContainsFeature(parser.FeatureSet_Feature_NULL_COALESCE_OP)
	}
	return true
}

// Java: RewriteNullishCoalesceOperator.visit(NodeTraversal, Node, Node)
func (r *RewriteNullishCoalesceOperator) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if n.IsNullishCoalesce() {
		r.VisitNullishCoalesce(t, n)
		t.ReportCodeChange()
	}
}

// Java: RewriteNullishCoalesceOperator.visitNullishCoalesce(NodeTraversal, Node)
func (r *RewriteNullishCoalesceOperator) VisitNullishCoalesce(t *NodeTraversal, n rhino.Node) {
	// a() ?? b()
	// let temp;
	// ((temp = a()) != null) ? temp : b()
	tempVarName := RewriteNullishCoalesceOperator_TEMP_VAR_NAME_PREFIX + r.uniqueIdSupplier.GetUniqueId(t.GetInput())
	enclosingStatement := NodeUtil_GetEnclosingStatement(n)

	left := n.RemoveFirstChild()
	right := n.GetLastChild().Detach()

	let := r.astFactory.CreateSingleLetNameDeclaration(tempVarName)
	assignName := r.astFactory.CreateName(tempVarName, AstFactory_Type_(left))
	assign := r.astFactory.CreateAssign(assignName, left)
	ne := r.astFactory.CreateNe(assign, r.astFactory.CreateNull())
	hookName := r.astFactory.CreateName(tempVarName, AstFactory_Type_(left))
	hook := r.astFactory.CreateHook(ne, hookName, right)

	let.SrcrefTreeIfMissing(left)
	assignName.SrcrefTreeIfMissing(left)
	assign.SrcrefTreeIfMissing(left)
	ne.SrcrefTreeIfMissing(left)
	hookName.SrcrefTreeIfMissing(left)

	let.InsertBefore(enclosingStatement)
	n.ReplaceWith(hook)

	NodeUtil_AddFeatureToScript(t.GetCurrentScript(), parser.FeatureSet_Feature_LET_DECLARATIONS, r.compiler)
	r.compiler.ReportChangeToEnclosingScope(hook)
}

// TranspilationUtil holds util functions for transpilation passes.
//
// Java: TranspilationUtil() is private (prevents instantiation); the class
// has only static members.
type TranspilationUtil struct{}

// Java: TranspilationUtil.CANNOT_CONVERT
var TranspilationUtil_CANNOT_CONVERT = DiagnosticType_Error("JSC_CANNOT_CONVERT", "This code cannot be transpiled. {0}")

// TODO(tbreisacher): Remove this once we have implemented transpilation for all the features
// we intend to support.
//
// Java: TranspilationUtil.CANNOT_CONVERT_YET
var TranspilationUtil_CANNOT_CONVERT_YET = DiagnosticType_Error(
	"JSC_CANNOT_CONVERT_YET", "Transpilation of ''{0}'' is not yet implemented.")

// Java: TranspilationUtil.cannotConvert(AbstractCompiler, Node, String)
func TranspilationUtil_CannotConvert(compiler AbstractCompiler, n rhino.Node, message string) {
	compiler.Report(JSError_MakeWithNode(n, TranspilationUtil_CANNOT_CONVERT, message))
}

// TranspilationUtil_CannotConvertYet warns the user that the given feature
// cannot be transpiled because the transpilation is not yet implemented. A
// call to this method is essentially a "TODO(tbreisacher): Implement
// feature" comment.
//
// Java: TranspilationUtil.cannotConvertYet(AbstractCompiler, Node, String)
func TranspilationUtil_CannotConvertYet(compiler AbstractCompiler, n rhino.Node, feature string) {
	compiler.Report(JSError_MakeWithNode(n, TranspilationUtil_CANNOT_CONVERT_YET, feature))
}

// Java: TranspilationUtil.preloadTranspilationRuntimeFunction(AbstractCompiler, String)
func TranspilationUtil_PreloadTranspilationRuntimeFunction(compiler AbstractCompiler, function string) {
	compiler.
		GetRuntimeJsLibManager().
		EnsureLibraryInjected("es6/util/"+jutil.StringToLowerCaseRoot(function), false)
}
