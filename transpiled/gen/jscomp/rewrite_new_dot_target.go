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
// Port of com/google/javascript/jscomp/RewriteNewDotTarget.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.RewriteNewDotTarget
type RewriteNewDotTarget struct {
	AbstractPeepholeTranspilationBase
	compiler AbstractCompiler `gson:"compiler"`
	astFactory *AstFactory `gson:"astFactory"`
}

func (this_ *RewriteNewDotTarget) BindSelf_(s any) {
	_ = s
}
func (this_ *RewriteNewDotTarget) IsJavaNil_() bool { return this_ == nil }
func (*RewriteNewDotTarget) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.RewriteNewDotTarget", "RewriteNewDotTarget" }
func (this_ *RewriteNewDotTarget) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.AbstractPeepholeTranspilationBase.IdentityHash_()
}
func (*RewriteNewDotTarget) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RewriteNewDotTarget", "com.google.javascript.jscomp.RewriteNewDotTarget", func(v_ any) bool {
		return j2gort.InstanceOf[*RewriteNewDotTarget](v_)
	})
}
func (*RewriteNewDotTarget) Implements_jscomp_AbstractPeepholeTranspilation() {}

func (this_ *RewriteNewDotTarget) initFields_RewriteNewDotTarget() {
	_ = this_
}

// Java constructor: RewriteNewDotTarget(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewRewriteNewDotTarget(compiler AbstractCompiler) *RewriteNewDotTarget {
	this_ := &RewriteNewDotTarget{}
	InitRewriteNewDotTarget(this_, compiler)
	return this_
}

func InitRewriteNewDotTarget(this_ *RewriteNewDotTarget, compiler AbstractCompiler) {
	_ = compiler
	InitAbstractPeepholeTranspilation(&this_.AbstractPeepholeTranspilationBase)
	this_.initFields_RewriteNewDotTarget()
	this_.compiler = compiler
	this_.astFactory = compiler.CreateAstFactory()
}

// Java method: RewriteNewDotTarget.getTranspiledAwayFeatures()
func (this_ *RewriteNewDotTarget) GetTranspiledAwayFeatures() *pkg_parser.FeatureSet {
	return pkg_parser.StaticGet_FeatureSet_BARE_MINIMUM().With_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_NEW_TARGET())
}

// Java method: RewriteNewDotTarget.transpileSubtree(com.google.javascript.rhino.Node n)
func (this_ *RewriteNewDotTarget) TranspileSubtree(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	if ((n.RawToken_() != uint8(135))) {
		return n
	}
	var enclosingNonArrowFunction pkg_rhino.Node = NodeUtil_GetEnclosingNonArrowFunction(n)
	_ = enclosingNonArrowFunction
	if (((enclosingNonArrowFunction != 0) && NodeUtil_IsEs6Constructor(enclosingNonArrowFunction))) {
		var enclosingClass pkg_rhino.Node = enclosingNonArrowFunction.GetParent().GetGrandparent()
		_ = enclosingClass
		var replacement pkg_rhino.Node = this_.astFactory.CreateGetProp(this_.astFactory.CreateThisForEs6Class(enclosingClass), "constructor", AstFactory_S_Type_1_Node(n)).SrcrefTree(n)
		_ = replacement
		n.ReplaceWith(replacement)
		this_.compiler.ReportChangeToEnclosingScope(replacement)
		return replacement
	} else {
		TranspilationUtil_CannotConvertYet(this_.compiler, n, "new.target")
	}
	return n
}

