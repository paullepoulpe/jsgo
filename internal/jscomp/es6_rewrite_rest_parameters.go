// Port of com/google/javascript/jscomp/Es6RewriteRestParameters.java @ closure-compiler v20261021

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
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

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
