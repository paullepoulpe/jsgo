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
// Port of com/google/javascript/jscomp/SynthesizeExplicitConstructors.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.SynthesizeExplicitConstructors
type SynthesizeExplicitConstructors struct {
	compiler   AbstractCompiler `gson:"compiler"`
	astFactory *AstFactory      `gson:"astFactory"`
	idHash_    atomic.Int32
}

func (this_ *SynthesizeExplicitConstructors) BindSelf_(s any) {
	_ = s
}
func (this_ *SynthesizeExplicitConstructors) IsJavaNil_() bool { return this_ == nil }
func (*SynthesizeExplicitConstructors) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.SynthesizeExplicitConstructors", "SynthesizeExplicitConstructors"
}
func (this_ *SynthesizeExplicitConstructors) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*SynthesizeExplicitConstructors) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.SynthesizeExplicitConstructors", "com.google.javascript.jscomp.SynthesizeExplicitConstructors", func(v_ any) bool {
		return j2gort.InstanceOf[*SynthesizeExplicitConstructors](v_)
	})
}

func (this_ *SynthesizeExplicitConstructors) initFields_SynthesizeExplicitConstructors() {
	_ = this_
}

// Java constructor: SynthesizeExplicitConstructors(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewSynthesizeExplicitConstructors(compiler AbstractCompiler) *SynthesizeExplicitConstructors {
	this_ := &SynthesizeExplicitConstructors{}
	InitSynthesizeExplicitConstructors(this_, compiler)
	return this_
}

func InitSynthesizeExplicitConstructors(this_ *SynthesizeExplicitConstructors, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_SynthesizeExplicitConstructors()
	this_.compiler = compiler
	this_.astFactory = compiler.CreateAstFactory()
}

// Java method: SynthesizeExplicitConstructors.synthesizeClassConstructorIfMissing(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node classNode)
func (this_ *SynthesizeExplicitConstructors) SynthesizeClassConstructorIfMissing(t *NodeTraversal, classNode pkg_rhino.Node) {
	_ = t
	_ = classNode
	if !(classNode.IsClass()) {
		j2gort.CheckArgument(false)
	}
	if NodeUtil_GetEs6ClassConstructorMemberFunctionDef(classNode) == 0 {
		this_.priv_SynthesizeExplicitConstructors_addSyntheticConstructor_5352(t, classNode)
	}
}

// Java method: SynthesizeExplicitConstructors.addSyntheticConstructor(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node classNode)
func (this_ *SynthesizeExplicitConstructors) priv_SynthesizeExplicitConstructors_addSyntheticConstructor_5352(t *NodeTraversal, classNode pkg_rhino.Node) {
	_ = t
	_ = classNode
	var superClass pkg_rhino.Node = classNode.GetSecondChild()
	_ = superClass
	var classMembers pkg_rhino.Node = classNode.GetLastChild()
	_ = classMembers
	var memberDef pkg_rhino.Node = pkg_rhino.Node(0)
	_ = memberDef
	if superClass.IsEmpty() {
		var function pkg_rhino.Node = this_.astFactory.CreateEmptyFunction(AstFactory_S_Type_1_Node(classNode))
		_ = function
		this_.compiler.ReportChangeToChangeScope(function)
		memberDef = this_.astFactory.CreateMemberFunctionDef("constructor", function)
	} else {
		if !(superClass.IsQualifiedName()) {
			j2gort.CheckStateFmt(false, "Expected Es6NormalizeClasses to make all extends clauses into qualified names, found %s", j2gort.BoxNode(superClass))
		}
		var body pkg_rhino.Node = pkg_rhino.IR_Block_0()
		_ = body
		if !(classNode.IsFromExterns()) && !(this_.priv_SynthesizeExplicitConstructors_isInterface_5353(classNode)) {
			var exprResult pkg_rhino.Node = pkg_rhino.IR_ExprResult(this_.astFactory.CreateConstructorCall(AstFactory_S_Type_1_Node(classNode), this_.astFactory.CreateSuper(AstFactory_S_Type_1_Node(superClass)), []pkg_rhino.Node{pkg_rhino.IR_IterSpread(this_.astFactory.CreateArgumentsReference())}))
			_ = exprResult
			body.AddChildToFront(exprResult)
			NodeUtil_AddFeatureToScript(t.GetCurrentScript(), pkg_parser.StaticGet_FeatureSet_Feature_SUPER(), this_.compiler)
			NodeUtil_AddFeatureToScript(t.GetCurrentScript(), pkg_parser.StaticGet_FeatureSet_Feature_SPREAD_EXPRESSIONS(), this_.compiler)
		}
		var constructor pkg_rhino.Node = this_.astFactory.CreateFunction("", pkg_rhino.IR_ParamList([]pkg_rhino.Node{}), body, AstFactory_S_Type_1_Node(classNode))
		_ = constructor
		memberDef = this_.astFactory.CreateMemberFunctionDef("constructor", constructor)
	}
	memberDef.SrcrefTreeIfMissing(classNode)
	memberDef.MakeNonIndexableRecursive()
	classMembers.AddChildToFront(memberDef)
	NodeUtil_AddFeatureToScript(t.GetCurrentScript(), pkg_parser.StaticGet_FeatureSet_Feature_MEMBER_DECLARATIONS(), this_.compiler)
	this_.compiler.ReportChangeToChangeScope(memberDef.GetOnlyChild())
	this_.compiler.ReportChangeToEnclosingScope(memberDef)
}

// Java method: SynthesizeExplicitConstructors.isInterface(com.google.javascript.rhino.Node classNode)
func (this_ *SynthesizeExplicitConstructors) priv_SynthesizeExplicitConstructors_isInterface_5353(classNode pkg_rhino.Node) bool {
	_ = classNode
	var classJsDocInfo *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(classNode)
	_ = classJsDocInfo
	return ((classJsDocInfo != nil) && classJsDocInfo.IsInterface())
}
