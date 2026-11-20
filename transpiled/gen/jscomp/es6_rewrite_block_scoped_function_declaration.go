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
// Port of com/google/javascript/jscomp/Es6RewriteBlockScopedFunctionDeclaration.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.Es6RewriteBlockScopedFunctionDeclaration
type Es6RewriteBlockScopedFunctionDeclaration struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler AbstractCompiler `gson:"compiler"`
}

func (this_ *Es6RewriteBlockScopedFunctionDeclaration) BindSelf_(s any) {
	_ = s
}
func (this_ *Es6RewriteBlockScopedFunctionDeclaration) IsJavaNil_() bool { return this_ == nil }
func (*Es6RewriteBlockScopedFunctionDeclaration) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Es6RewriteBlockScopedFunctionDeclaration", "Es6RewriteBlockScopedFunctionDeclaration"
}
func (this_ *Es6RewriteBlockScopedFunctionDeclaration) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*Es6RewriteBlockScopedFunctionDeclaration) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Es6RewriteBlockScopedFunctionDeclaration", "com.google.javascript.jscomp.Es6RewriteBlockScopedFunctionDeclaration", func(v_ any) bool {
		return j2gort.InstanceOf[*Es6RewriteBlockScopedFunctionDeclaration](v_)
	})
}
func (*Es6RewriteBlockScopedFunctionDeclaration) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {
}
func (*Es6RewriteBlockScopedFunctionDeclaration) Implements_jscomp_CompilerPass()           {}
func (*Es6RewriteBlockScopedFunctionDeclaration) Implements_jscomp_NodeTraversal_Callback() {}

var static_Es6RewriteBlockScopedFunctionDeclaration_transpiledFeatures *pkg_parser.FeatureSet

func StaticGet_Es6RewriteBlockScopedFunctionDeclaration_TranspiledFeatures() *pkg_parser.FeatureSet {
	clinit_Es6RewriteBlockScopedFunctionDeclaration()
	return static_Es6RewriteBlockScopedFunctionDeclaration_transpiledFeatures
}
func StaticSet_Es6RewriteBlockScopedFunctionDeclaration_TranspiledFeatures(v *pkg_parser.FeatureSet) {
	clinit_Es6RewriteBlockScopedFunctionDeclaration()
	static_Es6RewriteBlockScopedFunctionDeclaration_transpiledFeatures = v
}

var clinitDone_Es6RewriteBlockScopedFunctionDeclaration uint32
var clinitState_Es6RewriteBlockScopedFunctionDeclaration j2gort.ClinitState

func clinit_Es6RewriteBlockScopedFunctionDeclaration() {
	if atomic.LoadUint32(&clinitDone_Es6RewriteBlockScopedFunctionDeclaration) != 0 {
		return
	}
	clinitSlow_Es6RewriteBlockScopedFunctionDeclaration()
}

func clinitSlow_Es6RewriteBlockScopedFunctionDeclaration() {
	if !j2gort.EnterClinit(&clinitState_Es6RewriteBlockScopedFunctionDeclaration, &clinitDone_Es6RewriteBlockScopedFunctionDeclaration) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_Es6RewriteBlockScopedFunctionDeclaration, &clinitDone_Es6RewriteBlockScopedFunctionDeclaration)
	static_Es6RewriteBlockScopedFunctionDeclaration_transpiledFeatures = pkg_parser.StaticGet_FeatureSet_BARE_MINIMUM().With_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_BLOCK_SCOPED_FUNCTION_DECLARATION())
}

func (this_ *Es6RewriteBlockScopedFunctionDeclaration) initFields_Es6RewriteBlockScopedFunctionDeclaration() {
	_ = this_
}

// Java constructor: Es6RewriteBlockScopedFunctionDeclaration(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewEs6RewriteBlockScopedFunctionDeclaration(compiler AbstractCompiler) *Es6RewriteBlockScopedFunctionDeclaration {
	clinit_Es6RewriteBlockScopedFunctionDeclaration()
	this_ := &Es6RewriteBlockScopedFunctionDeclaration{}
	InitEs6RewriteBlockScopedFunctionDeclaration(this_, compiler)
	return this_
}

func InitEs6RewriteBlockScopedFunctionDeclaration(this_ *Es6RewriteBlockScopedFunctionDeclaration, compiler AbstractCompiler) {
	_ = compiler
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_Es6RewriteBlockScopedFunctionDeclaration()
	this_.compiler = compiler
}

// Java method: Es6RewriteBlockScopedFunctionDeclaration.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *Es6RewriteBlockScopedFunctionDeclaration) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	TranspilationPasses_ProcessTranspile(this_.compiler, root, StaticGet_Es6RewriteBlockScopedFunctionDeclaration_TranspiledFeatures(), []NodeTraversal_Callback{this_})
	TranspilationPasses_MaybeMarkFeaturesAsTranspiledAway(this_.compiler, root, StaticGet_Es6RewriteBlockScopedFunctionDeclaration_TranspiledFeatures())
}

// Java method: Es6RewriteBlockScopedFunctionDeclaration.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *Es6RewriteBlockScopedFunctionDeclaration) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if ((n.IsFunction() && (parent != 0)) && parent.IsBlock()) && !(parent.GetParent().IsFunction()) {
		this_.priv_Es6RewriteBlockScopedFunctionDeclaration_visitBlockScopedFunctionDeclaration_2605(t, n, parent)
	}
}

// Java method: Es6RewriteBlockScopedFunctionDeclaration.visitBlockScopedFunctionDeclaration(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *Es6RewriteBlockScopedFunctionDeclaration) priv_Es6RewriteBlockScopedFunctionDeclaration_visitBlockScopedFunctionDeclaration_2605(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	var oldNameNode pkg_rhino.Node = n.GetFirstChild()
	_ = oldNameNode
	var fnNameNode pkg_rhino.Node = oldNameNode.CloneNode_0()
	_ = fnNameNode
	var let pkg_rhino.Node = pkg_rhino.IR_Declaration_2(fnNameNode, pkg_rhino.StaticGet_Token_LET()).Srcref(n)
	_ = let
	NodeUtil_AddFeatureToScript(t.GetCurrentScript(), pkg_parser.StaticGet_FeatureSet_Feature_LET_DECLARATIONS(), this_.compiler)
	oldNameNode.SetString("")
	this_.compiler.ReportChangeToEnclosingScope(oldNameNode)
	n.Detach()
	parent.AddChildToFront(let)
	this_.compiler.ReportChangeToEnclosingScope(let)
	fnNameNode.AddChildToFront(n)
}
func (*Es6RewriteBlockScopedFunctionDeclaration) AlwaysShouldTraverse_()      {}
func (*Es6RewriteBlockScopedFunctionDeclaration) CanReleaseSyntacticScopes_() {}
