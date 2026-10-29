// Port of com/google/javascript/jscomp/Es6RewriteBlockScopedFunctionDeclaration.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/Es6RewriteRestParameters.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/SubstituteEs6Syntax.java @ closure-compiler v20261021

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

/*
 * Copyright 2015 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Es6RewriteBlockScopedFunctionDeclaration rewrites block-scoped function
// declarations as "let"s. This pass must happen before
// Es6RewriteBlockScopedDeclaration, which rewrites "let" to "var".
type Es6RewriteBlockScopedFunctionDeclaration struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler AbstractCompiler
}

// Java: Es6RewriteBlockScopedFunctionDeclaration.transpiledFeatures
var es6RewriteBlockScopedFunctionDeclaration_transpiledFeatures = parser.FeatureSet_BARE_MINIMUM.With(parser.FeatureSet_Feature_BLOCK_SCOPED_FUNCTION_DECLARATION)

// Java: Es6RewriteBlockScopedFunctionDeclaration(AbstractCompiler)
func NewEs6RewriteBlockScopedFunctionDeclaration(compiler AbstractCompiler) *Es6RewriteBlockScopedFunctionDeclaration {
	e := &Es6RewriteBlockScopedFunctionDeclaration{}
	e.compiler = compiler
	return e
}

// Java: Es6RewriteBlockScopedFunctionDeclaration.process(Node, Node)
func (e *Es6RewriteBlockScopedFunctionDeclaration) Process(externs rhino.Node, root rhino.Node) {
	TranspilationPasses_ProcessTranspile(e.compiler, root, es6RewriteBlockScopedFunctionDeclaration_transpiledFeatures, e)
	TranspilationPasses_MaybeMarkFeaturesAsTranspiledAway(e.compiler, root, es6RewriteBlockScopedFunctionDeclaration_transpiledFeatures)
}

// Java: Es6RewriteBlockScopedFunctionDeclaration.visit(NodeTraversal, Node, Node)
func (e *Es6RewriteBlockScopedFunctionDeclaration) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if n.IsFunction() &&
		!parent.IsNil() &&
		parent.IsBlock() &&
		!parent.GetParent().IsFunction() {
		// Only consider declarations (all expressions have non-block parents) that are not directly
		// within a function or top-level.
		e.VisitBlockScopedFunctionDeclaration(t, n, parent)
	}
}

// VisitBlockScopedFunctionDeclaration rewrites the function declaration
// from:
//
//	function f() {}
//	FUNCTION
//	  NAME x
//	  PARAM_LIST
//	  BLOCK
//
// to
//
//	let f = function() {};
//	LET
//	  NAME f
//	    FUNCTION
//	      NAME (w/ empty string)
//	      PARAM_LIST
//	      BLOCK
//
// This is similar to Normalize.NormalizeStatements#rewriteFunctionDeclaration
// but rewrites to "let" instead of "var".
//
// Java: Es6RewriteBlockScopedFunctionDeclaration.visitBlockScopedFunctionDeclaration(NodeTraversal, Node, Node)
func (e *Es6RewriteBlockScopedFunctionDeclaration) VisitBlockScopedFunctionDeclaration(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	// Prepare a spot for the function.
	oldNameNode := n.GetFirstChild()
	fnNameNode := oldNameNode.CloneNode()
	let := rhino.IR_Declaration(fnNameNode, rhino.Token_LET).Srcref(n)
	NodeUtil_AddFeatureToScript(t.GetCurrentScript(), parser.FeatureSet_Feature_LET_DECLARATIONS, e.compiler)

	// Prepare the function.
	oldNameNode.SetString("")
	e.compiler.ReportChangeToEnclosingScope(oldNameNode)

	// Move the function to the front of the parent.
	n.Detach()
	parent.AddChildToFront(let)
	e.compiler.ReportChangeToEnclosingScope(let)
	fnNameNode.AddChildToFront(n)
}

// Es6RewriteRestParameters converts REST parameters.
type Es6RewriteRestParameters struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	compiler         AbstractCompiler
	astFactory       *AstFactory
	namespace        rhino.StaticScope
	getRestArguments RuntimeJsLibManager_JsLibField
}

// Java: Es6RewriteRestParameters.TRANSPILED_FEATURES
var Es6RewriteRestParameters_TRANSPILED_FEATURES = parser.FeatureSet_BARE_MINIMUM.With(parser.FeatureSet_Feature_REST_PARAMETERS)

// Java: Es6RewriteRestParameters(AbstractCompiler)
func NewEs6RewriteRestParameters(compiler AbstractCompiler) *Es6RewriteRestParameters {
	e := &Es6RewriteRestParameters{}
	e.compiler = compiler
	e.astFactory = compiler.CreateAstFactory()
	e.namespace = compiler.GetTranspilationNamespace()
	e.getRestArguments =
		compiler.GetRuntimeJsLibManager().GetJsLibField("$jscomp.getRestArguments")
	return e
}

// Java: Es6RewriteRestParameters.process(Node, Node)
func (e *Es6RewriteRestParameters) Process(externs rhino.Node, root rhino.Node) {
	TranspilationPasses_ProcessTranspile(e.compiler, root, Es6RewriteRestParameters_TRANSPILED_FEATURES, e)
	TranspilationPasses_MaybeMarkFeaturesAsTranspiledAway(e.compiler, root, Es6RewriteRestParameters_TRANSPILED_FEATURES)
}

// Java: Es6RewriteRestParameters.visit(NodeTraversal, Node, Node)
func (e *Es6RewriteRestParameters) Visit(traversal *NodeTraversal, current rhino.Node, parent rhino.Node) {
	if current.GetToken() == rhino.Token_ITER_REST {
		e.VisitRestParam(traversal, current, parent)
	}
}

// VisitRestParam processes a rest parameter
//
// Java: Es6RewriteRestParameters.visitRestParam(NodeTraversal, Node, Node)
func (e *Es6RewriteRestParameters) VisitRestParam(t *NodeTraversal, restParam rhino.Node, paramList rhino.Node) {
	functionBody := paramList.GetNext()
	restIndex := paramList.GetIndexOfChild(restParam)
	nameNode := restParam.GetOnlyChild()
	paramName := nameNode.GetString()

	// Remove the existing param from the list, as it will be replaced with a declaration with the
	// same name.
	restParam.Detach()

	if !functionBody.HasChildren() {
		// If function has no body, we are done!
		t.ReportCodeChange()
		return
	}

	// Now that the restParam is deleted, create a let declaration by making a new NAME node of the
	// same name `paramName`
	let :=
		e.astFactory.
			CreateSingleLetNameDeclarationWithValue(
				paramName, // creates a new NAME node with name `paramName`
				e.astFactory.CreateCall(
					e.astFactory.CreateGetPropWithUnknownType(
						e.astFactory.CreateQNameOfField(e.namespace, e.getRestArguments), "apply"),
					AstFactory_Type_(nameNode),
					e.astFactory.CreateNumber(float64(restIndex)),
					e.astFactory.CreateArgumentsReference())).
			SrcrefTreeIfMissing(functionBody)
	insertBeforePoint :=
		NodeUtil_GetInsertionPointAfterAllInnerFunctionDeclarations(functionBody)
	if !insertBeforePoint.IsNil() {
		let.InsertBefore(insertBeforePoint)
	} else {
		// functionBody only contains hoisted function declarations
		functionBody.AddChildToBack(let)
	}
	NodeUtil_AddFeatureToScript(t.GetCurrentScript(), parser.FeatureSet_Feature_LET_DECLARATIONS, e.compiler)
	t.ReportCodeChange()
}

// SubstituteEs6Syntax is an optimization that does peephole optimizations
// of ES6 code.
type SubstituteEs6Syntax struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	compiler                       AbstractCompiler
	objectLiteralShorthandWasAdded bool
}

// Java: SubstituteEs6Syntax(AbstractCompiler)
func NewSubstituteEs6Syntax(compiler AbstractCompiler) *SubstituteEs6Syntax {
	s := &SubstituteEs6Syntax{}
	s.compiler = compiler
	return s
}

// Java: SubstituteEs6Syntax.process(Node, Node)
func (s *SubstituteEs6Syntax) Process(externs rhino.Node, root rhino.Node) {
	NodeTraversal_Traverse(s.compiler, root, s)
}

// Java: SubstituteEs6Syntax.visit(NodeTraversal, Node, Node)
func (s *SubstituteEs6Syntax) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_FUNCTION:
		if n.IsArrowFunction() {
			s.MaybeSimplifyArrowFunctionBody(n, n.GetLastChild())
		}
	case rhino.Token_STRING_KEY:
		if n.GetFirstChild().IsName() && n.GetFirstChild().GetStringID() == n.GetStringID() {
			n.SetShorthandProperty(true)
			s.objectLiteralShorthandWasAdded = true
		}
	case rhino.Token_SCRIPT:
		if s.objectLiteralShorthandWasAdded {
			NodeUtil_AddFeatureToScript(n, parser.FeatureSet_Feature_SHORTHAND_OBJECT_PROPERTIES, s.compiler)
			s.objectLiteralShorthandWasAdded = false // false for the next script
		}
	default:
	}
}

// MaybeSimplifyArrowFunctionBody: if possible, replace functions of the
// form ()=>{ return x; } with ()=>x
//
// Java: SubstituteEs6Syntax.maybeSimplifyArrowFunctionBody(Node, Node)
func (s *SubstituteEs6Syntax) MaybeSimplifyArrowFunctionBody(arrowFunction rhino.Node, body rhino.Node) {
	jutil.CheckArgument(arrowFunction.IsArrowFunction())
	if !body.IsBlock() || !body.HasOneChild() || !body.GetFirstChild().IsReturn() {
		return
	}
	returnValue := body.GetFirstChild().RemoveFirstChild()
	var replacement rhino.Node
	if !returnValue.IsNil() {
		replacement = returnValue
	} else {
		replacement = rhino.IR_Name("undefined")
	}
	body.ReplaceWith(replacement)
	s.compiler.ReportChangeToEnclosingScope(replacement)
}
