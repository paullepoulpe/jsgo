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
// Port of com/google/javascript/jscomp/RewriteCatchWithNoBinding.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.RewriteCatchWithNoBinding
type RewriteCatchWithNoBinding struct {
	AbstractPeepholeTranspilationBase
	compiler AbstractCompiler `gson:"compiler"`
	astFactory *AstFactory `gson:"astFactory"`
	uniqueIdSupplier *UniqueIdSupplier `gson:"uniqueIdSupplier"`
}

func (this_ *RewriteCatchWithNoBinding) BindSelf_(s any) {
	_ = s
}
func (this_ *RewriteCatchWithNoBinding) IsJavaNil_() bool { return this_ == nil }
func (*RewriteCatchWithNoBinding) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.RewriteCatchWithNoBinding", "RewriteCatchWithNoBinding" }
func (this_ *RewriteCatchWithNoBinding) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.AbstractPeepholeTranspilationBase.IdentityHash_()
}
func (*RewriteCatchWithNoBinding) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RewriteCatchWithNoBinding", "com.google.javascript.jscomp.RewriteCatchWithNoBinding", func(v_ any) bool {
		return j2gort.InstanceOf[*RewriteCatchWithNoBinding](v_)
	})
}
func (*RewriteCatchWithNoBinding) Implements_jscomp_AbstractPeepholeTranspilation() {}

func StaticGet_RewriteCatchWithNoBinding_BINDING_NAME() string { return "$jscomp$unused$catch$" }
func (this_ *RewriteCatchWithNoBinding) initFields_RewriteCatchWithNoBinding() {
	_ = this_
}

// Java constructor: RewriteCatchWithNoBinding(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewRewriteCatchWithNoBinding(compiler AbstractCompiler) *RewriteCatchWithNoBinding {
	this_ := &RewriteCatchWithNoBinding{}
	InitRewriteCatchWithNoBinding(this_, compiler)
	return this_
}

func InitRewriteCatchWithNoBinding(this_ *RewriteCatchWithNoBinding, compiler AbstractCompiler) {
	_ = compiler
	InitAbstractPeepholeTranspilation(&this_.AbstractPeepholeTranspilationBase)
	this_.initFields_RewriteCatchWithNoBinding()
	this_.compiler = compiler
	this_.astFactory = compiler.CreateAstFactory()
	this_.uniqueIdSupplier = compiler.GetUniqueIdSupplier()
}

// Java method: RewriteCatchWithNoBinding.getTranspiledAwayFeatures()
func (this_ *RewriteCatchWithNoBinding) GetTranspiledAwayFeatures() *pkg_parser.FeatureSet {
	return pkg_parser.StaticGet_FeatureSet_BARE_MINIMUM().With_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_OPTIONAL_CATCH_BINDING())
}

// Java method: RewriteCatchWithNoBinding.transpileSubtree(com.google.javascript.rhino.Node n)
func (this_ *RewriteCatchWithNoBinding) TranspileSubtree(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	if ((!(n.IsCatch()) || !(n.GetFirstChild().IsEmpty()))) {
		return n
	}
	var name pkg_rhino.Node = this_.astFactory.CreateNameWithUnknownType(("$jscomp$unused$catch$" + j2gort.StrValString(this_.uniqueIdSupplier.GetUniqueId(this_.compiler.GetInput(NodeUtil_GetInputId(n))))))
	_ = name
	n.GetFirstChild().ReplaceWith(name.SrcrefTree(n.GetFirstChild()))
	this_.compiler.ReportChangeToEnclosingScope(name)
	return n
}

