/*
 * Copyright 2018 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/J2clStringValueOfRewriterPass.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.J2clStringValueOfRewriterPass
type J2clStringValueOfRewriterPass struct {
	AbstractPeepholeOptimizationBase
	shouldRunJ2clPasses bool `gson:"shouldRunJ2clPasses"`
}

func (this_ *J2clStringValueOfRewriterPass) BindSelf_(s any) {
	_ = s
	this_.AbstractPeepholeOptimizationBase.BindSelf_(s)
}
func (this_ *J2clStringValueOfRewriterPass) IsJavaNil_() bool { return this_ == nil }
func (*J2clStringValueOfRewriterPass) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.J2clStringValueOfRewriterPass", "J2clStringValueOfRewriterPass"
}
func (this_ *J2clStringValueOfRewriterPass) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.AbstractPeepholeOptimizationBase.IdentityHash_()
}
func (*J2clStringValueOfRewriterPass) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.J2clStringValueOfRewriterPass", "com.google.javascript.jscomp.J2clStringValueOfRewriterPass", func(v_ any) bool {
		return j2gort.InstanceOf[*J2clStringValueOfRewriterPass](v_)
	})
}
func (*J2clStringValueOfRewriterPass) Implements_jscomp_AbstractPeepholeOptimization() {}

func StaticGet_J2clStringValueOfRewriterPass_METHOD_NAME() string {
	return "module$exports$java$lang$String$impl.m_valueOf__java_lang_Object__java_lang_String"
}
func (this_ *J2clStringValueOfRewriterPass) initFields_J2clStringValueOfRewriterPass() {
	_ = this_

}

// Java constructor: J2clStringValueOfRewriterPass()
func NewJ2clStringValueOfRewriterPass() *J2clStringValueOfRewriterPass {
	this_ := &J2clStringValueOfRewriterPass{}
	this_.BindSelf_(this_)
	InitJ2clStringValueOfRewriterPass(this_)
	return this_
}

func InitJ2clStringValueOfRewriterPass(this_ *J2clStringValueOfRewriterPass) {
	InitAbstractPeepholeOptimization(&this_.AbstractPeepholeOptimizationBase)
	this_.initFields_J2clStringValueOfRewriterPass()
}

// Java method: J2clStringValueOfRewriterPass.beginTraversal(com.google.javascript.jscomp.AbstractCompiler compiler)
func (this_ *J2clStringValueOfRewriterPass) BeginTraversal(compiler AbstractCompiler) {
	_ = compiler
	this_.AbstractPeepholeOptimizationBase.BeginTraversal(compiler)
	this_.shouldRunJ2clPasses = J2clSourceFileChecker_ShouldRunJ2clPasses(compiler)
}

// Java method: J2clStringValueOfRewriterPass.optimizeSubtree(com.google.javascript.rhino.Node node)
func (this_ *J2clStringValueOfRewriterPass) OptimizeSubtree(node pkg_rhino.Node) pkg_rhino.Node {
	_ = node
	if !(this_.shouldRunJ2clPasses) || !(J2clStringValueOfRewriterPass_IsStringValueOfCall(node)) {
		return node
	}
	var replacement pkg_rhino.Node = this_.priv_J2clStringValueOfRewriterPass_tryRewriteStringValueOfCall_3494(node)
	_ = replacement
	if replacement != node {
		replacement = replacement.SrcrefIfMissing(node)
		node.ReplaceWith(replacement)
		this_.ReportChangeToEnclosingScope(replacement)
	}
	return replacement
}

// Java method: J2clStringValueOfRewriterPass.tryRewriteStringValueOfCall(com.google.javascript.rhino.Node n)
func (this_ *J2clStringValueOfRewriterPass) priv_J2clStringValueOfRewriterPass_tryRewriteStringValueOfCall_3494(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	var param pkg_rhino.Node = n.GetSecondChild()
	_ = param
	if NodeUtil_IsUndefined(param) {
		return pkg_rhino.IR_String("null")
	} else if NodeUtil_IsDefinedValue(param) && !(param.IsArrayLit()) {
		return NodeUtil_NewCallNode(pkg_rhino.IR_Name("String").Srcref(n), []pkg_rhino.Node{param.Detach()})
	}
	return n
}

// Java static method: J2clStringValueOfRewriterPass.isStringValueOfCall(com.google.javascript.rhino.Node node)
func J2clStringValueOfRewriterPass_IsStringValueOfCall(node pkg_rhino.Node) bool {
	_ = node
	return ((node.IsCall() && node.HasXChildren(int32(2))) && J2clStringValueOfRewriterPass_IsStringValueOfMethodName(node.GetFirstChild()))
}

// Java static method: J2clStringValueOfRewriterPass.isStringValueOfMethodName(com.google.javascript.rhino.Node fnName)
func J2clStringValueOfRewriterPass_IsStringValueOfMethodName(fnName pkg_rhino.Node) bool {
	_ = fnName
	return (fnName.IsQualifiedName() && func() bool {
		sr_1 := fnName.GetOriginalQualifiedName()
		se_2 := "module$exports$java$lang$String$impl.m_valueOf__java_lang_Object__java_lang_String"
		return j2gort.CheckNotNullString(sr_1) == se_2
	}())
}
