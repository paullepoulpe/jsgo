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
// Port of com/google/javascript/jscomp/Es6RewriteRestParameters.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.Es6RewriteRestParameters
type Es6RewriteRestParameters struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler         AbstractCompiler               `gson:"compiler"`
	astFactory       *AstFactory                    `gson:"astFactory"`
	namespace        pkg_rhino.StaticScope          `gson:"namespace"`
	getRestArguments RuntimeJsLibManager_JsLibField `gson:"getRestArguments"`
}

func (this_ *Es6RewriteRestParameters) BindSelf_(s any) {
	_ = s
}
func (this_ *Es6RewriteRestParameters) IsJavaNil_() bool { return this_ == nil }
func (*Es6RewriteRestParameters) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Es6RewriteRestParameters", "Es6RewriteRestParameters"
}
func (this_ *Es6RewriteRestParameters) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*Es6RewriteRestParameters) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Es6RewriteRestParameters", "com.google.javascript.jscomp.Es6RewriteRestParameters", func(v_ any) bool {
		return j2gort.InstanceOf[*Es6RewriteRestParameters](v_)
	})
}
func (*Es6RewriteRestParameters) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {}
func (*Es6RewriteRestParameters) Implements_jscomp_CompilerPass()                            {}
func (*Es6RewriteRestParameters) Implements_jscomp_NodeTraversal_Callback()                  {}

var static_Es6RewriteRestParameters_TRANSPILED_FEATURES *pkg_parser.FeatureSet

func StaticGet_Es6RewriteRestParameters_TRANSPILED_FEATURES() *pkg_parser.FeatureSet {
	clinit_Es6RewriteRestParameters()
	return static_Es6RewriteRestParameters_TRANSPILED_FEATURES
}
func StaticSet_Es6RewriteRestParameters_TRANSPILED_FEATURES(v *pkg_parser.FeatureSet) {
	clinit_Es6RewriteRestParameters()
	static_Es6RewriteRestParameters_TRANSPILED_FEATURES = v
}

var clinitDone_Es6RewriteRestParameters uint32
var clinitState_Es6RewriteRestParameters j2gort.ClinitState

func clinit_Es6RewriteRestParameters() {
	if atomic.LoadUint32(&clinitDone_Es6RewriteRestParameters) != 0 {
		return
	}
	clinitSlow_Es6RewriteRestParameters()
}

func clinitSlow_Es6RewriteRestParameters() {
	if !j2gort.EnterClinit(&clinitState_Es6RewriteRestParameters, &clinitDone_Es6RewriteRestParameters) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_Es6RewriteRestParameters, &clinitDone_Es6RewriteRestParameters)
	static_Es6RewriteRestParameters_TRANSPILED_FEATURES = pkg_parser.StaticGet_FeatureSet_BARE_MINIMUM().With_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_REST_PARAMETERS())
}

func (this_ *Es6RewriteRestParameters) initFields_Es6RewriteRestParameters() {
	_ = this_
}

// Java constructor: Es6RewriteRestParameters(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewEs6RewriteRestParameters(compiler AbstractCompiler) *Es6RewriteRestParameters {
	clinit_Es6RewriteRestParameters()
	this_ := &Es6RewriteRestParameters{}
	InitEs6RewriteRestParameters(this_, compiler)
	return this_
}

func InitEs6RewriteRestParameters(this_ *Es6RewriteRestParameters, compiler AbstractCompiler) {
	_ = compiler
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_Es6RewriteRestParameters()
	this_.compiler = compiler
	this_.astFactory = compiler.CreateAstFactory()
	this_.namespace = compiler.GetTranspilationNamespace()
	this_.getRestArguments = compiler.GetRuntimeJsLibManager().GetJsLibField("$jscomp.getRestArguments")
}

// Java method: Es6RewriteRestParameters.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *Es6RewriteRestParameters) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	TranspilationPasses_ProcessTranspile(this_.compiler, root, StaticGet_Es6RewriteRestParameters_TRANSPILED_FEATURES(), []NodeTraversal_Callback{this_})
	TranspilationPasses_MaybeMarkFeaturesAsTranspiledAway(this_.compiler, root, StaticGet_Es6RewriteRestParameters_TRANSPILED_FEATURES())
}

// Java method: Es6RewriteRestParameters.visit(com.google.javascript.jscomp.NodeTraversal traversal, com.google.javascript.rhino.Node current, com.google.javascript.rhino.Node parent)
func (this_ *Es6RewriteRestParameters) Visit_3(traversal *NodeTraversal, current pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = traversal
	_ = current
	_ = parent
	if current.RawToken_() == uint8(125) {
		this_.priv_Es6RewriteRestParameters_visitRestParam_2768(traversal, current, parent)
	}
}

// Java method: Es6RewriteRestParameters.visitRestParam(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node restParam, com.google.javascript.rhino.Node paramList)
func (this_ *Es6RewriteRestParameters) priv_Es6RewriteRestParameters_visitRestParam_2768(t *NodeTraversal, restParam pkg_rhino.Node, paramList pkg_rhino.Node) {
	_ = t
	_ = restParam
	_ = paramList
	var functionBody pkg_rhino.Node = paramList.GetNext()
	_ = functionBody
	var restIndex int32 = paramList.GetIndexOfChild(restParam)
	_ = restIndex
	var nameNode pkg_rhino.Node = restParam.GetOnlyChild()
	_ = nameNode
	var paramName string = nameNode.GetString()
	_ = paramName
	restParam.Detach()
	if !(functionBody.HasChildren()) {
		t.ReportCodeChange_0()
		return
	}
	var let pkg_rhino.Node = this_.astFactory.CreateSingleLetNameDeclaration_2(paramName, this_.astFactory.CreateCall(this_.astFactory.CreateGetPropWithUnknownType(this_.astFactory.CreateQName_2_StaticScope_RuntimeJsLibManager_JsLibField(this_.namespace, this_.getRestArguments), "apply"), AstFactory_S_Type_1_Node(nameNode), []pkg_rhino.Node{this_.astFactory.CreateNumber(float64(restIndex)), this_.astFactory.CreateArgumentsReference()})).SrcrefTreeIfMissing(functionBody)
	_ = let
	var insertBeforePoint pkg_rhino.Node = NodeUtil_GetInsertionPointAfterAllInnerFunctionDeclarations(functionBody)
	_ = insertBeforePoint
	if insertBeforePoint != 0 {
		let.InsertBefore(insertBeforePoint)
	} else {
		functionBody.AddChildToBack(let)
	}
	NodeUtil_AddFeatureToScript(t.GetCurrentScript(), pkg_parser.StaticGet_FeatureSet_Feature_LET_DECLARATIONS(), this_.compiler)
	t.ReportCodeChange_0()
}
func (*Es6RewriteRestParameters) AlwaysShouldTraverse_()      {}
func (*Es6RewriteRestParameters) CanReleaseSyntacticScopes_() {}
